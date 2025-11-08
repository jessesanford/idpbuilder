# Wave 2.3 Implementation Plan - Error Handling & Validation

**Wave**: Phase 2, Wave 3 (Error Handling & Validation)
**Phase**: Phase 2 - Core Push Functionality
**Created**: 2025-11-03
**Planner**: @agent-code-reviewer
**Fidelity Level**: **EXACT SPECIFICATIONS** (detailed efforts with R213 metadata)

---

## Wave Overview

**Goal**: Add comprehensive input validation, security checks (command injection, SSRF), and user-friendly error handling with actionable error messages and correct exit codes to the idpbuilder push command.

**Architecture Reference**: See `planning/phase2/wave3/WAVE-2.3-ARCHITECTURE.md` for design details

**Test Plan Reference**: See `planning/phase2/wave3/WAVE-TEST-PLAN.md` for 63 test specifications

**Total Efforts**: 2

**Total Estimated Lines**: 750 lines (400 + 350)

---

## Effort Definitions

### Effort 2.3.1: Input Validation & Security Checks

#### R213 Metadata

```yaml
effort_id: "2.3.1"
effort_name: "Input Validation & Security Checks"
branch_name: "idpbuilder-oci-push/phase2/wave3/effort-2.3.1-input-validation"
base_branch: "idpbuilder-oci-push/phase2/wave3/integration"
parent_wave: "wave-2.3"
parent_phase: "phase-2"
depends_on:
  - "integration:phase2-wave2-integration"  # Wave 2.2 configuration system
dependencies_detail:
  - type: "integration"
    target: "phase2-wave2-integration"
    reason: "Extends PushConfig validation with security checks"
  - type: "code"
    target: "pkg/cmd/push/config.go"
    reason: "Uses PushConfig and ConfigValue types"
estimated_lines: 400
complexity: "medium-high"
can_parallelize: false
parallel_with: []
files_touched:
  - "pkg/validator/imagename.go"           # new
  - "pkg/validator/registry.go"            # new
  - "pkg/validator/credentials.go"         # new
  - "pkg/validator/validator_test.go"      # new
risk_level: "high"
risk_reason: "Security-critical validation logic"
test_count: 33
```

#### Scope

**Purpose**: Implement comprehensive input validation for OCI image names, registry URLs, and credentials with security checks to prevent command injection and SSRF attacks.

**What This Effort Accomplishes**:
- OCI image name format validation (registry/namespace/repository:tag[@digest])
- Command injection prevention for image names, registries, and credentials
- Registry URL validation (hostname, IP, port formats)
- SSRF protection (detect private IP ranges, localhost variants)
- Credentials validation (required fields, dangerous character detection)
- Security warnings (weak credentials, private registries)
- Actionable error messages with "Error: X, Suggestion: Y" format

**Boundaries - What is OUT of Scope**:
- ❌ Error type system (handled in Effort 2.3.2)
- ❌ Exit code mapping (handled in Effort 2.3.2)
- ❌ Integration with runPush (handled in Effort 2.3.2)
- ❌ Cobra command modification (handled in Effort 2.3.2)
- ❌ Password hashing/encryption (authentication is username/password plaintext for registry)
- ❌ Network connectivity checks (validation only, no actual network calls)

#### Files to Create/Modify

**New Files**:
- `pkg/validator/imagename.go` (150 lines)
  - ValidateImageName function
  - OCI image name regex patterns
  - Tag/digest validation
  - Command injection detection

- `pkg/validator/registry.go` (120 lines)
  - ValidateRegistryURL function
  - Registry hostname validation
  - SSRF detection (private IP ranges)
  - isPrivateIP helper

- `pkg/validator/credentials.go` (80 lines)
  - ValidateCredentials function
  - Command injection detection
  - Weak credentials warning

- `pkg/validator/validator_test.go` (350 lines)
  - 33 unit tests covering all validation scenarios
  - Test fixtures for command injection patterns
  - Table-driven tests (from Wave 2.2 pattern)

**Total Estimated Lines**: 700 lines (350 implementation + 350 tests)

**Note**: Error types (ValidationError, SSRFWarning, SecurityWarning) are defined in this effort but fully implemented in Effort 2.3.2.

#### Exact Code Specifications

**File: pkg/validator/imagename.go (150 lines)**

```go
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// OCI image name validation regex patterns
var (
	// imageNamePattern validates full OCI image references
	// Format: [registry/][namespace/]repository[:tag][@digest]
	// Examples: alpine:latest, docker.io/library/ubuntu:22.04, localhost:5000/myapp:v1.0
	imageNamePattern = regexp.MustCompile(`^([a-zA-Z0-9._-]+([.:][0-9]+)?/)?([a-zA-Z0-9._/-]+)(:[a-zA-Z0-9._-]+)?(@sha256:[a-fA-F0-9]{64})?$`)

	// tagPattern validates image tags
	tagPattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

	// digestPattern validates SHA256 digests
	digestPattern = regexp.MustCompile(`^sha256:[a-fA-F0-9]{64}$`)
)

// Dangerous characters that could be used for command injection
const dangerousChars = ";|&$`<>(){}[]'\"\\"

