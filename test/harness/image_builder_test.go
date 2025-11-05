package harness

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  ImageConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: 3,
				SizeMB: 10,
				Arch:   "amd64",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			config: ImageConfig{
				Name:   "",
				Tag:    "v1",
				Layers: 3,
				SizeMB: 10,
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "empty tag",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "",
				Layers: 3,
				SizeMB: 10,
			},
			wantErr: true,
			errMsg:  "tag cannot be empty",
		},
		{
			name: "zero layers",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: 0,
				SizeMB: 10,
			},
			wantErr: true,
			errMsg:  "layer count must be positive",
		},
		{
			name: "negative layers",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: -1,
				SizeMB: 10,
			},
			wantErr: true,
			errMsg:  "layer count must be positive",
		},
		{
			name: "negative size",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: 3,
				SizeMB: -5,
			},
			wantErr: true,
			errMsg:  "size cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGenerateDockerfile(t *testing.T) {
	tests := []struct {
		name   string
		config ImageConfig
		checks []string // Strings that should be in the Dockerfile
	}{
		{
			name: "single layer",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: 1,
				SizeMB: 5,
				Arch:   "amd64",
			},
			checks: []string{
				"FROM alpine:latest",
				"COPY layer-0.dat",
				"LABEL test.image.name=\"testapp\"",
				"LABEL test.image.layers=\"1\"",
				"CMD",
			},
		},
		{
			name: "multiple layers",
			config: ImageConfig{
				Name:   "multiapp",
				Tag:    "v2",
				Layers: 5,
				SizeMB: 100,
				Arch:   "arm64",
			},
			checks: []string{
				"FROM alpine:latest",
				"COPY layer-0.dat",
				"COPY layer-1.dat",
				"COPY layer-2.dat",
				"COPY layer-3.dat",
				"COPY layer-4.dat",
				"LABEL test.image.name=\"multiapp\"",
				"LABEL test.image.layers=\"5\"",
				"LABEL test.image.arch=\"arm64\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dockerfile := generateDockerfile(tt.config)

			// Check that all expected strings are present
			for _, check := range tt.checks {
				assert.Contains(t, dockerfile, check,
					"Dockerfile should contain: %s", check)
			}

			// Verify it's not empty
			assert.NotEmpty(t, dockerfile)

			// Verify it starts with FROM
			assert.True(t, strings.HasPrefix(dockerfile, "FROM "))
		})
	}
}

func TestGenerateTestFiles(t *testing.T) {
	tests := []struct {
		name   string
		config ImageConfig
	}{
		{
			name: "small files",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: 2,
				SizeMB: 1, // 1MB total
				Arch:   "amd64",
			},
		},
		{
			name: "zero size files",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: 3,
				SizeMB: 0, // Minimal files
				Arch:   "amd64",
			},
		},
		{
			name: "multiple layers",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: 5,
				SizeMB: 10,
				Arch:   "amd64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			buildDir, err := os.MkdirTemp("", "test-build-")
			require.NoError(t, err)
			defer os.RemoveAll(buildDir)

			// Generate test files
			err = generateTestFiles(buildDir, tt.config)
			require.NoError(t, err)

			// Verify files were created
			for i := 0; i < tt.config.Layers; i++ {
				filename := filepath.Join(buildDir, "layer-"+string(rune('0'+i))+".dat")
				info, err := os.Stat(filename)
				require.NoError(t, err, "File layer-%d.dat should exist", i)

				// If size specified, verify approximate size
				if tt.config.SizeMB > 0 {
					expectedSize := int64(tt.config.SizeMB) * 1024 * 1024 / int64(tt.config.Layers)
					assert.Equal(t, expectedSize, info.Size(),
						"File size should match expected size per layer")
				} else {
					// Minimal files should exist and not be empty
					assert.Greater(t, info.Size(), int64(0))
				}
			}
		})
	}
}

