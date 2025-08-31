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

func TestTrustStoreManager_ConfigureTransport(t *testing.T) {
	tempDir := t.TempDir()
	manager, err := NewTrustStoreManager(tempDir)
	require.NoError(t, err)

	registry := "test-registry"
	
	// Test with insecure registry
	err = manager.SetInsecureRegistry(registry, true)
	require.NoError(t, err)
	
	option, err := manager.ConfigureTransport(registry)
	assert.NoError(t, err)
	assert.NotNil(t, option)
}

func TestTrustStoreUtils_LoadCertificateFromPEM(t *testing.T) {
	utils := NewTrustStoreUtils()
	
	// Test with valid PEM data
	cert := createTestCertificate(t)
	utils2 := NewTrustStoreUtils()
	pemData, err := utils2.CertificateToPEM(cert)
	require.NoError(t, err)
	
	loadedCert, err := utils.LoadCertificateFromPEM(pemData)
	assert.NoError(t, err)
	assert.True(t, cert.Equal(loadedCert))
}

func TestTrustStoreUtils_ValidateCertificate(t *testing.T) {
	utils := NewTrustStoreUtils()
	
	// Test with valid certificate
	cert := createTestCertificate(t)
	err := utils.ValidateCertificate(cert)
	assert.NoError(t, err)
	
	// Test with nil certificate
	err = utils.ValidateCertificate(nil)
	assert.Error(t, err)
}

// createTestCertificate creates a test certificate for testing
func createTestCertificate(t *testing.T) *x509.Certificate {
	// This creates a minimal test certificate
	// In a real implementation, you'd want more comprehensive test certificates
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "test.example.com",
		},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour * 24),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
	}
	
	return cert
}