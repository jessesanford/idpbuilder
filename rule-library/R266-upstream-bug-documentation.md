# 🔴🔴🔴 RULE R266: Upstream Bug Documentation and Fix Protocol 🔴🔴🔴

## Rule Definition
**Criticality:** SUPREME - DOCUMENT THEN FIX
**Category:** Agent-Specific
**Applies To:** integration-agent (documentation), orchestrator (fix coordination)

## TWO-STEP PROTOCOL: DOCUMENT FIRST, FIX SECOND

### Step 1: The Integration Agent Documents
**INTEGRATION AGENT: YOU ARE AN INTEGRATOR, NOT A DEVELOPER**

- **NEVER** fix bugs in integrated code
- **NEVER** patch failing tests
- **NEVER** correct compilation errors
- **NEVER** modify source to make it work
- **NEVER** create adapter/wrapper code (R361)
- **NEVER** add new packages or files (R361)
- **NEVER** create "glue code" to make things work (R361)
- **ONLY** document what you find
- **ONLY** resolve conflicts between existing code

### Step 2: The Orchestrator Coordinates Fixes
**ORCHESTRATOR: BUGS MUST BE FIXED BEFORE PROCEEDING**

After Integration Agent documents bugs:
- **MUST** transition to SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING state
- **MUST** spawn Code Reviewer to create fix plans for all documented bugs
- **MUST** wait for plans in WAITING_FOR_PROJECT_FIX_PLANS state
- **MUST** spawn SW Engineers to fix in source branches
- **MUST** re-run integration after fixes complete
- **NEVER** proceed to validation with known bugs

### Why This Rule Exists
1. **Preserve Accountability** - Original authors must fix their code
2. **Maintain Audit Trail** - Changes must be traceable to source
3. **Prevent Hidden Issues** - Fixes might mask deeper problems
4. **Clear Responsibility** - Integration ≠ Development

## Bug Documentation Format

```markdown
## UPSTREAM BUGS IDENTIFIED
**⚠️ NOT FIXED - DOCUMENTATION ONLY ⚠️**

### Bug #1: [Descriptive Title]
**Severity**: CRITICAL / HIGH / MEDIUM / LOW
**Type**: COMPILATION / RUNTIME / TEST / LOGIC
**Source Branch**: feature-xyz
**Discovered During**: Integration of feature-xyz with main

#### Location
- File: `src/api/handler.go`
- Line: 234-237
- Function: `HandleUserRequest()`

#### Description
Clear description of what's wrong and why it fails.

#### Error Output
```
panic: runtime error: index out of range [5] with length 3
goroutine 1 [running]:
main.HandleUserRequest(0xc00010e000, 0xc00010e000)
    /src/api/handler.go:235 +0x92
```

#### Impact
- Affects: User authentication flow
- Frequency: Occurs on every third request
- Workaround: None available

#### Root Cause Analysis
The function assumes array has at least 6 elements but validation only ensures 3.

#### Recommended Fix
```go
// Instead of:
userData[5] = request.Extra  // PANIC HERE

// Should be:
if len(userData) > 5 {
    userData[5] = request.Extra
}
```

#### Integration Impact
- ❌ Blocks integration completion
- ⚠️ Can integrate but with known issues
- ✅ Does not affect integration

---

### Bug #2: [Next Bug Title]
[Same format...]
```

## Enforcement Protocol

```bash
# FATAL: Detect if integration agent modified source
detect_source_modifications() {
    local integration_branch="$1"
    local allowed_files="INTEGRATION-REPORT.md work-log.md INTEGRATION-PLAN.md"
    
    # Get modified files
    modified=$(git diff --name-only main.."$integration_branch")
    
    for file in $modified; do
        # Check if file is in allowed list
        if ! echo "$allowed_files" | grep -q "$file"; then
            # Check if it's a source file
            if [[ "$file" == *.go || "$file" == *.js || "$file" == *.py ]]; then
                echo "🔴🔴🔴 FATAL: Integration agent modified source file: $file"
                echo "Integration agents must NEVER fix bugs!"
                exit 1
            fi
        fi
    done
}
```

## Examples of Violations

```bash
# ❌❌❌ WRONG - Never do this!
# Found bug in test
vim src/auth/auth_test.go
# Fixed timeout issue  
git add src/auth/auth_test.go
git commit -m "fix: increase test timeout"  # VIOLATION!

# ❌❌❌ WRONG - Never patch compilation errors!
# Build fails
echo "Fixing undefined variable..."
vim src/api/types.go
# Added missing field
git commit -m "fix: add missing field"  # VIOLATION!

# ✅✅✅ CORRECT - Document only!
cat >> INTEGRATION-REPORT.md << 'EOF'
### Bug: Test Timeout Too Short
- File: src/auth/auth_test.go:45
- Issue: Timeout set to 1s, needs 5s minimum
- Impact: Tests fail intermittently
- Recommendation: Increase timeout to 5s
- STATUS: NOT FIXED (upstream bug)
EOF
```

## Bug Severity Classification

| Severity | Description | Example |
|----------|-------------|---------|
| CRITICAL | Prevents integration | Compilation fails |
| HIGH | Major functionality broken | Core API crashes |
| MEDIUM | Feature partially works | Edge case failures |
| LOW | Minor issues | Formatting problems |

## Orchestrator Responsibilities

When monitoring project integration:
1. **CHECK** for "UPSTREAM BUGS IDENTIFIED" section in report
2. **COUNT** number of bugs documented
3. **TRANSITION** to SPAWN_CODE_REVIEWER_PROJECT_FIX_PLANNING if bugs > 0
4. **SPAWN** Code Reviewer to create fix plans for each bug
5. **WAIT** for plans in WAITING_FOR_PROJECT_FIX_PLANS state
6. **SPAWN** SW Engineers to fix in source branches
6. **MONITOR** fix progress
7. **RE-RUN** integration after fixes complete

**VIOLATION**: Proceeding to validation with documented bugs = -100% FAILURE

## Grading Impact
- **-50% INSTANT PENALTY** if Integration Agent fixes any bug
- **-100% INSTANT FAILURE** if Orchestrator proceeds with known bugs
- **+10%** for comprehensive bug documentation
- **+5%** for clear reproduction steps
- **+5%** for actionable recommendations
- **+20%** for proper fix coordination by Orchestrator

## Related Rules
- R262 - Merge Operation Protocols (don't modify originals)
- R263 - Integration Documentation Requirements
- R267 - Integration Agent Grading Criteria
- R361 - Integration Conflict Resolution Only (no new code/packages)
- R321 - Immediate Backport During Integration (fixes go to source)