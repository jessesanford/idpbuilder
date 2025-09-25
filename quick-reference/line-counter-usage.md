# Line Counter Quick Reference

## 🚨 CRITICAL: Line Counter Location

The line counter is ALWAYS at: `${PROJECT_ROOT}/tools/line-counter.sh`

## Simple Usage Pattern

```bash
# Copy and use this exact snippet:

# 1. Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    if [ -f "$PROJECT_ROOT/orchestrator-state.json" ]; then
        break
    fi
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# 2. Run line counter (NO PARAMETERS!)
$PROJECT_ROOT/tools/line-counter.sh
```

## What NOT to Do

❌ NEVER look at `/workspaces/kcp-shared-tools/` (outdated)  
❌ NEVER use `./tools/line-counter.sh` (relative won't work)  
❌ NEVER use parameters like `-c branch-name` (auto-detects)  
❌ NEVER look in utilities folder (it's in tools)  

## Remember

- Tool is in: `PROJECT_ROOT/tools/` folder
- Project root: Where `orchestrator-state.json` exists
- Run from: Your effort directory
- Parameters: NONE (auto-detects everything)

## If Not Found

If `$PROJECT_ROOT/tools/line-counter.sh` doesn't exist:
- The project setup is incomplete
- Cannot proceed without measurement
- Report error and stop

## Example Output

```
Counting lines in phase1/wave1/api-types (excluding generated code)...
Total non-generated lines: 456
```