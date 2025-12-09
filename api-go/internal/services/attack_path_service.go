package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type AttackPathService struct {
	db *sql.DB
}

func NewAttackPathService(db *sql.DB) *AttackPathService {
	return &AttackPathService{db: db}
}

type AttackPath struct {
	PathID             string       `json:"path_id"`
	Name               string       `json:"name"`
	Steps              []AttackStep `json:"steps"`
	TotalLikelihood    float64      `json:"total_likelihood"`
	TotalImpact        float64      `json:"total_impact"`
	CriticalityScore   float64      `json:"criticality_score"`
	MitigationPriority string       `json:"mitigation_priority"`
	DetectionPoints    []string     `json:"detection_points"`
	PreventionControls []string     `json:"prevention_controls"`
	CreatedAt          string       `json:"created_at"`
}

type AttackStep struct {
	StepNumber          int      `json:"step_number"`
	Action              string   `json:"action"`
	Target              string   `json:"target"`
	TargetIP            string   `json:"target_ip,omitempty"`
	TargetHostname      string   `json:"target_hostname,omitempty"`
	Technique           string   `json:"technique"`
	TechniqueID         string   `json:"technique_id,omitempty"` // MITRE ATT&CK ID
	Likelihood          float64  `json:"likelihood"`
	Impact              float64  `json:"impact"`
	DetectionDifficulty string   `json:"detection_difficulty"`
	StepType            string   `json:"step_type"` // initial_access, lateral_movement, privilege_escalation, data_exfiltration, persistence
	CVEID               string   `json:"cve_id,omitempty"`
	VulnerabilityID     string   `json:"vulnerability_id,omitempty"`
	Proof               string   `json:"proof,omitempty"` // Exploit command or evidence
	MitigationControls  []string `json:"mitigation_controls,omitempty"`
}

