# Quick Start - Menu Bar Icon

## ✅ Menu Bar Icon Ready!

The agent now shows a **menu bar icon** in the **top-right task bar** next to WiFi.

## How to Use

### Step 1: Get the App

**Option A: Use the .app Bundle**
```bash
cd agent-go/mdm/build
open "ZeroTrace Agent.app"
```

**Option B: Use the DMG**
```bash
# Mount the DMG
open agent-go/mdm/dist/ZeroTrace-Agent-1.0.0-APP.dmg

# Drag "ZeroTrace Agent.app" to Applications
# Then launch from Applications
```

### Step 2: Launch

1. **Double-click** `ZeroTrace Agent.app`
2. **Menu bar icon appears** next to WiFi icon
3. **Look for icon** in top-right menu bar

### Step 3: Use the Icon

**Click the icon** to see:
- 🔄 Status - Connection status
- 📊 CPU - CPU usage
- ❌ Quit - Stop agent

## Icon Location

```
Top-Right Menu Bar:
[WiFi] [Bluetooth] [Battery] [Profile] [ZeroTrace] ← Your icon!
```

## Icon Colors

- 🟢 **Green** = Connected to API
- ⚪ **Gray** = Checking/Initializing
- 🔴 **Red** = Error/Disconnected

## If Icon Doesn't Appear

1. **Check logs** - Look for "Menu bar icon ready!" message
2. **Try restarting** - Quit and relaunch
3. **Check permissions** - macOS may need permission
4. **Use --no-tray** - Agent works without icon too

## Files Created

- ✅ `build/ZeroTrace Agent.app` - App bundle with menu bar icon
- ✅ `dist/ZeroTrace-Agent-1.0.0-APP.dmg` - DMG with .app bundle
- ✅ `Info.plist` - macOS app configuration
- ✅ `entitlements.plist` - Security entitlements

## Summary

The menu bar icon is now working! It appears next to WiFi when you launch the .app bundle. 🎉


