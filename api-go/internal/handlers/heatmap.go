package handlers

import (
	"net/http"

	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HeatmapHandler handles risk heatmap API endpoints
type HeatmapHandler struct {
	heatmapService *services.HeatmapService
}

// NewHeatmapHandler creates a new HeatmapHandler
func NewHeatmapHandler(heatmapService *services.HeatmapService) *HeatmapHandler {
	return &HeatmapHandler{
		heatmapService: heatmapService,
	}
}

// GenerateRiskHeatmap generates a risk heatmap for an organization
func (h *HeatmapHandler) GenerateRiskHeatmap(c *gin.Context) {
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

	// Get query parameters
	heatmapType := c.DefaultQuery("type", "comprehensive")
	timeRange := c.DefaultQuery("range", "30d")

	// Generate heatmap
	heatmapData, err := h.heatmapService.GenerateRiskHeatmap(organizationID, heatmapType, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "HEATMAP_GENERATION_FAILED",
				"message": "Failed to generate risk heatmap",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    heatmapData,
	})
}

// GetHeatmapHotspots gets hotspots for an organization
func (h *HeatmapHandler) GetHeatmapHotspots(c *gin.Context) {
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

	// Generate heatmap to get hotspots
	heatmapData, err := h.heatmapService.GenerateRiskHeatmap(organizationID, "comprehensive", "30d")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "HEATMAP_GENERATION_FAILED",
				"message": "Failed to generate heatmap for hotspots",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id": organizationID,
			"hotspots":        heatmapData.Hotspots,
			"total_hotspots":  len(heatmapData.Hotspots),
			"generated_at":    heatmapData.GeneratedAt,
		},
	})
}

// GetRiskDistribution gets risk distribution for an organization
func (h *HeatmapHandler) GetRiskDistribution(c *gin.Context) {
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

	// Generate heatmap to get risk distribution
	heatmapData, err := h.heatmapService.GenerateRiskHeatmap(organizationID, "comprehensive", "30d")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "HEATMAP_GENERATION_FAILED",
				"message": "Failed to generate heatmap for risk distribution",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":   organizationID,
			"risk_distribution": heatmapData.RiskDistribution,
			"generated_at":      heatmapData.GeneratedAt,
		},
	})
}

// GetHeatmapTrends gets trends for an organization's risk heatmap
func (h *HeatmapHandler) GetHeatmapTrends(c *gin.Context) {
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

	timeRange := c.DefaultQuery("range", "30d")

	// Generate heatmap to get trends
	heatmapData, err := h.heatmapService.GenerateRiskHeatmap(organizationID, "comprehensive", timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "HEATMAP_GENERATION_FAILED",
				"message": "Failed to generate heatmap for trends",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id": organizationID,
			"trends":          heatmapData.Trends,
			"time_range":      timeRange,
			"generated_at":    heatmapData.GeneratedAt,
		},
	})
}

// GetHeatmapRecommendations gets recommendations based on heatmap analysis
func (h *HeatmapHandler) GetHeatmapRecommendations(c *gin.Context) {
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

	// Generate heatmap to get recommendations
	heatmapData, err := h.heatmapService.GenerateRiskHeatmap(organizationID, "comprehensive", "30d")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "HEATMAP_GENERATION_FAILED",
				"message": "Failed to generate heatmap for recommendations",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":  organizationID,
			"recommendations":  heatmapData.Recommendations,
			"total_hotspots":   len(heatmapData.Hotspots),
			"confidence_score": heatmapData.ConfidenceScore,
			"generated_at":     heatmapData.GeneratedAt,
		},
	})
}

