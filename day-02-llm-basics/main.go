package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// LLMProvider represents different LLM providers
type LLMProvider string

const (
	ProviderOpenAI LLMProvider = "openai"
	// We'll add more providers in future days
)

// ModelConfig holds model-specific configuration
type ModelConfig struct {
	Name         string
	MaxTokens    int
	TokenCost    float64 // Cost per 1000 tokens
	ContextLimit int
}

// PredefinedModels contains configuration for common models
var PredefinedModels = map[string]ModelConfig{
	"gpt-3.5-turbo": {
		Name:         "gpt-3.5-turbo",
		MaxTokens:    4096,
		TokenCost:    0.002, // $0.002 per 1K tokens
		ContextLimit: 4096,
	},
	"gpt-4": {
		Name:         "gpt-4",
		MaxTokens:    8192,
		TokenCost:    0.03, // $0.03 per 1K tokens
		ContextLimit: 8192,
	},
	"gpt-4-turbo": {
		Name:         "gpt-4-turbo-preview",
		MaxTokens:    4096,
		TokenCost:    0.01, // $0.01 per 1K tokens
		ContextLimit: 128000,
	},
}

// Usage tracks API usage statistics
type Usage struct {
	TotalTokens   int
	TotalRequests int
	TotalCost     float64
	StartTime     time.Time
}

// AdvancedLLMClient provides enhanced LLM capabilities
type AdvancedLLMClient struct {
	client    *openai.Client
	config    ModelConfig
	usage     *Usage
	retryMax  int
	retryWait time.Duration
}

// NewAdvancedLLMClient creates a new advanced LLM client
func NewAdvancedLLMClient(apiKey string, modelName string) *AdvancedLLMClient {
	config, exists := PredefinedModels[modelName]
	if !exists {
		log.Printf("Unknown model %s, using default gpt-3.5-turbo", modelName)
		config = PredefinedModels["gpt-3.5-turbo"]
	}

	return &AdvancedLLMClient{
		client: openai.NewClient(apiKey),
		config: config,
		usage: &Usage{
			StartTime: time.Now(),
		},
		retryMax:  3,
		retryWait: time.Second,
	}
}

// ChatWithRetry sends a message with retry logic
func (c *AdvancedLLMClient) ChatWithRetry(ctx context.Context, message string, systemPrompt string) (string, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retryMax; attempt++ {
		if attempt > 0 {
			fmt.Printf("🔄 Retry attempt %d/%d...\n", attempt, c.retryMax)
			time.Sleep(c.retryWait * time.Duration(attempt)) // Exponential backoff
		}

		response, err := c.chat(ctx, message, systemPrompt)
		if err == nil {
			return response, nil
		}

		lastErr = err

		// Don't retry on certain errors
		if strings.Contains(err.Error(), "invalid_request_error") {
			break
		}
	}

	return "", fmt.Errorf("failed after %d retries: %w", c.retryMax, lastErr)
}

// chat performs the actual API call
func (c *AdvancedLLMClient) chat(ctx context.Context, message string, systemPrompt string) (string, error) {
	if systemPrompt == "" {
		systemPrompt = "You are a helpful AI assistant specializing in agentic AI and Go programming."
	}

	req := openai.ChatCompletionRequest{
		Model: c.config.Name,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: 0.7,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	// Update usage statistics
	c.updateUsage(resp.Usage)

	return resp.Choices[0].Message.Content, nil
}

// ChatStream handles streaming responses
func (c *AdvancedLLMClient) ChatStream(ctx context.Context, message string, systemPrompt string) error {
	if systemPrompt == "" {
		systemPrompt = "You are a helpful AI assistant specializing in agentic AI and Go programming."
	}

	req := openai.ChatCompletionRequest{
		Model: c.config.Name,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: 0.7,
		Stream:      true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}
	defer stream.Close()

	fmt.Print("AI: ")
	for {
		response, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("stream error: %w", err)
		}

		if len(response.Choices) > 0 {
			delta := response.Choices[0].Delta.Content
			fmt.Print(delta)
		}
	}
	fmt.Println()

	return nil
}

// updateUsage updates usage statistics
func (c *AdvancedLLMClient) updateUsage(usage openai.Usage) {
	c.usage.TotalTokens += usage.TotalTokens
	c.usage.TotalRequests++
	c.usage.TotalCost += float64(usage.TotalTokens) * c.config.TokenCost / 1000
}

// GetUsageStats returns current usage statistics
func (c *AdvancedLLMClient) GetUsageStats() Usage {
	return *c.usage
}

// EstimateCost estimates the cost for a given number of tokens
func (c *AdvancedLLMClient) EstimateCost(tokens int) float64 {
	return float64(tokens) * c.config.TokenCost / 1000
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create advanced LLM client
	fmt.Println("Available models:")
	for name, config := range PredefinedModels {
		fmt.Printf("- %s (Cost: $%.4f per 1K tokens)\n", name, config.TokenCost)
	}

	fmt.Print("\nSelect model (default: gpt-3.5-turbo): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	modelName := strings.TrimSpace(scanner.Text())
	if modelName == "" {
		modelName = "gpt-3.5-turbo"
	}

	client := NewAdvancedLLMClient(apiKey, modelName)
	ctx := context.Background()

	fmt.Printf("\n🤖 Advanced LLM Client using %s\n", client.config.Name)
	fmt.Println("Features: Retry logic, usage tracking, streaming")
	fmt.Println("Commands: 'stream <message>' for streaming, 'stats' for usage, 'quit' to exit")
	fmt.Println()

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if strings.ToLower(input) == "quit" {
			break
		}

		if strings.ToLower(input) == "stats" {
			stats := client.GetUsageStats()
			fmt.Printf("📊 Usage Statistics:\n")
			fmt.Printf("   Requests: %d\n", stats.TotalRequests)
			fmt.Printf("   Tokens: %d\n", stats.TotalTokens)
			fmt.Printf("   Estimated Cost: $%.4f\n", stats.TotalCost)
			fmt.Printf("   Session Time: %v\n", time.Since(stats.StartTime).Round(time.Second))
			continue
		}

		if strings.HasPrefix(strings.ToLower(input), "stream ") {
			message := input[7:] // Remove "stream " prefix
			if err := client.ChatStream(ctx, message, ""); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			fmt.Println()
		} else {
			response, err := client.ChatWithRetry(ctx, input, "")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("AI: %s\n\n", response)
		}
	}

	// Final statistics
	stats := client.GetUsageStats()
	fmt.Printf("\n📊 Final Session Statistics:\n")
	fmt.Printf("   Total Requests: %d\n", stats.TotalRequests)
	fmt.Printf("   Total Tokens: %d\n", stats.TotalTokens)
	fmt.Printf("   Total Cost: $%.4f\n", stats.TotalCost)
	fmt.Printf("   Session Duration: %v\n", time.Since(stats.StartTime).Round(time.Second))
	fmt.Println("👋 Thanks for using the Advanced LLM Client!")
}
