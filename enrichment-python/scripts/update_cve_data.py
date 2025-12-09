#!/usr/bin/env python3
"""
Automated CVE Database Update System
Fetches and processes CVE data from NVD API 2.0 with incremental updates
"""

import asyncio
import aiohttp
import json
import logging
import os
import time
from datetime import datetime, timedelta
from typing import List, Dict, Optional
from pathlib import Path
from collections import deque

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class CVEDatabaseUpdater:
    def __init__(self):
        self.nvd_api_key = os.getenv('NVD_API_KEY', '')
        self.base_url = 'https://services.nvd.nist.gov/rest/json/cves/2.0'
        # Rate limits: Without API key: 5 requests/30s, With API key: 50 requests/30s
        if self.nvd_api_key:
            self.max_concurrent = 10  # 10 concurrent requests with API key
            self.requests_per_window = 50
            self.window_seconds = 30
        else:
            self.max_concurrent = 2  # 2 concurrent requests without API key
            self.requests_per_window = 5
            self.window_seconds = 30
        
        # Rate limiting with sliding window
        self.request_times = deque()
        self.semaphore = asyncio.Semaphore(self.max_concurrent)
        
        self.cve_data_file = Path(__file__).parent.parent / 'cve_data.json'
        self.metadata_file = Path(__file__).parent.parent / 'cve_metadata.json'
        
    async def fetch_recent_cves(self, since_date: datetime) -> List[Dict]:
        """
        Fetch CVEs published since the given date using NVD API 2.0
        """
        cves = []
        start_index = 0
        results_per_page = 2000  # Maximum allowed by NVD API
        
        # Format date for NVD API 2.0 (ISO 8601 format with timezone)
        since_date_str = since_date.strftime('%Y-%m-%dT%H:%M:%S.000') + ' UTC'
        
        logger.info(f"Fetching CVEs published since {since_date_str}")
        
        async with aiohttp.ClientSession() as session:
            while True:
                params = {
                    'pubStartDate': since_date_str,
                    'resultsPerPage': results_per_page,
                    'startIndex': start_index
                }
                
                if self.nvd_api_key:
                    params['apiKey'] = self.nvd_api_key
                
                try:
                    # Rate limiting
                    await asyncio.sleep(self.rate_limit_delay)
                    
                    async with session.get(self.base_url, params=params) as response:
                        if response.status == 200:
                            data = await response.json()
                            
                            # Extract CVEs from response
                            cve_items = data.get('vulnerabilities', [])
                            if not cve_items:
                                logger.info("No more CVEs to fetch")
                                break
                            
                            # Process CVEs
                            for item in cve_items:
                                cve_data = item.get('cve', {})
                                if cve_data:
                                    cves.append(cve_data)
                            
                            logger.info(f"Fetched {len(cve_items)} CVEs (total: {len(cves)})")
                            
                            # Check if we have more pages
                            total_results = data.get('totalResults', 0)
                            if start_index + results_per_page >= total_results:
                                break
                            
                            start_index += results_per_page
                            
                        elif response.status == 403:
                            logger.error("NVD API rate limit exceeded. Consider using API key for higher limits.")
                            break
                        else:
                            logger.error(f"NVD API returned status {response.status}")
                            break
                
                except Exception as e:
                    logger.error(f"Error fetching CVEs: {e}")
                    break
        
        logger.info(f"Total CVEs fetched: {len(cves)}")
        return cves
    
    async def _wait_for_rate_limit(self):
        """Wait if we've hit the rate limit"""
        now = time.time()
        # Remove requests older than the window
        while self.request_times and self.request_times[0] < now - self.window_seconds:
            self.request_times.popleft()
        
        # If we're at the limit, wait
        if len(self.request_times) >= self.requests_per_window:
            sleep_time = self.window_seconds - (now - self.request_times[0]) + 0.1
            if sleep_time > 0:
                await asyncio.sleep(sleep_time)
                # Clean up again after waiting
                now = time.time()
                while self.request_times and self.request_times[0] < now - self.window_seconds:
                    self.request_times.popleft()
        
        # Record this request
        self.request_times.append(time.time())
    
    async def _fetch_page(self, session: aiohttp.ClientSession, start_index: int, results_per_page: int) -> Optional[Dict]:
        """Fetch a single page of CVEs with rate limiting"""
        async with self.semaphore:
            await self._wait_for_rate_limit()
            
            params = {
                'resultsPerPage': results_per_page,
                'startIndex': start_index
            }
            
            if self.nvd_api_key:
                params['apiKey'] = self.nvd_api_key
            
            try:
                async with session.get(self.base_url, params=params) as response:
                    if response.status == 200:
                        return await response.json()
                    elif response.status == 403 or response.status == 429:
                        # Rate limit exceeded - wait and return None to retry
                        wait_time = 120 if not self.nvd_api_key else 60
                        logger.warning(f"Rate limit hit at index {start_index}, waiting {wait_time}s...")
                        await asyncio.sleep(wait_time)
                        return None  # Signal to retry
                    elif response.status == 404:
                        logger.info(f"NVD API returned 404 at index {start_index}")
                        return {"vulnerabilities": [], "totalResults": 0}  # Signal end
                    else:
                        logger.error(f"NVD API returned status {response.status} at index {start_index}")
                        return None
            except Exception as e:
                logger.error(f"Error fetching page at index {start_index}: {e}")
                return None
    
    async def fetch_all_cves(self, max_cves: Optional[int] = None) -> List[Dict]:
        """
        Fetch all CVEs from NVD API 2.0 using concurrent requests
        
        Args:
            max_cves: Optional limit on number of CVEs to fetch (for testing)
        """
        cves = []
        results_per_page = 2000  # Maximum allowed by NVD API
        total_results = None
        
        logger.info(f"Fetching all CVEs from NVD API (using {self.max_concurrent} concurrent requests)")
        if max_cves:
            logger.info(f"Limited to {max_cves:,} CVEs")
        
        async with aiohttp.ClientSession() as session:
            # First, get total count
            first_page = await self._fetch_page(session, 0, 1)
            if first_page:
                total_results = first_page.get('totalResults', 0)
                logger.info(f"Total CVEs in NVD: {total_results:,}")
            
            if not total_results:
                logger.error("Could not determine total CVE count")
                return []
            
            # Calculate number of pages
            total_pages = (total_results + results_per_page - 1) // results_per_page
            if max_cves:
                total_pages = min(total_pages, (max_cves + results_per_page - 1) // results_per_page)
            
            logger.info(f"Fetching {total_pages} pages with {self.max_concurrent} concurrent requests")
            
            # Create tasks for all pages
            tasks = []
            for page in range(total_pages):
                start_index = page * results_per_page
                tasks.append(self._fetch_page(session, start_index, results_per_page))
            
            # Process results as they complete
            completed = 0
            for coro in asyncio.as_completed(tasks):
                try:
                    data = await coro
                    if data is None:
                        # Retry failed requests
                        continue
                    
                    cve_items = data.get('vulnerabilities', [])
                    if not cve_items:
                        continue
                    
                    # Process CVEs
                    for item in cve_items:
                        cve_data = item.get('cve', {})
                        if cve_data:
                            cves.append(cve_data)
                    
                    completed += 1
                    logger.info(f"Fetched page {completed}/{total_pages} - {len(cve_items)} CVEs (total: {len(cves):,}/{total_results:,})")
                    
                    # Check if we've reached max_cves limit
                    if max_cves and len(cves) >= max_cves:
                        logger.info(f"Reached max_cves limit of {max_cves:,}")
                        # Cancel remaining tasks
                        for task in tasks:
                            if not task.done():
                                task.cancel()
                        break
                    
                    # Progress update every 10,000 CVEs
                    if len(cves) % 10000 == 0:
                        percent = len(cves) * 100 // total_results
                        logger.info(f"Progress: {len(cves):,} CVEs downloaded ({percent}% complete)")
                
                except Exception as e:
                    logger.error(f"Error processing page: {e}")
                    continue
        
        logger.info(f"Total CVEs fetched: {len(cves):,}")
        return cves
                
    def load_existing_cves(self) -> List[Dict]:
        """Load existing CVE data from file"""
        if not self.cve_data_file.exists():
            logger.info("No existing CVE data file found")
            return []
        
        try:
            with open(self.cve_data_file, 'r') as f:
                data = json.load(f)
                if isinstance(data, list):
                    return data
                else:
                    logger.warning("CVE data file format unexpected")
                    return []
        except Exception as e:
            logger.error(f"Error loading existing CVE data: {e}")
            return []
    
    def load_metadata(self) -> Dict:
        """Load update metadata"""
        if not self.metadata_file.exists():
            return {
                'last_updated': None,
                'total_cves': 0,
                'update_duration': 0,
                'version': '1.0'
            }
        
        try:
            with open(self.metadata_file, 'r') as f:
                return json.load(f)
        except Exception as e:
            logger.error(f"Error loading metadata: {e}")
            return {
                'last_updated': None,
                'total_cves': 0,
                'update_duration': 0,
                'version': '1.0'
            }
    
    def save_metadata(self, metadata: Dict):
        """Save update metadata"""
        try:
            with open(self.metadata_file, 'w') as f:
                json.dump(metadata, f, indent=2)
        except Exception as e:
            logger.error(f"Error saving metadata: {e}")
    
    def merge_cves(self, existing_cves: List[Dict], new_cves: List[Dict]) -> List[Dict]:
        """Merge new CVEs with existing ones, removing duplicates"""
        # Create a map of existing CVEs by ID
        existing_map = {cve.get('id'): cve for cve in existing_cves}
        
        # Add new CVEs, updating existing ones
        for new_cve in new_cves:
            cve_id = new_cve.get('id')
            if cve_id:
                existing_map[cve_id] = new_cve
        
        # Convert back to list
        merged_cves = list(existing_map.values())
        
        logger.info(f"Merged CVEs: {len(existing_cves)} existing + {len(new_cves)} new = {len(merged_cves)} total")
        return merged_cves
    
    def save_cves(self, cves: List[Dict]):
        """Save CVE data to file"""
        try:
            # Create backup of existing file
            if self.cve_data_file.exists():
                backup_file = self.cve_data_file.with_suffix('.json.backup')
                self.cve_data_file.rename(backup_file)
                logger.info(f"Created backup: {backup_file}")
            
            # Save new data
            with open(self.cve_data_file, 'w') as f:
                json.dump(cves, f, indent=2)
            
            logger.info(f"Saved {len(cves)} CVEs to {self.cve_data_file}")
                
        except Exception as e:
            logger.error(f"Error saving CVE data: {e}")
            # Restore backup if save failed
            backup_file = self.cve_data_file.with_suffix('.json.backup')
            if backup_file.exists():
                backup_file.rename(self.cve_data_file)
                logger.info("Restored backup due to save error")
    
    async def update_database(self, force_full_update: bool = False, max_cves: Optional[int] = None):
        """Update the CVE database with latest data"""
        start_time = time.time()
        
        # Load existing data and metadata
        existing_cves = self.load_existing_cves()
        metadata = self.load_metadata()
        
        # Determine update strategy
        if force_full_update or not metadata.get('last_updated'):
            # Full update - fetch all CVEs (no date filter for initial population)
            logger.info("Performing full CVE database update (fetching all CVEs)")
            logger.info(f"Estimated total CVEs in NVD: ~320,000")
            if max_cves:
                logger.info(f"Limiting download to {max_cves:,} CVEs")
            new_cves = await self.fetch_all_cves(max_cves=max_cves)
        else:
            # Incremental update - fetch CVEs since last update
            last_updated = datetime.fromisoformat(metadata['last_updated'])
            since_date = last_updated
            logger.info(f"Performing incremental update since {last_updated}")
            new_cves = await self.fetch_recent_cves(since_date)
        
        if not new_cves:
            logger.info("No new CVEs found")
            return
        
        # Merge with existing data
        merged_cves = self.merge_cves(existing_cves, new_cves)
        
        # Save updated data
        self.save_cves(merged_cves)
        
        # Update metadata
        update_duration = time.time() - start_time
        new_metadata = {
            'last_updated': datetime.now().isoformat(),
            'total_cves': len(merged_cves),
            'update_duration': update_duration,
            'version': '1.0',
            'new_cves_count': len(new_cves),
            'api_key_used': bool(self.nvd_api_key)
        }
        
        self.save_metadata(new_metadata)
        
        logger.info(f"CVE database update completed in {update_duration:.2f}s")
        logger.info(f"Total CVEs: {len(merged_cves)}, New CVEs: {len(new_cves)}")

async def main():
    """Main entry point"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Update CVE database from NVD API')
    parser.add_argument('--force-full', action='store_true', 
                       help='Force full database update (all CVEs, ~390K)')
    parser.add_argument('--api-key', type=str, 
                       help='NVD API key for higher rate limits')
    parser.add_argument('--max-cves', type=int, default=None,
                       help='Limit number of CVEs to download (for testing)')
    
    args = parser.parse_args()
    
    # Set API key if provided
    if args.api_key:
        os.environ['NVD_API_KEY'] = args.api_key
    
    # Create updater and run
    updater = CVEDatabaseUpdater()
    await updater.update_database(force_full_update=args.force_full, max_cves=args.max_cves)

if __name__ == '__main__':
    asyncio.run(main())
