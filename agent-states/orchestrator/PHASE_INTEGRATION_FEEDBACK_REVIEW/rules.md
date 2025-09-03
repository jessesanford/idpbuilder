# Orchestrator - PHASE_INTEGRATION_FEEDBACK_REVIEW State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED PHASE_INTEGRATION_FEEDBACK_REVIEW STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_PHASE_INTEGRATION_FEEDBACK_REVIEW
echo "$(date +%s) - Rules read and acknowledged for PHASE_INTEGRATION_FEEDBACK_REVIEW" > .state_rules_read_orchestrator_PHASE_INTEGRATION_FEEDBACK_REVIEW
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY PHASE FEEDBACK WORK UNTIL RULES ARE READ:
- ❌ Parse phase integration reports
- ❌ Identify failed waves
- ❌ Extract conflict information
- ❌ Create fix request metadata
- ❌ Update state files
- ❌ Continue to next state
- ❌ Think about what to do in this state

### ✅ YOU MUST IMMEDIATELY:

## 🔴🔴🔴 MANDATORY STATE RULE READING AND ACKNOWLEDGMENT 🔴🔴🔴

### ⚠️⚠️⚠️ YOU MUST READ EACH RULE FILE LISTED IN PRIMARY DIRECTIVES. **I AM WATCHING YOUR TOOL CALLS FOR READ OPERATIONS** *YOU WILL FAIL* IF YOU DO NOT MAKE A READ FILE CALL FOR EACH RULE FILE IN PRIMARY DIRECTIVES!!! ⚠️⚠️⚠️

**AFTER READING, YOU MUST ACKNOWLEDGE ALL THE STATE RULES AND STATE THAT YOU WILL ABIDE BY THEM ONE AT A TIME GIVING THE RULE NUMBER AND DESCRIPTION.**

### PRIMARY DIRECTIVES - MANDATORY READING:

**USE THESE EXACT READ COMMANDS (IN THIS ORDER):**
1. Read: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. Read: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md
3. Read: $CLAUDE_PROJECT_DIR/rule-library/R290-state-rule-reading-verification-supreme-law.md
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R282-phase-integration-protocol.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R238-integration-report-evaluation.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R239-fix-plan-distribution.md
8. Read: $CLAUDE_PROJECT_DIR/rule-library/R206-state-machine-transition-validation.md

**WE ARE WATCHING EACH READ TOOL CALL**

### ❌ ANTI-PATTERNS THAT WILL CAUSE FAILURE:

1. **Fake Acknowledgment Without Reading**:
   ```
   ❌ WRONG: "I acknowledge R234, R208, R290..."
   (No Read tool calls detected = AUTOMATIC FAILURE)
   ```

2. **Bulk Acknowledgment**:
   ```
   ❌ WRONG: "I acknowledge all PHASE_INTEGRATION_FEEDBACK_REVIEW rules"
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
   ❌ WRONG: "I know R282 requires phase integration..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR PHASE_INTEGRATION_FEEDBACK_REVIEW:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute PHASE_INTEGRATION_FEEDBACK_REVIEW work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY phase feedback work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been created
4. ✅ You have stated readiness to execute PHASE_INTEGRATION_FEEDBACK_REVIEW work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY phase feedback work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

---

## 🔴🔴🔴 SUPREME DIRECTIVE: ANALYZE PHASE INTEGRATION FAILURES 🔴🔴🔴

**PARSE PHASE-LEVEL INTEGRATION FAILURES AND INITIATE FIXES!**

## State Overview

In PHASE_INTEGRATION_FEEDBACK_REVIEW, you analyze phase integration report to identify which waves/efforts failed at the phase level.

## Required Actions

### 1. Parse Phase Integration Report
```bash
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
REPORT_FILE="efforts/phase${PHASE}/phase-integration/PHASE_INTEGRATION_REPORT.md"

echo "📋 Parsing phase integration report: $REPORT_FILE"

# Extract failed waves and conflicts
FAILED_WAVES=()
CONFLICT_FILES=()
MISSING_DEPS=()

# Parse Failed Waves section
while IFS= read -r line; do
    if [[ "$line" =~ ^-[[:space:]]+wave([0-9]+):[[:space:]]+(.+) ]]; then
        WAVE_NUM="${BASH_REMATCH[1]}"
        ISSUE="${BASH_REMATCH[2]}"
        FAILED_WAVES+=("wave${WAVE_NUM}:${ISSUE}")
    fi
done < <(sed -n '/## Failed Waves/,/## /p' "$REPORT_FILE" | grep "^-")

# Parse Merge Conflicts section
while IFS= read -r line; do
    if [[ "$line" =~ ^-[[:space:]]+(.+) ]]; then
        CONFLICT_FILES+=("${BASH_REMATCH[1]}")
    fi
done < <(sed -n '/## Merge Conflicts/,/## /p' "$REPORT_FILE" | grep "^-")

