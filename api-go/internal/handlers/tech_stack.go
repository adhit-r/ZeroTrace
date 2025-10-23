package handlers

import (
	"net/http"

	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TechStackHandler handles technology stack analysis API endpoints
type TechStackHandler struct {
	techStackService *services.TechStackService
}

// NewTechStackHandler creates a new tech stack handler
func NewTechStackHandler(techStackService *services.TechStackService) *TechStackHandler {
	return &TechStackHandler{
		techStackService: techStackService,
	}
}

// AnalyzeTechStack analyzes technology stack for an organization
func (h *TechStackHandler) AnalyzeTechStack(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	analysis, err := h.techStackService.AnalyzeTechStackFromAssets(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ANALYSIS_FAILED",
				"message": "Failed to analyze technology stack",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    analysis,
	})
}

// GetTechStackRecommendations gets security recommendations based on tech stack
func (h *TechStackHandler) GetTechStackRecommendations(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	analysis, err := h.techStackService.AnalyzeTechStackFromAssets(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "ANALYSIS_FAILED",
				"message": "Failed to analyze technology stack",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":  organizationID,
			"recommendations":  analysis.Recommendations,
			"risk_factors":     analysis.RiskFactors,
			"confidence_score": analysis.ConfidenceScore,
		},
	})
}
