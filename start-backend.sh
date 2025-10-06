#!/bin/bash

echo "🚀 Starting AI Database Agent Backend..."
echo ""

cd backend

# Check if config exists
if [ ! -f config.yaml ]; then
    echo "❌ config.yaml not found!"
    echo "Please copy config.example.yaml to config.yaml and configure it."
    exit 1
fi

# Check if Ollama is running
if ! curl -s http://localhost:11434/api/tags > /dev/null 2>&1; then
    echo "⚠️  Warning: Ollama doesn't seem to be running on localhost:11434"
    echo "   Please start Ollama: ollama serve"
    echo ""
fi

# Check if air is installed for hot reloading
if command -v air &> /dev/null; then
    echo "🔥 Starting Go backend server with hot reloading (air)..."
    echo "   Files will auto-reload on changes"
    echo ""
    air
else
    echo "ℹ️  Air not installed. Running without hot reloading..."
    echo "   To enable hot reloading, install air: go install github.com/air-verse/air@latest"
    echo ""
    echo "Starting Go backend server..."
    go run cmd/server/main.go
fi
