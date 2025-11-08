# 🚨🚨🚨 BLOCKING: RULE R314 - Mandatory File Count Validation Protocol 🚨🚨🚨

**Category:** Split Implementation Control  
**Agents:** SW Engineer (PRIMARY), Code Reviewer (ENFORCER), Orchestrator  
**Criticality:** BLOCKING - Prevents catastrophic scope violations  
**Penalty:** -100% for >200% file count, automatic rejection

## 🔴🔴🔴 SUPREME DIRECTIVE: FILE COUNT IS THE PRIMARY CONTROL 🔴🔴🔴

**After a 2667% violation (80 files instead of 3), file count validation is MANDATORY**

## The Catastrophic Problem This Prevents

### 🚨 THE 2667% VIOLATION INCIDENT 🚨
```
Date: [REDACTED]
Split Plan: Implement 3 files in pkg/builder
SW Engineer Action: Implemented 80 files across 26 packages
Violation: 2667% of planned scope
Root Cause: Complete disregard for split plan
Impact: Total project failure, complete re-implementation required
```

This single incident proves that without strict file count control, SW Engineers can and will implement entire codebases instead of focused splits.

## 🔴 MANDATORY FILE COUNT VALIDATION PROTOCOL

### 1. SPLIT PLANS MUST SPECIFY EXACT FILE COUNTS

**Every split plan MUST include:**
```markdown
## 📊 SCOPE METRICS (MANDATORY)
- **EXACT FILES TO IMPLEMENT: 3**
- **MAXIMUM ALLOWED FILES: 6** (200% hard limit)
- **Functions to implement: 5**
- **Methods to implement: 2**
- **Tests to write: 3**

### Files to Implement (EXACTLY THESE):
1. pkg/builder/options.go
2. pkg/builder/options_test.go
3. pkg/builder/types.go
```

**Plans without explicit file counts are INVALID and must be rejected.**

### 2. PRE-IMPLEMENTATION FILE COUNT VALIDATION

**Before writing ANY code, SW Engineers MUST:**

```bash
# MANDATORY: Extract and validate file count
validate_file_count_before_start() {
    PLANNED_FILES=$(grep -c "^###.*File:\|^[0-9]\+\.\s.*\.go" SPLIT-PLAN-*.md)
    
    if [ -z "$PLANNED_FILES" ] || [ "$PLANNED_FILES" -eq 0 ]; then
        echo "🚨 FATAL: No file count in split plan!"
        echo "Cannot proceed without explicit file list"
        exit 314
    fi
    
    echo "🔴🔴🔴 FILE COUNT CONTRACT 🔴🔴🔴"
    echo "You are committing to implement EXACTLY $PLANNED_FILES files"
    echo "Exceeding by >2x = AUTOMATIC FAILURE"
    
    # Create enforcement file
    echo "$PLANNED_FILES" > .planned-file-count
    touch .implementation-start-marker
    
    # Require explicit acknowledgment
    read -p "Type 'I will implement exactly $PLANNED_FILES files' to proceed: " ack
    if [[ "$ack" != "I will implement exactly $PLANNED_FILES files" ]]; then
        echo "❌ Acknowledgment required to proceed"
        exit 1
    fi
}
```

### 3. INCREMENTAL FILE COUNT MONITORING_SWE_PROGRESS

**Every 100 lines of code added:**

```bash
# MANDATORY: Check file count incrementally
monitor_file_count() {
    PLANNED=$(cat .planned-file-count)
    ACTUAL=$(find . -name "*.go" -newer .implementation-start-marker | wc -l)
    RATIO=$((ACTUAL * 100 / PLANNED))
    
    echo "📊 FILE COUNT MONITOR:"
    echo "  Planned: $PLANNED files"
    echo "  Current: $ACTUAL files"
    echo "  Ratio: $RATIO%"
    
    # Escalating warnings
    if [ $RATIO -gt 150 ]; then
        echo "🚨🚨🚨 CRITICAL: Approaching 2x file limit!"
        echo "Current violation level: $RATIO%"
        echo "STOP IMMEDIATELY or face automatic failure!"
    elif [ $RATIO -gt 120 ]; then
        echo "⚠️ WARNING: File count exceeding plan by 20%"
    fi
    
    # Hard stop at 2x
    if [ $RATIO -gt 200 ]; then
        echo "❌❌❌ AUTOMATIC FAILURE: File count violation!"
        echo "You have $ACTUAL files but planned only $PLANNED"
        echo "This is a $RATIO% violation!"
        exit 314
    fi
}

# Hook this into git pre-commit
echo "monitor_file_count" >> .git/hooks/pre-commit
```

