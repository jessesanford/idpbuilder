package certs

import (
	"crypto/x509"
	"testing"
	"time"
)

// MockTrustStore implements TrustStoreProvider for testing
type MockTrustStore struct {
	trustedCerts map[string][]*x509.Certificate
	insecure     map[string]bool
}

func NewMockTrustStore() *MockTrustStore {
	return &MockTrustStore{
		trustedCerts: make(map[string][]*x509.Certificate),
		insecure:     make(map[string]bool),
	}
}

func (m *MockTrustStore) GetTrustedCerts(registry string) ([]*x509.Certificate, error) {
	return m.trustedCerts[registry], nil
}

func (m *MockTrustStore) IsInsecure(registry string) bool {
	return m.insecure[registry]
}

func (m *MockTrustStore) AddTrustedCert(registry string, cert *x509.Certificate) {
	m.trustedCerts[registry] = append(m.trustedCerts[registry], cert)
}

func (m *MockTrustStore) SetInsecure(registry string, insecure bool) {
	m.insecure[registry] = insecure
}

<<<<<<< HEAD
func TestNewChainValidator(t *testing.T) {
	trustStore := NewMockTrustStore()

=======

func TestNewChainValidator(t *testing.T) {
	trustStore := NewMockTrustStore()
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	validator := NewChainValidator(trustStore, StrictMode)
	if validator == nil {
		t.Fatal("NewChainValidator returned nil")
	}
<<<<<<< HEAD

	if validator.trustStore != trustStore {
		t.Error("Trust store not set correctly")
	}

=======
	
	if validator.trustStore != trustStore {
		t.Error("Trust store not set correctly")
	}
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	if validator.mode != StrictMode {
		t.Errorf("Expected StrictMode, got %v", validator.mode)
	}
}

func TestChainValidator_ValidateChain_EmptyChain(t *testing.T) {
	trustStore := NewMockTrustStore()
	validator := NewChainValidator(trustStore, StrictMode)
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	err := validator.ValidateChain([]*x509.Certificate{}, "test-registry", nil)
	if err == nil {
		t.Fatal("Expected error for empty chain")
	}
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	if validationErr, ok := err.(*ValidationError); ok {
		if validationErr.Type != ChainIncomplete {
			t.Errorf("Expected ChainIncomplete error, got %v", validationErr.Type)
		}
	} else {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestChainValidator_ValidateChain_SingleValidCert(t *testing.T) {
	trustStore := NewMockTrustStore()
	validator := NewChainValidator(trustStore, InsecureMode) // Use insecure mode to bypass trust checks
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	// Create a valid self-signed certificate
	cert := createTestCertificateWithTimes(t, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	cert.Subject.CommonName = "test.example.com"
	chain := []*x509.Certificate{cert}
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	options := &ChainValidationOptions{
		AllowSelfSigned: true,
		MaxChainLength:  1,
	}
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	err := validator.ValidateChain(chain, "test-registry", options)
	if err != nil {
		t.Fatalf("Expected no error for valid single certificate in insecure mode, got: %v", err)
	}
}

func TestChainValidator_ValidateChain_ChainTooLong(t *testing.T) {
	trustStore := NewMockTrustStore()
	validator := NewChainValidator(trustStore, StrictMode)
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	// Create multiple certificates to exceed limit
	certs := make([]*x509.Certificate, 5)
	for i := range certs {
		certs[i] = createTestCertificateWithTimes(t, time.Now().Add(-24*time.Hour), time.Now().Add(24*time.Hour))
	}
<<<<<<< HEAD

	options := &ChainValidationOptions{
		MaxChainLength: 3, // Set limit lower than chain length
	}

=======
	
	options := &ChainValidationOptions{
		MaxChainLength: 3, // Set limit lower than chain length
	}
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	err := validator.ValidateChain(certs, "test-registry", options)
	if err == nil {
		t.Fatal("Expected error for chain too long")
	}
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	if validationErr, ok := err.(*ValidationError); ok {
		if validationErr.Type != ChainTooLong {
			t.Errorf("Expected ChainTooLong error, got %v", validationErr.Type)
		}
	} else {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

func TestChainValidator_ValidationModes(t *testing.T) {
	tests := []struct {
		name string
		mode ValidationMode
		want string
	}{
		{"Strict Mode", StrictMode, "Strict"},
		{"Lenient Mode", LenientMode, "Lenient"},
		{"Insecure Mode", InsecureMode, "Insecure"},
		{"Unknown Mode", ValidationMode(99), "Unknown(99)"},
	}
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.String(); got != tt.want {
				t.Errorf("ValidationMode.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChainValidator_DefaultOptions(t *testing.T) {
	trustStore := NewMockTrustStore()
	validator := NewChainValidator(trustStore, StrictMode)
<<<<<<< HEAD

=======
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	// Test default max chain length for different modes
	strictValidator := NewChainValidator(trustStore, StrictMode)
	lenientValidator := NewChainValidator(trustStore, LenientMode)
	insecureValidator := NewChainValidator(trustStore, InsecureMode)
<<<<<<< HEAD

	if strictValidator.getMaxChainLengthForMode() != 4 {
		t.Errorf("Expected strict mode max chain length 4, got %d", strictValidator.getMaxChainLengthForMode())
	}

	if lenientValidator.getMaxChainLengthForMode() != 6 {
		t.Errorf("Expected lenient mode max chain length 6, got %d", lenientValidator.getMaxChainLengthForMode())
	}

	if insecureValidator.getMaxChainLengthForMode() != 10 {
		t.Errorf("Expected insecure mode max chain length 10, got %d", insecureValidator.getMaxChainLengthForMode())
	}

=======
	
	if strictValidator.getMaxChainLengthForMode() != 4 {
		t.Errorf("Expected strict mode max chain length 4, got %d", strictValidator.getMaxChainLengthForMode())
	}
	
	if lenientValidator.getMaxChainLengthForMode() != 6 {
		t.Errorf("Expected lenient mode max chain length 6, got %d", lenientValidator.getMaxChainLengthForMode())
	}
	
	if insecureValidator.getMaxChainLengthForMode() != 10 {
		t.Errorf("Expected insecure mode max chain length 10, got %d", insecureValidator.getMaxChainLengthForMode())
	}
	
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
	// Test default options
	options := validator.getDefaultOptions()
	if options == nil {
		t.Fatal("getDefaultOptions returned nil")
	}
<<<<<<< HEAD

	if options.MaxChainLength != 4 {
		t.Errorf("Expected default max chain length 4, got %d", options.MaxChainLength)
	}
}
=======
	
	if options.MaxChainLength != 4 {
		t.Errorf("Expected default max chain length 4, got %d", options.MaxChainLength)
	}
}
>>>>>>> origin/idpbuilder-oci-build-push/phase2/wave1/gitea-client-split-002
