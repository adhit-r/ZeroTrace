# Agentless Network Scanning - Quick Reference

## ⚠️ Important: Understanding "Agentless"

### The ZeroTrace Agent is the **SCANNING HOST**, not installed on targets!

```
┌─────────────────────────────────────┐
│  ZeroTrace Agent (Scanning Host)     │
│  Installed on: ONE machine           │
│  Role: Runs the scanner              │
└─────────────────────────────────────┘
              │
              │ Scans via network protocols
              │ (Nmap, Nuclei, HTTP requests)
              │ NO agent on targets!
              ▼
    ┌─────────────────────────────┐
    │  Target Devices             │
    │  • Routers                  │
    │  • Servers                  │
    │  • IoT Devices              │
    │  • Phones                   │
    │  • Switches                 │
    │                             │
    │  NO AGENT INSTALLED!        │
    └─────────────────────────────┘
```

## How It Works

1. **Install ZeroTrace Agent** on ONE machine (scanning host)
2. **Agent discovers local network** (e.g., 192.168.1.0/24)
3. **Agent scans all devices** using:
   - Nmap (network protocol scanning)
   - Nuclei (vulnerability scanning)
   - No installation needed on targets!
4. **Results sent to API**

## What Gets Scanned (Agentless)

✅ **Network Discovery**: All devices via network protocols
✅ **Port Scanning**: Open ports via TCP/UDP probes
✅ **Vulnerability Detection**: CVEs via HTTP/network requests
✅ **Configuration Auditing**: Insecure protocols, default creds
✅ **Service Detection**: Service identification via banners

**All WITHOUT installing anything on target devices!**

## Example

```bash
# Install agent on Machine A
./zerotrace-agent

# Agent automatically:
# 1. Discovers network: 192.168.1.0/24
# 2. Scans devices (agentless):
#    - Router (192.168.1.1) - NO agent on router
#    - Server (192.168.1.10) - NO agent on server
#    - IoT (192.168.1.50) - NO agent on IoT
# 3. Sends results to API
```

## Why This Architecture?

- ✅ Works on devices you don't control
- ✅ Works on IoT devices (can't run agents)
- ✅ Works on network equipment (switches/routers)
- ✅ Discovers unknown devices
- ✅ Similar to Tenable sensor architecture

## Documentation

See `docs/AGENTLESS_SCANNING_EXPLAINED.md` for detailed explanation.

