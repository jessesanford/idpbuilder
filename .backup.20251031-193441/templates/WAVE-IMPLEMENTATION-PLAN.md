# Phase [PHASE_NUMBER] Wave [WAVE_NUMBER] Implementation Plan

<!-- NOTE: [PROJECT_PREFIX/] should be replaced with the actual project prefix from target-repo-config.yaml 
     If project_prefix is "tmc-workspace", branches become: tmc-workspace/phase1/wave1/effort-name
     If project_prefix is empty, branches become: phase1/wave1/effort-name -->

## 📋 Wave Implementation Overview

**Phase**: [PHASE_NUMBER] - [PHASE_NAME]  
**Wave**: [WAVE_NUMBER] - [WAVE_NAME]  
**Created By**: Code Reviewer Agent  
**Based On**: PHASE-[PHASE]-WAVE-[WAVE]-ARCHITECTURE-PLAN.md  
**Date**: [DATE]  
**Total Efforts**: [NUMBER]  
**Estimated Total Lines**: [NUMBER]  

### Wave Implementation Mission
[How this implementation realizes the architectural vision from the Architecture Plan]

### Architecture Compliance Statement
✅ This implementation plan follows the Architecture Plan's:
- [ ] API contracts and interfaces
- [ ] Parallelization strategy  
- [ ] Code reuse requirements
- [ ] Priority ordering

## 📊 Implementation Context

### Required Inputs from Architecture
```yaml
from_architecture_plan:
  contracts_location: ./PHASE-[PHASE]-WAVE-[WAVE]-ARCHITECTURE-PLAN.md#contracts
  parallelization: ./PHASE-[PHASE]-WAVE-[WAVE]-ARCHITECTURE-PLAN.md#parallelization
  reuse_strategy: ./PHASE-[PHASE]-WAVE-[WAVE]-ARCHITECTURE-PLAN.md#reuse
```

### Available Resources from Previous Work
```yaml
previous_phases:
  phase_1:
    - pkg/phase1/common/auth/validator.go
    - pkg/phase1/common/config/loader.go
    - pkg/phase1/api/models.go
    
previous_waves:
  wave_1:
    - pkg/phase[PHASE]/wave1/api/service.go
    - pkg/phase[PHASE]/wave1/client/client.go
  wave_2:
    - pkg/phase[PHASE]/wave2/lib/utils.go
```

## 🚀 Parallelization Strategy

### Execution Groups
```yaml
group_1_sequential:  # MUST COMPLETE FIRST
  - effort_1: Contracts & Interfaces (blocking)
  
group_2_sequential:  # MUST COMPLETE SECOND
  - effort_2: Shared Libraries (blocking)
  
group_3_parallel:    # CAN RUN SIMULTANEOUSLY
  - effort_3: Feature A (independent)
  - effort_4: Feature B (independent)
  - effort_5: Feature C (independent)
  
group_4_sequential:  # MUST WAIT FOR GROUP 3
  - effort_6: Integration (requires all features)
  
group_5_sequential:  # FINAL
  - effort_7: Testing & Documentation
```

### Parallelization Rules
1. **Efforts 1-2**: MUST run sequentially - they block everything
2. **Efforts 3-5**: CAN run in parallel - assign to different engineers
3. **Effort 6**: MUST wait for ALL parallel efforts to complete
4. **Effort 7**: Final effort after integration

### Orchestrator Spawning Strategy
```bash
# Step 1: Spawn single Code Reviewer for Effort 1
spawn_code_reviewer effort_1

# Step 2: After Effort 1 complete, spawn for Effort 2
spawn_sw_engineer effort_2

# Step 3: After Effort 2 complete, spawn THREE engineers in PARALLEL
spawn_sw_engineer effort_3 &
spawn_sw_engineer effort_4 &
spawn_sw_engineer effort_5 &

# Step 4: After ALL parallel efforts complete, spawn for integration
spawn_sw_engineer effort_6

# Step 5: Final testing
spawn_sw_engineer effort_7
```

## 🎯 Effort Implementation Details

