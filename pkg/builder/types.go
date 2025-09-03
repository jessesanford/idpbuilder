package builder

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// BuildConfig represents the configuration for building container images.
type BuildConfig struct {
	BaseImage    string
	Platform     *PlatformConfig
	Labels       map[string]string
	WorkingDir   string
	Entrypoint   []string
	Cmd          []string
	Environment  map[string]string
	User         string
	ExposedPorts map[string]struct{}
}

// PlatformConfig represents platform-specific build configuration.
type PlatformConfig struct {
	Architecture string
	OS           string
	Variant      string
}

// Validate validates the build configuration.
func (bc *BuildConfig) Validate() error {
	if bc.BaseImage == "" {
		return fmt.Errorf("base image is required")
	}
	if bc.Platform == nil {
		return fmt.Errorf("platform configuration is required")
	}
	if bc.Platform.Architecture == "" {
		bc.Platform.Architecture = "amd64"
	}
	if bc.Platform.OS == "" {
		bc.Platform.OS = "linux"
	}
	return nil
}

// LayerType represents different types of layers in a container image.
type LayerType int

const (
	// LayerTypeEmpty represents an empty layer
	LayerTypeEmpty LayerType = iota
	// LayerTypeFile represents a layer containing files
	LayerTypeFile
	// LayerTypeTar represents a layer created from a tar archive
	LayerTypeTar
	// LayerTypeStream represents a streaming layer
	LayerTypeStream
)

// String returns the string representation of the layer type.
func (lt LayerType) String() string {
	switch lt {
	case LayerTypeEmpty:
		return "empty"
	case LayerTypeFile:
		return "file"
	case LayerTypeTar:
		return "tar"
	case LayerTypeStream:
		return "stream"
	default:
		return "unknown"
	}
}

// Layer interface extends v1.Layer with additional functionality.
type Layer interface {
	v1.Layer
	GetType() LayerType
}

// LayerCache provides caching functionality for layers.
type LayerCache interface {
	Get(key string) (v1.Layer, bool)
	Put(key string, layer v1.Layer)
	Clear()
	Size() int
}

// RegistryClient provides an interface for registry operations.
type RegistryClient interface {
	GetImage(ref string) (v1.Image, error)
	PushImage(ref string, image v1.Image) error
	CheckImageExists(ref string) (bool, error)
}

// Builder is the main struct for building container images.
type Builder struct {
	config   *BuildConfig
	registry RegistryClient
	cache    LayerCache
	layers   []v1.Layer
}

// BuildResult represents the result of a build operation.
type BuildResult struct {
	Image     v1.Image
	Digest    v1.Hash
	Size      int64
	Duration  time.Duration
	LayerInfo []LayerInfo
}

// BuildProgress represents progress information during build.
type BuildProgress struct {
	Step        string
	Progress    float64
	TotalSteps  int
	CurrentStep int
	Message     string
}

// LayerInfo provides information about a layer.
type LayerInfo struct {
	Digest      v1.Hash
	DiffID      v1.Hash
	Size        int64
	MediaType   string
	LayerType   LayerType
	Description string
}

// NewBuilder creates a new Builder instance.
func NewBuilder(config *BuildConfig, registry RegistryClient, cache LayerCache) (*Builder, error) {
	if config == nil {
		return nil, fmt.Errorf("build configuration is required")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	if registry == nil {
		return nil, fmt.Errorf("registry client is required")
	}

	if cache == nil {
		return nil, fmt.Errorf("layer cache is required")
	}

	return &Builder{
		config:   config,
		registry: registry,
		cache:    cache,
		layers:   make([]v1.Layer, 0),
	}, nil
}

// AddLayer adds a layer to the builder.
func (b *Builder) AddLayer(layer v1.Layer) {
	b.layers = append(b.layers, layer)
}

// Build builds the container image.
func (b *Builder) Build(ctx context.Context) (*BuildResult, error) {
	start := time.Now()

	// Get base image
	baseImage, err := b.registry.GetImage(b.config.BaseImage)
	if err != nil {
		return nil, fmt.Errorf("failed to get base image: %w", err)
	}

	// For this basic implementation, just return the base image
	// In a full implementation, this would apply layers, set configuration, etc.
	digest, err := baseImage.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get image digest: %w", err)
	}

	size, err := baseImage.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get image size: %w", err)
	}

	return &BuildResult{
		Image:    baseImage,
		Digest:   digest,
		Size:     size,
		Duration: time.Since(start),
	}, nil
}

