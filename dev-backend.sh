#!/bin/bash

echo "🔥 Starting AI Database Agent Backend in Development Mode..."
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

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "📦 Installing air for hot reloading..."
    go install github.com/air-verse/air@latest
    echo ""
fi

echo "🔥 Starting Go backend server with hot reloading..."
echo "   Files will auto-reload on changes in:"
echo "   - *.go files"
echo "   - *.yaml files"
echo ""
echo "Press Ctrl+C to stop"
echo ""

air
