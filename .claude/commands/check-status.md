# /check-status

╔═══════════════════════════════════════════════════════════════════════════════╗
║                        SOFTWARE FACTORY 2.0                                  ║
║                      STATUS CHECKING COMMAND                                  ║
║                                                                               ║
║ Rules: COMPREHENSIVE-DIAGNOSTICS + STATE-ANALYSIS + RECOVERY-GUIDANCE        ║
║ + AGENT-VERIFICATION + PROGRESS-ASSESSMENT                                    ║
╚═══════════════════════════════════════════════════════════════════════════════╝

## 🔍 COMPREHENSIVE STATUS ANALYSIS

This command provides complete diagnostic information about the current state of the Software Factory 2.0 system.

## 📊 SYSTEM OVERVIEW

### Environment Status
```bash
echo "🌐 ENVIRONMENT STATUS CHECK:"
echo "=================================="
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "Working Directory: $(pwd)"
echo "Expected Directory: /workspaces/[project]/"
echo "Git Repository: $(git rev-parse --is-inside-work-tree 2>/dev/null || echo 'Not a git repo')"
echo "Current Branch: $(git branch --show-current 2>/dev/null || echo 'No git branch')"
echo "Remote Tracking: $(git status -sb 2>/dev/null | head -1 || echo 'No remote tracking')"
echo "Uncommitted Changes: $(git status --porcelain 2>/dev/null | wc -l) files"
echo ""
```

### Project Configuration Status
```bash
echo "📁 PROJECT CONFIGURATION STATUS:"
echo "=================================="

# Check core directories
DIRS=("agent-configs" "agent-states" "expertise" "grading-rubrics" "state-machines" "rule-library")
for dir in "${DIRS[@]}"; do
    if [[ -d "./$dir" ]]; then
        echo "✅ $dir/ exists ($(ls -1 ./$dir 2>/dev/null | wc -l) items)"
    else
        echo "❌ $dir/ missing"
    fi
done

# Check agent config structure
if [[ -d "./agent-configs/[project]" ]]; then
    echo "✅ Project config exists: ./agent-configs/[project]/"
    echo "   Files: $(ls -1 ./agent-configs/[project]/*.md 2>/dev/null | wc -l) markdown files"
    echo "   TODOs: $(ls -1 ./agent-configs/[project]/todos/*.todo 2>/dev/null | wc -l) todo files"
else
    echo "❌ Project config missing: ./agent-configs/[project]/"
fi
echo ""
```

## 🎯 STATE MACHINE ANALYSIS

### Current State Detection
```bash
echo "🎯 STATE MACHINE ANALYSIS:"
echo "=================================="

# Read orchestrator state if available
STATE_FILE="./agent-configs/[project]/orchestrator-state.yaml"
if [[ -f "$STATE_FILE" ]]; then
    echo "📄 Orchestrator State File: EXISTS"
    echo "Current State: $(grep '^current_state:' "$STATE_FILE" 2>/dev/null | cut -d':' -f2 | tr -d ' ')"
    echo "Current Phase: $(grep '^current_phase:' "$STATE_FILE" 2>/dev/null | cut -d':' -f2 | tr -d ' ')"
    echo "Current Wave: $(grep '^current_wave:' "$STATE_FILE" 2>/dev/null | cut -d':' -f2 | tr -d ' ')"
    echo "Last Update: $(grep '^last_update:' "$STATE_FILE" 2>/dev/null | cut -d':' -f2- | tr -d '"')"
    echo ""
    
    # Check for efforts in progress
    echo "📋 Efforts Status:"
    if grep -q "efforts_in_progress:" "$STATE_FILE"; then
        echo "   In Progress: $(grep -A 20 "efforts_in_progress:" "$STATE_FILE" | grep -c "  [a-zA-Z]" || echo "0")"
    fi
    if grep -q "efforts_completed:" "$STATE_FILE"; then
        echo "   Completed: $(grep -A 20 "efforts_completed:" "$STATE_FILE" | grep -c "  [a-zA-Z]" || echo "0")"
    fi
    
else
    echo "❌ Orchestrator State File: MISSING"
    echo "   Expected: $STATE_FILE"
    echo "   Status: State machine not initialized or corrupted"
fi
echo ""
```

