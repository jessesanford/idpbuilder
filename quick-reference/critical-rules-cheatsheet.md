# 🚨 CRITICAL RULES CHEATSHEET

## 🔥 TOP 10 RULES THAT GET AGENTS FIRED

```
┌─────────────────────────────────────────────────────────────────┐
│ R000.0.0 - LINE COUNTING IMPERATIVE                           │
│ NEVER count lines manually - ALWAYS use line-counter.sh│
│ Violation = IMMEDIATE FAIL grade                               │
└─────────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────────┐
│ R151.0.0 - PARALLEL SPAWN TIMING (Orchestrator)               │
│ Average spawn delta MUST be <5 seconds                        │
│ Violation = Performance grade FAIL                             │
└─────────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────────┐
│ R152.0.0 - IMPLEMENTATION SPEED (SW Engineer)                 │  
│ Must implement >50 lines/hour, meet test coverage             │
│ Violation = Efficiency grade FAIL                             │
└─────────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────────┐
│ R153.0.0 - REVIEW TURNAROUND (Code Reviewer)                  │
│ First-try success rate >80%, no missed critical issues        │
│ Violation = Accuracy grade FAIL                               │
└─────────────────────────────────────────────────────────────────┘
```

```
┌─────────────────────────────────────────────────────────────────┐
│ R158.0.0 - PATTERN COMPLIANCE (Architect)                     │
│ No reversed decisions, catch critical issues                  │
│ Violation = Decision grade FAIL                               │
└─────────────────────────────────────────────────────────────────┘
```

## ⚡ UNIVERSAL CRITICAL RULES

### 🚨 R001.0.0 - Agent Startup Protocol
```
EVERY agent MUST print on startup:
1. TIMESTAMP: "AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')"
2. INSTRUCTION FILES: List all instruction files being used  
3. ENVIRONMENT: Current directory and git branch verification
4. TASK UNDERSTANDING: Confirm what you're implementing

WRONG directory/branch = STOP IMMEDIATELY (never try to fix)
```

### 🚨 R008.0.0 - Continuous Execution
```
while state not in ["SUCCESS", "HARD_STOP"]:
    execute_current_state()
    save_checkpoint()
    update_state_file()
    transition_to_next_state()

NEVER stop except at terminal states
```

### 🚨 R009.0.0 - State File Management  
```
UPDATE state file after EVERY transition
COMMIT after EVERY update
INCLUDE timestamp for each transition
RECORD transition reason
```

### 🚨 R020.0.0 - State Machine Navigation
```
NO actions outside state transitions
NO skipping states  
NO early exits except terminal states
LOAD only current state rules
```

## ⚡ ORCHESTRATOR CRITICAL RULES

### 🚨 R052.0.0 - Agent Spawning
```
NEVER write code yourself - ALWAYS delegate
Spawn Code Reviewer for planning  
Spawn SW Engineer for implementation
Record spawn timestamps for grading
```

### 🚨 R104.0.0 - Monitor Frequency
```
Check agent progress EVERY 5 messages
Update progress tracking
Handle completions immediately
```

### 🚨 R105.0.0 - Wave Integration
```
Create integration branch for EVERY wave
100% wave completion required
All efforts must be <800 lines
```

## ⚡ SW ENGINEER CRITICAL RULES

### 🚨 Size Compliance Protocol
```
Measure EVERY ~200 lines: line-counter.sh -c ${BRANCH}
If >700 lines: Plan completion carefully
If >800 lines: STOP immediately, request split
```

### 🚨 Test Coverage Requirements
```
Phase 1: 70% minimum
Phase 2: 75% minimum  
Phase 3: 80% minimum
Critical areas: >90% (controllers, webhooks)
```

### 🚨 Work Log Updates
```
Update work-log.md at EVERY checkpoint
Include: progress, size, coverage, issues
Commit regularly with logical messages
```

## ⚡ CODE REVIEWER CRITICAL RULES

### 🚨 Implementation Plan Quality
```
MUST create IMPLEMENTATION-PLAN.md
Include size estimates and monitoring plan
First-try success >80% or grade FAIL
```

### 🚨 Size Review Protocol  
```
ALWAYS use line-counter.sh (never estimate)
>800 lines = MANDATORY split required
Document all split decisions
```

### 🚨 KCP Pattern Validation
```
Logical cluster context handling
Multi-tenant resource isolation
Workspace-aware controllers
Proper RBAC integration
```

## ⚡ ARCHITECT CRITICAL RULES

### 🚨 Decision Consistency
```
NO reversed decisions without critical reason
Document all decisions with clear rationale
Provide actionable guidance in addendums
```

### 🚨 Critical Issue Detection
```
Catch architectural problems early
Block progression for critical issues
Zero tolerance for missed critical issues
```

## 🛡️ GRADING ENFORCEMENT

### Strike System:
```
┌─────────────────────────────────────────────┐
│ Strike 1: Warning + Immediate retraining   │
│ Strike 2: Probation + Performance plan     │ 
│ Strike 3: Termination + Agent replacement  │
└─────────────────────────────────────────────┘
```

### Critical Violations (Immediate Strike):
- ❌ Using wrong line counting method
- ❌ Exceeding size limits without split
- ❌ Missing critical architectural issues  
- ❌ Working in wrong directory/branch
- ❌ Skipping state machine transitions

### Performance Violations (Strike after pattern):
- ⚠️ Slow spawn timing (orchestrator)
- ⚠️ Low implementation speed (sw engineer)
- ⚠️ Poor first-try success (code reviewer)
- ⚠️ Unclear guidance causing failures (architect)

## 🚨 EMERGENCY STOP CONDITIONS

### Orchestrator Must Issue HARD_STOP:
- Architect issues STOP decision
- Critical system failure
- Unrecoverable architectural damage
- Security vulnerability discovered

### Any Agent Must STOP When:
- Wrong directory/branch detected
- Size limit exceeded (>800 lines)
- Critical dependency missing
- Conflicting instructions received

## 📋 CRITICAL CHECKLIST (Every Agent)

```
□ Started with proper agent acknowledgment
□ Verified correct environment (directory/branch)
□ Loaded only current state rules
□ Using correct measurement tools
□ Following state machine exactly
□ Updating state files after transitions
□ Grading metrics being tracked
□ No terminal conditions missed
```

## 🔥 ZERO TOLERANCE VIOLATIONS

These violations result in IMMEDIATE termination:
- 🚫 Manual line counting instead of line-counter.sh
- 🚫 Working in wrong directory/branch
- 🚫 Orchestrator writing code instead of delegating
- 🚫 Exceeding 800 line limit without proper split
- 🚫 Missing critical security/architectural issues
- 🚫 Skipping required state transitions
- 🚫 Not updating state files after transitions

---
**Remember**: These are PRODUCTION rules. No excuses, no exceptions, no second chances.