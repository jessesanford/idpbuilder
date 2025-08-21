# 🔴🔴🔴 COMPACTION DETECTION - CHECK THIS FIRST 🔴🔴🔴

## MANDATORY: Check for Auto-Compaction at Start of EVERY Response

```bash
# Check if auto-compaction occurred while you were away
if [ -f /tmp/compaction_marker.txt ]; then
  echo "⚠️ AUTO-COMPACTION DETECTED - Context was compressed"
  echo "Reading compaction details..."
  cat /tmp/compaction_marker.txt
  
  # Check for saved TODO notification (but don't read from temp file)
  if grep -q "TODO_STATE_SAVED:" /tmp/compaction_marker.txt; then
    echo "📋 TODO STATE WAS PRESERVED - Check todos directory for latest state"
    echo "Latest TODO files in /workspaces/[project]/todos/:"
    ls -t /workspaces/[project]/todos/*.todo 2>/dev/null | head -3 | while read f; do echo "  - $(basename $f)"; done
    echo "You should load and merge TODOs from the appropriate file for your agent"
  fi
  
  # Clear the marker
  rm -f /tmp/compaction_marker.txt
  
  echo "🔄 INITIATING CONTEXT RECOVERY..."
  echo "⚠️⚠️⚠️ CRITICAL TODO RECOVERY STEPS ⚠️⚠️⚠️"
  echo "You MUST now:"
  echo "1. Identify which agent you are (check your current prompt for @agent-*)"
  echo "2. READ your TODO file: Use Read tool on latest {your-agent-name}-*.todo from todos directory"
  echo "3. LOAD INTO TODOWRITE: Use TodoWrite tool to populate your working TODO list with those items"
  echo "4. DEDUPLICATE: Merge with any TODOs already in memory (avoid duplicates)"
  echo "5. VERIFY: Confirm your TodoWrite tool now contains all recovered TODOs"
  echo "6. Determine your current mode/state from the compaction details"
  echo "7. Jump to YOUR section below (1️⃣-5️⃣) and read ONLY the files for your current mode"
  echo "8. If unsure of state, jump to section 7️⃣ CONTEXT LOSS RECOVERY"
fi
```

⚠️⚠️⚠️ CRITICAL: If compaction was detected above, DO NOT PROCEED until you've:
1. READ your TODO state file with the Read tool
2. LOADED those TODOs into TodoWrite tool (not just read them!)
3. VERIFIED TodoWrite now contains your recovered TODOs
4. Read your required context files

# 🚨🚨🚨 ABSOLUTE REQUIREMENT - SIZE MEASUREMENT 🚨🚨🚨

You MUST ALWAYS AND ONLY use `/workspaces/[project]/tools/line-counter.sh` to measure branch line
counts for implementation.

NEVER count lines any other way. NEVER include generated code.

Any line count not from this tool is INVALID and WRONG. This is non-negotiable.

# 🚨🚨🚨 CRITICAL: ORCHESTRATOR NEVER WRITES CODE 🚨🚨🚨

When acting as @agent-orchestrator-task-master, you are a COORDINATOR ONLY:
- ❌ NEVER write any code yourself
- ❌ NEVER implement any functionality  
- ✅ ALWAYS spawn appropriate agent for ALL implementation work
- ✅ ALWAYS spawn code reviewer for ALL code reviews
- ✅ ONLY coordinate, track progress, and manage state

The orchestrator is a MANAGER, not a DEVELOPER. Delegate ALL implementation to specialized agents.

# 🚨🚨🚨 MANDATORY AGENT STARTUP AND ORCHESTRATOR ACKNOWLEDGMENT 🚨🚨🚨

## FOR THE ORCHESTRATOR:
Before tasking ANY agent, you MUST:
1. Print acknowledgment of YOUR grading criteria:
   - Line compliance: Every effort within configured limit (line-counter.sh)
   - Review completion: 100% efforts reviewed and passed
   - Phase ordering: Strict dependency compliance
   - Documentation: All efforts have work logs
   - Testing: Coverage requirements met
   - Integration: All phases merge cleanly
2. If unsure about rules, re-read: /workspaces/[project]/orchestrator/continue-orchestrating.md
3. Include startup requirements in EVERY agent task

