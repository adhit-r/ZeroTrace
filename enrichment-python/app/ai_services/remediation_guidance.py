"""
AI-Generated Remediation Guidance Service - LLM-based remediation plan generation
"""

import asyncio
import json
import logging
import os
import time
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any
from dataclasses import dataclass

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class RemediationStep:
    """Represents a single step in a remediation plan"""
    step_number: int
    title: str
    description: str
    estimated_time: str
    required_skills: List[str]
    prerequisites: List[str]
    testing_required: bool
    rollback_plan: str
    priority: str  # critical, high, medium, low

@dataclass
class RemediationPlan:
    """Complete remediation plan for a vulnerability"""
    vulnerability_id: str
    plan_id: str
    title: str
    description: str
    severity: str
    estimated_total_time: str
    complexity_level: str
    steps: List[RemediationStep]
    testing_plan: List[str]
    rollback_strategy: str
    success_criteria: List[str]
    risk_assessment: Dict[str, Any]
    industry_specific_notes: List[str]
    compliance_considerations: List[str]
    confidence_score: float
    generated_at: str
    ai_provider: str

class RemediationGuidanceService:
    """Service for AI-generated remediation guidance and plans"""
    
    def __init__(self, ai_service=None):
        self.ai_service = ai_service
        self.cache = {}
        self.cache_ttl = 7200  # 2 hours cache TTL
        self.plan_templates = self._load_remediation_templates()
        
    def _load_remediation_templates(self) -> Dict[str, Dict]:
        """Load industry and vulnerability type specific templates"""
        return {
            'web_application': {
                'critical': {
                    'immediate_actions': [
                        'Isolate affected web servers',
                        'Review access logs for exploitation attempts',
                        'Implement temporary workarounds if available'
                    ],
                    'testing_focus': ['Security testing', 'Performance testing', 'Regression testing'],
                    'rollback_priority': 'high'
                },
                'high': {
                    'immediate_actions': [
                        'Schedule maintenance window',
                        'Prepare rollback plan',
                        'Notify stakeholders'
                    ],
                    'testing_focus': ['Functional testing', 'Security testing'],
                    'rollback_priority': 'medium'
                }
            },
            'database': {
                'critical': {
                    'immediate_actions': [
                        'Backup database immediately',
                        'Review database access logs',
                        'Implement network segmentation'
                    ],
                    'testing_focus': ['Data integrity testing', 'Performance testing'],
                    'rollback_priority': 'critical'
                }
            },
            'authentication': {
                'critical': {
                    'immediate_actions': [
                        'Force password reset for all users',
                        'Review authentication logs',
                        'Implement additional monitoring'
                    ],
                    'testing_focus': ['Authentication testing', 'Authorization testing'],
                    'rollback_priority': 'high'
                }
            }
        }
    
    async def generate_remediation_plan(self, vulnerability_data: Dict, organization_context: Dict = None) -> RemediationPlan:
        """
        Generate comprehensive remediation plan for a vulnerability
        
        Args:
            vulnerability_data: Vulnerability details
            organization_context: Organization-specific context
        
        Returns:
            Complete remediation plan
        """
        cache_key = f"remediation_{vulnerability_data.get('id', 'unknown')}_{hash(str(organization_context))}"
        
        # Check cache first
        cached_plan = await self._get_cached_plan(cache_key)
        if cached_plan:
            return cached_plan
        
        try:
            # Generate plan using AI service if available
            if self.ai_service:
                plan_data = await self._generate_ai_plan(vulnerability_data, organization_context)
            else:
                plan_data = await self._generate_template_plan(vulnerability_data, organization_context)
            
            # Create remediation plan object
            plan = self._create_remediation_plan(vulnerability_data, plan_data, organization_context)
            
            # Cache the plan
            await self._cache_plan(cache_key, plan)
            
            return plan
            
        except Exception as e:
            logger.error(f"Error generating remediation plan: {e}")
            return self._create_fallback_plan(vulnerability_data, organization_context)
    
    async def _generate_ai_plan(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """Generate remediation plan using AI service"""
        prompt = self._build_remediation_prompt(vulnerability_data, organization_context)
        
        try:
            response = await self.ai_service.generate_remediation_guidance(vulnerability_data, organization_context)
            return response
        except Exception as e:
            logger.error(f"AI service error: {e}")
            raise
    
    async def _generate_template_plan(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict:
        """Generate remediation plan using templates"""
        vuln_type = self._classify_vulnerability_type(vulnerability_data)
        severity = vulnerability_data.get('severity', 'MEDIUM')
        
        # Get template for vulnerability type and severity
        template = self.plan_templates.get(vuln_type, {}).get(severity.lower(), {})
        
        # Generate steps based on template
        steps = self._generate_template_steps(vulnerability_data, template, organization_context)
        
        return {
            'immediate_actions': template.get('immediate_actions', []),
            'remediation_steps': steps,
            'testing_steps': template.get('testing_focus', []),
            'prevention_measures': self._generate_prevention_measures(vulnerability_data, organization_context),
            'timeline': self._estimate_timeline(severity, vuln_type),
            'risk_assessment': self._assess_risk(vulnerability_data, organization_context)
        }
    
    def _classify_vulnerability_type(self, vulnerability_data: Dict) -> str:
        """Classify vulnerability type based on data"""
        package_name = vulnerability_data.get('package_name', '').lower()
        vuln_type = vulnerability_data.get('type', '').lower()
        description = vulnerability_data.get('description', '').lower()
        
        # Web application vulnerabilities
        if any(keyword in package_name for keyword in ['web', 'http', 'server', 'api']):
            return 'web_application'
        elif any(keyword in vuln_type for keyword in ['xss', 'csrf', 'injection', 'web']):
            return 'web_application'
        
        # Database vulnerabilities
        elif any(keyword in package_name for keyword in ['database', 'db', 'sql', 'mysql', 'postgres']):
            return 'database'
        elif any(keyword in vuln_type for keyword in ['sql', 'database', 'db']):
            return 'database'
        
        # Authentication vulnerabilities
        elif any(keyword in package_name for keyword in ['auth', 'login', 'session', 'token']):
            return 'authentication'
        elif any(keyword in vuln_type for keyword in ['auth', 'session', 'token', 'credential']):
            return 'authentication'
        
        # Default classification
        else:
            return 'general'
    
    def _generate_template_steps(self, vulnerability_data: Dict, template: Dict, organization_context: Dict = None) -> List[Dict]:
        """Generate remediation steps from template"""
        steps = []
        step_number = 1
        
        # Immediate actions
        for action in template.get('immediate_actions', []):
            steps.append({
                'step': step_number,
                'description': action,
                'estimated_time': '1-4 hours',
                'priority': 'critical'
            })
            step_number += 1
        
        # Standard remediation steps
        standard_steps = [
            {
                'description': f"Update {vulnerability_data.get('package_name', 'affected package')} to latest secure version",
                'estimated_time': '2-8 hours',
                'priority': 'high'
            },
            {
                'description': 'Test the update in a non-production environment',
                'estimated_time': '4-16 hours',
                'priority': 'high'
            },
            {
                'description': 'Deploy update to production with monitoring',
                'estimated_time': '1-4 hours',
                'priority': 'high'
            },
            {
                'description': 'Verify the update resolves the vulnerability',
                'estimated_time': '2-8 hours',
                'priority': 'medium'
            }
        ]
        
        for step in standard_steps:
            steps.append({
                'step': step_number,
                'description': step['description'],
                'estimated_time': step['estimated_time'],
                'priority': step['priority']
            })
            step_number += 1
        
        return steps
    
    def _generate_prevention_measures(self, vulnerability_data: Dict, organization_context: Dict = None) -> List[str]:
        """Generate prevention measures based on vulnerability and organization context"""
        measures = [
            'Implement automated dependency scanning',
            'Set up security monitoring and alerting',
            'Regular security assessments and updates'
        ]
        
        # Industry-specific measures
        if organization_context:
            industry = organization_context.get('industry', '').lower()
            if industry == 'healthcare':
                measures.extend([
                    'HIPAA compliance monitoring',
                    'Patient data protection measures'
                ])
            elif industry == 'finance':
                measures.extend([
                    'PCI DSS compliance monitoring',
                    'Financial data encryption'
                ])
            elif industry == 'government':
                measures.extend([
                    'FISMA compliance monitoring',
                    'Security clearance verification'
                ])
        
        # Vulnerability-specific measures
        vuln_type = self._classify_vulnerability_type(vulnerability_data)
        if vuln_type == 'web_application':
            measures.extend([
                'Web Application Firewall (WAF) implementation',
                'Input validation and sanitization'
            ])
        elif vuln_type == 'database':
            measures.extend([
                'Database access controls',
                'Encryption at rest and in transit'
            ])
        elif vuln_type == 'authentication':
            measures.extend([
                'Multi-factor authentication',
                'Session management controls'
            ])
        
        return measures
    
    def _estimate_timeline(self, severity: str, vuln_type: str) -> str:
        """Estimate remediation timeline based on severity and type"""
        base_times = {
            'CRITICAL': '1-3 days',
            'HIGH': '3-7 days',
            'MEDIUM': '1-2 weeks',
            'LOW': '2-4 weeks'
        }
        
        # Adjust based on vulnerability type complexity
        if vuln_type in ['database', 'authentication']:
            # More complex, add time
            if severity == 'CRITICAL':
                return '2-5 days'
            elif severity == 'HIGH':
                return '5-10 days'
        
        return base_times.get(severity, '1-2 weeks')
    
    def _assess_risk(self, vulnerability_data: Dict, organization_context: Dict = None) -> Dict[str, Any]:
        """Assess risk factors for the vulnerability"""
        risk_factors = []
        risk_score = 0.0
        
        # Severity-based risk
        severity = vulnerability_data.get('severity', 'MEDIUM')
        if severity == 'CRITICAL':
            risk_score += 0.8
            risk_factors.append('Critical severity vulnerability')
        elif severity == 'HIGH':
            risk_score += 0.6
            risk_factors.append('High severity vulnerability')
        
        # Exploit availability
        if vulnerability_data.get('exploit_available'):
            risk_score += 0.3
            risk_factors.append('Public exploits available')
        
        # CVSS score
        cvss_score = vulnerability_data.get('cvss_score', 0.0) or 0.0
        if cvss_score >= 8.0:
            risk_score += 0.4
            risk_factors.append('High CVSS score')
        
        # Organization context
        if organization_context:
            industry = organization_context.get('industry', '').lower()
            if industry in ['healthcare', 'finance', 'government']:
                risk_score += 0.2
                risk_factors.append('Regulated industry')
        
        return {
            'risk_score': min(risk_score, 1.0),
            'risk_level': 'critical' if risk_score >= 0.8 else 'high' if risk_score >= 0.6 else 'medium',
            'risk_factors': risk_factors,
            'mitigation_priority': 'immediate' if risk_score >= 0.8 else 'high' if risk_score >= 0.6 else 'medium'
        }
    
    def _create_remediation_plan(self, vulnerability_data: Dict, plan_data: Dict, organization_context: Dict = None) -> RemediationPlan:
        """Create RemediationPlan object from data"""
        steps = []
        
        # Create remediation steps
        for i, step_data in enumerate(plan_data.get('remediation_steps', []), 1):
            step = RemediationStep(
                step_number=i,
                title=f"Step {i}",
                description=step_data.get('description', ''),
                estimated_time=step_data.get('estimated_time', 'Unknown'),
                required_skills=self._identify_required_skills(step_data, organization_context),
                prerequisites=self._identify_prerequisites(step_data, i),
                testing_required=step_data.get('priority', 'medium') in ['high', 'critical'],
                rollback_plan=self._generate_rollback_plan(step_data),
                priority=step_data.get('priority', 'medium')
            )
            steps.append(step)
        
        # Generate plan metadata
        plan_id = f"plan_{vulnerability_data.get('id', 'unknown')}_{int(time.time())}"
        
        return RemediationPlan(
            vulnerability_id=vulnerability_data.get('id', 'unknown'),
            plan_id=plan_id,
            title=f"Remediation Plan for {vulnerability_data.get('title', 'Vulnerability')}",
            description=plan_data.get('description', 'AI-generated remediation plan'),
            severity=vulnerability_data.get('severity', 'MEDIUM'),
            estimated_total_time=plan_data.get('timeline', '1-2 weeks'),
            complexity_level=self._assess_complexity(vulnerability_data, steps),
            steps=steps,
            testing_plan=plan_data.get('testing_steps', []),
            rollback_strategy=self._generate_rollback_strategy(vulnerability_data),
            success_criteria=self._define_success_criteria(vulnerability_data),
            risk_assessment=plan_data.get('risk_assessment', {}),
            industry_specific_notes=self._generate_industry_notes(organization_context),
            compliance_considerations=self._generate_compliance_notes(organization_context),
            confidence_score=0.8,  # Default confidence
            generated_at=datetime.utcnow().isoformat(),
            ai_provider='template' if not self.ai_service else 'ai'
        )
    
    def _identify_required_skills(self, step_data: Dict, organization_context: Dict = None) -> List[str]:
        """Identify required skills for a remediation step"""
        skills = ['System administration', 'Security patching']
        
        description = step_data.get('description', '').lower()
        
        if 'database' in description:
            skills.append('Database administration')
        if 'web' in description or 'server' in description:
            skills.append('Web server administration')
        if 'network' in description:
            skills.append('Network administration')
        if 'security' in description:
            skills.append('Security analysis')
        
        return skills
    
    def _identify_prerequisites(self, step_data: Dict, step_number: int) -> List[str]:
        """Identify prerequisites for a remediation step"""
        prerequisites = []
        
        if step_number > 1:
            prerequisites.append(f"Complete Step {step_number - 1}")
        
        if step_data.get('priority') == 'critical':
            prerequisites.append('Emergency approval')
        
        return prerequisites
    
    def _generate_rollback_plan(self, step_data: Dict) -> str:
        """Generate rollback plan for a step"""
        description = step_data.get('description', '').lower()
        
        if 'update' in description or 'patch' in description:
            return 'Restore previous version and restart services'
        elif 'configuration' in description:
            return 'Restore previous configuration from backup'
        else:
            return 'Reverse all changes made in this step'
    
    def _assess_complexity(self, vulnerability_data: Dict, steps: List[RemediationStep]) -> str:
        """Assess overall complexity of the remediation plan"""
        if len(steps) > 8:
            return 'high'
        elif len(steps) > 5:
            return 'medium'
        else:
            return 'low'
    
    def _generate_rollback_strategy(self, vulnerability_data: Dict) -> str:
        """Generate overall rollback strategy"""
        vuln_type = self._classify_vulnerability_type(vulnerability_data)
        
        if vuln_type == 'database':
            return 'Database backup restoration and service restart'
        elif vuln_type == 'web_application':
            return 'Application rollback and load balancer configuration'
        else:
            return 'System restoration from backup and service restart'
    
    def _define_success_criteria(self, vulnerability_data: Dict) -> List[str]:
        """Define success criteria for remediation"""
        criteria = [
            'Vulnerability is no longer detected in scans',
            'System functionality is restored',
            'No new security issues introduced',
            'Performance is maintained or improved'
        ]
        
        # Add type-specific criteria
        vuln_type = self._classify_vulnerability_type(vulnerability_data)
        if vuln_type == 'web_application':
            criteria.append('Web application security tests pass')
        elif vuln_type == 'database':
            criteria.append('Database integrity checks pass')
        
        return criteria
    
    def _generate_industry_notes(self, organization_context: Dict = None) -> List[str]:
        """Generate industry-specific notes"""
        if not organization_context:
            return []
        
        industry = organization_context.get('industry', '').lower()
        notes = []
        
        if industry == 'healthcare':
            notes.extend([
                'Ensure HIPAA compliance during remediation',
                'Consider patient data protection requirements',
                'Coordinate with compliance team'
            ])
        elif industry == 'finance':
            notes.extend([
                'Ensure PCI DSS compliance during remediation',
                'Consider financial data protection requirements',
                'Coordinate with risk management team'
            ])
        elif industry == 'government':
            notes.extend([
                'Ensure FISMA compliance during remediation',
                'Consider security clearance requirements',
                'Coordinate with security team'
            ])
        
        return notes
    
    def _generate_compliance_notes(self, organization_context: Dict = None) -> List[str]:
        """Generate compliance considerations"""
        if not organization_context:
            return []
        
        frameworks = organization_context.get('compliance_frameworks', [])
        notes = []
        
        for framework in frameworks:
            if framework == 'HIPAA':
                notes.append('Document all changes for HIPAA audit trail')
            elif framework == 'PCI DSS':
                notes.append('Ensure PCI DSS requirements are maintained')
            elif framework == 'SOX':
                notes.append('Document remediation for SOX compliance')
            elif framework == 'ISO27001':
                notes.append('Update ISO27001 documentation')
        
        return notes
    
    def _create_fallback_plan(self, vulnerability_data: Dict, organization_context: Dict = None) -> RemediationPlan:
        """Create fallback plan when AI service is unavailable"""
        steps = [
            RemediationStep(
                step_number=1,
                title="Assess Vulnerability",
                description="Review vulnerability details and impact",
                estimated_time="2-4 hours",
                required_skills=['Security analysis'],
                prerequisites=[],
                testing_required=True,
                rollback_plan="Document current state",
                priority="high"
            ),
            RemediationStep(
                step_number=2,
                title="Apply Patch",
                description=f"Update {vulnerability_data.get('package_name', 'affected package')} to secure version",
                estimated_time="4-8 hours",
                required_skills=['System administration'],
                prerequisites=["Complete Step 1"],
                testing_required=True,
                rollback_plan="Restore previous version",
                priority="high"
            ),
            RemediationStep(
                step_number=3,
                title="Test and Validate",
                description="Test the fix and validate security",
                estimated_time="2-4 hours",
                required_skills=['Security testing'],
                prerequisites=["Complete Step 2"],
                testing_required=True,
                rollback_plan="Restore previous version if issues found",
                priority="medium"
            )
        ]
        
        return RemediationPlan(
            vulnerability_id=vulnerability_data.get('id', 'unknown'),
            plan_id=f"fallback_{int(time.time())}",
            title=f"Fallback Remediation Plan for {vulnerability_data.get('title', 'Vulnerability')}",
            description="Basic remediation plan (AI service unavailable)",
            severity=vulnerability_data.get('severity', 'MEDIUM'),
            estimated_total_time="1-2 days",
            complexity_level="medium",
            steps=steps,
            testing_plan=["Functional testing", "Security testing"],
            rollback_strategy="Restore from backup",
            success_criteria=["Vulnerability resolved", "System functional"],
            risk_assessment={"risk_score": 0.5, "risk_level": "medium"},
            industry_specific_notes=[],
            compliance_considerations=[],
            confidence_score=0.6,
            generated_at=datetime.utcnow().isoformat(),
            ai_provider="fallback"
        )
    
    def _build_remediation_prompt(self, vulnerability_data: Dict, organization_context: Dict = None) -> str:
        """Build prompt for AI service"""
        base_prompt = f"""
        Generate a comprehensive remediation plan for the following vulnerability:
        
        Vulnerability Details:
        - ID: {vulnerability_data.get('id', 'Unknown')}
        - Title: {vulnerability_data.get('title', 'Unknown')}
        - Description: {vulnerability_data.get('description', 'No description available')}
        - Severity: {vulnerability_data.get('severity', 'Unknown')}
        - CVSS Score: {vulnerability_data.get('cvss_score', 'Unknown')}
        - Package: {vulnerability_data.get('package_name', 'Unknown')}
        - Version: {vulnerability_data.get('package_version', 'Unknown')}
        - CVE: {vulnerability_data.get('cve_id', 'Not available')}
        """
        
        if organization_context:
            base_prompt += f"""
        
        Organization Context:
        - Industry: {organization_context.get('industry', 'Unknown')}
        - Risk Tolerance: {organization_context.get('risk_tolerance', 'Unknown')}
        - Tech Stack: {organization_context.get('tech_stack', 'Unknown')}
        - Compliance Frameworks: {organization_context.get('compliance_frameworks', 'None')}
        """
        
        base_prompt += """
        
        Please provide a detailed remediation plan including:
        1. Immediate actions (if critical/high severity)
        2. Step-by-step remediation process with time estimates
        3. Required skills and prerequisites for each step
        4. Testing and validation requirements
        5. Rollback plans for each step
        6. Success criteria
        7. Industry-specific considerations
        8. Compliance requirements
        
        Format as JSON with detailed step information.
        """
        
        return base_prompt
    
    async def _get_cached_plan(self, cache_key: str) -> Optional[RemediationPlan]:
        """Get cached remediation plan"""
        if cache_key in self.cache:
            cached_data = self.cache[cache_key]
            if time.time() - cached_data['timestamp'] < self.cache_ttl:
                return cached_data['plan']
            else:
                del self.cache[cache_key]
        return None
    
    async def _cache_plan(self, cache_key: str, plan: RemediationPlan):
        """Cache remediation plan"""
        self.cache[cache_key] = {
            'plan': plan,
            'timestamp': time.time()
        }
    
    async def get_plan_effectiveness(self, plan_id: str, feedback_data: Dict) -> Dict:
        """Analyze effectiveness of a remediation plan based on feedback"""
        try:
            # Analyze feedback data
            effectiveness_score = self._calculate_effectiveness_score(feedback_data)
            
            # Generate improvement recommendations
            recommendations = self._generate_improvement_recommendations(feedback_data)
            
            return {
                'plan_id': plan_id,
                'effectiveness_score': effectiveness_score,
                'feedback_analysis': self._analyze_feedback(feedback_data),
                'improvement_recommendations': recommendations,
                'success_metrics': self._calculate_success_metrics(feedback_data),
                'analyzed_at': datetime.utcnow().isoformat()
            }
            
        except Exception as e:
            logger.error(f"Error analyzing plan effectiveness: {e}")
            return {
                'plan_id': plan_id,
                'effectiveness_score': 0.5,
                'error': str(e)
            }
    
    def _calculate_effectiveness_score(self, feedback_data: Dict) -> float:
        """Calculate effectiveness score from feedback"""
        score = 0.5  # Base score
        
        # Time to remediation
        if feedback_data.get('time_to_remediation'):
            if feedback_data['time_to_remediation'] <= 24:  # hours
                score += 0.2
            elif feedback_data['time_to_remediation'] <= 72:
                score += 0.1
        
        # Success rate
        if feedback_data.get('success_rate'):
            score += feedback_data['success_rate'] * 0.3
        
        # User satisfaction
        if feedback_data.get('user_satisfaction'):
            score += feedback_data['user_satisfaction'] * 0.2
        
        return min(score, 1.0)
    
    def _generate_improvement_recommendations(self, feedback_data: Dict) -> List[str]:
        """Generate improvement recommendations based on feedback"""
        recommendations = []
        
        if feedback_data.get('time_to_remediation', 0) > 72:
            recommendations.append("Reduce remediation time by automating more steps")
        
        if feedback_data.get('success_rate', 1.0) < 0.8:
            recommendations.append("Improve step clarity and add more detailed instructions")
        
        if feedback_data.get('user_satisfaction', 1.0) < 0.7:
            recommendations.append("Gather more user feedback to improve plan quality")
        
        return recommendations
    
    def _analyze_feedback(self, feedback_data: Dict) -> Dict[str, Any]:
        """Analyze feedback data for insights"""
        return {
            'time_analysis': {
                'average_time': feedback_data.get('time_to_remediation', 0),
                'time_variance': feedback_data.get('time_variance', 0)
            },
            'success_analysis': {
                'success_rate': feedback_data.get('success_rate', 0),
                'failure_reasons': feedback_data.get('failure_reasons', [])
            },
            'user_feedback': {
                'satisfaction_score': feedback_data.get('user_satisfaction', 0),
                'common_issues': feedback_data.get('common_issues', [])
            }
        }
    
    def _calculate_success_metrics(self, feedback_data: Dict) -> Dict[str, Any]:
        """Calculate success metrics from feedback"""
        return {
            'overall_effectiveness': self._calculate_effectiveness_score(feedback_data),
            'time_efficiency': 1.0 - min(feedback_data.get('time_to_remediation', 72) / 168, 1.0),  # Normalize to week
            'user_satisfaction': feedback_data.get('user_satisfaction', 0.5),
            'success_rate': feedback_data.get('success_rate', 0.5)
        }

# Global remediation guidance service instance
remediation_guidance_service = RemediationGuidanceService()
