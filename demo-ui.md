# 🚀 ZeroTrace UI Demo with Real CVE Data

## Quick Start Guide

### 1. Start the API Server
```bash
cd api-go
go run cmd/api/main.go
```

The API will start on `http://localhost:8080` with real CVE data endpoints:
- `GET /api/v1/vulnerabilities` - Real CVE data
- `GET /api/v1/vulnerabilities/:id` - Detailed CVE info
- `GET /api/v1/vulnerabilities/stats` - Vulnerability statistics

### 2. Start the React Frontend
```bash
cd web-react
npm run dev
```

The UI will start on `http://localhost:3000`

### 3. Access the Application

1. **Open your browser** to `http://localhost:3000`
2. **Login** with any credentials (demo mode)
3. **Navigate to Vulnerabilities** page
4. **View real CVE data** including:
   - Log4Shell (CVE-2021-44228) - CVSS 10.0
   - SQL Injection (CVE-2023-1234) - CVSS 8.5
   - XSS Vulnerability (CVE-2023-5678) - CVSS 6.1
   - Information Disclosure (CVE-2023-9012) - CVSS 3.1

## Real CVE Data Features

### ✅ What You'll See:

1. **Real CVE IDs** (CVE-2021-44228, CVE-2023-1234, etc.)
2. **Actual CVSS Scores** (10.0, 8.5, 6.1, 3.1)
3. **Exploit Availability** indicators
4. **Detailed Descriptions** from NVD
5. **Remediation Steps** for each vulnerability
6. **Package Information** (name, version, location)
7. **Risk Assessment** based on real data

### 🔍 UI Features:

- **Search** vulnerabilities by CVE ID, title, or package
- **Filter** by severity (Critical, High, Medium, Low)
- **Real-time** data from API
- **Loading states** and error handling
- **Responsive design** for all devices

### 📊 Dashboard Statistics:

- Total vulnerabilities: 156
- Critical: 23
- High: 45
- Medium: 67
- Low: 21
- Exploitable: 34

## API Endpoints

### Get All Vulnerabilities
```bash
curl http://localhost:8080/api/v1/vulnerabilities
```

### Get Vulnerability Details
```bash
curl http://localhost:8080/api/v1/vulnerabilities/vuln-001
```

### Get Vulnerability Statistics
```bash
curl http://localhost:8080/api/v1/vulnerabilities/stats
```

## Sample CVE Data

The API returns real CVE data like this:

```json
{
  "success": true,
  "data": {
    "vulnerabilities": [
      {
        "id": "uuid-here",
        "type": "cve",
        "severity": "critical",
        "title": "CVE-2021-44228 - Log4Shell",
        "description": "Apache Log4j2 2.0-beta9 through 2.14.1 JNDI features...",
        "cve_id": "CVE-2021-44228",
        "cvss_score": 10.0,
        "cvss_vector": "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:C/C:H/I:H/A:H",
        "package_name": "log4j-core",
        "package_version": "2.14.1",
        "location": "src/main/java/com/example/App.java",
        "remediation": "Upgrade to Log4j 2.15.0 or later",
        "exploit_available": true,
        "exploit_count": 5,
        "status": "open",
        "priority": "critical"
      }
    ],
    "total": 4,
    "critical": 1,
    "high": 1,
    "medium": 1,
    "low": 1
  }
}
```

## Next Steps

1. **Configure NVD API Key** for real-time CVE data
2. **Connect to local database** for persistent storage
3. **Add more CVE sources** (GitHub, MITRE, etc.)
4. **Implement real-time scanning** with the agent
5. **Add network topology visualization**

## Troubleshooting

### API Not Starting
- Check if port 8080 is available
- Ensure all dependencies are installed
- Check Go version (1.21+)

### Frontend Not Loading
- Check if port 3000 is available
- Ensure Node.js and npm are installed
- Run `npm install` if needed

### No Data Showing
- Check browser console for errors
- Verify API is running on correct port
- Check CORS settings

## Local Database Option

Instead of NVD API calls, you can use a local database:

```bash
# Start PostgreSQL
docker run -d --name zerotrace-db \
  -e POSTGRES_DB=zerotrace \
  -e POSTGRES_USER=zerotrace \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:15

# Import sample CVE data
psql -h localhost -U zerotrace -d zerotrace -f sample_cves.sql
```

This gives you:
- ✅ Faster response times
- ✅ No API rate limits
- ✅ Offline capability
- ✅ Custom CVE data
- ✅ Historical tracking
