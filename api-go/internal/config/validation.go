package config

import (
	"fmt"
	"os"
)

// Validate checks that required configuration is present
func (c *Config) Validate() error {
	// In production (non-debug mode), require Clerk JWT key
	if !c.Debug {
		if c.ClerkJWTVerificationKey == "" {
			return fmt.Errorf("CLERK_JWT_VERIFICATION_KEY is required in production mode")
		}
		if c.ClerkJWTVerificationKey == "dev-clerk-key-change-in-production" || c.ClerkJWTVerificationKey == "development-key" {
			return fmt.Errorf("CLERK_JWT_VERIFICATION_KEY must not use development default in production")
		}
	}

	// Validate database configuration
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}

	// Validate enrichment service URL
	if c.EnrichmentServiceURL == "" {
		return fmt.Errorf("ENRICHMENT_SERVICE_URL is required")
	}

	return nil
}

// ValidateEnvironment checks environment variables at startup
func ValidateEnvironment() error {
	requiredVars := []string{
		"DB_HOST",
		"DB_NAME",
		"DB_USER",
	}

	isDebug := os.Getenv("API_MODE") == "debug" || os.Getenv("DEBUG") == "true"
	if !isDebug {
		// In production, also require Clerk key
		requiredVars = append(requiredVars, "CLERK_JWT_VERIFICATION_KEY")
	}

	var missing []string
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	return nil
}

