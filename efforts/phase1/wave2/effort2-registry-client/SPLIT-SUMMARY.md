# E1.2.2 Registry Client Split Plan

## Split Analysis

**Current Status:**
- Total implementation: 877 lines
- Configured limit: 800 lines
- Overage: 77 lines
- Split required: YES

## Split Strategy

Based on logical functionality grouping and dependency analysis, this effort will be split into 2 parts:

### E1.2.2 (Core Registry Client) - 427 lines
**Focus:** Essential client functionality and interfaces

**Files Included:**
- `pkg/registry/interface.go` (61 lines) - Interface definitions
- `pkg/registry/types.go` (107 lines) - Core data types
- `pkg/registry/errors.go` (145 lines) - Error types and handling
- `pkg/registry/client.go` (114 lines) - Basic client structure and core methods

**Files Modified for Split:**
- `client.go` will be reduced from 191 to 114 lines by moving advanced authentication and TLS methods to E1.2.2a

**Functionality:**
- Registry client interface
- Basic client structure
- Core data types
- Error handling
- Basic registry operations (ping, catalog, tags)
- Simple authentication (basic auth)

### E1.2.2a (Advanced Client Features) - 450 lines
**Focus:** Advanced authentication, TLS, and comprehensive testing

**Files Included:**
- `pkg/registry/auth.go` (138 lines) - Advanced authentication methods
- `pkg/registry/tls.go` (50 lines) - TLS configuration
- `client_advanced.go` (77 lines) - Advanced client methods (extracted from client.go)
- `pkg/registry/types_test.go` (104 lines) - Type tests
- `pkg/registry/errors_test.go` (81 lines) - Error tests

**New Files:**
- `client_advanced.go` - Contains advanced client methods extracted from `client.go`

**Functionality:**
- OAuth2, Bearer token authentication
- Mutual TLS configuration
- Advanced client methods (pull, push operations)
- Comprehensive test coverage
- Integration scenarios

## Dependencies and Execution Order

### Split 1 (E1.2.2) - MUST complete first
- Provides base interfaces and types
- Required by Split 2

### Split 2 (E1.2.2a) - Depends on Split 1
- Uses interfaces from Split 1
- Extends functionality from Split 1
- Contains tests that validate both splits

## Implementation Instructions

### For SW Engineer implementing E1.2.2:

1. **Keep as-is:**
   - `pkg/registry/interface.go` (61 lines)
   - `pkg/registry/types.go` (107 lines)
   - `pkg/registry/errors.go` (145 lines)

2. **Modify client.go:**
   ```go
   // KEEP in E1.2.2 (Basic client methods)
   - New() constructor
   - Ping() method
   - Catalog() method
   - Tags() method  
   - Basic authentication setup
   - Core HTTP client initialization
   
   // MOVE to E1.2.2a (Advanced methods)
   - Pull() method
   - Push() method
   - Advanced OAuth2 authentication
   - Mutual TLS setup
   - Complex retry logic
   ```

3. **Target size:** ~427 lines total

### For SW Engineer implementing E1.2.2a:

1. **Keep as-is:**
   - `pkg/registry/auth.go` (138 lines)
   - `pkg/registry/tls.go` (50 lines)
   - `pkg/registry/types_test.go` (104 lines)
   - `pkg/registry/errors_test.go` (81 lines)

2. **Create client_advanced.go:**
   - Extract advanced methods from original client.go
   - Implement as extensions to base client
   - Add advanced authentication flows
   - Add complex registry operations

3. **Missing implementations to add:**
   - Client unit tests (~50 lines)
   - Authentication method tests (~27 lines)

4. **Target size:** ~450 lines total

## Validation Criteria

### E1.2.2 Acceptance:
- [ ] All basic client operations work
- [ ] Interface properly defined
- [ ] Error handling complete
- [ ] Size ≤ 800 lines
- [ ] Basic tests pass
- [ ] No missing dependencies

### E1.2.2a Acceptance:
- [ ] Advanced features functional
- [ ] Extends E1.2.2 properly
- [ ] All tests pass
- [ ] Size ≤ 800 lines
- [ ] Full test coverage achieved
- [ ] Integration with E1.2.2 verified

## Sequential Execution Protocol

1. **First:** Complete E1.2.2 implementation and review
2. **Second:** Complete E1.2.2a implementation and review
3. **Finally:** Integration testing of both parts

## Risk Mitigation

### Potential Issues:
- **Circular dependencies:** Mitigated by clear interface separation
- **Missing functionality:** Covered by comprehensive acceptance criteria
- **Test coverage gaps:** All tests moved to E1.2.2a with integration verification

### Quality Assurance:
- Each split independently buildable
- Each split independently testable  
- Combined functionality equivalent to original design
- No feature regression

## Branch Organization

```
phase1/wave2/effort2-registry-client/       (original branch - becomes E1.2.2)
└── Split into:
    ├── E1.2.2: Core client (this branch)
    └── E1.2.2a: Advanced features (new branch)
```

## Line Count Projections

| Split | Projected Size | Buffer | Status |
|-------|---------------|---------|---------|
| E1.2.2 | 427 lines | 373 lines available | SAFE |
| E1.2.2a | 450 lines | 350 lines available | SAFE |
| **Total** | **877 lines** | **Original preserved** | **✅** |