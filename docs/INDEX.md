# ZeroTrace Documentation Index

Comprehensive documentation for the ZeroTrace vulnerability detection and management platform.

## Table of Contents

### Getting Started

- [Main README](../README.md) - Project overview, architecture, quick start
- [Quick Start Guide](QUICK_START.md) - Get up and running in minutes
- [ROADMAP](../ROADMAP.md) - Development roadmap and milestones
- [CHANGELOG](../CHANGELOG.md) - Version history and changes
- [CONTRIBUTING](../CONTRIBUTING.md) - Contribution guidelines
- [Development Setup](development-setup.md) - Development environment setup

### Architecture & Design

- [Architecture Overview](architecture.md) - System architecture and design
- [Technology Stack](tech-stack.md) - Technologies and frameworks used
- [Database Schema](database-schema.md) - Database structure (v1)
- [Database Schema v2](database-schema-v2.md) - Enhanced database structure
- [API Endpoints](api-endpoints.md) - API v1 endpoint documentation
- [API v2 Documentation](api-v2-documentation.md) - Enhanced API v2 endpoints
- [OpenAPI Specification](openapi.yaml) - Complete OpenAPI/Swagger specification

### Components

#### API Service

- [API Service README](../api-go/README.md) - API service overview
- [Go API Documentation](go-api.md) - Detailed API implementation
- [API Endpoints](api-endpoints.md) - Endpoint reference
- [API v2 Documentation](api-v2-documentation.md) - Enhanced API features

#### Enrichment Service

- [Enrichment Service README](../enrichment-python/README.md) - Enrichment service overview
- [Python Enrichment Documentation](python-enrichment.md) - Enrichment implementation
- [Agent CVE Documentation](agent-cve.md) - CVE detection and enrichment

#### Agent Service

- [Agent Service README](../agent-go/README.md) - Agent service overview
- [Go Agent Documentation](go-agent.md) - Agent implementation
- [Scanner Modules Documentation](scanner-modules-documentation.md) - Scanner modules
- [Network Discovery Implementation](network-discovery-implementation.md) - Network scanning
- [MDM Agent Summary](../agent-go/MDM_AGENT_SUMMARY.md) - MDM integration

#### Frontend

- [Frontend README](../web-react/README.md) - Frontend overview
- [Web Implementation](web-implementation.md) - Frontend architecture
- [Frontend Technology Analysis](frontend-technology-analysis.md) - Technology choices

### Features & Capabilities

- [Compliance Frameworks](compliance-frameworks-documentation.md) - Compliance support
- [Organization Prioritization Design](org-prioritization-design.md) - Risk prioritization
- [Scalable Data Processing](scalable-data-processing.md) - Performance optimization
- [Monitoring Strategy](monitoring-strategy.md) - Monitoring and observability
- [Performance Monitoring Setup](performance/monitoring-setup.md) - Performance monitoring

### Deployment & Operations

- [Docker Compose](../docker-compose.yml) - Development deployment
- [Monitoring Setup](performance/monitoring-setup.md) - Production monitoring
- [MDM Agent Summary](../agent-go/MDM_AGENT_SUMMARY.md) - MDM deployment
- [GitHub Actions Workflows](../.github/workflows/README.md) - CI/CD pipelines

### Wiki

- [Wiki Home](../wiki/Home.md) - Wiki overview
- [Architecture Overview (Wiki)](../wiki/Architecture-Overview.md) - Architecture wiki
- [Installation Guide (Wiki)](../wiki/Installation-Guide.md) - Installation wiki
- [Contributing Guidelines (Wiki)](../wiki/Contributing-Guidelines.md) - Contributing wiki

## Quick Reference

### Component Quick Starts

- [API Quick Start](../api-go/README.md#quick-start)
- [Enrichment Quick Start](../enrichment-python/README.md#quick-start)
- [Agent Quick Start](../agent-go/README.md#quick-start)
- [Frontend Quick Start](../web-react/README.md#quick-start)

### API Quick Reference

- **Health Check**: `GET /health`
- **Agent Registration**: `POST /api/agents/register`
- **Agent Heartbeat**: `POST /api/agents/heartbeat`
- **Agent Results**: `POST /api/agents/results`
- **Vulnerabilities v2**: `GET /api/v2/vulnerabilities`
- **Dashboard**: `GET /api/dashboard/overview`

See [API v2 Documentation](api-v2-documentation.md) for complete endpoint reference.

### Configuration Files

- [API Environment Example](../api-go/env.example)
- [Enrichment Environment Example](../enrichment-python/env.example)
- [Frontend Environment Example](../web-react/.env.example)
- [Docker Compose](../docker-compose.yml)

## Documentation Standards

All documentation follows these standards:

- **Professional Format**: No emojis, consistent headers
- **Code Examples**: Complete, runnable examples
- **Cross-References**: Links between related documents
- **Last Updated**: Dates maintained on key documents
- **Versioning**: Clear version indicators for API docs

## Contributing to Documentation

See [CONTRIBUTING.md](../CONTRIBUTING.md) for documentation contribution guidelines.

---

**Last Updated**: January 2025

