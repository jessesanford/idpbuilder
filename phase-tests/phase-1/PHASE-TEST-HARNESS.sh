#!/bin/bash

# PHASE 1 TEST HARNESS - COMMAND SKELETON & FOUNDATION
# Project: idpbuilder-push
# Methodology: Test-Driven Development (TDD)
# Generated: 2025-09-22
#
# This harness runs all Phase 1 tests to verify command skeleton functionality.
# Tests will initially FAIL (no implementation) with helpful error messages.

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_DIR="./tests/phase1"
PROJECT_ROOT="/home/vscode/workspaces/idpbuilder-push"
COVERAGE_DIR="./coverage"
TIMESTAMP=$(date '+%Y%m%d-%H%M%S')

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
SKIPPED_TESTS=0

# Function to print colored output
print_color() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Function to print section header
print_header() {
    local title=$1
    echo ""
    print_color "$BLUE" "==============================================="
    print_color "$BLUE" "$title"
    print_color "$BLUE" "==============================================="
}

# Function to run a test category
run_test_category() {
    local category=$1
    local test_file=$2
    local description=$3

    print_color "$YELLOW" "Running: $description"

    # Check if test file exists
    if [ ! -f "$TEST_DIR/$test_file" ]; then
        print_color "$RED" "  ✗ Test file not found: $test_file"
        print_color "$RED" "    Expected at: $TEST_DIR/$test_file"
        print_color "$RED" "    This is expected before implementation begins (TDD)"
        echo ""
        ((SKIPPED_TESTS++))
        return 1
    fi

    # Run the test
    if go test "$TEST_DIR/$test_file" -v -count=1 2>&1 | tee -a "test-output-$TIMESTAMP.log"; then
        print_color "$GREEN" "  ✓ $category tests passed"
        ((PASSED_TESTS++))
    else
        print_color "$RED" "  ✗ $category tests failed"
        print_color "$YELLOW" "    This is expected in TDD - tests fail before implementation"
        ((FAILED_TESTS++))
    fi

    echo ""
    return 0
}

