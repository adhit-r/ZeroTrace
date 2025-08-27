# ZeroTrace Tech Stack

## Overview
ZeroTrace uses a modern, high-performance tech stack optimized for scalability, security, and developer productivity. All components are designed to run locally during development.

## Backend Technologies

### Go (Agent & API)
- **Version**: Go 1.21+
- **Framework**: Gin (HTTP framework)
- **Database**: GORM (ORM)
- **Authentication**: JWT-Go
- **Validation**: Go-playground/validator
- **Testing**: Testify
- **Configuration**: Viper
- **Logging**: Logrus
- **HTTP Client**: Resty
- **WebSocket**: Gorilla WebSocket

### Python (Enrichment)
- **Version**: Python 3.11+
- **Framework**: FastAPI
- **Async**: asyncio, aiohttp
- **Database**: SQLAlchemy (ORM)
- **Security**: PyJWT, passlib
- **Testing**: pytest, pytest-asyncio
- **Data Processing**: pandas, numpy
- **ML/AI**: scikit-learn, tensorflow
- **Vulnerability Analysis**: safety, bandit
- **Logging**: structlog

## Frontend Technologies

### React
- **Version**: React 18+
- **Build Tool**: Vite
- **Language**: TypeScript
- **State Management**: Zustand
- **Routing**: React Router v6
- **HTTP Client**: Axios
- **UI Components**: 
  - Headless UI
  - Tailwind CSS
  - Heroicons
- **Charts**: Recharts
- **Forms**: React Hook Form
- **Validation**: Zod
- **Testing**: Vitest, React Testing Library

## Database & Storage

### Primary Database
- **PostgreSQL 15+**
  - Multi-tenant architecture
  - JSONB for flexible data
  - Full-text search
  - Partitioning for large tables
  - Connection pooling

### Caching & Sessions
- **Redis 7+**
  - Session storage
  - API response caching
  - Job queues
  - Real-time data

### Local Storage
- **File System**
  - Scan results
  - Logs
  - Temporary files
  - Configuration backups

## Development Tools

### Containerization
- **Podman** (as per user preference)
- **Docker Compose** for local development
- **Multi-stage builds** for optimization

### Package Management
- **Go**: Go modules
- **Python**: Poetry
- **Node.js**: Bun (as per user preference)
- **System**: Homebrew (macOS)

### Development Environment
- **IDE**: VS Code with extensions
- **API Testing**: Insomnia
- **Database**: pgAdmin/DBeaver
- **Redis**: RedisInsight

## Testing & Quality

### Testing Frameworks
- **Go**: Testify, GoConvey
- **Python**: pytest, pytest-cov
- **Frontend**: Vitest, React Testing Library
- **E2E**: Playwright

### Code Quality
- **Linting**: 
  - Go: golangci-lint
  - Python: flake8, black, isort
  - TypeScript: ESLint, Prettier
- **Security**: 
  - Go: gosec
  - Python: bandit, safety
  - Frontend: npm audit

### Documentation
- **API**: OpenAPI/Swagger
- **Code**: GoDoc, Sphinx
- **Architecture**: Mermaid diagrams

## Performance & Monitoring

### Performance Tools
- **Profiling**: pprof (Go), cProfile (Python)
- **Benchmarking**: Go benchmarks, pytest-benchmark
- **Load Testing**: k6, Apache Bench

### Monitoring (Local)
- **Metrics**: Prometheus
- **Visualization**: Grafana
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Tracing**: Jaeger

## Security Stack

### Authentication & Authorization
- **JWT**: JSON Web Tokens
- **OAuth2**: For external integrations
- **RBAC**: Role-based access control
- **API Keys**: For agent authentication

### Security Tools
- **Dependency Scanning**: 
  - Go: nancy
  - Python: safety
  - Node.js: npm audit
- **Code Scanning**: 
  - Go: gosec
  - Python: bandit
  - TypeScript: ESLint security rules

### Encryption & Security
- **TLS**: HTTPS/WSS
- **Hashing**: bcrypt, Argon2
- **Encryption**: AES-256
- **Secrets**: Environment variables, .env files

## Local Development Setup

### Prerequisites
```bash
# Required software
- Go 1.21+
- Python 3.11+
- Node.js 18+ (or Bun)
- Podman
- PostgreSQL 15+
- Redis 7+
```

### Development Dependencies
```bash
# Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

# Python tools
pip install poetry
poetry install

# Node.js tools
bun install
```

## Build & Deployment

### Build Tools
- **Go**: Built-in compiler
- **Python**: Poetry build
- **Frontend**: Vite build
- **Docker**: Multi-stage builds

### Local Deployment
- **Development**: Docker Compose
- **Testing**: Local services
- **Staging**: Local Kubernetes (minikube)

## Performance Targets

### Response Times
- **API**: < 100ms (95th percentile)
- **Web**: < 2s page load
- **Scans**: < 30s for small repos
- **Enrichment**: < 5s per vulnerability

### Throughput
- **API**: 1000+ requests/second
- **Agents**: 100+ concurrent scans
- **Database**: 10,000+ operations/second
- **Cache**: 50,000+ operations/second

### Scalability
- **Horizontal**: All components scalable
- **Vertical**: Resource limits defined
- **Database**: Read replicas, sharding ready
- **Cache**: Redis cluster ready
