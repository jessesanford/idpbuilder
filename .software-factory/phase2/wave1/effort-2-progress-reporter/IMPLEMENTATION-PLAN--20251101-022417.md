# Effort 2.1.2: Progress Reporter & Output Formatting - Implementation Plan

**Status**: READY FOR IMPLEMENTATION
**Created**: 2025-11-01T02:24:17Z
**Planner**: Code Reviewer Agent (@agent-code-reviewer)
**Fidelity Level**: EXACT SPECIFICATIONS

---

## 🚨 CRITICAL EFFORT METADATA (R213 - FROM WAVE PLAN)

```yaml
---
effort_metadata:
  effort_id: "2.1.2"
  effort_name: "Progress Reporter & Output Formatting"
  estimated_lines: 300
  dependencies: ["2.1.1"]
  files_touched:
    - "pkg/progress/reporter.go"
    - "pkg/progress/interface.go"
    - "pkg/progress/reporter_test.go"
    - "pkg/cmd/push/push.go"
  branch_name: "idpbuilder-oci-push/phase2/wave1/effort-2-progress-reporter"
  base_branch: "idpbuilder-oci-push/phase2/wave1/effort-1-push-command-core"
  theme: "Enhanced progress tracking and console output formatting"
  scope: "Layer-by-layer progress, thread-safe updates, summary statistics, verbose/normal modes"
  complexity: "medium"
  can_parallelize: false
  parallel_with: []
---
```

---

## Overview

**Purpose**: Replace the basic progress callback from Effort 2.1.1 with a sophisticated reporter that tracks layer-by-layer progress, provides thread-safe updates, and displays formatted summaries.

**Wave Context**: Phase 2, Wave 1 - Second effort in command implementation wave
**Dependencies**: Effort 2.1.1 (Push Command Core) - MUST be complete first
**Callback Signature**: `func(update registry.ProgressUpdate)` from Phase 1 registry package

---

## Scope Definition

### ✅ IN SCOPE (What this effort MUST deliver)

1. **ProgressReporter Interface**
   - Define interface with 3 core methods
   - Match registry.ProgressCallback signature exactly
   - Return callback function for registry.Push()

2. **Thread-Safe Layer Progress Tracking**
   - Track progress for each layer individually
   - Use sync.Mutex for concurrent access
   - Support multiple layers being pushed simultaneously
   - Store layer metadata (digest, size, pushed bytes, status, timing)

3. **Display Modes**
   - **Normal mode**: Compact output (layer digest, percentage, status)
   - **Verbose mode**: Detailed output (full digest, rates, elapsed time, statistics)
   - Mode selection based on --verbose flag from Effort 2.1.1

4. **Summary Statistics**
   - Total layers pushed
   - Layers skipped (already exists)
   - Total bytes transferred
   - Total duration
   - Average transfer rate (MB/s)

5. **Integration with push.go**
   - Replace basic callback from Effort 2.1.1
   - Add progress package import
   - Call reporter.DisplaySummary() before success message
   - Maintain all 8 pipeline stages

6. **Testing**
   - 15 unit tests (T-2.1.2-01 through T-2.1.2-15)
   - Thread safety verification with -race flag
   - ≥85% code coverage

### ❌ OUT OF SCOPE (What this effort MUST NOT include)

1. **Terminal UI/TUI Components** - Not needed for Wave 1 text output
2. **Progress Bars with ANSI Escape Codes** - Future enhancement
3. **Color Output** - Future enhancement
4. **Real-time Updating (Overwriting Lines)** - Wave 2.3
5. **JSON Output Format** - Wave 2.3
6. **Progress Bar Libraries** - Keep dependencies minimal

---

## File Structure

### New Files to Create

#### 1. pkg/progress/interface.go (30 lines)
**Purpose**: Define ProgressReporter interface

**Implementation**:
```go
// Package progress provides progress tracking and reporting for image push operations
package progress

import "github.com/cnoe-io/idpbuilder/pkg/registry"

// ProgressReporter defines operations for tracking and displaying push progress
type ProgressReporter interface {
    // HandleProgress processes a progress update from registry push
    HandleProgress(update registry.ProgressUpdate)

    // DisplaySummary shows final statistics after push completes
    DisplaySummary()

    // GetCallback returns a callback function for registry.Push()
    GetCallback() registry.ProgressCallback
}
```

