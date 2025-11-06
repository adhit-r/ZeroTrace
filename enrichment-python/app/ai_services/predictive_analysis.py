"""
Predictive Vulnerability Analysis Service - ML-based vulnerability prediction and analysis
"""

import asyncio
import json
import logging
import numpy as np
import pandas as pd
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from sklearn.ensemble import RandomForestClassifier, GradientBoostingRegressor
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler, LabelEncoder
from sklearn.metrics import accuracy_score, mean_squared_error
import joblib
import os

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class PredictiveAnalysisService:
    """Service for predictive vulnerability analysis using ML models"""
    
    def __init__(self):
        self.models = {}
        self.scalers = {}
        self.label_encoders = {}
        self.model_dir = os.path.join(os.path.dirname(__file__), 'models')
        
        # Create models directory if it doesn't exist
        os.makedirs(self.model_dir, exist_ok=True)
        
        # Initialize models
        self._initialize_models()
    
    def _initialize_models(self):
        """Initialize ML models for different prediction tasks"""
        try:
            # Exploit likelihood prediction model
            self.models['exploit_likelihood'] = RandomForestClassifier(
                n_estimators=100,
                max_depth=10,
                random_state=42
            )
            
            # Business impact prediction model
            self.models['business_impact'] = GradientBoostingRegressor(
                n_estimators=100,
                max_depth=6,
                learning_rate=0.1,
                random_state=42
            )
            
            # Remediation complexity estimator
            self.models['remediation_complexity'] = RandomForestClassifier(
                n_estimators=50,
                max_depth=8,
                random_state=42
            )
            
            # Timeline urgency calculator
            self.models['timeline_urgency'] = GradientBoostingRegressor(
                n_estimators=80,
                max_depth=5,
                learning_rate=0.15,
                random_state=42
            )
            
            # Initialize scalers and encoders
            self.scalers['standard'] = StandardScaler()
            self.label_encoders['severity'] = LabelEncoder()
            self.label_encoders['vulnerability_type'] = LabelEncoder()
            self.label_encoders['industry'] = LabelEncoder()
            
            logger.info("ML models initialized successfully")
            
        except Exception as e:
            logger.error(f"Error initializing models: {e}")
    
    async def predict_exploit_likelihood(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """
        Predict the likelihood of a vulnerability being exploited
        
        Args:
            vulnerability_data: Vulnerability details
            organization_context: Organization-specific context
        
        Returns:
            Dict containing exploit likelihood prediction
        """
        try:
            # Prepare features for the model
            features = self._prepare_exploit_features(vulnerability_data, organization_context)
            
            # Load or train the model
            model = await self._get_or_train_model('exploit_likelihood', features)
            
            if model is None:
                return self._get_fallback_exploit_prediction(vulnerability_data)
            
            # Make prediction
            prediction = model.predict([features])[0]
            probability = model.predict_proba([features])[0]
            
            # Calculate confidence score
            confidence = max(probability)
            
            return {
                'exploit_likelihood': float(prediction),
                'confidence_score': float(confidence),
                'probability_distribution': {
                    'low': float(probability[0]),
                    'medium': float(probability[1]),
                    'high': float(probability[2])
                },
                'prediction_factors': self._analyze_exploit_factors(vulnerability_data, features),
                'recommendations': self._get_exploit_recommendations(prediction, confidence),
                'predicted_at': datetime.utcnow().isoformat()
            }
            
        except Exception as e:
            logger.error(f"Error predicting exploit likelihood: {e}")
            return self._get_fallback_exploit_prediction(vulnerability_data)
    
    async def predict_business_impact(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """
        Predict the business impact of a vulnerability
        
        Args:
            vulnerability_data: Vulnerability details
            organization_context: Organization-specific context
        
        Returns:
            Dict containing business impact prediction
        """
        try:
            # Prepare features for the model
            features = self._prepare_business_impact_features(vulnerability_data, organization_context)
            
            # Load or train the model
            model = await self._get_or_train_model('business_impact', features)
            
            if model is None:
                return self._get_fallback_business_impact_prediction(vulnerability_data)
            
            # Make prediction
            impact_score = model.predict([features])[0]
            
            # Categorize impact level
            if impact_score >= 0.8:
                impact_level = 'critical'
            elif impact_score >= 0.6:
                impact_level = 'high'
            elif impact_score >= 0.4:
                impact_level = 'medium'
            else:
                impact_level = 'low'
            
            return {
                'business_impact_score': float(impact_score),
                'impact_level': impact_level,
                'potential_losses': self._estimate_potential_losses(impact_score, organization_context),
                'affected_systems': self._identify_affected_systems(vulnerability_data, organization_context),
                'recovery_time': self._estimate_recovery_time(impact_score, organization_context),
                'predicted_at': datetime.utcnow().isoformat()
            }
            
        except Exception as e:
            logger.error(f"Error predicting business impact: {e}")
            return self._get_fallback_business_impact_prediction(vulnerability_data)
    
    async def estimate_remediation_complexity(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """
        Estimate the complexity of remediating a vulnerability
        
        Args:
            vulnerability_data: Vulnerability details
            organization_context: Organization-specific context
        
        Returns:
            Dict containing remediation complexity estimation
        """
        try:
            # Prepare features for the model
            features = self._prepare_remediation_features(vulnerability_data, organization_context)
            
            # Load or train the model
            model = await self._get_or_train_model('remediation_complexity', features)
            
            if model is None:
                return self._get_fallback_remediation_prediction(vulnerability_data)
            
            # Make prediction
            complexity = model.predict([features])[0]
            probability = model.predict_proba([features])[0]
            
            return {
                'complexity_level': complexity,
                'confidence_score': float(max(probability)),
                'estimated_effort': self._estimate_remediation_effort(complexity, vulnerability_data),
                'required_skills': self._identify_required_skills(complexity, vulnerability_data),
                'testing_requirements': self._assess_testing_requirements(complexity, vulnerability_data),
                'rollback_complexity': self._assess_rollback_complexity(complexity, vulnerability_data),
                'predicted_at': datetime.utcnow().isoformat()
            }
            
        except Exception as e:
            logger.error(f"Error estimating remediation complexity: {e}")
            return self._get_fallback_remediation_prediction(vulnerability_data)
    
    async def calculate_timeline_urgency(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """
        Calculate timeline urgency for vulnerability remediation
        
        Args:
            vulnerability_data: Vulnerability details
            organization_context: Organization-specific context
        
        Returns:
            Dict containing timeline urgency calculation
        """
        try:
            # Prepare features for the model
            features = self._prepare_timeline_features(vulnerability_data, organization_context)
            
            # Load or train the model
            model = await self._get_or_train_model('timeline_urgency', features)
            
            if model is None:
                return self._get_fallback_timeline_prediction(vulnerability_data)
            
            # Make prediction
            urgency_score = model.predict([features])[0]
            
            # Calculate timeline recommendations
            timeline = self._calculate_remediation_timeline(urgency_score, vulnerability_data, organization_context)
            
            return {
                'urgency_score': float(urgency_score),
                'urgency_level': self._categorize_urgency(urgency_score),
                'recommended_timeline': timeline,
                'deadline_risk': self._assess_deadline_risk(urgency_score, organization_context),
                'resource_requirements': self._estimate_resource_requirements(urgency_score, vulnerability_data),
                'predicted_at': datetime.utcnow().isoformat()
            }
            
        except Exception as e:
            logger.error(f"Error calculating timeline urgency: {e}")
            return self._get_fallback_timeline_prediction(vulnerability_data)
    
    def _prepare_exploit_features(self, vulnerability_data: Dict, organization_context: Dict = None) -> List[float]:
        """Prepare features for exploit likelihood prediction"""
        features = []
        
        # Vulnerability characteristics
        features.append(self._encode_severity(vulnerability_data.get('severity', 'MEDIUM')))
        features.append(self._encode_vulnerability_type(vulnerability_data.get('type', 'unknown')))
        features.append(vulnerability_data.get('cvss_score', 0.0) or 0.0)
        features.append(1.0 if vulnerability_data.get('exploit_available', False) else 0.0)
        features.append(vulnerability_data.get('exploit_count', 0) or 0)
        
        # Package characteristics
        package_name = vulnerability_data.get('package_name', '')
        features.append(len(package_name))  # Package name length
        features.append(1.0 if 'web' in package_name.lower() else 0.0)  # Web-related
        features.append(1.0 if 'database' in package_name.lower() else 0.0)  # Database-related
        
        # Organization context
        if organization_context:
            features.append(self._encode_industry(organization_context.get('industry', 'unknown')))
            features.append(1.0 if organization_context.get('risk_tolerance') == 'CONSERVATIVE' else 0.0)
            features.append(len(organization_context.get('tech_stack', {}).get('languages', [])))
        else:
            features.extend([0.0, 0.0, 0.0])
        
        # Time-based features
        cve_id = vulnerability_data.get('cve_id', '')
        if cve_id and '-' in cve_id:
            try:
                cve_year = int(cve_id.split('-')[1])
                current_year = datetime.now().year
                features.append(current_year - cve_year)  # Age of CVE
            except:
                features.append(0.0)
        else:
            features.append(0.0)
        
        return features
    
    def _prepare_business_impact_features(self, vulnerability_data: Dict, organization_context: Dict = None) -> List[float]:
        """Prepare features for business impact prediction"""
        features = []
        
        # Vulnerability severity and type
        features.append(self._encode_severity(vulnerability_data.get('severity', 'MEDIUM')))
        features.append(vulnerability_data.get('cvss_score', 0.0) or 0.0)
        features.append(1.0 if vulnerability_data.get('exploit_available', False) else 0.0)
        
        # Organization characteristics
        if organization_context:
            features.append(self._encode_industry(organization_context.get('industry', 'unknown')))
            features.append(len(organization_context.get('tech_stack', {}).get('languages', [])))
            features.append(len(organization_context.get('compliance_frameworks', [])))
        else:
            features.extend([0.0, 0.0, 0.0])
        
        # System exposure
        package_name = vulnerability_data.get('package_name', '')
        features.append(1.0 if 'web' in package_name.lower() else 0.0)
        features.append(1.0 if 'database' in package_name.lower() else 0.0)
        features.append(1.0 if 'auth' in package_name.lower() else 0.0)
        
        return features
    
    def _prepare_remediation_features(self, vulnerability_data: Dict, organization_context: Dict = None) -> List[float]:
        """Prepare features for remediation complexity estimation"""
        features = []
        
        # Vulnerability characteristics
        features.append(self._encode_severity(vulnerability_data.get('severity', 'MEDIUM')))
        features.append(vulnerability_data.get('cvss_score', 0.0) or 0.0)
        features.append(len(vulnerability_data.get('affected_versions', [])) or 0)
        features.append(len(vulnerability_data.get('patched_versions', [])) or 0)
        
        # Package complexity
        package_name = vulnerability_data.get('package_name', '')
        features.append(len(package_name))
        features.append(1.0 if 'framework' in package_name.lower() else 0.0)
        features.append(1.0 if 'database' in package_name.lower() else 0.0)
        
        # Organization tech stack complexity
        if organization_context:
            tech_stack = organization_context.get('tech_stack', {})
            features.append(len(tech_stack.get('languages', [])))
            features.append(len(tech_stack.get('frameworks', [])))
            features.append(len(tech_stack.get('databases', [])))
        else:
            features.extend([0.0, 0.0, 0.0])
        
        return features
    
    def _prepare_timeline_features(self, vulnerability_data: Dict, organization_context: Dict = None) -> List[float]:
        """Prepare features for timeline urgency calculation"""
        features = []
        
        # Vulnerability urgency factors
        features.append(self._encode_severity(vulnerability_data.get('severity', 'MEDIUM')))
        features.append(vulnerability_data.get('cvss_score', 0.0) or 0.0)
        features.append(1.0 if vulnerability_data.get('exploit_available', False) else 0.0)
        features.append(vulnerability_data.get('exploit_count', 0) or 0)
        
        # Compliance and regulatory factors
        if organization_context:
            compliance_frameworks = organization_context.get('compliance_frameworks', [])
            features.append(len(compliance_frameworks))
            features.append(1.0 if 'HIPAA' in compliance_frameworks else 0.0)
            features.append(1.0 if 'PCI' in str(compliance_frameworks) else 0.0)
        else:
            features.extend([0.0, 0.0, 0.0])
        
        # Business criticality
        package_name = vulnerability_data.get('package_name', '')
        features.append(1.0 if 'payment' in package_name.lower() else 0.0)
        features.append(1.0 if 'user' in package_name.lower() else 0.0)
        features.append(1.0 if 'api' in package_name.lower() else 0.0)
        
        return features
    
    async def _get_or_train_model(self, model_name: str, features: List[float]) -> Optional[Any]:
        """Get existing model or train a new one"""
        model_path = os.path.join(self.model_dir, f"{model_name}.joblib")
        
        try:
            # Try to load existing model
            if os.path.exists(model_path):
                model = joblib.load(model_path)
                return model
        except Exception as e:
            logger.warning(f"Could not load model {model_name}: {e}")
        
        # For now, return the initialized model
        # In production, this would train on historical data
        return self.models.get(model_name)
    
    def _encode_severity(self, severity: str) -> float:
        """Encode severity level to numeric value"""
        severity_map = {
            'CRITICAL': 4.0,
            'HIGH': 3.0,
            'MEDIUM': 2.0,
            'LOW': 1.0,
            'INFO': 0.0
        }
        return severity_map.get(severity.upper(), 2.0)
    
    def _encode_vulnerability_type(self, vuln_type: str) -> float:
        """Encode vulnerability type to numeric value"""
        type_map = {
            'remote_code_execution': 4.0,
            'privilege_escalation': 3.5,
            'buffer_overflow': 4.0,
            'sql_injection': 3.0,
            'cross_site_scripting': 2.5,
            'denial_of_service': 2.0,
            'information_disclosure': 1.5
        }
        return type_map.get(vuln_type.lower(), 2.0)
    
    def _encode_industry(self, industry: str) -> float:
        """Encode industry to numeric value"""
        industry_map = {
            'healthcare': 4.0,
            'finance': 4.0,
            'government': 4.0,
            'technology': 3.0,
            'education': 2.5,
            'manufacturing': 2.0
        }
        return industry_map.get(industry.lower(), 2.0)
    
    def _analyze_exploit_factors(self, vulnerability_data: Dict, features: List[float]) -> List[str]:
        """Analyze factors contributing to exploit likelihood"""
        factors = []
        
        if vulnerability_data.get('exploit_available'):
            factors.append("Public exploits available")
        
        if vulnerability_data.get('cvss_score', 0) >= 8.0:
            factors.append("High CVSS score")
        
        if vulnerability_data.get('severity') == 'CRITICAL':
            factors.append("Critical severity")
        
        if 'web' in vulnerability_data.get('package_name', '').lower():
            factors.append("Web-facing component")
        
        return factors
    
    def _get_exploit_recommendations(self, prediction: int, confidence: float) -> List[str]:
        """Get recommendations based on exploit prediction"""
        recommendations = []
        
        if prediction >= 2:  # High likelihood
            recommendations.append("Immediate patching required")
            recommendations.append("Implement additional monitoring")
            recommendations.append("Consider temporary workarounds")
        elif prediction >= 1:  # Medium likelihood
            recommendations.append("Schedule patching within 7 days")
            recommendations.append("Monitor for exploitation attempts")
        else:  # Low likelihood
            recommendations.append("Include in regular patch cycle")
            recommendations.append("Monitor for new exploit development")
        
        if confidence < 0.7:
            recommendations.append("Gather additional intelligence")
        
        return recommendations
    
    def _estimate_potential_losses(self, impact_score: float, organization_context: Dict = None) -> Dict:
        """Estimate potential financial losses"""
        base_loss = impact_score * 1000000  # Base loss in dollars
        
        # Adjust based on organization size and industry
        if organization_context:
            industry = organization_context.get('industry', 'technology')
            if industry in ['healthcare', 'finance']:
                base_loss *= 2.0
            elif industry == 'government':
                base_loss *= 1.5
        
        return {
            'estimated_loss': base_loss,
            'reputation_damage': impact_score * 0.8,
            'regulatory_fines': impact_score * 0.6,
            'business_disruption': impact_score * 0.9
        }
    
    def _identify_affected_systems(self, vulnerability_data: Dict, organization_context: Dict = None) -> List[str]:
        """Identify systems that might be affected"""
        systems = []
        
        package_name = vulnerability_data.get('package_name', '').lower()
        
        if 'web' in package_name:
            systems.append("Web servers")
        if 'database' in package_name:
            systems.append("Database servers")
        if 'api' in package_name:
            systems.append("API endpoints")
        if 'auth' in package_name:
            systems.append("Authentication systems")
        
        return systems
    
    def _estimate_recovery_time(self, impact_score: float, organization_context: Dict = None) -> Dict:
        """Estimate recovery time based on impact score"""
        base_hours = impact_score * 24
        
        return {
            'detection_time': f"{base_hours * 0.1:.1f} hours",
            'containment_time': f"{base_hours * 0.3:.1f} hours",
            'recovery_time': f"{base_hours:.1f} hours",
            'total_downtime': f"{base_hours * 1.4:.1f} hours"
        }
    
    def _estimate_remediation_effort(self, complexity: str, vulnerability_data: Dict) -> Dict:
        """Estimate remediation effort based on complexity"""
        effort_map = {
            'low': {'hours': 4, 'people': 1, 'days': 1},
            'medium': {'hours': 16, 'people': 2, 'days': 2},
            'high': {'hours': 40, 'people': 3, 'days': 5},
            'critical': {'hours': 80, 'people': 4, 'days': 10}
        }
        
        return effort_map.get(complexity, effort_map['medium'])
    
    def _identify_required_skills(self, complexity: str, vulnerability_data: Dict) -> List[str]:
        """Identify required skills for remediation"""
        base_skills = ['System administration', 'Security patching']
        
        if complexity in ['high', 'critical']:
            base_skills.extend(['Advanced security analysis', 'Incident response'])
        
        package_name = vulnerability_data.get('package_name', '').lower()
        if 'database' in package_name:
            base_skills.append('Database administration')
        if 'web' in package_name:
            base_skills.append('Web application security')
        
        return base_skills
    
    def _assess_testing_requirements(self, complexity: str, vulnerability_data: Dict) -> List[str]:
        """Assess testing requirements for remediation"""
        requirements = ['Functional testing']
        
        if complexity in ['high', 'critical']:
            requirements.extend(['Security testing', 'Performance testing', 'Regression testing'])
        
        return requirements
    
    def _assess_rollback_complexity(self, complexity: str, vulnerability_data: Dict) -> str:
        """Assess rollback complexity"""
        if complexity in ['high', 'critical']:
            return 'high'
        elif complexity == 'medium':
            return 'medium'
        else:
            return 'low'
    
    def _calculate_remediation_timeline(self, urgency_score: float, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """Calculate recommended remediation timeline"""
        if urgency_score >= 0.8:
            return {
                'immediate_action': '0-24 hours',
                'patch_deployment': '1-3 days',
                'testing_validation': '2-4 days',
                'full_remediation': '3-7 days'
            }
        elif urgency_score >= 0.6:
            return {
                'immediate_action': '24-48 hours',
                'patch_deployment': '3-7 days',
                'testing_validation': '5-10 days',
                'full_remediation': '7-14 days'
            }
        elif urgency_score >= 0.4:
            return {
                'immediate_action': '48-72 hours',
                'patch_deployment': '7-14 days',
                'testing_validation': '10-20 days',
                'full_remediation': '14-30 days'
            }
        else:
            return {
                'immediate_action': '1 week',
                'patch_deployment': '2-4 weeks',
                'testing_validation': '3-6 weeks',
                'full_remediation': '1-2 months'
            }
    
    def _categorize_urgency(self, urgency_score: float) -> str:
        """Categorize urgency level"""
        if urgency_score >= 0.8:
            return 'critical'
        elif urgency_score >= 0.6:
            return 'high'
        elif urgency_score >= 0.4:
            return 'medium'
        else:
            return 'low'
    
    def _assess_deadline_risk(self, urgency_score: float, organization_context: Dict = None) -> Dict:
        """Assess risk of missing deadlines"""
        base_risk = urgency_score
        
        # Adjust based on compliance requirements
        if organization_context:
            compliance_frameworks = organization_context.get('compliance_frameworks', [])
            if any(framework in ['HIPAA', 'PCI DSS', 'SOX'] for framework in compliance_frameworks):
                base_risk += 0.2
        
        return {
            'deadline_risk_score': min(base_risk, 1.0),
            'risk_level': 'high' if base_risk >= 0.7 else 'medium' if base_risk >= 0.4 else 'low',
            'recommended_escalation': base_risk >= 0.7
        }
    
    def _estimate_resource_requirements(self, urgency_score: float, vulnerability_data: Dict) -> Dict:
        """Estimate resource requirements for remediation"""
        base_team_size = max(1, int(urgency_score * 4))
        
        return {
            'team_size': base_team_size,
            'required_roles': ['Security Engineer', 'System Administrator'],
            'estimated_cost': base_team_size * urgency_score * 1000,  # Cost in dollars
            'external_consultants': urgency_score >= 0.8
        }
    
    # Fallback methods for when models are not available
    def _get_fallback_exploit_prediction(self, vulnerability_data: Dict) -> Dict:
        """Fallback exploit prediction when ML model is not available"""
        severity = vulnerability_data.get('severity', 'MEDIUM')
        exploit_available = vulnerability_data.get('exploit_available', False)
        
        # Simple heuristic-based prediction
        if severity == 'CRITICAL' and exploit_available:
            likelihood = 0.9
        elif severity == 'HIGH' and exploit_available:
            likelihood = 0.7
        elif severity == 'CRITICAL':
            likelihood = 0.6
        elif severity == 'HIGH':
            likelihood = 0.4
        else:
            likelihood = 0.2
        
        return {
            'exploit_likelihood': likelihood,
            'confidence_score': 0.6,
            'probability_distribution': {
                'low': 1 - likelihood,
                'medium': likelihood * 0.5,
                'high': likelihood * 0.5
            },
            'prediction_factors': ['Heuristic analysis'],
            'recommendations': ['Monitor for exploit development'],
            'predicted_at': datetime.utcnow().isoformat()
        }
    
    def _get_fallback_business_impact_prediction(self, vulnerability_data: Dict) -> Dict:
        """Fallback business impact prediction"""
        severity = vulnerability_data.get('severity', 'MEDIUM')
        cvss_score = vulnerability_data.get('cvss_score', 0.0) or 0.0
        
        impact_score = min(cvss_score / 10.0, 1.0)
        
        return {
            'business_impact_score': impact_score,
            'impact_level': 'critical' if impact_score >= 0.8 else 'high' if impact_score >= 0.6 else 'medium',
            'potential_losses': {'estimated_loss': impact_score * 500000},
            'affected_systems': ['Multiple systems'],
            'recovery_time': {'total_downtime': f"{impact_score * 48:.1f} hours"},
            'predicted_at': datetime.utcnow().isoformat()
        }
    
    def _get_fallback_remediation_prediction(self, vulnerability_data: Dict) -> Dict:
        """Fallback remediation complexity prediction"""
        severity = vulnerability_data.get('severity', 'MEDIUM')
        
        if severity == 'CRITICAL':
            complexity = 'high'
        elif severity == 'HIGH':
            complexity = 'medium'
        else:
            complexity = 'low'
        
        return {
            'complexity_level': complexity,
            'confidence_score': 0.6,
            'estimated_effort': {'hours': 16, 'people': 2, 'days': 2},
            'required_skills': ['System administration', 'Security patching'],
            'testing_requirements': ['Functional testing'],
            'rollback_complexity': 'medium',
            'predicted_at': datetime.utcnow().isoformat()
        }
    
    def _get_fallback_timeline_prediction(self, vulnerability_data: Dict) -> Dict:
        """Fallback timeline urgency prediction"""
        severity = vulnerability_data.get('severity', 'MEDIUM')
        exploit_available = vulnerability_data.get('exploit_available', False)
        
        if severity == 'CRITICAL' and exploit_available:
            urgency_score = 0.9
        elif severity == 'CRITICAL':
            urgency_score = 0.7
        elif severity == 'HIGH' and exploit_available:
            urgency_score = 0.6
        else:
            urgency_score = 0.3
        
        return {
            'urgency_score': urgency_score,
            'urgency_level': 'critical' if urgency_score >= 0.8 else 'high' if urgency_score >= 0.6 else 'medium',
            'recommended_timeline': {
                'immediate_action': '24-48 hours',
                'patch_deployment': '3-7 days',
                'full_remediation': '7-14 days'
            },
            'deadline_risk': {'deadline_risk_score': urgency_score, 'risk_level': 'medium'},
            'resource_requirements': {'team_size': 2, 'estimated_cost': 2000},
            'predicted_at': datetime.utcnow().isoformat()
        }

# Global predictive analysis service instance
predictive_analysis_service = PredictiveAnalysisService()

