package main
package main

import (
	"fmt"
	"log"
	"os"

	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/scanner"
)

func main() {
	// Set dummy config
	os.Setenv("AGENT_ID", "123e4567-e89b-12d3-a456-426614174000")
	os.Setenv("COMPANY_ID", "123e4567-e89b-12d3-a456-426614174001")

	cfg := &config.Config{
		AgentID:   "123e4567-e89b-12d3-a456-426614174000",
		CompanyID: "123e4567-e89b-12d3-a456-426614174001",
	}

	networkScanner := scanner.NewNetworkScanner(cfg)

	// Test scan on localhost
	result, err := networkScanner.Scan("127.0.0.1")
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	fmt.Printf("Scan completed in %v\n", result.EndTime.Sub(result.StartTime))
	fmt.Printf("Found %d network findings\n", len(result.NetworkFindings))

	for _, finding := range result.NetworkFindings {
		fmt.Printf("- %s: %s on %s:%d (%s)\n", finding.FindingType, finding.Description, finding.Host, finding.Port, finding.Severity)
	}
}