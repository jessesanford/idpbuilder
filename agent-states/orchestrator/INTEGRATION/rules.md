# Orchestrator - INTEGRATION State Rules

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


## 🔴🔴🔴 CRITICAL: CHECK YOUR CURRENT STATE FIRST! 🔴🔴🔴

**BEFORE ANYTHING ELSE, CHECK orchestrator-state.json:**
- If `current_state: INTEGRATION` → You are ALREADY in INTEGRATION state!
- If already in INTEGRATION → Skip to rule reading, then START INTEGRATION WORK IMMEDIATELY
- Integration is a VERB - it means "DO INTEGRATION NOW"

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED INTEGRATION STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INTEGRATION
echo "$(date +%s) - Rules read and acknowledged for INTEGRATION" > .state_rules_read_orchestrator_INTEGRATION
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INTEGRATION WORK UNTIL RULES ARE READ:
- ❌ Start create integration branch
- ❌ Start merge effort branches
- ❌ Start resolve conflicts
- ❌ Start run tests
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
   ❌ WRONG: "I acknowledge all INTEGRATION rules"
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

### ✅ CORRECT PATTERN FOR INTEGRATION:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute INTEGRATION work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY INTEGRATION work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute INTEGRATION work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with INTEGRATION work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY INTEGRATION work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## ⚠️⚠️⚠️ MANDATORY RULE READING AND ACKNOWLEDGMENT ⚠️⚠️⚠️

**YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. YOUR READ TOOL CALLS ARE BEING MONITORED.**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:
1. Fake acknowledgment without reading
2. Bulk acknowledgment
3. Reading from memory

### ✅ CORRECT PATTERN:
1. READ each rule file
2. Acknowledge individually with rule number and description

## 📋 PRIMARY DIRECTIVES FOR INTEGRATION STATE

### 🚨🚨🚨 R250 - Integration Isolation Requirement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R250-integration-isolation-requirement.md`
**Criticality**: BLOCKING - Integration must use separate target clone
**Summary**: Integration must happen under /efforts/ directory structure

### 🚨🚨🚨 R034 - Integration Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R034-integration-requirements.md`
**Criticality**: BLOCKING - Required for wave approval
**Summary**: Complete integration protocol with testing and validation

### 🚨🚨🚨 R296 - Deprecated Branch Marking Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R296-deprecated-branch-marking-protocol.md`
**Criticality**: BLOCKING - Prevents integration of wrong branches
**Summary**: Check for and prevent integration of deprecated split branches

### 🚨🚨🚨 R014 - Branch Naming Convention
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**Criticality**: BLOCKING - Mandatory project prefix for all branches
**Summary**: Use project prefix for all integration branches

### 🚨🚨🚨 R271 - Mandatory Production-Ready Validation
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R271-mandatory-production-ready-validation.md`
**Criticality**: BLOCKING - Full checkouts required for integration
**Summary**: Integration must use full repository clones, no sparse checkouts

### 🚨🚨🚨 R258 - Mandatory Wave Review Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R258-mandatory-wave-review-report.md`
**Criticality**: BLOCKING - Required for wave completion
**Summary**: Architect must create wave review report after integration

### 🚨🚨🚨 R006 - Orchestrator NEVER Writes Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md`
**Criticality**: BLOCKING - Any code operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator coordinates but NEVER implements or fixes code

### 🚨🚨🚨 R329 - Orchestrator NEVER Performs Git Merges [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R329-orchestrator-never-performs-merges.md`
**Criticality**: BLOCKING - Any merge operation = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator MUST spawn Integration Agent for ALL merges - NO EXCEPTIONS

### 🚨🚨🚨 R319 - Orchestrator NEVER Measures or Assesses Code [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R319-orchestrator-never-measures-code.md`
**Criticality**: BLOCKING - Any technical assessment = -100% IMMEDIATE FAILURE
**Summary**: Orchestrator NEVER runs tests, builds, or validates - spawns specialists

### 🚨🚨🚨 R269 - Code Reviewer Merge Plan No Execution
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R269-code-reviewer-merge-plan-no-execution.md`
**Criticality**: BLOCKING - Code Reviewer only plans, never executes
**Summary**: Code Reviewer creates plan, Integration Agent executes

### 🚨🚨🚨 R260 - Integration Agent Core Requirements
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R260-integration-agent-core-requirements.md`
**Criticality**: BLOCKING - Integration Agent must acknowledge INTEGRATION_DIR
**Summary**: Integration Agent must set and use INTEGRATION_DIR variable

