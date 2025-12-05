package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/middleware"
	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
)

// ErrorResponse creates a standardized error response
func ErrorResponse(c *gin.Context, statusCode int, code, message string, details interface{}) {
	correlationID := middleware.GetCorrelationID(c)
	
	response := models.APIResponse{
		Success: false,
		Error: &models.APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
		Timestamp: time.Now(),
	}

	// Add correlation ID to response if available
	if correlationID != "" {
		c.Header("X-Correlation-ID", correlationID)
	}

	c.JSON(statusCode, response)
}

// BadRequest creates a 400 Bad Request error response
func BadRequest(c *gin.Context, code, message string, details interface{}) {
	ErrorResponse(c, http.StatusBadRequest, code, message, details)
}

// Unauthorized creates a 401 Unauthorized error response
func Unauthorized(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusUnauthorized, code, message, nil)
}

// Forbidden creates a 403 Forbidden error response
func Forbidden(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusForbidden, code, message, nil)
}

// NotFound creates a 404 Not Found error response
func NotFound(c *gin.Context, code, message string) {
	ErrorResponse(c, http.StatusNotFound, code, message, nil)
}

// InternalServerError creates a 500 Internal Server Error response
func InternalServerError(c *gin.Context, code, message string, err error) {
	details := map[string]interface{}{}
	if err != nil {
		details["error"] = err.Error()
	}
	ErrorResponse(c, http.StatusInternalServerError, code, message, details)
}

// SuccessResponse creates a standardized success response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	correlationID := middleware.GetCorrelationID(c)
	
	response := models.APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}

	// Add correlation ID to response if available
	if correlationID != "" {
		c.Header("X-Correlation-ID", correlationID)
	}

	c.JSON(statusCode, response)
}

