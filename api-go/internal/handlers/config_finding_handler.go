package handlers

import (
	"net/http"

	"zerotrace/api/internal/constants"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConfigFindingHandler handles config finding API endpoints
type ConfigFindingHandler struct {
	configFindingService *services.ConfigFindingService
}

// NewConfigFindingHandler creates a new config finding handler
func NewConfigFindingHandler(configFindingService *services.ConfigFindingService) *ConfigFindingHandler {
	return &ConfigFindingHandler{
		configFindingService: configFindingService,
	}
}

// ListConfigFindings lists config findings with filters
func (h *ConfigFindingHandler) ListConfigFindings(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	var req models.ListConfigFindingsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse config_file_id if provided
	if configFileIDStr := c.Query("config_file_id"); configFileIDStr != "" {
		if id, err := uuid.Parse(configFileIDStr); err == nil {
			req.ConfigFileID = &id
		}
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = constants.DefaultPageSize
	}
	if req.SortBy == "" {
		req.SortBy = "severity"
	}
	if req.SortOrder == "" {
		req.SortOrder = "DESC"
	}

	response, err := h.configFindingService.ListConfigFindings(companyID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetConfigFinding retrieves a finding by ID
func (h *ConfigFindingHandler) GetConfigFinding(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid finding ID"})
		return
	}

	finding, err := h.configFindingService.GetConfigFinding(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "finding not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    finding,
	})
}

// UpdateFindingStatus updates finding status
func (h *ConfigFindingHandler) UpdateFindingStatus(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid finding ID"})
		return
	}

	var req models.UpdateFindingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status enum value
	isValidStatus := false
	for _, valid := range constants.ValidFindingStatuses {
		if req.Status == valid {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status value"})
		return
	}

	// Get user ID if available
	var resolvedBy *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if userID, err := uuid.Parse(userIDStr.(string)); err == nil {
			resolvedBy = &userID
		}
	}

	err = h.configFindingService.UpdateFindingStatus(id, companyID, req.Status, req.AssignedTo, resolvedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Finding status updated successfully",
	})
}

// GetFindingStats retrieves finding statistics
func (h *ConfigFindingHandler) GetFindingStats(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	configFileIDStr := c.Query("config_file_id")
	if configFileIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "config_file_id is required"})
		return
	}

	configFileID, err := uuid.Parse(configFileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid config_file_id"})
		return
	}

	stats, err := h.configFindingService.GetFindingStats(configFileID, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

