# Nuclei Integration - Why It's Essential

## ✅ Nuclei is NOT Removed - It's Very Useful!

Nuclei is **actively used** in the network scanner for vulnerability detection. It was **not removed** - we simplified the implementation to use the CLI instead of the Go library.

## What is Nuclei?

Nuclei is a fast, template-based vulnerability scanner that:
- Scans for **thousands of known vulnerabilities**
- Uses **YAML templates** for vulnerability detection
- Supports **web, network, and cloud** targets
- Has **active community** with regular template updates

## How We Use Nuclei

### In Network Scanner

```go
// agent-go/internal/scanner/network_scanner.go

// Step 4: Run Nuclei vulnerability scanning on discovered hosts
if len(hostsWithOpenPorts) > 0 {
    nucleiFindings, err := ns.nucleiScanner.ScanTargets(targets)
    // Finds CVEs, misconfigurations, exposed services
}
```

### What Nuclei Finds

1. **Web Vulnerabilities**
   - Exposed admin panels
   - Default credentials
   - SQL injection
   - XSS vulnerabilities
   - API security issues

2. **Network Vulnerabilities**
   - Exposed services
   - Weak SSL/TLS
   - Unpatched software
   - Known CVEs

3. **Configuration Issues**
   - Insecure headers
   - Missing security controls
   - Exposed sensitive data

## Why Use CLI Instead of Go Library?

### Current Implementation (CLI)

```go
// Uses Nuclei CLI
cmd := exec.Command("nuclei", args...)
```

**Benefits:**
- ✅ Always uses latest Nuclei version
- ✅ Access to all Nuclei features
- ✅ Automatic template updates
- ✅ Simpler integration
- ✅ Better performance

### Previous Attempt (Go Library)

We tried using the Go library but:
- ❌ Complex API changes
- ❌ Version compatibility issues
- ❌ Limited features
- ❌ Harder to maintain

## Nuclei Workflow

```
Network Scan Discovers Hosts
    ↓
Nmap finds open ports/services
    ↓
Nuclei scans each host:port
    ↓
Finds vulnerabilities (CVEs, misconfigs)
    ↓
Results sent to API
```

## Example Nuclei Findings

```json
{
  "finding_type": "vuln",
  "severity": "high",
  "host": "192.168.1.10",
  "port": 80,
  "description": "Exposed Admin Panel",
  "cve": "CVE-2023-XXXXX",
  "template_id": "exposed-panels/admin-panel"
}
```

## Installing Nuclei

### macOS
```bash
brew install nuclei
```

### Linux
```bash
go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
```

### Windows
```powershell
go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
```

## Nuclei Templates

Nuclei uses templates from:
- **Official templates**: Built-in with Nuclei
- **Community templates**: From GitHub
- **Custom templates**: You can create your own

### Update Templates
```bash
nuclei -update-templates
```

## Benefits of Nuclei

1. **Comprehensive**: Thousands of vulnerability checks
2. **Fast**: Parallel scanning
3. **Updated**: Regular template updates
4. **Flexible**: Custom templates
5. **Proven**: Used by security professionals worldwide

## Integration Points

### Network Scanner
- Scans discovered hosts automatically
- Finds web and network vulnerabilities
- Correlates with device classification

### Authenticated Scanning
- Can use credentials for deeper scans
- Checks authenticated endpoints
- Validates access controls

## Summary

**Nuclei is essential** for:
- ✅ Vulnerability detection
- ✅ CVE identification
- ✅ Security misconfiguration detection
- ✅ Comprehensive network security assessment

**It's NOT removed** - it's a core part of the agentless network scanner!

## See Also

- `agent-go/internal/scanner/nuclei_scanner.go` - Nuclei integration
- `agent-go/internal/scanner/network_scanner.go` - How it's used
- [Nuclei Documentation](https://docs.nuclei.sh/)

