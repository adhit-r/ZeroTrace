package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	requests   int
	window     time.Duration
	buckets    map[string]*bucket
	mu         sync.RWMutex
	cleanupTicker *time.Ticker
}

type bucket struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: requests,
		window:   window,
		buckets:  make(map[string]*bucket),
	}

	// Start cleanup goroutine to remove old buckets
	rl.cleanupTicker = time.NewTicker(5 * time.Minute)
	go rl.cleanup()

	return rl
}

// Allow checks if a request is allowed for the given key
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	b, exists := rl.buckets[key]
	if !exists {
		b = &bucket{
			tokens:     rl.requests,
			lastRefill: time.Now(),
		}
		rl.buckets[key] = b
	}
	rl.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	// Refill tokens based on elapsed time
	now := time.Now()
	elapsed := now.Sub(b.lastRefill)
	tokensToAdd := int(elapsed / (rl.window / time.Duration(rl.requests)))
	
	if tokensToAdd > 0 {
		b.tokens = min(b.tokens+tokensToAdd, rl.requests)
		b.lastRefill = now
	}

	// Check if we have tokens available
	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// cleanup removes old buckets to prevent memory leaks
func (rl *RateLimiter) cleanup() {
	for range rl.cleanupTicker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, b := range rl.buckets {
			b.mu.Lock()
			// Remove buckets that haven't been used in 1 hour
			if now.Sub(b.lastRefill) > time.Hour {
				delete(rl.buckets, key)
			}
			b.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(cfg *config.Config) gin.HandlerFunc {
	limiter := NewRateLimiter(cfg.RateLimitRequests, cfg.RateLimitWindow)

	return func(c *gin.Context) {
		// Get client identifier (IP address or user ID)
		clientID := c.ClientIP()
		if userID, exists := c.Get("user_id"); exists {
			clientID = fmt.Sprintf("user:%s", userID)
		}

		// Check rate limit
		if !limiter.Allow(clientID) {
			// Use standardized error response format
			c.JSON(http.StatusTooManyRequests, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: fmt.Sprintf("Rate limit exceeded. Maximum %d requests per %v", cfg.RateLimitRequests, cfg.RateLimitWindow),
				},
				Timestamp: time.Now(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

