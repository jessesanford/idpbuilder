# Effort Implementation Plan: Provider Interface Definition

<!-- STORAGE LOCATION: This plan should be saved in:
     efforts/phase1/wave1/P1W1-E1-provider-interface/IMPLEMENTATION-PLAN.md
     within the effort's working directory. This keeps plans organized and separate from code. -->

## 🚨 CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Effort**: P1W1-E1 - Provider Interface Definition
**Branch**: `phase1/wave1/P1W1-E1-provider-interface`
**Base Branch**: `main`
**Base Branch Reason**: First effort in Phase 1 Wave 1, establishing foundational interfaces
**Can Parallelize**: Yes
**Parallel With**: P1W1-E2, P1W1-E3, P1W1-E4 (all Wave 1 efforts can run in parallel)
**Size Estimate**: 150 lines (MUST be <800)
**Dependencies**: None (foundational effort)
**Dependent Efforts**: ALL subsequent efforts depend on these interfaces
**Atomic PR**: ✅ This effort = ONE PR to main (R220 REQUIREMENT)

## 📋 Source Information
**Wave Plan**: PROJECT-IMPLEMENTATION-PLAN.md (Phase 1, Wave 1, Section P1W1-E1)
**Effort Section**: P1W1-E1
**Created By**: Code Reviewer Agent
**Date**: 2025-09-27
**Extracted**: 2025-09-27T14:36:44Z

## 🔴 BASE BRANCH VALIDATION (R337 MANDATORY)
**The orchestrator-state.json is the SOLE SOURCE OF TRUTH for base branches!**
- Base branch MUST be explicitly specified above: `main`
- Base branch MUST match what's in orchestrator-state.json
- Reason MUST explain why this base: First effort in Phase 1, no prior integration branches exist
- Orchestrator MUST record this in state file before creating infrastructure

## 🚀 Parallelization Context
**Can Parallelize**: Yes
**Parallel With**: P1W1-E2 (OCI Package Format), P1W1-E3 (Registry Config), P1W1-E4 (CLI Contracts)
**Blocking Status**: No - This effort does not block other Wave 1 efforts
**Parallel Group**: Wave 1 Foundation Group (E1, E2, E3, E4)
**Orchestrator Guidance**: Can spawn all Wave 1 efforts immediately and simultaneously

## 🚨 EXPLICIT SCOPE DEFINITION (R311 MANDATORY)

### IMPLEMENT EXACTLY (BE SPECIFIC!)

#### Functions to Create (EXACTLY 0 - INTERFACES ONLY)
```go
// NO FUNCTIONS - This effort defines interfaces and types only
// All implementations will come in later efforts
```

#### Types/Structs to Define (EXACTLY 8)
```go
// Type 1: Core Provider interface
type Provider interface {
    // Push an OCI artifact to the registry
    Push(ctx context.Context, ref string, artifact Artifact) error
    // Pull an OCI artifact from the registry
    Pull(ctx context.Context, ref string) (Artifact, error)
    // List artifacts in a repository
    List(ctx context.Context, repository string) ([]ArtifactInfo, error)
    // Delete an artifact from the registry
    Delete(ctx context.Context, ref string) error
}

// Type 2: Artifact structure
type Artifact struct {
    MediaType string              // OCI media type
    Manifest  []byte              // Raw manifest content
    Layers    []Layer             // Artifact layers
    Config    []byte              // Configuration blob
    Annotations map[string]string // OCI annotations
}

// Type 3: Layer structure
type Layer struct {
    MediaType string // Layer media type
    Digest    string // Content digest
    Size      int64  // Layer size in bytes
    Data      []byte // Layer content
}

// Type 4: Artifact information
type ArtifactInfo struct {
    Reference   string            // Full reference (registry/repo:tag)
    Digest      string            // Manifest digest
    Tags        []string          // Associated tags
    Size        int64             // Total size
    Created     time.Time         // Creation timestamp
    Annotations map[string]string // Metadata annotations
}

// Type 5: Provider configuration
type ProviderConfig struct {
    Registry string            // Registry URL
    Auth     AuthConfig        // Authentication configuration
    Insecure bool              // Allow insecure connections
    Timeout  time.Duration     // Operation timeout
}

// Type 6: Authentication configuration
type AuthConfig struct {
    Username      string // Basic auth username
    Password      string // Basic auth password
    Token         string // Bearer token
    RegistryToken string // Registry-specific token
}

// Type 7: Provider errors
type ProviderError struct {
    Op   string // Operation that failed
    Ref  string // Reference involved
    Err  error  // Underlying error
}

// Type 8: Registry capabilities
type RegistryCapabilities struct {
    SupportsDelete bool // Registry supports delete
    SupportsOCI    bool // Full OCI compliance
    MaxLayerSize   int64 // Maximum layer size
}
```

