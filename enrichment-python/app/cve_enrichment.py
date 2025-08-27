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
import os

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
        
    async def enrich_software(self, software_list: List[Dict]) -> List[Dict]:
        """
        Enrich software list with CVE data
        """
        enriched_results = []
        
        for software in software_list:
            try:
                # Extract software info
                name = software.get('name', '')
                version = software.get('version', '')
                
                # Search for CVEs
                cves = await self.search_cves(name, version)
                
                # Enrich software data
                enriched_software = {
                    **software,
                    'cves': cves,
                    'vulnerability_count': len(cves),
                    'enriched_at': datetime.utcnow().isoformat()
                }
                
                enriched_results.append(enriched_software)
                logger.info(f"Enriched {name} {version} with {len(cves)} CVEs")
                
            except Exception as e:
                logger.error(f"Error enriching {software.get('name', 'unknown')}: {e}")
                enriched_results.append({
                    **software,
                    'cves': [],
                    'vulnerability_count': 0,
                    'enriched_at': datetime.utcnow().isoformat(),
                    'error': str(e)
                })
        
        return enriched_results
    
    async def search_cves(self, software_name: str, version: str = None) -> List[Dict]:
        """
        Search for CVEs related to software
        """
        cves = []
        
        # Search in multiple sources
        tasks = [
            self.search_nvd(software_name, version),
            self.search_cve_search(software_name, version)
        ]
        
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        for result in results:
            if isinstance(result, list):
                cves.extend(result)
        
        # Remove duplicates based on CVE ID
        unique_cves = {}
        for cve in cves:
            cve_id = cve.get('id') or cve.get('cve_id')
            if cve_id and cve_id not in unique_cves:
                unique_cves[cve_id] = cve
        
        return list(unique_cves.values())
    
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

# Global instance
cve_service = CVEEnrichmentService()

