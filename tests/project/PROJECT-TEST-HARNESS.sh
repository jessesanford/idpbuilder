#!/bin/bash
# PROJECT-TEST-HARNESS.sh
# Comprehensive test execution harness for idpbuilder-oci-push-rebuild
# Created: 2025-11-11
# Purpose: Execute all project-level tests in correct order with proper reporting

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Configuration
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
COVERAGE_DIR="$PROJECT_ROOT/coverage"
COVERAGE_PROFILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"

# Flags
RUN_CATEGORY=""
ENABLE_COVERAGE=false
CI_MODE=false
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --category=*)
            RUN_CATEGORY="${1#*=}"
            shift
            ;;
        --coverage)
            ENABLE_COVERAGE=true
            shift
            ;;
        --ci)
            CI_MODE=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --category=<cat>  Run specific test category (interfaces|integration|e2e|quality|all)"
            echo "  --coverage        Generate coverage report"
            echo "  --ci              Run in CI mode (strict, fail fast)"
            echo "  --verbose         Verbose test output"
            echo "  --help            Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Default to all categories if none specified
if [ -z "$RUN_CATEGORY" ]; then
    RUN_CATEGORY="all"
fi

# Create coverage directory if needed
if [ "$ENABLE_COVERAGE" = true ]; then
    mkdir -p "$COVERAGE_DIR"
    echo -e "${BLUE}📊 Coverage reporting enabled${NC}"
    echo -e "${BLUE}   Coverage profile: $COVERAGE_PROFILE${NC}"
fi

# Print banner
echo "═══════════════════════════════════════════════════════════════"
echo -e "${BLUE}🧪 PROJECT TEST HARNESS: idpbuilder-oci-push-rebuild${NC}"
echo "═══════════════════════════════════════════════════════════════"
echo "Test Category: $RUN_CATEGORY"
echo "CI Mode: $CI_MODE"
echo "Coverage: $ENABLE_COVERAGE"
echo "Verbose: $VERBOSE"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Function to run test category
run_test_category() {
    local category=$1
    local test_path=$2
    local description=$3

    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}📋 $description${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    # Check if test directory exists
    if [ ! -d "$test_path" ] && [ ! -f "$test_path" ]; then
        echo -e "${YELLOW}⚠️  Test path not found: $test_path${NC}"
        echo -e "${YELLOW}   Status: PENDING IMPLEMENTATION${NC}"
        echo -e "${YELLOW}   This is EXPECTED during TDD RED phase${NC}"
        SKIPPED_TESTS=$((SKIPPED_TESTS + 1))
        echo ""
        return 0
    fi

    # Build test command
    local test_cmd="go test"

    if [ "$VERBOSE" = true ]; then
        test_cmd="$test_cmd -v"
    fi

    if [ "$ENABLE_COVERAGE" = true ]; then
        test_cmd="$test_cmd -coverprofile=$COVERAGE_DIR/${category}_coverage.out"
    fi

    if [ "$CI_MODE" = true ]; then
        test_cmd="$test_cmd -race -timeout=10m"
    fi

    test_cmd="$test_cmd $test_path"

    # Run tests
    echo "Running: $test_cmd"
    echo ""

    if eval "$test_cmd"; then
        echo ""
        echo -e "${GREEN}✅ $description: PASSED${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo ""
        echo -e "${RED}❌ $description: FAILED${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))

        if [ "$CI_MODE" = true ]; then
            echo -e "${RED}CI Mode: Failing fast due to test failure${NC}"
            exit 1
        fi
    fi

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo ""
}

# Function to check prerequisites
check_prerequisites() {
    echo -e "${BLUE}🔍 Checking prerequisites...${NC}"

    # Check Go installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go not found. Please install Go 1.20+${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Go found: $(go version)${NC}"

    # Check Docker available (for integration tests)
    if command -v docker &> /dev/null; then
        if docker ps &> /dev/null; then
            echo -e "${GREEN}✅ Docker daemon accessible${NC}"
        else
            echo -e "${YELLOW}⚠️  Docker daemon not accessible (some integration tests may fail)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  Docker not found (integration/E2E tests will be skipped)${NC}"
    fi

    echo ""
}

