# 🔴🔴🔴 RULE R308 - Incremental Branching Strategy [CORE TENANT]

## Classification
- **Category**: Core Development Flow
- **Criticality Level**: 🔴🔴🔴 SUPREME LAW
- **Enforcement**: MANDATORY - Trunk-based development foundation
- **Penalty**: -100% for violating incremental flow

## The Rule

**EVERY effort MUST branch from the LATEST INTEGRATED CODE of the previous wave/phase. This is TRUE TRUNK-based development where each wave builds incrementally on the previous wave's integrated work.**

### 🔴 THE INCREMENTAL PRINCIPLE 🔴

```
Phase 1, Wave 1: Efforts branch from → main
Phase 1, Wave 2: Efforts branch from → phase1-wave1-integration
Phase 1, Wave 3: Efforts branch from → phase1-wave2-integration
Phase 2, Wave 1: Efforts branch from → phase1-integration  
Phase 2, Wave 2: Efforts branch from → phase2-wave1-integration
Phase 3, Wave 1: Efforts branch from → phase2-integration
```

**NO EFFORT MAY BRANCH FROM STALE CODE!**

## Why This Is CRITICAL

### Traditional (WRONG) Approach:
```
main
  ├─→ ALL Phase 1 efforts (stale base!)
  ├─→ ALL Phase 2 efforts (conflicts!)
  └─→ ALL Phase 3 efforts (massive conflicts!)
```
**Result**: Integration nightmare, "big bang" merges, conflicts discovered late

### Incremental (CORRECT) Approach:
```
main
  ├─→ P1W1 efforts → integrate → P1W1-integration
                                    ├─→ P1W2 efforts → integrate → P1W2-integration
                                                                     ├─→ P1W3 efforts
```
**Result**: Conflicts detected early, smooth integration, true CI/CD

## Mandatory Implementation

### 1. Orchestrator MUST Determine Base Branch

```bash
determine_effort_base_branch() {
    local PHASE=$1
    local WAVE=$2
    
    # Phase 1, Wave 1: Start from main
    if [[ $PHASE -eq 1 && $WAVE -eq 1 ]]; then
        echo "main"
        return
    fi
    
    # First wave of new phase: From previous phase integration
    if [[ $WAVE -eq 1 ]]; then
        PREV_PHASE=$((PHASE - 1))
        BASE="phase${PREV_PHASE}-integration"
        
        # Verify it exists
        if git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
            echo "$BASE"
        else
            echo "❌ FATAL: Previous phase integration not found: $BASE" >&2
            exit 1
        fi
        return
    fi
    
    # Subsequent waves: From previous wave integration
    PREV_WAVE=$((WAVE - 1))
    BASE="phase${PHASE}-wave${PREV_WAVE}-integration"
    
    # Verify it exists
    if git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
        echo "$BASE"
    else
        echo "❌ FATAL: Previous wave integration not found: $BASE" >&2
        exit 1
    fi
}
```

### 2. Clone MUST Use Correct Base

```bash
# WRONG - Always using main
git clone --branch main "$REPO" "$EFFORT_DIR"

# CORRECT - Using incremental base
BASE=$(determine_effort_base_branch $PHASE $WAVE)
git clone --branch "$BASE" "$REPO" "$EFFORT_DIR"
```

### 3. Agents MUST Verify Their Base

```bash
verify_incremental_base() {
    local expected_base="$1"
    local actual_base=$(git merge-base HEAD origin/main)
    local base_branch=$(git branch -r --contains "$actual_base" | grep -E "(main|integration)" | head -1)
    
    if [[ ! "$base_branch" =~ "$expected_base" ]]; then
        echo "❌ FATAL: Not based on incremental integration!"
        echo "Expected base: $expected_base"
        echo "Actual base: $base_branch"
        exit 1
    fi
}
```

## Integration Flows

### Wave Integration Creates Next Base
```
P1W1 efforts complete
    ↓
Create phase1-wave1-integration (merge all P1W1)
    ↓
P1W2 efforts START from phase1-wave1-integration
```

### Phase Integration Creates Phase Base
```
P1W3 efforts complete
    ↓
Create phase1-wave3-integration
    ↓
Create phase1-integration (merge all P1 waves)
    ↓
P2W1 efforts START from phase1-integration
```

## Split Branching (Sequential, Not Incremental)

**IMPORTANT**: Splits follow SEQUENTIAL branching, not incremental:
- Original effort: From wave/phase integration
- split-001: From SAME base as original
- split-002: From split-001
- split-003: From split-002

