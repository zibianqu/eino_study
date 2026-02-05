package main

import (
	"fmt"
	"log"

	"github.com/zibianqu/eino_study/internal/app/router"
	"github.com/zibianqu/eino_study/internal/config"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database connected successfully")

	// Setup router
	r := router.SetupRouter()

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}