# Orchestrator - ERROR_RECOVERY State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED ERROR_RECOVERY STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_ERROR_RECOVERY
echo "$(date +%s) - Rules read and acknowledged for ERROR_RECOVERY" > .state_rules_read_orchestrator_ERROR_RECOVERY
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY ERROR_RECOVERY WORK UNTIL RULES ARE READ:
- ❌ Start diagnose errors
- ❌ Start recover from failures
- ❌ Start restart failed efforts
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
   ❌ WRONG: "I acknowledge all ERROR_RECOVERY rules"
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

### ✅ CORRECT PATTERN FOR ERROR_RECOVERY:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute ERROR_RECOVERY work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY ERROR_RECOVERY work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute ERROR_RECOVERY work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with ERROR_RECOVERY work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY ERROR_RECOVERY work before reading and acknowledging rules:**
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

## 🔴🔴🔴 COMMON ERROR_RECOVERY TRIGGERS (R291 ENFORCEMENT) 🔴🔴🔴

**YOU ARE IN ERROR_RECOVERY BECAUSE ONE OF THESE LIKELY OCCURRED:**

1. **BUILD GATE FAILURE** (R291 Violation)
   - Build compilation failed (make/npm build/cargo build returned non-zero)
   - No build artifacts produced (missing dist/build/target)
   - Build process crashed or timed out

2. **TEST GATE FAILURE** (R291 Violation)
   - Unit tests failed
   - Integration tests failed
   - E2E tests failed
   - ANY test suite returned non-zero exit code

3. **DEMO GATE FAILURE** (R291 Violation)
   - Demo script failed (exit code != 0)
   - Features don't work as expected
   - Demo cannot be run due to missing components

4. **INTEGRATION FAILURES**
   - Merge conflicts that couldn't be resolved
   - Dependency issues blocking build
   - Incompatible changes between efforts

5. **ASSESSMENT FAILURES**
   - Architect phase assessment failed
   - Wave review identified critical issues
   - Quality gates not met

**YOUR FIRST ACTION:** Check orchestrator-state.yaml for `error_recovery.reason` to understand why you're here!

## 📋 PRIMARY DIRECTIVES FOR ERROR_RECOVERY STATE

### 🚨🚨🚨 R019 - Error Recovery Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R019-error-recovery-protocol.md`
**Criticality**: BLOCKING - Must follow recovery protocol
**Summary**: Systematic error assessment, preservation, and recovery

### 🔴🔴🔴 R021 - Orchestrator Never Stops (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R021-orchestrator-never-stops.md`
**Criticality**: SUPREME LAW - Violation = -100% failure
**Summary**: Continue operations until success or explicit stop

### ⚠️⚠️⚠️ R156 - Error Recovery Time Targets
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R156-error-recovery-time-targets.md`
**Criticality**: CRITICAL - Time targets for recovery by severity
**Summary**: CRITICAL <30min, HIGH <60min, MEDIUM <2hrs, LOW <4hrs

### ⚠️⚠️⚠️ R010 - Wrong Location Handling
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R010-wrong-location-handling.md`
**Criticality**: CRITICAL - Never attempt to fix location errors
**Summary**: Stop immediately, preserve state, wait for manual correction

### 🚨🚨🚨 R258 - Mandatory Wave Review Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R258-mandatory-wave-review-report.md`
**Criticality**: BLOCKING - Required for wave completion
**Summary**: Must read wave review report from specified location

### 🚨🚨🚨 R257 - Mandatory Phase Assessment Report
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R257-mandatory-phase-assessment-report.md`
**Criticality**: BLOCKING - Required for phase completion
**Summary**: Must read phase assessment report from specified location

### 🚨🚨🚨 R259 - Mandatory Phase Integration After Fixes
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R259-mandatory-phase-integration-after-fixes.md`
**Criticality**: BLOCKING - Must create integration branch after fixes
**Summary**: Return to PHASE_INTEGRATION after phase assessment fixes

### 🔴🔴🔴 R300 - Comprehensive Fix Management Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R300-comprehensive-fix-management-protocol.md`
**Criticality**: SUPREME LAW - Violation = -100% AUTOMATIC FAILURE  
**Summary**: ALL fixes MUST be applied to effort branches, NEVER to integration branches

## 🚨 ERROR_RECOVERY IS A VERB - START ERROR RECOVERY IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING ERROR_RECOVERY

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Diagnose the error from state file NOW
2. Create recovery plan immediately
3. Check TodoWrite for pending items and process them
4. Begin executing recovery steps

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in ERROR_RECOVERY" [stops]
- ❌ "Successfully entered ERROR_RECOVERY state" [waits]
- ❌ "Ready to start error recovery" [pauses]
- ❌ "I'm in ERROR_RECOVERY state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering ERROR_RECOVERY, Diagnose the error from state file NOW..."
- ✅ "START ERROR RECOVERY, create recovery plan immediately..."
- ✅ "ERROR_RECOVERY: Begin executing recovery steps..."

