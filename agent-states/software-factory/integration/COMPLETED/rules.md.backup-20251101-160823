# Integration Agent - COMPLETED State Rules

## State Definition
The COMPLETED state indicates successful completion of all integration tasks.

## Final Validation

### 1. Self-Assessment
```bash
echo "=== INTEGRATE_WAVE_EFFORTS COMPLETE ==="
echo "Integration Branch: $INTEGRATE_WAVE_EFFORTS_BRANCH"
echo "Branches Integrated: ${#MERGED_BRANCHES[@]}"
echo "Conflicts Resolved: $CONFLICT_COUNT"
echo "Documentation Complete: YES"
```

### 2. Grading Self-Evaluation
Review against R267 criteria:
- Completeness of Integration (50%)
  - Branch merging: ___/20%
  - Conflict resolution: ___/15%
  - Branch integrity: ___/10%
  - Final validation: ___/5%
- Documentation Quality (50%)
  - Work log: ___/25%
  - Integration report: ___/25%

### 3. Final Verification
```bash
# Verify no violations
git log --grep="cherry picked" && echo "❌ VIOLATION" || echo "✅ No cherry-picks"

# Check documentation
for doc in .software-factory/phase${PHASE}/wave${WAVE}/integration/INTEGRATE_WAVE_EFFORTS-PLAN--*.md .software-factory/phase${PHASE}/wave${WAVE}/integration/work-log--*.log .software-factory/phase${PHASE}/wave${WAVE}/integration/INTEGRATE_WAVE_EFFORTS-REPORT--*.md; do
    [[ -f "$doc" ]] && echo "✅ $doc exists" || echo "❌ Missing $doc"
done
```

## SF 3.0 Completion Protocol

When integration work is complete, this state:
1. Updates `orchestrator-state-v3.json` with completion status and final metrics
2. Sets `state_machine.current_state` to signal integration completion to orchestrator
3. Updates `integration-containers.json` with final container status (CONVERGED or issues)
4. Records all metadata locations in orchestrator-state-v3.json for downstream access
5. Commits all state changes atomically per R288

The orchestrator monitors these state files to determine when to proceed with next wave/phase operations or trigger error recovery if integration revealed cross-container issues.

## Success Indicators
- Integration branch exists and is pushed
- All documentation complete and committed
- No rule violations detected
- Ready for orchestrator review

## No Further Transitions
This is a terminal state.

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

