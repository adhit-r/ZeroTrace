package scanner

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// APIScanner handles API security scanning
type APIScanner struct {
	config *config.Config
}

// APIFinding represents an API security finding
type APIFinding struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`     // auth, authorization, rate_limit, cors, data_exposure
	Severity      string                 `json:"severity"` // critical, high, medium, low
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Endpoint      string                 `json:"endpoint"`
	Method        string                 `json:"method"`
	StatusCode    int                    `json:"status_code,omitempty"`
	CurrentValue  string                 `json:"current_value,omitempty"`
	RequiredValue string                 `json:"required_value,omitempty"`
	Remediation   string                 `json:"remediation"`
	DiscoveredAt  time.Time              `json:"discovered_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// APIEndpoint represents a discovered API endpoint
type APIEndpoint struct {
	URL             string            `json:"url"`
	Method          string            `json:"method"`
	Headers         map[string]string `json:"headers"`
	ResponseTime    int64             `json:"response_time_ms"`
	StatusCode      int               `json:"status_code"`
	ContentType     string            `json:"content_type"`
	ContentLength   int64             `json:"content_length"`
	IsAuthenticated bool              `json:"is_authenticated"`
	IsRateLimited   bool              `json:"is_rate_limited"`
	HasCORS         bool              `json:"has_cors"`
	APIVersion      string            `json:"api_version,omitempty"`
	Framework       string            `json:"framework,omitempty"`
}

// OWASPTop10 represents OWASP API Top 10 security issues
type OWASPTop10 struct {
	BrokenObjectLevelAuthorization   bool `json:"broken_object_level_authorization"`
	BrokenUserAuthentication         bool `json:"broken_user_authentication"`
	ExcessiveDataExposure            bool `json:"excessive_data_exposure"`
	LackOfResourcesAndRateLimiting   bool `json:"lack_of_resources_and_rate_limiting"`
	BrokenFunctionLevelAuthorization bool `json:"broken_function_level_authorization"`
	MassAssignment                   bool `json:"mass_assignment"`
	SecurityMisconfiguration         bool `json:"security_misconfiguration"`
	Injection                        bool `json:"injection"`
	ImproperAssetsManagement         bool `json:"improper_assets_management"`
	InsufficientLoggingAndMonitoring bool `json:"insufficient_logging_and_monitoring"`
}

// NewAPIScanner creates a new API security scanner
func NewAPIScanner(cfg *config.Config) *APIScanner {
	return &APIScanner{
		config: cfg,
	}
}

// Scan performs comprehensive API security scanning
func (as *APIScanner) Scan() ([]APIFinding, []APIEndpoint, OWASPTop10, error) {
	var findings []APIFinding
	var endpoints []APIEndpoint
	var owaspTop10 OWASPTop10

	// Discover API endpoints
	discoveredEndpoints := as.discoverAPIEndpoints()
	endpoints = append(endpoints, discoveredEndpoints...)

	// Scan each discovered endpoint
	for _, endpoint := range discoveredEndpoints {
		endpointFindings := as.scanAPIEndpoint(endpoint)
		findings = append(findings, endpointFindings...)
	}

	// Check for OWASP API Top 10 issues
	owaspFindings := as.checkOWASPTop10(endpoints)
	findings = append(findings, owaspFindings...)

	// Update OWASP Top 10 status
	owaspTop10 = as.analyzeOWASPTop10(findings)

	return findings, endpoints, owaspTop10, nil
}

// discoverAPIEndpoints discovers API endpoints
func (as *APIScanner) discoverAPIEndpoints() []APIEndpoint {
	var endpoints []APIEndpoint

	// Common API endpoints to test
	commonEndpoints := []string{
		"http://localhost:3000/api/",
		"http://localhost:8080/api/",
		"http://localhost:5000/api/",
		"http://localhost:3001/api/",
		"http://localhost:5173/api/",
	}

	// Common API paths
	commonPaths := []string{
		"",
		"vulnerabilities",
		"agents",
		"health",
		"status",
		"users",
		"auth",
		"login",
		"logout",
		"register",
		"profile",
		"settings",
	}

	// Common HTTP methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}

	for _, baseURL := range commonEndpoints {
		for _, path := range commonPaths {
			for _, method := range methods {
				url := baseURL + path
				endpoint := as.testAPIEndpoint(url, method)
				if endpoint != nil {
					endpoints = append(endpoints, *endpoint)
				}
			}
		}
	}

	return endpoints
}

