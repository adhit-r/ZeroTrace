"""
Attack Surface Mapping Service - Comprehensive security posture analysis
"""

import asyncio
import json
import logging
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass
from sklearn.cluster import DBSCAN
from sklearn.preprocessing import StandardScaler
import networkx as nx
import matplotlib.pyplot as plt
import seaborn as sns

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class AttackSurfaceMap:
    """Represents a comprehensive attack surface map"""
    organization_id: str
    map_id: str
    external_exposure: ExternalExposure
    internal_vulnerabilities: InternalVulnerabilities
    supply_chain_exposure: SupplyChainExposure
    attack_vectors: List[AttackVector]
    critical_paths: List[CriticalPath]
    risk_hotspots: List[RiskHotspot]
    defense_coverage: DefenseCoverage
    blind_spots: List[BlindSpot]
    recommendations: List[str]
    generated_at: str
    confidence_score: float

@dataclass
class ExternalExposure:
    """Represents external attack surface exposure"""
    public_assets: List[PublicAsset]
    exposed_services: List[ExposedService]
    open_ports: List[OpenPort]
    web_applications: List[WebApplication]
    api_endpoints: List[APIEndpoint]
    cloud_resources: List[CloudResource]
    dns_records: List[DNSRecord]
    ssl_certificates: List[SSLCertificate]
    exposure_score: float
    risk_level: str

@dataclass
class InternalVulnerabilities:
    """Represents internal vulnerability landscape"""
    network_segments: List[NetworkSegment]
    internal_services: List[InternalService]
    database_exposures: List[DatabaseExposure]
    file_shares: List[FileShare]
    internal_apis: List[InternalAPI]
    privilege_escalation_paths: List[PrivilegeEscalationPath]
    lateral_movement_paths: List[LateralMovementPath]
    vulnerability_score: float
    risk_level: str

@dataclass
class SupplyChainExposure:
    """Represents supply chain attack surface"""
    dependencies: List[Dependency]
    third_party_integrations: List[ThirdPartyIntegration]
    vendor_connections: List[VendorConnection]
    license_risks: List[LicenseRisk]
    malicious_packages: List[MaliciousPackage]
    transitive_vulnerabilities: List[TransitiveVulnerability]
    supply_chain_score: float
    risk_level: str

@dataclass
class AttackVector:
    """Represents a potential attack vector"""
    vector_id: str
    name: str
    type: str  # external, internal, supply_chain
    entry_point: str
    target: str
    likelihood: float
    impact: float
    complexity: str  # low, medium, high
    prerequisites: List[str]
    mitigation_controls: List[str]
    detection_difficulty: str
    exploitation_time: str

@dataclass
class CriticalPath:
    """Represents a critical attack path"""
    path_id: str
    name: str
    steps: List[AttackStep]
    total_likelihood: float
    total_impact: float
    criticality_score: float
    mitigation_priority: str
    detection_points: List[str]
    prevention_controls: List[str]

@dataclass
class AttackStep:
    """Represents a step in an attack path"""
    step_number: int
    action: str
    target: str
    technique: str
    likelihood: float
    impact: float
    detection_difficulty: str
    mitigation_controls: List[str]

@dataclass
class RiskHotspot:
    """Represents a high-risk area in the attack surface"""
    hotspot_id: str
    name: str
    location: str
    risk_score: float
    vulnerability_count: int
    exposure_level: str
    attack_vectors: List[str]
    critical_assets: List[str]
    immediate_actions: List[str]
    long_term_mitigations: List[str]

@dataclass
class DefenseCoverage:
    """Represents defense coverage analysis"""
    coverage_percentage: float
    protected_assets: int
    unprotected_assets: int
    defense_gaps: List[DefenseGap]
    security_controls: List[SecurityControl]
    monitoring_coverage: float
    incident_response_readiness: float
    overall_defense_score: float

@dataclass
class DefenseGap:
    """Represents a gap in defense coverage"""
    gap_id: str
    description: str
    affected_assets: List[str]
    risk_level: str
    recommended_controls: List[str]
    implementation_priority: str
    estimated_cost: str
    implementation_timeline: str

@dataclass
class SecurityControl:
    """Represents a security control"""
    control_id: str
    name: str
    type: str  # preventive, detective, responsive
    coverage_area: str
    effectiveness: float
    implementation_status: str
    maintenance_required: bool
    cost: str

@dataclass
class BlindSpot:
    """Represents a blind spot in security monitoring"""
    blind_spot_id: str
    description: str
    affected_area: str
    risk_implications: List[str]
    monitoring_recommendations: List[str]
    detection_techniques: List[str]
    priority: str

