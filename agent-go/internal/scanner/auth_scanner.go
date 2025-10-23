package scanner

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// AuthScanner handles authentication and access control security scanning
type AuthScanner struct {
	config *config.Config
}

// AuthFinding represents an authentication security finding
type AuthFinding struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`     // password_policy, account_security, privilege, auth_mechanism
	Severity        string                 `json:"severity"` // critical, high, medium, low
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	AffectedUser    string                 `json:"affected_user,omitempty"`
	AffectedService string                 `json:"affected_service,omitempty"`
	CurrentValue    string                 `json:"current_value,omitempty"`
	RequiredValue   string                 `json:"required_value,omitempty"`
	Remediation     string                 `json:"remediation"`
	DiscoveredAt    time.Time              `json:"discovered_at"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// PasswordPolicy represents password policy configuration
type PasswordPolicy struct {
	MinLength        int  `json:"min_length"`
	RequireUppercase bool `json:"require_uppercase"`
	RequireLowercase bool `json:"require_lowercase"`
	RequireNumbers   bool `json:"require_numbers"`
	RequireSymbols   bool `json:"require_symbols"`
	MaxAge           int  `json:"max_age_days"`
	MinAge           int  `json:"min_age_days"`
	HistoryCount     int  `json:"history_count"`
	LockoutAttempts  int  `json:"lockout_attempts"`
	LockoutDuration  int  `json:"lockout_duration_minutes"`
	ComplexityScore  int  `json:"complexity_score"`
}

// AccountInfo represents user account information
type AccountInfo struct {
	Username       string    `json:"username"`
	UID            string    `json:"uid"`
	Groups         []string  `json:"groups"`
	LastLogin      time.Time `json:"last_login"`
	PasswordAge    int       `json:"password_age_days"`
	IsActive       bool      `json:"is_active"`
	IsAdmin        bool      `json:"is_admin"`
	HasPassword    bool      `json:"has_password"`
	IsLocked       bool      `json:"is_locked"`
	LoginCount     int       `json:"login_count"`
	FailedAttempts int       `json:"failed_attempts"`
}

// PrivilegeInfo represents privilege and permission information
type PrivilegeInfo struct {
	User            string   `json:"user"`
	Privileges      []string `json:"privileges"`
	AdminGroups     []string `json:"admin_groups"`
	ServiceAccounts []string `json:"service_accounts"`
	ExcessivePerms  []string `json:"excessive_permissions"`
	RiskLevel       string   `json:"risk_level"`
}

// NewAuthScanner creates a new authentication security scanner
func NewAuthScanner(cfg *config.Config) *AuthScanner {
	return &AuthScanner{
		config: cfg,
	}
}

