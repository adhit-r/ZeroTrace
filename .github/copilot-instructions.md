# ZeroTrace Copilot Instructions

## Big Picture Architecture
- **ZeroTrace** is an enterprise vulnerability detection platform with four main components:
  - `api-go/`: Go-based REST API backend (multi-tenant, high-performance, JWT auth, real-time stats)
  - `agent-go/`: Universal Go agent for software discovery, vulnerability scanning, and MDM deployment
  - `enrichment-python/`: Python FastAPI service for CVE enrichment and batch processing
  - `web-react/`: React + Vite frontend for dashboards and management
- Data flows: Agents discover software, send results to API, which triggers enrichment and updates dashboards. Monitoring via Prometheus/Grafana.

## Developer Workflows
- **Build/Run All Services Locally:**
  - `docker-compose up -d` (root)
- **Individual Components:**
  - API: `cd api-go && go run cmd/api/main.go`
  - Agent: `cd agent-go && go build -o zerotrace-agent cmd/agent/main.go && ./zerotrace-agent`
  - Enrichment: `cd enrichment-python && pip install -r requirements.txt && uvicorn app.main:app --reload`
  - Frontend: `cd web-react && npm install && npm run dev`
- **Testing:**
  - Go: `go test ./...` (in `agent-go/` or `api-go/`)
  - Python: Use `pytest` in `enrichment-python/app/`
- **MDM Deployment:**
  - Build macOS package: `agent-go/mdm/build-macos-pkg.sh`

## Project-Specific Conventions
- **Universal Agent:** Single binary for all orgs, silent background operation for MDM, tray UI for dev
- **API Endpoints:** All agent communication via `/api/v1/agent/*` and `/api/enrollment/*`
- **Environment Variables:** See `env.example` in each component for required config
- **Multi-level Caching:** Memory, Redis, Memcached (API, enrichment)
- **Monitoring:** Prometheus metrics exposed by API and agent; dashboards in Grafana
- **Batch/Parallel Processing:** Python enrichment uses batch and ultra-optimized parallel code

## Integration Points & Dependencies
- **API <-> Agent:** Enrollment, heartbeat, results, registration (see `agent-go/internal/communicator/`)
- **API <-> Enrichment:** CVE lookups via HTTP (see `api-go/internal/handlers/`)
- **API <-> Frontend:** REST endpoints for dashboard data
- **Monitoring:** Prometheus scrapes API/agent; Grafana dashboards
- **MDM:** Agent packaged for Intune, Jamf, Azure AD, Workspace ONE (see `agent-go/mdm/README.md`)

## Key Files & Directories
- `api-go/cmd/api/main.go`: API entry point & route registration
- `agent-go/cmd/agent/main.go`: Full agent (tray UI)
- `agent-go/cmd/agent-simple/main.go`: Silent MDM agent
- `enrichment-python/app/ultra_optimized_enrichment.py`: High-performance enrichment logic
- `web-react/src/`: React app source
- `agent-go/mdm/README.md`: MDM deployment guide
- `PERFORMANCE_OPTIMIZATION_SUMMARY.md`: Performance strategies
- `SCALABLE_ARCHITECTURE_SUMMARY.md`: Architecture rationale

## Patterns & Examples
- **Adding API Endpoints:**
  1. Handler in `api-go/internal/handlers/`
  2. Logic in `api-go/internal/services/`
  3. Register route in `api-go/cmd/api/main.go`
- **Agent Communication:**
  - Use `agent-go/internal/communicator/` for API calls
- **Enrichment:**
  - Batch via `batch_enrichment.py`, ultra-optimized via `ultra_optimized_enrichment.py`
- **Testing:**
  - Go: `go test ./internal/<component>`
  - Python: `pytest app/<module>.py`

## Additional Resources
- `docs/` and `wiki/` for architecture, deployment, and troubleshooting
- `ROADMAP.md` for current and planned features

---
For unclear or missing conventions, consult the relevant README or ask for clarification in the project wiki or discussions.
