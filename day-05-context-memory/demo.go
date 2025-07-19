package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// SimpleSlidingWindow demonstrates basic sliding window memory
type SimpleSlidingWindow struct {
	messages []string
	maxSize  int
}

// NewSlidingWindow creates a new sliding window
func NewSlidingWindow(maxSize int) *SimpleSlidingWindow {
	return &SimpleSlidingWindow{
		messages: make([]string, 0),
		maxSize:  maxSize,
	}
}

// Add adds a message to the window
func (sw *SimpleSlidingWindow) Add(message string) {
	sw.messages = append(sw.messages, message)

	// Remove oldest messages if we exceed maxSize
	if len(sw.messages) > sw.maxSize {
		sw.messages = sw.messages[len(sw.messages)-sw.maxSize:]
	}
}

// GetMessages returns all messages in the window
func (sw *SimpleSlidingWindow) GetMessages() []string {
	return sw.messages
}

// GetContext returns formatted context for LLM
func (sw *SimpleSlidingWindow) GetContext() string {
	if len(sw.messages) == 0 {
		return "No previous context."
	}

	return "Previous conversation:\n" + strings.Join(sw.messages, "\n")
}

// TokenBudgetManager manages token allocation
type TokenBudgetManager struct {
	totalBudget    int
	systemTokens   int
	historyTokens  int
	responseTokens int
}

// NewTokenBudgetManager creates a new token budget manager
func NewTokenBudgetManager(totalBudget int) *TokenBudgetManager {
	return &TokenBudgetManager{
		totalBudget:    totalBudget,
		systemTokens:   200,                     // Reserve for system prompt
		responseTokens: 500,                     // Reserve for response
		historyTokens:  totalBudget - 200 - 500, // Remaining for history
	}
}

// CalculateAvailableHistory returns tokens available for history
func (tbm *TokenBudgetManager) CalculateAvailableHistory() int {
	return tbm.historyTokens
}

// AdjustForResponse adjusts budget based on expected response length
func (tbm *TokenBudgetManager) AdjustForResponse(expectedResponseTokens int) {
	tbm.responseTokens = expectedResponseTokens
	tbm.historyTokens = tbm.totalBudget - tbm.systemTokens - tbm.responseTokens

	if tbm.historyTokens < 0 {
		tbm.historyTokens = 100 // Minimum history
		tbm.responseTokens = tbm.totalBudget - tbm.systemTokens - tbm.historyTokens
	}
}

// SimpleFactExtractor extracts facts from conversation
type SimpleFactExtractor struct {
	patterns []string
}

// NewFactExtractor creates a new fact extractor
func NewFactExtractor() *SimpleFactExtractor {
	patterns := []string{
		"my name is",
		"i am",
		"i work",
		"i like",
		"i prefer",
		"i study",
		"i live",
		"i use",
		"i need",
		"i want",
	}

	return &SimpleFactExtractor{
		patterns: patterns,
	}
}

// ExtractFacts finds potential facts in a message
func (fe *SimpleFactExtractor) ExtractFacts(message string) []string {
	facts := []string{}

	// Split into sentences
	sentences := strings.Split(message, ".")

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		sentenceLower := strings.ToLower(sentence)

		// Check if sentence contains fact patterns
		for _, pattern := range fe.patterns {
			if strings.Contains(sentenceLower, pattern) && len(sentence) > 10 {
				facts = append(facts, sentence)
				break
			}
		}
	}

	return facts
}

// MemoryDemo demonstrates various memory concepts
type MemoryDemo struct {
	client        *openai.Client
	slidingWindow *SimpleSlidingWindow
	budgetManager *TokenBudgetManager
	factExtractor *SimpleFactExtractor
	learnedFacts  []string
}

// NewMemoryDemo creates a new memory demonstration
func NewMemoryDemo(apiKey string) *MemoryDemo {
	return &MemoryDemo{
		client:        openai.NewClient(apiKey),
		slidingWindow: NewSlidingWindow(5),
		budgetManager: NewTokenBudgetManager(2000),
		factExtractor: NewFactExtractor(),
		learnedFacts:  make([]string, 0),
	}
}

// estimateTokens provides rough token estimation
func (md *MemoryDemo) estimateTokens(text string) int {
	return len(strings.Fields(text)) // Roughly 1 token per word
}

// ProcessMessage processes a user message with memory features
func (md *MemoryDemo) ProcessMessage(ctx context.Context, userMessage string) (string, error) {
	// Extract facts from user message
	facts := md.factExtractor.ExtractFacts(userMessage)
	for _, fact := range facts {
		md.learnedFacts = append(md.learnedFacts, fact)
		fmt.Printf("💡 Learned: %s\n", fact)
	}

	// Add to sliding window
	md.slidingWindow.Add(fmt.Sprintf("User: %s", userMessage))

	// Build context with memory
	systemPrompt := md.buildSystemPrompt()
	contextHistory := md.slidingWindow.GetContext()

	// Check token budget
	systemTokens := md.estimateTokens(systemPrompt)
	contextTokens := md.estimateTokens(contextHistory)

	fmt.Printf("📊 Token usage: System=%d, Context=%d, Available=%d\n",
		systemTokens, contextTokens, md.budgetManager.CalculateAvailableHistory())

	// Prepare messages for LLM
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("%s\n\nCurrent message: %s", contextHistory, userMessage),
		},
	}

	// Make LLM call
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   500,
	}

	resp, err := md.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	response := resp.Choices[0].Message.Content

	// Add assistant response to sliding window
	md.slidingWindow.Add(fmt.Sprintf("Assistant: %s", response))

	return response, nil
}

