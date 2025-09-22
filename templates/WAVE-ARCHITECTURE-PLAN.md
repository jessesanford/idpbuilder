# Phase [PHASE_NUMBER] Wave [WAVE_NUMBER] Architecture Plan

## 🎯 Wave Vision Alignment

**Phase**: [PHASE_NUMBER] - [PHASE_NAME]  
**Wave**: [WAVE_NUMBER] - [WAVE_NAME]  
**Created By**: Architect Agent  
**Created After**: Wave [PREVIOUS_WAVE] Review PASSED  
**Date**: [DATE]  

### Phase Architecture Alignment
[How this wave contributes to the phase architecture goals]

### Wave Mission Statement
[1-2 sentences on what this wave achieves architecturally]

## 📊 Previous Wave Analysis

### Completed Implementations
```yaml
previous_waves:
  wave_1:
    - component: [COMPONENT_NAME]
      type: [API/Library/Service]
      location: pkg/phase[PHASE]/wave1/
      status: Complete
      reusable: true
      
  wave_2:
    - component: [COMPONENT_NAME]
      type: [API/Library/Service]
      location: pkg/phase[PHASE]/wave2/
      status: Complete
      reusable: true
```

### Available Contracts and APIs
```go
// Contracts from previous waves that this wave MUST use
import (
    wave1contracts "pkg/phase[PHASE]/wave1/api"
    wave2lib "pkg/phase[PHASE]/wave2/lib"
)

// Example of what's available
type AvailableService = wave1contracts.ServiceInterface
type AvailableClient = wave2lib.ClientImplementation
```

### Integration Requirements
- **From Wave [NUMBER]**: [WHAT_TO_INTEGRATE]
- **From Wave [NUMBER]**: [WHAT_TO_INTEGRATE]
- **Cross-Wave Shared Libraries**: [LIBRARIES]

## 🏗️ Wave [WAVE_NUMBER] Architecture

### Effort Dependency Graph
```
┌─────────────────────────────────────────────────────┐
│ CRITICAL PATH (Sequential - MUST be in order)      │
├─────────────────────────────────────────────────────┤
│ Effort 1: API Contracts & Interfaces               │
│    └─→ Defines ALL contracts for this wave         │
│    └─→ MUST complete first                         │
│                                                     │
│ Effort 2: Shared Libraries                         │
│    └─→ Depends on: Effort 1 contracts              │
│    └─→ Creates reusable implementations            │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ PARALLEL PATH (Can run simultaneously)             │
├─────────────────────────────────────────────────────┤
│ Effort 3: Feature A Implementation                 │
│    └─→ Depends on: Efforts 1 & 2                   │
│    └─→ Independent domain                          │
│                                                     │
│ Effort 4: Feature B Implementation                 │
│    └─→ Depends on: Efforts 1 & 2                   │
│    └─→ Independent domain                          │
│                                                     │
│ Effort 5: Feature C Implementation                 │
│    └─→ Depends on: Efforts 1 & 2                   │
│    └─→ Independent domain                          │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ INTEGRATION PATH (After parallel work)             │
├─────────────────────────────────────────────────────┤
│ Effort 6: Integration & Testing                    │
│    └─→ Depends on: Efforts 3, 4, 5                 │
│    └─→ Integrates all features                     │
│                                                     │
│ Effort 7: Documentation                            │
│    └─→ Depends on: Effort 6                        │
│    └─→ Documents entire wave                       │
└─────────────────────────────────────────────────────┘
```

### Effort Specifications

#### Effort 1: [NAME] - API Contracts & Interfaces
```go
// ============================================
// EFFORT 1: DEFINE ALL CONTRACTS FIRST
// This effort MUST be completed before ANY other work
// ============================================

package api

// Primary service contract for this wave
type [WaveService] interface {
    // Core operations
    [Operation1](ctx context.Context, req [Request]) ([Response], error)
    [Operation2](ctx context.Context, id string) ([Model], error)
    [Operation3](ctx context.Context, filter [Filter]) ([]*[Model], error)
}

// Data models
type [Request] struct {
    ID       string                 `json:"id"`
    Type     string                 `json:"type"`
    Payload  json.RawMessage        `json:"payload"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type [Response] struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// Repository contract
