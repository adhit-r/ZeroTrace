package discovery

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"zerotrace/agent/internal/models"
)

// NetworkDiscovery handles network asset discovery
type NetworkDiscovery struct {
	agentID   string
	companyID string
}

// NewNetworkDiscovery creates a new network discovery instance
func NewNetworkDiscovery(agentID, companyID string) *NetworkDiscovery {
	return &NetworkDiscovery{
		agentID:   agentID,
		companyID: companyID,
	}
}

// DiscoverLocalNetwork discovers assets on the local network
func (nd *NetworkDiscovery) DiscoverLocalNetwork(ctx context.Context) ([]models.NetworkAsset, error) {
	var assets []models.NetworkAsset

	// Get local network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

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
					// Discover assets in this subnet
					subnetAssets, err := nd.scanSubnet(ctx, ipnet)
					if err != nil {
						continue
					}
					assets = append(assets, subnetAssets...)
				}
			}
		}
	}

	// Perform network topology analysis using Fast SSSP
	if len(assets) > 0 {
		analyzer := NewNetworkPathAnalyzer()
		topology, err := analyzer.AnalyzeNetworkTopology(ctx, assets)
		if err != nil {
			log.Printf("Warning: Failed to analyze network topology: %v", err)
		} else {
			log.Printf("Network topology analysis complete: %d nodes, %d connections, %d critical paths found", 
				topology.TotalAssets, topology.TotalConnections, len(topology.CriticalPaths))
		}
	}

	return assets, nil
}

// scanSubnet scans a subnet for active hosts
func (nd *NetworkDiscovery) scanSubnet(ctx context.Context, ipnet *net.IPNet) ([]models.NetworkAsset, error) {
	var assets []models.NetworkAsset

	// Get network information
	network := ipnet.IP.Mask(ipnet.Mask)
	ones, bits := ipnet.Mask.Size()
	hosts := 1 << (bits - ones)

	// Scan first 254 hosts in the subnet (common for /24 networks)
	maxHosts := 254
	if hosts < maxHosts {
		maxHosts = hosts
	}

	for i := 1; i < maxHosts; i++ {
		select {
		case <-ctx.Done():
			return assets, ctx.Err()
		default:
		}

		ip := make(net.IP, len(network))
		copy(ip, network)
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j] += byte(i)
			if ip[j] != 0 {
				break
			}
		}

		// Quick ping to check if host is alive
		if nd.isHostAlive(ip.String()) {
			asset, err := nd.discoverHost(ctx, ip.String())
			if err != nil {
				continue
			}
			assets = append(assets, asset)
		}
	}

	return assets, nil
}

// isHostAlive checks if a host is alive using ping
func (nd *NetworkDiscovery) isHostAlive(ip string) bool {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("ping", "-n", "1", "-w", "1000", ip)
	case "darwin", "linux":
		cmd = exec.Command("ping", "-c", "1", "-W", "1", ip)
	default:
		return false
	}

	err := cmd.Run()
	return err == nil
}

// discoverHost discovers detailed information about a host
func (nd *NetworkDiscovery) discoverHost(ctx context.Context, ip string) (models.NetworkAsset, error) {
	asset := models.NetworkAsset{
		ID:         fmt.Sprintf("asset-%s-%s", nd.agentID, ip),
		AgentID:    nd.agentID,
		CompanyID:  nd.companyID,
		IPAddress:  ip,
		LastSeen:   time.Now(),
		IsMonitored: false,
		Metadata:   make(map[string]interface{}),
	}

	// Get hostname
	if hostname, err := net.LookupAddr(ip); err == nil && len(hostname) > 0 {
		asset.Hostname = strings.TrimSuffix(hostname[0], ".")
	}

	// Get MAC address (ARP table)
	if mac, err := nd.getMACAddress(ip); err == nil {
		asset.MACAddress = mac
	}

	// Detect OS and services
	asset.OS, asset.OSVersion = nd.detectOS(ip)
	asset.OpenPorts = nd.scanPorts(ip)
	asset.RunningServices = nd.detectServices(ip)
	asset.ConnectedPeers = nd.getConnectedPeers(ip)

	// Calculate risk score
	asset.RiskScore = nd.calculateRiskScore(asset)

	return asset, nil
}

