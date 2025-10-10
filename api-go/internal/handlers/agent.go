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
		// For public endpoint, get all agents without company filter
		agents := agentService.GetAllAgents()

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

		// For public endpoint, use a default company ID
		companyUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

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

// AgentResults handles scan results from agents
func AgentResults(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AgentID  string                 `json:"agent_id" binding:"required"`
			Results  []models.AgentScanResult `json:"results"`
			Metadata map[string]interface{} `json:"metadata"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success:   false,
				Message:   "Invalid request body",
				Timestamp: time.Now(),
			})
			return
		}

		// Update agent with results directly (AgentScanResult format)
		agentService.UpdateAgentResults(req.AgentID, req.Results, req.Metadata)

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "Scan results received successfully",
			Timestamp: time.Now(),
		})
	}
}

// AgentStatus handles agent status updates
func AgentStatus(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			AgentID  string                 `json:"agent_id" binding:"required"`
			Status   string                 `json:"status" binding:"required"`
			Metadata map[string]interface{} `json:"metadata"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success:   false,
				Message:   "Invalid request body",
				Timestamp: time.Now(),
			})
			return
		}

		// Update agent status
		agentService.UpdateAgentStatus(req.AgentID, req.Status, req.Metadata)

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "Agent status updated successfully",
			Timestamp: time.Now(),
		})
	}
}

// GetPublicDashboardOverview provides public dashboard data
func GetPublicDashboardOverview(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get all agents
		agents := agentService.GetAllAgents()

		// Calculate dashboard metrics
		totalAssets := 0
		onlineAgents := 0
		vulnerableAssets := 0
		criticalVulnerabilities := 0
		lastScan := time.Time{}

		for _, agent := range agents {
			// Count online agents (seen within last 5 minutes)
			if time.Since(agent.LastSeen) < 5*time.Minute {
				onlineAgents++
			}

			// Count actual scanned assets from agent metadata
			if agent.Metadata != nil {
				// Count total assets scanned by this agent
				if totalAssetsFromAgent, ok := agent.Metadata["total_assets"]; ok {
					if count, ok := totalAssetsFromAgent.(float64); ok {
						totalAssets += int(count)
					}
				}

				// Count vulnerable assets
				if vulns, ok := agent.Metadata["vulnerabilities_found"]; ok {
					if count, ok := vulns.(float64); ok && count > 0 {
						vulnerableAssets++
					}
				}

				// Count critical vulnerabilities
				if critical, ok := agent.Metadata["critical_vulnerabilities"]; ok {
					if count, ok := critical.(float64); ok {
						criticalVulnerabilities += int(count)
					}
				}
			}

			// Track last scan time
			if agent.LastSeen.After(lastScan) {
				lastScan = agent.LastSeen
			}
		}

		// Create dashboard response
		dashboardData := map[string]interface{}{
			"assets": map[string]interface{}{
				"total":      totalAssets,
				"vulnerable": vulnerableAssets,
				"critical":   criticalVulnerabilities,
				"high":       0, // Placeholder
				"medium":     0, // Placeholder
				"low":        0, // Placeholder
				"lastScan":   lastScan.Format(time.RFC3339),
			},
			"vulnerabilities": map[string]interface{}{
				"total":    criticalVulnerabilities,
				"critical": criticalVulnerabilities,
				"high":     0, // Placeholder
				"medium":   0, // Placeholder
				"low":      0, // Placeholder
			},
			"agents": map[string]interface{}{
				"total":  len(agents),
				"online": onlineAgents,
			},
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      dashboardData,
			Message:   "Dashboard overview retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// GetPublicAgentStats retrieves agent statistics for all agents (public endpoint)
func GetPublicAgentStats(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats := agentService.GetPublicAgentStats()

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      stats,
			Message:   "Agent statistics retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}
