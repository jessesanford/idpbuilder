# SW ENGINEER PRE-IMPLEMENTATION CHECKLIST - MANDATORY

## 🔴🔴🔴 CATASTROPHIC FAILURE PREVENTION 🔴🔴🔴

**THIS CHECKLIST PREVENTS 2667% SCOPE VIOLATIONS**

On [DATE], a SW Engineer implemented **80 files instead of 3** - a 2667% violation that caused complete split failure. This checklist ensures that NEVER happens again.

## 🛑 STOP! DO NOT WRITE ANY CODE UNTIL THIS IS COMPLETE

### 📊 SECTION 1: QUANTITATIVE SCOPE VERIFICATION

**Extract EXACT counts from the split plan:**

#### File Count Verification
```bash
# Count files explicitly listed in plan
PLANNED_FILES=$(grep -c "^###.*File:" SPLIT-PLAN-*.md)
echo "📁 FILES TO IMPLEMENT: $PLANNED_FILES (EXACTLY)"

# Set HARD limits
MAX_FILES=$((PLANNED_FILES * 2))  # 200% absolute maximum
STOP_FILES=$((PLANNED_FILES + 2))  # Stop if 2 files over

echo "🛑 HARD STOP at $MAX_FILES files"
echo "⚠️  WARNING at $STOP_FILES files"
```

- [ ] I will implement **EXACTLY** _______ files (write number)
- [ ] I understand exceeding this by 2x = AUTOMATIC FAILURE
- [ ] I have written these limits on paper next to my screen

#### Function/Method Count Verification
```bash
PLANNED_FUNCTIONS=$(grep -i "EXACTLY.*functions" SPLIT-PLAN-*.md | grep -o '[0-9]+' | head -1 || echo 0)
PLANNED_METHODS=$(grep -i "EXACTLY.*methods" SPLIT-PLAN-*.md | grep -o '[0-9]+' | head -1 || echo 0)
PLANNED_TESTS=$(grep -i "EXACTLY.*tests" SPLIT-PLAN-*.md | grep -o '[0-9]+' | head -1 || echo 0)

echo "📋 SCOPE LIMITS:"
echo "  Functions: ${PLANNED_FUNCTIONS}"
echo "  Methods: ${PLANNED_METHODS}"
echo "  Tests: ${PLANNED_TESTS}"
```

- [ ] Functions to implement: _______ (write number or "0")
- [ ] Methods to implement: _______ (write number or "0")
- [ ] Tests to write: _______ (write number or "0")

### 🚫 SECTION 2: FORBIDDEN ACTIONS (DO NOT LIST)

**Read and acknowledge what you MUST NOT do:**

```bash
# Extract DO NOT instructions
echo "═══════════════════════════════════════════"
echo "🚫 FORBIDDEN ACTIONS FROM SPLIT PLAN:"
echo "═══════════════════════════════════════════"
grep -A 10 "DO NOT\|FORBIDDEN\|STOP BOUNDARIES" SPLIT-PLAN-*.md
```

**Check ALL that apply from your split plan:**
- [ ] DO NOT add validation/Validate methods
- [ ] DO NOT add Clone/Copy methods
- [ ] DO NOT write comprehensive tests
- [ ] DO NOT refactor existing code
- [ ] DO NOT add helper functions not listed
- [ ] DO NOT implement error types not specified
- [ ] DO NOT create interfaces not requested
- [ ] DO NOT add configuration not in plan
- [ ] DO NOT write documentation beyond plan
- [ ] DO NOT optimize or improve performance

### 📝 SECTION 3: EXACT IMPLEMENTATION LIST

**List EVERY file you will create/modify (must match plan exactly):**

1. File: _________________________ (Action: create/modify)
2. File: _________________________ (Action: create/modify)
3. File: _________________________ (Action: create/modify)
4. File: _________________________ (Action: create/modify)
5. File: _________________________ (Action: create/modify)

**If more than 5 files, STOP and verify with plan!**

### 🎯 SECTION 4: SCOPE BOUNDARIES

**Where will you STOP implementation?**

```bash
# Display the boundaries
echo "🛑 STOP BOUNDARIES:"
grep -B2 -A5 "stop.*here\|boundary\|limit" SPLIT-PLAN-*.md -i
```

Write your STOP boundaries:
1. Stop after: _________________________________
2. Stop before: ________________________________
3. Will NOT continue to: _______________________

### 📏 SECTION 5: SIZE ESTIMATION

**Realistic line counts per element:**

| Element | Typical Lines | Your Estimate |
|---------|--------------|---------------|
| Simple function | 10-20 | _____ |
| Complex function | 30-50 | _____ |
| Method with validation | 20-40 | _____ |
| Basic test | 15-30 | _____ |
| Struct definition | 5-15 | _____ |
| Interface | 3-10 | _____ |

**Total estimated lines:** _______ (MUST be < 600 for safety)

### ✅ SECTION 6: ACKNOWLEDGMENTS

**Initial each statement:**

- _____ I have read the entire split plan
- _____ I have identified the EXACT files to implement
- _____ I understand what NOT to implement
- _____ I will STOP at the specified boundaries
- _____ I will NOT add "helpful" extras
- _____ I will NOT make it "complete"
- _____ I will check file count every 100 lines
- _____ I understand >200% scope = AUTOMATIC FAILURE
- _____ I will ASK if the plan seems insufficient
- _____ I have seen the 2667% violation example

### 🚨 SECTION 7: CONTINUOUS MONITORING COMMITMENT

**I commit to running these checks:**

#### Every 100 lines:
```bash
ACTUAL_FILES=$(find . -type f -name "*.go" -newer .start-marker | wc -l)
echo "Current files: $ACTUAL_FILES / $PLANNED_FILES"
[ $ACTUAL_FILES -gt $STOP_FILES ] && echo "🚨 STOP NOW!"
```

#### Before EVERY commit:
```bash
validate_split_scope() {
    ACTUAL=$(find . -type f -name "*.go" | wc -l)
    PLANNED=$(grep -c "^###.*File:" SPLIT-PLAN-*.md)
    
    if [[ $ACTUAL -gt $((PLANNED * 2)) ]]; then
        echo "🚨 FATAL: $ACTUAL files vs $PLANNED planned"
        echo "This is a $((ACTUAL * 100 / PLANNED))% violation!"
        exit 1
    fi
}
```

### 📅 TIMESTAMP AND SIGNATURE

**This checklist was completed at:**

Date: _____________________
Time: _____________________
Split: ____________________
Files to implement: ________

**By completing this checklist, I acknowledge that:**
- Implementing 80 files instead of 3 is UNACCEPTABLE
- I will be graded on scope adherence (>150% = FAIL)
- The split plan is my CONTRACT
- Incomplete is better than oversized

---

## 🔴 REMEMBER THE CATASTROPHIC FAILURE

**A SW Engineer ignored the split plan and implemented 80 files instead of 3.**

This was not a mistake - it was a complete disregard for instructions.

**DO NOT BE THAT ENGINEER.**

Your reputation depends on following the plan EXACTLY.