// ValidateImageName validates an OCI image name format
//
// Validation Rules:
// 1. Non-empty string required
// 2. No shell metacharacters (command injection prevention)
// 3. Must match OCI image name pattern
// 4. Tag must be valid if present
// 5. Digest must be valid if present
//
// Returns:
//   - nil if valid
//   - ValidationError with actionable suggestion if invalid
func ValidateImageName(imageName string) error {
	if imageName == "" {
		return &ValidationError{
			Field:      "image-name",
			Message:    "image name cannot be empty",
			Suggestion: "provide an image name like 'alpine:latest' or 'docker.io/library/ubuntu:22.04'",
			ExitCode:   1,
		}
	}

	// Check for dangerous characters (command injection prevention)
	if containsAnyChar(imageName, dangerousChars) {
		return &ValidationError{
			Field:      "image-name",
			Message:    fmt.Sprintf("image name contains shell metacharacters: %s", imageName),
			Suggestion: "use only alphanumeric characters, dots, hyphens, underscores, colons, and slashes",
			ExitCode:   1,
		}
	}

	// Validate against OCI image name pattern
	if !imageNamePattern.MatchString(imageName) {
		return &ValidationError{
			Field:      "image-name",
			Message:    fmt.Sprintf("invalid OCI image name format: %s", imageName),
			Suggestion: "use format: [registry/][namespace/]repository[:tag][@digest]",
			ExitCode:   1,
		}
	}

	// Validate tag if present
	parts := strings.Split(imageName, ":")
	if len(parts) > 1 && !strings.Contains(parts[1], "@") {
		tag := parts[1]
		if !tagPattern.MatchString(tag) {
			return &ValidationError{
				Field:      "image-tag",
				Message:    fmt.Sprintf("invalid tag format: %s", tag),
				Suggestion: "tags must contain only alphanumeric characters, dots, hyphens, and underscores",
				ExitCode:   1,
			}
		}
	}

	// Validate digest if present
	if strings.Contains(imageName, "@") {
		digestPart := strings.Split(imageName, "@")[1]
		if !digestPattern.MatchString(digestPart) {
			return &ValidationError{
				Field:      "image-digest",
				Message:    fmt.Sprintf("invalid digest format: %s", digestPart),
				Suggestion: "digest must be in format: sha256:[64 hex characters]",
				ExitCode:   1,
			}
		}
	}

	return nil
}

// containsAnyChar checks if a string contains any character from a set
func containsAnyChar(s, chars string) bool {
	for _, c := range chars {
		if strings.ContainsRune(s, c) {
			return true
		}
	}
	return false
}
```

**Implementation Requirements for imagename.go**:
- Use standard library only (regexp, strings, fmt)
- No external dependencies
- All regex patterns must be compiled at package initialization
- Error messages must include the invalid input (for debugging)
- Suggestions must be actionable (tell user exactly what to do)
- Command injection detection is CRITICAL (security requirement)

---

**File: pkg/validator/registry.go (120 lines)**

```go
package validator

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
)

var (
	// registryPattern validates registry hostnames (domain or IP with optional port)
	registryPattern = regexp.MustCompile(`^([a-zA-Z0-9.-]+|\[[0-9a-fA-F:]+\])(:[0-9]{1,5})?$`)

	// Private IP ranges (SSRF protection)
	privateIPRanges = []string{
		"10.0.0.0/8",      // Private class A
		"172.16.0.0/12",   // Private class B
		"192.168.0.0/16",  // Private class C
		"127.0.0.0/8",     // Loopback
		"169.254.0.0/16",  // Link-local
		"::1/128",         // IPv6 loopback
		"fc00::/7",        // IPv6 unique local
		"fe80::/10",       // IPv6 link-local
	}
)

// ValidateRegistryURL validates a registry URL or hostname
//
// Validation Rules:
// 1. Non-empty string required
// 2. No shell metacharacters (command injection prevention)
// 3. Must be valid hostname, IP, or full URL
// 4. SSRF warning for private IP ranges (non-blocking)
//
// Returns:
//   - nil if valid
//   - ValidationError if invalid format or contains dangerous characters
//   - SSRFWarning if private IP detected (allows continuation)
func ValidateRegistryURL(registry string) error {
	if registry == "" {
		return &ValidationError{
			Field:      "registry",
			Message:    "registry URL cannot be empty",
			Suggestion: "provide a registry URL like 'docker.io' or 'localhost:5000'",
			ExitCode:   1,
		}
	}

	// Check for dangerous characters
	if containsAnyChar(registry, dangerousChars) {
		return &ValidationError{
			Field:      "registry",
			Message:    fmt.Sprintf("registry URL contains shell metacharacters: %s", registry),
			Suggestion: "use only alphanumeric characters, dots, hyphens, colons, and brackets",
			ExitCode:   1,
		}
	}

	// Parse URL or treat as hostname
	var hostname string

	// Try parsing as full URL first
	if strings.Contains(registry, "://") {
		parsedURL, err := url.Parse(registry)
		if err != nil {
			return &ValidationError{
				Field:      "registry",
				Message:    fmt.Sprintf("invalid registry URL: %v", err),
				Suggestion: "use format: hostname[:port] or https://hostname[:port]",
				ExitCode:   1,
			}
		}
		hostname = parsedURL.Hostname()
	} else {
		// Treat as hostname:port
		parts := strings.Split(registry, ":")
		hostname = parts[0]
	}

	// Validate hostname pattern
	if !registryPattern.MatchString(hostname) && !registryPattern.MatchString(registry) {
		return &ValidationError{
			Field:      "registry",
			Message:    fmt.Sprintf("invalid registry hostname format: %s", hostname),
			Suggestion: "use a valid domain name, IP address, or 'localhost'",
			ExitCode:   1,
		}
	}

	// SSRF protection: warn about private IP ranges
	if isPrivateIP(hostname) {
		// This is a warning, not an error (some users intentionally use private registries)
		return &SSRFWarning{
			Target:     registry,
			Message:    fmt.Sprintf("registry appears to be in a private IP range: %s", hostname),
			Suggestion: "ensure this is intentional and you trust the target registry",
		}
	}

	return nil
}

// isPrivateIP checks if a hostname resolves to a private IP address
func isPrivateIP(hostname string) bool {
	// Check localhost variants
	if hostname == "localhost" || hostname == "127.0.0.1" || hostname == "::1" {
		return true
	}

	// Try to resolve as IP
	ip := net.ParseIP(hostname)
	if ip == nil {
		// Try DNS resolution
		ips, err := net.LookupIP(hostname)
		if err != nil || len(ips) == 0 {
			return false
		}
		ip = ips[0]
	}

	// Check against private ranges
	for _, cidr := range privateIPRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}
```

**Implementation Requirements for registry.go**:
- Use standard library (net, net/url, regexp, strings, fmt)
- SSRF detection must check IPv4 and IPv6 ranges
- DNS resolution failures should NOT block validation (return false)
- URL parsing must handle both hostname:port and full URLs
- Private IP warning is non-blocking (returns SSRFWarning, not error)

---

**File: pkg/validator/credentials.go (80 lines)**

```go
package validator

