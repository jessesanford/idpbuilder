# R054 Placement Analysis Report

## Executive Summary
R054 is currently required reading in orchestrator.md startup, but based on analysis, it should be moved to state-specific rules where it's actually needed.

## 1. R054 Rule Summary

### What R054 Says:
**Title:** Implementation Plan Creation
**Criticality:** BLOCKING
**Core Requirement:** Code Reviewers MUST create comprehensive IMPLEMENTATION-PLAN.md files for each effort before SW Engineers begin work.

### Key Requirements:
- Code Reviewers (not orchestrators) create implementation plans
- Plans must include technical architecture, implementation sequence, size management
- Plans must be created BEFORE SW Engineers start work
- Plans are created in effort directories during EFFORT_PLANNING state

### Who It Applies To:
- **Primary:** Code Reviewers (they CREATE the plans)
- **Secondary:** SW Engineers (they FOLLOW the plans)
- **Tertiary:** Orchestrator (needs to know plans exist but doesn't create them)

## 2. Current Placement Analysis

### Option A: In orchestrator.md startup rules (CURRENT)
**Current Location:** Line 365 in orchestrator.md startup sequence

**Pros:**
- Orchestrator is aware of R054 from the beginning
- Can reference it when planning wave structure

**Cons:**
- Adds unnecessary startup overhead (18th rule to read!)
- Not needed in many orchestrator states (INIT, PLANNING, ERROR_RECOVERY, etc.)
- Orchestrator NEVER creates implementation plans itself
- Violates just-in-time principle

**Verdict:** INEFFICIENT - Rule is loaded but rarely needed

## 3. Alternative Placement Options

### Option B: In SPAWN_AGENTS state rules
**Analysis:** TOO LATE
- By SPAWN_AGENTS, implementation plans should already exist
- SW Engineers need the plans to be ready when they spawn
- Would miss the critical planning phase

### Option C: In SPAWN_CODE_REVIEWERS_EFFORT_PLANNING state rules (RECOMMENDED)
**Analysis:** OPTIMAL PLACEMENT

**Why This Is Best:**
1. **Just-In-Time:** Read exactly when spawning code reviewers for planning
2. **Context-Appropriate:** Orchestrator needs to know about R054 when:
   - Spawning code reviewers for effort planning
   - Telling them what to create
   - Setting expectations for deliverables
3. **Reduced Overhead:** Not loaded during unrelated states
4. **Clear Purpose:** Direct connection between reading rule and using it

### Option D: In PLANNING state rules
**Analysis:** TOO EARLY
- PLANNING state is about high-level wave planning
- Effort-level implementation plans come later
- Would confuse orchestrator's role (it doesn't create these plans)

## 4. Current State Rule References

### Where R054 Is Already Referenced:

1. **code-reviewer/EFFORT_PLANNING/rules.md**
   - Line 72: INFO level reference
   - Appropriate: Code Reviewers need this when creating plans

2. **sw-engineer/IMPLEMENTATION/rules.md**
   - Line 331: MANDATORY level reference
   - Appropriate: SW Engineers need to follow the plans

3. **orchestrator/SPAWN_AGENTS/rules.md**
   - Line 218: Listed in rules to acknowledge
   - QUESTIONABLE: Too late, plans should already exist

4. **orchestrator/SPAWN_CODE_REVIEWERS_FOR_REVIEW/rules.md**
   - Line 213: Reference for report generation
   - DIFFERENT CONTEXT: This is about review reports, not implementation plans

## 5. Software Factory Philosophy

### Just-In-Time Knowledge Loading
The factory follows these principles:
1. **State-Aware Loading:** Rules loaded when entering relevant states
2. **Reduced Cognitive Load:** Don't front-load unnecessary rules
3. **Context-Specific:** Rules read when they're about to be used
4. **Efficiency:** Minimize startup overhead for faster agent initialization

### Evidence from State Machine:
- R203 (State-Aware Startup) explicitly supports state-specific rule loading
- State directories exist specifically for just-in-time rule delivery
- Each state has its own rules.md for context-specific requirements

## 6. Recommendation

### PRIMARY RECOMMENDATION: Move R054 to State-Specific Loading

**Remove From:**
- orchestrator.md startup sequence (line 365)

**Add To:**
- `/agent-states/orchestrator/SPAWN_CODE_REVIEWERS_EFFORT_PLANNING/rules.md`
  - Add as MANDATORY reading in PRIMARY DIRECTIVES section
  - Include clear explanation of WHY orchestrator needs to know this

**Implementation:**
```markdown
## PRIMARY DIRECTIVES - MANDATORY READING

### Implementation Plan Requirements
7. READ: $CLAUDE_PROJECT_DIR/rule-library/R054-implementation-plan-creation.md
   **R054** - Understanding what Code Reviewers must create
   **Why:** You're spawning reviewers to create these plans - know the deliverable
```

### SECONDARY ADJUSTMENTS:

1. **Keep in code-reviewer/EFFORT_PLANNING:**
   - Already there, elevate from INFO to BLOCKING criticality

2. **Keep in sw-engineer/IMPLEMENTATION:**
   - Already there with correct MANDATORY level

3. **Remove from orchestrator/SPAWN_AGENTS:**
   - Too late to be useful here

## 7. Benefits of This Change

1. **Faster Orchestrator Startup:**
   - One less rule to read at initialization
   - Reduces startup from ~18 rules to ~17

2. **Better Context:**
   - Rule read exactly when needed
   - Clear connection to spawning code reviewers

3. **Clearer Responsibilities:**
   - Orchestrator knows it's NOT creating plans
   - Just needs to understand what reviewers will create

4. **Follows Factory Philosophy:**
   - Aligns with R203 state-aware approach
   - Implements just-in-time knowledge loading

## 8. Implementation Steps

1. Edit `/agent-states/orchestrator/SPAWN_CODE_REVIEWERS_EFFORT_PLANNING/rules.md`
   - Add R054 to PRIMARY DIRECTIVES section

2. Edit `/.claude/agents/orchestrator.md`
   - Remove R054 from startup sequence
   - Update rule count from 9 to 8

3. Verify no other references need updating
   - Check all command files
   - Ensure no hard-coded dependencies

## Conclusion

R054 should be moved from orchestrator startup to SPAWN_CODE_REVIEWERS_EFFORT_PLANNING state rules. This follows the Software Factory's just-in-time philosophy, reduces startup overhead, and provides better context for when the rule is actually needed.

The orchestrator doesn't create implementation plans - it just needs to know they exist when spawning the Code Reviewers who will create them. Reading R054 at that moment makes perfect sense.