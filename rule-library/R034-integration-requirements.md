# 🚨🚨🚨 RULE R034: Integration Requirements

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for wave completion
- **Penalty**: -20% to -50% for violations

## The Rule

**Integration of efforts into wave branches MUST follow strict protocol to ensure code quality, test compliance, proper version control, AND demonstrate working functionality per R291.**

### 🔴🔴🔴 PARAMOUNT REQUIREMENT (R307) 🔴🔴🔴
**EVERY integrated effort MUST maintain independent mergeability:**
- Each effort must compile and work when merged alone
- No effort can break existing functionality
- Incomplete features must be behind feature flags
- The build must NEVER break during integration

## Pre-Integration Requirements

### 1. Mandatory Checklist
Before ANY integration attempt, verify:
- ✅ All efforts in wave marked as COMPLETED
- ✅ All code reviews PASSED (R222)
- ✅ All tests passing in individual efforts
- ✅ Line count verified under limits (<800 lines per effort)
- ✅ Integration workspace isolated per R250
- ✅ All effort code committed and pushed
- ✅ Build verification prepared (R291)
- ✅ Test harness ready (R291)
- ✅ Demo plan prepared (R291)
- ✅ No deprecated branches in integration list (R296)
- ✅ All SPLIT_DEPRECATED efforts using replacement splits

### 2. State Validation
```bash
# Verify all efforts completed
for effort in $(jq '.efforts_in_progress[]' orchestrator-state.json); do
    status=$(jq ".efforts_completed[] | select(.name == \"$effort\") | .status" orchestrator-state.json)
    if [ "$status" != "COMPLETED" ]; then
        echo "🚨 Cannot integrate - $effort not completed"
        exit 1
    fi
done

# R296: Check for deprecated branches
for effort in $(jq '.efforts_completed[].name' orchestrator-state.json); do
    status=$(jq ".efforts_completed[] | select(.name == \"$effort\") | .status" orchestrator-state.json)
    if [ "$status" == "SPLIT_DEPRECATED" ]; then
        echo "🚨 BLOCKED: Cannot integrate deprecated effort: $effort"
        SPLITS=$(jq ".efforts_completed[] | select(.name == \"$effort\") | .replacement_splits[]" orchestrator-state.json)
        echo "Use replacement splits instead:"
        echo "$SPLITS"
        exit 1
    fi
done
```

## Integration Protocol

### Step 1: Setup Isolated Workspace (R250 + R104)
```bash
PHASE=$(jq '.current_phase' orchestrator-state.json)
WAVE=$(jq '.current_wave' orchestrator-state.json)

# Per R104: Read target repository configuration
TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
TARGET_REPO_PATH=$(yq '.repository_path' "$TARGET_CONFIG")
TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")
DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")

# Per R250: Create isolated integration workspace
INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/wave${WAVE}/integration-workspace"
mkdir -p "$INTEGRATION_DIR"
cd "$INTEGRATION_DIR"

# Per R104: Clone TARGET repository (NOT software-factory!)
git clone "$TARGET_REPO_PATH" "$TARGET_REPO_NAME"
cd "$TARGET_REPO_NAME"

# CRITICAL SAFETY CHECK - Verify correct repository
REMOTE_URL=$(git remote get-url origin)
if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
    echo "Expected: Target project repository"
    echo "Got: $REMOTE_URL"
    exit 1
fi
```

### Step 2: Create Integration Branch (Per R104)
```bash
# Per R104: Create integration branch in TARGET repository
git checkout "$DEFAULT_BRANCH"
git pull origin "$DEFAULT_BRANCH"
git checkout -b "wave-${WAVE}-integration"  # R104 naming convention
```

### Step 3: Cherry-Pick or Merge Each Effort
```bash
# For each completed effort
for effort_dir in ../effort*/; do
    effort_name=$(basename "$effort_dir")
    echo "📦 Integrating $effort_name..."
    
    # R296: Skip deprecated branches
    if [[ "$effort_name" == *"-deprecated-split" ]]; then
        echo "⚠️ Skipping deprecated branch: $effort_name"
        continue
    fi
    
    # Add effort as remote
    git remote add "$effort_name" "$effort_dir"
    git fetch "$effort_name"
    
    # Get effort branch name
    effort_branch=$(cd "$effort_dir" && git branch --show-current)
    
    # R296: Final check - ensure branch isn't deprecated
    if [[ "$effort_branch" == *"-deprecated-split" ]]; then
        echo "❌ BLOCKED: Cannot integrate deprecated branch: $effort_branch"
        echo "Check state file for replacement splits"
        exit 1
    fi
    
    # Cherry-pick commits (preferred) or merge
    git cherry-pick $(git rev-list --reverse main.."$effort_name/$effort_branch")
    # OR: git merge "$effort_name/$effort_branch" --no-ff
    
    # Test after EACH integration
    if ! npm test || ! make test || ! pytest; then
        echo "🚨 Tests failed after integrating $effort_name"
        exit 1
    fi
done
```

