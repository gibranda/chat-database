package api

import (
"github.com/gin-contrib/cors"
"github.com/gin-gonic/gin"
)

func SetupRouter(handler *Handler, debug bool) *gin.Engine {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		api.GET("/health", handler.Health)
		
		// Connection management
		api.POST("/connection/test", handler.TestConnection)
		api.POST("/connection/connect", handler.Connect)
		
		// Query and schema
		api.POST("/query", handler.Query)
		api.GET("/schema", handler.GetSchema)
		api.POST("/schema/refresh", handler.RefreshSchema)
		api.GET("/tables", handler.GetTables)
		api.GET("/tables/:table", handler.GetTableInfo)
		api.POST("/history/clear", handler.ClearHistory)
	}

	return router
}
