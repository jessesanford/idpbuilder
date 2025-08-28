package certs

import (
	"crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
	cert := createTestCertificate(time.Now(), time.Now().Add(24*time.Hour))
	validator := NewValidator(cert)
	
	assert.NotNil(t, validator)
	assert.Equal(t, cert, validator.cert)
	assert.NotNil(t, validator.logger)
}

func TestValidator_ValidateChain(t *testing.T) {
	validator := NewValidator(nil)
	
	tests := []struct {
		name        string
		cert        *x509.Certificate
		expectError bool
	}{
		{
			name:        "nil certificate",
			cert:        nil,
			expectError: true,
		},
		{
			name:        "self-signed certificate",
			cert:        createTestCertificate(time.Now(), time.Now().Add(24*time.Hour)),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateChain(tt.cert)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.IsType(t, &CertificateError{}, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_CheckExpiry(t *testing.T) {
	validator := NewValidator(nil)
	now := time.Now()

	tests := []struct {
		name        string
		cert        *x509.Certificate
		expectError bool
		expectDuration bool
	}{
		{
			name:        "nil certificate",
			cert:        nil,
			expectError: true,
			expectDuration: false,
		},
		{
			name:        "expired certificate",
			cert:        createTestCertificate(now.Add(-48*time.Hour), now.Add(-24*time.Hour)),
			expectError: true,
			expectDuration: false,
		},
		{
			name:        "not yet valid certificate",
			cert:        createTestCertificate(now.Add(1*time.Hour), now.Add(25*time.Hour)),
			expectError: true,
			expectDuration: false,
		},
		{
			name:        "valid certificate",
			cert:        createTestCertificate(now.Add(-1*time.Hour), now.Add(24*time.Hour)),
			expectError: false,
			expectDuration: true,
		},
		{
			name:        "certificate expiring soon",
			cert:        createTestCertificate(now.Add(-1*time.Hour), now.Add(15*24*time.Hour)),
			expectError: false,
			expectDuration: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := validator.CheckExpiry(tt.cert)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, duration)
				assert.IsType(t, &CertificateError{}, err)
			} else {
				assert.NoError(t, err)
				if tt.expectDuration {
					assert.NotNil(t, duration)
					assert.True(t, *duration > 0)
				}
			}
		})
	}
}

func TestValidator_VerifyHostname(t *testing.T) {
	validator := NewValidator(nil)
	cert := createTestCertificate(time.Now(), time.Now().Add(24*time.Hour))

	tests := []struct {
		name        string
		cert        *x509.Certificate
		hostname    string
		expectError bool
	}{
		{
			name:        "nil certificate",
			cert:        nil,
			hostname:    "localhost",
			expectError: true,
		},
		{
			name:        "empty hostname",
			cert:        cert,
			hostname:    "",
			expectError: true,
		},
		{
			name:        "valid hostname",
			cert:        cert,
			hostname:    "localhost",
			expectError: false,
		},
		{
			name:        "DNS name in certificate",
			cert:        cert,
			hostname:    "gitea.gitea.svc.cluster.local",
			expectError: false,
		},
		{
			name:        "hostname not in certificate",
			cert:        cert,
			hostname:    "invalid.example.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.VerifyHostname(tt.cert, tt.hostname)
			
			if tt.expectError {
				assert.Error(t, err)
				if tt.cert != nil && tt.hostname != "" {
					assert.IsType(t, &CertificateError{}, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_GenerateDiagnostics(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name        string
		cert        *x509.Certificate
		expectError bool
		expectValid bool
		checkExpiry bool
	}{
		{
			name:        "nil certificate in validator",
			cert:        nil,
			expectError: true,
		},
		{
			name:        "valid certificate",
			cert:        createTestCertificate(now.Add(-1*time.Hour), now.Add(90*24*time.Hour)),
			expectError: false,
			expectValid: true,
			checkExpiry: false,
		},
		{
			name:        "expired certificate",
			cert:        createTestCertificate(now.Add(-48*time.Hour), now.Add(-24*time.Hour)),
			expectError: false,
			expectValid: false,
			checkExpiry: false,
		},
		{
			name:        "certificate expiring soon",
			cert:        createTestCertificate(now.Add(-1*time.Hour), now.Add(15*24*time.Hour)),
			expectError: false,
			expectValid: true,
			checkExpiry: true,
		},
		{
			name:        "not yet valid certificate",
			cert:        createTestCertificate(now.Add(1*time.Hour), now.Add(25*time.Hour)),
			expectError: false,
			expectValid: false,
			checkExpiry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.cert)
			
			diagnostics, err := validator.GenerateDiagnostics()
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, diagnostics)
				assert.IsType(t, &CertificateError{}, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, diagnostics)
				
				assert.Equal(t, tt.expectValid, diagnostics.IsValid)
				assert.NotEmpty(t, diagnostics.Subject)
				assert.NotEmpty(t, diagnostics.Issuer)
				assert.NotEmpty(t, diagnostics.DNSNames)
				
				if tt.checkExpiry {
					assert.Contains(t, diagnostics.Issues[0], "expires in")
				}
				
				if !tt.expectValid {
					assert.NotEmpty(t, diagnostics.Issues)
					assert.NotEmpty(t, diagnostics.Recommendations)
				}
			}
		})
	}
}

func TestValidator_ValidateForRegistry(t *testing.T) {
	now := time.Now()
	validCert := createTestCertificate(now.Add(-1*time.Hour), now.Add(24*time.Hour))
	expiredCert := createTestCertificate(now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	tests := []struct {
		name         string
		cert         *x509.Certificate
		registryHost string
		expectError  bool
	}{
		{
			name:         "valid certificate without hostname",
			cert:         validCert,
			registryHost: "",
			expectError:  false,
		},
		{
			name:         "valid certificate with valid hostname",
			cert:         validCert,
			registryHost: "localhost",
			expectError:  false,
		},
		{
			name:         "valid certificate with invalid hostname",
			cert:         validCert,
			registryHost: "invalid.example.com",
			expectError:  false, // Should log warning but not fail for MVP
		},
		{
			name:         "expired certificate",
			cert:         expiredCert,
			registryHost: "localhost",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(tt.cert)
			
			err := validator.ValidateForRegistry(tt.cert, tt.registryHost)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCertDiagnostics_Fields(t *testing.T) {
	now := time.Now()
	cert := createTestCertificate(now.Add(-1*time.Hour), now.Add(30*24*time.Hour))
	validator := NewValidator(cert)
	
	diagnostics, err := validator.GenerateDiagnostics()
	require.NoError(t, err)
	require.NotNil(t, diagnostics)
	
	// Check all fields are populated correctly
	assert.True(t, diagnostics.IsValid)
	assert.Equal(t, cert.NotAfter, diagnostics.ExpiresAt)
	assert.True(t, diagnostics.DaysUntilExpiry > 25 && diagnostics.DaysUntilExpiry < 35)
	assert.Equal(t, cert.Subject.String(), diagnostics.Subject)
	assert.Equal(t, cert.Issuer.String(), diagnostics.Issuer)
	assert.Equal(t, cert.DNSNames, diagnostics.DNSNames)
	assert.NotNil(t, diagnostics.Issues)
	assert.NotNil(t, diagnostics.Recommendations)
}