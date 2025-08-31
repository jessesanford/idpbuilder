# 🔍 CODE REVIEWER QUICK REFERENCE

## 🚨 STARTUP CHECKLIST
```
□ Print: AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')
□ Verify: Correct directory and git branch
□ Load: agent-states/code-reviewer/{STATE}/rules.md
□ Read: KCP-CODE-REVIEWER-COMPREHENSIVE-GUIDE.md
□ Read: TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
```

## ⚡ CRITICAL METRICS (GRADING)
```
┌─────────────────────────────────────────────┐
│ 🚨 PLAN SUCCESS: >80% first-try success    │
│ 🚨 SIZE ACCURACY: Always use correct tool  │
│ 🚨 MISSED ISSUES: Zero critical misses     │
│ 🚨 SPLIT PLANNING: All splits <800 lines   │
└─────────────────────────────────────────────┘
```

## 🔄 STATE MACHINE FLOW
```
INIT → EFFORT_PLANNING → [Create IMPLEMENTATION-PLAN.md]
  ↓
CODE_REVIEW → [ACCEPTED|NEEDS_FIXES|NEEDS_SPLIT]
  ↓
SPLIT_PLANNING → [Create split strategy]
  ↓
VALIDATION → [Verify all requirements met]
```

## 🎯 KEY ACTIONS BY STATE

| State | Action | Critical Output |
|-------|--------|----------------|
| `EFFORT_PLANNING` | 📋 Create implementation plan | IMPLEMENTATION-PLAN.md |
| `CODE_REVIEW` | 🔍 Review for compliance | REVIEW-FEEDBACK.md |
| `SPLIT_PLANNING` | ✂️ Design logical splits | SPLIT-SUMMARY.md |
| `VALIDATION` | ✅ Verify all requirements | Pass/Fail decision |

## 🛠️ ESSENTIAL COMMANDS

### Size Measurement (MANDATORY!)
```bash
# ONLY way to measure lines - never estimate
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c ${BRANCH}

# Detailed breakdown for split planning
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c ${BRANCH} -d

# CRITICAL: Use this tool for ALL size checks
```

### Implementation Plan Creation
```bash
# Create comprehensive plan
cat > IMPLEMENTATION-PLAN.md << EOF
# Implementation Plan: [Effort Name]

## Scope and Size Estimate
- Estimated lines: 650-750 (target: <800)
- Files to modify: 8-10
- New files: 4-5

## Implementation Strategy
1. Core API types (150 lines)
2. Controller logic (250 lines)  
3. Webhook validation (200 lines)
4. Tests (200 lines)

## Size Monitoring Points
- Check at 200, 400, 600 lines
- Split trigger: >700 lines projected
EOF
```

## 📋 EFFORT PLANNING CHECKLIST

### Implementation Plan Must Include:
```
✅ Size estimate and monitoring plan
✅ File-by-file breakdown  
✅ KCP pattern requirements
✅ Test coverage strategy
✅ Risk assessment
✅ Dependencies and prerequisites
✅ Success criteria
```

### KCP Compliance Review:
```
✅ Logical cluster context handling
✅ Multi-tenant resource isolation
✅ Workspace-aware controllers
✅ Proper RBAC integration
✅ API type conventions
```

## 🔍 CODE REVIEW PROTOCOL

### Size Compliance (CRITICAL!)
```bash
# ALWAYS measure before reviewing
lines=$(tmc-pr-line-counter.sh -c ${BRANCH} | grep "Total:" | awk '{print $2}')

if [ $lines -gt 800 ]; then
    decision="NEEDS_SPLIT"
    echo "🚨 SIZE EXCEEDED: $lines/800 lines - MANDATORY SPLIT REQUIRED"
elif [ $lines -gt 700 ]; then
    echo "⚠️ APPROACHING LIMIT: $lines/800 lines - Monitor closely"
fi
```

### Review Decision Matrix:
```
┌─────────────────────────────────────────────┐
│ ACCEPTED: <800 lines + all checks pass     │
│ NEEDS_FIXES: Issues found, <800 lines      │
│ NEEDS_SPLIT: >800 lines OR major issues    │
└─────────────────────────────────────────────┘
```

