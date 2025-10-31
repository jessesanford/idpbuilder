#!/bin/bash
# Apply SF 3.0 Critical Fixes - Automated State Name Replacements

# Fix #3: R516 State Naming - Replace all old names with R516-compliant names

FILES=(
    "docs/SOFTWARE-FACTORY-3.0-MIGRATION-PLAN.md"
    "docs/SOFTWARE-FACTORY-3.0-IMPLEMENTATION-CHECKLIST.md"
)

echo "Applying R516 state naming fixes..."

for file in "${FILES[@]}"; do
    echo "Processing $file..."

    # Wave level states
    sed -i 's/WAVE_1_1_SETUP/SETUP_WAVE_INFRASTRUCTURE/g' "$file"
    sed -i 's/WAVE_1_1_ITERATION_START/START_WAVE_ITERATION/g' "$file"
    sed -i 's/WAVE_1_1_INTEGRATE/INTEGRATE_WAVE_EFFORTS/g' "$file"
    sed -i 's/WAVE_1_1_REVIEW/REVIEW_WAVE_INTEGRATION/g' "$file"
    sed -i 's/WAVE_1_1_CREATE_FIX_PLAN/CREATE_WAVE_FIX_PLAN/g' "$file"
    sed -i 's/WAVE_1_1_FIX_UPSTREAM/FIX_WAVE_UPSTREAM_BUGS/g' "$file"
    sed -i 's/WAVE_1_1_ARCHITECT_REVIEW/REVIEW_WAVE_ARCHITECTURE/g' "$file"
    sed -i 's/WAVE_1_1_COMPLETE/COMPLETE_WAVE/g' "$file"

    # Phase level states
    sed -i 's/PHASE_1_SETUP/SETUP_PHASE_INFRASTRUCTURE/g' "$file"
    sed -i 's/PHASE_1_ITERATION_START/START_PHASE_ITERATION/g' "$file"
    sed -i 's/PHASE_1_INTEGRATE/INTEGRATE_PHASE_WAVES/g' "$file"
    sed -i 's/PHASE_1_REVIEW/REVIEW_PHASE_INTEGRATION/g' "$file"
    sed -i 's/PHASE_1_CREATE_FIX_PLAN/CREATE_PHASE_FIX_PLAN/g' "$file"
    sed -i 's/PHASE_1_FIX_UPSTREAM/FIX_PHASE_UPSTREAM_BUGS/g' "$file"
    sed -i 's/PHASE_1_ARCHITECT_REVIEW/REVIEW_PHASE_ARCHITECTURE/g' "$file"
    sed -i 's/PHASE_1_COMPLETE/COMPLETE_PHASE/g' "$file"

    # Project level states
    sed -i 's/PROJECT_SETUP/SETUP_PROJECT_INFRASTRUCTURE/g' "$file"
    sed -i 's/PROJECT_ITERATION_START/START_PROJECT_ITERATION/g' "$file"
    sed -i 's/PROJECT_INTEGRATE/INTEGRATE_PROJECT_PHASES/g' "$file"
    sed -i 's/PROJECT_REVIEW/REVIEW_PROJECT_INTEGRATION/g' "$file"
    sed -i 's/PROJECT_CREATE_FIX_PLAN/CREATE_PROJECT_FIX_PLAN/g' "$file"
    sed -i 's/PROJECT_FIX_UPSTREAM/FIX_PROJECT_UPSTREAM_BUGS/g' "$file"
    sed -i 's/PROJECT_ARCHITECT_REVIEW/REVIEW_PROJECT_ARCHITECTURE/g' "$file"
    sed -i 's/PROJECT_COMPLETE/COMPLETE_PROJECT/g' "$file"

    echo "✅ $file updated"
done

echo ""
echo "R516 naming fixes applied successfully!"
echo "Verifying changes..."

for file in "${FILES[@]}"; do
    OLD_COUNT=$(grep -c "WAVE_1_1\|PHASE_1_\|PROJECT_" "$file" || true)
    echo "  $file: $OLD_COUNT old-style names remaining"
done
