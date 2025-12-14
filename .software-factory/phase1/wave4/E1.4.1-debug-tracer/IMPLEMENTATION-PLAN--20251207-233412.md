---
created: "2025-12-07"
modified: "2025-12-07"
agent: "code-reviewer"
state: "EFFORT_PLAN_CREATION"
version: "1.0"
---

# EFFORT PLAN - E1.4.1 Debug Tracer

## Effort Identification

**Effort ID**: E1.4.1
**Effort Name**: debug-tracer
**Phase**: 1 - Core OCI Push Implementation
**Wave**: 4 - Debug Capabilities
**Created By**: Code Reviewer Agent
**Date Created**: 2025-12-07
**Assigned To**: SW Engineer
**Estimated Lines**: 150 (must be ≤800)

---

## 🎯 EFFORT THEME DECLARATION (R372)

**PRIMARY THEME**: Add structured debug logging infrastructure with HTTP request/response tracing for the idpbuilder push command

**THEME BOUNDARY**: Debug logging helpers, HTTP transport wrapper, credential resolution logging, request ID correlation

**THEME SPIRIT**: Enable developers to troubleshoot push operations without exposing credentials, using structured logging with correlation IDs

---

## 📋 SCOPE DEFINITION (R371)

### FILES IN SCOPE (ALLOWED)

**EXACT files to create/modify:**
```
pkg/cmd/push/tracer.go          # Debug logging helpers and DebugTransport
pkg/cmd/push/tracer_test.go     # Unit tests for tracer
tests/property/wave4_prop1_test.go   # Property test W1.4.1 (No Credential Logging)
tests/property/wave4_prop2_test.go   # Property test W1.4.2 (Request/Response Correlation)
tests/e2e/push/debug_test.go    # E2E integration tests
```

**Maximum files**: 5 files (well under 20)
**Estimated lines**: 150 lines total (well under 800)

### 🚫 OUT OF SCOPE (FORBIDDEN)

**DO NOT MODIFY:**
- Build system files (Makefile, go.mod - unless adding rapid dependency)
- Infrastructure files (.devcontainer/, docker/)
- CI/CD files (.github/, .gitlab-ci.yml)
- Documentation (unless adding --log-level flag docs)
- Wave 1-3 core implementations (only extend with logger injection)
- Test framework/infrastructure
- Configuration files

---

## 📊 SCOPE TRACEABILITY MATRIX

| File/Package | Requirement | Justification |
|--------------|-------------|---------------|
| pkg/cmd/push/tracer.go | REQ-005, REQ-020, REQ-025 | Core debug logging implementation |
| pkg/cmd/push/tracer_test.go | REQ-020 | Unit tests for credential redaction |
| tests/property/wave4_prop1_test.go | REQ-020 | Property test: No credential logging |
| tests/property/wave4_prop2_test.go | REQ-025 | Property test: Request/response correlation |
| tests/e2e/push/debug_test.go | REQ-005, REQ-006 | E2E validation of debug and info modes |

---

## VALIDATES REQUIREMENTS

> **Traceability:** This effort contributes to validating the following requirements and properties.

### Requirements Validated

| Requirement ID | Requirement Description | How This Effort Validates |
|----------------|------------------------|---------------------------|
| REQ-005 | When --log-level debug is set, the system shall log all HTTP requests and responses to stderr | Implements DebugTransport with HTTP request/response dumping |
| REQ-006 | When --log-level info is set (default), the system shall log only high-level operational steps | Implements configurable logging levels (info vs debug) |
| REQ-020 | The system shall never log credential values at any logging level | Implements Authorization header redaction and credential presence flags only |
| REQ-025 | When debug mode is enabled, the system shall log full HTTP request/response dumps with correlation IDs | Implements request ID generation and correlation between requests and responses |

### EARS Criteria Validated

