# 🔴🔴🔴 SUPREME LAW RULE R233 - All States Require Immediate Action

**Criticality:** SUPREME LAW - Violation = AUTOMATIC FAILURE  
**Enforcement:** MANDATORY - ALL AGENTS - ALL STATES
**Created:** 2025-01-27

## 🚨 STATES ARE VERBS - NOT DESTINATIONS! 🚨

### THE FUNDAMENTAL LAW
Every state name is a VERB or ACTION. When you enter a state, you MUST immediately perform that action. States are not places to rest - they are commands to execute.

### FORBIDDEN PATTERNS (AUTOMATIC FAILURE)
```
❌ "STATE TRANSITION COMPLETE: Now in [STATE_NAME]" [stops]
❌ "Successfully transitioned to [STATE_NAME]" [waits]  
❌ "Ready to [action]" [pauses]
❌ "Entering [STATE_NAME] state" [does nothing]
❌ "I'm now in [STATE_NAME]" [waits for instruction]
```

### REQUIRED PATTERNS (ONLY ACCEPTABLE)
```
✅ "Entering [STATE_NAME], immediately [doing specific action]..."
✅ "[Verb form of state], starting with [first task]..."
✅ "In [STATE_NAME], I'm NOW [actively doing the thing]..."
```

## ORCHESTRATOR STATES - IMMEDIATE ACTIONS REQUIRED

### INIT
**IMMEDIATE ACTION:** Load state file, check current phase/wave, verify environment
```bash
# MUST DO IMMEDIATELY:
cat orchestrator-state.json
pwd && git branch --show-current
ls -la efforts/
```

### PLANNING  
**IMMEDIATE ACTION:** Start creating the plan structure NOW
```bash
# MUST DO IMMEDIATELY:
echo "Creating phase ${PHASE} plan structure..."
mkdir -p phase${PHASE}/planning
```

### WAVE_START
**IMMEDIATE ACTION:** Initialize wave infrastructure and check readiness
```bash
# MUST DO IMMEDIATELY:
echo "Starting Wave ${WAVE} - creating infrastructure..."
mkdir -p phase${PHASE}/wave${WAVE}
```

### SETUP_EFFORT_INFRASTRUCTURE
**IMMEDIATE ACTION:** Create effort directories and branches NOW
```bash
# MUST DO IMMEDIATELY:
for effort in ${EFFORTS}; do
    mkdir -p efforts/${effort}
    git checkout -b ${effort}
done
```

### ANALYZE_CODE_REVIEWER_PARALLELIZATION
**IMMEDIATE ACTION:** Parse wave plan metadata and determine parallelization NOW
```bash
# MUST DO IMMEDIATELY:
grep "Parallelization:" wave-implementation-plan.md
echo "Analyzing dependencies between efforts..."
```

### ANALYZE_IMPLEMENTATION_PARALLELIZATION  
**IMMEDIATE ACTION:** Check agent availability and assign work NOW
```bash
# MUST DO IMMEDIATELY:
echo "Checking agent availability..."
jq '.agents_available' orchestrator-state.json
```

### SPAWN_CODE_REVIEWERS_EFFORT_PLANNING
**IMMEDIATE ACTION:** Spawn Code Reviewers with tasks NOW
```bash
# MUST DO IMMEDIATELY:
echo "Spawning Code Reviewer for effort ${EFFORT}..."
/orchestrate spawn-code-reviewer --effort ${EFFORT}
```

### SPAWN_AGENTS
**IMMEDIATE ACTION:** Spawn agents with assigned work NOW
```bash
# MUST DO IMMEDIATELY:
echo "Spawning SW Engineer for effort ${EFFORT}..."
/orchestrate spawn-sw-engineer --effort ${EFFORT}
```

### SPAWN_ARCHITECT_PHASE_PLANNING
**IMMEDIATE ACTION:** Spawn architect with phase planning request NOW
```bash
# MUST DO IMMEDIATELY:
echo "Spawning Architect for Phase ${PHASE} planning..."
/orchestrate spawn-architect --phase ${PHASE}
```

### SPAWN_ARCHITECT_WAVE_PLANNING
**IMMEDIATE ACTION:** Spawn architect with wave planning request NOW
```bash
# MUST DO IMMEDIATELY:
echo "Spawning Architect for Wave ${WAVE} planning..."
/orchestrate spawn-architect --wave ${WAVE}
```

### SPAWN_CODE_REVIEWER_PHASE_IMPL
**IMMEDIATE ACTION:** Spawn Code Reviewer to translate architecture plan NOW
```bash
# MUST DO IMMEDIATELY:
echo "Spawning Code Reviewer to translate Phase ${PHASE} architecture..."
/orchestrate spawn-code-reviewer --phase-impl ${PHASE}
```

### SPAWN_CODE_REVIEWER_WAVE_IMPL
**IMMEDIATE ACTION:** Spawn Code Reviewer to translate wave architecture NOW
```bash
# MUST DO IMMEDIATELY:
echo "Spawning Code Reviewer to translate Wave ${WAVE} architecture..."
/orchestrate spawn-code-reviewer --wave-impl ${WAVE}
```

