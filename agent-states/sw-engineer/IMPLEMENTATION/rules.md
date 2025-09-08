# SW Engineer - IMPLEMENTATION State Rules

## State Context
You are actively implementing code for an effort, following the implementation plan and maintaining quality standards.

## 🔴🔴🔴 PARAMOUNT: Repository Separation (R251 & R309) 🔴🔴🔴

### R251: Universal Repository Separation Law
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R251-REPOSITORY-SEPARATION-LAW.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation
**YOU ARE IN**: TARGET REPOSITORY CLONE (under /efforts/)
**NOT IN**: Software Factory repo (where orchestrator-state.json lives)

### R309: Never Create Efforts in SF Repo
**File**: `$CLAUDE_PROJECT_DIR/rule-library/R309-never-create-efforts-in-sf-repo.md`
**Criticality**: PARAMOUNT - Automatic -100% failure for violation

**VERIFY YOU'RE IN THE RIGHT PLACE:**
```bash
echo "🔴 R251/R309: Verifying repository location..."
if [ -f "orchestrator-state.json" ] || [ -f ".claude/CLAUDE.md" ]; then
    echo "🔴🔴🔴 FATAL: You're in Software Factory repo!"
    echo "NEVER WRITE CODE HERE! This is for planning only!"
    exit 309
fi

if [ -f "go.mod" ] || [ -f "Makefile" ] || [ -f "package.json" ]; then
    echo "✅ Confirmed: In TARGET repository clone"
    echo "✅ This is where code implementation happens"
else
    echo "❌ WARNING: Cannot confirm repository type"
    pwd
fi

# Verify we're under /efforts/
if [[ "$(pwd)" == *"/efforts/"* ]]; then
    echo "✅ Confirmed: Under /efforts/ directory structure"
else
    echo "❌ FATAL: Not in an effort directory!"
    exit 251
fi
```

**REMEMBER:**
- ✅ Write code in TARGET repo clones under /efforts/
- ❌ NEVER write code in Software Factory instance
- ❌ NEVER create orchestrator-state.json in TARGET repo
- ❌ NEVER create .claude/ configs in TARGET repo

## 🔴🔴🔴 MANDATORY: Effort Scope Strict Adherence (R311) 🔴🔴🔴

**CRITICAL: You MUST implement EXACTLY what's specified in the effort plan!**

### Before Starting ANY Code:
```bash
echo "📋 EXTRACTING SCOPE BOUNDARIES FROM EFFORT PLAN..."
echo "================================================"

# Extract DO NOT IMPLEMENT list
# Use the plan file found during INIT (stored in IMPLEMENTATION_PLAN env var)
PLAN_FILE="${IMPLEMENTATION_PLAN:-IMPLEMENTATION-PLAN.md}"
grep -A 15 "DO NOT IMPLEMENT\|SCOPE BOUNDARIES" "$PLAN_FILE"

# Count specified functions
FUNC_COUNT=$(grep -c "Function.*lines" "$PLAN_FILE")
echo "✅ Functions to implement: EXACTLY $FUNC_COUNT"

# Acknowledge scope
echo "I ACKNOWLEDGE:"
echo "  ✅ Will implement ONLY listed functions"
echo "  ✅ Will NOT add features not in plan"
echo "  ✅ Will document any essential additions"
echo "  ❌ Will NOT over-engineer or 'complete' the feature"
```

### During Implementation - Scope Monitoring:
```bash
# After each component, verify scope adherence
CHECK_SCOPE() {
    ACTUAL_FUNCS=$(grep -c "^func [A-Z]" *.go 2>/dev/null || echo 0)
    EXPECTED_FUNCS=$(grep -c "Function.*~.*lines" "$PLAN_FILE")
    
    if [ "$ACTUAL_FUNCS" -gt "$EXPECTED_FUNCS" ]; then
        echo "❌ SCOPE VIOLATION: $ACTUAL_FUNCS functions > $EXPECTED_FUNCS expected!"
        echo "Remove extra functions or document justification!"
        return 1
    fi
    echo "✅ Scope check passed: $ACTUAL_FUNCS/$EXPECTED_FUNCS functions"
}
```

### If Addition Seems Necessary:
```bash
# STOP and document BEFORE implementing
cat >> IMPLEMENTATION-REPORT.md << EOF
## Justified Addition
- **What**: [Function/feature name]
- **Why Essential**: [Not just "nice to have"]
- **What Breaks Without It**: [Specific failure]
- **Lines Added**: [Estimated count]
- **Decision**: [Proceeding with justification]
EOF
```

