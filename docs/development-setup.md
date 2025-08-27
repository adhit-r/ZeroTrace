# ZeroTrace Development Setup

## Overview
This guide provides step-by-step instructions for setting up the ZeroTrace development environment locally. All services will run on your local machine using Podman and Bun as specified in the user preferences.

## Prerequisites

### Required Software
```bash
# macOS (using Homebrew)
brew install go@1.21
brew install python@3.11
brew install node@18
brew install bun
brew install podman
brew install postgresql@15
brew install redis

# Verify installations
go version          # Should be 1.21+
python3 --version   # Should be 3.11+
node --version      # Should be 18+
bun --version       # Should be latest
podman --version    # Should be latest
psql --version      # Should be 15+
redis-server --version  # Should be 7+
```

### Development Tools
```bash
# Install development tools
brew install --cask visual-studio-code
brew install --cask insomnia
brew install --cask dbeaver-community
brew install --cask redisinsight

# VS Code Extensions
code --install-extension golang.go
code --install-extension ms-python.python
code --install-extension bradlc.vscode-tailwindcss
code --install-extension esbenp.prettier-vscode
code --install-extension ms-vscode.vscode-typescript-next
```

## Project Structure Setup

### 1. Clone and Initialize Project
```bash
# Clone the repository
git clone <repository-url>
cd ZeroTrace

# Create project structure
mkdir -p agent-go api-go enrichment-python web-react docker docs
```

### 2. Initialize Go Modules
```bash
# Agent Go Module
cd agent-go
go mod init zerotrace/agent
go mod tidy

# API Go Module
cd ../api-go
go mod init zerotrace/api
go mod tidy
```

### 3. Initialize Python Environment
```bash
# Python Enrichment
cd ../enrichment-python
python3 -m venv venv
source venv/bin/activate
pip install poetry
poetry init
poetry install
```

### 4. Initialize React Frontend
```bash
# React Frontend
cd ../web-react
bun create vite . --template react-ts
bun install
```

## Database Setup

### 1. PostgreSQL Configuration
```bash
# Start PostgreSQL
brew services start postgresql@15

# Create database
createdb zerotrace

# Create user (optional)
createuser zerotrace_user
psql -d zerotrace -c "GRANT ALL PRIVILEGES ON DATABASE zerotrace TO zerotrace_user;"
```

### 2. Redis Configuration
```bash
# Start Redis
brew services start redis

# Test Redis connection
redis-cli ping  # Should return PONG
```

### 3. Database Migrations
```bash
# Run migrations for API
cd api-go
go run cmd/migrate/main.go

# Run migrations for enrichment
cd ../enrichment-python
poetry run python scripts/migrate.py
```

## Service Configuration

### 1. Environment Variables

#### Agent Configuration
```bash
# agent-go/.env
ZEROTRACE_AGENT_ID=dev-agent-001
ZEROTRACE_COMPANY_ID=dev-company-123
ZEROTRACE_API_URL=http://localhost:8080
ZEROTRACE_API_TOKEN=dev-agent-token
ZEROTRACE_SCAN_DEPTH=5
ZEROTRACE_MAX_FILE_SIZE=10485760
ZEROTRACE_PARALLEL_WORKERS=2
ZEROTRACE_LOG_LEVEL=debug
```

#### API Configuration
```bash
# api-go/.env
API_PORT=8080
API_HOST=0.0.0.0
API_MODE=debug
DB_HOST=localhost
DB_PORT=5432
DB_NAME=zerotrace
DB_USER=postgres
DB_PASSWORD=
DB_SSL_MODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
JWT_SECRET=dev-secret-key-change-in-production
JWT_EXPIRY=24h
```

#### Enrichment Configuration
```bash
# enrichment-python/.env
ENRICHMENT_HOST=0.0.0.0
ENRICHMENT_PORT=8000
DEBUG=true
DATABASE_URL=postgresql://postgres@localhost/zerotrace
REDIS_URL=redis://localhost:6379/0
NVD_API_KEY=
GITHUB_TOKEN=
MODEL_PATH=./models
```

#### Frontend Configuration
```bash
# web-react/.env
VITE_API_URL=http://localhost:8080
VITE_ENRICHMENT_URL=http://localhost:8000
VITE_WS_URL=ws://localhost:8080/ws
```

### 2. Podman Compose Setup
```yaml
# docker-compose.yml (works with podman-compose)
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: zerotrace-postgres
    environment:
      POSTGRES_DB: zerotrace
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - zerotrace-network

  redis:
    image: redis:7-alpine
    container_name: zerotrace-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - zerotrace-network

  api:
    build:
      context: ./api-go
      dockerfile: Dockerfile
    container_name: zerotrace-api
    environment:
      - API_PORT=8080
      - DB_HOST=postgres
      - REDIS_HOST=redis
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    networks:
      - zerotrace-network

  enrichment:
    build:
      context: ./enrichment-python
      dockerfile: Dockerfile
    container_name: zerotrace-enrichment
    environment:
      - ENRICHMENT_PORT=8000
      - DATABASE_URL=postgresql://postgres:password@postgres/zerotrace
      - REDIS_URL=redis://redis:6379/0
    ports:
      - "8000:8000"
    depends_on:
      - postgres
      - redis
    networks:
      - zerotrace-network

  web:
    build:
      context: ./web-react
      dockerfile: Dockerfile
    container_name: zerotrace-web
    ports:
      - "3000:3000"
    depends_on:
      - api
      - enrichment
    networks:
      - zerotrace-network

volumes:
  postgres_data:
  redis_data:

networks:
  zerotrace-network:
    driver: bridge
```

## Development Workflow

