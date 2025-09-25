# 🔴🔴🔴 RULE R373: Mandatory Code Reuse and Interface Compliance (SUPREME LAW) 🔴🔴🔴

## Category
SUPREME LAW - ABSOLUTE HIGHEST PRIORITY

## Criticality
BLOCKING - Violation = -100% grading penalty = IMMEDIATE FAILURE

## Description
ABSOLUTELY FORBIDS creating duplicate or competing implementations of existing functionality. When interfaces, APIs, or implementations already exist, they MUST be reused EXACTLY. Creating alternative implementations when existing code could be reused is a CATASTROPHIC ARCHITECTURAL VIOLATION.

## Rationale
Multiple incompatible implementations of the same functionality:
1. Destroys architectural coherence
2. Creates impossible integration scenarios
3. Forces rebuilding of dependent branches
4. Violates fundamental DRY (Don't Repeat Yourself) principles
5. Makes the codebase unmaintainable

## Evidence of Violation
From actual transcript showing catastrophic failure:
- Branch 4: Created Registry interface with `Push(ctx, image string, content io.Reader) error`
- Branch 9: Created GiteaClient with INCOMPATIBLE `Push(ctx, image v1.Image, ref string, opts) error`
- Branch 11: Used GiteaClient's incompatible signature instead of Registry interface
- Result: 11 branches required complete rebuilding due to incompatible APIs

## Requirements

### MANDATORY Research Phase (BEFORE ANY PLANNING)
1. **Search for existing interfaces**:
   ```bash
   grep -r "type.*interface" --include="*.go" /path/to/effort-repo/
   grep -r "func.*Push\|func.*Upload\|func.*Store" --include="*.go" /path/to/effort-repo/
   ```

2. **Document all found abstractions**:
   - Interface names and signatures
   - Existing implementations
   - Package locations
   - Method signatures

3. **Check previous waves**:
   - Current wave implementations
   - Wave N-1 implementations
   - ALL previous phase implementations

### FORBIDDEN Actions
1. **NEVER create new interface when one exists**
2. **NEVER modify existing interface signatures**
3. **NEVER create "better" versions of existing APIs**
4. **NEVER implement competing abstractions**
5. **NEVER ignore existing implementations**

### REQUIRED Actions
1. **MUST use existing interfaces exactly**
2. **MUST implement all interface methods with exact signatures**
3. **MUST reuse existing utility functions**
4. **MUST import and use existing packages**
5. **MUST maintain API compatibility**

## Enforcement Mechanisms

### During Planning (Code Reviewer)
```markdown
## Existing Code Research (MANDATORY)
### Interfaces Found:
- Registry interface: Push(ctx, image string, content io.Reader) error
- [List all found interfaces]

### Implementations to Reuse:
- [List existing implementations]

### APIs to Implement:
- MUST implement Registry.Push with EXACT signature
- [List required implementations]

### FORBIDDEN:
- DO NOT create new Push method with different signature
- DO NOT create alternative Registry interface
```

### During Implementation (SW Engineer)
```go
// CORRECT: Implements existing interface
type GiteaClient struct{}

// MUST match Registry interface EXACTLY
func (g *GiteaClient) Push(ctx context.Context, image string, content io.Reader) error {
    // Implementation that conforms to existing interface
}

// WRONG: Creating incompatible method
func (g *GiteaClient) Push(ctx context.Context, image v1.Image, ref string, opts PushOptions) error {
    // THIS VIOLATES R373 - DIFFERENT SIGNATURE!
}
```

### During Review (Code Reviewer)
```bash
# Check for duplicate interfaces
grep -r "type.*interface.*{" --include="*.go" | sort | uniq -d

# Check for incompatible method signatures
grep -r "func.*Push" --include="*.go" | grep -v "Push(ctx context.Context, image string, content io.Reader)"

# Verify interface implementations
go test -run TestInterfaceCompliance ./...
```

### During Architecture Review (Architect)
- **REJECT** if duplicate interfaces found
- **REJECT** if incompatible implementations exist
- **REJECT** if existing code not reused
- **MANDATE** refactoring to use existing interfaces

## Example Violations and Fixes

### VIOLATION: Creating Competing Interface
```go
// Branch 4 already has:
type Registry interface {
    Push(ctx context.Context, image string, content io.Reader) error
}

// Branch 9 VIOLATES by creating:
type OCIRegistry interface {
    Push(ctx context.Context, img v1.Image, reference string, options ...Option) error
}
```

### FIX: Implement Existing Interface
```go
// Branch 9 MUST do:
import "branch4/registry"

type GiteaClient struct{}

// Implements the EXISTING Registry interface
func (g *GiteaClient) Push(ctx context.Context, image string, content io.Reader) error {
    // Convert parameters if needed internally
    // But maintain the existing interface contract
}
```

## Validation Scripts

### Pre-Planning Validation
```bash
#!/bin/bash
# R373-validate-no-duplicates.sh

echo "Checking for existing interfaces..."
INTERFACES=$(grep -r "type.*interface" --include="*.go" 2>/dev/null)

if [ ! -z "$INTERFACES" ]; then
    echo "FOUND EXISTING INTERFACES - MUST REUSE:"
    echo "$INTERFACES"
    echo ""
    echo "WARNING: Creating new interfaces when these exist violates R373"
fi

# Check for specific methods
for method in Push Upload Store Create Delete Get List; do
    echo "Checking for existing ${method} methods..."
    grep -r "func.*${method}(" --include="*.go" 2>/dev/null | head -5
done
```

### Post-Implementation Validation
```bash
#!/bin/bash
# R373-detect-violations.sh

# Find duplicate interface definitions
echo "Detecting duplicate interfaces..."
grep -r "type.*interface" --include="*.go" | \
    awk -F: '{print $2}' | \
    sort | uniq -d

# Find methods with same name but different signatures
echo "Detecting incompatible method signatures..."
grep -r "func.*Push(" --include="*.go" | \
    awk -F'Push' '{print $2}' | \
    sort | uniq -c | \
    awk '$1 > 1 {print "WARNING: Multiple Push signatures found"}'
```

## Grading Impact
- Creating duplicate interface: **-100% IMMEDIATE FAILURE**
- Creating incompatible implementation: **-100% IMMEDIATE FAILURE**
- Not researching existing code: **-50% penalty**
- Not documenting found interfaces: **-25% penalty**
- Not reusing existing utilities: **-30% penalty**

## Integration with Other Rules
- **R374**: Pre-Planning Research Protocol (companion rule)
- **R362**: No Architectural Rewrites (prevents changing existing interfaces)
- **R220**: Atomic PR Design (ensures compatible implementations)
- **R009**: Wave Integration (validates interface compliance during integration)

## State Machine Integration
- **EFFORT_PLAN_CREATION**: MUST include interface research
- **IMPLEMENTATION**: MUST use existing interfaces
- **CODE_REVIEW**: MUST verify no duplicates
- **WAVE_REVIEW**: MUST confirm interface compliance
- **INTEGRATION**: MUST validate API compatibility

## Recovery Protocol
If duplicate implementations are discovered:
1. **STOP all work immediately**
2. **Identify the canonical interface**
3. **Refactor ALL branches to use canonical version**
4. **Update all dependent code**
5. **Revalidate entire wave**

## Tags
#supreme-law #blocking #interface-compliance #code-reuse #architecture #r373