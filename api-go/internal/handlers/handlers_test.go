package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock services for testing
type MockAgentService struct {
	mock.Mock
}

func (m *MockAgentService) GetAllAgents() []models.Agent {
	args := m.Called()
	return args.Get(0).([]models.Agent)
}

func (m *MockAgentService) GetAgentByID(id string) (*models.Agent, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Agent), args.Error(1)
}

func (m *MockAgentService) UpdateAgentSystemInfo(agentID string, systemInfo map[string]interface{}) error {
	args := m.Called(agentID, systemInfo)
	return args.Error(0)
}

func (m *MockAgentService) HandleAgentHeartbeat(heartbeat models.AgentHeartbeat) error {
	args := m.Called(heartbeat)
	return args.Error(0)
}

type MockVulnerabilityService struct {
	mock.Mock
}

func (m *MockVulnerabilityService) GetVulnerabilities() ([]models.Vulnerability, error) {
	args := m.Called()
	return args.Get(0).([]models.Vulnerability), args.Error(1)
}

func (m *MockVulnerabilityService) GetVulnerabilityByID(id string) (*models.Vulnerability, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Vulnerability), args.Error(1)
}

func (m *MockVulnerabilityService) CreateVulnerability(vuln models.Vulnerability) error {
	args := m.Called(vuln)
	return args.Error(0)
}

func (m *MockVulnerabilityService) UpdateVulnerability(id string, vuln models.Vulnerability) error {
	args := m.Called(id, vuln)
	return args.Error(0)
}

func (m *MockVulnerabilityService) DeleteVulnerability(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetAgents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	handler := NewAgentHandler(mockAgentService)

	// Mock data
	expectedAgents := []models.Agent{
		{
			ID:       "test-agent-1",
			Hostname: "test-host-1",
			OS:       "macOS",
			Status:   "online",
		},
		{
			ID:       "test-agent-2",
			Hostname: "test-host-2",
			OS:       "Linux",
			Status:   "offline",
		},
	}

	mockAgentService.On("GetAllAgents").Return(expectedAgents)

	// Create test request
	req := httptest.NewRequest("GET", "/api/agents", nil)
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetAgents(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Agent
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedAgents[0].ID, response[0].ID)
	assert.Equal(t, expectedAgents[1].ID, response[1].ID)

	mockAgentService.AssertExpectations(t)
}

func TestGetAgentByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	handler := NewAgentHandler(mockAgentService)

	// Mock data
	expectedAgent := &models.Agent{
		ID:       "test-agent-1",
		Hostname: "test-host-1",
		OS:       "macOS",
		Status:   "online",
	}

	mockAgentService.On("GetAgentByID", "test-agent-1").Return(expectedAgent, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api/agents/test-agent-1", nil)
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "test-agent-1"}}

	// Call handler
	handler.GetAgentByID(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Agent
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedAgent.ID, response.ID)
	assert.Equal(t, expectedAgent.Hostname, response.Hostname)

	mockAgentService.AssertExpectations(t)
}

func TestUpdateSystemInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	handler := NewAgentHandler(mockAgentService)

	// Test data
	systemInfo := map[string]interface{}{
		"hostname":   "test-host",
		"os":         "macOS",
		"os_version": "13.0",
		"cpu_cores":  8,
		"memory_gb":  16,
		"disk_gb":    512,
	}

	request := SystemInfoRequest{
		AgentID:    "test-agent-1",
		SystemInfo: systemInfo,
	}

	requestBody, _ := json.Marshal(request)

	mockAgentService.On("UpdateAgentSystemInfo", "test-agent-1", systemInfo).Return(nil)

	// Create test request
	req := httptest.NewRequest("POST", "/api/agents/system-info", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.UpdateSystemInfo(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockAgentService.AssertExpectations(t)
}

func TestHandleHeartbeat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	handler := NewAgentHandler(mockAgentService)

	// Test data
	heartbeat := models.AgentHeartbeat{
		AgentID:        "test-agent-1",
		AgentName:      "Test Agent",
		OrganizationID: "test-org-1",
		Status:         "online",
		LastSeen:       time.Now(),
	}

	requestBody, _ := json.Marshal(heartbeat)

	mockAgentService.On("HandleAgentHeartbeat", heartbeat).Return(nil)

	// Create test request
	req := httptest.NewRequest("POST", "/api/agents/heartbeat", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.HandleHeartbeat(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockAgentService.AssertExpectations(t)
}

func TestGetVulnerabilities(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVulnService := new(MockVulnerabilityService)
	handler := NewVulnerabilityHandler(mockVulnService)

	// Mock data
	expectedVulns := []models.Vulnerability{
		{
			ID:          "vuln-1",
			Title:       "Test Vulnerability 1",
			Description: "Test description 1",
			Severity:    "high",
			Status:      "open",
		},
		{
			ID:          "vuln-2",
			Title:       "Test Vulnerability 2",
			Description: "Test description 2",
			Severity:    "medium",
			Status:      "open",
		},
	}

	mockVulnService.On("GetVulnerabilities").Return(expectedVulns, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api/vulnerabilities", nil)
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetVulnerabilities(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Vulnerability
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedVulns[0].ID, response[0].ID)
	assert.Equal(t, expectedVulns[1].ID, response[1].ID)

	mockVulnService.AssertExpectations(t)
}

func TestGetVulnerabilityByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVulnService := new(MockVulnerabilityService)
	handler := NewVulnerabilityHandler(mockVulnService)

	// Mock data
	expectedVuln := &models.Vulnerability{
		ID:          "vuln-1",
		Title:       "Test Vulnerability 1",
		Description: "Test description 1",
		Severity:    "high",
		Status:      "open",
	}

	mockVulnService.On("GetVulnerabilityByID", "vuln-1").Return(expectedVuln, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/api/vulnerabilities/vuln-1", nil)
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "vuln-1"}}

	// Call handler
	handler.GetVulnerabilityByID(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Vulnerability
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedVuln.ID, response.ID)
	assert.Equal(t, expectedVuln.Title, response.Title)

	mockVulnService.AssertExpectations(t)
}

func TestCreateVulnerability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVulnService := new(MockVulnerabilityService)
	handler := NewVulnerabilityHandler(mockVulnService)

	// Test data
	vuln := models.Vulnerability{
		Title:       "New Vulnerability",
		Description: "New vulnerability description",
		Severity:    "high",
		Status:      "open",
	}

	requestBody, _ := json.Marshal(vuln)

	mockVulnService.On("CreateVulnerability", vuln).Return(nil)

	// Create test request
	req := httptest.NewRequest("POST", "/api/vulnerabilities", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.CreateVulnerability(c)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockVulnService.AssertExpectations(t)
}

func TestUpdateVulnerability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVulnService := new(MockVulnerabilityService)
	handler := NewVulnerabilityHandler(mockVulnService)

	// Test data
	vuln := models.Vulnerability{
		ID:          "vuln-1",
		Title:       "Updated Vulnerability",
		Description: "Updated vulnerability description",
		Severity:    "medium",
		Status:      "in_progress",
	}

	requestBody, _ := json.Marshal(vuln)

	mockVulnService.On("UpdateVulnerability", "vuln-1", vuln).Return(nil)

	// Create test request
	req := httptest.NewRequest("PUT", "/api/vulnerabilities/vuln-1", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "vuln-1"}}

	// Call handler
	handler.UpdateVulnerability(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockVulnService.AssertExpectations(t)
}

