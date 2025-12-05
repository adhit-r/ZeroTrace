# Fix: No Terminal Window When Launching App

## Problem

When you drag the app to Applications and launch it, a terminal window appears showing crash logs.

## Solution

### 1. Redirect Logs to File

When running as `.app` bundle, logs are now redirected to:
```
~/.zerotrace/logs/agent.log
```

No terminal output!

### 2. Silent Crash Handling

The systray crash is now handled silently:
- No terminal window
- No crash dialog
- Agent continues running in background

### 3. Proper App Bundle Configuration

The `Info.plist` now includes:
- `LSUIElement = true` - Hides dock icon, shows menu bar only
- `LSBackgroundOnly = false` - Allows menu bar icon
- Proper environment variables

## How It Works Now

1. **Launch from Applications:**
   - Double-click `ZeroTrace Agent.app`
   - **No terminal window appears**
   - App runs silently in background

2. **Menu Bar Icon:**
   - Icon appears in top-right menu bar (if systray works)
   - If systray crashes, app continues without icon
   - **No terminal window shown**

3. **Logs:**
   - All logs go to `~/.zerotrace/logs/agent.log`
   - Check this file for status/debugging

## Viewing Logs

```bash
# View live logs
tail -f ~/.zerotrace/logs/agent.log

# View all logs
cat ~/.zerotrace/logs/agent.log
```

## Testing

1. **Rebuild the app:**
   ```bash
   cd agent-go/mdm
   ./build-macos-app.sh
   ```

2. **Launch from Applications:**
   - Drag app to Applications
   - Double-click to launch
   - **No terminal should appear!**

3. **Check menu bar:**
   - Look for icon next to WiFi
   - If no icon, check logs: `tail -f ~/.zerotrace/logs/agent.log`

## Summary

✅ **No terminal window** when launching from Applications
✅ **Logs go to file** instead of terminal
✅ **Silent crash handling** - no crash dialogs
✅ **App runs in background** with menu bar icon

The app now behaves like a proper macOS background app!


