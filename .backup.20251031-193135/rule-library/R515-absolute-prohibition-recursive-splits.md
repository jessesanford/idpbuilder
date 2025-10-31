# Rule R515: ABSOLUTE PROHIBITION ON RECURSIVE SPLITS (Splits of Splits)

## 🔴🔴🔴 SUPREME LAW: NO SPLITS OF SPLITS EVER 🔴🔴🔴

## Rule Statement
**RECURSIVE SPLITS ARE ABSOLUTELY FORBIDDEN.** If a split exceeds 800 lines, this indicates a FUNDAMENTAL DESIGN FAILURE that requires human intervention. Attempting to split a split destroys the entire split protocol architecture and causes cascading measurement failures.

## Criticality Level
🚨🚨🚨 **BLOCKING** - Recursive splits = IMMEDIATE -100% FAILURE 🚨🚨🚨

## Why This Is DEADLY SERIOUS

**SPLITS EXIST TO CREATE MANAGEABLE CHUNKS:**
- Original effort: Too large (>800 lines)
- Split-001: Should be <800 lines
- Split-002: Should be <800 lines
- Split-003: Should be <800 lines

**IF A SPLIT NEEDS SPLITTING, YOU'VE FAILED:**
- The original split plan was wrong
- The scope boundaries were incorrectly drawn
- The implementation violated the split plan
- Human architect review is MANDATORY

## What Happens With Recursive Splits

```
CATASTROPHIC CASCADE:
effort (2400 lines) →
  split-001 (1200 lines) →
    split-001-001 (600 lines) ← FORBIDDEN!
    split-001-002 (600 lines) ← FORBIDDEN!
  split-002 (800 lines) → OK
  split-003 (400 lines) → OK

RESULT:
- Dependency chains become incomprehensible
- Line counting becomes impossible
- Merge ordering breaks down
- Integration becomes a nightmare
- System enters undefined state
```

## MANDATORY PROTOCOL When Split Exceeds Limits

```bash
# SW Engineer detects split exceeding 800 lines
if [ $SPLIT_LINES -gt 800 ]; then
    echo "🔴🔴🔴 CRITICAL FAILURE DETECTED 🔴🔴🔴"
    echo "Split ${SPLIT_NUMBER} has ${SPLIT_LINES} lines"
    echo "This is a RECURSIVE SPLIT scenario"
    echo ""
    echo "MANDATORY ACTIONS:"
    echo "1. STOP all work immediately"
    echo "2. Set CONTINUE-SOFTWARE-FACTORY=FALSE"
    echo "3. Request human architect intervention"
    echo "4. DO NOT attempt to split the split"
    echo ""
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"  # R405 compliance
    exit 1
fi
```

## Root Cause Analysis Required

When a split exceeds limits, investigate:

1. **Was the original split plan wrong?**
   - Files incorrectly grouped
   - Dependencies misunderstood
   - Scope creep during implementation

2. **Did implementation violate the plan?**
   - Extra files added beyond plan
   - Scope expanded during coding
   - Generated code incorrectly included

3. **Is the feature inherently unsplittable?**
   - Core framework changes
   - Tightly coupled components
   - Fundamental architectural pieces

## Correct Resolution Paths

### Option 1: Re-plan ALL Splits
```bash
# Delete all split infrastructure
rm -rf *-SPLIT-*
# Code Reviewer creates NEW split plan
# More, smaller splits from the beginning
```

### Option 2: Accept Oversized PR (With Approval)
```bash
# Sometimes a large PR is unavoidable
# Requires explicit architect approval
# Document why splitting failed
# Add review requirements for large PR
```

### Option 3: Redesign the Feature
```bash
# The feature itself may be too monolithic
# Break into separate, independent features
# Each feature gets its own effort
# No single effort exceeds limits
```

## Integration with Other Rules

- **R202**: Single agent implements all splits sequentially
- **R204**: Orchestrator creates split infrastructure just-in-time
- **R207**: Split boundary validation
- **R405**: Sets CONTINUE-SOFTWARE-FACTORY=FALSE for recursive splits
- **R297**: Architect must detect and prevent split re-splitting

## Detection Points

### Code Reviewer
```bash
# When planning splits
if [ $ESTIMATED_SPLIT_SIZE -gt 700 ]; then
    echo "⚠️ WARNING: Split approaching limit"
    echo "Consider creating more, smaller splits"
fi
```

### SW Engineer
```bash
# After implementing split
ACTUAL_SIZE=$($PROJECT_ROOT/tools/line-counter.sh)
if [ $ACTUAL_SIZE -gt 800 ]; then
    echo "❌ RECURSIVE SPLIT DETECTED"
    # Follow mandatory protocol above
fi
```

### Orchestrator
```bash
# Never attempt to create split-of-split infrastructure
if [[ "$EFFORT_NAME" == *"-split-"* ]]; then
    echo "❌ FATAL: Attempting to split a split"
    echo "This is FORBIDDEN by R515"
    exit 1
fi
```

## Grading Impact

- **Attempting recursive split**: -100% (IMMEDIATE FAILURE)
- **Creating split-of-split infrastructure**: -100%
- **Not detecting oversized split**: -50%
- **Proceeding after recursive split detected**: -100%

## Remember

**SPLITS ARE NOT INFINITELY RECURSIVE!**

If you can't fit the implementation into <800 line chunks on the FIRST split attempt, you have a design problem, not an implementation problem. Stop and get human help.

**The split protocol is about creating reviewable, mergeable chunks - not about playing Tetris with code until it fits.**