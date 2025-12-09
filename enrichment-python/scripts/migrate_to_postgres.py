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

import ijson

import decimal

# ... (imports)

def decimal_default(obj):
    if isinstance(obj, decimal.Decimal):
        return float(obj)
    raise TypeError

async def migrate_data():
    """
    Migration Script (Streaming Version):
    1. Reads cve_data.json using ijson (streaming)
    2. Migrates raw CVEs to PostgreSQL (jsonb)
    3. Generates embeddings for CPEs and populates vector table
    """
    logger.info("Starting streaming migration to PostgreSQL...")
    
    # 1. Initialize Connections
    await db_manager.initialize()
    await semantic_matcher.initialize()
    
    # 2. Load JSON Data
    cve_file = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "cve_data.json")
    if not os.path.exists(cve_file):
        logger.error(f"File {cve_file} not found. Run update_cve_data.py first.")
        return

    logger.info(f"Streaming {cve_file}...")
    
    # 3. Process in Batches
    BATCH_SIZE = 100
    processed = 0
    
    # Buffers
    cve_records = []
    cpe_embedding_records = []
    
    try:
        with open(cve_file, 'rb') as f:
            # Assuming the file is a list of CVE objects
            parser = ijson.items(f, 'item')
            
            for cve_item in parser:
                # NVD API 2.0 format: {"cve": {...}, ...} or direct cve object
                cve_obj = cve_item.get('cve', cve_item) if isinstance(cve_item, dict) else cve_item
                cve_id = cve_obj.get('id') or cve_obj.get('CVE_data_meta', {}).get('ID')
                if not cve_id:
                    continue

                # Prepare CVE Record
                description = ""
                try:
                    # NVD API 2.0 format
                    desc_list = cve_obj.get('descriptions', [])
                    if desc_list and isinstance(desc_list, list):
                        for desc in desc_list:
                            if desc.get('lang') == 'en':
                                description = desc.get('value', '')
                                break
                        if not description and desc_list:
                            description = desc_list[0].get('value', '')
                    # Fallback to old format
                    elif cve_obj.get('description', {}).get('description_data'):
                        desc_list = cve_obj.get('description', {}).get('description_data', [])
                        if desc_list:
                            description = desc_list[0].get('value', '')
                except (KeyError, IndexError, AttributeError):
                    pass

                # Store the full CVE object as JSONB
                # FIX: Handle Decimal from ijson
                cve_records.append((cve_id, json.dumps(cve_obj, default=decimal_default), description))
                
                # Prepare CPE Embeddings
                cpes = extract_cpes(cve_obj)
                for cpe in cpes:
                    cpe_embedding_records.append({
                        "cve_id": cve_id,
                        "cpe": cpe
                    })

                # Flush Batch
                if len(cve_records) >= BATCH_SIZE:
                    await process_batch(cve_records, cpe_embedding_records)
                    processed += len(cve_records)
                    if processed % 1000 == 0:
                        logger.info(f"Progress: {processed} CVEs migrated")
                    cve_records = []
                    cpe_embedding_records = []

            # Process remaining
            if cve_records:
                await process_batch(cve_records, cpe_embedding_records)
                processed += len(cve_records)

    except Exception as e:
        logger.error(f"Streaming failed: {e}")
        raise

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
        # 1. Insert CVEs using psycopg2 for better JSONB handling
        import psycopg2
        import psycopg2.extras
        from app.core.config import settings
        
        # Get connection string
        conn_str = settings.get_database_dsn()
        conn = psycopg2.connect(conn_str)
        cur = conn.cursor()
        
        for cve_id, cve_json, description in cve_records:
            try:
                # Parse JSON to ensure it's valid
                import json as json_lib
                if isinstance(cve_json, str):
                    cve_data = json_lib.loads(cve_json)
                else:
                    cve_data = cve_json
                
                # Use psycopg2's Json adapter for JSONB
                cur.execute("""
                    INSERT INTO cves (id, data, description)
                    VALUES (%s, %s::jsonb, %s)
                    ON CONFLICT (id) DO UPDATE 
                    SET data = EXCLUDED.data, description = EXCLUDED.description, updated_at = NOW();
                """, (cve_id, psycopg2.extras.Json(cve_data), description or ''))
            except Exception as e:
                logger.warning(f"Failed to insert CVE {cve_id}: {e}")
                continue
        
        conn.commit()
        cur.close()
        conn.close()
        
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

