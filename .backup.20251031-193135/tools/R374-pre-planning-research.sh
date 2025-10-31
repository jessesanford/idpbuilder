#!/bin/bash
# R374-pre-planning-research.sh
# Performs mandatory pre-planning research to find existing code
# This tool enforces R374: Pre-Planning Research Protocol

set -e

EFFORT_REPO="${1:-$(pwd)}"
PHASE="${2:-1}"
WAVE="${3:-1}"
OUTPUT_FILE="${4:-pre-planning-research.md}"

echo "🔴🔴🔴 R374: Pre-Planning Research Protocol 🔴🔴🔴"
echo "Repository: $EFFORT_REPO"
echo "Phase: $PHASE, Wave: $WAVE"
echo "Output: $OUTPUT_FILE"
echo "================================================"

cd "$EFFORT_REPO"

# Initialize output file
cat > "$OUTPUT_FILE" << 'EOF'
# Pre-Planning Research Results (R374)

Generated: DATE_PLACEHOLDER
Phase: PHASE_PLACEHOLDER
Wave: WAVE_PLACEHOLDER

## 🔴 MANDATORY: This research MUST be included in effort plans

EOF

# Replace placeholders
sed -i "s/DATE_PLACEHOLDER/$(date -Iseconds)/" "$OUTPUT_FILE"
sed -i "s/PHASE_PLACEHOLDER/$PHASE/" "$OUTPUT_FILE"
sed -i "s/WAVE_PLACEHOLDER/$WAVE/" "$OUTPUT_FILE"

# Function to append to output
append_output() {
    echo "$1" >> "$OUTPUT_FILE"
}

# 1. Search current wave for interfaces
echo ""
echo "=== Searching Current Wave (Phase $PHASE, Wave $WAVE) ==="
append_output "## Current Wave Analysis (Phase $PHASE, Wave $WAVE)"
append_output ""

CURRENT_BRANCHES=$(git branch -r 2>/dev/null | grep "phase${PHASE}.*wave${WAVE}" | sed 's/origin\///' || true)

if [ ! -z "$CURRENT_BRANCHES" ]; then
    append_output "### Branches Analyzed:"
    echo "$CURRENT_BRANCHES" | while read branch; do
        echo "  - $branch"
        append_output "- $branch"
    done
    append_output ""

    append_output "### Interfaces Found:"
    append_output "| Interface | Location | Methods | Purpose |"
    append_output "|-----------|----------|---------|---------|"

    for branch in $CURRENT_BRANCHES; do
        echo "Checking $branch..."
        git checkout "$branch" 2>/dev/null || continue

        # Find interfaces
        INTERFACES=$(find . -name "*.go" -exec grep -l "type.*interface" {} \; 2>/dev/null || true)

        for file in $INTERFACES; do
            # Extract interface definitions
            grep -A 5 "^type.*interface" "$file" 2>/dev/null | while read -r line; do
                if [[ "$line" =~ ^type.*interface ]]; then
                    IFACE_NAME=$(echo "$line" | awk '{print $2}')
                    append_output "| $IFACE_NAME | $branch:$file | [methods] | [purpose] |"
                    echo "    Found interface: $IFACE_NAME in $file"
                fi
            done
        done
    done
    append_output ""
else
    append_output "No branches found for current wave"
    append_output ""
fi

# 2. Search previous waves
echo ""
echo "=== Searching Previous Waves ==="
append_output "## Previous Waves Analysis"
append_output ""

for prev_wave in $(seq 1 $((WAVE - 1))); do
    echo "Checking Wave $prev_wave..."
    append_output "### Wave $prev_wave:"

    PREV_BRANCHES=$(git branch -r 2>/dev/null | grep "phase${PHASE}.*wave${prev_wave}" | sed 's/origin\///' || true)

    if [ ! -z "$PREV_BRANCHES" ]; then
        append_output "| Component | Location | Type | Reusable |"
        append_output "|-----------|----------|------|----------|"

        for branch in $PREV_BRANCHES; do
            git checkout "$branch" 2>/dev/null || continue

            # Find key components
            COMPONENTS=$(find . -name "*.go" -exec grep -l "^type.*struct\|^type.*interface" {} \; 2>/dev/null | head -10 || true)

            for file in $COMPONENTS; do
                BASENAME=$(basename "$file" .go)
                append_output "| $BASENAME | $branch:$file | [type] | Yes |"
            done
        done
    else
        append_output "No branches found for wave $prev_wave"
    fi
    append_output ""
