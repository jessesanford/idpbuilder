# 🔴🔴🔴 SUPREME RULE R521: Integration Known Fixes & Conflict Resolution Scope

## Criticality: SUPREME LAW
**Violation = -100% AUTOMATIC FAILURE**

## Description
This rule defines the EXACT SCOPE of what integration agents CAN fix (known fixes, conflict resolution) versus what they CANNOT fix (new bugs, upstream issues). It provides the protocol for applying proven solutions from previous integrations while maintaining the integrity of R266 (document bugs) and R361 (no new code).

## 🎯 Integration Agent Decision Flow

When integration fails, follow this decision tree:

```
Integration Failure Detected
         |
         v
   ┌─────────────────────────────┐
   │ R521: Known Fix Available?  │
   │ (Check BUILD-FIX-SUMMARY)   │
   └─────────────────────────────┘
         |
         +--YES--> Apply known fix
         |         Document as conflict resolution
         |         CONTINUE (no new bug)
         |
         +--NO---> Continue to duplicate check
                   |
                   v
         ┌─────────────────────────────┐
         │ R522: Duplicate Bug Exists? │
         │ (Search bug_registry)       │
         └─────────────────────────────┘
                   |
                   +--YES--> R523: Link as duplicate
                   |         Update canonical bug
                   |         CONTINUE (no new bug)
                   |
                   +--NO---> Create new bug
                             |
                             v
                   ┌─────────────────────────────┐
                   │ R406: Track Cascade Deps    │
                   │ (Integration cascade)       │
                   └─────────────────────────────┘
                             |
                             v
                   Set CONTINUE-SOFTWARE-FACTORY=FALSE
                   Escalate to orchestrator

After Bug Fixed Upstream:
         |
         v
   ┌─────────────────────────────┐
   │ R350: Calculate Cascade     │
   │ R348: CASCADE_REINTEGRATION │
   │ R351: Execute Cascade       │
   └─────────────────────────────┘
         |
         v
   ┌─────────────────────────────┐
   │ R524: Propagate Status      │
   │ (Update all duplicates)     │
   └─────────────────────────────┘
         |
         v
   Integration Complete
```

**Key Decision Points**:
1. **Known fix available?** → Apply it (R521), no bug needed
2. **Duplicate bug exists?** → Link it (R522/R523), no new bug
3. **Neither?** → Create bug with cascade tracking (R406)
4. **Bug fixed?** → Integration cascade (R348/R350/R351) + status propagation (R524)

## 🔴🔴🔴 THE CRITICAL DISTINCTION 🔴🔴🔴

### Integration Agent CAN Fix (Conflict Resolution):
✅ **Known Fixes from BUILD-FIX-SUMMARY.md** - Apply proven solutions from previous integrations
✅ **Duplicate File Conflicts** - Delete duplicate files from merge conflicts
✅ **Import Path Conflicts** - Reconcile import statements (<10 lines)
✅ **Merge Marker Cleanup** - Remove conflict markers and select correct version
✅ **Demo Infrastructure** - Create demos per R291 at ALL levels (wave, phase, project) - R361 exception applies universally

### Integration Agent CANNOT Fix (Development):
❌ **New Bugs** - Bugs discovered for the first time during this integration
❌ **Upstream Issues** - Problems in effort branch code
❌ **Failing Tests** - Tests that fail due to code logic issues
❌ **Compilation Errors** - Errors from effort branch code
❌ **New Packages/Adapters** - Any code that doesn't exist in effort branches

## 🔍 THE KNOWN FIXES PROTOCOL