### 🔴🔴🔴 R321 - Immediate Backport During Integration Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R321-immediate-backport-during-integration.md`
**Criticality**: SUPREME LAW - Immediate backporting required
**Summary**: ANY fix during integration MUST be immediately backported to source branches before continuing

### 🚨🚨🚨 R280 - Main Branch Protection Protocol [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R280-main-branch-protection.md`
**Criticality**: BLOCKING - Direct commits to main/master are forbidden
**Summary**: All changes must go through PR process with proper reviews

### 🚨🚨🚨 R307 - Independent Branch Mergeability [BLOCKING]
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R307-independent-branch-mergeability.md`
**Criticality**: BLOCKING - Must verify branches are mergeable before attempting
**Summary**: Check for conflicts and mergeability before integration operations

## 🚨 INTEGRATION IS A VERB - COORDINATE INTEGRATION NOW! 🚨

### 🔴🔴🔴 CRITICAL: YOU ARE ALREADY IN INTEGRATION STATE! 🔴🔴🔴

**If current_state = "INTEGRATION" in orchestrator-state.json, you MUST:**
1. **IMMEDIATELY** check if integration infrastructure exists
2. **NO ANNOUNCEMENTS** - just start checking
3. **NO WAITING** - coordination begins NOW

### IMMEDIATE ACTIONS UPON ENTERING INTEGRATION

**THE MOMENT YOU SEE current_state: INTEGRATION, YOU MUST:**
1. Check if integration infrastructure exists NOW (no delay)
2. If NO infrastructure: Transition to SETUP_INTEGRATION_INFRASTRUCTURE
3. If infrastructure EXISTS: Transition to SPAWN_CODE_REVIEWER_MERGE_PLAN
4. Update state file with the appropriate next state
5. Stop per R322 for state transition

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in INTEGRATION" [stops]
- ❌ "Successfully entered INTEGRATION state" [waits]
- ❌ "Ready to start integrating" [pauses]
- ❌ "I'm in INTEGRATION state" [does nothing]
- ❌ "Preparing to setup integration..." [delays]
- ❌ "I see we're in INTEGRATION state..." [announces]
- ❌ Creating infrastructure yourself (INTEGRATION only coordinates!)

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "INTEGRATION STATE: Checking for existing integration infrastructure..."
- ✅ "No infrastructure found, transitioning to SETUP_INTEGRATION_INFRASTRUCTURE..."
- ✅ "Infrastructure exists, transitioning to SPAWN_CODE_REVIEWER_MERGE_PLAN..."

### ⚠️⚠️⚠️ RULE R020 - State Transitions
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R020-state-transitions.md`

## State Context

### 🔴🔴🔴 INTEGRATION IS A COORDINATION-ONLY STATE 🔴🔴🔴

**THIS STATE ONLY DETERMINES THE NEXT TRANSITION - IT DOES NOT CREATE INFRASTRUCTURE!**

The INTEGRATION state is a decision point that:
1. **CHECKS** if integration infrastructure exists
2. **TRANSITIONS** to SETUP_INTEGRATION_INFRASTRUCTURE if no infrastructure
3. **TRANSITIONS** to SPAWN_CODE_REVIEWER_MERGE_PLAN if infrastructure exists

**THIS STATE NEVER:**
- ❌ Creates integration workspace itself
- ❌ Sets up branches or directories itself
- ❌ Performs any actual integration work

You are the COORDINATOR of integration flow. Your ONLY responsibilities in this state:
1. **CHECK** if integration infrastructure already exists
2. **DECIDE** which state to transition to based on infrastructure status
3. **UPDATE** state file with appropriate next state
4. **STOP** per R322 for state transition

**YOU MUST NEVER (R329 + R006 ENFORCEMENT):**
- ❌ Execute git merges yourself (R329 VIOLATION = IMMEDIATE FAILURE)
- ❌ Resolve merge conflicts yourself (R329 VIOLATION)
- ❌ Perform "simple" merges yourself (R329: NO EXCEPTIONS)
- ❌ Run build commands yourself (R006 VIOLATION)
- ❌ Execute test suites yourself (R319 VIOLATION)
- ❌ Fix code issues yourself (R006 VIOLATION)
- ❌ Apply patches or cherry-picks yourself (R329 VIOLATION)

