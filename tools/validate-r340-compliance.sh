#!/bin/bash
# validate-r340-compliance.sh - R340 Compliance Checker
# Validates planning_files tracking matches filesystem reality
# Part of Software Factory 3.0 quality assurance

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

usage() {
    cat << EOF
Usage: $0 [options]

Validates R340 compliance - compares filesystem reality vs state file documentation.

Options:
  --fix       Attempt to auto-fix violations (runs discover-effort-metadata.sh)
  --strict    Exit with error on any violation (for CI/CD)
  --report    Generate detailed compliance report
  --quiet     Minimal output (errors only)

Examples:
  $0                 # Run compliance check
  $0 --fix           # Check and attempt auto-fix
  $0 --strict        # Fail build on violations
  $0 --report > r340-report.txt  # Generate report

Exit Codes:
  0 - Full compliance
  1 - Violations found (warnings)
  2 - Critical violations (blocking)
  3 - State file corrupted/missing

Part of R340 - Planning File Metadata Tracking enforcement.
EOF
    exit 0
}

# Parse arguments
FIX_MODE=false
STRICT_MODE=false
REPORT_MODE=false
QUIET_MODE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help) usage ;;
        --fix) FIX_MODE=true; shift ;;
        --strict) STRICT_MODE=true; shift ;;
        --report) REPORT_MODE=true; shift ;;
        --quiet) QUIET_MODE=true; shift ;;
        *) echo -e "${RED}Unknown option: $1${NC}"; usage ;;
    esac
done

cd "$PROJECT_ROOT"

log_info() {
    if [ "$QUIET_MODE" = "false" ]; then
        echo -e "${BLUE}$1${NC}"
    fi
}

log_success() {
    if [ "$QUIET_MODE" = "false" ]; then
        echo -e "${GREEN}$1${NC}"
    fi
}

log_warning() {
    echo -e "${YELLOW}$1${NC}"
}

log_error() {
    echo -e "${RED}$1${NC}"
}

# Check prerequisites
if [ ! -f "orchestrator-state-v3.json" ]; then
    log_error "❌ orchestrator-state-v3.json not found"
    exit 3
fi

if [ ! -d "efforts" ]; then
    log_info "ℹ️  No efforts/ directory - nothing to validate"
    exit 0
fi

log_info "════════════════════════════════════════════════════════════════"
log_info "  R340 Compliance Validation"
log_info "════════════════════════════════════════════════════════════════"
echo ""

# Initialize counters
TOTAL_EFFORTS=0
TRACKED_EFFORTS=0
UNTRACKED_EFFORTS=0
ORPHANED_METADATA=0
MISMATCHED_STATUS=0
VIOLATIONS=0

# Scan all phases and waves
for phase_dir in efforts/phase*; do
    [ -d "$phase_dir" ] || continue

    phase_num=$(basename "$phase_dir" | sed 's/phase//')

    for wave_dir in "$phase_dir"/wave*; do
        [ -d "$wave_dir" ] || continue

        wave_num=$(basename "$wave_dir" | sed 's/wave//')

        log_info "🔍 Checking Phase $phase_num Wave $wave_num..."

        # Count filesystem efforts
        fs_effort_count=$(find "$wave_dir" -maxdepth 1 -type d -name "effort-*" | wc -l)

        if [ "$fs_effort_count" -eq 0 ]; then
            log_info "  ℹ️  No efforts in filesystem"
            continue
        fi

        TOTAL_EFFORTS=$((TOTAL_EFFORTS + fs_effort_count))

        # Check state file tracking
        state_effort_json=$(jq -r ".planning_files.phases.phase${phase_num}.waves.wave${wave_num}.efforts // {}" orchestrator-state-v3.json)
        state_effort_count=$(echo "$state_effort_json" | jq '. | length')

        log_info "  Filesystem: $fs_effort_count efforts"
        log_info "  State file: $state_effort_count tracked"

        if [ "$state_effort_count" -eq 0 ]; then
            log_warning "  ⚠️  R340 VIOLATION: No efforts tracked (filesystem has $fs_effort_count)"
            UNTRACKED_EFFORTS=$((UNTRACKED_EFFORTS + fs_effort_count))
            VIOLATIONS=$((VIOLATIONS + 1))
        elif [ "$state_effort_count" -lt "$fs_effort_count" ]; then
            missing=$((fs_effort_count - state_effort_count))
            log_warning "  ⚠️  PARTIAL TRACKING: $missing efforts not tracked"
            UNTRACKED_EFFORTS=$((UNTRACKED_EFFORTS + missing))
            TRACKED_EFFORTS=$((TRACKED_EFFORTS + state_effort_count))
            VIOLATIONS=$((VIOLATIONS + 1))
        else
            TRACKED_EFFORTS=$((TRACKED_EFFORTS + state_effort_count))
            log_success "  ✅ All efforts tracked"
        fi

        # Check for orphaned metadata (completion markers not tracked as complete)
        for effort_dir in "$wave_dir"/effort-*; do
            [ -d "$effort_dir" ] || continue

            effort_name=$(basename "$effort_dir")
            complete_marker=$(find "$effort_dir" -name "IMPLEMENTATION-COMPLETE*.md" | head -1)

            if [ -n "$complete_marker" ]; then
                tracked_status=$(echo "$state_effort_json" | jq -r ".\"${effort_name}\".status // null")

                if [[ "$tracked_status" == "null" ]]; then
                    log_warning "  ⚠️  Orphaned completion marker: $effort_name (not tracked)"
                    ORPHANED_METADATA=$((ORPHANED_METADATA + 1))
                elif [[ "$tracked_status" != "completed" ]] && [[ "$tracked_status" != "approved" ]]; then
                    log_warning "  ⚠️  Status mismatch: $effort_name (marked complete, status=$tracked_status)"
                    MISMATCHED_STATUS=$((MISMATCHED_STATUS + 1))
                fi
            fi
        done

        echo ""
    done
