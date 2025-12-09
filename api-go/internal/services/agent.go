package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentService manages agent registration and heartbeats
type AgentService struct {
	agents map[uuid.UUID]*models.Agent
	mutex  sync.RWMutex
	db     *gorm.DB
}

// NewAgentService creates a new agent service
func NewAgentService(db *gorm.DB) *AgentService {
	// Restore agents from DB on startup
	agents := make(map[uuid.UUID]*models.Agent)
	var loadedAgents []models.Agent
	if err := db.Find(&loadedAgents).Error; err == nil {
		for _, agent := range loadedAgents {
			// Create a copy of the loop variable
			a := agent
			agents[agent.ID] = &a
		}
		log.Printf("[NewAgentService] Restored %d agents from database", len(agents))
	} else {
		log.Printf("[NewAgentService] Failed to load agents from DB: %v", err)
	}

	return &AgentService{
		agents: agents,
		db:     db,
	}
}

// RegisterAgent registers a new agent or updates an existing one
func (as *AgentService) RegisterAgent(agent models.Agent) (*models.Agent, error) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	// If agent ID is not provided, generate a new one
	if agent.ID == uuid.Nil {
		agent.ID = uuid.New()
	}

	agent.LastSeen = time.Now()
	as.agents[agent.ID] = &agent

	// Persist to DB
	if err := as.db.Save(&agent).Error; err != nil {
		log.Printf("Failed to persist registered agent %s: %v", agent.ID, err)
	}

	log.Printf("Agent registered or updated: %s", agent.ID)
	return &agent, nil
}

// UpdateAgentHeartbeat updates agent heartbeat
func (as *AgentService) UpdateAgentHeartbeat(heartbeat models.AgentHeartbeat) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	agent, exists := as.agents[heartbeat.AgentID]
	if !exists {
		// Create new agent if it doesn't exist
		agent = &models.Agent{
			ID:             heartbeat.AgentID,
			OrganizationID: heartbeat.OrganizationID,
			Name:           heartbeat.AgentName,
			Status:         "active",
			LastSeen:       time.Now(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		as.agents[heartbeat.AgentID] = agent
	}

	// Update agent status
	agent.LastSeen = time.Now()
	agent.CPUUsage = heartbeat.CPUUsage
	agent.MemoryUsage = heartbeat.MemoryUsage
	agent.Status = heartbeat.Status

	// Log metadata before merge
	log.Printf("[UpdateAgentHeartbeat] Metadata BEFORE merge: %v", getMetadataKeys(agent.Metadata))
	log.Printf("[UpdateAgentHeartbeat] Heartbeat metadata: %v", heartbeat.Metadata)

	// Merge heartbeat metadata instead of overwriting
	if agent.Metadata == nil {
		agent.Metadata = make(map[string]interface{})
	}
	for k, v := range heartbeat.Metadata {
		agent.Metadata[k] = v
	}

	// Log metadata after merge
	log.Printf("[UpdateAgentHeartbeat] Metadata AFTER merge: %v", getMetadataKeys(agent.Metadata))

	agent.UpdatedAt = time.Now()

	// Persist to DB
	if err := as.db.Save(agent).Error; err != nil {
		log.Printf("Failed to persist agent heartbeat %s: %v", agent.ID, err)
	}

	return nil
}

// GetAgent gets an agent by ID
func (as *AgentService) GetAgent(agentID uuid.UUID) (*models.Agent, bool) {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	agent, exists := as.agents[agentID]
	return agent, exists
}

// GetAgents gets all agents for an organization
func (as *AgentService) GetAgents(organizationID uuid.UUID) []*models.Agent {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	var agents []*models.Agent
	for _, agent := range as.agents {
		if agent.OrganizationID == organizationID {
			agents = append(agents, agent)
		}
	}
	return agents
}

// GetAllAgents gets all agents (for public endpoints)
func (as *AgentService) GetAllAgents() []*models.Agent {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	var agents []*models.Agent
	for _, agent := range as.agents {
		log.Printf("[GetAllAgents] Agent %s metadata keys: %v", agent.ID, getMetadataKeys(agent.Metadata))
		if deps, ok := agent.Metadata["dependencies"]; ok {
			log.Printf("[GetAllAgents] Dependencies type: %T, value: %v", deps, deps)
		}
		agents = append(agents, agent)
	}
	return agents
}

// GetOnlineAgents gets online agents for an organization
func (as *AgentService) GetOnlineAgents(organizationID uuid.UUID) []*models.Agent {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	var agents []*models.Agent
	offlineThreshold := time.Now().Add(-5 * time.Minute) // Consider offline after 5 minutes

	for _, agent := range as.agents {
		if agent.OrganizationID == organizationID && agent.LastSeen.After(offlineThreshold) {
			agents = append(agents, agent)
		}
	}
	return agents
}

// RemoveAgent removes an agent
func (as *AgentService) RemoveAgent(agentID uuid.UUID) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	delete(as.agents, agentID)
}

