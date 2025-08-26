// Package api provides interfaces and contracts for advanced OCI build capabilities
package api

import (
	"context"
	"time"
)

// StageOptimizer handles multi-stage build optimization and parallel execution.
type StageOptimizer interface {
	// AnalyzeStages examines a Dockerfile and identifies optimization opportunities.
	AnalyzeStages(dockerfile []byte) (*StageAnalysis, error)

	// BuildStages executes build stages with optimal parallelization.
	BuildStages(ctx context.Context, analysis *StageAnalysis, req *BuildRequest) (*StageResult, error)

	// OptimizeDependencies analyzes stage dependencies and creates an optimized dependency graph.
	OptimizeDependencies(stages []*Stage) (*DependencyGraph, error)

	// GetMetrics returns current build metrics and performance statistics.
	GetMetrics() *BuildMetrics

	// EstimateBuildTime predicts the total build time based on stage analysis.
	EstimateBuildTime(analysis *StageAnalysis, resources *BuildResources) (time.Duration, error)

	// ValidateOptimization verifies that the proposed optimization is safe.
	ValidateOptimization(original []*Stage, optimized *DependencyGraph) error
}

// StageAnalysis contains the results of analyzing a multi-stage Dockerfile.
type StageAnalysis struct {
	Stages              []*Stage              `json:"stages" yaml:"stages"`
	Dependencies        map[string][]string   `json:"dependencies" yaml:"dependencies"`
	ParallelGroups      [][]string            `json:"parallel_groups" yaml:"parallel_groups"`
	CacheableStages     []string              `json:"cacheable_stages" yaml:"cacheable_stages"`
	EstimatedTime       time.Duration         `json:"estimated_time" yaml:"estimated_time"`
	CriticalPath        []string              `json:"critical_path" yaml:"critical_path"`
	OptimizationScore   int                   `json:"optimization_score" yaml:"optimization_score"`
	Warnings            []string              `json:"warnings,omitempty" yaml:"warnings,omitempty"`
}

// BuildRequest contains parameters for a stage build operation.
type BuildRequest struct {
	BuildArgs   map[string]string `json:"build_args,omitempty" yaml:"build_args,omitempty"`
	Tags        []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Target      string            `json:"target,omitempty" yaml:"target,omitempty"`
	Platform    string            `json:"platform,omitempty" yaml:"platform,omitempty"`
	NoCache     bool              `json:"no_cache" yaml:"no_cache"`
	Pull        bool              `json:"pull" yaml:"pull"`
	MaxParallel int               `json:"max_parallel,omitempty" yaml:"max_parallel,omitempty"`
	Timeout     time.Duration     `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

// BuildMetrics contains performance metrics for build optimization analysis.
type BuildMetrics struct {
	TotalBuilds          int64          `json:"total_builds" yaml:"total_builds"`
	AverageBuildTime     time.Duration  `json:"average_build_time" yaml:"average_build_time"`
	ParallelizationRatio float64        `json:"parallelization_ratio" yaml:"parallelization_ratio"`
	CacheHitRate         float64        `json:"cache_hit_rate" yaml:"cache_hit_rate"`
	ResourceUtilization  *ResourceStats `json:"resource_utilization,omitempty" yaml:"resource_utilization,omitempty"`
	LastUpdated          time.Time      `json:"last_updated" yaml:"last_updated"`
}

// BuildResources defines available resources for build optimization.
type BuildResources struct {
	CPUs                int   `json:"cpus" yaml:"cpus"`
	Memory              int64 `json:"memory" yaml:"memory"`
	Storage             int64 `json:"storage" yaml:"storage"`
	MaxConcurrentBuilds int   `json:"max_concurrent_builds" yaml:"max_concurrent_builds"`
}

// ResourceStats contains resource utilization statistics.
type ResourceStats struct {
	CPUUsage    float64 `json:"cpu_usage" yaml:"cpu_usage"`
	MemoryUsage int64   `json:"memory_usage" yaml:"memory_usage"`
	DiskUsage   int64   `json:"disk_usage" yaml:"disk_usage"`
	NetworkIO   int64   `json:"network_io" yaml:"network_io"`
}