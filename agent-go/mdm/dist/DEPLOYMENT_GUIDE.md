# ZeroTrace Agent - MDM Deployment Guide

## 📦 Package Contents
- `ZeroTrace-Agent-1.0.0.pkg` - Main installation package
- `zerotrace-agent.mobileconfig` - MDM configuration profile

## 🚀 Deployment Steps

### Microsoft Intune
1. Upload `ZeroTrace-Agent-1.0.0.pkg` to Intune
2. Create macOS app with silent installation
3. Upload `zerotrace-agent.mobileconfig` as configuration profile
4. Assign to target device groups

### Jamf Pro
1. Upload `ZeroTrace-Agent-1.0.0.pkg` to Jamf
2. Create policy for package installation
3. Upload `zerotrace-agent.mobileconfig` as configuration profile
4. Configure scope and triggers

### Configuration Variables
Replace these variables in the configuration profile:
- `${ENROLLMENT_TOKEN}` - Your enrollment token
- `${API_URL}` - Your API endpoint
- `${ORG_ID}` - Your organization ID

## 🔍 Verification
After deployment, verify:
1. Agent is running: `sudo launchctl list | grep zerotrace`
2. Logs are generated: `tail -f /var/log/zerotrace-agent.log`
3. Tray icon appears (green = connected, gray = disconnected)

## 📞 Support
For deployment issues, contact: support@zerotrace.com