// CleanupOfflineAgents removes agents that haven't been seen for a while
func (as *AgentService) CleanupOfflineAgents() {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	offlineThreshold := time.Now().Add(-30 * time.Minute) // Remove after 30 minutes offline

	for agentID, agent := range as.agents {
		if agent.LastSeen.Before(offlineThreshold) {
			delete(as.agents, agentID)
		}
	}
}

// GetAgentStats gets agent statistics for an organization
func (as *AgentService) GetAgentStats(organizationID uuid.UUID) map[string]interface{} {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_agents":   0,
		"online_agents":  0,
		"offline_agents": 0,
		"total_cpu":      0.0,
		"total_memory":   0.0,
	}

	offlineThreshold := time.Now().Add(-5 * time.Minute)

	for _, agent := range as.agents {
		if agent.OrganizationID == organizationID {
			stats["total_agents"] = stats["total_agents"].(int) + 1
			stats["total_cpu"] = stats["total_cpu"].(float64) + agent.CPUUsage
			stats["total_memory"] = stats["total_memory"].(float64) + agent.MemoryUsage

			if agent.LastSeen.After(offlineThreshold) {
				stats["online_agents"] = stats["online_agents"].(int) + 1
			} else {
				stats["offline_agents"] = stats["offline_agents"].(int) + 1
			}
		}
	}

	// Calculate averages
	if stats["total_agents"].(int) > 0 {
		stats["avg_cpu"] = stats["total_cpu"].(float64) / float64(stats["total_agents"].(int))
		stats["avg_memory"] = stats["total_memory"].(float64) / float64(stats["total_agents"].(int))
	} else {
		stats["avg_cpu"] = 0.0
		stats["avg_memory"] = 0.0
	}

	return stats
}

// GetPublicAgentStats gets agent statistics for all agents (public endpoint)
func (as *AgentService) GetPublicAgentStats() map[string]interface{} {
	as.mutex.RLock()
	defer as.mutex.RUnlock()

	stats := map[string]interface{}{
		"total":     0,
		"online":    0,
		"offline":   0,
		"avgCpu":    0.0,
		"avgMemory": 0.0,
	}

	offlineThreshold := time.Now().Add(-5 * time.Minute)
	totalCpu := 0.0
	totalMemory := 0.0

	for _, agent := range as.agents {
		stats["total"] = stats["total"].(int) + 1
		totalCpu += agent.CPUUsage
		totalMemory += agent.MemoryUsage

		if agent.LastSeen.After(offlineThreshold) {
			stats["online"] = stats["online"].(int) + 1
		} else {
			stats["offline"] = stats["offline"].(int) + 1
		}
	}

	// Calculate averages
	if stats["total"].(int) > 0 {
		stats["avgCpu"] = totalCpu / float64(stats["total"].(int))
		stats["avgMemory"] = totalMemory / float64(stats["total"].(int))
	}

	return stats
}

