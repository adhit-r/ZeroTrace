package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
)

// GetPublicVulnerabilities retrieves vulnerabilities for all agents (public endpoint)
func GetPublicVulnerabilities(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get all agents to extract vulnerabilities
		agents := agentService.GetAllAgents()

		var vulnerabilities []map[string]interface{}

		// Extract vulnerabilities from agent metadata
		for _, agent := range agents {
			if agent.Metadata != nil {
				if vulns, ok := agent.Metadata["vulnerabilities"].([]interface{}); ok {
					for _, vuln := range vulns {
						if vulnMap, ok := vuln.(map[string]interface{}); ok {
							vulnMap["agent_id"] = agent.ID
							vulnMap["agent_name"] = agent.Name
							vulnMap["agent_hostname"] = agent.Hostname
							vulnerabilities = append(vulnerabilities, vulnMap)
						}
					}
				}
			}
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      vulnerabilities,
			Message:   "Vulnerabilities retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}