## State Context
You are recovering from critical errors that have blocked normal operation flow.


### 🔴🔴🔴 RULE R021 - Orchestrator Never Stops (SUPREME LAW)
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R021-orchestrator-never-stops.md`
**CRITICAL FOR ERROR_RECOVERY**: Never stop recovering until resolved or explicit stop

## Recovery Decision Matrix

```python
def classify_error_and_strategy(error_data):
    """Classify error and determine recovery strategy"""
    
    error_type = error_data['type']
    severity = error_data['severity']
    affected_scope = error_data['scope']
    
    strategies = {
        'SIZE_LIMIT_VIOLATION': {
            'severity': 'CRITICAL',
            'strategy': 'IMMEDIATE_SPLIT',
            'actions': [
                'Stop all implementation work',
                'Spawn Code Reviewer for split analysis',
                'Create split plan',
                'Execute split with size validation'
            ],
            'target_time': '30min'
        },
        
        'WAVE_REVIEW_CHANGES_REQUIRED': {
            'severity': 'HIGH',
            'strategy': 'ARCHITECT_DIRECTED_FIXES',
            'actions': [
                'Read wave review report per R258',
                'Extract Issues Identified section',
                'Parse Required Fixes from report',
                'Map each fix to appropriate agent type',
                'Spawn agents for each required fix',
                'Track completion against report requirements',
                'Return to INTEGRATION for re-review'
            ],
            'target_time': '120min',
            'report_location': 'wave-reviews/phase{N}/wave{W}/PHASE-{N}-WAVE-{W}-REVIEW-REPORT.md'
        },
        
        'PHASE_ASSESSMENT_NEEDS_WORK': {
            'severity': 'CRITICAL',
            'strategy': 'PHASE_LEVEL_REMEDIATION',
            'actions': [
                'Read phase assessment report per R257',
                'Extract all Issues Identified',
                'Parse Required Fixes section',
                'Create comprehensive fix plan',
                'Spawn multiple agents if needed',
                'Coordinate cross-wave fixes',
                'Track against assessment criteria',
                'Return to PHASE_INTEGRATION for integration branch creation (R259)'
            ],
            'target_time': '240min',
            'report_location': 'phase-assessments/phase{N}/PHASE-{N}-ASSESSMENT-REPORT.md'
        },
        
        'INTEGRATION_FAILURE': {
            'severity': 'HIGH',
            'strategy': 'CONFLICT_RESOLUTION',
            'actions': [
                'Analyze merge conflicts',
                'Spawn SW Engineers to fix IN EFFORT BRANCHES (R300)',
                'NEVER fix in integration branch - will be lost!',
                'Validate fixes in effort branches',
                'Retry integration with updated effort branches'
            ],
            'target_time': '60min'
        },
        
        'TEST_SUITE_FAILURE': {
            'severity': 'CRITICAL',
            'strategy': 'IMMEDIATE_FIX',
            'actions': [
                'Isolate failing tests',
                'Determine root cause',
                'Spawn SW Engineer for urgent fix',
                'Validate fix before resuming'
            ],
            'target_time': '30min'
        },
        
        'ARCHITECTURE_VIOLATION': {
            'severity': 'HIGH',
            'strategy': 'ARCHITECT_CONSULTATION',
            'actions': [
                'Document violation details',
                'Spawn Architect for guidance',
                'Create remediation plan',
                'Execute with continuous validation'
            ],
            'target_time': '60min'
        },
        
        'AGENT_COMMUNICATION_FAILURE': {
            'severity': 'HIGH',
            'strategy': 'AGENT_RESTART',
            'actions': [
                'Identify failed agent',
                'Save agent state',
                'Restart with full context',
                'Resume from last checkpoint'
            ],
            'target_time': '45min'
        }
    }
    
    return strategies.get(error_type, {
        'severity': 'UNKNOWN',
        'strategy': 'MANUAL_ANALYSIS',
        'actions': ['Escalate to human oversight'],
        'target_time': '120min'
    })
