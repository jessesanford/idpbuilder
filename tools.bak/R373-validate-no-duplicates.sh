#!/bin/bash
# R373-validate-no-duplicates.sh
# Validates that no duplicate interfaces or incompatible implementations exist
# This tool enforces R373: Mandatory Code Reuse and Interface Compliance

set -e

EFFORT_REPO="${1:-$(pwd)}"
PHASE="${2:-1}"
WAVE="${3:-1}"

echo "🔴🔴🔴 R373: Interface Compliance Validation 🔴🔴🔴"
echo "Repository: $EFFORT_REPO"
echo "Phase: $PHASE, Wave: $WAVE"
echo "================================================"

cd "$EFFORT_REPO"

# Track violations
VIOLATIONS=0
VIOLATION_REPORT=""

# Function to add violation
add_violation() {
    local severity="$1"
    local message="$2"
    VIOLATIONS=$((VIOLATIONS + 1))
    VIOLATION_REPORT="${VIOLATION_REPORT}\n${severity}: ${message}"
    echo "${severity}: ${message}"
}

# 1. Check for duplicate interface definitions
echo ""
echo "=== Checking for Duplicate Interface Definitions ==="
INTERFACES=$(find . -name "*.go" -exec grep -h "^type.*interface" {} \; 2>/dev/null | sort)

if [ ! -z "$INTERFACES" ]; then
    DUPLICATES=$(echo "$INTERFACES" | uniq -d)
    if [ ! -z "$DUPLICATES" ]; then
        add_violation "❌ CRITICAL" "Found duplicate interface definitions:"
        echo "$DUPLICATES"
    else
        echo "✅ No duplicate interface definitions found"
    fi
fi

# 2. Check for methods with same name but different signatures
echo ""
echo "=== Checking for Incompatible Method Signatures ==="

for method in Push Pull Upload Download Store Retrieve Create Delete Get List Update; do
    echo -n "Checking ${method}... "

    # Find all method signatures
    SIGNATURES=$(grep -r "func.*${method}(" --include="*.go" 2>/dev/null | \
        sed "s/.*func.*${method}//" | \
        sed 's/^\s*//' | \
        sort | uniq)

    UNIQUE_COUNT=$(echo "$SIGNATURES" | grep -v "^$" | wc -l)

    if [ "$UNIQUE_COUNT" -gt 1 ]; then
        add_violation "❌ WARNING" "Method '${method}' has $UNIQUE_COUNT different signatures:"
        echo "$SIGNATURES" | while read sig; do
            echo "    - ${method}${sig}"
        done
    elif [ "$UNIQUE_COUNT" -eq 1 ]; then
        echo "✅ Consistent (1 signature)"
    else
        echo "✅ Not found"
    fi
done

# 3. Check for interfaces with similar names (potential duplicates)
echo ""
echo "=== Checking for Similar Interface Names ==="

INTERFACE_NAMES=$(find . -name "*.go" -exec grep "^type.*interface" {} \; 2>/dev/null | \
    awk '{print $2}' | sort | uniq)

for name in Registry Storage Client Builder Processor Handler Manager Service; do
    SIMILAR=$(echo "$INTERFACE_NAMES" | grep -i "$name" || true)
    COUNT=$(echo "$SIMILAR" | grep -v "^$" | wc -l)

    if [ "$COUNT" -gt 1 ]; then
        add_violation "⚠️ WARNING" "Multiple interfaces containing '$name':"
        echo "$SIMILAR" | while read iface; do
            echo "    - $iface"
        done
    fi
done

# 4. Check for potential stub implementations
echo ""
echo "=== Checking for Stub Implementations (R355 Compliance) ==="

STUB_PATTERNS="TODO\|FIXME\|not implemented\|return nil.*// implement later\|panic.*not implemented"
STUBS=$(grep -r "$STUB_PATTERNS" --include="*.go" --exclude-dir=test 2>/dev/null || true)

if [ ! -z "$STUBS" ]; then
    add_violation "❌ CRITICAL" "Found stub implementations (R355 violation):"
    echo "$STUBS" | head -5
fi

# 5. Cross-branch interface comparison (if multiple branches exist)
echo ""
echo "=== Checking Cross-Branch Interface Consistency ==="

# Get all effort branches for this wave
BRANCHES=$(git branch -r 2>/dev/null | grep "phase${PHASE}.*wave${WAVE}" | sed 's/origin\///' || true)
BRANCH_COUNT=$(echo "$BRANCHES" | grep -v "^$" | wc -l)

if [ "$BRANCH_COUNT" -gt 1 ]; then
    echo "Found $BRANCH_COUNT branches to compare"

    # Create temporary directory for comparison
    TEMP_DIR=$(mktemp -d)

    # Extract interfaces from each branch
    for branch in $BRANCHES; do
        BRANCH_NAME=$(basename "$branch")
        echo -n "  Checking $BRANCH_NAME... "

        git checkout "$branch" 2>/dev/null || continue

        # Extract all interfaces
        find . -name "*.go" -exec grep -h "type.*interface {" {} \; 2>/dev/null | \
            sed 's/type //' | sed 's/ interface.*//' | sort > "$TEMP_DIR/$BRANCH_NAME.interfaces"

        # Extract all method signatures
        grep -r "^\s*[A-Z].*(" --include="*.go" 2>/dev/null | \
            grep -v "//" | \
            sed 's/.*:\s*//' | sort > "$TEMP_DIR/$BRANCH_NAME.methods"

        echo "✓"
    done

    # Compare interfaces across branches
    echo ""
    echo "  Interface comparison:"
    for file1 in "$TEMP_DIR"/*.interfaces; do
        for file2 in "$TEMP_DIR"/*.interfaces; do
            if [ "$file1" != "$file2" ]; then
                BRANCH1=$(basename "$file1" .interfaces)
                BRANCH2=$(basename "$file2" .interfaces)

                # Find interfaces in branch1 not in branch2
                MISSING=$(comm -23 "$file1" "$file2" 2>/dev/null | head -3)
                if [ ! -z "$MISSING" ]; then
                    echo "    Interfaces in $BRANCH1 but not $BRANCH2:"
                    echo "$MISSING" | sed 's/^/      - /'
                fi
            fi
        done
    done

    # Clean up
    rm -rf "$TEMP_DIR"
else
    echo "Single branch found - skipping cross-branch comparison"
fi

# 6. Generate compliance report
echo ""
echo "================================================"
echo "R373 COMPLIANCE REPORT"
echo "================================================"

if [ "$VIOLATIONS" -eq 0 ]; then
    echo "✅ PASSED: No R373 violations detected"
    echo "All interfaces and implementations are properly aligned"
    exit 0
else
    echo "❌ FAILED: Found $VIOLATIONS violations"
    echo ""
    echo "Violations Summary:"
    echo -e "$VIOLATION_REPORT"
    echo ""
    echo "REQUIRED ACTIONS:"
    echo "1. Refactor duplicate interfaces to use single definition"
    echo "2. Align all method signatures to match existing interfaces"
    echo "3. Remove stub implementations (R355 compliance)"
    echo "4. Import and reuse existing code instead of reimplementing"
    echo ""
    echo "See: rule-library/R373-mandatory-code-reuse-and-interface-compliance.md"
    exit 373
fi