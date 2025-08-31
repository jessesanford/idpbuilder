# Orchestrator - INTEGRATION_FEEDBACK_REVIEW State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED INTEGRATION_FEEDBACK_REVIEW STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INTEGRATION_FEEDBACK_REVIEW
echo "$(date +%s) - Rules read and acknowledged for INTEGRATION_FEEDBACK_REVIEW" > .state_rules_read_orchestrator_INTEGRATION_FEEDBACK_REVIEW
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INTEGRATION FEEDBACK WORK UNTIL RULES ARE READ:
- ❌ Parse integration reports
- ❌ Identify failed efforts
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
4. Read: $CLAUDE_PROJECT_DIR/rule-library/R291-integration-demo-requirement.md
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R292-integration-fixes-in-effort-branches.md
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
   ❌ WRONG: "I acknowledge all INTEGRATION_FEEDBACK_REVIEW rules"
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
   ❌ WRONG: "I know R238 requires parsing integration reports..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR INTEGRATION_FEEDBACK_REVIEW:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute INTEGRATION_FEEDBACK_REVIEW work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY INTEGRATION_FEEDBACK_REVIEW work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been CREATED
4. ✅ You have stated readiness to execute INTEGRATION_FEEDBACK_REVIEW work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY INTEGRATION_FEEDBACK_REVIEW work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**After reading ALL rules, acknowledge them:**
□ I have read R234 - Mandatory State Traversal (SUPREME LAW #1)
□ I have read R208 - No Implementation in Orchestrator (SUPREME LAW #2)
□ I have read R290 - State Rule Reading and Verification (SUPREME LAW #3)
□ I have read R238 - Integration Report Evaluation Protocol
□ I have read R239 - Fix Plan Distribution Protocol
□ I have read R206 - State Machine Transition Validation

**CRITICAL**: You must have made 7 actual Read tool calls. Count them!

---

## 🔴🔴🔴 SUPREME DIRECTIVE: ANALYZE AND FIX INTEGRATION FAILURES 🔴🔴🔴

**YOU MUST PARSE INTEGRATION FAILURES AND INITIATE FIX CYCLES!**

## State Overview

In INTEGRATION_FEEDBACK_REVIEW, you analyze the integration report to identify which efforts failed and what fixes are needed.

## Required Actions

### 1. Parse Integration Report for Failed Efforts
```bash
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
REPORT_FILE="efforts/phase${PHASE}/wave${WAVE}/integration-workspace/INTEGRATION_REPORT.md"

echo "📋 Parsing integration report: $REPORT_FILE"

# Extract failed branches and issues
FAILED_BRANCHES=()
ISSUES_BY_BRANCH=()

# Parse the Failed Branches section
while IFS= read -r line; do
    if [[ "$line" =~ ^-[[:space:]]+(.+):[[:space:]]+(.+) ]]; then
        BRANCH="${BASH_REMATCH[1]}"
        ISSUE="${BASH_REMATCH[2]}"
        FAILED_BRANCHES+=("$BRANCH")
        ISSUES_BY_BRANCH+=("$BRANCH:$ISSUE")
    fi
done < <(sed -n '/## Failed Branches/,/## /p' "$REPORT_FILE" | grep "^-")

echo "Found ${#FAILED_BRANCHES[@]} failed branches"
```

### 2. Identify Efforts Needing Fixes
```bash
# Map branches to efforts
EFFORTS_NEEDING_FIXES=()
for branch in "${FAILED_BRANCHES[@]}"; do
    # Extract effort name from branch (e.g., wave1-effort1 from phase1-wave1-effort1)
    EFFORT=$(echo "$branch" | sed "s/phase${PHASE}-wave${WAVE}-//")
    EFFORTS_NEEDING_FIXES+=("$EFFORT")
    
    # Record in state file
    yq eval ".integration_feedback.wave${WAVE}.efforts_needing_fixes += [\"$EFFORT\"]" -i orchestrator-state.yaml
done

# Check if we have build dependency issues
if grep -q "BLOCKED_BY_DEPENDENCIES" "$REPORT_FILE"; then
    echo "🚨 Build blocked by missing dependencies - need system-level fixes"
    
    # Extract missing dependencies
    MISSING_DEPS=$(grep -A5 "Missing Dependencies:" "$REPORT_FILE" | grep "^-" | sed 's/^- //')
    
    # Record in state file
    yq eval ".integration_feedback.wave${WAVE}.missing_dependencies = \"$MISSING_DEPS\"" -i orchestrator-state.yaml
    yq eval ".integration_feedback.wave${WAVE}.requires_system_fixes = true" -i orchestrator-state.yaml
fi
```

### 3. Create Fix Request Metadata
```bash
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
FIX_REQUEST_FILE="efforts/phase${PHASE}/wave${WAVE}/FIX_REQUEST_${TIMESTAMP}.yaml"

cat > "$FIX_REQUEST_FILE" << EOF
fix_request:
  timestamp: "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  phase: $PHASE
  wave: $WAVE
  integration_report: "$REPORT_FILE"
  failed_efforts:
EOF

for i in "${!EFFORTS_NEEDING_FIXES[@]}"; do
    EFFORT="${EFFORTS_NEEDING_FIXES[$i]}"
    ISSUE="${ISSUES_BY_BRANCH[$i]#*:}"
    cat >> "$FIX_REQUEST_FILE" << EOF
    - effort: "$EFFORT"
      issue: "$ISSUE"
      branch: "phase${PHASE}-wave${WAVE}-${EFFORT}"
      working_directory: "efforts/phase${PHASE}/wave${WAVE}/${EFFORT}"
EOF
done

if [ "${#EFFORTS_NEEDING_FIXES[@]}" -eq 0 ]; then
    echo "⚠️ No specific efforts identified but integration failed"
    echo "Possible system-level issue requiring investigation"
    UPDATE_STATE="ERROR_RECOVERY"
else
    echo "✅ Fix request metadata created"
    UPDATE_STATE="SPAWN_CODE_REVIEWER_FIX_PLAN"
fi
```

### 4. Update State File
```bash
# Update orchestrator state
yq eval ".current_state = \"$UPDATE_STATE\"" -i orchestrator-state.yaml
yq eval ".integration_feedback.wave${WAVE}.fix_request_file = \"$FIX_REQUEST_FILE\"" -i orchestrator-state.yaml
yq eval ".integration_feedback.wave${WAVE}.review_timestamp = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.yaml
yq eval ".state_transition_history += [{\"from\": \"INTEGRATION_FEEDBACK_REVIEW\", \"to\": \"$UPDATE_STATE\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"Identified ${#EFFORTS_NEEDING_FIXES[@]} efforts needing fixes\"}]" -i orchestrator-state.yaml

# Commit changes
git add orchestrator-state.yaml "$FIX_REQUEST_FILE"
git commit -m "feedback: Identified integration issues for wave ${WAVE} - ${#EFFORTS_NEEDING_FIXES[@]} efforts need fixes"
git push
```

## Valid Transitions

1. **FIX PATH**: `INTEGRATION_FEEDBACK_REVIEW` → `SPAWN_CODE_REVIEWER_FIX_PLAN`
   - When: Failed efforts identified, fix plans needed
   
2. **ERROR PATH**: `INTEGRATION_FEEDBACK_REVIEW` → `ERROR_RECOVERY`
   - When: No specific efforts identified or system-level issues

## Grading Criteria

- ✅ **+25%**: Correctly parse integration report
- ✅ **+25%**: Identify all failed efforts
- ✅ **+25%**: Create fix request metadata
- ✅ **+25%**: Update state file with feedback details

## Common Violations

- ❌ **-100%**: Not reading integration report
- ❌ **-50%**: Missing failed efforts
- ❌ **-30%**: Not creating fix request metadata
- ❌ **-30%**: Not recording in state file

## Integration Report Format Expected

```markdown
## Integration Report
Integration Status: FAILED
Build Status: BLOCKED_BY_DEPENDENCIES
Test Status: NOT_RUN

## Failed Branches
- phase1-wave1-effort1: Missing dependency libgpgme
- phase1-wave1-effort2: Build error in authentication module
- phase1-wave1-effort3: Merge conflict in shared config

## Missing Dependencies
- libgpgme-dev
- libbtrfs-dev
```

## Related Rules

- R238: Integration Report Evaluation Protocol
- R239: Fix Plan Distribution Protocol
- R260: Integration Agent Core Requirements
- R206: State Machine Transition Validation