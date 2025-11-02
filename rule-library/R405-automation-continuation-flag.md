# 🔴🔴🔴 RULE R405 - Automation Continuation Flag [SUPREME LAW]

# 🚨🚨🚨 CRITICAL ALERT - READ BEFORE USING THIS RULE 🚨🚨🚨

**SYSTEMIC VIOLATION DETECTED:**

Agents are repeatedly setting `CONTINUE-SOFTWARE-FACTORY=FALSE` for normal
operations, including PROJECT_DONE scenarios and R322 checkpoints. This is
**fundamentally wrong** and defeats the entire purpose of the Software Factory
automation system.

**MOST COMMON VIOLATION:**
> "R322 requires mandatory stop before state transition, therefore I should set
> CONTINUE-SOFTWARE-FACTORY=FALSE"

**THIS IS WRONG! R322 stops ≠ FALSE flag!**

**R322 checkpoints are NORMAL OPERATION:**
- State work completed ✅
- Waiting for /continue-orchestrating ✅
- Ready for next state transition ✅
- **Use TRUE because this is DESIGNED WORKFLOW!**

**OTHER VIOLATIONS INCLUDE:**
- Setting FALSE after successful integration (PROJECT_DONE!)
- Setting FALSE in WAITING states when results arrive (NORMAL!)
- Setting FALSE when spawning agents (normal workflow!)
- Setting FALSE when code review finds issues (designed process!)

**IMMEDIATE CORRECTION REQUIRED:**

Before setting this flag in ANY state, read the master guide:
**→ R405-CONTINUATION-FLAG-MASTER-GUIDE.md ←**

**If you set FALSE incorrectly:**
- First violation: -20% penalty
- Pattern (3+): -50% penalty
- Systematic misuse: -100% penalty

**The user is FRUSTRATED by unnecessary stops. Default to TRUE!**

---

## Purpose
Enable full automation of Software Factory 2.0 by providing a universal machine-parseable signal for continuation decisions.

## The Rule
**ALL agents in ALL states MUST output a continuation flag as their LAST output before completing the state.**

## Format (EXACT - NO VARIATIONS)

**CRITICAL: The flag MUST be emitted with the full variable name, NOT just the value!**

### ❌ WRONG - DO NOT USE:
```bash
# These are INVALID and will break automation
TRUE
FALSE
R405 Continuation Flag: TRUE
```

### ✅ CORRECT - Basic Format (REQUIRED):
```bash
# Success - continue to next state
CONTINUE-SOFTWARE-FACTORY=TRUE

# Failure/block - stop for human intervention
CONTINUE-SOFTWARE-FACTORY=FALSE
```

### Enhanced Format with Checkpoint Context (RECOMMENDED):
```bash
# R322 checkpoint (spawn/critical transition) - normal operation
CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322

# Normal operation continuation - no checkpoint
CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=NONE

# Error with specific reason - requires intervention
CONTINUE-SOFTWARE-FACTORY=FALSE REASON=ERROR_TYPE
```

**Why Use Enhanced Format?**
- Eliminates ambiguity about WHY agent stopped
- Allows test framework to auto-continue at R322 checkpoints
- Makes debugging easier (explicit checkpoint type in logs)
- Backward compatible (basic format still works)

## Critical Requirements
1. **UNIVERSAL**: Every single state for every single agent must have this
2. **CONSISTENT**: Always use exactly "CONTINUE-SOFTWARE-FACTORY" (not CONTINUE-ORCHESTRATING or any variant)
3. **COMPLETE**: Always include the full variable assignment "CONTINUE-SOFTWARE-FACTORY=TRUE" NOT just "TRUE"
4. **LAST OUTPUT**: Must be the absolute last line of output when agent completes state
5. **GREPPABLE**: Must be on its own line for easy grep/parsing
6. **MANDATORY**: No state can omit this flag - -100% penalty for omission

## Checkpoint Context Types (Enhanced Format)

When using the enhanced format with context, use these standardized types:

### Checkpoint Types (for TRUE flag):

