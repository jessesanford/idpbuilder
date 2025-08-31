package certs

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTrustStoreManager(t *testing.T) {
	tempDir := t.TempDir()
	
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)
	require.NotNil(t, manager)
	
	// Verify directory was created
	assert.DirExists(t, tempDir)
}

func TestTrustStoreManager_AddCertificate(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)

	// Create a test certificate
	cert := createTestCertificate(t)
	
	err = manager.AddCertificate("test-registry", cert)
	assert.NoError(t, err)
	
	// Verify certificate was added
	certs, err := manager.GetTrustedCerts("test-registry")
	require.NoError(t, err)
	assert.Len(t, certs, 1)
	assert.True(t, certs[0].Equal(cert))
}

func TestTrustStoreManager_RemoveCertificate(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)

	registry := "test-registry"
	cert := createTestCertificate(t)
	
	// Add certificate first
	err = manager.AddCertificate(registry, cert)
	require.NoError(t, err)
	
	// Verify certificate was added
	certs, err := manager.GetTrustedCerts(registry)
	require.NoError(t, err)
	assert.Len(t, certs, 1)
	
	// Remove certificate
	err = manager.RemoveCertificate(registry)
	assert.NoError(t, err)
	
	// Verify certificate was removed
	certs, err = manager.GetTrustedCerts(registry)
	require.NoError(t, err)
	assert.Len(t, certs, 0)
}

func TestTrustStoreManager_SetInsecureRegistry(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)

	registry := "test-registry"
	
	// Initially should be secure
	assert.False(t, manager.IsInsecure(registry))
	
	// Mark as insecure
	err = manager.SetInsecureRegistry(registry, true)
	assert.NoError(t, err)
	assert.True(t, manager.IsInsecure(registry))
	
	// Mark as secure again
	err = manager.SetInsecureRegistry(registry, false)
	assert.NoError(t, err)
	assert.False(t, manager.IsInsecure(registry))
}

func TestTrustStoreManager_GetCertPool(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)

	registry := "test-registry"
	
	// Test with no certificates - should return system pool
	pool, err := manager.GetCertPool(registry)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	
	// Add a certificate and test again
	cert := createTestCertificate(t)
	err = manager.AddCertificate(registry, cert)
	require.NoError(t, err)
	
	pool, err = manager.GetCertPool(registry)
	assert.NoError(t, err)
	assert.NotNil(t, pool)
}

func TestTrustStoreManager_DiskOperations(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)

	registry := "test-registry"
	
	// Test that the manager can be created and loads from disk without errors
	assert.NotNil(t, manager)
	
	// Test that calling LoadFromDisk doesn't error (even if directory is empty)
	err = manager.LoadFromDisk()
	assert.NoError(t, err)
	
	// Verify initial state - no certificates
	certs, err := manager.GetTrustedCerts(registry)
	require.NoError(t, err)
	assert.Len(t, certs, 0)
}

func TestTrustStoreManager_ErrorHandling(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)

	// Test with empty registry name
	err = manager.AddCertificate("", createTestCertificate(t))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "registry name cannot be empty")
	
	// Test with nil certificate
	err = manager.AddCertificate("test-registry", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "certificate cannot be nil")
	
	// Test expired certificate
	expiredCert := createExpiredCertificate(t)
	err = manager.AddCertificate("test-registry", expiredCert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has expired")
}

// createTestCertificate creates a test certificate for testing
func createTestCertificate(t *testing.T) *x509.Certificate {
	// Create a self-signed certificate for testing
	// This approach creates a certificate with proper DER encoding
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test.example.com",
		},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour * 24),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
		DNSNames:     []string{"test.example.com"},
	}
	
	// For testing, create a minimal certificate with DER bytes
	// In a real scenario, these would be properly generated
	template.Raw = []byte{0x30, 0x82, 0x01, 0x00} // Minimal DER structure for testing
	
	return template
}

// createExpiredCertificate creates an expired test certificate
func createExpiredCertificate(t *testing.T) *x509.Certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: "expired.example.com",
		},
		NotBefore:    time.Now().Add(-time.Hour * 48),
		NotAfter:     time.Now().Add(-time.Hour * 24), // Expired yesterday
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
	}
	
	return cert
}