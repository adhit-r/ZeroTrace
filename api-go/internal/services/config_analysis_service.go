package services

import (
	"encoding/json"
	"errors"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/google/uuid"
)

// ConfigAnalysisService handles config analysis result operations
type ConfigAnalysisService struct {
	configAnalysisRepo *repository.ConfigAnalysisRepository
	configFileRepo     *repository.ConfigFileRepository
}

// NewConfigAnalysisService creates a new config analysis service
func NewConfigAnalysisService(
	configAnalysisRepo *repository.ConfigAnalysisRepository,
	configFileRepo *repository.ConfigFileRepository,
) *ConfigAnalysisService {
	return &ConfigAnalysisService{
		configAnalysisRepo: configAnalysisRepo,
		configFileRepo:     configFileRepo,
	}
}

// GetAnalysisResults retrieves analysis results for a config file
func (s *ConfigAnalysisService) GetAnalysisResults(configFileID uuid.UUID, companyID uuid.UUID) (*models.ConfigAnalysisResult, error) {
	// Verify config file belongs to company
	configFile, err := s.configFileRepo.GetByID(configFileID)
	if err != nil {
		return nil, err
	}

	if configFile.CompanyID != companyID {
		return nil, errors.New("config file not found for this company")
	}

	// Get analysis result
	result, err := s.configAnalysisRepo.GetByConfigFileID(configFileID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetComplianceScores retrieves compliance scores for a config file
func (s *ConfigAnalysisService) GetComplianceScores(configFileID uuid.UUID, companyID uuid.UUID) (map[string]float64, error) {
	result, err := s.GetAnalysisResults(configFileID, companyID)
	if err != nil {
		return nil, err
	}

	// Parse compliance scores from JSONB
	var scores map[string]float64
	if result.ComplianceScores != nil {
		err = json.Unmarshal(result.ComplianceScores, &scores)
		if err != nil {
			return nil, err
		}
	}

	return scores, nil
}

// GetAnalysisStatus retrieves analysis status for a config file
func (s *ConfigAnalysisService) GetAnalysisStatus(configFileID uuid.UUID, companyID uuid.UUID) (string, error) {
	configFile, err := s.configFileRepo.GetByID(configFileID)
	if err != nil {
		return "", err
	}

	if configFile.CompanyID != companyID {
		return "", errors.New("config file not found for this company")
	}

	return configFile.AnalysisStatus, nil
}

