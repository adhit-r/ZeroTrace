package analytics

import (
	"fmt"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
)

// MaturityScore represents a comprehensive security maturity score
type MaturityScore struct {
	OrganizationID     uuid.UUID                 `json:"organization_id"`
	ScoreID            string                    `json:"score_id"`
	OverallScore       float64                   `json:"overall_score"`
	MaturityLevel      string                    `json:"maturity_level"`
	DimensionScores    map[string]DimensionScore `json:"dimension_scores"`
	IndustryBenchmark  IndustryBenchmark          `json:"industry_benchmark"`
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
func (s *AnalyticsService) CalculateMaturityScore(organizationID uuid.UUID) (*MaturityScore, error) {
	// Get vulnerability data
	vulnerabilities, err := s.GetVulnerabilitiesForOrganization(organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get vulnerabilities: %w", err)
	}

	// Get scan history
	scanHistory, err := s.GetScanHistory(organizationID, 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan history: %w", err)
	}

	// Calculate dimension scores
	dimensionScores := s.calculateDimensionScores(vulnerabilities, scanHistory)

	// Calculate overall score
	overallScore := s.calculateOverallMaturityScore(dimensionScores)

	// Determine maturity level
	maturityLevel := s.determineMaturityLevel(overallScore)

	// Get industry benchmark (simplified)
	industryBenchmark := IndustryBenchmark{
		Industry:           "Technology",
		IndustryAverage:    65.0,
		IndustryPercentile: 50.0,
		TopPerformers:      85.0,
		MarketLeaders:      90.0,
		CompetitiveGap:     overallScore - 65.0,
	}

	// Get peer comparison (simplified)
	peerComparison := PeerComparison{
		PeerCount:           100,
		PeerAverage:        65.0,
		PeerPercentile:     50.0,
		SimilarOrgs:         []string{},
		CompetitivePosition: "average",
	}

	// Generate improvement roadmap
	improvementRoadmap := s.generateImprovementRoadmap(dimensionScores)

	// Analyze trends
	trends := s.analyzeMaturityTrends(organizationID)

	return &MaturityScore{
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
		NextAssessment:     time.Now().Add(30 * 24 * time.Hour),
		ConfidenceScore:    0.85,
	}, nil
}

// Helper methods
func (s *AnalyticsService) calculateDimensionScores(vulnerabilities []models.Vulnerability, scanHistory []models.Scan) map[string]DimensionScore {
	dimensions := make(map[string]DimensionScore)
	
	// Simplified dimension calculations
	dimensions["vulnerability_management"] = DimensionScore{
		Dimension:       "vulnerability_management",
		Score:           70.0,
		Weight:          0.15,
		Level:           "intermediate",
		Description:     "Vulnerability management maturity",
		Strengths:       []string{"Regular scanning"},
		Weaknesses:      []string{"Slow remediation"},
		Recommendations: []string{"Implement automated remediation"},
	}
	
	return dimensions
}

func (s *AnalyticsService) calculateOverallMaturityScore(dimensions map[string]DimensionScore) float64 {
	total := 0.0
	totalWeight := 0.0
	
	for _, dim := range dimensions {
		total += dim.Score * dim.Weight
		totalWeight += dim.Weight
	}
	
	if totalWeight > 0 {
		return total / totalWeight
	}
	return 0.0
}

func (s *AnalyticsService) determineMaturityLevel(score float64) string {
	if score >= 80 {
		return "advanced"
	} else if score >= 60 {
		return "intermediate"
	} else if score >= 40 {
		return "basic"
	}
	return "initial"
}

func (s *AnalyticsService) generateImprovementRoadmap(dimensions map[string]DimensionScore) []ImprovementItem {
	roadmap := []ImprovementItem{}
	
	for _, dim := range dimensions {
		if dim.Score < 70 {
			roadmap = append(roadmap, ImprovementItem{
				ID:            fmt.Sprintf("improvement_%s", dim.Dimension),
				Title:         fmt.Sprintf("Improve %s", dim.Dimension),
				Description:   dim.Description,
				Priority:      "high",
				Impact:        0.8,
				Effort:        "medium",
				Timeline:      "3 months",
				Cost:          "medium",
				ROI:           2.5,
				Prerequisites: []string{},
			})
		}
	}
	
	return roadmap
}

func (s *AnalyticsService) analyzeMaturityTrends(organizationID uuid.UUID) []MaturityTrend {
	return []MaturityTrend{
		{
			Dimension:   "overall",
			Direction:   "improving",
			Magnitude:   5.0,
			Confidence:  0.8,
			Description: "Overall maturity is improving",
		},
	}
}

