import asyncpg
import asyncio
from typing import Optional, List, Dict, Any
from .config import settings
from .logging import get_logger

logger = get_logger(__name__)

class DatabaseManager:
    """
    Centralized Database Connection Pool
    Uses asyncpg for high-performance PostgreSQL interaction.
    """
    
    def __init__(self):
        self.pool: Optional[asyncpg.Pool] = None

    async def initialize(self):
        """Initialize the connection pool"""
        if self.pool:
            return

        try:
            self.pool = await asyncpg.create_pool(
                dsn=settings.get_database_dsn(),
                min_size=settings.db_pool_min_size,
                max_size=settings.db_pool_max_size,
                command_timeout=60,
                # Configure jsonb handling automatically
                init=self._init_connection
            )
            logger.info("Database pool initialized", 
                       min_size=settings.db_pool_min_size, 
                       max_size=settings.db_pool_max_size)
            
            # Initialize schema and extensions
            await self._init_schema()
            
        except Exception as e:
            logger.critical("Failed to initialize database pool", error=str(e))
            raise

    async def _init_connection(self, conn):
        """Setup connection state"""
        # Set codec for jsonb if needed, though asyncpg handles it well natively
        await conn.set_type_codec(
            'jsonb',
            encoder=None, # handled by asyncpg
            decoder=None,
            schema='pg_catalog',
            format='text'
        )

    async def _init_schema(self):
        """Initialize database schema and extensions"""
        async with self.pool.acquire() as conn:
            # 1. Enable pgvector extension (only if available)
            if settings.pgvector_enabled:
                try:
                    await conn.execute("CREATE EXTENSION IF NOT EXISTS vector;")
                except Exception as e:
                    logger.warning("pgvector extension not available, continuing without it", error=str(e))
                    # Disable pgvector for this session
                    settings.pgvector_enabled = False
            
            # 2. Create CVE Data table with JSONB support
            await conn.execute("""
                CREATE TABLE IF NOT EXISTS cves (
                    id TEXT PRIMARY KEY,
                    data JSONB NOT NULL,
                    description TEXT,
                    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
                );
                
                CREATE INDEX IF NOT EXISTS idx_cves_data ON cves USING gin (data);
            """)

            # 3. Create Embeddings table for Vector Search (if enabled)
            if settings.pgvector_enabled:
                await conn.execute("""
                    CREATE TABLE IF NOT EXISTS cpe_embeddings (
                        id SERIAL PRIMARY KEY,
                        cve_id TEXT REFERENCES cves(id) ON DELETE CASCADE,
                        cpe_string TEXT NOT NULL,
                        embedding vector(384), -- Dimension for all-MiniLM-L6-v2
                        UNIQUE(cve_id, cpe_string)
                    );
                    
                    -- HNSW Index for fast similarity search
                    CREATE INDEX IF NOT EXISTS idx_cpe_embeddings_vec 
                    ON cpe_embeddings 
                    USING hnsw (embedding vector_cosine_ops);
                """)

    async def close(self):
        """Close the connection pool"""
        if self.pool:
            await self.pool.close()
            logger.info("Database pool closed")

    async def fetch_one(self, query: str, *args) -> Optional[asyncpg.Record]:
        async with self.pool.acquire() as conn:
            return await conn.fetchrow(query, *args)

    async def fetch_all(self, query: str, *args) -> List[asyncpg.Record]:
        async with self.pool.acquire() as conn:
            return await conn.fetch(query, *args)

    async def execute(self, query: str, *args) -> str:
        async with self.pool.acquire() as conn:
            return await conn.execute(query, *args)
            
    async def executemany(self, query: str, args: List[Any]):
        async with self.pool.acquire() as conn:
            return await conn.executemany(query, args)

# Global database instance
db_manager = DatabaseManager()

