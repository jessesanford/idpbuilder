# Rule R305: SW Engineer Self-Monitoring Protocol

## Rule Statement
Software Engineers MUST continuously self-monitor code size during implementation using the line-counter.sh tool with correct parameters to prevent size violations proactively.

## Criticality Level
**BLOCKING** - Failure to self-monitor results in grading penalties and potential automatic failure

## Enforcement Mechanism
- **Technical**: Automated monitoring scripts and measurement logs
- **Behavioral**: Mandatory measurement at defined intervals
- **Grading**: -20% per missed measurement, -100% for exceeding limits

## 🔴🔴🔴 SUPREME REQUIREMENT: PROACTIVE SIZE MANAGEMENT 🔴🔴🔴

### Core Principle
SW Engineers are RESPONSIBLE for preventing size violations through continuous self-monitoring, not discovering them after the fact.

## Monitoring Requirements

### 1. Baseline Measurement (STARTUP)
```bash
# MANDATORY: First action when entering IMPLEMENTATION state
echo "📊 ESTABLISHING BASELINE MEASUREMENT"

# 🔴🔴🔴 CRITICAL: EFFORTS ARE SEPARATE GIT REPOSITORIES! 🔴🔴🔴
# You are working in your OWN git repository, not a branch of the main repo!

# Step 1: Verify you're in the effort repository
pwd  # Should show: /efforts/phase1/wave1/your-effort-name
ls -la .git  # MUST exist - this is YOUR git repository!

# Step 2: Ensure code is committed (git diff needs commits!)
git status  # Should be clean or commit first:
# git add -A && git commit -m "checkpoint: measuring progress"

# Step 3: Get ACTUAL branch name (NOT directory name!)
CURRENT_BRANCH=$(git branch --show-current)
echo "Current branch: $CURRENT_BRANCH"  # e.g., phase1-wave1-effort-name

# Step 4: Find base branch IN THIS REPOSITORY
git branch -a | grep -E "integration|main"  # See what exists
BASE="phase1/integration"  # Or whatever exists in THIS repo

# Step 5: Find project root and tool
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Step 6: Measure baseline (tool auto-detects correct base)
$PROJECT_ROOT/tools/line-counter.sh  # No parameters needed!

# Log baseline
echo "BASELINE: $(date '+%Y-%m-%d %H:%M:%S') - Starting size measurement" >> work-log.md
```

### 2. Continuous Monitoring (DURING)

**MANDATORY MEASUREMENT TRIGGERS:**

| Trigger | Frequency | Penalty for Missing |
|---------|-----------|-------------------|
| Time-based | Every 60 minutes | -15% |
| Change-based | Every ~100 lines | -20% |
| Commit-based | Before EVERY commit | -25% |
| Feature-based | After completing component | -20% |
| Warning-based | When approaching thresholds | -30% |

### 3. Threshold Actions

```bash
# Decision matrix based on measurements
measure_and_decide() {
    LINES=$1
    
    if [ "$LINES" -ge 800 ]; then
        echo "🚨🚨🚨 STOP: AUTOMATIC FAILURE - Size limit exceeded!"
        echo "ACTION: DO NOT CONTINUE - Transition to ERROR state"
        exit 1
    elif [ "$LINES" -ge 750 ]; then
        echo "🚨 CRITICAL: 50 lines from failure!"
        echo "ACTION: Complete immediately or split NOW"
        return 2
    elif [ "$LINES" -ge 700 ]; then
        echo "⚠️ WARNING: Approaching limit"
        echo "ACTION: Plan completion strategy, consider optimization"
        return 1
    elif [ "$LINES" -ge 600 ]; then
        echo "📊 CAUTION: Substantial size"
        echo "ACTION: Increase monitoring frequency to every 30 minutes"
        return 0
    else
        echo "✅ SAFE: Continue implementation"
        return 0
    fi
}
```

### 4. Base Branch Determination

**CRITICAL: Using wrong base branch = -100% AUTOMATIC FAILURE**

| Context | Correct Base Branch | NEVER Use |
|---------|-------------------|-----------|
| Normal Effort | phase[N]/integration | main, master |
| Split Effort | Original effort branch | main, phase/integration |
| Fix Work | Current integration | main, master |
| New Feature | Phase integration | main, develop |

### 5. Measurement Documentation

**EVERY measurement MUST be logged:**

```yaml
# work-log.md entry format
- [2025-01-20 14:30:00] Size measurement:
  Tool: ${PROJECT_ROOT}/tools/line-counter.sh
  Command: ./tools/line-counter.sh -b phase1/integration -c phase1/wave1/E1.1.1
  Base: phase1/integration
  Current: phase1/wave1/E1.1.1
  Result: 452/800 lines (56.5%)
  Growth: +87 lines since last measurement
  Decision: Continue with increased monitoring
```

