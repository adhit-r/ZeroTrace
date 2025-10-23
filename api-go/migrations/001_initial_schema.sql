-- 001_initial_schema.sql
-- Initial ZeroTrace database schema

BEGIN;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- ============================================================================
-- CORE TABLES
-- ============================================================================

-- Companies table
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE,
    settings JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Users table
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

-- Agents table
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    organization_id UUID,
    name VARCHAR(255) NOT NULL,
    hostname VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    version VARCHAR(50) NOT NULL,
    api_key VARCHAR(255) UNIQUE NOT NULL,
    capabilities JSONB DEFAULT '[]',
    location VARCHAR(100),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    current_scan_id UUID,
    performance_metrics JSONB DEFAULT '{}',
    
    -- Enhanced System Information
    os VARCHAR(100),
    os_name VARCHAR(100),
    os_version VARCHAR(100),
    os_build VARCHAR(100),
    kernel_version VARCHAR(100),
    cpu_model VARCHAR(255),
    cpu_cores INTEGER,
    memory_total_gb DECIMAL(10,2),
    storage_total_gb DECIMAL(10,2),
    gpu_model VARCHAR(255),
    ip_address INET,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Scans table
CREATE TABLE scans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
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

-- Vulnerabilities table
CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scan_id UUID NOT NULL REFERENCES scans(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id) ON DELETE SET NULL,
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

-- Dependencies table
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

-- ============================================================================
-- ENRICHMENT TABLES
-- ============================================================================

-- CVE Database table
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

-- Enrichment Results table
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

-- ============================================================================
-- REPORTING TABLES
-- ============================================================================

-- Reports table
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

-- ============================================================================
-- AUDIT AND LOGGING TABLES
-- ============================================================================

-- Audit Logs table
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

-- System Logs table
CREATE TABLE system_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level VARCHAR(20) NOT NULL,
    service VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Companies indexes
CREATE INDEX idx_companies_domain ON companies(domain);
CREATE INDEX idx_companies_status ON companies(status);
CREATE INDEX idx_companies_created_at ON companies(created_at);

-- Users indexes
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_last_login ON users(last_login);

-- Agents indexes
CREATE INDEX idx_agents_company_id ON agents(company_id);
CREATE INDEX idx_agents_organization_id ON agents(organization_id);
CREATE INDEX idx_agents_status ON agents(status);
CREATE INDEX idx_agents_api_key ON agents(api_key);
CREATE INDEX idx_agents_last_seen ON agents(last_seen);
CREATE INDEX idx_agents_current_scan ON agents(current_scan_id);

-- Scans indexes
CREATE INDEX idx_scans_company_id ON scans(company_id);
CREATE INDEX idx_scans_agent_id ON scans(agent_id);
CREATE INDEX idx_scans_status ON scans(status);
CREATE INDEX idx_scans_created_at ON scans(created_at);
CREATE INDEX idx_scans_start_time ON scans(start_time);
CREATE INDEX idx_scans_company_status ON scans(company_id, status);

-- Vulnerabilities indexes
CREATE INDEX idx_vulnerabilities_scan_id ON vulnerabilities(scan_id);
CREATE INDEX idx_vulnerabilities_company_id ON vulnerabilities(company_id);
CREATE INDEX idx_vulnerabilities_agent_id ON vulnerabilities(agent_id);
CREATE INDEX idx_vulnerabilities_severity ON vulnerabilities(severity);
CREATE INDEX idx_vulnerabilities_type ON vulnerabilities(type);
CREATE INDEX idx_vulnerabilities_cve_id ON vulnerabilities(cve_id);
CREATE INDEX idx_vulnerabilities_package ON vulnerabilities(package_name, package_version);
CREATE INDEX idx_vulnerabilities_status ON vulnerabilities(status);
CREATE INDEX idx_vulnerabilities_created_at ON vulnerabilities(created_at);
CREATE INDEX idx_vulnerabilities_company_severity ON vulnerabilities(company_id, severity);
CREATE INDEX idx_vulnerabilities_company_status ON vulnerabilities(company_id, status);

-- Dependencies indexes
CREATE INDEX idx_dependencies_scan_id ON dependencies(scan_id);
CREATE INDEX idx_dependencies_company_id ON dependencies(company_id);
CREATE INDEX idx_dependencies_name_version ON dependencies(name, version);
CREATE INDEX idx_dependencies_type ON dependencies(type);

