#!/bin/bash
# utilities/create-bug.sh
# R421 - Universal Bug Creation Utility
# Supports both cascade and non-cascade bugs

set -euo pipefail

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
STATE_FILE="${STATE_FILE:-$PROJECT_ROOT/orchestrator-state-v3.json}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Usage information
usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Create a bug in the bug_registry supporting both cascade and non-cascade sources.

REQUIRED OPTIONS:
    -s, --source SOURCE         Bug source: cascade|implementation|review|user_feedback|documentation|test|other
    -d, --description DESC      Human-readable bug description
    -S, --severity SEVERITY     Severity: CRITICAL|HIGH|MEDIUM|LOW
    -e, --efforts EFFORTS       Comma-separated list of affected effort IDs

CASCADE-SPECIFIC OPTIONS (required when source=cascade):
    -c, --cascade-id ID         Cascade ID (format: cascade-YYYYMMDD-HHMMSS)
    -l, --layer LAYER           Cascade layer number (integer >= 1)
    -i, --integration NAME      Integration name where bug was detected
    -t, --integration-type TYPE Integration type: wave|phase|project
    -p, --phase PHASE           Phase number (for wave/phase integrations)
    -w, --wave WAVE             Wave number (for wave integrations)

OPTIONAL OPTIONS:
    -C, --category CATEGORY     Category: build_failure|test_failure|lint_error|runtime_error|integration_conflict|other
    -D, --detected-by DETECTOR  What detected the bug (e.g., "R291 Build Gate", "SW Engineer", "User Report")
    -E, --error-message MSG     Error message text
    -f, --file-path PATH        File path where error occurred
    -n, --line-number NUM       Line number where error occurred
    -P, --primary-effort ID     Primary effort responsible for fix
    -r, --coordinated           Fix requires coordination across efforts (flag)
    -h, --help                  Show this help message

EXAMPLES:

    # CASCADE BUG (most common - integration failure)
    $0 --source cascade \\
       --cascade-id cascade-20251005-051000 \\
       --layer 1 \\
       --integration phase1_wave2 \\
       --integration-type wave \\
       --phase 1 \\
       --wave 2 \\
       --severity CRITICAL \\
       --category build_failure \\
       --description "PushCmd redeclared in multiple packages" \\
       --efforts "E1.2.1,E1.2.2" \\
       --primary-effort E1.2.1

    # IMPLEMENTATION BUG (found during coding)
    $0 --source implementation \\
       --severity HIGH \\
       --category runtime_error \\
       --description "Null pointer in auth.Validate()" \\
       --efforts "E1.2.1" \\
       --detected-by "SW Engineer"

    # USER FEEDBACK BUG (post-deployment report)
    $0 --source user_feedback \\
       --severity MEDIUM \\
       --description "Push command fails with 500 error for large files" \\
       --efforts "E1.2.3" \\
       --detected-by "User Report #42"

    # REVIEW BUG (found during code review)
    $0 --source review \\
       --severity LOW \\
       --category other \\
       --description "Inconsistent error handling in auth package" \\
       --efforts "E1.3.1" \\
       --detected-by "Code Reviewer"

EOF
    exit 1
}

# Parse arguments
BUG_SOURCE=""
DESCRIPTION=""
SEVERITY=""
AFFECTED_EFFORTS=""
CASCADE_ID=""
CASCADE_LAYER=""
INTEGRATE_WAVE_EFFORTS_NAME=""
INTEGRATE_WAVE_EFFORTS_TYPE=""
INTEGRATE_WAVE_EFFORTS_PHASE=""
INTEGRATE_WAVE_EFFORTS_WAVE=""
CATEGORY="other"
DETECTED_BY=""
ERROR_MESSAGE=""
FILE_PATH=""
LINE_NUMBER=""
PRIMARY_EFFORT=""
COORDINATED=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--source)
            BUG_SOURCE="$2"
            shift 2
            ;;
        -d|--description)
            DESCRIPTION="$2"
            shift 2
            ;;
        -S|--severity)
            SEVERITY="$2"
            shift 2
            ;;
        -e|--efforts)
            AFFECTED_EFFORTS="$2"
            shift 2
            ;;
        -c|--cascade-id)
            CASCADE_ID="$2"
            shift 2
            ;;
        -l|--layer)
            CASCADE_LAYER="$2"
            shift 2
            ;;
        -i|--integration)
            INTEGRATE_WAVE_EFFORTS_NAME="$2"
            shift 2
            ;;
        -t|--integration-type)
            INTEGRATE_WAVE_EFFORTS_TYPE="$2"
            shift 2
            ;;
        -p|--phase)
            INTEGRATE_WAVE_EFFORTS_PHASE="$2"
            shift 2
            ;;
        -w|--wave)
            INTEGRATE_WAVE_EFFORTS_WAVE="$2"
            shift 2
            ;;
        -C|--category)
            CATEGORY="$2"
            shift 2
            ;;
        -D|--detected-by)
            DETECTED_BY="$2"
            shift 2
            ;;
        -E|--error-message)
            ERROR_MESSAGE="$2"
            shift 2
            ;;
        -f|--file-path)
            FILE_PATH="$2"
            shift 2
            ;;
        -n|--line-number)
            LINE_NUMBER="$2"
            shift 2
            ;;
        -P|--primary-effort)
            PRIMARY_EFFORT="$2"
            shift 2
            ;;
        -r|--coordinated)
            COORDINATED=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo -e "${RED}Error: Unknown option $1${NC}"
            usage
            ;;
    esac