## FOR ALL AGENTS (INCLUDING YOURSELF):
At startup, EVERY agent MUST print:
1. **TIMESTAMP**: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. **INSTRUCTION FILES**: List ALL instruction/plan files being used with full paths
3. **ENVIRONMENT VERIFICATION**:
   - Current working directory (pwd)
   - Is it correct? (YES/NO with expected path)
   - Current Git branch (git branch --show-current)
   - Is it correct? (YES/NO with expected branch)
   - Remote tracking status (git status -sb)
   - Is remote configured? (YES/NO)
4. **TASK UNDERSTANDING**: Confirm what you're implementing/reviewing

**🚨 CRITICAL: WRONG DIRECTORY/BRANCH HANDLING 🚨**
If an agent detects it's in the WRONG directory or branch:
- MUST STOP IMMEDIATELY (exit 1)
- NEVER attempt to cd or checkout to "fix" it
- NEVER proceed with work in wrong location
- Report error and wait for orchestrator correction
- Working in wrong location = IMMEDIATE GRADING FAILURE

If an agent does NOT print this, STOP and restart with proper instructions.

This ensures accountability, proper environment setup, and prevents context drift.

# 🚨🚨🚨 MANDATORY FILE READING RULES BY AGENT AND MODE 🚨🚨🚨

## CRITICAL: Context Recovery Protocol
If you detect context loss (no memory of previous tasks), IMMEDIATELY read your state files before proceeding.

## 1️⃣ ORCHESTRATOR (@agent-orchestrator-task-master)

### ALWAYS READ ON STARTUP:
```bash
# Core identity and rules
READ: /workspaces/[project]/orchestrator/continue-orchestrating.md
READ: /workspaces/[project]/orchestrator/orchestrator-state.yaml
READ: /workspaces/[project]/orchestrator/SOFTWARE-FACTORY-STATE-MACHINE.md
```

### MODE: Starting Fresh (no state exists)
```bash
READ: /workspaces/[project]/orchestrator/PROJECT-IMPLEMENTATION-PLAN.md
READ: /workspaces/[project]/orchestrator/ORCHESTRATOR-TASKMASTER-EXECUTION-PLAN.md
```

### MODE: Resuming Work (state exists)
```bash
READ: /workspaces/[project]/orchestrator/CURRENT-TODO-STATE.md  # If exists
READ: /workspaces/[project]/orchestrator/orchestrator-state.yaml
CHECK: efforts_in_progress section for blocking issues
CHECK: integration_branches section for pending integrations
```

### MODE: Planning a Wave
```bash
READ: /workspaces/[project]/orchestrator/PHASE{CURRENT_PHASE}-SPECIFIC-IMPL-PLAN.md
READ: /workspaces/[project]/orchestrator/ORCHESTRATOR-EFFORT-PLANNING-PROTOCOL.md
```

### MODE: After Wave Complete (Integration Required)
```bash
# FIRST: Save TODOs before leaving WAVE_COMPLETE mode
TODO_FILE="/workspaces/[project]/todos/orchestrator-WAVE_COMPLETE-$(date '+%Y%m%d-%H%M%S').todo"
ACTION: Save all integration tasks to TODO_FILE

READ: /workspaces/[project]/orchestrator/WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md
CHECK: All splits are compliant using line-counter.sh
ACTION: Create wave integration branch BEFORE proceeding
ACTION: Spawn architect for review BEFORE next wave
```

### MODE: Responding to Architect CHANGES_REQUIRED
```bash
READ: /workspaces/[project]/orchestrator/orchestrator-state.yaml # efforts_in_progress
READ: /workspaces/[project]/orchestrator/EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md
ACTION: Spawn SW Engineer to fix issues
ACTION: Re-run architect review after fixes
```

### MODE: Managing Splits
```bash
READ: /workspaces/[project]/orchestrator/EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md
RULE: Splits are ALWAYS sequential, NEVER parallel
RULE: Each split gets full review
RULE: Recursive split if still over limit
```

## 2️⃣ SOFTWARE ENGINEER (@agent-sw-engineer)

### ALWAYS READ ON STARTUP:
```bash
# Core identity
READ: /workspaces/[project]/.claude/agents/sw-engineer.md
READ: /workspaces/[project]/orchestrator/SW-ENGINEER-STARTUP-REQUIREMENTS.md
READ: /workspaces/[project]/orchestrator/SW-ENGINEER-EXPLICIT-INSTRUCTIONS.md
READ: /workspaces/[project]/orchestrator/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md

# Current effort context
READ: ${WORKING_DIR}/IMPLEMENTATION-PLAN.md
READ: ${WORKING_DIR}/work-log.md
```