**PENALTY**: Adding unrequested features = -50% to -100% grade failure!

## 🔴🔴🔴 MANDATORY: Verify Incremental Base (R308)
**Your effort MUST be based on the latest integrated code!**
```bash
# Verify you're building on previous work
echo "🔍 Verifying incremental base branch..."
BRANCH=$(git rev-parse --abbrev-ref HEAD)
echo "Current branch: $BRANCH"

# Check what we're based on
MERGE_BASE=$(git merge-base HEAD origin/main)
COMMITS_SINCE_MAIN=$(git rev-list --count $MERGE_BASE..HEAD)

echo "📊 Commits ahead of main: $COMMITS_SINCE_MAIN"

# For Wave 2+, should include previous wave work
if [[ "$BRANCH" =~ wave([0-9]+) ]]; then
    WAVE_NUM="${BASH_REMATCH[1]}"
    if [[ $WAVE_NUM -gt 1 && $COMMITS_SINCE_MAIN -eq 0 ]]; then
        echo "⚠️ WARNING: Not based on previous wave integration!"
        echo "This violates R308 - Incremental Branching"
    else
        echo "✅ Building on previous integrated work"
    fi
fi
```

## 🔴 MANDATORY: Verify Directory Isolation (R209)
```bash
# VERIFY YOU'RE STILL IN YOUR EFFORT DIRECTORY
echo "📍 Current directory: $(pwd)"
if [ -z "$EFFORT_ISOLATION_DIR" ]; then
    echo "❌ FATAL: EFFORT_ISOLATION_DIR not set! Run INIT checks!"
    exit 1
fi

if [ "$(pwd)" != "$EFFORT_ISOLATION_DIR" ]; then
    echo "❌❌❌ R209 VIOLATION: You've left your effort directory!"
    echo "   Should be in: $EFFORT_ISOLATION_DIR"
    echo "   Currently in: $(pwd)"
    exit 1
fi

echo "✅ Confirmed in effort directory: $EFFORT_ISOLATION_DIR"
echo "🔒 Directory lock active - cannot cd out"
```

---
### ℹ️ RULE R106.0.0 - IMPLEMENTATION Rules
**Source:** rule-library/RULE-REGISTRY.md#R106
**Criticality:** INFO - Best practice

IMPLEMENTATION PROTOCOL:
1. Follow implementation plan exactly as specified
2. Measure size every 200 lines of code added
3. Write tests alongside implementation
4. Update work log with progress every 30-60 minutes
5. Stop immediately if approaching 800-line limit
6. Commit frequently with descriptive messages
---

## Size Monitoring During Implementation

---
### 🚨🚨🚨 RULE R007.0.0 - Size Limit Enforcement
**Source:** rule-library/RULE-REGISTRY.md#R007
**Criticality:** BLOCKING - Failure = Immediate Stop (exit 1)

SIZE MONITORING REQUIREMENTS:

⚠️ CRITICAL: Use ONLY line-counter.sh tool from PROJECT_ROOT/tools/
⚠️ NEVER count lines manually or with other tools
⚠️ Exclude generated code (zz_generated*, *.pb.go, CRDs)

MEASUREMENT SCHEDULE:
- Every 200 lines of implementation code added
- Before committing significant changes
- When implementation feels "substantial"
- If approaching complexity milestones

RESPONSE TO SIZE WARNINGS:
- At 600 lines: Start planning completion strategy
- At 700 lines: Consider code optimization
- At 750 lines: STOP and prepare for potential split
- At 800 lines: IMMEDIATE STOP - transition to MEASURE_SIZE
---

