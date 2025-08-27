package optimization

import (
	"context"
	"database/sql"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

// PerformanceOptimizer provides comprehensive performance optimizations
type PerformanceOptimizer struct {
	logger *zap.Logger

	// Connection pools
	dbPool    *sqlx.DB
	redisPool *redis.Client

	// Caches
	queryCache   *QueryCache
	resultCache  *ResultCache
	companyCache *CompanyCache

	// Rate limiters
	rateLimiters   map[string]*RateLimiter
	rateLimitersMu sync.RWMutex

	// Semaphores for concurrency control
	dbSemaphore    *semaphore.Weighted
	redisSemaphore *semaphore.Weighted

	// Metrics
	metrics *PerformanceMetrics
}

// PerformanceMetrics tracks performance metrics
type PerformanceMetrics struct {
	mu sync.RWMutex

	// Query performance
	queryCacheHits   int64
	queryCacheMisses int64
	queryDuration    time.Duration

	// Database performance
	dbConnections   int64
	dbQueryDuration time.Duration
	dbQueryErrors   int64

	// Redis performance
	redisHits     int64
	redisMisses   int64
	redisDuration time.Duration
	redisErrors   int64

	// Memory usage
	memoryAlloc uint64
	memoryHeap  uint64
	goroutines  int
}

// NewPerformanceOptimizer creates a new performance optimizer
func NewPerformanceOptimizer(logger *zap.Logger, dbPool *sqlx.DB, redisPool *redis.Client) *PerformanceOptimizer {
	po := &PerformanceOptimizer{
		logger:         logger,
		dbPool:         dbPool,
		redisPool:      redisPool,
		queryCache:     NewQueryCache(1000),   // 1000 cached queries
		resultCache:    NewResultCache(10000), // 10000 cached results
		companyCache:   NewCompanyCache(100),  // 100 companies
		rateLimiters:   make(map[string]*RateLimiter),
		dbSemaphore:    semaphore.NewWeighted(100), // Max 100 concurrent DB operations
		redisSemaphore: semaphore.NewWeighted(200), // Max 200 concurrent Redis operations
		metrics:        &PerformanceMetrics{},
	}

	// Start performance monitoring
	go po.monitorPerformance()

	return po
}

// QueryCache provides query result caching
type QueryCache struct {
	cache map[string]*CachedQuery
	mu    sync.RWMutex
	size  int
}

type CachedQuery struct {
	Query     string
	Result    interface{}
	ExpiresAt time.Time
	HitCount  int64
}

func NewQueryCache(size int) *QueryCache {
	return &QueryCache{
		cache: make(map[string]*CachedQuery),
		size:  size,
	}
}

func (qc *QueryCache) Get(key string) (interface{}, bool) {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	if cached, exists := qc.cache[key]; exists && time.Now().Before(cached.ExpiresAt) {
		cached.HitCount++
		return cached.Result, true
	}

	return nil, false
}

func (qc *QueryCache) Set(key string, result interface{}, ttl time.Duration) {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	// Evict if cache is full
	if len(qc.cache) >= qc.size {
		qc.evictLRU()
	}

	qc.cache[key] = &CachedQuery{
		Query:     key,
		Result:    result,
		ExpiresAt: time.Now().Add(ttl),
		HitCount:  0,
	}
}

func (qc *QueryCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, cached := range qc.cache {
		if oldestKey == "" || cached.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.ExpiresAt
		}
	}

	if oldestKey != "" {
		delete(qc.cache, oldestKey)
	}
}

// ResultCache provides result caching
type ResultCache struct {
	cache map[string]*CachedResult
	mu    sync.RWMutex
	size  int
}

type CachedResult struct {
	Data      interface{}
	ExpiresAt time.Time
	HitCount  int64
}

func NewResultCache(size int) *ResultCache {
	return &ResultCache{
		cache: make(map[string]*CachedResult),
		size:  size,
	}
}

func (rc *ResultCache) Get(key string) (interface{}, bool) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()

	if cached, exists := rc.cache[key]; exists && time.Now().Before(cached.ExpiresAt) {
		cached.HitCount++
		return cached.Data, true
	}

	return nil, false
}

func (rc *ResultCache) Set(key string, data interface{}, ttl time.Duration) {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	// Evict if cache is full
	if len(rc.cache) >= rc.size {
		rc.evictLRU()
	}

	rc.cache[key] = &CachedResult{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
		HitCount:  0,
	}
}

