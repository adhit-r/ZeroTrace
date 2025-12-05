package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/google/uuid"
)

// AppData represents application data from agents
type AppData struct {
	ID           string    `json:"id"`
	CompanyID    string    `json:"company_id"`
	AgentID      string    `json:"agent_id"`
	AppName      string    `json:"app_name"`
	AppVersion   string    `json:"app_version"`
	PackageType  string    `json:"package_type"`
	Architecture string    `json:"architecture"`
	Timestamp    time.Time `json:"timestamp"`
}

// ProcessedResult represents enriched app data
type ProcessedResult struct {
	AppID           string          `json:"app_id"`
	CompanyID       string          `json:"company_id"`
	AgentID         string          `json:"agent_id"`
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	EnrichedAt      time.Time       `json:"enriched_at"`
	Error           string          `json:"error,omitempty"`
}

// Vulnerability represents CVE data
type Vulnerability struct {
	CVEID       string  `json:"cve_id"`
	Severity    string  `json:"severity"`
	CVSSScore   float64 `json:"cvss_score"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
}

// QueueProcessor handles batch processing of app data
type QueueProcessor struct {
	redis     *redis.Client
	batchSize int
	workers   int
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	metrics   *QueueMetrics
}

// QueueMetrics tracks processing metrics
type QueueMetrics struct {
	mu             sync.RWMutex
	processedTotal int64
	processedToday int64
	errorsTotal    int64
	queueSize      int64
}

// NewQueueProcessor creates a new queue processor
func NewQueueProcessor(redis *redis.Client) *QueueProcessor {
	ctx, cancel := context.WithCancel(context.Background())

	return &QueueProcessor{
		redis:     redis,
		batchSize: 100, // Process 100 apps at a time
		workers:   10,  // 10 concurrent workers
		ctx:       ctx,
		cancel:    cancel,
		metrics:   &QueueMetrics{},
	}
}

// Start begins processing
func (qp *QueueProcessor) Start() {
	log.Printf("Starting queue processor with %d workers, batch size %d", qp.workers, qp.batchSize)

	// Start workers
	for i := 0; i < qp.workers; i++ {
		qp.wg.Add(1)
		go qp.worker(i)
	}

	// Start metrics collector
	go qp.metricsCollector()

	// Start cleanup routine
	go qp.cleanupRoutine()
}

// Stop gracefully stops processing
func (qp *QueueProcessor) Stop() {
	log.Println("Stopping queue processor...")
	qp.cancel()
	qp.wg.Wait()
	log.Println("Queue processor stopped")
}

// AddApp adds an app to the processing queue using Valkey Streams
func (qp *QueueProcessor) AddApp(app AppData) error {
	// Generate ID if not provided
	if app.ID == "" {
		app.ID = uuid.New().String()
	}

	// Set timestamp
	if app.Timestamp.IsZero() {
		app.Timestamp = time.Now()
	}

	// Serialize app data
	data, err := json.Marshal(app)
	if err != nil {
		return fmt.Errorf("failed to marshal app data: %w", err)
	}

	// Add to Valkey Stream (better than sorted sets for queues)
	// Stream name: app_queue
	// Fields: data (JSON), company_id, priority (timestamp for ordering)
	priority := app.Timestamp.Unix()
	err = qp.redis.XAdd(qp.ctx, &redis.XAddArgs{
		Stream: "app_queue",
		Values: map[string]interface{}{
			"data":       string(data),
			"company_id": app.CompanyID,
			"agent_id":   app.AgentID,
			"priority":   priority,
		},
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to add app to queue: %w", err)
	}

	// Update metrics
	qp.metrics.incrementQueueSize()

	return nil
}

// worker processes apps from the queue
func (qp *QueueProcessor) worker(id int) {
	defer qp.wg.Done()

	log.Printf("Worker %d started", id)

	for {
		select {
		case <-qp.ctx.Done():
			log.Printf("Worker %d stopping", id)
			return
		default:
			// Get batch of apps from queue
			apps := qp.getBatchFromQueue()

			if len(apps) > 0 {
				log.Printf("Worker %d processing batch of %d apps", id, len(apps))
				qp.processBatch(apps)
			} else {
				// No apps in queue, wait a bit
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

// getBatchFromQueue retrieves a batch of apps from the queue using Valkey Streams
func (qp *QueueProcessor) getBatchFromQueue() []AppData {
	var apps []AppData

	// Read from stream using XREADGROUP for consumer groups (better for distributed processing)
	// For now, use simple XREAD with COUNT
	streams, err := qp.redis.XRead(qp.ctx, &redis.XReadArgs{
		Streams: []string{"app_queue", "0"}, // Start from beginning, use "$" for new messages only
		Count:   int64(qp.batchSize),
		Block:   100 * time.Millisecond, // Block for 100ms if no messages
	}).Result()
	
	if err != nil && err != redis.Nil {
		log.Printf("Failed to read from queue stream: %v", err)
		return apps
	}

	if len(streams) == 0 {
		return apps
	}

	// Process messages from stream
	stream := streams[0]
	messageIDs := make([]string, 0, len(stream.Messages))
	
	for _, msg := range stream.Messages {
		messageIDs = append(messageIDs, msg.ID)
		
		// Extract app data
		dataStr, ok := msg.Values["data"].(string)
		if !ok {
			log.Printf("Invalid message data format")
			continue
		}
		
		var app AppData
		if err := json.Unmarshal([]byte(dataStr), &app); err != nil {
			log.Printf("Failed to unmarshal app data: %v", err)
			continue
		}
		
		apps = append(apps, app)
	}

	// Acknowledge messages (mark as processed) - in production, use XACK with consumer groups
	// For now, we'll delete processed messages
	if len(messageIDs) > 0 {
		// In a real implementation with consumer groups, use XACK
		// For simplicity, we'll keep messages in stream for now (can be trimmed later)
		// qp.redis.XAck(qp.ctx, "app_queue", "workers", messageIDs...)
	}

	// Update metrics
	qp.metrics.decrementQueueSize(int64(len(apps)))

	return apps
}

// processBatch processes a batch of apps
func (qp *QueueProcessor) processBatch(apps []AppData) {
	start := time.Now()

	// Group apps by company for efficient processing
	companyGroups := qp.groupAppsByCompany(apps)

	// Process each company's apps
	for companyID, companyApps := range companyGroups {
		go qp.processCompanyApps(companyID, companyApps)
	}

	// Wait for all companies to complete
	// In a real implementation, you'd use a WaitGroup here

	duration := time.Since(start)
	log.Printf("Processed batch of %d apps in %v", len(apps), duration)

	// Update metrics
	qp.metrics.incrementProcessedTotal(int64(len(apps)))
}

// groupAppsByCompany groups apps by company ID
func (qp *QueueProcessor) groupAppsByCompany(apps []AppData) map[string][]AppData {
	groups := make(map[string][]AppData)

	for _, app := range apps {
		groups[app.CompanyID] = append(groups[app.CompanyID], app)
	}

	return groups
}

// processCompanyApps processes apps for a specific company
func (qp *QueueProcessor) processCompanyApps(companyID string, apps []AppData) {
	log.Printf("Processing %d apps for company %s", len(apps), companyID)

	// Send to Python enrichment service
	results, err := qp.enrichApps(apps)
	if err != nil {
		log.Printf("Failed to enrich apps for company %s: %v", companyID, err)
		qp.metrics.incrementErrorsTotal(int64(len(apps)))
		return
	}

	// Store results in database
	err = qp.storeResults(results)
	if err != nil {
		log.Printf("Failed to store results for company %s: %v", companyID, err)
		qp.metrics.incrementErrorsTotal(int64(len(results)))
		return
	}

	log.Printf("Successfully processed %d apps for company %s", len(results), companyID)
}

// enrichApps sends apps to Python enrichment service
func (qp *QueueProcessor) enrichApps(apps []AppData) ([]ProcessedResult, error) {
	// In a real implementation, this would make HTTP calls to the Python service
	// For now, we'll simulate the enrichment process

	var results []ProcessedResult

	for _, app := range apps {
		// Simulate enrichment delay
		time.Sleep(10 * time.Millisecond)

		// Simulate CVE data
		vulnerabilities := qp.simulateCVEData(app)

		result := ProcessedResult{
			AppID:           app.ID,
			CompanyID:       app.CompanyID,
			AgentID:         app.AgentID,
			Vulnerabilities: vulnerabilities,
			EnrichedAt:      time.Now(),
		}

		results = append(results, result)
	}

	return results, nil
}

// simulateCVEData simulates CVE enrichment (replace with real implementation)
func (qp *QueueProcessor) simulateCVEData(app AppData) []Vulnerability {
	// Simulate finding vulnerabilities based on app name
	if app.AppName == "nginx" && app.AppVersion < "1.20.0" {
		return []Vulnerability{
			{
				CVEID:       "CVE-2021-23017",
				Severity:    "high",
				CVSSScore:   8.1,
				Title:       "Nginx vulnerability in " + app.AppName,
				Description: "Simulated vulnerability for " + app.AppName + " " + app.AppVersion,
			},
		}
	}

	return []Vulnerability{}
}

// storeResults stores processed results in database
func (qp *QueueProcessor) storeResults(results []ProcessedResult) error {
	// In a real implementation, this would store in PostgreSQL
	// For now, we'll just log the results

	for _, result := range results {
		log.Printf("Storing result for app %s: %d vulnerabilities found",
			result.AppID, len(result.Vulnerabilities))
	}

	return nil
}

// metricsCollector collects and reports metrics
func (qp *QueueProcessor) metricsCollector() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-qp.ctx.Done():
			return
		case <-ticker.C:
			qp.reportMetrics()
		}
	}
}

// reportMetrics reports current metrics
func (qp *QueueProcessor) reportMetrics() {
	metrics := qp.metrics.getMetrics()

	log.Printf("Queue Metrics - Processed: %d total, %d today, Errors: %d, Queue Size: %d",
		metrics.processedTotal, metrics.processedToday, metrics.errorsTotal, metrics.queueSize)
}

// cleanupRoutine performs periodic cleanup
func (qp *QueueProcessor) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-qp.ctx.Done():
			return
		case <-ticker.C:
			qp.cleanup()
		}
	}
}

// cleanup performs cleanup tasks
func (qp *QueueProcessor) cleanup() {
	// Trim stream to keep only recent messages (last 24 hours worth)
	// XTRIM removes old entries from stream
	cutoff := time.Now().Add(-24 * time.Hour)
	
	// Use MINID to trim entries older than cutoff
	// Note: This requires converting timestamp to message ID format
	// For simplicity, we'll use MAXLEN to keep stream size manageable
	removed, err := qp.redis.XTrimMaxLen(qp.ctx, "app_queue", 10000).Result()
	if err != nil {
		log.Printf("Failed to trim queue stream: %v", err)
	} else if removed > 0 {
		log.Printf("Trimmed %d old entries from queue stream", removed)
	}
	
	_ = cutoff // Keep for future MINID-based trimming

	// Reset daily metrics
	qp.metrics.resetDailyMetrics()
}

// QueueMetrics methods
func (qm *QueueMetrics) incrementProcessedTotal(count int64) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.processedTotal += count
	qm.processedToday += count
}

func (qm *QueueMetrics) incrementErrorsTotal(count int64) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.errorsTotal += count
}

func (qm *QueueMetrics) incrementQueueSize() {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.queueSize++
}

func (qm *QueueMetrics) decrementQueueSize(count int64) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.queueSize -= count
	if qm.queueSize < 0 {
		qm.queueSize = 0
	}
}

func (qm *QueueMetrics) getMetrics() struct {
	processedTotal int64
	processedToday int64
	errorsTotal    int64
	queueSize      int64
} {
	qm.mu.RLock()
	defer qm.mu.RUnlock()

	return struct {
		processedTotal int64
		processedToday int64
		errorsTotal    int64
		queueSize      int64
	}{
		processedTotal: qm.processedTotal,
		processedToday: qm.processedToday,
		errorsTotal:    qm.errorsTotal,
		queueSize:      qm.queueSize,
	}
}

func (qm *QueueMetrics) resetDailyMetrics() {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.processedToday = 0
}
