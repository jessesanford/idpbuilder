# Split Plan for E3.1.4-trust-store

## Executive Summary
**Current Size**: 1514 lines (89% over limit)
**Target**: 4 splits, each <400 lines (50% safety margin)
**Strategy**: Functional decomposition along natural boundaries
**Risk**: LOW - Clear interfaces between components

## Violation Analysis
- **Hard Limit**: 800 lines
- **Current Implementation**: 1514 lines
- **Excess**: 714 lines
- **Severity**: MAJOR - requires complete restructuring

## Split Architecture

### Split Boundaries Overview
```
Split 001: Core Interfaces & Types (350 lines)
  └── Defines all interfaces and shared types
  
Split 002: Storage Implementation (380 lines)
  └── Filesystem storage with watch capabilities
  
Split 003: Pool Management (390 lines)
  └── Certificate pool with hot-reload
  
Split 004: Configuration & Tests (394 lines)
  └── Config management and comprehensive testing
```

## Detailed Split Plans

### Split 001: Core Interfaces & Types
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.4-trust-store-split-001`
**Size Estimate**: 350 lines
**Dependencies**: None (foundational)

#### Files to Create:
```
pkg/oci/certificates/
├── types.go (120 lines)
│   ├── Certificate struct
│   ├── Event and EventType
│   ├── Validation interfaces
│   └── Common error types
├── interfaces.go (100 lines) 
│   ├── CertificateStore interface
│   ├── CertificateValidator interface
│   ├── CertPoolManager interface
│   └── ConfigLoader interface
├── errors.go (50 lines)
│   └── Custom error types and handling
└── types_test.go (80 lines)
    └── Basic type validation tests
```

#### Functionality:
- Define all public interfaces for the certificate system
- Core data structures (Certificate, Event, etc.)
- Error types and constants
- Interface contracts for all components

#### Implementation Notes:
- Pure interface definitions - no implementation
- All types must be exported for use by other splits
- Include comprehensive godoc comments
- No external dependencies beyond standard library

---

### Split 002: Storage Implementation
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.4-trust-store-split-002`
**Size Estimate**: 380 lines
**Dependencies**: Split 001 (interfaces and types)

#### Files to Create:
```
pkg/oci/certificates/
├── storage.go (250 lines)
│   ├── FilesystemStore implementation
│   ├── Save/Load/Delete operations
│   ├── List functionality
│   └── Basic validation
├── watch.go (80 lines)
│   ├── File watching implementation
│   ├── Event dispatch
│   └── Change detection
└── storage_test.go (50 lines)
    └── Core storage operations tests
```

#### Functionality:
- Implement CertificateStore interface
- Filesystem-based storage operations
- File watching for hot-reload support
- Certificate discovery from filesystem

#### Implementation Notes:
- Import interfaces from Split 001
- Focus on storage operations only
- Defer pool management to Split 003
- Use fsnotify for file watching

---

### Split 003: Pool Management
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.4-trust-store-split-003`
**Size Estimate**: 390 lines
**Dependencies**: Splits 001 & 002

#### Files to Create:
```
pkg/oci/certificates/
├── pool.go (200 lines)
│   ├── CertPoolManager implementation
│   ├── System/custom pool separation
│   ├── Hot-reload logic
│   └── Pool operations
├── validation.go (100 lines)
│   ├── Certificate validators
│   ├── Chain validation
│   ├── Expiry checking
│   └── Permission validation
└── pool_test.go (90 lines)
    └── Pool management tests
```

#### Functionality:
- Certificate pool management with hot-reload
- System and custom pool separation
- Certificate validation pipeline
- Dynamic pool updates

#### Implementation Notes:
- Import from Splits 001 and 002
- Focus on pool operations
- Implement validation strategies
- Handle concurrent pool access

---

### Split 004: Configuration & Integration Tests
**Branch**: `idpbuidler-oci-mgmt/phase3/wave1/E3.1.4-trust-store-split-004`
**Size Estimate**: 394 lines
**Dependencies**: Splits 001, 002, & 003

#### Files to Create:
```
pkg/oci/certificates/
├── config.go (200 lines)
│   ├── CertificateConfig struct
│   ├── Environment loading
│   ├── Config validation
│   └── Default configurations
├── loader.go (80 lines)
│   ├── Config file loading
│   ├── Environment override
│   └── Path resolution
└── integration_test.go (114 lines)
    ├── End-to-end tests
    ├── Hot-reload testing
    └── Error scenarios
