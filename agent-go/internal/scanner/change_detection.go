package scanner

import (
	"fmt"
	"time"
)

// DetectChanges compares two network scan results to identify changes in the attack surface.
func (ns *NetworkScanner) DetectChanges(currentScan, previousScan *NetworkScanResult) *ChangeDetection {
	detection := &ChangeDetection{
		ScanID:       currentScan.ID,
		PreviousScan: previousScan.ID,
		Timestamp:    time.Now(),
		Changes:      []Change{},
	}

	// 1. Compare Open Ports and Vulnerabilities
	prevFindings := mapFindings(previousScan.NetworkFindings)
	currFindings := mapFindings(currentScan.NetworkFindings)

	// Check for new and modified findings
	for key, currFinding := range currFindings {
		if _, exists := prevFindings[key]; !exists {
			// New finding
			change := Change{
				Type:        fmt.Sprintf("new-%s", currFinding.FindingType),
				Severity:    currFinding.Severity,
				Description: fmt.Sprintf("New finding: %s on %s:%d", currFinding.Description, currFinding.Host, currFinding.Port),
				Entity:      key,
				NewValue:    "present",
				OldValue:    "absent",
				RiskImpact:  calculateRiskImpact(currFinding),
			}
			detection.Changes = append(detection.Changes, change)
			detection.RiskDelta += change.RiskImpact
		} else {
			// Compare attributes of existing finding if needed (e.g., service version change)
			delete(prevFindings, key) // Mark as seen
		}
	}

	// Check for fixed/removed findings
	for key, prevFinding := range prevFindings {
		change := Change{
			Type:        fmt.Sprintf("fixed-%s", prevFinding.FindingType),
			Severity:    "info",
			Description: fmt.Sprintf("Finding removed: %s on %s:%d", prevFinding.Description, prevFinding.Host, prevFinding.Port),
			Entity:      key,
			NewValue:    "absent",
			OldValue:    "present",
			RiskImpact:  -calculateRiskImpact(prevFinding), // Negative impact for risk reduction
		}
		detection.Changes = append(detection.Changes, change)
		detection.RiskDelta += change.RiskImpact
	}

	return detection
}

// mapFindings creates a map of findings for easy lookup.
// The key is a unique identifier for the finding (e.g., "host:port:finding_type:description").
func mapFindings(findings []NetworkFinding) map[string]NetworkFinding {
	m := make(map[string]NetworkFinding)
	for _, f := range findings {
		key := fmt.Sprintf("%s:%d:%s:%s", f.Host, f.Port, f.FindingType, f.Description)
		m[key] = f
	}
	return m
}

// calculateRiskImpact provides a simplified risk score for a finding.
func calculateRiskImpact(finding NetworkFinding) float64 {
	switch finding.Severity {
	case "critical":
		return 100.0
	case "high":
		return 75.0
	case "medium":
		return 40.0
	case "low":
		return 10.0
	default:
		return 5.0
	}
}

// ShouldAlert determines if a change detection report warrants an immediate alert.
func (cd *ChangeDetection) ShouldAlert() bool {
	for _, change := range cd.Changes {
		if change.Severity == "critical" || change.Severity == "high" {
			return true
		}
		if change.RiskImpact > 50 {
			return true
		}
	}
	return false
}

// GenerateReport creates a human-readable summary of the changes.
func (cd *ChangeDetection) GenerateReport() string {
	newVulns := 0
	fixedVulns := 0
	newPorts := 0
	closedPorts := 0

	for _, change := range cd.Changes {
		switch change.Type {
		case "new-vuln":
			newVulns++
		case "fixed-vuln":
			fixedVulns++
		case "new-port":
			newPorts++
		case "closed-port":
			closedPorts++
		}
	}

	return fmt.Sprintf("Attack Surface Change Report: %d new ports, %d closed ports, %d new vulnerabilities, %d fixed vulnerabilities. Risk Delta: %.2f",
		newPorts, closedPorts, newVulns, fixedVulns, cd.RiskDelta)
}
