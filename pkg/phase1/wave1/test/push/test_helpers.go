package push

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// TestFixtures provides a complete test environment for push operations
type TestFixtures struct {
	Registry *MockRegistry       // Mock registry instance
	Client   *http.Client        // Test HTTP client
	TempDir  string              // Temporary test directory
}

// SetupTestFixtures initializes a complete test environment for push testing
func SetupTestFixtures(t *testing.T) *TestFixtures {
	// Create temporary directory for test artifacts
	tempDir, err := os.MkdirTemp("", "push-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create mock registry
	registry := NewMockRegistry()

	// Create HTTP client for testing
	client := &http.Client{
		Transport: http.DefaultTransport,
	}

	fixtures := &TestFixtures{
		Registry: registry,
		Client:   client,
		TempDir:  tempDir,
	}

	// Clean up on test completion
	t.Cleanup(func() {
		CleanupTestFixtures(fixtures)
	})

	return fixtures
}

// CreateTestImage generates a test OCI image with the specified name and tag
func CreateTestImage(name, tag string) v1.Image {
	return &testImage{
		name: name,
		tag:  tag,
		manifest: &v1.Manifest{
			SchemaVersion: 2,
			MediaType:     types.DockerManifestSchema2,
			Config: v1.Descriptor{
				MediaType: types.DockerConfigJSON,
				Size:      1024,
				Digest:    v1.Hash{Algorithm: "sha256", Hex: "abc123"},
			},
			Layers: []v1.Descriptor{
				{
					MediaType: types.DockerLayer,
					Size:      2048,
					Digest:    v1.Hash{Algorithm: "sha256", Hex: "def456"},
				},
			},
		},
		config: &v1.ConfigFile{
			Architecture: "amd64",
			OS:           "linux",
			RootFS: v1.RootFS{
				Type:    "layers",
				DiffIDs: []v1.Hash{{Algorithm: "sha256", Hex: "layer1"}},
			},
		},
	}
}

// CleanupTestFixtures removes all test artifacts and shuts down test services
func CleanupTestFixtures(fixtures *TestFixtures) {
	if fixtures == nil {
		return
	}

	// Cleanup mock registry
	if fixtures.Registry != nil {
		fixtures.Registry.Cleanup()
	}

	// Remove temporary directory
	if fixtures.TempDir != "" {
		os.RemoveAll(fixtures.TempDir)
	}
}

// testImage implements v1.Image interface for testing
type testImage struct {
	name     string
	tag      string
	manifest *v1.Manifest
	config   *v1.ConfigFile
}

// Digest returns the digest of the image
func (img *testImage) Digest() (v1.Hash, error) {
	return v1.Hash{Algorithm: "sha256", Hex: "testimage123"}, nil
}

// Manifest returns the manifest of the image
func (img *testImage) Manifest() (*v1.Manifest, error) {
	return img.manifest, nil
}

// RawManifest returns the raw manifest bytes
func (img *testImage) RawManifest() ([]byte, error) {
	return json.Marshal(img.manifest)
}

// Size returns the size of the image
func (img *testImage) Size() (int64, error) {
	return 3072, nil // config (1024) + layer (2048)
}

// ConfigName returns the digest of the config
func (img *testImage) ConfigName() (v1.Hash, error) {
	return img.manifest.Config.Digest, nil
}

// ConfigFile returns the config file
func (img *testImage) ConfigFile() (*v1.ConfigFile, error) {
	return img.config, nil
}

// RawConfigFile returns the raw config file bytes
func (img *testImage) RawConfigFile() ([]byte, error) {
	return json.Marshal(img.config)
}

// LayerByDigest returns a layer by its digest
func (img *testImage) LayerByDigest(hash v1.Hash) (v1.Layer, error) {
	return &testLayer{digest: hash}, nil
}

// LayerByDiffID returns a layer by its diff ID
func (img *testImage) LayerByDiffID(hash v1.Hash) (v1.Layer, error) {
	return &testLayer{digest: hash}, nil
}

// Layers returns all layers in the image
func (img *testImage) Layers() ([]v1.Layer, error) {
	layers := make([]v1.Layer, len(img.manifest.Layers))
	for i, desc := range img.manifest.Layers {
		layers[i] = &testLayer{digest: desc.Digest}
	}
	return layers, nil
}

// MediaType returns the media type of the image
func (img *testImage) MediaType() (types.MediaType, error) {
	return types.DockerManifestSchema2, nil
}

// testLayer implements v1.Layer interface for testing
type testLayer struct {
	digest v1.Hash
}

// Digest returns the digest of the layer
func (layer *testLayer) Digest() (v1.Hash, error) {
	return layer.digest, nil
}

// DiffID returns the diff ID of the layer
func (layer *testLayer) DiffID() (v1.Hash, error) {
	return layer.digest, nil
}

// Compressed returns the compressed layer content
func (layer *testLayer) Compressed() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte("compressed layer data"))), nil
}

// Uncompressed returns the uncompressed layer content
func (layer *testLayer) Uncompressed() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader([]byte("uncompressed layer data"))), nil
}

// Size returns the size of the layer
func (layer *testLayer) Size() (int64, error) {
	return 2048, nil
}

// MediaType returns the media type of the layer
func (layer *testLayer) MediaType() (types.MediaType, error) {
	return types.DockerLayer, nil
}

// CreateAuthenticatedFixtures creates test fixtures with authentication configured
func CreateAuthenticatedFixtures(t *testing.T, username, password string) *TestFixtures {
	fixtures := SetupTestFixtures(t)

	// Configure authentication on the mock registry
	fixtures.Registry.SetAuth(username, password, true)

	// Create authenticated transport
	authTransport := NewMockAuthTransport(username, password)
	fixtures.Client.Transport = authTransport

	return fixtures
}

// WriteTestFile creates a test file in the fixtures temp directory
func WriteTestFile(fixtures *TestFixtures, filename string, content []byte) (string, error) {
	fullPath := filepath.Join(fixtures.TempDir, filename)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fullPath, nil
}

// ReadTestFile reads a test file from the fixtures temp directory
func ReadTestFile(fixtures *TestFixtures, filename string) ([]byte, error) {
	fullPath := filepath.Join(fixtures.TempDir, filename)
	return os.ReadFile(fullPath)
}