# ZeroTrace API

The core backend API for the ZeroTrace vulnerability scanning platform.

## Overview

ZeroTrace API is a high-performance REST API built with Go and Gin framework. It provides comprehensive endpoints for vulnerability management, agent communication, dashboard data, and enterprise features.

## Features

- **Authentication**: Clerk-based authentication with role-based access control
- **Scan Management**: Create, read, update, and delete vulnerability scans
- **Agent Management**: Agent registration, enrollment, and heartbeat tracking
- **Vulnerability Tracking**: Comprehensive vulnerability detection and management
- **Dashboard Data**: Real-time vulnerability statistics and trends
- **Multi-tenant Support**: Organization isolation with secure data access
- **AI-Powered Analysis**: Advanced vulnerability analysis and remediation guidance
- **Compliance Reporting**: Automated compliance framework assessments
- **High Performance**: Optimized for handling 100,000+ data points

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Installation

1. **Clone and navigate to the API directory:**
   ```bash
   cd api-go
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Run the API:**
   ```bash
   go run cmd/api/main.go
   ```

The API will start on `http://localhost:8080`

## Configuration

### Environment Variables

See `env.example` for all available configuration options.

#### Required

- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: JWT signing key (or `CLERK_JWT_VERIFICATION_KEY` for Clerk)

#### Optional

- `API_PORT`: Server port (default: 8080)
- `API_HOST`: Server host (default: 0.0.0.0)
- `API_MODE`: Debug mode (default: debug)
- `LOG_LEVEL`: Logging level (default: info)
- `RATE_LIMIT_REQUESTS`: Rate limit requests per window (default: 100)
- `RATE_LIMIT_WINDOW`: Rate limit window (default: 1m)

## API Endpoints

### Health Check

- `GET /health` - API health status

**Example Request:**
```bash
curl http://localhost:8080/health
```

**Example Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

### Agent Operations

- `POST /api/agents/register` - Register new agent
- `POST /api/agents/heartbeat` - Send agent heartbeat
- `POST /api/agents/results` - Submit scan results
- `POST /api/agents/system-info` - Update system information
- `GET /api/agents` - List all agents
- `GET /api/agents/online` - Get online agents
- `GET /api/agents/stats` - Get agent statistics

**Example: Register Agent**
```bash
curl -X POST http://localhost:8080/api/agents/register \
  -H "Content-Type: application/json" \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "organization_id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "agent-001",
    "version": "1.0.0",
    "hostname": "server-01",
    "os": "Linux"
  }'
```

**Example: Send Heartbeat**
```bash
curl -X POST http://localhost:8080/api/agents/heartbeat \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "550e8400-e29b-41d4-a716-446655440000",
    "organization_id": "660e8400-e29b-41d4-a716-446655440001",
    "agent_name": "agent-001",
    "status": "active",
    "cpu_usage": 45.5,
    "memory_usage": 62.3,
    "timestamp": "2025-01-15T10:30:00Z"
  }'
```

**Example: Get All Agents**
```bash
curl http://localhost:8080/api/agents
```

### Vulnerabilities

- `GET /api/vulnerabilities` - List vulnerabilities
- `GET /api/v2/vulnerabilities` - List vulnerabilities (v2)
- `GET /api/v2/vulnerabilities/stats` - Get vulnerability statistics
- `GET /api/v2/vulnerabilities/export` - Export vulnerabilities

**Example: Get Vulnerabilities (v2)**
```bash
curl "http://localhost:8080/api/v2/vulnerabilities?severity=high&page=1&page_size=20"
```

**Example: Get Vulnerability Statistics**
```bash
curl http://localhost:8080/api/v2/vulnerabilities/stats
```

**Example Response:**
```json
{
  "total": 150,
  "by_severity": {
    "critical": 5,
    "high": 25,
    "medium": 80,
    "low": 40
  },
  "by_category": {
    "network": 30,
    "database": 20,
    "api": 10
  },
  "compliance_score": 85.5,
  "risk_score": 7.2
}
```

### Dashboard

- `GET /api/dashboard/overview` - Get dashboard overview
- `GET /api/v1/dashboard/overview` - Get protected dashboard overview
- `GET /api/v1/dashboard/trends` - Get vulnerability trends

**Example: Get Dashboard Overview**
```bash
curl http://localhost:8080/api/dashboard/overview
```