```bash
#!/bin/bash
# COMPREHENSIVE SIZE MONITORING WITH BASE BRANCH DETECTION

# Step 1: Find line counter tool
find_line_counter() {
    # Search up directory tree
    SEARCH_DIR=$(pwd)
    while [ "$SEARCH_DIR" != "/" ]; do
        if [ -f "$SEARCH_DIR/tools/line-counter.sh" ]; then
            echo "$SEARCH_DIR/tools/line-counter.sh"
            return 0
        fi
        SEARCH_DIR=$(dirname "$SEARCH_DIR")
    done
    
    # Fallback: search from home
    find /home -name "line-counter.sh" -path "*/tools/*" 2>/dev/null | head -1
}

# Step 2: Determine correct base branch
get_base_branch() {
    # Try to find orchestrator-state.json
    SEARCH_DIR=$(pwd)
    while [ "$SEARCH_DIR" != "/" ]; do
        if [ -f "$SEARCH_DIR/orchestrator-state.json" ]; then
            # Extract phase integration branch
            BASE=$(grep "current_phase_integration:" "$SEARCH_DIR/orchestrator-state.json" -A 2 | \
                   grep "branch:" | awk '{print $2}' | tr -d '"')
            if [ -n "$BASE" ]; then
                echo "$BASE"
                return 0
            fi
        fi
        SEARCH_DIR=$(dirname "$SEARCH_DIR")
    done
    
    # Fallback: look for phase pattern in current branch
    CURRENT=$(git branch --show-current)
    if [[ "$CURRENT" =~ phase([0-9]+)/ ]]; then
        echo "phase${BASH_REMATCH[1]}/integration"
    else
        echo "phase1/integration"  # Default fallback
    fi
}

# Step 3: Perform measurement
LC=$(find_line_counter)
if [ -z "$LC" ]; then
    echo "❌ FATAL: Cannot find line-counter.sh tool!"
    exit 1
fi

BASE=$(get_base_branch)
BRANCH=$(git branch --show-current)
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')

echo "📊 SIZE MEASUREMENT - $TIMESTAMP"
echo "Tool: $LC"
echo "Base: $BASE"
echo "Current: $BRANCH"

# Run measurement with proper parameters (R304 compliance)
RESULT=$("$LC" -b "$BASE" -c "$BRANCH")
echo "$RESULT"
LINES=$(echo "$RESULT" | grep "Total" | awk '{print $NF}')

if [ -z "$LINES" ]; then
    echo "⚠️ Could not extract line count, trying without base..."
    LINES=$("$LC" -c "$BRANCH" | grep "Total" | awk '{print $NF}')
fi

echo ""
echo "MEASUREMENT RESULT: $LINES/800 lines"

# Decision logic with R304 compliance
if [ "$LINES" -ge 800 ]; then
    echo "🚨🚨🚨 CRITICAL: SIZE LIMIT EXCEEDED ($LINES/800)"
    echo "ACTION: STOP IMMEDIATELY - AUTOMATIC FAILURE"
    exit 1
elif [ "$LINES" -ge 750 ]; then
    echo "🚨 DANGER: Approaching hard limit ($LINES/800)"
    echo "ACTION: Complete immediately or prepare split"
elif [ "$LINES" -ge 700 ]; then
    echo "⚠️ WARNING: High line count ($LINES/800)"
    echo "ACTION: Optimize code or plan completion"
elif [ "$LINES" -ge 600 ]; then
    echo "📊 CAUTION: Substantial size ($LINES/800)"
    echo "ACTION: Monitor closely, measure every hour"
else
    echo "✅ COMPLIANT: Safe to continue ($LINES/800)"
fi

# Log measurement with base branch for traceability
echo "- [$TIMESTAMP] Size: $LINES/800 (base: $BASE)" >> work-log.md
```

## 🔴🔴🔴 CRITICAL: CONTINUOUS SIZE MONITORING REQUIREMENTS 🔴🔴🔴

### Mandatory Monitoring Schedule
**YOU MUST MEASURE SIZE AT THESE INTERVALS OR FACE GRADING PENALTIES:**

1. **STARTUP**: Establish baseline when entering IMPLEMENTATION state
2. **EVERY 100 LINES**: After adding approximately 100 lines
3. **EVERY 60 MINUTES**: During active coding sessions
4. **BEFORE COMMITS**: Verify size before any git commit
5. **AFTER FEATURES**: When completing any component/feature

### R304 Compliance Requirements
**CRITICAL: Code Reviewer expects these EXACT standards:**

```bash
# NEVER DO THIS (AUTOMATIC -100% FAILURE):
wc -l *.go                                    # Manual counting = FAIL

# ALWAYS DO THIS (tool auto-detects correct base branch):
./tools/line-counter.sh   # No parameters needed - auto-detects everything!
```

### Base Branch Determination Rules
- **For Efforts**: Use phase integration branch (e.g., phase1/integration)
- **For Splits**: Use original effort branch (before split)
- **For Fixes**: Use current integration branch
- **NEVER**: Use "main" or "master" as base

