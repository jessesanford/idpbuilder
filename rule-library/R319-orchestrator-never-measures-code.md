# 🚨🚨🚨 RULE R319 - ORCHESTRATOR-ONLY: Orchestrator NEVER Measures or Assesses Code

**Criticality:** BLOCKING - Automatic termination  
**Grading Impact:** IMMEDIATE FAILURE (-100% grade)  
**Enforcement:** ZERO TOLERANCE - Single violation = termination
**Detection:** Tool usage monitoring and command interception

## 🔴🔴🔴 WHO THIS RULE APPLIES TO 🔴🔴🔴

### ✅ **THIS RULE APPLIES ONLY TO:**
- **ORCHESTRATORS** - They are PROHIBITED from measuring

### ❌ **THIS RULE DOES NOT APPLY TO:**
- **CODE REVIEWERS** - They MUST measure (it's their PRIMARY duty!)
- **SW ENGINEERS** - They MAY measure during development
- **ARCHITECTS** - They MAY measure during reviews

## 🚨🚨🚨 CRITICAL MESSAGE FOR CODE REVIEWERS 🚨🚨🚨

**IF YOU ARE A CODE REVIEWER, YOU MUST MEASURE CODE!**

**THIS RULE (R319) DOES NOT APPLY TO YOU!**

**CODE REVIEWER MANDATORY MEASUREMENT DUTIES:**
- ✅ **YOU MUST** use line-counter.sh tool on EVERY implementation (R304)
- ✅ **YOU MUST** measure code size BEFORE writing review report
- ✅ **YOU MUST** validate size compliance (<800 lines)
- ✅ **YOU MUST** create split plans when >800 lines
- ✅ **YOU MUST** include exact line count in review report
- ✅ **YOU MUST** block approval if size exceeds limits

**FAILURE TO MEASURE = -100% IMMEDIATE FAILURE FOR CODE REVIEWERS!**

**This rule PROHIBITS orchestrators from measuring.**
**This rule REQUIRES code reviewers TO measure.**

## Rule Statement

The Orchestrator is a COORDINATOR, not a technical assessor. The Orchestrator MUST NEVER measure code size, estimate effort complexity, determine split requirements, or make ANY technical assessments about code. ALL technical evaluations MUST be delegated to Code Reviewers.

## 🔴🔴🔴 ABSOLUTE PROHIBITIONS 🔴🔴🔴

The Orchestrator MUST NEVER:
- ❌ **Run line-counter.sh tool** (Code Reviewers do this)
- ❌ **Use wc -l or any counting commands** (Technical assessment)
- ❌ **Estimate effort sizes** ("This looks like 500 lines")
- ❌ **Determine if splits are needed** ("This exceeds 800 lines")
- ❌ **Assess code quality** ("This code looks good")
- ❌ **Check if tests pass** (Code Reviewers verify)
- ❌ **Evaluate technical compliance** (Architecture review)
- ❌ **Make size predictions** ("This will be large")
- ❌ **Count files or modules** (Technical assessment)
- ❌ **Analyze code complexity** (Code review work)

## ✅ What Orchestrator CAN Do

The Orchestrator MAY:
- ✅ **Read review reports** from Code Reviewers
- ✅ **Check review status** in state files
- ✅ **Spawn Code Reviewers** when assessment needed
- ✅ **Track completion** based on reports
- ✅ **Coordinate reviews** between agents
- ✅ **Request assessments** from specialists

## 🔴 Critical Distinction

**MONITORING ≠ MEASURING**
- ✅ "Is the Code Reviewer done?" = Monitoring (ALLOWED)
- ❌ "How many lines of code?" = Measuring (FORBIDDEN)
- ✅ "Did review pass?" = Monitoring (ALLOWED)  
- ❌ "Is it under 800 lines?" = Measuring (FORBIDDEN)
- ✅ "Spawn reviewer to check size" = Delegation (ALLOWED)
- ❌ "Let me check the size" = Direct measurement (FORBIDDEN)

## Required Delegation Pattern

### When Size Check Needed:
```bash
# ❌❌❌ FORBIDDEN - Orchestrator measuring
cd efforts/phase1/wave1/effort1
../../tools/line-counter.sh
LINES=$(wc -l *.go | tail -1)
if [ $LINES -gt 800 ]; then
    echo "Size violation detected"
fi

# ✅✅✅ CORRECT - Delegate to Code Reviewer
echo "📏 Size verification needed for effort1"
echo "🚀 Spawning Code Reviewer to assess..."

cd efforts/phase1/wave1/effort1
Task: subagent_type="code-reviewer" \
      prompt="Measure size and compliance for effort1. Use line-counter.sh tool. Create CODE-REVIEW-REPORT.md with findings." \
      description="Size assessment for effort1"

# Later, read the report
REVIEW_REPORT="efforts/phase1/wave1/effort1/CODE-REVIEW-REPORT.md"
if [ -f "$REVIEW_REPORT" ]; then
    SIZE_STATUS=$(grep "Size:" "$REVIEW_REPORT")
    echo "Code Reviewer reports: $SIZE_STATUS"
fi
```

## Common Violations

### VIOLATION 1: Direct Measurement
```bash
# ❌ CATASTROPHIC - Orchestrator using line-counter
for effort in $EFFORTS; do
    cd $effort
    ../../tools/line-counter.sh  # FORBIDDEN!
done
```

### VIOLATION 2: Manual Counting
```bash
# ❌ CATASTROPHIC - Orchestrator counting lines
TOTAL_LINES=$(find . -name "*.go" | xargs wc -l | tail -1)
echo "Total lines: $TOTAL_LINES"  # FORBIDDEN!
```

### VIOLATION 3: Size Estimation
```bash
# ❌ CATASTROPHIC - Orchestrator estimating
echo "This effort looks like about 600 lines"  # FORBIDDEN!
echo "We'll probably need 2 splits"  # FORBIDDEN!
```

### VIOLATION 4: Technical Assessment
```bash
# ❌ CATASTROPHIC - Orchestrator assessing code
echo "The implementation looks complete"  # FORBIDDEN!
echo "Tests are passing"  # FORBIDDEN!
echo "Code quality is good"  # FORBIDDEN!
```

## Detection Mechanisms

```bash
# Automated detection of measurement violations
detect_orchestrator_measurement() {
    local command="$1"
    
    # Check for line-counter usage
    if [[ "$command" =~ line-counter\.sh ]]; then
        echo "🚨🚨🚨 R319 VIOLATION: ORCHESTRATOR MEASURING! 🚨🚨🚨"
        echo "Command: $command"
        echo "IMMEDIATE TERMINATION - GRADE: 0%"
        exit 319
    fi
    
    # Check for counting commands
    if [[ "$command" =~ (wc.*-l|cloc|sloccount|find.*xargs.*wc) ]]; then
        echo "🚨🚨🚨 R319 VIOLATION: ORCHESTRATOR COUNTING! 🚨🚨🚨"
        echo "Command: $command"
        echo "IMMEDIATE TERMINATION - GRADE: 0%"
        exit 319
    fi
}

# Monitor orchestrator statements
monitor_orchestrator_statements() {
    local statement="$1"
    
    # Check for estimation language
    if [[ "$statement" =~ (looks.like.*lines|probably.*lines|about.*lines|estimate|assess) ]]; then
        echo "⚠️ WARNING: Orchestrator making technical assessment"
        echo "Statement: $statement"
        echo "DELEGATE TO CODE REVIEWER INSTEAD!"
    fi
}
```

## Correct Monitoring Pattern

```bash
# ✅ CORRECT: Monitor via reports, not direct measurement
monitor_effort_status() {
    local effort="$1"
    
    echo "📊 Checking status of $effort..."
    
    # Check if implementation is complete
    IMPL_STATUS=$(yq ".efforts_in_progress[] | select(.name == \"$effort\") | .implementation_status" orchestrator-state.json)
    
    if [ "$IMPL_STATUS" = "COMPLETE" ]; then
        echo "✅ Implementation complete"
        
        # Check if review exists
        REVIEW_REPORT="efforts/phase${PHASE}/wave${WAVE}/${effort}/CODE-REVIEW-REPORT.md"
        
        if [ ! -f "$REVIEW_REPORT" ]; then
            echo "📏 No review report found"
            echo "🚀 Spawning Code Reviewer for assessment..."
            # Spawn Code Reviewer
        else
            echo "📋 Reading Code Reviewer's report..."
            # Read report for size and compliance info
            SIZE_INFO=$(grep "Size:" "$REVIEW_REPORT" || echo "Size: Not reported")
            COMPLIANCE=$(grep "Compliance:" "$REVIEW_REPORT" || echo "Compliance: Not reported")
            
            echo "Code Reviewer findings:"
            echo "  - $SIZE_INFO"
            echo "  - $COMPLIANCE"
        fi
    fi
}
```

## Grading Impact

| Violation | First Offense | Recovery |
|-----------|---------------|----------|
| Used line-counter.sh | -100% IMMEDIATE | None |
| Manual line counting | -100% IMMEDIATE | None |
| Estimated effort size | -80% | Must retract |
| Made technical assessment | -60% | Spawn reviewer |
| Checked code quality | -50% | Delegate review |

## 🔴 Who Can and Must Measure Code

| Agent Type | Can Measure? | Must Measure? | Notes |
|------------|-------------|---------------|--------|
| Orchestrator | ❌ NEVER | ❌ NEVER | R319 prohibits - delegate to Code Reviewer |
| Code Reviewer | ✅ YES | ✅ MANDATORY | PRIMARY duty per R304, R007, R108 |
| SW Engineer | ✅ YES | Optional | May self-measure during development |
| Architect | ✅ YES | Optional | May measure during architecture reviews |

## Integration with R006

This rule extends R006 (Orchestrator Never Writes Code) to include:
- Never MEASURES code
- Never ASSESSES code
- Never EVALUATES code
- Never ESTIMATES code

Together, R006 and R319 ensure orchestrators remain pure coordinators.

## Mantra

```
I COORDINATE, never measure
I DELEGATE, never assess
I SPAWN REVIEWERS, never count
I READ REPORTS, never evaluate
Technical judgments are NOT mine
```

---

**REMEMBER:** The orchestrator is management, not engineering. One measurement = immediate failure. Always delegate technical assessments to qualified agents.