#### Endpoints/Handlers (if applicable)
```go
// NO ENDPOINTS in this effort - interface definitions only
```

### 🛑 DO NOT IMPLEMENT (SCOPE BOUNDARIES)

**EXPLICITLY FORBIDDEN IN THIS EFFORT:**
- ❌ DO NOT add any concrete implementations
- ❌ DO NOT implement actual registry communication
- ❌ DO NOT add HTTP client code
- ❌ DO NOT implement authentication logic
- ❌ DO NOT add validation beyond interface definitions
- ❌ DO NOT create helper/utility functions
- ❌ DO NOT add comprehensive error handling implementations
- ❌ DO NOT add logging implementations
- ❌ DO NOT write integration tests (unit tests for types only)
- ❌ DO NOT implement nice-to-have features

### 📊 REALISTIC SIZE CALCULATION

```
Component Breakdown:
- Provider interface:            15 lines
- Artifact type:                 10 lines
- Layer type:                     8 lines
- ArtifactInfo type:             10 lines
- ProviderConfig type:            8 lines
- AuthConfig type:                8 lines
- ProviderError type:             7 lines
- RegistryCapabilities type:      7 lines
- Package documentation:         20 lines
- Error method implementations:  15 lines
- Type documentation:            42 lines (6 lines per type average)

TOTAL ESTIMATE: 150 lines (must be <800)
BUFFER: 650 lines for unforeseen needs
```

## 🔴🔴🔴 PRE-PLANNING RESEARCH RESULTS (R374 MANDATORY) 🔴🔴🔴

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| None - This is the foundational effort | N/A | N/A | N/A |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| None - Creating new interfaces | N/A | N/A | N/A |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| None - Defining new APIs | N/A | N/A | First effort |

### FORBIDDEN DUPLICATIONS (R373)
- ✅ No existing provider interfaces to duplicate
- ✅ Creating foundational interfaces for first time

### REQUIRED INTEGRATIONS (R373)
- ✅ MUST be compatible with go-containerregistry types
- ✅ MUST support standard OCI media types
- ✅ MUST align with OCI Distribution Spec

## 📁 Files to Create

### Primary Implementation Files
```yaml
new_files:
  - path: pkg/providers/interface.go
    lines: ~60 MAX
    purpose: Core provider interface definition
    contains:
      - Provider interface
      - ProviderError type and methods
      - Package documentation

  - path: pkg/providers/types.go
    lines: ~60 MAX
    purpose: Common types and structures
    contains:
      - Artifact type
      - Layer type
      - ArtifactInfo type
      - RegistryCapabilities type

  - path: pkg/providers/errors.go
    lines: ~30 MAX
    purpose: Error definitions and helpers
    contains:
      - Error constants
      - ProviderError methods
      - Common error scenarios
```

### Test Files
```yaml
test_files:
  - path: pkg/providers/types_test.go
    lines: ~50 MAX
    coverage_target: 80%
    test_functions:
      - TestArtifactCreation  # ~15 lines
      - TestProviderErrorFormat  # ~15 lines
      - TestAuthConfigMasking  # ~20 lines
```

