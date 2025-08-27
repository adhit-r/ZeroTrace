package repository

import (
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScanRepository handles scan database operations
type ScanRepository struct {
	db *gorm.DB
}

// NewScanRepository creates a new scan repository
func NewScanRepository(db *gorm.DB) *ScanRepository {
	return &ScanRepository{db: db}
}

// Create creates a new scan
func (r *ScanRepository) Create(scan *models.Scan) error {
	scan.ID = uuid.New()
	scan.CreatedAt = time.Now()
	scan.UpdatedAt = time.Now()
	return r.db.Create(scan).Error
}

// GetByID retrieves a scan by ID
func (r *ScanRepository) GetByID(id uuid.UUID) (*models.Scan, error) {
	var scan models.Scan
	err := r.db.Where("id = ?", id).First(&scan).Error
	if err != nil {
		return nil, err
	}
	return &scan, nil
}

// GetByCompanyID retrieves scans by company ID with pagination
func (r *ScanRepository) GetByCompanyID(companyID uuid.UUID, page, limit int) ([]models.Scan, int64, error) {
	var scans []models.Scan
	var total int64

	// Get total count
	err := r.db.Model(&models.Scan{}).Where("company_id = ?", companyID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err = r.db.Where("company_id = ?", companyID).
		Preload("Vulnerabilities").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&scans).Error

	return scans, total, err
}

// GetByStatus retrieves scans by status
func (r *ScanRepository) GetByStatus(companyID uuid.UUID, status models.ScanStatus) ([]models.Scan, error) {
	var scans []models.Scan
	err := r.db.Where("company_id = ? AND status = ?", companyID, status).
		Order("created_at DESC").
		Find(&scans).Error
	return scans, err
}

// Update updates a scan
func (r *ScanRepository) Update(scan *models.Scan) error {
	scan.UpdatedAt = time.Now()
	return r.db.Save(scan).Error
}

// UpdateStatus updates scan status
func (r *ScanRepository) UpdateStatus(id uuid.UUID, status models.ScanStatus) error {
	return r.db.Model(&models.Scan{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// UpdateProgress updates scan progress
func (r *ScanRepository) UpdateProgress(id uuid.UUID, progress int) error {
	return r.db.Model(&models.Scan{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress":   progress,
			"updated_at": time.Now(),
		}).Error
}

// Delete deletes a scan
func (r *ScanRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Scan{}, id).Error
}

// GetStats retrieves scan statistics for a company
func (r *ScanRepository) GetStats(companyID uuid.UUID) (map[string]interface{}, error) {
	var stats struct {
		TotalScans     int64 `json:"total_scans"`
		CompletedScans int64 `json:"completed_scans"`
		ActiveScans    int64 `json:"active_scans"`
		FailedScans    int64 `json:"failed_scans"`
	}

	// Get total scans
	err := r.db.Model(&models.Scan{}).Where("company_id = ?", companyID).Count(&stats.TotalScans).Error
	if err != nil {
		return nil, err
	}

	// Get completed scans
	err = r.db.Model(&models.Scan{}).Where("company_id = ? AND status = ?", companyID, models.ScanStatusCompleted).Count(&stats.CompletedScans).Error
	if err != nil {
		return nil, err
	}

	// Get active scans
	err = r.db.Model(&models.Scan{}).Where("company_id = ? AND status = ?", companyID, models.ScanStatusScanning).Count(&stats.ActiveScans).Error
	if err != nil {
		return nil, err
	}

	// Get failed scans
	err = r.db.Model(&models.Scan{}).Where("company_id = ? AND status = ?", companyID, models.ScanStatusFailed).Count(&stats.FailedScans).Error
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_scans":     stats.TotalScans,
		"completed_scans": stats.CompletedScans,
		"active_scans":    stats.ActiveScans,
		"failed_scans":    stats.FailedScans,
	}, nil
}
