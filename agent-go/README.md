# ZeroTrace Agent

Enterprise-grade vulnerability scanning agent for the ZeroTrace platform with MDM deployment support.

## Overview

ZeroTrace Agent is a universal vulnerability scanning agent designed for enterprise deployment via MDM platforms. It provides silent, background operation with automatic software discovery and vulnerability detection.

## Features

### Core Capabilities

- **Software Discovery**: Automatically scans installed applications
- **Vulnerability Detection**: CVE checking via Python enrichment service
- **Real-time Communication**: Sends results to ZeroTrace API
- **System Monitoring**: CPU/Memory usage tracking
- **Multi-tenant Support**: Organization isolation with enrollment tokens

### Enterprise Features

- **MDM Deployment**: Ready for Intune, Jamf Pro, Azure AD, Workspace ONE
- **Silent Operation**: No UI, background service operation
- **Universal Agent**: Single binary for all organizations
- **Enrollment System**: Secure token-based organization identification
- **System Service**: LaunchDaemon integration for macOS

## Quick Start

### Development Mode (with UI)

```bash
# Clone and navigate
cd agent-go

# Install dependencies
go mod tidy

# Set up environment
cp env.example .env
# Edit .env with your configuration

# Run with tray icon (development)
go run cmd/agent/main.go
```

### MDM Deployment Mode (silent)

```bash
# Build MDM package
cd mdm
./build-macos-pkg.sh

# Generated files:
# - ZeroTrace-Agent-1.0.0.pkg (installation package)
# - zerotrace-agent.mobileconfig (MDM configuration)
# - DEPLOYMENT_GUIDE.md (deployment instructions)
```

## MDM Deployment

### Supported Platforms

- **Microsoft Intune** (Windows/macOS)
- **Jamf Pro** (macOS)
- **Azure AD** (Enterprise)
- **VMware Workspace ONE** (UEM)

### Deployment Process

1. **Build Package**: `./mdm/build-macos-pkg.sh`
2. **Upload to MDM**: Package and configuration profile
3. **Configure Settings**: Enrollment tokens and API endpoints
4. **Deploy to Devices**: Silent installation via MDM

### Configuration Variables

```bash
# Required for enrollment
ZEROTRACE_ENROLLMENT_TOKEN=<token>
ZEROTRACE_API_URL=<api-url>
ZEROTRACE_ORG_ID=<org-id>

# Optional
ZEROTRACE_SCAN_INTERVAL=24h
ZEROTRACE_LOG_LEVEL=info
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AGENT_ID` | Unique agent identifier | Auto-generated |
| `API_URL` | ZeroTrace API endpoint | `http://localhost:8080` |
| `ENROLLMENT_TOKEN` | Organization enrollment token | Required for MDM |
| `ORGANIZATION_ID` | Organization identifier | Set after enrollment |
| `SCAN_INTERVAL` | Time between scans | `5m` |
| `SCAN_DEPTH` | Directory scan depth | `3` |
| `LOG_LEVEL` | Logging level | `info` |

### Scanning Configuration

```bash
# File size limits
MAX_FILE_SIZE=10MB

# Include patterns
INCLUDE_PATTERNS=*.go,*.py,*.js,*.java,*.php

# Exclude patterns
EXCLUDE_PATTERNS=.git,node_modules,.DS_Store,*.log
```

## Project Structure

```
agent-go/
├── cmd/
│   ├── agent/           # Full agent with tray (development)
│   └── agent-simple/    # MDM agent without tray (deployment)
├── internal/
│   ├── communicator/    # API communication
│   ├── config/         # Configuration management
│   ├── monitor/        # System monitoring
│   ├── processor/      # Data processing
│   ├── scanner/        # Software scanning
│   └── tray/           # Tray functionality (dev only)
├── mdm/                # MDM deployment tools
│   ├── build-macos-pkg.sh
│   ├── README.md
│   └── dist/           # Generated packages
└── pkg/                # Shared packages
```

## Development

### Building

```bash
# Development agent (with tray)
go build -o zerotrace-agent cmd/agent/main.go

# MDM agent (silent)
go build -o zerotrace-agent-simple cmd/agent-simple/main.go

# Cross-platform builds
GOOS=darwin GOARCH=amd64 go build -o zerotrace-agent cmd/agent-simple/main.go
GOOS=darwin GOARCH=arm64 go build -o zerotrace-agent cmd/agent-simple/main.go
```

### Testing

```bash
# Run tests
go test ./...

# Test specific components
go test ./internal/scanner
go test ./internal/communicator
```

## Monitoring & Troubleshooting

### Agent Status

