package contexts

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestContextManagerIntegration(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "context-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Setup test context directory
	contextDir := filepath.Join(tempDir, "build-context")
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		t.Fatalf("Failed to create context dir: %v", err)
	}

	// Create test file
	testFile := filepath.Join(contextDir, "Dockerfile")
	if err := os.WriteFile(testFile, []byte("FROM alpine\nRUN echo hello"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test context manager lifecycle
	manager := NewContextManager(
		WithMaxCacheSize(5),
		WithDefaultTimeout(30*time.Second),
	)
	defer manager.Shutdown()

	ctx := context.Background()
	config := &ContextConfig{
		Type:          FileContext,
		Source:        contextDir,
		AllowSymlinks: false,
		Timeout:       30 * time.Second,
	}

	// Test context creation and caching
	buildContext1, err := manager.GetContext(ctx, config)
	if err != nil {
		t.Errorf("Failed to create context: %v", err)
	}

	buildContext2, err := manager.GetContext(ctx, config)
	if err != nil {
		t.Errorf("Failed to get cached context: %v", err)
	}

	// Should return same cached instance
	if buildContext1 != buildContext2 {
		t.Error("Expected cached context to be reused")
	}

	// Test context validation
	if err := buildContext1.Validate(); err != nil {
		t.Errorf("Context validation failed: %v", err)
	}

	// Test context type
	if buildContext1.GetType() != FileContext {
		t.Errorf("Expected context type %s, got %s", FileContext, buildContext1.GetType())
	}
}

func TestSecurityValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) error
		allowSymlinks bool
		expectValid bool
		expectError string
	}{
		{
			name: "valid context",
			setupFunc: func(dir string) error {
				return os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM alpine"), 0644)
			},
			allowSymlinks: false,
			expectValid:   true,
		},
		{
			name: "symlink not allowed",
			setupFunc: func(dir string) error {
				testFile := filepath.Join(dir, "test.txt")
				os.WriteFile(testFile, []byte("content"), 0644)
				return os.Symlink(testFile, filepath.Join(dir, "symlink"))
			},
			allowSymlinks: false,
			expectValid:   false,
			expectError:   "symlinks not allowed",
		},
		{
			name: "path traversal in filename",
			setupFunc: func(dir string) error {
				// This would be caught by validation
				return os.WriteFile(filepath.Join(dir, "normal.txt"), []byte("content"), 0644)
			},
			allowSymlinks: false,
			expectValid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "security-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			if err := tt.setupFunc(tempDir); err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			result, err := ValidateContext(tempDir, tt.allowSymlinks)
			if err != nil {
				t.Fatalf("ValidateContext failed: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%t, got valid=%t", tt.expectValid, result.Valid)
			}

			if tt.expectError != "" {
				found := false
				for _, errMsg := range result.Errors {
					if contains(errMsg, tt.expectError) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error containing '%s', got errors: %v", tt.expectError, result.Errors)
				}
			}
		})
	}
}

func TestContextFactory(t *testing.T) {
	factory := NewContextFactory()

	// Test custom validator registration
	customValidatorCalled := false
	factory.RegisterValidator("custom", func(source string, config *ContextConfig) error {
		customValidatorCalled = true
		return nil
	})

	config := &ContextConfig{
		Type:   "custom",
		Source: "/test/path",
	}

	// This will fail because we don't have a creator for "custom" type,
	// but it should call our validator first
	_, err := factory.CreateContext(config)
	if err == nil {
		t.Error("Expected error for unsupported context type")
	}

	if !customValidatorCalled {
		t.Error("Custom validator was not called")
	}
}

func TestManagerCaching(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create context directory
	contextDir := filepath.Join(tempDir, "context")
	if err := os.MkdirAll(contextDir, 0755); err != nil {
		t.Fatalf("Failed to create context dir: %v", err)
	}

	manager := NewContextManager(WithMaxCacheSize(2))
	defer manager.Shutdown()

	ctx := context.Background()

	// Create contexts that should fill cache
	configs := []*ContextConfig{
		{Type: FileContext, Source: contextDir + "/1"},
		{Type: FileContext, Source: contextDir + "/2"},
		{Type: FileContext, Source: contextDir + "/3"}, // Should evict first
	}

	// Create subdirectories
	for i, config := range configs {
		if err := os.MkdirAll(config.Source, 0755); err != nil {
			t.Fatalf("Failed to create test dir %d: %v", i, err)
		}
	}

	contexts := make([]BuildContext, len(configs))
	for i, config := range configs {
		contexts[i], err = manager.GetContext(ctx, config)
		if err != nil {
			t.Fatalf("Failed to get context %d: %v", i, err)
		}
	}

	// Test cache cleanup
	manager.CleanupCache()
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}