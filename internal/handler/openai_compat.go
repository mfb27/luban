package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/auth"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/provider"
	"github.com/mfb27/luban/internal/providerfactory"
	"github.com/mfb27/luban/internal/response"
	"gorm.io/gorm"
)

// OpenAICompat handles OpenAI-compatible API endpoints
type OpenAICompat struct {
	db *gorm.DB
}

// NewOpenAICompat creates a new OpenAI-compatible handler
func NewOpenAICompat(db *gorm.DB) *OpenAICompat {
	return &OpenAICompat{db: db}
}

// OpenAI model response format
type openAIModelResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type openAIModelsListResponse struct {
	Data []openAIModelResponse `json:"data"`
}

// OpenAI chat request format
type openAIChatRequest struct {
	Model    string          `json:"model"`
	Messages []provider.Message `json:"messages"`
	Stream   bool            `json:"stream"`
	SessionID string         `json:"session_id,omitempty"`
}

// OpenAI chat response format
type openAIChatResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []openAIChatChoice `json:"choices"`
	Usage   *openAIUsage      `json:"usage,omitempty"`
}

type openAIChatChoice struct {
	Index        int             `json:"index"`
	Message      *provider.Message `json:"message,omitempty"`
	Delta        *provider.Message `json:"delta,omitempty"`
	FinishReason string          `json:"finish_reason"`
}

type openAIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// OpenAI streaming chunk format
type openAIStreamChunk struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []openAIChatChoice `json:"choices"`
}

// listModels returns the list of available models in OpenAI format
func (o *OpenAICompat) listModels(c *gin.Context) {
	var models []model.Model
	if err := o.db.Where("status = ?", "active").Find(&models).Error; err != nil {
		response.NewResponseHelper(c).Error(response.CodeDatabaseError, "failed to load models")
		return
	}

	if models == nil {
		models = []model.Model{}
	}

	data := make([]openAIModelResponse, len(models))
	for i, m := range models {
		data[i] = openAIModelResponse{
			ID:      m.ModelID,
			Object:  "model",
			Created: m.CreatedAt.Unix(),
			OwnedBy: m.ProviderID,
		}
	}

	c.JSON(200, openAIModelsListResponse{Data: data})
}

// chatCompletions handles chat completion requests in OpenAI format
func (o *OpenAICompat) chatCompletions(c *gin.Context) {
	var req openAIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"message": "invalid request", "type": "invalid_request_error"}})
		return
	}

	if req.Model == "" {
		c.JSON(400, gin.H{"error": gin.H{"message": "model is required", "type": "invalid_request_error"}})
		return
	}

	if len(req.Messages) == 0 {
		c.JSON(400, gin.H{"error": gin.H{"message": "messages is required", "type": "invalid_request_error"}})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := c.GetString("user_id")

	// Get last message content (the current user input)
	lastMsg := req.Messages[len(req.Messages)-1]
	if lastMsg.Role != "user" {
		c.JSON(400, gin.H{"error": gin.H{"message": "last message must be from user", "type": "invalid_request_error"}})
		return
	}

	// Get model info to determine provider
	var modelInfo model.Model
	if err := o.db.Where("model_id = ? AND status = ?", req.Model, "active").First(&modelInfo).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			c.JSON(404, gin.H{"error": gin.H{"message": "model not found", "type": "invalid_request_error"}})
			return
		}
		c.JSON(500, gin.H{"error": gin.H{"message": "database error", "type": "server_error"}})
		return
	}

	// Get provider configuration by provider_id
	var providerConfig model.APIProvider
	if err := o.db.Where("id = ? AND status = ?", modelInfo.ProviderID, "active").First(&providerConfig).Error; err != nil {
		if strings.Contains(err.Error(), "record not found") {
			c.JSON(404, gin.H{"error": gin.H{"message": "provider not found or inactive", "type": "invalid_request_error"}})
			return
		}
		c.JSON(500, gin.H{"error": gin.H{"message": "database error", "type": "server_error"}})
		return
	}

	if providerConfig.GetBaseURL() == "" {
		c.JSON(400, gin.H{"error": gin.H{"message": "provider has no base URL configured", "type": "invalid_request_error"}})
		return
	}

	// Create provider client
	prov, err := providerfactory.NewProvider(providerConfig.ProviderType, &providerfactory.ProviderConfig{
		APIKey:  providerConfig.APIKey,
		BaseURL: providerConfig.GetBaseURL(),
	})
	if err != nil {
		c.JSON(500, gin.H{"error": gin.H{"message": fmt.Sprintf("failed to create provider: %v", err), "type": "server_error"}})
		return
	}

	// Track daily chat limit for authenticated users
	if userID != "" {
		now := time.Now()
		var user model.User
		if err := o.db.First(&user, "id = ?", userID).Error; err == nil {
			// Simplified limit check - use same pattern as chat.go
			if user.DailyChatLimit != -1 {
				today := now.Format("2006-01-02")
				var count model.UserDailyChatCount
				if err := o.db.Where("user_id = ? AND date = ?", userID, today).First(&count).Error; err == nil {
					if count.Count >= user.DailyChatLimit {
						c.JSON(429, gin.H{"error": gin.H{"message": fmt.Sprintf("daily chat limit exceeded. Limit: %d", user.DailyChatLimit), "type": "rate_limit_error"}})
						return
					}
					count.Count++
					o.db.Save(&count)
				} else if errors.Is(err, gorm.ErrRecordNotFound) {
					o.db.Create(&model.UserDailyChatCount{
						ID:        uuid.NewString(),
						UserID:    userID,
						Date:      today,
						Count:     1,
						CreatedAt: now,
						UpdatedAt: now,
					})
				}
			}
		}
	}

	if req.Stream {
		o.handleStreamChat(c, prov, req, providerConfig.ProviderType)
	} else {
		o.handleNonStreamChat(c, prov, req)
	}
}