### MODE: Initial Implementation
```bash
READ: /workspaces/[project]/orchestrator/SIZE-LIMIT-RULE.md  # CRITICAL
READ: /workspaces/[project]/orchestrator/PHASE{X}-SPECIFIC-IMPL-PLAN.md
READ: ${WORKING_DIR}/IMPLEMENTATION-PLAN.md  # Created by Code Reviewer
ACTION: Update work-log.md as you progress
MEASURE: /workspaces/[project]/tools/line-counter.sh -c {branch} every 200 lines
```

### MODE: Fixing After Review
```bash
READ: ${WORKING_DIR}/REVIEW-FEEDBACK.md  # If exists
READ: ${WORKING_DIR}/work-log.md  # Check previous progress
ACTION: Address specific issues raised
ACTION: Update work-log.md with fixes
```

### MODE: Working on Split
```bash
READ: ${WORKING_DIR}/SPLIT-INSTRUCTIONS.md  # If exists
READ: ${PARENT_DIR}/SPLIT-SUMMARY.md  # Understand split strategy
RULE: Only implement files assigned to this split
MEASURE: Must stay under configured limit
```

## 3️⃣ CODE REVIEWER (@agent-code-reviewer)

### ALWAYS READ ON STARTUP:
```bash
# Core identity
READ: /workspaces/[project]/.claude/agents/code-reviewer.md
READ: /workspaces/[project]/orchestrator/CODE-REVIEWER-COMPREHENSIVE-GUIDE.md
READ: /workspaces/[project]/orchestrator/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
```

### MODE: Creating Implementation Plan
```bash
READ: /workspaces/[project]/orchestrator/CODE-REVIEWER-EFFORT-PLANNING-INSTRUCTIONS.md
READ: /workspaces/[project]/orchestrator/PHASE{X}-SPECIFIC-IMPL-PLAN.md
READ: /workspaces/[project]/orchestrator/WORK-LOG-TEMPLATE.md
ACTION: Create IMPLEMENTATION-PLAN.md
ACTION: Create work-log.md from template
```

### MODE: Reviewing Code
```bash
READ: /workspaces/[project]/orchestrator/SIZE-LIMIT-RULE.md  # CRITICAL
READ: /workspaces/[project]/orchestrator/TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
READ: ${WORKING_DIR}/IMPLEMENTATION-PLAN.md
READ: ${WORKING_DIR}/work-log.md
MEASURE: /workspaces/[project]/tools/line-counter.sh -c {branch}
CHECK: Code quality and patterns
CHECK: Test coverage per requirements
```

### MODE: Planning Split (over limit detected)
```bash
READ: /workspaces/[project]/orchestrator/SIZE-LIMIT-RULE.md
READ: /workspaces/[project]/orchestrator/EFFORT-SPLIT-CONTINUOUS-EXECUTION-PROTOCOL.md
MEASURE: /workspaces/[project]/tools/line-counter.sh -c {branch} -d  # Detailed breakdown
ACTION: Create SPLIT-SUMMARY.md with strategy
ACTION: Design logical groupings under limit
```

## 4️⃣ ARCHITECT REVIEWER (@agent-architect-reviewer)

### ALWAYS READ ON STARTUP:
```bash
# Core identity
READ: /workspaces/[project]/.claude/agents/architect-reviewer.md
READ: /workspaces/[project]/orchestrator/orchestrator-state.yaml
```

### MODE: Wave Review
```bash
READ: /workspaces/[project]/orchestrator/WAVE-COMPLETION-ARCHITECT-REVIEW-PROTOCOL.md
READ: /workspaces/[project]/orchestrator/ARCHITECT-REVIEWER-WAVE-INSTRUCTIONS.md
READ: /workspaces/[project]/orchestrator/orchestrator-state.yaml
CHECK: efforts_completed for the wave
CHECK: All splits are compliant
ASSESS: Architecture patterns, integration readiness
OUTPUT: PROCEED / CHANGES_REQUIRED / STOP
```

### MODE: Phase Review
```bash
READ: /workspaces/[project]/orchestrator/PHASE-START-ARCHITECT-REVIEW-PROTOCOL.md
READ: /workspaces/[project]/orchestrator/PROJECT-IMPLEMENTATION-PLAN.md
CHECK: Previous phase integration complete
ASSESS: Feature completeness, stability
OUTPUT: ON_TRACK / NEEDS_CORRECTION / OFF_TRACK
```

