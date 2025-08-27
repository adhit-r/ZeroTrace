package models

import (
	"time"

	"github.com/google/uuid"
)

// ScanResult represents the result of a vulnerability scan
type ScanResult struct {
	ID              uuid.UUID       `json:"id"`
	AgentID         string          `json:"agent_id"`
	CompanyID       string          `json:"company_id"`
	Repository      string          `json:"repository,omitempty"`
	Branch          string          `json:"branch,omitempty"`
	Commit          string          `json:"commit,omitempty"`
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	Status          string          `json:"status"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Dependencies    []Dependency    `json:"dependencies"`
	Metadata        map[string]any  `json:"metadata"`
}

// Vulnerability represents a detected vulnerability
type Vulnerability struct {
	ID               string         `json:"id"`
	Type             string         `json:"type"`
	Severity         string         `json:"severity"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	CVEID            string         `json:"cve_id,omitempty"`
	CVSSScore        *float64       `json:"cvss_score,omitempty"`
	CVSSVector       string         `json:"cvss_vector,omitempty"`
	PackageName      string         `json:"package_name,omitempty"`
	PackageVersion   string         `json:"package_version,omitempty"`
	Location         string         `json:"location,omitempty"`
	Remediation      string         `json:"remediation,omitempty"`
	References       []string       `json:"references"`
	AffectedVersions []string       `json:"affected_versions"`
	PatchedVersions  []string       `json:"patched_versions"`
	ExploitAvailable bool           `json:"exploit_available"`
	ExploitCount     int            `json:"exploit_count"`
	Status           string         `json:"status"`
	Priority         string         `json:"priority"`
	Notes            string         `json:"notes,omitempty"`
	EnrichmentData   map[string]any `json:"enrichment_data"`
	CreatedAt        time.Time      `json:"created_at"`
}

// Dependency represents a detected dependency
type Dependency struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Version         string          `json:"version"`
	Type            string          `json:"type"`
	Location        string          `json:"location"`
	Path            string          `json:"path,omitempty"`
	InstallDate     time.Time       `json:"install_date,omitempty"`
	Size            int64           `json:"size,omitempty"`
	Vendor          string          `json:"vendor,omitempty"`
	Description     string          `json:"description,omitempty"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Metadata        map[string]any  `json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
}

// FileInfo represents information about a scanned file
type FileInfo struct {
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	ModTime     time.Time `json:"mod_time"`
	Language    string    `json:"language"`
	LinesOfCode int       `json:"lines_of_code"`
	Hash        string    `json:"hash"`
}

// ScanOptions represents scanning configuration
type ScanOptions struct {
	Depth            int      `json:"depth"`
	IncludePatterns  []string `json:"include_patterns"`
	ExcludePatterns  []string `json:"exclude_patterns"`
	MaxFileSize      int64    `json:"max_file_size"`
	MaxConcurrency   int      `json:"max_concurrency"`
	ScanDependencies bool     `json:"scan_dependencies"`
	ScanSecrets      bool     `json:"scan_secrets"`
	ScanSAST         bool     `json:"scan_sast"`
}

// ScanProgress represents scan progress information
type ScanProgress struct {
	TotalFiles      int       `json:"total_files"`
	ScannedFiles    int       `json:"scanned_files"`
	Vulnerabilities int       `json:"vulnerabilities"`
	Progress        float64   `json:"progress"`
	CurrentFile     string    `json:"current_file"`
	StartTime       time.Time `json:"start_time"`
	EstimatedEnd    time.Time `json:"estimated_end"`
}

// AgentStatus represents agent status information
type AgentStatus struct {
	AgentID              string         `json:"agent_id"`
	Status               string         `json:"status"`
	LastSeen             time.Time      `json:"last_seen"`
	CurrentScan          *string        `json:"current_scan,omitempty"`
	ScansCompleted       int            `json:"scans_completed"`
	VulnerabilitiesFound int            `json:"vulnerabilities_found"`
	PerformanceMetrics   map[string]any `json:"performance_metrics"`
}

// APIResponse represents API response structure
type APIResponse struct {
	Success   bool      `json:"success"`
	Data      any       `json:"data,omitempty"`
	Message   string    `json:"message,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// APIError represents API error structure
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// NetworkAsset represents a discovered network asset
type NetworkAsset struct {
	ID              string                 `json:"id"`
	AgentID         string                 `json:"agent_id"`
	CompanyID       string                 `json:"company_id"`
	IPAddress       string                 `json:"ip_address"`
	MACAddress      string                 `json:"mac_address"`
	Hostname        string                 `json:"hostname"`
	OS              string                 `json:"os"`
	OSVersion       string                 `json:"os_version"`
	DeviceType      string                 `json:"device_type"` // server, workstation, network_device, etc.
	Location        string                 `json:"location"`    // building, floor, room
	Department      string                 `json:"department"`
	Subnet          string                 `json:"subnet"`
	VLAN            string                 `json:"vlan"`
	OpenPorts       []PortInfo             `json:"open_ports"`
	RunningServices []ServiceInfo          `json:"running_services"`
	ConnectedPeers  []PeerInfo             `json:"connected_peers"`
	Vulnerabilities []Vulnerability        `json:"vulnerabilities"`
	RiskScore       float64                `json:"risk_score"`
	LastSeen        time.Time              `json:"last_seen"`
	IsMonitored     bool                   `json:"is_monitored"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// PortInfo represents information about an open port
type PortInfo struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Service  string `json:"service"`
	Version  string `json:"version"`
	Banner   string `json:"banner"`
	IsSecure bool   `json:"is_secure"`
}

