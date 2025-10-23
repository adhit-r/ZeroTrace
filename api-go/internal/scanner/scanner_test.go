package scanner

import (
	"testing"
	"time"
)

func TestSystemScanner_Scan(t *testing.T) {
	scanner := NewSystemScanner()

	// Test basic scan functionality
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("SystemScanner.Scan() failed: %v", err)
	}

	// Verify result structure
	if result == nil {
		t.Fatal("SystemScanner.Scan() returned nil result")
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
	scanner := NewSoftwareScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("SoftwareScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("SoftwareScanner.Scan() returned nil result")
	}

	// Verify applications are found
	if len(result.Applications) == 0 {
		t.Log("Warning: No applications found - this might be expected in test environment")
	}
}

func TestNetworkScanner_ScanPort(t *testing.T) {
	scanner := NewNetworkScanner()

	// Test scanning localhost
	result, err := scanner.ScanPort("127.0.0.1", 22)
	if err != nil {
		t.Logf("NetworkScanner.ScanPort() failed (expected in some environments): %v", err)
		return
	}

	if result == nil {
		t.Fatal("NetworkScanner.ScanPort() returned nil result")
	}
}

func TestConfigScanner_Scan(t *testing.T) {
	scanner := NewConfigScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("ConfigScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("ConfigScanner.Scan() returned nil result")
	}
}

func TestSystemVulnerabilityScanner_Scan(t *testing.T) {
	scanner := NewSystemVulnerabilityScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("SystemVulnerabilityScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("SystemVulnerabilityScanner.Scan() returned nil result")
	}
}

func TestAuthScanner_Scan(t *testing.T) {
	scanner := NewAuthScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("AuthScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("AuthScanner.Scan() returned nil result")
	}
}

func TestDatabaseScanner_Scan(t *testing.T) {
	scanner := NewDatabaseScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("DatabaseScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("DatabaseScanner.Scan() returned nil result")
	}
}

func TestAPIScanner_Scan(t *testing.T) {
	scanner := NewAPIScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("APIScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("APIScanner.Scan() returned nil result")
	}
}

func TestContainerScanner_Scan(t *testing.T) {
	scanner := NewContainerScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("ContainerScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("ContainerScanner.Scan() returned nil result")
	}
}

func TestAIMLScanner_Scan(t *testing.T) {
	scanner := NewAIMLScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("AIMLScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("AIMLScanner.Scan() returned nil result")
	}
}

func TestIoTOTScanner_Scan(t *testing.T) {
	scanner := NewIoTOTScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("IoTOTScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("IoTOTScanner.Scan() returned nil result")
	}
}

func TestPrivacyScanner_Scan(t *testing.T) {
	scanner := NewPrivacyScanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("PrivacyScanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("PrivacyScanner.Scan() returned nil result")
	}
}

func TestWeb3Scanner_Scan(t *testing.T) {
	scanner := NewWeb3Scanner()

	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Web3Scanner.Scan() failed: %v", err)
	}

	if result == nil {
		t.Fatal("Web3Scanner.Scan() returned nil result")
	}
}

// Benchmark tests for performance
func BenchmarkSystemScanner_Scan(b *testing.B) {
	scanner := NewSystemScanner()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.Scan()
		if err != nil {
			b.Fatalf("SystemScanner.Scan() failed: %v", err)
		}
	}
}

func BenchmarkSoftwareScanner_Scan(b *testing.B) {
	scanner := NewSoftwareScanner()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := scanner.Scan()
		if err != nil {
			b.Fatalf("SoftwareScanner.Scan() failed: %v", err)
		}
	}
}

func BenchmarkNetworkScanner_ScanPort(b *testing.B) {
	scanner := NewNetworkScanner()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = scanner.ScanPort("127.0.0.1", 22)
	}
}

// Integration tests
func TestScannerIntegration(t *testing.T) {
	scanners := []Scanner{
		NewSystemScanner(),
		NewSoftwareScanner(),
		NewConfigScanner(),
		NewSystemVulnerabilityScanner(),
		NewAuthScanner(),
		NewDatabaseScanner(),
		NewAPIScanner(),
		NewContainerScanner(),
		NewAIMLScanner(),
		NewIoTOTScanner(),
		NewPrivacyScanner(),
		NewWeb3Scanner(),
	}

	for _, scanner := range scanners {
		t.Run(scanner.GetName(), func(t *testing.T) {
			result, err := scanner.Scan()
			if err != nil {
				t.Errorf("Scanner %s failed: %v", scanner.GetName(), err)
				return
			}

			if result == nil {
				t.Errorf("Scanner %s returned nil result", scanner.GetName())
			}
		})
	}
}

// Test concurrent scanning
func TestConcurrentScanning(t *testing.T) {
	scanners := []Scanner{
		NewSystemScanner(),
		NewSoftwareScanner(),
		NewConfigScanner(),
	}

	results := make(chan ScanResult, len(scanners))
	errors := make(chan error, len(scanners))

	// Start all scanners concurrently
	for _, scanner := range scanners {
		go func(s Scanner) {
			result, err := s.Scan()
			if err != nil {
				errors <- err
				return
			}
			results <- result
		}(scanner)
	}

	// Collect results
	completed := 0
	timeout := time.After(60 * time.Second)

	for completed < len(scanners) {
		select {
		case result := <-results:
			if result == nil {
				t.Error("Received nil result from concurrent scan")
			}
			completed++
		case err := <-errors:
			t.Errorf("Concurrent scan failed: %v", err)
			completed++
		case <-timeout:
			t.Fatal("Concurrent scanning timed out")
		}
	}
}