### Effort 1: [EFFORT_NAME] - Contracts & Interfaces
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-contracts`  
**Can Parallelize**: No (MUST BE FIRST - blocks all other efforts)  
**Parallel With**: None  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: None (MUST BE FIRST!)  

#### Files to Create
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/api/interfaces.go
    lines: ~150
    purpose: Define all service interfaces
    content: |
      // Main service interface
      type Service interface {
          Method1(ctx context.Context) error
          Method2(ctx context.Context) error
      }
      
  - path: pkg/phase[PHASE]/wave[WAVE]/api/models.go
    lines: ~200
    purpose: Define all data models
    content: |
      // Core data structures
      type Request struct {
          ID   string
          Data interface{}
      }
      
  - path: pkg/phase[PHASE]/wave[WAVE]/api/errors.go
    lines: ~50
    purpose: Define error types
    
  - path: pkg/phase[PHASE]/wave[WAVE]/api/constants.go
    lines: ~30
    purpose: Define constants and enums
```

#### Implementation Instructions
```markdown
1. Create api/ directory structure
2. Define ALL interfaces from architecture plan
3. Create data models with proper JSON tags
4. Add validation tags where needed
5. Create mock implementations for testing
6. Document each interface method
```

#### Test Requirements
- Unit tests for model validation
- Mock implementation tests
- Interface compliance tests

---

### Effort 2: [EFFORT_NAME] - Shared Libraries
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-libraries`  
**Can Parallelize**: No (blocks feature efforts 3-5)  
**Parallel With**: None  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: Effort 1 (contracts)  

#### Files to Create
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/lib/client.go
    lines: ~300
    purpose: Implement client library
    imports:
      - pkg/phase[PHASE]/wave[WAVE]/api
      - pkg/phase1/common/http
      
  - path: pkg/phase[PHASE]/wave[WAVE]/lib/validator.go
    lines: ~150
    purpose: Validation utilities
    imports:
      - pkg/phase[PHASE]/wave[WAVE]/api
      
  - path: pkg/phase[PHASE]/wave[WAVE]/lib/converter.go
    lines: ~100
    purpose: Data conversion utilities
```

#### Files to Reuse/Import
```yaml
reuse_from_previous:
  - source: pkg/phase1/common/logger/logger.go
    usage: Import for logging
    
  - source: pkg/phase1/common/metrics/collector.go
    usage: Import for metrics collection
    
  - source: pkg/phase[PHASE]/wave1/lib/base.go
    usage: Extend base functionality
```

#### Implementation Instructions
```markdown
1. Import Effort 1 contracts
2. Implement client with retry logic
3. Add comprehensive validation
4. Include error handling
5. Add metrics collection
6. Implement caching where appropriate
```

---

### Efforts 3-5: [EFFORT_NAMES] - Parallel Features

#### Effort 3: Feature A Implementation
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-feature-a`  
**Can Parallelize**: Yes  
**Parallel With**: Efforts 4, 5 (Feature B, Feature C)  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: Efforts 1 & 2 (contracts & libraries)
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_a/handler.go
    lines: ~400
    purpose: Implement feature A logic
    imports:
      - pkg/phase[PHASE]/wave[WAVE]/api
      - pkg/phase[PHASE]/wave[WAVE]/lib
      
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_a/processor.go
    lines: ~200
    purpose: Process feature A data
    
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_a/repository.go
    lines: ~150
    purpose: Data access for feature A
```

#### Effort 4: Feature B Implementation
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-feature-b`  
**Can Parallelize**: Yes  
**Parallel With**: Efforts 3, 5 (Feature A, Feature C)  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: Efforts 1 & 2 (contracts & libraries)
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_b/service.go
    lines: ~350
    purpose: Feature B service implementation
    imports:
      - pkg/phase[PHASE]/wave[WAVE]/api
      - pkg/phase[PHASE]/wave[WAVE]/lib
      
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_b/worker.go
    lines: ~250
    purpose: Background processing
    
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_b/cache.go
    lines: ~150
    purpose: Caching layer
```

#### Effort 5: Feature C Implementation
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-feature-c`  
**Can Parallelize**: Yes  
**Parallel With**: Efforts 3, 4 (Feature A, Feature B)  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: Efforts 1 & 2 (contracts & libraries)
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_c/controller.go
    lines: ~300
    purpose: Feature C controller
    imports:
      - pkg/phase[PHASE]/wave[WAVE]/api
      - pkg/phase[PHASE]/wave[WAVE]/lib
      
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_c/middleware.go
    lines: ~200
    purpose: Request processing middleware
    
  - path: pkg/phase[PHASE]/wave[WAVE]/features/feature_c/validator.go
    lines: ~200
    purpose: Feature-specific validation
```

