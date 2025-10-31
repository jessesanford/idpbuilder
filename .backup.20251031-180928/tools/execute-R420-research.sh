#!/bin/bash
# R420 Research Executor
# Automates R420 Cross-Effort Planning Awareness research protocol

set -e

ORCHESTRATOR_STATE="${ORCHESTRATOR_STATE:-orchestrator-state-v3.json}"
OUTPUT_FILE="${OUTPUT_FILE:-R420-research-results.md}"

echo "==========================================="
echo "R420 PRIOR WORK RESEARCH EXECUTOR"
echo "==========================================="
echo "Orchestrator state: $ORCHESTRATOR_STATE"
echo "Output file: $OUTPUT_FILE"
echo ""

if [ ! -f "$ORCHESTRATOR_STATE" ]; then
    echo "❌ ERROR: Orchestrator state file not found: $ORCHESTRATOR_STATE"
    echo "   Set ORCHESTRATOR_STATE environment variable to correct path"
    exit 1
fi

# Initialize output file
cat > "$OUTPUT_FILE" << 'HEADER'
# 🔍 PRIOR WORK ANALYSIS (R420 MANDATORY)

**Research executed by:** [Your name/agent name]
**Research timestamp:** TIMESTAMP_PLACEHOLDER

This document contains comprehensive research of ALL previous implementations, plans, and sibling splits to prevent duplicate declarations, API mismatches, and method visibility errors.

HEADER

# Replace timestamp placeholder
TIMESTAMP=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
sed -i "s/TIMESTAMP_PLACEHOLDER/$TIMESTAMP/" "$OUTPUT_FILE"

# ============================================================================
# PHASE 1: DISCOVERY
# ============================================================================
echo "Phase 1: Discovery..."
echo "" >> "$OUTPUT_FILE"
echo "## Discovery Phase Results" >> "$OUTPUT_FILE"

# Discover previous efforts
echo "  → Discovering previous efforts..."
PREVIOUS_EFFORTS=$(jq -r '.efforts_completed[]? | .effort_id' "$ORCHESTRATOR_STATE" 2>/dev/null || echo "")

if [ -z "$PREVIOUS_EFFORTS" ]; then
    echo "- **Previous Efforts Reviewed**: NONE (this is the first effort)" >> "$OUTPUT_FILE"
    echo "    No previous efforts found (first effort in project)"
else
    echo "- **Previous Efforts Reviewed**: $(echo "$PREVIOUS_EFFORTS" | tr '\n' ', ' | sed 's/,$//')" >> "$OUTPUT_FILE"
    echo "    Found $(echo "$PREVIOUS_EFFORTS" | wc -l) previous effort(s)"
fi

# Discover sibling splits (if applicable)
echo "  → Discovering sibling splits..."
CURRENT_EFFORT=$(jq -r '.efforts_in_progress[0]? | .effort_id' "$ORCHESTRATOR_STATE" 2>/dev/null || echo "")
if [ -n "$CURRENT_EFFORT" ]; then
    PREVIOUS_SPLITS=$(jq -r ".efforts_in_progress[]? | select(.effort_id == \"$CURRENT_EFFORT\") | .splits_completed[]?.split_id" "$ORCHESTRATOR_STATE" 2>/dev/null || echo "")

    if [ -z "$PREVIOUS_SPLITS" ]; then
        echo "- **Sibling Splits Reviewed**: NONE (no previous splits for this effort)" >> "$OUTPUT_FILE"
    else
        echo "- **Sibling Splits Reviewed**: $(echo "$PREVIOUS_SPLITS" | tr '\n' ', ' | sed 's/,$//')" >> "$OUTPUT_FILE"
        echo "    Found $(echo "$PREVIOUS_SPLITS" | wc -l) previous split(s)"
    fi
fi

# Discover previous plans
echo "  → Discovering previous plans..."
PLAN_SEARCH_PATH="${PLAN_SEARCH_PATH:-.}"
PREVIOUS_PLANS=$(find "$PLAN_SEARCH_PATH" -name "*IMPLEMENTATION-PLAN*.md" -o -name "*SPLIT-PLAN*.md" 2>/dev/null || echo "")

if [ -z "$PREVIOUS_PLANS" ]; then
    echo "- **Previous Plans Reviewed**: NONE" >> "$OUTPUT_FILE"
