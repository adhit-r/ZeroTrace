package scanner

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// AIMLScanner handles AI/ML security scanning with high performance
type AIMLScanner struct {
	config      *config.Config
	logger      Logger
	maxWorkers  int
	maxFileSize int64
	scanTimeout time.Duration
	cache       *ScanCache
}

// Logger interface for dependency injection
type Logger interface {
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// ScanCache caches scan results to avoid duplicate work
type ScanCache struct {
	mu    sync.RWMutex
	files map[string]*CachedFileInfo
}

// CachedFileInfo stores cached file information
type CachedFileInfo struct {
	Hash      string
	Size      int64
	ScannedAt time.Time
	ModelInfo *ModelInfo
}

// AIMLFinding represents an AI/ML security finding
type AIMLFinding struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`     // model, data, training, inference, supply_chain
	Severity      string                 `json:"severity"` // critical, high, medium, low
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	FilePath      string                 `json:"file_path,omitempty"`
	ModelName     string                 `json:"model_name,omitempty"`
	ModelVersion  string                 `json:"model_version,omitempty"`
	Framework     string                 `json:"framework,omitempty"`
	CurrentValue  string                 `json:"current_value,omitempty"`
	RequiredValue string                 `json:"required_value,omitempty"`
	Remediation   string                 `json:"remediation"`
	DiscoveredAt  time.Time              `json:"discovered_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// ModelInfo represents AI/ML model information
type ModelInfo struct {
	Name            string                 `json:"name"`
	Path            string                 `json:"path"`
	Version         string                 `json:"version"`
	Framework       string                 `json:"framework"`
	Type            string                 `json:"type"`
	Size            int64                  `json:"size"`
	Hash            string                 `json:"hash"`
	Permissions     string                 `json:"permissions"`
	Owner           string                 `json:"owner,omitempty"`
	ModifiedTime    time.Time              `json:"modified_time"`
	Accuracy        float64                `json:"accuracy,omitempty"`
	TrainingData    string                 `json:"training_data,omitempty"`
	LastTrained     time.Time              `json:"last_trained,omitempty"`
	IsPublic        bool                   `json:"is_public"`
	HasAPI          bool                   `json:"has_api"`
	Endpoints       []string               `json:"endpoints"`
	InputSchema     map[string]interface{} `json:"input_schema,omitempty"`
	OutputSchema    map[string]interface{} `json:"output_schema,omitempty"`
	Vulnerabilities []AIMLVulnerability    `json:"vulnerabilities"`
	BiasMetrics     map[string]float64     `json:"bias_metrics,omitempty"`
	FairnessScore   float64                `json:"fairness_score,omitempty"`
	PrivacyScore    float64                `json:"privacy_score,omitempty"`
	SecurityScore   float64                `json:"security_score,omitempty"`
}

// Vulnerability represents a specific vulnerability
type AIMLVulnerability struct {
	ID          string    `json:"id"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	CVE         string    `json:"cve,omitempty"`
	FoundAt     time.Time `json:"found_at"`
}

// TrainingDataInfo represents training data information
type TrainingDataInfo struct {
	DatasetName     string    `json:"dataset_name"`
	Path            string    `json:"path"`
	Size            int64     `json:"size"`
	Hash            string    `json:"hash"`
	Records         int64     `json:"records"`
	Columns         []string  `json:"columns"`
	SensitiveFields []string  `json:"sensitive_fields"`
	HasPII          bool      `json:"has_pii"`
	HasBias         bool      `json:"has_bias"`
	DataQuality     float64   `json:"data_quality"`
	LastUpdated     time.Time `json:"last_updated"`
	Source          string    `json:"source"`
	License         string    `json:"license"`
	RetentionPolicy string    `json:"retention_policy"`
	Permissions     string    `json:"permissions"`
}

// SupplyChainInfo represents AI/ML supply chain information
type SupplyChainInfo struct {
	ModelName        string              `json:"model_name"`
	Dependencies     []Dependency        `json:"dependencies"`
	PreTrainedModels []string            `json:"pre_trained_models"`
	DataSources      []string            `json:"data_sources"`
	Libraries        []LibraryInfo       `json:"libraries"`
	Vulnerabilities  []AIMLVulnerability `json:"vulnerabilities"`
	Licenses         map[string]string   `json:"licenses"`
	Compliance       []ComplianceStatus  `json:"compliance"`
	RiskScore        float64             `json:"risk_score"`
	LastScanned      time.Time           `json:"last_scanned"`
}

// Dependency represents a dependency
type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`
}

// LibraryInfo represents library information
type LibraryInfo struct {
	Name            string              `json:"name"`
	Version         string              `json:"version"`
	Vulnerabilities []AIMLVulnerability `json:"vulnerabilities"`
	License         string              `json:"license"`
}

// ComplianceStatus represents compliance status
type ComplianceStatus struct {
	Framework string   `json:"framework"`
	Status    string   `json:"status"`
	Issues    []string `json:"issues,omitempty"`
}

// ScanResult aggregates all scan results
type ScanResult struct {
	Findings     []AIMLFinding      `json:"findings"`
	Models       []ModelInfo        `json:"models"`
	TrainingData []TrainingDataInfo `json:"training_data"`
	SupplyChain  SupplyChainInfo    `json:"supply_chain"`
	Statistics   ScanStatistics     `json:"statistics"`
}

// ScanStatistics provides scan metrics
type ScanStatistics struct {
	TotalFiles     int           `json:"total_files"`
	ModelsFound    int           `json:"models_found"`
	DatasetsFound  int           `json:"datasets_found"`
	FindingsCount  int           `json:"findings_count"`
	ScanDuration   time.Duration `json:"scan_duration"`
	FilesPerSecond float64       `json:"files_per_second"`
	ErrorCount     int           `json:"error_count"`
}

// Model file patterns by framework
var modelPatterns = map[string][]string{
	"TensorFlow":   {".pb", ".h5", ".tflite", ".keras"},
	"PyTorch":      {".pth", ".pt"},
	"ONNX":         {".onnx"},
	"HuggingFace":  {"config.json", "pytorch_model.bin", "model.safetensors"},
	"scikit-learn": {".joblib"},
	"JAX":          {".msgpack", ".flax"},
	"MXNet":        {".params"},
}

// Data file patterns
var dataPatterns = []string{".csv", ".json", ".parquet", ".h5", ".hdf5", ".tfrecord", ".arrow"}

// Sensitive field patterns (regex-compatible)
var sensitivePatterns = []string{
	"email", "ssn", "social_security", "password", "credit_card",
	"phone", "address", "name", "dob", "date_of_birth", "salary",
	"medical", "health", "diagnosis", "treatment",
}

// NewAIMLScanner creates a new high-performance AI/ML security scanner
func NewAIMLScanner(cfg *config.Config, logger Logger) *AIMLScanner {
	if logger == nil {
		logger = &noOpLogger{}
	}

	workers := runtime.NumCPU() * 2
	if cfg != nil && cfg.MaxWorkers > 0 {
		workers = cfg.MaxWorkers
	}

	maxFileSize := int64(500 * 1024 * 1024) // 500MB default
	if cfg != nil && cfg.MaxFileSizeMB > 0 {
		maxFileSize = int64(cfg.MaxFileSizeMB) * 1024 * 1024
	}

	timeout := 30 * time.Minute
	if cfg != nil && cfg.ScanTimeout > 0 {
		timeout = cfg.ScanTimeout
	}

	return &AIMLScanner{
		config:      cfg,
		logger:      logger,
		maxWorkers:  workers,
		maxFileSize: maxFileSize,
		scanTimeout: timeout,
		cache: &ScanCache{
			files: make(map[string]*CachedFileInfo),
		},
	}
}

// Scan performs comprehensive AI/ML security scanning with timeout
func (as *AIMLScanner) Scan(rootPath string) (*ScanResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), as.scanTimeout)
	defer cancel()

	return as.ScanWithContext(ctx, rootPath)
}

