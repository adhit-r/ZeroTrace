package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// CorrelationIDHeader is the HTTP header name for correlation ID
	CorrelationIDHeader = "X-Correlation-ID"
	// CorrelationIDKey is the context key for correlation ID
	CorrelationIDKey = "correlation_id"
)

// CorrelationID middleware adds a correlation ID to each request
func CorrelationID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get correlation ID from header or generate a new one
		correlationID := c.GetHeader(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
		}

		// Set in context for use in handlers
		c.Set(CorrelationIDKey, correlationID)

		// Set in response header
		c.Header(CorrelationIDHeader, correlationID)

		c.Next()
	}
}

// GetCorrelationID retrieves the correlation ID from context
func GetCorrelationID(c *gin.Context) string {
	if id, exists := c.Get(CorrelationIDKey); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

