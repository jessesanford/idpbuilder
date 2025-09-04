package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewBuilder tests the builder constructor
func TestNewBuilder(t *testing.T) {
	tests := []struct {
		name    string
		opts    BuildOptions
		wantErr bool
	}{
		{
			name: "default options",
			opts: BuildOptions{},
		},
		{
			name: "with platform",
			opts: BuildOptions{
				Platform: v1.Platform{
					OS:           "linux",
					Architecture: "amd64",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder, err := NewBuilder(tt.opts)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, builder)
			}
		})
	}
}

// TestBuild tests the core Build functionality
func TestBuild(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "builder-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test content
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	ctx := context.Background()
	
	t.Run("valid directory", func(t *testing.T) {
		builder, err := NewBuilder(DefaultBuildOptions())
		require.NoError(t, err)
		
		image, err := builder.Build(ctx, tempDir, DefaultBuildOptions())
		assert.NoError(t, err)
		assert.NotNil(t, image)
	})
	
	t.Run("nonexistent directory", func(t *testing.T) {
		builder, err := NewBuilder(DefaultBuildOptions())
		require.NoError(t, err)
		
		_, err = builder.Build(ctx, "/nonexistent", DefaultBuildOptions())
		assert.Error(t, err)
	})
}

// TestBuildTarball tests tarball export functionality
func TestBuildTarball(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "builder-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test content
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	ctx := context.Background()
	opts := DefaultBuildOptions()
	
	// Create output directory
	outputDir, err := os.MkdirTemp("", "output-*")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)
	
	outputPath := filepath.Join(outputDir, "test.tar")
	
	builder, err := NewBuilder(opts)
	require.NoError(t, err)
	
	// Test successful tarball creation
	err = builder.BuildTarball(ctx, tempDir, outputPath, opts)
	assert.NoError(t, err)
	
	// Verify output file exists
	info, err := os.Stat(outputPath)
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

// TestDefaultBuildOptions tests the default options
func TestDefaultBuildOptions(t *testing.T) {
	opts := DefaultBuildOptions()
	
	assert.Equal(t, "linux", opts.Platform.OS)
	assert.Equal(t, "amd64", opts.Platform.Architecture)
	assert.NotNil(t, opts.Labels)
	assert.NotNil(t, opts.FeatureFlags)
}

// TestLayerFactory tests layer creation functionality  
func TestLayerFactory(t *testing.T) {
	factory := NewLayerFactory()
	assert.NotNil(t, factory)
	
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "layer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test content
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	layer, err := factory.CreateLayer(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, layer)
}

// TestConfigFactory tests configuration generation
func TestConfigFactory(t *testing.T) {
	platform := v1.Platform{
		OS:           "linux",
		Architecture: "amd64",
	}
	
	factory := NewConfigFactory(platform)
	assert.NotNil(t, factory)
	
	opts := BuildOptions{
		Platform: platform,
		Env:      []string{"ENV_VAR=value"},
		Cmd:      []string{"echo", "hello"},
	}
	
	config, err := factory.GenerateConfig(opts)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	
	assert.Equal(t, "linux", config.OS)
	assert.Equal(t, "amd64", config.Architecture)
}

// TestTarballWriter tests tarball export functionality
func TestTarballWriter(t *testing.T) {
	writer := NewTarballWriter()
	assert.NotNil(t, writer)
	
	// Create a simple test image
	tempDir, err := os.MkdirTemp("", "tarball-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)
	
	// Build an image for testing
	ctx := context.Background()
	builder, err := NewBuilder(DefaultBuildOptions())
	require.NoError(t, err)
	
	image, err := builder.Build(ctx, tempDir, DefaultBuildOptions())
	require.NoError(t, err)
	
	// Test tarball writing
	outputPath := filepath.Join(tempDir, "test.tar")
	
	err = writer.Write(image, outputPath, "localhost/test:latest")
	assert.NoError(t, err)
	
	// Verify file was created
	info, err := os.Stat(outputPath)
	assert.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}