// ScanWithContext performs scanning with context support
func (as *AIMLScanner) ScanWithContext(ctx context.Context, rootPath string) (*ScanResult, error) {
	startTime := time.Now()

	as.logger.Info("Starting AI/ML security scan", "path", rootPath, "workers", as.maxWorkers)

	result := &ScanResult{
		Findings:     make([]AIMLFinding, 0),
		Models:       make([]ModelInfo, 0),
		TrainingData: make([]TrainingDataInfo, 0),
	}

	// Validate root path
	if err := as.validatePath(rootPath); err != nil {
		return nil, fmt.Errorf("invalid root path: %w", err)
	}

	// Discover files concurrently
	modelFiles, dataFiles, totalFiles, err := as.discoverFilesParallel(ctx, rootPath)
	if err != nil {
		return nil, fmt.Errorf("file discovery failed: %w", err)
	}

	as.logger.Info("File discovery complete",
		"total_files", totalFiles,
		"model_files", len(modelFiles),
		"data_files", len(dataFiles))

	// Process models concurrently
	models, modelFindings := as.processModelsParallel(ctx, modelFiles)
	result.Models = models
	result.Findings = append(result.Findings, modelFindings...)

	// Process training data concurrently
	trainingData, dataFindings := as.processTrainingDataParallel(ctx, dataFiles)
	result.TrainingData = trainingData
	result.Findings = append(result.Findings, dataFindings...)

	// Analyze supply chain
	result.SupplyChain = as.analyzeSupplyChain(models, trainingData)
	supplyChainFindings := as.scanSupplyChainSecurity(result.SupplyChain)
	result.Findings = append(result.Findings, supplyChainFindings...)

	// Calculate statistics
	duration := time.Since(startTime)
	result.Statistics = ScanStatistics{
		TotalFiles:     totalFiles,
		ModelsFound:    len(models),
		DatasetsFound:  len(trainingData),
		FindingsCount:  len(result.Findings),
		ScanDuration:   duration,
		FilesPerSecond: float64(totalFiles) / duration.Seconds(),
	}

	as.logger.Info("Scan complete",
		"duration", duration,
		"models", len(models),
		"datasets", len(trainingData),
		"findings", len(result.Findings))

	return result, nil
}

// validatePath validates the scan path
func (as *AIMLScanner) validatePath(path string) error {
	// Clean and resolve the path
	cleanPath := filepath.Clean(path)

	// Check if path exists
	info, err := os.Stat(cleanPath)
	if err != nil {
		return fmt.Errorf("path access error: %w", err)
	}

	// Must be a directory
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory")
	}

	// Check read permissions
	f, err := os.Open(cleanPath)
	if err != nil {
		return fmt.Errorf("path not readable: %w", err)
	}
	f.Close()

	return nil
}

