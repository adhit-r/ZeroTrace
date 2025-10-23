# ZeroTrace Database Schema v2

## Overview
The ZeroTrace database schema v2 extends the original schema with comprehensive security scanning capabilities across multiple categories. This enhanced schema supports network security, compliance frameworks, system vulnerabilities, authentication security, database security, API security, container security, AI/ML security, IoT/OT security, privacy compliance, and Web3 security.

## New Features in v2

### 🔒 **Comprehensive Security Categories**
- **Network Security**: Port scanning, service detection, SSL/TLS analysis
- **Compliance Frameworks**: CIS, PCI-DSS, HIPAA, GDPR, SOC 2, ISO 27001
- **System Vulnerabilities**: OS patches, kernel issues, EOL software
- **Authentication Security**: Password policies, account security, privilege escalation
- **Database Security**: Multi-database vulnerability scanning
- **API Security**: REST, GraphQL, OWASP API Top 10
- **Container Security**: Docker, Kubernetes, IaC scanning
- **AI/ML Security**: Model vulnerabilities, training data security
- **IoT/OT Security**: Device discovery, firmware analysis
- **Privacy Compliance**: PII detection, GDPR/CCPA compliance
- **Web3 Security**: Smart contracts, wallets, DApp vulnerabilities

### 📊 **Enhanced Analytics**
- Risk heatmap data aggregation
- Compliance scoring across frameworks
- Security statistics by category
- Trend analysis and reporting
- Multi-dimensional vulnerability analysis

### 🚀 **Performance Optimizations**
- Partitioned tables for scalability
- Optimized indexes for fast queries
- Full-text search capabilities
- Row-level security for multi-tenancy
- Efficient data aggregation views

## Database Architecture

### Multi-Tenant Design
- **Company-based isolation**: All data partitioned by company_id
- **Row-level security**: Database-level access control
- **Partitioning**: Large tables partitioned by company_id
- **Indexing**: Optimized indexes for query performance

## Core Tables

### 1. Enhanced Vulnerabilities (vulnerabilities_v2)
```sql
CREATE TABLE vulnerabilities_v2 (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    organization_id UUID,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low', 'info')),
    category VARCHAR(50) NOT NULL CHECK (category IN (
        'application', 'network', 'configuration', 'system', 'auth', 
        'database', 'api', 'container', 'ai', 'iot', 'privacy', 'web3'
    )),
    status VARCHAR(50) DEFAULT 'open',
    discovered_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    risk_score DECIMAL(3,2) DEFAULT 0.0,
    exploit_complexity VARCHAR(20),
    attack_vector VARCHAR(50),
    compliance_frameworks JSONB DEFAULT '[]',
    remediation TEXT,
    references JSONB DEFAULT '[]',
    tags JSONB DEFAULT '[]',
    metadata JSONB DEFAULT '{}',
    enrichment_data JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 2. Network Security Findings
```sql
CREATE TABLE network_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    scan_id UUID,
    host VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL,
    protocol VARCHAR(10) NOT NULL,
    service_name VARCHAR(100),
    service_version VARCHAR(100),
    banner TEXT,
    ssl_enabled BOOLEAN DEFAULT FALSE,
    ssl_version VARCHAR(20),
    ssl_cipher VARCHAR(100),
    ssl_certificate_issuer VARCHAR(255),
    ssl_certificate_expiry TIMESTAMP WITH TIME ZONE,
    vulnerability_count INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'open',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 3. Compliance Checks
```sql
CREATE TABLE compliance_checks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    framework VARCHAR(50) NOT NULL CHECK (framework IN ('CIS', 'PCI-DSS', 'HIPAA', 'GDPR', 'SOC2', 'ISO27001')),
    category VARCHAR(100) NOT NULL,
    requirement VARCHAR(500) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('pass', 'fail', 'not_applicable', 'error')),
    severity VARCHAR(20),
    description TEXT,
    remediation TEXT,
    evidence JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 4. System Vulnerabilities
```sql
CREATE TABLE system_vulnerabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    vulnerability_type VARCHAR(50) NOT NULL CHECK (vulnerability_type IN ('os_patch', 'kernel', 'driver', 'eol_software')),
    os_name VARCHAR(100),
    os_version VARCHAR(100),
    component_name VARCHAR(255),
    component_version VARCHAR(100),
    cve_id VARCHAR(20),
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    patch_available BOOLEAN DEFAULT FALSE,
    patch_url VARCHAR(500),
    eol_date DATE,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 5. Authentication Security Findings
