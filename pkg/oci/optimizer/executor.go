package optimizer

import (
	"context"
	"fmt"
	"github.com/jessesanford/idpbuilder/pkg/oci/api"
	"sync"
	"time"
)

type Executor struct {
	workers int
	pool    chan struct{}
}

func NewExecutor(workers int) *Executor {
	if workers <= 0 {
		workers = api.DefaultMaxParallel
	}
	return &Executor{workers: workers, pool: make(chan struct{}, workers)}
}

func (e *Executor) Execute(ctx context.Context, stages []*api.Stage) (*api.StageResult, error) {
	if len(stages) == 0 {
		return &api.StageResult{Success: true, StageResults: make(map[string]*api.StageExecution)}, nil
	}

	levels, err := e.scheduleStages(stages)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule stages: %w", err)
	}

	result := &api.StageResult{StageResults: make(map[string]*api.StageExecution), ParallelJobs: e.workers}
	start := time.Now()

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

func (e *Executor) executeLevel(ctx context.Context, stages []*api.Stage, result *api.StageResult, level int) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(stages))

	for _, stage := range stages {
		wg.Add(1)
		go e.executeStage(ctx, stage, &wg, result, level, errCh)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) executeStage(ctx context.Context, stage *api.Stage, wg *sync.WaitGroup, result *api.StageResult, level int, errCh chan<- error) {
	defer wg.Done()
	e.pool <- struct{}{}
	defer func() { <-e.pool }()

	start := time.Now()
	execution := &api.StageExecution{Name: stage.Name, ParallelGroup: level}

	select {
	case <-ctx.Done():
		execution.Success = false
		execution.Error = ctx.Err().Error()
		errCh <- ctx.Err()
		return
	default:
	}

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
	execution.CacheHit = stage.Cacheable

	mu := &sync.Mutex{}
	mu.Lock()
	result.StageResults[stage.Name] = execution
	if execution.CacheHit {
		result.CacheHits++
	}
	mu.Unlock()

	errCh <- nil
}

func (e *Executor) simulateStageExecution(ctx context.Context, stage *api.Stage) error {
	if err := stage.Validate(); err != nil {
		return err
	}
	buildTime := stage.EstimatedBuildTime
	if buildTime == 0 || buildTime > time.Minute {
		buildTime = 5 * time.Second
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(buildTime):
		return nil
	}
}

func (e *Executor) scheduleStages(stages []*api.Stage) ([][]*api.Stage, error) {
	if len(stages) == 0 {
		return nil, fmt.Errorf("no stages to schedule")
	}
	levels := make([][]*api.Stage, 0)
	processed := make(map[string]bool)
	inDegree := make(map[string]int)

	for _, stage := range stages {
		inDegree[stage.Name] = len(stage.Dependencies)
	}

	for len(processed) < len(stages) {
		currentLevel := make([]*api.Stage, 0)

		for _, stage := range stages {
			if !processed[stage.Name] && inDegree[stage.Name] == 0 {
				currentLevel = append(currentLevel, stage)
			}
		}

		if len(currentLevel) == 0 {
			return nil, fmt.Errorf("circular dependency detected")
		}

		levels = append(levels, currentLevel)

		for _, stage := range currentLevel {
			processed[stage.Name] = true
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
