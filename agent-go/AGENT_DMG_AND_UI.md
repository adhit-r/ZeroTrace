# Agent DMG and UI - Quick Answers

## ✅ Yes, the Agent Has a UI!

The ZeroTrace agent includes a **system tray application** that works on macOS, Windows, and Linux.

### System Tray Features

**Menu Options:**
- 🔄 Agent Status
- 📊 CPU/Memory Usage
- 🔍 Manual Scan
- 🌐 Open Web Dashboard
- ⚙️ Settings
- ❌ Quit

**Icon States:**
- 🟢 Green = Connected
- ⚪ Gray = Checking
- 🔴 Red = Error

## 📦 DMG Installer

### Build the DMG

```bash
cd agent-go/mdm
./build-macos-dmg.sh
```

This creates: `dist/ZeroTrace-Agent-1.0.0.dmg`

### DMG Contents

- `zerotrace-agent` - Agent binary (with tray UI)
- `README.txt` - Installation instructions
- `Applications` link - Drag to install

### Installation

1. Double-click the DMG
2. Drag `zerotrace-agent` to Applications
3. Open Applications and run the agent
4. Tray icon appears in menu bar

## 🔍 About Nuclei - It's NOT Removed!

### Nuclei is Essential and Still Used!

**Nuclei is actively used** for vulnerability scanning. We simplified the implementation to use the CLI instead of the Go library.

### What Nuclei Does

- ✅ Scans for **thousands of CVEs**
- ✅ Finds **web vulnerabilities**
- ✅ Detects **misconfigurations**
- ✅ Identifies **exposed services**

### How It Works

```
Network Scan → Discovers Hosts → Nuclei Scans → Finds Vulnerabilities
```

### Why CLI Instead of Library?

- ✅ Always latest version
- ✅ All features available
- ✅ Automatic template updates
- ✅ Simpler to maintain

### Installing Nuclei

```bash
# macOS
brew install nuclei

# Linux/Windows
go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
```

## Quick Start

### 1. Build DMG with UI

```bash
cd agent-go/mdm
./build-macos-dmg.sh
```

### 2. Install from DMG

- Open the DMG
- Drag agent to Applications
- Run the agent
- Tray icon appears!

### 3. Use the UI

- Click tray icon for menu
- Check status, monitor resources
- Trigger manual scans
- Open web dashboard

## Summary

| Question | Answer |
|----------|--------|
| **Does agent have UI?** | ✅ Yes - System tray on all platforms |
| **Can I get a DMG?** | ✅ Yes - Run `./build-macos-dmg.sh` |
| **Is Nuclei removed?** | ❌ No - It's essential and actively used |
| **Why use Nuclei?** | ✅ Finds thousands of vulnerabilities |
| **How to use Nuclei?** | ✅ Installed automatically, used in scans |

## Documentation

- `docs/AGENT_UI_GUIDE.md` - Complete UI guide
- `docs/NUCLEI_EXPLANATION.md` - Why Nuclei is essential
- `agent-go/mdm/build-macos-dmg.sh` - DMG builder script

