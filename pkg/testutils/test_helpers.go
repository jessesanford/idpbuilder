package testutils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// ImageConfig represents configuration for creating test images
type ImageConfig struct {
	Architecture string
	OS           string
	Environment  []string
	Labels       map[string]string
	WorkingDir   string
	Entrypoint   []string
	Cmd          []string
	ExposedPorts map[string]struct{}
}

// DefaultImageConfig returns a default configuration for test images
func DefaultImageConfig() *ImageConfig {
	return &ImageConfig{
		Architecture: "amd64",
		OS:           "linux",
		Environment:  []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		Labels:       make(map[string]string),
		WorkingDir:   "/",
		ExposedPorts: make(map[string]struct{}),
	}
}

// CreateTestImage creates a simple test image with the given name and tags
func CreateTestImage(name string, tags []string) (v1.Image, error) {
	return CreateTestImageWithConfig(name, tags, DefaultImageConfig())
}

// CreateTestImageWithConfig creates a test image with custom configuration
func CreateTestImageWithConfig(name string, tags []string, config *ImageConfig) (v1.Image, error) {
	// Start with empty image
	base := empty.Image

	// Create a simple layer with some test content
	layer, err := createTestLayer(fmt.Sprintf("test-content-for-%s", name))
	if err != nil {
		return nil, fmt.Errorf("failed to create test layer: %w", err)
	}

	// Add the layer to the image
	image, err := mutate.AppendLayers(base, layer)
	if err != nil {
		return nil, fmt.Errorf("failed to append layer: %w", err)
	}

	// Configure the image
	configFile := &v1.ConfigFile{
		Architecture: config.Architecture,
		OS:           config.OS,
		Config: v1.Config{
			Env:          config.Environment,
			Labels:       config.Labels,
			WorkingDir:   config.WorkingDir,
			Entrypoint:   config.Entrypoint,
			Cmd:          config.Cmd,
			ExposedPorts: config.ExposedPorts,
		},
		RootFS: v1.RootFS{
			Type: "layers",
		},
		History: []v1.History{
			{
				Created:   v1.Time{Time: time.Now()},
				CreatedBy: fmt.Sprintf("test-image-builder for %s", name),
				Comment:   "Test image layer",
			},
		},
	}

	// Apply configuration
	image, err = mutate.ConfigFile(image, configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to set config: %w", err)
	}

	return image, nil
}

// createTestLayer creates a simple tar.gz layer with test content
func createTestLayer(content string) (v1.Layer, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	// Create a simple file in the layer
	hdr := &tar.Header{
		Name: "test-file.txt",
		Mode: 0644,
		Size: int64(len(content)),
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return nil, fmt.Errorf("failed to write tar header: %w", err)
	}

	if _, err := tw.Write([]byte(content)); err != nil {
		return nil, fmt.Errorf("failed to write tar content: %w", err)
	}

	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close tar writer: %w", err)
	}

	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %w", err)
	}

	// Create layer from the tar.gz content
	layer, err := tarball.LayerFromReader(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create layer from tarball: %w", err)
	}

	return layer, nil
}

// CreateMultiArchImage creates a multi-architecture image index
func CreateMultiArchImage(name string, platforms []Platform) (v1.ImageIndex, error) {
	var index v1.ImageIndex = empty.Index

	var addenda []mutate.IndexAddendum
	for _, platform := range platforms {
		config := DefaultImageConfig()
		config.Architecture = platform.Architecture
		config.OS = platform.OS

		image, err := CreateTestImageWithConfig(
			fmt.Sprintf("%s-%s-%s", name, platform.OS, platform.Architecture),
			[]string{},
			config,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create image for platform %s/%s: %w", platform.OS, platform.Architecture, err)
		}

		// Add image to index with platform info
		desc, err := createDescriptorForImage(image, &platform)
		if err != nil {
			return nil, fmt.Errorf("failed to create descriptor: %w", err)
		}

		addenda = append(addenda, mutate.IndexAddendum{
			Add:        image,
			Descriptor: *desc,
		})
	}

	index = mutate.AppendManifests(index, addenda...)
	return index, nil
}

// Platform represents a target platform for multi-arch images
type Platform struct {
	Architecture string
	OS           string
	Variant      string
}

