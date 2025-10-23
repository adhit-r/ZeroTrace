package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
)

type TrayApp struct {
	apiURL       string
	dashboardURL string
}

func main() {
	app := &TrayApp{
		apiURL:       "http://localhost:8080",
		dashboardURL: "http://localhost:5173",
	}

	systray.Run(app.onReady, app.onExit)
}

func (app *TrayApp) onReady() {
	// Set icon and title
	systray.SetIcon(getIconData())
	systray.SetTitle("ZeroTrace")
	systray.SetTooltip("ZeroTrace Security Monitor")

	// Create menu items
	mDashboard := systray.AddMenuItem("Open Dashboard", "Open ZeroTrace Dashboard")
	mStatus := systray.AddMenuItem("Check Status", "Check Agent Status")
	mVulns := systray.AddMenuItem("View Vulnerabilities", "View Security Vulnerabilities")
	systray.AddSeparator()
	mAbout := systray.AddMenuItem("About", "About ZeroTrace")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit ZeroTrace")

	// Handle menu clicks
	go func() {
		for {
			select {
			case <-mDashboard.ClickedCh:
				app.openDashboard()
			case <-mStatus.ClickedCh:
				app.checkStatus()
			case <-mVulns.ClickedCh:
				app.viewVulnerabilities()
			case <-mAbout.ClickedCh:
				app.showAbout()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (app *TrayApp) onExit() {
	// Cleanup
}

func (app *TrayApp) openDashboard() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", app.dashboardURL)
	case "linux":
		cmd = exec.Command("xdg-open", app.dashboardURL)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", app.dashboardURL)
	default:
		log.Printf("Cannot open browser on %s", runtime.GOOS)
		return
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Failed to open dashboard: %v", err)
	}
}

func (app *TrayApp) checkStatus() {
	resp, err := http.Get(app.apiURL + "/api/agents/")
	if err != nil {
		app.showNotification("ZeroTrace Status", "Error", "Cannot connect to API")
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []struct {
			ID       string `json:"id"`
			Hostname string `json:"hostname"`
			Status   string `json:"status"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		app.showNotification("ZeroTrace Status", "Error", "Failed to parse response")
		return
	}

	if len(result.Data) > 0 {
		agent := result.Data[0]
		app.showNotification("ZeroTrace Status", fmt.Sprintf("Agent: %s", agent.Hostname), fmt.Sprintf("Status: %s", agent.Status))
	} else {
		app.showNotification("ZeroTrace Status", "No Agents", "No agents found")
	}
}

func (app *TrayApp) viewVulnerabilities() {
	resp, err := http.Get(app.apiURL + "/api/vulnerabilities/")
	if err != nil {
		app.showNotification("ZeroTrace Vulnerabilities", "Error", "Cannot fetch vulnerabilities")
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data []interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		app.showNotification("ZeroTrace Vulnerabilities", "Error", "Failed to parse response")
		return
	}

	if len(result.Data) > 0 {
		app.showNotification("ZeroTrace Vulnerabilities", fmt.Sprintf("Found %d vulnerabilities", len(result.Data)), "Check dashboard for details")
	} else {
		app.showNotification("ZeroTrace Vulnerabilities", "No vulnerabilities found", "System is secure")
	}
}

func (app *TrayApp) showAbout() {
	app.showNotification("ZeroTrace Security Monitor",
		"ZeroTrace Agent v1.0.0",
		"Real-time vulnerability monitoring")
}

func (app *TrayApp) showNotification(title, subtitle, message string) {
	// For macOS, use osascript to show notification
	if runtime.GOOS == "darwin" {
		script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s"`, message, title, subtitle)
		exec.Command("osascript", "-e", script).Run()
	} else {
		log.Printf("Notification: %s - %s: %s", title, subtitle, message)
	}
}

func getIconData() []byte {
	// Return a simple icon data (16x16 PNG)
	// This is a minimal PNG icon - in production, you'd use a proper icon file
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x10,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0xF3, 0xFF, 0x61, 0x00, 0x00, 0x00,
		0x19, 0x74, 0x45, 0x58, 0x74, 0x53, 0x6F, 0x66, 0x74, 0x77, 0x61, 0x72,
		0x65, 0x00, 0x41, 0x64, 0x6F, 0x62, 0x65, 0x20, 0x49, 0x6D, 0x61, 0x67,
		0x65, 0x52, 0x65, 0x61, 0x64, 0x79, 0x71, 0xC9, 0x65, 0x3C, 0x00, 0x00,
		0x00, 0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00,
		0x00, 0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00,
		0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}
