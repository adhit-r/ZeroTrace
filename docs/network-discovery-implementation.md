# Network Discovery & Topology Visualization Implementation

## Overview

This document outlines the implementation of the network asset discovery and topology visualization features for ZeroTrace, based on the requirements for mapping deployed agents, discovering network assets, and providing visual topology maps.

## ✅ Implemented Features

### 1. Agent-Based Network Discovery

**Location**: `agent-go/internal/discovery/discovery.go`

**Capabilities**:
- **Local Network Scanning**: Discovers assets on all local network interfaces
- **Subnet Discovery**: Scans subnets for active hosts using ping
- **Host Information Gathering**:
  - IP and MAC address detection
  - Hostname resolution
  - OS detection (Windows vs Linux/Unix)
  - Port scanning (common ports: 21, 22, 23, 25, 53, 80, 110, 143, 443, etc.)
  - Service identification
  - Risk scoring based on open ports and services

**Key Methods**:
- `DiscoverLocalNetwork()`: Main discovery entry point
- `scanSubnet()`: Scans individual subnets
- `discoverHost()`: Gathers detailed information about individual hosts
- `calculateRiskScore()`: Calculates risk scores based on vulnerabilities

### 2. OWASP Amass Integration

**Location**: `enrichment-python/app/amass_discovery.py`

**Capabilities**:
- **External Asset Discovery**: Uses OWASP Amass for external domain enumeration
- **Subdomain Discovery**: Finds subdomains and related infrastructure
- **Geographic Information**: ASN, country, city data
- **Service Detection**: Port and service identification
- **Concurrent Scanning**: Supports multiple domain scanning

**Key Features**:
- Async/await support for non-blocking operations
- JSON output parsing
- Error handling and retry logic
- Version checking and installation validation

### 3. Data Models

**Location**: `agent-go/internal/models/models.go`

**New Models Added**:
- `NetworkAsset`: Represents discovered network assets
- `PortInfo`: Information about open ports
- `ServiceInfo`: Running service details
- `PeerInfo`: Connected peer information
- `AmassResult`: OWASP Amass discovery results
- `NetworkTopology`: Complete topology structure
- `TopologyNode`: Individual nodes in the topology
- `TopologyLink`: Connections between nodes
- `Cluster`: Logical groupings of assets
- `AgentTelemetry`: Real-time agent status

### 4. React Frontend Visualization

**Location**: `web-react/src/components/NetworkTopology.tsx`

**Features**:
- **D3.js Integration**: Force-directed graph visualization
- **Interactive Nodes**: Clickable nodes with detailed tooltips
- **Real-time Updates**: Live data refresh capabilities
- **Multiple View Modes**: Network, floor plan, geographic, cluster views
- **Filtering**: By agent type, connection status
- **Export Functionality**: JSON export of topology data
- **Responsive Design**: Works on different screen sizes

**Visual Elements**:
- Color-coded nodes by type (agents, assets, amass discoveries)
- Node size based on risk score
- Connection lines with different colors for scan/network/external
- Glow effects and hover states
- Real-time statistics panel

### 5. Topology Page

**Location**: `web-react/src/pages/Topology.tsx`

**Features**:
- Full-screen topology visualization
- Mock data for demonstration
- Loading states and error handling
- Node click handlers for detailed views
- Auto-refresh simulation

## 🔧 Technical Implementation Details

### Agent Discovery Process

1. **Interface Detection**: Scans all network interfaces
2. **Subnet Identification**: Identifies IPv4 subnets
3. **Host Discovery**: Pings hosts in each subnet
4. **Asset Profiling**: Gathers detailed information about active hosts
5. **Risk Assessment**: Calculates risk scores
6. **Data Reporting**: Sends results to API

### Amass Integration Process

1. **Domain Input**: Receives domains to scan
2. **Amass Execution**: Runs OWASP Amass enumeration
3. **Output Parsing**: Parses JSON output files
4. **Data Extraction**: Extracts IPs, subdomains, services
5. **Result Processing**: Formats data for API consumption

### Frontend Visualization

1. **Data Loading**: Receives topology data from API
2. **D3 Setup**: Initializes force simulation
3. **Node Rendering**: Creates interactive SVG nodes
4. **Link Rendering**: Draws connection lines
5. **Event Handling**: Manages clicks, hovers, drags
6. **Real-time Updates**: Refreshes data periodically

## 🚀 Usage Instructions

### Starting the Services

```bash
# Using Podman (recommended)
podman-compose up -d

# Access the topology visualization
open http://localhost:3000/topology
```

### Agent Configuration

The agent will automatically:
- Discover local network assets
- Send telemetry data to the API
- Update topology information

### Amass Configuration

1. Install OWASP Amass:
   ```bash
   # macOS
   brew install amass
   
   # Linux
   sudo apt-get install amass
   ```

2. The Python enrichment service will automatically detect and use Amass

### Frontend Features

1. **View Modes**: Switch between different visualization modes
2. **Filtering**: Filter by agent type or connection status
3. **Node Interaction**: Click nodes for detailed information
4. **Data Export**: Export topology data as JSON
5. **Real-time Refresh**: Click refresh button for latest data

## 📊 Data Flow

```
Agent Discovery → API → Database → Frontend Visualization
     ↓
Amass Discovery → Enrichment Service → API → Database
     ↓
Real-time Updates → WebSocket → Frontend
```

## 🔮 Future Enhancements

### Planned Features

1. **Enhanced OS Detection**: More sophisticated OS fingerprinting
2. **Service Banner Grabbing**: Detailed service version detection
3. **Vulnerability Correlation**: Link discovered assets to vulnerabilities
4. **Geographic Visualization**: Map-based topology view
5. **Historical Tracking**: Asset discovery history and trends
6. **Alert System**: Notifications for new asset discoveries
7. **Integration APIs**: Connect with other security tools

### Performance Optimizations

1. **Parallel Scanning**: Concurrent subnet scanning
2. **Caching**: Cache discovery results
3. **Incremental Updates**: Only update changed assets
4. **WebSocket Optimization**: Efficient real-time updates
5. **Frontend Virtualization**: Handle large numbers of nodes

## 🛡️ Security Considerations

1. **Network Scanning**: Respect network policies and rate limits
2. **Data Privacy**: Ensure discovered asset data is properly secured
3. **Access Control**: Implement proper authentication for topology access
4. **Audit Logging**: Log all discovery activities
5. **Compliance**: Ensure compliance with data protection regulations

## 📝 API Endpoints

### Network Assets
- `GET /api/assets` - List discovered assets
- `POST /api/assets` - Report new asset discovery
- `GET /api/assets/{id}` - Get asset details
- `PUT /api/assets/{id}` - Update asset information

### Topology
- `GET /api/topology` - Get current topology
- `POST /api/topology/refresh` - Trigger topology refresh
- `GET /api/topology/export` - Export topology data

### Amass Results
- `GET /api/amass/results` - Get Amass discovery results
- `POST /api/amass/scan` - Trigger new Amass scan
- `GET /api/amass/status` - Get scan status

## 🎯 Success Metrics

1. **Asset Discovery Rate**: Percentage of network assets discovered
2. **Discovery Accuracy**: Correctness of asset information
3. **Real-time Performance**: Speed of topology updates
4. **User Engagement**: Usage of visualization features
5. **Security Coverage**: Percentage of assets with security monitoring

This implementation provides a solid foundation for network asset discovery and topology visualization, with room for expansion and enhancement based on specific requirements and use cases.
