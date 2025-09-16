package gitea

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImageLoader(t *testing.T) {
	loader, err := NewImageLoader()

	// Note: This test might fail if Docker daemon is not available
	// In CI/CD environments, this would be run with Docker daemon
	if err != nil {
		t.Skipf("Docker daemon not available: %v", err)
		return
	}

	assert.NotNil(t, loader)
	assert.NotNil(t, loader.client)

	// Clean up
	err = loader.Close()
	assert.NoError(t, err)
}

func TestCalculateDigest(t *testing.T) {
	loader := &ImageLoader{} // Don't need Docker client for this test

	content := []byte("test content")
	digest := loader.CalculateDigest(content)

	assert.NotEmpty(t, digest.String())
	assert.Contains(t, digest.String(), "sha256:")

	// Test that same content produces same digest
	digest2 := loader.CalculateDigest(content)
	assert.Equal(t, digest.String(), digest2.String())

	// Test that different content produces different digest
	differentContent := []byte("different content")
	differentDigest := loader.CalculateDigest(differentContent)
	assert.NotEqual(t, digest.String(), differentDigest.String())
}

func TestImageManifest_ToJSON(t *testing.T) {
	manifest := &ImageManifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
		Config: ManifestConfig{
			MediaType: "application/vnd.docker.container.image.v1+json",
			Size:      1234,
			Digest:    "sha256:testconfig",
		},
		Layers: []ManifestLayer{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      5678,
				Digest:    "sha256:testlayer",
			},
		},
		TotalSize: 6912,
	}

	jsonBytes, err := manifest.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify it's valid JSON with expected structure
	jsonStr := string(jsonBytes)
	assert.Contains(t, jsonStr, "schemaVersion")
	assert.Contains(t, jsonStr, "mediaType")
	assert.Contains(t, jsonStr, "config")
	assert.Contains(t, jsonStr, "layers")
}

func TestImageManifest_ToReader(t *testing.T) {
	manifest := &ImageManifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
		Config: ManifestConfig{
			MediaType: "application/vnd.docker.container.image.v1+json",
			Size:      1234,
			Digest:    "sha256:testconfig",
		},
		Layers: []ManifestLayer{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      5678,
				Digest:    "sha256:testlayer",
			},
		},
	}

	reader, err := manifest.ToReader()
	require.NoError(t, err)
	assert.NotNil(t, reader)
}

func TestImageManifest_ToOCIManifest(t *testing.T) {
	manifest := &ImageManifest{
		SchemaVersion: 2,
		MediaType:     "application/vnd.docker.distribution.manifest.v2+json",
		Config: ManifestConfig{
			MediaType: "application/vnd.docker.container.image.v1+json",
			Size:      1234,
			Digest:    "sha256:testconfig",
		},
		Layers: []ManifestLayer{
			{
				MediaType: "application/vnd.docker.image.rootfs.diff.tar.gzip",
				Size:      5678,
				Digest:    "sha256:testlayer",
			},
		},
	}

	ociManifest := manifest.ToOCIManifest()
	require.NotNil(t, ociManifest)

	// Note: OCI manifest doesn't include SchemaVersion field
	assert.Equal(t, manifest.MediaType, ociManifest.MediaType)
	assert.Equal(t, manifest.Config.Size, ociManifest.Config.Size)
	assert.Equal(t, len(manifest.Layers), len(ociManifest.Layers))
}
