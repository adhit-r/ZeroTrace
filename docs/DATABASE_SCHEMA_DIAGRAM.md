# ZeroTrace Database Schema Diagram

## Database Type

**PostgreSQL (RDBMS) with JSONB Support**

- **Primary Database**: PostgreSQL 15+
- **Data Model**: Relational Database Management System (RDBMS)
- **JSONB Usage**: Used for flexible/semi-structured data storage (metadata, enrichment_data, settings, etc.)
- **Architecture**: Hybrid approach - structured relational tables with JSONB columns for extensibility

## Schema Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ZEROTRACE DATABASE SCHEMA                            │
│                      PostgreSQL RDBMS with JSONB Support                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Core Entity Relationship Diagram

```
┌──────────────┐
│  companies   │
│──────────────│
│ id (PK)      │◄─────┐
│ name         │      │
│ domain       │      │
│ settings     │      │
│ status       │      │
│ created_at   │      │
│ updated_at   │      │
└──────────────┘      │
                       │
                       │ company_id (FK)
                       │
        ┌──────────────┼──────────────┬──────────────┐
        │              │              │              │
        ▼              ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│    users     │ │    agents    │ │    scans     │ │vulnerabilities│
│──────────────│ │──────────────│ │──────────────│ │──────────────│
│ id (PK)      │ │ id (PK)      │ │ id (PK)      │ │ id (PK)      │
│ company_id   │ │ company_id   │ │ company_id   │ │ scan_id      │
│ email        │ │ name         │ │ agent_id     │ │ company_id   │
│ password_hash│ │ hostname     │ │ scan_type    │ │ agent_id     │
│ name         │ │ status       │ │ status       │ │ type         │
│ role         │ │ version      │ │ progress     │ │ severity     │
│ status       │ │ api_key      │ │ start_time   │ │ title        │
│ last_login   │ │ capabilities │ │ end_time     │ │ description  │
│ created_at   │ │ (JSONB)      │ │ options      │ │ cve_id       │
│ updated_at   │ │ metrics     │ │ (JSONB)      │ │ cvss_score   │
└──────────────┘ │ (JSONB)      │ │ results      │ │ package_name │
                 │ os           │ │ (JSONB)      │ │ references   │
                 │ os_version   │ │ metadata     │ │ (JSONB)      │
                 │ cpu_cores    │ │ (JSONB)      │ │ enrichment_  │
                 │ memory_gb    │ │ created_at   │ │   data (JSONB)│
                 │ created_at   │ │ updated_at   │ │ created_at   │
                 │ updated_at   │ └──────────────┘ │ updated_at   │
                 └──────────────┘                  └──────────────┘
                          │                                │
                          │ agent_id (FK)                  │
                          └────────────────────────────────┘
```

## Enhanced Security Schema (v2)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    ENHANCED VULNERABILITY SCHEMA (V2)                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────────────────┐
│ vulnerabilities_v2   │
│──────────────────────│
│ id (PK)              │
│ agent_id (FK)        │
│ company_id (FK)      │
│ organization_id      │
│ title                │
│ description          │
│ severity             │
│ category             │
│ status               │
│ risk_score           │
│ compliance_frameworks│ (JSONB)
│ references           │ (JSONB)
│ tags                 │ (JSONB)
│ metadata             │ (JSONB)
│ enrichment_data      │ (JSONB)
│ discovered_at        │
│ last_seen            │
│ created_at           │
│ updated_at           │
└──────────────────────┘
         │
         │ (partitioned by company_id)
         │
         ├─────────────────────────────────────────────────────┐
         │                                                     │
         ▼                                                     ▼
┌──────────────────────┐                            ┌──────────────────────┐
│ network_findings     │                            │ compliance_checks    │
│──────────────────────│                            │──────────────────────│
│ id (PK)              │                            │ id (PK)              │
│ agent_id (FK)        │                            │ agent_id (FK)        │
│ company_id (FK)      │                            │ company_id (FK)      │
│ scan_id              │                            │ framework            │
│ host                 │                            │ category             │
│ port                 │                            │ requirement          │
│ protocol             │                            │ status               │
│ service_name         │                            │ severity             │
│ service_version      │                            │ description          │
│ ssl_enabled          │                            │ remediation          │
│ ssl_version          │                            │ evidence (JSONB)     │
│ ssl_cipher           │                            │ metadata (JSONB)     │
│ ssl_cert_expiry      │                            │ created_at           │
│ vulnerability_count  │                            │ updated_at           │
│ metadata (JSONB)     │                            └──────────────────────┘
│ created_at           │
│ updated_at           │
└──────────────────────┘

┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────────┐
│ system_vulnerabilities│  │ auth_findings        │  │ database_findings    │
│──────────────────────│  │──────────────────────│  │──────────────────────│
│ id (PK)              │  │ id (PK)              │  │ id (PK)              │
│ agent_id (FK)        │  │ agent_id (FK)        │  │ agent_id (FK)        │
│ company_id (FK)      │  │ company_id (FK)      │  │ company_id (FK)      │
│ vulnerability_type   │  │ finding_type         │  │ database_type        │
│ os_name              │  │ user_account         │  │ host                 │
│ os_version           │  │ severity             │  │ port                 │
│ component_name       │  │ title                │  │ database_name        │
│ component_version    │  │ description          │  │ finding_type         │
│ cve_id               │  │ risk_score           │  │ severity             │
│ severity             │  │ remediation          │  │ title                │
│ title                │  │ metadata (JSONB)     │  │ description          │
│ description          │  │ created_at           │  │ remediation          │
│ remediation          │  │ updated_at           │  │ metadata (JSONB)     │
│ patch_available      │  └──────────────────────┘  │ created_at           │
│ patch_url            │                            │ updated_at           │
│ eol_date             │                            └──────────────────────┘
│ metadata (JSONB)     │
│ created_at           │
│ updated_at           │
└──────────────────────┘

┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────────┐
│ api_findings         │  │ container_findings  │  │ ai_ml_findings       │
│──────────────────────│  │──────────────────────│  │──────────────────────│
│ id (PK)              │  │ id (PK)              │  │ id (PK)              │
│ agent_id (FK)        │  │ agent_id (FK)        │  │ agent_id (FK)        │
│ company_id (FK)      │  │ company_id (FK)      │  │ company_id (FK)      │
│ api_type             │  │ container_id         │  │ finding_type         │
│ endpoint             │  │ image_name           │  │ model_name           │
│ method               │  │ image_tag           │  │ model_version        │
│ finding_type         │  │ finding_type         │  │ framework            │
│ severity             │  │ severity             │  │ severity             │
│ title                │  │ title                │  │ title                │
│ description          │  │ description          │  │ description          │
│ remediation          │  │ remediation          │  │ remediation          │
│ metadata (JSONB)     │  │ metadata (JSONB)     │  │ metadata (JSONB)     │
│ created_at           │  │ created_at           │  │ created_at           │
│ updated_at           │  │ updated_at           │  │ updated_at           │
└──────────────────────┘  └──────────────────────┘  └──────────────────────┘

┌──────────────────────┐  ┌──────────────────────┐  ┌──────────────────────┐
│ iot_ot_findings      │  │ privacy_findings    │  │ web3_findings        │
│──────────────────────│  │──────────────────────│  │──────────────────────│
│ id (PK)              │  │ id (PK)              │  │ id (PK)              │
│ agent_id (FK)        │  │ agent_id (FK)        │  │ agent_id (FK)        │
│ company_id (FK)      │  │ company_id (FK)      │  │ company_id (FK)      │
│ device_type          │  │ finding_type         │  │ finding_type         │
│ device_id            │  │ data_type            │  │ contract_address     │
│ protocol             │  │ location             │  │ network              │
│ finding_type         │  │ severity             │  │ severity             │
│ severity             │  │ title                │  │ title                │
│ title                │  │ description          │  │ description          │
│ description          │  │ remediation          │  │ remediation          │
│ remediation          │  │ metadata (JSONB)     │  │ metadata (JSONB)     │
│ metadata (JSONB)     │  │ created_at           │  │ created_at           │
│ created_at           │  │ updated_at           │  │ updated_at           │
│ updated_at           │  └──────────────────────┘  └──────────────────────┘
└──────────────────────┘
```

## Enrichment & CVE Tables

```
┌──────────────────────┐
│ cve_database         │
│──────────────────────│
│ id (PK)              │
│ cve_id (UNIQUE)      │
│ title                │
│ description          │
│ cvss_score           │
│ cvss_vector          │
│ severity             │
│ published_date       │
│ last_modified_date   │
│ references (JSONB)   │
│ affected_products    │
│        (JSONB)       │
│ created_at           │
│ updated_at           │
└──────────────────────┘
         │
         │ cve_id
         │
         ▼