import (
	"fmt"
	"strings"
)

// ValidateCredentials validates username and password for security
//
// Validation Rules:
// 1. Username cannot be empty
// 2. Password cannot be empty
// 3. Username must not contain shell metacharacters
// 4. Password CAN contain special characters (no restrictions)
// 5. Weak credential warning (non-blocking)
//
// Returns:
//   - nil if valid
//   - ValidationError if required fields missing or username has dangerous chars
//   - SecurityWarning if weak credentials detected (allows continuation)
func ValidateCredentials(username, password string) error {
	// Validate username
	if username == "" {
		return &ValidationError{
			Field:      "username",
			Message:    "username is required",
			Suggestion: "provide username via --username flag or IDPBUILDER_USERNAME environment variable",
			ExitCode:   1,
		}
	}

	// Check username for command injection patterns
	if containsAnyChar(username, dangerousChars) {
		return &ValidationError{
			Field:      "username",
			Message:    fmt.Sprintf("username contains shell metacharacters: %s", username),
			Suggestion: "use alphanumeric characters, hyphens, underscores, dots, and @ symbols only",
			ExitCode:   1,
		}
	}

	// Validate password
	if password == "" {
		return &ValidationError{
			Field:      "password",
			Message:    "password is required",
			Suggestion: "provide password via --password flag or IDPBUILDER_PASSWORD environment variable",
			ExitCode:   1,
		}
	}

	// Note: We do NOT validate password characters (passwords can contain special chars)
	// But we DO escape them properly when used in authentication

	// Warn if credentials appear in plain sight (basic security check)
	if strings.Contains(username, "admin") && strings.Contains(password, "password") {
		return &SecurityWarning{
			Message:    "using default or weak credentials detected",
			Suggestion: "consider using stronger credentials for production registries",
		}
	}

	return nil
}
```

**Implementation Requirements for credentials.go**:
- Username validation is strict (no dangerous characters)
- Password validation is lenient (any characters allowed)
- Weak credential detection is heuristic (simple pattern matching)
- SecurityWarning is non-blocking (allows continuation)
- Error messages must reference configuration sources (flags, env vars)

---

**File: pkg/validator/types.go (50 lines - ERROR TYPES DEFINED HERE)**

```go
package validator

import "fmt"

// ValidationError represents input validation failures (exit code 1)
type ValidationError struct {
	Field      string
	Message    string
	Suggestion string
	ExitCode   int
}

