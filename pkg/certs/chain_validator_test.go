package certs

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

// mockValidator implements CertValidator for testing
type mockValidator struct {
	validationResult *ValidationResult
	expiryResult     *ExpiryResult
	shouldError      bool
}

func (m *mockValidator) ValidateCertificate(cert *x509.Certificate) (*ValidationResult, error) {
	if m.shouldError {
		return nil, &CertificateError{Code: "TEST_ERROR", Message: "test error"}
	}
	return m.validationResult, nil
}

func (m *mockValidator) CheckExpiry(cert *x509.Certificate, warnDays int) (*ExpiryResult, error) {
	if m.shouldError {
		return nil, &CertificateError{Code: "TEST_ERROR", Message: "test error"}
	}
	return m.expiryResult, nil
}

// createTestCertificate creates a test certificate
func createTestCertificate(t *testing.T, subject, issuer pkix.Name, isSelfSigned bool, dnsNames []string) *x509.Certificate {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      subject,
		Issuer:       issuer,
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     dnsNames,
	}

	var parent *x509.Certificate
	var parentKey *rsa.PrivateKey
	
	if isSelfSigned {
		parent = &template
		parentKey = privateKey
	} else {
		// For simplicity, use self-signed for parent too
		parent = &template
		parentKey = privateKey
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, parent, &privateKey.PublicKey, parentKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	return cert
}

func TestNewChainValidator(t *testing.T) {
	tests := []struct {
		name        string
		config      *ChainValidatorConfig
		shouldPanic bool
	}{
		{
			name:        "nil config should panic",
			config:      nil,
			shouldPanic: true,
		},
		{
			name: "nil BasicValidator should panic",
			config: &ChainValidatorConfig{
				BasicValidator: nil,
				TrustManager:   nil,
			},
			shouldPanic: true,
		},
		{
			name: "valid config should work",
			config: &ChainValidatorConfig{
				BasicValidator:  &mockValidator{},
				TrustManager:    nil,
				AllowSelfSigned: true,
			},
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but none occurred")
					}
				}()
			}
			validator := NewChainValidator(tt.config)
			if !tt.shouldPanic && validator == nil {
				t.Error("Expected validator but got nil")
			}
		})
	}
}

func TestValidateChain(t *testing.T) {
	subject := pkix.Name{CommonName: "test.example.com"}
	issuer := pkix.Name{CommonName: "Test CA"}
	
	testCert := createTestCertificate(t, subject, issuer, true, []string{"test.example.com"})

	tests := []struct {
		name          string
		cert          *x509.Certificate
		intermediates []*x509.Certificate
		validator     *mockValidator
		expectValid   bool
		expectError   bool
	}{
		{
			name:        "nil certificate should error",
			cert:        nil,
			expectError: true,
		},
		{
			name: "single valid self-signed cert should pass",
			cert: testCert,
			validator: &mockValidator{
				validationResult: &ValidationResult{Valid: true},
			},
			expectValid: true,
		},
		{
			name: "cert with validation issues should fail",
			cert: testCert,
			validator: &mockValidator{
				validationResult: &ValidationResult{
					Valid:  false,
					Issues: []string{"test issue"},
				},
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &ChainValidatorConfig{
				BasicValidator:  tt.validator,
				AllowSelfSigned: true,
			}
			if tt.validator == nil {
				config.BasicValidator = &mockValidator{
					validationResult: &ValidationResult{Valid: true},
				}
			}
			
			validator := NewChainValidator(config)
			result, err := validator.ValidateChain(context.Background(), tt.cert, tt.intermediates)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got %v", tt.expectValid, result.Valid)
			}
		})
	}
}

