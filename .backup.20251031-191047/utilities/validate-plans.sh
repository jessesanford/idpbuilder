#!/bin/bash
# validate-plans.sh - R502 Mandatory Plan Validation Helper
# This script validates that all required planning documents exist
# before allowing phase, wave, or effort work to begin

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️ $1${NC}"
}

print_critical() {
    echo -e "${RED}🚨🚨🚨 $1 🚨🚨🚨${NC}"
}

# Validate project plans
validate_project_plans() {
    echo "🔍 Validating PROJECT-level plans..."

    local PLANNING_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}/planning/project"
    local ALL_EXIST=true

    # Required project documents
    local REQUIRED_DOCS=(
        "PROJECT-ARCHITECTURE-PLAN.md"
        "PROJECT-IMPLEMENTATION-PLAN.md"
    )

    for doc in "${REQUIRED_DOCS[@]}"; do
        if [ -f "$PLANNING_DIR/$doc" ]; then
            print_success "Found: $doc"
        else
            print_error "Missing: $PLANNING_DIR/$doc"
            ALL_EXIST=false
        fi
    done

    if [ "$ALL_EXIST" = false ]; then
        print_critical "PROJECT CANNOT START - MISSING PLANS"
        echo "Action Required: Create project-level plans first"
        return 1
    fi

    print_success "Project plans validated"
    return 0
}

# Validate phase plans
validate_phase_plans() {
    local PHASE="$1"

    if [ -z "$PHASE" ]; then
        print_error "Phase number required"
        return 1
    fi

    echo "🔍 Validating PHASE $PHASE plans..."

    local PLANNING_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}/planning/phase${PHASE}"
    local ALL_EXIST=true

    # Required phase documents
    local REQUIRED_DOCS=(
        "PHASE-${PHASE}-ARCHITECTURE-PLAN.md"
        "PHASE-${PHASE}-IMPLEMENTATION-PLAN.md"
    )

    for doc in "${REQUIRED_DOCS[@]}"; do
        if [ -f "$PLANNING_DIR/$doc" ]; then
            print_success "Found: $doc"
        else
            print_error "Missing: $PLANNING_DIR/$doc"
            ALL_EXIST=false
        fi
    done

    if [ "$ALL_EXIST" = false ]; then
        print_critical "PHASE $PHASE CANNOT START - MISSING PLANS"
        echo "Action Required: Architect and Code Reviewer must create phase plans first"
        return 1
    fi

    print_success "Phase $PHASE plans validated"
    return 0
}

# Validate wave plans
validate_wave_plans() {
    local PHASE="$1"
    local WAVE="$2"

    if [ -z "$PHASE" ] || [ -z "$WAVE" ]; then
        print_error "Phase and wave numbers required"
        return 1
    fi

    echo "🔍 Validating PHASE $PHASE WAVE $WAVE plans..."

    # First validate phase plans exist
    if ! validate_phase_plans "$PHASE"; then
        return 1
    fi

    local PLANNING_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}/planning/phase${PHASE}/wave${WAVE}"
    local ALL_EXIST=true

    # Required wave documents
    local REQUIRED_DOCS=(
        "WAVE-${PHASE}-${WAVE}-ARCHITECTURE-PLAN.md"
        "WAVE-${PHASE}-${WAVE}-IMPLEMENTATION-PLAN.md"
    )

    for doc in "${REQUIRED_DOCS[@]}"; do
        if [ -f "$PLANNING_DIR/$doc" ]; then
            print_success "Found: $doc"
        else
            print_error "Missing: $PLANNING_DIR/$doc"
            ALL_EXIST=false
        fi
    done

    if [ "$ALL_EXIST" = false ]; then
        print_critical "WAVE ${PHASE}-${WAVE} CANNOT START - MISSING PLANS"
        echo "Action Required: Architect and Code Reviewer must create wave plans first"
        return 1
    fi

    print_success "Wave ${PHASE}-${WAVE} plans validated"
    return 0
}

