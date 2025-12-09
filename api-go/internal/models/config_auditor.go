package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ConfigFile represents an uploaded configuration file
type ConfigFile struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyID uuid.UUID      `json:"company_id" gorm:"type:uuid;not null;index"`
	UploadedBy *uuid.UUID    `json:"uploaded_by,omitempty" gorm:"type:uuid"`

	// File metadata
	Filename  string `json:"filename" gorm:"not null"`
	FilePath  string `json:"file_path" gorm:"not null"`
	FileSize  int64  `json:"file_size" gorm:"not null"`
	FileHash  string `json:"file_hash" gorm:"not null;size:64;index"`
	MimeType  string `json:"mime_type,omitempty" gorm:"size:100"`
	FileContent []byte `json:"-" gorm:"type:bytea"` // Store file content in PostgreSQL

	// Device information
	DeviceType      string `json:"device_type" gorm:"not null;size:50"`
	Manufacturer    string `json:"manufacturer" gorm:"not null;size:100;index"`
	Model           string `json:"model,omitempty" gorm:"size:100"`
	FirmwareVersion string `json:"firmware_version,omitempty" gorm:"size:100"`
	DeviceName      string `json:"device_name,omitempty" gorm:"size:255"`
	DeviceLocation  string `json:"device_location,omitempty" gorm:"size:255"`

	// Configuration metadata
	ConfigType  string `json:"config_type" gorm:"not null;size:50"`
	ConfigFormat string `json:"config_format,omitempty" gorm:"size:50"`
	ConfigVersion string `json:"config_version,omitempty" gorm:"size:50"`

	// Parsing status
	ParsingStatus string         `json:"parsing_status" gorm:"default:'pending';size:50;index"`
	ParsingError  string         `json:"parsing_error,omitempty" gorm:"type:text"`
	ParsedData    datatypes.JSON `json:"parsed_data,omitempty" gorm:"type:jsonb;default:'{}'"`

	// Analysis status
	AnalysisStatus    string     `json:"analysis_status" gorm:"default:'pending';size:50;index"`
	AnalysisStartedAt *time.Time `json:"analysis_started_at,omitempty"`
	AnalysisCompletedAt *time.Time `json:"analysis_completed_at,omitempty"`

	// Metadata
	Tags     datatypes.JSON `json:"tags,omitempty" gorm:"type:jsonb;default:'[]'"`
	Notes    string         `json:"notes,omitempty" gorm:"type:text"`
	Metadata datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb;default:'{}'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Company          Company                `json:"-" gorm:"foreignKey:CompanyID"`
	Uploader         *User                  `json:"uploader,omitempty" gorm:"foreignKey:UploadedBy"`
	Findings         []ConfigFinding        `json:"findings,omitempty" gorm:"foreignKey:ConfigFileID"`
	AnalysisResult   *ConfigAnalysisResult  `json:"analysis_result,omitempty" gorm:"foreignKey:ConfigFileID"`
}

// TableName specifies the table name
func (ConfigFile) TableName() string {
	return "config_files"
}

// IsParsed returns true if the config file has been parsed
func (c *ConfigFile) IsParsed() bool {
	return c.ParsingStatus == "parsed"
}

// IsAnalyzed returns true if the config file has been analyzed
func (c *ConfigFile) IsAnalyzed() bool {
	return c.AnalysisStatus == "completed"
}