## Integration with Other Rules

### Relationship to R304
- R304 defines Code Reviewer requirements
- R305 ensures SW Engineers MEET those requirements proactively
- Both use identical measurement standards

### Relationship to R007
- R007 defines the 800-line limit
- R305 ensures continuous compliance monitoring

### Relationship to R200
- R200 requires measuring only changeset
- R305 enforces correct base branch usage

## Implementation Automation

### Automated Monitoring Script
```bash
#!/bin/bash
# sw-engineer-monitor.sh - Add to implementation workflow

MONITOR_INTERVAL=3600  # 1 hour
LAST_MEASUREMENT_FILE=".last-measurement"
MEASUREMENT_LOG="work-log.md"

should_measure() {
    if [ ! -f "$LAST_MEASUREMENT_FILE" ]; then
        return 0  # First measurement
    fi
    
    LAST=$(cat "$LAST_MEASUREMENT_FILE")
    NOW=$(date +%s)
    ELAPSED=$((NOW - LAST))
    
    # Time trigger
    if [ $ELAPSED -gt $MONITOR_INTERVAL ]; then
        echo "⏰ Time-based measurement trigger"
        return 0
    fi
    
    # Change trigger
    CHANGES=$(git diff --numstat | awk '{sum+=$1} END {print sum}')
    if [ "$CHANGES" -gt 100 ]; then
        echo "📝 Change-based measurement trigger"
        return 0
    fi
    
    return 1
}

perform_measurement() {
    echo "📊 PERFORMING REQUIRED MEASUREMENT"
    
    # [Insert full measurement script from above]
    
    # Update last measurement
    date +%s > "$LAST_MEASUREMENT_FILE"
    
    # Log to work log
    echo "- [$(date '+%Y-%m-%d %H:%M:%S')] Automated measurement: $LINES/800" >> "$MEASUREMENT_LOG"
}

# Main monitoring loop
while true; do
    if should_measure; then
        perform_measurement
    fi
    sleep 300  # Check every 5 minutes
done
```

## Grading Impact

| Violation | Penalty | Recovery |
|-----------|---------|----------|
| No baseline measurement | -20% | Cannot recover |
| Missing hourly measurement | -15% each | Document reason |
| No measurement before commit | -25% | Add pre-commit hook |
| Wrong base branch used | -100% | AUTOMATIC FAILURE |
| Exceeding 800 lines | -100% | AUTOMATIC FAILURE |
| Poor monitoring discipline | -15% cumulative | Improve immediately |

## Common Failures to Avoid

### FAILURE 1: "I'll measure at the end"
- **Problem**: Discovers 1200 lines at completion
- **Result**: -100% AUTOMATIC FAILURE
- **Prevention**: Measure every hour minimum

### FAILURE 2: "I'll use main as base"
- **Problem**: Counts entire phase as effort
- **Result**: -100% (R304 violation)
- **Prevention**: Always use phase integration

### FAILURE 3: "The tool isn't working"
- **Problem**: Can't find line-counter.sh
- **Result**: Implementation proceeds blind
- **Prevention**: Use universal finder script

### FAILURE 4: "Using directory name as branch"
- **Problem**: ./line-counter.sh -b main -c my-effort-directory
- **Result**: Tool fails - "my-effort-directory" is NOT a branch!
- **Prevention**: Use git branch --show-current for actual branch name

### FAILURE 5: "Running from wrong location"
- **Problem**: Running from main repo, not effort repo
- **Result**: Can't find effort branches, wrong measurements
- **Prevention**: CD into effort directory first - it's a separate git repo!

### FAILURE 4: "I forgot to measure"
- **Problem**: 3 hours pass without measurement
- **Result**: -45% penalty (3 x -15%)
- **Prevention**: Automated monitoring script

## Success Pattern

```bash
# SW Engineer entering IMPLEMENTATION state:

1. STARTUP:
   - Measure baseline
   - Set up monitoring automation
   - Document in work-log.md

2. EVERY HOUR:
   - Run measurement
   - Check thresholds
   - Adjust strategy if needed

3. BEFORE COMMITS:
   - Final measurement
   - Verify under limits
   - Document in commit message

4. AT THRESHOLDS:
   - 600 lines: Increase monitoring
   - 700 lines: Plan completion
   - 750 lines: Emergency completion
   - 800 lines: NEVER REACHED (stopped at 750)
```

## Summary

**THE CARDINAL RULE**: It is better to measure 100 times unnecessarily than to exceed the limit once.

**REMEMBER**:
- You are RESPONSIBLE for size compliance
- Continuous monitoring is MANDATORY
- Wrong base branch = AUTOMATIC FAILURE
- Exceeding 800 lines = AUTOMATIC FAILURE
- Proactive > Reactive

This rule ensures SW Engineers succeed by preventing violations rather than discovering them.