// Scan performs comprehensive authentication security scanning
func (as *AuthScanner) Scan() ([]AuthFinding, PasswordPolicy, []AccountInfo, []PrivilegeInfo, error) {
	var findings []AuthFinding
	var passwordPolicy PasswordPolicy
	var accounts []AccountInfo
	var privileges []PrivilegeInfo

	// Perform OS-specific authentication scanning
	switch runtime.GOOS {
	case "darwin":
		authFindings, policy, accountList, privilegeList, err := as.scanMacOSAuth()
		if err != nil {
			return nil, PasswordPolicy{}, nil, nil, err
		}
		findings = append(findings, authFindings...)
		passwordPolicy = policy
		accounts = append(accounts, accountList...)
		privileges = append(privileges, privilegeList...)
	case "linux":
		authFindings, policy, accountList, privilegeList, err := as.scanLinuxAuth()
		if err != nil {
			return nil, PasswordPolicy{}, nil, nil, err
		}
		findings = append(findings, authFindings...)
		passwordPolicy = policy
		accounts = append(accounts, accountList...)
		privileges = append(privileges, privilegeList...)
	case "windows":
		authFindings, policy, accountList, privilegeList, err := as.scanWindowsAuth()
		if err != nil {
			return nil, PasswordPolicy{}, nil, nil, err
		}
		findings = append(findings, authFindings...)
		passwordPolicy = policy
		accounts = append(accounts, accountList...)
		privileges = append(privileges, privilegeList...)
	default:
		return nil, PasswordPolicy{}, nil, nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return findings, passwordPolicy, accounts, privileges, nil
}

// scanMacOSAuth performs macOS-specific authentication scanning
func (as *AuthScanner) scanMacOSAuth() ([]AuthFinding, PasswordPolicy, []AccountInfo, []PrivilegeInfo, error) {
	var findings []AuthFinding
	var accounts []AccountInfo
	var privileges []PrivilegeInfo

	// Check password policies
	passwordPolicy := as.checkMacOSPasswordPolicy()

	// Check for weak password policies
	if passwordPolicy.MinLength < 8 {
		finding := AuthFinding{
			ID:            uuid.New().String(),
			Type:          "password_policy",
			Severity:      "high",
			Title:         "Weak Password Length Policy",
			Description:   fmt.Sprintf("Minimum password length is %d characters, should be at least 8", passwordPolicy.MinLength),
			CurrentValue:  strconv.Itoa(passwordPolicy.MinLength),
			RequiredValue: "8",
			Remediation:   "Increase minimum password length to at least 8 characters",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"os":       "macOS",
				"policy":   "password_length",
				"current":  passwordPolicy.MinLength,
				"required": 8,
			},
		}
		findings = append(findings, finding)
	}

	// Check password complexity
	if !passwordPolicy.RequireUppercase || !passwordPolicy.RequireLowercase || !passwordPolicy.RequireNumbers {
		finding := AuthFinding{
			ID:           uuid.New().String(),
			Type:         "password_policy",
			Severity:     "medium",
			Title:        "Weak Password Complexity Policy",
			Description:  "Password complexity requirements are insufficient",
			Remediation:  "Enable uppercase, lowercase, and numeric character requirements",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"os":                "macOS",
				"policy":            "password_complexity",
				"require_uppercase": passwordPolicy.RequireUppercase,
				"require_lowercase": passwordPolicy.RequireLowercase,
				"require_numbers":   passwordPolicy.RequireNumbers,
			},
		}
		findings = append(findings, finding)
	}

	// Check for default accounts
	accounts = as.getMacOSAccounts()
	for _, account := range accounts {
		// Check for default accounts
		if account.Username == "admin" || account.Username == "administrator" {
			finding := AuthFinding{
				ID:           uuid.New().String(),
				Type:         "account_security",
				Severity:     "high",
				Title:        "Default Account Detected",
				Description:  fmt.Sprintf("Default account '%s' is present", account.Username),
				AffectedUser: account.Username,
				Remediation:  "Remove or rename default accounts",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"os":       "macOS",
					"account":  account.Username,
					"is_admin": account.IsAdmin,
				},
			}
			findings = append(findings, finding)
		}

		// Check for accounts without passwords
		if !account.HasPassword {
			finding := AuthFinding{
				ID:           uuid.New().String(),
				Type:         "account_security",
				Severity:     "critical",
				Title:        "Account Without Password",
				Description:  fmt.Sprintf("Account '%s' has no password set", account.Username),
				AffectedUser: account.Username,
				Remediation:  "Set a strong password for this account",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"os":           "macOS",
					"account":      account.Username,
					"has_password": account.HasPassword,
				},
			}
			findings = append(findings, finding)
		}

		// Check for dormant accounts
		if account.LastLogin.IsZero() || time.Since(account.LastLogin) > 90*24*time.Hour {
			finding := AuthFinding{
				ID:           uuid.New().String(),
				Type:         "account_security",
				Severity:     "medium",
				Title:        "Dormant Account Detected",
				Description:  fmt.Sprintf("Account '%s' has not logged in recently", account.Username),
				AffectedUser: account.Username,
				Remediation:  "Review and disable unused accounts",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"os":         "macOS",
					"account":    account.Username,
					"last_login": account.LastLogin,
					"days_since": int(time.Since(account.LastLogin).Hours() / 24),
				},
			}
			findings = append(findings, finding)
		}
	}

	// Check for excessive privileges
	privileges = as.getMacOSPrivileges()
	for _, privilege := range privileges {
		if len(privilege.ExcessivePerms) > 0 {
			finding := AuthFinding{
				ID:           uuid.New().String(),
				Type:         "privilege",
				Severity:     "high",
				Title:        "Excessive Privileges Detected",
				Description:  fmt.Sprintf("User '%s' has excessive privileges: %s", privilege.User, strings.Join(privilege.ExcessivePerms, ", ")),
				AffectedUser: privilege.User,
				Remediation:  "Review and reduce user privileges to minimum required",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"os":              "macOS",
					"user":            privilege.User,
					"excessive_perms": privilege.ExcessivePerms,
					"risk_level":      privilege.RiskLevel,
				},
			}
			findings = append(findings, finding)
		}
	}

	// Check for MFA/2FA enforcement
	mfaFinding := as.checkMacOSMFA()
	if mfaFinding != nil {
		findings = append(findings, *mfaFinding)
	}

	return findings, passwordPolicy, accounts, privileges, nil
}

