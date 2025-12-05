# Architecture Optimization Plan

## Current Complexity Issues

### 1. **Excessive Caching Layers** 🔴 HIGH PRIORITY
**Problem**: 3-4 levels of caching across services
- Go API: In-memory + Valkey
- Python Enrichment: L1 (memory) + L2 (Valkey) + L3 (Memcached)
- Multiple LRU caches in different modules
- CPE Guesser: Separate Valkey instance

**Impact**: 
- Increased memory usage
- Cache invalidation complexity
- Debugging difficulty
- Maintenance overhead

**Solution**: Consolidate to 2-tier caching
- **L1**: In-memory cache (fast, process-local)
- **L2**: Valkey (shared, persistent)
- **Remove**: Memcached (redundant with Valkey)

### 2. **Redundant Database Connections** 🔴 HIGH PRIORITY
**Problem**: 
- PostgreSQL configured in Python but not used
- Multiple Valkey database numbers (0, 8) - can be consolidated
- Mixed Redis/Valkey references causing confusion

**Impact**:
- Unnecessary connection overhead
- Configuration complexity
- Resource waste

**Solution**: 
- Use single Valkey instance with database separation (DB 0 for API, DB 8 for CPE)
- Remove unused PostgreSQL config from Python service
- Standardize on Valkey everywhere (no Redis) 

### 3. **Service Proliferation in Go API** 🟡 MEDIUM PRIORITY
**Problem**: 9 different services initialized
```
- scanService
- agentService  
- enrollmentService
- vulnerabilityV2Service
- organizationProfileService
- techStackService
- heatmapService
- maturityService
- complianceService
```

**Analysis**: Many services have overlapping responsibilities:
- `heatmapService`, `maturityService`, `complianceService` all analyze vulnerabilities
- `techStackService` and `organizationProfileService` both handle org data
- `vulnerabilityV2Service` and main vulnerability handling overlap

**Solution**: Consolidate related services
- Merge `heatmapService`, `maturityService`, `complianceService` into `analyticsService`
- Merge `techStackService` into `organizationProfileService`
- Simplify service boundaries

### 4. **CPE Guesser as Separate Service** 🟡 MEDIUM PRIORITY
**Problem**: CPE Guesser is a separate service but could be integrated

**Current**: 
- Separate Valkey connection
- Separate HTTP server
- Called via HTTP from enrichment service

**Solution Options**:
- **Option A**: Integrate CPE Guesser directly into enrichment service (recommended)
- **Option B**: Keep separate but share Valkey connection pool

### 5. **Multiple Config Files** 🟢 LOW PRIORITY
**Problem**: Config scattered across multiple files
- `api-go/internal/config/config.go`
- `enrichment-python/app/config.py`
- `enrichment-python/cpe-guesser/lib/config.py`

**Solution**: Centralize configuration or use shared config service

## Optimization Recommendations

### Phase 1: Quick Wins (Low Risk, High Impact)

#### 1.1 Remove Memcached
**Action**: Remove Memcached from Python enrichment service
- Valkey already provides all needed functionality
- Reduces infrastructure complexity
- Saves memory and connections

**Files to modify**:
- `enrichment-python/app/config.py` - Remove memcached config
- `enrichment-python/app/ultra_optimized_enrichment.py` - Remove L3 cache
- `enrichment-python/app/cve_enrichment.py` - Simplify caching

#### 1.2 Standardize on Valkey
**Action**: Use single Valkey instance with database separation
- DB 0: Go API caching
- DB 8: CPE Guesser (keep separate for data isolation)
- Replace all Redis references with Valkey everywhere
- Update Docker Compose to use valkey/valkey:latest image

**Files to modify**:
- `docker-compose.yml` - Change redis service to valkey
- `api-go/internal/config/config.go` - Update comments/docs to Valkey
- `enrichment-python/app/config.py` - Update Redis references to Valkey
- `enrichment-python/cpe-guesser/lib/config.py` - Already using Valkey (keep)
- `enrichment-python/cpe-guesser/lib/cpeguesser.py` - Already using Valkey (keep)
- All environment variable examples and documentation

