# Implementation Plan for Effort 1.1.2: Command Skeleton

**Created**: 2025-09-22T23:08:13Z
**Location**: efforts/phase1/wave1/effort-1.1.2-command-skeleton
**Phase**: 1
**Wave**: 1
**Effort**: 1.1.2
**Title**: Implement Command Skeleton

## Effort Metadata

**Branch**: idpbuilderpush/phase1/wave1/command-skeleton
**Can Parallelize**: No (depends on 1.1.1 tests)
**Parallel With**: None
**Size Estimate**: 200 lines
**Dependencies**: Effort 1.1.1 (write-command-tests)

## TDD GREEN Phase Context

This effort implements the MINIMAL code to make tests from effort 1.1.1 pass. This is the GREEN phase of TDD.

**Key Principle**: Write ONLY enough code to make tests pass - no more, no less.

## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| None | - | - | No existing interfaces to implement |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| cobra.Command | pkg/cmd/root.go | Root command structure | Import and use cobra pattern |
| helpers | pkg/cmd/helpers | Logging and output helpers | Import for consistency |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| Execute | pkg/cmd/root.go | Execute(ctx context.Context) | Main entry point |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create alternative command execution pattern
- DO NOT reimplement cobra command structure
- DO NOT create new logging/output helpers

### REQUIRED INTEGRATIONS (R373)
- MUST follow existing cobra command pattern from pkg/cmd/
- MUST use similar structure to create/delete/get commands
- MUST integrate with main.go execution flow

## Test Requirements Analysis

Based on effort 1.1.1's `cmd/push/root_test.go`, the following MUST be implemented:

### Required Variables/Functions (from test expectations):
1. **pushCmd** - *cobra.Command variable (tested by all test functions)
2. **PushConfig** struct with fields:
   - RegistryURL string
   - Username string
   - Password string
   - Namespace string
   - Dir string
   - Insecure bool
   - PlainHTTP bool

### Required Cobra Command Properties:
- **Use**: "push"
- **Name()**: "push"
- **Short**: Non-empty description
- **Long**: Non-empty description
- **Args**: Custom validation function

### Required Flags:
| Flag | Type | Shorthand | Default | Required |
|------|------|-----------|---------|----------|
| username | string | u | "" | Yes |
| password | string | p | "" | Yes |
| namespace | string | n | "idpbuilder" | Yes |
| dir | string | d | "." | Yes |
| insecure | bool | - | false | Yes |
| plain-http | bool | - | false | Yes |

### Test Coverage Requirements:
- All 7 tests in root_test.go MUST pass
- Tests verify structure only, not functionality
- No actual push operations needed yet

## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:

**File: cmd/push/root.go** (~120 lines)
- Variable: `pushCmd *cobra.Command` (~60 lines)
  - Use, Short, Long fields
  - Args validation function inline
  - RunE stub function (returns nil)
- Function: `init()` (~40 lines)
  - Register all 6 flags with descriptions
  - Set default values
  - Define shorthands
- Function: `validateArgs(cmd *cobra.Command, args []string) error` (~20 lines)
  - Check for exactly one argument
  - Return appropriate errors

**File: cmd/push/config.go** (~30 lines)
- Type: `PushConfig struct` (~10 lines)
  - 7 fields matching test expectations
  - NO methods
- Function: `NewPushConfig() *PushConfig` (~10 lines)
  - Return config with defaults
- Function: `parseFlags(cmd *cobra.Command) (*PushConfig, error)` (~10 lines)
  - Extract flag values
  - Return populated config

**File: cmd/push/push.go** (~30 lines) [OPTIONAL - only if needed for organization]
- Function: `runPush(cmd *cobra.Command, args []string) error` (~30 lines)
  - Parse flags
  - Validate registry URL format
  - Return nil (stub implementation)

**TOTAL**: ~180 lines (well under 200 limit)

### DO NOT IMPLEMENT:
- ❌ Registry client creation
- ❌ Authentication logic
- ❌ OCI push operations
- ❌ Complex error handling
- ❌ Validation beyond test requirements
- ❌ Logging or debug output
- ❌ Configuration file parsing
- ❌ Environment variable handling (tests handle this)
- ❌ Any actual push functionality
- ❌ Methods on PushConfig struct

## R355 PRODUCTION READINESS - ZERO TOLERANCE

This implementation MUST be production-ready from the first commit:
- ❌ NO STUBS except RunE returning nil (allowed for TDD GREEN)
- ❌ NO MOCKS except in test directories
- ❌ NO hardcoded credentials or secrets
- ❌ NO static configuration values
- ❌ NO TODO/FIXME markers in code
- ❌ NO panic("not implemented") patterns
- ❌ NO fake or dummy data

