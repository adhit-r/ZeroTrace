# ZeroTrace CPE Matching & CVE Enrichment Service

## Overview
The ZeroTrace CPE (Common Platform Enumeration) matching service provides intelligent mapping between raw application data from Go agents and official NVD CPE identifiers, enabling accurate CVE (Common Vulnerabilities and Exposures) enrichment. This service is critical for transforming raw OS package lists, Windows registry entries, and other agent-collected data into actionable vulnerability intelligence.

## Problem Statement
Vendor and application names from agents (raw OS package lists, Windows registry, etc.) often don't match NVD's official CPE format, leading to missed vulnerabilities and false positives.

## AI-Enhanced Matching Strategy
- **Primary**: Rule-based mapping with exact/fuzzy matching
- **Secondary**: AI/LLM model for complex cases only
- **Fallback**: Manual review for low-confidence matches
- **Cost Control**: Cache AI results to avoid repeated queries

## High-Level Architecture

### End-to-End Flow
```
Go Agent → Go API → Python CPE Matcher → DB → NVD Enrichment → Risk Scoring
```

### Detailed Flow:
1. **Agent (Go)** collects raw application data
2. **POST /ingest/topology** (Go API) normalizes & upserts assets + apps
3. **Enqueue** app records to `apps.to_enrich` (Redis Streams/Kafka)
4. **Python CPE Matcher** consumes `apps.to_enrich`:
   - Normalizes (cleanup/version strip)
   - Dictionary lookup (direct CPE mapping)
   - Fuzzy matching (rapidfuzz/Levenshtein)
   - Generates candidate CPEs + confidence scores
5. **Write** `app_cpe_links` (DB) and mark app as matched if confident
6. **Enqueue** matched CPEs to `cpe.to_nvd` topic
7. **Python NVD/Exploit fetcher** consumes `cpe.to_nvd`, resolves CVEs
8. **Risk-scorer** consumes app_cves and writes asset_risk rollups
9. **Web UI** reads asset_vuln_rollup, apps, nodes/tiles for display

## Design Goals & Constraints

### Core Objectives
- **Accurate**: Maximize correct CPE→CVE mapping while minimizing false positives
- **Explainable**: Store matching rationale + confidence for UI/analyst review
- **Fast & Scalable**: Support thousands of agent payloads per minute
- **Pluggable**: Allow seeding/overriding with manual mappings & vendor-supplied rules
- **Auditable**: Keep provenance (who/what matched, when, source of CPE)

### Performance Targets
- **Throughput**: 10,000+ apps processed per minute
- **Latency**: < 100ms per app matching
- **Accuracy**: > 95% precision, > 90% recall
- **Cache Hit Rate**: > 80% for repeated app patterns

## Data Model & Database Schema

### Core Tables

#### Applications Table
```sql
CREATE TABLE apps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    raw_name TEXT NOT NULL,                    -- raw string from agent
    raw_version TEXT,                          -- raw version string
    normalized_name TEXT,                      -- cleaned name used by matcher
    normalized_version TEXT,
    vendor_hint TEXT,                          -- if agent provided vendor
    package_type VARCHAR(50),                  -- rpm, deb, pip, npm, etc.
    architecture VARCHAR(20),                  -- x86_64, arm64, etc.
    first_seen TIMESTAMPTZ DEFAULT now(),
    last_seen TIMESTAMPTZ DEFAULT now(),
    match_status VARCHAR(20) DEFAULT 'pending', -- pending, matched, failed
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
```

#### CPE Links Table
```sql
CREATE TABLE app_cpe_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    cpe TEXT NOT NULL,                         -- cpe:2.3:...
    cpe_source VARCHAR(50) NOT NULL,           -- 'dictionary', 'fuzzy', 'manual', 'ai'
    confidence NUMERIC(3,2) NOT NULL,          -- 0.00 - 1.00
    is_primary BOOLEAN DEFAULT false,          -- primary match for this app
    matched_at TIMESTAMPTZ DEFAULT now(),
    rationale JSONB,                           -- details of matching steps
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(app_id, cpe)
);
```

#### CPE Dictionary Table
```sql
CREATE TABLE cpe_dictionary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cpe TEXT UNIQUE NOT NULL,
    vendor TEXT NOT NULL,
    product TEXT NOT NULL,
    version TEXT,
    title TEXT,
    description TEXT,
    deprecated BOOLEAN DEFAULT false,
    last_updated TIMESTAMPTZ DEFAULT now(),
    created_at TIMESTAMPTZ DEFAULT now()
);
```

