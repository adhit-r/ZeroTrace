package services

import (
	"fmt"
	"strings"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TechStackService handles technology stack analysis
type TechStackService struct {
	db *gorm.DB
}

// NewTechStackService creates a new tech stack service
func NewTechStackService(db *gorm.DB) *TechStackService {
	return &TechStackService{
		db: db,
	}
}

// TechStackAnalysis represents the result of tech stack analysis
type TechStackAnalysis struct {
	OrganizationID    uuid.UUID          `json:"organization_id"`
	DetectedTechStack models.TechStack   `json:"detected_tech_stack"`
	RelevanceScores   map[string]float64 `json:"relevance_scores"`
	RiskFactors       []string           `json:"risk_factors"`
	Recommendations   []string           `json:"recommendations"`
	LastAnalyzed      time.Time          `json:"last_analyzed"`
	ConfidenceScore   float64            `json:"confidence_score"`
}

// AnalyzeTechStackFromAssets analyzes technology stack from scanned assets
func (s *TechStackService) AnalyzeTechStackFromAssets(organizationID uuid.UUID) (*TechStackAnalysis, error) {
	// Get all agents for the organization
	var agents []models.Agent
	err := s.db.Where("organization_id = ?", organizationID).Find(&agents).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get agents: %w", err)
	}

	// Get all agent scan results for the organization
	var agentScanResults []models.AgentScanResult
	err = s.db.Where("company_id = ?", organizationID.String()).Find(&agentScanResults).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get agent scan results: %w", err)
	}

	// Analyze technology stack from scan results
	detectedTech := s.detectTechnologiesFromScans(agentScanResults)

	// Calculate relevance scores for different vulnerability types
	relevanceScores := s.calculateRelevanceScores(detectedTech)

	// Identify risk factors
	riskFactors := s.identifyRiskFactors(detectedTech, agentScanResults)

	// Generate recommendations
	recommendations := s.generateRecommendations(detectedTech, riskFactors)

	// Calculate confidence score
	confidenceScore := s.calculateConfidenceScore(detectedTech, len(agentScanResults))

	analysis := &TechStackAnalysis{
		OrganizationID:    organizationID,
		DetectedTechStack: detectedTech,
		RelevanceScores:   relevanceScores,
		RiskFactors:       riskFactors,
		Recommendations:   recommendations,
		LastAnalyzed:      time.Now(),
		ConfidenceScore:   confidenceScore,
	}

	return analysis, nil
}

// detectTechnologiesFromScans analyzes scan results to detect technologies
func (s *TechStackService) detectTechnologiesFromScans(scanResults []models.AgentScanResult) models.TechStack {
	techStack := models.TechStack{
		Languages:        make([]string, 0),
		Frameworks:       make([]string, 0),
		Databases:        make([]string, 0),
		CloudProviders:   make([]string, 0),
		OperatingSystems: make([]string, 0),
		Containers:       make([]string, 0),
		DevTools:         make([]string, 0),
		SecurityTools:    make([]string, 0),
	}

	// Technology detection maps
	languageMap := make(map[string]int)
	frameworkMap := make(map[string]int)
	databaseMap := make(map[string]int)
	cloudMap := make(map[string]int)
	osMap := make(map[string]int)
	containerMap := make(map[string]int)
	devToolMap := make(map[string]int)
	securityToolMap := make(map[string]int)

	for _, scanResult := range scanResults {
		// Analyze vulnerabilities for technology indicators
		for _, vuln := range scanResult.Vulnerabilities {
			tech := s.extractTechnologyFromVulnerability(vuln)
			if tech.Language != "" {
				languageMap[tech.Language]++
			}
			if tech.Framework != "" {
				frameworkMap[tech.Framework]++
			}
			if tech.Database != "" {
				databaseMap[tech.Database]++
			}
		}

		// Analyze dependencies for technology indicators
		for _, dep := range scanResult.Dependencies {
			tech := s.extractTechnologyFromDependency(dep)
			if tech.Language != "" {
				languageMap[tech.Language]++
			}
			if tech.Framework != "" {
				frameworkMap[tech.Framework]++
			}
			if tech.Database != "" {
				databaseMap[tech.Database]++
			}
		}

		// Analyze metadata for additional technology indicators
		if metadata, ok := scanResult.Metadata["technologies"]; ok {
			if techs, ok := metadata.([]interface{}); ok {
				for _, tech := range techs {
					if techStr, ok := tech.(string); ok {
						s.categorizeTechnology(techStr, languageMap, frameworkMap, databaseMap, cloudMap, osMap, containerMap, devToolMap, securityToolMap)
					}
				}
			}
		}
	}

	// Convert maps to slices, keeping only technologies with frequency > 1
	techStack.Languages = s.mapToSlice(languageMap, 1)
	techStack.Frameworks = s.mapToSlice(frameworkMap, 1)
	techStack.Databases = s.mapToSlice(databaseMap, 1)
	techStack.CloudProviders = s.mapToSlice(cloudMap, 1)
	techStack.OperatingSystems = s.mapToSlice(osMap, 1)
	techStack.Containers = s.mapToSlice(containerMap, 1)
	techStack.DevTools = s.mapToSlice(devToolMap, 1)
	techStack.SecurityTools = s.mapToSlice(securityToolMap, 1)

	return techStack
}

