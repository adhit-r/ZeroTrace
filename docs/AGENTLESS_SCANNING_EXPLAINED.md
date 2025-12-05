# Agentless Scanning Explained

## The Confusion: "Why use an agent for agentless scanning?"

### ✅ The Answer: Two Different Roles

**ZeroTrace Agent** = **Scanning Host** (like a Tenable sensor)
**Target Devices** = **Scanned Devices** (NO agent needed)

## Visual Architecture

```
┌─────────────────────────────────────────────────────────────┐
│  ZeroTrace Agent (Scanning Host)                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Network Scanner (Agentless)                          │  │
│  │  • Uses Nmap (network protocol)                       │  │
│  │  • Uses Nuclei (HTTP requests)                        │  │
│  │  • NO installation on targets                          │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                    │
                    │ Network Requests Only
                    │ (No agent on targets)
                    ▼
    ┌──────────────────────────────────────────┐
    │  Target Devices (NO AGENT INSTALLED)     │
    │                                          │
    │  • Router (192.168.1.1)                 │
    │  • Server (192.168.1.10)                │
    │  • IoT Device (192.168.1.50)             │
    │  • Phone (192.168.1.100)                 │
    │  • Switch (192.168.1.200)                │
    │                                          │
    │  All scanned WITHOUT installing anything │
    └──────────────────────────────────────────┘
```

## How Agentless Scanning Works

### 1. Network Discovery (Nmap)
```
ZeroTrace Agent sends:
  → TCP SYN packets
  → ICMP ping
  → ARP requests

Target Device responds:
  ← TCP SYN-ACK (if port open)
  ← ICMP reply
  ← ARP reply

NO AGENT INSTALLED on target!
```

### 2. Vulnerability Scanning (Nuclei)
```
ZeroTrace Agent sends:
  → HTTP GET requests
  → HTTP POST requests
  → Service probes

Target Device responds:
  ← HTTP responses
  ← Service banners
  ← Error messages

NO AGENT INSTALLED on target!
```

### 3. Configuration Auditing
```
ZeroTrace Agent:
  → Reads service banners
  → Detects insecure protocols
  → Checks for default credentials (if provided)

Target Device:
  ← Sends banners
  ← Responds to protocol queries

NO AGENT INSTALLED on target!
```

## Why We Need a Scanning Host

### Similar to Tenable Architecture

**Tenable:**
```
Tenable Sensor/Agent (Scanning Host)
    │
    ├─→ Scans network (agentless)
    ├─→ Scans servers (agentless)
    └─→ Sends results to Tenable.io
```

**ZeroTrace:**
```
ZeroTrace Agent (Scanning Host)
    │
    ├─→ Scans network (agentless)
    ├─→ Scans servers (agentless)
    └─→ Sends results to ZeroTrace API
```

### Why Not Scan from API Server?

1. **Network Access**: API server may not have access to internal networks
2. **Distributed Scanning**: Multiple agents can scan different networks
3. **Local Network Discovery**: Agent knows its local network
4. **Performance**: Scanning from within network is faster
5. **Security**: No need to expose internal networks to API

## Real-World Example

### Scenario: Office Network

```
Office Network (192.168.1.0/24)

┌─────────────────────────────────────────┐
│  Machine A: ZeroTrace Agent installed   │
│  IP: 192.168.1.5                        │
│  Role: Scanning Host                    │
└─────────────────────────────────────────┘
              │
              │ Agentless Scans
              │
    ┌─────────┴─────────┬───────────┬──────────┐
    │                   │           │          │
┌───▼───┐         ┌─────▼───┐  ┌───▼───┐  ┌──▼────┐
│Router │         │ Server  │  │ IoT   │  │Phone │
│1.1    │         │ 1.10     │  │ 1.50  │  │1.100 │
│       │         │          │  │       │  │      │
│NO     │         │NO        │  │NO     │  │NO    │
│AGENT  │         │AGENT     │  │AGENT  │  │AGENT │
└───────┘         └──────────┘  └───────┘  └──────┘
```

**What Happens:**
1. Agent on Machine A scans 192.168.1.0/24
2. Discovers Router, Server, IoT, Phone
3. **NO agent installed** on any of them
4. All scanning done via network protocols
5. Results sent to API

## Comparison: Agent-Based vs Agentless

### ❌ Agent-Based (NOT what we're doing)
```
Target Device
├─ Install ZeroTrace agent
├─ Agent reports to API
└─ Requires installation on EVERY device
```

**Problems:**
- Can't install on devices you don't control
- IoT devices often can't run agents
- Network devices don't support agents
- Requires access to install

### ✅ Agentless (What we're doing)
```
Scanning Host (ZeroTrace Agent)
├─ Scans network via protocols
├─ No installation on targets
└─ Discovers ALL devices
```

**Benefits:**
- Works on any network device
- No installation needed
- Discovers unknown devices
- Works on IoT devices
- Works on network equipment

## What Gets Scanned (Agentless)

### ✅ Scanned WITHOUT Agent on Target

- **Network Discovery**: All devices via ARP, ping, port scan
- **Port Scanning**: Open ports via TCP/UDP probes
- **Service Detection**: Service identification via banners
- **OS Detection**: OS fingerprinting via network responses
- **Vulnerability Scanning**: CVEs via HTTP/network probes
- **Configuration Errors**: Insecure protocols, default creds

### 🔐 Optional: Authenticated Scanning

If credentials provided (STILL agentless):
- **SSH**: Connect and check configs (no agent install)
- **SNMP**: Query device info (no agent install)
- **HTTP**: Access management interfaces (no agent install)
- **Database**: Check database configs (no agent install)

## Testing the Agentless Scanner

### Setup

1. **Install ZeroTrace Agent** on ONE machine (scanning host)
2. **Run network scan** - it discovers OTHER devices
3. **View results** - see all discovered devices

### Example Test

```bash
# On Machine A (has ZeroTrace Agent)
./zerotrace-agent

# Agent automatically:
# 1. Discovers local network (192.168.1.0/24)
# 2. Scans all devices (agentless)
# 3. Finds:
#    - Router at 192.168.1.1 (NO agent on router)
#    - Server at 192.168.1.10 (NO agent on server)
#    - IoT at 192.168.1.50 (NO agent on IoT)
# 4. Sends results to API
```

## Summary

| Component | Role | Agent Needed? |
|-----------|------|---------------|
| **ZeroTrace Agent** | Scanning Host | ✅ Yes (runs scanner) |
| **Target Devices** | Scanned Devices | ❌ No (scanned via network) |
| **Router** | Target Device | ❌ No |
| **Server** | Target Device | ❌ No |
| **IoT Device** | Target Device | ❌ No |
| **Phone** | Target Device | ❌ No |

**Key Point:** The ZeroTrace agent is the **scanning host**, not something installed on target devices. It performs **agentless scans** of other devices on the network.

## Why This Architecture?

1. **Discover Unknown Devices**: Can't install agents on devices you don't know about
2. **IoT Devices**: Many can't run agents
3. **Network Equipment**: Switches/routers typically don't support agents
4. **Quick Deployment**: Scan entire network from one host
5. **Compliance**: Some environments don't allow agent installation
6. **Similar to Tenable**: Industry-standard approach

The ZeroTrace agent is like a **Tenable sensor** - it performs **agentless scans** of your network without requiring installation on target devices.

