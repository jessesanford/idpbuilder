# Split Working Directory Guide - Critical Instructions

## 🚨 THE GOLDEN RULE OF SPLITS

**ALL WORK HAPPENS IN THE SPLIT DIRECTORY SPECIFIED IN THE SPLIT PLAN**

Once you navigate to a split directory, you NEVER leave it until that split is complete.

## For SW Engineers Working on Splits

### Step 1: Read the Split Plan Metadata
```bash
# You start in the too-large effort directory
pwd  # e.g., /efforts/phase1/wave1/api-types

# Find your split plan
SPLIT_PLAN=$(ls SPLIT-PLAN-*.md | head -1)

# Extract the working directory
WORKING_DIR=$(grep "\*\*WORKING_DIRECTORY\*\*:" "$SPLIT_PLAN" | cut -d: -f2- | xargs)
echo "Must work in: $WORKING_DIR"
```

### Step 2: Navigate to the Split Directory
```bash
# Change to the split directory
cd "$WORKING_DIR"

# Verify you're in the right place
pwd  # Should show: /efforts/phase1/wave1/api-types--split-001
```

### Step 3: ALL Work Happens Here

#### ✅ CORRECT - Working in Split Directory
```bash
# All file operations in split directory
ls ./pkg/                    # List split files
vim ./pkg/api/types.go      # Edit split files
grep -r "TODO" ./pkg/       # Search split files

# All git operations in split directory
git status                  # Shows split branch status
git add ./pkg/              # Stages split files
git commit -m "feat: implement split-001"
git push                    # Pushes split branch

# All build/test operations in split directory
go build ./pkg/...          # Build split code
go test ./pkg/...           # Test split code
make test                   # Run split tests

# Tool usage - all relative to split directory
Read: pkg/api/types.go      # Reads from split directory
Write: pkg/api/helpers.go   # Writes to split directory
Edit: pkg/api/types.go      # Edits in split directory
Bash: go test ./pkg/...     # Runs in split directory
```

#### ❌ WRONG - Working Outside Split Directory
```bash
# NEVER navigate away from split directory
cd ../api-types             # NO! Original effort directory
cd /workspace/main          # NO! Main repository
cd ../api-types--split-002 # NO! Different split

# NEVER work on files outside split directory
vim /workspace/main/pkg/api/types.go  # NO! Main repo file
cp ../api-types/pkg/api/* ./pkg/      # NO! Copying from wrong place

# NEVER commit from wrong directory
cd .. && git commit         # NO! Wrong directory
```

## Directory Structure During Splits

```
efforts/
├── phase1/
│   └── wave1/
│       ├── api-types/                    # Original too-large effort (ABANDONED)
│       │   ├── SPLIT-INVENTORY.md        # Master list of splits
│       │   ├── SPLIT-PLAN-001.md         # Plan for split 1
│       │   ├── SPLIT-PLAN-002.md         # Plan for split 2
│       │   └── SPLIT-PLAN-003.md         # Plan for split 3
│       │
│       ├── api-types--split-001/         # 🚨 WORK HERE for split 1
│       │   ├── SPLIT-MARKER.txt          # Confirms this is a split
│       │   ├── SPLIT-PLAN-001.md         # Has WORKING_DIRECTORY metadata
│       │   ├── pkg/                      # YOUR IMPLEMENTATION GOES HERE
│       │   │   ├── api/
│       │   │   │   ├── types.go         # Split 1 files only
│       │   │   │   └── helpers.go       # Split 1 files only
│       │   │   └── ...
│       │   └── .git/                     # On branch: phase1/wave1/api-types--split-from--phase1-wave1-api-types-001
│       │
│       ├── api-types--split-002/         # 🚨 WORK HERE for split 2 (AFTER split 1)
│       │   ├── SPLIT-MARKER.txt
│       │   ├── SPLIT-PLAN-002.md         # Has different WORKING_DIRECTORY
│       │   ├── pkg/                      # Different files than split 1
│       │   └── .git/                     # Different branch than split 1
│       │
│       └── api-types--split-003/         # 🚨 WORK HERE for split 3 (AFTER split 2)
│           └── ...
```