### Step 1: Search for Previous Fixes
```bash
# MANDATORY: Before documenting ANY bug, search for known fixes
search_known_fixes() {
    local bug_signature="$1"  # e.g., "BUG-007: Duplicate PushCmd"

    echo "🔍 R521: Searching for known fixes..."

    # Search all BUILD-FIX-SUMMARY.md files in current wave/phase
    for summary in .software-factory/*/BUILD-FIX-SUMMARY-*.md; do
        if [ -f "$summary" ]; then
            echo "📋 Checking: $summary"

            # Look for matching bug signature
            if grep -q "$bug_signature" "$summary"; then
                echo "✅ FOUND: Known fix exists in $summary"

                # Extract fix details
                grep -A 20 "$bug_signature" "$summary" > /tmp/known-fix.txt

                echo "📝 Fix details:"
                cat /tmp/known-fix.txt

                return 0  # Known fix found
            fi
        fi
    done

    echo "❌ No known fix found - this is a NEW bug"
    return 1  # New bug - must document per R266
}
```

### Step 2: Apply Known Fix (ALLOWED)
```bash
# If known fix found, integration agent CAN apply it
apply_known_fix() {
    local bug_signature="$1"
    local fix_file="/tmp/known-fix.txt"

    echo "🔧 R521: Applying known fix for $bug_signature"

    # Extract the fix from BUILD-FIX-SUMMARY.md
    # Example: "Fixed by deleting pkg/gitea/duplicate.go"
    FIX_ACTION=$(grep "Fixed by:" "$fix_file" | cut -d: -f2)

    echo "Applying fix: $FIX_ACTION"

    # Apply the EXACT same fix
    # This is conflict resolution, NOT development
    eval "$FIX_ACTION"

    # Document that we applied a known fix
    cat >> .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md << EOF

### Known Fix Applied: $bug_signature
**Source**: $summary
**Action**: $FIX_ACTION
**Justification**: This bug was already fixed in previous integration
**Classification**: Conflict resolution per R521 (NOT bug fixing per R266)
EOF

    echo "✅ Known fix applied successfully"
}
```

### Step 3: Document New Bug (REQUIRED)
```bash
# If NO known fix found, MUST document per R266
document_new_bug() {
    local bug_signature="$1"
    local bug_details="$2"

    echo "📝 R266: No known fix - documenting as NEW bug"

    # Follow R266 protocol exactly
    cat >> .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md << EOF

### NEW BUG DISCOVERED: $bug_signature
**Severity**: [CRITICAL/HIGH/MEDIUM/LOW]
**Type**: [COMPILATION/RUNTIME/TEST/LOGIC]
**Classification**: NEW (not in BUILD-FIX-SUMMARY.md)

$bug_details

**ACTION REQUIRED**:
- Integration agent CANNOT fix this (R266)
- Orchestrator MUST spawn fix cascade (R300)
- SW Engineers fix in effort branches
- Re-attempt integration after fixes
EOF

    echo "❌ New bug documented - integration agent stops here"
    return 1  # Signal that new bug found
}
```

## 📋 DECISION FLOWCHART

```
Integration encounters issue
         ↓
Is it a merge conflict?
    YES → Apply conflict resolution (select version) ✅
    NO → Continue
         ↓
Is it a duplicate file?
    YES → Delete duplicate ✅
    NO → Continue
         ↓
Is it a build/test failure?
    YES → Search BUILD-FIX-SUMMARY.md
         ↓
    Found in BUILD-FIX-SUMMARY.md?
        YES → Apply known fix per R521 ✅
        NO → Document as NEW bug per R266 ❌
         ↓
    Demo missing?
        YES → Create demo per R291 ✅
        NO → Continue
         ↓
Integration complete
```

## 🔄 CASCADE RETRY AWARENESS

### Detecting Cascade Context
```bash
# Check if this is a cascade retry (post-fixes)
detect_cascade_context() {
    # Check orchestrator state for cascade mode
    IN_CASCADE=$(jq -r '.cascade_mode // false' orchestrator-state-v3.json)

    # Check for previous integration attempts
    PREV_ATTEMPTS=$(ls -1 .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT-*.md 2>/dev/null | wc -l)

    # Check for recent BUILD-FIX-SUMMARY files
    RECENT_FIXES=$(find . -name "BUILD-FIX-SUMMARY-*.md" -mtime -1 | wc -l)

    if [[ "$IN_CASCADE" == "true" ]] || [[ $PREV_ATTEMPTS -gt 0 ]] || [[ $RECENT_FIXES -gt 0 ]]; then
        echo "✅ R521: Detected CASCADE context - known fixes likely available"
        return 0
    else
        echo "ℹ️ R521: First integration attempt - no known fixes expected"
        return 1
    fi
}
```

