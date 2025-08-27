package monitoring

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// APM provides comprehensive application performance monitoring
type APM struct {
	logger *zap.Logger

	// HTTP metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpRequestSize     *prometheus.HistogramVec
	httpResponseSize    *prometheus.HistogramVec

	// Business metrics
	scansTotal           *prometheus.CounterVec
	vulnerabilitiesFound *prometheus.CounterVec
	agentsOnline         *prometheus.GaugeVec
	agentHeartbeats      *prometheus.CounterVec

	// System metrics
	goroutines  prometheus.GaugeFunc
	memoryAlloc prometheus.GaugeFunc
	memoryHeap  prometheus.GaugeFunc
	gcDuration  prometheus.HistogramFunc

	// Database metrics
	dbConnections   *prometheus.GaugeVec
	dbQueryDuration *prometheus.HistogramVec
	dbQueryErrors   *prometheus.CounterVec

	// Queue metrics
	queueSize           *prometheus.GaugeVec
	queueProcessingTime *prometheus.HistogramVec
	queueErrors         *prometheus.CounterVec

	// Cache metrics
	cacheHits   *prometheus.CounterVec
	cacheMisses *prometheus.CounterVec
	cacheSize   *prometheus.GaugeVec

	// Enrichment metrics
	enrichmentDuration  *prometheus.HistogramVec
	enrichmentErrors    *prometheus.CounterVec
	enrichmentCacheHits *prometheus.CounterVec
}

// NewAPM creates a new APM instance
func NewAPM(logger *zap.Logger) *APM {
	apm := &APM{
		logger: logger,
	}

	apm.initMetrics()
	return apm
}

// initMetrics initializes all Prometheus metrics
func (a *APM) initMetrics() {
	// HTTP metrics
	a.httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status", "company_id"},
	)

	a.httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "company_id"},
	)

	a.httpRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "HTTP request size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	a.httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8),
		},
		[]string{"method", "endpoint"},
	)

	// Business metrics
	a.scansTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "scans_total",
			Help: "Total number of scans",
		},
		[]string{"status", "company_id"},
	)

	a.vulnerabilitiesFound = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "vulnerabilities_found_total",
			Help: "Total number of vulnerabilities found",
		},
		[]string{"severity", "company_id"},
	)

	a.agentsOnline = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "agents_online",
			Help: "Number of online agents",
		},
		[]string{"company_id"},
	)

	a.agentHeartbeats = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "agent_heartbeats_total",
			Help: "Total number of agent heartbeats",
		},
		[]string{"agent_id", "company_id"},
	)

	// System metrics
	a.goroutines = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "goroutines_total",
			Help: "Number of goroutines",
		},
		func() float64 {
			return float64(runtime.NumGoroutine())
		},
	)

	a.memoryAlloc = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "memory_alloc_bytes",
			Help: "Allocated memory in bytes",
		},
		func() float64 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return float64(m.Alloc)
		},
	)

	a.memoryHeap = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "memory_heap_bytes",
			Help: "Heap memory in bytes",
		},
		func() float64 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			return float64(m.HeapAlloc)
		},
	)

	// Database metrics
	a.dbConnections = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "db_connections",
			Help: "Database connections",
		},
		[]string{"status"},
	)

	a.dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	a.dbQueryErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_query_errors_total",
			Help: "Database query errors",
		},
		[]string{"operation", "table"},
	)

	// Queue metrics
	a.queueSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "queue_size",
			Help: "Queue size",
		},
		[]string{"queue_name", "company_id"},
	)

	a.queueProcessingTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "queue_processing_duration_seconds",
			Help:    "Queue processing duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"queue_name", "company_id"},
	)

	a.queueErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "queue_errors_total",
			Help: "Queue processing errors",
		},
		[]string{"queue_name", "company_id"},
	)

	// Cache metrics
	a.cacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Cache hits",
		},
		[]string{"cache_name"},
	)

	a.cacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Cache misses",
		},
		[]string{"cache_name"},
	)

	a.cacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_size",
			Help: "Cache size",
		},
		[]string{"cache_name"},
	)

	// Enrichment metrics
	a.enrichmentDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "enrichment_duration_seconds",
			Help:    "CVE enrichment duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"source", "company_id"},
	)

	a.enrichmentErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "enrichment_errors_total",
			Help: "CVE enrichment errors",
		},
		[]string{"source", "company_id"},
	)

	a.enrichmentCacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "enrichment_cache_hits_total",
			Help: "CVE enrichment cache hits",
		},
		[]string{"source", "company_id"},
	)
}