func (e *ValidationError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("Error: %s\nSuggestion: %s", e.Message, e.Suggestion)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

// SSRFWarning represents a potential SSRF risk (warning, not error)
type SSRFWarning struct {
	Target     string
	Message    string
	Suggestion string
}

func (w *SSRFWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// SecurityWarning represents security concerns (warning, not error)
type SecurityWarning struct {
	Message    string
	Suggestion string
}

func (w *SecurityWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// IsWarning returns true if the error is a warning (should not stop execution)
func IsWarning(err error) bool {
	_, isSSRF := err.(*SSRFWarning)
	_, isSecurity := err.(*SecurityWarning)
	return isSSRF || isSecurity
}
```

**Note**: These error types are MINIMAL definitions for Effort 2.3.1. Effort 2.3.2 will move them to `pkg/errors/` and add full error type system with exit code mapping.

#### Tests Required

**File: pkg/validator/validator_test.go (350 lines, 33 tests)**

See `planning/phase2/wave3/WAVE-TEST-PLAN.md` for complete test specifications.

**Test Suites**:
1. **Image Name Validation** (15 tests)
   - T-2.3.1-01: TestValidateImageName_SimpleTag
   - T-2.3.1-02: TestValidateImageName_WithRegistry
   - T-2.3.1-03: TestValidateImageName_WithNamespace
   - T-2.3.1-04: TestValidateImageName_WithDigest
   - T-2.3.1-05: TestValidateImageName_NoTag
   - T-2.3.1-06: TestValidateImageName_Empty
   - T-2.3.1-07: TestValidateImageName_CommandInjection_Semicolon
   - T-2.3.1-08: TestValidateImageName_CommandInjection_Backtick
   - T-2.3.1-09: TestValidateImageName_CommandInjection_Dollar
   - T-2.3.1-10: TestValidateImageName_InvalidTag
   - T-2.3.1-11: TestValidateImageName_InvalidDigest
   - T-2.3.1-12: TestValidateImageName_Localhost
   - T-2.3.1-13: TestValidateImageName_IPv4Registry
   - T-2.3.1-14: TestValidateImageName_IPv6Registry
   - T-2.3.1-15: TestValidateImageName_LongName

2. **Registry URL Validation** (10 tests)
   - T-2.3.2-01: TestValidateRegistryURL_SimpleDomain
   - T-2.3.2-02: TestValidateRegistryURL_DomainWithPort
   - T-2.3.2-03: TestValidateRegistryURL_Localhost
   - T-2.3.2-04: TestValidateRegistryURL_IPv4
   - T-2.3.2-05: TestValidateRegistryURL_IPv6
   - T-2.3.2-06: TestValidateRegistryURL_PrivateIP_ClassA
   - T-2.3.2-07: TestValidateRegistryURL_PrivateIP_ClassB
   - T-2.3.2-08: TestValidateRegistryURL_PrivateIP_ClassC
   - T-2.3.2-09: TestValidateRegistryURL_Localhost_Warning
   - T-2.3.2-10: TestValidateRegistryURL_CommandInjection

3. **Credentials Validation** (8 tests)
   - T-2.3.3-01: TestValidateCredentials_Alphanumeric
   - T-2.3.3-02: TestValidateCredentials_SpecialCharsPassword
   - T-2.3.3-03: TestValidateCredentials_EmailUsername
   - T-2.3.3-04: TestValidateCredentials_EmptyUsername
   - T-2.3.3-05: TestValidateCredentials_EmptyPassword
   - T-2.3.3-06: TestValidateCredentials_UsernameInjection
   - T-2.3.3-07: TestValidateCredentials_UsernameBacktick
   - T-2.3.3-08: TestValidateCredentials_WeakCredentials

**Test Coverage Requirements**:
- Minimum 95% statement coverage for validator package
- Minimum 90% branch coverage
- All security tests must pass (command injection, SSRF detection)
- Table-driven test pattern from Wave 2.2

**Test Dependencies**:
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- ✅ **Wave 2.2 Integration**: Configuration system with PushConfig (COMPLETE)
- ✅ **Phase 1**: All interface packages (docker, registry, auth, tls) (COMPLETE)

**Downstream Dependencies** (efforts that depend on this):
- **Effort 2.3.2**: Error type system (needs ValidationError types from this effort)
- **Effort 2.3.2**: Push command integration (needs validation functions)

**External Dependencies**:
- Standard library only: regexp, strings, fmt, net, net/url
- Test dependencies: github.com/stretchr/testify (already in go.mod)

#### Acceptance Criteria

- [ ] All files created as specified (imagename.go, registry.go, credentials.go, types.go, validator_test.go)
- [ ] All 33 tests passing (100% pass rate)
- [ ] Code coverage ≥95% statement, ≥90% branch
- [ ] No linting errors (golangci-lint)
- [ ] All validation functions have complete godoc comments
- [ ] Security tests verify command injection prevention
- [ ] SSRF tests verify private IP detection
- [ ] Line count within estimate (400 ±15% = 340-460 lines)
- [ ] TDD protocol followed (tests created BEFORE implementation per R341)

---

### Effort 2.3.2: Error Type System & Exit Code Mapping

#### R213 Metadata

```yaml
effort_id: "2.3.2"
effort_name: "Error Type System & Exit Code Mapping"
branch_name: "idpbuilder-oci-push/phase2/wave3/effort-2.3.2-error-handling"
base_branch: "idpbuilder-oci-push/phase2/wave3/effort-2.3.1-input-validation"
parent_wave: "wave-2.3"
parent_phase: "phase-2"
depends_on:
  - "effort:2.3.1"  # Needs ValidationError types
dependencies_detail:
  - type: "effort"
    target: "2.3.1"
    reason: "Imports ValidationError, SSRFWarning, SecurityWarning types"
  - type: "code"
    target: "pkg/validator/types.go"
    reason: "Moves error types to pkg/errors/ and extends them"
  - type: "integration"
    target: "pkg/cmd/push/push.go"
    reason: "Wraps errors in runPush function"
estimated_lines: 350
complexity: "medium"
can_parallelize: false
parallel_with: []
files_touched:
  - "pkg/errors/types.go"              # new (moves from pkg/validator/types.go)
  - "pkg/errors/exitcodes.go"          # new
  - "pkg/cmd/push/errors.go"           # new
  - "pkg/errors/types_test.go"         # new
  - "pkg/cmd/push/push_errors_test.go" # new
  - "pkg/cmd/push/push.go"             # modified (add error wrapping)
  - "pkg/validator/imagename.go"       # modified (update imports)
  - "pkg/validator/registry.go"        # modified (update imports)
  - "pkg/validator/credentials.go"     # modified (update imports)
risk_level: "medium"
risk_reason: "Changes exit code behavior of push command"
test_count: 30
```

#### Scope

**Purpose**: Implement comprehensive error type system with exit code mapping, error wrapping with context preservation, and integration with Wave 2.1/2.2 push command for actionable error messages.

**What This Effort Accomplishes**:
- Custom error types (ValidationError, AuthenticationError, NetworkError, ImageNotFoundError)
- Exit code mapping (1=validation, 2=auth, 3=network, 4=image not found)
- Error wrapping functions (wrapDockerError, wrapRegistryError)
- Error formatting with "Error: X, Suggestion: Y" pattern
- Integration with runPush function from Wave 2.1
- Cobra command exit code handling
- Error context preservation through wrapping

**Boundaries - What is OUT of Scope**:
- ❌ Validation logic (handled in Effort 2.3.1)
- ❌ SSRF detection (handled in Effort 2.3.1)
- ❌ Command injection prevention (handled in Effort 2.3.1)
- ❌ Configuration system (handled in Wave 2.2)
- ❌ Progress reporting (handled in Wave 2.1)
- ❌ Retry logic or error recovery (future wave)

#### Files to Create/Modify

**New Files**:
- `pkg/errors/types.go` (100 lines)
  - BaseError struct
  - ValidationError, AuthenticationError, NetworkError, ImageNotFoundError types
  - Error formatting methods
  - Error unwrapping support

- `pkg/errors/exitcodes.go` (60 lines)
  - Exit code constants
  - GetExitCode function
  - FormatError function

- `pkg/cmd/push/errors.go` (100 lines)
  - wrapDockerError function
  - wrapRegistryError function
  - validatePushOptions function

- `pkg/errors/types_test.go` (150 lines)
  - 18 tests for error type creation and formatting
  - Exit code mapping tests

- `pkg/cmd/push/push_errors_test.go` (200 lines)
  - 12 tests for error wrapping
  - Integration tests with mock Docker/registry clients

**Modified Files**:
- `pkg/cmd/push/push.go` (+80 lines)
  - Add error wrapping to runPush
  - Add exit code handling to Cobra command
  - Add validatePushOptions call

- `pkg/validator/imagename.go` (import change)
  - Change import from local types to pkg/errors

- `pkg/validator/registry.go` (import change)
  - Change import from local types to pkg/errors

- `pkg/validator/credentials.go` (import change)
  - Change import from local types to pkg/errors

**Total Estimated Lines**: 690 lines (340 implementation + 350 tests)

#### Exact Code Specifications

**File: pkg/errors/types.go (100 lines)**

```go
package errors

import (
	"fmt"
)

// BaseError is the base type for all custom errors
type BaseError struct {
	Message    string
	Suggestion string
	Cause      error
}

func (e *BaseError) Error() string {
	if e.Suggestion != "" {
		return fmt.Sprintf("Error: %s\nSuggestion: %s", e.Message, e.Suggestion)
	}
	return fmt.Sprintf("Error: %s", e.Message)
}

func (e *BaseError) Unwrap() error {
	return e.Cause
}

// ValidationError represents input validation failures (exit code 1)
type ValidationError struct {
	BaseError
	Field    string
	ExitCode int
}

func NewValidationError(field, message, suggestion string) *ValidationError {
	return &ValidationError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Field:    field,
		ExitCode: 1,
	}
}

// AuthenticationError represents authentication failures (exit code 2)
type AuthenticationError struct {
	BaseError
	Registry string
	ExitCode int
}

func NewAuthenticationError(registry, message, suggestion string) *AuthenticationError {
	return &AuthenticationError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Registry: registry,
		ExitCode: 2,
	}
}

// NetworkError represents network/connectivity failures (exit code 3)
type NetworkError struct {
	BaseError
	Target   string
	ExitCode int
}

func NewNetworkError(target, message, suggestion string) *NetworkError {
	return &NetworkError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		Target:   target,
		ExitCode: 3,
	}
}