## Message Formats & API Contracts

### Agent Payload Format
```json
{
  "agent_id": "agent-001",
  "company_id": "company-123",
  "timestamp": "2024-01-15T10:30:00Z",
  "host": {
    "hostname": "web-01",
    "site": "BLR-HQ",
    "os": "ubuntu",
    "os_version": "20.04"
  },
  "apps": [
    {
      "name": "nginx",
      "version": "1.18.0-ubuntu1.3",
      "vendor": "",
      "package_type": "deb",
      "architecture": "amd64",
      "candidates": ["nginx:1.18.0"]
    }
  ]
}
```

### Queue Message Format (apps.to_enrich)
```json
{
  "app_id": "app-001",
  "company_id": "company-123",
  "asset_id": "asset-001",
  "raw_name": "nginx",
  "raw_version": "1.18.0-ubuntu1.3",
  "normalized_name": "nginx",
  "normalized_version": "1.18.0",
  "vendor_hint": null,
  "package_type": "deb",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## CPE Matching Service Components

### 1. Application Normalizer
```python
class ApplicationNormalizer:
    """Normalizes application names and versions"""
    
    def normalize_name(self, raw_name: str) -> str:
        """Normalize application name"""
        name = raw_name.lower()
        name = re.sub(r'^lib', '', name)  # libjpeg -> jpeg
        name = re.sub(r'^python3?-', '', name)  # python3-requests -> requests
        name = re.sub(r'-[0-9]+\.[0-9]+.*$', '', name)  # nginx-1.18.0 -> nginx
        return name.strip()
    
    def normalize_version(self, raw_version: str) -> str:
        """Normalize version string"""
        if not raw_version:
            return ""
        version = re.sub(r'-[a-z]+[0-9]*\.[0-9]+$', '', raw_version)
        version = re.sub(r'\+[a-z0-9]+$', '', version)
        return version.strip()
```

### 2. CPE Matcher
```python
class CPEMatcher:
    """Main CPE matching engine"""
    
    async def match_application(self, app_data: dict) -> List[dict]:
        """Match application to CPE candidates"""
        candidates = []
        
        # 1. Manual mapping lookup
        manual_matches = await self._check_manual_mappings(app_data['raw_name'])
        if manual_matches:
            candidates.extend(manual_matches)
        
        # 2. Exact dictionary lookup
        exact_matches = await self._exact_match(app_data['normalized_name'])
        if exact_matches:
            candidates.extend(exact_matches)
        
        # 3. Fuzzy matching
        fuzzy_matches = await self._fuzzy_match(app_data['normalized_name'])
        candidates.extend(fuzzy_matches)
        
        # 4. Version matching and confidence scoring
        scored_candidates = await self._score_candidates(candidates, app_data)
        
        return scored_candidates
```

## Matching Algorithm

### Step-by-Step Process

#### 1. Normalization
```python
def normalize_application(raw_name: str, raw_version: str) -> dict:
    """Normalize application name and version"""
    normalized_name = normalizer.normalize_name(raw_name)
    normalized_version = normalizer.normalize_version(raw_version)
    
    return {
        'normalized_name': normalized_name,
        'normalized_version': normalized_version,
        'original_name': raw_name,
        'original_version': raw_version
    }
```

#### 2. Manual Mapping Lookup
```python
async def check_manual_mappings(raw_name: str, company_id: str) -> List[dict]:
    """Check manual mapping rules"""
    query = """
        SELECT cpe, confidence, description 
        FROM manual_mappings 
        WHERE company_id = $1 AND is_active = true
    """
    
    mappings = await db.fetch(query, company_id)
    matches = []
    
    for mapping in mappings:
        if re.match(mapping['raw_pattern'], raw_name, re.IGNORECASE):
            matches.append({
                'cpe': mapping['cpe'],
                'confidence': mapping['confidence'],
                'source': 'manual',
                'description': mapping['description']
            })
    
    return matches
```

#### 3. Fuzzy Matching
```python
async def fuzzy_match(normalized_name: str, limit: int = 5) -> List[dict]:
    """Fuzzy matching using RapidFuzz"""
    from rapidfuzz import process, fuzz
    
    # Get candidate CPEs (pre-filtered for performance)
    candidates = await get_cpe_candidates(normalized_name)
    
    # Perform fuzzy matching
    corpus = {
        i: f"{c['vendor']} {c['product']} {c.get('title', '')}"
        for i, c in enumerate(candidates)
    }
    
    matches = process.extract(
        normalized_name,
        corpus,
        scorer=fuzz.token_sort_ratio,
        limit=limit
    )
    
    return [
        {
            'cpe': candidates[idx]['cpe'],
            'confidence': min(score / 100.0, 0.8),  # Cap at 0.8 for fuzzy
            'source': 'fuzzy',
            'score': score,
            'vendor': candidates[idx]['vendor'],
            'product': candidates[idx]['product']
        }
        for _, score, idx in matches if score >= 70
    ]
