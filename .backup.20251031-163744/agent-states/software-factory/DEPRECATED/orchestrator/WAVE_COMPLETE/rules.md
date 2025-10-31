# WAVE_COMPLETE State Rules

## State Purpose
Finalize wave completion, validate all efforts integrated, prepare wave summary, and determine next phase/wave transition.

## Entry Conditions
- From INTEGRATE_WAVE_EFFORTS when wave integration branch successfully created
- From MONITORING_EFFORT_REVIEWS when all efforts reviewed and approved
- All efforts in wave marked as "completed"

## Mandatory Validations

### 1. Verify Wave Completion Status
```bash
# Get current wave information
PHASE=$(jq -r '.current_phase' orchestrator-state-v3.json)
WAVE=$(jq -r '.current_wave' orchestrator-state-v3.json)

echo "Validating Wave ${WAVE} completion (Phase ${PHASE})..."

# Check all efforts are complete
TOTAL_EFFORTS=$(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.efforts | length" orchestrator-state-v3.json)
COMPLETED_EFFORTS=$(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.efforts | map(select(.status == \"completed\")) | length" orchestrator-state-v3.json)

if [[ "$COMPLETED_EFFORTS" -ne "$TOTAL_EFFORTS" ]]; then
    echo "ERROR: Not all efforts complete (${COMPLETED_EFFORTS}/${TOTAL_EFFORTS})"

    # List incomplete efforts
    jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.efforts[] | select(.status != \"completed\") | \"  - \" + .id + \": \" + .status" orchestrator-state-v3.json

    update_state "ERROR_RECOVERY" "Wave incomplete - efforts pending"
    exit 1
fi

echo "✅ All ${TOTAL_EFFORTS} efforts completed"
```

### 2. Verify Integration Status
```bash
# Check wave integration branch exists
INTEGRATE_WAVE_EFFORTS_BRANCH="phase-${PHASE}-wave-${WAVE}-integration"

if ! git rev-parse --verify "$INTEGRATE_WAVE_EFFORTS_BRANCH" >/dev/null 2>&1; then
    echo "ERROR: Wave integration branch '${INTEGRATE_WAVE_EFFORTS_BRANCH}' not found"
    update_state "ERROR_RECOVERY" "Missing wave integration branch"
    exit 1
fi

echo "✅ Wave integration branch exists: ${INTEGRATE_WAVE_EFFORTS_BRANCH}"

# Verify all effort branches merged
for effort_id in $(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.efforts[].id" orchestrator-state-v3.json); do
    # Check if effort was integrated
    if ! git merge-base --is-ancestor "effort-${effort_id}" "$INTEGRATE_WAVE_EFFORTS_BRANCH" 2>/dev/null; then
        echo "WARNING: Effort ${effort_id} may not be fully integrated"
    fi
done
```

### 3. Generate Wave Metrics
```bash
# Collect wave statistics
WAVE_START=$(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.started_at" orchestrator-state-v3.json)
WAVE_END=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Count totals
TOTAL_FILES_CHANGED=$(git diff --stat "main...$INTEGRATE_WAVE_EFFORTS_BRANCH" | tail -1 | awk '{print $1}')
TOTAL_INSERTIONS=$(git diff --stat "main...$INTEGRATE_WAVE_EFFORTS_BRANCH" | tail -1 | awk '{print $4}')
TOTAL_DELETIONS=$(git diff --stat "main...$INTEGRATE_WAVE_EFFORTS_BRANCH" | tail -1 | awk '{print $6}')

echo "Wave Metrics:"
echo "  Duration: ${WAVE_START} to ${WAVE_END}"
echo "  Efforts: ${TOTAL_EFFORTS}"
echo "  Files Changed: ${TOTAL_FILES_CHANGED}"
echo "  Insertions: ${TOTAL_INSERTIONS}"
echo "  Deletions: ${TOTAL_DELETIONS}"
```

## State Actions

### 1. Create Wave Completion Report
```bash
# Generate detailed wave report
cat > "wave-${WAVE}-complete-$(date +%Y%m%d-%H%M%S).md" <<EOF
# Wave ${WAVE} Completion Report (Phase ${PHASE})

## Summary
- **Wave**: ${WAVE}
- **Phase**: ${PHASE}
- **Status**: COMPLETE
- **Started**: ${WAVE_START}
- **Completed**: ${WAVE_END}

## Efforts Completed (${TOTAL_EFFORTS})
$(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.efforts[] | \"- \" + .id + \" (\" + .type + \"): \" + .status" orchestrator-state-v3.json)

## Integration
- **Branch**: ${INTEGRATE_WAVE_EFFORTS_BRANCH}
- **Files Changed**: ${TOTAL_FILES_CHANGED}
- **Lines Added**: ${TOTAL_INSERTIONS}
- **Lines Removed**: ${TOTAL_DELETIONS}

## Quality Metrics
$(jq -r ".phases.phase_${PHASE}.waves.wave_${WAVE}.quality_metrics // {} | to_entries[] | \"- \" + .key + \": \" + (.value | tostring)" orchestrator-state-v3.json)

## Next Steps
$(determine_next_steps)

## Notes
Wave ${WAVE} successfully completed all planned efforts.
EOF

echo "Wave completion report created"
```

