# ⚠️ NO TIME ESTIMATES IN TEMPLATES

## Purpose
AI agents work continuously and don't follow human schedules. Time estimates waste tokens and add no value.

## What NOT to Include
❌ **NEVER include:**
- Weeks, days, hours, minutes
- Morning/afternoon schedules
- Time-based milestones
- Human time estimates
- Duration estimates
- Timeline charts with dates
- Schedule risks
- Day 1, Day 2 references

## What TO Include
✅ **ALWAYS include:**
- Logical sequence (Step 1, Step 2, etc.)
- Dependencies (what blocks what)
- Parallelization opportunities (what can run simultaneously)
- Integration points (when to merge)
- Priority ordering (what must come first)

## Examples

### ❌ WRONG - Time-based:
```
Day 1: Create interfaces
Day 2-3: Implement libraries
Week 2: Deploy features
```

### ✅ RIGHT - Sequence-based:
```
Step 1: Create interfaces (BLOCKING)
Step 2: Implement libraries (depends on Step 1)
Step 3: Deploy features (can parallelize after Step 2)
```

### ❌ WRONG - Duration estimates:
```
Duration: 3 weeks
Time required: 2 days
Estimated hours: 48
```

### ✅ RIGHT - Scope description:
```
Scope: Foundation infrastructure
Complexity: High (multiple integrations)
Dependencies: Requires completed API design
```

## Template Requirements

All planning templates MUST:
1. Focus on logical flow, not time
2. Describe dependencies clearly
3. Identify parallelization opportunities
4. Use step/phase/stage terminology
5. Never mention time units

## Enforcement

The factory manager will:
- Reject templates with time estimates
- Flag violations in reviews
- Update non-compliant templates
- Track token savings from removal

## Token Savings

Removing time estimates typically saves:
- 10-15% of template tokens
- 20-30% of plan generation tokens
- Clearer, more actionable plans

---

**Remember**: AI agents work continuously. Sequence matters, time doesn't.