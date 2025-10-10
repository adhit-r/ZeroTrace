// +build linux

package tray

import (
	"fmt"
	"log"
	"os/exec"
)

// LinuxOperations implements platform-specific operations for Linux
type LinuxOperations struct{}

// ShowNotification displays a Linux notification using notify-send or zenity
func (l *LinuxOperations) ShowNotification(title, subtitle, message string) error {
	// Combine subtitle and message for Linux notification
	fullMessage := message
	if subtitle != "" {
		fullMessage = fmt.Sprintf("%s: %s", subtitle, message)
	}

	// Try notify-send first (most common)
	cmd := exec.Command("notify-send", title, fullMessage, "-i", "dialog-information", "-a", "ZeroTrace Agent")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to zenity if notify-send is not available
	cmd = exec.Command("zenity", "--notification", "--text", fmt.Sprintf("%s: %s", title, fullMessage))
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to show notification: %v", err)
		return err
	}

	return nil
}

// OpenWebUI opens the web UI in the default browser
func (l *LinuxOperations) OpenWebUI(url string) error {
	// Try xdg-open first (works on most Linux distros)
	cmd := exec.Command("xdg-open", url)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to specific browsers
	browsers := []string{"firefox", "chromium-browser", "chromium", "google-chrome", "chrome"}
	for _, browser := range browsers {
		cmd = exec.Command(browser, url)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no suitable browser found")
}

// CheckProcessRunning checks if the agent process is running
func (l *LinuxOperations) CheckProcessRunning(processName string) bool {
	cmd := exec.Command("pgrep", "-f", processName)
	err := cmd.Run()
	return err == nil
}

// RestartAgent restarts the agent process
func (l *LinuxOperations) RestartAgent(appPath string) error {
	// Kill existing agent
	if err := l.KillProcess("zerotrace-agent"); err != nil {
		log.Printf("Failed to kill existing agent: %v", err)
	}

	// Wait a moment for cleanup
	cmd := exec.Command("sleep", "2")
	cmd.Run()

	// Start new agent in background
	cmd = exec.Command(appPath)
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to restart agent: %v", err)
		return err
	}

	return nil
}

// OpenSettings opens the settings file in the default text editor
func (l *LinuxOperations) OpenSettings(configPath string) error {
	// Try xdg-open first (opens with default editor)
	cmd := exec.Command("xdg-open", configPath)
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Fallback to common editors
	editors := []string{"gedit", "kate", "nano", "vim", "vi"}
	for _, editor := range editors {
		cmd = exec.Command(editor, configPath)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no suitable text editor found")
}

// KillProcess kills the specified process
func (l *LinuxOperations) KillProcess(processName string) error {
	cmd := exec.Command("pkill", "-f", processName)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to kill process %s: %v", processName, err)
		return err
	}
	return nil
}

// GetPlatformName returns the platform name
func (l *LinuxOperations) GetPlatformName() string {
	return "Linux"
}