// StartCleanupRoutine starts the cleanup routine
func (as *AgentService) StartCleanupRoutine() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			as.CleanupOfflineAgents()
		}
	}()
}

// UpdateAgentResults updates agent with scan results
func (as *AgentService) UpdateAgentResults(agentID string, results []models.AgentScanResult, metadata map[string]interface{}) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	log.Printf("[UpdateAgentResults] Received agent ID: '%s' (length: %d)", agentID, len(agentID))

	agentUUID, err := uuid.Parse(agentID)
	if err != nil {
		log.Printf("[UpdateAgentResults] Invalid agent ID format: %s, error: %v", agentID, err)
		return fmt.Errorf("invalid agent ID format: %w", err)
	}

	agent, exists := as.agents[agentUUID]
	if !exists {
		log.Printf("[UpdateAgentResults] Agent not found: %s", agentID)
		return fmt.Errorf("agent not found: %s", agentID)
	}

	// Update agent with scan results
	agent.LastSeen = time.Now()
	agent.UpdatedAt = time.Now()

	log.Printf("[UpdateAgentResults] Updating agent %s with %d scan results", agentID, len(results))
	log.Printf("[UpdateAgentResults] Results length: %d", len(results))

	// Initialize metadata if nil
	if agent.Metadata == nil {
		agent.Metadata = make(map[string]interface{})
	}

	// Store scan results in metadata
	log.Printf("[UpdateAgentResults] Checking if len(results) > 0: %d > 0 = %t", len(results), len(results) > 0)
	if len(results) > 0 {
		// Get existing dependencies and vulnerabilities to preserve them
		var existingDependencies []models.Dependency
		var existingVulnerabilities []models.Vulnerability

		if deps, ok := agent.Metadata["dependencies"].([]models.Dependency); ok {
			existingDependencies = deps
		}
		if vulns, ok := agent.Metadata["vulnerabilities"].([]models.Vulnerability); ok {
			existingVulnerabilities = vulns
		}

		// Count total vulnerabilities
		totalVulns := 0
		criticalVulns := 0
		highVulns := 0
		mediumVulns := 0
		lowVulns := 0
		totalAssets := 0

		// Collect all dependencies and vulnerabilities from new results
		var newDependencies []models.Dependency
		var newVulnerabilities []models.Vulnerability

		for _, result := range results {
			// Count dependencies as assets (agent sends Dependencies)
			totalAssets += len(result.Dependencies)
			newDependencies = append(newDependencies, result.Dependencies...)

			for _, vuln := range result.Vulnerabilities {
				totalVulns++
				newVulnerabilities = append(newVulnerabilities, vuln)
				switch vuln.Severity {
				case "critical":
					criticalVulns++
				case "high":
					highVulns++
				case "medium":
					mediumVulns++
				case "low":
					lowVulns++
				}
			}
		}

		// Merge with existing data
		allDependencies := append(existingDependencies, newDependencies...)
		allVulnerabilities := append(existingVulnerabilities, newVulnerabilities...)

		// Store actual data arrays
		log.Printf("[UpdateAgentResults] Storing %d dependencies and %d vulnerabilities in metadata (existing: %d deps, %d vulns)", len(allDependencies), len(allVulnerabilities), len(existingDependencies), len(existingVulnerabilities))
		agent.Metadata["dependencies"] = allDependencies
		agent.Metadata["vulnerabilities"] = allVulnerabilities
		log.Printf("[UpdateAgentResults] Dependencies stored successfully")

		// PERSISTENCE: Save Software/Dependencies to Database
		if len(allDependencies) > 0 {
			go func(agentID uuid.UUID, deps []models.Dependency) {
				for _, dep := range deps {
					software := models.Software{
						AgentID:   agentID,
						Name:      dep.Name,
						Version:   dep.Version,
						Type:      dep.Type,
						Status:    "active",
						Vendor:    dep.Description, // Using description as vendor for now
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
					// Upsert software based on AgentID + Name + Version
					if err := as.db.Where("agent_id = ? AND name = ? AND version = ?", agentID, dep.Name, dep.Version).
						FirstOrCreate(&software).Error; err != nil {
						log.Printf("Failed to persist software %s: %v", dep.Name, err)
					}
				}
				log.Printf("Persisted %d software items for agent %s", len(deps), agentID)
			}(agentUUID, allDependencies)
		}

		// Store counts in metadata
		agent.Metadata["total_vulnerabilities"] = totalVulns
		agent.Metadata["critical_vulnerabilities"] = criticalVulns
		agent.Metadata["high_vulnerabilities"] = highVulns
		agent.Metadata["medium_vulnerabilities"] = mediumVulns
		agent.Metadata["low_vulnerabilities"] = lowVulns
		agent.Metadata["total_assets"] = totalAssets
		agent.Metadata["last_scan_time"] = time.Now().Format(time.RFC3339)

		// Start async enrichment if we have dependencies
		if len(allDependencies) > 0 {
			log.Printf("[UpdateAgentResults] Starting async enrichment for agent %s with %d applications", agentID, len(allDependencies))
			// TODO: Call enrichment service asynchronously
			// This will be implemented when we integrate the enrichment service
		}
	}

	// Update with provided metadata (but preserve dependencies and vulnerabilities)
	if metadata != nil {
		for k, v := range metadata {
			// Don't overwrite dependencies and vulnerabilities that we just stored
			if k != "dependencies" && k != "vulnerabilities" {
				agent.Metadata[k] = v
			}
		}
	}

	log.Printf("[UpdateAgentResults] Final metadata keys: %v", getMetadataKeys(agent.Metadata))

	// Persist to DB
	if err := as.db.Save(agent).Error; err != nil {
		log.Printf("Failed to persist agent results %s: %v", agent.ID, err)
	}

	return nil
}

