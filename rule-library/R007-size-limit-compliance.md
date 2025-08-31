# 🚨 RULE R007 - Size Limit Compliance (800 Line Maximum)

**Criticality:** BLOCKING - Cannot proceed if violated  
**Grading Impact:** -20% per violation, -100% if merged  
**Enforcement:** CONTINUOUS - Check every 200 lines added

## Rule Statement

NO effort may EVER exceed 800 lines. Soft warning at 700 lines. Automatic split required at violation.

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

### What Counts
- All source code lines
- Comments in source files
- Test files
- Generated code that you modified
- Configuration as code

### What Doesn't Count
- Pure generated code (untouched)
- Markdown documentation
- YAML/JSON configs (unless executable)
- Binary files
- Vendor/node_modules

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

## The 800 Line Law

```
No effort exceeds 800 lines - EVER
Measure with the right tool - ALWAYS  
Split when needed - IMMEDIATELY
Monitor continuously - PROACTIVELY
Merge only compliant code - STRICTLY
```

---
**Remember:** 800 lines is not a suggestion, it's an absolute limit. Violate it and fail.