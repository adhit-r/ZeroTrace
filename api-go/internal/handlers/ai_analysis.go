package handlers

import (
	"net/http"

	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
)

// AIAnalysisHandler handles AI-powered vulnerability analysis endpoints
type AIAnalysisHandler struct {
	aiService *services.AIService
}

// NewAIAnalysisHandler creates a new AI analysis handler
func NewAIAnalysisHandler(aiService *services.AIService) *AIAnalysisHandler {
	return &AIAnalysisHandler{
		aiService: aiService,
	}
}

// AnalyzeVulnerabilityComprehensive performs comprehensive AI analysis
func (h *AIAnalysisHandler) AnalyzeVulnerabilityComprehensive(c *gin.Context) {
	vulnerabilityID := c.Param("id")
	if vulnerabilityID == "" {
		BadRequest(c, "MISSING_VULNERABILITY_ID", "Vulnerability ID is required", nil)
		return
	}

	// Create vulnerability data structure from provided ID
	// Note: Full implementation would fetch complete vulnerability details from database
	vulnData := services.VulnerabilityData{
		ID:          vulnerabilityID,
		Title:       "Vulnerability Analysis",
		Description: "Comprehensive AI analysis requested",
	}

	// Get organization context if available
	var orgContext *services.OrganizationContext
	if orgID, exists := c.Get("company_id"); exists {
		orgContext = &services.OrganizationContext{
			OrganizationID: orgID.(string),
		}
	}

	// Call AI service
	analysis, err := h.aiService.AnalyzeVulnerabilityComprehensive(vulnData, orgContext)
	if err != nil {
		InternalServerError(c, "AI_ANALYSIS_FAILED", "Failed to perform AI analysis", err)
		return
	}

	SuccessResponse(c, http.StatusOK, analysis, "Comprehensive analysis completed")
}

// AnalyzeVulnerabilityTrends analyzes vulnerability trends using AI
func (h *AIAnalysisHandler) AnalyzeVulnerabilityTrends(c *gin.Context) {
	// This endpoint requires vulnerability data from database
	// Return error indicating database integration is required
	ErrorResponse(c, http.StatusNotImplemented, "DATABASE_INTEGRATION_REQUIRED",
		"Trend analysis requires vulnerability data from database. This endpoint needs to be connected to vulnerability repository.",
		nil)
}

// GetExploitIntelligence retrieves exploit intelligence for a CVE
func (h *AIAnalysisHandler) GetExploitIntelligence(c *gin.Context) {
	cveID := c.Param("cve_id")
	if cveID == "" {
		BadRequest(c, "MISSING_CVE_ID", "CVE ID is required", nil)
		return
	}

	packageName := c.Query("package_name")

	// Call AI service
	intelligence, err := h.aiService.GetExploitIntelligence(cveID, packageName)
	if err != nil {
		InternalServerError(c, "EXPLOIT_INTELLIGENCE_FAILED", "Failed to retrieve exploit intelligence", err)
		return
	}

	SuccessResponse(c, http.StatusOK, intelligence, "Exploit intelligence retrieved")
}

// GetPredictiveAnalysis performs predictive analysis on vulnerabilities
func (h *AIAnalysisHandler) GetPredictiveAnalysis(c *gin.Context) {
	vulnerabilityID := c.Param("id")
	if vulnerabilityID == "" {
		BadRequest(c, "MISSING_VULNERABILITY_ID", "Vulnerability ID is required", nil)
		return
	}

	// Create vulnerability data
	vulnData := services.VulnerabilityData{
		ID:          vulnerabilityID,
		Title:       "Predictive Analysis Request",
		Description: "AI-powered predictive analysis requested",
	}

	// Get organization context if available
	var orgContext *services.OrganizationContext
	if orgID, exists := c.Get("company_id"); exists {
		orgContext = &services.OrganizationContext{
			OrganizationID: orgID.(string),
		}
	}

	// Call AI service
	predictions, err := h.aiService.GetPredictiveAnalysis(vulnData, orgContext)
	if err != nil {
		InternalServerError(c, "PREDICTIVE_ANALYSIS_FAILED", "Failed to perform predictive analysis", err)
		return
	}

	SuccessResponse(c, http.StatusOK, predictions, "Predictive analysis completed")
}

// GetRemediationPlan generates AI-powered remediation plans
func (h *AIAnalysisHandler) GetRemediationPlan(c *gin.Context) {
	vulnerabilityID := c.Param("id")
	if vulnerabilityID == "" {
		BadRequest(c, "MISSING_VULNERABILITY_ID", "Vulnerability ID is required", nil)
		return
	}

	// Create vulnerability data
	vulnData := services.VulnerabilityData{
		ID:          vulnerabilityID,
		Title:       "Remediation Plan Request",
		Description: "AI-powered remediation plan requested",
	}

	// Get organization context if available
	var orgContext *services.OrganizationContext
	if orgID, exists := c.Get("company_id"); exists {
		orgContext = &services.OrganizationContext{
			OrganizationID: orgID.(string),
		}
	}

	// Call AI service
	plan, err := h.aiService.GetRemediationPlan(vulnData, orgContext)
	if err != nil {
		InternalServerError(c, "REMEDIATION_PLAN_FAILED", "Failed to generate remediation plan", err)
		return
	}

	SuccessResponse(c, http.StatusOK, plan, "Remediation plan generated")
}

// GetBulkAnalysis performs bulk analysis on multiple vulnerabilities
func (h *AIAnalysisHandler) GetBulkAnalysis(c *gin.Context) {
	var req struct {
		VulnerabilityIDs []string `json:"vulnerability_ids" binding:"required"`
		AnalysisType     string   `json:"analysis_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	// Bulk analysis requires fetching vulnerabilities from database
	// Return error indicating database integration is required
	ErrorResponse(c, http.StatusNotImplemented, "DATABASE_INTEGRATION_REQUIRED",
		"Bulk analysis requires vulnerability data from database. This endpoint needs to be connected to vulnerability repository.",
		map[string]interface{}{
			"requested_count": len(req.VulnerabilityIDs),
			"analysis_type":   req.AnalysisType,
		})
}
