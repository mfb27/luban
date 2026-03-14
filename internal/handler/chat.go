package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mfb27/luban/internal/auth"
	"github.com/mfb27/luban/internal/middleware"
	"github.com/mfb27/luban/internal/model"
	"github.com/mfb27/luban/internal/zhipu"
	"go.uber.org/zap"
)

type chatReq struct {
	SessionID      string   `json:"session_id"`
	Content        string   `json:"content"`
	ModelID        string   `json:"model_id"`
	AttachmentURLs []string `json:"attachment_urls"`
}

// SSE event types
const (
	sseEventTypeSession = "session"
	sseEventTypeContent = "content"
	sseEventTypeDone    = "done"
	sseEventTypeError   = "error"
)

type sseEvent struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
	// For session event
	SessionID string `json:"session_id,omitempty"`
	UserMsgID string `json:"user_msg_id,omitempty"`
	// For error event
	Error string `json:"error,omitempty"`
}

func (a *App) postChat(c *gin.Context) {
	// Get request-aware logger
	reqLog := middleware.GetLoggerWithRequestID(c)
	if reqLog != nil {
		reqLog.Info("postChat started",
			middleware.GetRequestIDField(c),
		)
	}

	var req chatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		if reqLog != nil {
			reqLog.Error("failed to bind JSON",
				middleware.GetRequestIDField(c),
				middleware.GetErrorField(err),
			)
		}
		NewResponseHelper(c).Error(CodeInvalidParam, "invalid json")
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		NewResponseHelper(c).Error(CodeRequiredParam, "content required")
		return
	}
	if req.ModelID == "" {
		req.ModelID = "glm-4-flash"
	}

	now := time.Now()

	// Get user context - check if Authorization header is present
	userID := ""
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if claims, err := auth.ValidateToken(tokenString); err == nil {
			userID = claims.UserID
		}
	}

	// Ensure session exists.
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.NewString()
		s := model.Session{
			ID:        sessionID,
			UserID:    userID,
			Title:     titleFromContent(req.Content),
			ModelID:   req.ModelID,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := a.db.Create(&s).Error; err != nil {
			NewResponseHelper(c).Error(CodeInternal, err.Error())
			return
		}
	} else {
		// Check if session exists
		var session model.Session
		if err := a.db.First(&session, "id = ?", sessionID).Error; err != nil {
			NewResponseHelper(c).Error(CodeNotFound, "session not found")
			return
		}

		// Verify session ownership based on authentication state
		if userID != "" {
			// User is authenticated
			if session.UserID != "" && session.UserID != userID {
				// Session belongs to a different user
				NewResponseHelper(c).Error(CodeForbidden, "unauthorized to access this session")
				return
			} else if session.UserID == "" {
				// Convert anonymous session to authenticated session
				if err := a.db.Model(&session).Update("user_id", userID).Error; err != nil {
					NewResponseHelper(c).Error(CodeInternal, "failed to convert session to authenticated")
					return
				}
			}
		} else {
			// User is not authenticated
			if session.UserID != "" {
				// Trying to access an authenticated session anonymously
				NewResponseHelper(c).Error(CodeForbidden, "authentication required to access this session")
				return
			}
		}

		// Update session timestamp and convert to authenticated if needed
		updates := map[string]interface{}{"updated_at": now}
		if userID != "" && session.UserID == "" {
			// Convert anonymous session to authenticated session
			updates["user_id"] = userID
		}
		_ = a.db.Model(&model.Session{}).Where("id = ?", sessionID).Updates(updates).Error
	}

	userMsg := model.Message{
		ID:        uuid.NewString(),
		SessionID: sessionID,
		UserID:    userID,
		Role:      "user",
		Content:   req.Content,
		CreatedAt: now,
	}
	if err := a.db.Create(&userMsg).Error; err != nil {
		NewResponseHelper(c).Error(CodeDatabaseError, err.Error())
		return
	}

	// Check if this is a new conversation (no existing messages in the session)
	var messageCount int64
	isNewConversation := false

	// Count messages in the session
	if err := a.db.Table("messages").Where("session_id = ?", sessionID).Count(&messageCount).Error; err != nil {
		if reqLog != nil {
			reqLog.Error("failed to check message count",
				middleware.GetRequestIDField(c),
				middleware.GetErrorField(err),
			)
		}
		NewResponseHelper(c).Error(CodeDatabaseError, err.Error())
		return
	}

	// For new conversations, get conversation history for context
	var historyMessages []model.Message
	historyLimit := 20 // Limit history to prevent token overflow
	isNewConversation = messageCount == 0

	if isNewConversation {
		// Get history from this session
		if err := a.db.Table("messages").
			Where("session_id = ? AND id != ?", sessionID, userMsg.ID).
			Order("created_at ASC").
			Limit(historyLimit).
			Find(&historyMessages).Error; err != nil {
			if reqLog != nil {
				reqLog.Warn("failed to load history, continuing without it",
					middleware.GetRequestIDField(c),
					middleware.GetErrorField(err),
				)
			}
		}
	}

	// Build messages for API request
	messages := []zhipu.Message{}
	for _, msg := range historyMessages {
		messages = append(messages, zhipu.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	// Add current user message
	messages = append(messages, zhipu.Message{
		Role:    "user",
		Content: req.Content,
	})

	// Set SSE headers
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Stream(func(w io.Writer) bool {
		return false
	})

	// Helper function to send SSE event
	sendEvent := func(eventType string, data interface{}) error {
		var jsonStr string
		if str, ok := data.(string); ok {
			jsonStr = fmt.Sprintf(`{"type":"%s","data":%q}`, eventType, str)
		} else {
			jsonBytes, err := json.Marshal(data)
			if err != nil {
				return err
			}
			var obj map[string]interface{}
			_ = json.Unmarshal(jsonBytes, &obj)
			obj["type"] = eventType
			newBytes, _ := json.Marshal(obj)
			jsonStr = string(newBytes)
		}
		_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", jsonStr)
		if err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}

	// Send initial session info
	sessionEvent := map[string]interface{}{
		"type":       sseEventTypeSession,
		"session_id": sessionID,
		"user_msg_id": userMsg.ID,
	}
	sessionJson, _ := json.Marshal(sessionEvent)
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(sessionJson))
	c.Writer.Flush()

	// Prepare assistant message
	asstMsgID := uuid.NewString()
	var fullContent strings.Builder

	// Call ZhipuAI streaming API
	err := a.zhipu.ChatStream(&zhipu.ChatRequest{
		Model:    req.ModelID,
		Messages: messages,
		Stream:   true,
	}, func(chunk *zhipu.StreamChunk) error {
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil {
			delta := chunk.Choices[0].Delta
			if delta.Content != "" {
				fullContent.WriteString(delta.Content)
				// Send content chunk
				return sendEvent(sseEventTypeContent, delta.Content)
			}
		}
		return nil
	})

	if err != nil {
		if reqLog != nil {
			reqLog.Error("zhipuai streaming api call failed",
				middleware.GetRequestIDField(c),
				middleware.GetErrorField(err),
			)
		}
		sendEvent(sseEventTypeError, err.Error())
		return
	}

	// Save assistant message to database
	reply := fullContent.String()
	if reply == "" {
		reply = "抱歉，我没有收到有效的回复。"
	}

	asstMsg := model.Message{
		ID:        asstMsgID,
		SessionID: sessionID,
		UserID:    userID,
		Role:      "assistant",
		Content:   reply,
		CreatedAt: time.Now(),
	}
	if err := a.db.Create(&asstMsg).Error; err != nil {
		if reqLog != nil {
			reqLog.Error("failed to save assistant message",
				middleware.GetRequestIDField(c),
				middleware.GetErrorField(err),
			)
		}
	}

	// Send done event with assistant message ID and new conversation flag
	doneEvent := map[string]interface{}{
		"type":             sseEventTypeDone,
		"assistant_id":     asstMsgID,
		"session_id":       sessionID,
		"is_new_conversation": isNewConversation,
	}
	doneJson, _ := json.Marshal(doneEvent)
	fmt.Fprintf(c.Writer, "data: %s\n\n", string(doneJson))
	c.Writer.Flush()

	if reqLog != nil {
		reqLog.Info("postChat completed",
			middleware.GetRequestIDField(c),
			zap.String("session_id", sessionID),
		)
	}
}

func titleFromContent(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "新对话"
	}
	r := []rune(s)
	if len(r) > 20 {
		return string(r[:20]) + "…"
	}
	return s
}