else
    echo "- **Previous Plans Reviewed**: $(echo "$PREVIOUS_PLANS" | wc -l) plan file(s)" >> "$OUTPUT_FILE"
    echo "$PREVIOUS_PLANS" | while read plan; do
        echo "  - \`$plan\`" >> "$OUTPUT_FILE"
    done
fi

echo "- **Research Status**: IN_PROGRESS" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ============================================================================
# PHASE 2: READ IMPLEMENTATIONS
# ============================================================================
echo "Phase 2: Reading implementations..."
echo "## Implementation Analysis" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

if [ -z "$PREVIOUS_EFFORTS" ]; then
    echo "No previous efforts to analyze (this is the first effort)." >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
else
    echo "$PREVIOUS_EFFORTS" | while read effort_id; do
        echo "### Effort: $effort_id" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"

        # Get effort directory from state
        effort_dir=$(jq -r ".efforts_completed[]? | select(.effort_id == \"$effort_id\") | .effort_directory" "$ORCHESTRATOR_STATE" 2>/dev/null || echo "")

        if [ -z "$effort_dir" ] || [ ! -d "$effort_dir" ]; then
            echo "⚠️ **Directory not found or not tracked in state**" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
            continue
        fi

        echo "  → Analyzing $effort_dir..."
        echo "**Directory**: \`$effort_dir\`" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"

        # Exported functions
        echo "#### Exported Functions (can be imported)" >> "$OUTPUT_FILE"
        if grep -r "^func [A-Z]" --include="*.go" "$effort_dir" 2>/dev/null | head -10 >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
        else
            echo "None found" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        # Interfaces
        echo "#### Interfaces" >> "$OUTPUT_FILE"
        if grep -r "^type.*interface" --include="*.go" "$effort_dir" 2>/dev/null >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
        else
            echo "None found" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        # Exported structs
        echo "#### Exported Structs/Types (can be imported)" >> "$OUTPUT_FILE"
        if grep -r "^type [A-Z].*struct" --include="*.go" "$effort_dir" 2>/dev/null | head -10 >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
        else
            echo "None found" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        # Package structure
        echo "#### Package Structure" >> "$OUTPUT_FILE"
        if find "$effort_dir" -type d \( -name "pkg" -o -name "cmd" -o -name "internal" \) 2>/dev/null | head -20 >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
        else
            echo "No Go packages found" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        echo "---" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    done
fi

# Analyze sibling splits if applicable
if [ -n "$PREVIOUS_SPLITS" ]; then
    echo "## Sibling Split Analysis (CRITICAL for API compatibility!)" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"

    echo "$PREVIOUS_SPLITS" | while read split_id; do
        echo "### Split: $split_id" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"

        split_dir=$(jq -r ".efforts_in_progress[0]?.splits_completed[]? | select(.split_id == \"$split_id\") | .split_directory" "$ORCHESTRATOR_STATE" 2>/dev/null || echo "")

        if [ -z "$split_dir" ] || [ ! -d "$split_dir" ]; then
            echo "⚠️ **Directory not found**" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
            continue
        fi

        echo "  → Analyzing split: $split_dir..."
        echo "**Directory**: \`$split_dir\`" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"

        # Package exports (CRITICAL!)
        echo "#### Package Exports (APIs available for next splits)" >> "$OUTPUT_FILE"
        if grep -r "^func [A-Z]" --include="*.go" "$split_dir/pkg" 2>/dev/null >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
        else
            echo "None found" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        # Type definitions
        echo "#### Type Definitions" >> "$OUTPUT_FILE"
        if grep -r "^type [A-Z]" --include="*.go" "$split_dir/pkg" 2>/dev/null >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
        else
            echo "None found" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        # Exported methods (CRITICAL for method visibility!)
        echo "#### Exported Methods (uppercase = can access from other packages)" >> "$OUTPUT_FILE"
        if grep -r "func ([a-z]* \*[A-Z][a-zA-Z]*) [A-Z]" --include="*.go" "$split_dir" 2>/dev/null >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
        else
            echo "None found" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        # Unexported methods (CRITICAL - cannot access!)
        echo "#### Unexported Methods (lowercase = CANNOT access from other packages)" >> "$OUTPUT_FILE"
        if grep -r "func ([a-z]* \*[A-Z][a-zA-Z]*) [a-z]" --include="*.go" "$split_dir" 2>/dev/null | head -5 >> "$OUTPUT_FILE"; then
            echo "" >> "$OUTPUT_FILE"
            echo "⚠️ **These methods CANNOT be accessed from test packages or other splits!**" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        else
            echo "None found (all methods are exported)" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"
        fi

        echo "---" >> "$OUTPUT_FILE"
        echo "" >> "$OUTPUT_FILE"
    done