#### Parallel Execution Instructions
```markdown
## CAN BE EXECUTED SIMULTANEOUSLY BY DIFFERENT SW ENGINEERS

### Feature A Engineer:
1. Work in features/feature_a/ directory only
2. Import contracts from Effort 1
3. Use libraries from Effort 2
4. No dependencies on Features B or C

### Feature B Engineer:
1. Work in features/feature_b/ directory only
2. Import contracts from Effort 1
3. Use libraries from Effort 2
4. No dependencies on Features A or C

### Feature C Engineer:
1. Work in features/feature_c/ directory only
2. Import contracts from Effort 1
3. Use libraries from Effort 2
4. No dependencies on Features A or B
```

---

### Effort 6: [EFFORT_NAME] - Integration
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-integration`  
**Can Parallelize**: No (depends on all parallel features)  
**Parallel With**: None  
**Size Estimate**: [NUMBER] lines (MUST be <800)  
**Dependencies**: Efforts 3, 4, 5 (MUST WAIT for parallel work)  

#### Files to Create
```yaml
new_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/integration/orchestrator.go
    lines: ~300
    purpose: Orchestrate all features
    imports:
      - pkg/phase[PHASE]/wave[WAVE]/features/feature_a
      - pkg/phase[PHASE]/wave[WAVE]/features/feature_b
      - pkg/phase[PHASE]/wave[WAVE]/features/feature_c
      
  - path: pkg/phase[PHASE]/wave[WAVE]/integration/router.go
    lines: ~200
    purpose: Route requests to features
    
  - path: pkg/phase[PHASE]/wave[WAVE]/integration/config.go
    lines: ~100
    purpose: Unified configuration
```

#### Integration Instructions
```markdown
1. Import all feature implementations
2. Create unified service interface
3. Implement request routing
4. Add cross-feature coordination
5. Implement health checks
6. Add integration tests
```

---

### Effort 7: [EFFORT_NAME] - Testing & Documentation
**Branch**: `[PROJECT_PREFIX/]phase[PHASE]/wave[WAVE]/effort-testing`  
**Can Parallelize**: No (requires integration complete)  
**Parallel With**: None  
**Size Estimate**: N/A (tests and docs don't count)  
**Dependencies**: Effort 6 (integration complete)  

#### Files to Create
```yaml
test_files:
  - path: pkg/phase[PHASE]/wave[WAVE]/integration_test.go
    purpose: End-to-end tests
    
  - path: pkg/phase[PHASE]/wave[WAVE]/benchmark_test.go
    purpose: Performance benchmarks
    
documentation:
  - path: pkg/phase[PHASE]/wave[WAVE]/README.md
    purpose: Wave documentation
    
  - path: pkg/phase[PHASE]/wave[WAVE]/api/API.md
    purpose: API documentation
    
  - path: docs/phase[PHASE]-wave[WAVE]-guide.md
    purpose: User guide
```

## 🔄 Code Reuse Matrix

### Cross-Effort Dependencies
```yaml
effort_dependencies:
  effort_1_provides:
    - api.Service interface
    - api.Request model
    - api.Response model
    - api.Error types
    
  effort_2_provides:
    - lib.Client
    - lib.Validator
    - lib.Converter
    
  effort_3_consumes:
    - effort_1: all interfaces
    - effort_2: Client, Validator
    
  effort_4_consumes:
    - effort_1: all interfaces
    - effort_2: Client, Converter
    
  effort_5_consumes:
    - effort_1: all interfaces
    - effort_2: Validator, Converter
    
  effort_6_consumes:
    - efforts_3_5: all implementations
    - effort_2: all libraries
    - effort_1: all contracts
