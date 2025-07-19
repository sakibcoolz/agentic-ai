package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// Tool represents a function that the agent can call
type Tool struct {
	Definition openai.FunctionDefinition
	Handler    func(args map[string]interface{}) (string, error)
}

// AgentWithTools represents an AI agent that can use tools
type AgentWithTools struct {
	client       *openai.Client
	tools        map[string]Tool
	conversation []openai.ChatCompletionMessage
}

// NewAgentWithTools creates a new agent with tool capabilities
func NewAgentWithTools(apiKey string) *AgentWithTools {
	agent := &AgentWithTools{
		client:       openai.NewClient(apiKey),
		tools:        make(map[string]Tool),
		conversation: []openai.ChatCompletionMessage{},
	}

	// Add system message
	agent.conversation = append(agent.conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a helpful AI assistant with access to various tools. Use the available tools when needed to provide accurate and helpful responses.",
	})

	// Register built-in tools
	agent.registerBuiltinTools()

	return agent
}

// registerBuiltinTools adds the default tools to the agent
func (a *AgentWithTools) registerBuiltinTools() {
	// Calculator tool
	a.RegisterTool("calculator", Tool{
		Definition: openai.FunctionDefinition{
			Name:        "calculator",
			Description: "Perform mathematical calculations including basic arithmetic, trigonometry, and advanced math functions",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"operation": {
						Type:        jsonschema.String,
						Description: "The mathematical operation to perform (add, subtract, multiply, divide, power, sqrt, sin, cos, tan, log)",
					},
					"a": {
						Type:        jsonschema.Number,
						Description: "First number",
					},
					"b": {
						Type:        jsonschema.Number,
						Description: "Second number (optional for single-operand operations)",
					},
				},
				Required: []string{"operation", "a"},
			},
		},
		Handler: a.handleCalculator,
	})

	// Current time tool
	a.RegisterTool("get_current_time", Tool{
		Definition: openai.FunctionDefinition{
			Name:        "get_current_time",
			Description: "Get the current date and time",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"format": {
						Type:        jsonschema.String,
						Description: "Time format preference (default, iso, unix)",
						Enum:        []string{"default", "iso", "unix"},
					},
				},
			},
		},
		Handler: a.handleCurrentTime,
	})

	// Text analyzer tool
	a.RegisterTool("analyze_text", Tool{
		Definition: openai.FunctionDefinition{
			Name:        "analyze_text",
			Description: "Analyze text and provide statistics like word count, character count, and reading time",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"text": {
						Type:        jsonschema.String,
						Description: "The text to analyze",
					},
				},
				Required: []string{"text"},
			},
		},
		Handler: a.handleTextAnalysis,
	})
}

// RegisterTool adds a new tool to the agent
func (a *AgentWithTools) RegisterTool(name string, tool Tool) {
	a.tools[name] = tool
}

// handleCalculator implements the calculator tool
func (a *AgentWithTools) handleCalculator(args map[string]interface{}) (string, error) {
	operation, ok := args["operation"].(string)
	if !ok {
		return "", fmt.Errorf("operation must be a string")
	}

	aVal, ok := args["a"].(float64)
	if !ok {
		return "", fmt.Errorf("parameter 'a' must be a number")
	}

	var result float64

	switch operation {
	case "add":
		bVal, ok := args["b"].(float64)
		if !ok {
			return "", fmt.Errorf("parameter 'b' required for addition")
		}
		result = aVal + bVal
	case "subtract":
		bVal, ok := args["b"].(float64)
		if !ok {
			return "", fmt.Errorf("parameter 'b' required for subtraction")
		}
		result = aVal - bVal
	case "multiply":
		bVal, ok := args["b"].(float64)
		if !ok {
			return "", fmt.Errorf("parameter 'b' required for multiplication")
		}
		result = aVal * bVal
	case "divide":
		bVal, ok := args["b"].(float64)
		if !ok {
			return "", fmt.Errorf("parameter 'b' required for division")
		}
		if bVal == 0 {
			return "", fmt.Errorf("division by zero")
		}
		result = aVal / bVal
	case "power":
		bVal, ok := args["b"].(float64)
		if !ok {
			return "", fmt.Errorf("parameter 'b' required for power operation")
		}
		result = math.Pow(aVal, bVal)
	case "sqrt":
		if aVal < 0 {
			return "", fmt.Errorf("cannot take square root of negative number")
		}
		result = math.Sqrt(aVal)
	case "sin":
		result = math.Sin(aVal)
	case "cos":
		result = math.Cos(aVal)
	case "tan":
		result = math.Tan(aVal)
	case "log":
		if aVal <= 0 {
			return "", fmt.Errorf("logarithm requires positive number")
		}
		result = math.Log(aVal)
	default:
		return "", fmt.Errorf("unknown operation: %s", operation)
	}

	return fmt.Sprintf("%.6f", result), nil
}

