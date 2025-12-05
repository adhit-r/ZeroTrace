#!/bin/bash

# ZeroTrace Service Starter
# Fixes common issues and starts all services

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🚀 Starting ZeroTrace Services${NC}"
echo ""

# Add bun to PATH if not already there
export PATH="$HOME/.bun/bin:$PATH"

# Check if bun is installed
if ! command -v bun &> /dev/null; then
    echo -e "${YELLOW}Installing bun...${NC}"
    curl -fsSL https://bun.sh/install | bash
    export BUN_INSTALL="$HOME/.bun"
    export PATH="$BUN_INSTALL/bin:$PATH"
fi

# Kill existing processes on ports
echo -e "${YELLOW}Cleaning up existing processes...${NC}"
lsof -ti :8080 | xargs kill -9 2>/dev/null || true
lsof -ti :8000 | xargs kill -9 2>/dev/null || true
lsof -ti :5173 | xargs kill -9 2>/dev/null || true
sleep 1

# Start dependencies
echo -e "${YELLOW}Starting dependencies (PostgreSQL & Valkey)...${NC}"
cd /Users/adhi/axonome/ZeroTrace
docker-compose up -d postgres valkey 2>/dev/null || podman-compose up -d postgres valkey 2>/dev/null || echo "⚠️  Please start PostgreSQL and Valkey manually"
sleep 3

# Start Backend
echo -e "${YELLOW}Starting Backend (port 8080)...${NC}"
cd api-go
[ -f .env ] || cp env.example .env
go build -o zerotrace-api ./cmd/api
./zerotrace-api > ../api.log 2>&1 &
BACKEND_PID=$!
cd ..
sleep 3

# Check backend
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend running on http://localhost:8080${NC}"
else
    echo -e "${RED}✗ Backend failed to start. Check api.log${NC}"
fi

# Start Enrichment (already running, just verify)
echo -e "${YELLOW}Checking Enrichment (port 8000)...${NC}"
if curl -s http://localhost:8000/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Enrichment already running on http://localhost:8000${NC}"
else
    echo -e "${YELLOW}Starting Enrichment...${NC}"
    cd enrichment-python
    [ -f .env ] || cp env.example .env
    [ -d venv ] || [ -d .venv ] || uv venv
    # Use uv run directly (no need to activate)
    uv run uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload > ../enrichment.log 2>&1 &
    ENRICHMENT_PID=$!
    cd ..
    sleep 3
fi

# Start Frontend
echo -e "${YELLOW}Starting Frontend (port 5173)...${NC}"
cd web-react
[ -d node_modules ] || bun install
bun run dev > ../frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..
sleep 5

# Check frontend
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Frontend running on http://localhost:5173${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend starting... (may take a few more seconds)${NC}"
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Services Status${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}📍 Service URLs:${NC}"
echo -e "   Frontend:    ${BLUE}http://localhost:5173${NC}"
echo -e "   Backend:     ${BLUE}http://localhost:8080${NC}"
echo -e "   Enrichment:  ${BLUE}http://localhost:8000${NC}"
echo ""
echo -e "${GREEN}📋 Process IDs:${NC}"
echo -e "   Backend:    $BACKEND_PID"
echo -e "   Enrichment: $ENRICHMENT_PID"
echo -e "   Frontend:   $FRONTEND_PID"
echo ""
echo -e "${YELLOW}📝 View Logs:${NC}"
echo -e "   tail -f api.log"
echo -e "   tail -f enrichment.log"
echo -e "   tail -f frontend.log"
echo ""
echo -e "${YELLOW}🛑 To stop: pkill -f zerotrace-api && pkill -f uvicorn && pkill -f 'bun run dev'${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