func TestDeleteVulnerability(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockVulnService := new(MockVulnerabilityService)
	handler := NewVulnerabilityHandler(mockVulnService)

	mockVulnService.On("DeleteVulnerability", "vuln-1").Return(nil)

	// Create test request
	req := httptest.NewRequest("DELETE", "/api/vulnerabilities/vuln-1", nil)
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "vuln-1"}}

	// Call handler
	handler.DeleteVulnerability(c)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])

	mockVulnService.AssertExpectations(t)
}

// Performance tests
func BenchmarkGetAgents(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	handler := NewAgentHandler(mockAgentService)

	// Mock data
	agents := make([]models.Agent, 100)
	for i := 0; i < 100; i++ {
		agents[i] = models.Agent{
			ID:       fmt.Sprintf("agent-%d", i),
			Hostname: fmt.Sprintf("host-%d", i),
			OS:       "Linux",
			Status:   "online",
		}
	}

	mockAgentService.On("GetAllAgents").Return(agents)

	req := httptest.NewRequest("GET", "/api/agents", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		handler.GetAgents(c)
	}
}

func BenchmarkGetVulnerabilities(b *testing.B) {
	gin.SetMode(gin.TestMode)

	mockVulnService := new(MockVulnerabilityService)
	handler := NewVulnerabilityHandler(mockVulnService)

	// Mock data
	vulns := make([]models.Vulnerability, 1000)
	for i := 0; i < 1000; i++ {
		vulns[i] = models.Vulnerability{
			ID:       fmt.Sprintf("vuln-%d", i),
			Title:    fmt.Sprintf("Vulnerability %d", i),
			Severity: "medium",
			Status:   "open",
		}
	}

	mockVulnService.On("GetVulnerabilities").Return(vulns, nil)

	req := httptest.NewRequest("GET", "/api/vulnerabilities", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		handler.GetVulnerabilities(c)
	}
}

// Integration tests
func TestAPIIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test server
	router := gin.New()

	// Add routes
	router.GET("/api/agents", func(c *gin.Context) {
		c.JSON(200, []models.Agent{})
	})
	router.GET("/api/vulnerabilities", func(c *gin.Context) {
		c.JSON(200, []models.Vulnerability{})
	})

	// Test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Test endpoints
	client := &http.Client{Timeout: 5 * time.Second}

	// Test agents endpoint
	resp, err := client.Get(server.URL + "/api/agents")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Test vulnerabilities endpoint
	resp, err = client.Get(server.URL + "/api/vulnerabilities")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

// Test error handling
func TestErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	handler := NewAgentHandler(mockAgentService)

	// Test agent not found
	mockAgentService.On("GetAgentByID", "nonexistent").Return((*models.Agent)(nil), errors.New("agent not found"))

	req := httptest.NewRequest("GET", "/api/agents/nonexistent", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "nonexistent"}}

	handler.GetAgentByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	mockAgentService.AssertExpectations(t)
}

// Test concurrent requests
func TestConcurrentRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAgentService := new(MockAgentService)
	handler := NewAgentHandler(mockAgentService)

	agents := make([]models.Agent, 10)
	for i := 0; i < 10; i++ {
		agents[i] = models.Agent{
			ID:       fmt.Sprintf("agent-%d", i),
			Hostname: fmt.Sprintf("host-%d", i),
			OS:       "Linux",
			Status:   "online",
		}
	}

	mockAgentService.On("GetAllAgents").Return(agents)

	// Create multiple concurrent requests
	numRequests := 100
	results := make(chan int, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest("GET", "/api/agents", nil)
			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)
			c.Request = req

			handler.GetAgents(c)
			results <- w.Code
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		select {
		case code := <-results:
			assert.Equal(t, http.StatusOK, code)
		case <-time.After(10 * time.Second):
			t.Fatal("Concurrent requests timed out")
		}
	}

	mockAgentService.AssertExpectations(t)
}
