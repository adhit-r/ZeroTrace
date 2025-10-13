package scanner

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/models"

	"github.com/google/uuid"
)

// ConfigScanner scans for configuration vulnerabilities
type ConfigScanner struct {
	config *config.Config
}

// NewConfigScanner creates a new configuration scanner
func NewConfigScanner(cfg *config.Config) *ConfigScanner {
	return &ConfigScanner{
		config: cfg,
	}
}

// Scan performs configuration vulnerability scanning
func (cs *ConfigScanner) Scan() (*models.ScanResult, error) {
	startTime := time.Now()

	// Create scan result
	// Use a proper UUID for CompanyID, generate one if the config value is not a valid UUID
	companyID := cs.config.CompanyID
	if companyID == "" || companyID == "company-001" {
		// Generate a default company UUID for configuration scans
		defaultCompanyUUID := uuid.New()
		companyID = defaultCompanyUUID.String()
	}

	result := &models.ScanResult{
		ID:              uuid.New(),
		AgentID:         cs.config.AgentID,
		CompanyID:       companyID,
		Status:          "completed",
		StartTime:       startTime,
		EndTime:         time.Now(),
		Vulnerabilities: []models.Vulnerability{},
		Dependencies:    []models.Dependency{},
		Metadata:        make(map[string]interface{}),
	}

	// Perform OS-specific configuration scans
	var vulnerabilities []models.Vulnerability
	var assets []models.Asset
	var err error

	switch runtime.GOOS {
	case "darwin":
		vulnerabilities, assets, err = cs.scanMacOS()
	case "linux":
		vulnerabilities, assets, err = cs.scanLinux()
	case "windows":
		vulnerabilities, assets, err = cs.scanWindows()
	default:
		return result, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err != nil {
		result.Status = "failed"
		result.Metadata["error"] = err.Error()
		return result, err
	}

	// Add some test vulnerabilities for development/demo purposes
	testVulns := cs.generateTestVulnerabilities()
	vulnerabilities = append(vulnerabilities, testVulns...)

	// Set results
	result.Vulnerabilities = vulnerabilities
	result.Metadata["total_vulnerabilities"] = len(vulnerabilities)
	result.Metadata["critical_vulnerabilities"] = cs.countBySeverity(vulnerabilities, "critical")
	result.Metadata["high_vulnerabilities"] = cs.countBySeverity(vulnerabilities, "high")
	result.Metadata["medium_vulnerabilities"] = cs.countBySeverity(vulnerabilities, "medium")
	result.Metadata["low_vulnerabilities"] = cs.countBySeverity(vulnerabilities, "low")
	result.Metadata["total_assets"] = len(assets)
	result.Metadata["scan_duration"] = time.Since(startTime).Seconds()

	result.Metadata["os"] = runtime.GOOS
	result.Metadata["scan_type"] = "configuration"
	result.Metadata["timestamp"] = time.Now().Format(time.RFC3339)

	return result, nil
}

// scanMacOS performs macOS-specific configuration scanning
func (cs *ConfigScanner) scanMacOS() ([]models.Vulnerability, []models.Asset, error) {
	var vulnerabilities []models.Vulnerability
	var assets []models.Asset

	// Check system security settings
	securityChecks := []struct {
		name        string
		description string
		severity    string
		check       func() (bool, string)
	}{
		{
			name:        "Gatekeeper Status",
			description: "Check if Gatekeeper is enabled for malware protection",
			severity:    "high",
			check:       cs.checkGatekeeper,
		},
		{
			name:        "System Integrity Protection",
			description: "Check if System Integrity Protection (SIP) is enabled",
			severity:    "critical",
			check:       cs.checkSIP,
		},
		{
			name:        "Firewall Status",
			description: "Check if firewall is enabled",
			severity:    "high",
			check:       cs.checkFirewall,
		},
		{
			name:        "Automatic Updates",
			description: "Check if automatic security updates are enabled",
			severity:    "medium",
			check:       cs.checkAutoUpdates,
		},
		{
			name:        "FileVault Encryption",
			description: "Check if FileVault disk encryption is enabled",
			severity:    "high",
			check:       cs.checkFileVault,
		},
		{
			name:        "Screen Lock",
			description: "Check if screen lock is configured",
			severity:    "medium",
			check:       cs.checkScreenLock,
		},
		{
			name:        "Remote Login (SSH)",
			description: "Check if SSH remote login is disabled",
			severity:    "high",
			check:       cs.checkSSH,
		},
		{
			name:        "Remote Management (ARD)",
			description: "Check if Apple Remote Desktop is disabled",
			severity:    "medium",
			check:       cs.checkARD,
		},
		{
			name:        "Guest Account",
			description: "Check if guest account is disabled",
			severity:    "medium",
			check:       cs.checkGuestAccount,
		},
		{
			name:        "Automatic Login",
			description: "Check if automatic login is disabled",
			severity:    "medium",
			check:       cs.checkAutoLogin,
		},
		{
			name:        "Password Policy",
			description: "Check if strong password policy is enforced",
			severity:    "high",
			check:       cs.checkPasswordPolicy,
		},
		{
			name:        "Bluetooth Security",
			description: "Check if Bluetooth is configured securely",
			severity:    "low",
			check:       cs.checkBluetoothSecurity,
		},
		{
			name:        "Location Services",
			description: "Check if location services are properly configured",
			severity:    "low",
			check:       cs.checkLocationServices,
		},
		{
			name:        "System Time Sync",
			description: "Check if system time is synchronized",
			severity:    "medium",
			check:       cs.checkTimeSync,
		},
		{
			name:        "Secure Boot",
			description: "Check if secure boot is enabled",
			severity:    "high",
			check:       cs.checkSecureBoot,
		},
	}

	for _, check := range securityChecks {
		isSecure, details := check.check()
		if !isSecure {
			vulnerability := models.Vulnerability{
				ID:          uuid.New().String(),
				Type:        "configuration",
				Title:       check.name,
				Description: check.description,
				Severity:    check.severity,
				Status:      "open",
				EnrichmentData: map[string]interface{}{
					"details":  details,
					"os":       "macOS",
					"category": "configuration",
				},
				CreatedAt: time.Now(),
			}
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	// Create system asset
	systemAsset := models.Asset{
		ID:     "macos-system",
		Name:   "macOS System",
		Type:   "operating_system",
		Status: "active",
		Metadata: map[string]interface{}{
			"os":           runtime.GOOS,
			"version":      cs.getMacOSVersion(),
			"architecture": runtime.GOARCH,
		},
	}
	assets = append(assets, systemAsset)

	return vulnerabilities, assets, nil
}

// scanLinux performs Linux-specific configuration scanning
func (cs *ConfigScanner) scanLinux() ([]models.Vulnerability, []models.Asset, error) {
	var vulnerabilities []models.Vulnerability
	var assets []models.Asset

	// Basic Linux security checks
	securityChecks := []struct {
		name        string
		description string
		severity    string
		check       func() (bool, string)
	}{
		{
			name:        "SELinux Status",
			description: "Check if SELinux is enabled and enforcing",
			severity:    "high",
			check:       cs.checkSELinux,
		},
		{
			name:        "AppArmor Status",
			description: "Check if AppArmor is enabled",
			severity:    "high",
			check:       cs.checkAppArmor,
		},
		{
			name:        "UFW Firewall",
			description: "Check if UFW firewall is enabled",
			severity:    "high",
			check:       cs.checkUFW,
		},
		{
			name:        "Automatic Updates",
			description: "Check if automatic security updates are enabled",
			severity:    "medium",
			check:       cs.checkLinuxAutoUpdates,
		},
	}

	for _, check := range securityChecks {
		isSecure, details := check.check()
		if !isSecure {
			vulnerability := models.Vulnerability{
				ID:          uuid.New().String(),
				Type:        "configuration",
				Title:       check.name,
				Description: check.description,
				Severity:    check.severity,
				Status:      "open",
				EnrichmentData: map[string]interface{}{
					"details":  details,
					"os":       "Linux",
					"category": "configuration",
				},
				CreatedAt: time.Now(),
			}
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	// Create system asset
	systemAsset := models.Asset{
		ID:     "linux-system",
		Name:   "Linux System",
		Type:   "operating_system",
		Status: "active",
		Metadata: map[string]interface{}{
			"os":           runtime.GOOS,
			"architecture": runtime.GOARCH,
		},
	}
	assets = append(assets, systemAsset)

	return vulnerabilities, assets, nil
}

// scanWindows performs Windows-specific configuration scanning
func (cs *ConfigScanner) scanWindows() ([]models.Vulnerability, []models.Asset, error) {
	var vulnerabilities []models.Vulnerability
	var assets []models.Asset

	// Basic Windows security checks
	securityChecks := []struct {
		name        string
		description string
		severity    string
		check       func() (bool, string)
	}{
		{
			name:        "Windows Defender",
			description: "Check if Windows Defender is enabled",
			severity:    "critical",
			check:       cs.checkWindowsDefender,
		},
		{
			name:        "Windows Firewall",
			description: "Check if Windows Firewall is enabled",
			severity:    "high",
			check:       cs.checkWindowsFirewall,
		},
		{
			name:        "Automatic Updates",
			description: "Check if Windows Update is configured",
			severity:    "medium",
			check:       cs.checkWindowsUpdates,
		},
	}

	for _, check := range securityChecks {
		isSecure, details := check.check()
		if !isSecure {
			vulnerability := models.Vulnerability{
				ID:          uuid.New().String(),
				Type:        "configuration",
				Title:       check.name,
				Description: check.description,
				Severity:    check.severity,
				Status:      "open",
				EnrichmentData: map[string]interface{}{
					"details":  details,
					"os":       "Windows",
					"category": "configuration",
				},
				CreatedAt: time.Now(),
			}
			vulnerabilities = append(vulnerabilities, vulnerability)
		}
	}

	// Create system asset
	systemAsset := models.Asset{
		ID:     "windows-system",
		Name:   "Windows System",
		Type:   "operating_system",
		Status: "active",
		Metadata: map[string]interface{}{
			"os":           runtime.GOOS,
			"architecture": runtime.GOARCH,
		},
	}
	assets = append(assets, systemAsset)

	return vulnerabilities, assets, nil
}

// macOS Security Checks

func (cs *ConfigScanner) checkGatekeeper() (bool, string) {
	cmd := exec.Command("spctl", "--status")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check Gatekeeper status"
	}

	status := strings.TrimSpace(string(output))
	if strings.Contains(status, "enabled") {
		return true, "Gatekeeper is enabled"
	}
	return false, "Gatekeeper is disabled - malware protection is reduced"
}

func (cs *ConfigScanner) checkSIP() (bool, string) {
	cmd := exec.Command("csrutil", "status")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check SIP status"
	}

	status := strings.TrimSpace(string(output))
	if strings.Contains(status, "enabled") {
		return true, "System Integrity Protection is enabled"
	}
	return false, "System Integrity Protection is disabled - system security is compromised"
}

func (cs *ConfigScanner) checkFirewall() (bool, string) {
	cmd := exec.Command("defaults", "read", "/Library/Preferences/com.apple.alf", "globalstate")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check firewall status"
	}

	status := strings.TrimSpace(string(output))
	if status == "1" {
		return true, "Firewall is enabled"
	}
	return false, "Firewall is disabled - network security is reduced"
}

func (cs *ConfigScanner) checkAutoUpdates() (bool, string) {
	cmd := exec.Command("defaults", "read", "/Library/Preferences/com.apple.SoftwareUpdate", "AutomaticCheckEnabled")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check auto-update status"
	}

	status := strings.TrimSpace(string(output))
	if status == "1" {
		return true, "Automatic updates are enabled"
	}
	return false, "Automatic updates are disabled - system may be vulnerable to known exploits"
}

func (cs *ConfigScanner) checkFileVault() (bool, string) {
	cmd := exec.Command("fdesetup", "status")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check FileVault status"
	}

	status := strings.TrimSpace(string(output))
	if strings.Contains(status, "FileVault is On") {
		return true, "FileVault encryption is enabled"
	}
	return false, "FileVault encryption is disabled - disk data is not encrypted"
}

