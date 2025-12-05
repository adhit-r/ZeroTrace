package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"zerotrace/agent/internal/communicator"
	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/processor"
	"zerotrace/agent/internal/scanner"
	"zerotrace/agent/internal/tray"

	"github.com/joho/godotenv"
	"fyne.io/systray"
)

// noOpTrayManager is a no-op implementation for when tray is disabled
type noOpTrayManager struct{}

func (n *noOpTrayManager) Start() {}
func (n *noOpTrayManager) Stop()  {}

func main() {
	// Check if running as .app bundle (macOS)
	// If so, redirect logs to file to avoid showing terminal
	if runtime.GOOS == "darwin" {
		// Check if we're running from an .app bundle
		execPath, _ := os.Executable()
		if filepath.Ext(filepath.Dir(filepath.Dir(filepath.Dir(execPath)))) == ".app" {
			// Running as .app bundle - redirect logs to file
			logDir := filepath.Join(os.Getenv("HOME"), ".zerotrace", "logs")
			os.MkdirAll(logDir, 0755)
			logFile := filepath.Join(logDir, "agent.log")
			
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err == nil {
				log.SetOutput(file)
				// Also redirect stderr to avoid crash dialogs
				os.Stderr = file
			}
		}
	}

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize components
	softwareScanner := scanner.NewSoftwareScanner(cfg)
	systemScanner := scanner.NewSystemScanner(cfg)
	networkScanner := scanner.NewNetworkScanner(cfg)
	processor := processor.NewProcessor(cfg)
	communicator := communicator.NewCommunicator(cfg)

	// Parse flags
	disableTray := flag.Bool("no-tray", false, "Disable system tray UI")
	testTray := flag.Bool("test-tray", false, "Run in tray test mode")
	flag.Parse()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start agent initialization
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

	// Function to start all background agent work
	startAgentWork := func() {
		// Start software scanning in a goroutine
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Perform scan
					results, err := softwareScanner.Scan()
					if err != nil {
						log.Printf("Scan error: %v", err)
						time.Sleep(cfg.ScanInterval)
						continue
					}

					// Process results
					processedResults, err := processor.Process(results)
					if err != nil {
						log.Printf("Processing error: %v", err)
						time.Sleep(cfg.ScanInterval)
						continue
					}

					// Send results to API
					if err := communicator.SendResults(processedResults); err != nil {
						log.Printf("Communication error: %v", err)
					} else {
						log.Printf("Successfully sent software scan results to API")
					}

					// Wait before next scan
					log.Printf("Next scan in %v", cfg.ScanInterval)
					time.Sleep(cfg.ScanInterval)
				}
			}
		}()

		// Start system info scanning in a goroutine
		go func() {
			// Perform an initial scan right away
			sendSystemInfo(ctx, systemScanner, communicator)

			// Then scan on a longer interval
			ticker := time.NewTicker(1 * time.Hour)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					sendSystemInfo(ctx, systemScanner, communicator)
				}
			}
		}()

		// Start network scanning in a goroutine (if enabled)
		if cfg.NetworkScanEnabled {
			go func() {
				// Perform an initial scan after a short delay
				time.Sleep(30 * time.Second)
				sendNetworkScan(ctx, networkScanner, communicator)

				// Then scan on configured interval
				ticker := time.NewTicker(cfg.NetworkScanInterval)
				defer ticker.Stop()

				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						sendNetworkScan(ctx, networkScanner, communicator)
					}
				}
			}()
			log.Printf("Network scanning enabled (interval: %v)", cfg.NetworkScanInterval)
		} else {
			log.Println("Network scanning disabled")
		}

		// Start heartbeat in a goroutine
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					cpuUsage := 2.5
					memoryUsage := 45.2

					metadata := map[string]any{
						"scan_interval": cfg.ScanInterval.String(),
						"scan_depth":    cfg.ScanDepth,
						"version":       "1.0.0",
					}

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
	}

	// Handle tray UI
	var trayManager interface {
		Start()
		Stop()
	}

	if !*disableTray && runtime.GOOS == "darwin" {
		// macOS: systray MUST run on main thread
		// Lock the OS thread for systray
		runtime.LockOSThread()
		
		// Create tray manager
		trayMgr := tray.NewSimpleTrayManager()
		trayManager = trayMgr
		
		// Define onReady callback that starts agent work
		onReady := func() {
			// Call the tray manager's OnReady to set up the menu
			trayMgr.OnReady()
			
			// Start all background agent work from onReady
			// This ensures systray is initialized before we start background tasks
			startAgentWork()
		}
		
		onExit := func() {
			trayMgr.OnExit()
			cancel()
		}
		
		// Run systray on main thread - this blocks until systray.Quit() is called
		log.Println("Starting systray on main thread (macOS)...")
		systray.Run(onReady, onExit)
		
		// After systray exits, shutdown
		log.Println("Shutting down agent...")
		cancel()
		trayManager.Stop()
		
	} else if !*disableTray {
		// Non-macOS: can run systray in goroutine
		trayManager = tray.NewSimpleTrayManager()
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Tray UI crashed: %v", r)
					log.Println("Agent will continue running without tray icon")
				}
			}()
			trayManager.Start()
		}()
		
		// Start agent work normally
		startAgentWork()
		
		// Wait for interrupt signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		
		log.Println("Shutting down agent...")
		cancel()
		trayManager.Stop()
		
	} else {
		// Tray disabled
		log.Println("Tray UI disabled (--no-tray flag)")
		trayManager = &noOpTrayManager{}
		
		// Start agent work normally
		startAgentWork()
		
		// Wait for interrupt signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		
		log.Println("Shutting down agent...")
		cancel()
	}

	if *testTray {
		select {}
	}

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

func sendSystemInfo(ctx context.Context, systemScanner *scanner.SystemScanner, communicator *communicator.Communicator) {
	log.Println("Scanning for system information...")
	sysInfo, err := systemScanner.Scan()
	if err != nil {
		log.Printf("System info scan error: %v", err)
		return
	}

	if err := communicator.SendSystemInfo(sysInfo); err != nil {
		log.Printf("Failed to send system info: %v", err)
	} else {
		log.Println("Successfully sent system information to API.")
	}
}

// sendNetworkScan performs agentless network scanning
// Note: This is AGENTLESS scanning - the ZeroTrace agent is the SCANNING HOST,
// not something installed on target devices. It scans other devices on the network
// using network protocols (Nmap, Nuclei) without requiring any agent installation
// on the target devices. Similar to how Tenable sensors work.
func sendNetworkScan(ctx context.Context, networkScanner *scanner.NetworkScanner, communicator *communicator.Communicator) {
	log.Println("Starting agentless network scan...")
	scanResult, err := networkScanner.ScanLocalNetwork()
	if err != nil {
		log.Printf("Network scan error: %v", err)
		return
	}

	totalHosts := 0
	if hosts, ok := scanResult.Metadata["total_hosts"].(int); ok {
		totalHosts = hosts
	} else if hosts, ok := scanResult.Metadata["total_hosts"].(float64); ok {
		totalHosts = int(hosts)
	}
	log.Printf("Network scan completed: %d findings on %d hosts",
		len(scanResult.NetworkFindings),
		totalHosts)

	if err := communicator.SendNetworkScanResults(scanResult); err != nil {
		log.Printf("Failed to send network scan results: %v", err)
	} else {
		log.Println("Successfully sent network scan results to API.")
	}
}
