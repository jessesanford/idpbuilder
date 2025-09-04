package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestBuildCommand(t *testing.T) {
	// Set up test environment
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	// Create temporary directory for build context
	tempDir, err := os.MkdirTemp("", "build-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple Dockerfile
	dockerfile := filepath.Join(tempDir, "Dockerfile")
	dockerfileContent := `FROM alpine:latest
RUN echo "Hello World"
CMD ["echo", "test"]`
	
	if err := os.WriteFile(dockerfile, []byte(dockerfileContent), 0644); err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	// Test build command creation
	if buildCmd == nil {
		t.Fatal("buildCmd should not be nil when feature flag is enabled")
	}

	// Test command properties
	if buildCmd.Use != "build [CONTEXT]" {
		t.Errorf("Expected Use to be 'build [CONTEXT]', got: %s", buildCmd.Use)
	}

	if buildCmd.Short == "" {
		t.Error("Build command should have a short description")
	}

	if buildCmd.Long == "" {
		t.Error("Build command should have a long description")
	}
}

func TestValidateBuildOptions(t *testing.T) {
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	tests := []struct {
		name      string
		setupOpts func()
		wantError bool
	}{
		{
			name: "valid options",
			setupOpts: func() {
				tempDir, _ := os.MkdirTemp("", "valid-test-*")
				defer os.RemoveAll(tempDir)
				
				dockerfile := filepath.Join(tempDir, "Dockerfile")
				os.WriteFile(dockerfile, []byte("FROM alpine"), 0644)
				
				buildOpts.Context = tempDir
				buildOpts.Dockerfile = "Dockerfile"
				buildOpts.Platform = "linux/amd64"
				buildOpts.Tags = []string{"test:latest"}
			},
			wantError: false,
		},
		{
			name: "nonexistent context",
			setupOpts: func() {
				buildOpts.Context = "/nonexistent/path"
				buildOpts.Dockerfile = "Dockerfile"
				buildOpts.Platform = "linux/amd64"
				buildOpts.Tags = []string{"test:latest"}
			},
			wantError: true,
		},
		{
			name: "invalid platform",
			setupOpts: func() {
				tempDir, _ := os.MkdirTemp("", "platform-test-*")
				defer os.RemoveAll(tempDir)
				
				buildOpts.Context = tempDir
				buildOpts.Dockerfile = "Dockerfile"
				buildOpts.Platform = "invalid-platform"
				buildOpts.Tags = []string{"test:latest"}
			},
			wantError: true,
		},
		{
			name: "invalid tag",
			setupOpts: func() {
				tempDir, _ := os.MkdirTemp("", "tag-test-*")
				defer os.RemoveAll(tempDir)
				
				buildOpts.Context = tempDir
				buildOpts.Dockerfile = "Dockerfile"
				buildOpts.Platform = "linux/amd64"
				buildOpts.Tags = []string{"invalid tag with spaces"}
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset buildOpts
			buildOpts = &BuildOptions{
				Labels: make(map[string]string),
			}
			
			tt.setupOpts()
			
			err := validateBuildOptions()
			if (err != nil) != tt.wantError {
				t.Errorf("validateBuildOptions() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestIsFeatureEnabled(t *testing.T) {
	// Test enabled
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	if !isFeatureEnabled() {
		t.Error("isFeatureEnabled() should return true when ENABLE_CLI_TOOLS=true")
	}

	// Test disabled
	os.Setenv("ENABLE_CLI_TOOLS", "false")
	if isFeatureEnabled() {
		t.Error("isFeatureEnabled() should return false when ENABLE_CLI_TOOLS=false")
	}

	// Test unset
	os.Unsetenv("ENABLE_CLI_TOOLS")
	if isFeatureEnabled() {
		t.Error("isFeatureEnabled() should return false when ENABLE_CLI_TOOLS is unset")
	}
}

func TestQuietLogger(t *testing.T) {
	logger := &quietLogger{}
	
	// These should not panic or produce output
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	
	// Test with arguments
	logger.Debug("debug with args: %s %d", "test", 42)
	logger.Info("info with args: %s %d", "test", 42)
	logger.Warn("warn with args: %s %d", "test", 42)
	logger.Error("error with args: %s %d", "test", 42)
}

func TestBuildCommandFlags(t *testing.T) {
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	// Test that flags are properly defined
	flags := buildCmd.Flags()
	
	flagTests := []struct {
		name     string
		flagName string
	}{
		{"file flag", "file"},
		{"tag flag", "tag"},
		{"platform flag", "platform"},
		{"output flag", "output"},
		{"label flag", "label"},
		{"no-cache flag", "no-cache"},
		{"pull flag", "pull"},
		{"quiet flag", "quiet"},
	}

	for _, tt := range flagTests {
		t.Run(tt.name, func(t *testing.T) {
			flag := flags.Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("Flag %s should be defined", tt.flagName)
			}
		})
	}
}

func TestRunBuildWithoutFeatureFlag(t *testing.T) {
	// Disable feature flag
	os.Unsetenv("ENABLE_CLI_TOOLS")

	cmd := &cobra.Command{}
	err := runBuild(cmd, []string{"."})
	
	if err == nil {
		t.Error("runBuild should fail when feature flag is disabled")
	}

	if err != nil && err.Error() != "build command requires ENABLE_CLI_TOOLS=true" {
		t.Errorf("Expected specific error message, got: %v", err)
	}
}