// getMACAddress gets MAC address from ARP table
func (nd *NetworkDiscovery) getMACAddress(ip string) (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("arp", "-a", ip)
	case "darwin", "linux":
		cmd = exec.Command("arp", "-n", ip)
	default:
		return "", fmt.Errorf("unsupported OS")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse ARP output to extract MAC address
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ip) {
			fields := strings.Fields(line)
			for i, field := range fields {
				if strings.Contains(field, ":") && len(field) == 17 {
					return field, nil
				}
				if i > 0 && strings.Contains(fields[i-1], ip) && strings.Contains(field, "-") {
					return field, nil
				}
			}
		}
	}

	return "", fmt.Errorf("MAC address not found")
}

// detectOS attempts to detect the operating system
func (nd *NetworkDiscovery) detectOS(ip string) (string, string) {
	// This is a simplified OS detection
	// In a real implementation, you'd use more sophisticated techniques
	
	// Check for common ports to infer OS
	ports := []int{22, 3389, 135, 139, 445}
	
	sshOpen := false
	windowsOpen := false
	
	for _, port := range ports {
		if nd.isPortOpen(ip, port) {
			switch port {
			case 22:
				sshOpen = true
			case 3389, 135, 139, 445:
				windowsOpen = true
			}
		}
	}
	
	if windowsOpen {
		return "Windows", "Unknown"
	} else if sshOpen {
		return "Linux/Unix", "Unknown"
	}
	
	return "Unknown", "Unknown"
}

// isPortOpen checks if a port is open
func (nd *NetworkDiscovery) isPortOpen(ip string, port int) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// scanPorts scans for open ports
func (nd *NetworkDiscovery) scanPorts(ip string) []models.PortInfo {
	var ports []models.PortInfo
	
	// Common ports to scan
	commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 143, 443, 993, 995, 3389, 135, 139, 445, 1433, 1521, 3306, 5432, 6379, 8080, 8443}
	
	for _, port := range commonPorts {
		if nd.isPortOpen(ip, port) {
			portInfo := models.PortInfo{
				Port:     port,
				Protocol: "tcp",
				Service:  nd.getServiceName(port),
				IsSecure: nd.isSecurePort(port),
			}
			ports = append(ports, portInfo)
		}
	}
	
	return ports
}

// getServiceName returns the service name for a port
func (nd *NetworkDiscovery) getServiceName(port int) string {
	services := map[int]string{
		21: "ftp", 22: "ssh", 23: "telnet", 25: "smtp", 53: "dns",
		80: "http", 110: "pop3", 143: "imap", 443: "https", 993: "imaps",
		995: "pop3s", 3389: "rdp", 135: "rpc", 139: "netbios", 445: "smb",
		1433: "mssql", 1521: "oracle", 3306: "mysql", 5432: "postgresql",
		6379: "redis", 8080: "http-proxy", 8443: "https-alt",
	}
	
	if service, exists := services[port]; exists {
		return service
	}
	return "unknown"
}

// isSecurePort checks if a port is typically secure
func (nd *NetworkDiscovery) isSecurePort(port int) bool {
	securePorts := []int{443, 993, 995, 8443}
	for _, securePort := range securePorts {
		if port == securePort {
			return true
		}
	}
	return false
}

// detectServices detects running services
func (nd *NetworkDiscovery) detectServices(ip string) []models.ServiceInfo {
	// This is a placeholder - in a real implementation you'd use
	// more sophisticated service detection techniques
	return []models.ServiceInfo{}
}

// getConnectedPeers gets information about connected peers
func (nd *NetworkDiscovery) getConnectedPeers(ip string) []models.PeerInfo {
	// This is a placeholder - in a real implementation you'd parse
	// ARP tables, routing tables, etc.
	return []models.PeerInfo{}
}

// calculateRiskScore calculates a risk score for an asset
func (nd *NetworkDiscovery) calculateRiskScore(asset models.NetworkAsset) float64 {
	score := 0.0
	
	// Base score
	score += 10
	
	// Add score for each open port
	score += float64(len(asset.OpenPorts)) * 2
	
	// Add score for insecure services
	for _, port := range asset.OpenPorts {
		if !port.IsSecure && (port.Port == 21 || port.Port == 23 || port.Port == 135 || port.Port == 139) {
			score += 5
		}
	}
	
	// Cap the score at 100
	if score > 100 {
		score = 100
	}
	
	return score
}
