#!/bin/bash
# tools/validate-test-fixtures.sh
# Validation script for test fixtures to prevent hardcoded values
# Related to: Item #588, BLOCKER-20251014-180000 Phase 2 P1

set -euo pipefail

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Get the git root directory
GIT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || echo ".")"

# Function to print colored messages
print_error() {
    echo -e "${RED}${BOLD}❌ ERROR:${NC} ${RED}$1${NC}" >&2
}

print_success() {
    echo -e "${GREEN}${BOLD}✅ PROJECT_DONE:${NC} ${GREEN}$1${NC}"
}

print_info() {
    echo -e "${BLUE}${BOLD}ℹ️  INFO:${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}${BOLD}⚠️  WARNING:${NC} ${YELLOW}$1${NC}"
}

# Function to check for hardcoded values in test files
check_hardcoded_values() {
    local violations=0
    local test_dir="$GIT_ROOT/tests"

    # Hardcoded values to search for
    local hardcoded_patterns=(
        "hello-world-fullstack"
        "https://github.com/test/hello-world-fullstack.git"
        "git@github.com:test/hello-world-fullstack.git"
    )

    print_info "Checking for hardcoded test values in $test_dir..."
    echo ""

    # Check each pattern
    for pattern in "${hardcoded_patterns[@]}"; do
        # Search for pattern in test files (excluding fixtures directory for now)
        local matches=$(grep -rn "$pattern" "$test_dir" \
            --include="*.sh" \
            --include="*.bash" \
            --exclude-dir="fixtures" \
            2>/dev/null || true)

        if [ -n "$matches" ]; then
            violations=$((violations + 1))
            print_error "Found hardcoded value: \"$pattern\""
            echo ""
            echo "$matches" | while IFS= read -r match; do
                # Extract file, line number, and context
                local file=$(echo "$match" | cut -d: -f1)
                local line_num=$(echo "$match" | cut -d: -f2)
                local content=$(echo "$match" | cut -d: -f3-)

                echo -e "  ${RED}File:${NC} $(basename "$file"):${line_num}"
                echo -e "  ${RED}Code:${NC} $content"
                echo ""
            done

            # Suggest replacement
            case "$pattern" in
                "hello-world-fullstack")
                    echo -e "  ${YELLOW}Suggested fix:${NC} Replace with \$PROJECT_PREFIX variable"
                    echo -e "  ${YELLOW}Example:${NC} \"project_prefix\": \"\$PROJECT_PREFIX\""
                    ;;
                "https://github.com/test/hello-world-fullstack.git")
                    echo -e "  ${YELLOW}Suggested fix:${NC} Replace with \$TARGET_REPO variable"
                    echo -e "  ${YELLOW}Example:${NC} \"target_repo_url\": \"\$TARGET_REPO\""
                    ;;
            esac
            echo ""
        fi
    done

    # Also check for hardcoded paths in fixtures
    local fixture_dir="$test_dir/fixtures"
    if [ -d "$fixture_dir" ]; then
        print_info "Checking fixtures directory for hardcoded values..."

        local fixture_violations=$(grep -rn "hello-world-fullstack" "$fixture_dir" \
            --include="*.json" \
            2>/dev/null || true)

        if [ -n "$fixture_violations" ]; then
            print_warning "Found hardcoded values in fixture files"
            echo ""
            echo -e "${YELLOW}NOTE:${NC} Fixtures should use template system or PLACEHOLDER_UUID pattern"
            echo -e "${YELLOW}See:${NC} tests/fixtures/templates/README.md for template usage"
            echo ""
            echo "$fixture_violations" | while IFS= read -r match; do
                local file=$(echo "$match" | cut -d: -f1)
                local line_num=$(echo "$match" | cut -d: -f2)
                echo -e "  ${YELLOW}File:${NC} $(basename "$file"):${line_num}"
            done
            echo ""
            violations=$((violations + 1))
        fi
    fi

    return $violations
}

