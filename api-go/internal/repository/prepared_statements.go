package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PreparedStatementManager manages prepared statements for performance
type PreparedStatementManager struct {
	pool *pgxpool.Pool
}

// NewPreparedStatementManager creates a new prepared statement manager
func NewPreparedStatementManager(pool *pgxpool.Pool) *PreparedStatementManager {
	return &PreparedStatementManager{
		pool: pool,
	}
}

// PrepareStatement prepares a SQL statement for reuse
func (psm *PreparedStatementManager) PrepareStatement(ctx context.Context, name, query string) error {
	// pgx/v5 automatically uses prepared statements for repeated queries
	// This is a placeholder for explicit statement management if needed
	_, err := psm.pool.Exec(ctx, fmt.Sprintf("PREPARE %s AS %s", name, query))
	return err
}

// ExecutePrepared executes a prepared statement
func (psm *PreparedStatementManager) ExecutePrepared(ctx context.Context, name string, args ...interface{}) error {
	// Execute prepared statement
	_, err := psm.pool.Exec(ctx, fmt.Sprintf("EXECUTE %s", name), args...)
	return err
}

// Common prepared statements for frequently used queries
const (
	StmtGetVulnerabilitiesByCompany = "get_vulnerabilities_by_company"
	StmtGetScansByCompany           = "get_scans_by_company"
	StmtGetAgentsByCompany          = "get_agents_by_company"
)

// InitializePreparedStatements initializes common prepared statements
func (psm *PreparedStatementManager) InitializePreparedStatements(ctx context.Context) error {
	statements := map[string]string{
		StmtGetVulnerabilitiesByCompany: `
			SELECT * FROM vulnerabilities 
			WHERE company_id = $1 AND status = $2 
			ORDER BY created_at DESC 
			LIMIT $3 OFFSET $4
		`,
		StmtGetScansByCompany: `
			SELECT * FROM scans 
			WHERE company_id = $1 AND status = $2 
			ORDER BY created_at DESC 
			LIMIT $3 OFFSET $4
		`,
		StmtGetAgentsByCompany: `
			SELECT * FROM agents 
			WHERE company_id = $1 AND status = $2 
			ORDER BY last_seen DESC
		`,
	}

	for name, query := range statements {
		if err := psm.PrepareStatement(ctx, name, query); err != nil {
			return fmt.Errorf("failed to prepare statement %s: %w", name, err)
		}
	}

	return nil
}

