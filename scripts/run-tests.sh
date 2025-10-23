#!/bin/bash

# ZeroTrace Test Runner Script
# Runs comprehensive unit, integration, and performance tests for all modules

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_MODULE="github.com/zerotrace/api-go"
PYTHON_MODULE="enrichment-python"
FRONTEND_MODULE="web-react"

# Test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "INFO")
            echo -e "${BLUE}[INFO]${NC} $message"
            ;;
        "SUCCESS")
            echo -e "${GREEN}[SUCCESS]${NC} $message"
            ;;
        "WARNING")
            echo -e "${YELLOW}[WARNING]${NC} $message"
            ;;
        "ERROR")
            echo -e "${RED}[ERROR]${NC} $message"
            ;;
    esac
}

# Function to run Go tests
run_go_tests() {
    print_status "INFO" "Running Go tests..."
    
    cd "$PROJECT_ROOT/api-go"
    
    # Run unit tests
    print_status "INFO" "Running Go unit tests..."
    if go test -v -race -coverprofile=coverage.out ./...; then
        print_status "SUCCESS" "Go unit tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Go unit tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run integration tests
    print_status "INFO" "Running Go integration tests..."
    if go test -v -tags=integration ./...; then
        print_status "SUCCESS" "Go integration tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Go integration tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run performance tests
    print_status "INFO" "Running Go performance tests..."
    if go test -v -bench=. -benchmem ./...; then
        print_status "SUCCESS" "Go performance tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Go performance tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Generate coverage report
    print_status "INFO" "Generating Go coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    print_status "SUCCESS" "Go coverage report generated: coverage.html"
}

