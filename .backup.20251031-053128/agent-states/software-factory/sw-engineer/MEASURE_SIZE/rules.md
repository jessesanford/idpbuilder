# SW Engineer - MEASURE_SIZE State Rules

## State Context
You have reached a size checkpoint or warning threshold and need to measure exact size compliance and plan next steps.

## 🔴🔴🔴 CRITICAL: R340 Plan Location Tracking 🔴🔴🔴

**RULE R340: Planning File Metadata Tracking (BLOCKING)**
- Split plans (if any) are tracked in orchestrator-state-v3.json
- NEVER search for plans using `find` or `ls` commands
- The orchestrator tracks ALL planning files in the state file
- Violation of R340 = Integration delays and failures

## 🔴🔴🔴 CRITICAL: Finding the Line Counter Tool 🔴🔴🔴

**SPLIT EFFORTS ARE DIFFERENT** - They have their own git repositories!

```bash
# UNIVERSAL LINE COUNTER FINDER - Works for both splits and normal efforts
find_line_counter() {
    echo "🔍 Finding line counter tool..."
    
    # Check if this is a split effort (by directory name, not by searching for plans)
    if [[ "$(pwd)" == *"-split-"* ]] || [[ "$(pwd)" == *"-SPLIT-"* ]]; then
        echo "📂 Split effort detected (from directory name) - searching UP directory tree"
        IS_SPLIT=true
    else
        echo "📂 Normal effort - searching for project root"
        IS_SPLIT=false
    fi
    
    # Strategy 1: Search UP the directory tree for tools/line-counter.sh
    SEARCH_DIR=$(pwd)
    MAX_LEVELS=10
    LEVEL=0
    
    while [ "$SEARCH_DIR" != "/" ] && [ $LEVEL -lt $MAX_LEVELS ]; do
        if [ -f "$SEARCH_DIR/tools/line-counter.sh" ]; then
            LINE_COUNTER="$SEARCH_DIR/tools/line-counter.sh"
            echo "✅ Found line counter at: $LINE_COUNTER"
            return 0
        fi
        SEARCH_DIR=$(dirname "$SEARCH_DIR")
        LEVEL=$((LEVEL + 1))
    done
    
    # Strategy 2: Look for orchestrator-state-v3.json (main project marker)
    SEARCH_DIR=$(pwd)
    while [ "$SEARCH_DIR" != "/" ]; do
        if [ -f "$SEARCH_DIR/orchestrator-state-v3.json" ]; then
            if [ -f "$SEARCH_DIR/tools/line-counter.sh" ]; then
                LINE_COUNTER="$SEARCH_DIR/tools/line-counter.sh"
                echo "✅ Found via orchestrator-state-v3.json: $LINE_COUNTER"
                return 0
            fi
        fi
        SEARCH_DIR=$(dirname "$SEARCH_DIR")
    done
    
    # Strategy 3: Try relative paths from split directories
    for relative_path in \
        "../../../../../tools/line-counter.sh" \
        "../../../../tools/line-counter.sh" \
        "../../../tools/line-counter.sh" \
        "../../tools/line-counter.sh" \
        "../tools/line-counter.sh"; do
        if [ -f "$relative_path" ]; then
            LINE_COUNTER=$(realpath "$relative_path")
            echo "✅ Found via relative path: $LINE_COUNTER"
            return 0
        fi
    done
    
    echo "❌ ERROR: Could not find line-counter.sh"
    echo "   Searched from: $(pwd)"
    echo "   Try running: find /home -name 'line-counter.sh' 2>/dev/null"
    return 1
}

# Find the tool
find_line_counter

# Run measurement if found
if [ -n "$LINE_COUNTER" ]; then
    BRANCH=$(git branch --show-current)
    echo "📊 Measuring branch: $BRANCH"
    
    # For splits, always use explicit branch parameter
    if [ "$IS_SPLIT" = true ]; then
        "$LINE_COUNTER" -c "$BRANCH"
    else
        "$LINE_COUNTER"  # Normal efforts can auto-detect
    fi
else
    echo "❌ Cannot proceed without line counter"
    exit 1
fi
```

---
### ⚠️ RULE R107.0.0 - MEASURE_SIZE Rules
**Source:** rule-library/RULE-REGISTRY.md#R107
**Criticality:** IMPORTANT - Affects workflow

SIZE MEASUREMENT PROTOCOL:
1. Find line counter using the universal finder above
2. Generate detailed breakdown analysis
3. Determine exact compliance status
4. Calculate remaining capacity or overflow
5. Make transition decision based on measurement
6. Document findings and recommendations
---

## Mandatory Size Measurement

