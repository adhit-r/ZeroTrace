-- Migration 005: Performance Optimization
-- Adds additional indexes, query optimization, and connection pool settings

BEGIN;

-- ============================================================================
-- ADDITIONAL PERFORMANCE INDEXES
-- ============================================================================

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_company_severity_created 
    ON vulnerabilities(company_id, severity, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_vulnerabilities_company_status_created 
    ON vulnerabilities(company_id, status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_scans_company_status_created 
    ON scans(company_id, status, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_agents_company_status_last_seen 
    ON agents(company_id, status, last_seen DESC);

-- Partial indexes for active records (most common queries)
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_active 
    ON vulnerabilities(company_id, severity) 
    WHERE status = 'open';

CREATE INDEX IF NOT EXISTS idx_agents_active 
    ON agents(company_id, last_seen DESC) 
    WHERE status = 'active';

CREATE INDEX IF NOT EXISTS idx_scans_running 
    ON scans(company_id, created_at DESC) 
    WHERE status IN ('pending', 'running');

-- Indexes for organization_id (multi-tenant)
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_organization_id 
    ON vulnerabilities(organization_id) 
    WHERE organization_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_agents_organization_id 
    ON agents(organization_id) 
    WHERE organization_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_scans_organization_id 
    ON scans(organization_id) 
    WHERE organization_id IS NOT NULL;

-- Indexes for time-based queries
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_discovered_at 
    ON vulnerabilities(discovered_at DESC);

CREATE INDEX IF NOT EXISTS idx_vulnerabilities_last_seen 
    ON vulnerabilities(last_seen DESC);

-- Indexes for JSONB fields (common queries)
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_metadata_gin 
    ON vulnerabilities USING GIN (metadata);

CREATE INDEX IF NOT EXISTS idx_vulnerabilities_enrichment_data_gin 
    ON vulnerabilities USING GIN (enrichment_data);

CREATE INDEX IF NOT EXISTS idx_scans_options_gin 
    ON scans USING GIN (options);

CREATE INDEX IF NOT EXISTS idx_scans_results_gin 
    ON scans USING GIN (results);

-- ============================================================================
-- QUERY OPTIMIZATION
-- ============================================================================

-- Analyze tables for query planner
ANALYZE vulnerabilities;
ANALYZE scans;
ANALYZE agents;
ANALYZE companies;

-- Update statistics
VACUUM ANALYZE;

-- ============================================================================
-- CONNECTION POOL SETTINGS (via ALTER DATABASE - run separately)
-- ============================================================================

-- Note: These should be set at database level, not in migration
-- ALTER DATABASE zerotrace SET shared_buffers = '256MB';
-- ALTER DATABASE zerotrace SET effective_cache_size = '1GB';
-- ALTER DATABASE zerotrace SET maintenance_work_mem = '64MB';
-- ALTER DATABASE zerotrace SET checkpoint_completion_target = 0.9;
-- ALTER DATABASE zerotrace SET wal_buffers = '16MB';
-- ALTER DATABASE zerotrace SET default_statistics_target = 100;
-- ALTER DATABASE zerotrace SET random_page_cost = 1.1;
-- ALTER DATABASE zerotrace SET effective_io_concurrency = 200;
-- ALTER DATABASE zerotrace SET work_mem = '4MB';
-- ALTER DATABASE zerotrace SET min_wal_size = '1GB';
-- ALTER DATABASE zerotrace SET max_wal_size = '4GB';

COMMIT;

