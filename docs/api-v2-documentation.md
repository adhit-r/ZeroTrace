# ZeroTrace API v2 Documentation

## Overview

ZeroTrace API v2 provides comprehensive security scanning and vulnerability management capabilities across multiple security categories. The API supports real-time scanning, compliance monitoring, and advanced analytics.

## Base URL

```
http://localhost:8080/api/v2
```

## Authentication

All API endpoints require authentication via API key or JWT token:

```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     http://localhost:8080/api/v2/vulnerabilities
```

## Endpoints

### Vulnerabilities

#### Get Vulnerabilities
```http
GET /api/v2/vulnerabilities
```

**Query Parameters:**
- `category` (string): Filter by security category (network, compliance, system, auth, database, api, container, ai, iot, privacy, web3)
- `severity` (string): Filter by severity (critical, high, medium, low)
- `status` (string): Filter by status (open, in_progress, closed)
- `agent_id` (string): Filter by agent ID
- `limit` (int): Number of results (default: 50, max: 100)
- `offset` (int): Pagination offset (default: 0)

**Response:**
```json
{
  "vulnerabilities": [
    {
      "id": "vuln-123",
      "agent_id": "agent-456",
      "title": "SQL Injection Vulnerability",
      "description": "Application is vulnerable to SQL injection attacks",
      "severity": "high",
      "category": "database",
      "status": "open",
      "discovered_at": "2024-01-15T10:30:00Z",
      "last_seen": "2024-01-15T10:30:00Z",
      "risk_score": 8.5,
      "exploit_complexity": "low",
      "attack_vector": "network",
      "compliance_frameworks": ["PCI-DSS", "HIPAA"],
      "remediation": "Use parameterized queries",
      "references": ["https://owasp.org/www-community/attacks/SQL_Injection"],
      "tags": ["sql-injection", "database"],
      "metadata": {
        "cve_id": "CVE-2024-1234",
        "cvss_score": 8.5,
        "affected_software": "MySQL 5.7"
      },
      "enrichment_data": {
        "cpe": "cpe:2.3:a:mysql:mysql:5.7:*:*:*:*:*:*:*",
        "confidence": 0.95
      },
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "total": 150,
  "limit": 50,
  "offset": 0
}
```

#### Get Vulnerability Statistics
```http
GET /api/v2/vulnerabilities/stats
```

**Response:**
```json
{
  "total": 150,
  "by_severity": {
    "critical": 5,
    "high": 25,
    "medium": 80,
    "low": 40
  },
  "by_category": {
    "network": 30,
    "compliance": 20,
    "system": 25,
    "auth": 15,
    "database": 20,
    "api": 10,
    "container": 15,
    "ai": 5,
    "iot": 5,
    "privacy": 5,
    "web3": 5
  },
  "risk_score": 7.2,
  "compliance_score": 85.5,
  "last_updated": "2024-01-15T10:30:00Z"
}
```

#### Export Vulnerabilities
```http
GET /api/v2/vulnerabilities/export
```

**Query Parameters:**
- `format` (string): Export format (json, csv, pdf)
- `category` (string): Filter by category
- `severity` (string): Filter by severity

**Response:** File download

### Compliance

#### Get Compliance Status
```http
GET /api/v2/compliance/status
```

**Response:**
```json
{
  "frameworks": {
    "cis": {
      "score": 85.5,
      "total": 100,
      "passed": 85,
      "failed": 15,
      "last_assessment": "2024-01-15T10:30:00Z",
      "next_assessment": "2024-02-15T10:30:00Z"
    },
    "pci_dss": {
      "score": 92.0,
      "total": 50,
      "passed": 46,
      "failed": 4,
      "last_assessment": "2024-01-15T10:30:00Z",
      "next_assessment": "2024-02-15T10:30:00Z"
    },
    "hipaa": {
      "score": 78.0,
      "total": 75,
      "passed": 58,
      "failed": 17,
      "last_assessment": "2024-01-15T10:30:00Z",
      "next_assessment": "2024-02-15T10:30:00Z"
    },
    "gdpr": {
      "score": 88.0,
      "total": 60,
      "passed": 53,
      "failed": 7,
      "last_assessment": "2024-01-15T10:30:00Z",
      "next_assessment": "2024-02-15T10:30:00Z"
    },
    "soc2": {
      "score": 90.0,
      "total": 40,
      "passed": 36,
      "failed": 4,
      "last_assessment": "2024-01-15T10:30:00Z",
      "next_assessment": "2024-02-15T10:30:00Z"
    },
    "iso27001": {
      "score": 82.0,
      "total": 80,
      "passed": 66,
      "failed": 14,
      "last_assessment": "2024-01-15T10:30:00Z",
      "next_assessment": "2024-02-15T10:30:00Z"
    }
  },
  "overall_score": 85.9,
  "last_updated": "2024-01-15T10:30:00Z"
}
```

