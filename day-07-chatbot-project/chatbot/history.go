package chatbot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ConversationMessage represents a single message in a conversation
type ConversationMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// SavedConversation represents a complete saved conversation
type SavedConversation struct {
	Name      string                `json:"name"`
	Messages  []ConversationMessage `json:"messages"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// History manages conversation persistence
type History struct {
	saveDirectory string
}

// NewHistory creates a new history manager
func NewHistory(saveDirectory string) (*History, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(saveDirectory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create save directory: %w", err)
	}

	return &History{
		saveDirectory: saveDirectory,
	}, nil
}

// Save saves a conversation with the given name
func (h *History) Save(name string, messages []ConversationMessage) error {
	// Add timestamps to messages if they don't have them
	for i := range messages {
		if messages[i].Timestamp.IsZero() {
			messages[i].Timestamp = time.Now()
		}
	}

	conversation := SavedConversation{
		Name:      name,
		Messages:  messages,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check if conversation exists and preserve creation time
	existing, err := h.Load(name)
	if err == nil {
		conversation.CreatedAt = existing.CreatedAt
	}

	filename := h.getFilename(name)
	data, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write conversation file: %w", err)
	}

	return nil
}

// Load loads a conversation by name
func (h *History) Load(name string) (*SavedConversation, error) {
	filename := h.getFilename(name)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read conversation file: %w", err)
	}

	var conversation SavedConversation
	if err := json.Unmarshal(data, &conversation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation: %w", err)
	}

	return &conversation, nil
}

// List returns a list of all saved conversation names
func (h *History) List() []string {
	files, err := ioutil.ReadDir(h.saveDirectory)
	if err != nil {
		return nil
	}

	var conversations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			name := strings.TrimSuffix(file.Name(), ".json")
			conversations = append(conversations, name)
		}
	}

	return conversations
}

// Delete removes a saved conversation
func (h *History) Delete(name string) error {
	filename := h.getFilename(name)
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete conversation: %w", err)
	}
	return nil
}

// Exists checks if a conversation with the given name exists
func (h *History) Exists(name string) bool {
	filename := h.getFilename(name)
	_, err := os.Stat(filename)
	return err == nil
}

// getFilename returns the full path for a conversation file
func (h *History) getFilename(name string) string {
	// Sanitize the name to make it filesystem-safe
	safeName := strings.ReplaceAll(name, "/", "_")
	safeName = strings.ReplaceAll(safeName, "\\", "_")
	safeName = strings.ReplaceAll(safeName, ":", "_")
	safeName = strings.ReplaceAll(safeName, "*", "_")
	safeName = strings.ReplaceAll(safeName, "?", "_")
	safeName = strings.ReplaceAll(safeName, "\"", "_")
	safeName = strings.ReplaceAll(safeName, "<", "_")
	safeName = strings.ReplaceAll(safeName, ">", "_")
	safeName = strings.ReplaceAll(safeName, "|", "_")

	return filepath.Join(h.saveDirectory, safeName+".json")
}
