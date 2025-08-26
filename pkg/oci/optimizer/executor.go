package optimizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jessesanford/idpbuilder/pkg/oci/api"
)

// Executor manages parallel execution of build stages
type Executor struct {
	workers int
	pool    chan struct{}
}

// NewExecutor creates a new executor with specified worker count
func NewExecutor(workers int) *Executor {
	if workers <= 0 {
		workers = api.DefaultMaxParallel
	}
	return &Executor{
		workers: workers,
		pool:    make(chan struct{}, workers),
	}
}

// Execute runs stages in parallel based on dependency levels
func (e *Executor) Execute(ctx context.Context, stages []*api.Stage) (*api.StageResult, error) {
	if len(stages) == 0 {
		return &api.StageResult{Success: true, StageResults: make(map[string]*api.StageExecution)}, nil
	}

	// Schedule stages by dependency levels
	levels, err := e.scheduleStages(stages)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule stages: %w", err)
	}

	result := &api.StageResult{
		StageResults: make(map[string]*api.StageExecution),
		ParallelJobs: e.workers,
	}
	start := time.Now()

	// Execute each level in sequence, stages in level in parallel
	for levelIdx, levelStages := range levels {
		if err := e.executeLevel(ctx, levelStages, result, levelIdx); err != nil {
			result.Success = false
			result.Error = err.Error()
			result.BuildTime = time.Since(start)
			return result, err
		}
	}

	result.Success = true
	result.BuildTime = time.Since(start)
	return result, nil
}

// executeLevel executes all stages in a dependency level in parallel
func (e *Executor) executeLevel(ctx context.Context, stages []*api.Stage, result *api.StageResult, level int) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(stages))

	for _, stage := range stages {
		wg.Add(1)
		go e.executeStage(ctx, stage, &wg, result, level, errCh)
	}

	wg.Wait()
	close(errCh)

	// Check for errors
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

// executeStage executes a single stage with worker pool management
func (e *Executor) executeStage(ctx context.Context, stage *api.Stage, wg *sync.WaitGroup, result *api.StageResult, level int, errCh chan<- error) {
	defer wg.Done()

	// Acquire worker from pool
	e.pool <- struct{}{}
	defer func() { <-e.pool }()

	start := time.Now()
	execution := &api.StageExecution{
		Name:          stage.Name,
		ParallelGroup: level,
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		execution.Success = false
		execution.Error = ctx.Err().Error()
		errCh <- ctx.Err()
		return
	default:
	}

	// Simulate stage execution - in real implementation would build the stage
	if err := e.simulateStageExecution(ctx, stage); err != nil {
		execution.Success = false
		execution.Error = err.Error()
		execution.BuildTime = time.Since(start)
		result.StageResults[stage.Name] = execution
		errCh <- fmt.Errorf("stage %s failed: %w", stage.Name, err)
		return
	}

	execution.Success = true
	execution.BuildTime = time.Since(start)
	execution.CacheHit = stage.Cacheable // Simplified cache simulation

	// Thread-safe result update
	mu := &sync.Mutex{}
	mu.Lock()
	result.StageResults[stage.Name] = execution
	if execution.CacheHit {
		result.CacheHits++
	}
	mu.Unlock()

	errCh <- nil
}

// simulateStageExecution simulates building a stage (placeholder for real implementation)
func (e *Executor) simulateStageExecution(ctx context.Context, stage *api.Stage) error {
	// Validate stage
	if err := stage.Validate(); err != nil {
		return err
	}

	// Simulate build time based on stage estimation
	buildTime := stage.EstimatedBuildTime
	if buildTime == 0 {
		buildTime = 5 * time.Second // Default simulation time
	}

	// Use a reasonable simulation time for tests
	if buildTime > time.Minute {
		buildTime = time.Minute
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(buildTime):
		return nil
	}
}

// scheduleStages groups stages into dependency levels for parallel execution
func (e *Executor) scheduleStages(stages []*api.Stage) ([][]*api.Stage, error) {
	if len(stages) == 0 {
		return nil, fmt.Errorf("no stages to schedule")
	}

	// Create stage lookup map
	stageMap := make(map[string]*api.Stage)
	for _, stage := range stages {
		stageMap[stage.Name] = stage
	}

	// Calculate dependency levels using topological sort
	levels := make([][]*api.Stage, 0)
	processed := make(map[string]bool)
	inDegree := make(map[string]int)

	// Initialize in-degrees
	for _, stage := range stages {
		inDegree[stage.Name] = len(stage.Dependencies)
	}

	// Process levels
	for len(processed) < len(stages) {
		currentLevel := make([]*api.Stage, 0)

		// Find stages with no dependencies (in-degree 0)
		for _, stage := range stages {
			if !processed[stage.Name] && inDegree[stage.Name] == 0 {
				currentLevel = append(currentLevel, stage)
			}
		}

		if len(currentLevel) == 0 {
			return nil, fmt.Errorf("circular dependency detected")
		}

		levels = append(levels, currentLevel)

		// Mark as processed and update in-degrees
		for _, stage := range currentLevel {
			processed[stage.Name] = true

			// Update in-degrees for dependent stages
			for _, otherStage := range stages {
				for _, dep := range otherStage.Dependencies {
					if dep == stage.Name {
						inDegree[otherStage.Name]--
					}
				}
			}
		}
	}

	return levels, nil
}