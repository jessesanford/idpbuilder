#!/bin/bash
# R532: Validate Agent Templates Against R383 Metadata Path Requirements
# This script ensures NO templates have hardcoded paths that violate R383

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

VIOLATIONS=0
WARNINGS=0

# Color codes for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "🔍 R532 TEMPLATE VALIDATION - Starting..."
echo "Project Root: $PROJECT_ROOT"
echo ""

# Metadata file patterns that MUST be in .software-factory/
METADATA_PATTERNS=(
    "INTEGRATE_WAVE_EFFORTS-REPORT"
    "PHASE-MERGE-PLAN"
    "WAVE-MERGE-PLAN"
    "PROJECT-MERGE-PLAN"
    "BUILD-VALIDATION-REPORT"
    "SPLIT-INVENTORY"
    "SPLIT-PLAN"
    "CODE-REVIEW-REPORT"
    "IMPLEMENTATION-PLAN"
    "FIX-PLAN"
    "TEST-REPORT"
    "DEMO-PLAN"
    "DEMO-STATUS"
    "MASTER-PR-PLAN"
    "PR-BODY"
    "PR-VALIDATION-REPORT"
)

# Files that are ALLOWED to be created in root (not metadata)
ALLOWED_ROOT_FILES=(
    "demo-wave.sh"
    "demo-features.sh"
    "verify-prs.sh"
    "integration-demo.sh"
    "\.gitignore"
    "README\.md"
    "package\.json"
    "requirements\.txt"
    "go\.mod"
)

# Check if a filename matches allowed root files
is_allowed_root_file() {
    local filename="$1"
    for allowed in "${ALLOWED_ROOT_FILES[@]}"; do
        if [[ "$filename" =~ $allowed ]]; then
            return 0  # Allowed
        fi
    done
    return 1  # Not allowed
}

# Check if a line creates a metadata file
is_metadata_file() {
    local line="$1"
    for pattern in "${METADATA_PATTERNS[@]}"; do
        if [[ "$line" =~ $pattern ]]; then
            return 0  # Is metadata
        fi
    done
    return 1  # Not metadata
}

# Validate agent configs
echo "📋 Checking agent configs (.claude/agents/*.md)..."
for file in "$PROJECT_ROOT/.claude/agents/"*.md; do
    if [ ! -f "$file" ]; then continue; fi

    filename=$(basename "$file")
    echo "  Checking: $filename"

    # Find cat > or echo > commands that create files
    while IFS= read -r line_info; do
        line_num=$(echo "$line_info" | cut -d: -f1)
        line_content=$(echo "$line_info" | cut -d: -f2-)

        # Extract the target filename from the line
        if [[ "$line_content" =~ cat.*\>\ *(['\"]?)([^\ '\"]+) ]]; then
            target_file="${BASH_REMATCH[2]}"
        elif [[ "$line_content" =~ echo.*\>\ *(['\"]?)([^\ '\"]+) ]]; then
            target_file="${BASH_REMATCH[2]}"
        else
            continue
        fi

        # Skip if it's a variable or path expression
        if [[ "$target_file" =~ ^\$ ]] || [[ "$target_file" =~ ^/ ]] || [[ "$target_file" == "/tmp/"* ]]; then
            continue
        fi

        # Check if it's a metadata file
        if is_metadata_file "$target_file"; then
            # Check if it uses .software-factory/ or sf_metadata_path
            if [[ "$line_content" =~ \.software-factory/ ]] || [[ "$line_content" =~ sf_metadata_path ]]; then
                echo -e "    ${GREEN}✓${NC} Line $line_num: $target_file (uses .software-factory/)"
            else
                echo -e "    ${RED}✗${NC} Line $line_num: $target_file ${RED}← VIOLATION!${NC}"
                echo "       Content: $(echo "$line_content" | tr -s ' ')"
                VIOLATIONS=$((VIOLATIONS + 1))
            fi
        fi
    done < <(grep -n "cat\|echo.*>" "$file" 2>/dev/null || true)
done

# Validate agent state rules
echo ""
echo "📋 Checking agent state rules (agent-states/**/rules.md)..."
while IFS= read -r file; do
    relative_path="${file#$PROJECT_ROOT/}"
    echo "  Checking: $relative_path"

    while IFS= read -r line_info; do
        line_num=$(echo "$line_info" | cut -d: -f1)
        line_content=$(echo "$line_info" | cut -d: -f2-)

        # Extract target filename
        if [[ "$line_content" =~ cat.*\>\ *(['\"]?)([^\ '\"<]+) ]]; then
            target_file="${BASH_REMATCH[2]}"
        elif [[ "$line_content" =~ echo.*\>\ *(['\"]?)([^\ '\"<]+) ]]; then
            target_file="${BASH_REMATCH[2]}"
        else
            continue
        fi

        # Skip variables, absolute paths, /tmp
        if [[ "$target_file" =~ ^\$ ]] || [[ "$target_file" =~ ^/ ]] || [[ "$target_file" == "/tmp/"* ]]; then
            continue
        fi

        # Check if it's metadata
        if is_metadata_file "$target_file"; then
            if [[ "$line_content" =~ \.software-factory/ ]] || [[ "$line_content" =~ sf_metadata_path ]]; then
                echo -e "    ${GREEN}✓${NC} Line $line_num: $target_file"
            else
                echo -e "    ${RED}✗${NC} Line $line_num: $target_file ${RED}← VIOLATION!${NC}"
                echo "       Content: $(echo "$line_content" | tr -s ' ' | cut -c1-80)"
                VIOLATIONS=$((VIOLATIONS + 1))
            fi
        fi
    done < <(grep -n "cat\|echo.*>" "$file" 2>/dev/null || true)
done < <(find "$PROJECT_ROOT/agent-states" -name "rules.md" -type f)

# Summary
echo ""
echo "═══════════════════════════════════════════════"
if [ $VIOLATIONS -eq 0 ]; then
    echo -e "${GREEN}✅ PASSED: All templates comply with R383/R532${NC}"
    echo "   No hardcoded metadata paths found"
    echo "═══════════════════════════════════════════════"
    exit 0
else
    echo -e "${RED}❌ FAILED: $VIOLATIONS R383/R532 violations found${NC}"
    echo ""
    echo "FIX REQUIRED:"
    echo "1. Replace 'cat > METADATA-FILE.md' with sf_metadata_path()"
    echo "2. Ensure all metadata goes to .software-factory/phase*/wave*/[agent]/"
    echo "3. Add timestamps using --YYYYMMDD-HHMMSS format"
    echo ""
    echo "See: rule-library/R383-metadata-file-timestamp-requirements.md"
    echo "See: rule-library/R532-template-metadata-path-validation.md"
    echo "═══════════════════════════════════════════════"
    exit 1
fi