// BuildWithProgress builds the container image with progress reporting.
func (b *Builder) BuildWithProgress(ctx context.Context, progress chan<- BuildProgress) (*BuildResult, error) {
	defer close(progress)

	// Send initial progress
	progress <- BuildProgress{
		Step:        "Starting build",
		Progress:    0.0,
		TotalSteps:  3,
		CurrentStep: 0,
		Message:     "Initializing build process",
	}

	// Send progress update
	progress <- BuildProgress{
		Step:        "Fetching base image",
		Progress:    0.33,
		TotalSteps:  3,
		CurrentStep: 1,
		Message:     fmt.Sprintf("Fetching base image: %s", b.config.BaseImage),
	}

	// Build the image
	result, err := b.Build(ctx)
	if err != nil {
		return nil, err
	}

	// Send completion progress
	progress <- BuildProgress{
		Step:        "Build complete",
		Progress:    1.0,
		TotalSteps:  3,
		CurrentStep: 3,
		Message:     "Build completed successfully",
	}

	return result, nil
}

// EmptyLayer implements Layer interface for empty layers.
type EmptyLayer struct {
	digest   v1.Hash
	diffID   v1.Hash
	size     int64
	typeName LayerType
}

// NewEmptyLayer creates a new empty layer.
func NewEmptyLayer() *EmptyLayer {
	// Standard hash for empty content
	emptyHash := v1.Hash{
		Algorithm: "sha256",
		Hex:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	}
	
	return &EmptyLayer{
		digest:   emptyHash,
		diffID:   emptyHash,
		size:     0,
		typeName: LayerTypeEmpty,
	}
}

// GetType returns the type of the layer.
func (el *EmptyLayer) GetType() LayerType {
	return el.typeName
}

// Digest returns the hash of the compressed layer.
func (el *EmptyLayer) Digest() (v1.Hash, error) {
	return el.digest, nil
}

// DiffID returns the hash of the uncompressed layer.
func (el *EmptyLayer) DiffID() (v1.Hash, error) {
	return el.diffID, nil
}

// Size returns the compressed size of the layer.
func (el *EmptyLayer) Size() (int64, error) {
	return el.size, nil
}

// MediaType returns the media type of the layer.
func (el *EmptyLayer) MediaType() (types.MediaType, error) {
	return types.DockerLayer, nil
}

// Compressed returns a reader for the compressed layer content.
func (el *EmptyLayer) Compressed() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}

// Uncompressed returns a reader for the uncompressed layer content.  
func (el *EmptyLayer) Uncompressed() (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("")), nil
}

// FileLayer represents a layer created from files.
type FileLayer struct {
	*EmptyLayer
	sourcePath string
	targetPath string
}

// NewFileLayer creates a new file layer from source to target path.
func NewFileLayer(sourcePath, targetPath string) (*FileLayer, error) {
	if sourcePath == "" {
		return nil, fmt.Errorf("source path is required")
	}
	if targetPath == "" {
		return nil, fmt.Errorf("target path is required")
	}

	base := NewEmptyLayer()
	base.typeName = LayerTypeFile

	return &FileLayer{
		EmptyLayer: base,
		sourcePath: sourcePath,
		targetPath: targetPath,
	}, nil
}

// GetType returns the type of the layer.
func (fl *FileLayer) GetType() LayerType {
	return LayerTypeFile
}