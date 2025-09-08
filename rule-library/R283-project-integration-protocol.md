# 🚨🚨🚨 RULE R283: Project Integration Protocol

## Classification
- **Category**: Integration Management
- **Criticality Level**: 🚨🚨🚨 BLOCKING
- **Enforcement**: MANDATORY for project completion
- **Penalty**: -75% to -100% for violations

## The Rule

**Final project integration MUST occur in completely isolated workspace at `/efforts/integration/` with fresh clone of target repository, NEVER the software-factory repository.**

## Requirements

### 1. Project Integration Infrastructure (FOLLOWS R104)
```bash
# Setup for final project integration per R104
TOTAL_PHASES=$(yq '.total_phases' orchestrator-state.json)

# R104: Read target repository configuration
TARGET_CONFIG="$CLAUDE_PROJECT_DIR/target-repo-config.yaml"
TARGET_REPO_PATH=$(yq '.repository_path' "$TARGET_CONFIG")
TARGET_REPO_NAME=$(yq '.repository_name' "$TARGET_CONFIG")
DEFAULT_BRANCH=$(yq '.default_branch' "$TARGET_CONFIG")

# Create isolated project integration workspace
PROJECT_INTEGRATION_DIR="$CLAUDE_PROJECT_DIR/efforts/project/integration-workspace"
rm -rf "$PROJECT_INTEGRATION_DIR"
mkdir -p "$PROJECT_INTEGRATION_DIR"
cd "$PROJECT_INTEGRATION_DIR"

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

# Verify we're NOT in software-factory directory
if [[ "$PWD" == *"software-factory"* ]]; then
    echo "❌ CRITICAL: Integration happening in software-factory directory!"
    exit 104  # R104 violation
fi

# R104: Create final integration branch in TARGET repo
git checkout "$DEFAULT_BRANCH"
git pull origin "$DEFAULT_BRANCH"
INTEGRATION_BRANCH="project-integration"  # R104 naming convention
git checkout -b "$INTEGRATION_BRANCH"
```

### 2. Phase Integration Sequence
All phases must be integrated sequentially:
```bash
# For each phase in project
for phase in $(seq 1 $TOTAL_PHASES); do
    echo "📦 Integrating Phase $phase..."
    
    # Get phase integration branch (R104 naming)
    PHASE_BRANCH="phase-${phase}-integration"
    
    # Fetch and merge phase
    git fetch origin "$PHASE_BRANCH"
    if ! git merge "origin/$PHASE_BRANCH" --no-ff \
         -m "integrate: Phase $phase into final project"; then
        echo "🚨 Merge conflict in Phase $phase"
        # Document conflicts and resolve
        git status > "PHASE-${phase}-CONFLICTS.txt"
        exit 1
    fi
    
    # Test after EACH phase merge
    echo "🧪 Running tests after Phase $phase integration..."
    if ! npm test || ! make test || ! pytest; then
        echo "🚨 Tests failed after integrating Phase $phase"
        exit 1
    fi
    
    # Verify build still works
    if ! npm run build || ! make build; then
        echo "🚨 Build failed after integrating Phase $phase"
        exit 1
    fi
done
```

### 3. Workspace Isolation Requirements (Per R104)
- **Location**: `$CLAUDE_PROJECT_DIR/efforts/project/integration-workspace/[target-repo-name]/`
- **Repository**: Fresh clone of TARGET repository (R104 requirement)
- **Branch**: `project-integration` (R104 naming convention)
- **NEVER**: Use software-factory repository
- **NEVER**: Reuse phase integration workspaces
- **NEVER**: Work in software-factory-template directory

### 4. Continuous Validation
```bash
# Validation function for project integration
validate_project_integration() {
    # Check correct directory (must be in project integration workspace)
    if [[ "$PWD" != */efforts/project/integration-workspace/* ]]; then
        echo "❌ Not in project integration workspace!"
        echo "Expected: /efforts/integration/*"
        echo "Got: $PWD"
        return 1
    fi
    
    # Check NOT software-factory
    REMOTE=$(git remote get-url origin 2>/dev/null)
    if [[ "$REMOTE" == *"software-factory"* ]]; then
        echo "❌ CRITICAL: Working in software-factory repository!"
        echo "This will corrupt the orchestrator!"
        return 1
    fi
    
    # Check NOT in software-factory directory
    if [[ "$PWD" == *"software-factory"* ]]; then
        echo "❌ CRITICAL: In software-factory directory!"
        return 1
    fi
    
    # Check branch naming
    CURRENT_BRANCH=$(git branch --show-current)
    if [[ "$CURRENT_BRANCH" != *"final-integration"* ]]; then
        echo "⚠️ WARNING: Non-standard branch name: $CURRENT_BRANCH"
    fi
    
    echo "✅ Project integration validation passed"
    return 0
}

# Run validation before EVERY operation
validate_project_integration || exit 1
```

