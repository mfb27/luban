package anthropic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mfb27/luban/internal/provider"
)

const (
	defaultBaseURL = "https://api.anthropic.com/v1/messages"
	apiVersion     = "2023-06-01"
)

// Client represents an Anthropic API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Anthropic client
func NewClient(apiKey, baseURL string) *Client {
	url := defaultBaseURL
	if baseURL != "" {
		// Ensure the URL ends with /messages
		if !strings.HasSuffix(baseURL, "/messages") {
			if strings.HasSuffix(baseURL, "/v1") {
				url = baseURL + "/messages"
			} else if strings.HasSuffix(baseURL, "/") {
				url = baseURL + "v1/messages"
			} else {
				url = baseURL + "/v1/messages"
			}
		} else {
			url = baseURL
		}
	}
	return &Client{
		apiKey:     apiKey,
		baseURL:    url,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

// Name returns the provider name
func (c *Client) Name() string {
	return "anthropic"
}

// anthropicRequest represents Anthropic API request format
type anthropicRequest struct {
	Model     string                  `json:"model"`
	MaxTokens int                     `json:"max_tokens"` // Anthropic requires max_tokens
	Messages  []anthropicMessage      `json:"messages"`
	Stream    bool                    `json:"stream"`
}

// anthropicMessage represents Anthropic message format
type anthropicMessage struct {
	Role    string `json:"role"` // user/assistant
	Content string `json:"content"`
}

// anthropicResponse represents Anthropic API response format
type anthropicResponse struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Role       string                   `json:"role"`
	Model      string                   `json:"model"`
	Content    []anthropicContentBlock  `json:"content"`
	StopReason string                   `json:"stop_reason"`
	Usage      anthropicUsage           `json:"usage"`
	Error      *provider.ErrorResp      `json:"error,omitempty"`
}

// anthropicContentBlock represents a content block in Anthropic response
type anthropicContentBlock struct {
	Type string `json:"type"` // text
	Text string `json:"text"`
}

// anthropicUsage represents Anthropic usage information
type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// anthropicStreamChunk represents Anthropic streaming chunk format
type anthropicStreamChunk struct {
	Type         string                   `json:"type"`
	Index        int                      `json:"index,omitempty"`
	Delta        *anthropicStreamDelta    `json:"delta,omitempty"`
	Message      *anthropicStreamMessage  `json:"message,omitempty"`
	ContentBlock *anthropicContentBlock   `json:"content_block,omitempty"`
	Usage        *anthropicUsage          `json:"usage,omitempty"`
	Error        *provider.ErrorResp      `json:"error,omitempty"`
}

// anthropicStreamDelta represents delta in streaming response
type anthropicStreamDelta struct {
	Type string `json:"type"` // text_delta
	Text string `json:"text"`
}

// anthropicStreamMessage represents message start in streaming
type anthropicStreamMessage struct {
	ID    string         `json:"id"`
	Type  string         `json:"type"`
	Role  string         `json:"role"`
	Model string         `json:"model"`
	Usage anthropicUsage `json:"usage"`
}

// Chat sends a non-streaming chat completion request
func (c *Client) Chat(ctx context.Context, model string, messages []provider.Message) (*provider.ChatResponse, error) {
	// Convert messages to Anthropic format
	anthropicMessages := make([]anthropicMessage, len(messages))
	for i, msg := range messages {
		anthropicMessages[i] = anthropicMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	body, err := json.Marshal(&anthropicRequest{
		Model:     model,
		MaxTokens: 4096, // Default max tokens
		Messages:  anthropicMessages,
		Stream:    false,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var anthropicResp anthropicResponse
	if err := json.Unmarshal(respBody, &anthropicResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if anthropicResp.Error != nil {
		return nil, fmt.Errorf("api error: %s", anthropicResp.Error.Message)
	}

	// Extract text content
	var content string
	for _, block := range anthropicResp.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	// Convert to provider response
	return &provider.ChatResponse{
		ID:      anthropicResp.ID,
		Created: time.Now().Unix(),
		Model:   anthropicResp.Model,
		Choices: []provider.ChatChoice{
			{
				Index: 0,
				Message: &provider.Message{
					Role:    anthropicResp.Role,
					Content: content,
				},
				FinishReason: anthropicResp.StopReason,
			},
		},
		Usage: provider.ChatUsage{
			PromptTokens:     anthropicResp.Usage.InputTokens,
			CompletionTokens: anthropicResp.Usage.OutputTokens,
			TotalTokens:      anthropicResp.Usage.InputTokens + anthropicResp.Usage.OutputTokens,
		},
	}, nil
}

// ChatStream sends a streaming chat completion request
func (c *Client) ChatStream(ctx context.Context, model string, messages []provider.Message, callback provider.StreamCallback) error {
	// Convert messages to Anthropic format
	anthropicMessages := make([]anthropicMessage, len(messages))
	for i, msg := range messages {
		anthropicMessages[i] = anthropicMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	body, err := json.Marshal(&anthropicRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages:  anthropicMessages,
		Stream:    true,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", apiVersion)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var messageID string
	var modelName string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse SSE format: "data: {...}"
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimPrefix(line, "data:")
		data = strings.TrimSpace(data)

		var chunk anthropicStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // Skip invalid chunks
		}

		if chunk.Error != nil {
			return fmt.Errorf("api error: %s", chunk.Error.Message)
		}

		// Handle different event types
		switch chunk.Type {
		case "message_start":
			if chunk.Message != nil {
				messageID = chunk.Message.ID
				modelName = chunk.Message.Model
			}
		case "content_block_delta":
			if chunk.Delta != nil && chunk.Delta.Type == "text_delta" {
				if err := callback(&provider.StreamChunk{
					ID:      messageID,
					Created: time.Now().Unix(),
					Model:   modelName,
					Choices: []provider.ChatChoice{
						{
							Index: chunk.Index,
							Delta: &provider.Message{
								Content: chunk.Delta.Text,
							},
						},
					},
				}); err != nil {
					return err
				}
			}
		case "message_delta":
			// Message update, can contain stop_reason
		case "message_stop":
			// Stream ended - use return to exit the function
			return scanner.Err()
		}
	}

	return scanner.Err()
}