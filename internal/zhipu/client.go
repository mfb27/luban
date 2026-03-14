package zhipu

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
)

// Config holds the ZhipuAI configuration
type Config struct {
	APIKey  string
	BaseURL string
}

// Client represents a ZhipuAI API client
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new ZhipuAI client
func NewClient(cfg *Config) *Client {
	baseURL := defaultBaseURL
	if cfg.BaseURL != "" {
		baseURL = cfg.BaseURL
	}
	return &Client{
		config: &Config{
			APIKey:  cfg.APIKey,
			BaseURL: baseURL,
		},
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatChoice represents a choice in the response
type ChatChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Delta        *Message `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason"`
}

// ChatUsage represents token usage information
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatResponse represents the chat completion response
type ChatResponse struct {
	ID      string        `json:"id"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChatChoice  `json:"choices"`
	Usage   ChatUsage     `json:"usage"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// StreamChunk represents a streaming response chunk
type StreamChunk struct {
	ID      string       `json:"id"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

// ErrorResponse represents an API error
type ErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// generateToken creates a JWT token for API authentication
// Following ZhipuAI's token generation specification
func generateToken(apiKey string) (string, error) {
	parts := bytes.Split([]byte(apiKey), []byte("."))
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid api key format")
	}

	id := string(parts[0])
	secret := string(parts[1])

	now := time.Now()
	exp := now.Add(1 * time.Hour)

	// Header
	header := map[string]interface{}{
		"alg": "HS256",
		"sign_type": "SIGN",
	}

	// Payload
	payload := map[string]interface{}{
		"api_key": id,
		"exp":     exp.Unix(),
		"timestamp": now.Unix(),
	}

	// Encode header
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encode payload
	payloadJSON, _ := json.Marshal(payload)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create signature
	signingInput := headerB64 + "." + payloadB64
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	// Combine to form JWT
	token := signingInput + "." + signature
	return token, nil
}

// Chat sends a chat completion request
func (c *Client) Chat(req *ChatRequest) (*ChatResponse, error) {
	token, err := generateToken(c.config.APIKey)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.config.BaseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

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

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("api error: %s", chatResp.Error.Message)
	}

	return &chatResp, nil
}

// StreamCallback is called for each streaming chunk
type StreamCallback func(chunk *StreamChunk) error

// ChatStream sends a streaming chat completion request
func (c *Client) ChatStream(req *ChatRequest, callback StreamCallback) error {
	token, err := generateToken(c.config.APIKey)
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.config.BaseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(respBody))
	}

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

		// Check for [DONE] marker
		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // Skip invalid chunks
		}

		if chunk.Error != nil {
			return fmt.Errorf("api error: %s", chunk.Error.Message)
		}

		if err := callback(&chunk); err != nil {
			return err
		}
	}

	return scanner.Err()
}