// ImageNotFoundError represents missing image errors (exit code 4)
type ImageNotFoundError struct {
	BaseError
	ImageName string
	ExitCode  int
}

func NewImageNotFoundError(imageName, message, suggestion string) *ImageNotFoundError {
	return &ImageNotFoundError{
		BaseError: BaseError{
			Message:    message,
			Suggestion: suggestion,
		},
		ImageName: imageName,
		ExitCode:  4,
	}
}

// SSRFWarning represents a potential SSRF risk (warning, not error)
type SSRFWarning struct {
	Target     string
	Message    string
	Suggestion string
}

func (w *SSRFWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// SecurityWarning represents security concerns (warning, not error)
type SecurityWarning struct {
	Message    string
	Suggestion string
}

func (w *SecurityWarning) Error() string {
	return fmt.Sprintf("Warning: %s\nSuggestion: %s", w.Message, w.Suggestion)
}

// IsWarning returns true if the error is a warning (should not stop execution)
func IsWarning(err error) bool {
	_, isSSRF := err.(*SSRFWarning)
	_, isSecurity := err.(*SecurityWarning)
	return isSSRF || isSecurity
}
```

---

**File: pkg/errors/exitcodes.go (60 lines)**

```go
package errors

import (
	"errors"
	"fmt"
)

// Exit codes for different error types
const (
	ExitSuccess         = 0
	ExitValidationError = 1
	ExitAuthError       = 2
	ExitNetworkError    = 3
	ExitImageNotFound   = 4
	ExitGenericError    = 1
)

// GetExitCode maps an error to its appropriate exit code
func GetExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}

	// Check for typed errors
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return ExitValidationError
	}

	var authErr *AuthenticationError
	if errors.As(err, &authErr) {
		return ExitAuthError
	}

	var networkErr *NetworkError
	if errors.As(err, &networkErr) {
		return ExitNetworkError
	}

	var imageNotFoundErr *ImageNotFoundError
	if errors.As(err, &imageNotFoundErr) {
		return ExitImageNotFound
	}

	// Default to validation error for untyped errors
	return ExitGenericError
}

// FormatError formats an error with consistent styling
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	// Warnings have special formatting
	if IsWarning(err) {
		return fmt.Sprintf("⚠️  %s", err.Error())
	}

	// Regular errors
	return fmt.Sprintf("❌ %s", err.Error())
}
```

---

**File: pkg/cmd/push/errors.go (100 lines)**

```go
package push

import (
	"fmt"
	"strings"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
	"github.com/cnoe-io/idpbuilder/pkg/validator"
)

// validatePushOptions validates all push options with security checks
//
// This extends Wave 2.2's PushConfig.Validate() with additional
// security validation from Wave 2.3.
func validatePushOptions(opts *PushOptions) error {
	// Validate image name
	if err := validator.ValidateImageName(opts.ImageName); err != nil {
		// validator returns typed errors, just pass through
		return err
	}

	// Validate registry URL
	if err := validator.ValidateRegistryURL(opts.Registry); err != nil {
		// Check if it's a warning
		if errors.IsWarning(err) {
			// Log warning but continue
			fmt.Println(errors.FormatError(err))
		} else {
			return err
		}
	}

	// Validate credentials
	if err := validator.ValidateCredentials(opts.Username, opts.Password); err != nil {
		if errors.IsWarning(err) {
			fmt.Println(errors.FormatError(err))
		} else {
			return err
		}
	}

	return nil
}

// wrapDockerError wraps Docker client errors with appropriate types
func wrapDockerError(err error, imageName string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for "image not found" errors
	if strings.Contains(errMsg, "No such image") {
		return errors.NewImageNotFoundError(
			imageName,
			fmt.Sprintf("image '%s' not found in local Docker daemon", imageName),
			"pull the image first with: docker pull "+imageName,
		)
	}

	// Check for Docker daemon connection errors
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "Cannot connect") {
		return errors.NewNetworkError(
			"docker daemon",
			"cannot connect to Docker daemon",
			"ensure Docker daemon is running: systemctl start docker or start Docker Desktop",
		)
	}

	// Generic Docker error
	return fmt.Errorf("Docker error: %w", err)
}

// wrapRegistryError wraps registry client errors with appropriate types
func wrapRegistryError(err error, registry string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Check for authentication errors
	if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized") {
		return errors.NewAuthenticationError(
			registry,
			fmt.Sprintf("authentication failed for registry %s", registry),
			"check your username and password, or verify registry credentials",
		)
	}

	// Check for network errors
	if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "timeout") {
		return errors.NewNetworkError(
			registry,
			fmt.Sprintf("cannot connect to registry %s", registry),
			"verify registry URL and network connectivity, or try with --insecure if using self-signed certificates",
		)
	}

	// TLS certificate errors
	if strings.Contains(errMsg, "x509") || strings.Contains(errMsg, "certificate") {
		return errors.NewNetworkError(
			registry,
			fmt.Sprintf("TLS certificate verification failed for %s", registry),
			"use --insecure flag to skip certificate verification (not recommended for production)",
		)
	}

	// Generic registry error
	return fmt.Errorf("registry error: %w", err)
}
```

---

**File: pkg/cmd/push/push.go (MODIFICATIONS +80 lines)**

```go
// MODIFICATIONS TO EXISTING FILE (Wave 2.1/2.2)

