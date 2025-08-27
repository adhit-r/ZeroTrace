# ZeroTrace Scalable Data Processing Architecture

## 🎯 **Problem Statement**

**Scale Requirements:**
- **1000s of agents** per company
- **100s of companies** (clients)
- **100s of apps** per agent
- **Millions of app records** to process daily
- **Real-time CVE enrichment** for each app
- **Non-multi-tenant** but **organization isolation**
- **Single database** with company-based separation

## 🏗️ **Robust Architecture Design**

### **1. Data Flow Architecture**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Agent (1000s) │───▶│   API Gateway   │───▶│   Data Pipeline │
│   Per Company   │    │   (Rate Limit)  │    │   (Processing)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Queue System  │    │   CVE Enrichment│
                       │   (Redis/Rabbit)│    │   (Python)      │
                       └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Database      │    │   UI Dashboard  │
                       │   (PostgreSQL)  │    │   (Real-time)   │
                       └─────────────────┘    └─────────────────┘
```

### **2. Key Components**

#### **A. API Gateway with Rate Limiting**
```go
// api-go/internal/middleware/rate_limit.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "golang.org/x/time/rate"
    "sync"
    "time"
)

type RateLimiter struct {
    redis *redis.Client
    limiters map[string]*rate.Limiter
    mu sync.RWMutex
}

