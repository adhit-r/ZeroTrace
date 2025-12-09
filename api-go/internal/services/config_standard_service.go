package services

import (
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"
)

// ConfigStandardService handles config standard operations
type ConfigStandardService struct {
	configStandardRepo *repository.ConfigStandardRepository
}

// NewConfigStandardService creates a new config standard service
func NewConfigStandardService(configStandardRepo *repository.ConfigStandardRepository) *ConfigStandardService {
	return &ConfigStandardService{
		configStandardRepo: configStandardRepo,
	}
}

// GetStandardsForDevice retrieves standards for a specific manufacturer and device type
func (s *ConfigStandardService) GetStandardsForDevice(manufacturer, deviceType string) ([]models.ConfigStandard, error) {
	return s.configStandardRepo.GetByManufacturer(manufacturer, deviceType)
}

// GetStandardsByFramework retrieves standards by compliance framework
func (s *ConfigStandardService) GetStandardsByFramework(framework string) ([]models.ConfigStandard, error) {
	return s.configStandardRepo.GetByComplianceFramework(framework)
}

// CreateStandard creates a new standard
func (s *ConfigStandardService) CreateStandard(standard *models.ConfigStandard) error {
	return s.configStandardRepo.Create(standard)
}

// GetAllStandards retrieves all active standards
func (s *ConfigStandardService) GetAllStandards() ([]models.ConfigStandard, error) {
	return s.configStandardRepo.GetAll()
}

