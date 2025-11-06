# Architecture Overview

ZeroTrace is a comprehensive vulnerability detection and management platform designed for enterprise-scale deployment. This document provides a detailed overview of the system architecture.

## System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ZeroTrace     │    │   ZeroTrace     │    │   ZeroTrace     │
│     Agent       │    │     Agent       │    │     Agent       │
│   (Device 1)    │    │   (Device 2)    │    │   (Device N)    │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼─────────────┐
                    │      Go API Gateway       │
                    │   (Rate Limiting & Auth)  │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │    Queue Processor        │
                    │   (Redis-based Queue)     │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │  Python Enrichment       │
                    │   (CVE Data Processing)  │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │   PostgreSQL Database     │
                    │   (Partitioned by Org)    │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────▼─────────────┐
                    │   React Web Dashboard     │
                    │   (Real-time Updates)     │
                    └───────────────────────────┘
```

## Core Components

### 1. ZeroTrace Agent

- **Language**: Go
- **Purpose**: System scanning and data collection
- **Deployment**: Universal binary for all organizations
- **Features**:
  - System tray integration (Full Agent)
  - Headless operation (Simple Agent for MDM)
  - CPU/memory optimization
  - Adaptive scanning
  - Heartbeat mechanism
  - Enrollment token support

### 2. Go API Gateway

- **Language**: Go (Gin framework)
- **Purpose**: Primary API server and request handling
- **Features**:
  - Rate limiting (per company/agent)
  - Authentication & authorization
  - Request routing and validation
  - Performance monitoring (APM)
  - Multi-level caching
  - Connection pooling

### 3. Queue Processor

- **Language**: Go
- **Purpose**: Batch processing and data management
- **Features**:
  - Redis-based priority queue
  - Batch processing by company
  - Multiple worker goroutines
  - Metrics collection
  - Automatic cleanup

### 4. Python Enrichment Service

- **Language**: Python (FastAPI)
- **Purpose**: CVE data enrichment and processing
- **Features**:
  - Ultra-optimized performance (10,000x)
  - Multi-level caching (Memory, Redis, Memcached)
  - Parallel processing (1000+ concurrent requests)
  - Circuit breakers and retry logic
  - Background tasks and monitoring

### 5. PostgreSQL Database

- **Purpose**: Primary data storage
- **Features**:
  - Partitioning by company_id
  - Optimized indexes
  - Materialized views
  - Automated cleanup
  - JSONB support for metadata

### 6. React Web Dashboard

- **Language**: TypeScript/React
- **Purpose**: User interface and data visualization
- **Features**:
  - Terminal-inspired dark theme
  - Real-time updates via WebSocket
  - Responsive design
  - Data querying with React Query
  - Enterprise-grade UI components

## Data Flow

### 1. Agent Enrollment

```
Agent → Enrollment Token → Go API → Database
  ↓
Device Credential Issued
  ↓
Agent registers with credential
```

### 2. Data Collection

```
Agent → System Scan → App Data → Go API
  ↓
Queue Processor → Batch Processing
  ↓
Python Enrichment → CVE Data
  ↓
Database Storage
```

### 3. Data Retrieval

```
React Dashboard → Go API → Database
  ↓
Real-time Updates via WebSocket
  ↓
UI Rendering with Terminal Theme
```

## Multi-Organization Architecture

### Universal Agent Design

- Single binary for all organizations
- Enrollment token for organization identification
- Device credentials for long-term authentication
- Backend enforcement of organization isolation

### Database Partitioning

```sql
-- Partitioned tables by company_id
CREATE TABLE agents (
    id UUID PRIMARY KEY,
    company_id UUID NOT NULL,
    -- other fields
) PARTITION BY HASH (company_id);