// TechnologyInfo represents detected technology information
type TechnologyInfo struct {
	Language  string
	Framework string
	Database  string
}

// extractTechnologyFromVulnerability extracts technology information from vulnerability
func (s *TechStackService) extractTechnologyFromVulnerability(vuln models.Vulnerability) TechnologyInfo {
	tech := TechnologyInfo{}

	// Analyze package name and type
	packageName := strings.ToLower(vuln.PackageName)
	vulnType := strings.ToLower(vuln.Type)

	// Language detection
	if strings.Contains(packageName, "python") || strings.Contains(vulnType, "python") {
		tech.Language = "Python"
	} else if strings.Contains(packageName, "node") || strings.Contains(packageName, "npm") || strings.Contains(vulnType, "javascript") {
		tech.Language = "JavaScript"
	} else if strings.Contains(packageName, "java") || strings.Contains(vulnType, "java") {
		tech.Language = "Java"
	} else if strings.Contains(packageName, "go") || strings.Contains(vulnType, "golang") {
		tech.Language = "Go"
	} else if strings.Contains(packageName, "ruby") || strings.Contains(vulnType, "ruby") {
		tech.Language = "Ruby"
	} else if strings.Contains(packageName, "php") || strings.Contains(vulnType, "php") {
		tech.Language = "PHP"
	} else if strings.Contains(packageName, "csharp") || strings.Contains(packageName, "dotnet") || strings.Contains(vulnType, "csharp") {
		tech.Language = "C#"
	}

	// Framework detection
	if strings.Contains(packageName, "django") {
		tech.Framework = "Django"
	} else if strings.Contains(packageName, "flask") {
		tech.Framework = "Flask"
	} else if strings.Contains(packageName, "express") {
		tech.Framework = "Express.js"
	} else if strings.Contains(packageName, "react") {
		tech.Framework = "React"
	} else if strings.Contains(packageName, "angular") {
		tech.Framework = "Angular"
	} else if strings.Contains(packageName, "spring") {
		tech.Framework = "Spring"
	} else if strings.Contains(packageName, "rails") {
		tech.Framework = "Ruby on Rails"
	}

	// Database detection
	if strings.Contains(packageName, "mysql") {
		tech.Database = "MySQL"
	} else if strings.Contains(packageName, "postgres") {
		tech.Database = "PostgreSQL"
	} else if strings.Contains(packageName, "mongodb") {
		tech.Database = "MongoDB"
	} else if strings.Contains(packageName, "redis") {
		tech.Database = "Redis"
	} else if strings.Contains(packageName, "sqlite") {
		tech.Database = "SQLite"
	}

	return tech
}

// extractTechnologyFromDependency extracts technology information from dependency
func (s *TechStackService) extractTechnologyFromDependency(dep models.Dependency) TechnologyInfo {
	tech := TechnologyInfo{}

	depName := strings.ToLower(dep.Name)
	depType := strings.ToLower(dep.Type)

	// Language detection based on dependency type
	switch depType {
	case "npm", "yarn":
		tech.Language = "JavaScript"
	case "pip", "conda":
		tech.Language = "Python"
	case "maven", "gradle":
		tech.Language = "Java"
	case "go mod", "go get":
		tech.Language = "Go"
	case "gem", "bundler":
		tech.Language = "Ruby"
	case "composer":
		tech.Language = "PHP"
	case "nuget":
		tech.Language = "C#"
	}

	// Framework detection
	if strings.Contains(depName, "django") {
		tech.Framework = "Django"
	} else if strings.Contains(depName, "flask") {
		tech.Framework = "Flask"
	} else if strings.Contains(depName, "express") {
		tech.Framework = "Express.js"
	} else if strings.Contains(depName, "react") {
		tech.Framework = "React"
	} else if strings.Contains(depName, "angular") {
		tech.Framework = "Angular"
	} else if strings.Contains(depName, "spring") {
		tech.Framework = "Spring"
	} else if strings.Contains(depName, "rails") {
		tech.Framework = "Ruby on Rails"
	}

	return tech
}

