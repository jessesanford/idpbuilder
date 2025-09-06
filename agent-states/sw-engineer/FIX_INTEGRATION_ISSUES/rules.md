# SW Engineer - FIX_INTEGRATION_ISSUES State Rules

## 🔴🔴🔴 STATE CLARITY (R295) 🔴🔴🔴

**YOU ARE IN STATE: FIX_INTEGRATION_ISSUES**
**This means you should: Fix integration issues found during phase/wave/project integration testing**

## State Context
You are fixing integration issues that were discovered when multiple efforts were integrated together. These issues were not found during individual effort testing but emerged during integration.

## 🔴🔴🔴 MANDATORY MEASUREMENT RULE - R304 🔴🔴🔴

### R304: Mandatory Line Counter Tool Enforcement
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R304-mandatory-line-counter-enforcement.md`
**Criticality**: BLOCKING - Manual counting = AUTOMATIC -100% FAILURE

**ABSOLUTE REQUIREMENTS:**
- ✅ MUST use `$CLAUDE_PROJECT_DIR/tools/line-counter.sh` for ALL line counting
- ❌ NEVER use `wc -l` or any manual counting method
- ❌ NEVER count lines any other way - this is a -100% automatic failure
- ✅ MUST specify both -b (base branch) and -c (current branch) parameters
- ✅ Base branch MUST be phase integration branch (NOT "main")

**Failure to use the line counter tool = AUTOMATIC -100% GRADE**

## 🚨🚨🚨 CRITICAL RULES FOR THIS STATE 🚨🚨🚨

---
### 🔴🔴🔴 RULE R295 - SW Engineer Spawn Clarity Protocol
**Source:** rule-library/R295-sw-engineer-spawn-clarity-protocol.md
**Criticality:** SUPREME - You MUST know exactly what to follow

YOU MUST HAVE BEEN TOLD:
1. **Your exact state**: FIX_INTEGRATION_ISSUES
2. **Your exact plan**: INTEGRATION-REPORT.md
3. **Your exact location**: In your effort directory
4. **What to ignore**: Any *-COMPLETED-*.md files

IF YOU ARE CONFUSED ABOUT WHAT TO DO, STOP AND ASK FOR CLARIFICATION!
---

---
### 🚨🚨🚨 RULE R293 - Integration Report Distribution
**Source:** rule-library/R293-integration-report-distribution-protocol.md
**Criticality:** BLOCKING - Report must be in your directory

THE INTEGRATION-REPORT.md FILE:
- Has been distributed to your effort directory by the orchestrator
- Contains ONLY the fixes you need to make for your effort
- Is the ONLY plan you should follow in this state
- Supersedes any previous fix plans

CHECK: ls -la | grep -E "INTEGRATION-REPORT|COMPLETED"
---

---
### 🚨🚨🚨 RULE R294 - Fix Plan Archival Protocol
**Source:** rule-library/R294-fix-plan-archival-protocol.md
**Criticality:** BLOCKING - Must archive old plans

BEFORE STARTING INTEGRATION FIXES:
```bash
# Archive any old fix plans to prevent confusion
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
# Archive all split plans (timestamped and legacy)
for plan in SPLIT-PLAN*.md; do
    [[ -f "$plan" && ! "$plan" =~ COMPLETED ]] && mv "$plan" "${plan%.md}-COMPLETED-${TIMESTAMP}.md"
done
# Archive all review reports (timestamped and legacy)  
for report in CODE-REVIEW-REPORT*.md; do
    [[ -f "$report" && ! "$report" =~ COMPLETED ]] && mv "$report" "${report%.md}-COMPLETED-${TIMESTAMP}.md"
