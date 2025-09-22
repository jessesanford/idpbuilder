# 🔴🔴🔴 RULE R326: SPLIT FILE PLACEMENT PROTOCOL 🔴🔴🔴

**FILE**: `$CLAUDE_PROJECT_DIR/rule-library/R326-split-file-placement-protocol.md`
**CRITICALITY**: SUPREME LAW - Files in wrong directories = -100% measurement failure
**VIOLATIONS**: Split-001 creating split-001/pkg/ subdirectory = CATASTROPHIC FAILURE

## 🔴🔴🔴 THE CRITICAL BUG THAT BREAKS EVERYTHING 🔴🔴🔴

### ACTUAL FAILURE FROM TRANSCRIPT:
```
Split-001 created: split-001/pkg/registry/auth.go
Split-002 created: pkg/registry/auth.go
Git sees these as DIFFERENT files!
Result: 1989 lines measured instead of ~400!
```

**THE PROBLEM**: SW Engineers create subdirectories matching the split name!
**THE CONSEQUENCE**: Massive measurement errors, splits appear 3-10X oversized!

## 🔴🔴🔴 ABSOLUTE RULE: NO SPLIT SUBDIRECTORIES 🔴🔴🔴

### ❌ WRONG - CATASTROPHIC FAILURE:
```
Working directory: efforts/phase1/wave1/registry-SPLIT-001/
Files created:
└── split-001/              ❌❌❌ NO! NEVER CREATE THIS!
    └── pkg/
        └── registry/
            └── auth.go
```

### ✅ CORRECT - STANDARD PROJECT STRUCTURE:
```
Working directory: efforts/phase1/wave1/registry-SPLIT-001/
Files created:
├── pkg/                    ✅ Standard directory
│   └── registry/
│       └── auth.go
├── cmd/                    ✅ Standard directory
│   └── main.go
└── tests/                  ✅ Standard directory
    └── registry_test.go
```

## 🔴🔴🔴 ENFORCEMENT FOR SW ENGINEERS 🔴🔴🔴

### BEFORE CREATING ANY FILE:
```bash
# MANDATORY VERIFICATION
echo "🔴 R326: VERIFYING DIRECTORY STRUCTURE..."
pwd
# Should show: .../efforts/phaseX/waveY/effort-SPLIT-XXX/

echo "❌ I WILL NOT create split-001/ subdirectory"
echo "✅ I WILL create files in standard locations: pkg/, cmd/, tests/, etc."

# FORBIDDEN CHECK
if [ -d "split-*" ]; then
    echo "🔴🔴🔴 FATAL: Split subdirectory detected!"
    echo "DELETE IT IMMEDIATELY!"
    exit 326
fi
```

### DURING IMPLEMENTATION:
```bash
# When creating a file, ALWAYS use standard paths:
mkdir -p pkg/registry
touch pkg/registry/auth.go     # ✅ CORRECT

# NEVER do this:
mkdir -p split-001/pkg/registry  # ❌ WRONG!
touch split-001/pkg/registry/auth.go  # ❌ CATASTROPHIC!
```

### BEFORE COMMITTING:
```bash
# MANDATORY CHECK
echo "🔴 R326: Checking for split subdirectories..."
if find . -type d -name "split-*" | grep -q .; then
    echo "🔴🔴🔴 FATAL: Split subdirectory found!"
    find . -type d -name "split-*"
    echo "DELETE THESE DIRECTORIES BEFORE COMMITTING!"
    exit 326
fi

# Verify files are in standard locations
echo "✅ Files in standard directories:"
ls -la pkg/ cmd/ tests/ 2>/dev/null
```

## 🔴🔴🔴 ENFORCEMENT FOR CODE REVIEWERS 🔴🔴🔴

