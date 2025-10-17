package monitor

import (
	"log"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// Monitor tracks system and process metrics
type Monitor struct {
	agentPID   int32
	updateChan chan Metrics
	stopChan   chan bool
}

// Metrics contains system and process metrics
type Metrics struct {
	AgentCPU     float64
	AgentMemory  float64
	SystemCPU    float64
	SystemMemory float64
	Timestamp    time.Time
}

// NewMonitor creates a new monitor instance
func NewMonitor() *Monitor {
	return &Monitor{
		updateChan: make(chan Metrics, 1),
		stopChan:   make(chan bool, 1),
	}
}

// Start begins monitoring
func (m *Monitor) Start() {
	// Get current process PID
	m.agentPID = int32(runtime.NumCPU()) // Placeholder, will be updated

	go m.monitorLoop()
}

// Stop stops monitoring
func (m *Monitor) Stop() {
	m.stopChan <- true
}

// GetMetrics returns the latest metrics
func (m *Monitor) GetMetrics() Metrics {
	select {
	case metrics := <-m.updateChan:
		return metrics
	default:
		// Return default metrics if none available
		return Metrics{
			AgentCPU:     0.0,
			AgentMemory:  0.0,
			SystemCPU:    0.0,
			SystemMemory: 0.0,
			Timestamp:    time.Now(),
		}
	}
}

// monitorLoop continuously monitors system metrics
func (m *Monitor) monitorLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			metrics := m.collectMetrics()
			select {
			case m.updateChan <- metrics:
			default:
				// Channel full, skip this update
			}
		case <-m.stopChan:
			return
		}
	}
}

// collectMetrics gathers current system and process metrics
func (m *Monitor) collectMetrics() Metrics {
	metrics := Metrics{
		Timestamp: time.Now(),
	}

	// Get system CPU usage
	if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
		metrics.SystemCPU = cpuPercent[0]
	}

	// Get system memory usage
	if vmstat, err := mem.VirtualMemory(); err == nil {
		metrics.SystemMemory = vmstat.UsedPercent
	}

	// Get agent process metrics
	m.getAgentMetrics(&metrics)

	return metrics
}

// getAgentMetrics gets metrics for the agent process
func (m *Monitor) getAgentMetrics(metrics *Metrics) {
	// Find zerotrace-agent process
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error getting processes: %v", err)
		return
	}

	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}

		if name == "zerotrace-agent" {
			// Get CPU percentage
			if cpuPercent, err := proc.CPUPercent(); err == nil {
				metrics.AgentCPU = cpuPercent
			}

			// Get memory usage in MB
			if memInfo, err := proc.MemoryInfo(); err == nil {
				metrics.AgentMemory = float64(memInfo.RSS) / 1024 / 1024
			}

			break
		}
	}
}

// GetUpdateChannel returns the channel for receiving metric updates
func (m *Monitor) GetUpdateChannel() <-chan Metrics {
	return m.updateChan
}

// GetCPUUsage returns the current CPU usage percentage
func (m *Monitor) GetCPUUsage() float64 {
	metrics := m.GetMetrics()
	return metrics.SystemCPU
}

// GetMemoryUsage returns the current memory usage percentage
func (m *Monitor) GetMemoryUsage() float64 {
	metrics := m.GetMetrics()
	return metrics.SystemMemory
}
