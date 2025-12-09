#!/bin/bash
set -e

echo "=== ZeroTrace Enrichment Service Startup ==="
echo ""

# Step 1: Check if CVE data exists
if [ ! -f "cve_data.json" ]; then
    echo "️  cve_data.json not found. Generating it first..."
    echo "   This may take a few minutes..."
    python3 scripts/update_cve_data.py || {
        echo " Failed to generate cve_data.json"
        echo "   You can skip this and run migration later, but it's recommended."
        read -p "Continue anyway? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    }
else
    echo " cve_data.json found"
fi

# Step 2: Start Docker/Podman services
echo ""
echo " Starting services with docker-compose..."
echo "   (If you prefer podman, start Docker Desktop or use: podman compose up -d)"

# Try docker-compose first, fallback to podman compose
if command -v docker-compose &> /dev/null && docker info &> /dev/null; then
    docker-compose up -d
elif command -v podman &> /dev/null && podman info &> /dev/null; then
    echo "Using podman compose..."
    podman compose up -d
else
    echo " Neither Docker nor Podman is running."
    echo "   Please start Docker Desktop or Podman service first."
    exit 1
fi

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Step 3: Check service health
echo ""
echo " Checking service status..."
if command -v docker-compose &> /dev/null; then
    docker-compose ps
elif command -v podman &> /dev/null; then
    podman compose ps
fi

echo ""
echo " Services started!"
echo ""
echo "Next steps:"
echo "1. Run migration: python3 scripts/migrate_to_postgres.py"
echo "2. Start application: uvicorn app.main:app --reload"
echo ""

