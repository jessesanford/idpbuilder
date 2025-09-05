// Package builder provides OCI image building functionality.
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

// Builder defines the interface for building OCI images.
type Builder interface {
	Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error)
	BuildTarball(ctx context.Context, contextDir string, output string, opts BuildOptions) error
}

// SimpleBuilder implements the Builder interface.
type SimpleBuilder struct {
	platform      v1.Platform
	featureFlags  map[string]bool
	layerFactory  *LayerFactory
	configFactory *ConfigFactory
	tarballWriter *TarballWriter
}

// NewBuilder creates a new builder instance.
func NewBuilder(opts BuildOptions) (*SimpleBuilder, error) {
	platform := opts.Platform
	if platform.OS == "" {
		platform.OS = "linux"
	}
	if platform.Architecture == "" {
		platform.Architecture = "amd64"
	}

	featureFlags := opts.FeatureFlags
	if featureFlags == nil {
		featureFlags = make(map[string]bool)
	}

	return &SimpleBuilder{
		platform:      platform,
		featureFlags:  featureFlags,
		layerFactory:  &LayerFactory{preservePermissions: true, preserveTimestamps: false},
		configFactory: &ConfigFactory{platform: platform},
		tarballWriter: NewTarballWriter(),
	}, nil
}

// Build creates an OCI image from the context directory.
func (b *SimpleBuilder) Build(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error) {
	info, err := os.Stat(contextDir)
	if err != nil {
		return nil, fmt.Errorf("context directory not found: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("context path is not a directory: %s", contextDir)
	}

	// Check feature flags
	if opts.FeatureFlags["multi-stage-build"] || opts.FeatureFlags["buildkit-frontend"] || opts.FeatureFlags["base-image-support"] {
		return nil, fmt.Errorf("advanced features not yet implemented")
	}

	layer, err := b.layerFactory.CreateLayer(contextDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer: %w", err)
	}

	image, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to add layer: %w", err)
	}

	config, err := b.configFactory.GenerateConfig(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to generate config: %w", err)
	}

	image, err = b.configFactory.ApplyConfig(image, config)
	if err != nil {
		return nil, fmt.Errorf("failed to apply config: %w", err)
	}

	return mutate.MediaType(image, "application/vnd.oci.image.manifest.v1+json"), nil
}

// BuildTarball creates an OCI tarball from the context directory.
func (b *SimpleBuilder) BuildTarball(ctx context.Context, contextDir string, output string, opts BuildOptions) error {
	image, err := b.Build(ctx, contextDir, opts)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	ref := "localhost/built-image:latest"
	if opts.Labels["org.opencontainers.image.ref.name"] != "" {
		ref = opts.Labels["org.opencontainers.image.ref.name"]
	}

	return b.tarballWriter.Write(image, output, ref)
}
