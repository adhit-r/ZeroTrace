# Cross-Platform Tray Icon Implementation

## Overview

The ZeroTrace agent now includes a **production-ready, cross-platform system tray** implementation that works seamlessly on **Windows**, **Linux**, and **macOS**.

## Architecture

### Platform Interface Pattern

The implementation uses Go's build tag system to provide platform-specific functionality through a unified interface:

```
agent-go/internal/tray/
├── platform.go                    # Interface definition
├── platform_init_darwin.go        # macOS platform initialization
├── platform_init_windows.go       # Windows platform initialization
├── platform_init_linux.go         # Linux platform initialization
├── platform_init_fallback.go      # Fallback for unsupported platforms
├── platform_darwin.go             # macOS-specific implementations
├── platform_windows.go            # Windows-specific implementations
├── platform_linux.go              # Linux-specific implementations
├── platform_fallback.go           # Fallback implementations
├── icons.go                       # Icon management
├── assets/                        # Icon assets
│   ├── icon_green.png            # Connected state
│   ├── icon_gray.png             # Disconnected state
│   ├── icon_red.png              # Error state
│   └── icon_default.png          # Default state
├── tray.go                        # Full-featured tray manager
└── simple_tray.go                 # MDM-friendly simple tray
```

### Platform Operations Interface

```go
type PlatformOperations interface {
    ShowNotification(title, subtitle, message string) error
    OpenWebUI(url string) error
    CheckProcessRunning(processName string) bool
    RestartAgent(appPath string) error
    OpenSettings(configPath string) error
    KillProcess(processName string) error
    GetPlatformName() string
}
```

## Platform-Specific Features

### macOS (Darwin)
- **Notifications**: AppleScript-based notifications using `osascript`
- **Browser**: Opens URLs with `open` command
- **Process Management**: Uses `pgrep` and `pkill`
- **Settings**: Opens files with default text editor via `open -t`
- **Icon Format**: PNG (16x16 @1x, 32x32 @2x for Retina)

### Windows
- **Notifications**: PowerShell-based Windows Toast notifications
- **Browser**: Opens URLs with `cmd /c start`
- **Process Management**: Uses `tasklist` and `taskkill`
- **Settings**: Opens files with `notepad`
- **Icon Format**: ICO format recommended (16x16, 32x32, 48x48)

### Linux
- **Notifications**: `notify-send` with fallback to `zenity`
- **Browser**: `xdg-open` with fallback to common browsers (firefox, chromium, chrome)
- **Process Management**: Uses `pgrep` and `pkill`
- **Settings**: `xdg-open` with fallback to common editors (gedit, kate, nano, vim)
- **Icon Format**: PNG (22x22, 24x24 for most distros)

## Building for Different Platforms

### macOS (native)
```bash
cd agent-go
go build -o zerotrace-agent ./cmd/agent
```

### Windows (cross-compile from macOS/Linux)
```bash
GOOS=windows GOARCH=amd64 go build -o zerotrace-agent.exe ./cmd/agent
```

### Linux (cross-compile from macOS/Windows)
```bash
GOOS=linux GOARCH=amd64 go build -o zerotrace-agent ./cmd/agent
```

## Dependencies

### macOS
- No additional dependencies (uses built-in tools)
- `fyne.io/systray` (cross-platform tray library)

### Windows
- No additional dependencies (uses native Windows APIs via PowerShell)
- `fyne.io/systray` (cross-platform tray library)

### Linux Runtime Dependencies
- `libayatana-appindicator3-1` or `libappindicator3-1` (for tray icon)
- `notify-send` or `zenity` (for notifications, usually pre-installed)

**Ubuntu/Debian:**
```bash
sudo apt-get install libayatana-appindicator3-dev notify-osd
```

**Fedora/RHEL:**
```bash
sudo dnf install libayatana-appindicator-gtk3 libnotify
```

**Arch Linux:**
```bash
sudo pacman -S libappindicator-gtk3 libnotify
```

