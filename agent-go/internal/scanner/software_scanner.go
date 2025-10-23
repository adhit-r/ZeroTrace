package scanner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/models"

	"github.com/google/uuid"
)

// SoftwareScanner handles scanning for installed software applications
type SoftwareScanner struct {
	config *config.Config
}

// NewSoftwareScanner creates a new software scanner instance
func NewSoftwareScanner(cfg *config.Config) *SoftwareScanner {
	return &SoftwareScanner{
		config: cfg,
	}
}

// Scan performs a software vulnerability scan
func (s *SoftwareScanner) Scan() (*models.ScanResult, error) {
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

	// Detect installed software based on OS
	var installedApps []models.InstalledApp
	var err error

	switch runtime.GOOS {
	case "darwin":
		installedApps, err = s.scanMacOS()
	case "linux":
		installedApps, err = s.scanLinux()
	case "windows":
		installedApps, err = s.scanWindows()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err != nil {
		return nil, err
	}

	// Convert installed apps to dependencies for processing
	dependencies := s.convertAppsToDependencies(installedApps)
	result.Dependencies = dependencies
	result.EndTime = time.Now()
	result.Metadata["apps_scanned"] = len(installedApps)
	result.Metadata["scan_duration"] = result.EndTime.Sub(startTime).String()
	result.Metadata["os"] = runtime.GOOS
	result.Metadata["arch"] = runtime.GOARCH

	return result, nil
}

// scanMacOS scans for installed applications on macOS
func (s *SoftwareScanner) scanMacOS() ([]models.InstalledApp, error) {
	var apps []models.InstalledApp

	// Common application directories
	appDirs := []string{
		"/Applications",
		"/System/Applications",
		os.Getenv("HOME") + "/Applications",
	}

	for _, dir := range appDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && strings.HasSuffix(path, ".app") {
				app := s.extractMacOSAppInfo(path)
				if app.Name != "" {
					apps = append(apps, app)
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	// Also check Homebrew packages
	brewApps, err := s.scanHomebrew()
	if err == nil {
		apps = append(apps, brewApps...)
	}

	return apps, nil
}

// scanLinux scans for installed applications on Linux
func (s *SoftwareScanner) scanLinux() ([]models.InstalledApp, error) {
	var apps []models.InstalledApp

	// Check package managers
	packageManagers := []struct {
		name string
		cmd  string
		args []string
	}{
		{"apt", "dpkg", []string{"-l"}},
		{"yum", "rpm", []string{"-qa"}},
		{"pacman", "pacman", []string{"-Q"}},
		{"snap", "snap", []string{"list"}},
		{"flatpak", "flatpak", []string{"list"}},
	}

	for _, pm := range packageManagers {
		if _, err := exec.LookPath(pm.cmd); err == nil {
			pmApps, err := s.scanPackageManager(pm.name, pm.cmd, pm.args)
			if err == nil {
				apps = append(apps, pmApps...)
			}
		}
	}

	return apps, nil
}

// scanWindows scans for installed applications on Windows
func (s *SoftwareScanner) scanWindows() ([]models.InstalledApp, error) {
	var apps []models.InstalledApp

	// This would use Windows registry or PowerShell commands
	// For now, return empty list as we're on macOS/Linux
	return apps, nil
}

// extractMacOSAppInfo extracts information from a macOS .app bundle
func (s *SoftwareScanner) extractMacOSAppInfo(appPath string) models.InstalledApp {
	app := models.InstalledApp{
		Path: appPath,
		Type: "macos_app",
	}

	// Extract app name from path
	app.Name = strings.TrimSuffix(filepath.Base(appPath), ".app")

	// Try to get version from Info.plist
	infoPlistPath := filepath.Join(appPath, "Contents", "Info.plist")
	if _, err := os.Stat(infoPlistPath); err == nil {
		// Extract version using plutil command
		cmd := exec.Command("plutil", "-p", infoPlistPath)
		output, err := cmd.Output()
		if err == nil {
			// Parse the plist output to find CFBundleShortVersionString
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "CFBundleShortVersionString") {
					// Extract version from line like: "CFBundleShortVersionString" => "1.7.52"
					parts := strings.Split(line, "=>")
					if len(parts) > 1 {
						version := strings.TrimSpace(parts[1])
						version = strings.Trim(version, "\"")
						version = strings.TrimSpace(version)
						app.Version = version
						break
					}
				}
			}
		}
		if app.Version == "" {
			app.Version = "unknown"
		}
	} else {
		app.Version = "unknown"
	}

	// Get file info
	if info, err := os.Stat(appPath); err == nil {
		app.InstallDate = info.ModTime()
		app.Size = info.Size()
	}

	return app
}

// scanHomebrew scans for Homebrew packages
func (s *SoftwareScanner) scanHomebrew() ([]models.InstalledApp, error) {
	var apps []models.InstalledApp

	cmd := exec.Command("brew", "list", "--formula")
	output, err := cmd.Output()
	if err != nil {
		return apps, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Get version info
		versionCmd := exec.Command("brew", "list", "--versions", line)
		versionOutput, err := versionCmd.Output()
		if err != nil {
			continue
		}

		version := strings.TrimSpace(string(versionOutput))
		if version == "" {
			continue
		}

		app := models.InstalledApp{
			Name:    line,
			Version: version,
			Type:    "homebrew",
			Path:    "/opt/homebrew/bin/" + line, // Common Homebrew path
		}

		apps = append(apps, app)
	}

	return apps, nil
}

// scanPackageManager scans a specific package manager
func (s *SoftwareScanner) scanPackageManager(pmName, cmd string, args []string) ([]models.InstalledApp, error) {
	var apps []models.InstalledApp

	execCmd := exec.Command(cmd, args...)
	output, err := execCmd.Output()
	if err != nil {
		return apps, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse package info (simplified)
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			app := models.InstalledApp{
				Name:    parts[0],
				Version: parts[1],
				Type:    pmName,
			}
			apps = append(apps, app)
		}
	}

	return apps, nil
}

// convertAppsToDependencies converts installed apps to dependency format
func (s *SoftwareScanner) convertAppsToDependencies(apps []models.InstalledApp) []models.Dependency {
	var dependencies []models.Dependency

	for _, app := range apps {
		dep := models.Dependency{
			Name:        app.Name,
			Version:     app.Version,
			Type:        app.Type,
			Path:        app.Path,
			InstallDate: app.InstallDate,
			Size:        app.Size,
			Vendor:      s.detectVendor(app.Name),
			Description: fmt.Sprintf("Installed %s application", app.Type),
		}
		dependencies = append(dependencies, dep)
	}

	return dependencies
}

// detectVendor attempts to detect the vendor from the app name
func (s *SoftwareScanner) detectVendor(appName string) string {
	appName = strings.ToLower(appName)

	vendorMap := map[string]string{
		"chrome":    "Google",
		"firefox":   "Mozilla",
		"safari":    "Apple",
		"edge":      "Microsoft",
		"adobe":     "Adobe",
		"vlc":       "VideoLAN",
		"7zip":      "7-Zip",
		"notepad++": "Notepad++",
		"vscode":    "Microsoft",
		"intellij":  "JetBrains",
		"docker":    "Docker",
		"kubectl":   "Kubernetes",
		"git":       "Git",
		"node":      "Node.js",
		"python":    "Python",
		"java":      "Oracle",
		"maven":     "Apache",
		"gradle":    "Gradle",
	}

	for key, vendor := range vendorMap {
		if strings.Contains(appName, key) {
			return vendor
		}
	}

	return "Unknown"
}
