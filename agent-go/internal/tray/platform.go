package tray

// PlatformOperations defines the interface for platform-specific operations
type PlatformOperations interface {
	// ShowNotification displays a system notification
	ShowNotification(title, subtitle, message string) error

	// OpenWebUI opens the web UI in the default browser
	OpenWebUI(url string) error

	// CheckProcessRunning checks if the agent process is running
	CheckProcessRunning(processName string) bool

	// RestartAgent restarts the agent process
	RestartAgent(appPath string) error

	// OpenSettings opens the settings file in default editor
	OpenSettings(configPath string) error

	// KillProcess kills the agent process
	KillProcess(processName string) error

	// GetPlatformName returns the platform name
	GetPlatformName() string
}

// GetPlatformOperations returns the appropriate platform-specific implementation
// This function is implemented in platform-specific files:
// - platform_init_darwin.go for macOS
// - platform_init_windows.go for Windows
// - platform_init_linux.go for Linux
// - platform_init_fallback.go for other platforms
