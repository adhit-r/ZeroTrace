import json
import asyncio
from typing import Optional, Any, Union
from datetime import timedelta
from cachetools import TTLCache
import redis.asyncio as redis
from .config import settings
from .logging import get_logger

logger = get_logger(__name__)

class CacheManager:
    """
    Multi-Level Cache Strategy (L1 Memory + L2 Redis)
    
    L1: In-memory TTLCache for ultra-fast access (microsecond latency)
    L2: Redis for distributed caching (millisecond latency)
    """
    
    def __init__(self):
        # L1 Cache: In-memory, process-local
        # Holds 1000 items, expires in 5 minutes by default
        self.l1_cache = TTLCache(maxsize=1000, ttl=300)
        
        # L2 Cache: Redis, distributed
        self.redis_client: Optional[redis.Redis] = None
        self.enabled = settings.cache_enabled
        
    async def initialize(self):
        """Initialize Redis connection"""
        if not self.enabled:
            return

        try:
            self.redis_client = redis.from_url(
                settings.get_redis_dsn(),
                decode_responses=True,
                socket_timeout=5.0
            )
            await self.redis_client.ping()
            logger.info("CacheManager initialized", type="redis", host=settings.redis_host)
        except Exception as e:
            logger.error("Failed to connect to Redis cache", error=str(e))
            self.redis_client = None

    async def get(self, key: str) -> Optional[Any]:
        """Get item from cache (L1 -> L2)"""
        if not self.enabled:
            return None

        # Check L1
        if key in self.l1_cache:
            # logger.debug("L1 Cache Hit", key=key)
            return self.l1_cache[key]

        # Check L2
        if self.redis_client:
            try:
                data = await self.redis_client.get(key)
                if data:
                    # logger.debug("L2 Cache Hit", key=key)
                    # Deserialize and populate L1
                    try:
                        value = json.loads(data)
                        self.l1_cache[key] = value
                        return value
                    except json.JSONDecodeError:
                        return data # Return string if not valid JSON
            except Exception as e:
                logger.warning("L2 Cache error", error=str(e))
        
        return None

    async def set(self, key: str, value: Any, ttl: int = None):
        """Set item in cache (L1 + L2)"""
        if not self.enabled:
            return

        expiration = ttl or settings.cache_ttl

        # Set L1
        self.l1_cache[key] = value

        # Set L2
        if self.redis_client:
            try:
                # Serialize complex objects
                if isinstance(value, (dict, list)):
                    serialized = json.dumps(value)
                else:
                    serialized = str(value)
                
                await self.redis_client.setex(key, expiration, serialized)
            except Exception as e:
                logger.warning("L2 Cache set error", error=str(e))

    async def delete(self, key: str):
        """Delete item from cache"""
        if key in self.l1_cache:
            del self.l1_cache[key]
        
        if self.redis_client:
            try:
                await self.redis_client.delete(key)
            except Exception as e:
                logger.warning("L2 Cache delete error", error=str(e))

    async def close(self):
        """Close connections"""
        if self.redis_client:
            await self.redis_client.close()

# Global cache instance
cache_manager = CacheManager()

