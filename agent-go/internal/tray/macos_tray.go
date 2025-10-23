package tray

import (
	"log"
	"runtime"
	"time"
	"zerotrace/agent/internal/monitor"
)

// MacOSTrayManager provides a macOS-specific tray implementation
// that uses native macOS APIs instead of the problematic systray library
type MacOSTrayManager struct {
	monitor  *monitor.Monitor
	quitChan chan bool
	platform PlatformOperations
}

// NewMacOSTrayManager creates a new macOS-specific tray manager
func NewMacOSTrayManager() *MacOSTrayManager {
	return &MacOSTrayManager{
		quitChan: make(chan bool, 1),
		monitor:  monitor.NewMonitor(),
		platform: GetPlatformOperations(),
	}
}

// Start initializes the macOS tray (alternative implementation)
func (mtm *MacOSTrayManager) Start() {
	log.Printf("Starting macOS tray manager on %s", runtime.GOOS)

	// Start monitoring
	mtm.monitor.Start()

	// For now, we'll use a simple approach that shows the agent is running
	// without the problematic systray library
	go mtm.macOSStatusLoop()
}

// macOSStatusLoop provides status updates without systray
func (mtm *MacOSTrayManager) macOSStatusLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Println("ZeroTrace Agent is running in the background")
	log.Println("Status: Active - Monitoring system for vulnerabilities")
	log.Println("Note: Tray icon disabled on macOS due to library compatibility issues")
	log.Println("Use 'ps aux | grep zerotrace' to check if agent is running")

	for {
		select {
		case <-ticker.C:
			// Log status updates
			log.Printf("Agent status: Running - CPU: %.1f%%, Memory: %.1f%%",
				mtm.monitor.GetCPUUsage(), mtm.monitor.GetMemoryUsage())
		case <-mtm.quitChan:
			log.Println("macOS tray manager shutting down")
			return
		}
	}
}

// Stop stops the macOS tray manager
func (mtm *MacOSTrayManager) Stop() {
	select {
	case mtm.quitChan <- true:
	default:
	}
}

