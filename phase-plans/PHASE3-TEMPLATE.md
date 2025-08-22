# Phase 3: Core Implementation - Detailed Implementation Plan

## Phase Overview
**Duration:** [X] days  
**Critical Path:** [YES/NO] - Main business logic implementation  
**Base Branch:** `phase2-integration` (Phase 2 complete)  
**Target Integration Branch:** `phase3-integration`

---

## Wave 3.1: Core Business Logic

### E3.1.1: Sync Engine
**Branch:** `/phase3/wave1/effort1-sync-engine`  
**Duration:** [X] hours  
**Agent:** Single agent
**Dependencies:** Phase 2 controller framework

#### Source Material:
```markdown
# If porting existing sync logic
- Primary: `origin/feature/[sync-implementation]`
- Reference: `[external-project]/sync/`

# Cherry-pick specific algorithms if available
```

#### Requirements:
1. **MUST** implement bidirectional sync:
   - Source → Target sync
   - Conflict resolution
   - State tracking
   - Error recovery

2. **MUST** handle:
   - Partial failures
   - Network interruptions
   - Large datasets (pagination)
   - Concurrent modifications

#### Complex Implementation Warning:
```yaml
# This effort may exceed 800 lines if it includes:
complexity_factors:
  - State machine with 5+ states
  - Bidirectional sync logic
  - Conflict resolution algorithms
  - Retry mechanisms
  
if_exceeds_limit:
  consider_split:
    - Part 1: State machine core
    - Part 2: Sync algorithms
    - Part 3: Conflict resolution
    
  or_request_exception:
    reason: "Complex state machine - splitting breaks atomicity"
    justification: "States share transition logic"
```

#### Test Requirements (TDD):
```go
// test/sync/engine_test.go
func TestSyncEngine(t *testing.T) {
    t.Run("bidirectional_sync", func(t *testing.T) {
        source := NewMockSource(initialData)
        target := NewMockTarget()
        engine := NewSyncEngine(source, target)
        
        // When
        err := engine.Sync(context.TODO())
        
        // Then
        assert.NoError(t, err)
        assert.Equal(t, source.Data(), target.Data())
    })
    
    t.Run("conflict_resolution", func(t *testing.T) {
        // Test various conflict scenarios
        conflicts := []ConflictCase{
            {Type: "concurrent_update", Resolution: "last_write_wins"},
            {Type: "delete_vs_update", Resolution: "delete_wins"},
        }
        
        for _, c := range conflicts {
            testConflictResolution(t, c)
        }
    })
    
    t.Run("partial_failure_recovery", func(t *testing.T) {
        // Simulate failure midway through sync
        // Verify recovery and continuation
    })
}

// State machine tests
func TestStateMachine(t *testing.T) {
    states := []State{
        StateIdle,
        StateSyncing, 
        StateResolving,
        StateApplying,
        StateComplete,
    }
    
    transitions := []Transition{
        {From: StateIdle, Event: EventStart, To: StateSyncing},
        {From: StateSyncing, Event: EventConflict, To: StateResolving},
        // ... test all valid transitions
    }
    
    testAllTransitions(t, states, transitions)
}
```

#### Pseudo-Code Implementation:
```
FUNCTION implement_sync_engine():
    // Step 1: State machine
    STATE_MACHINE = {
        IDLE: {
            on_start: -> DISCOVERING
        },
        DISCOVERING: {
            on_resources_found: -> COMPARING
            on_error: -> ERROR
        },
        COMPARING: {
            on_no_diff: -> IDLE
            on_diff_found: -> SYNCING
        },
        SYNCING: {
            on_conflict: -> RESOLVING
            on_success: -> VALIDATING
        },
        RESOLVING: {
            on_resolved: -> SYNCING
            on_manual_required: -> WAITING
        },
        VALIDATING: {
            on_valid: -> IDLE
            on_invalid: -> ERROR
        }
    }
    
    // Step 2: Sync algorithm
    FUNCTION sync():
        source_items = FETCH_ALL(source, with_pagination)
        target_items = FETCH_ALL(target, with_pagination)
        
        diff = CALCULATE_DIFF(source_items, target_items)
        
        FOR change IN diff:
            TRY:
                APPLY_CHANGE(change)
            CATCH conflict:
                resolution = RESOLVE_CONFLICT(conflict)
                APPLY_RESOLUTION(resolution)
        
        VALIDATE_SYNC()
    
    // Step 3: Conflict resolution
    FUNCTION resolve_conflict(conflict):
        SWITCH conflict.type:
            CASE concurrent_modification:
                RETURN merge_changes(conflict)
            CASE delete_vs_update:
                RETURN prefer_delete(conflict)
            DEFAULT:
                RETURN last_write_wins(conflict)
```

