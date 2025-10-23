package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
)

// GetPublicVulnerabilities retrieves vulnerabilities for all agents (public endpoint)
func GetPublicVulnerabilities(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("DEBUG: GetPublicVulnerabilities called")
		// Get all agents to extract vulnerabilities
		agents := agentService.GetAllAgents()
		log.Printf("DEBUG: Found %d agents", len(agents))

		var vulnerabilities []map[string]interface{}

		// Extract vulnerabilities from agent metadata
		for _, agent := range agents {
			if agent.Metadata != nil {
				// Try to get vulnerabilities from the scan results
				if vulns, ok := agent.Metadata["vulnerabilities"]; ok && vulns != nil {
					// Debug logging
					fmt.Printf("DEBUG: Found vulnerabilities in agent %s, type: %T\n", agent.ID, vulns)

					// Handle both []models.Vulnerability and []interface{} types
					switch v := vulns.(type) {
					case []models.Vulnerability:
						fmt.Printf("DEBUG: Processing []models.Vulnerability, count: %d\n", len(v))
						// Convert models.Vulnerability to map[string]interface{}
						for _, vuln := range v {
							vulnMap := map[string]interface{}{
								"id":              vuln.ID,
								"cve_id":          vuln.CVEID,
								"title":           vuln.Title,
								"description":     vuln.Description,
								"severity":        vuln.Severity,
								"cvss_score":      vuln.CVSSScore,
								"package_name":    vuln.PackageName,
								"package_version": vuln.PackageVersion,
								"agent_id":        agent.ID,
								"agent_name":      agent.Name,
								"agent_hostname":  agent.Hostname,
							}
							vulnerabilities = append(vulnerabilities, vulnMap)
						}
					case []interface{}:
						fmt.Printf("DEBUG: Processing []interface{}, count: %d\n", len(v))
						// Handle interface{} format (most common case)
						for _, vuln := range v {
							if vulnMap, ok := vuln.(map[string]interface{}); ok {
								// Add agent information
								vulnMap["agent_id"] = agent.ID
								vulnMap["agent_name"] = agent.Name
								vulnMap["agent_hostname"] = agent.Hostname
								vulnerabilities = append(vulnerabilities, vulnMap)
							}
						}
					default:
						fmt.Printf("DEBUG: Unknown type: %T\n", v)
					}
				} else {
					fmt.Printf("DEBUG: No vulnerabilities found in agent %s\n", agent.ID)
				}
			}
		}

		// Add performance metrics to response
		fmt.Printf("DEBUG: Found %d vulnerabilities total\n", len(vulnerabilities))
		response := models.APIResponse{
			Success:   true,
			Data:      vulnerabilities,
			Message:   fmt.Sprintf("Vulnerabilities retrieved successfully - found %d vulnerabilities", len(vulnerabilities)),
			Timestamp: time.Now(),
		}

		c.JSON(http.StatusOK, response)
	}
}
