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

// TestDefaultBuildOptions tests the default options function
func TestDefaultBuildOptions(t *testing.T) {
	opts := DefaultBuildOptions()
	
	assert.Equal(t, "linux", opts.Platform.OS)
	assert.Equal(t, "amd64", opts.Platform.Architecture)
	assert.NotNil(t, opts.Labels)
	assert.NotNil(t, opts.FeatureFlags)
	assert.Contains(t, opts.Labels, "org.opencontainers.image.source")
	assert.Equal(t, "idpbuilder", opts.Labels["org.opencontainers.image.source"])
}

// TestBuildOptionsBuilders tests the fluent builder methods
func TestBuildOptionsBuilders(t *testing.T) {
	opts := DefaultBuildOptions()
	
	// Test WithPlatform
	opts = opts.WithPlatform("linux", "arm64")
	assert.Equal(t, "linux", opts.Platform.OS)
	assert.Equal(t, "arm64", opts.Platform.Architecture)
	
	// Test WithBaseImage
	opts = opts.WithBaseImage("alpine:latest")
	assert.Equal(t, "alpine:latest", opts.BaseImage)
	
	// Test WithLabels
	opts = opts.WithLabels(map[string]string{
		"test-label": "test-value",
	})
	assert.Equal(t, "test-value", opts.Labels["test-label"])
	
	// Test WithWorkingDir
	opts = opts.WithWorkingDir("/app")
	assert.Equal(t, "/app", opts.WorkingDir)
	
	// Test WithEntrypoint
	opts = opts.WithEntrypoint("/bin/sh", "-c")
	assert.Equal(t, []string{"/bin/sh", "-c"}, opts.Entrypoint)
	
	// Test WithCmd
	opts = opts.WithCmd("echo", "hello")
	assert.Equal(t, []string{"echo", "hello"}, opts.Cmd)
	
	// Test WithEnv
	opts = opts.WithEnv("ENV_VAR=value")
	assert.Contains(t, opts.Env, "ENV_VAR=value")
	
	// Test WithFeatureFlags
	opts = opts.WithFeatureFlags(map[string]bool{
		"test-flag": true,
	})
	assert.True(t, opts.FeatureFlags["test-flag"])
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
	
	// Create subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)
	
	subFile := filepath.Join(subDir, "sub.txt")
	err = os.WriteFile(subFile, []byte("sub content"), 0644)
	require.NoError(t, err)
	
	t.Run("create layer from directory", func(t *testing.T) {
		layer, err := factory.CreateLayer(tempDir)
		assert.NoError(t, err)
		assert.NotNil(t, layer)
		
		// Verify layer properties
		size, err := layer.Size()
		assert.NoError(t, err)
		assert.Greater(t, size, int64(0))
		
		diffID, err := layer.DiffID()
		assert.NoError(t, err)
		assert.NotEmpty(t, diffID.String())
		
		// Get layer info
		info, err := GetLayerInfo(layer)
		assert.NoError(t, err)
		assert.Equal(t, size, info.Size)
		assert.Equal(t, diffID, info.DiffID)
	})
	
	t.Run("empty context dir", func(t *testing.T) {
		_, err := factory.CreateLayer("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context directory cannot be empty")
	})
	
	t.Run("nonexistent directory", func(t *testing.T) {
		_, err := factory.CreateLayer("/nonexistent")
		assert.Error(t, err)
	})
}

