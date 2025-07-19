# Day 7 Labs: Simple Chatbot Project

## Lab Overview

These hands-on labs will guide you through building and extending the simple chatbot step by step. Each lab builds upon the previous one, introducing new concepts and features.

## 🧪 Lab 1: Basic Setup and First Chat

### Objective
Get the chatbot running and have your first conversation.

### Tasks
1. **Setup the environment:**
   ```bash
   cd day-07-chatbot-project
   ./setup.sh
   ```

2. **Configure your API key:**
   - Edit `.env` file
   - Add your OpenAI API key
   - Adjust other settings if needed

3. **Build and run:**
   ```bash
   go run main.go
   ```

4. **Test basic functionality:**
   - Start a conversation
   - Try different message types
   - Test the help command
   - Exit gracefully with `quit`

### Expected Output
```
🤖 Welcome to the Simple Chatbot!
Type 'help' for commands, 'quit' to exit.
Available modes: casual, assistant, creative
--------------------------------------------------

You: Hello!
Bot: Hello! How can I help you today?

You: What's the weather like?
Bot: I don't have access to real-time weather data, but I can help you think about weather-related topics or suggest ways to check the weather.
```

### Success Criteria
- ✅ Chatbot starts without errors
- ✅ Bot responds to messages
- ✅ Help command works
- ✅ Can exit gracefully

---

## 🧪 Lab 2: Conversation Modes and Memory

### Objective
Test different conversation modes and memory functionality.

### Tasks
1. **Test conversation modes:**
   ```
   You: /mode creative
   You: Write a short poem about coding
   
   You: /mode casual  
   You: How's it going?
   
   You: /mode assistant
   You: Can you help me understand arrays in Go?
   ```

2. **Test memory functionality:**
   - Have a multi-turn conversation
   - Reference previous messages
   - Check if bot remembers context

3. **Test memory limits:**
   - Have a long conversation (20+ messages)
   - Verify older messages are forgotten
   - Check that context is maintained for recent messages

### Expected Behavior
- Different modes should have distinct personalities
- Bot should remember recent conversation context
- Older messages should be pruned when limit is reached

### Success Criteria
- ✅ All three modes work differently
- ✅ Bot remembers recent context
- ✅ Memory management works correctly

---

## 🧪 Lab 3: Conversation Persistence

### Objective
Test saving and loading conversations.

### Tasks
1. **Save a conversation:**
   ```
   You: Let's talk about Go programming
   Bot: [response]
   You: What are goroutines?
   Bot: [response]
   You: /save golang-discussion
   ```

2. **Start a new topic and save it:**
   ```
   You: /clear
   You: I want to learn about cooking
   Bot: [response]
   You: /save cooking-chat
   ```

3. **List and load conversations:**
   ```
   You: /history
   You: /load golang-discussion
   You: Can you continue our discussion about goroutines?
   ```

4. **Verify persistence:**
   - Exit the chatbot
   - Restart it
   - Load a saved conversation
   - Verify context is restored

### Success Criteria
- ✅ Conversations save successfully
- ✅ History command lists saved conversations
- ✅ Loading restores conversation context
- ✅ Persistence survives app restart

---

## 🧪 Lab 4: Error Handling and Recovery

### Objective
Test error handling and retry mechanisms.

### Tasks
1. **Test with invalid API key:**
   - Temporarily change API key in .env
   - Start chatbot
   - Observe error handling

2. **Test network issues simulation:**
   - Use a rate-limited API key (if available)
   - Send multiple rapid requests
   - Observe retry behavior

3. **Test graceful degradation:**
   - Try very long messages
   - Test with special characters
   - Test empty messages

### Expected Behavior
- Clear error messages for configuration issues
- Automatic retries for transient errors
- Graceful handling of edge cases

### Success Criteria
- ✅ Clear error messages displayed
- ✅ Retries work for transient failures
- ✅ App doesn't crash on errors

---

## 🧪 Lab 5: Customization and Extensions

### Objective
Customize the chatbot and add new features.

### Tasks

#### 5.1: Custom Personality Mode
Add a new conversation mode called "philosopher":

1. **Edit `llm/prompts.go`:**
   ```go
   "philosopher": `You are a thoughtful philosopher. Respond with deep insights,
   thought-provoking questions, and references to philosophical concepts.
   Encourage reflection and critical thinking.`,
   ```

2. **Test the new mode:**
   ```
   You: /mode philosopher
   You: What is the meaning of life?
   ```

#### 5.2: Token Usage Tracking
Enhance the stats command to show more information:

1. **Test current stats:**
   ```
   You: /stats
   ```