```bash
# Check if agent is running
sudo launchctl list | grep zerotrace

# View logs
sudo log show --predicate 'process == "zerotrace-agent"' --last 1h

# Check configuration
sudo defaults read /Library/Preferences/com.zerotrace.agent
```

### MDM Status

```bash
# Check MDM enrollment
sudo profiles show -type configuration

# View MDM logs
sudo log show --predicate 'process == "mdm"' --last 1h
```

## API Integration

The agent communicates with the ZeroTrace API:

- **Enrollment**: POST `/api/enrollment/enroll`
- **Heartbeat**: POST `/api/agents/heartbeat`
- **Results**: POST `/api/agents/results`
- **Registration**: POST `/api/agents/register`

### Authentication

- **Enrollment**: Uses enrollment token (one-time)
- **Operations**: Uses agent credential (long-lived)
- **Legacy**: Uses API key (fallback)

### Code Examples

**Example: Agent Registration**
```go
package main

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type AgentRegistration struct {
    ID             string `json:"id"`
    OrganizationID string `json:"organization_id"`
    Name           string `json:"name"`
    Version        string `json:"version"`
    Hostname       string `json:"hostname"`
    OS             string `json:"os"`
}

func registerAgent(apiURL string, agent AgentRegistration) error {
    jsonData, _ := json.Marshal(agent)
    
    resp, err := http.Post(
        apiURL+"/api/agents/register",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

**Example: Sending Heartbeat**
```go
type Heartbeat struct {
    AgentID        string  `json:"agent_id"`
    OrganizationID string  `json:"organization_id"`
    AgentName      string  `json:"agent_name"`
    Status         string  `json:"status"`
    CPUUsage       float64 `json:"cpu_usage"`
    MemoryUsage    float64 `json:"memory_usage"`
    Timestamp      string  `json:"timestamp"`
}

func sendHeartbeat(apiURL string, heartbeat Heartbeat) error {
    jsonData, _ := json.Marshal(heartbeat)
    
    resp, err := http.Post(
        apiURL+"/api/agents/heartbeat",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

**Example: Submitting Scan Results**
```go
type ScanResult struct {
    AgentID  string      `json:"agent_id"`
    Results  []ScanData  `json:"results"`
    Metadata map[string]interface{} `json:"metadata"`
}

type ScanData struct {
    ScanType      string                `json:"scan_type"`
    Dependencies  []Dependency         `json:"dependencies"`
    Vulnerabilities []Vulnerability     `json:"vulnerabilities"`
    Timestamp     string                `json:"timestamp"`
}

func submitResults(apiURL string, result ScanResult) error {
    jsonData, _ := json.Marshal(result)
    
    resp, err := http.Post(
        apiURL+"/api/agents/results",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    return nil
}
```

## Agent Types

### Simple Agent (`cmd/agent-simple/`)

- **No UI**: Silent operation for MDM deployment
- **Background Service**: LaunchDaemon integration
- **Data Collection**: Software scanning and vulnerability detection
- **Enterprise Ready**: Perfect for large-scale deployment

### Full Agent (`cmd/agent/`)

- **Tray Icon**: Visual status indicator
- **Interactive Menu**: Status, CPU usage, quit options
- **Development Mode**: For testing and development
- **Visual Feedback**: Green/gray icon based on API connection

## Deployment Checklist

### Pre-Deployment

- [ ] Build agent packages
- [ ] Generate enrollment tokens
- [ ] Configure API endpoints
- [ ] Test in lab environment
- [ ] Create MDM policies

### Deployment

- [ ] Upload packages to MDM
- [ ] Configure installation policies
- [ ] Target device groups
- [ ] Deploy configuration profiles
- [ ] Monitor installation status

### Post-Deployment

- [ ] Verify agent enrollment
- [ ] Check API connectivity
- [ ] Monitor data collection
- [ ] Validate security policies
- [ ] Document deployment

## Next Steps

1. **Generate Enrollment Tokens**: Create tokens for each organization
2. **Upload to MDM**: Deploy packages to your MDM platform
3. **Configure Settings**: Set API endpoints and tokens
4. **Deploy to Devices**: Install on target devices
5. **Monitor Status**: Track installation and operation

## Support

- **Documentation**: `mdm/README.md`
- **Deployment Guide**: `mdm/dist/DEPLOYMENT_GUIDE.md`
- **Configuration**: `mdm/dist/zerotrace-agent.mobileconfig`
- **Summary**: `MDM_AGENT_SUMMARY.md`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

See [CONTRIBUTING.md](../CONTRIBUTING.md) for detailed guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.

---

**Last Updated**: January 2025
