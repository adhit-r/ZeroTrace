package constants

// Config Auditor Constants

// File size limits
const (
	MaxConfigFileSize = 10 * 1024 * 1024 // 10MB
	MinConfigFileSize = 1                 // Minimum 1 byte
)

// Pagination defaults
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
	MinPageSize     = 1
)

// Worker and queue configuration
const (
	DefaultWorkerCount     = 3
	DefaultQueueBufferSize = 100
)

// Status strings
const (
	// Parsing statuses
	StatusPending  = "pending"
	StatusParsing  = "parsing"
	StatusParsed   = "parsed"
	StatusFailed   = "failed"
	StatusPartial  = "partial"
	StatusAnalyzing = "analyzing"
	StatusCompleted = "completed"

	// Finding statuses
	StatusOpen          = "open"
	StatusAcknowledged  = "acknowledged"
	StatusMitigated     = "mitigated"
	StatusResolved      = "resolved"
	StatusFalsePositive = "false_positive"
	StatusAcceptedRisk  = "accepted_risk"

	// Severity levels
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
	SeverityInfo     = "info"
)

// Risk score thresholds
const (
	RiskThresholdCritical = 0.8
	RiskThresholdHigh    = 0.6
	RiskThresholdMedium  = 0.4
	RiskThresholdLow     = 0.2
)

// Severity weights for security score calculation
const (
	WeightCritical = 10.0
	WeightHigh     = 5.0
	WeightMedium   = 2.0
	WeightLow      = 1.0
	WeightInfo     = 0.5
)

// Risk scores by severity
const (
	RiskScoreCritical = 0.9
	RiskScoreHigh     = 0.7
	RiskScoreMedium   = 0.5
	RiskScoreLow      = 0.3
	RiskScoreInfo     = 0.1
)

// File path and storage
const (
	ConfigStoragePathTemplate = "configs"
)

// MIME types
const (
	MIMETypeTextPlain          = "text/plain"
	MIMETypeApplicationOctet   = "application/octet-stream"
	MIMETypeApplicationJSON    = "application/json"
	MIMETypeApplicationXML     = "application/xml"
)

// Analysis version
const (
	AnalysisVersion = "1.0"
)

// Default user accounts that should be flagged
var DefaultUserAccounts = []string{"admin", "root", "cisco", "user", "guest"}

// Device types (enum values)
var ValidDeviceTypes = []string{
	"firewall", "router", "switch", "load_balancer", "waf",
	"ids", "ips", "vpn_gateway", "wireless_controller", "other",
}

// Config types (enum values)
var ValidConfigTypes = []string{
	"running_config", "startup_config", "backup_config", "export_config", "other",
}

// Finding statuses (enum values)
var ValidFindingStatuses = []string{
	"open", "acknowledged", "mitigated", "resolved", "false_positive", "accepted_risk",
}

// Regex complexity limits (for ReDoS protection)
const (
	MaxRegexPatternLength = 1000
)

// Score calculation constants
const (
	MaxSecurityScore = 100.0
	MinSecurityScore = 0.0
)

