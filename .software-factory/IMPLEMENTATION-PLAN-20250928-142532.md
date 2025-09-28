# Implementation Plan for Progress Reporter Interface

**Created**: 2025-09-28T14:25:32Z
**Agent**: code-reviewer
**Phase**: 1
**Wave**: 2
**Effort**: E1.2.5 - Progress Reporter Interface
**Branch**: phase1/wave2/progress-reporter
**Base Branch**: phase1-wave1-integration

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

VIOLATION = -100% AUTOMATIC FAILURE

## Pre-Planning Research Results (R374 MANDATORY)

### Existing Interfaces Found

| Interface | Location | Signature | Must Implement |
|-----------|----------|-----------|----------------|
| ProgressReporter | /efforts/phase1/wave1/P1W1-E4-cli-contracts/pkg/cmd/interfaces/types.go:60 | ReportProgress(current, total int64, message string)<br>Start(message string)<br>Complete(message string)<br>Error(err error) | YES - EXACT MATCH |

### Existing Implementations to Reuse

| Component | Location | Purpose | How to Use |
|-----------|----------|---------|------------|
| CommandOptions | P1W1-E4/pkg/cmd/interfaces/types.go:8 | Common CLI options | Import for Verbose flag support |
| RegistryConfig | P1W1-E4/pkg/cmd/interfaces/types.go:86 | Registry configuration | Reference for progress context |

### APIs Already Defined

| API | Method | Signature | Notes |
|-----|--------|-----------|-------|
| ProgressReporter.Start | Start | Start(message string) | Begin progress session |
| ProgressReporter.ReportProgress | ReportProgress | ReportProgress(current, total int64, message string) | Update progress |
| ProgressReporter.Complete | Complete | Complete(message string) | Mark completion |
| ProgressReporter.Error | Error | Error(err error) | Report errors |

### FORBIDDEN DUPLICATIONS (R373)
- DO NOT create alternative ProgressReporter interface
- DO NOT reimplement the interface signatures
- DO NOT create competing progress reporting types
- DO NOT define different method names for same functionality

### REQUIRED INTEGRATIONS (R373)
- MUST implement interfaces.ProgressReporter from P1W1-E4 with EXACT signature
- MUST import "github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces"
- MUST provide concrete implementations for console, JSON, and multi-reporter

## EXPLICIT SCOPE (R311 MANDATORY)

### IMPLEMENT EXACTLY:

**Types (2 total, ~30 lines)**:
- Type: `Operation` struct with Name, Total, Current, StartTime fields (~15 lines)
- Type: `Result` struct with Success, Message, Duration fields (~15 lines)

**Implementations (3 total, ~170 lines)**:
- Function: `NewConsoleReporter() interfaces.ProgressReporter` (~60 lines)
- Function: `NewJSONReporter(writer io.Writer) interfaces.ProgressReporter` (~60 lines)
- Function: `NewMultiReporter(reporters ...interfaces.ProgressReporter) interfaces.ProgressReporter` (~50 lines)

**Helper Functions (2 total, ~30 lines)**:
- Function: `formatBytes(bytes int64) string` (~15 lines)
- Function: `formatDuration(d time.Duration) string` (~15 lines)

**Tests (4 test functions, ~120 lines)**:
- Test: `TestConsoleReporter` (~30 lines)
- Test: `TestJSONReporter` (~30 lines)
- Test: `TestMultiReporter` (~30 lines)
- Test: `TestFormatters` (~30 lines)

**TOTAL**: ~350 lines (well under 800 limit)

### DO NOT IMPLEMENT:
- ❌ Progress bar rendering (external library concern)
- ❌ Terminal UI widgets or animations
- ❌ Real-time streaming protocols
- ❌ File-based progress logging
- ❌ Network progress broadcasting
- ❌ Complex formatting beyond basic text
- ❌ Progress persistence/recovery
- ❌ Progress aggregation from multiple sources
- ❌ Historical progress tracking
- ❌ Progress metrics/analytics

## Size Limit Clarification (R359):
- The 800-line limit applies to NEW CODE YOU ADD
- Repository will grow by ~350 lines (EXPECTED)
- NEVER delete existing code to meet size limits
- Current repo size: ~100k lines
- Expected total after: ~100,350 lines

## File Structure

```
pkg/progress/
├── types.go             # Operation and Result types (~30 lines)
├── console_reporter.go  # Terminal output implementation (~60 lines)
├── json_reporter.go     # Structured JSON output (~60 lines)
├── multi.go            # Composite reporter pattern (~50 lines)
├── formatter.go        # Format helper functions (~30 lines)
├── console_reporter_test.go # Console tests (~30 lines)
├── json_reporter_test.go    # JSON tests (~30 lines)
├── multi_test.go           # Multi-reporter tests (~30 lines)
└── formatter_test.go       # Formatter tests (~30 lines)
```

## Implementation Details

### 1. types.go (~30 lines)

```go
package progress

import (
    "time"
)

// Operation represents a long-running operation being tracked
type Operation struct {
    Name      string    // Operation name
    Total     int64     // Total units of work
    Current   int64     // Current progress
    StartTime time.Time // When operation started
}

// Result represents the completion result of an operation
type Result struct {
    Success  bool          // Whether operation succeeded
    Message  string        // Completion or error message
    Duration time.Duration // How long operation took
}
```