// buildSystemPrompt creates a context-aware system prompt
func (md *MemoryDemo) buildSystemPrompt() string {
	basePrompt := "You are a helpful AI assistant with memory of our conversation."

	if len(md.learnedFacts) > 0 {
		basePrompt += "\n\nWhat I remember about you:"
		for i, fact := range md.learnedFacts {
			if i >= 5 { // Limit to 5 most recent facts
				break
			}
			basePrompt += fmt.Sprintf("\n- %s", fact)
		}
	}

	basePrompt += "\n\nUse this information to provide personalized and contextual responses."

	return basePrompt
}

// ShowMemoryState displays current memory state
func (md *MemoryDemo) ShowMemoryState() {
	fmt.Println("\n🧠 Current Memory State:")
	fmt.Printf("Sliding Window Size: %d/%d messages\n",
		len(md.slidingWindow.GetMessages()), md.slidingWindow.maxSize)

	fmt.Printf("Learned Facts: %d\n", len(md.learnedFacts))
	for i, fact := range md.learnedFacts {
		fmt.Printf("  %d. %s\n", i+1, fact)
	}

	fmt.Printf("Token Budget: %d total, %d available for history\n",
		md.budgetManager.totalBudget, md.budgetManager.CalculateAvailableHistory())

	fmt.Println("\nConversation Window:")
	for i, msg := range md.slidingWindow.GetMessages() {
		fmt.Printf("  %d. %s\n", i+1, msg)
	}
	fmt.Println()
}

// RunMemoryDemo demonstrates various memory concepts
func RunMemoryDemo() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get OpenAI API key
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	demo := NewMemoryDemo(apiKey)
	ctx := context.Background()

	fmt.Println("🔬 Memory Concepts Demonstration")
	fmt.Println("================================")
	fmt.Println()

	// Demo 1: Sliding Window
	fmt.Println("📺 Demo 1: Sliding Window Memory")
	fmt.Println("Adding messages to a sliding window with max size 3...")

	window := NewSlidingWindow(3)
	messages := []string{
		"Hello, how are you?",
		"I'm learning about AI",
		"This is interesting",
		"Can you help me with Go?",
		"I love programming",
		"Memory systems are cool",
	}

	for i, msg := range messages {
		window.Add(msg)
		fmt.Printf("After message %d: %v\n", i+1, window.GetMessages())
	}
	fmt.Println()

	// Demo 2: Token Budget Management
	fmt.Println("💰 Demo 2: Token Budget Management")
	budget := NewTokenBudgetManager(1500)
	fmt.Printf("Total budget: %d tokens\n", budget.totalBudget)
	fmt.Printf("System tokens: %d\n", budget.systemTokens)
	fmt.Printf("Response tokens: %d\n", budget.responseTokens)
	fmt.Printf("Available for history: %d\n", budget.CalculateAvailableHistory())

	// Adjust for longer response
	budget.AdjustForResponse(800)
	fmt.Printf("After adjusting for 800-token response:\n")
	fmt.Printf("Available for history: %d\n", budget.CalculateAvailableHistory())
	fmt.Println()

	// Demo 3: Fact Extraction
	fmt.Println("🔍 Demo 3: Fact Extraction")
	extractor := NewFactExtractor()
	testMessages := []string{
		"Hi, my name is John and I work as a software engineer.",
		"I like playing guitar in my free time.",
		"I am studying machine learning at university.",
		"The weather is nice today.",
		"I prefer Python but I'm learning Go now.",
	}

	for _, msg := range testMessages {
		facts := extractor.ExtractFacts(msg)
		fmt.Printf("Message: %s\n", msg)
		fmt.Printf("Facts: %v\n", facts)
		fmt.Println()
	}

	// Demo 4: Interactive Memory System
	fmt.Println("🎮 Demo 4: Interactive Memory System")
	fmt.Println("Try these example interactions:")

	exampleInteractions := []string{
		"Hi, my name is Alice and I'm a data scientist",
		"I'm working on a machine learning project",
		"What did I tell you about my profession?",
		"I prefer using Python for data analysis",
		"Can you help me with Go programming?",
	}

	for i, interaction := range exampleInteractions {
		fmt.Printf("\n--- Interaction %d ---\n", i+1)
		fmt.Printf("User: %s\n", interaction)

		response, err := demo.ProcessMessage(ctx, interaction)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Assistant: %s\n", response)

		// Show memory state every few interactions
		if i == 2 || i == 4 {
			demo.ShowMemoryState()
		}
	}

	fmt.Println("\n✨ Memory demonstration complete!")
	fmt.Println("Key concepts demonstrated:")
	fmt.Println("- Sliding window maintains recent context")
	fmt.Println("- Token budgets manage LLM input limits")
	fmt.Println("- Fact extraction learns about users")
	fmt.Println("- System prompts incorporate learned facts")
}
