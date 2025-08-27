package contexts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewContextResolver(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		resolver := NewContextResolver(nil)
		
		if resolver == nil {
			t.Fatal("NewContextResolver() returned nil")
		}
		
		if resolver.config == nil {
			t.Error("NewContextResolver() config is nil")
		}
		
		// Should use default config
		expectedMaxSize := int64(500 * 1024 * 1024)
		if resolver.config.MaxSize != expectedMaxSize {
			t.Errorf("NewContextResolver() MaxSize = %v, want %v", 
				resolver.config.MaxSize, expectedMaxSize)
		}
		
		if resolver.httpClient == nil {
			t.Error("NewContextResolver() httpClient is nil")
		}
		
		if resolver.cleanupFns == nil {
			t.Error("NewContextResolver() cleanupFns is nil")
		}
	})
	
	t.Run("with custom config", func(t *testing.T) {
		customConfig := &ContextConfig{
			MaxSize:      100 * 1024 * 1024,
			CacheEnabled: false,
			TempDir:      "/custom/temp",
			HTTPTimeout:  60 * time.Second,
		}
		
		resolver := NewContextResolver(customConfig)
		
		if resolver.config != customConfig {
			t.Error("NewContextResolver() did not use provided config")
		}
		
		if resolver.config.MaxSize != customConfig.MaxSize {
			t.Errorf("NewContextResolver() MaxSize = %v, want %v",
				resolver.config.MaxSize, customConfig.MaxSize)
		}
	})
}

func TestContextResolver_detectContextType(t *testing.T) {
	resolver := NewContextResolver(nil)
	
	tests := []struct {
		name     string
		source   string
		expected ContextType
	}{
		{"HTTP URL", "http://example.com/file.txt", URLContext},
		{"HTTPS URL", "https://example.com/file.txt", URLContext},
		{"Git URL with git://", "git://github.com/user/repo.git", GitContext},
		{"Git URL with git@", "git@github.com:user/repo.git", GitContext},
		{"HTTPS Git URL", "https://github.com/user/repo.git", URLContext}, // URLs are detected first
		{"Git URL ending with .git", "https://example.com/repo.git", URLContext}, // URLs are detected first
		{"Tar archive", "archive.tar", ArchiveContext},
		{"Tar.gz archive", "archive.tar.gz", ArchiveContext},
		{"Tgz archive", "archive.tgz", ArchiveContext},
		{"Tar.bz2 archive", "archive.tar.bz2", ArchiveContext},
		{"Tbz2 archive", "archive.tbz2", ArchiveContext},
		{"Zip archive", "archive.zip", ArchiveContext},
		{"Mixed case zip", "Archive.ZIP", ArchiveContext},
		{"Local path", "/path/to/directory", LocalContext},
		{"Relative path", "./directory", LocalContext},
		{"Current directory", ".", LocalContext},
		{"Unknown extension", "file.unknown", LocalContext},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.detectContextType(tt.source)
			if result != tt.expected {
				t.Errorf("detectContextType(%q) = %v, want %v", 
					tt.source, result, tt.expected)
			}
		})
	}
}

func TestContextResolver_isArchiveFile(t *testing.T) {
	resolver := NewContextResolver(nil)
	
	tests := []struct {
		name     string
		source   string
		expected bool
	}{
		{"tar file", "archive.tar", true},
		{"tar.gz file", "archive.tar.gz", true},
		{"tgz file", "archive.tgz", true},
		{"tar.bz2 file", "archive.tar.bz2", true},
		{"tbz2 file", "archive.tbz2", true},
		{"zip file", "archive.zip", true},
		{"mixed case ZIP", "ARCHIVE.ZIP", true},
		{"mixed case tar.GZ", "archive.tar.GZ", true},
		{"not archive", "file.txt", false},
		{"empty string", "", false},
		{"only extension", ".tar", true},
		{"no extension", "archive", false},
		{"partial match", "archive.tar.exe", false},
		{"path with archive ext", "/path/to/archive.zip", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolver.isArchiveFile(tt.source)
			if result != tt.expected {
				t.Errorf("isArchiveFile(%q) = %v, want %v", 
					tt.source, result, tt.expected)
			}
		})
	}
}

