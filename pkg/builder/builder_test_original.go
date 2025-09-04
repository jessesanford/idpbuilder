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
		{
			name: "with feature flags",
			opts: BuildOptions{
				FeatureFlags: map[string]bool{
					"test-flag": true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder, err := NewBuilder(tt.opts)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, builder)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, builder)
				
				// Verify platform defaults
				if tt.opts.Platform.OS == "" {
					assert.Equal(t, "linux", builder.platform.OS)
				}
				if tt.opts.Platform.Architecture == "" {
					assert.Equal(t, "amd64", builder.platform.Architecture)
				}
				
				// Verify factories are initialized
				assert.NotNil(t, builder.layerFactory)
				assert.NotNil(t, builder.configFactory)
				assert.NotNil(t, builder.tarballWriter)
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
	
	tests := []struct {
		name       string
		contextDir string
		opts       BuildOptions
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "valid directory",
			contextDir: tempDir,
			opts:       DefaultBuildOptions(),
		},
		{
			name:       "nonexistent directory",
			contextDir: "/nonexistent",
			opts:       DefaultBuildOptions(),
			wantErr:    true,
			errMsg:     "context directory not found",
		},
		{
			name:       "file instead of directory",
			contextDir: testFile,
			opts:       DefaultBuildOptions(),
			wantErr:    true,
			errMsg:     "context path is not a directory",
		},
		{
			name:       "empty context dir",
			contextDir: "",
			opts:       DefaultBuildOptions(),
			wantErr:    true,
		},
		{
			name:       "feature flag restriction",
			contextDir: tempDir,
			opts: BuildOptions{
				Platform: v1.Platform{OS: "linux", Architecture: "amd64"},
				FeatureFlags: map[string]bool{
					"multi-stage-build": true,
				},
			},
			wantErr: true,
			errMsg:  "multi-stage builds not yet implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder, err := NewBuilder(tt.opts)
			require.NoError(t, err)
			
			image, err := builder.Build(ctx, tt.contextDir, tt.opts)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, image)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, image)
				
				// Verify image properties
				manifest, err := image.Manifest()
				assert.NoError(t, err)
				assert.NotNil(t, manifest)
				
				// Should have at least one layer (our content)
				assert.Greater(t, len(manifest.Layers), 0)
				
				config, err := image.ConfigFile()
				assert.NoError(t, err)
				assert.NotNil(t, config)
				assert.Equal(t, tt.opts.Platform.OS, config.OS)
				assert.Equal(t, tt.opts.Platform.Architecture, config.Architecture)
			}
		})
	}
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
	t.Run("successful tarball", func(t *testing.T) {
		err = builder.BuildTarball(ctx, tempDir, outputPath, opts)
		assert.NoError(t, err)
		
		// Verify output file exists
		info, err := os.Stat(outputPath)
		assert.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))
	})
	
	// Test error cases
	t.Run("build failure", func(t *testing.T) {
		err = builder.BuildTarball(ctx, "/nonexistent", outputPath, opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to build image")
	})
}