```

## Report Reading Functions for Architect Feedback

```bash
# Function to read wave review issues (R258 compliance)
read_wave_review_issues() {
    local PHASE=$1
    local WAVE=$2
    local REPORT="wave-reviews/phase${PHASE}/wave${WAVE}/PHASE-${PHASE}-WAVE-${WAVE}-REVIEW-REPORT.md"
    
    if [ ! -f "$REPORT" ]; then
        echo "❌ CRITICAL: Wave review report not found!"
        echo "❌ Cannot process CHANGES_REQUIRED without report (R258 violation)"
        exit 1
    fi
    
    echo "📋 Reading architect's required fixes from wave review..."
    echo "📄 Report location: $REPORT"
    echo ""
    
    # Extract decision
    echo "🔍 Decision from report:"
    grep "^\*\*DECISION\*\*:" "$REPORT"
    echo ""
    
    # Extract Issues Identified section
    echo "🔍 Issues Identified:"
    sed -n '/## Issues Identified/,/## Required Actions/p' "$REPORT" | grep -v "^##"
    echo ""
    
    # Extract Required Fixes/Actions
    echo "🔍 Required Actions:"
    sed -n '/## Required Actions/,/## Recommendations/p' "$REPORT" | grep -v "^##"
}

# Function to read phase assessment issues (R257 compliance)
read_phase_assessment_issues() {
    local PHASE=$1
    local REPORT="phase-assessments/phase${PHASE}/PHASE-${PHASE}-ASSESSMENT-REPORT.md"
    
    if [ ! -f "$REPORT" ]; then
        echo "❌ CRITICAL: Phase assessment report not found!"
        echo "❌ Cannot process NEEDS_WORK without report (R257 violation)"
        exit 1
    fi
    
    echo "📋 Reading architect's required fixes from phase assessment..."
    echo "📄 Report location: $REPORT"
    echo ""
    
    # Extract decision
    echo "🔍 Decision from report:"
    grep "^\*\*DECISION\*\*:" "$REPORT"
    echo ""
    
    # Extract overall score
    echo "🔍 Overall Score:"
    grep "^\*\*TOTAL SCORE\*\*" "$REPORT"
    echo ""
    
    # Extract Issues Identified section
    echo "🔍 Issues Identified:"
    sed -n '/## Issues Identified/,/## Required Fixes/p' "$REPORT" | grep -v "^##"
    echo ""
    
    # Extract Required Fixes
    echo "🔍 Required Fixes (Priority 1 - Must Fix):"
    sed -n '/### Priority 1/,/### Priority 2/p' "$REPORT" | grep "^- \["
    echo ""
    
    echo "🔍 Required Fixes (Priority 2 - Should Fix):"
    sed -n '/### Priority 2/,/## Recommendations/p' "$REPORT" | grep "^- \["
}

# Function to parse and assign fixes to agents
parse_and_assign_fixes() {
    local REPORT_TYPE=$1  # WAVE or PHASE
    local REPORT_FILE=$2
    
    echo "📊 Parsing fixes and assigning to agents..."
    
    # Extract critical issues that need sw-engineer
    local CODE_FIXES=$(grep -A2 "\*\*\[CRITICAL\]\*\*\|Assigned To: sw-engineer" "$REPORT_FILE" | grep "Required Fix:" | cut -d: -f2-)
    
    # Extract review issues that need code-reviewer
    local REVIEW_FIXES=$(grep -A2 "Assigned To: code-reviewer" "$REPORT_FILE" | grep "Required Fix:" | cut -d: -f2-)
    
    # Extract architecture issues that need architect consultation
    local ARCH_FIXES=$(grep -A2 "Assigned To: architect" "$REPORT_FILE" | grep "Required Fix:" | cut -d: -f2-)
    
    # Create fix assignment plan
    cat > /tmp/fix-assignments.yaml << EOF
fix_assignments:
  report_type: "$REPORT_TYPE"
  report_file: "$REPORT_FILE"
  assignments:
    sw_engineer:
$(echo "$CODE_FIXES" | sed 's/^/      - /')
    code_reviewer:
$(echo "$REVIEW_FIXES" | sed 's/^/      - /')
    architect:
$(echo "$ARCH_FIXES" | sed 's/^/      - /')
EOF
    
    echo "✅ Fix assignments created: /tmp/fix-assignments.yaml"
}
```

## Recovery Execution Protocol

```yaml
# Example 1: Wave Review Changes Required (R258)
error_recovery_state:
  error_id: "ERR-WAVE-REVIEW-2025-08-27-001"
  detected_at: "2025-08-27T16:15:30Z"
  classification:
    type: "WAVE_REVIEW_CHANGES_REQUIRED"
    severity: "HIGH"
    source_report: "wave-reviews/phase3/wave2/PHASE-3-WAVE-2-REVIEW-REPORT.md"
    architect_decision: "CHANGES_REQUIRED"
    issues_count: 3
    
  recovery_strategy:
    name: "ARCHITECT_DIRECTED_FIXES"
    target_completion: "2025-08-27T18:15:30Z"  # 120min target
    
  architect_requirements:
    - issue_1: 
        description: "API inconsistency in effort2"
        assigned_to: "sw-engineer"
        status: "IN_PROGRESS"
    - issue_2: 
        description: "Missing integration tests"
        assigned_to: "sw-engineer"
        status: "PENDING"
    - issue_3: 
        description: "Size violation in effort3"
        assigned_to: "code-reviewer"
        status: "PENDING"
        
  recovery_steps:
    - step: 1
      action: "READ_WAVE_REVIEW_REPORT"
      status: "COMPLETED"
      completed_at: "2025-08-27T16:16:00Z"
      
    - step: 2
      action: "PARSE_REQUIRED_FIXES"
      status: "COMPLETED"
      completed_at: "2025-08-27T16:17:00Z"
      
    - step: 3
      action: "SPAWN_SW_ENGINEER_FOR_API_FIX"
      status: "IN_PROGRESS"
      started_at: "2025-08-27T16:18:00Z"
      
    - step: 4
      action: "SPAWN_CODE_REVIEWER_FOR_SIZE_FIX"
      status: "PENDING"
      depends_on: "step_3"
      
    - step: 5
      action: "VALIDATE_ALL_FIXES_COMPLETE"
      status: "PENDING"
      depends_on: "all_fixes"
      
    - step: 6
      action: "RETURN_TO_INTEGRATION"
      status: "PENDING"
      depends_on: "step_5"