### 1. Start Development Environment
```bash
# Option 1: Using Podman Compose (recommended for full stack)
podman-compose up -d

# Option 2: Individual services
# Start databases
brew services start postgresql@15
brew services start redis

# Start API
cd api-go
go run cmd/api/main.go

# Start enrichment (in new terminal)
cd enrichment-python
poetry run uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# Start frontend (in new terminal)
cd web-react
bun run dev
```

### 2. Development Commands

#### Go Services
```bash
# API Development
cd api-go
go run cmd/api/main.go
go test ./...
go build -o api cmd/api/main.go

# Agent Development
cd agent-go
go run cmd/agent/main.go
go test ./...
go build -o agent cmd/agent/main.go
```

#### Python Service
```bash
# Enrichment Development
cd enrichment-python
poetry run uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
poetry run pytest
poetry run python scripts/train_models.py
```

#### React Frontend
```bash
# Frontend Development
cd web-react
bun run dev          # Development server
bun run build        # Production build
bun run test         # Run tests
bun run lint         # Lint code
bun run type-check   # TypeScript check
```

### 3. Testing Setup

#### API Testing
```bash
# Run API tests
cd api-go
go test -v ./...

# Run integration tests
go test -v -tags=integration ./...

# Run benchmarks
go test -bench=. ./...
```

#### Enrichment Testing
```bash
# Run Python tests
cd enrichment-python
poetry run pytest

# Run with coverage
poetry run pytest --cov=app --cov-report=html

# Run specific test file
poetry run pytest tests/test_enrichment.py -v
```

#### Frontend Testing
```bash
# Run frontend tests
cd web-react
bun run test

# Run E2E tests
bun run test:e2e

# Run with coverage
bun run test:coverage
```

### 4. Database Management

#### Migrations
```bash
# Create new migration
cd api-go
go run cmd/migrate/main.go create migration_name

# Run migrations
go run cmd/migrate/main.go up

# Rollback migration
go run cmd/migrate/main.go down 1
```

#### Database Seeding
```bash
# Seed development data
cd api-go
go run cmd/seed/main.go

# Seed test data
go run cmd/seed/main.go --env=test
```

### 5. Monitoring and Debugging

#### Logs
```bash
# View API logs
cd api-go
tail -f logs/api.log

# View enrichment logs
cd enrichment-python
tail -f logs/enrichment.log

# View frontend logs (in browser console)
```

#### Health Checks
```bash
# API health
curl http://localhost:8080/health

# Enrichment health
curl http://localhost:8000/health

# Database health
psql -d zerotrace -c "SELECT version();"

# Redis health
redis-cli ping
```

#### Performance Monitoring
```bash
# Monitor API performance
curl http://localhost:8080/metrics

# Monitor enrichment performance
curl http://localhost:8000/metrics

# Database performance
psql -d zerotrace -c "SELECT * FROM pg_stat_activity;"
```

## Troubleshooting

### Common Issues

#### 1. Port Conflicts
```bash
# Check what's using a port
lsof -i :8080
lsof -i :8000
lsof -i :3000

# Kill process using port
kill -9 <PID>
```

#### 2. Database Connection Issues
```bash
# Check PostgreSQL status
brew services list | grep postgresql

# Restart PostgreSQL
brew services restart postgresql@15

# Check connection
psql -d zerotrace -c "SELECT 1;"
```

#### 3. Redis Connection Issues
```bash
# Check Redis status
brew services list | grep redis

# Restart Redis
brew services restart redis

# Check connection
redis-cli ping
```

#### 4. Go Module Issues
```bash
# Clean go mod cache
go clean -modcache

# Update go mod
go mod tidy
go mod download
```

#### 5. Python Environment Issues
```bash
# Recreate virtual environment
cd enrichment-python
rm -rf venv
python3 -m venv venv
source venv/bin/activate
poetry install
```

#### 6. Frontend Build Issues
```bash
# Clear node modules
cd web-react
rm -rf node_modules
bun install

# Clear build cache
rm -rf dist
bun run build
```

### Performance Optimization

#### 1. Development Performance
```bash
# Use Go modules proxy
go env -w GOPROXY=https://proxy.golang.org,direct

# Use Bun for faster package installation
bun install --frozen-lockfile

# Use Poetry for Python dependencies
poetry install --no-dev
```

#### 2. Database Performance
```bash
# Optimize PostgreSQL for development
psql -d zerotrace -c "ALTER SYSTEM SET shared_buffers = '256MB';"
psql -d zerotrace -c "ALTER SYSTEM SET effective_cache_size = '1GB';"
psql -d zerotrace -c "SELECT pg_reload_conf();"
```

#### 3. Memory Usage
```bash
# Monitor memory usage
top -o mem

# Check specific process memory
ps aux | grep zerotrace
```

## Production Preparation

### 1. Environment Variables
```bash
# Create production environment files
cp .env.example .env.production

# Update with production values
# - Strong JWT secrets
# - Production database credentials
# - API keys for external services
# - Proper logging levels
```

### 2. Security Configuration
```bash
# Generate strong secrets
openssl rand -hex 32  # JWT secret
openssl rand -hex 32  # API token

# Update firewall rules
# Configure SSL certificates
# Set up monitoring and alerting
```

### 3. Performance Tuning
```bash
# Database optimization
# Connection pooling
# Caching strategies
# Load balancing
# CDN configuration
```

## Next Steps

1. **Read the Architecture Documentation**: Understand the system design
2. **Review API Endpoints**: Familiarize with the API structure
3. **Explore Frontend Components**: Understand the UI architecture
4. **Set up Monitoring**: Configure logging and metrics
5. **Write Tests**: Add comprehensive test coverage
6. **Performance Testing**: Test with large datasets
7. **Security Review**: Conduct security assessment
8. **Deployment Planning**: Plan production deployment
