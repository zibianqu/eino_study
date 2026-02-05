package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/internal/app/handler"
	"github.com/zibianqu/eino_study/internal/app/service"
)

// SetupRouter sets up the Gin router with all routes
func SetupRouter(services *service.ServiceContainer) *gin.Engine {
	r := gin.Default()

	// Add CORS middleware if needed
	r.Use(corsMiddleware())

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	docHandler := handler.NewDocumentHandler(services.DocumentService)
	queryHandler := handler.NewQueryHandler(services.RAGService)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", healthHandler.Check)

		// Document management
		docs := v1.Group("/documents")
		{
			docs.POST("", docHandler.Upload)
			docs.GET("", docHandler.List)
			docs.GET("/:id", docHandler.Get)
			docs.DELETE("/:id", docHandler.Delete)
			docs.POST("/:id/process", docHandler.Process)
		}

		// Query
		v1.POST("/query", queryHandler.Query)
	}

	return r
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}