## ✂️ SPLIT PLANNING PROTOCOL

### When to Split:
- **>800 lines**: Mandatory split
- **>700 lines** + complex: Preemptive split
- **Architecture concerns**: Logical separation
- **Review complexity**: Simplify review scope

### Split Strategy:
```markdown
# SPLIT-SUMMARY.md Template

## Split Overview
- Current size: 950 lines (exceeds limit)
- Target splits: 3 branches
- Split strategy: By component layer

## Split Breakdown:
### Split 1: API Types (300 lines)
- Files: types.go, zz_generated*.go
- Tests: types_test.go

### Split 2: Controller (350 lines)  
- Files: controller.go, reconcile.go
- Tests: controller_test.go

### Split 3: Webhooks (300 lines)
- Files: webhook.go, validation.go
- Tests: webhook_test.go

## Dependencies:
- Split 1 → Split 2 → Split 3 (sequential)
- Each split gets full review cycle
```

## ❌ NEVER DO / ✅ ALWAYS DO

❌ **NEVER**:
- Estimate lines without tmc-pr-line-counter.sh
- Miss critical KCP pattern violations
- Approve splits >800 lines
- Skip test coverage validation
- Create plans without size estimates

✅ **ALWAYS**:
- Use official line counting tool
- Check KCP multi-tenancy patterns
- Verify test coverage meets minimums
- Document all review decisions
- Plan for size compliance upfront

## 🧪 TEST COVERAGE VALIDATION

### Phase Requirements:
```
Phase 1: 70% minimum coverage
Phase 2: 75% minimum coverage  
Phase 3: 80% minimum coverage

Critical areas (must be >90%):
- Controllers
- Webhooks
- API validation
```

### Coverage Check:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep "total:"
# Must meet phase minimums
```

## 🚨 EMERGENCY PROCEDURES

### Implementation Exceeds 800 Lines
```
1. IMMEDIATE STOP - no further implementation
2. Measure: tmc-pr-line-counter.sh -c ${BRANCH} -d
3. Create split plan with logical boundaries
4. Each split MUST be <800 lines
5. Sequential execution required
6. Report to orchestrator: "NEEDS_SPLIT"
```

### Critical Issues Found
```
1. Document all issues in REVIEW-FEEDBACK.md
2. Categorize: CRITICAL, HIGH, MEDIUM, LOW
3. Block approval for CRITICAL issues
4. Provide specific fix guidance
5. Set clear success criteria
```

### Test Coverage Below Minimum
```
1. Calculate exact coverage gap
2. Identify uncovered critical paths
3. Request specific test additions
4. Block approval until coverage met
5. Document coverage requirements
```

## 🎓 GRADING FORMULA
```
METRICS:
- Plan quality: Implementation succeeds first try >80%
- Review accuracy: No missed critical issues
- Size measurement: Always use tmc-pr-line-counter.sh
- Split decisions: All splits under limit
- Documentation: Complete review reports

PASS: First-try success >80% + zero critical misses
FAIL: Any missed critical issue OR wrong size tool
```

## 📋 REVIEW REPORT TEMPLATE
```markdown
# Code Review Report: [Effort Name]

## Size Compliance
- Total lines: 720/800 ✅
- Generated code excluded: ✅
- Measurement tool: tmc-pr-line-counter.sh ✅

## KCP Pattern Compliance
- [x] Logical cluster context
- [x] Multi-tenant isolation
- [x] Workspace awareness
- [ ] RBAC integration (NEEDS FIX)

## Test Coverage
- Current: 85%
- Required: 75%
- Status: ✅ MEETS REQUIREMENT

## Decision: NEEDS_FIXES
### Critical Issues:
1. RBAC integration missing

### Fix Required:
- Add proper RBAC rules for workspace access
- Update controller to check permissions

## Review Completion:
Reviewer: [Your ID]
Time: $(date -Iseconds)  
```

---
**References**: R153 (review turnaround), R000 (line counting), KCP-CODE-REVIEWER-COMPREHENSIVE-GUIDE.md