**Requirements**:
- Interface MUST match registry.ProgressCallback signature exactly
- GetCallback() MUST return a function compatible with registry.Push()
- Keep interface minimal (only 3 methods needed for Wave 1)
- All methods MUST be exported (capitalized)

---

#### 2. pkg/progress/reporter.go (200 lines)
**Purpose**: Concrete implementation of ProgressReporter with thread-safe tracking

**Structure**:
```go
package progress

import (
    "fmt"
    "sync"
    "time"

    "github.com/cnoe-io/idpbuilder/pkg/registry"
)

// Reporter tracks and displays progress for image push operations
type Reporter struct {
    verbose      bool
    startTime    time.Time
    layers       map[string]*LayerProgress
    mu           sync.Mutex
}

// LayerProgress tracks individual layer upload progress
type LayerProgress struct {
    Digest       string
    Size         int64
    Pushed       int64
    Status       string
    StartTime    time.Time
    CompleteTime *time.Time
}
```

**Functions to Implement**:

1. **NewReporter(verbose bool) *Reporter**
   - Initialize Reporter with verbose mode
   - Set startTime to current time
   - Initialize empty layers map
   - Return pointer to Reporter

2. **HandleProgress(update registry.ProgressUpdate)**
   - Lock mutex at start, defer unlock
   - Get or create layer progress entry
   - Update layer state (pushed bytes, status)
   - Set CompleteTime if status is "complete" or "exists"
   - Call displayNormal or displayVerbose based on mode

3. **displayNormal(layer *LayerProgress)** (internal, unexported)
   - Truncate digest to 12 characters for display
   - For "uploading": Show percentage (e.g., "abc123def456: Pushing [45.2%]")
   - For "complete": Show checkmark (e.g., "abc123def456: Pushed ✓")
   - For "exists": Show skip message (e.g., "abc123def456: Already exists ✓")

4. **displayVerbose(layer *LayerProgress)** (internal, unexported)
   - Show full digest
   - For "uploading": Display percentage, bytes ratio, rate (MB/s), elapsed time
   - For "complete": Display completion message with duration
   - For "exists": Display skip message
   - Calculate rates carefully (handle division by zero if elapsed < 0.001s)

5. **DisplaySummary()**
   - Lock mutex, defer unlock
   - Calculate total layers, pushed layers, skipped layers
   - Calculate total bytes transferred
   - Calculate duration since startTime
   - Calculate average rate (MB/s)
   - Display formatted summary:
     ```
     Push Summary:
       Total layers: X
       Pushed: Y, Skipped: Z
       Total size: XX.XX MB
       Duration: XX.Xs
       Average rate: XX.XX MB/s
     ```

6. **GetCallback() registry.ProgressCallback**
   - Return a closure that calls HandleProgress
   - Must match signature: `func(update registry.ProgressUpdate)`

**Critical Implementation Requirements**:
- MUST use sync.Mutex for thread-safe layer map access
- HandleProgress MUST handle concurrent calls correctly (multiple goroutines)
- displayNormal and displayVerbose MUST NOT be exported (lowercase)
- Digest truncation in displayNormal MUST be exactly 12 characters
- Rate calculations MUST handle division by zero gracefully
- Summary MUST calculate correct totals even with mixed status layers
- CompleteTime MUST be a pointer to handle nil case (layer still in progress)

---

#### 3. pkg/progress/reporter_test.go (70 lines)
**Purpose**: Unit tests for progress reporter (15 tests)

**Test Categories**:

1. **Reporter Creation** (T-2.1.2-01 to T-2.1.2-02)
   - TestNewReporter_Normal: Verify reporter created with verbose=false
   - TestNewReporter_Verbose: Verify reporter created with verbose=true

2. **Progress Tracking** (T-2.1.2-03 to T-2.1.2-07)
   - TestReporter_HandleProgress_Uploading: Update with uploading status
   - TestReporter_HandleProgress_Complete: Update from uploading to complete
   - TestReporter_HandleProgress_Exists: Handle already-exists status
   - TestReporter_HandleProgress_MultipleLayers: Track multiple layers simultaneously
   - TestReporter_HandleProgress_ThreadSafety: Concurrent updates with go test -race

