# SPLIT-PLAN-E1.2.2A: fallback-core

## Split 001 of 003: Core Fallback Infrastructure
**Planner**: Code Reviewer Agent
**Parent Effort**: E1.2.2 fallback-strategies
**Target Size**: ~650 lines
**Priority**: FIRST (foundational)

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries
- **Previous Split**: None (first split of fallback-strategies effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-core/
  - Branch: phase1/wave2/fallback-core
- **Next Splits**:
  - Split 002 (fallback-recommendations): efforts/phase1/wave2/fallback-recommendations/
  - Split 003 (fallback-security): efforts/phase1/wave2/fallback-security/

## Files in This Split (EXCLUSIVE)
```
pkg/fallback/
├── fallback.go         # 426 lines - Core fallback logic
├── types.go            # ~100 lines - Extracted type definitions
├── interfaces.go       # ~50 lines - Shared interface contracts
└── fallback_test.go    # ~74 lines - Unit tests
Total: ~650 lines
```

## Functionality Scope

### Core Components
1. **Fallback Manager** (`fallback.go`)
   - Primary fallback orchestration
   - Retry logic with exponential backoff
   - Circuit breaker implementation
   - Registry fallback strategies
   - Error categorization and handling

2. **Type System** (`types.go`)
   - `FallbackStrategy` struct
   - `FallbackConfig` configuration
   - `RetryPolicy` definitions
   - `ErrorCategory` enumeration
   - `FallbackResult` response types

3. **Interface Contracts** (`interfaces.go`)
   - `Fallbacker` interface for implementations
   - `RecommendationProvider` interface (for E1.2.2B)
   - `SecurityLogger` interface (for E1.2.2C)
   - `RegistryClient` interface adaptations

## Implementation Instructions

### Step 1: Repository Setup
```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2
mkdir -p fallback-core
cd fallback-core
git init
git remote add origin [same as parent]
git checkout -b phase1/wave2/fallback-core
```

### Step 2: Create Package Structure
```bash
mkdir -p pkg/fallback
cd pkg/fallback
```

### Step 3: Implement Core Files

#### 3.1 Create `types.go`
Extract all type definitions from the original implementation:
- Move struct definitions
- Move const/var declarations
- Move error definitions

#### 3.2 Create `interfaces.go`
Define clear contracts for other splits to implement:
```go
package fallback

// Fallbacker defines core fallback behavior
type Fallbacker interface {
    Execute(ctx context.Context, op Operation) (*FallbackResult, error)
    Configure(config *FallbackConfig) error
}

// RecommendationProvider will be implemented by E1.2.2B
type RecommendationProvider interface {
    GetRecommendations(failure *FailureContext) []Recommendation
}

// SecurityLogger will be implemented by E1.2.2C
type SecurityLogger interface {
    LogSecurityEvent(event *SecurityEvent) error
    EnableInsecureMode(reason string) error
}
```

#### 3.3 Implement `fallback.go`
Core fallback logic with:
- Manager struct implementation
- Retry mechanisms
- Circuit breaker logic
- Error handling
- Registry interaction fallbacks

#### 3.4 Write Tests
Create comprehensive tests in `fallback_test.go`:
- Unit tests for each fallback strategy
- Integration tests for manager
- Error scenario coverage
- Mock implementations for interfaces

### Step 4: Validation
```bash
# Ensure compilation
go build ./pkg/fallback

# Run tests
go test ./pkg/fallback -v

# Measure size (must be <800 lines)
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh
```

## Dependencies

### Internal Dependencies
- None (foundational split)

### External Dependencies (go.mod)
```go
require (
    github.com/spf13/cobra v1.7.0
    github.com/sirupsen/logrus v1.9.0
    k8s.io/apimachinery v0.27.3
    // ... existing project dependencies
)
```

## Integration Points

### Exports for E1.2.2B (fallback-recommendations)
- Type definitions from `types.go`
- `RecommendationProvider` interface
- `FailureContext` for analysis

### Exports for E1.2.2C (fallback-security)
- Type definitions from `types.go`
- `SecurityLogger` interface
- `SecurityEvent` structure

## Success Criteria

1. **Size Compliance**: Total lines < 800 (target ~650)
2. **Compilation**: Clean build with no errors
3. **Tests**: All unit tests passing
4. **Interfaces**: Clear contracts for other splits
5. **Documentation**: Inline comments for public APIs
6. **No Stubs**: All core functionality implemented

## Risk Mitigation

- **Risk**: Original `fallback.go` too large
- **Mitigation**: Extract more into types.go or defer to other splits

- **Risk**: Interface design blocks other splits
- **Mitigation**: Early review of interfaces before full implementation

## Review Focus Areas

1. Interface design adequacy
2. Type system completeness
3. Error handling robustness
4. Test coverage (>80%)
5. Size limit compliance