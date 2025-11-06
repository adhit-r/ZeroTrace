"""
ChatOps Integration - Slack/Teams integration with security chatbot
"""

import asyncio
import json
import logging
import requests
import time
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Any, Tuple
from dataclasses import dataclass
import re
from slack_sdk import WebClient
from slack_sdk.errors import SlackApiError
import openai
from msal import ConfidentialClientApplication

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

@dataclass
class ChatMessage:
    """Represents a chat message"""
    message_id: str
    channel: str
    user: str
    text: str
    timestamp: str
    thread_ts: Optional[str] = None
    bot_response: Optional[str] = None
    metadata: Dict[str, Any] = None

@dataclass
class SecurityAlert:
    """Represents a security alert for chat"""
    alert_id: str
    title: str
    description: str
    severity: str
    source: str
    timestamp: str
    actions: List[str] = None
    status: str = "active"

@dataclass
class ChatCommand:
    """Represents a chat command"""
    command: str
    description: str
    usage: str
    handler: str
    permissions: List[str] = None
    examples: List[str] = None

class SlackConnector:
    """Slack integration connector"""
    
    def __init__(self):
        self.client = None
        self.bot_token = None
        self.app_token = None
        self.workspace_id = None
    
    async def connect(self, config: Dict[str, Any]) -> bool:
        """Connect to Slack"""
        try:
            self.bot_token = config['bot_token']
            self.app_token = config.get('app_token')
            self.workspace_id = config.get('workspace_id')
            
            self.client = WebClient(token=self.bot_token)
            
            # Test connection
            response = self.client.auth_test()
            if response['ok']:
                logger.info("Successfully connected to Slack")
                return True
            else:
                logger.error(f"Failed to connect to Slack: {response['error']}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to Slack: {e}")
            return False
    
    async def send_message(self, channel: str, text: str, thread_ts: Optional[str] = None) -> bool:
        """Send message to Slack channel"""
        try:
            response = self.client.chat_postMessage(
                channel=channel,
                text=text,
                thread_ts=thread_ts
            )
            
            if response['ok']:
                logger.info(f"Message sent to {channel}")
                return True
            else:
                logger.error(f"Failed to send message: {response['error']}")
                return False
                
        except SlackApiError as e:
            logger.error(f"Error sending message to Slack: {e}")
            return False
    
    async def send_alert(self, channel: str, alert: SecurityAlert) -> bool:
        """Send security alert to Slack"""
        try:
            # Create rich message with blocks
            blocks = [
                {
                    "type": "header",
                    "text": {
                        "type": "plain_text",
                        "text": f"🚨 {alert.title}"
                    }
                },
                {
                    "type": "section",
                    "fields": [
                        {
                            "type": "mrkdwn",
                            "text": f"*Severity:* {alert.severity.upper()}"
                        },
                        {
                            "type": "mrkdwn",
                            "text": f"*Source:* {alert.source}"
                        },
                        {
                            "type": "mrkdwn",
                            "text": f"*Time:* {alert.timestamp}"
                        },
                        {
                            "type": "mrkdwn",
                            "text": f"*Status:* {alert.status}"
                        }
                    ]
                },
                {
                    "type": "section",
                    "text": {
                        "type": "mrkdwn",
                        "text": f"*Description:*\n{alert.description}"
                    }
                }
            ]
            
            # Add action buttons if available
            if alert.actions:
                actions = []
                for action in alert.actions:
                    actions.append({
                        "type": "button",
                        "text": {
                            "type": "plain_text",
                            "text": action
                        },
                        "action_id": f"action_{action.lower().replace(' ', '_')}"
                    })
                
                blocks.append({
                    "type": "actions",
                    "elements": actions
                })
            
            response = self.client.chat_postMessage(
                channel=channel,
                blocks=blocks
            )
            
            return response['ok']
            
        except SlackApiError as e:
            logger.error(f"Error sending alert to Slack: {e}")
            return False
    
    async def get_channel_messages(self, channel: str, limit: int = 100) -> List[ChatMessage]:
        """Get messages from Slack channel"""
        try:
            response = self.client.conversations_history(
                channel=channel,
                limit=limit
            )
            
            messages = []
            for msg in response['messages']:
                messages.append(ChatMessage(
                    message_id=msg['ts'],
                    channel=channel,
                    user=msg.get('user', ''),
                    text=msg.get('text', ''),
                    timestamp=msg.get('ts', ''),
                    thread_ts=msg.get('thread_ts'),
                    metadata=msg
                ))
            
            return messages
            
        except SlackApiError as e:
            logger.error(f"Error getting messages from Slack: {e}")
            return []
    
    async def create_channel(self, name: str, is_private: bool = False) -> str:
        """Create Slack channel"""
        try:
            response = self.client.conversations_create(
                name=name,
                is_private=is_private
            )
            
            if response['ok']:
                return response['channel']['id']
            else:
                logger.error(f"Failed to create channel: {response['error']}")
                return None
                
        except SlackApiError as e:
            logger.error(f"Error creating Slack channel: {e}")
            return None