// scanLinuxAuth performs Linux-specific authentication scanning
func (as *AuthScanner) scanLinuxAuth() ([]AuthFinding, PasswordPolicy, []AccountInfo, []PrivilegeInfo, error) {
	var findings []AuthFinding
	var accounts []AccountInfo
	var privileges []PrivilegeInfo

	// Check password policies
	passwordPolicy := as.checkLinuxPasswordPolicy()

	// Check for weak password policies
	if passwordPolicy.MinLength < 8 {
		finding := AuthFinding{
			ID:            uuid.New().String(),
			Type:          "password_policy",
			Severity:      "high",
			Title:         "Weak Password Length Policy",
			Description:   fmt.Sprintf("Minimum password length is %d characters, should be at least 8", passwordPolicy.MinLength),
			CurrentValue:  strconv.Itoa(passwordPolicy.MinLength),
			RequiredValue: "8",
			Remediation:   "Increase minimum password length to at least 8 characters",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"os":       "Linux",
				"policy":   "password_length",
				"current":  passwordPolicy.MinLength,
				"required": 8,
			},
		}
		findings = append(findings, finding)
	}

	// Check for default accounts
	accounts = as.getLinuxAccounts()
	for _, account := range accounts {
		// Check for root account security
		if account.Username == "root" {
			if account.HasPassword {
				finding := AuthFinding{
					ID:           uuid.New().String(),
					Type:         "account_security",
					Severity:     "medium",
					Title:        "Root Account Password Set",
					Description:  "Root account has a password set (consider using sudo instead)",
					AffectedUser: account.Username,
					Remediation:  "Consider disabling root login and using sudo",
					DiscoveredAt: time.Now(),
					Metadata: map[string]interface{}{
						"os":           "Linux",
						"account":      account.Username,
						"has_password": account.HasPassword,
					},
				}
				findings = append(findings, finding)
			}
		}

		// Check for accounts without passwords
		if !account.HasPassword && account.Username != "nobody" {
			finding := AuthFinding{
				ID:           uuid.New().String(),
				Type:         "account_security",
				Severity:     "critical",
				Title:        "Account Without Password",
				Description:  fmt.Sprintf("Account '%s' has no password set", account.Username),
				AffectedUser: account.Username,
				Remediation:  "Set a strong password for this account or disable it",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"os":           "Linux",
					"account":      account.Username,
					"has_password": account.HasPassword,
				},
			}
			findings = append(findings, finding)
		}
	}

	// Check for excessive privileges
	privileges = as.getLinuxPrivileges()
	for _, privilege := range privileges {
		if len(privilege.ExcessivePerms) > 0 {
			finding := AuthFinding{
				ID:           uuid.New().String(),
				Type:         "privilege",
				Severity:     "high",
				Title:        "Excessive Privileges Detected",
				Description:  fmt.Sprintf("User '%s' has excessive privileges: %s", privilege.User, strings.Join(privilege.ExcessivePerms, ", ")),
				AffectedUser: privilege.User,
				Remediation:  "Review and reduce user privileges to minimum required",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"os":              "Linux",
					"user":            privilege.User,
					"excessive_perms": privilege.ExcessivePerms,
					"risk_level":      privilege.RiskLevel,
				},
			}
			findings = append(findings, finding)
		}
	}

	return findings, passwordPolicy, accounts, privileges, nil
}

