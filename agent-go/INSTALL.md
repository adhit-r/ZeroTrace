# ZeroTrace Agent - macOS Installation

## 🍎 **macOS App Installation**

### **Option 1: DMG Installer (Recommended)**

1. **Download the DMG:**
   - `ZeroTrace-Agent-1.0.0.dmg` (4.9 MB)

2. **Install:**
   - Double-click the DMG file
   - Drag "ZeroTrace Agent" to your Applications folder
   - The app will start automatically

3. **First Run:**
   - The app runs in the background (no dock icon)
   - Check Activity Monitor to see it running
   - Look for "zerotrace-agent" process

### **Option 2: Manual Installation**

1. **Build from source:**
   ```bash
   cd agent-go
   go build -o zerotrace-agent cmd/agent/main.go
   ```

2. **Run directly:**
   ```bash
   ./zerotrace-agent
   ```

## 🔧 **Configuration**

The agent uses environment variables for configuration:

```bash
# Copy the example config
cp env.example .env

# Edit the configuration
nano .env
```

### **Required Settings:**

```bash
# Agent Configuration
AGENT_ID=agent-001
AGENT_NAME=ZeroTrace Agent
COMPANY_ID=company-001
API_KEY=your-api-key-here

# API Configuration
API_ENDPOINT=http://localhost:8080
API_TIMEOUT=30s

# Scanning Configuration
SCAN_INTERVAL=5m
SCAN_DEPTH=10
MAX_CONCURRENCY=4
```

## 🚀 **What the Agent Does**

The ZeroTrace Agent is a **software vulnerability scanner** that:

1. **Scans Installed Applications:**
   - macOS Applications (`/Applications`, `~/Applications`)
   - Homebrew packages (`brew list`)
   - System applications

2. **Detects Software:**
   - Browsers (Chrome, Firefox, Safari, Edge)
   - PDF readers (Adobe Acrobat)
   - Media players (VLC)
   - Development tools (VS Code, IntelliJ)
   - System utilities (7-Zip, Notepad++)

3. **Sends Data to Server:**
   - App names and versions
   - Installation dates
   - File sizes and paths
   - Vendor information

4. **Runs Continuously:**
   - Scans every 5 minutes (configurable)
   - Runs in background
   - No user interaction required

## 🔍 **Architecture: Agent vs Server**

**Agent Responsibilities:**
- ✅ Discovers installed software
- ✅ Collects app metadata
- ✅ Sends data to server
- ✅ Runs lightweight and fast

**Server Responsibilities:**
- ✅ Matches apps against CVE database
- ✅ Determines vulnerability status
- ✅ Provides remediation steps
- ✅ Centralized vulnerability intelligence

## 📊 **Data Flow**

```
Your Mac → Agent → API Server → Vulnerability Database
   ↓           ↓         ↓              ↓
Apps List → App Data → CVE Match → Vulnerability Report
```

## 🛡️ **Privacy & Security**

- **Local Discovery:** Agent only scans your local system
- **No Code Scanning:** Only installed applications, not source code
- **Configurable:** You control what gets scanned
- **Secure:** Data sent over HTTPS to your server

## 🎯 **Use Cases**

Perfect for:
- **Enterprise Security:** Monitor software vulnerabilities across company devices
- **Compliance:** Track software versions for security audits
- **IT Management:** Know what software is installed on endpoints
- **Security Teams:** Identify vulnerable software quickly

## 🔧 **Troubleshooting**

### **Agent Not Running:**
```bash
# Check if process is running
ps aux | grep zerotrace-agent

# Check logs
tail -f /var/log/system.log | grep zerotrace
```

### **Configuration Issues:**
```bash
# Verify environment variables
env | grep ZERO
```

### **Network Issues:**
```bash
# Test API connectivity
curl http://localhost:8080/health
```

## 📱 **macOS App Features**

- **Background Operation:** Runs without dock icon
- **Auto-start:** Can be configured to start on login
- **System Integration:** Uses macOS security features
- **Easy Updates:** Simple drag-and-drop installation

---

**🎉 Your ZeroTrace Agent is now running and scanning for software vulnerabilities!**
