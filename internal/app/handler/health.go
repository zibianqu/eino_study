package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zibianqu/eino_study/internal/pkg/database"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	// Check database connection
	db := database.GetDB()
	if db == nil {
		InternalError(c, "database not initialized")
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		InternalError(c, "failed to get database instance")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		InternalError(c, "database connection failed")
		return
	}

	Success(c, gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}