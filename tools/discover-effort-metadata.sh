#!/bin/bash
# discover-effort-metadata.sh - R340 State File Repair Tool
# Scans filesystem for orphaned metadata files and generates state file updates
# Part of Software Factory 3.0

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

usage() {
    cat << EOF
Usage: $0 [phase<N>] [wave<M>] [options]

Discovers orphaned effort metadata and generates state file updates.

Arguments:
  phase<N>  Optional: Scan specific phase (e.g., phase1)
  wave<M>   Optional: Scan specific wave (e.g., wave2)
            If phase given, wave filters within that phase
            If no args, scans ALL phases/waves

Options:
  --dry-run     Show what would be updated without making changes
  --apply       Apply updates to orchestrator-state-v3.json
  --output FILE Write JSON updates to FILE instead of stdout
  --format      Format output as pretty JSON (default: compact)

Examples:
  $0                      # Scan all phases/waves
  $0 phase2 wave2         # Scan only Phase 2 Wave 2
  $0 --dry-run            # Show updates without applying
  $0 --apply              # Discover and apply updates
  $0 phase1 --output p1.json  # Save Phase 1 updates to file

Exit Codes:
  0 - Success (updates found and/or applied)
  1 - Error (file not found, invalid args, etc.)
  2 - No orphaned metadata found
  3 - State file update failed

Part of R340 - Planning File Metadata Tracking enforcement.
EOF
    exit 0
}

# Parse arguments
DRY_RUN=false
APPLY=false
OUTPUT_FILE=""
FORMAT="compact"
FILTER_PHASE=""
FILTER_WAVE=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --apply)
            APPLY=true
            shift
            ;;
        --output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        --format)
            FORMAT="pretty"
            shift
            ;;
        phase[0-9]|phase[0-9][0-9])
            FILTER_PHASE="$1"
            shift
            ;;
        wave[0-9]|wave[0-9][0-9])
            FILTER_WAVE="$1"
            shift
            ;;
        *)
            echo -e "${RED}Unknown argument: $1${NC}"
            usage
            ;;
    esac
done

# Navigate to project root
cd "$PROJECT_ROOT"

echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  R340 Effort Metadata Discovery & State File Repair Tool${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════════${NC}"
echo ""

# Check if state file exists
if [ ! -f "orchestrator-state-v3.json" ]; then
    echo -e "${RED}❌ orchestrator-state-v3.json not found${NC}"
    exit 1
fi

# Check if efforts directory exists
if [ ! -d "efforts" ]; then
    echo -e "${YELLOW}⚠️  No efforts/ directory found${NC}"
    exit 2
fi

echo "🔍 Scanning for orphaned effort metadata..."
echo ""

# Determine scan scope
if [ -n "$FILTER_PHASE" ] && [ -n "$FILTER_WAVE" ]; then
    PHASE_NUM="${FILTER_PHASE#phase}"
    WAVE_NUM="${FILTER_WAVE#wave}"
    echo "📌 Scope: Phase $PHASE_NUM Wave $WAVE_NUM"
    SCAN_PATTERN="efforts/${FILTER_PHASE}/${FILTER_WAVE}/effort-*"
elif [ -n "$FILTER_PHASE" ]; then
    PHASE_NUM="${FILTER_PHASE#phase}"
    echo "📌 Scope: Phase $PHASE_NUM (all waves)"
    SCAN_PATTERN="efforts/${FILTER_PHASE}/wave*/effort-*"
else
    echo "📌 Scope: ALL phases and waves"
    SCAN_PATTERN="efforts/phase*/wave*/effort-*"
fi
echo ""

# Initialize update JSON
UPDATES_JSON="{}"

# Function to extract metadata from files
extract_effort_metadata() {
    local effort_dir="$1"
    local effort_name="$(basename "$effort_dir")"
    local phase_num=$(echo "$effort_dir" | sed -n 's|.*/phase\([0-9]*\)/.*|\1|p')
    local wave_num=$(echo "$effort_dir" | sed -n 's|.*/wave\([0-9]*\)/.*|\1|p')

    # Check if already tracked
    local is_tracked=$(jq -r ".planning_files.phases.phase${phase_num}.waves.wave${wave_num}.efforts[\"${effort_name}\"] // null" orchestrator-state-v3.json)

    if [[ "$is_tracked" != "null" ]]; then
        return 0  # Already tracked, skip
    fi

    echo -e "${YELLOW}📋 Found untracked effort: ${effort_name} (Phase ${phase_num} Wave ${wave_num})${NC}"

    # Search for metadata files
    local plan_file=$(find "$effort_dir" -name "IMPLEMENTATION-PLAN--*.md" | sort -r | head -1)
    local complete_file=$(find "$effort_dir" -name "IMPLEMENTATION-COMPLETE--*.md" | sort -r | head -1)
    local review_file=$(find "$effort_dir" -name "CODE-REVIEW-REPORT--*.md" | sort -r | head -1)

    # Initialize effort metadata
    local status="unknown"
    local reviewed=false
    local approved=false
    local implementation_lines=0
    local effort_id=""
    local effort_name_full=""
    local plan_created=""
    local impl_started=""
    local impl_completed=""
    local review_completed=""
    local review_decision=""
    local branch_name=""
    local base_branch=""
    local commit_hash=""

    # Extract from IMPLEMENTATION-PLAN
    if [ -n "$plan_file" ] && [ -f "$plan_file" ]; then
        echo "  📄 Found plan: $(basename "$plan_file")"
        effort_id=$(grep -m 1 "^**Effort ID**:" "$plan_file" | sed 's/.*: *//; s/ *$//' || echo "")
        effort_name_full=$(grep -m 1 "^**Effort Name**:" "$plan_file" | sed 's/.*: *//; s/ *$//' || echo "")
        plan_created=$(echo "$plan_file" | sed -n 's/.*--\([0-9]\{8\}\)-\([0-9]\{6\}\)\.md/\1T\2Z/p' | sed 's/\(..\)\(..\)\(..\)\(..\)\(..\)\(..\)/\1:\2:\3/')
        status="planned"
    fi

    # Extract from IMPLEMENTATION-COMPLETE
    if [ -n "$complete_file" ] && [ -f "$complete_file" ]; then
        echo "  ✅ Found completion: $(basename "$complete_file")"
        impl_completed=$(grep -m 1 "^**Completed**:" "$complete_file" | sed 's/.*: *//; s/ UTC.*//' | xargs -I {} date -d "{}" -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "")
        implementation_lines=$(grep -m 1 "^Total implementation lines:" "$complete_file" | sed 's/.*: *//; s/ .*//' || grep -m 1 "implementation_lines" "$complete_file" | sed 's/.*: *//; s/[^0-9].*//' || echo "0")
        branch_name=$(grep -m 1 "^**Branch**:" "$complete_file" | sed 's/.*: *//; s/ *$//' || echo "")
        base_branch=$(grep -m 1 "^**Base Branch**:" "$complete_file" | sed 's/.*: *//; s/ *$//' || echo "")
        commit_hash=$(grep -m 1 "^[0-9a-f]\{7\}" "$complete_file" | awk '{print $1}' || echo "")
        status="completed"
    fi

    # Extract from CODE-REVIEW-REPORT
    if [ -n "$review_file" ] && [ -f "$review_file" ]; then
        echo "  📝 Found review: $(basename "$review_file")"
        reviewed=true
        review_decision=$(grep -m 1 "^- \*\*Decision\*\*:" "$review_file" | sed 's/.*: *//; s/ *.*//' || grep -m 1 "^**Decision**:" "$review_file" | sed 's/.*: *//; s/ *.*//' || echo "UNKNOWN")
        review_completed=$(echo "$review_file" | sed -n 's/.*--\([0-9]\{8\}\)-\([0-9]\{6\}\)\.md/\1T\2Z/p' | sed 's/\(..\)\(..\)\(..\)\(..\)\(..\)\(..\)/\1:\2:\3/')

        if [[ "$review_decision" =~ "APPROVED" ]]; then
            approved=true
            status="approved"
        elif [[ "$review_decision" =~ "NEEDS_FIXES" ]]; then
            status="needs_fixes"
        fi

        # Try to extract implementation lines from review if not in completion
        if [ "$implementation_lines" -eq 0 ]; then
            implementation_lines=$(grep -m 1 "implementation_lines\|Implementation Lines\|Total.*lines" "$review_file" | sed 's/.*: *//; s/[^0-9].*//' || echo "0")
        fi
    fi

    # Build JSON for this effort
    local effort_json=$(jq -n \
        --arg eid "$effort_id" \
        --arg ename "$effort_name_full" \
        --arg status "$status" \
        --argjson reviewed "$reviewed" \
        --argjson approved "$approved" \
        --arg plan "$plan_file" \
        --arg complete "$complete_file" \
        --arg review "$review_file" \
        --argjson lines "$implementation_lines" \
        --arg plan_time "$plan_created" \
        --arg impl_time "$impl_completed" \
        --arg review_time "$review_completed" \
        --arg decision "$review_decision" \
        --arg branch "$branch_name" \
        --arg base "$base_branch" \
        --arg commit "$commit_hash" \
        '{
            effort_id: $eid,
            effort_name: $ename,
            status: $status,
            reviewed: $reviewed,
            approved: $approved,
            implementation_plan: (if $plan != "" then $plan else null end),
            implementation_complete: (if $complete != "" then $complete else null end),
            code_review_report: (if $review != "" then $review else null end),
            implementation_lines: $lines,
            plan_created_at: (if $plan_time != "" then $plan_time else null end),
            implementation_completed_at: (if $impl_time != "" then $impl_time else null end),
            review_completed_at: (if $review_time != "" then $review_time else null end),
            review_decision: (if $decision != "" then $decision else null end),
            branch_name: (if $branch != "" then $branch else null end),
            base_branch: (if $base != "" then $base else null end),
            commit_hash: (if $commit != "" then $commit else null end)
        }' | jq 'with_entries(select(.value != null))')

    # Add to updates
    UPDATES_JSON=$(echo "$UPDATES_JSON" | jq \
        --arg phase "$phase_num" \
        --arg wave "$wave_num" \
        --arg effort "$effort_name" \
        --argjson data "$effort_json" \
        '.phases["phase" + $phase].waves["wave" + $wave].efforts[$effort] = $data')

    echo "  ✅ Metadata extracted for $effort_name"
    echo ""
}

# Scan for untracked efforts
UNTRACKED_COUNT=0
for effort_dir in $SCAN_PATTERN; do
    if [ -d "$effort_dir" ]; then
        extract_effort_metadata "$effort_dir"
        ((UNTRACKED_COUNT++))
    fi
done

if [ "$UNTRACKED_COUNT" -eq 0 ]; then
    echo -e "${GREEN}✅ No orphaned metadata found - all efforts properly tracked!${NC}"
    exit 0
fi

echo "════════════════════════════════════════════════════════════════"
echo -e "${YELLOW}📊 Discovery Summary${NC}"
echo "════════════════════════════════════════════════════════════════"
echo "Untracked efforts found: $UNTRACKED_COUNT"
echo ""

# Format output based on request
if [ "$FORMAT" = "pretty" ]; then
    OUTPUT=$(echo "$UPDATES_JSON" | jq '.')
else
    OUTPUT=$(echo "$UPDATES_JSON" | jq -c '.')
fi

# Handle output destination
if [ -n "$OUTPUT_FILE" ]; then
    echo "$OUTPUT" > "$OUTPUT_FILE"
    echo -e "${GREEN}✅ Updates written to: $OUTPUT_FILE${NC}"
elif [ "$APPLY" = "false" ]; then
    echo "Discovered metadata (JSON):"
    echo "$OUTPUT"
fi

# Apply updates if requested
if [ "$APPLY" = "true" ]; then
    echo ""
    echo "📝 Applying updates to orchestrator-state-v3.json..."

    # Backup current state file
    cp orchestrator-state-v3.json "orchestrator-state-v3.json.backup-$(date +%Y%m%d-%H%M%S)"
    echo "  💾 Backup created"

    # Merge updates into state file
    jq \
        --argjson updates "$UPDATES_JSON" \
        '.planning_files.phases = (.planning_files.phases // {}) |
         ($updates.phases | to_entries[] |
          .key as $phase |
          .value.waves | to_entries[] |
          .key as $wave |
          .value.efforts | to_entries[] |
          . as $effort |
          .planning_files.phases[$phase].waves[$wave].efforts =
            ((.planning_files.phases[$phase].waves[$wave].efforts // {}) + {($effort.key): $effort.value})
         ) | .' \
        orchestrator-state-v3.json > orchestrator-state-v3.json.tmp

    if [ $? -eq 0 ]; then
        mv orchestrator-state-v3.json.tmp orchestrator-state-v3.json
        echo -e "  ${GREEN}✅ State file updated${NC}"
        echo ""
        echo "⚠️  IMPORTANT: Review changes and commit:"
        echo "   git diff orchestrator-state-v3.json"
        echo "   git add orchestrator-state-v3.json"
        echo "   git commit -m \"fix: Backfill effort metadata per R340 [discovered by tool]\""
    else
        echo -e "  ${RED}❌ Failed to update state file${NC}"
        rm -f orchestrator-state-v3.json.tmp
        exit 3
    fi
fi

echo ""
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}R340 Metadata Discovery Complete${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"

exit 0