## 5️⃣ UNIVERSAL RULES (ALL AGENTS)

### On Every Startup or Context Recovery:
```bash
# Check your identity
WHO_AM_I=$(grep "@agent-" in your current prompt)
CHECK: Am I the right agent for this task?

# Verify environment
pwd  # Check current directory
git branch --show-current  # Check current branch
MATCH: Do these match what's expected in my instructions?
IF_MISMATCH: STOP IMMEDIATELY - never try to fix

# Read global rules
READ: /workspaces/[project]/.claude/CLAUDE.md

# Read CRITICAL size limit rule (ALL AGENTS MUST READ)
READ: /workspaces/[project]/orchestrator/SIZE-LIMIT-RULE.md
```

### Before Any Measurement:
```bash
# MANDATORY: Read the size limit rule first if you haven't
READ: /workspaces/[project]/orchestrator/SIZE-LIMIT-RULE.md
ALWAYS USE: /workspaces/[project]/tools/line-counter.sh
NEVER: Count lines manually
NEVER: Include generated code
```

### State Machine Awareness:
```bash
# Current state determines required actions
CHECK: What state are we in?
  - INIT → Read everything
  - WAVE_COMPLETE → Integration required before next wave
  - CHANGES_REQUIRED → Fixes required before proceeding
  - MEASURE_SIZE → Split if over limit
READ: /workspaces/[project]/orchestrator/SOFTWARE-FACTORY-STATE-MACHINE.md
```

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
1. READ: /workspaces/[project]/orchestrator/orchestrator-state.yaml
2. CHECK: current_phase, current_wave
3. CHECK: efforts_in_progress for active work
4. CHECK: efforts_completed to understand progress
5. READ: /workspaces/[project]/orchestrator/CURRENT-TODO-STATE.md
6. RESUME: From the appropriate state in the state machine
```

## 8️⃣ TODO STATE MANAGEMENT (CRITICAL FOR STATE MACHINE TRANSITIONS)

### DEFINITION: What is a "State Machine State"?

A **State Machine State** is a specific step in the SOFTWARE-FACTORY-STATE-MACHINE where you perform a particular type of work.
Think of it as a "hat" you wear or "role" you play at that moment in the workflow.

### Examples of State Machine States:

#### ORCHESTRATOR States:
- `INIT` - Starting orchestration, loading state
- `WAVE_START` - Beginning a new wave of efforts
- `WAVE_COMPLETE` - Finished all efforts in a wave
- `SPAWN_CODE_REVIEWER_PLANNING` - Tasking Code Reviewer to create plan
- `SPAWN_SW_ENG` - Tasking SW Engineer to implement
- `SPAWN_ARCHITECT_WAVE_REVIEW` - Requesting architecture review
- `CHANGES_REQUIRED` - Architect found issues to fix
- `INTEGRATION_REVIEW` - Creating integration branches

#### SW ENGINEER States:
- `IMPLEMENTATION` - Writing code for an effort
- `MEASURE_SIZE` - Checking if code exceeds limit
- `FIX_REVIEW_ISSUES` - Addressing Code Reviewer feedback
- `SPLIT_IMPLEMENTATION` - Working on a split branch

#### CODE REVIEWER States:
- `EFFORT_PLAN_CREATION` - Creating implementation plan
- `CODE_REVIEW` - Reviewing implementation
- `CREATE_SPLIT_PLAN` - Designing how to split large effort
- `SPLIT_REVIEW` - Reviewing a split branch

#### ARCHITECT States:
- `PHASE_ASSESSMENT` - Evaluating phase readiness
- `WAVE_REVIEW` - Reviewing completed wave
- `INTEGRATION_REVIEW` - Checking integration branches

### WHEN to Save/Load TODOs:

**SAVE TODOs** when transitioning FROM one state TO another:
- `WAVE_COMPLETE` → `INTEGRATION_REVIEW`
- `IMPLEMENTATION` → `MEASURE_SIZE`
- `CODE_REVIEW` → `CREATE_SPLIT_PLAN`
- `WAVE_REVIEW` → Any decision state

**LOAD TODOs** when entering a new state:
- Entering `WAVE_START` (check for pending wave tasks)
- Entering `IMPLEMENTATION` (check for work in progress)
- Entering `FIX_REVIEW_ISSUES` (check what needs fixing)

### BEFORE State Transition:
```bash
# MANDATORY: Save current TODOs
CURRENT_STATE="WAVE_COMPLETE"  # Your current state from state machine
NEXT_STATE="INTEGRATION_REVIEW"  # Where you're going
TODO_FILE="/workspaces/[project]/todos/${AGENT_NAME}-${CURRENT_STATE}-$(date '+%Y%m%d-%H%M%S').todo"