func TestContextResolver_ResolveContext(t *testing.T) {
	resolver := NewContextResolver(nil)
	
	t.Run("empty source", func(t *testing.T) {
		ctx, err := resolver.ResolveContext("")
		if err == nil {
			t.Error("ResolveContext(\"\") should return error")
		}
		if ctx != nil {
			t.Error("ResolveContext(\"\") should return nil context")
		}
		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("ResolveContext(\"\") error = %v, should mention empty", err)
		}
	})
	
	t.Run("git context not supported", func(t *testing.T) {
		ctx, err := resolver.ResolveContext("git://github.com/user/repo.git")
		if err == nil {
			t.Error("ResolveContext(git) should return error")
		}
		if ctx != nil {
			t.Error("ResolveContext(git) should return nil context")
		}
		if !strings.Contains(err.Error(), "split 002") {
			t.Errorf("ResolveContext(git) error = %v, should mention split 002", err)
		}
	})
	
	t.Run("archive context not supported", func(t *testing.T) {
		ctx, err := resolver.ResolveContext("archive.tar.gz")
		if err == nil {
			t.Error("ResolveContext(archive) should return error")
		}
		if ctx != nil {
			t.Error("ResolveContext(archive) should return nil context")
		}
		if !strings.Contains(err.Error(), "split 002") {
			t.Errorf("ResolveContext(archive) error = %v, should mention split 002", err)
		}
	})
	
	t.Run("URL context", func(t *testing.T) {
		ctx, err := resolver.ResolveContext("http://example.com/file.txt")
		if err != nil {
			t.Errorf("ResolveContext(URL) error = %v", err)
		}
		if ctx == nil {
			t.Error("ResolveContext(URL) returned nil context")
		}
		if ctx != nil && ctx.Type() != URLContext {
			t.Errorf("ResolveContext(URL) type = %v, want %v", ctx.Type(), URLContext)
		}
	})
}

func TestContextResolver_createLocalContext(t *testing.T) {
	resolver := NewContextResolver(nil)
	
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_local_context")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	t.Run("valid directory", func(t *testing.T) {
		ctx, err := resolver.createLocalContext(tempDir)
		if err != nil {
			t.Errorf("createLocalContext() error = %v", err)
		}
		if ctx == nil {
			t.Error("createLocalContext() returned nil context")
		}
		if ctx != nil {
			if ctx.Type() != LocalContext {
				t.Errorf("createLocalContext() type = %v, want %v", ctx.Type(), LocalContext)
			}
			
			// Path should be absolute
			if !filepath.IsAbs(ctx.Path()) {
				t.Errorf("createLocalContext() path %q is not absolute", ctx.Path())
			}
		}
	})
	
	t.Run("non-existent directory", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "does-not-exist")
		ctx, err := resolver.createLocalContext(nonExistentPath)
		if err == nil {
			t.Error("createLocalContext() should return error for non-existent path")
		}
		if ctx != nil {
			t.Error("createLocalContext() should return nil for non-existent path")
		}
		if !strings.Contains(err.Error(), "does not exist") {
			t.Errorf("createLocalContext() error = %v, should mention non-existent", err)
		}
	})
	
	t.Run("file instead of directory", func(t *testing.T) {
		// Create a file
		testFile := filepath.Join(tempDir, "testfile.txt")
		if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		
		ctx, err := resolver.createLocalContext(testFile)
		if err == nil {
			t.Error("createLocalContext() should return error for file path")
		}
		if ctx != nil {
			t.Error("createLocalContext() should return nil for file path")
		}
		if !strings.Contains(err.Error(), "not a directory") {
			t.Errorf("createLocalContext() error = %v, should mention not directory", err)
		}
	})
	
	t.Run("relative path resolution", func(t *testing.T) {
		// Create a subdirectory in temp dir
		subDir := filepath.Join(tempDir, "subdir")
		if err := os.MkdirAll(subDir, 0755); err != nil {
			t.Fatalf("Failed to create subdir: %v", err)
		}
		
		// Change to temp dir to test relative path
		oldWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get working directory: %v", err)
		}
		defer os.Chdir(oldWd)
		
		if err := os.Chdir(tempDir); err != nil {
			t.Fatalf("Failed to change directory: %v", err)
		}
		
		ctx, err := resolver.createLocalContext("subdir")
		if err != nil {
			t.Errorf("createLocalContext() error = %v", err)
		}
		if ctx != nil && !filepath.IsAbs(ctx.Path()) {
			t.Errorf("createLocalContext() should return absolute path, got %q", ctx.Path())
		}
	})
}

