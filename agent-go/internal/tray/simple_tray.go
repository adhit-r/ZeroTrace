package tray

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"zerotrace/agent/internal/monitor"

	"fyne.io/systray"
)

// SimpleTrayManager provides a minimal tray interface for MDM deployment
type SimpleTrayManager struct {
	monitor  *monitor.Monitor
	quitChan chan bool
}

// NewSimpleTrayManager creates a new simple tray manager
func NewSimpleTrayManager() *SimpleTrayManager {
	return &SimpleTrayManager{
		quitChan: make(chan bool, 1),
	}
}

// Start initializes and runs the simple tray icon
func (stm *SimpleTrayManager) Start() {
	if runtime.GOOS != "darwin" {
		log.Println("Tray icon only supported on macOS")
		return
	}

	// Start monitoring
	stm.monitor.Start()

	go func() {
		systray.Run(stm.onReady, stm.onExit)
	}()
}

// onReady is called when the tray is ready
func (stm *SimpleTrayManager) onReady() {
	// Set initial icon (gray - offline)
	systray.SetIcon(getGrayIcon())
	systray.SetTitle("ZeroTrace Agent")
	systray.SetTooltip("ZeroTrace Vulnerability Agent")

	// Create minimal menu
	mStatus := systray.AddMenuItem("🔄 Status: Checking...", "Agent status")
	mCPU := systray.AddMenuItem("📊 CPU: --", "CPU usage")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("❌ Quit", "Quit agent")

	// Start status monitoring
	go stm.monitorStatus(mStatus, mCPU)

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mStatus.ClickedCh:
				stm.showStatus()
			case <-mCPU.ClickedCh:
				stm.showCPUInfo()
			case <-mQuit.ClickedCh:
				stm.quitAgent()
			case <-stm.quitChan:
				return
			}
		}
	}()
}

// onExit is called when the tray is exiting
func (stm *SimpleTrayManager) onExit() {
	log.Println("Simple tray icon exiting")
}

// monitorStatus continuously updates the tray with status
func (stm *SimpleTrayManager) monitorStatus(mStatus, mCPU *systray.MenuItem) {
	ticker := time.NewTicker(10 * time.Second) // Less frequent updates for MDM
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check API connectivity
			apiConnected := stm.checkAPIConnectivity()

			// Update icon based on API status
			if apiConnected {
				systray.SetIcon(getGreenIcon())
				mStatus.SetTitle("🟢 Connected")
			} else {
				systray.SetIcon(getGrayIcon())
				mStatus.SetTitle("⚫ Disconnected")
			}

			// Update CPU usage
			metrics := stm.monitor.GetMetrics()
			mCPU.SetTitle(fmt.Sprintf("📊 CPU: %.1f%%", metrics.SystemCPU))

		case <-stm.quitChan:
			return
		}
	}
}

// checkAPIConnectivity checks if the API is reachable
func (stm *SimpleTrayManager) checkAPIConnectivity() bool {
	cmd := exec.Command("curl", "-s", "--connect-timeout", "3", "http://localhost:8080/health")
	return cmd.Run() == nil
}

// showStatus displays agent status notification
func (stm *SimpleTrayManager) showStatus() {
	apiConnected := stm.checkAPIConnectivity()
	status := "Disconnected"
	if apiConnected {
		status = "Connected"
	}

	stm.showNotification("ZeroTrace Agent", "Status", status)
}

// showCPUInfo displays CPU information
func (stm *SimpleTrayManager) showCPUInfo() {
	metrics := stm.monitor.GetMetrics()
	info := fmt.Sprintf("System CPU: %.1f%%", metrics.SystemCPU)
	stm.showNotification("ZeroTrace Agent", "CPU Usage", info)
}

// quitAgent quits the agent
func (stm *SimpleTrayManager) quitAgent() {
	stm.quitChan <- true
	systray.Quit()
}

// showNotification displays a macOS notification
func (stm *SimpleTrayManager) showNotification(title, subtitle, message string) {
	script := fmt.Sprintf(`
		display notification "%s" with title "%s" subtitle "%s"
	`, message, title, subtitle)

	exec.Command("osascript", "-e", script).Run()
}

// Stop stops the simple tray manager
func (stm *SimpleTrayManager) Stop() {
	stm.quitChan <- true
	stm.monitor.Stop()
}

// GetMonitor returns the monitor instance
func (stm *SimpleTrayManager) GetMonitor() *monitor.Monitor {
	return stm.monitor
}

// getGreenIcon returns the green ZeroTrace icon (API connected)
func getGreenIcon() []byte {
	// Green ZeroTrace logo icon data
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF, 0x61, 0x00, 0x00, 0x00,
		0x0C, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x60, 0x18, 0x05, 0x03,
		0x00, 0x00, 0x30, 0x00, 0x00, 0x01, 0x57, 0x6D, 0xB7, 0x4A, 0x00, 0x00,
		0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}

// getGrayIcon returns the gray ZeroTrace icon (API disconnected)
func getGrayIcon() []byte {
	// Gray ZeroTrace logo icon data
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF, 0x61, 0x00, 0x00, 0x00,
		0x0C, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x60, 0x18, 0x05, 0x03,
		0x00, 0x00, 0x30, 0x00, 0x00, 0x01, 0x57, 0x6D, 0xB7, 0x4A, 0x00, 0x00,
		0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}
