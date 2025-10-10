package communicator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/models"
)

// Communicator handles communication with the API
type Communicator struct {
	config *config.Config
	client *http.Client
}

// NewCommunicator creates a new communicator instance
func NewCommunicator(cfg *config.Config) *Communicator {
	return &Communicator{
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.APITimeout) * time.Second,
		},
	}
}

// SendResults sends scan results to the API
func (c *Communicator) SendResults(result *models.ScanResult) error {
	// Prepare request payload
	payload := map[string]any{
		"agent_id": c.config.AgentID,
		"results":  []models.ScanResult{*result},
		"metadata": map[string]interface{}{
			"status": result.Status,
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal scan results: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/api/agents/results", c.config.APIEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ZeroTrace-Agent/1.0")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}

// SendStatus sends agent status to the API
func (c *Communicator) SendStatus(status *models.AgentStatus) error {
	// Prepare request payload
	payload := map[string]any{
		"agent_status": status,
		"timestamp":    time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal agent status: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/api/agents/status", c.config.APIEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("User-Agent", "ZeroTrace-Agent/1.0")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}

// GetScanTasks retrieves scan tasks from the API
func (c *Communicator) GetScanTasks() ([]map[string]any, error) {
	// Create request
	url := fmt.Sprintf("%s/api/v1/agent/tasks", c.config.APIEndpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("User-Agent", "ZeroTrace-Agent/1.0")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse response
	var response models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract tasks
	if response.Success && response.Data != nil {
		if tasks, ok := response.Data.([]map[string]any); ok {
			return tasks, nil
		}
	}

	return []map[string]any{}, nil
}

// SendHeartbeat sends agent heartbeat to the API
func (c *Communicator) SendHeartbeat(cpuUsage, memoryUsage float64, metadata map[string]any) error {
	// Prepare heartbeat payload
	heartbeat := map[string]any{
		"agent_id":     c.config.AgentID,
		"status":       "online",
		"cpu_usage":    cpuUsage,
		"memory_usage": memoryUsage,
		"metadata":     metadata,
		"timestamp":    time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/api/agents/heartbeat", c.config.APIEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("User-Agent", "ZeroTrace-Agent/1.0")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d for heartbeat", resp.StatusCode)
	}

	return nil
}

// RegisterAgent registers the agent with the API
func (c *Communicator) RegisterAgent() error {
	// Prepare registration payload
	registration := map[string]any{
		"agent_id": c.config.AgentID,
		"name":     "ZeroTrace Agent",
		"version":  "1.0.0",
		"hostname": c.config.Hostname,
		"os":       c.config.OS,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(registration)
	if err != nil {
		return fmt.Errorf("failed to marshal registration: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/api/agents/register", c.config.APIEndpoint)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create registration request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ZeroTrace-Agent/1.0")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send registration: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API returned status %d for registration", resp.StatusCode)
	}

	return nil
}

// EnrollAgent enrolls the agent using an enrollment token
func (c *Communicator) EnrollAgent() error {
	// Check if we have an enrollment token
	if !c.config.HasEnrollmentToken() {
		return fmt.Errorf("no enrollment token available")
	}

	// Prepare enrollment payload
	enrollment := map[string]any{
		"enrollment_token": c.config.EnrollmentToken,
		"agent_info": map[string]any{
			"hostname":     c.config.Hostname,
			"os":           c.config.OS,
			"version":      "1.0.0",
			"architecture": "unknown", // TODO: detect architecture
			"metadata": map[string]any{
				"company_id":   c.config.CompanyID,
				"company_name": c.config.CompanyName,
				"company_slug": c.config.CompanySlug,
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(enrollment)
	if err != nil {
		return fmt.Errorf("failed to marshal enrollment: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/api/enrollment/enroll", c.config.APIURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create enrollment request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ZeroTrace-Agent/1.0")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send enrollment: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d for enrollment", resp.StatusCode)
	}

	// Parse response
	var response struct {
		Success bool `json:"success"`
		Data    struct {
			AgentID        string    `json:"agent_id"`
			OrganizationID string    `json:"organization_id"`
			Credential     string    `json:"credential"`
			ExpiresAt      time.Time `json:"expires_at"`
		} `json:"data"`
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode enrollment response: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("enrollment failed: %s", response.Message)
	}

	// Update configuration with enrollment data
	c.config.AgentID = response.Data.AgentID
	c.config.OrganizationID = response.Data.OrganizationID
	c.config.AgentCredential = response.Data.Credential

	// Clear the enrollment token since it's been used
	c.config.EnrollmentToken = ""

	return nil
}

// SendHeartbeatWithCredential sends heartbeat using agent credential
func (c *Communicator) SendHeartbeatWithCredential(cpuUsage, memoryUsage float64, metadata map[string]any) error {
	// Check if we have a credential
	if !c.config.IsEnrolled() {
		return fmt.Errorf("agent not enrolled")
	}

	// Prepare heartbeat payload
	heartbeat := map[string]any{
		"agent_id":        c.config.AgentID,
		"organization_id": c.config.OrganizationID,
		"status":          "online",
		"cpu_usage":       cpuUsage,
		"memory_usage":    memoryUsage,
		"metadata":        metadata,
		"timestamp":       time.Now(),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		return fmt.Errorf("failed to marshal heartbeat: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/api/agents/heartbeat", c.config.APIURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create heartbeat request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.AgentCredential))
	req.Header.Set("User-Agent", "ZeroTrace-Agent/1.0")

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send heartbeat: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d for heartbeat", resp.StatusCode)
	}

	return nil
}