| Type | When to Use | Example |
|------|-------------|---------|
| `CHECKPOINT=R322` | R322 mandatory stop (spawn/critical transition) | After spawning multiple agents |
| `CHECKPOINT=NONE` | Normal operation, no checkpoint | Regular state completion |

### Reason Types (for FALSE flag):

| Type | When to Use | Example |
|------|-------------|---------|
| `REASON=ERROR` | Generic unrecoverable error | Unknown system failure |
| `REASON=STATE_CORRUPTION` | State file corrupted or invalid | JSON schema validation failed |
| `REASON=MISSING_FILES` | Critical files missing | orchestrator-state-v3.json not found |
| `REASON=RECURSIVE_SPLIT` | Split of split detected | Effort split still >800 lines |
| `REASON=ITERATION_OVERFLOW` | Exceeded max iterations | Integration failed >10 times |
| `REASON=WRONG_LOCATION` | Agent in wrong directory/branch | cwd != expected workspace |

### Automation Framework Behavior

Based on the flag context, the test/automation framework will:

- `TRUE CHECKPOINT=R322` → **Auto-continue** (R322 is normal operation, designed workflow)
- `TRUE CHECKPOINT=NONE` → **Auto-continue** (normal progression)
- `TRUE` (no context) → **Auto-continue** (backward compatible)
- `FALSE REASON=*` → **STOP** for human intervention
- `FALSE` (no context) → **STOP** for human intervention
- Missing flag → **STOP** with warning

## DEFAULT BEHAVIOR: TRUE FOR NORMAL OPERATIONS
**CRITICAL**: The default should be TRUE! Only set FALSE for EXCEPTIONAL situations.

**NOTE**: Even at R322 checkpoints, use TRUE! R322 checkpoints are NORMAL workflow, not errors.

## When to Output FALSE (EXCEPTIONAL CASES ONLY)
**ONLY set FALSE for these TRULY EXCEPTIONAL situations:**
- **Unrecoverable errors** - System cannot proceed at all
- **Missing CRITICAL files** - Required configuration or state files missing/corrupted
- **Invalid state machine** - Current state not in state machine
- **Wrong working directory/branch** - Agent in completely wrong location
- **Recursive splits** - When a split needs to be split AGAIN (split of a split)
- **Multiple fix cascades** - More than 3 cascades with recurring failures (4+ is exceptional)
- **Integration process crashed** - Integration agent died without producing report

**NOT reasons to set FALSE:**
- ❌ Integration build failures (NORMAL - triggers fix cascade)
- ❌ Integration test failures (NORMAL - triggers fix cascade)
- ❌ Integration demo failures (NORMAL - triggers fix cascade)
- ❌ Failed effort with 0 lines (NORMAL - can spawn reviewer)
- ❌ Normal state transitions (monitoring → review)
- ❌ Routine review failures (spawn fixes)
- ❌ First-time split needed (normal operation)
- ❌ Tests failing once, twice, or three times (can be fixed)
- ❌ Normal monitoring transitions
- ❌ First, second, or third fix cascade (system handles automatically)

## When to Output TRUE (DEFAULT FOR NORMAL OPERATIONS)
**This should be your DEFAULT choice! Set TRUE for:**
- ✅ State completed (even if some efforts failed)
- ✅ Integration failed with build/test/demo issues (spawn fix cascade - NORMAL!)
- ✅ Normal transitions (monitor → review, review → fix, etc.)
- ✅ Ready for next state
- ✅ Splits needed (first time - this is normal)
- ✅ Reviews found issues (spawn fixes - normal flow)
- ✅ Some tests failing (can be fixed through fix cascades)
- ✅ Implementation complete (spawn review)
- ✅ Effort produced 0 lines (spawn reviewer anyway)
- ✅ Any routine state transition
- ✅ First, second, or third fix cascade (system handles automatically)
- ✅ Need to spawn Code Reviewer for fix plans (part of normal flow)
- ✅ Need to transition to IMMEDIATE_BACKPORT_REQUIRED (normal for integration failures)

