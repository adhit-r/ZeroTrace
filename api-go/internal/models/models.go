package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	Email          string     `json:"email" db:"email"`
	Password       string     `json:"-" db:"password_hash"`
	Name           string     `json:"name" db:"name"`
	Role           UserRole   `json:"role" db:"role"`
	CompanyID      uuid.UUID  `json:"company_id" db:"company_id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	Status         string     `json:"status" db:"status"`
	LastLogin      *time.Time `json:"last_login,omitempty" db:"last_login"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// UserRole represents user roles
type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)

// Company represents a company/organization
type Company struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Domain    string         `json:"domain" db:"domain"`
	Settings  map[string]any `json:"settings" db:"settings" gorm:"type:jsonb"`
	Status    string         `json:"status" db:"status"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// Organization represents a tenant/organization within a company
type Organization struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	CompanyID   uuid.UUID      `json:"company_id" db:"company_id"`
	Name        string         `json:"name" db:"name"`
	Slug        string         `json:"slug" db:"slug" gorm:"uniqueIndex"`
	Description string         `json:"description" db:"description"`
	Settings    map[string]any `json:"settings" db:"settings" gorm:"type:jsonb"`
	Status      string         `json:"status" db:"status"`
	CreatedAt   time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" db:"updated_at"`
}

// EnrollmentToken represents a token for agent enrollment
type EnrollmentToken struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	Token          string     `json:"token" db:"token" gorm:"uniqueIndex"`
	TokenHash      string     `json:"-" db:"token_hash"`
	IssuedBy       uuid.UUID  `json:"issued_by" db:"issued_by"`
	IssuedAt       time.Time  `json:"issued_at" db:"issued_at"`
	ExpiresAt      time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt         *time.Time `json:"used_at,omitempty" db:"used_at"`
	UsedBy         *uuid.UUID `json:"used_by,omitempty" db:"used_by"`
	Status         string     `json:"status" db:"status"` // active, used, expired, revoked
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// AgentCredential represents a long-lived credential for an agent
type AgentCredential struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	AgentID        uuid.UUID  `json:"agent_id" db:"agent_id"`
	OrganizationID uuid.UUID  `json:"organization_id" db:"organization_id"`
	CredentialHash string     `json:"-" db:"credential_hash"`
	IssuedAt       time.Time  `json:"issued_at" db:"issued_at"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	LastUsedAt     time.Time  `json:"last_used_at" db:"last_used_at"`
	Status         string     `json:"status" db:"status"` // active, expired, revoked
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// Scan represents a vulnerability scan
type Scan struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	CompanyID      uuid.UUID      `json:"company_id" db:"company_id"`
	OrganizationID uuid.UUID      `json:"organization_id" db:"organization_id"`
	AgentID        *uuid.UUID     `json:"agent_id,omitempty" db:"agent_id"`
	Repository     string         `json:"repository" db:"repository"`
	Branch         string         `json:"branch" db:"branch"`
	Commit         string         `json:"commit,omitempty" db:"commit"`
	ScanType       string         `json:"scan_type" db:"scan_type"`
	Status         ScanStatus     `json:"status" db:"status"`
	Progress       int            `json:"progress" db:"progress"`
	StartTime      *time.Time     `json:"start_time,omitempty" db:"start_time"`
	EndTime        *time.Time     `json:"end_time,omitempty" db:"end_time"`
	Options        map[string]any `json:"options" db:"options" gorm:"type:jsonb"`
	Results        map[string]any `json:"results" db:"results" gorm:"type:jsonb"`
	Metadata       map[string]any `json:"metadata" db:"metadata" gorm:"type:jsonb"`
	Notes          string         `json:"notes,omitempty" db:"notes"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

// ScanStatus represents scan status
type ScanStatus string

const (
	ScanStatusPending   ScanStatus = "pending"
	ScanStatusScanning  ScanStatus = "scanning"
	ScanStatusCompleted ScanStatus = "completed"
	ScanStatusFailed    ScanStatus = "failed"
	ScanStatusCancelled ScanStatus = "cancelled"
)

