package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/handlers"
	"zerotrace/api/internal/middleware"
	"zerotrace/api/internal/repository"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := repository.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	scanRepo := repository.NewScanRepository(db.DB)

	// Initialize services
	scanService := services.NewScanService(cfg, scanRepo)
	agentService := services.NewAgentService()
	enrollmentService := services.NewEnrollmentService(cfg)

	// Setup router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Setup middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())

	// Setup routes
	setupRoutes(router, scanService, agentService, enrollmentService)

	// Create server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting ZeroTrace API server on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRoutes(router *gin.Engine, scanService *services.ScanService, agentService *services.AgentService, enrollmentService *services.EnrollmentService) {
	// Health check
	router.GET("/health", handlers.HealthCheck)

	// Agent routes (public - no auth required)
	agents := router.Group("/api/agents")
	{
		agents.POST("/register", handlers.RegisterAgent(agentService))
		agents.POST("/heartbeat", handlers.AgentHeartbeat(agentService))
		agents.POST("/results", handlers.AgentResults(agentService))
		agents.POST("/status", handlers.AgentStatus(agentService))
		agents.GET("/", handlers.GetAgents(agentService))
		agents.GET("/online", handlers.GetOnlineAgents(agentService))
		agents.GET("/stats", handlers.GetAgentStats(agentService))
		agents.GET("/stats/public", handlers.GetPublicAgentStats(agentService))
	}

	// Public dashboard routes (no auth required)
	dashboard := router.Group("/api/dashboard")
	{
		dashboard.GET("/overview", handlers.GetPublicDashboardOverview(agentService))
	}

	// Public vulnerabilities route (no auth required)
	vulnerabilities := router.Group("/api/vulnerabilities")
	{
		vulnerabilities.GET("/", handlers.GetPublicVulnerabilities(agentService))
	}

	// Enrollment routes (public - no auth required)
	enrollment := router.Group("/api/enrollment")
	{
		enrollment.POST("/enroll", handlers.EnrollAgent(enrollmentService))
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Note: Authentication is now handled by Clerk
		// No custom auth routes needed - users authenticate via Clerk frontend

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.ClerkAuth())
		{
			// Scan routes
			scans := protected.Group("/scans")
			{
				scans.GET("/", handlers.GetScans(scanService))
				scans.POST("/", handlers.CreateScan(scanService))
				scans.GET("/:id", handlers.GetScan(scanService))
				scans.PUT("/:id", handlers.UpdateScan(scanService))
				scans.DELETE("/:id", handlers.DeleteScan(scanService))
			}

			// Company routes
			companies := protected.Group("/companies")
			{
				companies.GET("/:id", handlers.GetCompany)
				companies.PUT("/:id", handlers.UpdateCompany)
			}

			// Vulnerability routes (commented out until handlers are implemented)
			// vulnerabilities := protected.Group("/vulnerabilities")
			// {
			// 	vulnerabilities.GET("/", handlers.GetVulnerabilities)
			// 	vulnerabilities.GET("/:id", handlers.GetVulnerabilityDetails)
			// 	vulnerabilities.GET("/stats", handlers.GetVulnerabilityStats)
			// }

			// Dashboard routes
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/overview", handlers.GetDashboardOverview)
				dashboard.GET("/trends", handlers.GetVulnerabilityTrends)
			}

			// Enrollment management routes (protected)
			enrollment := protected.Group("/enrollment")
			{
				enrollment.POST("/tokens", handlers.GenerateEnrollmentToken(enrollmentService))
				enrollment.DELETE("/tokens/:id", handlers.RevokeEnrollmentToken(enrollmentService))
				enrollment.DELETE("/credentials/:id", handlers.RevokeAgentCredential(enrollmentService))
			}
		}
	}
}
