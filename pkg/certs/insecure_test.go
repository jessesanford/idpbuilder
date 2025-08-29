package certs

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewInsecureMode(t *testing.T) {
	im := NewInsecureMode()
	if im == nil {
		t.Fatal("NewInsecureMode returned nil")
	}
	
	if im.logger == nil {
		t.Error("Expected logger to be initialized")
	}
	
	if im.auditLogger == nil {
		t.Error("Expected auditLogger to be initialized")
	}
	
	if im.warningsIssued == nil {
		t.Error("Expected warningsIssued map to be initialized")
	}
}

func TestIsInsecureModeRequested(t *testing.T) {
	im := NewInsecureMode()
	
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "no insecure flag",
			args:     []string{"cmd", "subcmd"},
			expected: false,
		},
		{
			name:     "with --insecure flag",
			args:     []string{"cmd", "--insecure", "subcmd"},
			expected: true,
		},
		{
			name:     "with -k flag",
			args:     []string{"cmd", "-k", "subcmd"},
			expected: true,
		},
		{
			name:     "insecure in middle of args",
			args:     []string{"cmd", "subcmd", "--insecure", "target"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			result := im.IsInsecureModeRequested()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestApplyInsecureMode(t *testing.T) {
	im := NewInsecureMode()
	ctx := context.Background()

	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	t.Run("without requiring explicit flag", func(t *testing.T) {
		config := &InsecureConfig{
			Registry:        "test-registry",
			Operation:       "push",
			Duration:        5 * time.Minute,
			Reason:          "testing",
			RequireExplicit: false,
		}
		
		err := im.ApplyInsecureMode(ctx, config)
		if err != nil {
			t.Fatalf("ApplyInsecureMode failed: %v", err)
		}
	})

	t.Run("requiring explicit flag with flag present", func(t *testing.T) {
		os.Args = []string{"cmd", "--insecure"}
		
		config := &InsecureConfig{
			Registry:        "test-registry",
			Operation:       "push",
			Duration:        5 * time.Minute,
			Reason:          "testing",
			RequireExplicit: true,
		}
		
		err := im.ApplyInsecureMode(ctx, config)
		if err != nil {
			t.Fatalf("ApplyInsecureMode failed: %v", err)
		}
	})

	t.Run("requiring explicit flag without flag", func(t *testing.T) {
		os.Args = []string{"cmd"}
		
		config := &InsecureConfig{
			Registry:        "test-registry",
			Operation:       "push",
			Duration:        5 * time.Minute,
			Reason:          "testing",
			RequireExplicit: true,
		}
		
		err := im.ApplyInsecureMode(ctx, config)
		if err == nil {
			t.Error("Expected error when insecure flag required but not present")
		}
		
		expectedMsg := "insecure mode requires explicit --insecure flag"
		if !strings.Contains(err.Error(), expectedMsg) {
			t.Errorf("Expected error message to contain %q, got %q", expectedMsg, err.Error())
		}
	})
}

func TestPromptUserConsent(t *testing.T) {
	im := NewInsecureMode()
	
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	config := &InsecureConfig{
		Registry:  "test-registry",
		Operation: "push",
		Duration:  5 * time.Minute,
		Reason:    "testing",
	}

	t.Run("with insecure flag", func(t *testing.T) {
		os.Args = []string{"cmd", "--insecure"}
		
		consent, err := im.PromptUserConsent(config)
		if err != nil {
			t.Fatalf("PromptUserConsent failed: %v", err)
		}
		
		if !consent {
			t.Error("Expected consent to be true when --insecure flag present")
		}
	})

	t.Run("without insecure flag", func(t *testing.T) {
		os.Args = []string{"cmd"}
		
		consent, err := im.PromptUserConsent(config)
		if err == nil {
			t.Error("Expected error when no insecure flag present")
		}
		
		if consent {
			t.Error("Expected consent to be false when no insecure flag")
		}
	})
}

func TestCreateTimeLimitedConfig(t *testing.T) {
	im := NewInsecureMode()
	
	registry := "test-registry"
	operation := "push"
	reason := "testing"
	duration := 10 * time.Minute
	
	config := im.CreateTimeLimitedConfig(registry, operation, reason, duration)
	
	if config.Registry != registry {
		t.Errorf("Expected registry %s, got %s", registry, config.Registry)
	}
	
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	
	if config.Reason != reason {
		t.Errorf("Expected reason %s, got %s", reason, config.Reason)
	}
	
	if config.Duration != duration {
		t.Errorf("Expected duration %v, got %v", duration, config.Duration)
	}
	
	if !config.RequireExplicit {
		t.Error("Expected RequireExplicit to be true")
	}
}

func TestValidateInsecureConfig(t *testing.T) {
	im := NewInsecureMode()

	tests := []struct {
		name        string
		config      *InsecureConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &InsecureConfig{
				Registry:  "test-registry",
				Operation: "push",
				Duration:  5 * time.Minute,
				Reason:    "testing",
			},
			expectError: false,
		},
		{
			name: "missing registry",
			config: &InsecureConfig{
				Operation: "push",
				Duration:  5 * time.Minute,
				Reason:    "testing",
			},
			expectError: true,
			errorMsg:    "registry must be specified",
		},
		{
			name: "missing operation",
			config: &InsecureConfig{
				Registry: "test-registry",
				Duration: 5 * time.Minute,
				Reason:   "testing",
			},
			expectError: true,
			errorMsg:    "operation must be specified",
		},
		{
			name: "zero duration gets default",
			config: &InsecureConfig{
				Registry:  "test-registry",
				Operation: "push",
				Duration:  0,
				Reason:    "testing",
			},
			expectError: false,
		},
		{
			name: "excessive duration",
			config: &InsecureConfig{
				Registry:  "test-registry",
				Operation: "push",
				Duration:  48 * time.Hour,
				Reason:    "testing",
			},
			expectError: true,
			errorMsg:    "duration cannot exceed 24 hours",
		},
		{
			name: "empty reason gets default",
			config: &InsecureConfig{
				Registry:  "test-registry",
				Operation: "push",
				Duration:  5 * time.Minute,
				Reason:    "",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := im.ValidateInsecureConfig(tt.config)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				
				// Check defaults are set
				if tt.config.Duration == 0 {
					// Should have been set to default
					if tt.config.Duration != 5*time.Minute {
						t.Error("Expected default duration to be set")
					}
				}
				
				if tt.config.Reason == "" {
					// Should have been set to default
					if tt.config.Reason != "No reason provided" {
						t.Error("Expected default reason to be set")
					}
				}
			}
		})
	}
}