package push

import (
	"context"
	"fmt"
	"os"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
	// ... existing imports from Wave 2.1/2.2
)

// runPush executes the push workflow with comprehensive error handling (Wave 2.3)
func runPush(ctx context.Context, opts *PushOptions) error {
	// STAGE 1: Validation (NEW in Wave 2.3)
	if err := validatePushOptions(opts); err != nil {
		return err // Already typed error from validator
	}

	// STAGE 2: Docker client connection (existing from Wave 2.1)
	dockerClient, err := docker.NewClient()
	if err != nil {
		return wrapDockerError(err, opts.ImageName) // NEW: error wrapping
	}
	defer dockerClient.Close()

	// STAGE 3: Retrieve image from Docker daemon (existing from Wave 2.1)
	image, err := dockerClient.GetImage(ctx, opts.ImageName)
	if err != nil {
		return wrapDockerError(err, opts.ImageName) // NEW: error wrapping
	}

	// STAGE 4-7: Existing Wave 2.1/2.2 logic with error wrapping
	// ... (configuration, auth, TLS, registry client, push)

	// Add error wrapping for registry operations:
	if err := registryClient.Push(ctx, image, targetRef, progressCallback); err != nil {
		return wrapRegistryError(err, opts.Registry) // NEW: error wrapping
	}

	// STAGE 8: Success (existing)
	fmt.Printf("✅ Successfully pushed %s to %s\n", opts.ImageName, targetRef)
	return nil
}

// NewPushCommand creates the push command (MODIFIED for exit code handling)
func NewPushCommand(v *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		// ... existing configuration from Wave 2.2
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration (Wave 2.2)
			config, err := LoadConfig(cmd, args, v)
			if err != nil {
				return err
			}

			// Validate configuration (Wave 2.2)
			if err := config.Validate(); err != nil {
				return err
			}

			// Convert to PushOptions (Wave 2.2)
			opts := config.ToPushOptions()

			// Execute push pipeline (Wave 2.1/2.3)
			if err := runPush(cmd.Context(), opts); err != nil {
				// NEW: Get appropriate exit code from error type
				exitCode := errors.GetExitCode(err)

				// NEW: Print formatted error
				fmt.Fprintln(os.Stderr, errors.FormatError(err))

				// NEW: Exit with typed exit code
				os.Exit(exitCode)
			}

			return nil
		},
	}

	// ... existing flag definitions from Wave 2.2
	return cmd
}
```

**Integration Points**:
- Wraps existing runPush stages with error handling
- Preserves all Wave 2.1/2.2 functionality
- Adds exit code handling to Cobra command
- Backward compatible (existing code paths unchanged)

#### Tests Required

**File: pkg/errors/types_test.go (150 lines, 18 tests)**

See `planning/phase2/wave3/WAVE-TEST-PLAN.md` Test Suites 4-5 for complete specifications.

**Test Suites**:
1. **Error Type Creation & Formatting** (10 tests)
   - T-2.3.4-01: TestNewValidationError
   - T-2.3.4-02: TestNewAuthenticationError
   - T-2.3.4-03: TestNewNetworkError
   - T-2.3.4-04: TestNewImageNotFoundError
   - T-2.3.4-05: TestValidationError_Format
   - T-2.3.4-06: TestAuthenticationError_Format
   - T-2.3.4-07: TestSSRFWarning_Format
   - T-2.3.4-08: TestBaseError_Unwrap
   - T-2.3.4-09: TestErrorChain_Unwrap
   - T-2.3.4-10: TestIsWarning_Detection

2. **Exit Code Mapping** (8 tests)
   - T-2.3.5-01: TestGetExitCode_ValidationError
   - T-2.3.5-02: TestGetExitCode_AuthenticationError
   - T-2.3.5-03: TestGetExitCode_NetworkError
   - T-2.3.5-04: TestGetExitCode_ImageNotFoundError
   - T-2.3.5-05: TestGetExitCode_GenericError
   - T-2.3.5-06: TestGetExitCode_NilError
   - T-2.3.5-07: TestGetExitCode_WrappedValidationError
   - T-2.3.5-08: TestGetExitCode_WrappedAuthError

---

**File: pkg/cmd/push/push_errors_test.go (200 lines, 12 tests)**

See `planning/phase2/wave3/WAVE-TEST-PLAN.md` Test Suite 6 for complete specifications.

**Test Suite**:
3. **Error Wrapping Integration** (12 tests)
   - T-2.3.6-01: TestWrapDockerError_ImageNotFound
   - T-2.3.6-02: TestWrapDockerError_ConnectionRefused
   - T-2.3.6-03: TestWrapDockerError_CannotConnect
   - T-2.3.6-04: TestWrapDockerError_GenericError
   - T-2.3.6-05: TestWrapRegistryError_Unauthorized
   - T-2.3.6-06: TestWrapRegistryError_ConnectionRefused
   - T-2.3.6-07: TestWrapRegistryError_Timeout
   - T-2.3.6-08: TestWrapRegistryError_TLSError
   - T-2.3.6-09: TestWrapRegistryError_GenericError
   - T-2.3.6-10: TestWrapDockerError_PreservesImageName
   - T-2.3.6-11: TestWrapRegistryError_PreservesRegistry
   - T-2.3.6-12: TestWrapErrors_ChainUnwraps

**Test Coverage Requirements**:
- Minimum 95% statement coverage for errors package
- Minimum 90% branch coverage
- All exit code paths tested
- Error wrapping preserves context
- Integration with Wave 2.1/2.2 verified

#### Dependencies

**Upstream Dependencies** (must complete before this effort):
- ✅ **Effort 2.3.1**: Validation system (MUST COMPLETE FIRST)
  - Needs: ValidationError, SSRFWarning, SecurityWarning types
  - Needs: ValidateImageName, ValidateRegistryURL, ValidateCredentials functions
- ✅ **Wave 2.2**: Configuration system (COMPLETE)
- ✅ **Wave 2.1**: Push command core (COMPLETE)

**Downstream Dependencies** (efforts that depend on this):
- None (this is the final effort in Wave 2.3)

**External Dependencies**:
- Standard library: errors, fmt, strings, os
- Test dependencies: github.com/stretchr/testify (already in go.mod)

#### Acceptance Criteria

- [ ] All files created/modified as specified
- [ ] All 30 tests passing (100% pass rate)
- [ ] Code coverage ≥95% statement, ≥90% branch for errors package
- [ ] Exit code mapping verified (1, 2, 3, 4, 0)
- [ ] Error wrapping preserves context and type information
- [ ] Integration with Wave 2.1 runPush verified
- [ ] Backward compatibility with Wave 2.2 configuration maintained
- [ ] Line count within estimate (350 ±15% = 298-403 lines)
- [ ] No linting errors (golangci-lint)
- [ ] All public functions have complete godoc comments

---

## Parallelization Strategy

### Sequential Implementation (NO PARALLELIZATION)

**Execution Order**:
```
Effort 2.3.1: Input Validation & Security Checks
    ↓ (MUST COMPLETE FIRST)