### State Machine Consistency Check
```bash
echo "🔄 STATE MACHINE CONSISTENCY:"
echo "=================================="

# Check for state machine definition files
SM_FILES=("orchestrator.md" "sw-engineer.md" "code-reviewer.md" "architect.md")
echo "State Machine Definitions:"
for sm_file in "${SM_FILES[@]}"; do
    if [[ -f "./state-machines/$sm_file" ]]; then
        echo "✅ $sm_file exists"
    else
        echo "❌ $sm_file missing"
    fi
done
echo ""

# Check current state validity
if [[ -f "$STATE_FILE" ]]; then
    CURRENT_STATE=$(grep '^current_state:' "$STATE_FILE" 2>/dev/null | cut -d':' -f2 | tr -d ' ')
    if [[ -n "$CURRENT_STATE" ]]; then
        echo "Current State: $CURRENT_STATE"
        echo "State Validity: $(grep -q "$CURRENT_STATE" ./state-machines/*.md && echo "✅ Valid" || echo "❌ Invalid or undefined")"
    fi
fi
echo ""
```

## 📋 TODO SYSTEM STATUS

### TODO File Analysis
```bash
echo "📋 TODO SYSTEM STATUS:"
echo "=================================="

TODO_DIR="./agent-configs/[project]/todos"
if [[ -d "$TODO_DIR" ]]; then
    echo "TODO Directory: EXISTS"
    
    # Count TODO files by agent
    AGENTS=("orchestrator" "sw-eng" "code-reviewer" "architect")
    for agent in "${AGENTS[@]}"; do
        COUNT=$(ls -1 "$TODO_DIR"/${agent}-*.todo 2>/dev/null | wc -l)
        LATEST=$(ls -t "$TODO_DIR"/${agent}-*.todo 2>/dev/null | head -1)
        if [[ $COUNT -gt 0 ]]; then
            echo "✅ $agent: $COUNT files (latest: $(basename "$LATEST" 2>/dev/null))"
            if [[ -f "$LATEST" ]]; then
                MODIFIED=$(stat -f %Sm -t '%Y-%m-%d %H:%M' "$LATEST" 2>/dev/null || stat -c '%y' "$LATEST" 2>/dev/null | cut -d'.' -f1)
                echo "   Last modified: $MODIFIED"
            fi
        else
            echo "❌ $agent: No TODO files"
        fi
    done
    
    # Check for old files (cleanup needed)
    OLD_FILES=$(find "$TODO_DIR" -name "*.todo" -mtime +1 2>/dev/null | wc -l)
    if [[ $OLD_FILES -gt 0 ]]; then
        echo "⚠️  Cleanup needed: $OLD_FILES files older than 24 hours"
    fi
    
else
    echo "❌ TODO Directory: MISSING ($TODO_DIR)"
fi
echo ""
```

### TODO State Recovery Check
```bash
echo "🔄 TODO RECOVERY STATUS:"
echo "=================================="

# Check for compaction markers
if [[ -f "/tmp/compaction_marker.txt" ]]; then
    echo "⚠️  COMPACTION MARKER DETECTED"
    echo "Compaction Details:"
    cat /tmp/compaction_marker.txt | head -10
    echo ""
    
    if grep -q "TODO_STATE_SAVED:" /tmp/compaction_marker.txt; then
        echo "📋 TODO STATE WAS PRESERVED"
        echo "Recovery Actions Required:"
        echo "1. Read latest TODO files from todos directory"
        echo "2. Use TodoWrite tool to load recovered TODOs"
        echo "3. Deduplicate with existing TODOs"
        echo "4. Clear compaction marker when done"
    fi
else
    echo "✅ No active compaction markers"
fi
echo ""
```

