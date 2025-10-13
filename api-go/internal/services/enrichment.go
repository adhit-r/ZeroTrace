package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"zerotrace/api/internal/models"
)

// EnrichmentService handles async enrichment of applications
type EnrichmentService struct {
	enrichmentURL string
	httpClient    *http.Client
}

// NewEnrichmentService creates a new enrichment service
func NewEnrichmentService(enrichmentURL string) *EnrichmentService {
	return &EnrichmentService{
		enrichmentURL: enrichmentURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for API-level enrichment
		},
	}
}

// EnrichApplicationsAsync enriches applications asynchronously
func (es *EnrichmentService) EnrichApplicationsAsync(agentID string, dependencies []models.Dependency) {
	go func() {
		log.Printf("[EnrichmentService] Starting async enrichment for agent %s with %d applications", agentID, len(dependencies))

		// Convert dependencies to enrichment format
		softwareItems := make([]map[string]interface{}, len(dependencies))
		for i, dep := range dependencies {
			softwareItems[i] = map[string]interface{}{
				"name":    dep.Name,
				"version": dep.Version,
				"type":    dep.Type,
			}
		}

		// Prepare enrichment request
		request := map[string]interface{}{
			"software": softwareItems,
		}

		// Send to enrichment service
		jsonData, err := json.Marshal(request)
		if err != nil {
			log.Printf("[EnrichmentService] Failed to marshal enrichment request: %v", err)
			return
		}

		url := fmt.Sprintf("%s/enrich/software", es.enrichmentURL)
		resp, err := es.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("[EnrichmentService] Failed to call enrichment service: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("[EnrichmentService] Enrichment service returned status %d", resp.StatusCode)
			return
		}

		// Parse response
		var enrichmentResponse struct {
			Vulnerabilities []models.Vulnerability `json:"vulnerabilities"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&enrichmentResponse); err != nil {
			log.Printf("[EnrichmentService] Failed to decode enrichment response: %v", err)
			return
		}

		log.Printf("[EnrichmentService] Enrichment completed for agent %s: found %d vulnerabilities",
			agentID, len(enrichmentResponse.Vulnerabilities))

		// TODO: Update agent metadata with enriched vulnerabilities
		// This would require access to the agent service
	}()
}
