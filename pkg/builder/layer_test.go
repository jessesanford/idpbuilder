package builder

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLayerManager(t *testing.T) {
	// Set up test environment
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	opts := DefaultBuildOptions().WithTags("test:latest")
	lm := NewLayerManager(opts)

	if lm == nil {
		t.Fatal("NewLayerManager returned nil")
	}

	if lm.options != opts {
		t.Error("LayerManager options not set correctly")
	}
}

func TestCreateLayerFromDirectory(t *testing.T) {
	// Set up test environment  
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "layer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	opts := DefaultBuildOptions().WithTags("test:latest")
	lm := NewLayerManager(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	layer, err := lm.CreateLayerFromDirectory(ctx, tempDir)
	if err != nil {
		t.Fatalf("Failed to create layer from directory: %v", err)
	}

	if layer == nil {
		t.Fatal("Created layer is nil")
	}

	// Verify layer properties
	digest, err := layer.Digest()
	if err != nil {
		t.Fatalf("Failed to get layer digest: %v", err)
	}

	if digest.String() == "" {
		t.Error("Layer digest is empty")
	}

	size, err := layer.Size()
	if err != nil {
		t.Fatalf("Failed to get layer size: %v", err)
	}

	if size <= 0 {
		t.Error("Layer size should be greater than 0")
	}
}

func TestCreateLayerFromFiles(t *testing.T) {
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "layer-files-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	
	if err := os.WriteFile(file1, []byte("content 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(file2, []byte("content 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}

	opts := DefaultBuildOptions().WithTags("test:latest")
	lm := NewLayerManager(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	files := []string{file1, file2}
	layer, err := lm.CreateLayerFromFiles(ctx, files, tempDir)
	if err != nil {
		t.Fatalf("Failed to create layer from files: %v", err)
	}

	if layer == nil {
		t.Fatal("Created layer is nil")
	}
}

func TestCreateEmptyLayer(t *testing.T) {
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	opts := DefaultBuildOptions().WithTags("test:latest")
	lm := NewLayerManager(opts)

	layer, err := lm.CreateEmptyLayer()
	if err != nil {
		t.Fatalf("Failed to create empty layer: %v", err)
	}

	if layer == nil {
		t.Fatal("Empty layer is nil")
	}

	// Verify it's actually empty
	size, err := layer.Size()
	if err != nil {
		t.Fatalf("Failed to get empty layer size: %v", err)
	}

	// Empty compressed tar should be very small but not zero
	if size >= 1024 {
		t.Errorf("Empty layer size too large: %d bytes", size)
	}
}

func TestGetLayerInfo(t *testing.T) {
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	opts := DefaultBuildOptions().WithTags("test:latest")
	lm := NewLayerManager(opts)

	layer, err := lm.CreateEmptyLayer()
	if err != nil {
		t.Fatalf("Failed to create empty layer: %v", err)
	}

	info, err := lm.GetLayerInfo(layer)
	if err != nil {
		t.Fatalf("Failed to get layer info: %v", err)
	}

	if info == nil {
		t.Fatal("Layer info is nil")
	}

	if info.Digest == "" {
		t.Error("Layer info digest is empty")
	}

	if info.Size < 0 {
		t.Error("Layer info size is negative")
	}

	if info.MediaType == "" {
		t.Error("Layer info media type is empty")
	}
}

func TestExtractLayer(t *testing.T) {
	os.Setenv("ENABLE_CLI_TOOLS", "true")
	defer os.Unsetenv("ENABLE_CLI_TOOLS")

	// Create source directory
	srcDir, err := os.MkdirTemp("", "extract-src-*")
	if err != nil {
		t.Fatalf("Failed to create source temp dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	// Create test file in source
	testContent := "extract test content"
	testFile := filepath.Join(srcDir, "extract-test.txt")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create destination directory
	destDir, err := os.MkdirTemp("", "extract-dest-*")
	if err != nil {
		t.Fatalf("Failed to create dest temp dir: %v", err)
	}
	defer os.RemoveAll(destDir)

	opts := DefaultBuildOptions().WithTags("test:latest")
	lm := NewLayerManager(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create layer from source
	layer, err := lm.CreateLayerFromDirectory(ctx, srcDir)
	if err != nil {
		t.Fatalf("Failed to create layer: %v", err)
	}

	// Extract layer to destination
	if err := lm.ExtractLayer(ctx, layer, destDir); err != nil {
		t.Fatalf("Failed to extract layer: %v", err)
	}

	// Verify extracted file
	extractedFile := filepath.Join(destDir, "extract-test.txt")
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Extracted content mismatch. Expected: %s, Got: %s", testContent, string(content))
	}
}

func TestFeatureFlagDisabled(t *testing.T) {
	// Ensure feature flag is disabled
	os.Unsetenv("ENABLE_CLI_TOOLS")

	opts := DefaultBuildOptions().WithTags("test:latest")
	lm := NewLayerManager(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// All operations should fail when feature flag is disabled
	_, err := lm.CreateLayerFromDirectory(ctx, ".")
	if err == nil {
		t.Error("CreateLayerFromDirectory should fail when feature flag is disabled")
	}

	_, err = lm.CreateLayerFromFiles(ctx, []string{"test.txt"}, ".")
	if err == nil {
		t.Error("CreateLayerFromFiles should fail when feature flag is disabled")
	}

	_, err = lm.CreateEmptyLayer()
	if err == nil {
		t.Error("CreateEmptyLayer should fail when feature flag is disabled")
	}
}