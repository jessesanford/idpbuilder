# Orchestrator - INJECT_WAVE_METADATA State Rules

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED INJECT_WAVE_METADATA STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_INJECT_WAVE_METADATA
echo "$(date +%s) - Rules read and acknowledged for INJECT_WAVE_METADATA" > .state_rules_read_orchestrator_INJECT_WAVE_METADATA
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY INJECT_WAVE_METADATA WORK UNTIL RULES ARE READ:
- ❌ Start inject wave metadata
- ❌ Start update tracking files
- ❌ Start configure wave settings
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
   ❌ WRONG: "I acknowledge all INJECT_WAVE_METADATA rules"
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

### ✅ CORRECT PATTERN FOR INJECT_WAVE_METADATA:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute INJECT_WAVE_METADATA work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY INJECT_WAVE_METADATA work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute INJECT_WAVE_METADATA work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with INJECT_WAVE_METADATA work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY INJECT_WAVE_METADATA work before reading and acknowledging rules:**
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

## 📋 PRIMARY DIRECTIVES FOR INJECT_WAVE_METADATA STATE

### 🚨🚨🚨 R213 - Wave and Effort Metadata Injection
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R213-wave-and-effort-metadata-protocol.md`
**Criticality**: BLOCKING - Must inject metadata before spawning
**Summary**: Inject parallelization metadata into wave implementation plans

### 🔴🔴🔴 R234 - Mandatory State Traversal (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R234-mandatory-state-traversal-supreme-law.md`
**Criticality**: SUPREME LAW - Violation = -100% automatic failure
**Summary**: Must traverse all states in sequence, no skipping allowed

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml on all state changes

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push within 60 seconds
**Summary**: Commit and push state file immediately after updates

### 🔴🔴🔴 R232 - TodoWrite Pending Items Override (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R232-todowrite-pending-items-override.md`
**Criticality**: SUPREME LAW - Pending items are COMMANDS
**Summary**: Any pending TODO items must be executed immediately

## 🚨 INJECT_WAVE_METADATA IS A VERB - START INJECTING METADATA IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING INJECT_WAVE_METADATA

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Open wave implementation plan for editing NOW
2. Insert R213 parallelization metadata immediately
3. Check TodoWrite for pending items and process them
4. Save and validate metadata without delay

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in INJECT_WAVE_METADATA" [stops]
- ❌ "Successfully entered INJECT_WAVE_METADATA state" [waits]
- ❌ "Ready to start injecting metadata" [pauses]
- ❌ "I'm in INJECT_WAVE_METADATA state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering INJECT_WAVE_METADATA, opening wave implementation plan for editing NOW..."
- ✅ "START INJECTING METADATA per R213, inserting parallelization metadata immediately..."
- ✅ "INJECT_WAVE_METADATA: Saving and validating metadata without delay..."

## State Context
This state is responsible for injecting R213 parallelization metadata into wave implementation plans before any agent spawning occurs.

## R213 Metadata Injection Protocol

**CRITICAL**: This metadata MUST be injected BEFORE spawning Code Reviewers or SW Engineers!

```bash
# Example metadata to inject into wave plan:
EFFORT_METADATA:
  effort_id: "E3.1.1"
  name: "sync-engine-foundation"
  can_parallelize: false
  blocks: ["E3.1.2", "E3.1.3", "E3.1.4", "E3.1.5"]
  dependencies: []
  estimated_lines: 600
  assigned_to: "sw-engineer-1"
  working_directory: "/efforts/phase3/wave1/sync-engine-foundation"
  branch: "phase3/wave1/sync-engine-foundation"
```

## Critical Requirements

1. **READ** the wave implementation plan with Read tool
2. **IDENTIFY** all efforts in the wave
3. **INJECT** parallelization metadata for EACH effort
4. **SAVE** the updated plan
5. **VERIFY** metadata is present before proceeding
6. **TRANSITION** to next state in mandatory sequence

## State Transitions

- **FROM**: Previous state in mandatory sequence
- **TO**: Next state per R234 mandatory traversal
- **CANNOT SKIP**: This state is part of mandatory sequence

## Validation Before Transition

```bash
validate_metadata_injection() {
    echo "🔍 Validating R213 metadata injection..."
    
    # Check each effort has metadata
    for effort in efforts/phase${PHASE}/wave${WAVE}/*/; do
        if ! grep -q "can_parallelize:" "${effort}/IMPLEMENTATION-PLAN.md"; then
            echo "❌ FATAL: Missing parallelization metadata in ${effort}"
            exit 213
        fi
    done
    
    echo "✅ All efforts have R213 metadata"
}
```

## Next Steps

After successfully injecting metadata:
1. Update orchestrator-state.yaml (R288)
2. Commit and push changes (R288)
3. Check TodoWrite for pending items (R232)
4. Transition to next mandatory state (R234)
