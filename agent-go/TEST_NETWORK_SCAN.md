# Testing Network Scanning Locally

This guide helps you test the network scanning feature locally.

## Prerequisites

1. **Nmap** - Installed via `brew install nmap` ✅
2. **Nuclei** - Installed via `brew install nuclei` ✅
3. **Go dependencies** - Run `go mod tidy`

## Quick Test

### 1. Setup Environment

```bash
cd agent-go
cp env.example .env
```

Edit `.env` and ensure:
```bash
NETWORK_SCAN_ENABLED=true
NETWORK_SCAN_INTERVAL=10m  # Shorter interval for testing
API_ENDPOINT=http://localhost:8080  # Your API endpoint
```

### 2. Build and Run

```bash
# Build the agent
go build -o zerotrace-agent cmd/agent/main.go

# Run the agent
./zerotrace-agent
```

Or run directly:
```bash
go run cmd/agent/main.go
```

### 3. What to Expect

The agent will:
1. **Start scanning** after 30 seconds (initial delay)
2. **Discover devices** on your local network using Nmap
3. **Classify devices** (switches, routers, IoT, phones, servers)
4. **Detect configuration errors** (insecure protocols, default credentials, etc.)
5. **Run vulnerability scans** using Nuclei on discovered hosts
6. **Send results** to the API endpoint

### 4. Monitor Output

Look for log messages like:
```
Starting network scan...
Network scan completed: X findings on Y hosts
Successfully sent network scan results to API.
```

### 5. Test Without API

If you don't have the API running, the agent will still scan but fail to send results. You can:
- Check the logs for scan results
- Modify the code to print findings to console
- Run the API server separately

## Testing Specific Features

### Test Device Classification

The scanner identifies:
- **Switches**: Devices with SNMP (port 161) and switch-specific services
- **Routers**: Devices with routing protocol ports (BGP 179, OSPF 88, RIP 520)
- **IoT Devices**: Devices with MQTT (1883, 8883), CoAP (5683, 5684)
- **Phones**: Mobile devices with Android/iOS OS detection
- **Servers**: Devices with common server ports (22, 80, 443, 3306, etc.)

### Test Configuration Auditing

The scanner checks for:
- Default credentials (admin/admin, root/password, etc.)
- Insecure protocols (Telnet, FTP, SNMP v1/v2)
- Open management interfaces without authentication
- Unnecessary open ports
- Weak SSL/TLS configurations

### Test Vulnerability Scanning

Nuclei will scan discovered hosts for:
- Web vulnerabilities
- Network vulnerabilities
- Service misconfigurations
- Known CVEs

## Troubleshooting

### Nmap not found
```bash
brew install nmap
```

### Nuclei not found
```bash
brew install nuclei
```

### Build errors
```bash
go mod tidy
go clean -cache
go build -o zerotrace-agent cmd/agent/main.go
```

### Network scan not running
- Check `NETWORK_SCAN_ENABLED=true` in `.env`
- Check logs for errors
- Verify you have network permissions (may need sudo on some systems)

### No devices found
- Ensure you're on a network with other devices
- Check firewall settings
- Try scanning a specific subnet: modify code to scan `192.168.1.0/24`

## Manual Testing

You can also test components individually:

### Test Nmap Integration
```bash
nmap -sn 192.168.1.0/24  # Ping scan
nmap -sV 192.168.1.1     # Service version detection
```

### Test Nuclei
```bash
nuclei -u http://localhost -json -silent
```

### Test Device Classification
The device classifier uses:
- Port patterns
- Service banners
- OS detection results
- Protocol signatures

## Expected Results

A successful scan should produce:
- **Port findings**: Open ports with service information
- **Device classifications**: Device types for each host
- **Configuration findings**: Security misconfigurations
- **Vulnerability findings**: CVEs and security issues from Nuclei

All findings are sent to the API endpoint `/api/agents/network-scan-results`.

## Next Steps

1. Review scan results in API logs
2. Check device classifications
3. Verify configuration error detection
4. Test with credentials for authenticated scanning
5. Monitor scan performance and adjust intervals