// testAPIEndpoint tests a specific API endpoint
func (as *APIScanner) testAPIEndpoint(url, method string) *APIEndpoint {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil
	}

	// Add common headers
	req.Header.Set("User-Agent", "ZeroTrace-API-Scanner/1.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	responseTime := time.Since(start).Milliseconds()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	endpoint := &APIEndpoint{
		URL:           url,
		Method:        method,
		Headers:       make(map[string]string),
		ResponseTime:  responseTime,
		StatusCode:    resp.StatusCode,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: int64(len(body)),
	}

	// Extract headers
	for name, values := range resp.Header {
		if len(values) > 0 {
			endpoint.Headers[name] = values[0]
		}
	}

	// Check for authentication
	endpoint.IsAuthenticated = as.checkAuthentication(resp)

	// Check for rate limiting
	endpoint.IsRateLimited = as.checkRateLimiting(resp)

	// Check for CORS
	endpoint.HasCORS = as.checkCORS(resp)

	// Detect API version
	endpoint.APIVersion = as.detectAPIVersion(resp)

	// Detect framework
	endpoint.Framework = as.detectFramework(resp)

	return endpoint
}

// checkAuthentication checks if endpoint requires authentication
func (as *APIScanner) checkAuthentication(resp *http.Response) bool {
	// Check for 401 Unauthorized
	if resp.StatusCode == 401 {
		return true
	}

	// Check for WWW-Authenticate header
	if resp.Header.Get("WWW-Authenticate") != "" {
		return true
	}

	// Check for Authorization header in request
	// This would need to be passed from the request
	return false
}

// checkRateLimiting checks if endpoint has rate limiting
func (as *APIScanner) checkRateLimiting(resp *http.Response) bool {
	// Check for rate limiting headers
	rateLimitHeaders := []string{
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
		"Retry-After",
	}

	for _, header := range rateLimitHeaders {
		if resp.Header.Get(header) != "" {
			return true
		}
	}

	return false
}

// checkCORS checks if endpoint has CORS configuration
func (as *APIScanner) checkCORS(resp *http.Response) bool {
	// Check for CORS headers
	corsHeaders := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Credentials",
	}

	for _, header := range corsHeaders {
		if resp.Header.Get(header) != "" {
			return true
		}
	}

	return false
}

// detectAPIVersion detects API version
func (as *APIScanner) detectAPIVersion(resp *http.Response) string {
	// Check for API version in headers
	version := resp.Header.Get("API-Version")
	if version != "" {
		return version
	}

	// Check for version in response body
	// This would require parsing the response body
	return ""
}

// detectFramework detects the API framework
func (as *APIScanner) detectFramework(resp *http.Response) string {
	// Check for framework-specific headers
	server := resp.Header.Get("Server")
	if server != "" {
		server = strings.ToLower(server)
		if strings.Contains(server, "express") {
			return "Express.js"
		}
		if strings.Contains(server, "nginx") {
			return "Nginx"
		}
		if strings.Contains(server, "apache") {
			return "Apache"
		}
	}

	// Check for framework-specific response patterns
	// This would require analyzing the response body
	return ""
}

