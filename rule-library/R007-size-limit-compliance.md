# 🚨 RULE R007 - Size Limit Compliance (800 Line Maximum)

**Criticality:** BLOCKING - Cannot proceed if violated  
**Grading Impact:** -20% per violation, -100% if merged  
**Enforcement:** CONTINUOUS - Check every 200 lines added

## Rule Statement

NO effort may EVER exceed 800 lines of IMPLEMENTATION CODE. Soft warning at 700 lines. Automatic split required at violation.

**CRITICAL**: Line counts ONLY include critical path implementation files. Tests, demos, docs, configs, and generated code do NOT count toward the limit.

### 🔴🔴🔴 PARAMOUNT REQUIREMENT (Per R307) 🔴🔴🔴
**EVERY effort AND split must be independently mergeable!**
- Must compile and pass tests when merged alone
- Must NOT break existing functionality
- Must use feature flags for incomplete features
- Must be mergeable months/years later

## Size Limits

| Threshold | Action Required | Grading Impact |
|-----------|----------------|----------------|
| 0-600 lines | Continue normally | None |
| 600-700 lines | Heightened monitoring | None |
| 700-750 lines | Warning - prepare to split | -5% if ignored |
| 750-800 lines | Critical - must plan split | -10% if delayed |
| >800 lines | VIOLATION - must split NOW | -20% immediate |
| >800 merged | CATASTROPHIC FAILURE | -100% FAIL |

## Measurement Requirements

### 🔴🔴🔴 CRITICAL: ONLY IMPLEMENTATION CODE COUNTS 🔴🔴🔴
**The 800-line limit applies ONLY to critical path implementation code!**
- Tests do NOT count (write as many as needed!)
- Demos do NOT count (R330 demo requirements are separate)
- Documentation does NOT count
- Configuration files do NOT count
- Generated code does NOT count

### 🚨 CRITICAL: Check Splits First (R297)
**BEFORE measuring ANY effort, check if it was already split!**
- Check `split_count` in orchestrator-state.json
- If > 0: Effort is COMPLIANT (already split)
- Measure ORIGINAL effort branches, NOT integration branches
- Integration branches merge all splits (will exceed limits - EXPECTED)

### MANDATORY: Use Official Tool
```bash
# ONLY valid measurement method
$CLAUDE_PROJECT_DIR/tools/line-counter.sh -c <branch>

# NEVER use:
# - wc -l
# - git diff --stat
# - manual counting
# - find | xargs wc -l
```

### Measurement Frequency
- Every 200 lines of new code
- Before EVERY commit
- After EVERY major feature
- Before requesting review
- After addressing review comments

## Violation Response Protocol

### When >800 Lines Detected

```bash
# IMMEDIATE ACTIONS REQUIRED:
handle_size_violation() {
    local branch="$1"
    local line_count="$2"
    
    echo "🚨🚨🚨 SIZE VIOLATION DETECTED! 🚨🚨🚨"
    echo "Branch: $branch"
    echo "Lines: $line_count (limit: 800)"
    echo ""
    echo "MANDATORY ACTIONS:"
    echo "1. STOP all implementation immediately"
    echo "2. Commit current work"
    echo "3. Create split plan"
    echo "4. Execute split protocol"
    
    # Orchestrator must:
    # 1. Update state to MEASURE_SIZE_VIOLATION
    # 2. Spawn Code Reviewer for split planning
    # 3. Create split infrastructure
    # 4. Resume with split implementation
}
```

### Split Planning Requirements

When splitting is required:
1. Code Reviewer creates SPLIT-PLAN.md
2. Identify logical boundaries
3. Each split must be <700 lines
4. Splits execute SEQUENTIALLY
5. Each split gets full review

### 🔴 Split Independence Requirements (R307)
**CRITICAL**: Each split must be independently mergeable!
1. **Split 1** must work alone (even if Split 2 never merges)
2. **Split 2** must work with just Split 1 merged
3. **No split** can break existing functionality
4. **Feature flags** hide incomplete features across splits
5. **Example**: Authentication split sequence
   - Split 1: Basic auth interface + stub (works alone)
   - Split 2: OAuth implementation (enhances Split 1)
   - Split 3: Advanced features (optional enhancement)

## Common Violations

### VIOLATION: Ignoring Warning
```bash
# Line count: 720
# ❌ WRONG: Continuing to add code
echo "Just a few more functions..."
# Adds 100 more lines → 820 total
```

### VIOLATION: Wrong Measurement
```bash
# ❌ WRONG: Using wrong tool
wc -l *.go
# Shows 650 lines (incorrect)

# Reality: line-counter.sh shows 850 lines
```

