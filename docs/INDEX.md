# ZeroTrace Documentation Index

Essential documentation for the ZeroTrace vulnerability detection and management platform.

## Core Documentation

### Getting Started
- [Main README](../README.md) - Project overview, architecture, quick start
- [ROADMAP](../ROADMAP.md) - Development roadmap and milestones
- [CHANGELOG](../CHANGELOG.md) - Version history and changes
- [CONTRIBUTING](../CONTRIBUTING.md) - Contribution guidelines

### Architecture & Design
- [Architecture Overview](architecture.md) - System architecture and design
- [Database Schema Diagram](DATABASE_SCHEMA_DIAGRAM.md) - Visual database schema with relationships
- [OpenAPI Specification](openapi.yaml) - Complete OpenAPI/Swagger specification

### Features
- [KEV Integration](KEV_INTEGRATION.md) - CISA Known Exploited Vulnerabilities integration
- [Configuration Auditor](CONFIG_AUDITOR_NIPPER_STUDIO.md) - Firewall/network device config auditing (Nipper Studio-like)

### Component READMEs
- [API Service](../api-go/README.md) - API service overview
- [Enrichment Service](../enrichment-python/README.md) - Enrichment service overview
- [Agent Service](../agent-go/README.md) - Agent service overview
- [Frontend](../web-react/README.md) - Frontend overview

## Quick Reference

### API Endpoints
- **Health Check**: `GET /health`
- **Agent Registration**: `POST /api/agents/register`
- **Agent Heartbeat**: `POST /api/agents/heartbeat`
- **Agent Results**: `POST /api/agents/results`
- **Vulnerabilities v2**: `GET /api/v2/vulnerabilities`
- **Dashboard**: `GET /api/dashboard/overview`

### Configuration Files
- [API Environment Example](../api-go/env.example)
- [Enrichment Environment Example](../enrichment-python/env.example)
- [Docker Compose](../docker-compose.yml)

## Archived Documentation

Detailed implementation documentation, component-specific guides, and historical documentation have been moved to `archive/docs/` for reference.

**Last Updated**: December 2025