func TestContextResolver_createURLContext(t *testing.T) {
	resolver := NewContextResolver(nil)
	
	t.Run("valid URL", func(t *testing.T) {
		ctx, err := resolver.createURLContext("http://example.com/file.txt")
		if err != nil {
			t.Errorf("createURLContext() error = %v", err)
		}
		if ctx == nil {
			t.Error("createURLContext() returned nil context")
		}
		if ctx != nil && ctx.Type() != URLContext {
			t.Errorf("createURLContext() type = %v, want %v", ctx.Type(), URLContext)
		}
		
		// Check that cleanup function was registered
		if len(resolver.cleanupFns) == 0 {
			t.Error("createURLContext() should register cleanup function")
		}
	})
	
	t.Run("invalid URL", func(t *testing.T) {
		ctx, err := resolver.createURLContext("ht!tp://invalid url with spaces")
		if err == nil {
			t.Error("createURLContext() should return error for invalid URL")
		}
		if ctx != nil {
			t.Error("createURLContext() should return nil for invalid URL")
		}
		if err != nil && !strings.Contains(err.Error(), "invalid URL") {
			t.Errorf("createURLContext() error = %v, should mention invalid URL", err)
		}
	})
	
	t.Run("URL with special characters", func(t *testing.T) {
		ctx, err := resolver.createURLContext("https://example.com/path with spaces")
		// URL parsing is actually quite lenient, so this might succeed
		if err != nil {
			// If it fails, that's also acceptable behavior
			t.Logf("createURLContext() properly rejected malformed URL: %v", err)
		} else if ctx == nil {
			t.Error("createURLContext() returned nil context without error")
		}
	})
	
	t.Run("properly encoded URL", func(t *testing.T) {
		ctx, err := resolver.createURLContext("https://example.com/path%20with%20spaces")
		if err != nil {
			t.Errorf("createURLContext() error = %v", err)
		}
		if ctx == nil {
			t.Error("createURLContext() returned nil for valid encoded URL")
		}
	})
}

func TestContextResolver_Cleanup(t *testing.T) {
	resolver := NewContextResolver(nil)
	
	t.Run("no cleanup functions", func(t *testing.T) {
		err := resolver.Cleanup()
		if err != nil {
			t.Errorf("Cleanup() with no functions error = %v", err)
		}
	})
	
	t.Run("successful cleanup", func(t *testing.T) {
		called := false
		resolver.addCleanupFn(func() error {
			called = true
			return nil
		})
		
		err := resolver.Cleanup()
		if err != nil {
			t.Errorf("Cleanup() error = %v", err)
		}
		if !called {
			t.Error("Cleanup() did not call registered function")
		}
	})
	
	t.Run("cleanup with errors", func(t *testing.T) {
		resolver.cleanupFns = nil // Reset
		
		resolver.addCleanupFn(func() error {
			return fmt.Errorf("error1")
		})
		resolver.addCleanupFn(func() error {
			return fmt.Errorf("error2")
		})
		
		err := resolver.Cleanup()
		if err == nil {
			t.Error("Cleanup() should return error when cleanup functions fail")
		}
		if !strings.Contains(err.Error(), "error1") || !strings.Contains(err.Error(), "error2") {
			t.Errorf("Cleanup() error = %v, should contain both errors", err)
		}
	})
	
	t.Run("partial cleanup failures", func(t *testing.T) {
		resolver.cleanupFns = nil // Reset
		
		success := false
		resolver.addCleanupFn(func() error {
			success = true
			return nil
		})
		resolver.addCleanupFn(func() error {
			return fmt.Errorf("cleanup failed")
		})
		
		err := resolver.Cleanup()
		if err == nil {
			t.Error("Cleanup() should return error when some functions fail")
		}
		if !success {
			t.Error("Cleanup() should still call successful functions")
		}
	})
}

func TestLocalContextImpl(t *testing.T) {
	testPath := "/test/path"
	ctx := &LocalContextImpl{path: testPath}
	
	t.Run("Path", func(t *testing.T) {
		if ctx.Path() != testPath {
			t.Errorf("LocalContextImpl.Path() = %v, want %v", ctx.Path(), testPath)
		}
	})
	
	t.Run("Type", func(t *testing.T) {
		if ctx.Type() != LocalContext {
			t.Errorf("LocalContextImpl.Type() = %v, want %v", ctx.Type(), LocalContext)
		}
	})
	
	t.Run("Cleanup", func(t *testing.T) {
		err := ctx.Cleanup()
		if err != nil {
			t.Errorf("LocalContextImpl.Cleanup() error = %v, want nil", err)
		}
	})
}

func TestContextResolver_addCleanupFn(t *testing.T) {
	resolver := NewContextResolver(nil)
	
	initialLen := len(resolver.cleanupFns)
	
	resolver.addCleanupFn(func() error { return nil })
	
	if len(resolver.cleanupFns) != initialLen+1 {
		t.Errorf("addCleanupFn() did not add function, len = %d, want %d", 
			len(resolver.cleanupFns), initialLen+1)
	}
}

// Test error edge cases
func TestContextResolver_EdgeCases(t *testing.T) {
	t.Run("resolver with zero timeout", func(t *testing.T) {
		config := &ContextConfig{HTTPTimeout: 0}
		resolver := NewContextResolver(config)
		
		if resolver.config.HTTPTimeout != 0 {
			t.Error("NewContextResolver should preserve zero timeout")
		}
	})
	
	t.Run("resolver with negative max size", func(t *testing.T) {
		config := &ContextConfig{MaxSize: -1}
		resolver := NewContextResolver(config)
		
		if resolver.config.MaxSize != -1 {
			t.Error("NewContextResolver should preserve negative max size")
		}
	})
}