### Self-Monitoring Enforcement
```bash
# Add to your implementation workflow:
LAST_MEASUREMENT=$(date +%s)

check_if_measurement_needed() {
    CURRENT_TIME=$(date +%s)
    TIME_SINCE_LAST=$((CURRENT_TIME - LAST_MEASUREMENT))
    
    # Measure if more than 1 hour
    if [ $TIME_SINCE_LAST -gt 3600 ]; then
        echo "⏰ Time trigger: Measuring size (>1 hour since last)"
        return 0
    fi
    
    # Check if significant changes made
    LINES_CHANGED=$(git diff --numstat | awk '{sum+=$1+$2} END {print sum}')
    if [ "$LINES_CHANGED" -gt 100 ]; then
        echo "📝 Change trigger: Measuring size (>100 lines changed)"
        return 0
    fi
    
    return 1
}

# Call this regularly during implementation
if check_if_measurement_needed; then
    # Run the comprehensive monitoring script above
    LAST_MEASUREMENT=$(date +%s)
fi
```

## Test-Driven Implementation

---
### ℹ️ RULE R060.0.0 - Test Implementation
**Source:** rule-library/RULE-REGISTRY.md#R060
**Criticality:** INFO - Best practice

TEST IMPLEMENTATION REQUIREMENTS:

TESTING STRATEGY:
1. Write unit tests alongside implementation
2. Test coverage must meet minimum requirements
3. Integration tests for cross-component interactions
4. Edge case and error condition testing

COVERAGE REQUIREMENTS:
- Core business logic: 90%+ coverage
- Controller methods: 85%+ coverage
- Utility functions: 80%+ coverage
- Integration points: 85%+ coverage

TEST STRUCTURE:
- Unit tests in same package as implementation
- Integration tests in separate test package
- Test files follow _test.go convention
- Table-driven tests for multiple scenarios
---

```go
// Example test structure for implementation
package api_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestResourceController_Reconcile(t *testing.T) {
    tests := []struct {
        name           string
        inputResource  *Resource
        expectedResult reconcile.Result
        expectedError  bool
        setupMocks     func(*MockClient)
    }{
        {
            name: "successful_reconciliation",
            inputResource: &Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "test-resource",
                    Namespace: "default",
                },
                Spec: ResourceSpec{
                    Replicas: 3,
                },
            },
            expectedResult: reconcile.Result{},
            expectedError:  false,
            setupMocks: func(client *MockClient) {
                client.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
                client.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
            },
        },
        {
            name: "resource_not_found_creates_new",
            inputResource: &Resource{
                ObjectMeta: metav1.ObjectMeta{
                    Name:      "new-resource",
                    Namespace: "default",
                },
            },
            expectedResult: reconcile.Result{RequeueAfter: time.Minute * 5},
            expectedError:  false,
            setupMocks: func(client *MockClient) {
                client.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(apierrors.NewNotFound(schema.GroupResource{}, "new-resource"))
                client.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            mockClient := NewMockClient(t)
            tt.setupMocks(mockClient)
            
            controller := &ResourceController{
                Client: mockClient,
                Log:    logr.Discard(),
            }

            // Execute
            result, err := controller.Reconcile(context.TODO(), ctrl.Request{
                NamespacedName: types.NamespacedName{
                    Name:      tt.inputResource.Name,
                    Namespace: tt.inputResource.Namespace,
                },
            })

            // Assert
            if tt.expectedError {
                require.Error(t, err)
            } else {
                require.NoError(t, err)
            }
            assert.Equal(t, tt.expectedResult, result)
        })
    }
}

// Integration test example
func TestResourceController_Integration(t *testing.T) {
    testEnv := envtest.Environment{
        CRDDirectoryPaths: []string{
            filepath.Join("..", "config", "crd", "bases"),
        },
    }

    cfg, err := testEnv.Start()
    require.NoError(t, err)
    defer testEnv.Stop()

    // Test real reconciliation with real API server
    // ...
}
```

## Work Log Creation and Maintenance (R301)

### 🚨🚨 MANDATORY: Create Timestamped Work Log First (R301)
**Prevents naming collisions when multiple agents work simultaneously**

```bash
# R301: MANDATORY - Create work log with unique timestamp
EFFORT_NAME=$(basename "$(pwd)")
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
WORKLOG_NAME="worklog-${EFFORT_NAME}-${TIMESTAMP}.txt"

echo "📝 Creating timestamped work log per R301: $WORKLOG_NAME"

cat > "$WORKLOG_NAME" << EOF
# Work Log - ${EFFORT_NAME}
Created: $(date)
Branch: $(git branch --show-current)
Plan: ${PLAN_FILE:-IMPLEMENTATION-PLAN.md}

## Scope from Plan:
$(grep -A 5 "EXACTLY\|DO NOT" "${PLAN_FILE:-IMPLEMENTATION-PLAN.md}" | head -15 2>/dev/null)

## Session Log:
### [$(date '+%H:%M')] Implementation Start
- Verified scope boundaries
- Starting implementation
EOF

echo "✅ Work log created with timestamp to prevent collisions"

# Create convenience symlink (optional)
ln -sf "$WORKLOG_NAME" "work-log-current.md"
```

