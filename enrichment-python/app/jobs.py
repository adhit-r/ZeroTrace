"""
Background Job System using ARQ
Replaces basic Redis sorted sets with proper job queue
Uses Valkey (Redis-compatible) for job queue storage
"""

import asyncio
import logging
from typing import Dict, Any, Optional
from datetime import datetime
from arq import create_pool
from arq.connections import RedisSettings
from arq.worker import Worker
from .core.config import settings

logger = logging.getLogger(__name__)


class JobQueue:
    """
    Background job queue using ARQ (Async Redis Queue)
    Uses Valkey (Redis-compatible) for job queue storage
    Provides retries, dead-letter queues, and job scheduling
    """
    
    def __init__(self):
        self.redis_pool = None
        self.worker: Optional[Worker] = None
        
    async def initialize(self):
        """Initialize ARQ connection pool"""
        try:
            redis_url = settings.get_redis_dsn() if hasattr(settings, 'get_redis_dsn') else f"redis://{settings.redis_host}:{settings.redis_port}/{settings.redis_db}"
            
            # Parse Redis URL for ARQ
            redis_settings = RedisSettings.from_dsn(redis_url)
            
            self.redis_pool = await create_pool(redis_settings)
            logger.info("ARQ job queue initialized")
            
        except Exception as e:
            logger.error(f"Failed to initialize ARQ: {e}")
            raise
    
    async def enqueue_job(
        self,
        job_name: str,
        job_data: Dict[str, Any],
        job_id: Optional[str] = None,
        delay: Optional[float] = None,
        max_retries: int = 3
    ) -> str:
        """
        Enqueue a background job
        
        Args:
            job_name: Name of the job function
            job_data: Job payload data
            job_id: Optional job ID
            delay: Optional delay in seconds
            max_retries: Maximum retry attempts
        """
        if not self.redis_pool:
            await self.initialize()
        
        try:
            job = await self.redis_pool.enqueue_job(
                job_name,
                **job_data,
                _job_id=job_id,
                _defer_by=delay,
                _max_retries=max_retries
            )
            return job.job_id
        except Exception as e:
            logger.error(f"Failed to enqueue job {job_name}: {e}")
            raise
    
    async def get_job_status(self, job_id: str) -> Dict[str, Any]:
        """Get job status"""
        if not self.redis_pool:
            await self.initialize()
        
        try:
            job = await self.redis_pool.get_job_result(job_id)
            return {
                "job_id": job_id,
                "status": "completed" if job else "pending",
                "result": job if job else None
            }
        except Exception as e:
            logger.error(f"Failed to get job status {job_id}: {e}")
            return {"job_id": job_id, "status": "error", "error": str(e)}
    
    async def close(self):
        """Close ARQ connection pool"""
        if self.redis_pool:
            await self.redis_pool.close()


from .services.enrichment import enrichment_service

async def startup(ctx):
    """Initialize services on worker startup"""
    logger.info("ARQ Worker starting up...")
    await enrichment_service.initialize()

async def shutdown(ctx):
    """Cleanup on worker shutdown"""
    logger.info("ARQ Worker shutting down...")
    await enrichment_service.close()

async def enrichment_job(ctx, app_data: Dict[str, Any]) -> Dict[str, Any]:
    """
    Background job for enrichment processing
    
    Args:
        ctx: ARQ context
        app_data: Application data to enrich (expects 'packages' list or similar)
    """
    logger.info(f"Processing enrichment job for app: {app_data.get('app_name', 'unknown')}")
    
    try:
        # Extract software list from app_data
        # Supports multiple formats: 
        # 1. Direct list of packages
        # 2. 'packages' or 'dependencies' key
        software_list = []
        if isinstance(app_data, list):
            software_list = app_data
        elif "packages" in app_data:
            software_list = app_data["packages"]
        elif "dependencies" in app_data:
            software_list = app_data["dependencies"]
        else:
            # Fallback: treat app_data as single item if it has name/version
            if "name" in app_data:
                software_list = [app_data]
        
        if not software_list:
            logger.warning("No software items found in enrichment job payload")
            return {"status": "skipped", "reason": "no_packages"}

        # Perform enrichment
        results = await enrichment_service.enrich_software(software_list)
        
        # In a real app, we might save this to DB here if not handled by service
        # For now, we return detailed results which ARQ stores
        
        return {
            "status": "completed",
            "app_id": app_data.get("id"),
            "processed_count": len(results),
            "vulnerabilities_found": sum(len(item.get("vulnerabilities", [])) for item in results),
            "processed_at": datetime.utcnow().isoformat(),
            "results": results 
        }
    except Exception as e:
        logger.error(f"Enrichment job failed: {e}")
        return {
            "status": "failed",
            "error": str(e)
        }


async def batch_enrichment_job(ctx, batch_data: Dict[str, Any]) -> Dict[str, Any]:
    """
    Background job for batch enrichment processing
    """
    apps = batch_data.get("apps", [])
    logger.info(f"Processing batch enrichment job: {len(apps)} apps")
    
    processed_count = 0
    errors = 0
    
    for app in apps:
        try:
            # We can reusing the logic by calling enqueue or directly processing
            # Processing directly to save queue overhead for batch
            await enrichment_job(ctx, app)
            processed_count += 1
        except Exception as e:
            logger.error(f"Failed to process app in batch: {e}")
            errors += 1
    
    return {
        "status": "completed",
        "processed_count": processed_count,
        "error_count": errors,
        "processed_at": datetime.utcnow().isoformat()
    }


# ARQ worker functions mapping
class WorkerFunctions:
    """ARQ worker function registry"""
    
    functions = [
        enrichment_job,
        batch_enrichment_job,
    ]
    
    on_startup = startup
    on_shutdown = shutdown


# Global job queue instance
_job_queue: Optional[JobQueue] = None


async def get_job_queue() -> JobQueue:
    """Get or create global job queue instance"""
    global _job_queue
    
    if _job_queue is None:
        _job_queue = JobQueue()
        await _job_queue.initialize()
    
    return _job_queue