## 🔴🔴🔴 CRITICAL: RE-RUNNING INTEGRATION AFTER FIXES 🔴🔴🔴

**When returning to INTEGRATION after fixes (from MONITORING_FIX_PROGRESS flow):**
1. **VERIFY** all fixes are in source effort branches (R321 mandatory)
2. **DELETE** the old integration workspace completely
3. **CREATE** a new clean integration workspace
4. **RE-EXECUTE** the entire integration plan from scratch
5. **MERGE** ALL branches again (not just fixed ones)
6. **RUN** full test suite on newly integrated code

**NEVER DO THESE (AUTOMATIC FAILURE):**
- ❌ Reuse the old integration workspace
- ❌ Manually copy fixed files over
- ❌ Cherry-pick fixes without proper merging
- ❌ Skip any branches in the merge sequence
- ❌ Assume previous merges are still valid
- ❌ Apply fixes directly to integration branch (R321 violation)

**The integration MUST be completely fresh after fixes!**

## 🔴🔴🔴 R321 ENFORCEMENT: IMMEDIATE BACKPORTING 🔴🔴🔴

### CRITICAL: Integration Branches Are READ-ONLY
**Per R321, if ANY issue is found during integration:**

1. **STOP IMMEDIATELY** - Do not continue integration
2. **IDENTIFY** which effort branch needs fixing
3. **SPAWN SW ENGINEER** to fix the source effort branch
4. **WAIT** for fix to be applied and pushed to effort branch
5. **VERIFY** effort branch works independently
6. **ONLY THEN** retry integration with fixed source

### Validation Before Completing Integration
```bash
# R321 MANDATORY: Verify all source branches work independently
validate_source_branches_before_completion() {
    echo "🔍 R321 Validation: Checking all source branches"
    
    for effort_dir in /efforts/phase${PHASE}/wave${WAVE}/*/; do
        cd "$effort_dir"
        BRANCH=$(git branch --show-current)
        
        echo "Testing $BRANCH independently..."
        if ! npm run build; then
            echo "❌ R321 VIOLATION: $BRANCH doesn't build!"
            echo "Must fix in source before integration can complete"
            exit 1
        fi
        
        if ! npm test; then
            echo "❌ R321 VIOLATION: $BRANCH tests fail!"
            echo "Must fix in source before integration can complete"
            exit 1
        fi
    done
    
    echo "✅ R321 Validated: All sources work independently"
}
```

### Detection of Integration Branch Modifications
```bash
# R321 ENFORCEMENT: Check for forbidden direct edits
check_integration_branch_purity() {
    NON_MERGE=$(git log --oneline --no-merges origin/main..HEAD)
    
    if [ -n "$NON_MERGE" ]; then
        echo "🔴🔴🔴 R321 VIOLATION DETECTED!"
        echo "Direct edits found in integration branch:"
        echo "$NON_MERGE"
        echo "ALL fixes must go to source branches!"
        exit 1
    fi
}
```

## 🔴🔴🔴 CRITICAL: INTEGRATION LOCATION 🔴🔴🔴

### 🚨🚨🚨 RULE R250 - Integration Isolation Requirements
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R250-integration-isolation-requirement.md`
**ALSO SEE**: `$CLAUDE_PROJECT_DIR/agent-states/orchestrator/INTEGRATION/RULE-R250-INTEGRATION-ISOLATION.md`

### 🚨🚨🚨 RULE R034 - Integration Protocol Requirements
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R034-integration-requirements.md`

## Branch Creation Strategy

