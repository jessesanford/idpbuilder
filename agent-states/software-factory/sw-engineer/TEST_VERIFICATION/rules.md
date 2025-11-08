# SW Engineer - TEST_WRITING State Rules

## 🔴🔴🔴 ABSOLUTE REQUIREMENT: STATE MANAGER CONSULTATION 🔴🔴🔴

**SUPREME LAW - R517 - UNIVERSAL STATE MANAGER CONSULTATION LAW**

**BEFORE EXITING THIS STATE, YOU MUST:**

1. **SPAWN STATE MANAGER FOR SHUTDOWN_CONSULTATION** (MANDATORY - NO EXCEPTIONS)
2. **NEVER UPDATE STATE FILES DIRECTLY** (orchestrator-state-v3.json, bug-tracking.json, etc.)
3. **NEVER COMMIT STATE FILES YOURSELF** (State Manager does this atomically)

**FAILURE TO CONSULT STATE MANAGER = IMMEDIATE SYSTEM HALT (-100% GRADE)**

### Enforcement Mechanism

If you attempt to exit this state without spawning State Manager:
- ❌ Pre-commit hooks will REJECT your commit
- ❌ Validation tools will FAIL the build
- ❌ Grading system will assign -100% penalty
- ❌ System will transition to ERROR_RECOVERY

### Required Pattern (COPY THIS EXACTLY)

```bash
# At end of state work, BEFORE any state file updates:

echo "🔴 MANDATORY: Spawning State Manager for SHUTDOWN_CONSULTATION"

# Spawn State Manager (REQUIRED - NOT OPTIONAL)
# Task: state-manager
# State: SHUTDOWN_CONSULTATION
# Current State: [YOUR_CURRENT_STATE]
# Proposed Next State: [YOUR_PROPOSED_NEXT_STATE]
# Work Summary: [SUMMARY_OF_WORK_COMPLETED]

# State Manager will:
# 1. Validate proposed transition against state machine
# 2. Update all 4 state files atomically
# 3. Commit with [R288] tag
# 4. Return REQUIRED next state (may differ from proposal)

# Wait for State Manager response
# Follow State Manager's directive (REQUIRED next state)
# DO NOT proceed until State Manager confirms
```

**YOU MUST NEVER:**
- ❌ Update orchestrator-state-v3.json yourself
- ❌ Update bug-tracking.json yourself
- ❌ Update integration-containers.json yourself
- ❌ Use `jq` to modify state files
- ❌ Use `sed/awk` to modify state files
- ❌ Set `validated_by: "orchestrator"` (must be "state-manager")
- ❌ Skip State Manager consultation "just this once"
- ❌ Think "I'll validate it manually"

**ONLY State Manager may update state files. This is NON-NEGOTIABLE.**

See: `rule-library/R517-universal-state-manager-consultation-law.md`
## State Context
You are focused on writing comprehensive tests to meet coverage requirements and ensure code quality.

## SF 3.0 Test Verification Context

In this state, the SW Engineer verifies test quality within SF 3.0:
- Reads test requirements from effort plan tracked in orchestrator-state-v3.json `metadata_locations.effort_plans`
- Verifies test coverage meets phase/wave requirements defined in orchestrator-state-v3.json
- Documents test results that may be recorded in `bug-tracking.json` if test failures reveal bugs
- Reports test completion status for orchestrator to update in orchestrator-state-v3.json per R288
- All test artifacts and coverage reports stored with proper metadata for tracking

## 🔴🔴🔴 SUPREME LAW R355: PRODUCTION READY TEST CODE 🔴🔴🔴

### EVEN TESTS MUST BE PRODUCTION READY:
- ❌ **No Hardcoded Test Credentials** - Use test configs
- ❌ **No Incomplete Test Stubs** - Every test must run
- ❌ **No TODO/FIXME in Tests** - Complete all test cases
- ❌ **No Static Test Data** - Use configurable fixtures
- ✅ **Mocks ALLOWED in test files** - This is the ONLY exception

### TEST CODE VERIFICATION:
```bash
echo "🔴 R355: VERIFYING TEST CODE QUALITY"
cd $EFFORT_DIR
# Check test files don't have production violations
VIOLATIONS=0
# Note: We check test files but ALLOW mocks in them
grep -r "password.*=.*['\"].*admin\|password123" --include="*_test.go" --include="*_test.py" --include="*.test.js" && VIOLATIONS=1
grep -r "TODO\|FIXME\|INCOMPLETE" --include="*_test.go" --include="*_test.py" --include="*.test.js" && VIOLATIONS=1
grep -r "not.*implemented.*test" --include="*_test.go" --include="*_test.py" --include="*.test.js" && VIOLATIONS=1

if [ $VIOLATIONS -eq 1 ]; then
    echo "🚨 R355 VIOLATION: Test code not production ready!"
    exit 355
fi
echo "✅ R355: Test code is production ready"
echo "✅ Note: Mocks are ALLOWED in test files only"
```

---
### ℹ️ RULE R108.0.0 - TEST_WRITING Rules
**Source:** rule-library/RULE-REGISTRY.md#R108
**Criticality:** INFO - Best practice