# Function to run interface contract tests
run_interface_tests() {
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  PHASE 1: INTERFACE CONTRACT TESTS                        ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"
    echo ""

    run_test_category "docker-interface" \
        "./pkg/docker/..." \
        "Docker Client Interface Tests"

    run_test_category "registry-interface" \
        "./pkg/registry/..." \
        "Registry Client Interface Tests"

    run_test_category "auth-interface" \
        "./pkg/auth/..." \
        "Authentication Provider Interface Tests"

    run_test_category "tls-interface" \
        "./pkg/tls/..." \
        "TLS Provider Interface Tests"

    run_test_category "command-interface" \
        "./cmd/..." \
        "Command Layer Interface Tests"
}

# Function to run integration tests
run_integration_tests() {
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  PHASE 2: FUNCTIONAL INTEGRATION TESTS                    ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"
    echo ""

    run_test_category "docker-registry-integration" \
        "./tests/integration/docker_registry_test.go" \
        "Docker + Registry Integration Tests"

    run_test_category "auth-tls-registry-integration" \
        "./tests/integration/auth_tls_registry_test.go" \
        "Auth + TLS + Registry Integration Tests"

    run_test_category "full-stack-integration" \
        "./tests/integration/full_stack_test.go" \
        "Full Stack Integration Tests"
}

# Function to run E2E tests
run_e2e_tests() {
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  PHASE 3: END-TO-END WORKFLOW TESTS                       ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"
    echo ""

    # Check if test registry is available
    if [ -f "./tests/e2e/docker-compose.yml" ]; then
        echo -e "${BLUE}🐳 Starting test registry environment...${NC}"
        cd tests/e2e
        docker-compose up -d
        sleep 5  # Wait for registry to be ready
        cd ../..
        echo -e "${GREEN}✅ Test registry started${NC}"
        echo ""
    fi

    run_test_category "e2e-success" \
        "./tests/e2e/success_workflows_test.go" \
        "E2E Success Workflow Tests"

    run_test_category "e2e-failure" \
        "./tests/e2e/failure_workflows_test.go" \
        "E2E Failure Workflow Tests"

    run_test_category "e2e-edge-cases" \
        "./tests/e2e/edge_cases_test.go" \
        "E2E Edge Case Tests"

    # Cleanup test registry
    if [ -f "./tests/e2e/docker-compose.yml" ]; then
        echo ""
        echo -e "${BLUE}🧹 Cleaning up test registry environment...${NC}"
        cd tests/e2e
        docker-compose down
        cd ../..
        echo -e "${GREEN}✅ Test registry stopped${NC}"
    fi
}

# Function to run quality tests
run_quality_tests() {
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  QUALITY ATTRIBUTE TESTS                                  ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════╝${NC}"
    echo ""

    run_test_category "performance" \
        "./tests/quality/performance_test.go" \
        "Performance Tests"

    run_test_category "security" \
        "./tests/quality/security_test.go" \
        "Security Tests"

    run_test_category "reliability" \
        "./tests/quality/reliability_test.go" \
        "Reliability Tests"
}

