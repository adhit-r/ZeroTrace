package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"zerotrace/api/internal/constants"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ConfigFileHandler handles config file API endpoints
type ConfigFileHandler struct {
	configFileService *services.ConfigFileService
}

// NewConfigFileHandler creates a new config file handler
func NewConfigFileHandler(configFileService *services.ConfigFileService) *ConfigFileHandler {
	return &ConfigFileHandler{
		configFileService: configFileService,
	}
}

// UploadConfigFile handles config file upload
func (h *ConfigFileHandler) UploadConfigFile(c *gin.Context) {
	// Get company ID from context (set by middleware)
	companyIDStr, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "company_id not found in context"})
		return
	}

	companyID, err := uuid.Parse(companyIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid company_id"})
		return
	}

	// Get uploaded by (user ID) if available
	var uploadedBy *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if userID, err := uuid.Parse(userIDStr.(string)); err == nil {
			uploadedBy = &userID
		}
	}

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Sanitize filename to prevent path traversal attacks
	sanitizedFilename := filepath.Base(file.Filename)
	if sanitizedFilename == "." || sanitizedFilename == "/" || strings.Contains(sanitizedFilename, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filename"})
		return
	}

	// Validate file size (max 10MB)
	maxSize := int64(constants.MaxConfigFileSize)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB limit"})
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// Read file content
	fileContent, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// Parse form data
	var req models.UploadConfigFileRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Upload config file (use sanitized filename)
	configFile, err := h.configFileService.UploadConfigFile(
		fileContent,
		sanitizedFilename,
		req,
		companyID,
		uploadedBy,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Don't return file content in response
	configFile.FileContent = nil

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    configFile,
		"message": "Config file uploaded successfully",
	})
}

// GetConfigFile retrieves a config file by ID
func (h *ConfigFileHandler) GetConfigFile(c *gin.Context) {
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

	configFile, err := h.configFileService.GetConfigFile(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
		return
	}

	// Don't return file content in detail view
	configFile.FileContent = nil

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    configFile,
	})
}

// ListConfigFiles lists config files with filters
func (h *ConfigFileHandler) ListConfigFiles(c *gin.Context) {
	companyID, ok := getCompanyIDOrError(c)
	if !ok {
		return
	}

	var req models.ListConfigFilesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = constants.DefaultPageSize
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "DESC"
	}

	response, err := h.configFileService.ListConfigFiles(companyID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}

// GetConfigFileContent downloads the config file content
func (h *ConfigFileHandler) GetConfigFileContent(c *gin.Context) {
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

	configFile, err := h.configFileService.GetConfigFile(id, companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "config file not found"})
		return
	}

	content, err := h.configFileService.GetConfigFileContent(id, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set headers for file download (sanitize filename in header)
	safeFilename := filepath.Base(configFile.Filename)
	c.Header("Content-Disposition", "attachment; filename="+safeFilename)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", content)
}

// DeleteConfigFile deletes a config file
func (h *ConfigFileHandler) DeleteConfigFile(c *gin.Context) {
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

	err = h.configFileService.DeleteConfigFile(id, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Config file deleted successfully",
	})
}

// TriggerAnalysis manually triggers analysis for a config file
func (h *ConfigFileHandler) TriggerAnalysis(c *gin.Context) {
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

	err = h.configFileService.TriggerAnalysis(id, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analysis triggered successfully",
	})
}