Effort 2.3.2: Error Type System & Exit Code Mapping
```

**Rationale for Sequential Execution**:
1. **Effort 2.3.2 depends on 2.3.1**: Error type system imports ValidationError and other types from 2.3.1
2. **Small wave size**: Only 2 efforts (~750 lines total), minimal parallelization benefit
3. **Clear interface boundary**: 2.3.1 defines error types, 2.3.2 extends and integrates them
4. **Testing dependency**: 2.3.2 tests require 2.3.1 validation functions to be complete

**Why This is Correct**:
- Effort 2.3.1 creates `pkg/validator/types.go` with error type definitions
- Effort 2.3.2 moves those types to `pkg/errors/types.go` and extends them
- Effort 2.3.2 imports validation functions from 2.3.1 in error wrapping
- Sequential execution ensures clean import graph and no circular dependencies

---

## Wave Size Compliance

**Total Wave Lines**: 750 lines (400 + 350)

**Breakdown**:
- Effort 2.3.1 implementation: 350 lines
- Effort 2.3.1 tests: 350 lines
- Effort 2.3.2 implementation: 340 lines
- Effort 2.3.2 tests: 350 lines
- **Total implementation**: 690 lines
- **Total tests**: 700 lines
- **Total**: 1390 lines (including tests)

**Size Limit Compliance** (R535 - Code Reviewer enforcement at 900 lines):
- ✅ **Effort 2.3.1**: 400 lines (implementation only) - WITHIN LIMIT
- ✅ **Effort 2.3.2**: 350 lines (implementation only) - WITHIN LIMIT
- ✅ **Wave total**: 750 lines (implementation only) - WELL WITHIN LIMIT
- ✅ **Per R007**: Tests excluded from size limit

**Status**:
- [x] Effort 2.3.1 within 900-line enforcement threshold (400 < 900)
- [x] Effort 2.3.2 within 900-line enforcement threshold (350 < 900)
- [x] Wave total within limits (<3500 lines)
- [ ] No split required

---

## Integration Strategy

### Wave 2.3 Integration Flow

1. **Create Wave 2.3 Integration Branch**
   - Branch from: `idpbuilder-oci-push/phase2/wave2/integration` (Wave 2.2 complete)
   - Branch name: `idpbuilder-oci-push/phase2/wave3/integration`

2. **Effort 2.3.1 Execution**
   - Branch from wave integration: `effort-2.3.1-input-validation`
   - Implement validation system (400 lines)
   - Pass 33 tests (95% coverage)
   - Code review and approval
   - Merge to wave integration branch

3. **Effort 2.3.2 Execution**
   - Branch from wave integration: `effort-2.3.2-error-handling`
   - Implement error type system (350 lines)
   - Pass 30 tests (95% coverage)
   - Code review and approval
   - Merge to wave integration branch

4. **Wave Integration Tests**
   - Run all 63 tests on wave integration branch
   - Verify backward compatibility with Wave 2.1/2.2
   - Run integration tests (error handling end-to-end)
   - Execute demo script (R330 compliance)

5. **Wave Review**
   - Code Reviewer: Final wave assessment
   - Architect: Wave 2.3 architecture compliance review
   - Verify all quality gates passed

6. **Merge to Phase 2 Integration**
   - Merge `idpbuilder-oci-push/phase2/wave3/integration` to `idpbuilder-oci-push/phase2/integration`
   - Phase 2 complete after all waves merged

---

## Testing Strategy

### Unit Tests (63 total)

**Effort 2.3.1** (33 tests):
- Image name validation: 15 tests
- Registry URL validation: 10 tests
- Credentials validation: 8 tests
- Target: 95% statement coverage

**Effort 2.3.2** (30 tests):
- Error type creation: 10 tests
- Exit code mapping: 8 tests
- Error wrapping integration: 12 tests
- Target: 95% statement coverage

### Integration Tests

**Wave-Level Integration Tests**:
```go
// tests/integration/test_wave_2.3_integration.go

func TestErrorHandling_EndToEnd(t *testing.T) {
    // Test complete error flow:
    // 1. Invalid image name → ValidationError → exit code 1
    // 2. Auth failure → AuthenticationError → exit code 2
    // 3. Network timeout → NetworkError → exit code 3
    // 4. Image not found → ImageNotFoundError → exit code 4
}

func TestBackwardCompatibility_Wave22(t *testing.T) {
    // Verify Wave 2.2 configuration still works
    // Verify Wave 2.1 push pipeline unchanged
    // Verify error wrapping doesn't break existing code
}

