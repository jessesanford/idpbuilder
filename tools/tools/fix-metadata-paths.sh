#!/bin/bash

# Fix Metadata Path Violations Script
# Enforces R383 (timestamps) and R343 (directory structure)

echo "🔧 Starting comprehensive metadata path fixes..."

# Helper function to add to files
add_helper_function() {
    local file="$1"

    # Check if helper already exists
    if grep -q "sf_metadata_path()" "$file"; then
        echo "  ✓ Helper function already exists in $file"
        return
    fi

    # Find a good spot to add it (after State Context or after initial rules)
    if grep -q "## State Context" "$file"; then
        sed -i '/## State Context/i\
## Helper Functions\
\
```bash\
# Generate metadata file path with R383/R343 compliance\
sf_metadata_path() {\
    local file_type="$1"  # IMPLEMENTATION-PLAN, CODE-REVIEW-REPORT, etc.\
    local phase="$2"\
    local wave="$3"\
    local effort="$4"\
    local timestamp="${5:-$(date +%Y%m%d-%H%M%S)}"\
\
    echo ".software-factory/phase${phase}/wave${wave}/${effort}/${file_type}--${timestamp}.md"\
}\
```\
' "$file"
        echo "  ✓ Added helper function to $file"
    fi
}

# Fix WAITING_FOR_PROJECT_FIX_PLANS
echo "Fixing WAITING_FOR_PROJECT_FIX_PLANS..."
file="/home/vscode/software-factory-template/agent-states/software-factory/orchestrator/WAITING_FOR_PROJECT_FIX_PLANS/rules.md"
if [ -f "$file" ]; then
    # Replace all PROJECT-FIX-PLAN.md references
    sed -i 's|"project-integration/PROJECT-FIX-PLAN\.md"|"project-integration/.software-factory/PROJECT-FIX-PLAN--${TIMESTAMP}.md"|g' "$file"
    sed -i 's|PROJECT-FIX-PLAN\.md|PROJECT-FIX-PLAN--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed WAITING_FOR_PROJECT_FIX_PLANS"
fi

# Fix Integration states
echo "Fixing Integration state violations..."

# Fix MERGING
file="/home/vscode/software-factory-template/agent-states/software-factory/integration/MERGING/rules.md"
if [ -f "$file" ]; then
    sed -i 's|work-log\.md|.software-factory/work-log--${TIMESTAMP}.log|g' "$file"
    sed -i 's|WAVE-MERGE-PLAN\.md|.software-factory/WAVE-MERGE-PLAN--${TIMESTAMP}.md|g' "$file"
    sed -i 's|PHASE-MERGE-PLAN\.md|.software-factory/PHASE-MERGE-PLAN--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed integration/MERGING"
fi

# Fix REPORTING
file="/home/vscode/software-factory-template/agent-states/software-factory/integration/REPORTING/rules.md"
if [ -f "$file" ]; then
    sed -i 's|work-log\.md|.software-factory/work-log--${TIMESTAMP}.log|g' "$file"
    sed -i 's|INTEGRATE_WAVE_EFFORTS-PLAN\.md|.software-factory/INTEGRATE_WAVE_EFFORTS-PLAN--${TIMESTAMP}.md|g' "$file"
    sed -i 's|INTEGRATE_WAVE_EFFORTS-REPORT\.md|.software-factory/INTEGRATE_WAVE_EFFORTS-REPORT--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed integration/REPORTING"
fi

# Fix INIT
file="/home/vscode/software-factory-template/agent-states/software-factory/integration/INIT/rules.md"
if [ -f "$file" ]; then
    sed -i 's|work-log\.md|.software-factory/work-log--${TIMESTAMP}.log|g' "$file"
    sed -i 's|WAVE-MERGE-PLAN\.md|.software-factory/WAVE-MERGE-PLAN--${TIMESTAMP}.md|g' "$file"
    sed -i 's|PHASE-MERGE-PLAN\.md|.software-factory/PHASE-MERGE-PLAN--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed integration/INIT"
