#!/usr/bin/env python3
"""
ZeroTrace macOS Tray Icon
Creates a proper system tray icon for macOS using rumps
"""

import rumps
import webbrowser
import subprocess
import json
from datetime import datetime
import urllib.request
import urllib.parse

class ZeroTraceTray(rumps.App):
    def __init__(self):
        super(ZeroTraceTray, self).__init__("ZeroTrace", icon="🔒", template=True)
        self.api_url = "http://localhost:8080"
        self.dashboard_url = "http://localhost:5173"
        
        # Create menu items
        self.menu = [
            "ZeroTrace Security Monitor",
            None,
            rumps.MenuItem("Open Dashboard", callback=self.open_dashboard),
            rumps.MenuItem("Check Status", callback=self.check_status),
            rumps.MenuItem("View Vulnerabilities", callback=self.view_vulnerabilities),
            None,
            rumps.MenuItem("About", callback=self.show_about),
            rumps.MenuItem("Quit", callback=self.quit_app)
        ]
    
    def open_dashboard(self, _):
        """Open the ZeroTrace dashboard in browser"""
        webbrowser.open(self.dashboard_url)
        rumps.notification("ZeroTrace", "Dashboard", "Opening ZeroTrace Dashboard...")
    
    def check_status(self, _):
        """Check agent status"""
        try:
            with urllib.request.urlopen(f"{self.api_url}/api/agents/", timeout=5) as response:
                if response.status == 200:
                    data = json.loads(response.read().decode())
                    agents = data.get('data', [])
                    if agents:
                        agent = agents[0]
                        status = agent.get('status', 'Unknown')
                        hostname = agent.get('hostname', 'Unknown')
                        rumps.notification("ZeroTrace Status", f"Agent: {hostname}", f"Status: {status}")
                    else:
                        rumps.notification("ZeroTrace Status", "No Agents", "No agents found")
                else:
                    rumps.notification("ZeroTrace Status", "Error", "Cannot connect to API")
        except Exception as e:
            rumps.notification("ZeroTrace Status", "Error", f"Connection failed: {str(e)}")
    
    def view_vulnerabilities(self, _):
        """View vulnerabilities"""
        try:
            with urllib.request.urlopen(f"{self.api_url}/api/vulnerabilities/", timeout=5) as response:
                if response.status == 200:
                    data = json.loads(response.read().decode())
                    vulns = data.get('data', [])
                    if vulns:
                        rumps.notification("ZeroTrace Vulnerabilities", f"Found {len(vulns)} vulnerabilities", "Check dashboard for details")
                    else:
                        rumps.notification("ZeroTrace Vulnerabilities", "No vulnerabilities found", "System is secure")
                else:
                    rumps.notification("ZeroTrace Vulnerabilities", "Error", "Cannot fetch vulnerabilities")
        except Exception as e:
            rumps.notification("ZeroTrace Vulnerabilities", "Error", f"Connection failed: {str(e)}")
    
    def show_about(self, _):
        """Show about information"""
        rumps.alert("ZeroTrace Security Monitor", 
                   "ZeroTrace Agent v1.0.0\n"
                   "Real-time vulnerability monitoring\n"
                   "Dashboard: http://localhost:5173\n"
                   "API: http://localhost:8080")
    
    def quit_app(self, _):
        """Quit the application"""
        rumps.quit_application()

if __name__ == "__main__":
    app = ZeroTraceTray()
    app.run()
