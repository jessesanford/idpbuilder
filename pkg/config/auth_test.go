package config

import (
	"os"
	"testing"
)

func TestGetAuthConfig(t *testing.T) {
	// Test nil config
	_, err := GetAuthConfig(nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}

	// Test empty auth config (no auth)
	config := &RegistryConfig{
		URL: "registry.example.com",
		Auth: AuthConfig{},
	}
	authConfig, err := GetAuthConfig(config)
	if err != nil {
		t.Fatalf("Failed to get empty auth config: %v", err)
	}
	if authConfig.Type != "" {
		t.Errorf("Expected empty auth type, got '%s'", authConfig.Type)
	}

	// Test basic auth
	config.Auth = AuthConfig{
		Type:     "basic",
		Username: "testuser",
		Password: "testpass",
	}
	authConfig, err = GetAuthConfig(config)
	if err != nil {
		t.Fatalf("Failed to get basic auth config: %v", err)
	}
	if authConfig.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", authConfig.Username)
	}

	// Test token auth
	config.Auth = AuthConfig{
		Type:  "token",
		Token: "abc123",
	}
	authConfig, err = GetAuthConfig(config)
	if err != nil {
		t.Fatalf("Failed to get token auth config: %v", err)
	}
	if authConfig.Token != "abc123" {
		t.Errorf("Expected token 'abc123', got '%s'", authConfig.Token)
	}

	// Test OAuth2 auth
	config.Auth = AuthConfig{
		Type:  "oauth2",
		Token: "oauth-token-123",
	}
	authConfig, err = GetAuthConfig(config)
	if err != nil {
		t.Fatalf("Failed to get oauth2 auth config: %v", err)
	}
	if authConfig.Token != "oauth-token-123" {
		t.Errorf("Expected token 'oauth-token-123', got '%s'", authConfig.Token)
	}

	// Test environment variable resolution
	os.Setenv("REGISTRY_USERNAME", "envuser")
	os.Setenv("REGISTRY_PASSWORD", "envpass")
	defer func() {
		os.Unsetenv("REGISTRY_USERNAME")
		os.Unsetenv("REGISTRY_PASSWORD")
	}()

	config.Auth = AuthConfig{
		Type: "basic",
		// Username and Password left empty to test env resolution
	}
	authConfig, err = GetAuthConfig(config)
	if err != nil {
		t.Fatalf("Failed to get auth config with env vars: %v", err)
	}
	if authConfig.Username != "envuser" {
		t.Errorf("Expected env username 'envuser', got '%s'", authConfig.Username)
	}
	if authConfig.Password != "envpass" {
		t.Errorf("Expected env password 'envpass', got '%s'", authConfig.Password)
	}

	// Test validation errors
	config.Auth = AuthConfig{
		Type: "basic",
		// Missing username and password
	}
	os.Unsetenv("REGISTRY_USERNAME")
	os.Unsetenv("REGISTRY_PASSWORD")
	_, err = GetAuthConfig(config)
	if err == nil {
		t.Error("Expected error for basic auth without credentials")
	}

	config.Auth = AuthConfig{
		Type: "token",
		// Missing token
	}
	_, err = GetAuthConfig(config)
	if err == nil {
		t.Error("Expected error for token auth without token")
	}

	config.Auth = AuthConfig{
		Type: "unsupported",
	}
	_, err = GetAuthConfig(config)
	if err == nil {
		t.Error("Expected error for unsupported auth type")
	}
}

func TestTLSConfig(t *testing.T) {
	// Test TLS config structure
	tlsConfig := TLSConfig{
		CACert:     "/path/to/ca.crt",
		ClientCert: "/path/to/client.crt",
		ClientKey:  "/path/to/client.key",
		SkipVerify: false,
	}

	if tlsConfig.CACert != "/path/to/ca.crt" {
		t.Errorf("Expected CA cert path '/path/to/ca.crt', got '%s'", tlsConfig.CACert)
	}
	if tlsConfig.ClientCert != "/path/to/client.crt" {
		t.Errorf("Expected client cert path '/path/to/client.crt', got '%s'", tlsConfig.ClientCert)
	}
	if tlsConfig.ClientKey != "/path/to/client.key" {
		t.Errorf("Expected client key path '/path/to/client.key', got '%s'", tlsConfig.ClientKey)
	}
	if tlsConfig.SkipVerify != false {
		t.Errorf("Expected SkipVerify false, got %v", tlsConfig.SkipVerify)
	}

	// Test skip verify enabled
	tlsConfigSkip := TLSConfig{
		SkipVerify: true,
	}
	if tlsConfigSkip.SkipVerify != true {
		t.Errorf("Expected SkipVerify true, got %v", tlsConfigSkip.SkipVerify)
	}
}