---
### ℹ️ RULE R018.0.0 - Progress Reporting  
**Source:** rule-library/RULE-REGISTRY.md#R018
**Criticality:** INFO - Best practice

WORK LOG UPDATE REQUIREMENTS:

UPDATE FREQUENCY:
- Every 30-60 minutes during active implementation
- After completing each implementation plan section
- After each significant commit
- When encountering and resolving issues
- Before any size measurements

LOG ENTRY CONTENT:
- Timestamp of work session
- Specific tasks completed
- Current line count and size status
- Test coverage achieved
- Issues encountered and resolutions
- Next planned activities
---

---
### 🚨🚨 RULE R170.0.0 - Work Log Template Usage
**Source:** rule-library/RULE-REGISTRY.md#R170
**Criticality:** MANDATORY - Required for approval

MANDATORY WORK LOG USAGE WITH R301 NAMING:

On Starting Implementation:
1. CREATE timestamped work log (R301 MANDATORY):
```bash
# ALWAYS use timestamp to prevent collisions
EFFORT=$(basename $(pwd))
WORKLOG="worklog-${EFFORT}-$(date +%Y%m%d-%H%M%S).txt"
echo "Creating: $WORKLOG"

# Copy template if available
if [ -f "$CLAUDE_PROJECT_DIR/templates/WORK-LOG-TEMPLATE.md" ]; then
    cp "$CLAUDE_PROJECT_DIR/templates/WORK-LOG-TEMPLATE.md" "$WORKLOG"
else
    echo "# Work Log - $EFFORT" > "$WORKLOG"
    echo "Date: $(date)" >> "$WORKLOG"
fi
```

2. NEVER use generic "work-log.md" name (collision risk)

3. Initialize work log:
- Fill in effort ID and name
- Set target metrics from implementation plan
- Record start date/time

Daily Updates (from template):
- Use "Daily Log" section format
- Update "Progress Summary" table
- Fill "Implementation Checkpoints" at 25/50/75/100%
- Record in "Size Tracking" table
- Document in "Issues and Resolutions" section

Completion Requirements:
- "Completion Checklist" fully checked
- "Final Metrics" section filled
- "Review Feedback" section ready
---

```markdown
# Work Log Template Entry

## [2025-08-23 14:30] Implementation Session
**Duration**: 1.5 hours
**Focus**: Controller reconciliation logic

### Completed Tasks
- ✅ Implemented ResourceController struct and basic scaffolding
- ✅ Added reconciliation logic for CREATE operations  
- ✅ Implemented status update handling
- ✅ Added error handling and retry logic

### Implementation Progress
- **Lines Added**: ~150 lines (total: 387/800)
- **Files Modified**: 
  - `pkg/controllers/resource_controller.go` (new, 298 lines)
  - `pkg/api/v1/resource_types.go` (modified, +89 lines)
- **Test Coverage**: 87% (unit tests for reconcile logic)

### Quality Metrics
- Size Check: ✅ 387/800 lines (48% of limit)
- Tests: ✅ 15 unit tests written and passing
- Linting: ✅ No linting issues
- Build: ✅ Clean build

### Issues Encountered
1. **Issue**: Client injection pattern unclear from implementation plan
   - **Resolution**: Reviewed existing controller patterns in codebase
   - **Time**: 20 minutes investigation + 15 minutes implementation

2. **Issue**: Status update causing infinite reconciliation loop
   - **Resolution**: Added status comparison before updating
   - **Time**: 30 minutes debugging + 10 minutes fix

### Next Session Plans
- [ ] Implement UPDATE operation reconciliation
- [ ] Add DELETE operation handling with finalizers
- [ ] Implement resource status conditions
- [ ] Add integration tests

### Notes
- Controller pattern consistent with existing codebase conventions
- May need to add webhook validation in later session
- Consider performance optimization if resource count grows large
```

## Implementation Plan Adherence

