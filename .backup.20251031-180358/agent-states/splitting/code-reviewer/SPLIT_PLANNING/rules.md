# SPLIT_PLANNING State Rules

## 🔴🔴🔴 CRITICAL: SPLITS PARTITION NEW WORK, NOT EXISTING CODE 🔴🔴🔴

### R359 COMPLIANCE IN SPLIT PLANNING
**Creating split plans that would violate R359 = IMMEDIATE FAILURE**

## ✅ WHAT A SPLIT PLAN MUST DO

**A split plan divides NEW ADDITIONS into reviewable chunks:**
- Each split implements a PORTION of the NEW functionality
- All splits start from the SAME base (usually main)
- Splits are about making PR review manageable
- NO split should require deleting existing code

## 🎯 CORRECT SPLIT PLANNING APPROACH

### Step 1: Analyze the Oversized Implementation
```bash
# Measure ONLY the NEW additions
git checkout [oversized-branch]
added_lines=$(git diff --numstat main..HEAD | awk '{sum+=$1} END {print sum}')
deleted_lines=$(git diff --numstat main..HEAD | awk '{sum+=$2} END {print sum}')

echo "This effort adds $added_lines NEW lines"
echo "This effort deletes $deleted_lines lines"

# If significant deletions, investigate immediately
if [ "$deleted_lines" -gt 100 ]; then
    echo "⚠️ WARNING: Large deletions detected - verify these are intentional refactors"
fi
```

### Step 2: Identify Logical Boundaries in NEW Code
```bash
# Examine what NEW functionality was added
git diff main..HEAD --name-status | grep "^A"  # New files
git diff main..HEAD --stat  # See which files grew

# Group related NEW additions:
# - New feature A: 600 lines
# - New feature B: 500 lines
# - New feature C: 700 lines
# Total: 1800 lines (needs 3 splits)
```

### Step 3: Create the Split Plan

## 📝 SPLIT PLAN TEMPLATE

```markdown
# Split Plan for [Effort-Name]

## Overview
Original effort added [X] lines of NEW functionality, exceeding the 800-line limit.
This plan divides the NEW additions into [N] reviewable splits.

## Base Branch
All splits will branch from: `main` (or specify parent branch)
Current HEAD commit: [commit-hash]

## Split Strategy

### Split 001: [Feature/Component Name]
**Adds:** ~[XXX] lines of NEW code
**Description:** Implements [specific new functionality]
**New files:**
- path/to/new/file1.go
- path/to/new/file2.go

**Modified files (additions only):**
- existing/file.go (+50 lines for new methods)

**Dependencies:** None (first split)
**Can parallelize:** Yes/No

### Split 002: [Feature/Component Name]
**Adds:** ~[XXX] lines of NEW code
**Description:** Implements [specific new functionality]
**New files:**
- path/to/new/file3.go

**Modified files (additions only):**
- existing/file2.go (+30 lines for integration)

**Dependencies:** Independent of split-001
**Can parallelize:** Yes (with split-001)

### Split 003: [Feature/Component Name]
**Adds:** ~[XXX] lines of NEW code
**Description:** Implements [specific new functionality]
**Files:** [List files that will be created/extended]
**Dependencies:** Requires splits 001 and 002 for full functionality
**Can parallelize:** No (depends on previous splits)

## Integration Order
1. Splits 001 and 002 can be implemented in parallel
2. Split 003 must wait for both to complete
3. Final integration merges all splits sequentially

## Validation Checklist
- [ ] Each split adds ≤ 800 lines of NEW code
- [ ] NO split requires deleting existing code
- [ ] All NEW functionality is accounted for
- [ ] Sum of all splits = original NEW additions
- [ ] Each split can be reviewed independently
- [ ] Integration order is clearly defined

## R359 Compliance Statement
This split plan has been verified to comply with R359:
- No existing code will be deleted
- Each split ADDS new functionality
- The existing codebase remains intact
```

## 🚫 WHAT A SPLIT PLAN MUST NEVER DO

### ❌ NEVER Plan to Delete Code
```markdown
# WRONG SPLIT PLAN:
Split 001: Keep authentication module, delete everything else
Split 002: Keep logging module, delete everything else
# THIS VIOLATES R359!
```

### ❌ NEVER Partition Existing Code
```markdown
# WRONG SPLIT PLAN:
Split 001: Move files A-M to this branch, delete N-Z
Split 002: Move files N-Z to this branch, delete A-M
# THIS VIOLATES R359!
```

### ✅ CORRECT: Partition NEW Additions
```markdown
# RIGHT SPLIT PLAN:
Split 001: Add new authentication module (400 lines)
Split 002: Add new authorization module (400 lines)
Split 003: Add new audit logging (400 lines)
# Each builds on main, adds different NEW functionality
```

## 🔍 VALIDATION REQUIREMENTS

Before finalizing any split plan:

### 1. Verify No Deletions Required
```bash
# Each split command should look like:
git checkout main
git checkout -b split-001
# ADD new files and code
git add [new-files]
git commit -m "Add [new functionality]"

# NEVER:
git rm [existing-files]  # R359 VIOLATION!
```

### 2. Verify Total Functionality Preserved
```yaml
Original NEW additions:
  - Feature A: 600 lines
  - Feature B: 500 lines
  - Feature C: 700 lines
  Total: 1800 lines

After all splits merged:
  - Split 001: +600 lines (Feature A)
  - Split 002: +500 lines (Feature B)
  - Split 003: +700 lines (Feature C)
  Total: +1800 lines (MATCHES!)
```

### 3. Verify Each Split is Independent
- Split 001 can be reviewed without split 002
- Split 002 can be reviewed without split 001
- Each adds complete, testable functionality

## 📊 SIZE CALCULATION FOR SPLITS

```bash
# For each planned split, calculate:
# Size = NEW lines only, not total repository size

# Example calculation:
echo "Repository before: 10,000 lines"
echo "Split 001 adds: 600 lines"
echo "Repository after split 001: 10,600 lines"  # This is fine!
echo "Size counted for limit: 600 lines"  # Only NEW additions count
```

## ⚠️ WARNING SIGNS OF R359 VIOLATION

**Immediately revise the plan if you see:**
1. Any mention of "removing" or "deleting" files
2. Splits that "take" different parts of existing code
3. Plans to "divide" the current codebase
4. Instructions to "exclude" existing functionality
5. Any use of `git rm` commands

## 🎓 SPLIT PLAN REVIEW CHECKLIST

Before submitting your split plan:
- [ ] Does each split ADD new code?
- [ ] Is existing code preserved in ALL splits?
- [ ] Can splits be reviewed independently?
- [ ] Do all splits together = original NEW work?
- [ ] No `git rm` commands anywhere?
- [ ] No deletions required?
- [ ] Clear integration strategy?
- [ ] R359 compliance verified?

## 📋 MANDATORY DECLARATION

Include this in every split plan:
```markdown
## R359 Compliance Declaration
I certify that this split plan:
- ✅ Requires NO deletion of existing code
- ✅ Each split ADDS new functionality
- ✅ All splits preserve the existing codebase
- ✅ The 800-line limit applies to NEW code only
- ✅ Violating R359 = IMMEDIATE TERMINATION (-1000%)
```

---

**Remember:** Split plans make large NEW features reviewable by breaking them into smaller additions. They NEVER involve deleting or dividing existing code. The repository GROWS with each split.