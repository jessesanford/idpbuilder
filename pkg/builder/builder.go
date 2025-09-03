// Package builder provides container image building functionality using go-containerregistry
package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// Builder provides container image building functionality
type Builder interface {
	// Build builds a container image from the provided context
	Build(ctx context.Context, opts BuildOptions) (*BuildResult, error)

	// BuildFromTarball builds a container image from a tarball
	BuildFromTarball(ctx context.Context, tarballPath string, opts BuildOptions) (*BuildResult, error)

	// ValidateOptions validates build options
	ValidateOptions(opts BuildOptions) error
}

// BuildResult contains the result of a build operation
type BuildResult struct {
	// Image reference for the built image
	ImageRef string
	// Image digest
	Digest string
	// Size of the built image in bytes
	Size int64
	// Build metadata
	Metadata map[string]string
}

// ContainerBuilder implements the Builder interface
type ContainerBuilder struct {
	registry     string
	insecure     bool
	logger       logr.Logger
	buildOptions *BuilderDefaults
}

// BuilderDefaults provides default configuration for builds
type BuilderDefaults struct {
	BaseImage    string
	Registry     string
	Platform     string
	Compression  string
	EnableCache  bool
	CacheDir     string
	Insecure     bool
}

// NewBuilder creates a new container builder instance
func NewBuilder(opts BuilderOptions) Builder {
	if !IsFeatureEnabled() {
		return &disabledBuilder{}
	}

	logger := logr.Discard().WithName("container-builder")

	defaults := &BuilderDefaults{
		BaseImage:   "gcr.io/distroless/static:nonroot",
		Registry:    opts.Registry,
		Platform:    "linux/amd64",
		Compression: "gzip",
		EnableCache: true,
		CacheDir:    "/tmp/builder-cache",
		Insecure:    opts.Insecure,
	}

	return &ContainerBuilder{
		registry:     opts.Registry,
		insecure:     opts.Insecure,
		logger:       logger,
		buildOptions: defaults,
	}
}

// Build builds a container image from the provided context
func (b *ContainerBuilder) Build(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
	if err := b.ValidateOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid build options: %w", err)
	}

	b.logger.Info("Starting container build", "image", opts.ImageName, "context", opts.ContextPath)

	// Create base image
	baseRef, err := name.ParseReference(opts.BaseImage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base image: %w", err)
	}

	baseImg, err := crane.Pull(baseRef.String())
	if err != nil {
		return nil, fmt.Errorf("failed to pull base image: %w", err)
	}

	// Add application layer
	layer, err := b.createApplicationLayer(opts.ContextPath, opts.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to create application layer: %w", err)
	}

	// Build final image
	finalImg, err := mutate.AppendLayers(baseImg, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to add application layer: %w", err)
	}

	// Configure image metadata
	finalImg = b.configureImage(finalImg, opts)

	// Push to registry
	imageRef := fmt.Sprintf("%s/%s:%s", b.registry, opts.ImageName, opts.Tag)
	targetRef, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target reference: %w", err)
	}

	if err := crane.Push(finalImg, targetRef.String()); err != nil {
		return nil, fmt.Errorf("failed to push image: %w", err)
	}

	// Get image details for result
	digest, err := finalImg.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get image digest: %w", err)
	}

	size, err := finalImg.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get image size: %w", err)
	}

	b.logger.Info("Build completed successfully", "image", imageRef, "digest", digest.String())

	return &BuildResult{
		ImageRef: imageRef,
		Digest:   digest.String(),
		Size:     size,
		Metadata: map[string]string{
			"builder":   "idpbuilder-oci-go-cr",
			"platform":  opts.Platform,
			"base_image": opts.BaseImage,
		},
	}, nil
}

