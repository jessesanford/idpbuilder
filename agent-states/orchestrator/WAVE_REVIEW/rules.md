# Orchestrator - WAVE_REVIEW State Rules

## 🛑🛑🛑 R322 MANDATORY STOP BEFORE STATE TRANSITIONS 🛑🛑🛑

**SUPREME LAW - VIOLATION = -100% IMMEDIATE FAILURE**

### YOU MUST STOP AFTER:
1. ✅ Completing all TODOs for this state
2. ✅ Updating orchestrator-state.json with new state
3. ✅ Committing and pushing the state file  
4. ✅ Providing work summary

### YOU MUST NOT:
- ❌ Continue to the next state automatically
- ❌ Start work for the new state
- ❌ Spawn agents for the new state
- ❌ Assume permission to continue

### STOP PROTOCOL:
```markdown
## 🛑 STATE TRANSITION CHECKPOINT: CURRENT_STATE → NEXT_STATE

### ✅ Current State Work Completed:
- [List completed work]

### 📊 Current Status:
- Current State: CURRENT_STATE
- Next State: NEXT_STATE
- TODOs Completed: X/Y
- State Files: Updated and committed ✅

### ⏸️ STOPPED - Awaiting User Continuation
Ready to transition to NEXT_STATE. Please use /continue-orchestrating.
```

**STOP MEANS STOP - Exit and wait for /continue-orchestrating**

---


## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED WAVE_REVIEW STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_WAVE_REVIEW
echo "$(date +%s) - Rules read and acknowledged for WAVE_REVIEW" > .state_rules_read_orchestrator_WAVE_REVIEW
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY WAVE_REVIEW WORK UNTIL RULES ARE READ:
- ❌ Start spawn architecture reviewer
- ❌ Start request wave assessment
- ❌ Start collect review feedback
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R151, R208, R053..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all WAVE_REVIEW rules"
   (YOU Must READ AND ACKNOWLEDGE EACH rule individually)
   ```

3. **Silent Reading**:
   ```
   ❌ WRONG: [Reads rules but doesn't acknowledge]
   "Now I've read the rules, let me start work..."
   (MUST explicitly acknowledge EACH rule)
   ```

4. **Reading From Memory**:
   ```
   ❌ WRONG: "I know R208 requires CD before spawn..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR WAVE_REVIEW:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute WAVE_REVIEW work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY WAVE_REVIEW work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute WAVE_REVIEW work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with WAVE_REVIEW work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY WAVE_REVIEW work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 🚨 WAVE_REVIEW IS A VERB - START REVIEW PROCESS IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING WAVE_REVIEW

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Spawn Architect for wave review NOW
2. Provide integration branch for review
3. Check TodoWrite for pending items and process them
4. Include all effort completion status

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in WAVE_REVIEW" [stops]
- ❌ "Successfully entered WAVE_REVIEW state" [waits]
- ❌ "Ready to start review process" [pauses]
- ❌ "I'm in WAVE_REVIEW state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering WAVE_REVIEW, Spawn Architect for wave review NOW..."
- ✅ "START REVIEW PROCESS, provide integration branch for review..."
- ✅ "WAVE_REVIEW: Include all effort completion status..."

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## State Context
You are requesting an architect review of a completed wave's integration. The architect will assess the integrated code for architectural compliance, patterns, and readiness to proceed.

## 🔴🔴🔴 CRITICAL: MANDATORY STATE FILE UPDATE (R288) 🔴🔴🔴

### IMMEDIATELY upon entering WAVE_REVIEW state:

```bash
# 1. Update state machine
update_orchestrator_state "WAVE_REVIEW" "Requesting architect review of wave integration"

# 2. Record review request in state file
jq '.current_review.type = \"WAVE_INTEGRATION\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.phase = $PHASE' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.wave = $WAVE' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.requested_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.integration_branch = \"$INTEGRATION_BRANCH\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.status = \"PENDING\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
```

## Primary Purpose

The WAVE_REVIEW state is for:
1. Spawning an architect reviewer to assess wave integration
2. Waiting for architect feedback
3. Processing architect decision (APPROVED/CHANGES_REQUIRED/REJECTED)
4. Determining next state based on review outcome

## Architect Review Scope

The architect should review:
- **Integration Quality**: Clean merges, no conflicts remaining
- **Architectural Compliance**: Patterns followed correctly
- **Multi-tenancy**: Proper workspace/cluster isolation
- **API Consistency**: Consistent API design across efforts
- **Test Coverage**: Adequate testing at integration level
- **Performance**: No obvious performance issues introduced
- **Dependencies**: Proper dependency management

## Spawning the Architect

```bash
# Prepare architect review request
REVIEW_CONTEXT="Wave $WAVE Integration Review for Phase $PHASE"
INTEGRATION_BRANCH=$(get_wave_integration_branch_name "$PHASE" "$WAVE")

# Key information for architect
REVIEW_PROMPT="Review the integrated wave $WAVE code in branch $INTEGRATION_BRANCH. 
Focus on:
1. Architectural patterns compliance
2. Multi-tenancy implementation
3. API consistency
4. Integration quality
5. Test coverage adequacy

Provide decision: APPROVED, CHANGES_REQUIRED, or REJECTED with detailed feedback."

# Spawn architect reviewer
Task: subagent_type="architect-reviewer" \
      prompt="$REVIEW_PROMPT" \
      description="Architect review of Phase $PHASE Wave $WAVE integration"
```

## 🚨🚨🚨 RULE R258 - Mandatory Wave Review Report Verification [BLOCKING]

**BEFORE processing ANY architect decision, you MUST verify the wave review report exists:**

```bash
# MANDATORY: Verify wave review report exists (R258)
verify_wave_review_report() {
    local PHASE=$1
    local WAVE=$2
    local REPORT_FILE="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"
    
    if [ ! -f "$REPORT_FILE" ]; then
        echo "❌ CRITICAL: No wave review report found!"
        echo "❌ Expected: $REPORT_FILE"
        echo "❌ R258 VIOLATION: Cannot proceed without wave review report"
        transition_to "ERROR_RECOVERY"
        exit 1
    fi
    
    # Extract decision from report
    local DECISION=$(grep "^\*\*DECISION\*\*:" "$REPORT_FILE" | cut -d: -f2 | xargs)
    
    case "$DECISION" in
        PROCEED_NEXT_WAVE|PROCEED_PHASE_ASSESSMENT|CHANGES_REQUIRED|WAVE_FAILED)
            echo "✅ Wave review report verified with decision: $DECISION"
            ;;
        *)
            echo "❌ Invalid decision in report: $DECISION"
            transition_to "ERROR_RECOVERY"
            exit 1
            ;;
    esac
    
    # Update state with report location
    jq '.wave_review.report_file = \"$REPORT_FILE\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    jq '.wave_review.decision = \"$DECISION\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    return 0
}

# MUST call this FIRST before processing any decision
verify_wave_review_report "$CURRENT_PHASE" "$CURRENT_WAVE"
```

## Processing Architect Decision (Based on Report)

### If PROCEED_NEXT_WAVE:
```bash
# Architect approved wave, proceed to next wave
jq '.current_review.status = \"APPROVED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.decision = \"PROCEED_NEXT_WAVE\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# Transition to plan next wave
transition_to "PLANNING"  # Plan next wave
```

### If PROCEED_PHASE_ASSESSMENT:
```bash
# Last wave approved, ready for phase integration then assessment
jq '.current_review.status = \"APPROVED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.decision = \"PROCEED_PHASE_ASSESSMENT\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# CRITICAL R285: Must integrate phase BEFORE assessment!
# Phase integration required to create single integrated branch for architect review
transition_to "PHASE_INTEGRATION"  # Integrate all waves before assessment
```

### If CHANGES_REQUIRED:
```bash
# Wave needs fixes before proceeding
jq '.current_review.status = \"CHANGES_REQUIRED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.decision = \"CHANGES_REQUIRED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# CRITICAL: Setup ERROR_RECOVERY context with report location (R258)
REPORT_FILE=$(jq '.wave_review.report_file' orchestrator-state.json)
echo "📋 Setting up ERROR_RECOVERY context for wave review fixes..."

# Extract specific issues from report for ERROR_RECOVERY to process
ISSUES=$(sed -n '/## Issues Identified/,/## Required Actions/p' "$REPORT_FILE")
REQUIRED_ACTIONS=$(sed -n '/## Required Actions/,/## Recommendations/p' "$REPORT_FILE")

# Store in state file for ERROR_RECOVERY
jq '.error_recovery.source = \"WAVE_REVIEW_CHANGES_REQUIRED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.error_recovery.report_to_read = \"$REPORT_FILE\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.error_recovery.phase = $CURRENT_PHASE' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.error_recovery.wave = $CURRENT_WAVE' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.wave_review.issues_identified = \"$ISSUES\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.wave_review.required_actions = \"$REQUIRED_ACTIONS\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

echo "❌ Wave review returned CHANGES_REQUIRED - transitioning to ERROR_RECOVERY"
echo "📄 ERROR_RECOVERY must read report: $REPORT_FILE"

# Transition to fix issues - ERROR_RECOVERY will read report and coordinate fixes
transition_to "ERROR_RECOVERY"
```

### If WAVE_FAILED:
```bash
# Major architectural problems, cannot proceed
jq '.current_review.status = \"FAILED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.completed_at = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.current_review.decision = \"WAVE_FAILED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

# CRITICAL: Setup ERROR_RECOVERY context for major rework (R258)
REPORT_FILE=$(jq '.wave_review.report_file' orchestrator-state.json)
echo "❌❌❌ Wave review returned WAVE_FAILED - major architectural issues detected"

# Extract all issues for comprehensive rework
ISSUES=$(sed -n '/## Issues Identified/,/## Required Actions/p' "$REPORT_FILE")
CRITICAL_ISSUES=$(grep -A2 "\*\*\[CRITICAL\]\*\*" "$REPORT_FILE")

# Store in state file for ERROR_RECOVERY major rework
jq '.error_recovery.source = \"WAVE_FAILED\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.error_recovery.severity = \"CRITICAL\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.error_recovery.report_to_read = \"$REPORT_FILE\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.error_recovery.phase = $CURRENT_PHASE' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.error_recovery.wave = $CURRENT_WAVE' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.wave_review.critical_issues = \"$CRITICAL_ISSUES\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
jq '.wave_review.all_issues = \"$ISSUES\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json

echo "📄 ERROR_RECOVERY must read full report for major rework: $REPORT_FILE"

# Major issues require significant rework
transition_to "ERROR_RECOVERY"
```

## State Transitions

From WAVE_REVIEW (based on R258 report decision):
- **PROCEED_NEXT_WAVE** → PLANNING (plan and start next wave)
- **PROCEED_PHASE_ASSESSMENT** → PHASE_INTEGRATION (integrate waves then assess - R285)
- **CHANGES_REQUIRED** → ERROR_RECOVERY or SPAWN_AGENTS (fix identified issues)
- **WAVE_FAILED** → ERROR_RECOVERY (major architectural rework needed)
- **NO REPORT** → ERROR_RECOVERY (R258 violation - cannot proceed without report)
- **TIMEOUT** → MONITOR (check status)

**CRITICAL R285**: Never transition directly to SUCCESS from WAVE_REVIEW. If this is the last wave of the phase, MUST transition to PHASE_INTEGRATION first to create integrated branch, then to SPAWN_ARCHITECT_PHASE_ASSESSMENT for phase-level validation on the integrated code.

## Review Tracking

```yaml
# State file structure for reviews
wave_reviews:
  phase1_wave1:
    requested_at: "2025-08-25T22:00:00Z"
    architect_spawned: "architect-123456"
    integration_branch: "tmc-workspace/phase1/wave1-integration"
    decision: "APPROVED"
    feedback: |
      - Good architectural patterns
      - Multi-tenancy properly implemented
      - Minor suggestion: improve error handling in webhook controller
    completed_at: "2025-08-25T22:15:00Z"
    next_action: "PLANNING"  # Plan wave 2
```

## Timeout Handling

If architect doesn't respond within reasonable time:
```bash
# Check timeout (e.g., 30 minutes)
REQUESTED_AT=$(jq '.current_review.requested_at' orchestrator-state.json)
if [ timeout_exceeded ]; then
    # Log timeout
    jq '.current_review.status = \"TIMEOUT\"' orchestrator-state.json > tmp.json && mv tmp.json orchestrator-state.json
    
    # Transition to MONITOR to check manually
    transition_to "MONITOR"
fi
```

## Success Criteria

The wave review is successful when:
- [ ] Architect provides clear APPROVED decision
- [ ] No blocking architectural issues found
- [ ] Integration quality validated
- [ ] All patterns properly followed
- [ ] Multi-tenancy correctly implemented
- [ ] Test coverage adequate

## Common Issues to Watch For

1. **Missing Patterns**: Not following established architectural patterns
2. **Isolation Violations**: Multi-tenancy boundaries not respected
3. **API Inconsistency**: Different efforts using different API styles
4. **Test Gaps**: Integration tests missing critical paths
5. **Performance Issues**: Obvious performance problems in integrated code
6. **Dependency Conflicts**: Incompatible dependencies between efforts

## Next Steps After Approval

1. Record approval in state file
2. Tag integration branch as reviewed
3. Determine if more waves exist
4. Either plan next wave or mark phase complete
5. Update metrics and reporting

## Required State File Updates

Before leaving WAVE_REVIEW state:
- [ ] Review decision recorded
- [ ] Feedback captured
- [ ] Completion timestamp set
- [ ] Next action determined
- [ ] Any issues documented
- [ ] Metrics updated

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**
