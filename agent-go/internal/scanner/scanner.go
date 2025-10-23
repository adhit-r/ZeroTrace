package scanner

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/models"

	"github.com/google/uuid"
)

// Scanner handles vulnerability scanning
type Scanner struct {
	config *config.Config
}

// NewScanner creates a new scanner instance
func NewScanner(cfg *config.Config) *Scanner {
	return &Scanner{
		config: cfg,
	}
}

// Scan performs a vulnerability scan
func (s *Scanner) Scan() (*models.ScanResult, error) {
	startTime := time.Now()

	// Create scan result
	result := &models.ScanResult{
		ID:              uuid.New(),
		AgentID:         s.config.AgentID,
		CompanyID:       s.config.CompanyID,
		StartTime:       startTime,
		EndTime:         time.Now(),
		Status:          "completed",
		Vulnerabilities: []models.Vulnerability{},
		Dependencies:    []models.Dependency{},
		Metadata:        make(map[string]any),
	}

	// Performance metrics
	scanMetrics := map[string]interface{}{
		"scan_start_time": startTime.Unix(),
		"agent_id":        s.config.AgentID,
		"scan_type":       "vulnerability_scan",
	}

	// Scan current directory
	fileScanStart := time.Now()
	files, err := s.scanFiles(".")
	fileScanDuration := time.Since(fileScanStart)
	if err != nil {
		return nil, err
	}

	// Process files for vulnerabilities
	analysisStart := time.Now()
	vulnerabilities, dependencies, err := s.analyzeFiles(files)
	analysisDuration := time.Since(analysisStart)
	if err != nil {
		return nil, err
	}

	result.Vulnerabilities = vulnerabilities
	result.Dependencies = dependencies
	result.EndTime = time.Now()

	// Enhanced performance metrics
	totalDuration := result.EndTime.Sub(startTime)
	scanMetrics["file_scan_duration_ms"] = fileScanDuration.Milliseconds()
	scanMetrics["analysis_duration_ms"] = analysisDuration.Milliseconds()
	scanMetrics["total_duration_ms"] = totalDuration.Milliseconds()
	scanMetrics["files_scanned"] = len(files)
	scanMetrics["vulnerabilities_found"] = len(vulnerabilities)
	scanMetrics["dependencies_found"] = len(dependencies)
	scanMetrics["scan_end_time"] = result.EndTime.Unix()

	result.Metadata["performance_metrics"] = scanMetrics
	result.Metadata["files_scanned"] = len(files)
	result.Metadata["scan_duration"] = totalDuration.String()

	return result, nil
}

// scanFiles scans the directory for files matching include/exclude patterns
func (s *Scanner) scanFiles(root string) ([]models.FileInfo, error) {
	var files []models.FileInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Use a semaphore to limit concurrent file processing
	semaphore := make(chan struct{}, 10) // Limit to 10 concurrent files

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check file size
		if info.Size() > s.config.MaxFileSize {
			return nil
		}

		// Check include/exclude patterns
		if !s.shouldScanFile(path) {
			return nil
		}

		// Process file in goroutine for better performance
		wg.Add(1)
		go func(p string, i os.FileInfo) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Get file hash
			hash, err := s.getFileHash(p)
			if err != nil {
				return
			}

			fileInfo := models.FileInfo{
				Path:        p,
				Size:        i.Size(),
				ModTime:     i.ModTime(),
				Language:    s.detectLanguage(p),
				LinesOfCode: s.countLines(p),
				Hash:        hash,
			}

			mu.Lock()
			files = append(files, fileInfo)
			mu.Unlock()
		}(path, info)

		return nil
	})

	// Wait for all goroutines to complete
	wg.Wait()

	return files, err
}

// shouldScanFile determines if a file should be scanned based on patterns
func (s *Scanner) shouldScanFile(path string) bool {
	// Check exclude patterns first
	for _, pattern := range s.config.ExcludePatterns {
		if strings.Contains(path, pattern) {
			return false
		}
	}

	// Check include patterns
	for _, pattern := range s.config.IncludePatterns {
		if strings.HasSuffix(path, pattern) {
			return true
		}
	}

	return false
}

// getFileHash calculates SHA256 hash of a file
func (s *Scanner) getFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// detectLanguage detects the programming language of a file
func (s *Scanner) detectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go":
		return "go"
	case ".py":
		return "python"
	case ".js":
		return "javascript"
	case ".ts":
		return "typescript"
	case ".java":
		return "java"
	case ".php":
		return "php"
	case ".rb":
		return "ruby"
	case ".rs":
		return "rust"
	case ".cpp", ".cc", ".cxx":
		return "cpp"
	case ".c":
		return "c"
	case ".cs":
		return "csharp"
	default:
		return "unknown"
	}
}

// countLines counts the number of lines in a file
func (s *Scanner) countLines(path string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return 0
	}

	return strings.Count(string(content), "\n") + 1
}

