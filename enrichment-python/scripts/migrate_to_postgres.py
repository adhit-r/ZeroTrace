import asyncio
import json
import os
import sys
from typing import List, Dict

# Add parent directory to path to allow importing app modules
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from app.core.config import settings
from app.core.database import db_manager
from app.core.logging import configure_logging, get_logger
from app.services.cpe_matcher import semantic_matcher

# Configure logging
configure_logging()
logger = get_logger(__name__)

async def migrate_data():
    """
    Migration Script:
    1. Reads cve_data.json
    2. Migrates raw CVEs to PostgreSQL (jsonb)
    3. Generates embeddings for CPEs and populates vector table
    """
    logger.info("Starting migration to PostgreSQL...")
    
    # 1. Initialize Connections
    await db_manager.initialize()
    await semantic_matcher.initialize()
    
    # 2. Load JSON Data
    cve_file = "cve_data.json"
    if not os.path.exists(cve_file):
        logger.error(f"File {cve_file} not found. Run update_cve_data.py first.")
        return

    logger.info(f"Reading {cve_file}...")
    with open(cve_file, 'r') as f:
        cve_data = json.load(f)
    
    logger.info(f"Loaded {len(cve_data)} CVEs from file")

    # 3. Process in Batches
    BATCH_SIZE = 100
    total_cves = len(cve_data)
    processed = 0
    
    # Buffers
    cve_records = []
    cpe_embedding_records = []
    
    for cve in cve_data:
        cve_id = cve.get('id') or cve.get('cve', {}).get('CVE_data_meta', {}).get('ID')
        if not cve_id:
            continue

        # Prepare CVE Record
        # Extract description for potential future search
        description = ""
        try:
            desc_list = cve.get('cve', {}).get('description', {}).get('description_data', [])
            if desc_list:
                description = desc_list[0].get('value', '')
        except (KeyError, IndexError):
            pass

        cve_records.append((cve_id, json.dumps(cve), description))
        
        # Prepare CPE Embeddings
        # Extract CPEs from configurations
        cpes = extract_cpes(cve)
        for cpe in cpes:
            # We will generate embedding later in batch
            cpe_embedding_records.append({
                "cve_id": cve_id,
                "cpe": cpe
            })

        # Flush Batch
        if len(cve_records) >= BATCH_SIZE:
            await process_batch(cve_records, cpe_embedding_records)
            processed += len(cve_records)
            logger.info(f"Progress: {processed}/{total_cves} CVEs migrated")
            cve_records = []
            cpe_embedding_records = []

    # Process remaining
    if cve_records:
        await process_batch(cve_records, cpe_embedding_records)
        processed += len(cve_records)

    logger.info(f"Migration completed. {processed} CVEs processed.")
    await db_manager.close()

def extract_cpes(cve_data: Dict) -> List[str]:
    """Extract CPE strings from CVE JSON"""
    cpes = set()
    try:
        configurations = cve_data.get('configurations', {})
        nodes = configurations.get('nodes', [])
        for node in nodes:
            cpe_match = node.get('cpe_match', [])
            for match in cpe_match:
                if match.get('vulnerable'):
                    cpe_str = match.get('cpe23Uri')
                    if cpe_str:
                        cpes.add(cpe_str)
    except Exception:
        pass
    return list(cpes)

async def process_batch(cve_records: List, cpe_records: List[Dict]):
    """Insert CVEs and CPE Embeddings"""
    try:
        # 1. Insert CVEs (Upsert)
        await db_manager.executemany("""
            INSERT INTO cves (id, data, description)
            VALUES ($1, $2::jsonb, $3)
            ON CONFLICT (id) DO UPDATE 
            SET data = EXCLUDED.data, description = EXCLUDED.description, updated_at = NOW();
        """, cve_records)
        
        # 2. Generate Embeddings & Insert
        if settings.pgvector_enabled and semantic_matcher.model:
            embedding_batch = []
            
            # Deduplicate CPEs to avoid duplicate embedding work
            unique_cpes = {rec['cpe'] for rec in cpe_records}
            
            # This part can be slow on CPU - consider smaller sub-batches
            for cpe_str in unique_cpes:
                # "Productize" the CPE for better embedding
                # cpe:2.3:a:apache:http_server:2.4.50 -> "apache http server 2.4.50"
                text_to_encode = parse_cpe_to_text(cpe_str)
                vector = semantic_matcher.encode(text_to_encode)
                
                # Find all CVEs linked to this CPE in current batch
                linked_cves = [r['cve_id'] for r in cpe_records if r['cpe'] == cpe_str]
                
                for linked_cve in linked_cves:
                    embedding_batch.append((linked_cve, cpe_str, str(vector)))
            
            if embedding_batch:
                await db_manager.executemany("""
                    INSERT INTO cpe_embeddings (cve_id, cpe_string, embedding)
                    VALUES ($1, $2, $3::vector)
                    ON CONFLICT (cve_id, cpe_string) DO NOTHING;
                """, embedding_batch)

    except Exception as e:
        logger.error(f"Batch insert failed: {e}")

def parse_cpe_to_text(cpe: str) -> str:
    """Convert technical CPE string to natural language for embedding"""
    # cpe:2.3:a:vendor:product:version:...
    parts = cpe.split(':')
    if len(parts) >= 6:
        vendor = parts[3].replace('_', ' ')
        product = parts[4].replace('_', ' ')
        version = parts[5]
        if version == '*': version = ''
        return f"{vendor} {product} {version}".strip()
    return cpe

if __name__ == "__main__":
    asyncio.run(migrate_data())

