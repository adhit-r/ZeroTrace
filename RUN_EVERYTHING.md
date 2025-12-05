# 🚀 Run Everything - ZeroTrace

## Quick Command

```bash
./run-local.sh
```

This single command will:
1. ✅ Start PostgreSQL & Valkey (via docker-compose/podman-compose)
2. ✅ Start Backend API (Go) on port 8080
3. ✅ Start Enrichment Service (Python) on port 8000
4. ✅ Start Frontend (React) on port 5173

## What Gets Started

### 1. Dependencies (PostgreSQL & Valkey)
- Automatically detected and started via docker-compose or podman-compose
- PostgreSQL: `localhost:5432`
- Valkey: `localhost:6379`

### 2. Backend API (Go)
- **URL**: http://localhost:8080
- **Logs**: `api.log`
- **Location**: `api-go/`

### 3. Enrichment Service (Python)
- **URL**: http://localhost:8000
- **Logs**: `enrichment.log`
- **Location**: `enrichment-python/`

### 4. Frontend (React with Bun)
- **URL**: http://localhost:5173
- **Logs**: `frontend.log`
- **Location**: `web-react/`

## Service URLs

| Service | URL | Status Check |
|---------|-----|--------------|
| Frontend | http://localhost:5173 | Open in browser |
| Backend API | http://localhost:8080 | `curl http://localhost:8080/health` |
| Enrichment | http://localhost:8000 | `curl http://localhost:8000/health` |
| PostgreSQL | localhost:5432 | `pg_isready -h localhost` |
| Valkey | localhost:6379 | `valkey-cli ping` or `redis-cli ping` |

## View Logs

```bash
# Backend logs
tail -f api.log

# Enrichment logs
tail -f enrichment.log

# Frontend logs
tail -f frontend.log

# All logs at once
tail -f api.log enrichment.log frontend.log
```

## Stop Everything

Press `Ctrl+C` in the terminal where you ran `./run-local.sh`, or:

```bash
# Stop containers
docker-compose down
# OR
podman-compose down

# Kill processes
pkill -f zerotrace-api
pkill -f uvicorn
pkill -f "bun run dev"
```

## Manual Start (Alternative)

If you prefer to start services individually:

### 1. Start Dependencies
```bash
docker-compose up -d postgres valkey
# OR
podman-compose up -d postgres valkey
```

### 2. Start Backend
```bash
cd api-go
cp env.example .env
go build -o zerotrace-api ./cmd/api
./zerotrace-api
```

### 3. Start Enrichment
```bash
cd enrichment-python
cp env.example .env
uv venv
source venv/bin/activate
uv pip install -r requirements.txt
uv run uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

### 4. Start Frontend
```bash
cd web-react
bun install
bun run dev
```

## Troubleshooting

### Port Already in Use
```bash
# Find what's using the port
lsof -i :8080  # Backend
lsof -i :8000  # Enrichment
lsof -i :5173  # Frontend

# Kill the process or change port in .env
```

### Services Not Starting
1. Check logs: `tail -f *.log`
2. Verify dependencies: `docker-compose ps`
3. Check environment files: Ensure `.env` exists in each service directory

### Database Connection Issues
```bash
# Restart PostgreSQL
docker-compose restart postgres

# Check connection
psql -h localhost -U postgres -d zerotrace
```

## Next Steps

1. Open http://localhost:5173 in your browser
2. The frontend will connect to the backend automatically
3. Register an agent or start scanning
4. View results in the dashboard

## Agent (Separate)

To run the agent separately:

```bash
# macOS .app bundle
open agent-go/mdm/build/ZeroTrace\ Agent.app

# OR use the script
./run-agent.sh

# OR run binary
cd agent-go
cp env.example .env
# Edit .env: API_ENDPOINT=http://localhost:8080
./agent
```


