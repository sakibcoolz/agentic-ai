package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// AIClient wraps the OpenAI client with our custom functionality
type AIClient struct {
	client *openai.Client
	model  string
}

// NewAIClient creates a new AI client instance
func NewAIClient(apiKey string) *AIClient {
	client := openai.NewClient(apiKey)
	return &AIClient{
		client: client,
		model:  openai.GPT3Dot5Turbo, // Using GPT-3.5-turbo for cost efficiency
	}
}

// Chat sends a message to the AI and returns the response
func (ai *AIClient) Chat(ctx context.Context, message string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: ai.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful AI assistant teaching about agentic AI and Go programming.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
		MaxTokens:   500,
		Temperature: 0.7,
	}

	resp, err := ai.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}

// ValidateSetup checks if the AI client is properly configured
func (ai *AIClient) ValidateSetup(ctx context.Context) error {
	_, err := ai.Chat(ctx, "Hello! Can you confirm you're working?")
	if err != nil {
		return fmt.Errorf("setup validation failed: %w", err)
	}
	return nil
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

	// Create AI client
	aiClient := NewAIClient(apiKey)

	// Validate setup
	ctx := context.Background()
	fmt.Println("🔍 Validating setup...")
	if err := aiClient.ValidateSetup(ctx); err != nil {
		log.Fatalf("Setup validation failed: %v", err)
	}
	fmt.Println("✅ Setup validated successfully!")

	// Start interactive chat
	fmt.Println("\n🤖 Welcome to your first AI Agent!")
	fmt.Println("Ask me anything about AI, agents, or Go programming.")
	fmt.Println("Type 'quit' to exit.\n")

	scanner := bufio.NewScanner(os.Stdin)

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
			fmt.Println("👋 Goodbye!")
			break
		}

		fmt.Print("AI: ")
		response, err := aiClient.Chat(ctx, input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println(response)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}
