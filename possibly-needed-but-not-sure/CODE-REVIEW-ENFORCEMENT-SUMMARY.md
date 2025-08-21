# Code Review Enforcement Summary

## Overview
This document summarizes the comprehensive code review system for ensuring only production-ready, KCP-compliant code enters the TMC implementation.

## Key Enforcement Points

### 1. Size Measurement (UNIVERSAL)
```bash
# ONLY METHOD - Never measure any other way
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c ${BRANCH_NAME}
```
- Limit: 800 lines
- Exceptions: Only reviewer can grant
- Measurement: Universal across project

### 2. Phase/Wave/Effort Ordering (STRICT)

| Phase | Must Contain | Must NOT Contain |
|-------|--------------|------------------|
| 1 | APIs, Types, CRDs | Any implementation |
| 2 | Controllers | Syncer |
| 3 | Syncer (phase7 only) | Duplicate syncers |
| 4 | Features | Skip cross-workspace fix |
| 5 | Tests | < 95% coverage |

### 3. Code Quality Standards

#### KCP Style Compliance
- Package naming: `lowercase` (no underscores)
- Types: `CamelCase` for exported
- Functions: `camelCase` (private), `CamelCase` (public)
- Documentation: All exports must have godoc

#### No Hardcoded Values
```go
// Every magic number/string must be a constant
const DefaultTimeout = 30 * time.Second
const MaxRetries = 5
```

#### Decomposition Rules
- Functions: 50 lines max
- Files: 500 lines max
- Cyclomatic complexity: 10 max

### 4. Commit Requirements

#### Format
```
type(scope): description

Types: feat, fix, docs, test, refactor, chore
```

#### Structure
- Linear history (no merges except base)
- Atomic commits (< 200 lines preferred)
- Tell coherent story: APIs → Implementation → Tests

### 5. Testing Requirements by Phase

| Phase | Coverage | Critical Areas |
|-------|----------|----------------|
| 1 | 80% | Type validation |
| 2 | 85% | Controller logic |
| 3 | 90% | Syncer (CRITICAL) |
| 4 | 85% | Feature completeness |
| 5 | 95% | Full integration |

### 6. Build/Lint Requirements

Must pass ALL:
- `go build ./...` - No warnings
- `go test ./... -race` - No data races
- `golangci-lint run ./...` - No errors
- `make generate` - No uncommitted changes

## Review Decision Flow

```
1. SIZE CHECK
   ├─ ≤800 lines → Continue
   └─ >800 lines
      ├─ Exception warranted? → Document → Continue
      └─ No exception → Create split plan → STOP

2. PHASE COMPLIANCE
   ├─ Correct phase content → Continue
   └─ Wrong phase → REJECT

3. CODE QUALITY
   ├─ All standards met → Continue
   └─ Issues found → Document fixes → NEEDS FIXES

4. ARCHITECTURE
   ├─ No duplicates, uses existing → Continue
   └─ Duplicates or reinvents → REJECT

5. COMMITS
   ├─ Linear, atomic, well-described → Continue
   └─ Issues → Request rebase → NEEDS FIXES

6. TESTING
   ├─ Coverage met → Continue
   └─ Below requirement → NEEDS FIXES

7. BUILD/LINT
   ├─ All pass → ACCEPT ✅
   └─ Failures → NEEDS FIXES
```

## Exception Criteria (Reviewer Authority Only)

### Can Grant Exception If:
1. Complex state machine that cannot be decomposed
2. Atomic transaction spanning components
3. Generated code (protobuf, deepcopy)
4. Would require >30% code duplication to split

### Never Grant If:
1. Can separate APIs/Implementation/Tests
2. Contains unrelated features
3. Simple concatenation of functions

### Exception Documentation Required:
```yaml
exception_granted: true
reason: "Complex bidirectional sync state machine"
justification: |
  7-state machine with shared transitions.
  Splitting would duplicate 400+ lines.
risk_mitigation:
  - 2 additional reviewers required
  - Extended test suite required
  - Performance benchmarks required
```

## Review Output Standards

### Accepted
```markdown
✅ ACCEPTED
- Size: XXX lines ✅
- Phase compliance ✅
- Coverage: XX% ✅
- Ready for PR
```

### Needs Fixes
```markdown
❌ NEEDS FIXES
Critical Issues:
1. [Specific issue] at file:line
   Fix: [Exact fix required]
2. [Issue 2]...

Required Actions:
- [ ] Fix hardcoded values
- [ ] Increase coverage to XX%
- [ ] Rebase to remove merge commits
```

### Needs Split
```markdown
🔄 NEEDS SPLIT
Current: XXXX lines
Plan:
1. APIs (~XXX lines)
2. Implementation (~XXX lines)
3. Tests (~XXX lines)

[Detailed split instructions]
```

## Common Rejection Reasons

1. **Wrong Phase** - Implementation in API phase
2. **Size Violation** - >800 lines without valid exception
3. **Duplicate Code** - Reimplemented existing functionality
4. **Poor Testing** - Below coverage requirement
5. **Bad Commits** - Non-atomic, merge commits, poor messages
6. **Style Violations** - Not matching KCP conventions
7. **Hardcoded Values** - Magic numbers/strings
8. **Missing Documentation** - Exported symbols undocumented
9. **Giant Functions** - Functions >50 lines
10. **Wrong Base Branch** - Not based on correct integration

## Enforcement Hierarchy

1. **Orchestrator** - Ensures reviews happen
2. **Reviewer** - Enforces all standards
3. **Developer Agent** - Follows phase plans
4. **Split Protocol** - Handles oversized branches

## Key Reminders

- **Measurement**: ONLY use tmc-pr-line-counter.sh
- **Phases**: Strictly enforce ordering
- **Quality**: No compromises on standards
- **Exceptions**: Only reviewer can grant
- **Documentation**: Every issue needs specific fix
- **Continuous**: Don't stop during splits
- **Linear**: No merge commits in history

## Success Metrics

When all efforts complete:
- 100% pass code review
- 0 branches >800 lines (without exception)
- 100% meet coverage requirements
- 0 duplicate implementations
- 100% KCP style compliant
- 100% have linear history

This comprehensive review system ensures only clean, maintainable, production-ready code enters the KCP+TMC project.