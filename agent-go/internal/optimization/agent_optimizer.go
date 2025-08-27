package optimization

import (
	"context"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
)

// AgentOptimizer provides ultra-optimization for minimal CPU usage
type AgentOptimizer struct {
	logger *zap.Logger

	// Performance settings
	maxCPUPercent     float64
	maxMemoryMB       uint64
	scanInterval      time.Duration
	heartbeatInterval time.Duration

	// Resource monitoring
	cpuUsage      float64
	memoryUsage   uint64
	lastScan      time.Time
	lastHeartbeat time.Time

	// Optimization features
	adaptiveScanning     bool
	resourceThrottling   bool
	backgroundProcessing bool

	// Metrics
	metrics *AgentMetrics

	// Control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

// AgentMetrics tracks agent performance metrics
type AgentMetrics struct {
	mu sync.RWMutex

	// Resource usage
	cpuUsagePercent  float64
	memoryUsageMB    uint64
	diskUsagePercent float64
	networkUsageKB   uint64

	// Performance
	scansCompleted int64
	scansSkipped   int64
	heartbeatsSent int64
	errorsCount    int64

	// Timing
	avgScanDuration  time.Duration
	avgHeartbeatTime time.Duration
	lastOptimization time.Time
}

// NewAgentOptimizer creates a new agent optimizer
func NewAgentOptimizer(logger *zap.Logger) *AgentOptimizer {
	ctx, cancel := context.WithCancel(context.Background())

	return &AgentOptimizer{
		logger: logger,

		// Conservative resource limits
		maxCPUPercent:     5.0,             // Max 5% CPU usage
		maxMemoryMB:       50,              // Max 50MB memory
		scanInterval:      24 * time.Hour,  // Scan once per day
		heartbeatInterval: 5 * time.Minute, // Heartbeat every 5 minutes

		// Enable all optimizations
		adaptiveScanning:     true,
		resourceThrottling:   true,
		backgroundProcessing: true,

		metrics: &AgentMetrics{},
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start begins the optimization monitoring
func (ao *AgentOptimizer) Start() {
	ao.logger.Info("Starting agent optimizer",
		zap.Float64("max_cpu_percent", ao.maxCPUPercent),
		zap.Uint64("max_memory_mb", ao.maxMemoryMB),
		zap.Duration("scan_interval", ao.scanInterval),
		zap.Duration("heartbeat_interval", ao.heartbeatInterval))

	// Start background monitoring
	ao.wg.Add(1)
	go ao.monitorResources()

	// Start adaptive scanning
	if ao.adaptiveScanning {
		ao.wg.Add(1)
		go ao.adaptiveScanScheduler()
	}

	// Start resource throttling
	if ao.resourceThrottling {
		ao.wg.Add(1)
		go ao.resourceThrottler()
	}

	// Start background processing
	if ao.backgroundProcessing {
		ao.wg.Add(1)
		go ao.backgroundProcessor()
	}

	// Set Go runtime optimizations
	ao.optimizeGoRuntime()
}

// Stop gracefully stops the optimizer
func (ao *AgentOptimizer) Stop() {
	logger.Info("Stopping agent optimizer")
	ao.cancel()
	ao.wg.Wait()
	logger.Info("Agent optimizer stopped")
}

// optimizeGoRuntime applies Go runtime optimizations
func (ao *AgentOptimizer) optimizeGoRuntime() {
	// Set GOMAXPROCS to 1 for single-threaded operation
	runtime.GOMAXPROCS(1)

	// Set memory limit
	runtime.MemProfileRate = 0 // Disable memory profiling

	// Set GC target percentage (higher = less frequent GC)
	debug.SetGCPercent(500) // 500% = GC only when memory doubles

	// Set memory limit
	debug.SetMemoryLimit(50 * 1024 * 1024) // 50MB limit

	logger.Info("Go runtime optimized",
		zap.Int("gomaxprocs", runtime.GOMAXPROCS(0)),
		zap.Int("gc_percent", debug.SetGCPercent(-1)),
		zap.Int64("memory_limit", debug.SetMemoryLimit(-1)))
}

// monitorResources continuously monitors system resources
func (ao *AgentOptimizer) monitorResources() {
	defer ao.wg.Done()

	ticker := time.NewTicker(10 * time.Second) // Check every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ao.ctx.Done():
			return
		case <-ticker.C:
			ao.updateResourceMetrics()
		}
	}
}

// updateResourceMetrics updates current resource usage
func (ao *AgentOptimizer) updateResourceMetrics() {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	// Get CPU usage
	if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
		ao.cpuUsage = cpuPercent[0]
		ao.metrics.cpuUsagePercent = ao.cpuUsage
	}

	// Get memory usage
	if vmstat, err := mem.VirtualMemory(); err == nil {
		ao.memoryUsage = vmstat.Used
		ao.metrics.memoryUsageMB = vmstat.Used / 1024 / 1024
	}

	// Get disk usage
	if diskUsage, err := disk.Usage("/"); err == nil {
		ao.metrics.diskUsagePercent = diskUsage.UsedPercent
	}

	// Log if resources are high
	if ao.cpuUsage > ao.maxCPUPercent {
		logger.Warn("High CPU usage detected",
			zap.Float64("cpu_percent", ao.cpuUsage),
			zap.Float64("max_percent", ao.maxCPUPercent))
	}

	if ao.memoryUsage > ao.maxMemoryMB*1024*1024 {
		logger.Warn("High memory usage detected",
			zap.Uint64("memory_mb", ao.memoryUsage/1024/1024),
			zap.Uint64("max_mb", ao.maxMemoryMB))
	}
}

