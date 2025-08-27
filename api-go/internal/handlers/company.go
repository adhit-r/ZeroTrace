package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCompany retrieves company information
func GetCompany(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_COMPANY_ID",
				Message: "Invalid company ID",
			},
			Timestamp: time.Now(),
		})
		return
	}

	// TODO: Implement actual company lookup
	company := &models.Company{
		ID:        companyID,
		Name:      "Example Company",
		Domain:    "example.com",
		Status:    "active",
		Settings:  make(map[string]any),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      company,
		Message:   "Company retrieved successfully",
		Timestamp: time.Now(),
	})
}

// UpdateCompany updates company information
func UpdateCompany(c *gin.Context) {
	companyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_COMPANY_ID",
				Message: "Invalid company ID",
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

	// TODO: Implement actual company update
	company := &models.Company{
		ID:        companyID,
		Name:      "Updated Company",
		Domain:    "updated.com",
		Status:    "active",
		Settings:  updates,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      company,
		Message:   "Company updated successfully",
		Timestamp: time.Now(),
	})
}
