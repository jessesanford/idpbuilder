# 🚨🚨🚨 RULE R282: Phase Integration Protocol

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for phase completion
- **Penalty**: -50% to -100% for violations

## The Rule

**Phase integration MUST occur in completely isolated workspace with fresh clone of target repository, NOT the software-factory repository.**

## Requirements

### 1. Phase Integration Infrastructure (FOLLOWS R104)
```bash
# Setup for phase integration per R104
PHASE=$(yq '.current_phase' orchestrator-state.yaml)

# R104: Read target repository configuration
TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
TARGET_REPO_PATH=$(yq '.repository_path' "$TARGET_CONFIG")
TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")
DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")

# Create isolated phase integration workspace
PHASE_INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/efforts/phase${PHASE}/integration-workspace"
rm -rf "$PHASE_INTEGRATION_DIR"
mkdir -p "$PHASE_INTEGRATION_DIR"
cd "$PHASE_INTEGRATION_DIR"

# R104: Clone TARGET repository (NOT software-factory!)
git clone "$TARGET_REPO_PATH" "$TARGET_REPO_NAME"
cd "$TARGET_REPO_NAME"

# CRITICAL SAFETY CHECK - Verify correct repository
REMOTE_URL=$(git remote get-url origin)
if [[ "$REMOTE_URL" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Cloned orchestrator repository instead of target!"
    echo "Expected: Target project repository"
    echo "Got: $REMOTE_URL"
    exit 104  # R104 violation
fi

# R104: Create phase integration branch in TARGET repo
git checkout "$DEFAULT_BRANCH"
git pull origin "$DEFAULT_BRANCH"
PHASE_BRANCH="phase-${PHASE}-integration"  # R104 naming convention
git checkout -b "$PHASE_BRANCH"
```

### 2. Wave Integration Sequence
All waves in the phase must be integrated sequentially:
```bash
# For each wave in phase
for wave in $(seq 1 $TOTAL_WAVES); do
    echo "📦 Integrating Wave $wave..."
    
    # Get wave integration branch (R104 naming)
    WAVE_BRANCH="wave-${wave}-integration"
    
    # Fetch and merge wave
    git fetch origin "$WAVE_BRANCH"
    git merge "origin/$WAVE_BRANCH" --no-ff \
        -m "integrate: Wave $wave into Phase $PHASE"
    
    # Test after EACH wave merge
    if ! npm test || ! make test || ! pytest; then
        echo "🚨 Tests failed after integrating Wave $wave"
        exit 1
    fi
done
```

### 3. Workspace Isolation Requirements (Per R104)
- **Location**: `$CLAUDE_PROJECT_DIR/efforts/phase{X}/integration-workspace/[target-repo-name]/`
- **Repository**: Fresh clone of TARGET repository (R104 requirement)
- **Branch**: `phase-{X}-integration` (R104 naming convention)
- **NEVER**: Use software-factory repository
- **NEVER**: Reuse wave integration workspaces

### 4. Validation Checks
```bash
# Continuous validation during phase integration
validate_phase_integration() {
    # Check correct directory
    if [[ "$PWD" != */phase-integration* ]]; then
        echo "❌ Not in phase integration workspace!"
        return 1
    fi
    
    # Check NOT software-factory
    if git remote get-url origin | grep -q "software-factory"; then
        echo "❌ CRITICAL: Working in wrong repository!"
        return 1
    fi
    
    # Check branch naming
    if ! git branch --show-current | grep -q "phase.*integration"; then
        echo "❌ Wrong branch naming convention!"
        return 1
    fi
    
    return 0
}
```

## Integration Process

### Step 1: Pre-Integration Validation
```bash
# Verify all waves completed
for wave in $(seq 1 $TOTAL_WAVES); do
    wave_status=$(yq ".waves.wave_${wave}.status" orchestrator-state.yaml)
    if [ "$wave_status" != "INTEGRATED" ]; then
        echo "🚨 Cannot integrate phase - Wave $wave not integrated"
        exit 1
    fi
done
```

### Step 2: Sequential Wave Merging
- Merge each wave branch in order
- Run full test suite after each merge
- Document any conflicts
- Stop on first failure

### Step 3: Phase Validation & Demo (R291)

#### 🏗️ Build Verification
```bash
# MANDATORY: Verify build after phase integration
echo "🏗️ Running phase build verification..."
npm run build:prod 2>&1 | tee phase${PHASE}-build.log
if [ ! -d "dist" ] && [ ! -d "build" ]; then
    echo "❌ Build failed - phase integration incomplete!"
    exit 1
fi
echo "✅ Phase build successful"
```

#### 🧪 Test Harness
```bash
# MANDATORY: Create phase test harness
cat > phase${PHASE}-test-harness.sh << 'EOF'
#!/bin/bash
echo "🧪 Phase Integration Test Suite"
echo "================================"

# Full test suite
npm test || exit 1
npm run test:integration || exit 1
npm run test:e2e || exit 1

# Phase-specific tests
if [ -f "tests/phase${PHASE}/run-all.sh" ]; then
    bash tests/phase${PHASE}/run-all.sh || exit 1
fi

echo "✅ All phase tests passed!"
EOF

chmod +x phase${PHASE}-test-harness.sh
./phase${PHASE}-test-harness.sh
```