// categorizeTechnology categorizes a technology string into appropriate categories
func (s *TechStackService) categorizeTechnology(tech string, languageMap, frameworkMap, databaseMap, cloudMap, osMap, containerMap, devToolMap, securityToolMap map[string]int) {
	techLower := strings.ToLower(tech)

	// Language detection
	languages := []string{"python", "javascript", "java", "go", "ruby", "php", "csharp", "c++", "c", "rust", "swift", "kotlin", "scala", "typescript"}
	for _, lang := range languages {
		if strings.Contains(techLower, lang) {
			languageMap[strings.Title(lang)]++
			return
		}
	}

	// Framework detection
	frameworks := []string{"django", "flask", "express", "react", "angular", "vue", "spring", "rails", "laravel", "symfony", "asp.net", "fastapi"}
	for _, framework := range frameworks {
		if strings.Contains(techLower, framework) {
			frameworkMap[strings.Title(framework)]++
			return
		}
	}

	// Database detection
	databases := []string{"mysql", "postgresql", "mongodb", "redis", "sqlite", "oracle", "sql server", "cassandra", "elasticsearch"}
	for _, db := range databases {
		if strings.Contains(techLower, db) {
			databaseMap[strings.Title(db)]++
			return
		}
	}

	// Cloud provider detection
	clouds := []string{"aws", "azure", "gcp", "google cloud", "amazon web services", "microsoft azure"}
	for _, cloud := range clouds {
		if strings.Contains(techLower, cloud) {
			cloudMap[strings.ToUpper(cloud)]++
			return
		}
	}

	// Operating system detection
	oses := []string{"linux", "windows", "macos", "ubuntu", "centos", "debian", "redhat", "suse"}
	for _, os := range oses {
		if strings.Contains(techLower, os) {
			osMap[strings.Title(os)]++
			return
		}
	}

	// Container detection
	containers := []string{"docker", "kubernetes", "podman", "containerd", "rkt"}
	for _, container := range containers {
		if strings.Contains(techLower, container) {
			containerMap[strings.Title(container)]++
			return
		}
	}

	// Development tool detection
	devTools := []string{"git", "jenkins", "github", "gitlab", "bitbucket", "jira", "confluence", "slack", "vscode", "intellij"}
	for _, tool := range devTools {
		if strings.Contains(techLower, tool) {
			devToolMap[strings.Title(tool)]++
			return
		}
	}

	// Security tool detection
	securityTools := []string{"nessus", "qualys", "burp", "nmap", "wireshark", "metasploit", "owasp", "snyk", "veracode", "checkmarx"}
	for _, tool := range securityTools {
		if strings.Contains(techLower, tool) {
			securityToolMap[strings.Title(tool)]++
			return
		}
	}
}

// calculateRelevanceScores calculates relevance scores for different vulnerability types
func (s *TechStackService) calculateRelevanceScores(techStack models.TechStack) map[string]float64 {
	scores := make(map[string]float64)

	// Calculate scores based on technology diversity and coverage
	totalTechnologies := len(techStack.Languages) + len(techStack.Frameworks) + len(techStack.Databases)

	if totalTechnologies > 0 {
		// Higher diversity = higher relevance for comprehensive security
		scores["technology_diversity"] = float64(totalTechnologies) / 20.0 // Normalize to 0-1
		if scores["technology_diversity"] > 1.0 {
			scores["technology_diversity"] = 1.0
		}
	} else {
		scores["technology_diversity"] = 0.1 // Minimum score
	}

	// Calculate risk scores based on technology types
	scores["web_application_risk"] = s.calculateWebAppRisk(techStack)
	scores["database_risk"] = s.calculateDatabaseRisk(techStack)
	scores["cloud_risk"] = s.calculateCloudRisk(techStack)
	scores["container_risk"] = s.calculateContainerRisk(techStack)

	return scores
}