done

# Validate required arguments
if [[ -z "$BUG_SOURCE" ]]; then
    echo -e "${RED}Error: --source is required${NC}"
    usage
fi

if [[ -z "$DESCRIPTION" ]]; then
    echo -e "${RED}Error: --description is required${NC}"
    usage
fi

if [[ -z "$SEVERITY" ]]; then
    echo -e "${RED}Error: --severity is required${NC}"
    usage
fi

if [[ -z "$AFFECTED_EFFORTS" ]]; then
    echo -e "${RED}Error: --efforts is required${NC}"
    usage
fi

# Validate bug source
case "$BUG_SOURCE" in
    cascade|implementation|review|user_feedback|documentation|test|other)
        ;;
    *)
        echo -e "${RED}Error: Invalid source '$BUG_SOURCE'. Must be: cascade|implementation|review|user_feedback|documentation|test|other${NC}"
        exit 1
        ;;
esac

# Validate severity
case "$SEVERITY" in
    CRITICAL|HIGH|MEDIUM|LOW)
        ;;
    *)
        echo -e "${RED}Error: Invalid severity '$SEVERITY'. Must be: CRITICAL|HIGH|MEDIUM|LOW${NC}"
        exit 1
        ;;
esac

# Validate category
case "$CATEGORY" in
    build_failure|test_failure|lint_error|runtime_error|integration_conflict|other)
        ;;
    *)
        echo -e "${RED}Error: Invalid category '$CATEGORY'. Must be: build_failure|test_failure|lint_error|runtime_error|integration_conflict|other${NC}"
        exit 1
        ;;
esac

# CASCADE-SPECIFIC VALIDATION
if [[ "$BUG_SOURCE" == "cascade" ]]; then
    if [[ -z "$CASCADE_ID" ]]; then
        echo -e "${RED}Error: --cascade-id is required for cascade bugs${NC}"
        exit 1
    fi
    if [[ -z "$CASCADE_LAYER" ]]; then
        echo -e "${RED}Error: --layer is required for cascade bugs${NC}"
        exit 1
    fi
    if [[ -z "$INTEGRATE_WAVE_EFFORTS_NAME" ]]; then
        echo -e "${RED}Error: --integration is required for cascade bugs${NC}"
        exit 1
    fi
    if [[ -z "$INTEGRATE_WAVE_EFFORTS_TYPE" ]]; then
        echo -e "${RED}Error: --integration-type is required for cascade bugs${NC}"
        exit 1
    fi

    # Validate cascade_id format
    if ! [[ "$CASCADE_ID" =~ ^cascade-[0-9]{8}-[0-9]{6}$ ]]; then
        echo -e "${RED}Error: Invalid cascade_id format. Expected: cascade-YYYYMMDD-HHMMSS${NC}"
        exit 1
    fi

    # Validate integration type
    case "$INTEGRATE_WAVE_EFFORTS_TYPE" in
        wave|phase|project)
            ;;
        *)
            echo -e "${RED}Error: Invalid integration type '$INTEGRATE_WAVE_EFFORTS_TYPE'. Must be: wave|phase|project${NC}"
            exit 1
            ;;
    esac
fi

# Check state file exists
if [[ ! -f "$STATE_FILE" ]]; then
    echo -e "${RED}Error: State file not found: $STATE_FILE${NC}"
    exit 1
fi

