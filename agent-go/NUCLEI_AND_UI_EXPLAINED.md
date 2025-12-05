# Nuclei & UI - Quick Answers

## ✅ Yes, the Agent Has a UI!

The ZeroTrace agent includes a **system tray application** that works on macOS, Windows, and Linux.

### System Tray Features

**Menu Options:**
- 🔄 Agent Status - Check connection and health
- 📊 CPU/Memory Usage - Real-time resource monitoring
- 🔍 Check Now - Manual vulnerability scan
- 🌐 Open Web Dashboard - Launch web UI
- ⚙️ Settings - Configure agent
- ❌ Quit - Stop agent

**Icon States:**
- 🟢 Green = Connected to API
- ⚪ Gray = Checking/Initializing
- 🔴 Red = Error/Disconnected

### How to Use

1. **Run the agent** - Tray icon appears in menu bar
2. **Click icon** - See menu with options
3. **Monitor status** - CPU/Memory updates in real-time
4. **Trigger scans** - Manual vulnerability checks
5. **Open dashboard** - Launch web interface

## 🔍 Nuclei is NOT Removed - It's Essential!

### Nuclei is Still Used and Very Useful!

**Nuclei is actively used** for vulnerability scanning. We use the CLI version instead of the Go library.

### What Nuclei Does

- ✅ Scans for **thousands of CVEs**
- ✅ Finds **web vulnerabilities**
- ✅ Detects **misconfigurations**
- ✅ Identifies **exposed services**
- ✅ Checks **default credentials**
- ✅ Tests **API endpoints**

### How It Works in ZeroTrace

```
Network Scan Flow:
1. Nmap discovers devices and ports
2. Nuclei scans discovered hosts for vulnerabilities
3. Results combined and sent to API
```

### Why CLI Instead of Library?

- ✅ Always latest version (auto-updates)
- ✅ All features available
- ✅ Automatic template updates
- ✅ Simpler to maintain
- ✅ Better performance

### Where Nuclei is Used

**In `network_scanner.go`:**
```go
// Step 4: Run Nuclei vulnerability scanning on discovered hosts
nucleiFindings, err := ns.nucleiScanner.ScanTargets(targets)
```

**Scan Methods:**
- `nmap+nuclei` - Primary method
- `naabu+nuclei` - Fallback method

### Installing Nuclei

```bash
# macOS
brew install nuclei

# Linux/Windows
go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest

# Or download from:
# https://github.com/projectdiscovery/nuclei/releases
```

### Nuclei Templates

Nuclei uses YAML templates to scan for vulnerabilities:
- **CVE templates** - Known vulnerabilities
- **Misconfiguration templates** - Security issues
- **Exposed services** - Open services
- **Default credentials** - Weak passwords

Templates are automatically updated when you run Nuclei.

## Summary

| Question | Answer |
|----------|--------|
| **Does agent have UI?** | ✅ Yes - System tray on all platforms |
| **Is Nuclei removed?** | ❌ No - Still used for vulnerability scanning |
| **Is Nuclei useful?** | ✅ Yes - Essential for finding CVEs and vulnerabilities |
| **How is Nuclei used?** | CLI version - scans discovered hosts after Nmap |
| **Can I get a DMG?** | ✅ Yes - Run `./build-macos-dmg.sh` |

## Building the DMG

```bash
cd agent-go/mdm
./build-macos-dmg.sh
```

This creates: `dist/ZeroTrace-Agent-1.0.0.dmg`

The DMG includes:
- Agent binary (with tray UI)
- Installation instructions
- Applications link

## Quick Start

1. **Build DMG**: `cd agent-go/mdm && ./build-macos-dmg.sh`
2. **Install**: Open DMG, drag to Applications
3. **Run**: Agent starts with tray icon
4. **Use UI**: Click tray icon for menu
5. **Nuclei scans**: Automatically runs on network scans


