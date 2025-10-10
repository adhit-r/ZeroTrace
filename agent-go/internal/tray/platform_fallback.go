// +build !windows,!linux,!darwin

package tray

import (
	"fmt"
	"log"
)

// FallbackOperations implements platform-specific operations for unsupported platforms
// This provides basic logging functionality when running on unsupported platforms
type FallbackOperations struct{}

// ShowNotification logs the notification message
func (f *FallbackOperations) ShowNotification(title, subtitle, message string) error {
	log.Printf("Notification: %s - %s: %s", title, subtitle, message)
	return fmt.Errorf("notifications not supported on this platform")
}

// OpenWebUI logs the URL that would be opened
func (f *FallbackOperations) OpenWebUI(url string) error {
	log.Printf("Open URL: %s", url)
	return fmt.Errorf("opening browser not supported on this platform")
}

// CheckProcessRunning always returns false on unsupported platforms
func (f *FallbackOperations) CheckProcessRunning(processName string) bool {
	log.Printf("Process check not supported for: %s", processName)
	return false
}

// RestartAgent returns an error on unsupported platforms
func (f *FallbackOperations) RestartAgent(appPath string) error {
	return fmt.Errorf("restart not supported on this platform")
}

// OpenSettings returns an error on unsupported platforms
func (f *FallbackOperations) OpenSettings(configPath string) error {
	log.Printf("Cannot open settings: %s", configPath)
	return fmt.Errorf("opening settings not supported on this platform")
}

// KillProcess returns an error on unsupported platforms
func (f *FallbackOperations) KillProcess(processName string) error {
	return fmt.Errorf("kill process not supported on this platform")
}

// GetPlatformName returns "Unsupported"
func (f *FallbackOperations) GetPlatformName() string {
	return "Unsupported Platform"
}