// TestConfigFactory tests configuration generation
func TestConfigFactory(t *testing.T) {
	platform := v1.Platform{
		OS:           "linux",
		Architecture: "amd64",
	}
	
	factory := NewConfigFactory(platform)
	assert.NotNil(t, factory)
	
	t.Run("generate basic config", func(t *testing.T) {
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
		assert.Contains(t, config.Config.Env, "ENV_VAR=value")
		assert.Equal(t, []string{"echo", "hello"}, config.Config.Cmd)
	})
	
	t.Run("config with labels", func(t *testing.T) {
		opts := BuildOptions{
			Platform: platform,
			Labels: map[string]string{
				"org.opencontainers.image.created": "", // Should be set to build time
				"test-label":                          "test-value",
			},
		}
		
		config, err := factory.GenerateConfig(opts)
		assert.NoError(t, err)
		assert.NotEmpty(t, config.Config.Labels["org.opencontainers.image.created"])
		assert.Equal(t, "test-value", config.Config.Labels["test-label"])
	})
	
	t.Run("validation errors", func(t *testing.T) {
		// Test invalid working directory
		opts := BuildOptions{
			Platform:   platform,
			WorkingDir: "relative/path", // Should be absolute
		}
		
		_, err := factory.GenerateConfig(opts)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "working directory must be an absolute path")
	})
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
	t.Run("write tarball", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "test.tar")
		
		err := writer.Write(image, outputPath, "localhost/test:latest")
		assert.NoError(t, err)
		
		// Verify file was created
		info, err := os.Stat(outputPath)
		assert.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))
	})
	
	t.Run("invalid parameters", func(t *testing.T) {
		outputPath := filepath.Join(tempDir, "test.tar")
		
		// Nil image
		err := writer.Write(nil, outputPath, "test:latest")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "image cannot be nil")
		
		// Empty output path
		err = writer.Write(image, "", "test:latest")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "output path cannot be empty")
		
		// Empty reference
		err = writer.Write(image, outputPath, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "image reference cannot be empty")
	})
}

