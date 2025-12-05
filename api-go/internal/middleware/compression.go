package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CompressionMiddleware adds gzip compression to responses
func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip compression for certain content types
		if shouldSkipCompression(c.Request) {
			c.Next()
			return
		}

		// Check if client accepts gzip
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Set response headers
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Create gzip writer
		gz := gzip.NewWriter(c.Writer)
		defer gz.Close()

		// Replace writer with gzip writer
		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			Writer:         gz,
		}

		c.Next()
	}
}

// gzipWriter wraps the response writer with gzip compression
type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.Writer.Write(data)
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.Writer.Write([]byte(s))
}

// shouldSkipCompression checks if compression should be skipped
func shouldSkipCompression(r *http.Request) bool {
	// Skip for already compressed content
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "image/") ||
		strings.Contains(contentType, "video/") ||
		strings.Contains(contentType, "application/zip") ||
		strings.Contains(contentType, "application/gzip") {
		return true
	}
	return false
}

