package scanner

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/Ullaakut/nmap/v2"
	"github.com/google/uuid"
	"github.com/projectdiscovery/naabu/v2/pkg/result"
	"github.com/projectdiscovery/naabu/v2/pkg/runner"
)

// NetworkScanner handles network security scanning using Nmap, Naabu, and Nuclei
type NetworkScanner struct {
	config           *config.Config
	deviceClassifier *DeviceClassifier
	configAuditor    *ConfigAuditor
	nucleiScanner    *NucleiScanner
}

// NewNetworkScanner creates a new NetworkScanner
func NewNetworkScanner(cfg *config.Config) *NetworkScanner {
	return &NetworkScanner{
		config:           cfg,
		deviceClassifier: NewDeviceClassifier(),
		configAuditor:    NewConfigAuditor(),
		nucleiScanner:    NewNucleiScanner(),
	}
}

// Scan performs a comprehensive network scan using Nmap for device discovery,
// device classification, configuration auditing, and Nuclei for vulnerability scanning
func (ns *NetworkScanner) Scan(target string) (*NetworkScanResult, error) {
	scanID := uuid.New()
	startTime := time.Now()

	var allFindings []NetworkFinding
	var hostsWithOpenPorts []string

	// Step 1: Use Nmap for comprehensive device discovery and fingerprinting
	nmapResults, err := ns.scanWithNmap(target)
	if err != nil {
		// Fallback to Naabu if Nmap fails
		return ns.scanWithNaabu(target, scanID, startTime)
	}

	// Step 2: Process Nmap results and classify devices
	for _, host := range nmapResults {
		ports := []int{}
		services := make(map[int]string)
		banners := make(map[int]string)
		osInfo := ""
		osVersion := ""

		// Extract port, service, and banner information
		for _, port := range host.Ports {
			if port.State.State == "open" {
				ports = append(ports, int(port.ID))
				if port.Service.Name != "" {
					services[int(port.ID)] = port.Service.Name
				}
				if port.Service.Product != "" {
					serviceStr := port.Service.Product
					if port.Service.Version != "" {
						serviceStr = fmt.Sprintf("%s %s", port.Service.Product, port.Service.Version)
					}
					services[int(port.ID)] = serviceStr
				}
				// Note: Banner field may not be available in this nmap library version
				// We'll use service name/product as banner info
				if port.Service.Name != "" {
					banners[int(port.ID)] = port.Service.Name
				}
			}
		}

		// Extract OS information
		if len(host.OS.Matches) > 0 {
			osInfo = host.OS.Matches[0].Name
			// OS version may be in different fields depending on nmap library version
			if len(host.OS.Matches[0].Classes) > 0 {
				// Try to get OS type/version from class
				osClass := host.OS.Matches[0].Classes[0]
				if osClass.Type != "" {
					osVersion = osClass.Type
				}
			}
		}

		// Classify device type
		deviceType := ns.deviceClassifier.ClassifyDevice(host.Addresses[0].Addr, ports, services, osInfo, banners)

		// Create port findings
		for _, port := range ports {
			serviceName := services[port]
			if serviceName == "" {
				serviceName = "unknown"
			}

			banner := banners[port]

			finding := NetworkFinding{
				ID:             uuid.New(),
				FindingType:    "port",
				Severity:       "info",
				Host:           host.Addresses[0].Addr,
				Port:           port,
				Protocol:       "tcp",
				ServiceName:    serviceName,
				ServiceVersion: "",
				Banner:         banner,
				Description:    fmt.Sprintf("Open port %d (%s) discovered on %s", port, serviceName, deviceType),
				Remediation:    "Review if this service is necessary and secure it properly",
				DiscoveredAt:   time.Now(),
				Status:         "open",
				DeviceType:     deviceType,
				OS:             osInfo,
				OSVersion:      osVersion,
				Metadata: map[string]interface{}{
					"confidence": ns.deviceClassifier.GetDeviceConfidence(deviceType, ports, services, osInfo),
				},
			}

			allFindings = append(allFindings, finding)
			hostsWithOpenPorts = append(hostsWithOpenPorts, fmt.Sprintf("%s:%d", host.Addresses[0].Addr, port))
		}

		// Step 3: Perform configuration auditing
		credentials := make(map[string]string) // Would be populated from config if available
		configFindings := ns.configAuditor.AuditConfiguration(
			host.Addresses[0].Addr,
			ports,
			services,
			banners,
			credentials,
		)
		allFindings = append(allFindings, configFindings...)

		// Check for missing patches
		serviceVersions := make(map[string]string)
		for port, service := range services {
			serviceVersions[fmt.Sprintf("%d", port)] = service
		}
		patchFindings := ns.configAuditor.CheckMissingPatches(
			host.Addresses[0].Addr,
			osInfo,
			serviceVersions,
		)
		allFindings = append(allFindings, patchFindings...)
	}

	// Step 4: Run Nuclei vulnerability scanning on discovered hosts
	var vulnFindings []NetworkFinding
	if len(hostsWithOpenPorts) > 0 {
		// Prepare targets for Nuclei (unique hosts)
		uniqueHosts := make(map[string]bool)
		targets := []string{}
		for _, hostPort := range hostsWithOpenPorts {
			parts := strings.Split(hostPort, ":")
			if len(parts) > 0 {
				host := parts[0]
				if !uniqueHosts[host] {
					uniqueHosts[host] = true
					targets = append(targets, host)
				}
			}
		}

		// Run Nuclei scan
		nucleiFindings, err := ns.nucleiScanner.ScanTargets(targets)
		if err != nil {
			// Try CLI fallback
			nucleiFindings, err = ns.nucleiScanner.ScanUsingCLI(targets)
			if err != nil {
				fmt.Printf("Warning: Nuclei vulnerability scanning failed: %v\n", err)
			}
		}
		vulnFindings = append(vulnFindings, nucleiFindings...)
	}

	// Combine all findings
	allFindings = append(allFindings, vulnFindings...)

	result := &NetworkScanResult{
		ID:              scanID,
		AgentID:         uuid.MustParse(ns.config.AgentID),
		CompanyID:       uuid.MustParse(ns.config.CompanyID),
		StartTime:       startTime,
		EndTime:         time.Now(),
		Status:          "completed",
		NetworkFindings: allFindings,
		Metadata: map[string]interface{}{
			"total_hosts":    len(nmapResults),
			"total_findings": len(allFindings),
			"port_findings":  len(allFindings) - len(vulnFindings),
			"vuln_findings":  len(vulnFindings),
			"scan_method":    "nmap+nuclei",
		},
	}

	return result, nil
}