// handleNonStreamChat handles non-streaming chat completion
func (o *OpenAICompat) handleNonStreamChat(c *gin.Context, prov provider.Provider, req openAIChatRequest) {
	resp, err := prov.Chat(context.Background(), req.Model, req.Messages)
	if err != nil {
		c.JSON(500, gin.H{"error": gin.H{"message": err.Error(), "type": "server_error"}})
		return
	}

	choices := make([]openAIChatChoice, len(resp.Choices))
	for i, ch := range resp.Choices {
		choices[i] = openAIChatChoice{
			Index:        ch.Index,
			Message:      ch.Message,
			FinishReason: ch.FinishReason,
		}
	}

	usage := &openAIUsage{
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}

	c.JSON(200, openAIChatResponse{
		ID:      resp.ID,
		Object:  "chat.completion",
		Created: resp.Created,
		Model:   resp.Model,
		Choices: choices,
		Usage:   usage,
	})
}

// handleStreamChat handles streaming chat completion
func (o *OpenAICompat) handleStreamChat(c *gin.Context, prov provider.Provider, req openAIChatRequest, providerType string) {
	// Get user ID from context (set by auth middleware)
	userID := c.GetString("user_id")

	now := time.Now()

	// Create or find session
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.NewString()
		s := model.Session{
			ID:        sessionID,
			UserID:    userID,
			Title:     titleFromContent(req.Messages[len(req.Messages)-1].Content),
			ModelID:   req.Model,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := o.db.Create(&s).Error; err != nil {
			c.JSON(500, gin.H{"error": gin.H{"message": "failed to create session", "type": "server_error"}})
			return
		}
	} else {
		// Check if session exists
		var session model.Session
		if err := o.db.First(&session, "id = ?", sessionID).Error; err != nil {
			c.JSON(404, gin.H{"error": gin.H{"message": "session not found", "type": "not_found_error"}})
			return
		}

		// Update session timestamp
		_ = o.db.Model(&model.Session{}).Where("id = ?", sessionID).Update("updated_at", now).Error
	}

	// Save user message (last message from user)
	lastMsg := req.Messages[len(req.Messages)-1]
	userMsg := model.Message{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "user",
		Content:   lastMsg.Content,
		CreatedAt: now,
	}
	if err := o.db.Create(&userMsg).Error; err != nil {
		// Log error but continue
		fmt.Printf("failed to save user message: %v\n", err)
	}

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Stream(func(w io.Writer) bool {
		return false
	})

	chatID := "chatcmpl-" + uuid.NewString()[:8]
	created := time.Now().Unix()

	var fullContent strings.Builder

	err := prov.ChatStream(context.Background(), req.Model, req.Messages, func(chunk *provider.StreamChunk) error {
		if len(chunk.Choices) > 0 {
			choices := make([]openAIChatChoice, len(chunk.Choices))
			for i, ch := range chunk.Choices {
				choices[i] = openAIChatChoice{
					Index:        ch.Index,
					Delta:        ch.Delta,
					FinishReason: ch.FinishReason,
				}

				// Collect content for saving
				if ch.Delta != nil && ch.Delta.Content != "" {
					fullContent.WriteString(ch.Delta.Content)
				}
			}

			streamChunk := openAIStreamChunk{
				ID:      chatID,
				Object:  "chat.completion.chunk",
				Created: created,
				Model:   req.Model,
				Choices: choices,
			}

			data, err := json.Marshal(streamChunk)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			if err != nil {
				return err
			}
			c.Writer.Flush()
		}
		return nil
	})

	if err != nil {
		// Send error chunk
		errorChunk := gin.H{
			"error": gin.H{
				"message": err.Error(),
				"type":    "server_error",
			},
		}
		data, _ := json.Marshal(errorChunk)
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		c.Writer.Flush()
		return
	}

	// Save assistant message
	reply := fullContent.String()
	if reply == "" {
		reply = "抱歉，我没有收到有效的回复。"
	}

	asstMsg := model.Message{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   reply,
		CreatedAt: time.Now(),
	}
	if err := o.db.Create(&asstMsg).Error; err != nil {
		fmt.Printf("failed to save assistant message: %v\n", err)
	}

	// Send done marker with session_id
	doneData := map[string]interface{}{
		"session_id": sessionID,
	}
	doneJSON, _ := json.Marshal(doneData)
	fmt.Fprintf(c.Writer, "data: [DONE] %s\n\n", doneJSON)
	c.Writer.Flush()
}

// RegisterOpenAICompatRoutes registers OpenAI-compatible API routes
func (a *App) registerOpenAICompatRoutes() {
	openaiCompat := NewOpenAICompat(a.db)

	// Public route: list models (no authentication required, per OpenAI spec)
	a.Engine.GET("/v1/models", openaiCompat.listModels)

	// Protected routes: chat completions require authentication
	v1 := a.Engine.Group("/v1")
	v1.Use(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": gin.H{"message": "missing authorization header", "type": "invalid_request_error"}})
			c.Abort()
			return
		}

		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims, err := auth.ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": gin.H{"message": "invalid token", "type": "invalid_request_error"}})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	})
	{
		v1.POST("/chat/completions", openaiCompat.chatCompletions)
	}
}