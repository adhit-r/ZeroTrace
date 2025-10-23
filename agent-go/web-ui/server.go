package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Get the directory where the server is running
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Serve static files from the current directory
	http.Handle("/", http.FileServer(http.Dir(dir)))

	port := "3002"
	if envPort := os.Getenv("AGENT_UI_PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("🛡️ ZeroTrace Agent UI starting on http://localhost:%s\n", port)
	fmt.Printf("📊 Agent status and monitoring interface\n")
	fmt.Printf("⏹️  Press Ctrl+C to stop\n")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
