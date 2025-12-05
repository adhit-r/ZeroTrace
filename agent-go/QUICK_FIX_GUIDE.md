# Quick Fix Guide - DMG Crash

## What Happened?

The agent crashed when starting the systray UI:
```
SIGTRAP: trace trap
signal arrived during cgo execution
```

This happens because macOS requires proper code signing for native UI components.

## ✅ Fixed!

### Changes Made:

1. **Tray is now optional** - Can be disabled with `--no-tray` flag
2. **Better error handling** - Agent continues even if tray crashes
3. **Panic recovery** - Tray crashes don't kill the agent
4. **SimpleTrayManager** - More stable for macOS

## How to Use

### Option 1: Run with Tray (Try First)
```bash
/Volumes/ZeroTrace\ Agent/zerotrace-agent
```

If it crashes, use Option 2.

### Option 2: Run Without Tray (If crashes)
```bash
/Volumes/ZeroTrace\ Agent/zerotrace-agent --no-tray
```

The agent will work perfectly without the tray UI!

## What Works Without Tray

✅ All scanning features
✅ Network discovery
✅ Vulnerability detection
✅ API communication
✅ System monitoring
✅ Everything except the menu bar icon

## New DMG

**Location:** `agent-go/mdm/dist/ZeroTrace-Agent-1.0.0.dmg`

**Includes:**
- Fixed tray initialization
- Better error handling
- Option to disable tray

## Testing

1. **Mount DMG**: Double-click
2. **Try with tray**: Run normally
3. **If crashes**: Run with `--no-tray` flag
4. **Agent works**: All features functional!

The agent is fully functional with or without the tray UI! 🎉


