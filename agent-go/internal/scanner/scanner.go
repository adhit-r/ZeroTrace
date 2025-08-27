package scanner

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	// Scan current directory
	files, err := s.scanFiles(".")
	if err != nil {
		return nil, err
	}

	// Process files for vulnerabilities
	vulnerabilities, dependencies, err := s.analyzeFiles(files)
	if err != nil {
		return nil, err
	}

	result.Vulnerabilities = vulnerabilities
	result.Dependencies = dependencies
	result.EndTime = time.Now()
	result.Metadata["files_scanned"] = len(files)
	result.Metadata["scan_duration"] = result.EndTime.Sub(startTime).String()

	return result, nil
}

// scanFiles scans the directory for files matching include/exclude patterns
func (s *Scanner) scanFiles(root string) ([]models.FileInfo, error) {
	var files []models.FileInfo

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

		// Get file hash
		hash, err := s.getFileHash(path)
		if err != nil {
			return err
		}

		fileInfo := models.FileInfo{
			Path:        path,
			Size:        info.Size(),
			ModTime:     info.ModTime(),
			Language:    s.detectLanguage(path),
			LinesOfCode: s.countLines(path),
			Hash:        hash,
		}

		files = append(files, fileInfo)
		return nil
	})

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
	// TODO: Implement actual vulnerability detection
	// For now, return mock vulnerabilities for testing

	var vulnerabilities []models.Vulnerability

	// Mock vulnerability for Go files
	if file.Language == "go" {
		vuln := models.Vulnerability{
			ID:          generateVulnID(),
			Type:        "dependency",
			Severity:    "medium",
			Title:       "Outdated Go Module",
			Description: "Go module is using an outdated version",
			Location:    file.Path,
			Status:      "open",
			Priority:    "medium",
			CreatedAt:   time.Now(),
		}
		vulnerabilities = append(vulnerabilities, vuln)
	}

	return vulnerabilities, nil
}

// analyzeDependencies analyzes a file for dependencies
func (s *Scanner) analyzeDependencies(file models.FileInfo) ([]models.Dependency, error) {
	// TODO: Implement actual dependency detection
	// For now, return mock dependencies for testing

	var dependencies []models.Dependency

	// Mock dependency for Go files
	if file.Language == "go" {
		dep := models.Dependency{
			ID:        generateDepID(),
			Name:      "github.com/example/module",
			Version:   "v1.0.0",
			Type:      "go",
			Location:  file.Path,
			CreatedAt: time.Now(),
		}
		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
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
