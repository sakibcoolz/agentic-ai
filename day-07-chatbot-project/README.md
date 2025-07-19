# Day 7: Week 1 Project - Simple Chatbot

## 🎯 Project Overview

Build a complete chatbot application that integrates all the concepts learned in Week 1:
- LLM integration with OpenAI API
- Prompt engineering techniques
- Context and memory management
- Error handling and retries
- Clean Go architecture

## 🚀 What You'll Build

A command-line chatbot with the following features:
- Interactive conversation loop
- Conversation memory (remembers previous messages)
- Multiple conversation modes (casual, assistant, creative)
- Graceful error handling with retries
- Conversation history saving/loading
- Rate limiting and API optimization

## 📋 Learning Objectives

By completing this project, you will:
- Apply Week 1 concepts in a real application
- Understand chatbot architecture patterns
- Implement conversation state management
- Practice error handling in production scenarios
- Learn to optimize API usage and costs

## 🛠 Prerequisites

- Completed Days 1-6
- OpenAI API key configured
- Go 1.21+ installed
- Basic understanding of JSON and file I/O

## 📁 Project Structure

```
day-07-chatbot-project/
├── README.md           # This file
├── main.go            # Entry point
├── go.mod             # Go module
├── go.sum             # Dependencies
├── config/
│   └── config.go      # Configuration management
├── chatbot/
│   ├── bot.go         # Main chatbot logic
│   ├── memory.go      # Conversation memory
│   ├── modes.go       # Conversation modes
│   └── history.go     # Conversation persistence
├── llm/
│   ├── client.go      # OpenAI client wrapper
│   └── prompts.go     # Prompt templates
├── utils/
│   ├── errors.go      # Error handling utilities
│   └── retry.go       # Retry mechanisms
└── data/
    └── conversations/ # Saved conversations
```

## 🏗 Implementation Steps

### Step 1: Project Setup
- Initialize Go module
- Set up basic project structure
- Configure environment variables

### Step 2: LLM Integration
- Create OpenAI client wrapper
- Implement prompt templates
- Add response parsing

### Step 3: Conversation Management
- Design conversation state
- Implement memory system
- Add conversation modes

### Step 4: Error Handling
- Implement retry mechanisms
- Add graceful error recovery
- Handle rate limiting

### Step 5: User Interface
- Create interactive CLI
- Add conversation commands
- Implement help system

### Step 6: Persistence
- Save conversation history
- Load previous conversations
- Export/import functionality

## 🎮 Features to Implement

### Core Features
- [x] Basic chat loop
- [x] OpenAI API integration
- [x] Conversation memory
- [x] Error handling with retries
- [x] Multiple conversation modes

### Advanced Features
- [ ] Conversation persistence
- [ ] Chat history browsing
- [ ] Custom system prompts
- [ ] Token usage tracking
- [ ] Export conversations
- [ ] Conversation search

### Bonus Features
- [ ] Web interface (optional)
- [ ] Multiple LLM provider support
- [ ] Plugin system for custom commands
- [ ] Voice input/output
- [ ] Conversation analytics

## 🚀 Quick Start

1. **Setup the project:**
```bash
cd day-07-chatbot-project
go mod init chatbot
go mod tidy
```

2. **Configure environment:**
```bash
export OPENAI_API_KEY="your-api-key-here"
```

3. **Run the chatbot:**
```bash
go run main.go
```

4. **Start chatting:**
```
Welcome to the Simple Chatbot!
Type 'help' for commands, 'quit' to exit.

You: Hello!
Bot: Hello! How can I help you today?
```

## 🧪 Testing

Run the test suite:
```bash
go test ./...
```

Test specific components:
```bash
go test ./chatbot
go test ./llm
```

## 📊 Example Usage

### Basic Conversation
```
You: What's the weather like?
Bot: I don't have access to real-time weather data, but I can help you think about weather-related topics or suggest ways to check the weather.

You: Can you remember what we talked about?
Bot: Yes! We were just discussing weather information and how I can help you with weather-related topics.
```

### Changing Modes
```
You: /mode creative
Bot: Switched to creative mode! I'm now ready to help with creative writing, brainstorming, and imaginative tasks.

You: Write a short poem about coding
Bot: Here's a short poem about coding:

    In lines of logic, neat and clean,
    Where algorithms dance unseen,
    We craft solutions, byte by byte,
    Turning problems into light.
```

### Managing Conversations
```
You: /save my-coding-chat
Bot: Conversation saved as 'my-coding-chat'

You: /load my-coding-chat
Bot: Conversation 'my-coding-chat' loaded successfully!

You: /history
Bot: Recent conversations:
    - my-coding-chat (5 messages)
    - weather-chat (3 messages)
    - creative-session (12 messages)
```

## 🎯 Learning Challenges

### Beginner Challenges
1. **Memory Optimization**: Limit conversation history to last N messages
2. **Cost Tracking**: Track and display token usage
3. **Prompt Variants**: Create different personality modes

### Intermediate Challenges
1. **Conversation Branching**: Allow multiple conversation threads
2. **Custom Commands**: Add special bot commands (weather, time, etc.)
3. **Response Streaming**: Stream responses as they're generated

### Advanced Challenges
1. **Plugin System**: Allow loading custom functionality
2. **Multi-LLM Support**: Switch between different AI providers
3. **Conversation Analytics**: Analyze conversation patterns

## 🔧 Configuration Options

The chatbot supports various configuration options:

```go
type Config struct {
    OpenAIAPIKey     string
    Model            string // gpt-3.5-turbo, gpt-4, etc.
    MaxTokens        int
    Temperature      float64
    MaxHistory       int     // Maximum messages to remember
    RetryAttempts    int
    RetryDelay       time.Duration
    SaveDirectory    string
}
```

## 📈 Extending the Project

### Week 2 Preview
This chatbot will serve as the foundation for Week 2 projects:
- Add RAG capabilities for knowledge retrieval
- Integrate vector databases for semantic search
- Implement document processing for custom knowledge bases

### Integration Ideas
- Connect to external APIs (weather, news, etc.)
- Add scheduling and reminder capabilities
- Implement multi-user support
- Create web or mobile interfaces

## 🐛 Troubleshooting

### Common Issues

1. **API Key Issues**
```
Error: invalid API key
Solution: Check your OPENAI_API_KEY environment variable
```

2. **Rate Limiting**
```
Error: rate limit exceeded
Solution: The bot includes automatic retry with exponential backoff
```

3. **Memory Issues**
```
Error: conversation too long
Solution: Adjust MaxHistory in configuration
```

## 📚 Additional Resources

- [OpenAI API Documentation](https://platform.openai.com/docs)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Conversation Design](https://developers.google.com/assistant/conversation-design)

## 🎉 Completion Criteria

You've successfully completed Day 7 when your chatbot:
- ✅ Responds to user messages using OpenAI API
- ✅ Maintains conversation context/memory
- ✅ Handles errors gracefully with retries
- ✅ Supports multiple conversation modes
- ✅ Can save and load conversation history
- ✅ Includes comprehensive error handling
- ✅ Has clean, modular code architecture

## 🚀 Next Steps

Ready for Week 2? Your chatbot will become the foundation for building RAG-powered knowledge agents that can:
- Answer questions from custom documents
- Retrieve relevant information from vector databases
- Process and understand complex document collections

---

**Happy coding! 🤖✨**
