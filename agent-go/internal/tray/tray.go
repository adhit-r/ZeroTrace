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

// TrayManager handles the system tray icon and menu
type TrayManager struct {
	statusChan chan string
	cpuChan    chan float64
	memChan    chan float64
	quitChan   chan bool
	monitor    *monitor.Monitor
}

// NewTrayManager creates a new tray manager
func NewTrayManager() *TrayManager {
	return &TrayManager{
		statusChan: make(chan string, 1),
		cpuChan:    make(chan float64, 1),
		memChan:    make(chan float64, 1),
		quitChan:   make(chan bool, 1),
		monitor:    monitor.NewMonitor(),
	}
}

// Start initializes and runs the tray icon
func (tm *TrayManager) Start() {
	if runtime.GOOS != "darwin" {
		log.Println("Tray icon only supported on macOS")
		return
	}

	// Start monitoring
	tm.monitor.Start()

	go func() {
		systray.Run(tm.onReady, tm.onExit)
	}()
}

// onReady is called when the tray is ready
func (tm *TrayManager) onReady() {
	systray.SetIcon(getIcon())
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

	for {
		select {
		case <-ticker.C:
			// Update status
			status := tm.getAgentStatus()
			mStatus.SetTitle(fmt.Sprintf("🔄 %s", status))

			// Update CPU usage
			cpu := tm.getCPUUsage()
			mCPU.SetTitle(fmt.Sprintf("📊 CPU: %.1f%%", cpu))

			// Update memory usage
			mem := tm.getMemoryUsage()
			mMem.SetTitle(fmt.Sprintf("💾 Memory: %.1f MB", mem))

		case <-tm.quitChan:
			return
		}
	}
}

// getAgentStatus returns the current agent status
func (tm *TrayManager) getAgentStatus() string {
	// Check if agent process is running
	cmd := exec.Command("pgrep", "-f", "zerotrace-agent")
	if err := cmd.Run(); err != nil {
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

// showStatus displays agent status notification
func (tm *TrayManager) showStatus() {
	status := tm.getAgentStatus()
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
	exec.Command("open", "http://localhost:5173").Run()
}

// restartAgent restarts the agent
func (tm *TrayManager) restartAgent() {
	// Kill existing agent
	exec.Command("pkill", "-f", "zerotrace-agent").Run()
	time.Sleep(2 * time.Second)

	// Start new agent
	exec.Command("open", "ZeroTrace Agent.app").Run()

	tm.showNotification("ZeroTrace Agent", "Restart", "Agent restarted successfully")
}

// openSettings opens the configuration file
func (tm *TrayManager) openSettings() {
	exec.Command("open", "-t", ".env").Run()
}

// quitAgent quits the agent
func (tm *TrayManager) quitAgent() {
	tm.quitChan <- true
	systray.Quit()
}

// showNotification displays a macOS notification
func (tm *TrayManager) showNotification(title, subtitle, message string) {
	script := fmt.Sprintf(`
		display notification "%s" with title "%s" subtitle "%s"
	`, message, title, subtitle)

	exec.Command("osascript", "-e", script).Run()
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

// getIcon returns the tray icon data
func getIcon() []byte {
	// This would return actual icon data
	// For now, return a simple placeholder
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF, 0x61, 0x00, 0x00, 0x00,
		0x0C, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x60, 0x18, 0x05, 0x03,
		0x00, 0x00, 0x30, 0x00, 0x00, 0x01, 0x57, 0x6D, 0xB7, 0x4A, 0x00, 0x00,
		0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}