# Supporting data classes
@dataclass
class PublicAsset:
    ip_address: str
    hostname: str
    services: List[str]
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class ExposedService:
    service_name: str
    port: int
    protocol: str
    version: str
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class OpenPort:
    port: int
    protocol: str
    service: str
    state: str
    risk_level: str

@dataclass
class WebApplication:
    url: str
    technology_stack: List[str]
    vulnerabilities: List[str]
    authentication_required: bool
    risk_score: float

@dataclass
class APIEndpoint:
    endpoint: str
    method: str
    authentication: str
    rate_limiting: bool
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class CloudResource:
    resource_type: str
    provider: str
    region: str
    configuration: Dict[str, Any]
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class DNSRecord:
    domain: str
    record_type: str
    value: str
    ttl: int
    risk_indicators: List[str]

@dataclass
class SSLCertificate:
    domain: str
    issuer: str
    expiry_date: str
    key_size: int
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class NetworkSegment:
    segment_name: str
    ip_range: str
    assets: List[str]
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class InternalService:
    service_name: str
    host: str
    port: int
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class DatabaseExposure:
    database_name: str
    host: str
    port: int
    authentication: str
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class FileShare:
    share_name: str
    path: str
    permissions: str
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class InternalAPI:
    api_name: str
    endpoint: str
    authentication: str
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class PrivilegeEscalationPath:
    path_id: str
    start_user: str
    target_privileges: str
    escalation_steps: List[str]
    likelihood: float
    impact: float

@dataclass
class LateralMovementPath:
    path_id: str
    start_host: str
    target_hosts: List[str]
    movement_steps: List[str]
    likelihood: float
    impact: float

@dataclass
class Dependency:
    name: str
    version: str
    type: str
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class ThirdPartyIntegration:
    service_name: str
    provider: str
    integration_type: str
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class VendorConnection:
    vendor_name: str
    connection_type: str
    data_shared: List[str]
    vulnerabilities: List[str]
    risk_score: float

@dataclass
class LicenseRisk:
    component_name: str
    license_type: str
    compliance_risk: str
    legal_implications: List[str]
    risk_score: float

@dataclass
class MaliciousPackage:
    package_name: str
    version: str
    threat_type: str
    indicators: List[str]
    risk_score: float

@dataclass
class TransitiveVulnerability:
    direct_dependency: str
    transitive_dependency: str
    vulnerability: str
    risk_score: float

