package analytics

import (
	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AnalyticsService provides unified analytics capabilities
// Consolidates heatmap, maturity, and compliance services
type AnalyticsService struct {
	db *gorm.DB
}

// NewAnalyticsService creates a new unified analytics service
func NewAnalyticsService(db *gorm.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

// GetVulnerabilitiesForOrganization retrieves vulnerabilities for analytics
func (s *AnalyticsService) GetVulnerabilitiesForOrganization(organizationID uuid.UUID) ([]models.Vulnerability, error) {
	var vulnerabilities []models.Vulnerability
	
	err := s.db.Where("organization_id = ?", organizationID).
		Order("severity DESC, created_at DESC").
		Find(&vulnerabilities).Error
	
	return vulnerabilities, err
}

// GetScanHistory retrieves scan history for analytics
func (s *AnalyticsService) GetScanHistory(organizationID uuid.UUID, limit int) ([]models.Scan, error) {
	var scans []models.Scan
	
	err := s.db.Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Limit(limit).
		Find(&scans).Error
	
	return scans, err
}

