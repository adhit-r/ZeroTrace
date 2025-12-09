package handlers

import (
	"net/http"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
)

type AttackPathHandler struct {
	attackPathService *services.AttackPathService
}

func NewAttackPathHandler(service *services.AttackPathService) *AttackPathHandler {
	return &AttackPathHandler{
		attackPathService: service,
	}
}

// GetAttackPaths retrieves all attack paths for an organization
func (h *AttackPathHandler) GetAttackPaths(c *gin.Context) {
	organizationID := c.Query("organization_id")
	if organizationID == "" {
		// Use default organization for now
		organizationID = "00000000-0000-0000-0000-000000000001"
	}

	paths, err := h.attackPathService.GetAttackPaths(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve attack paths",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    paths,
		"count":   len(paths),
	})
}

// GetAttackPath retrieves a specific attack path
func (h *AttackPathHandler) GetAttackPath(c *gin.Context) {
	pathID := c.Param("path_id")
	if pathID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path ID is required"})
		return
	}

	path, err := h.attackPathService.GetAttackPath(pathID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Attack path not found",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    path,
	})
}

// GenerateAttackPaths generates attack paths from current vulnerabilities and network data
func (h *AttackPathHandler) GenerateAttackPaths(c *gin.Context) {
	organizationID := c.Query("organization_id")
	if organizationID == "" {
		organizationID = "00000000-0000-0000-0000-000000000001"
	}

	paths, err := h.attackPathService.GenerateAttackPaths(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate attack paths",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    paths,
		"count":   len(paths),
		"message": "Attack paths generated successfully",
	})
}