// Vulnerability represents a vulnerability
type Vulnerability struct {
	ID               string         `json:"id" db:"id"`
	ScanID           uuid.UUID      `json:"scan_id" db:"scan_id"`
	CompanyID        uuid.UUID      `json:"company_id" db:"company_id"`
	OrganizationID   uuid.UUID      `json:"organization_id" db:"organization_id"`
	Type             string         `json:"type" db:"type"`
	Severity         SeverityLevel  `json:"severity" db:"severity"`
	Title            string         `json:"title" db:"title"`
	Description      string         `json:"description" db:"description"`
	CVEID            string         `json:"cve_id,omitempty" db:"cve_id"`
	CVSSScore        *float64       `json:"cvss_score,omitempty" db:"cvss_score"`
	CVSSVector       string         `json:"cvss_vector,omitempty" db:"cvss_vector"`
	PackageName      string         `json:"package_name,omitempty" db:"package_name"`
	PackageVersion   string         `json:"package_version,omitempty" db:"package_version"`
	Location         string         `json:"location,omitempty" db:"location"`
	Remediation      string         `json:"remediation,omitempty" db:"remediation"`
	References       []string       `json:"references" db:"references" gorm:"type:jsonb"`
	AffectedVersions []string       `json:"affected_versions" db:"affected_versions" gorm:"type:jsonb"`
	PatchedVersions  []string       `json:"patched_versions" db:"patched_versions" gorm:"type:jsonb"`
	ExploitAvailable bool           `json:"exploit_available" db:"exploit_available"`
	ExploitCount     int            `json:"exploit_count" db:"exploit_count"`
	Status           string         `json:"status" db:"status"`
	Priority         string         `json:"priority" db:"priority"`
	Notes            string         `json:"notes,omitempty" db:"notes"`
	EnrichmentData   map[string]any `json:"enrichment_data" db:"enrichment_data" gorm:"type:jsonb"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// SeverityLevel represents vulnerability severity
type SeverityLevel string

const (
	SeverityCritical SeverityLevel = "CRITICAL"
	SeverityHigh     SeverityLevel = "HIGH"
	SeverityMedium   SeverityLevel = "MEDIUM"
	SeverityLow      SeverityLevel = "LOW"
	SeverityInfo     SeverityLevel = "INFO"
)

// Agent represents an active agent
type Agent struct {
	ID             uuid.UUID `json:"id" db:"id"`
	CompanyID      uuid.UUID `json:"company_id" db:"company_id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Name           string    `json:"name" db:"name"`
	Status         string    `json:"status" db:"status"`
	Version        string    `json:"version" db:"version"`
	LastSeen       time.Time `json:"last_seen" db:"last_seen"`
	CPUUsage       float64   `json:"cpu_usage" db:"cpu_usage"`
	MemoryUsage    float64   `json:"memory_usage" db:"memory_usage"`
	IPAddress      string    `json:"ip_address" db:"ip_address"`
	Hostname       string    `json:"hostname" db:"hostname"`
	OS             string    `json:"os" db:"os"`

	// Enhanced System Information
	OSName         string  `json:"os_name" db:"os_name"`
	OSVersion      string  `json:"os_version" db:"os_version"`
	OSBuild        string  `json:"os_build" db:"os_build"`
	KernelVersion  string  `json:"kernel_version" db:"kernel_version"`
	CPUModel       string  `json:"cpu_model" db:"cpu_model"`
	CPUCores       int     `json:"cpu_cores" db:"cpu_cores"`
	MemoryTotalGB  float64 `json:"memory_total_gb" db:"memory_total_gb"`
	StorageTotalGB float64 `json:"storage_total_gb" db:"storage_total_gb"`
	GPUModel       string  `json:"gpu_model" db:"gpu_model"`
	SerialNumber   string  `json:"serial_number" db:"serial_number"`
	Platform       string  `json:"platform" db:"platform"`
	MACAddress     string  `json:"mac_address" db:"mac_address"`
	City           string  `json:"city" db:"city"`
	Region         string  `json:"region" db:"region"`
	Country        string  `json:"country" db:"country"`
	Timezone       string  `json:"timezone" db:"timezone"`
	RiskScore      float64 `json:"risk_score" db:"risk_score"`
	Tags           string  `json:"tags" db:"tags"` // JSON array as string

	Metadata  map[string]any `json:"metadata" db:"metadata" gorm:"type:jsonb"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// AgentStatus represents the current status of an agent
type AgentStatus struct {
	AgentID        uuid.UUID `json:"agent_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Status         string    `json:"status"`
	LastSeen       time.Time `json:"last_seen"`
	CPUUsage       float64   `json:"cpu_usage"`
	MemoryUsage    float64   `json:"memory_usage"`
	IsOnline       bool      `json:"is_online"`
}

