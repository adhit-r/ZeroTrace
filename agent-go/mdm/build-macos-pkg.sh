#!/bin/bash

# ZeroTrace Agent - macOS Package Builder for MDM
# Builds .pkg files for deployment via MDM (Intune, Jamf, etc.)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
AGENT_NAME="ZeroTrace Agent"
BUNDLE_ID="com.zerotrace.agent"
VERSION="1.0.0"
BUILD_DIR="build"
PACKAGE_DIR="package"
DIST_DIR="dist"

echo -e "${BLUE}️  Building ZeroTrace Agent for MDM Deployment${NC}"
echo "=================================================="
echo ""

# Function to check dependencies
check_dependencies() {
    echo -e "${YELLOW} Checking dependencies...${NC}"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED} Go is not installed${NC}"
        exit 1
    fi
    
    # Check if pkgbuild is available
    if ! command -v pkgbuild &> /dev/null; then
        echo -e "${RED} pkgbuild is not available (requires Xcode Command Line Tools)${NC}"
        exit 1
    fi
    
    echo -e "${GREEN} All dependencies available${NC}"
}

# Function to clean build directories
clean_build() {
    echo -e "${YELLOW} Cleaning build directories...${NC}"
    rm -rf "$BUILD_DIR" "$PACKAGE_DIR" "$DIST_DIR"
    mkdir -p "$BUILD_DIR" "$PACKAGE_DIR" "$DIST_DIR"
}

# Function to build the agent
build_agent() {
    echo -e "${YELLOW} Building agent binary...${NC}"
    
    cd ..
    
    # Build for macOS (simple agent without tray)
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "mdm/$BUILD_DIR/zerotrace-agent" cmd/agent-simple/main.go
    
    # Build for Apple Silicon (simple agent without tray)
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "mdm/$BUILD_DIR/zerotrace-agent-arm64" cmd/agent-simple/main.go
    
    cd mdm
    
    echo -e "${GREEN} Agent binaries built${NC}"
}

# Function to create LaunchDaemon plist
create_launchdaemon() {
    echo -e "${YELLOW} Creating LaunchDaemon plist...${NC}"
    
    cat > "$BUILD_DIR/com.zerotrace.agent.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.zerotrace.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/zerotrace-agent</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/var/log/zerotrace-agent.log</string>
    <key>StandardErrorPath</key>
    <string>/var/log/zerotrace-agent.error.log</string>
    <key>WorkingDirectory</key>
    <string>/usr/local/bin</string>
    <key>ProcessType</key>
    <string>Background</string>
</dict>
</plist>
EOF
    
    echo -e "${GREEN} LaunchDaemon plist created${NC}"
}

# Function to create post-install script
create_postinstall() {
    echo -e "${YELLOW} Creating post-install script...${NC}"
    
    cat > "$BUILD_DIR/postinstall" << 'EOF'
#!/bin/bash

# ZeroTrace Agent Post-Install Script
# Sets up the agent as a system service

set -e

# Set permissions
chmod 755 /usr/local/bin/zerotrace-agent
chown root:wheel /usr/local/bin/zerotrace-agent

# Load LaunchDaemon
launchctl load /Library/LaunchDaemons/com.zerotrace.agent.plist

# Create log directory
mkdir -p /var/log
touch /var/log/zerotrace-agent.log
touch /var/log/zerotrace-agent.error.log
chmod 644 /var/log/zerotrace-agent*.log

echo "ZeroTrace Agent installed successfully"
exit 0
EOF
    
    chmod +x "$BUILD_DIR/postinstall"
    echo -e "${GREEN} Post-install script created${NC}"
}

# Function to create post-uninstall script
create_postuninstall() {
    echo -e "${YELLOW} Creating post-uninstall script...${NC}"
    
    cat > "$BUILD_DIR/postuninstall" << 'EOF'
#!/bin/bash

# ZeroTrace Agent Post-Uninstall Script
# Cleans up the agent service

set -e

# Unload LaunchDaemon
launchctl unload /Library/LaunchDaemons/com.zerotrace.agent.plist 2>/dev/null || true

# Remove files
rm -f /usr/local/bin/zerotrace-agent
rm -f /Library/LaunchDaemons/com.zerotrace.agent.plist
rm -f /var/log/zerotrace-agent*.log

echo "ZeroTrace Agent uninstalled successfully"
exit 0
EOF
    
    chmod +x "$BUILD_DIR/postuninstall"
    echo -e "${GREEN} Post-uninstall script created${NC}"
}

# Function to create package structure
create_package_structure() {
    echo -e "${YELLOW} Creating package structure...${NC}"
    
    # Create directory structure
    mkdir -p "$PACKAGE_DIR/usr/local/bin"
    mkdir -p "$PACKAGE_DIR/Library/LaunchDaemons"
    mkdir -p "$PACKAGE_DIR/var/log"
    
    # Copy files
    cp "$BUILD_DIR/zerotrace-agent" "$PACKAGE_DIR/usr/local/bin/"
    cp "$BUILD_DIR/com.zerotrace.agent.plist" "$PACKAGE_DIR/Library/LaunchDaemons/"
    
    # Create log files
    touch "$PACKAGE_DIR/var/log/zerotrace-agent.log"
    touch "$PACKAGE_DIR/var/log/zerotrace-agent.error.log"
    
    echo -e "${GREEN} Package structure created${NC}"
}