done
```

AFTER COMPLETING INTEGRATION FIXES:
```bash
# Archive the integration report when done
mv INTEGRATION-REPORT.md "INTEGRATION-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md"
echo "✅ Integration fixes complete and report archived"
```
---

---
### 🔴🔴🔴 RULE R300 - Comprehensive Fix Management Protocol
**Source:** rule-library/R300-comprehensive-fix-management-protocol.md
**Criticality:** BLOCKING - NEVER modify integration branch

ALL FIXES MUST BE MADE IN YOUR EFFORT BRANCH:
- Make changes ONLY in your effort branch
- NEVER switch to or modify the integration branch
- Your fixes will be re-integrated after completion
- The integration branch is READ-ONLY to you
---

---
### 🔴🔴🔴 RULE R300 - Comprehensive Fix Management Protocol (Verification)
**Source:** rule-library/R300-comprehensive-fix-management-protocol.md
**Criticality:** SUPREME LAW - Must verify fixes are in effort branch

YOU MUST VERIFY YOUR FIXES:
- After applying fixes, verify they're committed to effort branch
- Push fixes to remote effort branch IMMEDIATELY
- Verify remote has your fixes before marking complete
- Integration will fail if fixes aren't in effort branches

VERIFICATION STEPS:
```bash
# After fixing, MUST verify:
git branch --show-current  # Must be effort branch, not integration
git log --oneline -1       # Must show your fix commit
git push origin $(git branch --show-current)  # Push to remote
git fetch origin
git log origin/$(git branch --show-current) --oneline -1  # Verify remote has fix
```
---

## Required Actions in FIX_INTEGRATION_ISSUES State

### 1. Verify Your Environment and Plan
```bash
# Confirm you're in the right place
pwd  # Should be in your effort directory
git branch --show-current  # Should be your effort branch

# Check for the integration report
if [[ ! -f "INTEGRATION-REPORT.md" ]]; then
    echo "❌ ERROR: INTEGRATION-REPORT.md not found!"
    echo "The orchestrator should have distributed this file."
    exit 1
fi

# Archive old plans (R294) - handle both formats
for pattern in "SPLIT-PLAN*.md" "CODE-REVIEW-REPORT*.md" "FIX-INSTRUCTIONS*.md"; do
    for old_plan in $pattern; do
        if [[ -f "$old_plan" && ! "$old_plan" =~ COMPLETED ]]; then
        mv "$old_plan" "${old_plan%.md}-COMPLETED-$(date +%Y%m%d-%H%M%S).md"
        echo "📦 Archived old plan: $old_plan"
    fi
done
```

### 2. Read and Understand Integration Report
```bash
# Read the integration report
cat INTEGRATION-REPORT.md

# Find your effort's section
grep -A 20 "effort-$(basename $(pwd))" INTEGRATION-REPORT.md || \
    grep -A 20 "$(basename $(pwd))" INTEGRATION-REPORT.md
```

### 3. Implement Integration Fixes

**Focus Areas for Integration Issues:**
1. **Dependency conflicts**: Version mismatches, missing dependencies
2. **API incompatibilities**: Interface changes, parameter mismatches
3. **Configuration conflicts**: Overlapping configs, missing settings
4. **Build issues**: Compilation errors when combined
5. **Test failures**: Integration test failures
6. **Resource conflicts**: Port conflicts, file collisions

```python
def implement_integration_fixes():
    """Implement fixes from INTEGRATION-REPORT.md"""
    
    # Read the integration report
    with open('INTEGRATION-REPORT.md', 'r') as f:
        report = f.read()
    
    # Extract fixes for this effort
    effort_name = os.path.basename(os.getcwd())
    fixes_section = extract_effort_section(report, effort_name)
    
    if not fixes_section:
        print(f"✅ No integration issues for {effort_name}")
        return True
    
    # Parse required fixes
    fixes = parse_required_fixes(fixes_section)
    
    for fix in fixes:
        print(f"🔧 Applying fix: {fix['description']}")
        
        if fix['type'] == 'dependency':
            fix_dependency_issue(fix)
        elif fix['type'] == 'api_compatibility':
            fix_api_compatibility(fix)
        elif fix['type'] == 'configuration':
            fix_configuration_conflict(fix)
        elif fix['type'] == 'build':
            fix_build_issue(fix)
        elif fix['type'] == 'test':
            fix_test_failure(fix)
        else:
            apply_generic_fix(fix)
        
        # Verify fix was successful
        if not verify_fix(fix):
            print(f"❌ Fix failed: {fix['description']}")
            return False
    
    print("✅ All integration fixes applied successfully")
    return True
