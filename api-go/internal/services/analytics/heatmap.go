package analytics

import (
	"fmt"
	"math"
	"sort"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
)

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
func (s *AnalyticsService) GenerateRiskHeatmap(organizationID uuid.UUID, heatmapType string, timeRange string) (*HeatmapData, error) {
	// Get vulnerability data
	vulnerabilities, err := s.GetVulnerabilitiesForOrganization(organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities: %w", err)
	}

	// Calculate risk distribution
	riskDist := s.calculateRiskDistribution(vulnerabilities)

	// Generate data points based on heatmap type
	var dataPoints []HeatmapDataPoint
	var dimensions []HeatmapDimension

	switch heatmapType {
	case "severity_trend":
		dataPoints, dimensions = s.generateSeverityTrendHeatmap(vulnerabilities, timeRange)
	case "compliance_risk":
		dataPoints, dimensions = s.generateComplianceRiskHeatmap(vulnerabilities)
	case "technology":
		dataPoints, dimensions = s.generateTechnologyHeatmap(vulnerabilities)
	default:
		dataPoints, dimensions = s.generateDefaultHeatmap(vulnerabilities)
	}

	// Identify hotspots
	hotspots := s.identifyHotspots(vulnerabilities, dataPoints)

	// Calculate trends
	trends := s.calculateTrends(vulnerabilities, timeRange)

	// Generate recommendations
	recommendations := s.generateHeatmapRecommendations(hotspots, riskDist)

	return &HeatmapData{
		OrganizationID:   organizationID,
		HeatmapType:      heatmapType,
		Dimensions:       dimensions,
		DataPoints:       dataPoints,
		Hotspots:         hotspots,
		RiskDistribution: riskDist,
		Trends:           trends,
		Recommendations:  recommendations,
		GeneratedAt:      time.Now(),
		ConfidenceScore:  0.85,
	}, nil
}

// Helper methods for heatmap generation
func (s *AnalyticsService) calculateRiskDistribution(vulnerabilities []models.Vulnerability) RiskDistribution {
	dist := RiskDistribution{}
	
	for _, vuln := range vulnerabilities {
		dist.Total++
		switch string(vuln.Severity) {
		case "critical":
			dist.Critical++
		case "high":
			dist.High++
		case "medium":
			dist.Medium++
		case "low":
			dist.Low++
		}
	}
	
	if dist.Total > 0 {
		dist.Average = float64(dist.Critical*10+dist.High*7+dist.Medium*4+dist.Low*1) / float64(dist.Total)
	}
	
	return dist
}

func (s *AnalyticsService) generateSeverityTrendHeatmap(vulnerabilities []models.Vulnerability, timeRange string) ([]HeatmapDataPoint, []HeatmapDimension) {
	// Implementation for severity trend heatmap
	dataPoints := []HeatmapDataPoint{}
	dimensions := []HeatmapDimension{
		{Name: "Severity", Type: "severity", Categories: []string{"critical", "high", "medium", "low"}},
		{Name: "Time", Type: "trend", Categories: []string{"week1", "week2", "week3", "week4"}},
	}
	
	// Group by severity and time period
	severityCounts := make(map[string]map[string]int)
	for _, vuln := range vulnerabilities {
		severity := string(vuln.Severity)
		period := s.getTimePeriod(vuln.CreatedAt, timeRange)
		
		if severityCounts[severity] == nil {
			severityCounts[severity] = make(map[string]int)
		}
		severityCounts[severity][period]++
	}
	
	// Create data points
	for severity, periods := range severityCounts {
		for period, count := range periods {
			dataPoints = append(dataPoints, HeatmapDataPoint{
				X:         severity,
				Y:         period,
				Value:     float64(count),
				Count:     count,
				RiskLevel: severity,
				Trend:     "stable",
			})
		}
	}
	
	return dataPoints, dimensions
}

func (s *AnalyticsService) generateComplianceRiskHeatmap(vulnerabilities []models.Vulnerability) ([]HeatmapDataPoint, []HeatmapDimension) {
	// Implementation for compliance risk heatmap
	return []HeatmapDataPoint{}, []HeatmapDimension{}
}

func (s *AnalyticsService) generateTechnologyHeatmap(vulnerabilities []models.Vulnerability) ([]HeatmapDataPoint, []HeatmapDimension) {
	// Implementation for technology heatmap
	return []HeatmapDataPoint{}, []HeatmapDimension{}
}

func (s *AnalyticsService) generateDefaultHeatmap(vulnerabilities []models.Vulnerability) ([]HeatmapDataPoint, []HeatmapDimension) {
	// Default heatmap implementation
	return []HeatmapDataPoint{}, []HeatmapDimension{}
}

func (s *AnalyticsService) identifyHotspots(vulnerabilities []models.Vulnerability, dataPoints []HeatmapDataPoint) []Hotspot {
	hotspots := []Hotspot{}
	
	// Group by severity and count
	severityCounts := make(map[string]int)
	for _, vuln := range vulnerabilities {
		severityCounts[string(vuln.Severity)]++
	}
	
	// Create hotspots for high-risk areas
	for severity, count := range severityCounts {
		if count > 10 || severity == "critical" || severity == "high" {
			hotspots = append(hotspots, Hotspot{
				ID:          fmt.Sprintf("hotspot-%s", severity),
				Name:        fmt.Sprintf("%s Vulnerabilities", severity),
				RiskScore:   s.calculateRiskScore(severity, count),
				Severity:    severity,
				Count:       count,
				Trend:       "stable",
				Description: fmt.Sprintf("High concentration of %s severity vulnerabilities", severity),
				Actions:     []string{"Review and prioritize remediation", "Implement additional monitoring"},
			})
		}
	}
	
	// Sort by risk score
	sort.Slice(hotspots, func(i, j int) bool {
		return hotspots[i].RiskScore > hotspots[j].RiskScore
	})
	
	return hotspots
}

func (s *AnalyticsService) calculateTrends(vulnerabilities []models.Vulnerability, timeRange string) []Trend {
	// Calculate trends over time
	return []Trend{}
}

func (s *AnalyticsService) generateHeatmapRecommendations(hotspots []Hotspot, riskDist RiskDistribution) []string {
	recommendations := []string{}
	
	if riskDist.Critical > 0 {
		recommendations = append(recommendations, "Immediate action required: Address critical vulnerabilities")
	}
	if riskDist.High > 10 {
		recommendations = append(recommendations, "High number of high-severity vulnerabilities detected")
	}
	if len(hotspots) > 5 {
		recommendations = append(recommendations, "Multiple risk hotspots identified - consider comprehensive security review")
	}
	
	return recommendations
}

func (s *AnalyticsService) calculateRiskScore(severity string, count int) float64 {
	severityWeights := map[string]float64{
		"critical": 10.0,
		"high":     7.0,
		"medium":   4.0,
		"low":      1.0,
	}
	
	weight := severityWeights[severity]
	return weight * math.Log(float64(count+1))
}

func (s *AnalyticsService) getTimePeriod(t time.Time, timeRange string) string {
	// Simple implementation - returns week number
	_, week := t.ISOWeek()
	return fmt.Sprintf("week%d", week%4+1)
}

