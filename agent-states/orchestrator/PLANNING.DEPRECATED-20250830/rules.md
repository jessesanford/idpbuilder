# Orchestrator - PLANNING State Rules

## 🔴🔴🔴 R290 ENFORCEMENT: READ THESE RULES FIRST! 🔴🔴🔴

**SUPREME LAW #3 (R290): STATE RULES MUST BE READ BEFORE STATE ACTIONS**

## 🔴🔴🔴 STOP! STATE RULE READING IS ABSOLUTELY FIRST! 🔴🔴🔴

**YOU HAVE ENTERED PLANNING STATE - YOU MUST READ AND ACKNOWLEDGE ALL STATE RULES BEFORE DOING ANY STATE WORK!**

## 🔴🔴🔴 R290 VERIFICATION REQUIREMENT 🔴🔴🔴

**R290 ENFORCEMENT: CREATE VERIFICATION MARKER AFTER READING**

After reading and acknowledging all state rules, you MUST create a verification marker:

```bash
# MANDATORY: Create verification marker after reading rules
touch .state_rules_read_orchestrator_PLANNING
echo "$(date +%s) - Rules read and acknowledged for PLANNING" > .state_rules_read_orchestrator_PLANNING
```

**FAILURE TO CREATE MARKER = AUTOMATIC -100% PENALTY**

The system will check for this marker. No marker = Immediate failure.

### ❌ DO NOT DO ANY PLANNING WORK UNTIL RULES ARE READ:
- ❌ Start load planning templates
- ❌ Start spawn architects
- ❌ Start request implementation plans
- ❌ Start create phase plans
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
   ❌ WRONG: "I acknowledge all PLANNING rules"
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

### ✅ CORRECT PATTERN FOR PLANNING:
```
1. READ: $CLAUDE_PROJECT_DIR/rule-library/[first-rule-file].md
2. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
3. READ: $CLAUDE_PROJECT_DIR/rule-library/[second-rule-file].md  
4. "I acknowledge [Rule#] - [Rule Name]: [Brief description]"
[Continue for EVERY rule in PRIMARY DIRECTIVES...]
5. "Ready to execute PLANNING work"
```

### 🚨 NO WORK UNTIL ACKNOWLEDGMENT COMPLETE 🚨
**You may NOT begin ANY PLANNING work until:**
1. ✅ ALL rules in PRIMARY DIRECTIVES have been READ
2. ✅ ALL rules have been individually ACKNOWLEDGED
3. ✅ You have stated readiness to execute PLANNING work
1. **READ** every rule file listed in PRIMARY DIRECTIVES below
2. **ACKNOWLEDGE** each rule individually with number and description
3. **ONLY THEN** proceed with PLANNING work

### 🚨 FAILURE TO READ STATE RULES FIRST = IMMEDIATE EXIT 🚨
**If you do ANY PLANNING work before reading and acknowledging rules:**
- **STOP ALL WORK IMMEDIATELY**
- **EXIT WITH FAILURE STATUS**
- **YOU HAVE VIOLATED STATE COMPLIANCE**

**THE SYSTEM IS MONITORING YOUR READ TOOL CALLS!**

## 📋 PRIMARY DIRECTIVES FOR PLANNING STATE

### 🚨🚨🚨 R109 - Planning Rules
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R109-planning-rules.md`
**Criticality**: BLOCKING - Must follow planning templates
**Summary**: Use correct templates for master, phase, and effort planning

### 🚨🚨🚨 R287 - TODO Persistence Suite
**Files**:
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-triggers.md`
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-save-frequency.md`
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-commit-protocol.md`
- R287: `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-recovery-verification.md`
**Criticality**: BLOCKING - TODO loss = -50% to -100% penalty
**Summary**: Save TODOs within 30s, every 10 messages/15 min, commit within 60s

### 🔴🔴🔴 R288 - State File Update and Commit Protocol (SUPREME LAW)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: SUPREME LAW - Update on every transition
**Summary**: Update orchestrator-state.yaml on all state changes

### 🚨🚨🚨 R288 - State File Update and Commit Protocol
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R288-state-file-update-and-commit-protocol.md`
**Criticality**: BLOCKING - Push within 60 seconds
**Summary**: Commit and push state file immediately after updates

## 🚨 PLANNING IS A VERB - START PLANNING IMMEDIATELY! 🚨

### IMMEDIATE ACTIONS UPON ENTERING PLANNING

**THE MOMENT YOU ENTER THIS STATE, YOU MUST:**
1. Determine which level of planning is needed NOW (master/phase/wave)
2. Copy the appropriate template file NOW
3. Start filling in the plan sections immediately
4. Check TodoWrite for any pending planning tasks
5. Begin populating placeholders without delay

**FORBIDDEN - AUTOMATIC FAILURE:**
- ❌ "STATE TRANSITION COMPLETE: Now in PLANNING" [stops]
- ❌ "Successfully entered PLANNING state" [waits]
- ❌ "Ready to create plans" [pauses]
- ❌ "I'm in PLANNING state" [does nothing]

**REQUIRED - IMMEDIATE ACTION:**
- ✅ "Entering PLANNING, copying master template now..."
- ✅ "Planning Phase 2, loading phase template..."
- ✅ "PLANNING: Creating wave structure for Phase ${PHASE}..."

## State Context
You are creating or reviewing implementation plans at various levels (master, phase, wave) - START IMMEDIATELY!

## Planning Protocol

### Template Locations:
- **Master Plan**: `templates/MASTER-IMPLEMENTATION-PLAN.md`
- **Phase Plan**: `templates/PHASE-IMPLEMENTATION-PLAN.md`  
- **Effort Plan**: `templates/EFFORT-PLANNING-TEMPLATE.md`
- **Work Log**: `templates/WORK-LOG-TEMPLATE.md`

### Planning Workflow:
1. Determine planning level (master/phase/wave/effort)
2. Copy appropriate template to target location
3. Fill in all placeholders with actual values
4. Validate completeness before proceeding
5. Transition to appropriate next state

## Planning Examples

### Example: Starting New Project
```bash
# Orchestrator in PLANNING state for new project

1. Copy master template:
   cp templates/MASTER-IMPLEMENTATION-PLAN.md ./IMPLEMENTATION-PLAN.md

2. Fill in project details:
   - Replace [PROJECT_NAME] with actual name
   - Define 5 phases with targets
   - Estimate total lines and efforts

3. Create phase plans:
   for i in {1..5}; do
     cp templates/PHASE-IMPLEMENTATION-PLAN.md phase-plans/PHASE-$i-PLAN.md
   done

4. Transition to WAVE_START when ready
```

### Example: Planning a Phase
```bash
# Orchestrator planning Phase 2

1. Load phase template:
   cp templates/PHASE-IMPLEMENTATION-PLAN.md phase-plans/PHASE-2-PLAN.md

2. Define waves (typically 3-5):
   - Wave 1: Core Business Logic
   - Wave 2: API Implementation  
   - Wave 3: Data Layer

3. Define efforts per wave:
   - E2.1.1: User Service (600 lines)
   - E2.1.2: Auth Service (500 lines)
   - E2.1.3: Session Management (400 lines)

4. Spawn Code Reviewer for each effort plan
```

## State Transitions

From PLANNING state:
- **PLAN_COMPLETE** → WAVE_START (Begin wave execution)
- **SPAWN_REVIEWER** → SPAWN_AGENTS (Spawn Code Reviewer for effort plans)