---
### 🚨🚨 RULE R007.0.0 - Size Limit Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R007
**Criticality:** MANDATORY - Required for approval

MEASUREMENT REQUIREMENTS:

⚠️ CRITICAL: Use ONLY the project's line-counter.sh tool
⚠️ NEVER count lines manually or with other tools
⚠️ Exclude generated code automatically

REQUIRED MEASUREMENTS:
1. Total effort size (primary metric)
2. Size breakdown by package/directory
3. Growth trend analysis
4. Remaining capacity calculation

DECISION MATRIX:
- ≤700 lines: CONTINUE implementation
- 701-750 lines: OPTIMIZE or prepare for completion
- 751-800 lines: COMPLETE immediately or split
- >800 lines: MANDATORY split required
---

## Comprehensive Size Measurement Script

```bash
#!/bin/bash
# Comprehensive size measurement for any effort type

measure_size() {
    local BRANCH=$(git branch --show-current)
    local WORKING_DIR=$(pwd)
    local TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
    
    echo "🔍 COMPREHENSIVE SIZE MEASUREMENT - $TIMESTAMP"
    echo "Branch: $BRANCH"
    echo "Working Directory: $WORKING_DIR"
    echo ""
    
    # Find line counter first
    LINE_COUNTER=""
    IS_SPLIT=false
    
    # Check if this is a split (by directory name, not by searching for plans - R340)
    if [[ "$(pwd)" == *"-split-"* ]] || [[ "$(pwd)" == *"-SPLIT-"* ]]; then
        IS_SPLIT=true
        echo "📂 Split effort detected (from directory name per R340)"
    fi
    
    # Search for line counter
    echo "🔍 Locating line counter tool..."
    SEARCH_DIR=$(pwd)
    MAX_LEVELS=10
    LEVEL=0
    
    while [ "$SEARCH_DIR" != "/" ] && [ $LEVEL -lt $MAX_LEVELS ]; do
        if [ -f "$SEARCH_DIR/tools/line-counter.sh" ]; then
            LINE_COUNTER="$SEARCH_DIR/tools/line-counter.sh"
            break
        fi
        SEARCH_DIR=$(dirname "$SEARCH_DIR")
        LEVEL=$((LEVEL + 1))
    done
    
    if [ -z "$LINE_COUNTER" ]; then
        echo "❌ FATAL: Cannot find line-counter.sh"
        echo "   Searched up to $LEVEL levels from $(pwd)"
        exit 1
    fi
    
    echo "✅ Found: $LINE_COUNTER"
    echo ""
    
    # Primary measurement
    echo "📏 PRIMARY SIZE MEASUREMENT"
    echo "Tool: $LINE_COUNTER"
    
    if [ "$IS_SPLIT" = true ]; then
        echo "Command: $LINE_COUNTER -c \"$BRANCH\""
        SIZE_RESULT=$("$LINE_COUNTER" -c "$BRANCH")
    else
        echo "Command: $LINE_COUNTER"
        SIZE_RESULT=$("$LINE_COUNTER")
    fi
    
    echo "$SIZE_RESULT"
    
    # Extract total lines from result
    TOTAL_LINES=$(echo "$SIZE_RESULT" | grep "Total non-generated lines:" | grep -oE '[0-9]+' | head -1)
    
    if [ -z "$TOTAL_LINES" ]; then
        echo "⚠️ Could not extract line count, trying alternative parsing..."
        TOTAL_LINES=$(echo "$SIZE_RESULT" | grep -oE '[0-9]+ lines' | head -1 | grep -oE '[0-9]+')
    fi
    
    if [ -z "$TOTAL_LINES" ]; then
        echo "❌ ERROR: Could not determine line count from output"
        exit 1
    fi
    
    echo ""
    
    # Detailed breakdown
    echo "📊 DETAILED SIZE BREAKDOWN"
    if [ "$IS_SPLIT" = true ]; then
        DETAILED_RESULT=$("$LINE_COUNTER" -c "$BRANCH" -d)
    else
        DETAILED_RESULT=$("$LINE_COUNTER" -d)
    fi
    echo "$DETAILED_RESULT"
    echo ""
    
    # Calculate compliance metrics
    echo "⚖️ COMPLIANCE ANALYSIS"
    LIMIT=800
    UTILIZATION=$(awk "BEGIN {printf \"%.2f\", ($TOTAL_LINES * 100) / $LIMIT}")
    REMAINING=$((LIMIT - TOTAL_LINES))
    
    echo "Total Lines: $TOTAL_LINES"
    echo "Limit: $LIMIT"
    echo "Utilization: ${UTILIZATION}%"
    echo "Remaining Capacity: $REMAINING lines"
    echo ""
    
    # Determine status
    if [ "$TOTAL_LINES" -gt 800 ]; then
        STATUS="VIOLATION"
        SEVERITY="CRITICAL"
        echo "🚨 STATUS: SIZE LIMIT VIOLATION"
        echo "🚨 SEVERITY: CRITICAL - Immediate split required"
        EXIT_CODE=3
    elif [ "$TOTAL_LINES" -gt 750 ]; then
        STATUS="DANGER"
        SEVERITY="HIGH"
        echo "⚠️ STATUS: DANGER ZONE"
        echo "⚠️ SEVERITY: HIGH - Complete immediately or split"
        EXIT_CODE=2
    elif [ "$TOTAL_LINES" -gt 700 ]; then
        STATUS="WARNING"
        SEVERITY="MEDIUM"
        echo "⚠️ STATUS: WARNING"
        echo "⚠️ SEVERITY: MEDIUM - Optimize and plan completion"
        EXIT_CODE=1
    else
        STATUS="COMPLIANT"
        SEVERITY="LOW"
        echo "✅ STATUS: COMPLIANT"
        echo "✅ SEVERITY: LOW - Safe to continue"
        EXIT_CODE=0
    fi
    echo ""
    
    # Log measurement to work log
    echo "- [$TIMESTAMP] Size measurement: $TOTAL_LINES/$LIMIT lines (${UTILIZATION}%) - $STATUS" >> work-log.md
    
    return $EXIT_CODE
}

# Execute measurement
measure_size
RESULT=$?

# Make decision based on result
case $RESULT in
    0) echo "→ Decision: Continue implementation safely";;
    1) echo "→ Decision: Continue with caution, monitor closely";;
    2) echo "→ Decision: Complete immediately or plan split";;
    3) echo "→ Decision: Stop work, create split plan NOW";;
esac

exit $RESULT
```

