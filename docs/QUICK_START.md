# ZeroTrace Quick Start Guide

Get up and running with ZeroTrace in minutes.

## Prerequisites

- Docker and Docker Compose (recommended)
- OR: Go 1.21+, Python 3.9+, Node.js 24.8.0+, PostgreSQL 15+, Redis 7+

## Option 1: Docker Compose (Recommended)

### Step 1: Clone Repository

```bash
git clone https://github.com/adhit-r/ZeroTrace.git
cd ZeroTrace
```

### Step 2: Start Services

```bash
# Start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### Step 3: Access Services

- **API**: http://localhost:8080
- **Enrichment Service**: http://localhost:8000
- **Frontend**: http://localhost:3000
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3001

### Step 4: Verify Installation

```bash
# Check API health
curl http://localhost:8080/health

# Check enrichment service
curl http://localhost:8000/health

# Check dashboard
curl http://localhost:8080/api/dashboard/overview
```

## Option 2: Manual Setup

### Step 1: Start Database and Redis

```bash
# Start PostgreSQL
docker run -d \
  --name zerotrace-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=zerotrace \
  -p 5432:5432 \
  postgres:15

# Start Redis
docker run -d \
  --name zerotrace-redis \
  -p 6379:6379 \
  redis:7
```

### Step 2: Set Up API Service

```bash
cd api-go

# Install dependencies
go mod tidy

# Copy environment template
cp env.example .env

# Edit .env with your configuration
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/zerotrace?sslmode=disable
# REDIS_URL=redis://localhost:6379

# Run migrations (if needed)
# go run migrations/migrate.go

# Start API
go run cmd/api/main.go
```

### Step 3: Set Up Enrichment Service

```bash
cd enrichment-python

# Install dependencies with uv (recommended)
uv pip install -r requirements.txt

# Copy environment template
cp env.example .env

# Edit .env with your configuration
# NVD_API_KEY=your-nvd-api-key (optional but recommended)

# Start enrichment service
uv run python app/main.py
```

### Step 4: Set Up Frontend

```bash
cd web-react

# Install dependencies with Bun (recommended)
bun install

# Copy environment template
cp .env.example .env

# Edit .env with your configuration
# VITE_API_URL=http://localhost:8080

# Start development server
bun run dev
```

### Step 5: Set Up Agent (Optional)

```bash
cd agent-go

# Install dependencies
go mod tidy

# Copy environment template
cp env.example .env

# Edit .env with your configuration
# API_URL=http://localhost:8080
# ENROLLMENT_TOKEN=your-enrollment-token

# Run agent (development mode with tray)
go run cmd/agent/main.go
```

## First Steps

### 1. Register an Agent

```bash
curl -X POST http://localhost:8080/api/agents/register \
  -H "Content-Type: application/json" \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "organization_id": "660e8400-e29b-41d4-a716-446655440001",
    "name": "my-agent",
    "version": "1.0.0",
    "hostname": "my-computer",
    "os": "Linux"
  }'
```

### 2. Send Heartbeat

```bash
curl -X POST http://localhost:8080/api/agents/heartbeat \
  -H "Content-Type: application/json" \
  -d '{
    "agent_id": "550e8400-e29b-41d4-a716-446655440000",
    "organization_id": "660e8400-e29b-41d4-a716-446655440001",
    "agent_name": "my-agent",
    "status": "active",
    "cpu_usage": 45.5,
    "memory_usage": 62.3,
    "timestamp": "2025-01-15T10:30:00Z"
  }'
```

### 3. View Dashboard

Open http://localhost:3000 in your browser to see the dashboard.

### 4. Check Vulnerabilities

```bash
# Get vulnerabilities
curl "http://localhost:8080/api/v2/vulnerabilities?severity=high"

# Get vulnerability statistics
curl http://localhost:8080/api/v2/vulnerabilities/stats
```

## Configuration

### Environment Variables

#### API Service (`api-go/.env`)

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:5432/zerotrace?sslmode=disable
REDIS_URL=redis://localhost:6379
API_PORT=8080
API_HOST=0.0.0.0
LOG_LEVEL=info
```

#### Enrichment Service (`enrichment-python/.env`)

```bash
ENRICHMENT_PORT=8000
NVD_API_KEY=your-nvd-api-key
DATABASE_URL=postgres://postgres:postgres@localhost:5432/zerotrace?sslmode=disable
REDIS_URL=redis://localhost:6379
LOG_LEVEL=info
```

#### Frontend (`web-react/.env`)

```bash
VITE_API_URL=http://localhost:8080
VITE_ENRICHMENT_URL=http://localhost:8000
VITE_ENV=development
```

#### Agent (`agent-go/.env`)

```bash
API_URL=http://localhost:8080
ENROLLMENT_TOKEN=your-enrollment-token
ORGANIZATION_ID=your-organization-id
SCAN_INTERVAL=5m
LOG_LEVEL=info
```

## Troubleshooting

### API Not Starting

```bash
# Check if port 8080 is available
lsof -i :8080

# Check database connection
psql postgres://postgres:postgres@localhost:5432/zerotrace

# Check Redis connection
redis-cli ping
```

### Enrichment Service Not Starting

```bash
# Check if port 8000 is available
lsof -i :8000

# Verify Python dependencies
uv pip list
```

### Frontend Not Loading

```bash
# Check if port 3000 is available
lsof -i :3000

# Verify dependencies
bun install

# Check environment variables
cat web-react/.env
```

### Agent Not Connecting

```bash
# Verify API URL
curl http://localhost:8080/health

# Check enrollment token
echo $ENROLLMENT_TOKEN

# View agent logs
tail -f agent.log
```

## Next Steps

1. **Read Documentation**: See [Documentation Index](INDEX.md)
2. **Explore API**: See [API v2 Documentation](api-v2-documentation.md)
3. **Deploy Agent**: See [Agent README](../agent-go/README.md)
4. **Configure Monitoring**: See [Monitoring Setup](performance/monitoring-setup.md)

## Component Quick Starts

- [API Quick Start](../api-go/README.md#quick-start)
- [Enrichment Quick Start](../enrichment-python/README.md#quick-start)
- [Frontend Quick Start](../web-react/README.md#quick-start)
- [Agent Quick Start](../agent-go/README.md#quick-start)

## Getting Help

- **Documentation**: [docs/INDEX.md](INDEX.md)
- **Issues**: [GitHub Issues](https://github.com/adhit-r/ZeroTrace/issues)
- **Contributing**: [CONTRIBUTING.md](../CONTRIBUTING.md)

---

**Last Updated**: January 2025

