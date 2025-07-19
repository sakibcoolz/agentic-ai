package chatbot

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"

	"chatbot/config"
	"chatbot/llm"
)

// Bot represents the main chatbot instance
type Bot struct {
	llmClient *llm.Client
	config    *Config
	memory    *Memory
	history   *History
	stats     *Stats
}

// Config holds bot-specific configuration
type Config struct {
	MaxTokens     int
	Temperature   float64
	MaxHistory    int
	RetryAttempts int
	RetryDelay    time.Duration
	SaveDirectory string
}

// Stats tracks bot usage statistics
type Stats struct {
	MessageCount int
	TokensUsed   int
	CurrentMode  string
	StartTime    time.Time
}

// New creates a new chatbot instance
func New(llmClient *llm.Client, cfg *config.Config) (*Bot, error) {
	botConfig := &Config{
		MaxTokens:     cfg.MaxTokens,
		Temperature:   cfg.Temperature,
		MaxHistory:    cfg.MaxHistory,
		RetryAttempts: cfg.RetryAttempts,
		RetryDelay:    cfg.RetryDelay,
		SaveDirectory: cfg.SaveDirectory,
	}

	memory := NewMemory(cfg.MaxHistory)
	history, err := NewHistory(cfg.SaveDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize history: %w", err)
	}

	stats := &Stats{
		MessageCount: 0,
		TokensUsed:   0,
		CurrentMode:  "assistant",
		StartTime:    time.Now(),
	}

	bot := &Bot{
		llmClient: llmClient,
		config:    botConfig,
		memory:    memory,
		history:   history,
		stats:     stats,
	}

	// Set initial system message
	bot.memory.SetSystemMessage(llm.GetSystemPrompt("assistant"))

	return bot, nil
}

// ProcessMessage processes a user message and returns the bot's response
func (b *Bot) ProcessMessage(ctx context.Context, message string) (string, error) {
	// Add user message to memory
	b.memory.AddMessage("user", message)
	b.stats.MessageCount++

	// Get conversation messages for the API
	messages := b.memory.GetMessages()

	// Try to get response with retries
	var response *openai.ChatCompletionResponse
	var err error

	for attempt := 0; attempt < b.config.RetryAttempts; attempt++ {
		response, err = b.llmClient.ChatCompletion(
			ctx,
			messages,
			b.config.MaxTokens,
			b.config.Temperature,
		)

		if err == nil {
			break
		}

		if attempt < b.config.RetryAttempts-1 {
			time.Sleep(b.config.RetryDelay * time.Duration(attempt+1))
		}
	}

	if err != nil {
		return "", fmt.Errorf("failed to get response after %d attempts: %w", b.config.RetryAttempts, err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	botResponse := response.Choices[0].Message.Content

	// Add bot response to memory
	b.memory.AddMessage("assistant", botResponse)

	// Update token usage
	b.stats.TokensUsed += response.Usage.TotalTokens

	return botResponse, nil
}

// SetMode changes the conversation mode
func (b *Bot) SetMode(mode string) error {
	availableModes := llm.GetAvailableModes()
	valid := false
	for _, m := range availableModes {
		if m == mode {
			valid = true
			break
		}
	}

	if !valid {
		return fmt.Errorf("invalid mode '%s'. Available modes: %v", mode, availableModes)
	}

	b.stats.CurrentMode = mode
	b.memory.SetSystemMessage(llm.GetSystemPrompt(mode))
	return nil
}

// ClearMemory clears the conversation memory
func (b *Bot) ClearMemory() {
	b.memory.Clear()
	b.memory.SetSystemMessage(llm.GetSystemPrompt(b.stats.CurrentMode))
}

// SaveConversation saves the current conversation
func (b *Bot) SaveConversation(name string) error {
	conversation := b.memory.GetConversation()
	return b.history.Save(name, conversation)
}

// LoadConversation loads a saved conversation
func (b *Bot) LoadConversation(name string) error {
	conversation, err := b.history.Load(name)
	if err != nil {
		return err
	}

	b.memory.LoadConversation(conversation.Messages)
	return nil
}

// ListConversations returns a list of saved conversations
func (b *Bot) ListConversations() []string {
	return b.history.List()
}

// GetStats returns current bot statistics
func (b *Bot) GetStats() Stats {
	return *b.stats
}
