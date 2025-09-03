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
	Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error)
	GetSupportedFeatures() []string
	IsFeatureSupported(feature string) bool
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
	defaultOpts := NewBuildOptions()
	builder := &SimpleBuilder{
		configFactory: NewConfigFactory(),
		featureFlags:  make(map[string]bool),
		buildOptions:  defaultOpts,
	}

	// Set default feature flags (all disabled in Split 001)
	builder.featureFlags[FeatureTarballExport] = false
	builder.featureFlags[FeatureLayerCaching] = false
	builder.featureFlags[FeatureMultiLayer] = false

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

// WithFeatureFlag sets a feature flag for the builder.
func WithFeatureFlag(feature string, enabled bool) BuilderOption {
	return func(b *SimpleBuilder) error {
		if feature == "" {
			return fmt.Errorf("feature name cannot be empty")
		}
		b.featureFlags[feature] = enabled
		return nil
	}
}

// WithDefaultBuildOptions sets the default build options for the builder.
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

// Build creates a container image - now implemented in Split 002b.
func (b *SimpleBuilder) Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error) {
	// Delegate to the full implementation
	return b.BuildFromContext(ctx, contextDir, opts)
}

// GetSupportedFeatures returns a list of supported features.
func (b *SimpleBuilder) GetSupportedFeatures() []string {
	var features []string
	for feature, enabled := range b.featureFlags {
		if enabled {
			features = append(features, feature)
		}
	}
	return features
}

// IsFeatureSupported checks if a specific feature is supported.
func (b *SimpleBuilder) IsFeatureSupported(feature string) bool {
	return b.featureFlags[feature]
}

// GetConfig returns the current configuration factory.
func (b *SimpleBuilder) GetConfig() *ConfigFactory {
	return b.configFactory
}

// GetDefaultBuildOptions returns a copy of the default build options.
func (b *SimpleBuilder) GetDefaultBuildOptions() *BuildOptions {
	if b.buildOptions == nil {
		return NewBuildOptions()
	}
	// Return a deep copy to prevent mutation
	opts := *b.buildOptions
	
	// Deep copy maps
	if b.buildOptions.Labels != nil {
		opts.Labels = make(map[string]string)
		for k, v := range b.buildOptions.Labels {
			opts.Labels[k] = v
		}
	}
	if b.buildOptions.Environment != nil {
		opts.Environment = make(map[string]string)
		for k, v := range b.buildOptions.Environment {
			opts.Environment[k] = v
		}
	}
	
	// Deep copy maps
	if b.buildOptions.BuildArgs != nil {
		opts.BuildArgs = make(map[string]string)
		for k, v := range b.buildOptions.BuildArgs {
			opts.BuildArgs[k] = v
		}
	}
	if b.buildOptions.FeatureFlags != nil {
		opts.FeatureFlags = make(map[string]bool)
		for k, v := range b.buildOptions.FeatureFlags {
			opts.FeatureFlags[k] = v
		}
	}
	
	// Deep copy slices
	if b.buildOptions.Tags != nil {
		opts.Tags = make([]string, len(b.buildOptions.Tags))
		copy(opts.Tags, b.buildOptions.Tags)
	}
	if b.buildOptions.Entrypoint != nil {
		opts.Entrypoint = make([]string, len(b.buildOptions.Entrypoint))
		copy(opts.Entrypoint, b.buildOptions.Entrypoint)
	}
	if b.buildOptions.Cmd != nil {
		opts.Cmd = make([]string, len(b.buildOptions.Cmd))
		copy(opts.Cmd, b.buildOptions.Cmd)
	}
	if b.buildOptions.ExposedPorts != nil {
		opts.ExposedPorts = make([]string, len(b.buildOptions.ExposedPorts))
		copy(opts.ExposedPorts, b.buildOptions.ExposedPorts)
	}
	
	return &opts
}

// SetBuildOption allows setting individual build options on the builder.
func (b *SimpleBuilder) SetBuildOption(key string, value interface{}) error {
	if b.buildOptions == nil {
		b.buildOptions = NewBuildOptions()
	}
	
	switch key {
	case "working_dir":
		if workDir, ok := value.(string); ok {
			b.buildOptions.WorkingDir = workDir
		} else {
			return fmt.Errorf("working_dir must be a string")
		}
	case "user":
		if user, ok := value.(string); ok {
			b.buildOptions.User = user
		} else {
			return fmt.Errorf("user must be a string")
		}
	case "context_path":
		if contextPath, ok := value.(string); ok {
			b.buildOptions.ContextPath = contextPath
		} else {
			return fmt.Errorf("context_path must be a string")
		}
	default:
		return fmt.Errorf("unsupported build option: %s", key)
	}
	
	return nil
}

// GetBuildOption retrieves a build option value by key.
func (b *SimpleBuilder) GetBuildOption(key string) (interface{}, error) {
	if b.buildOptions == nil {
		return nil, fmt.Errorf("build options not initialized")
	}
	
	switch key {
	case "working_dir":
		return b.buildOptions.WorkingDir, nil
	case "user":
		return b.buildOptions.User, nil
	case "context_path":
		return b.buildOptions.ContextPath, nil
	case "platform":
		return b.buildOptions.Platform, nil
	default:
		return nil, fmt.Errorf("unsupported build option: %s", key)
	}
}

// ValidateContext validates the build context directory.
func (b *SimpleBuilder) ValidateContext(contextDir string) error {
	if contextDir == "" {
		return fmt.Errorf("context directory cannot be empty")
	}
	return nil
}

// MergeBuildOptions merges the provided options with the default build options.
func (b *SimpleBuilder) MergeBuildOptions(opts BuildOptions) (*BuildOptions, error) {
	defaults := b.GetDefaultBuildOptions()
	
	// Merge fields from opts into defaults
	if opts.Platform != nil {
		defaults.Platform = opts.Platform
	}
	if len(opts.Tags) > 0 {
		defaults.Tags = opts.Tags
	}
	if len(opts.Labels) > 0 {
		if defaults.Labels == nil {
			defaults.Labels = make(map[string]string)
		}
		for k, v := range opts.Labels {
			defaults.Labels[k] = v
		}
	}
	if len(opts.Environment) > 0 {
		if defaults.Environment == nil {
			defaults.Environment = make(map[string]string)
		}
		for k, v := range opts.Environment {
			defaults.Environment[k] = v
		}
	}
	if len(opts.ExposedPorts) > 0 {
		defaults.ExposedPorts = opts.ExposedPorts
	}
	if len(opts.BuildArgs) > 0 {
		if defaults.BuildArgs == nil {
			defaults.BuildArgs = make(map[string]string)
		}
		for k, v := range opts.BuildArgs {
			defaults.BuildArgs[k] = v
		}
	}
	if len(opts.FeatureFlags) > 0 {
		if defaults.FeatureFlags == nil {
			defaults.FeatureFlags = make(map[string]bool)
		}
		for k, v := range opts.FeatureFlags {
			defaults.FeatureFlags[k] = v
		}
	}
	if opts.WorkingDir != "" {
		defaults.WorkingDir = opts.WorkingDir
	}
	if len(opts.Entrypoint) > 0 {
		defaults.Entrypoint = opts.Entrypoint
	}
	if len(opts.Cmd) > 0 {
		defaults.Cmd = opts.Cmd
	}
	if opts.User != "" {
		defaults.User = opts.User
	}
	if opts.ContextPath != "" {
		defaults.ContextPath = opts.ContextPath
	}
	
	return defaults, nil
}