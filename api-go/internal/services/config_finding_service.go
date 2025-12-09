package services

import (
	"errors"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/google/uuid"
)

// ConfigFindingService handles config finding operations
type ConfigFindingService struct {
	configFindingRepo *repository.ConfigFindingRepository
}

// NewConfigFindingService creates a new config finding service
func NewConfigFindingService(configFindingRepo *repository.ConfigFindingRepository) *ConfigFindingService {
	return &ConfigFindingService{
		configFindingRepo: configFindingRepo,
	}
}

// ListConfigFindings lists config findings with filters and pagination
func (s *ConfigFindingService) ListConfigFindings(companyID uuid.UUID, req models.ListConfigFindingsRequest) (*models.PaginationResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	filters := map[string]interface{}{
		"config_file_id": req.ConfigFileID,
		"severity":       req.Severity,
		"category":       req.Category,
		"status":         req.Status,
		"finding_type":   req.FindingType,
		"sort_by":        req.SortBy,
		"sort_order":     req.SortOrder,
	}

	findings, total, err := s.configFindingRepo.GetByCompanyID(companyID, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	findingPointers := make([]*models.ConfigFinding, len(findings))
	for i := range findings {
		findingPointers[i] = &findings[i]
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	response := &models.PaginationResponse{
		Data:       findingPointers,
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return response, nil
}

// GetConfigFinding retrieves a finding by ID
func (s *ConfigFindingService) GetConfigFinding(id uuid.UUID, companyID uuid.UUID) (*models.ConfigFinding, error) {
	finding, err := s.configFindingRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify company ID matches
	if finding.CompanyID != companyID {
		return nil, errors.New("finding not found for this company")
	}

	return finding, nil
}

// UpdateFindingStatus updates finding status
func (s *ConfigFindingService) UpdateFindingStatus(id uuid.UUID, companyID uuid.UUID, status string, assignedTo *uuid.UUID, resolvedBy *uuid.UUID) error {
	// Verify ownership
	finding, err := s.GetConfigFinding(id, companyID)
	if err != nil {
		return err
	}

	// Update status
	return s.configFindingRepo.UpdateStatus(finding.ID, status, resolvedBy)
}

// GetFindingStats retrieves finding statistics for a config file
func (s *ConfigFindingService) GetFindingStats(configFileID uuid.UUID, companyID uuid.UUID) (map[string]int, error) {
	// Verify config file belongs to company (would need config file repo for this)
	// For now, just get stats
	stats, err := s.configFindingRepo.GetStatsByConfigFile(configFileID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

