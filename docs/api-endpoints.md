# ZeroTrace API Endpoints

## Overview
The ZeroTrace API provides RESTful endpoints for managing vulnerability scans, user authentication, company management, and real-time data access. All endpoints follow REST conventions and return JSON responses.

## Base URL
```
Development: http://localhost:8080/api/v1
Production: https://api.zerotrace.com/api/v1
```

## Authentication
All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Response Format
All API responses follow this standard format:
```json
{
  "success": true,
  "data": {},
  "message": "Operation completed successfully",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

Error responses:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": {}
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Authentication Endpoints

### 1. User Registration
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@company.com",
  "password": "securepassword123",
  "name": "John Doe",
  "company_id": "company-123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user-456",
      "email": "user@company.com",
      "name": "John Doe",
      "role": "USER",
      "company_id": "company-123",
      "created_at": "2024-01-15T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "User registered successfully"
}
```

### 2. User Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@company.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "user-456",
      "email": "user@company.com",
      "name": "John Doe",
      "role": "USER",
      "company_id": "company-123"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-16T10:30:00Z"
  },
  "message": "Login successful"
}
```

### 3. Refresh Token
```http
POST /auth/refresh
Authorization: Bearer <refresh_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-16T10:30:00Z"
  },
  "message": "Token refreshed successfully"
}
```

### 4. Logout
```http
POST /auth/logout
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

## Scan Management Endpoints

### 1. Create Scan
```http
POST /scans
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "repository": "https://github.com/company/repo",
  "branch": "main",
  "scan_type": "full",
  "options": {
    "depth": 10,
    "include_patterns": ["*.go", "*.py", "*.js"],
    "exclude_patterns": ["vendor/", "node_modules/"]
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "scan": {
      "id": "scan-789",
      "company_id": "company-123",
      "repository": "https://github.com/company/repo",
      "branch": "main",
      "status": "pending",
      "progress": 0,
      "created_at": "2024-01-15T10:30:00Z",
      "estimated_completion": "2024-01-15T10:35:00Z"
    }
  },
  "message": "Scan created successfully"
}
```

### 2. Get Scans (Paginated)
```http
GET /scans?page=1&limit=20&status=completed&repository=github.com/company/repo
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "scans": [
      {
        "id": "scan-789",
        "company_id": "company-123",
        "repository": "https://github.com/company/repo",
        "branch": "main",
        "status": "completed",
        "progress": 100,
        "start_time": "2024-01-15T10:30:00Z",
        "end_time": "2024-01-15T10:35:00Z",
        "vulnerability_count": 15,
        "critical_count": 2,
        "high_count": 5,
        "medium_count": 8
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  },
  "message": "Scans retrieved successfully"
}
```

### 3. Get Scan Details
```http
GET /scans/{scan_id}
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "scan": {
      "id": "scan-789",
      "company_id": "company-123",
      "agent_id": "agent-001",
      "repository": "https://github.com/company/repo",
      "branch": "main",
      "commit": "abc123def456",
      "status": "completed",
      "progress": 100,
      "start_time": "2024-01-15T10:30:00Z",
      "end_time": "2024-01-15T10:35:00Z",
      "metadata": {
        "files_scanned": 1250,
        "dependencies_analyzed": 45,
        "scan_duration": "5m 30s"
      },
      "vulnerabilities": [
        {
          "id": "vuln-001",
          "type": "dependency",
          "severity": "critical",
          "title": "CVE-2021-1234: Remote Code Execution",
          "description": "A critical vulnerability in package X allows remote code execution",
          "cve_id": "CVE-2021-1234",
          "cvss_score": 9.8,
          "package_name": "vulnerable-package",
          "package_version": "1.2.3",
          "location": "go.mod",
          "remediation": "Update to version 1.2.4 or later"
        }
      ]
    }
  },
  "message": "Scan details retrieved successfully"
}
```

### 4. Update Scan
```http
PUT /scans/{scan_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "status": "cancelled",
  "notes": "Cancelled due to maintenance"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "scan": {
      "id": "scan-789",
      "status": "cancelled",
      "notes": "Cancelled due to maintenance",
      "updated_at": "2024-01-15T10:32:00Z"
    }
  },
  "message": "Scan updated successfully"
}
```

### 5. Delete Scan
```http
DELETE /scans/{scan_id}
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "message": "Scan deleted successfully"
}
```

## Vulnerability Endpoints

### 1. Get Vulnerabilities
```http
GET /vulnerabilities?severity=critical&package_name=vulnerable-package&page=1&limit=50
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "vulnerabilities": [
      {
        "id": "vuln-001",
        "scan_id": "scan-789",
        "type": "dependency",
        "severity": "critical",
        "title": "CVE-2021-1234: Remote Code Execution",
        "description": "A critical vulnerability in package X allows remote code execution",
        "cve_id": "CVE-2021-1234",
        "cvss_score": 9.8,
        "package_name": "vulnerable-package",
        "package_version": "1.2.3",
        "location": "go.mod",
        "remediation": "Update to version 1.2.4 or later",
        "references": [
          "https://nvd.nist.gov/vuln/detail/CVE-2021-1234"
        ],
        "created_at": "2024-01-15T10:35:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 50,
      "total": 150,
      "total_pages": 3
    }
  },
  "message": "Vulnerabilities retrieved successfully"
}
```