// Company-based rate limiting
func CompanyRateLimit(redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        companyID := c.GetString("company_id")
        agentID := c.GetString("agent_id")
        
        // Rate limit per company: 10,000 requests/minute
        key := fmt.Sprintf("rate_limit:company:%s", companyID)
        allowed := checkRateLimit(redis, key, 10000, time.Minute)
        
        if !allowed {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        // Rate limit per agent: 100 requests/minute
        agentKey := fmt.Sprintf("rate_limit:agent:%s", agentID)
        allowed = checkRateLimit(redis, key, 100, time.Minute)
        
        if !allowed {
            c.JSON(429, gin.H{"error": "Agent rate limit exceeded"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

#### **B. Queue-Based Processing**
```go
// api-go/internal/queue/processor.go
package queue

import (
    "encoding/json"
    "log"
    "time"
    "github.com/go-redis/redis/v8"
)

type AppData struct {
    ID          string    `json:"id"`
    CompanyID   string    `json:"company_id"`
    AgentID     string    `json:"agent_id"`
    AppName     string    `json:"app_name"`
    AppVersion  string    `json:"app_version"`
    PackageType string    `json:"package_type"`
    Timestamp   time.Time `json:"timestamp"`
}

type QueueProcessor struct {
    redis *redis.Client
    batchSize int
    workers   int
}

func NewQueueProcessor(redis *redis.Client) *QueueProcessor {
    return &QueueProcessor{
        redis: redis,
        batchSize: 100, // Process 100 apps at a time
        workers: 10,    // 10 concurrent workers
    }
}

// Process apps in batches
func (qp *QueueProcessor) ProcessApps() {
    for i := 0; i < qp.workers; i++ {
        go qp.worker()
    }
}

func (qp *QueueProcessor) worker() {
    for {
        // Get batch of apps from queue
        apps := qp.getBatchFromQueue()
        
        if len(apps) > 0 {
            // Send to Python enrichment service
            qp.enrichApps(apps)
        }
        
        time.Sleep(100 * time.Millisecond)
    }
}

func (qp *QueueProcessor) getBatchFromQueue() []AppData {
    var apps []AppData
    
    // Get up to batchSize items from queue
    for i := 0; i < qp.batchSize; i++ {
        result, err := qp.redis.BLPop(ctx, 1*time.Second, "app_queue").Result()
        if err != nil {
            break
        }
        
        var app AppData
        json.Unmarshal([]byte(result[1]), &app)
        apps = append(apps, app)
    }
    
    return apps
}
```

#### **C. Optimized Database Schema**
```sql
-- Optimized schema for massive scale
-- 1. Partitioned tables by company_id
CREATE TABLE apps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    agent_id UUID NOT NULL,
    app_name VARCHAR(255) NOT NULL,
    app_version VARCHAR(100),
    package_type VARCHAR(50),
    architecture VARCHAR(20),
    first_seen TIMESTAMPTZ DEFAULT now(),
    last_seen TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
) PARTITION BY HASH (company_id);

-- Create partitions for each company (hash-based)
CREATE TABLE apps_partition_0 PARTITION OF apps FOR VALUES WITH (modulus 100, remainder 0);
CREATE TABLE apps_partition_1 PARTITION OF apps FOR VALUES WITH (modulus 100, remainder 1);
-- ... up to partition_99

-- 2. Optimized indexes
CREATE INDEX CONCURRENTLY idx_apps_company_agent ON apps(company_id, agent_id);
CREATE INDEX CONCURRENTLY idx_apps_name_version ON apps(app_name, app_version);
CREATE INDEX CONCURRENTLY idx_apps_last_seen ON apps(last_seen);

-- 3. Vulnerabilities table (partitioned)
CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    agent_id UUID NOT NULL,
    app_id UUID NOT NULL,
    cve_id VARCHAR(20),
    severity VARCHAR(20),
    cvss_score DECIMAL(3,1),
    title TEXT,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
) PARTITION BY HASH (company_id);

-- 4. Materialized view for fast queries
CREATE MATERIALIZED VIEW company_vulnerability_summary AS
SELECT 
    company_id,
    COUNT(*) as total_vulnerabilities,
    COUNT(*) FILTER (WHERE severity = 'critical') as critical_count,
    COUNT(*) FILTER (WHERE severity = 'high') as high_count,
    COUNT(*) FILTER (WHERE severity = 'medium') as medium_count,
    COUNT(*) FILTER (WHERE severity = 'low') as low_count,
    MAX(created_at) as last_updated
FROM vulnerabilities
GROUP BY company_id;

-- Refresh every 5 minutes
CREATE OR REPLACE FUNCTION refresh_vulnerability_summary()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY company_vulnerability_summary;
END;
$$ LANGUAGE plpgsql;

-- 5. Automated cleanup (keep last 90 days)
CREATE OR REPLACE FUNCTION cleanup_old_data()
RETURNS void AS $$
BEGIN
    DELETE FROM apps WHERE last_seen < now() - interval '90 days';
    DELETE FROM vulnerabilities WHERE created_at < now() - interval '90 days';
END;
$$ LANGUAGE plpgsql;
```

#### **D. Python Enrichment Service (Optimized)**
```python
# enrichment-python/app/batch_enrichment.py
import asyncio
import aiohttp
import json
from typing import List, Dict
from dataclasses import dataclass
from datetime import datetime
import logging

@dataclass
class AppData:
    id: str
    company_id: str
    agent_id: str
    app_name: str
    app_version: str
    package_type: str

class BatchEnrichmentService:
    def __init__(self):
        self.session = None
        self.batch_size = 100
        self.max_concurrent = 10
        self.cache_ttl = 3600  # 1 hour
        
    async def process_batch(self, apps: List[AppData]) -> List[Dict]:
        """Process a batch of apps efficiently"""
        
        # 1. Check cache first
        cached_results = await self.get_cached_results(apps)
        uncached_apps = [app for app in apps if app.id not in cached_results]
        
        # 2. Process uncached apps in parallel
        if uncached_apps:
            enriched_results = await self.enrich_apps_parallel(uncached_apps)
            
            # 3. Cache results
            await self.cache_results(enriched_results)
            
            # 4. Combine results
            all_results = {**cached_results, **enriched_results}
        else:
            all_results = cached_results
            
        return [all_results[app.id] for app in apps]
    
    async def enrich_apps_parallel(self, apps: List[AppData]) -> Dict[str, Dict]:
        """Enrich apps in parallel with rate limiting"""
        
        semaphore = asyncio.Semaphore(self.max_concurrent)
        
        async def enrich_single(app: AppData) -> tuple:
            async with semaphore:
                try:
                    # Rate limiting: 100 requests per second to NVD
                    await asyncio.sleep(0.01)  # 10ms delay
                    
                    cve_data = await self.get_cve_data(app.app_name, app.app_version)
                    return (app.id, {
                        'app_id': app.id,
                        'company_id': app.company_id,
                        'agent_id': app.agent_id,
                        'cve_data': cve_data,
                        'enriched_at': datetime.utcnow().isoformat()
                    })
                except Exception as e:
                    logging.error(f"Failed to enrich app {app.id}: {e}")
                    return (app.id, {
                        'app_id': app.id,
                        'error': str(e),
                        'enriched_at': datetime.utcnow().isoformat()
                    })
        
        # Process all apps concurrently
        tasks = [enrich_single(app) for app in apps]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        return dict(results)
    
    async def get_cve_data(self, app_name: str, app_version: str) -> List[Dict]:
        """Get CVE data from multiple sources"""
        
        # 1. NVD API
        nvd_cves = await self.get_nvd_cves(app_name, app_version)
        
        # 2. CVE Search API
        cve_search_cves = await self.get_cve_search_cves(app_name, app_version)
        
        # 3. Merge and deduplicate
        all_cves = self.merge_cve_data(nvd_cves, cve_search_cves)
        
        return all_cves

# FastAPI endpoint for batch processing
@app.post("/enrich/batch")
async def enrich_batch(apps: List[AppData]):
    """Process batch of apps efficiently"""
    
    service = BatchEnrichmentService()
    results = await service.process_batch(apps)
    
    return {
        "processed": len(results),
        "results": results,
        "timestamp": datetime.utcnow().isoformat()
    }
```

#### **E. Real-time UI Updates**
```typescript
// web-react/src/hooks/useRealTimeData.ts
import { useEffect, useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';

export const useRealTimeData = (companyId: string) => {
    const queryClient = useQueryClient();
    const [lastUpdate, setLastUpdate] = useState<Date>(new Date());

    // Poll for updates every 30 seconds
    const { data: vulnerabilities } = useQuery({
        queryKey: ['vulnerabilities', companyId],
        queryFn: () => fetchVulnerabilities(companyId),
        refetchInterval: 30000, // 30 seconds
        staleTime: 10000, // 10 seconds
    });

    // WebSocket for real-time updates
    useEffect(() => {
        const ws = new WebSocket(`ws://localhost:8080/ws/company/${companyId}`);
        
        ws.onmessage = (event) => {
            const update = JSON.parse(event.data);
            
            // Update specific data based on type
            switch (update.type) {
                case 'new_vulnerability':
                    queryClient.invalidateQueries(['vulnerabilities', companyId]);
                    break;
                case 'agent_status':
                    queryClient.invalidateQueries(['agents', companyId]);
                    break;
                case 'scan_complete':
                    queryClient.invalidateQueries(['scans', companyId]);
                    break;
            }
            
            setLastUpdate(new Date());
        };

        return () => ws.close();
    }, [companyId, queryClient]);

    return { vulnerabilities, lastUpdate };
};
```

### **3. Performance Optimizations**

#### **A. Database Optimizations**
```sql
-- 1. Connection pooling
-- In postgresql.conf
max_connections = 200
shared_buffers = 1GB
effective_cache_size = 4GB
work_mem = 16MB
maintenance_work_mem = 256MB

-- 2. Query optimization
-- Use prepared statements
-- Implement query result caching
-- Use materialized views for complex aggregations

-- 3. Partitioning strategy
-- Partition by company_id (hash-based)
-- Partition by date for time-series data
-- Automatic partition management
```

#### **B. Caching Strategy**
```go
// api-go/internal/cache/cache.go
package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/go-redis/redis/v8"
)

type CacheManager struct {
    redis *redis.Client
}

// Cache vulnerability data for 1 hour
func (cm *CacheManager) CacheVulnerabilities(companyID string, data interface{}) error {
    key := fmt.Sprintf("vulns:company:%s", companyID)
    return cm.redis.Set(context.Background(), key, data, time.Hour).Err()
}

// Cache app data for 30 minutes
func (cm *CacheManager) CacheApps(companyID string, data interface{}) error {
    key := fmt.Sprintf("apps:company:%s", companyID)
    return cm.redis.Set(context.Background(), key, data, 30*time.Minute).Err()
}

// Cache CVE data for 24 hours (rarely changes)
func (cm *CacheManager) CacheCVEData(cveID string, data interface{}) error {
    key := fmt.Sprintf("cve:%s", cveID)
    return cm.redis.Set(context.Background(), key, data, 24*time.Hour).Err()
}
```

#### **C. Batch Processing**
```go
// api-go/internal/processor/batch.go
package processor

import (
    "time"
    "sync"
)

type BatchProcessor struct {
    batchSize    int
    batchTimeout time.Duration
    workers      int
    queue        chan AppData
    results      chan ProcessedResult
}

func NewBatchProcessor() *BatchProcessor {
    return &BatchProcessor{
        batchSize:    100,
        batchTimeout: 5 * time.Second,
        workers:      10,
        queue:        make(chan AppData, 10000),
        results:      make(chan ProcessedResult, 10000),
    }
}

func (bp *BatchProcessor) Start() {
    // Start workers
    for i := 0; i < bp.workers; i++ {
        go bp.worker()
    }
    
    // Start batch collector
    go bp.batchCollector()
}

func (bp *BatchProcessor) batchCollector() {
    var batch []AppData
    timer := time.NewTimer(bp.batchTimeout)
    
    for {
        select {
        case app := <-bp.queue:
            batch = append(batch, app)
            
            if len(batch) >= bp.batchSize {
                bp.processBatch(batch)
                batch = batch[:0]
                timer.Reset(bp.batchTimeout)
            }
            
        case <-timer.C:
            if len(batch) > 0 {
                bp.processBatch(batch)
                batch = batch[:0]
            }
            timer.Reset(bp.batchTimeout)
        }
    }
}
```

### **4. Monitoring & Alerting**

#### **A. Performance Metrics**
```go
// api-go/internal/metrics/performance.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Processing metrics
    appsProcessedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "apps_processed_total",
            Help: "Total number of apps processed",
        },
        []string{"company_id", "status"},
    )
    
    processingDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "app_processing_duration_seconds",
            Help:    "Time taken to process apps",
            Buckets: prometheus.DefBuckets,
        },
        []string{"company_id"},
    )
    
    // Queue metrics
    queueSize = promauto.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "app_queue_size",
            Help: "Number of apps in processing queue",
        },
        []string{"company_id"},
    )
    
    // CVE enrichment metrics
    cveEnrichmentDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "cve_enrichment_duration_seconds",
            Help:    "Time taken to enrich CVE data",
            Buckets: prometheus.DefBuckets,
        },
        []string{"company_id"},
    )
)
```

#### **B. Alerting Rules**
```yaml
# prometheus/alerts.yml
groups:
  - name: zerotrace_processing
    rules:
      # High queue size
      - alert: HighQueueSize
        expr: app_queue_size > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High app processing queue"
          description: "Queue size is {{ $value }} apps"

      # Slow processing
      - alert: SlowProcessing
        expr: histogram_quantile(0.95, app_processing_duration_seconds) > 30
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Slow app processing"
          description: "95th percentile processing time is {{ $value }}s"

      # CVE enrichment failures
      - alert: CVEEnrichmentFailures
        expr: rate(cve_enrichment_failures_total[5m]) > 0.1
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "CVE enrichment failures"
          description: "{{ $value }} failures per second"
