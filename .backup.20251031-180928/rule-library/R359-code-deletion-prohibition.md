# 🔴🔴🔴 RULE R359 - ABSOLUTE PROHIBITION ON DELETING APPROVED CODE

## SUPREME LAW #6 - PENALTY: IMMEDIATE TERMINATION (-1000%)

### ABSOLUTE PROHIBITION
You MUST NEVER delete code that:
- Was already merged to main/master
- Was previously approved in code review
- Exists in the base branch you're working from
- Is part of the existing codebase

### THE 800-LINE LIMIT APPLIES ONLY TO NEW CODE
- Measure ONLY the lines YOU ADD
- NEVER delete existing code to fit the limit
- If your NEW additions exceed 800 lines, SPLIT the work
- The limit is about YOUR CHANGES, not total repository size

### WHAT "SPLITTING" ACTUALLY MEANS
❌ WRONG: Delete everything except your 800-line portion
✅ CORRECT: Break your NEW work into 800-line increments

Example of CORRECT splitting:
- You need to add 2000 lines of NEW functionality
- Split into: 800 lines (split-001), 800 lines (split-002), 400 lines (split-003)
- Each split ADDS to the codebase, building on the previous
- NEVER delete existing code to make room

### MANDATORY APPROVAL FOR ANY DELETIONS
Before deleting ANY existing code:
1. STOP IMMEDIATELY
2. Request ARCHITECT review
3. Request CODE REVIEWER approval
4. Document justification in detail
5. ONLY proceed with explicit written approval

### EMERGENCY STOP CONDITIONS
If you find yourself:
- Deleting files to meet size limits → STOP AND EXIT 359
- Removing packages for line count → STOP AND EXIT 359
- Deleting main.go, LICENSE, README → STOP AND EXIT 359
- Removing existing functionality → STOP AND EXIT 359

Exit with code 359 immediately!

### CRITICAL EXAMPLES OF VIOLATIONS

#### Example 1: CATASTROPHIC VIOLATION
```bash
# WRONG - This is what happened in the disaster
git rm -rf pkg/build/
git rm -rf pkg/cmd/
git rm -rf pkg/controllers/
git rm -rf pkg/k8s/
git rm -rf pkg/resources/
git rm main.go
git rm Makefile
git rm LICENSE
git rm README.md
# "Now my 595 lines fits in the 800 limit!"
# THIS DELETED 9,552 LINES OF APPROVED CODE!
```

#### Example 2: CORRECT APPROACH
```bash
# RIGHT - Split NEW work only
# Starting from main with existing 10,000 lines
git checkout -b effort-001
# Add 800 lines of NEW code
# Total repo now: 10,800 lines (that's fine!)

git checkout main
git checkout -b effort-002
# Add another 800 lines of NEW code
# Total repo now: 10,800 lines

# Later merge both efforts
# Final repo: 11,600 lines (10,000 original + 1,600 new)
```

### ENFORCEMENT MECHANISMS
```bash
# MANDATORY CHECK IN ALL IMPLEMENTATION STATES
deleted_lines=$(git diff --numstat "$BASE..$HEAD" | awk '{sum+=$2} END {print sum}')
if [ "$deleted_lines" -gt 100 ]; then
    echo "🔴🔴🔴 R359 VIOLATION DETECTED!"
    echo "You are attempting to delete $deleted_lines lines!"
    echo "NEVER DELETE EXISTING CODE TO MEET SIZE LIMITS!"
    exit 359
fi

# CHECK FOR CRITICAL FILE DELETION
if git diff --name-status "$BASE" | grep -E "^D.*main\.(go|py|js|ts)|^D.*LICENSE|^D.*README"; then
    echo "🔴🔴🔴 R359 CRITICAL VIOLATION!"
    echo "Attempting to delete critical project files!"
    exit 359
fi
```

### GRADING IMPACT
- Any violation of R359: **IMMEDIATE -1000% (FACTORY SHUTDOWN)**
- This is SUPREME LAW #6 - highest possible penalty
- There is NO tolerance for this behavior
- Violating agents will be terminated immediately

### WHY THIS RULE EXISTS
On [DATE], agents misunderstood "splitting" and deleted 9,552 lines of already-approved code, including entire packages, main.go, LICENSE, and README, thinking they needed to make the repository fit within 800 lines. This catastrophic misunderstanding nearly destroyed an entire project.

### RELATED RULES
- R220/R221: Size limits (applies to NEW code only)
- R304: Line counting methodology
- R251: Split planning protocols
- R309: Split execution procedures

### MANDATORY ACKNOWLEDGMENT
Every agent MUST acknowledge this rule on startup:
```bash
echo "📋 Acknowledging R359: NEVER delete approved code for size limits"
echo "  ✅ Size limits apply to NEW code only"
echo "  ✅ Splitting means breaking NEW work into pieces"
echo "  ✅ Deleting existing code = IMMEDIATE TERMINATION"
```

---

**THIS IS SUPREME LAW #6 - ABSOLUTE AND INVIOLABLE**