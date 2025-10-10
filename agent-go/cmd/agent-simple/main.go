package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zerotrace/agent/internal/communicator"
	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/processor"
	"zerotrace/agent/internal/scanner"

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
	softwareScanner := scanner.NewSoftwareScanner(cfg)
	configScanner := scanner.NewConfigScanner(cfg)
	processor := processor.NewProcessor(cfg)
	communicator := communicator.NewCommunicator(cfg)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start agent
	log.Println("Starting ZeroTrace Software Vulnerability Agent (MDM Mode)...")
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
				softwareResults, err := softwareScanner.Scan()
				if err != nil {
					log.Printf("Software scan error: %v", err)
					time.Sleep(30 * time.Second)
					continue
				}

				log.Printf("Found %d installed applications", len(softwareResults.Dependencies))

				// Scan for configuration vulnerabilities
				configResults, err := configScanner.Scan()
				if err != nil {
					log.Printf("Configuration scan error: %v", err)
				} else {
					log.Printf("Found %d configuration vulnerabilities", len(configResults.Vulnerabilities))
				}

				// Process software results
				processedResults, err := processor.Process(softwareResults)
				if err != nil {
					log.Printf("Processing error: %v", err)
					continue
				}

				// Send software results to API
				if err := communicator.SendResults(processedResults); err != nil {
					log.Printf("Failed to send software results: %v", err)
				} else {
					log.Printf("Software results sent successfully")
				}

				// Send configuration results to API
				if configResults != nil {
					if err := communicator.SendResults(configResults); err != nil {
						log.Printf("Failed to send configuration results: %v", err)
					} else {
						log.Printf("Configuration results sent successfully")
					}
				}

				// Wait for next scan interval
				time.Sleep(cfg.ScanInterval)
			}
		}
	}()

	// Start heartbeat in a goroutine
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				var err error
				// Simple metrics for MDM mode
				cpuUsage := 0.0
				memoryUsage := 0.0
				metadata := map[string]any{
					"mode": "mdm",
					"os":   cfg.OS,
				}

				if cfg.IsEnrolled() {
					err = communicator.SendHeartbeatWithCredential(cpuUsage, memoryUsage, metadata)
				} else {
					err = communicator.SendHeartbeat(cpuUsage, memoryUsage, metadata)
				}

				if err != nil {
					log.Printf("Heartbeat error: %v", err)
				} else {
					log.Printf("Heartbeat sent successfully")
				}
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutting down agent...")

	// Cancel context
	cancel()

	// Graceful shutdown
	log.Println("Graceful shutdown completed")
}
