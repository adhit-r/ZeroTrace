package repository

import (
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigAnalysisRepository handles config analysis result database operations
type ConfigAnalysisRepository struct {
	db *gorm.DB
}

// NewConfigAnalysisRepository creates a new config analysis repository
func NewConfigAnalysisRepository(db *gorm.DB) *ConfigAnalysisRepository {
	return &ConfigAnalysisRepository{db: db}
}

// Create creates a new analysis result
func (r *ConfigAnalysisRepository) Create(result *models.ConfigAnalysisResult) error {
	result.ID = uuid.New()
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()
	return r.db.Create(result).Error
}

// GetByConfigFileID retrieves analysis result by config file ID
func (r *ConfigAnalysisRepository) GetByConfigFileID(configFileID uuid.UUID) (*models.ConfigAnalysisResult, error) {
	var result models.ConfigAnalysisResult
	err := r.db.Where("config_file_id = ?", configFileID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates an analysis result
func (r *ConfigAnalysisRepository) Update(result *models.ConfigAnalysisResult) error {
	result.UpdatedAt = time.Now()
	return r.db.Save(result).Error
}

// GetByID retrieves an analysis result by ID
func (r *ConfigAnalysisRepository) GetByID(id uuid.UUID) (*models.ConfigAnalysisResult, error) {
	var result models.ConfigAnalysisResult
	err := r.db.Where("id = ?", id).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes an analysis result
func (r *ConfigAnalysisRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ConfigAnalysisResult{}, id).Error
}

// DeleteByConfigFileID deletes analysis result by config file ID
func (r *ConfigAnalysisRepository) DeleteByConfigFileID(configFileID uuid.UUID) error {
	return r.db.Where("config_file_id = ?", configFileID).Delete(&models.ConfigAnalysisResult{}).Error
}