# Parse Missing Dependencies
if grep -q "Missing Dependencies:" "$REPORT_FILE"; then
    while IFS= read -r dep; do
        MISSING_DEPS+=("$dep")
    done < <(sed -n '/## Missing Dependencies/,/## /p' "$REPORT_FILE" | grep "^-" | sed 's/^- //')
fi

echo "Failed waves: ${#FAILED_WAVES[@]}"
echo "Conflict files: ${#CONFLICT_FILES[@]}"
echo "Missing dependencies: ${#MISSING_DEPS[@]}"
```

### 2. Create Phase Fix Request
```bash
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
PHASE_FIX_REQUEST="efforts/phase${PHASE}/PHASE_FIX_REQUEST_${TIMESTAMP}.yaml"

cat > "$PHASE_FIX_REQUEST" << EOF
phase_fix_request:
  timestamp: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  phase: $PHASE
  integration_report: "$REPORT_FILE"
  
  failed_waves:
EOF

for wave_issue in "${FAILED_WAVES[@]}"; do
    WAVE="${wave_issue%%:*}"
    ISSUE="${wave_issue#*:}"
    cat >> "$PHASE_FIX_REQUEST" << EOF
    - wave: "$WAVE"
      issue: "$ISSUE"
      directory: "efforts/phase${PHASE}/${WAVE}"
EOF
done

if [ ${#CONFLICT_FILES[@]} -gt 0 ]; then
    cat >> "$PHASE_FIX_REQUEST" << EOF
  
  merge_conflicts:
EOF
    for file in "${CONFLICT_FILES[@]}"; do
        echo "    - \"$file\"" >> "$PHASE_FIX_REQUEST"
    done
fi

if [ ${#MISSING_DEPS[@]} -gt 0 ]; then
    cat >> "$PHASE_FIX_REQUEST" << EOF
  
  missing_dependencies:
EOF
    for dep in "${MISSING_DEPS[@]}"; do
        echo "    - \"$dep\"" >> "$PHASE_FIX_REQUEST"
    done
fi

cat >> "$PHASE_FIX_REQUEST" << EOF
  
  fix_strategy:
    - resolve_conflicts: ${#CONFLICT_FILES[@]} files
    - install_dependencies: ${#MISSING_DEPS[@]} packages
    - fix_failed_waves: ${#FAILED_WAVES[@]} waves
EOF

echo "✅ Created phase fix request: $PHASE_FIX_REQUEST"
```

### 3. Update State and Transition
```bash
# Record in state file
yq eval ".phase_integration_feedback.phase${PHASE}.report_parsed = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.yaml
yq eval ".phase_integration_feedback.phase${PHASE}.fix_request_file = \"$PHASE_FIX_REQUEST\"" -i orchestrator-state.yaml
yq eval ".phase_integration_feedback.phase${PHASE}.failed_waves = ${#FAILED_WAVES[@]}" -i orchestrator-state.yaml
yq eval ".phase_integration_feedback.phase${PHASE}.conflicts = ${#CONFLICT_FILES[@]}" -i orchestrator-state.yaml

# Determine next state
if [ ${#FAILED_WAVES[@]} -gt 0 ] || [ ${#CONFLICT_FILES[@]} -gt 0 ] || [ ${#MISSING_DEPS[@]} -gt 0 ]; then
    echo "🔧 Phase integration issues found - spawning Code Reviewer for fix planning"
    UPDATE_STATE="SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN"
    UPDATE_REASON="Phase integration issues require fix plans"
else
    echo "✅ No phase integration issues - proceeding to next phase"
    UPDATE_STATE="PHASE_COMPLETE"
    UPDATE_REASON="Phase integration successful"
fi

# Update state
yq eval ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.yaml
yq eval ".state_transition_history += [{\"from\": \"PHASE_INTEGRATION_FEEDBACK_REVIEW\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"$UPDATE_REASON\"}]" -i orchestrator-state.yaml

# Commit
git add orchestrator-state.yaml "$PHASE_FIX_REQUEST"
git commit -m "state: Phase feedback review complete - $UPDATE_REASON"
git push
```

## Valid Transitions

1. **FIX NEEDED Path**: `PHASE_INTEGRATION_FEEDBACK_REVIEW` → `SPAWN_CODE_REVIEWER_PHASE_FIX_PLAN`
   - When: Phase integration has failures requiring fixes
   
2. **SUCCESS Path**: `PHASE_INTEGRATION_FEEDBACK_REVIEW` → `PHASE_COMPLETE`
   - When: Phase integration is successful with no issues

## Grading Criteria

- ✅ **+25%**: Parse phase integration report correctly
- ✅ **+25%**: Create comprehensive fix request
- ✅ **+25%**: Update state file properly
- ✅ **+25%**: Transition to correct next state

## Common Violations

- ❌ **-100%**: Not parsing integration report
- ❌ **-50%**: Missing fix request creation
- ❌ **-50%**: Wrong state transition
- ❌ **-30%**: Not recording metadata

## Related Rules

- R282: Phase Integration Protocol
- R238: Integration Report Evaluation Protocol
- R239: Fix Plan Distribution Protocol
- R206: State Machine Transition Validation