### Enhanced Search in Cascade Context
```bash
# In cascade context, search MORE aggressively for known fixes
cascade_fix_search() {
    local bug_pattern="$1"

    echo "🔍 R521: CASCADE MODE - Enhanced fix search"

    # Search in multiple locations
    for location in \
        ".software-factory/*/BUILD-FIX-SUMMARY-*.md" \
        "../*/BUILD-FIX-SUMMARY-*.md" \
        ".software-factory/phase*/wave*/BUILD-FIX-SUMMARY-*.md"; do

        if grep -r "$bug_pattern" $location 2>/dev/null; then
            echo "✅ Found known fix in: $location"
            return 0
        fi
    done

    echo "❌ No known fix found even in cascade context"
    return 1
}
```

## 🎬 DEMO CREATION = CONFLICT RESOLUTION (NOT NEW CODE)

### R291 Integration Infrastructure Exception
```bash
# Creating demos is EXPLICITLY ALLOWED per R291
# This is NOT a violation of R361 because:
create_integration_demos() {
    echo "🎬 R521: Creating demo infrastructure per R291"

    # This is conflict resolution because:
    # 1. R291 REQUIRES demos at integration level
    # 2. R504 PRE-PLANNED the demo locations
    # 3. R330 ensured demos were designed
    # 4. Integration agent is IMPLEMENTING pre-existing plan

    # NOT new code because:
    # - Demos were planned before integration (R504)
    # - Demo paths pre-defined in orchestrator-state-v3.json
    # - Integration agent executes plan, doesn't create plan

    # Create demos using PRE-PLANNED information
    DEMO_DIR=$(jq -r '.pre_planned_infrastructure.integrations.wave_integrations.phase1_wave1.demo_script_file' orchestrator-state-v3.json | xargs dirname)

    mkdir -p "$DEMO_DIR"
    # ... create demo artifacts per R291 ...

    echo "✅ Demo infrastructure created per R521 exception"
    echo "   This is integration infrastructure, not new code (R361 compliant)"
}
```

## ⚖️ R266 vs R521 RECONCILIATION

### R266 Says: "Document bugs, don't fix"
**APPLIES TO**: New bugs discovered during integration

### R521 Says: "Apply known fixes"
**APPLIES TO**: Bugs already documented and fixed in previous integrations

### How They Work Together:
1. Integration encounters issue
2. R521: Search for known fix
3. If found → R521: Apply it (conflict resolution)
4. If NOT found → R266: Document it (new bug)

**NO CONFLICT** - They cover different scenarios!

## 🚨 ENFORCEMENT AND VALIDATION

### Pre-Integration Check
```bash
# Before starting integration, verify BUILD-FIX-SUMMARY availability
pre_integration_check_r521() {
    echo "🔍 R521: Pre-integration validation"

    # Check if previous integrations exist
    if [ -f ".software-factory/wave*/BUILD-FIX-SUMMARY-*.md" ]; then
        echo "✅ Previous BUILD-FIX-SUMMARY files found"
        echo "   Will search for known fixes if issues arise"
    else
        echo "ℹ️ No previous BUILD-FIX-SUMMARY files"
        echo "   All issues will be NEW bugs per R266"
    fi
}
```

