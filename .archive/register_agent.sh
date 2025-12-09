#!/bin/bash

# ZeroTrace Agent Registration Script
echo "🔧 Registering ZeroTrace Agent..."

# Generate a proper UUID for the agent
AGENT_ID=$(python3 -c "import uuid; print(uuid.uuid4())")
echo "Agent ID: $AGENT_ID"

# Register the agent
curl -X POST "http://localhost:8080/api/agents/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"agent_id\": \"$AGENT_ID\",
    \"name\": \"ZeroTrace Agent\",
    \"version\": \"1.0.0\",
    \"platform\": \"$(uname -s)\",
    \"architecture\": \"$(uname -m)\",
    \"status\": \"active\"
  }"

echo ""
echo "✅ Agent registration complete!"
echo "Agent ID: $AGENT_ID"
echo "You can now see this agent in the dashboard."