#### Get Compliance Gaps
```http
GET /api/v2/compliance/gaps
```

**Response:**
```json
{
  "gaps": [
    {
      "framework": "CIS",
      "requirement": "CIS-1.1.1",
      "category": "Access Control",
      "priority": "high",
      "status": "non_compliant",
      "gap": "Default passwords not changed",
      "remediation": "Change all default passwords",
      "effort": "low",
      "timeline": "1 week"
    }
  ],
  "total": 25,
  "critical": 5,
  "high": 10,
  "medium": 7,
  "low": 3
}
```

### Network Scanning

#### Initiate Network Scan
```http
POST /api/v2/network/scan
```

**Request Body:**
```json
{
  "targets": ["192.168.1.0/24", "10.0.0.0/8"],
  "scan_type": "comprehensive",
  "options": {
    "port_range": "1-65535",
    "scan_timeout": 300,
    "parallel_scans": 10
  }
}
```

**Response:**
```json
{
  "scan_id": "scan-123",
  "status": "initiated",
  "estimated_duration": "15 minutes",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### Get Scan Status
```http
GET /api/v2/network/scan/{scan_id}/status
```

**Response:**
```json
{
  "scan_id": "scan-123",
  "status": "running",
  "progress": 65,
  "hosts_scanned": 130,
  "total_hosts": 200,
  "vulnerabilities_found": 15,
  "started_at": "2024-01-15T10:30:00Z",
  "estimated_completion": "2024-01-15T10:45:00Z"
}
```

#### Get Scan Results
```http
GET /api/v2/network/scan/{scan_id}/results
```

**Response:**
```json
{
  "scan_id": "scan-123",
  "status": "completed",
  "hosts": [
    {
      "id": "host-456",
      "ip": "192.168.1.100",
      "hostname": "server-01",
      "os": "Linux Ubuntu 20.04",
      "status": "online",
      "open_ports": [22, 80, 443, 3306],
      "services": [
        {
          "port": 22,
          "protocol": "tcp",
          "service": "ssh",
          "version": "OpenSSH 8.2",
          "ssl_enabled": false
        }
      ],
      "vulnerabilities": [
        {
          "id": "vuln-789",
          "title": "SSH Weak Encryption",
          "severity": "medium",
          "port": 22,
          "service": "ssh"
        }
      ],
      "risk_score": 6.5,
      "last_seen": "2024-01-15T10:30:00Z"
    }
  ],
  "summary": {
    "total_hosts": 200,
    "online_hosts": 150,
    "vulnerabilities_found": 25,
    "critical_vulnerabilities": 2,
    "high_vulnerabilities": 8,
    "medium_vulnerabilities": 10,
    "low_vulnerabilities": 5
  }
}
```

### Agents

#### Get Agent Processing Status
```http
GET /api/v2/agents/processing-status
```

**Response:**
```json
{
  "agents": [
    {
      "id": "agent-123",
      "name": "Production Server",
      "status": "online",
      "processing": {
        "current_app": "Google Chrome",
        "total_apps": 150,
        "processed_apps": 75,
        "vulnerabilities_found": 12,
        "processing_rate": 5.2
      },
      "logs": [
        {
          "timestamp": "2024-01-15T10:30:00Z",
          "app_name": "Google Chrome",
          "cpe": "cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*",
          "vulnerabilities": 3,
          "confidence": 0.95
        }
      ]
    }
  ],
  "total_agents": 5,
  "online_agents": 4,
  "processing_apps": 150,
  "total_vulnerabilities": 25
}
```

## Error Handling

All API endpoints return consistent error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": {
      "field": "severity",
      "issue": "Invalid severity level. Must be one of: critical, high, medium, low"
    }
  },
  "timestamp": "2024-01-15T10:30:00Z",
  "request_id": "req-123"
}
```

