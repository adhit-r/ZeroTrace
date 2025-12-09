package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"zerotrace/api/internal/constants"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/google/uuid"
)

// ConfigAnalyzerService handles configuration analysis against standards
type ConfigAnalyzerService struct {
	configFileRepo      *repository.ConfigFileRepository
	configFindingRepo   *repository.ConfigFindingRepository
	configStandardRepo  *repository.ConfigStandardRepository
	configAnalysisRepo  *repository.ConfigAnalysisRepository
}

// NewConfigAnalyzerService creates a new config analyzer service
func NewConfigAnalyzerService(
	configFileRepo *repository.ConfigFileRepository,
	configFindingRepo *repository.ConfigFindingRepository,
	configStandardRepo *repository.ConfigStandardRepository,
	configAnalysisRepo *repository.ConfigAnalysisRepository,
) *ConfigAnalyzerService {
	return &ConfigAnalyzerService{
		configFileRepo:     configFileRepo,
		configFindingRepo:   configFindingRepo,
		configStandardRepo:  configStandardRepo,
		configAnalysisRepo:  configAnalysisRepo,
	}
}

// AnalyzeConfigFile analyzes a configuration file against standards
func (s *ConfigAnalyzerService) AnalyzeConfigFile(configFileID uuid.UUID) error {
	// Get config file
	configFile, err := s.configFileRepo.GetByID(configFileID)
	if err != nil {
		return fmt.Errorf("failed to get config file: %w", err)
	}

	// Update analysis status
	err = s.configFileRepo.UpdateAnalysisStatus(configFileID, constants.StatusAnalyzing)
	if err != nil {
		return err
	}

	// Check if parsed
	if configFile.ParsingStatus != constants.StatusParsed {
		return fmt.Errorf("config file must be parsed before analysis")
	}

	// Get standards for this manufacturer/device type
	standards, err := s.configStandardRepo.GetByManufacturer(configFile.Manufacturer, configFile.DeviceType)
	if err != nil {
		return fmt.Errorf("failed to get standards: %w", err)
	}

	// Parse the parsed_data JSONB
	var parsedConfig map[string]interface{}
	err = json.Unmarshal(configFile.ParsedData, &parsedConfig)
	if err != nil {
		return fmt.Errorf("failed to parse config data for file %s: %w", configFileID, err)
	}

	// Check against standards
	findings, err := s.CheckAgainstStandards(parsedConfig, standards, configFile)
	if err != nil {
		s.configFileRepo.UpdateAnalysisStatus(configFileID, constants.StatusFailed)
		return fmt.Errorf("failed to check standards: %w", err)
	}

	// Save findings
	if len(findings) > 0 {
		err = s.configFindingRepo.CreateBatch(findings)
		if err != nil {
			return fmt.Errorf("failed to save findings: %w", err)
		}
	}

	// Calculate scores
	complianceScores := s.CalculateComplianceScores(findings, standards)
	securityScore := s.CalculateSecurityScore(findings)

	// Generate analysis result
	analysisResult, err := s.GenerateAnalysisResult(configFileID, configFile.CompanyID, findings, complianceScores, securityScore, standards)
	if err != nil {
		return fmt.Errorf("failed to generate analysis result: %w", err)
	}

	// Save analysis result
	err = s.configAnalysisRepo.Create(analysisResult)
	if err != nil {
		return fmt.Errorf("failed to save analysis result: %w", err)
	}

	// Update analysis status to completed
	err = s.configFileRepo.UpdateAnalysisStatus(configFileID, constants.StatusCompleted)
	if err != nil {
		return err
	}

	return nil
}