TEST WRITING PROTOCOL:
1. Focus exclusively on test development
2. Achieve minimum coverage requirements
3. Write meaningful tests, not just coverage fillers
4. Include unit, integration, and edge case tests
5. Maintain test quality and maintainability
6. Monitor size impact of test code
---

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

## Test Coverage Requirements

---
### 🚨 RULE R032.0.0 - Test Coverage Requirements
**Source:** rule-library/RULE-REGISTRY.md#R032
**Criticality:** CRITICAL - Major impact on grading

COVERAGE REQUIREMENTS BY COMPONENT:

CORE BUSINESS LOGIC: 90%+ coverage
- Domain models and business rules
- Critical calculation and processing logic
- State management and validation

CONTROLLERS: 85%+ coverage
- Reconciliation logic
- CRUD operations
- Error handling paths

API TYPES: 80%+ coverage
- Validation methods
- Default value setting
- Conversion functions

UTILITIES: 85%+ coverage
- Helper functions
- Common operations
- Integration points
---

```go
// Example comprehensive test structure
package controllers_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/types"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/client/fake"
    
    myapi "github.com/example/api/v1"
    "github.com/example/controllers"
)

func TestResourceController_Reconcile_CompleteScenarios(t *testing.T) {
    tests := []struct {
        name              string
        existingObjects   []client.Object
        inputRequest      ctrl.Request
        expectedResult    ctrl.Result
        expectedError     bool
        expectedActions   []string
        setupMocks        func(*testing.T) client.Client
        verifyAssertions  func(*testing.T, client.Client)
    }{
        {
            name: "create_new_resource_success",
            existingObjects: []client.Object{},
            inputRequest: ctrl.Request{
                NamespacedName: types.NamespacedName{
                    Name:      "test-resource",
                    Namespace: "default",
                },
            },
            expectedResult: ctrl.Result{RequeueAfter: time.Minute * 5},
            expectedError:  false,
            expectedActions: []string{"create", "status_update"},
            setupMocks: func(t *testing.T) client.Client {
                scheme := runtime.NewScheme()
                require.NoError(t, myapi.AddToScheme(scheme))
                
                return fake.NewClientBuilder().
                    WithScheme(scheme).
                    Build()
            },
            verifyAssertions: func(t *testing.T, client client.Client) {
                var resource myapi.Resource
                err := client.Get(context.TODO(), types.NamespacedName{
                    Name: "test-resource", Namespace: "default",
                }, &resource)
                
                require.NoError(t, err)
                assert.Equal(t, myapi.ResourcePhaseReady, resource.Status.Phase)
                assert.Len(t, resource.Status.Conditions, 1)
                assert.Equal(t, myapi.ConditionTypeReady, resource.Status.Conditions[0].Type)
            },
        },
        {
            name: "update_existing_resource_spec_change",
            existingObjects: []client.Object{
                &myapi.Resource{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      "existing-resource",
                        Namespace: "default",
                    },
                    Spec: myapi.ResourceSpec{
                        Replicas: 1,
                    },
                    Status: myapi.ResourceStatus{
                        Phase: myapi.ResourcePhaseReady,
                        Conditions: []metav1.Condition{
                            {
                                Type:   myapi.ConditionTypeReady,
                                Status: metav1.ConditionTrue,
                            },
                        },
                    },
                },
            },
            inputRequest: ctrl.Request{
                NamespacedName: types.NamespacedName{
                    Name:      "existing-resource",
                    Namespace: "default",
                },
            },
            expectedResult: ctrl.Result{},
            expectedError:  false,
            setupMocks: func(t *testing.T) client.Client {
                scheme := runtime.NewScheme()
                require.NoError(t, myapi.AddToScheme(scheme))
                
                // Simulate spec change
                resource := &myapi.Resource{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      "existing-resource",
                        Namespace: "default",
                    },
                    Spec: myapi.ResourceSpec{
                        Replicas: 3, // Changed from 1 to 3
                    },
                }
                
                return fake.NewClientBuilder().
                    WithScheme(scheme).
                    WithObjects(resource).
                    Build()
            },
            verifyAssertions: func(t *testing.T, client client.Client) {
                var resource myapi.Resource
                err := client.Get(context.TODO(), types.NamespacedName{
                    Name: "existing-resource", Namespace: "default",
                }, &resource)
                
                require.NoError(t, err)
                assert.Equal(t, int32(3), resource.Spec.Replicas)
                
                // Verify status was updated to reflect spec change
                readyCondition := findCondition(resource.Status.Conditions, myapi.ConditionTypeReady)
                require.NotNil(t, readyCondition)
                assert.Equal(t, metav1.ConditionTrue, readyCondition.Status)
            },
        },
        {
            name: "resource_deletion_with_finalizer",
            existingObjects: []client.Object{
                &myapi.Resource{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:              "delete-resource",
                        Namespace:         "default",
                        DeletionTimestamp: &metav1.Time{Time: time.Now()},
                        Finalizers:        []string{myapi.ResourceFinalizer},
                    },
                    Spec: myapi.ResourceSpec{
                        Replicas: 1,
                    },
                },
            },
            inputRequest: ctrl.Request{
                NamespacedName: types.NamespacedName{
                    Name:      "delete-resource",
                    Namespace: "default",
                },
            },
            expectedResult: ctrl.Result{},
            expectedError:  false,
            setupMocks: func(t *testing.T) client.Client {
                scheme := runtime.NewScheme()
                require.NoError(t, myapi.AddToScheme(scheme))
                
                return fake.NewClientBuilder().
                    WithScheme(scheme).
                    Build()
            },
            verifyAssertions: func(t *testing.T, client client.Client) {
                var resource myapi.Resource
                err := client.Get(context.TODO(), types.NamespacedName{
                    Name: "delete-resource", Namespace: "default",
                }, &resource)
                
                // Resource should be cleaned up (finalizer removed)
                if err == nil {
                    assert.NotContains(t, resource.Finalizers, myapi.ResourceFinalizer)
                }
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            client := tt.setupMocks(t)
            
            // Add existing objects if any
            for _, obj := range tt.existingObjects {
                require.NoError(t, client.Create(context.TODO(), obj))
            }
            
            controller := &controllers.ResourceController{
                Client: client,
                Log:    ctrl.Log.WithName("test"),
                Scheme: runtime.NewScheme(),
            }
            
            // Execute
            result, err := controller.Reconcile(context.TODO(), tt.inputRequest)
            
            // Assert
            if tt.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            assert.Equal(t, tt.expectedResult, result)
            
            // Custom assertions
            if tt.verifyAssertions != nil {
                tt.verifyAssertions(t, client)
            }
        })
    }
}

// Edge case and error condition tests
func TestResourceController_EdgeCases(t *testing.T) {
    tests := []struct {
        name           string
        setupScenario  func(*testing.T) (client.Client, ctrl.Request)
        expectedError  bool
        errorContains  string
    }{
        {
            name: "invalid_resource_spec_validation_error",
            setupScenario: func(t *testing.T) (client.Client, ctrl.Request) {
                scheme := runtime.NewScheme()
                require.NoError(t, myapi.AddToScheme(scheme))
                
                // Create resource with invalid spec
                resource := &myapi.Resource{
                    ObjectMeta: metav1.ObjectMeta{
                        Name:      "invalid-resource",
                        Namespace: "default",
                    },
                    Spec: myapi.ResourceSpec{
                        Replicas: -1, // Invalid negative replicas
                    },
                }
                
                client := fake.NewClientBuilder().
                    WithScheme(scheme).
                    WithObjects(resource).
                    Build()
                
                return client, ctrl.Request{
                    NamespacedName: types.NamespacedName{
                        Name: "invalid-resource", Namespace: "default",
                    },
                }
            },
            expectedError: true,
            errorContains: "validation",
        },
        {
            name: "client_api_error_handling",
            setupScenario: func(t *testing.T) (client.Client, ctrl.Request) {
                // Use a client that returns errors
                errorClient := &ErrorInjectingClient{
                    Client: fake.NewClientBuilder().Build(),
                    GetError: errors.New("API server unavailable"),
                }
                
                return errorClient, ctrl.Request{
                    NamespacedName: types.NamespacedName{
                        Name: "test-resource", Namespace: "default",
                    },
                }
            },
            expectedError: true,
            errorContains: "API server unavailable",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, request := tt.setupScenario(t)
            
            controller := &controllers.ResourceController{
                Client: client,
                Log:    ctrl.Log.WithName("test"),
                Scheme: runtime.NewScheme(),
            }
            
            _, err := controller.Reconcile(context.TODO(), request)
            
            if tt.expectedError {
                require.Error(t, err)
                if tt.errorContains != "" {
                    assert.Contains(t, err.Error(), tt.errorContains)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// Integration tests with real Kubernetes API
func TestResourceController_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests in short mode")
    }
    
    testEnv := &envtest.Environment{
        CRDDirectoryPaths: []string{
            filepath.Join("..", "..", "config", "crd", "bases"),
        },
    }
    
    cfg, err := testEnv.Start()
    require.NoError(t, err)
    defer func() {
        assert.NoError(t, testEnv.Stop())
    }()
    
    // Setup client
    k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
    require.NoError(t, err)
    
    // Test real reconciliation against live API server
    ctx := context.Background()
    
    // Create a resource
    resource := &myapi.Resource{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "integration-test-resource",
            Namespace: "default",
        },
        Spec: myapi.ResourceSpec{
            Replicas: 2,
        },
    }
    
    err = k8sClient.Create(ctx, resource)
    require.NoError(t, err)
    
    // Setup and start controller
    mgr, err := ctrl.NewManager(cfg, ctrl.Options{
        Scheme: scheme.Scheme,
        Port:   0, // Disable webhook server
    })
    require.NoError(t, err)
    
    controller := &controllers.ResourceController{
        Client: k8sClient,
        Log:    ctrl.Log.WithName("integration-test"),
        Scheme: scheme.Scheme,
    }
    
    err = controller.SetupWithManager(mgr)
    require.NoError(t, err)
    
    // Run reconciliation
    result, err := controller.Reconcile(ctx, ctrl.Request{
        NamespacedName: types.NamespacedName{
            Name: "integration-test-resource", Namespace: "default",
        },
    })
    
    assert.NoError(t, err)
    assert.False(t, result.Requeue)
    
    // Verify resource was updated
    var updatedResource myapi.Resource
    err = k8sClient.Get(ctx, types.NamespacedName{
        Name: "integration-test-resource", Namespace: "default",
    }, &updatedResource)
    
    require.NoError(t, err)
    assert.Equal(t, myapi.ResourcePhaseReady, updatedResource.Status.Phase)
    
    // Cleanup
    err = k8sClient.Delete(ctx, &updatedResource)
    assert.NoError(t, err)
}

// Helper functions and test utilities
func findCondition(conditions []metav1.Condition, conditionType string) *metav1.Condition {
    for i := range conditions {
        if conditions[i].Type == conditionType {
            return &conditions[i]
        }
    }
    return nil
}

// ErrorInjectingClient wraps a client to inject errors for testing
type ErrorInjectingClient struct {
    client.Client
    GetError    error
    CreateError error
    UpdateError error
    DeleteError error
}

func (c *ErrorInjectingClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
    if c.GetError != nil {
        return c.GetError
    }
    return c.Client.Get(ctx, key, obj)
}

func (c *ErrorInjectingClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
    if c.CreateError != nil {
        return c.CreateError
    }
    return c.Client.Create(ctx, obj, opts...)
}
```

