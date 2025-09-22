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
        }' orchestrator-state.json > tmp.json
    
    mv tmp.json orchestrator-state.json
    
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
    
    # Update state file
    jq '.current_state = "SPAWN_CODE_REVIEWER_WAVE_IMPL"' orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
    
    # Record wave test planning completion
    jq '.wave_test_planning.completed_at = "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"' \
        orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
    
    # Link tests to implementation planning
    jq '.wave_test_planning.test_directory = "wave-tests/phase-'${PHASE_NUM}'/wave-'${WAVE_NUM}'"' \
        orchestrator-state.json > tmp.json
    mv tmp.json orchestrator-state.json
    
    git add -A
    git commit -m "tdd: wave ${WAVE_NUM} tests ready - proceeding to implementation planning
    
- Captured wave test locations in state file (R340)
- Test-to-effort assignments tracked
- SW Engineers can find tests from state"
    git push
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

## Related Rules
- R341: Test-Driven Development Enforcement
- R330: Demo Planning Requirements
- R291: Integration Demo Requirement
- R211: Wave Implementation Planning