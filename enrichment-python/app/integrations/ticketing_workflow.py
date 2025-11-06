"""
Automated Ticketing & Workflow - Intelligent ticket creation and routing
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
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import re

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class Ticket:
    """Represents a ticket"""
    ticket_id: str
    title: str
    description: str
    priority: str  # critical, high, medium, low
    status: str    # open, in_progress, resolved, closed
    assignee: Optional[str] = None
    category: str = "security"
    tags: List[str] = None
    created_at: str = ""
    updated_at: str = ""
    due_date: Optional[str] = None
    resolution: Optional[str] = None
    metadata: Dict[str, Any] = None

@dataclass
class WorkflowRule:
    """Represents a workflow rule"""
    rule_id: str
    name: str
    conditions: List[Dict[str, Any]]
    actions: List[Dict[str, Any]]
    priority: int
    enabled: bool = True
    created_at: str = ""
    last_triggered: Optional[str] = None
    trigger_count: int = 0

@dataclass
class WorkflowAction:
    """Represents a workflow action"""
    action_id: str
    name: str
    type: str  # assign, escalate, notify, create_subtask, update_status
    parameters: Dict[str, Any]
    conditions: List[Dict[str, Any]] = None

@dataclass
class TicketTemplate:
    """Represents a ticket template"""
    template_id: str
    name: str
    category: str
    title_template: str
    description_template: str
    default_priority: str
    default_assignee: str
    tags: List[str] = None
    metadata: Dict[str, Any] = None

class TicketingConnector(ABC):
    """Abstract base class for ticketing system connectors"""
    
    @abstractmethod
    async def connect(self, config: Dict[str, Any]) -> bool:
        """Connect to ticketing system"""
        pass
    
    @abstractmethod
    async def create_ticket(self, ticket: Ticket) -> str:
        """Create ticket in system"""
        pass
    
    @abstractmethod
    async def update_ticket(self, ticket_id: str, updates: Dict[str, Any]) -> bool:
        """Update ticket in system"""
        pass
    
    @abstractmethod
    async def get_ticket(self, ticket_id: str) -> Optional[Ticket]:
        """Get ticket from system"""
        pass
    
    @abstractmethod
    async def search_tickets(self, filters: Dict[str, Any]) -> List[Ticket]:
        """Search tickets in system"""
        pass

class JiraConnector(TicketingConnector):
    """Jira ticketing system connector"""
    
    def __init__(self):
        self.session = None
        self.base_url = None
        self.auth_token = None
    
    async def connect(self, config: Dict[str, Any]) -> bool:
        """Connect to Jira"""
        try:
            self.base_url = config['endpoint']
            self.auth_token = config['api_token']
            self.session = requests.Session()
            
            # Test connection
            url = f"{self.base_url}/rest/api/2/myself"
            headers = {
                'Authorization': f'Basic {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            response = self.session.get(url, headers=headers)
            if response.status_code == 200:
                logger.info("Successfully connected to Jira")
                return True
            else:
                logger.error(f"Failed to connect to Jira: {response.status_code}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to Jira: {e}")
            return False
    
    async def create_ticket(self, ticket: Ticket) -> str:
        """Create ticket in Jira"""
        try:
            url = f"{self.base_url}/rest/api/2/issue"
            headers = {
                'Authorization': f'Basic {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            # Map priority to Jira priority
            priority_map = {
                'critical': 'Highest',
                'high': 'High',
                'medium': 'Medium',
                'low': 'Low'
            }
            
            jira_ticket = {
                'fields': {
                    'project': {'key': 'SEC'},  # Security project
                    'issuetype': {'name': 'Task'},
                    'summary': ticket.title,
                    'description': ticket.description,
                    'priority': {'name': priority_map.get(ticket.priority, 'Medium')},
                    'labels': ticket.tags or [],
                    'assignee': {'name': ticket.assignee} if ticket.assignee else None
                }
            }
            
            response = self.session.post(url, headers=headers, json=jira_ticket)
            if response.status_code == 201:
                ticket_id = response.json().get('key')
                logger.info(f"Created Jira ticket: {ticket_id}")
                return ticket_id
            else:
                logger.error(f"Failed to create Jira ticket: {response.status_code}")
                return None
                
        except Exception as e:
            logger.error(f"Error creating Jira ticket: {e}")
            return None
    
    async def update_ticket(self, ticket_id: str, updates: Dict[str, Any]) -> bool:
        """Update ticket in Jira"""
        try:
            url = f"{self.base_url}/rest/api/2/issue/{ticket_id}"
            headers = {
                'Authorization': f'Basic {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            # Map updates to Jira format
            jira_updates = {}
            if 'status' in updates:
                jira_updates['status'] = {'name': updates['status']}
            if 'assignee' in updates:
                jira_updates['assignee'] = {'name': updates['assignee']}
            if 'priority' in updates:
                priority_map = {
                    'critical': 'Highest',
                    'high': 'High',
                    'medium': 'Medium',
                    'low': 'Low'
                }
                jira_updates['priority'] = {'name': priority_map.get(updates['priority'], 'Medium')}
            
            data = {'fields': jira_updates}
            response = self.session.put(url, headers=headers, json=data)
            return response.status_code == 204
            
        except Exception as e:
            logger.error(f"Error updating Jira ticket: {e}")
            return False
    
    async def get_ticket(self, ticket_id: str) -> Optional[Ticket]:
        """Get ticket from Jira"""
        try:
            url = f"{self.base_url}/rest/api/2/issue/{ticket_id}"
            headers = {
                'Authorization': f'Basic {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            response = self.session.get(url, headers=headers)
            if response.status_code == 200:
                data = response.json()
                fields = data.get('fields', {})
                
                # Map Jira priority back to our format
                priority_map = {
                    'Highest': 'critical',
                    'High': 'high',
                    'Medium': 'medium',
                    'Low': 'low'
                }
                
                return Ticket(
                    ticket_id=ticket_id,
                    title=fields.get('summary', ''),
                    description=fields.get('description', ''),
                    priority=priority_map.get(fields.get('priority', {}).get('name', ''), 'medium'),
                    status=fields.get('status', {}).get('name', 'open'),
                    assignee=fields.get('assignee', {}).get('displayName') if fields.get('assignee') else None,
                    created_at=fields.get('created', ''),
                    updated_at=fields.get('updated', ''),
                    tags=fields.get('labels', [])
                )
            else:
                logger.error(f"Failed to get Jira ticket: {response.status_code}")
                return None
                
        except Exception as e:
            logger.error(f"Error getting Jira ticket: {e}")
            return None
    
    async def search_tickets(self, filters: Dict[str, Any]) -> List[Ticket]:
        """Search tickets in Jira"""
        try:
            url = f"{self.base_url}/rest/api/2/search"
            headers = {
                'Authorization': f'Basic {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            # Build JQL query
            jql_parts = []
            if 'status' in filters:
                jql_parts.append(f"status = '{filters['status']}'")
            if 'priority' in filters:
                jql_parts.append(f"priority = '{filters['priority']}'")
            if 'assignee' in filters:
                jql_parts.append(f"assignee = '{filters['assignee']}'")
            
            jql = " AND ".join(jql_parts) if jql_parts else "project = SEC"
            
            data = {
                'jql': jql,
                'maxResults': filters.get('limit', 50)
            }
            
            response = self.session.post(url, headers=headers, json=data)
            if response.status_code == 200:
                tickets = []
                for issue in response.json().get('issues', []):
                    fields = issue.get('fields', {})
                    
                    priority_map = {
                        'Highest': 'critical',
                        'High': 'high',
                        'Medium': 'medium',
                        'Low': 'low'
                    }
                    
                    tickets.append(Ticket(
                        ticket_id=issue.get('key', ''),
                        title=fields.get('summary', ''),
                        description=fields.get('description', ''),
                        priority=priority_map.get(fields.get('priority', {}).get('name', ''), 'medium'),
                        status=fields.get('status', {}).get('name', 'open'),
                        assignee=fields.get('assignee', {}).get('displayName') if fields.get('assignee') else None,
                        created_at=fields.get('created', ''),
                        updated_at=fields.get('updated', ''),
                        tags=fields.get('labels', [])
                    ))
                
                return tickets
            else:
                logger.error(f"Failed to search Jira tickets: {response.status_code}")
                return []
                
        except Exception as e:
            logger.error(f"Error searching Jira tickets: {e}")
            return []

class ServiceNowConnector(TicketingConnector):
    """ServiceNow ticketing system connector"""
    
    def __init__(self):
        self.session = None
        self.base_url = None
        self.auth_token = None
    
    async def connect(self, config: Dict[str, Any]) -> bool:
        """Connect to ServiceNow"""
        try:
            self.base_url = config['endpoint']
            self.auth_token = config['api_token']
            self.session = requests.Session()
            
            # Test connection
            url = f"{self.base_url}/api/now/table/sys_user"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            response = self.session.get(url, headers=headers)
            if response.status_code == 200:
                logger.info("Successfully connected to ServiceNow")
                return True
            else:
                logger.error(f"Failed to connect to ServiceNow: {response.status_code}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to ServiceNow: {e}")
            return False
    
    async def create_ticket(self, ticket: Ticket) -> str:
        """Create ticket in ServiceNow"""
        try:
            url = f"{self.base_url}/api/now/table/incident"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            servicenow_ticket = {
                'short_description': ticket.title,
                'description': ticket.description,
                'priority': ticket.priority,
                'category': ticket.category,
                'assigned_to': ticket.assignee,
                'state': '1' if ticket.status == 'open' else '2'  # 1=New, 2=In Progress
            }
            
            response = self.session.post(url, headers=headers, json=servicenow_ticket)
            if response.status_code == 201:
                ticket_id = response.json().get('result', {}).get('sys_id')
                logger.info(f"Created ServiceNow ticket: {ticket_id}")
                return ticket_id
            else:
                logger.error(f"Failed to create ServiceNow ticket: {response.status_code}")
                return None
                
        except Exception as e:
            logger.error(f"Error creating ServiceNow ticket: {e}")
            return None
    
    async def update_ticket(self, ticket_id: str, updates: Dict[str, Any]) -> bool:
        """Update ticket in ServiceNow"""
        try:
            url = f"{self.base_url}/api/now/table/incident/{ticket_id}"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            response = self.session.put(url, headers=headers, json=updates)
            return response.status_code == 200
            
        except Exception as e:
            logger.error(f"Error updating ServiceNow ticket: {e}")
            return False
    
    async def get_ticket(self, ticket_id: str) -> Optional[Ticket]:
        """Get ticket from ServiceNow"""
        try:
            url = f"{self.base_url}/api/now/table/incident/{ticket_id}"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            response = self.session.get(url, headers=headers)
            if response.status_code == 200:
                data = response.json().get('result', {})
                
                return Ticket(
                    ticket_id=ticket_id,
                    title=data.get('short_description', ''),
                    description=data.get('description', ''),
                    priority=data.get('priority', 'medium'),
                    status='open' if data.get('state') == '1' else 'in_progress',
                    assignee=data.get('assigned_to'),
                    created_at=data.get('sys_created_on', ''),
                    updated_at=data.get('sys_updated_on', '')
                )
            else:
                logger.error(f"Failed to get ServiceNow ticket: {response.status_code}")
                return None
                
        except Exception as e:
            logger.error(f"Error getting ServiceNow ticket: {e}")
            return None
    
    async def search_tickets(self, filters: Dict[str, Any]) -> List[Ticket]:
        """Search tickets in ServiceNow"""
        try:
            url = f"{self.base_url}/api/now/table/incident"
            headers = {
                'Authorization': f'Bearer {self.auth_token}',
                'Content-Type': 'application/json'
            }
            
            # Build query parameters
            params = {}
            if 'status' in filters:
                params['state'] = '1' if filters['status'] == 'open' else '2'
            if 'priority' in filters:
                params['priority'] = filters['priority']
            if 'assignee' in filters:
                params['assigned_to'] = filters['assignee']
            
            response = self.session.get(url, headers=headers, params=params)
            if response.status_code == 200:
                tickets = []
                for data in response.json().get('result', []):
                    tickets.append(Ticket(
                        ticket_id=data.get('sys_id', ''),
                        title=data.get('short_description', ''),
                        description=data.get('description', ''),
                        priority=data.get('priority', 'medium'),
                        status='open' if data.get('state') == '1' else 'in_progress',
                        assignee=data.get('assigned_to'),
                        created_at=data.get('sys_created_on', ''),
                        updated_at=data.get('sys_updated_on', '')
                    ))
                
                return tickets
            else:
                logger.error(f"Failed to search ServiceNow tickets: {response.status_code}")
                return []
                
        except Exception as e:
            logger.error(f"Error searching ServiceNow tickets: {e}")
            return []

class IntelligentTicketRouter:
    """Intelligent ticket routing and assignment system"""
    
    def __init__(self):
        self.vectorizer = TfidfVectorizer(max_features=1000, stop_words='english')
        self.team_expertise = {}
        self.routing_rules = []
        self.ticket_templates = {}
        
        # Initialize routing system
        self._initialize_routing_system()
    
    def _initialize_routing_system(self):
        """Initialize the routing system"""
        try:
            # Define team expertise areas
            self.team_expertise = {
                'security_team': {
                    'expertise': ['vulnerability', 'security', 'threat', 'incident', 'malware', 'firewall'],
                    'capacity': 10,
                    'current_load': 0
                },
                'network_team': {
                    'expertise': ['network', 'routing', 'switching', 'firewall', 'vpn', 'dns'],
                    'capacity': 8,
                    'current_load': 0
                },
                'application_team': {
                    'expertise': ['application', 'web', 'api', 'database', 'code', 'development'],
                    'capacity': 12,
                    'current_load': 0
                },
                'infrastructure_team': {
                    'expertise': ['server', 'cloud', 'infrastructure', 'monitoring', 'backup', 'storage'],
                    'capacity': 6,
                    'current_load': 0
                }
            }
            
            # Define routing rules
            self.routing_rules = [
                {
                    'name': 'Critical Security Issues',
                    'conditions': [
                        {'field': 'priority', 'operator': 'equals', 'value': 'critical'},
                        {'field': 'category', 'operator': 'equals', 'value': 'security'}
                    ],
                    'actions': [
                        {'type': 'assign', 'target': 'security_team'},
                        {'type': 'escalate', 'level': 'immediate'}
                    ]
                },
                {
                    'name': 'Network Issues',
                    'conditions': [
                        {'field': 'description', 'operator': 'contains', 'value': 'network'}
                    ],
                    'actions': [
                        {'type': 'assign', 'target': 'network_team'}
                    ]
                },
                {
                    'name': 'Application Issues',
                    'conditions': [
                        {'field': 'description', 'operator': 'contains', 'value': 'application'}
                    ],
                    'actions': [
                        {'type': 'assign', 'target': 'application_team'}
                    ]
                }
            ]
            
            logger.info("Intelligent ticket routing system initialized")
            
        except Exception as e:
            logger.error(f"Error initializing routing system: {e}")
    
    async def route_ticket(self, ticket: Ticket) -> str:
        """Route ticket to appropriate team"""
        try:
            # Check routing rules
            for rule in self.routing_rules:
                if await self._evaluate_rule_conditions(rule['conditions'], ticket):
                    return await self._execute_rule_actions(rule['actions'], ticket)
            
            # Use ML-based routing as fallback
            return await self._ml_based_routing(ticket)
            
        except Exception as e:
            logger.error(f"Error routing ticket: {e}")
            return 'security_team'  # Default fallback
    
    async def _evaluate_rule_conditions(self, conditions: List[Dict[str, Any]], ticket: Ticket) -> bool:
        """Evaluate rule conditions"""
        try:
            for condition in conditions:
                field = condition['field']
                operator = condition['operator']
                value = condition['value']
                
                ticket_value = getattr(ticket, field, '')
                
                if operator == 'equals':
                    if ticket_value != value:
                        return False
                elif operator == 'contains':
                    if value.lower() not in str(ticket_value).lower():
                        return False
                elif operator == 'greater_than':
                    if not (isinstance(ticket_value, (int, float)) and ticket_value > value):
                        return False
                elif operator == 'less_than':
                    if not (isinstance(ticket_value, (int, float)) and ticket_value < value):
                        return False
            
            return True
            
        except Exception as e:
            logger.error(f"Error evaluating rule conditions: {e}")
            return False
    
    async def _execute_rule_actions(self, actions: List[Dict[str, Any]], ticket: Ticket) -> str:
        """Execute rule actions"""
        try:
            for action in actions:
                if action['type'] == 'assign':
                    return action['target']
                elif action['type'] == 'escalate':
                    # Handle escalation logic
                    pass
            
            return 'security_team'  # Default fallback
            
        except Exception as e:
            logger.error(f"Error executing rule actions: {e}")
            return 'security_team'
    
    async def _ml_based_routing(self, ticket: Ticket) -> str:
        """ML-based ticket routing"""
        try:
            # Prepare text for analysis
            text = f"{ticket.title} {ticket.description}"
            
            # Calculate similarity with team expertise
            similarities = {}
            for team, info in self.team_expertise.items():
                expertise_text = ' '.join(info['expertise'])
                similarity = self._calculate_text_similarity(text, expertise_text)
                similarities[team] = similarity
            
            # Consider team capacity
            for team in similarities:
                capacity_ratio = self.team_expertise[team]['current_load'] / self.team_expertise[team]['capacity']
                similarities[team] *= (1 - capacity_ratio)  # Reduce score for overloaded teams
            
            # Return team with highest similarity
            best_team = max(similarities, key=similarities.get)
            return best_team
            
        except Exception as e:
            logger.error(f"Error in ML-based routing: {e}")
            return 'security_team'
    
    def _calculate_text_similarity(self, text1: str, text2: str) -> float:
        """Calculate text similarity using TF-IDF and cosine similarity"""
        try:
            # Combine texts
            texts = [text1, text2]
            
            # Fit and transform
            tfidf_matrix = self.vectorizer.fit_transform(texts)
            
            # Calculate cosine similarity
            similarity = cosine_similarity(tfidf_matrix[0:1], tfidf_matrix[1:2])[0][0]
            
            return similarity
            
        except Exception as e:
            logger.error(f"Error calculating text similarity: {e}")
            return 0.0
    
    async def create_ticket_from_template(self, template_id: str, context: Dict[str, Any]) -> Ticket:
        """Create ticket from template"""
        try:
            template = self.ticket_templates.get(template_id)
            if not template:
                raise ValueError(f"Template {template_id} not found")
            
            # Replace template variables
            title = template.title_template.format(**context)
            description = template.description_template.format(**context)
            
            # Create ticket
            ticket = Ticket(
                ticket_id=f"TEMP_{int(time.time())}",
                title=title,
                description=description,
                priority=template.default_priority,
                assignee=template.default_assignee,
                category=template.category,
                tags=template.tags or [],
                created_at=datetime.utcnow().isoformat()
            )
            
            return ticket
            
        except Exception as e:
            logger.error(f"Error creating ticket from template: {e}")
            return None
    
    async def update_team_load(self, team: str, change: int):
        """Update team load"""
        try:
            if team in self.team_expertise:
                self.team_expertise[team]['current_load'] += change
                self.team_expertise[team]['current_load'] = max(0, self.team_expertise[team]['current_load'])
                
        except Exception as e:
            logger.error(f"Error updating team load: {e}")

class AutomatedWorkflowEngine:
    """Automated workflow engine for ticket management"""
    
    def __init__(self):
        self.workflow_rules = []
        self.active_workflows = {}
        self.ticketing_connectors = {}
        
        # Initialize workflow engine
        self._initialize_workflow_engine()
    
    def _initialize_workflow_engine(self):
        """Initialize workflow engine"""
        try:
            # Define default workflow rules
            self.workflow_rules = [
                {
                    'rule_id': 'auto_assign_critical',
                    'name': 'Auto-assign Critical Tickets',
                    'conditions': [
                        {'field': 'priority', 'operator': 'equals', 'value': 'critical'}
                    ],
                    'actions': [
                        {'type': 'assign', 'target': 'security_team'},
                        {'type': 'notify', 'channels': ['email', 'slack']},
                        {'type': 'escalate', 'level': 'immediate'}
                    ],
                    'priority': 1,
                    'enabled': True
                },
                {
                    'rule_id': 'auto_close_resolved',
                    'name': 'Auto-close Resolved Tickets',
                    'conditions': [
                        {'field': 'status', 'operator': 'equals', 'value': 'resolved'},
                        {'field': 'updated_at', 'operator': 'older_than', 'value': '7 days'}
                    ],
                    'actions': [
                        {'type': 'update_status', 'status': 'closed'},
                        {'type': 'notify', 'channels': ['email']}
                    ],
                    'priority': 2,
                    'enabled': True
                },
                {
                    'rule_id': 'escalate_overdue',
                    'name': 'Escalate Overdue Tickets',
                    'conditions': [
                        {'field': 'due_date', 'operator': 'past_due'},
                        {'field': 'status', 'operator': 'not_equals', 'value': 'closed'}
                    ],
                    'actions': [
                        {'type': 'escalate', 'level': 'manager'},
                        {'type': 'notify', 'channels': ['email', 'slack']}
                    ],
                    'priority': 1,
                    'enabled': True
                }
            ]
            
            logger.info("Automated workflow engine initialized")
            
        except Exception as e:
            logger.error(f"Error initializing workflow engine: {e}")
    
    async def process_ticket(self, ticket: Ticket) -> bool:
        """Process ticket through workflow rules"""
        try:
            # Sort rules by priority
            sorted_rules = sorted(self.workflow_rules, key=lambda x: x['priority'])
            
            for rule in sorted_rules:
                if not rule['enabled']:
                    continue
                
                if await self._evaluate_workflow_conditions(rule['conditions'], ticket):
                    await self._execute_workflow_actions(rule['actions'], ticket)
                    rule['last_triggered'] = datetime.utcnow().isoformat()
                    rule['trigger_count'] += 1
            
            return True
            
        except Exception as e:
            logger.error(f"Error processing ticket: {e}")
            return False
    
    async def _evaluate_workflow_conditions(self, conditions: List[Dict[str, Any]], ticket: Ticket) -> bool:
        """Evaluate workflow conditions"""
        try:
            for condition in conditions:
                field = condition['field']
                operator = condition['operator']
                value = condition['value']
                
                ticket_value = getattr(ticket, field, '')
                
                if operator == 'equals':
                    if ticket_value != value:
                        return False
                elif operator == 'not_equals':
                    if ticket_value == value:
                        return False
                elif operator == 'contains':
                    if value.lower() not in str(ticket_value).lower():
                        return False
                elif operator == 'older_than':
                    if not self._is_older_than(ticket_value, value):
                        return False
                elif operator == 'past_due':
                    if not self._is_past_due(ticket_value):
                        return False
            
            return True
            
        except Exception as e:
            logger.error(f"Error evaluating workflow conditions: {e}")
            return False
    
    async def _execute_workflow_actions(self, actions: List[Dict[str, Any]], ticket: Ticket):
        """Execute workflow actions"""
        try:
            for action in actions:
                if action['type'] == 'assign':
                    ticket.assignee = action['target']
                elif action['type'] == 'update_status':
                    ticket.status = action['status']
                elif action['type'] == 'notify':
                    await self._send_notification(action['channels'], ticket)
                elif action['type'] == 'escalate':
                    await self._escalate_ticket(action['level'], ticket)
            
        except Exception as e:
            logger.error(f"Error executing workflow actions: {e}")
    
    def _is_older_than(self, timestamp: str, days: str) -> bool:
        """Check if timestamp is older than specified days"""
        try:
            if not timestamp:
                return False
            
            ticket_time = datetime.fromisoformat(timestamp.replace('Z', '+00:00'))
            cutoff_time = datetime.utcnow() - timedelta(days=int(days.split()[0]))
            
            return ticket_time < cutoff_time
            
        except Exception as e:
            logger.error(f"Error checking if older than: {e}")
            return False
    
    def _is_past_due(self, due_date: str) -> bool:
        """Check if due date is past due"""
        try:
            if not due_date:
                return False
            
            due_time = datetime.fromisoformat(due_date.replace('Z', '+00:00'))
            return due_time < datetime.utcnow()
            
        except Exception as e:
            logger.error(f"Error checking if past due: {e}")
            return False
    
    async def _send_notification(self, channels: List[str], ticket: Ticket):
        """Send notification through specified channels"""
        try:
            message = f"Ticket {ticket.ticket_id}: {ticket.title} - Status: {ticket.status}"
            
            for channel in channels:
                if channel == 'email':
                    # Send email notification
                    logger.info(f"Email notification: {message}")
                elif channel == 'slack':
                    # Send Slack notification
                    logger.info(f"Slack notification: {message}")
                elif channel == 'teams':
                    # Send Teams notification
                    logger.info(f"Teams notification: {message}")
            
        except Exception as e:
            logger.error(f"Error sending notification: {e}")
    
    async def _escalate_ticket(self, level: str, ticket: Ticket):
        """Escalate ticket to specified level"""
        try:
            if level == 'immediate':
                ticket.priority = 'critical'
                ticket.assignee = 'security_manager'
            elif level == 'manager':
                ticket.assignee = 'team_manager'
            
            logger.info(f"Escalated ticket {ticket.ticket_id} to {level}")
            
        except Exception as e:
            logger.error(f"Error escalating ticket: {e}")

# Global ticketing and workflow system
ticketing_workflow_system = {
    'connectors': {
        'jira': JiraConnector(),
        'servicenow': ServiceNowConnector()
    },
    'router': IntelligentTicketRouter(),
    'workflow_engine': AutomatedWorkflowEngine()
}