**Example Response:**
```json
{
  "success": true,
  "data": {
    "assets": {
      "total": 100,
      "vulnerable": 45,
      "critical": 5,
      "high": 25,
      "medium": 80,
      "low": 40
    },
    "vulnerabilities": {
      "total": 150,
      "critical": 5,
      "high": 25,
      "medium": 80,
      "low": 40
    },
    "agents": {
      "total": 100,
      "online": 85
    }
  }
}
```

### AI Analysis

- `GET /api/ai-analysis/vulnerabilities/:id/comprehensive` - Comprehensive vulnerability analysis
- `GET /api/ai-analysis/vulnerabilities/trends` - Analyze vulnerability trends
- `GET /api/ai-analysis/exploit-intelligence/:cve_id` - Get exploit intelligence
- `GET /api/ai-analysis/vulnerabilities/:id/predictive` - Predictive analysis
- `GET /api/ai-analysis/vulnerabilities/:id/remediation-plan` - Get remediation plan
- `POST /api/ai-analysis/bulk-analysis` - Bulk vulnerability analysis

### Compliance

- `GET /api/compliance/organizations/:id/report` - Generate compliance report
- `GET /api/compliance/organizations/:id/score` - Get compliance score
- `GET /api/compliance/organizations/:id/findings` - Get compliance findings
- `GET /api/v2/compliance/status` - Get compliance status

### Organization Profile

- `POST /api/organizations/profile` - Create organization profile
- `GET /api/organizations/:id/profile` - Get organization profile
- `PUT /api/organizations/:id/profile` - Update organization profile
- `DELETE /api/organizations/:id/profile` - Delete organization profile

### Enrollment

- `POST /api/enrollment/enroll` - Enroll agent
- `POST /api/v1/enrollment/tokens` - Generate enrollment token (protected)
- `DELETE /api/v1/enrollment/tokens/:id` - Revoke enrollment token (protected)

Full API documentation available in `/docs/api-v2-documentation.md`

## Development

### Project Structure

```
api-go/
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── handlers/      # HTTP request handlers
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Data models
│   ├── services/      # Business logic
│   ├── repository/    # Data access layer
│   ├── monitoring/    # APM system
│   ├── optimization/  # Performance optimizations
│   └── queue/         # Queue processing
├── migrations/         # Database migrations
└── tests/             # Test files
```

### Adding New Endpoints

1. **Create handler** in `internal/handlers/`
2. **Add service logic** in `internal/services/`
3. **Register route** in `cmd/api/main.go`
4. **Add tests** in `tests/`

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./internal/handlers/ -v

# Run with coverage
go test ./... -cover
```

## Performance

### Optimizations

- **Multi-level Caching**: Memory + Redis + Memcached
- **Connection Pooling**: Optimized database connections
- **Request Batching**: Bulk operations for efficiency
- **Query Optimization**: Indexed database queries
- **100x Performance**: Comprehensive caching improvements

### Performance Metrics

- **API Response Time**: < 100ms (95th percentile)
- **Throughput**: 1000+ requests/second
- **Database Operations**: 10,000+ operations/second
- **Cache Hit Rate**: > 80%

## Monitoring

### Health Checks

- `GET /health` - Service health status

### Metrics

Metrics available via Prometheus:

- `zerotrace_api_requests_total` - Total API requests
- `zerotrace_api_request_duration_seconds` - Request latency
- `zerotrace_vulnerabilities_total` - Total vulnerabilities
- `zerotrace_assets_total` - Total assets scanned

## Deployment

### Docker

```bash
# Build image
docker build -t zerotrace-api .

# Run container
docker run -p 8080:8080 --env-file .env zerotrace-api
```

### Docker Compose

```bash
# Start with docker-compose
docker-compose up api
```

### Production

```bash
# Build for production
go build -o zerotrace-api cmd/api/main.go

# Run with systemd or supervisor
./zerotrace-api
```

## Authentication

ZeroTrace API uses Clerk for production-grade authentication:

- **Multi-Organization Support**: Secure multi-tenant architecture
- **Role-Based Access Control**: Granular permission system
- **JWT Tokens**: Secure token-based authentication
- **SSO Integration**: Enterprise identity provider support

### Clerk Setup

1. Create a Clerk account at [clerk.com](https://clerk.com)
2. Create a new application
3. Set `CLERK_JWT_VERIFICATION_KEY` in API environment
4. Configure webhook endpoints for user/org changes

## Documentation

- [API v2 Documentation](../docs/api-v2-documentation.md)
- [Go API Documentation](../docs/go-api.md)
- [Architecture Documentation](../docs/architecture.md)
- [Development Setup](../docs/development-setup.md)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.

---

**Last Updated**: January 2025
