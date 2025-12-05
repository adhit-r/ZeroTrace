#!/bin/bash

# ZeroTrace Agent - macOS App Bundle Builder
# Creates a proper macOS .app bundle with menu bar icon support

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

AGENT_NAME="ZeroTrace Agent"
BUNDLE_ID="com.zerotrace.agent"
VERSION="1.0.0"
APP_NAME="ZeroTrace Agent.app"
BUILD_DIR="build"
APP_DIR="${BUILD_DIR}/${APP_NAME}"
CONTENTS_DIR="${APP_DIR}/Contents"
MACOS_DIR="${CONTENTS_DIR}/MacOS"
RESOURCES_DIR="${CONTENTS_DIR}/Resources"

echo -e "${BLUE}🏗️  Building ZeroTrace Agent macOS App Bundle${NC}"
echo "=================================================="
echo ""

# Clean
echo -e "${YELLOW}🧹 Cleaning...${NC}"
rm -rf "$BUILD_DIR" "$APP_DIR"
mkdir -p "$MACOS_DIR" "$RESOURCES_DIR"

# Build agent
echo -e "${YELLOW}🔨 Building agent binary...${NC}"
cd ..
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "mdm/$MACOS_DIR/zerotrace-agent" cmd/agent/main.go
cd mdm

# Copy Info.plist
echo -e "${YELLOW}📝 Creating app bundle structure...${NC}"
cp Info.plist "$CONTENTS_DIR/"

# Create PkgInfo
echo "APPL????" > "$CONTENTS_DIR/PkgInfo"

# Copy icons to Resources
if [ -d "assets" ]; then
    cp -r assets/* "$RESOURCES_DIR/" 2>/dev/null || true
fi

# Make executable
chmod +x "$MACOS_DIR/zerotrace-agent"

# Create app bundle structure
cat > "$CONTENTS_DIR/Info.plist" << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleDevelopmentRegion</key>
	<string>en</string>
	<key>CFBundleExecutable</key>
	<string>zerotrace-agent</string>
	<key>CFBundleIdentifier</key>
	<string>${BUNDLE_ID}</string>
	<key>CFBundleInfoDictionaryVersion</key>
	<string>6.0</string>
	<key>CFBundleName</key>
	<string>${AGENT_NAME}</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>${VERSION}</string>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>LSMinimumSystemVersion</key>
	<string>10.15</string>
	<key>NSHighResolutionCapable</key>
	<true/>
	<key>LSUIElement</key>
	<true/>
	<key>NSAppTransportSecurity</key>
	<dict>
		<key>NSAllowsArbitraryLoads</key>
		<true/>
	</dict>
</dict>
</plist>
EOF

echo -e "${GREEN}✅ App bundle created: ${APP_DIR}${NC}"
echo ""
echo "To use:"
echo "1. Double-click the .app to launch"
echo "2. Menu bar icon will appear next to WiFi"
echo "3. Click icon for menu"


