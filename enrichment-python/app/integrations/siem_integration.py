"""
SIEM/SOAR Integration Hub - Comprehensive security tool integration
"""

import asyncio
import json
import logging
import requests
import time
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass
import pandas as pd
import numpy as np
from abc import ABC, abstractmethod

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class SIEMEvent:
    """Represents a SIEM event"""
    event_id: str
    timestamp: str
    source: str
    event_type: str
    severity: str
    description: str
    raw_data: Dict[str, Any]
    normalized_data: Dict[str, Any]
    correlation_id: Optional[str] = None
    tags: List[str] = None

@dataclass
class SOARPlaybook:
    """Represents a SOAR playbook"""
    playbook_id: str
    name: str
    description: str
    trigger_conditions: List[str]
    actions: List[SOARAction]
    status: str  # active, inactive, testing
    last_executed: Optional[str] = None
    success_rate: float = 0.0
    execution_count: int = 0

@dataclass
class SOARAction:
    """Represents a SOAR action"""
    action_id: str
    name: str
    type: str  # api_call, script, notification, enrichment
    parameters: Dict[str, Any]
    timeout: int = 300
    retry_count: int = 3
    success_condition: str = "status_code == 200"

@dataclass
class IntegrationConfig:
    """Represents integration configuration"""
    integration_id: str
    name: str
    type: str  # siem, soar, ticketing, chatops
    endpoint: str
    credentials: Dict[str, str]
    settings: Dict[str, Any]
    status: str  # active, inactive, error
    last_sync: Optional[str] = None