## Test Quality Standards

---
### ℹ️ RULE R037.0.0 - Pattern Compliance
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

TEST QUALITY REQUIREMENTS:

TEST STRUCTURE:
- Use table-driven tests for multiple scenarios
- Clear test names describing exact scenario
- Proper setup, execution, assertion structure
- Independent tests with no cross-dependencies

ASSERTION QUALITY:
- Specific assertions, not just "no error"
- Verify both positive and negative cases
- Check all relevant output properties
- Include boundary condition testing

MOCKING AND ISOLATION:
- Mock external dependencies appropriately
- Test units in isolation where possible
- Use dependency injection for testability
- Avoid testing implementation details
---

```go
// Example of comprehensive API validation tests
func TestResourceValidation_ComprehensiveCoverage(t *testing.T) {
    tests := []struct {
        name          string
        resource      *myapi.Resource
        expectedValid bool
        expectedError string
    }{
        {
            name: "valid_resource_all_fields",
            resource: &myapi.Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "valid-resource",
                    Namespace: "default",
                },
                Spec: myapi.ResourceSpec{
                    Replicas:    3,
                    Image:       "nginx:1.21",
                    Environment: map[string]string{"ENV": "test"},
                    Resources: corev1.ResourceRequirements{
                        Requests: corev1.ResourceList{
                            corev1.ResourceCPU:    resource.MustParse("100m"),
                            corev1.ResourceMemory: resource.MustParse("128Mi"),
                        },
                    },
                },
            },
            expectedValid: true,
        },
        {
            name: "invalid_negative_replicas",
            resource: &myapi.Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "invalid-replicas",
                    Namespace: "default",
                },
                Spec: myapi.ResourceSpec{
                    Replicas: -1,
                    Image:    "nginx:1.21",
                },
            },
            expectedValid: false,
            expectedError: "replicas must be non-negative",
        },
        {
            name: "invalid_empty_image",
            resource: &myapi.Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "invalid-image",
                    Namespace: "default",
                },
                Spec: myapi.ResourceSpec{
                    Replicas: 1,
                    Image:    "",
                },
            },
            expectedValid: false,
            expectedError: "image cannot be empty",
        },
        {
            name: "boundary_max_replicas",
            resource: &myapi.Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "max-replicas",
                    Namespace: "default",
                },
                Spec: myapi.ResourceSpec{
                    Replicas: 100, // At boundary
                    Image:    "nginx:1.21",
                },
            },
            expectedValid: true,
        },
        {
            name: "invalid_exceeds_max_replicas",
            resource: &myapi.Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "too-many-replicas",
                    Namespace: "default",
                },
                Spec: myapi.ResourceSpec{
                    Replicas: 101, // Exceeds limit
                    Image:    "nginx:1.21",
                },
            },
            expectedValid: false,
            expectedError: "replicas cannot exceed 100",
        },
        {
            name: "default_values_applied",
            resource: &myapi.Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "default-values",
                    Namespace: "default",
                },
                Spec: myapi.ResourceSpec{
                    // Only required fields, test defaults
                    Image: "nginx:1.21",
                },
            },
            expectedValid: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Execute validation
            err := tt.resource.ValidateCreate()
            
            if tt.expectedValid {
                assert.NoError(t, err, "Expected resource to be valid")
                
                // Verify defaults were applied
                if tt.name == "default_values_applied" {
                    assert.Equal(t, int32(1), tt.resource.Spec.Replicas, "Default replicas should be 1")
                    assert.NotEmpty(t, tt.resource.Spec.Environment["DEFAULT_ENV"], "Default environment should be set")
                }
            } else {
                require.Error(t, err, "Expected validation error")
                assert.Contains(t, err.Error(), tt.expectedError, "Error message should contain expected text")
            }
        })
    }
}

// Performance and load testing
func TestResourceController_Performance(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping performance tests in short mode")
    }
    
    scheme := runtime.NewScheme()
    require.NoError(t, myapi.AddToScheme(scheme))
    
    client := fake.NewClientBuilder().
        WithScheme(scheme).
        Build()
    
    controller := &controllers.ResourceController{
        Client: client,
        Log:    ctrl.Log.WithName("perf-test"),
        Scheme: scheme,
    }
    
    // Test reconciliation performance with many resources
    resourceCount := 100
    
    // Create resources
    for i := 0; i < resourceCount; i++ {
        resource := &myapi.Resource{
            ObjectMeta: metav1.ObjectMeta{
                Name:      fmt.Sprintf("perf-resource-%d", i),
                Namespace: "default",
            },
            Spec: myapi.ResourceSpec{
                Replicas: 1,
                Image:    "nginx:1.21",
            },
        }
        require.NoError(t, client.Create(context.TODO(), resource))
    }
    
    // Measure reconciliation performance
    start := time.Now()
    
    for i := 0; i < resourceCount; i++ {
        request := ctrl.Request{
            NamespacedName: types.NamespacedName{
                Name:      fmt.Sprintf("perf-resource-%d", i),
                Namespace: "default",
            },
        }
        
        result, err := controller.Reconcile(context.TODO(), request)
        assert.NoError(t, err)
        assert.False(t, result.Requeue)
    }
    
    elapsed := time.Since(start)
    
    // Performance assertions
    avgReconcileTime := elapsed / time.Duration(resourceCount)
    assert.Less(t, avgReconcileTime, 100*time.Millisecond, 
        "Average reconciliation time should be under 100ms")
    
    t.Logf("Reconciled %d resources in %v (avg: %v per resource)", 
        resourceCount, elapsed, avgReconcileTime)
}
```

