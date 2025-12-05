# Quick Test Guide - Network Scanning

## ✅ Setup Complete!

Dependencies installed:
- ✅ Nmap 7.98
- ✅ Nuclei v3.5.1
- ✅ Agent built successfully

## Quick Start

### 1. Configure Environment

```bash
cd agent-go
cp env.example .env
```

Edit `.env` and set:
```bash
NETWORK_SCAN_ENABLED=true
NETWORK_SCAN_INTERVAL=10m  # Test with shorter interval
API_ENDPOINT=http://localhost:8080
```

### 2. Run the Agent

```bash
# Option 1: Run the built binary
./zerotrace-agent

# Option 2: Run directly with go
go run cmd/agent/main.go
```

### 3. What Happens

1. **30 seconds after start**: First network scan begins
2. **Nmap discovery**: Scans local network for devices
3. **Device classification**: Identifies switches, routers, IoT, phones, servers
4. **Config auditing**: Checks for security misconfigurations
5. **Nuclei scanning**: Runs vulnerability scans on discovered hosts
6. **Results sent**: Sends findings to API (if running)

### 4. Monitor Logs

Watch for:
```
Starting network scan...
Network scan completed: X findings on Y hosts
Successfully sent network scan results to API.
```

## Test Without API

If API is not running, the agent will:
- Still perform scans ✅
- Log results to console
- Retry sending when API is available

## Manual Testing

### Test Nmap Directly
```bash
nmap -sn 192.168.1.0/24  # Discover hosts
nmap -sV 192.168.1.1     # Service detection
```

### Test Nuclei Directly
```bash
nuclei -u http://localhost -json -silent
```

## Troubleshooting

**No devices found?**
- Check your network: `ifconfig` or `ip addr`
- Try scanning specific subnet in code
- Check firewall settings

**Build errors?**
```bash
go mod tidy
go clean -cache
go build -o zerotrace-agent cmd/agent/main.go
```

**Permission errors?**
- Nmap may need sudo for some scans
- Try: `sudo ./zerotrace-agent`

## Expected Output

You should see:
- Device discoveries with types (switch/router/IoT/phone/server)
- Open ports and services
- Configuration findings (insecure protocols, etc.)
- Vulnerability findings from Nuclei

All findings are structured and ready for API submission!