// createDescriptorForImage creates a descriptor for an image with platform info
func createDescriptorForImage(image v1.Image, platform *Platform) (*v1.Descriptor, error) {
	digest, err := image.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get image digest: %w", err)
	}

	size, err := image.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get image size: %w", err)
	}

	mediaType, err := image.MediaType()
	if err != nil {
		return nil, fmt.Errorf("failed to get media type: %w", err)
	}

	return &v1.Descriptor{
		MediaType: mediaType,
		Size:      size,
		Digest:    digest,
		Platform: &v1.Platform{
			Architecture: platform.Architecture,
			OS:           platform.OS,
			Variant:      platform.Variant,
		},
	}, nil
}

// PushTestImage pushes a test image to the mock registry
func PushTestImage(ctx context.Context, registry *MockRegistry, ref string, image v1.Image) error {
	if registry == nil {
		return fmt.Errorf("registry cannot be nil")
	}

	if image == nil {
		return fmt.Errorf("image cannot be nil")
	}

	// Store the image in the mock registry
	return registry.StoreImage(ref, image)
}

// CreateImageWithLayers creates an image with multiple layers
func CreateImageWithLayers(name string, layerContents []string) (v1.Image, error) {
	base := empty.Image

	for i, content := range layerContents {
		layer, err := createTestLayer(fmt.Sprintf("%s-layer-%d", content, i))
		if err != nil {
			return nil, fmt.Errorf("failed to create layer %d: %w", i, err)
		}

		base, err = mutate.AppendLayers(base, layer)
		if err != nil {
			return nil, fmt.Errorf("failed to append layer %d: %w", i, err)
		}
	}

	// Set basic configuration
	config := DefaultImageConfig()
	configFile := &v1.ConfigFile{
		Architecture: config.Architecture,
		OS:           config.OS,
		Config: v1.Config{
			Env:        config.Environment,
			Labels:     config.Labels,
			WorkingDir: config.WorkingDir,
		},
		RootFS: v1.RootFS{
			Type: "layers",
		},
	}

	// Add history for each layer
	for i := range layerContents {
		configFile.History = append(configFile.History, v1.History{
			Created:   v1.Time{Time: time.Now()},
			CreatedBy: fmt.Sprintf("test-layer-%d", i),
			Comment:   fmt.Sprintf("Test layer %d for %s", i, name),
		})
	}

	image, err := mutate.ConfigFile(base, configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to set config: %w", err)
	}

	return image, nil
}

// CompareImages compares two images and returns true if they're equivalent
func CompareImages(img1, img2 v1.Image) (bool, error) {
	if img1 == nil || img2 == nil {
		return false, fmt.Errorf("images cannot be nil")
	}

	digest1, err := img1.Digest()
	if err != nil {
		return false, fmt.Errorf("failed to get digest for first image: %w", err)
	}

	digest2, err := img2.Digest()
	if err != nil {
		return false, fmt.Errorf("failed to get digest for second image: %w", err)
	}

	return digest1 == digest2, nil
}

// GetImageInfo returns basic information about an image
func GetImageInfo(image v1.Image) (*ImageInfo, error) {
	if image == nil {
		return nil, fmt.Errorf("image cannot be nil")
	}

	digest, err := image.Digest()
	if err != nil {
		return nil, fmt.Errorf("failed to get digest: %w", err)
	}

	size, err := image.Size()
	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	layers, err := image.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to get layers: %w", err)
	}

	configFile, err := image.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	return &ImageInfo{
		Digest:       digest,
		Size:         size,
		LayerCount:   len(layers),
		Architecture: configFile.Architecture,
		OS:           configFile.OS,
		Created:      configFile.Created.Time,
	}, nil
}

// ImageInfo contains basic information about an image
type ImageInfo struct {
	Digest       v1.Hash
	Size         int64
	LayerCount   int
	Architecture string
	OS           string
	Created      time.Time
}

// String returns a string representation of the image info
func (info *ImageInfo) String() string {
	return fmt.Sprintf("Image{digest=%s, size=%d, layers=%d, arch=%s, os=%s}",
		info.Digest.String()[:12], info.Size, info.LayerCount, info.Architecture, info.OS)
}