```

#### Functionality:
- Configuration management system
- Environment variable support
- YAML/JSON config loading
- Comprehensive integration testing

#### Implementation Notes:
- Import all previous splits
- Focus on configuration and wiring
- Include integration tests for full system
- Test hot-reload scenarios

---

## Implementation Sequence

### Phase 1: Foundation (Split 001)
1. Create all interfaces and types
2. Define error types
3. Establish contracts
4. Unit test type validation

### Phase 2: Storage (Split 002)
1. Implement filesystem storage
2. Add file watching
3. Test storage operations
4. Verify isolation from pool logic

### Phase 3: Pool Management (Split 003)
1. Implement pool manager
2. Add validation pipeline
3. Wire to storage layer
4. Test concurrent operations

### Phase 4: Configuration & Integration (Split 004)
1. Implement configuration system
2. Add environment support
3. Create integration tests
4. Verify full system operation

## Integration Strategy

### Interface Contracts
Each split communicates through well-defined interfaces from Split 001:
- Storage ← CertificateStore interface
- Pool ← CertPoolManager interface  
- Config ← ConfigLoader interface

### Dependency Graph
```
Split 001 (Interfaces)
    ↓
Split 002 (Storage)
    ↓
Split 003 (Pool) 
    ↓
Split 004 (Config & Tests)
```

### Merge Strategy
1. Each split creates a PR to its own branch
2. Sequential merging in dependency order
3. Integration testing after each merge
4. Final merge to parent effort branch

## Risk Mitigation

### Size Management
- Each split targets <400 lines (50% margin)
- Regular measurement during implementation
- Clear stop points if approaching limits

### Quality Assurance
- Each split must compile independently
- Unit tests in each split
- Integration tests in final split
- Code review for each split

### Coordination
- Clear interface boundaries defined upfront
- No circular dependencies
- Sequential implementation order
- Integration points well-documented

## Success Criteria

### Per-Split Criteria
- [ ] Under 400 lines actual code
- [ ] Compiles independently  
- [ ] Unit tests pass
- [ ] No circular dependencies
- [ ] Clear documentation

### Overall Criteria
- [ ] Total implementation unchanged
- [ ] All tests passing
- [ ] Full feature parity
- [ ] Clean integration
- [ ] Each split reviewable independently

## Validation Checklist

### Pre-Implementation
- [x] Current size verified: 1514 lines
- [x] Split boundaries identified
- [x] Dependencies mapped
- [x] Interfaces defined

### During Implementation
- [ ] Regular size checks with line-counter.sh
- [ ] Interface compliance verified
- [ ] Test coverage maintained
- [ ] Documentation updated

### Post-Implementation
- [ ] All splits under limit
- [ ] Integration tests pass
- [ ] Performance unchanged
- [ ] No functionality lost

## Notes for SW Engineers

### Critical Instructions
1. **MUST** implement splits in order (001 → 002 → 003 → 004)
2. **MUST** check size every 100 lines with tools/line-counter.sh
3. **MUST** maintain interface contracts exactly
4. **MUST** ensure each split compiles independently

### Common Pitfalls to Avoid
- Don't combine splits even if under limit
- Don't skip dependency ordering
- Don't modify interfaces after Split 001
- Don't include generated code in measurements

### Success Tips
- Start with interfaces (Split 001) to establish contracts
- Keep implementations focused on single responsibility
- Test each split thoroughly before moving on
- Use dependency injection for clean separation

## Orchestrator Metadata Placeholder
<!-- The orchestrator will add infrastructure metadata here for each split -->
<!-- Including: WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH -->
<!-- SW Engineers MUST read this metadata when implementing -->