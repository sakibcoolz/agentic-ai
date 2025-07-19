package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

// Message represents a single conversation message
type Message struct {
	ID         string                 `json:"id"`
	Role       string                 `json:"role"`
	Content    string                 `json:"content"`
	Timestamp  time.Time              `json:"timestamp"`
	Metadata   map[string]interface{} `json:"metadata"`
	TokensUsed int                    `json:"tokens_used"`
}

// ConversationSummary represents a summarized conversation segment
type ConversationSummary struct {
	ID             string    `json:"id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Summary        string    `json:"summary"`
	KeyTopics      []string  `json:"key_topics"`
	ImportantFacts []string  `json:"important_facts"`
	MessageCount   int       `json:"message_count"`
	TokensUsed     int       `json:"tokens_used"`
}

// UserMemory stores persistent information about a user
type UserMemory struct {
	UserID      string                 `json:"user_id"`
	Profile     map[string]interface{} `json:"profile"`
	Preferences map[string]interface{} `json:"preferences"`
	Facts       []MemoryFact           `json:"facts"`
	LastSeen    time.Time              `json:"last_seen"`
	Sessions    int                    `json:"sessions"`
}

// MemoryFact represents a learned fact about the user or conversation
type MemoryFact struct {
	ID         string                 `json:"id"`
	Fact       string                 `json:"fact"`
	Confidence float64                `json:"confidence"`
	Source     string                 `json:"source"`
	Timestamp  time.Time              `json:"timestamp"`
	Category   string                 `json:"category"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// ContextWindow manages the conversation context for LLM calls
type ContextWindow struct {
	Messages     []Message `json:"messages"`
	TokenLimit   int       `json:"token_limit"`
	TokensUsed   int       `json:"tokens_used"`
	SystemPrompt string    `json:"system_prompt"`
}

// MemoryManager handles all aspects of conversation memory
type MemoryManager struct {
	client              *openai.Client
	conversationHistory []Message
	summaries           []ConversationSummary
	userMemory          *UserMemory
	contextWindow       *ContextWindow
	config              MemoryConfig
}

// MemoryConfig holds configuration for memory management
type MemoryConfig struct {
	MaxMessages         int     `json:"max_messages"`
	MaxTokens           int     `json:"max_tokens"`
	SummaryThreshold    int     `json:"summary_threshold"`
	RelevanceThreshold  float64 `json:"relevance_threshold"`
	MemoryRetentionDays int     `json:"memory_retention_days"`
}

// NewMemoryManager creates a new memory management system
func NewMemoryManager(apiKey string, userID string) *MemoryManager {
	config := MemoryConfig{
		MaxMessages:         50,
		MaxTokens:           3000,
		SummaryThreshold:    20,
		RelevanceThreshold:  0.7,
		MemoryRetentionDays: 30,
	}

	contextWindow := &ContextWindow{
		Messages:     make([]Message, 0),
		TokenLimit:   config.MaxTokens,
		TokensUsed:   0,
		SystemPrompt: "You are a helpful AI assistant with memory of our conversation history.",
	}

	userMemory := &UserMemory{
		UserID:      userID,
		Profile:     make(map[string]interface{}),
		Preferences: make(map[string]interface{}),
		Facts:       make([]MemoryFact, 0),
		LastSeen:    time.Now(),
		Sessions:    1,
	}

	return &MemoryManager{
		client:              openai.NewClient(apiKey),
		conversationHistory: make([]Message, 0),
		summaries:           make([]ConversationSummary, 0),
		userMemory:          userMemory,
		contextWindow:       contextWindow,
		config:              config,
	}
}

// AddMessage adds a new message to the conversation
func (mm *MemoryManager) AddMessage(role, content string) {
	message := Message{
		ID:         fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		Role:       role,
		Content:    content,
		Timestamp:  time.Now(),
		Metadata:   make(map[string]interface{}),
		TokensUsed: mm.estimateTokens(content),
	}

	mm.conversationHistory = append(mm.conversationHistory, message)

	// Check if we need to summarize old messages
	if len(mm.conversationHistory) > mm.config.SummaryThreshold {
		mm.createSummary()
	}

	// Update context window
	mm.updateContextWindow()
}

// estimateTokens provides a rough token count estimate
func (mm *MemoryManager) estimateTokens(text string) int {
	// Rough estimation: ~4 characters per token
	return len(text) / 4
}

// createSummary creates a summary of older conversation messages
func (mm *MemoryManager) createSummary() {
	if len(mm.conversationHistory) < mm.config.SummaryThreshold {
		return
	}

	// Take the first half of messages for summarization
	splitPoint := len(mm.conversationHistory) / 2
	messagesToSummarize := mm.conversationHistory[:splitPoint]

	// Create conversation text for summarization
	conversationText := mm.buildConversationText(messagesToSummarize)

	// Generate summary using LLM
	summary, err := mm.generateSummary(context.Background(), conversationText)
	if err != nil {
		log.Printf("Failed to generate summary: %v", err)
		return
	}

	// Create summary object
	summaryObj := ConversationSummary{
		ID:             fmt.Sprintf("summary_%d", time.Now().UnixNano()),
		StartTime:      messagesToSummarize[0].Timestamp,
		EndTime:        messagesToSummarize[len(messagesToSummarize)-1].Timestamp,
		Summary:        summary,
		KeyTopics:      mm.extractTopics(conversationText),
		ImportantFacts: mm.extractFacts(summary),
		MessageCount:   len(messagesToSummarize),
		TokensUsed:     mm.calculateTokens(messagesToSummarize),
	}

	// Store summary and remove old messages
	mm.summaries = append(mm.summaries, summaryObj)
	mm.conversationHistory = mm.conversationHistory[splitPoint:]

	fmt.Printf("📝 Created conversation summary covering %d messages\n", len(messagesToSummarize))
}

// buildConversationText creates a text representation of messages
func (mm *MemoryManager) buildConversationText(messages []Message) string {
	var builder strings.Builder

	for _, msg := range messages {
		builder.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	return builder.String()
}

// generateSummary creates a summary using the LLM
func (mm *MemoryManager) generateSummary(ctx context.Context, conversationText string) (string, error) {
	prompt := fmt.Sprintf(`Please summarize the following conversation, highlighting:
1. Key topics discussed
2. Important decisions made
3. User preferences revealed
4. Any facts learned about the user
5. Action items or follow-ups

Conversation:
%s

Summary:`, conversationText)

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Temperature: 0.3,
		MaxTokens:   500,
	}

	resp, err := mm.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no summary generated")
	}

	return resp.Choices[0].Message.Content, nil
}

