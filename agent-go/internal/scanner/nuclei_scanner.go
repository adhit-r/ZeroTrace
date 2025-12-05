package scanner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
)

// NucleiScanner handles vulnerability scanning using Nuclei
type NucleiScanner struct{}

// NewNucleiScanner creates a new Nuclei scanner
func NewNucleiScanner() *NucleiScanner {
	return &NucleiScanner{}
}

// ScanTargets performs Nuclei vulnerability scanning on given targets using CLI
func (ns *NucleiScanner) ScanTargets(targets []string) ([]NetworkFinding, error) {
	// Use CLI method as the primary method
	return ns.ScanUsingCLI(targets)
}

// ScanTargetsWithCredentials performs authenticated Nuclei scanning
func (ns *NucleiScanner) ScanTargetsWithCredentials(targets []string, credentials map[string]string) ([]NetworkFinding, error) {
	// For authenticated scanning, we would pass credentials as headers or auth options
	// For now, use standard CLI scan
	// In production, you'd add: -H "Authorization: Bearer ..." or similar
	return ns.ScanUsingCLI(targets)
}

// ScanUsingCLI performs Nuclei scan using CLI (fallback method)
func (ns *NucleiScanner) ScanUsingCLI(targets []string) ([]NetworkFinding, error) {
	if len(targets) == 0 {
		return []NetworkFinding{}, nil
	}

	var findings []NetworkFinding

	// Build Nuclei command
	args := []string{
		"-json",
		"-silent",
		"-no-color",
		"-rate-limit", "150",
		"-timeout", "10",
	}

	// Add targets
	for _, target := range targets {
		args = append(args, "-target", target)
	}

	// Execute Nuclei CLI
	cmd := exec.Command("nuclei", args...)
	output, err := cmd.Output()
	if err != nil {
		// Nuclei may return non-zero exit code even with findings
		// Check if we got any output
		if len(output) == 0 {
			return nil, fmt.Errorf("nuclei scan failed: %w", err)
		}
	}

	// Parse JSON output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result nucleiResult
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			continue
		}

		// Convert Nuclei result to NetworkFinding
		finding := ns.convertNucleiResult(result)
		if finding != nil {
			findings = append(findings, *finding)
		}
	}

	return findings, nil
}

// convertNucleiResult converts Nuclei JSON result to NetworkFinding
func (ns *NucleiScanner) convertNucleiResult(result nucleiResult) *NetworkFinding {
	// Map Nuclei severity to our severity system
	severity := ns.mapNucleiSeverity(result.Info.Severity)

	// Extract host and port from matched-at URL
	host, port := ns.extractHostPort(result.MatchedAt)

	return &NetworkFinding{
		ID:             uuid.New(),
		FindingType:    "vuln",
		Severity:       severity,
		Host:           host,
		Port:           port,
		Protocol:       "tcp",
		ServiceName:    result.Info.Name,
		ServiceVersion: "",
		Banner:         "",
		Description:    result.Info.Description,
		Remediation:    ns.generateRemediation(result),
		DiscoveredAt:   time.Now(),
		Status:         "open",
		Metadata: map[string]interface{}{
			"template_id":       result.TemplateID,
			"template_path":     result.TemplatePath,
			"matched_at":        result.MatchedAt,
			"extracted_results": result.ExtractedResults,
			"curl_command":      result.CurlCommand,
		},
	}
}

// mapNucleiSeverity maps Nuclei severity levels to our system
func (ns *NucleiScanner) mapNucleiSeverity(nucleiSeverity string) string {
	nucleiSeverity = strings.ToLower(nucleiSeverity)
	switch nucleiSeverity {
	case "critical":
		return "critical"
	case "high":
		return "high"
	case "medium":
		return "medium"
	case "low":
		return "low"
	case "info":
		return "info"
	default:
		return "medium" // Default to medium if unknown
	}
}

// extractHostPort extracts host and port from URL
func (ns *NucleiScanner) extractHostPort(url string) (string, int) {
	// Simple extraction - in production, use proper URL parsing
	// Format: http://host:port/path or https://host:port/path
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	parts := strings.Split(url, "/")
	hostPort := parts[0]

	hostPortParts := strings.Split(hostPort, ":")
	host := hostPortParts[0]

	port := 80
	if len(hostPortParts) > 1 {
		fmt.Sscanf(hostPortParts[1], "%d", &port)
	} else if strings.HasPrefix(url, "https://") {
		port = 443
	}

	return host, port
}

// generateRemediation generates remediation advice based on Nuclei result
func (ns *NucleiScanner) generateRemediation(result nucleiResult) string {
	if result.Info.Remediation != "" {
		return result.Info.Remediation
	}

	// Generate generic remediation based on severity
	severity := strings.ToLower(result.Info.Severity)
	switch severity {
	case "critical", "high":
		return "Address this vulnerability immediately. Apply security patches or updates. If no patch is available, implement compensating controls."
	case "medium":
		return "Review and address this vulnerability. Apply patches or updates when available."
	case "low", "info":
		return "Review this finding. Consider addressing if it poses a risk to your environment."
	default:
		return "Review this finding and determine appropriate remediation."
	}
}

// nucleiResult represents a Nuclei JSON output result
type nucleiResult struct {
	TemplateID       string     `json:"template-id"`
	TemplatePath     string     `json:"template-path"`
	Info             nucleiInfo `json:"info"`
	Type             string     `json:"type"`
	Host             string     `json:"host"`
	MatchedAt        string     `json:"matched-at"`
	ExtractedResults []string   `json:"extracted-results"`
	CurlCommand      string     `json:"curl-command"`
	Request          string     `json:"request"`
	Response         string     `json:"response"`
	IP               string     `json:"ip"`
	Timestamp        string     `json:"timestamp"`
}

// nucleiInfo represents Nuclei template info
type nucleiInfo struct {
	Name           string               `json:"name"`
	Author         []string             `json:"author"`
	Tags           []string             `json:"tags"`
	Description    string               `json:"description"`
	Severity       string               `json:"severity"`
	Remediation    string               `json:"remediation"`
	Reference      []string             `json:"reference"`
	Classification nucleiClassification `json:"classification"`
}

// nucleiClassification represents Nuclei classification
type nucleiClassification struct {
	CVEID     []string `json:"cve-id"`
	CWEID     []string `json:"cwe-id"`
	CVSSScore float64  `json:"cvss-score"`
	EPSSScore float64  `json:"epss-score"`
	CPEScore  float64  `json:"cpe-score"`
}