// GetAttackPaths retrieves all attack paths for an organization
func (s *AttackPathService) GetAttackPaths(organizationID string) ([]AttackPath, error) {
	// Query vulnerabilities and network scans to build attack paths
	query := `
		SELECT 
			v.id,
			v.cve_id,
			v.severity,
			v.title,
			v.description,
			v.risk_score,
			v.agent_id,
			a.hostname,
			a.os,
			a.metadata
		FROM vulnerabilities v
		LEFT JOIN agents a ON v.agent_id = a.id
		WHERE v.severity IN ('critical', 'high')
		ORDER BY v.risk_score DESC
		LIMIT 100
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query vulnerabilities: %w", err)
	}
	defer rows.Close()

	var vulnerabilities []struct {
		ID          string
		CVEID       sql.NullString
		Severity    string
		Title       string
		Description sql.NullString
		RiskScore   float64
		AgentID     sql.NullString
		Hostname    sql.NullString
		OS          sql.NullString
		Metadata    sql.NullString
	}

	for rows.Next() {
		var v struct {
			ID          string
			CVEID       sql.NullString
			Severity    string
			Title       string
			Description sql.NullString
			RiskScore   float64
			AgentID     sql.NullString
			Hostname    sql.NullString
			OS          sql.NullString
			Metadata    sql.NullString
		}

		if err := rows.Scan(&v.ID, &v.CVEID, &v.Severity, &v.Title, &v.Description, &v.RiskScore, &v.AgentID, &v.Hostname, &v.OS, &v.Metadata); err != nil {
			log.Printf("Error scanning vulnerability: %v", err)
			continue
		}

		vulnerabilities = append(vulnerabilities, v)
	}

	// Build attack paths from vulnerabilities
	paths := s.buildAttackPathsFromVulnerabilities(vulnerabilities)

	return paths, nil
}

// buildAttackPathsFromVulnerabilities creates attack paths from vulnerability data
func (s *AttackPathService) buildAttackPathsFromVulnerabilities(vulns []struct {
	ID          string
	CVEID       sql.NullString
	Severity    string
	Title       string
	Description sql.NullString
	RiskScore   float64
	AgentID     sql.NullString
	Hostname    sql.NullString
	OS          sql.NullString
	Metadata    sql.NullString
}) []AttackPath {
	var paths []AttackPath

	// Group vulnerabilities by agent/host to create paths
	hostVulns := make(map[string][]struct {
		ID          string
		CVEID       sql.NullString
		Severity    string
		Title       string
		Description sql.NullString
		RiskScore   float64
		AgentID     sql.NullString
		Hostname    sql.NullString
		OS          sql.NullString
		Metadata    sql.NullString
	})

	for _, vuln := range vulns {
		hostname := "Unknown"
		if vuln.Hostname.Valid {
			hostname = vuln.Hostname.String
		} else if vuln.AgentID.Valid {
			hostname = vuln.AgentID.String[:8]
		}
		hostVulns[hostname] = append(hostVulns[hostname], vuln)
	}

	// Create attack paths
	pathID := 1
	for hostname, hostVulns := range hostVulns {
		if len(hostVulns) == 0 {
			continue
		}

		// Determine step type based on vulnerability
		steps := []AttackStep{}
		for i, vuln := range hostVulns {
			stepType := "initial_access"
			if i > 0 {
				stepType = "lateral_movement"
			}
			if i == len(hostVulns)-1 {
				stepType = "data_exfiltration"
			}

			cveID := ""
			if vuln.CVEID.Valid {
				cveID = vuln.CVEID.String
			}

			// Calculate likelihood and impact from risk score
			likelihood := vuln.RiskScore / 10.0
			if likelihood > 1.0 {
				likelihood = 1.0
			}
			impact := 0.9
			if vuln.Severity == "critical" {
				impact = 0.9
			} else if vuln.Severity == "high" {
				impact = 0.7
			} else {
				impact = 0.5
			}

			step := AttackStep{
				StepNumber:          i + 1,
				Action:              fmt.Sprintf("Exploit %s", vuln.Title),
				Target:              hostname,
				TargetHostname:      hostname,
				Technique:           s.inferTechnique(vuln.Title, vuln.Description),
				Likelihood:          likelihood,
				Impact:              impact,
				DetectionDifficulty: s.getDetectionDifficulty(vuln.Severity),
				StepType:            stepType,
				CVEID:               cveID,
				VulnerabilityID:     vuln.ID,
				Proof:               s.generateProof(vuln, cveID),
				MitigationControls:  s.getMitigationControls(vuln.Severity),
			}

			// Extract IP from metadata if available
			if vuln.Metadata.Valid {
				var metadata map[string]interface{}
				if err := json.Unmarshal([]byte(vuln.Metadata.String), &metadata); err == nil {
					if ip, ok := metadata["ip_address"].(string); ok {
						step.TargetIP = ip
					}
				}
			}

			steps = append(steps, step)
		}

		if len(steps) > 0 {
			// Calculate path metrics
			totalLikelihood := 1.0
			for _, step := range steps {
				totalLikelihood *= step.Likelihood
			}
			totalImpact := steps[len(steps)-1].Impact // Use last step's impact
			criticalityScore := totalLikelihood * totalImpact

			mitigationPriority := "medium"
			if criticalityScore > 0.5 {
				mitigationPriority = "high"
			} else if criticalityScore < 0.1 {
				mitigationPriority = "low"
			}

			path := AttackPath{
				PathID:             fmt.Sprintf("ap_%03d", pathID),
				Name:               fmt.Sprintf("Attack Path to %s", hostname),
				Steps:              steps,
				TotalLikelihood:    totalLikelihood,
				TotalImpact:        totalImpact,
				CriticalityScore:   criticalityScore,
				MitigationPriority: mitigationPriority,
				DetectionPoints:    []string{"Network logs", "Application logs", "IDS/IPS"},
				PreventionControls: []string{"Patch management", "Network segmentation", "Access controls"},
			}

			paths = append(paths, path)
			pathID++
		}
	}

	return paths
}

// inferTechnique determines MITRE ATT&CK technique from vulnerability
func (s *AttackPathService) inferTechnique(title string, description sql.NullString) string {
	titleLower := title
	descLower := ""
	if description.Valid {
		descLower = description.String
	}

	combined := titleLower + " " + descLower

	if containsAttackPath(combined, "sql injection") || containsAttackPath(combined, "sqli") {
		return "T1059.003 - Command and Scripting Interpreter"
	}
	if containsAttackPath(combined, "xss") || containsAttackPath(combined, "cross-site") {
		return "T1059.007 - JavaScript"
	}
	if containsAttackPath(combined, "rce") || containsAttackPath(combined, "remote code execution") {
		return "T1059 - Command and Scripting Interpreter"
	}
	if containsAttackPath(combined, "privilege") || containsAttackPath(combined, "escalation") {
		return "T1068 - Exploitation for Privilege Escalation"
	}
	if containsAttackPath(combined, "authentication") || containsAttackPath(combined, "bypass") {
		return "T1078 - Valid Accounts"
	}

	return "T1190 - Exploit Public-Facing Application"
}

func containsAttackPath(s, substr string) bool {
	return len(s) >= len(substr) && containsSubstring(strings.ToLower(s), strings.ToLower(substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (s *AttackPathService) getDetectionDifficulty(severity string) string {
	switch severity {
	case "critical":
		return "low" // Critical vulns are easier to detect
	case "high":
		return "medium"
	default:
		return "high"
	}
}

func (s *AttackPathService) generateProof(vuln struct {
	ID          string
	CVEID       sql.NullString
	Severity    string
	Title       string
	Description sql.NullString
	RiskScore   float64
	AgentID     sql.NullString
	Hostname    sql.NullString
	OS          sql.NullString
	Metadata    sql.NullString
}, cveID string) string {
	if cveID != "" {
		return fmt.Sprintf("Exploit available for %s. Check Exploit-DB, Metasploit, or GitHub for proof-of-concept.", cveID)
	}
	return fmt.Sprintf("Vulnerability %s detected. Manual verification recommended.", vuln.Title)
}

func (s *AttackPathService) getMitigationControls(severity string) []string {
	controls := []string{"Apply security patches", "Implement network segmentation"}
	if severity == "critical" {
		controls = append(controls, "Immediate remediation required", "Enable additional monitoring")
	}
	return controls
}

// GetAttackPath retrieves a specific attack path
func (s *AttackPathService) GetAttackPath(pathID string) (*AttackPath, error) {
	paths, err := s.GetAttackPaths("")
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		if path.PathID == pathID {
			return &path, nil
		}
	}

	return nil, fmt.Errorf("attack path not found: %s", pathID)
}

// GenerateAttackPaths generates new attack paths (same as GetAttackPaths for now)
func (s *AttackPathService) GenerateAttackPaths(organizationID string) ([]AttackPath, error) {
	return s.GetAttackPaths(organizationID)
}
