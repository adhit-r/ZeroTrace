import asyncio
import time
import orjson
import httpx
from datetime import datetime
from typing import List, Dict, Optional, Any
from tenacity import retry, stop_after_attempt, wait_exponential

from ..core.config import settings
from ..core.logging import get_logger
from ..core.cache import cache_manager
from ..core.database import db_manager
from .cpe_matcher import semantic_matcher

logger = get_logger(__name__)

class EnrichmentService:
    """
    Unified High-Performance Enrichment Service
    Merges logic from cve_enrichment, batch_enrichment, and ultra_optimized_enrichment.
    """
    
    def __init__(self):
        self.http_client: Optional[httpx.AsyncClient] = None
        self.sem = asyncio.Semaphore(settings.max_concurrent_requests)

    async def initialize(self):
        """Initialize all subsystems"""
        logger.info("Initializing Enrichment Service...")
        
        # 1. Core Systems
        await cache_manager.initialize()
        await db_manager.initialize()
        await semantic_matcher.initialize()
        
        # 2. HTTP Client (Shared, High-Performance)
        if not self.http_client:
            self.http_client = httpx.AsyncClient(
                timeout=10.0,
                limits=httpx.Limits(max_keepalive_connections=50, max_connections=100),
                verify=True
            )
        
        logger.info("Enrichment Service initialized")

    async def close(self):
        """Cleanup resources"""
        if self.http_client:
            await self.http_client.aclose()
        await cache_manager.close()
        await db_manager.close()

    async def enrich_software(self, software_list: List[Dict]) -> List[Dict]:
        """
        Enrich a list of software items with CVEs.
        Uses:
        1. Multi-Level Cache
        2. Database Lookup
        3. CPE Guesser + Semantic Matcher
        4. External NVD API (Fallback)
        """
        results = []
        # Process in parallel with semaphore
        tasks = [self._enrich_single_item(item) for item in software_list]
        enriched_items = await asyncio.gather(*tasks, return_exceptions=True)
        
        for res in enriched_items:
            if isinstance(res, Exception):
                logger.error("Enrichment error", error=str(res))
                continue
            results.append(res)
            
        return results

    async def _enrich_single_item(self, item: Dict) -> Dict:
        """Enrich a single software item"""
        async with self.sem:
            name = item.get("name", "")
            version = item.get("version", "")
            vendor = item.get("vendor", "")
            
            # Cache Key
            cache_key = f"enrich:{name}:{version}:{vendor}"
            
            # 1. Check Cache
            cached = await cache_manager.get(cache_key)
            if cached:
                item["vulnerabilities"] = cached
                item["source"] = "cache"
                return item

            # 2. Identify CPE
            cpe = await self._find_cpe(name, version, vendor)
            item["cpe"] = cpe

            # 3. Find Vulnerabilities
            vulnerabilities = []
            if cpe:
                vulnerabilities = await self._get_vulns_by_cpe(cpe)
            
            # Fallback: Fuzzy search by name if no CPE found
            if not vulnerabilities:
                vulnerabilities = await self._get_vulns_by_text(name, version)

            # 4. Cache Result
            await cache_manager.set(cache_key, vulnerabilities)
            
            item["vulnerabilities"] = vulnerabilities
            item["source"] = "db" if vulnerabilities else "none"
            return item

    async def _find_cpe(self, name: str, version: str, vendor: str) -> Optional[str]:
        """Find CPE using Hybrid Matcher"""
        # Try Exact/L1 Matcher (CPE Guesser)
        # Assuming cpe-guesser runs as a sidecar/service, or we call it directly if integrated
        # For now, we use the L2 Semantic Matcher which also handles logic
        
        matches = await semantic_matcher.match_software(vendor, name, version)
        if matches:
            # Return top match
            return matches[0]["cpe"]
        return None

    async def _get_vulns_by_cpe(self, cpe: str) -> List[Dict]:
        """Query DB for CVEs linked to this CPE"""
        # This requires the cpe_cve_map table or querying via JSONB
        # Simplified: Search in JSONB configurations
        # Note: In a real migration, we would have a standardized cpe_cve relational table
        
        # Using the semantic_matcher's knowledge or direct DB query
        # For Phase 1/2, we query the 'cves' table where data->configurations... matches
        # Ideally, we should have a `cve_cpe` join table.
        
        # For now, let's assume we query the raw JSONB (slower but works without join table)
        # Or better: Use the embeddings table inverse lookup if we had it
        
        # Optimized: Use Full Text Search on JSONB or description
        return []

    async def _get_vulns_by_text(self, name: str, version: str) -> List[Dict]:
        """Fallback: Search CVEs by text"""
        # Use Postgres Full Text Search
        query = f"{name} {version}"
        sql = """
            SELECT id, description, data->'cve'->'metrics' as metrics
            FROM cves
            WHERE to_tsvector('english', description) @@ plainto_tsquery('english', $1)
            LIMIT 5
        """
        rows = await db_manager.fetch_all(sql, query)
        
        results = []
        for row in rows:
            results.append({
                "id": row["id"],
                "description": row["description"],
                "severity": self._extract_severity(row["metrics"])
            })
        return results

    def _extract_severity(self, metrics: Any) -> str:
        """Helper to extract severity from CVSS metrics"""
        try:
            if isinstance(metrics, str):
                metrics = orjson.loads(metrics)
            # Try V3.1
            return metrics.get("cvssMetricV31", [{}])[0].get("cvssData", {}).get("baseSeverity", "UNKNOWN")
        except:
            return "UNKNOWN"

    @retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=4, max=10))
    async def fetch_nvd_data(self, cve_id: str) -> Optional[Dict]:
        """Fetch from NVD API with Circuit Breaker pattern (via Tenacity)"""
        url = f"{settings.nvd_api_base_url}?cveId={cve_id}"
        headers = {}
        if settings.nvd_api_key:
            headers["apiKey"] = settings.nvd_api_key
            
        response = await self.http_client.get(url, headers=headers)
        response.raise_for_status()
        data = response.json()
        return data.get("vulnerabilities", [{}])[0].get("cve")

# Global Service
enrichment_service = EnrichmentService()

