package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// MockRegistryClient implements RegistryClient for testing.
type MockRegistryClient struct {
	images   map[string]v1.Image
	pushCalled bool
	getImageError error
	pushImageError error
}

// GetImage returns a mock image for testing.
func (m *MockRegistryClient) GetImage(ref string) (v1.Image, error) {
	if m.getImageError != nil {
		return nil, m.getImageError
	}
	if img, exists := m.images[ref]; exists {
		return img, nil
	}
	return nil, fmt.Errorf("image not found: %s", ref)
}

// PushImage simulates pushing an image to registry.
func (m *MockRegistryClient) PushImage(ref string, image v1.Image) error {
	m.pushCalled = true
	if m.pushImageError != nil {
		return m.pushImageError
	}
	if m.images == nil {
		m.images = make(map[string]v1.Image)
	}
	m.images[ref] = image
	return nil
}

// CheckImageExists checks if an image exists in the mock registry.
func (m *MockRegistryClient) CheckImageExists(ref string) (bool, error) {
	_, exists := m.images[ref]
	return exists, nil
}

// MockLayerCache implements LayerCache for testing.
type MockLayerCache struct {
	layers map[string]v1.Layer
	getHits int
	putHits int
}

// Get retrieves a layer from the mock cache.
func (m *MockLayerCache) Get(key string) (v1.Layer, bool) {
	m.getHits++
	if layer, exists := m.layers[key]; exists {
		return layer, true
	}
	return nil, false
}

// Put stores a layer in the mock cache.
func (m *MockLayerCache) Put(key string, layer v1.Layer) {
	m.putHits++
	if m.layers == nil {
		m.layers = make(map[string]v1.Layer)
	}
	m.layers[key] = layer
}

// Clear removes all layers from the mock cache.
func (m *MockLayerCache) Clear() {
	m.layers = make(map[string]v1.Layer)
}

// Size returns the number of layers in the mock cache.
func (m *MockLayerCache) Size() int {
	return len(m.layers)
}

