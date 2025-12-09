"""
Threat Intelligence Service
Integrates with MITRE ATT&CK, AlienVault OTX, and OpenCVE for enhanced vulnerability context
"""

import httpx
import json
import logging
from typing import Dict, List, Optional
from pathlib import Path

from ..core.config import settings
from ..core.logging import get_logger

logger = get_logger(__name__)


class ThreatIntelService:
    """Base class for threat intelligence services"""
    
    def __init__(self):
        self.enabled = False
        self.http_client: Optional[httpx.AsyncClient] = None
    
    async def initialize(self):
        """Initialize HTTP client"""
        if not self.http_client:
            self.http_client = httpx.AsyncClient(timeout=10.0)
    
    async def close(self):
        """Close HTTP client"""
        if self.http_client:
            await self.http_client.aclose()


class MITREATTACKService(ThreatIntelService):
    """
    MITRE ATT&CK Integration
    Provides attack pattern context for vulnerabilities
    Data Source: https://github.com/mitre/cti
    """
    
    def __init__(self):
        super().__init__()
        self.enabled = getattr(settings, 'mitre_attack_enabled', False)
        self.data_url = getattr(settings, 'mitre_attack_data_url', 
                                'https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json')
        self.attack_data: Optional[Dict] = None
    
    async def initialize(self):
        """Load MITRE ATT&CK data"""
        await super().initialize()
        if not self.enabled:
            return
        
        try:
            # Load MITRE ATT&CK data (can be cached locally)
            if self.http_client:
                response = await self.http_client.get(self.data_url)
                if response.status_code == 200:
                    self.attack_data = response.json()
                    logger.info("MITRE ATT&CK data loaded successfully")
        except Exception as e:
            logger.warning(f"Failed to load MITRE ATT&CK data: {e}")
            self.enabled = False
    
    async def get_attack_patterns(self, cve_id: str, cwe_id: Optional[str] = None) -> List[Dict]:
        """
        Get MITRE ATT&CK attack patterns related to a CVE or CWE
        
        Args:
            cve_id: CVE identifier
            cwe_id: Optional CWE identifier
            
        Returns:
            List of attack patterns with techniques and tactics
        """
        if not self.enabled or not self.attack_data:
            return []
        
        try:
            patterns = []
            # Search for attack patterns related to the CVE/CWE
            # This is a simplified implementation - full implementation would
            # map CVEs/CWEs to MITRE ATT&CK techniques based on vulnerability descriptions
            if self.attack_data and 'objects' in self.attack_data:
                for obj in self.attack_data['objects']:
                    if obj.get('type') == 'attack-pattern':
                        # Check if pattern relates to the CVE/CWE
                        # (Simplified - would need proper mapping logic)
                        external_refs = obj.get('external_references', [])
                        for ref in external_refs:
                            if ref.get('source_name') == 'cve' and ref.get('external_id') == cve_id:
                                patterns.append({
                                    'id': obj.get('id'),
                                    'name': obj.get('name'),
                                    'description': obj.get('description', ''),
                                    'kill_chain_phases': obj.get('kill_chain_phases', [])
                                })
            
            return patterns
        except Exception as e:
            logger.debug(f"Error fetching MITRE ATT&CK patterns: {e}")
            return []


