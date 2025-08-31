# Orchestrator - SPAWN_ENGINEERS_FOR_FIXES State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED SPAWN_ENGINEERS_FOR_FIXES STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_SPAWN_ENGINEERS_FOR_FIXES
echo "$(date +%s) - Rules read and acknowledged for SPAWN_ENGINEERS_FOR_FIXES" > .state_rules_read_orchestrator_SPAWN_ENGINEERS_FOR_FIXES
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY SPAWNING WORK UNTIL RULES ARE READ:
- ❌ Identify efforts needing engineers
- ❌ Create command files
- ❌ Spawn Software Engineer agents  
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
5. Read: $CLAUDE_PROJECT_DIR/rule-library/R239-fix-plan-distribution.md
6. Read: $CLAUDE_PROJECT_DIR/rule-library/R197-one-agent-per-effort.md
7. Read: $CLAUDE_PROJECT_DIR/rule-library/R209-effort-directory-isolation-protocol.md
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
   ❌ WRONG: "I acknowledge all SPAWN_ENGINEERS_FOR_FIXES rules"
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
   ❌ WRONG: "I know R239 requires fix plan distribution..."
   (Must READ from file, not recall from memory)
   ```

5. **Skipping Rules in PRIMARY DIRECTIVES**:
   ```
   ❌ WRONG: Reading only some rules from the list
   (ALL rules in PRIMARY DIRECTIVES are MANDATORY)
   ```

### ✅ CORRECT PATTERN FOR SPAWN_ENGINEERS_FOR_FIXES:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md
2. "I acknowledge R234 - Mandatory State Traversal: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/R006-orchestrator-never-writes-code.md  
4. "I acknowledge R006 - Orchestrator Never Writes Code: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. Create verification marker
6. "Ready to execute SPAWN_ENGINEERS_FOR_FIXES work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY spawning work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ Verification marker has been created
4. ✅ You have stated readiness to execute SPAWN_ENGINEERS_FOR_FIXES work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY spawning work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

---

## 🔴🔴🔴 SUPREME DIRECTIVE: SPAWN ENGINEERS TO EXECUTE FIXES 🔴🔴🔴

**SPAWN SW ENGINEERS TO IMPLEMENT INTEGRATION FIXES!**

## State Overview

In SPAWN_ENGINEERS_FOR_FIXES, you spawn Software Engineer agents to implement the fixes specified in the distributed fix plans.

## Required Actions

### 1. Identify Efforts Needing Engineers
```bash
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

echo "🔧 Identifying efforts that need fix engineers"

# Get list of efforts with distributed fix plans
EFFORTS_TO_FIX=()
while IFS= read -r effort; do
    if [ "$effort" != "null" ]; then
        NEEDS_FIXES=$(yq ".efforts_in_progress.\"$effort\".needs_fixes" orchestrator-state.yaml)
        if [ "$NEEDS_FIXES" = "true" ]; then
            EFFORTS_TO_FIX+=("$effort")
            echo "  - $effort needs fixes"
        fi
    fi
done < <(yq '.efforts_in_progress | keys | .[]' orchestrator-state.yaml)

echo "Total efforts needing fixes: ${#EFFORTS_TO_FIX[@]}"
```

### 2. Spawn Engineers for Each Effort
```bash
SPAWN_LOG="efforts/phase${PHASE}/wave${WAVE}/FIX_SPAWN_LOG_${TIMESTAMP}.md"
echo "# Fix Engineer Spawn Log" > "$SPAWN_LOG"
echo "Timestamp: $(date -u +%Y-%m-%dT%H:%M:%SZ)" >> "$SPAWN_LOG"
echo "" >> "$SPAWN_LOG"

for effort in "${EFFORTS_TO_FIX[@]}"; do
    EFFORT_DIR="efforts/phase${PHASE}/wave${WAVE}/${effort}"
    COMMAND_FILE="${EFFORT_DIR}/sw-engineer-fix-command.md"
    
    echo "🚀 Spawning SW Engineer for: $effort"
    
    # Create command file
    cat > "$COMMAND_FILE" << EOF
# SOFTWARE ENGINEER FIX IMPLEMENTATION TASK

## Your State
You are in state: FIX_INTEGRATION_ISSUES

