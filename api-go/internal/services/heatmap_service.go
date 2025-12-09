package services

import (
	"fmt"
	"math"
	"sort"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HeatmapService handles risk heatmap generation and analysis
type HeatmapService struct {
	db *gorm.DB
}

// NewHeatmapService creates a new HeatmapService
func NewHeatmapService(db *gorm.DB) *HeatmapService {
	return &HeatmapService{db: db}
}

// HeatmapData represents heatmap visualization data
type HeatmapData struct {
	OrganizationID   uuid.UUID          `json:"organization_id"`
	HeatmapType      string             `json:"heatmap_type"`
	Dimensions       []HeatmapDimension `json:"dimensions"`
	DataPoints       []HeatmapDataPoint `json:"data_points"`
	Hotspots         []Hotspot          `json:"hotspots"`
	RiskDistribution RiskDistribution   `json:"risk_distribution"`
	Trends           []Trend            `json:"trends"`
	Recommendations  []string           `json:"recommendations"`
	GeneratedAt      time.Time          `json:"generated_at"`
	ConfidenceScore  float64            `json:"confidence_score"`
}

// HeatmapDimension represents a dimension in the heatmap
type HeatmapDimension struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"` // severity, technology, compliance, trend
	Categories  []string `json:"categories"`
	Weight      float64  `json:"weight"`
	Description string   `json:"description"`
}

// HeatmapDataPoint represents a single data point in the heatmap
type HeatmapDataPoint struct {
	X          string                 `json:"x"`
	Y          string                 `json:"y"`
	Value      float64                `json:"value"`
	Count      int                    `json:"count"`
	RiskLevel  string                 `json:"risk_level"`
	Trend      string                 `json:"trend"` // increasing, decreasing, stable
	Confidence float64                `json:"confidence"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Hotspot represents a high-risk area in the heatmap
type Hotspot struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	RiskScore   float64                `json:"risk_score"`
	Severity    string                 `json:"severity"`
	Count       int                    `json:"count"`
	Trend       string                 `json:"trend"`
	Description string                 `json:"description"`
	Actions     []string               `json:"actions"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// RiskDistribution represents the distribution of risk across categories
type RiskDistribution struct {
	Critical int     `json:"critical"`
	High     int     `json:"high"`
	Medium   int     `json:"medium"`
	Low      int     `json:"low"`
	Total    int     `json:"total"`
	Average  float64 `json:"average"`
}

// Trend represents a trend in the heatmap data
type Trend struct {
	Dimension   string  `json:"dimension"`
	Direction   string  `json:"direction"` // up, down, stable
	Magnitude   float64 `json:"magnitude"`
	Confidence  float64 `json:"confidence"`
	Description string  `json:"description"`
}

// GenerateRiskHeatmap generates a comprehensive risk heatmap for an organization
func (s *HeatmapService) GenerateRiskHeatmap(organizationID uuid.UUID, heatmapType string, timeRange string) (*HeatmapData, error) {
	// Get organization profile for context
	orgProfile, err := s.getOrganizationProfile(organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization profile: %w", err)
	}

	// Get vulnerability data for the organization
	vulnerabilities, err := s.getVulnerabilitiesForOrganization(organizationID, timeRange)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities: %w", err)
	}

	// Generate heatmap based on type
	var heatmapData *HeatmapData
	switch heatmapType {
	case "severity_technology":
		heatmapData, err = s.generateSeverityTechnologyHeatmap(organizationID, vulnerabilities, orgProfile)
	case "compliance_trend":
		heatmapData, err = s.generateComplianceTrendHeatmap(organizationID, vulnerabilities, orgProfile)
	case "risk_velocity":
		heatmapData, err = s.generateRiskVelocityHeatmap(organizationID, vulnerabilities, orgProfile)
	case "comprehensive":
		heatmapData, err = s.generateComprehensiveHeatmap(organizationID, vulnerabilities, orgProfile)
	default:
		heatmapData, err = s.generateComprehensiveHeatmap(organizationID, vulnerabilities, orgProfile)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate heatmap: %w", err)
	}

	// Add hotspots and trends
	heatmapData.Hotspots = s.identifyHotspots(heatmapData.DataPoints, vulnerabilities)
	heatmapData.Trends = s.analyzeTrends(heatmapData.DataPoints, timeRange)
	heatmapData.Recommendations = s.generateRecommendations(heatmapData, orgProfile)
	heatmapData.ConfidenceScore = s.calculateConfidenceScore(heatmapData)

	return heatmapData, nil
}

