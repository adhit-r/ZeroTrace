package repository

import (
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigFindingRepository handles config finding database operations
type ConfigFindingRepository struct {
	db *gorm.DB
}

// NewConfigFindingRepository creates a new config finding repository
func NewConfigFindingRepository(db *gorm.DB) *ConfigFindingRepository {
	return &ConfigFindingRepository{db: db}
}

// Create creates a new config finding
func (r *ConfigFindingRepository) Create(finding *models.ConfigFinding) error {
	finding.ID = uuid.New()
	finding.CreatedAt = time.Now()
	finding.UpdatedAt = time.Now()
	return r.db.Create(finding).Error
}

// CreateBatch creates multiple config findings in a single transaction
func (r *ConfigFindingRepository) CreateBatch(findings []models.ConfigFinding) error {
	if len(findings) == 0 {
		return nil
	}

	now := time.Now()
	for i := range findings {
		findings[i].ID = uuid.New()
		findings[i].CreatedAt = now
		findings[i].UpdatedAt = now
	}

	return r.db.CreateInBatches(findings, 100).Error
}

// GetByConfigFileID retrieves findings by config file ID with filters
func (r *ConfigFindingRepository) GetByConfigFileID(configFileID uuid.UUID, filters map[string]interface{}) ([]models.ConfigFinding, error) {
	var findings []models.ConfigFinding

	query := r.db.Where("config_file_id = ?", configFileID)

	// Apply filters
	if severity, ok := filters["severity"].(string); ok && severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if findingType, ok := filters["finding_type"].(string); ok && findingType != "" {
		query = query.Where("finding_type = ?", findingType)
	}

	sortBy := "severity"
	if sb, ok := filters["sort_by"].(string); ok && sb != "" {
		sortBy = sb
	}
	sortOrder := "DESC"
	if so, ok := filters["sort_order"].(string); ok && so != "" {
		sortOrder = so
	}

	err := query.Order(sortBy + " " + sortOrder).Find(&findings).Error
	return findings, err
}

// GetByCompanyID retrieves findings by company ID with pagination and filters
func (r *ConfigFindingRepository) GetByCompanyID(companyID uuid.UUID, page, limit int, filters map[string]interface{}) ([]models.ConfigFinding, int64, error) {
	var findings []models.ConfigFinding
	var total int64

	query := r.db.Model(&models.ConfigFinding{}).Where("company_id = ?", companyID)

	// Apply filters
	if configFileID, ok := filters["config_file_id"].(*uuid.UUID); ok && configFileID != nil {
		query = query.Where("config_file_id = ?", *configFileID)
	}
	if severity, ok := filters["severity"].(string); ok && severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if category, ok := filters["category"].(string); ok && category != "" {
		query = query.Where("category = ?", category)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if findingType, ok := filters["finding_type"].(string); ok && findingType != "" {
		query = query.Where("finding_type = ?", findingType)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	sortBy := "severity"
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
		Find(&findings).Error

	return findings, total, err
}

// GetByID retrieves a finding by ID
func (r *ConfigFindingRepository) GetByID(id uuid.UUID) (*models.ConfigFinding, error) {
	var finding models.ConfigFinding
	err := r.db.Where("id = ?", id).First(&finding).Error
	if err != nil {
		return nil, err
	}
	return &finding, nil
}

// UpdateStatus updates finding status
func (r *ConfigFindingRepository) UpdateStatus(id uuid.UUID, status string, resolvedBy *uuid.UUID) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if status == "resolved" {
		now := time.Now()
		updates["resolved_at"] = now
		if resolvedBy != nil {
			updates["resolved_by"] = resolvedBy
		}
	}
	return r.db.Model(&models.ConfigFinding{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetStatsByConfigFile retrieves finding statistics for a config file
func (r *ConfigFindingRepository) GetStatsByConfigFile(configFileID uuid.UUID) (map[string]int, error) {
	var stats struct {
		Total     int64 `json:"total"`
		Critical  int64 `json:"critical"`
		High      int64 `json:"high"`
		Medium    int64 `json:"medium"`
		Low       int64 `json:"low"`
		Info      int64 `json:"info"`
		Open      int64 `json:"open"`
		Resolved  int64 `json:"resolved"`
	}

	// Get total
	err := r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ?", configFileID).Count(&stats.Total).Error
	if err != nil {
		return nil, err
	}

	// Get by severity
	err = r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ? AND severity = ?", configFileID, "critical").Count(&stats.Critical).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ? AND severity = ?", configFileID, "high").Count(&stats.High).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ? AND severity = ?", configFileID, "medium").Count(&stats.Medium).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ? AND severity = ?", configFileID, "low").Count(&stats.Low).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ? AND severity = ?", configFileID, "info").Count(&stats.Info).Error
	if err != nil {
		return nil, err
	}

	// Get by status
	err = r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ? AND status = ?", configFileID, "open").Count(&stats.Open).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Model(&models.ConfigFinding{}).Where("config_file_id = ? AND status = ?", configFileID, "resolved").Count(&stats.Resolved).Error
	if err != nil {
		return nil, err
	}

	return map[string]int{
		"total":    int(stats.Total),
		"critical": int(stats.Critical),
		"high":     int(stats.High),
		"medium":   int(stats.Medium),
		"low":      int(stats.Low),
		"info":     int(stats.Info),
		"open":     int(stats.Open),
		"resolved": int(stats.Resolved),
	}, nil
}

// DeleteByConfigFileID deletes all findings for a config file
func (r *ConfigFindingRepository) DeleteByConfigFileID(configFileID uuid.UUID) error {
	return r.db.Where("config_file_id = ?", configFileID).Delete(&models.ConfigFinding{}).Error
}

