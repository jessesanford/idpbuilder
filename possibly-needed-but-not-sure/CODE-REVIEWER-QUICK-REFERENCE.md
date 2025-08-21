# Code Reviewer Quick Reference Card

## Review Order (ALWAYS)
1. **Size** → 2. **Structure** → 3. **Quality** → 4. **Architecture** → 5. **Commits** → 6. **Tests** → 7. **Build/Lint**

## Size Measurement (ONLY METHOD)
```bash
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c ${BRANCH_NAME}
```
- Limit: 800 lines
- Exception: Only if split breaks logic
- Never measure any other way

## Phase Rules Enforcement

| Phase | Allowed | Forbidden |
|-------|---------|-----------|
| 1 | APIs, Types, CRDs | Implementations |
| 2 | Controllers | Syncer |
| 3 | Syncer (phase7 ONLY) | Duplicate syncers |
| 4 | Features | Skip cross-workspace fix |
| 5 | Tests | < 95% coverage |

## Critical Checkpoints

### KCP Style
```go
package names     // lowercase_no_underscores
TypeNames        // CamelCase
functionNames    // camelCase (private)
FunctionNames    // CamelCase (public)  
variableNames    // camelCase
CONSTANTS        // UPPER_SNAKE or CamelCase
```

### No Hardcoded Values
```go
// ❌ BAD
timeout := 30
retries := 5

// ✅ GOOD  
const DefaultTimeout = 30
const MaxRetries = 5
```

### Documentation Required
- Every exported type
- Every exported function
- Every exported const/var
- Package level doc.go

### Decomposition Limits
- Function: 50 lines max
- File: 500 lines max
- Package: 10 files max
- Complexity: 10 cyclomatic

## Commit Requirements

### Message Format
```
type(scope): description

- feat: new feature
- fix: bug fix
- docs: documentation
- test: tests only
- refactor: code restructure
- chore: maintenance
```

### Commit Size
- Preferred: < 200 lines
- Warning: > 500 lines
- Must be atomic

### History
- Linear (no merges)
- Tells coherent story
- APIs → Impl → Tests order

## Test Coverage by Phase

| Phase | Required | Critical Areas |
|-------|----------|----------------|
| 1 | 80% | Type validation |
| 2 | 85% | Controller logic |
| 3 | 90% | Syncer (CRITICAL) |
| 4 | 85% | Feature completeness |
| 5 | 95% | Integration tests |

## Build/Lint Requirements

### Must Pass
- `go build ./...` - No warnings
- `go test ./... -race` - No races
- `golangci-lint run` - No errors
- `make generate` - No uncommitted

### Critical Linters
- gofmt
- goimports
- govet
- gosec (security)
- ineffassign
- misspell

## Review Decision Tree

```
Size Check
├─ ≤800 lines → Continue
└─ >800 lines
   ├─ Can grant exception? → Document → Continue
   └─ Cannot grant → Create split plan → STOP

Structure Check  
├─ Correct phase? → Continue
└─ Wrong phase → REJECT

Quality Check
├─ All pass? → Continue
└─ Issues? → Document fixes → NEEDS_FIXES

Architecture Check
├─ No duplicates? → Continue
└─ Duplicates? → REJECT

Commit Check
├─ Linear & atomic? → Continue
└─ Issues? → Request rebase → NEEDS_FIXES

Test Check
├─ Coverage met? → Continue
└─ Below minimum? → NEEDS_FIXES

Build/Lint Check
├─ All pass? → ACCEPT
└─ Failures? → NEEDS_FIXES
```

## Output Templates

### Quick Accept
```markdown
✅ ACCEPTED
- Size: XXX lines ✅
- Phase compliance ✅  
- Coverage: XX% ✅
- Ready for PR
```

### Quick Reject
```markdown
❌ NEEDS FIXES
Critical:
1. [Issue] at file:line
2. [Issue] at file:line

Required: [List fixes]
```

### Quick Split
```markdown
🔄 NEEDS SPLIT
- Current: XXXX lines
- Plan: 3 parts
  1. APIs (XXX lines)
  2. Impl (XXX lines)
  3. Tests (XXX lines)
```

## Common Issues Checklist

- [ ] Size > 800 without exception
- [ ] Implementation in API phase
- [ ] Duplicate SyncEngine
- [ ] Coverage below requirement
- [ ] Hardcoded values
- [ ] Missing godoc
- [ ] Merge commits in history
- [ ] Wrong base branch
- [ ] Giant function (>50 lines)
- [ ] Not using phase7 syncer

## Key Branches to Remember

- Phase 1 APIs: Use main as base
- Phase 2 Controllers: Use phase1-integration
- Phase 3 Syncer: Use phase2-integration, cherry-pick from phase7
- Phase 4 Features: Use phase3-integration, fix cross-workspace first
- Phase 5 Tests: Use phase4-integration

## Exception Criteria

Grant exception ONLY if:
1. Complex state machine
2. Atomic transaction required
3. Generated code
4. Would duplicate >30% code to split

Never grant if:
1. Can separate APIs/Impl/Tests
2. Contains unrelated features
3. Simple concatenation

## Remember

- **Measure**: Only use tmc-pr-line-counter.sh
- **Enforce**: Phase ordering strictly
- **Document**: Every issue with fix
- **Quality**: Better to fix now than later
- **Linear**: No merge commits except base