class TeamsConnector:
    """Microsoft Teams integration connector"""
    
    def __init__(self):
        self.app = None
        self.access_token = None
        self.tenant_id = None
        self.client_id = None
        self.client_secret = None
    
    async def connect(self, config: Dict[str, Any]) -> bool:
        """Connect to Microsoft Teams"""
        try:
            self.tenant_id = config['tenant_id']
            self.client_id = config['client_id']
            self.client_secret = config['client_secret']
            
            # Initialize MSAL app
            self.app = ConfidentialClientApplication(
                client_id=self.client_id,
                client_credential=self.client_secret,
                authority=f"https://login.microsoftonline.com/{self.tenant_id}"
            )
            
            # Get access token
            scopes = ["https://graph.microsoft.com/.default"]
            result = self.app.acquire_token_silent(scopes, account=None)
            
            if not result:
                result = self.app.acquire_token_for_client(scopes=scopes)
            
            if "access_token" in result:
                self.access_token = result["access_token"]
                logger.info("Successfully connected to Microsoft Teams")
                return True
            else:
                logger.error(f"Failed to get access token: {result.get('error_description')}")
                return False
                
        except Exception as e:
            logger.error(f"Error connecting to Teams: {e}")
            return False
    
    async def send_message(self, channel_id: str, text: str) -> bool:
        """Send message to Teams channel"""
        try:
            url = f"https://graph.microsoft.com/v1.0/teams/{channel_id}/channels/{channel_id}/messages"
            headers = {
                'Authorization': f'Bearer {self.access_token}',
                'Content-Type': 'application/json'
            }
            
            data = {
                'body': {
                    'content': text
                }
            }
            
            response = requests.post(url, headers=headers, json=data)
            return response.status_code == 201
            
        except Exception as e:
            logger.error(f"Error sending message to Teams: {e}")
            return False
    
    async def send_alert(self, channel_id: str, alert: SecurityAlert) -> bool:
        """Send security alert to Teams"""
        try:
            # Create adaptive card for Teams
            card = {
                "type": "AdaptiveCard",
                "version": "1.0",
                "body": [
                    {
                        "type": "TextBlock",
                        "text": f"🚨 {alert.title}",
                        "weight": "Bolder",
                        "size": "Large"
                    },
                    {
                        "type": "FactSet",
                        "facts": [
                            {"title": "Severity", "value": alert.severity.upper()},
                            {"title": "Source", "value": alert.source},
                            {"title": "Time", "value": alert.timestamp},
                            {"title": "Status", "value": alert.status}
                        ]
                    },
                    {
                        "type": "TextBlock",
                        "text": f"**Description:**\n{alert.description}",
                        "wrap": True
                    }
                ]
            }
            
            # Add action buttons if available
            if alert.actions:
                actions = []
                for action in alert.actions:
                    actions.append({
                        "type": "Action.Submit",
                        "title": action,
                        "data": {
                            "action": action,
                            "alert_id": alert.alert_id
                        }
                    })
                
                card["actions"] = actions
            
            url = f"https://graph.microsoft.com/v1.0/teams/{channel_id}/channels/{channel_id}/messages"
            headers = {
                'Authorization': f'Bearer {self.access_token}',
                'Content-Type': 'application/json'
            }
            
            data = {
                'body': {
                    'contentType': 'html',
                    'content': f'<div>{json.dumps(card)}</div>'
                }
            }
            
            response = requests.post(url, headers=headers, json=data)
            return response.status_code == 201
            
        except Exception as e:
            logger.error(f"Error sending alert to Teams: {e}")
            return False