// extractTopics extracts key topics from conversation text
func (mm *MemoryManager) extractTopics(text string) []string {
	// Simple keyword extraction - in production, use more sophisticated NLP
	keywords := []string{}

	// Common technical topics
	techTopics := []string{"programming", "go", "golang", "code", "function", "api", "database", "web", "server"}

	textLower := strings.ToLower(text)
	for _, topic := range techTopics {
		if strings.Contains(textLower, topic) {
			keywords = append(keywords, topic)
		}
	}

	return keywords
}

// extractFacts extracts facts from summary text
func (mm *MemoryManager) extractFacts(summary string) []string {
	// Simple fact extraction - look for declarative sentences
	facts := []string{}

	sentences := strings.Split(summary, ". ")
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) > 20 && !strings.Contains(sentence, "?") {
			facts = append(facts, sentence)
		}
	}

	return facts
}

// calculateTokens sums up tokens used in messages
func (mm *MemoryManager) calculateTokens(messages []Message) int {
	total := 0
	for _, msg := range messages {
		total += msg.TokensUsed
	}
	return total
}

// updateContextWindow optimizes the context window for the next LLM call
func (mm *MemoryManager) updateContextWindow() {
	mm.contextWindow.Messages = make([]Message, 0)
	mm.contextWindow.TokensUsed = mm.estimateTokens(mm.contextWindow.SystemPrompt)

	// Add relevant summaries first
	relevantSummaries := mm.getRelevantSummaries(3)
	for _, summary := range relevantSummaries {
		summaryText := fmt.Sprintf("Previous conversation summary: %s", summary.Summary)
		tokens := mm.estimateTokens(summaryText)

		if mm.contextWindow.TokensUsed+tokens < mm.contextWindow.TokenLimit {
			mm.contextWindow.Messages = append(mm.contextWindow.Messages, Message{
				Role:       "system",
				Content:    summaryText,
				TokensUsed: tokens,
			})
			mm.contextWindow.TokensUsed += tokens
		}
	}

	// Add recent messages
	for i := len(mm.conversationHistory) - 1; i >= 0; i-- {
		message := mm.conversationHistory[i]
		if mm.contextWindow.TokensUsed+message.TokensUsed < mm.contextWindow.TokenLimit {
			mm.contextWindow.Messages = append([]Message{message}, mm.contextWindow.Messages...)
			mm.contextWindow.TokensUsed += message.TokensUsed
		} else {
			break
		}
	}
}

