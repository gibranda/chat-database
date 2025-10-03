package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gibranda/chat-with-database/internal/api"
	"github.com/gibranda/chat-with-database/internal/config"
	"github.com/gibranda/chat-with-database/internal/llm"
)

func main() {
	// Load configuration
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize LLM client
	log.Printf("Connecting to Ollama at %s...\n", cfg.Ollama.Host)
	llmClient := llm.NewOllamaClient(
		cfg.Ollama.Host,
		cfg.Ollama.Model,
		cfg.Ollama.Temperature,
		cfg.Ollama.Timeout,
	)
	log.Printf("✓ Ollama client initialized (model: %s)\n", cfg.Ollama.Model)

	// Setup API without database connection
	// Database will be connected when user provides credentials from UI
	log.Println("⏳ Waiting for database connection from UI...")
	handler := api.NewHandler(nil, nil, llmClient, cfg)
	router := api.SetupRouter(handler, cfg.Server.Debug)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("\n🚀 Server starting on %s\n", addr)
	log.Printf("🤖 LLM Model: %s\n", cfg.Ollama.Model)
	log.Printf("📊 Database: Not connected (waiting for UI input)\n")
	log.Println("\nEndpoints:")
	log.Println("  POST   /api/connection/test    - Test database connection")
	log.Println("  POST   /api/connection/connect - Connect to database")
	log.Println("  POST   /api/query              - Ask questions in natural language")
	log.Println("  GET    /api/schema             - Get database schema")
	log.Println("  GET    /api/health             - Health check")
	log.Println()

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