### IN SPLIT PLAN TEMPLATES:
```markdown
## 🔴🔴🔴 CRITICAL: FILE LOCATIONS (R326) 🔴🔴🔴

⚠️⚠️⚠️ FILES MUST GO DIRECTLY IN STANDARD DIRECTORIES ⚠️⚠️⚠️
❌ NEVER create split-001/ subdirectory!
✅ Use standard project structure: pkg/, cmd/, tests/

### CORRECT File Locations for This Split:
- `pkg/registry/auth.go` - Authentication logic
- `pkg/registry/client.go` - Client implementation  
- `tests/registry_test.go` - Unit tests

### ❌ WRONG (CATASTROPHIC FAILURE):
- `split-001/pkg/registry/auth.go` ← NEVER DO THIS!
- `split-001/cmd/main.go` ← AUTOMATIC FAILURE!

### Directory Structure After This Split:
```
. (working directory: efforts/phase1/wave1/registry-SPLIT-001/)
├── pkg/
│   └── registry/
│       ├── auth.go      ✅ CORRECT LOCATION
│       └── client.go    ✅ CORRECT LOCATION
└── tests/
    └── registry_test.go ✅ CORRECT LOCATION

NOT THIS:
└── split-001/           ❌❌❌ NEVER CREATE THIS!
    └── pkg/
        └── registry/
```
```

## 🔴🔴🔴 ENFORCEMENT FOR ORCHESTRATORS 🔴🔴🔴

### WHEN SPAWNING SW ENGINEERS FOR SPLITS:
```markdown
## 🔴 CRITICAL WARNING (R326): FILE PLACEMENT

**DO NOT CREATE split-XXX/ SUBDIRECTORIES!**
- ❌ WRONG: split-001/pkg/registry/auth.go
- ✅ RIGHT: pkg/registry/auth.go

Files go DIRECTLY in standard project directories (pkg/, cmd/, tests/).
Creating split subdirectories causes MASSIVE measurement errors!
```

## 🔴🔴🔴 WHY THIS MATTERS 🔴🔴🔴

### Measurement Impact:
```bash
# With split subdirectory (WRONG):
git diff --name-only split-001 split-002
split-001/pkg/registry/auth.go    # Split-001's file
pkg/registry/auth.go              # Split-002's file
# Git sees 2 DIFFERENT files = counts BOTH!

# Without split subdirectory (CORRECT):
git diff --name-only split-001 split-002
pkg/registry/auth.go              # Modified by Split-002
# Git sees 1 file modified = counts only changes!
```

### The Catastrophic Result:
- Split-001: Creates 400 lines in split-001/pkg/
- Split-002: Creates 400 lines in pkg/ (different location!)
- Split-003: Creates 400 lines in pkg/ (same as Split-002)
- **MEASUREMENT**: 1989 lines (sees all as different files!)
- **ACTUAL**: Should be ~400 lines per split

## 🔴🔴🔴 VALIDATION COMMANDS 🔴🔴🔴

### For SW Engineers:
```bash
validate_no_split_dirs() {
    if find . -type d -name "split-*" | grep -q .; then
        echo "🔴🔴🔴 R326 VIOLATION: Split subdirectories found!"
        find . -type d -name "split-*"
        return 326
    fi
    echo "✅ R326: No split subdirectories - CORRECT"
    return 0
}
```

### For Code Reviewers:
```bash
# After reviewing split implementation
if git ls-tree -r HEAD | grep "^[0-9]* blob.*split-[0-9]*/"; then
    echo "🔴🔴🔴 R326 VIOLATION: Files in split subdirectory!"
    echo "REJECT THIS SPLIT - Files in wrong location!"
    exit 326
fi
```

## PENALTIES

- Creating split-XXX/ subdirectory: **-100% IMMEDIATE FAILURE**
- Files in wrong location: **SPLIT REJECTED, MUST RE-IMPLEMENT**
- Measurement errors from wrong placement: **ENTIRE EFFORT FAILS**

## REMEMBER

**SPLITS USE THE SAME DIRECTORY STRUCTURE AS THE MAIN PROJECT!**
- Each split is in its own working directory (effort-SPLIT-XXX/)
- But files go in STANDARD locations (pkg/, cmd/, tests/)
- NEVER create subdirectories matching the split name!

---
**Rule Created**: 2025-01-08
**Reason**: Split-001 created split-001/pkg/ causing 1989-line measurement error
**Impact**: CATASTROPHIC - Breaks entire measurement system