**Note**: The RunE function returning nil is NOT a stub - it's the correct GREEN phase implementation that makes tests pass without implementing the full feature yet.

## Configuration Requirements (R355 Mandatory)

### CORRECT Implementation Examples:

```go
// ✅ Flag with default from constant
namespace := cmd.Flags().StringP("namespace", "n", DefaultNamespace, "Registry namespace")

// ✅ Validation with clear error
if len(args) != 1 {
    return fmt.Errorf("exactly one registry URL required, got %d", len(args))
}

// ✅ Config with defaults
func NewPushConfig() *PushConfig {
    return &PushConfig{
        Namespace: DefaultNamespace,
        Dir:       ".",
        Insecure:  false,
        PlainHTTP: false,
    }
}
```

## Implementation Steps

### Step 1: Import Tests from Effort 1.1.1
```bash
# Copy test file from effort 1.1.1
cp ../effort-1.1.1-write-command-tests/cmd/push/root_test.go cmd/push/

# Verify tests fail initially (TDD RED verification)
go test ./cmd/push/... -v
```

### Step 2: Create Command Structure
1. Create `cmd/push/` directory
2. Implement `root.go` with pushCmd variable
3. Add all required cobra command fields
4. Implement Args validation inline

### Step 3: Implement Configuration
1. Create `config.go` with PushConfig struct
2. Add NewPushConfig factory function
3. Implement parseFlags helper

### Step 4: Register Flags
1. Add init() function to root.go
2. Register all 6 flags with proper types
3. Set default values
4. Configure shorthands

### Step 5: Verify Tests Pass
```bash
# All 7 tests should now pass
go test ./cmd/push/... -v

# Verify test coverage
go test ./cmd/push/... -cover
```

### Step 6: Integration
1. Add push command to main root command
2. Update imports in pkg/cmd/root.go
3. Verify command appears in help

## Size Management

**Measurement Command**:
```bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure with official tool
$PROJECT_ROOT/tools/line-counter.sh
```

**Expected Size Breakdown**:
- cmd/push/root.go: ~120 lines
- cmd/push/config.go: ~30 lines
- cmd/push/push.go: ~30 lines (if needed)
- **Total**: ~180 lines

**Size Compliance**: ✅ Well under 200-line budget

## Atomic PR Design (R220)

```yaml
effort_atomic_pr_design:
  pr_summary: "feat: implement push command skeleton (TDD GREEN phase)"
  can_merge_to_main_alone: true  # Tests included, no dependencies

  r355_production_ready_checklist:
    no_hardcoded_values: true
    all_config_from_env: false  # Flags used instead
    no_stub_implementations: true  # RunE returns nil is valid GREEN
    no_todo_markers: true
    all_functions_complete: true

  feature_flags_needed: []  # Not needed - tests ensure compatibility

  pr_verification:
    tests_pass_alone: true
    build_remains_working: true
    flags_tested_both_ways: N/A
    no_external_dependencies: true
    backward_compatible: true

  example_pr_structure:
    files_added:
      - "cmd/push/root.go"
      - "cmd/push/config.go"
      - "cmd/push/root_test.go"  # From effort 1.1.1
    tests_included:
      - "7 unit tests for command structure"
      - "Flag registration tests"
      - "Argument validation tests"
    documentation:
      - "Command help text included"
```

## Success Criteria

1. ✅ All 7 tests from effort 1.1.1 pass
2. ✅ Total implementation ≤ 200 lines
3. ✅ Follows existing cobra patterns
4. ✅ No stubs or TODOs in code
5. ✅ Clean, minimal implementation
6. ✅ Command integrates with root

## Notes for SW Engineer

1. **Start with tests**: Copy root_test.go from effort 1.1.1 first
2. **Minimal implementation**: Resist adding extra features
3. **Follow patterns**: Look at create/delete/get commands for examples
4. **TDD discipline**: Make tests pass with minimum code
5. **No push logic**: This effort is structure only

## Size Limit Clarification (R359)

- The 200-line limit applies to NEW CODE YOU ADD
- Repository will grow by ~180-200 lines (EXPECTED)
- NEVER delete existing code to meet size limits
- Current repository has existing commands - we ADD to it

---

**Review Checkpoint**: Before implementing, verify:
- [ ] Tests copied from effort 1.1.1
- [ ] Tests fail initially (RED phase)
- [ ] Implementation scope is minimal
- [ ] No actual push logic included
- [ ] Following cobra patterns from existing commands