**RULE OF THUMB**: If the system CAN proceed automatically, it SHOULD (TRUE).

## Implementation Pattern
```bash
# At the very end of state execution, after ALL other work:

# DEFAULT TO TRUE unless EXCEPTIONAL situation
CONTINUE_FLAG="TRUE"  # Start with TRUE as default

# Only set FALSE for TRULY EXCEPTIONAL cases
if [[ "$UNRECOVERABLE_ERROR" == true ]] || \
   [[ "$MISSING_CRITICAL_FILES" == true ]] || \
   [[ "$RECURSIVE_SPLIT_NEEDED" == true ]] || \
   [[ "$MULTIPLE_FIX_CASCADES" == true ]]; then
    CONTINUE_FLAG="FALSE"
    echo "❌ Exceptional situation detected - manual intervention required"
fi

echo "CONTINUE-SOFTWARE-FACTORY=$CONTINUE_FLAG"
```

## Automation Usage
This enables external automation scripts to:
```bash
# Run agent and capture output
OUTPUT=$(agent_execution_command)

# Check continuation flag
if echo "$OUTPUT" | grep -q "CONTINUE-SOFTWARE-FACTORY=TRUE$"; then
    # Automatically continue to next state
    continue_to_next_state
else
    # Stop and alert human
    alert_human "Manual intervention required"
    exit 1
fi
```

## State Rule Implementation
Every state rule file (`agent-states/*/[STATE]/rules.md`) must include:
```markdown
### 🔴 Automation Continuation Flag (R405)
**MANDATORY - LAST OUTPUT**: Before completing this state, output automation flag:
- If state completed successfully and SF2.0 should continue: `CONTINUE-SOFTWARE-FACTORY=TRUE`
- If error/block/manual review needed: `CONTINUE-SOFTWARE-FACTORY=FALSE`

This MUST be the absolute last line of output before state completion.
```

## Enforcement
- **Grading Impact**: -100% for omitting this flag
- **Validation**: Every state transition must be preceded by this flag
- **Monitoring**: Automation scripts will fail if flag is missing
- **Recovery**: Missing flag = immediate stop, manual intervention required

## Examples

### R322 Checkpoint - After Spawning Agents (TRUE CHECKPOINT=R322)
```bash
echo "✅ Spawned 3 SW Engineers for implementation"
echo "📊 Agents working on efforts: E1.1, E1.2, E1.3"
update_state "MONITORING_SWE_PROGRESS"
commit_state()

echo "🛑 R322: Checkpoint after spawn (context preservation)"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=R322"  # Enhanced format!
exit 0  # Stop this conversation turn
```

**What happens next:**
1. Framework sees `TRUE CHECKPOINT=R322`
2. Framework knows this is normal R322 checkpoint
3. Framework auto-continues with `/continue-software-factory`
4. Orchestrator restarts, reads state=MONITORING_SWE_PROGRESS
5. Orchestrator executes monitoring state

### Normal Operations - Integration Failed (TRUE CHECKPOINT=NONE)
```bash
echo "🔴 Integration build failed with duplicate symbol errors"
echo "🔴 Tests failed: 5 test suites failing"
echo "🔴 Demo not working due to build failures"
echo "⚠️ This is NORMAL - spawning fix cascade"
echo "📝 Transitioning to SPAWN_CODE_REVIEWER_FIX_PLAN"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE CHECKPOINT=NONE"  # Normal flow, no checkpoint
```

### Normal Operations - Implementation Complete (TRUE)
```bash
echo "✅ Implementation monitoring complete"
echo "✅ Some efforts succeeded, one produced 0 lines"
echo "✅ Ready to spawn code reviewers"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Normal flow continues
```

### Normal Operations - Review Found Issues (TRUE)
```bash
echo "⚠️ Review found issues requiring fixes"
echo "🔄 Will spawn SW Engineers for fixes"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Normal fix cycle
```

### Normal Operations - Split Needed (TRUE)
```bash
echo "📏 Effort exceeds size limit, split required"
echo "🔄 Will execute split implementation"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # First split is normal
```

