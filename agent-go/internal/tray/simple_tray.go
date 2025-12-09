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
// NOTE: On macOS, this should NOT be called directly.
// Instead, systray.Run() must be called from main() on the main thread.
// This method is kept for non-macOS platforms or when called from main().
func (stm *SimpleTrayManager) Start() {
	log.Printf("Starting simple tray icon on %s (%s)", runtime.GOOS, stm.platform.GetPlatformName())

	// Start monitoring
	stm.monitor.Start()

	// On macOS, systray.Run() must be called from main thread in main()
	// For other platforms, we can run it in a goroutine
	if runtime.GOOS == "darwin" {
		log.Println("WARNING: On macOS, systray.Run() must be called from main() on main thread")
		log.Println("This Start() method should not be used on macOS")
		return
	}

	// Non-macOS: can run in goroutine
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Systray crashed: %v", r)
			}
		}()
		systray.Run(stm.OnReady, stm.OnExit)
	}()
}

// OnReady is called when the tray is ready - exported for use from main()
func (stm *SimpleTrayManager) OnReady() {
	log.Println(" Menu bar icon ready! Check top-right menu bar next to WiFi")

	// Start monitoring FIRST (required for CPU metrics)
	stm.monitor.Start()

	// Set initial icon (gray - offline)
	systray.SetIcon(GetGrayIcon())
	systray.SetTitle("ZeroTrace Agent")
	systray.SetTooltip("ZeroTrace Vulnerability Agent - Click for menu")

	// Create minimal menu
	mStatus := systray.AddMenuItem(" Status: Checking...", "Agent status")
	mCPU := systray.AddMenuItem(" CPU: --", "CPU usage")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem(" Quit", "Quit agent")

	// Start status monitoring (updates every 10 seconds)
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

// OnExit is called when the tray is exiting - exported for use from main()
func (stm *SimpleTrayManager) OnExit() {
	log.Println("Simple tray icon exiting")
}

// onExit is kept for backward compatibility (non-macOS)
func (stm *SimpleTrayManager) onExit() {
	stm.OnExit()
}

// monitorStatus continuously updates the tray with status
func (stm *SimpleTrayManager) monitorStatus(mStatus, mCPU *systray.MenuItem) {
	// Initial update after 2 seconds
	time.Sleep(2 * time.Second)

	ticker := time.NewTicker(10 * time.Second) // Update every 10 seconds
	defer ticker.Stop()

	// Initial update
	func() {
		// Check API connectivity
		apiConnected := stm.checkAPIConnectivity()

		// Update icon based on API status
		if apiConnected {
			systray.SetIcon(GetGreenIcon())
			mStatus.SetTitle(" Connected")
		} else {
			systray.SetIcon(GetGrayIcon())
			mStatus.SetTitle(" Disconnected")
		}

		// Update CPU usage
		metrics := stm.monitor.GetMetrics()
		if metrics.SystemCPU > 0 || metrics.AgentCPU > 0 {
			mCPU.SetTitle(fmt.Sprintf(" CPU: %.1f%% / %.1f%%", metrics.AgentCPU, metrics.SystemCPU))
		} else {
			mCPU.SetTitle(" CPU: Calculating...")
		}
	}()

	for {
		select {
		case <-ticker.C:
			// Check API connectivity
			apiConnected := stm.checkAPIConnectivity()

			// Update icon based on API status
			if apiConnected {
				systray.SetIcon(GetGreenIcon())
				mStatus.SetTitle(" Connected")
			} else {
				systray.SetIcon(GetGrayIcon())
				mStatus.SetTitle(" Disconnected")
			}

			// Update CPU usage
			metrics := stm.monitor.GetMetrics()
			if metrics.SystemCPU > 0 || metrics.AgentCPU > 0 {
				mCPU.SetTitle(fmt.Sprintf(" CPU: %.1f%% / %.1f%%", metrics.AgentCPU, metrics.SystemCPU))
			} else {
				mCPU.SetTitle(" CPU: Calculating...")
			}

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
	info := fmt.Sprintf("Agent CPU: %.1f%%\nSystem CPU: %.1f%%", metrics.AgentCPU, metrics.SystemCPU)
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
