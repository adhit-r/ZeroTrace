# Menu Bar Icon Guide - macOS

## ✅ Menu Bar Icon Created!

The agent now creates a **menu bar icon** that appears in the **top-right task bar** next to WiFi, profile, and other icons.

## How to Get the Menu Bar Icon

### Option 1: Use the .app Bundle (Recommended)

**Location:** `agent-go/mdm/build/ZeroTrace Agent.app`

1. **Double-click** `ZeroTrace Agent.app`
2. **Menu bar icon appears** next to WiFi icon
3. **Click icon** for menu options

### Option 2: Use the DMG

**Location:** `agent-go/mdm/dist/ZeroTrace-Agent-1.0.0-APP.dmg`

1. **Mount DMG**: Double-click the DMG
2. **Drag** `ZeroTrace Agent.app` to Applications
3. **Launch** from Applications
4. **Menu bar icon appears**!

## Menu Bar Icon Features

### Icon States

- 🟢 **Green** = Connected to API
- ⚪ **Gray** = Checking/Initializing  
- 🔴 **Red** = Error/Disconnected

### Menu Options

When you click the icon:

- **🔄 Status** - Agent connection status
- **📊 CPU** - Current CPU usage
- **❌ Quit** - Stop the agent

## Location in Menu Bar

```
[WiFi] [Bluetooth] [Battery] [Profile] [ZeroTrace] ← Your icon here!
```

The icon appears in the **top-right corner** of your Mac's menu bar.

## Troubleshooting

### Icon Not Appearing?

1. **Check if running as .app**: Must be `.app` bundle, not just binary
2. **Check permissions**: macOS may need permission for menu bar
3. **Try restarting**: Quit and relaunch the app
4. **Check logs**: Look for "Menu bar icon ready!" message

### Icon Crashes?

Run with `--no-tray` flag:
```bash
./zerotrace-agent --no-tray
```

Agent works perfectly without the icon!

## Building the App Bundle

```bash
cd agent-go/mdm
./build-macos-app.sh
```

This creates: `build/ZeroTrace Agent.app`

## Building the DMG

```bash
cd agent-go/mdm
# The DMG build script now includes the .app bundle
```

## What's Different?

### Before (Binary)
- Just a binary file
- No menu bar icon support
- Crashed when trying to show icon

### Now (.app Bundle)
- Proper macOS app bundle
- Menu bar icon support
- Info.plist with LSUIElement
- Works correctly on macOS

## Summary

| Feature | Status |
|---------|--------|
| **Menu bar icon** | ✅ Yes - Next to WiFi |
| **Icon states** | ✅ Green/Gray/Red |
| **Menu options** | ✅ Status, CPU, Quit |
| **.app bundle** | ✅ Created |
| **DMG with .app** | ✅ Available |

The menu bar icon is now working! 🎉