# Validate effort prerequisites
validate_effort_prerequisites() {
    local PHASE="$1"
    local WAVE="$2"
    local EFFORT="$3"

    if [ -z "$PHASE" ] || [ -z "$WAVE" ] || [ -z "$EFFORT" ]; then
        print_error "Phase, wave, and effort name required"
        return 1
    fi

    echo "🔍 Validating EFFORT $EFFORT prerequisites..."

    # Wave plans MUST exist
    if ! validate_wave_plans "$PHASE" "$WAVE"; then
        return 1
    fi

    # Check for effort-specific plan (may not exist yet)
    local EFFORT_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}/efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"

    if [ -d "$EFFORT_DIR/.software-factory" ]; then
        local EFFORT_PLAN=$(find "$EFFORT_DIR/.software-factory" -name "IMPLEMENTATION-PLAN--*.md" 2>/dev/null | head -1)
        if [ -n "$EFFORT_PLAN" ]; then
            print_success "Found effort plan: $(basename "$EFFORT_PLAN")"
        else
            print_warning "No effort-specific implementation plan found (Code Reviewer may not have created it yet)"
        fi
    else
        print_warning "Effort directory not yet created: $EFFORT_DIR"
    fi

    print_success "Effort prerequisites validated (phase/wave plans exist)"
    return 0
}

# Validate all plans for current state
validate_current_state() {
    local STATE_FILE="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}/orchestrator-state-v3.json"

    if [ ! -f "$STATE_FILE" ]; then
        print_error "orchestrator-state-v3.json not found"
        return 1
    fi

    local CURRENT_PHASE=$(jq -r '.current_phase' "$STATE_FILE")
    local CURRENT_WAVE=$(jq -r '.current_wave' "$STATE_FILE")

    echo "📊 Current State: Phase $CURRENT_PHASE, Wave $CURRENT_WAVE"
    echo ""

    # Validate based on current phase/wave
    if [ "$CURRENT_PHASE" = "null" ] || [ "$CURRENT_PHASE" = "0" ]; then
        echo "Project not started yet - validating project plans..."
        validate_project_plans
    elif [ "$CURRENT_WAVE" = "null" ] || [ "$CURRENT_WAVE" = "0" ]; then
        echo "Phase $CURRENT_PHASE not started yet - validating phase plans..."
        validate_phase_plans "$CURRENT_PHASE"
    else
        echo "Phase $CURRENT_PHASE Wave $CURRENT_WAVE - validating wave plans..."
        validate_wave_plans "$CURRENT_PHASE" "$CURRENT_WAVE"
    fi
}

# Check for missing plans across all phases
audit_all_plans() {
    echo "📊 COMPREHENSIVE PLAN AUDIT"
    echo "=========================="

    local PLANNING_DIR="${CLAUDE_PROJECT_DIR:-/home/vscode/software-factory-template}/planning"

    # Check project plans
    echo ""
    echo "PROJECT PLANS:"
    validate_project_plans || true

    # Check all phase directories
    for phase_dir in "$PLANNING_DIR"/phase*; do
        if [ -d "$phase_dir" ]; then
            local PHASE=$(basename "$phase_dir" | sed 's/phase//')
            echo ""
            echo "PHASE $PHASE PLANS:"
            validate_phase_plans "$PHASE" || true

            # Check all wave directories
            for wave_dir in "$phase_dir"/wave*; do
                if [ -d "$wave_dir" ]; then
                    local WAVE=$(basename "$wave_dir" | sed 's/wave//')
                    echo ""
                    echo "WAVE ${PHASE}-${WAVE} PLANS:"
                    validate_wave_plans "$PHASE" "$WAVE" || true
                fi
            done
        fi
    done
}

# Main execution
main() {
    echo "🔒 R502 PLAN VALIDATION TOOL"
    echo "============================"
    echo ""

    case "${1:-current}" in
        project)
            validate_project_plans
            ;;
        phase)
            if [ -z "$2" ]; then
                print_error "Usage: $0 phase PHASE_NUMBER"
                exit 1
            fi
            validate_phase_plans "$2"
            ;;
        wave)
            if [ -z "$2" ] || [ -z "$3" ]; then
                print_error "Usage: $0 wave PHASE_NUMBER WAVE_NUMBER"
                exit 1
            fi
            validate_wave_plans "$2" "$3"
            ;;
        effort)
            if [ -z "$2" ] || [ -z "$3" ] || [ -z "$4" ]; then
                print_error "Usage: $0 effort PHASE_NUMBER WAVE_NUMBER EFFORT_NAME"
                exit 1
            fi
            validate_effort_prerequisites "$2" "$3" "$4"
            ;;
        current)
            validate_current_state
            ;;
        audit)
            audit_all_plans
            ;;
        *)
            echo "Usage: $0 [project|phase|wave|effort|current|audit] [args...]"
            echo ""
            echo "Commands:"
            echo "  project                - Validate project-level plans"
            echo "  phase PHASE_NUM       - Validate phase plans"
            echo "  wave PHASE WAVE       - Validate wave plans"
            echo "  effort PHASE WAVE NAME - Validate effort prerequisites"
            echo "  current               - Validate based on current state (default)"
            echo "  audit                 - Audit all plans across project"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"