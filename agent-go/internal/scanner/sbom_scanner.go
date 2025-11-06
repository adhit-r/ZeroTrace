package scanner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

// SBOMScanner handles SBOM generation and vulnerability scanning.
type SBOMScanner struct {
	syftPath  string
	grypePath string
}

// NewSBOMScanner creates a new SBOM scanner instance.
func NewSBOMScanner() *SBOMScanner {
	return &SBOMScanner{
		syftPath:  "syft",  // Assumes syft is in PATH
		grypePath: "grype", // Assumes grype is in PATH
	}
}

// Scan performs SBOM generation and vulnerability scanning for a given target.
// The target can be a container image, directory, or other source supported by Syft.
func (s *SBOMScanner) Scan(target string) (*SBOMFinding, error) {
	// 1. Create the syft and grype commands
	syftCmd := exec.Command(s.syftPath, target, "-o", "cyclonedx-json")
	grypeCmd := exec.Command(s.grypePath, "sbom:stdin", "-o", "json")

	// 2. Pipe syft's stdout to grype's stdin
	var err error
	grypeCmd.Stdin, err = syftCmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe for syft: %w", err)
	}

	// 3. Capture grype's output and errors
	var grypeOutput, grypeErr bytes.Buffer
	grypeCmd.Stdout = &grypeOutput
	grypeCmd.Stderr = &grypeErr

	// 4. Start both commands
	if err := syftCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start syft: %w", err)
	}
	if err := grypeCmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start grype: %w", err)
	}

	// 5. Wait for commands to finish
	if err := syftCmd.Wait(); err != nil {
		return nil, fmt.Errorf("syft command failed: %w", err)
	}
	if err := grypeCmd.Wait(); err != nil {
		// Grype exits with a non-zero status code if vulnerabilities are found.
		// This is not a fatal error for the execution itself.
		// We will proceed to parse the output anyway.
	}

	// 6. Check for empty output, which can happen if grype errors out before producing JSON
	if grypeOutput.Len() == 0 {
		return nil, fmt.Errorf("grype produced no output. Stderr: %s", grypeErr.String())
	}

	// 7. Parse Grype JSON output into our structs
	finding, err := s.parseOutput(grypeOutput.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse scanner output: %w", err)
	}
	finding.ServiceName = target

	return finding, nil
}

// parseOutput transforms the output from Grype into our data model.
func (s *SBOMScanner) parseOutput(grypeJSON []byte) (*SBOMFinding, error) {
	finding := &SBOMFinding{
		Components:      []Component{},
		Vulnerabilities: []Vulnerability{},
	}

	// Define a struct that matches the Grype JSON output format
	var grypeResult struct {
		Matches []struct {
			Vulnerability struct {
				ID          string `json:"id"`
				Severity    string `json:"severity"`
				Description string `json:"description"`
				Fix         struct {
					State    string   `json:"state"`
					Versions []string `json:"versions"`
				} `json:"fix"`
			} `json:"vulnerability"`
			Artifact struct {
				Name    string `json:"name"`
				Version string `json:"version"`
				Type    string `json:"type"`
				PURL    string `json:"purl"`
			} `json:"artifact"`
		} `json:"matches"`
	}

	if err := json.Unmarshal(grypeJSON, &grypeResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grype json: %w. Raw output: %s", err, string(grypeJSON))
	}

	components := make(map[string]Component)
	for _, match := range grypeResult.Matches {
		// Add vulnerability
		vuln := Vulnerability{
			CVE:         match.Vulnerability.ID,
			Severity:    match.Vulnerability.Severity,
			Description: match.Vulnerability.Description,
		}
		if match.Vulnerability.Fix.State == "fixed" && len(match.Vulnerability.Fix.Versions) > 0 {
			vuln.FixAvailable = true
			vuln.FixVersion = match.Vulnerability.Fix.Versions[0]
		}
		finding.Vulnerabilities = append(finding.Vulnerabilities, vuln)

		// Add component (if not already seen)
		if _, exists := components[match.Artifact.PURL]; !exists && match.Artifact.PURL != "" {
			comp := Component{
				Name:    match.Artifact.Name,
				Version: match.Artifact.Version,
				Type:    match.Artifact.Type,
				PURL:    match.Artifact.PURL,
			}
			components[match.Artifact.PURL] = comp
		}
	}

	for _, comp := range components {
		finding.Components = append(finding.Components, comp)
	}

	return finding, nil
}
