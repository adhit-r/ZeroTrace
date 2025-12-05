package middleware

import (
	"net/http"
	"strings"
	"time"

	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InputValidationMiddleware validates and sanitizes input
func InputValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate UUID parameters
		if id := c.Param("id"); id != "" {
			if _, err := uuid.Parse(id); err != nil {
				c.JSON(http.StatusBadRequest, models.APIResponse{
					Success: false,
					Error: &models.APIError{
						Code:    "INVALID_UUID",
						Message: "Invalid UUID format",
						Details: err.Error(),
					},
					Timestamp: time.Now(),
				})
				c.Abort()
				return
			}
		}

		// Validate query parameters
		for _, values := range c.Request.URL.Query() {
			for _, value := range values {
				// Basic XSS prevention - remove script tags
				if strings.Contains(strings.ToLower(value), "<script") {
					c.JSON(http.StatusBadRequest, models.APIResponse{
						Success: false,
						Error: &models.APIError{
							Code:    "INVALID_INPUT",
							Message: "Invalid characters in query parameter",
						},
						Timestamp: time.Now(),
					})
					c.Abort()
					return
				}
			}
		}

		// Limit request body size (handled by Gin's MaxMultipartMemory, but we can add additional checks)
		if c.Request.ContentLength > 10*1024*1024 { // 10MB limit
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "REQUEST_TOO_LARGE",
					Message: "Request body exceeds maximum size of 10MB",
				},
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UUIDValidation validates UUID path parameters
func UUIDValidation(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(paramName)
		if id == "" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "MISSING_UUID",
					Message: "UUID parameter is required",
				},
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		if _, err := uuid.Parse(id); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_UUID",
					Message: "Invalid UUID format",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Remove control characters except newlines and tabs
	var result strings.Builder
	for _, r := range input {
		if r >= 32 || r == '\n' || r == '\t' {
			result.WriteRune(r)
		}
	}
	
	return strings.TrimSpace(result.String())
}

// ValidateContentType validates Content-Type header
func ValidateContentType(allowedTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Next()
			return
		}

		contentType := c.GetHeader("Content-Type")
		if contentType == "" {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "MISSING_CONTENT_TYPE",
					Message: "Content-Type header is required",
				},
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		// Check if content type is allowed
		allowed := false
		for _, allowedType := range allowedTypes {
			if strings.HasPrefix(contentType, allowedType) {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_CONTENT_TYPE",
					Message: "Content-Type must be one of: " + strings.Join(allowedTypes, ", "),
				},
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

