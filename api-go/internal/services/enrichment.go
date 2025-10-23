package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"zerotrace/api/internal/models"
)

// EnrichmentService handles communication with the Python enrichment service
type EnrichmentService struct {
	enrichmentURL string
	httpClient    *http.Client
}

// NewEnrichmentService creates a new enrichment service
func NewEnrichmentService(enrichmentURL string) *EnrichmentService {
	return &EnrichmentService{
		enrichmentURL: enrichmentURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Enrichment can take time
		},
	}
}

// SoftwareItem represents a software item to be enriched
type SoftwareItem struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Vendor  string `json:"vendor,omitempty"`
	Type    string `json:"type,omitempty"`
}

// EnrichedSoftware represents enriched software with CVE data
type EnrichedSoftware struct {
	Name               string    `json:"name"`
	Version            string    `json:"version"`
	CVEs               []CVEData `json:"cves"`
	VulnerabilityCount int       `json:"vulnerability_count"`
	CPEIdentifier      string    `json:"cpe_identifier,omitempty"`
	CPEConfidence      string    `json:"cpe_confidence,omitempty"`
}

// CVEData represents a CVE from the enrichment service
type CVEData struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Severity    string  `json:"severity"`
	CVSSScore   float64 `json:"cvss_score"`
	Published   string  `json:"published_date"`
	Modified    string  `json:"last_modified"`
	Source      string  `json:"source"`
}

// EnrichmentResponse represents the response from the enrichment service
type EnrichmentResponse struct {
	Success bool               `json:"success"`
	Data    []EnrichedSoftware `json:"data"`
	Message string             `json:"message"`
}

// EnrichDependencies enriches dependencies with CVE data from the enrichment service
func (e *EnrichmentService) EnrichDependencies(dependencies []models.Dependency) ([]models.Vulnerability, error) {
	if len(dependencies) == 0 {
		log.Printf("[Enrichment] No dependencies to enrich")
		return []models.Vulnerability{}, nil
	}

	log.Printf("[Enrichment] Starting enrichment for %d dependencies", len(dependencies))

	// Convert dependencies to enrichment request format
	software := make([]SoftwareItem, 0, len(dependencies))
	for _, dep := range dependencies {
		log.Printf("[Enrichment] Processing dependency: %s %s", dep.Name, dep.Version)
		software = append(software, SoftwareItem{
			Name:    dep.Name,
			Version: dep.Version,
			Vendor:  "", // Dependency model doesn't have vendor field
			Type:    dep.Type,
		})
	}

	// Marshal request to JSON
	jsonData, err := json.Marshal(software)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal enrichment request: %w", err)
	}

	// Make HTTP request to enrichment service
	url := fmt.Sprintf("%s/enrich/software", e.enrichmentURL)
	log.Printf("[Enrichment] Sending %d software items to %s", len(software), url)

	resp, err := e.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[Enrichment] Failed to connect to enrichment service: %v", err)
		// Return empty vulnerabilities instead of error to allow API to continue
		return []models.Vulnerability{}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Enrichment] Enrichment service returned status %d: %s", resp.StatusCode, string(body))
		return []models.Vulnerability{}, fmt.Errorf("enrichment service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Enrichment] Failed to read enrichment response: %v", err)
		return []models.Vulnerability{}, nil
	}

	var enrichmentResp EnrichmentResponse
	if err := json.Unmarshal(body, &enrichmentResp); err != nil {
		log.Printf("[Enrichment] Failed to parse enrichment response: %v", err)
		return []models.Vulnerability{}, nil
	}

	if !enrichmentResp.Success {
		log.Printf("[Enrichment] Enrichment service returned error: %s", enrichmentResp.Message)
		return []models.Vulnerability{}, nil
	}

	// Convert enriched data to vulnerabilities
	vulnerabilities := []models.Vulnerability{}

	for _, enriched := range enrichmentResp.Data {
		for _, cve := range enriched.CVEs {
			vuln := models.Vulnerability{
				ID:             cve.ID,
				Type:           "cve",
				Title:          cve.ID,
				Description:    cve.Description,
				Severity:       models.SeverityLevel(cve.Severity),
				CVEID:          cve.ID,
				CVSSScore:      &cve.CVSSScore,
				PackageName:    enriched.Name,
				PackageVersion: enriched.Version,
				Status:         "open",
				Priority:       getPriorityFromCVSS(cve.CVSSScore),
				EnrichmentData: map[string]interface{}{
					"published_date":   cve.Published,
					"last_modified":    cve.Modified,
					"software_name":    enriched.Name,
					"software_version": enriched.Version,
					"source":           cve.Source,
					"cpe_identifier":   enriched.CPEIdentifier,
					"cpe_confidence":   enriched.CPEConfidence,
				},
				CreatedAt: time.Now(),
			}
			vulnerabilities = append(vulnerabilities, vuln)
		}
	}

	log.Printf("[Enrichment] Found %d vulnerabilities across %d software items", len(vulnerabilities), len(enrichmentResp.Data))

	return vulnerabilities, nil
}

// getPriorityFromCVSS converts CVSS score to priority level
func getPriorityFromCVSS(score float64) string {
	switch {
	case score >= 9.0:
		return "critical"
	case score >= 7.0:
		return "high"
	case score >= 4.0:
		return "medium"
	default:
		return "low"
	}
}
