package builder

import (
	"context"
	"fmt"
	"sync"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

// AdvancedBuilder extends the basic Builder with advanced functionality.
type AdvancedBuilder struct {
	*Builder
	multiStageSupport bool
	parallelLayers    bool
	buildCache        BuildCache
	metrics           *BuildMetrics
	mutex             sync.RWMutex
}

// BuildCache provides advanced caching functionality.
type BuildCache interface {
	LayerCache
	GetBuildResult(key string) (*BuildResult, bool)
	PutBuildResult(key string, result *BuildResult)
	GetLayerByContentHash(hash string) (v1.Layer, bool)
	PutLayerWithContentHash(hash string, layer v1.Layer)
	InvalidatePattern(pattern string) int
}

// BuildMetrics tracks build performance and statistics.
type BuildMetrics struct {
	StartTime        time.Time
	EndTime          time.Time
	TotalDuration    time.Duration
	LayerCount       int
	LayerSizes       []int64
	CacheHits        int
	CacheMisses      int
	RegistryPulls    int
	RegistryPushes   int
	ParallelJobs     int
	PeakMemoryUsage  int64
	NetworkTransfer  int64
	CompressionRatio float64
	mutex            sync.RWMutex
}

// NewAdvancedBuilder creates a new AdvancedBuilder with enhanced features.
func NewAdvancedBuilder(config *BuildConfig, registry RegistryClient, cache BuildCache) (*AdvancedBuilder, error) {
	baseBuilder, err := NewBuilder(config, registry, cache)
	if err != nil {
		return nil, fmt.Errorf("failed to create base builder: %w", err)
	}

	return &AdvancedBuilder{
		Builder:           baseBuilder,
		multiStageSupport: true,
		parallelLayers:    true,
		buildCache:        cache,
		metrics:           &BuildMetrics{},
	}, nil
}

// BuildWithOptimizations builds an image with advanced optimizations.
func (ab *AdvancedBuilder) BuildWithOptimizations(ctx context.Context, opts *BuildOptimizations) (*BuildResult, error) {
	ab.mutex.Lock()
	defer ab.mutex.Unlock()

	ab.metrics.StartTime = time.Now()
	defer func() {
		ab.metrics.EndTime = time.Now()
		ab.metrics.TotalDuration = ab.metrics.EndTime.Sub(ab.metrics.StartTime)
	}()

	// Check build cache first
	if opts.EnableBuildCache {
		cacheKey := ab.generateCacheKey()
		if result, found := ab.buildCache.GetBuildResult(cacheKey); found {
			ab.metrics.CacheHits++
			return result, nil
		}
		ab.metrics.CacheMisses++
	}

	// Get base image with optimization
	baseImage, err := ab.getOptimizedBaseImage(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get optimized base image: %w", err)
	}

	// Set parallel jobs metric
	maxJobs := opts.MaxParallelJobs
	if maxJobs <= 0 {
		maxJobs = 4
	}
	ab.metrics.ParallelJobs = maxJobs

	// Apply layers with parallelization if enabled
	var finalImage v1.Image = baseImage
	if ab.parallelLayers && opts.ParallelLayerProcessing {
		finalImage, err = ab.applyLayersParallel(ctx, baseImage, opts)
	} else {
		finalImage, err = ab.applyLayersSequential(ctx, baseImage, opts)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to apply layers: %w", err)
	}

	// Create build result with metrics
	result, err := ab.createBuildResult(finalImage)
	if err != nil {
		return nil, fmt.Errorf("failed to create build result: %w", err)
	}

	// Cache the result if enabled
	if opts.EnableBuildCache {
		cacheKey := ab.generateCacheKey()
		ab.buildCache.PutBuildResult(cacheKey, result)
	}

	return result, nil
}

// BuildOptimizations contains options for optimized building.
type BuildOptimizations struct {
	EnableBuildCache        bool
	ParallelLayerProcessing bool
	LayerCompression        CompressionType
	LayerDeduplication      bool
	ImageSquashing          bool
	RegistryOptimizations   bool
	MaxParallelJobs         int
	MemoryLimit             int64
	TimeLimit               time.Duration
}

// MultiStageBuild handles multi-stage Dockerfile builds.
func (ab *AdvancedBuilder) MultiStageBuild(ctx context.Context, stages []BuildStage) (*BuildResult, error) {
	if !ab.multiStageSupport {
		return nil, fmt.Errorf("multi-stage builds not enabled")
	}

	stageResults := make(map[string]*BuildResult)
	
	for i, stage := range stages {
		stageCtx := context.WithValue(ctx, "stage", stage.Name)
		
		// Build stage with dependencies
		stageBuilder, err := ab.createStageBuilder(stage, stageResults)
		if err != nil {
			return nil, fmt.Errorf("failed to create stage builder for %s: %w", stage.Name, err)
		}

		result, err := stageBuilder.Build(stageCtx)
		if err != nil {
			return nil, fmt.Errorf("failed to build stage %s: %w", stage.Name, err)
		}

		stageResults[stage.Name] = result
		ab.metrics.LayerCount += len(stage.Layers)

		// Only keep the final stage result
		if i == len(stages)-1 {
			return result, nil
		}
	}

	return nil, fmt.Errorf("no final stage found")
}

// BuildStage represents a single stage in a multi-stage build.
type BuildStage struct {
	Name     string
	BaseRef  string
	Layers   []v1.Layer
	Config   *v1.Config
	CopyFrom []StageCopy
}

// StageCopy represents a COPY --from instruction.
type StageCopy struct {
	FromStage string
	SrcPath   string
	DestPath  string
}

// GetMetrics returns the current build metrics.
func (ab *AdvancedBuilder) GetMetrics() *BuildMetrics {
	ab.metrics.mutex.RLock()
	defer ab.metrics.mutex.RUnlock()
	
	// Return a copy to prevent race conditions
	return &BuildMetrics{
		StartTime:        ab.metrics.StartTime,
		EndTime:          ab.metrics.EndTime,
		TotalDuration:    ab.metrics.TotalDuration,
		LayerCount:       ab.metrics.LayerCount,
		LayerSizes:       append([]int64{}, ab.metrics.LayerSizes...),
		CacheHits:        ab.metrics.CacheHits,
		CacheMisses:      ab.metrics.CacheMisses,
		RegistryPulls:    ab.metrics.RegistryPulls,
		RegistryPushes:   ab.metrics.RegistryPushes,
		ParallelJobs:     ab.metrics.ParallelJobs,
		PeakMemoryUsage:  ab.metrics.PeakMemoryUsage,
		NetworkTransfer:  ab.metrics.NetworkTransfer,
		CompressionRatio: ab.metrics.CompressionRatio,
	}
}

// ResetMetrics clears all build metrics.
func (ab *AdvancedBuilder) ResetMetrics() {
	ab.metrics.mutex.Lock()
	defer ab.metrics.mutex.Unlock()
	
	ab.metrics = &BuildMetrics{}
}

// generateCacheKey creates a cache key for the current build configuration.
func (ab *AdvancedBuilder) generateCacheKey() string {
	// This would be implemented with a proper hash of the build configuration
	return fmt.Sprintf("build-%s-%d", ab.config.BaseImage, time.Now().Unix())
}

// getOptimizedBaseImage retrieves and optimizes the base image.
func (ab *AdvancedBuilder) getOptimizedBaseImage(ctx context.Context, opts *BuildOptimizations) (v1.Image, error) {
	ab.metrics.RegistryPulls++
	
	image, err := ab.registry.GetImage(ab.config.BaseImage)
	if err != nil {
		return nil, err
	}

	// Apply registry optimizations if enabled
	if opts.RegistryOptimizations {
		// This would implement registry-specific optimizations
		return image, nil
	}

	return image, nil
}

// applyLayersParallel applies layers using parallel processing.
func (ab *AdvancedBuilder) applyLayersParallel(ctx context.Context, base v1.Image, opts *BuildOptimizations) (v1.Image, error) {
	maxJobs := opts.MaxParallelJobs
	if maxJobs <= 0 {
		maxJobs = 4
	}
	ab.metrics.ParallelJobs = maxJobs

	// For this implementation, fall back to sequential
	// In a real implementation, this would use worker pools
	return ab.applyLayersSequential(ctx, base, opts)
}

// applyLayersSequential applies layers sequentially.
func (ab *AdvancedBuilder) applyLayersSequential(ctx context.Context, base v1.Image, opts *BuildOptimizations) (v1.Image, error) {
	image := base

	for _, layer := range ab.layers {
		var err error
		
		// Apply layer deduplication if enabled
		if opts.LayerDeduplication {
			if existingLayer, found := ab.findDuplicateLayer(layer); found {
				layer = existingLayer
				ab.metrics.CacheHits++
			}
		}

		// Record layer size
		if size, err := layer.Size(); err == nil {
			ab.metrics.LayerSizes = append(ab.metrics.LayerSizes, size)
		}

		// Apply the layer
		image, err = mutate.AppendLayers(image, layer)
		if err != nil {
			return nil, fmt.Errorf("failed to append layer: %w", err)
		}
	}

	// Apply image squashing if enabled
	if opts.ImageSquashing {
		image, err := ab.squashImage(image)
		if err != nil {
			return nil, fmt.Errorf("failed to squash image: %w", err)
		}
		return image, nil
	}

	return image, nil
}

// createStageBuilder creates a builder for a specific multi-stage build stage.
func (ab *AdvancedBuilder) createStageBuilder(stage BuildStage, stageResults map[string]*BuildResult) (*Builder, error) {
	config := &BuildConfig{
		BaseImage: stage.BaseRef,
		Platform:  ab.config.Platform,
	}

	builder, err := NewBuilder(config, ab.registry, ab.buildCache)
	if err != nil {
		return nil, err
	}

	// Add stage-specific layers
	for _, layer := range stage.Layers {
		builder.AddLayer(layer)
	}

	return builder, nil
}

// findDuplicateLayer looks for an existing layer with the same content.
func (ab *AdvancedBuilder) findDuplicateLayer(layer v1.Layer) (v1.Layer, bool) {
	digest, err := layer.Digest()
	if err != nil {
		return nil, false
	}
	
	return ab.buildCache.GetLayerByContentHash(digest.String())
}

// squashImage combines all layers into a single layer.
func (ab *AdvancedBuilder) squashImage(image v1.Image) (v1.Image, error) {
	// This would implement actual image squashing
	// For now, return the image as-is
	return image, nil
}

// createBuildResult creates a comprehensive build result with metrics.
func (ab *AdvancedBuilder) createBuildResult(image v1.Image) (*BuildResult, error) {
	digest, err := image.Digest()
	if err != nil {
		return nil, err
	}

	size, err := image.Size()
	if err != nil {
		return nil, err
	}

	layers, err := image.Layers()
	if err != nil {
		return nil, err
	}

	layerInfo := make([]LayerInfo, len(layers))
	for i, layer := range layers {
		digest, _ := layer.Digest()
		diffID, _ := layer.DiffID()
		size, _ := layer.Size()
		mediaType, _ := layer.MediaType()

		layerInfo[i] = LayerInfo{
			Digest:      digest,
			DiffID:      diffID,
			Size:        size,
			MediaType:   string(mediaType),
			LayerType:   LayerTypeFile, // Default type
			Description: fmt.Sprintf("Layer %d", i),
		}
	}

	return &BuildResult{
		Image:     image,
		Digest:    digest,
		Size:      size,
		Duration:  ab.metrics.TotalDuration,
		LayerInfo: layerInfo,
	}, nil
}