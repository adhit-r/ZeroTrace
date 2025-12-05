#!/bin/bash

# Test script for network scanning feature
# This script helps test the network scanning functionality locally

set -e

echo "=== ZeroTrace Network Scanner Test ==="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "Error: Please run this script from the agent-go directory"
    exit 1
fi

# Check dependencies
echo "Checking dependencies..."
if ! command -v nmap &> /dev/null; then
    echo "❌ Nmap not found. Installing..."
    if command -v brew &> /dev/null; then
        brew install nmap
    else
        echo "Please install nmap manually: https://nmap.org/download.html"
        exit 1
    fi
else
    echo "✅ Nmap found: $(nmap --version | head -n 1)"
fi

if ! command -v nuclei &> /dev/null; then
    echo "❌ Nuclei not found. Installing..."
    if command -v brew &> /dev/null; then
        brew install nuclei
    else
        echo "Please install nuclei manually: https://github.com/projectdiscovery/nuclei"
        exit 1
    fi
else
    echo "✅ Nuclei found: $(nuclei -version 2>&1 | head -n 1)"
fi

# Check Go dependencies
echo ""
echo "Installing Go dependencies..."
go mod tidy

# Check if .env exists
if [ ! -f ".env" ]; then
    echo ""
    echo "Creating .env file from env.example..."
    cp env.example .env
    echo "⚠️  Please edit .env and set your API_ENDPOINT and other configuration"
fi

# Build the agent
echo ""
echo "Building agent..."
go build -o zerotrace-agent cmd/agent/main.go

if [ $? -eq 0 ]; then
    echo "✅ Agent built successfully"
else
    echo "❌ Build failed"
    exit 1
fi

echo ""
echo "=== Setup Complete ==="
echo ""
echo "To test network scanning:"
echo "1. Make sure your API server is running (if testing with API)"
echo "2. Edit .env file and set NETWORK_SCAN_ENABLED=true"
echo "3. Run: ./zerotrace-agent"
echo ""
echo "Or test network scanning directly:"
echo "  go run cmd/agent/main.go"
echo ""
echo "The agent will:"
echo "- Scan your local network automatically"
echo "- Identify devices (switches, routers, IoT, phones, servers)"
echo "- Detect configuration errors"
echo "- Run vulnerability scans with Nuclei"
echo "- Send results to API (if configured)"
echo ""

