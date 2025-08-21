# Phase Completion Functional Testing Protocol

## Overview
After completing each phase, the orchestrator MUST create a functional test environment that allows the user to see and interact with the TMC features built in that phase.

## Requirements for Each Phase

### Phase 1 Completion - API Types and CRDs
**What to Test**: API definitions, CRD installation, validation
```bash
# Test Harness: phase1-functional-test.sh
- Install CRDs in test cluster
- Create sample TMC resources
- Validate API schemas work
- Show kubectl get/describe commands
```

### Phase 2 Completion - Controllers
**What to Test**: Controller reconciliation, resource management
```bash
# Test Harness: phase2-functional-test.sh
- Deploy controllers
- Create TMC resources that trigger controllers
- Watch reconciliation loops
- Show controller logs and events
```

### Phase 3 Completion - Syncer
**What to Test**: Bidirectional synchronization
```bash
# Test Harness: phase3-functional-test.sh
- Setup virtual and physical clusters
- Deploy syncer
- Create resources to sync
- Demonstrate bidirectional sync
- Show sync status and health
```

### Phase 4 Completion - Features
**What to Test**: All TMC features including cross-workspace
```bash
# Test Harness: phase4-functional-test.sh
- Full TMC deployment
- Cross-workspace placement
- Tunneling functionality
- DNS resolution
- Quota management
```

### Phase 5 Completion - Final Integration
**What to Test**: Complete TMC system
```bash
# Test Harness: phase5-functional-test.sh
- Full multi-cluster demo
- All features integrated
- Performance metrics
- Observability dashboards
```

## Test Harness Creation Process

### 1. Create Test Working Copy
```bash
# After phase integration branch is created
PHASE_NUM=1  # Current phase
TEST_DIR="/workspaces/tests/phase${PHASE_NUM}-functional"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Clone with full checkout for testing
git clone https://github.com/jessesanford/kcp.git .
git checkout phase${PHASE_NUM}-integration

# Build binaries
make build-all
```

### 2. Task @agent-kcp-kubernetes-code-reviewer to Design Test Harness

```markdown
Task for @agent-kcp-kubernetes-code-reviewer:

Design a functional test harness for Phase ${PHASE_NUM} TMC features.

Reference: /workspaces/agent-configs/example-functional-test-scripts/tmc-multi-cluster-demo.sh

Requirements:
1. Interactive script that guides user through features
2. Sets up test environment (kind clusters if needed)
3. Deploys phase-specific TMC components
4. Creates sample resources
5. Demonstrates functionality with pauses for user observation
6. Includes cleanup

Output to: /workspaces/tests/phase${PHASE_NUM}-functional/test-harness.sh

The test should:
- Be runnable on local machine
- Show real TMC features working
- Allow user to interact and explore
- Provide clear output and explanations
- Clean up after completion
```

### 3. Test Harness Template

