package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

type Migration struct {
	Version string
	File    string
	Content string
}

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "zerotrace")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create migrations table if not exists
	if err := createMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get all migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		log.Fatalf("Failed to get migration files: %v", err)
	}

	// Get applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		log.Fatalf("Failed to get applied migrations: %v", err)
	}

	// Apply new migrations
	for _, migration := range migrations {
		if !isMigrationApplied(migration.Version, appliedMigrations) {
			fmt.Printf("Applying migration: %s\n", migration.File)

			if err := applyMigration(db, migration); err != nil {
				log.Fatalf("Failed to apply migration %s: %v", migration.File, err)
			}

			fmt.Printf("Successfully applied migration: %s\n", migration.File)
		} else {
			fmt.Printf("Migration already applied: %s\n", migration.File)
		}
	}

	fmt.Println("All migrations completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationFiles() ([]Migration, error) {
	migrationsDir := "."
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") && strings.Contains(file.Name(), "_") {
			// Extract version from filename (e.g., "001_initial_schema.sql" -> "001")
			parts := strings.Split(file.Name(), "_")
			if len(parts) >= 2 {
				version := parts[0]

				content, err := ioutil.ReadFile(filepath.Join(migrationsDir, file.Name()))
				if err != nil {
					return nil, fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
				}

				migrations = append(migrations, Migration{
					Version: version,
					File:    file.Name(),
					Content: string(content),
				})
			}
		}
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) ([]string, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applied []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied = append(applied, version)
	}

	return applied, nil
}

func isMigrationApplied(version string, appliedMigrations []string) bool {
	for _, applied := range appliedMigrations {
		if applied == version {
			return true
		}
	}
	return false
}

func applyMigration(db *sql.DB, migration Migration) error {
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(migration.Content); err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}

	// Record migration as applied
	query := "INSERT INTO schema_migrations (version) VALUES ($1)"
	if _, err := tx.Exec(query, migration.Version); err != nil {
		return fmt.Errorf("failed to record migration: %v", err)
	}

	// Commit transaction
	return tx.Commit()
}
