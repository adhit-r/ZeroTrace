# ZeroTrace Database Schema

## Overview
The ZeroTrace database uses PostgreSQL 15+ with a multi-tenant architecture designed to handle 100,000+ data points efficiently. The schema includes tables for companies, users, scans, vulnerabilities, agents, and enrichment data.

## Database Architecture

### Multi-Tenant Design
- **Company-based isolation**: All data is partitioned by company_id
- **Row-level security**: Database-level access control
- **Partitioning**: Large tables partitioned by company_id and date
- **Indexing**: Optimized indexes for query performance

### Schema Version
```sql
-- Check schema version
SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;
```

## Core Tables

### 1. Companies Table
```sql
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    settings JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_companies_domain ON companies(domain);
CREATE INDEX idx_companies_status ON companies(status);
CREATE INDEX idx_companies_created_at ON companies(created_at);

-- Row-level security
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
CREATE POLICY companies_isolation ON companies
    FOR ALL USING (id = current_setting('app.current_company_id')::UUID);
```

### 2. Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'USER',
    status VARCHAR(50) DEFAULT 'active',
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(company_id, email)
);

-- Indexes
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_last_login ON users(last_login);

-- Row-level security
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_isolation ON users
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

### 3. Agents Table
```sql
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    api_key VARCHAR(255) UNIQUE NOT NULL,
    capabilities JSONB DEFAULT '[]',
    location VARCHAR(100),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    current_scan_id UUID,
    performance_metrics JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_agents_company_id ON agents(company_id);
CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_api_key ON agents(api_key);
CREATE INDEX idx_agents_last_seen ON agents(last_seen);
CREATE INDEX idx_agents_current_scan ON agents(current_scan_id);

-- Row-level security
ALTER TABLE agents ENABLE ROW LEVEL SECURITY;
CREATE POLICY agents_isolation ON agents
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

### 4. Scans Table
```sql
CREATE TABLE scans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
    repository VARCHAR(500) NOT NULL,
    branch VARCHAR(100) DEFAULT 'main',
    commit_hash VARCHAR(100),
    scan_type VARCHAR(50) DEFAULT 'full',
    status VARCHAR(50) DEFAULT 'pending',
    progress INTEGER DEFAULT 0,
    start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    options JSONB DEFAULT '{}',
    results JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partition by company_id and date
CREATE TABLE scans_partitioned (
    LIKE scans INCLUDING ALL
) PARTITION BY HASH (company_id);

-- Create partitions
CREATE TABLE scans_partition_0 PARTITION OF scans_partitioned
    FOR VALUES WITH (modulus 4, remainder 0);
CREATE TABLE scans_partition_1 PARTITION OF scans_partitioned
    FOR VALUES WITH (modulus 4, remainder 1);
CREATE TABLE scans_partition_2 PARTITION OF scans_partitioned
    FOR VALUES WITH (modulus 4, remainder 2);
CREATE TABLE scans_partition_3 PARTITION OF scans_partitioned
    FOR VALUES WITH (modulus 4, remainder 3);

-- Indexes
CREATE INDEX idx_scans_company_id ON scans_partitioned(company_id);
CREATE INDEX idx_scans_agent_id ON scans_partitioned(agent_id);
CREATE INDEX idx_scans_status ON scans_partitioned(status);
CREATE INDEX idx_scans_repository ON scans_partitioned(repository);
CREATE INDEX idx_scans_created_at ON scans_partitioned(created_at);
CREATE INDEX idx_scans_start_time ON scans_partitioned(start_time);
CREATE INDEX idx_scans_company_status ON scans_partitioned(company_id, status);

-- Row-level security
ALTER TABLE scans_partitioned ENABLE ROW LEVEL SECURITY;
CREATE POLICY scans_isolation ON scans_partitioned
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

