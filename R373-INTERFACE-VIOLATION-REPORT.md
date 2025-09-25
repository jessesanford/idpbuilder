# R373 Interface Violation Report - Critical Architectural Failure

## Executive Summary

A catastrophic architectural failure was discovered where 11 effort branches created multiple incompatible implementations of the same functionality, completely ignoring interfaces that were already established in earlier branches. This violation caused a cascade of integration failures requiring complete rebuilding of all affected branches.

## The Violation: Three Incompatible Push Methods

### Branch 4 - Original Interface Definition
```go
// Branch: phase2-wave1-effort-4
// File: pkg/registry/interface.go

type Registry interface {
    Push(ctx context.Context, image string, content io.Reader) error
    Pull(ctx context.Context, image string) (io.Reader, error)
}
```

### Branch 9 - Incompatible Implementation
```go
// Branch: phase2-wave1-effort-9
// File: pkg/gitea/client.go

type GiteaClient struct{}

// VIOLATION: Different signature, not implementing Registry interface!
func (g *GiteaClient) Push(ctx context.Context, image v1.Image, ref string, opts PushOptions) error {
    // Implementation with completely different parameters
}
```

### Branch 11 - Using Wrong Implementation
```go
// Branch: phase2-wave1-effort-11
// File: pkg/builder/oci.go

// Using GiteaClient's incompatible Push signature instead of Registry interface
client.Push(ctx, v1Image, "latest", opts) // Won't work with Registry interface!
```

## Impact Analysis

### Immediate Consequences
1. **11 branches** required complete rebuilding
2. **3 incompatible APIs** for the same functionality
3. **Integration impossible** without major refactoring
4. **Testing failures** across all dependent branches
5. **Architecture completely broken**

### Root Causes
1. **No pre-planning research** - Agents didn't check for existing interfaces
2. **No interface compliance validation** - Code reviewers didn't enforce implementation
3. **No cross-branch awareness** - Engineers created new APIs instead of using existing ones
4. **No architectural governance** - Architect didn't catch duplicates in review

## The Solution: R373 and R374 Rules

### R373: Mandatory Code Reuse and Interface Compliance
- **FORBIDS** creating duplicate implementations
- **REQUIRES** implementing existing interfaces exactly
- **PENALTY**: -100% for violations (immediate failure)

### R374: Pre-Planning Research Protocol
- **REQUIRES** searching all existing code before planning
- **DOCUMENTS** all interfaces and APIs found
- **SPECIFIES** what must be reused
- **LISTS** forbidden duplications

## Enforcement Mechanisms

### 1. During Planning (Code Reviewer)
```markdown
## Pre-Planning Research Results (R374 MANDATORY)
### Existing Interfaces Found:
- Registry interface: Push(ctx, image string, content io.Reader) error
### REQUIRED INTEGRATIONS:
- MUST implement Registry interface with EXACT signature
### FORBIDDEN:
- DO NOT create new Push method with different signature
```

### 2. During Implementation (SW Engineer)
```go
// CORRECT: Implements existing interface
type GiteaClient struct{}
func (g *GiteaClient) Push(ctx context.Context, image string, content io.Reader) error {
    // Matches Registry interface exactly
}
```

### 3. During Review
```bash
# Run validation tool
./tools/R373-validate-no-duplicates.sh

# Output:
❌ CRITICAL: Method 'Push' has 3 different signatures:
    - Push(ctx context.Context, image string, content io.Reader) error
    - Push(ctx context.Context, image v1.Image, ref string, opts PushOptions) error
    - Push(image string) error
```

### 4. During Architecture Review
- **REJECT** entire wave if duplicates found
- **REQUIRE** refactoring to use existing interfaces
- **VALIDATE** all implementations match interface contracts

## Correct Pattern Example

### Phase 2, Wave 1, Effort 4
```go
// Define the interface
type Registry interface {
    Push(ctx context.Context, image string, content io.Reader) error
}
```

### Phase 2, Wave 1, Effort 9
```go
// Import and implement the EXISTING interface
import "github.com/project/effort-repo/branch4/registry"

type GiteaClient struct{}

// Implement Registry interface EXACTLY
func (g *GiteaClient) Push(ctx context.Context, image string, content io.Reader) error {
    // Internal conversion if needed, but interface stays consistent
    v1Image := convertToV1Image(image, content)
    return g.internalPush(v1Image)
}
```

### Phase 2, Wave 1, Effort 11
```go
// Use the Registry interface
var registry Registry = NewGiteaClient()
registry.Push(ctx, "my-image:latest", imageContent) // Works perfectly!
```

## Prevention Checklist

### For Code Reviewers (Planning)
- [ ] Run R374 research tool before creating any plan
- [ ] Document ALL existing interfaces found
- [ ] Specify which interfaces MUST be implemented
- [ ] List what CANNOT be created (already exists)
- [ ] Include exact method signatures to implement

### For SW Engineers (Implementation)
- [ ] Read "REQUIRED INTEGRATIONS" section of plan
- [ ] Import existing interfaces/packages
- [ ] Implement interfaces with EXACT signatures
- [ ] Never create "better" versions of existing APIs
- [ ] Run R373 validation before committing

### For Architects (Review)
- [ ] Check for duplicate interface definitions
- [ ] Verify all implementations match interfaces
- [ ] Confirm no competing APIs exist
- [ ] Validate proper code reuse

## Key Takeaways

1. **Research BEFORE Planning** - Always check what already exists
2. **Reuse BEFORE Creating** - Never duplicate existing functionality
3. **Implement EXACTLY** - Match interface signatures perfectly
4. **Validate CONTINUOUSLY** - Check compliance at every stage

## Tools and Resources

- **Research Tool**: `tools/R374-pre-planning-research.sh`
- **Validation Tool**: `tools/R373-validate-no-duplicates.sh`
- **Rule R373**: `rule-library/R373-mandatory-code-reuse-and-interface-compliance.md`
- **Rule R374**: `rule-library/R374-pre-planning-research-protocol.md`

## Penalty Structure

| Violation | Penalty | Recovery Required |
|-----------|---------|-------------------|
| Creating duplicate interface | -100% FAIL | Complete refactoring |
| Not researching existing code | -50% | Add research, update plan |
| Incompatible implementation | -100% FAIL | Reimplement correctly |
| Not documenting interfaces | -25% | Update documentation |

## Conclusion

The creation of three incompatible Push methods for the same functionality represents a complete breakdown of architectural governance. With R373 and R374 now in place as SUPREME LAWS, this type of catastrophic failure will be prevented through:

1. **Mandatory research** before any planning
2. **Explicit documentation** of what must be reused
3. **Strict validation** of interface compliance
4. **Immediate rejection** of duplicate implementations

This is not optional - it's a fundamental requirement for maintaining a coherent, maintainable codebase.