// getRelevantSummaries returns the most relevant conversation summaries
func (mm *MemoryManager) getRelevantSummaries(limit int) []ConversationSummary {
	if len(mm.summaries) == 0 {
		return []ConversationSummary{}
	}

	// Sort by recency for now - in production, use semantic similarity
	summaries := make([]ConversationSummary, len(mm.summaries))
	copy(summaries, mm.summaries)

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].EndTime.After(summaries[j].EndTime)
	})

	if limit > len(summaries) {
		limit = len(summaries)
	}

	return summaries[:limit]
}

// Chat processes a user message and generates a response
func (mm *MemoryManager) Chat(ctx context.Context, userMessage string) (string, error) {
	// Add user message to history
	mm.AddMessage("user", userMessage)

	// Build messages for LLM call
	messages := make([]openai.ChatCompletionMessage, 0)

	// Add system prompt
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: mm.buildSystemPrompt(),
	})

	// Add context messages
	for _, msg := range mm.contextWindow.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Make LLM call
	req := openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   800,
	}

	resp, err := mm.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	response := resp.Choices[0].Message.Content

	// Add assistant response to history
	mm.AddMessage("assistant", response)

	// Extract and store any new facts about the user
	mm.extractAndStoreFacts(userMessage, response)

	return response, nil
}

// buildSystemPrompt creates a context-aware system prompt
func (mm *MemoryManager) buildSystemPrompt() string {
	basePrompt := "You are a helpful AI assistant with memory of our conversation history."

	// Add user information if available
	if len(mm.userMemory.Facts) > 0 {
		basePrompt += "\n\nWhat I know about you:"
		for _, fact := range mm.userMemory.Facts {
			if fact.Confidence > 0.7 {
				basePrompt += fmt.Sprintf("\n- %s", fact.Fact)
			}
		}
	}

	// Add user preferences
	if len(mm.userMemory.Preferences) > 0 {
		basePrompt += "\n\nYour preferences:"
		for key, value := range mm.userMemory.Preferences {
			basePrompt += fmt.Sprintf("\n- %s: %v", key, value)
		}
	}

	return basePrompt
}

// extractAndStoreFacts extracts facts from the conversation
func (mm *MemoryManager) extractAndStoreFacts(userMessage, assistantResponse string) {
	// Simple fact extraction - look for "I am", "I like", "I work", etc.
	factPatterns := []string{
		"I am ", "I like ", "I work ", "I study ", "I live ",
		"My name is ", "I prefer ", "I use ", "I need ",
	}

	userLower := strings.ToLower(userMessage)

	for _, pattern := range factPatterns {
		if strings.Contains(userLower, pattern) {
			// Extract the sentence containing the fact
			sentences := strings.Split(userMessage, ".")
			for _, sentence := range sentences {
				if strings.Contains(strings.ToLower(sentence), pattern) {
					fact := MemoryFact{
						ID:         fmt.Sprintf("fact_%d", time.Now().UnixNano()),
						Fact:       strings.TrimSpace(sentence),
						Confidence: 0.8,
						Source:     "user_statement",
						Timestamp:  time.Now(),
						Category:   "personal",
						Metadata:   make(map[string]interface{}),
					}
					mm.userMemory.Facts = append(mm.userMemory.Facts, fact)
					break
				}
			}
		}
	}
}

