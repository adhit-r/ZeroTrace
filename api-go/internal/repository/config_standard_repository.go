package repository

import (
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigStandardRepository handles config standard database operations
type ConfigStandardRepository struct {
	db *gorm.DB
}

// NewConfigStandardRepository creates a new config standard repository
func NewConfigStandardRepository(db *gorm.DB) *ConfigStandardRepository {
	return &ConfigStandardRepository{db: db}
}

// GetByManufacturer retrieves standards by manufacturer and device type
func (r *ConfigStandardRepository) GetByManufacturer(manufacturer, deviceType string) ([]models.ConfigStandard, error) {
	var standards []models.ConfigStandard
	err := r.db.Where("manufacturer = ? AND device_type = ? AND status = ?", manufacturer, deviceType, "active").
		Order("requirement_id ASC").
		Find(&standards).Error
	return standards, err
}

// GetByComplianceFramework retrieves standards by compliance framework
func (r *ConfigStandardRepository) GetByComplianceFramework(framework string) ([]models.ConfigStandard, error) {
	var standards []models.ConfigStandard
	err := r.db.Where("compliance_frameworks @> ? AND status = ?", `["`+framework+`"]`, "active").
		Order("manufacturer, device_type, requirement_id ASC").
		Find(&standards).Error
	return standards, err
}

// GetByID retrieves a standard by ID
func (r *ConfigStandardRepository) GetByID(id uuid.UUID) (*models.ConfigStandard, error) {
	var standard models.ConfigStandard
	err := r.db.Where("id = ?", id).First(&standard).Error
	if err != nil {
		return nil, err
	}
	return &standard, nil
}

// Create creates a new standard
func (r *ConfigStandardRepository) Create(standard *models.ConfigStandard) error {
	standard.ID = uuid.New()
	standard.CreatedAt = time.Now()
	standard.UpdatedAt = time.Now()
	return r.db.Create(standard).Error
}

// CreateBatch creates multiple standards in a single transaction
func (r *ConfigStandardRepository) CreateBatch(standards []models.ConfigStandard) error {
	if len(standards) == 0 {
		return nil
	}

	now := time.Now()
	for i := range standards {
		standards[i].ID = uuid.New()
		standards[i].CreatedAt = now
		standards[i].UpdatedAt = now
	}

	return r.db.CreateInBatches(standards, 100).Error
}

// GetAll retrieves all active standards
func (r *ConfigStandardRepository) GetAll() ([]models.ConfigStandard, error) {
	var standards []models.ConfigStandard
	err := r.db.Where("status = ?", "active").
		Order("manufacturer, device_type, requirement_id ASC").
		Find(&standards).Error
	return standards, err
}

// Update updates a standard
func (r *ConfigStandardRepository) Update(standard *models.ConfigStandard) error {
	standard.UpdatedAt = time.Now()
	return r.db.Save(standard).Error
}

// Delete deletes a standard
func (r *ConfigStandardRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ConfigStandard{}, id).Error
}

