package repository

import (
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigFileRepository handles config file database operations
type ConfigFileRepository struct {
	db *gorm.DB
}

// NewConfigFileRepository creates a new config file repository
func NewConfigFileRepository(db *gorm.DB) *ConfigFileRepository {
	return &ConfigFileRepository{db: db}
}

// Create creates a new config file
func (r *ConfigFileRepository) Create(configFile *models.ConfigFile) error {
	configFile.ID = uuid.New()
	configFile.CreatedAt = time.Now()
	configFile.UpdatedAt = time.Now()
	return r.db.Create(configFile).Error
}

// GetByID retrieves a config file by ID
func (r *ConfigFileRepository) GetByID(id uuid.UUID) (*models.ConfigFile, error) {
	var configFile models.ConfigFile
	err := r.db.Where("id = ?", id).First(&configFile).Error
	if err != nil {
		return nil, err
	}
	return &configFile, nil
}

// GetByCompanyID retrieves config files by company ID with pagination and filters
func (r *ConfigFileRepository) GetByCompanyID(companyID uuid.UUID, page, limit int, filters map[string]interface{}) ([]models.ConfigFile, int64, error) {
	var configFiles []models.ConfigFile
	var total int64

	query := r.db.Model(&models.ConfigFile{}).Where("company_id = ?", companyID)

	// Apply filters
	if manufacturer, ok := filters["manufacturer"].(string); ok && manufacturer != "" {
		query = query.Where("manufacturer = ?", manufacturer)
	}
	if deviceType, ok := filters["device_type"].(string); ok && deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("analysis_status = ?", status)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	sortBy := "created_at"
	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := filters["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}

	err = query.Order(sortBy + " " + sortOrder).
		Offset(offset).
		Limit(limit).
		Find(&configFiles).Error

	return configFiles, total, err
}

// GetByHash retrieves a config file by hash and company ID (for deduplication)
func (r *ConfigFileRepository) GetByHash(hash string, companyID uuid.UUID) (*models.ConfigFile, error) {
	var configFile models.ConfigFile
	err := r.db.Where("file_hash = ? AND company_id = ?", hash, companyID).First(&configFile).Error
	if err != nil {
		return nil, err
	}
	return &configFile, nil
}

// UpdateStatus updates config file status
func (r *ConfigFileRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&models.ConfigFile{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"analysis_status": status,
			"updated_at":      time.Now(),
		}).Error
}

// UpdateParsingStatus updates parsing status and parsed data
func (r *ConfigFileRepository) UpdateParsingStatus(id uuid.UUID, status string, parsedData interface{}, parsingError string) error {
	updates := map[string]interface{}{
		"parsing_status": status,
		"updated_at":     time.Now(),
	}
	if parsedData != nil {
		updates["parsed_data"] = parsedData
	}
	if parsingError != "" {
		updates["parsing_error"] = parsingError
	}
	return r.db.Model(&models.ConfigFile{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// UpdateAnalysisStatus updates analysis status
func (r *ConfigFileRepository) UpdateAnalysisStatus(id uuid.UUID, status string) error {
	updates := map[string]interface{}{
		"analysis_status": status,
		"updated_at":      time.Now(),
	}
	if status == "analyzing" {
		now := time.Now()
		updates["analysis_started_at"] = now
	} else if status == "completed" || status == "failed" {
		now := time.Now()
		updates["analysis_completed_at"] = now
	}
	return r.db.Model(&models.ConfigFile{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete deletes a config file
func (r *ConfigFileRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.ConfigFile{}, id).Error
}

// GetStats retrieves config file statistics for a company
func (r *ConfigFileRepository) GetStats(companyID uuid.UUID) (map[string]interface{}, error) {
	var stats struct {
		TotalConfigs      int64 `json:"total_configs"`
		ParsedConfigs     int64 `json:"parsed_configs"`
		AnalyzedConfigs   int64 `json:"analyzed_configs"`
		FailedConfigs     int64 `json:"failed_configs"`
		PendingAnalysis   int64 `json:"pending_analysis"`
	}

	// Get total configs
	err := r.db.Model(&models.ConfigFile{}).Where("company_id = ?", companyID).Count(&stats.TotalConfigs).Error
	if err != nil {
		return nil, err
	}

	// Get parsed configs
	err = r.db.Model(&models.ConfigFile{}).Where("company_id = ? AND parsing_status = ?", companyID, "parsed").Count(&stats.ParsedConfigs).Error
	if err != nil {
		return nil, err
	}

	// Get analyzed configs
	err = r.db.Model(&models.ConfigFile{}).Where("company_id = ? AND analysis_status = ?", companyID, "completed").Count(&stats.AnalyzedConfigs).Error
	if err != nil {
		return nil, err
	}

	// Get failed configs
	err = r.db.Model(&models.ConfigFile{}).Where("company_id = ? AND analysis_status = ?", companyID, "failed").Count(&stats.FailedConfigs).Error
	if err != nil {
		return nil, err
	}

	// Get pending analysis
	err = r.db.Model(&models.ConfigFile{}).Where("company_id = ? AND analysis_status = ?", companyID, "pending").Count(&stats.PendingAnalysis).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_configs":    stats.TotalConfigs,
		"parsed_configs":   stats.ParsedConfigs,
		"analyzed_configs": stats.AnalyzedConfigs,
		"failed_configs":   stats.FailedConfigs,
		"pending_analysis": stats.PendingAnalysis,
	}, nil
}

