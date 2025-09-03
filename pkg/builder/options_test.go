package builder

import (
	"reflect"
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

func TestBuildOptions_AddBuildArg(t *testing.T) {
	opts := NewBuildOptions()
	
	opts.AddBuildArg("TEST_ARG", "test-value")
	
	if value, exists := opts.BuildArgs["TEST_ARG"]; !exists || value != "test-value" {
		t.Errorf("Expected build arg 'TEST_ARG' with value 'test-value', got %s (exists: %v)", value, exists)
	}
}

func TestBuildOptions_AddTag(t *testing.T) {
	opts := NewBuildOptions()
	
	// Test adding valid tag
	err := opts.AddTag("my-app:latest")
	if err != nil {
		t.Errorf("AddTag() error = %v", err)
	}
	
	if len(opts.Tags) != 1 || opts.Tags[0] != "my-app:latest" {
		t.Errorf("Expected tag 'my-app:latest' to be added")
	}
	
	// Test adding invalid tag
	err = opts.AddTag("invalid tag")
	if err == nil {
		t.Error("AddTag with invalid tag should return error")
	}
}

func TestBuildOptions_AddExposedPort(t *testing.T) {
	opts := NewBuildOptions()
	
	// Test adding valid port
	err := opts.AddExposedPort("8080/tcp")
	if err != nil {
		t.Errorf("AddExposedPort() error = %v", err)
	}
	
	if len(opts.ExposedPorts) != 1 || opts.ExposedPorts[0] != "8080/tcp" {
		t.Errorf("Expected port '8080/tcp' to be added")
	}
	
	// Test adding invalid port
	err = opts.AddExposedPort("invalid/port/format")
	if err == nil {
		t.Error("AddExposedPort with invalid port should return error")
	}
}

func TestBuildOptions_SetPlatform(t *testing.T) {
	opts := NewBuildOptions()
	
	opts.SetPlatform("windows", "arm64")
	
	if opts.Platform.OS != "windows" {
		t.Errorf("Expected OS 'windows', got %s", opts.Platform.OS)
	}
	if opts.Platform.Architecture != "arm64" {
		t.Errorf("Expected architecture 'arm64', got %s", opts.Platform.Architecture)
	}
}

func TestBuildOptions_Clone(t *testing.T) {
	opts := NewBuildOptions()
	opts.AddLabel("test.label", "test-value")
	opts.AddEnvironment("TEST_ENV", "env-value")
	opts.AddBuildArg("TEST_ARG", "arg-value")
	opts.SetFeatureFlag("test_feature", true)
	opts.WorkingDir = "/test"
	opts.User = "testuser"
	err := opts.AddTag("test:latest")
	if err != nil {
		t.Fatalf("AddTag() error = %v", err)
	}
	err = opts.AddExposedPort("8080/tcp")
	if err != nil {
		t.Fatalf("AddExposedPort() error = %v", err)
	}
	opts.Entrypoint = []string{"./app"}
	opts.Cmd = []string{"run"}
	
	clone := opts.Clone()
	
	// Verify all values are copied
	if clone.WorkingDir != "/test" {
		t.Errorf("Expected WorkingDir '/test', got '%s'", clone.WorkingDir)
	}
	if clone.User != "testuser" {
		t.Errorf("Expected User 'testuser', got '%s'", clone.User)
	}
	if clone.Labels["test.label"] != "test-value" {
		t.Error("Label should be copied")
	}
	if clone.Environment["TEST_ENV"] != "env-value" {
		t.Error("Environment variable should be copied")
	}
	if clone.BuildArgs["TEST_ARG"] != "arg-value" {
		t.Error("Build arg should be copied")
	}
	if !clone.FeatureFlags["test_feature"] {
		t.Error("Feature flag should be copied")
	}
	if len(clone.Tags) != 1 || clone.Tags[0] != "test:latest" {
		t.Error("Tags should be copied")
	}
	if len(clone.ExposedPorts) != 1 || clone.ExposedPorts[0] != "8080/tcp" {
		t.Error("ExposedPorts should be copied")
	}
	if len(clone.Entrypoint) != 1 || clone.Entrypoint[0] != "./app" {
		t.Error("Entrypoint should be copied")
	}
	if len(clone.Cmd) != 1 || clone.Cmd[0] != "run" {
		t.Error("Cmd should be copied")
	}
	
	// Verify they are different objects
	if reflect.ValueOf(opts).Pointer() == reflect.ValueOf(clone).Pointer() {
		t.Error("Clone should return a different object")
	}
	
	// Modify original and ensure clone is not affected
	opts.WorkingDir = "/modified"
	opts.AddLabel("new.label", "new-value")
	
	if clone.WorkingDir != "/test" {
		t.Error("Clone should not be affected by modifications to original")
	}
	if _, exists := clone.Labels["new.label"]; exists {
		t.Error("Clone should not have new label from original")
	}
}

func TestBuildOptions_GetAllLabels(t *testing.T) {
	opts := NewBuildOptions()
	opts.AddLabel("label1", "value1")
	opts.AddLabel("label2", "value2")
	
	labels := opts.GetAllLabels()
	
	if len(labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(labels))
	}
	if labels["label1"] != "value1" {
		t.Error("label1 should have value1")
	}
	if labels["label2"] != "value2" {
		t.Error("label2 should have value2")
	}
	
	// Modify returned map and ensure original is not affected
	labels["new"] = "new-value"
	if _, exists := opts.Labels["new"]; exists {
		t.Error("Original labels should not be affected")
	}
}

func TestBuildOptions_GetAllEnvironment(t *testing.T) {
	opts := NewBuildOptions()
	opts.AddEnvironment("ENV1", "value1")
	opts.AddEnvironment("ENV2", "value2")
	
	env := opts.GetAllEnvironment()
	
	if len(env) != 2 {
		t.Errorf("Expected 2 environment variables, got %d", len(env))
	}
	if env["ENV1"] != "value1" {
		t.Error("ENV1 should have value1")
	}
	if env["ENV2"] != "value2" {
		t.Error("ENV2 should have value2")
	}
	
	// Modify returned map and ensure original is not affected
	env["NEW"] = "new-value"
	if _, exists := opts.Environment["NEW"]; exists {
		t.Error("Original environment should not be affected")
	}
}