// ConfigFinding represents a security finding from configuration analysis
type ConfigFinding struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConfigFileID uuid.UUID `json:"config_file_id" gorm:"type:uuid;not null;index"`
	CompanyID   uuid.UUID `json:"company_id" gorm:"type:uuid;not null;index"`

	// Finding details
	FindingType string `json:"finding_type" gorm:"not null;size:50;index"`
	Severity    string `json:"severity" gorm:"not null;size:20;index"`
	Category    string `json:"category" gorm:"not null;size:50;index"`

	// Finding information
	Title            string `json:"title" gorm:"not null;size:500"`
	Description      string `json:"description" gorm:"not null;type:text"`
	AffectedComponent string `json:"affected_component,omitempty" gorm:"size:255"`
	ConfigSnippet    string `json:"config_snippet,omitempty" gorm:"type:text"`
	LineNumbers      datatypes.JSON `json:"line_numbers,omitempty" gorm:"type:jsonb"` // Array of integers

	// Standards and compliance
	StandardID          *uuid.UUID    `json:"standard_id,omitempty" gorm:"type:uuid;index"`
	ComplianceFrameworks datatypes.JSON `json:"compliance_frameworks,omitempty" gorm:"type:jsonb;default:'[]'"`
	CVEID               string        `json:"cve_id,omitempty" gorm:"size:20"`
	CVSSScore           *float64      `json:"cvss_score,omitempty" gorm:"type:decimal(3,1)"`

	// Remediation
	Remediation        string         `json:"remediation,omitempty" gorm:"type:text"`
	RemediationSteps   datatypes.JSON `json:"remediation_steps,omitempty" gorm:"type:jsonb;default:'[]'"`
	RemediationPriority string        `json:"remediation_priority" gorm:"default:'medium';size:20"`
	EstimatedEffort    string         `json:"estimated_effort,omitempty" gorm:"size:50"`

	// Risk assessment
	RiskScore    float64 `json:"risk_score" gorm:"type:decimal(3,2);default:0.0"`
	Exploitability string `json:"exploitability,omitempty" gorm:"size:20"`
	Impact        string `json:"impact,omitempty" gorm:"size:20"`

	// Status
	Status      string     `json:"status" gorm:"default:'open';size:50;index"`
	AssignedTo  *uuid.UUID `json:"assigned_to,omitempty" gorm:"type:uuid"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy  *uuid.UUID `json:"resolved_by,omitempty" gorm:"type:uuid"`

	// Evidence and references
	Evidence    datatypes.JSON `json:"evidence,omitempty" gorm:"type:jsonb;default:'{}'"`
	References  datatypes.JSON `json:"references,omitempty" gorm:"type:jsonb;default:'[]'"`
	Tags        datatypes.JSON `json:"tags,omitempty" gorm:"type:jsonb;default:'[]'"`
	Metadata    datatypes.JSON `json:"metadata,omitempty" gorm:"type:jsonb;default:'{}'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	ConfigFile ConfigFile      `json:"-" gorm:"foreignKey:ConfigFileID"`
	Company    Company         `json:"-" gorm:"foreignKey:CompanyID"`
	Standard   *ConfigStandard `json:"standard,omitempty" gorm:"foreignKey:StandardID"`
	Assignee   *User           `json:"assignee,omitempty" gorm:"foreignKey:AssignedTo"`
	Resolver   *User           `json:"resolver,omitempty" gorm:"foreignKey:ResolvedBy"`
}

// TableName specifies the table name
func (ConfigFinding) TableName() string {
	return "config_findings"
}

// IsOpen returns true if the finding is open
func (c *ConfigFinding) IsOpen() bool {
	return c.Status == "open"
}

// IsResolved returns true if the finding is resolved
func (c *ConfigFinding) IsResolved() bool {
	return c.Status == "resolved"
}

