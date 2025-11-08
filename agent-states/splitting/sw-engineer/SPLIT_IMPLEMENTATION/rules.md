# SPLIT_IMPLEMENTATION State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## 🔴🔴🔴 ABSOLUTE RULE: NO CODE DELETION ALLOWED 🔴🔴🔴

### R359 SUPREME LAW ENFORCEMENT
**PENALTY: IMMEDIATE TERMINATION (-1000%)**

## ✅ WHAT SPLIT IMPLEMENTATION MEANS

**Splitting is about PARTITIONING NEW WORK, not dividing existing code!**

When you receive a split plan, you are implementing a PORTION of NEW functionality:
- Each split ADDS new code to the existing codebase
- Each split builds ON TOP of existing code
- NO split should EVER delete existing functionality
- ALL splits combined = the original oversized implementation

## 🔴 FORBIDDEN ACTIONS IN THIS STATE

**NEVER DO ANY OF THESE:**
```bash
# 🚨 THESE ARE ALL R359 VIOLATIONS - IMMEDIATE TERMINATION
git rm -rf pkg/        # NEVER delete packages
git rm main.go         # NEVER delete main files
git rm LICENSE         # NEVER delete project files
git rm README.md       # NEVER delete documentation

# 🚨 WRONG UNDERSTANDING OF SPLITS
# Thinking: "I need to remove code to make my split fit"
# THIS IS COMPLETELY WRONG!
```

## ✅ CORRECT SPLIT IMPLEMENTATION PATTERN

### Step 1: Understand Your Split Assignment
```bash
# Read your split plan to understand what NEW code to add
cat .software-factory/phase${PHASE}/wave${WAVE}/${EFFORT}/SPLIT-PLAN--*.md
# Look for: "Split 001: Add authentication module (400 lines)"
# This means ADD 400 lines of NEW authentication code
```

### Step 2: Start From the Correct Base
```bash
# Splits build on existing code, not replace it
git checkout [base-branch]  # Usually main or previous split
git checkout -b [split-branch]
```

### Step 3: Add Your Assigned NEW Code
```bash
# Only ADD the new functionality assigned to your split
# Example: If split-001 is "Add user authentication"
touch pkg/auth/user.go      # NEW file
echo "// New auth code" >> pkg/auth/user.go

# The existing codebase remains intact!
```

### Step 4: Verify No Deletions
```bash
# MANDATORY CHECK BEFORE EVERY COMMIT
deleted_lines=$(git diff --numstat HEAD | awk '{sum+=$2} END {print sum}')
if [ "$deleted_lines" -gt 0 ]; then
    echo "🔴🔴🔴 R359 VIOLATION DETECTED!"
    echo "You are attempting to delete $deleted_lines lines!"
    echo "Splits must ADD code, not DELETE code!"
    exit 359
fi

# Check for file deletions
if git status --porcelain | grep -E "^D "; then
    echo "🔴🔴🔴 R359 VIOLATION: Files marked for deletion!"
    echo "NEVER delete files during splits!"
    exit 359
fi
```

## 📊 SPLIT SIZE CALCULATION

**The 800-line limit applies to NEW CODE ONLY:**
```bash
# Correct measurement: Count only ADDITIONS
added_lines=$(git diff --numstat main..HEAD | awk '{sum+=$1} END {print sum}')
echo "This split adds $added_lines NEW lines"

# The limit is on NEW lines, not total repository size
if [ "$added_lines" -gt 800 ]; then
    echo "Split too large, needs further splitting"
fi
```

## 🎯 CORRECT SPLIT EXAMPLE

**Scenario:** Original effort added 2400 lines (too large)
**Solution:** Split into 3 parts

### Split 001 (800 lines of NEW code)
```bash
git checkout main
git checkout -b effort-E1.1-split-001
# Add authentication module (800 NEW lines)
# Repository: 10,000 (existing) + 800 (new) = 10,800 total
```

### Split 002 (800 lines of NEW code)
```bash
git checkout main  # Start from same base
git checkout -b effort-E1.1-split-002
# Add authorization module (800 NEW lines)
# Repository: 10,000 (existing) + 800 (new) = 10,800 total
```

### Split 003 (800 lines of NEW code)
```bash
git checkout main  # Start from same base
git checkout -b effort-E1.1-split-003
# Add logging module (800 NEW lines)
# Repository: 10,000 (existing) + 800 (new) = 10,800 total
```

### Integration Result
When all splits are merged:
- Repository has 10,000 (original) + 2400 (new from all splits) = 12,400 total lines
- NO code was deleted
- All functionality preserved

## ⚠️ COMMON MISUNDERSTANDINGS TO AVOID

### ❌ WRONG: "Divide the existing code between splits"
```bash
# WRONG THINKING:
# "The repo has 10,000 lines, I'll take 800 and delete the rest"
git rm -rf [most of the repo]  # R359 VIOLATION!
```

### ❌ WRONG: "Remove code from split-002 that's in split-001"
```bash
# WRONG THINKING:
# "Split-001 has auth code, so split-002 shouldn't have it"
git rm -rf pkg/auth/  # R359 VIOLATION!
```

### ✅ CORRECT: "Each split adds its portion of NEW work"
```bash
# RIGHT THINKING:
# "Split-001 adds auth, split-002 adds logging, both build on main"
# Each split is independent, adding its NEW functionality
```

## 🛑 EMERGENCY STOP CONDITIONS

**IMMEDIATELY STOP and exit 359 if you find yourself:**
1. Running `git rm` on any existing file
2. Deleting code to "make room" for your split
3. Removing packages or modules
4. Thinking about "dividing" existing code
5. Deleting main.go, LICENSE, README, or core files

```bash
echo "🔴🔴🔴 R359 EMERGENCY STOP!"
echo "I was about to delete existing code!"
echo "Splits ADD new code, never DELETE existing code!"
exit 359
```

## 📋 MANDATORY ACKNOWLEDGMENT

Before starting ANY split implementation:
```bash
echo "📋 Acknowledging SPLIT_IMPLEMENTATION rules:"
echo "  ✅ Splits ADD new code, never DELETE existing code"
echo "  ✅ Each split implements a PORTION of NEW functionality"
echo "  ✅ The 800-line limit applies to NEW additions only"
echo "  ✅ All splits build ON TOP of existing code"
echo "  ✅ R359 violation = IMMEDIATE TERMINATION"
```

## 🎓 FINAL UNDERSTANDING CHECK

**Ask yourself:**
- Am I ADDING new code? ✅
- Is the existing codebase intact? ✅
- Will all functionality still work? ✅
- Am I only implementing my assigned NEW features? ✅

If ANY answer is NO, STOP IMMEDIATELY!

---

**Remember:** Splits make large NEW additions reviewable by breaking them into smaller pieces. They NEVER remove existing, approved code. The repository GROWS with each split, it doesn't get divided.