### 2. Get Vulnerability Details
```http
GET /vulnerabilities/{vulnerability_id}
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "vulnerability": {
      "id": "vuln-001",
      "scan_id": "scan-789",
      "type": "dependency",
      "severity": "critical",
      "title": "CVE-2021-1234: Remote Code Execution",
      "description": "A critical vulnerability in package X allows remote code execution",
      "cve_id": "CVE-2021-1234",
      "cvss_score": 9.8,
      "cvss_vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H",
      "package_name": "vulnerable-package",
      "package_version": "1.2.3",
      "location": "go.mod",
      "remediation": "Update to version 1.2.4 or later",
      "references": [
        "https://nvd.nist.gov/vuln/detail/CVE-2021-1234"
      ],
      "affected_versions": ["<1.2.4"],
      "patched_versions": ["1.2.4", "1.3.0"],
      "exploit_available": true,
      "exploit_count": 3,
      "created_at": "2024-01-15T10:35:00Z",
      "enrichment": {
        "risk_score": 0.95,
        "trend": "increasing",
        "recommendations": [
          "Update package immediately",
          "Review affected code paths",
          "Monitor for exploitation attempts"
        ]
      }
    }
  },
  "message": "Vulnerability details retrieved successfully"
}
```

### 3. Update Vulnerability Status
```http
PUT /vulnerabilities/{vulnerability_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "status": "acknowledged",
  "notes": "Working on fix",
  "priority": "high"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "vulnerability": {
      "id": "vuln-001",
      "status": "acknowledged",
      "notes": "Working on fix",
      "priority": "high",
      "updated_at": "2024-01-15T10:40:00Z"
    }
  },
  "message": "Vulnerability updated successfully"
}
```

## Report Endpoints

### 1. Generate Report
```http
POST /reports
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "scan_id": "scan-789",
  "format": "pdf",
  "include_charts": true,
  "include_recommendations": true
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "report": {
      "id": "report-456",
      "scan_id": "scan-789",
      "format": "pdf",
      "status": "generating",
      "download_url": null,
      "created_at": "2024-01-15T10:45:00Z",
      "estimated_completion": "2024-01-15T10:47:00Z"
    }
  },
  "message": "Report generation started"
}
```

### 2. Get Report Status
```http
GET /reports/{report_id}
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "report": {
      "id": "report-456",
      "scan_id": "scan-789",
      "format": "pdf",
      "status": "completed",
      "download_url": "https://api.zerotrace.com/reports/report-456.pdf",
      "file_size": "2.5MB",
      "created_at": "2024-01-15T10:45:00Z",
      "completed_at": "2024-01-15T10:47:00Z"
    }
  },
  "message": "Report status retrieved successfully"
}
```

### 3. Download Report
```http
GET /reports/{report_id}/download
Authorization: Bearer <jwt_token>
```

**Response:** Binary file (PDF, CSV, etc.)

## Agent Management Endpoints

### 1. Get Agents
```http
GET /agents?status=active&company_id=company-123
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "agents": [
      {
        "id": "agent-001",
        "company_id": "company-123",
        "name": "Development Agent",
        "status": "active",
        "version": "1.0.0",
        "last_seen": "2024-01-15T10:30:00Z",
        "current_scan": "scan-789",
        "performance": {
          "scans_completed": 45,
          "avg_scan_duration": "5m 30s",
          "success_rate": 98.5
        }
      }
    ]
  },
  "message": "Agents retrieved successfully"
}
```

### 2. Register Agent
```http
POST /agents
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "New Agent",
  "version": "1.0.0",
  "capabilities": ["go", "python", "nodejs"],
  "location": "us-east-1"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "agent": {
      "id": "agent-002",
      "company_id": "company-123",
      "name": "New Agent",
      "status": "active",
      "version": "1.0.0",
      "api_key": "agent-api-key-123",
      "created_at": "2024-01-15T10:50:00Z"
    }
  },
  "message": "Agent registered successfully"
}
```

### 3. Update Agent Status
```http
PUT /agents/{agent_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "status": "maintenance",
  "notes": "Scheduled maintenance"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "agent": {
      "id": "agent-001",
      "status": "maintenance",
      "notes": "Scheduled maintenance",
      "updated_at": "2024-01-15T10:55:00Z"
    }
  },
  "message": "Agent updated successfully"
}
```

## Company Management Endpoints

### 1. Get Company Details
```http
GET /companies/{company_id}
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "company": {
      "id": "company-123",
      "name": "Acme Corporation",
      "domain": "acme.com",
      "settings": {
        "scan_frequency": "daily",
        "notification_preferences": {
          "email": true,
          "slack": false
        },
        "security_policies": {
          "auto_block_critical": true,
          "require_approval": false
        }
      },
      "statistics": {
        "total_scans": 150,
        "total_vulnerabilities": 45,
        "critical_vulnerabilities": 5,
        "last_scan": "2024-01-15T10:30:00Z"
      },
      "created_at": "2024-01-01T00:00:00Z"
    }
  },
  "message": "Company details retrieved successfully"
}
```