### 5. Vulnerabilities Table
```sql
CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scan_id UUID NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    cve_id VARCHAR(20),
    cvss_score DECIMAL(3,1),
    cvss_vector VARCHAR(100),
    package_name VARCHAR(255),
    package_version VARCHAR(100),
    location VARCHAR(500),
    remediation TEXT,
    references JSONB DEFAULT '[]',
    affected_versions JSONB DEFAULT '[]',
    patched_versions JSONB DEFAULT '[]',
    exploit_available BOOLEAN DEFAULT FALSE,
    exploit_count INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'open',
    priority VARCHAR(20) DEFAULT 'medium',
    notes TEXT,
    enrichment_data JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partition by company_id
CREATE TABLE vulnerabilities_partitioned (
    LIKE vulnerabilities INCLUDING ALL
) PARTITION BY HASH (company_id);

-- Create partitions
CREATE TABLE vulnerabilities_partition_0 PARTITION OF vulnerabilities_partitioned
    FOR VALUES WITH (modulus 4, remainder 0);
CREATE TABLE vulnerabilities_partition_1 PARTITION OF vulnerabilities_partitioned
    FOR VALUES WITH (modulus 4, remainder 1);
CREATE TABLE vulnerabilities_partition_2 PARTITION OF vulnerabilities_partitioned
    FOR VALUES WITH (modulus 4, remainder 2);
CREATE TABLE vulnerabilities_partition_3 PARTITION OF vulnerabilities_partitioned
    FOR VALUES WITH (modulus 4, remainder 3);

-- Indexes
CREATE INDEX idx_vulnerabilities_scan_id ON vulnerabilities_partitioned(scan_id);
CREATE INDEX idx_vulnerabilities_company_id ON vulnerabilities_partitioned(company_id);
CREATE INDEX idx_vulnerabilities_severity ON vulnerabilities_partitioned(severity);
CREATE INDEX idx_vulnerabilities_type ON vulnerabilities_partitioned(type);
CREATE INDEX idx_vulnerabilities_cve_id ON vulnerabilities_partitioned(cve_id);
CREATE INDEX idx_vulnerabilities_package ON vulnerabilities_partitioned(package_name, package_version);
CREATE INDEX idx_vulnerabilities_status ON vulnerabilities_partitioned(status);
CREATE INDEX idx_vulnerabilities_created_at ON vulnerabilities_partitioned(created_at);
CREATE INDEX idx_vulnerabilities_company_severity ON vulnerabilities_partitioned(company_id, severity);
CREATE INDEX idx_vulnerabilities_company_status ON vulnerabilities_partitioned(company_id, status);

-- Full-text search
CREATE INDEX idx_vulnerabilities_search ON vulnerabilities_partitioned 
    USING gin(to_tsvector('english', title || ' ' || description));

-- Row-level security
ALTER TABLE vulnerabilities_partitioned ENABLE ROW LEVEL SECURITY;
CREATE POLICY vulnerabilities_isolation ON vulnerabilities_partitioned
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

### 6. Dependencies Table
```sql
CREATE TABLE dependencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scan_id UUID NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    version VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(500),
    license VARCHAR(100),
    vulnerabilities JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partition by company_id
CREATE TABLE dependencies_partitioned (
    LIKE dependencies INCLUDING ALL
) PARTITION BY HASH (company_id);

-- Create partitions
CREATE TABLE dependencies_partition_0 PARTITION OF dependencies_partitioned
    FOR VALUES WITH (modulus 4, remainder 0);
CREATE TABLE dependencies_partition_1 PARTITION OF dependencies_partitioned
    FOR VALUES WITH (modulus 4, remainder 1);
CREATE TABLE dependencies_partition_2 PARTITION OF dependencies_partitioned
    FOR VALUES WITH (modulus 4, remainder 2);
CREATE TABLE dependencies_partition_3 PARTITION OF dependencies_partitioned
    FOR VALUES WITH (modulus 4, remainder 3);

-- Indexes
CREATE INDEX idx_dependencies_scan_id ON dependencies_partitioned(scan_id);
CREATE INDEX idx_dependencies_company_id ON dependencies_partitioned(company_id);
CREATE INDEX idx_dependencies_name_version ON dependencies_partitioned(name, version);
CREATE INDEX idx_dependencies_type ON dependencies_partitioned(type);

-- Row-level security
ALTER TABLE dependencies_partitioned ENABLE ROW LEVEL SECURITY;
CREATE POLICY dependencies_isolation ON dependencies_partitioned
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

