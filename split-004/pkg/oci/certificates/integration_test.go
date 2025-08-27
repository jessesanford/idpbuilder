package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestIntegrationConfigLoaderComplete tests complete ConfigLoader implementation.
func TestIntegrationConfigLoaderComplete(t *testing.T) {
	ctx := context.Background()
	
	tempDir, err := ioutil.TempDir("", "config_integration_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	loader := NewConfigLoader()

	// Test GetDefaultConfig
	defaultConfig, err := loader.GetDefaultConfig(ctx)
	if err != nil {
		t.Fatalf("failed to get default config: %v", err)
	}
	if defaultConfig.StorageBackend == nil {
		t.Error("default config should have storage backend")
	}

	// Test ValidateConfig
	result, err := loader.ValidateConfig(ctx, defaultConfig)
	if err != nil {
		t.Fatalf("failed to validate config: %v", err)
	}
	if !result.Valid {
		t.Errorf("default config should be valid: %v", result.Errors)
	}

	// Test SaveConfig and LoadConfig
	configPath := filepath.Join(tempDir, "test-config.yaml")
	err = loader.SaveConfig(ctx, defaultConfig, configPath)
	if err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	loadedConfig, err := loader.LoadConfig(ctx, configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if loadedConfig.StorageBackend.Type != defaultConfig.StorageBackend.Type {
		t.Error("loaded config storage type mismatch")
	}
}

// TestIntegrationConfigurationFormats tests YAML and JSON configuration loading.
func TestIntegrationConfigurationFormats(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "format_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test YAML config
	yamlConfig := `
trustStore:
  storagePath: ` + tempDir + `
  pools:
    system:
      enabled: true
      path: /etc/ssl/certs
  validation:
    checkExpiry: true
    checkChain: false
  events:
    enabled: false
    handlers: []
`
	
	yamlPath := filepath.Join(tempDir, "test.yaml")
	err = ioutil.WriteFile(yamlPath, []byte(yamlConfig), 0644)
	if err != nil {
		t.Fatalf("failed to write YAML config: %v", err)
	}

	config1, err := LoadYAMLConfig(yamlPath)
	if err != nil {
		t.Fatalf("failed to load YAML config: %v", err)
	}

	// Test JSON config
	jsonPath := filepath.Join(tempDir, "test.json")
	loader := NewConfigLoader()
	err = loader.SaveConfig(context.Background(), config1, jsonPath)
	if err != nil {
		t.Fatalf("failed to save JSON config: %v", err)
	}

	config2, err := LoadJSONConfig(jsonPath)
	if err != nil {
		t.Fatalf("failed to load JSON config: %v", err)
	}

	if config1.StorageBackend != nil && config2.StorageBackend != nil &&
		config1.StorageBackend.ConnectionString != config2.StorageBackend.ConnectionString {
		t.Error("YAML and JSON configs should have same storage path")
	}
}

// TestIntegrationConfigMerging tests configuration merging functionality.
func TestIntegrationConfigMerging(t *testing.T) {
	// Base configuration
	base := &Config{
		StorageBackend: &StorageConfig{
			Type:             "filesystem",
			ConnectionString: "/base/path",
		},
		DefaultPools: []string{"system"},
		ValidationRules: []*ValidationRule{
			{Name: "base_rule", Enabled: true},
		},
	}

	// Override configuration
	override := &Config{
		StorageBackend: &StorageConfig{
			Type:             "database",
			ConnectionString: "/override/path",
		},
		DefaultPools: []string{"custom"},
		ValidationRules: []*ValidationRule{
			{Name: "override_rule", Enabled: false},
		},
		EventHandlers: &EventConfig{
			Enabled: true,
		},
	}

	merged := MergeConfigs(base, override)

	// Verify override took precedence
	if merged.StorageBackend.Type != "database" {
		t.Error("merged config should use override storage type")
	}
	if merged.StorageBackend.ConnectionString != "/override/path" {
		t.Error("merged config should use override storage path")
	}
	if len(merged.DefaultPools) != 1 || merged.DefaultPools[0] != "custom" {
		t.Error("merged config should use override default pools")
	}
	if merged.EventHandlers == nil || !merged.EventHandlers.Enabled {
		t.Error("merged config should include override event handlers")
	}
}

// TestIntegrationEnvironmentOverrides tests environment variable overrides.
func TestIntegrationEnvironmentOverrides(t *testing.T) {
	// Set environment variables
	os.Setenv("TRUSTSTORE_STORAGE_PATH", "/env/override/path")
	os.Setenv("TRUSTSTORE_EVENTS_ENABLED", "true")
	defer func() {
		os.Unsetenv("TRUSTSTORE_STORAGE_PATH")
		os.Unsetenv("TRUSTSTORE_EVENTS_ENABLED")
	}()

	config := &Config{
		StorageBackend: &StorageConfig{
			Type:             "filesystem",
			ConnectionString: "/original/path",
		},
		EventHandlers: &EventConfig{
			Enabled: false,
		},
	}

	err := ApplyEnvironmentOverrides(config)
	if err != nil {
		t.Fatalf("failed to apply environment overrides: %v", err)
	}

	if config.StorageBackend.ConnectionString != "/env/override/path" {
		t.Errorf("expected storage path to be overridden, got: %s", config.StorageBackend.ConnectionString)
	}
	if !config.EventHandlers.Enabled {
		t.Error("expected events to be enabled via environment override")
	}
}

// generateTestCertificate generates a test certificate.
func generateTestCertificate(t *testing.T, commonName string, notAfter time.Time) (*x509.Certificate, []byte) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("failed to parse certificate: %v", err)
	}

	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}

	pemData := pem.EncodeToMemory(pemBlock)
	return cert, pemData
}