func (rc *ResultCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time

	for key, cached := range rc.cache {
		if oldestKey == "" || cached.ExpiresAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = cached.ExpiresAt
		}
	}

	if oldestKey != "" {
		delete(rc.cache, oldestKey)
	}
}

// CompanyCache provides company data caching
type CompanyCache struct {
	cache map[string]*CachedCompany
	mu    sync.RWMutex
	size  int
}

type CachedCompany struct {
	CompanyID  string
	Data       interface{}
	ExpiresAt  time.Time
	LastAccess time.Time
}

func NewCompanyCache(size int) *CompanyCache {
	return &CompanyCache{
		cache: make(map[string]*CachedCompany),
		size:  size,
	}
}

func (cc *CompanyCache) Get(companyID string) (interface{}, bool) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	if cached, exists := cc.cache[companyID]; exists && time.Now().Before(cached.ExpiresAt) {
		cached.LastAccess = time.Now()
		return cached.Data, true
	}

	return nil, false
}

func (cc *CompanyCache) Set(companyID string, data interface{}, ttl time.Duration) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// Evict if cache is full
	if len(cc.cache) >= cc.size {
		cc.evictLRU()
	}

	cc.cache[companyID] = &CachedCompany{
		CompanyID:  companyID,
		Data:       data,
		ExpiresAt:  time.Now().Add(ttl),
		LastAccess: time.Now(),
	}
}

func (cc *CompanyCache) evictLRU() {
	var oldestKey string
	var oldestAccess time.Time

	for key, cached := range cc.cache {
		if oldestKey == "" || cached.LastAccess.Before(oldestAccess) {
			oldestKey = key
			oldestAccess = cached.LastAccess
		}
	}

	if oldestKey != "" {
		delete(cc.cache, oldestKey)
	}
}

// RateLimiter provides rate limiting
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Clean old requests
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}
		rl.requests[key] = validRequests
	}

	// Check if limit exceeded
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// OptimizedQuery executes optimized database queries
func (po *PerformanceOptimizer) OptimizedQuery(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()

	// Check query cache first
	cacheKey := fmt.Sprintf("%s:%v", query, args)
	if cached, hit := po.queryCache.Get(cacheKey); hit {
		po.metrics.incrementQueryCacheHits()
		return cached.(*sql.Rows), nil
	}

	po.metrics.incrementQueryCacheMisses()

	// Acquire semaphore for database access
	if err := po.dbSemaphore.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("failed to acquire DB semaphore: %w", err)
	}
	defer po.dbSemaphore.Release(1)

	// Execute query with timeout
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := po.dbPool.QueryContext(queryCtx, query, args...)
	if err != nil {
		po.metrics.incrementDBQueryErrors()
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Cache result for 5 minutes
	po.queryCache.Set(cacheKey, rows, 5*time.Minute)

	po.metrics.recordDBQueryDuration(time.Since(start))
	return rows, nil
}

// OptimizedGet executes optimized GET operations with caching
func (po *PerformanceOptimizer) OptimizedGet(ctx context.Context, key string) (string, error) {
	start := time.Now()

	// Check result cache first
	if cached, hit := po.resultCache.Get(key); hit {
		po.metrics.incrementRedisHits()
		return cached.(string), nil
	}

	po.metrics.incrementRedisMisses()

	// Acquire semaphore for Redis access
	if err := po.redisSemaphore.Acquire(ctx, 1); err != nil {
		return "", fmt.Errorf("failed to acquire Redis semaphore: %w", err)
	}
	defer po.redisSemaphore.Release(1)

	// Execute Redis GET with timeout
	redisCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	result, err := po.redisPool.Get(redisCtx, key).Result()
	if err != nil {
		po.metrics.incrementRedisErrors()
		return "", fmt.Errorf("Redis GET failed: %w", err)
	}

	// Cache result for 1 minute
	po.resultCache.Set(key, result, 1*time.Minute)

	po.metrics.recordRedisDuration(time.Since(start))
	return result, nil
}

// OptimizedSet executes optimized SET operations
func (po *PerformanceOptimizer) OptimizedSet(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// Acquire semaphore for Redis access
	if err := po.redisSemaphore.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire Redis semaphore: %w", err)
	}
	defer po.redisSemaphore.Release(1)

	// Execute Redis SET with timeout
	redisCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	err := po.redisPool.Set(redisCtx, key, value, ttl).Err()
	if err != nil {
		po.metrics.incrementRedisErrors()
		return fmt.Errorf("Redis SET failed: %w", err)
	}

	// Update result cache
	po.resultCache.Set(key, value, ttl)

	return nil
}