3. **Display Formatting** (T-2.1.2-08 to T-2.1.2-10)
   - TestReporter_DisplayNormal_Format: Verify normal mode output format
   - TestReporter_DisplayVerbose_Format: Verify verbose mode output format
   - TestReporter_DisplayVerbose_RateCalculation: Test rate calculation logic

4. **Summary Statistics** (T-2.1.2-11 to T-2.1.2-13)
   - TestReporter_DisplaySummary_SingleLayer: Summary with one layer
   - TestReporter_DisplaySummary_MultipleLayers: Summary with multiple layers
   - TestReporter_DisplaySummary_MixedStatus: Summary with pushed + skipped layers

5. **Integration** (T-2.1.2-14 to T-2.1.2-15)
   - TestReporter_GetCallback: Verify callback matches registry.ProgressCallback
   - TestReporter_DigestTruncation: Verify 12-character truncation

**Test Pattern Example**:
```go
package progress_test

import (
    "fmt"
    "sync"
    "testing"

    "github.com/cnoe-io/idpbuilder/pkg/progress"
    "github.com/cnoe-io/idpbuilder/pkg/registry"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Example: T-2.1.2-04 - TestReporter_HandleProgress_Complete
func TestReporter_HandleProgress_Complete(t *testing.T) {
    // Given: Reporter
    reporter := progress.NewReporter(false)

    // When: Receiving progress updates from uploading to complete
    reporter.HandleProgress(registry.ProgressUpdate{
        LayerDigest: "sha256:abc123def456",
        LayerSize:   1024,
        BytesPushed: 0,
        Status:      "uploading",
    })

    reporter.HandleProgress(registry.ProgressUpdate{
        LayerDigest: "sha256:abc123def456",
        LayerSize:   1024,
        BytesPushed: 1024,
        Status:      "complete",
    })

    // Then: No panics, summary works
    reporter.DisplaySummary()
}

// Example: T-2.1.2-07 - TestReporter_HandleProgress_ThreadSafety
func TestReporter_HandleProgress_ThreadSafety(t *testing.T) {
    // Given: Reporter
    reporter := progress.NewReporter(false)

    // When: Concurrent progress updates
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(layerNum int) {
            defer wg.Done()
            digest := fmt.Sprintf("sha256:layer%d", layerNum)
            reporter.HandleProgress(registry.ProgressUpdate{
                LayerDigest: digest,
                LayerSize:   1000,
                BytesPushed: 500,
                Status:      "uploading",
            })
        }(i)
    }

    wg.Wait()

    // Then: No race conditions, all updates processed
    reporter.DisplaySummary()
    // Race detector would catch issues: go test -race
}
```

**Test Coverage Requirements**:
- Minimum 85% code coverage
- All 15 tests from test plan MUST be implemented
- Thread safety MUST be verified with -race flag
- Use testify/assert for assertions
- Use testify/require for critical failures

---

### Modifications to Existing Files

#### 4. pkg/cmd/push/push.go (~10 lines changed)
**Purpose**: Replace basic callback with reporter

**Location**: In runPush() function, Stage 7 (lines ~192-202 in Effort 2.1.1)

**Original Code** (from Effort 2.1.1):
```go
// STAGE 7: Create progress callback (basic implementation for Wave 1)
progressCallback := func(update registry.ProgressUpdate) {
    if opts.Verbose {
        fmt.Printf("Layer %s: %d/%d bytes (%s)\n",
            truncateDigest(update.LayerDigest, 12),
            update.BytesPushed,
            update.LayerSize,
            update.Status)
    }
}
```

**New Code** (for this effort):
```go
// Add import at top of file
import "github.com/cnoe-io/idpbuilder/pkg/progress"

// STAGE 7: Create progress reporter (replaces basic callback)
reporter := progress.NewReporter(opts.Verbose)

// STAGE 8: Execute push with reporter (modify this stage)
fmt.Printf("Pushing to %s...\n", targetRef)
if err := registryClient.Push(ctx, image, targetRef, reporter.GetCallback()); err != nil {
    return fmt.Errorf("push failed: %w", err)
}

// Display final summary (NEW - add before success message)
reporter.DisplaySummary()

fmt.Printf("✓ Successfully pushed %s to %s\n", opts.ImageName, opts.Registry)
return nil
```

