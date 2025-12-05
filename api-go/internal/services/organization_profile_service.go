package services

import (
	"fmt"
	"strings"
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OrganizationProfileService handles organization profile operations
type OrganizationProfileService struct {
	db *gorm.DB
}

// NewOrganizationProfileService creates a new organization profile service
func NewOrganizationProfileService(db *gorm.DB) *OrganizationProfileService {
	return &OrganizationProfileService{
		db: db,
	}
}

// CreateOrganizationProfile creates a new organization profile
func (s *OrganizationProfileService) CreateOrganizationProfile(req *models.CreateOrganizationProfileRequest) (*models.OrganizationProfile, error) {
	// Check if profile already exists
	var existingProfile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", req.OrganizationID).First(&existingProfile).Error
	if err == nil {
		return nil, fmt.Errorf("organization profile already exists for organization %s", req.OrganizationID)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check existing profile: %w", err)
	}

	// Create new profile
	profile := &models.OrganizationProfile{
		ID:                   uuid.New(),
		OrganizationID:       req.OrganizationID,
		Industry:             req.Industry,
		RiskTolerance:        req.RiskTolerance,
		TechStack:            req.TechStack,
		ComplianceFrameworks: req.ComplianceFrameworks,
		SecurityPolicies:     req.SecurityPolicies,
		RiskWeights:          req.RiskWeights,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Set default risk weights if not provided
	if profile.RiskWeights == nil {
		profile.RiskWeights = s.getDefaultRiskWeights(req.Industry, req.RiskTolerance)
	}

	// Set default security policies if not provided
	if profile.SecurityPolicies == nil {
		profile.SecurityPolicies = s.getDefaultSecurityPolicies(req.Industry)
	}

	err = s.db.Create(profile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create organization profile: %w", err)
	}

	return profile, nil
}

// GetOrganizationProfile retrieves an organization profile by organization ID
func (s *OrganizationProfileService) GetOrganizationProfile(organizationID uuid.UUID) (*models.OrganizationProfile, error) {
	var profile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization profile not found for organization %s", organizationID)
		}
		return nil, fmt.Errorf("failed to get organization profile: %w", err)
	}

	return &profile, nil
}

// UpdateOrganizationProfile updates an existing organization profile
func (s *OrganizationProfileService) UpdateOrganizationProfile(organizationID uuid.UUID, req *models.UpdateOrganizationProfileRequest) (*models.OrganizationProfile, error) {
	var profile models.OrganizationProfile
	err := s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("organization profile not found for organization %s", organizationID)
		}
		return nil, fmt.Errorf("failed to get organization profile: %w", err)
	}

	// Update fields if provided
	updates := make(map[string]interface{})

	if req.Industry != nil {
		updates["industry"] = *req.Industry
	}
	if req.RiskTolerance != nil {
		updates["risk_tolerance"] = *req.RiskTolerance
	}
	if req.TechStack != nil {
		updates["tech_stack"] = *req.TechStack
	}
	if req.ComplianceFrameworks != nil {
		updates["compliance_frameworks"] = *req.ComplianceFrameworks
	}
	if req.SecurityPolicies != nil {
		updates["security_policies"] = req.SecurityPolicies
	}
	if req.RiskWeights != nil {
		updates["risk_weights"] = req.RiskWeights
	}

	updates["updated_at"] = time.Now()

	err = s.db.Model(&profile).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update organization profile: %w", err)
	}

	// Reload the profile to get updated data
	err = s.db.Where("organization_id = ?", organizationID).First(&profile).Error
	if err != nil {
		return nil, fmt.Errorf("failed to reload organization profile: %w", err)
	}

	return &profile, nil
}