### Step 4: Final Validation
```bash
# Run complete test suite
npm test || make test || pytest

# Verify line count of integration
total_lines=$(git diff main --numstat | awk '{sum+=$1+$2} END {print sum}')
if [ "$total_lines" -gt 800 ]; then
    echo "🚨 Integration exceeds 800 lines: $total_lines"
    exit 1
fi
```

### Step 5: Push Integration Branch
```bash
git push -u origin "$INTEGRATION_BRANCH"
```

## Commit Message Standards

### For Cherry-Picks
```
cherry-pick: effort-{name} - {description}

Source: {effort_branch}
Changes: {X} lines
Tests: PASS
Review: {reviewer_agent}
```

### For Merge Commits
```
integrate: merge {effort_name} into wave{N}

Effort: {effort_description}
Changes: {X} lines added, {Y} lines removed
Tests: ALL PASSING
Review: Approved by {reviewer}
Conflicts: None|Resolved
```

## Post-Integration Requirements

1. **Create Pull Request**
```bash
gh pr create \
  --title "Wave ${WAVE} Integration - ${total_lines} lines" \
  --body "Integration of all Wave ${WAVE} efforts\n\nEfforts included:\n- effort1\n- effort2\n\nAll tests passing."
```

2. **Update State File**
```yaml
integration:
  wave_{N}:
    branch: wave-1-integration  # Per R104 naming
    repository: target-repo-name  # TARGET repository, not software-factory
    efforts_merged: [effort1, effort2]
    total_lines: 750
    tests: PASS
    pr_number: 123
    completed_at: 2024-01-29T10:00:00Z
```

3. **Archive Integration Artifacts**
- Save test results
- Document any conflict resolutions
- Record integration metrics

## Failure Recovery Protocol

If integration fails:
1. **STOP immediately** - Do not force or skip
2. **Document failure**:
   ```yaml
   integration_failure:
     wave: 1
     failed_at: effort2_merge
     reason: "Test failures after merge"
     error: "Unit test X failed"
   ```
3. **Transition to ERROR_RECOVERY**
4. **Create recovery plan** before retry
5. **NEVER force merge or skip tests**

## 🎬 Mandatory Demo Requirements (R291)

### Build Verification
```bash
# REQUIRED: Verify build after integration
echo "🏗️ Running build verification..."
npm run build 2>&1 | tee wave-build.log
[ -d "dist" ] || [ -d "build" ] || { echo "❌ Build failed!"; exit 1; }
```

### Test Harness
```bash
# REQUIRED: Create and run test harness
cat > wave-test-harness.sh << 'EOF'
#!/bin/bash
echo "🧪 Wave Integration Test Suite"
npm test || exit 1
npm run test:integration || exit 1
./verify-wave-features.sh || exit 1
echo "✅ All tests passed!"
EOF
chmod +x wave-test-harness.sh
./wave-test-harness.sh
```

### Demo Documentation
```bash
# REQUIRED: Create demo artifacts
cat > WAVE-DEMO.md << EOF
# Wave ${WAVE} Demo
- Build: ✅ PASSING
- Tests: ✅ ALL PASSING
- Features: [List integrated features]
- How to run: ./demo-wave.sh
EOF
```

## Success Criteria
- ✅ All efforts integrated without unresolved conflicts
- ✅ All tests passing on integration branch
- ✅ Total line count under limits
- ✅ Integration branch pushed to remote
- ✅ PR created for integration
- ✅ State file updated with results
- ✅ Build verified working (R291)
- ✅ Test harness executed (R291)
- ✅ Demo created and verified (R291)

## Penalties
- Skipping tests: **-40% grade**
- Force merging conflicts: **-50% grade**
- Missing efforts: **-30% grade per effort**
- No PR created: **-20% grade**
- State not updated: **-25% grade**
- Exceeding line limits: **-35% grade**

## Related Rules
- R250: Integration Isolation Requirement
- R222: Code Review Gate
- R271: Full Checkout Requirement
- R014: Branch Naming Convention
- R288: State File Updates
- R296: Deprecated Branch Marking Protocol

## Common Failures to Avoid

### ❌ Force Merging
```bash
git merge --force  # NEVER DO THIS
```

### ❌ Skipping Tests
```bash
# Merging without testing
git merge effort1
git merge effort2  # NO TEST = FAILURE
```

### ✅ Correct Approach
```bash
git merge effort1 --no-ff
npm test  # TEST AFTER EACH
git merge effort2 --no-ff
npm test  # TEST AGAIN
```

## Remember
Integration is where parallel efforts converge. **Every merge must be validated** before proceeding. Failed integration is better than corrupted main branch.