```

## Confidence Scoring

### Confidence Calculation
```python
def calculate_confidence(candidate: dict, app_data: dict) -> float:
    """Calculate confidence score for CPE candidate"""
    base_score = candidate.get('confidence', 0.0)
    
    # Vendor hint bonus
    vendor_bonus = 0.0
    if app_data.get('vendor_hint') and candidate.get('vendor'):
        if app_data['vendor_hint'].lower() == candidate['vendor'].lower():
            vendor_bonus = 0.15
    
    # Version matching bonus
    version_bonus = 0.0
    if app_data.get('normalized_version') and candidate.get('version'):
        if app_data['normalized_version'] == candidate['version']:
            version_bonus = 0.2
    
    # Source bonus
    source_bonus = {
        'manual': 0.2,
        'dictionary': 0.1,
        'fuzzy': 0.0
    }.get(candidate.get('source', ''), 0.0)
    
    total_confidence = base_score + vendor_bonus + version_bonus + source_bonus
    return min(total_confidence, 1.0)
```

### Confidence Thresholds
- **≥ 0.95**: Automatic match (very safe; manual mapping or exact)
- **≥ 0.85**: Auto-enrich & auto-assign (good confidence)
- **0.5–0.85**: Candidate; show in UI with "low confidence" warning
- **< 0.5**: Record for analytics; don't auto-assign

## Performance Optimizations

### 1. Caching Strategy
```python
class CPEMatcherCache:
    """Multi-level caching for CPE matching"""
    
    def __init__(self):
        self.memory_cache = {}  # LRU cache for hot items
        self.redis_cache = redis.Redis()  # Distributed cache
        self.cache_ttl = 3600  # 1 hour
    
    async def get_cached_result(self, app_key: str) -> Optional[dict]:
        """Get cached matching result"""
        # Try memory cache first
        if app_key in self.memory_cache:
            return self.memory_cache[app_key]
        
        # Try Redis cache
        cached = await self.redis_cache.get(f"cpe_match:{app_key}")
        if cached:
            result = json.loads(cached)
            self.memory_cache[app_key] = result
            return result
        
        return None
```

### 2. Bulk Processing
```python
async def process_app_batch(apps: List[dict]) -> List[dict]:
    """Process multiple apps in batch for efficiency"""
    results = []
    
    # Group apps by normalized name for batch processing
    app_groups = {}
    for app in apps:
        key = app['normalized_name']
        if key not in app_groups:
            app_groups[key] = []
        app_groups[key].append(app)
    
    # Process each group
    for normalized_name, app_group in app_groups.items():
        # Get CPE candidates once for the group
        candidates = await get_cpe_candidates(normalized_name)
        
        # Match each app in the group
        for app in app_group:
            app_results = await match_app_with_candidates(app, candidates)
            results.extend(app_results)
    
    return results
```

## Integration with NVD Enrichment

### NVD Fetcher Service
```python
class NVDFetcher:
    """Fetches CVE data from NVD for matched CPEs"""
    
    def __init__(self, db_connection, nvd_api_key: str = None):
        self.db = db_connection
        self.nvd_api_key = nvd_api_key
        self.api_base = "https://services.nvd.nist.gov/rest/json/cves/2.0"
    
    async def fetch_cves_for_cpe(self, cpe: str) -> List[dict]:
        """Fetch CVEs for a specific CPE"""
        params = {
            'cpeName': cpe,
            'resultsPerPage': 2000
        }
        
        if self.nvd_api_key:
            params['apiKey'] = self.nvd_api_key
        
        async with aiohttp.ClientSession() as session:
            async with session.get(self.api_base, params=params) as response:
                if response.status == 200:
                    data = await response.json()
                    return data.get('vulnerabilities', [])
                else:
                    logger.error(f"NVD API error: {response.status}")
                    return []
```

## UI/UX Considerations

### Confidence Display
```typescript
interface CPEMatchDisplay {
  cpe: string;
  confidence: number;
  source: 'manual' | 'dictionary' | 'fuzzy' | 'ai';
  rationale: string;
  isPrimary: boolean;
}

