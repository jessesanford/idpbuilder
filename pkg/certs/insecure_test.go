package certs

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"
)

func TestDefaultInsecureConfig(t *testing.T) {
	config := DefaultInsecureConfig()
	
	if config == nil {
		t.Fatal("DefaultInsecureConfig returned nil")
	}
	
	// Verify secure defaults
	if config.Enabled {
		t.Error("Expected insecure mode to be disabled by default")
	}
	
	if !config.LogWarnings {
		t.Error("Expected warning logging to be enabled by default")
	}
	
	if config.AllowProduction {
		t.Error("Expected production mode to be disallowed by default")
	}
	
	if !config.RequireConsent {
		t.Error("Expected consent to be required by default")
	}
	
	if config.MaxWarnings != 10 {
		t.Errorf("Expected MaxWarnings to be 10, got %d", config.MaxWarnings)
	}
	
	if len(config.AllowedRegistries) != 0 {
		t.Error("Expected empty allowlist by default")
	}
}

func TestNewInsecureMode(t *testing.T) {
	// Test with nil config (should use defaults)
	mode := NewInsecureMode(nil)
	if mode == nil {
		t.Fatal("NewInsecureMode returned nil")
	}
	
	if mode.IsEnabled() {
		t.Error("Expected insecure mode to be disabled with default config")
	}
	
	// Test with enabled config
	config := &InsecureConfig{
		Enabled:     true,
		LogWarnings: true,
	}
	
	mode = NewInsecureMode(config)
	if mode == nil {
		t.Fatal("NewInsecureMode with config returned nil")
	}
}

func TestInsecureModeEnabled(t *testing.T) {
	testCases := []struct {
		name     string
		config   *InsecureConfig
		expected bool
	}{
		{
			name: "Disabled by default",
			config: &InsecureConfig{
				Enabled: false,
			},
			expected: false,
		},
		{
			name: "Explicitly enabled",
			config: &InsecureConfig{
				Enabled: true,
			},
			expected: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mode := NewInsecureMode(tc.config)
			if mode.IsEnabled() != tc.expected {
				t.Errorf("Expected IsEnabled() to be %t, got %t", tc.expected, mode.IsEnabled())
			}
		})
	}
}

func TestShouldBypass(t *testing.T) {
	config := &InsecureConfig{
		Enabled:     true,
		LogWarnings: false, // Disable to avoid log noise in tests
	}
	mode := NewInsecureMode(config)
	
	testCases := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "Nil error should not bypass",
			error:    nil,
			expected: false,
		},
		{
			name:     "Self-signed certificate should bypass",
			error:    errors.New("self signed certificate"),
			expected: true,
		},
		{
			name:     "Unknown authority should bypass",
			error:    errors.New("certificate signed by unknown authority"),
			expected: true,
		},
		{
			name:     "Expired certificate should bypass",
			error:    errors.New("certificate has expired"),
			expected: true,
		},
		{
			name:     "Name mismatch should bypass",
			error:    errors.New("certificate name does not match"),
			expected: true,
		},
		{
			name:     "TLS verify failed should bypass",
			error:    errors.New("tls: failed to verify certificate"),
			expected: true,
		},
		{
			name:     "Non-certificate error should not bypass",
			error:    errors.New("network connection refused"),
			expected: false,
		},
		{
			name:     "Random error should not bypass",
			error:    errors.New("some random error"),
			expected: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := mode.ShouldBypass(tc.error)
			if result != tc.expected {
				t.Errorf("Expected ShouldBypass(%v) to be %t, got %t", tc.error, tc.expected, result)
			}
		})
	}
}

func TestShouldBypassWhenDisabled(t *testing.T) {
	config := &InsecureConfig{
		Enabled: false, // Disabled
	}
	mode := NewInsecureMode(config)
	
	// Even bypassable errors should not bypass when disabled
	err := errors.New("self signed certificate")
	if mode.ShouldBypass(err) {
		t.Error("ShouldBypass should return false when insecure mode is disabled")
	}
}

func TestLogSecurityWarning(t *testing.T) {
	config := &InsecureConfig{
		Enabled:          true,
		LogWarnings:      true,
		MaxWarnings:      2,
		WarningThreshold: 100 * time.Millisecond,
	}
	
	mode := NewInsecureMode(config)
	
	// Log a few warnings
	mode.LogSecurityWarning("Test warning 1")
	mode.LogSecurityWarning("Test warning 2")
	
	// This should be throttled
	mode.LogSecurityWarning("Test warning 3 - should be throttled")
	
	// Wait for threshold and try again
	time.Sleep(150 * time.Millisecond)
	mode.LogSecurityWarning("Test warning 4 - after threshold")
}

