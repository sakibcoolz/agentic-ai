package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// PromptOptimizer helps test and improve prompt effectiveness
type PromptOptimizer struct {
	client *openai.Client
}

// TestResult represents the results of a prompt test
type TestResult struct {
	Prompt     string  `json:"prompt"`
	Response   string  `json:"response"`
	TokensUsed int     `json:"tokens_used"`
	Score      float64 `json:"score"`
	Feedback   string  `json:"feedback"`
}

// NewPromptOptimizer creates a new prompt optimization tool
func NewPromptOptimizer(apiKey string) *PromptOptimizer {
	return &PromptOptimizer{
		client: openai.NewClient(apiKey),
	}
}

// ABTestPrompts compares two different prompts for the same task
func (po *PromptOptimizer) ABTestPrompts(ctx context.Context, promptA, promptB string) (TestResult, TestResult, error) {
	resultA, err := po.testSinglePrompt(ctx, promptA)
	if err != nil {
		return TestResult{}, TestResult{}, fmt.Errorf("testing prompt A failed: %w", err)
	}

	resultB, err := po.testSinglePrompt(ctx, promptB)
	if err != nil {
		return TestResult{}, TestResult{}, fmt.Errorf("testing prompt B failed: %w", err)
	}

	// Score the prompts based on response quality indicators
	resultA.Score = po.scoreResponse(resultA.Response, resultA.TokensUsed)
	resultB.Score = po.scoreResponse(resultB.Response, resultB.TokensUsed)

	return resultA, resultB, nil
}

// testSinglePrompt executes a single prompt and returns results
func (po *PromptOptimizer) testSinglePrompt(ctx context.Context, prompt string) (TestResult, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	resp, err := po.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return TestResult{}, err
	}

	if len(resp.Choices) == 0 {
		return TestResult{}, fmt.Errorf("no response from model")
	}

	return TestResult{
		Prompt:     prompt,
		Response:   resp.Choices[0].Message.Content,
		TokensUsed: resp.Usage.TotalTokens,
	}, nil
}

// scoreResponse provides a simple quality score for responses
func (po *PromptOptimizer) scoreResponse(response string, tokensUsed int) float64 {
	score := 0.0

	// Length score (reasonable length is good)
	length := len(response)
	if length > 100 && length < 1000 {
		score += 0.3
	}

	// Token efficiency (lower tokens for same quality is better)
	if tokensUsed < 500 {
		score += 0.2
	}

	// Structure indicators (lists, paragraphs, etc.)
	if containsStructure(response) {
		score += 0.3
	}

	// Completeness indicators
	if appearsComplete(response) {
		score += 0.2
	}

	return score
}

// containsStructure checks if response has good structure
func containsStructure(response string) bool {
	// Simple heuristics for structured content
	indicators := []string{"1.", "2.", "-", "*", "\n\n", "**", "##"}
	for _, indicator := range indicators {
		if contains(response, indicator) {
			return true
		}
	}
	return false
}

// appearsComplete checks if response seems complete
func appearsComplete(response string) bool {
	// Simple heuristics for completeness
	if len(response) < 50 {
		return false
	}

	// Ends with proper punctuation
	lastChar := response[len(response)-1]
	return lastChar == '.' || lastChar == '!' || lastChar == '?'
}

// contains checks if a string contains a substring
func contains(str, substr string) bool {
	return len(str) >= len(substr) &&
		(len(substr) == 0 ||
			findSubstring(str, substr))
}

// findSubstring simple substring search
func findSubstring(str, substr string) bool {
	if len(substr) > len(str) {
		return false
	}

	for i := 0; i <= len(str)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if str[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// RunOptimizationLab demonstrates prompt A/B testing
func RunOptimizationLab() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	optimizer := NewPromptOptimizer(apiKey)
	ctx := context.Background()

	fmt.Println("🔬 Prompt A/B Testing Lab")
	fmt.Println("=========================")

	// Example: Testing two different approaches to code explanation
	promptA := `Explain how this Go code works:

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}`

	promptB := `You are a Go programming instructor. Explain the following code to a beginner programmer, including:
1. What the function does
2. How the logic works step by step
3. Any potential issues or improvements

Code:
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}`

	fmt.Println("Testing two different prompt approaches for code explanation...")
	fmt.Println()

	resultA, resultB, err := optimizer.ABTestPrompts(ctx, promptA, promptB)
	if err != nil {
		log.Fatalf("A/B test failed: %v", err)
	}

	// Display results
	fmt.Println("📊 A/B Test Results")
	fmt.Println("===================")

	fmt.Printf("\n🅰️  Prompt A (Basic):\n")
	fmt.Printf("Score: %.2f\n", resultA.Score)
	fmt.Printf("Tokens: %d\n", resultA.TokensUsed)
	fmt.Printf("Response:\n%s\n", resultA.Response)

	fmt.Printf("\n🅱️  Prompt B (Structured):\n")
	fmt.Printf("Score: %.2f\n", resultB.Score)
	fmt.Printf("Tokens: %d\n", resultB.TokensUsed)
	fmt.Printf("Response:\n%s\n", resultB.Response)

	// Determine winner
	fmt.Printf("\n🏆 Winner: ")
	if resultA.Score > resultB.Score {
		fmt.Println("Prompt A (Basic approach)")
		fmt.Printf("Advantage: %.2f points\n", resultA.Score-resultB.Score)
	} else if resultB.Score > resultA.Score {
		fmt.Println("Prompt B (Structured approach)")
		fmt.Printf("Advantage: %.2f points\n", resultB.Score-resultA.Score)
	} else {
		fmt.Println("Tie!")
	}

	fmt.Println("\n💡 Key Insights:")
	fmt.Println("- Structured prompts often provide more comprehensive responses")
	fmt.Println("- Clear instructions help the model understand expectations")
	fmt.Println("- Role-based prompting can improve response quality")
	fmt.Println("- Token efficiency vs. response quality is a key trade-off")
}
