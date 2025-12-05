package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger middleware logs HTTP requests with correlation ID
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Get correlation ID from context if available
		correlationID := "unknown"
		if param.Keys != nil {
			if id, exists := param.Keys[CorrelationIDKey]; exists {
				if str, ok := id.(string); ok {
					correlationID = str
				}
			}
		}

		return fmt.Sprintf("[%s] %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			correlationID,
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
