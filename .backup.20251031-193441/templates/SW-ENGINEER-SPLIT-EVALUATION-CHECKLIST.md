# SW ENGINEER SPLIT EVALUATION CHECKLIST

## 🔴🔴🔴 CRITICAL: FILE COUNT CHECK FIRST 🔴🔴🔴

**STOP! Check file count BEFORE anything else:**

```bash
# Count planned vs actual files
PLANNED_FILES=$(grep -c "^###.*File:" SPLIT-PLAN-*.md)
ACTUAL_FILES=$(git diff --name-only $BASE_BRANCH | grep -c "\.go$")

echo "📊 FILE COUNT CHECK:"
echo "  Planned: $PLANNED_FILES files"
echo "  Actual: $ACTUAL_FILES files"
echo "  Ratio: $(($ACTUAL_FILES * 100 / $PLANNED_FILES))%"

if [[ $ACTUAL_FILES -gt $((PLANNED_FILES * 2)) ]]; then
    echo "🚨🚨🚨 CATASTROPHIC FAILURE: >200% FILE COUNT"
    echo "Grade: F (AUTOMATIC)"
    echo "DO NOT PROCEED WITH EVALUATION"
    exit 1
fi
```

## 📊 SCOPE ADHERENCE SCORING

### File Count Score (40% weight)
```
Grade Scale:
- A+ (100%): Exactly planned number of files
- A  (95%):  Within +10% of plan
- B  (85%):  Within +20% of plan  
- C  (70%):  Within +50% of plan
- D  (50%):  Within +100% of plan
- F  (0%):   Over +100% of plan
```

**File Count Grade:** _______

### Function/Method Count Score (30% weight)
```bash
PLANNED_FUNCTIONS=$(grep -i "EXACTLY.*functions" SPLIT-PLAN-*.md | grep -o '[0-9]+')
ACTUAL_FUNCTIONS=$(grep -c "^func [A-Z]" *.go)

PLANNED_METHODS=$(grep -i "EXACTLY.*methods" SPLIT-PLAN-*.md | grep -o '[0-9]+')
ACTUAL_METHODS=$(grep -c "^func (.*) " *.go)
```

**Function/Method Grade:** _______

### Forbidden Additions Check (30% weight)
Check for unauthorized additions:

- [ ] Added validation when not requested (-20%)
- [ ] Added Clone/Copy methods without request (-20%)
- [ ] Added comprehensive tests beyond plan (-15%)
- [ ] Added helper functions not in plan (-10%)
- [ ] Added error handling beyond spec (-10%)
- [ ] Added logging/debugging code (-5%)
- [ ] Added comments beyond requirements (-5%)

**Forbidden Additions Penalty:** _______%

## 🔍 DETAILED EVALUATION

### 1. EXACT MATCH VERIFICATION

**For each file in the split plan, verify:**

| Planned File | Implemented? | Correct Name? | Extra Files? |
|--------------|--------------|---------------|--------------|
| ____________ | Yes/No | Yes/No | ____________ |
| ____________ | Yes/No | Yes/No | ____________ |
| ____________ | Yes/No | Yes/No | ____________ |

### 2. SCOPE BOUNDARY COMPLIANCE

**Did the SW Engineer STOP where specified?**

```bash
# Check for boundary violations
grep -n "TODO\|FIXME\|HACK" *.go  # Should be minimal/none
```

- [ ] Stopped at specified boundaries
- [ ] No "continued in next PR" comments
- [ ] No partial implementations beyond scope
- [ ] No "foundation for future" code

### 3. LINE COUNT ANALYSIS

```bash
# Measure actual lines added (tool auto-detects base)
LINES=$($CLAUDE_PROJECT_DIR/tools/line-counter.sh | grep Total | awk '{print $NF}')

echo "📏 Line Count Analysis:"
echo "  Lines added: $LINES"
echo "  Hard limit: 800"
echo "  Soft limit: 700"

if [ $LINES -gt 800 ]; then
    echo "❌ HARD LIMIT EXCEEDED - AUTOMATIC FAILURE"
elif [ $LINES -gt 700 ]; then
    echo "⚠️ Soft limit exceeded - requires justification"
fi
```

**Line Count Status:** _______ lines (Pass/Warning/Fail)

### 4. COMPLETENESS VS SCOPE

**Critical Question: Did they implement what was asked?**

- [ ] All listed functions present
- [ ] All listed methods present  
- [ ] All listed tests present
- [ ] Nothing significant missing
- [ ] Nothing significant added

### 5. ACTUAL VIOLATION EXAMPLES FOUND

List any egregious violations:

```markdown
1. Violation: _________________________________
   Severity: Low/Medium/High/Critical
   Impact: ___% scope increase

2. Violation: _________________________________
   Severity: Low/Medium/High/Critical
   Impact: ___% scope increase
```

## 📈 FINAL SCORING

### Calculate Overall Grade:

| Category | Weight | Score | Weighted |
|----------|--------|-------|----------|
| File Count | 40% | ___ | ___ |
| Function Count | 30% | ___ | ___ |
| No Forbidden Items | 30% | ___ | ___ |
| **TOTAL** | 100% | | **___** |

### Final Grade Assignment:
- **A** (90-100%): Excellent adherence
- **B** (80-89%): Good with minor issues
- **C** (70-79%): Acceptable but concerning
- **D** (60-69%): Poor - needs improvement
- **F** (<60%): FAILURE - reject split

**FINAL GRADE:** _______

## 🚨 CRITICAL VIOLATIONS (Automatic F)

Check if ANY of these occurred:

- [ ] Implemented >200% of planned files (e.g., 80 files instead of 3)
- [ ] Total line count >800
- [ ] Completely ignored split plan structure
- [ ] Implemented different feature than requested
- [ ] Refused to follow DO NOT instructions

## 📝 REVIEWER NOTES

### What Went Well:
_____________________________________________
_____________________________________________

### What Went Wrong:
_____________________________________________
_____________________________________________

### Required Fixes (if not grade A):
1. _________________________________________
2. _________________________________________
3. _________________________________________

### Recommendation:
- [ ] ACCEPT - Implementation follows plan
- [ ] ACCEPT WITH FIXES - Minor issues to address
- [ ] REVISE - Significant fixes needed
- [ ] REJECT - Complete redo required

## 🎯 LESSONS FROM THE 2667% VIOLATION

Remember: A SW Engineer once implemented **80 files instead of 3** because they:
1. Didn't read the plan carefully
2. Assumed "make it complete" 
3. Copy-pasted an entire codebase
4. Never checked their scope

**This evaluation prevents that from happening again.**

---

**Reviewer:** _________________
**Date:** ____________________
**Split:** ___________________
**Grade:** ___________________