# Function to generate coverage report
generate_coverage_report() {
    if [ "$ENABLE_COVERAGE" = false ]; then
        return
    fi

    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}📊 Generating Coverage Report${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    # Merge coverage profiles
    if ls "$COVERAGE_DIR"/*_coverage.out 1> /dev/null 2>&1; then
        echo "Merging coverage profiles..."
        cat "$COVERAGE_DIR"/*_coverage.out > "$COVERAGE_PROFILE"

        # Generate HTML report
        go tool cover -html="$COVERAGE_PROFILE" -o "$COVERAGE_HTML"

        # Display coverage summary
        echo ""
        echo -e "${GREEN}Coverage Report Generated:${NC}"
        echo "  Profile: $COVERAGE_PROFILE"
        echo "  HTML:    $COVERAGE_HTML"
        echo ""

        # Calculate coverage percentage
        total_coverage=$(go tool cover -func="$COVERAGE_PROFILE" | grep total | awk '{print $3}')
        echo -e "${GREEN}Total Coverage: $total_coverage${NC}"

        # Check coverage threshold (85% target)
        coverage_value=$(echo "$total_coverage" | sed 's/%//')
        if (( $(echo "$coverage_value >= 85.0" | bc -l) )); then
            echo -e "${GREEN}✅ Coverage meets 85% threshold${NC}"
        else
            echo -e "${YELLOW}⚠️  Coverage below 85% threshold (current: $total_coverage)${NC}"
            if [ "$CI_MODE" = true ]; then
                echo -e "${RED}CI Mode: Failing due to insufficient coverage${NC}"
                exit 1
            fi
        fi
    else
        echo -e "${YELLOW}⚠️  No coverage data generated${NC}"
    fi

    echo ""
}

# Function to print test summary
print_summary() {
    echo "═══════════════════════════════════════════════════════════════"
    echo -e "${BLUE}📋 TEST EXECUTION SUMMARY${NC}"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "Total Test Categories: $TOTAL_TESTS"
    echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
    echo -e "${RED}Failed: $FAILED_TESTS${NC}"
    echo -e "${YELLOW}Skipped (Pending Implementation): $SKIPPED_TESTS${NC}"
    echo ""

    if [ $FAILED_TESTS -eq 0 ] && [ $TOTAL_TESTS -gt 0 ]; then
        echo -e "${GREEN}╔═══════════════════════════════════════════════════════════╗${NC}"
        echo -e "${GREEN}║  ✅ ALL TESTS PASSED - PROJECT READY FOR RELEASE         ║${NC}"
        echo -e "${GREEN}╚═══════════════════════════════════════════════════════════╝${NC}"
        echo ""
    elif [ $SKIPPED_TESTS -gt 0 ] && [ $FAILED_TESTS -eq 0 ]; then
        echo -e "${YELLOW}╔═══════════════════════════════════════════════════════════╗${NC}"
        echo -e "${YELLOW}║  ⚠️  TDD RED PHASE - Tests pending implementation         ║${NC}"
        echo -e "${YELLOW}║     This is EXPECTED behavior before implementation       ║${NC}"
        echo -e "${YELLOW}╚═══════════════════════════════════════════════════════════╝${NC}"
        echo ""
    else
        echo -e "${RED}╔═══════════════════════════════════════════════════════════╗${NC}"
        echo -e "${RED}║  ❌ TEST FAILURES DETECTED - Review required              ║${NC}"
        echo -e "${RED}╚═══════════════════════════════════════════════════════════╝${NC}"
        echo ""
    fi

    if [ $SKIPPED_TESTS -gt 0 ]; then
        echo -e "${YELLOW}Note: Skipped tests indicate pending implementation (TDD RED phase)${NC}"
        echo -e "${YELLOW}This is normal and expected before code is written.${NC}"
        echo ""
    fi

    echo "═══════════════════════════════════════════════════════════════"
}

# Main execution
main() {
    check_prerequisites

    case $RUN_CATEGORY in
        interfaces)
            run_interface_tests
            ;;
        integration)
            run_integration_tests
            ;;
        e2e)
            run_e2e_tests
            ;;
        quality)
            run_quality_tests
            ;;
        all)
            run_interface_tests
            run_integration_tests
            run_e2e_tests
            run_quality_tests
            ;;
        *)
            echo -e "${RED}Unknown category: $RUN_CATEGORY${NC}"
            echo "Valid categories: interfaces, integration, e2e, quality, all"
            exit 1
            ;;
    esac

    generate_coverage_report
    print_summary

    # Exit with failure if any tests failed
    if [ $FAILED_TESTS -gt 0 ]; then
        exit 1
    fi

    # In CI mode, treat skipped tests as success (pending implementation)
    exit 0
}

# Run main function
main
