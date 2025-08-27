package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestFilesystemStore_NewFilesystemStore tests the creation of a new filesystem store.
func TestFilesystemStore_NewFilesystemStore(t *testing.T) {
	tempDir := t.TempDir()

	store, err := NewFilesystemStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create filesystem store: %v", err)
	}
	defer store.Close()

	if store.basePath != tempDir {
		t.Errorf("Expected basePath %s, got %s", tempDir, store.basePath)
	}

	// Verify directory was created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Base directory was not created")
	}
}

// TestFilesystemStore_SaveAndLoad tests saving and loading certificates.
func TestFilesystemStore_SaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewFilesystemStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create filesystem store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	cert := createTestCertificate(t)

	// Save certificate
	if err := store.Save(ctx, "test-cert", cert); err != nil {
		t.Fatalf("Failed to save certificate: %v", err)
	}

	// Load certificate
	loaded, err := store.Load(ctx, "test-cert")
	if err != nil {
		t.Fatalf("Failed to load certificate: %v", err)
	}

	if loaded.ID != "test-cert" {
		t.Errorf("Expected ID 'test-cert', got '%s'", loaded.ID)
	}

	if string(loaded.Data) != string(cert.Data) {
		t.Error("Certificate data mismatch")
	}
}

// TestFilesystemStore_Delete tests certificate deletion.
func TestFilesystemStore_Delete(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewFilesystemStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create filesystem store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	cert := createTestCertificate(t)

	// Save certificate
	if err := store.Save(ctx, "test-cert", cert); err != nil {
		t.Fatalf("Failed to save certificate: %v", err)
	}

	// Delete certificate
	if err := store.Delete(ctx, "test-cert"); err != nil {
		t.Fatalf("Failed to delete certificate: %v", err)
	}

	// Verify certificate is gone
	_, err = store.Load(ctx, "test-cert")
	if err == nil {
		t.Error("Expected error when loading deleted certificate")
	}
}

// TestFilesystemStore_List tests listing certificates.
func TestFilesystemStore_List(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewFilesystemStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create filesystem store: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	cert1 := createTestCertificate(t)
	cert2 := createTestCertificate(t)

	// Save certificates
	if err := store.Save(ctx, "cert1", cert1); err != nil {
		t.Fatalf("Failed to save cert1: %v", err)
	}
	if err := store.Save(ctx, "cert2", cert2); err != nil {
		t.Fatalf("Failed to save cert2: %v", err)
	}

	// List certificates
	certIDs, err := store.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list certificates: %v", err)
	}

	if len(certIDs) != 2 {
		t.Errorf("Expected 2 certificates, got %d", len(certIDs))
	}

	expectedIDs := map[string]bool{"cert1": true, "cert2": true}
	for _, id := range certIDs {
		if !expectedIDs[id] {
			t.Errorf("Unexpected certificate ID: %s", id)
		}
	}
}

// TestCertificateConfig_LoadConfig tests loading configuration.
func TestCertificateConfig_LoadConfig(t *testing.T) {
	// Test default config
	config, err := LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	if config.ValidationMode != ValidationModeStrict {
		t.Errorf("Expected strict validation mode, got %s", config.ValidationMode)
	}

	if !config.AutoDiscovery {
		t.Error("Expected auto-discovery to be enabled by default")
	}
}

// TestCertificateConfig_Validation tests configuration validation.
func TestCertificateConfig_Validation(t *testing.T) {
	config := DefaultConfig()

	// Valid config should pass
	if err := config.Validate(); err != nil {
		t.Errorf("Default config validation failed: %v", err)
	}

	// Invalid validation mode should fail
	config.ValidationMode = "invalid"
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for invalid validation mode")
	}

	// Reset and test empty storage path
	config = DefaultConfig()
	config.StoragePath = ""
	if err := config.Validate(); err == nil {
		t.Error("Expected validation to fail for empty storage path")
	}
}

// TestCertificateConfig_EnvironmentVariables tests environment variable loading.
func TestCertificateConfig_EnvironmentVariables(t *testing.T) {
	// Set test environment variables
	os.Setenv(EnvCertPath, "/test/path")
	os.Setenv(EnvCertAutoDiscover, "false")
	os.Setenv(EnvCertValidation, ValidationModePermissive)
	defer func() {
		os.Unsetenv(EnvCertPath)
		os.Unsetenv(EnvCertAutoDiscover)
		os.Unsetenv(EnvCertValidation)
	}()

	config, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("Failed to load config from environment: %v", err)
	}

	if !filepath.IsAbs(config.StoragePath) || filepath.Clean(config.StoragePath) != config.StoragePath {
		t.Errorf("Expected absolute path, got %s", config.StoragePath)
	}

	if config.AutoDiscovery {
		t.Error("Expected auto-discovery to be disabled via environment variable")
	}

	if config.ValidationMode != ValidationModePermissive {
		t.Errorf("Expected permissive validation mode, got %s", config.ValidationMode)
	}
}

// createTestCertificate creates a test certificate for testing purposes.
func createTestCertificate(t *testing.T) *Certificate {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Test Org"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Test City"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: nil,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Parse the certificate
	x509Cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	// Encode certificate as PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	// Encode private key as PEM
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyDER})

	return &Certificate{
		ID:         "test-cert",
		Data:       certPEM,
		PrivateKey: keyPEM,
		X509:       x509Cert,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