---
### 🚨🚨 RULE R054.0.0 - Implementation Plan Creation
**Source:** rule-library/RULE-REGISTRY.md#R054
**Criticality:** MANDATORY - Required for approval

PLAN ADHERENCE REQUIREMENTS:

MANDATORY COMPLIANCE:
- Follow implementation plan section order exactly
- Complete each planned task fully before moving to next
- Document any deviations with rationale
- Update plan if discoveries require scope changes

SCOPE CHANGE PROTOCOL:
1. Identify scope change need during implementation
2. Document current progress and stopping point
3. Calculate impact on size estimates
4. Request code reviewer guidance for significant changes
5. Update implementation plan before proceeding
---

```python
def validate_implementation_plan_adherence(work_dir):
    """Validate current implementation follows the plan"""
    
    plan_path = os.path.join(work_dir, 'IMPLEMENTATION-PLAN.md')
    worklog_path = os.path.join(work_dir, 'work-log.md')
    
    if not os.path.exists(plan_path):
        return {
            'valid': False,
            'error': 'Implementation plan not found',
            'action_required': 'Locate or create implementation plan'
        }
    
    plan_tasks = parse_implementation_plan_tasks(plan_path)
    completed_tasks = parse_work_log_completed_tasks(worklog_path)
    
    adherence_analysis = {
        'plan_tasks_total': len(plan_tasks),
        'tasks_completed': len(completed_tasks),
        'completion_percentage': (len(completed_tasks) / len(plan_tasks)) * 100 if plan_tasks else 0,
        'out_of_order_tasks': [],
        'scope_deviations': [],
        'missing_tasks': []
    }
    
    # Check for out-of-order completion
    planned_order = [task['order'] for task in plan_tasks]
    completed_order = [task['order'] for task in completed_tasks if task.get('order')]
    
    for i, completed_task_order in enumerate(completed_order[1:], 1):
        if completed_task_order < completed_order[i-1]:
            adherence_analysis['out_of_order_tasks'].append({
                'task_order': completed_task_order,
                'expected_after': completed_order[i-1]
            })
    
    # Check for scope deviations
    planned_task_names = {task['name'] for task in plan_tasks}
    completed_task_names = {task['name'] for task in completed_tasks}
    
    adherence_analysis['scope_deviations'] = list(completed_task_names - planned_task_names)
    adherence_analysis['missing_tasks'] = [
        task for task in plan_tasks 
        if task['name'] not in completed_task_names
    ]
    
    # Determine overall adherence status
    if adherence_analysis['scope_deviations'] or adherence_analysis['out_of_order_tasks']:
        adherence_status = 'DEVIATIONS_DETECTED'
    elif adherence_analysis['completion_percentage'] < 50:
        adherence_status = 'ON_TRACK_EARLY'
    elif adherence_analysis['completion_percentage'] < 90:
        adherence_status = 'ON_TRACK_PROGRESSING'
    else:
        adherence_status = 'ON_TRACK_NEAR_COMPLETION'
    
    return {
        'valid': len(adherence_analysis['scope_deviations']) == 0,
        'status': adherence_status,
        'analysis': adherence_analysis,
        'recommendations': generate_adherence_recommendations(adherence_analysis)
    }

def generate_adherence_recommendations(analysis):
    """Generate recommendations based on plan adherence analysis"""
    
    recommendations = []
    
    if analysis['out_of_order_tasks']:
        recommendations.append({
            'type': 'ORDER_CORRECTION',
            'message': 'Complete remaining tasks in planned order',
            'priority': 'MEDIUM'
        })
    
    if analysis['scope_deviations']:
        recommendations.append({
            'type': 'SCOPE_VALIDATION',
            'message': 'Document rationale for scope deviations or update plan',
            'priority': 'HIGH'
        })
    
    if analysis['completion_percentage'] > 80 and analysis['missing_tasks']:
        recommendations.append({
            'type': 'COMPLETION_FOCUS',
            'message': 'Focus on completing remaining planned tasks',
            'priority': 'HIGH'
        })
    
    return recommendations
```

## Code Quality Standards

---
### ℹ️ RULE R037.0.0 - Pattern Compliance
**Source:** rule-library/RULE-REGISTRY.md#R037
**Criticality:** INFO - Best practice

PATTERN COMPLIANCE REQUIREMENTS:

KUBERNETES PATTERNS:
- Controllers follow controller-runtime patterns
- Resources implement runtime.Object interface
- Status conditions follow standard conventions
- Finalizers implemented correctly for cleanup