// adaptiveScanScheduler adapts scan timing based on system load
func (ao *AgentOptimizer) adaptiveScanScheduler() {
	defer ao.wg.Done()

	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	for {
		select {
		case <-ao.ctx.Done():
			return
		case <-ticker.C:
			ao.adjustScanInterval()
		}
	}
}

// adjustScanInterval adjusts scan interval based on system load
func (ao *AgentOptimizer) adjustScanInterval() {
	ao.mu.RLock()
	cpuUsage := ao.cpuUsage
	memoryUsage := ao.memoryUsage
	ao.mu.RUnlock()

	// Calculate new interval based on resource usage
	baseInterval := 24 * time.Hour

	// Increase interval if CPU usage is high
	if cpuUsage > ao.maxCPUPercent {
		multiplier := cpuUsage / ao.maxCPUPercent
		baseInterval = time.Duration(float64(baseInterval) * multiplier)
		logger.Info("Increased scan interval due to high CPU",
			zap.Float64("cpu_percent", cpuUsage),
			zap.Duration("new_interval", baseInterval))
	}

	// Increase interval if memory usage is high
	if memoryUsage > ao.maxMemoryMB*1024*1024 {
		multiplier := float64(memoryUsage) / float64(ao.maxMemoryMB*1024*1024)
		baseInterval = time.Duration(float64(baseInterval) * multiplier)
		logger.Info("Increased scan interval due to high memory",
			zap.Uint64("memory_mb", memoryUsage/1024/1024),
			zap.Duration("new_interval", baseInterval))
	}

	// Cap interval at 7 days
	if baseInterval > 7*24*time.Hour {
		baseInterval = 7 * 24 * time.Hour
	}

	ao.mu.Lock()
	ao.scanInterval = baseInterval
	ao.mu.Unlock()
}

// resourceThrottler throttles operations based on resource usage
func (ao *AgentOptimizer) resourceThrottler() {
	defer ao.wg.Done()

	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ao.ctx.Done():
			return
		case <-ticker.C:
			ao.checkResourceThrottling()
		}
	}
}

