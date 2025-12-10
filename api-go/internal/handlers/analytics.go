package handlers

import (
	"net/http"

	analytics "zerotrace/api/internal/services/analytics"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AnalyticsHandler handles unified analytics API endpoints (heatmap, maturity, compliance)
type AnalyticsHandler struct {
	analyticsService *analytics.AnalyticsService
}

// NewAnalyticsHandler creates a new AnalyticsHandler
func NewAnalyticsHandler(analyticsService *analytics.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetDashboardHistory returns dashboard history snapshots
func (h *AnalyticsHandler) GetDashboardHistory(c *gin.Context) {
	organizationIDStr := c.Query("organization_id")
	if organizationIDStr == "" {
		BadRequest(c, "MISSING_PARAM", "Organization ID is required", nil)
		return
	}

	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	days := 30 // Default to 30 days
	if c.Query("days") != "" {
		// Parse days
	}

	snapshots, err := h.analyticsService.GetDashboardHistory(organizationID, days)
	if err != nil {
		InternalServerError(c, "HISTORY_RETRIEVAL_FAILED", "Failed to retrieve dashboard history", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"history": snapshots}, "Dashboard history retrieved successfully")
}

// Heatmap endpoints

// GenerateRiskHeatmap generates a risk heatmap for an organization
func (h *AnalyticsHandler) GenerateRiskHeatmap(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	heatmapType := c.DefaultQuery("type", "comprehensive")
	timeRange := c.DefaultQuery("range", "30d")

	heatmapData, err := h.analyticsService.GenerateRiskHeatmap(organizationID, heatmapType, timeRange)
	if err != nil {
		InternalServerError(c, "HEATMAP_GENERATION_FAILED", "Failed to generate heatmap", err)
		return
	}

	SuccessResponse(c, http.StatusOK, heatmapData, "Heatmap generated successfully")
}

// GetHeatmapHotspots returns hotspots from heatmap
func (h *AnalyticsHandler) GetHeatmapHotspots(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	heatmapType := c.DefaultQuery("type", "comprehensive")
	timeRange := c.DefaultQuery("range", "30d")

	heatmapData, err := h.analyticsService.GenerateRiskHeatmap(organizationID, heatmapType, timeRange)
	if err != nil {
		InternalServerError(c, "HEATMAP_GENERATION_FAILED", "Failed to generate heatmap", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"hotspots": heatmapData.Hotspots}, "Hotspots retrieved successfully")
}

// GetRiskDistribution returns risk distribution from heatmap
func (h *AnalyticsHandler) GetRiskDistribution(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	heatmapType := c.DefaultQuery("type", "comprehensive")
	timeRange := c.DefaultQuery("range", "30d")

	heatmapData, err := h.analyticsService.GenerateRiskHeatmap(organizationID, heatmapType, timeRange)
	if err != nil {
		InternalServerError(c, "HEATMAP_GENERATION_FAILED", "Failed to generate heatmap", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"risk_distribution": heatmapData.RiskDistribution}, "Risk distribution retrieved successfully")
}

// GetHeatmapTrends returns trends from heatmap
func (h *AnalyticsHandler) GetHeatmapTrends(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	heatmapType := c.DefaultQuery("type", "comprehensive")
	timeRange := c.DefaultQuery("range", "30d")

	heatmapData, err := h.analyticsService.GenerateRiskHeatmap(organizationID, heatmapType, timeRange)
	if err != nil {
		InternalServerError(c, "HEATMAP_GENERATION_FAILED", "Failed to generate heatmap", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"trends": heatmapData.Trends}, "Trends retrieved successfully")
}

// GetHeatmapRecommendations returns recommendations from heatmap
func (h *AnalyticsHandler) GetHeatmapRecommendations(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	heatmapType := c.DefaultQuery("type", "comprehensive")
	timeRange := c.DefaultQuery("range", "30d")

	heatmapData, err := h.analyticsService.GenerateRiskHeatmap(organizationID, heatmapType, timeRange)
	if err != nil {
		InternalServerError(c, "HEATMAP_GENERATION_FAILED", "Failed to generate heatmap", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"recommendations": heatmapData.Recommendations}, "Recommendations retrieved successfully")
}

// Maturity endpoints

// CalculateMaturityScore calculates maturity score for an organization
func (h *AnalyticsHandler) CalculateMaturityScore(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	score, err := h.analyticsService.CalculateMaturityScore(organizationID)
	if err != nil {
		InternalServerError(c, "MATURITY_CALCULATION_FAILED", "Failed to calculate maturity score", err)
		return
	}

	SuccessResponse(c, http.StatusOK, score, "Maturity score calculated successfully")
}

// GetMaturityBenchmark returns maturity benchmark
func (h *AnalyticsHandler) GetMaturityBenchmark(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	score, err := h.analyticsService.CalculateMaturityScore(organizationID)
	if err != nil {
		InternalServerError(c, "MATURITY_CALCULATION_FAILED", "Failed to calculate maturity score", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"benchmark": score.IndustryBenchmark}, "Benchmark retrieved successfully")
}

// GetImprovementRoadmap returns improvement roadmap
func (h *AnalyticsHandler) GetImprovementRoadmap(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	score, err := h.analyticsService.CalculateMaturityScore(organizationID)
	if err != nil {
		InternalServerError(c, "MATURITY_CALCULATION_FAILED", "Failed to calculate maturity score", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"roadmap": score.ImprovementRoadmap}, "Roadmap retrieved successfully")
}

// GetMaturityTrends returns maturity trends
func (h *AnalyticsHandler) GetMaturityTrends(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	score, err := h.analyticsService.CalculateMaturityScore(organizationID)
	if err != nil {
		InternalServerError(c, "MATURITY_CALCULATION_FAILED", "Failed to calculate maturity score", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"trends": score.Trends}, "Trends retrieved successfully")
}

// GetDimensionScores returns dimension scores
func (h *AnalyticsHandler) GetDimensionScores(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	score, err := h.analyticsService.CalculateMaturityScore(organizationID)
	if err != nil {
		InternalServerError(c, "MATURITY_CALCULATION_FAILED", "Failed to calculate maturity score", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"dimensions": score.DimensionScores}, "Dimension scores retrieved successfully")
}

// Compliance endpoints

// GenerateComplianceReport generates compliance report
func (h *AnalyticsHandler) GenerateComplianceReport(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")
	reportType := c.DefaultQuery("type", "full")
	reportPeriod := c.DefaultQuery("period", "quarterly")

	report, err := h.analyticsService.GenerateComplianceReport(organizationID, framework, reportType, reportPeriod)
	if err != nil {
		InternalServerError(c, "COMPLIANCE_REPORT_GENERATION_FAILED", "Failed to generate compliance report", err)
		return
	}

	SuccessResponse(c, http.StatusOK, report, "Compliance report generated successfully")
}

// GetComplianceScore returns compliance score
func (h *AnalyticsHandler) GetComplianceScore(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")
	reportType := c.DefaultQuery("type", "full")
	reportPeriod := c.DefaultQuery("period", "quarterly")

	report, err := h.analyticsService.GenerateComplianceReport(organizationID, framework, reportType, reportPeriod)
	if err != nil {
		InternalServerError(c, "COMPLIANCE_REPORT_GENERATION_FAILED", "Failed to generate compliance report", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"score": report.OverallScore, "level": report.ComplianceLevel}, "Compliance score retrieved successfully")
}

// GetComplianceFindings returns compliance findings
func (h *AnalyticsHandler) GetComplianceFindings(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")
	reportType := c.DefaultQuery("type", "full")
	reportPeriod := c.DefaultQuery("period", "quarterly")

	report, err := h.analyticsService.GenerateComplianceReport(organizationID, framework, reportType, reportPeriod)
	if err != nil {
		InternalServerError(c, "COMPLIANCE_REPORT_GENERATION_FAILED", "Failed to generate compliance report", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"findings": report.Findings}, "Findings retrieved successfully")
}

// GetComplianceRecommendations returns compliance recommendations
func (h *AnalyticsHandler) GetComplianceRecommendations(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")
	reportType := c.DefaultQuery("type", "full")
	reportPeriod := c.DefaultQuery("period", "quarterly")

	report, err := h.analyticsService.GenerateComplianceReport(organizationID, framework, reportType, reportPeriod)
	if err != nil {
		InternalServerError(c, "COMPLIANCE_REPORT_GENERATION_FAILED", "Failed to generate compliance report", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"recommendations": report.Recommendations}, "Recommendations retrieved successfully")
}

// GetComplianceEvidence returns compliance evidence
func (h *AnalyticsHandler) GetComplianceEvidence(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")
	reportType := c.DefaultQuery("type", "full")
	reportPeriod := c.DefaultQuery("period", "quarterly")

	report, err := h.analyticsService.GenerateComplianceReport(organizationID, framework, reportType, reportPeriod)
	if err != nil {
		InternalServerError(c, "COMPLIANCE_REPORT_GENERATION_FAILED", "Failed to generate compliance report", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"evidence": report.EvidenceItems}, "Evidence retrieved successfully")
}

// GetExecutiveSummary returns executive summary
func (h *AnalyticsHandler) GetExecutiveSummary(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		BadRequest(c, "INVALID_UUID", "Invalid organization ID format", err.Error())
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")
	reportType := c.DefaultQuery("type", "full")
	reportPeriod := c.DefaultQuery("period", "quarterly")

	report, err := h.analyticsService.GenerateComplianceReport(organizationID, framework, reportType, reportPeriod)
	if err != nil {
		InternalServerError(c, "COMPLIANCE_REPORT_GENERATION_FAILED", "Failed to generate compliance report", err)
		return
	}

	SuccessResponse(c, http.StatusOK, gin.H{"executive_summary": report.ExecutiveSummary}, "Executive summary retrieved successfully")
}