```sql
CREATE TABLE auth_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    finding_type VARCHAR(50) NOT NULL CHECK (finding_type IN ('password_policy', 'account_security', 'privilege_escalation', 'auth_bypass')),
    user_account VARCHAR(255),
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    risk_score DECIMAL(3,2) DEFAULT 0.0,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 6. Database Security Findings
```sql
CREATE TABLE database_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    database_type VARCHAR(50) NOT NULL CHECK (database_type IN ('postgresql', 'mysql', 'mongodb', 'redis', 'sqlserver', 'oracle')),
    host VARCHAR(255) NOT NULL,
    port INTEGER,
    database_name VARCHAR(100),
    finding_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 7. API Security Findings
```sql
CREATE TABLE api_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    api_type VARCHAR(20) NOT NULL CHECK (api_type IN ('REST', 'GraphQL', 'SOAP', 'gRPC')),
    endpoint VARCHAR(500) NOT NULL,
    method VARCHAR(10),
    finding_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 8. Container Security Findings
```sql
CREATE TABLE container_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    container_id VARCHAR(100),
    image_name VARCHAR(255),
    image_tag VARCHAR(100),
    finding_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 9. AI/ML Security Findings
```sql
CREATE TABLE ai_ml_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    finding_type VARCHAR(50) NOT NULL CHECK (finding_type IN ('model_vulnerability', 'training_data', 'llm_security', 'bias_fairness')),
    model_name VARCHAR(255),
    model_version VARCHAR(100),
    framework VARCHAR(100),
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 10. IoT/OT Security Findings
```sql
CREATE TABLE iot_ot_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    device_id VARCHAR(100),
    protocol VARCHAR(20),
    finding_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 11. Privacy and Compliance Findings
```sql
CREATE TABLE privacy_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    finding_type VARCHAR(50) NOT NULL CHECK (finding_type IN ('pii_detection', 'gdpr_compliance', 'ccpa_compliance', 'data_retention')),
    data_type VARCHAR(50),
    location VARCHAR(500),
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 12. Web3 Security Findings
```sql
CREATE TABLE web3_findings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL,
    company_id UUID NOT NULL,
    finding_type VARCHAR(50) NOT NULL CHECK (finding_type IN ('smart_contract', 'wallet_security', 'dapp_security', 'defi_risks')),
    contract_address VARCHAR(100),
    network VARCHAR(50),
    severity VARCHAR(20) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    remediation TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Partitioning Strategy

### Hash Partitioning by Company ID
All major tables are partitioned by `company_id` using hash partitioning for optimal performance:

```sql
-- Example: vulnerabilities_v2 partitioning
CREATE TABLE vulnerabilities_v2_partitioned (
    LIKE vulnerabilities_v2 INCLUDING ALL
) PARTITION BY HASH (company_id);

-- Create 4 partitions for load distribution
CREATE TABLE vulnerabilities_v2_partition_0 PARTITION OF vulnerabilities_v2_partitioned
    FOR VALUES WITH (modulus 4, remainder 0);
CREATE TABLE vulnerabilities_v2_partition_1 PARTITION OF vulnerabilities_v2_partitioned
    FOR VALUES WITH (modulus 4, remainder 1);
CREATE TABLE vulnerabilities_v2_partition_2 PARTITION OF vulnerabilities_v2_partitioned
    FOR VALUES WITH (modulus 4, remainder 2);
CREATE TABLE vulnerabilities_v2_partition_3 PARTITION OF vulnerabilities_v2_partitioned
    FOR VALUES WITH (modulus 4, remainder 3);
```

