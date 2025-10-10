// +build darwin

package tray

import (
	"fmt"
	"log"
	"os/exec"
)

// DarwinOperations implements platform-specific operations for macOS
type DarwinOperations struct{}

// ShowNotification displays a macOS notification using AppleScript
func (d *DarwinOperations) ShowNotification(title, subtitle, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s"`,
		message, title, subtitle)

	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to show notification: %v", err)
		return err
	}
	return nil
}

// OpenWebUI opens the web UI in the default browser
func (d *DarwinOperations) OpenWebUI(url string) error {
	cmd := exec.Command("open", url)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to open web UI: %v", err)
		return err
	}
	return nil
}

// CheckProcessRunning checks if the agent process is running
func (d *DarwinOperations) CheckProcessRunning(processName string) bool {
	cmd := exec.Command("pgrep", "-f", processName)
	err := cmd.Run()
	return err == nil
}

// RestartAgent restarts the agent process
func (d *DarwinOperations) RestartAgent(appPath string) error {
	// Kill existing agent
	if err := d.KillProcess("zerotrace-agent"); err != nil {
		log.Printf("Failed to kill existing agent: %v", err)
	}

	// Wait a moment for cleanup
	cmd := exec.Command("sleep", "2")
	cmd.Run()

	// Start new agent
	cmd = exec.Command("open", appPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to restart agent: %v", err)
		return err
	}

	return nil
}

// OpenSettings opens the settings file in the default text editor
func (d *DarwinOperations) OpenSettings(configPath string) error {
	cmd := exec.Command("open", "-t", configPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to open settings: %v", err)
		return err
	}
	return nil
}

// KillProcess kills the specified process
func (d *DarwinOperations) KillProcess(processName string) error {
	cmd := exec.Command("pkill", "-f", processName)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to kill process %s: %v", processName, err)
		return err
	}
	return nil
}

// GetPlatformName returns the platform name
func (d *DarwinOperations) GetPlatformName() string {
	return "macOS"
}
