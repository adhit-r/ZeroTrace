# Valkey Usage in ZeroTrace

## Overview

ZeroTrace uses **Valkey** (a Redis-compatible fork) as the primary in-memory data store. Valkey is fully compatible with Redis clients and protocols, providing the same functionality with active open-source development.

## Why Valkey?

- **Redis-compatible**: All Redis clients and libraries work with Valkey
- **Open-source**: Active community development
- **Performance**: Same performance characteristics as Redis
- **Protocol compatibility**: Uses the same RESP protocol

## Services Using Valkey

### 1. Go API (`api-go/`)
- **Cache**: `internal/cache/cache.go` - Unified caching layer
- **Job Queue**: `internal/jobs/job.go` - Background jobs using asynq
- **Queue Processor**: `internal/queue/processor.go` - Streams-based processing
- **Performance Optimizer**: `internal/optimization/performance.go` - Caching layer

**Client Library**: `github.com/redis/go-redis/v9` (works with Valkey)

### 2. Python Enrichment Service (`enrichment-python/`)
- **Cache Manager**: `app/cache_manager.py` - Unified cache with invalidation
- **Job Queue**: `app/jobs.py` - Background jobs using ARQ
- **CPE Guesser**: `cpe-guesser/lib/` - CPE dictionary storage (DB 8)
- **Ultra Optimized Enrichment**: `app/ultra_optimized_enrichment.py` - L2 cache

**Client Libraries**: 
- `redis[hiredis]>=5.2.0` - Synchronous Redis client
- `aioredis>=2.0.1` - Async Redis client

### 3. CPE Guesser
- **Storage**: Uses Valkey DB 8 for CPE dictionary indexing
- **Library**: `cpe-guesser/lib/cpeguesser_async.py` - Direct Valkey access

## Configuration

### Environment Variables

For compatibility, we use `REDIS_*` environment variables (Valkey is Redis-compatible):

```bash
# Go API
REDIS_HOST=valkey          # Service name in docker-compose
REDIS_PORT=6379
REDIS_PASSWORD=           # Optional
REDIS_DB=0                # Default DB

# Python Enrichment
REDIS_URL=redis://valkey:6379/0
REDIS_HOST=valkey
REDIS_PORT=6379
REDIS_DB=0

# CPE Guesser (uses DB 8)
CPE_GUESSER_VALKEY_HOST=valkey
CPE_GUESSER_VALKEY_PORT=6379
CPE_GUESSER_VALKEY_DB=8
```

### Docker Compose

```yaml
services:
  valkey:
    image: valkey/valkey:latest
    container_name: zerotrace-valkey
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "valkey-cli", "ping"]
```

## Database Allocation

Valkey uses multiple databases for data isolation:

- **DB 0**: Main cache and job queues (Go API, Python Enrichment)
- **DB 8**: CPE Guesser dictionary data (167,573+ CPE entries)

## Usage Patterns

### 1. Caching
- **L1 Cache**: In-memory (Python dict, Go map)
- **L2 Cache**: Valkey (distributed, persistent)
- **Cache Invalidation**: Pub/Sub pattern for distributed invalidation

### 2. Job Queues
- **Go**: asynq library (priority queues, retries, dead-letter queues)
- **Python**: ARQ library (async job processing)

### 3. Streams
- **Queue Processing**: Redis Streams for app data processing
- **Consumer Groups**: For parallel processing

### 4. CPE Dictionary
- **Storage**: Valkey DB 8
- **Indexing**: Keyword-based inverted index
- **Access**: Direct library access (no HTTP overhead)

## Migration Notes

All Redis clients work with Valkey without changes:
- `github.com/redis/go-redis/v9` ✅
- `redis[hiredis]` (Python) ✅
- `aioredis` (Python) ✅
- `asynq` (Go) ✅
- `ARQ` (Python) ✅

## Health Checks

```bash
# Check Valkey health
valkey-cli ping

# Check database sizes
valkey-cli -n 0 dbsize  # Main cache
valkey-cli -n 8 dbsize  # CPE Guesser

# Monitor connections
valkey-cli info clients
```

## Performance

Valkey provides the same performance as Redis:
- **Latency**: Sub-millisecond for cache operations
- **Throughput**: 100K+ ops/sec per instance
- **Memory**: Efficient data structures (same as Redis)

## References

- [Valkey GitHub](https://github.com/valkey-io/valkey)
- [Valkey Documentation](https://valkey.io/)
- [Redis Compatibility](https://valkey.io/docs/about/compatibility/)

