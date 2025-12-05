package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ETagMiddleware adds ETag support for HTTP caching
func ETagMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to GET and HEAD requests
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.Next()
			return
		}

		// Get If-None-Match header
		ifNoneMatch := c.GetHeader("If-None-Match")

		// Create response recorder to capture body
		recorder := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           make([]byte, 0),
		}
		c.Writer = recorder

		c.Next()

		// Generate ETag from response body
		etag := generateETag(recorder.body)

		// Set ETag header
		c.Header("ETag", etag)

		// Check if client has matching ETag
		if ifNoneMatch != "" && ifNoneMatch == etag {
			c.Status(http.StatusNotModified)
			c.Writer = recorder.ResponseWriter
			return
		}

		// Write the captured body
		c.Writer = recorder.ResponseWriter
		if len(recorder.body) > 0 {
			c.Writer.Write(recorder.body)
		}
	}
}

// responseRecorder captures response body for ETag generation
type responseRecorder struct {
	gin.ResponseWriter
	body []byte
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}

func (r *responseRecorder) WriteString(s string) (int, error) {
	r.body = append(r.body, []byte(s)...)
	return len(s), nil
}

// generateETag generates ETag from content
func generateETag(content []byte) string {
	hash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf(`"%s"`, hashStr[:16]) // Use first 16 chars
}

