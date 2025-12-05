#!/bin/bash

# ZeroTrace Agent - macOS DMG Builder
# Creates a user-friendly DMG installer for macOS

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
DIST_DIR="dist"
DMG_NAME="ZeroTrace-Agent-${VERSION}.dmg"
DMG_PATH="${DIST_DIR}/${DMG_NAME}"

echo -e "${BLUE}🏗️  Building ZeroTrace Agent DMG for macOS${NC}"
echo "=================================================="
echo ""

# Function to check dependencies
check_dependencies() {
    echo -e "${YELLOW}🔍 Checking dependencies...${NC}"
    
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go is not installed${NC}"
        exit 1
    fi
    
    if ! command -v hdiutil &> /dev/null; then
        echo -e "${RED}❌ hdiutil is not available (macOS only)${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ All dependencies available${NC}"
}

# Function to clean build directories
clean_build() {
    echo -e "${YELLOW}🧹 Cleaning build directories...${NC}"
    rm -rf "$BUILD_DIR" "$DIST_DIR" "dmg_temp"
    mkdir -p "$BUILD_DIR" "$DIST_DIR" "dmg_temp"
}

# Function to build the agent with tray UI
build_agent() {
    echo -e "${YELLOW}🔨 Building agent binary with tray UI...${NC}"
    
    cd ..
    
    # Build for macOS (with tray UI)
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "mdm/$BUILD_DIR/zerotrace-agent" cmd/agent/main.go
    
    # Build for Apple Silicon (with tray UI)
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "mdm/$BUILD_DIR/zerotrace-agent-arm64" cmd/agent/main.go
    
    # Use universal binary if on macOS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo -e "${YELLOW}📦 Creating universal binary...${NC}"
        lipo -create \
            "mdm/$BUILD_DIR/zerotrace-agent" \
            "mdm/$BUILD_DIR/zerotrace-agent-arm64" \
            -output "mdm/$BUILD_DIR/zerotrace-agent-universal" 2>/dev/null || {
            echo -e "${YELLOW}⚠️  Universal binary creation skipped (using arm64)${NC}"
            cp "mdm/$BUILD_DIR/zerotrace-agent-arm64" "mdm/$BUILD_DIR/zerotrace-agent-universal"
        }
    else
        cp "mdm/$BUILD_DIR/zerotrace-agent-arm64" "mdm/$BUILD_DIR/zerotrace-agent-universal"
    fi
    
    cd mdm
    
    echo -e "${GREEN}✅ Agent binaries built${NC}"
}

# Function to create DMG structure
create_dmg_structure() {
    echo -e "${YELLOW}📁 Creating DMG structure...${NC}"
    
    # Create Applications link
    ln -s /Applications "dmg_temp/Applications"
    
    # Copy .app bundle if it exists, otherwise copy binary
    if [ -d "$BUILD_DIR/ZeroTrace Agent.app" ]; then
        echo -e "${YELLOW}📦 Using .app bundle for menu bar icon support${NC}"
        cp -R "$BUILD_DIR/ZeroTrace Agent.app" "dmg_temp/"
    else
        echo -e "${YELLOW}📦 Using binary (no menu bar icon)${NC}"
        cp "$BUILD_DIR/zerotrace-agent-universal" "dmg_temp/zerotrace-agent"
        chmod +x "dmg_temp/zerotrace-agent"
    fi
    
    # Create README
    cat > "dmg_temp/README.txt" << EOF
ZeroTrace Agent ${VERSION}
=======================

Installation:
1. Drag "ZeroTrace Agent.app" (or "zerotrace-agent") to Applications
2. Open Applications folder
3. Double-click to launch

The agent will:
- Run in the background
- Show a menu bar icon next to WiFi (if .app bundle)
- Automatically scan for vulnerabilities
- Send results to ZeroTrace API

Menu Bar Icon:
- Appears in top-right menu bar next to WiFi
- Click icon for status, CPU usage, and options
- Green = Connected, Gray = Checking, Red = Error

Configuration:
- Edit ~/.zerotrace/.env to configure API endpoint
- The agent will create this file on first run

System Requirements:
- macOS 10.15 or later
- Internet connection for API communication

For support, visit: https://zerotrace.com
EOF
    
    # Create .DS_Store for better appearance
    # This is optional but makes the DMG look better
    echo -e "${YELLOW}📝 DMG structure created${NC}"
}

# Function to create DMG
create_dmg() {
    echo -e "${YELLOW}💿 Creating DMG file...${NC}"
    
    # Calculate size (add 20MB overhead)
    SIZE=$(du -sm dmg_temp | cut -f1)
    SIZE=$((SIZE + 20))
    
    # Create DMG
    hdiutil create -srcfolder "dmg_temp" \
        -volname "ZeroTrace Agent ${VERSION}" \
        -fs HFS+ \
        -fsargs "-c c=64,a=16,e=16" \
        -format UDRW \
        -size ${SIZE}M \
        "${DMG_PATH}.temp" || {
        echo -e "${RED}❌ Failed to create DMG${NC}"
        exit 1
    }
    
    # Mount DMG
    MOUNT_DIR=$(hdiutil attach -readwrite -noverify -noautoopen "${DMG_PATH}.temp" | \
        egrep '^/dev/' | sed 1q | awk '{print $3}')
    
    # Set volume icon (optional)
    # cp icon.icns "$MOUNT_DIR/.VolumeIcon.icns" 2>/dev/null || true
    
    # Set background (optional - requires background image)
    # osascript << EOF
    # tell application "Finder"
    #     tell disk "ZeroTrace Agent ${VERSION}"
    #         open
    #         set current view of container window to icon view
    #         set toolbar visible of container window to false
    #         set statusbar visible of container window to false
    #         set bounds of container window to {400, 100, 920, 420}
    #         set viewOptions to the icon view options of container window
    #         set arrangement of viewOptions to not arranged
    #         set icon size of viewOptions to 72
    #         set background picture of viewOptions to file ".background:background.png"
    #         set position of item "zerotrace-agent" of container window to {160, 205}
    #         set position of item "Applications" of container window to {360, 205}
    #         close
    #         open
    #         update without registering applications
    #         delay 2
    #     end tell
    # end tell
    # EOF
    
    # Unmount
    hdiutil detach "$MOUNT_DIR"
    
    # Convert to compressed read-only DMG
    hdiutil convert "${DMG_PATH}.temp" \
        -format UDZO \
        -imagekey zlib-level=9 \
        -o "$DMG_PATH"
    
    # Remove temp file
    rm -f "${DMG_PATH}.temp"
    
    echo -e "${GREEN}✅ DMG created: ${DMG_PATH}${NC}"
}

# Function to show summary
show_summary() {
    echo ""
    echo -e "${GREEN}✅ Build Complete!${NC}"
    echo "=================================================="
    echo -e "${BLUE}DMG Location:${NC} ${DMG_PATH}"
    echo -e "${BLUE}Size:${NC} $(du -h "$DMG_PATH" | cut -f1)"
    echo ""
    echo -e "${YELLOW}Next Steps:${NC}"
    echo "1. Test the DMG by double-clicking it"
    echo "2. Verify the agent runs with tray icon"
    echo "3. Distribute to users"
    echo ""
}

# Main execution
main() {
    check_dependencies
    clean_build
    build_agent
    create_dmg_structure
    create_dmg
    show_summary
}

# Run main
main