func TestVerifyHostname(t *testing.T) {
	subject := pkix.Name{CommonName: "example.com"}
	issuer := pkix.Name{CommonName: "Test CA"}
	
	// Create certificates with different DNS names
	exactCert := createTestCertificate(t, subject, issuer, true, []string{"example.com"})
	wildcardCert := createTestCertificate(t, 
		pkix.Name{CommonName: "*.example.com"}, issuer, true, []string{"*.example.com"})
	sanCert := createTestCertificate(t, subject, issuer, true, 
		[]string{"example.com", "www.example.com", "api.example.com"})

	validator := NewChainValidator(&ChainValidatorConfig{
		BasicValidator:  &mockValidator{validationResult: &ValidationResult{Valid: true}},
		AllowSelfSigned: true,
	})

	tests := []struct {
		name        string
		cert        *x509.Certificate
		hostname    string
		expectError bool
	}{
		{
			name:        "nil certificate should error",
			cert:        nil,
			hostname:    "example.com",
			expectError: true,
		},
		{
			name:        "empty hostname should error",
			cert:        exactCert,
			hostname:    "",
			expectError: true,
		},
		{
			name:        "exact match should pass",
			cert:        exactCert,
			hostname:    "example.com",
			expectError: false,
		},
		{
			name:        "wildcard match should pass",
			cert:        wildcardCert,
			hostname:    "www.example.com",
			expectError: false,
		},
		{
			name:        "SAN match should pass",
			cert:        sanCert,
			hostname:    "api.example.com",
			expectError: false,
		},
		{
			name:        "no match should fail",
			cert:        exactCert,
			hostname:    "other.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.VerifyHostname(tt.cert, tt.hostname)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestCheckChainExpiry(t *testing.T) {
	// Create certificates with different expiry dates
	expiredCert := createExpiredCertificate(t)
	validCert := createValidCertificate(t)
	expiringSoonCert := createExpiringSoonCertificate(t, 15) // expires in 15 days

	validator := NewChainValidator(&ChainValidatorConfig{
		BasicValidator:  &mockValidator{validationResult: &ValidationResult{Valid: true}},
		AllowSelfSigned: true,
	})

	tests := []struct {
		name         string
		chain        []*x509.Certificate
		warnDays     int
		expectValid  bool
		expectError  bool
		expiredCount int
		expiringCount int
	}{
		{
			name:        "empty chain should error",
			chain:       []*x509.Certificate{},
			expectError: true,
		},
		{
			name:        "valid chain should pass",
			chain:       []*x509.Certificate{validCert},
			warnDays:    30,
			expectValid: true,
		},
		{
			name:         "expired certificate should fail",
			chain:        []*x509.Certificate{expiredCert},
			warnDays:     30,
			expectValid:  false,
			expiredCount: 1,
		},
		{
			name:          "expiring soon certificate should warn",
			chain:         []*x509.Certificate{expiringSoonCert},
			warnDays:      30,
			expectValid:   true,
			expiringCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.CheckChainExpiry(tt.chain, tt.warnDays)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.ChainValid != tt.expectValid {
				t.Errorf("Expected chainValid=%v, got %v", tt.expectValid, result.ChainValid)
			}

			if len(result.ExpiredCerts) != tt.expiredCount {
				t.Errorf("Expected %d expired certs, got %d", tt.expiredCount, len(result.ExpiredCerts))
			}

			if len(result.ExpiringCerts) != tt.expiringCount {
				t.Errorf("Expected %d expiring certs, got %d", tt.expiringCount, len(result.ExpiringCerts))
			}
		})
	}
}

func TestGenerateDiagnostics(t *testing.T) {
	subject := pkix.Name{CommonName: "test.example.com"}
	issuer := pkix.Name{CommonName: "Test CA"}
	testCert := createTestCertificate(t, subject, issuer, true, []string{"test.example.com"})

	validator := NewChainValidator(&ChainValidatorConfig{
		BasicValidator:  &mockValidator{validationResult: &ValidationResult{Valid: true}},
		AllowSelfSigned: true,
	})

	tests := []struct {
		name        string
		cert        *x509.Certificate
		hostname    string
		expectError bool
	}{
		{
			name:        "nil certificate should error",
			cert:        nil,
			hostname:    "test.example.com",
			expectError: true,
		},
		{
			name:        "valid certificate should generate report",
			cert:        testCert,
			hostname:    "test.example.com",
			expectError: false,
		},
		{
			name:        "valid certificate without hostname should generate report",
			cert:        testCert,
			hostname:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := validator.GenerateDiagnostics(context.Background(), tt.cert, tt.hostname)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if report == nil {
				t.Error("Expected report but got nil")
				return
			}

			if report.CertificateDetails == nil {
				t.Error("Expected certificate details but got nil")
			}

			if report.ChainAnalysis == nil {
				t.Error("Expected chain analysis but got nil")
			}

			if tt.hostname != "" && report.HostnameValidation == nil {
				t.Error("Expected hostname validation but got nil")
			}
		})
	}
}

// Helper functions for creating test certificates with specific expiry dates

func createExpiredCertificate(t *testing.T) *x509.Certificate {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "expired.example.com"},
		Issuer:       pkix.Name{CommonName: "expired.example.com"},
		NotBefore:    time.Now().Add(-48 * time.Hour),
		NotAfter:     time.Now().Add(-24 * time.Hour), // Expired 24 hours ago
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	return cert
}

func createValidCertificate(t *testing.T) *x509.Certificate {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "valid.example.com"},
		Issuer:       pkix.Name{CommonName: "valid.example.com"},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(90 * 24 * time.Hour), // Valid for 90 days
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	return cert
}

func createExpiringSoonCertificate(t *testing.T, daysUntilExpiry int) *x509.Certificate {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "expiring.example.com"},
		Issuer:       pkix.Name{CommonName: "expiring.example.com"},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().Add(time.Duration(daysUntilExpiry) * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		t.Fatalf("Failed to parse certificate: %v", err)
	}

	return cert
}