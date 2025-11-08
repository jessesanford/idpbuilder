# R216: Bash Execution Syntax Rules

## Critical: Single-Line vs Multi-Line Execution

### Rule Statement
All agents MUST use proper bash syntax when executing commands, distinguishing between documentation style (with backslashes) and execution style (with semicolons or multi-line).

### Related Supreme Laws
- **R221 (SUPREME LAW #1)**: Remember that bash ALWAYS resets to home directory - you MUST use `cd $DIR &&` prefix
- **R208 (SUPREME LAW #2)**: Orchestrator must CD before spawning agents

### When Using the Bash Tool

**Option 1: Multi-Line Execution (Preferred)**
Execute each line separately as natural bash:
```bash
echo "Checking for Master Implementation Plan..."
if [ -f "./IMPLEMENTATION-PLAN.md" ]; then
    echo "✓ Master Implementation Plan found"
    echo "Decision: PROCEED_WITH_ORCHESTRATION"
else
    echo "⚠️ No Master Implementation Plan found!"
fi
```

**Option 2: Single-Line Execution**
If you must run as a single line, use semicolons (`;`) to separate statements:
```bash
echo "Checking for Master Implementation Plan..."; if [ -f "./IMPLEMENTATION-PLAN.md" ]; then echo "✓ Master Implementation Plan found"; echo "Decision: PROCEED_WITH_ORCHESTRATION"; else echo "⚠️ No Master Implementation Plan found!"; fi
```

### Important: About Backslash Continuations

The backslash (`\`) line continuations in markdown documentation are for:
1. Making long commands readable in documentation
2. Allowing copy-paste of multi-line commands

They are NOT for single-line execution!

### Example of Documentation Style (with backslashes):
```bash
# This is for DOCUMENTATION - shows line continuation
if [ -f "config.json" ]; then
    echo "Config file exists";
    cat config.json;
fi
```

### Example of Execution Style (no backslashes):
```bash
# For ACTUAL EXECUTION - Option 1 (multi-line)
if [ -f "config.json" ]; then
    echo "Config file exists"
    cat config.json
fi

# For ACTUAL EXECUTION - Option 2 (single-line with semicolons)
if [ -f "config.json" ]; then echo "Config file exists"; cat config.json; fi
```

## Rules for Agents

1. **When using the Bash tool**: 
   - Prefer multi-line format (natural bash syntax)
   - Each statement on its own line
   - No backslashes needed

2. **When you must use single-line**:
   - Use semicolons between statements
   - Remove any backslashes from documentation examples
   - Test the syntax before running

3. **Common Pitfalls to Avoid**:
   - ❌ `echo "test" if [ -f "file" ]; then` (missing semicolon)
   - ✅ `echo "test"; if [ -f "file" ]; then`
   
   - ❌ `if [ -f "file" ]; then \ echo "found" \ fi` (backslashes in single line)
   - ✅ `if [ -f "file" ]; then echo "found"; fi`

## Example: Checking for Files

### Good Multi-Line Approach:
```bash
echo "Checking for implementation plan..."
if [ -f "./IMPLEMENTATION-PLAN.md" ]; then
    echo "Found implementation plan"
    plan_exists=true
elif [ -f "./ARCHITECT-PROMPT-IDPBUILDER-OCI.md" ]; then
    echo "Found architect prompt"
    need_architect=true
else
    echo "No planning files found"
    missing_requirements=true
fi
```

### Good Single-Line Approach:
```bash
echo "Checking for implementation plan..."; if [ -f "./IMPLEMENTATION-PLAN.md" ]; then echo "Found implementation plan"; plan_exists=true; elif [ -f "./ARCHITECT-PROMPT-IDPBUILDER-OCI.md" ]; then echo "Found architect prompt"; need_architect=true; else echo "No planning files found"; missing_requirements=true; fi
```

## Summary

- **Backslashes (`\`)**: For documentation line continuation only
- **Semicolons (`;`)**: For separating statements in single-line execution
- **Multi-line**: Preferred for the Bash tool - more readable and less error-prone

## Validation

To validate your bash syntax before execution:
```bash
# Test with bash -n (syntax check only)
bash -n -c 'your command here'

# Or write to a temp file and check
echo 'your command here' > /tmp/test.sh
bash -n /tmp/test.sh
```

## Application

This rule applies to:
- Orchestrator agents executing state checks
- SW Engineers running build and test commands
- Code Reviewers executing validation scripts
- Architect Reviewers running analysis commands
- Any agent using the Bash tool for command execution