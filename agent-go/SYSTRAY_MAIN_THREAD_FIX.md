# ✅ Systray Main Thread Fix

## Problem Solved

The systray was crashing with `SIGTRAP: trace trap` because `systray.Run()` was being called from a goroutine instead of the main thread on macOS.

## Root Cause

On macOS (Apple Silicon), `systray.Run()` **MUST** be called from the **main thread**. The Cocoa AppKit event loop requires this.

## Solution Implemented

### 1. Main Thread Execution

In `main.go`, for macOS:
```go
if !*disableTray && runtime.GOOS == "darwin" {
    // Lock OS thread for systray
    runtime.LockOSThread()
    
    // Create tray manager
    trayMgr := tray.NewSimpleTrayManager()
    
    // Define onReady callback
    onReady := func() {
        trayMgr.OnReady()  // Set up menu
        startAgentWork()   // Start background work
    }
    
    // Run systray on MAIN THREAD - blocks until systray.Quit()
    systray.Run(onReady, onExit)
}
```

### 2. Background Work in onReady

All agent background work (scanning, heartbeats) now starts **inside** the `onReady()` callback:
- Ensures systray is initialized first
- Runs in goroutines (which is fine)
- Main thread stays with systray

### 3. Exported Methods

Changed `onReady()` and `onExit()` to `OnReady()` and `OnExit()` so they can be called from `main()`.

## Key Changes

| Before | After |
|--------|-------|
| `systray.Run()` in goroutine | `systray.Run()` on main thread |
| Background work in `main()` | Background work in `onReady()` |
| Crashed with SIGTRAP | Works correctly |

## Testing

1. **Build:**
   ```bash
   cd agent-go
   GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -o mdm/build/ZeroTrace\ Agent.app/Contents/MacOS/zerotrace-agent cmd/agent/main.go
   ```

2. **Launch:**
   - Double-click `ZeroTrace Agent.app`
   - **No crash!**
   - Menu bar icon appears

3. **Check logs:**
   ```bash
   tail -f ~/.zerotrace/logs/agent.log
   ```

## Summary

✅ **Systray runs on main thread** (macOS requirement)
✅ **Background work in onReady()** (after systray init)
✅ **No more SIGTRAP crashes**
✅ **Menu bar icon works**

The systray now works correctly on macOS! 🎉