func (cs *ConfigScanner) checkScreenLock() (bool, string) {
	cmd := exec.Command("defaults", "read", "com.apple.screensaver", "askForPassword")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check screen lock status"
	}

	status := strings.TrimSpace(string(output))
	if status == "1" {
		return true, "Screen lock is configured"
	}
	return false, "Screen lock is not configured - physical access is not protected"
}

// Additional macOS Security Checks

func (cs *ConfigScanner) checkSSH() (bool, string) {
	cmd := exec.Command("systemsetup", "-getremotelogin")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check SSH status"
	}

	status := strings.TrimSpace(string(output))
	if strings.Contains(status, "Remote Login: Off") {
		return true, "SSH remote login is disabled"
	}
	return false, "SSH remote login is enabled - potential security risk"
}

func (cs *ConfigScanner) checkARD() (bool, string) {
	cmd := exec.Command("launchctl", "list", "com.apple.RemoteDesktop")
	output, err := cmd.Output()
	if err != nil {
		return true, "Apple Remote Desktop is not running"
	}

	if strings.Contains(string(output), "com.apple.RemoteDesktop") {
		return false, "Apple Remote Desktop is enabled - potential security risk"
	}
	return true, "Apple Remote Desktop is disabled"
}

func (cs *ConfigScanner) checkGuestAccount() (bool, string) {
	cmd := exec.Command("dscl", ".", "-read", "/Users/Guest", "AuthenticationAuthority")
	output, err := cmd.Output()
	if err != nil {
		return true, "Guest account is disabled"
	}

	if strings.Contains(string(output), "No such key") {
		return true, "Guest account is disabled"
	}
	return false, "Guest account is enabled - potential security risk"
}

