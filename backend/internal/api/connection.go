package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gibranda/chat-with-database/internal/agent"
	"github.com/gibranda/chat-with-database/internal/database"
)

type ConnectionRequest struct {
	Type     string `json:"type" binding:"required"`     // postgres, mysql, sqlite3
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database" binding:"required"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  string `json:"sslmode"`
	Path     string `json:"path"` // for SQLite
}

type ConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Tables  int    `json:"tables,omitempty"`
}

func (h *Handler) TestConnection(c *gin.Context) {
	var req ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build connection string
	connStr := buildConnectionString(req)

	// Try to connect
	testDB, err := database.New(req.Type, connStr)
	if err != nil {
		c.JSON(http.StatusOK, ConnectionResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	defer testDB.Close()

	// Get table count to verify connection
	tables, err := testDB.GetTables()
	if err != nil {
		c.JSON(http.StatusOK, ConnectionResponse{
			Success: false,
			Message: "Connected but failed to read schema: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ConnectionResponse{
		Success: true,
		Message: "Connection successful",
		Tables:  len(tables),
	})
}

func (h *Handler) Connect(c *gin.Context) {
	var req ConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build connection string
	connStr := buildConnectionString(req)

	// Connect to new database
	newDB, err := database.New(req.Type, connStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Close old database connection
	if h.db != nil {
		h.db.Close()
	}

	// Update handler with new connection
	h.db = newDB

	// Initialize agent with new database
	h.agent = agent.NewAgent(
		h.llmClient,
		newDB,
		h.config.Agent.MaxIterations,
		h.config.Agent.EnableQueryValidation,
		h.config.Agent.ReadonlyMode,
		h.config.Agent.MaxResults,
	)

	log.Printf("✓ Connected to %s database: %s", req.Type, req.Database)

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Connected successfully",
		"type":     req.Type,
		"database": req.Database,
	})
}

func buildConnectionString(req ConnectionRequest) string {
	switch req.Type {
	case "postgres":
		sslmode := req.SSLMode
		if sslmode == "" {
			sslmode = "disable"
		}
		return database.BuildPostgresConnString(
			req.Host,
			req.Port,
			req.Database,
			req.User,
			req.Password,
			sslmode,
		)
	case "mysql":
		return database.BuildMySQLConnString(
			req.Host,
			req.Port,
			req.Database,
			req.User,
			req.Password,
		)
	case "sqlite3":
		return req.Path
	default:
		return ""
	}
}