KCP PATTERNS:
- Multi-tenancy implemented with logical clusters
- Authorization uses RBAC with cluster-scoped permissions
- API grouping follows KCP conventions
- Resource scheduling aware of workspace context

GO CODE PATTERNS:
- Error handling with wrapped errors
- Context propagation through call chains
- Structured logging with contextual fields
- Interface-based design for testability
---

```go
// Example of KCP-aware controller implementation
type ResourceController struct {
    client.Client
    Log    logr.Logger
    Scheme *runtime.Scheme
}

func (r *ResourceController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := r.Log.WithValues("resource", req.NamespacedName)
    log.V(1).Info("Starting reconciliation")

    // Get the resource
    var resource myapi.Resource
    if err := r.Get(ctx, req.NamespacedName, &resource); err != nil {
        if apierrors.IsNotFound(err) {
            // Resource deleted, cleanup if needed
            log.V(1).Info("Resource not found, assuming deleted")
            return ctrl.Result{}, nil
        }
        return ctrl.Result{}, fmt.Errorf("failed to get resource: %w", err)
    }

    // Check if resource is being deleted
    if !resource.DeletionTimestamp.IsZero() {
        return r.handleDeletion(ctx, &resource, log)
    }

    // Add finalizer if not present
    if !controllerutil.ContainsFinalizer(&resource, myapi.ResourceFinalizer) {
        controllerutil.AddFinalizer(&resource, myapi.ResourceFinalizer)
        return r.updateResource(ctx, &resource, log)
    }

    // Main reconciliation logic
    result, err := r.reconcileResource(ctx, &resource, log)
    if err != nil {
        // Update status with error condition
        r.updateStatusCondition(&resource, myapi.ConditionTypeReady, metav1.ConditionFalse, 
            myapi.ReasonReconcileError, err.Error())
        if statusErr := r.Status().Update(ctx, &resource); statusErr != nil {
            log.Error(statusErr, "Failed to update status")
        }
        return result, fmt.Errorf("reconciliation failed: %w", err)
    }

    // Update status with success condition
    r.updateStatusCondition(&resource, myapi.ConditionTypeReady, metav1.ConditionTrue,
        myapi.ReasonReconcileSuccess, "Resource reconciled successfully")
    
    if err := r.Status().Update(ctx, &resource); err != nil {
        return ctrl.Result{}, fmt.Errorf("failed to update status: %w", err)
    }

    log.V(1).Info("Reconciliation completed successfully")
    return result, nil
}

func (r *ResourceController) reconcileResource(ctx context.Context, resource *myapi.Resource, log logr.Logger) (ctrl.Result, error) {
    // KCP-aware logic: Handle multi-tenant scenarios
    clusterName := logicalcluster.From(resource)
    log = log.WithValues("cluster", clusterName)
    
    // Implement main business logic here
    // ...
    
    return ctrl.Result{}, nil
}

func (r *ResourceController) updateStatusCondition(resource *myapi.Resource, conditionType string, status metav1.ConditionStatus, reason, message string) {
    condition := metav1.Condition{
        Type:               conditionType,
        Status:             status,
        LastTransitionTime: metav1.Now(),
        Reason:             reason,
        Message:            message,
    }
    
    meta.SetStatusCondition(&resource.Status.Conditions, condition)
}
```

## Commit Strategy

---
### ℹ️ RULE R015.0.0 - Commit Message Format
**Source:** rule-library/RULE-REGISTRY.md#R015
**Criticality:** INFO - Best practice

COMMIT FREQUENCY AND FORMAT:

COMMIT FREQUENCY:
- Every 1-2 hours of active development
- After completing each implementation plan section
- Before any significant refactoring
- When reaching size measurement milestones

COMMIT MESSAGE FORMAT:
feat|fix|refactor|test: brief description (≤50 chars)

Optional body with details:
- What was implemented
- Size impact
- Test coverage added
- Any architectural decisions
---

```bash
# Example commit messages during implementation

# Feature commits
git commit -m "feat: add ResourceController reconciliation logic

- Implement basic CRUD reconciliation for Resource CRD
- Add error handling and retry mechanisms  
- Include status condition updates
- Size impact: +150 lines (387/800 total)
- Test coverage: 87% with 12 unit tests"

# Test commits  
git commit -m "test: add integration tests for ResourceController

- Add envtest-based integration test suite
- Test real reconciliation against k8s API server
- Cover edge cases and error conditions
- Size impact: +89 lines (476/800 total)
- Test coverage: 91% overall"

# Fix commits
git commit -m "fix: prevent infinite reconciliation loop in status updates

- Add status comparison before updating
- Only update status when actual changes detected
- Improves controller performance and reduces API calls
- No size impact (refactoring existing code)"
```

