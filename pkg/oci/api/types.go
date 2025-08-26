package api

import (
	"context"
	"fmt"
	"time"
)

// Constants for build configuration
const (
	DefaultMaxParallel = 4
	MaxDockerfileSize  = 10 * 1024 * 1024 // 10MB
	DefaultTimeout     = 30 * time.Minute
	MaxStages         = 50
)

// Stage represents a single build stage in a multi-stage Dockerfile.
type Stage struct {
	Name               string            `json:"name"`
	BaseImage          string            `json:"base_image"`
	Instructions       []string          `json:"instructions"`
	Dependencies       []string          `json:"dependencies"`
	BuildArgs          map[string]string `json:"build_args,omitempty"`
	Cacheable          bool              `json:"cacheable"`
	EstimatedBuildTime time.Duration     `json:"estimated_build_time"`
}

// Validate performs basic validation on a Stage.
func (s *Stage) Validate() error {
	if s.Name == "" || s.BaseImage == "" {
		return fmt.Errorf("stage name and base image required")
	}
	return nil
}

// StageAnalysis contains comprehensive analysis of a multi-stage build.
type StageAnalysis struct {
	Stages            []*Stage          `json:"stages"`
	Dependencies      map[string][]string `json:"dependencies"`
	ParallelGroups    [][]string        `json:"parallel_groups"`
	CacheableStages   []string          `json:"cacheable_stages"`
	CriticalPath      []string          `json:"critical_path"`
	EstimatedTime     time.Duration     `json:"estimated_time"`
	OptimizationScore int               `json:"optimization_score"`
	Warnings          []string          `json:"warnings,omitempty"`
}

// DependencyGraph represents the dependency relationships between stages.
type DependencyGraph struct {
	Nodes map[string]*GraphNode `json:"nodes"`
	Edges map[string][]string   `json:"edges"`
}

// GraphNode represents a node in the dependency graph.
type GraphNode struct {
	Stage        *Stage   `json:"stage"`
	Dependencies []string `json:"dependencies"`
	Level        int      `json:"level"`
}

// Validate checks the validity of a dependency graph.
func (dg *DependencyGraph) Validate() error {
	if len(dg.Nodes) == 0 {
		return fmt.Errorf("dependency graph cannot be empty")
	}
	for node, edges := range dg.Edges {
		if _, exists := dg.Nodes[node]; !exists {
			return fmt.Errorf("edge references non-existent node: %s", node)
		}
		for _, edge := range edges {
			if _, exists := dg.Nodes[edge]; !exists {
				return fmt.Errorf("edge %s references non-existent node: %s", node, edge)
			}
		}
	}
	return nil
}

// BuildRequest contains parameters for executing a multi-stage build.
type BuildRequest struct {
	Context     context.Context   `json:"-"`
	BuildArgs   map[string]string `json:"build_args,omitempty"`
	MaxParallel int              `json:"max_parallel"`
	Timeout     time.Duration    `json:"timeout"`
	CacheEnabled bool            `json:"cache_enabled"`
	Tags        []string         `json:"tags,omitempty"`
}

// BuildResources defines the resources available for building.
type BuildResources struct {
	MaxConcurrentBuilds int           `json:"max_concurrent_builds"`
	MemoryLimit        int64         `json:"memory_limit"`
	CPULimit           int           `json:"cpu_limit"`
	Timeout            time.Duration `json:"timeout"`
}

// StageResult contains the results of executing build stages.
type StageResult struct {
	Success       bool                    `json:"success"`
	StageResults  map[string]*StageExecution `json:"stage_results"`
	BuildTime     time.Duration           `json:"build_time"`
	ParallelJobs  int                     `json:"parallel_jobs"`
	CacheHits     int                     `json:"cache_hits"`
	Warnings      []string                `json:"warnings,omitempty"`
	Error         string                  `json:"error,omitempty"`
}

// StageExecution contains the execution details for a single stage.
type StageExecution struct {
	Name          string        `json:"name"`
	Success       bool          `json:"success"`
	BuildTime     time.Duration `json:"build_time"`
	CacheHit      bool          `json:"cache_hit"`
	ParallelGroup int           `json:"parallel_group"`
	Error         string        `json:"error,omitempty"`
}

// BuildMetrics tracks performance metrics for the optimizer.
type BuildMetrics struct {
	TotalBuilds          int           `json:"total_builds"`
	AverageBuildTime     time.Duration `json:"average_build_time"`
	ParallelizationRatio float64       `json:"parallelization_ratio"`
	CacheHitRate         float64       `json:"cache_hit_rate"`
	ResourceUtilization  float64       `json:"resource_utilization"`
	LastUpdated          time.Time     `json:"last_updated"`
}

// StageOptimizer defines the interface for multi-stage build optimization.
type StageOptimizer interface {
	AnalyzeStages(dockerfile []byte) (*StageAnalysis, error)
	BuildStages(ctx context.Context, analysis *StageAnalysis, req *BuildRequest) (*StageResult, error)
	OptimizeDependencies(stages []*Stage) (*DependencyGraph, error)
	EstimateBuildTime(analysis *StageAnalysis, resources *BuildResources) (time.Duration, error)
	ValidateOptimization(original []*Stage, optimized *DependencyGraph) error
	GetMetrics() *BuildMetrics
}