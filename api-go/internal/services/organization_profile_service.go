package services

import (
	"fmt"
	"strings"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationProfileService handles organization profile operations
type OrganizationProfileService struct {
	db *gorm.DB
}

// NewOrganizationProfileService creates a new organization profile service
func NewOrganizationProfileService(db *gorm.DB) *OrganizationProfileService {
	return &OrganizationProfileService{
		db: db,
	}
}

// CreateOrganizationProfile creates a new organization profile
func (s *OrganizationProfileService) CreateOrganizationProfile(req *models.CreateOrganizationProfileRequest) (*models.OrganizationProfile, error) {
	// Check if profile already exists
	var existingProfile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", req.OrganizationID).First(&existingProfile).Error
	if err == nil {
		return nil, fmt.Errorf("organization profile already exists for organization %s", req.OrganizationID)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing profile: %w", err)
	}

	// Create new profile
	profile := &models.OrganizationProfile{
		ID:                   uuid.New(),
		OrganizationID:       req.OrganizationID,
		Industry:             req.Industry,
		RiskTolerance:        req.RiskTolerance,
		TechStack:            req.TechStack,
		ComplianceFrameworks: req.ComplianceFrameworks,
		SecurityPolicies:     req.SecurityPolicies,
		RiskWeights:          req.RiskWeights,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Set default risk weights if not provided
	if profile.RiskWeights == nil {
		profile.RiskWeights = s.getDefaultRiskWeights(req.Industry, req.RiskTolerance)
	}

	// Set default security policies if not provided
	if profile.SecurityPolicies == nil {
		profile.SecurityPolicies = s.getDefaultSecurityPolicies(req.Industry)
	}

	err = s.db.Create(profile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create organization profile: %w", err)
	}

	return profile, nil
}

// GetOrganizationProfile retrieves an organization profile by organization ID
func (s *OrganizationProfileService) GetOrganizationProfile(organizationID uuid.UUID) (*models.OrganizationProfile, error) {
	var profile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization profile not found for organization %s", organizationID)
		}
		return nil, fmt.Errorf("failed to get organization profile: %w", err)
	}

	return &profile, nil
}

// UpdateOrganizationProfile updates an existing organization profile
func (s *OrganizationProfileService) UpdateOrganizationProfile(organizationID uuid.UUID, req *models.UpdateOrganizationProfileRequest) (*models.OrganizationProfile, error) {
	var profile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization profile not found for organization %s", organizationID)
		}
		return nil, fmt.Errorf("failed to get organization profile: %w", err)
	}

	// Update fields if provided
	updates := make(map[string]interface{})

	if req.Industry != nil {
		updates["industry"] = *req.Industry
	}
	if req.RiskTolerance != nil {
		updates["risk_tolerance"] = *req.RiskTolerance
	}
	if req.TechStack != nil {
		updates["tech_stack"] = *req.TechStack
	}
	if req.ComplianceFrameworks != nil {
		updates["compliance_frameworks"] = *req.ComplianceFrameworks
	}
	if req.SecurityPolicies != nil {
		updates["security_policies"] = req.SecurityPolicies
	}
	if req.RiskWeights != nil {
		updates["risk_weights"] = req.RiskWeights
	}

	updates["updated_at"] = time.Now()

	err = s.db.Model(&profile).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update organization profile: %w", err)
	}

	// Reload the profile to get updated data
	err = s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to reload organization profile: %w", err)
	}

	return &profile, nil
}

