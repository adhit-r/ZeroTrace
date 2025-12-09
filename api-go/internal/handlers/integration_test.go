package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"zerotrace/api/internal/middleware"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestAgentEndpointsIntegration tests agent endpoints integration
func TestAgentEndpointsIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	agentService := services.NewAgentService(nil) // Mock DB
	router := gin.New()
	router.Use(middleware.CorrelationID())

	// Setup routes
	router.GET("/api/agents", GetAgents(agentService))
	router.GET("/api/agents/:id", GetAgent(agentService))
	router.POST("/api/agents/heartbeat", AgentHeartbeat(agentService))

	// Test GetAgents
	t.Run("GetAgents", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/agents", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Contains(t, w.Header().Get("X-Correlation-ID"), "")
	})

	// Test GetAgent with invalid ID
	t.Run("GetAgent_InvalidID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/agents/invalid-id", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.NotNil(t, response.Error)
	})

	// Test AgentHeartbeat with valid data
	t.Run("AgentHeartbeat_Valid", func(t *testing.T) {
		heartbeat := map[string]interface{}{
			"agent_id":        uuid.New().String(),
			"organization_id": uuid.New().String(),
			"agent_name":      "Test Agent",
			"status":          "online",
			"cpu_usage":       10.5,
			"memory_usage":    25.0,
			"metadata":        map[string]interface{}{},
			"timestamp":       "2024-01-01T00:00:00Z",
		}

		body, _ := json.Marshal(heartbeat)
		req := httptest.NewRequest("POST", "/api/agents/heartbeat", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test AgentHeartbeat with invalid JSON
	t.Run("AgentHeartbeat_InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/agents/heartbeat", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestErrorHandlingIntegration tests error handling across endpoints
func TestErrorHandlingIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CorrelationID())

	// Test correlation ID is added
	t.Run("CorrelationID_Added", func(t *testing.T) {
		router.GET("/test", func(c *gin.Context) {
			SuccessResponse(c, http.StatusOK, nil, "test")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		correlationID := w.Header().Get("X-Correlation-ID")
		assert.NotEmpty(t, correlationID)
	})

	// Test correlation ID from header is preserved
	t.Run("CorrelationID_Preserved", func(t *testing.T) {
		router.GET("/test", func(c *gin.Context) {
			SuccessResponse(c, http.StatusOK, nil, "test")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Correlation-ID", "test-correlation-id")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		correlationID := w.Header().Get("X-Correlation-ID")
		assert.Equal(t, "test-correlation-id", correlationID)
	})
}

// TestStandardizedErrorResponses tests standardized error response format
func TestStandardizedErrorResponses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(CorrelationID())

	t.Run("BadRequest", func(t *testing.T) {
		router.GET("/test", func(c *gin.Context) {
			BadRequest(c, "TEST_ERROR", "Test error message", "test details")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.NotNil(t, response.Error)
		assert.Equal(t, "TEST_ERROR", response.Error.Code)
		assert.NotEmpty(t, w.Header().Get("X-Correlation-ID"))
	})

	t.Run("NotFound", func(t *testing.T) {
		router.GET("/test", func(c *gin.Context) {
			NotFound(c, "NOT_FOUND", "Resource not found")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		router.GET("/test", func(c *gin.Context) {
			testErr := fmt.Errorf("test error")
			InternalServerError(c, "INTERNAL_ERROR", "Internal server error", testErr)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
