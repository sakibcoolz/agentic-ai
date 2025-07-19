#!/bin/bash

# Day 7 Chatbot Project Setup Script

echo "🤖 Setting up Day 7 Chatbot Project..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Found Go version: $GO_VERSION"

# Create data directory
echo "📁 Creating data directory..."
mkdir -p ./data/conversations

# Copy environment file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from example..."
    cp .env.example .env
    echo "⚠️  Please edit .env and add your OpenAI API key!"
else
    echo "✅ .env file already exists"
fi

# Initialize go module and download dependencies
echo "📦 Installing dependencies..."
go mod tidy

# Build the project
echo "🔨 Building the project..."
if go build -o chatbot main.go; then
    echo "✅ Build successful!"
else
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env and add your OpenAI API key"
echo "2. Run the chatbot with: ./chatbot"
echo "   Or: go run main.go"
echo ""
echo "Commands to try:"
echo "  help                 - Show available commands"
echo "  /mode creative       - Switch to creative mode"
echo "  /save my-chat        - Save current conversation"
echo "  /clear               - Clear conversation memory"
echo "  quit                 - Exit the chatbot"
echo ""
echo "Happy chatting! 🚀"
