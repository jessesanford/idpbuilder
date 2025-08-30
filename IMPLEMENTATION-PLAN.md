# Error Messaging Enhancement Implementation Plan

## CRITICAL EFFORT METADATA (FROM WAVE PLAN)
**Branch**: `phase3/wave1/error-messaging`
**Can Parallelize**: Yes
**Parallel With**: [cert-integration-manager, security-features]
**Size Estimate**: 500 lines
**Dependencies**: None (defines error types)

## Overview
- **Effort**: Error Messaging Enhancement - Clear, actionable certificate error messages
- **Phase**: 3, Wave: 1
- **Estimated Size**: 500 lines total
- **Implementation Time**: 4-6 hours

## Purpose
Create a comprehensive error messaging system specifically for certificate-related issues in the IDPBuilder OCI MVP. The system will provide clear, actionable error messages with resolution guidance to help users quickly identify and fix certificate problems.

## File Structure
```
pkg/
└── errors/
    ├── certificate_errors.go    # Error type definitions (100 lines)
    ├── messages.go              # Message template system (100 lines)
    ├── guidance.go              # Resolution guidance (100 lines)
    ├── context.go               # Error context enrichment (100 lines)
    ├── formatter.go             # User-friendly formatter (100 lines)
    └── tests/
        ├── certificate_errors_test.go  # Unit tests
        ├── messages_test.go            # Template tests
        ├── guidance_test.go            # Guidance tests
        ├── context_test.go             # Context tests
        └── formatter_test.go           # Formatter tests
```

## Implementation Steps

### Step 1: Define Certificate Error Types (100 lines)
**File**: `pkg/errors/certificate_errors.go`

#### Error Taxonomy
```go
// Core error types for certificate issues
const (
    ErrCertNotFound      ErrorType = "CERT_NOT_FOUND"       // Certificate file missing
    ErrCertInvalid       ErrorType = "CERT_INVALID"         // Malformed certificate
    ErrCertExpired       ErrorType = "CERT_EXPIRED"         // Past expiration date
    ErrCertUntrusted     ErrorType = "CERT_UNTRUSTED"       // Not in trust store
    ErrCertMismatch      ErrorType = "CERT_MISMATCH"        // Domain/registry mismatch
    ErrCertPermission    ErrorType = "CERT_PERMISSION"      // Access denied
    ErrCertChainBroken   ErrorType = "CERT_CHAIN_BROKEN"    // Incomplete chain
    ErrCertFormat        ErrorType = "CERT_FORMAT"          // Unsupported format
)
```

#### Error Structure
- Base `CertificateError` struct with type, message, details
- Error wrapping support for preserving original errors
- Method chaining for adding context (`WithDetail()`, `WithResolution()`)
- Severity levels (INFO, WARNING, ERROR, CRITICAL)

#### Implementation Tasks
1. Define ErrorType enumeration
2. Create CertificateError struct with fields
3. Implement Error() interface
4. Add builder methods for context
5. Create constructor functions for each error type

### Step 2: Create Error Message Templates (100 lines)
**File**: `pkg/errors/messages.go`

#### Template System
- Go text/template based message generation
- Placeholder support for dynamic values
- Consistent message format across all errors
- Localization preparation (future enhancement)

#### Message Templates
```go
var errorTemplates = map[ErrorType]string{
    ErrCertNotFound:    "Certificate not found at expected location: {{.Path}}",
    ErrCertInvalid:     "Certificate validation failed: {{.Reason}}",
    ErrCertExpired:     "Certificate expired on {{.ExpiryDate}} ({{.DaysAgo}} days ago)",
    ErrCertUntrusted:   "Certificate issuer '{{.Issuer}}' not in trust store",
    ErrCertMismatch:    "Certificate CN '{{.CN}}' doesn't match registry '{{.Registry}}'",
    ErrCertPermission:  "Permission denied accessing certificate at {{.Path}}: {{.Error}}",
    ErrCertChainBroken: "Certificate chain incomplete: missing {{.MissingLink}}",
    ErrCertFormat:      "Unsupported certificate format: expected {{.Expected}}, got {{.Actual}}",
}
```

#### Implementation Tasks
1. Define template map for all error types
2. Create MessageBuilder struct
3. Implement template parsing and caching
4. Add BuildMessage() method with data injection
5. Handle template execution errors gracefully

### Step 3: Implement Resolution Guidance (100 lines)
**File**: `pkg/errors/guidance.go`

#### Resolution Framework
- Step-by-step resolution instructions
- Command examples with actual values
- Documentation links for detailed help
- Workarounds for development scenarios

#### Guidance Structure
```go
type ResolutionGuide struct {
    Steps       []string              // Ordered resolution steps
    Examples    []string              // Command examples
    DocLinks    []string              // Documentation URLs
    Workaround  string               // Quick workaround if available
    Severity    string               // How urgent is this?
    AutoFix     func() error         // Optional auto-fix function
}
```

#### Resolution Guides
- **ErrCertNotFound**: Extract certificates, verify paths, check cluster
- **ErrCertExpired**: Regenerate certificates, update rotation config
- **ErrCertUntrusted**: Add to trust store, verify issuer
- **ErrCertMismatch**: Update certificate CN, verify registry URL
- **ErrCertPermission**: Fix file permissions, check user access
- **ErrCertChainBroken**: Get intermediate certificates, rebuild chain
- **ErrCertFormat**: Convert certificate format, check encoding

#### Implementation Tasks
1. Define ResolutionGuide struct
2. Create resolution map for each error type
3. Implement GetResolutionGuide() lookup
4. Add Format() method for display
5. Include auto-fix suggestions where possible

### Step 4: Add Error Context Enrichment (100 lines)
**File**: `pkg/errors/context.go`