This maintains split atomicity while building incrementally.

## Verification Protocol

### Pre-Clone Verification
```bash
echo "🔍 Verifying incremental base branch..."
PHASE=$(yq '.current_phase' orchestrator-state.yaml)
WAVE=$(yq '.current_wave' orchestrator-state.yaml)
BASE=$(determine_effort_base_branch $PHASE $WAVE)

echo "📌 Phase $PHASE, Wave $WAVE"
echo "📌 Incremental base: $BASE"

if ! git ls-remote --heads origin "$BASE" > /dev/null 2>&1; then
    echo "❌ FATAL: Integration branch not ready!"
    echo "Previous wave must be integrated first"
    exit 1
fi
```

### Post-Clone Verification
```bash
cd "$EFFORT_DIR"
COMMITS_SINCE_MAIN=$(git rev-list --count origin/main..HEAD)
echo "📊 Commits ahead of main: $COMMITS_SINCE_MAIN"

# Should include previous wave's work
if [[ $WAVE -gt 1 && $COMMITS_SINCE_MAIN -eq 0 ]]; then
    echo "⚠️ WARNING: No previous wave commits found!"
    echo "May not be properly based on integration"
fi
```

## Examples

### ✅ CORRECT: Phase 2, Wave 2 Effort
```bash
# Determine base (will return phase2-wave1-integration)
BASE=$(determine_effort_base_branch 2 2)

# Clone from integration
git clone --branch "$BASE" "$REPO" "efforts/phase2/wave2/new-feature"

# Create effort branch
cd "efforts/phase2/wave2/new-feature"
git checkout -b "phase2/wave2/new-feature"

# Verify includes Wave 1 work
git log --oneline | grep "wave1" || echo "⚠️ Missing wave1 commits!"
```

### ❌ WRONG: Always Using Main
```bash
# VIOLATION - Not incremental!
git clone --branch main "$REPO" "efforts/phase2/wave2/new-feature"
# Missing all Phase 1 and Phase 2 Wave 1 work!
```

### ❌ WRONG: Using Target Config Base Only
```bash
# VIOLATION - Ignoring integration!
BASE=$(yq '.base_branch' target-repo-config.yaml)  # Always "main"
git clone --branch "$BASE" "$REPO" "$EFFORT_DIR"
# Not building on previous work!
```

## Agent-Specific Requirements

### Orchestrator
- MUST determine correct incremental base
- MUST verify integration branch exists
- MUST document base branch decision
- MUST refuse to proceed if integration missing

### SW Engineer  
- MUST verify working from correct base
- MUST confirm includes previous wave work
- MUST refuse work if not incremental

### Code Reviewer
- MUST verify effort has correct base
- MUST check for integration conflicts
- MUST validate incremental development

### Integration Agent
- MUST create integration branches
- MUST push immediately for next wave
- MUST verify branch accessible

## Visual Flow Diagram

```
main (Phase 1 start)
    │
    ├─→ P1W1-E1 ─┐
    ├─→ P1W1-E2 ─┼─→ phase1-wave1-integration
    └─→ P1W1-E3 ─┘            │
                              ├─→ P1W2-E1 ─┐
                              ├─→ P1W2-E2 ─┼─→ phase1-wave2-integration
                              └─→ P1W2-E3 ─┘            │
                                                        └─→ phase1-integration
                                                                    │
                                                                    ├─→ P2W1-E1
                                                                    └─→ P2W1-E2
```

## Grading Impact

- **Using stale base (main when should use integration)**: -100%
- **Skipping incremental flow**: -100%  
- **Not verifying base branch**: -50%
- **Proceeding without integration branch**: -75%
- **Not documenting base decision**: -25%

## Integration with Other Rules

- **R196**: Base branch selection (enhanced by R308)
- **R271**: Single-branch full checkout (uses R308 base)
- **R034**: Integration requirements (creates R308 bases)
- **R302**: Split tracking (sequential within incremental)
- **R304**: Line counting (counts from R308 base)

## Stop Work Conditions

**IMMEDIATE STOP if:**
1. Previous wave integration branch doesn't exist
2. Effort not based on correct incremental branch
3. Integration conflicts detected with base
4. Base branch has pending unmerged work

## Rule Authority

This rule is a CORE TENANT of trunk-based development and SUPERSEDES any configuration that would cause efforts to branch from stale code. The incremental principle is NON-NEGOTIABLE.

---
**Remember**: Every effort builds on what came before. This is how we achieve TRUE continuous integration!