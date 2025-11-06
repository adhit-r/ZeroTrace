package handlers

import (
	"net/http"

	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MaturityHandler handles security maturity score API endpoints
type MaturityHandler struct {
	maturityService *services.MaturityService
}

// NewMaturityHandler creates a new MaturityHandler
func NewMaturityHandler(maturityService *services.MaturityService) *MaturityHandler {
	return &MaturityHandler{
		maturityService: maturityService,
	}
}

// CalculateMaturityScore calculates security maturity score for an organization
func (h *MaturityHandler) CalculateMaturityScore(c *gin.Context) {
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

	// Calculate maturity score
	maturityScore, err := h.maturityService.CalculateMaturityScore(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MATURITY_CALCULATION_FAILED",
				"message": "Failed to calculate maturity score",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    maturityScore,
	})
}

// GetMaturityBenchmark gets industry benchmark for an organization
func (h *MaturityHandler) GetMaturityBenchmark(c *gin.Context) {
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

	// Get maturity score to extract benchmark data
	maturityScore, err := h.maturityService.CalculateMaturityScore(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MATURITY_CALCULATION_FAILED",
				"message": "Failed to calculate maturity score for benchmark",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":    organizationID,
			"industry_benchmark": maturityScore.IndustryBenchmark,
			"peer_comparison":    maturityScore.PeerComparison,
			"competitive_gap":    maturityScore.IndustryBenchmark.CompetitiveGap,
			"generated_at":       maturityScore.GeneratedAt,
		},
	})
}

// GetImprovementRoadmap gets improvement roadmap for an organization
func (h *MaturityHandler) GetImprovementRoadmap(c *gin.Context) {
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

	// Get maturity score to extract roadmap data
	maturityScore, err := h.maturityService.CalculateMaturityScore(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MATURITY_CALCULATION_FAILED",
				"message": "Failed to calculate maturity score for roadmap",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":     organizationID,
			"improvement_roadmap": maturityScore.ImprovementRoadmap,
			"total_improvements":  len(maturityScore.ImprovementRoadmap),
			"priority_items":      h.filterPriorityItems(maturityScore.ImprovementRoadmap),
			"generated_at":        maturityScore.GeneratedAt,
		},
	})
}

// GetMaturityTrends gets maturity trends for an organization
func (h *MaturityHandler) GetMaturityTrends(c *gin.Context) {
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

	// Get maturity score to extract trends data
	maturityScore, err := h.maturityService.CalculateMaturityScore(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MATURITY_CALCULATION_FAILED",
				"message": "Failed to calculate maturity score for trends",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id": organizationID,
			"trends":          maturityScore.Trends,
			"overall_score":   maturityScore.OverallScore,
			"maturity_level":  maturityScore.MaturityLevel,
			"generated_at":    maturityScore.GeneratedAt,
		},
	})
}

// GetDimensionScores gets detailed dimension scores for an organization
func (h *MaturityHandler) GetDimensionScores(c *gin.Context) {
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

	// Get maturity score to extract dimension data
	maturityScore, err := h.maturityService.CalculateMaturityScore(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "MATURITY_CALCULATION_FAILED",
				"message": "Failed to calculate maturity score for dimensions",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":  organizationID,
			"dimension_scores": maturityScore.DimensionScores,
			"overall_score":    maturityScore.OverallScore,
			"maturity_level":   maturityScore.MaturityLevel,
			"generated_at":     maturityScore.GeneratedAt,
		},
	})
}

// Helper method to filter priority items
func (h *MaturityHandler) filterPriorityItems(roadmap []services.ImprovementItem) []services.ImprovementItem {
	var priorityItems []services.ImprovementItem

	for _, item := range roadmap {
		if item.Priority == "High" {
			priorityItems = append(priorityItems, item)
		}
	}

	return priorityItems
}

