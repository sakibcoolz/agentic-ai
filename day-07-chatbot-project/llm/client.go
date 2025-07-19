package llm

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Client wraps the OpenAI client with additional functionality
type Client struct {
	client *openai.Client
	model  string
}

// NewClient creates a new LLM client
func NewClient(apiKey, model string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if model == "" {
		model = openai.GPT3Dot5Turbo
	}

	client := openai.NewClient(apiKey)

	return &Client{
		client: client,
		model:  model,
	}, nil
}

// ChatCompletion sends a chat completion request to OpenAI
func (c *Client) ChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage, maxTokens int, temperature float64) (*openai.ChatCompletionResponse, error) {
	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: float32(temperature),
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	return &resp, nil
}

// GetModel returns the current model being used
func (c *Client) GetModel() string {
	return c.model
}
