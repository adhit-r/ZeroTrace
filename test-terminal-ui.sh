#!/bin/bash

echo "🎨 Testing Terminal-Inspired Enterprise UI..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔧 Testing VulnDetect Terminal UI${NC}"
echo "Theme: Terminal-Inspired Dark Theme"
echo "Design: Enterprise-grade with sharp edges"
echo ""

# Function to check if Node.js is installed
check_node() {
    echo -e "${YELLOW}🔍 Checking Node.js installation...${NC}"
    if command -v node &> /dev/null; then
        echo -e "${GREEN}✅ Node.js is installed${NC}"
        node --version
    else
        echo -e "${RED}❌ Node.js is not installed${NC}"
        echo "Please install Node.js to run the UI"
        exit 1
    fi
}

# Function to check if bun is installed
check_bun() {
    echo -e "${YELLOW}🔍 Checking Bun installation...${NC}"
    if command -v bun &> /dev/null; then
        echo -e "${GREEN}✅ Bun is installed${NC}"
        bun --version
    else
        echo -e "${RED}❌ Bun is not installed${NC}"
        echo "Please install Bun to run the UI"
        exit 1
    fi
}

# Function to install dependencies
install_deps() {
    echo -e "${YELLOW}📦 Installing dependencies...${NC}"
    cd web-react
    
    if [ -f "bun.lock" ]; then
        echo "Using Bun for package management..."
        bun install
    else
        echo "Using npm for package management..."
        npm install
    fi
    
    cd ..
}

# Function to start the development server
start_dev_server() {
    echo -e "${YELLOW}🚀 Starting development server...${NC}"
    cd web-react
    
    if command -v bun &> /dev/null; then
        echo "Starting with Bun..."
        bun run dev &
    else
        echo "Starting with npm..."
        npm run dev &
    fi
    
    DEV_PID=$!
    cd ..
    
    echo -e "${GREEN}✅ Development server started (PID: $DEV_PID)${NC}"
    echo "The UI should be available at: http://localhost:5173"
    echo ""
    echo -e "${BLUE}📋 Terminal UI Features:${NC}"
    echo "  ✅ Terminal-inspired dark theme"
    echo "  ✅ Sharp edges and modern design"
    echo "  ✅ Gold accent colors"
    echo "  ✅ Monospace typography"
    echo "  ✅ Scanline effects"
    echo "  ✅ Glow animations"
    echo "  ✅ Enterprise-grade components"
    echo ""
    echo -e "${YELLOW}🎯 Test the following:${NC}"
    echo "  1. Login page with terminal styling"
    echo "  2. Dashboard with terminal cards"
    echo "  3. Agent monitoring page"
    echo "  4. Navigation with terminal effects"
    echo "  5. Responsive design on mobile"
    echo ""
    echo -e "${GREEN}🌐 Open http://localhost:5173 in your browser${NC}"
    echo ""
    echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
    
    # Wait for user to stop
    wait $DEV_PID
}

# Function to check UI files
check_ui_files() {
    echo -e "${YELLOW}📁 Checking UI files...${NC}"
    
    if [ -f "web-react/src/styles/terminal-theme.css" ]; then
        echo -e "${GREEN}✅ Terminal theme CSS exists${NC}"
    else
        echo -e "${RED}❌ Terminal theme CSS missing${NC}"
        return 1
    fi
    
    if [ -f "web-react/src/components/Layout.tsx" ]; then
        echo -e "${GREEN}✅ Layout component exists${NC}"
    else
        echo -e "${RED}❌ Layout component missing${NC}"
        return 1
    fi
    
    if [ -f "web-react/src/pages/Dashboard.tsx" ]; then
        echo -e "${GREEN}✅ Dashboard component exists${NC}"
    else
        echo -e "${RED}❌ Dashboard component missing${NC}"
        return 1
    fi
    
    if [ -f "web-react/src/pages/Login.tsx" ]; then
        echo -e "${GREEN}✅ Login component exists${NC}"
    else
        echo -e "${RED}❌ Login component missing${NC}"
        return 1
    fi
    
    if [ -f "web-react/src/pages/Agents.tsx" ]; then
        echo -e "${GREEN}✅ Agents component exists${NC}"
    else
        echo -e "${RED}❌ Agents component missing${NC}"
        return 1
    fi
    
    echo ""
}

# Main execution
main() {
    echo -e "${BLUE}🚀 Starting Terminal UI Test${NC}"
    echo "=================================="
    echo ""
    
    # Check dependencies
    check_node
    check_bun
    
    # Check UI files
    if ! check_ui_files; then
        echo -e "${RED}❌ UI files are missing${NC}"
        exit 1
    fi
    
    # Install dependencies
    install_deps
    
    # Start development server
    start_dev_server
}

# Run main function
main
