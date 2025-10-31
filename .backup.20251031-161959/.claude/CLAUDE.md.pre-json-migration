# SOFTWARE FACTORY 2.0 - PROJECT CLAUDE.MD

This file contains PROJECT-SPECIFIC configurations and overrides.
Most rules are now managed centrally:
- **Agent configs**: `$CLAUDE_PROJECT_DIR/.claude/agents/[agent].md`
- **Rules library**: `$CLAUDE_PROJECT_DIR/rule-library/`
- **State rules**: `$CLAUDE_PROJECT_DIR/agent-states/[agent]/[state]/rules.md`

# 🔴🔴🔴 COMPACTION DETECTION - CHECK THIS FIRST 🔴🔴🔴

## Quick Compaction Check
```bash
# The utilities/check-compaction.sh script handles detection and recovery guidance
bash $CLAUDE_PROJECT_DIR/utilities/check-compaction.sh
```

**Note**: TODO persistence is handled by rule R287.
PreCompact hooks require `.claude/settings.json` - see CRITICAL-SETTINGS-JSON.md

# ⚠️⚠️⚠️ CRITICAL TODO PERSISTENCE REMINDER ⚠️⚠️⚠️

**YOU MUST SAVE TODOs FREQUENTLY OR LOSE ALL WORK:**
- 🚨 Save within 30 seconds after TodoWrite
- 🚨 Save every 10 messages OR 15 minutes (whichever comes first)
- 🚨 Save before EVERY state transition
- 🚨 Commit and push within 60 seconds of saving
- **See Section 8️⃣ below for MANDATORY enforcement details**
- **Failure to save = -20% to -100% grading penalty**

# 🚨🚨🚨 ** ACKNOWLEDGE YOUR GRADING CRITERIA** 🚨🚨🚨 
Before tasking any agent, acknowledge:
```markdown
📊 ORCHESTRATOR GRADING CRITERIA:
### Grading Criteria

You will be graded on:

1. WORKSPACE ISOLATION (20%)
   - ✓ ALL agents confined to assigned working copies
   - ✓ No agents escape or modify other working copies
   - ✓ No branch/remote tracking deviations
   - ✓ Protection verification before any work
   
2. WORKFLOW COMPLIANCE (25%)
   - ✓ Code Reviewer agent spawned to perform CODE REVIEW PROTOCOL
   - ✓ Measure after every logical change group
   - ✓ ALL code committed and pushed (zero uncommitted)
   - ✓ Immediate review after development
   - ✓ Sequential handling of split branches
   - ✓ CODE REVIEW REPORT/INSTRUCTIONS markdown created for all findings for Orchestrator agent to give to SWEs

3. SIZE COMPLIANCE (20%)
   - ✓ Code Reviewer agent spawned to perform SPLIT PLANNING PROTOCOL
   - ✓ Zero PRs >700 ~soft (> 800 !HARD!) lines committed
   - ✓ Line counter run regularly during development
   - ✓ Immediate stop and split when violations detected
   - ✓ SPLIT PLAN INSTRUCTIONS / REPORT created before splitting for Orchestrator agent to give to SWEs

4. PARALLELIZATION (15%)
   - ✓ Multiple SWE agents spawned for independent work (if allowed by plan)
   - ✓ Single SWE for sequential split work (again check the plan)
   - ✓ **ALL AGENTS MUST emit the current timestamp as their first task on startup, YOU WILL BE GRADED ON HOW CLOSE THEIR TIMESTAMPS ARE. >5s DEVIATION BETWEEN PARALLEL AGENTS IS A FAILURE!** [R151]
   - ✓ Immediate reviewer deployment after SWE completion
   

5. QUALITY ASSURANCE (20%)
   - ✓ ALL tests passing before completion
   - ✓ ALL review issues resolved
   - ✓ No features left incomplete
   - ✓ Proper review-fix cycles until clean
   - ✓ State file updates after EVERY transition
   - ✓ Monitoring progress every 5 messages
   - ✓ TODO Persistence: R287 full compliance

FAILURE CONDITIONS:
- Agents corrupt/pollute working copies = FAIL
- Any PR >800 lines merged = FAIL
- Agents deployed sequentially when parallelization allowed = FAIL  
- Todo list not maintained = FAIL
- Reviews ignored = FAIL
- Orchestrator writes ANY code = FAIL
```

# 🚨🚨🚨 ABSOLUTE REQUIREMENT - SIZE MEASUREMENT 🚨🚨🚨

