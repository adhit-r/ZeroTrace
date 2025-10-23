package handlers

import (
	"net/http"

	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ComplianceHandler handles compliance reporting API endpoints
type ComplianceHandler struct {
	complianceService *services.ComplianceService
}

// NewComplianceHandler creates a new ComplianceHandler
func NewComplianceHandler(complianceService *services.ComplianceService) *ComplianceHandler {
	return &ComplianceHandler{
		complianceService: complianceService,
	}
}

// GenerateComplianceReport generates a compliance report for an organization
func (h *ComplianceHandler) GenerateComplianceReport(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	// Get query parameters
	framework := c.DefaultQuery("framework", "SOC2")
	reportType := c.DefaultQuery("type", "full")
	reportPeriod := c.DefaultQuery("period", "quarterly")

	// Generate compliance report
	report, err := h.complianceService.GenerateComplianceReport(organizationID, framework, reportType, reportPeriod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "COMPLIANCE_REPORT_FAILED",
				"message": "Failed to generate compliance report",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    report,
	})
}

// GetComplianceScore gets compliance score for an organization
func (h *ComplianceHandler) GetComplianceScore(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")

	// Generate compliance report to get score
	report, err := h.complianceService.GenerateComplianceReport(organizationID, framework, "summary", "current")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "COMPLIANCE_SCORE_FAILED",
				"message": "Failed to generate compliance score",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":  organizationID,
			"framework":        framework,
			"overall_score":    report.OverallScore,
			"compliance_level": report.ComplianceLevel,
			"control_scores":   report.ControlScores,
			"generated_at":     report.GeneratedAt,
		},
	})
}

// GetComplianceFindings gets compliance findings for an organization
func (h *ComplianceHandler) GetComplianceFindings(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")

	// Generate compliance report to get findings
	report, err := h.complianceService.GenerateComplianceReport(organizationID, framework, "findings", "current")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "COMPLIANCE_FINDINGS_FAILED",
				"message": "Failed to generate compliance findings",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id": organizationID,
			"framework":       framework,
			"findings":        report.Findings,
			"total_findings":  len(report.Findings),
			"critical_count":  h.countFindingsBySeverity(report.Findings, "critical"),
			"high_count":      h.countFindingsBySeverity(report.Findings, "high"),
			"generated_at":    report.GeneratedAt,
		},
	})
}

// GetComplianceRecommendations gets compliance recommendations for an organization
func (h *ComplianceHandler) GetComplianceRecommendations(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")

	// Generate compliance report to get recommendations
	report, err := h.complianceService.GenerateComplianceReport(organizationID, framework, "recommendations", "current")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "COMPLIANCE_RECOMMENDATIONS_FAILED",
				"message": "Failed to generate compliance recommendations",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":       organizationID,
			"framework":             framework,
			"recommendations":       report.Recommendations,
			"total_recommendations": len(report.Recommendations),
			"high_priority":         h.filterRecommendationsByPriority(report.Recommendations, "high"),
			"generated_at":          report.GeneratedAt,
		},
	})
}

// GetComplianceEvidence gets compliance evidence for an organization
func (h *ComplianceHandler) GetComplianceEvidence(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")

	// Generate compliance report to get evidence
	report, err := h.complianceService.GenerateComplianceReport(organizationID, framework, "evidence", "current")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "COMPLIANCE_EVIDENCE_FAILED",
				"message": "Failed to generate compliance evidence",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id": organizationID,
			"framework":       framework,
			"evidence_items":  report.EvidenceItems,
			"total_evidence":  len(report.EvidenceItems),
			"valid_evidence":  h.countValidEvidence(report.EvidenceItems),
			"generated_at":    report.GeneratedAt,
		},
	})
}

// GetExecutiveSummary gets executive summary for compliance
func (h *ComplianceHandler) GetExecutiveSummary(c *gin.Context) {
	organizationIDStr := c.Param("id")
	organizationID, err := uuid.Parse(organizationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "INVALID_UUID",
				"message": "Invalid organization ID format",
				"details": err.Error(),
			},
		})
		return
	}

	framework := c.DefaultQuery("framework", "SOC2")

	// Generate compliance report to get executive summary
	report, err := h.complianceService.GenerateComplianceReport(organizationID, framework, "executive", "current")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "EXECUTIVE_SUMMARY_FAILED",
				"message": "Failed to generate executive summary",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"organization_id":   organizationID,
			"framework":         framework,
			"executive_summary": report.ExecutiveSummary,
			"overall_score":     report.OverallScore,
			"compliance_level":  report.ComplianceLevel,
			"generated_at":      report.GeneratedAt,
		},
	})
}

// Helper methods

func (h *ComplianceHandler) countFindingsBySeverity(findings []services.ComplianceFinding, severity string) int {
	count := 0
	for _, finding := range findings {
		if finding.Severity == severity {
			count++
		}
	}
	return count
}

func (h *ComplianceHandler) filterRecommendationsByPriority(recommendations []services.ComplianceRecommendation, priority string) []services.ComplianceRecommendation {
	var filtered []services.ComplianceRecommendation
	for _, rec := range recommendations {
		if rec.Priority == priority {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}

func (h *ComplianceHandler) countValidEvidence(evidenceItems []services.EvidenceItem) int {
	count := 0
	for _, evidence := range evidenceItems {
		if evidence.Status == "valid" {
			count++
		}
	}
	return count
}
