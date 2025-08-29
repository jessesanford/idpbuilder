package certs

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInsecureMode(t *testing.T) {
	im := NewInsecureMode()
	if im == nil || im.logger == nil || im.warningsIssued == nil {
		t.Fatal("NewInsecureMode initialization failed")
	}

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test flag detection
	tests := []struct {
		args     []string
		expected bool
	}{
		{[]string{"cmd"}, false},
		{[]string{"cmd", "--insecure"}, true},
		{[]string{"cmd", "-k"}, true},
	}

	for _, tt := range tests {
		os.Args = tt.args
		if im.IsInsecureModeRequested() != tt.expected {
			t.Errorf("Flag detection failed for args: %v", tt.args)
		}
	}

	// Test insecure mode application
	ctx := context.Background()
	config := &InsecureConfig{
		Registry: "test", Operation: "push", Duration: 5 * time.Minute, Reason: "testing",
	}
	
	if err := im.ApplyInsecureMode(ctx, config); err != nil {
		t.Error("ApplyInsecureMode failed")
	}

	// Test explicit flag requirement
	config.RequireExplicit = true
	os.Args = []string{"cmd"}
	if err := im.ApplyInsecureMode(ctx, config); err == nil {
		t.Error("Expected error when explicit flag required but not present")
	}

	os.Args = []string{"cmd", "--insecure"}
	if err := im.ApplyInsecureMode(ctx, config); err != nil {
		t.Error("ApplyInsecureMode failed with explicit flag")
	}
}

func TestInsecureConfig(t *testing.T) {
	im := NewInsecureMode()
	
	// Test config creation
	config := im.CreateTimeLimitedConfig("test-reg", "push", "testing", 10*time.Minute)
	if config.Registry != "test-reg" || !config.RequireExplicit {
		t.Error("CreateTimeLimitedConfig failed")
	}

	// Test validation
	tests := []struct {
		config      *InsecureConfig
		expectError bool
		errorMsg    string
	}{
		{&InsecureConfig{Registry: "test", Operation: "push", Duration: 5 * time.Minute}, false, ""},
		{&InsecureConfig{Operation: "push"}, true, "registry must be specified"},
		{&InsecureConfig{Registry: "test"}, true, "operation must be specified"},
		{&InsecureConfig{Registry: "test", Operation: "push", Duration: 48 * time.Hour}, true, "duration cannot exceed 24 hours"},
	}

	for _, tt := range tests {
		err := im.ValidateInsecureConfig(tt.config)
		if (err != nil) != tt.expectError {
			t.Errorf("Validation test failed: %v", err)
		}
		if tt.expectError && !strings.Contains(err.Error(), tt.errorMsg) {
			t.Errorf("Expected error message %q, got %q", tt.errorMsg, err.Error())
		}
	}
}

func TestInsecureAllowedRegistries(t *testing.T) {
	im := NewInsecureMode()
	
	allowedTests := []struct {
		registry string
		allowed  bool
	}{
		{"localhost:5000", true},
		{"kind-registry", true},
		{"registry.dev", true},
		{"production.example.com", false},
	}

	for _, tt := range allowedTests {
		if im.IsInsecureAllowed(tt.registry) != tt.allowed {
			t.Errorf("IsInsecureAllowed(%s) should be %v", tt.registry, tt.allowed)
		}
	}

	// Test recommendations
	recs := im.GetInsecureRecommendations("test-registry")
	if len(recs) == 0 {
		t.Error("Expected recommendations")
	}
	
	hasHighPrio := false
	for _, rec := range recs {
		if rec.Priority == PriorityHigh {
			hasHighPrio = true
		}
	}
	if !hasHighPrio {
		t.Error("Expected high priority recommendation")
	}

	// Test warning generation
	config := &InsecureConfig{Registry: "test", Operation: "push", Duration: 5 * time.Minute, Reason: "testing"}
	warning := im.GenerateSecurityWarning(config)
	if !strings.Contains(warning, "SECURITY WARNING") || !strings.Contains(warning, "test") {
		t.Error("Security warning generation failed")
	}
}