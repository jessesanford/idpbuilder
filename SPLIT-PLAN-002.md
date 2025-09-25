# Split Plan 002 - OCI Client Core

## 🔴 MINDSET REMINDER FOR SW ENGINEER 🔴
**This is a PARTIAL implementation. It SHOULD feel incomplete.**
**Your job is to stay within budget, NOT to complete the feature.**
**If you think "this needs X to work properly" - STOP. X is probably in another split.**

## Split Metadata
- **Split Number**: 002
- **Parent Effort**: client-interface-tests
- **Original Branch**: idpbuilderpush/phase3/wave1/client-interface-tests
- **Target Size**: 600 lines (max 800)
- **Created**: 2025-09-25 02:58:56

## Infrastructure Requirements
- **Directory Name**: client-interface-tests-SPLIT-002
- **Location**: efforts/phase3/wave1/client-interface-tests-SPLIT-002/
- **Clone Required**: Yes - separate clone of target repository
- **Branch Base**: idpbuilderpush/phase3/wave1/client-interface-tests-split-001 (from Split-001!)
- **Branch Name**: idpbuilderpush/phase3/wave1/client-interface-tests-split-002

## 🔴🔴🔴 CRITICAL: FILE PLACEMENT (R326) 🔴🔴🔴

**⚠️⚠️⚠️ DO NOT CREATE split-002/ SUBDIRECTORY! ⚠️⚠️⚠️**

Files MUST go directly in standard project directories:
- ✅ CORRECT: pkg/oci/auth.go
- ❌ WRONG: split-002/pkg/oci/auth.go

Creating split subdirectories causes CATASTROPHIC measurement errors!
Your working directory is already split-specific: efforts/.../client-interface-tests-SPLIT-002/

## Implementation Scope

### 🚨 EXPLICIT SCOPE DEFINITION (MANDATORY PER R310)

#### MINIMUM VIABLE SCOPE (Your Exact Contract)

**FILES TO CREATE/MODIFY (COMPLETE LIST):**
```
1. pkg/oci/types.go (CREATE) - 65 lines MAX
2. pkg/oci/errors.go (CREATE) - 35 lines MAX
3. pkg/oci/auth.go (CREATE) - 335 lines MAX
4. pkg/oci/flow.go (CREATE) - 151 lines MAX
TOTAL: 4 files, ~586 lines (implementing any other file = AUTOMATIC FAILURE)
```

**TYPES/STRUCTS TO IMPLEMENT:**
```go
// In types.go:
type RegistryAuth struct {
    Username string
    Password string
    Token    string
    // Essential fields only
}

type PushOptions struct {
    ImageRef string
    Auth     *RegistryAuth
    Insecure bool
    // Essential fields only
}

// In errors.go:
type RegistryError struct {
    Code    int
    Message string
}

// In auth.go:
type Authenticator interface {
    // Define interface methods
}

type BasicAuthenticator struct {
    // Implementation fields
}

type TokenAuthenticator struct {
    // Implementation fields
}

// In flow.go:
type PushFlow struct {
    // Flow control fields
}
```

**FUNCTIONS TO IMPLEMENT (BY EXACT SIGNATURE):**
```go
// In errors.go:
func (e *RegistryError) Error() string // 5 lines
func NewRegistryError(code int, msg string) *RegistryError // 10 lines

// In auth.go:
func NewBasicAuthenticator(username, password string) Authenticator // 15 lines
func NewTokenAuthenticator(token string) Authenticator // 15 lines
func (a *BasicAuthenticator) Authenticate(req *http.Request) error // 25 lines
func (a *TokenAuthenticator) Authenticate(req *http.Request) error // 20 lines
func DetectAuthScheme(endpoint string) (string, error) // 50 lines

// In flow.go:
func NewPushFlow(opts *PushOptions) (*PushFlow, error) // 30 lines
func (f *PushFlow) Execute() error // 80 lines
func (f *PushFlow) validateOptions() error // 25 lines
```

**TESTS TO WRITE:**
```
ZERO tests in this split (tests come in Split-004)
```

### 🚨🚨🚨 R355 PRODUCTION READY REQUIREMENTS (SUPREME LAW) 🚨🚨🚨

**ALL CODE MUST BE PRODUCTION READY - NO EXCEPTIONS**

