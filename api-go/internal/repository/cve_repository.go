package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// CVERepository handles CVE data operations
type CVERepository struct {
	pool *pgxpool.Pool
}

// NewCVERepository creates a new CVE repository
func NewCVERepository(pool *pgxpool.Pool) *CVERepository {
	return &CVERepository{
		pool: pool,
	}
}

// SearchCVEs searches for CVEs by software name, version, or CPE
func (r *CVERepository) SearchCVEs(
	ctx context.Context,
	softwareName string,
	version *string,
	cpeString *string,
	limit int,
) ([]map[string]interface{}, error) {
	query := `
		SELECT DISTINCT cd.cve_id, cd.cve_data, cd.severity, cd.cvss_score
		FROM cve_data cd
		LEFT JOIN cve_cpe_mapping ccm ON cd.cve_id = ccm.cve_id
		WHERE (
			($1::text IS NOT NULL AND ccm.cpe_string = $1) OR
			($2::text IS NOT NULL AND (
				LOWER(ccm.vendor) LIKE '%' || LOWER($2) || '%' OR
				LOWER(ccm.product) LIKE '%' || LOWER($2) || '%'
			)) OR
			($2::text IS NOT NULL AND cd.search_vector @@ plainto_tsquery('english', $2))
		)
		ORDER BY cd.cvss_score DESC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, cpeString, softwareName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query CVEs: %w", err)
	}
	defer rows.Close()

	var cves []map[string]interface{}
	for rows.Next() {
		var cveID string
		var cveDataJSON []byte
		var severity sql.NullString
		var cvssScore sql.NullFloat64

		err := rows.Scan(&cveID, &cveDataJSON, &severity, &cvssScore)
		if err != nil {
			continue
		}

		var cveData map[string]interface{}
		if err := json.Unmarshal(cveDataJSON, &cveData); err != nil {
			continue
		}

		cve := make(map[string]interface{})
		cve["id"] = cveID
		cve["cve_id"] = cveID
		for k, v := range cveData {
			cve[k] = v
		}
		if severity.Valid {
			cve["severity"] = severity.String
		}
		if cvssScore.Valid {
			cve["cvss_score"] = cvssScore.Float64
		}

		cves = append(cves, cve)
	}

	return cves, nil
}

// GetCVEByID retrieves a specific CVE by ID
func (r *CVERepository) GetCVEByID(ctx context.Context, cveID string) (map[string]interface{}, error) {
	query := `
		SELECT cve_id, cve_data, severity, cvss_score
		FROM cve_data
		WHERE cve_id = $1
	`

	var cveDataJSON []byte
	var severity sql.NullString
	var cvssScore sql.NullFloat64

	err := r.pool.QueryRow(ctx, query, cveID).Scan(&cveID, &cveDataJSON, &severity, &cvssScore)
	if err != nil {
		return nil, fmt.Errorf("failed to get CVE: %w", err)
	}

	var cveData map[string]interface{}
	if err := json.Unmarshal(cveDataJSON, &cveData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CVE data: %w", err)
	}

	cve := make(map[string]interface{})
	cve["id"] = cveID
	cve["cve_id"] = cveID
	for k, v := range cveData {
		cve[k] = v
	}
	if severity.Valid {
		cve["severity"] = severity.String
	}
	if cvssScore.Valid {
		cve["cvss_score"] = cvssScore.Float64
	}

	return cve, nil
}

// GetCVECount returns the total number of CVEs in the database
func (r *CVERepository) GetCVECount(ctx context.Context) (int64, error) {
	var count int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM cve_data").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get CVE count: %w", err)
	}
	return count, nil
}

