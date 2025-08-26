// Package optimizer provides multi-stage build optimization and parallel execution capabilities.
package optimizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// Optimizer implements the api.StageOptimizer interface for multi-stage Docker build optimization.
// It provides comprehensive stage analysis, parallel execution, and dependency optimization.
type Optimizer struct {
	analyzer   *Analyzer
	executor   *Executor  
	graph      *GraphBuilder
	metrics    *api.BuildMetrics
	mu         sync.RWMutex
	maxWorkers int
}

// NewOptimizer creates a new multi-stage build optimizer with default configuration.
func NewOptimizer() api.StageOptimizer {
	return &Optimizer{
		analyzer:   NewAnalyzer(),
		executor:   NewExecutor(api.DefaultMaxParallel),
		graph:      NewGraphBuilder(),
		metrics:    &api.BuildMetrics{
			LastUpdated: time.Now(),
		},
		maxWorkers: api.DefaultMaxParallel,
	}
}

// NewOptimizerWithConfig creates a new optimizer with custom configuration.
func NewOptimizerWithConfig(maxWorkers int) api.StageOptimizer {
	if maxWorkers <= 0 {
		maxWorkers = api.DefaultMaxParallel
	}

	return &Optimizer{
		analyzer:   NewAnalyzer(),
		executor:   NewExecutor(maxWorkers),
		graph:      NewGraphBuilder(),
		metrics:    &api.BuildMetrics{
			LastUpdated: time.Now(),
		},
		maxWorkers: maxWorkers,
	}
}

// AnalyzeStages examines a Dockerfile and identifies optimization opportunities.
// It parses the multi-stage structure, analyzes dependencies, and identifies parallel execution groups.
func (o *Optimizer) AnalyzeStages(dockerfile []byte) (*api.StageAnalysis, error) {
	if len(dockerfile) == 0 {
		return nil, fmt.Errorf("dockerfile cannot be empty")
	}

	if len(dockerfile) > api.MaxDockerfileSize {
		return nil, fmt.Errorf("dockerfile size exceeds maximum limit of %d bytes", api.MaxDockerfileSize)
	}

	// Use the analyzer to parse and analyze the Dockerfile
	analysis, err := o.analyzer.Analyze(dockerfile)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze dockerfile: %w", err)
	}

	// Validate the analysis results
	if err := o.validateAnalysis(analysis); err != nil {
		return nil, fmt.Errorf("invalid analysis result: %w", err)
	}

	o.mu.Lock()
	o.metrics.LastUpdated = time.Now()
	o.mu.Unlock()

	return analysis, nil
}

// BuildStages executes build stages with optimal parallelization.
// It uses the stage analysis to determine execution order and maximize parallel execution.
func (o *Optimizer) BuildStages(ctx context.Context, analysis *api.StageAnalysis, req *api.BuildRequest) (*api.StageResult, error) {
	if analysis == nil {
		return nil, fmt.Errorf("stage analysis cannot be nil")
	}

	if req == nil {
		return nil, fmt.Errorf("build request cannot be nil")
	}

	// Validate the build request
	if err := o.validateBuildRequest(req); err != nil {
		return nil, fmt.Errorf("invalid build request: %w", err)
	}

	startTime := time.Now()

	// Create dependency graph for execution planning
	graph, err := o.OptimizeDependencies(analysis.Stages)
	if err != nil {
		return nil, fmt.Errorf("failed to create dependency graph: %w", err)
	}

	// Execute stages using the parallel executor
	result, err := o.executor.Execute(ctx, graph, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute stages: %w", err)
	}

	// Update metrics with build results
	buildTime := time.Since(startTime)
	o.updateMetrics(buildTime, len(analysis.Stages), len(analysis.ParallelGroups))

	return result, nil
}

