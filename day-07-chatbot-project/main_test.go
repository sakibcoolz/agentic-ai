package main

import (
	"testing"
	"time"

	"chatbot/chatbot"
	"chatbot/config"
	"chatbot/llm"
)

func TestChatbotInitialization(t *testing.T) {
	// Test LLM client creation
	client, err := llm.NewClient("test-key", "gpt-3.5-turbo")
	if err != nil {
		t.Fatalf("Failed to create LLM client: %v", err)
	}

	if client.GetModel() != "gpt-3.5-turbo" {
		t.Errorf("Expected model gpt-3.5-turbo, got %s", client.GetModel())
	}
}

func TestSystemPrompts(t *testing.T) {
	modes := llm.GetAvailableModes()
	expectedModes := []string{"casual", "assistant", "creative"}

	if len(modes) != len(expectedModes) {
		t.Errorf("Expected %d modes, got %d", len(expectedModes), len(modes))
	}

	for _, mode := range expectedModes {
		prompt := llm.GetSystemPrompt(mode)
		if prompt == "" {
			t.Errorf("Empty prompt for mode: %s", mode)
		}
	}

	// Test default prompt for invalid mode
	defaultPrompt := llm.GetSystemPrompt("invalid")
	assistantPrompt := llm.GetSystemPrompt("assistant")
	if defaultPrompt != assistantPrompt {
		t.Error("Invalid mode should return assistant prompt as default")
	}
}

func TestMemoryManagement(t *testing.T) {
	memory := chatbot.NewMemory(3)

	// Test system message
	memory.SetSystemMessage("You are a test bot")
	messages := memory.GetMessages()
	if len(messages) != 1 || messages[0].Role != "system" {
		t.Error("System message not set correctly")
	}

	// Test adding messages
	memory.AddMessage("user", "Hello")
	memory.AddMessage("assistant", "Hi there!")

	if memory.GetMessageCount() != 2 {
		t.Errorf("Expected 2 messages, got %d", memory.GetMessageCount())
	}

	// Test memory limit
	memory.AddMessage("user", "Message 3")
	memory.AddMessage("assistant", "Response 3")
	memory.AddMessage("user", "Message 4")
	memory.AddMessage("assistant", "Response 4")

	// Should have system + max 3 pairs = 7 total messages
	if len(memory.GetMessages()) > 7 {
		t.Errorf("Memory exceeded limit: %d messages", len(memory.GetMessages()))
	}
}

func TestConversationHistory(t *testing.T) {
	// Create temporary directory for testing
	tempDir := "/tmp/chatbot-test"
	history, err := chatbot.NewHistory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create history: %v", err)
	}

	// Test saving and loading
	messages := []chatbot.ConversationMessage{
		{Role: "user", Content: "Hello", Timestamp: time.Now()},
		{Role: "assistant", Content: "Hi!", Timestamp: time.Now()},
	}

	err = history.Save("test-conversation", messages)
	if err != nil {
		t.Fatalf("Failed to save conversation: %v", err)
	}

	loaded, err := history.Load("test-conversation")
	if err != nil {
		t.Fatalf("Failed to load conversation: %v", err)
	}

	if len(loaded.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(loaded.Messages))
	}

	// Test listing conversations
	conversations := history.List()
	found := false
	for _, conv := range conversations {
		if conv == "test-conversation" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Saved conversation not found in list")
	}
}

func TestConfigLoading(t *testing.T) {
	// Set test environment variables
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("MAX_TOKENS", "100")
	t.Setenv("TEMPERATURE", "0.5")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.OpenAIAPIKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", cfg.OpenAIAPIKey)
	}

	if cfg.MaxTokens != 100 {
		t.Errorf("Expected MaxTokens 100, got %d", cfg.MaxTokens)
	}

	if cfg.Temperature != 0.5 {
		t.Errorf("Expected Temperature 0.5, got %f", cfg.Temperature)
	}
}

func TestErrorHandling(t *testing.T) {
	// Test invalid API key
	_, err := llm.NewClient("", "gpt-3.5-turbo")
	if err == nil {
		t.Error("Expected error for empty API key")
	}

	// Test loading non-existent conversation
	history, _ := chatbot.NewHistory("/tmp/chatbot-test")
	_, err = history.Load("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent conversation")
	}
}