### Error Codes

- `VALIDATION_ERROR`: Invalid request parameters
- `AUTHENTICATION_ERROR`: Invalid or missing authentication
- `AUTHORIZATION_ERROR`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `RATE_LIMIT_EXCEEDED`: Too many requests
- `INTERNAL_ERROR`: Server error

## Rate Limiting

API requests are rate limited:
- **Authenticated users**: 1000 requests per hour
- **Unauthenticated users**: 100 requests per hour

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642248600
```

## Webhooks

ZeroTrace supports webhooks for real-time notifications:

### Webhook Events

- `vulnerability.created`: New vulnerability discovered
- `vulnerability.updated`: Vulnerability status changed
- `scan.completed`: Network scan completed
- `compliance.updated`: Compliance score changed
- `agent.offline`: Agent went offline

### Webhook Payload

```json
{
  "event": "vulnerability.created",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "vulnerability_id": "vuln-123",
    "title": "SQL Injection Vulnerability",
    "severity": "high",
    "agent_id": "agent-456"
  }
}
```

## SDKs

### Go SDK

```go
package main

import (
    "github.com/zerotrace/sdk-go"
)

func main() {
    client := zerotrace.NewClient("http://localhost:8080", "your-api-key")
    
    vulnerabilities, err := client.GetVulnerabilities(&zerotrace.VulnerabilityFilter{
        Category: "database",
        Severity: "high",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, vuln := range vulnerabilities {
        fmt.Printf("Vulnerability: %s (Severity: %s)\n", vuln.Title, vuln.Severity)
    }
}
```

### Python SDK

```python
from zerotrace import ZeroTraceClient

client = ZeroTraceClient("http://localhost:8080", "your-api-key")

# Get vulnerabilities
vulnerabilities = client.get_vulnerabilities(
    category="database",
    severity="high"
)

for vuln in vulnerabilities:
    print(f"Vulnerability: {vuln.title} (Severity: {vuln.severity})")
```

### JavaScript SDK

```javascript
const ZeroTraceClient = require('@zerotrace/sdk-js');

const client = new ZeroTraceClient('http://localhost:8080', 'your-api-key');

// Get vulnerabilities
client.getVulnerabilities({
    category: 'database',
    severity: 'high'
}).then(vulnerabilities => {
    vulnerabilities.forEach(vuln => {
        console.log(`Vulnerability: ${vuln.title} (Severity: ${vuln.severity})`);
    });
});
```

## Examples

### Complete Vulnerability Management Workflow

```bash
# 1. Get all high-severity vulnerabilities
curl -H "Authorization: Bearer YOUR_API_KEY" \
     "http://localhost:8080/api/v2/vulnerabilities?severity=high"

# 2. Get compliance status
curl -H "Authorization: Bearer YOUR_API_KEY" \
     "http://localhost:8080/api/v2/compliance/status"

# 3. Initiate network scan
curl -X POST \
     -H "Authorization: Bearer YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"targets": ["192.168.1.0/24"], "scan_type": "comprehensive"}' \
     "http://localhost:8080/api/v2/network/scan"

# 4. Get scan results
curl -H "Authorization: Bearer YOUR_API_KEY" \
     "http://localhost:8080/api/v2/network/scan/scan-123/results"
```

### Real-time Monitoring

```bash
# Monitor agent processing status
curl -H "Authorization: Bearer YOUR_API_KEY" \
     "http://localhost:8080/api/v2/agents/processing-status"

# Get vulnerability statistics
curl -H "Authorization: Bearer YOUR_API_KEY" \
     "http://localhost:8080/api/v2/vulnerabilities/stats"
```

## Changelog

### v2.1.0 (2024-01-15)
- Added network topology visualization
- Enhanced compliance framework support
- Improved vulnerability categorization
- Added real-time processing status

### v2.0.0 (2024-01-01)
- Complete API redesign
- Added comprehensive security categories
- Enhanced compliance monitoring
- Improved performance and scalability

## Support

For API support and questions:
- **Documentation**: https://docs.zerotrace.com/api-v2
- **Support Email**: support@zerotrace.com
- **GitHub Issues**: https://github.com/zerotrace/zerotrace/issues