class SecurityChatbot:
    """AI-powered security chatbot"""
    
    def __init__(self):
        self.openai_client = None
        self.commands = {}
        self.conversation_history = {}
        
        # Initialize chatbot
        self._initialize_chatbot()
    
    def _initialize_chatbot(self):
        """Initialize the security chatbot"""
        try:
            # Initialize OpenAI client
            self.openai_client = openai.OpenAI(api_key="your-openai-api-key")
            
            # Define security commands
            self.commands = {
                'vulnerability_scan': ChatCommand(
                    command='scan',
                    description='Start vulnerability scan',
                    usage='scan [target] [options]',
                    handler='handle_scan_command',
                    permissions=['security_team'],
                    examples=['scan server1', 'scan --full network']
                ),
                'security_status': ChatCommand(
                    command='status',
                    description='Get security status',
                    usage='status [component]',
                    handler='handle_status_command',
                    permissions=['security_team', 'management'],
                    examples=['status', 'status firewall']
                ),
                'incident_report': ChatCommand(
                    command='incident',
                    description='Report security incident',
                    usage='incident [description]',
                    handler='handle_incident_command',
                    permissions=['security_team'],
                    examples=['incident suspicious activity detected']
                ),
                'compliance_check': ChatCommand(
                    command='compliance',
                    description='Check compliance status',
                    usage='compliance [framework]',
                    handler='handle_compliance_command',
                    permissions=['security_team', 'compliance_team'],
                    examples=['compliance SOC2', 'compliance PCI']
                ),
                'help': ChatCommand(
                    command='help',
                    description='Show available commands',
                    usage='help [command]',
                    handler='handle_help_command',
                    permissions=['all'],
                    examples=['help', 'help scan']
                )
            }
            
            logger.info("Security chatbot initialized")
            
        except Exception as e:
            logger.error(f"Error initializing chatbot: {e}")
    
    async def process_message(self, message: ChatMessage) -> str:
        """Process incoming chat message"""
        try:
            # Extract command from message
            command, args = self._parse_command(message.text)
            
            if command in self.commands:
                # Check permissions
                if not self._check_permissions(command, message.user):
                    return "❌ You don't have permission to use this command."
                
                # Execute command
                handler = getattr(self, self.commands[command].handler)
                response = await handler(args, message)
                return response
            else:
                # Use AI to generate response
                return await self._generate_ai_response(message)
            
        except Exception as e:
            logger.error(f"Error processing message: {e}")
            return "❌ Sorry, I encountered an error processing your message."
    
    def _parse_command(self, text: str) -> Tuple[str, List[str]]:
        """Parse command from text"""
        try:
            # Remove bot mention if present
            text = re.sub(r'<@[^>]+>', '', text).strip()
            
            # Split into command and arguments
            parts = text.split()
            if not parts:
                return '', []
            
            command = parts[0].lower()
            args = parts[1:] if len(parts) > 1 else []
            
            return command, args
            
        except Exception as e:
            logger.error(f"Error parsing command: {e}")
            return '', []
    
    def _check_permissions(self, command: str, user: str) -> bool:
        """Check if user has permission to use command"""
        try:
            command_info = self.commands.get(command)
            if not command_info:
                return False
            
            # Mock permission check - in real implementation, this would check user roles
            user_roles = self._get_user_roles(user)
            required_permissions = command_info.permissions
            
            if 'all' in required_permissions:
                return True
            
            return any(role in user_roles for role in required_permissions)
            
        except Exception as e:
            logger.error(f"Error checking permissions: {e}")
            return False
    
    def _get_user_roles(self, user: str) -> List[str]:
        """Get user roles - mock implementation"""
        # Mock user roles - in real implementation, this would query user database
        role_map = {
            'U1234567890': ['security_team', 'management'],
            'U0987654321': ['security_team'],
            'U1122334455': ['compliance_team']
        }
        
        return role_map.get(user, ['user'])
    
    async def handle_scan_command(self, args: List[str], message: ChatMessage) -> str:
        """Handle vulnerability scan command"""
        try:
            target = args[0] if args else 'all'
            
            # Mock scan execution
            response = f"🔍 Starting vulnerability scan on {target}...\n"
            response += "⏳ Scan in progress...\n"
            response += "📊 Results will be available in 5-10 minutes.\n"
            response += "🔔 I'll notify you when the scan is complete."
            
            return response
            
        except Exception as e:
            logger.error(f"Error handling scan command: {e}")
            return "❌ Error executing scan command."
    
    async def handle_status_command(self, args: List[str], message: ChatMessage) -> str:
        """Handle security status command"""
        try:
            component = args[0] if args else 'overall'
            
            # Mock status response
            if component == 'overall':
                response = "🛡️ **Security Status Overview**\n"
                response += "• Firewall: ✅ Active\n"
                response += "• Antivirus: ✅ Updated\n"
                response += "• Vulnerability Scanner: ✅ Running\n"
                response += "• Last Scan: 2 hours ago\n"
                response += "• Critical Issues: 0\n"
                response += "• High Issues: 2\n"
                response += "• Medium Issues: 5"
            else:
                response = f"📊 **{component.title()} Status**\n"
                response += f"• Status: ✅ Operational\n"
                response += f"• Last Check: 1 hour ago\n"
                response += f"• Issues: 0"
            
            return response
            
        except Exception as e:
            logger.error(f"Error handling status command: {e}")
            return "❌ Error getting security status."
    
    async def handle_incident_command(self, args: List[str], message: ChatMessage) -> str:
        """Handle incident report command"""
        try:
            description = ' '.join(args) if args else 'No description provided'
            
            # Mock incident creation
            incident_id = f"INC-{int(time.time())}"
            
            response = f"🚨 **Security Incident Created**\n"
            response += f"• Incident ID: {incident_id}\n"
            response += f"• Description: {description}\n"
            response += f"• Status: Open\n"
            response += f"• Assigned to: Security Team\n"
            response += f"• Priority: High\n"
            response += f"• Created: {datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S')} UTC"
            
            return response
            
        except Exception as e:
            logger.error(f"Error handling incident command: {e}")
            return "❌ Error creating incident report."
    
    async def handle_compliance_command(self, args: List[str], message: ChatMessage) -> str:
        """Handle compliance check command"""
        try:
            framework = args[0] if args else 'SOC2'
            
            # Mock compliance status
            response = f"📋 **{framework} Compliance Status**\n"
            response += f"• Overall Score: 85%\n"
            response += f"• Status: ✅ Compliant\n"
            response += f"• Last Assessment: 1 week ago\n"
            response += f"• Next Assessment: 3 months\n"
            response += f"• Critical Findings: 0\n"
            response += f"• High Findings: 1\n"
            response += f"• Medium Findings: 3"
            
            return response
            
        except Exception as e:
            logger.error(f"Error handling compliance command: {e}")
            return "❌ Error checking compliance status."
    
    async def handle_help_command(self, args: List[str], message: ChatMessage) -> str:
        """Handle help command"""
        try:
            if args:
                command = args[0]
                if command in self.commands:
                    cmd_info = self.commands[command]
                    response = f"📖 **{command} Command Help**\n"
                    response += f"• Description: {cmd_info.description}\n"
                    response += f"• Usage: `{cmd_info.usage}`\n"
                    if cmd_info.examples:
                        response += f"• Examples:\n"
                        for example in cmd_info.examples:
                            response += f"  - `{example}`\n"
                else:
                    response = f"❌ Command '{command}' not found."
            else:
                response = "📖 **Available Commands**\n"
                for cmd_name, cmd_info in self.commands.items():
                    response += f"• `{cmd_name}` - {cmd_info.description}\n"
                response += "\nUse `help [command]` for detailed information."
            
            return response
            
        except Exception as e:
            logger.error(f"Error handling help command: {e}")
            return "❌ Error displaying help."
    
    async def _generate_ai_response(self, message: ChatMessage) -> str:
        """Generate AI-powered response"""
        try:
            # Get conversation history
            history = self.conversation_history.get(message.user, [])
            
            # Prepare context for AI
            context = "You are a security chatbot for ZeroTrace. Help users with security-related questions and tasks."
            if history:
                context += f"\n\nPrevious conversation:\n{chr(10).join(history[-5:])}"
            
            # Generate response using OpenAI
            response = self.openai_client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=[
                    {"role": "system", "content": context},
                    {"role": "user", "content": message.text}
                ],
                max_tokens=200,
                temperature=0.7
            )
            
            ai_response = response.choices[0].message.content
            
            # Update conversation history
            history.append(f"User: {message.text}")
            history.append(f"Bot: {ai_response}")
            self.conversation_history[message.user] = history[-10:]  # Keep last 10 exchanges
            
            return ai_response
            
        except Exception as e:
            logger.error(f"Error generating AI response: {e}")
            return "❌ Sorry, I'm having trouble generating a response right now."