```bash
#!/bin/bash
# Phase ${PHASE_NUM} TMC Functional Test Harness
# Generated: $(date)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Phase ${PHASE_NUM} TMC Functional Test ===${NC}"
echo "This test will demonstrate the TMC features completed in Phase ${PHASE_NUM}"
echo ""

# Function to pause and wait for user
pause() {
    echo -e "${GREEN}$1${NC}"
    echo "Press Enter to continue..."
    read
}

# Setup test environment
setup_environment() {
    echo -e "${BLUE}Setting up test environment...${NC}"
    
    # Create kind clusters if needed
    if [ ${PHASE_NUM} -ge 3 ]; then
        kind create cluster --name kcp-control
        kind create cluster --name physical-1
    fi
    
    # Build and prepare binaries
    echo "Building TMC components..."
    make build-all
}

# Deploy phase-specific components
deploy_components() {
    echo -e "${BLUE}Deploying Phase ${PHASE_NUM} components...${NC}"
    
    case ${PHASE_NUM} in
        1)
            echo "Installing CRDs..."
            kubectl apply -f config/crds/
            pause "CRDs installed. Let's verify they're registered."
            kubectl get crds | grep tmc
            ;;
        2)
            echo "Deploying controllers..."
            kubectl apply -f config/controllers/
            pause "Controllers deployed. Let's check their status."
            kubectl get pods -n tmc-system
            ;;
        3)
            echo "Deploying syncer..."
            ./bin/syncer --config test-config.yaml &
            SYNCER_PID=$!
            pause "Syncer started. Let's verify it's running."
            ;;
        4)
            echo "Deploying all TMC features..."
            kubectl apply -f config/features/
            pause "Features deployed. Let's test cross-workspace."
            ;;
        5)
            echo "Full TMC system deployment..."
            ./scripts/deploy-all.sh
            pause "Complete system deployed."
            ;;
    esac
}

# Demonstrate functionality
demonstrate_features() {
    echo -e "${BLUE}Demonstrating Phase ${PHASE_NUM} features...${NC}"
    
    # Phase-specific demonstrations
    case ${PHASE_NUM} in
        1) demo_apis ;;
        2) demo_controllers ;;
        3) demo_syncer ;;
        4) demo_features ;;
        5) demo_full_system ;;
    esac
}

# Cleanup
cleanup() {
    echo -e "${BLUE}Cleaning up...${NC}"
    
    if [ ${PHASE_NUM} -ge 3 ]; then
        kind delete cluster --name kcp-control
        kind delete cluster --name physical-1
    fi
    
    if [ ! -z "$SYNCER_PID" ]; then
        kill $SYNCER_PID 2>/dev/null || true
    fi
}

# Main execution
main() {
    trap cleanup EXIT
    
    setup_environment
    pause "Environment ready. Let's deploy TMC components."
    
    deploy_components
    pause "Components deployed. Let's see them in action."
    
    demonstrate_features
    
    echo -e "${GREEN}Test completed successfully!${NC}"
    echo "The Phase ${PHASE_NUM} TMC features are working correctly."
}

# Run the test
main
```

## Integration with Orchestrator Workflow

### After Phase Completion
```python
def complete_phase(phase_num):
    # Regular integration
    create_integration_branch(phase_num)
    
    # NEW: Functional testing
    test_dir = create_test_working_copy(phase_num)
    build_binaries(test_dir)
    
    # Task reviewer to create test harness
    test_harness = task_reviewer_for_test_design(phase_num)
    
    # Notify user
    print(f"""
    Phase {phase_num} Complete!
    
    Integration branch: phase{phase_num}-integration
    Test directory: {test_dir}
    
    To see your TMC features in action:
    cd {test_dir}
    ./test-harness.sh
    
    This will guide you through all Phase {phase_num} features.
    """)
    
    # Wait for user confirmation before proceeding
    input("Press Enter after testing to proceed to Phase {}...".format(phase_num + 1))
```

## User Experience

### What You'll See
1. **Clear Progress**: "Phase 1 Complete! Ready for testing."
2. **Test Location**: Working directory with built binaries
3. **Interactive Script**: Step-by-step demonstration
4. **Real Functionality**: Actual TMC features running
5. **User Control**: Pause points to explore and verify

### Example User Flow
```
Orchestrator: Phase 1 complete! All APIs and CRDs implemented.
Orchestrator: Test harness created at /workspaces/tests/phase1-functional/
User: cd /workspaces/tests/phase1-functional/
User: ./test-harness.sh
Script: "Installing TMC CRDs... Press Enter to continue"
User: [explores with kubectl]
Script: "Creating sample TMC resources..."
User: [verifies resources work]
Script: "Phase 1 testing complete!"
User: [returns to orchestrator to continue]
```

## Success Criteria

Each phase test must:
- ✅ Build successfully from integration branch
- ✅ Deploy without errors
- ✅ Demonstrate phase-specific features
- ✅ Allow user interaction
- ✅ Clean up properly
- ✅ Provide clear feedback

## Test Output Artifacts

Save for each phase:
- `test-harness.sh` - The interactive test script
- `test-results.log` - Output from test run
- `screenshots/` - If UI components exist
- `metrics.json` - Performance/health metrics

This ensures you can see and verify TMC functionality at every phase before proceeding.