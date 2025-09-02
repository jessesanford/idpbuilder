// Package builder provides OCI image building functionality using go-containerregistry.
// This package supports daemonless image building with OCI compliance and tarball export.
package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

// Builder defines the interface for building OCI images from directory contents.
// It supports both direct image objects and OCI tarball output for offline distribution.
type Builder interface {
	// Build creates an OCI image from a context directory.
	// The context directory contents are packaged into a single layer.
	Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error)
	
	// BuildTarball creates an OCI tarball from a context directory.
	// The tarball can be imported into container runtimes or registries.
	BuildTarball(ctx context.Context, contextDir string, output string, opts BuildOptions) error
}

// SimpleBuilder implements the Builder interface using go-containerregistry.
// It creates images with single layers from directory contents.
type SimpleBuilder struct {
	platform     v1.Platform
	featureFlags map[string]bool
	layerFactory *LayerFactory
	configFactory *ConfigFactory
	tarballWriter *TarballWriter
}

// NewBuilder creates a new builder instance with the specified options.
// It initializes internal factories and validates the build configuration.
func NewBuilder(opts BuildOptions) (*SimpleBuilder, error) {
	// Set platform defaults if not specified
	platform := opts.Platform
	if platform.OS == "" {
		platform.OS = "linux"
	}
	if platform.Architecture == "" {
		platform.Architecture = "amd64"
	}
	
	// Validate feature flags for R307 compliance
	featureFlags := opts.FeatureFlags
	if featureFlags == nil {
		featureFlags = make(map[string]bool)
	}
	
	// Initialize factories
	layerFactory := &LayerFactory{
		preservePermissions: true,
		preserveTimestamps: false, // Normalize for reproducible builds
	}
	
	configFactory := &ConfigFactory{
		platform: platform,
	}
	
	tarballWriter := NewTarballWriter()
	
	return &SimpleBuilder{
		platform:      platform,
		featureFlags:  featureFlags,
		layerFactory:  layerFactory,
		configFactory: configFactory,
		tarballWriter: tarballWriter,
	}, nil
}

// Build creates an OCI image from the specified context directory.
// It validates the directory, creates a layer, generates configuration,
// and combines them into a complete image.
func (b *SimpleBuilder) Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error) {
	// Validate context directory
	info, err := os.Stat(contextDir)
	if err != nil {
		return nil, fmt.Errorf("context directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("context path is not a directory: %s", contextDir)
	}
	
	// Check for feature flag restrictions (R307)
	if opts.FeatureFlags["multi-stage-build"] {
		return nil, fmt.Errorf("multi-stage builds not yet implemented")
	}
	if opts.FeatureFlags["buildkit-frontend"] {
		return nil, fmt.Errorf("BuildKit frontend not yet implemented")
	}
	
	// Create layer from context directory
	layer, err := b.layerFactory.CreateLayer(contextDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from context: %w", err)
	}
	
	// Start with empty base image or use specified base
	var baseImage v1.Image = empty.Image
	if opts.BaseImage != "" {
		// TODO: Load base image from reference in future implementation
		// For now, we only support empty base images
		if opts.FeatureFlags["base-image-support"] {
			return nil, fmt.Errorf("base image support not yet implemented")
		}
	}
	
	// Add our layer to the base image
	image, err := mutate.AppendLayers(baseImage, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to add layer to image: %w", err)
	}
	
	// Generate and apply OCI configuration
	config, err := b.configFactory.GenerateConfig(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image config: %w", err)
	}
	
	image, err = b.configFactory.ApplyConfig(image, config)
	if err != nil {
		return nil, fmt.Errorf("failed to apply config to image: %w", err)
	}
	
	// Set platform information
	image = mutate.MediaType(image, "application/vnd.oci.image.manifest.v1+json")
	
	return image, nil
}

// BuildTarball creates an OCI image from the context directory and exports it as a tarball.
// This provides offline distribution capabilities for environments without registry access.
func (b *SimpleBuilder) BuildTarball(ctx context.Context, contextDir string, output string, opts BuildOptions) error {
	// Build the image first
	image, err := b.Build(ctx, contextDir, opts)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	
	// Ensure output directory exists
	outputDir := filepath.Dir(output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// Generate a reference name for the image
	ref := "localhost/built-image:latest"
	if opts.Labels["org.opencontainers.image.ref.name"] != "" {
		ref = opts.Labels["org.opencontainers.image.ref.name"]
	}
	
	// Export to tarball
	if err := b.tarballWriter.Write(image, output, ref); err != nil {
		return fmt.Errorf("failed to write tarball: %w", err)
	}
	
	return nil
}