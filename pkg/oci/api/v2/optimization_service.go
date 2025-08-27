package v2

import (
	"context"
	"time"
)

// OptimizationService defines the contract for performance optimization
type OptimizationService interface {
	// AnalyzeBuild analyzes a build for optimization opportunities
	AnalyzeBuild(ctx context.Context, dockerfile string) (*BuildAnalysis, error)

	// OptimizeBuildOrder optimizes build step ordering
	OptimizeBuildOrder(steps []BuildStep) []BuildStep

	// PredictCacheHit predicts cache hit probability
	PredictCacheHit(layer LayerInfo) float64

	// GetCacheStatistics returns cache performance statistics
	GetCacheStatistics() CacheStats

	// EnablePersistentCache enables persistent caching
	EnablePersistentCache(config CacheConfig) error

	// ClearCache clears optimization caches
	ClearCache(cacheType CacheType) error

	// GetPerformanceMetrics returns performance metrics
	GetPerformanceMetrics() PerformanceMetrics
}

// BuildAnalysis represents build optimization analysis
type BuildAnalysis struct {
	OptimizationOpportunities []Optimization
	EstimatedTimeSaved        time.Duration
	CacheableStages           []string
	ParallelizableStages      [][]string
}

// Optimization represents an optimization opportunity
type Optimization struct {
	Type        OptimizationType
	Description string
	Impact      ImpactLevel
	Stage       string
}

// OptimizationType defines types of optimizations
type OptimizationType string

const (
	OptimizationTypeLayerReorder    OptimizationType = "layer-reorder"
	OptimizationTypeParallelization OptimizationType = "parallelization"
	OptimizationTypeCacheReuse      OptimizationType = "cache-reuse"
	OptimizationTypeMultiStage      OptimizationType = "multi-stage"
)

// ImpactLevel represents optimization impact
type ImpactLevel string

const (
	ImpactLevelHigh   ImpactLevel = "high"
	ImpactLevelMedium ImpactLevel = "medium"
	ImpactLevelLow    ImpactLevel = "low"
)

// BuildStep represents a build step
type BuildStep struct {
	Command       string
	Dependencies  []string
	Cacheable     bool
	EstimatedTime time.Duration
}

// LayerInfo represents layer information
type LayerInfo struct {
	Digest  string
	Size    int64
	Created time.Time
	Command string
}

// CacheConfig defines cache configuration
type CacheConfig struct {
	PersistentPath     string
	MaxSize            int64
	TTL                time.Duration
	CompressionEnabled bool
}

// CacheType defines cache types
type CacheType string

const (
	CacheTypeLayer    CacheType = "layer"
	CacheTypeBuild    CacheType = "build"
	CacheTypeMetadata CacheType = "metadata"
)

// CacheStats provides cache statistics
type CacheStats struct {
	HitRate       float64
	MissRate      float64
	TotalHits     int64
	TotalMisses   int64
	CacheSize     int64
	EvictionCount int64
}

// PerformanceMetrics provides performance metrics
type PerformanceMetrics struct {
	AverageBuildTime    time.Duration
	CacheHitRate        float64
	OptimizationSavings time.Duration
	ResourceUsage       ResourceMetrics
}

// ResourceMetrics tracks resource usage
type ResourceMetrics struct {
	CPUUsage    float64
	MemoryUsage int64
	DiskIO      int64
}