// analyzeFiles analyzes files for vulnerabilities and dependencies
func (s *Scanner) analyzeFiles(files []models.FileInfo) ([]models.Vulnerability, []models.Dependency, error) {
	var vulnerabilities []models.Vulnerability
	var dependencies []models.Dependency

	for _, file := range files {
		// Analyze for vulnerabilities
		fileVulns, err := s.analyzeVulnerabilities(file)
		if err != nil {
			continue
		}
		vulnerabilities = append(vulnerabilities, fileVulns...)

		// Analyze for dependencies
		fileDeps, err := s.analyzeDependencies(file)
		if err != nil {
			continue
		}
		dependencies = append(dependencies, fileDeps...)
	}

	return vulnerabilities, dependencies, nil
}

// analyzeVulnerabilities analyzes a file for vulnerabilities
func (s *Scanner) analyzeVulnerabilities(file models.FileInfo) ([]models.Vulnerability, error) {
	// Agent no longer performs local vulnerability detection
	// Dependencies are sent to API, which handles enrichment via Python service
	// This ensures consistent CVE detection across all agents

	return []models.Vulnerability{}, nil
}

// analyzeDependencies analyzes a file for dependencies
func (s *Scanner) analyzeDependencies(file models.FileInfo) ([]models.Dependency, error) {
	// Scan actual package managers for real dependencies
	var dependencies []models.Dependency

	// Scan Go modules
	if file.Language == "go" && strings.HasSuffix(file.Path, "go.mod") {
		deps := s.scanGoMod(file.Path)
		dependencies = append(dependencies, deps...)
	}

	// Scan Node.js packages
	if file.Language == "javascript" && strings.HasSuffix(file.Path, "package.json") {
		deps := s.scanPackageJson(file.Path)
		dependencies = append(dependencies, deps...)
	}

	// Scan Python requirements
	if file.Language == "python" && (strings.HasSuffix(file.Path, "requirements.txt") || strings.HasSuffix(file.Path, "pyproject.toml")) {
		deps := s.scanPythonDeps(file.Path)
		dependencies = append(dependencies, deps...)
	}

	return dependencies, nil
}

// scanGoMod scans go.mod file for dependencies
func (s *Scanner) scanGoMod(filePath string) []models.Dependency {
	var dependencies []models.Dependency

	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require ") {
			// Parse require line: require module version
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				dep := models.Dependency{
					ID:        generateDepID(),
					Name:      parts[1],
					Version:   parts[2],
					Type:      "go",
					Location:  filePath,
					CreatedAt: time.Now(),
				}
				dependencies = append(dependencies, dep)
			}
		}
	}

	return dependencies
}

// scanPackageJson scans package.json file for dependencies
func (s *Scanner) scanPackageJson(filePath string) []models.Dependency {
	var dependencies []models.Dependency

	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies
	}

	// Simple JSON parsing for dependencies
	lines := strings.Split(string(content), "\n")
	inDependencies := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "\"dependencies\"") {
			inDependencies = true
			continue
		}
		if inDependencies && strings.Contains(line, "}") {
			break
		}
		if inDependencies && strings.Contains(line, "\"") {
			// Parse dependency line: "name": "version"
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				name := strings.Trim(strings.Trim(parts[0], " "), "\"")
				version := strings.Trim(strings.Trim(parts[1], " ,"), "\"")
				if name != "" && version != "" {
					dep := models.Dependency{
						ID:        generateDepID(),
						Name:      name,
						Version:   version,
						Type:      "javascript",
						Location:  filePath,
						CreatedAt: time.Now(),
					}
					dependencies = append(dependencies, dep)
				}
			}
		}
	}

	return dependencies
}

// scanPythonDeps scans Python requirements files
func (s *Scanner) scanPythonDeps(filePath string) []models.Dependency {
	var dependencies []models.Dependency

	content, err := os.ReadFile(filePath)
	if err != nil {
		return dependencies
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			// Parse requirement line: package==version or package>=version
			parts := strings.FieldsFunc(line, func(c rune) bool {
				return c == '=' || c == '>' || c == '<' || c == '!'
			})
			if len(parts) >= 1 {
				name := strings.TrimSpace(parts[0])
				version := "unknown"
				if len(parts) >= 2 {
					version = strings.TrimSpace(parts[1])
				}
				dep := models.Dependency{
					ID:        generateDepID(),
					Name:      name,
					Version:   version,
					Type:      "python",
					Location:  filePath,
					CreatedAt: time.Now(),
				}
				dependencies = append(dependencies, dep)
			}
		}
	}

	return dependencies
}

// Helper functions for generating IDs
func generateScanID() string {
	return fmt.Sprintf("scan-%d", time.Now().Unix())
}

func generateVulnID() string {
	return fmt.Sprintf("vuln-%d", time.Now().UnixNano())
}

func generateDepID() string {
	return fmt.Sprintf("dep-%d", time.Now().UnixNano())
}
