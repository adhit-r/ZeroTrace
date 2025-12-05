package handlers

import (
	"net/http"

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
			BadRequest(c, "INVALID_REQUEST", "Invalid request body", err.Error())
			return
		}

		// Get the user ID from the context (set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			Unauthorized(c, "UNAUTHORIZED", "User not authenticated")
			return
		}

		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			BadRequest(c, "INVALID_USER_ID", "Invalid user ID", err.Error())
			return
		}

		// Generate the enrollment token
		enrollmentToken, err := enrollmentService.GenerateEnrollmentToken(&req, userUUID)
		if err != nil {
			InternalServerError(c, "TOKEN_GENERATION_FAILED", "Failed to generate enrollment token", err)
			return
		}

		SuccessResponse(c, http.StatusOK, map[string]interface{}{
			"token":       enrollmentToken.Token,
			"expires_at":  enrollmentToken.ExpiresAt,
			"issued_at":   enrollmentToken.IssuedAt,
			"issued_by":   enrollmentToken.IssuedBy,
			"description": req.Description,
		}, "Enrollment token generated successfully")
	}
}

// EnrollAgent enrolls an agent using an enrollment token
func EnrollAgent(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.AgentEnrollmentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			BadRequest(c, "INVALID_REQUEST", "Invalid request body", err.Error())
			return
		}

		// Enroll the agent
		response, err := enrollmentService.EnrollAgent(&req)
		if err != nil {
			BadRequest(c, "ENROLLMENT_FAILED", "Failed to enroll agent", err.Error())
			return
		}

		SuccessResponse(c, http.StatusOK, response, "Agent enrolled successfully")
	}
}

// RevokeEnrollmentToken revokes an enrollment token
func RevokeEnrollmentToken(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			BadRequest(c, "INVALID_TOKEN_ID", "Invalid token ID", err.Error())
			return
		}

		// Revoke the token
		err = enrollmentService.RevokeEnrollmentToken(tokenID)
		if err != nil {
			InternalServerError(c, "REVOKE_FAILED", "Failed to revoke enrollment token", err)
			return
		}

		SuccessResponse(c, http.StatusOK, nil, "Enrollment token revoked successfully")
	}
}

// RevokeAgentCredential revokes an agent's credential
func RevokeAgentCredential(enrollmentService *services.EnrollmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		credentialID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			BadRequest(c, "INVALID_CREDENTIAL_ID", "Invalid credential ID", err.Error())
			return
		}

		// Revoke the credential
		err = enrollmentService.RevokeAgentCredential(credentialID)
		if err != nil {
			InternalServerError(c, "REVOKE_FAILED", "Failed to revoke agent credential", err)
			return
		}

		SuccessResponse(c, http.StatusOK, nil, "Agent credential revoked successfully")
	}
}


