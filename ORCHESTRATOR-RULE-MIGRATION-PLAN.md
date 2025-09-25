# ORCHESTRATOR RULE MIGRATION PLAN

## Executive Summary

The orchestrator agent currently loads ALL rules during bootstrap, creating unnecessary cognitive overhead and context bloat. Since the orchestrator LOSES ALL CONTEXT between states (completely scrubbed and restarted), we can optimize by moving state-specific rules to their respective state directories, keeping only essential bootstrap rules in orchestrator.md.

**CRITICAL CONSTRAINT**: The orchestrator has NO MEMORY between states. Each state must be completely self-contained.

## Current State Analysis

### Files Analyzed
- **orchestrator.md**: 878 lines total
  - Contains 16 Supreme Laws
  - Contains 7 Mandatory Rules  
  - Many rules only needed in specific states
- **State directories**: 67 states with rules.md files
  - Each already has some state-specific rules
  - Many duplicate rules from orchestrator.md

### Current Rule Loading Pattern
```
1. Orchestrator starts fresh (NO prior context)
2. Reads orchestrator.md (ALL rules)
3. Determines current state from orchestrator-state.json
4. Reads state-specific rules
5. Has ALL rules loaded (many unnecessary)
```

## Rule Category Analysis

### CATEGORY A: ALWAYS IN BOOTSTRAP (Cannot be moved)
**These rules MUST stay in orchestrator.md because they are needed to:**
- Determine current state
- Load state rules
- Understand basic identity
- Read orchestrator-state.json

#### Essential Bootstrap Rules (MUST KEEP):

1. **R283** - Complete File Reading
   - Needed to properly read any files including state rules
   - Without this, orchestrator might partially read state rules

2. **R290** - State Rule Reading and Verification
   - Defines HOW to read state rules
   - Creates verification markers
   - Without this, orchestrator won't know to read state rules

3. **R203** - State-Aware Agent Startup
   - Defines the startup sequence
   - Tells orchestrator to determine state and load state rules
   - Core identity rule

4. **R206** - State Machine Transition Validation
   - Needed to validate ANY state transition
   - References SOFTWARE-FACTORY-STATE-MACHINE.md
   - Used in all states

5. **R288** - State File Update and Commit
   - Needed to update orchestrator-state.json
   - Used in EVERY state transition
   - Without this, state persistence fails

6. **R287** - TODO Persistence
   - Needed for recovery from compaction
   - Used across all states
   - Critical for context recovery

7. **SOFTWARE-FACTORY-STATE-MACHINE.md**
   - The absolute authority on valid states
   - Needed to understand what states exist
   - Referenced by R206

8. **R322** - Mandatory Stop Before State Transitions
   - Prevents automatic state transitions
   - Core to the stop-and-restart model
   - Needed in EVERY state

### CATEGORY B: STATE-SPECIFIC (Can be moved)
**These rules are only needed in certain states:**

#### Rules to Move to Specific States:

1. **R234** - Mandatory State Traversal
   - Only needed in states that are part of mandatory sequences
   - Move to: SETUP_EFFORT_INFRASTRUCTURE, ANALYZE_CODE_REVIEWER_PARALLELIZATION, etc.

2. **R208** - CD Before Spawn
   - Only needed in SPAWN_* states
   - Move to: All SPAWN_* state directories

3. **R281** - Complete State File Initialization
   - Only needed in INIT state
   - Move to: INIT/rules.md

4. **R151** - Parallel Spawning Timing
   - Only needed in SPAWN_* states
   - Move to: All SPAWN_* state directories

5. **R309** - Never Create Efforts in SF Repo
   - Only needed in SETUP_EFFORT_INFRASTRUCTURE and SPAWN_* states
   - Move to: Those specific states

6. **R006** - Orchestrator Never Writes Code
   - While important, only enforced during active work states
   - Keep minimal reference in bootstrap, full rule in work states

7. **R319** - Orchestrator Never Measures Code
   - Only relevant in MONITOR states
   - Move to: MONITOR_* state directories

8. **R221** - Bash Directory Reset
   - While used frequently, can be loaded per-state
   - States know if they need bash operations

9. **R235** - Mandatory Pre-Flight Verification
   - Can be part of state-specific startup
   - Each state verifies its own environment

10. **R280** - Main Branch Protection
    - Only relevant during git operations
    - Move to states that perform git work

