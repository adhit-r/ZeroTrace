"""
CPE Guesser Client Library
Async client wrapper for CPE Guesser service with caching and fallback support
"""

import asyncio
import logging
from typing import List, Dict, Optional, Tuple
from collections import OrderedDict
import time
import httpx
from .core.config import settings

logger = logging.getLogger(__name__)


class LRUCache:
    """LRU Cache for CPE results"""
    def __init__(self, max_size: int = 10000, ttl: int = 3600):
        self.max_size = max_size
        self.ttl = ttl
        self.cache = OrderedDict()
        self.timestamps = {}
    
    def get(self, key: str) -> Optional[List[Tuple[Optional[int], str]]]:
        """Get item from cache if not expired"""
        if key not in self.cache:
            return None
        
        # Check if expired
        if time.time() - self.timestamps.get(key, 0) > self.ttl:
            self._remove(key)
            return None
        
        # Move to end (most recently used)
        self.cache.move_to_end(key)
        return self.cache[key]
    
    def set(self, key: str, value: List[Tuple[Optional[int], str]]):
        """Set item in cache with LRU eviction"""
        if key in self.cache:
            self.cache.move_to_end(key)
        else:
            # Evict oldest if at capacity
            if len(self.cache) >= self.max_size:
                oldest_key = next(iter(self.cache))
                self._remove(oldest_key)
        
        self.cache[key] = value
        self.timestamps[key] = time.time()
    
    def _remove(self, key: str):
        """Remove item from cache"""
        self.cache.pop(key, None)
        self.timestamps.pop(key, None)


class CPEGuesserClient:
    """
    Async client for CPE Guesser service.
    
    Supports both HTTP API mode and direct Valkey access mode.
    Includes caching, retry logic, and fallback support.
    """
    
    def __init__(
        self,
        base_url: Optional[str] = None,
        timeout: float = 5.0,
        max_retries: int = 3,
        cache_enabled: bool = True,
        cache_ttl: int = 3600,
        cache_max_size: int = 10000
    ):
        """
        Initialize CPE Guesser client.
        
        Args:
            base_url: Base URL for CPE Guesser API (defaults to config)
            timeout: Request timeout in seconds
            max_retries: Maximum number of retry attempts
            cache_enabled: Enable result caching
            cache_ttl: Cache TTL in seconds
            cache_max_size: Maximum cache size
        """
        self.base_url = base_url or getattr(settings, 'cpe_guesser_url', 'http://localhost:8000')
        self.timeout = timeout
        self.max_retries = max_retries
        self.cache_enabled = cache_enabled
        self.cache = LRUCache(max_size=cache_max_size, ttl=cache_ttl) if cache_enabled else None
        self.http_client: Optional[httpx.AsyncClient] = None
        self._enabled = getattr(settings, 'cpe_guesser_enabled', True)
    
    async def _ensure_client(self):
        """Ensure HTTP client is initialized"""
        if not self.http_client:
            self.http_client = httpx.AsyncClient(
                base_url=self.base_url,
                timeout=self.timeout,
                limits=httpx.Limits(max_keepalive_connections=20, max_connections=100)
            )
    
    async def close(self):
        """Close HTTP client"""
        if self.http_client:
            await self.http_client.aclose()
            self.http_client = None
    
    def _cache_key(self, words: List[str]) -> str:
        """Generate cache key from words"""
        return "|".join(sorted(word.lower() for word in words))
    
    async def guess_cpe(
        self,
        words: List[str],
        use_cache: bool = True
    ) -> List[Tuple[Optional[int], str]]:
        """
        Guess CPE identifiers from keywords.
        
        Args:
            words: List of keywords to search for
            use_cache: Whether to use cache (default: True)
        
        Returns:
            List of tuples (rank, cpe_string) sorted by rank
        """
        if not self._enabled:
            logger.debug("CPE Guesser is disabled")
            return []
        
        if not words:
            return []
        
        # Check cache
        if use_cache and self.cache:
            cache_key = self._cache_key(words)
            cached_result = self.cache.get(cache_key)
            if cached_result is not None:
                logger.debug(f"Cache hit for CPE search: {words}")
                return cached_result
        
        # Make API request with retry
        for attempt in range(self.max_retries + 1):
            try:
                await self._ensure_client()
                
                response = await self.http_client.post(
                    "/search",
                    json={"query": words},
                    headers={"Content-Type": "application/json"}
                )
                response.raise_for_status()
                
                data = response.json()
                
                # Parse response (new format)
                if isinstance(data, dict) and "results" in data:
                    results = [
                        (r.get("rank"), r.get("cpe"))
                        for r in data["results"]
                    ]
                # Legacy format: [[rank, cpe], ...]
                elif isinstance(data, list):
                    results = [(r[0] if len(r) > 0 else None, r[1] if len(r) > 1 else "") for r in data]
                else:
                    logger.warning(f"Unexpected response format: {data}")
                    results = []
                
                # Cache result
                if use_cache and self.cache:
                    cache_key = self._cache_key(words)
                    self.cache.set(cache_key, results)
                
                logger.debug(f"CPE search completed: {len(results)} results for {words}")
                return results
                
            except httpx.HTTPStatusError as e:
                if e.response.status_code == 404:
                    logger.warning(f"CPE Guesser endpoint not found: {e}")
                    return []
                if attempt < self.max_retries:
                    wait_time = 0.5 * (2 ** attempt)
                    logger.warning(f"Retry {attempt + 1}/{self.max_retries} after {wait_time}s")
                    await asyncio.sleep(wait_time)
                    continue
                logger.error(f"HTTP error in CPE search: {e}")
                return []
            except httpx.RequestError as e:
                if attempt < self.max_retries:
                    wait_time = 0.5 * (2 ** attempt)
                    logger.warning(f"Request error, retry {attempt + 1}/{self.max_retries} after {wait_time}s: {e}")
                    await asyncio.sleep(wait_time)
                    continue
                logger.error(f"Request error in CPE search: {e}")
                return []
            except Exception as e:
                logger.error(f"Unexpected error in CPE search: {e}", exc_info=True)
                return []
        
        return []
    
    async def get_unique_cpe(self, words: List[str]) -> Optional[str]:
        """
        Get the best matching CPE identifier.
        
        Args:
            words: List of keywords to search for
        
        Returns:
            Best matching CPE string or None
        """
        results = await self.guess_cpe(words)
        if results and len(results) > 0:
            return results[0][1]
        return None
    
    async def health_check(self) -> bool:
        """
        Check if CPE Guesser service is healthy.
        
        Returns:
            True if service is healthy, False otherwise
        """
        if not self._enabled:
            return False
        
        try:
            await self._ensure_client()
            response = await self.http_client.get("/health", timeout=2.0)
            return response.status_code == 200
        except Exception as e:
            logger.debug(f"CPE Guesser health check failed: {e}")
            return False


# Global client instance
_cpe_guesser_client: Optional[CPEGuesserClient] = None


async def get_cpe_guesser_client() -> CPEGuesserClient:
    """Get or create global CPE Guesser client instance"""
    global _cpe_guesser_client
    
    if _cpe_guesser_client is None:
        _cpe_guesser_client = CPEGuesserClient(
            base_url=getattr(settings, 'cpe_guesser_url', None),
            timeout=getattr(settings, 'cpe_guesser_timeout', 5.0),
            cache_enabled=getattr(settings, 'cpe_guesser_cache_enabled', True),
            cache_ttl=getattr(settings, 'cpe_guesser_cache_ttl', 3600)
        )
    
    return _cpe_guesser_client