## Tool Usage Examples for Splits

### Using Claude Code Tools in Split Directory

```python
# When using Read tool
Read: ./pkg/api/types.go           # Reads from current split directory
Read: SPLIT-PLAN-001.md            # Reads split plan in current directory

# When using Write tool  
Write: ./pkg/api/new_file.go       # Creates in split directory
Write: ./pkg/api/v1/types.go       # Creates subdirectories in split

# When using Edit tool
Edit: ./pkg/api/types.go           # Edits file in split directory

# When using Bash tool
Bash: pwd                          # Should show split directory
Bash: git branch --show-current    # Should show split branch
Bash: go test ./pkg/...            # Tests split code only

# When using Grep tool
Grep: "pattern" ./pkg/             # Searches only in split directory
```

### Verification Commands

Always verify you're in the right place:

```bash
# Check current directory
pwd | grep -- "--split-"            # Should match
echo $?                             # Should be 0

# Check current branch
git branch --show-current | grep -- "--split-from--"  # Should match
echo $?                             # Should be 0

# Check for split marker
test -f SPLIT-MARKER.txt && echo "✅ In split directory" || echo "❌ NOT in split!"

# Check split plan exists
test -f SPLIT-PLAN-*.md && echo "✅ Split plan found" || echo "❌ No split plan!"
```

## Common Mistakes and How to Avoid Them

### Mistake 1: Running Preflight Before Navigation
```bash
# ❌ WRONG
if [[ $(pwd) != */efforts/* ]]; then exit 1; fi  # Fails in wrong directory

# ✅ CORRECT  
# First navigate per R205
cd "$WORKING_DIR"
# THEN run preflight
if [[ $(pwd) != */efforts/* ]]; then exit 1; fi  # Now passes
```

### Mistake 2: Forgetting You're in a Split
```bash
# ❌ WRONG - Trying to access main repo files
vim /workspace/kcp/pkg/apis/...    # Not in split!

# ✅ CORRECT - Work with split files
vim ./pkg/apis/...                 # In split directory
```

### Mistake 3: Switching Between Splits
```bash
# ❌ WRONG - Jumping between splits
cd ../api-types--split-001
# do some work
cd ../api-types--split-002  # NO! Finish split-001 first!

# ✅ CORRECT - Complete one split at a time
# Stay in split-001 until COMPLETELY done
# Commit, push, verify
# ONLY THEN move to split-002
```

## The Split Workflow

1. **Start**: Read SPLIT-INVENTORY.md in too-large directory
2. **Navigate**: Read SPLIT-PLAN-001.md metadata, cd to that directory
3. **Verify**: Check you're in right directory and branch
4. **Work**: Implement EVERYTHING for split-001 in that directory
5. **Complete**: Commit, push, measure size
6. **Next**: Navigate to split-002 directory (repeat from step 2)

## Critical Reminders

- 📁 **One Directory Per Split**: Each split has its own complete workspace
- 🌿 **One Branch Per Split**: Each split has its own git branch
- 📝 **One Plan Per Split**: Each split has its own SPLIT-PLAN-XXX.md
- 🚫 **Never Mix Splits**: Complete one split entirely before starting next
- ✅ **Always Verify Location**: Check pwd and branch before ANY work
- 🔧 **All Tools in Context**: Every command runs in the split directory

## Summary

The split directory specified in the split plan's metadata is your ENTIRE UNIVERSE while working on that split. You:
- Navigate there FIRST (before preflight)
- Work there EXCLUSIVELY
- Complete everything there
- Push from there
- ONLY leave when the split is 100% complete

This ensures clean, mergeable, testable splits that integrate properly.