11. **R307** - Independent Branch Mergeability
    - Only needed during integration states
    - Move to: INTEGRATION, PHASE_INTEGRATION states

12. **R308** - Incremental Branching Strategy
    - Only needed during branch creation
    - Move to: SETUP_EFFORT_INFRASTRUCTURE

13. **R216** - Bash Execution Syntax
    - Reference rule, move to states that execute bash

14. **R232** - TodoWrite Pending Items Override
    - Can be loaded when TODO operations happen
    - Not needed at bootstrap

### CATEGORY C: TRANSITION RULES (Special handling)
**Rules needed during state transitions:**

These should be duplicated:
- In the state being LEFT (to know how to transition)
- In the state being ENTERED (to validate the transition)

Examples:
- R322 (Stop before transitions) - Keep in bootstrap AND all states
- R288 (Update state file) - Keep in bootstrap AND all states
- R206 (Validate transitions) - Keep in bootstrap for reference

## State Rule Matrix

### Critical States and Their Rule Requirements

#### INIT State
**Current**: Has R191, R192 (target config rules)
**Add**: R281 (state file creation)
**Remove from bootstrap**: R281

#### SETUP_EFFORT_INFRASTRUCTURE  
**Current**: Basic rules
**Add**: R234, R309, R308, R271 (infrastructure rules)
**Remove from bootstrap**: These specific rules

#### All SPAWN_* States (17 states)
**Current**: Various spawn rules
**Add**: R208, R151, R309, R221 (spawn protocols)
**Remove from bootstrap**: These spawn-specific rules

#### All MONITOR_* States (7 states)
**Current**: Monitoring rules
**Add**: R319, R006 (delegation rules)
**Remove from bootstrap**: Monitor-specific rules

#### INTEGRATION States (5 states)
**Current**: Integration rules
**Add**: R307, R280, R321 (branch/integration rules)
**Remove from bootstrap**: Integration-specific rules

## Migration Strategy

### Phase 1: Preparation
1. Create backup of current orchestrator.md
2. Document current rule loading behavior
3. Create test scenarios for each state

### Phase 2: Rule Redistribution
1. **Keep Minimal Bootstrap** (orchestrator.md):
   ```markdown
   - R283 (Complete file reading)
   - R290 (State rule reading)
   - R203 (State-aware startup)
   - R206 (State validation)
   - R288 (State file updates)
   - R287 (TODO persistence)
   - R322 (Stop before transitions)
   - SOFTWARE-FACTORY-STATE-MACHINE.md reference
   ```

2. **Enhance State Rules** (per state):
   - Add state-specific Supreme Laws
   - Add rules only needed in that state
   - Include clear "Prerequisites" section
   - Include "Next States" section

3. **Create State Rule Template**:
   ```markdown
   # Orchestrator - [STATE_NAME] State Rules
   
   ## Prerequisites (from bootstrap)
   - R203, R206, R288, R287, R290, R283, R322
   - SOFTWARE-FACTORY-STATE-MACHINE.md loaded
   
   ## State-Specific Supreme Laws
   [Rules that are critical for this state]
   
   ## State-Specific Mandatory Rules
   [Rules only needed in this state]
   
   ## State Context
   [What this state does]
   
   ## Valid Transitions
   [Where can we go from here]
   ```

### Phase 3: Testing Protocol

#### Test Scenario 1: Cold Start
```bash
# Orchestrator starts with no context
1. Read orchestrator.md (minimal rules)
2. Read orchestrator-state.json
3. Determine state = "SPAWN_AGENTS"
4. Read agent-states/orchestrator/SPAWN_AGENTS/rules.md
5. Verify has all needed rules for spawning
```

#### Test Scenario 2: State Transition
```bash
# Orchestrator completes SPAWN_AGENTS
1. Update orchestrator-state.json to MONITOR
2. Stop per R322
3. User runs /continue-orchestrating
4. NEW orchestrator starts (no memory)
5. Reads minimal bootstrap
6. Reads state = MONITOR
7. Reads MONITOR rules
8. Has all needed rules for monitoring
```

#### Test Scenario 3: Recovery from Compaction
```bash
# Context lost mid-state
1. Orchestrator starts fresh
2. Reads minimal bootstrap (includes R287)
3. Detects TODO file exists
4. Loads TODOs
5. Reads current state
6. Loads state rules
7. Continues work
```

