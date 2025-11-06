"""
Integrated AI Service - Main service that coordinates all AI-powered features
"""

import asyncio
import json
import logging
import time
from datetime import datetime
from typing import Dict, List, Optional, Any

from .ai_service import ai_service
from .exploit_intelligence import exploit_intel_service
from .predictive_analysis import predictive_analysis_service
from .remediation_guidance import remediation_guidance_service, RemediationPlan

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class IntegratedAIService:
    """Main AI service that coordinates all AI-powered features"""
    
    def __init__(self):
        self.ai_service = ai_service
        self.exploit_intel_service = exploit_intel_service
        self.predictive_analysis_service = predictive_analysis_service
        self.remediation_guidance_service = remediation_guidance_service
        
        # Initialize remediation guidance with AI service
        self.remediation_guidance_service.ai_service = self.ai_service
    
    async def analyze_vulnerability_comprehensive(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """
        Perform comprehensive AI-powered vulnerability analysis
        
        Args:
            vulnerability_data: Vulnerability details
            organization_context: Organization-specific context
        
        Returns:
            Comprehensive analysis including exploit intelligence, predictions, and remediation
        """
        try:
            logger.info(f"Starting comprehensive analysis for vulnerability {vulnerability_data.get('id', 'unknown')}")
            
            # Run all analyses concurrently for better performance
            tasks = [
                self._get_exploit_intelligence(vulnerability_data),
                self._get_predictive_analysis(vulnerability_data, organization_context),
                self._get_remediation_plan(vulnerability_data, organization_context)
            ]
            
            results = await asyncio.gather(*tasks, return_exceptions=True)
            
            # Process results
            exploit_intel = results[0] if not isinstance(results[0], Exception) else {}
            predictions = results[1] if not isinstance(results[1], Exception) else {}
            remediation_plan = results[2] if not isinstance(results[2], Exception) else None
            
            # Combine all analysis results
            comprehensive_analysis = {
                'vulnerability_id': vulnerability_data.get('id'),
                'analysis_timestamp': datetime.utcnow().isoformat(),
                'exploit_intelligence': exploit_intel,
                'predictive_analysis': predictions,
                'remediation_plan': remediation_plan.__dict__ if remediation_plan else None,
                'overall_risk_score': self._calculate_overall_risk_score(exploit_intel, predictions),
                'priority_recommendation': self._generate_priority_recommendation(exploit_intel, predictions),
                'action_items': self._generate_action_items(exploit_intel, predictions, remediation_plan),
                'confidence_score': self._calculate_confidence_score(exploit_intel, predictions)
            }
            
            logger.info(f"Completed comprehensive analysis for vulnerability {vulnerability_data.get('id', 'unknown')}")
            return comprehensive_analysis
            
        except Exception as e:
            logger.error(f"Error in comprehensive vulnerability analysis: {e}")
            return self._get_fallback_analysis(vulnerability_data, organization_context)
    
    async def _get_exploit_intelligence(self, vulnerability_data: Dict) -> Dict:
        """Get exploit intelligence for vulnerability"""
        try:
            cve_id = vulnerability_data.get('cve_id', '')
            package_name = vulnerability_data.get('package_name', '')
            
            if cve_id:
                return await self.exploit_intel_service.get_exploit_intelligence(cve_id, package_name)
            else:
                return {'error': 'No CVE ID available for exploit intelligence'}
        except Exception as e:
            logger.error(f"Error getting exploit intelligence: {e}")
            return {'error': str(e)}
    
    async def _get_predictive_analysis(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """Get predictive analysis for vulnerability"""
        try:
            # Run all predictive analyses concurrently
            tasks = [
                self.predictive_analysis_service.predict_exploit_likelihood(vulnerability_data, organization_context),
                self.predictive_analysis_service.predict_business_impact(vulnerability_data, organization_context),
                self.predictive_analysis_service.estimate_remediation_complexity(vulnerability_data, organization_context),
                self.predictive_analysis_service.calculate_timeline_urgency(vulnerability_data, organization_context)
            ]
            
            results = await asyncio.gather(*tasks, return_exceptions=True)
            
            return {
                'exploit_likelihood': results[0] if not isinstance(results[0], Exception) else {},
                'business_impact': results[1] if not isinstance(results[1], Exception) else {},
                'remediation_complexity': results[2] if not isinstance(results[2], Exception) else {},
                'timeline_urgency': results[3] if not isinstance(results[3], Exception) else {}
            }
        except Exception as e:
            logger.error(f"Error getting predictive analysis: {e}")
            return {'error': str(e)}
    
    async def _get_remediation_plan(self, vulnerability_data: Dict, organization_context: Dict = None) -> Optional[RemediationPlan]:
        """Get remediation plan for vulnerability"""
        try:
            return await self.remediation_guidance_service.generate_remediation_plan(vulnerability_data, organization_context)
        except Exception as e:
            logger.error(f"Error getting remediation plan: {e}")
            return None
    
    def _calculate_overall_risk_score(self, exploit_intel: Dict, predictions: Dict) -> float:
        """Calculate overall risk score from all analyses"""
        risk_factors = []
        
        # Exploit intelligence factors
        if exploit_intel.get('exploit_availability'):
            risk_factors.append(0.3)
        if exploit_intel.get('cisa_kev'):
            risk_factors.append(0.2)
        if exploit_intel.get('exploit_likelihood', 0) > 0.7:
            risk_factors.append(0.2)
        
        # Predictive analysis factors
        exploit_likelihood = predictions.get('exploit_likelihood', {})
        if exploit_likelihood.get('exploit_likelihood', 0) > 0.7:
            risk_factors.append(0.2)
        
        business_impact = predictions.get('business_impact', {})
        if business_impact.get('business_impact_score', 0) > 0.7:
            risk_factors.append(0.2)
        
        timeline_urgency = predictions.get('timeline_urgency', {})
        if timeline_urgency.get('urgency_score', 0) > 0.7:
            risk_factors.append(0.1)
        
        # Calculate overall risk score
        if risk_factors:
            return min(sum(risk_factors), 1.0)
        else:
            return 0.5  # Default medium risk
    
    def _generate_priority_recommendation(self, exploit_intel: Dict, predictions: Dict) -> Dict:
        """Generate priority recommendation based on analysis"""
        risk_score = self._calculate_overall_risk_score(exploit_intel, predictions)
        
        if risk_score >= 0.8:
            priority = 'critical'
            recommendation = 'Immediate action required'
            timeline = '0-24 hours'
        elif risk_score >= 0.6:
            priority = 'high'
            recommendation = 'Urgent action required'
            timeline = '24-72 hours'
        elif risk_score >= 0.4:
            priority = 'medium'
            recommendation = 'Schedule remediation'
            timeline = '1-2 weeks'
        else:
            priority = 'low'
            recommendation = 'Include in regular patch cycle'
            timeline = '2-4 weeks'
        
        return {
            'priority': priority,
            'recommendation': recommendation,
            'timeline': timeline,
            'risk_score': risk_score,
            'justification': self._generate_priority_justification(exploit_intel, predictions, risk_score)
        }
    
    def _generate_priority_justification(self, exploit_intel: Dict, predictions: Dict, risk_score: float) -> List[str]:
        """Generate justification for priority recommendation"""
        justifications = []
        
        if exploit_intel.get('exploit_availability'):
            justifications.append('Public exploits are available')
        if exploit_intel.get('cisa_kev'):
            justifications.append('Listed in CISA Known Exploited Vulnerabilities')
        if exploit_intel.get('exploit_likelihood', 0) > 0.7:
            justifications.append('High likelihood of exploitation')
        
        exploit_likelihood = predictions.get('exploit_likelihood', {})
        if exploit_likelihood.get('exploit_likelihood', 0) > 0.7:
            justifications.append('AI predicts high exploit likelihood')
        
        business_impact = predictions.get('business_impact', {})
        if business_impact.get('business_impact_score', 0) > 0.7:
            justifications.append('High predicted business impact')
        
        timeline_urgency = predictions.get('timeline_urgency', {})
        if timeline_urgency.get('urgency_score', 0) > 0.7:
            justifications.append('High timeline urgency')
        
        if not justifications:
            justifications.append('Standard vulnerability requiring attention')
        
        return justifications
    
    def _generate_action_items(self, exploit_intel: Dict, predictions: Dict, remediation_plan: Optional[RemediationPlan]) -> List[Dict]:
        """Generate actionable items based on analysis"""
        action_items = []
        
        # Immediate actions based on exploit intelligence
        if exploit_intel.get('exploit_availability'):
            action_items.append({
                'action': 'Review access logs for exploitation attempts',
                'priority': 'critical',
                'timeline': 'immediate',
                'category': 'monitoring'
            })
        
        if exploit_intel.get('cisa_kev'):
            action_items.append({
                'action': 'Implement emergency patching procedures',
                'priority': 'critical',
                'timeline': '0-24 hours',
                'category': 'remediation'
            })
        
        # Actions based on predictive analysis
        exploit_likelihood = predictions.get('exploit_likelihood', {})
        if exploit_likelihood.get('exploit_likelihood', 0) > 0.7:
            action_items.append({
                'action': 'Implement additional monitoring and alerting',
                'priority': 'high',
                'timeline': '24-48 hours',
                'category': 'security'
            })
        
        business_impact = predictions.get('business_impact', {})
        if business_impact.get('business_impact_score', 0) > 0.7:
            action_items.append({
                'action': 'Prepare incident response procedures',
                'priority': 'high',
                'timeline': '24-48 hours',
                'category': 'preparation'
            })
        
        # Remediation actions
        if remediation_plan:
            action_items.append({
                'action': f'Execute remediation plan: {remediation_plan.title}',
                'priority': 'high',
                'timeline': remediation_plan.estimated_total_time,
                'category': 'remediation',
                'plan_id': remediation_plan.plan_id
            })
        
        return action_items
    
    def _calculate_confidence_score(self, exploit_intel: Dict, predictions: Dict) -> float:
        """Calculate overall confidence score for the analysis"""
        confidence_factors = []
        
        # Exploit intelligence confidence
        if 'confidence_score' in exploit_intel:
            confidence_factors.append(exploit_intel['confidence_score'])
        
        # Predictive analysis confidence
        for analysis_type in ['exploit_likelihood', 'business_impact', 'remediation_complexity', 'timeline_urgency']:
            analysis = predictions.get(analysis_type, {})
            if 'confidence_score' in analysis:
                confidence_factors.append(analysis['confidence_score'])
        
        # Calculate average confidence
        if confidence_factors:
            return sum(confidence_factors) / len(confidence_factors)
        else:
            return 0.7  # Default confidence
    
    def _get_fallback_analysis(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """Get fallback analysis when AI services are unavailable"""
        severity = vulnerability_data.get('severity', 'MEDIUM')
        
        # Simple heuristic-based analysis
        risk_score = 0.5
        if severity == 'CRITICAL':
            risk_score = 0.9
        elif severity == 'HIGH':
            risk_score = 0.7
        elif severity == 'MEDIUM':
            risk_score = 0.5
        else:
            risk_score = 0.3
        
        return {
            'vulnerability_id': vulnerability_data.get('id'),
            'analysis_timestamp': datetime.utcnow().isoformat(),
            'exploit_intelligence': {
                'exploit_availability': False,
                'confidence_score': 0.5
            },
            'predictive_analysis': {
                'exploit_likelihood': {'exploit_likelihood': risk_score, 'confidence_score': 0.6},
                'business_impact': {'business_impact_score': risk_score, 'confidence_score': 0.6},
                'remediation_complexity': {'complexity_level': 'medium', 'confidence_score': 0.6},
                'timeline_urgency': {'urgency_score': risk_score, 'confidence_score': 0.6}
            },
            'remediation_plan': None,
            'overall_risk_score': risk_score,
            'priority_recommendation': {
                'priority': 'high' if risk_score >= 0.7 else 'medium',
                'recommendation': 'Standard remediation required',
                'timeline': '1-2 weeks',
                'risk_score': risk_score
            },
            'action_items': [
                {
                    'action': 'Review vulnerability details',
                    'priority': 'medium',
                    'timeline': '24-48 hours',
                    'category': 'assessment'
                }
            ],
            'confidence_score': 0.6,
            'fallback_mode': True
        }
    
    async def analyze_vulnerability_trends(self, vulnerabilities: List[Dict], time_period: str = "30d") -> Dict:
        """Analyze vulnerability trends using AI"""
        try:
            # Use AI service for trend analysis
            if self.ai_service:
                return await self.ai_service.analyze_vulnerability_trends(vulnerabilities, time_period)
            else:
                return self._get_fallback_trend_analysis(vulnerabilities)
        except Exception as e:
            logger.error(f"Error analyzing vulnerability trends: {e}")
            return self._get_fallback_trend_analysis(vulnerabilities)
    
    def _get_fallback_trend_analysis(self, vulnerabilities: List[Dict]) -> Dict:
        """Get fallback trend analysis when AI is unavailable"""
        severity_counts = {}
        type_counts = {}
        
        for vuln in vulnerabilities:
            severity = vuln.get('severity', 'UNKNOWN')
            vuln_type = vuln.get('type', 'UNKNOWN')
            
            severity_counts[severity] = severity_counts.get(severity, 0) + 1
            type_counts[vuln_type] = type_counts.get(vuln_type, 0) + 1
        
        total_vulns = len(vulnerabilities)
        critical_count = severity_counts.get('CRITICAL', 0)
        high_count = severity_counts.get('HIGH', 0)
        
        return {
            'trends': [
                f"Total of {total_vulns} vulnerabilities analyzed",
                f"Critical: {critical_count}, High: {high_count}",
                "Most common types: " + ", ".join(list(type_counts.keys())[:3])
            ],
            'high_risk_areas': [
                "Critical and high severity vulnerabilities",
                "Frequently affected packages",
                "Unpatched systems"
            ],
            'recommendations': [
                "Prioritize critical and high severity vulnerabilities",
                "Implement automated patching for common vulnerabilities",
                "Regular security scanning and monitoring"
            ],
            'risk_score': min((critical_count * 1.0 + high_count * 0.7) / max(total_vulns, 1), 1.0),
            'confidence': 0.6,
            'summary': f"Analysis of {total_vulns} vulnerabilities with {critical_count} critical and {high_count} high severity issues",
            'generated_at': datetime.utcnow().isoformat(),
            'vulnerability_count': total_vulns,
            'ai_provider': 'fallback'
        }
    
    async def get_plan_effectiveness(self, plan_id: str, feedback_data: Dict) -> Dict:
        """Analyze effectiveness of a remediation plan"""
        try:
            return await self.remediation_guidance_service.get_plan_effectiveness(plan_id, feedback_data)
        except Exception as e:
            logger.error(f"Error analyzing plan effectiveness: {e}")
            return {
                'plan_id': plan_id,
                'effectiveness_score': 0.5,
                'error': str(e)
            }
    
    async def get_bulk_analysis(self, vulnerabilities: List[Dict], organization_context: Dict = None) -> Dict[str, Dict]:
        """Perform bulk analysis for multiple vulnerabilities"""
        try:
            results = {}
            
            # Process vulnerabilities in batches to respect rate limits
            batch_size = 3
            for i in range(0, len(vulnerabilities), batch_size):
                batch = vulnerabilities[i:i + batch_size]
                
                tasks = [self.analyze_vulnerability_comprehensive(vuln, organization_context) for vuln in batch]
                batch_results = await asyncio.gather(*tasks, return_exceptions=True)
                
                for j, result in enumerate(batch_results):
                    if not isinstance(result, Exception):
                        results[batch[j].get('id', f'vuln_{i+j}')] = result
                    else:
                        logger.error(f"Error processing {batch[j].get('id', f'vuln_{i+j}')}: {result}")
                        results[batch[j].get('id', f'vuln_{i+j}')] = {'error': str(result)}
                
                # Rate limiting delay between batches
                if i + batch_size < len(vulnerabilities):
                    await asyncio.sleep(2)
            
            return results
            
        except Exception as e:
            logger.error(f"Error in bulk analysis: {e}")
            return {}

# Global integrated AI service instance
integrated_ai_service = IntegratedAIService()

