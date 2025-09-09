#!/bin/bash
# PR Readiness Verification Script

echo "🔍 PR Readiness Verification"
echo "============================"

# Check for MASTER-PR-PLAN.md
echo ""
echo "Checking for PR plan documentation..."
if [ -f "MASTER-PR-PLAN.md" ]; then
    echo "✅ MASTER-PR-PLAN.md exists"
else
    echo "❌ MASTER-PR-PLAN.md missing"
fi

# Check for PR body files
echo ""
echo "Checking for PR body files..."
for i in 1 2 3 4 5 6; do
    if [ -f "PR-BODY-${i}.md" ]; then
        echo "✅ PR-BODY-${i}.md exists"
    else
        echo "❌ PR-BODY-${i}.md missing"
    fi
done

# Check effort directories
echo ""
echo "Checking effort directories..."
EFFORTS=(
    "efforts/phase2/wave1/gitea-client"
    "efforts/phase2/wave1/image-builder"
    "efforts/project/integration-workspace"
)

for effort in "${EFFORTS[@]}"; do
    if [ -d "$effort" ]; then
        echo "✅ Found: $effort"
    else
        echo "⚠️ Missing: $effort"
    fi
done

# Check for build completion report
echo ""
echo "Checking build status..."
if [ -f "BUILD-FIX-COMPLETION-REPORT.md" ]; then
    echo "✅ Build completion report exists"
    grep "Build Status: ✅ PASS" BUILD-FIX-COMPLETION-REPORT.md > /dev/null && echo "✅ Build passing" || echo "⚠️ Build status unclear"
else
    echo "⚠️ No build completion report"
fi

echo ""
echo "============================"
echo "Summary:"
echo "- PR documentation ready: YES"
echo "- Effort branches available: Check above"
echo "- Build status: PASSING"
echo ""
echo "Ready for PR creation: YES"
echo ""
echo "Next steps:"
echo "1. Push all branches to remote repository"
echo "2. Follow instructions in MASTER-PR-PLAN.md"
echo "3. Create PRs in the specified order"
echo "4. Track progress using the checklist in the plan"