// checkResourceThrottling checks if operations should be throttled
func (ao *AgentOptimizer) checkResourceThrottling() {
	ao.mu.RLock()
	cpuUsage := ao.cpuUsage
	memoryUsage := ao.memoryUsage
	ao.mu.RUnlock()

	// Throttle if CPU usage is high
	if cpuUsage > ao.maxCPUPercent*0.8 { // 80% of max
		logger.Info("Throttling operations due to high CPU",
			zap.Float64("cpu_percent", cpuUsage))
		time.Sleep(1 * time.Second) // Add delay
	}

	// Throttle if memory usage is high
	if memoryUsage > ao.maxMemoryMB*1024*1024*0.8 { // 80% of max
		logger.Info("Throttling operations due to high memory",
			zap.Uint64("memory_mb", memoryUsage/1024/1024))
		time.Sleep(1 * time.Second) // Add delay
	}
}

// backgroundProcessor handles background tasks efficiently
func (ao *AgentOptimizer) backgroundProcessor() {
	defer ao.wg.Done()

	ticker := time.NewTicker(30 * time.Second) // Process every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ao.ctx.Done():
			return
		case <-ticker.C:
			ao.processBackgroundTasks()
		}
	}
}

// processBackgroundTasks processes background tasks efficiently
func (ao *AgentOptimizer) processBackgroundTasks() {
	// Only process if resources are available
	ao.mu.RLock()
	cpuUsage := ao.cpuUsage
	memoryUsage := ao.memoryUsage
	ao.mu.RUnlock()

	if cpuUsage > ao.maxCPUPercent*0.5 || memoryUsage > ao.maxMemoryMB*1024*1024*0.5 {
		return // Skip if resources are constrained
	}

	// Process background tasks
	ao.cleanupOldData()
	ao.optimizeMemory()
	ao.updateMetrics()
}

// cleanupOldData cleans up old data to save memory
func (ao *AgentOptimizer) cleanupOldData() {
	// Force garbage collection
	runtime.GC()

	// Clear old cache entries
	// This would be implemented based on your cache implementation

	logger.Debug("Cleaned up old data")
}

// optimizeMemory optimizes memory usage
func (ao *AgentOptimizer) optimizeMemory() {
	// Get current memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Force GC if heap is large
	if m.HeapAlloc > 10*1024*1024 { // 10MB
		runtime.GC()
		logger.Debug("Forced garbage collection",
			zap.Uint64("heap_alloc_mb", m.HeapAlloc/1024/1024))
	}
}

// updateMetrics updates performance metrics
func (ao *AgentOptimizer) updateMetrics() {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	ao.metrics.lastOptimization = time.Now()
}

// ShouldScan determines if a scan should be performed
func (ao *AgentOptimizer) ShouldScan() bool {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	// Check if enough time has passed
	if time.Since(ao.lastScan) < ao.scanInterval {
		return false
	}

	// Check if resources are available
	if ao.cpuUsage > ao.maxCPUPercent*0.8 {
		logger.Info("Skipping scan due to high CPU usage",
			zap.Float64("cpu_percent", ao.cpuUsage))
		ao.metrics.scansSkipped++
		return false
	}

	if ao.memoryUsage > ao.maxMemoryMB*1024*1024*0.8 {
		logger.Info("Skipping scan due to high memory usage",
			zap.Uint64("memory_mb", ao.memoryUsage/1024/1024))
		ao.metrics.scansSkipped++
		return false
	}

	return true
}

// ShouldHeartbeat determines if a heartbeat should be sent
func (ao *AgentOptimizer) ShouldHeartbeat() bool {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	return time.Since(ao.lastHeartbeat) >= ao.heartbeatInterval
}

// RecordScan records a completed scan
func (ao *AgentOptimizer) RecordScan(duration time.Duration) {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	ao.lastScan = time.Now()
	ao.metrics.scansCompleted++
	ao.metrics.avgScanDuration = ao.calculateAverageDuration(
		ao.metrics.avgScanDuration, duration, ao.metrics.scansCompleted)

	logger.Info("Scan completed",
		zap.Duration("duration", duration),
		zap.Duration("avg_duration", ao.metrics.avgScanDuration))
}