func TestRequiresConsent(t *testing.T) {
	testCases := []struct {
		name           string
		requireConsent bool
		expected       bool
	}{
		{
			name:           "Consent required by default",
			requireConsent: true,
			expected:       true,
		},
		{
			name:           "Consent not required",
			requireConsent: false,
			expected:       false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &InsecureConfig{
				Enabled:        true,
				RequireConsent: tc.requireConsent,
			}
			mode := NewInsecureMode(config)
			
			if mode.RequiresConsent() != tc.expected {
				t.Errorf("Expected RequiresConsent() to be %t, got %t", tc.expected, mode.RequiresConsent())
			}
		})
	}
}

func TestGetWarningMessage(t *testing.T) {
	mode := NewInsecureMode(&InsecureConfig{Enabled: true})
	
	message := mode.GetWarningMessage("push")
	
	if !strings.Contains(message, "SECURITY WARNING") {
		t.Error("Warning message should contain SECURITY WARNING")
	}
	
	if !strings.Contains(message, "push") {
		t.Error("Warning message should contain the operation")
	}
	
	if !strings.Contains(message, "Man-in-the-middle") {
		t.Error("Warning message should mention security risks")
	}
}

func TestDetectEnvironment(t *testing.T) {
	// Save original environment
	originalEnv := os.Getenv("GO_ENV")
	defer func() {
		if originalEnv != "" {
			os.Setenv("GO_ENV", originalEnv)
		} else {
			os.Unsetenv("GO_ENV")
		}
	}()
	
	testCases := []struct {
		name     string
		envVar   string
		envValue string
		expected string
	}{
		{
			name:     "Production environment",
			envVar:   "GO_ENV",
			envValue: "production",
			expected: "production",
		},
		{
			name:     "Development environment", 
			envVar:   "GO_ENV",
			envValue: "development",
			expected: "development",
		},
		{
			name:     "Staging environment",
			envVar:   "NODE_ENV", 
			envValue: "staging",
			expected: "staging",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear environment first
			os.Unsetenv("GO_ENV")
			os.Unsetenv("NODE_ENV")
			os.Unsetenv("ENVIRONMENT")
			os.Unsetenv("ENV")
			
			// Set test environment
			os.Setenv(tc.envVar, tc.envValue)
			
			env := detectEnvironment()
			if env != tc.expected {
				t.Errorf("Expected environment %s, got %s", tc.expected, env)
			}
		})
	}
}

func TestNewInsecureStrategy(t *testing.T) {
	// Test with nil mode (should create default)
	strategy := NewInsecureStrategy(nil)
	if strategy == nil {
		t.Fatal("NewInsecureStrategy returned nil")
	}
	
	if strategy.Name() != "insecure-bypass" {
		t.Errorf("Expected strategy name 'insecure-bypass', got %s", strategy.Name())
	}
	
	if strategy.Priority() != 15 {
		t.Errorf("Expected priority 15 (higher than tertiary), got %d", strategy.Priority())
	}
	
	// Test with custom mode
	config := &InsecureConfig{Enabled: true}
	mode := NewInsecureMode(config)
	strategy = NewInsecureStrategy(mode)
	
	if strategy == nil {
		t.Fatal("NewInsecureStrategy with custom mode returned nil")
	}
}

func TestInsecureStrategyCanHandle(t *testing.T) {
	config := &InsecureConfig{
		Enabled:     true,
		LogWarnings: false, // Disable to avoid log noise
	}
	mode := NewInsecureMode(config)
	strategy := NewInsecureStrategy(mode)
	
	testCases := []struct {
		name     string
		error    error
		expected bool
	}{
		{
			name:     "Can handle self-signed certificate",
			error:    errors.New("self signed certificate"),
			expected: true,
		},
		{
			name:     "Can handle unknown authority",
			error:    errors.New("certificate signed by unknown authority"),
			expected: true,
		},
		{
			name:     "Cannot handle nil error",
			error:    nil,
			expected: false,
		},
		{
			name:     "Cannot handle non-certificate error",
			error:    errors.New("network error"),
			expected: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := strategy.CanHandle(tc.error)
			if result != tc.expected {
				t.Errorf("Expected CanHandle(%v) to be %t, got %t", tc.error, tc.expected, result)
			}
		})
	}
}