// TestBuilderNew tests the creation of a new Builder instance.
func TestBuilderNew(t *testing.T) {
	config := &BuildConfig{
		BaseImage: "alpine:latest",
		Platform: &PlatformConfig{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	registry := &MockRegistryClient{}
	cache := &MockLayerCache{}

	builder, err := NewBuilder(config, registry, cache)
	if err != nil {
		t.Fatalf("failed to create builder: %v", err)
	}

	if builder == nil {
		t.Fatal("builder should not be nil")
	}

	if builder.config != config {
		t.Error("builder config not set correctly")
	}
}

// TestBuilderNewWithNilConfig tests builder creation with nil config.
func TestBuilderNewWithNilConfig(t *testing.T) {
	registry := &MockRegistryClient{}
	cache := &MockLayerCache{}

	builder, err := NewBuilder(nil, registry, cache)
	if err == nil {
		t.Fatal("expected error for nil config")
	}

	if builder != nil {
		t.Error("builder should be nil when config is nil")
	}
}

// TestBuilderNewWithInvalidConfig tests builder creation with invalid config.
func TestBuilderNewWithInvalidConfig(t *testing.T) {
	config := &BuildConfig{
		// Missing required BaseImage
		Platform: &PlatformConfig{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	registry := &MockRegistryClient{}
	cache := &MockLayerCache{}

	builder, err := NewBuilder(config, registry, cache)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}

	if builder != nil {
		t.Error("builder should be nil when config is invalid")
	}
}

// TestBuilderBuild tests the basic build functionality.
func TestBuilderBuild(t *testing.T) {
	config := &BuildConfig{
		BaseImage: "alpine:latest",
		Platform: &PlatformConfig{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	registry := &MockRegistryClient{
		images: map[string]v1.Image{
			"alpine:latest": &mockImage{},
		},
	}
	cache := &MockLayerCache{}

	builder, err := NewBuilder(config, registry, cache)
	if err != nil {
		t.Fatalf("failed to create builder: %v", err)
	}

	ctx := context.Background()
	result, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	if result.Image == nil {
		t.Error("build result should contain an image")
	}
}

// TestBuilderBuildWithLayers tests building with additional layers.
func TestBuilderBuildWithLayers(t *testing.T) {
	config := &BuildConfig{
		BaseImage: "alpine:latest",
		Platform: &PlatformConfig{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	registry := &MockRegistryClient{
		images: map[string]v1.Image{
			"alpine:latest": &mockImage{},
		},
	}
	cache := &MockLayerCache{}

	builder, err := NewBuilder(config, registry, cache)
	if err != nil {
		t.Fatalf("failed to create builder: %v", err)
	}

	// Add some layers
	emptyLayer := NewEmptyLayer()
	builder.AddLayer(emptyLayer)

	ctx := context.Background()
	result, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("build with layers failed: %v", err)
	}

	if result.Image == nil {
		t.Error("build result should contain an image")
	}
}

// TestLayerCreation tests various layer creation methods.
func TestLayerCreation(t *testing.T) {
	t.Run("EmptyLayer", func(t *testing.T) {
		layer := NewEmptyLayer()
		if layer == nil {
			t.Fatal("empty layer should not be nil")
		}

		if layer.GetType() != LayerTypeEmpty {
			t.Errorf("expected layer type %v, got %v", LayerTypeEmpty, layer.GetType())
		}
	})

	t.Run("FileLayer", func(t *testing.T) {
		// Create a temporary file for testing
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		layer, err := NewFileLayer(tmpDir, "/app/")
		if err != nil {
			t.Fatalf("failed to create file layer: %v", err)
		}

		if layer == nil {
			t.Fatal("file layer should not be nil")
		}

		if layer.GetType() != LayerTypeFile {
			t.Errorf("expected layer type %v, got %v", LayerTypeFile, layer.GetType())
		}
	})
}

// TestTarballOperations tests tarball creation and extraction.
func TestTarballOperations(t *testing.T) {
	t.Run("CreateTarball", func(t *testing.T) {
		// Create a temporary directory with test files
		tmpDir := t.TempDir()
		testFile1 := filepath.Join(tmpDir, "file1.txt")
		testFile2 := filepath.Join(tmpDir, "file2.txt")
		
		if err := os.WriteFile(testFile1, []byte("content1"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
		if err := os.WriteFile(testFile2, []byte("content2"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Create tarball
		tarballPath := filepath.Join(t.TempDir(), "test.tar")
		err := CreateTarball(tmpDir, tarballPath)
		if err != nil {
			t.Fatalf("failed to create tarball: %v", err)
		}

		// Verify tarball exists and has content
		info, err := os.Stat(tarballPath)
		if err != nil {
			t.Fatalf("tarball not created: %v", err)
		}

		if info.Size() == 0 {
			t.Error("tarball should not be empty")
		}
	})

	t.Run("ExtractTarball", func(t *testing.T) {
		// Create source directory and tarball
		srcDir := t.TempDir()
		testFile := filepath.Join(srcDir, "test.txt")
		content := "test content for extraction"
		
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		tarballPath := filepath.Join(t.TempDir(), "extract_test.tar")
		if err := CreateTarball(srcDir, tarballPath); err != nil {
			t.Fatalf("failed to create tarball: %v", err)
		}

		// Extract tarball
		destDir := t.TempDir()
		if err := ExtractTarball(tarballPath, destDir); err != nil {
			t.Fatalf("failed to extract tarball: %v", err)
		}

		// Verify extracted content
		extractedFile := filepath.Join(destDir, "test.txt")
		extractedContent, err := os.ReadFile(extractedFile)
		if err != nil {
			t.Fatalf("failed to read extracted file: %v", err)
		}

		if string(extractedContent) != content {
			t.Errorf("extracted content mismatch: got %s, want %s", extractedContent, content)
		}
	})
}

// TestCompressionSupport tests various compression options.
func TestCompressionSupport(t *testing.T) {
	t.Run("GzipCompression", func(t *testing.T) {
		// Create source directory
		srcDir := t.TempDir()
		testFile := filepath.Join(srcDir, "test.txt")
		
		if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Create compressed tarball
		tarballPath := filepath.Join(t.TempDir(), "test.tar.gz")
		err := CreateTarball(srcDir, tarballPath, WithCompression(GzipCompression))
		if err != nil {
			t.Fatalf("failed to create compressed tarball: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(tarballPath); err != nil {
			t.Fatalf("compressed tarball not created: %v", err)
		}
	})

	t.Run("NoCompression", func(t *testing.T) {
		srcDir := t.TempDir()
		testFile := filepath.Join(srcDir, "test.txt")
		
		if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		tarballPath := filepath.Join(t.TempDir(), "test.tar")
		err := CreateTarball(srcDir, tarballPath, WithCompression(NoCompression))
		if err != nil {
			t.Fatalf("failed to create uncompressed tarball: %v", err)
		}

		if _, err := os.Stat(tarballPath); err != nil {
			t.Fatalf("uncompressed tarball not created: %v", err)
		}
	})
}

// TestBuilderWithProgress tests build operations with progress reporting.
func TestBuilderWithProgress(t *testing.T) {
	config := &BuildConfig{
		BaseImage: "alpine:latest",
		Platform: &PlatformConfig{
			Architecture: "amd64",
			OS:           "linux",
		},
	}

	registry := &MockRegistryClient{
		images: map[string]v1.Image{
			"alpine:latest": &mockImage{},
		},
	}
	cache := &MockLayerCache{}

	builder, err := NewBuilder(config, registry, cache)
	if err != nil {
		t.Fatalf("failed to create builder: %v", err)
	}

	progressChan := make(chan BuildProgress, 10)
	ctx := context.Background()

	result, err := builder.BuildWithProgress(ctx, progressChan)
	if err != nil {
		t.Fatalf("build with progress failed: %v", err)
	}

	if result.Image == nil {
		t.Error("build result should contain an image")
	}

	// Check that we received some progress updates
	// The channel is already closed by BuildWithProgress
	progressCount := 0
	for range progressChan {
		progressCount++
	}

	if progressCount == 0 {
		t.Error("expected progress updates")
	}
}

// mockImage implements v1.Image for testing.
type mockImage struct{}

func (m *mockImage) Layers() ([]v1.Layer, error) {
	return []v1.Layer{}, nil
}

func (m *mockImage) MediaType() (types.MediaType, error) {
	return types.OCIManifestSchema1, nil
}

func (m *mockImage) Size() (int64, error) {
	return 0, nil
}

func (m *mockImage) ConfigName() (v1.Hash, error) {
	return v1.Hash{}, nil
}

func (m *mockImage) ConfigFile() (*v1.ConfigFile, error) {
	return &v1.ConfigFile{}, nil
}

func (m *mockImage) RawConfigFile() ([]byte, error) {
	return []byte("{}"), nil
}

func (m *mockImage) Digest() (v1.Hash, error) {
	return v1.Hash{Algorithm: "sha256", Hex: "test"}, nil
}

func (m *mockImage) Manifest() (*v1.Manifest, error) {
	return &v1.Manifest{}, nil
}

func (m *mockImage) RawManifest() ([]byte, error) {
	return []byte("{}"), nil
}

func (m *mockImage) LayerByDigest(v1.Hash) (v1.Layer, error) {
	return NewEmptyLayer(), nil
}

func (m *mockImage) LayerByDiffID(v1.Hash) (v1.Layer, error) {
	return NewEmptyLayer(), nil
}