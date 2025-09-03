package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

// Builder is the main struct for building container images.
type Builder struct {
	// config contains the build configuration
	config *BuildConfig

	// baseImage is the base image to build upon
	baseImage v1.Image

	// layers contains the accumulated layers
	layers []v1.Layer

	// imageConfig holds the image configuration
	imageConfig *v1.Config

	// buildContext is the path to the build context
	buildContext string

	// cache provides layer caching functionality
	cache LayerCache

	// registry provides access to remote registries
	registry RegistryClient
}

// RegistryClient provides an interface for registry operations.
type RegistryClient interface {
	// GetImage retrieves an image from the registry
	GetImage(ref string) (v1.Image, error)
	
	// PushImage pushes an image to the registry
	PushImage(ref string, image v1.Image) error
	
	// CheckImageExists checks if an image exists in the registry
	CheckImageExists(ref string) (bool, error)
}

// LayerCache provides caching functionality for layers.
type LayerCache interface {
	// GetLayer retrieves a cached layer by its digest
	GetLayer(digest v1.Hash) (v1.Layer, error)
	
	// PutLayer stores a layer in the cache
	PutLayer(digest v1.Hash, layer v1.Layer) error
	
	// HasLayer checks if a layer exists in the cache
	HasLayer(digest v1.Hash) bool
	
	// Clear clears the cache
	Clear() error
}

// BuildResult represents the result of a build operation.
type BuildResult struct {
	// Image is the built container image
	Image v1.Image
	
	// Digest is the SHA256 digest of the built image
	Digest v1.Hash
	
	// Size is the uncompressed size of the image
	Size int64
	
	// LayerCount is the number of layers in the image
	LayerCount int
	
	// BuildTime is how long the build took
	BuildTime time.Duration
	
	// Tags are the tags applied to the image
	Tags []string
}

// NewBuilder creates a new Builder instance with the given configuration and options.
func NewBuilder(config *BuildConfig, opts ...BuilderOption) (*Builder, error) {
	if config == nil {
		return nil, fmt.Errorf("build config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid build config: %w", err)
	}

	builder := &Builder{
		config:       config.Clone(),
		layers:       make([]v1.Layer, 0),
		imageConfig:  &v1.Config{},
		buildContext: config.ContextPath,
	}

	// Apply builder options
	for _, opt := range opts {
		if err := opt(builder); err != nil {
			return nil, fmt.Errorf("failed to apply builder option: %w", err)
		}
	}

	// Initialize default components if not provided
	if builder.cache == nil {
		builder.cache = NewMemoryCache()
	}

	if builder.registry == nil {
		builder.registry = NewDefaultRegistryClient(config.Registry)
	}

	return builder, nil
}

// BuilderOption defines a functional option for configuring a Builder.
type BuilderOption func(*Builder) error

// WithCache sets a custom layer cache for the builder.
func WithCache(cache LayerCache) BuilderOption {
	return func(b *Builder) error {
		if cache == nil {
			return fmt.Errorf("cache cannot be nil")
		}
		b.cache = cache
		return nil
	}
}

// WithRegistry sets a custom registry client for the builder.
func WithRegistry(registry RegistryClient) BuilderOption {
	return func(b *Builder) error {
		if registry == nil {
			return fmt.Errorf("registry client cannot be nil")
		}
		b.registry = registry
		return nil
	}
}

// WithBaseImage sets a custom base image for the build.
func WithBaseImage(image v1.Image) BuilderOption {
	return func(b *Builder) error {
		if image == nil {
			return fmt.Errorf("base image cannot be nil")
		}
		b.baseImage = image
		return nil
	}
}

// Build executes the build process and returns the built image.
func (b *Builder) Build(ctx context.Context) (*BuildResult, error) {
	startTime := time.Now()

	// Validate build context
	if err := b.validateBuildContext(); err != nil {
		return nil, fmt.Errorf("build context validation failed: %w", err)
	}

	// Load or create base image
	if err := b.prepareBaseImage(ctx); err != nil {
		return nil, fmt.Errorf("failed to prepare base image: %w", err)
	}

	// Initialize image configuration from base
	if err := b.initializeImageConfig(); err != nil {
		return nil, fmt.Errorf("failed to initialize image config: %w", err)
	}

	// Apply build arguments and labels
	b.applyBuildSettings()

	// Build layers (this would be extended in SPLIT-003 for actual Dockerfile processing)
	if err := b.buildLayers(ctx); err != nil {
		return nil, fmt.Errorf("failed to build layers: %w", err)
	}

	// Create the final image
	image, err := b.createImage()
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %w", err)
	}

	// Calculate image digest and size
	digest, err := image.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate image digest: %w", err)
	}

	size, err := image.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate image size: %w", err)
	}

	layers, err := image.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to get image layers: %w", err)
	}

	buildTime := time.Since(startTime)

	return &BuildResult{
		Image:      image,
		Digest:     digest,
		Size:       size,
		LayerCount: len(layers),
		BuildTime:  buildTime,
		Tags:       b.config.Tags,
	}, nil
}

// validateBuildContext validates that the build context exists and is accessible.
func (b *Builder) validateBuildContext() error {
	contextDir, err := b.config.ContextDir()
	if err != nil {
		return fmt.Errorf("failed to resolve context directory: %w", err)
	}

	// Check if context directory exists
	if _, err := os.Stat(contextDir); os.IsNotExist(err) {
		return fmt.Errorf("build context directory does not exist: %s", contextDir)
	}

	// Check if Dockerfile exists
	dockerfilePath, err := b.config.DockerfilePath()
	if err != nil {
		return fmt.Errorf("failed to resolve Dockerfile path: %w", err)
	}

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("Dockerfile does not exist: %s", dockerfilePath)
	}

	return nil
}

