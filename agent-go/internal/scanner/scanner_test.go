package scanner

import (
	"os"
	"testing"
	"time"
	"zerotrace/agent/internal/config"
)

func setupTestConfig() *config.Config {
	// Set dummy env vars for config
	os.Setenv("AGENT_ID", "123e4567-e89b-12d3-a456-426614174000")
	os.Setenv("ZEROTRACE_ORGANIZATION_ID", "123e4567-e89b-12d3-a456-426614174001")
	return config.Load()
}

func TestSystemScanner_Scan(t *testing.T) {
	cfg := setupTestConfig()
	scanner := NewSystemScanner(cfg)

	// Test basic scan functionality
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("SystemScanner.Scan() failed: %v", err)
	}

	// Verify result structure
	if result == nil {
		t.Fatal("SystemScanner.Scan() returned nil result")
	}

	if result.Hostname == "" {
		t.Error("Hostname should not be empty")
	}

	// Test that scan completes within reasonable time
	start := time.Now()
	_, err = scanner.Scan()
	duration := time.Since(start)
	if err != nil {
		t.Fatalf("SystemScanner.Scan() failed on second run: %v", err)
	}

	if duration > 30*time.Second {
		t.Errorf("SystemScanner.Scan() took too long: %v", duration)
	}
}

func TestSoftwareScanner_Scan(t *testing.T) {
	cfg := setupTestConfig()
	scanner := NewSoftwareScanner(cfg)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("SoftwareScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("SoftwareScanner.Scan() returned nil result")
	}

	// Verify result structure
	// Note: result type depends on SoftwareScanner implementation.
	// Assuming it returns something non-nil is a good start.
}

func TestNetworkScanner_Scan(t *testing.T) {
	cfg := setupTestConfig()
	scanner := NewNetworkScanner(cfg)

	// We might not want to actually scan network in unit tests as it requires root usually for nmap
	// But we can test initialization
	if scanner == nil {
		t.Fatal("NewNetworkScanner returned nil")
	}

	// Skipping actual Scan execution to avoid CI/environment issues with nmap presence
	/*
		result, err := scanner.Scan("127.0.0.1")
		if err != nil {
			t.Logf("NetworkScanner.Scan() failed (expected if nmap missing): %v", err)
			return
		}
		if result == nil {
			t.Fatal("NetworkScanner.Scan() returned nil result")
		}
	*/
}

func TestConfigScanner_Scan(t *testing.T) {
	cfg := setupTestConfig()
	scanner := NewConfigScanner(cfg)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("ConfigScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("ConfigScanner.Scan() returned nil result")
	}
}