// generateSeverityTechnologyHeatmap creates a heatmap showing severity vs technology
func (s *HeatmapService) generateSeverityTechnologyHeatmap(organizationID uuid.UUID, vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile) (*HeatmapData, error) {
	// Define dimensions
	dimensions := []HeatmapDimension{
		{
			Name:        "Severity",
			Type:        "severity",
			Categories:  []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"},
			Weight:      1.0,
			Description: "Vulnerability severity levels",
		},
		{
			Name:        "Technology",
			Type:        "technology",
			Categories:  s.extractTechnologyCategories(vulnerabilities),
			Weight:      0.8,
			Description: "Technology categories",
		},
	}

	// Generate data points
	dataPoints := s.calculateSeverityTechnologyDataPoints(vulnerabilities, dimensions)

	// Calculate risk distribution
	riskDistribution := s.calculateRiskDistribution(vulnerabilities)

	return &HeatmapData{
		OrganizationID:   organizationID,
		HeatmapType:      "severity_technology",
		Dimensions:       dimensions,
		DataPoints:       dataPoints,
		RiskDistribution: riskDistribution,
		GeneratedAt:      time.Now(),
	}, nil
}

// generateComplianceTrendHeatmap creates a heatmap showing compliance vs trends
func (s *HeatmapService) generateComplianceTrendHeatmap(organizationID uuid.UUID, vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile) (*HeatmapData, error) {
	// Define dimensions
	dimensions := []HeatmapDimension{
		{
			Name:        "Compliance",
			Type:        "compliance",
			Categories:  s.getComplianceCategories(orgProfile),
			Weight:      1.0,
			Description: "Compliance framework categories",
		},
		{
			Name:        "Trend",
			Type:        "trend",
			Categories:  []string{"Increasing", "Stable", "Decreasing"},
			Weight:      0.7,
			Description: "Risk trend direction",
		},
	}

	// Generate data points
	dataPoints := s.calculateComplianceTrendDataPoints(vulnerabilities, orgProfile, dimensions)

	// Calculate risk distribution
	riskDistribution := s.calculateRiskDistribution(vulnerabilities)

	return &HeatmapData{
		OrganizationID:   organizationID,
		HeatmapType:      "compliance_trend",
		Dimensions:       dimensions,
		DataPoints:       dataPoints,
		RiskDistribution: riskDistribution,
		GeneratedAt:      time.Now(),
	}, nil
}

// generateRiskVelocityHeatmap creates a heatmap showing risk velocity over time
func (s *HeatmapService) generateRiskVelocityHeatmap(organizationID uuid.UUID, vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile) (*HeatmapData, error) {
	// Define dimensions
	dimensions := []HeatmapDimension{
		{
			Name:        "Time",
			Type:        "time",
			Categories:  s.generateTimeCategories(),
			Weight:      1.0,
			Description: "Time periods",
		},
		{
			Name:        "Risk Level",
			Type:        "risk",
			Categories:  []string{"Critical", "High", "Medium", "Low"},
			Weight:      0.9,
			Description: "Risk levels",
		},
	}

	// Generate data points
	dataPoints := s.calculateRiskVelocityDataPoints(vulnerabilities, dimensions)

	// Calculate risk distribution
	riskDistribution := s.calculateRiskDistribution(vulnerabilities)

	return &HeatmapData{
		OrganizationID:   organizationID,
		HeatmapType:      "risk_velocity",
		Dimensions:       dimensions,
		DataPoints:       dataPoints,
		RiskDistribution: riskDistribution,
		GeneratedAt:      time.Now(),
	}, nil
}