// DeleteOrganizationProfile deletes an organization profile
func (s *OrganizationProfileService) DeleteOrganizationProfile(organizationID uuid.UUID) error {
	result := s.db.Where("organization_id = ?", organizationID).Delete(&models.OrganizationProfile{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete organization profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("organization profile not found for organization %s", organizationID)
	}

	return nil
}

// GetTechStackRelevance calculates relevance score for vulnerabilities based on tech stack
func (s *OrganizationProfileService) GetTechStackRelevance(organizationID uuid.UUID, vulnerability *models.Vulnerability) (float64, error) {
	profile, err := s.GetOrganizationProfile(organizationID)
	if err != nil {
		return 0, err
	}

	// Calculate relevance based on tech stack
	relevanceScore := 0.0
	techStack := profile.TechStack

	// Check if vulnerability affects technologies in the organization's stack
	if vulnerability.PackageName != "" {
		// Check against languages
		for _, lang := range techStack.Languages {
			if s.isTechnologyMatch(vulnerability.PackageName, lang) {
				relevanceScore += 0.3
			}
		}

		// Check against frameworks
		for _, framework := range techStack.Frameworks {
			if s.isTechnologyMatch(vulnerability.PackageName, framework) {
				relevanceScore += 0.4
			}
		}

		// Check against databases
		for _, db := range techStack.Databases {
			if s.isTechnologyMatch(vulnerability.PackageName, db) {
				relevanceScore += 0.3
			}
		}
	}

	// Apply risk tolerance multiplier
	switch profile.RiskTolerance {
	case models.RiskToleranceConservative:
		relevanceScore *= 1.2 // Higher relevance for conservative organizations
	case models.RiskToleranceModerate:
		relevanceScore *= 1.0 // No change
	case models.RiskToleranceAggressive:
		relevanceScore *= 0.8 // Lower relevance for aggressive organizations
	}

	// Cap at 1.0
	if relevanceScore > 1.0 {
		relevanceScore = 1.0
	}

	return relevanceScore, nil
}

// GetIndustryRiskWeights returns industry-specific risk weights
func (s *OrganizationProfileService) GetIndustryRiskWeights(organizationID uuid.UUID) (map[string]float64, error) {
	profile, err := s.GetOrganizationProfile(organizationID)
	if err != nil {
		return nil, err
	}

	// Convert risk weights to float64 map
	weights := make(map[string]float64)
	for key, value := range profile.RiskWeights {
		if floatVal, ok := value.(float64); ok {
			weights[key] = floatVal
		}
	}

	return weights, nil
}

// isTechnologyMatch checks if a vulnerability package matches a technology
func (s *OrganizationProfileService) isTechnologyMatch(packageName, technology string) bool {
	// Simple string matching - can be enhanced with fuzzy matching
	packageName = strings.ToLower(packageName)
	technology = strings.ToLower(technology)

	return strings.Contains(packageName, technology) || strings.Contains(technology, packageName)
}

// getDefaultRiskWeights returns default risk weights based on industry and risk tolerance
func (s *OrganizationProfileService) getDefaultRiskWeights(industry string, riskTolerance models.RiskTolerance) map[string]any {
	baseWeights := map[string]float64{
		"critical": 1.0,
		"high":     0.8,
		"medium":   0.6,
		"low":      0.4,
		"info":     0.2,
	}

	// Industry-specific adjustments
	switch industry {
	case "healthcare":
		baseWeights["critical"] = 1.2
		baseWeights["high"] = 1.0
	case "finance":
		baseWeights["critical"] = 1.1
		baseWeights["high"] = 0.9
	case "government":
		baseWeights["critical"] = 1.3
		baseWeights["high"] = 1.1
	}

	// Risk tolerance adjustments
	switch riskTolerance {
	case models.RiskToleranceConservative:
		for key := range baseWeights {
			baseWeights[key] *= 1.2
		}
	case models.RiskToleranceAggressive:
		for key := range baseWeights {
			baseWeights[key] *= 0.8
		}
	}

	// Convert to interface{} map
	weights := make(map[string]any)
	for key, value := range baseWeights {
		weights[key] = value
	}

	return weights
}

// getDefaultSecurityPolicies returns default security policies based on industry
func (s *OrganizationProfileService) getDefaultSecurityPolicies(industry string) map[string]any {
	policies := map[string]any{
		"patch_management": map[string]any{
			"critical_patches": "immediate",
			"high_patches":     "24_hours",
			"medium_patches":   "7_days",
			"low_patches":      "30_days",
		},
		"vulnerability_management": map[string]any{
			"scan_frequency": "daily",
			"reporting":      "real_time",
		},
		"compliance": map[string]any{
			"enabled":    true,
			"frameworks": []string{},
		},
	}

	// Industry-specific policies
	switch industry {
	case "healthcare":
		policies["compliance"] = map[string]any{
			"enabled":    true,
			"frameworks": []string{"HIPAA", "HITECH"},
		}
		policies["data_protection"] = map[string]any{
			"encryption_required": true,
			"access_controls":     "strict",
		}
	case "finance":
		policies["compliance"] = map[string]any{
			"enabled":    true,
			"frameworks": []string{"PCI DSS", "SOX", "Basel III"},
		}
		policies["audit_requirements"] = map[string]any{
			"logging":    "comprehensive",
			"retention":  "7_years",
			"monitoring": "continuous",
		}
	case "government":
		policies["compliance"] = map[string]any{
			"enabled":    true,
			"frameworks": []string{"FISMA", "FedRAMP", "NIST"},
		}
		policies["security_clearance"] = map[string]any{
			"required": true,
			"levels":   []string{"public", "confidential", "secret"},
		}
	}

	return policies
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

// TechnologyInfo represents detected technology information
type TechnologyInfo struct {
	Language  string
	Framework string
	Database  string
}

// AnalyzeTechStackFromAssets analyzes technology stack from scanned assets
func (s *OrganizationProfileService) AnalyzeTechStackFromAssets(organizationID uuid.UUID) (*TechStackAnalysis, error) {
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

// GetTechStackRecommendations returns recommendations based on tech stack analysis
func (s *OrganizationProfileService) GetTechStackRecommendations(organizationID uuid.UUID) ([]string, error) {
	analysis, err := s.AnalyzeTechStackFromAssets(organizationID)
	if err != nil {
		return nil, err
	}
	return analysis.Recommendations, nil
}

// detectTechnologiesFromScans analyzes scan results to detect technologies
func (s *OrganizationProfileService) detectTechnologiesFromScans(scanResults []models.AgentScanResult) models.TechStack {
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

// extractTechnologyFromVulnerability extracts technology information from vulnerability
func (s *OrganizationProfileService) extractTechnologyFromVulnerability(vuln models.Vulnerability) TechnologyInfo {
	tech := TechnologyInfo{}

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
func (s *OrganizationProfileService) extractTechnologyFromDependency(dep models.Dependency) TechnologyInfo {
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
func (s *OrganizationProfileService) categorizeTechnology(tech string, languageMap, frameworkMap, databaseMap, cloudMap, osMap, containerMap, devToolMap, securityToolMap map[string]int) {
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
func (s *OrganizationProfileService) calculateRelevanceScores(techStack models.TechStack) map[string]float64 {
	scores := make(map[string]float64)

	totalTechnologies := len(techStack.Languages) + len(techStack.Frameworks) + len(techStack.Databases)

	if totalTechnologies > 0 {
		scores["technology_diversity"] = float64(totalTechnologies) / 20.0
		if scores["technology_diversity"] > 1.0 {
			scores["technology_diversity"] = 1.0
		}
	} else {
		scores["technology_diversity"] = 0.1
	}

	scores["web_application_risk"] = s.calculateWebAppRisk(techStack)
	scores["database_risk"] = s.calculateDatabaseRisk(techStack)
	scores["cloud_risk"] = s.calculateCloudRisk(techStack)
	scores["container_risk"] = s.calculateContainerRisk(techStack)

	return scores
}

// calculateWebAppRisk calculates risk score for web applications
func (s *OrganizationProfileService) calculateWebAppRisk(techStack models.TechStack) float64 {
	risk := 0.0
	risk += float64(len(techStack.Languages)) * 0.1
	risk += float64(len(techStack.Frameworks)) * 0.15

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
func (s *OrganizationProfileService) calculateDatabaseRisk(techStack models.TechStack) float64 {
	risk := float64(len(techStack.Databases)) * 0.2
	if len(techStack.Databases) > 3 {
		risk += 0.3
	}
	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// calculateCloudRisk calculates risk score for cloud technologies
func (s *OrganizationProfileService) calculateCloudRisk(techStack models.TechStack) float64 {
	risk := float64(len(techStack.CloudProviders)) * 0.3
	if len(techStack.CloudProviders) > 1 {
		risk += 0.2
	}
	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// calculateContainerRisk calculates risk score for container technologies
func (s *OrganizationProfileService) calculateContainerRisk(techStack models.TechStack) float64 {
	risk := float64(len(techStack.Containers)) * 0.4
	if len(techStack.Containers) > 0 {
		risk += 0.3
	}
	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

// identifyRiskFactors identifies potential risk factors in the technology stack
func (s *OrganizationProfileService) identifyRiskFactors(techStack models.TechStack, scanResults []models.AgentScanResult) []string {
	riskFactors := make([]string, 0)

	if len(techStack.Languages) > 5 {
		riskFactors = append(riskFactors, "High technology diversity increases attack surface")
	}

	if len(techStack.Databases) > 2 {
		riskFactors = append(riskFactors, "Multiple database systems increase complexity and risk")
	}

	if len(techStack.CloudProviders) > 1 {
		riskFactors = append(riskFactors, "Multi-cloud environment increases security complexity")
	}

	for _, framework := range techStack.Frameworks {
		if strings.Contains(strings.ToLower(framework), "legacy") {
			riskFactors = append(riskFactors, "Legacy frameworks may have unpatched vulnerabilities")
		}
	}

	if len(techStack.Containers) > 0 {
		riskFactors = append(riskFactors, "Containerized environments require additional security measures")
	}

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
func (s *OrganizationProfileService) generateRecommendations(techStack models.TechStack, riskFactors []string) []string {
	recommendations := make([]string, 0)

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

	if len(techStack.Databases) > 0 {
		recommendations = append(recommendations, "Implement database security scanning and access controls")
	}

	if len(techStack.CloudProviders) > 0 {
		recommendations = append(recommendations, "Implement cloud security posture management")
	}

	if len(techStack.Containers) > 0 {
		recommendations = append(recommendations, "Implement container security scanning and runtime protection")
	}

	if len(riskFactors) > 3 {
		recommendations = append(recommendations, "Consider consolidating technology stack to reduce complexity")
	}

	return recommendations
}

// calculateConfidenceScore calculates confidence score for the analysis
func (s *OrganizationProfileService) calculateConfidenceScore(techStack models.TechStack, scanResultCount int) float64 {
	confidence := 0.0

	if scanResultCount > 0 {
		confidence += 0.3
	}
	if scanResultCount > 10 {
		confidence += 0.2
	}
	if scanResultCount > 50 {
		confidence += 0.2
	}

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
func (s *OrganizationProfileService) mapToSlice(techMap map[string]int, minFreq int) []string {
	result := make([]string, 0)
	for tech, freq := range techMap {
		if freq >= minFreq {
			result = append(result, tech)
		}
	}
	return result
}
