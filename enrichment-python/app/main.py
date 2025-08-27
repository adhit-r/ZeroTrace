from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import os
import json
from typing import List, Dict
from .cve_enrichment import cve_service

app = FastAPI(
    title="ZeroTrace Enrichment Service",
    description="Python service for vulnerability data enrichment",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def root():
    return {"message": "ZeroTrace Enrichment Service"}

@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "enrichment"}

@app.post("/enrich/software")
async def enrich_software(software_list: List[Dict]):
    """Enrich software list with CVE data"""
    try:
        enriched_data = await cve_service.enrich_software(software_list)
        return {
            "success": True,
            "data": enriched_data,
            "message": f"Enriched {len(enriched_data)} software items"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/enrich/batch")
async def enrich_batch(background_tasks: BackgroundTasks, software_list: List[Dict]):
    """Enrich software list in background"""
    try:
        # Start background enrichment
        background_tasks.add_task(cve_service.enrich_software, software_list)
        return {
            "success": True,
            "message": f"Started enrichment for {len(software_list)} software items"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.get("/enrich/status/{job_id}")
async def get_enrichment_status(job_id: str):
    """Get enrichment job status"""
    # This would track job status in a real implementation
    return {
        "job_id": job_id,
        "status": "completed",
        "message": "Enrichment completed"
    }

if __name__ == "__main__":
    port = int(os.getenv("ENRICHMENT_PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)
