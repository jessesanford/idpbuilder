package builder

import (
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func TestNewBuildOptions(t *testing.T) {
	opts := NewBuildOptions()

	if opts == nil {
		t.Fatal("NewBuildOptions returned nil")
	}

	if opts.Platform == nil {
		t.Error("Platform should not be nil")
	}

	if opts.Platform.OS != "linux" {
		t.Errorf("Expected default OS to be 'linux', got %s", opts.Platform.OS)
	}

	if opts.Platform.Architecture != "amd64" {
		t.Errorf("Expected default architecture to be 'amd64', got %s", opts.Platform.Architecture)
	}

	if opts.Labels == nil {
		t.Error("Labels map should be initialized")
	}

	if opts.Environment == nil {
		t.Error("Environment map should be initialized")
	}

	if opts.FeatureFlags == nil {
		t.Error("FeatureFlags map should be initialized")
	}
}

func TestBuildOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *BuildOptions
		wantErr bool
	}{
		{
			name:    "valid options",
			opts:    NewBuildOptions(),
			wantErr: false,
		},
		{
			name: "nil platform",
			opts: &BuildOptions{
				Platform: nil,
			},
			wantErr: true,
		},
		{
			name: "empty OS",
			opts: &BuildOptions{
				Platform: &v1.Platform{
					Architecture: "amd64",
					OS:           "",
				},
			},
			wantErr: true,
		},
		{
			name: "empty architecture",
			opts: &BuildOptions{
				Platform: &v1.Platform{
					Architecture: "",
					OS:           "linux",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tag",
			opts: &BuildOptions{
				Platform: &v1.Platform{
					Architecture: "amd64",
					OS:           "linux",
				},
				Tags: []string{"invalid tag"},
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			opts: &BuildOptions{
				Platform: &v1.Platform{
					Architecture: "amd64",
					OS:           "linux",
				},
				ExposedPorts: []string{"invalid/port/format"},
			},
			wantErr: true,
		},
		{
			name: "relative context path",
			opts: &BuildOptions{
				Platform: &v1.Platform{
					Architecture: "amd64",
					OS:           "linux",
				},
				ContextPath: "relative/path",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildOptions.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBuildOptions_SetFeatureFlag(t *testing.T) {
	opts := NewBuildOptions()

	opts.SetFeatureFlag("test_feature", true)

	if !opts.IsFeatureEnabled("test_feature") {
		t.Error("Feature flag should be enabled")
	}

	opts.SetFeatureFlag("test_feature", false)

	if opts.IsFeatureEnabled("test_feature") {
		t.Error("Feature flag should be disabled")
	}
}

func TestBuildOptions_AddLabel(t *testing.T) {
	opts := NewBuildOptions()

	opts.AddLabel("test.label", "test-value")

	if value, exists := opts.Labels["test.label"]; !exists || value != "test-value" {
		t.Errorf("Expected label 'test.label' with value 'test-value', got %s (exists: %v)", value, exists)
	}
}

func TestBuildOptions_AddEnvironment(t *testing.T) {
	opts := NewBuildOptions()

	opts.AddEnvironment("TEST_VAR", "test-value")

	if value, exists := opts.Environment["TEST_VAR"]; !exists || value != "test-value" {
		t.Errorf("Expected environment variable 'TEST_VAR' with value 'test-value', got %s (exists: %v)", value, exists)
	}
}

func TestValidateTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		wantErr bool
	}{
		{"valid tag", "my-app:latest", false},
		{"empty tag", "", true},
		{"tag with spaces", "my app:latest", true},
		{"tag starting with dash", "-invalid", true},
		{"simple tag", "latest", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		wantErr bool
	}{
		{"tcp port", "8080/tcp", false},
		{"udp port", "53/udp", false},
		{"port without protocol", "8080", false},
		{"empty port", "", true},
		{"invalid protocol", "8080/http", true},
		{"too many parts", "8080/tcp/extra", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}