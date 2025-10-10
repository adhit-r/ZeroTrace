// +build windows

package tray

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// WindowsOperations implements platform-specific operations for Windows
type WindowsOperations struct{}

// ShowNotification displays a Windows Toast notification using PowerShell
func (w *WindowsOperations) ShowNotification(title, subtitle, message string) error {
	// Combine subtitle and message for Windows toast
	fullMessage := message
	if subtitle != "" {
		fullMessage = fmt.Sprintf("%s: %s", subtitle, message)
	}

	// Use PowerShell to create a Windows Toast notification
	script := fmt.Sprintf(`
		[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null;
		[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null;
		$template = '<toast><visual><binding template="ToastText02"><text id="1">%s</text><text id="2">%s</text></binding></visual></toast>';
		$xml = New-Object Windows.Data.Xml.Dom.XmlDocument;
		$xml.LoadXml($template);
		$toast = [Windows.UI.Notifications.ToastNotification]::new($xml);
		[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier('ZeroTrace Agent').Show($toast);
	`, title, fullMessage)

	cmd := exec.Command("powershell", "-WindowStyle", "Hidden", "-Command", script)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to show notification: %v", err)
		return err
	}
	return nil
}

// OpenWebUI opens the web UI in the default browser
func (w *WindowsOperations) OpenWebUI(url string) error {
	cmd := exec.Command("cmd", "/c", "start", url)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to open web UI: %v", err)
		return err
	}
	return nil
}

// CheckProcessRunning checks if the agent process is running
func (w *WindowsOperations) CheckProcessRunning(processName string) bool {
	// Add .exe extension if not present
	if !strings.HasSuffix(processName, ".exe") {
		processName = processName + ".exe"
	}

	cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s", processName))
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), processName)
}

// RestartAgent restarts the agent process
func (w *WindowsOperations) RestartAgent(appPath string) error {
	// Kill existing agent
	if err := w.KillProcess("zerotrace-agent"); err != nil {
		log.Printf("Failed to kill existing agent: %v", err)
	}

	// Wait a moment for cleanup
	cmd := exec.Command("timeout", "/t", "2", "/nobreak")
	cmd.Run()

	// Start new agent
	cmd = exec.Command("cmd", "/c", "start", "", appPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to restart agent: %v", err)
		return err
	}

	return nil
}

// OpenSettings opens the settings file in the default text editor
func (w *WindowsOperations) OpenSettings(configPath string) error {
	cmd := exec.Command("notepad", configPath)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to open settings: %v", err)
		return err
	}
	return nil
}

// KillProcess kills the specified process
func (w *WindowsOperations) KillProcess(processName string) error {
	// Add .exe extension if not present
	if !strings.HasSuffix(processName, ".exe") {
		processName = processName + ".exe"
	}

	cmd := exec.Command("taskkill", "/F", "/IM", processName)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to kill process %s: %v", processName, err)
		return err
	}
	return nil
}

// GetPlatformName returns the platform name
func (w *WindowsOperations) GetPlatformName() string {
	return "Windows"
}