| EARS ID | EARS Criterion | How This Effort Validates |
|---------|---------------|---------------------------|
| REQ-005 | WHEN --log-level debug THEN system SHALL log HTTP requests/responses | DebugTransport.RoundTrip() logs requests and responses when debug level |
| REQ-006 | WHEN --log-level info THEN system SHALL log operational steps only | NewDebugLogger() configures level, DebugTransport only active in debug mode |
| REQ-020 | WHEN logging credentials THEN system SHALL NOT log values | LogCredentialResolution() logs has_* flags only, DebugTransport redacts Authorization header |
| REQ-025 | WHEN debug mode THEN system SHALL log with correlation IDs | generateRequestID() creates unique IDs, logRequest/logResponse use same ID |

---

## VALIDATES PROPERTIES

> **Reference:** `docs/CORRECTNESS-PROPERTIES-GUIDE.md`

### Properties Implemented/Validated by This Effort

| Property ID | Property Statement | Test File | Status |
|-------------|-------------------|-----------|--------|
| W1.4.1 | For any log output at any level, credential values SHALL NOT appear | `tests/property/wave4_prop1_test.go` | [ ] Pending |
| W1.4.2 | For any HTTP request logged in debug mode, a corresponding response SHALL be logged with the same request_id | `tests/property/wave4_prop2_test.go` | [ ] Pending |

### Property Test Implementation

**Property W1.4.1: No Credential Logging**

This effort implements/validates the following property:

> For any log output at any level (debug, info, warn), credential values SHALL NOT appear. Only presence/absence flags may be logged (has_username, has_password, has_token).

**Test File:** `tests/property/wave4_prop1_test.go`

**Generator Notes:**
- Generate random credentials (username, password, token) as strings
- Use rapid.String() for arbitrary string generation
- Test with empty strings, short strings, long strings
- Boundary conditions: empty credentials, Unicode characters, special characters

**Test Skeleton (TDD - WRITE THIS FIRST):**
```go
package property_test

import (
    "bytes"
    "log/slog"
    "testing"

    "github.com/stretchr/testify/assert"
    "pgregory.net/rapid"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
)

func TestProperty_W1_4_1_NoCredentialLogging(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // Generator: Random credentials
        username := rapid.String().Draw(t, "username")
        password := rapid.String().Draw(t, "password")
        token := rapid.String().Draw(t, "token")

        // Capture log output
        var buf bytes.Buffer
        testLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        }))

        // Action: Log credential resolution
        push.LogCredentialResolution(
            testLogger,
            "flags",
            len(username) > 0,
            len(password) > 0,
            len(token) > 0,
        )

        output := buf.String()

        // Property: No credential values in output
        if len(username) > 0 {
            assert.NotContains(t, output, username,
                "username value should NOT be logged")
        }
        if len(password) > 0 {
            assert.NotContains(t, output, password,
                "password value should NOT be logged")
        }
        if len(token) > 0 {
            assert.NotContains(t, output, token,
                "token value should NOT be logged")
        }

        // Should only log presence/absence flags
        assert.Contains(t, output, "has_username=")
        assert.Contains(t, output, "has_password=")
        assert.Contains(t, output, "has_token=")
    })
}
```

**Property W1.4.2: Request/Response Correlation**

This effort implements/validates the following property:

> For any HTTP request logged in debug mode, a corresponding response SHALL be logged with the same request_id. This enables tracing request/response pairs.

**Test File:** `tests/property/wave4_prop2_test.go`

**Generator Notes:**
- Generate random HTTP request parameters (method, status code)
- Use rapid.SampledFrom() for valid HTTP methods
- Use rapid.IntRange() for status codes (200-599)
- Test various method/status combinations

