# ZeroTrace Quick Start Guide

## 🚀 Run Everything Locally

### Prerequisites Check

```bash
# Check if you have all required tools
go version      # Need 1.23+
python3 --version  # Need 3.11+
bun --version   # Need latest
podman --version # For PostgreSQL/Valkey
```

### 1️⃣ Start Dependencies (PostgreSQL & Valkey)

```bash
# Start PostgreSQL and Valkey using podman-compose
podman-compose up -d postgres valkey

# Verify they're running
pg_isready -h localhost -p 5432
valkey-cli -h localhost -p 6379 ping
```

### 2️⃣ Backend (Go API) - Port 8080

```bash
cd api-go

# Setup environment
cp env.example .env

# Build and run
go build -o zerotrace-api ./cmd/api
./zerotrace-api

# OR run directly
go run ./cmd/api
```

**Backend URL**: http://localhost:8080

### 3️⃣ Enrichment Service (Python) - Port 8000

```bash
cd enrichment-python

# Setup environment
cp env.example .env

# Create virtual environment and install dependencies
uv venv
source venv/bin/activate  # On macOS/Linux
uv pip install -r requirements.txt

# Run the service
uv run uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

**Enrichment URL**: http://localhost:8000

### 4️⃣ Frontend (React) - Port 5173

```bash
cd web-react

# Install dependencies (first time only)
bun install

# Run development server
bun run dev
```

**Frontend URL**: http://localhost:5173

### 5️⃣ Agent (macOS)

#### Option A: Use the .app Bundle (Recommended for macOS)

```bash
cd agent-go

# Open the macOS .app bundle
open mdm/build/ZeroTrace\ Agent.app

# OR use the script
./run-agent.sh
```

#### Option B: Run Binary Directly

```bash
cd agent-go

# Setup environment
cp env.example .env
# Edit .env with your API endpoint: API_ENDPOINT=http://localhost:8080

# Build if needed
go build -o agent ./cmd/agent

# Run
./agent
```

#### Option C: Run from Source

```bash
cd agent-go
cp env.example .env
go run ./cmd/agent
```

## 📋 Service URLs Summary

| Service | URL | Port |
|---------|-----|------|
| Frontend | http://localhost:5173 | 5173 |
| Backend API | http://localhost:8080 | 8080 |
| Enrichment | http://localhost:8000 | 8000 |
| PostgreSQL | localhost:5432 | 5432 |
| Valkey | localhost:6379 | 6379 |

## 🔧 Environment Files

Each service needs a `.env` file. Copy from `env.example`:

```bash
# Backend
cd api-go && cp env.example .env

# Enrichment
cd enrichment-python && cp env.example .env

# Agent
cd agent-go && cp env.example .env
```

## 🛑 Stopping Services

Press `Ctrl+C` in each terminal window, or:

```bash
# Stop containers
podman-compose down

# Kill processes
pkill -f zerotrace-api
pkill -f uvicorn
pkill -f "bun run dev"
```

## 📝 Quick Test

1. Start all services (see above)
2. Open http://localhost:5173 in your browser
3. The frontend should connect to the backend
4. Register an agent via the API or frontend
5. Start scanning!

## 🐛 Troubleshooting

**Port already in use?**
- Change the port in the respective `.env` file
- Or stop the service using that port

**Database connection error?**
- Make sure PostgreSQL is running: `podman-compose up -d postgres`
- Check credentials in `.env` files

**Valkey connection error?**
- Make sure Valkey is running: `podman-compose up -d valkey`
- Check `REDIS_HOST` and `REDIS_PORT` in `.env` files

## 📚 More Information

See `LOCAL_SETUP.md` for detailed setup instructions.


