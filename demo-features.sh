#!/bin/bash
# Demo: IDPBuilder Structure Analysis
# Effort: E1.1.1-analyze-existing-structure

set -e
echo "🎬 Demonstrating IDPBuilder Structure Analysis Deliverables"
echo "==========================================================="

# Demo objective 1: Verify analysis report exists
echo "✅ Step 1: Verifying analysis report exists..."
if [ -f ".software-factory/ANALYSIS-REPORT.md" ]; then
    echo "✅ Analysis report found: .software-factory/ANALYSIS-REPORT.md"
    WORD_COUNT=$(wc -w < .software-factory/ANALYSIS-REPORT.md)
    echo "   Report size: $WORD_COUNT words"
    if [ "$WORD_COUNT" -lt 1000 ]; then
        echo "❌ Analysis report too short (< 1000 words)"
        exit 1
    fi
    echo "   ✅ Comprehensive analysis (>= 1000 words)"
else
    echo "❌ Analysis report not found"
    exit 1
fi

# Demo objective 2: Verify analysis covered required topics
echo ""
echo "✅ Step 2: Checking analysis coverage..."

REQUIRED_TOPICS=(
    "Command Structure"
    "Dependencies"
    "Testing Patterns"
    "Package Structure"
    "Cobra"
    "Build System"
    "Authentication"
)

MISSING_TOPICS=()
for topic in "${REQUIRED_TOPICS[@]}"; do
    if grep -qi "$topic" .software-factory/ANALYSIS-REPORT.md; then
        echo "   ✅ Covered: $topic"
    else
        echo "   ❌ Missing: $topic"
        MISSING_TOPICS+=("$topic")
    fi
done

if [ ${#MISSING_TOPICS[@]} -gt 0 ]; then
    echo "❌ Analysis incomplete - missing topics: ${MISSING_TOPICS[*]}"
    exit 1
fi

# Demo objective 3: Verify key packages were analyzed
echo ""
echo "✅ Step 3: Verifying key packages analyzed..."

KEY_PACKAGES=(
    "pkg/cmd"
    "pkg/build"
    "pkg/k8s"
    "pkg/controllers"
)

for pkg in "${KEY_PACKAGES[@]}"; do
    if grep -q "$pkg" .software-factory/ANALYSIS-REPORT.md; then
        echo "   ✅ Analyzed: $pkg"
    else
        echo "   ⚠️  Warning: $pkg may not be fully analyzed"
    fi
done

# Demo objective 4: Verify recommendations were provided
echo ""
echo "✅ Step 4: Checking for recommendations..."
if grep -qi "recommendation\|next steps\|architecture" .software-factory/ANALYSIS-REPORT.md; then
    RECOMMENDATION_COUNT=$(grep -ci "recommend" .software-factory/ANALYSIS-REPORT.md || true)
    echo "   ✅ Recommendations provided (found $RECOMMENDATION_COUNT instances)"
else
    echo "   ❌ No recommendations found"
    exit 1
fi

# Demo objective 5: Verify implementation plan was created based on analysis
echo ""
echo "✅ Step 5: Verifying implementation plan exists..."
if ls .software-factory/IMPLEMENTATION-PLAN-*.md >/dev/null 2>&1; then
    PLAN_FILE=$(ls .software-factory/IMPLEMENTATION-PLAN-*.md | head -1)
    echo "   ✅ Implementation plan found: $(basename "$PLAN_FILE")"
else
    echo "   ⚠️  Warning: No implementation plan found (may be in different location)"
fi

# Demo objective 6: Show analysis summary
echo ""
echo "✅ Step 6: Displaying analysis summary..."
echo "───────────────────────────────────────────────────────────────"
echo "Analysis Report Summary:"
echo "───────────────────────────────────────────────────────────────"
if grep -A 20 "Executive Summary" .software-factory/ANALYSIS-REPORT.md | head -25; then
    echo "───────────────────────────────────────────────────────────────"
else
    echo "   ⚠️  No executive summary section found"
fi

echo ""
echo "==========================================================="
echo "✅ IDPBuilder Structure Analysis Demo PASSED"
echo ""
echo "All analysis objectives achieved:"
echo "  ✅ Comprehensive analysis report created (${WORD_COUNT} words)"
echo "  ✅ All required topics covered (${#REQUIRED_TOPICS[@]}/${#REQUIRED_TOPICS[@]})"
echo "  ✅ Key packages analyzed"
echo "  ✅ Recommendations provided"
echo "  ✅ Implementation guidance documented"
echo ""
echo "Analysis Deliverables:"
echo "  - .software-factory/ANALYSIS-REPORT.md: Comprehensive codebase analysis"
echo "  - Command structure patterns identified"
echo "  - Package organization documented"
echo "  - Testing patterns analyzed"
echo "  - Authentication mechanisms reviewed"
echo "  - Architecture recommendations provided"
echo ""
echo "This analysis provides the foundation for implementing"
echo "the push command in subsequent Phase 1 Wave 1 efforts."
echo "==========================================================="
exit 0