```

### 4. Verify Integration Fixes
```bash
# Run build to ensure compilation works
make build || npm run build || go build ./... || cargo build

# Run tests to verify fixes
make test || npm test || go test ./... || cargo test

# Check for any remaining integration markers
grep -r "INTEGRATION_ISSUE" . --exclude-dir=.git
grep -r "TODO.*integration" . --exclude-dir=.git

# Verify specific integration requirements from report
# (Check the INTEGRATION-REPORT.md for specific verification steps)
```

### 5. Complete Integration Fix Cycle
```bash
# Archive the integration report (R294)
mv INTEGRATION-REPORT.md "INTEGRATION-REPORT-COMPLETED-$(date +%Y%m%d-%H%M%S).md"

# Create completion marker
cat > FIX_COMPLETE.flag << EOF
Integration fixes completed: $(date)
State: FIX_INTEGRATION_ISSUES
Report: INTEGRATION-REPORT-COMPLETED-*.md
All issues resolved: YES
Build status: PASSING
Test status: PASSING
EOF

# Commit all changes
git add -A
git commit -m "fix(integration): resolve integration issues from INTEGRATION-REPORT

- Fixed dependency conflicts
- Resolved API compatibility issues
- Updated configurations
- All tests passing"

git push

echo "✅ Integration fixes complete and pushed"
```

## Common Integration Issues and Solutions

### Dependency Conflicts
```bash
# Example: Resolve conflicting package versions
# If integration report says: "Conflict: package X requires Y@1.0, but Z requires Y@2.0"

# Solution 1: Update to compatible version
npm install Y@1.5.0  # Version that satisfies both

# Solution 2: Use resolutions (package.json)
"resolutions": {
  "Y": "1.5.0"
}
```

### API Compatibility
```python
# Example: Fix API parameter mismatch
# If integration report says: "API mismatch: service A expects 'user_id', service B sends 'userId'"

# Solution: Add compatibility layer
def normalize_user_request(request):
    """Normalize user request parameters for compatibility"""
    if 'userId' in request and 'user_id' not in request:
        request['user_id'] = request['userId']
    return request
```

### Configuration Conflicts
```yaml
# Example: Resolve port conflicts
# If integration report says: "Port conflict: services A and B both use port 8080"

# Solution: Update service B configuration
server:
  port: 8081  # Changed from 8080 to avoid conflict
```

## State Transitions

From FIX_INTEGRATION_ISSUES state:
- **FIXES_COMPLETE** → COMPLETED (All integration issues resolved)
- **BUILD_STILL_FAILING** → Continue fixing
- **TESTS_STILL_FAILING** → Continue fixing or TEST_WRITING
- **BLOCKED_BY_OTHER_EFFORT** → Report blockage to orchestrator

## Success Criteria

✅ **Integration fixes are complete when:**
1. All issues from INTEGRATION-REPORT.md are resolved
2. Build passes with all efforts integrated
3. Integration tests pass
4. No new issues introduced
5. INTEGRATION-REPORT.md is archived as COMPLETED
6. Changes committed and pushed

## Violations to Avoid

❌ **Common mistakes:**
- Following old fix plans instead of INTEGRATION-REPORT.md
- Modifying the integration branch directly (R300 violation)
- Not archiving completed plans (R294 violation)
- Implementing features instead of just fixes
- Not verifying fixes work in integration context

---

**REMEMBER**: 
- You are in FIX_INTEGRATION_ISSUES state
- Follow ONLY INTEGRATION-REPORT.md
- Make fixes ONLY in your effort branch
- Archive the report when complete
