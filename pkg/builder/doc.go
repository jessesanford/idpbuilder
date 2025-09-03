/*
Package builder provides comprehensive container image building capabilities for the IDP Builder project.

This package implements a complete container image builder that supports:
- Basic image building from configurations
- Advanced multi-stage builds
- Layer caching and optimization  
- Registry integration
- Progress reporting and metrics
- Parallel processing capabilities

# Basic Usage

To create a simple container image:

	config := &BuildConfig{
		BaseImage: "alpine:latest",
		Platform: &PlatformConfig{
			Architecture: "amd64", 
			OS:          "linux",
		},
	}

	registry := NewRegistryClient(registryOptions)
	cache := NewLayerCache()
	
	builder, err := NewBuilder(config, registry, cache)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	result, err := builder.Build(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Built image: %s\n", result.Digest)

# Advanced Usage with Optimizations

For production workloads requiring advanced features:

	cache := NewAdvancedBuildCache()
	builder, err := NewAdvancedBuilder(config, registry, cache)
	if err != nil {
		log.Fatal(err)
	}

	opts := &BuildOptimizations{
		EnableBuildCache:        true,
		ParallelLayerProcessing: true,
		LayerDeduplication:      true,
		MaxParallelJobs:         4,
	}

	result, err := builder.BuildWithOptimizations(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	metrics := builder.GetMetrics()
	fmt.Printf("Build completed in %v with %d cache hits\n", 
		metrics.TotalDuration, metrics.CacheHits)

# Multi-Stage Builds

The builder supports complex multi-stage builds similar to Docker:

	stages := []BuildStage{
		{
			Name:    "builder",
			BaseRef: "golang:1.21",
			Layers:  builderLayers,
		},
		{
			Name:    "runtime", 
			BaseRef: "alpine:latest",
			Layers:  runtimeLayers,
			CopyFrom: []StageCopy{{
				FromStage: "builder",
				SrcPath:   "/app/binary",
				DestPath:  "/usr/local/bin/app",
			}},
		},
	}

	result, err := builder.MultiStageBuild(ctx, stages)

# Layer Management

The package provides several layer types for different use cases:

- EmptyLayer: For creating directory structures
- FileLayer: For adding files from the build context  
- TarLayer: For adding entire archives
- StreamLayer: For efficient streaming of large content

# Performance and Monitoring

Built-in metrics collection provides insight into build performance:

	metrics := builder.GetMetrics()
	fmt.Printf("Layers: %d, Cache hits: %d, Duration: %v\n",
		metrics.LayerCount, metrics.CacheHits, metrics.TotalDuration)

# Registry Integration

The package abstracts registry operations behind a clean interface:

	type RegistryClient interface {
		GetImage(ref string) (v1.Image, error)
		PushImage(ref string, image v1.Image) error  
		CheckImageExists(ref string) (bool, error)
	}

This allows for testing with mock registries and integration with various
registry implementations.

# Caching Strategy

The advanced cache implementation supports:
- Layer-level caching by digest
- Build result caching  
- Content-based deduplication
- Pattern-based cache invalidation

This significantly improves build performance for repeated builds and
continuous integration scenarios.

# Thread Safety

All public APIs are thread-safe and support concurrent usage. The package
uses appropriate synchronization primitives to ensure data consistency
across goroutines.

# Error Handling

The package follows Go conventions for error handling, providing detailed
error messages with context about what operation failed and why.

# Compatibility

This package is designed to be compatible with the go-containerregistry
ecosystem and produces OCI-compliant container images.
*/
package builder