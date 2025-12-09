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
from .threat_intel import (
    mitre_attack_service,
    alienvault_otx_service,
    opencve_service,
    cisa_kev_service,
    initialize_threat_intel_services
)

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
        # Initialize threat intel services
        await initialize_threat_intel_services()
        
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
            
            # Final Fallback: Use NVD API if database is empty (works without API key, but rate-limited)
            if not vulnerabilities:
                try:
                    # Try to fetch from NVD API using CPE or product name
                    nvd_vulns = await self._get_vulns_from_nvd(name, version, cpe)
                    if nvd_vulns:
                        vulnerabilities = nvd_vulns
                        logger.info(f"Found {len(vulnerabilities)} CVEs from NVD API for {name} {version}")
                except Exception as e:
                    logger.debug(f"NVD API fallback failed: {e}")

            # 4. Enrich with Threat Intel Context (if vulnerabilities found)
            if vulnerabilities:
                item["vulnerabilities"] = await self._enrich_with_threat_intel(vulnerabilities)
            else:
                item["vulnerabilities"] = []
            
            # 5. Cache Result
            await cache_manager.set(cache_key, item["vulnerabilities"])
            
            item["source"] = "db" if item["vulnerabilities"] else ("nvd" if item["vulnerabilities"] else "none")
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
        try:
            # Optimized Strategy:
            # 1. Try 'cpe_embeddings' table first (fastest, relational match)
            # 2. Fallback to JSONB query on 'cves' table (slower, but comprehensive)
            
            # 1. CPE Embeddings Lookup
            sql = """
                SELECT c.id, 
                       c.description, 
                       c.data->'cve'->'metrics' as metrics
                FROM cves c
                JOIN cpe_embeddings ce ON c.id = ce.cve_id
                WHERE ce.cpe_string = $1
                LIMIT 20
            """
            
            rows = await db_manager.fetch_all(sql, cpe)
            
            # 2. Fallback: JSONB Search (only if no matches found via embeddings)
            if not rows:
                # Use JSONB existence check. 
                # Note: This assumes standard NVD 2.0 structure in 'data' column
                # Path: configurations -> nodes -> cpeMatch -> criteria
                # We use @> operator which needs GIN index for performance
                # Constructing a partial JSON to match against is complex due to array structure.
                # Alternative: EXISTS with jsonb_array_elements (requires lateral join or subquery)
                sql_jsonb = """
                    SELECT id, description, data->'cve'->'metrics' as metrics
                    FROM cves
                    WHERE EXISTS (
                        SELECT 1
                        FROM jsonb_array_elements(data->'cve'->'configurations') as configs,
                             jsonb_array_elements(configs->'nodes') as nodes,
                             jsonb_array_elements(nodes->'cpeMatch') as match
                        WHERE match->>'criteria' = $1
                    )
                    LIMIT 20
                """
                rows = await db_manager.fetch_all(sql_jsonb, cpe)

            results = []
            for row in rows:
                results.append({
                    "id": row.get("id", ""),
                    "description": row.get("description", ""),
                    "severity": self._extract_severity(row.get("metrics"))
                })
            
            if results:
                logger.info(f"Found {len(results)} vulnerabilities for CPE: {cpe}")
                
            return results
            
        except Exception as e:
            logger.error(f"Database CPE search failed: {e}")
            return []

    async def _get_vulns_by_text(self, name: str, version: str) -> List[Dict]:
        """Fallback: Search CVEs by text"""
        try:
            # Use Postgres Full Text Search
            # Fix: Use COALESCE to handle NULL descriptions and cast properly
            query = f"{name} {version}"
            sql = """
                SELECT id, 
                       COALESCE(description, '') as description, 
                       data->'cve'->'metrics' as metrics
                FROM cves
                WHERE to_tsvector('english', COALESCE(description, '')) @@ plainto_tsquery('english', $1::text)
                LIMIT 5
            """
            rows = await db_manager.fetch_all(sql, query)
            
            results = []
            for row in rows:
                results.append({
                    "id": row.get("id", ""),
                    "description": row.get("description", ""),
                    "severity": self._extract_severity(row.get("metrics"))
                })
            return results
        except Exception as e:
            logger.debug(f"Database text search failed: {e}")
            return []

    async def _enrich_with_threat_intel(self, vulnerabilities: List[Dict]) -> List[Dict]:
        """
        Enrich vulnerabilities with threat intelligence context
        Adds MITRE ATT&CK patterns, OTX IOCs, and OpenCVE data
        """
        enriched = []
        for vuln in vulnerabilities:
            cve_id = vuln.get("id")
            if not cve_id:
                enriched.append(vuln)
                continue
            
            # Create enriched vulnerability object
            enriched_vuln = vuln.copy()
            enriched_vuln["threat_intel"] = {}
            
            # Add CISA KEV (Known Exploited Vulnerabilities) - HIGH PRIORITY
            if cisa_kev_service.enabled:
                try:
                    kev_info = await cisa_kev_service.get_kev_info(cve_id)
                    if kev_info:
                        enriched_vuln["threat_intel"]["cisa_kev"] = kev_info
                        enriched_vuln["known_exploited"] = True
                        # Mark as high priority if in KEV
                        if "severity" not in enriched_vuln or enriched_vuln.get("severity") == "medium":
                            enriched_vuln["severity"] = "critical"  # KEV = critical priority
                except Exception as e:
                    logger.debug(f"Failed to fetch CISA KEV data for {cve_id}: {e}")
            
            # Add MITRE ATT&CK patterns
            if mitre_attack_service.enabled:
                try:
                    attack_patterns = await mitre_attack_service.get_attack_patterns(cve_id)
                    if attack_patterns:
                        enriched_vuln["threat_intel"]["mitre_attack"] = attack_patterns
                except Exception as e:
                    logger.debug(f"Failed to fetch MITRE ATT&CK data for {cve_id}: {e}")
            
            # Add AlienVault OTX IOCs
            if alienvault_otx_service.enabled:
                try:
                    iocs = await alienvault_otx_service.get_iocs(cve_id)
                    if iocs:
                        enriched_vuln["threat_intel"]["otx_iocs"] = iocs
                except Exception as e:
                    logger.debug(f"Failed to fetch OTX IOCs for {cve_id}: {e}")
            
            # Add OpenCVE enhanced context
            if opencve_service.enabled:
                try:
                    opencve_data = await opencve_service.get_cve_details(cve_id)
                    if opencve_data:
                        enriched_vuln["threat_intel"]["opencve"] = opencve_data
                except Exception as e:
                    logger.debug(f"Failed to fetch OpenCVE data for {cve_id}: {e}")
            
            enriched.append(enriched_vuln)
        
        return enriched

    def _extract_severity(self, metrics: Any) -> str:
        """Helper to extract severity from CVSS metrics"""
        try:
            if isinstance(metrics, str):
                metrics = orjson.loads(metrics)
            # Try V3.1
            return metrics.get("cvssMetricV31", [{}])[0].get("cvssData", {}).get("baseSeverity", "UNKNOWN")
        except:
            return "UNKNOWN"

    async def _get_vulns_from_nvd(self, name: str, version: str, cpe: Optional[str] = None) -> List[Dict]:
        """Fetch vulnerabilities from NVD API using keyword search"""
        try:
            # Search by keyword (product name)
            keyword = f"{name} {version}".strip()
            url = f"{settings.nvd_api_base_url}?keywordSearch={keyword}"
            headers = {"Content-Type": "application/json"}
            if settings.nvd_api_key:
                headers["apiKey"] = settings.nvd_api_key
            
            response = await self.http_client.get(url, headers=headers, timeout=10.0)
            response.raise_for_status()
            data = response.json()
            
            vulnerabilities = []
            for vuln in data.get("vulnerabilities", []):
                cve_data = vuln.get("cve", {})
                cve_id = cve_data.get("id", "")
                
                # Extract CVSS score
                metrics = cve_data.get("metrics", {})
                cvss_v31 = metrics.get("cvssMetricV31", [{}])[0] if metrics.get("cvssMetricV31") else {}
                cvss_v30 = metrics.get("cvssMetricV30", [{}])[0] if metrics.get("cvssMetricV30") else {}
                cvss_v2 = metrics.get("cvssMetricV2", [{}])[0] if metrics.get("cvssMetricV2") else {}
                
                cvss_data = cvss_v31.get("cvssData", {}) or cvss_v30.get("cvssData", {}) or cvss_v2.get("cvssData", {})
                base_score = cvss_data.get("baseScore", 0.0)
                severity = cvss_data.get("baseSeverity", "UNKNOWN")
                
                vulnerabilities.append({
                    "id": cve_id,
                    "cve_id": cve_id,
                    "description": cve_data.get("descriptions", [{}])[0].get("value", ""),
                    "severity": severity.lower() if severity else "unknown",
                    "cvss_score": float(base_score) if base_score else 0.0,
                    "source": "nvd"
                })
            
            return vulnerabilities[:10]  # Limit to top 10
        except Exception as e:
            logger.debug(f"NVD API search failed: {e}")
            return []

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

