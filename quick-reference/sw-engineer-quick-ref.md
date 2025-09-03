# ⚡ SW ENGINEER QUICK REFERENCE

## 🚨 STARTUP CHECKLIST
```
□ Print: AGENT STARTUP: $(date '+%Y-%m-%d %H:%M:%S %Z')
□ Verify: Correct directory and git branch  
□ Load: agent-states/sw-engineer/{STATE}/rules.md
□ Read: ${WORKING_DIR}/IMPLEMENTATION-PLAN.md
□ Read: TEST-DRIVEN-VALIDATION-REQUIREMENTS.md
```

## ⚡ CRITICAL METRICS (GRADING)
```
┌─────────────────────────────────────────────┐
│ 🚨 SPEED: >50 lines/hour implementation    │
│ 🚨 SIZE: NEVER exceed 800 lines total      │
│ 🚨 TEST COVERAGE: Meet phase minimums      │
│ 🚨 WORK LOG: Update every checkpoint       │
└─────────────────────────────────────────────┘
```

## 🔄 STATE MACHINE FLOW
```
INIT → IMPLEMENTATION → MEASURE_SIZE
  ↓
MEASURE_SIZE → [COMPLETE|FIX_ISSUES|SPLIT_WORK]
  ↓
FIX_ISSUES → CODE_REVIEW → IMPLEMENTATION
  ↓  
SPLIT_WORK → IMPLEMENTATION (for each split)
```

## 🎯 KEY ACTIONS BY STATE

| State | Action | Critical Rule |
|-------|--------|---------------|
| `IMPLEMENTATION` | 💻 Write code >50 lines/hr | R152 |
| `MEASURE_SIZE` | 📏 Use tmc-pr-line-counter.sh | R000 |
| `FIX_ISSUES` | 🔧 Address review feedback | R152 |
| `TEST_WRITING` | 🧪 Meet coverage minimums | R152 |

## 🛠️ ESSENTIAL COMMANDS

### Size Measurement (MANDATORY TOOL!)
```bash
# ONLY way to measure lines - NEVER count manually
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c ${BRANCH}

# With detailed breakdown if approaching limit
/workspaces/kcp-shared-tools/tmc-pr-line-counter.sh -c ${BRANCH} -d

# Check every ~200 lines of implementation
```

### Work Log Updates
```bash
# Update progress regularly
echo "## $(date '+%Y-%m-%d %H:%M'): Implemented X feature" >> work-log.md
echo "- Lines added: ~150" >> work-log.md
echo "- Test coverage: 85%" >> work-log.md
echo "- Size check: 420/800 lines" >> work-log.md
```

### Git Hygiene
```bash
# Logical, small commits
git add specific-feature.go specific-feature_test.go
git commit -m "feat: add specific feature with tests

- Implements core functionality
- Adds comprehensive test coverage
- Maintains KCP patterns"
```

## 🧪 TEST COVERAGE REQUIREMENTS

| Phase | Minimum Coverage | Critical Areas |
|-------|------------------|----------------|
| Phase 1 | 70% | Core APIs, Controllers |
| Phase 2 | 75% | Webhooks, Validation |
| Phase 3 | 80% | Integration, E2E |

```go
// Always include table-driven tests
func TestMyFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    MyInput
        expected MyOutput
    }{
        // Test cases...
    }
    // Run tests...
}
```

## 📏 SIZE MONITORING PROTOCOL

### Every 200 Lines
```bash
lines=$(tmc-pr-line-counter.sh -c ${BRANCH} | grep "Total:" | awk '{print $2}')
if [ $lines -gt 600 ]; then
    echo "⚠️ APPROACHING LIMIT: $lines/800 lines"
    # Plan completion strategy
fi
if [ $lines -gt 800 ]; then
    echo "🚨 LIMIT EXCEEDED - STOP IMMEDIATELY"
    # Request split from Code Reviewer
fi
```

## ❌ NEVER DO / ✅ ALWAYS DO

❌ **NEVER**:
- Count lines manually (use tmc-pr-line-counter.sh)
- Exceed 800 line limit
- Skip test coverage requirements
- Forget work log updates
- Implement without plan

✅ **ALWAYS**:
- Use KCP patterns and conventions
- Update work-log.md at checkpoints
- Measure size every ~200 lines
- Write tests alongside implementation
- Follow implementation plan exactly

## 🏗️ KCP IMPLEMENTATION PATTERNS

### Multi-Tenant Controller
```go
// Always check workspace context
func (r *MyReconciler) Reconcile(ctx context.Context, req ctrl.Request) {
    cluster := kcplogicalcluster.From(ctx)
    // Implementation with cluster context
}
```

### API Types
```go
// Include workspace fields
type MySpec struct {
    // Core fields
    WorkspaceRef corev1alpha1.WorkspaceReference `json:"workspaceRef"`
}
```

## 🚨 EMERGENCY PROCEDURES

### Size Limit Exceeded
```
1. STOP coding immediately
2. Measure: tmc-pr-line-counter.sh -c ${BRANCH} -d  
3. Commit current work
4. Report to orchestrator: "NEEDS_SPLIT"
5. Wait for Code Reviewer split plan
```

### Test Coverage Below Minimum
```
1. Identify uncovered areas: go test -cover
2. Prioritize critical paths
3. Add table-driven tests
4. Update work log with coverage %
```

### Review Feedback Received
```
1. Read: ${WORKING_DIR}/REVIEW-FEEDBACK.md
2. Address each issue systematically  
3. Update work-log.md with fixes
4. Re-measure size after fixes
```

## 🎓 GRADING FORMULA
```
score = (lines_per_hour/50) * 0.3 +
        (test_coverage/required) * 0.3 +
        (1 if under_limit else 0) * 0.2 +
        (work_log_frequency) * 0.1 +
        (commit_quality) * 0.1

PASS: score >= 0.8
FAIL: score < 0.8 → Warning → Retraining → Termination
```

## 📋 WORK LOG TEMPLATE
```markdown
# Implementation Log - [Effort Name]

## Progress Tracking
- Started: 2025-08-23 14:30
- Current size: 420/800 lines  
- Test coverage: 85%
- Phase: IMPLEMENTATION

## Completed Features
- [x] Core API types (150 lines, 90% coverage)
- [x] Controller logic (200 lines, 80% coverage)
- [ ] Webhook validation (planned: 100 lines)

## Size History
| Time | Lines | Notes |
|------|-------|-------|
| 14:30 | 150   | API types complete |
| 15:45 | 350   | Controller added |
| 16:30 | 420   | Tests added |
```

---
**References**: R152 (implementation speed), R000 (line counting), TEST-DRIVEN-VALIDATION-REQUIREMENTS.md