// TestConfigFactoryValidation tests more validation scenarios
func TestConfigFactoryValidation(t *testing.T) {
	platform := v1.Platform{OS: "linux", Architecture: "amd64"}
	factory := NewConfigFactory(platform)
	
	tests := []struct {
		name    string
		opts    BuildOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid port format",
			opts: BuildOptions{
				Platform: platform,
				ExposedPorts: map[string]struct{}{
					"invalid": {},
				},
			},
			wantErr: true,
			errMsg:  "invalid exposed port format",
		},
		{
			name: "invalid environment variable",
			opts: BuildOptions{
				Platform: platform,
				Env:      []string{"INVALID_ENV_VAR"},
			},
			wantErr: true,
			errMsg:  "environment variable must be in KEY=value format",
		},
		{
			name: "invalid user format",
			opts: BuildOptions{
				Platform: platform,
				User:     "user:group:extra",
			},
			wantErr: true,
			errMsg:  "user specification can have at most one colon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := factory.GenerateConfig(tt.opts)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestConfigFactoryHelpers tests helper functions
func TestConfigFactoryHelpers(t *testing.T) {
	t.Run("DefaultLabels", func(t *testing.T) {
		labels := DefaultLabels("test-source")
		assert.Contains(t, labels, "org.opencontainers.image.source")
		assert.Equal(t, "test-source", labels["org.opencontainers.image.source"])
	})
	
	t.Run("MergeConfigs", func(t *testing.T) {
		config1 := &v1.ConfigFile{
			Architecture: "amd64",
			OS:           "linux",
			Config: v1.Config{
				Env: []string{"VAR1=value1"},
				Labels: map[string]string{
					"label1": "value1",
				},
			},
		}
		
		config2 := &v1.ConfigFile{
			Architecture: "arm64", // Override
			Config: v1.Config{
				Env: []string{"VAR2=value2"},
				Labels: map[string]string{
					"label2": "value2",
				},
			},
		}
		
		merged := MergeConfigs(config1, config2)
		assert.Equal(t, "arm64", merged.Architecture) // Should be overridden
		assert.Equal(t, "linux", merged.OS)           // Should be preserved
		assert.Contains(t, merged.Config.Env, "VAR1=value1")
		assert.Contains(t, merged.Config.Env, "VAR2=value2")
		assert.Equal(t, "value1", merged.Config.Labels["label1"])
		assert.Equal(t, "value2", merged.Config.Labels["label2"])
	})
}

// TestLayerFactoryWithOptions tests layer factory configuration options
func TestLayerFactoryWithOptions(t *testing.T) {
	factory := NewLayerFactory().
		WithPermissionPreservation(false).
		WithTimestampPreservation(true)
	
	assert.NotNil(t, factory)
	assert.False(t, factory.preservePermissions)
	assert.True(t, factory.preserveTimestamps)
}

// TestTarballWriterWithOptions tests tarball writer configuration
func TestTarballWriterWithOptions(t *testing.T) {
	t.Run("with options", func(t *testing.T) {
		opts := TarballOptions{
			Compress:        true,
			IncludeManifest: false,
		}
		writer := NewTarballWriterWithOptions(opts)
		assert.NotNil(t, writer)
		assert.True(t, writer.options.Compress)
		assert.False(t, writer.options.IncludeManifest)
	})
	
	t.Run("with compression", func(t *testing.T) {
		writer := NewTarballWriter().WithCompression(true)
		assert.True(t, writer.options.Compress)
	})
}

// TestTarballWriterMultiple tests multi-image tarball export
func TestTarballWriterMultiple(t *testing.T) {
	writer := NewTarballWriter()
	
	t.Run("empty images map", func(t *testing.T) {
		outputPath := filepath.Join(os.TempDir(), "empty.tar")
		err := writer.WriteMultiple(map[string]v1.Image{}, outputPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no images provided for export")
	})
	
	t.Run("invalid reference", func(t *testing.T) {
		outputPath := filepath.Join(os.TempDir(), "invalid.tar")
		images := map[string]v1.Image{
			"invalid::reference": nil,
		}
		err := writer.WriteMultiple(images, outputPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid reference")
	})
}

// TestTarballInfo tests tarball information functions
func TestTarballInfo(t *testing.T) {
	t.Run("nonexistent file", func(t *testing.T) {
		_, err := GetTarballInfo("/nonexistent/file.tar")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tarball file not found")
	})
	
	t.Run("valid file", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "test*.tar")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		
		// Write some content
		_, err = tempFile.WriteString("test content")
		require.NoError(t, err)
		tempFile.Close()
		
		info, err := GetTarballInfo(tempFile.Name())
		assert.NoError(t, err)
		assert.Equal(t, tempFile.Name(), info.Path)
		assert.Greater(t, info.Size, int64(0))
	})
}

// TestValidateTarball tests tarball validation
func TestValidateTarball(t *testing.T) {
	t.Run("nonexistent file", func(t *testing.T) {
		err := ValidateTarball("/nonexistent/file.tar")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot open tarball file")
	})
	
	t.Run("empty file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "empty*.tar")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()
		
		err = ValidateTarball(tempFile.Name())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tarball file is empty")
	})
	
	t.Run("valid file", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "valid*.tar")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		
		_, err = tempFile.WriteString("valid content")
		require.NoError(t, err)
		tempFile.Close()
		
		err = ValidateTarball(tempFile.Name())
		assert.NoError(t, err)
	})
}

// TestCompressTarball tests tarball compression
func TestCompressTarball(t *testing.T) {
	err := CompressTarball("input.tar", "output.tar.gz")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tarball compression not yet implemented")
}

// Benchmark tests for performance validation
func BenchmarkBuild(b *testing.B) {
	// Create test directory
	tempDir, err := os.MkdirTemp("", "bench-*")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)
	
	// Create some test files
	for i := 0; i < 10; i++ {
		err := os.WriteFile(filepath.Join(tempDir, "file"+string(rune(i))+".txt"), 
			[]byte("content "+string(rune(i))), 0644)
		require.NoError(b, err)
	}
	
	builder, err := NewBuilder(DefaultBuildOptions())
	require.NoError(b, err)
	
	ctx := context.Background()
	opts := DefaultBuildOptions()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := builder.Build(ctx, tempDir, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}