done

# 3. Search for key method signatures
echo ""
echo "=== Searching for Key Method Signatures ==="
append_output "## Key Method Signatures Found"
append_output ""
append_output "### Critical Methods to Check:"
append_output "| Method | Signature | Location | Notes |"
append_output "|--------|-----------|----------|-------|"

for method in Push Pull Upload Download Store Retrieve Create Delete Get List Update Build Process; do
    echo "Searching for ${method}..."

    FOUND=$(grep -r "func.*${method}(" --include="*.go" 2>/dev/null | head -3 || true)

    if [ ! -z "$FOUND" ]; then
        echo "$FOUND" | while IFS=: read -r file signature; do
            # Clean up signature
            SIG=$(echo "$signature" | sed 's/func.*//' | sed 's/^\s*//')
            append_output "| ${method} | ${signature} | ${file} | Check |"
            echo "    Found: ${method} in ${file}"
        done
    fi
done

append_output ""

# 4. Generate required sections for effort plans
echo ""
echo "=== Generating Template Sections ==="
append_output "## Template Sections for Effort Plan (COPY TO PLAN)"
append_output ""
append_output '```markdown'
append_output "## Pre-Planning Research Results (R374 MANDATORY)"
append_output ""
append_output "### Existing Interfaces Found"
append_output "| Interface | Location | Signature | Must Implement |"
append_output "|-----------|----------|-----------|----------------|"
append_output "| [Copy from research above] | | | |"
append_output ""
append_output "### Existing Implementations to Reuse"
append_output "| Component | Location | Purpose | How to Use |"
append_output "|-----------|----------|---------|------------|"
append_output "| [Copy from research above] | | | |"
append_output ""
append_output "### APIs Already Defined"
append_output "| API | Method | Signature | Notes |"
append_output "|-----|--------|-----------|-------|"
append_output "| [Copy from research above] | | | |"
append_output ""
append_output "### FORBIDDEN DUPLICATIONS (R373)"
append_output "- DO NOT create [interfaces found above]"
append_output "- DO NOT reimplement [existing functionality]"
append_output "- DO NOT create alternative [method] signatures"
append_output ""
append_output "### REQUIRED INTEGRATE_WAVE_EFFORTSS (R373)"
append_output "- MUST implement [interface] from [location] with EXACT signature"
append_output "- MUST reuse [component] from [location]"
append_output "- MUST import and use [package] for [functionality]"
append_output '```'

# 5. Generate validation checklist
append_output ""
append_output "## Validation Checklist"
append_output ""
append_output "Before creating effort plan, verify:"
append_output "- [ ] Searched current wave for existing interfaces"
append_output "- [ ] Searched previous waves for reusable code"
append_output "- [ ] Identified all existing method signatures"
append_output "- [ ] Documented interfaces that must be implemented"
append_output "- [ ] Listed code that must be reused"
append_output "- [ ] Specified forbidden duplications"
append_output "- [ ] Defined required integrations"

# 6. Summary
echo ""
echo "================================================"
echo "R374 Pre-Planning Research Complete"
echo "================================================"
echo "✅ Research saved to: $OUTPUT_FILE"
echo ""
echo "NEXT STEPS:"
echo "1. Review the research results in $OUTPUT_FILE"
echo "2. Copy the template sections into your effort plan"
echo "3. Fill in the specific details from the research"
echo "4. Ensure NO duplicates are created (R373 compliance)"
echo ""
echo "REMEMBER:"
echo "- Creating duplicate interfaces = -100% FAILURE"
echo "- Not researching existing code = -50% PENALTY"
echo "- All effort plans MUST include research results"

exit 0