┌──────────────────────┐
│ enrichment_results   │
│──────────────────────│
│ id (PK)              │
│ vulnerability_id (FK)│
│ company_id (FK)      │
│ enrichment_type      │
│ risk_score           │
│ trend                │
│ recommendations      │
│        (JSONB)       │
│ threat_intelligence  │
│        (JSONB)       │
│ ml_predictions       │
│        (JSONB)       │
│ confidence_score     │
│ created_at           │
└──────────────────────┘
         │
         │ vulnerability_id
         │
         ▼
┌──────────────────────┐
│ vulnerabilities      │
│   (v1 table)         │
└──────────────────────┘
```

## Supporting Tables

```
┌──────────────────────┐      ┌──────────────────────┐
│ dependencies         │      │ reports              │
│──────────────────────│      │──────────────────────│
│ id (PK)              │      │ id (PK)              │
│ scan_id (FK)         │      │ company_id (FK)      │
│ company_id (FK)      │      │ scan_id (FK)         │
│ name                 │      │ name                 │
│ version              │      │ format               │
│ type                 │      │ status               │
│ location             │      │ file_path             │
│ license              │      │ file_size            │
│ vulnerabilities      │      │ download_url         │
│        (JSONB)       │      │ options (JSONB)      │
│ created_at           │      │ generated_by (FK)    │
└──────────────────────┘      │ created_at           │
                               │ completed_at          │
                               └──────────────────────┘

┌──────────────────────┐      ┌──────────────────────┐
│ audit_logs           │      │ system_logs          │
│──────────────────────│      │──────────────────────│
│ id (PK)              │      │ id (PK)              │
│ company_id (FK)      │      │ level                │
│ user_id (FK)         │      │ service              │
│ action               │      │ message              │
│ resource_type        │      │ metadata (JSONB)     │
│ resource_id          │      │ created_at           │
│ old_values (JSONB)   │      └──────────────────────┘
│ new_values (JSONB)   │
│ ip_address           │
│ user_agent           │
│ created_at           │
└──────────────────────┘
```

## Key Design Patterns

### 1. Multi-Tenancy
- All tables include `company_id` for data isolation
- Row-Level Security (RLS) policies enforce isolation
- Partitioning by `company_id` for performance

### 2. JSONB Usage
JSONB columns are used for:
- **Flexible metadata**: `metadata JSONB` - stores variable attributes
- **Enrichment data**: `enrichment_data JSONB` - CVE enrichment results
- **Settings**: `settings JSONB` - company/user preferences
- **Arrays**: `capabilities JSONB[]`, `references JSONB[]` - lists of objects
- **Threat intelligence**: `threat_intelligence JSONB` - threat feed data
- **ML predictions**: `ml_predictions JSONB` - ML model outputs

### 3. Partitioning Strategy
- **Hash Partitioning**: Large tables partitioned by `company_id`
- **Partitions**: 4 partitions for load distribution
- **Tables Partitioned**: 
  - `vulnerabilities_v2_partitioned`
  - `network_findings_partitioned`
  - `compliance_checks_partitioned`

### 4. Indexing Strategy
- **Primary Keys**: UUID with `gen_random_uuid()`
- **Foreign Keys**: `company_id`, `agent_id`, `scan_id`
- **Composite Indexes**: `(company_id, severity)`, `(company_id, category)`
- **Full-Text Search**: GIN indexes on title/description
- **Performance Indexes**: severity, category, status, created_at

### 5. Row-Level Security (RLS)
- All tables have RLS enabled
- Policies enforce `company_id` isolation
- Context function: `set_company_context(company_id)`

## Data Flow

```
Agent Scan
    │
    ▼
┌──────────┐
│  scans   │
└──────────┘
    │
    ├──► vulnerabilities (v1)
    │         │
    │         └──► enrichment_results
    │                    │
    │                    └──► cve_database
    │
    └──► vulnerabilities_v2 (v2)
             │
             ├──► network_findings
             ├──► compliance_checks
             ├──► system_vulnerabilities
             ├──► auth_findings
             ├──► database_findings
             ├──► api_findings
             ├──► container_findings
             ├──► ai_ml_findings
             ├──► iot_ot_findings
             ├──► privacy_findings
             └──► web3_findings
