package builder

import (
	"context"
	"fmt"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

// BuildResult contains the result of a build operation.
type BuildResult struct {
	Image     v1.Image    `json:"-"`
	Digest    v1.Hash     `json:"digest"`
	Size      int64       `json:"size"`
	Layers    []LayerInfo `json:"layers"`
	Config    v1.Config   `json:"config"`
	Tags      []string    `json:"tags"`
}

// LayerInfo contains information about a layer in the built image.
type LayerInfo struct {
	Digest v1.Hash `json:"digest"`
	Size   int64   `json:"size"`
}

// BuildFromContext creates a container image from the specified context directory and options.
func (b *SimpleBuilder) BuildFromContext(ctx context.Context, contextDir string, opts BuildOptions) (v1.Image, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if contextDir == "" {
		return nil, fmt.Errorf("context directory cannot be empty")
	}

	// Validate build options
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid build options: %w", err)
	}

	// Validate context directory
	if err := b.ValidateContext(contextDir); err != nil {
		return nil, fmt.Errorf("invalid context directory: %w", err)
	}

	// Check if tarball export feature is enabled
	if !b.featureFlags[FeatureTarballExport] {
		return nil, fmt.Errorf("tarball export feature is not enabled")
	}

	// Merge build options with defaults
	mergedOpts, err := b.MergeBuildOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to merge build options: %w", err)
	}

	// Create base image
	baseImage := empty.Image

	// Create configuration
	config, err := b.createImageConfig(*mergedOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create image config: %w", err)
	}

	// Apply configuration to base image
	image, err := mutate.Config(baseImage, config)
	if err != nil {
		return nil, fmt.Errorf("failed to apply config to image: %w", err)
	}

	// Create layer from context directory
	layer, err := b.createLayerFromContext(ctx, contextDir, *mergedOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from context: %w", err)
	}

	// Add layer to image
	image, err = mutate.AppendLayers(image, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to add layer to image: %w", err)
	}

	// Apply any additional mutations based on options
	finalImage, err := b.applyImageMutations(image, *mergedOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to apply image mutations: %w", err)
	}

	return finalImage, nil
}

// BuildWithResult creates a container image and returns detailed build result.
func (b *SimpleBuilder) BuildWithResult(ctx context.Context, contextDir string, opts BuildOptions) (*BuildResult, error) {
	image, err := b.Build(ctx, contextDir, opts)
	if err != nil {
		return nil, err
	}

	// Calculate image digest
	digest, err := image.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate image digest: %w", err)
	}

	// Calculate image size
	manifest, err := image.Manifest()
	if err != nil {
		return nil, fmt.Errorf("failed to get image manifest: %w", err)
	}

	var totalSize int64
	var layerInfos []LayerInfo

	for _, layer := range manifest.Layers {
		layerInfos = append(layerInfos, LayerInfo{
			Digest: layer.Digest,
			Size:   layer.Size,
		})
		totalSize += layer.Size
	}

	// Get image config
	configFile, err := image.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get image config: %w", err)
	}

	return &BuildResult{
		Image:  image,
		Digest: digest,
		Size:   totalSize,
		Layers: layerInfos,
		Config: configFile.Config,
		Tags:   opts.Tags,
	}, nil
}

// createImageConfig creates a v1.Config from build options.
func (b *SimpleBuilder) createImageConfig(opts BuildOptions) (v1.Config, error) {
	config := v1.Config{
		Labels:       make(map[string]string),
		Env:          make([]string, 0),
		ExposedPorts: make(map[string]struct{}),
	}

	// Set working directory
	if opts.WorkingDir != "" {
		config.WorkingDir = opts.WorkingDir
	}

	// Set user
	if opts.User != "" {
		config.User = opts.User
	}

	// Set entrypoint
	if len(opts.Entrypoint) > 0 {
		config.Entrypoint = make([]string, len(opts.Entrypoint))
		copy(config.Entrypoint, opts.Entrypoint)
	}

	// Set cmd
	if len(opts.Cmd) > 0 {
		config.Cmd = make([]string, len(opts.Cmd))
		copy(config.Cmd, opts.Cmd)
	}

	// Add labels
	for key, value := range opts.Labels {
		config.Labels[key] = value
	}

	// Add environment variables
	for key, value := range opts.Environment {
		config.Env = append(config.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Add exposed ports
	for _, port := range opts.ExposedPorts {
		config.ExposedPorts[port] = struct{}{}
	}

	return config, nil
}

// createLayerFromContext creates a layer from the context directory.
func (b *SimpleBuilder) createLayerFromContext(ctx context.Context, contextDir string, opts BuildOptions) (v1.Layer, error) {
	// Create tarball writer with appropriate options
	tarballOpts := TarballOptions{
		TimestampPolicy: TimestampEpoch, // For reproducible builds
		PreserveOwner:   false,
		DefaultMode:     0644,
	}

	writer := NewTarballWriter(tarballOpts)

	// Create layer from the entire context directory
	layer, err := writer.CreateLayerFromDirectory(ctx, contextDir, "/")
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from directory: %w", err)
	}

	return layer, nil
}

// applyImageMutations applies any additional mutations to the image based on options.
func (b *SimpleBuilder) applyImageMutations(image v1.Image, opts BuildOptions) (v1.Image, error) {
	// For now, return the image as-is
	// Future enhancements could apply additional mutations here
	return image, nil
}

// ValidateFeatureFlags validates that required features are enabled for the build.
func (b *SimpleBuilder) ValidateFeatureFlags(requiredFeatures []string) error {
	var missingFeatures []string
	for _, feature := range requiredFeatures {
		if !b.featureFlags[feature] {
			missingFeatures = append(missingFeatures, feature)
		}
	}
	if len(missingFeatures) > 0 {
		return fmt.Errorf("missing required features: %s", strings.Join(missingFeatures, ", "))
	}
	return nil
}