// scanWithNmap performs network scanning using Nmap
func (ns *NetworkScanner) scanWithNmap(target string) ([]nmap.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Determine if target is a single IP, IP range, or CIDR
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithContext(ctx),
		nmap.WithTimingTemplate(nmap.TimingAggressive), // Faster scanning
		nmap.WithOSDetection(),                         // OS detection
		nmap.WithServiceInfo(),                         // Service version detection
		nmap.WithScripts("default,safe"),               // Safe scripts
		nmap.WithSkipHostDiscovery(),                   // Skip ping scan if target is specific
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Nmap scanner: %w", err)
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		return nil, fmt.Errorf("Nmap scan failed: %w", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("Nmap warnings: %v\n", warnings)
	}

	if result == nil || len(result.Hosts) == 0 {
		return nil, fmt.Errorf("no hosts found in Nmap scan")
	}

	return result.Hosts, nil
}

// scanWithNaabu is a fallback method using Naabu (original implementation)
func (ns *NetworkScanner) scanWithNaabu(target string, scanID uuid.UUID, startTime time.Time) (*NetworkScanResult, error) {
	var portFindings []NetworkFinding
	var hostsWithOpenPorts []string

	// Run Naabu to discover open ports
	naabuOptions := &runner.Options{
		Host:   []string{target},
		Silent: true,
	}

	var naabuResults []*result.HostResult
	naabuOptions.OnResult = func(hr *result.HostResult) {
		naabuResults = append(naabuResults, hr)
	}

	naabuRunner, err := runner.NewRunner(naabuOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create naabu runner: %w", err)
	}

	err = naabuRunner.RunEnumeration(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to run naabu scan: %w", err)
	}

	for _, result := range naabuResults {
		for _, port := range result.Ports {
			portFindings = append(portFindings, NetworkFinding{
				ID:           uuid.New(),
				FindingType:  "port",
				Severity:     "info",
				Host:         result.IP,
				Port:         port.Port,
				Protocol:     port.Protocol.String(),
				Description:  fmt.Sprintf("Open port %d discovered", port.Port),
				Remediation:  "Review if this service is necessary and secure it properly",
				DiscoveredAt: time.Now(),
				Status:       "open",
			})
			hostsWithOpenPorts = append(hostsWithOpenPorts, fmt.Sprintf("%s:%d", result.IP, port.Port))
		}
	}

	// Run Nuclei on discovered hosts
	var vulnFindings []NetworkFinding
	if len(hostsWithOpenPorts) > 0 {
		uniqueHosts := make(map[string]bool)
		targets := []string{}
		for _, hostPort := range hostsWithOpenPorts {
			parts := strings.Split(hostPort, ":")
			if len(parts) > 0 {
				host := parts[0]
				if !uniqueHosts[host] {
					uniqueHosts[host] = true
					targets = append(targets, host)
				}
			}
		}

		nucleiFindings, err := ns.nucleiScanner.ScanTargets(targets)
		if err != nil {
			nucleiFindings, err = ns.nucleiScanner.ScanUsingCLI(targets)
			if err != nil {
				fmt.Printf("Warning: Nuclei vulnerability scanning failed: %v\n", err)
			}
		}
		vulnFindings = append(vulnFindings, nucleiFindings...)
	}

	allFindings := append(portFindings, vulnFindings...)

	return &NetworkScanResult{
		ID:              scanID,
		AgentID:         uuid.MustParse(ns.config.AgentID),
		CompanyID:       uuid.MustParse(ns.config.CompanyID),
		StartTime:       startTime,
		EndTime:         time.Now(),
		Status:          "completed",
		NetworkFindings: allFindings,
		Metadata: map[string]interface{}{
			"scan_method": "naabu+nuclei",
		},
	}, nil
}

// ScanLocalNetwork scans the local network for devices
func (ns *NetworkScanner) ScanLocalNetwork() (*NetworkScanResult, error) {
	// Get local network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	var targets []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					// Convert to CIDR notation for scanning
					ones, _ := ipnet.Mask.Size()
					cidr := fmt.Sprintf("%s/%d", ipnet.IP.String(), ones)
					targets = append(targets, cidr)
					break // Only scan one network per interface
				}
			}
		}
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no network interfaces found for scanning")
	}

	// Scan the first network (can be extended to scan all)
	return ns.Scan(targets[0])
}