// OptimizeDependencies analyzes stage dependencies and creates an optimized dependency graph.
// It performs topological sorting and identifies parallel execution opportunities.
func (o *Optimizer) OptimizeDependencies(stages []*api.Stage) (*api.DependencyGraph, error) {
	if len(stages) == 0 {
		return nil, fmt.Errorf("stages list cannot be empty")
	}

	// Validate all stages
	for i, stage := range stages {
		if stage == nil {
			return nil, fmt.Errorf("stage %d cannot be nil", i)
		}
		if err := stage.Validate(); err != nil {
			return nil, fmt.Errorf("stage %d validation failed: %w", i, err)
		}
	}

	// Use graph builder to create optimized dependency graph
	graph, err := o.graph.BuildGraph(stages)
	if err != nil {
		return nil, fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Validate the resulting graph
	if err := graph.Validate(); err != nil {
		return nil, fmt.Errorf("invalid dependency graph: %w", err)
	}

	return graph, nil
}

// GetMetrics returns current build metrics and performance statistics.
func (o *Optimizer) GetMetrics() *api.BuildMetrics {
	o.mu.RLock()
	defer o.mu.RUnlock()

	// Return a copy to prevent external modification
	return &api.BuildMetrics{
		TotalBuilds:          o.metrics.TotalBuilds,
		AverageBuildTime:     o.metrics.AverageBuildTime,
		ParallelizationRatio: o.metrics.ParallelizationRatio,
		CacheHitRate:         o.metrics.CacheHitRate,
		ResourceUtilization:  o.metrics.ResourceUtilization,
		LastUpdated:          o.metrics.LastUpdated,
	}
}

// EstimateBuildTime predicts the total build time based on stage analysis.
// It considers parallelization opportunities and resource constraints.
func (o *Optimizer) EstimateBuildTime(analysis *api.StageAnalysis, resources *api.BuildResources) (time.Duration, error) {
	if analysis == nil {
		return 0, fmt.Errorf("stage analysis cannot be nil")
	}

	if resources == nil {
		return 0, fmt.Errorf("build resources cannot be nil")
	}

	if len(analysis.Stages) == 0 {
		return 0, fmt.Errorf("no stages to estimate")
	}

	// Calculate sequential build time
	var sequentialTime time.Duration
	for _, stage := range analysis.Stages {
		sequentialTime += stage.EstimatedBuildTime
	}

	// If no parallelization is possible, return sequential time
	if len(analysis.ParallelGroups) <= 1 {
		return sequentialTime, nil
	}

	// Calculate parallel execution time based on critical path
	criticalPathTime := o.calculateCriticalPathTime(analysis)

	// Apply resource constraints and parallelization efficiency
	maxParallel := resources.MaxConcurrentBuilds
	if maxParallel <= 0 {
		maxParallel = o.maxWorkers
	}

	parallelizationFactor := o.calculateParallelizationFactor(analysis, maxParallel)
	estimatedTime := time.Duration(float64(criticalPathTime) / parallelizationFactor)

	// Add overhead for coordination and synchronization (5-10%)
	overhead := time.Duration(float64(estimatedTime) * 0.075)
	totalEstimatedTime := estimatedTime + overhead

	return totalEstimatedTime, nil
}

// ValidateOptimization verifies that the proposed optimization is safe.
// It ensures that the dependency graph maintains build correctness.
func (o *Optimizer) ValidateOptimization(original []*api.Stage, optimized *api.DependencyGraph) error {
	if len(original) == 0 {
		return fmt.Errorf("original stages list cannot be empty")
	}

	if optimized == nil {
		return fmt.Errorf("optimized dependency graph cannot be nil")
	}

	// Validate the optimized graph
	if err := optimized.Validate(); err != nil {
		return fmt.Errorf("optimized graph is invalid: %w", err)
	}

	// Ensure all original stages are present in the optimized graph
	if len(original) != len(optimized.Nodes) {
		return fmt.Errorf("stage count mismatch: original %d, optimized %d", len(original), len(optimized.Nodes))
	}

	// Check that all original stages exist in the optimized graph
	for _, stage := range original {
		if _, exists := optimized.Nodes[stage.Name]; !exists {
			return fmt.Errorf("stage %s missing from optimized graph", stage.Name)
		}
	}

	// Validate dependency preservation
	return o.validateDependencyPreservation(original, optimized)
}

// Helper methods

func (o *Optimizer) validateAnalysis(analysis *api.StageAnalysis) error {
	if analysis == nil {
		return fmt.Errorf("analysis cannot be nil")
	}

	if len(analysis.Stages) == 0 {
		return fmt.Errorf("analysis must contain at least one stage")
	}

	// Validate each stage
	for i, stage := range analysis.Stages {
		if err := stage.Validate(); err != nil {
			return fmt.Errorf("stage %d validation failed: %w", i, err)
		}
	}

	return nil
}

func (o *Optimizer) validateBuildRequest(req *api.BuildRequest) error {
	if req.MaxParallel < 0 {
		return fmt.Errorf("max parallel cannot be negative")
	}

	if req.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}

	return nil
}