### Post-Integration Audit
```bash
# Verify integration agent followed R521 correctly
post_integration_audit_r521() {
    echo "📊 R521: Post-integration compliance audit"

    # Check for applied known fixes
    KNOWN_FIXES=$(grep -c "Known Fix Applied:" .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md || echo 0)

    # Check for new bugs documented
    NEW_BUGS=$(grep -c "NEW BUG DISCOVERED:" .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md || echo 0)

    echo "Results:"
    echo "  Known fixes applied: $KNOWN_FIXES"
    echo "  New bugs documented: $NEW_BUGS"

    # Verify known fixes were actually from BUILD-FIX-SUMMARY
    if [ $KNOWN_FIXES -gt 0 ]; then
        for fix in $(grep "Known Fix Applied:" .software-factory/INTEGRATE_WAVE_EFFORTS-REPORT.md | cut -d: -f2); do
            if ! grep -q "$fix" .software-factory/*/BUILD-FIX-SUMMARY-*.md; then
                echo "❌ R521 VIOLATION: Applied 'known fix' not found in BUILD-FIX-SUMMARY!"
                exit 521
            fi
        done
        echo "✅ All known fixes verified in BUILD-FIX-SUMMARY files"
    fi
}
```

## 📏 SIZE LIMITS STILL APPLY

**CRITICAL**: Even when applying known fixes:
- Maximum 50 lines total changes (R361)
- No new packages/directories (R361)
- No new files that don't exist in efforts (R361)

**If known fix requires >50 lines:**
```bash
if [ $KNOWN_FIX_LINES -gt 50 ]; then
    echo "❌ R361 VIOLATION: Known fix exceeds 50-line limit!"
    echo "   This fix must be applied in effort branches first"
    echo "   Cannot apply during integration"
    exit 361
fi
```

## Common Scenarios

### Scenario 1: BUG-007 Duplicate Found Again (CASCADE)
```bash
# ✅ CORRECT: Apply known fix
search_known_fixes "BUG-007: Duplicate PushCmd"  # Found in BUILD-FIX-SUMMARY
apply_known_fix "BUG-007"  # Delete duplicate file
# This is conflict resolution, not bug fixing
```

### Scenario 2: New Compilation Error (FIRST INTEGRATE_WAVE_EFFORTS)
```bash
# ✅ CORRECT: Document as new bug
search_known_fixes "undefined: registry.NewClient"  # Not found
document_new_bug "undefined: registry.NewClient" "..."  # Per R266
# Stop - need upstream fix
```

### Scenario 3: Missing Demo (ALWAYS)
```bash
# ✅ CORRECT: Create demo infrastructure
create_integration_demos()  # Per R291, not R361 violation
# This was pre-planned in R504
```

### Scenario 4: Adapter Needed (ALWAYS)
```bash
# ❌ WRONG: Create adapter
mkdir pkg/gitea  # R361 VIOLATION!

# ✅ CORRECT: Document need
document_new_bug "INTEGRATE_WAVE_EFFORTS_BLOCKED" "Needs pkg/gitea adapter"
# Stop - create adapter in effort branch
```

## Grading Impact

### COMPLIANCE BONUS (+30%)
- Correctly applying known fixes from BUILD-FIX-SUMMARY
- Proper cascade awareness
- Clean R266/R521 distinction

### MAJOR VIOLATIONS (-100%)
- Applying "known fix" that isn't documented
- Fixing new bugs under guise of "known fix"
- Creating new code claiming "conflict resolution"
- Bypassing R266 by misusing R521

## Related Rules
- R266: Upstream Bug Documentation (for NEW bugs)
- R361: Integration Conflict Resolution Only (size limits)
- R291: Integration Demo Requirement (demo exception)
- R300: Comprehensive Fix Management (cascade protocol)
- R504: Pre-Infrastructure Planning (demo pre-planning)
- R330: Demo Planning Requirements (demo design)

## Remember

**"Known fixes are conflict resolution, not development"**
**"Search BUILD-FIX-SUMMARY before documenting"**
**"Cascade context = known fixes expected"**
**"Demos are infrastructure, not new code"**

The integration agent is a SMART integrator - it applies PROVEN solutions from previous work, but it NEVER fixes NEW problems. That's the job of SW Engineers in effort branches.

## Date Added
2025-10-06

## Changelog
- 2025-10-06: Initial creation based on orchestrator analysis of integration agent false-negative R405 flags