## Critical Information
- **Working Directory**: $EFFORT_DIR
- **Branch**: phase${PHASE}-wave${WAVE}-${effort}
- **Fix Plan Location**: Check FIX_PLAN_LOCATION.txt for the fix plan file

## Required Actions

1. **Read the fix plan**:
   - The file path is in FIX_PLAN_LOCATION.txt
   - Follow ALL instructions in the fix plan

2. **Implement fixes**:
   - Apply all code changes specified
   - Install any missing dependencies
   - Resolve conflicts if any
   - Fix build issues

3. **Verify fixes**:
   - Run all verification commands from the fix plan
   - Ensure build passes
   - Run tests to confirm fixes work

4. **Update status**:
   - Remove FIX_REQUIRED.flag when complete
   - Create FIX_COMPLETE.flag with summary
   - Commit all changes with clear message

5. **Important**:
   - Do NOT create new features
   - ONLY fix the specified issues
   - Follow the fix plan exactly

## Success Criteria
- All issues from fix plan resolved
- Build passes successfully
- Tests pass (if applicable)
- FIX_COMPLETE.flag created
EOF

    # Spawn the engineer
    echo "@agent-sw-engineer Please execute the fix task in: $COMMAND_FILE"
    echo "  Working directory: $EFFORT_DIR"
    
    # Record spawn
    echo "- Spawned engineer for: $effort" >> "$SPAWN_LOG"
    echo "  Command: $COMMAND_FILE" >> "$SPAWN_LOG"
    
    # Update state file
    yq eval ".agents_spawned += [{\"type\": \"sw-engineer\", \"task\": \"fix_integration\", \"effort\": \"$effort\", \"state\": \"FIX_INTEGRATION_ISSUES\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"command_file\": \"$COMMAND_FILE\"}]" -i orchestrator-state.yaml
    yq eval ".efforts_in_progress.\"$effort\".fix_engineer_spawned = \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"" -i orchestrator-state.yaml
done

echo "" >> "$SPAWN_LOG"
echo "Total engineers spawned: ${#EFFORTS_TO_FIX[@]}" >> "$SPAWN_LOG"
```

### 3. Transition to Monitoring
```bash
# Update state
yq eval ".current_state = \"MONITORING_FIX_PROGRESS\"" -i orchestrator-state.yaml
yq eval ".integration_feedback.wave${WAVE}.fix_engineers_spawned = ${#EFFORTS_TO_FIX[@]}" -i orchestrator-state.yaml
yq eval ".state_transition_history += [{\"from\": \"SPAWN_ENGINEERS_FOR_FIXES\", \"to\": \"MONITORING_FIX_PROGRESS\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\", \"reason\": \"Spawned ${#EFFORTS_TO_FIX[@]} engineers for fixes\"}]" -i orchestrator-state.yaml

# Commit
git add orchestrator-state.yaml "$SPAWN_LOG"
for effort in "${EFFORTS_TO_FIX[@]}"; do
    git add "efforts/phase${PHASE}/wave${WAVE}/${effort}/sw-engineer-fix-command.md"
done
git commit -m "spawn: ${#EFFORTS_TO_FIX[@]} SW Engineers for integration fixes"
git push
```

## Valid Transitions

1. **ALWAYS**: `SPAWN_ENGINEERS_FOR_FIXES` → `MONITORING_FIX_PROGRESS`
   - Transition after all engineers spawned

## Spawn Requirements

For each engineer spawned:
1. Create clear command file with fix instructions
2. Reference the distributed fix plan
3. Specify working directory and branch
4. Set state to FIX_INTEGRATION_ISSUES
5. Record spawn in state file

## Grading Criteria

- ✅ **+25%**: Spawn engineer for each effort
- ✅ **+25%**: Create proper command files
- ✅ **+25%**: Reference fix plans correctly
- ✅ **+25%**: Update state file properly

## Common Violations

- ❌ **-100%**: Not spawning any engineers
- ❌ **-50%**: Missing fix plan references
- ❌ **-50%**: Wrong working directories
- ❌ **-30%**: Not recording spawns

## Related Rules

- R239: Fix Plan Distribution Protocol
- R197: One Agent Per Effort
- R209: Effort Directory Isolation
- R206: State Machine Transition Validation