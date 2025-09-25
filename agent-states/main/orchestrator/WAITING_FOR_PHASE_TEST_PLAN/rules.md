# WAITING_FOR_PHASE_TEST_PLAN State Rules

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

## 📋 PRIMARY DIRECTIVES FOR WAITING_FOR_PHASE_TEST_PLAN STATE

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
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-requirements.md`
   - Criticality: SUPREME - State updates required for all transitions
   - Summary: MUST update orchestrator-state.json before EVERY state transition

4. **🔴🔴🔴 R322 Part A** - Mandatory Stop After Spawn States
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-mandatory-stop-after-spawn.md`
   - Criticality: SUPREME LAW - Must stop after spawning
   - Summary: ALL spawn states require STOP after spawning agents

### State-Specific Rules:

5. **🔴🔴🔴 R233** - Immediate Action On State Entry
   - File: `$CLAUDE_PROJECT_DIR/rule-library/R233-immediate-action-on-state-entry.md`
   - Criticality: SUPREME LAW - Must act immediately on entering state
   - Summary: WAITING states require active monitoring, not passive waiting

## State Purpose
Monitor Code Reviewer creating phase-level tests and demo plans. Ensure all tests are ready BEFORE implementation planning begins (TDD enforcement).

## Entry Conditions
- Code Reviewer spawned for phase test planning
- Current state is SPAWN_CODE_REVIEWER_PHASE_TEST_PLANNING
- Waiting for test deliverables

## Required Actions

### 1. Monitor Test Creation Progress
```bash
# Check for test deliverables
check_phase_test_deliverables() {
    local phase_dir="phase-tests/phase-${PHASE_NUM}"
    
    echo "🔍 Checking for phase test deliverables..."
    
    # Required files
    local files_needed=(
        "PHASE-TEST-PLAN.md"
        "PHASE-TEST-HARNESS.sh"
        "PHASE-DEMO-PLAN.md"
    )
    
    for file in "${files_needed[@]}"; do
        if [ -f "$phase_dir/$file" ]; then
            echo "✅ Found: $file"
        else
            echo "⏳ Waiting for: $file"
            return 1
        fi
    done
    
    # Check for actual test files
    if ls "$phase_dir/tests/phase/"*.test.* >/dev/null 2>&1; then
        echo "✅ Found test files"
    else
        echo "⏳ Waiting for test files"
        return 1
    fi
    
    echo "✅ All phase test deliverables ready!"
    return 0
}
```

### 2. Validate Test Completeness
```bash
validate_phase_tests() {
    local phase_dir="phase-tests/phase-${PHASE_NUM}"
    
    echo "🧪 Validating phase test completeness..."
    
    # Test plan must reference architecture
    if grep -q "Architecture Coverage" "$phase_dir/PHASE-TEST-PLAN.md"; then
        echo "✅ Test plan covers architecture"
    else
        echo "❌ Test plan missing architecture coverage"
        return 1
    fi
    
    # Test harness must be executable
    if [ -x "$phase_dir/PHASE-TEST-HARNESS.sh" ]; then
        echo "✅ Test harness is executable"
    else
        echo "❌ Test harness not executable"
        chmod +x "$phase_dir/PHASE-TEST-HARNESS.sh"
    fi
    
    # Demo plan must include scenarios (R330)
    if grep -q "Demo Scenarios" "$phase_dir/PHASE-DEMO-PLAN.md"; then
        echo "✅ Demo scenarios defined"
    else
        echo "❌ Demo scenarios missing (R330 violation)"
        return 1
    fi
    
    # Tests should fail (no implementation yet)
    echo "🔴 Running tests (expecting failures - no implementation yet)..."
    if cd "$phase_dir" && ./PHASE-TEST-HARNESS.sh; then
        echo "⚠️ WARNING: Tests passing without implementation?"
    else
        echo "✅ Tests failing as expected (TDD - red phase)"
    fi
}
```

### 3. Capture Test Metadata (R340 Compliance)
```bash
# When Code Reviewer reports test locations
capture_test_metadata() {
    local phase_num="${PHASE_NUM}"
    local phase_key="phase${phase_num}"
    
    echo "📋 Capturing phase test metadata for state tracking..."
    
    # Update test_plans section with reported metadata
    jq --arg key "$phase_key" \
       --arg test_plan "/phase-tests/phase-${phase_num}/PHASE-TEST-PLAN.md" \
       --arg harness "/phase-tests/phase-${phase_num}/PHASE-TEST-HARNESS.sh" \
       --arg demo "/phase-tests/phase-${phase_num}/PHASE-DEMO-PLAN.md" \
       --arg test_dir "/phase-tests/phase-${phase_num}/tests/phase" \
       --arg created "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
       '.test_plans.phase[$key] = {
            "test_plan_path": $test_plan,
            "test_harness_path": $harness,
            "demo_plan_path": $demo,
            "test_dir": $test_dir,
            "created_at": $created,
            "created_by": "code-reviewer",
            "phase": '$phase_num',
            "status": "active",
            "tdd_phase": "red"
        }' orchestrator-state.json > tmp.json
    
    mv tmp.json orchestrator-state.json
    
    echo "✅ Test metadata captured in state file"
    echo "   SW Engineers can now find tests at: $test_dir"
}
```

### 4. Transition When Ready
```bash
# When all tests are ready
transition_to_implementation_planning() {
    echo "📋 Phase tests ready - proceeding to implementation planning"
    
    # Capture test metadata if not already done
    capture_test_metadata
    
    # Update state file
    jq '.current_state = "SPAWN_CODE_REVIEWER_PHASE_IMPL"' orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
    
    # Record test planning completion
    jq '.phase_test_planning.completed_at = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' \
        orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
    
    git add -A
    git commit -m "tdd: phase ${PHASE_NUM} tests ready - proceeding to implementation planning
    
- Captured test locations in state file (R340)
- Tests available at tracked locations
- SW Engineers can find tests from state"
    git push
}
```

## Exit Conditions
- All test deliverables created and validated
- Transition to SPAWN_CODE_REVIEWER_PHASE_IMPL
- Tests documented and ready for implementation

## Success Criteria
- ✅ PHASE-TEST-PLAN.md exists and covers architecture
- ✅ PHASE-TEST-HARNESS.sh is executable
- ✅ Test files created in tests/phase/
- ✅ PHASE-DEMO-PLAN.md includes scenarios (R330)
- ✅ Tests fail initially (TDD red phase)

## Related Rules
- R341: Test-Driven Development Enforcement
- R330: Demo Planning Requirements
- R291: Integration Demo Requirement
- R211: Implementation Planning