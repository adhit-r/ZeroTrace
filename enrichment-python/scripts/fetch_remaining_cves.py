#!/usr/bin/env python3
"""
Fetch remaining CVEs that were missed during parallel download
"""

import asyncio
import aiohttp
import json
import logging
import os
from pathlib import Path
from typing import List, Dict, Set

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

async def fetch_page(session: aiohttp.ClientSession, start_index: int, results_per_page: int = 2000, api_key: str = None) -> Dict:
    """Fetch a single page of CVEs"""
    base_url = 'https://services.nvd.nist.gov/rest/json/cves/2.0'
    params = {
        'resultsPerPage': results_per_page,
        'startIndex': start_index
    }
    
    if api_key:
        params['apiKey'] = api_key
    
    try:
        async with session.get(base_url, params=params) as response:
            if response.status == 200:
                return await response.json()
            elif response.status == 429 or response.status == 403:
                logger.warning(f"Rate limit at index {start_index}, waiting 60s...")
                await asyncio.sleep(60)
                return await fetch_page(session, start_index, results_per_page, api_key)
            else:
                logger.error(f"Error {response.status} at index {start_index}")
                return None
    except Exception as e:
        logger.error(f"Exception at index {start_index}: {e}")
        return None

async def fetch_remaining_cves():
    """Fetch CVEs that are missing"""
    api_key = os.getenv('NVD_API_KEY', '')
    cve_data_file = Path(__file__).parent.parent / 'cve_data.json'
    
    # Load existing CVEs
    existing_cve_ids = set()
    if cve_data_file.exists():
        with open(cve_data_file, 'r') as f:
            existing_data = json.load(f)
            for cve in existing_data:
                cve_id = cve.get('id')
                if cve_id:
                    existing_cve_ids.add(cve_id)
        logger.info(f"Loaded {len(existing_cve_ids):,} existing CVE IDs")
    
    # Get total count from NVD
    async with aiohttp.ClientSession() as session:
        first_page = await fetch_page(session, 0, 1, api_key)
        if not first_page:
            logger.error("Could not fetch total count")
            return
        
        total_results = first_page.get('totalResults', 0)
        logger.info(f"Total CVEs in NVD: {total_results:,}")
        logger.info(f"Existing CVEs: {len(existing_cve_ids):,}")
        logger.info(f"Missing: {total_results - len(existing_cve_ids):,}")
        
        if len(existing_cve_ids) >= total_results:
            logger.info(" All CVEs already downloaded!")
            return
        
        # Fetch all pages and find missing CVEs
        results_per_page = 2000
        total_pages = (total_results + results_per_page - 1) // results_per_page
        logger.info(f"Fetching {total_pages} pages to find missing CVEs...")
        
        new_cves = []
        rate_limit_delay = 0.6 if api_key else 6.0
        
        for page in range(total_pages):
            start_index = page * results_per_page
            
            # Rate limiting
            await asyncio.sleep(rate_limit_delay)
            
            data = await fetch_page(session, start_index, results_per_page, api_key)
            if not data:
                logger.warning(f"Failed to fetch page {page + 1}/{total_pages}")
                continue
            
            cve_items = data.get('vulnerabilities', [])
            for item in cve_items:
                cve_data = item.get('cve', {})
                if cve_data:
                    cve_id = cve_data.get('id')
                    if cve_id and cve_id not in existing_cve_ids:
                        new_cves.append(cve_data)
                        existing_cve_ids.add(cve_id)
            
            if (page + 1) % 10 == 0:
                logger.info(f"Processed {page + 1}/{total_pages} pages, found {len(new_cves):,} new CVEs")
        
        logger.info(f"Found {len(new_cves):,} missing CVEs")
        
        if new_cves:
            # Merge with existing
            if cve_data_file.exists():
                with open(cve_data_file, 'r') as f:
                    existing_data = json.load(f)
                all_cves = existing_data + new_cves
            else:
                all_cves = new_cves
            
            # Save
            backup_file = cve_data_file.with_suffix('.json.backup2')
            if cve_data_file.exists():
                import shutil
                shutil.copy(cve_data_file, backup_file)
                logger.info(f"Created backup: {backup_file}")
            
            with open(cve_data_file, 'w') as f:
                json.dump(all_cves, f, indent=2)
            
            logger.info(f" Saved {len(all_cves):,} total CVEs to {cve_data_file}")
        else:
            logger.info(" No missing CVEs found - all up to date!")

if __name__ == '__main__':
    asyncio.run(fetch_remaining_cves())