### 🚨🚨🚨 RULE R014 - Branch Naming Convention
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R014-branch-naming-convention.md`
**NOTE**: Use utilities/branch-naming-helpers.sh for automatic prefix handling

## Integration Coordination Protocol

### 🔴🔴🔴 CRITICAL: INTEGRATION NOW DELEGATES INFRASTRUCTURE CREATION 🔴🔴🔴
**Infrastructure creation has been moved to SETUP_INTEGRATION_INFRASTRUCTURE state (R308 enforced)**

```bash
# Check if integration infrastructure exists
check_integration_infrastructure() {
    PHASE=$(jq '.current_phase' orchestrator-state.json)
    WAVE=$(jq '.current_wave' orchestrator-state.json)
    
    echo "🔍 Checking for existing integration infrastructure..."
    
    # Check state file for infrastructure metadata
    INFRA_EXISTS=$(jq '.integration_infrastructure // null' orchestrator-state.json)
    
    if [ "$INFRA_EXISTS" = "null" ]; then
        echo "❌ No integration infrastructure found in state file"
        echo "➡️ Must transition to SETUP_INTEGRATION_INFRASTRUCTURE"
        return 1
    fi
    
    # Verify infrastructure directory exists
    INFRA_DIR=$(jq -r '.integration_infrastructure.directory' orchestrator-state.json)
    if [ ! -d "$INFRA_DIR" ]; then
        echo "❌ Infrastructure directory missing: $INFRA_DIR"
        echo "➡️ Must recreate infrastructure via SETUP_INTEGRATION_INFRASTRUCTURE"
        return 1
    fi
    
    # Verify branch exists on remote
    INFRA_BRANCH=$(jq -r '.integration_infrastructure.branch' orchestrator-state.json)
    if ! git ls-remote --heads origin "$INFRA_BRANCH" > /dev/null 2>&1; then
        echo "❌ Infrastructure branch not on remote: $INFRA_BRANCH"
        echo "➡️ Must recreate infrastructure via SETUP_INTEGRATION_INFRASTRUCTURE"
        return 1
    fi
    
    # Verify R308 compliance
    BASE_BRANCH=$(jq -r '.integration_infrastructure.base_branch' orchestrator-state.json)
    echo "✅ Integration infrastructure exists and is valid"
    echo "   - Directory: $INFRA_DIR"
    echo "   - Branch: $INFRA_BRANCH"
    echo "   - Base Branch: $BASE_BRANCH (R308 compliant)"
    echo "   - Type: $(jq -r '.integration_infrastructure.type' orchestrator-state.json)"
    echo "➡️ Can proceed to SPAWN_CODE_REVIEWER_MERGE_PLAN"
    return 0
}

# Main coordination logic
coordinate_integration() {
    echo "🔧 INTEGRATION STATE: Coordinating integration process"
    echo "📋 This state determines next action but does NOT create infrastructure"
    
    # Load state helpers
    source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"
    
    if check_integration_infrastructure; then
        echo "📋 Infrastructure ready, proceeding to merge planning"
        
        # Update state to spawn Code Reviewer
        safe_state_transition "SPAWN_CODE_REVIEWER_MERGE_PLAN" "Infrastructure exists, need merge plan"
        
        echo "✅ Transitioning to SPAWN_CODE_REVIEWER_MERGE_PLAN"
    else
        echo "🏗️ Infrastructure needed, delegating to setup state"
        echo "🔴 R308 will be enforced in SETUP_INTEGRATION_INFRASTRUCTURE"
        
        # Update state to create infrastructure
        safe_state_transition "SETUP_INTEGRATION_INFRASTRUCTURE" "Need to create integration infrastructure with R308 compliance"
        
        echo "✅ Transitioning to SETUP_INTEGRATION_INFRASTRUCTURE"
    fi
    
    echo "✅ Integration coordination complete"
    echo "🛑 Stopping for state transition per R322"
}

# Execute coordination immediately
coordinate_integration
```

## Spawning Code Reviewer for Merge Plan

```bash
# Setup integration infrastructure first
PHASE=$(jq '.current_phase' orchestrator-state.json)
WAVE=$(jq '.current_wave' orchestrator-state.json)
INTEGRATION_DIR="/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"

# Ensure we're in integration directory
cd "$INTEGRATION_DIR"

# Create and push integration branch
git checkout -b "phase${PHASE}-wave${WAVE}-integration-$(date +%Y%m%d-%H%M%S)"
git push -u origin HEAD

# Spawn Code Reviewer for MERGE PLAN
Task: subagent_type="code-reviewer" \
      prompt="Create MERGE PLAN for Phase ${PHASE} Wave ${WAVE} integration.
      
      CRITICAL REQUIREMENTS:
      1. Use ONLY original effort branches - NO integration branches!
      2. Analyze branch bases to determine correct merge order
      3. Exclude 'too-large' branches, include only splits
      4. Create WAVE-MERGE-PLAN.md with exact merge instructions
      5. DO NOT execute merges - only plan them!
      
      CRITICAL LOCATION REQUIREMENT:
      - CD to integration directory FIRST: cd ${INTEGRATION_DIR}
      - Create WAVE-MERGE-PLAN.md IN the integration directory
      - Full path for the file: ${INTEGRATION_DIR}/WAVE-MERGE-PLAN.md
      
      Integration Directory: ${INTEGRATION_DIR}
      Target Branch: $(git branch --show-current)" \
      description="Create Wave ${WAVE} Merge Plan"
