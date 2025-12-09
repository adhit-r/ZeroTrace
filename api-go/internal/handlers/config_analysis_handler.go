package handlers

import (
	"net/http"

	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConfigAnalysisHandler handles config analysis API endpoints
type ConfigAnalysisHandler struct {
	configAnalysisService *services.ConfigAnalysisService
}

// NewConfigAnalysisHandler creates a new config analysis handler
func NewConfigAnalysisHandler(configAnalysisService *services.ConfigAnalysisService) *ConfigAnalysisHandler {
	return &ConfigAnalysisHandler{
		configAnalysisService: configAnalysisService,
	}
}

// GetAnalysisResults retrieves analysis results for a config file
func (h *ConfigAnalysisHandler) GetAnalysisResults(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config file ID"})
		return
	}

	result, err := h.configAnalysisService.GetAnalysisResults(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "analysis results not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetComplianceScores retrieves compliance scores for a config file
func (h *ConfigAnalysisHandler) GetComplianceScores(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config file ID"})
		return
	}

	scores, err := h.configAnalysisService.GetComplianceScores(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "compliance scores not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    scores,
	})
}

// GetAnalysisStatus retrieves analysis status for a config file
func (h *ConfigAnalysisHandler) GetAnalysisStatus(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config file ID"})
		return
	}

	status, err := h.configAnalysisService.GetAnalysisStatus(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"status": status,
		},
	})
}