## 🔍 AGENT-SPECIFIC DIAGNOSTICS

### Current Agent Identity
```bash
echo "👤 AGENT IDENTITY CHECK:"
echo "=================================="

# Try to detect current agent from various sources
AGENT_PROMPT_MATCH=$(env | grep -i "agent" 2>/dev/null || echo "No agent environment variables")
echo "Environment Check: $AGENT_PROMPT_MATCH"

# Check current working directory context
CURRENT_DIR=$(basename "$(pwd)")
if [[ "$CURRENT_DIR" =~ (orchestrator|sw-engineer|code-reviewer|architect) ]]; then
    echo "Directory Context: Suggests $(echo $CURRENT_DIR | sed 's/-/ /g') agent"
else
    echo "Directory Context: No specific agent context detected"
fi

echo ""
echo "🔍 IDENTITY VERIFICATION REQUIRED:"
echo "Check your current prompt for @agent-* mentions to confirm identity"
echo ""
```

### Agent State Directory Analysis
```bash
echo "📁 AGENT STATE DIRECTORIES:"
echo "=================================="

AGENT_STATES=("orchestrator" "sw-engineer" "code-reviewer" "architect")
for agent in "${AGENT_STATES[@]}"; do
    STATE_DIR="./agent-states/$agent"
    if [[ -d "$STATE_DIR" ]]; then
        echo "✅ $agent state directory exists"
        STATES=$(ls -1 "$STATE_DIR" 2>/dev/null | grep -v ".md" | wc -l)
        echo "   States defined: $STATES"
        
        # Check for specific state files with checkpoints
        CHECKPOINTS=$(find "$STATE_DIR" -name "checkpoint.md" 2>/dev/null | wc -l)
        GRADING=$(find "$STATE_DIR" -name "grading.md" 2>/dev/null | wc -l)
        RULES=$(find "$STATE_DIR" -name "rules.md" 2>/dev/null | wc -l)
        echo "   Checkpoints: $CHECKPOINTS, Grading: $GRADING, Rules: $RULES"
    else
        echo "❌ $agent state directory missing"
    fi
done
echo ""
```

## 🔧 WORK PROGRESS ASSESSMENT

### Git Branch Analysis
```bash
echo "🌿 GIT BRANCH ANALYSIS:"
echo "=================================="

# List all branches with categorization
echo "Branch Categories:"
ALL_BRANCHES=$(git branch -a 2>/dev/null | sed 's/^[* ] //' | grep -v '^remotes/origin/HEAD')

# Count different types of branches
EFFORT_BRANCHES=$(echo "$ALL_BRANCHES" | grep -c "effort-" 2>/dev/null || echo "0")
SPLIT_BRANCHES=$(echo "$ALL_BRANCHES" | grep -c "split-" 2>/dev/null || echo "0")
INTEGRATION_BRANCHES=$(echo "$ALL_BRANCHES" | grep -c "integration" 2>/dev/null || echo "0")
WAVE_BRANCHES=$(echo "$ALL_BRANCHES" | grep -c "wave" 2>/dev/null || echo "0")
PHASE_BRANCHES=$(echo "$ALL_BRANCHES" | grep -c "phase" 2>/dev/null || echo "0")

echo "  Effort branches: $EFFORT_BRANCHES"
echo "  Split branches: $SPLIT_BRANCHES"
echo "  Integration branches: $INTEGRATION_BRANCHES"
echo "  Wave branches: $WAVE_BRANCHES"
echo "  Phase branches: $PHASE_BRANCHES"
echo ""

# Check for uncommitted work
UNCOMMITTED=$(git status --porcelain 2>/dev/null | wc -l)
if [[ $UNCOMMITTED -gt 0 ]]; then
    echo "⚠️  UNCOMMITTED CHANGES: $UNCOMMITTED files"
    echo "Recent changes:"
    git status --porcelain 2>/dev/null | head -5
    echo ""
fi
```

