package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PoolMonitor monitors database connection pool health
type PoolMonitor struct {
	pool *pgxpool.Pool
}

// NewPoolMonitor creates a new pool monitor
func NewPoolMonitor(pool *pgxpool.Pool) *PoolMonitor {
	return &PoolMonitor{
		pool: pool,
	}
}

// PoolStats represents connection pool statistics
type PoolStats struct {
	TotalConns     int32
	AcquiredConns  int32
	IdleConns      int32
	Constructing   int32
	MaxConns       int32
	AcquireCount   int64
	AcquireDuration time.Duration
	AcquiredConnsCount int64
	CanceledAcquireCount int64
}

// GetStats retrieves current pool statistics
func (pm *PoolMonitor) GetStats(ctx context.Context) (*PoolStats, error) {
	stats := pm.pool.Stat()
	
	return &PoolStats{
		TotalConns:     stats.TotalConns(),
		AcquiredConns:  stats.AcquiredConns(),
		IdleConns:      stats.IdleConns(),
		Constructing:   stats.ConstructingConns(),
		MaxConns:       stats.MaxConns(),
		AcquireCount:   stats.AcquireCount(),
		AcquireDuration: stats.AcquireDuration(),
		AcquiredConnsCount: int64(stats.AcquiredConns()),
		CanceledAcquireCount: stats.CanceledAcquireCount(),
	}, nil
}

// HealthCheck performs a health check on the connection pool
func (pm *PoolMonitor) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	stats, err := pm.GetStats(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pool stats: %w", err)
	}

	// Check if pool is healthy
	if stats.TotalConns == 0 && stats.Constructing == 0 {
		return fmt.Errorf("pool has no connections and none are being constructed")
	}

	// Check if pool is near capacity
	if stats.TotalConns >= stats.MaxConns*9/10 {
		return fmt.Errorf("pool is near capacity: %d/%d connections", stats.TotalConns, stats.MaxConns)
	}

	// Test connection
	conn, err := pm.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()

	// Ping database
	if err := conn.Ping(ctx); err != nil {
		return fmt.Errorf("connection ping failed: %w", err)
	}

	return nil
}

// MonitorPool continuously monitors pool health
func (pm *PoolMonitor) MonitorPool(ctx context.Context, interval time.Duration, callback func(*PoolStats, error)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			stats, err := pm.GetStats(ctx)
			callback(stats, err)
		}
	}
}