func TestCreateRandomFile(t *testing.T) {
	tests := []struct {
		name string
		size int64
	}{
		{"small file", 1024},       // 1KB
		{"medium file", 1024 * 100}, // 100KB
		{"zero size", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "test-random-")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			filename := filepath.Join(tmpDir, "random.dat")

			// Create random file
			err = createRandomFile(filename, tt.size)
			require.NoError(t, err)

			// Verify file exists and has correct size
			info, err := os.Stat(filename)
			require.NoError(t, err)
			assert.Equal(t, tt.size, info.Size())

			// If size > 0, verify file has non-zero content
			if tt.size > 0 {
				data, err := os.ReadFile(filename)
				require.NoError(t, err)
				assert.Len(t, data, int(tt.size))

				// Check that it's not all zeros (random data)
				hasNonZero := false
				for _, b := range data {
					if b != 0 {
						hasNonZero = true
						break
					}
				}
				assert.True(t, hasNonZero, "Random file should contain non-zero bytes")
			}
		})
	}
}

func TestCleanupBuildDir(t *testing.T) {
	// Create temporary directory
	buildDir, err := os.MkdirTemp("", "test-cleanup-")
	require.NoError(t, err)

	// Create some files in it
	testFile := filepath.Join(buildDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	// Verify directory exists
	_, err = os.Stat(buildDir)
	require.NoError(t, err)

	// Cleanup
	err = cleanupBuildDir(buildDir)
	require.NoError(t, err)

	// Verify directory is gone
	_, err = os.Stat(buildDir)
	assert.True(t, os.IsNotExist(err), "Build directory should be removed")
}

// Integration test - requires Docker daemon
func TestBuildTestImage_Integration(t *testing.T) {
	// Skip if Docker not available
	client, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		t.Skip("Docker client not available, skipping integration test")
	}
	defer client.Close()

	// Ping Docker to ensure it's running
	ctx := context.Background()
	_, err = client.Ping(ctx)
	if err != nil {
		t.Skip("Docker daemon not running, skipping integration test")
	}

	env := &BuilderTestEnvironment{
		DockerClient: client,
	}

	tests := []struct {
		name   string
		config ImageConfig
	}{
		{
			name: "small image",
			config: ImageConfig{
				Name:   "test-small",
				Tag:    "latest",
				Layers: 2,
				SizeMB: 1,
				Arch:   "amd64",
			},
		},
		{
			name: "medium image",
			config: ImageConfig{
				Name:   "test-medium",
				Tag:    "v1",
				Layers: 5,
				SizeMB: 10,
				Arch:   "amd64",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build image
			result, err := env.BuildTestImage(ctx, tt.config)
			require.NoError(t, err)
			require.NotNil(t, result)

			// Verify result
			expectedRef := tt.config.Name + ":" + tt.config.Tag
			assert.Equal(t, expectedRef, result.ImageRef)
			assert.NotEmpty(t, result.ImageID)
			assert.Greater(t, result.LayerCount, 0)
			assert.Greater(t, result.SizeBytes, int64(0))

			// Cleanup: remove the test image
			_, err = client.ImageRemove(ctx, result.ImageID, types.ImageRemoveOptions{
				Force: true,
			})
			assert.NoError(t, err)
		})
	}
}

func TestBuildTestImage_InvalidConfig(t *testing.T) {
	// Create mock environment (don't need real Docker for validation tests)
	env := &BuilderTestEnvironment{
		DockerClient: nil, // Will fail before using client
	}

	ctx := context.Background()

	tests := []struct {
		name    string
		config  ImageConfig
		errMsg  string
	}{
		{
			name: "empty name",
			config: ImageConfig{
				Name:   "",
				Tag:    "v1",
				Layers: 3,
			},
			errMsg: "name cannot be empty",
		},
		{
			name: "empty tag",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "",
				Layers: 3,
			},
			errMsg: "tag cannot be empty",
		},
		{
			name: "invalid layer count",
			config: ImageConfig{
				Name:   "testapp",
				Tag:    "v1",
				Layers: -1,
			},
			errMsg: "layer count must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := env.BuildTestImage(ctx, tt.config)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}
