package processor

import (
	"time"

	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/models"
)

// Processor handles scan result processing
type Processor struct {
	config *config.Config
}

// NewProcessor creates a new processor instance
func NewProcessor(cfg *config.Config) *Processor {
	return &Processor{
		config: cfg,
	}
}

// Process processes scan results
func (p *Processor) Process(result *models.ScanResult) (*models.ScanResult, error) {
	// Add processing metadata
	result.Metadata["processed_at"] = time.Now()
	result.Metadata["processor_version"] = "1.0.0"
	result.Metadata["agent_id"] = p.config.AgentID

	// Agent sends raw dependencies to API
	// API handles enrichment asynchronously via Python enrichment service
	// No local enrichment processing needed

	// Process vulnerabilities (if any from local scanning)
	for i := range result.Vulnerabilities {
		p.processVulnerability(&result.Vulnerabilities[i])
	}

	// Process dependencies (add metadata only)
	for i := range result.Dependencies {
		p.processDependency(&result.Dependencies[i])
	}

	return result, nil
}

// processVulnerability processes a single vulnerability
func (p *Processor) processVulnerability(vuln *models.Vulnerability) {
	// Add processing metadata
	if vuln.EnrichmentData == nil {
		vuln.EnrichmentData = make(map[string]any)
	}
	vuln.EnrichmentData["processed"] = true
	vuln.EnrichmentData["processor_timestamp"] = time.Now()
	vuln.EnrichmentData["agent_id"] = p.config.AgentID

	// Vulnerability is already enriched by enrichment service
	// Additional processing can be added here
}

// processDependency processes a single dependency
func (p *Processor) processDependency(dep *models.Dependency) {
	// Add processing metadata
	if dep.Metadata == nil {
		dep.Metadata = make(map[string]any)
	}
	dep.Metadata["processed"] = true
	dep.Metadata["processor_timestamp"] = time.Now()

	// Dependencies are enriched by enrichment service
	// Additional processing can be added here
}