# Function to build the package
build_package() {
    echo -e "${YELLOW} Building macOS package...${NC}"
    
    # Build component package
    pkgbuild \
        --root "$PACKAGE_DIR" \
        --scripts "$BUILD_DIR" \
        --identifier "$BUNDLE_ID" \
        --version "$VERSION" \
        --install-location "/" \
        "$DIST_DIR/zerotrace-agent-component.pkg"
    
    # Build distribution package
    productbuild \
        --distribution distribution.xml \
        --package-path "$DIST_DIR" \
        --version "$VERSION" \
        "$DIST_DIR/ZeroTrace-Agent-$VERSION.pkg"
    
    echo -e "${GREEN} Package built: $DIST_DIR/ZeroTrace-Agent-$VERSION.pkg${NC}"
}

# Function to create distribution XML
create_distribution_xml() {
    echo -e "${YELLOW} Creating distribution XML...${NC}"
    
    cat > distribution.xml << EOF
<?xml version="1.0" encoding="utf-8"?>
<installer-gui-script minSpecVersion="1">
    <title>ZeroTrace Agent</title>
    <organization>com.zerotrace</organization>
    <domains enable_localSystem="true"/>
    <options customize="never" require-scripts="true" rootVolumeOnly="true"/>
    <pkg-ref id="$BUNDLE_ID"/>
    <choices-outline>
        <line choice="$BUNDLE_ID"/>
    </choices-outline>
    <choice id="$BUNDLE_ID" title="ZeroTrace Agent">
        <pkg-ref id="$BUNDLE_ID"/>
    </choice>
    <pkg-ref id="$BUNDLE_ID" version="$VERSION" onConclusion="none">zerotrace-agent-component.pkg</pkg-ref>
</installer-gui-script>
EOF
    
    echo -e "${GREEN} Distribution XML created${NC}"
}

# Function to create MDM configuration
create_mdm_config() {
    echo -e "${YELLOW} Creating MDM configuration...${NC}"
    
    cat > "$DIST_DIR/zerotrace-agent.mobileconfig" << EOF
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
                                <string>\${ENROLLMENT_TOKEN}</string>
                                <key>APIURL</key>
                                <string>\${API_URL}</string>
                                <key>OrganizationID</key>
                                <string>\${ORG_ID}</string>
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
    <string>$(uuidgen)</string>
    <key>PayloadDisplayName</key>
    <string>ZeroTrace Agent</string>
    <key>PayloadDescription</key>
    <string>ZeroTrace Agent Configuration</string>
    <key>PayloadOrganization</key>
    <string>ZeroTrace</string>
</dict>
</plist>
EOF
    
    echo -e "${GREEN} MDM configuration created${NC}"
}

# Function to create deployment guide
create_deployment_guide() {
    echo -e "${YELLOW} Creating deployment guide...${NC}"
    
    cat > "$DIST_DIR/DEPLOYMENT_GUIDE.md" << 'EOF'
# ZeroTrace Agent - MDM Deployment Guide

##  Package Contents
- `ZeroTrace-Agent-1.0.0.pkg` - Main installation package
- `zerotrace-agent.mobileconfig` - MDM configuration profile

##  Deployment Steps

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

##  Verification
After deployment, verify:
1. Agent is running: `sudo launchctl list | grep zerotrace`
2. Logs are generated: `tail -f /var/log/zerotrace-agent.log`
3. Tray icon appears (green = connected, gray = disconnected)

##  Support
For deployment issues, contact: support@zerotrace.com
EOF
    
    echo -e "${GREEN} Deployment guide created${NC}"
}

# Main execution
main() {
    echo -e "${BLUE} Starting ZeroTrace Agent MDM Package Build${NC}"
    echo "=================================================="
    echo ""
    
    # Check dependencies
    check_dependencies
    
    # Clean and prepare
    clean_build
    
    # Build components
    build_agent
    create_launchdaemon
    create_postinstall
    create_postuninstall
    create_package_structure
    create_distribution_xml
    
    # Build package
    build_package
    
    # Create MDM assets
    create_mdm_config
    create_deployment_guide
    
    echo ""
    echo -e "${GREEN} Build completed successfully!${NC}"
    echo ""
    echo -e "${BLUE} Generated files:${NC}"
    echo "  • $DIST_DIR/ZeroTrace-Agent-$VERSION.pkg"
    echo "  • $DIST_DIR/zerotrace-agent.mobileconfig"
    echo "  • $DIST_DIR/DEPLOYMENT_GUIDE.md"
    echo ""
    echo -e "${YELLOW} Next steps:${NC}"
    echo "  1. Upload package to your MDM platform"
    echo "  2. Configure enrollment tokens"
    echo "  3. Deploy to target devices"
    echo "  4. Monitor installation status"
}

# Run main function
main
