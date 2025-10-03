package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gibranda/chat-with-database/internal/agent"
	"github.com/gibranda/chat-with-database/internal/config"
	"github.com/gibranda/chat-with-database/internal/database"
	"github.com/gibranda/chat-with-database/internal/llm"
)

type Handler struct {
	agent     *agent.Agent
	db        *database.Database
	llmClient *llm.OllamaClient
	config    *config.Config
}

type QueryRequest struct {
	Question string `json:"question" binding:"required"`
}

type SchemaResponse struct {
	Schema *database.SchemaInfo `json:"schema"`
}

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	LLM      string `json:"llm"`
}

func NewHandler(agentInstance *agent.Agent, db *database.Database, llmClient *llm.OllamaClient, cfg *config.Config) *Handler {
	return &Handler{
		agent:     agentInstance,
		db:        db,
		llmClient: llmClient,
		config:    cfg,
	}
}

func (h *Handler) Health(c *gin.Context) {
	dbStatus := "not_connected"
	if h.db != nil {
		dbStatus = "connected"
	}
	
	c.JSON(http.StatusOK, HealthResponse{
		Status:   "healthy",
		Database: dbStatus,
		LLM:      "connected",
	})
}

func (h *Handler) Query(c *gin.Context) {
	// Check if database is connected
	if h.db == nil || h.agent == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Database not connected. Please connect to a database first.",
		})
		return
	}

	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Processing query: %s", req.Question)

	response, err := h.agent.ProcessQuery(req.Question)
	if err != nil {
		log.Printf("Error processing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Query processed successfully")
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSchema(c *gin.Context) {
	schema, err := h.agent.GetSchema()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, SchemaResponse{Schema: schema})
}

func (h *Handler) RefreshSchema(c *gin.Context) {
	if err := h.agent.RefreshSchema(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schema refreshed successfully"})
}

func (h *Handler) ClearHistory(c *gin.Context) {
	h.agent.ClearHistory()
	c.JSON(http.StatusOK, gin.H{"message": "Conversation history cleared"})
}

func (h *Handler) GetTables(c *gin.Context) {
	tables, err := h.db.GetTables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tables": tables})
}

func (h *Handler) GetTableInfo(c *gin.Context) {
	tableName := c.Param("table")
	if tableName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "table name is required"})
		return
	}

	info, err := h.db.GetTableInfo(tableName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}