// AgentHeartbeat represents a heartbeat from an agent
type AgentHeartbeat struct {
	AgentID        uuid.UUID      `json:"agent_id"`
	OrganizationID uuid.UUID      `json:"organization_id"`
	AgentName      string         `json:"agent_name"`
	Status         string         `json:"status"`
	CPUUsage       float64        `json:"cpu_usage"`
	MemoryUsage    float64        `json:"memory_usage"`
	Metadata       map[string]any `json:"metadata" gorm:"type:jsonb"`
	Timestamp      time.Time      `json:"timestamp"`
}

// AgentEnrollmentRequest represents an agent enrollment request
type AgentEnrollmentRequest struct {
	EnrollmentToken string              `json:"enrollment_token" binding:"required"`
	AgentInfo       AgentEnrollmentInfo `json:"agent_info" binding:"required"`
}

// AgentEnrollmentInfo represents agent information during enrollment
type AgentEnrollmentInfo struct {
	Hostname     string         `json:"hostname" binding:"required"`
	OS           string         `json:"os" binding:"required"`
	Version      string         `json:"version"`
	Architecture string         `json:"architecture"`
	Metadata     map[string]any `json:"metadata"`
}

// AgentEnrollmentResponse represents the response to agent enrollment
type AgentEnrollmentResponse struct {
	AgentID        uuid.UUID `json:"agent_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Credential     string    `json:"credential"`
	ExpiresAt      time.Time `json:"expires_at"`
}

// Request/Response Models

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=8"`
	Name      string    `json:"name" binding:"required"`
	CompanyID uuid.UUID `json:"company_id" binding:"required"`
}

// CreateScanRequest represents scan creation request
type CreateScanRequest struct {
	Repository string         `json:"repository" binding:"required,url"`
	Branch     string         `json:"branch" binding:"required"`
	ScanType   string         `json:"scan_type"`
	Options    map[string]any `json:"options"`
}