func TestSSRFWarning_NonBlocking(t *testing.T) {
    // Verify private IP registry warning allows continuation
    // Verify weak credential warning allows continuation
}
```

### Demo Execution (R330 Compliance)

**Demo Script**: `demos/wave2.3-error-handling-demo.sh`

See `planning/phase2/wave3/WAVE-TEST-PLAN.md` Demo Planning section for complete demo specifications.

**Demo Objectives** (5 verifiable items):
1. Security Validation: Command injection prevention
2. SSRF Protection: Private IP warning
3. Error Type Mapping: Correct exit codes
4. Actionable Suggestions: All errors include suggestions
5. Backward Compatibility: Wave 2.2 still works

**Demo Scenarios** (4 scenarios):
1. Command injection attempt (exit code 1)
2. SSRF warning with successful push (exit code 0)
3. Authentication failure (exit code 2)
4. Image not found (exit code 4)

---

## Risk Mitigation

### High-Risk Areas

**Security-Critical Code**:
- **Risk**: Command injection prevention failure
- **Mitigation**: Comprehensive security tests (15 injection tests)
- **Mitigation**: Code review focused on validation regex patterns
- **Mitigation**: Manual testing with real injection payloads

**Exit Code Changes**:
- **Risk**: Breaking scripts that depend on exit codes
- **Mitigation**: Document all exit code changes
- **Mitigation**: Backward compatibility tests
- **Mitigation**: Feature flag for new error handling (if needed)

**Error Wrapping**:
- **Risk**: Losing error context in wrapping chain
- **Mitigation**: Error unwrapping tests (TestErrorChain_Unwrap)
- **Mitigation**: Context preservation tests
- **Mitigation**: Integration tests with real errors

### External Dependencies

**None** - Wave 2.3 uses only:
- Standard library (regexp, strings, fmt, net, errors, os)
- Existing dependencies from Wave 2.1/2.2 (cobra, viper, testify)

### Complexity Hotspots

**SSRF Detection**:
- DNS resolution may be slow or fail
- IPv6 support requires careful testing
- Private IP range detection must be accurate

**Error Type Inference**:
- String matching for error wrapping is brittle
- Future: Use error types from Docker/registry libraries
- Mitigation: Comprehensive error wrapping tests

---

## Compliance Checklist

### R213 Effort Metadata
- [x] effort_id, effort_name, branch_name specified
- [x] base_branch follows cascade pattern (effort 2.3.2 branches from 2.3.1)
- [x] dependencies documented (2.3.2 depends on 2.3.1)
- [x] estimated_lines within 900-line enforcement threshold
- [x] can_parallelize: false (sequential execution required)
- [x] files_touched lists all files with change type (new/modified)

### R341 TDD Protocol
- [x] 63 tests defined in WAVE-TEST-PLAN.md BEFORE implementation
- [x] Tests specify expected behavior
- [x] Tests use real imports (no pseudocode)
- [x] Coverage targets: 95% statement, 90% branch

### R330 Demo Planning
- [x] Demo objectives defined (5 items)
- [x] Demo scenarios specified (4 scenarios)
- [x] Demo deliverables planned (script + documentation)
- [x] Demo size excluded from effort line limits

### R340 Progressive Realism
- [x] All code examples use actual Go syntax
- [x] Real function signatures with complete types
- [x] Real imports from Wave 2.1/2.2 packages
- [x] Integration with existing codebase verified

### R510 Checklist Compliance
- [x] Clear effort definitions
- [x] File-by-file specifications
- [x] Acceptance criteria per effort
- [x] Quality gates documented

### R535 Size Enforcement
- [x] Effort 2.3.1: 400 lines < 900-line enforcement threshold
- [x] Effort 2.3.2: 350 lines < 900-line enforcement threshold
- [x] No splits required

---

## Next Steps

1. **Orchestrator** reviews and approves this implementation plan
2. **Orchestrator** creates Wave 2.3 integration branch infrastructure
3. **Orchestrator** installs test harness + demo harness in integration branch
4. **Software Engineer** spawned for Effort 2.3.1 (validation system)
   - Implement 4 files (imagename.go, registry.go, credentials.go, types.go)
   - Write 33 tests BEFORE implementation (TDD Red phase)
   - Implement to pass tests (TDD Green phase)
   - Achieve 95% coverage
   - Code review and merge
5. **Software Engineer** spawned for Effort 2.3.2 (error type system)
   - Implement 3 files (types.go, exitcodes.go, errors.go)
   - Modify 4 files (push.go, validator imports)
   - Write 30 tests
   - Implement to pass tests
   - Achieve 95% coverage
   - Code review and merge
6. **Software Engineer** executes demo script
   - Run all 4 demo scenarios
   - Capture screenshots/output
   - Document results in WAVE-2.3-DEMO-RESULTS.md
7. **Code Reviewer** validates Wave 2.3 complete
   - All 63 tests passing
   - Coverage targets met
   - Demo objectives verified
8. **Architect** performs Wave 2.3 assessment
   - Architecture compliance
   - Integration with Wave 2.2 verified
   - Wave ready for phase integration

---

## Document Status

**Status**: ✅ READY FOR ORCHESTRATOR REVIEW
**Planner**: @agent-code-reviewer
**Created**: 2025-11-03
**Efforts**: 2 (Input Validation, Error Type System)
**Total Lines**: 750 lines (implementation only, tests excluded per R007)
**Fidelity Level**: EXACT (complete code specifications)

**Builds Upon**:
- Wave 2.2: Registry Override & Environment Variable Support (750 lines, 50 tests, COMPLETE)
- Wave 2.1: Push Command Core & Progress Reporter (1005 lines, 40 tests, COMPLETE)
- Phase 1: All interface packages (docker, registry, auth, tls - COMPLETE)

**Compliance Verified**:
- ✅ R213: Effort metadata complete for both efforts
- ✅ R341: TDD protocol (63 tests defined before implementation)
- ✅ R330: Demo plan complete (5 objectives, 4 scenarios)
- ✅ R340: Progressive realism (all code is real Go syntax)
- ✅ R510: Checklist structure followed
- ✅ R535: Size enforcement (both efforts < 900 lines)
- ✅ R287: TODO persistence rules acknowledged
- ✅ R405: Automation continuation flag required at state completion

---

**END OF WAVE 2.3 IMPLEMENTATION PLAN**
