# Quick Start Guide - Day 7 Chatbot

## 🚀 5-Minute Setup

### Prerequisites
- Go 1.21+ installed
- OpenAI API key

### Step 1: Setup
```bash
cd day-07-chatbot-project
./setup.sh
```

### Step 2: Configure
Edit `.env` file and add your OpenAI API key:
```bash
OPENAI_API_KEY=your-actual-api-key-here
```

### Step 3: Run
```bash
go run main.go
```

### Step 4: Chat!
```
You: Hello!
Bot: Hello! How can I help you today?

You: Tell me a joke
Bot: Why do programmers prefer dark mode? Because light attracts bugs! 🐛

You: /mode creative
Bot: Switched to creative mode! 🎭

You: Write a haiku about coffee
Bot: Steam rises gently,
     Dark brew awakens the mind,
     Morning's first blessing.
```

## 🎮 Essential Commands

| Command | Description | Example |
|---------|-------------|---------|
| `help` | Show all commands | `help` |
| `/mode <type>` | Change conversation style | `/mode creative` |
| `/save <name>` | Save current chat | `/save my-coding-chat` |
| `/load <name>` | Load saved chat | `/load my-coding-chat` |
| `/clear` | Clear memory | `/clear` |
| `/history` | List saved chats | `/history` |
| `/stats` | Show usage stats | `/stats` |
| `quit` | Exit chatbot | `quit` |

## 🎭 Conversation Modes

- **assistant** (default): Professional, helpful responses
- **casual**: Relaxed, friendly conversation  
- **creative**: Imaginative, artistic responses

## 💡 Pro Tips

1. **Memory Management**: Bot remembers last 10 message pairs
2. **Save Important Chats**: Use `/save` for conversations you want to revisit
3. **Switch Modes**: Try different modes for different types of conversations
4. **Check Stats**: Use `/stats` to track token usage

## 🐛 Troubleshooting

**Bot won't start?**
- Check your API key in `.env`
- Run `go mod tidy` to install dependencies

**Getting rate limited?**
- The bot has automatic retry with backoff
- Consider upgrading your OpenAI plan

**Responses seem slow?**
- Check your internet connection
- Try reducing `MAX_TOKENS` in `.env`

## 📈 What's Next?

Once you're comfortable with the basic chatbot:
1. Complete the labs in `LABS.md`
2. Customize conversation modes
3. Add new features
4. Prepare for Week 2's RAG capabilities!

---

**Happy chatting! 🤖**