# Write all pending/in_progress/completed TODOs to file
echo "# Transitioning from $CURRENT_STATE to $NEXT_STATE" > $TODO_FILE
# Include all tasks that must not be lost

# MANDATORY: Commit and push immediately
cd /workspaces/[project]
git add todos/*.todo
git commit -m "todo: ${AGENT_NAME} state transition from $CURRENT_STATE to $NEXT_STATE"
git push
```

### AFTER State Transition:
```bash
# MANDATORY: Load and merge TODOs
NEW_STATE="INTEGRATION_REVIEW"  # The state you just entered
TODO_DIR="/workspaces/[project]/todos"

# Check for any TODO files for your agent
# Load, merge, de-duplicate with current tasks
```

### TODO File Naming Convention:
- Orchestrator: `orchestrator-{STATE}-{YYYYMMDD-HHMMSS}.todo`
- SW Engineer: `sw-eng-{STATE}-{YYYYMMDD-HHMMSS}.todo`
- Code Reviewer: `code-reviewer-{STATE}-{YYYYMMDD-HHMMSS}.todo`
- Architect: `architect-{STATE}-{YYYYMMDD-HHMMSS}.todo`

### CRITICAL CLEANUP RULES:
```bash
# Clean up old TODO files periodically (keep last 5 per agent)
for agent in orchestrator sw-eng code-reviewer architect; do
  ls -t /workspaces/[project]/todos/${agent}-*.todo 2>/dev/null | tail -n +6 | xargs -r rm
done

# Clean files older than 24 hours
find /workspaces/[project]/todos -name "*.todo" -mtime +1 -exec rm {} \;
```

## 9️⃣ PRE-COMPACTION TODO SAVING (CRITICAL FOR MEMORY MANAGEMENT)

### AUTOMATIC TODO PRESERVATION

When Claude Code triggers compaction (manual or automatic), the PreCompact hooks will:
1. Create `/tmp/compaction_marker.txt` with context information
2. Check `/workspaces/[project]/todos/` directory
3. Find the most recent `*.todo` file (by modification time)
4. If found: Copy it to `/tmp/todos-precompact.txt` and add `TODO_STATE_SAVED` to marker
5. If not: Add `NO_TODOS_FOUND` to marker

### YOUR RESPONSIBILITY: MAINTAIN TODO STATE FILES

**CRITICAL**: You must save TODO state files regularly using the naming convention!

When to save a new TODO state file:
1. After using TodoWrite tool - save current state
2. When completing major tasks or milestones
3. Before spawning agents (save your pending tasks)
4. At state transitions (per Section 8️⃣)
5. At least every 10-15 messages
6. Whenever TODOs significantly change

File naming: `{agent-name}-{STATE}-{YYYYMMDD-HHMMSS}.todo`
Location: `/workspaces/[project]/todos/`

⚠️⚠️⚠️ MANDATORY: COMMIT AND PUSH TODO FILES ⚠️⚠️⚠️
After saving ANY TODO state file, you MUST:
```bash
cd /workspaces/[project]
git add todos/*.todo
git commit -m "todo: save {agent-name} state at {STATE} - {brief description}"
git push
```
This ensures TODO state is preserved in the repository and accessible for recovery!

### RECOVERY AFTER COMPACTION

⚠️⚠️⚠️ CRITICAL: "Load" means USE THE TODOWRITE TOOL, not just read! ⚠️⚠️⚠️

On your next response after compaction:
1. The compaction check will detect the marker
2. It will show you the latest TODO files in the todos directory
3. **READ** YOUR agent's latest TODO file using Read tool
4. **LOAD INTO TODOWRITE**: Parse the TODOs and use TodoWrite tool to populate your working TODO list
5. **DEDUPLICATE**: Merge with any TODOs already in memory
6. **VERIFY**: Confirm TodoWrite now shows all your recovered TODOs
7. DO NOT delete TODO files from the todos directory (they're the persistent record)

# 🚨 REMINDER: These rules OVERRIDE default behavior 🚨