# Generate bug ID
if [[ "$BUG_SOURCE" == "cascade" ]]; then
    # Cascade bug: BUG-cascade-YYYYMMDD-HHMMSS-NNN
    # Extract timestamp from cascade_id
    CASCADE_TIMESTAMP="${CASCADE_ID#cascade-}"

    # Count existing cascade bugs with this cascade_id to get next number
    EXISTING_COUNT=$(jq -r --arg cid "$CASCADE_ID" '[.bug_registry[] | select(.cascade_id == $cid)] | length' "$STATE_FILE")
    BUG_NUMBER=$(printf "%03d" $((EXISTING_COUNT + 1)))
    BUG_ID="BUG-cascade-${CASCADE_TIMESTAMP}-${BUG_NUMBER}"
else
    # Non-cascade bug: BUG-YYYYMMDD-NNN
    TODAY=$(date +%Y%m%d)

    # Count existing non-cascade bugs from today to get next number
    EXISTING_COUNT=$(jq -r --arg today "$TODAY" '[.bug_registry[] | select(.bug_source != "cascade" and (.bug_id | startswith("BUG-" + $today)))] | length' "$STATE_FILE")
    BUG_NUMBER=$(printf "%03d" $((EXISTING_COUNT + 1)))
    BUG_ID="BUG-${TODAY}-${BUG_NUMBER}"
fi

# Convert comma-separated efforts to JSON array
EFFORTS_ARRAY=$(echo "$AFFECTED_EFFORTS" | jq -R 'split(",") | map(ltrimstr(" ") | rtrimstr(" "))')

# Set primary effort if not specified
if [[ -z "$PRIMARY_EFFORT" ]]; then
    PRIMARY_EFFORT=$(echo "$EFFORTS_ARRAY" | jq -r '.[0]')
fi

# Set detected_by if not specified
if [[ -z "$DETECTED_BY" ]]; then
    case "$BUG_SOURCE" in
        cascade)
            DETECTED_BY="R291 Build Gate"
            ;;
        implementation)
            DETECTED_BY="SW Engineer"
            ;;
        review)
            DETECTED_BY="Code Reviewer"
            ;;
        user_feedback)
            DETECTED_BY="User Report"
            ;;
        documentation)
            DETECTED_BY="Documentation Review"
            ;;
        test)
            DETECTED_BY="Test Suite"
            ;;
        *)
            DETECTED_BY="Unknown"
            ;;
    esac
fi

# Build the bug object
echo -e "${BLUE}Creating bug: $BUG_ID${NC}"

# Start building jq command
JQ_FILTER='.bug_registry += [{'
JQ_FILTER+='  "bug_id": $bug_id,'
JQ_FILTER+='  "bug_source": $bug_source,'
JQ_FILTER+='  "detected_at": (now | todate),'
JQ_FILTER+='  "detected_by": $detected_by,'
JQ_FILTER+='  "severity": $severity,'
JQ_FILTER+='  "category": $category,'
JQ_FILTER+='  "description": $description,'
JQ_FILTER+='  "affected_efforts": $efforts,'
JQ_FILTER+='  "primary_effort": $primary_effort,'
JQ_FILTER+='  "requires_coordinated_fix": $coordinated,'
JQ_FILTER+='  "fix_status": "pending",'
JQ_FILTER+='  "fix_attempts": [],'
JQ_FILTER+='  "current_attempt": 0,'
JQ_FILTER+='  "integration_status": "not_integrated",'
JQ_FILTER+='  "integrated_at": null,'
JQ_FILTER+='  "integration_commit": null,'
JQ_FILTER+='  "blocked_by": [],'
JQ_FILTER+='  "blocks": [],'
JQ_FILTER+='  "resolution_notes": "",'

# CASCADE-SPECIFIC FIELDS
if [[ "$BUG_SOURCE" == "cascade" ]]; then
    JQ_FILTER+='  "cascade_id": $cascade_id,'
    JQ_FILTER+='  "cascade_layer": ($cascade_layer | tonumber),'
    JQ_FILTER+='  "detected_in_integration": {'
    JQ_FILTER+='    "name": $integration_name,'
    JQ_FILTER+='    "type": $integration_type,'

    if [[ -n "$INTEGRATE_WAVE_EFFORTS_PHASE" ]]; then
        JQ_FILTER+='    "phase": ($integration_phase | tonumber),'
    else
        JQ_FILTER+='    "phase": null,'
    fi

    if [[ -n "$INTEGRATE_WAVE_EFFORTS_WAVE" ]]; then
        JQ_FILTER+='    "wave": ($integration_wave | tonumber)'
    else
        JQ_FILTER+='    "wave": null'
    fi

    JQ_FILTER+='  },'