```

## Spawning Integration Agent for Execution

```bash
# After Code Reviewer creates MERGE PLAN
cd "$INTEGRATION_DIR"

# Verify merge plan exists
if [ ! -f "WAVE-MERGE-PLAN.md" ]; then
    echo "❌ Cannot spawn Integration Agent - no merge plan!"
    exit 1
fi

# Spawn Integration Agent
Task: subagent_type="integration-agent" \
      prompt="Execute integration merges for Phase ${PHASE} Wave ${WAVE}.
      
      CRITICAL REQUIREMENTS:
      1. You are in INTEGRATION_DIR: ${INTEGRATION_DIR}
      2. Acknowledge and set INTEGRATION_DIR variable
      3. Read and follow WAVE-MERGE-PLAN.md EXACTLY
      4. Execute merges in specified order
      5. Handle conflicts as directed in plan
      6. Run tests after each merge
      
      Your working directory has been set to: ${INTEGRATION_DIR}
      The merge plan is: WAVE-MERGE-PLAN.md" \
      description="Execute Wave ${WAVE} Integration"
```

## 🏗️ MANDATORY BUILD VERIFICATION (Via Code Reviewer)

**NO INTEGRATION IS COMPLETE WITHOUT A WORKING BUILD:**
1. Code Reviewer must verify build success
2. Code Reviewer must verify artifacts
3. Code Reviewer must capture build logs
4. Any build failures = integration incomplete
5. Build must be runnable/executable

```bash
# ✅ CORRECT: Spawn Code Reviewer for build verification
echo "🏗️ Build verification needed after integration"
echo "🚀 Spawning Code Reviewer to validate build..."

Task: subagent_type="code-reviewer" \
      prompt="Validate integrated code builds successfully. Run full build, verify artifacts, capture logs. Create BUILD-VALIDATION-REPORT.md" \
      workspace="$INTEGRATION_DIR" \
      description="Validate integration build"

# ❌❌❌ FORBIDDEN - Orchestrator CANNOT run builds!
# if npm run build; then  # VIOLATION OF R006/R319!
#     echo "Build successful"  # IMMEDIATE FAILURE!
# fi
```

## 🧪 MANDATORY TEST HARNESS

**EVERY INTEGRATION MUST HAVE A TEST HARNESS:**
1. Create test-harness.sh (or appropriate for language)
2. Must test the integrated functionality
3. Must be automated and repeatable
4. Must clearly show pass/fail
5. Results must be captured in logs

```bash
# Create test harness script
cat > test-harness.sh << 'EOF'
#!/bin/bash
# Wave Integration Test Harness
echo "🧪 Starting Wave Integration Test Suite"
echo "========================================="

# Unit tests
echo "📦 Running unit tests..."
if npm test 2>&1 | tee unit-tests.log; then
    echo "✅ Unit tests passed"
else
    echo "❌ Unit tests failed"
    exit 1
fi

# Integration tests
echo "🔗 Running integration tests..."
if npm run test:integration 2>&1 | tee integration-tests.log; then
    echo "✅ Integration tests passed"
else
    echo "❌ Integration tests failed"
    exit 1
fi

# Feature verification
echo "🎯 Verifying new features..."
# Add feature-specific tests here
./verify-wave-features.sh

echo "========================================="
echo "✅ ALL TESTS PASSED - Integration verified!"
EOF

chmod +x test-harness.sh
./test-harness.sh
```

## 🎬 MANDATORY DEMO

**DEMONSTRATE THE WORKING INTEGRATION:**
1. Create demo script or documentation
2. Show actual functionality working
3. Capture output/screenshots/video as appropriate
4. Prove the integration delivers value
5. Document what features are now available

```bash
# Create demo documentation
cat > WAVE-DEMO.md << 'EOF'
# Wave ${WAVE} Integration Demo

## Build Status
- Build: ✅ PASSING
- Tests: ✅ ALL PASSING
- Integration: ✅ COMPLETE

## Features Demonstrated
1. [Feature 1]: Working implementation
2. [Feature 2]: Integration verified
3. [Feature 3]: Tests passing

## How to Run Demo
```bash
# Start the application
npm start

# Run the demo script
./demo-wave-features.sh

# Verify outputs
curl http://localhost:3000/api/new-feature
```

