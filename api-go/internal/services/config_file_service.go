package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"path/filepath"
	"regexp"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/constants"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigFileService handles config file operations
type ConfigFileService struct {
	config           *config.Config
	configFileRepo   *repository.ConfigFileRepository
	parserService    *ConfigParserService
	analyzerService  *ConfigAnalyzerService
	jobService       *ConfigJobService
}

// NewConfigFileService creates a new config file service
func NewConfigFileService(
	cfg *config.Config,
	configFileRepo *repository.ConfigFileRepository,
	parserService *ConfigParserService,
	analyzerService *ConfigAnalyzerService,
	jobService *ConfigJobService,
) *ConfigFileService {
	return &ConfigFileService{
		config:          cfg,
		configFileRepo:  configFileRepo,
		parserService:   parserService,
		analyzerService: analyzerService,
		jobService:      jobService,
	}
}

// UploadConfigFile uploads and stores a configuration file
func (s *ConfigFileService) UploadConfigFile(
	fileContent []byte,
	filename string,
	req models.UploadConfigFileRequest,
	companyID uuid.UUID,
	uploadedBy *uuid.UUID,
) (*models.ConfigFile, error) {
	// Validate file content is not empty
	if len(fileContent) < constants.MinConfigFileSize {
		return nil, errors.New("file content cannot be empty")
	}

	// Validate file size (use config value)
	maxFileSize := s.config.ConfigAuditorMaxFileSize
	if maxFileSize == 0 {
		maxFileSize = constants.MaxConfigFileSize // Fallback to constant
	}
	if len(fileContent) > maxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxFileSize)
	}

	// Calculate SHA-256 hash
	hash := sha256.Sum256(fileContent)
	hashString := hex.EncodeToString(hash[:])

	// Check for duplicate (properly handle database errors)
	existing, err := s.configFileRepo.GetByHash(hashString, companyID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check for duplicates: %w", err)
	}
	if existing != nil {
		return nil, errors.New("duplicate file: this configuration has already been uploaded")
	}

	// Detect MIME type
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = constants.MIMETypeTextPlain // Default for config files
	}

	// Validate file content (magic bytes/content-based validation)
	if err := s.validateFileContent(fileContent, filename); err != nil {
		return nil, fmt.Errorf("file content validation failed: %w", err)
	}

	// Detect config format
	configFormat := s.detectConfigFormat(fileContent, filename)

	// Validate required fields
	if req.DeviceType == "" {
		return nil, errors.New("device_type is required")
	}
	if req.Manufacturer == "" {
		return nil, errors.New("manufacturer is required")
	}
	if req.ConfigType == "" {
		return nil, errors.New("config_type is required")
	}

	// Validate enum values
	if !s.isValidDeviceType(req.DeviceType) {
		return nil, fmt.Errorf("invalid device_type: %s", req.DeviceType)
	}
	if !s.isValidConfigType(req.ConfigType) {
		return nil, fmt.Errorf("invalid config_type: %s", req.ConfigType)
	}

	// Validate and sanitize file path components
	if !s.isValidUUID(companyID.String()) {
		return nil, errors.New("invalid company_id format")
	}
	if !s.isValidHash(hashString) {
		return nil, errors.New("invalid file hash format")
	}

	// Create config file record with sanitized path
	storagePath := s.config.ConfigAuditorStoragePath
	if storagePath == "" {
		storagePath = constants.ConfigStoragePathTemplate // Fallback to constant
	}
	filePath := filepath.Join(storagePath, companyID.String(), hashString)
	configFile := &models.ConfigFile{
		CompanyID:       companyID,
		UploadedBy:      uploadedBy,
		Filename:        filename,
		FilePath:        filePath,
		FileSize:        int64(len(fileContent)),
		FileHash:        hashString,
		MimeType:        mimeType,
		FileContent:     fileContent, // Store in PostgreSQL BYTEA
		DeviceType:      req.DeviceType,
		Manufacturer:    req.Manufacturer,
		Model:           req.Model,
		FirmwareVersion: req.FirmwareVersion,
		DeviceName:      req.DeviceName,
		DeviceLocation:  req.DeviceLocation,
		ConfigType:      req.ConfigType,
		ConfigFormat:    configFormat,
		ParsingStatus:   constants.StatusPending,
		AnalysisStatus:  constants.StatusPending,
	}

	// Convert tags to JSON with proper error handling
	if len(req.Tags) > 0 {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		configFile.Tags = tagsJSON
	}

	if req.Notes != "" {
		configFile.Notes = req.Notes
	}

	// Save to database
	err = s.configFileRepo.Create(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to save config file: %w", err)
	}

	// Trigger parsing and analysis asynchronously with error logging
	// Note: This is a fire-and-forget operation. The job service manages workers
	// and will be stopped gracefully on application shutdown.
	go func() {
		if err := s.jobService.QueueConfigAnalysis(configFile.ID); err != nil {
			log.Printf("Failed to queue config analysis for %s: %v", configFile.ID, err)
		}
	}()

	return configFile, nil
}