#### 🎬 Demo Creation
```bash
# MANDATORY: Create phase demo
cat > PHASE-${PHASE}-DEMO.md << EOF
# Phase $PHASE Integration Demo

## Build & Test Status
- Build: ✅ PASSING (see phase${PHASE}-build.log)
- Tests: ✅ ALL PASSING
- Test Harness: phase${PHASE}-test-harness.sh

## Waves Integrated
$(for w in $(seq 1 $TOTAL_WAVES); do echo "- Wave $w: INTEGRATED"; done)

## Features Delivered
[List major features from phase]

## How to Run Demo
\`\`\`bash
./phase${PHASE}-demo.sh
\`\`\`

## Metrics
- Total Lines: $(git diff main --numstat | awk '{sum+=$1+$2} END {print sum}')
- Test Coverage: $(npm run coverage --silent | grep "All files" | awk '{print $10}')
- Build Size: $(du -sh dist/ 2>/dev/null | awk '{print $1}')
EOF

# Create demo script
cat > phase${PHASE}-demo.sh << 'EOF'
#!/bin/bash
echo "🎬 Phase Integration Demo"
# Add actual demo commands
npm start &
sleep 5
curl http://localhost:3000/api/phase-features
EOF
chmod +x phase${PHASE}-demo.sh
```

### Step 4: Phase Report Generation
```bash
# Generate comprehensive phase report
cat > PHASE-${PHASE}-INTEGRATION.md << EOF
# Phase $PHASE Integration Report

## Waves Integrated
$(for w in $(seq 1 $TOTAL_WAVES); do echo "- Wave $w: INTEGRATED"; done)

## Build & Test Results
- Build: ✅ SUCCESSFUL
- Unit Tests: ✅ PASSING
- Integration Tests: ✅ PASSING
- E2E Tests: ✅ PASSING
- Test Harness: phase${PHASE}-test-harness.sh

## Demo Artifacts
- Demo Documentation: PHASE-${PHASE}-DEMO.md
- Demo Script: phase${PHASE}-demo.sh
- Build Log: phase${PHASE}-build.log

## Metrics
- Total Lines: $total_lines
- Tests: ALL PASSING
- Conflicts: None/Resolved

## Validation
- Repository: Verified (not software-factory)
- Workspace: Isolated at $PWD
- Branch: $PHASE_BRANCH
EOF
```

### Step 4: Push and Create PR
```bash
# Push phase integration
git push -u origin "$PHASE_BRANCH"

# Create phase PR
gh pr create \
  --title "Phase ${PHASE} Integration - ${total_lines} lines" \
  --body "$(cat PHASE-${PHASE}-INTEGRATION.md)"
```

## Failure Conditions

### Critical Failures (Immediate Stop)
- 🚨 Cloning software-factory repository = FAIL
- 🚨 Using non-isolated workspace = FAIL
- 🚨 Missing wave integrations = FAIL
- 🚨 Test failures after merge = FAIL

### Recovery Protocol
1. Document failure point
2. Transition to ERROR_RECOVERY
3. Clean workspace completely
4. Start fresh with new clone
5. Never continue from corrupted state

## Success Criteria
- ✅ All waves integrated sequentially
- ✅ All tests passing
- ✅ Correct repository verified
- ✅ Isolated workspace used
- ✅ Phase branch pushed
- ✅ PR created
- ✅ State file updated
- ✅ Build verified working (R291)
- ✅ Test harness executed (R291)
- ✅ Demo created and functional (R291)
- ✅ All artifacts documented

## Penalties
- Wrong repository: **-100% grade** (CRITICAL)
- Non-isolated workspace: **-75% grade**
- Missing waves: **-50% grade**
- Test failures ignored: **-60% grade**
- No PR created: **-30% grade**

## Related Rules
- R034: Integration Requirements
- R250: Integration Isolation
- R259: Phase Integration After Fixes
- R256: Phase Assessment Gate
- R288/R288: State File Updates
- R283: Project Integration Protocol

## Common Violations

### ❌ WRONG: Using software-factory repo
```bash
cd /home/vscode/software-factory-template
git checkout -b phase-integration  # CRITICAL ERROR!
```

### ❌ WRONG: Reusing wave workspace
```bash
cd /efforts/phase1/wave1/integration-workspace
git checkout -b phase-integration  # CONTAMINATION!
```

### ✅ CORRECT: Fresh isolated workspace
```bash
rm -rf /efforts/phase1/phase-integration
git clone [target-repo] /efforts/phase1/phase-integration
cd /efforts/phase1/phase-integration
# Verify NOT software-factory
git checkout -b project-phase1-integration
```

## Remember
**Phase integration is a critical consolidation point.** Using the wrong repository or contaminated workspace will corrupt the entire project. Always verify repository before proceeding.