type [Repository] interface {
    Create(ctx context.Context, entity [Entity]) error
    Read(ctx context.Context, id string) (*[Entity], error)
    Update(ctx context.Context, id string, entity [Entity]) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, opts ListOptions) ([]*[Entity], error)
}
```

#### Effort 2: [NAME] - Shared Libraries
```go
// ============================================
// EFFORT 2: SHARED LIBRARY IMPLEMENTATIONS
// Depends on Effort 1 contracts
// ============================================

package lib

import (
    "pkg/phase[PHASE]/wave[WAVE]/effort1/api"
    "pkg/common/logger"
    "pkg/common/metrics"
)

// Shared client implementation
type Client struct {
    endpoint string
    logger   logger.Logger
    metrics  metrics.Collector
}

func NewClient(endpoint string) *Client {
    return &Client{
        endpoint: endpoint,
        logger:   logger.Get("wave[WAVE]"),
        metrics:  metrics.NewCollector("wave[WAVE]"),
    }
}

// Shared validation logic
func Validate[Request](req [Request]) error {
    // Common validation used by all efforts
}

// Shared error handling
func HandleError(err error) error {
    // Common error handling logic
}
```

#### Effort 3-5: [NAME] - Parallel Implementations
```go
// ============================================
// EFFORTS 3-5: PARALLEL FEATURE IMPLEMENTATIONS
// Can run simultaneously after Efforts 1 & 2
// ============================================

// Each effort works in its own domain
// No dependencies between efforts 3, 4, and 5

// Effort 3: Feature A
package feature_a
import "pkg/phase[PHASE]/wave[WAVE]/effort1/api"
import "pkg/phase[PHASE]/wave[WAVE]/effort2/lib"

// Effort 4: Feature B  
package feature_b
import "pkg/phase[PHASE]/wave[WAVE]/effort1/api"
import "pkg/phase[PHASE]/wave[WAVE]/effort2/lib"

// Effort 5: Feature C
package feature_c
import "pkg/phase[PHASE]/wave[WAVE]/effort1/api"
import "pkg/phase[PHASE]/wave[WAVE]/effort2/lib"
```

### Code Reuse Strategy

#### From Previous Phases
```go
// Phase 1 utilities (established foundation)
import (
    "pkg/phase1/common/auth"
    "pkg/phase1/common/config"
    "pkg/phase1/common/database"
)

// Usage in this wave
func NewService() *Service {
    return &Service{
        auth: auth.NewValidator(),        // Phase 1
        db:   database.NewConnection(),   // Phase 1
        cfg:  config.Load(),              // Phase 1
    }
}
```

#### From Previous Waves (Same Phase)
```go
// Earlier waves in this phase
import (
    wave1api "pkg/phase[PHASE]/wave1/api"
    wave2client "pkg/phase[PHASE]/wave2/client"
)

// Building on previous work
type EnhancedService struct {
    base wave1api.ServiceInterface    // Wave 1 contract
    client wave2client.Client          // Wave 2 implementation
}
```

#### Within This Wave
```yaml
reuse_map:
  effort_1_provides:
    - All API contracts
    - Data models
    - Interface definitions
    
  effort_2_provides:
    - Client library
    - Validation functions
    - Error handlers
    
  efforts_3_5_consume:
    - Effort 1 contracts (implement interfaces)
    - Effort 2 libraries (use utilities)
    
  effort_6_integrates:
    - All effort outputs
    - Creates unified service
```

## 📈 Implementation Priorities

### Implementation Order
```
Priority 1: APIs/Interfaces/Contracts (Effort 1)
    └─→ Unlocks all other work
    └─→ MUST be 100% complete first
    
Priority 2: Shared Libraries/Utilities (Effort 2)
    └─→ Reduces duplication
    └─→ Ensures consistency
    
Priority 3: Core Features (Efforts 3-5)
    └─→ Can parallelize
    └─→ Independent domains
    
Priority 4: Integration Layer (Effort 6)
    └─→ Combines all features
    └─→ End-to-end testing
    
Priority 5: Documentation (Effort 7)
    └─→ API docs
    └─→ Usage examples
