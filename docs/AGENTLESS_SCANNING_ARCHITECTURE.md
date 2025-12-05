# Agentless Network Scanning Architecture

## Understanding the Architecture

### Key Distinction

**Agentless Scanner** = No agent installed on **target devices** being scanned
**Scanning Host** = The ZeroTrace agent that **runs** the scanner

## How It Works

```
┌─────────────────────────────────────────────────────────────┐
│                    ZeroTrace Agent (Scanning Host)           │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Network Scanner (Agentless)                          │  │
│  │  - Uses Nmap to discover devices                      │  │
│  │  - Uses Nuclei to scan vulnerabilities                │  │
│  │  - No installation needed on targets                   │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                          │
                          │ Network Requests
                          │ (No agent on targets)
                          ▼
        ┌─────────────────────────────────────┐
        │   Target Devices (No Agent Needed)  │
        │                                     │
        │  • Switches                         │
        │  • Routers                          │
        │  • IoT Devices                      │
        │  • Phones                           │
        │  • Servers                          │
        │  • Workstations                     │
        └─────────────────────────────────────┘
```

## Agentless Scanning Methods

### 1. Network Discovery (Nmap)
- **SYN Scan**: Sends TCP SYN packets, no agent needed
- **Port Scanning**: Discovers open ports without installing anything
- **OS Detection**: Fingerprints OS from network responses
- **Service Detection**: Identifies services from banners

### 2. Vulnerability Scanning (Nuclei)
- **HTTP Requests**: Scans web services via HTTP/HTTPS
- **Template-Based**: Uses YAML templates, no agent needed
- **Banner Grabbing**: Reads service banners
- **CVE Detection**: Matches known vulnerabilities

### 3. Configuration Auditing
- **Banner Analysis**: Reads service banners
- **Protocol Analysis**: Detects insecure protocols (Telnet, FTP)
- **SNMP Queries**: If credentials provided (still agentless)
- **HTTP Interface Checks**: Scans web management interfaces

## Why Use an Agent as Scanning Host?

### Similar to Tenable Architecture

```
Tenable Sensor/Agent (Scanning Host)
    │
    ├─→ Scans network devices (agentless)
    ├─→ Scans servers (agentless)
    ├─→ Scans IoT devices (agentless)
    └─→ Sends results to Tenable.io
```

```
ZeroTrace Agent (Scanning Host)
    │
    ├─→ Scans network devices (agentless)
    ├─→ Scans servers (agentless)
    ├─→ Scans IoT devices (agentless)
    └─→ Sends results to ZeroTrace API
```

### Benefits

1. **Centralized Scanning**: One agent can scan entire network
2. **No Target Installation**: No need to install on every device
3. **Network Visibility**: Discovers devices you don't control
4. **Credential-Based Deep Scan**: Optional credentials for authenticated scanning
5. **Continuous Monitoring**: Agent runs scans on schedule

## What Gets Scanned (Agentless)

### ✅ Scanned Without Agent on Target

- **Network Discovery**: All devices on network
- **Port Scanning**: Open ports and services
- **Vulnerability Detection**: Known CVEs and misconfigurations
- **Banner Grabbing**: Service identification
- **OS Fingerprinting**: Operating system detection
- **Configuration Errors**: Insecure protocols, default credentials

### 🔐 Optional: Authenticated Scanning

If credentials provided (still agentless):
- **SSH**: Connect and check configurations
- **SNMP**: Query device information
- **HTTP**: Access management interfaces
- **Database**: Check database configurations

## Testing the Agentless Scanner

### Current Setup

The ZeroTrace agent **IS** the scanning host. It:
1. Runs on a machine in your network
2. Performs agentless scans of other devices
3. Sends results to the API

### To Test

1. **Install ZeroTrace Agent** on one machine (scanning host)
2. **Run network scan** - it will discover other devices
3. **View results** - see all discovered devices (no agents on them)

### Example Scenario

```
Your Network:
├─ Machine A: ZeroTrace Agent installed (scanning host)
├─ Machine B: No agent (will be scanned)
├─ Router: No agent (will be scanned)
├─ IoT Device: No agent (will be scanned)
└─ Server: No agent (will be scanned)

Agent on Machine A scans B, Router, IoT, Server → Agentless!
```

## Comparison: Agent-Based vs Agentless

### Agent-Based Scanning (NOT what we're doing)
```
Target Device
├─ Install agent
├─ Agent reports to central server
└─ Requires installation on every device
```

### Agentless Scanning (What we're doing)
```
Scanning Host (ZeroTrace Agent)
├─ Scans network without installing on targets
├─ Discovers devices via network protocols
└─ No installation needed on target devices
```

## Why This Architecture?

1. **Discover Unknown Devices**: Can't install agents on devices you don't control
2. **IoT Devices**: Many IoT devices can't run agents
3. **Network Devices**: Switches/routers typically don't support agents
4. **Quick Deployment**: Scan entire network from one host
5. **Compliance**: Some environments don't allow agent installation

## Summary

- ✅ **Agentless**: No agent needed on target devices
- ✅ **Scanning Host**: ZeroTrace agent runs the scanner
- ✅ **Network Discovery**: Finds all devices via network protocols
- ✅ **Vulnerability Scanning**: Scans without installing anything
- ✅ **Similar to Tenable**: Sensor/agent performs agentless scans

The ZeroTrace agent is like a Tenable sensor - it performs agentless scans of your network without requiring installation on target devices.