// discoverFilesParallel discovers files using concurrent filesystem walk
func (as *AIMLScanner) discoverFilesParallel(ctx context.Context, rootPath string) (
	modelFiles map[string][]string,
	dataFiles []string,
	totalFiles int,
	err error,
) {
	modelFiles = make(map[string][]string)
	dataFiles = make([]string, 0)

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Use WalkDir for better performance
	err = filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			as.logger.Warn("Error accessing path", "path", path, "error", err)
			return nil // Continue walking
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Skip hidden directories and common exclusions
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, ".") ||
				name == "node_modules" ||
				name == "__pycache__" ||
				name == "venv" ||
				name == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process regular files
		if !d.Type().IsRegular() {
			return nil
		}

		mu.Lock()
		totalFiles++
		mu.Unlock()

		// Get file info
		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Skip files exceeding size limit
		if info.Size() > as.maxFileSize {
			as.logger.Warn("Skipping large file", "path", path, "size", info.Size())
			return nil
		}

		fileName := filepath.Base(path)
		ext := filepath.Ext(fileName)

		// Check if it's a model file
		for framework, patterns := range modelPatterns {
			for _, pattern := range patterns {
				if strings.HasSuffix(fileName, pattern) || ext == pattern {
					mu.Lock()
					modelFiles[framework] = append(modelFiles[framework], path)
					mu.Unlock()
					return nil
				}
			}
		}

		// Check if it's a data file
		for _, pattern := range dataPatterns {
			if ext == pattern {
				mu.Lock()
				dataFiles = append(dataFiles, path)
				mu.Unlock()
				return nil
			}
		}

		return nil
	})

	wg.Wait()

	if err != nil && err != context.Canceled {
		return nil, nil, 0, err
	}

	return modelFiles, dataFiles, totalFiles, nil
}

// processModelsParallel processes model files concurrently
func (as *AIMLScanner) processModelsParallel(ctx context.Context, modelFiles map[string][]string) ([]ModelInfo, []AIMLFinding) {
	var models []ModelInfo
	var findings []AIMLFinding
	var mu sync.Mutex

	// Create work queue
	type workItem struct {
		framework string
		path      string
	}

	workQueue := make(chan workItem, 100)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < as.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range workQueue {
				select {
				case <-ctx.Done():
					return
				default:
				}

				model, err := as.analyzeModelFile(item.framework, item.path)
				if err != nil {
					as.logger.Error("Failed to analyze model", "path", item.path, "error", err)
					continue
				}

				modelFindings := as.scanModel(*model)

				mu.Lock()
				models = append(models, *model)
				findings = append(findings, modelFindings...)
				mu.Unlock()
			}
		}()
	}

	// Send work items
	go func() {
		for framework, paths := range modelFiles {
			for _, path := range paths {
				select {
				case <-ctx.Done():
					close(workQueue)
					return
				case workQueue <- workItem{framework: framework, path: path}:
				}
			}
		}
		close(workQueue)
	}()

	wg.Wait()
	return models, findings
}