## Integration Process

### Step 1: Pre-Integration Validation
```bash
# Verify all phases completed
for phase in $(seq 1 $TOTAL_PHASES); do
    phase_status=$(yq ".phases.phase_${phase}.status" orchestrator-state.json)
    if [ "$phase_status" != "INTEGRATED" ]; then
        echo "🚨 Cannot integrate project - Phase $phase not integrated"
        exit 1
    fi
done

# Verify architect approval
architect_approval=$(yq '.architect_approval.project_ready' orchestrator-state.json)
if [ "$architect_approval" != "true" ]; then
    echo "🚨 Cannot integrate - Architect approval required"
    exit 1
fi
```

### Step 2: Sequential Phase Merging
- Merge each phase branch in order
- Run full test suite after each merge
- Run build verification after each merge
- Document all conflicts
- Stop on first failure

### Step 3: Final Project Validation & Demo

#### 🏗️ MANDATORY BUILD VERIFICATION
```bash
# After all phases merged
echo "🏗️ Running final project build verification..."

# Clean build
rm -rf dist/ build/ out/ target/
echo "📦 Creating production build..."

# Production build with full logs
if npm run build:prod 2>&1 | tee project-build.log; then
    echo "✅ Production build successful"
else
    echo "❌ Production build failed!"
    exit 1
fi

# Verify build artifacts
if [ -d "dist" ] || [ -d "build" ] || [ -d "out" ] || [ -d "target" ]; then
    echo "✅ Build artifacts created"
    ls -la dist/ build/ out/ target/ 2>/dev/null | tee project-artifacts.log
else
    echo "❌ No build artifacts found!"
    exit 1
fi

# Verify executable
echo "🔍 Verifying executable/deployable artifacts..."
if [ -f "dist/index.js" ] || [ -f "build/app" ] || [ -f "target/*.jar" ]; then
    echo "✅ Deployable artifacts verified"
else
    echo "⚠️ Warning: Verify deployable artifacts manually"
fi
```

#### 🧪 MANDATORY TEST HARNESS
```bash
# Create comprehensive project test harness
cat > project-test-harness.sh << 'EOF'
#!/bin/bash
echo "🧪 PROJECT INTEGRATION TEST SUITE"
echo "=================================="

# Unit tests
echo "📦 Running ALL unit tests..."
if npm test 2>&1 | tee project-unit-tests.log; then
    echo "✅ Unit tests: PASS"
    UNIT_PASS=true
else
    echo "❌ Unit tests: FAIL"
    UNIT_PASS=false
fi

# Integration tests
echo "🔗 Running integration tests..."
if npm run test:integration 2>&1 | tee project-integration-tests.log; then
    echo "✅ Integration tests: PASS"
    INT_PASS=true
else
    echo "❌ Integration tests: FAIL"
    INT_PASS=false
fi

# E2E tests
echo "🌐 Running E2E tests..."
if npm run test:e2e 2>&1 | tee project-e2e-tests.log; then
    echo "✅ E2E tests: PASS"
    E2E_PASS=true
else
    echo "❌ E2E tests: FAIL"
    E2E_PASS=false
fi

# Performance tests
echo "⚡ Running performance tests..."
if npm run test:performance 2>&1 | tee project-performance-tests.log; then
    echo "✅ Performance tests: PASS"
else
    echo "⚠️ Performance tests: CHECK RESULTS"
fi

# Security scan
echo "🔒 Running security scan..."
npm audit --audit-level=moderate 2>&1 | tee project-security.log

echo "=================================="
if $UNIT_PASS && $INT_PASS && $E2E_PASS; then
    echo "✅ ALL CRITICAL TESTS PASSED!"
    exit 0
else
    echo "❌ TEST FAILURES DETECTED!"
    exit 1
fi
EOF

chmod +x project-test-harness.sh
./project-test-harness.sh
```

