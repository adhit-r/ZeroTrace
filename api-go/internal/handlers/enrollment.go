package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GenerateEnrollmentToken generates a new enrollment token for an organization
func GenerateEnrollmentToken(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.GenerateEnrollmentTokenRequest
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

		// Get the user ID from the context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "UNAUTHORIZED",
					Message: "User not authenticated",
				},
				Timestamp: time.Now(),
			})
			return
		}

		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_USER_ID",
					Message: "Invalid user ID",
				},
				Timestamp: time.Now(),
			})
			return
		}

		// Generate the enrollment token
		enrollmentToken, err := enrollmentService.GenerateEnrollmentToken(&req, userUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "TOKEN_GENERATION_FAILED",
					Message: "Failed to generate enrollment token",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success: true,
			Data: map[string]interface{}{
				"token":       enrollmentToken.Token,
				"expires_at":  enrollmentToken.ExpiresAt,
				"issued_at":   enrollmentToken.IssuedAt,
				"issued_by":   enrollmentToken.IssuedBy,
				"description": req.Description,
			},
			Message:   "Enrollment token generated successfully",
			Timestamp: time.Now(),
		})
	}
}

// EnrollAgent enrolls an agent using an enrollment token
func EnrollAgent(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AgentEnrollmentRequest
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

		// Enroll the agent
		response, err := enrollmentService.EnrollAgent(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "ENROLLMENT_FAILED",
					Message: "Failed to enroll agent",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      response,
			Message:   "Agent enrolled successfully",
			Timestamp: time.Now(),
		})
	}
}

// RevokeEnrollmentToken revokes an enrollment token
func RevokeEnrollmentToken(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_TOKEN_ID",
					Message: "Invalid token ID",
				},
				Timestamp: time.Now(),
			})
			return
		}

		// Revoke the token
		err = enrollmentService.RevokeEnrollmentToken(tokenID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "REVOKE_FAILED",
					Message: "Failed to revoke enrollment token",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "Enrollment token revoked successfully",
			Timestamp: time.Now(),
		})
	}
}

// RevokeAgentCredential revokes an agent's credential
func RevokeAgentCredential(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		credentialID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_CREDENTIAL_ID",
					Message: "Invalid credential ID",
				},
				Timestamp: time.Now(),
			})
			return
		}

		// Revoke the credential
		err = enrollmentService.RevokeAgentCredential(credentialID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "REVOKE_FAILED",
					Message: "Failed to revoke agent credential",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "Agent credential revoked successfully",
			Timestamp: time.Now(),
		})
	}
}


