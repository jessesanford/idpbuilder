package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadRegistryConfig(t *testing.T) {
	// Test YAML config loading
	yamlContent := `
url: "registry.example.com"
type: "harbor"
insecure: false
auth:
  type: "basic"
  username: "testuser"
  password: "testpass"
options:
  timeout: "60s"
`
	tmpDir, err := ioutil.TempDir("", "registry_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	yamlFile := filepath.Join(tmpDir, "config.yaml")
	if err := ioutil.WriteFile(yamlFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write YAML file: %v", err)
	}

	config, err := LoadRegistryConfig(yamlFile)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if config.URL != "registry.example.com" {
		t.Errorf("Expected URL 'registry.example.com', got '%s'", config.URL)
	}
	if config.Type != "harbor" {
		t.Errorf("Expected type 'harbor', got '%s'", config.Type)
	}
	if config.Insecure != false {
		t.Errorf("Expected insecure false, got %v", config.Insecure)
	}

	// Test JSON config loading
	jsonContent := `{
		"url": "registry.json.com",
		"type": "dockerhub",
		"insecure": true
	}`
	jsonFile := filepath.Join(tmpDir, "config.json")
	if err := ioutil.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write JSON file: %v", err)
	}

	jsonConfig, err := LoadRegistryConfig(jsonFile)
	if err != nil {
		t.Fatalf("Failed to load JSON config: %v", err)
	}

	if jsonConfig.URL != "registry.json.com" {
		t.Errorf("Expected URL 'registry.json.com', got '%s'", jsonConfig.URL)
	}

	// Test environment variable override
	os.Setenv("REGISTRY_URL", "env.registry.com")
	defer os.Unsetenv("REGISTRY_URL")

	envConfig, err := LoadRegistryConfig(yamlFile)
	if err != nil {
		t.Fatalf("Failed to load config with env override: %v", err)
	}

	if envConfig.URL != "env.registry.com" {
		t.Errorf("Expected env URL 'env.registry.com', got '%s'", envConfig.URL)
	}

	// Test error cases
	if _, err := LoadRegistryConfig(""); err == nil {
		t.Error("Expected error for empty path")
	}

	if _, err := LoadRegistryConfig("nonexistent.yaml"); err == nil {
		t.Error("Expected error for nonexistent file")
	}

	unsupportedFile := filepath.Join(tmpDir, "config.txt")
	if err := ioutil.WriteFile(unsupportedFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write unsupported file: %v", err)
	}
	if _, err := LoadRegistryConfig(unsupportedFile); err == nil {
		t.Error("Expected error for unsupported file format")
	}
}

func TestToConnectionString(t *testing.T) {
	// Test HTTPS connection string
	config := &RegistryConfig{
		URL:      "registry.example.com",
		Insecure: false,
	}
	connStr := ToConnectionString(config)
	expected := "https://registry.example.com"
	if connStr != expected {
		t.Errorf("Expected '%s', got '%s'", expected, connStr)
	}

	// Test HTTP connection string
	config.Insecure = true
	connStr = ToConnectionString(config)
	expected = "http://registry.example.com"
	if connStr != expected {
		t.Errorf("Expected '%s', got '%s'", expected, connStr)
	}

	// Test nil config
	connStr = ToConnectionString(nil)
	if connStr != "" {
		t.Errorf("Expected empty string for nil config, got '%s'", connStr)
	}
}