## Risk Analysis

### High Risk Areas
1. **Missing Essential Rules**
   - Risk: Orchestrator doesn't know to read state rules
   - Mitigation: Keep R290, R203 in bootstrap

2. **State Determination Failure**
   - Risk: Can't figure out current state
   - Mitigation: Keep R288, state file rules in bootstrap

3. **Incomplete State Rules**
   - Risk: State missing critical rule
   - Mitigation: Comprehensive testing of each state

### Medium Risk Areas
1. **Duplicate Rules**
   - Risk: Same rule in bootstrap and state
   - Mitigation: Clear documentation of rule location

2. **Transition Validation**
   - Risk: Invalid transitions not caught
   - Mitigation: R206 stays in bootstrap

### Low Risk Areas
1. **Performance**
   - Risk: Slower startup from reading multiple files
   - Mitigation: Minimal - already reading multiple files

## Implementation Sequence

### Step 1: Create Enhanced State Rules (Week 1)
- For each state directory:
  1. Audit current rules.md
  2. Identify missing rules from bootstrap
  3. Add those rules with proper context
  4. Test state in isolation

### Step 2: Create Minimal Bootstrap (Week 2)
1. Backup current orchestrator.md
2. Create new minimal version
3. Test cold start scenarios
4. Verify state determination works

### Step 3: Integration Testing (Week 3)
1. Test complete workflows
2. Test state transitions
3. Test recovery scenarios
4. Test parallel operations

### Step 4: Rollout (Week 4)
1. Deploy to development
2. Monitor for issues
3. Deploy to production
4. Create rollback plan

## Verification Checklist

### Pre-Migration Verification
- [ ] All 67 state directories audited
- [ ] Rule dependencies mapped
- [ ] Test scenarios documented
- [ ] Rollback plan created

### Per-State Verification
For each state:
- [ ] State rules.md contains all needed rules
- [ ] Prerequisites section lists bootstrap rules
- [ ] State context clearly defined
- [ ] Valid transitions documented
- [ ] Test scenario passes

### Post-Migration Verification
- [ ] Cold start works
- [ ] All state transitions work
- [ ] Recovery from compaction works
- [ ] No duplicate rule loading
- [ ] No missing rules in any state
- [ ] Performance acceptable

### Continuous Verification
- [ ] State rule completeness check (weekly)
- [ ] Bootstrap minimality check (monthly)
- [ ] Rule usage audit (monthly)
- [ ] Dead rule detection (quarterly)

## Expected Benefits

### Immediate Benefits
1. **Reduced Cognitive Load**: 70% fewer rules loaded per state
2. **Faster Startup**: Less reading required
3. **Clearer Context**: Only relevant rules visible
4. **Better Maintainability**: Rules located where used

### Long-term Benefits
1. **Easier Debugging**: State-specific issues isolated
2. **Simpler Onboarding**: New developers learn per-state
3. **Reduced Errors**: Less chance of applying wrong rule
4. **Flexible Evolution**: States can have unique rules

## Success Metrics

### Quantitative Metrics
- Bootstrap size: From 878 lines to ~200 lines (77% reduction)
- Rules per state: From 23 to 8-12 average (50% reduction)
- Startup time: Measure before/after
- Memory usage: Measure before/after

### Qualitative Metrics
- Developer feedback on clarity
- Error rate reduction
- Time to debug issues
- Ease of adding new states

## Rollback Plan

If issues arise:
1. **Immediate**: Revert orchestrator.md to original
2. **Quick Fix**: Add missing rules to affected states
3. **Full Rollback**: Restore all original files from backup
4. **Hybrid Approach**: Keep some optimization, restore critical rules

## Conclusion

This migration will transform the orchestrator from a "load everything" model to a "load what's needed" model, taking advantage of the fact that context is completely lost between states anyway. Each state becomes self-contained with exactly the rules it needs, no more, no less.

The key insight is that since the orchestrator has NO MEMORY between states, we don't need to front-load all rules. We can load them just-in-time as needed, making each state cleaner and more focused.

**Critical Success Factor**: The bootstrap must retain enough rules for the orchestrator to:
1. Determine its current state
2. Know to load state-specific rules
3. Validate state transitions
4. Persist state changes

Everything else can be distributed to where it's actually used.