// Helper function to get metadata keys for debugging
func getMetadataKeys(metadata map[string]interface{}) []string {
	keys := make([]string, 0, len(metadata))
	for k := range metadata {
		keys = append(keys, k)
	}
	return keys
}

// UpdateAgentStatus updates agent status
func (as *AgentService) UpdateAgentStatus(agentID string, status string, metadata map[string]interface{}) {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	agentUUID, err := uuid.Parse(agentID)
	if err != nil {
		return
	}

	agent, exists := as.agents[agentUUID]
	if !exists {
		return
	}

	// Update agent status
	agent.Status = status
	agent.LastSeen = time.Now()
	if metadata != nil {
		agent.Metadata = metadata
	}

	// Persist to DB
	if err := as.db.Save(agent).Error; err != nil {
		log.Printf("Failed to persist agent status %s: %v", agent.ID, err)
	}
}

// UpdateAgentSystemInfo updates agent with system information
func (as *AgentService) UpdateAgentSystemInfo(agentID string, systemInfo map[string]interface{}) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	log.Printf("[UpdateAgentSystemInfo] Received system info for agent %s: %+v", agentID, systemInfo)

	agentUUID, err := uuid.Parse(agentID)
	if err != nil {
		return fmt.Errorf("invalid agent ID format: %w", err)
	}

	agent, exists := as.agents[agentUUID]
	if !exists {
		log.Printf("[UpdateAgentSystemInfo] Agent not found: %s", agentID)
		return fmt.Errorf("agent not found: %s", agentID)
	}

	log.Printf("[UpdateAgentSystemInfo] Found agent, current hostname: %s, new hostname: %v", agent.Hostname, systemInfo["hostname"])

	// Update system information fields
	if hostname, ok := systemInfo["hostname"].(string); ok {
		log.Printf("[UpdateAgentSystemInfo] Setting hostname to: %s", hostname)
		agent.Hostname = hostname
	}
	if osName, ok := systemInfo["os_name"].(string); ok {
		log.Printf("[UpdateAgentSystemInfo] Setting OS to: %s", osName)
		agent.OS = osName
		agent.OSName = osName
	}
	if osVersion, ok := systemInfo["os_version"].(string); ok {
		agent.OSVersion = osVersion
	}
	if osBuild, ok := systemInfo["os_build"].(string); ok {
		agent.OSBuild = osBuild
	}
	if kernelVersion, ok := systemInfo["kernel_version"].(string); ok {
		agent.KernelVersion = kernelVersion
	}
	if cpuModel, ok := systemInfo["cpu_model"].(string); ok {
		agent.CPUModel = cpuModel
	}
	if cpuCores, ok := systemInfo["cpu_cores"].(float64); ok {
		agent.CPUCores = int(cpuCores)
	}
	if memoryTotalGB, ok := systemInfo["memory_total_gb"].(float64); ok {
		agent.MemoryTotalGB = memoryTotalGB
	}
	if storageTotalGB, ok := systemInfo["storage_total_gb"].(float64); ok {
		agent.StorageTotalGB = storageTotalGB
	}
	if gpuModel, ok := systemInfo["gpu_model"].(string); ok {
		agent.GPUModel = gpuModel
	}
	if serialNumber, ok := systemInfo["serial_number"].(string); ok {
		agent.SerialNumber = serialNumber
	}
	if platform, ok := systemInfo["platform"].(string); ok {
		agent.Platform = platform
	}
	if macAddress, ok := systemInfo["mac_address"].(string); ok {
		agent.MACAddress = macAddress
	}
	if city, ok := systemInfo["city"].(string); ok {
		agent.City = city
	}
	if region, ok := systemInfo["region"].(string); ok {
		agent.Region = region
	}
	if country, ok := systemInfo["country"].(string); ok {
		agent.Country = country
	}
	if timezone, ok := systemInfo["timezone"].(string); ok {
		agent.Timezone = timezone
	}
	if riskScore, ok := systemInfo["risk_score"].(float64); ok {
		agent.RiskScore = riskScore
	}
	if tags, ok := systemInfo["tags"].([]string); ok {
		// Convert tags array to JSON string
		if tagsJSON, err := json.Marshal(tags); err == nil {
			agent.Tags = string(tagsJSON)
		}
	}

	agent.LastSeen = time.Now()
	agent.UpdatedAt = time.Now()

	// Update the agent in the map to ensure changes persist
	as.agents[agentUUID] = agent

	// Persist to DB
	if err := as.db.Save(agent).Error; err != nil {
		log.Printf("Failed to persist agent system info %s: %v", agent.ID, err)
	}

	log.Printf("[UpdateAgentSystemInfo] Updated system info for agent %s - hostname: %s, os: %s, cpu: %s", agentID, agent.Hostname, agent.OS, agent.CPUModel)
	return nil
}