class AlienVaultOTXService(ThreatIntelService):
    """
    AlienVault OTX Integration
    Provides IOCs (Indicators of Compromise) related to CVEs
    API: https://otx.alienvault.com/api/v1/
    """
    
    def __init__(self):
        super().__init__()
        self.enabled = getattr(settings, 'alienvault_otx_enabled', False)
        self.api_key = getattr(settings, 'alienvault_otx_api_key', None)
        self.base_url = 'https://otx.alienvault.com/api/v1'
    
    async def initialize(self):
        """Initialize OTX service"""
        await super().initialize()
        if not self.enabled or not self.api_key:
            self.enabled = False
            return
        
        # Test API key
        try:
            if self.http_client:
                headers = {'X-OTX-API-KEY': self.api_key}
                response = await self.http_client.get(f"{self.base_url}/user/me", headers=headers)
                if response.status_code == 200:
                    logger.info("AlienVault OTX API key validated")
                else:
                    logger.warning(f"AlienVault OTX API key validation failed: {response.status_code}")
                    self.enabled = False
        except Exception as e:
            logger.warning(f"Failed to validate AlienVault OTX API: {e}")
            self.enabled = False
    
    async def get_iocs(self, cve_id: str) -> List[Dict]:
        """
        Get IOCs (Indicators of Compromise) related to a CVE
        
        Args:
            cve_id: CVE identifier
            
        Returns:
            List of IOCs (IPs, domains, hashes, etc.)
        """
        if not self.enabled or not self.api_key:
            return []
        
        try:
            headers = {'X-OTX-API-KEY': self.api_key}
            # Search for pulses related to the CVE
            response = await self.http_client.get(
                f"{self.base_url}/search/pulses",
                headers=headers,
                params={'q': cve_id}
            )
            
            if response.status_code == 200:
                data = response.json()
                iocs = []
                for pulse in data.get('results', []):
                    # Extract IOCs from pulse
                    for indicator in pulse.get('indicators', []):
                        iocs.append({
                            'type': indicator.get('type'),
                            'indicator': indicator.get('indicator'),
                            'title': pulse.get('name'),
                            'description': pulse.get('description', '')
                        })
                return iocs
        except Exception as e:
            logger.debug(f"Error fetching OTX IOCs: {e}")
        
        return []


class CISAKEVService(ThreatIntelService):
    """
    CISA KEV (Known Exploited Vulnerabilities) Integration
    Provides mapping of CVEs to known exploits from CISA
    Data Source: https://www.cisa.gov/known-exploited-vulnerabilities-catalog
    """
    
    def __init__(self):
        super().__init__()
        self.enabled = getattr(settings, 'cisa_kev_enabled', True)  # Enabled by default
        self.kev_url = 'https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json'
        self.kev_data: Optional[Dict] = None
        self.kev_mapping: Dict[str, Dict] = {}  # CVE ID -> KEV data
    
    async def initialize(self):
        """Load CISA KEV catalog"""
        await super().initialize()
        if not self.enabled:
            return
        
        try:
            if self.http_client:
                response = await self.http_client.get(self.kev_url)
                if response.status_code == 200:
                    data = response.json()
                    self.kev_data = data
                    
                    # Build CVE -> KEV mapping
                    for vuln in data.get('vulnerabilities', []):
                        cve_id = vuln.get('cveID')
                        if cve_id:
                            self.kev_mapping[cve_id] = {
                                'cveID': cve_id,
                                'vendorProject': vuln.get('vendorProject'),
                                'product': vuln.get('product'),
                                'vulnerabilityName': vuln.get('vulnerabilityName'),
                                'dateAdded': vuln.get('dateAdded'),
                                'shortDescription': vuln.get('shortDescription'),
                                'requiredAction': vuln.get('requiredAction'),
                                'dueDate': vuln.get('dueDate'),
                                'knownRansomwareCampaignUse': vuln.get('knownRansomwareCampaignUse', 'Unknown')
                            }
                    
                    logger.info(f"CISA KEV catalog loaded: {len(self.kev_mapping)} CVEs with known exploits")
                else:
                    logger.warning(f"Failed to load CISA KEV catalog: {response.status_code}")
                    self.enabled = False
        except Exception as e:
            logger.warning(f"Failed to load CISA KEV catalog: {e}")
            self.enabled = False
    
    async def get_kev_info(self, cve_id: str) -> Optional[Dict]:
        """
        Get KEV information for a CVE
        
        Args:
            cve_id: CVE identifier
            
        Returns:
            KEV data if CVE is in catalog, None otherwise
        """
        if not self.enabled or not self.kev_mapping:
            return None
        
        return self.kev_mapping.get(cve_id)
    
    async def is_known_exploited(self, cve_id: str) -> bool:
        """Check if a CVE is in the KEV catalog"""
        return cve_id in self.kev_mapping if self.kev_mapping else False
    
    async def get_all_kev_cves(self) -> List[str]:
        """Get list of all CVE IDs in KEV catalog"""
        return list(self.kev_mapping.keys()) if self.kev_mapping else []


