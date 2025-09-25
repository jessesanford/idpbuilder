# Orchestrator - SPAWN_AGENTS State Rules (OPTIMIZED)

## State Entry Requirements
**Prerequisites**: Bootstrap rules already loaded (R203, R006, R319, R322, R288)

## SPAWN_AGENTS-Specific Rules (4 FILES)

### 1. 🚨🚨🚨 R151 - Parallel Spawning Requirements (BLOCKING - 50% GRADE)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R151-parallel-spawning-requirements.md`
**Critical**: Agents MUST be spawned with <5 second timestamp delta
**Implementation**: Use multiple Task tools in SINGLE message

### 2. 🚨🚨🚨 R208 - CD Before Spawn Protocol (SUPREME LAW #2)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R208-orchestrator-spawn-cd-protocol.md`
**Critical**: MUST CD to effort directory before EVERY spawn
**Violation**: Wrong directory = -100% immediate failure

### 3. ⚠️⚠️⚠️ R218 - Parallel Code Reviewer Spawning (WARNING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R218-orchestrator-parallel-code-reviewer-spawning.md`
**When**: Spawning multiple code reviewers for different efforts
**Requirement**: Parallel spawn for independent reviews

### 4. ⚠️⚠️⚠️ R322 Part A - Stop After Spawn (WARNING)
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R322 Part A-orchestrator-stop-after-spawn.md`
**Requirement**: After spawning, transition to MONITOR and STOP
**Purpose**: Let agents work before checking progress

## SPAWN_AGENTS State Actions

### Phase 1: Pre-Spawn Setup
```bash
# Verify all effort directories exist
for effort in efforts_in_progress; do
    if [ ! -d "efforts/phase$P/wave$W/$effort" ]; then
        echo "ERROR: Missing effort directory"
        exit 1
    fi
done
```

### Phase 2: Parallel Agent Spawning (R151 CRITICAL)
```markdown
MUST spawn all agents in SINGLE message:

Task: software-engineer
Working directory: efforts/phase1/wave1/effort1
[Instructions]

Task: software-engineer  
Working directory: efforts/phase1/wave1/effort2
[Instructions]

Task: code-reviewer
Working directory: efforts/phase1/wave1
[Instructions]
```

### Phase 3: Post-Spawn Protocol
```yaml
Next State: MONITOR
Stop Required: YES (R322 + R322 Part A)
Action: Update state file, commit, STOP
```

## Spawn Validation Checklist
- [ ] All agents spawned in ONE message (R151)
- [ ] Each spawn has CD to correct directory (R208)
- [ ] Timestamp delta <5 seconds average (R151)
- [ ] Independent work spawned in parallel (R218)
- [ ] State updated to MONITOR (R322 Part A)
- [ ] STOPPED after spawn (R322 Part A)

## Common Violations to Avoid
1. **Sequential Spawning**: Separate messages = FAIL
2. **Missing CD**: No working directory = FAIL
3. **Wrong Directory**: Not in effort dir = FAIL
4. **Continuing After Spawn**: Not stopping = FAIL

## Verification Marker
```bash
touch .state_rules_read_orchestrator_SPAWN_AGENTS
echo "$(date +%s) - SPAWN rules acknowledged" > .state_rules_read_orchestrator_SPAWN_AGENTS
```

## State Exit Protocol (R322 + R322 Part A)
1. Complete all spawns in single message
2. Update current_state to MONITOR
3. Commit orchestrator-state.json
4. STOP immediately (don't monitor yet)
5. Wait for /continue-orchestrating

---
*SPAWN_AGENTS State Rules - Optimized Version*
*Critical: R151 is 50% of orchestrator grade*