// DeleteOrganizationProfile deletes an organization profile
func (s *OrganizationProfileService) DeleteOrganizationProfile(organizationID uuid.UUID) error {
	result := s.db.Where("organization_id = ?", organizationID).Delete(&models.OrganizationProfile{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete organization profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization profile not found for organization %s", organizationID)
	}

	return nil
}

// GetTechStackRelevance calculates relevance score for vulnerabilities based on tech stack
func (s *OrganizationProfileService) GetTechStackRelevance(organizationID uuid.UUID, vulnerability *models.Vulnerability) (float64, error) {
	profile, err := s.GetOrganizationProfile(organizationID)
	if err != nil {
		return 0, err
	}

	// Calculate relevance based on tech stack
	relevanceScore := 0.0
	techStack := profile.TechStack

	// Check if vulnerability affects technologies in the organization's stack
	if vulnerability.PackageName != "" {
		// Check against languages
		for _, lang := range techStack.Languages {
			if s.isTechnologyMatch(vulnerability.PackageName, lang) {
				relevanceScore += 0.3
			}
		}

		// Check against frameworks
		for _, framework := range techStack.Frameworks {
			if s.isTechnologyMatch(vulnerability.PackageName, framework) {
				relevanceScore += 0.4
			}
		}

		// Check against databases
		for _, db := range techStack.Databases {
			if s.isTechnologyMatch(vulnerability.PackageName, db) {
				relevanceScore += 0.3
			}
		}
	}

	// Apply risk tolerance multiplier
	switch profile.RiskTolerance {
	case models.RiskToleranceConservative:
		relevanceScore *= 1.2 // Higher relevance for conservative organizations
	case models.RiskToleranceModerate:
		relevanceScore *= 1.0 // No change
	case models.RiskToleranceAggressive:
		relevanceScore *= 0.8 // Lower relevance for aggressive organizations
	}

	// Cap at 1.0
	if relevanceScore > 1.0 {
		relevanceScore = 1.0
	}

	return relevanceScore, nil
}

// GetIndustryRiskWeights returns industry-specific risk weights
func (s *OrganizationProfileService) GetIndustryRiskWeights(organizationID uuid.UUID) (map[string]float64, error) {
	profile, err := s.GetOrganizationProfile(organizationID)
	if err != nil {
		return nil, err
	}

	// Convert risk weights to float64 map
	weights := make(map[string]float64)
	for key, value := range profile.RiskWeights {
		if floatVal, ok := value.(float64); ok {
			weights[key] = floatVal
		}
	}

	return weights, nil
}

// isTechnologyMatch checks if a vulnerability package matches a technology
func (s *OrganizationProfileService) isTechnologyMatch(packageName, technology string) bool {
	// Simple string matching - can be enhanced with fuzzy matching
	packageName = strings.ToLower(packageName)
	technology = strings.ToLower(technology)

	return strings.Contains(packageName, technology) || strings.Contains(technology, packageName)
}

// getDefaultRiskWeights returns default risk weights based on industry and risk tolerance
func (s *OrganizationProfileService) getDefaultRiskWeights(industry string, riskTolerance models.RiskTolerance) map[string]any {
	baseWeights := map[string]float64{
		"critical": 1.0,
		"high":     0.8,
		"medium":   0.6,
		"low":      0.4,
		"info":     0.2,
	}

	// Industry-specific adjustments
	switch industry {
	case "healthcare":
		baseWeights["critical"] = 1.2
		baseWeights["high"] = 1.0
	case "finance":
		baseWeights["critical"] = 1.1
		baseWeights["high"] = 0.9
	case "government":
		baseWeights["critical"] = 1.3
		baseWeights["high"] = 1.1
	}

	// Risk tolerance adjustments
	switch riskTolerance {
	case models.RiskToleranceConservative:
		for key := range baseWeights {
			baseWeights[key] *= 1.2
		}
	case models.RiskToleranceAggressive:
		for key := range baseWeights {
			baseWeights[key] *= 0.8
		}
	}

	// Convert to interface{} map
	weights := make(map[string]any)
	for key, value := range baseWeights {
		weights[key] = value
	}

	return weights
}

// getDefaultSecurityPolicies returns default security policies based on industry
func (s *OrganizationProfileService) getDefaultSecurityPolicies(industry string) map[string]any {
	policies := map[string]any{
		"patch_management": map[string]any{
			"critical_patches": "immediate",
			"high_patches":     "24_hours",
			"medium_patches":   "7_days",
			"low_patches":      "30_days",
		},
		"vulnerability_management": map[string]any{
			"scan_frequency": "daily",
			"reporting":      "real_time",
		},
		"compliance": map[string]any{
			"enabled":    true,
			"frameworks": []string{},
		},
	}

	// Industry-specific policies
	switch industry {
	case "healthcare":
		policies["compliance"] = map[string]any{
			"enabled":    true,
			"frameworks": []string{"HIPAA", "HITECH"},
		}
		policies["data_protection"] = map[string]any{
			"encryption_required": true,
			"access_controls":     "strict",
		}
	case "finance":
		policies["compliance"] = map[string]any{
			"enabled":    true,
			"frameworks": []string{"PCI DSS", "SOX", "Basel III"},
		}
		policies["audit_requirements"] = map[string]any{
			"logging":    "comprehensive",
			"retention":  "7_years",
			"monitoring": "continuous",
		}
	case "government":
		policies["compliance"] = map[string]any{
			"enabled":    true,
			"frameworks": []string{"FISMA", "FedRAMP", "NIST"},
		}
		policies["security_clearance"] = map[string]any{
			"required": true,
			"levels":   []string{"public", "confidential", "secret"},
		}
	}

	return policies
}