// HTTPMiddleware provides HTTP monitoring middleware
func (a *APM) HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Get request size
		reqSize := c.Request.ContentLength
		if reqSize < 0 {
			reqSize = 0
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get company ID from context
		companyID := c.GetString("company_id")
		if companyID == "" {
			companyID = "unknown"
		}

		// Record metrics
		a.httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			fmt.Sprintf("%d", c.Writer.Status()),
			companyID,
		).Inc()

		a.httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
			companyID,
		).Observe(duration)

		a.httpRequestSize.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(float64(reqSize))

		a.httpResponseSize.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(float64(c.Writer.Size()))

		// Log slow requests
		if duration > 1.0 {
			a.logger.Warn("Slow HTTP request",
				zap.String("method", c.Request.Method),
				zap.String("path", c.FullPath()),
				zap.Duration("duration", time.Duration(duration*float64(time.Second))),
				zap.String("company_id", companyID),
			)
		}
	}
}

// MetricsHandler returns Prometheus metrics handler
func (a *APM) MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// RecordScan records scan metrics
func (a *APM) RecordScan(status, companyID string) {
	a.scansTotal.WithLabelValues(status, companyID).Inc()
}

// RecordVulnerability records vulnerability metrics
func (a *APM) RecordVulnerability(severity, companyID string) {
	a.vulnerabilitiesFound.WithLabelValues(severity, companyID).Inc()
}

// UpdateAgentStatus updates agent status metrics
func (a *APM) UpdateAgentStatus(companyID string, count int) {
	a.agentsOnline.WithLabelValues(companyID).Set(float64(count))
}

// RecordAgentHeartbeat records agent heartbeat
func (a *APM) RecordAgentHeartbeat(agentID, companyID string) {
	a.agentHeartbeats.WithLabelValues(agentID, companyID).Inc()
}

// RecordDBQuery records database query metrics
func (a *APM) RecordDBQuery(operation, table string, duration time.Duration, err error) {
	a.dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())

	if err != nil {
		a.dbQueryErrors.WithLabelValues(operation, table).Inc()
	}
}

// UpdateQueueSize updates queue size metrics
func (a *APM) UpdateQueueSize(queueName, companyID string, size int) {
	a.queueSize.WithLabelValues(queueName, companyID).Set(float64(size))
}

// RecordQueueProcessing records queue processing metrics
func (a *APM) RecordQueueProcessing(queueName, companyID string, duration time.Duration, err error) {
	a.queueProcessingTime.WithLabelValues(queueName, companyID).Observe(duration.Seconds())

	if err != nil {
		a.queueErrors.WithLabelValues(queueName, companyID).Inc()
	}
}

// RecordCacheHit records cache hit
func (a *APM) RecordCacheHit(cacheName string) {
	a.cacheHits.WithLabelValues(cacheName).Inc()
}

// RecordCacheMiss records cache miss
func (a *APM) RecordCacheMiss(cacheName string) {
	a.cacheMisses.WithLabelValues(cacheName).Inc()
}

// UpdateCacheSize updates cache size
func (a *APM) UpdateCacheSize(cacheName string, size int) {
	a.cacheSize.WithLabelValues(cacheName).Set(float64(size))
}

// RecordEnrichment records enrichment metrics
func (a *APM) RecordEnrichment(source, companyID string, duration time.Duration, err error) {
	a.enrichmentDuration.WithLabelValues(source, companyID).Observe(duration.Seconds())

	if err != nil {
		a.enrichmentErrors.WithLabelValues(source, companyID).Inc()
	}
}

// RecordEnrichmentCacheHit records enrichment cache hit
func (a *APM) RecordEnrichmentCacheHit(source, companyID string) {
	a.enrichmentCacheHits.WithLabelValues(source, companyID).Inc()
}

// StartMetricsServer starts the metrics server
func (a *APM) StartMetricsServer(addr string) {
	http.Handle("/metrics", a.MetricsHandler())

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			a.logger.Error("Failed to start metrics server", zap.Error(err))
		}
	}()

	a.logger.Info("Metrics server started", zap.String("addr", addr))
}

// GetMetrics returns current metrics summary
func (a *APM) GetMetrics() map[string]interface{} {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return map[string]interface{}{
		"goroutines": runtime.NumGoroutine(),
		"memory": map[string]interface{}{
			"alloc":      m.Alloc,
			"heap_alloc": m.HeapAlloc,
			"heap_sys":   m.HeapSys,
			"heap_idle":  m.HeapIdle,
			"heap_inuse": m.HeapInuse,
		},
		"gc": map[string]interface{}{
			"num_gc":      m.NumGC,
			"pause_total": m.PauseTotalNs,
		},
	}
}