## Indexing Strategy

### Performance Indexes
- **Primary keys**: UUID with gen_random_uuid()
- **Foreign keys**: company_id, agent_id, scan_id
- **Query optimization**: severity, category, status, created_at
- **Composite indexes**: company_id + severity, company_id + category
- **Full-text search**: GIN indexes on title and description fields

### Example Indexes
```sql
-- Vulnerabilities v2 indexes
CREATE INDEX idx_vulnerabilities_v2_agent_id ON vulnerabilities_v2_partitioned(agent_id);
CREATE INDEX idx_vulnerabilities_v2_company_id ON vulnerabilities_v2_partitioned(company_id);
CREATE INDEX idx_vulnerabilities_v2_severity ON vulnerabilities_v2_partitioned(severity);
CREATE INDEX idx_vulnerabilities_v2_category ON vulnerabilities_v2_partitioned(category);
CREATE INDEX idx_vulnerabilities_v2_status ON vulnerabilities_v2_partitioned(status);
CREATE INDEX idx_vulnerabilities_v2_risk_score ON vulnerabilities_v2_partitioned(risk_score);
CREATE INDEX idx_vulnerabilities_v2_discovered_at ON vulnerabilities_v2_partitioned(discovered_at);
CREATE INDEX idx_vulnerabilities_v2_company_severity ON vulnerabilities_v2_partitioned(company_id, severity);
CREATE INDEX idx_vulnerabilities_v2_company_category ON vulnerabilities_v2_partitioned(company_id, category);
CREATE INDEX idx_vulnerabilities_v2_company_status ON vulnerabilities_v2_partitioned(company_id, status);

-- Full-text search
CREATE INDEX idx_vulnerabilities_v2_search ON vulnerabilities_v2_partitioned 
    USING gin(to_tsvector('english', title || ' ' || COALESCE(description, '')));
```

## Row-Level Security (RLS)

### Multi-Tenant Data Isolation
All tables implement row-level security to ensure data isolation between companies:

```sql
-- Enable RLS
ALTER TABLE vulnerabilities_v2_partitioned ENABLE ROW LEVEL SECURITY;

-- Create isolation policy
CREATE POLICY vulnerabilities_v2_isolation ON vulnerabilities_v2_partitioned
    FOR ALL USING (company_id = current_setting('app.current_company_id')::UUID);
```

### Context Setting Function
```sql
CREATE OR REPLACE FUNCTION set_company_context(company_id UUID)
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.current_company_id', company_id::text, false);
END;
$$ LANGUAGE plpgsql;
```

## Analytics Views

### 1. Comprehensive Vulnerability Summary
```sql
CREATE VIEW vulnerability_summary_v2 AS
SELECT 
    v.id,
    v.agent_id,
    v.company_id,
    v.title,
    v.severity,
    v.category,
    v.status,
    v.risk_score,
    v.discovered_at,
    v.last_seen,
    a.hostname,
    a.os,
    a.os_version
FROM vulnerabilities_v2_partitioned v
LEFT JOIN agents a ON v.agent_id = a.id;
```

### 2. Compliance Framework Summary
```sql
CREATE VIEW compliance_summary AS
SELECT 
    c.company_id,
    c.framework,
    COUNT(*) as total_checks,
    COUNT(CASE WHEN c.status = 'pass' THEN 1 END) as passed_checks,
    COUNT(CASE WHEN c.status = 'fail' THEN 1 END) as failed_checks,
    COUNT(CASE WHEN c.status = 'not_applicable' THEN 1 END) as not_applicable_checks,
    ROUND(
        (COUNT(CASE WHEN c.status = 'pass' THEN 1 END)::DECIMAL / 
         NULLIF(COUNT(CASE WHEN c.status IN ('pass', 'fail') THEN 1 END), 0)) * 100, 2
    ) as compliance_score
FROM compliance_checks_partitioned c
GROUP BY c.company_id, c.framework;
```

