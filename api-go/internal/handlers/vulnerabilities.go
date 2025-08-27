package handlers

import (
	"net/http"
	"time"

	"zerotrace/api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetVulnerabilities returns software vulnerabilities found on the system
func GetVulnerabilities(c *gin.Context) {
	// Mock software vulnerability data - real implementation would query database
	vulnerabilities := []models.Vulnerability{
		{
			ID:               uuid.New(),
			ScanID:           uuid.New(),
			CompanyID:        uuid.New(),
			Type:             "software",
			Severity:         "Critical",
			Title:            "Adobe Acrobat Reader DC - CVE-2023-21608",
			Description:      "Use-after-free vulnerability in Adobe Acrobat Reader DC allows remote code execution",
			CVEID:            "CVE-2023-21608",
			CVSSScore:        &[]float64{9.8}[0],
			CVSSVector:       "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
			PackageName:      "Adobe Acrobat Reader DC",
			PackageVersion:   "23.008.20470",
			Location:         "/Applications/Adobe Acrobat Reader DC/Adobe Acrobat Reader.app",
			Remediation:      "Update to version 23.008.20533 or later",
			References:       []string{"https://helpx.adobe.com/security/products/acrobat/apsb23-01.html"},
			AffectedVersions: []string{"<23.008.20533"},
			PatchedVersions:  []string{"23.008.20533", "24.001.20604"},
			ExploitAvailable: true,
			ExploitCount:     2,
			Status:           "open",
			Priority:         "high",
			Notes:            "Critical vulnerability affecting PDF processing",
			EnrichmentData: map[string]any{
				"installation_date": "2023-12-15",
				"last_used":         "2024-01-20",
				"file_size":         "245MB",
				"vendor":            "Adobe Inc.",
			},
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:               uuid.New(),
			ScanID:           uuid.New(),
			CompanyID:        uuid.New(),
			Type:             "software",
			Severity:         "High",
			Title:            "Google Chrome - CVE-2023-7024",
			Description:      "Type confusion vulnerability in V8 JavaScript engine allows remote code execution",
			CVEID:            "CVE-2023-7024",
			CVSSScore:        &[]float64{8.8}[0],
			CVSSVector:       "CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:U/C:H/I:H/A:H",
			PackageName:      "Google Chrome",
			PackageVersion:   "120.0.6099.109",
			Location:         "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			Remediation:      "Update to version 120.0.6099.129 or later",
			References:       []string{"https://chromereleases.googleblog.com/2023/12/stable-channel-update-for-desktop_19.html"},
			AffectedVersions: []string{"<120.0.6099.129"},
			PatchedVersions:  []string{"120.0.6099.129", "121.0.6167.85"},
			ExploitAvailable: true,
			ExploitCount:     1,
			Status:           "open",
			Priority:         "high",
			Notes:            "Browser vulnerability affecting web browsing security",
			EnrichmentData: map[string]any{
				"installation_date": "2023-11-10",
				"last_used":         "2024-01-21",
				"file_size":         "156MB",
				"vendor":            "Google LLC",
			},
			CreatedAt: time.Now().Add(-12 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:               uuid.New(),
			ScanID:           uuid.New(),
			CompanyID:        uuid.New(),
			Type:             "software",
			Severity:         "Medium",
			Title:            "7-Zip - CVE-2023-31102",
			Description:      "Buffer overflow vulnerability in 7-Zip allows local privilege escalation",
			CVEID:            "CVE-2023-31102",
			CVSSScore:        &[]float64{7.8}[0],
			CVSSVector:       "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:H/I:H/A:H",
			PackageName:      "7-Zip",
			PackageVersion:   "23.01",
			Location:         "/usr/local/bin/7z",
			Remediation:      "Update to version 23.02 or later",
			References:       []string{"https://www.7-zip.org/history.txt"},
			AffectedVersions: []string{"<23.02"},
			PatchedVersions:  []string{"23.02", "24.01"},
			ExploitAvailable: false,
			ExploitCount:     0,
			Status:           "open",
			Priority:         "medium",
			Notes:            "Compression utility vulnerability",
			EnrichmentData: map[string]any{
				"installation_date": "2023-08-05",
				"last_used":         "2024-01-18",
				"file_size":         "1.2MB",
				"vendor":            "Igor Pavlov",
			},
			CreatedAt: time.Now().Add(-6 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:               uuid.New(),
			ScanID:           uuid.New(),
			CompanyID:        uuid.New(),
			Type:             "software",
			Severity:         "Medium",
			Title:            "VLC Media Player - CVE-2023-29547",
			Description:      "Heap-based buffer overflow in VLC media player allows remote code execution",
			CVEID:            "CVE-2023-29547",
			CVSSScore:        &[]float64{6.5}[0],
			CVSSVector:       "CVSS:3.1/AV:N/AC:L/PR:N/UI:R/S:U/C:N/I:N/A:H",
			PackageName:      "VLC Media Player",
			PackageVersion:   "3.0.18",
			Location:         "/Applications/VLC.app/Contents/MacOS/VLC",
			Remediation:      "Update to version 3.0.19 or later",
			References:       []string{"https://www.videolan.org/security/sa2301.html"},
			AffectedVersions: []string{"<3.0.19"},
			PatchedVersions:  []string{"3.0.19", "3.0.20"},
			ExploitAvailable: false,
			ExploitCount:     0,
			Status:           "open",
			Priority:         "medium",
			Notes:            "Media player vulnerability affecting file processing",
			EnrichmentData: map[string]any{
				"installation_date": "2023-09-22",
				"last_used":         "2024-01-19",
				"file_size":         "89MB",
				"vendor":            "VideoLAN",
			},
			CreatedAt: time.Now().Add(-3 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:               uuid.New(),
			ScanID:           uuid.New(),
			CompanyID:        uuid.New(),
			Type:             "software",
			Severity:         "Low",
			Title:            "Notepad++ - CVE-2023-40031",
			Description:      "Information disclosure vulnerability in Notepad++",
			CVEID:            "CVE-2023-40031",
			CVSSScore:        &[]float64{3.3}[0],
			CVSSVector:       "CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:L/I:N/A:N",
			PackageName:      "Notepad++",
			PackageVersion:   "8.6.2",
			Location:         "/Applications/Notepad++.app/Contents/MacOS/Notepad++",
			Remediation:      "Update to version 8.6.3 or later",
			References:       []string{"https://github.com/notepad-plus-plus/notepad-plus-plus/releases"},
			AffectedVersions: []string{"<8.6.3"},
			PatchedVersions:  []string{"8.6.3", "8.6.4"},
			ExploitAvailable: false,
			ExploitCount:     0,
			Status:           "open",
			Priority:         "low",
			Notes:            "Text editor vulnerability with minimal impact",
			EnrichmentData: map[string]any{
				"installation_date": "2023-10-15",
				"last_used":         "2024-01-21",
				"file_size":         "12MB",
				"vendor":            "Notepad++ Team",
			},
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"vulnerabilities": vulnerabilities,
			"total":           len(vulnerabilities),
			"critical":        1,
			"high":            1,
			"medium":          2,
			"low":             1,
		},
		"message": "Software vulnerabilities retrieved successfully",
	})
}

// GetVulnerabilityDetails returns details of a specific software vulnerability
func GetVulnerabilityDetails(c *gin.Context) {
	vulnID := c.Param("id")

	// Mock detailed vulnerability data
	vulnerability := models.Vulnerability{
		ID:             uuid.MustParse(vulnID),
		ScanID:         uuid.New(),
		CompanyID:      uuid.New(),
		Type:           "software",
		Severity:       "Critical",
		Title:          "Adobe Acrobat Reader DC - CVE-2023-21608",
		Description:    "A critical use-after-free vulnerability exists in Adobe Acrobat Reader DC that could allow an attacker to execute arbitrary code on the affected system. This vulnerability affects the PDF parsing engine and can be triggered when processing specially crafted PDF files.",
		CVEID:          "CVE-2023-21608",
		CVSSScore:      &[]float64{9.8}[0],
		CVSSVector:     "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
		PackageName:    "Adobe Acrobat Reader DC",
		PackageVersion: "23.008.20470",
		Location:       "/Applications/Adobe Acrobat Reader DC/Adobe Acrobat Reader.app",
		Remediation:    "1. Update Adobe Acrobat Reader DC to version 23.008.20533 or later\n2. Enable automatic updates in Adobe Reader preferences\n3. Consider using alternative PDF readers if update is not possible",
		References: []string{
			"https://helpx.adobe.com/security/products/acrobat/apsb23-01.html",
			"https://nvd.nist.gov/vuln/detail/CVE-2023-21608",
			"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-21608",
		},
		AffectedVersions: []string{"<23.008.20533"},
		PatchedVersions:  []string{"23.008.20533", "24.001.20604"},
		ExploitAvailable: true,
		ExploitCount:     2,
		Status:           "open",
		Priority:         "high",
		Notes:            "Critical vulnerability affecting PDF processing. Multiple exploit attempts detected in the wild.",
		EnrichmentData: map[string]any{
			"installation_date":  "2023-12-15",
			"last_used":          "2024-01-20",
			"file_size":          "245MB",
			"vendor":             "Adobe Inc.",
			"update_available":   true,
			"latest_version":     "24.001.20604",
			"download_url":       "https://get.adobe.com/reader/",
			"affected_platforms": []string{"Windows", "macOS", "Linux"},
		},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vulnerability,
		"message": "Software vulnerability details retrieved successfully",
	})
}

// GetVulnerabilityStats returns statistics about software vulnerabilities
func GetVulnerabilityStats(c *gin.Context) {
	stats := gin.H{
		"total_vulnerabilities": 6,
		"critical":              1,
		"high":                  1,
		"medium":                2,
		"low":                   1,
		"open":                  6,
		"resolved":              0,
		"software_categories": gin.H{
			"browsers":      1,
			"pdf_readers":   1,
			"media_players": 1,
			"utilities":     2,
			"office_tools":  0,
			"development":   1,
		},
		"trends": gin.H{
			"last_7_days": gin.H{
				"new":      2,
				"resolved": 0,
				"critical": 1,
			},
			"last_30_days": gin.H{
				"new":      4,
				"resolved": 1,
				"critical": 1,
			},
		},
		"top_vulnerable_software": []gin.H{
			{"name": "Adobe Acrobat Reader DC", "vulnerabilities": 1, "severity": "Critical"},
			{"name": "Google Chrome", "vulnerabilities": 1, "severity": "High"},
			{"name": "7-Zip", "vulnerabilities": 1, "severity": "Medium"},
			{"name": "VLC Media Player", "vulnerabilities": 1, "severity": "Medium"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
		"message": "Software vulnerability statistics retrieved successfully",
	})
}