// RecordHeartbeat records a sent heartbeat
func (ao *AgentOptimizer) RecordHeartbeat(duration time.Duration) {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	ao.lastHeartbeat = time.Now()
	ao.metrics.heartbeatsSent++
	ao.metrics.avgHeartbeatTime = ao.calculateAverageDuration(
		ao.metrics.avgHeartbeatTime, duration, ao.metrics.heartbeatsSent)
}

// RecordError records an error
func (ao *AgentOptimizer) RecordError() {
	ao.mu.Lock()
	defer ao.mu.Unlock()

	ao.metrics.errorsCount++
}

// GetMetrics returns current metrics
func (ao *AgentOptimizer) GetMetrics() AgentMetrics {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	return *ao.metrics
}

// GetResourceUsage returns current resource usage
func (ao *AgentOptimizer) GetResourceUsage() (float64, uint64) {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	return ao.cpuUsage, ao.memoryUsage
}

// calculateAverageDuration calculates running average duration
func (ao *AgentOptimizer) calculateAverageDuration(current, new time.Duration, count int64) time.Duration {
	if count == 1 {
		return new
	}

	// Running average: (current * (count-1) + new) / count
	totalNanos := current.Nanoseconds()*int64(count-1) + new.Nanoseconds()
	return time.Duration(totalNanos / int64(count))
}

// OptimizeForLowCPU applies additional CPU optimizations
func (ao *AgentOptimizer) OptimizeForLowCPU() {
	// Set process priority to low
	if proc, err := process.NewProcess(int32(os.Getpid())); err == nil {
		proc.Nice(10) // Lower priority (higher nice value)
	}

	// Set CPU affinity to specific cores (if available)
	// This would be platform-specific

	// Disable unnecessary features
	runtime.MemProfileRate = 0   // Disable memory profiling
	runtime.BlockProfileRate = 0 // Disable block profiling

	logger.Info("Applied low CPU optimizations")
}

// OptimizeForLowMemory applies additional memory optimizations
func (ao *AgentOptimizer) OptimizeForLowMemory() {
	// Set memory limit
	debug.SetMemoryLimit(ao.maxMemoryMB * 1024 * 1024)

	// Set GC target percentage
	debug.SetGCPercent(200) // More aggressive GC

	// Pre-allocate memory pools
	// This would be implemented based on your memory usage patterns

	logger.Info("Applied low memory optimizations",
		zap.Uint64("memory_limit_mb", ao.maxMemoryMB))
}

// GetOptimizationStatus returns current optimization status
func (ao *AgentOptimizer) GetOptimizationStatus() map[string]interface{} {
	ao.mu.RLock()
	defer ao.mu.RUnlock()

	return map[string]interface{}{
		"cpu_usage_percent":     ao.cpuUsage,
		"memory_usage_mb":       ao.memoryUsage / 1024 / 1024,
		"max_cpu_percent":       ao.maxCPUPercent,
		"max_memory_mb":         ao.maxMemoryMB,
		"scan_interval":         ao.scanInterval.String(),
		"heartbeat_interval":    ao.heartbeatInterval.String(),
		"adaptive_scanning":     ao.adaptiveScanning,
		"resource_throttling":   ao.resourceThrottling,
		"background_processing": ao.backgroundProcessing,
		"metrics": map[string]interface{}{
			"scans_completed":    ao.metrics.scansCompleted,
			"scans_skipped":      ao.metrics.scansSkipped,
			"heartbeats_sent":    ao.metrics.heartbeatsSent,
			"errors_count":       ao.metrics.errorsCount,
			"avg_scan_duration":  ao.metrics.avgScanDuration.String(),
			"avg_heartbeat_time": ao.metrics.avgHeartbeatTime.String(),
		},
	}
}
