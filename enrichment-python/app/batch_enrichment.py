import asyncio
import aiohttp
import json
import logging
import time
from typing import List, Dict, Optional
from dataclasses import dataclass, asdict
from datetime import datetime, timedelta
from fastapi import HTTPException
import redis.asyncio as redis

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class AppData:
    id: str
    company_id: str
    agent_id: str
    app_name: str
    app_version: str
    package_type: str
    architecture: Optional[str] = None

@dataclass
class Vulnerability:
    cve_id: str
    severity: str
    cvss_score: float
    title: str
    description: str
    references: List[str] = None

@dataclass
class EnrichmentResult:
    app_id: str
    company_id: str
    agent_id: str
    vulnerabilities: List[Vulnerability]
    enriched_at: datetime
    error: Optional[str] = None

class BatchEnrichmentService:
    def __init__(self):
        self.session: Optional[aiohttp.ClientSession] = None
        self.redis_client: Optional[redis.Redis] = None
        self.batch_size = 100
        self.max_concurrent = 10
        self.cache_ttl = 3600  # 1 hour
        self.rate_limit_delay = 0.01  # 10ms between requests
        
        # NVD API configuration
        self.nvd_api_key = None  # Set from environment
        self.nvd_base_url = "https://services.nvd.nist.gov/rest/json/cves/2.0"
        
        # CVE Search API configuration
        self.cve_search_url = "https://cve.circl.lu/api/search"
        
    async def initialize(self):
        """Initialize the service"""
        # Create HTTP session
        timeout = aiohttp.ClientTimeout(total=30)
        self.session = aiohttp.ClientSession(timeout=timeout)
        
        # Initialize Redis connection
        self.redis_client = redis.Redis(
            host='localhost',
            port=6379,
            db=0,
            decode_responses=True
        )
        
        logger.info("BatchEnrichmentService initialized")
    
    async def cleanup(self):
        """Cleanup resources"""
        if self.session:
            await self.session.close()
        if self.redis_client:
            await self.redis_client.close()
    
    async def process_batch(self, apps: List[AppData]) -> List[EnrichmentResult]:
        """Process a batch of apps efficiently"""
        start_time = time.time()
        logger.info(f"Processing batch of {len(apps)} apps")
        
        try:
            # 1. Check cache first
            cached_results = await self.get_cached_results(apps)
            uncached_apps = [app for app in apps if app.id not in cached_results]
            
            logger.info(f"Found {len(cached_results)} cached results, {len(uncached_apps)} need enrichment")
            
            # 2. Process uncached apps in parallel
            if uncached_apps:
                enriched_results = await self.enrich_apps_parallel(uncached_apps)
                
                # 3. Cache results
                await self.cache_results(enriched_results)
                
                # 4. Combine results
                all_results = {**cached_results, **enriched_results}
            else:
                all_results = cached_results
            
            # 5. Convert to list and sort by app_id
            results = [all_results[app.id] for app in apps]
            
            duration = time.time() - start_time
            logger.info(f"Batch processing completed in {duration:.2f}s")
            
            return results
            
        except Exception as e:
            logger.error(f"Batch processing failed: {e}")
            raise HTTPException(status_code=500, detail=f"Enrichment failed: {str(e)}")
    
    async def get_cached_results(self, apps: List[AppData]) -> Dict[str, EnrichmentResult]:
        """Get cached enrichment results"""
        cached_results = {}
        
        if not self.redis_client:
            return cached_results
        
        try:
            # Get cache keys for all apps
            cache_keys = [f"cve:app:{app.id}" for app in apps]
            
            # Get cached data in batch
            cached_data = await self.redis_client.mget(cache_keys)
            
            for app, cached_json in zip(apps, cached_data):
                if cached_json:
                    try:
                        data = json.loads(cached_json)
                        
                        # Convert back to EnrichmentResult
                        vulnerabilities = [
                            Vulnerability(**vuln) for vuln in data.get('vulnerabilities', [])
                        ]
                        
                        result = EnrichmentResult(
                            app_id=data['app_id'],
                            company_id=data['company_id'],
                            agent_id=data['agent_id'],
                            vulnerabilities=vulnerabilities,
                            enriched_at=datetime.fromisoformat(data['enriched_at']),
                            error=data.get('error')
                        )
                        
                        cached_results[app.id] = result
                        
                    except (json.JSONDecodeError, KeyError) as e:
                        logger.warning(f"Failed to parse cached data for app {app.id}: {e}")
                        continue
        
        except Exception as e:
            logger.error(f"Failed to get cached results: {e}")
        
        return cached_results
    
    async def cache_results(self, results: Dict[str, EnrichmentResult]):
        """Cache enrichment results"""
        if not self.redis_client:
            return
        
        try:
            # Prepare cache data
            cache_data = {}
            for app_id, result in results.items():
                # Convert to dict for JSON serialization
                data = {
                    'app_id': result.app_id,
                    'company_id': result.company_id,
                    'agent_id': result.agent_id,
                    'vulnerabilities': [asdict(v) for v in result.vulnerabilities],
                    'enriched_at': result.enriched_at.isoformat(),
                    'error': result.error
                }
                
                cache_data[f"cve:app:{app_id}"] = json.dumps(data)
            
            # Set cache with TTL
            pipeline = self.redis_client.pipeline()
            for key, value in cache_data.items():
                pipeline.setex(key, self.cache_ttl, value)
            
            await pipeline.execute()
            logger.info(f"Cached {len(cache_data)} results")
            
        except Exception as e:
            logger.error(f"Failed to cache results: {e}")
    
    async def enrich_apps_parallel(self, apps: List[AppData]) -> Dict[str, EnrichmentResult]:
        """Enrich apps in parallel with rate limiting"""
        semaphore = asyncio.Semaphore(self.max_concurrent)
        
        async def enrich_single(app: AppData) -> tuple:
            async with semaphore:
                try:
                    # Rate limiting
                    await asyncio.sleep(self.rate_limit_delay)
                    
                    # Get CVE data
                    vulnerabilities = await self.get_cve_data(app)
                    
                    result = EnrichmentResult(
                        app_id=app.id,
                        company_id=app.company_id,
                        agent_id=app.agent_id,
                        vulnerabilities=vulnerabilities,
                        enriched_at=datetime.utcnow()
                    )
                    
                    return (app.id, result)
                    
                except Exception as e:
                    logger.error(f"Failed to enrich app {app.id}: {e}")
                    result = EnrichmentResult(
                        app_id=app.id,
                        company_id=app.company_id,
                        agent_id=app.agent_id,
                        vulnerabilities=[],
                        enriched_at=datetime.utcnow(),
                        error=str(e)
                    )
                    return (app.id, result)
        
        # Process all apps concurrently
        tasks = [enrich_single(app) for app in apps]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        # Convert to dict, handling exceptions
        enriched_results = {}
        for result in results:
            if isinstance(result, Exception):
                logger.error(f"Task failed with exception: {result}")
                continue
            app_id, enrichment_result = result
            enriched_results[app_id] = enrichment_result
        
        return enriched_results
    
    async def get_cve_data(self, app: AppData) -> List[Vulnerability]:
        """Get CVE data from multiple sources"""
        vulnerabilities = []
        
        try:
            # 1. NVD API
            nvd_vulns = await self.get_nvd_cves(app)
            vulnerabilities.extend(nvd_vulns)
            
            # 2. CVE Search API
            cve_search_vulns = await self.get_cve_search_cves(app)
            vulnerabilities.extend(cve_search_vulns)
            
            # 3. Remove duplicates based on CVE ID
            unique_vulns = self.deduplicate_vulnerabilities(vulnerabilities)
            
            return unique_vulns
            
        except Exception as e:
            logger.error(f"Failed to get CVE data for {app.app_name} {app.app_version}: {e}")
            return []
    
    async def get_nvd_cves(self, app: AppData) -> List[Vulnerability]:
        """Get CVE data from NVD API"""
        try:
            # Build search query
            query = f"{app.app_name} {app.app_version}"
            
            params = {
                'keywordSearch': query,
                'resultsPerPage': 20
            }
            
            if self.nvd_api_key:
                params['apiKey'] = self.nvd_api_key
            
            async with self.session.get(self.nvd_base_url, params=params) as response:
                if response.status != 200:
                    logger.warning(f"NVD API returned status {response.status}")
                    return []
                
                data = await response.json()
                
                vulnerabilities = []
                for vuln in data.get('vulnerabilities', []):
                    cve = vuln.get('cve', {})
                    
                    # Extract CVSS score
                    cvss_score = 0.0
                    if 'metrics' in cve:
                        cvss_v3 = cve['metrics'].get('cvssMetricV31', [{}])[0]
                        if 'cvssData' in cvss_v3:
                            cvss_score = cvss_v3['cvssData'].get('baseScore', 0.0)
                    
                    # Determine severity
                    severity = self.get_severity(cvss_score)
                    
                    vulnerability = Vulnerability(
                        cve_id=cve.get('id', ''),
                        severity=severity,
                        cvss_score=cvss_score,
                        title=cve.get('descriptions', [{}])[0].get('value', ''),
                        description=cve.get('descriptions', [{}])[0].get('value', ''),
                        references=[ref.get('url', '') for ref in cve.get('references', [])]
                    )
                    
                    vulnerabilities.append(vulnerability)
                
                return vulnerabilities
                
        except Exception as e:
            logger.error(f"Failed to get NVD CVEs: {e}")
            return []
    
    async def get_cve_search_cves(self, app: AppData) -> List[Vulnerability]:
        """Get CVE data from CVE Search API"""
        try:
            # Build search query
            query = f"{app.app_name}:{app.app_version}"
            
            params = {'q': query}
            
            async with self.session.get(self.cve_search_url, params=params) as response:
                if response.status != 200:
                    logger.warning(f"CVE Search API returned status {response.status}")
                    return []
                
                data = await response.json()
                
                vulnerabilities = []
                for vuln in data:
                    # Extract CVSS score
                    cvss_score = float(vuln.get('cvss', 0))
                    severity = self.get_severity(cvss_score)
                    
                    vulnerability = Vulnerability(
                        cve_id=vuln.get('id', ''),
                        severity=severity,
                        cvss_score=cvss_score,
                        title=vuln.get('summary', ''),
                        description=vuln.get('summary', ''),
                        references=vuln.get('references', [])
                    )
                    
                    vulnerabilities.append(vulnerability)
                
                return vulnerabilities
                
        except Exception as e:
            logger.error(f"Failed to get CVE Search CVEs: {e}")
            return []
    
    def get_severity(self, cvss_score: float) -> str:
        """Convert CVSS score to severity level"""
        if cvss_score >= 9.0:
            return 'critical'
        elif cvss_score >= 7.0:
            return 'high'
        elif cvss_score >= 4.0:
            return 'medium'
        elif cvss_score >= 0.1:
            return 'low'
        else:
            return 'info'
    
    def deduplicate_vulnerabilities(self, vulnerabilities: List[Vulnerability]) -> List[Vulnerability]:
        """Remove duplicate vulnerabilities based on CVE ID"""
        seen = set()
        unique_vulns = []
        
        for vuln in vulnerabilities:
            if vuln.cve_id not in seen:
                seen.add(vuln.cve_id)
                unique_vulns.append(vuln)
        
        return unique_vulns