## Test Coverage Measurement

---
### ℹ️ RULE R154.0.0 - Test Coverage Achievement
**Source:** rule-library/RULE-REGISTRY.md#R154
**Criticality:** INFO - Best practice

COVERAGE MEASUREMENT PROTOCOL:

MEASUREMENT TOOLS:
- Use go test -cover for basic coverage
- Use go test -coverprofile for detailed analysis
- Generate HTML coverage reports for visualization
- Track coverage trends over time

MEASUREMENT FREQUENCY:
- After each test writing session
- Before transitioning back to implementation
- When coverage targets are not met
- Before code review requests
---

```bash
#!/bin/bash
# Comprehensive test coverage measurement script

PACKAGE_PATH=${1:-"./..."}
COVERAGE_TARGET=${2:-85}
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

echo "🧪 COMPREHENSIVE TEST COVERAGE ANALYSIS - $TIMESTAMP"
echo "Package Path: $PACKAGE_PATH"
echo "Coverage Target: ${COVERAGE_TARGET}%"
echo ""

# Run tests with coverage
echo "📊 RUNNING TESTS WITH COVERAGE MEASUREMENT"
go test -v -race -coverprofile="$COVERAGE_FILE" -covermode=atomic "$PACKAGE_PATH"
TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -ne 0 ]; then
    echo "❌ TESTS FAILED - Coverage analysis aborted"
    echo "Fix failing tests before analyzing coverage"
    exit $TEST_EXIT_CODE
fi

# Generate coverage report
echo ""
echo "📈 COVERAGE ANALYSIS"
COVERAGE_RESULT=$(go tool cover -func="$COVERAGE_FILE")
echo "$COVERAGE_RESULT"

# Extract total coverage percentage
TOTAL_COVERAGE=$(echo "$COVERAGE_RESULT" | tail -1 | awk '{print $3}' | sed 's/%//')

echo ""
echo "📋 COVERAGE SUMMARY"
echo "Total Coverage: ${TOTAL_COVERAGE}%"
echo "Target Coverage: ${COVERAGE_TARGET}%"

# Generate HTML report
echo ""
echo "🌐 GENERATING HTML COVERAGE REPORT"
go tool cover -html="$COVERAGE_FILE" -o "$COVERAGE_HTML"
echo "HTML Report: $COVERAGE_HTML"

# Package-level coverage analysis
echo ""
echo "📦 PACKAGE-LEVEL COVERAGE BREAKDOWN"
go tool cover -func="$COVERAGE_FILE" | grep -E "\.(go):" | \
    awk '{print $1, $3}' | sort -t/ -k1,1 | \
    awk -F/ '{
        package = $1; 
        for(i=2; i<NF; i++) package = package "/" $i; 
        gsub(/:[0-9]+$/, "", package);
        coverage[package] += $NF; 
        count[package]++
    } 
    END {
        for(p in coverage) 
            printf "  %-40s %.1f%%\n", p, coverage[p]/count[p]
    }' | sort

# Coverage quality assessment
echo ""
echo "🎯 COVERAGE QUALITY ASSESSMENT"

# Check if target is met
if (( $(echo "$TOTAL_COVERAGE >= $COVERAGE_TARGET" | bc -l) )); then
    echo "✅ COVERAGE TARGET MET: ${TOTAL_COVERAGE}% >= ${COVERAGE_TARGET}%"
    COVERAGE_STATUS="PASSED"
else
    SHORTAGE=$(echo "$COVERAGE_TARGET - $TOTAL_COVERAGE" | bc -l)
    echo "❌ COVERAGE TARGET NOT MET: ${TOTAL_COVERAGE}% < ${COVERAGE_TARGET}%"
    echo "   Coverage shortage: ${SHORTAGE}%"
    COVERAGE_STATUS="FAILED"
fi

# Analyze uncovered critical functions
echo ""
echo "🔍 CRITICAL UNCOVERED CODE ANALYSIS"
UNCOVERED_CRITICAL=$(go tool cover -func="$COVERAGE_FILE" | \
    grep -E "(reconcile|validate|controller|webhook)" | \
    awk '$3 == "0.0%" {print $1, $3}')

if [ -n "$UNCOVERED_CRITICAL" ]; then
    echo "⚠️ CRITICAL FUNCTIONS WITH 0% COVERAGE:"
    echo "$UNCOVERED_CRITICAL"
else
    echo "✅ All critical functions have some test coverage"
fi

# Log coverage to work log
echo "- [$TIMESTAMP] Test coverage: ${TOTAL_COVERAGE}% (target: ${COVERAGE_TARGET}%) - $COVERAGE_STATUS" >> work-log.md

# Return appropriate exit code
if [ "$COVERAGE_STATUS" = "PASSED" ]; then
    exit 0
else
    exit 1
fi
```

