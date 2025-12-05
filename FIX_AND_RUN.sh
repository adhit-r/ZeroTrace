#!/bin/bash

# Quick Fix and Run Script for ZeroTrace
# This script fixes common issues and starts all services

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔧 Fixing and Starting ZeroTrace Services${NC}"
echo ""

# Add bun to PATH
export PATH="$HOME/.bun/bin:$PATH"

# Step 1: Kill existing processes
echo -e "${YELLOW}Step 1: Cleaning up existing processes...${NC}"
pkill -f zerotrace-api 2>/dev/null || true
pkill -f "uvicorn app.main:app" 2>/dev/null || true
pkill -f "bun run dev" 2>/dev/null || true
pkill -f vite 2>/dev/null || true
lsof -ti :8080 | xargs kill -9 2>/dev/null || true
lsof -ti :8000 | xargs kill -9 2>/dev/null || true
lsof -ti :5173 | xargs kill -9 2>/dev/null || true
sleep 2
echo -e "${GREEN}✓ Cleaned up${NC}"
echo ""

# Step 2: Start dependencies
echo -e "${YELLOW}Step 2: Starting dependencies...${NC}"
cd /Users/adhi/axonome/ZeroTrace
docker-compose up -d postgres valkey 2>/dev/null || echo "⚠️  Dependencies may already be running"
sleep 2
echo -e "${GREEN}✓ Dependencies ready${NC}"
echo ""

# Step 3: Start Backend
echo -e "${YELLOW}Step 3: Starting Backend (Go API)...${NC}"
cd api-go
[ -f .env ] || cp env.example .env
go build -o zerotrace-api ./cmd/api
./zerotrace-api > ../api.log 2>&1 &
BACKEND_PID=$!
cd ..
sleep 3
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend running on http://localhost:8080${NC}"
else
    echo -e "${RED}✗ Backend failed. Check api.log${NC}"
fi
echo ""

# Step 4: Start Enrichment
echo -e "${YELLOW}Step 4: Starting Enrichment Service...${NC}"
cd enrichment-python
[ -f .env ] || cp env.example .env
# Use uv run directly (works with both venv and .venv)
uv run uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload > ../enrichment.log 2>&1 &
ENRICHMENT_PID=$!
cd ..
sleep 3
if curl -s http://localhost:8000/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Enrichment running on http://localhost:8000${NC}"
else
    echo -e "${RED}✗ Enrichment failed. Check enrichment.log${NC}"
fi
echo ""

# Step 5: Start Frontend
echo -e "${YELLOW}Step 5: Starting Frontend...${NC}"
cd web-react
# Install dependencies if needed
if [ ! -d "node_modules" ] || [ ! -f "node_modules/.bin/vite" ]; then
    echo "Installing dependencies..."
    bun install
fi
bun run dev > ../frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..
sleep 5
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Frontend running on http://localhost:5173${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend starting... (check frontend.log if issues persist)${NC}"
fi
echo ""

# Final Status
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ Services Started${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${GREEN}📍 URLs:${NC}"
echo -e "   Frontend:    ${BLUE}http://localhost:5173${NC}"
echo -e "   Backend:     ${BLUE}http://localhost:8080${NC}"
echo -e "   Enrichment:  ${BLUE}http://localhost:8000${NC}"
echo ""
echo -e "${YELLOW}📝 Logs:${NC}"
echo -e "   tail -f api.log"
echo -e "   tail -f enrichment.log"
echo -e "   tail -f frontend.log"
echo ""
echo -e "${YELLOW}🛑 Stop: pkill -f 'zerotrace-api|uvicorn|bun run dev'${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"


