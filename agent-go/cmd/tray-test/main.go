package main

import (
	"log"

	"zerotrace/agent/internal/tray"
)

func main() {
	log.Println("Starting ZeroTrace Tray Test...")

	// Create and start tray manager
	trayManager := tray.NewTrayManager()
	trayManager.Start()

	log.Println("Tray icon should now be visible in the menu bar")
	log.Println("Press Ctrl+C to exit")

	// Keep the program running
	select {}
}
