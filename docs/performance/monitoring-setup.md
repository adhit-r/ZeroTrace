# Performance Monitoring & Alerting Setup

## Overview
This document outlines the performance monitoring and alerting system for ZeroTrace components.

## Monitoring Architecture

### 1. Prometheus Metrics Exporters

#### Agent Metrics Exporter
```go
// agent-go/internal/monitoring/metrics.go
package monitoring

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Scan metrics
    scanDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "zerotrace_scan_duration_seconds",
            Help: "Duration of vulnerability scans",
        },
        []string{"scan_type", "status"},
    )
    
    filesScanned = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "zerotrace_files_scanned_total",
            Help: "Total number of files scanned",
        },
        []string{"file_type"},
    )
    
    vulnerabilitiesFound = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "zerotrace_vulnerabilities_found_total",
            Help: "Total number of vulnerabilities found",
        },
        []string{"severity", "type"},
    )
    
    // Resource metrics
    cpuUsage = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "zerotrace_agent_cpu_usage_percent",
            Help: "CPU usage percentage",
        },
        []string{"agent_id"},
    )
    
    memoryUsage = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "zerotrace_agent_memory_usage_bytes",
            Help: "Memory usage in bytes",
        },
        []string{"agent_id"},
    )
)
```

#### Enrichment Service Metrics
```python
# enrichment-python/app/monitoring.py
from prometheus_client import Counter, Histogram, Gauge, start_http_server
import time

# Metrics
enrichment_requests = Counter('zerotrace_enrichment_requests_total', 'Total enrichment requests', ['status'])
enrichment_duration = Histogram('zerotrace_enrichment_duration_seconds', 'Enrichment processing time')
cve_cache_hits = Counter('zerotrace_cve_cache_hits_total', 'CVE cache hits')
cve_cache_misses = Counter('zerotrace_cve_cache_misses_total', 'CVE cache misses')
active_connections = Gauge('zerotrace_enrichment_active_connections', 'Active connections')

class MetricsCollector:
    def __init__(self, port=8001):
        self.port = port
        start_http_server(port)
    
    def record_enrichment_request(self, status: str):
        enrichment_requests.labels(status=status).inc()
    
    def record_enrichment_duration(self, duration: float):
        enrichment_duration.observe(duration)
    
    def record_cache_hit(self):
        cve_cache_hits.inc()
    
    def record_cache_miss(self):
        cve_cache_misses.inc()
```

#### API Metrics
```go
// api-go/internal/monitoring/metrics.go
var (
    httpRequests = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "zerotrace_api_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "zerotrace_api_request_duration_seconds",
            Help: "HTTP request duration",
        },
        []string{"method", "endpoint"},
    )
    
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "zerotrace_api_active_connections",
            Help: "Active API connections",
        },
    )
)
```

### 2. Alert Rules

#### Prometheus Alert Rules
```yaml
# monitoring/alert-rules.yml
groups:
- name: zerotrace.rules
  rules:
  - alert: HighCPUUsage
    expr: zerotrace_agent_cpu_usage_percent > 80
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High CPU usage detected"
      description: "Agent {{ $labels.agent_id }} has high CPU usage: {{ $value }}%"
  
  - alert: HighMemoryUsage
    expr: zerotrace_agent_memory_usage_bytes > 100 * 1024 * 1024
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High memory usage detected"
      description: "Agent {{ $labels.agent_id }} has high memory usage: {{ $value }} bytes"
  
  - alert: SlowScanPerformance
    expr: zerotrace_scan_duration_seconds > 300
    for: 2m
    labels:
      severity: critical
    annotations:
      summary: "Slow scan performance"
      description: "Scan {{ $labels.scan_type }} is taking too long: {{ $value }}s"
  
  - alert: EnrichmentServiceDown
    expr: up{job="enrichment"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Enrichment service is down"
      description: "Enrichment service has been down for more than 1 minute"
  
  - alert: HighEnrichmentLatency
    expr: zerotrace_enrichment_duration_seconds > 30
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High enrichment latency"
      description: "Enrichment service is experiencing high latency: {{ $value }}s"
  
  - alert: APIDown
    expr: up{job="api"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "API service is down"
      description: "API service has been down for more than 1 minute"
  
  - alert: HighAPIResponseTime
    expr: zerotrace_api_request_duration_seconds > 5
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High API response time"
      description: "API is experiencing high response times: {{ $value }}s"
```

### 3. Grafana Dashboards