// CheckAgainstStandards checks parsed config against standards
func (s *ConfigAnalyzerService) CheckAgainstStandards(
	parsedConfig map[string]interface{},
	standards []models.ConfigStandard,
	configFile *models.ConfigFile,
) ([]models.ConfigFinding, error) {
	var findings []models.ConfigFinding

	// Convert config to string for pattern matching
	configContent := string(configFile.FileContent)
	configLines := strings.Split(configContent, "\n")

	for _, standard := range standards {
		if !standard.IsActive() {
			continue
		}

		violated, lineNumbers, snippet := s.checkStandard(parsedConfig, standard, configContent, configLines)
		if violated {
			finding := models.ConfigFinding{
				ConfigFileID:        configFile.ID,
				CompanyID:          configFile.CompanyID,
				FindingType:        "compliance_violation",
				Severity:           standard.DefaultSeverity,
				Category:           standard.Category,
				Title:              standard.RequirementTitle,
				Description:        standard.RequirementDescription,
				AffectedComponent:  standard.CheckConfigPath,
				ConfigSnippet:      snippet,
				LineNumbers:        lineNumbers,
				StandardID:         &standard.ID,
				ComplianceFrameworks: standard.ComplianceFrameworks,
				Remediation:         standard.RemediationGuidance,
				RemediationSteps:    s.parseRemediationSteps(standard.RemediationGuidance),
				RemediationPriority: standard.Priority,
				RiskScore:          s.calculateRiskScore(standard.DefaultSeverity),
				Status:             constants.StatusOpen,
			}

			// Set exploitability and impact based on severity
			switch standard.DefaultSeverity {
			case constants.SeverityCritical:
				finding.Exploitability = "high"
				finding.Impact = "critical"
			case constants.SeverityHigh:
				finding.Exploitability = "medium"
				finding.Impact = "high"
			case constants.SeverityMedium:
				finding.Exploitability = "low"
				finding.Impact = "medium"
			default:
				finding.Exploitability = "low"
				finding.Impact = "low"
			}

			findings = append(findings, finding)
		}
	}

	// Also perform basic security checks
	basicFindings := s.performBasicSecurityChecks(parsedConfig, configFile, configContent, configLines)
	findings = append(findings, basicFindings...)

	return findings, nil
}

// checkStandard checks if a standard is violated
func (s *ConfigAnalyzerService) checkStandard(
	parsedConfig map[string]interface{},
	standard models.ConfigStandard,
	configContent string,
	configLines []string,
) (bool, []byte, string) {
	var lineNumbers []int
	var snippet string

	switch standard.CheckType {
	case "presence":
		// Check if required config is present
		value := s.getConfigValue(parsedConfig, standard.CheckConfigPath)
		if value == nil {
			// Find line numbers where this should be
			for i, line := range configLines {
				if strings.Contains(strings.ToLower(line), strings.ToLower(standard.CheckConfigPath)) {
					lineNumbers = append(lineNumbers, i+1)
				}
			}
			return true, s.intArrayToJSON(lineNumbers), snippet
		}
		return false, nil, ""

	case "absence":
		// Check if prohibited config is absent
		value := s.getConfigValue(parsedConfig, standard.CheckConfigPath)
		if value != nil {
			// Find line numbers
			for i, line := range configLines {
				if strings.Contains(strings.ToLower(line), strings.ToLower(standard.CheckConfigPath)) {
					lineNumbers = append(lineNumbers, i+1)
					snippet += line + "\n"
				}
			}
			return true, s.intArrayToJSON(lineNumbers), snippet
		}
		return false, nil, ""

	case "pattern_match":
		// Check if config matches pattern
		if standard.CheckPattern != "" {
			// Validate regex pattern complexity to prevent ReDoS
			if len(standard.CheckPattern) > constants.MaxRegexPatternLength {
				// Pattern too long, skip to prevent ReDoS
				return false, nil, ""
			}
			regex, err := regexp.Compile(standard.CheckPattern)
			if err != nil {
				// Invalid regex pattern, log and skip
				return false, nil, ""
			}
			matches := regex.FindAllString(configContent, -1)
			if len(matches) == 0 {
				return true, nil, "" // Pattern not found when it should be
			}
			// Check if matches are correct (with error handling)
			if standard.ExpectedValue != "" {
				expectedRegex, err := regexp.Compile(standard.ExpectedValue)
				if err != nil {
					// Invalid expected value regex, skip this check
					return false, nil, ""
				}
				for _, match := range matches {
					if !expectedRegex.MatchString(match) {
						// Find line numbers
						for i, line := range configLines {
							if strings.Contains(line, match) {
								lineNumbers = append(lineNumbers, i+1)
								snippet += line + "\n"
							}
						}
						return true, s.intArrayToJSON(lineNumbers), snippet
					}
				}
			}
		}
		return false, nil, ""

	case "value_match":
		// Check if value matches expected
		value := s.getConfigValue(parsedConfig, standard.CheckConfigPath)
		if value != nil {
			valueStr := fmt.Sprintf("%v", value)
			if valueStr != standard.ExpectedValue {
				// Find line numbers
				for i, line := range configLines {
					if strings.Contains(line, standard.CheckConfigPath) {
						lineNumbers = append(lineNumbers, i+1)
						snippet += line + "\n"
					}
				}
				return true, s.intArrayToJSON(lineNumbers), snippet
			}
		}
		return false, nil, ""

	default:
		return false, nil, ""
	}
}

// getConfigValue gets a value from parsed config using JSON path
func (s *ConfigAnalyzerService) getConfigValue(config map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := config

	for i, part := range parts {
		if i == len(parts)-1 {
			return current[part]
		}
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}
	return nil
}

