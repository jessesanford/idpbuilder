# 🚨 RULE R022 - Architect Size Verification Protocol

**Criticality:** BLOCKING - Incorrect measurement = invalid review  
**Grading Impact:** -30% for measurement errors  
**Enforcement:** EVERY effort review requires proper measurement

## Rule Statement

Architects MUST verify effort sizes by navigating to each effort directory and running the line counter with NO PARAMETERS. The tool is located at `${PROJECT_ROOT}/tools/line-counter.sh` where PROJECT_ROOT is the orchestrator's directory (containing `orchestrator-state.yaml`).

## 🚨🚨🚨 CRITICAL: Architects Don't Pass Parameters! 🚨🚨🚨

### ❌ WRONG - What Architects Are Doing Wrong
```bash
# ❌ CATASTROPHIC ERROR - Passing effort names as parameters!
./tools/line-counter.sh phase2/wave2/effort2-optimizer  # WRONG!
$CLAUDE_PROJECT_DIR/tools/line-counter.sh phase2/wave2/effort2-optimizer  # WRONG!

# ❌ ALSO WRONG - Using non-existent tool names
$PROJECT_ROOT/tools/line-counter.sh  # This tool doesn't exist!
```

### ✅ CORRECT - How Architects Must Measure

```bash
# For each effort you're reviewing:

# 1. Navigate TO the effort directory
cd /home/vscode/workspaces/idpbuilder-oci-mgmt/efforts/phase2/wave2/effort2-optimizer

# 2. Verify you're in the right place
pwd  # Should show the effort directory
git branch --show-current  # Should show effort branch

# 3. Find the orchestrator's project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.yaml" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
echo "Project root: $PROJECT_ROOT"

# 4. Run the line counter with NO PARAMETERS
$PROJECT_ROOT/tools/line-counter.sh
# Output: Total non-generated lines: XXX
```

## Complete Architect Workflow for Wave Review

### 🚨🚨🚨 CRITICAL: Check Split Count FIRST (R297) 🚨🚨🚨
**BEFORE measuring any effort, check if it was already split!**
See: rule-library/R297-architect-split-detection-protocol.md

### Measuring All Efforts in a Wave

```bash
verify_wave_sizes() {
    local PHASE=$1
    local WAVE=$2
    local PROJECT_ROOT=$3
    
    echo "📊 ARCHITECT: Verifying sizes for Phase $PHASE Wave $WAVE"
    
    # Get list of efforts from state file
    EFFORTS=$(yq ".waves.wave${WAVE}.efforts[]" $PROJECT_ROOT/orchestrator-state.yaml)
    
    for effort_name in $EFFORTS; do
        echo ""
        echo "🔍 Checking effort: $effort_name"
        
        # R297: CHECK SPLIT_COUNT FIRST!
        SPLIT_COUNT=$(yq ".efforts_completed.\"${effort_name}\".split_count" $PROJECT_ROOT/orchestrator-state.yaml 2>/dev/null || echo "0")
        
        if [ "$SPLIT_COUNT" -gt 0 ]; then
            echo "   ✅ COMPLIANT: Already split into $SPLIT_COUNT parts"
            echo "   📦 Original branches are size-compliant"
            echo "   ℹ️ Integration will exceed limit (expected behavior)"
            continue  # Skip to next effort
        fi
        
        # Only measure if not already split
        effort_dir="$PROJECT_ROOT/efforts/phase${PHASE}/wave${WAVE}/${effort_name}"
        
        if [ ! -d "$effort_dir" ]; then
            echo "   ❌ ERROR: Effort directory not found: $effort_dir"
            continue
        fi
        
        # Navigate TO the effort directory
        cd "$effort_dir"
        
        # Verify it's a git repo with correct branch
        if ! git status &>/dev/null; then
            echo "   ❌ Not a git repository: $effort_dir"
            continue
        fi
        
        branch=$(git branch --show-current)
        echo "   Branch: $branch"
        
        # Run line counter with NO PARAMETERS
        size=$($PROJECT_ROOT/tools/line-counter.sh | grep "Total" | awk '{print $NF}')
        
        if [ -z "$size" ]; then
            echo "   ❌ ERROR: Could not measure size"
        elif [ "$size" -gt 800 ]; then
            echo "   ❌ VIOLATION: $size lines (>800 limit)"
        else
            echo "   ✅ COMPLIANT: $size lines"
        fi
    done
    
    cd "$PROJECT_ROOT"  # Return to project root
}

# Example usage:
PROJECT_ROOT="/home/vscode/workspaces/idpbuilder-oci-mgmt"
verify_wave_sizes 2 2 "$PROJECT_ROOT"
```