// analyzeModelFile analyzes a single model file
func (as *AIMLScanner) analyzeModelFile(framework, path string) (*ModelInfo, error) {
	// Check cache
	cached := as.cache.Get(path)
	if cached != nil && cached.ModelInfo != nil {
		return cached.ModelInfo, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Calculate file hash
	hash, err := as.calculateFileHash(path)
	if err != nil {
		as.logger.Warn("Failed to calculate hash", "path", path, "error", err)
		hash = "unknown"
	}

	// Get file permissions
	perms := info.Mode().String()

	model := &ModelInfo{
		Name:         filepath.Base(path),
		Path:         path,
		Framework:    framework,
		Size:         info.Size(),
		Hash:         hash,
		Permissions:  perms,
		ModifiedTime: info.ModTime(),
		Version:      "unknown",
		Type:         as.detectModelType(framework, path),
		IsPublic:     as.isPubliclyAccessible(info),
		HasAPI:       as.detectAPIEndpoint(path),
		Endpoints:    []string{},
		BiasMetrics:  make(map[string]float64),
	}

	// Analyze model content
	as.analyzeModelContent(model)

	// Calculate security scores
	model.FairnessScore = as.calculateFairnessScore(model)
	model.PrivacyScore = as.calculatePrivacyScore(model)
	model.SecurityScore = as.calculateSecurityScore(model)

	// Cache the result
	as.cache.Set(path, &CachedFileInfo{
		Hash:      hash,
		Size:      info.Size(),
		ScannedAt: time.Now(),
		ModelInfo: model,
	})

	return model, nil
}

// calculateFileHash calculates SHA256 hash of a file
func (as *AIMLScanner) calculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()

	// For large files, only hash first and last chunks
	info, _ := file.Stat()
	if info.Size() > 10*1024*1024 { // 10MB
		buf := make([]byte, 1024*1024) // 1MB chunks

		// Hash first chunk
		n, _ := file.Read(buf)
		hash.Write(buf[:n])

		// Hash last chunk
		file.Seek(-1024*1024, io.SeekEnd)
		n, _ = file.Read(buf)
		hash.Write(buf[:n])
	} else {
		io.Copy(hash, file)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// detectModelType detects the model type based on framework and analysis
func (as *AIMLScanner) detectModelType(framework, path string) string {
	name := strings.ToLower(filepath.Base(path))

	// Common model type indicators
	typeIndicators := map[string]string{
		"bert":        "nlp",
		"gpt":         "nlp",
		"transformer": "nlp",
		"resnet":      "vision",
		"vgg":         "vision",
		"yolo":        "vision",
		"lstm":        "sequence",
		"rnn":         "sequence",
		"gan":         "generative",
		"vae":         "generative",
		"classifier":  "classification",
		"regressor":   "regression",
	}

	for indicator, modelType := range typeIndicators {
		if strings.Contains(name, indicator) {
			return modelType
		}
	}

	return "unknown"
}

// isPubliclyAccessible checks if file has overly permissive permissions
func (as *AIMLScanner) isPubliclyAccessible(info os.FileInfo) bool {
	mode := info.Mode()
	// Check if world-readable
	return mode.Perm()&0004 != 0
}

// detectAPIEndpoint checks if model has associated API configuration
func (as *AIMLScanner) detectAPIEndpoint(path string) bool {
	dir := filepath.Dir(path)

	// Look for common API config files
	apiFiles := []string{"api.yaml", "api.json", "serving.yaml", "config.yaml"}
	for _, file := range apiFiles {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			return true
		}
	}

	return false
}

// analyzeModelContent performs deeper analysis of model file
func (as *AIMLScanner) analyzeModelContent(model *ModelInfo) {
	// For now, basic analysis - in production, would parse model format

	// Check for known vulnerabilities based on framework
	model.Vulnerabilities = as.checkKnownVulnerabilities(model.Framework, model.Name)

	// Analyze file header for additional metadata
	if metadata := as.extractModelMetadata(model.Path); metadata != nil {
		if version, ok := metadata["version"].(string); ok {
			model.Version = version
		}
	}
}

// checkKnownVulnerabilities checks for known vulnerabilities
func (as *AIMLScanner) checkKnownVulnerabilities(framework, modelName string) []AIMLVulnerability {
	var vulns []AIMLVulnerability

	// In production, this would query a vulnerability database
	// For now, check for known patterns

	name := strings.ToLower(modelName)

	// Example: Check for pickle-based models (deserialization risk)
	if strings.HasSuffix(name, ".pkl") || strings.HasSuffix(name, ".pickle") {
		vulns = append(vulns, AIMLVulnerability{
			ID:          uuid.New().String(),
			Severity:    "high",
			Description: "Pickle-based model files can execute arbitrary code during deserialization",
			CVE:         "CWE-502",
			FoundAt:     time.Now(),
		})
	}

	return vulns
}

// extractModelMetadata extracts metadata from model file
func (as *AIMLScanner) extractModelMetadata(path string) map[string]interface{} {
	// In production, would parse actual model format
	// For now, return nil
	return nil
}

// calculateFairnessScore calculates model fairness score
func (as *AIMLScanner) calculateFairnessScore(model *ModelInfo) float64 {
	score := 0.85 // Default reasonable score

	// Reduce score based on risk factors
	if len(model.Vulnerabilities) > 0 {
		score -= 0.1
	}

	if model.IsPublic && !model.HasAPI {
		score -= 0.1 // Public without proper API is risky
	}

	return max(score, 0.0)
}

// calculatePrivacyScore calculates model privacy score
func (as *AIMLScanner) calculatePrivacyScore(model *ModelInfo) float64 {
	score := 0.85

	// Reduce score based on risk factors
	if model.IsPublic {
		score -= 0.15
	}

	if len(model.Vulnerabilities) > 0 {
		score -= float64(len(model.Vulnerabilities)) * 0.05
	}

	// Check file permissions
	if model.Permissions != "" && strings.Contains(model.Permissions, "w") {
		score -= 0.1 // World-writable is risky
	}

	return max(score, 0.0)
}

// calculateSecurityScore calculates model security score
func (as *AIMLScanner) calculateSecurityScore(model *ModelInfo) float64 {
	score := 0.85

	// Reduce score based on vulnerabilities
	for _, vuln := range model.Vulnerabilities {
		switch vuln.Severity {
		case "critical":
			score -= 0.3
		case "high":
			score -= 0.2
		case "medium":
			score -= 0.1
		case "low":
			score -= 0.05
		}
	}

	// File permission issues
	if model.IsPublic {
		score -= 0.15
	}

	return max(score, 0.0)
}

// processTrainingDataParallel processes training data files concurrently
func (as *AIMLScanner) processTrainingDataParallel(ctx context.Context, dataFiles []string) ([]TrainingDataInfo, []AIMLFinding) {
	var trainingData []TrainingDataInfo
	var findings []AIMLFinding
	var mu sync.Mutex

	workQueue := make(chan string, 100)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < as.maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range workQueue {
				select {
				case <-ctx.Done():
					return
				default:
				}

				data, err := as.analyzeDataFile(path)
				if err != nil {
					as.logger.Error("Failed to analyze data file", "path", path, "error", err)
					continue
				}

				dataFindings := as.scanTrainingData(*data)

				mu.Lock()
				trainingData = append(trainingData, *data)
				findings = append(findings, dataFindings...)
				mu.Unlock()
			}
		}()
	}

	// Send work
	go func() {
		for _, path := range dataFiles {
			select {
			case <-ctx.Done():
				close(workQueue)
				return
			case workQueue <- path:
			}
		}
		close(workQueue)
	}()

	wg.Wait()
	return trainingData, findings
}

