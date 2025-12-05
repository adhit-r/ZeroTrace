# 🎯 ZeroTrace Services Status

## Current Status

Run this to check all services:

```bash
./run-local.sh
```

## Service URLs

✅ **Backend API**: http://localhost:8080
- Health check: `curl http://localhost:8080/health`
- Status: Running ✓

✅ **Enrichment Service**: http://localhost:8000  
- Health check: `curl http://localhost:8000/health`
- Status: Running ✓

✅ **Frontend**: http://localhost:5173
- Open in browser: http://localhost:5173
- Status: Starting...

## Quick Commands

### Start Everything
```bash
./run-local.sh
```

### Check Status
```bash
# Backend
curl http://localhost:8080/health

# Enrichment
curl http://localhost:8000/health

# Frontend (open in browser)
open http://localhost:5173
```

### View Logs
```bash
tail -f api.log          # Backend
tail -f enrichment.log   # Enrichment
tail -f frontend.log    # Frontend
```

### Stop Everything
```bash
# Press Ctrl+C in the terminal running ./run-local.sh
# OR manually:
pkill -f zerotrace-api
pkill -f uvicorn
pkill -f "bun run dev"
docker-compose down
```

## Agent (macOS)

To run the agent:

```bash
./run-agent.sh
# OR
open agent-go/mdm/build/ZeroTrace\ Agent.app
```