## Test Development Strategy

```python
def develop_comprehensive_test_suite(implementation_analysis):
    """Develop comprehensive test suite based on implementation analysis"""
    
    print("🧪 DEVELOPING COMPREHENSIVE TEST SUITE")
    
    test_plan = {
        'unit_tests': plan_unit_tests(implementation_analysis),
        'integration_tests': plan_integration_tests(implementation_analysis),
        'edge_case_tests': plan_edge_case_tests(implementation_analysis),
        'performance_tests': plan_performance_tests(implementation_analysis),
        'coverage_targets': calculate_coverage_targets(implementation_analysis)
    }
    
    # Generate test implementation order
    test_implementation_order = prioritize_test_implementation(test_plan)
    
    print(f"Test plan generated with {len(test_implementation_order)} test groups")
    print("\n📋 TEST IMPLEMENTATION ORDER:")
    
    for i, test_group in enumerate(test_implementation_order, 1):
        print(f"{i}. {test_group['category']}: {test_group['description']}")
        print(f"   Estimated lines: {test_group['estimated_lines']}")
        print(f"   Coverage impact: +{test_group['coverage_increase']}%")
        print(f"   Priority: {test_group['priority']}")
        print()
    
    return {
        'test_plan': test_plan,
        'implementation_order': test_implementation_order,
        'estimated_total_lines': sum(group['estimated_lines'] for group in test_implementation_order),
        'projected_coverage': calculate_projected_coverage(test_plan, implementation_analysis)
    }

def plan_unit_tests(analysis):
    """Plan comprehensive unit tests based on implementation analysis"""
    
    unit_tests = []
    
    # Controller tests
    controllers = analysis.get('controllers', [])
    for controller in controllers:
        unit_tests.append({
            'category': 'CONTROLLER_UNIT',
            'target': controller['name'],
            'scenarios': [
                'successful_reconciliation',
                'resource_creation',
                'resource_update', 
                'resource_deletion_with_finalizers',
                'error_handling',
                'status_updates',
                'requeue_scenarios'
            ],
            'estimated_lines': 200,
            'coverage_increase': 15
        })
    
    # API type validation tests  
    api_types = analysis.get('api_types', [])
    for api_type in api_types:
        unit_tests.append({
            'category': 'API_VALIDATION',
            'target': api_type['name'],
            'scenarios': [
                'valid_spec_acceptance',
                'invalid_spec_rejection',
                'boundary_value_testing',
                'default_value_application',
                'conversion_functions'
            ],
            'estimated_lines': 120,
            'coverage_increase': 10
        })
    
    # Business logic tests
    business_logic = analysis.get('business_logic', [])
    for logic_component in business_logic:
        unit_tests.append({
            'category': 'BUSINESS_LOGIC',
            'target': logic_component['name'],
            'scenarios': [
                'normal_operation',
                'edge_cases',
                'error_conditions',
                'state_transitions',
                'calculation_accuracy'
            ],
            'estimated_lines': 150,
            'coverage_increase': 20
        })
    
    return unit_tests

def plan_integration_tests(analysis):
    """Plan integration tests for component interactions"""
    
    integration_tests = []
    
    # Controller-to-API integration
    if analysis.get('controllers') and analysis.get('api_types'):
        integration_tests.append({
            'category': 'CONTROLLER_API_INTEGRATE_WAVE_EFFORTS',
            'description': 'End-to-end controller reconciliation with real API objects',
            'test_environment': 'envtest',
            'scenarios': [
                'full_lifecycle_management',
                'cross_resource_dependencies',
                'webhook_validation_integration',
                'status_propagation'
            ],
            'estimated_lines': 180,
            'coverage_increase': 12
        })
    
    # Webhook integration
    webhooks = analysis.get('webhooks', [])
    for webhook in webhooks:
        integration_tests.append({
            'category': 'WEBHOOK_INTEGRATE_WAVE_EFFORTS',
            'target': webhook['name'],
            'description': f'Integration testing for {webhook["name"]} webhook',
            'scenarios': [
                'admission_validation',
                'mutation_logic',
                'error_response_handling',
                'certificate_management'
            ],
            'estimated_lines': 140,
            'coverage_increase': 8
        })
    
    return integration_tests

def prioritize_test_implementation(test_plan):
    """Prioritize test implementation based on coverage impact and risk"""
    
    all_tests = []
    
    # Add all test categories with priority scoring
    for unit_test in test_plan['unit_tests']:
        priority_score = calculate_test_priority(unit_test)
        unit_test['priority_score'] = priority_score
        unit_test['priority'] = get_priority_level(priority_score)
        all_tests.append(unit_test)
    
    for integration_test in test_plan['integration_tests']:
        priority_score = calculate_test_priority(integration_test)
        integration_test['priority_score'] = priority_score
        integration_test['priority'] = get_priority_level(priority_score)
        all_tests.append(integration_test)
    
    # Sort by priority score (highest first)
    all_tests.sort(key=lambda x: x['priority_score'], reverse=True)
    
    return all_tests

def calculate_test_priority(test_item):
    """Calculate priority score for test implementation"""
    
    score = 0
    
    # Coverage impact weight (0-40 points)
    coverage_increase = test_item.get('coverage_increase', 0)
    score += min(40, coverage_increase * 2)
    
    # Risk mitigation weight (0-30 points)
    if test_item.get('category') in ['CONTROLLER_UNIT', 'BUSINESS_LOGIC']:
        score += 30  # High-risk components
    elif test_item.get('category') in ['API_VALIDATION']:
        score += 25  # Medium-high risk
    elif test_item.get('category') in ['WEBHOOK_INTEGRATE_WAVE_EFFORTS']:
        score += 20  # Medium risk
    else:
        score += 15  # Lower risk
    
    # Implementation complexity weight (0-20 points, inversely related)
    estimated_lines = test_item.get('estimated_lines', 100)
    if estimated_lines < 100:
        score += 20  # Quick wins
    elif estimated_lines < 150:
        score += 15
    elif estimated_lines < 200:
        score += 10
    else:
        score += 5   # Complex implementations
    
    # Test category bonus (0-10 points)
    if test_item.get('category') == 'CONTROLLER_UNIT':
        score += 10  # Controllers are critical
    elif test_item.get('category') == 'API_VALIDATION':
        score += 8   # API validation is important
    elif test_item.get('category') == 'BUSINESS_LOGIC':
        score += 9   # Business logic is critical
    
    return score

def get_priority_level(score):
    """Convert priority score to level"""
    
    if score >= 80:
        return 'CRITICAL'
    elif score >= 65:
        return 'HIGH'
    elif score >= 50:
        return 'MEDIUM'
    else:
        return 'LOW'
```

