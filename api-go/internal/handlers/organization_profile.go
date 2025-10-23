package handlers

import (
	"net/http"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrganizationProfileHandler handles organization profile API endpoints
type OrganizationProfileHandler struct {
	profileService *services.OrganizationProfileService
}

// NewOrganizationProfileHandler creates a new organization profile handler
func NewOrganizationProfileHandler(profileService *services.OrganizationProfileService) *OrganizationProfileHandler {
	return &OrganizationProfileHandler{
		profileService: profileService,
	}
}

// CreateOrganizationProfile creates a new organization profile
func (h *OrganizationProfileHandler) CreateOrganizationProfile(c *gin.Context) {
	var req models.CreateOrganizationProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request payload",
				Details: err.Error(),
			},
		})
		return
	}

	profile, err := h.profileService.CreateOrganizationProfile(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "CREATE_FAILED",
				Message: "Failed to create organization profile",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Data:    profile,
		Message: "Organization profile created successfully",
	})
}

// GetOrganizationProfile retrieves an organization profile
func (h *OrganizationProfileHandler) GetOrganizationProfile(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_UUID",
				Message: "Invalid organization ID format",
				Details: err.Error(),
			},
		})
		return
	}

	profile, err := h.profileService.GetOrganizationProfile(organizationID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization profile not found for organization "+organizationID.String() {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "PROFILE_NOT_FOUND",
				Message: "Failed to retrieve organization profile",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    profile,
	})
}

// UpdateOrganizationProfile updates an organization profile
func (h *OrganizationProfileHandler) UpdateOrganizationProfile(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_UUID",
				Message: "Invalid organization ID format",
				Details: err.Error(),
			},
		})
		return
	}

	var req models.UpdateOrganizationProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request payload",
				Details: err.Error(),
			},
		})
		return
	}

	profile, err := h.profileService.UpdateOrganizationProfile(organizationID, &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization profile not found for organization "+organizationID.String() {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "UPDATE_FAILED",
				Message: "Failed to update organization profile",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    profile,
		Message: "Organization profile updated successfully",
	})
}

// DeleteOrganizationProfile deletes an organization profile
func (h *OrganizationProfileHandler) DeleteOrganizationProfile(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_UUID",
				Message: "Invalid organization ID format",
				Details: err.Error(),
			},
		})
		return
	}

	err = h.profileService.DeleteOrganizationProfile(organizationID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization profile not found for organization "+organizationID.String() {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "DELETE_FAILED",
				Message: "Failed to delete organization profile",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Organization profile deleted successfully",
	})
}

// GetTechStackRelevance calculates tech stack relevance for a vulnerability
func (h *OrganizationProfileHandler) GetTechStackRelevance(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_UUID",
				Message: "Invalid organization ID format",
				Details: err.Error(),
			},
		})
		return
	}

	// Get vulnerability ID from query parameter
	vulnerabilityID := c.Query("vulnerability_id")
	if vulnerabilityID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "MISSING_VULNERABILITY_ID",
				Message: "vulnerability_id query parameter is required",
			},
		})
		return
	}

	// For now, we'll create a mock vulnerability - in real implementation,
	// this would fetch from the database
	vulnerability := &models.Vulnerability{
		ID:          vulnerabilityID,
		PackageName: c.Query("package_name"),
		Severity:    models.SeverityLevel(c.Query("severity")),
	}

	relevanceScore, err := h.profileService.GetTechStackRelevance(organizationID, vulnerability)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "RELEVANCE_CALCULATION_FAILED",
				Message: "Failed to calculate tech stack relevance",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"organization_id":  organizationID,
			"vulnerability_id": vulnerabilityID,
			"relevance_score":  relevanceScore,
			"package_name":     vulnerability.PackageName,
			"severity":         vulnerability.Severity,
		},
	})
}

// GetIndustryRiskWeights retrieves industry-specific risk weights
func (h *OrganizationProfileHandler) GetIndustryRiskWeights(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_UUID",
				Message: "Invalid organization ID format",
				Details: err.Error(),
			},
		})
		return
	}

	weights, err := h.profileService.GetIndustryRiskWeights(organizationID)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "organization profile not found for organization "+organizationID.String() {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "WEIGHTS_RETRIEVAL_FAILED",
				Message: "Failed to retrieve industry risk weights",
				Details: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"organization_id": organizationID,
			"risk_weights":    weights,
		},
	})
}
