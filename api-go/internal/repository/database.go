package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents the database connection and repositories
type Database struct {
	DB          *gorm.DB
	PgxPool     *pgxpool.Pool // Direct pgx pool for hot paths
	Redis       *redis.Client
	QueryCache  *QueryCache
	StmtMgr     *PreparedStatementManager
	PoolMonitor *PoolMonitor
	TxManager   *TransactionManager
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	// Configure GORM logger
	gormLogger := logger.Default
	if cfg.Debug {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create pgx pool for direct PostgreSQL access (bypass GORM for hot paths)
	pgxConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	// Configure pool settings
	pgxConfig.MaxConns = 100
	pgxConfig.MinConns = 10
	pgxConfig.MaxConnLifetime = time.Hour
	pgxConfig.MaxConnIdleTime = 30 * time.Minute

	pgxPool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	// Test pgx pool connection
	if err := pgxPool.Ping(context.Background()); err != nil {
		pgxPool.Close()
		return nil, fmt.Errorf("failed to ping pgx pool: %w", err)
	}

	log.Println("Database connected successfully (GORM + pgx pool)")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		// We don't return error here to allow running without Redis (it will just fail cache operations)
	} else {
		log.Println("Redis connected successfully")
	}

	// Initialize query cache with Redis client
	queryCache := NewQueryCache(pgxPool, redisClient, 5*time.Minute)

	// Initialize prepared statement manager
	stmtMgr := NewPreparedStatementManager(pgxPool)

	if err := stmtMgr.InitializePreparedStatements(ctx); err != nil {
		log.Printf("Warning: Failed to initialize prepared statements: %v", err)
	}

	// Initialize pool monitor
	poolMonitor := NewPoolMonitor(pgxPool)

	// Initialize transaction manager
	txManager := NewTransactionManager(pgxPool)

	return &Database{
		DB:          db,
		PgxPool:     pgxPool,
		Redis:       redisClient,
		QueryCache:  queryCache,
		StmtMgr:     stmtMgr,
		PoolMonitor: poolMonitor,
		TxManager:   txManager,
	}, nil
}

// AutoMigrate runs database migrations
func (d *Database) AutoMigrate() error {
	log.Println("Running database migrations...")

	err := d.DB.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Organization{},
		&models.Scan{},
		&models.Vulnerability{},
		&models.Agent{},
		&models.Software{},
		&models.NetworkHost{},
		&models.EnrollmentToken{},
		&models.AgentCredential{},
		&models.DashboardSnapshot{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	var errs []error

	// Close GORM connection
	if sqlDB, err := d.DB.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close GORM connection: %w", err))
		}
	} else {
		errs = append(errs, fmt.Errorf("failed to get GORM sql.DB: %w", err))
	}

	// Close pgx pool
	if d.PgxPool != nil {
		d.PgxPool.Close()
	}

	// Close Redis
	if d.Redis != nil {
		if err := d.Redis.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close Redis: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing database: %v", errs)
	}

	return nil
}