### 2. console_reporter.go (~60 lines)

```go
package progress

import (
    "fmt"
    "io"
    "os"
    "time"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces"
)

type consoleReporter struct {
    writer    io.Writer
    operation *Operation
    verbose   bool
}

// NewConsoleReporter creates a progress reporter for terminal output
func NewConsoleReporter() interfaces.ProgressReporter {
    return &consoleReporter{
        writer: os.Stdout,
    }
}

func (c *consoleReporter) Start(message string) {
    c.operation = &Operation{
        Name:      message,
        StartTime: time.Now(),
    }
    fmt.Fprintf(c.writer, "Starting: %s\n", message)
}

func (c *consoleReporter) ReportProgress(current, total int64, message string) {
    if c.operation != nil {
        c.operation.Current = current
        c.operation.Total = total
    }

    if total > 0 {
        percentage := float64(current) / float64(total) * 100
        fmt.Fprintf(c.writer, "[%.1f%%] %s\n", percentage, message)
    } else {
        fmt.Fprintf(c.writer, "[%d] %s\n", current, message)
    }
}

func (c *consoleReporter) Complete(message string) {
    if c.operation != nil {
        duration := time.Since(c.operation.StartTime)
        fmt.Fprintf(c.writer, "Completed: %s (took %s)\n", message, formatDuration(duration))
    } else {
        fmt.Fprintf(c.writer, "Completed: %s\n", message)
    }
}

func (c *consoleReporter) Error(err error) {
    fmt.Fprintf(c.writer, "Error: %v\n", err)
}
```

### 3. json_reporter.go (~60 lines)

```go
package progress

import (
    "encoding/json"
    "io"
    "os"
    "time"

    "github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces"
)

type jsonReporter struct {
    writer    io.Writer
    encoder   *json.Encoder
    operation *Operation
}

// NewJSONReporter creates a progress reporter with JSON output
func NewJSONReporter(writer io.Writer) interfaces.ProgressReporter {
    if writer == nil {
        writer = os.Stdout
    }
    return &jsonReporter{
        writer:  writer,
        encoder: json.NewEncoder(writer),
    }
}

func (j *jsonReporter) Start(message string) {
    j.operation = &Operation{
        Name:      message,
        StartTime: time.Now(),
    }

    j.encoder.Encode(map[string]interface{}{
        "event":     "start",
        "message":   message,
        "timestamp": j.operation.StartTime,
    })
}

func (j *jsonReporter) ReportProgress(current, total int64, message string) {
    if j.operation != nil {
        j.operation.Current = current
        j.operation.Total = total
    }

    j.encoder.Encode(map[string]interface{}{
        "event":   "progress",
        "current": current,
        "total":   total,
        "message": message,
    })
}

func (j *jsonReporter) Complete(message string) {
    var duration time.Duration
    if j.operation != nil {
        duration = time.Since(j.operation.StartTime)
    }

    j.encoder.Encode(map[string]interface{}{
        "event":    "complete",
        "message":  message,
        "duration": duration.Seconds(),
    })
}

func (j *jsonReporter) Error(err error) {
    j.encoder.Encode(map[string]interface{}{
        "event":   "error",
        "message": err.Error(),
    })
}
```

### 4. multi.go (~50 lines)

```go
package progress

import (
    "github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces"
)

type multiReporter struct {
    reporters []interfaces.ProgressReporter
}

// NewMultiReporter creates a composite reporter that forwards to multiple reporters
func NewMultiReporter(reporters ...interfaces.ProgressReporter) interfaces.ProgressReporter {
    return &multiReporter{
        reporters: reporters,
    }
}

func (m *multiReporter) Start(message string) {
    for _, r := range m.reporters {
        r.Start(message)
    }
}

func (m *multiReporter) ReportProgress(current, total int64, message string) {
    for _, r := range m.reporters {
        r.ReportProgress(current, total, message)
    }
}

func (m *multiReporter) Complete(message string) {
    for _, r := range m.reporters {
        r.Complete(message)
    }
}

func (m *multiReporter) Error(err error) {
    for _, r := range m.reporters {
        r.Error(err)
    }
}
```

### 5. formatter.go (~30 lines)

```go
package progress

import (
    "fmt"
    "time"
)

// formatBytes formats byte count for human readability
func formatBytes(bytes int64) string {
    const unit = 1024
    if bytes < unit {
        return fmt.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// formatDuration formats duration for human readability
func formatDuration(d time.Duration) string {
    if d < time.Second {
        return fmt.Sprintf("%dms", d.Milliseconds())
    }
    if d < time.Minute {
        return fmt.Sprintf("%.1fs", d.Seconds())
    }
    return fmt.Sprintf("%.1fm", d.Minutes())
}
```

## Configuration Requirements (R355 Mandatory)