### 4. COMMIT-TIME ENFORCEMENT

**Before EVERY commit:**

```bash
# MANDATORY: Final file count validation
pre_commit_file_validation() {
    PLANNED=$(cat .planned-file-count 2>/dev/null)
    if [ -z "$PLANNED" ]; then
        echo "❌ No planned file count found!"
        echo "Run validate_file_count_before_start first"
        exit 314
    fi
    
    # Count actual files in commit
    STAGED=$(git diff --cached --name-only | grep "\.go$" | wc -l)
    TOTAL=$(find . -name "*.go" -newer .implementation-start-marker | wc -l)
    
    echo "📊 COMMIT FILE VALIDATION:"
    echo "  Files in this commit: $STAGED"
    echo "  Total files created: $TOTAL"
    echo "  Planned files: $PLANNED"
    
    if [ $TOTAL -gt $((PLANNED * 2)) ]; then
        echo "🚨 COMMIT BLOCKED: File count violation!"
        echo "Remove excess files before committing"
        git status --short
        exit 314
    fi
}
```

### 5. CODE REVIEWER ENFORCEMENT

**Code Reviewers MUST check file count FIRST:**

```bash
# FIRST thing in any split review
review_file_count() {
    echo "🔴 R314: FILE COUNT VALIDATION (CHECK THIS FIRST)"
    
    PLANNED=$(grep -c "^###.*File:\|^[0-9]\+\.\s.*\.go" SPLIT-PLAN-*.md)
    ACTUAL=$(git diff --name-only $BASE_BRANCH | grep "\.go$" | wc -l)
    RATIO=$((ACTUAL * 100 / PLANNED))
    
    echo "Planned files: $PLANNED"
    echo "Actual files: $ACTUAL"
    echo "Violation ratio: $RATIO%"
    
    if [ $RATIO -gt 200 ]; then
        echo "❌ AUTOMATIC REJECTION: >200% file count"
        echo "Grade: F"
        echo "No further review needed"
        exit 314
    elif [ $RATIO -gt 150 ]; then
        echo "⚠️ SEVERE WARNING: >150% file count"
        echo "Maximum grade: D"
    fi
}
```

## 📈 File Count Grading Scale

| Violation % | Grade | Action |
|------------|-------|--------|
| 100% exact | A+ | Perfect adherence |
| 101-110% | A | Excellent |
| 111-120% | B | Good |
| 121-150% | C | Concerning |
| 151-199% | D | Poor |
| 200%+ | F | AUTOMATIC FAILURE |
| 2667% | ☠️ | Career-ending |

## 🛡️ Prevention Mechanisms

### Split Plan Template Requirements
All split plans MUST include:
1. Exact file count at the top
2. List of specific files to implement
3. "DO NOT IMPLEMENT" section for excluded files
4. Maximum file count warning

### Automatic Safeguards
1. Pre-implementation checklist (mandatory)
2. File count validation before first line of code
3. Incremental monitoring every 100 lines
4. Commit hooks for validation
5. Code reviewer file count check (first action)

## 🔴 Enforcement Across Agents

### SW Engineer Responsibilities
- Complete pre-implementation checklist
- Validate file count before starting
- Monitor incrementally
- Stop if approaching 2x limit

### Code Reviewer Responsibilities
- Check file count BEFORE anything else
- Automatic F grade for >200% violations
- Document violation percentage in review

### Orchestrator Responsibilities
- Ensure split plans have explicit file counts
- Track file count violations in state
- Reject efforts with >200% violations

## Examples

### ✅ Good: Exact Adherence
```
Plan: 3 files
Implemented: 3 files
Ratio: 100%
Grade: A+
```

### ⚠️ Warning: Minor Overage
```
Plan: 3 files
Implemented: 4 files
Ratio: 133%
Grade: C
```

### ❌ Failure: Major Violation
```
Plan: 3 files
Implemented: 7 files
Ratio: 233%
Grade: F (AUTOMATIC)
```

### ☠️ Catastrophic: The 2667% Incident
```
Plan: 3 files
Implemented: 80 files
Ratio: 2667%
Grade: Project termination
```

## Summary

**File count is the #1 indicator of split scope violations.**

After the 2667% incident, file count validation is no longer optional - it's the primary control mechanism preventing catastrophic failures.

**Key Points:**
- Split plans MUST specify exact file counts
- SW Engineers MUST validate before implementation
- Incremental monitoring every 100 lines is MANDATORY
- >200% file count = AUTOMATIC FAILURE
- Code reviewers check file count FIRST

**Remember: Someone implemented 80 files instead of 3. This rule ensures that NEVER happens again.**