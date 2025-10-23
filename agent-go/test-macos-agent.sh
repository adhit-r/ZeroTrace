#!/bin/bash

# ZeroTrace Agent macOS Test Script
# This script helps test the agent on macOS

echo "🔍 ZeroTrace Agent macOS Test"
echo "=============================="

# Check if agent binary exists
if [ ! -f "./agent" ]; then
    echo "❌ Agent binary not found. Building agent..."
    go build -o agent ./cmd/agent/
    if [ $? -ne 0 ]; then
        echo "❌ Failed to build agent"
        exit 1
    fi
    echo "✅ Agent built successfully"
fi

# Check if agent is already running
if pgrep -f "zerotrace" > /dev/null; then
    echo "⚠️  Agent is already running. Stopping existing agent..."
    pkill -f "zerotrace"
    sleep 2
fi

echo ""
echo "🚀 Starting ZeroTrace Agent..."
echo "=============================="

# Set environment variables for testing
export ZEROTRACE_API_ENDPOINT="http://localhost:8080"
export ZEROTRACE_ORGANIZATION_ID="test-org-123"
export ZEROTRACE_AGENT_ID="test-agent-$(date +%s)"

echo "📋 Configuration:"
echo "   API Endpoint: $ZEROTRACE_API_ENDPOINT"
echo "   Organization ID: $ZEROTRACE_ORGANIZATION_ID"
echo "   Agent ID: $ZEROTRACE_AGENT_ID"
echo ""

# Start agent in background
echo "🔄 Starting agent in background..."
./agent &
AGENT_PID=$!

echo "✅ Agent started with PID: $AGENT_PID"
echo ""

# Wait a moment for agent to initialize
sleep 3

# Check if agent is running
if ps -p $AGENT_PID > /dev/null; then
    echo "✅ Agent is running successfully!"
    echo ""
    echo "📊 Agent Status:"
    echo "   PID: $AGENT_PID"
    echo "   Process: $(ps -p $AGENT_PID -o comm=)"
    echo "   Memory: $(ps -p $AGENT_PID -o rss= | awk '{print $1/1024 " MB"}')"
    echo ""
    echo "🔍 To check agent logs:"
    echo "   tail -f agent.log"
    echo ""
    echo "🛑 To stop agent:"
    echo "   kill $AGENT_PID"
    echo "   or"
    echo "   pkill -f zerotrace"
    echo ""
    echo "📱 Note: Tray icon is disabled on macOS due to library compatibility issues"
    echo "   The agent runs in the background and logs to agent.log"
    echo ""
    echo "🌐 Check the web dashboard at: http://localhost:3000"
    echo "   (Make sure the API server is running on port 8080)"
else
    echo "❌ Agent failed to start or crashed"
    echo "Check agent.log for error details"
    exit 1
fi

