# Split Implementation Instructions for SW Engineer

## CRITICAL: You Are Implementing a SPLIT, Not the Original Effort

**IMPORTANT**: The original oci-types effort (974 lines) exceeded the 800-line limit. You are now implementing SPLITS of this effort. Each split must be handled as a separate, sequential task.

## Overview
- **Original Effort**: oci-types (974 lines - OVER LIMIT)
- **Split Count**: 2 splits
- **Your Task**: Implement splits SEQUENTIALLY (not in parallel)
- **Size Limit**: Each split MUST be <800 lines

## Implementation Order (MANDATORY)

### ⚠️ CRITICAL SEQUENCING ⚠️
**Splits MUST be implemented in order. DO NOT start Split 002 until Split 001 is merged!**

1. **FIRST**: Implement Split 001 (OCI Package)
2. **WAIT**: For Split 001 to be merged to main
3. **THEN**: Implement Split 002 (Stack Package)

## Split 001: OCI Package (START THIS FIRST)

### Branch Setup
```bash
cd /home/vscode/workspaces/idpbuilder-oci-mgmt
git checkout main
git pull origin main
git checkout -b phase1/wave1/oci-types-split-001
```

### Create Effort Directory
```bash
mkdir -p efforts/phase1/wave1/oci-types-split-001/pkg/oci
cd efforts/phase1/wave1/oci-types-split-001
```

### Files to Implement (622 lines total)
Copy these files from the original effort:
- `pkg/oci/types.go` (121 lines)
- `pkg/oci/manifest.go` (124 lines)
- `pkg/oci/constants.go` (56 lines)
- `pkg/oci/types_test.go` (130 lines)
- `pkg/oci/manifest_test.go` (191 lines)

### Verification Steps
1. **Compile**: `go build ./pkg/oci/...`
2. **Test**: `go test ./pkg/oci/... -cover` (must be >80%)
3. **Measure**: Run line counter - MUST be ≤622 lines
4. **Commit**: Use proper commit message from SPLIT-PLAN-001.md

### DO NOT Include
- ❌ Any files from pkg/stack/
- ❌ pkg/doc.go
- ❌ Any files not listed above

---

## Split 002: Stack Package (ONLY AFTER Split 001 is MERGED)

### Pre-requisite Check
```bash
# MUST verify Split 001 is merged first!
git checkout main
git pull origin main
ls pkg/oci/types.go || echo "ERROR: Split 001 not merged yet!"
```

### Branch Setup
```bash
git checkout main
git pull origin main
git checkout -b phase1/wave1/oci-types-split-002
```

### Create Effort Directory
```bash
mkdir -p efforts/phase1/wave1/oci-types-split-002/pkg/stack
cd efforts/phase1/wave1/oci-types-split-002
```

### Files to Implement (352 lines total)
Copy these files from the original effort:
- `pkg/stack/types.go` (107 lines)
- `pkg/stack/constants.go` (42 lines)
- `pkg/stack/types_test.go` (164 lines)
- `pkg/doc.go` (39 lines)

### Important: Update Imports
The stack package references OCI types. Ensure imports are correct:
```go
import "github.com/idpbuilder/idpbuilder-oci-mgmt/pkg/oci"
```

### Verification Steps
1. **Compile**: `go build ./pkg/stack/...`
2. **Test**: `go test ./pkg/stack/... -cover` (must be >80%)
3. **Measure**: Run line counter - MUST be ≤352 lines
4. **Commit**: Use proper commit message from SPLIT-PLAN-002.md

### DO NOT Include
- ❌ Any files from pkg/oci/ (already in Split 001)
- ❌ Any files not listed above

---

## Line Counter Usage

### Always Measure from Effort Directory
```bash
cd efforts/phase1/wave1/oci-types-split-XXX
/home/vscode/workspaces/idpbuilder-oci-mgmt/tools/line-counter.sh
```

### Expected Output
- Split 001: Must show ≤622 lines
- Split 002: Must show ≤352 lines

### If Over Limit
- STOP immediately
- Do NOT commit
- Report to orchestrator for further split planning

---

## Common Mistakes to Avoid

### ❌ DO NOT:
- Start both splits in parallel
- Include files from both splits in one branch
- Start Split 002 before Split 001 is merged
- Exceed line limits
- Skip tests
- Modify files outside your assigned split

### ✅ DO:
- Implement splits sequentially
- Verify each split compiles independently
- Run all tests with >80% coverage
- Measure size with line-counter.sh
- Use sparse checkout (only assigned files)
- Wait for merge confirmation between splits

---

## Communication Protocol

### After Completing Split 001:
```markdown
STATUS: Split 001 Complete
- Branch: phase1/wave1/oci-types-split-001
- Size: [actual] lines (limit: 622)
- Tests: [coverage]% (required: 80%)
- Compilation: ✅ Successful
- Ready for review
```

### After Review Approval:
```markdown
STATUS: Awaiting Split 001 Merge
- Cannot start Split 002 until merge confirmed
- Blocking on: main branch integration
```

### After Split 001 Merged:
```markdown
STATUS: Starting Split 002
- Split 001 merged to main ✅
- Branch: phase1/wave1/oci-types-split-002
- Beginning implementation
```

---

## Final Checklist

Before marking any split complete:
- [ ] Correct branch name used
- [ ] Only assigned files included
- [ ] Compilation successful
- [ ] Tests pass with >80% coverage
- [ ] Size under limit (verified with line-counter.sh)
- [ ] Proper commit message used
- [ ] No files from other splits included
- [ ] Dependencies handled correctly

## Questions?
If anything is unclear:
1. Check SPLIT-PLAN-XXX.md for your specific split
2. Check SPLIT-INVENTORY.md for overall strategy
3. Report blockers immediately to orchestrator

Remember: Success depends on SEQUENTIAL execution and STRICT adherence to file boundaries!