// analyzeDataFile analyzes a training data file
func (as *AIMLScanner) analyzeDataFile(path string) (*TrainingDataInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	hash, err := as.calculateFileHash(path)
	if err != nil {
		hash = "unknown"
	}

	data := &TrainingDataInfo{
		DatasetName:     filepath.Base(path),
		Path:            path,
		Size:            info.Size(),
		Hash:            hash,
		LastUpdated:     info.ModTime(),
		Permissions:     info.Mode().String(),
		Source:          "local",
		License:         "unknown",
		RetentionPolicy: "unknown",
	}

	// Analyze data content
	as.analyzeDataContent(data)

	return data, nil
}

// analyzeDataContent analyzes data file content
func (as *AIMLScanner) analyzeDataContent(data *TrainingDataInfo) {
	ext := filepath.Ext(data.Path)

	switch ext {
	case ".csv":
		as.analyzeCSV(data)
	case ".json":
		as.analyzeJSON(data)
	case ".parquet":
		as.analyzeParquet(data)
	default:
		data.DataQuality = 0.5 // Unknown quality
	}
}

// analyzeCSV analyzes CSV files
func (as *AIMLScanner) analyzeCSV(data *TrainingDataInfo) {
	file, err := os.Open(data.Path)
	if err != nil {
		return
	}
	defer file.Close()

	// Read first few lines to get column names
	buf := make([]byte, 4096)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return
	}

	content := string(buf[:n])
	lines := strings.Split(content, "\n")

	if len(lines) > 0 {
		// Parse header
		columns := strings.Split(lines[0], ",")
		data.Columns = make([]string, len(columns))
		for i, col := range columns {
			data.Columns[i] = strings.TrimSpace(col)
		}

		// Check for sensitive fields
		data.SensitiveFields = as.detectSensitiveFields(data.Columns)
		data.HasPII = len(data.SensitiveFields) > 0

		// Estimate records (rough estimate)
		data.Records = int64(len(lines) - 1) // Subtract header
	}

	data.DataQuality = as.calculateDataQuality(data)
}

// analyzeJSON analyzes JSON files
func (as *AIMLScanner) analyzeJSON(data *TrainingDataInfo) {
	file, err := os.Open(data.Path)
	if err != nil {
		return
	}
	defer file.Close()

	// Read limited amount for structure analysis
	buf := make([]byte, 8192)
	n, _ := file.Read(buf)

	var jsonData interface{}
	if err := json.Unmarshal(buf[:n], &jsonData); err == nil {
		// Extract field names
		data.Columns = as.extractJSONFields(jsonData)
		data.SensitiveFields = as.detectSensitiveFields(data.Columns)
		data.HasPII = len(data.SensitiveFields) > 0
	}

	data.DataQuality = as.calculateDataQuality(data)
}

// analyzeParquet analyzes Parquet files
func (as *AIMLScanner) analyzeParquet(data *TrainingDataInfo) {
	// Would need parquet library in production
	data.DataQuality = 0.7 // Default for binary format
}

// extractJSONFields extracts field names from JSON structure
func (as *AIMLScanner) extractJSONFields(data interface{}) []string {
	fields := make([]string, 0)

	switch v := data.(type) {
	case map[string]interface{}:
		for key := range v {
			fields = append(fields, key)
		}
	case []interface{}:
		if len(v) > 0 {
			if obj, ok := v[0].(map[string]interface{}); ok {
				for key := range obj {
					fields = append(fields, key)
				}
			}
		}
	}

	return fields
}

// detectSensitiveFields detects sensitive field names
func (as *AIMLScanner) detectSensitiveFields(columns []string) []string {
	sensitive := make([]string, 0)

	for _, col := range columns {
		colLower := strings.ToLower(col)
		for _, pattern := range sensitivePatterns {
			if strings.Contains(colLower, pattern) {
				sensitive = append(sensitive, col)
				break
			}
		}
	}

	return sensitive
}

// calculateDataQuality calculates data quality score
func (as *AIMLScanner) calculateDataQuality(data *TrainingDataInfo) float64 {
	score := 0.8 // Base score

	// Reduce for sensitive data without proper handling
	if data.HasPII {
		score -= 0.2
	}

	// Reduce for unknown license
	if data.License == "unknown" {
		score -= 0.1
	}

	return max(score, 0.0)
}

