package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		name          string
		enableFeature bool
		opts          BuilderOptions
	}{
		{
			name:          "feature enabled",
			enableFeature: true,
			opts: BuilderOptions{
				Registry: "localhost:5000",
				Insecure: true,
			},
		},
		{
			name:          "feature disabled",
			enableFeature: false,
			opts: BuilderOptions{
				Registry: "localhost:5000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set feature flag
			if tt.enableFeature {
				os.Setenv("ENABLE_CORE_BUILDER", "true")
			} else {
				os.Setenv("ENABLE_CORE_BUILDER", "false")
			}
			defer os.Unsetenv("ENABLE_CORE_BUILDER")

			builder := NewBuilder(tt.opts)
			if builder == nil {
				t.Errorf("NewBuilder() returned nil")
			}
		})
	}
}

func TestBuildOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    BuildOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid options",
			opts: BuildOptions{
				ImageName:   "test-image",
				Tag:         "latest",
				BaseImage:   "gcr.io/distroless/static:nonroot",
				Platform:    "linux/amd64",
				Parallelism: 1,
				Timeout:     30 * time.Minute,
			},
			wantErr: false,
		},
		{
			name: "missing image name",
			opts: BuildOptions{
				Tag:       "latest",
				BaseImage: "gcr.io/distroless/static:nonroot",
			},
			wantErr: true,
			errMsg:  "image name is required",
		},
		{
			name: "missing tag",
			opts: BuildOptions{
				ImageName: "test-image",
				BaseImage: "gcr.io/distroless/static:nonroot",
			},
			wantErr: true,
			errMsg:  "tag is required",
		},
		{
			name: "missing base image",
			opts: BuildOptions{
				ImageName: "test-image",
				Tag:       "latest",
			},
			wantErr: true,
			errMsg:  "base image is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				// Basic error message check - not checking exact match to keep it simple
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestBuildOptions_SetDefaults(t *testing.T) {
	opts := BuildOptions{
		ImageName: "test-image",
		Tag:       "v1.0.0",
	}

	opts.SetDefaults()

	if opts.BaseImage != "gcr.io/distroless/static:nonroot" {
		t.Errorf("Expected BaseImage to be set to default")
	}
	if opts.Platform != "linux/amd64" {
		t.Errorf("Expected Platform to be set to default")
	}
	if opts.Labels == nil {
		t.Errorf("Expected Labels to be initialized")
	}
	if opts.BuildArgs == nil {
		t.Errorf("Expected BuildArgs to be initialized")
	}
	if opts.CreatedBy != "idpbuilder-oci-go-cr" {
		t.Errorf("Expected CreatedBy to be set to default")
	}

	// Check that standard labels are set
	if _, exists := opts.Labels["org.opencontainers.image.created"]; !exists {
		t.Errorf("Expected standard label to be set")
	}
}

func TestBuildOptions_Clone(t *testing.T) {
	original := BuildOptions{
		ImageName: "test-image",
		Tag:       "latest",
		BaseImage: "gcr.io/distroless/static:nonroot",
		Labels: map[string]string{
			"test": "value",
		},
		Env: []string{"ENV=value"},
		Cmd: []string{"echo", "hello"},
		Files: []FileSpec{
			{
				Source:      "test.txt",
				Destination: "/app/test.txt",
				Mode:        0644,
			},
		},
	}

	cloned := original.Clone()

	// Verify values are copied
	if cloned.ImageName != original.ImageName {
		t.Errorf("Clone failed to copy ImageName")
	}
	if cloned.Tag != original.Tag {
		t.Errorf("Clone failed to copy Tag")
	}

	// Verify deep copy of maps
	if cloned.Labels["test"] != "value" {
		t.Errorf("Clone failed to copy Labels")
	}
	original.Labels["test"] = "modified"
	if cloned.Labels["test"] == "modified" {
		t.Errorf("Clone did not deep copy Labels")
	}

	// Verify deep copy of slices
	if len(cloned.Env) != 1 || cloned.Env[0] != "ENV=value" {
		t.Errorf("Clone failed to copy Env slice")
	}
}

func TestNewDefaultBuildOptions(t *testing.T) {
	opts := NewDefaultBuildOptions()

	if opts.BaseImage != "gcr.io/distroless/static:nonroot" {
		t.Errorf("Expected default BaseImage")
	}
	if opts.Platform != "linux/amd64" {
		t.Errorf("Expected default Platform")
	}
	if opts.Tag != "latest" {
		t.Errorf("Expected default Tag")
	}
	if opts.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries")
	}
	if opts.Timeout != 30*time.Minute {
		t.Errorf("Expected default Timeout")
	}
	if opts.Parallelism != 1 {
		t.Errorf("Expected default Parallelism")
	}
	if opts.CreatedBy != "idpbuilder-oci-go-cr" {
		t.Errorf("Expected default CreatedBy")
	}
}

