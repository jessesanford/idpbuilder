# SPLIT-PLAN-003.md
## Split 003 of 3: Main Operation Orchestration
**Planner**: Code Reviewer Agent
**Parent Effort**: E1.2.3-image-push-operations

<!-- ORCHESTRATOR METADATA PLACEHOLDER - DO NOT REMOVE -->
<!-- The orchestrator will add infrastructure metadata below: -->
<!-- WORKING_DIRECTORY, BRANCH, REMOTE, BASE_BRANCH, etc. -->
<!-- SW Engineers MUST read this metadata to navigate to the correct directory -->
<!-- END PLACEHOLDER -->

## Boundaries (CRITICAL: All splits MUST reference SAME effort!)
- **Previous Split**: Split 002 of phase1/wave2/E1.2.3-image-push-operations
  - Path: efforts/phase1/wave2/E1.2.3-image-push-operations/split-002/
  - Branch: phase1/wave2/image-push-operations-split-002
  - Summary: Implemented discovery and pusher components
- **This Split**: Split 003 of phase1/wave2/E1.2.3-image-push-operations
  - Path: efforts/phase1/wave2/E1.2.3-image-push-operations/split-003/
  - Branch: phase1/wave2/image-push-operations-split-003
- **Next Split**: None (final split of THIS effort)

## Files in This Split (EXCLUSIVE - no overlap with other splits)
- **pkg/push/operations.go** (389 lines)
  - PushOperation configuration struct
  - PushOperationResult with metrics
  - Command integration methods
  - Main orchestration logic
  - Filter criteria implementation

**Total Lines**: 389 lines (well within 800 line limit)

## Functionality
### Main Operation Orchestration (operations.go)
- **PushOperation struct**: Complete operation configuration
  - Registry credentials and URL
  - Authentication configuration
  - Transport settings (insecure, TLS)
  - Concurrency and retry settings
  - Filter criteria for selective pushing
  - Integration of logger and progress reporter

- **PushOperationResult struct**: Operation results
  - Timing metrics (start, end, duration)
  - Image counts (found, pushed, failed)
  - Total bytes transferred
  - Individual push results
  - Aggregated error collection
  - Performance metrics

- **Core Methods**:
  - NewPushOperationFromCommand: Create from cobra command
  - Execute: Main execution orchestrator
  - ExecuteWithContext: Context-aware execution
  - ValidateConfiguration: Pre-flight checks
  - ApplyFilters: Image filtering logic
  - CollectMetrics: Performance measurement

- **Integration Functions**:
  - Orchestrates discovery from Split 002
  - Coordinates pusher from Split 002
  - Uses logging from Split 001
  - Manages progress reporting from Split 001
  - Handles concurrent operations
  - Aggregates results and errors

## Dependencies
- **From Split 001**:
  - PushLogger (for operation logging)
  - ProgressReporter (for operation progress)
- **From Split 002**:
  - DiscoverLocalImages functions
  - ImagePusher implementation
  - LocalImage struct
- **External Libraries**:
  - github.com/spf13/cobra (for command integration)
  - github.com/google/go-containerregistry/pkg/authn
  - github.com/google/go-containerregistry/pkg/name
  - Standard library for concurrency (sync, context)

## Implementation Instructions
1. **Setup**:
   ```bash
   cd efforts/phase1/wave2/E1.2.3-image-push-operations
   git checkout phase1/wave2/image-push-operations-split-002
   git checkout -b phase1/wave2/image-push-operations-split-003
   ```

2. **Implement operations.go**:
   - Define PushOperation struct with all fields
   - Define PushOperationResult with metrics
   - Implement NewPushOperationFromCommand for CLI integration
   - Create Execute method as main orchestrator
   - Implement filter criteria logic
   - Add concurrent push coordination
   - Implement result aggregation
   - Add metric collection

3. **Orchestration Logic**:
   ```go
   // Pseudo-code for main execution flow
   func (po *PushOperation) Execute(ctx context.Context) (*PushOperationResult, error) {
       // 1. Validate configuration
       // 2. Setup authentication
       // 3. Discover images (from Split 002)
       // 4. Apply filters
       // 5. Create pusher (from Split 002)
       // 6. Setup progress reporter (from Split 001)
       // 7. Push images concurrently with rate limiting
       // 8. Collect results and metrics
       // 9. Return aggregated result
   }
   ```

4. **Concurrency Management**:
   - Use worker pool pattern for concurrent pushes
   - Implement rate limiting with semaphore
   - Handle context cancellation properly
   - Aggregate errors from workers

5. **Testing**:
   - Unit tests for orchestration logic
   - Integration tests with all splits combined
   - Test concurrent execution scenarios
   - Test error aggregation

6. **Validation**:
   ```bash
   # Ensure compilation with all splits
   go build ./pkg/push/...

   # Run tests
   go test ./pkg/push/...

   # Measure size
   ${PROJECT_ROOT}/tools/line-counter.sh
   ```

## Split Branch Strategy
- Branch: `phase1/wave2/image-push-operations-split-003`
- Base: `phase1/wave2/image-push-operations-split-002`
- Must merge to: `phase1/wave2/image-push-operations` after review

## Success Criteria
- [ ] File compiles without errors
- [ ] Proper integration with Split 001 and 002
- [ ] Main execution flow works end-to-end
- [ ] Concurrent operations properly coordinated
- [ ] Command integration functional
- [ ] Metrics and results properly collected
- [ ] Implementation size <800 lines (target: ~390 lines)
- [ ] Complete functionality when combined with other splits

## Integration Testing Requirements
After this split is complete, test the full effort:
1. All three splits merged into parent branch
2. Full end-to-end push operation works
3. Progress reporting displays correctly
4. Logging captures all operations
5. Concurrent pushes execute properly
6. Error handling works throughout

## Notes for SW Engineer
- This is the orchestration layer - tie everything together
- Import and use components from Split 001 and 002
- Focus on clean orchestration, not reimplementation
- Worker pool pattern recommended for concurrency
- Use context for cancellation support
- Ensure proper error aggregation from concurrent operations
- Filter criteria should be extensible
- Consider adding dry-run support
- Metrics should be useful for monitoring
- Command integration should follow cobra patterns