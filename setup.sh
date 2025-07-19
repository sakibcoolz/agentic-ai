#!/bin/bash

# Agentic AI Learning Repository Setup Script

echo "🚀 Setting up Agentic AI Learning Repository"
echo "==========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Found Go version: $GO_VERSION"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "⚠️  Please edit .env and add your API keys!"
else
    echo "✅ .env file already exists"
fi

# Initialize Go modules and install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi

# Create additional directories
echo "📁 Creating project directories..."
mkdir -p logs tmp data

echo ""
echo "🎉 Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit .env file and add your OpenAI API key"
echo "2. Start with Day 1: cd day-01-setup && go run main.go"
echo "3. Follow the README.md for the complete learning path"
echo ""
echo "Happy learning! 🤖"
