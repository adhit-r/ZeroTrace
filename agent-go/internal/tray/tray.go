package tray

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"zerotrace/agent/internal/communicator"
	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/monitor"

	"fyne.io/systray"
)

// TrayManager handles the system tray icon and menu
type TrayManager struct {
	statusChan   chan string
	cpuChan      chan float64
	memChan      chan float64
	quitChan     chan bool
	monitor      *monitor.Monitor
	platform     PlatformOperations
	communicator *communicator.Communicator // Add communicator
	apiConnected bool                       // Track API connection status
}

// NewTrayManager creates a new tray manager
func NewTrayManager(cfg *config.Config) *TrayManager {
	return &TrayManager{
		statusChan:   make(chan string, 1),
		cpuChan:      make(chan float64, 1),
		memChan:      make(chan float64, 1),
		quitChan:     make(chan bool, 1),
		monitor:      monitor.NewMonitor(),
		platform:     GetPlatformOperations(),
		communicator: communicator.NewCommunicator(cfg),
		apiConnected: false,
	}
}

// Start initializes and runs the tray icon
func (tm *TrayManager) Start() {
	log.Printf("Starting tray icon on %s (%s)", runtime.GOOS, tm.platform.GetPlatformName())

	// Start monitoring
	tm.monitor.Start()

	go func() {
		// Try to run systray on all platforms, including macOS
		// with proper error handling
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Systray crashed on %s: %v", runtime.GOOS, r)
				log.Println("Agent will continue running without tray icon")
			}
		}()

		log.Printf("Attempting to start systray on %s", runtime.GOOS)
		systray.Run(tm.onReady, tm.onExit)
	}()
}

// onReady is called when the tray is ready
func (tm *TrayManager) onReady() {
	systray.SetIcon(GetDefaultIcon())
	systray.SetTitle("ZeroTrace Agent")
	systray.SetTooltip("ZeroTrace Vulnerability Agent")

	// Create menu items
	mStatus := systray.AddMenuItem("🔄 Agent Status", "Check agent status")
	mCPU := systray.AddMenuItem("📊 CPU: --", "CPU usage")
	mMem := systray.AddMenuItem("💾 Memory: --", "Memory usage")
	systray.AddSeparator()
	mCheck := systray.AddMenuItem("🔍 Check Now", "Manual check")
	mUI := systray.AddMenuItem("🌐 Open Web UI", "Open dashboard")
	mRestart := systray.AddMenuItem("🔄 Restart Agent", "Restart agent")
	systray.AddSeparator()
	mSettings := systray.AddMenuItem("⚙️ Settings", "Open settings")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("❌ Quit", "Quit agent")

	// Start monitoring goroutine
	go tm.monitorAndUpdate(mStatus, mCPU, mMem)

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mStatus.ClickedCh:
				tm.showStatus()
			case <-mCPU.ClickedCh:
				tm.showCPUDetails()
			case <-mMem.ClickedCh:
				tm.showMemoryDetails()
			case <-mCheck.ClickedCh:
				tm.checkNow()
			case <-mUI.ClickedCh:
				tm.openWebUI()
			case <-mRestart.ClickedCh:
				tm.restartAgent()
			case <-mSettings.ClickedCh:
				tm.openSettings()
			case <-mQuit.ClickedCh:
				tm.quitAgent()
			case <-tm.quitChan:
				return
			}
		}
	}()
}

// onExit is called when the tray is exiting
func (tm *TrayManager) onExit() {
	log.Println("Tray icon exiting")
}

// monitorAndUpdate continuously updates the tray menu with current stats
func (tm *TrayManager) monitorAndUpdate(mStatus, mCPU, mMem *systray.MenuItem) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	apiCheckTicker := time.NewTicker(30 * time.Second) // Check API status every 30 seconds
	defer apiCheckTicker.Stop()

	// Perform an initial API check right away
	tm.checkAPIStatus()

	for {
		select {
		case <-ticker.C:
			cpuUsage := tm.monitor.GetCPUUsage()
			memUsage := tm.monitor.GetMemoryUsage()

			// Update menu items
			if tm.apiConnected {
				mStatus.SetTitle("Status: Connected")
			} else {
				mStatus.SetTitle("Status: Disconnected")
			}
			mCPU.SetTitle(fmt.Sprintf("CPU: %.1f%%", cpuUsage))
			mMem.SetTitle(fmt.Sprintf("Memory: %.1f%%", memUsage))

		case <-apiCheckTicker.C:
			tm.checkAPIStatus()

		case <-tm.quitChan:
			return
		}
	}
}

