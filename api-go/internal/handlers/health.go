package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/gin-gonic/gin"
)

// Root handles root route requests
func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "ZeroTrace API",
		"version":     "1.0.0",
		"status":      "running",
		"description": "ZeroTrace Security Platform API",
		"endpoints": gin.H{
			"health":          "/health",
			"agents":          "/api/agents",
			"vulnerabilities": "/api/vulnerabilities",
			"dashboard":       "/api/dashboard/overview",
		},
		"frontend":  "http://localhost:5173",
		"timestamp": time.Now(),
	})
}

// HealthCheck handles health check requests
func HealthCheck(db *repository.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbStatus := "unknown"

		// Check database health
		if sqlDB, err := db.DB.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				dbStatus = "healthy"
			} else {
				dbStatus = "unhealthy"
			}
		} else {
			dbStatus = "error"
		}

		// Check Redis health
		redisStatus := "unknown"
		if db.Redis != nil {
			if err := db.Redis.Ping(c.Request.Context()).Err(); err == nil {
				redisStatus = "healthy"
			} else {
				redisStatus = "unhealthy"
			}
		} else {
			redisStatus = "disabled"
		}

		response := models.HealthCheckResponse{
			Status:    "healthy",
			Timestamp: time.Now(),
			Services: map[string]string{
				"api":      "healthy",
				"database": dbStatus,
				"redis":    redisStatus,
			},
		}

		// If DB or Redis is down, overall status is degraded
		if dbStatus != "healthy" || (redisStatus != "healthy" && redisStatus != "disabled") {
			response.Status = "degraded"
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      response,
			Message:   "Service status retrieved",
			Timestamp: time.Now(),
		})
	}
}