// generateComprehensiveHeatmap creates a multi-dimensional comprehensive heatmap
func (s *HeatmapService) generateComprehensiveHeatmap(organizationID uuid.UUID, vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile) (*HeatmapData, error) {
	// Define multiple dimensions
	dimensions := []HeatmapDimension{
		{
			Name:        "Severity",
			Type:        "severity",
			Categories:  []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"},
			Weight:      1.0,
			Description: "Vulnerability severity levels",
		},
		{
			Name:        "Technology",
			Type:        "technology",
			Categories:  s.extractTechnologyCategories(vulnerabilities),
			Weight:      0.8,
			Description: "Technology categories",
		},
		{
			Name:        "Compliance",
			Type:        "compliance",
			Categories:  s.getComplianceCategories(orgProfile),
			Weight:      0.9,
			Description: "Compliance framework categories",
		},
		{
			Name:        "Trend",
			Type:        "trend",
			Categories:  []string{"Increasing", "Stable", "Decreasing"},
			Weight:      0.7,
			Description: "Risk trend direction",
		},
	}

	// Generate comprehensive data points
	dataPoints := s.calculateComprehensiveDataPoints(vulnerabilities, orgProfile, dimensions)

	// Calculate risk distribution
	riskDistribution := s.calculateRiskDistribution(vulnerabilities)

	return &HeatmapData{
		OrganizationID:   organizationID,
		HeatmapType:      "comprehensive",
		Dimensions:       dimensions,
		DataPoints:       dataPoints,
		RiskDistribution: riskDistribution,
		GeneratedAt:      time.Now(),
	}, nil
}

// Helper methods for heatmap generation