## Test Size Management

```python
def monitor_test_code_size_impact(test_session_data):
    """Monitor the size impact of test code during development"""
    
    print("📏 MONITORING_SWE_PROGRESS TEST CODE SIZE IMPACT")
    
    # Measure current total size
    total_size = measure_current_total_size()
    
    # Separate implementation and test sizes
    implementation_size = measure_implementation_code_size()
    test_size = measure_test_code_size()
    
    size_analysis = {
        'total_size': total_size,
        'implementation_size': implementation_size,
        'test_size': test_size,
        'test_ratio': (test_size / implementation_size * 100) if implementation_size > 0 else 0,
        'size_limit': 800,
        'remaining_capacity': 800 - total_size,
        'test_efficiency': calculate_test_efficiency(test_size, test_session_data)
    }
    
    # Assessment
    if size_analysis['total_size'] > 800:
        status = 'SIZE_VIOLATION'
        urgency = 'CRITICAL'
    elif size_analysis['total_size'] > 750:
        status = 'SIZE_DANGER'
        urgency = 'HIGH'
    elif size_analysis['test_ratio'] > 50:  # Test code > 50% of implementation
        status = 'TEST_HEAVY'
        urgency = 'MEDIUM'
    else:
        status = 'COMPLIANT'
        urgency = 'LOW'
    
    size_analysis['status'] = status
    size_analysis['urgency'] = urgency
    
    print(f"Total Size: {size_analysis['total_size']}/800 lines")
    print(f"Implementation: {implementation_size} lines")
    print(f"Tests: {test_size} lines ({size_analysis['test_ratio']:.1f}% ratio)")
    print(f"Status: {status}")
    
    if status != 'COMPLIANT':
        print(f"⚠️ SIZE CONCERN: {status}")
        recommendations = generate_test_size_recommendations(size_analysis)
        for rec in recommendations:
            print(f"  • {rec}")
    
    return size_analysis

def calculate_test_efficiency(test_lines, session_data):
    """Calculate test efficiency metrics"""
    
    coverage_achieved = session_data.get('coverage_increase', 0)
    test_cases_written = session_data.get('test_cases_written', 0)
    
    if test_lines == 0:
        return {'lines_per_coverage_percent': 0, 'lines_per_test_case': 0}
    
    return {
        'lines_per_coverage_percent': test_lines / coverage_achieved if coverage_achieved > 0 else 0,
        'lines_per_test_case': test_lines / test_cases_written if test_cases_written > 0 else 0,
        'efficiency_score': calculate_efficiency_score(test_lines, coverage_achieved, test_cases_written)
    }

def calculate_efficiency_score(test_lines, coverage_increase, test_cases):
    """Calculate overall test efficiency score (0-100)"""
    
    if test_lines == 0:
        return 0
    
    # Ideal metrics (based on Go testing best practices)
    ideal_lines_per_coverage = 8   # 8 lines of test per 1% coverage
    ideal_lines_per_test_case = 15 # 15 lines per test case average
    
    # Calculate efficiency components
    coverage_efficiency = min(100, (ideal_lines_per_coverage / (test_lines / coverage_increase)) * 100) if coverage_increase > 0 else 0
    case_efficiency = min(100, (ideal_lines_per_test_case / (test_lines / test_cases)) * 100) if test_cases > 0 else 0
    
    # Weighted average
    efficiency_score = (coverage_efficiency * 0.6) + (case_efficiency * 0.4)
    
    return efficiency_score
```