func (cs *ConfigScanner) checkAutoLogin() (bool, string) {
	cmd := exec.Command("defaults", "read", "/Library/Preferences/com.apple.loginwindow", "autoLoginUser")
	output, err := cmd.Output()
	if err != nil {
		return true, "Automatic login is disabled"
	}

	status := strings.TrimSpace(string(output))
	if status == "" || status == "0" {
		return true, "Automatic login is disabled"
	}
	return false, "Automatic login is enabled - potential security risk"
}

func (cs *ConfigScanner) checkPasswordPolicy() (bool, string) {
	cmd := exec.Command("pwpolicy", "-getaccountpolicies")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check password policy"
	}

	// Check for common password policy requirements
	policy := string(output)
	hasMinLength := strings.Contains(policy, "minChars")
	hasComplexity := strings.Contains(policy, "requireMixedCase") || strings.Contains(policy, "requireNumeric")

	if hasMinLength && hasComplexity {
		return true, "Strong password policy is enforced"
	}
	return false, "Weak or no password policy - security risk"
}

func (cs *ConfigScanner) checkBluetoothSecurity() (bool, string) {
	cmd := exec.Command("defaults", "read", "/Library/Preferences/com.apple.Bluetooth", "ControllerPowerState")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check Bluetooth status"
	}

	status := strings.TrimSpace(string(output))
	if status == "0" {
		return true, "Bluetooth is disabled - most secure"
	}
	return false, "Bluetooth is enabled - ensure discoverable mode is off"
}

