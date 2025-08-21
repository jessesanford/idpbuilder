# Phase Completion Functional Testing Protocol

## Overview
After completing each phase, the orchestrator MUST create a functional test environment that allows the user to see and interact with the features built in that phase.

## Purpose
- **Validate Phase Completeness**: Ensure all phase objectives are met
- **Integration Testing**: Verify components work together
- **User Acceptance**: Demonstrate working features to stakeholders
- **Regression Prevention**: Ensure nothing broke from previous phases

## Requirements for Each Phase

### Phase 1 Completion - Core APIs/Contracts
**What to Test**: API definitions, data models, contracts
```bash
# Test Harness: phase1-functional-test.sh
- Validate all API endpoints respond correctly
- Test data model CRUD operations
- Verify contract schemas and validation
- Demonstrate API documentation works
- Show sample requests/responses
```

### Phase 2 Completion - Business Logic/Services
**What to Test**: Service layer, business rules, workflows
```bash
# Test Harness: phase2-functional-test.sh
- Deploy core services
- Execute business workflows
- Validate business rule enforcement
- Show service logs and metrics
- Test error handling and recovery
```

### Phase 3 Completion - Integration Layer
**What to Test**: External integrations, data flow
```bash
# Test Harness: phase3-functional-test.sh
- Setup integration endpoints
- Test data synchronization
- Validate external API connections
- Demonstrate data consistency
- Show integration health status
```

### Phase 4 Completion - User Features
**What to Test**: All user-facing features
```bash
# Test Harness: phase4-functional-test.sh
- Full application deployment
- User authentication/authorization
- Complete user workflows
- Performance under load
- Feature completeness check
```

### Phase 5 Completion - Final Integration
**What to Test**: Complete system with all components
```bash
# Test Harness: phase5-functional-test.sh
- Full end-to-end demo
- All features integrated
- Performance metrics
- Monitoring/observability
- Production readiness check
```

## Test Harness Creation Process

### 1. Create Test Working Copy
```bash
# After phase integration branch is created
PHASE_NUM=1  # Current phase
TEST_DIR="/workspaces/tests/phase${PHASE_NUM}-functional"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Clone/copy with full checkout for testing
git checkout phase${PHASE_NUM}-integration

# Build/prepare test environment
[BUILD_COMMAND]  # e.g., make build, npm build, cargo build
```

### 2. Task Code Reviewer to Design Test Harness

```markdown
Task for @agent-code-reviewer:

Design a functional test harness for Phase ${PHASE_NUM} features.

Requirements:
1. Interactive script that guides user through features
2. Sets up test environment (containers/services if needed)
3. Deploys phase-specific components
4. Creates sample data/resources
5. Demonstrates functionality with pauses for user observation
6. Includes cleanup

Output to: /workspaces/tests/phase${PHASE_NUM}-functional/test-harness.sh

The test should:
- Be runnable on local machine
- Show real features working
- Allow user to interact and explore
- Provide clear output and explanations
- Clean up after completion
```

### 3. Test Harness Template

```bash
#!/bin/bash
# Phase ${PHASE_NUM} Functional Test Harness
# Generated: $(date)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Phase ${PHASE_NUM} Functional Test ===${NC}"
echo "This test will demonstrate the features completed in Phase ${PHASE_NUM}"
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
    
    # Create test infrastructure if needed
    # Examples: Docker containers, test databases, mock services
    if [ ${PHASE_NUM} -ge 3 ]; then
        # Setup integration test environment
        docker-compose up -d test-env
    fi
    
    # Build and prepare application
    echo "Building application components..."
    [BUILD_COMMAND]  # make build, npm build, etc.
}

# Deploy phase-specific components
deploy_components() {
    echo -e "${BLUE}Deploying Phase ${PHASE_NUM} components...${NC}"
    
    case ${PHASE_NUM} in
        1)
            echo "Setting up APIs/Contracts..."
            # Deploy API definitions, schemas, contracts
            ./deploy-apis.sh
            pause "APIs deployed. Let's verify they're accessible."
            ;;
        2)
            echo "Deploying business logic..."
            # Deploy services, workers, processors
            ./deploy-services.sh
            pause "Services deployed. Let's check their status."
            ;;
        3)
            echo "Setting up integrations..."
            # Deploy integration components
            ./deploy-integrations.sh
            pause "Integrations configured. Let's verify connections."
            ;;
        4)
            echo "Deploying user features..."
            # Deploy UI, user-facing components
            ./deploy-features.sh
            pause "Features deployed. Let's test user workflows."
            ;;
        5)
            echo "Full system deployment..."
            ./deploy-all.sh
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
    
    # Stop test infrastructure
    if [ ${PHASE_NUM} -ge 3 ]; then
        docker-compose down
    fi
    
    # Kill any background processes
    if [ ! -z "$BG_PID" ]; then
        kill $BG_PID 2>/dev/null || true
    fi
    
    # Clean test data
    rm -rf test-data/
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
    echo "The Phase ${PHASE_NUM} features are working correctly."
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
2. **Test Location**: Working directory with built application
3. **Interactive Script**: Step-by-step demonstration
4. **Real Functionality**: Actual features running
5. **User Control**: Pause points to explore and verify

### Example User Flow
```
Orchestrator: Phase 1 complete! All APIs implemented.
Orchestrator: Test harness created at /workspaces/tests/phase1-functional/
User: cd /workspaces/tests/phase1-functional/
User: ./test-harness.sh
Script: "Setting up API endpoints... Press Enter to continue"
User: [explores with API client]
Script: "Creating sample data..."
User: [verifies APIs work]
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

This ensures you can see and verify functionality at every phase before proceeding.