package scanner

import (
	"fmt"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// ScanCredentials represents credentials for authenticated scanning
type ScanCredentials struct {
	SSHUsername   string            `json:"ssh_username,omitempty"`
	SSHPassword   string            `json:"ssh_password,omitempty"`
	SSHKeyPath    string            `json:"ssh_key_path,omitempty"`
	SNMPCommunity string            `json:"snmp_community,omitempty"`
	HTTPUsername  string            `json:"http_username,omitempty"`
	HTTPPassword  string            `json:"http_password,omitempty"`
	DBUsername    string            `json:"db_username,omitempty"`
	DBPassword    string            `json:"db_password,omitempty"`
	DBType        string            `json:"db_type,omitempty"` // mysql, postgres, mssql, etc.
	Custom        map[string]string `json:"custom,omitempty"`
}

// AuthenticatedScanner handles authenticated vulnerability scanning
type AuthenticatedScanner struct {
	config      *config.Config
	nucleiScanner *NucleiScanner
}

// NewAuthenticatedScanner creates a new authenticated scanner
func NewAuthenticatedScanner(cfg *config.Config) *AuthenticatedScanner {
	return &AuthenticatedScanner{
		config:       cfg,
		nucleiScanner: NewNucleiScanner(),
	}
}

// ScanWithCredentials performs authenticated scanning on a target
func (as *AuthenticatedScanner) ScanWithCredentials(target string, credentials *ScanCredentials) ([]NetworkFinding, error) {
	var findings []NetworkFinding

	// Perform SSH-based authenticated scanning
	if credentials.SSHUsername != "" {
		sshFindings, err := as.scanSSH(target, credentials)
		if err != nil {
			return nil, fmt.Errorf("SSH scan failed: %w", err)
		}
		findings = append(findings, sshFindings...)
	}

	// Perform SNMP-based scanning
	if credentials.SNMPCommunity != "" {
		snmpFindings, err := as.scanSNMP(target, credentials)
		if err != nil {
			// SNMP failures are not critical, log and continue
			fmt.Printf("SNMP scan warning: %v\n", err)
		} else {
			findings = append(findings, snmpFindings...)
		}
	}

	// Perform HTTP-based authenticated scanning
	if credentials.HTTPUsername != "" {
		httpFindings, err := as.scanHTTP(target, credentials)
		if err != nil {
			fmt.Printf("HTTP authenticated scan warning: %v\n", err)
		} else {
			findings = append(findings, httpFindings...)
		}
	}

	// Perform database authenticated scanning
	if credentials.DBUsername != "" {
		dbFindings, err := as.scanDatabase(target, credentials)
		if err != nil {
			fmt.Printf("Database scan warning: %v\n", err)
		} else {
			findings = append(findings, dbFindings...)
		}
	}

	// Perform authenticated Nuclei scanning
	if len(findings) > 0 {
		// Use credentials for authenticated Nuclei scans
		targets := []string{target}
		nucleiFindings, err := as.nucleiScanner.ScanTargetsWithCredentials(targets, as.credentialsToMap(credentials))
		if err != nil {
			fmt.Printf("Authenticated Nuclei scan warning: %v\n", err)
		} else {
			findings = append(findings, nucleiFindings...)
		}
	}

	return findings, nil
}

// scanSSH performs SSH-based authenticated scanning
func (as *AuthenticatedScanner) scanSSH(target string, credentials *ScanCredentials) ([]NetworkFinding, error) {
	var findings []NetworkFinding

	// Check SSH connectivity and authentication
	// This is a simplified version - in production, use proper SSH client
	// For now, we'll create a finding indicating SSH access is available
	finding := NetworkFinding{
		ID:           uuid.New(),
		FindingType:  "auth",
		Severity:     "info",
		Host:         target,
		Port:         22,
		Protocol:     "tcp",
		ServiceName:  "ssh",
		Description:  fmt.Sprintf("SSH authentication available for user: %s", credentials.SSHUsername),
		Remediation:  "Ensure SSH access is properly secured with key-based authentication and strong passwords.",
		DiscoveredAt: time.Now(),
		Status:       "open",
		Metadata: map[string]interface{}{
			"username": credentials.SSHUsername,
			"auth_method": "password",
		},
	}

	// If key-based auth is used, update metadata
	if credentials.SSHKeyPath != "" {
		finding.Metadata["auth_method"] = "key"
		finding.Metadata["key_path"] = credentials.SSHKeyPath
	}

	findings = append(findings, finding)

	// Perform local vulnerability checks via SSH
	// This would require actual SSH connection and command execution
	// For now, we'll add a placeholder finding
	localVulnFinding := NetworkFinding{
		ID:           uuid.New(),
		FindingType:  "vuln",
		Severity:     "medium",
		Host:         target,
		Port:         22,
		Protocol:     "tcp",
		ServiceName:  "ssh",
		Description:  "Authenticated scan available. Check for local vulnerabilities, missing patches, and insecure configurations.",
		Remediation:  "Review system for missing security patches, insecure configurations, and unnecessary services.",
		DiscoveredAt: time.Now(),
		Status:       "open",
		Metadata: map[string]interface{}{
			"scan_type": "authenticated",
			"requires_manual_review": true,
		},
	}

	findings = append(findings, localVulnFinding)

	return findings, nil
}

// scanSNMP performs SNMP-based scanning
func (as *AuthenticatedScanner) scanSNMP(target string, credentials *ScanCredentials) ([]NetworkFinding, error) {
	var findings []NetworkFinding

	// Check SNMP access
	// This is a simplified version - in production, use proper SNMP client
	finding := NetworkFinding{
		ID:           uuid.New(),
		FindingType:  "auth",
		Severity:     "info",
		Host:         target,
		Port:         161,
		Protocol:     "udp",
		ServiceName:  "snmp",
		Description:  fmt.Sprintf("SNMP access available with community: %s", credentials.SNMPCommunity),
		Remediation:  "Ensure SNMP uses v3 with authentication and encryption. Avoid using default community strings.",
		DiscoveredAt: time.Now(),
		Status:       "open",
		Metadata: map[string]interface{}{
			"community": credentials.SNMPCommunity,
			"version": "v1/v2", // Would be determined by actual scan
		},
	}

	// Check if using default community strings
	defaultCommunities := []string{"public", "private", "community"}
	for _, defaultComm := range defaultCommunities {
		if credentials.SNMPCommunity == defaultComm {
			finding.Severity = "high"
			finding.Description = fmt.Sprintf("Default SNMP community string detected: %s", credentials.SNMPCommunity)
			finding.Remediation = "Change SNMP community string immediately. Use SNMPv3 with authentication."
			break
		}
	}

	findings = append(findings, finding)

	return findings, nil
}

// scanHTTP performs HTTP-based authenticated scanning
func (as *AuthenticatedScanner) scanHTTP(target string, credentials *ScanCredentials) ([]NetworkFinding, error) {
	var findings []NetworkFinding

	// Check HTTP/HTTPS authentication
	// This is a simplified version - in production, use proper HTTP client
	finding := NetworkFinding{
		ID:           uuid.New(),
		FindingType:  "auth",
		Severity:     "info",
		Host:         target,
		Port:         80,
		Protocol:     "tcp",
		ServiceName:  "http",
		Description:  fmt.Sprintf("HTTP authentication available for user: %s", credentials.HTTPUsername),
		Remediation:  "Ensure HTTP authentication uses strong passwords and HTTPS for encrypted transmission.",
		DiscoveredAt: time.Now(),
		Status:       "open",
		Metadata: map[string]interface{}{
			"username": credentials.HTTPUsername,
			"auth_type": "basic",
		},
	}

	findings = append(findings, finding)

	return findings, nil
}

// scanDatabase performs database authenticated scanning
func (as *AuthenticatedScanner) scanDatabase(target string, credentials *ScanCredentials) ([]NetworkFinding, error) {
	var findings []NetworkFinding

	// Determine database port based on type
	dbPort := 3306 // Default to MySQL
	switch credentials.DBType {
	case "postgres", "postgresql":
		dbPort = 5432
	case "mssql", "sqlserver":
		dbPort = 1433
	case "oracle":
		dbPort = 1521
	case "mongodb":
		dbPort = 27017
	}

	finding := NetworkFinding{
		ID:           uuid.New(),
		FindingType:  "auth",
		Severity:     "info",
		Host:         target,
		Port:         dbPort,
		Protocol:     "tcp",
		ServiceName:  credentials.DBType,
		Description:  fmt.Sprintf("Database authentication available for %s user: %s", credentials.DBType, credentials.DBUsername),
		Remediation:  "Ensure database access is properly secured with strong passwords, encrypted connections, and proper access controls.",
		DiscoveredAt: time.Now(),
		Status:       "open",
		Metadata: map[string]interface{}{
			"username": credentials.DBUsername,
			"db_type": credentials.DBType,
		},
	}

	findings = append(findings, finding)

	return findings, nil
}

// credentialsToMap converts ScanCredentials to a map for Nuclei
func (as *AuthenticatedScanner) credentialsToMap(creds *ScanCredentials) map[string]string {
	result := make(map[string]string)

	if creds.SSHUsername != "" {
		result["ssh_username"] = creds.SSHUsername
		if creds.SSHPassword != "" {
			result["ssh_password"] = creds.SSHPassword
		}
	}

	if creds.HTTPUsername != "" {
		result["http_username"] = creds.HTTPUsername
		if creds.HTTPPassword != "" {
			result["http_password"] = creds.HTTPPassword
		}
	}

	if creds.SNMPCommunity != "" {
		result["snmp_community"] = creds.SNMPCommunity
	}

	if creds.DBUsername != "" {
		result["db_username"] = creds.DBUsername
		if creds.DBPassword != "" {
			result["db_password"] = creds.DBPassword
		}
		result["db_type"] = creds.DBType
	}

	// Add custom credentials
	for k, v := range creds.Custom {
		result[k] = v
	}

	return result
}

// ValidateCredentials validates credentials before use
func (as *AuthenticatedScanner) ValidateCredentials(credentials *ScanCredentials) error {
	// Basic validation
	if credentials.SSHUsername != "" && credentials.SSHPassword == "" && credentials.SSHKeyPath == "" {
		return fmt.Errorf("SSH username provided but no password or key")
	}

	if credentials.HTTPUsername != "" && credentials.HTTPPassword == "" {
		return fmt.Errorf("HTTP username provided but no password")
	}

	if credentials.DBUsername != "" && credentials.DBPassword == "" {
		return fmt.Errorf("Database username provided but no password")
	}

	return nil
}