func (o *Optimizer) updateMetrics(buildTime time.Duration, stageCount, parallelGroups int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.metrics.TotalBuilds++
	
	// Update average build time
	if o.metrics.TotalBuilds == 1 {
		o.metrics.AverageBuildTime = buildTime
	} else {
		total := o.metrics.AverageBuildTime * time.Duration(o.metrics.TotalBuilds-1)
		o.metrics.AverageBuildTime = (total + buildTime) / time.Duration(o.metrics.TotalBuilds)
	}

	// Update parallelization ratio
	if stageCount > 0 {
		o.metrics.ParallelizationRatio = float64(parallelGroups) / float64(stageCount)
	}

	o.metrics.LastUpdated = time.Now()
}

func (o *Optimizer) calculateCriticalPathTime(analysis *api.StageAnalysis) time.Duration {
	var maxTime time.Duration
	
	// Find the longest path through the critical path stages
	for _, stageName := range analysis.CriticalPath {
		for _, stage := range analysis.Stages {
			if stage.Name == stageName && stage.EstimatedBuildTime > maxTime {
				maxTime = stage.EstimatedBuildTime
			}
		}
	}

	return maxTime
}

func (o *Optimizer) calculateParallelizationFactor(analysis *api.StageAnalysis, maxParallel int) float64 {
	if len(analysis.ParallelGroups) == 0 {
		return 1.0
	}

	// Calculate ideal parallelization factor
	totalStages := len(analysis.Stages)
	maxGroupSize := 0
	for _, group := range analysis.ParallelGroups {
		if len(group) > maxGroupSize {
			maxGroupSize = len(group)
		}
	}

	idealFactor := float64(totalStages) / float64(len(analysis.ParallelGroups))
	resourceConstrainedFactor := float64(maxParallel)

	// Use the minimum of ideal and resource-constrained factors
	if idealFactor < resourceConstrainedFactor {
		return idealFactor
	}
	return resourceConstrainedFactor
}

func (o *Optimizer) validateDependencyPreservation(original []*api.Stage, optimized *api.DependencyGraph) error {
	// Create a map of original dependencies for easy lookup
	originalDeps := make(map[string][]string)
	for _, stage := range original {
		originalDeps[stage.Name] = stage.Dependencies
	}

	// Check that all dependencies are preserved in the optimized graph
	for stageName, deps := range originalDeps {
		optimizedDeps, exists := optimized.Edges[stageName]
		if !exists && len(deps) > 0 {
			return fmt.Errorf("stage %s dependencies not preserved in optimization", stageName)
		}

		// Verify all original dependencies are present
		for _, dep := range deps {
			found := false
			for _, optimizedDep := range optimizedDeps {
				if dep == optimizedDep {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("dependency %s -> %s not preserved in optimization", stageName, dep)
			}
		}
	}

	return nil
}