## 📦 Files to Import/Reuse

### From Previous Efforts (This Wave)
```yaml
this_wave_imports:
  # None - this is the first effort, others depend on us
```

### From Previous Waves/Phases
```yaml
previous_work_imports:
  # None - this is Phase 1 Wave 1 Effort 1
```

## 🔗 Dependencies

### Effort Dependencies
- **Must Complete First**: None (first effort)
- **Can Run in Parallel With**: P1W1-E2, P1W1-E3, P1W1-E4
- **Blocks**: All Phase 1 Wave 2 efforts need these interfaces

### Technical Dependencies
- Standard library only (context, time, errors packages)
- No external dependencies in this effort

## 🔴 ATOMIC PR REQUIREMENTS (R220 - SUPREME LAW)

### 🔴🔴🔴 PARAMOUNT: Independent Mergeability (R307) 🔴🔴🔴
**This effort MUST be mergeable at ANY time, even YEARS later:**
- ✅ Must compile when merged alone to main
- ✅ Must NOT break any existing functionality
- ✅ Interfaces only - no runtime impact
- ✅ Must work even if next effort merges 6 months later
- ✅ No external dependencies to break

### Feature Flags for This Effort
```yaml
feature_flags:
  # No feature flags needed - interface definitions only
  # Implementations in later efforts will use flags
```

### 🚨🚨🚨 R355 PRODUCTION READY CODE (SUPREME LAW #5) 🚨🚨🚨

**ALL CODE MUST BE PRODUCTION READY - NO EXCEPTIONS**

#### ✅ REQUIRED PATTERNS:
```go
// Provider interface is complete and production-ready
// No stub methods or placeholders
// All methods have clear contracts
// Error types are fully defined
```

### Interface Implementations (Instead of Stubs)
```yaml
interfaces:
  - name: "Provider"
    implements: "N/A - defining the interface"
    type: "Interface Definition"
    notes: "Concrete implementations in Wave 2"
    production_ready: true
```

### PR Mergeability Checklist
- [x] PR can merge to main independently
- [x] Build passes with just this PR
- [x] All tests pass in isolation
- [x] No feature flags needed (interfaces only)
- [x] No external dependencies
- [x] No breaking changes to existing code
- [x] Backward compatible with main

## 🔴 MANDATORY ADHERENCE CHECKPOINTS (R311)

### Before Starting:
```bash
echo "EFFORT SCOPE LOCKED:"
echo "✓ Functions: EXACTLY 0 (interfaces only)"
echo "✓ Types: EXACTLY 8 (Provider, Artifact, Layer, etc.)"
echo "✓ Endpoints: EXACTLY 0 (no handlers)"
echo "✓ Tests: EXACTLY 3 basic type tests"
echo "✗ Implementations: NONE"
echo "✗ Extra features: NONE"
echo "✗ Optimizations: NONE"
```

### During Implementation:
```bash
# Check scope adherence after each component
TYPE_COUNT=$(grep -c "^type " pkg/providers/*.go 2>/dev/null || echo 0)
if [ "$TYPE_COUNT" -gt 8 ]; then
    echo "⚠️ WARNING: Exceeding type count! Stop adding!"
fi
```

## 📝 Implementation Instructions

### Step-by-Step Guide
1. **Scope Acknowledgment**
   - Read and acknowledge DO NOT IMPLEMENT section
   - Confirm only 8 types, 0 functions
   - Create tracking checklist

2. **Implementation Order**
   - Start with `pkg/providers/interface.go` - Provider interface
   - Create `pkg/providers/types.go` - Core types
   - Add `pkg/providers/errors.go` - Error definitions
   - Write minimal type tests only

3. **Key Implementation Details**
   ```go
   // pkg/providers/interface.go
   package providers

   import (
       "context"
   )

   // Provider defines the contract for OCI registry providers
   type Provider interface {
       Push(ctx context.Context, ref string, artifact Artifact) error
       Pull(ctx context.Context, ref string) (Artifact, error)
       List(ctx context.Context, repository string) ([]ArtifactInfo, error)
       Delete(ctx context.Context, ref string) error
   }
   ```