```

## Configuration Auditor Tables (Nipper Studio-like)

```
┌──────────────────────┐
│ config_files         │
│──────────────────────│
│ id (PK)              │
│ company_id (FK)      │
│ uploaded_by (FK)     │
│ filename             │
│ file_path            │
│ file_hash            │
│ device_type          │
│ manufacturer         │
│ model                │
│ firmware_version     │
│ config_type          │
│ parsing_status       │
│ parsed_data (JSONB)  │
│ analysis_status      │
│ created_at           │
│ updated_at           │
└──────────────────────┘
         │
         │ config_file_id
         │
         ├──────────────────────────────┐
         │                              │
         ▼                              ▼
┌──────────────────────┐      ┌──────────────────────┐
│ config_findings      │      │ config_analysis_     │
│──────────────────────│      │   results            │
│ id (PK)              │      │──────────────────────│
│ config_file_id (FK)  │      │ id (PK)              │
│ company_id (FK)      │      │ config_file_id (FK)  │
│ finding_type         │      │ company_id (FK)      │
│ severity             │      │ total_findings       │
│ category             │      │ critical_findings    │
│ title                │      │ compliance_scores   │
│ description          │      │   (JSONB)            │
│ affected_component   │      │ overall_security_    │
│ config_snippet       │      │   score              │
│ line_numbers         │      │ overall_risk_score   │
│ standard_id (FK)     │      │ report_path          │
│ compliance_          │      │ created_at           │
│   frameworks (JSONB) │      │ updated_at           │
│ remediation          │      └──────────────────────┘
│ risk_score           │
│ status               │
│ created_at           │
│ updated_at           │
└──────────────────────┘
         │
         │ standard_id
         │
         ▼
┌──────────────────────┐
│ config_standards     │
│──────────────────────│
│ id (PK)              │
│ standard_name        │
│ manufacturer         │
│ device_type          │
│ requirement_id       │
│ requirement_title    │
│ check_type           │
│ check_pattern        │
│ expected_value       │
│ default_severity     │
│ remediation_guidance │
│ compliance_          │
│   frameworks (JSONB) │
│ status               │
│ created_at           │
│ updated_at           │
└──────────────────────┘
```

## Exploit Storage

**Exploits are stored in multiple locations:**

1. **Vulnerabilities Table (v1)**:
   - `exploit_available BOOLEAN` - Whether exploits exist
   - `exploit_count INTEGER` - Number of available exploits

2. **Vulnerabilities V2 Table**:
   - `exploit_complexity VARCHAR(20)` - Complexity of exploitation
   - `attack_vector VARCHAR(50)` - Attack vector information

3. **Enrichment Results Table**:
   - `threat_intelligence JSONB` - Contains exploit data from:
     - Exploit-DB API
     - CISA KEV (Known Exploited Vulnerabilities)
     - Other threat intelligence feeds
   - Structure:
     ```json
     {
       "exploit_availability": true,
       "public_exploits": [
         {
           "id": "12345",
           "title": "Exploit Title",
           "description": "...",
           "date_published": "2024-01-01",
           "verified": true,
           "platform": "linux",
           "type": "remote"
         }
       ],
       "exploit_count": 3,
       "cisa_kev": {
         "cveID": "CVE-2024-1234",
         "knownRansomwareCampaignUse": "Yes"
       }
     }
     ```

4. **Exploit Intelligence Service**:
   - Located: `enrichment-python/app/ai_services/exploit_intelligence.py`
   - Fetches from Exploit-DB API
   - Caches in Redis
   - Integrated into enrichment pipeline

**Note**: There is no dedicated `exploits` table. Exploit data is embedded in vulnerability and enrichment records for better performance and to maintain relationship with CVEs.

## Summary

**Database Type**: PostgreSQL RDBMS  
**JSONB Usage**: Extensively used for flexible data storage  
**Architecture**: Hybrid relational + JSONB for extensibility  
**Multi-Tenancy**: Company-based isolation with RLS  
**Performance**: Partitioning, indexing, and full-text search  
**Scalability**: Designed for 1000+ agents, 100+ companies  
**Configuration Auditing**: Nipper Studio-like config file analysis  
**Exploit Storage**: Embedded in vulnerabilities and enrichment data

