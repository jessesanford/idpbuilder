#!/bin/bash
# Fast validation script for R405 Continuation Flag compliance
# Checks that all files use CONTINUE-SOFTWARE-FACTORY=TRUE, not CONTINUE=TRUE

set -e

PROJECT_ROOT="${CLAUDE_PROJECT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"

echo "R405 Continuation Flag Fast Validator"
echo "====================================="

# Quick grep for violations
# Look for CONTINUE=TRUE but exclude CONTINUE-SOFTWARE-FACTORY=TRUE
VIOLATIONS=$(grep -r "CONTINUE=TRUE" \
    --include="*.md" \
    "$PROJECT_ROOT/agent-states" \
    "$PROJECT_ROOT/.claude/agents" \
    2>/dev/null | \
    grep -v "CONTINUE-SOFTWARE-FACTORY=TRUE" | \
    grep -v "CONTINUE=TRUE/FALSE" | \
    grep -v "Documentation" || true)

if [ -n "$VIOLATIONS" ]; then
    echo "❌ VIOLATIONS FOUND:"
    echo ""
    echo "$VIOLATIONS"
    echo ""
    echo "VIOLATION COUNT: $(echo "$VIOLATIONS" | wc -l)"
    echo ""
    echo "REQUIRED: Use CONTINUE-SOFTWARE-FACTORY=TRUE (not CONTINUE=TRUE)"
    echo "See: rule-library/R405-automation-continuation-flag.md"
    exit 1
else
    echo "✅ NO VIOLATIONS FOUND"
    echo ""
    echo "All files correctly use CONTINUE-SOFTWARE-FACTORY=TRUE"
    echo ""
    exit 0
fi