class SIEMConnector(ABC):
    """Abstract base class for SIEM connectors"""
    
    @abstractmethod
    async def connect(self, config: IntegrationConfig) -> bool:
        """Connect to SIEM system"""
        pass
    
    @abstractmethod
    async def send_event(self, event: SIEMEvent) -> bool:
        """Send event to SIEM"""
        pass
    
    @abstractmethod
    async def query_events(self, query: str, time_range: Tuple[str, str]) -> List[SIEMEvent]:
        """Query events from SIEM"""
        pass
    
    @abstractmethod
    async def get_alerts(self, filters: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Get alerts from SIEM"""
        pass

class SplunkConnector(SIEMConnector):
    """Splunk SIEM connector"""
    
    def __init__(self):
        self.session = None
        self.base_url = None
        self.auth_token = None
    
    async def connect(self, config: IntegrationConfig) -> bool:
        """Connect to Splunk"""
        try:
            self.base_url = config.endpoint
            self.session = requests.Session()
            
            # Authenticate with Splunk
            auth_url = f"{self.base_url}/services/auth/login"
            auth_data = {
                'username': config.credentials['username'],
                'password': config.credentials['password']
            }
            
            response = self.session.post(auth_url, data=auth_data)
            if response.status_code == 200:
                self.auth_token = response.text.strip()
                logger.info("Successfully connected to Splunk")
                return True
            else:
                logger.error(f"Failed to connect to Splunk: {response.status_code}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to Splunk: {e}")
            return False
    
    async def send_event(self, event: SIEMEvent) -> bool:
        """Send event to Splunk"""
        try:
            url = f"{self.base_url}/services/receivers/simple"
            headers = {
                'Authorization': f'Splunk {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            # Format event for Splunk
            splunk_event = {
                'time': int(time.time()),
                'host': event.source,
                'sourcetype': 'zerotrace',
                'event': event.normalized_data
            }
            
            response = self.session.post(url, headers=headers, json=splunk_event)
            return response.status_code == 200
            
        except Exception as e:
            logger.error(f"Error sending event to Splunk: {e}")
            return False
    
    async def query_events(self, query: str, time_range: Tuple[str, str]) -> List[SIEMEvent]:
        """Query events from Splunk"""
        try:
            url = f"{self.base_url}/services/search/jobs/export"
            headers = {
                'Authorization': f'Splunk {self.auth_token}',
                'Content-Type': 'application/x-www-form-urlencoded'
            }
            
            data = {
                'search': query,
                'earliest_time': time_range[0],
                'latest_time': time_range[1],
                'output_mode': 'json'
            }
            
            response = self.session.post(url, headers=headers, data=data)
            if response.status_code == 200:
                events = []
                for line in response.text.strip().split('\n'):
                    if line:
                        event_data = json.loads(line)
                        events.append(SIEMEvent(
                            event_id=event_data.get('_cd', ''),
                            timestamp=event_data.get('_time', ''),
                            source=event_data.get('host', ''),
                            event_type=event_data.get('sourcetype', ''),
                            severity=event_data.get('severity', 'medium'),
                            description=event_data.get('_raw', ''),
                            raw_data=event_data,
                            normalized_data=event_data
                        ))
                return events
            else:
                logger.error(f"Failed to query Splunk: {response.status_code}")
                return []
                
        except Exception as e:
            logger.error(f"Error querying Splunk: {e}")
            return []
    
    async def get_alerts(self, filters: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Get alerts from Splunk"""
        try:
            # Mock alert retrieval
            return [
                {
                    'alert_id': 'alert_001',
                    'title': 'High severity vulnerability detected',
                    'severity': 'high',
                    'timestamp': datetime.utcnow().isoformat(),
                    'source': 'zerotrace_scanner'
                }
            ]
        except Exception as e:
            logger.error(f"Error getting alerts from Splunk: {e}")
            return []

class QRadarConnector(SIEMConnector):
    """IBM QRadar SIEM connector"""
    
    def __init__(self):
        self.session = None
        self.base_url = None
        self.auth_token = None
    
    async def connect(self, config: IntegrationConfig) -> bool:
        """Connect to QRadar"""
        try:
            self.base_url = config.endpoint
            self.session = requests.Session()
            
            # Authenticate with QRadar
            auth_url = f"{self.base_url}/api/auth/login"
            auth_data = {
                'username': config.credentials['username'],
                'password': config.credentials['password']
            }
            
            response = self.session.post(auth_url, json=auth_data)
            if response.status_code == 200:
                self.auth_token = response.json().get('token')
                logger.info("Successfully connected to QRadar")
                return True
            else:
                logger.error(f"Failed to connect to QRadar: {response.status_code}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to QRadar: {e}")
            return False
    
    async def send_event(self, event: SIEMEvent) -> bool:
        """Send event to QRadar"""
        try:
            url = f"{self.base_url}/api/events"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            # Format event for QRadar
            qradar_event = {
                'timestamp': int(time.time() * 1000),
                'source': event.source,
                'event_type': event.event_type,
                'severity': event.severity,
                'description': event.description,
                'raw_data': event.raw_data
            }
            
            response = self.session.post(url, headers=headers, json=qradar_event)
            return response.status_code == 201
            
        except Exception as e:
            logger.error(f"Error sending event to QRadar: {e}")
            return False
    
    async def query_events(self, query: str, time_range: Tuple[str, str]) -> List[SIEMEvent]:
        """Query events from QRadar"""
        try:
            url = f"{self.base_url}/api/events/search"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            data = {
                'query': query,
                'start_time': time_range[0],
                'end_time': time_range[1]
            }
            
            response = self.session.post(url, headers=headers, json=data)
            if response.status_code == 200:
                events = []
                for event_data in response.json().get('events', []):
                    events.append(SIEMEvent(
                        event_id=event_data.get('id', ''),
                        timestamp=event_data.get('timestamp', ''),
                        source=event_data.get('source', ''),
                        event_type=event_data.get('event_type', ''),
                        severity=event_data.get('severity', 'medium'),
                        description=event_data.get('description', ''),
                        raw_data=event_data,
                        normalized_data=event_data
                    ))
                return events
            else:
                logger.error(f"Failed to query QRadar: {response.status_code}")
                return []
                
        except Exception as e:
            logger.error(f"Error querying QRadar: {e}")
            return []
    
    async def get_alerts(self, filters: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Get alerts from QRadar"""
        try:
            # Mock alert retrieval
            return [
                {
                    'alert_id': 'alert_002',
                    'title': 'Suspicious network activity',
                    'severity': 'medium',
                    'timestamp': datetime.utcnow().isoformat(),
                    'source': 'qradar_network_monitor'
                }
            ]
        except Exception as e:
            logger.error(f"Error getting alerts from QRadar: {e}")
            return []

class AzureSentinelConnector(SIEMConnector):
    """Azure Sentinel SIEM connector"""
    
    def __init__(self):
        self.session = None
        self.base_url = None
        self.auth_token = None
    
    async def connect(self, config: IntegrationConfig) -> bool:
        """Connect to Azure Sentinel"""
        try:
            self.base_url = config.endpoint
            self.session = requests.Session()
            
            # Authenticate with Azure Sentinel
            auth_url = f"https://login.microsoftonline.com/{config.credentials['tenant_id']}/oauth2/v2.0/token"
            auth_data = {
                'client_id': config.credentials['client_id'],
                'client_secret': config.credentials['client_secret'],
                'scope': 'https://management.azure.com/.default',
                'grant_type': 'client_credentials'
            }
            
            response = self.session.post(auth_url, data=auth_data)
            if response.status_code == 200:
                self.auth_token = response.json().get('access_token')
                logger.info("Successfully connected to Azure Sentinel")
                return True
            else:
                logger.error(f"Failed to connect to Azure Sentinel: {response.status_code}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to Azure Sentinel: {e}")
            return False
    
    async def send_event(self, event: SIEMEvent) -> bool:
        """Send event to Azure Sentinel"""
        try:
            url = f"{self.base_url}/api/logs"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            # Format event for Azure Sentinel
            sentinel_event = {
                'timestamp': event.timestamp,
                'source': event.source,
                'event_type': event.event_type,
                'severity': event.severity,
                'description': event.description,
                'raw_data': event.raw_data
            }
            
            response = self.session.post(url, headers=headers, json=sentinel_event)
            return response.status_code == 200
            
        except Exception as e:
            logger.error(f"Error sending event to Azure Sentinel: {e}")
            return False
    
    async def query_events(self, query: str, time_range: Tuple[str, str]) -> List[SIEMEvent]:
        """Query events from Azure Sentinel"""
        try:
            url = f"{self.base_url}/api/query"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            data = {
                'query': query,
                'start_time': time_range[0],
                'end_time': time_range[1]
            }
            
            response = self.session.post(url, headers=headers, json=data)
            if response.status_code == 200:
                events = []
                for event_data in response.json().get('events', []):
                    events.append(SIEMEvent(
                        event_id=event_data.get('id', ''),
                        timestamp=event_data.get('timestamp', ''),
                        source=event_data.get('source', ''),
                        event_type=event_data.get('event_type', ''),
                        severity=event_data.get('severity', 'medium'),
                        description=event_data.get('description', ''),
                        raw_data=event_data,
                        normalized_data=event_data
                    ))
                return events
            else:
                logger.error(f"Failed to query Azure Sentinel: {response.status_code}")
                return []
                
        except Exception as e:
            logger.error(f"Error querying Azure Sentinel: {e}")
            return []
    
    async def get_alerts(self, filters: Dict[str, Any]) -> List[Dict[str, Any]]:
        """Get alerts from Azure Sentinel"""
        try:
            # Mock alert retrieval
            return [
                {
                    'alert_id': 'alert_003',
                    'title': 'Azure security center alert',
                    'severity': 'high',
                    'timestamp': datetime.utcnow().isoformat(),
                    'source': 'azure_security_center'
                }
            ]
        except Exception as e:
            logger.error(f"Error getting alerts from Azure Sentinel: {e}")
            return []

class SOARConnector(ABC):
    """Abstract base class for SOAR connectors"""
    
    @abstractmethod
    async def connect(self, config: IntegrationConfig) -> bool:
        """Connect to SOAR system"""
        pass
    
    @abstractmethod
    async def execute_playbook(self, playbook: SOARPlaybook, context: Dict[str, Any]) -> bool:
        """Execute SOAR playbook"""
        pass
    
    @abstractmethod
    async def get_playbooks(self) -> List[SOARPlaybook]:
        """Get available playbooks"""
        pass

class CortexXSOARConnector(SOARConnector):
    """Palo Alto Cortex XSOAR connector"""
    
    def __init__(self):
        self.session = None
        self.base_url = None
        self.api_key = None
    
    async def connect(self, config: IntegrationConfig) -> bool:
        """Connect to Cortex XSOAR"""
        try:
            self.base_url = config.endpoint
            self.api_key = config.credentials['api_key']
            self.session = requests.Session()
            
            # Test connection
            url = f"{self.base_url}/api/v2/system/info"
            headers = {
                'Authorization': f'Bearer {self.api_key}',
                'Content-Type': 'application/json'
            }
            
            response = self.session.get(url, headers=headers)
            if response.status_code == 200:
                logger.info("Successfully connected to Cortex XSOAR")
                return True
            else:
                logger.error(f"Failed to connect to Cortex XSOAR: {response.status_code}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to Cortex XSOAR: {e}")
            return False
    
    async def execute_playbook(self, playbook: SOARPlaybook, context: Dict[str, Any]) -> bool:
        """Execute SOAR playbook"""
        try:
            url = f"{self.base_url}/api/v2/playbook/execute"
            headers = {
                'Authorization': f'Bearer {self.api_key}',
                'Content-Type': 'application/json'
            }
            
            data = {
                'playbook_id': playbook.playbook_id,
                'context': context
            }
            
            response = self.session.post(url, headers=headers, json=data)
            return response.status_code == 200
            
        except Exception as e:
            logger.error(f"Error executing playbook: {e}")
            return False
    
    async def get_playbooks(self) -> List[SOARPlaybook]:
        """Get available playbooks"""
        try:
            url = f"{self.base_url}/api/v2/playbook/list"
            headers = {
                'Authorization': f'Bearer {self.api_key}',
                'Content-Type': 'application/json'
            }
            
            response = self.session.get(url, headers=headers)
            if response.status_code == 200:
                playbooks = []
                for pb_data in response.json().get('playbooks', []):
                    playbooks.append(SOARPlaybook(
                        playbook_id=pb_data.get('id', ''),
                        name=pb_data.get('name', ''),
                        description=pb_data.get('description', ''),
                        trigger_conditions=pb_data.get('triggers', []),
                        actions=[],  # Would be populated from actual data
                        status=pb_data.get('status', 'inactive'),
                        last_executed=pb_data.get('last_executed'),
                        success_rate=pb_data.get('success_rate', 0.0),
                        execution_count=pb_data.get('execution_count', 0)
                    ))
                return playbooks
            else:
                logger.error(f"Failed to get playbooks: {response.status_code}")
                return []
                
        except Exception as e:
            logger.error(f"Error getting playbooks: {e}")
            return []

class SIEMSOARIntegrationHub:
    """Main integration hub for SIEM/SOAR systems"""
    
    def __init__(self):
        self.siem_connectors = {}
        self.soar_connectors = {}
        self.integration_configs = {}
        self.event_cache = {}
        self.cache_ttl = 3600  # 1 hour cache TTL
        
        # Initialize connectors
        self._initialize_connectors()
    
    def _initialize_connectors(self):
        """Initialize SIEM and SOAR connectors"""
        try:
            # SIEM connectors
            self.siem_connectors['splunk'] = SplunkConnector()
            self.siem_connectors['qradar'] = QRadarConnector()
            self.siem_connectors['azure_sentinel'] = AzureSentinelConnector()
            
            # SOAR connectors
            self.soar_connectors['cortex_xsoar'] = CortexXSOARConnector()
            
            logger.info("SIEM/SOAR integration hub initialized successfully")
            
        except Exception as e:
            logger.error(f"Error initializing connectors: {e}")
    
    async def register_integration(self, config: IntegrationConfig) -> bool:
        """Register a new integration"""
        try:
            self.integration_configs[config.integration_id] = config
            
            # Connect to the system
            if config.type == 'siem':
                connector = self.siem_connectors.get(config.name.lower().replace(' ', '_'))
                if connector:
                    success = await connector.connect(config)
                    if success:
                        config.status = 'active'
                        logger.info(f"Successfully registered {config.name} integration")
                        return True
                    else:
                        config.status = 'error'
                        logger.error(f"Failed to connect to {config.name}")
                        return False
                else:
                    logger.error(f"No connector found for {config.name}")
                    return False
            
            elif config.type == 'soar':
                connector = self.soar_connectors.get(config.name.lower().replace(' ', '_'))
                if connector:
                    success = await connector.connect(config)
                    if success:
                        config.status = 'active'
                        logger.info(f"Successfully registered {config.name} integration")
                        return True
                    else:
                        config.status = 'error'
                        logger.error(f"Failed to connect to {config.name}")
                        return False
                else:
                    logger.error(f"No connector found for {config.name}")
                    return False
            
            return False
            
        except Exception as e:
            logger.error(f"Error registering integration: {e}")
            return False
    
    async def send_event_to_siem(self, event: SIEMEvent, siem_type: str = 'all') -> Dict[str, bool]:
        """Send event to SIEM systems"""
        try:
            results = {}
            
            for integration_id, config in self.integration_configs.items():
                if config.type == 'siem' and config.status == 'active':
                    if siem_type == 'all' or config.name.lower() == siem_type.lower():
                        connector = self.siem_connectors.get(config.name.lower().replace(' ', '_'))
                        if connector:
                            success = await connector.send_event(event)
                            results[integration_id] = success
                            logger.info(f"Sent event to {config.name}: {success}")
            
            return results
            
        except Exception as e:
            logger.error(f"Error sending event to SIEM: {e}")
            return {}
    
    async def query_events_from_siem(self, query: str, time_range: Tuple[str, str], siem_type: str = 'all') -> Dict[str, List[SIEMEvent]]:
        """Query events from SIEM systems"""
        try:
            results = {}
            
            for integration_id, config in self.integration_configs.items():
                if config.type == 'siem' and config.status == 'active':
                    if siem_type == 'all' or config.name.lower() == siem_type.lower():
                        connector = self.siem_connectors.get(config.name.lower().replace(' ', '_'))
                        if connector:
                            events = await connector.query_events(query, time_range)
                            results[integration_id] = events
                            logger.info(f"Queried {len(events)} events from {config.name}")
            
            return results
            
        except Exception as e:
            logger.error(f"Error querying events from SIEM: {e}")
            return {}
    
    async def get_alerts_from_siem(self, filters: Dict[str, Any], siem_type: str = 'all') -> Dict[str, List[Dict[str, Any]]]:
        """Get alerts from SIEM systems"""
        try:
            results = {}
            
            for integration_id, config in self.integration_configs.items():
                if config.type == 'siem' and config.status == 'active':
                    if siem_type == 'all' or config.name.lower() == siem_type.lower():
                        connector = self.siem_connectors.get(config.name.lower().replace(' ', '_'))
                        if connector:
                            alerts = await connector.get_alerts(filters)
                            results[integration_id] = alerts
                            logger.info(f"Retrieved {len(alerts)} alerts from {config.name}")
            
            return results
            
        except Exception as e:
            logger.error(f"Error getting alerts from SIEM: {e}")
            return {}
    
    async def execute_soar_playbook(self, playbook_id: str, context: Dict[str, Any], soar_type: str = 'all') -> Dict[str, bool]:
        """Execute SOAR playbook"""
        try:
            results = {}
            
            for integration_id, config in self.integration_configs.items():
                if config.type == 'soar' and config.status == 'active':
                    if soar_type == 'all' or config.name.lower() == soar_type.lower():
                        connector = self.soar_connectors.get(config.name.lower().replace(' ', '_'))
                        if connector:
                            # Get playbook
                            playbooks = await connector.get_playbooks()
                            playbook = next((pb for pb in playbooks if pb.playbook_id == playbook_id), None)
                            
                            if playbook:
                                success = await connector.execute_playbook(playbook, context)
                                results[integration_id] = success
                                logger.info(f"Executed playbook {playbook_id} on {config.name}: {success}")
                            else:
                                logger.error(f"Playbook {playbook_id} not found on {config.name}")
                                results[integration_id] = False
            
            return results
            
        except Exception as e:
            logger.error(f"Error executing SOAR playbook: {e}")
            return {}
    
    async def get_soar_playbooks(self, soar_type: str = 'all') -> Dict[str, List[SOARPlaybook]]:
        """Get SOAR playbooks"""
        try:
            results = {}
            
            for integration_id, config in self.integration_configs.items():
                if config.type == 'soar' and config.status == 'active':
                    if soar_type == 'all' or config.name.lower() == soar_type.lower():
                        connector = self.soar_connectors.get(config.name.lower().replace(' ', '_'))
                        if connector:
                            playbooks = await connector.get_playbooks()
                            results[integration_id] = playbooks
                            logger.info(f"Retrieved {len(playbooks)} playbooks from {config.name}")
            
            return results
            
        except Exception as e:
            logger.error(f"Error getting SOAR playbooks: {e}")
            return {}
    
    async def correlate_events(self, events: List[SIEMEvent]) -> List[Dict[str, Any]]:
        """Correlate events across SIEM systems"""
        try:
            correlations = []
            
            # Group events by correlation criteria
            event_groups = {}
            for event in events:
                # Simple correlation by source and event type
                key = f"{event.source}_{event.event_type}"
                if key not in event_groups:
                    event_groups[key] = []
                event_groups[key].append(event)
            
            # Analyze correlations
            for key, group_events in event_groups.items():
                if len(group_events) > 1:
                    correlation = {
                        'correlation_id': f"corr_{int(time.time())}",
                        'event_count': len(group_events),
                        'source': group_events[0].source,
                        'event_type': group_events[0].event_type,
                        'severity': max(event.severity for event in group_events),
                        'time_range': {
                            'start': min(event.timestamp for event in group_events),
                            'end': max(event.timestamp for event in group_events)
                        },
                        'events': [event.event_id for event in group_events],
                        'correlation_score': len(group_events) * 0.2
                    }
                    correlations.append(correlation)
            
            return correlations
            
        except Exception as e:
            logger.error(f"Error correlating events: {e}")
            return []
    
    async def generate_integration_report(self) -> Dict[str, Any]:
        """Generate integration status report"""
        try:
            report = {
                'total_integrations': len(self.integration_configs),
                'active_integrations': len([c for c in self.integration_configs.values() if c.status == 'active']),
                'siem_integrations': len([c for c in self.integration_configs.values() if c.type == 'siem']),
                'soar_integrations': len([c for c in self.integration_configs.values() if c.type == 'soar']),
                'integration_details': []
            }
            
            for integration_id, config in self.integration_configs.items():
                report['integration_details'].append({
                    'integration_id': integration_id,
                    'name': config.name,
                    'type': config.type,
                    'status': config.status,
                    'last_sync': config.last_sync
                })
            
            return report
            
        except Exception as e:
            logger.error(f"Error generating integration report: {e}")
            return {}

# Global SIEM/SOAR integration hub instance
siem_soar_hub = SIEMSOARIntegrationHub()