func TestInsecureStrategyExecute(t *testing.T) {
	config := &InsecureConfig{
		Enabled:        true,
		RequireConsent: false, // Disable consent for easier testing
		LogWarnings:    false, // Disable to avoid log noise
	}
	mode := NewInsecureMode(config)
	strategy := NewInsecureStrategy(mode)
	ctx := context.Background()
	
	t.Run("Execute with nil error", func(t *testing.T) {
		input := &ValidationInput{
			Registry:  "test.registry.com",
			Operation: "pull",
			Error:     nil,
		}
		
		result, err := strategy.Execute(ctx, input)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		
		if !result.Success {
			t.Error("Expected success for nil error")
		}
		
		if result.SecurityLevel != SecurityNone {
			t.Errorf("Expected SecurityNone, got %v", result.SecurityLevel)
		}
	})
	
	t.Run("Execute with bypassable error", func(t *testing.T) {
		input := &ValidationInput{
			Registry:  "test.registry.com",
			Operation: "pull", 
			Error:     errors.New("self signed certificate"),
		}
		
		result, err := strategy.Execute(ctx, input)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		
		if !result.Success {
			t.Error("Expected success for bypassable error")
		}
		
		if result.SecurityLevel != SecurityNone {
			t.Error("Expected SecurityNone for bypassed certificate")
		}
		
		if result.Strategy != "insecure-bypass" {
			t.Errorf("Expected strategy 'insecure-bypass', got %s", result.Strategy)
		}
		
		// Check that NewConfig contains expected fields
		newConfig := result.NewConfig
		if newConfig["insecure_mode"] != true {
			t.Error("Expected insecure_mode to be true in NewConfig")
		}
		
		if newConfig["bypass_certificates"] != true {
			t.Error("Expected bypass_certificates to be true in NewConfig")
		}
	})
	
	t.Run("Execute with non-bypassable error", func(t *testing.T) {
		input := &ValidationInput{
			Registry:  "test.registry.com",
			Operation: "pull",
			Error:     errors.New("network connection refused"),
		}
		
		result, err := strategy.Execute(ctx, input)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}
		
		if result.Success {
			t.Error("Expected failure for non-bypassable error")
		}
	})
}

func TestInsecureStrategyExecuteWithConsent(t *testing.T) {
	config := &InsecureConfig{
		Enabled:        true,
		RequireConsent: true, // Consent required
		LogWarnings:    false,
	}
	mode := NewInsecureMode(config)
	strategy := NewInsecureStrategy(mode)
	ctx := context.Background()
	
	input := &ValidationInput{
		Registry:  "test.registry.com",
		Operation: "pull",
		Error:     errors.New("self signed certificate"),
	}
	
	result, err := strategy.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	// Should fail because consent is required but not given
	if result.Success {
		t.Error("Expected failure when consent is required but not given")
	}
	
	if !strings.Contains(result.Message, "consent required") {
		t.Error("Expected message to mention consent requirement")
	}
}

func TestInsecureStrategyDisabled(t *testing.T) {
	config := &InsecureConfig{
		Enabled: false, // Disabled
	}
	mode := NewInsecureMode(config)
	strategy := NewInsecureStrategy(mode)
	ctx := context.Background()
	
	input := &ValidationInput{
		Registry:  "test.registry.com",
		Operation: "pull",
		Error:     errors.New("self signed certificate"),
	}
	
	result, err := strategy.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	if result.Success {
		t.Error("Expected failure when insecure mode is disabled")
	}
	
	if !strings.Contains(result.Message, "disabled") {
		t.Error("Expected message to mention that insecure mode is disabled")
	}
}

func TestIsRegistryAllowed(t *testing.T) {
	config := &InsecureConfig{
		Enabled: true,
		AllowedRegistries: []string{
			"localhost:5000",
			"*.test.com",
			"internal.registry.dev",
		},
	}
	mode := NewInsecureMode(config)
	strategy := NewInsecureStrategy(mode).(*InsecureStrategy)
	
	testCases := []struct {
		name     string
		registry string
		expected bool
	}{
		{
			name:     "Exact match allowed",
			registry: "localhost:5000",
			expected: true,
		},
		{
			name:     "Wildcard match allowed",
			registry: "sub.test.com",
			expected: true,
		},
		{
			name:     "Another wildcard match",
			registry: "api.test.com",
			expected: true,
		},
		{
			name:     "Internal registry allowed",
			registry: "internal.registry.dev",
			expected: true,
		},
		{
			name:     "Public registry not allowed",
			registry: "docker.io",
			expected: false,
		},
		{
			name:     "Random registry not allowed",
			registry: "malicious.registry.com",
			expected: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := strategy.IsRegistryAllowed(tc.registry)
			if result != tc.expected {
				t.Errorf("Expected IsRegistryAllowed(%s) to be %t, got %t", tc.registry, tc.expected, result)
			}
		})
	}
}

func TestInsecureStrategyEmptyAllowlist(t *testing.T) {
	config := &InsecureConfig{
		Enabled:           true,
		AllowedRegistries: []string{}, // Empty allowlist
	}
	mode := NewInsecureMode(config)
	strategy := NewInsecureStrategy(mode).(*InsecureStrategy)
	
	// With empty allowlist, all registries should be allowed (risky but explicit)
	if !strategy.IsRegistryAllowed("docker.io") {
		t.Error("Expected all registries to be allowed with empty allowlist")
	}
	
	if !strategy.IsRegistryAllowed("any.registry.com") {
		t.Error("Expected all registries to be allowed with empty allowlist")
	}
}