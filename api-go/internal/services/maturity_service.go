package services

import (
	"fmt"
	"math"
	"strings"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MaturityService handles security maturity scoring and industry benchmarking
type MaturityService struct {
	db *gorm.DB
}

// NewMaturityService creates a new MaturityService
func NewMaturityService(db *gorm.DB) *MaturityService {
	return &MaturityService{db: db}
}

// MaturityScore represents a comprehensive security maturity score
type MaturityScore struct {
	OrganizationID     uuid.UUID                 `json:"organization_id"`
	ScoreID            string                    `json:"score_id"`
	OverallScore       float64                   `json:"overall_score"`
	MaturityLevel      string                    `json:"maturity_level"`
	DimensionScores    map[string]DimensionScore `json:"dimension_scores"`
	IndustryBenchmark  IndustryBenchmark         `json:"industry_benchmark"`
	PeerComparison     PeerComparison            `json:"peer_comparison"`
	ImprovementRoadmap []ImprovementItem         `json:"improvement_roadmap"`
	Trends             []MaturityTrend           `json:"trends"`
	GeneratedAt        time.Time                 `json:"generated_at"`
	NextAssessment     time.Time                 `json:"next_assessment"`
	ConfidenceScore    float64                   `json:"confidence_score"`
}

// DimensionScore represents a score for a specific maturity dimension
type DimensionScore struct {
	Dimension       string   `json:"dimension"`
	Score           float64  `json:"score"`
	Weight          float64  `json:"weight"`
	Level           string   `json:"level"`
	Description     string   `json:"description"`
	Strengths       []string `json:"strengths"`
	Weaknesses      []string `json:"weaknesses"`
	Recommendations []string `json:"recommendations"`
}

// IndustryBenchmark represents industry comparison data
type IndustryBenchmark struct {
	Industry           string  `json:"industry"`
	IndustryAverage    float64 `json:"industry_average"`
	IndustryPercentile float64 `json:"industry_percentile"`
	TopPerformers      float64 `json:"top_performers"`
	MarketLeaders      float64 `json:"market_leaders"`
	CompetitiveGap     float64 `json:"competitive_gap"`
}

// PeerComparison represents peer organization comparison
type PeerComparison struct {
	PeerCount           int      `json:"peer_count"`
	PeerAverage         float64  `json:"peer_average"`
	PeerPercentile      float64  `json:"peer_percentile"`
	SimilarOrgs         []string `json:"similar_orgs"`
	CompetitivePosition string   `json:"competitive_position"`
}

// ImprovementItem represents an improvement recommendation
type ImprovementItem struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Priority      string   `json:"priority"`
	Impact        float64  `json:"impact"`
	Effort        string   `json:"effort"`
	Timeline      string   `json:"timeline"`
	Cost          string   `json:"cost"`
	ROI           float64  `json:"roi"`
	Prerequisites []string `json:"prerequisites"`
}

