#!/bin/bash

# ZeroTrace Agent Runner for macOS
# This script builds and runs the agent locally

set -e

echo " Starting ZeroTrace Agent"
echo "============================"

cd agent-go

# Check if .env exists
if [ ! -f .env ]; then
    echo "Creating .env from env.example..."
    cp env.example .env
    echo "️  Please edit .env with your API endpoint and credentials"
fi

# Check for macOS .app bundle
if [ -d "mdm/build/ZeroTrace Agent.app" ]; then
    echo "Found macOS .app bundle!"
    echo "Opening ZeroTrace Agent.app..."
    open "mdm/build/ZeroTrace Agent.app"
    exit 0
fi

# Check if agent binary exists
if [ ! -f agent ] && [ ! -f zerotrace-agent ]; then
    echo "Building agent..."
    go build -o agent ./cmd/agent
fi

# Run the agent
echo "Starting agent..."
if [ -f zerotrace-agent ]; then
    ./zerotrace-agent
elif [ -f agent ]; then
    ./agent
else
    go run ./cmd/agent
fi

