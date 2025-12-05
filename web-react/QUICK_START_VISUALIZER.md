# Quick Start - Network Topology Visualizer

## ✅ Implementation Complete!

An n8n-style network topology visualizer has been added to your app.

## Access the Visualizer

1. **Start the frontend:**
   ```bash
   cd web-react
   npm run dev
   ```

2. **Navigate to:**
   - `/topology` - Network topology page
   - `/network-topology` - Alternative route

## What You'll See

### Node Types

1. **Device Nodes** (Rectangular cards)
   - 🟢 Green: Low risk (0-40)
   - 🟡 Yellow: Medium risk (40-70)
   - 🔴 Red: High risk (70-100)
   - Shows: IP, device type, OS, risk score, vulnerabilities

2. **Vulnerability Nodes** (Colored by severity)
   - 🔴 Critical
   - 🟠 High
   - 🟡 Medium
   - 🔵 Low

3. **Service Nodes** (Purple)
   - Shows open ports and services

### Interactive Features

- **Zoom**: Mouse wheel or controls
- **Pan**: Click and drag
- **Search**: Filter by device name or IP
- **Filter**: Show only devices/vulnerabilities/services
- **Click Nodes**: See detailed information
- **Export**: Download topology as JSON

## Data Flow

```
Network Scanner → Agent Metadata → API → Visualizer
```

The visualizer automatically fetches:
- Network scan results from agent metadata
- Device classifications (switch/router/IoT/phone/server)
- Vulnerabilities from Nuclei scans
- Configuration findings
- Service information

## Testing

1. **Run network scan** on your agent
2. **Wait for results** to be sent to API
3. **Open visualizer** at `/topology`
4. **See your network** visualized!

## Customization

Edit `NetworkFlowVisualizer.tsx` to:
- Change node colors
- Add new node types
- Modify layout algorithms
- Add custom interactions

## Future: Workflow Builder

Coming soon - n8n-style workflow builder for:
- Creating remediation workflows
- Automating security responses
- Multi-agent orchestration
- Custom automation rules

See `docs/OPENAGENTS_INTEGRATION_PLAN.md` for details.