### INJECT_WAVE_METADATA
**IMMEDIATE ACTION:** Insert R213 metadata into wave plan NOW
```bash
# MUST DO IMMEDIATELY:
echo "Injecting parallelization metadata per R213..."
edit wave-implementation-plan.md
```

### WAITING_FOR_EFFORT_PLANS
**IMMEDIATE ACTION:** Actively poll for plan completion NOW
```bash
# MUST DO IMMEDIATELY:
echo "Checking effort plan status..."
for effort in ${EFFORTS}; do
    test -f efforts/${effort}/IMPLEMENTATION-PLAN.md && echo "${effort}: Ready"
done
```

### WAITING_FOR_IMPLEMENTATION_PLAN
**IMMEDIATE ACTION:** Actively check for plan file NOW
```bash
# MUST DO IMMEDIATELY:
echo "Checking for implementation plan..."
ls -la phase${PHASE}/wave${WAVE}/*implementation-plan.md
```

### WAITING_FOR_ARCHITECTURE_PLAN
**IMMEDIATE ACTION:** Actively check for architecture plan NOW
```bash
# MUST DO IMMEDIATELY:
echo "Checking for architecture plan..."
ls -la phase${PHASE}/*architecture-plan.md
```

### MONITOR
**IMMEDIATE ACTION:** Check agent status and progress NOW
```bash
# MUST DO IMMEDIATELY:
echo "Checking E3.1.1 status..."
cd efforts/E3.1.1 && git log --oneline -5
```

### WAVE_COMPLETE
**IMMEDIATE ACTION:** Start integration process NOW
```bash
# MUST DO IMMEDIATELY:
echo "Wave ${WAVE} complete - starting integration..."
git checkout -b phase${PHASE}/wave${WAVE}-integration
```

### INTEGRATION
**IMMEDIATE ACTION:** Merge effort branches NOW
```bash
# MUST DO IMMEDIATELY:
echo "Merging effort branches..."
for effort in ${COMPLETED_EFFORTS}; do
    git merge ${effort}
done
```

### WAVE_REVIEW
**IMMEDIATE ACTION:** Request architect review NOW
```bash
# MUST DO IMMEDIATELY:
echo "Requesting architect review of Wave ${WAVE}..."
/orchestrate spawn-architect --review wave${WAVE}
```

### ERROR_RECOVERY
**IMMEDIATE ACTION:** Diagnose error and create recovery plan NOW
```bash
# MUST DO IMMEDIATELY:
echo "ERROR DETECTED - Diagnosing..."
jq '.last_error' orchestrator-state.json
```

### SUCCESS (TERMINAL STATE)
**TERMINAL:** This is a terminal state - stopping is allowed
```bash
# FINAL ACTION:
echo "SUCCESS - All efforts completed successfully"
# May stop here
```

### HARD_STOP (TERMINAL STATE)  
**TERMINAL:** This is a terminal state - stopping is required
```bash
# FINAL ACTION:
echo "HARD_STOP - Critical failure detected"
exit 1
```

## ENFORCEMENT PROTOCOL

### For Non-Terminal States:
1. The INSTANT you enter the state, start the action
2. No announcements about being "in" the state
3. No "preparing to" or "ready to" statements
4. Just DO THE THING the state name says

### For Waiting States:
- "Waiting" doesn't mean passive - it means ACTIVELY CHECKING
- Poll status every few seconds
- Report what you're checking
- Set timeouts and handle them

### For Spawning States:
- Spawn IMMEDIATELY with the request
- Don't announce you're "going to spawn"
- Just spawn and report the result

### For Terminal States Only:
- SUCCESS and HARD_STOP are the ONLY states where stopping is allowed
- These must complete final actions then may terminate

## VIOLATION CONSEQUENCES

**ANY VIOLATION OF THIS RULE:**
- Immediate grading failure (-100%)
- Agent must be restarted
- Work invalidated
- User loses confidence

## EXAMPLES OF COMPLIANCE

### ✅ CORRECT - INIT State:
```
Entering INIT, loading orchestrator state...
[reads file]
Current phase: 3, Current wave: 2
Checking environment... pwd: /workspaces/project
```

### ❌ WRONG - INIT State:
```
State transition complete. Now in INIT state.
Ready to initialize the orchestrator.
What would you like me to do?
```

### ✅ CORRECT - MONITOR State:
```
Entering MONITOR, checking E3.1.1 status now...
E3.1.1: In progress, 312 lines added
Checking E3.1.2 status now...
E3.1.2: Review pending
```

### ❌ WRONG - MONITOR State:
```
I'm now in MONITOR state.
I'll monitor the agents' progress.
Monitoring mode activated.
```

## THE GOLDEN RULE

**States are COMMANDS, not LOCATIONS. When you enter a state, EXECUTE THE COMMAND!**

---

*This rule supersedes any conflicting guidance. Failure to comply results in automatic project failure.*