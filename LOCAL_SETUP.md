# ZeroTrace Local Development Setup

This guide will help you run all ZeroTrace services locally.

## Prerequisites

1. **Go** (1.23+ for API, 1.24+ for Agent)
2. **Python** (3.11+) with `uv` package manager
3. **Bun** (latest) for frontend
4. **PostgreSQL** (15+) - can run via podman-compose
5. **Valkey** (Redis-compatible) - can run via podman-compose
6. **Podman** (for running PostgreSQL and Valkey)

## Quick Start

### Option 1: Run All Services at Once

```bash
./run-local.sh
```

This will:
- Check and start PostgreSQL & Valkey (via podman-compose)
- Start the Go API backend (port 8080)
- Start the Python enrichment service (port 8000)
- Start the React frontend (port 5173)

### Option 2: Run Services Individually

#### 1. Start Dependencies (PostgreSQL & Valkey)

```bash
podman-compose up -d postgres valkey
```

Verify they're running:
```bash
pg_isready -h localhost -p 5432
valkey-cli -h localhost -p 6379 ping
```

#### 2. Backend (Go API)

```bash
cd api-go

# Create .env file if it doesn't exist
cp env.example .env

# Build (optional, can use go run)
go build -o zerotrace-api ./cmd/api

# Run
./zerotrace-api
# OR
go run ./cmd/api
```

The API will be available at: http://localhost:8080

#### 3. Enrichment Service (Python)

```bash
cd enrichment-python

# Create .env file if it doesn't exist
cp env.example .env

# Create virtual environment
uv venv

# Activate and install dependencies
source venv/bin/activate  # On macOS/Linux
# OR
venv\Scripts\activate  # On Windows

uv pip install -r requirements.txt

# Run
uv run uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

The enrichment service will be available at: http://localhost:8000

#### 4. Frontend (React)

```bash
cd web-react

# Install dependencies (first time only)
bun install

# Run development server
bun run dev
```

The frontend will be available at: http://localhost:5173

#### 5. Agent (macOS)

```bash
# Option A: Run the pre-built agent
cd agent-go
./agent

# Option B: Build and run
cd agent-go
go build -o agent ./cmd/agent
./agent

# Option C: Use the script
./run-agent.sh
```

**For macOS .app bundle:**

The agent can be packaged as a macOS .app bundle. Check the `agent-go/mdm/` directory for the .app bundle, or build it:

```bash
cd agent-go
# The .app bundle should be in mdm/zerotrace-agent.app/
open mdm/zerotrace-agent.app
```

## Environment Configuration

### Backend (.env in api-go/)

```env
API_PORT=8080
API_HOST=0.0.0.0
API_MODE=debug
DB_HOST=localhost
DB_PORT=5432
DB_NAME=zerotrace
DB_USER=postgres
DB_PASSWORD=password
DB_SSL_MODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
JWT_SECRET=dev-secret-key-change-in-production
```

### Enrichment (.env in enrichment-python/)

```env
ENRICHMENT_PORT=8000
DATABASE_URL=postgresql://postgres:password@localhost:5432/zerotrace
REDIS_URL=redis://localhost:6379/0
CPE_GUESSER_ENABLED=true
```

### Agent (.env in agent-go/)

```env
AGENT_ID=agent-001
AGENT_NAME=ZeroTrace Agent
COMPANY_ID=company-001
API_KEY=your-api-key-here
API_ENDPOINT=http://localhost:8080
```

## Service URLs

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080
- **Enrichment Service**: http://localhost:8000
- **PostgreSQL**: localhost:5432
- **Valkey**: localhost:6379

## Troubleshooting

### PostgreSQL Connection Issues

```bash
# Check if PostgreSQL is running
pg_isready -h localhost -p 5432

# If not running, start with podman-compose
podman-compose up -d postgres

# Check logs
podman-compose logs postgres
```

### Valkey Connection Issues

```bash
# Check if Valkey is running
valkey-cli -h localhost -p 6379 ping

# If not running, start with podman-compose
podman-compose up -d valkey

# Check logs
podman-compose logs valkey
```

### Port Already in Use

If a port is already in use, either:
1. Stop the service using that port
2. Change the port in the respective `.env` file

### Database Migrations

The Go API will automatically run migrations on startup. If you need to run them manually:

```bash
cd api-go
go run ./migrations
```

## Stopping Services

Press `Ctrl+C` in each terminal, or:

```bash
# Stop all podman containers
podman-compose down

# Find and kill processes
pkill -f zerotrace-api
pkill -f uvicorn
pkill -f "bun run dev"
```

## Development Tips

1. **Hot Reload**: 
   - Frontend: Automatic with Vite
   - Enrichment: Use `--reload` flag with uvicorn
   - Backend: Rebuild or use `air` for hot reload

2. **Logs**:
   - Backend: Check console output or `api.log`
   - Enrichment: Check console output
   - Frontend: Check browser console and terminal

3. **Database Access**:
   ```bash
   psql -h localhost -U postgres -d zerotrace
   ```

4. **Valkey Access**:
   ```bash
   valkey-cli -h localhost -p 6379
   ```

## Next Steps

1. Register an agent via the API
2. Configure organization profile
3. Start scanning for vulnerabilities
4. View results in the frontend dashboard