## Features

### Full Tray Manager (`tray.go`)
- Real-time system metrics (CPU, Memory)
- Agent status monitoring
- Manual vulnerability checks
- Web UI launcher
- Agent restart capability
- Settings editor
- Platform-specific notifications

### Simple Tray Manager (`simple_tray.go`)
- Minimal interface for MDM deployments
- Connection status indicator
- Basic CPU monitoring
- Platform-appropriate notifications
- Green/gray icon states (connected/disconnected)

## Icon Management

Icons are embedded at compile time using `go:embed` directives in `icons.go`:

```go
//go:embed assets/icon_green.png
var iconGreenData []byte
```

### Customizing Icons

Replace the PNG files in `agent-go/internal/tray/assets/` with your own icons:

- **icon_green.png**: Connected/healthy state
- **icon_gray.png**: Disconnected/offline state
- **icon_red.png**: Error/warning state
- **icon_default.png**: Default/startup state

**Recommended sizes:**
- macOS: 16x16 and 32x32 (retina)
- Windows: 16x16, 32x32, 48x48
- Linux: 22x22 or 24x24

## Testing

### Test on macOS
```bash
./zerotrace-agent
# Tray icon should appear in the menu bar
```

### Test on Windows
```cmd
zerotrace-agent.exe
# Tray icon should appear in the system tray
```

### Test on Linux
```bash
./zerotrace-agent
# Tray icon should appear in the system tray
# Requires a desktop environment with tray support
```

## Production Deployment

### macOS
1. Build the agent: `go build -o zerotrace-agent ./cmd/agent`
2. Package as `.app` bundle (optional)
3. Distribute via MDM or manual installation
4. Agent runs in menu bar with full tray functionality

### Windows
1. Build the agent: `GOOS=windows GOARCH=amd64 go build -o zerotrace-agent.exe ./cmd/agent`
2. Sign the executable (recommended for production)
3. Distribute via GPO, SCCM, or installer
4. Agent runs in system tray with native Windows notifications

### Linux
1. Build the agent: `GOOS=linux GOARCH=amd64 go build -o zerotrace-agent ./cmd/agent`
2. Create systemd service (optional for auto-start)
3. Install dependencies: `libayatana-appindicator3-1`
4. Distribute via package manager or manual installation
5. Agent runs in system tray with desktop notifications

## Troubleshooting

### macOS Issues
- **Tray icon not appearing**: Check macOS permissions for the application
- **Notifications not showing**: Enable notifications in System Preferences

### Windows Issues
- **Tray icon not appearing**: Check Windows system tray settings
- **Toast notifications not showing**: Enable notifications in Windows Settings
- **PowerShell execution error**: Check PowerShell execution policy

### Linux Issues
- **Tray icon not appearing**:
  - Install `libayatana-appindicator3-1` or `libappindicator3-1`
  - Ensure desktop environment supports system tray (GNOME requires Shell extension)
- **Notifications not showing**: Install `notify-send` or `zenity`
- **Browser not opening**: Install `xdg-utils`

## Future Enhancements

- [ ] Custom icon colors/themes
- [ ] More detailed system metrics in tray
- [ ] Context menu customization via config
- [ ] Multi-language notification support
- [ ] Animated icons for scan status
- [ ] Integration with desktop notification centers

## References

- [fyne.io/systray](https://github.com/fyne-io/systray) - Cross-platform tray library
- [Go Build Tags](https://pkg.go.dev/cmd/go#hdr-Build_constraints) - Platform-specific compilation
- Platform-specific documentation:
  - [macOS Notification Center](https://developer.apple.com/design/human-interface-guidelines/notification-center)
  - [Windows Toast Notifications](https://learn.microsoft.com/en-us/windows/apps/design/shell/tiles-and-notifications/toast-notifications-overview)
  - [Linux Desktop Notifications](https://specifications.freedesktop.org/notification-spec/latest/)