#### Agent Performance Dashboard
```json
{
  "dashboard": {
    "title": "ZeroTrace Agent Performance",
    "panels": [
      {
        "title": "CPU Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "zerotrace_agent_cpu_usage_percent",
            "legendFormat": "{{agent_id}}"
          }
        ]
      },
      {
        "title": "Memory Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "zerotrace_agent_memory_usage_bytes",
            "legendFormat": "{{agent_id}}"
          }
        ]
      },
      {
        "title": "Scan Duration",
        "type": "graph",
        "targets": [
          {
            "expr": "zerotrace_scan_duration_seconds",
            "legendFormat": "{{scan_type}}"
          }
        ]
      },
      {
        "title": "Vulnerabilities Found",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(zerotrace_vulnerabilities_found_total[5m])",
            "legendFormat": "{{severity}}"
          }
        ]
      }
    ]
  }
}
```

#### Enrichment Service Dashboard
```json
{
  "dashboard": {
    "title": "ZeroTrace Enrichment Service",
    "panels": [
      {
        "title": "Enrichment Requests",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(zerotrace_enrichment_requests_total[5m])",
            "legendFormat": "{{status}}"
          }
        ]
      },
      {
        "title": "Enrichment Duration",
        "type": "graph",
        "targets": [
          {
            "expr": "zerotrace_enrichment_duration_seconds",
            "legendFormat": "Duration"
          }
        ]
      },
      {
        "title": "Cache Hit Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(zerotrace_cve_cache_hits_total[5m]) / (rate(zerotrace_cve_cache_hits_total[5m]) + rate(zerotrace_cve_cache_misses_total[5m]))",
            "legendFormat": "Hit Rate"
          }
        ]
      }
    ]
  }
}
```

### 4. Performance Runbooks

#### Agent Performance Issues
```markdown
# Agent Performance Troubleshooting

## High CPU Usage
1. Check scan configuration
2. Reduce scan depth
3. Increase scan interval
4. Check for stuck processes

## High Memory Usage
1. Check for memory leaks
2. Reduce cache size
3. Increase garbage collection frequency
4. Check for large file processing

## Slow Scan Performance
1. Check file system performance
2. Reduce concurrent file processing
3. Optimize include/exclude patterns
4. Check network connectivity
```

#### Enrichment Service Issues
```markdown
# Enrichment Service Troubleshooting

## High Latency
1. Check CVE database connectivity
2. Verify cache configuration
3. Check API rate limits
4. Monitor external service health

## Cache Issues
1. Check cache hit rate
2. Verify cache TTL settings
3. Check cache storage
4. Monitor cache size
```

### 5. Deployment Configuration

#### Docker Compose with Monitoring
```yaml
# docker-compose.monitoring.yml
version: '3.8'
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      - ./monitoring/alert-rules.yml:/etc/prometheus/alert-rules.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
      - '--web.enable-admin-api'

  grafana:
    image: grafana/grafana
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana-storage:/var/lib/grafana
      - ./monitoring/grafana/dashboards:/etc/grafana/provisioning/dashboards
      - ./monitoring/grafana/datasources:/etc/grafana/provisioning/datasources

  alertmanager:
    image: prom/alertmanager
    ports:
      - "9093:9093"
    volumes:
      - ./monitoring/alertmanager.yml:/etc/alertmanager/alertmanager.yml

volumes:
  grafana-storage:
```

#### Prometheus Configuration
```yaml
# monitoring/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "alert-rules.yml"

scrape_configs:
  - job_name: 'zerotrace-agent'
    static_configs:
      - targets: ['agent:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s

  - job_name: 'zerotrace-enrichment'
    static_configs:
      - targets: ['enrichment:8001']
    metrics_path: '/metrics'
    scrape_interval: 30s

  - job_name: 'zerotrace-api'
    static_configs:
      - targets: ['api:8080']
    metrics_path: '/metrics'
    scrape_interval: 30s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093
```

### 6. Performance Benchmarks

#### Target Performance Metrics
- **Agent CPU Usage**: < 10% average
- **Agent Memory Usage**: < 100MB
- **Scan Duration**: < 5 minutes for typical codebase
- **Enrichment Latency**: < 30 seconds per request
- **API Response Time**: < 1 second for 95th percentile
- **Cache Hit Rate**: > 80% for enrichment service

#### Performance Testing
```bash
# Run performance tests
./scripts/performance-test.sh

# Benchmark agent scanning
./scripts/benchmark-agent.sh

# Load test enrichment service
./scripts/load-test-enrichment.sh

# API performance testing
./scripts/api-performance-test.sh
```

## Implementation Checklist

- [ ] Add Prometheus metrics to all components
- [ ] Configure alert rules
- [ ] Set up Grafana dashboards
- [ ] Deploy monitoring stack
- [ ] Configure alerting channels
- [ ] Create performance runbooks
- [ ] Set up performance testing
- [ ] Document monitoring procedures