// GetCompanyData gets company data with caching
func (po *PerformanceOptimizer) GetCompanyData(companyID string) (interface{}, bool) {
	return po.companyCache.Get(companyID)
}

// SetCompanyData sets company data with caching
func (po *PerformanceOptimizer) SetCompanyData(companyID string, data interface{}) {
	po.companyCache.Set(companyID, data, 30*time.Minute)
}

// CheckRateLimit checks rate limit for a key
func (po *PerformanceOptimizer) CheckRateLimit(key string, limit int, window time.Duration) bool {
	limiterKey := fmt.Sprintf("%s:%d:%v", key, limit, window)

	po.rateLimitersMu.Lock()
	if _, exists := po.rateLimiters[limiterKey]; !exists {
		po.rateLimiters[limiterKey] = NewRateLimiter(limit, window)
	}
	po.rateLimitersMu.Unlock()

	return po.rateLimiters[limiterKey].Allow(key)
}

// PerformanceMiddleware provides performance monitoring middleware
func (po *PerformanceOptimizer) PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Add performance context
		c.Set("performance_optimizer", po)

		// Process request
		c.Next()

		// Record performance metrics
		duration := time.Since(start)

		// Log slow requests
		if duration > 100*time.Millisecond {
			po.logger.Warn("Slow request detected",
				zap.String("method", c.Request.Method),
				zap.String("path", c.FullPath()),
				zap.Duration("duration", duration),
			)
		}
	}
}

// monitorPerformance monitors system performance
func (po *PerformanceOptimizer) monitorPerformance() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		po.metrics.updateMemoryStats(m)
		po.metrics.updateGoroutines(runtime.NumGoroutine())

		// Log performance metrics
		po.logger.Info("Performance metrics",
			zap.Int64("query_cache_hits", po.metrics.getQueryCacheHits()),
			zap.Int64("query_cache_misses", po.metrics.getQueryCacheMisses()),
			zap.Duration("avg_query_duration", po.metrics.getAvgQueryDuration()),
			zap.Int64("redis_hits", po.metrics.getRedisHits()),
			zap.Int64("redis_misses", po.metrics.getRedisMisses()),
			zap.Uint64("memory_alloc", m.Alloc),
			zap.Uint64("memory_heap", m.HeapAlloc),
			zap.Int("goroutines", runtime.NumGoroutine()),
		)

		// Force garbage collection if memory usage is high
		if m.HeapAlloc > 100*1024*1024 { // 100MB
			runtime.GC()
		}
	}
}

// PerformanceMetrics methods
func (pm *PerformanceMetrics) incrementQueryCacheHits() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.queryCacheHits++
}

func (pm *PerformanceMetrics) incrementQueryCacheMisses() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.queryCacheMisses++
}

func (pm *PerformanceMetrics) recordDBQueryDuration(duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.queryDuration = duration
}

func (pm *PerformanceMetrics) incrementDBQueryErrors() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.dbQueryErrors++
}

func (pm *PerformanceMetrics) incrementRedisHits() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.redisHits++
}

func (pm *PerformanceMetrics) incrementRedisMisses() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.redisMisses++
}

func (pm *PerformanceMetrics) recordRedisDuration(duration time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.redisDuration = duration
}

func (pm *PerformanceMetrics) incrementRedisErrors() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.redisErrors++
}

func (pm *PerformanceMetrics) updateMemoryStats(m runtime.MemStats) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.memoryAlloc = m.Alloc
	pm.memoryHeap = m.HeapAlloc
}

func (pm *PerformanceMetrics) updateGoroutines(count int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.goroutines = count
}

func (pm *PerformanceMetrics) getQueryCacheHits() int64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.queryCacheHits
}

func (pm *PerformanceMetrics) getQueryCacheMisses() int64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.queryCacheMisses
}

func (pm *PerformanceMetrics) getAvgQueryDuration() time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.queryDuration
}

func (pm *PerformanceMetrics) getRedisHits() int64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.redisHits
}

func (pm *PerformanceMetrics) getRedisMisses() int64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.redisMisses
}
