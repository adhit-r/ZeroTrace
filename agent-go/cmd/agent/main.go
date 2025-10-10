package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zerotrace/agent/internal/communicator"
	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/processor"
	"zerotrace/agent/internal/scanner"
	"zerotrace/agent/internal/tray"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize components
	scanner := scanner.NewSoftwareScanner(cfg) // Use software scanner instead of code scanner
	processor := processor.NewProcessor(cfg)
	communicator := communicator.NewCommunicator(cfg)

	// Initialize simple tray manager (MDM-friendly)
	trayManager := tray.NewSimpleTrayManager()
	trayManager.Start()

	testTray := flag.Bool("test-tray", false, "Run in tray test mode")
	flag.Parse()

	if *testTray {
		select {}
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start agent
	log.Println("Starting ZeroTrace Software Vulnerability Agent...")
	log.Printf("Agent ID: %s", cfg.AgentID)
	log.Printf("API Endpoint: %s", cfg.APIEndpoint)
	log.Printf("Organization ID: %s", cfg.OrganizationID)

	// Check if agent needs enrollment
	if !cfg.IsEnrolled() {
		if cfg.HasEnrollmentToken() {
			log.Println("Enrolling agent with enrollment token...")
			if err := communicator.EnrollAgent(); err != nil {
				log.Printf("Failed to enroll agent: %v", err)
				log.Println("Agent will continue with legacy registration...")
			} else {
				log.Println("Agent enrolled successfully")
				log.Printf("Organization ID: %s", cfg.OrganizationID)
			}
		} else {
			log.Println("No enrollment token found, using legacy registration...")
		}
	} else {
		log.Println("Agent already enrolled")
	}

	// Legacy registration (fallback)
	if !cfg.IsEnrolled() {
		if err := communicator.RegisterAgent(); err != nil {
			log.Printf("Warning: Failed to register agent: %v", err)
		} else {
			log.Println("Agent registered successfully (legacy)")
		}
	}

	// Start scanning in a goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Scan for installed software
				results, err := scanner.Scan()
				if err != nil {
					log.Printf("Scan error: %v", err)
					time.Sleep(30 * time.Second)
					continue
				}

				log.Printf("Found %d installed applications", len(results.Dependencies))

				// Process results
				processedResults, err := processor.Process(results)
				if err != nil {
					log.Printf("Processing error: %v", err)
					continue
				}

				// Send results to API
				if err := communicator.SendResults(processedResults); err != nil {
					log.Printf("Communication error: %v", err)
				} else {
					log.Printf("Successfully sent scan results to API")
				}

				// Wait before next scan
				log.Printf("Next scan in %v", cfg.ScanInterval)
				time.Sleep(cfg.ScanInterval)
			}
		}
	}()

	// Start heartbeat in a goroutine
	go func() {
		ticker := time.NewTicker(30 * time.Second) // Send heartbeat every 30 seconds
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Get CPU and memory usage from monitor
				cpuUsage := 2.5     // Placeholder - would get from monitor
				memoryUsage := 45.2 // Placeholder - would get from monitor

				metadata := map[string]any{
					"scan_interval": cfg.ScanInterval.String(),
					"scan_depth":    cfg.ScanDepth,
					"version":       "1.0.0",
				}

				// Use enrollment-based heartbeat if enrolled, otherwise legacy
				if cfg.IsEnrolled() {
					if err := communicator.SendHeartbeatWithCredential(cpuUsage, memoryUsage, metadata); err != nil {
						log.Printf("Enrollment heartbeat error: %v", err)
					} else {
						log.Printf("Enrollment heartbeat sent successfully")
					}
				} else {
					if err := communicator.SendHeartbeat(cpuUsage, memoryUsage, metadata); err != nil {
						log.Printf("Legacy heartbeat error: %v", err)
					} else {
						log.Printf("Legacy heartbeat sent successfully")
					}
				}
			}
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down agent...")
	cancel()

	// Stop tray manager
	trayManager.Stop()

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		log.Println("Forced shutdown")
	case <-time.After(5 * time.Second):
		log.Println("Graceful shutdown completed")
	}
}
