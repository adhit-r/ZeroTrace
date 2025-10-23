"""
CVE Enrichment Service
Fetches and processes CVE data from multiple sources
"""

import asyncio
import aiohttp
import json
import logging
from typing import List, Dict, Optional
from datetime import datetime
import time
import os
from .ai_matching.cpe_matcher import cpe_matcher

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class CVEEnrichmentService:
    def __init__(self):
        self.nvd_api_key = os.getenv('NVD_API_KEY', '')
        self.base_urls = {
            'nvd': 'https://services.nvd.nist.gov/rest/json/cves/2.0',
            'vulndb': 'https://vuldb.com/api/v1/search',
            'cve_search': 'https://cve.circl.lu/api/search'
        }
        self.local_cve_data = self.load_local_cve_data()
        self.cache = {}  # In-memory cache for recent searches
        self.cache_ttl = 3600  # 1 hour cache TTL
        
    def load_local_cve_data(self) -> List[Dict]:
        """Load CVE data from local file"""
        try:
            cve_file = os.path.join(os.path.dirname(__file__), '..', 'cve_data.json')
            if os.path.exists(cve_file):
                with open(cve_file, 'r') as f:
                    data = json.load(f)
                    logger.info(f"Loaded {len(data)} CVEs from local database")
                    return data
            else:
                logger.warning("No local CVE data file found")
                return []
        except Exception as e:
            logger.error(f"Error loading local CVE data: {e}")
            return []
        
    async def enrich_software(self, software_list: List[Dict]) -> List[Dict]:
        """
        Enrich software list with CVE data
        """
        import time
        start_time = time.time()
        enriched_results = []
        
        for software in software_list:
            software_start = time.time()
            try:
                # Extract software info
                name = software.get('name', '')
                version = software.get('version', '')
                
                # Get CPE identifier for better matching with enhanced version support
                cpe_start = time.time()
                cpe_identifier = cpe_matcher.get_cpe_for_software(name, version, software.get('vendor'))
                cpe_duration = time.time() - cpe_start
                
                # Get detailed CPE matches for confidence scoring
                cpe_matches = cpe_matcher.match_software_to_cpe(name, version, software.get('vendor'))
                best_cpe_match = cpe_matches[0] if cpe_matches else None
                
                # Search for CVEs
                cve_search_start = time.time()
                cves = await self.search_cves(name, version, cpe_identifier)
                cve_search_duration = time.time() - cve_search_start
                
                # Enrich software data with enhanced CPE information
                enriched_software = {
                    **software,
                    'cves': cves,
                    'vulnerability_count': len(cves),
                    'cpe_identifier': cpe_identifier,
                    'cpe_confidence': best_cpe_match.get('confidence', 'UNKNOWN') if best_cpe_match else 'UNKNOWN',
                    'cpe_version_match': best_cpe_match.get('version_match', False) if best_cpe_match else False,
                    'cpe_similarity_score': best_cpe_match.get('similarity_score', 0.0) if best_cpe_match else 0.0,
                    'enriched_at': datetime.utcnow().isoformat(),
                    'performance_metrics': {
                        'cpe_matching_duration_ms': cpe_duration * 1000,
                        'cve_search_duration_ms': cve_search_duration * 1000,
                        'software_processing_duration_ms': (time.time() - software_start) * 1000
                    }
                }
                
                enriched_results.append(enriched_software)
                logger.info(f"Enriched {name} {version} with {len(cves)} CVEs in {cve_search_duration:.3f}s")
                
            except Exception as e:
                logger.error(f"Error enriching {software.get('name', 'unknown')}: {e}")
                enriched_results.append({
                    **software,
                    'cves': [],
                    'vulnerability_count': 0,
                    'enriched_at': datetime.utcnow().isoformat(),
                    'error': str(e),
                    'performance_metrics': {
                        'software_processing_duration_ms': (time.time() - software_start) * 1000
                    }
                })
        
        total_duration = time.time() - start_time
        logger.info(f"Enriched {len(software_list)} software items in {total_duration:.3f}s")
        
        return enriched_results
    
    async def search_cves(self, software_name: str, version: str = None, cpe_identifier: str = None) -> List[Dict]:
        """
        Search for CVEs related to software using hybrid approach with caching
        """
        # Check cache first
        cache_key = f"{software_name}_{version or 'any'}_{cpe_identifier or 'none'}"
        if cache_key in self.cache:
            cached_data = self.cache[cache_key]
            if time.time() - cached_data['timestamp'] < self.cache_ttl:
                logger.info(f"Found {len(cached_data['cves'])} CVEs in cache for {software_name}")
                return cached_data['cves']
        
        cves = []
        search_method = "local"
        
        # First, search local CVE database
        local_cves = self.search_local_cve_data(software_name, version)
        if local_cves:
            logger.info(f"Found {len(local_cves)} CVEs in local database for {software_name}")
            cves.extend(local_cves)
        else:
            # Fallback to online sources if no local data
            logger.info(f"No local CVEs found for {software_name}, searching online sources...")
            search_method = "online"
            
            tasks = [
                self.search_nvd(software_name, version),
                self.search_cve_search(software_name, version)
            ]
            
            results = await asyncio.gather(*tasks, return_exceptions=True)
            
            for result in results:
                if isinstance(result, list):
                    cves.extend(result)
            
            # Cache online results locally for future use
            if cves:
                await self._cache_online_results(software_name, version, cves)
        
        # Remove duplicates based on CVE ID
        unique_cves = {}
        for cve in cves:
            cve_id = cve.get('id') or cve.get('cve_id')
            if cve_id and cve_id not in unique_cves:
                unique_cves[cve_id] = cve
        
        # Add search metadata
        final_cves = list(unique_cves.values())
        for cve in final_cves:
            cve['search_method'] = search_method
            cve['searched_at'] = datetime.utcnow().isoformat()
        
        # Cache results
        self.cache[cache_key] = {
            'cves': final_cves,
            'timestamp': time.time()
        }
        
        return final_cves
    
    def search_local_cve_data(self, software_name: str, version: str = None) -> List[Dict]:
        """
        Search local CVE database for software vulnerabilities
        """
        matching_cves = []
        software_name_lower = software_name.lower()
        
        for cve_entry in self.local_cve_data:
            try:
                cve_data = cve_entry.get('cve', {})
                descriptions = cve_data.get('descriptions', [])
                
                # Check if software name appears in any description
                for desc in descriptions:
                    desc_text = desc.get('value', '').lower()
                    if software_name_lower in desc_text:
                        # Extract CVE information
                        cve_id = cve_data.get('id', '')
                        if cve_id:
                            # Get CVSS score
                            cvss_score = 0.0
                            severity = 'UNKNOWN'
                            
                            metrics = cve_data.get('metrics', {})
                            cvss_v31 = metrics.get('cvssMetricV31', [])
                            if cvss_v31:
                                cvss_data = cvss_v31[0].get('cvssData', {})
                                cvss_score = cvss_data.get('baseScore', 0.0)
                                severity = cvss_data.get('baseSeverity', 'UNKNOWN')
                            
                            # Create CVE entry
                            cve_entry = {
                                'id': cve_id,
                                'description': desc.get('value', ''),
                                'severity': severity,
                                'cvss_score': cvss_score,
                                'published_date': cve_data.get('published', ''),
                                'last_modified': cve_data.get('lastModified', ''),
                                'source': 'local_database'
                            }
                            
                            matching_cves.append(cve_entry)
                            break  # Found a match, move to next CVE
                            
            except Exception as e:
                logger.error(f"Error processing CVE entry: {e}")
                continue
        
        return matching_cves
    
    async def search_nvd(self, software_name: str, version: str = None) -> List[Dict]:
        """
        Search NVD database for CVEs
        """
        try:
            async with aiohttp.ClientSession() as session:
                # Build query
                query = f'"{software_name}"'
                if version:
                    query += f' "{version}"'
                
                params = {
                    'keywordSearch': query,
                    'resultsPerPage': 20
                }
                
                if self.nvd_api_key:
                    params['apiKey'] = self.nvd_api_key
                
                async with session.get(self.base_urls['nvd'], params=params) as response:
                    if response.status == 200:
                        data = await response.json()
                        return self.parse_nvd_response(data)
                    else:
                        logger.warning(f"NVD API returned {response.status}")
                        return []
                        
        except Exception as e:
            logger.error(f"Error searching NVD: {e}")
            return []
    
    async def search_cve_search(self, software_name: str, version: str = None) -> List[Dict]:
        """
        Search CVE Search database
        """
        try:
            async with aiohttp.ClientSession() as session:
                query = software_name
                if version:
                    query += f' {version}'
                
                params = {'q': query}
                
                async with session.get(self.base_urls['cve_search'], params=params) as response:
                    if response.status == 200:
                        data = await response.json()
                        return self.parse_cve_search_response(data)
                    else:
                        logger.warning(f"CVE Search API returned {response.status}")
                        return []
                        
        except Exception as e:
            logger.error(f"Error searching CVE Search: {e}")
            return []
    
    def parse_nvd_response(self, data: Dict) -> List[Dict]:
        """
        Parse NVD API response
        """
        cves = []
        
        try:
            vulnerabilities = data.get('vulnerabilities', [])
            
            for vuln in vulnerabilities:
                cve_data = vuln.get('cve', {})
                
                cve = {
                    'id': cve_data.get('id'),
                    'description': self.get_description(cve_data),
                    'severity': self.get_severity(cve_data),
                    'cvss_score': self.get_cvss_score(cve_data),
                    'published_date': cve_data.get('published'),
                    'last_modified': cve_data.get('lastModified'),
                    'source': 'NVD'
                }
                
                if cve['id']:
                    cves.append(cve)
                    
        except Exception as e:
            logger.error(f"Error parsing NVD response: {e}")
        
        return cves
    
    def parse_cve_search_response(self, data: Dict) -> List[Dict]:
        """
        Parse CVE Search API response
        """
        cves = []
        
        try:
            for cve_data in data:
                cve = {
                    'id': cve_data.get('id'),
                    'description': cve_data.get('summary', ''),
                    'severity': self.map_severity(cve_data.get('cvss', 0)),
                    'cvss_score': cve_data.get('cvss', 0),
                    'published_date': cve_data.get('Published'),
                    'last_modified': cve_data.get('Modified'),
                    'source': 'CVE Search'
                }
                
                if cve['id']:
                    cves.append(cve)
                    
        except Exception as e:
            logger.error(f"Error parsing CVE Search response: {e}")
        
        return cves
    
    def get_description(self, cve_data: Dict) -> str:
        """
        Extract description from CVE data
        """
        descriptions = cve_data.get('descriptions', [])
        for desc in descriptions:
            if desc.get('lang') == 'en':
                return desc.get('value', '')
        return ''
    
    def get_severity(self, cve_data: Dict) -> str:
        """
        Extract severity from CVE data
        """
        metrics = cve_data.get('metrics', {})
        
        # Check CVSS v3
        cvss_v3 = metrics.get('cvssMetricV31', []) or metrics.get('cvssMetricV30', [])
        if cvss_v3:
            return cvss_v3[0].get('cvssData', {}).get('baseSeverity', 'UNKNOWN')
        
        # Check CVSS v2
        cvss_v2 = metrics.get('cvssMetricV2', [])
        if cvss_v2:
            return self.map_cvss_v2_severity(cvss_v2[0].get('cvssData', {}).get('baseScore', 0))
        
        return 'UNKNOWN'
    
    def get_cvss_score(self, cve_data: Dict) -> float:
        """
        Extract CVSS score from CVE data
        """
        metrics = cve_data.get('metrics', {})
        
        # Check CVSS v3
        cvss_v3 = metrics.get('cvssMetricV31', []) or metrics.get('cvssMetricV30', [])
        if cvss_v3:
            return cvss_v3[0].get('cvssData', {}).get('baseScore', 0)
        
        # Check CVSS v2
        cvss_v2 = metrics.get('cvssMetricV2', [])
        if cvss_v2:
            return cvss_v2[0].get('cvssData', {}).get('baseScore', 0)
        
        return 0.0
    
    def map_cvss_v2_severity(self, score: float) -> str:
        """
        Map CVSS v2 score to severity
        """
        if score >= 7.0:
            return 'HIGH'
        elif score >= 4.0:
            return 'MEDIUM'
        else:
            return 'LOW'
    
    def map_severity(self, score: float) -> str:
        """
        Map CVSS score to severity
        """
        if score >= 7.0:
            return 'HIGH'
        elif score >= 4.0:
            return 'MEDIUM'
        else:
            return 'LOW'
    
    async def _cache_online_results(self, software_name: str, version: str, cves: List[Dict]):
        """
        Cache online search results locally for future use
        """
        try:
            cache_file = self.data_dir / 'cve_cache.json'
            
            # Load existing cache
            cache_data = {}
            if cache_file.exists():
                with open(cache_file, 'r') as f:
                    cache_data = json.load(f)
            
            # Add new results to cache
            cache_key = f"{software_name}_{version or 'any'}"
            cache_data[cache_key] = {
                'cves': cves,
                'cached_at': datetime.utcnow().isoformat(),
                'software_name': software_name,
                'version': version
            }
            
            # Save updated cache
            with open(cache_file, 'w') as f:
                json.dump(cache_data, f, indent=2)
            
            logger.info(f"Cached {len(cves)} CVEs for {software_name} {version}")
            
        except Exception as e:
            logger.error(f"Error caching online results: {e}")

# Global instance
cve_service = CVEEnrichmentService()