# Global service instance
enrichment_service = BatchEnrichmentService()

# FastAPI endpoints
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

router = APIRouter()

class BatchEnrichmentRequest(BaseModel):
    apps: List[AppData]

class BatchEnrichmentResponse(BaseModel):
    processed: int
    results: List[EnrichmentResult]
    timestamp: str
    duration_ms: float

@router.post("/enrich/batch", response_model=BatchEnrichmentResponse)
async def enrich_batch(request: BatchEnrichmentRequest):
    """Process batch of apps efficiently"""
    start_time = time.time()
    
    try:
        # Initialize service if needed
        if not enrichment_service.session:
            await enrichment_service.initialize()
        
        # Process batch
        results = await enrichment_service.process_batch(request.apps)
        
        duration_ms = (time.time() - start_time) * 1000
        
        return BatchEnrichmentResponse(
            processed=len(results),
            results=results,
            timestamp=datetime.utcnow().isoformat(),
            duration_ms=duration_ms
        )
        
    except Exception as e:
        logger.error(f"Batch enrichment failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/enrich/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "timestamp": datetime.utcnow().isoformat(),
        "service": "batch_enrichment"
    }

@router.get("/enrich/stats")
async def get_stats():
    """Get enrichment statistics"""
    try:
        if enrichment_service.redis_client:
            # Get cache statistics
            cache_stats = await enrichment_service.redis_client.info('memory')
            return {
                "cache_stats": cache_stats,
                "batch_size": enrichment_service.batch_size,
                "max_concurrent": enrichment_service.max_concurrent,
                "timestamp": datetime.utcnow().isoformat()
            }
        else:
            return {
                "error": "Redis not connected",
                "timestamp": datetime.utcnow().isoformat()
            }
    except Exception as e:
        return {
            "error": str(e),
            "timestamp": datetime.utcnow().isoformat()
        }