### WRONG - Will fail review:
```go
// ❌ VIOLATION - Hardcoded output location
outputFile := "/var/log/progress.json"

// ❌ VIOLATION - Stub implementation
func (r *reporter) ReportProgress(current, total int64, message string) {
    // TODO: implement progress reporting
}

// ❌ VIOLATION - Static configuration
updateInterval := 100 // milliseconds
```

### CORRECT - Production ready:
```go
// ✅ Configurable output
writer := os.Getenv("PROGRESS_OUTPUT")
if writer == "" {
    writer = os.Stdout
}

// ✅ Full implementation required
func (r *reporter) ReportProgress(current, total int64, message string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    r.current = current
    r.total = total
    r.writeProgress(message)
}

// ✅ No hardcoded values
updateInterval := time.Duration(0) // Immediate updates by default
if v := os.Getenv("PROGRESS_UPDATE_MS"); v != "" {
    if ms, err := strconv.Atoi(v); err == nil {
        updateInterval = time.Duration(ms) * time.Millisecond
    }
}
```

## Testing Requirements

### Unit Test Coverage (60% minimum)
- Test all three reporter implementations
- Test format helper functions
- Test error conditions
- Test concurrent usage safety

### Test File Structure
```
pkg/progress/
├── console_reporter_test.go
├── json_reporter_test.go
├── multi_test.go
└── formatter_test.go
```

### Example Test Case
```go
func TestConsoleReporter(t *testing.T) {
    var buf bytes.Buffer
    reporter := NewConsoleReporterWithWriter(&buf)

    reporter.Start("Test operation")
    reporter.ReportProgress(50, 100, "Processing")
    reporter.Complete("Done")

    output := buf.String()
    if !strings.Contains(output, "Starting: Test operation") {
        t.Errorf("Expected start message, got: %s", output)
    }
    if !strings.Contains(output, "[50.0%]") {
        t.Errorf("Expected progress percentage, got: %s", output)
    }
}
```

## Atomic PR Design

### PR Summary
Single PR implementing progress reporting abstraction layer with console, JSON, and multi-reporter implementations.

### Can Merge to Main Alone
**TRUE** - This effort has no runtime dependencies on other Wave 2 efforts.

### R355 Production Ready Checklist
- ✅ no_hardcoded_values: true
- ✅ all_config_from_env: true (where applicable)
- ✅ no_stub_implementations: true
- ✅ no_todo_markers: true
- ✅ all_functions_complete: true

### Feature Flags Needed
None - This is a foundational interface implementation that will be consumed by Wave 3.

### Interface Implementations
- interface: "interfaces.ProgressReporter"
- implementations:
  - ConsoleReporter (terminal output)
  - JSONReporter (structured logging)
  - MultiReporter (composite pattern)
- production_ready: true
- notes: "All implementations fully functional, not stubs"

### PR Verification
- tests_pass_alone: true
- build_remains_working: true
- no_external_dependencies: true (beyond Wave 1)
- backward_compatible: true

## Dependencies

### Wave 1 Dependencies
- `github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces` - For ProgressReporter interface

### External Dependencies
- Standard library only (fmt, io, encoding/json, time)
- NO external progress bar libraries

## Success Criteria

- ✅ Implements interfaces.ProgressReporter exactly as defined in Wave 1
- ✅ Provides three concrete implementations (console, JSON, multi)
- ✅ All implementations are production-ready (no stubs)
- ✅ Helper functions for common formatting needs
- ✅ 60% test coverage minimum
- ✅ Under 800 lines total
- ✅ Can be merged independently to main
- ✅ No breaking changes to existing code

## Next Wave Dependencies

Wave 3 efforts that will use this progress reporter:
- CLI command implementations (push, pull, list, inspect)
- Build operations progress
- Registry operations progress
- Stack processing progress

## Implementation Order

1. Create types.go with Operation and Result structs
2. Implement console_reporter.go with basic terminal output
3. Implement json_reporter.go for structured logging
4. Implement multi.go for composite reporting
5. Add formatter.go with helper functions
6. Write comprehensive tests for all components
7. Verify implementation matches interface exactly
8. Run line counter to confirm under limit

## Notes for SW Engineer

1. **CRITICAL**: You MUST import the ProgressReporter interface from Wave 1:
   ```go
   import "github.com/cnoe-io/idpbuilder/pkg/cmd/interfaces"
   ```

2. **Thread Safety**: Consider adding mutex protection if reporters might be used concurrently.

3. **Writer Injection**: Allow writers to be injected for testability (especially for ConsoleReporter).

4. **Error Handling**: The Error() method should handle nil errors gracefully.

5. **Zero Values**: Ensure reporters work correctly even if Start() is never called.

## Out of Scope (DO NOT ADD)

- Animation or spinner implementations
- Color output or terminal styling
- Progress persistence across restarts
- Network-based progress reporting
- File-based progress logs
- Progress history tracking
- Complex progress calculations
- Rate limiting or throttling

---

**Plan Created By**: code-reviewer
**Plan Location**: /home/vscode/workspaces/idpbuilder-gitea-push/efforts/phase1/wave2/progress-reporter/.software-factory/IMPLEMENTATION-PLAN-20250928-142532.md
**Review Status**: Ready for SW Engineer Implementation