// calculateWebAppRisk calculates risk score for web applications
func (s *TechStackService) calculateWebAppRisk(techStack models.TechStack) float64 {
	risk := 0.0

	// Higher risk for more web technologies
	risk += float64(len(techStack.Languages)) * 0.1
	risk += float64(len(techStack.Frameworks)) * 0.15

	// Specific high-risk technologies
	for _, framework := range techStack.Frameworks {
		if strings.Contains(strings.ToLower(framework), "django") ||
			strings.Contains(strings.ToLower(framework), "flask") ||
			strings.Contains(strings.ToLower(framework), "express") {
			risk += 0.2
		}
	}

	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// calculateDatabaseRisk calculates risk score for databases
func (s *TechStackService) calculateDatabaseRisk(techStack models.TechStack) float64 {
	risk := float64(len(techStack.Databases)) * 0.2

	// Higher risk for more databases
	if len(techStack.Databases) > 3 {
		risk += 0.3
	}

	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// calculateCloudRisk calculates risk score for cloud technologies
func (s *TechStackService) calculateCloudRisk(techStack models.TechStack) float64 {
	risk := float64(len(techStack.CloudProviders)) * 0.3

	// Higher risk for multiple cloud providers
	if len(techStack.CloudProviders) > 1 {
		risk += 0.2
	}

	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// calculateContainerRisk calculates risk score for container technologies
func (s *TechStackService) calculateContainerRisk(techStack models.TechStack) float64 {
	risk := float64(len(techStack.Containers)) * 0.4

	// Higher risk for containerized environments
	if len(techStack.Containers) > 0 {
		risk += 0.3
	}

	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// identifyRiskFactors identifies potential risk factors in the technology stack
func (s *TechStackService) identifyRiskFactors(techStack models.TechStack, scanResults []models.AgentScanResult) []string {
	riskFactors := make([]string, 0)

	// High-risk technology combinations
	if len(techStack.Languages) > 5 {
		riskFactors = append(riskFactors, "High technology diversity increases attack surface")
	}

	if len(techStack.Databases) > 2 {
		riskFactors = append(riskFactors, "Multiple database systems increase complexity and risk")
	}

	if len(techStack.CloudProviders) > 1 {
		riskFactors = append(riskFactors, "Multi-cloud environment increases security complexity")
	}

	// Outdated or vulnerable technologies
	for _, framework := range techStack.Frameworks {
		if strings.Contains(strings.ToLower(framework), "legacy") {
			riskFactors = append(riskFactors, "Legacy frameworks may have unpatched vulnerabilities")
		}
	}

	// Container security risks
	if len(techStack.Containers) > 0 {
		riskFactors = append(riskFactors, "Containerized environments require additional security measures")
	}

	// High vulnerability count
	totalVulns := 0
	for _, result := range scanResults {
		totalVulns += len(result.Vulnerabilities)
	}
	if totalVulns > 100 {
		riskFactors = append(riskFactors, "High number of vulnerabilities detected")
	}

	return riskFactors
}

// generateRecommendations generates security recommendations based on tech stack
func (s *TechStackService) generateRecommendations(techStack models.TechStack, riskFactors []string) []string {
	recommendations := make([]string, 0)

	// Technology-specific recommendations
	for _, language := range techStack.Languages {
		switch strings.ToLower(language) {
		case "python":
			recommendations = append(recommendations, "Implement Python security best practices and dependency scanning")
		case "javascript":
			recommendations = append(recommendations, "Use npm audit and implement Node.js security measures")
		case "java":
			recommendations = append(recommendations, "Implement Java security scanning and dependency management")
		}
	}

	// Framework-specific recommendations
	for _, framework := range techStack.Frameworks {
		switch strings.ToLower(framework) {
		case "django":
			recommendations = append(recommendations, "Configure Django security settings and middleware")
		case "flask":
			recommendations = append(recommendations, "Implement Flask security extensions and best practices")
		case "express":
			recommendations = append(recommendations, "Use Express.js security middleware and input validation")
		}
	}

	// Database recommendations
	if len(techStack.Databases) > 0 {
		recommendations = append(recommendations, "Implement database security scanning and access controls")
	}

	// Cloud recommendations
	if len(techStack.CloudProviders) > 0 {
		recommendations = append(recommendations, "Implement cloud security posture management")
	}

	// Container recommendations
	if len(techStack.Containers) > 0 {
		recommendations = append(recommendations, "Implement container security scanning and runtime protection")
	}

	// General recommendations
	if len(riskFactors) > 3 {
		recommendations = append(recommendations, "Consider consolidating technology stack to reduce complexity")
	}

	return recommendations
}

// calculateConfidenceScore calculates confidence score for the analysis
func (s *TechStackService) calculateConfidenceScore(techStack models.TechStack, scanResultCount int) float64 {
	confidence := 0.0

	// Base confidence from scan result count
	if scanResultCount > 0 {
		confidence += 0.3
	}
	if scanResultCount > 10 {
		confidence += 0.2
	}
	if scanResultCount > 50 {
		confidence += 0.2
	}

	// Technology detection confidence
	totalTechnologies := len(techStack.Languages) + len(techStack.Frameworks) + len(techStack.Databases)
	if totalTechnologies > 0 {
		confidence += 0.3
	}
	if totalTechnologies > 5 {
		confidence += 0.2
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// mapToSlice converts a map to a slice, filtering by minimum frequency
func (s *TechStackService) mapToSlice(techMap map[string]int, minFreq int) []string {
	result := make([]string, 0)
	for tech, freq := range techMap {
		if freq >= minFreq {
			result = append(result, tech)
		}
	}
	return result
}
