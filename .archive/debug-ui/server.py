#!/usr/bin/env python3
"""
ZeroTrace Debug UI Server
A lightweight HTTP server for the debug monitoring interface
"""

import http.server
import socketserver
import webbrowser
import threading
import time
from pathlib import Path

class DebugServer:
    def __init__(self, port=3001):
        self.port = port
        self.handler = http.server.SimpleHTTPRequestHandler
        self.httpd = None
        
    def start(self):
        """Start the debug server"""
        try:
            self.httpd = socketserver.TCPServer(("", self.port), self.handler)
            print(f"🔍 ZeroTrace Debug Monitor starting on http://localhost:{self.port}")
            print(f"📊 Monitoring agent data processing pipeline...")
            print(f"🔄 Auto-refresh enabled for real-time updates")
            print(f"⏹️  Press Ctrl+C to stop")
            
            # Open browser automatically
            threading.Timer(1.0, lambda: webbrowser.open(f'http://localhost:{self.port}')).start()
            
            self.httpd.serve_forever()
        except KeyboardInterrupt:
            print(f"\n🛑 Debug monitor stopped")
            self.stop()
        except Exception as e:
            print(f"❌ Error starting debug server: {e}")
            
    def stop(self):
        """Stop the debug server"""
        if self.httpd:
            self.httpd.shutdown()
            self.httpd.server_close()

if __name__ == "__main__":
    # Change to the debug-ui directory
    import os
    os.chdir(Path(__file__).parent)
    
    server = DebugServer(port=3001)
    server.start()