**Test Skeleton (TDD - WRITE THIS FIRST):**
```go
package property_test

import (
    "bytes"
    "fmt"
    "io"
    "log/slog"
    "net/http"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
    "pgregory.net/rapid"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/push"
)

func TestProperty_W1_4_2_RequestResponseCorrelation(t *testing.T) {
    rapid.Check(t, func(t *rapid.T) {
        // Generator: Random HTTP request
        method := rapid.SampledFrom([]string{"GET", "POST", "PUT", "HEAD"}).Draw(t, "method")
        statusCode := rapid.IntRange(200, 599).Draw(t, "status_code")

        var buf bytes.Buffer
        testLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
            Level: slog.LevelDebug,
        }))

        // Mock round tripper
        mockRT := &mockRoundTripper{
            response: &http.Response{
                StatusCode: statusCode,
                Status:     fmt.Sprintf("%d Status", statusCode),
                Body:       io.NopCloser(bytes.NewBufferString("")),
            },
        }

        transport := &push.DebugTransport{
            Base:   mockRT,
            Logger: testLogger,
        }

        req, _ := http.NewRequest(method, "https://registry.example.com/v2/", nil)
        _, _ = transport.RoundTrip(req)

        output := buf.String()
        lines := strings.Split(output, "\n")

        // Extract request_id from request log
        var requestID string
        for _, line := range lines {
            if strings.Contains(line, "HTTP request") && strings.Contains(line, "request_id=") {
                parts := strings.Split(line, "request_id=")
                if len(parts) > 1 {
                    requestID = strings.Fields(parts[1])[0]
                    break
                }
            }
        }

        // Property: Request ID must be present
        assert.NotEmpty(t, requestID, "request should have request_id")

        // Property: Response log has same request_id
        responseFound := false
        for _, line := range lines {
            if strings.Contains(line, "HTTP response") && strings.Contains(line, requestID) {
                responseFound = true
                break
            }
        }

        assert.True(t, responseFound,
            "response should be logged with same request_id as request")
    })
}

type mockRoundTripper struct {
    response *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
    return m.response, nil
}
```

---

## PROPERTY TEST FILES

> **Convention:** Property test files follow `*_property_test.[ext]` naming.

| Property | Test File Location | Effort Responsibility |
|----------|-------------------|----------------------|
| W1.4.1 | `tests/property/wave4_prop1_test.go` | Create (TDD - write FIRST) |
| W1.4.2 | `tests/property/wave4_prop2_test.go` | Create (TDD - write FIRST) |

### Property Test Checklist

- [ ] Property test W1.4.1 created FIRST (TDD red phase)
- [ ] Property test W1.4.2 created FIRST (TDD red phase)
- [ ] Both tests FAILING initially (confirms they check correct behavior)
- [ ] Generator implemented with correct constraints (genCredentials, genHTTPRequest)
- [ ] Tests run with 100 examples in CI configuration
- [ ] Seed logging enabled for reproducibility
- [ ] After implementation, all property tests PASS (TDD green phase)

---

## 🎬 Demo Requirements (R330 MANDATORY)

### Demo Objectives (3-5 specific, verifiable objectives)
- [ ] Demonstrate debug mode shows HTTP request/response dumps with --log-level debug
- [ ] Show proper credential redaction (Authorization: [REDACTED]) in debug output
- [ ] Verify info mode shows only high-level operations, no HTTP dumps
- [ ] Prove request ID correlation between request and response logs
- [ ] Display credential source detection (flags vs env vs anonymous)

**Success Criteria**: All objectives checked = demo passes

### Demo Scenarios (IMPLEMENT EXACTLY THESE - 3 scenarios)