// prepareBaseImage loads or creates the base image for the build.
func (b *Builder) prepareBaseImage(ctx context.Context) error {
	if b.baseImage != nil {
		// Base image already provided
		return nil
	}

	// For now, use an empty base image. In a full implementation,
	// this would parse the Dockerfile and determine the base image
	b.baseImage = empty.Image

	return nil
}

// initializeImageConfig initializes the image configuration from the base image.
func (b *Builder) initializeImageConfig() error {
	config, err := b.baseImage.ConfigFile()
	if err != nil {
		return fmt.Errorf("failed to get base image config: %w", err)
	}

	// Clone the base config
	b.imageConfig = config.Config.DeepCopy()

	return nil
}

// applyBuildSettings applies build arguments, labels, and other settings to the image config.
func (b *Builder) applyBuildSettings() {
	// Apply labels
	if b.imageConfig.Labels == nil {
		b.imageConfig.Labels = make(map[string]string)
	}

	for key, value := range b.config.Labels {
		b.imageConfig.Labels[key] = value
	}

	// Add build metadata
	b.imageConfig.Labels["io.cnoe.build.timestamp"] = time.Now().UTC().Format(time.RFC3339)
	b.imageConfig.Labels["io.cnoe.build.context"] = b.config.ContextPath
	b.imageConfig.Labels["io.cnoe.build.dockerfile"] = b.config.Dockerfile

	// Apply environment variables from build args (simplified)
	for key, value := range b.config.BuildArgs {
		envVar := fmt.Sprintf("%s=%s", key, value)
		b.imageConfig.Env = append(b.imageConfig.Env, envVar)
	}
}

// buildLayers processes the build context and creates the necessary layers.
// This is a simplified implementation that would be extended in SPLIT-003.
func (b *Builder) buildLayers(ctx context.Context) error {
	// This is a placeholder implementation. In a full builder, this would:
	// 1. Parse the Dockerfile
	// 2. Process each instruction
	// 3. Create layers for each RUN, COPY, ADD instruction
	// 4. Handle multi-stage builds
	// 5. Apply layer caching

	// For now, create a simple layer with build metadata
	layer, err := b.createMetadataLayer()
	if err != nil {
		return fmt.Errorf("failed to create metadata layer: %w", err)
	}

	b.layers = append(b.layers, layer)
	return nil
}

// createMetadataLayer creates a layer containing build metadata.
func (b *Builder) createMetadataLayer() (v1.Layer, error) {
	// This would be implemented in SPLIT-003 with proper tarball operations
	// For now, return an empty layer
	return NewEmptyLayer(), nil
}

// createImage creates the final container image from the base image and built layers.
func (b *Builder) createImage() (v1.Image, error) {
	image := b.baseImage

	// Add all built layers to the image
	for _, layer := range b.layers {
		var err error
		image, err = mutate.AppendLayers(image, layer)
		if err != nil {
			return nil, fmt.Errorf("failed to append layer: %w", err)
		}
	}

	// Apply the final configuration
	image, err := mutate.Config(image, *b.imageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to apply image config: %w", err)
	}

	return image, nil
}

// Push pushes the built image to the configured registry with all specified tags.
func (b *Builder) Push(ctx context.Context, image v1.Image) error {
	if b.registry == nil {
		return fmt.Errorf("no registry client configured")
	}

	for _, tag := range b.config.Tags {
		if err := b.registry.PushImage(tag, image); err != nil {
			return fmt.Errorf("failed to push image with tag %s: %w", tag, err)
		}
	}

	return nil
}

// GetBuildContext returns the path to the build context directory.
func (b *Builder) GetBuildContext() string {
	return b.buildContext
}

// GetConfig returns a copy of the build configuration.
func (b *Builder) GetConfig() *BuildConfig {
	return b.config.Clone()
}

// GetLayers returns a copy of the current layers.
func (b *Builder) GetLayers() []v1.Layer {
	layers := make([]v1.Layer, len(b.layers))
	copy(layers, b.layers)
	return layers
}

// AddLayer adds a custom layer to the build.
func (b *Builder) AddLayer(layer v1.Layer) error {
	if layer == nil {
		return fmt.Errorf("layer cannot be nil")
	}

	b.layers = append(b.layers, layer)
	return nil
}

// SetImageConfig allows customization of the image configuration.
func (b *Builder) SetImageConfig(config v1.Config) {
	b.imageConfig = &config
}

// Close cleans up any resources used by the builder.
func (b *Builder) Close() error {
	if b.cache != nil {
		return b.cache.Clear()
	}
	return nil
}

// BuildWithProgress builds an image with progress reporting.
func (b *Builder) BuildWithProgress(ctx context.Context, progress io.Writer) (*BuildResult, error) {
	if progress != nil {
		fmt.Fprintf(progress, "Starting build at %s\n", time.Now().Format(time.RFC3339))
	}

	result, err := b.Build(ctx)
	if err != nil {
		if progress != nil {
			fmt.Fprintf(progress, "Build failed: %v\n", err)
		}
		return nil, err
	}

	if progress != nil {
		fmt.Fprintf(progress, "Build completed in %v\n", result.BuildTime)
		fmt.Fprintf(progress, "Image digest: %s\n", result.Digest.String())
		fmt.Fprintf(progress, "Image size: %d bytes\n", result.Size)
		fmt.Fprintf(progress, "Layers: %d\n", result.LayerCount)
	}

	return result, nil
}