// UpdateAgentMetadata updates agent metadata
func (as *AgentService) UpdateAgentMetadata(agentID string, metadata map[string]interface{}) error {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	agentUUID, err := uuid.Parse(agentID)
	if err != nil {
		return fmt.Errorf("invalid agent ID: %v", err)
	}

	agent, exists := as.agents[agentUUID]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	// Update metadata fields
	if osName, ok := metadata["os_name"].(string); ok {
		agent.OS = osName
	}
	if osVersion, ok := metadata["os_version"].(string); ok {
		agent.OSVersion = osVersion
	}
	if hostname, ok := metadata["hostname"].(string); ok {
		agent.Hostname = hostname
	}
	if ipAddress, ok := metadata["ip_address"].(string); ok {
		agent.IPAddress = ipAddress
	}
	if macAddress, ok := metadata["mac_address"].(string); ok {
		agent.MACAddress = macAddress
	}
	if serialNumber, ok := metadata["serial_number"].(string); ok {
		agent.SerialNumber = serialNumber
	}
	if city, ok := metadata["city"].(string); ok {
		agent.City = city
	}
	if region, ok := metadata["region"].(string); ok {
		agent.Region = region
	}
	if country, ok := metadata["country"].(string); ok {
		agent.Country = country
	}
	if timezone, ok := metadata["timezone"].(string); ok {
		agent.Timezone = timezone
	}
	if riskScore, ok := metadata["risk_score"].(float64); ok {
		agent.RiskScore = riskScore
	}
	if tags, ok := metadata["tags"].([]string); ok {
		// Convert tags array to JSON string
		if tagsJSON, err := json.Marshal(tags); err == nil {
			agent.Tags = string(tagsJSON)
		}
	}

	// Update hardware fields
	if cpuModel, ok := metadata["cpu_model"].(string); ok {
		agent.CPUModel = cpuModel
	}
	if cpuCores, ok := metadata["cpu_cores"].(int); ok {
		agent.CPUCores = cpuCores
	}
	if memoryTotalGB, ok := metadata["memory_total_gb"].(float64); ok {
		agent.MemoryTotalGB = memoryTotalGB
	}
	if storageTotalGB, ok := metadata["storage_total_gb"].(float64); ok {
		agent.StorageTotalGB = storageTotalGB
	}
	if gpuModel, ok := metadata["gpu_model"].(string); ok {
		agent.GPUModel = gpuModel
	}
	if platform, ok := metadata["platform"].(string); ok {
		agent.Platform = platform
	}
	if osBuild, ok := metadata["os_build"].(string); ok {
		agent.OSBuild = osBuild
	}
	if kernelVersion, ok := metadata["kernel_version"].(string); ok {
		agent.KernelVersion = kernelVersion
	}

	// Update the agent's metadata field with all the new information
	if agent.Metadata == nil {
		agent.Metadata = make(map[string]interface{})
	}

	// Merge new metadata with existing metadata
	for key, value := range metadata {
		agent.Metadata[key] = value
	}

	agent.LastSeen = time.Now()
	agent.UpdatedAt = time.Now()

	// PERSISTENCE: Save NetworkHost to Database
	if networkResults, ok := metadata["network_scan_result"].(map[string]interface{}); ok {
		go func(agentID uuid.UUID, results map[string]interface{}) {
			// Iterate over found hosts (assuming standard Nmap/Naabu structure)
			// This depends on the exact JSON structure of scan_result
			// For now, checks if there's a "hosts" array or similar

			// Example structure handling (adjust based on actual agent output)
			if hosts, ok := results["hosts"].([]interface{}); ok {
				for _, h := range hosts {
					if hostMap, ok := h.(map[string]interface{}); ok {
						ip, _ := hostMap["ip"].(string)
						if ip == "" {
							continue
						}

						hostname, _ := hostMap["hostname"].(string)
						// ... map other fields

						networkHost := models.NetworkHost{
							AgentID:   agentID,
							IPAddress: ip,
							Hostname:  hostname,
							Status:    "active",
							LastSeen:  time.Now(),
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}

						if ports, ok := hostMap["ports"].([]interface{}); ok {
							var openPorts []int
							for _, p := range ports {
								if pNum, ok := p.(float64); ok {
									openPorts = append(openPorts, int(pNum))
								}
							}
							networkHost.OpenPorts = openPorts
						}

						// Upsert
						if err := as.db.Where("agent_id = ? AND ip_address = ?", agentID, ip).
							FirstOrCreate(&networkHost).Error; err != nil {
							log.Printf("Failed to persist network host %s: %v", ip, err)
						}
					}
				}
			}
		}(agentUUID, networkResults)
	}

	log.Printf("[UpdateAgentMetadata] Updated metadata for agent %s", agentID)

	// Persist to DB (Update the agent record itself with new metadata)
	if err := as.db.Save(agent).Error; err != nil {
		log.Printf("Failed to persist agent metadata update %s: %v", agent.ID, err)
	}

	return nil
}
