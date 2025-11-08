# 🔴🔴🔴 RULE R295: SW ENGINEER SPAWN CLARITY PROTOCOL (SUPREME)

## PRIORITY: SUPREME/ABSOLUTE

This rule is **ABSOLUTELY MANDATORY** to prevent SW Engineer confusion and wasted effort. Violations directly cause implementation failures.

## WHEN THIS RULE APPLIES

This rule applies **EVERY TIME** the orchestrator spawns a SW Engineer, especially for:
1. Initial implementation work
2. Code review fixes
3. Split implementation fixes
4. Integration fixes (phase, wave, or project)
5. Any fix or rework cycles

## MANDATORY SPAWN MESSAGE COMPONENTS

### 1. STATE SPECIFICATION (CRITICAL)

**EVERY spawn command MUST explicitly specify the state:**

```markdown
🔴🔴🔴 CRITICAL STATE INFORMATION:
YOU ARE IN STATE: [EXACT_STATE_NAME]
This means you should: [brief description of what this state entails]
🔴🔴🔴
```

Valid states:
- `IMPLEMENTATION` - Initial feature implementation
- `FIX_ISSUES` - Fixing code review issues
- `SPLIT_IMPLEMENTATION` - Working on a specific split
- `FIX_INTEGRATE_WAVE_EFFORTS_ISSUES` - Fixing integration problems
- `TEST_WRITING` - Writing tests

### 2. PLAN FILE SPECIFICATION (CRITICAL)

**EVERY spawn command MUST specify the EXACT plan file to follow:**

```markdown
📋 YOUR INSTRUCTIONS:
FOLLOW ONLY: [EXACT-FILENAME.md]
LOCATION: [exact path, e.g., "in your effort directory"]
IGNORE: Any files named *-COMPLETED-*.md or other plan files
```

### 3. CONTEXT CLARITY (CRITICAL)

**EVERY spawn command MUST provide clear context:**

```markdown
🎯 CONTEXT:
- EFFORT: [effort name]
- WAVE: [wave number]
- PHASE: [phase name if applicable]
- PREVIOUS WORK: [brief summary of what was done before]
- YOUR TASK: [specific task for this spawn]
```

### 4. WARNING ABOUT OLD PLANS

**For fix work, MUST include warning:**

```markdown
⚠️⚠️⚠️ IMPORTANT:
- Old fix plans have been archived as *-COMPLETED-*.md
- DO NOT follow these archived plans
- ONLY follow the plan specified above
⚠️⚠️⚠️
```

## SPAWN COMMAND TEMPLATES

### Template 1: Initial Implementation
```bash
/orchestrate spawn sw-engineer effort-user-auth \
    --working-dir "phase-1/wave-1/effort-user-auth" \
    --message "
🔴🔴🔴 CRITICAL STATE INFORMATION:
YOU ARE IN STATE: IMPLEMENTATION
This means you should: Implement the features defined in EFFORT-PLAN.md
🔴🔴🔴

📋 YOUR INSTRUCTIONS:
FOLLOW ONLY: EFFORT-PLAN.md
LOCATION: In your effort directory
IGNORE: Any other plan files

🎯 CONTEXT:
- EFFORT: user-auth
- WAVE: 1
- PHASE: 1
- YOUR TASK: Implement user authentication as specified in EFFORT-PLAN.md
"
```

### Template 2: Code Review Fixes
```bash
/orchestrate spawn sw-engineer effort-database \
    --working-dir "phase-1/wave-2/effort-database" \
    --message "
🔴🔴🔴 CRITICAL STATE INFORMATION:
YOU ARE IN STATE: FIX_ISSUES
This means you should: Fix the issues identified in CODE-REVIEW-REPORT.md
🔴🔴🔴

📋 YOUR INSTRUCTIONS:
FOLLOW ONLY: CODE-REVIEW-REPORT.md
LOCATION: In your effort directory
IGNORE: EFFORT-PLAN.md and any *-COMPLETED-*.md files

⚠️⚠️⚠️ IMPORTANT:
- Previous implementation is complete
- Only fix the specific issues listed in CODE-REVIEW-REPORT.md
- Do not add new features or refactor beyond the fixes
⚠️⚠️⚠️

🎯 CONTEXT:
- EFFORT: database
- WAVE: 2
- PHASE: 1
- PREVIOUS WORK: Initial implementation complete, review found issues
- YOUR TASK: Fix all MUST_FIX issues in CODE-REVIEW-REPORT.md
"
```

