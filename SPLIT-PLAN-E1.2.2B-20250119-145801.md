# SPLIT-PLAN-E1.2.2B: fallback-recommendations

## Split 002 of 003: Recommendation Engine
**Planner**: Code Reviewer Agent
**Parent Effort**: E1.2.2 fallback-strategies
**Target Size**: ~550 lines
**Priority**: SECOND (can parallelize with E1.2.2C after E1.2.2A)

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries
- **Previous Split**: Split 001 (fallback-core) of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-core/
  - Branch: phase1/wave2/fallback-core
  - Summary: Core fallback logic and type definitions
- **This Split**: Split 002 of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-recommendations/
  - Branch: phase1/wave2/fallback-recommendations
- **Next Split**: Split 003 (fallback-security) of phase1/wave2/fallback-strategies
  - Path: efforts/phase1/wave2/fallback-security/
  - Branch: phase1/wave2/fallback-security

## Files in This Split (EXCLUSIVE)
```
pkg/fallback/recommendations/
├── recommendations.go       # 440 lines - Recommendation engine
├── recommendations_test.go  # ~110 lines - Unit tests
Total: ~550 lines
```

## Functionality Scope

### Core Components
1. **Recommendation Engine** (`recommendations.go`)
   - Failure pattern analysis
   - Recommendation generation algorithms
   - Priority scoring system
   - Context-aware suggestions
   - Historical data analysis

2. **Recommendation Types**
   - Registry alternatives
   - Retry strategy adjustments
   - Configuration changes
   - Manual intervention steps
   - Escalation procedures

## Implementation Instructions

### Step 1: Repository Setup
```bash
cd /home/vscode/workspaces/idpbuilder-oci-build-push/efforts/phase1/wave2
mkdir -p fallback-recommendations
cd fallback-recommendations
git init
git remote add origin [same as parent]
git checkout -b phase1/wave2/fallback-recommendations
```

### Step 2: Import Dependencies from E1.2.2A
```bash
# After E1.2.2A is complete and reviewed
mkdir -p pkg/fallback/recommendations
cd pkg/fallback/recommendations

# Import types from fallback-core
# Note: Will need to set up proper go.mod with replace directive
```

### Step 3: Implement Recommendation Engine

#### 3.1 Create `recommendations.go`
```go
package recommendations

import (
    "context"
    "sort"

    "github.com/idpbuilder/idpbuilder/pkg/fallback" // from E1.2.2A
)

// RecommendationEngine implements the RecommendationProvider interface
type RecommendationEngine struct {
    config           *Config
    analyzer         *FailureAnalyzer
    prioritizer      *Prioritizer
    historyTracker   *HistoryTracker
}

// Core functionality sections:
// 1. Failure Analysis (100 lines)
//    - Pattern detection
//    - Root cause analysis
//    - Failure categorization

// 2. Recommendation Generation (150 lines)
//    - Strategy alternatives
//    - Configuration adjustments
//    - Manual interventions

// 3. Prioritization (90 lines)
//    - Scoring algorithms
//    - Context weighting
//    - Success probability

// 4. History Tracking (60 lines)
//    - Previous recommendations
//    - Success rates
//    - Pattern learning

// 5. Public API (40 lines)
//    - GetRecommendations implementation
//    - Configuration methods
//    - Utility functions
```

#### 3.2 Key Features to Implement

1. **Pattern Detection**
   ```go
   func (e *RecommendationEngine) detectPattern(failures []FailureEvent) Pattern {
       // Analyze failure sequence
       // Identify common causes
       // Return detected pattern
   }
   ```

2. **Alternative Strategies**
   ```go
   func (e *RecommendationEngine) generateAlternatives(failure *FailureContext) []Strategy {
       // Generate registry alternatives
       // Suggest retry modifications
       // Propose configuration changes
   }
   ```

3. **Priority Scoring**
   ```go
   func (e *RecommendationEngine) scoreRecommendation(rec Recommendation, context Context) float64 {
       // Calculate success probability
       // Weight by impact
       // Consider implementation effort
   }
   ```

### Step 4: Write Comprehensive Tests
Create `recommendations_test.go`:
- Unit tests for each component
- Integration tests with mock fallback types
- Edge case coverage
- Performance benchmarks

### Step 5: Validation
```bash
# Ensure compilation
go build ./pkg/fallback/recommendations

# Run tests
go test ./pkg/fallback/recommendations -v -cover

# Measure size
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done
$PROJECT_ROOT/tools/line-counter.sh
```

## Dependencies

### Internal Dependencies
- E1.2.2A (fallback-core):
  - Import type definitions
  - Implement `RecommendationProvider` interface
  - Use `FailureContext` for analysis

### External Dependencies (go.mod)
```go
require (
    github.com/idpbuilder/idpbuilder v0.0.0 // local replace
    // ... existing project dependencies
)

replace github.com/idpbuilder/idpbuilder => ../fallback-core
```

## Integration Points

### From E1.2.2A (fallback-core)
- `RecommendationProvider` interface implementation
- `FailureContext` type usage
- `Recommendation` type definition

### To Main Fallback System
- Pluggable recommendation engine
- Configurable priority weights
- Historical learning capabilities

## Success Criteria

1. **Size Compliance**: Total lines < 800 (target ~550)
2. **Interface Implementation**: Fully implements `RecommendationProvider`
3. **Functionality**: All recommendation features working
4. **Tests**: >85% code coverage
5. **Performance**: Recommendations generated < 100ms
6. **No Stubs**: Complete implementation, no TODOs

## Risk Mitigation

- **Risk**: Recommendation logic too complex
- **Mitigation**: Focus on core patterns, defer advanced ML features

- **Risk**: Integration with E1.2.2A types
- **Mitigation**: Early validation of interface compatibility

## Review Focus Areas

1. Algorithm effectiveness
2. Performance characteristics
3. Test coverage quality
4. Interface compliance
5. Size limit adherence

## Parallel Development Notes

This split can be developed in parallel with E1.2.2C (fallback-security) once E1.2.2A is complete. Key considerations:
- No direct dependencies on E1.2.2C
- Share only through E1.2.2A interfaces
- Coordinate integration testing after both complete