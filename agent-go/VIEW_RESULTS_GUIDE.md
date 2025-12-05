# View Results Guide

## Ports and URLs

### API Server (Backend)
- **Port:** `8080`
- **URL:** `http://localhost:8080`
- **Health Check:** `http://localhost:8080/health`
- **API Docs:** `http://localhost:8080/api/docs` (if available)

### Web UI (Frontend)
- **Port:** `3000`
- **URL:** `http://localhost:3000`
- **This is where you see all results!**

### Enrichment Service
- **Port:** `8000`
- **URL:** `http://localhost:8000`

## How to View Results

### Step 1: Start the API Server

```bash
cd api-go
go run cmd/api/main.go
```

Or with Docker:
```bash
podman-compose up api
```

The API will start on **http://localhost:8080**

### Step 2: Start the Web UI

```bash
cd web-react
bun run dev
```

The web UI will start on **http://localhost:3000**

### Step 3: View Results

1. **Open browser:** Go to `http://localhost:3000`
2. **Login/Register** (if authentication is enabled)
3. **View Dashboard:**
   - Agents list
   - Vulnerabilities
   - Network scan results
   - System information

## Menu Bar Icon Status

The menu bar icon shows:
- **🟢 Green** = Connected to API (port 8080)
- **⚫ Gray** = Disconnected (API not running)

**Status updates every 10 seconds**

## API Endpoints to Check Results

### List All Agents
```bash
curl http://localhost:8080/api/agents
```

### Get Agent Results
```bash
curl http://localhost:8080/api/agents/{agent_id}/results
```

### Get Vulnerabilities
```bash
curl http://localhost:8080/api/vulnerabilities
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Troubleshooting

### Status Not Updating in Menu Bar

1. **Check if API is running:**
   ```bash
   curl http://localhost:8080/health
   ```

2. **Check agent logs:**
   ```bash
   tail -f ~/.zerotrace/logs/agent.log
   ```

3. **Restart the agent** if needed

### CPU Not Showing

- CPU metrics update every 10 seconds
- First update happens after 2 seconds
- If still not showing, check monitor is started

### Can't See Results

1. **Make sure API is running** on port 8080
2. **Make sure Web UI is running** on port 3000
3. **Check agent is connected** (green icon in menu bar)
4. **Wait for scan to complete** (first scan takes time)

## Quick Start

```bash
# Terminal 1: Start API
cd api-go && go run cmd/api/main.go

# Terminal 2: Start Web UI
cd web-react && bun run dev

# Terminal 3: Launch Agent
open "agent-go/mdm/build/ZeroTrace Agent.app"

# Browser: View Results
open http://localhost:3000
```

## Summary

| Service | Port | URL | Purpose |
|---------|------|-----|---------|
| **API** | 8080 | http://localhost:8080 | Backend API |
| **Web UI** | 3000 | http://localhost:3000 | **View results here!** |
| **Enrichment** | 8000 | http://localhost:8000 | Vulnerability enrichment |

**Main URL to view results: http://localhost:3000** 🎯