### EXCEPTIONAL - Missing Critical File (FALSE REASON=MISSING_FILES)
```bash
echo "❌ ERROR: Required state file is corrupted/missing"
echo "❌ Cannot proceed without valid state file"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=MISSING_FILES"  # Enhanced format with reason
```

**What happens next:**
1. Framework sees `FALSE REASON=MISSING_FILES`
2. Framework knows specific error type
3. Framework STOPS and alerts human
4. Human investigates missing files before continuing

### EXCEPTIONAL - Recursive Split (FALSE REASON=RECURSIVE_SPLIT)
```bash
echo "❌ ERROR: Split E3.1.2-split-1 STILL exceeds limits"
echo "❌ Recursive split required - needs human review"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=RECURSIVE_SPLIT"  # Enhanced format
```

### EXCEPTIONAL - Multiple Cascades (FALSE REASON=ITERATION_OVERFLOW)
```bash
echo "❌ ERROR: This is the 4th fix cascade for same issues"
echo "❌ Pattern indicates fundamental problem"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=ITERATION_OVERFLOW"  # Enhanced format
```

### EXCEPTIONAL - State Corruption (FALSE REASON=STATE_CORRUPTION)
```bash
echo "❌ ERROR: State file JSON schema validation failed"
echo "❌ current_state field is missing or invalid"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE REASON=STATE_CORRUPTION"  # Enhanced format
```

## Universality Requirement
This rule applies to:
- All orchestrator states
- All sw-engineer states
- All code-reviewer states
- All architect states
- All integration states
- All software-factory-manager states
- All future agents and states
- All contexts (main, integration, pr-ready, fix-cascade, splitting, initialization)

**NO EXCEPTIONS**

## Related Rules
- R234: Mandatory State Traversal
- R231: Continuous Operation Through Transitions
- R288: Mandatory State File Updates

## Compliance Verification
```bash
# Check all state rule files have R405
for file in agent-states/*/*/rules.md; do
    if ! grep -q "R405" "$file"; then
        echo "VIOLATION: $file missing R405"
    fi
done

# Check agent outputs include flag with FULL FORMAT
if ! echo "$AGENT_OUTPUT" | grep -q "^CONTINUE-SOFTWARE-FACTORY=\(TRUE\|FALSE\)"; then
    echo "VIOLATION: Agent did not output continuation flag with full format"
fi

# Detect incorrect short format
if echo "$AGENT_OUTPUT" | grep -q "^R405.*: \(TRUE\|FALSE\)$"; then
    echo "VIOLATION: Agent used short format 'R405...: TRUE' instead of 'CONTINUE-SOFTWARE-FACTORY=TRUE'"
fi
```

## Date Added
2025-01-23

## Rationale
Full automation of Software Factory 2.0 requires a reliable, universal signal that external scripts can use to determine whether to continue processing or stop for human intervention. This flag provides that signal in a consistent, greppable format that works across all agents and states.

**CRITICAL DESIGN PRINCIPLE**: The system should run AUTONOMOUSLY through NORMAL operations. Human intervention should ONLY be required for TRULY EXCEPTIONAL situations that the system cannot handle. A failed effort, a needed split, or a review finding issues are ALL NORMAL and can be handled automatically. Default to TRUE unless you have a TRULY EXCEPTIONAL reason to stop.
## State Manager Coordination (SF 3.0)

State Manager integration with R405 continuation flag:
- **Shutdown consultation** returns `validation_result` with continuation recommendation
- **CONTINUE-SOFTWARE-FACTORY=TRUE** when state successfully updated and pushed
- **CONTINUE-SOFTWARE-FACTORY=FALSE** when validation fails or rollback occurs
- **Orchestrator** sets final flag based on State Manager shutdown result

The bookend pattern ensures continuation flag accuracy matches actual state persistence.

See: `.claude/agents/state-manager.md`, `agent-states/state-manager/SHUTDOWN_CONSULTATION/rules.md`
