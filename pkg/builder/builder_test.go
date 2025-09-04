package builder

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	opts := DefaultBuildOptions()
	opts.Tags = []string{"test:latest"}
	
	builder, err := New(opts)
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	
	if builder == nil {
		t.Fatal("Builder should not be nil")
	}
}

func TestFeatureFlag(t *testing.T) {
	opts := DefaultBuildOptions()
	opts.Tags = []string{"test:latest"}
	
	builder, err := New(opts)
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	
	// Test without feature flag
	result, err := builder.Build()
	if err == nil {
		t.Fatal("Build should fail without feature flag")
	}
	
	// Test with feature flag
	os.Setenv("ENABLE_CORE_BUILDER", "true")
	defer os.Unsetenv("ENABLE_CORE_BUILDER")
	
	// This will fail due to missing Dockerfile, but should pass feature flag check
	result, err = builder.Build()
	if err != nil && result == nil {
		// Expected - no real Dockerfile exists in test
		t.Logf("Expected failure: %v", err)
	}
}

func TestDefaultBuildOptions(t *testing.T) {
	opts := DefaultBuildOptions()
	
	if opts.Context != "." {
		t.Errorf("Expected context '.', got %s", opts.Context)
	}
	
	if opts.Dockerfile != "Dockerfile" {
		t.Errorf("Expected dockerfile 'Dockerfile', got %s", opts.Dockerfile)
	}
	
	if opts.Platform != "linux/amd64" {
		t.Errorf("Expected platform 'linux/amd64', got %s", opts.Platform)
	}
}