done

# Generate summary
echo "════════════════════════════════════════════════════════════════"
log_info "R340 Compliance Summary"
echo "════════════════════════════════════════════════════════════════"
echo "Total efforts found: $TOTAL_EFFORTS"
echo "Tracked in state file: $TRACKED_EFFORTS"

if [ "$UNTRACKED_EFFORTS" -gt 0 ]; then
    log_error "Untracked efforts: $UNTRACKED_EFFORTS ❌"
else
    log_success "Untracked efforts: 0 ✅"
fi

if [ "$ORPHANED_METADATA" -gt 0 ]; then
    log_warning "Orphaned metadata files: $ORPHANED_METADATA ⚠️"
fi

if [ "$MISMATCHED_STATUS" -gt 0 ]; then
    log_warning "Status mismatches: $MISMATCHED_STATUS ⚠️"
fi

if [ "$VIOLATIONS" -gt 0 ]; then
    log_error "Total violations: $VIOLATIONS ❌"
else
    log_success "Total violations: 0 ✅"
fi

echo ""

# Calculate compliance score
if [ "$TOTAL_EFFORTS" -gt 0 ]; then
    compliance_pct=$((TRACKED_EFFORTS * 100 / TOTAL_EFFORTS))
    echo "Compliance score: ${compliance_pct}%"

    if [ "$compliance_pct" -eq 100 ]; then
        log_success "Grade: A+ (Perfect compliance)"
    elif [ "$compliance_pct" -ge 90 ]; then
        log_success "Grade: A (Excellent)"
    elif [ "$compliance_pct" -ge 80 ]; then
        log_warning "Grade: B (Good, but needs improvement)"
    elif [ "$compliance_pct" -ge 70 ]; then
        log_warning "Grade: C (Marginal compliance)"
    else
        log_error "Grade: F (Poor compliance)"
    fi
fi

# Penalty calculation
if [ "$VIOLATIONS" -gt 0 ]; then
    penalty_per_untracked=$((UNTRACKED_EFFORTS * 20))
    penalty_orphaned=$((ORPHANED_METADATA > 0 ? 50 : 0))
    total_penalty=$((penalty_per_untracked + penalty_orphaned))

    echo ""
    log_error "R340 Penalty Assessment:"
    echo "  Untracked efforts: -${penalty_per_untracked}% (${UNTRACKED_EFFORTS} × 20%)"
    [ "$penalty_orphaned" -gt 0 ] && echo "  Orphaned metadata: -${penalty_orphaned}%"
    echo "  TOTAL PENALTY: -${total_penalty}%"
fi

echo ""

# Auto-fix if requested
if [ "$FIX_MODE" = "true" ] && [ "$VIOLATIONS" -gt 0 ]; then
    log_info "🔧 Attempting auto-fix..."
    if [ -f "tools/discover-effort-metadata.sh" ]; then
        bash tools/discover-effort-metadata.sh --apply
        log_success "✅ Auto-fix complete - please review and commit changes"
    else
        log_error "❌ discover-effort-metadata.sh not found"
        exit 2
    fi
fi

# Report mode
if [ "$REPORT_MODE" = "true" ]; then
    cat << EOF

════════════════════════════════════════════════════════════════
R340 COMPLIANCE REPORT
════════════════════════════════════════════════════════════════

Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
Project: $(basename "$PROJECT_ROOT")

OVERALL STATUS: $([ "$VIOLATIONS" -eq 0 ] && echo "COMPLIANT ✅" || echo "NON-COMPLIANT ❌")

METRICS:
  Total Efforts: $TOTAL_EFFORTS
  Tracked: $TRACKED_EFFORTS ($compliance_pct%)
  Untracked: $UNTRACKED_EFFORTS
  Orphaned Metadata: $ORPHANED_METADATA
  Status Mismatches: $MISMATCHED_STATUS
  Violations: $VIOLATIONS

COMPLIANCE SCORE: $compliance_pct%
GRADE: $(if [ "$compliance_pct" -eq 100 ]; then echo "A+"; elif [ "$compliance_pct" -ge 90 ]; then echo "A"; elif [ "$compliance_pct" -ge 80 ]; then echo "B"; elif [ "$compliance_pct" -ge 70 ]; then echo "C"; else echo "F"; fi)

$([ "$VIOLATIONS" -gt 0 ] && echo "PENALTY: -${total_penalty}%")

RECOMMENDATIONS:
$([ "$UNTRACKED_EFFORTS" -gt 0 ] && echo "  • Run: bash tools/discover-effort-metadata.sh --apply")
$([ "$ORPHANED_METADATA" -gt 0 ] && echo "  • Review completion markers and update effort status")
$([ "$MISMATCHED_STATUS" -gt 0 ] && echo "  • Sync effort status with actual completion state")

═══════════════════════════════════════════════════════════════
EOF
fi

# Exit code logic
if [ "$VIOLATIONS" -eq 0 ]; then
    exit 0  # Full compliance
elif [ "$STRICT_MODE" = "true" ]; then
    exit 2  # Strict mode - violations are critical
else
    exit 1  # Violations found (warnings)
fi
