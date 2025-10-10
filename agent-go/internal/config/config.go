package config

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// Config holds application configuration
type Config struct {
	// Agent Configuration
	AgentID     string `json:"agent_id"`
	APIURL      string `json:"api_url"`
	APIEndpoint string `json:"api_endpoint"`
	APIKey      string `json:"api_key"`
	APITimeout  int    `json:"api_timeout"`
	LogLevel    string `json:"log_level"`
	Debug       bool   `json:"debug"`

	// Enrollment Configuration
	EnrollmentToken string `json:"enrollment_token"`
	AgentCredential string `json:"agent_credential"`
	OrganizationID  string `json:"organization_id"`

	// Company-specific Configuration (legacy - will be replaced by enrollment)
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`
	CompanySlug string `json:"company_slug"`

	// System Information
	Hostname string `json:"hostname"`
	OS       string `json:"os"`

	// API Configuration
	APIPort int `json:"api_port"`

	// Scan Configuration
	ScanInterval    time.Duration `json:"scan_interval"`
	ScanDepth       int           `json:"scan_depth"`
	MaxFileSize     int64         `json:"max_file_size"`
	ExcludePatterns []string      `json:"exclude_patterns"`
	IncludePatterns []string      `json:"include_patterns"`

	// Database Configuration
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBName     string `json:"db_name"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBSSLMode  string `json:"db_ssl_mode"`
}

// Load loads configuration from environment variables
func Load() *Config {
	apiPort, _ := strconv.Atoi(getEnv("API_PORT", "8080"))
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	debug, _ := strconv.ParseBool(getEnv("DEBUG", "false"))

	// Generate UUID for agent ID if not set or not a valid UUID
	agentID := getEnv("AGENT_ID", "")
	if agentID == "" {
		agentID = uuid.New().String()
	} else {
		// Validate if it's a UUID, if not generate a new one
		if _, err := uuid.Parse(agentID); err != nil {
			agentID = uuid.New().String()
		}
	}

	return &Config{
		// Agent Configuration
		AgentID:     agentID,
		APIURL:      getEnv("API_URL", "http://localhost:8080"),
		APIEndpoint: getEnv("API_ENDPOINT", "http://localhost:8080"),
		APIKey:      getEnv("API_KEY", ""),
		APITimeout:  30, // 30 seconds default
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Debug:       debug,

		// Enrollment Configuration
		EnrollmentToken: getEnv("ENROLLMENT_TOKEN", ""),
		AgentCredential: getEnv("AGENT_CREDENTIAL", ""),
		OrganizationID:  getEnv("ORGANIZATION_ID", ""),

		// Company-specific Configuration (legacy)
		CompanyID:   getEnv("COMPANY_ID", ""),
		CompanyName: getEnv("COMPANY_NAME", ""),
		CompanySlug: getEnv("COMPANY_SLUG", ""),

		// System Information
		Hostname: getEnv("HOSTNAME", getHostname()),
		OS:       getEnv("OS", getOS()),

		// API Configuration
		APIPort: apiPort,

		// Scan Configuration
		ScanInterval:    5 * time.Minute,  // Default 5 minutes
		ScanDepth:       3,                // Default depth 3
		MaxFileSize:     10 * 1024 * 1024, // 10MB default
		ExcludePatterns: []string{".git", "node_modules", ".DS_Store", "*.log"},
		IncludePatterns: []string{".go", ".py", ".js", ".ts", ".java", ".php", ".rb", ".rs", ".cpp", ".c", ".cs"},

		// Database Configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBName:     getEnv("DB_NAME", "zerotrace"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// IsEnrolled checks if the agent is enrolled with an organization
func (c *Config) IsEnrolled() bool {
	return c.AgentCredential != "" && c.OrganizationID != ""
}

// HasEnrollmentToken checks if an enrollment token is available
func (c *Config) HasEnrollmentToken() bool {
	return c.EnrollmentToken != ""
}

// IsCompanyConfigured checks if company-specific configuration is set (legacy)
func (c *Config) IsCompanyConfigured() bool {
	return c.CompanyID != "" && c.CompanyName != "" && c.CompanySlug != ""
}

// GetCompanyIdentifier returns a unique identifier for the company (legacy)
func (c *Config) GetCompanyIdentifier() string {
	if c.CompanySlug != "" {
		return c.CompanySlug
	}
	return c.CompanyID
}

// GetOrganizationIdentifier returns the organization identifier
func (c *Config) GetOrganizationIdentifier() string {
	return c.OrganizationID
}

// getHostname gets the system hostname
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// getOS gets the operating system name
func getOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return runtime.GOOS
	}
}