-- CVE Database indexes
CREATE INDEX idx_cve_database_cve_id ON cve_database(cve_id);
CREATE INDEX idx_cve_database_severity ON cve_database(severity);
CREATE INDEX idx_cve_database_cvss_score ON cve_database(cvss_score);
CREATE INDEX idx_cve_database_published_date ON cve_database(published_date);

-- Enrichment Results indexes
CREATE INDEX idx_enrichment_vulnerability_id ON enrichment_results(vulnerability_id);
CREATE INDEX idx_enrichment_company_id ON enrichment_results(company_id);
CREATE INDEX idx_enrichment_type ON enrichment_results(enrichment_type);
CREATE INDEX idx_enrichment_risk_score ON enrichment_results(risk_score);

-- Reports indexes
CREATE INDEX idx_reports_company_id ON reports(company_id);
CREATE INDEX idx_reports_scan_id ON reports(scan_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_created_at ON reports(created_at);

-- Audit Logs indexes
CREATE INDEX idx_audit_logs_company_id ON audit_logs(company_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);

-- System Logs indexes
CREATE INDEX idx_system_logs_level ON system_logs(level);
CREATE INDEX idx_system_logs_service ON system_logs(service);
CREATE INDEX idx_system_logs_created_at ON system_logs(created_at);

-- ============================================================================
-- FULL-TEXT SEARCH INDEXES
-- ============================================================================

-- Full-text search for vulnerabilities
CREATE INDEX idx_vulnerabilities_search ON vulnerabilities 
    USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- Full-text search for CVE database
CREATE INDEX idx_cve_database_search ON cve_database 
    USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- ============================================================================
-- ROW-LEVEL SECURITY
-- ============================================================================

-- Enable RLS on all tables
ALTER TABLE companies ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE agents ENABLE ROW LEVEL SECURITY;
ALTER TABLE scans ENABLE ROW LEVEL SECURITY;
ALTER TABLE vulnerabilities ENABLE ROW LEVEL SECURITY;
ALTER TABLE dependencies ENABLE ROW LEVEL SECURITY;
ALTER TABLE enrichment_results ENABLE ROW LEVEL SECURITY;
ALTER TABLE reports ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_logs ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY companies_isolation ON companies
    FOR ALL USING (id = current_setting('app.current_company_id')::UUID);

CREATE POLICY users_isolation ON users
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY agents_isolation ON agents
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY scans_isolation ON scans
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY vulnerabilities_isolation ON vulnerabilities
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY dependencies_isolation ON dependencies
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY enrichment_isolation ON enrichment_results
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY reports_isolation ON reports
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY audit_logs_isolation ON audit_logs
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

-- ============================================================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================================================

-- Create trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to all tables with updated_at
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_agents_updated_at BEFORE UPDATE ON agents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_scans_updated_at BEFORE UPDATE ON scans
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_vulnerabilities_updated_at BEFORE UPDATE ON vulnerabilities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_cve_database_updated_at BEFORE UPDATE ON cve_database
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- VIEWS FOR ANALYTICS
-- ============================================================================

-- Scan Summary View
CREATE VIEW scan_summary AS
SELECT 
    s.id,
    s.company_id,
    s.agent_id,
    s.scan_type,
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
FROM scans s
LEFT JOIN vulnerabilities v ON s.id = v.scan_id
GROUP BY s.id, s.company_id, s.agent_id, s.scan_type, s.status, s.progress, 
         s.start_time, s.end_time, s.created_at;

-- Company Statistics View
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
LEFT JOIN scans s ON c.id = s.company_id
LEFT JOIN vulnerabilities v ON s.id = v.scan_id
GROUP BY c.id, c.name;

-- ============================================================================
-- FUNCTIONS FOR ANALYTICS
-- ============================================================================

-- Function to set company context
CREATE OR REPLACE FUNCTION set_company_context(company_id UUID)
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.current_company_id', company_id::text, false);
END;
$$ LANGUAGE plpgsql;

-- Function to get vulnerability statistics
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
    FROM vulnerabilities v
    WHERE v.company_id = p_company_id
    AND (p_date_from IS NULL OR v.created_at >= p_date_from)
    AND (p_date_to IS NULL OR v.created_at <= p_date_to);
END;
$$ LANGUAGE plpgsql;

COMMIT;
