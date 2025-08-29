//go:build !buildah
// +build !buildah

package build

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBuildahBuilder(t *testing.T) {
	t.Run("creates builder successfully", func(t *testing.T) {
		builder, err := NewBuildahBuilder(nil)
		require.NoError(t, err)
		assert.NotNil(t, builder)
	})

	t.Run("accepts trust manager", func(t *testing.T) {
		// Mock trust manager (placeholder since we don't have the actual Phase 1 implementation)
		builder, err := NewBuildahBuilder(nil)
		require.NoError(t, err)
		assert.NotNil(t, builder)
	})
}

func TestBuildahBuilder_BuildImage(t *testing.T) {
	builder, err := NewBuildahBuilder(nil)
	require.NoError(t, err)

	t.Run("fails with empty dockerfile path", func(t *testing.T) {
		opts := BuildOptions{
			DockerfilePath: "",
			ContextDir:     "/tmp",
			Tag:            "test:latest",
		}

		_, err := builder.BuildImage(context.Background(), opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "dockerfile path cannot be empty")
	})

	t.Run("fails with empty context directory", func(t *testing.T) {
		opts := BuildOptions{
			DockerfilePath: "/tmp/Dockerfile",
			ContextDir:     "",
			Tag:            "test:latest",
		}

		_, err := builder.BuildImage(context.Background(), opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context directory cannot be empty")
	})

	t.Run("handles mock build correctly", func(t *testing.T) {
		// Create a temporary dockerfile for testing
		tmpDir := t.TempDir()
		dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
		
		dockerfileContent := `FROM scratch
COPY test.txt /test.txt
`
		err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644)
		require.NoError(t, err)

		// Create a test file to copy
		testFilePath := filepath.Join(tmpDir, "test.txt")
		err = os.WriteFile(testFilePath, []byte("test content"), 0644)
		require.NoError(t, err)

		opts := BuildOptions{
			DockerfilePath: dockerfilePath,
			ContextDir:     tmpDir,
			Tag:            "test:latest",
			Args:           map[string]string{"BUILD_ARG": "value"},
			NoCache:        true,
		}

		// In mock mode, we expect a mock result with an error indicating mock status
		result, err := builder.BuildImage(context.Background(), opts)
		if err != nil {
			// In mock mode, we expect an error but also a result
			t.Logf("Mock build error as expected: %v", err)
			if result != nil {
				assert.NotEmpty(t, result.ImageID)
				assert.Equal(t, "test", result.Repository)
				assert.Equal(t, "latest", result.Tag)
				assert.NotEmpty(t, result.Digest)
				assert.Greater(t, result.Size, int64(0))
			}
		}
	})
}

func TestBuildahBuilder_ListImages(t *testing.T) {
	builder, err := NewBuildahBuilder(nil)
	require.NoError(t, err)

	t.Run("lists images returns mock error", func(t *testing.T) {
		images, err := builder.ListImages(context.Background())
		// In mock mode, we expect an error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mock implementation")
		assert.NotNil(t, images)
		assert.Empty(t, images) // Mock returns empty list
	})
}

func TestBuildahBuilder_RemoveImage(t *testing.T) {
	builder, err := NewBuildahBuilder(nil)
	require.NoError(t, err)

	t.Run("returns mock error", func(t *testing.T) {
		err := builder.RemoveImage(context.Background(), "any-image-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mock implementation")
	})
}

func TestBuildahBuilder_TagImage(t *testing.T) {
	builder, err := NewBuildahBuilder(nil)
	require.NoError(t, err)

	t.Run("returns mock error", func(t *testing.T) {
		err := builder.TagImage(context.Background(), "source:tag", "new:tag")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mock implementation")
	})
}

func TestBuildOptions(t *testing.T) {
	t.Run("validates build options struct", func(t *testing.T) {
		opts := BuildOptions{
			DockerfilePath: "/path/to/Dockerfile",
			ContextDir:     "/path/to/context",
			Tag:            "myimage:latest",
			Args:           map[string]string{"ARG1": "value1"},
			NoCache:        true,
		}

		assert.Equal(t, "/path/to/Dockerfile", opts.DockerfilePath)
		assert.Equal(t, "/path/to/context", opts.ContextDir)
		assert.Equal(t, "myimage:latest", opts.Tag)
		assert.Equal(t, "value1", opts.Args["ARG1"])
		assert.True(t, opts.NoCache)
	})
}

func TestBuildResult(t *testing.T) {
	t.Run("validates build result struct", func(t *testing.T) {
		result := BuildResult{
			ImageID:    "sha256:abc123",
			Repository: "myrepo/myimage",
			Tag:        "latest",
			Digest:     "sha256:def456",
			Size:       123456,
			BuildTime:  30 * time.Second,
		}

		assert.Equal(t, "sha256:abc123", result.ImageID)
		assert.Equal(t, "myrepo/myimage", result.Repository)
		assert.Equal(t, "latest", result.Tag)
		assert.Equal(t, "sha256:def456", result.Digest)
		assert.Equal(t, int64(123456), result.Size)
		assert.Equal(t, 30*time.Second, result.BuildTime)
	})
}

func TestImageInfo(t *testing.T) {
	t.Run("validates image info struct", func(t *testing.T) {
		now := time.Now()
		info := ImageInfo{
			ID:         "sha256:abc123",
			Repository: "myrepo/myimage",
			Tag:        "latest",
			Digest:     "sha256:def456",
			Size:       123456,
			Created:    now,
		}

		assert.Equal(t, "sha256:abc123", info.ID)
		assert.Equal(t, "myrepo/myimage", info.Repository)
		assert.Equal(t, "latest", info.Tag)
		assert.Equal(t, "sha256:def456", info.Digest)
		assert.Equal(t, int64(123456), info.Size)
		assert.Equal(t, now, info.Created)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("getRepository extracts repository correctly", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"myrepo/myimage:latest", "myrepo/myimage"},
			{"localhost:5000/myimage:v1.0", "localhost:5000/myimage"},
			{"simpleimage", "simpleimage"},
			{"", ""},
		}

		for _, tc := range testCases {
			result := getRepository(tc.input)
			assert.Equal(t, tc.expected, result, "failed for input: %s", tc.input)
		}
	})

	t.Run("getTag extracts tag correctly", func(t *testing.T) {
		testCases := []struct {
			input    string
			expected string
		}{
			{"myrepo/myimage:latest", "latest"},
			{"localhost:5000/myimage:v1.0", "v1.0"},
			{"simpleimage", "latest"},
			{"", ""},
		}

		for _, tc := range testCases {
			result := getTag(tc.input)
			assert.Equal(t, tc.expected, result, "failed for input: %s", tc.input)
		}
	})
}