#### Validation Commands:
```bash
# Unit tests with coverage
go test ./pkg/sync/... -cover -race

# Integration test with real resources
make test-sync-integration

# Stress test with large datasets
go test ./pkg/sync/... -run=TestStress -timeout=30m

# Check for goroutine leaks
go test ./pkg/sync/... -trace=trace.out
go tool trace trace.out

# Line count check (may need split)
/home/vscode/workspaces/idpbuilder/tools/line-counter.sh -c $(git branch --show-current)
```

---

### E3.1.2: Resource Transformation
**Branch:** `/phase3/wave1/effort2-transformation`  
**Duration:** [X] hours  
**Dependencies:** Can run parallel with E3.1.1

#### Requirements:
1. **MUST** transform between formats:
   - Internal → External representation
   - Version conversion
   - Schema migration

---

## Wave 3.2: Advanced Features

### E3.2.1: Caching Layer
**Branch:** `/phase3/wave2/effort1-caching`  
**Duration:** [X] hours  
**Dependencies:** E3.1.1, E3.1.2

#### Requirements:
1. **MUST** implement:
   - LRU cache
   - TTL support
   - Invalidation strategies
   - Cache warming

#### Performance Requirements:
```yaml
performance_targets:
  cache_hit_ratio: ">80%"
  lookup_latency: "<1ms p99"
  memory_usage: "<100MB for 10k objects"
```

---

## Split Strategy (If Needed)

### When to Split:
```markdown
If effort > 800 lines:
1. Review for logical boundaries
2. Create SPLIT-PLAN.md
3. Implement splits sequentially

Example split for sync engine:
- effort1-sync-engine-part1: State machine (400 lines)
- effort1-sync-engine-part2: Sync algorithm (350 lines)  
- effort1-sync-engine-part3: Conflict resolution (300 lines)
```

### Split Validation:
```bash
# Each split must:
- Build independently
- Have its own tests
- Be under line limit
- Merge cleanly with other splits
```

## Performance Gates

### Required Benchmarks:
```go
// benchmark/sync_bench_test.go
func BenchmarkSyncEngine(b *testing.B) {
    benchmarks := []struct {
        name string
        size int
    }{
        {"small_10", 10},
        {"medium_1000", 1000},
        {"large_100000", 100000},
    }
    
    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            // Setup
            engine := setupEngine(bm.size)
            b.ResetTimer()
            
            // Benchmark
            for i := 0; i < b.N; i++ {
                engine.Sync(context.TODO())
            }
        })
    }
}
```

### Performance Criteria:
| Operation | Requirement | Measurement |
|-----------|-------------|-------------|
| Sync 1K objects | <1 second | Benchmark |
| Memory per object | <1KB | pprof |
| Concurrent syncs | 100+ | Load test |

## Integration Points

### With Phase 2:
- Uses controller framework from E2.1.1
- Uses client library from E2.2.1
- Reports metrics via E2.3.1

### With Phase 4:
- Provides sync interface for advanced features
- Exposes transformation APIs
- Enables caching for performance features

## Notes for Orchestrator

1. **Potential Splits**: E3.1.1 likely needs splitting
2. **Performance Critical**: This phase needs extensive benchmarking
3. **Integration Testing**: Requires test cluster with sample data
4. **Risk Area**: Sync engine is most complex component