package scanner

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ConfigAuditor checks for configuration errors and security misconfigurations
type ConfigAuditor struct{}

// NewConfigAuditor creates a new configuration auditor
func NewConfigAuditor() *ConfigAuditor {
	return &ConfigAuditor{}
}

// AuditConfiguration performs configuration auditing on discovered hosts
func (ca *ConfigAuditor) AuditConfiguration(host string, ports []int, services map[int]string, banners map[int]string, credentials map[string]string) []NetworkFinding {
	var findings []NetworkFinding

	// Check for default credentials
	findings = append(findings, ca.checkDefaultCredentials(host, ports, services, credentials)...)

	// Check for insecure protocols
	findings = append(findings, ca.checkInsecureProtocols(host, ports, services)...)

	// Check for open management interfaces without auth
	findings = append(findings, ca.checkOpenManagementInterfaces(host, ports, services)...)

	// Check for unnecessary open ports
	findings = append(findings, ca.checkUnnecessaryPorts(host, ports, services)...)

	// Check for weak SSL/TLS (this would typically require deeper scanning)
	findings = append(findings, ca.checkWeakEncryption(host, ports, services)...)

	return findings
}

// checkDefaultCredentials checks for common default credentials
func (ca *ConfigAuditor) checkDefaultCredentials(host string, ports []int, services map[int]string, credentials map[string]string) []NetworkFinding {
	var findings []NetworkFinding

	// Common default credential patterns
	defaultCreds := map[string][]string{
		"admin":   {"admin", "password", "1234", "default"},
		"root":    {"root", "toor", "password"},
		"cisco":   {"cisco", "cisco"},
		"user":    {"user", "user"},
		"guest":   {"guest", "guest"},
		"service": {"service", "service"},
		"support": {"support", "support"},
		"test":    {"test", "test"},
	}

	// Check if any default credentials are being used
	for username, passwords := range defaultCreds {
		if providedCreds, exists := credentials[username]; exists {
			for _, defaultPass := range passwords {
				if providedCreds == defaultPass {
					findings = append(findings, NetworkFinding{
						ID:           uuid.New(),
						FindingType:  "config",
						Severity:     "high",
						Host:         host,
						Port:         0,
						Protocol:     "tcp",
						ServiceName:  "authentication",
						Description:  fmt.Sprintf("Default credentials detected: %s/%s", username, defaultPass),
						Remediation:  "Change default credentials immediately. Use strong, unique passwords.",
						DiscoveredAt: time.Now(),
						Status:       "open",
						Metadata: map[string]interface{}{
							"username": username,
							"password": "***",
						},
					})
					break
				}
			}
		}
	}

	return findings
}

// checkInsecureProtocols checks for insecure protocol usage
func (ca *ConfigAuditor) checkInsecureProtocols(host string, ports []int, services map[int]string) []NetworkFinding {
	var findings []NetworkFinding

	// Insecure protocols and their ports
	insecureProtocols := map[int]string{
		23:   "telnet",  // Unencrypted
		21:   "ftp",     // Often unencrypted
		161:  "snmp-v1", // SNMP v1/v2 without encryption
		162:  "snmp-v2", // SNMP v2 without encryption
		1433: "mssql",   // If unencrypted
		3306: "mysql",   // If unencrypted
	}

	for _, port := range ports {
		if protocol, isInsecure := insecureProtocols[port]; isInsecure {
			severity := "medium"
			if port == 23 { // Telnet is particularly insecure
				severity = "high"
			}

			serviceName := services[port]
			if serviceName == "" {
				serviceName = protocol
			}

			findings = append(findings, NetworkFinding{
				ID:           uuid.New(),
				FindingType:  "config",
				Severity:     severity,
				Host:         host,
				Port:         port,
				Protocol:     "tcp",
				ServiceName:  serviceName,
				Description:  fmt.Sprintf("Insecure protocol detected: %s on port %d. Data transmitted may be unencrypted.", protocol, port),
				Remediation:  fmt.Sprintf("Disable %s or use encrypted alternatives (SSH instead of Telnet, SFTP/FTPS instead of FTP, SNMPv3 instead of SNMPv1/v2).", protocol),
				DiscoveredAt: time.Now(),
				Status:       "open",
				Metadata: map[string]interface{}{
					"protocol": protocol,
					"risk":     "unencrypted_communication",
				},
			})
		}
	}

	return findings
}

// checkOpenManagementInterfaces checks for open management interfaces without authentication
func (ca *ConfigAuditor) checkOpenManagementInterfaces(host string, ports []int, services map[int]string) []NetworkFinding {
	var findings []NetworkFinding

	// Management interface ports that should require authentication
	managementPorts := map[int]string{
		22:   "SSH",
		23:   "Telnet",
		80:   "HTTP",
		443:  "HTTPS",
		161:  "SNMP",
		3389: "RDP",
		5985: "WinRM",
		5986: "WinRM-HTTPS",
	}

	for _, port := range ports {
		if mgmtName, isManagement := managementPorts[port]; isManagement {
			// Check if service is accessible without authentication
			// This is a simplified check - in reality, we'd need to test authentication
			serviceName := services[port]
			if serviceName == "" {
				serviceName = mgmtName
			}

			// For HTTP/HTTPS, check if it's a management interface
			if port == 80 || port == 443 {
				// This would require deeper scanning to determine if auth is required
				// For now, we'll flag it as a potential issue
				findings = append(findings, NetworkFinding{
					ID:           uuid.New(),
					FindingType:  "config",
					Severity:     "medium",
					Host:         host,
					Port:         port,
					Protocol:     "tcp",
					ServiceName:  serviceName,
					Description:  fmt.Sprintf("Management interface (%s) accessible on port %d. Verify authentication is required.", mgmtName, port),
					Remediation:  "Ensure all management interfaces require strong authentication. Use multi-factor authentication when possible.",
					DiscoveredAt: time.Now(),
					Status:       "open",
					Metadata: map[string]interface{}{
						"interface_type": mgmtName,
						"requires_auth":  "unknown",
					},
				})
			} else if port == 23 { // Telnet is always insecure
				findings = append(findings, NetworkFinding{
					ID:           uuid.New(),
					FindingType:  "config",
					Severity:     "high",
					Host:         host,
					Port:         port,
					Protocol:     "tcp",
					ServiceName:  serviceName,
					Description:  fmt.Sprintf("Unencrypted management interface (Telnet) detected on port %d.", port),
					Remediation:  "Disable Telnet and use SSH instead for secure remote management.",
					DiscoveredAt: time.Now(),
					Status:       "open",
					Metadata: map[string]interface{}{
						"interface_type": "Telnet",
						"encrypted":      false,
					},
				})
			}
		}
	}

	return findings
}