## State Transitions

From TEST_WRITING state:
- **COVERAGE_ACHIEVED** → IMPLEMENTATION (Return to implementation with adequate tests)
- **SIZE_LIMIT_APPROACHED** → MEASURE_SIZE (Check size compliance)
- **TEST_FAILURES** → FIX_ISSUES (Address test issues before continuing)
- **COVERAGE_INSUFFICIENT** → Continue TEST_WRITING (More tests needed)
- **OPTIMIZATION_NEEDED** → FIX_ISSUES (Refactor verbose tests)

## Test Quality Validation

```python
def validate_test_quality(test_files):
    """Validate the quality of written tests"""
    
    quality_metrics = {
        'test_file_count': len(test_files),
        'total_test_cases': 0,
        'assertion_count': 0,
        'mock_usage': 0,
        'edge_case_coverage': 0,
        'error_path_coverage': 0,
        'quality_issues': []
    }
    
    for test_file in test_files:
        file_analysis = analyze_test_file_quality(test_file)
        
        quality_metrics['total_test_cases'] += file_analysis['test_cases']
        quality_metrics['assertion_count'] += file_analysis['assertions']
        quality_metrics['mock_usage'] += file_analysis['mocks']
        quality_metrics['edge_case_coverage'] += file_analysis['edge_cases']
        quality_metrics['error_path_coverage'] += file_analysis['error_paths']
        quality_metrics['quality_issues'].extend(file_analysis['issues'])
    
    # Calculate quality score
    quality_score = calculate_test_quality_score(quality_metrics)
    
    return {
        'metrics': quality_metrics,
        'quality_score': quality_score,
        'recommendations': generate_test_quality_recommendations(quality_metrics)
    }

def analyze_test_file_quality(test_file_path):
    """Analyze individual test file for quality indicators"""
    
    with open(test_file_path, 'r') as f:
        content = f.read()
    
    analysis = {
        'test_cases': len(re.findall(r'func Test\w+', content)),
        'assertions': len(re.findall(r'assert\.|require\.', content)),
        'mocks': len(re.findall(r'mock|Mock|fake|Fake', content)),
        'edge_cases': count_edge_case_tests(content),
        'error_paths': count_error_path_tests(content),
        'issues': identify_test_quality_issues(content)
    }
    
    return analysis


## 🔴🔴🔴 R405 - MANDATORY AUTOMATION CONTINUATION FLAG 🔴🔴🔴

**SUPREME LAW - PENALTY FOR VIOLATION: -100% GRADE**

### YOU MUST OUTPUT THE CONTINUATION FLAG AS YOUR LAST ACTION

**EVERY STATE COMPLETION MUST END WITH EXACTLY ONE OF:**
```bash
# If state completed successfully and factory should continue:
echo "CONTINUE-SOFTWARE-FACTORY=TRUE"

