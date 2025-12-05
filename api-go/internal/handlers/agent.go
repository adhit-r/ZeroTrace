package handlers

import (
	"bytes"
	"io"
	"log"
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

		SuccessResponse(c, http.StatusOK, agents, "Agents retrieved successfully")
	}
}

// GetOnlineAgents retrieves online agents for a company
func GetOnlineAgents(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		agents := agentService.GetOnlineAgents(companyUUID)

		SuccessResponse(c, http.StatusOK, agents, "Online agents retrieved successfully")
	}
}

// GetAgent retrieves a specific agent
func GetAgent(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		agentID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			BadRequest(c, "INVALID_AGENT_ID", "Invalid agent ID", err.Error())
			return
		}

		agent, exists := agentService.GetAgent(agentID)
		if !exists {
			NotFound(c, "AGENT_NOT_FOUND", "Agent not found")
			return
		}

		SuccessResponse(c, http.StatusOK, agent, "Agent retrieved successfully")
	}
}

// GetAgentStats retrieves agent statistics for a company
func GetAgentStats(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		stats := agentService.GetAgentStats(companyUUID)

		SuccessResponse(c, http.StatusOK, stats, "Agent statistics retrieved successfully")
	}
}

// AgentHeartbeat handles agent heartbeat updates
func AgentHeartbeat(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Temporary struct to bind the request payload with string IDs
		var req struct {
			AgentID        string                 `json:"agent_id"`
			OrganizationID string                 `json:"organization_id"`
			AgentName      string                 `json:"agent_name"`
			Status         string                 `json:"status"`
			CPUUsage       float64                `json:"cpu_usage"`
			MemoryUsage    float64                `json:"memory_usage"`
			Metadata       map[string]interface{} `json:"metadata"`
			Timestamp      time.Time              `json:"timestamp"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			BadRequest(c, "INVALID_HEARTBEAT", "Invalid heartbeat data", err.Error())
			return
		}

		// Parse string IDs into UUIDs
		agentUUID, err := uuid.Parse(req.AgentID)
		if err != nil {
			BadRequest(c, "INVALID_AGENT_ID", "Invalid agent ID format", err.Error())
			return
		}

		// Handle missing or invalid organization ID
		var orgUUID uuid.UUID
		if req.OrganizationID == "" {
			// Use default organization ID if not provided
			orgUUID, _ = uuid.Parse("00000000-0000-0000-0000-000000000001")
		} else {
			orgUUID, err = uuid.Parse(req.OrganizationID)
			if err != nil {
				BadRequest(c, "INVALID_ORGANIZATION_ID", "Invalid organization ID format", err.Error())
				return
			}
		}

		// Create the heartbeat model
		heartbeat := models.AgentHeartbeat{
			AgentID:        agentUUID,
			OrganizationID: orgUUID,
			AgentName:      req.AgentName,
			Status:         req.Status,
			CPUUsage:       req.CPUUsage,
			MemoryUsage:    req.MemoryUsage,
			Metadata:       req.Metadata,
			Timestamp:      req.Timestamp,
		}

		// Set timestamp if not provided
		if heartbeat.Timestamp.IsZero() {
			heartbeat.Timestamp = time.Now()
		}

		// Debug: log heartbeat metadata
		log.Printf("[Heartbeat Handler] Heartbeat metadata: %v", heartbeat.Metadata)

		err = agentService.UpdateAgentHeartbeat(heartbeat)
		if err != nil {
			InternalServerError(c, "HEARTBEAT_UPDATE_FAILED", "Failed to update agent heartbeat", err)
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
		// Temporary struct to bind the request payload with string IDs
		var req struct {
			ID             string `json:"id"`
			CompanyID      string `json:"company_id"`
			OrganizationID string `json:"organization_id"`
			Name           string `json:"name"`
			Status         string `json:"status"`
			Version        string `json:"version"`
			Hostname       string `json:"hostname"`
			OS             string `json:"os"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		// Parse string IDs into UUIDs
		agentUUID, err := uuid.Parse(req.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID format"})
			return
		}

		// Handle missing or invalid organization ID
		var orgUUID uuid.UUID
		if req.OrganizationID == "" {
			// Use default organization ID if not provided
			orgUUID, _ = uuid.Parse("00000000-0000-0000-0000-000000000001")
		} else {
			orgUUID, err = uuid.Parse(req.OrganizationID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID format"})
				return
			}
		}

		// Create the agent model
		agent := models.Agent{
			ID:             agentUUID,
			OrganizationID: orgUUID,
			Name:           req.Name,
			Status:         "active", // Set initial status
			Version:        req.Version,
			Hostname:       req.Hostname,
			OS:             req.OS,
		}

		// Use a default company ID for now
		defaultCompanyUUID, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
		agent.CompanyID = defaultCompanyUUID

		registeredAgent, err := agentService.RegisterAgent(agent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register agent"})
			return
		}
		c.JSON(http.StatusOK, registeredAgent)
	}
}

