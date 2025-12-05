# ZeroTrace Enrichment Service

Python FastAPI service for CVE enrichment and batch processing with ultra-optimized performance.

## Overview

The ZeroTrace Enrichment Service provides comprehensive vulnerability data enrichment by:

- Fetching CVE data from multiple sources (NVD, MITRE, ExploitDB)
- Batch processing for high-performance enrichment
- AI-powered vulnerability analysis
- Ultra-optimized algorithms for 10,000x performance improvement
- Real-time CVE matching and scoring

## Features

### CVE Enrichment
- Multi-source CVE data fetching (NVD, MITRE, ExploitDB)
- Local CVE database caching
- Real-time CVE matching
- CVSS scoring and severity classification
- **Hybrid CPE Matching**: Combines official NVD CPE dictionary (via CPE Guesser) with ML-based matching for maximum accuracy

### Batch Processing
- Ultra-optimized batch processing (500+ apps per batch)
- Parallel processing (1000+ concurrent requests)
- Smart caching (Memory + Redis + Memcached)
- Rate limiting and retry logic

### AI Services
- AI-powered vulnerability analysis
- Exploit intelligence gathering
- Predictive vulnerability modeling
- Remediation plan generation

## Quick Start

### Prerequisites

- Python 3.9+
- uv (recommended) or pip
- Redis/Valkey (optional, for caching and CPE Guesser)
- PostgreSQL (optional, for data storage)
- CPE Guesser service (optional, for enhanced CPE matching - included in docker-compose)

### Installation

```bash
# Clone repository
git clone https://github.com/adhit-r/ZeroTrace.git
cd ZeroTrace/enrichment-python

# Install dependencies with uv (recommended)
uv pip install -r requirements.txt

# Or with pip
pip install -r requirements.txt
```

### Configuration

```bash
# Copy environment template
cp env.example .env

# Edit .env file with your configuration
# Required: NVD_API_KEY (for NVD API access)
```

### Running the Service

```bash
# Run with uv (recommended)
uv run python app/main.py

# Or with uvicorn
uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload

# Or with python
python -m uvicorn app.main:app --host 0.0.0.0 --port 8000
```

## Configuration

### Environment Variables

See `env.example` for all available configuration options.

#### Required

- `ENRICHMENT_PORT`: Service port (default: 8000)
- `NVD_API_KEY`: NVD API key for CVE data (recommended)

#### Optional

- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string for caching
- `OPENAI_API_KEY`: OpenAI API key for AI services
- `ANTHROPIC_API_KEY`: Anthropic API key for AI services
- `LOG_LEVEL`: Logging level (default: info)
- `CACHE_TTL`: Cache TTL in seconds (default: 3600)
- `CPE_GUESSER_ENABLED`: Enable CPE Guesser integration (default: true)
- `CPE_GUESSER_URL`: CPE Guesser service URL (default: http://localhost:8001)
- `CPE_GUESSER_FALLBACK_TO_ML`: Fallback to ML matcher if CPE Guesser unavailable (default: true)

### Configuration Files

- `env.example`: Environment variable template
- `requirements.txt`: Python dependencies
- `cve_data.json`: Local CVE database cache

## API Endpoints

### Health Check

```bash
curl http://localhost:8000/health
```

**Example Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

### Enrich Software

```bash
curl -X POST http://localhost:8000/enrich/software \
  -H "Content-Type: application/json" \
  -d '[
    {
      "name": "nginx",
      "version": "1.21.0",
      "vendor": "nginx"
    }
  ]'
```

**Example Response:**
```json
[
  {
    "name": "nginx",
    "version": "1.21.0",
    "vendor": "nginx",
    "cves": [
      {
        "cve_id": "CVE-2021-23017",
        "severity": "high",
        "cvss_score": 7.5,
        "description": "Security vulnerability in nginx"
      }
    ],
    "enrichment_status": "completed"
  }
]
```

### Batch Enrichment

```bash
curl -X POST http://localhost:8000/enrich/batch \
  -H "Content-Type: application/json" \
  -d '[
    {
      "name": "nginx",
      "version": "1.21.0"
    },
    {
      "name": "python",
      "version": "3.9.0"
    }
  ]'
```

**Example Response:**
```json
{
  "job_id": "job-12345",
  "status": "processing",
  "total_items": 2,
  "processed": 0,
  "estimated_completion": "2025-01-15T10:35:00Z"
}
```

### Enrichment Status

```bash
curl http://localhost:8000/enrich/status/job-12345
```

**Example Response:**
```json
{
  "job_id": "job-12345",
  "status": "completed",
  "total_items": 2,
  "processed": 2,
  "results": [
    {
      "name": "nginx",
      "version": "1.21.0",
      "cves": [...],
      "enrichment_status": "completed"
    }
  ]
}
```

### Python Client Example

```python
import requests

# Enrich single software
response = requests.post(
    "http://localhost:8000/enrich/software",
    json=[{
        "name": "nginx",
        "version": "1.21.0",
        "vendor": "nginx"
    }]
)
data = response.json()
print(f"Found {len(data[0]['cves'])} CVEs")

# Batch enrichment
batch_response = requests.post(
    "http://localhost:8000/enrich/batch",
    json=[
        {"name": "nginx", "version": "1.21.0"},
        {"name": "python", "version": "3.9.0"}
    ]
)
job_id = batch_response.json()["job_id"]

# Check status
status_response = requests.get(
    f"http://localhost:8000/enrich/status/{job_id}"
)
print(status_response.json())
```

## Development

### Project Structure

```
enrichment-python/
├── app/
│   ├── main.py                    # FastAPI application
│   ├── cve_enrichment.py          # CVE enrichment service
│   ├── batch_enrichment.py        # Batch processing service
│   ├── ultra_optimized_enrichment.py  # Ultra-optimized service
│   ├── ai_services/               # AI-powered services
│   │   ├── ai_service.py
│   │   ├── exploit_intelligence.py
│   │   └── ...
│   ├── ai_matching/               # AI-based CVE matching
│   ├── analytics/                 # Analytics services
│   ├── compliance/                # Compliance services
│   └── integrations/              # Third-party integrations
├── scripts/
│   ├── update_cve_data.py         # CVE database updater
│   ├── cron_update_cve.py        # Scheduled CVE updates
│   └── ...
├── tests/
│   ├── test_cve_enrichment.py
│   └── conftest.py
├── requirements.txt
├── env.example
└── README.md
```

### Running Tests

```bash
# Run all tests
pytest tests/ -v

# Run specific test
pytest tests/test_cve_enrichment.py -v

# Run with coverage
pytest tests/ -v --cov=app --cov-report=html
```

### Updating CVE Data

```bash
# Manual update
python scripts/update_cve_data.py

# Scheduled update (systemd)
# Install cve-update.service and cve-update.timer
sudo systemctl enable cve-update.timer
sudo systemctl start cve-update.timer
```

### Data Files

**CPE Dictionary** (`data/nvdcpe-2.0.tar`, ~682MB):
- Required for seeding CPE database
- Auto-downloaded by `cpe-guesser/bin/import.py --download` if missing
- Kept locally to avoid re-downloading (in .gitignore)

**CVE Data** (`cve_data.json`, ~9.8MB):
- Local CVE database for fast lookups
- Generated by `scripts/update_cve_data.py`
- Can be updated periodically

See `DATA_FILES.md` for details.

## Performance

### Optimizations

- **Ultra-optimized algorithms**: 10,000x performance improvement
- **Multi-level caching**: Memory + Redis + Memcached
- **Parallel processing**: 1000+ concurrent requests
- **Batch processing**: 500+ apps per batch
- **Connection pooling**: Optimized HTTP connections

### Performance Metrics

- **Enrichment Processing**: < 30ms per application
- **Batch Processing**: 500+ apps per second
- **Cache Hit Rate**: > 80%
- **API Response Time**: < 100ms (95th percentile)

## Deployment

### Docker

```bash
# Build image
docker build -t zerotrace-enrichment .

# Run container
docker run -p 8000:8000 --env-file .env zerotrace-enrichment
```

### Docker Compose

```bash
# Start with docker-compose (includes CPE Guesser)
docker-compose up -d

# CPE Guesser will automatically seed on first run
# This may take 5-15 minutes depending on system
```

**Note**: The CPE Guesser service is included in docker-compose and will automatically:
- Download the NVD CPE dictionary on first run
- Seed the Valkey database
- Start the service once ready

### Production

```bash
# Run with gunicorn
gunicorn app.main:app -w 4 -k uvicorn.workers.UvicornWorker --bind 0.0.0.0:8000
```

## CPE Guesser Integration

The enrichment service uses a hybrid CPE matcher that combines:
1. **CPE Guesser** (Primary) - Official NVD CPE dictionary via fast Valkey lookups
2. **ML Matcher** (Fallback) - TF-IDF similarity matching for edge cases

### Setup

CPE Guesser is automatically included in docker-compose. For manual setup:

```bash
# Start CPE Guesser service
cd cpe-guesser
python3 bin/import.py --download  # Seed database (first time only)
python3 bin/server.py             # Start server
```

### Configuration

Configure CPE Guesser integration via environment variables:

```bash
CPE_GUESSER_ENABLED=true
CPE_GUESSER_URL=http://localhost:8001
CPE_GUESSER_FALLBACK_TO_ML=true
```

### Benefits

- **Accuracy**: Official NVD CPE dictionary vs. CVE-derived data
- **Performance**: Sub-10ms lookups vs. 100ms+ for ML matching
- **Coverage**: Complete NVD CPE dictionary
- **Reliability**: Automatic fallback to ML matcher if unavailable

## Monitoring

### Health Checks

- `GET /health`: Service health status
- `GET /metrics`: Prometheus metrics (if enabled)
- `GET http://localhost:8001/health`: CPE Guesser health check

### Logging

- Structured JSON logging
- Configurable log levels
- Log rotation support

## Documentation

- [API Documentation](../docs/api-v2-documentation.md)
- [Python Enrichment Documentation](../docs/python-enrichment.md)
- [Architecture Documentation](../docs/architecture.md)

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.

---

**Last Updated**: January 2025