### Line Count Assessment
```bash
echo "📏 LINE COUNT ASSESSMENT:"
echo "=================================="

# Check for line counter tool
LINE_COUNTER="./tools/[project]-line-counter.sh"
if [[ -f "$LINE_COUNTER" ]]; then
    echo "✅ Line counter tool exists: $LINE_COUNTER"
    
    # Get current branch line count
    CURRENT_BRANCH=$(git branch --show-current 2>/dev/null)
    if [[ -n "$CURRENT_BRANCH" && "$CURRENT_BRANCH" != "main" && "$CURRENT_BRANCH" != "master" ]]; then
        echo "📊 Current branch ($CURRENT_BRANCH) analysis:"
        if $LINE_COUNTER -c "$CURRENT_BRANCH" 2>/dev/null; then
            echo ""
        else
            echo "   ❌ Failed to measure current branch"
        fi
    else
        echo "📊 On main branch - no effort measurement needed"
    fi
    
    # Check effort branches for compliance
    echo "🔍 Effort branches compliance check:"
    git branch -a | grep "effort-" | head -5 | while read -r branch; do
        branch_name=$(echo "$branch" | sed 's/^[* ] //' | sed 's/remotes\/origin\///')
        count=$($LINE_COUNTER -c "$branch_name" 2>/dev/null || echo "ERROR")
        if [[ "$count" == "ERROR" ]]; then
            echo "   ❌ $branch_name: Measurement failed"
        elif [[ $count -le 800 ]]; then
            echo "   ✅ $branch_name: $count lines (compliant)"
        else
            echo "   ❌ $branch_name: $count lines (EXCEEDS LIMIT)"
        fi
    done
    
else
    echo "❌ Line counter tool missing: $LINE_COUNTER"
    echo "   Cannot verify line count compliance"
fi
echo ""
```

## 🎯 RECOVERY RECOMMENDATIONS

### Based on Current State
```bash
echo "💡 RECOVERY RECOMMENDATIONS:"
echo "=================================="

# Generate recommendations based on findings
RECOMMENDATIONS=""

if [[ ! -f "$STATE_FILE" ]]; then
    RECOMMENDATIONS="$RECOMMENDATIONS
🔥 CRITICAL: Missing orchestrator state file
   → Use /reset-state (Level 2) to reinitialize state machine
   → Or manually create orchestrator-state.yaml"
fi

if [[ ! -d "$TODO_DIR" ]] || [[ $(ls -1 "$TODO_DIR"/*.todo 2>/dev/null | wc -l) -eq 0 ]]; then
    RECOMMENDATIONS="$RECOMMENDATIONS
📋 TODO Recovery Needed:
   → Check for compaction markers with context recovery protocol
   → Initialize fresh TODO lists for active agents
   → Use appropriate /continue-* commands with startup sequences"
fi

UNCOMMITTED=$(git status --porcelain 2>/dev/null | wc -l)
if [[ $UNCOMMITTED -gt 5 ]]; then
    RECOMMENDATIONS="$RECOMMENDATIONS
🌿 Git Cleanup Needed:
   → Review uncommitted changes: git status
   → Commit work in progress: git add . && git commit
   → Consider creating backup branch before major operations"
fi

if [[ -f "/tmp/compaction_marker.txt" ]]; then
    RECOMMENDATIONS="$RECOMMENDATIONS
🔄 Compaction Recovery Required:
   → Follow context recovery protocol in CLAUDE.md
   → Read and load TODO state files using TodoWrite tool
   → Clear compaction marker after recovery"
fi

if [[ -n "$RECOMMENDATIONS" ]]; then
    echo "$RECOMMENDATIONS"
else
    echo "✅ No immediate recovery actions needed"
    echo "💡 System appears to be in good state"
fi
echo ""
```

