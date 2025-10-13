package main

import (
	"fmt"
	"zerotrace/agent/internal/config"
	"zerotrace/agent/internal/scanner"
)

func main() {
	cfg := &config.Config{}
	ss := scanner.NewSystemScanner(cfg)
	info, err := ss.Scan()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("OS: %s %s\n", info.OSName, info.OSVersion)
		fmt.Printf("CPU: %s (%d cores)\n", info.CPUModel, info.CPUCores)
		fmt.Printf("Memory: %.1f GB\n", info.MemoryTotalGB)
		fmt.Printf("Hostname: %s\n", info.Hostname)
	}
}
