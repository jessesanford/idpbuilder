# SPLIT-PLAN-001.md
## Split 001 of 3: Core Infrastructure (Logging & Progress)
**Planner**: Code Reviewer Agent
**Parent Effort**: E1.2.3-image-push-operations

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: None (first split of THIS effort)
  - Path: N/A (this is Split 001)
  - Branch: N/A
- **This Split**: Split 001 of phase1/wave2/E1.2.3-image-push-operations
  - Path: efforts/phase1/wave2/E1.2.3-image-push-operations/split-001/
  - Branch: phase1/wave2/image-push-operations-split-001
- **Next Split**: Split 002 of phase1/wave2/E1.2.3-image-push-operations
  - Path: efforts/phase1/wave2/E1.2.3-image-push-operations/split-002/
  - Branch: phase1/wave2/image-push-operations-split-002
- **File Boundaries**:
  - This Split Start: pkg/push/logging.go (line 1)
  - This Split End: pkg/push/progress.go (line 303)
  - Next Split Start: pkg/push/discovery.go (line 1)

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- **pkg/push/logging.go** (249 lines)
  - PushLogger wrapper for logr.Logger
  - Push-specific logging methods
  - Structured logging for push operations
- **pkg/push/progress.go** (303 lines)
  - ProgressReporter interface definition
  - ConsoleProgressReporter implementation
  - Layer and image progress tracking

**Total Lines**: 552 lines (well within 800 line limit)

## Functionality
### Logging Infrastructure
- **PushLogger**: Wrapper around logr.Logger with push-specific methods
  - LogPushStart: Log beginning of push operation
  - LogPushComplete: Log successful completion with metrics
  - LogPushFailed: Log failures with error details
  - LogLayerProgress: Log individual layer progress
  - LogRetryAttempt: Log retry attempts with backoff info
  - LogDiscoveryStart/Complete: Log image discovery operations
  - LogAuthenticationAttempt: Log registry authentication

### Progress Reporting
- **ProgressReporter Interface**: Contract for progress implementations
  - StartImage/FinishImage: Track image-level progress
  - UpdateLayer/FinishLayer: Track layer-level progress
  - SetError: Report errors during operations
- **ConsoleProgressReporter**: Console output implementation
  - Real-time progress bars
  - Throughput calculations
  - Multi-image concurrent progress tracking
  - Human-readable size formatting

## Dependencies
- **External Libraries**:
  - github.com/go-logr/logr (for logging interface)
  - Standard library only for progress reporting
- **No dependencies on other splits** (foundational split)

## Implementation Instructions
1. **Setup**:
   ```bash
   cd efforts/phase1/wave2/E1.2.3-image-push-operations
   git checkout -b phase1/wave2/image-push-operations-split-001
   ```

2. **Create pkg/push directory**:
   ```bash
   mkdir -p pkg/push
   ```

3. **Implement logging.go**:
   - Create PushLogger struct wrapping logr.Logger
   - Implement all logging methods with structured fields
   - Add log levels (Debug, Info, Warn, Error)
   - Include timing and metrics in log entries

4. **Implement progress.go**:
   - Define ProgressReporter interface
   - Implement ConsoleProgressReporter with mutex for thread safety
   - Add progress tracking for images and layers
   - Implement human-readable formatting helpers

5. **Testing**:
   - Unit tests for all logging methods
   - Unit tests for progress reporter
   - Mock implementations for testing

6. **Validation**:
   ```bash
   # Ensure compilation
   go build ./pkg/push/...

   # Run tests
   go test ./pkg/push/...

   # Measure size
   ${PROJECT_ROOT}/tools/line-counter.sh
   ```

## Split Branch Strategy
- Branch: `phase1/wave2/image-push-operations-split-001`
- Base: `phase1/wave2/image-push-operations`
- Must merge to: `phase1/wave2/image-push-operations` after review

## Success Criteria
- [ ] Both files compile without errors
- [ ] All types and interfaces properly defined
- [ ] No external dependencies beyond approved libraries
- [ ] Unit tests passing
- [ ] Implementation size <800 lines (target: ~550 lines)
- [ ] No references to discovery, pusher, or operations code

## Notes for SW Engineer
- This is the foundation split - focus on clean interfaces
- Keep ProgressReporter interface simple but extensible
- PushLogger should wrap, not replace, the base logger
- Thread safety is critical for ConsoleProgressReporter
- Consider using sync.Mutex for concurrent access protection
- Human-readable formatting should handle bytes, KB, MB, GB