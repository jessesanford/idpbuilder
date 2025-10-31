#!/bin/bash
# R420 Compliance Validation Script
# Validates that implementation/split plans include required R420 research

set -e

PLAN_FILE="$1"

if [ -z "$PLAN_FILE" ]; then
    echo "Usage: $0 <plan-file.md>"
    echo ""
    echo "Validates R420 Cross-Effort Planning Awareness compliance in plan files."
    exit 1
fi

if [ ! -f "$PLAN_FILE" ]; then
    echo "❌ ERROR: Plan file not found: $PLAN_FILE"
    exit 1
fi

echo "==========================================="
echo "R420 COMPLIANCE VALIDATION"
echo "==========================================="
echo "Plan file: $PLAN_FILE"
echo ""

VALIDATION_FAILED=0

# Check for mandatory R420 section
echo "Checking for R420 Prior Work Analysis section..."
if ! grep -q "PRIOR WORK ANALYSIS (R420 MANDATORY)" "$PLAN_FILE"; then
    echo "❌ BLOCKING: Missing mandatory 'PRIOR WORK ANALYSIS (R420 MANDATORY)' section"
    echo "   This section is REQUIRED by R420 to prevent integration failures"
    VALIDATION_FAILED=1
else
    echo "✅ Found R420 section header"
fi

# Check for required subsections
echo ""
echo "Checking for required subsections..."

REQUIRED_SECTIONS=(
    "Discovery Phase Results"
    "File Structure Findings"
    "Interface/API Findings"
    "Conflicts Detected"
    "Required Integrations"
    "Forbidden Actions"
)

for section in "${REQUIRED_SECTIONS[@]}"; do
    if ! grep -q "$section" "$PLAN_FILE"; then
        echo "❌ BLOCKING: Missing required section: '$section'"
        VALIDATION_FAILED=1
    else
        echo "✅ Found section: $section"
    fi
done

# Check for substantive content (not just templates)
echo ""
echo "Checking for substantive research content..."

# Count finding rows in tables
FINDINGS_COUNT=$(grep -c "| .* | .* |" "$PLAN_FILE" 2>/dev/null || echo "0")
if [ "$FINDINGS_COUNT" -lt 3 ]; then
    echo "⚠️ WARNING: Research may be incomplete"
    echo "   Found only $FINDINGS_COUNT table rows, expected at least 3"
    echo "   This suggests research was not thorough"
    VALIDATION_FAILED=1
else
    echo "✅ Found $FINDINGS_COUNT table rows (substantive research)"
fi

# Check that research timestamp is present
echo ""
echo "Checking for research timestamp..."
if ! grep -q "Research Timestamp" "$PLAN_FILE"; then
    echo "⚠️ WARNING: No research timestamp found"
    echo "   Cannot verify when research was conducted"
else
    TIMESTAMP_LINE=$(grep "Research Timestamp" "$PLAN_FILE")
    echo "✅ $TIMESTAMP_LINE"
fi

# Check for split-specific requirements (if this is a split plan)
if echo "$PLAN_FILE" | grep -qi "split"; then
    echo ""
    echo "Detected split plan - checking split-specific requirements..."

    if ! grep -q "API Compatibility Findings" "$PLAN_FILE"; then
        echo "❌ BLOCKING: Split plan missing 'API Compatibility Findings'"
        echo "   API compatibility is CRITICAL for splits (prevents retry.DefaultBackoff errors)"
        VALIDATION_FAILED=1
    else
        echo "✅ Found API Compatibility Findings"
    fi

    if ! grep -q "Method Visibility Findings" "$PLAN_FILE"; then
        echo "❌ BLOCKING: Split plan missing 'Method Visibility Findings'"
        echo "   Method visibility is CRITICAL (prevents unexported method access)"
        VALIDATION_FAILED=1
    else
        echo "✅ Found Method Visibility Findings"
    fi

    if ! grep -q "API Assumptions Verified" "$PLAN_FILE"; then
        echo "❌ BLOCKING: Split plan missing 'API Assumptions Verified'"
        echo "   API verification is MANDATORY to prevent assumptions"
        VALIDATION_FAILED=1
    else
        echo "✅ Found API Assumptions Verified"
    fi
fi

# Check that conflicts section is populated
echo ""
echo "Checking conflicts analysis..."
if grep -q "Conflicts Detected" "$PLAN_FILE"; then
    # Check if conflicts are documented or if "NO conflicts" is stated
    if grep -A5 "Conflicts Detected" "$PLAN_FILE" | grep -q "NO.*detected\|✅"; then
        echo "✅ Conflicts analysis performed (no conflicts found)"
    elif grep -A5 "Conflicts Detected" "$PLAN_FILE" | grep -q "❌\|CONFLICT"; then
        echo "⚠️ CONFLICTS FOUND - These MUST be resolved before implementation:"
        grep -A10 "Conflicts Detected" "$PLAN_FILE" | grep "❌\|CONFLICT"
    else
        echo "⚠️ WARNING: Conflicts section present but appears empty"
        echo "   R420 requires explicit conflict analysis"
    fi
fi

# Check that forbidden actions are specified
echo ""
echo "Checking forbidden actions..."
if grep -A5 "Forbidden Actions" "$PLAN_FILE" | grep -q "❌.*DO NOT"; then
    FORBIDDEN_COUNT=$(grep -A20 "Forbidden Actions" "$PLAN_FILE" | grep -c "❌.*DO NOT" || echo "0")
    echo "✅ Found $FORBIDDEN_COUNT forbidden action(s) documented"
else
    echo "⚠️ WARNING: No forbidden actions specified"
    echo "   R420 requires documenting what NOT to do based on prior work"
fi

# Final validation result
echo ""
echo "==========================================="
if [ $VALIDATION_FAILED -eq 0 ]; then
    echo "✅ PLAN PASSES R420 VALIDATION"
    echo "==========================================="
    echo ""
    echo "Plan includes required R420 research and is ready for approval."
    exit 0
else
    echo "❌ PLAN FAILS R420 VALIDATION"
    echo "==========================================="
    echo ""
    echo "BLOCKING: Plan cannot be approved until R420 research is complete."
    echo ""
    echo "Required actions:"
    echo "1. Execute R420 research protocol (see state rules)"
    echo "2. Add missing sections to plan"
    echo "3. Ensure all tables have substantive content"
    echo "4. Re-run this validation script"
    echo ""
    echo "See: rule-library/R420-cross-effort-planning-awareness-protocol.md"
    echo "See: agent-states/code-reviewer/EFFORT_PLAN_CREATION/rules.md"
    echo "See: agent-states/code-reviewer/CREATE_SPLIT_PLAN/rules.md"
    exit 1
fi