// handleCurrentTime implements the current time tool
func (a *AgentWithTools) handleCurrentTime(args map[string]interface{}) (string, error) {
	format := "default"
	if f, ok := args["format"].(string); ok {
		format = f
	}

	now := time.Now()

	switch format {
	case "iso":
		return now.Format(time.RFC3339), nil
	case "unix":
		return strconv.FormatInt(now.Unix(), 10), nil
	default:
		return now.Format("Monday, January 2, 2006 at 3:04 PM MST"), nil
	}
}

// handleTextAnalysis implements the text analysis tool
func (a *AgentWithTools) handleTextAnalysis(args map[string]interface{}) (string, error) {
	text, ok := args["text"].(string)
	if !ok {
		return "", fmt.Errorf("text parameter must be a string")
	}

	words := strings.Fields(text)
	chars := len(text)
	charsNoSpaces := len(strings.ReplaceAll(text, " ", ""))
	lines := len(strings.Split(text, "\n"))

	// Estimate reading time (average 200 words per minute)
	readingTime := float64(len(words)) / 200.0

	analysis := fmt.Sprintf(`Text Analysis Results:
- Characters: %d (including spaces), %d (excluding spaces)
- Words: %d
- Lines: %d
- Estimated reading time: %.1f minutes`,
		chars, charsNoSpaces, len(words), lines, readingTime)

	return analysis, nil
}

// Chat processes a user message and handles any function calls
func (a *AgentWithTools) Chat(ctx context.Context, message string) (string, error) {
	// Add user message to conversation
	a.conversation = append(a.conversation, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	// Convert tools to OpenAI function definitions
	var functions []openai.FunctionDefinition
	for _, tool := range a.tools {
		functions = append(functions, tool.Definition)
	}

	for {
		req := openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Messages:    a.conversation,
			Functions:   functions,
			Temperature: 0.7,
		}

		resp, err := a.client.CreateChatCompletion(ctx, req)
		if err != nil {
			return "", fmt.Errorf("API call failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", fmt.Errorf("no response choices returned")
		}

		choice := resp.Choices[0]

		// Add assistant's response to conversation
		a.conversation = append(a.conversation, choice.Message)

		// Check if the model wants to call a function
		if choice.Message.FunctionCall != nil {
			funcCall := choice.Message.FunctionCall

			fmt.Printf("🔧 Calling tool: %s\n", funcCall.Name)

			// Parse function arguments
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(funcCall.Arguments), &args); err != nil {
				return "", fmt.Errorf("failed to parse function arguments: %w", err)
			}

			// Execute the function
			tool, exists := a.tools[funcCall.Name]
			if !exists {
				return "", fmt.Errorf("unknown function: %s", funcCall.Name)
			}

			result, err := tool.Handler(args)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			// Add function result to conversation
			a.conversation = append(a.conversation, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleFunction,
				Name:    funcCall.Name,
				Content: result,
			})

			// Continue the loop to get the model's response to the function result
			continue
		}

		// No function call, return the response
		return choice.Message.Content, nil
	}
}

// GetConversationHistory returns the current conversation
func (a *AgentWithTools) GetConversationHistory() []openai.ChatCompletionMessage {
	return a.conversation
}

// ClearConversation resets the conversation history
func (a *AgentWithTools) ClearConversation() {
	a.conversation = []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a helpful AI assistant with access to various tools. Use the available tools when needed to provide accurate and helpful responses.",
		},
	}
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

	// Create agent with tools
	agent := NewAgentWithTools(apiKey)

	fmt.Println("🤖 Function-Calling Agent Ready!")
	fmt.Println("\nAvailable tools:")
	for name, tool := range agent.tools {
		fmt.Printf("- %s: %s\n", name, tool.Definition.Description)
	}
	fmt.Println("\nTry asking me to:")
	fmt.Println("- Calculate something: 'What is 15 * 23?'")
	fmt.Println("- Get current time: 'What time is it?'")
	fmt.Println("- Analyze text: 'Analyze this text: Hello world'")
	fmt.Println("- Complex tasks: 'Calculate the area of a circle with radius 5'")
	fmt.Println("\nCommands: 'clear' to reset conversation, 'quit' to exit")

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

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

		if strings.ToLower(input) == "clear" {
			agent.ClearConversation()
			fmt.Println("🗑️ Conversation cleared!")
			continue
		}

		response, err := agent.Chat(ctx, input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("AI: %s\n\n", response)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}