// GetConfigFile retrieves a config file by ID
func (s *ConfigFileService) GetConfigFile(id uuid.UUID, companyID uuid.UUID) (*models.ConfigFile, error) {
	configFile, err := s.configFileRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Verify company ID matches
	if configFile.CompanyID != companyID {
		return nil, errors.New("config file not found for this company")
	}

	return configFile, nil
}

// ListConfigFiles lists config files with filters and pagination
func (s *ConfigFileService) ListConfigFiles(companyID uuid.UUID, req models.ListConfigFilesRequest) (*models.PaginationResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < constants.MinPageSize {
		pageSize = s.config.ConfigAuditorDefaultPageSize
		if pageSize == 0 {
			pageSize = constants.DefaultPageSize // Fallback to constant
		}
	}
	maxPageSize := s.config.ConfigAuditorMaxPageSize
	if maxPageSize == 0 {
		maxPageSize = constants.MaxPageSize // Fallback to constant
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	filters := map[string]interface{}{
		"manufacturer": req.Manufacturer,
		"device_type":  req.DeviceType,
		"status":       req.Status,
		"sort_by":      req.SortBy,
		"sort_order":   req.SortOrder,
	}

	configFiles, total, err := s.configFileRepo.GetByCompanyID(companyID, page, pageSize, filters)
	if err != nil {
		return nil, err
	}

	// Convert to pointer slice
	configFilePointers := make([]*models.ConfigFile, len(configFiles))
	for i := range configFiles {
		// Don't return file content in list
		configFiles[i].FileContent = nil
		configFilePointers[i] = &configFiles[i]
	}

	totalPages := (int(total) + pageSize - 1) / pageSize
	response := &models.PaginationResponse{
		Data:       configFilePointers,
		Total:      total,
		Page:       page,
		Limit:      pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	return response, nil
}

// GetConfigFileContent retrieves the file content
func (s *ConfigFileService) GetConfigFileContent(id uuid.UUID, companyID uuid.UUID) ([]byte, error) {
	configFile, err := s.GetConfigFile(id, companyID)
	if err != nil {
		return nil, err
	}

	return configFile.FileContent, nil
}

// DeleteConfigFile deletes a config file
func (s *ConfigFileService) DeleteConfigFile(id uuid.UUID, companyID uuid.UUID) error {
	// Verify ownership
	configFile, err := s.GetConfigFile(id, companyID)
	if err != nil {
		return err
	}

	// Delete from database (cascade will delete findings and analysis results)
	return s.configFileRepo.Delete(configFile.ID)
}

// detectConfigFormat detects the configuration file format
func (s *ConfigFileService) detectConfigFormat(content []byte, filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".xml":
		return "xml"
	case ".json":
		return "json"
	case ".txt", ".cfg", ".conf":
		return "text"
	default:
		// Try to detect from content (with bounds checking)
		if len(content) == 0 {
			return "text"
		}
		contentStr := string(content[:min(100, len(content))])
		if len(contentStr) > 0 {
			if contentStr[0] == '<' {
				return "xml"
			}
			if contentStr[0] == '{' || contentStr[0] == '[' {
				return "json"
			}
		}
		return "text"
	}
}

// isValidDeviceType validates device type enum value
func (s *ConfigFileService) isValidDeviceType(deviceType string) bool {
	for _, valid := range constants.ValidDeviceTypes {
		if deviceType == valid {
			return true
		}
	}
	return false
}

// isValidConfigType validates config type enum value
func (s *ConfigFileService) isValidConfigType(configType string) bool {
	for _, valid := range constants.ValidConfigTypes {
		if configType == valid {
			return true
		}
	}
	return false
}

// isValidUUID validates UUID format
func (s *ConfigFileService) isValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

// isValidHash validates SHA-256 hash format (64 hex characters)
func (s *ConfigFileService) isValidHash(hash string) bool {
	if len(hash) != 64 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-fA-F0-9]{64}$", hash)
	return matched
}

// validateFileContent performs content-based validation (magic bytes)
func (s *ConfigFileService) validateFileContent(content []byte, filename string) error {
	if len(content) == 0 {
		return errors.New("file content is empty")
	}

	// Check for XML format
	if len(content) >= 5 && string(content[:5]) == "<?xml" {
		return nil // Valid XML
	}

	// Check for JSON format
	if len(content) > 0 {
		firstChar := content[0]
		if firstChar == '{' || firstChar == '[' {
			// Try to parse as JSON to validate
			var test interface{}
			if err := json.Unmarshal(content, &test); err == nil {
				return nil // Valid JSON
			}
		}
	}

	// For text-based configs, check for common config file signatures
	// Allow text files (most configs are text-based)
	// Additional validation can be added here for specific manufacturers

	return nil // Default: allow text-based configs
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TriggerAnalysis manually triggers analysis for a config file
func (s *ConfigFileService) TriggerAnalysis(id uuid.UUID, companyID uuid.UUID) error {
	configFile, err := s.GetConfigFile(id, companyID)
	if err != nil {
		return err
	}

	return s.jobService.QueueConfigAnalysis(configFile.ID)
}

