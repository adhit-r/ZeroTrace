# Database Architecture Summary

## Overview
ZeroTrace uses multiple databases for different purposes across the microservices architecture.

## 1. Go API Service (`api-go`)

### PostgreSQL (Primary Database)
- **Purpose**: Main application database for all persistent data
- **ORM**: GORM (Go ORM)
- **Connection**: Via `gorm.io/driver/postgres`
- **Configuration**:
  - Host: `DB_HOST` (default: localhost)
  - Port: `DB_PORT` (default: 5432)
  - Database: `DB_NAME` (default: zerotrace)
  - User: `DB_USER` (default: postgres)
  - Password: `DB_PASSWORD`
  - SSL Mode: `DB_SSL_MODE` (default: disable)
- **Connection Pool**:
  - Max Idle Connections: 10
  - Max Open Connections: 100
  - Connection Max Lifetime: 1 hour
- **Stores**:
  - Users, Companies, Organizations
  - Agents and Agent Credentials
  - Scans and Scan Results
  - Vulnerabilities
  - Organization Profiles
  - Compliance Data
  - Maturity Assessments

### Valkey
- **Purpose**: Caching, session storage, rate limiting
- **Configuration**:
  - Host: `VALKEY_HOST` or `REDIS_HOST` (default: localhost) - using REDIS_* for compatibility
  - Port: `VALKEY_PORT` or `REDIS_PORT` (default: 6379)
  - Password: `VALKEY_PASSWORD` or `REDIS_PASSWORD`
  - Database: `VALKEY_DB` or `REDIS_DB` (default: 0)

## 2. Python Enrichment Service (`enrichment-python/app`)

### PostgreSQL (Optional)
- **Purpose**: Optional database for enrichment data storage
- **Configuration**:
  - Host: `DB_HOST` (default: localhost)
  - Port: `DB_PORT` (default: 5432)
  - Database: `DB_NAME` (default: zerotrace)
  - User: `DB_USER` (default: postgres)
  - Password: `DB_PASSWORD`
  - SSL Mode: `DB_SSL_MODE` (default: disable)
- **Note**: Currently configured but may not be actively used

### Valkey
- **Purpose**: Caching enrichment results
- **Configuration**:
  - Host: `VALKEY_HOST` or `REDIS_HOST` (default: localhost) - using REDIS_* for compatibility
  - Port: `VALKEY_PORT` or `REDIS_PORT` (default: 6379)
  - Password: `VALKEY_PASSWORD` or `REDIS_PASSWORD`
  - Database: `VALKEY_DB` or `REDIS_DB` (default: 0)

### Memcached (Optional)
- **Purpose**: Additional caching layer
- **Configuration**:
  - Host: `MEMCACHED_HOST` (default: localhost)
  - Port: `MEMCACHED_PORT` (default: 11211)

## 3. CPE Guesser Service (`enrichment-python/cpe-guesser`)

### Valkey
- **Purpose**: Stores CPE dictionary data with inverse indexing for fast keyword matching
- **What is Valkey?**: Open-source fork of Redis, fully Redis-compatible, now used everywhere
- **Configuration**:
  - Host: `CPE_GUESSER_VALKEY_HOST` (default: 127.0.0.1)
  - Port: `CPE_GUESSER_VALKEY_PORT` (default: 6379)
  - Password: `CPE_GUESSER_VALKEY_PASSWORD`
  - Database: `CPE_GUESSER_VALKEY_DB` (default: 8) - **Separate DB from main Redis**
  - Socket Timeout: 5.0s
  - Max Connections: 50
- **Data Stored**:
  - NVD CPE Dictionary (downloaded from NIST)
  - Inverse indexes for keyword matching
  - Ranked sets for popularity scoring
- **Initialization**: Must be seeded using `bin/import.py` before use
- **Why Separate DB?**: Uses database 8 to avoid conflicts with other Valkey usage

## Database Summary Table

| Service | Database | Purpose | Default Port | DB Number |
|---------|----------|---------|--------------|-----------|
| Go API | PostgreSQL | Primary application data | 5432 | N/A |
| Go API | Valkey | Caching, sessions | 6379 | 0 |
| Python Enrichment | PostgreSQL | Optional data storage | 5432 | N/A |
| Python Enrichment | Valkey | Caching | 6379 | 0 |
| CPE Guesser | Valkey | CPE dictionary storage | 6379 | 8 |

## Environment Variables

### Go API
```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=zerotrace
DB_USER=postgres
DB_PASSWORD=your_password
DB_SSL_MODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### Python Enrichment
```bash
DB_HOST=localhost
DB_PORT=5432
DB_NAME=zerotrace
DB_USER=postgres
DB_PASSWORD=your_password

VALKEY_HOST=localhost  # or REDIS_HOST for compatibility
VALKEY_PORT=6379       # or REDIS_PORT for compatibility
VALKEY_PASSWORD=        # or REDIS_PASSWORD for compatibility
VALKEY_DB=0            # or REDIS_DB for compatibility
```

### CPE Guesser
```bash
CPE_GUESSER_VALKEY_HOST=127.0.0.1
CPE_GUESSER_VALKEY_PORT=6379
CPE_GUESSER_VALKEY_PASSWORD=
CPE_GUESSER_VALKEY_DB=8
```

## Notes

1. **Valkey Standardization**: All services now use Valkey instead of Redis. Valkey is Redis-compatible, so existing Redis clients work without changes.

2. **Database Separation**: CPE Guesser uses Valkey DB 8 to avoid conflicts with other services using DB 0.

3. **PostgreSQL**: Primary database for all persistent application data. All services can share the same PostgreSQL instance with different databases or schemas.

4. **Valkey**: Used for caching and can be shared across services using different database numbers (0 for main, 8 for CPE Guesser). Environment variables can use `REDIS_*` prefix for backward compatibility with existing Redis clients.

