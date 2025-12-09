-- 006_config_auditor_schema.sql
-- Configuration file auditing schema (Nipper Studio-like functionality)
-- Supports firewall, router, switch configuration file upload and analysis

BEGIN;

-- ============================================================================
-- CONFIGURATION FILE UPLOAD TABLES
-- ============================================================================

-- Configuration files table - stores uploaded config files
CREATE TABLE config_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- File metadata
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    file_hash VARCHAR(64) NOT NULL, -- SHA-256 hash for deduplication
    mime_type VARCHAR(100),
    
    -- Device information
    device_type VARCHAR(50) NOT NULL CHECK (device_type IN (
        'firewall', 'router', 'switch', 'load_balancer', 'waf', 
        'ids', 'ips', 'vpn_gateway', 'wireless_controller', 'other'
    )),
    manufacturer VARCHAR(100) NOT NULL, -- e.g., 'Cisco', 'Palo Alto', 'Fortinet', 'Juniper'
    model VARCHAR(100),
    firmware_version VARCHAR(100),
    device_name VARCHAR(255),
    device_location VARCHAR(255),
    
    -- Configuration metadata
    config_type VARCHAR(50) NOT NULL CHECK (config_type IN (
        'running_config', 'startup_config', 'backup_config', 'export_config', 'other'
    )),
    config_format VARCHAR(50), -- 'text', 'xml', 'json', 'binary'
    config_version VARCHAR(50),
    
    -- Parsing status
    parsing_status VARCHAR(50) DEFAULT 'pending' CHECK (parsing_status IN (
        'pending', 'parsing', 'parsed', 'failed', 'partial'
    )),
    parsing_error TEXT,
    parsed_data JSONB DEFAULT '{}', -- Parsed configuration structure
    
    -- Analysis status
    analysis_status VARCHAR(50) DEFAULT 'pending' CHECK (analysis_status IN (
        'pending', 'analyzing', 'completed', 'failed'
    )),
    analysis_started_at TIMESTAMP WITH TIME ZONE,
    analysis_completed_at TIMESTAMP WITH TIME ZONE,
    
    -- Metadata
    tags JSONB DEFAULT '[]',
    notes TEXT,
    metadata JSONB DEFAULT '{}',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Indexes
    UNIQUE(file_hash, company_id) -- Prevent duplicate uploads
);