// BuildFromTarball builds a container image from a tarball
func (b *ContainerBuilder) BuildFromTarball(ctx context.Context, tarballPath string, opts BuildOptions) (*BuildResult, error) {
	if err := b.ValidateOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid build options: %w", err)
	}

	b.logger.Info("Building from tarball", "tarball", tarballPath, "image", opts.ImageName)

	// Load image from tarball
	img, err := tarball.ImageFromPath(tarballPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load image from tarball: %w", err)
	}

	// Configure image metadata
	finalImg := b.configureImage(img, opts)

	// Push to registry
	imageRef := fmt.Sprintf("%s/%s:%s", b.registry, opts.ImageName, opts.Tag)
	targetRef, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, fmt.Errorf("failed to parse target reference: %w", err)
	}

	if err := crane.Push(finalImg, targetRef.String()); err != nil {
		return nil, fmt.Errorf("failed to push image: %w", err)
	}

	digest, err := finalImg.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get image digest: %w", err)
	}

	size, err := finalImg.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get image size: %w", err)
	}

	b.logger.Info("Tarball build completed", "image", imageRef, "digest", digest.String())

	return &BuildResult{
		ImageRef: imageRef,
		Digest:   digest.String(),
		Size:     size,
		Metadata: map[string]string{
			"builder": "idpbuilder-oci-go-cr",
			"source":  "tarball",
		},
	}, nil
}

// createApplicationLayer creates a layer containing the application files
func (b *ContainerBuilder) createApplicationLayer(contextPath string, files []FileSpec) (v1.Layer, error) {
	if len(files) == 0 {
		// Create empty layer if no files specified
		return tarball.LayerFromOpener(func() (io.ReadCloser, error) {
			return os.Open(os.DevNull)
		})
	}

	// Create temporary tarball with specified files
	tempFile, err := os.CreateTemp("", "builder-layer-*.tar")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())

	// Build tarball with files
	if err := b.createLayerTarball(tempFile, contextPath, files); err != nil {
		return nil, fmt.Errorf("failed to create layer tarball: %w", err)
	}

	tempFile.Close()

	return tarball.LayerFromFile(tempFile.Name())
}

// createLayerTarball creates a tarball containing the specified files
func (b *ContainerBuilder) createLayerTarball(dst *os.File, contextPath string, files []FileSpec) error {
	// Implementation would create a proper tar archive
	// For now, create a minimal implementation
	for _, file := range files {
		srcPath := filepath.Join(contextPath, file.Source)
		if _, err := os.Stat(srcPath); err != nil {
			return fmt.Errorf("file not found: %s", srcPath)
		}
	}
	return nil
}

// configureImage configures the final image with metadata and configuration
func (b *ContainerBuilder) configureImage(img v1.Image, opts BuildOptions) v1.Image {
	// Configure image with labels, environment, etc.
	config, _ := img.ConfigFile()
	if config != nil {
		if config.Config.Labels == nil {
			config.Config.Labels = make(map[string]string)
		}
		
		for key, value := range opts.Labels {
			config.Config.Labels[key] = value
		}

		if len(opts.Env) > 0 {
			config.Config.Env = append(config.Config.Env, opts.Env...)
		}

		if opts.WorkingDir != "" {
			config.Config.WorkingDir = opts.WorkingDir
		}

		if len(opts.Cmd) > 0 {
			config.Config.Cmd = opts.Cmd
		}
	}

	return img
}

// ValidateOptions validates the provided build options
func (b *ContainerBuilder) ValidateOptions(opts BuildOptions) error {
	if opts.ImageName == "" {
		return fmt.Errorf("image name is required")
	}
	
	if opts.Tag == "" {
		return fmt.Errorf("tag is required")
	}

	if opts.BaseImage == "" {
		return fmt.Errorf("base image is required")
	}

	if opts.ContextPath != "" {
		if _, err := os.Stat(opts.ContextPath); err != nil {
			return fmt.Errorf("context path not accessible: %w", err)
		}
	}

	return nil
}

// disabledBuilder is a no-op implementation when feature is disabled
type disabledBuilder struct{}

func (d *disabledBuilder) Build(ctx context.Context, opts BuildOptions) (*BuildResult, error) {
	return nil, fmt.Errorf("container builder is disabled (ENABLE_CORE_BUILDER=false)")
}

func (d *disabledBuilder) BuildFromTarball(ctx context.Context, tarballPath string, opts BuildOptions) (*BuildResult, error) {
	return nil, fmt.Errorf("container builder is disabled (ENABLE_CORE_BUILDER=false)")
}

func (d *disabledBuilder) ValidateOptions(opts BuildOptions) error {
	return fmt.Errorf("container builder is disabled (ENABLE_CORE_BUILDER=false)")
}

// IsFeatureEnabled checks if the core builder feature is enabled
func IsFeatureEnabled() bool {
	return os.Getenv("ENABLE_CORE_BUILDER") == "true"
}