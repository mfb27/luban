package openai

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
	defaultBaseURL = "https://api.openai.com/v1/chat/completions"
)

// Client represents an OpenAI API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new OpenAI client
func NewClient(apiKey, baseURL string) *Client {
	url := defaultBaseURL
	if baseURL != "" {
		// Ensure the URL ends with /chat/completions
		if !strings.HasSuffix(baseURL, "/chat/completions") {
			if strings.HasSuffix(baseURL, "/v1") {
				url = baseURL + "/chat/completions"
			} else if strings.HasSuffix(baseURL, "/") {
				url = baseURL + "v1/chat/completions"
			} else {
				url = baseURL + "/v1/chat/completions"
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
	return "openai"
}

// openAIRequest represents OpenAI API request format
type openAIRequest struct {
	Model    string            `json:"model"`
	Messages []provider.Message `json:"messages"`
	Stream   bool              `json:"stream"`
}

// openAIResponse represents OpenAI API response format
type openAIResponse struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Choices []openAIChoice      `json:"choices"`
	Usage   provider.ChatUsage  `json:"usage,omitempty"`
	Error   *provider.ErrorResp `json:"error,omitempty"`
}

// openAIChoice represents a choice in OpenAI response
type openAIChoice struct {
	Index        int               `json:"index"`
	Message      *provider.Message `json:"message,omitempty"`
	Delta        *provider.Message `json:"delta,omitempty"`
	FinishReason string            `json:"finish_reason"`
}

// openAIStreamChunk represents OpenAI streaming chunk format
type openAIStreamChunk struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Choices []openAIChoice      `json:"choices"`
	Error   *provider.ErrorResp `json:"error,omitempty"`
}

// Chat sends a non-streaming chat completion request
func (c *Client) Chat(ctx context.Context, model string, messages []provider.Message) (*provider.ChatResponse, error) {
	body, err := json.Marshal(&openAIRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

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

	var openAIResp openAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("api error: %s", openAIResp.Error.Message)
	}

	// Convert to provider response
	choices := make([]provider.ChatChoice, len(openAIResp.Choices))
	for i, ch := range openAIResp.Choices {
		choices[i] = provider.ChatChoice{
			Index:        ch.Index,
			Message:      ch.Message,
			FinishReason: ch.FinishReason,
		}
	}

	return &provider.ChatResponse{
		ID:      openAIResp.ID,
		Created: openAIResp.Created,
		Model:   openAIResp.Model,
		Choices: choices,
		Usage:   openAIResp.Usage,
	}, nil
}

// ChatStream sends a streaming chat completion request
func (c *Client) ChatStream(ctx context.Context, model string, messages []provider.Message, callback provider.StreamCallback) error {
	body, err := json.Marshal(&openAIRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

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

		var chunk openAIStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue // Skip invalid chunks
		}

		if chunk.Error != nil {
			return fmt.Errorf("api error: %s", chunk.Error.Message)
		}

		// Convert to provider chunk
		choices := make([]provider.ChatChoice, len(chunk.Choices))
		for i, ch := range chunk.Choices {
			choices[i] = provider.ChatChoice{
				Index:        ch.Index,
				Delta:        ch.Delta,
				FinishReason: ch.FinishReason,
			}
		}

		if err := callback(&provider.StreamChunk{
			ID:      chunk.ID,
			Created: chunk.Created,
			Model:   chunk.Model,
			Choices: choices,
		}); err != nil {
			return err
		}
	}

	return scanner.Err()
}