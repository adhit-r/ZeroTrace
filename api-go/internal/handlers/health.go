package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"

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
func HealthCheck(c *gin.Context) {
	response := models.HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Services: map[string]string{
			"api":      "healthy",
			"database": "healthy", // TODO: Add actual DB health check
			"redis":    "healthy", // TODO: Add actual Redis health check
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      response,
		Message:   "Service is healthy",
		Timestamp: time.Now(),
	})
}
