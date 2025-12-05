import contextlib
import uvicorn
from fastapi import FastAPI, HTTPException, Request, BackgroundTasks
from fastapi.responses import JSONResponse
from fastapi.middleware.cors import CORSMiddleware
from typing import List, Dict

# Unified Architecture Imports
from .core.config import settings
from .core.logging import configure_logging, get_logger
from .services.enrichment import enrichment_service

# Configure Logger (structlog)
configure_logging()
logger = get_logger(__name__)

@contextlib.asynccontextmanager
async def lifespan(app: FastAPI):
    """
    Unified Application Lifespan
    Initializes all subsystems: Cache, DB, AI, HTTP Client
    """
    logger.info("Application starting up...")
    
    # Initialize Unified Service (which in turn inits Cache, DB, AI)
    await enrichment_service.initialize()
    
    yield
    
    logger.info("Application shutting down...")
    await enrichment_service.close()

# FastAPI App Definition
app = FastAPI(
    title=settings.app_name,
    version="2.0.0", # Bump for Enterprise Release
    lifespan=lifespan,
    debug=settings.debug
)

# Middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.get_cors_origins_list(),
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# --- Routes ---

@app.get("/")
async def root():
    """Root endpoint with service information"""
    return {
        "name": settings.app_name,
        "version": "2.0.0",
        "status": "running",
        "description": "ZeroTrace Enrichment Service - CVE and CPE enrichment for software dependencies",
        "endpoints": {
            "health": "/health",
            "enrich": "/enrich/software",
            "docs": "/docs",
            "redoc": "/redoc"
        },
        "environment": settings.environment
    }

@app.get("/health")
async def health_check():
    """Enterprise Health Check"""
    # In a real enterprise app, we should check DB/Redis connectivity status here
    return {
        "status": "healthy",
        "service": settings.app_name,
        "version": "2.0.0",
        "environment": settings.environment
    }

@app.post("/enrich/software")
async def enrich_software(software_list: List[Dict]):
    """
    Main Enrichment Endpoint
    Accepts list of software items: [{"name": "nginx", "version": "1.18", "vendor": "nginx"}]
    Returns enriched data with CVEs and CPEs.
    """
    if not software_list:
        raise HTTPException(status_code=400, detail="Empty software list")

    if len(software_list) > settings.max_concurrent_requests:
        raise HTTPException(status_code=400, detail=f"Batch size limit exceeded ({settings.max_concurrent_requests})")

    try:
        results = await enrichment_service.enrich_software(software_list)
        return {
            "success": True,
            "count": len(results),
            "data": results
        }
    except Exception as e:
        logger.error("Enrichment API failed", error=str(e))
        raise HTTPException(status_code=500, detail="Internal Server Error")

# --- Training / Admin Routes (Optional, protected in prod) ---

@app.post("/admin/migrate")
async def trigger_migration(background_tasks: BackgroundTasks):
    """Trigger data migration from JSON to Postgres (Async)"""
    # Import script function dynamically to avoid load time
    from scripts.migrate_to_postgres import migrate_data
    background_tasks.add_task(migrate_data)
    return {"status": "Migration started in background"}

if __name__ == "__main__":
    # Use settings for host/port
    uvicorn.run(
        "app.main:app", 
        host=settings.enrichment_host, 
        port=settings.enrichment_port,
        reload=settings.debug
    )
