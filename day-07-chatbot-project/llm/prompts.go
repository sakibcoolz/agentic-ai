package llm

// SystemPrompts contains predefined system prompts for different conversation modes
var SystemPrompts = map[string]string{
	"casual": `You are a friendly, casual chatbot. Respond in a relaxed, conversational tone. 
Keep responses concise and engaging. Use emojis occasionally but don't overdo it.
Be helpful and show genuine interest in the conversation.`,

	"assistant": `You are a helpful AI assistant. Provide clear, accurate, and useful information.
Be professional but approachable. Structure your responses well and provide actionable advice when appropriate.
If you don't know something, say so rather than guessing.`,

	"creative": `You are a creative AI companion. Think outside the box and provide imaginative responses.
Help with creative writing, brainstorming, and artistic endeavors. Be playful with language and ideas.
Encourage creativity and offer unique perspectives.`,
}

// GetSystemPrompt returns the system prompt for a given mode
func GetSystemPrompt(mode string) string {
	if prompt, exists := SystemPrompts[mode]; exists {
		return prompt
	}
	return SystemPrompts["assistant"] // Default to assistant mode
}

// GetAvailableModes returns a list of available conversation modes
func GetAvailableModes() []string {
	modes := make([]string, 0, len(SystemPrompts))
	for mode := range SystemPrompts {
		modes = append(modes, mode)
	}
	return modes
}