func TestGenerateSecurityWarning(t *testing.T) {
	im := NewInsecureMode()
	
	config := &InsecureConfig{
		Registry:  "test-registry",
		Operation: "push",
		Duration:  5 * time.Minute,
		Reason:    "testing",
	}
	
	warning := im.GenerateSecurityWarning(config)
	
	if !strings.Contains(warning, config.Registry) {
		t.Error("Warning should contain registry name")
	}
	
	if !strings.Contains(warning, config.Operation) {
		t.Error("Warning should contain operation")
	}
	
	if !strings.Contains(warning, config.Reason) {
		t.Error("Warning should contain reason")
	}
	
	if !strings.Contains(warning, "SECURITY WARNING") {
		t.Error("Warning should contain security warning text")
	}
}

func TestIsInsecureAllowed(t *testing.T) {
	im := NewInsecureMode()

	tests := []struct {
		registry string
		expected bool
	}{
		{"localhost:5000", true},
		{"127.0.0.1:5000", true},
		{"kind-registry:5000", true},
		{"gitea-http:3000", true},
		{"registry.local", true},
		{"registry.dev", true},
		{"production.example.com", false},
		{"registry.io", false},
		{"docker.io", false},
	}

	for _, tt := range tests {
		t.Run(tt.registry, func(t *testing.T) {
			result := im.IsInsecureAllowed(tt.registry)
			if result != tt.expected {
				t.Errorf("IsInsecureAllowed(%s) = %v, expected %v", tt.registry, result, tt.expected)
			}
		})
	}
}

func TestGetInsecureRecommendations(t *testing.T) {
	im := NewInsecureMode()
	
	registry := "test-registry"
	recommendations := im.GetInsecureRecommendations(registry)
	
	if len(recommendations) == 0 {
		t.Fatal("Expected at least one recommendation")
	}
	
	// Check that recommendations contain the registry name
	foundRegistryRef := false
	for _, rec := range recommendations {
		if strings.Contains(rec.Command, registry) || strings.Contains(rec.Description, registry) {
			foundRegistryRef = true
			break
		}
	}
	
	if !foundRegistryRef {
		t.Error("Expected at least one recommendation to reference the registry")
	}
	
	// Check priority ordering (high priority recommendations should be first)
	highPriorityFound := false
	for _, rec := range recommendations {
		if rec.Priority == PriorityHigh || rec.Priority == PriorityCritical {
			highPriorityFound = true
			break
		}
	}
	
	if !highPriorityFound {
		t.Error("Expected at least one high priority recommendation")
	}
}

func TestWarningSuppressionPerRegistry(t *testing.T) {
	im := NewInsecureMode()
	ctx := context.Background()
	
	config1 := &InsecureConfig{
		Registry:        "registry1",
		Operation:       "push",
		Duration:        5 * time.Minute,
		Reason:          "testing",
		RequireExplicit: false,
	}
	
	config2 := &InsecureConfig{
		Registry:        "registry2",
		Operation:       "push",
		Duration:        5 * time.Minute,
		Reason:          "testing",
		RequireExplicit: false,
	}
	
	// First call to registry1 should issue warning
	err1 := im.ApplyInsecureMode(ctx, config1)
	if err1 != nil {
		t.Fatalf("First ApplyInsecureMode failed: %v", err1)
	}
	
	// Second call to same registry should not issue warning again
	err2 := im.ApplyInsecureMode(ctx, config1)
	if err2 != nil {
		t.Fatalf("Second ApplyInsecureMode failed: %v", err2)
	}
	
	// Call to different registry should issue warning
	err3 := im.ApplyInsecureMode(ctx, config2)
	if err3 != nil {
		t.Fatalf("Third ApplyInsecureMode failed: %v", err3)
	}
	
	// Verify warnings were tracked
	key1 := "registry1:push"
	key2 := "registry2:push"
	
	if !im.warningsIssued[key1] {
		t.Error("Expected warning to be tracked for registry1")
	}
	
	if !im.warningsIssued[key2] {
		t.Error("Expected warning to be tracked for registry2")
	}
}