### 3. Network Security Summary
```sql
CREATE VIEW network_security_summary AS
SELECT 
    n.company_id,
    n.agent_id,
    COUNT(*) as total_findings,
    COUNT(CASE WHEN n.ssl_enabled = false THEN 1 END) as non_ssl_services,
    COUNT(CASE WHEN n.vulnerability_count > 0 THEN 1 END) as vulnerable_services,
    COUNT(DISTINCT n.host) as unique_hosts,
    COUNT(DISTINCT n.port) as unique_ports
FROM network_findings_partitioned n
GROUP BY n.company_id, n.agent_id;
```

## Analytics Functions

### 1. Comprehensive Security Statistics
```sql
CREATE OR REPLACE FUNCTION get_security_statistics(
    p_company_id UUID,
    p_date_from TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    p_date_to TIMESTAMP WITH TIME ZONE DEFAULT NULL
)
RETURNS TABLE (
    total_vulnerabilities BIGINT,
    vulnerabilities_by_category JSONB,
    vulnerabilities_by_severity JSONB,
    compliance_scores JSONB,
    network_findings_count BIGINT,
    system_vulnerabilities_count BIGINT,
    auth_findings_count BIGINT,
    database_findings_count BIGINT,
    api_findings_count BIGINT,
    container_findings_count BIGINT,
    ai_ml_findings_count BIGINT,
    iot_ot_findings_count BIGINT,
    privacy_findings_count BIGINT,
    web3_findings_count BIGINT
) AS $$
-- Implementation provides comprehensive security statistics
-- across all security categories with date filtering
$$ LANGUAGE plpgsql;
```

### 2. Risk Heatmap Data
```sql
CREATE OR REPLACE FUNCTION get_risk_heatmap_data(
    p_company_id UUID,
    p_date_from TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    p_date_to TIMESTAMP WITH TIME ZONE DEFAULT NULL
)
RETURNS TABLE (
    category VARCHAR(50),
    severity VARCHAR(20),
    count BIGINT,
    risk_score DECIMAL(3,2)
) AS $$
-- Implementation provides risk heatmap data
-- for visualization across categories and severities
$$ LANGUAGE plpgsql;
```

## Migration Strategy

### 1. Database Setup
```bash
# Run the database setup script
./scripts/setup-database.sh
```

### 2. Migration Execution
```bash
# Navigate to migrations directory
cd api-go/migrations

# Run migrations
go run migrate.go
```

### 3. Environment Variables
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=zerotrace
export DB_SSLMODE=disable
```

## Performance Considerations

### 1. Connection Pooling
- Configure connection pool size based on expected load
- Use connection pooling libraries (pgxpool, sqlx)
- Monitor connection usage and adjust accordingly

### 2. Query Optimization
- Use EXPLAIN ANALYZE for query optimization
- Monitor slow queries with pg_stat_statements
- Implement query result caching where appropriate

### 3. Maintenance
- Regular VACUUM and ANALYZE operations
- Monitor table sizes and partition accordingly
- Implement automated maintenance procedures

## Security Considerations

### 1. Data Encryption
- Enable SSL/TLS for database connections
- Consider encryption at rest for sensitive data
- Implement proper key management

### 2. Access Control
- Use least privilege principle for database users
- Implement role-based access control
- Regular security audits and access reviews

### 3. Backup and Recovery
- Implement regular automated backups
- Test backup and recovery procedures
- Consider point-in-time recovery capabilities

## Monitoring and Alerting

### 1. Database Metrics
- Connection count and usage
- Query performance and slow queries
- Table sizes and growth rates
- Index usage and effectiveness

### 2. Application Metrics
- API response times
- Error rates and types
- Data processing throughput
- User activity patterns

### 3. Security Metrics
- Failed authentication attempts
- Unusual data access patterns
- Compliance check failures
- Vulnerability detection rates

This enhanced database schema provides the foundation for comprehensive security scanning and analysis across all categories, with optimized performance and multi-tenant isolation for enterprise-scale deployments.
