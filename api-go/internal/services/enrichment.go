package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EnrichmentService handles communication with the Python enrichment server
type EnrichmentService struct {
	baseURL    string
	httpClient *http.Client
}

// NewEnrichmentService creates a new enrichment service
func NewEnrichmentService(baseURL string) *EnrichmentService {
	return &EnrichmentService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SoftwareItem represents a software item to be enriched
type SoftwareItem struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	PackageType string `json:"package_type,omitempty"`
	Path        string `json:"path,omitempty"`
}

// EnrichedSoftware represents enriched software with CVE data
type EnrichedSoftware struct {
	SoftwareItem
	CVEs               []CVE  `json:"cves"`
	VulnerabilityCount int    `json:"vulnerability_count"`
	EnrichedAt         string `json:"enriched_at"`
	Error              string `json:"error,omitempty"`
}

// CVE represents a Common Vulnerability and Exposure
type CVE struct {
	ID            string  `json:"id"`
	Description   string  `json:"description"`
	Severity      string  `json:"severity"`
	CVSSScore     float64 `json:"cvss_score"`
	PublishedDate string  `json:"published_date"`
	LastModified  string  `json:"last_modified"`
	Source        string  `json:"source"`
}

// EnrichmentResponse represents the response from the enrichment service
type EnrichmentResponse struct {
	Success bool               `json:"success"`
	Data    []EnrichedSoftware `json:"data"`
	Message string             `json:"message"`
}

// EnrichSoftware enriches a list of software items with CVE data
func (es *EnrichmentService) EnrichSoftware(software []SoftwareItem) ([]EnrichedSoftware, error) {
	// Prepare request
	requestBody, err := json.Marshal(software)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal software data: %w", err)
	}

	// Make request to enrichment service
	url := fmt.Sprintf("%s/enrich/software", es.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := es.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to enrichment service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("enrichment service returned status %d", resp.StatusCode)
	}

	// Parse response
	var enrichmentResp EnrichmentResponse
	if err := json.NewDecoder(resp.Body).Decode(&enrichmentResp); err != nil {
		return nil, fmt.Errorf("failed to decode enrichment response: %w", err)
	}

	if !enrichmentResp.Success {
		return nil, fmt.Errorf("enrichment service error: %s", enrichmentResp.Message)
	}

	return enrichmentResp.Data, nil
}

// EnrichSoftwareBatch starts a background enrichment job
func (es *EnrichmentService) EnrichSoftwareBatch(software []SoftwareItem) (string, error) {
	// Prepare request
	requestBody, err := json.Marshal(software)
	if err != nil {
		return "", fmt.Errorf("failed to marshal software data: %w", err)
	}

	// Make request to enrichment service
	url := fmt.Sprintf("%s/enrich/batch", es.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := es.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to enrichment service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("enrichment service returned status %d", resp.StatusCode)
	}

	// Parse response
	var batchResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		JobID   string `json:"job_id,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&batchResp); err != nil {
		return "", fmt.Errorf("failed to decode batch response: %w", err)
	}

	if !batchResp.Success {
		return "", fmt.Errorf("enrichment service error: %s", batchResp.Message)
	}

	return batchResp.JobID, nil
}

// GetEnrichmentStatus gets the status of an enrichment job
func (es *EnrichmentService) GetEnrichmentStatus(jobID string) (map[string]interface{}, error) {
	// Make request to enrichment service
	url := fmt.Sprintf("%s/enrich/status/%s", es.baseURL, jobID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := es.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to enrichment service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("enrichment service returned status %d", resp.StatusCode)
	}

	// Parse response
	var status map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %w", err)
	}

	return status, nil
}

// HealthCheck checks if the enrichment service is healthy
func (es *EnrichmentService) HealthCheck() error {
	// Make request to enrichment service
	url := fmt.Sprintf("%s/health", es.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// Send request
	resp, err := es.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send health check request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("enrichment service health check returned status %d", resp.StatusCode)
	}

	return nil
}

