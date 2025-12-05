-- Migration 004: CVE Database Storage
-- Migrates CVE data from JSON file to PostgreSQL with JSONB storage
-- Supports incremental updates and fast searches with GIN indexes

-- Create CVE data table with JSONB storage
CREATE TABLE IF NOT EXISTS cve_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cve_id VARCHAR(50) NOT NULL UNIQUE,
    cve_data JSONB NOT NULL,
    -- Extracted fields for faster queries (indexed)
    description TEXT GENERATED ALWAYS AS (cve_data->>'description') STORED,
    severity VARCHAR(20) GENERATED ALWAYS AS (cve_data->>'severity') STORED,
    cvss_score NUMERIC(3,1) GENERATED ALWAYS AS ((cve_data->>'cvss_score')::NUMERIC) STORED,
    published_date TIMESTAMP WITH TIME ZONE GENERATED ALWAYS AS ((cve_data->>'published_date')::TIMESTAMP WITH TIME ZONE) STORED,
    last_modified_date TIMESTAMP WITH TIME ZONE GENERATED ALWAYS AS ((cve_data->>'last_modified_date')::TIMESTAMP WITH TIME ZONE) STORED,
    -- Full-text search vector
    search_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('english', 
            COALESCE(cve_data->>'description', '') || ' ' ||
            COALESCE(cve_data->>'title', '') || ' ' ||
            COALESCE(cve_data->>'cve_id', '')
        )
    ) STORED,
    -- Metadata
    source VARCHAR(50) DEFAULT 'nvd',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_cve_data_cve_id ON cve_data(cve_id);
CREATE INDEX IF NOT EXISTS idx_cve_data_severity ON cve_data(severity);
CREATE INDEX IF NOT EXISTS idx_cve_data_cvss_score ON cve_data(cvss_score DESC);
CREATE INDEX IF NOT EXISTS idx_cve_data_published_date ON cve_data(published_date DESC);
CREATE INDEX IF NOT EXISTS idx_cve_data_last_modified ON cve_data(last_modified_date DESC);
CREATE INDEX IF NOT EXISTS idx_cve_data_source ON cve_data(source);

-- GIN index for JSONB queries
CREATE INDEX IF NOT EXISTS idx_cve_data_jsonb_gin ON cve_data USING GIN (cve_data);

-- GIN index for full-text search
CREATE INDEX IF NOT EXISTS idx_cve_data_search_vector ON cve_data USING GIN (search_vector);

-- Index for CPE matching (common query pattern)
CREATE INDEX IF NOT EXISTS idx_cve_data_cpe ON cve_data USING GIN ((cve_data->'cpe_list'));

-- Create CPE mapping table for fast CPE-to-CVE lookups
CREATE TABLE IF NOT EXISTS cve_cpe_mapping (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cve_id VARCHAR(50) NOT NULL REFERENCES cve_data(cve_id) ON DELETE CASCADE,
    cpe_string TEXT NOT NULL,
    cpe_part VARCHAR(10), -- 'a' for application, 'o' for operating system, 'h' for hardware
    vendor TEXT,
    product TEXT,
    version TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(cve_id, cpe_string)
);

-- Indexes for CPE mapping
CREATE INDEX IF NOT EXISTS idx_cve_cpe_mapping_cve_id ON cve_cpe_mapping(cve_id);
CREATE INDEX IF NOT EXISTS idx_cve_cpe_mapping_cpe_string ON cve_cpe_mapping(cpe_string);
CREATE INDEX IF NOT EXISTS idx_cve_cpe_mapping_vendor_product ON cve_cpe_mapping(vendor, product);
CREATE INDEX IF NOT EXISTS idx_cve_cpe_mapping_part ON cve_cpe_mapping(cpe_part);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_cve_data_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for updated_at
DROP TRIGGER IF EXISTS trigger_update_cve_data_updated_at ON cve_data;
CREATE TRIGGER trigger_update_cve_data_updated_at
    BEFORE UPDATE ON cve_data
    FOR EACH ROW
    EXECUTE FUNCTION update_cve_data_updated_at();

-- Create materialized view for CVE statistics (refreshed periodically)
CREATE MATERIALIZED VIEW IF NOT EXISTS cve_statistics AS
SELECT 
    severity,
    COUNT(*) as count,
    AVG(cvss_score) as avg_cvss,
    MIN(published_date) as earliest_cve,
    MAX(published_date) as latest_cve,
    COUNT(*) FILTER (WHERE published_date > NOW() - INTERVAL '30 days') as recent_count
FROM cve_data
GROUP BY severity;

CREATE UNIQUE INDEX IF NOT EXISTS idx_cve_statistics_severity ON cve_statistics(severity);

-- Create function to refresh statistics
CREATE OR REPLACE FUNCTION refresh_cve_statistics()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY cve_statistics;
END;
$$ LANGUAGE plpgsql;

-- Add comment
COMMENT ON TABLE cve_data IS 'CVE data stored as JSONB for flexibility with extracted fields for performance';
COMMENT ON TABLE cve_cpe_mapping IS 'Mapping table for fast CPE-to-CVE lookups';
COMMENT ON MATERIALIZED VIEW cve_statistics IS 'Cached CVE statistics for dashboard queries';