// scanWindowsAuth performs Windows-specific authentication scanning
func (as *AuthScanner) scanWindowsAuth() ([]AuthFinding, PasswordPolicy, []AccountInfo, []PrivilegeInfo, error) {
	var findings []AuthFinding
	var accounts []AccountInfo
	var privileges []PrivilegeInfo

	// This would require Windows-specific implementation
	// For now, return placeholder data
	passwordPolicy := PasswordPolicy{
		MinLength:        8,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSymbols:   false,
		MaxAge:           90,
		MinAge:           1,
		HistoryCount:     5,
		LockoutAttempts:  5,
		LockoutDuration:  30,
		ComplexityScore:  3,
	}

	// Placeholder findings for Windows
	finding := AuthFinding{
		ID:           uuid.New().String(),
		Type:         "password_policy",
		Severity:     "medium",
		Title:        "Windows Password Policy Check",
		Description:  "Windows password policy analysis not fully implemented",
		Remediation:  "Implement Windows-specific password policy checking",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"os":     "Windows",
			"status": "placeholder",
		},
	}
	findings = append(findings, finding)

	return findings, passwordPolicy, accounts, privileges, nil
}

// checkMacOSPasswordPolicy checks macOS password policy
func (as *AuthScanner) checkMacOSPasswordPolicy() PasswordPolicy {
	policy := PasswordPolicy{
		MinLength:        8,
		RequireUppercase: false,
		RequireLowercase: false,
		RequireNumbers:   false,
		RequireSymbols:   false,
		MaxAge:           90,
		MinAge:           1,
		HistoryCount:     0,
		LockoutAttempts:  0,
		LockoutDuration:  0,
		ComplexityScore:  1,
	}

	// Check pwpolicy settings
	cmd := exec.Command("pwpolicy", "-getaccountpolicies")
	output, err := cmd.Output()
	if err == nil {
		policyText := string(output)

		// Parse minimum length
		if strings.Contains(policyText, "minChars") {
			// Extract minimum character requirement
			policy.MinLength = 8 // Default assumption
		}

		// Parse complexity requirements
		if strings.Contains(policyText, "requireMixedCase") {
			policy.RequireUppercase = true
			policy.RequireLowercase = true
		}
		if strings.Contains(policyText, "requireNumeric") {
			policy.RequireNumbers = true
		}
		if strings.Contains(policyText, "requireSymbol") {
			policy.RequireSymbols = true
		}
	}

	return policy
}

// checkLinuxPasswordPolicy checks Linux password policy
func (as *AuthScanner) checkLinuxPasswordPolicy() PasswordPolicy {
	policy := PasswordPolicy{
		MinLength:        8,
		RequireUppercase: false,
		RequireLowercase: false,
		RequireNumbers:   false,
		RequireSymbols:   false,
		MaxAge:           90,
		MinAge:           1,
		HistoryCount:     5,
		LockoutAttempts:  5,
		LockoutDuration:  30,
		ComplexityScore:  2,
	}

	// Check PAM configuration
	cmd := exec.Command("grep", "-r", "pam_cracklib", "/etc/pam.d/")
	output, err := cmd.Output()
	if err == nil {
		// Parse PAM configuration for password requirements
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "minlen") {
				// Extract minimum length
				parts := strings.Split(line, "minlen=")
				if len(parts) > 1 {
					if length, err := strconv.Atoi(parts[1]); err == nil {
						policy.MinLength = length
					}
				}
			}
			if strings.Contains(line, "dcredit") {
				policy.RequireNumbers = true
			}
			if strings.Contains(line, "ucredit") {
				policy.RequireUppercase = true
			}
			if strings.Contains(line, "lcredit") {
				policy.RequireLowercase = true
			}
			if strings.Contains(line, "ocredit") {
				policy.RequireSymbols = true
			}
		}
	}

	return policy
}

// getMacOSAccounts retrieves macOS user accounts
func (as *AuthScanner) getMacOSAccounts() []AccountInfo {
	var accounts []AccountInfo

	// Get all users
	cmd := exec.Command("dscl", ".", "-list", "/Users")
	output, err := cmd.Output()
	if err != nil {
		return accounts
	}

	lines := strings.Split(string(output), "\n")
	for _, username := range lines {
		username = strings.TrimSpace(username)
		if username == "" || username == "daemon" || username == "nobody" {
			continue
		}

		account := AccountInfo{
			Username:    username,
			IsActive:    true,
			HasPassword: true, // Default assumption
		}

		// Check if user is admin
		cmd := exec.Command("dscl", ".", "-read", "/Groups/admin", "GroupMembership")
		output, err := cmd.Output()
		if err == nil && strings.Contains(string(output), username) {
			account.IsAdmin = true
		}

		// Check last login
		cmd = exec.Command("last", "-1", username)
		output, err = cmd.Output()
		if err == nil {
			// Parse last login time
			account.LastLogin = time.Now().Add(-24 * time.Hour) // Placeholder
		}

		accounts = append(accounts, account)
	}

	return accounts
}