#### Context Collection
- System information (OS, architecture, versions)
- Environment variables relevant to certificates
- Component and operation that triggered error
- Timestamp and correlation IDs
- Diagnostic command outputs

#### Context Structure
```go
type ErrorContext struct {
    Timestamp    time.Time
    Component    string                // Which component failed
    Operation    string                // What was being attempted
    Environment  map[string]string     // Relevant env vars
    SystemInfo   SystemInfo            // OS, arch, versions
    Diagnostics  map[string]string     // Command outputs
    CorrelationID string               // For tracking related errors
}

type SystemInfo struct {
    OS           string
    Arch         string
    GoVersion    string
    KindVersion  string
    DockerInfo   string
    KubeVersion  string
}
```

#### Diagnostic Collection
- Automatic capture of relevant system state
- Safe collection (no sensitive data)
- Lightweight to avoid performance impact
- Optional verbose mode for detailed diagnostics

#### Implementation Tasks
1. Define ErrorContext and SystemInfo structs
2. Implement CaptureContext() function
3. Add safe environment variable collection
4. Create diagnostic command runner
5. Implement context serialization for logs

### Step 5: Create User-Friendly Formatter (100 lines)
**File**: `pkg/errors/formatter.go`

#### Formatting Features
- Console output with color coding
- Structured output (JSON/YAML) for automation
- Severity-based formatting (colors, symbols)
- Error aggregation for multiple issues
- Progress indicators for resolution steps

#### Format Options
```go
type FormatOptions struct {
    UseColor     bool              // Enable terminal colors
    Verbose      bool              // Include full context
    Format       OutputFormat      // console, json, yaml
    ShowGuidance bool              // Include resolution steps
    ShowContext  bool              // Include system context
}

type OutputFormat string
const (
    FormatConsole OutputFormat = "console"
    FormatJSON    OutputFormat = "json"
    FormatYAML    OutputFormat = "yaml"
)
```

#### Console Output Example
```
❌ CERTIFICATE ERROR [CERT_EXPIRED]
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Certificate expired on 2024-12-01 (30 days ago)

📍 Details:
  • Path: /tmp/certs/registry.crt
  • CN: registry.local
  • Issuer: IDPBuilder CA
  • Not After: 2024-12-01 15:30:00 UTC

💡 Resolution Steps:
  1. Regenerate certificates in the cluster
  2. Update certificate rotation configuration
  3. Restart affected services

📋 Example Commands:
  $ kubectl delete secret -n gitea gitea-tls-cert
  $ kubectl rollout restart deployment -n gitea gitea

⚠️  Workaround: Use --insecure flag for development (NOT for production)

📚 Documentation: https://github.com/idpbuilder/docs/cert-rotation
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Implementation Tasks
1. Define FormatOptions struct
2. Create ErrorFormatter interface
3. Implement ConsoleFormatter with colors
4. Add JSONFormatter for structured output
5. Create error aggregation logic

## Size Management
- **Estimated Lines**: 500 total
  - certificate_errors.go: 100 lines
  - messages.go: 100 lines
  - guidance.go: 100 lines
  - context.go: 100 lines
  - formatter.go: 100 lines
- **Measurement Tool**: /home/vscode/workspaces/idpbuilder-oci-mvp/tools/line-counter.sh
- **Check Frequency**: After each component completion
- **Split Threshold**: 700 lines (warning), 800 lines (stop)

## Test Requirements
- **Unit Tests**: 90% coverage minimum
- **Test Scenarios**:
  - All error types properly formatted
  - Template rendering with various data
  - Resolution guide accuracy
  - Context capture completeness
  - Formatter output validation
- **Test Files**:
  - pkg/errors/tests/certificate_errors_test.go
  - pkg/errors/tests/messages_test.go
  - pkg/errors/tests/guidance_test.go
  - pkg/errors/tests/context_test.go
  - pkg/errors/tests/formatter_test.go

## Pattern Compliance
- **Error Handling**: Using Go's error interface with rich context
- **Builder Pattern**: For constructing complex error objects
- **Template Pattern**: For message generation
- **Factory Pattern**: For creating specific error types
- **Strategy Pattern**: For different formatting options

## Integration Points
- **FROM All Components**: Any component can create certificate errors
- **TO CLI**: Formatted errors displayed to users
- **TO Logs**: Structured errors for debugging
- **TO Monitoring**: Error metrics and alerts (future)

## Success Criteria
1. ✅ All 8 certificate error types defined and tested
2. ✅ Message templates render correctly with data
3. ✅ Resolution guides provide actionable steps
4. ✅ Context capture includes relevant diagnostics
5. ✅ Console formatter produces clear, readable output
6. ✅ Total implementation stays under 500 lines
7. ✅ Test coverage exceeds 90%
8. ✅ No external dependencies beyond standard library

## Risk Mitigation
- **Risk**: Message templates becoming too complex
  - **Mitigation**: Keep templates simple, move logic to code
- **Risk**: Context collection impacting performance
  - **Mitigation**: Lazy evaluation, optional verbose mode
- **Risk**: Resolution guides becoming outdated
  - **Mitigation**: Version-specific guides, documentation links

## Development Sequence
1. Start with certificate_errors.go (foundation)
2. Implement messages.go (templates)
3. Add guidance.go (resolution)
4. Create context.go (diagnostics)
5. Finish with formatter.go (output)
6. Write comprehensive tests throughout

## Notes for SW Engineer
- Focus on clarity over brevity in error messages
- Always provide actionable next steps
- Include examples with real commands
- Test with actual certificate scenarios
- Consider user experience first
- Keep messages consistent across all errors
- Use color sparingly but effectively
- Ensure errors are grep-friendly