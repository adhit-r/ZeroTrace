"""
Security DNA Analysis Service - Pattern recognition and vulnerability trend analysis
"""

import asyncio
import json
import logging
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass
from sklearn.cluster import KMeans
from sklearn.preprocessing import StandardScaler
from sklearn.decomposition import PCA
from sklearn.metrics import silhouette_score
import matplotlib.pyplot as plt
import seaborn as sns

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class SecurityDNAProfile:
    """Represents a security DNA profile for an organization"""
    organization_id: str
    profile_id: str
    dna_signature: Dict[str, Any]
    vulnerability_patterns: List[Dict[str, Any]]
    remediation_velocity: Dict[str, Any]
    risk_acceptance_patterns: Dict[str, Any]
    technology_evolution: Dict[str, Any]
    security_maturity_score: float
    confidence_score: float
    generated_at: str
    recommendations: List[str]

@dataclass
class VulnerabilityPattern:
    """Represents a pattern in vulnerability data"""
    pattern_id: str
    pattern_type: str
    frequency: int
    severity_distribution: Dict[str, int]
    technology_focus: List[str]
    temporal_pattern: Dict[str, Any]
    risk_score: float
    description: str

@dataclass
class RemediationVelocity:
    """Represents remediation velocity metrics"""
    average_time_to_patch: float
    patch_velocity_trend: str
    critical_patch_time: float
    high_patch_time: float
    medium_patch_time: float
    low_patch_time: float
    velocity_score: float