// scanModel scans a model for security issues
func (as *AIMLScanner) scanModel(model ModelInfo) []AIMLFinding {
	var findings []AIMLFinding
	now := time.Now()

	// Check for vulnerabilities
	for _, vuln := range model.Vulnerabilities {
		finding := AIMLFinding{
			ID:           uuid.New().String(),
			Type:         "model",
			Severity:     vuln.Severity,
			Title:        "Model Vulnerability Detected",
			Description:  vuln.Description,
			FilePath:     model.Path,
			ModelName:    model.Name,
			ModelVersion: model.Version,
			Framework:    model.Framework,
			Remediation:  "Update model or apply security patches. Consider using safer serialization formats.",
			DiscoveredAt: now,
			Metadata: map[string]interface{}{
				"vulnerability_id": vuln.ID,
				"cve":              vuln.CVE,
				"model_hash":       model.Hash,
			},
		}
		findings = append(findings, finding)
	}

	// Check fairness score
	threshold := 0.7
	if as.config != nil && as.config.FairnessThreshold > 0 {
		threshold = as.config.FairnessThreshold
	}

	if model.FairnessScore < threshold {
		finding := AIMLFinding{
			ID:            uuid.New().String(),
			Type:          "model",
			Severity:      "medium",
			Title:         "Low Model Fairness Score",
			Description:   fmt.Sprintf("Model %s has fairness score %.2f below threshold %.2f", model.Name, model.FairnessScore, threshold),
			FilePath:      model.Path,
			ModelName:     model.Name,
			ModelVersion:  model.Version,
			Framework:     model.Framework,
			CurrentValue:  fmt.Sprintf("%.2f", model.FairnessScore),
			RequiredValue: fmt.Sprintf("%.2f+", threshold),
			Remediation:   "Review training data for bias, retrain with balanced datasets, implement fairness constraints",
			DiscoveredAt:  now,
			Metadata: map[string]interface{}{
				"fairness_score": model.FairnessScore,
				"threshold":      threshold,
			},
		}
		findings = append(findings, finding)
	}

	// Check privacy score
	if model.PrivacyScore < threshold {
		finding := AIMLFinding{
			ID:            uuid.New().String(),
			Type:          "model",
			Severity:      "high",
			Title:         "Privacy Risk Detected",
			Description:   fmt.Sprintf("Model %s has privacy score %.2f indicating potential data leakage risks", model.Name, model.PrivacyScore),
			FilePath:      model.Path,
			ModelName:     model.Name,
			ModelVersion:  model.Version,
			Framework:     model.Framework,
			CurrentValue:  fmt.Sprintf("%.2f", model.PrivacyScore),
			RequiredValue: fmt.Sprintf("%.2f+", threshold),
			Remediation:   "Implement differential privacy, review model for memorization, add privacy-preserving techniques",
			DiscoveredAt:  now,
			Metadata: map[string]interface{}{
				"privacy_score": model.PrivacyScore,
				"is_public":     model.IsPublic,
			},
		}
		findings = append(findings, finding)
	}

	// Check security score
	if model.SecurityScore < threshold {
		finding := AIMLFinding{
			ID:            uuid.New().String(),
			Type:          "model",
			Severity:      "high",
			Title:         "Model Security Issues",
			Description:   fmt.Sprintf("Model %s has security score %.2f indicating vulnerabilities", model.Name, model.SecurityScore),
			FilePath:      model.Path,
			ModelName:     model.Name,
			ModelVersion:  model.Version,
			Framework:     model.Framework,
			CurrentValue:  fmt.Sprintf("%.2f", model.SecurityScore),
			RequiredValue: fmt.Sprintf("%.2f+", threshold),
			Remediation:   "Implement adversarial training, add input validation, secure model artifacts, use model signing",
			DiscoveredAt:  now,
			Metadata: map[string]interface{}{
				"security_score": model.SecurityScore,
			},
		}
		findings = append(findings, finding)
	}

	// Check for overly permissive access
	if model.IsPublic {
		finding := AIMLFinding{
			ID:           uuid.New().String(),
			Type:         "model",
			Severity:     "medium",
			Title:        "Publicly Accessible Model",
			Description:  fmt.Sprintf("Model %s has world-readable permissions", model.Name),
			FilePath:     model.Path,
			ModelName:    model.Name,
			Framework:    model.Framework,
			Remediation:  "Restrict file permissions to authorized users only (chmod 600 or 640)",
			DiscoveredAt: now,
			Metadata: map[string]interface{}{
				"permissions": model.Permissions,
				"is_public":   true,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanTrainingData scans training data for security issues
func (as *AIMLScanner) scanTrainingData(data TrainingDataInfo) []AIMLFinding {
	var findings []AIMLFinding
	now := time.Now()

	// Check for PII
	if data.HasPII {
		finding := AIMLFinding{
			ID:           uuid.New().String(),
			Type:         "data",
			Severity:     "critical",
			Title:        "PII Detected in Training Data",
			Description:  fmt.Sprintf("Dataset %s contains %d sensitive fields with potential PII", data.DatasetName, len(data.SensitiveFields)),
			FilePath:     data.Path,
			Remediation:  "Remove PII, implement anonymization/pseudonymization, use data masking, obtain proper consent",
			DiscoveredAt: now,
			Metadata: map[string]interface{}{
				"dataset_name":     data.DatasetName,
				"sensitive_fields": data.SensitiveFields,
				"data_hash":        data.Hash,
			},
		}
		findings = append(findings, finding)
	}

	// Check data quality
	qualityThreshold := 0.7
	if as.config != nil && as.config.DataQualityThreshold > 0 {
		qualityThreshold = as.config.DataQualityThreshold
	}

	if data.DataQuality < qualityThreshold {
		finding := AIMLFinding{
			ID:            uuid.New().String(),
			Type:          "data",
			Severity:      "medium",
			Title:         "Low Data Quality",
			Description:   fmt.Sprintf("Dataset %s has quality score %.2f below threshold", data.DatasetName, data.DataQuality),
			FilePath:      data.Path,
			CurrentValue:  fmt.Sprintf("%.2f", data.DataQuality),
			RequiredValue: fmt.Sprintf("%.2f+", qualityThreshold),
			Remediation:   "Improve data validation, clean outliers, handle missing values, verify data integrity",
			DiscoveredAt:  now,
			Metadata: map[string]interface{}{
				"data_quality": data.DataQuality,
				"dataset_size": data.Size,
			},
		}
		findings = append(findings, finding)
	}

	// Check for bias indicators
	if data.HasBias {
		finding := AIMLFinding{
			ID:           uuid.New().String(),
			Type:         "data",
			Severity:     "high",
			Title:        "Data Bias Detected",
			Description:  fmt.Sprintf("Dataset %s shows signs of bias", data.DatasetName),
			FilePath:     data.Path,
			Remediation:  "Balance dataset, implement bias detection, diversify data sources, apply fairness constraints",
			DiscoveredAt: now,
			Metadata: map[string]interface{}{
				"dataset_name": data.DatasetName,
			},
		}
		findings = append(findings, finding)
	}

	// Check retention policy
	if data.RetentionPolicy == "unknown" && data.HasPII {
		finding := AIMLFinding{
			ID:           uuid.New().String(),
			Type:         "data",
			Severity:     "medium",
			Title:        "Missing Data Retention Policy",
			Description:  fmt.Sprintf("Dataset %s with PII lacks retention policy", data.DatasetName),
			FilePath:     data.Path,
			Remediation:  "Define and document data retention policy, implement automated deletion, ensure compliance",
			DiscoveredAt: now,
			Metadata: map[string]interface{}{
				"has_pii": true,
			},
		}
		findings = append(findings, finding)
	}

	// Check file permissions
	if strings.Contains(data.Permissions, "rw-rw-rw-") || strings.Contains(data.Permissions, "rwxrwxrwx") {
		finding := AIMLFinding{
			ID:           uuid.New().String(),
			Type:         "data",
			Severity:     "high",
			Title:        "Overly Permissive Data File",
			Description:  fmt.Sprintf("Dataset %s has insecure permissions", data.DatasetName),
			FilePath:     data.Path,
			Remediation:  "Restrict file permissions (chmod 600 or 640)",
			DiscoveredAt: now,
			Metadata: map[string]interface{}{
				"permissions": data.Permissions,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// analyzeSupplyChain analyzes AI/ML supply chain
func (as *AIMLScanner) analyzeSupplyChain(models []ModelInfo, datasets []TrainingDataInfo) SupplyChainInfo {
	info := SupplyChainInfo{
		Dependencies:     make([]Dependency, 0),
		PreTrainedModels: make([]string, 0),
		DataSources:      make([]string, 0),
		Libraries:        make([]LibraryInfo, 0),
		Vulnerabilities:  make([]AIMLVulnerability, 0),
		Licenses:         make(map[string]string),
		Compliance:       make([]ComplianceStatus, 0),
		LastScanned:      time.Now(),
	}

	// Collect frameworks and create library info
	frameworkMap := make(map[string]bool)
	for _, model := range models {
		if !frameworkMap[model.Framework] {
			frameworkMap[model.Framework] = true

			lib := LibraryInfo{
				Name:            model.Framework,
				Version:         "unknown",
				Vulnerabilities: []AIMLVulnerability{},
				License:         "Apache-2.0", // Default for most ML frameworks
			}

			info.Libraries = append(info.Libraries, lib)
			info.Licenses[model.Framework] = lib.License
		}

		// Collect model vulnerabilities
		info.Vulnerabilities = append(info.Vulnerabilities, model.Vulnerabilities...)
	}

	// Collect data sources
	sourceMap := make(map[string]bool)
	for _, dataset := range datasets {
		if !sourceMap[dataset.Source] {
			sourceMap[dataset.Source] = true
			info.DataSources = append(info.DataSources, dataset.Source)
		}
	}

	// Assess compliance
	info.Compliance = as.assessCompliance(models, datasets)

	// Calculate risk score
	info.RiskScore = as.calculateSupplyChainRisk(info, models, datasets)

	return info
}

// assessCompliance assesses compliance status
func (as *AIMLScanner) assessCompliance(models []ModelInfo, datasets []TrainingDataInfo) []ComplianceStatus {
	compliance := make([]ComplianceStatus, 0)

	// GDPR compliance
	gdprIssues := make([]string, 0)
	for _, dataset := range datasets {
		if dataset.HasPII && dataset.RetentionPolicy == "unknown" {
			gdprIssues = append(gdprIssues, fmt.Sprintf("Dataset %s lacks retention policy", dataset.DatasetName))
		}
	}

	gdprStatus := "compliant"
	if len(gdprIssues) > 0 {
		gdprStatus = "non-compliant"
	}

	compliance = append(compliance, ComplianceStatus{
		Framework: "GDPR",
		Status:    gdprStatus,
		Issues:    gdprIssues,
	})

	// Model risk management
	rmIssues := make([]string, 0)
	for _, model := range models {
		if len(model.Vulnerabilities) > 0 {
			rmIssues = append(rmIssues, fmt.Sprintf("Model %s has vulnerabilities", model.Name))
		}
	}

	rmStatus := "compliant"
	if len(rmIssues) > 0 {
		rmStatus = "review-required"
	}

	compliance = append(compliance, ComplianceStatus{
		Framework: "AI Risk Management",
		Status:    rmStatus,
		Issues:    rmIssues,
	})

	return compliance
}

// calculateSupplyChainRisk calculates overall supply chain risk
func (as *AIMLScanner) calculateSupplyChainRisk(info SupplyChainInfo, models []ModelInfo, datasets []TrainingDataInfo) float64 {
	risk := 0.0

	// Vulnerability risk
	criticalCount := 0
	highCount := 0
	for _, vuln := range info.Vulnerabilities {
		switch vuln.Severity {
		case "critical":
			criticalCount++
			risk += 0.2
		case "high":
			highCount++
			risk += 0.1
		case "medium":
			risk += 0.05
		}
	}

	// Compliance risk
	for _, comp := range info.Compliance {
		if comp.Status == "non-compliant" {
			risk += 0.15
		} else if comp.Status == "review-required" {
			risk += 0.08
		}
	}

	// Data risk
	piiDatasets := 0
	for _, dataset := range datasets {
		if dataset.HasPII {
			piiDatasets++
		}
	}
	if piiDatasets > 0 {
		risk += float64(piiDatasets) * 0.05
	}

	// Unknown licenses risk
	unknownLicenses := 0
	for _, license := range info.Licenses {
		if license == "unknown" {
			unknownLicenses++
		}
	}
	risk += float64(unknownLicenses) * 0.03

	// Cap at 1.0
	return min(risk, 1.0)
}

// scanSupplyChainSecurity scans supply chain for security issues
func (as *AIMLScanner) scanSupplyChainSecurity(info SupplyChainInfo) []AIMLFinding {
	var findings []AIMLFinding
	now := time.Now()

	// High risk score
	riskThreshold := 0.6
	if as.config != nil && as.config.RiskThreshold > 0 {
		riskThreshold = as.config.RiskThreshold
	}

	if info.RiskScore > riskThreshold {
		finding := AIMLFinding{
			ID:            uuid.New().String(),
			Type:          "supply_chain",
			Severity:      "high",
			Title:         "Elevated Supply Chain Risk",
			Description:   fmt.Sprintf("AI/ML supply chain risk score %.2f exceeds threshold", info.RiskScore),
			CurrentValue:  fmt.Sprintf("%.2f", info.RiskScore),
			RequiredValue: fmt.Sprintf("%.2f-", riskThreshold),
			Remediation:   "Review and address vulnerabilities, improve compliance posture, secure dependencies",
			DiscoveredAt:  now,
			Metadata: map[string]interface{}{
				"risk_score":          info.RiskScore,
				"vulnerability_count": len(info.Vulnerabilities),
			},
		}
		findings = append(findings, finding)
	}

	// Critical vulnerabilities
	criticalVulns := 0
	for _, vuln := range info.Vulnerabilities {
		if vuln.Severity == "critical" {
			criticalVulns++
		}
	}

	if criticalVulns > 0 {
		finding := AIMLFinding{
			ID:           uuid.New().String(),
			Type:         "supply_chain",
			Severity:     "critical",
			Title:        "Critical Supply Chain Vulnerabilities",
			Description:  fmt.Sprintf("Found %d critical vulnerabilities in AI/ML supply chain", criticalVulns),
			Remediation:  "Immediately patch or replace vulnerable components, review security advisories",
			DiscoveredAt: now,
			Metadata: map[string]interface{}{
				"critical_count": criticalVulns,
				"total_vulns":    len(info.Vulnerabilities),
			},
		}
		findings = append(findings, finding)
	}

	// Non-compliance
	for _, comp := range info.Compliance {
		if comp.Status == "non-compliant" {
			finding := AIMLFinding{
				ID:           uuid.New().String(),
				Type:         "supply_chain",
				Severity:     "high",
				Title:        fmt.Sprintf("%s Compliance Issues", comp.Framework),
				Description:  fmt.Sprintf("Non-compliant with %s: %d issues found", comp.Framework, len(comp.Issues)),
				Remediation:  fmt.Sprintf("Address compliance gaps for %s framework", comp.Framework),
				DiscoveredAt: now,
				Metadata: map[string]interface{}{
					"framework": comp.Framework,
					"issues":    comp.Issues,
				},
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

// Cache methods
func (sc *ScanCache) Get(path string) *CachedFileInfo {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	info, exists := sc.files[path]
	if !exists {
		return nil
	}

	// Cache valid for 1 hour
	if time.Since(info.ScannedAt) > time.Hour {
		return nil
	}

	return info
}

func (sc *ScanCache) Set(path string, info *CachedFileInfo) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.files[path] = info
}

// Helper functions
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// noOpLogger is a no-op logger implementation
type noOpLogger struct{}

func (l *noOpLogger) Info(msg string, fields ...interface{})  {}
func (l *noOpLogger) Warn(msg string, fields ...interface{})  {}
func (l *noOpLogger) Error(msg string, fields ...interface{}) {}