```

### **5. Scalability Features**

#### **A. Horizontal Scaling**
```yaml
# docker-compose.scale.yml
version: '3.8'

services:
  api:
    image: zerotrace/api
    deploy:
      replicas: 5  # Scale API instances
    environment:
      - WORKER_ID=${HOSTNAME}
  
  enrichment:
    image: zerotrace/enrichment
    deploy:
      replicas: 10  # Scale enrichment workers
    environment:
      - WORKER_ID=${HOSTNAME}
  
  redis:
    image: redis:7-alpine
    deploy:
      replicas: 3  # Redis cluster
```

#### **B. Load Balancing**
```nginx
# nginx/load-balancer.conf
upstream api_backend {
    least_conn;  # Least connections algorithm
    server api1:8080 max_fails=3 fail_timeout=30s;
    server api2:8080 max_fails=3 fail_timeout=30s;
    server api3:8080 max_fails=3 fail_timeout=30s;
    server api4:8080 max_fails=3 fail_timeout=30s;
    server api5:8080 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name api.zerotrace.com;
    
    location / {
        proxy_pass http://api_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        
        # Rate limiting per company
        limit_req_zone $http_x_company_id zone=company:10m rate=1000r/m;
        limit_req zone=company burst=2000 nodelay;
    }
}
```

## 🎯 **Key Benefits**

### **1. Scalability**
- ✅ **Handles 1000s of agents** per company
- ✅ **Processes 100s of companies** simultaneously
- ✅ **Manages millions of app records** efficiently
- ✅ **Real-time processing** with minimal latency

### **2. Reliability**
- ✅ **Queue-based processing** prevents data loss
- ✅ **Batch processing** optimizes throughput
- ✅ **Caching** reduces redundant work
- ✅ **Error handling** and retry mechanisms

### **3. Performance**
- ✅ **Parallel processing** for CVE enrichment
- ✅ **Database partitioning** for fast queries
- ✅ **Connection pooling** for database efficiency
- ✅ **Rate limiting** prevents overload

### **4. Monitoring**
- ✅ **Real-time metrics** for all components
- ✅ **Alerting** for performance issues
- ✅ **Dashboard** for system health
- ✅ **Logging** for debugging

## 🚀 **Implementation Timeline**

### **Week 1: Core Infrastructure**
- Set up queue system (Redis)
- Implement rate limiting
- Create database partitions
- Basic batch processing

### **Week 2: Processing Pipeline**
- Optimize Python enrichment service
- Implement caching strategy
- Add monitoring metrics
- Create alerting rules

### **Week 3: UI & Real-time Updates**
- Implement WebSocket connections
- Add real-time data updates
- Optimize UI performance
- Add data visualization

### **Week 4: Testing & Optimization**
- Load testing with realistic data
- Performance optimization
- Error handling improvements
- Documentation

This architecture can handle **millions of app records** from **thousands of agents** across **hundreds of companies** while maintaining **real-time performance** and **enterprise reliability**.