#### Scenario 1: Debug Mode HTTP Tracing (Tests REQ-005, REQ-025)
- **Requirement**: REQ-005 - Debug logging enabled, REQ-025 - HTTP correlation
- **Setup**: Test image built, Gitea registry accessible, credentials configured
- **Input**: `idpbuilder push myimage:latest --log-level debug --registry gitea.cnoe.localtest.me --username admin --password <password>`
- **Action**: Execute push command with debug logging
- **Expected Output**:
  ```
  level=DEBUG msg="credential resolution" source=flags has_username=true has_password=true has_token=false
  level=DEBUG msg="HTTP request" request_id=req-1733614892543210 method=PUT url=https://gitea.cnoe.localtest.me/v2/myimage/manifests/latest dump="PUT /v2/myimage/manifests/latest HTTP/1.1\r\nHost: gitea.cnoe.localtest.me\r\nAuthorization: [REDACTED]\r\n..."
  level=DEBUG msg="HTTP response" request_id=req-1733614892543210 status_code=201 duration=45ms dump="HTTP/1.1 201 Created\r\n..."
  ```
- **Verification**:
  - Stderr contains "HTTP request" and "HTTP response" logs
  - Same request_id appears in both logs
  - Authorization header shows [REDACTED], NOT actual password
  - credential resolution shows has_* flags, NOT actual values
- **R335 Compliance**: Directly tests REQ-005, REQ-025
- **Script Lines**: ~30 lines

#### Scenario 2: Info Mode Selective Logging (Tests REQ-006)
- **Requirement**: REQ-006 - Info level shows only operational steps
- **Setup**: Same test image and credentials
- **Input**: `idpbuilder push myimage:latest --log-level info --registry gitea.cnoe.localtest.me --username admin --password <password>`
- **Action**: Execute push command with info logging
- **Expected Output**:
  ```
  level=INFO msg="Pushing image myimage:latest to gitea.cnoe.localtest.me"
  level=INFO msg="Layer pushed" digest=sha256:abc123... size=1234567
  level=INFO msg="Push complete" ref=gitea.cnoe.localtest.me/myimage:latest
  ```
- **Verification**:
  - Stderr does NOT contain "HTTP request" or "HTTP response" (debug only)
  - Stderr does NOT contain "credential resolution" (debug only)
  - Stderr contains "Pushing image" and "Push complete" (info level)
  - Push completes successfully without verbose output
- **R335 Compliance**: Directly tests REQ-006
- **Script Lines**: ~20 lines

#### Scenario 3: Credential Source Detection (Tests REQ-020)
- **Requirement**: REQ-020 - No credential exposure in logs
- **Setup**: Test image, configure credentials via environment variables instead of flags
- **Input**:
  ```bash
  export REGISTRY_USERNAME=admin
  export REGISTRY_PASSWORD=<password>
  idpbuilder push myimage:latest --log-level debug --registry gitea.cnoe.localtest.me
  ```
- **Action**: Execute push with credentials from environment
- **Expected Output**:
  ```
  level=DEBUG msg="credential resolution" source=env has_username=true has_password=true has_token=false
  ```
- **Verification**:
  - Debug output shows source=env (not source=flags)
  - Debug output shows has_username=true, has_password=true
  - Actual password value does NOT appear anywhere in output
  - Environment variable names do NOT appear in output (only "source=env")
- **R335 Compliance**: Directly tests REQ-020 (no credential exposure)
- **Script Lines**: ~25 lines

**TOTAL SCENARIO LINES**: ~75 lines

### Demo Size Planning

#### Demo Artifacts (Excluded from line count per R007)
```
demo-debug-tracer.sh:     75 lines  # Executable script with 3 scenarios
DEMO.md:                  40 lines  # Documentation
test-data/test-image:     10 lines  # Sample Dockerfile
────────────────────────────────
TOTAL DEMO FILES:        125 lines (NOT counted toward 800)
```

#### Effort Size Summary
```
Implementation:     150 lines  # ← ONLY this counts toward 800
────────────────────────────────
Tests:              200 lines  # Excluded per R007 (property + unit + e2e)
Demos:              125 lines  # Excluded per R007
────────────────────────────────
Implementation:    150/800 ✅ (well within limit)
```

**NOTE**: While demos don't count toward the line limit, they MUST still be planned and implemented as specified!

### Demo Deliverables

