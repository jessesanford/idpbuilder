package certs

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/fs"
	"math/big"
	"testing"
	"time"
)

// createTestCertificate creates a test certificate for testing purposes
func createTestCertificate(t *testing.T, subject string) []byte {
	t.Helper()
	
	// Generate a private key
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}
	
	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: subject,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	
	// Create the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}
	
	// Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	return certPEM
}

func TestNewTrustManager(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	if manager == nil {
		t.Fatal("Expected non-nil trust manager")
	}
}

func TestDefaultTrustStoreConfig(t *testing.T) {
	config := DefaultTrustStoreConfig()
	
	if config.Location != UserTrustStore {
		t.Errorf("Expected UserTrustStore, got %v", config.Location)
	}
	
	if config.DirPermissions != fs.FileMode(0755) {
		t.Errorf("Expected 0755 dir permissions, got %v", config.DirPermissions)
	}
	
	if config.FilePermissions != fs.FileMode(0600) {
		t.Errorf("Expected 0600 file permissions, got %v", config.FilePermissions)
	}
}

func TestAddCertificate(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	ctx := context.Background()
	
	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}
	
	err := manager.AddCertificate(ctx, "test-registry", cert)
	if err != nil {
		t.Fatalf("Failed to add certificate: %v", err)
	}
	
	// Verify certificate was stored
	if cert.Info.Subject == "" {
		t.Error("Expected certificate info to be populated")
	}
	
	if cert.Info.Fingerprint == "" {
		t.Error("Expected certificate fingerprint to be populated")
	}
}

func TestAddCertificateValidation(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	ctx := context.Background()
	
	tests := []struct {
		name        string
		registry    string
		cert        *Certificate
		expectError bool
	}{
		{
			name:        "empty registry",
			registry:    "",
			cert:        &Certificate{Data: createTestCertificate(t, "test")},
			expectError: true,
		},
		{
			name:        "nil certificate",
			registry:    "test-registry",
			cert:        nil,
			expectError: true,
		},
		{
			name:        "valid inputs",
			registry:    "test-registry",
			cert:        &Certificate{Data: createTestCertificate(t, "test")},
			expectError: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.AddCertificate(ctx, tt.registry, tt.cert)
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestListCertificates(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	ctx := context.Background()
	
	// Add test certificates
	cert1 := &Certificate{Data: createTestCertificate(t, "test1.example.com")}
	cert2 := &Certificate{Data: createTestCertificate(t, "test2.example.com")}
	
	err := manager.AddCertificate(ctx, "test-registry", cert1)
	if err != nil {
		t.Fatalf("Failed to add certificate 1: %v", err)
	}
	
	err = manager.AddCertificate(ctx, "test-registry", cert2)
	if err != nil {
		t.Fatalf("Failed to add certificate 2: %v", err)
	}
	
	// List certificates
	certs, err := manager.ListCertificates(ctx, "test-registry")
	if err != nil {
		t.Fatalf("Failed to list certificates: %v", err)
	}
	
	if len(certs) != 2 {
		t.Errorf("Expected 2 certificates, got %d", len(certs))
	}
}

func TestRemoveCertificate(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	ctx := context.Background()
	
	// Add a certificate
	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}
	
	err := manager.AddCertificate(ctx, "test-registry", cert)
	if err != nil {
		t.Fatalf("Failed to add certificate: %v", err)
	}
	
	// Remove the certificate
	err = manager.RemoveCertificate(ctx, "test-registry", cert.Info.Fingerprint)
	if err != nil {
		t.Fatalf("Failed to remove certificate: %v", err)
	}
	
	// Verify certificate is gone
	certs, err := manager.ListCertificates(ctx, "test-registry")
	if err != nil {
		t.Fatalf("Failed to list certificates: %v", err)
	}
	
	if len(certs) != 0 {
		t.Errorf("Expected 0 certificates after removal, got %d", len(certs))
	}
}

func TestSetInsecureRegistry(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	ctx := context.Background()
	
	// Set registry as insecure
	err := manager.SetInsecureRegistry(ctx, "insecure-registry", true)
	if err != nil {
		t.Fatalf("Failed to set insecure registry: %v", err)
	}
	
	// Get registry config
	regConfig, err := manager.GetRegistryConfig(ctx, "insecure-registry")
	if err != nil {
		t.Fatalf("Failed to get registry config: %v", err)
	}
	
	if !regConfig.Insecure {
		t.Error("Expected registry to be marked as insecure")
	}
	
	// Unset insecure
	err = manager.SetInsecureRegistry(ctx, "insecure-registry", false)
	if err != nil {
		t.Fatalf("Failed to unset insecure registry: %v", err)
	}
	
	// Verify it's no longer insecure
	regConfig, err = manager.GetRegistryConfig(ctx, "insecure-registry")
	if err != nil {
		t.Fatalf("Failed to get registry config: %v", err)
	}
	
	if regConfig.Insecure {
		t.Error("Expected registry to not be marked as insecure")
	}
}

func TestTrustStoreLocationString(t *testing.T) {
	tests := []struct {
		location TrustStoreLocation
		expected string
	}{
		{UserTrustStore, "user"},
		{SystemTrustStore, "system"},
		{TrustStoreLocation(999), "unknown"},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.location.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestValidateCertificate(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	ctx := context.Background()
	
	// Create and add a test certificate
	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}
	
	err := manager.AddCertificate(ctx, "test-registry", cert)
	if err != nil {
		t.Fatalf("Failed to add certificate: %v", err)
	}
	
	// Parse the certificate for validation
	block, _ := pem.Decode(testCert)
	if block == nil {
		t.Fatal("Failed to decode test certificate")
	}
	
	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse X.509 certificate: %v", err)
	}
	
	// Validate the certificate - should pass
	err = manager.ValidateCertificate(ctx, "test-registry", x509Cert)
	if err != nil {
		t.Fatalf("Certificate validation failed: %v", err)
	}
	
	// Create a different certificate and try to validate - should fail
	otherCert := createTestCertificate(t, "other.example.com")
	otherBlock, _ := pem.Decode(otherCert)
	otherX509Cert, err := x509.ParseCertificate(otherBlock.Bytes)
	if err != nil {
		t.Fatalf("Failed to parse other certificate: %v", err)
	}
	
	err = manager.ValidateCertificate(ctx, "test-registry", otherX509Cert)
	if err == nil {
		t.Error("Expected validation to fail for untrusted certificate")
	}
}

func TestGetRegistryConfig(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}
	
	manager := NewTrustManager(config)
	ctx := context.Background()
	
	// Add a certificate and set as insecure
	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}
	
	err := manager.AddCertificate(ctx, "test-registry", cert)
	if err != nil {
		t.Fatalf("Failed to add certificate: %v", err)
	}
	
	err = manager.SetInsecureRegistry(ctx, "test-registry", true)
	if err != nil {
		t.Fatalf("Failed to set insecure registry: %v", err)
	}
	
	// Get registry config
	regConfig, err := manager.GetRegistryConfig(ctx, "test-registry")
	if err != nil {
		t.Fatalf("Failed to get registry config: %v", err)
	}
	
	if regConfig.Registry != "test-registry" {
		t.Errorf("Expected registry name 'test-registry', got '%s'", regConfig.Registry)
	}
	
	if !regConfig.Insecure {
		t.Error("Expected registry to be marked as insecure")
	}
	
	if len(regConfig.Certificates) != 1 {
		t.Errorf("Expected 1 certificate, got %d", len(regConfig.Certificates))
	}
}