package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDefaultBuildOptions tests the default build options
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