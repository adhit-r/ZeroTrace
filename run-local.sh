#!/bin/bash

# ZeroTrace Local Development Runner
# This script starts all services locally

set -e

# Add bun to PATH if it exists
if [ -d "$HOME/.bun/bin" ]; then
    export PATH="$HOME/.bun/bin:$PATH"
fi

echo "🚀 Starting ZeroTrace Local Development Environment"
echo "=================================================="

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Detect container runtime
detect_runtime() {
    if command -v podman-compose &> /dev/null; then
        echo "podman-compose"
    elif command -v docker-compose &> /dev/null; then
        echo "docker-compose"
    elif command -v docker &> /dev/null && docker compose version &> /dev/null; then
        echo "docker compose"
    else
        echo ""
    fi
}

RUNTIME=$(detect_runtime)

# Check if dependencies are running
check_dependencies() {
    echo -e "${YELLOW}Checking dependencies...${NC}"
    
    # Check PostgreSQL
    if ! pg_isready -h localhost -p 5432 -U postgres > /dev/null 2>&1; then
        if [ -z "$RUNTIME" ]; then
            echo -e "${RED}❌ PostgreSQL not running and no container runtime found${NC}"
            echo -e "${YELLOW}Please install podman-compose or docker-compose, or start PostgreSQL manually${NC}"
            exit 1
        fi
        echo -e "${YELLOW}⚠️  PostgreSQL not running. Starting with $RUNTIME...${NC}"
        $RUNTIME up -d postgres valkey
        echo "Waiting for services to be ready..."
        sleep 5
    else
        echo -e "${GREEN}✓ PostgreSQL is running${NC}"
    fi
    
    # Check Valkey
    if ! command -v valkey-cli &> /dev/null; then
        # Try redis-cli as fallback
        if command -v redis-cli &> /dev/null; then
            if ! redis-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
                echo -e "${YELLOW}⚠️  Valkey/Redis not running. Already started with $RUNTIME above.${NC}"
            else
                echo -e "${GREEN}✓ Valkey/Redis is running${NC}"
            fi
        else
            echo -e "${YELLOW}⚠️  Cannot check Valkey (valkey-cli not found), assuming it's running${NC}"
        fi
    else
        if ! valkey-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
            echo -e "${YELLOW}⚠️  Valkey not running. Already started with $RUNTIME above.${NC}"
        else
            echo -e "${GREEN}✓ Valkey is running${NC}"
        fi
    fi
}

# Start Backend (Go API)
start_backend() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Starting Backend (Go API)...${NC}"
    cd api-go
    
    # Check if .env exists
    if [ ! -f .env ]; then
        echo "Creating .env from env.example..."
        cp env.example .env
        echo -e "${YELLOW}⚠️  Please review api-go/.env and update if needed${NC}"
    fi
    
    # Build if needed
    if [ ! -f api ] && [ ! -f zerotrace-api ]; then
        echo "Building Go API..."
        go build -o zerotrace-api ./cmd/api
    fi
    
    # Run the API
    echo -e "${GREEN}✓ Starting API on http://localhost:8080${NC}"
    if [ -f zerotrace-api ]; then
        ./zerotrace-api > ../api.log 2>&1 &
    elif [ -f api ]; then
        ./api > ../api.log 2>&1 &
    else
        go run ./cmd/api > ../api.log 2>&1 &
    fi
    
    BACKEND_PID=$!
    cd ..
    echo -e "${GREEN}Backend started (PID: $BACKEND_PID)${NC}"
    sleep 2
}

# Start Enrichment Service
start_enrichment() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Starting Enrichment Service (Python)...${NC}"
    cd enrichment-python
    
    # Check if .env exists
    if [ ! -f .env ]; then
        echo "Creating .env from env.example..."
        cp env.example .env
        echo -e "${YELLOW}⚠️  Please review enrichment-python/.env and update if needed${NC}"
    fi
    
    # Check if virtual environment exists
    if [ ! -d ".venv" ]; then
        echo "Creating virtual environment..."
        uv venv
    fi
    
    # Install dependencies if needed
    if [ ! -f .venv/.installed ]; then
        echo "Installing dependencies..."
        uv pip install -r requirements.txt
        touch .venv/.installed
    fi
    
    # Run the service
    echo -e "${GREEN}✓ Starting Enrichment Service on http://localhost:8000${NC}"
    export PGVECTOR_ENABLED=false
    uv run uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload > ../enrichment.log 2>&1 &
    
    ENRICHMENT_PID=$!
    cd ..
    echo -e "${GREEN}Enrichment started (PID: $ENRICHMENT_PID)${NC}"
    sleep 2
}

# Start Frontend
start_frontend() {
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Starting Frontend (React)...${NC}"
    cd web-react
    
    # Check if bun is available
    if ! command -v bun &> /dev/null; then
        echo -e "${RED}❌ bun not found. Installing bun...${NC}"
        curl -fsSL https://bun.sh/install | bash
        export PATH="$HOME/.bun/bin:$PATH"
    fi
    
    # Install dependencies if needed
    if [ ! -d "node_modules" ]; then
        echo "Installing dependencies with bun..."
        export PATH="$HOME/.bun/bin:$PATH"
        bun install
    fi
    
    # Run the frontend
    echo -e "${GREEN}✓ Starting Frontend on http://localhost:5173${NC}"
    export PATH="$HOME/.bun/bin:$PATH"
    bun run dev > ../frontend.log 2>&1 &
    
    FRONTEND_PID=$!
    cd ..
    echo -e "${GREEN}Frontend started (PID: $FRONTEND_PID)${NC}"
    sleep 2
}

# Main execution
main() {
    check_dependencies
    
    # Start services
    start_backend
    sleep 2
    
    start_enrichment
    sleep 2
    
    start_frontend
    
    echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✅ All services started!${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${GREEN}📍 Service URLs:${NC}"
    echo -e "   Backend API:    ${BLUE}http://localhost:8080${NC}"
    echo -e "   Enrichment:     ${BLUE}http://localhost:8000${NC}"
    echo -e "   Frontend:       ${BLUE}http://localhost:5173${NC}"
    echo ""
    echo -e "${GREEN}📋 Process IDs:${NC}"
    echo -e "   Backend:    $BACKEND_PID"
    echo -e "   Enrichment: $ENRICHMENT_PID"
    echo -e "   Frontend:   $FRONTEND_PID"
    echo ""
    echo -e "${YELLOW}📝 Logs:${NC}"
    echo -e "   Backend:    tail -f api.log"
    echo -e "   Enrichment: tail -f enrichment.log"
    echo -e "   Frontend:   tail -f frontend.log"
    echo ""
    echo -e "${YELLOW}🛑 To stop all services, press Ctrl+C${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    # Wait for user interrupt
    wait
}

# Trap Ctrl+C and cleanup
cleanup() {
    echo -e "\n${YELLOW}Stopping services...${NC}"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$ENRICHMENT_PID" ]; then
        kill $ENRICHMENT_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    # Kill any remaining processes
    pkill -f "zerotrace-api" 2>/dev/null || true
    pkill -f "uvicorn app.main:app" 2>/dev/null || true
    pkill -f "bun run dev" 2>/dev/null || true
    echo -e "${GREEN}✓ All services stopped${NC}"
    exit 0
}

trap cleanup INT TERM

main