2. **Have a longer conversation and check stats again**

#### 5.3: Custom Commands
Add a new command `/export` that saves conversation in a readable format:

1. **Implement in `main.go`:**
   - Add case for `/export <filename>`
   - Save conversation as formatted text

2. **Test the feature:**
   ```
   You: /export my-conversation.txt
   ```

### Success Criteria
- ✅ New philosopher mode works
- ✅ Stats tracking is accurate
- ✅ Export functionality works

---

## 🧪 Lab 6: Advanced Features

### Objective
Implement advanced chatbot features.

### Tasks

#### 6.1: Conversation Search
Add ability to search through saved conversations:

1. **Add search functionality to history module**
2. **Implement `/search <term>` command**
3. **Test searching across multiple conversations**

#### 6.2: Conversation Metadata
Add metadata tracking:

1. **Track conversation statistics:**
   - Message count
   - Total tokens used
   - Duration

2. **Display in history listing:**
   ```
   You: /history
   Saved conversations:
     - golang-discussion (12 messages, 450 tokens)
     - cooking-chat (5 messages, 200 tokens)
   ```

#### 6.3: Rate Limiting Protection
Add intelligent rate limiting:

1. **Track API usage per minute**
2. **Automatically adjust request timing**
3. **Show rate limit warnings**

### Success Criteria
- ✅ Search finds relevant conversations
- ✅ Metadata is tracked and displayed
- ✅ Rate limiting prevents API errors

---

## 🧪 Lab 7: Performance and Testing

### Objective
Optimize performance and add comprehensive testing.

### Tasks

#### 7.1: Performance Testing
1. **Measure response times:**
   - Add timing to bot responses
   - Track average response time

2. **Memory usage optimization:**
   - Test with long conversations
   - Monitor memory usage
   - Optimize message storage

#### 7.2: Unit Testing
1. **Run existing tests:**
   ```bash
   go test ./...
   ```

2. **Add new test cases:**
   - Test conversation export
   - Test search functionality
   - Test custom modes

#### 7.3: Integration Testing
1. **Test complete workflows:**
   - End-to-end conversation flows
   - Save/load/search sequences
   - Error recovery scenarios

2. **Load testing:**
   - Multiple rapid conversations
   - Large conversation loading
   - Concurrent usage simulation

### Success Criteria
- ✅ All tests pass
- ✅ Performance is acceptable
- ✅ Memory usage is optimized

---

## 🎯 Challenge Labs

### Challenge 1: Web Interface
Create a simple web interface for the chatbot:
- Use Go's `net/http` package
- Create a basic HTML chat interface
- Support all existing commands

### Challenge 2: Multi-User Support
Extend the chatbot to support multiple users:
- Add user identification
- Separate conversation histories
- User-specific settings

### Challenge 3: Plugin System
Create a plugin architecture:
- Define plugin interface
- Implement sample plugins (weather, calculator)
- Dynamic plugin loading

### Challenge 4: Voice Integration
Add voice capabilities:
- Text-to-speech for bot responses
- Speech-to-text for user input
- Voice mode toggle

---

## 🏆 Lab Completion Checklist

### Basic Functionality
- [ ] Chatbot starts and responds correctly
- [ ] All three conversation modes work
- [ ] Memory management functions properly
- [ ] Conversation persistence works
- [ ] Error handling is robust

### Advanced Features
- [ ] Custom conversation modes
- [ ] Enhanced statistics tracking
- [ ] Conversation search capability
- [ ] Rate limiting protection
- [ ] Performance optimization

### Quality Assurance
- [ ] All tests pass
- [ ] Code is well-documented
- [ ] Error messages are user-friendly
- [ ] Performance is acceptable

### Optional Challenges
- [ ] Web interface implemented
- [ ] Multi-user support added
- [ ] Plugin system created
- [ ] Voice integration added

---

## 📚 Learning Takeaways

After completing these labs, you should understand:

1. **Chatbot Architecture**: How to structure a conversational AI application
2. **State Management**: Managing conversation context and memory
3. **API Integration**: Working with LLM APIs effectively
4. **Error Handling**: Building resilient applications
5. **Persistence**: Saving and loading application state
6. **Testing**: Comprehensive testing strategies
7. **Performance**: Optimizing for real-world usage

## 🚀 Next Steps

Ready for Week 2? Your chatbot foundation will support:
- RAG (Retrieval-Augmented Generation)
- Document processing and search
- Vector database integration
- Knowledge base querying
- Multi-modal interactions

---

**Great job completing Day 7! Your chatbot is ready to evolve into a knowledge-powered agent! 🤖✨**