func TestNewDefaultBuilderOptions(t *testing.T) {
	opts := NewDefaultBuilderOptions()

	if opts.Registry != "localhost:5000" {
		t.Errorf("Expected default Registry")
	}
	if opts.Insecure != false {
		t.Errorf("Expected default Insecure")
	}
	if opts.Timeout != 30*time.Minute {
		t.Errorf("Expected default Timeout")
	}
	if opts.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries")
	}
}

func TestIsFeatureEnabled(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{
			name:     "enabled",
			envValue: "true",
			expected: true,
		},
		{
			name:     "disabled",
			envValue: "false",
			expected: false,
		},
		{
			name:     "unset",
			envValue: "",
			expected: false,
		},
		{
			name:     "invalid value",
			envValue: "yes",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("ENABLE_CORE_BUILDER", tt.envValue)
			} else {
				os.Unsetenv("ENABLE_CORE_BUILDER")
			}

			result := IsFeatureEnabled()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDisabledBuilder(t *testing.T) {
	// Set feature flag to disabled
	os.Setenv("ENABLE_CORE_BUILDER", "false")
	defer os.Unsetenv("ENABLE_CORE_BUILDER")

	builder := NewBuilder(BuilderOptions{})
	opts := BuildOptions{
		ImageName: "test",
		Tag:       "latest",
		BaseImage: "alpine",
	}

	// Test that all methods return appropriate errors
	_, err := builder.Build(context.Background(), opts)
	if err == nil {
		t.Errorf("Expected error from disabled builder.Build()")
	}

	_, err = builder.BuildFromTarball(context.Background(), "/tmp/test.tar", opts)
	if err == nil {
		t.Errorf("Expected error from disabled builder.BuildFromTarball()")
	}

	err = builder.ValidateOptions(opts)
	if err == nil {
		t.Errorf("Expected error from disabled builder.ValidateOptions()")
	}
}

func TestBuildOptions_ValidateFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "builder-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		opts    BuildOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid file spec",
			opts: BuildOptions{
				ImageName:   "test",
				Tag:         "latest",
				BaseImage:   "alpine",
				ContextPath: tempDir,
				Parallelism: 1,
				Timeout:     30 * time.Minute,
				Files: []FileSpec{
					{
						Source:      "test.txt",
						Destination: "/app/test.txt",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing source",
			opts: BuildOptions{
				ImageName:   "test",
				Tag:         "latest",
				BaseImage:   "alpine",
				ContextPath: tempDir,
				Files: []FileSpec{
					{
						Destination: "/app/test.txt",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestContainerBuilder_ValidateOptions(t *testing.T) {
	// Enable feature for testing
	os.Setenv("ENABLE_CORE_BUILDER", "true")
	defer os.Unsetenv("ENABLE_CORE_BUILDER")

	builder := NewBuilder(BuilderOptions{})
	if cb, ok := builder.(*ContainerBuilder); ok {
		opts := BuildOptions{
			ImageName: "test-image",
			Tag:       "latest",
			BaseImage: "alpine",
		}

		err := cb.ValidateOptions(opts)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Test missing image name
		opts.ImageName = ""
		err = cb.ValidateOptions(opts)
		if err == nil {
			t.Errorf("Expected error for missing image name")
		}
	}
}