```

### Critical Path Sequence
```
Step 1:   Effort 1 (Contracts)      ████████ BLOCKING
Step 2:   Effort 2 (Libraries)      ████████████████
Step 3:   Efforts 3-5 (Parallel)    ████████████████████████
Step 4:   Effort 6 (Integration)    ████████
Step 5:   Effort 7 (Documentation)  ████████
```

## 🔧 Technical Specifications

### Effort Size Constraints
```yaml
effort_sizes:
  effort_1_contracts:
    target: 400 lines
    maximum: 600 lines
    includes: interfaces, models, constants
    
  effort_2_libraries:
    target: 600 lines
    maximum: 750 lines
    includes: shared implementations
    
  efforts_3_5_features:
    target: 700 lines each
    maximum: 800 lines each
    includes: feature implementation
    
  effort_6_integration:
    target: 500 lines
    maximum: 700 lines
    includes: integration logic
    
  effort_7_docs:
    target: N/A
    maximum: N/A
    includes: markdown documentation
```

### Performance Targets (Per Effort)
```yaml
effort_performance:
  api_latency:
    p50: [NUMBER]ms
    p95: [NUMBER]ms
    p99: [NUMBER]ms
    
  throughput:
    minimum: [NUMBER] ops/sec
    target: [NUMBER] ops/sec
    
  resource_usage:
    memory: <[NUMBER]MB
    cpu: <[NUMBER]%
```

### Testing Requirements
```yaml
test_coverage:
  effort_1: N/A (interfaces only)
  effort_2: ≥80% (library code)
  effort_3: ≥75% (feature code)
  effort_4: ≥75% (feature code)
  effort_5: ≥75% (feature code)
  effort_6: ≥85% (integration)
  effort_7: N/A (documentation)
```

## ✅ Wave Success Criteria

### Architectural Success
- [ ] Contracts defined in Effort 1 before any implementation
- [ ] Parallel efforts (3-5) don't block each other
- [ ] All efforts successfully reuse previous code
- [ ] No circular dependencies between efforts
- [ ] Clean separation of concerns

### Implementation Success
- [ ] All efforts under 800 lines
- [ ] Parallelization achieved as designed
- [ ] Integration completed successfully
- [ ] Test coverage meets requirements
- [ ] Documentation complete

### Quality Gates
- [ ] Code review passed for each effort
- [ ] No architectural drift from wave plan
- [ ] Performance targets met
- [ ] Security requirements satisfied
- [ ] All APIs backward compatible

## 🚨 Risk Mitigation

### Risk: Effort 1 Incomplete Contracts
- **Impact**: Blocks ALL other efforts
- **Mitigation**: 
  - Architect reviews Effort 1 BEFORE implementation starts
  - Mock implementations provided
  - Contract changes require architect approval

### Risk: Parallel Effort Conflicts
- **Impact**: Integration failures
- **Mitigation**:
  - Clear domain boundaries defined
  - Separate package structures
  - Integration tests run continuously

### Risk: Size Limit Exceeded
- **Impact**: Requires splitting, delays wave
- **Mitigation**:
  - Monitor size continuously
  - Pre-plan splitting strategy
  - Keep buffer (target 700, max 800)

## 📚 Appendix

### File Structure Expected
```
pkg/phase[PHASE]/wave[WAVE]/
├── effort1/
│   ├── api/
│   │   ├── interfaces.go
│   │   ├── models.go
│   │   └── constants.go
│   └── README.md
├── effort2/
│   ├── lib/
│   │   ├── client.go
│   │   ├── validator.go
│   │   └── errors.go
│   └── lib_test.go
├── effort3/
│   ├── feature_a.go
│   └── feature_a_test.go
├── effort4/
│   ├── feature_b.go
│   └── feature_b_test.go
├── effort5/
│   ├── feature_c.go
│   └── feature_c_test.go
├── effort6/
│   ├── integration.go
│   └── integration_test.go
└── effort7/
    └── README.md
```

### References
- [Phase [PHASE] Architecture Plan](../PHASE-[PHASE]-ARCHITECTURE-PLAN.md)
- [Previous Wave Architecture](./PHASE-[PHASE]-WAVE-[PREV]-ARCHITECTURE-PLAN.md)
- [Master Implementation Plan](../../IMPLEMENTATION-PLAN.md)

### Document History
| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | [DATE] | Architect Agent | Initial wave architecture |