### VIOLATION: Merging Oversized
```bash
# ❌ CATASTROPHIC: Merging >800 lines
git merge effort-branch  # 850 lines
# IMMEDIATE FAILURE - Agent terminated
```

## Correct Patterns

### GOOD: Proactive Monitoring
```bash
# At 650 lines
echo "📊 Current size: 650/800 lines"
echo "📋 Planning for potential split at 700"
echo "🎯 Identifying logical boundaries now"
```

### GOOD: Immediate Response
```bash
# At 801 lines
echo "🚨 SIZE VIOLATION: 801 lines!"
echo "🛑 STOPPING all implementation"
echo "📋 Creating split plan..."
echo "🚀 Spawning Code Reviewer for split"
```

### GOOD: Clean Splits
```bash
# Split Plan Execution
Split 1: Core logic (650 lines)
Split 2: API endpoints (500 lines)  
Split 3: Tests (400 lines)
# Total: Same functionality, compliant sizes
```

## Size Calculation Details

### ✅ What Counts (IMPLEMENTATION ONLY)
- Core business logic source code
- API endpoint implementations  
- Service layer code
- Critical algorithms and data structures
- Production code comments
- Modified generated code (if you edited it)

### ❌ What NEVER Counts (EXCLUDED)
- **Demo files**: demos/*, demo-*, DEMO.md, example-*
- **Test files**: *_test.go, test/*, tests/*, *.test.*, fixtures/*
- **Documentation**: *.md, docs/*, README*, LICENSE*
- **Generated code**: *.pb.go, *_generated.*, *.gen.go
- **Configuration**: *.json, *.yaml, *.yml, *.toml, *.ini
- **Dependencies**: vendor/*, node_modules/*, .cache/*
- **Build artifacts**: bin/*, dist/*, build/*, *.o, *.so
- **Lock files**: *.lock, go.sum, package-lock.json
- **CI/CD**: .github/*, Jenkinsfile, .gitlab-ci.yml
- **Temporary**: *.tmp, *.bak, *.swp

## Grading Formula

```python
def calculate_size_compliance_grade(efforts):
    violations = 0
    warnings_ignored = 0
    
    for effort in efforts:
        if effort.final_size > 800:
            violations += 1
            if effort.had_warning_at_700:
                warnings_ignored += 1
    
    # Base grade
    grade = 100
    
    # Deductions
    grade -= violations * 20
    grade -= warnings_ignored * 5
    
    # Catastrophic failure
    if any(e.merged_oversized for e in efforts):
        grade = 0
    
    return max(grade, 0)
```

## Prevention Strategies

### 1. Architecture First
```bash
# Before implementation
echo "📐 Designing with size limits in mind"
echo "📦 Breaking into <600 line components"
echo "🎯 Each effort naturally under limit"
```

### 2. Continuous Monitoring
```bash
# During implementation
watch -n 300 "line-counter.sh -c current-branch"
# Check every 5 minutes
```

### 3. Early Splitting
```bash
# At 650 lines
echo "🤔 Currently 650 lines"
echo "📊 Estimating 200 more lines needed"
echo "✂️ Initiating proactive split now"
```

## Recovery from Violation

If you've exceeded 800 lines:

1. **DO NOT PANIC**
2. **DO NOT MERGE**
3. Save work: `git add . && git commit -m "WIP: Pre-split checkpoint"`
4. Measure exactly: `line-counter.sh -c branch -d`
5. Create split plan immediately
6. Execute splits sequentially
7. Each split gets reviewed
8. Only merge compliant splits

## Split Effort Compliance (R297)

### Already-Split Efforts Are Compliant
When an effort has been split:
- Each split branch must be ≤800 lines
- The integration branch WILL exceed 800 lines (EXPECTED)
- This is COMPLIANT because PRs come from split branches
- Architects must check `split_count` before measuring

### Example: Split Effort Compliance
```yaml
# E1.1.2 was split into 2 parts:
split_1_branch: 450 lines  # ✅ Compliant
split_2_branch: 454 lines  # ✅ Compliant
integration_branch: 904 lines  # ✅ Still compliant (integration expected to exceed)

# orchestrator-state.json shows:
efforts_completed:
  E1.1.2:
    split_count: 2  # This means it's compliant!
```

## The 800 Line Law

```
No effort exceeds 800 lines - EVER
Check split_count first - ALWAYS (R297)
Measure original branches - NOT integration
Measure with the right tool - ALWAYS  
Split when needed - IMMEDIATELY
Monitor continuously - PROACTIVELY
Merge only compliant code - STRICTLY
```

---
**Remember:** 800 lines is not a suggestion, it's an absolute limit. But already-split efforts are compliant even if their integration exceeds the limit!