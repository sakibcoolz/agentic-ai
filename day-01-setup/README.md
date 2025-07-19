# Day 1: Go Setup & AI Fundamentals

Welcome to Day 1 of your Agentic AI journey! Today we'll set up your development environment and understand the core concepts of AI agents.

## 🎯 Learning Goals

- Set up Go development environment for AI projects
- Understand what Agentic AI means
- Learn the difference between traditional AI and Agent-based AI
- Set up basic project structure
- Create your first simple AI interaction

## 📖 Theory: What is Agentic AI?

### Traditional AI vs Agentic AI

**Traditional AI**: 
- Input → Model → Output
- Stateless interactions
- Single-turn responses

**Agentic AI**:
- Can use tools and external resources
- Maintains state and memory
- Can plan and execute multi-step tasks
- Can interact with external systems
- Can learn and adapt over time

### Key Components of AI Agents

1. **Perception**: Understanding input and context
2. **Planning**: Deciding what actions to take
3. **Action**: Executing tools and functions
4. **Memory**: Maintaining context across interactions
5. **Learning**: Improving performance over time

## 🛠 Setup Instructions

### 1. Verify Go Installation

```bash
go version
```

You should see Go 1.21 or higher.

### 2. Set up the project

```bash
cd /Users/sakibmulla/Documents/Development/learnings/agentic-ai
cp .env.example .env
```

### 3. Edit your .env file

Add your OpenAI API key:
```
OPENAI_API_KEY=sk-your-actual-api-key-here
```

### 4. Install dependencies

```bash
go mod tidy
```

## 💻 Practical Exercise

Let's create a simple AI client to test our setup:

### Step 1: Create the basic client structure

Run the example:
```bash
go run day-01-setup/main.go
```

### Step 2: Test different prompts

Try these prompts:
- "What is artificial intelligence?"
- "Explain the concept of agents in AI"
- "How do you build software agents?"

## 🧪 Hands-on Lab

### Lab 1: Environment Validation

Create a simple health check system that validates all your API keys and connections.

### Lab 2: Basic Chat Interface

Build a command-line chat interface that maintains conversation history.

## 📝 Key Concepts Learned

1. **Environment Setup**: Proper project structure for AI applications
2. **API Integration**: How to connect to LLM services
3. **Error Handling**: Graceful handling of API failures
4. **Basic Conversation**: Simple request-response patterns

## 🔄 Next Steps

Tomorrow we'll dive deeper into LLM integration and learn about:
- Different LLM providers and APIs
- Token management and cost optimization
- Response streaming and async operations
- Model selection strategies

## 📚 Additional Reading

- [OpenAI API Documentation](https://platform.openai.com/docs)
- [Go Best Practices](https://golang.org/doc/effective_go.html)
- [Introduction to AI Agents](https://en.wikipedia.org/wiki/Software_agent)

## 🤔 Reflection Questions

1. What makes an AI system "agentic"?
2. How might agents be different from traditional chatbots?
3. What are some real-world applications where agentic AI would be beneficial?

---

**Congratulations! You've completed Day 1. Tomorrow we'll start building more sophisticated LLM integrations! 🚀**