// checkUnnecessaryPorts identifies potentially unnecessary open ports
func (ca *ConfigAuditor) checkUnnecessaryPorts(host string, ports []int, services map[int]string) []NetworkFinding {
	var findings []NetworkFinding

	// Ports that are often unnecessary or should be restricted
	unnecessaryPorts := map[int]string{
		135:   "RPC Endpoint Mapper",
		139:   "NetBIOS Session Service",
		445:   "SMB (if not needed)",
		1433:  "MSSQL (if not needed)",
		3306:  "MySQL (if not needed)",
		5432:  "PostgreSQL (if not needed)",
		6379:  "Redis (if not needed)",
		27017: "MongoDB (if not needed)",
	}

	for _, port := range ports {
		if reason, isUnnecessary := unnecessaryPorts[port]; isUnnecessary {
			serviceName := services[port]
			if serviceName == "" {
				serviceName = reason
			}

			findings = append(findings, NetworkFinding{
				ID:           uuid.New(),
				FindingType:  "config",
				Severity:     "low",
				Host:         host,
				Port:         port,
				Protocol:     "tcp",
				ServiceName:  serviceName,
				Description:  fmt.Sprintf("Potentially unnecessary port %d (%s) is open. Verify if this service is required.", port, reason),
				Remediation:  fmt.Sprintf("Review if port %d is necessary. If not, disable the service or restrict access using firewall rules.", port),
				DiscoveredAt: time.Now(),
				Status:       "open",
				Metadata: map[string]interface{}{
					"reason": reason,
				},
			})
		}
	}

	return findings
}

// checkWeakEncryption checks for weak SSL/TLS configurations
func (ca *ConfigAuditor) checkWeakEncryption(host string, ports []int, services map[int]string) []NetworkFinding {
	var findings []NetworkFinding

	// HTTPS ports that should use strong encryption
	httpsPorts := []int{443, 8443, 9443}

	for _, port := range ports {
		for _, httpsPort := range httpsPorts {
			if port == httpsPort {
				serviceName := services[port]
				if serviceName == "" {
					serviceName = "HTTPS"
				}

				// This is a placeholder - actual SSL/TLS checking would require
				// deeper scanning to check cipher suites, TLS versions, etc.
				// For now, we'll flag it as something to verify
				findings = append(findings, NetworkFinding{
					ID:           uuid.New(),
					FindingType:  "config",
					Severity:     "medium",
					Host:         host,
					Port:         port,
					Protocol:     "tcp",
					ServiceName:  serviceName,
					Description:  fmt.Sprintf("HTTPS service detected on port %d. Verify SSL/TLS configuration uses strong ciphers and TLS 1.2+.", port),
					Remediation:  "Ensure TLS 1.2 or higher is used. Disable weak cipher suites. Use strong certificate key sizes (2048+ bits).",
					DiscoveredAt: time.Now(),
					Status:       "open",
					Metadata: map[string]interface{}{
						"encryption_type":     "TLS/SSL",
						"verification_needed": true,
					},
				})
				break
			}
		}
	}

	return findings
}

// CheckMissingPatches identifies potential missing security patches
// This is a simplified version - real implementation would require CVE matching
func (ca *ConfigAuditor) CheckMissingPatches(host string, osInfo string, serviceVersions map[string]string) []NetworkFinding {
	var findings []NetworkFinding

	// This is a placeholder - actual patch checking would require:
	// 1. OS version detection
	// 2. Service version detection
	// 3. CVE database matching
	// 4. Patch availability checking

	if osInfo != "" {
		osLower := strings.ToLower(osInfo)
		// Check for known end-of-life or unsupported OS versions
		if strings.Contains(osLower, "windows xp") ||
			strings.Contains(osLower, "windows 7") ||
			strings.Contains(osLower, "windows server 2003") {
			findings = append(findings, NetworkFinding{
				ID:           uuid.New(),
				FindingType:  "config",
				Severity:     "critical",
				Host:         host,
				Port:         0,
				Protocol:     "tcp",
				ServiceName:  "operating_system",
				Description:  fmt.Sprintf("Unsupported or end-of-life operating system detected: %s", osInfo),
				Remediation:  "Upgrade to a supported operating system version that receives security updates.",
				DiscoveredAt: time.Now(),
				Status:       "open",
				Metadata: map[string]interface{}{
					"os":  osInfo,
					"eol": true,
				},
			})
		}
	}

	return findings
}
