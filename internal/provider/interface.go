package provider

import "context"

// Provider defines the interface for AI providers
type Provider interface {
	// Chat sends a non-streaming chat completion request
	Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error)

	// ChatStream sends a streaming chat completion request
	ChatStream(ctx context.Context, model string, messages []Message, callback StreamCallback) error

	// Name returns the provider name
	Name() string
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`    // user/assistant/system
	Content string `json:"content"`
}

// ChatChoice represents a choice in the response
type ChatChoice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Message `json:"delta,omitempty"`
	FinishReason string   `json:"finish_reason"`
}

// ChatUsage represents token usage information
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatResponse represents the chat completion response
type ChatResponse struct {
	ID      string       `json:"id"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Usage   ChatUsage    `json:"usage,omitempty"`
	Error   *ErrorResp   `json:"error,omitempty"`
}

// StreamChunk represents a streaming response chunk
type StreamChunk struct {
	ID      string       `json:"id"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []ChatChoice `json:"choices"`
	Error   *ErrorResp   `json:"error,omitempty"`
}

// ErrorResp represents an API error
type ErrorResp struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// StreamCallback is called for each streaming chunk
type StreamCallback func(chunk *StreamChunk) error