// AgentResults handles scan results from agents
func AgentResults(agentService *services.AgentService, enrichmentService *services.EnrichmentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[AgentResults] *** REQUEST RECEIVED *** from %s", c.ClientIP())

		// Log raw request body for debugging
		bodyBytes, _ := c.GetRawData()
		log.Printf("[AgentResults] Received request from agent, body length: %d bytes", len(bodyBytes))

		// Restore body for binding
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var req struct {
			AgentID  string                   `json:"agent_id" binding:"required"`
			Results  []models.AgentScanResult `json:"results"`
			Metadata map[string]interface{}   `json:"metadata"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("[AgentResults] JSON binding error: %v", err)
			log.Printf("[AgentResults] Raw body: %s", string(bodyBytes))
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success:   false,
				Message:   "Invalid request body: " + err.Error(),
				Timestamp: time.Now(),
			})
			return
		}

		log.Printf("[AgentResults] Successfully parsed request for agent %s with %d results", req.AgentID, len(req.Results))

		// Extract dependencies from scan results for enrichment
		var allDependencies []models.Dependency
		for _, result := range req.Results {
			allDependencies = append(allDependencies, result.Dependencies...)
		}

		// Enrich dependencies with CVE data
		if len(allDependencies) > 0 {
			log.Printf("[AgentResults] Enriching %d dependencies with CVE data", len(allDependencies))
			enrichedVulns, err := enrichmentService.EnrichDependencies(allDependencies)
			if err != nil {
				log.Printf("[AgentResults] Enrichment failed: %v", err)
				// Continue without enrichment rather than failing
			} else {
				log.Printf("[AgentResults] Found %d vulnerabilities from enrichment", len(enrichedVulns))
				// Store enriched vulnerabilities in metadata
				if req.Metadata == nil {
					req.Metadata = make(map[string]interface{})
				}
				req.Metadata["enriched_vulnerabilities"] = enrichedVulns
				req.Metadata["enrichment_timestamp"] = time.Now()
			}
		}

		// Update agent with results (including enriched vulnerabilities)
		err := agentService.UpdateAgentResults(req.AgentID, req.Results, req.Metadata)
		if err != nil {
			log.Printf("[AgentResults] Failed to update agent results: %v", err)
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success:   false,
				Message:   "Failed to update agent results: " + err.Error(),
				Timestamp: time.Now(),
			})
			return
		}

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
		totalApplications := 0
		onlineAgents := 0
		vulnerableAssets := 0
		criticalVulns := 0
		highVulns := 0
		mediumVulns := 0
		lowVulns := 0
		totalVulns := 0
		lastScan := time.Time{}

		for _, agent := range agents {
			// Count online agents (seen within last 5 minutes)
			if time.Since(agent.LastSeen) < 5*time.Minute {
				onlineAgents++
			}

			// Count actual assets (devices/systems) - each agent represents one asset
			if agent.Hostname != "" {
				totalAssets++
			}

			// Count applications from agent metadata
			if agent.Metadata != nil {
				// Count applications processed by this agent
				if applicationsProcessed, ok := agent.Metadata["applications_processed"]; ok && applicationsProcessed != nil {
					if count, ok := applicationsProcessed.(float64); ok {
						totalApplications += int(count)
					}
				}

				// Count vulnerabilities by severity
				agentCritical := 0
				agentHigh := 0
				agentMedium := 0
				agentLow := 0
				agentTotal := 0

				// Handle null values properly - if field doesn't exist or is null, treat as 0
				if critical, ok := agent.Metadata["critical_vulnerabilities"]; ok && critical != nil {
					if count, ok := critical.(float64); ok {
						agentCritical = int(count)
						criticalVulns += agentCritical
					}
				}
				if high, ok := agent.Metadata["high_vulnerabilities"]; ok && high != nil {
					if count, ok := high.(float64); ok {
						agentHigh = int(count)
						highVulns += agentHigh
					}
				}
				if medium, ok := agent.Metadata["medium_vulnerabilities"]; ok && medium != nil {
					if count, ok := medium.(float64); ok {
						agentMedium = int(count)
						mediumVulns += agentMedium
					}
				}
				if low, ok := agent.Metadata["low_vulnerabilities"]; ok && low != nil {
					if count, ok := low.(float64); ok {
						agentLow = int(count)
						lowVulns += agentLow
					}
				}
				if total, ok := agent.Metadata["total_vulnerabilities"]; ok && total != nil {
					if count, ok := total.(float64); ok {
						agentTotal = int(count)
						totalVulns += agentTotal
					}
				}

				// Count vulnerable assets (agents with vulnerabilities)
				if agentTotal > 0 {
					vulnerableAssets++
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
				"critical":   criticalVulns,
				"high":       highVulns,
				"medium":     mediumVulns,
				"low":        lowVulns,
				"lastScan":   lastScan.Format(time.RFC3339),
			},
			"applications": map[string]interface{}{
				"total": totalApplications,
			},
			"vulnerabilities": map[string]interface{}{
				"total":    totalVulns,
				"critical": criticalVulns,
				"high":     highVulns,
				"medium":   mediumVulns,
				"low":      lowVulns,
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

// SystemInfoRequest represents the system info request payload
type SystemInfoRequest struct {
	AgentID    string                 `json:"agent_id" binding:"required"`
	SystemInfo map[string]interface{} `json:"system_info"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// UpdateSystemInfo updates agent system information
func UpdateSystemInfo(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SystemInfoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success:   false,
				Message:   "Invalid request payload",
				Timestamp: time.Now(),
			})
			return
		}

		// Update agent with system information
		err := agentService.UpdateAgentSystemInfo(req.AgentID, req.SystemInfo)
		if err != nil {
			log.Printf("Error updating agent system info: %v", err)
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success:   false,
				Message:   "Failed to update system information",
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "System information updated successfully",
			Timestamp: time.Now(),
		})
	}
}

// NetworkScanResults handles network scan results from agents
func NetworkScanResults(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[NetworkScanResults] Request received from %s", c.ClientIP())

		var req struct {
			AgentID    string                 `json:"agent_id" binding:"required"`
			ScanResult map[string]interface{} `json:"scan_result" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("[NetworkScanResults] JSON binding error: %v", err)
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success:   false,
				Message:   "Invalid request body: " + err.Error(),
				Timestamp: time.Now(),
			})
			return
		}

		log.Printf("[NetworkScanResults] Successfully parsed network scan results for agent %s", req.AgentID)

		// Store network scan results in agent metadata
		// In a full implementation, this would be stored in a separate table
		err := agentService.UpdateAgentMetadata(req.AgentID, map[string]interface{}{
			"last_network_scan": time.Now(),
			"network_scan_result": req.ScanResult,
		})
		if err != nil {
			log.Printf("[NetworkScanResults] Failed to update agent metadata: %v", err)
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success:   false,
				Message:   "Failed to store network scan results: " + err.Error(),
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "Network scan results received successfully",
			Timestamp: time.Now(),
		})
	}
}
