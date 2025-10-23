package scanner

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// PrivacyScanner handles privacy and data protection scanning
type PrivacyScanner struct {
	config *config.Config
}

// PrivacyFinding represents a privacy security finding
type PrivacyFinding struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`     // pii, gdpr, ccpa, data_retention, consent
	Severity      string                 `json:"severity"` // critical, high, medium, low
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	DataType      string                 `json:"data_type,omitempty"`
	Location      string                 `json:"location,omitempty"`
	CurrentValue  string                 `json:"current_value,omitempty"`
	RequiredValue string                 `json:"required_value,omitempty"`
	Remediation   string                 `json:"remediation"`
	DiscoveredAt  time.Time              `json:"discovered_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// PIIInfo represents personally identifiable information
type PIIInfo struct {
	Type            string                 `json:"type"` // email, phone, ssn, credit_card, address
	Value           string                 `json:"value"`
	Location        string                 `json:"location"`
	IsEncrypted     bool                   `json:"is_encrypted"`
	IsMasked        bool                   `json:"is_masked"`
	AccessLevel     string                 `json:"access_level"` // public, internal, restricted, confidential
	RetentionPeriod string                 `json:"retention_period"`
	LastAccessed    time.Time              `json:"last_accessed"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// ComplianceInfo represents compliance framework information
type ComplianceInfo struct {
	Framework    string                 `json:"framework"` // GDPR, CCPA, HIPAA, SOX
	Status       string                 `json:"status"`    // compliant, non_compliant, partial
	Score        float64                `json:"score"`
	Requirements []string               `json:"requirements"`
	Violations   []string               `json:"violations"`
	LastAudit    time.Time              `json:"last_audit"`
	NextAudit    time.Time              `json:"next_audit"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// DataRetentionInfo represents data retention information
type DataRetentionInfo struct {
	DataType        string                 `json:"data_type"`
	Location        string                 `json:"location"`
	RetentionPeriod string                 `json:"retention_period"`
	IsExpired       bool                   `json:"is_expired"`
	ExpiryDate      time.Time              `json:"expiry_date"`
	IsDeleted       bool                   `json:"is_deleted"`
	DeletionDate    time.Time              `json:"deletion_date"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// NewPrivacyScanner creates a new privacy security scanner
func NewPrivacyScanner(cfg *config.Config) *PrivacyScanner {
	return &PrivacyScanner{
		config: cfg,
	}
}

// Scan performs comprehensive privacy and data protection scanning
func (ps *PrivacyScanner) Scan() ([]PrivacyFinding, []PIIInfo, []ComplianceInfo, []DataRetentionInfo, error) {
	var findings []PrivacyFinding
	var piiData []PIIInfo
	var compliance []ComplianceInfo
	var retention []DataRetentionInfo

	// Scan for PII
	discoveredPII := ps.scanPII()
	piiData = append(piiData, discoveredPII...)

	// Scan PII for issues
	for _, pii := range discoveredPII {
		piiFindings := ps.scanPIIIssues(pii)
		findings = append(findings, piiFindings...)
	}

	// Scan compliance
	discoveredCompliance := ps.scanCompliance()
	compliance = append(compliance, discoveredCompliance...)

	// Scan compliance for issues
	for _, comp := range discoveredCompliance {
		compFindings := ps.scanComplianceIssues(comp)
		findings = append(findings, compFindings...)
	}

	// Scan data retention
	discoveredRetention := ps.scanDataRetention()
	retention = append(retention, discoveredRetention...)

	// Scan retention for issues
	for _, ret := range discoveredRetention {
		retFindings := ps.scanRetentionIssues(ret)
		findings = append(findings, retFindings...)
	}

	return findings, piiData, compliance, retention, nil
}

// scanPII scans for personally identifiable information
func (ps *PrivacyScanner) scanPII() []PIIInfo {
	var piiData []PIIInfo

	// Scan for email addresses
	emails := ps.scanForEmails()
	piiData = append(piiData, emails...)

	// Scan for phone numbers
	phones := ps.scanForPhoneNumbers()
	piiData = append(piiData, phones...)

	// Scan for SSNs
	ssns := ps.scanForSSNs()
	piiData = append(piiData, ssns...)

	// Scan for credit card numbers
	creditCards := ps.scanForCreditCards()
	piiData = append(piiData, creditCards...)

	// Scan for addresses
	addresses := ps.scanForAddresses()
	piiData = append(piiData, addresses...)

	return piiData
}

// scanForEmails scans for email addresses
func (ps *PrivacyScanner) scanForEmails() []PIIInfo {
	var piiData []PIIInfo

	// Email regex pattern
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// Scan common locations for emails
	locations := []string{
		"/var/log/",
		"/tmp/",
		"/home/",
		"/Users/",
	}

	for _, location := range locations {
		if ps.isDirectory(location) {
			emails := ps.scanDirectoryForPattern(location, emailPattern)
			for _, email := range emails {
				pii := PIIInfo{
					Type:            "email",
					Value:           email,
					Location:        location,
					IsEncrypted:     false,
					IsMasked:        false,
					AccessLevel:     "internal",
					RetentionPeriod: "unknown",
					LastAccessed:    time.Now(),
					Metadata: map[string]interface{}{
						"pattern":  "email",
						"location": location,
					},
				}
				piiData = append(piiData, pii)
			}
		}
	}

	return piiData
}

// scanForPhoneNumbers scans for phone numbers
func (ps *PrivacyScanner) scanForPhoneNumbers() []PIIInfo {
	var piiData []PIIInfo

	// Phone number regex pattern
	phonePattern := regexp.MustCompile(`(\+?1[-.\s]?)?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})`)

	// Scan common locations for phone numbers
	locations := []string{
		"/var/log/",
		"/tmp/",
		"/home/",
		"/Users/",
	}

	for _, location := range locations {
		if ps.isDirectory(location) {
			phones := ps.scanDirectoryForPattern(location, phonePattern)
			for _, phone := range phones {
				pii := PIIInfo{
					Type:            "phone",
					Value:           phone,
					Location:        location,
					IsEncrypted:     false,
					IsMasked:        false,
					AccessLevel:     "internal",
					RetentionPeriod: "unknown",
					LastAccessed:    time.Now(),
					Metadata: map[string]interface{}{
						"pattern":  "phone",
						"location": location,
					},
				}
				piiData = append(piiData, pii)
			}
		}
	}

	return piiData
}

// scanForSSNs scans for Social Security Numbers
func (ps *PrivacyScanner) scanForSSNs() []PIIInfo {
	var piiData []PIIInfo

	// SSN regex pattern
	ssnPattern := regexp.MustCompile(`\b\d{3}-?\d{2}-?\d{4}\b`)

	// Scan common locations for SSNs
	locations := []string{
		"/var/log/",
		"/tmp/",
		"/home/",
		"/Users/",
	}

	for _, location := range locations {
		if ps.isDirectory(location) {
			ssns := ps.scanDirectoryForPattern(location, ssnPattern)
			for _, ssn := range ssns {
				pii := PIIInfo{
					Type:            "ssn",
					Value:           ssn,
					Location:        location,
					IsEncrypted:     false,
					IsMasked:        false,
					AccessLevel:     "restricted",
					RetentionPeriod: "unknown",
					LastAccessed:    time.Now(),
					Metadata: map[string]interface{}{
						"pattern":  "ssn",
						"location": location,
					},
				}
				piiData = append(piiData, pii)
			}
		}
	}

	return piiData
}

// scanForCreditCards scans for credit card numbers
func (ps *PrivacyScanner) scanForCreditCards() []PIIInfo {
	var piiData []PIIInfo

	// Credit card regex pattern
	ccPattern := regexp.MustCompile(`\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b`)

	// Scan common locations for credit cards
	locations := []string{
		"/var/log/",
		"/tmp/",
		"/home/",
		"/Users/",
	}

	for _, location := range locations {
		if ps.isDirectory(location) {
			creditCards := ps.scanDirectoryForPattern(location, ccPattern)
			for _, cc := range creditCards {
				pii := PIIInfo{
					Type:            "credit_card",
					Value:           cc,
					Location:        location,
					IsEncrypted:     false,
					IsMasked:        false,
					AccessLevel:     "restricted",
					RetentionPeriod: "unknown",
					LastAccessed:    time.Now(),
					Metadata: map[string]interface{}{
						"pattern":  "credit_card",
						"location": location,
					},
				}
				piiData = append(piiData, pii)
			}
		}
	}

	return piiData
}

// scanForAddresses scans for addresses
func (ps *PrivacyScanner) scanForAddresses() []PIIInfo {
	var piiData []PIIInfo

	// Address regex pattern (simplified)
	addressPattern := regexp.MustCompile(`\b\d+\s+[A-Za-z\s]+(?:Street|St|Avenue|Ave|Road|Rd|Boulevard|Blvd|Drive|Dr|Lane|Ln|Way|Circle|Cir)\b`)

	// Scan common locations for addresses
	locations := []string{
		"/var/log/",
		"/tmp/",
		"/home/",
		"/Users/",
	}

	for _, location := range locations {
		if ps.isDirectory(location) {
			addresses := ps.scanDirectoryForPattern(location, addressPattern)
			for _, address := range addresses {
				pii := PIIInfo{
					Type:            "address",
					Value:           address,
					Location:        location,
					IsEncrypted:     false,
					IsMasked:        false,
					AccessLevel:     "internal",
					RetentionPeriod: "unknown",
					LastAccessed:    time.Now(),
					Metadata: map[string]interface{}{
						"pattern":  "address",
						"location": location,
					},
				}
				piiData = append(piiData, pii)
			}
		}
	}

	return piiData
}

// scanDirectoryForPattern scans a directory for a regex pattern
func (ps *PrivacyScanner) scanDirectoryForPattern(location string, pattern *regexp.Regexp) []string {
	var matches []string

	// This would involve recursively scanning files in the directory
	// For now, return a placeholder
	return matches
}

// isDirectory checks if a path is a directory
func (ps *PrivacyScanner) isDirectory(path string) bool {
	cmd := exec.Command("test", "-d", path)
	err := cmd.Run()
	return err == nil
}

// scanCompliance scans for compliance frameworks
func (ps *PrivacyScanner) scanCompliance() []ComplianceInfo {
	var compliance []ComplianceInfo

	// GDPR compliance
	gdpr := ComplianceInfo{
		Framework:    "GDPR",
		Status:       "partial",
		Score:        0.6,
		Requirements: []string{"data_protection", "consent", "right_to_be_forgotten", "data_portability"},
		Violations:   []string{"missing_consent", "no_data_portability"},
		LastAudit:    time.Now().AddDate(0, -1, 0),
		NextAudit:    time.Now().AddDate(0, 1, 0),
		Metadata: map[string]interface{}{
			"region": "EU",
			"scope":  "personal_data",
		},
	}
	compliance = append(compliance, gdpr)

	// CCPA compliance
	ccpa := ComplianceInfo{
		Framework:    "CCPA",
		Status:       "partial",
		Score:        0.7,
		Requirements: []string{"consumer_rights", "opt_out", "data_disclosure"},
		Violations:   []string{"missing_opt_out"},
		LastAudit:    time.Now().AddDate(0, -2, 0),
		NextAudit:    time.Now().AddDate(0, 1, 0),
		Metadata: map[string]interface{}{
			"region": "California",
			"scope":  "consumer_data",
		},
	}
	compliance = append(compliance, ccpa)

	// HIPAA compliance
	hipaa := ComplianceInfo{
		Framework:    "HIPAA",
		Status:       "non_compliant",
		Score:        0.4,
		Requirements: []string{"administrative_safeguards", "physical_safeguards", "technical_safeguards"},
		Violations:   []string{"missing_encryption", "no_access_controls"},
		LastAudit:    time.Now().AddDate(0, -3, 0),
		NextAudit:    time.Now().AddDate(0, 1, 0),
		Metadata: map[string]interface{}{
			"region": "US",
			"scope":  "healthcare_data",
		},
	}
	compliance = append(compliance, hipaa)

	return compliance
}

// scanDataRetention scans for data retention policies
func (ps *PrivacyScanner) scanDataRetention() []DataRetentionInfo {
	var retention []DataRetentionInfo

	// Scan for expired data
	expiredData := DataRetentionInfo{
		DataType:        "personal_data",
		Location:        "/var/data/personal/",
		RetentionPeriod: "7 years",
		IsExpired:       true,
		ExpiryDate:      time.Now().AddDate(-1, 0, 0),
		IsDeleted:       false,
		DeletionDate:    time.Time{},
		Metadata: map[string]interface{}{
			"data_classification": "sensitive",
		},
	}
	retention = append(retention, expiredData)

	// Scan for data without retention policy
	noPolicyData := DataRetentionInfo{
		DataType:        "log_data",
		Location:        "/var/log/",
		RetentionPeriod: "unknown",
		IsExpired:       false,
		ExpiryDate:      time.Time{},
		IsDeleted:       false,
		DeletionDate:    time.Time{},
		Metadata: map[string]interface{}{
			"data_classification": "internal",
		},
	}
	retention = append(retention, noPolicyData)

	return retention
}

// scanPIIIssues scans PII for security issues
func (ps *PrivacyScanner) scanPIIIssues(pii PIIInfo) []PrivacyFinding {
	var findings []PrivacyFinding

	// Check for unencrypted PII
	if !pii.IsEncrypted {
		finding := PrivacyFinding{
			ID:            uuid.New().String(),
			Type:          "pii",
			Severity:      "critical",
			Title:         "Unencrypted PII",
			Description:   fmt.Sprintf("PII of type %s is not encrypted", pii.Type),
			DataType:      pii.Type,
			Location:      pii.Location,
			CurrentValue:  "unencrypted",
			RequiredValue: "encrypted",
			Remediation:   "Encrypt PII data",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"pii_type":  pii.Type,
				"location":  pii.Location,
				"encrypted": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for unmasked PII
	if !pii.IsMasked {
		finding := PrivacyFinding{
			ID:            uuid.New().String(),
			Type:          "pii",
			Severity:      "high",
			Title:         "Unmasked PII",
			Description:   fmt.Sprintf("PII of type %s is not masked", pii.Type),
			DataType:      pii.Type,
			Location:      pii.Location,
			CurrentValue:  "unmasked",
			RequiredValue: "masked",
			Remediation:   "Mask PII data",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"pii_type": pii.Type,
				"location": pii.Location,
				"masked":   false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for public access
	if pii.AccessLevel == "public" {
		finding := PrivacyFinding{
			ID:            uuid.New().String(),
			Type:          "pii",
			Severity:      "critical",
			Title:         "Public PII Access",
			Description:   fmt.Sprintf("PII of type %s has public access", pii.Type),
			DataType:      pii.Type,
			Location:      pii.Location,
			CurrentValue:  "public",
			RequiredValue: "restricted",
			Remediation:   "Restrict PII access",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"pii_type":     pii.Type,
				"location":     pii.Location,
				"access_level": pii.AccessLevel,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanComplianceIssues scans compliance for issues
func (ps *PrivacyScanner) scanComplianceIssues(compliance ComplianceInfo) []PrivacyFinding {
	var findings []PrivacyFinding

	// Check for non-compliance
	if compliance.Status == "non_compliant" {
		finding := PrivacyFinding{
			ID:            uuid.New().String(),
			Type:          "gdpr",
			Severity:      "critical",
			Title:         fmt.Sprintf("%s Non-Compliance", compliance.Framework),
			Description:   fmt.Sprintf("%s compliance status is non-compliant", compliance.Framework),
			CurrentValue:  compliance.Status,
			RequiredValue: "compliant",
			Remediation:   fmt.Sprintf("Implement %s compliance requirements", compliance.Framework),
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"framework": compliance.Framework,
				"status":    compliance.Status,
				"score":     compliance.Score,
			},
		}
		findings = append(findings, finding)
	}

	// Check for low compliance score
	if compliance.Score < 0.7 {
		finding := PrivacyFinding{
			ID:            uuid.New().String(),
			Type:          "gdpr",
			Severity:      "high",
			Title:         fmt.Sprintf("Low %s Compliance Score", compliance.Framework),
			Description:   fmt.Sprintf("%s compliance score is %.2f", compliance.Framework, compliance.Score),
			CurrentValue:  fmt.Sprintf("%.2f", compliance.Score),
			RequiredValue: "0.7+",
			Remediation:   fmt.Sprintf("Improve %s compliance", compliance.Framework),
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"framework": compliance.Framework,
				"score":     compliance.Score,
			},
		}
		findings = append(findings, finding)
	}

	// Check for violations
	if len(compliance.Violations) > 0 {
		finding := PrivacyFinding{
			ID:           uuid.New().String(),
			Type:         "gdpr",
			Severity:     "high",
			Title:        fmt.Sprintf("%s Violations", compliance.Framework),
			Description:  fmt.Sprintf("%s has %d violations: %s", compliance.Framework, len(compliance.Violations), strings.Join(compliance.Violations, ", ")),
			Remediation:  fmt.Sprintf("Address %s violations", compliance.Framework),
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"framework":  compliance.Framework,
				"violations": compliance.Violations,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanRetentionIssues scans data retention for issues
func (ps *PrivacyScanner) scanRetentionIssues(retention DataRetentionInfo) []PrivacyFinding {
	var findings []PrivacyFinding

	// Check for expired data
	if retention.IsExpired && !retention.IsDeleted {
		finding := PrivacyFinding{
			ID:            uuid.New().String(),
			Type:          "data_retention",
			Severity:      "high",
			Title:         "Expired Data Not Deleted",
			Description:   fmt.Sprintf("Data of type %s has expired but not been deleted", retention.DataType),
			DataType:      retention.DataType,
			Location:      retention.Location,
			CurrentValue:  "expired_not_deleted",
			RequiredValue: "deleted",
			Remediation:   "Delete expired data",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"data_type": retention.DataType,
				"location":  retention.Location,
				"expired":   retention.IsExpired,
				"deleted":   retention.IsDeleted,
			},
		}
		findings = append(findings, finding)
	}

	// Check for missing retention policy
	if retention.RetentionPeriod == "unknown" {
		finding := PrivacyFinding{
			ID:            uuid.New().String(),
			Type:          "data_retention",
			Severity:      "medium",
			Title:         "Missing Data Retention Policy",
			Description:   fmt.Sprintf("Data of type %s has no retention policy", retention.DataType),
			DataType:      retention.DataType,
			Location:      retention.Location,
			CurrentValue:  "no_policy",
			RequiredValue: "defined_policy",
			Remediation:   "Define data retention policy",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"data_type":        retention.DataType,
				"location":         retention.Location,
				"retention_period": retention.RetentionPeriod,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}
