#!/bin/bash

echo "🚀 Setting up AI Database Agent..."
echo ""

# Check prerequisites
echo "Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi
echo "✅ Go found: $(go version)"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or higher."
    exit 1
fi
echo "✅ Node.js found: $(node --version)"

# Check Ollama
if ! command -v ollama &> /dev/null; then
    echo "⚠️  Ollama is not installed. Please install Ollama from https://ollama.ai"
    echo "   After installation, run: ollama pull llama3.1"
else
    echo "✅ Ollama found"
fi

echo ""
echo "Setting up backend..."
cd backend

# Copy config if not exists
if [ ! -f config.yaml ]; then
    cp config.example.yaml config.yaml
    echo "✅ Created config.yaml (please edit with your database credentials)"
else
    echo "✅ config.yaml already exists"
fi

# Download Go dependencies
echo "📦 Downloading Go dependencies..."
go mod download
if [ $? -eq 0 ]; then
    echo "✅ Go dependencies installed"
else
    echo "❌ Failed to download Go dependencies"
    exit 1
fi

cd ..

echo ""
echo "Setting up frontend..."
cd frontend

# Copy env file
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env file"
else
    echo "✅ .env already exists"
fi

# Install dependencies
echo "📦 Installing Node.js dependencies..."
npm install
if [ $? -eq 0 ]; then
    echo "✅ Node.js dependencies installed"
else
    echo "❌ Failed to install Node.js dependencies"
    exit 1
fi

cd ..

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit backend/config.yaml with your database credentials"
echo "2. Make sure Ollama is running: ollama serve"
echo "3. Pull a model if needed: ollama pull llama3.1"
echo "4. Start the backend: cd backend && go run cmd/server/main.go"
echo "5. Start the frontend: cd frontend && npm run dev"
echo ""
echo "Happy querying! 🎉"
