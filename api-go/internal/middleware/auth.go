package middleware

import (
	"net/http"
	"strings"

	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
)

// Auth middleware validates JWT tokens
func Auth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "NO_TOKEN",
					"message": "No authorization token provided",
				},
			})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix
		if strings.HasPrefix(token, "Bearer ") {
			token = token[7:]
		}

		// Demo mode - accept demo tokens
		if token == "demo-valid-token" {
			c.Set("user_id", "demo-user-1")
			c.Set("company_id", "demo-company-1")
			c.Set("role", "admin")
			c.Next()
			return
		}

		// Validate token
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Set user context
		if userID, exists := claims["user_id"]; exists {
			c.Set("user_id", userID)
		}
		if companyID, exists := claims["company_id"]; exists {
			c.Set("company_id", companyID)
		}
		if role, exists := claims["role"]; exists {
			c.Set("role", role)
		}

		c.Next()
	}
}
