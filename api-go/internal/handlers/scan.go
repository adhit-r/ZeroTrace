package handlers

import (
	"net/http"
	"strconv"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetScans retrieves scans with pagination
func GetScans(scanService *services.ScanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		scans, err := scanService.GetScans(companyUUID, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "SCAN_FETCH_FAILED",
					Message: "Failed to fetch scans",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      scans,
			Message:   "Scans retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// CreateScan creates a new scan
func CreateScan(scanService *services.ScanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateScanRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_REQUEST",
					Message: "Invalid request body",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		scan, err := scanService.CreateScan(req, companyUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "SCAN_CREATION_FAILED",
					Message: "Failed to create scan",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusCreated, models.APIResponse{
			Success:   true,
			Data:      scan,
			Message:   "Scan created successfully",
			Timestamp: time.Now(),
		})
	}
}

// GetScan retrieves a specific scan
func GetScan(scanService *services.ScanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scanID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_SCAN_ID",
					Message: "Invalid scan ID",
				},
				Timestamp: time.Now(),
			})
			return
		}

		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		scan, err := scanService.GetScan(scanID, companyUUID)
		if err != nil {
			c.JSON(http.StatusNotFound, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "SCAN_NOT_FOUND",
					Message: "Scan not found",
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      scan,
			Message:   "Scan retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// UpdateScan updates a scan
func UpdateScan(scanService *services.ScanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scanID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_SCAN_ID",
					Message: "Invalid scan ID",
				},
				Timestamp: time.Now(),
			})
			return
		}

		var updates map[string]any
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_REQUEST",
					Message: "Invalid request body",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		scan, err := scanService.UpdateScan(scanID, companyUUID, updates)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "SCAN_UPDATE_FAILED",
					Message: "Failed to update scan",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      scan,
			Message:   "Scan updated successfully",
			Timestamp: time.Now(),
		})
	}
}

// DeleteScan deletes a scan
func DeleteScan(scanService *services.ScanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		scanID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_SCAN_ID",
					Message: "Invalid scan ID",
				},
				Timestamp: time.Now(),
			})
			return
		}

		companyID, _ := c.Get("company_id")
		companyUUID, _ := uuid.Parse(companyID.(string))

		err = scanService.DeleteScan(scanID, companyUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "SCAN_DELETION_FAILED",
					Message: "Failed to delete scan",
					Details: err.Error(),
				},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Message:   "Scan deleted successfully",
			Timestamp: time.Now(),
		})
	}
}
