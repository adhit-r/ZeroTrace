package services

import (
	"errors"
	"time"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/google/uuid"
)

// ScanService handles scan operations
type ScanService struct {
	config            *config.Config
	scanRepo          *repository.ScanRepository
	enrichmentService *EnrichmentService
}

// NewScanService creates a new scan service
func NewScanService(cfg *config.Config, scanRepo *repository.ScanRepository) *ScanService {
	enrichmentService := NewEnrichmentService(cfg.EnrichmentServiceURL)
	return &ScanService{
		config:            cfg,
		scanRepo:          scanRepo,
		enrichmentService: enrichmentService,
	}
}

// CreateScan creates a new scan with transaction management
func (s *ScanService) CreateScan(req models.CreateScanRequest, companyID uuid.UUID) (*models.Scan, error) {
	scan := &models.Scan{
		ID:         uuid.New(),
		CompanyID:  companyID,
		Repository: req.Repository,
		Branch:     req.Branch,
		ScanType:   req.ScanType,
		Status:     models.ScanStatusPending,
		Progress:   0,
		Options:    req.Options,
		Results:    make(map[string]any),
		Metadata:   make(map[string]any),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save to database using repository (repository handles transactions internally)
	err := s.scanRepo.Create(scan)
	if err != nil {
		return nil, err
	}

	// TODO: Queue scan for processing

	return scan, nil
}

// GetScan retrieves a scan by ID
func (s *ScanService) GetScan(scanID, companyID uuid.UUID) (*models.Scan, error) {
	if scanID == uuid.Nil {
		return nil, errors.New("invalid scan ID")
	}

	// Query from database using repository
	scan, err := s.scanRepo.GetByID(scanID)
	if err != nil {
		return nil, err
	}

	// Verify company ID matches
	if scan.CompanyID != companyID {
		return nil, errors.New("scan not found for this company")
	}

	return scan, nil
}

// GetScans retrieves scans for a company with pagination
func (s *ScanService) GetScans(companyID uuid.UUID, page, limit int) (*models.PaginationResponse, error) {
	// Query from database using repository
	scans, total, err := s.scanRepo.GetByCompanyID(companyID, page, limit)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	scanPointers := make([]*models.Scan, len(scans))
	for i := range scans {
		scanPointers[i] = &scans[i]
	}

	totalPages := (int(total) + limit - 1) / limit
	response := &models.PaginationResponse{
		Data:       scanPointers,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return response, nil
}

// UpdateScan updates a scan with transaction management
func (s *ScanService) UpdateScan(scanID, companyID uuid.UUID, updates map[string]any) (*models.Scan, error) {
	scan, err := s.GetScan(scanID, companyID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	for key, value := range updates {
		switch key {
		case "status":
			if status, ok := value.(string); ok {
				scan.Status = models.ScanStatus(status)
			}
		case "progress":
			if progress, ok := value.(int); ok {
				scan.Progress = progress
			}
		case "notes":
			if notes, ok := value.(string); ok {
				scan.Notes = notes
			}
		}
	}

	scan.UpdatedAt = time.Now()
	
	// Update in database using repository (handles transactions)
	err = s.scanRepo.Update(scan)
	if err != nil {
		return nil, err
	}

	return scan, nil
}

// DeleteScan deletes a scan
func (s *ScanService) DeleteScan(scanID, companyID uuid.UUID) error {
	// TODO: Implement actual database deletion
	return nil
}
