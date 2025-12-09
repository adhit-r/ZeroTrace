# ZeroTrace Agent - MDM Deployment Guide

## **Supported MDM Platforms**

### **Microsoft Intune**
- **Windows**: `.intunewin` package
- **macOS**: `.pkg` package
- **Configuration**: JSON policies

### **Jamf Pro**
- **macOS**: `.pkg` package
- **Configuration**: Configuration Profiles
- **Smart Groups**: Device targeting

### **Azure AD**
- **Enterprise Application**: App registration
- **Conditional Access**: Policy enforcement
- **Device Management**: Intune integration

### **VMware Workspace ONE**
- **UEM**: Unified endpoint management
- **macOS**: `.pkg` deployment
- **Windows**: `.msi` deployment

## **Deployment Packages**

### **macOS Package (.pkg)**
```bash
# Build macOS package
./build-macos-pkg.sh

# Package includes:
# - ZeroTrace Agent binary
# - LaunchDaemon plist
# - Configuration files
# - Post-install scripts
```

### **Windows Package (.msi)**
```powershell
# Build Windows package
.\build-windows-msi.ps1

# Package includes:
# - ZeroTrace Agent executable
# - Windows Service
# - Registry entries
# - Configuration files
```

## **Configuration**

### **Environment Variables**
```bash
# Required for enrollment
ZEROTRACE_ENROLLMENT_TOKEN=<token>
ZEROTRACE_API_URL=<api-url>
ZEROTRACE_ORG_ID=<org-id>

# Optional
ZEROTRACE_SCAN_INTERVAL=24h
ZEROTRACE_LOG_LEVEL=info
```

### **MDM Configuration Profile**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>PayloadContent</key>
    <array>
        <dict>
            <key>PayloadType</key>
            <string>com.apple.ManagedClient.preferences</string>
            <key>PayloadVersion</key>
            <integer>1</integer>
            <key>PayloadIdentifier</key>
            <string>com.zerotrace.agent.config</string>
            <key>PayloadDisplayName</key>
            <string>ZeroTrace Agent Configuration</string>
            <key>PayloadDescription</key>
            <string>Configures ZeroTrace Agent settings</string>
            <key>PayloadOrganization</key>
            <string>ZeroTrace</string>
            <key>PayloadScope</key>
            <string>System</string>
            <key>PayloadRemovalDisallowed</key>
            <false/>
            <key>PayloadContent</key>
            <dict>
                <key>com.zerotrace.agent</key>
                <dict>
                    <key>Forced</key>
                    <array>
                        <dict>
                            <key>mcx_preference_settings</key>
                            <dict>
                                <key>EnrollmentToken</key>
                                <string>${ENROLLMENT_TOKEN}</string>
                                <key>APIURL</key>
                                <string>${API_URL}</string>
                                <key>OrganizationID</key>
                                <string>${ORG_ID}</string>
                            </dict>
                        </dict>
                    </array>
                </dict>
            </dict>
        </dict>
    </array>
    <key>PayloadRemovalDisallowed</key>
    <false/>
    <key>PayloadType</key>
    <string>Configuration</string>
    <key>PayloadVersion</key>
    <integer>1</integer>
    <key>PayloadIdentifier</key>
    <string>com.zerotrace.agent</string>
    <key>PayloadUUID</key>
    <string>12345678-1234-1234-1234-123456789012</string>
    <key>PayloadDisplayName</key>
    <string>ZeroTrace Agent</string>
    <key>PayloadDescription</key>
    <string>ZeroTrace Agent Configuration</string>
    <key>PayloadOrganization</key>
    <string>ZeroTrace</string>
</dict>
</plist>
```

## **Deployment Steps**

### **1. Microsoft Intune**

#### **macOS Deployment**
1. **Upload Package**: Upload `.pkg` file to Intune
2. **Create App**: Configure as macOS app
3. **Assign**: Target devices/groups
4. **Deploy**: Install automatically

#### **Windows Deployment**
1. **Upload Package**: Upload `.intunewin` file
2. **Create App**: Configure as Windows app
3. **Assign**: Target devices/groups
4. **Deploy**: Install automatically

### **2. Jamf Pro**

#### **Package Deployment**
```bash
# 1. Upload package to Jamf
# 2. Create policy
# 3. Configure scope (Smart Groups)
# 4. Set trigger (Install)
```

#### **Configuration Profile**
1. **Upload Profile**: Upload `.mobileconfig`
2. **Scope**: Target devices
3. **Deploy**: Install automatically

### **3. Azure AD**

#### **Enterprise Application**
1. **Register App**: In Azure AD
2. **Configure SSO**: SAML/OIDC
3. **Assign Users**: Target groups
4. **Conditional Access**: Enforce policies

## **Deployment Checklist**

### **Pre-Deployment**
- [ ] Build agent packages
- [ ] Generate enrollment tokens
- [ ] Configure API endpoints
- [ ] Test in lab environment
- [ ] Create MDM policies

### **Deployment**
- [ ] Upload packages to MDM
- [ ] Configure installation policies
- [ ] Target device groups
- [ ] Deploy configuration profiles
- [ ] Monitor installation status

### **Post-Deployment**
- [ ] Verify agent enrollment
- [ ] Check API connectivity
- [ ] Monitor data collection
- [ ] Validate security policies
- [ ] Document deployment

## **Monitoring & Troubleshooting**

### **Agent Status**
```bash
# Check agent status
sudo launchctl list | grep zerotrace

# View logs
sudo log show --predicate 'process == "zerotrace-agent"' --last 1h

# Check configuration
sudo defaults read /Library/Preferences/com.zerotrace.agent
```

### **MDM Status**
```bash
# Check MDM enrollment
sudo profiles show -type configuration

# View MDM logs
sudo log show --predicate 'process == "mdm"' --last 1h
```

### **Common Issues**
1. **Enrollment Token Invalid**: Regenerate token
2. **API Connection Failed**: Check network/firewall
3. **Package Installation Failed**: Check MDM logs
4. **Configuration Not Applied**: Verify profile scope

## **Support**

For MDM deployment support:
- **Documentation**: [docs.zerotrace.com/mdm](https://docs.zerotrace.com/mdm)
- **Support**: [support@zerotrace.com](mailto:support@zerotrace.com)
- **Community**: [community.zerotrace.com](https://community.zerotrace.com)
