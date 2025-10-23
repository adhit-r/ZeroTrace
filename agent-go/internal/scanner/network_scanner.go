package scanner

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// NetworkScanner handles network security scanning
type NetworkScanner struct {
	config *config.Config
}

// NetworkScanResult represents the result of a network scan
type NetworkScanResult struct {
	ID              uuid.UUID              `json:"id"`
	AgentID         uuid.UUID              `json:"agent_id"`
	CompanyID       uuid.UUID              `json:"company_id"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         time.Time              `json:"end_time"`
	Status          string                 `json:"status"`
	NetworkFindings []NetworkFinding       `json:"network_findings"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// NetworkFinding represents a network security finding
type NetworkFinding struct {
	ID             uuid.UUID `json:"id"`
	FindingType    string    `json:"finding_type"` // port, service, ssl, protocol
	Severity       string    `json:"severity"`     // critical, high, medium, low, info
	Host           string    `json:"host"`
	Port           int       `json:"port"`
	Protocol       string    `json:"protocol"` // tcp, udp
	ServiceName    string    `json:"service_name"`
	ServiceVersion string    `json:"service_version"`
	Banner         string    `json:"banner"`
	Description    string    `json:"description"`
	Remediation    string    `json:"remediation"`
	DiscoveredAt   time.Time `json:"discovered_at"`
	Status         string    `json:"status"` // open, filtered, closed
}

// SSLFinding represents SSL/TLS security findings
type SSLFinding struct {
	Host          string    `json:"host"`
	Port          int       `json:"port"`
	Protocol      string    `json:"protocol"`
	Certificate   string    `json:"certificate"`
	Issuer        string    `json:"issuer"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidTo       time.Time `json:"valid_to"`
	WeakCiphers   []string  `json:"weak_ciphers"`
	WeakProtocols []string  `json:"weak_protocols"`
	SelfSigned    bool      `json:"self_signed"`
	Expired       bool      `json:"expired"`
	ExpiresSoon   bool      `json:"expires_soon"`
	Description   string    `json:"description"`
	Remediation   string    `json:"remediation"`
	Severity      string    `json:"severity"`
}

// NewNetworkScanner creates a new network scanner instance
func NewNetworkScanner(cfg *config.Config) *NetworkScanner {
	return &NetworkScanner{
		config: cfg,
	}
}

// Scan performs a comprehensive network security scan
func (ns *NetworkScanner) Scan() (*NetworkScanResult, error) {
	startTime := time.Now()

	// Create scan result
	agentID, _ := uuid.Parse(ns.config.AgentID)
	companyID, _ := uuid.Parse(ns.config.CompanyID)
	result := &NetworkScanResult{
		ID:              uuid.New(),
		AgentID:         agentID,
		CompanyID:       companyID,
		StartTime:       startTime,
		EndTime:         time.Now(),
		Status:          "completed",
		NetworkFindings: []NetworkFinding{},
		Metadata:        make(map[string]interface{}),
	}

	// Get local network interfaces
	interfaces, err := ns.getNetworkInterfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Scan each interface
	for _, iface := range interfaces {
		findings, err := ns.scanInterface(iface)
		if err != nil {
			continue // Continue with other interfaces
		}
		result.NetworkFindings = append(result.NetworkFindings, findings...)
	}

	// Update metadata
	result.Metadata["total_findings"] = len(result.NetworkFindings)
	result.Metadata["scan_duration"] = time.Since(startTime).Seconds()
	result.Metadata["interfaces_scanned"] = len(interfaces)

	return result, nil
}

// getNetworkInterfaces returns a list of network interfaces to scan
func (ns *NetworkScanner) getNetworkInterfaces() ([]string, error) {
	var interfaces []string

	// Get all network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		// Skip loopback and inactive interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Get interface addresses
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					// IPv4 address
					interfaces = append(interfaces, ipnet.IP.String())
				}
			}
		}
	}

	return interfaces, nil
}

// scanInterface scans a specific network interface
func (ns *NetworkScanner) scanInterface(host string) ([]NetworkFinding, error) {
	var findings []NetworkFinding

	// Common ports to scan
	commonPorts := []int{
		21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 993, 995, 1433, 3389, 5432, 5900, 6379, 8080, 8443,
	}

	// Scan common ports
	portFindings := ns.scanPorts(host, commonPorts)
	findings = append(findings, portFindings...)

	// Scan for SSL/TLS services
	sslFindings := ns.scanSSL(host, portFindings)
	for _, ssl := range sslFindings {
		findings = append(findings, NetworkFinding{
			ID:             uuid.New(),
			FindingType:    "ssl",
			Severity:       ssl.Severity,
			Host:           ssl.Host,
			Port:           ssl.Port,
			Protocol:       ssl.Protocol,
			ServiceName:    "ssl/tls",
			ServiceVersion: ssl.Protocol,
			Description:    ssl.Description,
			Remediation:    ssl.Remediation,
			DiscoveredAt:   time.Now(),
			Status:         "open",
		})
	}

	return findings, nil
}

// scanPorts scans a list of ports on a host
func (ns *NetworkScanner) scanPorts(host string, ports []int) []NetworkFinding {
	var findings []NetworkFinding
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Limit concurrent scans
	semaphore := make(chan struct{}, 50)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			finding := ns.scanPort(host, p)
			if finding != nil {
				mu.Lock()
				findings = append(findings, *finding)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	return findings
}

// scanPort scans a single port
func (ns *NetworkScanner) scanPort(host string, port int) *NetworkFinding {
	address := fmt.Sprintf("%s:%d", host, port)

	// Try TCP connection
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return nil // Port is closed or filtered
	}
	defer conn.Close()

	// Port is open, try to identify service
	service, version, banner := ns.identifyService(host, port)

	severity := ns.determineSeverity(port, service)

	return &NetworkFinding{
		ID:             uuid.New(),
		FindingType:    "port",
		Severity:       severity,
		Host:           host,
		Port:           port,
		Protocol:       "tcp",
		ServiceName:    service,
		ServiceVersion: version,
		Banner:         banner,
		Description:    fmt.Sprintf("Open %s service on port %d", service, port),
		Remediation:    ns.getRemediation(port, service),
		DiscoveredAt:   time.Now(),
		Status:         "open",
	}
}

// identifyService attempts to identify the service running on a port
func (ns *NetworkScanner) identifyService(host string, port int) (string, string, string) {
	address := fmt.Sprintf("%s:%d", host, port)

	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return "unknown", "", ""
	}
	defer conn.Close()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Try to read banner
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "unknown", "", ""
	}

	banner := string(buffer[:n])
	service, version := ns.parseBanner(banner, port)

	return service, version, banner
}

// parseBanner parses service banner to identify service and version
func (ns *NetworkScanner) parseBanner(banner string, port int) (string, string) {
	banner = strings.ToLower(banner)

	// Common service patterns
	patterns := map[string][]string{
		"ssh":      {"openssh", "ssh-2.0", "ssh-1.99"},
		"http":     {"http/", "server:", "apache", "nginx", "iis"},
		"ftp":      {"220", "ftp"},
		"smtp":     {"220", "esmtp", "postfix", "sendmail"},
		"pop3":     {"+ok", "pop3"},
		"imap":     {"* ok", "imap"},
		"telnet":   {"login:", "password:"},
		"mysql":    {"mysql"},
		"postgres": {"postgresql"},
		"redis":    {"redis"},
		"rdp":      {"microsoft terminal services"},
	}

	// Check for service patterns
	for service, patterns := range patterns {
		for _, pattern := range patterns {
			if strings.Contains(banner, pattern) {
				version := ns.extractVersion(banner)
				return service, version
			}
		}
	}

	// Default based on port
	portServices := map[int]string{
		22:  "ssh",
		23:  "telnet",
		25:  "smtp",
		53:  "dns",
		80:  "http",
		110: "pop3",
		143: "imap",
		443: "https",
		993: "imaps",
		995: "pop3s",
	}

	if service, exists := portServices[port]; exists {
		return service, ""
	}

	return "unknown", ""
}

// extractVersion attempts to extract version information from banner
func (ns *NetworkScanner) extractVersion(banner string) string {
	// Look for version patterns
	versionPatterns := []string{
		"version", "v", "release", "build",
	}

	for _, pattern := range versionPatterns {
		if strings.Contains(banner, pattern) {
			// Simple version extraction
			parts := strings.Fields(banner)
			for i, part := range parts {
				if strings.Contains(part, pattern) && i+1 < len(parts) {
					return parts[i+1]
				}
			}
		}
	}

	return ""
}

// determineSeverity determines the severity of an open port
func (ns *NetworkScanner) determineSeverity(port int, service string) string {
	// Critical ports
	criticalPorts := []int{22, 23, 3389, 5900}
	for _, p := range criticalPorts {
		if port == p {
			return "critical"
		}
	}

	// High risk ports
	highRiskPorts := []int{21, 25, 135, 139, 1433, 5432, 6379}
	for _, p := range highRiskPorts {
		if port == p {
			return "high"
		}
	}

	// Medium risk ports
	mediumRiskPorts := []int{80, 110, 143, 443, 993, 995, 8080, 8443}
	for _, p := range mediumRiskPorts {
		if port == p {
			return "medium"
		}
	}

	return "low"
}

// getRemediation provides remediation advice for open ports
func (ns *NetworkScanner) getRemediation(port int, service string) string {
	remediations := map[int]string{
		21:   "Disable FTP service or use SFTP/FTPS",
		22:   "Ensure SSH is properly configured with key-based authentication",
		23:   "Disable Telnet service, use SSH instead",
		25:   "Configure SMTP server securely with authentication",
		80:   "Use HTTPS instead of HTTP",
		110:  "Use POP3S (port 995) instead of POP3",
		143:  "Use IMAPS (port 993) instead of IMAP",
		135:  "Disable RPC endpoint mapper if not needed",
		139:  "Disable NetBIOS if not needed",
		443:  "Ensure SSL/TLS is properly configured",
		1433: "Secure SQL Server with proper authentication",
		3389: "Secure RDP with Network Level Authentication",
		5432: "Secure PostgreSQL with proper authentication",
		5900: "Secure VNC with authentication",
		6379: "Secure Redis with authentication",
	}

	if remediation, exists := remediations[port]; exists {
		return remediation
	}

	return "Review if this service is necessary and secure it properly"
}

// scanSSL scans for SSL/TLS services and analyzes certificates
func (ns *NetworkScanner) scanSSL(host string, portFindings []NetworkFinding) []SSLFinding {
	var sslFindings []SSLFinding

	// Check for SSL/TLS on common ports
	sslPorts := []int{443, 993, 995, 8443, 465, 587}

	for _, finding := range portFindings {
		for _, sslPort := range sslPorts {
			if finding.Port == sslPort {
				sslFinding := ns.analyzeSSL(host, finding.Port)
				if sslFinding != nil {
					sslFindings = append(sslFindings, *sslFinding)
				}
			}
		}
	}

	return sslFindings
}

// analyzeSSL analyzes SSL/TLS configuration
func (ns *NetworkScanner) analyzeSSL(host string, port int) *SSLFinding {
	// This is a simplified SSL analysis
	// In a real implementation, you would use crypto/tls package

	address := fmt.Sprintf("%s:%d", host, port)

	// Try to establish SSL connection
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return nil
	}
	defer conn.Close()

	// Basic SSL analysis (simplified)
	sslFinding := &SSLFinding{
		Host:          host,
		Port:          port,
		Protocol:      "tls",
		SelfSigned:    false, // Would need actual certificate analysis
		Expired:       false, // Would need actual certificate analysis
		ExpiresSoon:   false, // Would need actual certificate analysis
		WeakCiphers:   []string{},
		WeakProtocols: []string{},
	}

	// Determine severity based on common SSL issues
	severity := "medium"
	description := fmt.Sprintf("SSL/TLS service detected on port %d", port)
	remediation := "Ensure SSL/TLS is properly configured with strong ciphers and valid certificates"

	sslFinding.Severity = severity
	sslFinding.Description = description
	sslFinding.Remediation = remediation

	return sslFinding
}

// GetNetworkTopology attempts to map network topology
func (ns *NetworkScanner) GetNetworkTopology() (map[string]interface{}, error) {
	topology := make(map[string]interface{})

	// Get local network information
	interfaces, err := ns.getNetworkInterfaces()
	if err != nil {
		return nil, err
	}

	topology["local_interfaces"] = interfaces
	topology["scan_timestamp"] = time.Now()

	// In a real implementation, you would:
	// - Discover other hosts on the network
	// - Map network topology
	// - Identify network devices
	// - Analyze routing tables

	return topology, nil
}
