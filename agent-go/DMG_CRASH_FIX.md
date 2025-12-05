# DMG Crash Fix

## Problem

The agent was crashing when running from DMG with:
```
SIGTRAP: trace trap
signal arrived during cgo execution
fyne.io/systray._Cfunc_nativeLoop()
```

## Root Cause

The systray library requires:
- CGO enabled
- Native macOS frameworks
- Proper code signing/entitlements

When running from DMG without proper signing, macOS blocks native calls.

## Solution

### 1. Made Tray Optional

Added `--no-tray` flag to disable tray UI:

```bash
./zerotrace-agent --no-tray
```

### 2. Better Error Handling

Tray now starts in a goroutine with panic recovery:
- If tray crashes, agent continues running
- Logs error but doesn't crash agent

### 3. Use SimpleTrayManager

Switched to `SimpleTrayManager` for macOS (more stable than full tray).

## Usage

### With Tray (Default)
```bash
./zerotrace-agent
```

### Without Tray (If crashes)
```bash
./zerotrace-agent --no-tray
```

## New DMG

The updated DMG includes:
- Fixed tray initialization
- Better error handling
- Option to disable tray

**Location:** `agent-go/mdm/dist/ZeroTrace-Agent-1.0.0.dmg`

## Testing

1. **Mount DMG**: Double-click the DMG
2. **Run agent**: `/Volumes/ZeroTrace\ Agent/zerotrace-agent`
3. **If crashes**: Try `./zerotrace-agent --no-tray`

The agent will work with or without the tray UI!