-- Configuration findings table - stores security findings from config analysis
CREATE TABLE config_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_file_id UUID NOT NULL REFERENCES config_files(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    
    -- Finding details
    finding_type VARCHAR(50) NOT NULL CHECK (finding_type IN (
        'security_misconfiguration', 'weak_password', 'default_credentials',
        'insecure_protocol', 'missing_encryption', 'weak_cipher',
        'excessive_permissions', 'missing_logging', 'weak_access_control',
        'vulnerable_service', 'outdated_firmware', 'missing_patch',
        'compliance_violation', 'best_practice_violation', 'other'
    )),
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low', 'info')),
    category VARCHAR(50) NOT NULL CHECK (category IN (
        'authentication', 'authorization', 'encryption', 'network', 'logging',
        'access_control', 'compliance', 'firmware', 'services', 'other'
    )),
    
    -- Finding information
    title VARCHAR(500) NOT NULL,
    description TEXT NOT NULL,
    affected_component VARCHAR(255), -- e.g., 'interface eth0', 'rule 42', 'user admin'
    config_snippet TEXT, -- Relevant config lines
    line_numbers INTEGER[], -- Line numbers in original config
    
    -- Standards and compliance
    standard_id UUID REFERENCES config_standards(id) ON DELETE SET NULL,
    compliance_frameworks JSONB DEFAULT '[]', -- ['CIS', 'PCI-DSS', 'NIST']
    cve_id VARCHAR(20), -- If related to a CVE
    cvss_score DECIMAL(3,1),
    
    -- Remediation
    remediation TEXT,
    remediation_steps JSONB DEFAULT '[]', -- Step-by-step remediation
    remediation_priority VARCHAR(20) DEFAULT 'medium',
    estimated_effort VARCHAR(50), -- 'low', 'medium', 'high', 'critical'
    
    -- Risk assessment
    risk_score DECIMAL(3,2) DEFAULT 0.0,
    exploitability VARCHAR(20), -- 'low', 'medium', 'high', 'critical'
    impact VARCHAR(20), -- 'low', 'medium', 'high', 'critical'
    
    -- Status
    status VARCHAR(50) DEFAULT 'open' CHECK (status IN (
        'open', 'acknowledged', 'mitigated', 'resolved', 'false_positive', 'accepted_risk'
    )),
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Evidence and references
    evidence JSONB DEFAULT '{}',
    references JSONB DEFAULT '[]',
    tags JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Configuration standards table - manufacturer-specific security standards
CREATE TABLE config_standards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Standard identification
    standard_name VARCHAR(255) NOT NULL,
    standard_version VARCHAR(50),
    manufacturer VARCHAR(100) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    model_family VARCHAR(100), -- e.g., 'ASA', 'IOS', 'NX-OS' for Cisco
    
    -- Standard details
    category VARCHAR(50) NOT NULL,
    requirement_id VARCHAR(100) NOT NULL, -- e.g., 'CIS-1.1', 'NIST-AC-3'
    requirement_title VARCHAR(500) NOT NULL,
    requirement_description TEXT,
    
    -- Compliance mapping
    compliance_frameworks JSONB DEFAULT '[]', -- ['CIS', 'PCI-DSS', 'NIST', 'ISO27001']
    compliance_requirement VARCHAR(500), -- Specific compliance requirement text
    
    -- Configuration check
    check_type VARCHAR(50) NOT NULL CHECK (check_type IN (
        'presence', 'absence', 'value_match', 'value_range', 'pattern_match',
        'complex_rule', 'custom_script'
    )),
    check_config_path VARCHAR(500), -- JSON path or config location
    check_pattern TEXT, -- Regex or pattern to match
    expected_value TEXT, -- Expected value or pattern
    check_script TEXT, -- Custom validation script (if needed)
    
    -- Severity and priority
    default_severity VARCHAR(20) NOT NULL,
    priority VARCHAR(20) DEFAULT 'medium',
    
    -- Remediation guidance
    remediation_guidance TEXT,
    remediation_example TEXT, -- Example of correct configuration
    remediation_script TEXT, -- Script to auto-remediate (if possible)
    
    -- References
    references JSONB DEFAULT '[]',
    documentation_url VARCHAR(500),
    
    -- Status
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'deprecated', 'draft')),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Unique constraint
    UNIQUE(manufacturer, device_type, requirement_id, standard_version)
);

-- Configuration analysis results table - stores overall analysis results
CREATE TABLE config_analysis_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_file_id UUID NOT NULL REFERENCES config_files(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    
    -- Analysis summary
    total_findings INTEGER DEFAULT 0,
    critical_findings INTEGER DEFAULT 0,
    high_findings INTEGER DEFAULT 0,
    medium_findings INTEGER DEFAULT 0,
    low_findings INTEGER DEFAULT 0,
    info_findings INTEGER DEFAULT 0,
    
    -- Compliance scores
    compliance_scores JSONB DEFAULT '{}', -- {'CIS': 85, 'PCI-DSS': 72, 'NIST': 90}
    overall_security_score DECIMAL(5,2), -- 0-100 security score
    
    -- Analysis details
    analysis_version VARCHAR(50),
    standards_checked JSONB DEFAULT '[]', -- List of standards checked
    checks_performed INTEGER DEFAULT 0,
    checks_passed INTEGER DEFAULT 0,
    checks_failed INTEGER DEFAULT 0,
    
    -- Risk assessment
    overall_risk_score DECIMAL(3,2) DEFAULT 0.0,
    risk_level VARCHAR(20), -- 'low', 'medium', 'high', 'critical'
    
    -- Report
    report_path VARCHAR(500), -- Path to generated report
    report_format VARCHAR(20), -- 'html', 'pdf', 'json', 'xml'
    
    -- Metadata
    analysis_metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- INDEXES
-- ============================================================================

-- Config files indexes
CREATE INDEX idx_config_files_company_id ON config_files(company_id);
CREATE INDEX idx_config_files_device_type ON config_files(device_type);
CREATE INDEX idx_config_files_manufacturer ON config_files(manufacturer);
CREATE INDEX idx_config_files_parsing_status ON config_files(parsing_status);
CREATE INDEX idx_config_files_analysis_status ON config_files(analysis_status);
CREATE INDEX idx_config_files_created_at ON config_files(created_at);
CREATE INDEX idx_config_files_file_hash ON config_files(file_hash);

-- Config findings indexes
CREATE INDEX idx_config_findings_config_file_id ON config_findings(config_file_id);
CREATE INDEX idx_config_findings_company_id ON config_findings(company_id);
CREATE INDEX idx_config_findings_severity ON config_findings(severity);
CREATE INDEX idx_config_findings_category ON config_findings(category);
CREATE INDEX idx_config_findings_finding_type ON config_findings(finding_type);
CREATE INDEX idx_config_findings_status ON config_findings(status);
CREATE INDEX idx_config_findings_standard_id ON config_findings(standard_id);
CREATE INDEX idx_config_findings_created_at ON config_findings(created_at);
CREATE INDEX idx_config_findings_company_severity ON config_findings(company_id, severity);
CREATE INDEX idx_config_findings_company_status ON config_findings(company_id, status);

-- Config standards indexes
CREATE INDEX idx_config_standards_manufacturer ON config_standards(manufacturer);
CREATE INDEX idx_config_standards_device_type ON config_standards(device_type);
CREATE INDEX idx_config_standards_category ON config_standards(category);
CREATE INDEX idx_config_standards_status ON config_standards(status);
CREATE INDEX idx_config_standards_compliance ON config_standards USING gin(compliance_frameworks);

-- Config analysis results indexes
CREATE INDEX idx_config_analysis_config_file_id ON config_analysis_results(config_file_id);
CREATE INDEX idx_config_analysis_company_id ON config_analysis_results(company_id);
CREATE INDEX idx_config_analysis_created_at ON config_analysis_results(created_at);

-- ============================================================================
-- FULL-TEXT SEARCH INDEXES
-- ============================================================================

-- Full-text search for config findings
CREATE INDEX idx_config_findings_search ON config_findings 
    USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));

