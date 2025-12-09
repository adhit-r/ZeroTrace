# ZeroTrace

<div align="center">

**Enterprise-Grade Vulnerability Detection and Management Platform**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8.svg)](https://golang.org/)
[![Python Version](https://img.shields.io/badge/Python-3.9+-3776AB.svg)](https://python.org/)
[![React Version](https://img.shields.io/badge/React-19+-61DAFB.svg)](https://reactjs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CVE Database](https://img.shields.io/badge/CVE%20Database-320K+-red.svg)]()
[![KEV Integration](https://img.shields.io/badge/KEV-1.4K+%20CVEs-orange.svg)]()

[Quick Start](#quick-start) • [Features](#features) • [Architecture](#architecture) • [Documentation](#documentation)

</div>

---

## Overview

ZeroTrace is a high-performance vulnerability detection and management platform engineered for enterprise-scale deployments. Unlike traditional vulnerability scanners that require agent installation on every target device, ZeroTrace employs a **hybrid architecture** combining agent-based software scanning with **agentless network discovery**, enabling comprehensive security assessment without device-level dependencies.

### Key Differentiators

- **Agentless Network Scanning** - Single agent scans entire network segments (no installation on targets)
- **Ultra-Optimized Performance** - 10,000x faster enrichment, 95% CPU reduction
- **Hybrid CPE Matching** - Two-tier system: exact match + semantic search
- **Multi-Tenant Enterprise** - Organization isolation with universal agent binary
- **AI-Powered Analysis** - Predictive vulnerability analysis and exploit intelligence
- **KEV Integration** - CISA Known Exploited Vulnerabilities with 1,468+ CVEs

---

## Features

### 1. Agentless Network Scanning

**How It Works:**

```
┌─────────────────────────────────────────────────────────────────┐
│                    ZeroTrace Agent (Sensor)                     │
│              Installed on ONE machine in network                 │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Step 1: Network Discovery                               │  │
│  │  • Detects local subnet (e.g., 192.168.1.0/24)          │  │
│  │  • Uses Nmap for comprehensive device discovery           │  │
│  │  • Scans all IPs in subnet range                         │  │
│  └──────────────────────────────────────────────────────────┘  │
│                            │                                     │
│                            ▼                                     │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Step 2: Device Classification                           │  │
│  │  • Mobile Devices (Android/iOS via OS fingerprinting)    │  │
│  │  • IoT Devices (MQTT, CoAP, UPnP ports)                  │  │
│  │  • Servers (SSH, HTTP, Database ports)                   │  │
│  │  • Network Infrastructure (Switches, Routers)           │  │
│  │  • Laptops/Desktops (Windows, macOS, Linux)             │  │
│  └──────────────────────────────────────────────────────────┘  │
│                            │                                     │
│                            ▼                                     │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Step 3: Vulnerability Scanning                          │  │
│  │  • Nuclei templates for CVE detection                    │  │
│  │  • Port scanning and service detection                   │  │
│  │  • Configuration auditing                                │  │
│  │  • SSL/TLS certificate analysis                          │  │
│  └──────────────────────────────────────────────────────────┘  │
│                            │                                     │
│                            ▼                                     │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │  Step 4: Results Aggregation                             │  │
│  │  • Device classification and risk scoring                │  │
│  │  • Vulnerability correlation with CVE database           │  │
│  │  • Network topology mapping                              │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │   ZeroTrace API        │
              │   (Port 8080)          │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │   PostgreSQL Database   │
              │   + Network Topology    │
              └─────────────────────────┘
                            │
                            ▼
              ┌─────────────────────────┐
              │   React Frontend        │
              │   • Network Topology    │
              │   • Attack Paths        │
              │   • Device Inventory    │
              └─────────────────────────┘
```

**What Gets Scanned:**

| Device Type | Detection Method | Examples |
|------------|------------------|----------|
| **Mobile Devices** | OS fingerprinting, AirPlay/AirPrint services | iPhones, Android phones, Tablets |
| **IoT Devices** | MQTT (1883), CoAP (5683), UPnP (1900) ports | Smart bulbs, Nest, Echo, Smart TVs |
| **Servers** | Common server ports (22, 80, 443, 3306) | Web servers, Database servers, File servers |
| **Laptops/Desktops** | OS detection, file sharing ports | Windows, macOS, Linux workstations |
| **Network Infrastructure** | SNMP (161), routing protocols | Switches, Routers, Firewalls |

**Benefits:**
- No installation required on target devices
- Works with network equipment, IoT devices, and unmanaged endpoints
- Discovers unknown devices on the network
- Automatic scanning every 6 hours (configurable)
- Real-time network topology visualization

### 2. Software Vulnerability Detection

**Flow:**

```
┌──────────────┐
│  Go Agent    │
│  (Installed) │
└──────┬───────┘
       │
       ▼
┌─────────────────────────────────────┐
│  Software Discovery                 │
│  • Installed applications           │
│  • Package managers (npm, pip, etc) │
│  • System dependencies              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Dependency Extraction              │
│  • Parse package.json, requirements  │
│  • Extract versions and vendors      │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Go API Gateway                     │
│  • Batch processing                  │
│  • Queue management                  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Python Enrichment Service           │
│  ┌──────────────────────────────┐   │
│  │ Tier 1: CPE Guesser (L1)     │   │
│  │ • Exact match via Valkey     │   │
│  │ • 167K+ CPE entries          │   │
│  │ • Sub-millisecond lookup     │   │
│  └──────────────┬───────────────┘   │
│                 │                    │
│                 ▼                    │
│  ┌──────────────────────────────┐   │
│  │ Tier 2: Semantic Matcher (L2)│   │
│  │ • Vector embeddings          │   │
│  │ • pgvector similarity       │   │
│  │ • Handles name variations   │   │
│  └──────────────┬───────────────┘   │
│                 │                    │
│                 ▼                    │
│  ┌──────────────────────────────┐   │
│  │ CVE Database Lookup          │   │
│  │ • 320,407 CVEs (local)       │   │
│  │ • NVD API fallback           │   │
│  │ • KEV integration (1,468 CVEs) │   │
│  └──────────────────────────────┘   │
└──────────────┬───────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Risk Assessment & Scoring          │
│  • CVSS scoring                     │
│  • Exploit availability             │
│  • KEV prioritization              │
│  • Threat intelligence             │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  React Frontend                      │
│  • Real-time dashboard              │
│  • Vulnerability management         │
│  • Attack path visualization        │
└──────────────────────────────────────┘
```

**Capabilities:**
- Automatic discovery of installed applications (macOS, Linux, Windows)
- Dependency analysis for npm, pip, gem, cargo, go mod
- Real-time CVE matching with 320K+ CVEs
- KEV integration (1,468 known exploited CVEs)
- CVSS scoring and exploit availability tracking
- Vulnerability prioritization based on severity

### 3. Attack Path Visualization

**NodeZero-Style Exploit Chain Analysis:**

```
┌─────────────────────────────────────────────────────────────┐
│                    Attack Path Flow                         │
└─────────────────────────────────────────────────────────────┘

    [Entry Point]         [Step 1]              [Step 2]              [Step 3]           [Target]
    External Network  ──▶ Initial Access   ──▶ Lateral Move    ──▶ Privilege Escal  ──▶ Data Exfil
    (Green)                (Red)                (Orange)              (Yellow)            (Red)
    
    Features:
    • Multiple path visualization
    • Branching paths support
    • CVE details per step
    • MITRE ATT&CK technique mapping
    • Proof/evidence for exploits
    • Risk scoring (likelihood × impact)
```

**Features:**
- **Multiple Paths** - Visualize multiple attack paths simultaneously
- **Branching Support** - Shows shared attack points
- **Risk Scoring** - Likelihood, impact, and criticality scores
- **CVE Integration** - CVE IDs, proof, and mitigation controls
- **Color-Coded** - Visual risk indicators
- **Interactive** - Click nodes for detailed information

### 4. Hybrid CPE Matching System

**Two-Tier Architecture:**

```
Software Name: "nginx 1.18"
       │
       ▼
┌──────────────────────────────────────┐
│  Tier 1: CPE Guesser (L1)           │
│  • Valkey/Redis lookup               │
│  • 167,000+ CPE entries              │
│  • Exact match: "nginx:nginx:1.18"   │
│  • Speed: < 1ms                      │
└──────────────┬───────────────────────┘
               │ (if not found)
               ▼
┌──────────────────────────────────────┐
│  Tier 2: Semantic Matcher (L2)      │
│  • SentenceTransformers embeddings   │
│  • pgvector similarity search        │
│  • Handles: "nginx", "Nginx", "NGINX"│
│  • Version variations                │
│  • Speed: < 30ms                     │
└──────────────┬───────────────────────┘
               │ (if not found)
               ▼
┌──────────────────────────────────────┐
│  Fallback: NVD API                   │
│  • Real-time CVE lookup              │
│  • Rate-limited (6s) or fast (0.6s) │
│  • With API key                      │
└──────────────────────────────────────┘
```

**Performance:**
- **10,000x faster** than traditional NVD API lookups
- **95%+ match accuracy** with hybrid approach
- **1M+ applications/hour** processing capacity
- **320,407 CVEs** in local database

### 5. Threat Intelligence Integration

**KEV (Known Exploited Vulnerabilities):**

```
CVE Detection → KEV Check → Risk Elevation
     │              │              │
     ▼              ▼              ▼
  CVE-2021-12345  In KEV?    Critical Priority
                    │
                    ▼
              ┌─────────────────┐
              │ KEV Metadata    │
              │ • Date added    │
              │ • Required action│
              │ • Due date      │
              │ • Ransomware use│
              └─────────────────┘
```

**Integrated Sources:**
- **CISA KEV** - 1,468 known exploited CVEs (enabled by default)
- **MITRE ATT&CK** - Attack pattern mapping (optional)
- **AlienVault OTX** - Threat intelligence (optional)
- **OpenCVE** - Multi-source aggregation (optional)

### 6. Network Topology Visualization

**Interactive React Flow Diagram:**

```
┌─────────────────────────────────────────────────────────────┐
│              Network Topology View                          │
│                                                              │
│    [Router] ──▶ [Switch] ──▶ [Server]                      │
│       │            │            │                            │
│       │            ▼            ▼                            │
│       │        [IoT Device] [Database]                      │
│       │            │            │                            │
│       └────────────┴────────────┘                            │
│                    │                                         │
│                    ▼                                         │
│              [Vulnerabilities]                               │
│                                                              │
│  Features:                                                   │
│  • Device classification (color-coded)                     │
│  • Risk score visualization                                 │
│  • Interactive zoom, pan, search                             │
│  • Vulnerability mapping                                    │
│  • Export topology data                                     │
└─────────────────────────────────────────────────────────────┘
```

### 7. Configuration Auditing

**Nipper Studio-like Network Device Auditing:**

- Firewall rule analysis
- Default credential detection
- Insecure protocol identification
- Compliance checking (CIS, NIST, ISO 27001)
- Security posture scoring

---

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────────┐
│                         ZeroTrace Platform                      │
└─────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐
│   Agent 1    │    │   Agent 2    │    │   Agent N    │
│  (Sensor)    │    │  (Sensor)    │    │  (Sensor)    │
│              │    │              │    │              │
│ • Software   │    │ • Software   │    │ • Software   │
│   Scanner    │    │   Scanner    │    │   Scanner    │
│ • Network    │    │ • Network    │    │ • Network    │
│   Scanner    │    │   Scanner    │    │   Scanner    │
│ (Agentless)  │    │ (Agentless)  │    │ (Agentless)  │
└──────┬───────┘    └──────┬───────┘    └──────┬───────┘
       │                   │                   │
       └───────────────────┼───────────────────┘
                           │
                           ▼
              ┌────────────────────────┐
              │   Go API Gateway       │
              │   (Port 8080)          │
              │                        │
              │ • Multi-tenant Support │
              │ • Rate Limiting        │
              │ • Request Routing      │
              │ • Authentication       │
              └───────────┬────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
        ▼                 ▼                 ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  PostgreSQL  │  │  Valkey/     │  │   Queue      │
│   Database   │  │  Redis       │  │  Processor   │
│              │  │              │  │              │
│ • Assets     │  │ • CPE Cache  │  │ • Batch      │
│ • Scans      │  │ • Sessions   │  │   Processing │
│ • Vulns      │  │ • Queue      │  │ • Priority   │
│ • pgvector   │  │              │  │   Queue      │
│ • 320K CVEs  │  └──────────────┘  └──────┬───────┘
└──────┬───────┘                           │
       │                                    │
       └────────────────┬───────────────────┘
                        │
                        ▼
              ┌────────────────────────┐
              │   Python Enrichment    │
              │   Service (Port 8000)  │
              │                        │
              │ • CVE Lookup           │
              │ • CPE Matching        │
              │ • Batch Processing     │
              │ • AI Analysis          │
              │ • KEV Integration      │
              └───────────┬────────────┘
                          │
                          ▼
              ┌────────────────────────┐
              │   React Frontend       │
              │   (Port 5173)          │
              │                        │
              │ • Dashboard            │
              │ • Network Topology     │
              │ • Attack Paths         │
              │ • Vulnerability Mgmt  │
              │ • Real-time Updates   │
              └────────────────────────┘
```

### Data Flow Diagrams

**Software Scanning Flow:**
```
Agent Scanner 
    │
    ▼
Dependency Extraction (npm, pip, gem, cargo, go mod)
    │
    ▼
API Gateway (Batch Processing)
    │
    ▼
Queue Processor (Priority Queue)
    │
    ▼
Enrichment Service
    ├─▶ CPE Guesser (L1: Exact Match)
    ├─▶ Semantic Matcher (L2: Vector Search)
    └─▶ CVE Database (320K+ CVEs)
    │
    ▼
Risk Scoring (CVSS + KEV + Threat Intel)
    │
    ▼
Database Update
    │
    ▼
Real-time Frontend Update
```

**Network Scanning Flow:**
```
Agent Network Scanner (Sensor)
    │
    ▼
Network Discovery (Detect subnet: 192.168.1.0/24)
    │
    ▼
Nmap Scanning (All IPs in subnet)
    ├─▶ Port Scanning
    ├─▶ Service Detection
    ├─▶ OS Fingerprinting
    └─▶ Banner Grabbing
    │
    ▼
Device Classification
    ├─▶ Mobile Devices (Android/iOS)
    ├─▶ IoT Devices (MQTT, CoAP)
    ├─▶ Servers (SSH, HTTP, DB)
    ├─▶ Network Infrastructure (Switches, Routers)
    └─▶ Laptops/Desktops (Windows, macOS, Linux)
    │
    ▼
Nuclei Vulnerability Scanning
    │
    ▼
Configuration Auditing
    │
    ▼
API Gateway → Database
    │
    ▼
Network Topology Visualizer
    │
    ▼
Attack Path Analysis
```

**Enrichment Flow:**
```
Software List
    │
    ▼
CPE Guesser (L1) - Valkey/Redis
    │ (167K+ CPE entries)
    ├─▶ Found? → CVE Lookup
    └─▶ Not Found? → Continue
    │
    ▼
Semantic Matcher (L2) - pgvector
    │ (Vector embeddings)
    ├─▶ Found? → CVE Lookup
    └─▶ Not Found? → Continue
    │
    ▼
NVD API Fallback
    │ (Real-time lookup)
    ├─▶ With API Key: 0.6s delay
    └─▶ Without Key: 6s delay
    │
    ▼
CVE Database (PostgreSQL)
    │ (320,407 CVEs)
    ├─▶ KEV Check (1,468 CVEs)
    ├─▶ Threat Intel Enrichment
    └─▶ Risk Scoring
    │
    ▼
Results with Full Context
```

---

## Frontend Features

### Dashboard
- Real-time metrics and statistics
- Vulnerability trends and analytics
- Risk heatmaps
- Asset inventory overview

### Network Topology
- Interactive network visualization (React Flow)
- Color-coded device classification
- Search and filter capabilities
- Export topology data

### Attack Paths
- NodeZero-style exploit chain visualization
- Multiple paths and branching support
- Risk scoring and prioritization
- CVE details and proof per step

### Vulnerability Management
- Comprehensive vulnerability listing
- Advanced filtering and search
- Severity-based prioritization
- Export capabilities (JSON, CSV, PDF)

### AI Analytics
- Predictive vulnerability analysis
- Trend forecasting
- Exploit intelligence
- Remediation recommendations

---

## Quick Start

### Prerequisites

```bash
# Required Software
- Docker/Podman & Docker Compose
- Git
- 8GB+ RAM, 50GB+ storage
- Node.js 24.8.0+ (or Bun)
- Python 3.9+
- Go 1.21+
```

### Installation

```bash
# 1. Clone repository
git clone https://github.com/adhit-r/ZeroTrace.git
cd ZeroTrace

# 2. Set up environment variables
cp api-go/env.example api-go/.env
cp agent-go/env.example agent-go/.env
cp enrichment-python/env.example enrichment-python/.env
cp web-react/.env.example web-react/.env

# 3. Start services with Docker Compose
docker-compose up -d

# 4. Verify installation
curl http://localhost:8080/health
curl http://localhost:8000/health
```

### Local Development

```bash
# Start all services locally
./run-local.sh

# Services available at:
# - API: http://localhost:8080
# - Enrichment: http://localhost:8000
# - Frontend: http://localhost:5173
```

### Agent Setup

```bash
# 1. Navigate to agent directory
cd agent-go

# 2. Configure agent
cp env.example .env
# Edit .env with your API URL and enrollment token

# 3. Run agent
go run cmd/agent/main.go

# Agent will:
# - Register with API
# - Start software scanning
# - Begin network scanning (every 6 hours)
# - Send results to API
```

---

## Configuration

### Environment Variables

**API Configuration (`api-go/.env`):**
```bash
DATABASE_URL=postgresql://user:password@localhost:5432/zerotrace
REDIS_URL=redis://localhost:6379
API_PORT=8080
```

**Enrichment Configuration (`enrichment-python/.env`):**
```bash
NVD_API_KEY=your-nvd-api-key  # Optional, improves rate limits
DATABASE_URL=postgresql://user:password@localhost:5432/zerotrace
REDIS_URL=redis://localhost:6379
CISA_KEV_ENABLED=true  # Enabled by default
```

**Agent Configuration (`agent-go/.env`):**
```bash
API_URL=http://localhost:8080
NETWORK_SCAN_ENABLED=true
NETWORK_SCAN_INTERVAL=6h  # Default: 6 hours
```

**Frontend Configuration (`web-react/.env`):**
```bash
VITE_API_URL=http://localhost:8080
```

---

## Performance Metrics

| Metric | Target | Status |
|--------|--------|--------|
| API Response Time | < 100ms (95th percentile) | Yes |
| Enrichment Processing | < 30ms per application | Yes |
| Agent CPU Usage | < 5% average | Yes |
| Data Processing | 1M+ applications/hour | Yes |
| Concurrent Requests | 1000+ per second | Yes |
| CVE Database | 320,407 CVEs | Yes |
| KEV Integration | 1,468 CVEs | Yes |

---

## API Endpoints

### Agent Operations

```bash
# Register agent
POST /api/agents/register

# Send heartbeat
POST /api/agents/heartbeat

# Submit scan results
POST /api/agents/results

# Submit network scan results
POST /api/agents/network-scan-results

# List all agents
GET /api/agents/
```

### Network Scanning

```bash
# Initiate network scan
POST /api/v2/scans/network
{
  "agent_id": "uuid",
  "targets": ["192.168.1.0/24"],
  "scan_type": "tcp"
}

# Get scan status
GET /api/v2/scans/:scan_id/status

# Get scan results
GET /api/v2/scans/:scan_id/results
```

### Attack Paths

```bash
# Get all attack paths
GET /api/v2/attack-paths

# Get specific attack path
GET /api/v2/attack-paths/:path_id

# Generate attack paths
POST /api/v2/attack-paths/generate
```

### Vulnerabilities

```bash
# List vulnerabilities
GET /api/v2/vulnerabilities?severity=high&page=1

# Get vulnerability stats
GET /api/v2/vulnerabilities/stats

# Export vulnerabilities
GET /api/v2/vulnerabilities/export?format=csv
```

Full API documentation: [docs/openapi.yaml](docs/openapi.yaml)

---

## Project Structure

```
ZeroTrace/
├── api-go/                 # Go API server
│   ├── cmd/api/           # API entry point
│   ├── internal/          # Internal packages
│   │   ├── handlers/      # HTTP handlers
│   │   ├── services/      # Business logic
│   │   └── repository/   # Data access layer
│   └── migrations/       # Database migrations
│
├── agent-go/              # Go agent (sensor)
│   ├── cmd/agent/         # Agent binary
│   ├── internal/         # Internal packages
│   │   ├── scanner/      # Scanning modules
│   │   │   ├── network_scanner.go
│   │   │   └── device_classifier.go
│   │   └── communicator/ # API communication
│   └── mdm/              # MDM deployment files
│
├── enrichment-python/     # Python enrichment service
│   ├── app/              # FastAPI application
│   │   ├── services/     # Enrichment services
│   │   └── core/         # Configuration
│   └── scripts/          # CVE update scripts
│
├── web-react/            # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   │   └── network/  # Network visualization
│   │   ├── pages/        # Route pages
│   │   └── services/      # API integration
│
├── docs/                  # Documentation
│   ├── INDEX.md          # Documentation index
│   ├── KEV_INTEGRATION.md
│   └── openapi.yaml      # API specification
│
└── scripts/              # Utility scripts
```

---

## Documentation

### Getting Started
- [Documentation Index](docs/INDEX.md)
- [Architecture Overview](docs/architecture.md)
- [Database Schema](docs/DATABASE_SCHEMA_DIAGRAM.md)

### Features
- [KEV Integration](docs/KEV_INTEGRATION.md) - CISA Known Exploited Vulnerabilities
- [Config Auditor](docs/CONFIG_AUDITOR_NIPPER_STUDIO.md) - Network device auditing

### Component READMEs
- [API Service](api-go/README.md)
- [Enrichment Service](enrichment-python/README.md)
- [Agent Service](agent-go/README.md)
- [Frontend](web-react/README.md)

---

## Technology Stack

### Backend
- **Go API**: Gin framework, PostgreSQL, Redis/Valkey
- **Python Enrichment**: FastAPI, pgvector, SentenceTransformers
- **Agent**: Go with Nmap, Naabu, Nuclei integration

### Frontend
- **React 19** with TypeScript
- **Vite 7** for optimized builds
- **Tailwind CSS** with neobrutalist design
- **React Flow** for network visualization
- **React Query** for state management

### Infrastructure
- **PostgreSQL 15+** with pgvector extension
- **Redis/Valkey** for caching
- **Docker/Podman** for containerization
- **Prometheus + Grafana** for monitoring

---

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development setup
- Code standards
- Pull request process
- Issue reporting

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Repository**: https://github.com/adhit-r/ZeroTrace  
**Maintainer**: [adhit-r](https://github.com/adhit-r)

Made for enterprise security teams

</div>