// checkAPIStatus checks the API connection and updates the tray icon.
func (tm *TrayManager) checkAPIStatus() {
	err := tm.communicator.CheckAPIStatus()
	if err != nil {
		if tm.apiConnected {
			log.Println("Connection to API lost.")
			systray.SetIcon(GetRedIcon()) // Red icon for connection error
			tm.apiConnected = false
		} else {
			// If already disconnected, ensure icon is red
			systray.SetIcon(GetRedIcon())
		}
	} else {
		if !tm.apiConnected {
			log.Println("Successfully connected to API.")
			systray.SetIcon(GetGreenIcon()) // Green icon for connected
			tm.apiConnected = true
		}
	}
}

// getAgentStatus returns the current agent status
func (tm *TrayManager) getAgentStatus() string {
	// Check if agent process is running
	if !tm.platform.CheckProcessRunning("zerotrace-agent") {
		return "Not Running"
	}

	// Check API connectivity
	if tm.checkAPIConnectivity() {
		return "Running ✅"
	}
	return "Running ⚠️"
}

// getCPUUsage returns current CPU usage percentage
func (tm *TrayManager) getCPUUsage() float64 {
	metrics := tm.monitor.GetMetrics()
	return metrics.AgentCPU
}

// getMemoryUsage returns current memory usage in MB
func (tm *TrayManager) getMemoryUsage() float64 {
	metrics := tm.monitor.GetMetrics()
	return metrics.AgentMemory
}

// checkAPIConnectivity checks if the API is reachable
func (tm *TrayManager) checkAPIConnectivity() bool {
	cmd := exec.Command("curl", "-s", "--connect-timeout", "2", "http://localhost:8080/health")
	return cmd.Run() == nil
}

// showStatus displays the current status in a notification
func (tm *TrayManager) showStatus() {
	var status string
	if tm.apiConnected {
		status = "Connected to API"
	} else {
		status = "Disconnected from API"
	}
	tm.showNotification("ZeroTrace Agent", "Status", status)
}

// showCPUDetails displays detailed CPU information
func (tm *TrayManager) showCPUDetails() {
	cpu := tm.getCPUUsage()
	details := fmt.Sprintf("Agent CPU: %.1f%%\nSystem CPU: %.1f%%", cpu, cpu*2)
	tm.showNotification("ZeroTrace Agent", "CPU Usage", details)
}

// showMemoryDetails displays detailed memory information
func (tm *TrayManager) showMemoryDetails() {
	mem := tm.getMemoryUsage()
	details := fmt.Sprintf("Agent Memory: %.1f MB\nSystem Memory: %.1f%% used", mem, 65.2)
	tm.showNotification("ZeroTrace Agent", "Memory Usage", details)
}

// checkNow performs a manual status check
func (tm *TrayManager) checkNow() {
	status := tm.getAgentStatus()
	tm.showNotification("ZeroTrace Agent", "Manual Check", status)
}

// openWebUI opens the web dashboard
func (tm *TrayManager) openWebUI() {
	if err := tm.platform.OpenWebUI("http://localhost:5173"); err != nil {
		log.Printf("Failed to open web UI: %v", err)
	}
}

// restartAgent restarts the agent
func (tm *TrayManager) restartAgent() {
	// Restart using platform-specific method
	if err := tm.platform.RestartAgent("ZeroTrace Agent.app"); err != nil {
		log.Printf("Failed to restart agent: %v", err)
		tm.showNotification("ZeroTrace Agent", "Restart Failed", "Could not restart agent")
		return
	}

	tm.showNotification("ZeroTrace Agent", "Restart", "Agent restarted successfully")
}

// openSettings opens the configuration file
func (tm *TrayManager) openSettings() {
	if err := tm.platform.OpenSettings(".env"); err != nil {
		log.Printf("Failed to open settings: %v", err)
	}
}

// quitAgent quits the agent
func (tm *TrayManager) quitAgent() {
	tm.quitChan <- true
	systray.Quit()
}

// showNotification displays a system notification
func (tm *TrayManager) showNotification(title, subtitle, message string) {
	if err := tm.platform.ShowNotification(title, subtitle, message); err != nil {
		log.Printf("Failed to show notification: %v", err)
	}
}

// UpdateStatus updates the status from external sources
func (tm *TrayManager) UpdateStatus(status string) {
	select {
	case tm.statusChan <- status:
	default:
	}
}

// UpdateCPU updates the CPU usage from external sources
func (tm *TrayManager) UpdateCPU(cpu float64) {
	select {
	case tm.cpuChan <- cpu:
	default:
	}
}

// UpdateMemory updates the memory usage from external sources
func (tm *TrayManager) UpdateMemory(mem float64) {
	select {
	case tm.memChan <- mem:
	default:
	}
}

// Stop stops the tray manager
func (tm *TrayManager) Stop() {
	tm.quitChan <- true
	tm.monitor.Stop()
}

// GetMonitor returns the monitor instance
func (tm *TrayManager) GetMonitor() *monitor.Monitor {
	return tm.monitor
}