**Requirements**:
- Add progress import at top of file
- Replace lines 193-200 (basic callback) with reporter creation
- Modify line 204 to use reporter.GetCallback() instead of progressCallback
- Add reporter.DisplaySummary() call after push succeeds (line ~208)
- DO NOT change other parts of runPush() pipeline
- Keep modifications minimal (~10 lines total: 1 import, 9 in function)
- Remove truncateDigest() function (no longer needed - reporter handles it)

---

## Implementation Steps (Sequential)

### Step 1: Create Package Structure (5 minutes)

```bash
# Create package directory
mkdir -p pkg/progress

# Verify directory structure
ls -la pkg/progress
```

**Verification**: Directory exists and is empty

---

### Step 2: Implement interface.go (10 minutes)

**File**: `pkg/progress/interface.go`

**Implementation Order**:
1. Create package declaration
2. Import registry package for ProgressUpdate and ProgressCallback types
3. Define ProgressReporter interface with 3 methods
4. Add godoc comments for interface and each method

**Verification**:
- File compiles without errors
- Interface exported (capitalized)
- All methods exported
- Godoc comments present

---

### Step 3: Implement reporter.go - Data Structures (15 minutes)

**File**: `pkg/progress/reporter.go`

**Implementation Order**:
1. Package declaration and imports
2. Define Reporter struct with 4 fields
3. Define LayerProgress struct with 6 fields
4. Add godoc comments for both structs and all fields

**Code Template**:
```go
package progress

import (
    "fmt"
    "sync"
    "time"

    "github.com/cnoe-io/idpbuilder/pkg/registry"
)

// Reporter tracks and displays progress for image push operations
type Reporter struct {
    verbose      bool                        // Enable detailed output
    startTime    time.Time                   // Start time for duration calculation
    layers       map[string]*LayerProgress   // Layer digest -> progress
    mu           sync.Mutex                  // Protects layers map
}

// LayerProgress tracks individual layer upload progress
type LayerProgress struct {
    Digest       string     // Layer digest (e.g., sha256:abc123...)
    Size         int64      // Total layer size in bytes
    Pushed       int64      // Bytes pushed so far
    Status       string     // Current status (uploading, complete, exists)
    StartTime    time.Time  // When layer push started
    CompleteTime *time.Time // When layer push completed (nil if in progress)
}
```

**Verification**:
- Structs compile
- Field types correct
- Mutex field present for thread safety

---

### Step 4: Implement reporter.go - NewReporter (10 minutes)

**Function**: `NewReporter(verbose bool) *Reporter`

**Code Template**:
```go
// NewReporter creates a progress reporter
func NewReporter(verbose bool) *Reporter {
    return &Reporter{
        verbose:   verbose,
        startTime: time.Now(),
        layers:    make(map[string]*LayerProgress),
    }
}
```

**Verification**:
- Function exported
- Returns initialized Reporter
- Godoc comment present

---

### Step 5: Implement reporter.go - HandleProgress (30 minutes)

**Function**: `HandleProgress(update registry.ProgressUpdate)`

**Implementation Order**:
1. Lock mutex, defer unlock
2. Check if layer exists in map, create if not
3. Update layer fields (pushed, status)
4. Set CompleteTime if status is complete/exists
5. Call display function based on verbose mode

**Code Template**:
```go
// HandleProgress processes a progress update from registry.Push()
// This matches the registry.ProgressCallback signature from Phase 1
func (r *Reporter) HandleProgress(update registry.ProgressUpdate) {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Get or create layer progress
    layer, exists := r.layers[update.LayerDigest]
    if !exists {
        layer = &LayerProgress{
            Digest:    update.LayerDigest,
            Size:      update.LayerSize,
            StartTime: time.Now(),
        }
        r.layers[update.LayerDigest] = layer
    }

    // Update layer state
    layer.Pushed = update.BytesPushed
    layer.Status = update.Status
    if update.Status == "complete" || update.Status == "exists" {
        now := time.Now()
        layer.CompleteTime = &now
    }

    // Display progress
    if r.verbose {
        r.displayVerbose(layer)
    } else {
        r.displayNormal(layer)
    }
}
```

**Verification**:
- Mutex correctly guards map access
- Layer creation logic correct
- CompleteTime pointer handling correct
- Display routing works

---

### Step 6: Implement reporter.go - Display Functions (40 minutes)

**Functions**:
- `displayNormal(layer *LayerProgress)` (internal)
- `displayVerbose(layer *LayerProgress)` (internal)

