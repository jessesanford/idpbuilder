package builder

import (
	"context"
	"fmt"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// Feature flags for controlling functionality availability
const (
	FeatureTarballExport = "tarball_export" // Disabled in Split 001
	FeatureLayerCaching  = "layer_caching"  // Disabled in Split 001
	FeatureMultiLayer    = "multi_layer"    // Disabled in Split 001
)

// Builder defines the interface for building OCI container images.
type Builder interface {
	// Build creates a container image from the specified context directory
	// using the provided build options.
	Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error)

	// GetSupportedFeatures returns a list of features supported by this builder implementation.
	GetSupportedFeatures() []string

	// IsFeatureSupported checks if a specific feature is supported.
	IsFeatureSupported(feature string) bool

	// GetConfig returns the current configuration factory used by the builder.
	GetConfig() *ConfigFactory
}

// SimpleBuilder is a basic implementation of the Builder interface.
type SimpleBuilder struct {
	configFactory *ConfigFactory
	featureFlags  map[string]bool
	buildOptions  *BuildOptions
}

// BuilderOption defines a function type for configuring a SimpleBuilder.
type BuilderOption func(*SimpleBuilder) error

// NewBuilder creates a new SimpleBuilder with the specified options.
func NewBuilder(opts ...BuilderOption) (Builder, error) {
	builder := &SimpleBuilder{
		configFactory: NewConfigFactory(),
		featureFlags:  make(map[string]bool),
		buildOptions:  NewBuildOptions(),
	}

	// Set default feature flags (all disabled in Split 001)
	builder.featureFlags[FeatureTarballExport] = false
	builder.featureFlags[FeatureLayerCaching] = false
	builder.featureFlags[FeatureMultiLayer] = false

	// Apply options
	for _, opt := range opts {
		if err := opt(builder); err != nil {
			return nil, fmt.Errorf("failed to apply builder option: %w", err)
		}
	}

	return builder, nil
}

// WithConfigFactory sets a custom configuration factory for the builder.
func WithConfigFactory(factory *ConfigFactory) BuilderOption {
	return func(b *SimpleBuilder) error {
		if factory == nil {
			return fmt.Errorf("config factory cannot be nil")
		}
		b.configFactory = factory
		return nil
	}
}

// WithFeatureFlag enables or disables a specific feature flag.
func WithFeatureFlag(feature string, enabled bool) BuilderOption {
	return func(b *SimpleBuilder) error {
		if feature == "" {
			return fmt.Errorf("feature name cannot be empty")
		}
		b.featureFlags[feature] = enabled
		return nil
	}
}

// WithDefaultBuildOptions sets default build options for the builder.
func WithDefaultBuildOptions(opts *BuildOptions) BuilderOption {
	return func(b *SimpleBuilder) error {
		if opts == nil {
			return fmt.Errorf("build options cannot be nil")
		}
		if err := opts.Validate(); err != nil {
			return fmt.Errorf("invalid build options: %w", err)
		}
		b.buildOptions = opts
		return nil
	}
}

// Build creates a container image from the specified context directory.
// This is a stub implementation for Split 001 - full implementation in Split 002.
func (b *SimpleBuilder) Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	if contextDir == "" {
		return nil, fmt.Errorf("context directory cannot be empty")
	}

	// Validate options
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid build options: %w", err)
	}

	// Check if required features are enabled
	if !b.featureFlags[FeatureTarballExport] {
		return nil, fmt.Errorf("tarball export not enabled - full implementation in Split 002")
	}

	// Stub for now - will be completed in Split 002
	return nil, fmt.Errorf("Build method not fully implemented in Split 001 - will be completed in Split 002")
}

// GetSupportedFeatures returns a list of features supported by this builder implementation.
func (b *SimpleBuilder) GetSupportedFeatures() []string {
	var features []string
	for feature, enabled := range b.featureFlags {
		if enabled {
			features = append(features, feature)
		}
	}
	return features
}

// IsFeatureSupported checks if a specific feature is supported and enabled.
func (b *SimpleBuilder) IsFeatureSupported(feature string) bool {
	return b.featureFlags[feature]
}

// GetConfig returns the current configuration factory used by the builder.
func (b *SimpleBuilder) GetConfig() *ConfigFactory {
	return b.configFactory
}

// ValidateContext performs basic validation on the build context directory.
func (b *SimpleBuilder) ValidateContext(contextDir string) error {
	if contextDir == "" {
		return fmt.Errorf("context directory cannot be empty")
	}
	return nil
}

// MergeBuildOptions merges the provided options with the builder's default options.
func (b *SimpleBuilder) MergeBuildOptions(opts BuildOptions) (*BuildOptions, error) {
	result := NewBuildOptions()

	// Copy defaults if they exist
	if b.buildOptions != nil {
		*result = *b.buildOptions
		// Deep copy slices and maps
		result.Labels = make(map[string]string)
		result.Environment = make(map[string]string)
		result.FeatureFlags = make(map[string]bool)
		result.BuildArgs = make(map[string]string)
		for k, v := range b.buildOptions.Labels {
			result.Labels[k] = v
		}
		for k, v := range b.buildOptions.Environment {
			result.Environment[k] = v
		}
		for k, v := range b.buildOptions.FeatureFlags {
			result.FeatureFlags[k] = v
		}
		for k, v := range b.buildOptions.BuildArgs {
			result.BuildArgs[k] = v
		}
	}

	// Override with provided options
	if opts.Platform != nil {
		result.Platform = opts.Platform
	}
	if opts.WorkingDir != "" {
		result.WorkingDir = opts.WorkingDir
	}
	if opts.User != "" {
		result.User = opts.User
	}
	if len(opts.Entrypoint) > 0 {
		result.Entrypoint = append([]string(nil), opts.Entrypoint...)
	}
	if len(opts.Cmd) > 0 {
		result.Cmd = append([]string(nil), opts.Cmd...)
	}
	if len(opts.Tags) > 0 {
		result.Tags = append([]string(nil), opts.Tags...)
	}
	if len(opts.ExposedPorts) > 0 {
		result.ExposedPorts = append([]string(nil), opts.ExposedPorts...)
	}
	if opts.ContextPath != "" {
		result.ContextPath = opts.ContextPath
	}

	// Merge maps
	for k, v := range opts.Labels {
		result.Labels[k] = v
	}
	for k, v := range opts.Environment {
		result.Environment[k] = v
	}
	for k, v := range opts.FeatureFlags {
		result.FeatureFlags[k] = v
	}
	for k, v := range opts.BuildArgs {
		result.BuildArgs[k] = v
	}

	return result, nil
}

// GetDefaultBuildOptions returns a copy of the builder's default build options.
func (b *SimpleBuilder) GetDefaultBuildOptions() *BuildOptions {
	if b.buildOptions == nil {
		return NewBuildOptions()
	}
	
	// Use MergeBuildOptions with empty opts to get a clean copy
	result, _ := b.MergeBuildOptions(BuildOptions{})
	return result
}