// scanAPIEndpoint scans a specific API endpoint for security issues
func (as *APIScanner) scanAPIEndpoint(endpoint APIEndpoint) []APIFinding {
	var findings []APIFinding

	// Check for missing authentication
	if !endpoint.IsAuthenticated && endpoint.StatusCode == 200 {
		finding := APIFinding{
			ID:           uuid.New().String(),
			Type:         "auth",
			Severity:     "high",
			Title:        "Missing Authentication",
			Description:  fmt.Sprintf("API endpoint %s %s is accessible without authentication", endpoint.Method, endpoint.URL),
			Endpoint:     endpoint.URL,
			Method:       endpoint.Method,
			StatusCode:   endpoint.StatusCode,
			Remediation:  "Implement authentication for this API endpoint",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"endpoint":      endpoint.URL,
				"method":        endpoint.Method,
				"auth_required": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for missing rate limiting
	if !endpoint.IsRateLimited {
		finding := APIFinding{
			ID:           uuid.New().String(),
			Type:         "rate_limit",
			Severity:     "medium",
			Title:        "Missing Rate Limiting",
			Description:  fmt.Sprintf("API endpoint %s %s has no rate limiting", endpoint.Method, endpoint.URL),
			Endpoint:     endpoint.URL,
			Method:       endpoint.Method,
			StatusCode:   endpoint.StatusCode,
			Remediation:  "Implement rate limiting for this API endpoint",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"endpoint":     endpoint.URL,
				"method":       endpoint.Method,
				"rate_limited": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for CORS misconfiguration
	if endpoint.HasCORS {
		origin := endpoint.Headers["Access-Control-Allow-Origin"]
		if origin == "*" {
			finding := APIFinding{
				ID:            uuid.New().String(),
				Type:          "cors",
				Severity:      "medium",
				Title:         "CORS Misconfiguration",
				Description:   fmt.Sprintf("API endpoint %s %s allows all origins (*)", endpoint.Method, endpoint.URL),
				Endpoint:      endpoint.URL,
				Method:        endpoint.Method,
				StatusCode:    endpoint.StatusCode,
				CurrentValue:  "*",
				RequiredValue: "Specific domain",
				Remediation:   "Restrict CORS to specific domains instead of allowing all origins",
				DiscoveredAt:  time.Now(),
				Metadata: map[string]interface{}{
					"endpoint":    endpoint.URL,
					"method":      endpoint.Method,
					"cors_origin": origin,
				},
			}
			findings = append(findings, finding)
		}
	}

	// Check for information disclosure
	if endpoint.StatusCode == 200 && endpoint.ContentLength > 1000 {
		finding := APIFinding{
			ID:           uuid.New().String(),
			Type:         "data_exposure",
			Severity:     "low",
			Title:        "Potential Information Disclosure",
			Description:  fmt.Sprintf("API endpoint %s %s returns large response (%d bytes)", endpoint.Method, endpoint.URL, endpoint.ContentLength),
			Endpoint:     endpoint.URL,
			Method:       endpoint.Method,
			StatusCode:   endpoint.StatusCode,
			Remediation:  "Review response size and implement pagination if needed",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"endpoint":       endpoint.URL,
				"method":         endpoint.Method,
				"content_length": endpoint.ContentLength,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// checkOWASPTop10 checks for OWASP API Top 10 issues
func (as *APIScanner) checkOWASPTop10(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// A01:2021 – Broken Object Level Authorization
	findings = append(findings, as.checkBrokenObjectLevelAuthorization(endpoints)...)

	// A02:2021 – Broken User Authentication
	findings = append(findings, as.checkBrokenUserAuthentication(endpoints)...)

	// A03:2021 – Excessive Data Exposure
	findings = append(findings, as.checkExcessiveDataExposure(endpoints)...)

	// A04:2021 – Lack of Resources & Rate Limiting
	findings = append(findings, as.checkLackOfResourcesAndRateLimiting(endpoints)...)

	// A05:2021 – Broken Function Level Authorization
	findings = append(findings, as.checkBrokenFunctionLevelAuthorization(endpoints)...)

	// A06:2021 – Mass Assignment
	findings = append(findings, as.checkMassAssignment(endpoints)...)

	// A07:2021 – Security Misconfiguration
	findings = append(findings, as.checkSecurityMisconfiguration(endpoints)...)

	// A08:2021 – Injection
	findings = append(findings, as.checkInjection(endpoints)...)

	// A09:2021 – Improper Assets Management
	findings = append(findings, as.checkImproperAssetsManagement(endpoints)...)

	// A10:2021 – Insufficient Logging & Monitoring
	findings = append(findings, as.checkInsufficientLoggingAndMonitoring(endpoints)...)

	return findings
}

// checkBrokenObjectLevelAuthorization checks for A01:2021
func (as *APIScanner) checkBrokenObjectLevelAuthorization(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// This would require testing for IDOR vulnerabilities
	// For now, return a placeholder finding
	finding := APIFinding{
		ID:           uuid.New().String(),
		Type:         "authorization",
		Severity:     "high",
		Title:        "Potential Broken Object Level Authorization",
		Description:  "API endpoints may be vulnerable to IDOR (Insecure Direct Object Reference) attacks",
		Remediation:  "Implement proper object-level authorization checks",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"owasp_top10": "A01:2021",
			"category":    "authorization",
		},
	}
	findings = append(findings, finding)

	return findings
}

// checkBrokenUserAuthentication checks for A02:2021
func (as *APIScanner) checkBrokenUserAuthentication(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// Check for endpoints without authentication
	for _, endpoint := range endpoints {
		if !endpoint.IsAuthenticated && endpoint.StatusCode == 200 {
			finding := APIFinding{
				ID:           uuid.New().String(),
				Type:         "auth",
				Severity:     "high",
				Title:        "Broken User Authentication",
				Description:  fmt.Sprintf("API endpoint %s %s lacks proper authentication", endpoint.Method, endpoint.URL),
				Endpoint:     endpoint.URL,
				Method:       endpoint.Method,
				Remediation:  "Implement proper user authentication",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"owasp_top10": "A02:2021",
					"endpoint":    endpoint.URL,
					"method":      endpoint.Method,
				},
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

// checkExcessiveDataExposure checks for A03:2021
func (as *APIScanner) checkExcessiveDataExposure(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// Check for endpoints returning large amounts of data
	for _, endpoint := range endpoints {
		if endpoint.ContentLength > 10000 { // 10KB threshold
			finding := APIFinding{
				ID:           uuid.New().String(),
				Type:         "data_exposure",
				Severity:     "medium",
				Title:        "Excessive Data Exposure",
				Description:  fmt.Sprintf("API endpoint %s %s returns large response (%d bytes)", endpoint.Method, endpoint.URL, endpoint.ContentLength),
				Endpoint:     endpoint.URL,
				Method:       endpoint.Method,
				Remediation:  "Implement data filtering and pagination",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"owasp_top10":    "A03:2021",
					"endpoint":       endpoint.URL,
					"method":         endpoint.Method,
					"content_length": endpoint.ContentLength,
				},
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

// checkLackOfResourcesAndRateLimiting checks for A04:2021
func (as *APIScanner) checkLackOfResourcesAndRateLimiting(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// Check for endpoints without rate limiting
	for _, endpoint := range endpoints {
		if !endpoint.IsRateLimited {
			finding := APIFinding{
				ID:           uuid.New().String(),
				Type:         "rate_limit",
				Severity:     "medium",
				Title:        "Lack of Resources & Rate Limiting",
				Description:  fmt.Sprintf("API endpoint %s %s has no rate limiting", endpoint.Method, endpoint.URL),
				Endpoint:     endpoint.URL,
				Method:       endpoint.Method,
				Remediation:  "Implement rate limiting and resource management",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"owasp_top10": "A04:2021",
					"endpoint":    endpoint.URL,
					"method":      endpoint.Method,
				},
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

// checkBrokenFunctionLevelAuthorization checks for A05:2021
func (as *APIScanner) checkBrokenFunctionLevelAuthorization(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// This would require testing for function-level authorization
	// For now, return a placeholder finding
	finding := APIFinding{
		ID:           uuid.New().String(),
		Type:         "authorization",
		Severity:     "high",
		Title:        "Potential Broken Function Level Authorization",
		Description:  "API endpoints may lack proper function-level authorization",
		Remediation:  "Implement proper function-level authorization checks",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"owasp_top10": "A05:2021",
			"category":    "authorization",
		},
	}
	findings = append(findings, finding)

	return findings
}

// checkMassAssignment checks for A06:2021
func (as *APIScanner) checkMassAssignment(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// This would require testing for mass assignment vulnerabilities
	// For now, return a placeholder finding
	finding := APIFinding{
		ID:           uuid.New().String(),
		Type:         "data_exposure",
		Severity:     "medium",
		Title:        "Potential Mass Assignment Vulnerability",
		Description:  "API endpoints may be vulnerable to mass assignment attacks",
		Remediation:  "Implement proper input validation and filtering",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"owasp_top10": "A06:2021",
			"category":    "input_validation",
		},
	}
	findings = append(findings, finding)

	return findings
}

// checkSecurityMisconfiguration checks for A07:2021
func (as *APIScanner) checkSecurityMisconfiguration(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// Check for security headers
	for _, endpoint := range endpoints {
		securityHeaders := []string{
			"X-Content-Type-Options",
			"X-Frame-Options",
			"X-XSS-Protection",
			"Strict-Transport-Security",
		}

		missingHeaders := []string{}
		for _, header := range securityHeaders {
			if endpoint.Headers[header] == "" {
				missingHeaders = append(missingHeaders, header)
			}
		}

		if len(missingHeaders) > 0 {
			finding := APIFinding{
				ID:           uuid.New().String(),
				Type:         "config",
				Severity:     "medium",
				Title:        "Security Misconfiguration",
				Description:  fmt.Sprintf("API endpoint %s %s missing security headers: %s", endpoint.Method, endpoint.URL, strings.Join(missingHeaders, ", ")),
				Endpoint:     endpoint.URL,
				Method:       endpoint.Method,
				Remediation:  "Add missing security headers",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"owasp_top10":     "A07:2021",
					"endpoint":        endpoint.URL,
					"method":          endpoint.Method,
					"missing_headers": missingHeaders,
				},
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

// checkInjection checks for A08:2021
func (as *APIScanner) checkInjection(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// This would require testing for injection vulnerabilities
	// For now, return a placeholder finding
	finding := APIFinding{
		ID:           uuid.New().String(),
		Type:         "injection",
		Severity:     "high",
		Title:        "Potential Injection Vulnerability",
		Description:  "API endpoints may be vulnerable to injection attacks",
		Remediation:  "Implement proper input validation and parameterized queries",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"owasp_top10": "A08:2021",
			"category":    "injection",
		},
	}
	findings = append(findings, finding)

	return findings
}

// checkImproperAssetsManagement checks for A09:2021
func (as *APIScanner) checkImproperAssetsManagement(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// Check for API versioning
	for _, endpoint := range endpoints {
		if endpoint.APIVersion == "" {
			finding := APIFinding{
				ID:           uuid.New().String(),
				Type:         "config",
				Severity:     "low",
				Title:        "Improper Assets Management",
				Description:  fmt.Sprintf("API endpoint %s %s has no version information", endpoint.Method, endpoint.URL),
				Endpoint:     endpoint.URL,
				Method:       endpoint.Method,
				Remediation:  "Implement proper API versioning",
				DiscoveredAt: time.Now(),
				Metadata: map[string]interface{}{
					"owasp_top10": "A09:2021",
					"endpoint":    endpoint.URL,
					"method":      endpoint.Method,
				},
			}
			findings = append(findings, finding)
		}
	}

	return findings
}

// checkInsufficientLoggingAndMonitoring checks for A10:2021
func (as *APIScanner) checkInsufficientLoggingAndMonitoring(endpoints []APIEndpoint) []APIFinding {
	var findings []APIFinding

	// This would require checking for logging and monitoring
	// For now, return a placeholder finding
	finding := APIFinding{
		ID:           uuid.New().String(),
		Type:         "monitoring",
		Severity:     "medium",
		Title:        "Insufficient Logging & Monitoring",
		Description:  "API endpoints may lack proper logging and monitoring",
		Remediation:  "Implement comprehensive logging and monitoring",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"owasp_top10": "A10:2021",
			"category":    "monitoring",
		},
	}
	findings = append(findings, finding)

	return findings
}

// analyzeOWASPTop10 analyzes findings to determine OWASP Top 10 status
func (as *APIScanner) analyzeOWASPTop10(findings []APIFinding) OWASPTop10 {
	owasp := OWASPTop10{}

	for _, finding := range findings {
		if finding.Metadata["owasp_top10"] != nil {
			top10 := finding.Metadata["owasp_top10"].(string)
			switch top10 {
			case "A01:2021":
				owasp.BrokenObjectLevelAuthorization = true
			case "A02:2021":
				owasp.BrokenUserAuthentication = true
			case "A03:2021":
				owasp.ExcessiveDataExposure = true
			case "A04:2021":
				owasp.LackOfResourcesAndRateLimiting = true
			case "A05:2021":
				owasp.BrokenFunctionLevelAuthorization = true
			case "A06:2021":
				owasp.MassAssignment = true
			case "A07:2021":
				owasp.SecurityMisconfiguration = true
			case "A08:2021":
				owasp.Injection = true
			case "A09:2021":
				owasp.ImproperAssetsManagement = true
			case "A10:2021":
				owasp.InsufficientLoggingAndMonitoring = true
			}
		}
	}

	return owasp
}