class OpenCVEService(ThreatIntelService):
    """
    OpenCVE Integration
    Provides enhanced CVE details from multiple sources (NVD, MITRE, RedHat, CISA)
    """
    
    def __init__(self):
        super().__init__()
        self.enabled = getattr(settings, 'opencve_enabled', False)
        self.api_url = getattr(settings, 'opencve_api_url', 'http://localhost:8000')
        self.api_key = getattr(settings, 'opencve_api_key', None)
    
    async def initialize(self):
        """Initialize OpenCVE service"""
        await super().initialize()
        if not self.enabled:
            return
        
        # Test connection
        try:
            if self.http_client:
                headers = {}
                if self.api_key:
                    headers['Authorization'] = f'Bearer {self.api_key}'
                
                response = await self.http_client.get(f"{self.api_url}/api/health", headers=headers)
                if response.status_code == 200:
                    logger.info("OpenCVE connection validated")
                else:
                    logger.warning(f"OpenCVE connection failed: {response.status_code}")
                    self.enabled = False
        except Exception as e:
            logger.warning(f"Failed to connect to OpenCVE: {e}")
            self.enabled = False
    
    async def get_cve_details(self, cve_id: str) -> Optional[Dict]:
        """
        Get enhanced CVE details from OpenCVE
        
        Args:
            cve_id: CVE identifier
            
        Returns:
            Enhanced CVE details with multi-source data
        """
        if not self.enabled:
            return None
        
        try:
            headers = {}
            if self.api_key:
                headers['Authorization'] = f'Bearer {self.api_key}'
            
            response = await self.http_client.get(
                f"{self.api_url}/api/cves/{cve_id}",
                headers=headers
            )
            
            if response.status_code == 200:
                return response.json()
        except Exception as e:
            logger.debug(f"Error fetching OpenCVE details: {e}")
        
        return None
    
    async def search_cves(self, query: str, filters: Optional[Dict] = None) -> List[Dict]:
        """
        Search CVEs in OpenCVE
        
        Args:
            query: Search query
            filters: Optional filters (cvss, vendor, product, etc.)
            
        Returns:
            List of matching CVEs
        """
        if not self.enabled:
            return []
        
        try:
            headers = {}
            if self.api_key:
                headers['Authorization'] = f'Bearer {self.api_key}'
            
            params = {'q': query}
            if filters:
                params.update(filters)
            
            response = await self.http_client.get(
                f"{self.api_url}/api/cves",
                headers=headers,
                params=params
            )
            
            if response.status_code == 200:
                data = response.json()
                return data.get('data', [])
        except Exception as e:
            logger.debug(f"Error searching OpenCVE: {e}")
        
        return []


# Global instances
mitre_attack_service = MITREATTACKService()
alienvault_otx_service = AlienVaultOTXService()
opencve_service = OpenCVEService()
cisa_kev_service = CISAKEVService()


async def initialize_threat_intel_services():
    """Initialize all threat intelligence services"""
    await mitre_attack_service.initialize()
    await alienvault_otx_service.initialize()
    await opencve_service.initialize()
    await cisa_kev_service.initialize()  # KEV is important, load it


async def close_threat_intel_services():
    """Close all threat intelligence services"""
    await mitre_attack_service.close()
    await alienvault_otx_service.close()
    await opencve_service.close()
    await cisa_kev_service.close()