4. **Integration Points**
   - No integrations in this effort
   - Other efforts will import these interfaces
   - Keep interfaces minimal and focused

## ✅ Test Requirements

### Coverage Requirements
- **Minimum Coverage**: 80% (for types only)
- **Critical Paths**: Type creation and error formatting
- **Error Handling**: ProviderError.Error() method

### Test Categories
```yaml
required_tests:
  unit_tests:
    - Artifact type creation
    - ProviderError formatting
    - AuthConfig field masking

  # NO integration tests in this effort
  # NO performance tests in this effort
```

## 📏 Size Constraints
**Target Size**: 150 lines (from wave plan)
**Maximum Size**: 800 lines (HARD LIMIT)
**Current Size**: 0 lines (to be updated)

### Size Monitoring Protocol
```bash
# Check size after creating each file
cd /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave1/P1W1-E1-provider-interface
find pkg/providers -name "*.go" | xargs wc -l

# If approaching 200 lines (well under limit):
# Should not happen with just interfaces
```

## 🏁 Completion Criteria

### Implementation Checklist
- [ ] Created pkg/providers/interface.go with Provider interface
- [ ] Created pkg/providers/types.go with all 8 types
- [ ] Created pkg/providers/errors.go with error handling
- [ ] Size verified under 150 lines
- [ ] No concrete implementations added

### Quality Checklist
- [ ] Test coverage ≥80% for types
- [ ] All tests passing
- [ ] No linting errors
- [ ] Error() method for ProviderError
- [ ] Comments for all exported types

### Documentation Checklist
- [ ] Package documentation in interface.go
- [ ] Comments for all interface methods
- [ ] Comments for all exported types
- [ ] Clear contract definitions

### Review Checklist
- [ ] Self-review completed
- [ ] Code committed and pushed
- [ ] Ready for Code Reviewer assessment
- [ ] No blocking issues

## 📊 Progress Tracking

### Work Log
```markdown
## 2025-09-27 - Session 1
- Created IMPLEMENTATION-PLAN.md
- Defined scope and boundaries
- Ready for SW Engineer implementation

[SW Engineer will update during implementation]
```

## ⚠️ Important Notes

### Parallelization Reminder
- This effort can run simultaneously with P1W1-E2, E3, E4
- No need to wait for other Wave 1 efforts
- All Wave 1 efforts are independent foundations

### Common Pitfalls to Avoid (R311 ENFORCEMENT)
1. **SCOPE CREEP**: Adding implementations = FAILURE
2. **OVER-ENGINEERING**: Complex types = overrun
3. **ASSUMPTIONS**: Adding "helpful" methods = VIOLATION
4. **Dependencies**: Keep it stdlib only
5. **Test Coverage**: Type tests only, no integration
6. **Isolation**: Stay in pkg/providers/
7. **Parallelization**: No dependencies on other Wave 1 efforts

### Success Criteria Checklist
- [ ] Created EXACTLY 8 types (no more)
- [ ] Created EXACTLY 0 implementations
- [ ] Wrote EXACTLY 3 type tests
- [ ] Total lines under 150
- [ ] NO unauthorized features added
- [ ] Followed all scope boundaries

## 📚 References

### Source Documents
- [Master Plan](/home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase2/wave1/integration-workspace/planning/project/PROJECT-IMPLEMENTATION-PLAN.md)
- [OCI Distribution Spec](https://github.com/opencontainers/distribution-spec)
- [OCI Image Spec](https://github.com/opencontainers/image-spec)

### Standards
- Go interface best practices
- OCI specification compliance
- Error handling patterns

---

**Remember**: This is the foundational effort for the entire OCI implementation. These interfaces will be used by ALL subsequent efforts. Keep them clean, minimal, and focused on the essential contract for registry operations.

CONTINUE-SOFTWARE-FACTORY=TRUE