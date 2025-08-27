package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"zerotrace/agent/internal/models"
)

// CVESource represents a CVE data source
type CVESource interface {
	GetCVE(cveID string) (*models.Vulnerability, error)
	SearchCVEs(query string) ([]models.Vulnerability, error)
	GetRecentCVEs(limit int) ([]models.Vulnerability, error)
}

// NVDSource implements CVE data from NIST National Vulnerability Database
type NVDSource struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewNVDSource creates a new NVD CVE source
func NewNVDSource(apiKey string) *NVDSource {
	return &NVDSource{
		baseURL: "https://services.nvd.nist.gov/rest/json/cves/2.0",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCVE retrieves a specific CVE from NVD
func (n *NVDSource) GetCVE(cveID string) (*models.Vulnerability, error) {
	url := fmt.Sprintf("%s?cveId=%s", n.baseURL, cveID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add API key if available
	if n.apiKey != "" {
		req.Header.Set("apiKey", n.apiKey)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CVE: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var nvdResponse struct {
		Vulnerabilities []struct {
			CVE struct {
				ID          string `json:"id"`
				Description struct {
					DescriptionData []struct {
						Lang  string `json:"lang"`
						Value string `json:"value"`
					} `json:"description_data"`
				} `json:"description"`
				Metrics struct {
					CvssMetricV31 []struct {
						CvssData struct {
							BaseScore    float64 `json:"baseScore"`
							BaseSeverity string  `json:"baseSeverity"`
							VectorString string  `json:"vectorString"`
						} `json:"cvssData"`
					} `json:"cvssMetricV31"`
				} `json:"metrics"`
				Published    string `json:"published"`
				LastModified string `json:"lastModifiedDate"`
			} `json:"cve"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(body, &nvdResponse); err != nil {
		return nil, fmt.Errorf("failed to parse NVD response: %w", err)
	}

	if len(nvdResponse.Vulnerabilities) == 0 {
		return nil, fmt.Errorf("CVE %s not found", cveID)
	}

	nvdCVE := nvdResponse.Vulnerabilities[0].CVE

	// Convert to our vulnerability model
	vuln := &models.Vulnerability{
		ID:          nvdCVE.ID,
		Type:        "cve",
		Severity:    strings.ToLower(nvdCVE.Metrics.CvssMetricV31[0].CvssData.BaseSeverity),
		Title:       nvdCVE.ID,
		Description: nvdCVE.Description.DescriptionData[0].Value,
		CVEID:       nvdCVE.ID,
		CVSSScore:   &nvdCVE.Metrics.CvssMetricV31[0].CvssData.BaseScore,
		CVSSVector:  nvdCVE.Metrics.CvssMetricV31[0].CvssData.VectorString,
		Status:      "open",
		Priority:    getPriorityFromCVSS(nvdCVE.Metrics.CvssMetricV31[0].CvssData.BaseScore),
		CreatedAt:   time.Now(),
	}

	return vuln, nil
}

// SearchCVEs searches for CVEs by keyword
func (n *NVDSource) SearchCVEs(query string) ([]models.Vulnerability, error) {
	url := fmt.Sprintf("%s?keywordSearch=%s", n.baseURL, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if n.apiKey != "" {
		req.Header.Set("apiKey", n.apiKey)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search CVEs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var nvdResponse struct {
		Vulnerabilities []struct {
			CVE struct {
				ID          string `json:"id"`
				Description struct {
					DescriptionData []struct {
						Lang  string `json:"lang"`
						Value string `json:"value"`
					} `json:"description_data"`
				} `json:"description"`
				Metrics struct {
					CvssMetricV31 []struct {
						CvssData struct {
							BaseScore    float64 `json:"baseScore"`
							BaseSeverity string  `json:"baseSeverity"`
							VectorString string  `json:"vectorString"`
						} `json:"cvssData"`
					} `json:"cvssMetricV31"`
				} `json:"metrics"`
			} `json:"cve"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(body, &nvdResponse); err != nil {
		return nil, fmt.Errorf("failed to parse NVD response: %w", err)
	}

	var vulnerabilities []models.Vulnerability
	for _, vuln := range nvdResponse.Vulnerabilities {
		if len(vuln.CVE.Metrics.CvssMetricV31) > 0 {
			baseScore := vuln.CVE.Metrics.CvssMetricV31[0].CvssData.BaseScore
			vulnerabilities = append(vulnerabilities, models.Vulnerability{
				ID:          vuln.CVE.ID,
				Type:        "cve",
				Severity:    strings.ToLower(vuln.CVE.Metrics.CvssMetricV31[0].CvssData.BaseSeverity),
				Title:       vuln.CVE.ID,
				Description: vuln.CVE.Description.DescriptionData[0].Value,
				CVEID:       vuln.CVE.ID,
				CVSSScore:   &baseScore,
				CVSSVector:  vuln.CVE.Metrics.CvssMetricV31[0].CvssData.VectorString,
				Status:      "open",
				Priority:    getPriorityFromCVSS(vuln.CVE.Metrics.CvssMetricV31[0].CvssData.BaseScore),
				CreatedAt:   time.Now(),
			})
		}
	}

	return vulnerabilities, nil
}

// GetRecentCVEs gets recent CVEs from NVD
func (n *NVDSource) GetRecentCVEs(limit int) ([]models.Vulnerability, error) {
	// Get CVEs from the last 7 days
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02T00:00:00:000 UTC-05:00")
	url := fmt.Sprintf("%s?pubStartDate=%s&resultsPerPage=%d", n.baseURL, startDate, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if n.apiKey != "" {
		req.Header.Set("apiKey", n.apiKey)
	}

	resp, err := n.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent CVEs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var nvdResponse struct {
		Vulnerabilities []struct {
			CVE struct {
				ID          string `json:"id"`
				Description struct {
					DescriptionData []struct {
						Lang  string `json:"lang"`
						Value string `json:"value"`
					} `json:"description_data"`
				} `json:"description"`
				Metrics struct {
					CvssMetricV31 []struct {
						CvssData struct {
							BaseScore    float64 `json:"baseScore"`
							BaseSeverity string  `json:"baseSeverity"`
							VectorString string  `json:"vectorString"`
						} `json:"cvssData"`
					} `json:"cvssMetricV31"`
				} `json:"metrics"`
			} `json:"cve"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(body, &nvdResponse); err != nil {
		return nil, fmt.Errorf("failed to parse NVD response: %w", err)
	}

	var vulnerabilities []models.Vulnerability
	for _, vuln := range nvdResponse.Vulnerabilities {
		if len(vuln.CVE.Metrics.CvssMetricV31) > 0 {
			baseScore := vuln.CVE.Metrics.CvssMetricV31[0].CvssData.BaseScore
			vulnerabilities = append(vulnerabilities, models.Vulnerability{
				ID:          vuln.CVE.ID,
				Type:        "cve",
				Severity:    strings.ToLower(vuln.CVE.Metrics.CvssMetricV31[0].CvssData.BaseSeverity),
				Title:       vuln.CVE.ID,
				Description: vuln.CVE.Description.DescriptionData[0].Value,
				CVEID:       vuln.CVE.ID,
				CVSSScore:   &baseScore,
				CVSSVector:  vuln.CVE.Metrics.CvssMetricV31[0].CvssData.VectorString,
				Status:      "open",
				Priority:    getPriorityFromCVSS(vuln.CVE.Metrics.CvssMetricV31[0].CvssData.BaseScore),
				CreatedAt:   time.Now(),
			})
		}
	}

	return vulnerabilities, nil
}

// GitHubAdvisorySource implements CVE data from GitHub Security Advisories
type GitHubAdvisorySource struct {
	baseURL    string
	httpClient *http.Client
}

// NewGitHubAdvisorySource creates a new GitHub advisory source
func NewGitHubAdvisorySource() *GitHubAdvisorySource {
	return &GitHubAdvisorySource{
		baseURL: "https://api.github.com/advisories",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetCVE retrieves a CVE from GitHub advisories
func (g *GitHubAdvisorySource) GetCVE(cveID string) (*models.Vulnerability, error) {
	url := fmt.Sprintf("%s?cve_id=%s", g.baseURL, cveID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch CVE: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var advisories []struct {
		GHSAID          string `json:"ghsa_id"`
		CVEID           string `json:"cve_id"`
		Summary         string `json:"summary"`
		Description     string `json:"description"`
		Severity        string `json:"severity"`
		PublishedAt     string `json:"published_at"`
		UpdatedAt       string `json:"updated_at"`
		Vulnerabilities []struct {
			Package struct {
				Name string `json:"name"`
			} `json:"package"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(body, &advisories); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	if len(advisories) == 0 {
		return nil, fmt.Errorf("CVE %s not found", cveID)
	}

	advisory := advisories[0]

	vuln := &models.Vulnerability{
		ID:          advisory.CVEID,
		Type:        "cve",
		Severity:    strings.ToLower(advisory.Severity),
		Title:       advisory.Summary,
		Description: advisory.Description,
		CVEID:       advisory.CVEID,
		Status:      "open",
		Priority:    getPriorityFromSeverity(advisory.Severity),
		CreatedAt:   time.Now(),
	}

	if len(advisory.Vulnerabilities) > 0 {
		vuln.PackageName = advisory.Vulnerabilities[0].Package.Name
	}

	return vuln, nil
}

// SearchCVEs searches for CVEs in GitHub advisories
func (g *GitHubAdvisorySource) SearchCVEs(query string) ([]models.Vulnerability, error) {
	// GitHub doesn't have a direct search endpoint, so we'll use a different approach
	// For now, return empty slice - implement based on specific needs
	return []models.Vulnerability{}, nil
}

// GetRecentCVEs gets recent CVEs from GitHub advisories
func (g *GitHubAdvisorySource) GetRecentCVEs(limit int) ([]models.Vulnerability, error) {
	url := fmt.Sprintf("%s?per_page=%d&sort=published", g.baseURL, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent CVEs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var advisories []struct {
		GHSAID          string `json:"ghsa_id"`
		CVEID           string `json:"cve_id"`
		Summary         string `json:"summary"`
		Description     string `json:"description"`
		Severity        string `json:"severity"`
		PublishedAt     string `json:"published_at"`
		Vulnerabilities []struct {
			Package struct {
				Name string `json:"name"`
			} `json:"package"`
		} `json:"vulnerabilities"`
	}

	if err := json.Unmarshal(body, &advisories); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	var vulnerabilities []models.Vulnerability
	for _, advisory := range advisories {
		if advisory.CVEID != "" {
			vuln := models.Vulnerability{
				ID:          advisory.CVEID,
				Type:        "cve",
				Severity:    strings.ToLower(advisory.Severity),
				Title:       advisory.Summary,
				Description: advisory.Description,
				CVEID:       advisory.CVEID,
				Status:      "open",
				Priority:    getPriorityFromSeverity(advisory.Severity),
				CreatedAt:   time.Now(),
			}

			if len(advisory.Vulnerabilities) > 0 {
				vuln.PackageName = advisory.Vulnerabilities[0].Package.Name
			}

			vulnerabilities = append(vulnerabilities, vuln)
		}
	}

	return vulnerabilities, nil
}

// Helper functions
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

func getPriorityFromSeverity(severity string) string {
	switch strings.ToLower(severity) {
	case "critical":
		return "critical"
	case "high":
		return "high"
	case "medium":
		return "medium"
	case "low":
		return "low"
	default:
		return "medium"
	}
}
