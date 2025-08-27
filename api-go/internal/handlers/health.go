package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
)

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