fi

# Fix SW-ENGINEER SPLIT_IMPLEMENTATION
echo "Fixing SW-ENGINEER SPLIT_IMPLEMENTATION..."
file="/home/vscode/software-factory-template/agent-states/software-factory/sw-engineer/SPLIT_IMPLEMENTATION/rules.md"
if [ -f "$file" ]; then
    # Fix SPLIT-PLAN references
    sed -i 's|"${SPLIT_PLAN:-SPLIT-PLAN\.md}"|"${SPLIT_PLAN:-.software-factory/SPLIT-PLAN--*.md}"|g' "$file"
    sed -i 's|SPLIT-PLAN\.md|.software-factory/SPLIT-PLAN--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed SW-ENGINEER SPLIT_IMPLEMENTATION"
fi

# Fix ERROR_RECOVERY
echo "Fixing ERROR_RECOVERY..."
file="/home/vscode/software-factory-template/agent-states/software-factory/orchestrator/ERROR_RECOVERY/rules.md"
if [ -f "$file" ]; then
    sed -i 's|PHASE-\([0-9]\+\)-WAVE-\([0-9]\+\)-REVIEW-REPORT\.md|.software-factory/PHASE-\1-WAVE-\2-REVIEW-REPORT--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed ERROR_RECOVERY"
fi

# Fix WAVE_COMPLETE
echo "Fixing WAVE_COMPLETE..."
file="/home/vscode/software-factory-template/agent-states/software-factory/orchestrator/WAVE_COMPLETE/rules.md"
if [ -f "$file" ]; then
    sed -i 's|/CODE-REVIEW-REPORT\.md"|/.software-factory/CODE-REVIEW-REPORT--${TIMESTAMP}.md"|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed WAVE_COMPLETE"
fi

# Fix MONITORING_EFFORT_REVIEWS
echo "Fixing MONITORING_EFFORT_REVIEWS..."
file="/home/vscode/software-factory-template/agent-states/software-factory/orchestrator/MONITORING_EFFORT_REVIEWS/rules.md"
if [ -f "$file" ]; then
    sed -i 's|CODE-REVIEW-REPORT\.md|.software-factory/CODE-REVIEW-REPORT--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed MONITORING_EFFORT_REVIEWS"
fi

# Fix other files with violations
echo "Fixing additional violations..."

# Fix SPAWN_SW_ENGINEERS
file="/home/vscode/software-factory-template/agent-states/software-factory/orchestrator/SPAWN_SW_ENGINEERS/rules.md"
if [ -f "$file" ]; then
    sed -i 's|SPLIT-PLAN\.md|.software-factory/SPLIT-PLAN--${TIMESTAMP}.md|g' "$file"
    sed -i 's|CODE-REVIEW-REPORT\.md|.software-factory/CODE-REVIEW-REPORT--${TIMESTAMP}.md|g' "$file"
    add_helper_function "$file"
    echo "  ✓ Fixed SPAWN_SW_ENGINEERS"
fi

echo "✅ Metadata path fixes complete!"
echo ""
echo "Running validation..."

# Validate no more violations
echo "Checking for remaining violations..."
violations=$(grep -r "IMPLEMENTATION-PLAN\.md\|SPLIT-PLAN\.md\|REVIEW-REPORT\.md\|work-log\.md" /home/vscode/software-factory-template/agent-states --include="*.md" | grep -v "\.software-factory" | grep -v "\-\-" | wc -l)

if [ "$violations" -gt 0 ]; then
    echo "⚠️  WARNING: $violations violations still remaining"
    echo "Showing sample violations:"
    grep -r "IMPLEMENTATION-PLAN\.md\|SPLIT-PLAN\.md\|REVIEW-REPORT\.md\|work-log\.md" /home/vscode/software-factory-template/agent-states --include="*.md" | grep -v "\.software-factory" | grep -v "\-\-" | head -5
else
    echo "✅ All violations fixed!"
fi