**Implementation Order**:

1. **displayNormal** (15 minutes):
   - Truncate digest to 12 characters
   - Switch on status
   - Format output for uploading (percentage)
   - Format output for complete (checkmark)
   - Format output for exists (skip message)

2. **displayVerbose** (25 minutes):
   - Calculate elapsed time
   - Switch on status
   - For uploading: show full digest, percentage, bytes, rate, elapsed
   - For complete: show completion message with duration
   - For exists: show skip message
   - Handle edge cases (elapsed too small, division by zero)

**Code Templates**:

```go
// displayNormal shows compact progress (for normal mode)
func (r *Reporter) displayNormal(layer *LayerProgress) {
    shortDigest := layer.Digest
    if len(shortDigest) > 12 {
        shortDigest = shortDigest[:12]
    }

    switch layer.Status {
    case "uploading":
        percent := float64(layer.Pushed) / float64(layer.Size) * 100
        fmt.Printf("  %s: Pushing [%.1f%%]\n", shortDigest, percent)
    case "complete":
        fmt.Printf("  %s: Pushed ✓\n", shortDigest)
    case "exists":
        fmt.Printf("  %s: Already exists ✓\n", shortDigest)
    }
}

// displayVerbose shows detailed progress (for --verbose)
func (r *Reporter) displayVerbose(layer *LayerProgress) {
    elapsed := time.Since(layer.StartTime).Seconds()

    switch layer.Status {
    case "uploading":
        percent := float64(layer.Pushed) / float64(layer.Size) * 100
        rate := 0.0
        if elapsed > 0.001 {  // Avoid division by zero
            rate = float64(layer.Pushed) / elapsed / 1024 / 1024 // MB/s
        }
        fmt.Printf("  Layer %s:\n", layer.Digest)
        fmt.Printf("    Status: Uploading [%.1f%%] (%d/%d bytes)\n",
            percent, layer.Pushed, layer.Size)
        fmt.Printf("    Rate: %.2f MB/s\n", rate)
        fmt.Printf("    Elapsed: %.1fs\n", elapsed)
    case "complete":
        duration := elapsed
        if layer.CompleteTime != nil {
            duration = layer.CompleteTime.Sub(layer.StartTime).Seconds()
        }
        fmt.Printf("  Layer %s: Complete (%.1fs)\n", layer.Digest, duration)
    case "exists":
        fmt.Printf("  Layer %s: Already exists (skipped)\n", layer.Digest)
    }
}
```

**Verification**:
- Functions are unexported (lowercase)
- Digest truncation works
- Division by zero handled
- Output formats match specification

---

### Step 7: Implement reporter.go - DisplaySummary (25 minutes)

**Function**: `DisplaySummary()`

**Implementation Order**:
1. Lock mutex, defer unlock
2. Initialize counters (total, pushed, skipped, bytes)
3. Iterate through layers, calculate totals
4. Calculate duration and rate
5. Format and print summary

**Code Template**:
```go
// DisplaySummary shows final statistics after push completes
func (r *Reporter) DisplaySummary() {
    r.mu.Lock()
    defer r.mu.Unlock()

    totalLayers := len(r.layers)
    totalBytes := int64(0)
    pushedLayers := 0
    skippedLayers := 0

    for _, layer := range r.layers {
        totalBytes += layer.Size
        if layer.Status == "complete" {
            pushedLayers++
        } else if layer.Status == "exists" {
            skippedLayers++
        }
    }

    duration := time.Since(r.startTime)
    avgRate := 0.0
    if duration.Seconds() > 0.001 {  // Avoid division by zero
        avgRate = float64(totalBytes) / duration.Seconds() / 1024 / 1024 // MB/s
    }

    fmt.Println("\nPush Summary:")
    fmt.Printf("  Total layers: %d\n", totalLayers)
    fmt.Printf("  Pushed: %d, Skipped: %d\n", pushedLayers, skippedLayers)
    fmt.Printf("  Total size: %.2f MB\n", float64(totalBytes)/1024/1024)
    fmt.Printf("  Duration: %.1fs\n", duration.Seconds())
    fmt.Printf("  Average rate: %.2f MB/s\n", avgRate)
}
```

**Verification**:
- Mutex protects map iteration
- All counters calculated correctly
- Division by zero handled
- Output format matches specification