### 2. Update Company Settings
```http
PUT /companies/{company_id}
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "settings": {
    "scan_frequency": "weekly",
    "notification_preferences": {
      "email": true,
      "slack": true,
      "slack_webhook": "https://hooks.slack.com/..."
    }
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "company": {
      "id": "company-123",
      "settings": {
        "scan_frequency": "weekly",
        "notification_preferences": {
          "email": true,
          "slack": true,
          "slack_webhook": "https://hooks.slack.com/..."
        }
      },
      "updated_at": "2024-01-15T11:00:00Z"
    }
  },
  "message": "Company settings updated successfully"
}
```

## Dashboard Endpoints

### 1. Get Dashboard Overview
```http
GET /dashboard/overview
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "overview": {
      "total_scans": 150,
      "active_scans": 3,
      "total_vulnerabilities": 45,
      "critical_vulnerabilities": 5,
      "high_vulnerabilities": 12,
      "medium_vulnerabilities": 20,
      "low_vulnerabilities": 8,
      "trends": {
        "vulnerabilities_last_week": 15,
        "vulnerabilities_this_week": 8,
        "trend": "decreasing"
      },
      "recent_activity": [
        {
          "type": "scan_completed",
          "scan_id": "scan-789",
          "repository": "github.com/company/repo",
          "vulnerabilities_found": 3,
          "timestamp": "2024-01-15T10:35:00Z"
        }
      ]
    }
  },
  "message": "Dashboard overview retrieved successfully"
}
```

### 2. Get Vulnerability Trends
```http
GET /dashboard/trends?period=30d&severity=critical
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "trends": {
      "period": "30d",
      "severity": "critical",
      "data": [
        {
          "date": "2024-01-01",
          "count": 2,
          "new": 1,
          "resolved": 0
        },
        {
          "date": "2024-01-02",
          "count": 3,
          "new": 2,
          "resolved": 1
        }
      ],
      "summary": {
        "total_new": 15,
        "total_resolved": 8,
        "net_change": 7
      }
    }
  },
  "message": "Vulnerability trends retrieved successfully"
}
```

## WebSocket Endpoints

### 1. Real-time Updates
```http
GET /ws
Authorization: Bearer <jwt_token>
```

**WebSocket Events:**

#### Scan Status Update
```json
{
  "type": "scan_update",
  "data": {
    "scan_id": "scan-789",
    "status": "scanning",
    "progress": 75,
    "current_file": "src/main.go",
    "timestamp": "2024-01-15T10:32:00Z"
  }
}
```

#### Vulnerability Found
```json
{
  "type": "vulnerability_found",
  "data": {
    "scan_id": "scan-789",
    "vulnerability": {
      "id": "vuln-001",
      "severity": "critical",
      "title": "CVE-2021-1234: Remote Code Execution",
      "package_name": "vulnerable-package"
    },
    "timestamp": "2024-01-15T10:33:00Z"
  }
}
```

#### Agent Status Update
```json
{
  "type": "agent_update",
  "data": {
    "agent_id": "agent-001",
    "status": "idle",
    "current_scan": null,
    "timestamp": "2024-01-15T10:35:00Z"
  }
}
```

## Error Codes

### Common Error Codes
```json
{
  "VALIDATION_ERROR": "Input validation failed",
  "AUTHENTICATION_ERROR": "Invalid or missing authentication",
  "AUTHORIZATION_ERROR": "Insufficient permissions",
  "NOT_FOUND": "Resource not found",
  "RATE_LIMIT_EXCEEDED": "Too many requests",
  "INTERNAL_ERROR": "Internal server error",
  "SERVICE_UNAVAILABLE": "Service temporarily unavailable"
}
```

## Rate Limiting

### Rate Limits
- **Authentication endpoints**: 10 requests per minute
- **Scan endpoints**: 100 requests per minute
- **Vulnerability endpoints**: 200 requests per minute
- **Report endpoints**: 20 requests per minute
- **Dashboard endpoints**: 300 requests per minute

### Rate Limit Headers
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642248600
```

## Pagination

### Pagination Parameters
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)
- `sort`: Sort field (default: created_at)
- `order`: Sort order (asc/desc, default: desc)

### Pagination Response
```json
{
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8,
    "has_next": true,
    "has_prev": false
  }
}
```

## Filtering

### Common Filters
- `status`: Filter by status (pending, scanning, completed, failed)
- `severity`: Filter by severity (critical, high, medium, low)
- `type`: Filter by type (dependency, code, config)
- `repository`: Filter by repository URL
- `date_from`: Filter by start date
- `date_to`: Filter by end date
- `company_id`: Filter by company ID

### Filter Examples
```
GET /scans?status=completed&severity=critical&date_from=2024-01-01
GET /vulnerabilities?type=dependency&package_name=vulnerable-package
GET /agents?status=active&company_id=company-123
```