# If error/block/manual review needed:
echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
```

### CRITICAL REQUIREMENTS:
1. **ABSOLUTE LAST OUTPUT**: This MUST be the very last line of output before state completion
2. **NO TEXT AFTER**: No explanations, summaries, or any text after the flag
3. **EXACTLY AS SHOWN**: Use exact format - no variations like CONTINUE-ORCHESTRATING
4. **ALWAYS REQUIRED**: Every single state must output this flag
5. **GREPPABLE**: Must be on its own line for automation parsing

### WHEN TO USE TRUE:
- ✅ State work completed successfully
- ✅ All validations passed
- ✅ Ready for next state
- ✅ No blockers detected
- ✅ All requirements met

### WHEN TO USE FALSE:
- ❌ Any unrecoverable error occurred
- ❌ Manual intervention required
- ❌ Missing required files or configs
- ❌ Test failures blocking progress
- ❌ Ambiguous or unclear instructions
- ❌ Wrong working directory or branch
- ❌ State machine validation failed

### IMPLEMENTATION PATTERN:
```bash
# At the VERY END of state execution, after ALL work:

# Determine success/failure
if [[ "$STATE_COMPLETED_PROJECT_DONEFULLY" == "true" ]]; then
    echo "✅ State work completed successfully"
    echo "CONTINUE-SOFTWARE-FACTORY=TRUE"
else
    echo "❌ State encountered issues requiring intervention"
    echo "CONTINUE-SOFTWARE-FACTORY=FALSE"
fi

# DO NOT OUTPUT ANYTHING AFTER THE FLAG!
```

### GRADING IMPACT:
- **Missing flag**: -100% AUTOMATIC FAILURE
- **Text after flag**: -50% penalty
- **Wrong format**: -100% AUTOMATIC FAILURE
- **Multiple flags**: -50% penalty

**See: $CLAUDE_PROJECT_DIR/rule-library/R405-automation-continuation-flag.md**