else
    # Non-cascade bugs have null for cascade fields
    JQ_FILTER+='  "cascade_id": null,'
    JQ_FILTER+='  "cascade_layer": null,'
    JQ_FILTER+='  "detected_in_integration": null,'
fi

# ERROR DETAILS (optional)
if [[ -n "$ERROR_MESSAGE" ]] || [[ -n "$FILE_PATH" ]]; then
    JQ_FILTER+='  "error_details": {'

    if [[ -n "$ERROR_MESSAGE" ]]; then
        JQ_FILTER+='    "error_message": $error_message,'
    fi

    if [[ -n "$FILE_PATH" ]]; then
        JQ_FILTER+='    "file_path": $file_path,'
    fi

    if [[ -n "$LINE_NUMBER" ]]; then
        JQ_FILTER+='    "line_number": ($line_number | tonumber)'
    else
        JQ_FILTER+='    "line_number": null'
    fi

    JQ_FILTER+='  },'
fi

# TIMESTAMPS
JQ_FILTER+='  "created_at": (now | todate),'
JQ_FILTER+='  "updated_at": (now | todate)'

JQ_FILTER+='}]'

# Build jq arguments
JQ_ARGS=(
    --arg bug_id "$BUG_ID"
    --arg bug_source "$BUG_SOURCE"
    --arg detected_by "$DETECTED_BY"
    --arg severity "$SEVERITY"
    --arg category "$CATEGORY"
    --arg description "$DESCRIPTION"
    --argjson efforts "$EFFORTS_ARRAY"
    --arg primary_effort "$PRIMARY_EFFORT"
    --argjson coordinated "$COORDINATED"
)

if [[ "$BUG_SOURCE" == "cascade" ]]; then
    JQ_ARGS+=(
        --arg cascade_id "$CASCADE_ID"
        --arg cascade_layer "$CASCADE_LAYER"
        --arg integration_name "$INTEGRATE_WAVE_EFFORTS_NAME"
        --arg integration_type "$INTEGRATE_WAVE_EFFORTS_TYPE"
    )

    if [[ -n "$INTEGRATE_WAVE_EFFORTS_PHASE" ]]; then
        JQ_ARGS+=(--arg integration_phase "$INTEGRATE_WAVE_EFFORTS_PHASE")
    fi

    if [[ -n "$INTEGRATE_WAVE_EFFORTS_WAVE" ]]; then
        JQ_ARGS+=(--arg integration_wave "$INTEGRATE_WAVE_EFFORTS_WAVE")
    fi
fi

if [[ -n "$ERROR_MESSAGE" ]]; then
    JQ_ARGS+=(--arg error_message "$ERROR_MESSAGE")
fi

if [[ -n "$FILE_PATH" ]]; then
    JQ_ARGS+=(--arg file_path "$FILE_PATH")
fi

if [[ -n "$LINE_NUMBER" ]]; then
    JQ_ARGS+=(--arg line_number "$LINE_NUMBER")
fi

# Execute jq to add bug
jq "${JQ_ARGS[@]}" "$JQ_FILTER" "$STATE_FILE" > "${STATE_FILE}.tmp"
mv "${STATE_FILE}.tmp" "$STATE_FILE"

# Print summary
echo -e "${GREEN}✅ Bug created successfully!${NC}"
echo ""
echo "Bug Details:"
echo "  ID: $BUG_ID"
echo "  Source: $BUG_SOURCE"
echo "  Severity: $SEVERITY"
echo "  Category: $CATEGORY"
echo "  Description: $DESCRIPTION"
echo "  Affected Efforts: $AFFECTED_EFFORTS"

if [[ "$BUG_SOURCE" == "cascade" ]]; then
    echo "  Cascade ID: $CASCADE_ID"
    echo "  Cascade Layer: $CASCADE_LAYER"
    echo "  Integration: $INTEGRATE_WAVE_EFFORTS_NAME ($INTEGRATE_WAVE_EFFORTS_TYPE)"
fi

echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Assign bug to appropriate engineer/agent"
echo "  2. Update bug status as fix progresses"
echo "  3. Track fix commits in fix_attempts"

# Return bug ID for scripting
echo "$BUG_ID"