# Function to validate staged test files
validate_staged_files() {
    local violations=0

    # Get list of staged test files
    local staged_tests=$(git diff --cached --name-only --diff-filter=ACM | grep "^tests/" | grep -E '\.(sh|bash|json)$' || true)

    if [ -z "$staged_tests" ]; then
        print_info "No staged test files to validate"
        return 0
    fi

    print_info "Validating $(echo "$staged_tests" | wc -l) staged test file(s)..."
    echo ""

    # Check each staged file for hardcoded patterns
    while IFS= read -r file; do
        if [ -f "$GIT_ROOT/$file" ]; then
            local file_violations=$(grep -n "hello-world-fullstack\|github.com/test/hello-world-fullstack" "$GIT_ROOT/$file" || true)

            if [ -n "$file_violations" ]; then
                violations=$((violations + 1))
                print_error "Hardcoded values found in staged file: $file"
                echo ""
                echo "$file_violations" | while IFS= read -r match; do
                    local line_num=$(echo "$match" | cut -d: -f1)
                    local content=$(echo "$match" | cut -d: -f2-)
                    echo -e "  ${RED}Line $line_num:${NC} $content"
                done
                echo ""
            fi
        fi
    done <<< "$staged_tests"

    return $violations
}

# Function to suggest using template system
suggest_template_usage() {
    echo ""
    echo -e "${BLUE}${BOLD}📝 TEST FIXTURE BEST PRACTICES:${NC}"
    echo ""
    echo "1. Use the fixture template system for dynamic values:"
    echo "   - Template files: tests/fixtures/templates/*.json.template"
    echo "   - Use tokens: {{PROJECT_PREFIX}}, {{TARGET_REPO}}, {{TIMESTAMP}}"
    echo "   - Generate at runtime: tests/lib/fixture-template-engine.sh"
    echo ""
    echo "2. For inline fixtures, use framework variables:"
    echo "   - \$PROJECT_PREFIX (unique per test run)"
    echo "   - \$TARGET_REPO (configured repository URL)"
    echo "   - \$TEST_WORKSPACE (isolated test directory)"
    echo ""
    echo "3. For static fixtures, use placeholder pattern:"
    echo "   - \"project_prefix\": \"PLACEHOLDER_UUID\""
    echo "   - Replace at runtime with framework substitution"
    echo ""
    echo "See: tests/fixtures/templates/README.md for complete documentation"
    echo ""
}

# Main validation function
main() {
    local exit_code=0
    local mode="${1:-all}"

    echo ""
    echo -e "${BOLD}===========================================${NC}"
    echo -e "${BOLD}Test Fixture Validation${NC}"
    echo -e "${BOLD}===========================================${NC}"
    echo ""

    case "$mode" in
        "staged")
            # Validate only staged files (for pre-commit hook)
            if ! validate_staged_files; then
                print_error "Hardcoded values found in staged test files"
                suggest_template_usage
                exit_code=1
            else
                print_success "All staged test files validated"
            fi
            ;;

        "all")
            # Validate all test files (for manual checks)
            if ! check_hardcoded_values; then
                print_error "Hardcoded values found in test files"
                suggest_template_usage
                exit_code=1
            else
                print_success "No hardcoded values found in test files"
            fi
            ;;

        "help")
            echo "Usage: $0 [MODE]"
            echo ""
            echo "Modes:"
            echo "  all     - Check all test files for hardcoded values (default)"
            echo "  staged  - Check only staged test files (for pre-commit hook)"
            echo "  help    - Show this help message"
            echo ""
            echo "Purpose:"
            echo "  Prevents hardcoded project names and repository URLs in test files"
            echo "  Enforces use of framework variables (\$PROJECT_PREFIX, \$TARGET_REPO)"
            echo "  Ensures test isolation and concurrent execution compatibility"
            echo ""
            echo "Examples:"
            echo "  $0              # Check all test files"
            echo "  $0 staged       # Check staged files (pre-commit)"
            echo ""
            exit 0
            ;;

        *)
            print_error "Unknown mode: $mode"
            echo "Run '$0 help' for usage information"
            exit 1
            ;;
    esac

    echo ""
    echo -e "${BOLD}===========================================${NC}"

    if [ $exit_code -eq 0 ]; then
        print_success "Test fixture validation passed!"
        echo ""
    else
        print_error "Test fixture validation failed!"
        echo ""
        echo -e "${RED}${BOLD}⚠️  CRITICAL:${NC} Test isolation requires dynamic values"
        echo -e "${RED}Fix these violations before committing${NC}"
        echo ""
    fi

    exit $exit_code
}

# Run main function
main "$@"
