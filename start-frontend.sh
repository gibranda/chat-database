#!/bin/bash

echo "🚀 Starting AI Database Agent Frontend..."
echo ""

cd frontend

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "⚠️  Dependencies not installed. Run setup.sh first."
    exit 1
fi

# Check if backend is running
if ! curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "⚠️  Warning: Backend doesn't seem to be running on localhost:8080"
    echo "   Please start the backend first: ./start-backend.sh"
    echo ""
fi

echo "Starting Nuxt.js frontend..."
npm run dev