class AttackSurfaceMapper:
    """Service for comprehensive attack surface mapping and analysis"""
    
    def __init__(self):
        self.cache = {}
        self.cache_ttl = 3600  # 1 hour cache TTL
        
    async def generate_attack_surface_map(self, organization_id: str, 
                                        vulnerability_data: List[Dict] = None,
                                        network_data: List[Dict] = None,
                                        asset_data: List[Dict] = None,
                                        organization_context: Dict = None) -> AttackSurfaceMap:
        """
        Generate comprehensive attack surface map for an organization
        
        Args:
            organization_id: Organization identifier
            vulnerability_data: Vulnerability scan results
            network_data: Network discovery data
            asset_data: Asset inventory data
            organization_context: Organization profile data
        
        Returns:
            AttackSurfaceMap with comprehensive analysis
        """
        try:
            logger.info(f"Generating attack surface map for organization {organization_id}")
            
            # Analyze external exposure
            external_exposure = await self._analyze_external_exposure(
                organization_id, network_data, asset_data, vulnerability_data
            )
            
            # Analyze internal vulnerabilities
            internal_vulnerabilities = await self._analyze_internal_vulnerabilities(
                organization_id, vulnerability_data, network_data
            )
            
            # Analyze supply chain exposure
            supply_chain_exposure = await self._analyze_supply_chain_exposure(
                organization_id, vulnerability_data, organization_context
            )
            
            # Identify attack vectors
            attack_vectors = await self._identify_attack_vectors(
                external_exposure, internal_vulnerabilities, supply_chain_exposure
            )
            
            # Find critical attack paths
            critical_paths = await self._find_critical_paths(
                attack_vectors, external_exposure, internal_vulnerabilities
            )
            
            # Identify risk hotspots
            risk_hotspots = await self._identify_risk_hotspots(
                external_exposure, internal_vulnerabilities, supply_chain_exposure
            )
            
            # Analyze defense coverage
            defense_coverage = await self._analyze_defense_coverage(
                external_exposure, internal_vulnerabilities, organization_context
            )
            
            # Identify blind spots
            blind_spots = await self._identify_blind_spots(
                external_exposure, internal_vulnerabilities, defense_coverage
            )
            
            # Generate recommendations
            recommendations = await self._generate_attack_surface_recommendations(
                risk_hotspots, defense_coverage, blind_spots, critical_paths
            )
            
            # Calculate confidence score
            confidence_score = await self._calculate_attack_surface_confidence(
                vulnerability_data, network_data, asset_data
            )
            
            # Create attack surface map
            attack_surface_map = AttackSurfaceMap(
                organization_id=organization_id,
                map_id=f"asm_{organization_id}_{int(datetime.now().timestamp())}",
                external_exposure=external_exposure,
                internal_vulnerabilities=internal_vulnerabilities,
                supply_chain_exposure=supply_chain_exposure,
                attack_vectors=attack_vectors,
                critical_paths=critical_paths,
                risk_hotspots=risk_hotspots,
                defense_coverage=defense_coverage,
                blind_spots=blind_spots,
                recommendations=recommendations,
                generated_at=datetime.utcnow().isoformat(),
                confidence_score=confidence_score
            )
            
            logger.info(f"Completed attack surface map for organization {organization_id}")
            return attack_surface_map
            
        except Exception as e:
            logger.error(f"Error generating attack surface map: {e}")
            return self._get_fallback_attack_surface_map(organization_id)
    
    async def _analyze_external_exposure(self, organization_id: str, network_data: List[Dict],
                                       asset_data: List[Dict], vulnerability_data: List[Dict]) -> ExternalExposure:
        """Analyze external attack surface exposure"""
        try:
            # Mock external exposure analysis
            # In a real implementation, this would analyze actual network scans, DNS records, etc.
            
            public_assets = [
                PublicAsset(
                    ip_address="203.0.113.1",
                    hostname="web.example.com",
                    services=["HTTP", "HTTPS", "SSH"],
                    vulnerabilities=["CVE-2023-1234", "CVE-2023-5678"],
                    risk_score=0.7
                ),
                PublicAsset(
                    ip_address="203.0.113.2",
                    hostname="api.example.com",
                    services=["HTTPS", "API"],
                    vulnerabilities=["CVE-2023-9012"],
                    risk_score=0.5
                )
            ]
            
            exposed_services = [
                ExposedService(
                    service_name="Apache HTTP Server",
                    port=80,
                    protocol="TCP",
                    version="2.4.41",
                    vulnerabilities=["CVE-2023-1234"],
                    risk_score=0.6
                ),
                ExposedService(
                    service_name="OpenSSH",
                    port=22,
                    protocol="TCP",
                    version="8.2",
                    vulnerabilities=[],
                    risk_score=0.3
                )
            ]
            
            open_ports = [
                OpenPort(port=80, protocol="TCP", service="HTTP", state="open", risk_level="medium"),
                OpenPort(port=443, protocol="TCP", service="HTTPS", state="open", risk_level="low"),
                OpenPort(port=22, protocol="TCP", service="SSH", state="open", risk_level="medium"),
                OpenPort(port=3389, protocol="TCP", service="RDP", state="open", risk_level="high")
            ]
            
            web_applications = [
                WebApplication(
                    url="https://web.example.com",
                    technology_stack=["Apache", "PHP", "MySQL"],
                    vulnerabilities=["CVE-2023-1234", "CVE-2023-5678"],
                    authentication_required=False,
                    risk_score=0.8
                )
            ]
            
            api_endpoints = [
                APIEndpoint(
                    endpoint="/api/v1/users",
                    method="GET",
                    authentication="Bearer Token",
                    rate_limiting=True,
                    vulnerabilities=["CVE-2023-9012"],
                    risk_score=0.4
                )
            ]
            
            cloud_resources = [
                CloudResource(
                    resource_type="EC2 Instance",
                    provider="AWS",
                    region="us-east-1",
                    configuration={"instance_type": "t3.medium", "security_groups": ["sg-12345"]},
                    vulnerabilities=["CVE-2023-3456"],
                    risk_score=0.5
                )
            ]
            
            dns_records = [
                DNSRecord(
                    domain="example.com",
                    record_type="A",
                    value="203.0.113.1",
                    ttl=3600,
                    risk_indicators=[]
                )
            ]
            
            ssl_certificates = [
                SSLCertificate(
                    domain="example.com",
                    issuer="Let's Encrypt",
                    expiry_date="2024-12-31",
                    key_size=2048,
                    vulnerabilities=[],
                    risk_score=0.2
                )
            ]
            
            # Calculate exposure score
            total_assets = len(public_assets)
            vulnerable_assets = len([asset for asset in public_assets if asset.vulnerabilities])
            exposure_score = vulnerable_assets / total_assets if total_assets > 0 else 0
            
            # Determine risk level
            if exposure_score > 0.7:
                risk_level = "high"
            elif exposure_score > 0.4:
                risk_level = "medium"
            else:
                risk_level = "low"
            
            return ExternalExposure(
                public_assets=public_assets,
                exposed_services=exposed_services,
                open_ports=open_ports,
                web_applications=web_applications,
                api_endpoints=api_endpoints,
                cloud_resources=cloud_resources,
                dns_records=dns_records,
                ssl_certificates=ssl_certificates,
                exposure_score=exposure_score,
                risk_level=risk_level
            )
            
        except Exception as e:
            logger.error(f"Error analyzing external exposure: {e}")
            return ExternalExposure(
                public_assets=[],
                exposed_services=[],
                open_ports=[],
                web_applications=[],
                api_endpoints=[],
                cloud_resources=[],
                dns_records=[],
                ssl_certificates=[],
                exposure_score=0.0,
                risk_level="unknown"
            )
    
    async def _analyze_internal_vulnerabilities(self, organization_id: str, 
                                              vulnerability_data: List[Dict], 
                                              network_data: List[Dict]) -> InternalVulnerabilities:
        """Analyze internal vulnerability landscape"""
        try:
            # Mock internal vulnerability analysis
            network_segments = [
                NetworkSegment(
                    segment_name="DMZ",
                    ip_range="192.168.1.0/24",
                    assets=["web-server-01", "web-server-02"],
                    vulnerabilities=["CVE-2023-1234"],
                    risk_score=0.6
                ),
                NetworkSegment(
                    segment_name="Internal",
                    ip_range="192.168.2.0/24",
                    assets=["db-server-01", "app-server-01"],
                    vulnerabilities=["CVE-2023-5678", "CVE-2023-9012"],
                    risk_score=0.8
                )
            ]
            
            internal_services = [
                InternalService(
                    service_name="MySQL",
                    host="192.168.2.10",
                    port=3306,
                    vulnerabilities=["CVE-2023-3456"],
                    risk_score=0.7
                )
            ]
            
            database_exposures = [
                DatabaseExposure(
                    database_name="user_db",
                    host="192.168.2.10",
                    port=3306,
                    authentication="password",
                    vulnerabilities=["CVE-2023-3456"],
                    risk_score=0.8
                )
            ]
            
            file_shares = [
                FileShare(
                    share_name="shared_docs",
                    path="\\\\fileserver\\shared",
                    permissions="everyone:full",
                    vulnerabilities=["CVE-2023-7890"],
                    risk_score=0.9
                )
            ]
            
            internal_apis = [
                InternalAPI(
                    api_name="internal_auth",
                    endpoint="/auth/internal",
                    authentication="API Key",
                    vulnerabilities=["CVE-2023-2345"],
                    risk_score=0.5
                )
            ]
            
            privilege_escalation_paths = [
                PrivilegeEscalationPath(
                    path_id="pe_001",
                    start_user="user",
                    target_privileges="admin",
                    escalation_steps=["Exploit CVE-2023-1234", "Access admin panel"],
                    likelihood=0.6,
                    impact=0.9
                )
            ]
            
            lateral_movement_paths = [
                LateralMovementPath(
                    path_id="lm_001",
                    start_host="web-server-01",
                    target_hosts=["db-server-01", "app-server-01"],
                    movement_steps=["Exploit web vulnerability", "Access internal network"],
                    likelihood=0.5,
                    impact=0.8
                )
            ]
            
            # Calculate vulnerability score
            total_internal_assets = len(network_segments) + len(internal_services) + len(database_exposures)
            vulnerable_internal_assets = len([seg for seg in network_segments if seg.vulnerabilities]) + \
                                       len([svc for svc in internal_services if svc.vulnerabilities]) + \
                                       len([db for db in database_exposures if db.vulnerabilities])
            
            vulnerability_score = vulnerable_internal_assets / total_internal_assets if total_internal_assets > 0 else 0
            
            # Determine risk level
            if vulnerability_score > 0.7:
                risk_level = "high"
            elif vulnerability_score > 0.4:
                risk_level = "medium"
            else:
                risk_level = "low"
            
            return InternalVulnerabilities(
                network_segments=network_segments,
                internal_services=internal_services,
                database_exposures=database_exposures,
                file_shares=file_shares,
                internal_apis=internal_apis,
                privilege_escalation_paths=privilege_escalation_paths,
                lateral_movement_paths=lateral_movement_paths,
                vulnerability_score=vulnerability_score,
                risk_level=risk_level
            )
            
        except Exception as e:
            logger.error(f"Error analyzing internal vulnerabilities: {e}")
            return InternalVulnerabilities(
                network_segments=[],
                internal_services=[],
                database_exposures=[],
                file_shares=[],
                internal_apis=[],
                privilege_escalation_paths=[],
                lateral_movement_paths=[],
                vulnerability_score=0.0,
                risk_level="unknown"
            )
    
    async def _analyze_supply_chain_exposure(self, organization_id: str, 
                                           vulnerability_data: List[Dict],
                                           organization_context: Dict) -> SupplyChainExposure:
        """Analyze supply chain attack surface"""
        try:
            # Mock supply chain analysis
            dependencies = [
                Dependency(
                    name="express",
                    version="4.18.2",
                    type="npm",
                    vulnerabilities=["CVE-2023-1234"],
                    risk_score=0.6
                ),
                Dependency(
                    name="lodash",
                    version="4.17.21",
                    type="npm",
                    vulnerabilities=[],
                    risk_score=0.2
                )
            ]
            
            third_party_integrations = [
                ThirdPartyIntegration(
                    service_name="Stripe",
                    provider="Stripe Inc.",
                    integration_type="Payment Processing",
                    vulnerabilities=[],
                    risk_score=0.3
                ),
                ThirdPartyIntegration(
                    service_name="AWS S3",
                    provider="Amazon",
                    integration_type="Cloud Storage",
                    vulnerabilities=["CVE-2023-5678"],
                    risk_score=0.5
                )
            ]
            
            vendor_connections = [
                VendorConnection(
                    vendor_name="DataCorp",
                    connection_type="Data Sharing",
                    data_shared=["Customer Data", "Analytics"],
                    vulnerabilities=["CVE-2023-9012"],
                    risk_score=0.7
                )
            ]
            
            license_risks = [
                LicenseRisk(
                    component_name="jQuery",
                    license_type="MIT",
                    compliance_risk="low",
                    legal_implications=[],
                    risk_score=0.1
                )
            ]
            
            malicious_packages = [
                MaliciousPackage(
                    package_name="malicious-package",
                    version="1.0.0",
                    threat_type="Backdoor",
                    indicators=["Suspicious network activity", "Data exfiltration"],
                    risk_score=0.9
                )
            ]
            
            transitive_vulnerabilities = [
                TransitiveVulnerability(
                    direct_dependency="express",
                    transitive_dependency="qs",
                    vulnerability="CVE-2023-3456",
                    risk_score=0.4
                )
            ]
            
            # Calculate supply chain score
            total_components = len(dependencies) + len(third_party_integrations) + len(vendor_connections)
            vulnerable_components = len([dep for dep in dependencies if dep.vulnerabilities]) + \
                                 len([tpi for tpi in third_party_integrations if tpi.vulnerabilities]) + \
                                 len([vc for vc in vendor_connections if vc.vulnerabilities])
            
            supply_chain_score = vulnerable_components / total_components if total_components > 0 else 0
            
            # Determine risk level
            if supply_chain_score > 0.7:
                risk_level = "high"
            elif supply_chain_score > 0.4:
                risk_level = "medium"
            else:
                risk_level = "low"
            
            return SupplyChainExposure(
                dependencies=dependencies,
                third_party_integrations=third_party_integrations,
                vendor_connections=vendor_connections,
                license_risks=license_risks,
                malicious_packages=malicious_packages,
                transitive_vulnerabilities=transitive_vulnerabilities,
                supply_chain_score=supply_chain_score,
                risk_level=risk_level
            )
            
        except Exception as e:
            logger.error(f"Error analyzing supply chain exposure: {e}")
            return SupplyChainExposure(
                dependencies=[],
                third_party_integrations=[],
                vendor_connections=[],
                license_risks=[],
                malicious_packages=[],
                transitive_vulnerabilities=[],
                supply_chain_score=0.0,
                risk_level="unknown"
            )
    
    async def _identify_attack_vectors(self, external_exposure: ExternalExposure,
                                     internal_vulnerabilities: InternalVulnerabilities,
                                     supply_chain_exposure: SupplyChainExposure) -> List[AttackVector]:
        """Identify potential attack vectors"""
        try:
            attack_vectors = []
            
            # External attack vectors
            for asset in external_exposure.public_assets:
                if asset.vulnerabilities:
                    attack_vectors.append(AttackVector(
                        vector_id=f"ext_{asset.ip_address}",
                        name=f"External attack on {asset.hostname}",
                        type="external",
                        entry_point=asset.ip_address,
                        target=asset.hostname,
                        likelihood=0.7,
                        impact=0.8,
                        complexity="medium",
                        prerequisites=["Internet access", "Vulnerability knowledge"],
                        mitigation_controls=["Firewall rules", "Patch management"],
                        detection_difficulty="medium",
                        exploitation_time="hours to days"
                    ))
            
            # Internal attack vectors
            for segment in internal_vulnerabilities.network_segments:
                if segment.vulnerabilities:
                    attack_vectors.append(AttackVector(
                        vector_id=f"int_{segment.segment_name}",
                        name=f"Internal attack on {segment.segment_name}",
                        type="internal",
                        entry_point="Internal network",
                        target=segment.segment_name,
                        likelihood=0.5,
                        impact=0.9,
                        complexity="high",
                        prerequisites=["Internal access", "Network knowledge"],
                        mitigation_controls=["Network segmentation", "Access controls"],
                        detection_difficulty="high",
                        exploitation_time="days to weeks"
                    ))
            
            # Supply chain attack vectors
            for dependency in supply_chain_exposure.dependencies:
                if dependency.vulnerabilities:
                    attack_vectors.append(AttackVector(
                        vector_id=f"sc_{dependency.name}",
                        name=f"Supply chain attack via {dependency.name}",
                        type="supply_chain",
                        entry_point="Dependency",
                        target=dependency.name,
                        likelihood=0.6,
                        impact=0.7,
                        complexity="low",
                        prerequisites=["Package installation", "Vulnerability knowledge"],
                        mitigation_controls=["Dependency scanning", "Version pinning"],
                        detection_difficulty="low",
                        exploitation_time="minutes to hours"
                    ))
            
            return attack_vectors
            
        except Exception as e:
            logger.error(f"Error identifying attack vectors: {e}")
            return []
    
    async def _find_critical_paths(self, attack_vectors: List[AttackVector],
                                 external_exposure: ExternalExposure,
                                 internal_vulnerabilities: InternalVulnerabilities) -> List[CriticalPath]:
        """Find critical attack paths"""
        try:
            critical_paths = []
            
            # Mock critical path analysis
            # In a real implementation, this would use graph algorithms to find attack paths
            
            if attack_vectors:
                # Create a critical path from external to internal
                critical_paths.append(CriticalPath(
                    path_id="cp_001",
                    name="External to Internal Data Exfiltration",
                    steps=[
                        AttackStep(
                            step_number=1,
                            action="Exploit external vulnerability",
                            target="Web server",
                            technique="SQL Injection",
                            likelihood=0.7,
                            impact=0.6,
                            detection_difficulty="medium",
                            mitigation_controls=["Input validation", "WAF"]
                        ),
                        AttackStep(
                            step_number=2,
                            action="Lateral movement",
                            target="Database server",
                            technique="Privilege escalation",
                            likelihood=0.5,
                            impact=0.8,
                            detection_difficulty="high",
                            mitigation_controls=["Network segmentation", "Access controls"]
                        ),
                        AttackStep(
                            step_number=3,
                            action="Data exfiltration",
                            target="Sensitive data",
                            technique="Data export",
                            likelihood=0.4,
                            impact=0.9,
                            detection_difficulty="high",
                            mitigation_controls=["Data loss prevention", "Monitoring"]
                        )
                    ],
                    total_likelihood=0.14,  # 0.7 * 0.5 * 0.4
                    total_impact=0.9,
                    criticality_score=0.126,  # 0.14 * 0.9
                    mitigation_priority="high",
                    detection_points=["Web server logs", "Database logs", "Network traffic"],
                    prevention_controls=["Input validation", "Network segmentation", "Monitoring"]
                ))
            
            return critical_paths
            
        except Exception as e:
            logger.error(f"Error finding critical paths: {e}")
            return []
    
    async def _identify_risk_hotspots(self, external_exposure: ExternalExposure,
                                    internal_vulnerabilities: InternalVulnerabilities,
                                    supply_chain_exposure: SupplyChainExposure) -> List[RiskHotspot]:
        """Identify risk hotspots in the attack surface"""
        try:
            hotspots = []
            
            # External hotspots
            high_risk_external_assets = [asset for asset in external_exposure.public_assets if asset.risk_score > 0.7]
            if high_risk_external_assets:
                hotspots.append(RiskHotspot(
                    hotspot_id="hs_ext_001",
                    name="High-Risk External Assets",
                    location="External network",
                    risk_score=0.8,
                    vulnerability_count=len([v for asset in high_risk_external_assets for v in asset.vulnerabilities]),
                    exposure_level="high",
                    attack_vectors=["External exploitation", "Data breach"],
                    critical_assets=[asset.hostname for asset in high_risk_external_assets],
                    immediate_actions=["Patch vulnerabilities", "Implement WAF", "Review access controls"],
                    long_term_mitigations=["Security architecture review", "Regular penetration testing"]
                ))
            
            # Internal hotspots
            high_risk_internal_segments = [seg for seg in internal_vulnerabilities.network_segments if seg.risk_score > 0.7]
            if high_risk_internal_segments:
                hotspots.append(RiskHotspot(
                    hotspot_id="hs_int_001",
                    name="High-Risk Internal Segments",
                    location="Internal network",
                    risk_score=0.9,
                    vulnerability_count=len([v for seg in high_risk_internal_segments for v in seg.vulnerabilities]),
                    exposure_level="critical",
                    attack_vectors=["Lateral movement", "Privilege escalation"],
                    critical_assets=[seg.segment_name for seg in high_risk_internal_segments],
                    immediate_actions=["Network segmentation", "Access control review", "Vulnerability patching"],
                    long_term_mitigations=["Zero trust architecture", "Continuous monitoring"]
                ))
            
            # Supply chain hotspots
            high_risk_dependencies = [dep for dep in supply_chain_exposure.dependencies if dep.risk_score > 0.7]
            if high_risk_dependencies:
                hotspots.append(RiskHotspot(
                    hotspot_id="hs_sc_001",
                    name="High-Risk Dependencies",
                    location="Supply chain",
                    risk_score=0.7,
                    vulnerability_count=len([v for dep in high_risk_dependencies for v in dep.vulnerabilities]),
                    exposure_level="high",
                    attack_vectors=["Supply chain attack", "Dependency exploitation"],
                    critical_assets=[dep.name for dep in high_risk_dependencies],
                    immediate_actions=["Update dependencies", "Review third-party access"],
                    long_term_mitigations=["Dependency management", "Supply chain security"]
                ))
            
            return hotspots
            
        except Exception as e:
            logger.error(f"Error identifying risk hotspots: {e}")
            return []
    
    async def _analyze_defense_coverage(self, external_exposure: ExternalExposure,
                                      internal_vulnerabilities: InternalVulnerabilities,
                                      organization_context: Dict) -> DefenseCoverage:
        """Analyze defense coverage"""
        try:
            # Mock defense coverage analysis
            total_assets = len(external_exposure.public_assets) + len(internal_vulnerabilities.network_segments)
            protected_assets = int(total_assets * 0.7)  # 70% coverage
            unprotected_assets = total_assets - protected_assets
            
            defense_gaps = [
                DefenseGap(
                    gap_id="gap_001",
                    description="Unprotected internal network segments",
                    affected_assets=["Internal segment 1", "Internal segment 2"],
                    risk_level="high",
                    recommended_controls=["Network segmentation", "Access controls"],
                    implementation_priority="high",
                    estimated_cost="$50K - $100K",
                    implementation_timeline="3-6 months"
                )
            ]
            
            security_controls = [
                SecurityControl(
                    control_id="ctrl_001",
                    name="Web Application Firewall",
                    type="preventive",
                    coverage_area="External web applications",
                    effectiveness=0.8,
                    implementation_status="implemented",
                    maintenance_required=True,
                    cost="$10K - $20K annually"
                ),
                SecurityControl(
                    control_id="ctrl_002",
                    name="Network Monitoring",
                    type="detective",
                    coverage_area="Internal network",
                    effectiveness=0.7,
                    implementation_status="partial",
                    maintenance_required=True,
                    cost="$30K - $50K annually"
                )
            ]
            
            monitoring_coverage = 0.75
            incident_response_readiness = 0.6
            overall_defense_score = (monitoring_coverage + incident_response_readiness) / 2
            
            return DefenseCoverage(
                coverage_percentage=0.7,
                protected_assets=protected_assets,
                unprotected_assets=unprotected_assets,
                defense_gaps=defense_gaps,
                security_controls=security_controls,
                monitoring_coverage=monitoring_coverage,
                incident_response_readiness=incident_response_readiness,
                overall_defense_score=overall_defense_score
            )
            
        except Exception as e:
            logger.error(f"Error analyzing defense coverage: {e}")
            return DefenseCoverage(
                coverage_percentage=0.0,
                protected_assets=0,
                unprotected_assets=0,
                defense_gaps=[],
                security_controls=[],
                monitoring_coverage=0.0,
                incident_response_readiness=0.0,
                overall_defense_score=0.0
            )
    
    async def _identify_blind_spots(self, external_exposure: ExternalExposure,
                                  internal_vulnerabilities: InternalVulnerabilities,
                                  defense_coverage: DefenseCoverage) -> List[BlindSpot]:
        """Identify blind spots in security monitoring"""
        try:
            blind_spots = []
            
            # Mock blind spot identification
            blind_spots.append(BlindSpot(
                blind_spot_id="bs_001",
                description="Unmonitored internal network segments",
                affected_area="Internal network",
                risk_implications=["Undetected lateral movement", "Privilege escalation"],
                monitoring_recommendations=["Implement network monitoring", "Deploy endpoint detection"],
                detection_techniques=["Network traffic analysis", "Behavioral analytics"],
                priority="high"
            ))
            
            blind_spots.append(BlindSpot(
                blind_spot_id="bs_002",
                description="Third-party integration monitoring gaps",
                affected_area="Supply chain",
                risk_implications=["Undetected supply chain attacks", "Data exfiltration"],
                monitoring_recommendations=["Implement API monitoring", "Deploy third-party risk management"],
                detection_techniques=["API call analysis", "Anomaly detection"],
                priority="medium"
            ))
            
            return blind_spots
            
        except Exception as e:
            logger.error(f"Error identifying blind spots: {e}")
            return []
    
    async def _generate_attack_surface_recommendations(self, risk_hotspots: List[RiskHotspot],
                                                    defense_coverage: DefenseCoverage,
                                                    blind_spots: List[BlindSpot],
                                                    critical_paths: List[CriticalPath]) -> List[str]:
        """Generate attack surface recommendations"""
        try:
            recommendations = []
            
            # Recommendations based on risk hotspots
            for hotspot in risk_hotspots:
                if hotspot.risk_score > 0.8:
                    recommendations.append(f"Immediate action required for {hotspot.name}")
                    recommendations.extend(hotspot.immediate_actions)
            
            # Recommendations based on defense gaps
            for gap in defense_coverage.defense_gaps:
                if gap.risk_level == "high":
                    recommendations.append(f"Address defense gap: {gap.description}")
                    recommendations.extend(gap.recommended_controls)
            
            # Recommendations based on blind spots
            for blind_spot in blind_spots:
                if blind_spot.priority == "high":
                    recommendations.append(f"Implement monitoring for {blind_spot.description}")
                    recommendations.extend(blind_spot.monitoring_recommendations)
            
            # Recommendations based on critical paths
            for path in critical_paths:
                if path.criticality_score > 0.1:
                    recommendations.append(f"Mitigate critical path: {path.name}")
                    recommendations.extend(path.prevention_controls)
            
            # General recommendations
            recommendations.extend([
                "Implement continuous attack surface monitoring",
                "Regular security assessments and penetration testing",
                "Develop incident response procedures",
                "Establish security awareness training"
            ])
            
            return recommendations
            
        except Exception as e:
            logger.error(f"Error generating attack surface recommendations: {e}")
            return [
                "Conduct comprehensive security assessment",
                "Implement basic security controls",
                "Establish monitoring and incident response"
            ]
    
    async def _calculate_attack_surface_confidence(self, vulnerability_data: List[Dict],
                                                network_data: List[Dict],
                                                asset_data: List[Dict]) -> float:
        """Calculate confidence score for attack surface analysis"""
        try:
            base_confidence = 0.6
            
            # Adjust based on data availability
            if vulnerability_data and len(vulnerability_data) > 10:
                base_confidence += 0.1
            
            if network_data and len(network_data) > 5:
                base_confidence += 0.1
            
            if asset_data and len(asset_data) > 5:
                base_confidence += 0.1
            
            return min(base_confidence, 1.0)
            
        except Exception as e:
            logger.error(f"Error calculating attack surface confidence: {e}")
            return 0.5
    
    def _get_fallback_attack_surface_map(self, organization_id: str) -> AttackSurfaceMap:
        """Get fallback attack surface map when analysis fails"""
        return AttackSurfaceMap(
            organization_id=organization_id,
            map_id=f"fallback_{organization_id}_{int(datetime.now().timestamp())}",
            external_exposure=ExternalExposure(
                public_assets=[],
                exposed_services=[],
                open_ports=[],
                web_applications=[],
                api_endpoints=[],
                cloud_resources=[],
                dns_records=[],
                ssl_certificates=[],
                exposure_score=0.0,
                risk_level="unknown"
            ),
            internal_vulnerabilities=InternalVulnerabilities(
                network_segments=[],
                internal_services=[],
                database_exposures=[],
                file_shares=[],
                internal_apis=[],
                privilege_escalation_paths=[],
                lateral_movement_paths=[],
                vulnerability_score=0.0,
                risk_level="unknown"
            ),
            supply_chain_exposure=SupplyChainExposure(
                dependencies=[],
                third_party_integrations=[],
                vendor_connections=[],
                license_risks=[],
                malicious_packages=[],
                transitive_vulnerabilities=[],
                supply_chain_score=0.0,
                risk_level="unknown"
            ),
            attack_vectors=[],
            critical_paths=[],
            risk_hotspots=[],
            defense_coverage=DefenseCoverage(
                coverage_percentage=0.0,
                protected_assets=0,
                unprotected_assets=0,
                defense_gaps=[],
                security_controls=[],
                monitoring_coverage=0.0,
                incident_response_readiness=0.0,
                overall_defense_score=0.0
            ),
            blind_spots=[],
            recommendations=[
                "Conduct comprehensive security assessment",
                "Implement basic security controls",
                "Establish monitoring and incident response"
            ],
            generated_at=datetime.utcnow().isoformat(),
            confidence_score=0.3
        )

# Global attack surface mapper instance
attack_surface_mapper = AttackSurfaceMapper()