// GetMemoryStats returns statistics about the memory system
func (mm *MemoryManager) GetMemoryStats() map[string]interface{} {
	return map[string]interface{}{
		"total_messages":       len(mm.conversationHistory),
		"summaries_created":    len(mm.summaries),
		"facts_learned":        len(mm.userMemory.Facts),
		"context_window_usage": fmt.Sprintf("%d/%d tokens", mm.contextWindow.TokensUsed, mm.contextWindow.TokenLimit),
		"user_sessions":        mm.userMemory.Sessions,
		"last_interaction":     mm.userMemory.LastSeen.Format("2006-01-02 15:04:05"),
	}
}

// GetConversationHistory returns the current conversation history
func (mm *MemoryManager) GetConversationHistory() []Message {
	return mm.conversationHistory
}

// GetUserFacts returns learned facts about the user
func (mm *MemoryManager) GetUserFacts() []MemoryFact {
	return mm.userMemory.Facts
}

// ClearMemory resets the memory system
func (mm *MemoryManager) ClearMemory() {
	mm.conversationHistory = make([]Message, 0)
	mm.summaries = make([]ConversationSummary, 0)
	mm.userMemory.Facts = make([]MemoryFact, 0)
	mm.updateContextWindow()
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

	// Create memory manager for a user
	userID := "demo_user_001"
	memoryManager := NewMemoryManager(apiKey, userID)
	ctx := context.Background()

	fmt.Println("🧠 Context Management & Memory System")
	fmt.Println("=====================================")
	fmt.Printf("User ID: %s\n", userID)
	fmt.Printf("Memory Config: %d messages, %d tokens max\n",
		memoryManager.config.MaxMessages, memoryManager.config.MaxTokens)
	fmt.Println()

	fmt.Println("💡 This AI assistant has memory! Try:")
	fmt.Println("- Tell it your name and preferences")
	fmt.Println("- Ask follow-up questions")
	fmt.Println("- Reference things you mentioned earlier")
	fmt.Println("- Have a long conversation to see summarization")
	fmt.Println()
	fmt.Println("Commands: 'stats' for memory info, 'facts' for learned facts, 'clear' to reset, 'quit' to exit")
	fmt.Println()

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
			fmt.Println("👋 Goodbye! Your memories are preserved for next time.")
			break
		}

		if strings.ToLower(input) == "stats" {
			stats := memoryManager.GetMemoryStats()
			fmt.Println("\n📊 Memory Statistics:")
			for key, value := range stats {
				fmt.Printf("  %s: %v\n", key, value)
			}
			fmt.Println()
			continue
		}

		if strings.ToLower(input) == "facts" {
			facts := memoryManager.GetUserFacts()
			fmt.Printf("\n🧠 Facts I've learned about you (%d):\n", len(facts))
			for i, fact := range facts {
				fmt.Printf("  %d. %s (confidence: %.2f)\n", i+1, fact.Fact, fact.Confidence)
			}
			fmt.Println()
			continue
		}

		if strings.ToLower(input) == "clear" {
			memoryManager.ClearMemory()
			fmt.Println("🗑️ Memory cleared!")
			continue
		}

		// Process chat message
		response, err := memoryManager.Chat(ctx, input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("AI: %s\n\n", response)

		// Show memory update if facts were learned
		currentFacts := len(memoryManager.GetUserFacts())
		if currentFacts > 0 {
			fmt.Printf("💭 [Learned %d facts about you so far]\n\n", currentFacts)
		}
	}

	// Final memory statistics
	fmt.Println("\n📊 Final Session Statistics:")
	stats := memoryManager.GetMemoryStats()
	for key, value := range stats {
		fmt.Printf("  %s: %v\n", key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading input: %v", err)
	}
}
