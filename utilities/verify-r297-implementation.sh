#!/bin/bash

# R297 Implementation Verification Script
# Ensures architect split detection protocol is properly implemented

echo "═══════════════════════════════════════════════════════════════"
echo "🔍 R297 IMPLEMENTATION VERIFICATION"
echo "═══════════════════════════════════════════════════════════════"
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo ""

ERRORS=0
WARNINGS=0

# Check 1: R297 rule file exists
echo "📋 Check 1: R297 rule file exists"
if [ -f "rule-library/R297-architect-split-detection-protocol.md" ]; then
    echo "   ✅ R297 rule file found"
else
    echo "   ❌ R297 rule file missing!"
    ((ERRORS++))
fi
echo ""

# Check 2: R297 in RULE-REGISTRY
echo "📋 Check 2: R297 registered in RULE-REGISTRY"
if grep -q "R297" rule-library/RULE-REGISTRY.md; then
    echo "   ✅ R297 found in registry"
else
    echo "   ❌ R297 not in registry!"
    ((ERRORS++))
fi
echo ""

# Check 3: R297 in architect config
echo "📋 Check 3: R297 in architect configuration"
if grep -q "R297" .claude/agents/architect.md; then
    echo "   ✅ R297 referenced in architect config"
else
    echo "   ❌ R297 missing from architect config!"
    ((ERRORS++))
fi
echo ""

# Check 4: R297 in all architect states
echo "📋 Check 4: R297 in architect state rules"
STATES=("WAVE_REVIEW" "PHASE_ASSESSMENT" "INTEGRATION_REVIEW")
for state in "${STATES[@]}"; do
    if grep -q "R297" "agent-states/architect/${state}/rules.md" 2>/dev/null; then
        echo "   ✅ R297 found in ${state} state"
    else
        echo "   ⚠️ R297 missing from ${state} state"
        ((WARNINGS++))
    fi
done
echo ""

# Check 5: R022 references R297
echo "📋 Check 5: R022 (size verification) references R297"
if grep -q "R297" rule-library/R022-architect-size-verification.md; then
    echo "   ✅ R022 properly references R297"
else
    echo "   ❌ R022 doesn't reference R297!"
    ((ERRORS++))
fi
echo ""

# Check 6: Split detection logic
echo "📋 Check 6: Split detection implementation"
echo "   Checking for key patterns in R297..."

PATTERNS=(
    "split_count"
    "orchestrator-state.json"
    "[Oo]riginal effort branches"
    "NOT integration branches"
    "PRs come from effort branches"
)

for pattern in "${PATTERNS[@]}"; do
    if grep -q "$pattern" rule-library/R297-architect-split-detection-protocol.md; then
        echo "   ✅ Found: '$pattern'"
    else
        echo "   ❌ Missing: '$pattern'"
        ((ERRORS++))
    fi
done
echo ""

# Check 7: Integration clarification
echo "📋 Check 7: Integration branch clarification"
if grep -q "Integration branches merge all splits" rule-library/R297-architect-split-detection-protocol.md; then
    echo "   ✅ Integration behavior clearly explained"
else
    echo "   ⚠️ Integration behavior could be clearer"
    ((WARNINGS++))
fi
echo ""

# Summary
echo "═══════════════════════════════════════════════════════════════"
echo "📊 VERIFICATION SUMMARY"
echo "═══════════════════════════════════════════════════════════════"
echo "   Errors: $ERRORS"
echo "   Warnings: $WARNINGS"
echo ""

if [ $ERRORS -eq 0 ]; then
    if [ $WARNINGS -eq 0 ]; then
        echo "✅ R297 FULLY IMPLEMENTED - All checks passed!"
    else
        echo "✅ R297 IMPLEMENTED - Minor warnings to address"
    fi
    echo ""
    echo "The architect will now:"
    echo "1. Check split_count BEFORE measuring any effort"
    echo "2. Measure ORIGINAL effort branches (not integration)"
    echo "3. Recognize already-split efforts as compliant"
    echo "4. Understand integration branches exceed limits by design"
    exit 0
else
    echo "❌ R297 IMPLEMENTATION INCOMPLETE - Critical errors found!"
    echo ""
    echo "Fix the errors above to ensure architects:"
    echo "- Don't demand re-splits of already-split efforts"
    echo "- Measure the correct branches for compliance"
    exit 1
fi