// MaturityTrend represents a trend in maturity scores
type MaturityTrend struct {
	Dimension   string  `json:"dimension"`
	Direction   string  `json:"direction"` // improving, declining, stable
	Magnitude   float64 `json:"magnitude"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// CalculateMaturityScore calculates comprehensive security maturity score
func (s *MaturityService) CalculateMaturityScore(organizationID uuid.UUID) (*MaturityScore, error) {
	// Get organization profile
	orgProfile, err := s.getOrganizationProfile(organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization profile: %w", err)
	}

	// Get vulnerability data
	vulnerabilities, err := s.getVulnerabilitiesForOrganization(organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities: %w", err)
	}

	// Get scan history
	scanHistory, err := s.getScanHistory(organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan history: %w", err)
	}

	// Calculate dimension scores
	dimensionScores := s.calculateDimensionScores(vulnerabilities, scanHistory, orgProfile)

	// Calculate overall score
	overallScore := s.calculateOverallScore(dimensionScores)

	// Determine maturity level
	maturityLevel := s.determineMaturityLevel(overallScore)

	// Get industry benchmark
	industryBenchmark := s.getIndustryBenchmark(orgProfile.Industry, overallScore)

	// Get peer comparison
	peerComparison := s.getPeerComparison(organizationID, overallScore)

	// Generate improvement roadmap
	improvementRoadmap := s.generateImprovementRoadmap(dimensionScores, orgProfile)

	// Analyze trends
	trends := s.analyzeMaturityTrends(organizationID)

	// Calculate confidence score
	confidenceScore := s.calculateConfidenceScore(vulnerabilities, scanHistory)

	// Create maturity score
	score := &MaturityScore{
		OrganizationID:     organizationID,
		ScoreID:            fmt.Sprintf("maturity_%s_%d", organizationID.String(), time.Now().Unix()),
		OverallScore:       overallScore,
		MaturityLevel:      maturityLevel,
		DimensionScores:    dimensionScores,
		IndustryBenchmark:  industryBenchmark,
		PeerComparison:     peerComparison,
		ImprovementRoadmap: improvementRoadmap,
		Trends:             trends,
		GeneratedAt:        time.Now(),
		NextAssessment:     time.Now().Add(30 * 24 * time.Hour), // 30 days
		ConfidenceScore:    confidenceScore,
	}

	return score, nil
}

// calculateDimensionScores calculates scores for each maturity dimension
func (s *MaturityService) calculateDimensionScores(vulnerabilities []models.Vulnerability, scanHistory []models.ScanResult, orgProfile *models.OrganizationProfile) map[string]DimensionScore {
	dimensions := make(map[string]DimensionScore)

	// Vulnerability Management Maturity
	dimensions["vulnerability_management"] = s.calculateVulnerabilityManagementScore(vulnerabilities, scanHistory)

	// Patch Management Maturity
	dimensions["patch_management"] = s.calculatePatchManagementScore(vulnerabilities, scanHistory)

	// Risk Management Maturity
	dimensions["risk_management"] = s.calculateRiskManagementScore(vulnerabilities, orgProfile)

	// Compliance Maturity
	dimensions["compliance"] = s.calculateComplianceScore(orgProfile, vulnerabilities)

	// Security Awareness Maturity
	dimensions["security_awareness"] = s.calculateSecurityAwarenessScore(orgProfile, vulnerabilities)

	// Technology Security Maturity
	dimensions["technology_security"] = s.calculateTechnologySecurityScore(orgProfile, vulnerabilities)

	// Incident Response Maturity
	dimensions["incident_response"] = s.calculateIncidentResponseScore(scanHistory, vulnerabilities)

	// Security Governance Maturity
	dimensions["security_governance"] = s.calculateSecurityGovernanceScore(orgProfile, vulnerabilities)

	return dimensions
}

// calculateVulnerabilityManagementScore calculates vulnerability management maturity
func (s *MaturityService) calculateVulnerabilityManagementScore(vulnerabilities []models.Vulnerability, scanHistory []models.ScanResult) DimensionScore {
	score := 0.0
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Calculate vulnerability discovery rate
	totalVulns := len(vulnerabilities)
	if totalVulns > 0 {
		criticalVulns := s.countVulnerabilitiesBySeverity(vulnerabilities, "CRITICAL")
		highVulns := s.countVulnerabilitiesBySeverity(vulnerabilities, "HIGH")

		// Score based on vulnerability distribution
		criticalRatio := float64(criticalVulns) / float64(totalVulns)
		highRatio := float64(highVulns) / float64(totalVulns)

		if criticalRatio < 0.1 {
			score += 0.3
			strengths = append(strengths, "Low critical vulnerability ratio")
		} else {
			weaknesses = append(weaknesses, "High critical vulnerability ratio")
			recommendations = append(recommendations, "Prioritize critical vulnerability remediation")
		}

		if highRatio < 0.3 {
			score += 0.2
			strengths = append(strengths, "Controlled high severity vulnerabilities")
		} else {
			weaknesses = append(weaknesses, "High severity vulnerability concentration")
			recommendations = append(recommendations, "Implement systematic high severity vulnerability management")
		}
	}

	// Calculate scan frequency and consistency
	scanFrequency := s.calculateScanFrequency(scanHistory)
	if scanFrequency > 0.8 {
		score += 0.2
		strengths = append(strengths, "Regular vulnerability scanning")
	} else if scanFrequency < 0.4 {
		weaknesses = append(weaknesses, "Irregular vulnerability scanning")
		recommendations = append(recommendations, "Implement automated vulnerability scanning")
	}

	// Calculate vulnerability age
	avgAge := s.calculateAverageVulnerabilityAge(vulnerabilities)
	if avgAge < 30 {
		score += 0.3
		strengths = append(strengths, "Quick vulnerability identification")
	} else if avgAge > 90 {
		weaknesses = append(weaknesses, "Old vulnerabilities present")
		recommendations = append(recommendations, "Implement vulnerability lifecycle management")
	}

	// Determine level
	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("vulnerability_management", score, level)

	return DimensionScore{
		Dimension:       "vulnerability_management",
		Score:           score,
		Weight:          0.15,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// calculatePatchManagementScore calculates patch management maturity
func (s *MaturityService) calculatePatchManagementScore(vulnerabilities []models.Vulnerability, scanHistory []models.ScanResult) DimensionScore {
	score := 0.0
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Calculate patch velocity
	patchVelocity := s.calculatePatchVelocity(vulnerabilities, scanHistory)
	if patchVelocity > 0.8 {
		score += 0.4
		strengths = append(strengths, "Fast patch deployment")
	} else if patchVelocity < 0.4 {
		weaknesses = append(weaknesses, "Slow patch deployment")
		recommendations = append(recommendations, "Implement automated patch management")
	}

	// Calculate patch coverage
	patchCoverage := s.calculatePatchCoverage(vulnerabilities)
	if patchCoverage > 0.9 {
		score += 0.3
		strengths = append(strengths, "High patch coverage")
	} else if patchCoverage < 0.6 {
		weaknesses = append(weaknesses, "Low patch coverage")
		recommendations = append(recommendations, "Improve patch coverage and tracking")
	}

	// Calculate patch testing
	patchTesting := s.calculatePatchTestingMaturity(scanHistory)
	if patchTesting > 0.7 {
		score += 0.3
		strengths = append(strengths, "Comprehensive patch testing")
	} else {
		weaknesses = append(weaknesses, "Insufficient patch testing")
		recommendations = append(recommendations, "Implement patch testing procedures")
	}

	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("patch_management", score, level)

	return DimensionScore{
		Dimension:       "patch_management",
		Score:           score,
		Weight:          0.15,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// calculateRiskManagementScore calculates risk management maturity
func (s *MaturityService) calculateRiskManagementScore(vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile) DimensionScore {
	score := 0.0
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Calculate risk assessment maturity
	riskAssessment := s.calculateRiskAssessmentMaturity(vulnerabilities, orgProfile)
	score += riskAssessment * 0.4

	if riskAssessment > 0.7 {
		strengths = append(strengths, "Comprehensive risk assessment")
	} else {
		weaknesses = append(weaknesses, "Limited risk assessment")
		recommendations = append(recommendations, "Implement formal risk assessment processes")
	}

	// Calculate risk tolerance alignment
	riskTolerance := s.calculateRiskToleranceAlignment(orgProfile, vulnerabilities)
	score += riskTolerance * 0.3

	if riskTolerance > 0.7 {
		strengths = append(strengths, "Aligned risk tolerance")
	} else {
		weaknesses = append(weaknesses, "Misaligned risk tolerance")
		recommendations = append(recommendations, "Review and align risk tolerance")
	}

	// Calculate risk monitoring
	riskMonitoring := s.calculateRiskMonitoringMaturity(vulnerabilities)
	score += riskMonitoring * 0.3

	if riskMonitoring > 0.7 {
		strengths = append(strengths, "Effective risk monitoring")
	} else {
		weaknesses = append(weaknesses, "Insufficient risk monitoring")
		recommendations = append(recommendations, "Implement continuous risk monitoring")
	}

	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("risk_management", score, level)

	return DimensionScore{
		Dimension:       "risk_management",
		Score:           score,
		Weight:          0.12,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// calculateComplianceScore calculates compliance maturity
func (s *MaturityService) calculateComplianceScore(orgProfile *models.OrganizationProfile, vulnerabilities []models.Vulnerability) DimensionScore {
	score := 0.0
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Calculate compliance framework coverage
	complianceCoverage := s.calculateComplianceCoverage(orgProfile)
	score += complianceCoverage * 0.4

	if complianceCoverage > 0.7 {
		strengths = append(strengths, "Comprehensive compliance coverage")
	} else {
		weaknesses = append(weaknesses, "Limited compliance coverage")
		recommendations = append(recommendations, "Expand compliance framework coverage")
	}

	// Calculate compliance monitoring
	complianceMonitoring := s.calculateComplianceMonitoringMaturity(vulnerabilities)
	score += complianceMonitoring * 0.3

	if complianceMonitoring > 0.7 {
		strengths = append(strengths, "Effective compliance monitoring")
	} else {
		weaknesses = append(weaknesses, "Insufficient compliance monitoring")
		recommendations = append(recommendations, "Implement automated compliance monitoring")
	}

	// Calculate compliance reporting
	complianceReporting := s.calculateComplianceReportingMaturity(orgProfile)
	score += complianceReporting * 0.3

	if complianceReporting > 0.7 {
		strengths = append(strengths, "Automated compliance reporting")
	} else {
		weaknesses = append(weaknesses, "Manual compliance reporting")
		recommendations = append(recommendations, "Implement automated compliance reporting")
	}

	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("compliance", score, level)

	return DimensionScore{
		Dimension:       "compliance",
		Score:           score,
		Weight:          0.13,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// calculateSecurityAwarenessScore calculates security awareness maturity
func (s *MaturityService) calculateSecurityAwarenessScore(orgProfile *models.OrganizationProfile, vulnerabilities []models.Vulnerability) DimensionScore {
	score := 0.5 // Default baseline
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Mock security awareness calculation
	// In a real implementation, this would analyze training records, phishing tests, etc.

	strengths = append(strengths, "Security awareness program in place")
	recommendations = append(recommendations, "Enhance security awareness training")

	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("security_awareness", score, level)

	return DimensionScore{
		Dimension:       "security_awareness",
		Score:           score,
		Weight:          0.10,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// calculateTechnologySecurityScore calculates technology security maturity
func (s *MaturityService) calculateTechnologySecurityScore(orgProfile *models.OrganizationProfile, vulnerabilities []models.Vulnerability) DimensionScore {
	score := 0.0
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Calculate technology stack security
	techStackSecurity := s.calculateTechStackSecurityMaturity(orgProfile, vulnerabilities)
	score += techStackSecurity * 0.5

	if techStackSecurity > 0.7 {
		strengths = append(strengths, "Secure technology stack")
	} else {
		weaknesses = append(weaknesses, "Technology stack security gaps")
		recommendations = append(recommendations, "Implement secure development practices")
	}

	// Calculate security tooling
	securityTooling := s.calculateSecurityToolingMaturity(orgProfile)
	score += securityTooling * 0.5

	if securityTooling > 0.7 {
		strengths = append(strengths, "Comprehensive security tooling")
	} else {
		weaknesses = append(weaknesses, "Limited security tooling")
		recommendations = append(recommendations, "Expand security tooling coverage")
	}

	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("technology_security", score, level)

	return DimensionScore{
		Dimension:       "technology_security",
		Score:           score,
		Weight:          0.15,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// calculateIncidentResponseScore calculates incident response maturity
func (s *MaturityService) calculateIncidentResponseScore(scanHistory []models.ScanResult, vulnerabilities []models.Vulnerability) DimensionScore {
	score := 0.5 // Default baseline
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Mock incident response calculation
	// In a real implementation, this would analyze incident response procedures, response times, etc.

	strengths = append(strengths, "Incident response procedures in place")
	recommendations = append(recommendations, "Enhance incident response capabilities")

	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("incident_response", score, level)

	return DimensionScore{
		Dimension:       "incident_response",
		Score:           score,
		Weight:          0.10,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// calculateSecurityGovernanceScore calculates security governance maturity
func (s *MaturityService) calculateSecurityGovernanceScore(orgProfile *models.OrganizationProfile, vulnerabilities []models.Vulnerability) DimensionScore {
	score := 0.0
	strengths := []string{}
	weaknesses := []string{}
	recommendations := []string{}

	// Calculate governance structure
	governanceStructure := s.calculateGovernanceStructureMaturity(orgProfile)
	score += governanceStructure * 0.4

	if governanceStructure > 0.7 {
		strengths = append(strengths, "Strong governance structure")
	} else {
		weaknesses = append(weaknesses, "Weak governance structure")
		recommendations = append(recommendations, "Strengthen security governance")
	}

	// Calculate policy maturity
	policyMaturity := s.calculatePolicyMaturity(orgProfile)
	score += policyMaturity * 0.3

	if policyMaturity > 0.7 {
		strengths = append(strengths, "Comprehensive security policies")
	} else {
		weaknesses = append(weaknesses, "Limited security policies")
		recommendations = append(recommendations, "Develop comprehensive security policies")
	}

	// Calculate oversight maturity
	oversightMaturity := s.calculateOversightMaturity(orgProfile)
	score += oversightMaturity * 0.3

	if oversightMaturity > 0.7 {
		strengths = append(strengths, "Effective security oversight")
	} else {
		weaknesses = append(weaknesses, "Insufficient security oversight")
		recommendations = append(recommendations, "Implement security oversight mechanisms")
	}

	level := s.determineDimensionLevel(score)
	description := s.generateDimensionDescription("security_governance", score, level)

	return DimensionScore{
		Dimension:       "security_governance",
		Score:           score,
		Weight:          0.10,
		Level:           level,
		Description:     description,
		Strengths:       strengths,
		Weaknesses:      weaknesses,
		Recommendations: recommendations,
	}
}

// Helper methods for maturity calculations

func (s *MaturityService) getOrganizationProfile(organizationID uuid.UUID) (*models.OrganizationProfile, error) {
	var profile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *MaturityService) getVulnerabilitiesForOrganization(organizationID uuid.UUID) ([]models.Vulnerability, error) {
	// Mock implementation - in real scenario, this would query vulnerabilities
	score1 := 8.5
	score2 := 6.2
	return []models.Vulnerability{
		{ID: uuid.New().String(), Title: "SQL Injection", Severity: models.SeverityHigh, CVSSScore: &score1},
		{ID: uuid.New().String(), Title: "XSS Vulnerability", Severity: models.SeverityMedium, CVSSScore: &score2},
	}, nil
}

func (s *MaturityService) getScanHistory(organizationID uuid.UUID) ([]models.ScanResult, error) {
	// Mock implementation - in real scenario, this would query scan history
	return []models.ScanResult{
		{ID: uuid.New(), Status: "completed", CreatedAt: time.Now().Add(-24 * time.Hour)},
		{ID: uuid.New(), Status: "completed", CreatedAt: time.Now().Add(-48 * time.Hour)},
	}, nil
}

func (s *MaturityService) countVulnerabilitiesBySeverity(vulnerabilities []models.Vulnerability, severity string) int {
	count := 0
	for _, vuln := range vulnerabilities {
		if string(vuln.Severity) == severity {
			count++
		}
	}
	return count
}

func (s *MaturityService) calculateScanFrequency(scanHistory []models.ScanResult) float64 {
	// Mock implementation
	return 0.8
}

func (s *MaturityService) calculateAverageVulnerabilityAge(vulnerabilities []models.Vulnerability) float64 {
	// Mock implementation
	return 45.0 // days
}

func (s *MaturityService) calculatePatchVelocity(vulnerabilities []models.Vulnerability, scanHistory []models.ScanResult) float64 {
	// Mock implementation
	return 0.7
}

func (s *MaturityService) calculatePatchCoverage(vulnerabilities []models.Vulnerability) float64 {
	// Mock implementation
	return 0.85
}

func (s *MaturityService) calculatePatchTestingMaturity(scanHistory []models.ScanResult) float64 {
	// Mock implementation
	return 0.6
}

func (s *MaturityService) calculateRiskAssessmentMaturity(vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile) float64 {
	// Mock implementation
	return 0.7
}

func (s *MaturityService) calculateRiskToleranceAlignment(orgProfile *models.OrganizationProfile, vulnerabilities []models.Vulnerability) float64 {
	// Mock implementation
	return 0.8
}

func (s *MaturityService) calculateRiskMonitoringMaturity(vulnerabilities []models.Vulnerability) float64 {
	// Mock implementation
	return 0.6
}

func (s *MaturityService) calculateComplianceCoverage(orgProfile *models.OrganizationProfile) float64 {
	if orgProfile == nil {
		return 0.3
	}
	return math.Min(float64(len(orgProfile.ComplianceFrameworks))/5.0, 1.0)
}

func (s *MaturityService) calculateComplianceMonitoringMaturity(vulnerabilities []models.Vulnerability) float64 {
	// Mock implementation
	return 0.7
}

func (s *MaturityService) calculateComplianceReportingMaturity(orgProfile *models.OrganizationProfile) float64 {
	// Mock implementation
	return 0.6
}

func (s *MaturityService) calculateTechStackSecurityMaturity(orgProfile *models.OrganizationProfile, vulnerabilities []models.Vulnerability) float64 {
	// Mock implementation
	return 0.8
}

func (s *MaturityService) calculateSecurityToolingMaturity(orgProfile *models.OrganizationProfile) float64 {
	// Mock implementation
	return 0.7
}

func (s *MaturityService) calculateGovernanceStructureMaturity(orgProfile *models.OrganizationProfile) float64 {
	// Mock implementation
	return 0.6
}

func (s *MaturityService) calculatePolicyMaturity(orgProfile *models.OrganizationProfile) float64 {
	// Mock implementation
	return 0.7
}

func (s *MaturityService) calculateOversightMaturity(orgProfile *models.OrganizationProfile) float64 {
	// Mock implementation
	return 0.5
}

func (s *MaturityService) calculateOverallScore(dimensionScores map[string]DimensionScore) float64 {
	totalScore := 0.0
	totalWeight := 0.0

	for _, dimension := range dimensionScores {
		totalScore += dimension.Score * dimension.Weight
		totalWeight += dimension.Weight
	}

	if totalWeight > 0 {
		return totalScore / totalWeight
	}
	return 0.0
}

func (s *MaturityService) determineMaturityLevel(score float64) string {
	if score >= 0.8 {
		return "Advanced"
	} else if score >= 0.6 {
		return "Intermediate"
	} else if score >= 0.4 {
		return "Basic"
	} else {
		return "Initial"
	}
}

func (s *MaturityService) determineDimensionLevel(score float64) string {
	if score >= 0.8 {
		return "Advanced"
	} else if score >= 0.6 {
		return "Intermediate"
	} else if score >= 0.4 {
		return "Basic"
	} else {
		return "Initial"
	}
}

func (s *MaturityService) generateDimensionDescription(dimension string, score float64, level string) string {
	return fmt.Sprintf("%s maturity is at %s level with a score of %.2f",
		dimension, level, score)
}

func (s *MaturityService) getIndustryBenchmark(industry string, score float64) IndustryBenchmark {
	// Mock industry benchmark data
	industryAverages := map[string]float64{
		"technology":    0.65,
		"healthcare":    0.70,
		"finance":       0.75,
		"government":    0.60,
		"manufacturing": 0.55,
	}

	industryAverage := industryAverages[strings.ToLower(industry)]
	if industryAverage == 0 {
		industryAverage = 0.60 // Default
	}

	percentile := (score - industryAverage) / industryAverage * 100
	if percentile < 0 {
		percentile = 0
	} else if percentile > 100 {
		percentile = 100
	}

	return IndustryBenchmark{
		Industry:           industry,
		IndustryAverage:    industryAverage,
		IndustryPercentile: percentile,
		TopPerformers:      industryAverage + 0.2,
		MarketLeaders:      industryAverage + 0.3,
		CompetitiveGap:     math.Max(0, industryAverage+0.2-score),
	}
}

func (s *MaturityService) getPeerComparison(organizationID uuid.UUID, score float64) PeerComparison {
	// Mock peer comparison data
	return PeerComparison{
		PeerCount:           25,
		PeerAverage:         0.65,
		PeerPercentile:      75.0,
		SimilarOrgs:         []string{"TechCorp", "SecureInc", "DataSafe"},
		CompetitivePosition: "Above Average",
	}
}

func (s *MaturityService) generateImprovementRoadmap(dimensionScores map[string]DimensionScore, orgProfile *models.OrganizationProfile) []ImprovementItem {
	var roadmap []ImprovementItem

	// Generate improvement items based on lowest scoring dimensions
	for dimension, score := range dimensionScores {
		if score.Score < 0.6 {
			roadmap = append(roadmap, ImprovementItem{
				ID:            fmt.Sprintf("improvement_%s", dimension),
				Title:         fmt.Sprintf("Improve %s", dimension),
				Description:   fmt.Sprintf("Enhance %s capabilities to reach intermediate level", dimension),
				Priority:      "High",
				Impact:        score.Score * 0.3,
				Effort:        "Medium",
				Timeline:      "3-6 months",
				Cost:          "$10K - $50K",
				ROI:           2.5,
				Prerequisites: []string{"Management approval", "Resource allocation"},
			})
		}
	}

	return roadmap
}

func (s *MaturityService) analyzeMaturityTrends(organizationID uuid.UUID) []MaturityTrend {
	// Mock trend analysis
	return []MaturityTrend{
		{
			Dimension:   "Overall",
			Direction:   "improving",
			Magnitude:   0.1,
			Confidence:  0.8,
			Description: "Overall maturity is improving gradually",
		},
	}
}

func (s *MaturityService) calculateConfidenceScore(vulnerabilities []models.Vulnerability, scanHistory []models.ScanResult) float64 {
	baseConfidence := 0.7

	// Adjust based on data availability
	if len(vulnerabilities) > 50 {
		baseConfidence += 0.1
	}
	if len(scanHistory) > 10 {
		baseConfidence += 0.1
	}

	return math.Min(baseConfidence, 1.0)
}
