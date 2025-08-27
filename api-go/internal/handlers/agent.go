package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAgents retrieves all agents for a company
func GetAgents(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		agents := agentService.GetAgents(companyUUID)

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      agents,
			Message:   "Agents retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// GetOnlineAgents retrieves online agents for a company
func GetOnlineAgents(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		agents := agentService.GetOnlineAgents(companyUUID)

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      agents,
			Message:   "Online agents retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// GetAgent retrieves a specific agent
func GetAgent(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_AGENT_ID",
					Message: "Invalid agent ID",
				},
				Timestamp: time.Now(),
			})
			return
		}

		agent, exists := agentService.GetAgent(agentID)
		if !exists {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "AGENT_NOT_FOUND",
					Message: "Agent not found",
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      agent,
			Message:   "Agent retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// GetAgentStats retrieves agent statistics for a company
func GetAgentStats(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		stats := agentService.GetAgentStats(companyUUID)

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      stats,
			Message:   "Agent statistics retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// AgentHeartbeat handles agent heartbeat updates
func AgentHeartbeat(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var heartbeat models.AgentHeartbeat
		if err := c.ShouldBindJSON(&heartbeat); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_HEARTBEAT",
					Message: "Invalid heartbeat data",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		// Set timestamp if not provided
		if heartbeat.Timestamp.IsZero() {
			heartbeat.Timestamp = time.Now()
		}

		err := agentService.UpdateAgentHeartbeat(heartbeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "HEARTBEAT_UPDATE_FAILED",
					Message: "Failed to update agent heartbeat",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "Heartbeat updated successfully",
			Timestamp: time.Now(),
		})
	}
}

// RegisterAgent handles agent registration
func RegisterAgent(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AgentID  uuid.UUID `json:"agent_id" binding:"required"`
			Name     string    `json:"name" binding:"required"`
			Version  string    `json:"version" binding:"required"`
			Hostname string    `json:"hostname" binding:"required"`
			OS       string    `json:"os" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_REGISTRATION",
					Message: "Invalid registration data",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		agent := agentService.RegisterAgent(
			req.AgentID,
			companyUUID,
			req.Name,
			req.Version,
			req.Hostname,
			req.OS,
		)

		c.JSON(http.StatusCreated, models.APIResponse{
			Success:   true,
			Data:      agent,
			Message:   "Agent registered successfully",
			Timestamp: time.Now(),
		})
	}
}

