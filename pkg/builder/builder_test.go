package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	opts := BuildOptions{}
	builder, err := NewBuilder(opts)
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	if builder == nil {
		t.Fatal("Builder should not be nil")
	}
}

func TestBuild(t *testing.T) {
	// Create test directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	opts := DefaultBuildOptions()
	builder, err := NewBuilder(opts)
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	
	ctx := context.Background()
	image, err := builder.Build(ctx, tmpDir, opts)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	if image == nil {
		t.Fatal("Image should not be nil")
	}
}

func TestBuildTarball(t *testing.T) {
	// Create test directory
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	outputPath := filepath.Join(tmpDir, "output.tar")
	opts := DefaultBuildOptions()
	builder, err := NewBuilder(opts)
	if err != nil {
		t.Fatalf("Failed to create builder: %v", err)
	}
	
	ctx := context.Background()
	err = builder.BuildTarball(ctx, tmpDir, outputPath, opts)
	if err != nil {
		t.Fatalf("BuildTarball failed: %v", err)
	}
	
	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Output tarball file was not created")
	}
}