package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestValidateCredentials(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	tests := []struct {
		name        string
		credentials *Credentials
		expectError bool
		errorType   error
	}{
		{
			name:        "nil credentials",
			credentials: nil,
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "valid basic auth",
			credentials: &Credentials{
				Registry:   "registry.example.com:5000",
				Username:   "testuser",
				Password:   "testpass",
				AuthMethod: AuthMethodBasic,
				CreatedAt:  now,
			},
			expectError: false,
		},
		{
			name: "basic auth missing username",
			credentials: &Credentials{
				Registry:   "registry.example.com:5000",
				Password:   "testpass",
				AuthMethod: AuthMethodBasic,
				CreatedAt:  now,
			},
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "basic auth missing password",
			credentials: &Credentials{
				Registry:   "registry.example.com:5000",
				Username:   "testuser",
				AuthMethod: AuthMethodBasic,
				CreatedAt:  now,
			},
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "valid token auth",
			credentials: &Credentials{
				Registry:   "registry.example.com",
				AuthMethod: AuthMethodToken,
				Token: &Token{
					Value:     "test-token-value",
					Type:      TokenTypeBearer,
					IssuedAt:  now,
					ExpiresAt: future,
				},
				CreatedAt: now,
			},
			expectError: false,
		},
		{
			name: "token auth missing token",
			credentials: &Credentials{
				Registry:   "registry.example.com",
				AuthMethod: AuthMethodToken,
				CreatedAt:  now,
			},
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "expired credentials",
			credentials: &Credentials{
				Registry:   "registry.example.com",
				Username:   "testuser",
				Password:   "testpass",
				AuthMethod: AuthMethodBasic,
				CreatedAt:  now,
				ExpiresAt:  &past,
			},
			expectError: true,
			errorType:   ErrExpiredCredentials,
		},
		{
			name: "invalid registry format",
			credentials: &Credentials{
				Registry:   "invalid registry format!",
				Username:   "testuser",
				Password:   "testpass",
				AuthMethod: AuthMethodBasic,
				CreatedAt:  now,
			},
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "unsupported auth method",
			credentials: &Credentials{
				Registry:   "registry.example.com",
				Username:   "testuser",
				Password:   "testpass",
				AuthMethod: AuthMethod("unsupported"),
				CreatedAt:  now,
			},
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentials(tt.credentials)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				// Check if error is of expected type (simplified check)
				if tt.errorType != nil && !isErrorType(err, tt.errorType) {
					t.Errorf("expected error type %v but got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name        string
		token       *Token
		expectError bool
		errorType   error
	}{
		{
			name:        "nil token",
			token:       nil,
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "valid bearer token",
			token: &Token{
				Value:     "valid-token",
				Type:      TokenTypeBearer,
				IssuedAt:  now,
				ExpiresAt: future,
			},
			expectError: false,
		},
		{
			name: "empty token value",
			token: &Token{
				Value:     "",
				Type:      TokenTypeBearer,
				IssuedAt:  now,
				ExpiresAt: future,
			},
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "empty token type",
			token: &Token{
				Value:     "valid-token",
				Type:      "",
				IssuedAt:  now,
				ExpiresAt: future,
			},
			expectError: true,
			errorType:   ErrInvalidCredentials,
		},
		{
			name: "expired token",
			token: &Token{
				Value:     "expired-token",
				Type:      TokenTypeBearer,
				IssuedAt:  past,
				ExpiresAt: past,
			},
			expectError: true,
			errorType:   ErrExpiredCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateToken(tt.token)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !isErrorType(err, tt.errorType) {
					t.Errorf("expected error type %v but got %v", tt.errorType, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateRegistryURL(t *testing.T) {
	tests := []struct {
		name        string
		registry    string
		expectError bool
	}{
		{
			name:        "empty registry",
			registry:    "",
			expectError: true,
		},
		{
			name:        "valid hostname",
			registry:    "registry.example.com",
			expectError: false,
		},
		{
			name:        "valid hostname with port",
			registry:    "registry.example.com:5000",
			expectError: false,
		},
		{
			name:        "valid HTTP URL",
			registry:    "http://registry.example.com",
			expectError: false,
		},
		{
			name:        "valid HTTPS URL",
			registry:    "https://registry.example.com:5000",
			expectError: false,
		},
		{
			name:        "invalid hostname format",
			registry:    "invalid registry format!",
			expectError: true,
		},
		{
			name:        "URL without scheme",
			registry:    "://registry.example.com",
			expectError: true,
		},
		{
			name:        "URL without host",
			registry:    "https://",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRegistryURL(tt.registry)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateCertificate(t *testing.T) {
	// Generate a valid test certificate
	validCert := generateTestCertificate(t, time.Now().Add(24*time.Hour))
	expiredCert := generateTestCertificate(t, time.Now().Add(-24*time.Hour))
	futureCert := generateTestCertificate(t, time.Now().Add(-48*time.Hour), time.Now().Add(24*time.Hour))

	tests := []struct {
		name        string
		certData    []byte
		expectError bool
		errorType   error
	}{
		{
			name:        "empty certificate data",
			certData:    []byte{},
			expectError: true,
			errorType:   ErrInvalidCertificate,
		},
		{
			name:        "invalid PEM data",
			certData:    []byte("invalid pem data"),
			expectError: true,
			errorType:   ErrInvalidCertificate,
		},
		{
			name:        "valid certificate",
			certData:    validCert,
			expectError: false,
		},
		{
			name:        "expired certificate",
			certData:    expiredCert,
			expectError: true,
			errorType:   ErrInvalidCertificate,
		},
		{
			name:        "future certificate",
			certData:    futureCert,
			expectError: true,
			errorType:   ErrInvalidCertificate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ValidateCertificate(tt.certData)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errorType != nil && !isErrorType(err, tt.errorType) {
					t.Errorf("expected error type %v but got %v", tt.errorType, err)
				}
				if info != nil {
					t.Errorf("expected nil info on error but got %v", info)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if info == nil {
					t.Errorf("expected certificate info but got nil")
					return
				}
				// Basic validation of certificate info
				if info.Subject == "" {
					t.Errorf("expected non-empty subject")
				}
				if info.Issuer == "" {
					t.Errorf("expected non-empty issuer")
				}
			}
		})
	}
}

func TestValidateHostnamePort(t *testing.T) {
	// Test custom validator directly
	tests := []struct {
		hostname string
		valid    bool
	}{
		{"registry.example.com", true},
		{"registry.example.com:5000", true},
		{"localhost", true},
		{"localhost:8080", true},
		{"", false},
		{"invalid hostname!", false},
		{"registry:invalid-port", false},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			// Test hostname validation using our validateHostnamePort function
			valid := validateHostnamePort(tt.hostname)
			if tt.valid && !valid {
				t.Errorf("expected valid hostname %s but validation failed", tt.hostname)
			}
			if !tt.valid && valid {
				t.Errorf("expected invalid hostname %s but validation passed", tt.hostname)
			}
		})
	}
}

// Helper function to check error types
func isErrorType(err, expectedType error) bool {
	return err == expectedType || (err != nil && expectedType != nil && err.Error() != expectedType.Error())
}

// generateTestCertificate creates a test certificate with specified validity period
func generateTestCertificate(t *testing.T, notAfter time.Time, notBefore ...time.Time) []byte {
	t.Helper()

	// Generate a private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	// Set up certificate template
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
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"test.example.com"},
	}

	// Set NotBefore
	if len(notBefore) > 0 {
		template.NotBefore = notBefore[0]
	} else {
		template.NotBefore = time.Now().Add(-1 * time.Hour)
	}

	// Create the certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	// Encode to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	return certPEM
}