// performBasicSecurityChecks performs basic security checks
func (s *ConfigAnalyzerService) performBasicSecurityChecks(
	parsedConfig map[string]interface{},
	configFile *models.ConfigFile,
	configContent string,
	configLines []string,
) []models.ConfigFinding {
	var findings []models.ConfigFinding

	// Check for default credentials
	if users, ok := parsedConfig["user_accounts"].([]map[string]interface{}); ok {
		for _, user := range users {
			if username, ok := user["username"].(string); ok {
				for _, defaultUser := range constants.DefaultUserAccounts {
					if strings.ToLower(username) == defaultUser {
						finding := models.ConfigFinding{
							ConfigFileID:       configFile.ID,
							CompanyID:          configFile.CompanyID,
							FindingType:        "default_credentials",
							Severity:           constants.SeverityHigh,
							Category:           "authentication",
							Title:              "Default User Account Detected",
							Description:        fmt.Sprintf("Default user account '%s' is present in configuration", username),
							AffectedComponent:  fmt.Sprintf("user: %s", username),
							Remediation:        "Remove default user accounts or change default passwords",
							RemediationPriority: "high",
							RiskScore:          constants.RiskScoreHigh,
							Status:             constants.StatusOpen,
						}
						findings = append(findings, finding)
					}
				}
			}
		}
	}

	// Check for insecure protocols (Telnet, FTP, SNMP v1/v2)
	configLower := strings.ToLower(configContent)
	if strings.Contains(configLower, "telnet") && !strings.Contains(configLower, "no telnet") {
		finding := models.ConfigFinding{
			ConfigFileID:       configFile.ID,
			CompanyID:          configFile.CompanyID,
			FindingType:        "insecure_protocol",
			Severity:           constants.SeverityHigh,
			Category:           "network",
			Title:              "Telnet Protocol Enabled",
			Description:        "Telnet is an insecure protocol that transmits data in plaintext",
			Remediation:        "Disable Telnet and use SSH instead",
			RemediationPriority: "high",
			RiskScore:          constants.RiskScoreHigh,
			Status:             constants.StatusOpen,
		}
		findings = append(findings, finding)
	}

	// Check for weak encryption
	if crypto, ok := parsedConfig["crypto"].(map[string]interface{}); ok {
		if config, ok := crypto["config"].([]string); ok {
			for _, line := range config {
				if strings.Contains(strings.ToLower(line), "md5") || strings.Contains(strings.ToLower(line), "des") {
					finding := models.ConfigFinding{
						ConfigFileID:       configFile.ID,
						CompanyID:          configFile.CompanyID,
						FindingType:        "weak_cipher",
						Severity:           constants.SeverityMedium,
						Category:           "encryption",
						Title:              "Weak Encryption Algorithm Detected",
						Description:        fmt.Sprintf("Weak encryption algorithm found: %s", line),
						ConfigSnippet:      line,
						Remediation:        "Use strong encryption algorithms (AES-256, SHA-256)",
						RemediationPriority: "medium",
						RiskScore:          constants.RiskScoreMedium,
						Status:             constants.StatusOpen,
					}
					findings = append(findings, finding)
				}
			}
		}
	}

	return findings
}

// CalculateComplianceScores calculates compliance scores by framework
func (s *ConfigAnalyzerService) CalculateComplianceScores(
	findings []models.ConfigFinding,
	standards []models.ConfigStandard,
) map[string]float64 {
	scores := make(map[string]float64)
	frameworkTotals := make(map[string]int)
	frameworkPassed := make(map[string]int)

	// Count total checks and passed checks per framework
	for _, standard := range standards {
		if !standard.IsActive() {
			continue
		}

		// Get compliance frameworks from standard
		var frameworks []string
		if standard.ComplianceFrameworks != nil {
			json.Unmarshal(standard.ComplianceFrameworks, &frameworks)
		}

		for _, framework := range frameworks {
			frameworkTotals[framework]++

			// Check if this standard has a violation
			violated := false
			for _, finding := range findings {
				if finding.StandardID != nil && *finding.StandardID == standard.ID {
					violated = true
					break
				}
			}

			if !violated {
				frameworkPassed[framework]++
			}
		}
	}

	// Calculate scores
	for framework, total := range frameworkTotals {
		passed := frameworkPassed[framework]
		if total > 0 {
			scores[framework] = float64(passed) / float64(total) * 100.0
		}
	}

	return scores
}