---

### Step 8: Implement reporter.go - GetCallback (10 minutes)

**Function**: `GetCallback() registry.ProgressCallback`

**Code Template**:
```go
// GetCallback returns a registry.ProgressCallback that uses this reporter
func (r *Reporter) GetCallback() registry.ProgressCallback {
    return func(update registry.ProgressUpdate) {
        r.HandleProgress(update)
    }
}
```

**Verification**:
- Returns correct function signature
- Closure captures reporter correctly
- Matches registry.ProgressCallback type

---

### Step 9: Implement Tests (60 minutes)

**File**: `pkg/progress/reporter_test.go`

**Test Implementation Order**:

1. **Setup** (5 minutes):
   - Package declaration
   - Imports (testing, testify, sync, registry)
   - Helper functions if needed

2. **Reporter Creation Tests** (10 minutes) - T-2.1.2-01, T-2.1.2-02:
   - TestNewReporter_Normal
   - TestNewReporter_Verbose

3. **Progress Tracking Tests** (20 minutes) - T-2.1.2-03 to T-2.1.2-07:
   - TestReporter_HandleProgress_Uploading
   - TestReporter_HandleProgress_Complete
   - TestReporter_HandleProgress_Exists
   - TestReporter_HandleProgress_MultipleLayers
   - TestReporter_HandleProgress_ThreadSafety (with sync.WaitGroup)

4. **Display Tests** (15 minutes) - T-2.1.2-08 to T-2.1.2-10:
   - TestReporter_DisplayNormal_Format (capture stdout if possible)
   - TestReporter_DisplayVerbose_Format
   - TestReporter_DisplayVerbose_RateCalculation

5. **Summary Tests** (10 minutes) - T-2.1.2-11 to T-2.1.2-13:
   - TestReporter_DisplaySummary_SingleLayer
   - TestReporter_DisplaySummary_MultipleLayers
   - TestReporter_DisplaySummary_MixedStatus

6. **Integration Tests** (10 minutes) - T-2.1.2-14 to T-2.1.2-15:
   - TestReporter_GetCallback
   - TestReporter_DigestTruncation

**Verification**:
```bash
# Run tests
go test ./pkg/progress/... -v

# Check coverage
go test ./pkg/progress/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
# Must show ≥85%

# Thread safety check
go test ./pkg/progress/... -race -v
# Must show 0 data races
```

---

### Step 10: Modify push.go (15 minutes)

**File**: `pkg/cmd/push/push.go`

**Modification Order**:
1. Add progress import at top
2. Locate Stage 7 (progress callback)
3. Replace basic callback with reporter creation
4. Modify Stage 8 to use reporter.GetCallback()
5. Add reporter.DisplaySummary() call
6. Remove truncateDigest() function (no longer needed)

**Verification**:
```bash
# Compile check
go build ./pkg/cmd/push/...

# Verify push command still works
go test ./pkg/cmd/push/... -v
```

---

## Size Management

### Line Count Tracking

**Tool**: `$PROJECT_ROOT/tools/line-counter.sh`

**Measurement Points**:
1. After implementing interface.go (~30 lines)
2. After implementing reporter.go structs (~50 lines)
3. After implementing NewReporter + HandleProgress (~100 lines)
4. After implementing display functions (~170 lines)
5. After implementing DisplaySummary + GetCallback (~230 lines)
6. After implementing tests (~300 lines)
7. After modifying push.go (~310 lines)
8. Final measurement before review

**Commands**:
```bash
# Find project root
PROJECT_ROOT=$(pwd)
while [ "$PROJECT_ROOT" != "/" ]; do
    [ -f "$PROJECT_ROOT/orchestrator-state-v3.json" ] && break
    PROJECT_ROOT=$(dirname "$PROJECT_ROOT")
done

# Measure implementation lines (R221: CD first!)
cd /home/vscode/workspaces/idpbuilder-oci-push-planning/efforts/phase2/wave1/effort-2-progress-reporter
$PROJECT_ROOT/tools/line-counter.sh

# Expected output: ~300 lines
```

**Thresholds**:
- ✅ <700 lines: COMPLIANT (continue)
- ⚠️ 700-800 lines: WARNING (approaching limit, finish and review)
- ❌ ≥800 lines: STOP IMMEDIATELY (requires split plan)

