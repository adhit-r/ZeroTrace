# ZeroTrace Agent UI Guide

## ✅ Yes, the Agent Has a UI!

The ZeroTrace agent includes a **system tray application** that provides a user-friendly interface on macOS, Windows, and Linux.

## System Tray Features

### macOS Tray Icon

The agent shows a tray icon in the menu bar with:

**Menu Options:**
- 🔄 **Agent Status** - Check connection and scan status
- 📊 **CPU Usage** - Real-time CPU monitoring
- 💾 **Memory Usage** - Real-time memory monitoring
- 🔍 **Check Now** - Manual vulnerability scan
- 🌐 **Open Web UI** - Launch dashboard in browser
- 🔄 **Restart Agent** - Restart the agent
- ⚙️ **Settings** - Open configuration
- ❌ **Quit** - Stop the agent

**Icon States:**
- 🟢 **Green** - Connected and running
- ⚪ **Gray** - Disconnected or checking
- 🔴 **Red** - Error state

### Windows Tray Icon

Similar features with Windows-native notifications and controls.

### Linux Tray Icon

Works with most desktop environments (GNOME, KDE, XFCE).

## How to Use

### Starting the Agent

```bash
# The agent automatically starts the tray UI
./zerotrace-agent
```

### Accessing the Tray

1. **macOS**: Look for icon in menu bar (top right)
2. **Windows**: Look for icon in system tray (bottom right)
3. **Linux**: Look for icon in system tray/notification area

### Tray Menu Actions

**Check Status:**
- Click "Agent Status" to see connection status
- Shows agent ID, hostname, and API connection

**Monitor Resources:**
- CPU and Memory usage updated in real-time
- Click to see detailed information

**Manual Scan:**
- Click "Check Now" to trigger immediate scan
- Useful for testing or on-demand scanning

**Open Dashboard:**
- Click "Open Web UI" to launch browser
- Opens the ZeroTrace web dashboard

**Restart Agent:**
- Click "Restart Agent" to restart
- Useful after configuration changes

**Settings:**
- Click "Settings" to open config file
- Edit `.env` file for configuration

## UI Components

### Full Tray Manager

**Location:** `agent-go/internal/tray/tray.go`

**Features:**
- Real-time monitoring
- Manual scan triggers
- Web UI launcher
- Agent restart
- Settings editor

### Simple Tray Manager

**Location:** `agent-go/internal/tray/simple_tray.go`

**Features:**
- Minimal interface
- Status only
- Basic system info
- MDM-friendly

## Integration

The tray UI is **automatically started** when you run the agent:

```go
// agent-go/cmd/agent/main.go

// Initialize tray manager
trayManager := tray.NewTrayManager(cfg)
trayManager.Start()
```

## Building with UI

### Standard Build (with UI)
```bash
go build -o zerotrace-agent cmd/agent/main.go
```

### Simple Build (no UI, for MDM)
```bash
go build -o zerotrace-agent cmd/agent-simple/main.go
```

## DMG Installer

The DMG includes the agent **with tray UI**:

```bash
cd agent-go/mdm
./build-macos-dmg.sh
```

This creates: `dist/ZeroTrace-Agent-1.0.0.dmg`

### DMG Contents

- `zerotrace-agent` - Agent binary (with UI)
- `README.txt` - Installation instructions
- `Applications` link - For easy installation

## Platform Support

| Platform | Tray UI | Notifications | Status |
|----------|---------|---------------|--------|
| macOS    | ✅ Yes  | ✅ Yes        | ✅ Full Support |
| Windows  | ✅ Yes  | ✅ Yes        | ✅ Full Support |
| Linux    | ✅ Yes  | ✅ Yes        | ✅ Full Support |

## Troubleshooting

### Tray Icon Not Showing

**macOS:**
- Check if agent is running: `ps aux | grep zerotrace`
- Check menu bar permissions in System Preferences
- Restart the agent

**Windows:**
- Check system tray settings
- Ensure agent has proper permissions
- Check Windows Defender/antivirus

**Linux:**
- Ensure desktop environment supports system tray
- Check notification daemon is running
- Verify permissions

### UI Not Responding

1. Check agent logs: `~/.zerotrace/agent.log`
2. Restart the agent
3. Check API connectivity
4. Verify configuration

## Customization

### Change Icon

Edit: `agent-go/internal/tray/icons.go`

### Modify Menu

Edit: `agent-go/internal/tray/tray.go`

### Add Features

Extend: `agent-go/internal/tray/platform.go`

## Summary

✅ **Agent has full UI** - System tray on all platforms
✅ **Easy to use** - Click menu for actions
✅ **Real-time monitoring** - CPU, memory, status
✅ **User-friendly** - No command line needed
✅ **DMG available** - Easy installation

The agent is **not just a background service** - it has a complete user interface!

