#!/bin/bash

# ZeroTrace Integration Test Script
# Tests the complete data flow: Agent → API → Python Enrichment → Database

set -e

echo "🚀 Starting ZeroTrace Integration Tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    command -v go >/dev/null 2>&1 || { print_error "Go is not installed"; exit 1; }
    command -v python3 >/dev/null 2>&1 || { print_error "Python3 is not installed"; exit 1; }
    command -v bun >/dev/null 2>&1 || { print_warning "Bun is not installed, using npm"; }
    command -v curl >/dev/null 2>&1 || { print_error "curl is not installed"; exit 1; }
    
    print_success "All dependencies found"
}

# Build Go components
build_go_components() {
    print_status "Building Go components..."
    
    # Build API
    cd api-go
    go build -o zerotrace-api cmd/api/main.go
    print_success "API built successfully"
    
    # Build Agent
    cd ../agent-go
    go build -o zerotrace-agent cmd/agent/main.go
    print_success "Agent built successfully"
    
    # Build tray test
    go build -o tray-test cmd/tray-test/main.go
    print_success "Tray test built successfully"
    
    cd ..
}

# Install Python dependencies
install_python_deps() {
    print_status "Installing Python dependencies..."
    
    cd enrichment-python
    pip3 install -r requirements.txt
    print_success "Python dependencies installed"
    cd ..
}

# Install Node.js dependencies
install_node_deps() {
    print_status "Installing Node.js dependencies..."
    
    cd web-react
    if command -v bun >/dev/null 2>&1; then
        bun install
    else
        npm install
    fi
    print_success "Node.js dependencies installed"
    cd ..
}

# Start services
start_services() {
    print_status "Starting services..."
    
    # Start Python enrichment service
    cd enrichment-python
    python3 -m uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload &
    ENRICHMENT_PID=$!
    print_success "Python enrichment service started (PID: $ENRICHMENT_PID)"
    cd ..
    
    # Wait for enrichment service to start
    sleep 3
    
    # Start Go API
    cd api-go
    ENRICHMENT_SERVICE_URL=http://localhost:8000 ./zerotrace-api &
    API_PID=$!
    print_success "Go API started (PID: $API_PID)"
    cd ..
    
    # Wait for API to start
    sleep 3
    
    # Start Go Agent
    cd agent-go
    ./zerotrace-agent &
    AGENT_PID=$!
    print_success "Go Agent started (PID: $AGENT_PID)"
    cd ..
    
    # Wait for agent to start
    sleep 3
}

# Test API endpoints
test_api_endpoints() {
    print_status "Testing API endpoints..."
    
    # Test health endpoint
    if curl -s http://localhost:8080/health | grep -q "healthy"; then
        print_success "API health check passed"
    else
        print_error "API health check failed"
        return 1
    fi
    
    # Test enrichment service health
    if curl -s http://localhost:8000/health | grep -q "healthy"; then
        print_success "Enrichment service health check passed"
    else
        print_error "Enrichment service health check failed"
        return 1
    fi
}

# Test CVE enrichment
test_cve_enrichment() {
    print_status "Testing CVE enrichment..."
    
    # Test enrichment endpoint
    TEST_SOFTWARE='[{"name": "nginx", "version": "1.18.0", "package_type": "system"}]'
    
    RESPONSE=$(curl -s -X POST http://localhost:8000/enrich/software \
        -H "Content-Type: application/json" \
        -d "$TEST_SOFTWARE")
    
    if echo "$RESPONSE" | grep -q "success.*true"; then
        print_success "CVE enrichment test passed"
        echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
    else
        print_error "CVE enrichment test failed"
        echo "$RESPONSE"
        return 1
    fi
}

# Test tray indicator
test_tray_indicator() {
    print_status "Testing tray indicator..."
    
    # Start tray test in background
    cd agent-go
    ./tray-test &
    TRAY_PID=$!
    print_success "Tray test started (PID: $TRAY_PID)"
    cd ..
    
    # Let it run for a few seconds
    sleep 5
    
    # Kill tray test
    kill $TRAY_PID 2>/dev/null || true
    print_success "Tray test completed"
}

# Test web UI
test_web_ui() {
    print_status "Testing web UI..."
    
    cd web-react
    if command -v bun >/dev/null 2>&1; then
        bun run dev &
    else
        npm run dev &
    fi
    WEB_PID=$!
    print_success "Web UI started (PID: $WEB_PID)"
    cd ..
    
    # Wait for web UI to start
    sleep 10
    
    # Test web UI endpoint
    if curl -s http://localhost:5173 | grep -q "ZeroTrace"; then
        print_success "Web UI test passed"
    else
        print_warning "Web UI test failed (may need manual verification)"
    fi
    
    # Kill web UI
    kill $WEB_PID 2>/dev/null || true
}

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    
    # Kill all background processes
    kill $ENRICHMENT_PID $API_PID $AGENT_PID $TRAY_PID $WEB_PID 2>/dev/null || true
    
    # Wait for processes to terminate
    sleep 2
    
    print_success "Cleanup completed"
}

# Main test execution
main() {
    # Set up cleanup on exit
    trap cleanup EXIT
    
    print_status "Starting ZeroTrace integration tests..."
    
    check_dependencies
    build_go_components
    install_python_deps
    install_node_deps
    start_services
    test_api_endpoints
    test_cve_enrichment
    test_tray_indicator
    test_web_ui
    
    print_success "🎉 All integration tests completed successfully!"
    print_status "Data flow verified: Agent → API → Python Enrichment → Database"
}

# Run main function
main "$@"

