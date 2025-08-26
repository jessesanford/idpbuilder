package api

import (
	"context"
	"testing"
	"time"
)

// MockStageOptimizer implements StageOptimizer for testing
type MockStageOptimizer struct {
	AnalyzeStagesFunc      func([]byte) (*StageAnalysis, error)
	BuildStagesFunc        func(context.Context, *StageAnalysis, *BuildRequest) (*StageResult, error)
	OptimizeDependenciesFunc func([]*Stage) (*DependencyGraph, error)
	GetMetricsFunc         func() *BuildMetrics
	EstimateBuildTimeFunc  func(*StageAnalysis, *BuildResources) (time.Duration, error)
	ValidateOptimizationFunc func([]*Stage, *DependencyGraph) error
}

func (m *MockStageOptimizer) AnalyzeStages(dockerfile []byte) (*StageAnalysis, error) {
	if m.AnalyzeStagesFunc != nil {
		return m.AnalyzeStagesFunc(dockerfile)
	}
	return &StageAnalysis{}, nil
}

func (m *MockStageOptimizer) BuildStages(ctx context.Context, analysis *StageAnalysis, req *BuildRequest) (*StageResult, error) {
	if m.BuildStagesFunc != nil {
		return m.BuildStagesFunc(ctx, analysis, req)
	}
	return &StageResult{}, nil
}

func (m *MockStageOptimizer) OptimizeDependencies(stages []*Stage) (*DependencyGraph, error) {
	if m.OptimizeDependenciesFunc != nil {
		return m.OptimizeDependenciesFunc(stages)
	}
	return &DependencyGraph{}, nil
}

func (m *MockStageOptimizer) GetMetrics() *BuildMetrics {
	if m.GetMetricsFunc != nil {
		return m.GetMetricsFunc()
	}
	return &BuildMetrics{}
}

func (m *MockStageOptimizer) EstimateBuildTime(analysis *StageAnalysis, resources *BuildResources) (time.Duration, error) {
	if m.EstimateBuildTimeFunc != nil {
		return m.EstimateBuildTimeFunc(analysis, resources)
	}
	return time.Minute, nil
}

func (m *MockStageOptimizer) ValidateOptimization(original []*Stage, optimized *DependencyGraph) error {
	if m.ValidateOptimizationFunc != nil {
		return m.ValidateOptimizationFunc(original, optimized)
	}
	return nil
}

// Test interface compliance
func TestStageOptimizerInterface(t *testing.T) {
	var _ StageOptimizer = &MockStageOptimizer{}
}

func TestStageAnalysis(t *testing.T) {
	analysis := &StageAnalysis{
		Stages: []*Stage{
			{Name: "base", BaseImage: "alpine", Instructions: []string{"RUN apk add curl"}},
		},
		Dependencies:      map[string][]string{"base": {}},
		ParallelGroups:    [][]string{{"base"}},
		CacheableStages:   []string{"base"},
		EstimatedTime:     time.Minute,
		CriticalPath:      []string{"base"},
		OptimizationScore: 85,
	}

	if len(analysis.Stages) != 1 {
		t.Errorf("Expected 1 stage, got %d", len(analysis.Stages))
	}
	if analysis.OptimizationScore != 85 {
		t.Errorf("Expected optimization score 85, got %d", analysis.OptimizationScore)
	}
}

func TestBuildRequest(t *testing.T) {
	req := &BuildRequest{
		BuildArgs:   map[string]string{"VERSION": "1.0"},
		Tags:        []string{"myapp:latest", "myapp:1.0"},
		Target:      "production",
		Platform:    "linux/amd64",
		NoCache:     false,
		Pull:        true,
		MaxParallel: 4,
		Timeout:     30 * time.Minute,
	}

	if req.Target != "production" {
		t.Errorf("Expected target 'production', got '%s'", req.Target)
	}
	if req.MaxParallel != 4 {
		t.Errorf("Expected max parallel 4, got %d", req.MaxParallel)
	}
}