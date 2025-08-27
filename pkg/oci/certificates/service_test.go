// Package certificates provides comprehensive tests for certificate service implementation.
// This file tests all major functionality including service creation, verification modes,
// thread safety, and integration with Gitea discovery.
package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServiceCreation tests proper service initialization
func TestServiceCreation(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful service creation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewCertificateService()
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				
				// Verify default settings
				assert.Equal(t, v2.VerificationModeStrict, service.verificationMode)
				assert.NotNil(t, service.certPool)
				assert.NotNil(t, service.systemPool)
				assert.NotNil(t, service.certificates)
				assert.NotNil(t, service.giteaDiscovery)
				assert.NotNil(t, service.verificationMgr)
			}
		})
	}
}

// TestVerificationModes tests all verification modes and switching
func TestVerificationModes(t *testing.T) {
	service, err := NewCertificateService()
	require.NoError(t, err)

	ctx := context.Background()

	tests := []struct {
		name     string
		mode     v2.VerificationMode
		wantErr  bool
	}{
		{
			name:    "switch to strict mode",
			mode:    v2.VerificationModeStrict,
			wantErr: false,
		},
		{
			name:    "switch to custom CA mode",
			mode:    v2.VerificationModeCustomCA,
			wantErr: false,
		},
		{
			name:    "switch to skip mode",
			mode:    v2.VerificationModeSkip,
			wantErr: false,
		},
		{
			name:    "invalid mode",
			mode:    "invalid-mode",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.SetVerificationMode(ctx, tt.mode)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mode, service.verificationMode)
			}
		})
	}
}

// TestThreadSafety tests concurrent operations for race conditions
func TestThreadSafety(t *testing.T) {
	service, err := NewCertificateService()
	require.NoError(t, err)

	ctx := context.Background()
	numGoroutines := 10
	numOperations := 100

	// Generate test certificate
	cert, err := generateTestCertificate()
	require.NoError(t, err)

	// Test concurrent certificate operations
	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3)

	// Concurrent certificate additions
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				testCert, _ := generateTestCertificate()
				service.AddCertificate(ctx, testCert)
			}
		}(i)
	}

	// Concurrent mode switching
	modes := []v2.VerificationMode{
		v2.VerificationModeStrict,
		v2.VerificationModeCustomCA,
		v2.VerificationModeSkip,
	}
	
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				mode := modes[j%len(modes)]
				service.SetVerificationMode(ctx, mode)
			}
		}(i)
	}

	// Concurrent certificate validations
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				service.ValidateCertificate(ctx, cert)
			}
		}(i)
	}

	// Wait for all operations to complete
	wg.Wait()

	// Verify service is still functional
	pool := service.GetCertPool()
	assert.NotNil(t, pool)

	result, err := service.ValidateCertificate(ctx, cert)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestCertificateValidation tests certificate validation logic
func TestCertificateValidation(t *testing.T) {
	service, err := NewCertificateService()
	require.NoError(t, err)

	ctx := context.Background()

	// Generate valid certificate
	validCert, err := generateTestCertificate()
	require.NoError(t, err)

	// Generate expired certificate
	expiredCert, err := generateExpiredCertificate()
	require.NoError(t, err)

	tests := []struct {
		name      string
		cert      *x509.Certificate
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "valid certificate",
			cert:      validCert,
			wantValid: false, // Will fail chain validation in test environment
			wantErr:   false,
		},
		{
			name:      "expired certificate",
			cert:      expiredCert,
			wantValid: false,
			wantErr:   false,
		},
		{
			name:      "nil certificate",
			cert:      nil,
			wantValid: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ValidateCertificate(ctx, tt.cert)
			
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantValid, result.Valid)
				
				if tt.cert != nil {
					assert.NotNil(t, result.Certificate)
					assert.Equal(t, tt.cert.Subject.String(), result.Certificate.Subject)
				}
			}
		})
	}
}

// TestPoolOperations tests certificate pool management
func TestPoolOperations(t *testing.T) {
	service, err := NewCertificateService()
	require.NoError(t, err)

	ctx := context.Background()

	// Generate test certificate
	cert, err := generateTestCertificate()
	require.NoError(t, err)

	// Test adding certificate
	err = service.AddCertificate(ctx, cert)
	assert.NoError(t, err)

	// Test getting pool
	pool := service.GetCertPool()
	assert.NotNil(t, pool)

	// Test removing certificate
	fingerprint := service.getCertificateFingerprint(cert)
	err = service.RemoveCertificate(ctx, fingerprint)
	assert.NoError(t, err)

	// Test removing non-existent certificate
	err = service.RemoveCertificate(ctx, "non-existent-fingerprint")
	assert.Error(t, err)
}

// generateTestCertificate creates a test certificate for testing
func generateTestCertificate() (*x509.Certificate, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Certificate template
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
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost", "test.example.com"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	// Parse certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

// generateExpiredCertificate creates an expired test certificate
func generateExpiredCertificate() (*x509.Certificate, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Certificate template (expired)
	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Test Org"},
			Country:      []string{"US"},
		},
		NotBefore:             time.Now().Add(-365 * 24 * time.Hour), // 1 year ago
		NotAfter:              time.Now().Add(-24 * time.Hour),        // 1 day ago
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	// Parse certificate
	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}

	return cert, nil
}