Required Files:
- [ ] `demo-debug-tracer.sh` - Main demo script (executable, 3 scenarios)
- [ ] `DEMO.md` - Demo documentation per template
- [ ] `test-data/test-image/Dockerfile` - Test image for push demos
- [ ] `.demo-config` - Demo environment settings (Gitea URL, default credentials)

Integration Hooks:
- [ ] Export DEMO_READY=true when complete
- [ ] Provide integration point for wave demo (wave4-integration will aggregate)
- [ ] Include cleanup function (remove test images after demo)

---

## Assigned Requirements (R335 Scoped - MANDATORY)

### Requirements This Effort Must Demonstrate

| REQ ID | Requirement Description | Demo Scenario | Priority |
|--------|------------------------|---------------|----------|
| REQ-005 | When --log-level debug is set, the system shall log all HTTP requests and responses to stderr | Scenario 1: Debug mode HTTP tracing | CRITICAL |
| REQ-006 | When --log-level info is set (default), the system shall log only high-level operational steps | Scenario 2: Info mode selective logging | CRITICAL |
| REQ-020 | The system shall never log credential values at any logging level | Scenario 3: Credential source detection | CRITICAL |
| REQ-025 | When debug mode is enabled, the system shall log full HTTP request/response dumps with correlation IDs | Scenario 1: Debug mode HTTP tracing | CRITICAL |

### Coverage Target
- **In-Scope Requirements**: 4 (REQ-005, REQ-006, REQ-020, REQ-025)
- **Coverage Target**: 100% (4/4)
- **Minimum Threshold**: 100% for effort-level demos

### NOT In Scope (Other Efforts)
These requirements are assigned to OTHER waves:
- REQ-001: Assigned to Wave 1 (credential resolution)
- REQ-002: Assigned to Wave 2 (registry selection)
- REQ-003: Assigned to Wave 2 (default registry)
- REQ-010: Assigned to Wave 3 (retry logic)

**DO NOT** create demos for out-of-scope requirements.
**DO** ensure 100% coverage of in-scope requirements.

---

## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found
| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| http.RoundTripper | Go standard library | RoundTrip(req *http.Request) (*http.Response, error) | YES - DebugTransport implements this |
| slog.Logger | Go log/slog package | Debug(msg string, args ...any) | YES - use for structured logging |

### Existing Implementations to Reuse
| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| slog.NewTextHandler | log/slog | Create structured logger | Call in NewDebugLogger() |
| http.DefaultTransport | net/http | Base HTTP transport | Wrap with DebugTransport when debug mode |
| httputil.DumpRequest | net/http/httputil | Dump HTTP request | Call in logRequest() method |
| httputil.DumpResponse | net/http/httputil | Dump HTTP response | Call in logResponse() method |

### APIs Already Defined
| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| slog.Logger.Debug | Debug | Debug(msg string, args ...any) | Use for debug-level logging |
| slog.Logger.Info | Info | Info(msg string, args ...any) | Use for info-level logging |
| http.RoundTripper | RoundTrip | RoundTrip(req *http.Request) (*http.Response, error) | DebugTransport implements this interface |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create custom logging framework (use slog from stdlib)
- DO NOT reimplement HTTP transport from scratch (wrap http.DefaultTransport)
- DO NOT create alternative request dumping (use httputil.DumpRequest/DumpResponse)

### REQUIRED INTEGRATIONS (R373)
- MUST implement http.RoundTripper interface with EXACT signature
- MUST reuse slog.Logger from log/slog package
- MUST import and use httputil.DumpRequest/DumpResponse for HTTP dumps

---

## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:
- Function: NewDebugLogger(level slog.Level, phase string) *slog.Logger (~20 lines)
- Function: LogCredentialResolution(logger *slog.Logger, source string, hasUsername, hasPassword, hasToken bool) (~10 lines)
- Type: DebugTransport struct with Base http.RoundTripper, Logger *slog.Logger (~5 lines)
- Method: (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) (~30 lines)
- Method: (t *DebugTransport) logRequest(req *http.Request, requestID string) (~25 lines)
- Method: (t *DebugTransport) logResponse(resp *http.Response, requestID string, duration time.Duration) (~20 lines)
- Function: generateRequestID() string (~5 lines)

TOTAL Implementation: ~115 lines (core tracer.go)
TOTAL Tests: ~200 lines (property + unit + e2e, NOT counted toward 800)
TOTAL Effort: ~150 lines counted (implementation only per R007)

### DO NOT IMPLEMENT:
- ❌ Logging configuration files (use slog defaults)
- ❌ Log rotation or persistence (logs to stderr only)
- ❌ Log aggregation or shipping (out of scope)
- ❌ Metrics collection (logging only, not metrics)
- ❌ Advanced filtering or sampling (basic debug/info levels only)
- ❌ Custom log formatters (use slog.NewTextHandler)
- ❌ Async/buffered logging (synchronous only for simplicity)
- ❌ Log levels beyond debug/info/warn (3 levels sufficient)

---

## 🚨🚨🚨 R355 PRODUCTION READINESS - ZERO TOLERANCE 🚨🚨🚨

This implementation MUST be production-ready from the first commit:
- ❌ NO STUBS or placeholder implementations
- ❌ NO MOCKS except in test directories
- ❌ NO hardcoded credentials or secrets
- ❌ NO static configuration values
- ❌ NO TODO/FIXME markers in code
- ❌ NO returning nil or empty for "later implementation"
- ❌ NO panic("not implemented") patterns
- ❌ NO fake or dummy data

**VIOLATION = -100% AUTOMATIC FAILURE**

### Configuration Requirements (R355 Mandatory)

**WRONG - Will fail review:**
```go
// ❌ VIOLATION - Hardcoded level
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// ❌ VIOLATION - Stub implementation
func LogCredentialResolution(logger *slog.Logger, source string) {
    // TODO: implement later
}

// ❌ VIOLATION - Logging credential values
logger.Debug("credentials", slog.String("password", password))
```

**CORRECT - Production ready:**
```go
// ✅ Configurable level from parameter
func NewDebugLogger(level slog.Level, phase string) *slog.Logger {
    handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
        Level: level,
    })
    return slog.New(handler).With(slog.String("phase", phase))
}

// ✅ Full implementation with all parameters
func LogCredentialResolution(logger *slog.Logger, source string, hasUsername, hasPassword, hasToken bool) {
    logger.Debug("credential resolution",
        slog.String("source", source),
        slog.Bool("has_username", hasUsername),
        slog.Bool("has_password", hasPassword),
        slog.Bool("has_token", hasToken),
    )
}

// ✅ NO credential values, only presence flags
// ✅ Authorization header redacted
reqCopy := req.Clone(req.Context())
if reqCopy.Header.Get("Authorization") != "" {
    reqCopy.Header.Set("Authorization", "[REDACTED]")
}
```

---

## Size Limit Clarification (R359):
- The 800-line limit applies to NEW CODE YOU ADD
- Repository will grow by ~150 lines (EXPECTED)
- NEVER delete existing code to meet size limits
- Example: If repo has 5,000 lines and you add 150, total will be 5,150

## Implementation Size Estimate:
- NEW code to be added: ~150 lines
- Existing codebase: ~5,000 lines (from Wave 1-3)
- Expected total after: ~5,150 lines

---

## ✅ ACCEPTANCE CRITERIA

1. [ ] All files in scope created (tracer.go, tests)
2. [ ] NO files outside scope touched
3. [ ] Theme coherence maintained: single debug logging theme (>95%)
4. [ ] Property tests pass (W1.4.1, W1.4.2) with 100 examples
5. [ ] Unit tests pass (tracer_test.go)
6. [ ] E2E tests pass (debug mode, info mode, credential detection)
7. [ ] Line count <800 (actual: ~150 lines)
8. [ ] No stubs/mocks/TODOs in production code (R355)
9. [ ] Demo requirements satisfied (R330) - 3 scenarios executable
10. [ ] All 4 assigned requirements validated (REQ-005, REQ-006, REQ-020, REQ-025)

