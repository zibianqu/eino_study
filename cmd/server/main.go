package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/internal/app/repository"
	"github.com/zibianqu/eino_study/internal/app/router"
	"github.com/zibianqu/eino_study/internal/app/service"
	"github.com/zibianqu/eino_study/internal/config"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("✓ Database connected successfully")

	// Initialize repositories
	db := database.GetDB()
	docRepo := repository.NewDocumentRepository(db)
	chunkRepo := repository.NewChunkRepository(db)
	entityRepo := repository.NewEntityRepository(db)

	// Initialize services with Eino components
	log.Println("Initializing Eino components...")
	services, err := service.InitServices(cfg, docRepo, chunkRepo, entityRepo)
	if err != nil {
		log.Fatalf("Failed to initialize services: %v", err)
	}
	log.Println("✓ Eino components initialized successfully")

	// Setup router
	r := router.SetupRouter(services)

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 Server starting on %s", addr)
	log.Printf("📖 API documentation: http://%s:%d/api/v1/health", cfg.Server.Host, cfg.Server.Port)
	
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}