# Function to generate coverage report
generate_coverage_report() {
    print_header "GENERATING COVERAGE REPORT"

    # Create coverage directory
    mkdir -p "$COVERAGE_DIR"

    # Check if any test files exist
    if ! ls "$TEST_DIR"/*.go 2>/dev/null | grep -q .; then
        print_color "$YELLOW" "No test files found yet. This is expected before implementation."
        return
    fi

    # Generate coverage
    if go test "$TEST_DIR/..." -coverprofile="$COVERAGE_DIR/phase1-coverage.out" 2>/dev/null; then
        go tool cover -html="$COVERAGE_DIR/phase1-coverage.out" -o "$COVERAGE_DIR/phase1-coverage.html"
        print_color "$GREEN" "Coverage report generated: $COVERAGE_DIR/phase1-coverage.html"
    else
        print_color "$YELLOW" "Coverage report will be available once tests can compile"
    fi
}

# Function to check prerequisites
check_prerequisites() {
    print_header "CHECKING PREREQUISITES"

    # Check Go installation
    if ! command -v go &> /dev/null; then
        print_color "$RED" "Go is not installed. Please install Go 1.22+"
        exit 1
    fi

    GO_VERSION=$(go version | awk '{print $3}')
    print_color "$GREEN" "✓ Go installed: $GO_VERSION"

    # Check project structure
    if [ ! -d "$PROJECT_ROOT" ]; then
        print_color "$RED" "Project root not found: $PROJECT_ROOT"
        exit 1
    fi
    print_color "$GREEN" "✓ Project root exists: $PROJECT_ROOT"

    # Create test directory if it doesn't exist
    if [ ! -d "$TEST_DIR" ]; then
        mkdir -p "$TEST_DIR"
        print_color "$YELLOW" "Created test directory: $TEST_DIR"
    else
        print_color "$GREEN" "✓ Test directory exists: $TEST_DIR"
    fi
}

# Function to run demo validation
run_demo_validation() {
    print_header "DEMO SCENARIO VALIDATION"

    if [ -f "./PHASE-DEMO-PLAN.md" ]; then
        print_color "$GREEN" "✓ Demo plan exists"

        # Check if demo scripts are executable
        if [ -d "./demos" ]; then
            for demo in ./demos/*.sh; do
                if [ -f "$demo" ]; then
                    if [ -x "$demo" ]; then
                        print_color "$GREEN" "  ✓ $(basename $demo) is executable"
                    else
                        print_color "$YELLOW" "  ! $(basename $demo) is not executable"
                        chmod +x "$demo"
                        print_color "$GREEN" "    Fixed: made executable"
                    fi
                fi
            done
        else
            print_color "$YELLOW" "Demo scripts directory not created yet"
        fi
    else
        print_color "$YELLOW" "Demo plan not found. Expected at: ./PHASE-DEMO-PLAN.md"
    fi
}

# Function to print summary
print_summary() {
    print_header "TEST EXECUTION SUMMARY"

    echo "Test Results:"
    print_color "$GREEN" "  Passed:  $PASSED_TESTS"
    print_color "$RED" "  Failed:  $FAILED_TESTS"
    print_color "$YELLOW" "  Skipped: $SKIPPED_TESTS"
    echo ""

    TOTAL_TESTS=$((PASSED_TESTS + FAILED_TESTS + SKIPPED_TESTS))
    echo "Total Tests Categories: $TOTAL_TESTS"

    if [ $FAILED_TESTS -gt 0 ] || [ $SKIPPED_TESTS -gt 0 ]; then
        echo ""
        print_color "$YELLOW" "NOTE: Test failures and skipped tests are EXPECTED in TDD!"
        print_color "$YELLOW" "Tests should fail before implementation begins."
        print_color "$YELLOW" "These failures guide the implementation process."
    fi

    if [ $PASSED_TESTS -eq $TOTAL_TESTS ] && [ $TOTAL_TESTS -gt 0 ]; then
        echo ""
        print_color "$GREEN" "🎉 All tests passing! Phase 1 implementation is complete!"
    fi

    echo ""
    echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "Log file: test-output-$TIMESTAMP.log"
}

# Function to show implementation hints
show_implementation_hints() {
    print_header "IMPLEMENTATION HINTS"

    print_color "$BLUE" "Next Steps for Implementation:"
    echo "1. Start with Wave 1.1, Effort 1.1.1: Command Tests"
    echo "2. Run this harness to see failing tests"
    echo "3. Implement minimal code to make tests pass"
    echo "4. Refactor once tests are green"
    echo ""

    print_color "$BLUE" "TDD Workflow:"
    echo "  RED:      Write/run failing tests (current state)"
    echo "  GREEN:    Write minimal code to pass"
    echo "  REFACTOR: Improve code quality"
    echo ""

    print_color "$BLUE" "Files to Create (in order):"
    echo "  1. cmd/push/root.go         - Command definition"
    echo "  2. cmd/push/config.go       - Configuration structures"
    echo "  3. pkg/oci/validation.go    - Input validation"
    echo "  4. pkg/oci/errors.go        - Error handling"
}

# Main execution
main() {
    print_header "PHASE 1 TEST HARNESS - TDD EXECUTION"
    print_color "$YELLOW" "Project: idpbuilder-push"
    print_color "$YELLOW" "Phase: 1 - Command Skeleton & Foundation"
    echo ""

    # Check prerequisites
    check_prerequisites

    # Change to project root
    cd "$PROJECT_ROOT/phase-tests/phase-1" || exit 1

    print_header "RUNNING PHASE 1 TEST SUITE"

    # Run test categories in order
    run_test_category "Command Registration" "command_test.go" "Command registration and existence tests"
    run_test_category "Flag Definition" "flags_test.go" "Command flag definition and parsing tests"
    run_test_category "Argument Validation" "validation_test.go" "Argument validation and error handling tests"
    run_test_category "Help Text" "help_test.go" "Help text and documentation tests"
    run_test_category "Integration" "integration_test.go" "CLI integration and environment tests"

    # Generate coverage report
    generate_coverage_report

    # Run demo validation
    run_demo_validation

    # Print summary
    print_summary

    # Show implementation hints
    show_implementation_hints
}

# Handle script arguments
case "${1:-}" in
    --help|-h)
        echo "Phase 1 Test Harness for idpbuilder-push"
        echo ""
        echo "Usage: $0 [OPTION]"
        echo ""
        echo "Options:"
        echo "  --help, -h     Show this help message"
        echo "  --verbose, -v  Run with verbose output"
        echo "  --coverage     Only generate coverage report"
        echo "  --quick        Run tests without coverage"
        echo ""
        echo "This harness runs all Phase 1 tests following TDD methodology."
        echo "Tests should initially fail (RED phase) before implementation."
        exit 0
        ;;
    --coverage)
        cd "$PROJECT_ROOT/phase-tests/phase-1" || exit 1
        generate_coverage_report
        exit 0
        ;;
    --quick)
        # Set flag to skip coverage
        SKIP_COVERAGE=true
        main
        ;;
    --verbose|-v)
        set -x
        main
        ;;
    *)
        main
        ;;
esac