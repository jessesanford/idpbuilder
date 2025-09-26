package buildah

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBuildContext(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	if ctx == nil {
		t.Fatal("NewBuildContext returned nil")
	}

	if ctx.GetRoot() != tmpDir {
		t.Errorf("Expected root %s, got %s", tmpDir, ctx.GetRoot())
	}
}

func TestBuildContext_AddFile(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	// Test adding a simple file
	err := ctx.AddFile("test.txt", []byte("test content"))
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	// Verify file exists on filesystem
	content, err := os.ReadFile(filepath.Join(tmpDir, "test.txt"))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != "test content" {
		t.Errorf("Expected 'test content', got %s", string(content))
	}

	// Test adding file in subdirectory
	err = ctx.AddFile("subdir/nested.txt", []byte("nested content"))
	if err != nil {
		t.Fatalf("AddFile with subdirectory failed: %v", err)
	}

	// Verify nested file exists
	nestedContent, err := os.ReadFile(filepath.Join(tmpDir, "subdir", "nested.txt"))
	if err != nil {
		t.Fatalf("Failed to read nested file: %v", err)
	}

	if string(nestedContent) != "nested content" {
		t.Errorf("Expected 'nested content', got %s", string(nestedContent))
	}
}

func TestBuildContext_AddFile_ErrorCases(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	// Test empty path
	err := ctx.AddFile("", []byte("content"))
	if err == nil {
		t.Error("Expected error for empty path")
	}

	// Test absolute path
	err = ctx.AddFile("/absolute/path", []byte("content"))
	if err == nil {
		t.Error("Expected error for absolute path")
	}
}

func TestBuildContext_AddDirectory(t *testing.T) {
	// Create source directory structure
	srcDir := t.TempDir()
	os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
	os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644)

	// Create destination context
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	// Add directory
	err := ctx.AddDirectory(srcDir, "copied")
	if err != nil {
		t.Fatalf("AddDirectory failed: %v", err)
	}

	// Verify files were copied
	content1, err := os.ReadFile(filepath.Join(tmpDir, "copied", "file1.txt"))
	if err != nil {
		t.Fatalf("Failed to read copied file1: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("Expected 'content1', got %s", string(content1))
	}

	content2, err := os.ReadFile(filepath.Join(tmpDir, "copied", "subdir", "file2.txt"))
	if err != nil {
		t.Fatalf("Failed to read copied file2: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("Expected 'content2', got %s", string(content2))
	}
}

func TestBuildContext_AddDirectory_ErrorCases(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	// Test empty source
	err := ctx.AddDirectory("", "dest")
	if err == nil {
		t.Error("Expected error for empty source")
	}

	// Test empty destination
	err = ctx.AddDirectory("src", "")
	if err == nil {
		t.Error("Expected error for empty destination")
	}

	// Test absolute destination
	err = ctx.AddDirectory("src", "/absolute/dest")
	if err == nil {
		t.Error("Expected error for absolute destination")
	}

	// Test non-existent source
	err = ctx.AddDirectory("/nonexistent", "dest")
	if err == nil {
		t.Error("Expected error for non-existent source")
	}
}

func TestBuildContext_GetSize(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	// Add test files
	err := ctx.AddFile("file1.txt", []byte("content1"))
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	err = ctx.AddFile("file2.txt", []byte("content2"))
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	size, err := ctx.GetSize()
	if err != nil {
		t.Fatalf("GetSize failed: %v", err)
	}

	expectedSize := int64(len("content1") + len("content2"))
	if size != expectedSize {
		t.Errorf("Expected size %d, got %d", expectedSize, size)
	}
}

func TestBuildContext_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	// Should fail without Dockerfile
	err := ctx.Validate()
	if err == nil {
		t.Error("Expected validation to fail without Dockerfile")
	}

	// Add Dockerfile and validate
	err = ctx.AddFile("Dockerfile", []byte("FROM alpine"))
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	err = ctx.Validate()
	if err != nil {
		t.Errorf("Validation failed with Dockerfile: %v", err)
	}

	// Test with different Dockerfile names
	tmpDir2 := t.TempDir()
	ctx2 := NewBuildContext(tmpDir2)

	err = ctx2.AddFile("dockerfile", []byte("FROM ubuntu"))
	if err != nil {
		t.Fatalf("AddFile failed: %v", err)
	}

	err = ctx2.Validate()
	if err != nil {
		t.Errorf("Validation failed with lowercase dockerfile: %v", err)
	}
}

func TestBuildContext_Clean(t *testing.T) {
	// Test cleaning non-temp directory (should not remove)
	tmpDir := t.TempDir()
	ctx := NewBuildContext(tmpDir)

	err := ctx.Clean()
	if err != nil {
		t.Errorf("Clean failed: %v", err)
	}

	// Directory should still exist
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		t.Error("Directory was removed when it shouldn't have been")
	}
}

func TestNewBuildContextManager(t *testing.T) {
	workDir := "/tmp/work"
	manager := NewBuildContextManager(workDir)

	if manager == nil {
		t.Fatal("NewBuildContextManager returned nil")
	}
}

func TestBuildContextManager_CreateContext(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewBuildContextManager(tmpDir)

	// Create test Dockerfile
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	err := os.WriteFile(dockerfilePath, []byte("FROM alpine"), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	ctx, err := manager.CreateContext(context.Background(), dockerfilePath)
	if err != nil {
		t.Fatalf("CreateContext failed: %v", err)
	}

	if ctx == nil {
		t.Error("Expected non-nil context")
	}

	if ctx.GetRoot() != tmpDir {
		t.Errorf("Expected root %s, got %s", tmpDir, ctx.GetRoot())
	}
}

func TestBuildContextManager_CreateContext_ErrorCases(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewBuildContextManager(tmpDir)

	// Test empty dockerfile path
	_, err := manager.CreateContext(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty dockerfile path")
	}

	// Test non-existent dockerfile
	_, err = manager.CreateContext(context.Background(), "/nonexistent/Dockerfile")
	if err == nil {
		t.Error("Expected error for non-existent dockerfile")
	}
}

func TestBuildContextManager_CreateFromDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewBuildContextManager("/tmp/work")

	// Create some files in the directory
	err := os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	ctx, err := manager.CreateFromDirectory(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("CreateFromDirectory failed: %v", err)
	}

	if ctx == nil {
		t.Error("Expected non-nil context")
	}

	if ctx.GetRoot() != tmpDir {
		t.Errorf("Expected root %s, got %s", tmpDir, ctx.GetRoot())
	}
}

func TestBuildContextManager_CreateFromDirectory_ErrorCases(t *testing.T) {
	manager := NewBuildContextManager("/tmp/work")

	// Test empty directory path
	_, err := manager.CreateFromDirectory(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty directory path")
	}

	// Test non-existent directory
	_, err = manager.CreateFromDirectory(context.Background(), "/nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestBuildContextManager_CreateTarball(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewBuildContextManager("/tmp/work")
	ctx := NewBuildContext(tmpDir)

	tarballPath, err := manager.CreateTarball(ctx)
	if err != nil {
		t.Fatalf("CreateTarball failed: %v", err)
	}

	// For this basic implementation, it should return the context root
	if tarballPath != tmpDir {
		t.Errorf("Expected tarball path %s, got %s", tmpDir, tarballPath)
	}

	// Test nil context
	_, err = manager.CreateTarball(nil)
	if err == nil {
		t.Error("Expected error for nil context")
	}
}