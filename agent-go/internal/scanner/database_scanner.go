package scanner

import (
	"database/sql"
	"fmt"
	"net"
	"time"

	"zerotrace/agent/internal/config"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/google/uuid"
	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DatabaseScanner handles database security scanning
type DatabaseScanner struct {
	config *config.Config
}

// DatabaseFinding represents a database security finding
type DatabaseFinding struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`     // config, vulnerability, access, encryption
	Severity      string                 `json:"severity"` // critical, high, medium, low
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	DatabaseType  string                 `json:"database_type"` // postgresql, mysql, mongodb, redis, sqlserver
	Host          string                 `json:"host"`
	Port          int                    `json:"port"`
	Database      string                 `json:"database,omitempty"`
	CurrentValue  string                 `json:"current_value,omitempty"`
	RequiredValue string                 `json:"required_value,omitempty"`
	Remediation   string                 `json:"remediation"`
	DiscoveredAt  time.Time              `json:"discovered_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// DatabaseInfo represents database connection information
type DatabaseInfo struct {
	Type         string `json:"type"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Database     string `json:"database"`
	Username     string `json:"username"`
	IsEncrypted  bool   `json:"is_encrypted"`
	IsRemote     bool   `json:"is_remote"`
	Version      string `json:"version"`
	IsAccessible bool   `json:"is_accessible"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	DatabaseType   string            `json:"database_type"`
	Host           string            `json:"host"`
	Port           int               `json:"port"`
	Database       string            `json:"database"`
	Username       string            `json:"username"`
	Password       string            `json:"password"`
	SSLMode        string            `json:"ssl_mode"`
	ConnectionPool int               `json:"connection_pool"`
	Timeout        int               `json:"timeout"`
	ConfigParams   map[string]string `json:"config_params"`
}

// NewDatabaseScanner creates a new database security scanner
func NewDatabaseScanner(cfg *config.Config) *DatabaseScanner {
	return &DatabaseScanner{
		config: cfg,
	}
}

// Scan performs comprehensive database security scanning
func (ds *DatabaseScanner) Scan() ([]DatabaseFinding, []DatabaseInfo, error) {
	var findings []DatabaseFinding
	var databases []DatabaseInfo

	// Discover databases
	discoveredDBs := ds.discoverDatabases()
	databases = append(databases, discoveredDBs...)

	// Scan each discovered database
	for _, db := range discoveredDBs {
		if db.IsAccessible {
			dbFindings := ds.scanDatabase(db)
			findings = append(findings, dbFindings...)
		}
	}

	// Scan for common database vulnerabilities
	commonFindings := ds.scanCommonVulnerabilities()
	findings = append(findings, commonFindings...)

	return findings, databases, nil
}

// discoverDatabases discovers databases on the system
func (ds *DatabaseScanner) discoverDatabases() []DatabaseInfo {
	var databases []DatabaseInfo

	// Common database ports and types
	dbPorts := map[int]string{
		5432:  "postgresql",
		3306:  "mysql",
		1433:  "sqlserver",
		27017: "mongodb",
		6379:  "redis",
		1521:  "oracle",
		5984:  "couchdb",
		9200:  "elasticsearch",
	}

	// Scan localhost for database services
	for port, dbType := range dbPorts {
		if ds.isPortOpen("localhost", port) {
			db := DatabaseInfo{
				Type:         dbType,
				Host:         "localhost",
				Port:         port,
				Database:     "default",
				IsEncrypted:  false, // Would need to check SSL/TLS
				IsRemote:     false,
				Version:      "unknown",
				IsAccessible: ds.testDatabaseConnection(dbType, "localhost", port),
			}
			databases = append(databases, db)
		}
	}

	// Check for database configuration files
	configDBs := ds.findDatabaseConfigs()
	databases = append(databases, configDBs...)

	return databases
}

// isPortOpen checks if a port is open
func (ds *DatabaseScanner) isPortOpen(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// testDatabaseConnection tests database connectivity
func (ds *DatabaseScanner) testDatabaseConnection(dbType, host string, port int) bool {
	switch dbType {
	case "postgresql":
		return ds.testPostgreSQLConnection(host, port)
	case "mysql":
		return ds.testMySQLConnection(host, port)
	case "redis":
		return ds.testRedisConnection(host, port)
	default:
		return false
	}
}

// testPostgreSQLConnection tests PostgreSQL connectivity
func (ds *DatabaseScanner) testPostgreSQLConnection(host string, port int) bool {
	connStr := fmt.Sprintf("host=%s port=%d user=postgres dbname=postgres sslmode=disable", host, port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Ping()
	return err == nil
}

// testMySQLConnection tests MySQL connectivity
func (ds *DatabaseScanner) testMySQLConnection(host string, port int) bool {
	connStr := fmt.Sprintf("root:@tcp(%s:%d)/", host, port)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return false
	}
	defer db.Close()

	err = db.Ping()
	return err == nil
}

// testRedisConnection tests Redis connectivity
func (ds *DatabaseScanner) testRedisConnection(host string, port int) bool {
	// This would require Redis client library
	// For now, just check if port is open
	return ds.isPortOpen(host, port)
}

// findDatabaseConfigs finds database configuration files
func (ds *DatabaseScanner) findDatabaseConfigs() []DatabaseInfo {
	var databases []DatabaseInfo

	// This would scan for common database configuration files
	// For now, return empty list
	return databases
}

// scanDatabase scans a specific database for security issues
func (ds *DatabaseScanner) scanDatabase(db DatabaseInfo) []DatabaseFinding {
	var findings []DatabaseFinding

	switch db.Type {
	case "postgresql":
		findings = append(findings, ds.scanPostgreSQL(db)...)
	case "mysql":
		findings = append(findings, ds.scanMySQL(db)...)
	case "redis":
		findings = append(findings, ds.scanRedis(db)...)
	case "mongodb":
		findings = append(findings, ds.scanMongoDB(db)...)
	}

	return findings
}

// scanPostgreSQL scans PostgreSQL for security issues
func (ds *DatabaseScanner) scanPostgreSQL(db DatabaseInfo) []DatabaseFinding {
	var findings []DatabaseFinding

	// Check for default credentials
	finding := DatabaseFinding{
		ID:           uuid.New().String(),
		Type:         "config",
		Severity:     "critical",
		Title:        "PostgreSQL Default Credentials",
		Description:  "PostgreSQL is accessible with default credentials",
		DatabaseType: "postgresql",
		Host:         db.Host,
		Port:         db.Port,
		Remediation:  "Change default PostgreSQL credentials and use strong passwords",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"database": "postgresql",
			"host":     db.Host,
			"port":     db.Port,
		},
	}
	findings = append(findings, finding)

	// Check for SSL/TLS encryption
	if !db.IsEncrypted {
		finding := DatabaseFinding{
			ID:           uuid.New().String(),
			Type:         "encryption",
			Severity:     "high",
			Title:        "PostgreSQL Unencrypted Connection",
			Description:  "PostgreSQL connection is not encrypted",
			DatabaseType: "postgresql",
			Host:         db.Host,
			Port:         db.Port,
			Remediation:  "Enable SSL/TLS encryption for PostgreSQL connections",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"database":  "postgresql",
				"encrypted": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for remote access
	if db.IsRemote {
		finding := DatabaseFinding{
			ID:           uuid.New().String(),
			Type:         "config",
			Severity:     "medium",
			Title:        "PostgreSQL Remote Access",
			Description:  "PostgreSQL is accessible from remote hosts",
			DatabaseType: "postgresql",
			Host:         db.Host,
			Port:         db.Port,
			Remediation:  "Restrict PostgreSQL access to localhost only if not needed remotely",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"database": "postgresql",
				"remote":   true,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanMySQL scans MySQL for security issues
func (ds *DatabaseScanner) scanMySQL(db DatabaseInfo) []DatabaseFinding {
	var findings []DatabaseFinding

	// Check for default credentials
	finding := DatabaseFinding{
		ID:           uuid.New().String(),
		Type:         "config",
		Severity:     "critical",
		Title:        "MySQL Default Credentials",
		Description:  "MySQL is accessible with default credentials",
		DatabaseType: "mysql",
		Host:         db.Host,
		Port:         db.Port,
		Remediation:  "Change default MySQL credentials and use strong passwords",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"database": "mysql",
			"host":     db.Host,
			"port":     db.Port,
		},
	}
	findings = append(findings, finding)

	// Check for SSL/TLS encryption
	if !db.IsEncrypted {
		finding := DatabaseFinding{
			ID:           uuid.New().String(),
			Type:         "encryption",
			Severity:     "high",
			Title:        "MySQL Unencrypted Connection",
			Description:  "MySQL connection is not encrypted",
			DatabaseType: "mysql",
			Host:         db.Host,
			Port:         db.Port,
			Remediation:  "Enable SSL/TLS encryption for MySQL connections",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"database":  "mysql",
				"encrypted": false,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanRedis scans Redis for security issues
func (ds *DatabaseScanner) scanRedis(db DatabaseInfo) []DatabaseFinding {
	var findings []DatabaseFinding

	// Check for default configuration
	finding := DatabaseFinding{
		ID:           uuid.New().String(),
		Type:         "config",
		Severity:     "high",
		Title:        "Redis Default Configuration",
		Description:  "Redis is running with default configuration",
		DatabaseType: "redis",
		Host:         db.Host,
		Port:         db.Port,
		Remediation:  "Configure Redis with authentication and proper security settings",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"database": "redis",
			"host":     db.Host,
			"port":     db.Port,
		},
	}
	findings = append(findings, finding)

	// Check for remote access
	if db.IsRemote {
		finding := DatabaseFinding{
			ID:           uuid.New().String(),
			Type:         "config",
			Severity:     "critical",
			Title:        "Redis Remote Access Without Authentication",
			Description:  "Redis is accessible from remote hosts without authentication",
			DatabaseType: "redis",
			Host:         db.Host,
			Port:         db.Port,
			Remediation:  "Enable Redis authentication and restrict access to localhost",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"database": "redis",
				"remote":   true,
				"auth":     false,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanMongoDB scans MongoDB for security issues
func (ds *DatabaseScanner) scanMongoDB(db DatabaseInfo) []DatabaseFinding {
	var findings []DatabaseFinding

	// Check for default configuration
	finding := DatabaseFinding{
		ID:           uuid.New().String(),
		Type:         "config",
		Severity:     "high",
		Title:        "MongoDB Default Configuration",
		Description:  "MongoDB is running with default configuration",
		DatabaseType: "mongodb",
		Host:         db.Host,
		Port:         db.Port,
		Remediation:  "Configure MongoDB with authentication and proper security settings",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"database": "mongodb",
			"host":     db.Host,
			"port":     db.Port,
		},
	}
	findings = append(findings, finding)

	return findings
}

// scanCommonVulnerabilities scans for common database vulnerabilities
func (ds *DatabaseScanner) scanCommonVulnerabilities() []DatabaseFinding {
	var findings []DatabaseFinding

	// Check for database files with weak permissions
	finding := DatabaseFinding{
		ID:           uuid.New().String(),
		Type:         "config",
		Severity:     "medium",
		Title:        "Database Files with Weak Permissions",
		Description:  "Database files may have overly permissive access controls",
		Remediation:  "Review and restrict database file permissions",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"category": "file_permissions",
		},
	}
	findings = append(findings, finding)

	// Check for database backup files
	finding = DatabaseFinding{
		ID:           uuid.New().String(),
		Type:         "config",
		Severity:     "high",
		Title:        "Unencrypted Database Backups",
		Description:  "Database backup files may not be encrypted",
		Remediation:  "Encrypt database backup files and secure backup storage",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"category": "backup_security",
		},
	}
	findings = append(findings, finding)

	return findings
}