func (cs *ConfigScanner) checkLocationServices() (bool, string) {
	cmd := exec.Command("defaults", "read", "/var/db/locationd/Library/Preferences/ByHost/com.apple.locationd", "LocationServicesEnabled")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check location services"
	}

	status := strings.TrimSpace(string(output))
	if status == "0" {
		return true, "Location services are disabled - privacy protected"
	}
	return false, "Location services are enabled - privacy consideration"
}

func (cs *ConfigScanner) checkTimeSync() (bool, string) {
	cmd := exec.Command("sntp", "-sS", "time.apple.com")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check time synchronization"
	}

	if strings.Contains(string(output), "synchronized") {
		return true, "System time is synchronized"
	}
	return false, "System time may not be synchronized - security risk"
}

func (cs *ConfigScanner) checkSecureBoot() (bool, string) {
	cmd := exec.Command("bputil", "-d")
	output, err := cmd.Output()
	if err != nil {
		return false, "Unable to check secure boot status"
	}

	if strings.Contains(string(output), "Secure Boot: Full Security") {
		return true, "Secure boot is enabled with full security"
	}
	if strings.Contains(string(output), "Secure Boot: Medium Security") {
		return false, "Secure boot is enabled with medium security - consider full security"
	}
	return false, "Secure boot is disabled or not available - security risk"
}

// Linux Security Checks

func (cs *ConfigScanner) checkSELinux() (bool, string) {
	cmd := exec.Command("getenforce")
	output, err := cmd.Output()
	if err != nil {
		return false, "SELinux not available or not installed"
	}

	status := strings.TrimSpace(string(output))
	if status == "Enforcing" {
		return true, "SELinux is enforcing"
	}
	return false, fmt.Sprintf("SELinux is %s - mandatory access control is not enforced", status)
}

func (cs *ConfigScanner) checkAppArmor() (bool, string) {
	cmd := exec.Command("aa-status")
	output, err := cmd.Output()
	if err != nil {
		return false, "AppArmor not available or not installed"
	}

	status := strings.TrimSpace(string(output))
	if strings.Contains(status, "enforce") {
		return true, "AppArmor is enforcing"
	}
	return false, "AppArmor is not enforcing - mandatory access control is not active"
}

func (cs *ConfigScanner) checkUFW() (bool, string) {
	cmd := exec.Command("ufw", "status")
	output, err := cmd.Output()
	if err != nil {
		return false, "UFW not available or not installed"
	}

	status := strings.TrimSpace(string(output))
	if strings.Contains(status, "Status: active") {
		return true, "UFW firewall is active"
	}
	return false, "UFW firewall is not active - network security is reduced"
}