---

## 🔍 VALIDATION CHECKLIST

**BEFORE STARTING:**
- [ ] Theme is single and focused: Debug logging only
- [ ] File list is explicit and complete: 5 files total
- [ ] OUT OF SCOPE section reviewed: No build system, CI/CD changes
- [ ] <20 files total: YES (5 files)
- [ ] No mixed concerns: Only debug logging, no retry logic, no registry selection

**BEFORE COMMITTING:**
- [ ] Run: `tools/line-counter.sh` (verify <800 lines)
- [ ] Run: Property tests with 100 examples (verify both W1.4.1, W1.4.2 pass)
- [ ] All changes trace to requirements (REQ-005, REQ-006, REQ-020, REQ-025)
- [ ] No scope creep occurred: Only debug logging added
- [ ] Theme purity verified: 100% debug logging infrastructure

---

## 📝 IMPLEMENTATION NOTES

### TDD Workflow (R341 MANDATORY)

**CRITICAL: Write tests BEFORE implementation!**

1. **RED Phase - Write Property Tests First:**
   - Create `tests/property/wave4_prop1_test.go` (W1.4.1)
   - Create `tests/property/wave4_prop2_test.go` (W1.4.2)
   - Run tests: EXPECT FAILURES (no implementation yet)
   - This confirms tests check correct behavior

2. **GREEN Phase - Implement to Make Tests Pass:**
   - Create `pkg/cmd/push/tracer.go` with all functions
   - Implement LogCredentialResolution (makes W1.4.1 pass)
   - Implement DebugTransport (makes W1.4.2 pass)
   - Run property tests: EXPECT PASSES (100 examples)

3. **REFACTOR Phase - Add Unit and E2E Tests:**
   - Create `pkg/cmd/push/tracer_test.go` (unit tests)
   - Create `tests/e2e/push/debug_test.go` (E2E tests)
   - Verify all tests still pass
   - No regressions from refactoring

### Libraries/Dependencies
- `log/slog` (Go standard library) - structured logging
- `net/http/httputil` (Go standard library) - HTTP request/response dumping
- `pgregory.net/rapid` (add to go.mod) - property-based testing framework
- `github.com/stretchr/testify/assert` (add to go.mod) - test assertions

### Patterns to Follow
- Use slog structured logging (not fmt.Printf)
- Use correlation IDs for request/response tracing
- Redact sensitive headers (Authorization, X-Auth-Token, etc.)
- Log to stderr (not stdout - stdout is for command output)
- Use descriptive field names in slog (source, has_username, request_id, etc.)

### Integration Points
- Inject logger into credential resolver (modify ResolveWithLogger signature)
- Wrap HTTP transport when --log-level debug (in push.go RunE function)
- Add --log-level flag to push command (cobra flag)

---

## ⚠️ CRITICAL REMINDERS

- **R371**: The file list above is LAW - no additions allowed (5 files only)
- **R372**: Maintain single theme - debug logging infrastructure only
- **R355**: Production-ready code only - no stubs, no TODOs, no hardcoded values
- **R359**: Never delete existing code for size (only ADD ~150 lines)
- **R362**: Follow approved architecture exactly (from Wave 4 Architecture Plan)
- **R341**: TDD MANDATORY - write property tests BEFORE implementation
- **R765**: PBT MANDATORY - W1.4.1 and W1.4.2 must pass with 100 examples

---

**Generated**: 2025-12-07
**Effort State**: PLANNED
**Review Status**: PENDING
**R604 Compliance**: Plan will be committed with git hash before completion
