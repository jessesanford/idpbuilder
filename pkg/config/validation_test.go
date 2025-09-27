package config

import (
	"os"
	"testing"
)

func TestValidateRegistryConfig(t *testing.T) {
	// Test nil config
	result := ValidateRegistryConfig(nil)
	if result.Valid {
		t.Error("Expected invalid result for nil config")
	}
	if len(result.Errors) == 0 {
		t.Error("Expected errors for nil config")
	}

	// Test valid config
	config := &RegistryConfig{
		URL:  "https://registry.example.com",
		Type: "harbor",
		Auth: AuthConfig{
			Type:     "basic",
			Username: "testuser",
			Password: "testpass",
		},
	}
	result = ValidateRegistryConfig(config)
	if !result.Valid {
		t.Errorf("Expected valid config, got errors: %v", result.Errors)
	}

	// Test missing URL
	config.URL = ""
	result = ValidateRegistryConfig(config)
	if result.Valid {
		t.Error("Expected invalid result for missing URL")
	}
	found := false
	for _, err := range result.Errors {
		if err == "URL is required" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'URL is required' error")
	}

	// Test invalid URL
	config.URL = "ht tp://invalid url with spaces"
	result = ValidateRegistryConfig(config)
	if result.Valid {
		t.Error("Expected invalid result for malformed URL")
	}

	// Test unknown registry type (should be warning, not error)
	config.URL = "https://registry.example.com"
	config.Type = "unknown"
	result = ValidateRegistryConfig(config)
	if !result.Valid {
		t.Errorf("Expected valid config with unknown type, got errors: %v", result.Errors)
	}
	if len(result.Warnings) == 0 {
		t.Error("Expected warning for unknown registry type")
	}

	// Test insecure connection warning
	config.Type = "harbor"
	config.Insecure = true
	result = ValidateRegistryConfig(config)
	if !result.Valid {
		t.Errorf("Expected valid config with insecure, got errors: %v", result.Errors)
	}
	warningFound := false
	for _, warning := range result.Warnings {
		if warning == "insecure connections enabled - not recommended for production" {
			warningFound = true
			break
		}
	}
	if !warningFound {
		t.Error("Expected insecure connection warning")
	}

	// Test invalid auth configs
	config.Insecure = false
	config.Auth = AuthConfig{
		Type: "basic",
		// Missing username and password
	}
	result = ValidateRegistryConfig(config)
	if result.Valid {
		t.Error("Expected invalid result for incomplete basic auth")
	}

	config.Auth = AuthConfig{
		Type: "token",
		// Missing token
	}
	result = ValidateRegistryConfig(config)
	if result.Valid {
		t.Error("Expected invalid result for missing token")
	}

	config.Auth = AuthConfig{
		Type: "invalid",
	}
	result = ValidateRegistryConfig(config)
	if result.Valid {
		t.Error("Expected invalid result for unsupported auth type")
	}
}

func TestSetDefaultValues(t *testing.T) {
	// Test nil config
	SetDefaultValues(nil)
	// Should not panic

	// Test setting defaults
	config := &RegistryConfig{}
	SetDefaultValues(config)

	if config.Type != "generic" {
		t.Errorf("Expected default type 'generic', got '%s'", config.Type)
	}
	if config.URL != "docker.io" {
		t.Errorf("Expected default URL 'docker.io', got '%s'", config.URL)
	}
	if config.Options == nil {
		t.Error("Expected options map to be initialized")
	}
	if config.Options["timeout"] != "30s" {
		t.Errorf("Expected default timeout '30s', got '%s'", config.Options["timeout"])
	}

	// Test environment variable overrides
	os.Setenv("REGISTRY_URL", "env.registry.com")
	os.Setenv("REGISTRY_TIMEOUT", "60s")
	defer func() {
		os.Unsetenv("REGISTRY_URL")
		os.Unsetenv("REGISTRY_TIMEOUT")
	}()

	config2 := &RegistryConfig{}
	SetDefaultValues(config2)

	if config2.URL != "env.registry.com" {
		t.Errorf("Expected env URL 'env.registry.com', got '%s'", config2.URL)
	}
	if config2.Options["timeout"] != "60s" {
		t.Errorf("Expected env timeout '60s', got '%s'", config2.Options["timeout"])
	}

	// Test preserving existing values
	config3 := &RegistryConfig{
		URL:  "existing.registry.com",
		Type: "harbor",
		Options: map[string]string{
			"timeout": "120s",
		},
	}
	SetDefaultValues(config3)

	if config3.URL != "existing.registry.com" {
		t.Errorf("Expected preserved URL 'existing.registry.com', got '%s'", config3.URL)
	}
	if config3.Type != "harbor" {
		t.Errorf("Expected preserved type 'harbor', got '%s'", config3.Type)
	}
	if config3.Options["timeout"] != "120s" {
		t.Errorf("Expected preserved timeout '120s', got '%s'", config3.Options["timeout"])
	}
}

func TestMergeConfigs(t *testing.T) {
	// Test both nil
	result := MergeConfigs(nil, nil)
	if result != nil {
		t.Error("Expected nil result for both nil configs")
	}

	// Test base nil
	override := &RegistryConfig{URL: "override.com"}
	result = MergeConfigs(nil, override)
	if result == nil || result.URL != "override.com" {
		t.Error("Expected override config when base is nil")
	}

	// Test override nil
	base := &RegistryConfig{URL: "base.com"}
	result = MergeConfigs(base, nil)
	if result == nil || result.URL != "base.com" {
		t.Error("Expected base config when override is nil")
	}

	// Test normal merge
	base = &RegistryConfig{
		URL:      "base.registry.com",
		Type:     "harbor",
		Insecure: false,
		Auth: AuthConfig{
			Type:     "basic",
			Username: "baseuser",
		},
		Options: map[string]string{
			"timeout":        "30s",
			"max_connections": "5",
		},
	}

	override = &RegistryConfig{
		URL:      "override.registry.com",
		Insecure: true,
		Auth: AuthConfig{
			Password: "overridepass",
		},
		Options: map[string]string{
			"timeout": "60s",
			"retries": "3",
		},
	}

	result = MergeConfigs(base, override)

	if result.URL != "override.registry.com" {
		t.Errorf("Expected override URL 'override.registry.com', got '%s'", result.URL)
	}
	if result.Type != "harbor" {
		t.Errorf("Expected base type 'harbor', got '%s'", result.Type)
	}
	if result.Insecure != true {
		t.Errorf("Expected override insecure true, got %v", result.Insecure)
	}
	if result.Auth.Username != "baseuser" {
		t.Errorf("Expected base username 'baseuser', got '%s'", result.Auth.Username)
	}
	if result.Auth.Password != "overridepass" {
		t.Errorf("Expected override password 'overridepass', got '%s'", result.Auth.Password)
	}
	if result.Options["timeout"] != "60s" {
		t.Errorf("Expected override timeout '60s', got '%s'", result.Options["timeout"])
	}
	if result.Options["max_connections"] != "5" {
		t.Errorf("Expected base max_connections '5', got '%s'", result.Options["max_connections"])
	}
	if result.Options["retries"] != "3" {
		t.Errorf("Expected override retries '3', got '%s'", result.Options["retries"])
	}

	// Test auth type override
	override.Auth.Type = "token"
	override.Auth.Token = "abc123"
	result = MergeConfigs(base, override)

	if result.Auth.Type != "token" {
		t.Errorf("Expected override auth type 'token', got '%s'", result.Auth.Type)
	}
	if result.Auth.Token != "abc123" {
		t.Errorf("Expected override token 'abc123', got '%s'", result.Auth.Token)
	}
}