## Key Points for Architects

### 1. NEVER Pass Parameters to Line Counter
```bash
# ❌ WRONG - Don't do this:
./tools/line-counter.sh effort-name
./tools/line-counter.sh -c branch-name
./tools/line-counter.sh phase2/wave2/effort

# ✅ RIGHT - Just run it:
$PROJECT_ROOT/tools/line-counter.sh  # NOTHING ELSE!
```

### 2. The Tool Auto-Detects Everything
- Auto-detects current branch from git
- Auto-compares against base branch
- Auto-excludes generated code
- You just need to be IN the effort directory

### 3. $CLAUDE_PROJECT_DIR Is Often Not Set
```bash
# ❌ Don't rely on $CLAUDE_PROJECT_DIR
$CLAUDE_PROJECT_DIR/tools/line-counter.sh  # Often undefined!

# ✅ Find PROJECT_ROOT yourself
PROJECT_ROOT=$(find /home -name "orchestrator-state.yaml" -type f 2>/dev/null | head -1 | xargs dirname)
$PROJECT_ROOT/tools/line-counter.sh
```

### 4. The Tool Name Is line-counter.sh
```bash
# ❌ WRONG tool names:
line-counter.sh  # Doesn't exist!
kcp-line-counter.sh     # Doesn't exist!

# ✅ CORRECT tool name:
line-counter.sh  # This is the only tool
```

## Common Architect Errors and Solutions

### Error: "/bin/bash: line 1: /tools/line-counter.sh: No such file or directory"
```bash
# CAUSE: $CLAUDE_PROJECT_DIR is empty/unset
# SOLUTION: Find project root manually
PROJECT_ROOT=$(find /home -name "orchestrator-state.yaml" -type f 2>/dev/null | head -1 | xargs dirname)
```

### Error: "Not in a git repository or cannot detect current branch"
```bash
# CAUSE: Not in the effort directory when running tool
# SOLUTION: CD to effort directory FIRST, then run tool
cd /path/to/effort
$PROJECT_ROOT/tools/line-counter.sh
```

### Error: "fatal: ambiguous argument 'phase2/wave2/effort2-optimizer'"
```bash
# CAUSE: Passing effort name as parameter
# SOLUTION: Don't pass ANY parameters!
$PROJECT_ROOT/tools/line-counter.sh  # Just this!
```

## Integration with Other Rules

### R297 (Architect Split Detection Protocol)
**CRITICAL**: Check split_count BEFORE measuring ANY effort!
- If split_count > 0: Effort is COMPLIANT (already split)
- Integration branches will exceed limits (EXPECTED)
- Measure ORIGINAL effort branches, NOT integration

### R076 (Size Compliance Verification)
This rule supports R076 by ensuring architects measure correctly:

```yaml
wave_review_checklist:
  size_compliance:
    - Check split_count in orchestrator-state.yaml (R297)
    - If already split, mark as compliant
    - Otherwise, navigate to each effort directory
    - Run line counter with NO parameters
    - Record actual size from output
    - Flag any >800 line violations
    - Document in review report
```

## Grading Impact

- **Passing parameters to line counter**: -30% (Measurement error)
- **Not navigating to effort directory**: -20% (Context error)
- **Using wrong tool name**: -15% (Tool confusion)
- **Missing size violations**: -40% (Review failure)

## Quick Reference for Architects

```bash
# ONE-LINER for measuring an effort:
cd /path/to/effort && $(find /home -name orchestrator-state.yaml -type f 2>/dev/null | head -1 | xargs dirname)/tools/line-counter.sh

# COMPLETE WORKFLOW:
1. cd to effort directory
2. Find project root with orchestrator-state.yaml
3. Run $PROJECT_ROOT/tools/line-counter.sh (NO PARAMETERS!)
4. Check output for "Total non-generated lines"
5. Verify ≤800 lines
```

---
**Remember:** The line counter knows what to do. Just run it from the effort directory with NO parameters!