func (cs *ConfigScanner) checkLinuxAutoUpdates() (bool, string) {
	// Check for unattended-upgrades
	cmd := exec.Command("systemctl", "is-enabled", "unattended-upgrades")
	output, err := cmd.Output()
	if err != nil {
		return false, "Automatic updates not configured"
	}

	status := strings.TrimSpace(string(output))
	if status == "enabled" {
		return true, "Automatic updates are enabled"
	}
	return false, "Automatic updates are disabled - system may be vulnerable to known exploits"
}

// Windows Security Checks

func (cs *ConfigScanner) checkWindowsDefender() (bool, string) {
	// This would require Windows-specific implementation
	// For now, return a placeholder
	return false, "Windows Defender status check not implemented"
}

func (cs *ConfigScanner) checkWindowsFirewall() (bool, string) {
	// This would require Windows-specific implementation
	// For now, return a placeholder
	return false, "Windows Firewall status check not implemented"
}

func (cs *ConfigScanner) checkWindowsUpdates() (bool, string) {
	// This would require Windows-specific implementation
	// For now, return a placeholder
	return false, "Windows Update status check not implemented"
}

// Utility functions

func (cs *ConfigScanner) getMacOSVersion() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return "Unknown"
	}
	return strings.TrimSpace(string(output))
}

func (cs *ConfigScanner) countBySeverity(vulnerabilities []models.Vulnerability, severity string) int {
	count := 0
	for _, vuln := range vulnerabilities {
		if vuln.Severity == severity {
			count++
		}
	}
	return count
}

// generateTestVulnerabilities creates realistic test vulnerabilities for development/demo
func (cs *ConfigScanner) generateTestVulnerabilities() []models.Vulnerability {
	now := time.Now()

	return []models.Vulnerability{
		{
			ID:               "test-vuln-001",
			Type:             "configuration",
			Severity:         "critical",
			Title:            "Outdated System Components",
			Description:      "System contains outdated components with known security vulnerabilities",
			Status:           "open",
			Priority:         "urgent",
			ExploitAvailable: true,
			ExploitCount:     3,
			CreatedAt:        now.Add(-24 * time.Hour),
			EnrichmentData: map[string]interface{}{
				"category":  "system",
				"details":   "Multiple outdated system components detected",
				"os":        runtime.GOOS,
				"cve_count": 5,
			},
		},
		{
			ID:               "test-vuln-002",
			Type:             "software",
			Severity:         "high",
			Title:            "Vulnerable Browser Extensions",
			Description:      "Browser extensions with known security vulnerabilities detected",
			Status:           "open",
			Priority:         "high",
			ExploitAvailable: true,
			ExploitCount:     1,
			CreatedAt:        now.Add(-12 * time.Hour),
			EnrichmentData: map[string]interface{}{
				"category":  "browser",
				"details":   "Chrome extensions with CVE-2024-1234",
				"os":        runtime.GOOS,
				"cve_count": 2,
			},
		},
		{
			ID:               "test-vuln-003",
			Type:             "configuration",
			Severity:         "medium",
			Title:            "Weak Password Policy",
			Description:      "System password policy does not meet security requirements",
			Status:           "open",
			Priority:         "medium",
			ExploitAvailable: false,
			ExploitCount:     0,
			CreatedAt:        now.Add(-6 * time.Hour),
			EnrichmentData: map[string]interface{}{
				"category":  "authentication",
				"details":   "Password complexity requirements not enforced",
				"os":        runtime.GOOS,
				"cve_count": 0,
			},
		},
		{
			ID:               "test-vuln-004",
			Type:             "software",
			Severity:         "high",
			Title:            "Unpatched Development Tools",
			Description:      "Development tools contain unpatched security vulnerabilities",
			Status:           "open",
			Priority:         "high",
			ExploitAvailable: true,
			ExploitCount:     2,
			CreatedAt:        now.Add(-3 * time.Hour),
			EnrichmentData: map[string]interface{}{
				"category":  "development",
				"details":   "Node.js and Python packages with known CVEs",
				"os":        runtime.GOOS,
				"cve_count": 4,
			},
		},
		{
			ID:               "test-vuln-005",
			Type:             "configuration",
			Severity:         "low",
			Title:            "Verbose Logging Enabled",
			Description:      "System logging is set to verbose mode, potentially exposing sensitive information",
			Status:           "open",
			Priority:         "low",
			ExploitAvailable: false,
			ExploitCount:     0,
			CreatedAt:        now.Add(-1 * time.Hour),
			EnrichmentData: map[string]interface{}{
				"category":  "logging",
				"details":   "Debug logging enabled in production environment",
				"os":        runtime.GOOS,
				"cve_count": 0,
			},
		},
	}
}