// CalculateSecurityScore calculates overall security score (0-100)
func (s *ConfigAnalyzerService) CalculateSecurityScore(findings []models.ConfigFinding) float64 {
	if len(findings) == 0 {
		return constants.MaxSecurityScore
	}

	// Weight findings by severity
	totalWeight := 0.0
	weightedScore := 0.0

	severityWeights := map[string]float64{
		constants.SeverityCritical: constants.WeightCritical,
		constants.SeverityHigh:     constants.WeightHigh,
		constants.SeverityMedium:   constants.WeightMedium,
		constants.SeverityLow:      constants.WeightLow,
		constants.SeverityInfo:     constants.WeightInfo,
	}

	for _, finding := range findings {
		weight := severityWeights[finding.Severity]
		totalWeight += weight
		weightedScore += weight * (1.0 - finding.RiskScore)
	}

	if totalWeight == 0 {
		return constants.MaxSecurityScore
	}

	score := (weightedScore / totalWeight) * 100.0
	if score < constants.MinSecurityScore {
		score = constants.MinSecurityScore
	}
	if score > constants.MaxSecurityScore {
		score = constants.MaxSecurityScore
	}

	return score
}

// GenerateAnalysisResult generates analysis result from findings
func (s *ConfigAnalyzerService) GenerateAnalysisResult(
	configFileID uuid.UUID,
	companyID uuid.UUID,
	findings []models.ConfigFinding,
	complianceScores map[string]float64,
	securityScore float64,
	standards []models.ConfigStandard,
) (*models.ConfigAnalysisResult, error) {
	// Count findings by severity
	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0
	infoCount := 0

	for _, finding := range findings {
		switch finding.Severity {
		case constants.SeverityCritical:
			criticalCount++
		case constants.SeverityHigh:
			highCount++
		case constants.SeverityMedium:
			mediumCount++
		case constants.SeverityLow:
			lowCount++
		case constants.SeverityInfo:
			infoCount++
		}
	}

	// Calculate overall risk score
	overallRiskScore := 0.0
	if len(findings) > 0 {
		for _, finding := range findings {
			overallRiskScore += finding.RiskScore
		}
		overallRiskScore = overallRiskScore / float64(len(findings))
	}

	// Determine risk level
	riskLevel := constants.SeverityLow
	if overallRiskScore >= constants.RiskThresholdCritical {
		riskLevel = constants.SeverityCritical
	} else if overallRiskScore >= constants.RiskThresholdHigh {
		riskLevel = constants.SeverityHigh
	} else if overallRiskScore >= constants.RiskThresholdMedium {
		riskLevel = constants.SeverityMedium
	}

	// Get standards checked
	standardsChecked := make([]string, 0, len(standards))
	for _, std := range standards {
		standardsChecked = append(standardsChecked, std.RequirementID)
	}
	standardsJSON, _ := json.Marshal(standardsChecked)

	// Convert compliance scores to JSONB
	complianceJSON, _ := json.Marshal(complianceScores)

	result := &models.ConfigAnalysisResult{
		ConfigFileID:        configFileID,
		CompanyID:          companyID,
		TotalFindings:      len(findings),
		CriticalFindings:   criticalCount,
		HighFindings:       highCount,
		MediumFindings:     mediumCount,
		LowFindings:        lowCount,
		InfoFindings:       infoCount,
		ComplianceScores:   complianceJSON,
		OverallSecurityScore: &securityScore,
		ChecksPerformed:    len(standards),
		ChecksPassed:       len(standards) - len(findings),
		ChecksFailed:       len(findings),
		OverallRiskScore:   overallRiskScore,
		RiskLevel:          riskLevel,
		StandardsChecked:   standardsJSON,
		AnalysisVersion:    "1.0",
	}

	return result, nil
}

// Helper functions

func (s *ConfigAnalyzerService) calculateRiskScore(severity string) float64 {
	switch severity {
	case constants.SeverityCritical:
		return constants.RiskScoreCritical
	case constants.SeverityHigh:
		return constants.RiskScoreHigh
	case constants.SeverityMedium:
		return constants.RiskScoreMedium
	case constants.SeverityLow:
		return constants.RiskScoreLow
	default:
		return constants.RiskScoreInfo
	}
}

func (s *ConfigAnalyzerService) parseRemediationSteps(guidance string) []byte {
	// Simple parsing - split by newlines or numbered lists
	steps := []string{}
	lines := strings.Split(guidance, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && (strings.HasPrefix(line, "1.") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*")) {
			steps = append(steps, line)
		}
	}
	if len(steps) == 0 {
		steps = []string{guidance}
	}
	json, _ := json.Marshal(steps)
	return json
}

func (s *ConfigAnalyzerService) intArrayToJSON(arr []int) []byte {
	json, _ := json.Marshal(arr)
	return json
}

