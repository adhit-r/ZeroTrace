package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
)

// GetDashboardOverview returns dashboard overview data
func GetDashboardOverview(c *gin.Context) {
	// TODO: Implement actual dashboard data aggregation
	overview := map[string]any{
		"total_scans":              150,
		"active_scans":             5,
		"total_vulnerabilities":    1250,
		"critical_vulnerabilities": 45,
		"high_vulnerabilities":     120,
		"medium_vulnerabilities":   350,
		"low_vulnerabilities":      735,
		"recent_scans": []map[string]any{
			{
				"id":              "scan-1",
				"repository":      "https://github.com/example/repo1",
				"status":          "completed",
				"vulnerabilities": 25,
				"created_at":      time.Now().Add(-2 * time.Hour),
			},
			{
				"id":              "scan-2",
				"repository":      "https://github.com/example/repo2",
				"status":          "scanning",
				"vulnerabilities": 0,
				"created_at":      time.Now().Add(-1 * time.Hour),
			},
		},
		"top_vulnerabilities": []map[string]any{
			{
				"cve_id":     "CVE-2023-1234",
				"title":      "SQL Injection Vulnerability",
				"severity":   "CRITICAL",
				"count":      15,
				"cvss_score": 9.8,
			},
			{
				"cve_id":     "CVE-2023-5678",
				"title":      "Cross-Site Scripting",
				"severity":   "HIGH",
				"count":      12,
				"cvss_score": 8.5,
			},
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      overview,
		Message:   "Dashboard overview retrieved successfully",
		Timestamp: time.Now(),
	})
}

// GetVulnerabilityTrends returns vulnerability trends data
func GetVulnerabilityTrends(c *gin.Context) {
	// TODO: Implement actual trends calculation
	trends := map[string]any{
		"daily_trends": []map[string]any{
			{
				"date":                  time.Now().AddDate(0, 0, -6).Format("2006-01-02"),
				"total_vulnerabilities": 120,
				"critical":              5,
				"high":                  15,
				"medium":                45,
				"low":                   55,
			},
			{
				"date":                  time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
				"total_vulnerabilities": 135,
				"critical":              8,
				"high":                  18,
				"medium":                52,
				"low":                   57,
			},
			{
				"date":                  time.Now().AddDate(0, 0, -4).Format("2006-01-02"),
				"total_vulnerabilities": 142,
				"critical":              12,
				"high":                  22,
				"medium":                58,
				"low":                   50,
			},
			{
				"date":                  time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
				"total_vulnerabilities": 138,
				"critical":              10,
				"high":                  20,
				"medium":                55,
				"low":                   53,
			},
			{
				"date":                  time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
				"total_vulnerabilities": 145,
				"critical":              15,
				"high":                  25,
				"medium":                60,
				"low":                   45,
			},
			{
				"date":                  time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
				"total_vulnerabilities": 150,
				"critical":              18,
				"high":                  28,
				"medium":                62,
				"low":                   42,
			},
			{
				"date":                  time.Now().Format("2006-01-02"),
				"total_vulnerabilities": 155,
				"critical":              20,
				"high":                  30,
				"medium":                65,
				"low":                   40,
			},
		},
		"severity_distribution": map[string]int{
			"critical": 20,
			"high":     30,
			"medium":   65,
			"low":      40,
		},
		"trend_analysis": map[string]any{
			"trend_direction":  "increasing",
			"trend_percentage": 12.5,
			"critical_trend":   "increasing",
			"high_trend":       "stable",
			"medium_trend":     "increasing",
			"low_trend":        "decreasing",
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      trends,
		Message:   "Vulnerability trends retrieved successfully",
		Timestamp: time.Now(),
	})
}