// GenerateEnrollmentTokenRequest represents enrollment token generation request
type GenerateEnrollmentTokenRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
	ExpiresIn      int       `json:"expires_in"` // minutes, default 60
	Description    string    `json:"description"`
}

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Asset represents a scanned asset
type Asset struct {
	ID        uuid.UUID              `json:"id"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// ScanResult represents a scan result from an agent
type ScanResult struct {
	ID              uuid.UUID              `json:"id"`
	AgentID         uuid.UUID              `json:"agent_id"`
	ScanType        string                 `json:"scan_type"`
	Status          string                 `json:"status"`
	Results         map[string]interface{} `json:"results"`
	Vulnerabilities []Vulnerability        `json:"vulnerabilities,omitempty"`
	Assets          []Asset                `json:"assets,omitempty"`
	Metadata        map[string]interface{} `json:"metadata"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// AgentScanResult represents scan results from agents (matches agent format)
type AgentScanResult struct {
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

// Dependency represents a software dependency
type Dependency struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// APIResponse represents standard API response
type APIResponse struct {
	Success   bool      `json:"success"`
	Data      any       `json:"data,omitempty"`
	Message   string    `json:"message,omitempty"`
	Error     *APIError `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// APIError represents API error
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// HealthCheckResponse represents health check response
type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// Organization Profile Models

// OrganizationProfile represents an organization's security profile and settings
type OrganizationProfile struct {
	ID                   uuid.UUID      `json:"id" db:"id"`
	OrganizationID       uuid.UUID      `json:"organization_id" db:"organization_id"`
	Industry             string         `json:"industry" db:"industry"`
	RiskTolerance        RiskTolerance  `json:"risk_tolerance" db:"risk_tolerance"`
	TechStack            TechStack      `json:"tech_stack" db:"tech_stack" gorm:"type:jsonb"`
	ComplianceFrameworks []string       `json:"compliance_frameworks" db:"compliance_frameworks" gorm:"type:jsonb"`
	SecurityPolicies     map[string]any `json:"security_policies" db:"security_policies" gorm:"type:jsonb"`
	RiskWeights          map[string]any `json:"risk_weights" db:"risk_weights" gorm:"type:jsonb"`
	CreatedAt            time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at" db:"updated_at"`
}

// RiskTolerance represents organization's risk tolerance level
type RiskTolerance string

const (
	RiskToleranceConservative RiskTolerance = "CONSERVATIVE"
	RiskToleranceModerate     RiskTolerance = "MODERATE"
	RiskToleranceAggressive   RiskTolerance = "AGGRESSIVE"
)

// TechStack represents organization's technology stack
type TechStack struct {
	Languages        []string `json:"languages"`
	Frameworks       []string `json:"frameworks"`
	Databases        []string `json:"databases"`
	CloudProviders   []string `json:"cloud_providers"`
	OperatingSystems []string `json:"operating_systems"`
	Containers       []string `json:"containers"`
	DevTools         []string `json:"dev_tools"`
	SecurityTools    []string `json:"security_tools"`
}

// CreateOrganizationProfileRequest represents request to create organization profile
type CreateOrganizationProfileRequest struct {
	OrganizationID       uuid.UUID      `json:"organization_id" binding:"required"`
	Industry             string         `json:"industry" binding:"required"`
	RiskTolerance        RiskTolerance  `json:"risk_tolerance" binding:"required"`
	TechStack            TechStack      `json:"tech_stack" binding:"required"`
	ComplianceFrameworks []string       `json:"compliance_frameworks"`
	SecurityPolicies     map[string]any `json:"security_policies"`
	RiskWeights          map[string]any `json:"risk_weights"`
}

// UpdateOrganizationProfileRequest represents request to update organization profile
type UpdateOrganizationProfileRequest struct {
	Industry             *string        `json:"industry,omitempty"`
	RiskTolerance        *RiskTolerance `json:"risk_tolerance,omitempty"`
	TechStack            *TechStack     `json:"tech_stack,omitempty"`
	ComplianceFrameworks *[]string      `json:"compliance_frameworks,omitempty"`
	SecurityPolicies     map[string]any `json:"security_policies,omitempty"`
	RiskWeights          map[string]any `json:"risk_weights,omitempty"`
}

// Software represents installed software on an agent
type Software struct {
	ID        uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AgentID   uuid.UUID `json:"agent_id" db:"agent_id"`
	Name      string    `json:"name" db:"name"`
	Version   string    `json:"version" db:"version"`
	Type      string    `json:"type" db:"type"`
	Vendor    string    `json:"vendor" db:"vendor"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NetworkHost represents a host discovered during network scanning
type NetworkHost struct {
	ID         uuid.UUID      `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AgentID    uuid.UUID      `json:"agent_id" db:"agent_id"`
	IPAddress  string         `json:"ip_address" db:"ip_address"`
	Hostname   string         `json:"hostname" db:"hostname"`
	MACAddress string         `json:"mac_address" db:"mac_address"`
	OS         string         `json:"os" db:"os"`
	Status     string         `json:"status" db:"status"`
	OpenPorts  []int          `json:"open_ports" db:"open_ports" gorm:"type:jsonb"`
	Metadata   map[string]any `json:"metadata" db:"metadata" gorm:"type:jsonb"`
	LastSeen   time.Time      `json:"last_seen" db:"last_seen"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
}

// DashboardSnapshot represents a historical snapshot of dashboard metrics
type DashboardSnapshot struct {
	ID                   uuid.UUID `json:"id" db:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID       uuid.UUID `json:"organization_id" db:"organization_id"`
	Date                 time.Time `json:"date" db:"date"`
	TotalAssets          int       `json:"total_assets" db:"total_assets"`
	TotalVulnerabilities int       `json:"total_vulnerabilities" db:"total_vulnerabilities"`
	CriticalVulns        int       `json:"critical_vulns" db:"critical_vulns"`
	HighVulns            int       `json:"high_vulns" db:"high_vulns"`
	MediumVulns          int       `json:"medium_vulns" db:"medium_vulns"`
	LowVulns             int       `json:"low_vulns" db:"low_vulns"`
	AvgRiskScore         float64   `json:"avg_risk_score" db:"avg_risk_score"`
	ComplianceScore      float64   `json:"compliance_score" db:"compliance_score"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}