#### Configuration Examples for This Split:
```go
// ❌ WRONG - Hardcoded values (AUTOMATIC FAILURE)
const RegistryURL = "https://registry.example.com"
const RetryCount = 3

// ✅ CORRECT - Configuration-driven
func getRegistryURL() string {
    url := os.Getenv("REGISTRY_URL")
    if url == "" {
        url = "https://registry-1.docker.io" // Default ONLY
    }
    return url
}

// ❌ WRONG - Stub implementation (AUTOMATIC FAILURE)
func (f *PushFlow) Execute() error {
    // TODO: implement push
    return nil
}

// ✅ CORRECT - Complete implementation (even if simple)
func (f *PushFlow) Execute() error {
    if err := f.validateOptions(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Even basic implementation is better than stub
    auth := f.options.Auth
    if auth == nil {
        return NewRegistryError(401, "authentication required")
    }

    // Actual push logic (simplified but complete)
    return f.performPush()
}
```

### 🛑 STOP BOUNDARIES - DO NOT IMPLEMENT
**EXPLICITLY FORBIDDEN IN THIS SPLIT:**
- ❌ DO NOT add test files (*_test.go)
- ❌ DO NOT add mock implementations
- ❌ DO NOT add testutil or testdata packages
- ❌ DO NOT implement the K8s client
- ❌ DO NOT implement Kind integration
- ❌ DO NOT implement the push command
- ❌ DO NOT add complex retry logic (basic only)
- ❌ DO NOT add comprehensive logging (basic only)
- ❌ DO NOT leave ANY TODOs or stubs (R355)
- ❌ DO NOT hardcode ANY values (R355)

## Technical Requirements

### Dependencies
- External dependencies:
  - Standard library (net/http, fmt, errors, etc.)
  - github.com/google/go-containerregistry v0.16.1+ (for OCI operations)
- From previous splits:
  - API types from Split-001 (import as needed)

### Interfaces to Provide
- Authenticator interface for registry authentication
- PushFlow for OCI push operations
- RegistryError for error handling

### Interfaces to Consume
- None directly (may use types from Split-001 if needed)

## Implementation Instructions

### Step 1: Setup
1. Verify you're in the split directory: `efforts/phase3/wave1/client-interface-tests-SPLIT-002/`
2. Create branch: `idpbuilderpush/phase3/wave1/client-interface-tests-split-002`
3. Branch from `idpbuilderpush/phase3/wave1/client-interface-tests-split-001`

### Step 2: Implementation
1. Create the pkg/oci directory
2. Implement types.go with basic types (no methods beyond essentials)
3. Implement errors.go with RegistryError type
4. Implement auth.go with authentication logic
5. Implement flow.go with push flow control

### Step 3: Testing
- NO tests in this split
- Ensure code compiles: `go build ./...`

### Step 4: Integration
- Commit all changes
- Push branch
- Measure with line-counter.sh to verify under 600 lines

## Size Management with REALISTIC Calculations
- Target: 600 lines (MAX 600 to leave buffer)
- Hard Stop: 700 lines (better incomplete than oversized)
- Measurement: Use line-counter.sh before committing

### Realistic Line Estimates:
```
types.go: ~65 lines (structs + minimal docs)
errors.go: ~35 lines (error type + methods)
auth.go: ~335 lines (interface + 2 implementations + detection)
flow.go: ~151 lines (flow control + validation)
TOTAL: ~586 lines (within 600 limit)
```

## Success Criteria
- [ ] All 4 specified files implemented
- [ ] Size under 600 lines (measured with line-counter.sh)
- [ ] Code compiles without errors
- [ ] Complete implementations (no TODOs or stubs)
- [ ] No hardcoded values (R355)
- [ ] Authentication logic works

## Notes for SW Engineer
- This split focuses on OCI client core functionality
- Keep implementations simple but complete
- Don't add extra features or optimizations
- Basic error handling is sufficient
- Remember: You have Split-001's types available

### 🔴 ADHERENCE REMINDER (R310):
- Implement EXACTLY what's listed - no more
- If it seems incomplete, that's intentional
- Do NOT add "helpful" extras
- Do NOT "complete" the implementation
- STOP at the boundaries specified above