# Example 2: Phase Assessment Needs Work (R257)
error_recovery_state:
  error_id: "ERR-PHASE-ASSESS-2025-08-27-002"
  detected_at: "2025-08-27T14:30:00Z"
  classification:
    type: "PHASE_ASSESSMENT_NEEDS_WORK"
    severity: "CRITICAL"
    source_report: "phase-assessments/phase2/PHASE-2-ASSESSMENT-REPORT.md"
    architect_decision: "NEEDS_WORK"
    score: 68  # Below passing threshold
    
  recovery_strategy:
    name: "PHASE_LEVEL_REMEDIATION"
    target_completion: "2025-08-27T18:30:00Z"  # 240min target
    
  assessment_requirements:
    priority_1_must_fix:
      - "KCP multi-tenancy patterns missing in wave 2"
      - "API backwards compatibility broken in wave 3"
      - "Integration tests coverage below 60%"
    priority_2_should_fix:
      - "Documentation incomplete for new APIs"
      - "Performance degradation in controller reconciliation"
      
  recovery_steps:
    - step: 1
      action: "READ_PHASE_ASSESSMENT_REPORT"
      status: "COMPLETED"
      completed_at: "2025-08-27T14:31:00Z"
      
    - step: 2
      action: "CREATE_COMPREHENSIVE_FIX_PLAN"
      status: "COMPLETED"
      completed_at: "2025-08-27T14:35:00Z"
      
    - step: 3
      action: "SPAWN_SW_ENGINEER_MULTI_TENANCY_FIX"
      status: "IN_PROGRESS"
      target_wave: 2
      
    - step: 4
      action: "SPAWN_SW_ENGINEER_API_COMPATIBILITY"
      status: "PENDING"
      target_wave: 3
      
    - step: 5
      action: "SPAWN_CODE_REVIEWER_TEST_COVERAGE"
      status: "PENDING"
      target_waves: [1, 2, 3]
      
    - step: 6
      action: "VALIDATE_ALL_PRIORITY_1_FIXED"
      status: "PENDING"
      depends_on: ["step_3", "step_4", "step_5"]
      
    - step: 7
      action: "RETURN_TO_PHASE_INTEGRATION"
      status: "PENDING"
      depends_on: "step_6"
      note: "Per R259 - Must create phase integration branch before reassessment"

# Example 3: Standard Size Violation
error_recovery_state:
  error_id: "ERR-SIZE-2025-08-27-003"
  detected_at: "2025-08-27T12:00:00Z"
  classification:
    type: "SIZE_LIMIT_VIOLATION"
    severity: "CRITICAL"
    affected_efforts: ["effort2-controller"]
    measured_size: 1250
    limit: 800
    
  recovery_strategy:
    name: "IMMEDIATE_SPLIT"
    target_completion: "2025-08-27T12:30:00Z"  # 30min target
    
  recovery_steps:
    - step: 1
      action: "STOP_ALL_WORK"
      status: "COMPLETED"
      
    - step: 2
      action: "SPAWN_CODE_REVIEWER_SPLIT_ANALYSIS"
      status: "IN_PROGRESS"
      
    - step: 3
      action: "CREATE_SPLIT_PLAN"
      status: "PENDING"
      
    - step: 4
      action: "EXECUTE_SPLITS_SEQUENTIALLY"
      status: "PENDING"
      
    - step: 5
      action: "VALIDATE_ALL_SPLITS_UNDER_800"
      status: "PENDING"
```

## Critical Recovery Actions

### ⚠️⚠️⚠️ RULE R010 - Wrong Location Handling
**SEE**: `$CLAUDE_PROJECT_DIR/rule-library/R010-wrong-location-handling.md`