// ConfigStandard represents a manufacturer-specific security standard
type ConfigStandard struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Standard identification
	StandardName    string `json:"standard_name" gorm:"not null;size:255"`
	StandardVersion string `json:"standard_version,omitempty" gorm:"size:50"`
	Manufacturer    string `json:"manufacturer" gorm:"not null;size:100;index"`
	DeviceType      string `json:"device_type" gorm:"not null;size:50;index"`
	ModelFamily     string `json:"model_family,omitempty" gorm:"size:100"`

	// Standard details
	Category              string `json:"category" gorm:"not null;size:50;index"`
	RequirementID         string `json:"requirement_id" gorm:"not null;size:100"`
	RequirementTitle      string `json:"requirement_title" gorm:"not null;size:500"`
	RequirementDescription string `json:"requirement_description,omitempty" gorm:"type:text"`

	// Compliance mapping
	ComplianceFrameworks  datatypes.JSON `json:"compliance_frameworks,omitempty" gorm:"type:jsonb;default:'[]';index"`
	ComplianceRequirement string         `json:"compliance_requirement,omitempty" gorm:"size:500"`

	// Configuration check
	CheckType      string `json:"check_type" gorm:"not null;size:50"`
	CheckConfigPath string `json:"check_config_path,omitempty" gorm:"size:500"`
	CheckPattern   string `json:"check_pattern,omitempty" gorm:"type:text"`
	ExpectedValue  string `json:"expected_value,omitempty" gorm:"type:text"`
	CheckScript    string `json:"check_script,omitempty" gorm:"type:text"`

	// Severity and priority
	DefaultSeverity string `json:"default_severity" gorm:"not null;size:20"`
	Priority        string `json:"priority" gorm:"default:'medium';size:20"`

	// Remediation guidance
	RemediationGuidance string `json:"remediation_guidance,omitempty" gorm:"type:text"`
	RemediationExample  string `json:"remediation_example,omitempty" gorm:"type:text"`
	RemediationScript   string `json:"remediation_script,omitempty" gorm:"type:text"`

	// References
	References      datatypes.JSON `json:"references,omitempty" gorm:"type:jsonb;default:'[]'"`
	DocumentationURL string         `json:"documentation_url,omitempty" gorm:"size:500"`

	// Status
	Status string `json:"status" gorm:"default:'active';size:50;index"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Findings []ConfigFinding `json:"-" gorm:"foreignKey:StandardID"`
}

// TableName specifies the table name
func (ConfigStandard) TableName() string {
	return "config_standards"
}

// IsActive returns true if the standard is active
func (c *ConfigStandard) IsActive() bool {
	return c.Status == "active"
}

// ConfigAnalysisResult represents overall analysis results for a config file
type ConfigAnalysisResult struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConfigFileID uuid.UUID `json:"config_file_id" gorm:"type:uuid;not null;index"`
	CompanyID   uuid.UUID `json:"company_id" gorm:"type:uuid;not null;index"`

	// Analysis summary
	TotalFindings   int `json:"total_findings" gorm:"default:0"`
	CriticalFindings int `json:"critical_findings" gorm:"default:0"`
	HighFindings     int `json:"high_findings" gorm:"default:0"`
	MediumFindings   int `json:"medium_findings" gorm:"default:0"`
	LowFindings      int `json:"low_findings" gorm:"default:0"`
	InfoFindings     int `json:"info_findings" gorm:"default:0"`

	// Compliance scores
	ComplianceScores    datatypes.JSON `json:"compliance_scores,omitempty" gorm:"type:jsonb;default:'{}'"`
	OverallSecurityScore *float64      `json:"overall_security_score,omitempty" gorm:"type:decimal(5,2)"`

	// Analysis details
	AnalysisVersion  string         `json:"analysis_version,omitempty" gorm:"size:50"`
	StandardsChecked datatypes.JSON `json:"standards_checked,omitempty" gorm:"type:jsonb;default:'[]'"`
	ChecksPerformed  int            `json:"checks_performed" gorm:"default:0"`
	ChecksPassed     int            `json:"checks_passed" gorm:"default:0"`
	ChecksFailed     int            `json:"checks_failed" gorm:"default:0"`

	// Risk assessment
	OverallRiskScore float64 `json:"overall_risk_score" gorm:"type:decimal(3,2);default:0.0"`
	RiskLevel        string  `json:"risk_level,omitempty" gorm:"size:20"`

	// Report
	ReportPath   string `json:"report_path,omitempty" gorm:"size:500"`
	ReportFormat string `json:"report_format,omitempty" gorm:"size:20"`

	// Metadata
	AnalysisMetadata datatypes.JSON `json:"analysis_metadata,omitempty" gorm:"type:jsonb;default:'{}'"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	ConfigFile ConfigFile `json:"-" gorm:"foreignKey:ConfigFileID"`
	Company    Company    `json:"-" gorm:"foreignKey:CompanyID"`
}

// TableName specifies the table name
func (ConfigAnalysisResult) TableName() string {
	return "config_analysis_results"
}

// GetComplianceScore returns the compliance score for a specific framework
func (c *ConfigAnalysisResult) GetComplianceScore(framework string) float64 {
	// This would need to be implemented with proper JSON parsing
	// For now, return 0 as placeholder
	return 0.0
}

// Request/Response DTOs

// UploadConfigFileRequest represents a request to upload a config file
type UploadConfigFileRequest struct {
	DeviceType      string `form:"device_type" binding:"required"`
	Manufacturer    string `form:"manufacturer" binding:"required"`
	Model           string `form:"model"`
	FirmwareVersion string `form:"firmware_version"`
	DeviceName      string `form:"device_name"`
	DeviceLocation  string `form:"device_location"`
	ConfigType      string `form:"config_type" binding:"required"`
	ConfigFormat    string `form:"config_format"`
	Tags            []string `form:"tags"`
	Notes           string `form:"notes"`
}

// ListConfigFilesRequest represents filters for listing config files
type ListConfigFilesRequest struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
	Manufacturer string `form:"manufacturer"`
	DeviceType   string `form:"device_type"`
	Status       string `form:"status"`
	SortBy       string `form:"sort_by"`
	SortOrder    string `form:"sort_order"`
}

// ListConfigFindingsRequest represents filters for listing config findings
type ListConfigFindingsRequest struct {
	Page         int       `form:"page"`
	PageSize     int       `form:"page_size"`
	ConfigFileID *uuid.UUID `form:"config_file_id"`
	Severity     string    `form:"severity"`
	Category     string    `form:"category"`
	Status       string    `form:"status"`
	FindingType  string    `form:"finding_type"`
	SortBy       string    `form:"sort_by"`
	SortOrder    string    `form:"sort_order"`
}

// UpdateFindingStatusRequest represents a request to update finding status
type UpdateFindingStatusRequest struct {
	Status     string     `json:"status" binding:"required"`
	AssignedTo *uuid.UUID `json:"assigned_to,omitempty"`
	Notes      string     `json:"notes,omitempty"`
}

