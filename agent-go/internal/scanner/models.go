package scanner

import (
	"time"

	"github.com/google/uuid"
)

// NetworkScanResult represents the result of a comprehensive network scan.
type NetworkScanResult struct {
	ID              uuid.UUID              `json:"id"`
	AgentID         uuid.UUID              `json:"agent_id"`
	CompanyID       uuid.UUID              `json:"company_id"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         time.Time              `json:"end_time"`
	Status          string                 `json:"status"`
	NetworkFindings []NetworkFinding       `json:"network_findings"`
	SSLResults      *SSLAudit              `json:"ssl_audit_results,omitempty"`
	SBOMResults     *SBOMFinding           `json:"sbom_results,omitempty"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// NetworkFinding represents a single network security finding from Naabu or Nuclei.
type NetworkFinding struct {
	ID             uuid.UUID              `json:"id"`
	FindingType    string                 `json:"finding_type"` // e.g., "port", "vuln", "config"
	Severity       string                 `json:"severity"`     // critical, high, medium, low, info
	Host           string                 `json:"host"`
	Port           int                    `json:"port"`
	Protocol       string                 `json:"protocol"` // tcp, udp
	ServiceName    string                 `json:"service_name"`
	ServiceVersion string                 `json:"service_version"`
	Banner         string                 `json:"banner"`
	Description    string                 `json:"description"`
	Remediation    string                 `json:"remediation"`
	DiscoveredAt   time.Time              `json:"discovered_at"`
	Status         string                 `json:"status"` // open, filtered, closed
	DeviceType     string                 `json:"device_type,omitempty"` // switch, router, iot, phone, server, unknown
	OS             string                 `json:"os,omitempty"`
	OSVersion      string                 `json:"os_version,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// SBOMFinding represents the result of an SBOM scan.
type SBOMFinding struct {
	ServiceName     string          `json:"service_name"`
	Components      []Component     `json:"components"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	RiskScore       float64         `json:"risk_score"`
}

// Component represents a software component in the SBOM.
type Component struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Type     string   `json:"type"`
	Licenses []string `json:"licenses"`
	PURL     string   `json:"purl"`
}

// Vulnerability represents a CVE found in a component.
type Vulnerability struct {
	CVE          string  `json:"cve"`
	Severity     string  `json:"severity"`
	CVSS         float64 `json:"cvss"`
	Description  string  `json:"description"`
	FixAvailable bool    `json:"fix_available"`
	FixVersion   string  `json:"fix_version"`
}

// SSLAudit represents the result of an SSL/TLS security audit.
type SSLAudit struct {
	Host             string        `json:"host"`
	Port             int           `json:"port"`
	CertificateChain []Certificate `json:"certificate_chain"`
	Issues           []SSLIssue    `json:"issues"`
	Grade            string        `json:"grade"` // A+, A, B, C, F
	Recommendations  []string      `json:"recommendations"`
}

// Certificate represents details of a TLS certificate.
type Certificate struct {
	Subject      string    `json:"subject"`
	Issuer       string    `json:"issuer"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
	SelfSigned   bool      `json:"self_signed"`
	KeySize      int       `json:"key_size"`
	SignatureAlg string    `json:"signature_alg"`
	SANs         []string  `json:"sans"`
}

// SSLIssue represents a specific security issue found during the audit.
type SSLIssue struct {
	Severity    string `json:"severity"`
	Type        string `json:"type"` // expired, weak-cipher, protocol-vuln
	Description string `json:"description"`
	CVE         string `json:"cve,omitempty"`
	Remediation string `json:"remediation"`
}

// ChangeDetection represents the result of an attack surface change analysis.
type ChangeDetection struct {
	ScanID       uuid.UUID `json:"scan_id"`
	PreviousScan uuid.UUID `json:"previous_scan_id"`
	Timestamp    time.Time `json:"timestamp"`
	Changes      []Change  `json:"changes"`
	RiskDelta    float64   `json:"risk_delta"` // Change in overall risk score
}

// Change represents a single detected change between two scans.
type Change struct {
	Type        string  `json:"type"` // new-port, closed-port, new-vuln, fixed-vuln, config-change
	Severity    string  `json:"severity"`
	Description string  `json:"description"`
	Entity      string  `json:"entity"` // What changed (e.g., "port-22", "ssl-cert")
	OldValue    string  `json:"old_value,omitempty"`
	NewValue    string  `json:"new_value,omitempty"`
	RiskImpact  float64 `json:"risk_impact"`
}