### Next Steps Guide
```bash
echo "🎯 NEXT STEPS GUIDE:"
echo "=================================="

# Determine appropriate next actions
if [[ -f "$STATE_FILE" ]]; then
    CURRENT_STATE=$(grep '^current_state:' "$STATE_FILE" 2>/dev/null | cut -d':' -f2 | tr -d ' ')
    case "$CURRENT_STATE" in
        "INIT")
            echo "🚀 State: INIT - Ready for orchestration startup"
            echo "   → Use: /continue-orchestrating"
            echo "   → Read startup requirements and begin planning"
            ;;
        "WAVE_START"|"WAVE_COMPLETE")
            echo "🌊 State: $CURRENT_STATE - Wave management needed"
            echo "   → Use: /continue-orchestrating"
            echo "   → Check efforts_in_progress and integration branches"
            ;;
        "")
            echo "❌ State: UNDEFINED - State file exists but current_state is empty"
            echo "   → Consider: /reset-state (Level 2)"
            ;;
        *)
            echo "🔄 State: $CURRENT_STATE"
            echo "   → Use appropriate /continue-* command based on agent role"
            echo "   → Check state machine definition for valid transitions"
            ;;
    esac
else
    echo "🏗️  No state file - Fresh initialization needed"
    echo "   → Use: /reset-state (Level 2) or manually create orchestrator-state.yaml"
    echo "   → Begin with orchestrator agent using /continue-orchestrating"
fi

echo ""
echo "🤖 Agent-Specific Actions:"
echo "   Orchestrator: /continue-orchestrating"
echo "   SW Engineer: /continue-implementing"
echo "   Code Reviewer: /continue-reviewing"
echo "   Architect: /continue-architecting"
echo ""
echo "🆘 Emergency Actions:"
echo "   Full diagnostics: /check-status (this command)"
echo "   State corruption: /reset-state"
echo "   Context loss: Follow CLAUDE.md recovery protocol"
```

## 📊 SYSTEM HEALTH SUMMARY

### Overall Status
```bash
echo "📊 SYSTEM HEALTH SUMMARY:"
echo "=================================="

HEALTH_SCORE=100
ISSUES_FOUND=""

# Deduct points for various issues
[[ ! -f "$STATE_FILE" ]] && HEALTH_SCORE=$((HEALTH_SCORE - 30)) && ISSUES_FOUND="$ISSUES_FOUND State-file-missing"
[[ ! -d "$TODO_DIR" ]] && HEALTH_SCORE=$((HEALTH_SCORE - 20)) && ISSUES_FOUND="$ISSUES_FOUND TODO-dir-missing"
[[ $(git status --porcelain 2>/dev/null | wc -l) -gt 10 ]] && HEALTH_SCORE=$((HEALTH_SCORE - 10)) && ISSUES_FOUND="$ISSUES_FOUND Many-uncommitted-changes"
[[ -f "/tmp/compaction_marker.txt" ]] && HEALTH_SCORE=$((HEALTH_SCORE - 15)) && ISSUES_FOUND="$ISSUES_FOUND Compaction-recovery-needed"

# Determine health status
if [[ $HEALTH_SCORE -ge 90 ]]; then
    echo "🟢 SYSTEM HEALTH: EXCELLENT ($HEALTH_SCORE/100)"
    echo "   System is ready for normal operation"
elif [[ $HEALTH_SCORE -ge 70 ]]; then
    echo "🟡 SYSTEM HEALTH: GOOD ($HEALTH_SCORE/100)"
    echo "   Minor issues present: $ISSUES_FOUND"
elif [[ $HEALTH_SCORE -ge 50 ]]; then
    echo "🟠 SYSTEM HEALTH: FAIR ($HEALTH_SCORE/100)"
    echo "   Significant issues: $ISSUES_FOUND"
else
    echo "🔴 SYSTEM HEALTH: POOR ($HEALTH_SCORE/100)"
    echo "   Critical issues: $ISSUES_FOUND"
    echo "   Consider reset or manual recovery"
fi

echo ""
echo "Generated: $(date '+%Y-%m-%d %H:%M:%S %Z')"
echo "Status check complete."
```

This comprehensive status command provides detailed diagnostics to help determine the exact state of the Software Factory 2.0 system and provides specific guidance for recovery and next steps.