"""
AI Service for ZeroTrace - Provides AI-powered vulnerability analysis and remediation guidance
"""

import asyncio
import json
import logging
import os
import time
from typing import Dict, List, Optional, Any
from datetime import datetime, timedelta
import aiohttp
import openai
from anthropic import Anthropic

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class AIService:
    """Main AI service class that handles AI-powered operations"""
    
    def __init__(self):
        self.openai_client = None
        self.anthropic_client = None
        self.cache = {}
        self.cache_ttl = 3600  # 1 hour cache TTL
        self.rate_limits = {
            'openai': {'requests': 0, 'reset_time': 0},
            'anthropic': {'requests': 0, 'reset_time': 0}
        }
        self.max_requests_per_minute = 60
        
        # Initialize AI clients
        self._initialize_clients()
    
    def _initialize_clients(self):
        """Initialize AI clients with API keys"""
        try:
            # OpenAI client
            openai_api_key = os.getenv('OPENAI_API_KEY')
            if openai_api_key:
                openai.api_key = openai_api_key
                self.openai_client = openai
                logger.info("OpenAI client initialized")
            else:
                logger.warning("OpenAI API key not found")
            
            # Anthropic client
            anthropic_api_key = os.getenv('ANTHROPIC_API_KEY')
            if anthropic_api_key:
                self.anthropic_client = Anthropic(api_key=anthropic_api_key)
                logger.info("Anthropic client initialized")
            else:
                logger.warning("Anthropic API key not found")
                
        except Exception as e:
            logger.error(f"Failed to initialize AI clients: {e}")
    
    async def _check_rate_limit(self, provider: str) -> bool:
        """Check if we're within rate limits for the provider"""
        current_time = time.time()
        rate_limit = self.rate_limits[provider]
        
        # Reset counter if minute has passed
        if current_time > rate_limit['reset_time']:
            rate_limit['requests'] = 0
            rate_limit['reset_time'] = current_time + 60
        
        return rate_limit['requests'] < self.max_requests_per_minute
    
    async def _increment_rate_limit(self, provider: str):
        """Increment rate limit counter for provider"""
        self.rate_limits[provider]['requests'] += 1
    
    async def _get_cached_response(self, cache_key: str) -> Optional[Dict]:
        """Get cached response if available and not expired"""
        if cache_key in self.cache:
            cached_data = self.cache[cache_key]
            if time.time() - cached_data['timestamp'] < self.cache_ttl:
                logger.debug(f"Returning cached response for {cache_key}")
                return cached_data['response']
            else:
                # Remove expired cache entry
                del self.cache[cache_key]
        return None
    
    async def _cache_response(self, cache_key: str, response: Dict):
        """Cache response with timestamp"""
        self.cache[cache_key] = {
            'response': response,
            'timestamp': time.time()
        }
        logger.debug(f"Cached response for {cache_key}")
    
    async def generate_remediation_guidance(self, vulnerability: Dict, organization_context: Dict = None) -> Dict:
        """
        Generate AI-powered remediation guidance for a vulnerability
        
        Args:
            vulnerability: Vulnerability details
            organization_context: Organization-specific context (industry, tech stack, etc.)
        
        Returns:
            Dict containing remediation guidance
        """
        cache_key = f"remediation_{vulnerability.get('id', 'unknown')}_{hash(str(organization_context))}"
        
        # Check cache first
        cached_response = await self._get_cached_response(cache_key)
        if cached_response:
            return cached_response
        
        # Prepare prompt
        prompt = self._build_remediation_prompt(vulnerability, organization_context)
        
        # Try OpenAI first, fallback to Anthropic
        response = None
        try:
            if self.openai_client and await self._check_rate_limit('openai'):
                response = await self._generate_with_openai(prompt)
                await self._increment_rate_limit('openai')
            elif self.anthropic_client and await self._check_rate_limit('anthropic'):
                response = await self._generate_with_anthropic(prompt)
                await self._increment_rate_limit('anthropic')
            else:
                logger.warning("No AI provider available or rate limited")
                return self._get_fallback_remediation(vulnerability)
        except Exception as e:
            logger.error(f"AI generation failed: {e}")
            return self._get_fallback_remediation(vulnerability)
        
        # Parse and structure response
        guidance = self._parse_remediation_response(response, vulnerability)
        
        # Cache the response
        await self._cache_response(cache_key, guidance)
        
        return guidance
    
    def _build_remediation_prompt(self, vulnerability: Dict, organization_context: Dict = None) -> str:
        """Build prompt for remediation guidance"""
        base_prompt = f"""
        You are a cybersecurity expert providing remediation guidance for a vulnerability.
        
        Vulnerability Details:
        - ID: {vulnerability.get('id', 'Unknown')}
        - Title: {vulnerability.get('title', 'Unknown')}
        - Description: {vulnerability.get('description', 'No description available')}
        - Severity: {vulnerability.get('severity', 'Unknown')}
        - CVSS Score: {vulnerability.get('cvss_score', 'Unknown')}
        - Package: {vulnerability.get('package_name', 'Unknown')}
        - Version: {vulnerability.get('package_version', 'Unknown')}
        - CVE: {vulnerability.get('cve_id', 'Not available')}
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
        
        Please provide comprehensive remediation guidance including:
        1. Immediate actions (if critical/high severity)
        2. Step-by-step remediation process
        3. Testing and validation steps
        4. Prevention measures for similar vulnerabilities
        5. Timeline recommendations
        6. Risk assessment of the vulnerability
        
        Format your response as JSON with the following structure:
        {
            "immediate_actions": ["action1", "action2"],
            "remediation_steps": [
                {"step": 1, "description": "description", "estimated_time": "time"},
                {"step": 2, "description": "description", "estimated_time": "time"}
            ],
            "testing_steps": ["test1", "test2"],
            "prevention_measures": ["measure1", "measure2"],
            "timeline": "estimated_timeline",
            "risk_assessment": "risk_description",
            "confidence_score": 0.95
        }
        """
        
        return base_prompt
    
    async def _generate_with_openai(self, prompt: str) -> str:
        """Generate response using OpenAI"""
        try:
            response = await asyncio.to_thread(
                self.openai_client.ChatCompletion.create,
                model="gpt-4",
                messages=[
                    {"role": "system", "content": "You are a cybersecurity expert providing detailed remediation guidance."},
                    {"role": "user", "content": prompt}
                ],
                max_tokens=2000,
                temperature=0.3
            )
            return response.choices[0].message.content
        except Exception as e:
            logger.error(f"OpenAI API error: {e}")
            raise
    
    async def _generate_with_anthropic(self, prompt: str) -> str:
        """Generate response using Anthropic"""
        try:
            response = await asyncio.to_thread(
                self.anthropic_client.messages.create,
                model="claude-3-sonnet-20240229",
                max_tokens=2000,
                messages=[{"role": "user", "content": prompt}]
            )
            return response.content[0].text
        except Exception as e:
            logger.error(f"Anthropic API error: {e}")
            raise
    
    def _parse_remediation_response(self, response: str, vulnerability: Dict) -> Dict:
        """Parse AI response and structure it"""
        try:
            # Try to extract JSON from response
            start_idx = response.find('{')
            end_idx = response.rfind('}') + 1
            
            if start_idx != -1 and end_idx != -1:
                json_str = response[start_idx:end_idx]
                parsed_response = json.loads(json_str)
                
                # Add metadata
                parsed_response['generated_at'] = datetime.utcnow().isoformat()
                parsed_response['vulnerability_id'] = vulnerability.get('id')
                parsed_response['ai_provider'] = 'openai' if self.openai_client else 'anthropic'
                
                return parsed_response
            else:
                # Fallback if JSON parsing fails
                return self._get_fallback_remediation(vulnerability)
                
        except json.JSONDecodeError as e:
            logger.error(f"Failed to parse AI response as JSON: {e}")
            return self._get_fallback_remediation(vulnerability)
    
    def _get_fallback_remediation(self, vulnerability: Dict) -> Dict:
        """Get fallback remediation guidance when AI is unavailable"""
        severity = vulnerability.get('severity', 'MEDIUM').upper()
        
        immediate_actions = []
        if severity in ['CRITICAL', 'HIGH']:
            immediate_actions = [
                "Isolate affected systems if possible",
                "Review access logs for potential exploitation",
                "Implement temporary workarounds if available"
            ]
        
        remediation_steps = [
            {
                "step": 1,
                "description": f"Update {vulnerability.get('package_name', 'affected package')} to latest secure version",
                "estimated_time": "1-4 hours"
            },
            {
                "step": 2,
                "description": "Test the update in a non-production environment",
                "estimated_time": "2-8 hours"
            },
            {
                "step": 3,
                "description": "Deploy update to production with monitoring",
                "estimated_time": "1-2 hours"
            }
        ]
        
        return {
            "immediate_actions": immediate_actions,
            "remediation_steps": remediation_steps,
            "testing_steps": [
                "Verify the update resolves the vulnerability",
                "Run security tests to ensure no new issues",
                "Monitor system stability post-update"
            ],
            "prevention_measures": [
                "Implement automated dependency scanning",
                "Set up security monitoring and alerting",
                "Regular security assessments and updates"
            ],
            "timeline": "1-3 days" if severity in ['CRITICAL', 'HIGH'] else "1-2 weeks",
            "risk_assessment": f"High risk due to {severity} severity vulnerability",
            "confidence_score": 0.7,
            "generated_at": datetime.utcnow().isoformat(),
            "vulnerability_id": vulnerability.get('id'),
            "ai_provider": "fallback"
        }
    
    async def analyze_vulnerability_trends(self, vulnerabilities: List[Dict], time_period: str = "30d") -> Dict:
        """
        Analyze vulnerability trends using AI
        
        Args:
            vulnerabilities: List of vulnerability data
            time_period: Time period for analysis (7d, 30d, 90d)
        
        Returns:
            Dict containing trend analysis
        """
        cache_key = f"trends_{hash(str(vulnerabilities))}_{time_period}"
        
        # Check cache first
        cached_response = await self._get_cached_response(cache_key)
        if cached_response:
            return cached_response
        
        # Build trend analysis prompt
        prompt = self._build_trend_analysis_prompt(vulnerabilities, time_period)
        
        try:
            if self.openai_client and await self._check_rate_limit('openai'):
                response = await self._generate_with_openai(prompt)
                await self._increment_rate_limit('openai')
            elif self.anthropic_client and await self._check_rate_limit('anthropic'):
                response = await self._generate_with_anthropic(prompt)
                await self._increment_rate_limit('anthropic')
            else:
                return self._get_fallback_trend_analysis(vulnerabilities)
        except Exception as e:
            logger.error(f"Trend analysis failed: {e}")
            return self._get_fallback_trend_analysis(vulnerabilities)
        
        # Parse response
        analysis = self._parse_trend_response(response, vulnerabilities)
        
        # Cache the response
        await self._cache_response(cache_key, analysis)
        
        return analysis
    
    def _build_trend_analysis_prompt(self, vulnerabilities: List[Dict], time_period: str) -> str:
        """Build prompt for trend analysis"""
        vuln_summary = {
            'total': len(vulnerabilities),
            'by_severity': {},
            'by_type': {},
            'by_package': {}
        }
        
        for vuln in vulnerabilities:
            severity = vuln.get('severity', 'UNKNOWN')
            vuln_type = vuln.get('type', 'UNKNOWN')
            package = vuln.get('package_name', 'UNKNOWN')
            
            vuln_summary['by_severity'][severity] = vuln_summary['by_severity'].get(severity, 0) + 1
            vuln_summary['by_type'][vuln_type] = vuln_summary['by_type'].get(vuln_type, 0) + 1
            vuln_summary['by_package'][package] = vuln_summary['by_package'].get(package, 0) + 1
        
        prompt = f"""
        Analyze the following vulnerability data for trends and patterns over the last {time_period}:
        
        Vulnerability Summary:
        - Total vulnerabilities: {vuln_summary['total']}
        - By severity: {vuln_summary['by_severity']}
        - By type: {vuln_summary['by_type']}
        - By package: {vuln_summary['by_package']}
        
        Please provide analysis including:
        1. Key trends and patterns
        2. High-risk areas requiring attention
        3. Recommendations for improvement
        4. Risk score and confidence level
        
        Format as JSON:
        {{
            "trends": ["trend1", "trend2"],
            "high_risk_areas": ["area1", "area2"],
            "recommendations": ["rec1", "rec2"],
            "risk_score": 0.75,
            "confidence": 0.9,
            "summary": "analysis_summary"
        }}
        """
        
        return prompt
    
    def _parse_trend_response(self, response: str, vulnerabilities: List[Dict]) -> Dict:
        """Parse trend analysis response"""
        try:
            start_idx = response.find('{')
            end_idx = response.rfind('}') + 1
            
            if start_idx != -1 and end_idx != -1:
                json_str = response[start_idx:end_idx]
                parsed_response = json.loads(json_str)
                
                # Add metadata
                parsed_response['generated_at'] = datetime.utcnow().isoformat()
                parsed_response['vulnerability_count'] = len(vulnerabilities)
                parsed_response['ai_provider'] = 'openai' if self.openai_client else 'anthropic'
                
                return parsed_response
            else:
                return self._get_fallback_trend_analysis(vulnerabilities)
                
        except json.JSONDecodeError as e:
            logger.error(f"Failed to parse trend response as JSON: {e}")
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
        
        # Calculate basic risk score
        total_vulns = len(vulnerabilities)
        critical_count = severity_counts.get('CRITICAL', 0)
        high_count = severity_counts.get('HIGH', 0)
        
        risk_score = (critical_count * 1.0 + high_count * 0.7) / max(total_vulns, 1)
        
        return {
            "trends": [
                f"Total of {total_vulns} vulnerabilities detected",
                f"Critical: {critical_count}, High: {high_count}",
                "Most common types: " + ", ".join(list(type_counts.keys())[:3])
            ],
            "high_risk_areas": [
                "Critical and high severity vulnerabilities",
                "Frequently affected packages",
                "Unpatched systems"
            ],
            "recommendations": [
                "Prioritize critical and high severity vulnerabilities",
                "Implement automated patching for common vulnerabilities",
                "Regular security scanning and monitoring"
            ],
            "risk_score": min(risk_score, 1.0),
            "confidence": 0.6,
            "summary": f"Analysis of {total_vulns} vulnerabilities with {critical_count} critical and {high_count} high severity issues",
            "generated_at": datetime.utcnow().isoformat(),
            "vulnerability_count": total_vulns,
            "ai_provider": "fallback"
        }

# Global AI service instance
ai_service = AIService()