### 2. Update Wave Status
```bash
# Mark wave as complete
jq --arg phase "$PHASE" \
   --arg wave "$WAVE" \
   --arg timestamp "$WAVE_END" \
   '.phases["phase_" + $phase].waves["wave_" + $wave].status = "completed" |
    .phases["phase_" + $phase].waves["wave_" + $wave].completed_at = $timestamp |
    .efforts_completed += .phases["phase_" + $phase].waves["wave_" + $wave].efforts |
    .efforts_in_progress = []' \
   orchestrator-state-v3.json > /tmp/state.json

mv /tmp/state.json orchestrator-state-v3.json
echo "Wave status updated to completed"
```

### 3. Determine Next Transition
```bash
determine_next_steps() {
    # Check if more waves in current phase
    TOTAL_WAVES=$(jq -r ".phases.phase_${PHASE}.total_waves" orchestrator-state-v3.json)

    if [[ "$WAVE" -lt "$TOTAL_WAVES" ]]; then
        # More waves in this phase
        NEXT_WAVE=$((WAVE + 1))
        echo "- Proceed to Wave ${NEXT_WAVE}"
        echo "- Continue Phase ${PHASE} implementation"

        # Update current_wave
        jq --arg wave "$NEXT_WAVE" '.current_wave = ($wave | tonumber)' orchestrator-state-v3.json > /tmp/state.json
        mv /tmp/state.json orchestrator-state-v3.json

        NEXT_STATE="WAVE_START"
    else
        # Phase complete - all waves done
        echo "- All waves complete for Phase ${PHASE}"
        echo "- Prepare for phase integration"

        NEXT_STATE="COMPLETE_PHASE"
    fi
}

# Execute determination
determine_next_steps
echo "Next state: ${NEXT_STATE}"
```

### 4. Archive Wave Artifacts
```bash
# Archive wave-specific artifacts
ARCHIVE_DIR="archives/phase-${PHASE}/wave-${WAVE}"
mkdir -p "$ARCHIVE_DIR"

# Copy wave reports
cp wave-${WAVE}-complete-*.md "$ARCHIVE_DIR/" 2>/dev/null || true
cp "planning/phase${PHASE}/wave${WAVE}/"*.md "$ARCHIVE_DIR/" 2>/dev/null || true

# Save final wave state
jq ".phases.phase_${PHASE}.waves.wave_${WAVE}" orchestrator-state-v3.json > "$ARCHIVE_DIR/wave-state.json"

echo "Wave artifacts archived to: ${ARCHIVE_DIR}"
```

## Exit Conditions

### Success Criteria
- All efforts verified complete
- Wave integration branch exists
- Wave metrics collected
- Completion report generated
- Next transition determined

### State Transitions
- **WAVE_START**: If more waves in current phase
- **COMPLETE_PHASE**: If all waves in phase complete
- **ERROR_RECOVERY**: If validation fails

### State Update Requirements
```bash
# Update state for transition
update_state() {
    local next_state="$1"
    local notes="${2:-Wave ${WAVE} complete}"

    jq --arg state "$next_state" \
       --arg notes "$notes" \
       --arg timestamp "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.state_machine.current_state = $state |
        .last_transition = {
            from: "WAVE_COMPLETE",
            to: $state,
            timestamp: $timestamp,
            notes: $notes
        }' orchestrator-state-v3.json > /tmp/state.json

    mv /tmp/state.json orchestrator-state-v3.json
    echo "State updated to: ${next_state}"
}

# Transition to next state
update_state "${NEXT_STATE}" "Wave ${WAVE} of Phase ${PHASE} complete"
```

## Associated Rules
- **R290**: State rule reading verification (SUPREME LAW)
- **R313**: Wave completion requirements
- **R233**: Single operation per state (SUPREME LAW)
- **R258**: Wave integration requirements
- **R288**: State file update requirements

## Prohibitions
- ❌ Mark wave complete with pending efforts
- ❌ Skip integration verification
- ❌ Proceed without completion report
- ❌ Lose effort completion data
- ❌ Start new wave without updating state

## Automation Flag

```bash
# After successful wave completion:
if [ "$MORE_WAVES" = true ]; then
    echo "✅ Wave ${WAVE} complete, proceeding to next wave"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to WAVE_START for next wave
elif [ "$MORE_PHASES" = true ]; then
    echo "✅ Phase ${PHASE} complete, proceeding to next phase"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to START_PHASE_ITERATION for next phase
else
    echo "🎉 All phases complete!"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Continue to PROJECT_DONE state
fi
```

## Notes
- Wave completion is a critical checkpoint
- All efforts must be fully integrated
- Metrics help track project progress
- Determines phase completion automatically
- Archives preserve wave history
## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

**MUST output as LAST line before state completion:**
- `CONTINUE-SOFTWARE-FACTORY=TRUE` - State completed, automation can continue
- `CONTINUE-SOFTWARE-FACTORY=FALSE` - Error/block, manual intervention required

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-CONTINUATION-FLAG-MASTER-GUIDE.md**