You MUST ALWAYS AND ONLY use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` to measure branch line
counts for implementation (if it exists).

NEVER count lines any other way. NEVER include generated code.

Any line count not from this tool is INVALID and WRONG. This is non-negotiable.

# 🚨🚨🚨 CRITICAL: ORCHESTRATOR NEVER WRITES CODE 🚨🚨🚨

When acting as @agent-orchestrator, you are a COORDINATOR ONLY:
- ❌ NEVER write any code yourself
- ❌ NEVER implement any functionality  
- ✅ ALWAYS spawn appropriate agent for ALL implementation work
- ✅ ALWAYS spawn code reviewer for ALL code reviews
- ✅ ONLY coordinate, track progress, and manage state

The orchestrator is a MANAGER, not a DEVELOPER. Delegate ALL implementation to specialized agents.

# 🚨🚨🚨 MANDATORY AGENT STARTUP 🚨🚨🚨

**All agents follow R203 (State-Aware Startup) which requires:**
1. Load core agent config
2. Read TODO persistence rule (R287)
3. Determine current state
4. Load state-specific rules
5. Acknowledge all rules

**Critical for parallelization:**
- Timestamps must be within 5s for parallel agents (R151)
- Wrong directory/branch = immediate stop (never try to fix)

See: `$CLAUDE_PROJECT_DIR/rule-library/R203-state-aware-agent-startup.md`

# 🚨🚨🚨 PROJECT-SPECIFIC CONFIGURATION 🚨🚨🚨

## State-Specific Rules
Each agent has detailed state-specific rules in:
`$CLAUDE_PROJECT_DIR/agent-states/[agent-name]/[STATE]/rules.md`

**Available States:**
- **Orchestrator**: INIT, PLANNING, SETUP_EFFORT_INFRASTRUCTURE, SPAWN_AGENTS, MONITOR, WAVE_COMPLETE, INTEGRATION, ERROR_RECOVERY, etc.
- **SW-Engineer**: INIT, IMPLEMENTATION, SPLIT_IMPLEMENTATION, MEASURE_SIZE, FIX_ISSUES, TEST_WRITING, etc.
- **Code-Reviewer**: INIT, EFFORT_PLAN_CREATION, CODE_REVIEW, CREATE_SPLIT_PLAN, SPLIT_REVIEW, etc.
- **Architect**: INIT, WAVE_REVIEW, PHASE_ASSESSMENT, INTEGRATION_REVIEW, etc.

## Available Resources in This Template

**Core Files:**
- `$CLAUDE_PROJECT_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md` - State transitions
- `$CLAUDE_PROJECT_DIR/orchestrator-state.json` - Example state file
- `$CLAUDE_PROJECT_DIR/.claude/agents/` - Agent configurations
- `$CLAUDE_PROJECT_DIR/.claude/commands/` - Continuation commands
- `$CLAUDE_PROJECT_DIR/rule-library/` - All rule definitions
- `$CLAUDE_PROJECT_DIR/agent-states/` - State-specific rules
- `$CLAUDE_PROJECT_DIR/templates/` - Planning templates
- `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` - Line counting tool
- `$CLAUDE_PROJECT_DIR/utilities/` - Recovery scripts

**Create for Your Project:**
- `orchestrator-state.json` - Copy from .example
- `protocols/` - Project-specific protocols
- `phase-plans/` - Phase planning documents
- `PROJECT-IMPLEMENTATION-PLAN.md` - Master plan
- `todos/` - Auto-created when saving TODOs

## 📁 AGENT-SPECIFIC RULES

**Each agent has detailed state-specific rules that MUST be followed:**

### 1️⃣ ORCHESTRATOR (@agent-orchestrator)
- **Config**: `$CLAUDE_PROJECT_DIR/.claude/agents/orchestrator.md`
- **State Rules**: `$CLAUDE_PROJECT_DIR/agent-states/orchestrator/[STATE]/rules.md`
- **States**: INIT, PLANNING, SETUP_EFFORT_INFRASTRUCTURE, SPAWN_AGENTS, MONITOR, WAVE_COMPLETE, INTEGRATION, ERROR_RECOVERY

### 2️⃣ SOFTWARE ENGINEER (@agent-sw-engineer)
- **Config**: `$CLAUDE_PROJECT_DIR/.claude/agents/sw-engineer.md`
- **State Rules**: `$CLAUDE_PROJECT_DIR/agent-states/sw-engineer/[STATE]/rules.md`
- **States**: INIT, IMPLEMENTATION, SPLIT_IMPLEMENTATION, MEASURE_SIZE, FIX_ISSUES, TEST_WRITING

### 3️⃣ CODE REVIEWER (@agent-code-reviewer)
- **Config**: `$CLAUDE_PROJECT_DIR/.claude/agents/code-reviewer.md`
- **State Rules**: `$CLAUDE_PROJECT_DIR/agent-states/code-reviewer/[STATE]/rules.md`
- **States**: INIT, EFFORT_PLAN_CREATION, CODE_REVIEW, CREATE_SPLIT_PLAN, SPLIT_REVIEW

### 4️⃣ ARCHITECT (@agent-architect)
- **Config**: `$CLAUDE_PROJECT_DIR/.claude/agents/architect.md`
- **State Rules**: `$CLAUDE_PROJECT_DIR/agent-states/architect/[STATE]/rules.md`
- **States**: INIT, WAVE_REVIEW, PHASE_ASSESSMENT, INTEGRATION_REVIEW

**IMPORTANT**: Agents MUST read their state-specific rules based on their current state. The rules in these files override any general guidance.

## 5️⃣ UNIVERSAL RULES (ALL AGENTS)

**Key Rules:**
- R203: State-aware startup and acknowledgment
- R287: Comprehensive TODO persistence requirements
- R206: State machine validation
- R216: Bash execution syntax
- R220/R221: Size limits and CD requirements

**Line Counting:**
- ALWAYS use: `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` (if exists)
- NEVER count manually or include generated code

**State Machine:**
- Read: `$CLAUDE_PROJECT_DIR/SOFTWARE-FACTORY-STATE-MACHINE.md`  # In root directory
- Validate all transitions with R206

## 6️⃣ CRITICAL GATES (ENFORCEMENT POINTS)

### Wave Completion Gate:
```
BEFORE starting next wave, MUST have:
✅ All effort splits within limit
✅ Wave integration branch created
✅ Architect review completed
✅ State file updated
```

### Phase Transition Gate:
```
BEFORE starting next phase, MUST have:
✅ All waves integrated
✅ Phase integration branch created
✅ Architect phase assessment
✅ No OFF_TRACK status
```

### Split Gate:
```
WHEN effort exceeds limit, MUST:
✅ Stop implementation
✅ Create split plan
✅ Execute splits sequentially
✅ Each split gets review
```

## 7️⃣ CONTEXT LOSS RECOVERY

If you lose context and don't remember previous work:
```bash
1. READ: $CLAUDE_PROJECT_DIR/orchestrator-state.json
2. CHECK: current_phase, current_wave
3. CHECK: efforts_in_progress for active work
4. CHECK: efforts_completed to understand progress
5. READ: $CLAUDE_PROJECT_DIR/CURRENT-TODO-STATE.md
6. RESUME: From the appropriate state in the state machine
```

## 8️⃣ TODO STATE MANAGEMENT - CRITICAL ENFORCEMENT

### 🚨 MANDATORY TODO PERSISTENCE RULES 🚨

**ALL AGENTS MUST FOLLOW THESE TODO RULES OR FACE GRADING PENALTIES:**

#### R287: Comprehensive TODO Persistence (BLOCKING)
**MUST SAVE when ANY of these occur:**
- ✅ After using TodoWrite tool (+30 seconds MAX)
- ✅ Before ANY state machine transition
- ✅ Before spawning another agent
- ✅ After completing effort/wave/review
- ✅ When encountering errors or blocks

**Penalty:** -20% for each missed trigger, -50% for lost TODOs

Consolidated into R287 above - includes save frequency requirements
**MUST SAVE at these intervals:**
- ⏰ Every 10 messages exchanged
- ⏰ Every 15 minutes of work (MANDATORY)
- ⏰ After 200 lines of code written
- ⏰ After 3 files modified
- ⏰ Before high-memory operations

**Penalty:** -15% per violation, -100% for total TODO loss

Consolidated into R287 above - includes commit protocol
**After EVERY save:**
```bash
cd $CLAUDE_PROJECT_DIR
git add todos/*.todo
git commit -m "todo: ${AGENT_NAME} - ${TRIGGER_REASON}"
git push
```
**Penalty:** -10% for delayed commits, -30% for uncommitted TODOs

Consolidated into R287 above - includes recovery verification
**After compaction/recovery:**
- READ saved TODO file
- LOAD into TodoWrite tool (MANDATORY - not just read!)
- VERIFY all tasks recovered
- DELETE processed recovery files

### 📝 TODO File Format
```
todos/${AGENT_NAME}-${STATE}-${YYYYMMDD-HHMMSS}.todo
```
Examples:
- `orchestrator-WAVE_COMPLETE-20250120-143000.todo`
- `sw-engineer-IMPLEMENTATION-20250120-145500.todo`

### ⚠️ ENFORCEMENT TRACKING
Agents MUST self-monitor:
```bash
# Check time since last save
LAST_SAVE_AGO=$(($(date +%s) - LAST_TODO_SAVE))
if [ $LAST_SAVE_AGO -gt 900 ]; then  # >15 minutes
    echo "🚨 VIOLATION: Exceeding R287 frequency requirement!"
    save_todos "OVERDUE_CHECKPOINT"
fi
```

### 📊 Grading Impact Summary
- Missing trigger save: -20%
- Missing frequency save: -15%
- Missing commit: -10%
- Lost TODOs in compaction: -50%
- Total TODO loss: -100% (IMMEDIATE FAILURE)

**Full rule details:** `$CLAUDE_PROJECT_DIR/rule-library/R287-todo-persistence-comprehensive.md`
**Agent functions:** `$CLAUDE_PROJECT_DIR/.claude/agents/[agent-name].md`

# 🚨 REMINDER: These rules OVERRIDE default behavior 🚨