# WAITING_FOR_WAVE_TEST_PLAN State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_WAVE_TEST_PLAN STATE

### Core Mandatory Rules (ALL orchestrator states must have these):

1. **🚨🚨🚨 R006** - ORCHESTRATOR NEVER WRITES CODE OR PERFORMS FILE OPERATIONS (BLOCKING)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
   - Criticality: BLOCKING - Automatic termination, 0% grade
   - Summary: NEVER write, copy, move, or manipulate ANY code files - delegate ALL to agents

2. **🔴🔴🔴 R287** - TODO PERSISTENCE COMPREHENSIVE (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
   - Criticality: SUPREME - -20% to -100% penalty for violations
   - Summary: MUST save TODOs within 30s after write, every 10 messages, before transitions

3. **🔴🔴🔴 R288** - STATE FILE UPDATE REQUIREMENTS (SUPREME LAW)
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state-v3.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322-mandatory-stop-before-state-transitions.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-all-states-immediate-action.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: WAITING states require active monitoring, not passive waiting

## State Purpose
Monitor Code Reviewer creating wave-level tests and demo plans. Ensure all tests are ready BEFORE wave implementation planning begins (TDD enforcement).

## Entry Conditions
- Code Reviewer spawned for wave test planning
- Current state is SPAWN_CODE_REVIEWER_WAVE_TEST_PLANNING
- Waiting for wave test deliverables

## Required Actions

### 1. Monitor Wave Test Creation Progress
```bash
# Check for wave test deliverables
check_wave_test_deliverables() {
    local wave_dir="wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}"
    
    echo "🔍 Checking for wave test deliverables..."
    
    # Required files
    local files_needed=(
        "WAVE-TEST-PLAN.md"
        "WAVE-TEST-HARNESS.sh"
        "WAVE-DEMO-PLAN.md"
    )
    
    for file in "${files_needed[@]}"; do
        if [ -f "$wave_dir/$file" ]; then
            echo "✅ Found: $file"
        else
            echo "⏳ Waiting for: $file"
            return 1
        fi
    done
    
    # Check for actual test files
    if ls "$wave_dir/tests/wave/"*.test.* >/dev/null 2>&1; then
        echo "✅ Found wave test files"
    else
        echo "⏳ Waiting for wave test files"
        return 1
    fi
    
    echo "✅ All wave test deliverables ready!"
    return 0
}
```

### 2. Validate Wave Test Completeness
```bash
validate_wave_tests() {
    local wave_dir="wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}"
    
    echo "🧪 Validating wave test completeness..."
    
    # Test plan must reference wave architecture
    if grep -q "Wave Architecture Coverage" "$wave_dir/WAVE-TEST-PLAN.md"; then
        echo "✅ Test plan covers wave architecture"
    else
        echo "❌ Test plan missing wave architecture coverage"
        return 1
    fi
    
    # Test harness must be executable
    if [ -x "$wave_dir/WAVE-TEST-HARNESS.sh" ]; then
        echo "✅ Wave test harness is executable"
    else
        echo "❌ Wave test harness not executable"
        chmod +x "$wave_dir/WAVE-TEST-HARNESS.sh"
    fi
    
    # Demo plan must include wave scenarios (R330)
    if grep -q "Wave Demo Scenarios" "$wave_dir/WAVE-DEMO-PLAN.md"; then
        echo "✅ Wave demo scenarios defined"
    else
        echo "❌ Wave demo scenarios missing (R330 violation)"
        return 1
    fi
    
    # Test assignment mapping for efforts
    if grep -q "Effort Test Assignments" "$wave_dir/WAVE-TEST-PLAN.md"; then
        echo "✅ Test assignments for efforts defined"
    else
        echo "⚠️ Warning: No explicit effort test assignments"
    fi
    
    # Tests should fail (no implementation yet)
    echo "🔴 Running wave tests (expecting failures - no implementation yet)..."
    if cd "$wave_dir" && ./WAVE-TEST-HARNESS.sh; then
        echo "⚠️ WARNING: Wave tests passing without implementation?"
    else
        echo "✅ Wave tests failing as expected (TDD - red phase)"
    fi
}
```

### 3. Prepare for Implementation Planning
```bash
prepare_test_context_for_implementation() {
    local wave_dir="wave-tests/phase-${PHASE_NUM}/wave-${WAVE_NUM}"
    
    echo "📋 Preparing test context for implementation planning..."
    
    # Extract test assignments for efforts
    grep -A 20 "Effort Test Assignments" "$wave_dir/WAVE-TEST-PLAN.md" > test-assignments.txt
    
    # Create test reference for implementation planner
    cat > test-reference.md << EOF
# WAVE ${WAVE_NUM} TEST REFERENCE (TDD)

## Tests Created (Must Pass Before Wave Complete)
$(ls "$wave_dir/tests/wave/"*.test.* | xargs -I {} basename {})

## Demo Requirements (R330/R291)
- Wave demo plan: $wave_dir/WAVE-DEMO-PLAN.md
- Integration demo required even for single-effort waves

## Implementation Guidance
- Each effort must reference which tests it targets
- Implementation success = assigned tests pass
- No new features without corresponding tests
EOF
}
```

### 4. Capture Wave Test Metadata (R340 Compliance)
```bash
# When Code Reviewer reports wave test locations
capture_wave_test_metadata() {
    local phase_num="${PHASE_NUM}"
    local wave_num="${WAVE_NUM}"
    local wave_key="phase${phase_num}_wave${wave_num}"
    
    echo "📋 Capturing wave test metadata for state tracking..."
    
    # Read effort assignments from test-effort-mapping.json if available
    local mapping_file="wave-tests/phase-${phase_num}/wave-${wave_num}/test-effort-mapping.json"
    local effort_assignments="{}"
    if [ -f "$mapping_file" ]; then
        effort_assignments=$(jq '.test_assignments | to_entries | map({key: .key, value: .value.tests}) | from_entries' "$mapping_file")
    fi
    
    # Update test_plans section with wave test metadata
    jq --arg key "$wave_key" \
       --arg test_plan "/wave-tests/phase-${phase_num}/wave-${wave_num}/WAVE-TEST-PLAN.md" \
       --arg harness "/wave-tests/phase-${phase_num}/wave-${wave_num}/WAVE-TEST-HARNESS.sh" \
       --arg demo "/wave-tests/phase-${phase_num}/wave-${wave_num}/WAVE-DEMO-PLAN.md" \
       --arg test_dir "/wave-tests/phase-${phase_num}/wave-${wave_num}/tests/wave" \
       --arg mapping "/wave-tests/phase-${phase_num}/wave-${wave_num}/test-effort-mapping.json" \
       --arg created "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       --argjson assignments "$effort_assignments" \
       '.test_plans.wave[$key] = {
            "test_plan_path": $test_plan,
            "test_harness_path": $harness,
            "demo_plan_path": $demo,
            "test_dir": $test_dir,
            "test_effort_mapping_path": $mapping,
            "created_at": $created,
            "created_by": "code-reviewer",
            "phase": '$phase_num',
            "wave": '$wave_num',
            "status": "active",
            "tdd_phase": "red",
            "effort_assignments": $assignments
        }' orchestrator-state-v3.json > tmp.json
    
    mv tmp.json orchestrator-state-v3.json
    
    echo "✅ Wave test metadata captured in state file"
    echo "   SW Engineers can find their assigned tests from state"
}
```

### 5. Transition to Implementation Planning
```bash
# When all wave tests are ready
transition_to_wave_implementation_planning() {
    echo "📋 Wave tests ready - proceeding to implementation planning"

    # Capture test metadata if not already done
    capture_wave_test_metadata

    # Update tracking fields (ALLOWED - orchestrator maintains this data)
    jq '.wave_test_planning.completed_at = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' \
        orchestrator-state-v3.json > tmp.json
    mv tmp.json orchestrator-state-v3.json

    jq '.wave_test_planning.test_directory = "wave-tests/phase-'${PHASE_NUM}'/wave-'${WAVE_NUM}'"' \
        orchestrator-state-v3.json > tmp.json
    mv tmp.json orchestrator-state-v3.json

    # Set proposed next state (State Manager will update state_machine fields)
    PROPOSED_NEXT_STATE="SPAWN_CODE_REVIEWER_WAVE_IMPL"
    TRANSITION_REASON="Wave tests ready - captured locations in state file (R340)"
    # State Manager consultation happens in Step 3 of completion checklist
}
```

## Exit Conditions
- All wave test deliverables created and validated
- Transition to SPAWN_CODE_REVIEWER_WAVE_IMPL
- Tests documented and ready for effort planning

## Success Criteria
- ✅ WAVE-TEST-PLAN.md exists with architecture coverage
- ✅ WAVE-TEST-HARNESS.sh is executable
- ✅ Test files created in tests/wave/
- ✅ WAVE-DEMO-PLAN.md includes scenarios (R330)
- ✅ Tests fail initially (TDD red phase)
- ✅ Test assignments ready for effort planning



## 🔴🔴🔴 STATE COMPLETION CHECKLIST 🔴🔴🔴

**Execute these steps IN ORDER to properly complete WAITING_FOR_WAVE_TEST_PLAN:**

### ✅ Step 1: Complete State-Specific Work
**Refer to "Primary Actions" and "Mandatory Validations" sections above for state-specific tasks.**

Once all state work is complete, proceed to mandatory exit protocol:

---

### ✅ Step 2: Set Proposed Next State and Transition Reason
```bash
# Based on state work results, determine proposed next state
PROPOSED_NEXT_STATE="[DETERMINE_FROM_STATE_LOGIC]"
TRANSITION_REASON="WAITING_FOR_WAVE_TEST_PLAN complete - [accomplishment description]"
echo "Proposed next state: $PROPOSED_NEXT_STATE"
echo "Transition reason: $TRANSITION_REASON"
```

---

### ✅ Step 3: Spawn State Manager for SHUTDOWN_CONSULTATION
```bash
# State Manager validates transition and updates state files (SF 3.0 Pattern)
echo "🔄 Spawning State Manager for SHUTDOWN_CONSULTATION..."

# Prepare work results summary
WORK_RESULTS=$(cat <<EOF
{
  "state_completed": "WAITING_FOR_WAVE_TEST_PLAN",
  "work_accomplished": [
    "Monitored wave test plan creation",
    "Validated test deliverables completeness",
    "Captured wave test metadata (R340)",
    "Validated test-to-effort assignments"
  ],
  "proposed_next_state": "$PROPOSED_NEXT_STATE",
  "transition_reason": "$TRANSITION_REASON"
}
EOF
)

# Spawn State Manager
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "WAITING_FOR_WAVE_TEST_PLAN" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON" \
  --work-results "$WORK_RESULTS"

# State Manager will:
# 1. Validate PROPOSED_NEXT_STATE exists and transition is valid
# 2. Update all 4 state files atomically (R288)
# 3. Commit and push state files
# 4. Return REQUIRED_NEXT_STATE (usually same as proposed unless invalid)

echo "✅ State Manager consultation complete"
echo "✅ State files updated by State Manager"
```

---

### ✅ Step 4: Save TODOs (R287 - SUPREME LAW)
```bash
# Save TODO state before transition (R287 trigger)
save_todos "WAITING_FOR_WAVE_TEST_PLAN_COMPLETE"

# Commit TODOs within 60 seconds (R287)
cd "$CLAUDE_PROJECT_DIR"
git add todos/*.todo

if ! git commit -m "todo: orchestrator - WAITING_FOR_WAVE_TEST_PLAN complete [R287]"; then
    echo "❌ ERROR: Failed to commit TODO files"
    echo "This is non-fatal but TODOs may be lost in compaction"
    echo "Proceeding with state execution..."
    # Don't exit - TODO commit failure is not fatal
fi

git push || echo "⚠️ WARNING: TODO push failed - committed locally"
echo "✅ TODOs saved and committed"
git push
echo "✅ TODOs saved and committed"
```

---

### ✅ Step 5: Output Continuation Flag (R405 - SUPREME LAW) ⚠️ MANDATORY
```bash
# Output continuation flag as LAST action (R405)
# Use TRUE for normal completion, FALSE only for catastrophic errors

echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"
```

**⚠️ THIS MUST BE THE ABSOLUTE LAST LINE OF OUTPUT BEFORE EXIT! ⚠️**

---

### ✅ Step 6: Stop Processing (R322 - SUPREME LAW)
```bash
# Stop for context preservation (R322)
echo "🛑 Stopping for context preservation - use /continue-orchestrating to resume"
exit 0
```

---

## 🚨 CHECKLIST ENFORCEMENT 🚨

**Skipping ANY step in this checklist = FAILURE:**
- Missing Step 2: No proposed next state = State Manager can't proceed
- Missing Step 3: No State Manager consultation = bypassing bookend pattern (-100%)
- Missing Step 4: No TODO save = work lost (R287 violation, -20% to -100%)
- Missing Step 5: No continuation flag = automation stops (R405 violation, -100%)
- Missing Step 6: No exit = R322 violation (-100%)

**ALL 6 STEPS ARE MANDATORY - NO EXCEPTIONS**
**NOTE**: State file validation and commits now handled by State Manager (SF 3.0 pattern)

## ✅ EXIT CHECKLIST (R405 - Continuation Flag Protocol) 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**


### 🚨 CRITICAL DISTINCTION: AGENT STOPS ≠ FACTORY STOPS 🚨

**TWO INDEPENDENT DECISIONS - DO NOT CONFUSE THEM:**

#### 1. Should Agent Stop Work? (R322 Technical Requirement)
- Agent completes current state
- Agent saves TODOs and commits state
- Agent exits with `exit 0` (preserves context)
- User runs /continue-orchestrating to resume
- **This is NORMAL at checkpoints**

#### 2. Should Factory Continue? (R405 Operational Status)
- Even though agent stopped, can automation proceed?
- TRUE = Healthy completion, automation can continue
- FALSE = Catastrophic failure, must halt everything
- **R322 checkpoints = TRUE (99.9% of cases)**

### THE PATTERN AT R322 CHECKPOINTS (SF 3.0)

```bash
# 1. Complete state work
echo "✅ State work complete"

# 2. Set proposed next state
PROPOSED_NEXT_STATE="NEXT_STATE"
TRANSITION_REASON="State work complete"

# 3. Spawn State Manager for state transition
/spawn state-manager SHUTDOWN_CONSULTATION \
  --current-state "CURRENT_STATE" \
  --proposed-next-state "$PROPOSED_NEXT_STATE" \
  --transition-reason "$TRANSITION_REASON"
# State Manager updates all 4 state files atomically

# 4. Save TODOs
save_todos "R322_CHECKPOINT"

# 5. Factory continues (operational status)
echo "CONTINUE-SOFTWARE-FACTORY=TRUE REASON=STATE_COMPLETE"

# 6. Agent stops (technical requirement)
exit 0
```

**Both happen together! Agent stops AND factory continues!**

### WHEN TO USE EACH FLAG VALUE

**TRUE (99.9%):**
- ✅ R322 checkpoint reached
- ✅ State work completed successfully
- ✅ Ready for /continue-orchestrating
- ✅ Waiting for user to continue (NORMAL)
- ✅ Plan ready for review (agent done, factory proceeds)

**FALSE (0.1%):**
- ❌ CATASTROPHIC unrecoverable error
- ❌ Data corruption spreading
- ❌ Critical security violation
- ❌ NOT for R322 checkpoints
- ❌ NOT for user review needs
### 🚨 WAITING STATE PATTERN - CRITICAL UNDERSTANDING 🚨

**This is a WAITING state. Common source of incorrect FALSE usage!**

**WRONG interpretation:**
> "R322 mandates stop before transition"
> "State work is complete (validation done)"
> "User needs to invoke /continue-orchestrating"
> "Therefore I must set CONTINUE-SOFTWARE-FACTORY=FALSE"

**CORRECT interpretation:**
> "R322 checkpoint is NORMAL procedure for context preservation"
> "State work completed successfully = NORMAL outcome"
> "Waiting for /continue is DESIGNED user experience"
> "System KNOWS next step from state file"
> "NO manual intervention required, just normal continuation"
> "Therefore set CONTINUE-SOFTWARE-FACTORY=TRUE"

**The key distinction:**
- **Stopping inference** (`exit 0`) = Context management (ALWAYS at R322 points)
- **Continuation flag** = Can automation proceed? (TRUE unless catastrophic failure)

**ONLY use FALSE if:**
- ❌ The thing we're waiting for completely disappeared (agents crashed with no recovery)
- ❌ Results arrived but are completely corrupted/unreadable
- ❌ State file corruption prevents determining what to wait for
- ❌ System deadlock with no automated resolution
**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**


## Related Rules
- R341: Test-Driven Development Enforcement
- R330: Demo Planning Requirements
- R291: Integration Demo Requirement
- R211: Wave Implementation Planning