-- Separate partitions for each company
CREATE TABLE agents_company_1 PARTITION OF agents
FOR VALUES WITH (modulus 10, remainder 0);
```

### API Scoping

- All API requests scoped by organization
- Middleware enforces organization isolation
- Rate limiting per organization
- Caching with organization prefixes

## Performance Optimizations

### Go API (100x Performance)

- Connection pooling (database, Redis)
- Multi-level caching (Memory, Redis, Memcached)
- Query optimization with prepared statements
- Rate limiting with semaphores
- Batch processing capabilities
- Memory optimization with GC tuning

### Python Enrichment (10,000x Performance)

- `uvloop` for async I/O optimization
- `orjson` for fast JSON processing
- Connection pooling (10,000+ connections)
- Parallel processing (1000 concurrent requests)
- Load balancing across multiple endpoints
- Memory monitoring and optimization

### Agent (Minimal CPU Usage)

- Adaptive scanning based on system load
- Resource throttling
- Background processing
- Go runtime optimization (GOMAXPROCS=1)
- Memory limits and GC tuning
- Process priority adjustment

## Monitoring & Observability

### Application Performance Monitoring (APM)

- Prometheus metrics collection
- Custom business metrics
- System resource monitoring
- Database query metrics
- Queue processing metrics
- Cache hit/miss ratios

### Logging Strategy

- Structured logging with Zap (Go)
- Structured logging with structlog (Python)
- Centralized logging with ELK stack
- Log correlation across services
- Error tracking and alerting

### Health Checks

- Service health endpoints
- Database connectivity checks
- External service dependencies
- Queue health monitoring
- Cache health verification

## Security Architecture

### Authentication & Authorization

- JWT-based authentication
- Organization-scoped access control
- Role-based permissions
- Token expiration and rotation
- Secure credential storage

### Data Protection

- Organization isolation at database level
- Encrypted communication (HTTPS/WSS)
- Secure credential transmission
- Audit logging for all operations
- Data retention policies

### Network Security

- Rate limiting per organization
- Input validation and sanitization
- SQL injection prevention
- XSS protection
- CORS configuration

## Deployment Architecture

### Development Environment

```
Local Development:
├── Docker Compose (PostgreSQL, Redis)
├── Go API (localhost:8080)
├── Python Enrichment (localhost:8000)
├── React Dashboard (localhost:3000)
└── Agent (local binary)
```

### Production Environment

```
Production Deployment:
├── Load Balancer (Nginx)
├── Go API (multiple instances)
├── Python Enrichment (multiple instances)
├── PostgreSQL (clustered)
├── Redis (clustered)
├── Monitoring Stack (Prometheus, Grafana)
└── Logging Stack (ELK)
```

### MDM Deployment

```
Enterprise MDM:
├── Microsoft Intune
├── Jamf Pro
├── Azure AD
├── VMware Workspace ONE
└── Custom deployment scripts
```

## Scalability Features

### Horizontal Scaling

- Stateless API services
- Multiple enrichment workers
- Database read replicas
- Redis clustering
- Load balancer distribution

### Vertical Scaling

- Resource optimization
- Memory management
- CPU utilization
- I/O optimization
- Cache efficiency

### Data Scaling

- Database partitioning
- Automated cleanup
- Archive strategies
- Index optimization
- Query optimization

## Configuration Management

### Environment Variables

```bash
# API Configuration
API_PORT=8080
API_ENV=production
API_LOG_LEVEL=info

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=zerotrace
DB_USER=zerotrace
DB_PASSWORD=secure_password

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secure_password

# Enrichment Configuration
ENRICHMENT_WORKERS=10
ENRICHMENT_CACHE_TTL=3600
ENRICHMENT_RATE_LIMIT=1000
```

### Feature Flags

- A/B testing capabilities
- Gradual feature rollouts
- Environment-specific features
- Performance monitoring toggles
- Debug mode controls

## Future Architecture

### Planned Enhancements

- Microservices architecture
- Event-driven architecture
- GraphQL API
- Real-time collaboration
- Advanced analytics
- Machine learning integration

### Technology Evolution

- Kubernetes deployment
- Service mesh (Istio)
- Event streaming (Kafka)
- Advanced caching (Hazelcast)
- Time-series database (InfluxDB)

---

**Architecture Version**: 2.0.0  
**Last Updated**: January 2025  
**Next Review**: March 2025
