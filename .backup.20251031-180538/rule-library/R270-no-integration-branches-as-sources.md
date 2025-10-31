# 🔴🔴🔴 RULE R270: No Integration Branches as Sources 🔴🔴🔴

## Rule Definition
**Criticality:** SUPREME - Violation = Recursive integration chaos
**Category:** Integration
**Applies To:** code-reviewer, integration-agent, orchestrator
**Created:** 2025-08-27
**Modified:** 2025-01-20 (Clarified integration branches NEVER merge to main)

## SUPREME LAW

**INTEGRATE_WAVE_EFFORTS BRANCHES ARE TARGETS, NEVER SOURCES!**
**INTEGRATE_WAVE_EFFORTS BRANCHES NEVER MERGE TO MAIN!** (See R363/R364)

When creating merge plans or executing integrations:
- ✅ Merge FROM original effort branches TO main
- ✅ Merge FROM split branches TO main
- ✅ Merge FROM ERROR_RECOVERY fix branches TO main
- ❌ NEVER merge FROM wave integration branches
- ❌ NEVER merge FROM phase integration branches
- ❌ NEVER merge FROM any "*-integration-*" branch
- ❌ NEVER merge integration branches TO main

## Why This Is SUPREME

Using integration branches as sources creates:
1. **Recursive Integration Hell** - Integration of integrations of integrations
2. **Lost Commit Attribution** - Original authors obscured by integration commits
3. **Conflict Multiplication** - Same conflicts resolved multiple times
4. **History Corruption** - Non-linear story becomes incomprehensible
5. **Untraceable Changes** - Cannot determine which effort introduced what

## Examples

### ✅ CORRECT: Original Branches Only
```yaml
Wave 2 Integration Sources:
  - phase3/wave2/effort1-api-gateway      # Original effort
  - phase3/wave2/effort2-controller-split1 # Split of original
  - phase3/wave2/effort2-controller-split2 # Split of original
  - phase3/wave2/effort3-webhooks         # Original effort
```

### ❌ WRONG: Integration Branches as Sources
```yaml
Phase 3 Integration Sources:
  - phase3-wave1-integration-20250826  # ❌ NO! Integration branch!
  - phase3-wave2-integration-20250827  # ❌ NO! Integration branch!
  - phase3-wave3-integration-20250828  # ❌ NO! Integration branch!
```

### ✅ CORRECT: Phase Integration Sources
```yaml
Phase 3 Integration Sources:
  # Wave 1 originals
  - phase3/wave1/effort1-api-types
  - phase3/wave1/effort2-controller-split1
  - phase3/wave1/effort2-controller-split2
  # Wave 2 originals  
  - phase3/wave2/effort1-api-gateway
  - phase3/wave2/effort2-validation
  # Wave 3 originals
  - phase3/wave3/effort1-reconciler
  # ERROR_RECOVERY fixes
  - phase3-fix-kcp-patterns-20250827
  - phase3-fix-api-compatibility-20250827
```

## Enforcement

```bash
# Validate merge plan has no integration branches as sources
validate_no_integration_sources() {
    local merge_plan="$1"
    
    # Check for integration branches in merge commands
    if grep "git merge.*integration" "$merge_plan" | grep -v "^#"; then
        echo "🔴🔴🔴 SUPREME VIOLATION: Integration branch used as source!"
        grep "git merge.*integration" "$merge_plan" | grep -v "^#"
        return 1
    fi
    
    # Check for wave integration patterns
    if grep -E "wave[0-9]-integration" "$merge_plan" | grep -i "merge"; then
        echo "🔴🔴🔴 SUPREME VIOLATION: Wave integration branch as source!"
        return 1
    fi
    
    # Check for phase integration patterns
    if grep -E "phase[0-9]-integration" "$merge_plan" | grep -i "merge"; then
        echo "🔴🔴🔴 SUPREME VIOLATION: Phase integration branch as source!"
        return 1
    fi
    
    echo "✅ No integration branches used as sources"
    return 0
}
```

## The Correct Flow (Updated per R363/R364)

### Wave Integration (TESTING ONLY)
```
main
 ├─→ effort1 branch ─┐
 ├─→ effort2 branch ─┼─→ wave-integration branch (FOR TESTING ONLY)
 └─→ effort3 branch ─┘     ❌ NEVER merges to main
                           ✅ Gets deleted after testing
```

### Sequential Merging to Main (THE CORRECT WAY)
```
main
 ←── effort1 branch (merges directly to main)
 ←── effort2 branch (merges directly to main AFTER effort1)
 ←── effort3 branch (merges directly to main AFTER effort2)

Integration branches test but NEVER merge to main!
```

### Phase Integration (TESTING ONLY)
```
main
 ├─→ wave1/effort1 ─┐
 ├─→ wave1/effort2 ─┤
 ├─→ wave2/effort1 ─┼─→ phase-integration branch (FOR TESTING ONLY)
 ├─→ wave2/effort2 ─┤     ❌ NEVER merges to main
 └─→ fix branches  ─┘     ✅ Gets deleted after testing
```

**ABSOLUTELY NOT THIS:**
```
main
 ←── wave1-integration  ❌ NEVER! Integration branches don't merge to main
 ←── wave2-integration  ❌ NEVER! Violates R363/R364
 ←── phase-integration  ❌ NEVER! Testing branches stay separate
```

## Common Violations

### 1. Lazy Phase Integration
```bash
# ❌ WRONG - Using wave integration branches
for wave in 1 2 3 4; do
    git merge origin/phase3-wave${wave}-integration
done
```

### 2. Shortcut Thinking
"Wave integration branches already have everything merged, so I'll just merge those"
**NO!** This creates recursive integration and loses the clean history.

### 3. Misunderstanding Integration Purpose
Integration branches are RESULTS, not INPUTS:
- Wave integration = Result of merging wave efforts
- Phase integration = Result of merging ALL original efforts + fixes
- They are endpoints, not waypoints

## Grading Impact
- **AUTOMATIC FAILURE** for using integration branches as sources
- -50% for each integration branch used as source
- -100% for systematic use of integration branches

## Exception Cases
**NONE** - This is a SUPREME rule with no exceptions.

## Related Rules
- R363 - Sequential Direct Mergeability (FUNDAMENTAL)
- R364 - Integration Testing Only Branches (FUNDAMENTAL)
- R250 - Integration Isolation
- R261 - Integration Planning Requirements
- R262 - Merge Operation Protocols
- R269 - Code Reviewer Merge Plan No Execution