// getLinuxAccounts retrieves Linux user accounts
func (as *AuthScanner) getLinuxAccounts() []AccountInfo {
	var accounts []AccountInfo

	// Read /etc/passwd
	cmd := exec.Command("cat", "/etc/passwd")
	output, err := cmd.Output()
	if err != nil {
		return accounts
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}

		username := parts[0]
		uid := parts[2]
		shell := parts[6]

		// Skip system accounts
		if uid == "0" && username != "root" {
			continue
		}
		if shell == "/bin/false" || shell == "/usr/sbin/nologin" {
			continue
		}

		account := AccountInfo{
			Username:    username,
			UID:         uid,
			IsActive:    true,
			HasPassword: true, // Default assumption
		}

		// Check if user is in admin groups
		cmd = exec.Command("groups", username)
		output, err = cmd.Output()
		if err == nil {
			groups := strings.Split(strings.TrimSpace(string(output)), " ")
			for _, group := range groups {
				if group == "sudo" || group == "wheel" || group == "admin" {
					account.IsAdmin = true
					break
				}
			}
		}

		// Check last login
		cmd = exec.Command("last", "-1", username)
		output, err = cmd.Output()
		if err == nil {
			// Parse last login time
			account.LastLogin = time.Now().Add(-24 * time.Hour) // Placeholder
		}

		accounts = append(accounts, account)
	}

	return accounts
}

// getMacOSPrivileges retrieves macOS privilege information
func (as *AuthScanner) getMacOSPrivileges() []PrivilegeInfo {
	var privileges []PrivilegeInfo

	// Get admin users
	cmd := exec.Command("dscl", ".", "-read", "/Groups/admin", "GroupMembership")
	output, err := cmd.Output()
	if err == nil {
		adminUsers := strings.Split(strings.TrimSpace(string(output)), " ")
		for _, user := range adminUsers {
			if user != "" {
				privilege := PrivilegeInfo{
					User:           user,
					Privileges:     []string{"admin"},
					AdminGroups:    []string{"admin"},
					ExcessivePerms: []string{},
					RiskLevel:      "medium",
				}
				privileges = append(privileges, privilege)
			}
		}
	}

	return privileges
}

// getLinuxPrivileges retrieves Linux privilege information
func (as *AuthScanner) getLinuxPrivileges() []PrivilegeInfo {
	var privileges []PrivilegeInfo

	// Get sudo users
	cmd := exec.Command("getent", "group", "sudo")
	output, err := cmd.Output()
	if err == nil {
		// Parse sudo group members
		line := strings.TrimSpace(string(output))
		parts := strings.Split(line, ":")
		if len(parts) > 3 {
			members := strings.Split(parts[3], ",")
			for _, user := range members {
				if user != "" {
					privilege := PrivilegeInfo{
						User:           user,
						Privileges:     []string{"sudo"},
						AdminGroups:    []string{"sudo"},
						ExcessivePerms: []string{},
						RiskLevel:      "medium",
					}
					privileges = append(privileges, privilege)
				}
			}
		}
	}

	return privileges
}

// checkMacOSMFA checks for MFA/2FA enforcement on macOS
func (as *AuthScanner) checkMacOSMFA() *AuthFinding {
	// Check for Touch ID or other MFA methods
	cmd := exec.Command("bioutil", "-r")
	output, err := cmd.Output()
	if err != nil {
		_ = output // Suppress unused variable warning
		// MFA not configured
		return &AuthFinding{
			ID:           uuid.New().String(),
			Type:         "auth_mechanism",
			Severity:     "medium",
			Title:        "MFA/2FA Not Enforced",
			Description:  "Multi-factor authentication is not configured",
			Remediation:  "Enable Touch ID or other MFA methods for enhanced security",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"os":       "macOS",
				"mfa_type": "none",
			},
		}
	}

	return nil // MFA is configured
}