## Enrichment Tables

### 7. CVE Database Table
```sql
CREATE TABLE cve_database (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cve_id VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    cvss_score DECIMAL(3,1),
    cvss_vector VARCHAR(100),
    severity VARCHAR(20),
    published_date DATE,
    last_modified_date DATE,
    references JSONB DEFAULT '[]',
    affected_products JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_cve_database_cve_id ON cve_database(cve_id);
CREATE INDEX idx_cve_database_severity ON cve_database(severity);
CREATE INDEX idx_cve_database_cvss_score ON cve_database(cvss_score);
CREATE INDEX idx_cve_database_published_date ON cve_database(published_date);

-- Full-text search
CREATE INDEX idx_cve_database_search ON cve_database 
    USING gin(to_tsvector('english', title || ' ' || description));
```

### 8. Enrichment Results Table
```sql
CREATE TABLE enrichment_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vulnerability_id UUID NOT NULL REFERENCES vulnerabilities(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    enrichment_type VARCHAR(50) NOT NULL,
    risk_score DECIMAL(3,2),
    trend VARCHAR(20),
    recommendations JSONB DEFAULT '[]',
    threat_intelligence JSONB DEFAULT '{}',
    ml_predictions JSONB DEFAULT '{}',
    confidence_score DECIMAL(3,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partition by company_id
CREATE TABLE enrichment_results_partitioned (
    LIKE enrichment_results INCLUDING ALL
) PARTITION BY HASH (company_id);

-- Create partitions
CREATE TABLE enrichment_results_partition_0 PARTITION OF enrichment_results_partitioned
    FOR VALUES WITH (modulus 4, remainder 0);
CREATE TABLE enrichment_results_partition_1 PARTITION OF enrichment_results_partitioned
    FOR VALUES WITH (modulus 4, remainder 1);
CREATE TABLE enrichment_results_partition_2 PARTITION OF enrichment_results_partitioned
    FOR VALUES WITH (modulus 4, remainder 2);
CREATE TABLE enrichment_results_partition_3 PARTITION OF enrichment_results_partitioned
    FOR VALUES WITH (modulus 4, remainder 3);

-- Indexes
CREATE INDEX idx_enrichment_vulnerability_id ON enrichment_results_partitioned(vulnerability_id);
CREATE INDEX idx_enrichment_company_id ON enrichment_results_partitioned(company_id);
CREATE INDEX idx_enrichment_type ON enrichment_results_partitioned(enrichment_type);
CREATE INDEX idx_enrichment_risk_score ON enrichment_results_partitioned(risk_score);

-- Row-level security
ALTER TABLE enrichment_results_partitioned ENABLE ROW LEVEL SECURITY;
CREATE POLICY enrichment_isolation ON enrichment_results_partitioned
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

## Reporting Tables

### 9. Reports Table
```sql
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    scan_id UUID REFERENCES scans(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    format VARCHAR(20) NOT NULL,
    status VARCHAR(50) DEFAULT 'generating',
    file_path VARCHAR(500),
    file_size BIGINT,
    download_url VARCHAR(500),
    options JSONB DEFAULT '{}',
    generated_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_reports_company_id ON reports(company_id);
CREATE INDEX idx_reports_scan_id ON reports(scan_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_created_at ON reports(created_at);

-- Row-level security
ALTER TABLE reports ENABLE ROW LEVEL SECURITY;
CREATE POLICY reports_isolation ON reports
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

## Audit and Logging Tables

### 10. Audit Logs Table
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partition by company_id and date
CREATE TABLE audit_logs_partitioned (
    LIKE audit_logs INCLUDING ALL
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE audit_logs_2024_01 PARTITION OF audit_logs_partitioned
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
CREATE TABLE audit_logs_2024_02 PARTITION OF audit_logs_partitioned
    FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

-- Indexes
CREATE INDEX idx_audit_logs_company_id ON audit_logs_partitioned(company_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs_partitioned(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs_partitioned(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs_partitioned(created_at);
CREATE INDEX idx_audit_logs_resource ON audit_logs_partitioned(resource_type, resource_id);

-- Row-level security
ALTER TABLE audit_logs_partitioned ENABLE ROW LEVEL SECURITY;
CREATE POLICY audit_logs_isolation ON audit_logs_partitioned
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

### 11. System Logs Table
```sql
CREATE TABLE system_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level VARCHAR(20) NOT NULL,
    service VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Partition by date
CREATE TABLE system_logs_partitioned (
    LIKE system_logs INCLUDING ALL
) PARTITION BY RANGE (created_at);

-- Create daily partitions
CREATE TABLE system_logs_2024_01_15 PARTITION OF system_logs_partitioned
    FOR VALUES FROM ('2024-01-15') TO ('2024-01-16');

-- Indexes
CREATE INDEX idx_system_logs_level ON system_logs_partitioned(level);
CREATE INDEX idx_system_logs_service ON system_logs_partitioned(service);
CREATE INDEX idx_system_logs_created_at ON system_logs_partitioned(created_at);
```

## Views

### 1. Scan Summary View
```sql
CREATE VIEW scan_summary AS
SELECT 
    s.id,
    s.company_id,
    s.repository,
    s.branch,
    s.status,
    s.progress,
    s.start_time,
    s.end_time,
    s.created_at,
    COUNT(v.id) as vulnerability_count,
    COUNT(CASE WHEN v.severity = 'critical' THEN 1 END) as critical_count,
    COUNT(CASE WHEN v.severity = 'high' THEN 1 END) as high_count,
    COUNT(CASE WHEN v.severity = 'medium' THEN 1 END) as medium_count,
    COUNT(CASE WHEN v.severity = 'low' THEN 1 END) as low_count
FROM scans_partitioned s
LEFT JOIN vulnerabilities_partitioned v ON s.id = v.scan_id
GROUP BY s.id, s.company_id, s.repository, s.branch, s.status, s.progress, 
         s.start_time, s.end_time, s.created_at;
```

### 2. Company Statistics View
```sql
CREATE VIEW company_statistics AS
SELECT 
    c.id as company_id,
    c.name as company_name,
    COUNT(DISTINCT s.id) as total_scans,
    COUNT(DISTINCT v.id) as total_vulnerabilities,
    COUNT(CASE WHEN v.severity = 'critical' THEN 1 END) as critical_vulnerabilities,
    COUNT(CASE WHEN v.severity = 'high' THEN 1 END) as high_vulnerabilities,
    COUNT(CASE WHEN v.severity = 'medium' THEN 1 END) as medium_vulnerabilities,
    COUNT(CASE WHEN v.severity = 'low' THEN 1 END) as low_vulnerabilities,
    MAX(s.created_at) as last_scan_date
FROM companies c
LEFT JOIN scans_partitioned s ON c.id = s.company_id
LEFT JOIN vulnerabilities_partitioned v ON s.id = v.scan_id
GROUP BY c.id, c.name;
```

## Functions and Triggers

### 1. Update Timestamp Function
```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables with updated_at
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_agents_updated_at BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_scans_updated_at BEFORE UPDATE ON scans_partitioned
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vulnerabilities_updated_at BEFORE UPDATE ON vulnerabilities_partitioned
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### 2. Company Isolation Function
```sql
CREATE OR REPLACE FUNCTION set_company_context(company_id UUID)
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.current_company_id', company_id::text, false);
END;
$$ LANGUAGE plpgsql;
```

### 3. Vulnerability Statistics Function
```sql
CREATE OR REPLACE FUNCTION get_vulnerability_statistics(
    p_company_id UUID,
    p_date_from DATE DEFAULT NULL,
    p_date_to DATE DEFAULT NULL
)
RETURNS TABLE (
    total_vulnerabilities BIGINT,
    critical_count BIGINT,
    high_count BIGINT,
    medium_count BIGINT,
    low_count BIGINT,
    new_vulnerabilities BIGINT,
    resolved_vulnerabilities BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        COUNT(v.id) as total_vulnerabilities,
        COUNT(CASE WHEN v.severity = 'critical' THEN 1 END) as critical_count,
        COUNT(CASE WHEN v.severity = 'high' THEN 1 END) as high_count,
        COUNT(CASE WHEN v.severity = 'medium' THEN 1 END) as medium_count,
        COUNT(CASE WHEN v.severity = 'low' THEN 1 END) as low_count,
        COUNT(CASE WHEN v.created_at >= COALESCE(p_date_from, CURRENT_DATE - INTERVAL '30 days') THEN 1 END) as new_vulnerabilities,
        COUNT(CASE WHEN v.status = 'resolved' AND v.updated_at >= COALESCE(p_date_from, CURRENT_DATE - INTERVAL '30 days') THEN 1 END) as resolved_vulnerabilities
    FROM vulnerabilities_partitioned v
    WHERE v.company_id = p_company_id
    AND (p_date_from IS NULL OR v.created_at >= p_date_from)
    AND (p_date_to IS NULL OR v.created_at <= p_date_to);
END;
$$ LANGUAGE plpgsql;
```

## Performance Optimizations

### 1. Connection Pooling
```sql
-- Configure connection pooling
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET work_mem = '4MB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
```

### 2. Vacuum and Maintenance
```sql
-- Create maintenance function
CREATE OR REPLACE FUNCTION maintenance_cleanup()
RETURNS VOID AS $$
BEGIN
    -- Vacuum tables
    VACUUM ANALYZE;
    
    -- Update statistics
    ANALYZE;
    
    -- Clean up old logs (keep 90 days)
    DELETE FROM system_logs_partitioned 
    WHERE created_at < CURRENT_DATE - INTERVAL '90 days';
    
    DELETE FROM audit_logs_partitioned 
    WHERE created_at < CURRENT_DATE - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- Schedule maintenance (run daily)
SELECT cron.schedule('maintenance-cleanup', '0 2 * * *', 'SELECT maintenance_cleanup();');
```

### 3. Monitoring Queries
```sql
-- Table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Slow queries
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

## Backup and Recovery

### 1. Backup Strategy
```bash
#!/bin/bash
# backup.sh

# Full backup
pg_dump -h localhost -U postgres -d zerotrace -F c -f /backups/zerotrace_$(date +%Y%m%d_%H%M%S).dump

# WAL archiving
# In postgresql.conf:
# archive_mode = on
# archive_command = 'cp %p /backups/wal/%f'
```

### 2. Point-in-Time Recovery
```sql
-- Create restore point
SELECT pg_create_restore_point('before_migration_2024_01_15');

-- Restore to point
-- pg_restore -h localhost -U postgres -d zerotrace --clean --if-exists backup.dump
```

## Security Considerations

### 1. Row-Level Security
All tables with company_id have RLS enabled to ensure data isolation.

### 2. Encryption
```sql
-- Enable encryption at rest
ALTER SYSTEM SET ssl = on;
ALTER SYSTEM SET ssl_cert_file = '/etc/ssl/certs/server.crt';
ALTER SYSTEM SET ssl_key_file = '/etc/ssl/private/server.key';
```

### 3. Access Control
```sql
-- Create read-only user
CREATE USER readonly_user WITH PASSWORD 'secure_password';
GRANT CONNECT ON DATABASE zerotrace TO readonly_user;
GRANT USAGE ON SCHEMA public TO readonly_user;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;
```

## Migration Scripts

### 1. Initial Migration
```sql
-- 001_initial_schema.sql
BEGIN;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create tables
-- (All table creation statements from above)

-- Create indexes
-- (All index creation statements from above)

-- Create functions and triggers
-- (All function and trigger creation statements from above)

COMMIT;
```

### 2. Add Partitioning
```sql
-- 002_add_partitioning.sql
BEGIN;

-- Create partitioned tables
-- (Partitioning statements from above)

-- Migrate existing data
INSERT INTO scans_partitioned SELECT * FROM scans;
INSERT INTO vulnerabilities_partitioned SELECT * FROM vulnerabilities;

-- Drop old tables and rename
DROP TABLE scans;
ALTER TABLE scans_partitioned RENAME TO scans;

DROP TABLE vulnerabilities;
ALTER TABLE vulnerabilities_partitioned RENAME TO vulnerabilities;

COMMIT;
```
