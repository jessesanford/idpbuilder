# ORCHESTRATOR RULE DISTRIBUTION MATRIX

## Overview
This matrix shows the current location of rules and their proposed new locations after optimization.

**Legend:**
- 🟢 Keep in bootstrap (orchestrator.md)
- 🔵 Move to specific states only
- 🟡 Duplicate (keep in bootstrap + specific states)
- 🔴 Remove from bootstrap entirely

## Supreme Laws Distribution

| Rule | Current Location | Proposed Location | Status | Reason |
|------|-----------------|-------------------|---------|---------|
| R234 - Mandatory State Traversal | orchestrator.md | SETUP_EFFORT_INFRASTRUCTURE, ANALYZE_* states | 🔵 | Only needed in those specific sequences |
| R221 - Bash Directory Reset | orchestrator.md | States with bash operations | 🔵 | Not all states use bash |
| R208 - CD Before Spawn | orchestrator.md | All SPAWN_* states (17 states) | 🔵 | Only spawning states need this |
| R235 - Pre-Flight Verification | orchestrator.md | States doing work operations | 🔵 | Each state can verify its own environment |
| SOFTWARE-FACTORY-STATE-MACHINE.md | orchestrator.md | Bootstrap (referenced by R206) | 🟢 | Needed to understand valid states |
| R280 - Main Branch Protection | orchestrator.md | Git operation states | 🔵 | Only git states need this |
| R290 - State Rule Reading | orchestrator.md | Bootstrap | 🟢 | Defines HOW to read state rules |
| R288 - State File Updates | orchestrator.md | Bootstrap | 🟢 | Every state needs to update state file |
| R281 - State File Initialization | orchestrator.md | INIT state only | 🔵 | Only INIT creates state file |
| R283 - Complete File Reading | orchestrator.md | Bootstrap | 🟢 | Needed to properly read all files |
| R307 - Branch Mergeability | orchestrator.md | Integration states | 🔵 | Only integration needs this |
| R308 - Incremental Branching | orchestrator.md | SETUP_EFFORT_INFRASTRUCTURE | 🔵 | Only infrastructure setup needs this |
| R309 - No Efforts in SF Repo | orchestrator.md | SETUP_*, SPAWN_* states | 🔵 | Only those states create/use efforts |
| R322 - Stop Before Transitions | orchestrator.md | Bootstrap + All states | 🟡 | Critical for stop model |

## Mandatory Rules Distribution

| Rule | Current Location | Proposed Location | Status | Reason |
|------|-----------------|-------------------|---------|---------|
| R006 - Never Writes Code | orchestrator.md | Work states (MONITOR_*, etc) | 🔵 | Only enforced during active work |
| R319 - Never Measures Code | orchestrator.md | MONITOR_* states | 🔵 | Only monitoring states would measure |
| R287 - TODO Persistence | orchestrator.md | Bootstrap | 🟢 | Needed for recovery in any state |
| R216 - Bash Syntax | orchestrator.md | States using bash | 🔵 | Reference rule for bash operations |
| R206 - State Validation | orchestrator.md | Bootstrap | 🟢 | Every state needs transition validation |
| R203 - State-Aware Startup | orchestrator.md | Bootstrap | 🟢 | Defines the startup sequence |

## State-Specific Rules Already in States

| State | Existing Unique Rules | Rules to Add | Rules to Remove |
|-------|----------------------|--------------|-----------------|
| INIT | R191, R192, R304 | R281 | None |
| SETUP_EFFORT_INFRASTRUCTURE | R176, R271, R191 | R234, R309, R308, R221 | None |
| ANALYZE_CODE_REVIEWER_PARALLELIZATION | Basic rules | R234, R053 | None |
| SPAWN_CODE_REVIEWERS_EFFORT_PLANNING | R051, R052, R197 | R208, R151, R221, R309 | None |
| WAITING_FOR_EFFORT_PLANS | Basic wait rules | None | None |
| ANALYZE_IMPLEMENTATION_PARALLELIZATION | R053 | R234 | None |
| SPAWN_AGENTS | R006, R151, R295, R052, R197, R255 | R208, R221, R309 | R006 (duplicate) |
| MONITOR_IMPLEMENTATION | Monitoring rules | R319, R006 | None |
| All other MONITOR_* states | Various monitoring | R319, R006 | None |
| INTEGRATION | Integration rules | R307, R280, R321 | None |
| PHASE_INTEGRATION | Phase rules | R307, R280 | None |

## Rules by Frequency of Use

| Rule | States Using It | Frequency | Location Decision |
|------|----------------|-----------|-------------------|
| R322 - Stop Before Transitions | ALL (67 states) | 100% | Bootstrap + States |
| R288 - State File Updates | ALL (67 states) | 100% | Bootstrap |
| R206 - State Validation | ALL (67 states) | 100% | Bootstrap |
| R287 - TODO Persistence | ALL (67 states) | 100% | Bootstrap |
| R203 - State-Aware Startup | ALL (67 states) | 100% | Bootstrap |
| R290 - State Rule Reading | ALL (67 states) | 100% | Bootstrap |
| R283 - Complete File Reading | ALL (67 states) | 100% | Bootstrap |
| R208 - CD Before Spawn | SPAWN_* (17 states) | 25% | Spawn states only |
| R151 - Parallel Timing | SPAWN_* (17 states) | 25% | Spawn states only |
| R309 - No SF Efforts | SETUP, SPAWN (20 states) | 30% | Those states only |
| R006 - Never Writes Code | Work states (40 states) | 60% | Work states only |
| R319 - Never Measures | MONITOR_* (7 states) | 10% | Monitor states only |
| R280 - Main Protection | Git states (15 states) | 22% | Git states only |
| R281 - State Init | INIT (1 state) | 1.5% | INIT only |