---

## Testing Strategy

### Unit Tests (15 tests - No External Dependencies)

**Run Command**:
```bash
go test -short ./pkg/progress/... -v
```

**Expected Results**:
- 15 tests passing
- 0 failures
- <2 seconds execution time

### Thread Safety Verification

**Command**:
```bash
go test ./pkg/progress/... -race -v
```

**Expected Results**:
- All tests passing
- 0 data races detected
- Race detector confirms thread safety

### Coverage Verification

**Command**:
```bash
go test ./pkg/progress/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Check coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if (( $(echo "$COVERAGE < 85" | bc -l) )); then
    echo "❌ Coverage below 85%: $COVERAGE%"
    exit 1
else
    echo "✅ Coverage: $COVERAGE%"
fi
```

**Minimum**: 85% coverage required

### Integration with push.go

**Manual Test**:
```bash
# Build idpbuilder
go build -o idpbuilder ./...

# Test with normal mode (compact output)
./idpbuilder push alpine:latest --username admin --password password --insecure

# Test with verbose mode (detailed output)
./idpbuilder push alpine:latest --username admin --password password --insecure --verbose

# Verify summary displays after push completes
```

**Expected Output** (normal mode):
```
Pushing to gitea.cnoe.localtest.me:8443/alpine:latest...
  abc123def456: Pushing [25.5%]
  abc123def456: Pushing [78.3%]
  abc123def456: Pushed ✓
  def789ghi012: Already exists ✓
  ...

Push Summary:
  Total layers: 5
  Pushed: 3, Skipped: 2
  Total size: 7.45 MB
  Duration: 2.3s
  Average rate: 3.24 MB/s
✓ Successfully pushed alpine:latest to gitea.cnoe.localtest.me:8443
```

---

## Phase 1 Integration

### Registry Client Integration

**Interface** (from Phase 1 Wave 2):
```go
// ProgressCallback is called during layer push to report progress
type ProgressCallback func(update ProgressUpdate)

// ProgressUpdate contains information about layer push progress
type ProgressUpdate struct {
    LayerDigest string  // Layer digest (e.g., sha256:abc123...)
    LayerSize   int64   // Total layer size in bytes
    BytesPushed int64   // Bytes pushed so far
    Status      string  // Current status (uploading, complete, exists)
}
```

**Usage in reporter.go**:
```go
// HandleProgress matches the ProgressCallback signature
func (r *Reporter) HandleProgress(update registry.ProgressUpdate) {
    // Process update...
}

// GetCallback returns a ProgressCallback
func (r *Reporter) GetCallback() registry.ProgressCallback {
    return func(update registry.ProgressUpdate) {
        r.HandleProgress(update)
    }
}
```

**Integration Point**: The reporter MUST work with registry.Client.Push() from Phase 1:
```go
// From Effort 2.1.1 push.go
if err := registryClient.Push(ctx, image, targetRef, reporter.GetCallback()); err != nil {
    return fmt.Errorf("push failed: %w", err)
}
```

---

## Dependencies

### Upstream Dependencies (MUST be complete)

**Effort 2.1.1 (Push Command Core)** - ✅ REQUIRED:
- Defines PushOptions with Verbose flag
- Defines runPush() pipeline with Stage 7 integration point
- Provides callback signature pattern
- This effort MUST be complete and merged before starting Effort 2.1.2

**Phase 1 (Registry Package)** - ✅ Complete:
- registry.ProgressCallback type definition
- registry.ProgressUpdate struct definition

### Downstream Dependencies

**None** - This completes Wave 2.1

**Future Enhancements** (Wave 2.2 and beyond):
- Wave 2.2 will enhance reporter with registry auto-detection feedback
- Wave 2.3 may add JSON output format
- Future waves may add progress bars or TUI

---

## Acceptance Criteria

### Implementation Checklist

- [ ] All 4 files created/modified as specified
  - [ ] pkg/progress/interface.go (30 lines)
  - [ ] pkg/progress/reporter.go (200 lines)
  - [ ] pkg/progress/reporter_test.go (70 lines)
  - [ ] pkg/cmd/push/push.go (~10 lines modified)