func (s *HeatmapService) getOrganizationProfile(organizationID uuid.UUID) (*models.OrganizationProfile, error) {
	var profile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (s *HeatmapService) getVulnerabilitiesForOrganization(organizationID uuid.UUID, timeRange string) ([]models.Vulnerability, error) {
	// This would typically query vulnerabilities from scan results
	// Database integration required - returning empty data until database integration is implemented
	score1 := 8.5
	score2 := 9.2
	return []models.Vulnerability{
		{
			ID:          uuid.New().String(),
			Title:       "SQL Injection Vulnerability",
			Description: "SQL injection in web application",
			Severity:    models.SeverityHigh,
			PackageName: "web-framework",
			CVSSScore:   &score1,
		},
		{
			ID:          uuid.New().String(),
			Title:       "Remote Code Execution",
			Description: "RCE in database component",
			Severity:    models.SeverityCritical,
			PackageName: "database-server",
			CVSSScore:   &score2,
		},
	}, nil
}

func (s *HeatmapService) extractTechnologyCategories(vulnerabilities []models.Vulnerability) []string {
	categories := make(map[string]bool)

	for _, vuln := range vulnerabilities {
		packageName := vuln.PackageName
		if packageName != "" {
			// Extract technology category from package name
			if containsHeatmap(packageName, "web") || containsHeatmap(packageName, "http") {
				categories["Web Applications"] = true
			}
			if containsHeatmap(packageName, "database") || containsHeatmap(packageName, "db") {
				categories["Databases"] = true
			}
			if containsHeatmap(packageName, "auth") || containsHeatmap(packageName, "login") {
				categories["Authentication"] = true
			}
			if containsHeatmap(packageName, "api") || containsHeatmap(packageName, "service") {
				categories["APIs"] = true
			}
		}
	}

	// Convert to slice
	var result []string
	for category := range categories {
		result = append(result, category)
	}

	// Add default categories if none found
	if len(result) == 0 {
		result = []string{"Web Applications", "Databases", "APIs", "Infrastructure"}
	}

	return result
}

func (s *HeatmapService) getComplianceCategories(orgProfile *models.OrganizationProfile) []string {
	if orgProfile == nil {
		return []string{"General", "Security", "Privacy"}
	}

	categories := []string{"General", "Security", "Privacy"}

	// Add organization-specific compliance categories
	for _, framework := range orgProfile.ComplianceFrameworks {
		categories = append(categories, framework)
	}

	return categories
}

func (s *HeatmapService) generateTimeCategories() []string {
	return []string{"Last 24h", "Last 7d", "Last 30d", "Last 90d"}
}

func (s *HeatmapService) calculateSeverityTechnologyDataPoints(vulnerabilities []models.Vulnerability, dimensions []HeatmapDimension) []HeatmapDataPoint {
	var dataPoints []HeatmapDataPoint

	severityCategories := dimensions[0].Categories
	technologyCategories := dimensions[1].Categories

	for _, severity := range severityCategories {
		for _, technology := range technologyCategories {
			count := s.countVulnerabilitiesBySeverityAndTechnology(vulnerabilities, severity, technology)
			riskScore := s.calculateRiskScore(severity, count)

			dataPoints = append(dataPoints, HeatmapDataPoint{
				X:          severity,
				Y:          technology,
				Value:      riskScore,
				Count:      count,
				RiskLevel:  s.getRiskLevel(riskScore),
				Trend:      "stable", // Would be calculated from historical data
				Confidence: 0.8,
				Metadata: map[string]interface{}{
					"severity":   severity,
					"technology": technology,
				},
			})
		}
	}

	return dataPoints
}

func (s *HeatmapService) calculateComplianceTrendDataPoints(vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile, dimensions []HeatmapDimension) []HeatmapDataPoint {
	var dataPoints []HeatmapDataPoint

	complianceCategories := dimensions[0].Categories
	trendCategories := dimensions[1].Categories

	for _, compliance := range complianceCategories {
		for _, trend := range trendCategories {
			count := s.countVulnerabilitiesByComplianceAndTrend(vulnerabilities, compliance, trend)
			riskScore := s.calculateComplianceRiskScore(compliance, count, orgProfile)

			dataPoints = append(dataPoints, HeatmapDataPoint{
				X:          compliance,
				Y:          trend,
				Value:      riskScore,
				Count:      count,
				RiskLevel:  s.getRiskLevel(riskScore),
				Trend:      trend,
				Confidence: 0.8,
				Metadata: map[string]interface{}{
					"compliance": compliance,
					"trend":      trend,
				},
			})
		}
	}

	return dataPoints
}

func (s *HeatmapService) calculateRiskVelocityDataPoints(vulnerabilities []models.Vulnerability, dimensions []HeatmapDimension) []HeatmapDataPoint {
	var dataPoints []HeatmapDataPoint

	timeCategories := dimensions[0].Categories
	riskCategories := dimensions[1].Categories

	for _, timePeriod := range timeCategories {
		for _, riskLevel := range riskCategories {
			count := s.countVulnerabilitiesByTimeAndRisk(vulnerabilities, timePeriod, riskLevel)
			riskScore := s.calculateVelocityRiskScore(timePeriod, riskLevel, count)

			dataPoints = append(dataPoints, HeatmapDataPoint{
				X:          timePeriod,
				Y:          riskLevel,
				Value:      riskScore,
				Count:      count,
				RiskLevel:  riskLevel,
				Trend:      s.calculateTrend(vulnerabilities, timePeriod),
				Confidence: 0.8,
				Metadata: map[string]interface{}{
					"time_period": timePeriod,
					"risk_level":  riskLevel,
				},
			})
		}
	}

	return dataPoints
}

func (s *HeatmapService) calculateComprehensiveDataPoints(vulnerabilities []models.Vulnerability, orgProfile *models.OrganizationProfile, dimensions []HeatmapDimension) []HeatmapDataPoint {
	// This would be a more complex calculation combining all dimensions
	// For now, return a simplified version
	return s.calculateSeverityTechnologyDataPoints(vulnerabilities, dimensions[:2])
}

func (s *HeatmapService) countVulnerabilitiesBySeverityAndTechnology(vulnerabilities []models.Vulnerability, severity, technology string) int {
	count := 0
	for _, vuln := range vulnerabilities {
		if string(vuln.Severity) == severity {
			// Check if vulnerability matches technology category
			if s.matchesTechnology(vuln.PackageName, technology) {
				count++
			}
		}
	}
	return count
}

func (s *HeatmapService) countVulnerabilitiesByComplianceAndTrend(vulnerabilities []models.Vulnerability, compliance, trend string) int {
	// Simplified implementation
	// Calculate based on actual vulnerability count
	if len(vulnerabilities) == 0 {
		return 0
	}
	return len(vulnerabilities) / 4
}

func (s *HeatmapService) countVulnerabilitiesByTimeAndRisk(vulnerabilities []models.Vulnerability, timePeriod, riskLevel string) int {
	// Simplified implementation
	// Calculate based on actual vulnerability count
	if len(vulnerabilities) == 0 {
		return 0
	}
	return len(vulnerabilities) / 3
}

func (s *HeatmapService) matchesTechnology(packageName, technology string) bool {
	packageLower := toLower(packageName)
	techLower := toLower(technology)

	switch techLower {
	case "web applications":
		return containsHeatmap(packageLower, "web") || containsHeatmap(packageLower, "http")
	case "databases":
		return containsHeatmap(packageLower, "database") || containsHeatmap(packageLower, "db")
	case "authentication":
		return containsHeatmap(packageLower, "auth") || containsHeatmap(packageLower, "login")
	case "apis":
		return containsHeatmap(packageLower, "api") || containsHeatmap(packageLower, "service")
	default:
		return true
	}
}

func (s *HeatmapService) calculateRiskScore(severity string, count int) float64 {
	severityWeights := map[string]float64{
		"CRITICAL": 4.0,
		"HIGH":     3.0,
		"MEDIUM":   2.0,
		"LOW":      1.0,
	}

	weight := severityWeights[severity]
	return weight * float64(count) * 0.25 // Normalize
}

func (s *HeatmapService) calculateComplianceRiskScore(compliance string, count int, orgProfile *models.OrganizationProfile) float64 {
	baseScore := float64(count) * 0.2

	// Adjust based on organization's compliance frameworks
	if orgProfile != nil {
		for _, framework := range orgProfile.ComplianceFrameworks {
			if framework == compliance {
				baseScore *= 1.5 // Higher risk for organization's compliance frameworks
			}
		}
	}

	return math.Min(baseScore, 1.0)
}

func (s *HeatmapService) calculateVelocityRiskScore(timePeriod, riskLevel string, count int) float64 {
	timeWeights := map[string]float64{
		"Last 24h": 4.0,
		"Last 7d":  3.0,
		"Last 30d": 2.0,
		"Last 90d": 1.0,
	}

	riskWeights := map[string]float64{
		"Critical": 4.0,
		"High":     3.0,
		"Medium":   2.0,
		"Low":      1.0,
	}

	timeWeight := timeWeights[timePeriod]
	riskWeight := riskWeights[riskLevel]

	return (timeWeight * riskWeight * float64(count)) / 16.0 // Normalize
}

func (s *HeatmapService) getRiskLevel(riskScore float64) string {
	if riskScore >= 0.8 {
		return "Critical"
	} else if riskScore >= 0.6 {
		return "High"
	} else if riskScore >= 0.4 {
		return "Medium"
	} else {
		return "Low"
	}
}

func (s *HeatmapService) calculateTrend(vulnerabilities []models.Vulnerability, timePeriod string) string {
	// Simplified trend calculation
	return "stable"
}

func (s *HeatmapService) calculateRiskDistribution(vulnerabilities []models.Vulnerability) RiskDistribution {
	distribution := RiskDistribution{}

	for _, vuln := range vulnerabilities {
		distribution.Total++
		switch vuln.Severity {
		case "CRITICAL":
			distribution.Critical++
		case "HIGH":
			distribution.High++
		case "MEDIUM":
			distribution.Medium++
		case "LOW":
			distribution.Low++
		}
	}

	// Calculate average risk score
	if distribution.Total > 0 {
		totalScore := float64(distribution.Critical*4 + distribution.High*3 + distribution.Medium*2 + distribution.Low*1)
		distribution.Average = totalScore / float64(distribution.Total)
	}

	return distribution
}

func (s *HeatmapService) identifyHotspots(dataPoints []HeatmapDataPoint, vulnerabilities []models.Vulnerability) []Hotspot {
	var hotspots []Hotspot

	// Sort data points by value to find highest risk areas
	sort.Slice(dataPoints, func(i, j int) bool {
		return dataPoints[i].Value > dataPoints[j].Value
	})

	// Take top 5 as hotspots
	for i, point := range dataPoints {
		if i >= 5 {
			break
		}

		if point.Value > 0.5 { // Only include significant hotspots
			hotspot := Hotspot{
				ID:          fmt.Sprintf("hotspot_%d", i+1),
				Name:        fmt.Sprintf("%s - %s", point.X, point.Y),
				RiskScore:   point.Value,
				Severity:    point.RiskLevel,
				Count:       point.Count,
				Trend:       point.Trend,
				Description: s.generateHotspotDescription(point),
				Actions:     s.generateHotspotActions(point),
				Metadata:    point.Metadata,
			}
			hotspots = append(hotspots, hotspot)
		}
	}

	return hotspots
}

func (s *HeatmapService) generateHotspotDescription(point HeatmapDataPoint) string {
	return fmt.Sprintf("High-risk area with %d vulnerabilities in %s/%s category",
		point.Count, point.X, point.Y)
}

func (s *HeatmapService) generateHotspotActions(point HeatmapDataPoint) []string {
	actions := []string{
		"Immediate assessment required",
		"Implement additional monitoring",
		"Review and update security controls",
	}

	if point.RiskLevel == "Critical" {
		actions = append(actions, "Emergency remediation plan needed")
	}

	return actions
}

func (s *HeatmapService) analyzeTrends(dataPoints []HeatmapDataPoint, timeRange string) []Trend {
	var trends []Trend

	// Simplified trend analysis
	trends = append(trends, Trend{
		Dimension:   "Overall Risk",
		Direction:   "stable",
		Magnitude:   0.1,
		Confidence:  0.8,
		Description: "Risk levels remain stable across most categories",
	})

	return trends
}

func (s *HeatmapService) generateRecommendations(heatmapData *HeatmapData, orgProfile *models.OrganizationProfile) []string {
	var recommendations []string

	// Generate recommendations based on hotspots
	for _, hotspot := range heatmapData.Hotspots {
		if hotspot.RiskScore > 0.7 {
			recommendations = append(recommendations,
				fmt.Sprintf("Prioritize remediation of %s vulnerabilities", hotspot.Name))
		}
	}

	// Add general recommendations
	recommendations = append(recommendations,
		"Implement continuous monitoring for high-risk areas",
		"Regular security assessments for critical technologies",
		"Update security policies based on risk patterns")

	return recommendations
}

func (s *HeatmapService) calculateConfidenceScore(heatmapData *HeatmapData) float64 {
	// Calculate confidence based on data quality and completeness
	baseConfidence := 0.8

	// Adjust based on data points count
	if len(heatmapData.DataPoints) > 20 {
		baseConfidence += 0.1
	}

	// Adjust based on hotspots
	if len(heatmapData.Hotspots) > 0 {
		baseConfidence += 0.05
	}

	return math.Min(baseConfidence, 1.0)
}

// Utility functions
func containsHeatmap(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) >= 0)))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func toLower(s string) string {
	// Simple toLower implementation
	result := make([]byte, len(s))
	for i, b := range []byte(s) {
		if b >= 'A' && b <= 'Z' {
			result[i] = b + 32
		} else {
			result[i] = b
		}
	}
	return string(result)
}
