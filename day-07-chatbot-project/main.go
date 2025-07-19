package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"chatbot/chatbot"
	"chatbot/config"
	"chatbot/llm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize LLM client
	llmClient, err := llm.NewClient(cfg.OpenAIAPIKey, cfg.Model)
	if err != nil {
		fmt.Printf("Error initializing LLM client: %v\n", err)
		os.Exit(1)
	}

	// Initialize chatbot
	bot, err := chatbot.New(llmClient, cfg)
	if err != nil {
		fmt.Printf("Error initializing chatbot: %v\n", err)
		os.Exit(1)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down gracefully...")
		cancel()
	}()

	// Start the chat loop
	if err := runChatLoop(ctx, bot); err != nil {
		fmt.Printf("Chat loop error: %v\n", err)
		os.Exit(1)
	}
}

func runChatLoop(ctx context.Context, bot *chatbot.Bot) error {
	scanner := bufio.NewScanner(os.Stdin)

	// Print welcome message
	fmt.Println("🤖 Welcome to the Simple Chatbot!")
	fmt.Println("Type 'help' for commands, 'quit' to exit.")
	fmt.Println("Available modes: casual, assistant, creative")
	fmt.Println(strings.Repeat("-", 50))

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Print("\nYou: ")
			if !scanner.Scan() {
				return scanner.Err()
			}

			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}

			// Handle special commands
			if handled, err := handleCommand(input, bot); err != nil {
				fmt.Printf("Command error: %v\n", err)
				continue
			} else if handled {
				continue
			}

			// Get bot response
			response, err := bot.ProcessMessage(ctx, input)
			if err != nil {
				fmt.Printf("Bot error: %v\n", err)
				continue
			}

			fmt.Printf("Bot: %s\n", response)
		}
	}
}

func handleCommand(input string, bot *chatbot.Bot) (bool, error) {
	if !strings.HasPrefix(input, "/") && input != "help" && input != "quit" {
		return false, nil
	}

	switch {
	case input == "quit" || input == "/quit":
		fmt.Println("Goodbye! 👋")
		os.Exit(0)
		return true, nil

	case input == "help" || input == "/help":
		printHelp()
		return true, nil

	case strings.HasPrefix(input, "/mode "):
		mode := strings.TrimPrefix(input, "/mode ")
		if err := bot.SetMode(mode); err != nil {
			return true, err
		}
		fmt.Printf("Switched to %s mode! 🎭\n", mode)
		return true, nil

	case input == "/clear":
		bot.ClearMemory()
		fmt.Println("Conversation memory cleared! 🧹")
		return true, nil

	case strings.HasPrefix(input, "/save "):
		name := strings.TrimPrefix(input, "/save ")
		if err := bot.SaveConversation(name); err != nil {
			return true, err
		}
		fmt.Printf("Conversation saved as '%s' 💾\n", name)
		return true, nil

	case strings.HasPrefix(input, "/load "):
		name := strings.TrimPrefix(input, "/load ")
		if err := bot.LoadConversation(name); err != nil {
			return true, err
		}
		fmt.Printf("Conversation '%s' loaded! 📂\n", name)
		return true, nil

	case input == "/history":
		conversations := bot.ListConversations()
		if len(conversations) == 0 {
			fmt.Println("No saved conversations found.")
		} else {
			fmt.Println("Saved conversations:")
			for _, conv := range conversations {
				fmt.Printf("  - %s\n", conv)
			}
		}
		return true, nil

	case input == "/stats":
		stats := bot.GetStats()
		fmt.Printf("Session stats:\n")
		fmt.Printf("  Messages: %d\n", stats.MessageCount)
		fmt.Printf("  Tokens used: %d\n", stats.TokensUsed)
		fmt.Printf("  Current mode: %s\n", stats.CurrentMode)
		return true, nil

	default:
		fmt.Printf("Unknown command: %s\n", input)
		return true, nil
	}
}

func printHelp() {
	fmt.Println("\n📚 Available Commands:")
	fmt.Println("  help                 - Show this help message")
	fmt.Println("  quit                 - Exit the chatbot")
	fmt.Println("  /mode <mode>         - Change conversation mode (casual/assistant/creative)")
	fmt.Println("  /clear               - Clear conversation memory")
	fmt.Println("  /save <name>         - Save current conversation")
	fmt.Println("  /load <name>         - Load a saved conversation")
	fmt.Println("  /history             - List saved conversations")
	fmt.Println("  /stats               - Show session statistics")
	fmt.Println("\n💡 Tips:")
	fmt.Println("  - The bot remembers your conversation within the session")
	fmt.Println("  - Try different modes for different conversation styles")
	fmt.Println("  - Save important conversations for later reference")
}