- [ ] All 15 tests from test plan implemented and passing
  - [ ] T-2.1.2-01 to T-2.1.2-02: Reporter creation tests
  - [ ] T-2.1.2-03 to T-2.1.2-07: Progress tracking tests
  - [ ] T-2.1.2-08 to T-2.1.2-10: Display formatting tests
  - [ ] T-2.1.2-11 to T-2.1.2-13: Summary statistics tests
  - [ ] T-2.1.2-14 to T-2.1.2-15: Integration tests

- [ ] Code quality
  - [ ] Code coverage ≥85%
  - [ ] Thread safety verified with `go test -race` (zero data races)
  - [ ] No linting errors (golangci-lint)
  - [ ] All exported functions/types have godoc comments
  - [ ] Line count within estimate (300 ± 45 lines, 15% tolerance)

- [ ] Manual testing
  - [ ] `idpbuilder push alpine:latest --verbose` shows detailed progress
  - [ ] `idpbuilder push alpine:latest` shows compact progress
  - [ ] Summary displays correct statistics (layers, bytes, duration, rate)
  - [ ] Integration with push.go works correctly (verified by end-to-end test)

- [ ] Integration verification
  - [ ] Reporter works with registry.Client.Push()
  - [ ] Callback signature matches exactly
  - [ ] No changes to other pipeline stages
  - [ ] Effort 2.1.1 tests still pass after modification

---

## Risk Mitigation

### High-Risk Areas

1. **Thread Safety**
   - **Risk**: Concurrent progress updates could cause race conditions
   - **Mitigation**: Use sync.Mutex, verify with -race flag, stress tests
   - **Detection**: go test -race detects data races

2. **Division by Zero**
   - **Risk**: Rate calculation if elapsed time = 0
   - **Mitigation**: Check elapsed > 0.001s before division
   - **Detection**: Unit tests with immediate completion

3. **Display Output Verification**
   - **Risk**: Hard to test console output in automated tests
   - **Mitigation**: Manual testing required, verify format with examples
   - **Detection**: Visual inspection during manual testing

4. **Integration with Effort 2.1.1**
   - **Risk**: Breaking changes to push.go pipeline
   - **Mitigation**: Minimize modifications, keep changes localized to Stage 7
   - **Detection**: Effort 2.1.1 tests must still pass

---

## Quality Gates (R502 Compliance)

### Implementation Plan Quality

- [x] **R213 metadata**: Complete effort metadata from wave plan
- [x] **Exact code specifications**: Real Go code (not pseudocode) for all implementations
- [x] **Complete file lists**: All 4 files specified with line counts
- [x] **Detailed task breakdowns**: 10 sequential implementation steps
- [x] **Test specifications**: References to 15 concrete tests in test plan
- [x] **Dependency documentation**: Integration with Effort 2.1.1 and Phase 1
- [x] **Scope clarity**: Clear IN SCOPE and OUT OF SCOPE sections
- [x] **Size management**: Line counting strategy with thresholds

---

## Next Steps (After Implementation)

1. **Size Measurement**: Run line-counter.sh to verify <800 lines
2. **Code Review**: Spawn Code Reviewer agent for review
3. **Integration**: Merge to wave integration branch after approval
4. **Wave Complete**: Wave 2.1 is complete (both efforts done)
5. **Wave Integration Tests**: Run full wave test suite
6. **Architect Assessment**: Spawn Architect for R340 wave review

---

## Document Status

**Status**: ✅ READY FOR SOFTWARE ENGINEER
**Created**: 2025-11-01T02:24:17Z
**Planner**: Code Reviewer Agent (@agent-code-reviewer)
**Estimated Lines**: 300 lines
**Fidelity Level**: EXACT (step-by-step implementation guidance)

**Compliance Verified**:
- ✅ R213: Effort metadata complete
- ✅ R219: Dependency effort plan analyzed (read Effort 2.1.1 plan)
- ✅ R303: Timestamped filename (IMPLEMENTATION-PLAN--20251101-022417.md)
- ✅ R383: Stored in .software-factory/phase2/wave1/effort-2-progress-reporter/
- ✅ R502: Implementation plan quality gates (exact specifications)
- ✅ R341: References test plan (15 tests)

**Next Action**: Orchestrator spawns Software Engineer with this plan

---

**🚨 CRITICAL R405 AUTOMATION FLAG 🚨**

CONTINUE-SOFTWARE-FACTORY=TRUE

---

**END OF IMPLEMENTATION PLAN**