class SecurityDNAAnalyzer:
    """Service for analyzing security DNA patterns and trends"""
    
    def __init__(self):
        self.cache = {}
        self.cache_ttl = 3600  # 1 hour cache TTL
        
    async def analyze_security_dna(self, organization_id: str, vulnerability_data: List[Dict], 
                                 scan_history: List[Dict] = None, organization_context: Dict = None) -> SecurityDNAProfile:
        """
        Analyze security DNA for an organization
        
        Args:
            organization_id: Organization identifier
            vulnerability_data: List of vulnerability records
            scan_history: Historical scan data
            organization_context: Organization profile data
        
        Returns:
            SecurityDNAProfile with comprehensive analysis
        """
        try:
            logger.info(f"Starting security DNA analysis for organization {organization_id}")
            
            # Analyze vulnerability patterns
            vulnerability_patterns = await self._analyze_vulnerability_patterns(vulnerability_data)
            
            # Calculate remediation velocity
            remediation_velocity = await self._calculate_remediation_velocity(vulnerability_data, scan_history)
            
            # Analyze risk acceptance patterns
            risk_acceptance = await self._analyze_risk_acceptance_patterns(vulnerability_data, organization_context)
            
            # Analyze technology evolution
            technology_evolution = await self._analyze_technology_evolution(vulnerability_data, scan_history)
            
            # Generate DNA signature
            dna_signature = await self._generate_dna_signature(
                vulnerability_patterns, remediation_velocity, risk_acceptance, technology_evolution
            )
            
            # Calculate security maturity score
            maturity_score = await self._calculate_security_maturity_score(
                vulnerability_patterns, remediation_velocity, risk_acceptance, technology_evolution
            )
            
            # Generate recommendations
            recommendations = await self._generate_dna_recommendations(
                vulnerability_patterns, remediation_velocity, risk_acceptance, technology_evolution
            )
            
            # Create security DNA profile
            profile = SecurityDNAProfile(
                organization_id=organization_id,
                profile_id=f"dna_{organization_id}_{int(datetime.now().timestamp())}",
                dna_signature=dna_signature,
                vulnerability_patterns=vulnerability_patterns,
                remediation_velocity=remediation_velocity,
                risk_acceptance_patterns=risk_acceptance,
                technology_evolution=technology_evolution,
                security_maturity_score=maturity_score,
                confidence_score=await self._calculate_confidence_score(vulnerability_data),
                generated_at=datetime.utcnow().isoformat(),
                recommendations=recommendations
            )
            
            logger.info(f"Completed security DNA analysis for organization {organization_id}")
            return profile
            
        except Exception as e:
            logger.error(f"Error in security DNA analysis: {e}")
            return self._get_fallback_profile(organization_id)
    
    async def _analyze_vulnerability_patterns(self, vulnerability_data: List[Dict]) -> List[Dict[str, Any]]:
        """Analyze patterns in vulnerability data"""
        patterns = []
        
        try:
            # Convert to DataFrame for analysis
            df = pd.DataFrame(vulnerability_data)
            
            if df.empty:
                return patterns
            
            # Pattern 1: Severity distribution patterns
            severity_pattern = self._analyze_severity_patterns(df)
            if severity_pattern:
                patterns.append(severity_pattern)
            
            # Pattern 2: Technology clustering patterns
            tech_pattern = self._analyze_technology_patterns(df)
            if tech_pattern:
                patterns.append(tech_pattern)
            
            # Pattern 3: Temporal patterns
            temporal_pattern = self._analyze_temporal_patterns(df)
            if temporal_pattern:
                patterns.append(temporal_pattern)
            
            # Pattern 4: CVSS score patterns
            cvss_pattern = self._analyze_cvss_patterns(df)
            if cvss_pattern:
                patterns.append(cvss_pattern)
            
            # Pattern 5: Package vulnerability patterns
            package_pattern = self._analyze_package_patterns(df)
            if package_pattern:
                patterns.append(package_pattern)
            
        except Exception as e:
            logger.error(f"Error analyzing vulnerability patterns: {e}")
        
        return patterns
    
    def _analyze_severity_patterns(self, df: pd.DataFrame) -> Optional[Dict[str, Any]]:
        """Analyze severity distribution patterns"""
        if 'severity' not in df.columns:
            return None
        
        severity_counts = df['severity'].value_counts()
        total_vulns = len(df)
        
        # Calculate severity ratios
        severity_ratios = {
            'critical_ratio': severity_counts.get('CRITICAL', 0) / total_vulns,
            'high_ratio': severity_counts.get('HIGH', 0) / total_vulns,
            'medium_ratio': severity_counts.get('MEDIUM', 0) / total_vulns,
            'low_ratio': severity_counts.get('LOW', 0) / total_vulns
        }
        
        # Identify patterns
        patterns = []
        if severity_ratios['critical_ratio'] > 0.2:
            patterns.append("High critical vulnerability ratio")
        if severity_ratios['high_ratio'] > 0.4:
            patterns.append("High severity vulnerability concentration")
        if severity_ratios['low_ratio'] > 0.6:
            patterns.append("Low severity vulnerability dominance")
        
        return {
            'pattern_type': 'severity_distribution',
            'severity_ratios': severity_ratios,
            'patterns': patterns,
            'risk_level': 'high' if severity_ratios['critical_ratio'] > 0.1 else 'medium',
            'description': f"Severity distribution shows {', '.join(patterns) if patterns else 'balanced distribution'}"
        }
    
    def _analyze_technology_patterns(self, df: pd.DataFrame) -> Optional[Dict[str, Any]]:
        """Analyze technology clustering patterns"""
        if 'package_name' not in df.columns:
            return None
        
        # Extract technology categories from package names
        tech_categories = []
        for package in df['package_name'].dropna():
            if 'web' in str(package).lower() or 'http' in str(package).lower():
                tech_categories.append('Web')
            elif 'database' in str(package).lower() or 'db' in str(package).lower():
                tech_categories.append('Database')
            elif 'auth' in str(package).lower() or 'login' in str(package).lower():
                tech_categories.append('Authentication')
            elif 'api' in str(package).lower():
                tech_categories.append('API')
            else:
                tech_categories.append('Other')
        
        tech_counts = pd.Series(tech_categories).value_counts()
        total_tech = len(tech_categories)
        
        # Calculate technology ratios
        tech_ratios = {tech: count / total_tech for tech, count in tech_counts.items()}
        
        # Identify dominant technologies
        dominant_tech = tech_counts.index[0] if len(tech_counts) > 0 else 'Unknown'
        dominant_ratio = tech_ratios.get(dominant_tech, 0)
        
        patterns = []
        if dominant_ratio > 0.5:
            patterns.append(f"Strong focus on {dominant_tech} technologies")
        if len(tech_counts) > 5:
            patterns.append("High technology diversity")
        if 'Web' in tech_ratios and tech_ratios['Web'] > 0.3:
            patterns.append("Web application security focus")
        
        return {
            'pattern_type': 'technology_clustering',
            'technology_ratios': tech_ratios,
            'dominant_technology': dominant_tech,
            'patterns': patterns,
            'risk_level': 'high' if dominant_ratio > 0.7 else 'medium',
            'description': f"Technology patterns show {', '.join(patterns) if patterns else 'balanced technology usage'}"
        }
    
    def _analyze_temporal_patterns(self, df: pd.DataFrame) -> Optional[Dict[str, Any]]:
        """Analyze temporal patterns in vulnerabilities"""
        if 'created_at' not in df.columns:
            return None
        
        try:
            # Convert to datetime
            df['created_at'] = pd.to_datetime(df['created_at'])
            
            # Group by time periods
            df['month'] = df['created_at'].dt.to_period('M')
            monthly_counts = df.groupby('month').size()
            
            # Calculate trends
            if len(monthly_counts) > 1:
                trend = 'increasing' if monthly_counts.iloc[-1] > monthly_counts.iloc[0] else 'decreasing'
                volatility = monthly_counts.std() / monthly_counts.mean() if monthly_counts.mean() > 0 else 0
            else:
                trend = 'stable'
                volatility = 0
            
            # Identify patterns
            patterns = []
            if volatility > 0.5:
                patterns.append("High temporal volatility in vulnerability discovery")
            if trend == 'increasing':
                patterns.append("Increasing vulnerability discovery trend")
            elif trend == 'decreasing':
                patterns.append("Decreasing vulnerability discovery trend")
            
            return {
                'pattern_type': 'temporal',
                'trend': trend,
                'volatility': volatility,
                'monthly_counts': monthly_counts.to_dict(),
                'patterns': patterns,
                'risk_level': 'high' if volatility > 0.7 else 'medium',
                'description': f"Temporal patterns show {', '.join(patterns) if patterns else 'stable vulnerability discovery'}"
            }
            
        except Exception as e:
            logger.error(f"Error analyzing temporal patterns: {e}")
            return None
    
    def _analyze_cvss_patterns(self, df: pd.DataFrame) -> Optional[Dict[str, Any]]:
        """Analyze CVSS score patterns"""
        if 'cvss_score' not in df.columns:
            return None
        
        cvss_scores = df['cvss_score'].dropna()
        if len(cvss_scores) == 0:
            return None
        
        # Calculate CVSS statistics
        mean_cvss = cvss_scores.mean()
        std_cvss = cvss_scores.std()
        high_cvss_ratio = (cvss_scores >= 7.0).sum() / len(cvss_scores)
        critical_cvss_ratio = (cvss_scores >= 9.0).sum() / len(cvss_scores)
        
        # Identify patterns
        patterns = []
        if mean_cvss > 7.0:
            patterns.append("High average CVSS scores")
        if high_cvss_ratio > 0.3:
            patterns.append("High concentration of high-severity vulnerabilities")
        if critical_cvss_ratio > 0.1:
            patterns.append("Significant critical vulnerability presence")
        if std_cvss > 2.0:
            patterns.append("High CVSS score variability")
        
        return {
            'pattern_type': 'cvss_distribution',
            'mean_cvss': mean_cvss,
            'std_cvss': std_cvss,
            'high_cvss_ratio': high_cvss_ratio,
            'critical_cvss_ratio': critical_cvss_ratio,
            'patterns': patterns,
            'risk_level': 'high' if mean_cvss > 7.0 or critical_cvss_ratio > 0.1 else 'medium',
            'description': f"CVSS patterns show {', '.join(patterns) if patterns else 'moderate severity levels'}"
        }
    
    def _analyze_package_patterns(self, df: pd.DataFrame) -> Optional[Dict[str, Any]]:
        """Analyze package vulnerability patterns"""
        if 'package_name' not in df.columns:
            return None
        
        package_counts = df['package_name'].value_counts()
        total_packages = len(package_counts)
        
        # Calculate package vulnerability ratios
        high_vuln_packages = package_counts[package_counts > 1]
        package_concentration = len(high_vuln_packages) / total_packages if total_packages > 0 else 0
        
        # Identify patterns
        patterns = []
        if package_concentration > 0.3:
            patterns.append("High vulnerability concentration in specific packages")
        if len(high_vuln_packages) > 0:
            most_vulnerable = package_counts.index[0]
            patterns.append(f"Package '{most_vulnerable}' has highest vulnerability count")
        if total_packages > 50:
            patterns.append("High package diversity")
        
        return {
            'pattern_type': 'package_vulnerability',
            'total_packages': total_packages,
            'package_concentration': package_concentration,
            'most_vulnerable_packages': high_vuln_packages.head(5).to_dict(),
            'patterns': patterns,
            'risk_level': 'high' if package_concentration > 0.5 else 'medium',
            'description': f"Package patterns show {', '.join(patterns) if patterns else 'distributed vulnerability across packages'}"
        }
    
    async def _calculate_remediation_velocity(self, vulnerability_data: List[Dict], scan_history: List[Dict] = None) -> Dict[str, Any]:
        """Calculate remediation velocity metrics"""
        try:
            # Mock remediation velocity calculation
            # In a real implementation, this would analyze actual remediation times
            
            avg_critical_time = 2.5  # days
            avg_high_time = 7.0      # days
            avg_medium_time = 14.0    # days
            avg_low_time = 30.0      # days
            
            # Calculate overall average
            avg_time = (avg_critical_time + avg_high_time + avg_medium_time + avg_low_time) / 4
            
            # Determine trend (would be calculated from historical data)
            trend = "improving"  # Mock trend
            
            # Calculate velocity score (0-1, higher is better)
            velocity_score = max(0, 1 - (avg_time / 30))  # Normalize to 30 days
            
            return {
                'average_time_to_patch': avg_time,
                'patch_velocity_trend': trend,
                'critical_patch_time': avg_critical_time,
                'high_patch_time': avg_high_time,
                'medium_patch_time': avg_medium_time,
                'low_patch_time': avg_low_time,
                'velocity_score': velocity_score
            }
            
        except Exception as e:
            logger.error(f"Error calculating remediation velocity: {e}")
            return {
                'average_time_to_patch': 15.0,
                'patch_velocity_trend': 'stable',
                'critical_patch_time': 5.0,
                'high_patch_time': 10.0,
                'medium_patch_time': 20.0,
                'low_patch_time': 30.0,
                'velocity_score': 0.5
            }
    
    async def _analyze_risk_acceptance_patterns(self, vulnerability_data: List[Dict], organization_context: Dict = None) -> Dict[str, Any]:
        """Analyze risk acceptance patterns"""
        try:
            # Analyze vulnerability age and acceptance patterns
            current_time = datetime.now()
            accepted_vulns = 0
            total_vulns = len(vulnerability_data)
            
            # Mock analysis of risk acceptance
            for vuln in vulnerability_data:
                # Simulate risk acceptance based on age and severity
                vuln_age = 30  # Mock age in days
                severity = vuln.get('severity', 'MEDIUM')
                
                # Risk acceptance logic
                if vuln_age > 90 and severity in ['LOW', 'MEDIUM']:
                    accepted_vulns += 1
                elif vuln_age > 30 and severity == 'HIGH':
                    accepted_vulns += 1
            
            acceptance_ratio = accepted_vulns / total_vulns if total_vulns > 0 else 0
            
            # Determine risk tolerance level
            if acceptance_ratio > 0.3:
                risk_tolerance = 'high'
            elif acceptance_ratio > 0.1:
                risk_tolerance = 'medium'
            else:
                risk_tolerance = 'low'
            
            return {
                'risk_tolerance_level': risk_tolerance,
                'acceptance_ratio': acceptance_ratio,
                'accepted_vulnerabilities': accepted_vulns,
                'total_vulnerabilities': total_vulns,
                'risk_acceptance_trend': 'stable',  # Would be calculated from historical data
                'compliance_impact': 'medium' if risk_tolerance == 'high' else 'low'
            }
            
        except Exception as e:
            logger.error(f"Error analyzing risk acceptance patterns: {e}")
            return {
                'risk_tolerance_level': 'medium',
                'acceptance_ratio': 0.2,
                'accepted_vulnerabilities': 0,
                'total_vulnerabilities': 0,
                'risk_acceptance_trend': 'stable',
                'compliance_impact': 'medium'
            }
    
    async def _analyze_technology_evolution(self, vulnerability_data: List[Dict], scan_history: List[Dict] = None) -> Dict[str, Any]:
        """Analyze technology evolution patterns"""
        try:
            # Extract technology trends from vulnerability data
            tech_evolution = {
                'emerging_technologies': [],
                'legacy_technologies': [],
                'technology_adoption_rate': 0.5,
                'security_maturity_by_tech': {},
                'technology_risk_profile': 'medium'
            }
            
            # Mock technology evolution analysis
            if vulnerability_data:
                # Analyze package names for technology trends
                packages = [vuln.get('package_name', '') for vuln in vulnerability_data]
                
                # Identify emerging vs legacy technologies
                web_tech_count = sum(1 for pkg in packages if 'web' in str(pkg).lower())
                db_tech_count = sum(1 for pkg in packages if 'database' in str(pkg).lower() or 'db' in str(pkg).lower())
                
                if web_tech_count > len(packages) * 0.3:
                    tech_evolution['emerging_technologies'].append('Web Technologies')
                if db_tech_count > len(packages) * 0.2:
                    tech_evolution['legacy_technologies'].append('Database Systems')
                
                # Calculate adoption rate (mock)
                tech_evolution['technology_adoption_rate'] = min(1.0, len(set(packages)) / 10)
                
                # Security maturity by technology
                tech_evolution['security_maturity_by_tech'] = {
                    'Web': 0.7,
                    'Database': 0.8,
                    'API': 0.6,
                    'Infrastructure': 0.9
                }
                
                # Overall technology risk profile
                if tech_evolution['technology_adoption_rate'] > 0.7:
                    tech_evolution['technology_risk_profile'] = 'high'
                elif tech_evolution['technology_adoption_rate'] > 0.4:
                    tech_evolution['technology_risk_profile'] = 'medium'
                else:
                    tech_evolution['technology_risk_profile'] = 'low'
            
            return tech_evolution
            
        except Exception as e:
            logger.error(f"Error analyzing technology evolution: {e}")
            return {
                'emerging_technologies': [],
                'legacy_technologies': [],
                'technology_adoption_rate': 0.5,
                'security_maturity_by_tech': {},
                'technology_risk_profile': 'medium'
            }
    
    async def _generate_dna_signature(self, vulnerability_patterns: List[Dict], remediation_velocity: Dict, 
                                   risk_acceptance: Dict, technology_evolution: Dict) -> Dict[str, Any]:
        """Generate a unique DNA signature for the organization"""
        try:
            # Create DNA signature based on all analysis components
            signature = {
                'vulnerability_dna': {
                    'pattern_count': len(vulnerability_patterns),
                    'dominant_patterns': [p.get('pattern_type') for p in vulnerability_patterns[:3]],
                    'risk_level': self._calculate_overall_risk_level(vulnerability_patterns)
                },
                'remediation_dna': {
                    'velocity_score': remediation_velocity.get('velocity_score', 0.5),
                    'trend': remediation_velocity.get('patch_velocity_trend', 'stable'),
                    'efficiency_level': 'high' if remediation_velocity.get('velocity_score', 0.5) > 0.7 else 'medium'
                },
                'risk_acceptance_dna': {
                    'tolerance_level': risk_acceptance.get('risk_tolerance_level', 'medium'),
                    'acceptance_ratio': risk_acceptance.get('acceptance_ratio', 0.2),
                    'compliance_impact': risk_acceptance.get('compliance_impact', 'medium')
                },
                'technology_dna': {
                    'adoption_rate': technology_evolution.get('technology_adoption_rate', 0.5),
                    'risk_profile': technology_evolution.get('technology_risk_profile', 'medium'),
                    'maturity_distribution': technology_evolution.get('security_maturity_by_tech', {})
                },
                'overall_signature': self._generate_overall_signature(
                    vulnerability_patterns, remediation_velocity, risk_acceptance, technology_evolution
                )
            }
            
            return signature
            
        except Exception as e:
            logger.error(f"Error generating DNA signature: {e}")
            return {'error': str(e)}
    
    def _calculate_overall_risk_level(self, vulnerability_patterns: List[Dict]) -> str:
        """Calculate overall risk level from vulnerability patterns"""
        high_risk_count = sum(1 for pattern in vulnerability_patterns if pattern.get('risk_level') == 'high')
        total_patterns = len(vulnerability_patterns)
        
        if total_patterns == 0:
            return 'medium'
        
        high_risk_ratio = high_risk_count / total_patterns
        
        if high_risk_ratio > 0.5:
            return 'high'
        elif high_risk_ratio > 0.2:
            return 'medium'
        else:
            return 'low'
    
    def _generate_overall_signature(self, vulnerability_patterns: List[Dict], remediation_velocity: Dict,
                                  risk_acceptance: Dict, technology_evolution: Dict) -> str:
        """Generate a human-readable overall signature"""
        signatures = []
        
        # Vulnerability signature
        vuln_risk = self._calculate_overall_risk_level(vulnerability_patterns)
        signatures.append(f"Vulnerability risk: {vuln_risk}")
        
        # Remediation signature
        velocity_score = remediation_velocity.get('velocity_score', 0.5)
        if velocity_score > 0.7:
            signatures.append("Fast remediation")
        elif velocity_score > 0.4:
            signatures.append("Moderate remediation")
        else:
            signatures.append("Slow remediation")
        
        # Risk acceptance signature
        tolerance = risk_acceptance.get('risk_tolerance_level', 'medium')
        signatures.append(f"Risk tolerance: {tolerance}")
        
        # Technology signature
        adoption_rate = technology_evolution.get('technology_adoption_rate', 0.5)
        if adoption_rate > 0.7:
            signatures.append("High technology adoption")
        elif adoption_rate > 0.4:
            signatures.append("Moderate technology adoption")
        else:
            signatures.append("Conservative technology adoption")
        
        return "; ".join(signatures)
    
    async def _calculate_security_maturity_score(self, vulnerability_patterns: List[Dict], remediation_velocity: Dict,
                                               risk_acceptance: Dict, technology_evolution: Dict) -> float:
        """Calculate overall security maturity score"""
        try:
            scores = []
            
            # Vulnerability management maturity (0-1)
            vuln_score = 1.0 - (len([p for p in vulnerability_patterns if p.get('risk_level') == 'high']) / max(len(vulnerability_patterns), 1))
            scores.append(vuln_score)
            
            # Remediation maturity (0-1)
            remediation_score = remediation_velocity.get('velocity_score', 0.5)
            scores.append(remediation_score)
            
            # Risk management maturity (0-1)
            risk_score = 1.0 - risk_acceptance.get('acceptance_ratio', 0.2)
            scores.append(risk_score)
            
            # Technology maturity (0-1)
            tech_score = technology_evolution.get('security_maturity_by_tech', {})
            if tech_score:
                avg_tech_maturity = sum(tech_score.values()) / len(tech_score)
            else:
                avg_tech_maturity = 0.5
            scores.append(avg_tech_maturity)
            
            # Calculate weighted average
            weights = [0.3, 0.3, 0.2, 0.2]  # Vulnerability, Remediation, Risk, Technology
            maturity_score = sum(score * weight for score, weight in zip(scores, weights))
            
            return min(max(maturity_score, 0.0), 1.0)
            
        except Exception as e:
            logger.error(f"Error calculating security maturity score: {e}")
            return 0.5
    
    async def _generate_dna_recommendations(self, vulnerability_patterns: List[Dict], remediation_velocity: Dict,
                                         risk_acceptance: Dict, technology_evolution: Dict) -> List[str]:
        """Generate recommendations based on DNA analysis"""
        recommendations = []
        
        try:
            # Vulnerability pattern recommendations
            high_risk_patterns = [p for p in vulnerability_patterns if p.get('risk_level') == 'high']
            if high_risk_patterns:
                recommendations.append("Address high-risk vulnerability patterns immediately")
                recommendations.append("Implement additional monitoring for high-risk areas")
            
            # Remediation velocity recommendations
            velocity_score = remediation_velocity.get('velocity_score', 0.5)
            if velocity_score < 0.5:
                recommendations.append("Improve patch management processes")
                recommendations.append("Implement automated remediation where possible")
            
            # Risk acceptance recommendations
            acceptance_ratio = risk_acceptance.get('acceptance_ratio', 0.2)
            if acceptance_ratio > 0.3:
                recommendations.append("Review and tighten risk acceptance criteria")
                recommendations.append("Implement stricter vulnerability management policies")
            
            # Technology evolution recommendations
            adoption_rate = technology_evolution.get('technology_adoption_rate', 0.5)
            if adoption_rate > 0.7:
                recommendations.append("Implement security-first approach for new technologies")
                recommendations.append("Establish technology security review processes")
            
            # General recommendations
            recommendations.extend([
                "Regular security DNA analysis to track improvements",
                "Benchmark against industry standards",
                "Continuous improvement of security processes"
            ])
            
        except Exception as e:
            logger.error(f"Error generating DNA recommendations: {e}")
            recommendations = ["Conduct regular security assessments", "Implement continuous improvement processes"]
        
        return recommendations
    
    async def _calculate_confidence_score(self, vulnerability_data: List[Dict]) -> float:
        """Calculate confidence score for the DNA analysis"""
        try:
            base_confidence = 0.7
            
            # Adjust based on data volume
            if len(vulnerability_data) > 100:
                base_confidence += 0.1
            elif len(vulnerability_data) > 50:
                base_confidence += 0.05
            
            # Adjust based on data quality
            complete_records = sum(1 for vuln in vulnerability_data 
                                 if vuln.get('severity') and vuln.get('package_name'))
            if len(vulnerability_data) > 0:
                quality_ratio = complete_records / len(vulnerability_data)
                base_confidence += quality_ratio * 0.2
            
            return min(base_confidence, 1.0)
            
        except Exception as e:
            logger.error(f"Error calculating confidence score: {e}")
            return 0.5
    
    def _get_fallback_profile(self, organization_id: str) -> SecurityDNAProfile:
        """Get fallback profile when analysis fails"""
        return SecurityDNAProfile(
            organization_id=organization_id,
            profile_id=f"fallback_{organization_id}_{int(datetime.now().timestamp())}",
            dna_signature={'error': 'Analysis failed'},
            vulnerability_patterns=[],
            remediation_velocity={
                'average_time_to_patch': 15.0,
                'patch_velocity_trend': 'stable',
                'velocity_score': 0.5
            },
            risk_acceptance_patterns={
                'risk_tolerance_level': 'medium',
                'acceptance_ratio': 0.2
            },
            technology_evolution={
                'technology_adoption_rate': 0.5,
                'technology_risk_profile': 'medium'
            },
            security_maturity_score=0.5,
            confidence_score=0.3,
            generated_at=datetime.utcnow().isoformat(),
            recommendations=["Conduct comprehensive security assessment", "Implement regular monitoring"]
        )

# Global security DNA analyzer instance
security_dna_analyzer = SecurityDNAAnalyzer()