## Minimal Bootstrap Configuration

### Rules to Keep in orchestrator.md (7 rules + 1 reference)
```markdown
1. R283 - Complete File Reading (read files properly)
2. R290 - State Rule Reading (know to read state rules)
3. R203 - State-Aware Startup (startup sequence)
4. R206 - State Validation (validate transitions)
5. R288 - State File Updates (update state file)
6. R287 - TODO Persistence (recover from loss)
7. R322 - Stop Before Transitions (stop model)
8. SOFTWARE-FACTORY-STATE-MACHINE.md (state reference)
```

### Estimated Size Reduction
- Current orchestrator.md: 878 lines
- After migration: ~200 lines
- Reduction: 77%

## Migration Priority by Risk

### Priority 1: Low Risk Moves (Do First)
| Rule | Move To | Risk Level | Why Low Risk |
|------|---------|------------|--------------|
| R281 | INIT only | Low | Only used once |
| R319 | MONITOR_* only | Low | Clear scope |
| R309 | SETUP/SPAWN only | Low | Clear scope |
| R308 | SETUP_EFFORT only | Low | Single state |

### Priority 2: Medium Risk Moves
| Rule | Move To | Risk Level | Why Medium Risk |
|------|---------|------------|-----------------|
| R208 | SPAWN_* states | Medium | Many states affected |
| R151 | SPAWN_* states | Medium | Grading critical |
| R234 | Sequence states | Medium | Complex sequences |
| R280 | Git states | Medium | Multiple states |

### Priority 3: High Risk Moves (Do Last)
| Rule | Move To | Risk Level | Why High Risk |
|------|---------|------------|---------------|
| R221 | Bash states | High | Used frequently |
| R235 | Work states | High | Pre-flight critical |
| R006 | Work states | High | Core identity rule |
| R216 | Bash states | High | Syntax critical |

## State Categories and Rule Clusters

### Category 1: Initialization States
**States**: INIT
**Rule Cluster**:
- R191, R192 (repository config)
- R281 (state file creation)
- R304 (line counter)

### Category 2: Infrastructure States
**States**: SETUP_EFFORT_INFRASTRUCTURE, CREATE_NEXT_SPLIT_INFRASTRUCTURE
**Rule Cluster**:
- R234 (mandatory traversal)
- R309 (no SF efforts)
- R308 (incremental branching)
- R271 (single branch checkout)
- R176 (workspace isolation)

### Category 3: Spawning States
**States**: All SPAWN_* (17 states)
**Rule Cluster**:
- R208 (CD before spawn)
- R151 (parallel timing)
- R052 (spawn protocol)
- R197 (one agent per effort)
- R295 (spawn clarity)
- R309 (no SF efforts)
- R221 (bash reset)

### Category 4: Monitoring States
**States**: All MONITOR_* (7 states)
**Rule Cluster**:
- R319 (never measures code)
- R006 (never writes code)
- R008 (monitoring frequency)

### Category 5: Integration States
**States**: INTEGRATION, PHASE_INTEGRATION, PROJECT_INTEGRATION, FINAL_INTEGRATION
**Rule Cluster**:
- R307 (branch mergeability)
- R280 (main protection)
- R321 (immediate backport)
- R009 (integration branch creation)

### Category 6: Waiting States
**States**: All WAITING_* states
**Rule Cluster**:
- Minimal rules (mostly wait logic)
- R287 (TODO persistence)

### Category 7: Review States
**States**: WAVE_REVIEW, INTEGRATION_FEEDBACK_REVIEW
**Rule Cluster**:
- R031 (mandatory review)
- R057 (review authority)

## Validation Requirements

### For Bootstrap (orchestrator.md)
Must be able to:
- [ ] Determine current state from orchestrator-state.json
- [ ] Know to read state-specific rules (R290)
- [ ] Validate state transitions (R206)
- [ ] Update and persist state (R288)
- [ ] Recover from context loss (R287)
- [ ] Stop before transitions (R322)

### For Each State Directory
Must contain:
- [ ] All rules needed for that state's work
- [ ] Clear prerequisites from bootstrap
- [ ] Valid transition targets
- [ ] State-specific verification markers (R290)

## Testing Matrix

| Test Scenario | States to Test | Rules to Verify | Expected Result |
|---------------|----------------|-----------------|-----------------|
| Cold start | INIT | R203, R290, R281 | Loads state, creates file |
| Spawn operation | SPAWN_AGENTS | R208, R151 | Proper CD, timing |
| Monitoring | MONITOR_IMPLEMENTATION | R319, R006 | No code operations |
| Integration | INTEGRATION | R307, R280 | Branch protection |
| State transition | Any -> Any | R322, R288, R206 | Stop, update, validate |
| Recovery | Any | R287, R290 | TODO recovery, rule reload |

## Summary

### Total Rules to Migrate: 14 Supreme Laws + 6 Mandatory Rules = 20 rules
### Rules Staying in Bootstrap: 7 rules + 1 reference
### Rules Moving to States: 13 rules
### States Affected: All 67 states will be updated
### Expected Reduction: 77% fewer lines in orchestrator.md
### Risk Level: Medium (with proper testing)

This distribution ensures each state has exactly what it needs while keeping the bootstrap minimal but functional for state determination and rule loading.