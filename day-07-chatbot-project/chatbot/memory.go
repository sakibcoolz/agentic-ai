package chatbot

import (
	"time"

	"github.com/sashabaranov/go-openai"
)

// Memory manages conversation history and context
type Memory struct {
	messages   []openai.ChatCompletionMessage
	maxHistory int
}

// NewMemory creates a new memory instance
func NewMemory(maxHistory int) *Memory {
	return &Memory{
		messages:   make([]openai.ChatCompletionMessage, 0),
		maxHistory: maxHistory,
	}
}

// AddMessage adds a message to memory
func (m *Memory) AddMessage(role, content string) {
	message := openai.ChatCompletionMessage{
		Role:    role,
		Content: content,
	}

	m.messages = append(m.messages, message)

	// Keep only the most recent messages (plus system message)
	if len(m.messages) > m.maxHistory+1 { // +1 for system message
		// Keep system message (first) and trim user/assistant messages
		systemMsg := m.messages[0]
		recentMessages := m.messages[len(m.messages)-m.maxHistory:]
		m.messages = append([]openai.ChatCompletionMessage{systemMsg}, recentMessages...)
	}
}

// SetSystemMessage sets or updates the system message
func (m *Memory) SetSystemMessage(content string) {
	systemMsg := openai.ChatCompletionMessage{
		Role:    "system",
		Content: content,
	}

	// If we already have messages and the first is a system message, replace it
	if len(m.messages) > 0 && m.messages[0].Role == "system" {
		m.messages[0] = systemMsg
	} else {
		// Insert system message at the beginning
		m.messages = append([]openai.ChatCompletionMessage{systemMsg}, m.messages...)
	}
}

// GetMessages returns all messages for API calls
func (m *Memory) GetMessages() []openai.ChatCompletionMessage {
	return m.messages
}

// Clear clears all messages from memory
func (m *Memory) Clear() {
	m.messages = make([]openai.ChatCompletionMessage, 0)
}

// GetConversation returns the conversation without system message for saving
func (m *Memory) GetConversation() []ConversationMessage {
	var conversation []ConversationMessage

	for _, msg := range m.messages {
		if msg.Role != "system" {
			conversation = append(conversation, ConversationMessage{
				Role:      msg.Role,
				Content:   msg.Content,
				Timestamp: time.Now(),
			})
		}
	}

	return conversation
}

// LoadConversation loads a conversation into memory
func (m *Memory) LoadConversation(conversation []ConversationMessage) {
	// Keep system message if it exists
	var systemMsg *openai.ChatCompletionMessage
	if len(m.messages) > 0 && m.messages[0].Role == "system" {
		systemMsg = &m.messages[0]
	}

	// Clear and reload
	m.messages = make([]openai.ChatCompletionMessage, 0)

	// Add system message back
	if systemMsg != nil {
		m.messages = append(m.messages, *systemMsg)
	}

	// Add conversation messages
	for _, msg := range conversation {
		m.AddMessage(msg.Role, msg.Content)
	}
}

// GetMessageCount returns the number of messages (excluding system)
func (m *Memory) GetMessageCount() int {
	count := len(m.messages)
	if count > 0 && m.messages[0].Role == "system" {
		count--
	}
	return count
}