### Template 3: Integration Fixes
```bash
/orchestrate spawn sw-engineer effort-api-gateway \
    --working-dir "phase-2/wave-1/effort-api-gateway" \
    --message "
🔴🔴🔴 CRITICAL STATE INFORMATION:
YOU ARE IN STATE: FIX_INTEGRATE_WAVE_EFFORTS_ISSUES
This means you should: Fix integration issues found during phase integration
🔴🔴🔴

📋 YOUR INSTRUCTIONS:
FOLLOW ONLY: INTEGRATE_WAVE_EFFORTS-REPORT.md
LOCATION: In your effort directory (just copied there)
IGNORE: Any files named *-COMPLETED-*.md (these are from previous fix cycles)

⚠️⚠️⚠️ IMPORTANT:
- SPLIT-PLAN-COMPLETED-*.md = old, already done
- CODE-REVIEW-REPORT-COMPLETED-*.md = old, already done
- ONLY follow INTEGRATE_WAVE_EFFORTS-REPORT.md
⚠️⚠️⚠️

🎯 CONTEXT:
- EFFORT: api-gateway
- WAVE: 1
- PHASE: 2
- PREVIOUS WORK: Phase integration revealed compatibility issues
- YOUR TASK: Fix ONLY the issues for your effort listed in INTEGRATE_WAVE_EFFORTS-REPORT.md
"
```

## VERIFICATION CHECKLIST

Before spawning any SW Engineer:
- [ ] State explicitly specified in spawn message
- [ ] Exact plan filename provided
- [ ] Location of plan file clarified
- [ ] Context information included
- [ ] Warning about old plans (if applicable)
- [ ] Task clearly defined

## VIOLATIONS

**This rule is violated if:**
- ❌ Spawn message doesn't specify the state
- ❌ Spawn message doesn't specify which plan to follow
- ❌ Spawn message is ambiguous about the task
- ❌ Multiple plan files exist without clear instruction
- ❌ SW Engineer has to guess what to do

## PENALTIES

- Missing state specification = -25% penalty
- Missing plan specification = -30% penalty
- Ambiguous instructions = -40% penalty
- SW Engineer confusion leading to wrong work = -75% penalty
- Complete implementation failure due to confusion = -100% penalty

## COMMON MISTAKES TO AVOID

### ❌ BAD: Ambiguous spawn
```bash
/orchestrate spawn sw-engineer effort-logging \
    --message "Fix the issues"  # WRONG: Which issues? What plan?
```

### ❌ BAD: Missing state
```bash
/orchestrate spawn sw-engineer effort-auth \
    --message "Follow INTEGRATE_WAVE_EFFORTS-REPORT.md"  # WRONG: What state?
```

### ❌ BAD: Unclear plan location
```bash
/orchestrate spawn sw-engineer effort-db \
    --message "Fix issues in the report"  # WRONG: Which report? Where?
```

### ✅ GOOD: Crystal clear
```bash
/orchestrate spawn sw-engineer effort-api \
    --message "
🔴 STATE: FIX_INTEGRATE_WAVE_EFFORTS_ISSUES
📋 FOLLOW: INTEGRATE_WAVE_EFFORTS-REPORT.md (in your effort directory)
⚠️ IGNORE: *-COMPLETED-*.md files
🎯 TASK: Fix the 3 API compatibility issues listed for your effort
"
```

## ENFORCEMENT IN ORCHESTRATOR

**Orchestrator MUST have a spawn clarity check:**

```python
def validate_spawn_message(message: str) -> bool:
    """Ensure spawn message meets clarity requirements"""
    
    required_elements = [
        "STATE:",           # Must specify state
        "FOLLOW",           # Must specify what to follow
        ".md",              # Must reference specific file
        "TASK:",            # Must define the task
    ]
    
    for element in required_elements:
        if element not in message:
            print(f"❌ Spawn message missing: {element}")
            return False
    
    return True
```

## RELATED RULES

- R293: Integration Report Distribution Protocol
- R294: Fix Plan Archival Protocol
- R052: Agent Spawning Protocol
- R203: State-Aware Agent Startup

---

**REMEMBER**: Confusion wastes time and causes failures. ALWAYS provide crystal-clear instructions when spawning agents!