// ServiceInfo represents information about a running service
type ServiceInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	ProcessID   int    `json:"process_id"`
	User        string `json:"user"`
	CommandLine string `json:"command_line"`
}

// PeerInfo represents information about connected peers
type PeerInfo struct {
	IPAddress      string    `json:"ip_address"`
	MACAddress     string    `json:"mac_address"`
	Hostname       string    `json:"hostname"`
	ConnectionType string    `json:"connection_type"` // arp, ndp, routing, etc.
	RiskScore      float64   `json:"risk_score"`
	LastSeen       time.Time `json:"last_seen"`
}

// AmassResult represents OWASP Amass discovery results
type AmassResult struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	Domain       string    `json:"domain"`
	Subdomain    string    `json:"subdomain"`
	IPAddress    string    `json:"ip_address"`
	ASN          int       `json:"asn"`
	ASName       string    `json:"as_name"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
	Service      string    `json:"service"`
	Port         int       `json:"port"`
	Protocol     string    `json:"protocol"`
	Source       string    `json:"source"` // amass, dns, certificate, etc.
	DiscoveredAt time.Time `json:"discovered_at"`
	IsRelated    bool      `json:"is_related"` // related to company assets
}

// NetworkTopology represents the network topology map
type NetworkTopology struct {
	ID               string         `json:"id"`
	CompanyID        string         `json:"company_id"`
	Nodes            []TopologyNode `json:"nodes"`
	Links            []TopologyLink `json:"links"`
	Clusters         []Cluster      `json:"clusters"`
	CriticalPaths    []*NetworkPath `json:"critical_paths"`
	TotalAssets      int            `json:"total_assets"`
	TotalConnections int            `json:"total_connections"`
	LastUpdated      time.Time      `json:"last_updated"`
}

// TopologyNode represents a node in the topology map
type TopologyNode struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // agent, asset, amass_discovery
	AssetID     string                 `json:"asset_id,omitempty"`
	AgentID     string                 `json:"agent_id,omitempty"`
	Name        string                 `json:"name"`
	IPAddress   string                 `json:"ip_address"`
	Location    string                 `json:"location"`
	Department  string                 `json:"department"`
	ClusterID   string                 `json:"cluster_id,omitempty"`
	RiskScore   float64                `json:"risk_score"`
	IsMonitored bool                   `json:"is_monitored"`
	Status      string                 `json:"status"` // active, inactive, discovered
	Metadata    map[string]interface{} `json:"metadata"`
}

// TopologyLink represents a connection between nodes
type TopologyLink struct {
	Source     string  `json:"source"`
	Target     string  `json:"target"`
	Type       string  `json:"type"` // network, scan, external
	Strength   float64 `json:"strength"`
	Weight     float64 `json:"weight"`
	IsCritical bool    `json:"is_critical"`
	Protocol   string  `json:"protocol,omitempty"`
	Port       int     `json:"port,omitempty"`
}

// Cluster represents a logical grouping of assets
type Cluster struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // subnet, department, floor, geographic
	NodeIDs     []string `json:"node_ids"`
	RiskScore   float64  `json:"risk_score"`
	Description string   `json:"description"`
}

// AgentTelemetry represents real-time agent status and metrics
type AgentTelemetry struct {
	AgentID       string                 `json:"agent_id"`
	CompanyID     string                 `json:"company_id"`
	Status        string                 `json:"status"` // online, offline, scanning
	LastHeartbeat time.Time              `json:"last_heartbeat"`
	CPUUsage      float64                `json:"cpu_usage"`
	MemoryUsage   float64                `json:"memory_usage"`
	DiskUsage     float64                `json:"disk_usage"`
	NetworkStats  NetworkStats           `json:"network_stats"`
	ActiveScans   int                    `json:"active_scans"`
	AssetsScanned int                    `json:"assets_scanned"`
	VulnsFound    int                    `json:"vulns_found"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// NetworkPath represents a path between two network assets
type NetworkPath struct {
	Source      string    `json:"source"`
	Destination string    `json:"destination"`
	Path        []string  `json:"path"`
	Distance    float64   `json:"distance"`
	Hops        int       `json:"hops"`
	Latency     float64   `json:"latency_ms"`
	RiskScore   float64   `json:"risk_score"`
	Discovered  time.Time `json:"discovered_at"`
}

// NetworkStats represents network statistics
type NetworkStats struct {
	BytesSent       int64 `json:"bytes_sent"`
	BytesReceived   int64 `json:"bytes_received"`
	PacketsSent     int64 `json:"packets_sent"`
	PacketsReceived int64 `json:"packets_received"`
	Connections     int   `json:"connections"`
}

// InstalledApp represents an installed application on the system
type InstalledApp struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Type        string    `json:"type"` // macos_app, homebrew, apt, etc.
	Path        string    `json:"path"`
	InstallDate time.Time `json:"install_date"`
	Size        int64     `json:"size"`
	Vendor      string    `json:"vendor"`
	Description string    `json:"description"`
}
