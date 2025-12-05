package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// AIService handles communication with the Python AI service
type AIService struct {
	aiServiceURL string
	httpClient   *http.Client
}

// NewAIService creates a new AI service client
func NewAIService(aiServiceURL string) *AIService {
	return &AIService{
		aiServiceURL: aiServiceURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // AI analysis can take longer
		},
	}
}

// VulnerabilityData represents vulnerability data for AI analysis
type VulnerabilityData struct {
	ID          string                 `json:"id"`
	CVEID       string                 `json:"cve_id,omitempty"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	CVSSScore   *float64               `json:"cvss_score,omitempty"`
	PackageName string                 `json:"package_name,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// OrganizationContext represents organization context for AI analysis
type OrganizationContext struct {
	OrganizationID string                 `json:"organization_id"`
	Industry       string                 `json:"industry,omitempty"`
	RiskTolerance  string                 `json:"risk_tolerance,omitempty"`
	TechStack      map[string]interface{} `json:"tech_stack,omitempty"`
}

// ComprehensiveAnalysisRequest represents request for comprehensive analysis
type ComprehensiveAnalysisRequest struct {
	VulnerabilityData  VulnerabilityData  `json:"vulnerability_data"`
	OrganizationContext *OrganizationContext `json:"organization_context,omitempty"`
}

// ComprehensiveAnalysisResponse represents comprehensive analysis response
type ComprehensiveAnalysisResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}

// AnalyzeVulnerabilityComprehensive performs comprehensive AI analysis
func (a *AIService) AnalyzeVulnerabilityComprehensive(vulnData VulnerabilityData, orgContext *OrganizationContext) (map[string]interface{}, error) {
	req := ComprehensiveAnalysisRequest{
		VulnerabilityData:  vulnData,
		OrganizationContext: orgContext,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal AI analysis request: %w", err)
	}

	url := fmt.Sprintf("%s/ai/analyze/comprehensive", a.aiServiceURL)
	log.Printf("[AI Service] Sending comprehensive analysis request to %s", url)

	resp, err := a.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[AI Service] Failed to connect to AI service: %v", err)
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[AI Service] AI service returned status %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read AI service response: %w", err)
	}

	var aiResp ComprehensiveAnalysisResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI service response: %w", err)
	}

	if !aiResp.Success {
		return nil, fmt.Errorf("AI service returned error: %s", aiResp.Message)
	}

	return aiResp.Data, nil
}

// GetExploitIntelligence retrieves exploit intelligence for a CVE
func (a *AIService) GetExploitIntelligence(cveID string, packageName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/ai/exploit-intelligence/%s", a.aiServiceURL, cveID)
	if packageName != "" {
		url += fmt.Sprintf("?package_name=%s", packageName)
	}

	log.Printf("[AI Service] Requesting exploit intelligence for %s from %s", cveID, url)

	resp, err := a.httpClient.Get(url)
	if err != nil {
		log.Printf("[AI Service] Failed to connect to AI service: %v", err)
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read AI service response: %w", err)
	}

	var aiResp ComprehensiveAnalysisResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI service response: %w", err)
	}

	if !aiResp.Success {
		return nil, fmt.Errorf("AI service returned error: %s", aiResp.Message)
	}

	return aiResp.Data, nil
}

// GetPredictiveAnalysis performs predictive analysis
func (a *AIService) GetPredictiveAnalysis(vulnData VulnerabilityData, orgContext *OrganizationContext) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"vulnerability_data":  vulnData,
		"organization_context": orgContext,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal predictive analysis request: %w", err)
	}

	url := fmt.Sprintf("%s/ai/analyze/predictive", a.aiServiceURL)
	resp, err := a.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read AI service response: %w", err)
	}

	var aiResp ComprehensiveAnalysisResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI service response: %w", err)
	}

	return aiResp.Data, nil
}

// GetRemediationPlan generates remediation plan
func (a *AIService) GetRemediationPlan(vulnData VulnerabilityData, orgContext *OrganizationContext) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"vulnerability_data":  vulnData,
		"organization_context": orgContext,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal remediation plan request: %w", err)
	}

	url := fmt.Sprintf("%s/ai/remediation-plan", a.aiServiceURL)
	resp, err := a.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read AI service response: %w", err)
	}

	var aiResp ComprehensiveAnalysisResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI service response: %w", err)
	}

	return aiResp.Data, nil
}

// AnalyzeVulnerabilityTrends analyzes trends
func (a *AIService) AnalyzeVulnerabilityTrends(vulnerabilities []VulnerabilityData, orgContext *OrganizationContext) (map[string]interface{}, error) {
	req := map[string]interface{}{
		"vulnerabilities":      vulnerabilities,
		"organization_context": orgContext,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trend analysis request: %w", err)
	}

	url := fmt.Sprintf("%s/ai/analyze/trends", a.aiServiceURL)
	resp, err := a.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read AI service response: %w", err)
	}

	var aiResp ComprehensiveAnalysisResponse
	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse AI service response: %w", err)
	}

	return aiResp.Data, nil
}