-- Full-text search for config standards
CREATE INDEX idx_config_standards_search ON config_standards 
    USING gin(to_tsvector('english', requirement_title || ' ' || COALESCE(requirement_description, '')));

-- ============================================================================
-- ROW-LEVEL SECURITY
-- ============================================================================

-- Enable RLS on all tables
ALTER TABLE config_files ENABLE ROW LEVEL SECURITY;
ALTER TABLE config_findings ENABLE ROW LEVEL SECURITY;
ALTER TABLE config_analysis_results ENABLE ROW LEVEL SECURITY;

-- Create RLS policies
CREATE POLICY config_files_isolation ON config_files
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY config_findings_isolation ON config_findings
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

CREATE POLICY config_analysis_isolation ON config_analysis_results
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);

-- ============================================================================
-- TRIGGERS FOR UPDATED_AT
-- ============================================================================

CREATE TRIGGER update_config_files_updated_at BEFORE UPDATE ON config_files
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_config_findings_updated_at BEFORE UPDATE ON config_findings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_config_standards_updated_at BEFORE UPDATE ON config_standards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_config_analysis_updated_at BEFORE UPDATE ON config_analysis_results
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- VIEWS FOR ANALYTICS
-- ============================================================================

-- Configuration security summary view
CREATE VIEW config_security_summary AS
SELECT 
    cf.company_id,
    cf.manufacturer,
    cf.device_type,
    COUNT(DISTINCT cf.id) as total_configs,
    COUNT(DISTINCT cfind.id) as total_findings,
    COUNT(CASE WHEN cfind.severity = 'critical' THEN 1 END) as critical_findings,
    COUNT(CASE WHEN cfind.severity = 'high' THEN 1 END) as high_findings,
    COUNT(CASE WHEN cfind.severity = 'medium' THEN 1 END) as medium_findings,
    COUNT(CASE WHEN cfind.severity = 'low' THEN 1 END) as low_findings,
    AVG(car.overall_security_score) as avg_security_score,
    MAX(cf.created_at) as last_config_uploaded
FROM config_files cf
LEFT JOIN config_findings cfind ON cf.id = cfind.config_file_id
LEFT JOIN config_analysis_results car ON cf.id = car.config_file_id
GROUP BY cf.company_id, cf.manufacturer, cf.device_type;

-- Compliance score by framework view
CREATE VIEW config_compliance_summary AS
SELECT 
    car.company_id,
    jsonb_object_keys(car.compliance_scores) as framework,
    AVG((car.compliance_scores->>jsonb_object_keys(car.compliance_scores))::DECIMAL) as avg_score,
    COUNT(*) as configs_analyzed
FROM config_analysis_results car
WHERE car.compliance_scores IS NOT NULL
GROUP BY car.company_id, jsonb_object_keys(car.compliance_scores);

COMMIT;