class ChatOpsIntegrationHub:
    """Main ChatOps integration hub"""
    
    def __init__(self):
        self.slack_connector = SlackConnector()
        self.teams_connector = TeamsConnector()
        self.security_chatbot = SecurityChatbot()
        self.integration_configs = {}
        
        # Initialize integration hub
        self._initialize_integration_hub()
    
    def _initialize_integration_hub(self):
        """Initialize ChatOps integration hub"""
        try:
            logger.info("ChatOps integration hub initialized")
            
        except Exception as e:
            logger.error(f"Error initializing ChatOps hub: {e}")
    
    async def register_integration(self, config: Dict[str, Any]) -> bool:
        """Register ChatOps integration"""
        try:
            integration_type = config.get('type')
            integration_id = config.get('integration_id')
            
            if integration_type == 'slack':
                success = await self.slack_connector.connect(config)
                if success:
                    self.integration_configs[integration_id] = config
                    logger.info("Slack integration registered successfully")
                    return True
                else:
                    logger.error("Failed to register Slack integration")
                    return False
            
            elif integration_type == 'teams':
                success = await self.teams_connector.connect(config)
                if success:
                    self.integration_configs[integration_id] = config
                    logger.info("Teams integration registered successfully")
                    return True
                else:
                    logger.error("Failed to register Teams integration")
                    return False
            
            return False
            
        except Exception as e:
            logger.error(f"Error registering ChatOps integration: {e}")
            return False
    
    async def send_security_alert(self, alert: SecurityAlert, channels: List[str] = None) -> Dict[str, bool]:
        """Send security alert to chat channels"""
        try:
            results = {}
            
            for integration_id, config in self.integration_configs.items():
                if channels and config.get('channel') not in channels:
                    continue
                
                if config.get('type') == 'slack':
                    success = await self.slack_connector.send_alert(
                        config.get('channel'), alert
                    )
                    results[integration_id] = success
                
                elif config.get('type') == 'teams':
                    success = await self.teams_connector.send_alert(
                        config.get('channel_id'), alert
                    )
                    results[integration_id] = success
            
            return results
            
        except Exception as e:
            logger.error(f"Error sending security alert: {e}")
            return {}
    
    async def process_chat_message(self, message: ChatMessage, integration_type: str) -> str:
        """Process chat message through chatbot"""
        try:
            # Process through security chatbot
            response = await self.security_chatbot.process_message(message)
            
            # Send response back to chat
            if integration_type == 'slack':
                await self.slack_connector.send_message(
                    message.channel, response, message.thread_ts
                )
            elif integration_type == 'teams':
                await self.teams_connector.send_message(
                    message.channel, response
                )
            
            return response
            
        except Exception as e:
            logger.error(f"Error processing chat message: {e}")
            return "❌ Error processing your message."
    
    async def create_security_channel(self, name: str, integration_type: str) -> str:
        """Create security-focused chat channel"""
        try:
            if integration_type == 'slack':
                channel_id = await self.slack_connector.create_channel(name)
                return channel_id
            elif integration_type == 'teams':
                # Teams channel creation would require additional API calls
                logger.info(f"Teams channel creation not implemented yet")
                return None
            
            return None
            
        except Exception as e:
            logger.error(f"Error creating security channel: {e}")
            return None
    
    async def get_integration_status(self) -> Dict[str, Any]:
        """Get ChatOps integration status"""
        try:
            status = {
                'total_integrations': len(self.integration_configs),
                'slack_integrations': len([c for c in self.integration_configs.values() if c.get('type') == 'slack']),
                'teams_integrations': len([c for c in self.integration_configs.values() if c.get('type') == 'teams']),
                'active_commands': len(self.security_chatbot.commands),
                'integration_details': []
            }
            
            for integration_id, config in self.integration_configs.items():
                status['integration_details'].append({
                    'integration_id': integration_id,
                    'type': config.get('type'),
                    'channel': config.get('channel', config.get('channel_id')),
                    'status': 'active'
                })
            
            return status
            
        except Exception as e:
            logger.error(f"Error getting integration status: {e}")
            return {}

# Global ChatOps integration hub
chatops_hub = ChatOpsIntegrationHub()

