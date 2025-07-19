# Day 5 Quick Start Guide

## 🚀 Running the Memory System

### Option 1: Full Memory Manager
```bash
cd day-05-context-memory
go run main.go
```

### Option 2: Memory Concepts Demo
```bash
go run demo.go
```

## 💬 Try These Conversations

### Test 1: Personal Information Learning
```
You: Hi, my name is Sarah and I work as a data scientist
AI: [Responds and learns your name and profession]

You: I really enjoy working with Python and machine learning
AI: [Learns your preferences]

You: What do you know about me so far?
AI: [Recalls learned facts about you]
```

### Test 2: Long Conversation with Summarization
Start a conversation and keep it going with 25+ messages to see automatic summarization in action.

### Test 3: Context Continuity
```
You: I'm working on a Go project for web development
AI: [Provides relevant help]

[Several messages later...]
You: Can you remind me about that project we discussed?
AI: [References the Go web development project]
```

## 🔧 Key Commands

- `stats` - View memory statistics
- `facts` - See learned facts about you
- `clear` - Reset memory
- `quit` - Exit

## 📊 What to Observe

1. **Fact Learning**: Watch as the system extracts and stores facts about you
2. **Token Management**: Notice how context is optimized for token limits
3. **Summarization**: See conversation compression after 20+ messages
4. **Personalization**: Observe how responses become more tailored to you

## 🎯 Learning Objectives Achieved

- ✅ **Conversation Memory**: Maintains context across interactions
- ✅ **User Learning**: Automatically extracts and remembers user facts
- ✅ **Context Optimization**: Manages token budgets efficiently
- ✅ **Progressive Summarization**: Compresses long conversations
- ✅ **Persistent Storage**: Remembers information across sessions

## 🔄 Next Steps

Ready for Day 6? We'll focus on making these systems production-ready with robust error handling and monitoring!
