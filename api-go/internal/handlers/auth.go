package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
)

// Login handles user login
func Login(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_REQUEST",
					Message: "Invalid request body",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		user, token, err := authService.Login(req.Email, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "AUTHENTICATION_FAILED",
					Message: "Invalid credentials",
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data: models.LoginResponse{
				User:  user,
				Token: token,
			},
			Message:   "Login successful",
			Timestamp: time.Now(),
		})
	}
}

// Register handles user registration
func Register(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_REQUEST",
					Message: "Invalid request body",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		user, err := authService.Register(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "REGISTRATION_FAILED",
					Message: "Registration failed",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusCreated, models.APIResponse{
			Success:   true,
			Data:      user,
			Message:   "User registered successfully",
			Timestamp: time.Now(),
		})
	}
}