fi

# ============================================================================
# PHASE 3: ANALYSIS TEMPLATES
# ============================================================================
echo "Phase 3: Creating analysis templates..."

echo "## 📋 ANALYSIS TEMPLATES (Fill these out based on above findings)" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "### File Structure Findings" >> "$OUTPUT_FILE"
echo "| File Path | Source Effort | Status | Action Required |" >> "$OUTPUT_FILE"
echo "|-----------|---------------|--------|-----------------|" >> "$OUTPUT_FILE"
echo "| [example: pkg/cmd/push/root.go] | [effort ID] | EXISTS/NEW | MUST (NOT) create |" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "**Instructions:** List files that already exist and mark whether new files conflict." >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "### Interface/API Findings" >> "$OUTPUT_FILE"
echo "| Interface/API | Source | Signature | Action Required |" >> "$OUTPUT_FILE"
echo "|---------------|--------|-----------|-----------------|" >> "$OUTPUT_FILE"
echo "| [example: Registry] | [effort ID] | [full signature] | MUST implement exactly |" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "**Instructions:** Document interfaces found above with EXACT signatures." >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "### Method Visibility Findings (CRITICAL for tests!)" >> "$OUTPUT_FILE"
echo "| Method | Type | Visibility | Can Access? | Action Required |" >> "$OUTPUT_FILE"
echo "|--------|------|------------|-------------|-----------------|" >> "$OUTPUT_FILE"
echo "| [example: Push()] | [MockRegistry] | EXPORTED | YES | Safe to use in tests |" >> "$OUTPUT_FILE"
echo "| [example: validate()] | [MockRegistry] | UNEXPORTED | NO | CANNOT use |" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "**Instructions:** Check method capitalization - uppercase=EXPORTED, lowercase=UNEXPORTED." >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "### API Assumptions Verified" >> "$OUTPUT_FILE"
echo "- ✅ VERIFIED: [API/type/method] exists in [source]" >> "$OUTPUT_FILE"
echo "- ❌ INCORRECT: [assumed API] does NOT exist (use [alternative] instead)" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "**Instructions:** For EVERY API you plan to use, verify it actually exists above." >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "### Conflicts Detected" >> "$OUTPUT_FILE"
echo "- ✅ NO duplicate file paths detected" >> "$OUTPUT_FILE"
echo "- ✅ NO API mismatches detected" >> "$OUTPUT_FILE"
echo "- ✅ NO method visibility violations detected" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"
echo "OR (if conflicts found):" >> "$OUTPUT_FILE"
echo "- ❌ CONFLICT: [describe conflict]" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "### Required Integrations" >> "$OUTPUT_FILE"
echo "1. MUST [requirement from R420 research]" >> "$OUTPUT_FILE"
echo "2. MUST [requirement from R420 research]" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

echo "### Forbidden Actions" >> "$OUTPUT_FILE"
echo "- ❌ DO NOT [forbidden action with reason]" >> "$OUTPUT_FILE"
echo "- ❌ DO NOT [forbidden action with reason]" >> "$OUTPUT_FILE"
echo "" >> "$OUTPUT_FILE"

# ============================================================================
# COMPLETION
# ============================================================================
echo ""
echo "==========================================="
echo "✅ R420 RESEARCH COMPLETE"
echo "==========================================="
echo ""
echo "Research results saved to: $OUTPUT_FILE"
echo ""
echo "NEXT STEPS:"
echo "1. Review the Implementation Analysis section"
echo "2. Fill out the Analysis Templates section"
echo "3. Copy completed templates into your IMPLEMENTATION-PLAN.md"
echo "4. Run validation: tools/validate-R420-compliance.sh <your-plan.md>"
echo ""
echo "The analysis templates MUST be completed based on the findings above."
echo "DO NOT leave them as examples - use actual data from the research!"
echo ""
echo "See: rule-library/R420-cross-effort-planning-awareness-protocol.md"