#### 1.3 Remove Unused PostgreSQL from Python
**Action**: Remove PostgreSQL config from Python service if not used
- Check if actually used
- If not, remove config to reduce complexity

### Phase 2: Service Consolidation (Medium Risk, Medium Impact)

#### 2.1 Merge Analytics Services
**Action**: Combine heatmap, maturity, compliance into single analytics service

**Benefits**:
- Single database connection
- Shared caching
- Reduced code duplication
- Simpler API surface

**New Structure**:
```go
type AnalyticsService struct {
    db *gorm.DB
    // Handles:
    // - Heatmap generation
    // - Maturity scoring
    // - Compliance assessment
}
```

#### 2.2 Merge Organization Services
**Action**: Combine techStackService into organizationProfileService

**Benefits**:
- Tech stack is part of organization profile
- Single source of truth
- Reduced service count

### Phase 3: Integration (Higher Risk, Higher Impact)

#### 3.1 Integrate CPE Guesser
**Action**: Make CPE Guesser a library instead of separate service

**Benefits**:
- No HTTP overhead
- Direct function calls
- Shared connection pools
- Simpler deployment

**Implementation**:
- Keep CPE Guesser as library
- Import directly in enrichment service
- Remove HTTP server wrapper

#### 3.2 Unified Caching Strategy
**Action**: Implement shared caching interface

**Benefits**:
- Consistent cache behavior
- Easier debugging
- Better cache hit rates

## Simplified Architecture

### Before (Current)
```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Go API    │    │   Python    │    │  CPE        │
│             │    │ Enrichment  │    │  Guesser    │
│ - 9 Services│    │             │    │             │
│ - In-mem    │    │ - L1 Cache  │    │ - Valkey    │
│ - Valkey    │    │ - L2 Valkey │    │             │
└─────────────┘    │ - L3 Memcach│    └─────────────┘
                   └─────────────┘
```

### After (Optimized)
```
┌─────────────┐    ┌─────────────┐
│   Go API    │    │   Python    │
│             │    │ Enrichment  │
│ - 5 Services│    │             │
│ - In-mem    │    │ - In-mem    │
│ - Valkey         │ - Valkey    │
└─────────────┘    │ (CPE lib)   │
                   └─────────────┘
```

## Implementation Priority

### Immediate (This Week)
1. ✅ Remove Memcached
2. ⏳ Standardize on Valkey (replace all Redis with Valkey)
3. ✅ Remove unused PostgreSQL config from Python

### Short Term (This Month)
4. Merge analytics services (heatmap, maturity, compliance)
5. Merge organization services (techStack into organizationProfile)

### Long Term (Next Quarter)
6. Integrate CPE Guesser as library
7. Implement unified caching interface

## Expected Benefits

### Performance
- **Reduced latency**: Fewer cache layers = faster lookups
- **Lower memory**: Remove redundant caches
- **Better hit rates**: Unified caching strategy

### Maintainability
- **Less code**: Fewer services to maintain
- **Simpler debugging**: Fewer moving parts
- **Easier testing**: Smaller service surface area

### Operations
- **Fewer dependencies**: Remove Memcached
- **Simpler deployment**: Fewer services
- **Lower costs**: Less infrastructure

## Risk Assessment

| Change | Risk | Impact | Recommendation |
|--------|------|--------|----------------|
| Remove Memcached | Low | High | ✅ Do immediately |
| Standardize on Valkey | Low | High | ✅ Do immediately |
| Merge Analytics Services | Medium | High | ⚠️ Do after testing |
| Use Official CPE Guesser | Medium | Medium | ⚠️ Do after Phase 1 |

## Migration Strategy

1. **Phase 1**: Remove redundant components (low risk)
2. **Phase 2**: Consolidate services (test thoroughly)
3. **Phase 3**: Integrate components (requires refactoring)

Each phase should be:
- Tested in development
- Deployed to staging
- Monitored for issues
- Rolled back if problems occur

