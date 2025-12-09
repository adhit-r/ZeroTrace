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
	analytics "zerotrace/api/internal/services/analytics"

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

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Validate environment variables
	if err := config.ValidateEnvironment(); err != nil {
		log.Fatalf("Environment validation failed: %v", err)
	}

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

	// Initialize config auditor repositories
	configFileRepo := repository.NewConfigFileRepository(db.DB)
	configFindingRepo := repository.NewConfigFindingRepository(db.DB)
	configStandardRepo := repository.NewConfigStandardRepository(db.DB)
	configAnalysisRepo := repository.NewConfigAnalysisRepository(db.DB)

	// Initialize services
	scanService := services.NewScanService(cfg, scanRepo)
	agentService := services.NewAgentService(db.DB)
	enrollmentService := services.NewEnrollmentService(cfg, db)
	vulnerabilityV2Service := services.NewVulnerabilityV2Service()
	organizationProfileService := services.NewOrganizationProfileService(db.DB)
	analyticsService := analytics.NewAnalyticsService(db.DB)
	enrichmentService := services.NewEnrichmentService(cfg.EnrichmentServiceURL)
	aiService := services.NewAIService(cfg.AIServiceURL)

	// Initialize config auditor services
	configParserService := services.NewConfigParserService(configFileRepo)
	configAnalyzerService := services.NewConfigAnalyzerService(configFileRepo, configFindingRepo, configStandardRepo, configAnalysisRepo)
	configJobService := services.NewConfigJobService(configFileRepo, configParserService, configAnalyzerService, cfg)
	configFileService := services.NewConfigFileService(cfg, configFileRepo, configParserService, configAnalyzerService, configJobService)
	configFindingService := services.NewConfigFindingService(configFindingRepo)
	configAnalysisService := services.NewConfigAnalysisService(configAnalysisRepo, configFileRepo)

	// Get underlying sql.DB for AttackPathService
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	attackPathService := services.NewAttackPathService(sqlDB)

	// Setup router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Setup middleware (order matters - correlation ID should be first)
	router.Use(middleware.CorrelationID())
	router.Use(middleware.CORS())
	router.Use(middleware.CompressionMiddleware()) // Add compression
	router.Use(middleware.ETagMiddleware())        // Add ETag support
	router.Use(middleware.InputValidationMiddleware())
	router.Use(middleware.RateLimitMiddleware(cfg))
	router.Use(middleware.RequestLogger())

	// Setup routes
	setupRoutes(router, db, scanService, agentService, enrollmentService, vulnerabilityV2Service, organizationProfileService, analyticsService, enrichmentService, aiService, configFileService, configFindingService, configAnalysisService, attackPathService)

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

	// Graceful shutdown - stop background workers first
	configJobService.Stop()

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRoutes(router *gin.Engine, db *repository.Database, scanService *services.ScanService, agentService *services.AgentService, enrollmentService *services.EnrollmentService, vulnerabilityV2Service *services.VulnerabilityV2Service, organizationProfileService *services.OrganizationProfileService, analyticsService *analytics.AnalyticsService, enrichmentService *services.EnrichmentService, aiService *services.AIService, configFileService *services.ConfigFileService, configFindingService *services.ConfigFindingService, configAnalysisService *services.ConfigAnalysisService, attackPathService *services.AttackPathService) {
	// Root route
	// router.GET("/", handlers.Root)

	// Health check
	router.GET("/health", handlers.HealthCheck(db))

	// Agent routes (public - no auth required)
	agents := router.Group("/api/agents")
	{
		agents.POST("/register", handlers.RegisterAgent(agentService))
		agents.POST("/heartbeat", handlers.AgentHeartbeat(agentService))
		agents.POST("/results", handlers.AgentResults(agentService, enrichmentService))
		agents.POST("/status", handlers.AgentStatus(agentService))
		agents.POST("/system-info", handlers.UpdateSystemInfo(agentService))
		agents.POST("/network-scan-results", handlers.NetworkScanResults(agentService))
		agents.GET("/", handlers.GetAgents(agentService))
		agents.GET("/:id", handlers.GetAgent(agentService))
		agents.GET("/online", handlers.GetOnlineAgents(agentService))
		agents.GET("/stats", handlers.GetAgentStats(agentService))
		agents.GET("/stats/public", handlers.GetPublicAgentStats(agentService))
		agents.GET("/processing-status", handlers.GetProcessingStatus(agentService))
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

	// Organization profile routes (public for now)
	organizationProfileHandler := handlers.NewOrganizationProfileHandler(organizationProfileService)
	organizations := router.Group("/api/organizations")
	{
		organizations.POST("/profile", organizationProfileHandler.CreateOrganizationProfile)
		organizations.GET("/:id/profile", organizationProfileHandler.GetOrganizationProfile)
		organizations.PUT("/:id/profile", organizationProfileHandler.UpdateOrganizationProfile)
		organizations.DELETE("/:id/profile", organizationProfileHandler.DeleteOrganizationProfile)
		organizations.GET("/:id/tech-stack/relevance", organizationProfileHandler.GetTechStackRelevance)
		organizations.GET("/:id/risk-weights", organizationProfileHandler.GetIndustryRiskWeights)
	}

	// Technology stack analysis routes (merged into organization profile)
	techStack := router.Group("/api/tech-stack")
	{
		techStack.GET("/organizations/:id/analyze", organizationProfileHandler.AnalyzeTechStack)
		techStack.GET("/organizations/:id/recommendations", organizationProfileHandler.GetTechStackRecommendations)
	}

	// AI-powered analysis routes (public for now)
	aiAnalysisHandler := handlers.NewAIAnalysisHandler(aiService)
	aiAnalysis := router.Group("/api/ai-analysis")
	{
		aiAnalysis.GET("/vulnerabilities/:id/comprehensive", aiAnalysisHandler.AnalyzeVulnerabilityComprehensive)
		aiAnalysis.GET("/vulnerabilities/trends", aiAnalysisHandler.AnalyzeVulnerabilityTrends)
		aiAnalysis.GET("/exploit-intelligence/:cve_id", aiAnalysisHandler.GetExploitIntelligence)
		aiAnalysis.GET("/vulnerabilities/:id/predictive", aiAnalysisHandler.GetPredictiveAnalysis)
		aiAnalysis.GET("/vulnerabilities/:id/remediation-plan", aiAnalysisHandler.GetRemediationPlan)
		aiAnalysis.POST("/bulk-analysis", aiAnalysisHandler.GetBulkAnalysis)
	}

	// Analytics routes (unified service for heatmap, maturity, compliance)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Risk heatmap routes (public for now)
	heatmaps := router.Group("/api/heatmaps")
	{
		heatmaps.GET("/organizations/:id", analyticsHandler.GenerateRiskHeatmap)
		heatmaps.GET("/organizations/:id/hotspots", analyticsHandler.GetHeatmapHotspots)
		heatmaps.GET("/organizations/:id/risk-distribution", analyticsHandler.GetRiskDistribution)
		heatmaps.GET("/organizations/:id/trends", analyticsHandler.GetHeatmapTrends)
		heatmaps.GET("/organizations/:id/recommendations", analyticsHandler.GetHeatmapRecommendations)
	}

	// Security maturity score routes (public for now)
	maturity := router.Group("/api/maturity")
	{
		maturity.GET("/organizations/:id/score", analyticsHandler.CalculateMaturityScore)
		maturity.GET("/organizations/:id/benchmark", analyticsHandler.GetMaturityBenchmark)
		maturity.GET("/organizations/:id/roadmap", analyticsHandler.GetImprovementRoadmap)
		maturity.GET("/organizations/:id/trends", analyticsHandler.GetMaturityTrends)
		maturity.GET("/organizations/:id/dimensions", analyticsHandler.GetDimensionScores)
	}

	// Compliance reporting routes (public for now)
	compliance := router.Group("/api/compliance")
	{
		compliance.GET("/organizations/:id/report", analyticsHandler.GenerateComplianceReport)
		compliance.GET("/organizations/:id/score", analyticsHandler.GetComplianceScore)
		compliance.GET("/organizations/:id/findings", analyticsHandler.GetComplianceFindings)
		compliance.GET("/organizations/:id/recommendations", analyticsHandler.GetComplianceRecommendations)
		compliance.GET("/organizations/:id/evidence", analyticsHandler.GetComplianceEvidence)
		compliance.GET("/organizations/:id/executive-summary", analyticsHandler.GetExecutiveSummary)
	}

	// API v2 routes (public - no auth required for now)
	v2 := router.Group("/api/v2")
	{
		// Vulnerability v2 routes
		vulnerabilityV2Handler := handlers.NewVulnerabilityV2Handler(vulnerabilityV2Service, agentService)
		v2Vulns := v2.Group("/vulnerabilities")
		{
			v2Vulns.GET("/", vulnerabilityV2Handler.GetVulnerabilitiesV2)
			v2Vulns.GET("/stats", vulnerabilityV2Handler.GetVulnerabilityStats)
			v2Vulns.GET("/export", vulnerabilityV2Handler.ExportVulnerabilities)
		}

		// Compliance routes
		v2Compliance := v2.Group("/compliance")
		{
			v2Compliance.GET("/status", vulnerabilityV2Handler.GetComplianceStatus)
		}

		// Scan routes
		v2Scans := v2.Group("/scans")
		{
			v2Scans.POST("/network", vulnerabilityV2Handler.InitiateNetworkScan)
			v2Scans.GET("/:scan_id/status", vulnerabilityV2Handler.GetScanStatus)
			v2Scans.GET("/:scan_id/results", vulnerabilityV2Handler.GetScanResults)
		}

		// Attack Path routes
		attackPathHandler := handlers.NewAttackPathHandler(attackPathService)
		v2AttackPaths := v2.Group("/attack-paths")
		{
			v2AttackPaths.GET("/", attackPathHandler.GetAttackPaths)
			v2AttackPaths.GET("/:path_id", attackPathHandler.GetAttackPath)
			v2AttackPaths.POST("/generate", attackPathHandler.GenerateAttackPaths)
		}

		// Config Auditor routes (public for now)
		configFileHandler := handlers.NewConfigFileHandler(configFileService)
		configFindingHandler := handlers.NewConfigFindingHandler(configFindingService)
		configAnalysisHandler := handlers.NewConfigAnalysisHandler(configAnalysisService)

		v2ConfigFiles := v2.Group("/config-files")
		{
			v2ConfigFiles.POST("/upload", configFileHandler.UploadConfigFile)
			v2ConfigFiles.GET("/", configFileHandler.ListConfigFiles)
			v2ConfigFiles.GET("/:id", configFileHandler.GetConfigFile)
			v2ConfigFiles.GET("/:id/content", configFileHandler.GetConfigFileContent)
			v2ConfigFiles.DELETE("/:id", configFileHandler.DeleteConfigFile)
			v2ConfigFiles.POST("/:id/analyze", configFileHandler.TriggerAnalysis)
		}

		v2ConfigFindings := v2.Group("/config-findings")
		{
			v2ConfigFindings.GET("/", configFindingHandler.ListConfigFindings)
			v2ConfigFindings.GET("/:id", configFindingHandler.GetConfigFinding)
			v2ConfigFindings.PATCH("/:id/status", configFindingHandler.UpdateFindingStatus)
			v2ConfigFindings.GET("/stats", configFindingHandler.GetFindingStats)
		}

		v2ConfigAnalysis := v2.Group("/config-files/:id")
		{
			v2ConfigAnalysis.GET("/analysis", configAnalysisHandler.GetAnalysisResults)
			v2ConfigAnalysis.GET("/compliance", configAnalysisHandler.GetComplianceScores)
			v2ConfigAnalysis.GET("/analysis/status", configAnalysisHandler.GetAnalysisStatus)
		}
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
