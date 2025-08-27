package services

import (
	"sync"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
)

// AgentService manages agent registration and heartbeats
type AgentService struct {
	agents map[uuid.UUID]*models.Agent
	mutex  sync.RWMutex
}

// NewAgentService creates a new agent service
func NewAgentService() *AgentService {
	as := &AgentService{
		agents: make(map[uuid.UUID]*models.Agent),
	}
	as.StartCleanupRoutine()
	return as
}

// RegisterAgent registers a new agent
func (as *AgentService) RegisterAgent(agentID uuid.UUID, organizationID uuid.UUID, name, version, hostname, os string) *models.Agent {
	as.mutex.Lock()
	defer as.mutex.Unlock()

	agent := &models.Agent{
		ID:             agentID,
		OrganizationID: organizationID,
		Name:           name,
		Status:         "active",
		Version:        version,
		LastSeen:       time.Now(),
		Hostname:       hostname,
		OS:             os,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	as.agents[agentID] = agent
	return agent
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
			Name:           "Unknown Agent",
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
	agent.Metadata = heartbeat.Metadata
	agent.UpdatedAt = time.Now()

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

