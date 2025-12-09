package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ClerkAuth middleware validates Clerk JWT tokens
func ClerkAuth() gin.HandlerFunc {
	// Get Clerk JWT verification key from environment
	clerkSecret := os.Getenv("CLERK_JWT_VERIFICATION_KEY")
	isDebug := os.Getenv("API_MODE") == "debug" || os.Getenv("DEBUG") == "true"

	if clerkSecret == "" {
		if isDebug {
			// Only allow development key in debug mode
			clerkSecret = "development-key"
		} else {
			// In production, fail if no key is provided
			// In production, fail if no key is provided
			// Use log.Fatal inside init/setup instead of panic
			// Note: We avoid importing "log" and use fmt.Println + os.Exit to match the style or add import.
			// Ideally we should import "log". Let's update imports too.
			fmt.Printf("FATAL: CLERK_JWT_VERIFICATION_KEY environment variable is required in production\n")
			os.Exit(1)
		}
	}

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

		// Demo mode only in debug/development mode
		if isDebug && token == "demo-valid-token" {
			c.Set("user_id", "demo-user-1")
			c.Set("company_id", "demo-company-1")
			c.Set("role", "admin")
			c.Next()
			return
		}

		// Validate Clerk JWT token
		claims, err := validateClerkToken(token, clerkSecret)
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

		// Set user context from Clerk claims
		if sub, exists := claims["sub"]; exists {
			c.Set("user_id", sub)
		}
		if orgID, exists := claims["org_id"]; exists {
			c.Set("company_id", orgID)
		}
		if orgRole, exists := claims["org_role"]; exists {
			c.Set("role", orgRole)
		}

		c.Next()
	}
}

// validateClerkToken validates a Clerk JWT token
func validateClerkToken(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