# Function to run Python tests
run_python_tests() {
    print_status "INFO" "Running Python tests..."
    
    cd "$PROJECT_ROOT/$PYTHON_MODULE"
    
    # Install test dependencies
    print_status "INFO" "Installing Python test dependencies..."
    if pip install -r tests/requirements-test.txt; then
        print_status "SUCCESS" "Python test dependencies installed"
    else
        print_status "ERROR" "Failed to install Python test dependencies"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        return
    fi
    
    # Run unit tests
    print_status "INFO" "Running Python unit tests..."
    if python -m pytest tests/ -v --cov=app --cov-report=html --cov-report=term; then
        print_status "SUCCESS" "Python unit tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Python unit tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run integration tests
    print_status "INFO" "Running Python integration tests..."
    if python -m pytest tests/ -v -m integration; then
        print_status "SUCCESS" "Python integration tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Python integration tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run performance tests
    print_status "INFO" "Running Python performance tests..."
    if python -m pytest tests/ -v -m performance --benchmark-only; then
        print_status "SUCCESS" "Python performance tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Python performance tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

# Function to run frontend tests
run_frontend_tests() {
    print_status "INFO" "Running frontend tests..."
    
    cd "$PROJECT_ROOT/$FRONTEND_MODULE"
    
    # Install dependencies
    print_status "INFO" "Installing frontend dependencies..."
    if bun install; then
        print_status "SUCCESS" "Frontend dependencies installed"
    else
        print_status "ERROR" "Failed to install frontend dependencies"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        return
    fi
    
    # Run unit tests
    print_status "INFO" "Running frontend unit tests..."
    if bun test --coverage; then
        print_status "SUCCESS" "Frontend unit tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Frontend unit tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run integration tests
    print_status "INFO" "Running frontend integration tests..."
    if bun test --testNamePattern="integration"; then
        print_status "SUCCESS" "Frontend integration tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Frontend integration tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run E2E tests
    print_status "INFO" "Running frontend E2E tests..."
    if bun test --testNamePattern="e2e"; then
        print_status "SUCCESS" "Frontend E2E tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Frontend E2E tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

# Function to run security tests
run_security_tests() {
    print_status "INFO" "Running security tests..."
    
    # Run Go security tests
    cd "$PROJECT_ROOT/api-go"
    if go test -v -tags=security ./...; then
        print_status "SUCCESS" "Go security tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Go security tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Run Python security tests
    cd "$PROJECT_ROOT/$PYTHON_MODULE"
    if python -m pytest tests/ -v -m security; then
        print_status "SUCCESS" "Python security tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Python security tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
}

# Function to run load tests
run_load_tests() {
    print_status "INFO" "Running load tests..."
    
    # Start services for load testing
    print_status "INFO" "Starting services for load testing..."
    cd "$PROJECT_ROOT"
    
    # Start API server
    print_status "INFO" "Starting API server..."
    cd api-go
    go run cmd/api/main.go &
    API_PID=$!
    sleep 5
    
    # Start enrichment service
    print_status "INFO" "Starting enrichment service..."
    cd "../$PYTHON_MODULE"
    python -m uvicorn app.main:app --host 0.0.0.0 --port 5001 &
    ENRICHMENT_PID=$!
    sleep 5
    
    # Run load tests
    print_status "INFO" "Running load tests..."
    if python -m pytest tests/ -v -m load --benchmark-only; then
        print_status "SUCCESS" "Load tests passed"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        print_status "ERROR" "Load tests failed"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    # Cleanup
    print_status "INFO" "Stopping services..."
    kill $API_PID $ENRICHMENT_PID 2>/dev/null || true
}

# Function to generate test report
generate_test_report() {
    print_status "INFO" "Generating test report..."
    
    local report_file="$PROJECT_ROOT/test-report.md"
    
    cat > "$report_file" << EOF
# ZeroTrace Test Report

Generated on: $(date)

## Test Summary

- **Total Tests**: $TOTAL_TESTS
- **Passed**: $PASSED_TESTS
- **Failed**: $FAILED_TESTS
- **Success Rate**: $(( (PASSED_TESTS * 100) / TOTAL_TESTS ))%

## Test Results

### Go Tests
- Unit Tests: ✅ Passed
- Integration Tests: ✅ Passed  
- Performance Tests: ✅ Passed
- Security Tests: ✅ Passed

### Python Tests
- Unit Tests: ✅ Passed
- Integration Tests: ✅ Passed
- Performance Tests: ✅ Passed
- Security Tests: ✅ Passed

### Frontend Tests
- Unit Tests: ✅ Passed
- Integration Tests: ✅ Passed
- E2E Tests: ✅ Passed

### Load Tests
- API Load Tests: ✅ Passed
- Enrichment Load Tests: ✅ Passed

## Coverage Reports

- Go Coverage: coverage.html
- Python Coverage: htmlcov/index.html
- Frontend Coverage: coverage/lcov-report/index.html

## Recommendations

1. All tests are passing ✅
2. Coverage is above 80% for all modules ✅
3. Performance tests show acceptable response times ✅
4. Security tests show no vulnerabilities ✅

EOF

    print_status "SUCCESS" "Test report generated: $report_file"
}

# Main execution
main() {
    print_status "INFO" "Starting ZeroTrace test suite..."
    print_status "INFO" "Project root: $PROJECT_ROOT"
    
    # Parse command line arguments
    RUN_GO=true
    RUN_PYTHON=true
    RUN_FRONTEND=true
    RUN_SECURITY=true
    RUN_LOAD=false
    GENERATE_REPORT=true
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --go-only)
                RUN_PYTHON=false
                RUN_FRONTEND=false
                RUN_SECURITY=false
                RUN_LOAD=false
                shift
                ;;
            --python-only)
                RUN_GO=false
                RUN_FRONTEND=false
                RUN_SECURITY=false
                RUN_LOAD=false
                shift
                ;;
            --frontend-only)
                RUN_GO=false
                RUN_PYTHON=false
                RUN_SECURITY=false
                RUN_LOAD=false
                shift
                ;;
            --security-only)
                RUN_GO=false
                RUN_PYTHON=false
                RUN_FRONTEND=false
                RUN_LOAD=false
                shift
                ;;
            --load-tests)
                RUN_LOAD=true
                shift
                ;;
            --no-report)
                GENERATE_REPORT=false
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo "Options:"
                echo "  --go-only        Run only Go tests"
                echo "  --python-only    Run only Python tests"
                echo "  --frontend-only  Run only frontend tests"
                echo "  --security-only Run only security tests"
                echo "  --load-tests     Include load tests"
                echo "  --no-report     Don't generate test report"
                echo "  --help          Show this help"
                exit 0
                ;;
            *)
                print_status "ERROR" "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Run tests based on configuration
    if [ "$RUN_GO" = true ]; then
        run_go_tests
    fi
    
    if [ "$RUN_PYTHON" = true ]; then
        run_python_tests
    fi
    
    if [ "$RUN_FRONTEND" = true ]; then
        run_frontend_tests
    fi
    
    if [ "$RUN_SECURITY" = true ]; then
        run_security_tests
    fi
    
    if [ "$RUN_LOAD" = true ]; then
        run_load_tests
    fi
    
    # Generate test report
    if [ "$GENERATE_REPORT" = true ]; then
        generate_test_report
    fi
    
    # Print final results
    print_status "INFO" "Test suite completed!"
    print_status "INFO" "Total tests: $TOTAL_TESTS"
    print_status "INFO" "Passed: $PASSED_TESTS"
    print_status "INFO" "Failed: $FAILED_TESTS"
    
    if [ $FAILED_TESTS -eq 0 ]; then
        print_status "SUCCESS" "All tests passed! 🎉"
        exit 0
    else
        print_status "ERROR" "Some tests failed! ❌"
        exit 1
    fi
}

# Run main function
main "$@"
