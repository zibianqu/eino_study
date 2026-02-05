package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/internal/app/handler"
	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/app/service"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize repositories
	db := database.GetDB()
	docRepo := repository.NewDocumentRepository(db)
	chunkRepo := repository.NewChunkRepository(db)
	entityRepo := repository.NewEntityRepository(db)

	// Initialize services
	docService := service.NewDocumentService(docRepo, chunkRepo, entityRepo)
	ragService := service.NewRAGService(chunkRepo, docRepo)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	docHandler := handler.NewDocumentHandler(docService)
	queryHandler := handler.NewQueryHandler(ragService)

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