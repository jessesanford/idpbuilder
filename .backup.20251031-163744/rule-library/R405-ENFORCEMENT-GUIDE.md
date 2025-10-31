# R405 ENFORCEMENT GUIDE: CONTINUE-SOFTWARE-FACTORY Flag Usage

## Quick Reference Card

### ✅ USE TRUE FOR:
- Code review finding issues (NORMAL!)
- Spawning engineers for fixes (NORMAL!)
- Spawning code reviewers (NORMAL!)
- Transitioning between states (NORMAL!)
- Multiple efforts needing work (NORMAL!)
- Size violations detected (NORMAL!)
- Any designed workflow process (NORMAL!)
- After spawning agents (NORMAL!)

### 🔴 USE FALSE ONLY FOR:
- Files missing/corrupt (ERROR!)
- State machine corruption (ERROR!)
- Infrastructure completely broken (ERROR!)
- Cannot proceed without human fix (ERROR!)
- Unrecoverable system failure (ERROR!)

## Decision Algorithm

```
START
  |
  v
Is there a file/system corruption or unrecoverable error?
  |
  ├─ YES → CONTINUE-SOFTWARE-FACTORY=FALSE
  |         (Human intervention required)
  |
  └─ NO → Is this a normal workflow process?
           (spawning, reviewing, fixing, transitioning)
           |
           ├─ YES → CONTINUE-SOFTWARE-FACTORY=TRUE
           |         (Automation continues)
           |
           └─ UNSURE → DEFAULT TO TRUE
                       (System designed for autonomy)
```

## Common Scenarios with Answers

| Scenario | Flag | Reason |
|----------|------|--------|
| Code review found 5 blocking issues | TRUE | Normal fix workflow |
| Spawning engineer to fix issues | TRUE | Standard operation |
| Transitioning to SPAWN_SW_ENGINEERS | TRUE | Designed transition |
| Multiple efforts need fixes | TRUE | System handles multiple |
| Size violation detected | TRUE | Split workflow handles it |
| Review report file missing | FALSE | Cannot proceed (corruption) |
| Infrastructure directory missing | FALSE | System broken |
| State file corrupt/unreadable | FALSE | Cannot determine state |
| R322 requires stop | TRUE | Stop conversation, not automation |

## State-Specific Guidance

### MONITORING_EFFORT_REVIEWS
- Reviews in progress → TRUE
- Reviews found issues → TRUE (this is NORMAL!)
- Reviews approved → TRUE
- Size violations → TRUE (split workflow handles)
- Review file corrupt → FALSE (only this case!)

### SPAWN_SW_ENGINEERS
- Fix instructions exist → TRUE (always!)
- Engineers spawned → TRUE (always!)
- Fix instructions missing → FALSE (corruption!)
- Effort directory missing → FALSE (infrastructure broken!)

### SPAWN_CODE_REVIEWER_FIX_PLAN
- Review reports exist → TRUE (always!)
- Reviewer spawned → TRUE (always!)
- Review report missing → FALSE (corruption!)

### CREATE_NEXT_INFRASTRUCTURE
- Infrastructure created → TRUE
- Pre-plan missing → FALSE (violation!)
- Directories created → TRUE
- Git errors → FALSE (if unrecoverable)

## R322 vs R405 Clarification

**R322 (Stop Inference):**
- Purpose: Preserve context by ending conversation
- Action: `exit 0`
- When: After spawning, at state boundaries
- Effect: Stops THIS conversation turn

**R405 (Continuation Flag):**
- Purpose: Tell automation if it can restart
- Action: `echo "CONTINUE-SOFTWARE-FACTORY=TRUE/FALSE"`
- When: BEFORE exiting
- Effect: Controls whether automation restarts

**KEY INSIGHT:**
```bash
# Normal spawn operation does BOTH:
echo "Spawned agents successfully"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Automation can restart
exit 0                                   # Stop this conversation
```

## Grading Matrix

| Usage Pattern | Penalty |
|---------------|---------|
| TRUE for normal operations | 0% (correct) |
| FALSE for genuine errors | 0% (correct) |
| FALSE for code review issues | -20% |
| FALSE for "user might want to see" | -20% |
| FALSE for normal spawns | -20% |
| Pattern of excessive FALSE (3+) | -50% |
| Confusing R322 with flag | -15% |

## Self-Check Questions

Before setting the flag, ask:

1. **Is the system broken?** (files corrupt, infrastructure missing)
   - YES → FALSE
   - NO → Continue to Q2

2. **Can the system proceed automatically?** (designed workflow)
   - YES → TRUE
   - NO → Continue to Q3

3. **Is human intervention truly required?** (cannot recover)
   - YES → FALSE
   - NO → TRUE (default to automation)

## Examples from Real Issues

### ❌ RECENT VIOLATION (E2.2.2 Case):
```bash
# WRONG - Orchestrator set FALSE when review found issues
echo "Code review found issues in E2.2.2"
echo "Transitioning to SPAWN_SW_ENGINEERS"
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # VIOLATION!
```

### ✅ CORRECT APPROACH:
```bash
# RIGHT - Review finding issues is NORMAL
echo "Code review found issues in E2.2.2"
echo "Fix instructions created"
echo "Transitioning to SPAWN_SW_ENGINEERS"
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"  # Normal workflow!
exit 0  # Stop conversation for context, but automation continues
```

## Red Flags (Incorrect FALSE Usage)

Watch for these phrases that often indicate wrong FALSE usage:

- "User might want to review this" → Usually should be TRUE
- "Code review found issues" → Definitely TRUE
- "Need to spawn engineer" → Definitely TRUE
- "R322 requires stop" → Confusing stop with flag, use TRUE
- "Multiple efforts need work" → Definitely TRUE
- "This is important" → Importance ≠ exceptional, probably TRUE

## Final Principle

**The Software Factory is designed to operate autonomously.**

- Code review finding issues = System working correctly → TRUE
- Needing fixes = Expected workflow → TRUE  
- Complex scenarios = Designed for this → TRUE
- Multiple efforts = Built to handle → TRUE

**Only genuine system failures warrant FALSE.**

---

**See also:**
- R405-automation-flag-continuation-principle.md (full rule)
- R322-mandatory-stop-before-state-transitions.md (stop vs flag)
- Agent state rules for state-specific guidance
