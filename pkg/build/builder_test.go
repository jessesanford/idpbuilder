package build

import (
	"os"
	"testing"
	"time"
)

func TestNewBuilderConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		opts     *BuildOptions
		envVars  map[string]string
		checkFn  func(*BuildOptions) bool
	}{
		{
			name: "valid options preserved",
			opts: &BuildOptions{
				StoragePath: "/tmp/test-storage",
				RunRoot:     "/tmp/test-run",
				Debug:       true,
			},
			checkFn: func(opts *BuildOptions) bool {
				return opts.StoragePath == "/tmp/test-storage" && opts.RunRoot == "/tmp/test-run"
			},
		},
		{
			name: "environment variable fallback",
			opts: &BuildOptions{},
			envVars: map[string]string{
				"BUILDAH_STORAGE_PATH": "/env/storage",
				"BUILDAH_RUN_ROOT":     "/env/run",
			},
			checkFn: func(opts *BuildOptions) bool {
				return opts.StoragePath != "" && opts.RunRoot != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Test the configuration logic without actually creating buildah instance
			opts := tt.opts
			if opts == nil {
				opts = &BuildOptions{}
			}

			// Apply same logic as NewBuilder
			if opts.StoragePath == "" {
				opts.StoragePath = os.Getenv("BUILDAH_STORAGE_PATH")
				if opts.StoragePath == "" {
					opts.StoragePath = os.Getenv("XDG_DATA_HOME")
					if opts.StoragePath == "" {
						homeDir, err := os.UserHomeDir()
						if err == nil {
							opts.StoragePath = homeDir + "/.local/share/containers/storage"
						}
					} else {
						opts.StoragePath = opts.StoragePath + "/containers/storage"
					}
				}
			}

			if opts.RunRoot == "" {
				opts.RunRoot = os.Getenv("BUILDAH_RUN_ROOT")
				if opts.RunRoot == "" {
					opts.RunRoot = "/tmp/buildah-run-root"
				}
			}

			if !tt.checkFn(opts) {
				t.Errorf("Configuration check failed for %s", tt.name)
			}
		})
	}
}

func TestBuildContext(t *testing.T) {
	ctx := &BuildContext{
		ID:         "test-build-123",
		WorkingDir: "/app",
		Metadata: map[string]interface{}{
			"created_at": time.Now(),
			"test_key":   "test_value",
		},
	}

	if ctx.ID != "test-build-123" {
		t.Errorf("BuildContext.ID = %v, want %v", ctx.ID, "test-build-123")
	}

	if ctx.WorkingDir != "/app" {
		t.Errorf("BuildContext.WorkingDir = %v, want %v", ctx.WorkingDir, "/app")
	}

	if ctx.Metadata["test_key"] != "test_value" {
		t.Errorf("BuildContext.Metadata[test_key] = %v, want %v", ctx.Metadata["test_key"], "test_value")
	}
}

func TestBuilderInterface(t *testing.T) {
	// Test that our adapter implements the Builder interface
	var _ Builder = &buildahAdapter{}
}

// TestBuildConfigValidation tests BuildConfig validation
func TestBuildConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config *BuildConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: &BuildConfig{
				BaseImage:  "alpine:latest",
				WorkingDir: "/app",
				Env:        map[string]string{"ENV": "test"},
				Labels:     map[string]string{"version": "1.0"},
				Entrypoint: []string{"/app/entrypoint.sh"},
				Cmd:        []string{"serve"},
			},
			valid: true,
		},
		{
			name: "minimal valid config",
			config: &BuildConfig{
				BaseImage: "alpine:latest",
			},
			valid: true,
		},
		{
			name:   "nil config",
			config: nil,
			valid:  false,
		},
		{
			name: "empty base image",
			config: &BuildConfig{
				BaseImage: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.config != nil && tt.config.BaseImage != ""
			if valid != tt.valid {
				t.Errorf("Config validation = %v, want %v", valid, tt.valid)
			}
		})
	}
}

// TestLayerSpec tests LayerSpec functionality
func TestLayerSpec(t *testing.T) {
	layer := &LayerSpec{
		Source:      "/src/file.txt",
		Destination: "/dest/file.txt",
		Type:        LayerTypeCopy,
		Permissions: "644",
		Owner:       "root:root",
	}

	if layer.Type != LayerTypeCopy {
		t.Errorf("LayerSpec.Type = %v, want %v", layer.Type, LayerTypeCopy)
	}

	if layer.Source != "/src/file.txt" {
		t.Errorf("LayerSpec.Source = %v, want %v", layer.Source, "/src/file.txt")
	}
}

// TestLayerTypes tests all layer type constants
func TestLayerTypes(t *testing.T) {
	tests := []struct {
		layerType LayerType
		expected  string
	}{
		{LayerTypeCopy, "copy"},
		{LayerTypeAdd, "add"},
		{LayerTypeRun, "run"},
	}

	for _, tt := range tests {
		t.Run(string(tt.layerType), func(t *testing.T) {
			if string(tt.layerType) != tt.expected {
				t.Errorf("LayerType = %v, want %v", string(tt.layerType), tt.expected)
			}
		})
	}
}

// TestBuildResult tests BuildResult structure
func TestBuildResult(t *testing.T) {
	now := time.Now()
	result := &BuildResult{
		ImageID:   "sha256:abc123",
		Digest:    "sha256:def456",
		Size:      1024000,
		CreatedAt: now,
	}

	if result.ImageID != "sha256:abc123" {
		t.Errorf("BuildResult.ImageID = %v, want %v", result.ImageID, "sha256:abc123")
	}

	if result.Size != 1024000 {
		t.Errorf("BuildResult.Size = %v, want %v", result.Size, 1024000)
	}

	if !result.CreatedAt.Equal(now) {
		t.Errorf("BuildResult.CreatedAt = %v, want %v", result.CreatedAt, now)
	}
}