const ConfidenceBar: React.FC<{ confidence: number }> = ({ confidence }) => {
  const getColor = (conf: number) => {
    if (conf >= 0.95) return 'bg-green-500';
    if (conf >= 0.85) return 'bg-yellow-500';
    if (conf >= 0.5) return 'bg-orange-500';
    return 'bg-red-500';
  };
  
  return (
    <div className="flex items-center space-x-2">
      <div className="w-20 bg-gray-200 rounded-full h-2">
        <div 
          className={`h-2 rounded-full ${getColor(confidence)}`}
          style={{ width: `${confidence * 100}%` }}
        />
      </div>
      <span className="text-sm font-medium">
        {Math.round(confidence * 100)}%
      </span>
    </div>
  );
};
```

## Testing & Quality Assurance

### Unit Tests
```python
import pytest
from unittest.mock import AsyncMock, patch

class TestCPEMatcher:
    @pytest.mark.asyncio
    async def test_exact_match(self):
        """Test exact dictionary matching"""
        matcher = CPEMatcher(mock_db, mock_normalizer, mock_dictionary)
        
        with patch.object(matcher, '_exact_match') as mock_exact:
            mock_exact.return_value = [
                {
                    'cpe': 'cpe:2.3:a:nginx:nginx:1.18.0:*:*:*:*:*:*:*',
                    'confidence': 0.9,
                    'source': 'dictionary'
                }
            ]
            
            result = await matcher.match_application({
                'normalized_name': 'nginx',
                'raw_name': 'nginx',
                'raw_version': '1.18.0'
            })
            
            assert len(result) == 1
            assert result[0]['cpe'] == 'cpe:2.3:a:nginx:nginx:1.18.0:*:*:*:*:*:*:*'
            assert result[0]['confidence'] >= 0.9
```

## Observability & Metrics

### Key Metrics
```python
class CPEMatchingMetrics:
    """Metrics collection for CPE matching"""
    
    def __init__(self):
        self.match_rate = Counter('cpe_match_rate', 'Apps with any CPE match')
        self.auto_assign_rate = Counter('cpe_auto_assign_rate', 'Apps auto-assigned')
        self.human_review_rate = Counter('cpe_human_review_rate', 'Apps requiring review')
        self.processing_time = Histogram('cpe_processing_time', 'Time to process app')
        self.queue_backlog = Gauge('cpe_queue_backlog', 'Queue backlog size')
    
    def record_match(self, app_id: str, confidence: float):
        """Record successful match"""
        self.match_rate.inc()
        
        if confidence >= 0.85:
            self.auto_assign_rate.inc()
        elif confidence >= 0.5:
            self.human_review_rate.inc()
```

## Implementation Roadmap

### Phase 1: Core Matching (Week 1-2)
- [ ] Implement CPE dictionary loader
- [ ] Build basic normalization logic
- [ ] Implement exact matching
- [ ] Create database schema
- [ ] Basic API endpoints

### Phase 2: Advanced Matching (Week 3-4)
- [ ] Implement fuzzy matching
- [ ] Add confidence scoring
- [ ] Build manual mapping system
- [ ] Create queue integration
- [ ] Basic monitoring

### Phase 3: Integration (Week 5-6)
- [ ] Integrate with Go API
- [ ] Deploy NVD fetcher
- [ ] Connect to risk scoring
- [ ] Update web UI
- [ ] End-to-end testing

### Phase 4: Optimization (Week 7-8)
- [ ] Performance optimization
- [ ] Advanced caching
- [ ] Load testing
- [ ] Production deployment
- [ ] Documentation

## Operational Notes

### Deployment Checklist
1. **Database Setup**
   - Create CPE dictionary tables
   - Load initial NVD CPE data
   - Set up indexes and partitions
   - Configure row-level security

2. **Service Deployment**
   - Deploy CPE matcher service
   - Configure Redis/Kafka queues
   - Set up monitoring and alerting
   - Configure confidence thresholds

3. **Integration**
   - Update Go API to enqueue apps
   - Deploy NVD fetcher service
   - Configure risk scoring pipeline
   - Update web UI for confidence display

### Security Considerations
- **Input Validation**: Validate all agent inputs
- **Rate Limiting**: Limit API requests to prevent abuse
- **Access Control**: Implement proper authentication/authorization
- **Data Sanitization**: Sanitize all data before processing
- **Audit Logging**: Log all matching decisions for audit

This enhanced CPE matching service provides a robust, scalable solution for accurately mapping agent-collected application data to official CPE identifiers, enabling precise CVE enrichment and risk assessment in the ZeroTrace platform. 