## 🔴🔴🔴 MANDATORY: Create IMPLEMENTATION-COMPLETE Marker 🔴🔴🔴

### CRITICAL REQUIREMENT - CANNOT CONSIDER WORK COMPLETE WITHOUT THIS

**When implementation is complete, you MUST create a completion marker:**

```bash
# MANDATORY when implementation is complete:
echo "🔴 Creating MANDATORY IMPLEMENTATION-COMPLETE marker..."
touch IMPLEMENTATION-COMPLETE.marker
cat > IMPLEMENTATION-COMPLETE.marker << EOF
Completed at: $(date '+%Y-%m-%d %H:%M:%S %Z')
Effort: $(basename $(pwd))
Branch: $(git branch --show-current)
Total lines: $(./tools/line-counter.sh | grep Total | awk '{print $NF}') lines
Final commit: $(git log --oneline -1)
Status: IMPLEMENTATION COMPLETE
EOF

# MUST add, commit and push
git add IMPLEMENTATION-COMPLETE.marker
git commit -m "marker: implementation complete - MANDATORY for orchestrator monitoring"
git push

echo "✅ IMPLEMENTATION-COMPLETE.marker created and pushed"
echo "📋 Orchestrator can now detect completion"
```

**THIS IS NOT OPTIONAL:**
- ❌ WITHOUT marker = Orchestrator cannot detect completion
- ❌ WITHOUT marker = Work is NOT considered complete
- ❌ WITHOUT marker = Grading penalty for incomplete work
- ✅ WITH marker = Clear signal work is done
- ✅ WITH marker = Orchestrator can proceed with next steps

### Validation Before Stopping:
```bash
# MANDATORY check before considering work complete
if [ ! -f IMPLEMENTATION-COMPLETE.marker ]; then
    echo "🔴 ERROR: Cannot stop without creating IMPLEMENTATION-COMPLETE.marker"
    echo "Implementation is NOT complete without marker!"
    exit 1
fi
echo "✅ Completion marker exists - implementation is complete"
```

## State Transitions

From IMPLEMENTATION state:
- **SIZE_LIMIT_APPROACHED** → MEASURE_SIZE (Check exact size and plan next steps)
- **IMPLEMENTATION_COMPLETE** → TEST_WRITING (Focus on comprehensive testing)
- **ISSUES_ENCOUNTERED** → FIX_ISSUES (Address blocking problems)
- **SPLIT_REQUIRED** → SPLIT_WORK (Effort too large, needs splitting)
- **PLAN_DEVIATION** → CODE_REVIEW (Get guidance on scope changes)

## Early Warning System

```python
def monitor_implementation_health():
    """Monitor implementation progress for early warning signs"""
    
    current_dir = os.getcwd()
    
    health_indicators = {
        'size_trajectory': analyze_size_growth_trend(),
        'test_coverage': measure_current_test_coverage(),
        'plan_adherence': check_plan_adherence_status(),
        'commit_frequency': analyze_commit_frequency(),
        'issue_resolution_time': track_issue_resolution()
    }
    
    warnings = []
    
    # Size trajectory warning
    if health_indicators['size_trajectory']['projected_final_size'] > 750:
        warnings.append({
            'type': 'SIZE_WARNING',
            'message': f"Projected final size {health_indicators['size_trajectory']['projected_final_size']} lines may exceed limit",
            'action': 'Consider scope reduction or prepare for split'
        })
    
    # Test coverage warning  
    if health_indicators['test_coverage']['current_percentage'] < 80:
        warnings.append({
            'type': 'COVERAGE_WARNING',
            'message': f"Test coverage at {health_indicators['test_coverage']['current_percentage']:.1f}% below target",
            'action': 'Increase test coverage before proceeding'
        })
    
    # Plan adherence warning
    if health_indicators['plan_adherence']['deviation_score'] > 20:
        warnings.append({
            'type': 'PLAN_WARNING',
            'message': 'Significant deviation from implementation plan detected',
            'action': 'Review plan adherence or update plan'
        })
    
    return {
        'health_status': 'HEALTHY' if not warnings else 'WARNINGS_PRESENT',
        'indicators': health_indicators,
        'warnings': warnings,
        'next_checkpoint': 'In 200 lines or 1 hour of work'
    }