```

### Import Map
```go
// Standard import hierarchy for all efforts
import (
    // Standard library
    "context"
    "fmt"
    
    // Phase 1 common (foundation)
    "pkg/phase1/common/auth"
    "pkg/phase1/common/logger"
    
    // Current phase, previous waves
    "pkg/phase[PHASE]/wave1/api"
    wave2lib "pkg/phase[PHASE]/wave2/lib"
    
    // Current wave contracts (Effort 1)
    "pkg/phase[PHASE]/wave[WAVE]/api"
    
    // Current wave libraries (Effort 2)
    "pkg/phase[PHASE]/wave[WAVE]/lib"
)
```

## 📈 Implementation Sequence

### Execution Order
```
Step 1:    Effort 1 (Contracts)     ████ BLOCKING
Step 2:    Effort 2 (Libraries)     ████████ 
Step 3:    Efforts 3-5 (Parallel)   ████████████ (3 engineers)
Step 4:    Effort 6 (Integration)   ████
Step 5:    Effort 7 (Testing/Docs)  ████
```

### Critical Path
1. **First**: Complete ALL contracts (Effort 1) - BLOCKS everything
2. **Second**: Complete shared libraries (Effort 2) - BLOCKS features
3. **Third**: Parallel feature work (Efforts 3-5) - 3 engineers simultaneously
4. **Fourth**: Integration (Effort 6) - Requires all features
5. **Fifth**: Testing and documentation (Effort 7)

### Resource Allocation
```yaml
engineer_assignments:
  engineer_1:
    - first: Effort 1 (contracts)
    - second: Effort 2 (libraries)
    - third: Effort 3 (feature A)
    - fourth: Effort 6 (integration)
    
  engineer_2:
    - first: Review Effort 1
    - third: Effort 4 (feature B)
    - fifth: Effort 7 (testing)
    
  engineer_3:
    - first: Review Effort 1
    - third: Effort 5 (feature C)
    - fifth: Effort 7 (documentation)
```

## ✅ Implementation Checklist

### Pre-Implementation
- [ ] Architecture plan reviewed and understood
- [ ] All contracts from architecture plan identified
- [ ] Previous phase/wave code identified for reuse
- [ ] Branches created for each effort
- [ ] Engineers assigned to parallel efforts

### During Implementation
- [ ] Effort 1 contracts match architecture plan exactly
- [ ] Effort 2 libraries are genuinely reusable
- [ ] Efforts 3-5 have no interdependencies
- [ ] Size monitoring (must stay <800 lines)
- [ ] Continuous integration tests running

### Post-Implementation
- [ ] All efforts under 800 lines (verified by line counter)
- [ ] Integration tests passing
- [ ] Code reviews completed
- [ ] Documentation complete
- [ ] Performance benchmarks met

## 🚨 Common Pitfalls to Avoid

### Pitfall 1: Incomplete Contracts
**Problem**: Effort 1 doesn't define all interfaces  
**Impact**: Later efforts blocked or need rework  
**Prevention**: Architect reviews Effort 1 before implementation  

### Pitfall 2: Parallel Effort Dependencies
**Problem**: Efforts 3-5 depend on each other  
**Impact**: Can't actually parallelize  
**Prevention**: Strict directory/package separation  

### Pitfall 3: Size Limit Violations
**Problem**: Effort exceeds 800 lines  
**Impact**: Requires splitting, delays wave  
**Prevention**: Monitor continuously, plan for 700 max  

### Pitfall 4: Ignoring Reuse Requirements
**Problem**: Reimplementing existing functionality  
**Impact**: Inconsistency, wasted effort  
**Prevention**: Explicit import requirements in plan  

## 📚 References

### Required Reading
- [Architecture Plan](./PHASE-[PHASE]-WAVE-[WAVE]-ARCHITECTURE-PLAN.md) - MUST FOLLOW
- [Phase Architecture](../PHASE-[PHASE]-ARCHITECTURE-PLAN.md) - For context
- [Master Plan](../../IMPLEMENTATION-PLAN.md) - For vision alignment

### Code Examples
- [Previous Wave Implementation](../wave[PREV]/README.md)
- [Phase 1 Examples](../../phase1/examples/)

### Standards and Guidelines
- [Coding Standards](../../docs/coding-standards.md)
- [Testing Requirements](../../docs/testing-requirements.md)
- [API Design Guide](../../docs/api-design.md)

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | [DATE] | Code Reviewer Agent | Initial implementation plan from architecture |