#### 🎬 MANDATORY PROJECT DEMO
```bash
# Create comprehensive project demo
cat > PROJECT-DEMO.md << 'EOF'
# COMPLETE PROJECT INTEGRATION DEMO

## 🏗️ Build Status
- Production Build: ✅ SUCCESSFUL
- Build Artifacts: Available in dist/
- Build Log: project-build.log
- Deployable: READY

## 🧪 Test Results
- Unit Tests: ✅ ALL PASSING
- Integration Tests: ✅ ALL PASSING  
- E2E Tests: ✅ ALL PASSING
- Performance Tests: ✅ ACCEPTABLE
- Security Scan: ✅ NO HIGH/CRITICAL
- Test Harness: project-test-harness.sh

## 📊 Project Metrics
- Total Phases Integrated: $(yq '.total_phases' orchestrator-state.json)
- Total Features Delivered: $(yq '.features_completed | length' orchestrator-state.json)
- Total Lines of Code: $(git diff main --numstat | awk '{sum+=$1+$2} END {print sum}')
- Test Coverage: $(npm run coverage --silent | grep "All files" | awk '{print $10}')
- Build Size: $(du -sh dist/ 2>/dev/null | awk '{print $1}')

## 🎯 Feature Demonstration

### Phase 1 Features
$(grep "^- " phase-plans/phase1/PHASE-PLAN.md | head -5)

### Phase 2 Features  
$(grep "^- " phase-plans/phase2/PHASE-PLAN.md | head -5)

### Phase 3 Features
$(grep "^- " phase-plans/phase3/PHASE-PLAN.md | head -5)

## 🚀 How to Run Complete Demo

### Quick Start
\`\`\`bash
# 1. Install and build
npm install
npm run build:prod

# 2. Start application
npm start

# 3. Run automated demo
./project-demo.sh

# 4. Verify all features
./verify-project-features.sh
\`\`\`

### Manual Testing
1. Navigate to http://localhost:3000
2. Login with test credentials
3. Verify each major feature area
4. Check performance metrics
5. Review audit logs

## 📁 Demo Artifacts
- Demo Script: project-demo.sh
- Feature Verification: verify-project-features.sh
- Test Results: project-*-tests.log
- Screenshots: demos/project-final/
- Video Walkthrough: demos/project-demo.mp4 (if created)
EOF

# Create automated demo script
cat > project-demo.sh << 'EOF'
#!/bin/bash
echo "🎬 COMPLETE PROJECT DEMONSTRATION"
echo "=================================="
echo ""

# Ensure application is running
if ! curl -s http://localhost:3000/health > /dev/null 2>&1; then
    echo "📦 Starting application..."
    npm start &
    APP_PID=$!
    sleep 10
fi

echo "🔍 Demonstrating Core Features:"
echo ""

# Phase 1 features
echo "📍 Phase 1: Foundation Features"
curl -X GET http://localhost:3000/api/v1/status
echo ""

# Phase 2 features
echo "📍 Phase 2: Core Business Logic"
curl -X POST http://localhost:3000/api/v1/process -d '{"action": "demo"}'
echo ""

# Phase 3 features
echo "📍 Phase 3: Advanced Features"
curl -X GET http://localhost:3000/api/v1/analytics
echo ""

echo "=================================="
echo "✅ PROJECT DEMO COMPLETE!"
echo ""
echo "All features have been demonstrated."
echo "Check logs and outputs for verification."
EOF

chmod +x project-demo.sh

# Create feature verification script
cat > verify-project-features.sh << 'EOF'
#!/bin/bash
echo "🔍 PROJECT FEATURE VERIFICATION"
echo "================================"

TOTAL_FEATURES=0
PASSED_FEATURES=0

verify_feature() {
    local name="$1"
    local check_cmd="$2"
    
    echo -n "Checking $name... "
    ((TOTAL_FEATURES++))
    
    if eval "$check_cmd" > /dev/null 2>&1; then
        echo "✅"
        ((PASSED_FEATURES++))
    else
        echo "❌"
    fi
}

# Verify each major feature
verify_feature "API Health" "curl -s http://localhost:3000/health | grep -q 'ok'"
verify_feature "Authentication" "curl -s http://localhost:3000/api/auth | grep -q 'ready'"
verify_feature "Core Processing" "curl -s http://localhost:3000/api/v1/process"
verify_feature "Data Management" "curl -s http://localhost:3000/api/v1/data"
verify_feature "Analytics" "curl -s http://localhost:3000/api/v1/analytics"
verify_feature "Admin Portal" "curl -s http://localhost:3000/admin | grep -q 'Admin'"

echo "================================"
echo "Results: $PASSED_FEATURES/$TOTAL_FEATURES features working"

if [ $PASSED_FEATURES -eq $TOTAL_FEATURES ]; then
    echo "✅ ALL FEATURES VERIFIED!"
    exit 0
else
    echo "❌ Some features need attention"
    exit 1
fi
EOF

chmod +x verify-project-features.sh

# Line count verification
total_lines=$(git diff main --numstat | awk '{sum+=$1+$2} END {print sum}')
echo "📊 Project total: $total_lines lines"

# Generate final report
cat > PROJECT-INTEGRATION-REPORT.md << EOF
# Project Integration Report

## Repository Verification
- Target Repository: $(git remote get-url origin)
- Verified NOT software-factory: ✅
- Integration Directory: $PWD
- Branch: $(git branch --show-current)

## Phases Integrated
$(for p in $(seq 1 $TOTAL_PHASES); do echo "- Phase $p: INTEGRATED"; done)

## Build & Test Results
- Production Build: ✅ SUCCESSFUL
- Unit Tests: ✅ PASSING
- Integration Tests: ✅ PASSING
- E2E Tests: ✅ PASSING
- Performance: ✅ ACCEPTABLE
- Security: ✅ CLEAR

## Demo & Verification
- Test Harness: ✅ project-test-harness.sh
- Demo Script: ✅ project-demo.sh
- Feature Verification: ✅ verify-project-features.sh
- Demo Documentation: ✅ PROJECT-DEMO.md

## Metrics
- Total Lines: $total_lines
- Test Coverage: $(npm run coverage --silent | grep "All files" | awk '{print $10}')
- Build Size: $(du -sh dist/ 2>/dev/null | awk '{print $1}')
- Features Delivered: $(yq '.features_completed | length' orchestrator-state.json)

## Validation Checklist
- [x] All phases merged successfully
- [x] No unresolved conflicts
- [x] All tests passing
- [x] Production build successful
- [x] Security scan clear
- [x] Repository isolation verified
- [x] Demo artifacts created
- [x] Features demonstrated

## Approval
- Orchestrator: COMPLETE
- Demo Status: READY
- Architect Review: PENDING
EOF
```