## Quick Measurement Commands

For quick measurement without the full analysis:

```bash
# Quick measure for normal efforts
LC=$(d=$(pwd); while [ "$d" != "/" ]; do [ -f "$d/tools/line-counter.sh" ] && echo "$d/tools/line-counter.sh" && break; d=$(dirname "$d"); done)
[ -n "$LC" ] && "$LC" || echo "Tool not found"

# Quick measure for split efforts (with branch)
LC=$(d=$(pwd); while [ "$d" != "/" ]; do [ -f "$d/tools/line-counter.sh" ] && echo "$d/tools/line-counter.sh" && break; d=$(dirname "$d"); done)
[ -n "$LC" ] && "$LC" -c "$(git branch --show-current)" || echo "Tool not found"
```

## Decision Making Framework

---
### 🚨🚨 RULE R020.0.0 - State Transitions
**Source:** rule-library/RULE-REGISTRY.md#R020
**Criticality:** MANDATORY - Required for approval

SIZE-BASED DECISION MATRIX:

>800 LINES (VIOLATION):
→ SPLIT_WORK (mandatory)
→ Document violation and split strategy

751-800 LINES (DANGER):
→ Evaluate completion feasibility
→ IMPLEMENTATION if <10% work remaining
→ SPLIT_WORK if >10% work remaining

701-750 LINES (WARNING):
→ Optimize current code if opportunities exist
→ IMPLEMENTATION with careful monitoring
→ Consider completion strategy

≤700 LINES (COMPLIANT):
→ IMPLEMENTATION (safe to continue)
→ Continue normal development pace
---

## State Transitions

From MEASURE_SIZE state:
- **SIZE_COMPLIANT** → IMPLEMENTATION (Continue development safely)
- **SIZE_WARNING** → IMPLEMENTATION (Continue with enhanced monitoring)
- **SIZE_DANGER** → IMPLEMENTATION or SPLIT_WORK (Based on completion analysis)
- **SIZE_VIOLATION** → SPLIT_WORK (Mandatory split required)
- **OPTIMIZATION_IDENTIFIED** → FIX_ISSUES (Refactor before continuing)

## Common Issues and Solutions

### Issue: "Line counter not found"
**Solution**: Use the universal finder script above

### Issue: "Not in a git repository"
**For Splits**: You ARE in a git repo (sparse checkout). Use `-c branch-name` parameter

### Issue: "Branch doesn't exist"
**Check**: Run `git branch --show-current` to get exact branch name

### Issue: "Fatal: bad revision"
**Fix**: The branch name might have typos. Copy from `git branch --show-current`

## Remember

1. **Splits are different** - They're sparse checkouts with separate git repos
2. **Search UP for tools** - The main project is above split directories  
3. **Use branch parameter for splits** - Auto-detection may fail
4. **Never hardcode paths** - Projects can be in different locations
5. **Log all measurements** - Update work-log.md with results

## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

