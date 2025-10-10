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
	platform PlatformOperations
}

// NewSimpleTrayManager creates a new simple tray manager
func NewSimpleTrayManager() *SimpleTrayManager {
	return &SimpleTrayManager{
		quitChan: make(chan bool, 1),
		monitor:  monitor.NewMonitor(),
		platform: GetPlatformOperations(),
	}
}

// Start initializes and runs the simple tray icon
func (stm *SimpleTrayManager) Start() {
	log.Printf("Starting simple tray icon on %s (%s)", runtime.GOOS, stm.platform.GetPlatformName())

	// Start monitoring
	stm.monitor.Start()

	go func() {
		if runtime.GOOS == "darwin" {
			log.Println("Systray temporarily disabled on macOS due to native crash. Agent running without tray icon.")
			// For macOS, we'll just keep the goroutine alive without systray.Run()
			// In a real scenario, we'd use a different tray library or fix the systray issue.
			select {}
		} else {
			systray.Run(stm.onReady, stm.onExit)
		}
	}()
}

// onReady is called when the tray is ready
func (stm *SimpleTrayManager) onReady() {
	// Set initial icon (gray - offline)
	systray.SetIcon(GetGrayIcon())
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
				systray.SetIcon(GetGreenIcon())
				mStatus.SetTitle("🟢 Connected")
			} else {
				systray.SetIcon(GetGrayIcon())
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

// showNotification displays a system notification
func (stm *SimpleTrayManager) showNotification(title, subtitle, message string) {
	err := stm.platform.ShowNotification(title, subtitle, message)
	if err != nil {
		log.Printf("Failed to show notification: %v", err)
	}
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
