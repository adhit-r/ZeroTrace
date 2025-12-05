# Network Topology Visualizer - n8n Style

## Overview

The Network Flow Visualizer provides an interactive, n8n-style node-based visualization of your network infrastructure. It displays devices, services, and vulnerabilities in an intuitive flow diagram.

## Features

### Visual Elements

1. **Device Nodes** (Rectangular cards)
   - Color-coded by risk score (green/yellow/red)
   - Shows device type icon (switch, router, IoT, phone, server)
   - Displays IP address, OS, risk score, and vulnerability count
   - Status indicators (online/offline)

2. **Vulnerability Nodes** (Colored cards by severity)
   - Critical: Red
   - High: Orange
   - Medium: Yellow
   - Low: Blue
   - Shows CVE, description, and affected host

3. **Service Nodes** (Purple cards)
   - Shows service name and port
   - Connected to their host devices

4. **Connections** (Animated edges)
   - Devices → Services (purple)
   - Devices → Vulnerabilities (severity-colored)
   - Animated flow for better visualization

### Interactive Features

- **Zoom & Pan**: Use mouse wheel to zoom, drag to pan
- **Node Selection**: Click nodes to see details
- **Search**: Filter nodes by name or IP
- **Type Filter**: Filter by device/vulnerability/service
- **Export**: Download topology as JSON
- **Mini Map**: Overview of entire topology

## Usage

### Access the Visualizer

Navigate to `/topology` or `/network-topology` in your app.

### Controls

**Left Panel:**
- Search bar: Find devices by name or IP
- Type filter: Show only specific node types
- Statistics: Node and connection counts

**Top Right:**
- Export button: Download topology data
- Refresh button: Reload network data

**Bottom Right:**
- Node details panel (when node is selected)

### Node Interaction

1. **Click a node** to see detailed information
2. **Drag nodes** to rearrange layout
3. **Hover over edges** to see connection details
4. **Use minimap** to navigate large topologies

## Data Source

The visualizer fetches data from:
- Agent metadata: `agent.metadata.network_scan_result`
- Network scan findings: Ports, vulnerabilities, configurations
- Device classifications: Switch, router, IoT, phone, server

## Customization

### Node Styling

Edit `NetworkFlowVisualizer.tsx` to customize:
- Node colors and sizes
- Icons and labels
- Risk score thresholds
- Layout algorithms

### Layout Options

Currently uses grid-based positioning. Can be enhanced with:
- Force-directed layout
- Hierarchical layout
- Custom positioning algorithms

## Integration with Network Scanner

The visualizer automatically updates when:
- Network scans complete
- New devices are discovered
- Vulnerabilities are found
- Configuration errors are detected

## Future Enhancements

1. **Workflow Builder**: Create remediation workflows (n8n-style)
2. **Real-time Updates**: WebSocket for live topology changes
3. **Grouping**: Cluster devices by subnet/VLAN
4. **Timeline View**: See topology changes over time
5. **Export Formats**: PNG, SVG, PDF exports
6. **Custom Layouts**: Save and load custom arrangements

