# Rule R201: Line Counter Tool Location and Usage

## Rule Statement
The line counter tool is located in the PROJECT ROOT's tools folder. Find the project root (where orchestrator-state-v3.json lives), then use `tools/line-counter.sh` from there. The tool must be run WITH NO PARAMETERS from within the effort directory.

## Criticality Level
**BLOCKING** - Using wrong tool or wrong location causes measurement failures

## Enforcement Mechanism
- **Technical**: Tool exists in project tools folder
- **Behavioral**: STOP if tool not found in expected location
- **Grading**: -30% for using wrong measurement methods

## Core Principle

```
Line Counter Location: ${PROJECT_ROOT}/tools/line-counter.sh
Project Root: Where orchestrator-state-v3.json exists
NOT: /workspaces/kcp-shared-tools/ (outdated)
NOT: ./tools/ relative paths (won't work from efforts)
NOT: utilities/ folder (wrong folder)
Usage: NO PARAMETERS (auto-detects everything)
```

## Detailed Requirements

### CORRECT Tool Location and Usage

```bash
# ✅✅✅ CORRECT - Find project root, then use utilities folder
measure_effort_size() {
    # Step 1: Find the project root (where orchestrator-state-v3.json exists)
    # Start from current directory and search upward
    PROJECT_ROOT=$(pwd)
    while [ "$PROJECT_ROOT" != "/" ]; do
        if [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ]; then
            echo "Found project root: $PROJECT_ROOT"
            break
        fi
        PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
    done
    
    # Step 2: Verify tool exists in tools folder
    LINE_COUNTER="$PROJECT_ROOT/tools/line-counter.sh"
    if [ ! -f "$LINE_COUNTER" ]; then
        echo "❌ ERROR: Line counter not found at $LINE_COUNTER"
        echo "FATAL: Cannot find line-counter.sh in tools folder"
        echo "Expected location: $PROJECT_ROOT/tools/line-counter.sh"
        exit 1
    fi
    echo "✅ Found line counter at: $LINE_COUNTER"
    
    # Step 3: Run from effort directory with NO parameters
    # MUST be IN the effort directory when running
    cd /home/vscode/workspaces/project/efforts/phase1/wave1/my-effort
    $LINE_COUNTER  # NO PARAMETERS!
    # Output: Counting lines in phase1/wave1/my-effort (excluding generated code)...
    # Total non-generated lines: 456
}
```

### WRONG Locations and Usage

```bash
# ❌❌❌ WRONG - Old location (doesn't exist)
$PROJECT_ROOT/tools/line-counter.sh

# ❌❌❌ WRONG - Wrong tool name
./tools/line-counter.sh

# ❌❌❌ WRONG - Using parameters
./tools/line-counter.sh -c branch-name  # NO! No parameters!
./tools/line-counter.sh --detailed      # NO! No flags!
./tools/line-counter.sh phase1/wave1    # NO! No arguments!
```

## Tool Behavior

The line counter tool (`./tools/line-counter.sh`):
1. Auto-detects the current git branch
2. Automatically excludes generated code patterns:
   - `zz_generated*.go`
   - `*.pb.go`
   - `*_generated.go`
   - CRD YAML files
   - SDK client code
3. Counts only human-written code
4. Requires NO parameters or flags

## If Tool Not Found

```bash
# Check if orchestrator copied it
ls -la ./tools/

# If missing, orchestrator failed to set up workspace properly
echo "ERROR: Orchestrator did not copy line-counter.sh to ./tools/"
echo "This is an orchestrator workspace setup failure"
echo "Cannot proceed without proper tool"
exit 1
```

## Manual Fallback (Emergency Only)

```bash
# ONLY if tool is completely unavailable and you need rough count
# This is NOT the proper method and may be inaccurate
find pkg/ -name "*.go" -type f ! -name "*generated*" ! -name "*.pb.go" | 
    xargs wc -l | tail -1
```

## Common Errors

### Error: "kcp-shared-tools not found"
**Cause**: Looking in wrong location
**Fix**: Use `./tools/line-counter.sh` in effort directory

### Error: "fatal: bad revision"
**Cause**: Passing parameters to the tool
**Fix**: Run with NO parameters

### Error: "not a git repository"
**Cause**: Not in the effort directory
**Fix**: cd to effort directory first

### Error: Tool not found
**Cause**: Orchestrator didn't copy it
**Fix**: This is an orchestrator failure - report and stop

## Integration with Other Rules

- **R198**: Line counter usage (NO parameters)
- **R200**: Measure only changeset (tool handles this)
- **R007**: Size limits (tool provides the measurement)

## Grading Impact

- **Using wrong tool location**: -30% (Major error)
- **Using parameters with tool**: -20% (Usage violation)
- **Manual counting instead of tool**: -15% (Process violation)
- **Not measuring at all**: -40% (Critical failure)

## Summary

**Remember**:
- Tool is ALWAYS at `./tools/line-counter.sh`
- NEVER at `/workspaces/kcp-shared-tools/`
- Run with NO parameters whatsoever
- Tool auto-detects everything it needs
- If tool is missing, it's an orchestrator failure