### Step 4: Push and Create Final PR
```bash
# Push final integration
git push -u origin "$INTEGRATION_BRANCH"

# Create final project PR
gh pr create \
  --title "${PROJECT_PREFIX} - Complete Project Integration" \
  --body "$(cat PROJECT-INTEGRATION-REPORT.md)" \
  --base main
```

## Failure Conditions

### Critical Failures (Immediate Stop)
- 🚨 Cloning software-factory repository = IMMEDIATE FAIL
- 🚨 Working in software-factory directory = IMMEDIATE FAIL
- 🚨 Using non-isolated workspace = IMMEDIATE FAIL
- 🚨 Missing phase integrations = FAIL
- 🚨 Test failures = FAIL
- 🚨 Build failures = FAIL

### Recovery Protocol
1. **STOP ALL WORK IMMEDIATELY**
2. Document exact failure point
3. Transition to ERROR_RECOVERY
4. Alert architect for assessment
5. Clean ALL integration workspaces
6. Start fresh with new strategy
7. **NEVER continue from corrupted state**

## Success Criteria
- ✅ Correct target repository cloned
- ✅ Working in `/efforts/integration/`
- ✅ NOT in software-factory repository
- ✅ All phases integrated sequentially
- ✅ All tests passing (unit, integration, E2E)
- ✅ Production build successful
- ✅ Security scan clear
- ✅ Final branch pushed
- ✅ PR created to main
- ✅ State file updated

## Penalties
- Wrong repository: **-100% grade** (AUTOMATIC FAILURE)
- Software-factory contamination: **-100% grade** (CRITICAL)
- Non-isolated workspace: **-80% grade**
- Missing phases: **-60% grade**
- Test failures ignored: **-70% grade**
- Build failures: **-65% grade**
- No final PR: **-40% grade**

## Related Rules
- R282: Phase Integration Protocol
- R034: Integration Requirements
- R250: Integration Isolation
- R288/R288: State File Updates
- R256: Phase Assessment Gate

## Critical Examples

### ❌ CATASTROPHIC: Using software-factory repo
```bash
cd /home/vscode/software-factory-template
git checkout -b final-integration  # DESTROYS ORCHESTRATOR!
```

### ❌ CRITICAL: Wrong directory
```bash
cd /efforts/phase3/phase-integration
git checkout -b final-integration  # WRONG LOCATION!
```

### ❌ SEVERE: Reusing workspace
```bash
cd /efforts/phase2/wave3/integration-workspace
# Trying to do project integration in wave workspace
```

### ✅ CORRECT: Proper isolation
```bash
rm -rf ${CLAUDE_PROJECT_DIR}/efforts/integration
mkdir -p ${CLAUDE_PROJECT_DIR}/efforts/integration
git clone [target-repo] ${CLAUDE_PROJECT_DIR}/efforts/integration/project
cd ${CLAUDE_PROJECT_DIR}/efforts/integration/project
# Verify NOT software-factory
REMOTE=$(git remote get-url origin)
[[ "$REMOTE" != *"software-factory"* ]] || exit 1
git checkout -b project-final-integration
```

## Remember
**This is the FINAL integration** - all project code converges here. Using the wrong repository will DESTROY the orchestrator and fail the entire project. ALWAYS verify you're in the target repository, NOT software-factory!