## Evidence
- Build log: build.log
- Test results: test-results.log
- Screenshots: demos/wave-${WAVE}/
EOF

# Create demo script
cat > demo-wave-features.sh << 'EOF'
#!/bin/bash
echo "🎬 Wave Integration Demo"
echo "Demonstrating new functionality..."
# Add actual demo commands here
EOF
chmod +x demo-wave-features.sh
```

## Integration Validation

```yaml
# Update orchestrator-state.json
integration_records:
  phase{X}_wave{Y}:
    started_at: "2025-08-23T15:00:00Z"
    efforts_included:
      - "effort1-api-types"
      - "effort2-controller"
      - "effort3-webhooks"
    conflicts_detected: 0
    resolution_time: "0s"
    build_status: "PASS"
    build_artifacts: ["dist/", "build.log"]
    test_harness: "test-harness.sh"
    test_results:
      unit_tests: "PASS"
      integration_tests: "PASS"
      feature_tests: "PASS"
      build: "PASS"
    demo_completed: true
    demo_artifacts: ["WAVE-DEMO.md", "demo-wave-features.sh"]
    final_size_check: "742 lines"
    completed_at: "2025-08-23T15:05:22Z"
    grade: "PASS"
```

## Phase Transition Protocol

## Size Validation During Integration

```python
def validate_integration_size(branch_name):
    """Validate size limits maintained during integration"""
    
    result = subprocess.run([
        '$PROJECT_ROOT/tools/line-counter.sh',
        '-c', branch_name
    ], capture_output=True, text=True)
    
    total_lines = int(result.stdout.split()[-1])
    
    # Integration should not exceed sum of constituent efforts
    # with reasonable overhead for integration code
    max_allowed = calculate_effort_sum() * 1.05  # 5% overhead
    
    if total_lines > max_allowed:
        return {
            'valid': False,
            'total_lines': total_lines,
            'max_allowed': max_allowed,
            'action': 'SPLIT_INTEGRATION'
        }
    
    return {
        'valid': True,
        'total_lines': total_lines,
        'grade': 'PASS'
    }
```

## Integration Failure Recovery

If integration fails:
1. Stop all integration work
2. Transition to ERROR_RECOVERY state
3. Analyze failure cause
4. Create recovery plan
5. Execute fixes before retrying

## State Transitions

From INTEGRATION state:
- **SETUP_COMPLETE** → SPAWN_CODE_REVIEWER_MERGE_PLAN
- **MERGE_PLAN_READY** → SPAWN_INTEGRATION_AGENT  
- **INTEGRATION_COMPLETE** → WAVE_REVIEW (R258: Architect must create wave review report)
- **FAILURE** → ERROR_RECOVERY

New intermediate states:
- **SPAWN_CODE_REVIEWER_MERGE_PLAN** - Spawning Code Reviewer for merge planning
- **WAITING_FOR_MERGE_PLAN** - Waiting for Code Reviewer to complete plan
- **SPAWN_INTEGRATION_AGENT** - Spawning Integration Agent for execution
- **MONITORING_INTEGRATION** - Monitoring Integration Agent progress

**IMPORTANT**: The architect MUST create a wave review report (R258) with one of these decisions:
- PROCEED_NEXT_WAVE - Wave approved, start next wave
- PROCEED_PHASE_ASSESSMENT - Last wave complete, trigger phase assessment
- CHANGES_REQUIRED - Fixes needed before progression
- WAVE_FAILED - Major issues, cannot proceed

## R322 VIOLATION DETECTION

If you find yourself:
- Starting work for a new state without /continue-orchestrating
- Transitioning without stopping after state file commit
- Continuing after completing state work

**STOP IMMEDIATELY - You are violating R322!**


### 🔴🔴🔴 MANDATORY VALIDATION REQUIREMENT 🔴🔴🔴

**Per R288 and R324**: ALL state file updates MUST be validated before commit:

```bash
# After ANY update to orchestrator-state.json:
"$CLAUDE_PROJECT_DIR/tools/validate-state.sh" orchestrator-state.json || {
    echo "❌ State file validation failed!"
    exit 288
}
```

**Use helper functions for automatic validation:**
```bash
# Source the helper functions
source "$CLAUDE_PROJECT_DIR/utilities/state-file-update-functions.sh"

# Use safe functions that include validation:
safe_state_transition "NEW_STATE" "reason"
safe_update_field "field_name" "value"
```
