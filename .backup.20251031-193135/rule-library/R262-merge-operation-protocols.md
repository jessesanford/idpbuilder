# 🔴🔴🔴 RULE R262: Merge Operation Protocols 🔴🔴🔴

## Rule Definition
**Criticality:** SUPREME - NEVER VIOLATE
**Category:** Agent-Specific
**Applies To:** integration-agent

## ABSOLUTE REQUIREMENTS

### 1. NEVER MODIFY ORIGINAL BRANCHES
**THIS IS ABSOLUTE - NO EXCEPTIONS:**
- **NEVER** checkout and modify original feature branches
- **NEVER** force push to original branches
- **NEVER** rebase original branches
- **NEVER** amend commits on original branches
- Original branches must remain AS IS for audit trail

### 2. NO CHERRY PICKING
**PRESERVE COMPLETE HISTORY:**
- **NEVER** use cherry-pick between branches
- Preserve originals AS IS with full commit history
- Maintain commit-to-commit relationships
- Keep author information intact
- Preserve timestamps and metadata

### 3. SYNTHESIS BRANCH CREATION
When branches cannot merge cleanly:
- Create NEW synthesis branches
- Use naming: `{original-prefix}-synthesis-{timestamp}`
- Document synthesis rationale in commit message
- Link to original branches in documentation
- Preserve all original branches unchanged

### 4. CONFLICT RESOLUTION
The agent MUST:
- Resolve ALL conflicts in integration branch
- Document resolution decisions in work-log
- Test resolved code for compilation
- NEVER leave conflicts unresolved
- Create COMPLETELY MERGED WORKING COPY

## Enforcement Protocol

```bash
# CRITICAL: Verify no original branch modification
verify_branch_integrity() {
    local original_branch="$1"
    local before_sha="$2"
    local after_sha=$(git rev-parse "$original_branch")
    
    if [[ "$before_sha" != "$after_sha" ]]; then
        echo "🔴🔴🔴 FATAL: Original branch $original_branch was modified!"
        echo "Before: $before_sha"
        echo "After: $after_sha"
        exit 1
    fi
}

# Verify no cherry-picks
detect_cherry_picks() {
    local integration_branch="$1"
    
    if git log "$integration_branch" --grep="cherry picked from" > /dev/null; then
        echo "🔴🔴🔴 FATAL: Cherry-pick detected in $integration_branch!"
        exit 1
    fi
}

# Verify complete merge
verify_complete_merge() {
    local integration_branch="$1"
    
    # Check for conflict markers
    if grep -r "<<<<<<< HEAD" . --exclude-dir=.git; then
        echo "🔴🔴🔴 FATAL: Unresolved conflicts remain!"
        exit 1
    fi
    
    # Verify builds
    if ! make build 2>/dev/null; then
        echo "⚠️ WARNING: Build fails - document in INTEGRATE_WAVE_EFFORTS-REPORT.md"
    fi
}
```

## Example Correct Workflow

```bash
# CORRECT: Create integration branch, never modify originals
INTEGRATE_WAVE_EFFORTS_BRANCH="integration-$(date +%Y%m%d-%H%M%S)"
git checkout -b "$INTEGRATE_WAVE_EFFORTS_BRANCH" main

# Merge each branch in order (from INTEGRATE_WAVE_EFFORTS-PLAN.md)
for branch in branch1 branch2 branch3; do
    echo "Merging $branch..."
    git merge "$branch" --no-ff -m "integrate: $branch into $INTEGRATE_WAVE_EFFORTS_BRANCH"
    
    # Resolve any conflicts
    if [[ $? -ne 0 ]]; then
        echo "Resolving conflicts from $branch..."
        # Manual conflict resolution
        git add -A
        git commit -m "resolve: conflicts from merging $branch"
    fi
done
```

## Grading Impact
- 20% - AUTOMATIC FAILURE if original branches modified
- 20% - AUTOMATIC FAILURE if cherry-pick detected
- 10% - Complete conflict resolution

## Related Rules
- R260 - Integration Agent Core Requirements
- R263 - Integration Documentation Requirements
- R265 - Integration Testing Requirements