package certs

import (
	"crypto/x509"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs/testdata"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// mockTrustStoreManager implements TrustStoreManager for testing
type mockTrustStoreManager struct {
	trustedCerts map[string][]*x509.Certificate
	insecure     map[string]bool
}

func newMockTrustStoreManager() *mockTrustStoreManager {
	return &mockTrustStoreManager{
		trustedCerts: make(map[string][]*x509.Certificate),
		insecure:     make(map[string]bool),
	}
}

func (m *mockTrustStoreManager) AddCertificate(registry string, cert *x509.Certificate) error {
	if m.trustedCerts[registry] == nil {
		m.trustedCerts[registry] = make([]*x509.Certificate, 0)
	}
	m.trustedCerts[registry] = append(m.trustedCerts[registry], cert)
	return nil
}

func (m *mockTrustStoreManager) RemoveCertificate(registry string) error {
	delete(m.trustedCerts, registry)
	return nil
}

func (m *mockTrustStoreManager) SetInsecureRegistry(registry string, insecure bool) error {
	m.insecure[registry] = insecure
	return nil
}

func (m *mockTrustStoreManager) GetTrustedCerts(registry string) ([]*x509.Certificate, error) {
	return m.trustedCerts[registry], nil
}

func (m *mockTrustStoreManager) GetCertPool(registry string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	if certs, ok := m.trustedCerts[registry]; ok {
		for _, cert := range certs {
			pool.AddCert(cert)
		}
	}
	return pool, nil
}

func (m *mockTrustStoreManager) IsInsecure(registry string) bool {
	return m.insecure[registry]
}

func (m *mockTrustStoreManager) LoadFromDisk() error                        { return nil }
func (m *mockTrustStoreManager) SaveToDisk(string, *x509.Certificate) error { return nil }
func (m *mockTrustStoreManager) ConfigureTransport(registry string) (remote.Option, error) {
	return nil, nil
}
func (m *mockTrustStoreManager) ConfigureTransportWithConfig(registry string, config *TransportConfig) (remote.Option, error) {
	return nil, nil
}
func (m *mockTrustStoreManager) CreateHTTPClient(registry string) (*http.Client, error) {
	return nil, nil
}
func (m *mockTrustStoreManager) CreateHTTPClientWithConfig(registry string, config *TransportConfig) (*http.Client, error) {
	return nil, nil
}

func TestNewValidator(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()

	validator, err := NewValidator(mockTrustStore)
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	if validator == nil {
		t.Fatal("Expected validator, got nil")
	}

	if validator.trustStore != mockTrustStore {
		t.Error("Trust store not set correctly")
	}

	if validator.expiryWarning != 30*24*time.Hour {
		t.Errorf("Expected default expiry warning of 30 days, got %v", validator.expiryWarning)
	}
}

func TestNewValidator_NilTrustStore(t *testing.T) {
	validator, err := NewValidator(nil)
	if err == nil {
		t.Error("Expected error with nil trust store")
	}
	if validator != nil {
		t.Error("Expected nil validator with error")
	}
}

func TestNewValidatorWithWarningThreshold(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	customThreshold := 7 * 24 * time.Hour // 7 days

	validator, err := NewValidatorWithWarningThreshold(mockTrustStore, customThreshold)
	if err != nil {
		t.Fatalf("NewValidatorWithWarningThreshold failed: %v", err)
	}

	if validator.expiryWarning != customThreshold {
		t.Errorf("Expected expiry warning of %v, got %v", customThreshold, validator.expiryWarning)
	}
}

func TestValidateChain_ValidCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	// Add self-signed cert to trust store
	mockTrustStore.AddCertificate("gitea.local", validCert)

	err = validator.ValidateChain(validCert)
	if err != nil {
		t.Errorf("ValidateChain failed for valid cert: %v", err)
	}
}

func TestValidateChain_NilCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	err := validator.ValidateChain(nil)
	if err == nil {
		t.Error("Expected error for nil certificate")
	}
	if !strings.Contains(err.Error(), "certificate cannot be nil") {
		t.Errorf("Expected 'certificate cannot be nil' error, got: %v", err)
	}
}

func TestCheckExpiry_ValidCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	duration, err := validator.CheckExpiry(validCert)
	if err != nil {
		t.Errorf("CheckExpiry failed for valid cert: %v", err)
	}
	if duration == nil {
		t.Error("Expected duration, got nil")
	}
	if *duration <= 0 {
		t.Errorf("Expected positive duration, got %v", *duration)
	}
}

func TestCheckExpiry_ExpiredCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	expiredCert, err := testdata.GenerateExpiredCert()
	if err != nil {
		t.Fatalf("Failed to generate expired cert: %v", err)
	}

	duration, err := validator.CheckExpiry(expiredCert)
	if err == nil {
		t.Error("Expected error for expired certificate")
	}
	if duration != nil {
		t.Error("Expected nil duration for expired certificate")
	}
	if !strings.Contains(err.Error(), "expired") {
		t.Errorf("Expected 'expired' in error message, got: %v", err)
	}
}

func TestCheckExpiry_ExpiringSoonCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	expiringSoonCert, err := testdata.GenerateExpiringSoonCert()
	if err != nil {
		t.Fatalf("Failed to generate expiring soon cert: %v", err)
	}

	duration, err := validator.CheckExpiry(expiringSoonCert)
	if err == nil {
		t.Error("Expected warning for expiring soon certificate")
	}
	if duration == nil {
		t.Error("Expected duration, got nil")
	}
	if !strings.Contains(err.Error(), "expires soon") {
		t.Errorf("Expected 'expires soon' in error message, got: %v", err)
	}
}

func TestCheckExpiry_NotYetValidCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	notYetValidCert, err := testdata.GenerateNotYetValidCert()
	if err != nil {
		t.Fatalf("Failed to generate not yet valid cert: %v", err)
	}

	duration, err := validator.CheckExpiry(notYetValidCert)
	if err == nil {
		t.Error("Expected error for not yet valid certificate")
	}
	if duration != nil {
		t.Error("Expected nil duration for not yet valid certificate")
	}
	if !strings.Contains(err.Error(), "not valid until") {
		t.Errorf("Expected 'not valid until' in error message, got: %v", err)
	}
}

func TestCheckExpiry_NilCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	duration, err := validator.CheckExpiry(nil)
	if err == nil {
		t.Error("Expected error for nil certificate")
	}
	if duration != nil {
		t.Error("Expected nil duration for nil certificate")
	}
}

func TestVerifyHostname_ValidMatch(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	// Test exact match with SAN
	err = validator.VerifyHostname(validCert, "gitea.local")
	if err != nil {
		t.Errorf("VerifyHostname failed for valid hostname: %v", err)
	}

	// Test alternative SAN
	err = validator.VerifyHostname(validCert, "registry.local")
	if err != nil {
		t.Errorf("VerifyHostname failed for alternative SAN: %v", err)
	}
}

func TestVerifyHostname_WildcardMatch(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	wildcardCert, err := testdata.GenerateWildcardCert()
	if err != nil {
		t.Fatalf("Failed to generate wildcard cert: %v", err)
	}

	// Test wildcard subdomain match
	err = validator.VerifyHostname(wildcardCert, "api.example.local")
	if err != nil {
		t.Errorf("VerifyHostname failed for wildcard match: %v", err)
	}

	// Test exact domain match
	err = validator.VerifyHostname(wildcardCert, "example.local")
	if err != nil {
		t.Errorf("VerifyHostname failed for exact domain match: %v", err)
	}
}

func TestVerifyHostname_NoMatch(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	err = validator.VerifyHostname(validCert, "wrong.domain")
	if err == nil {
		t.Error("Expected error for hostname mismatch")
	}
	if !strings.Contains(err.Error(), "does not match certificate") {
		t.Errorf("Expected hostname mismatch error, got: %v", err)
	}
}

func TestVerifyHostname_NilCertOrEmptyHostname(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	// Test nil certificate
	err = validator.VerifyHostname(nil, "test.local")
	if err == nil || !strings.Contains(err.Error(), "certificate cannot be nil") {
		t.Error("Expected 'certificate cannot be nil' error")
	}

	// Test empty hostname
	err = validator.VerifyHostname(validCert, "")
	if err == nil || !strings.Contains(err.Error(), "hostname cannot be empty") {
		t.Error("Expected 'hostname cannot be empty' error")
	}
}

func TestGenerateDiagnostics_ValidCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	diag, err := validator.GenerateDiagnostics(validCert)
	if err != nil {
		t.Fatalf("GenerateDiagnostics failed: %v", err)
	}

	if diag == nil {
		t.Fatal("Expected diagnostics, got nil")
	}

	if diag.Subject == "" {
		t.Error("Expected subject to be set")
	}

	if diag.Issuer == "" {
		t.Error("Expected issuer to be set")
	}

	if len(diag.DNSNames) == 0 {
		t.Error("Expected DNS names to be populated")
	}

	// Should have warnings about self-signed certificate
	foundSelfSignedWarning := false
	for _, warning := range diag.Warnings {
		if strings.Contains(warning, "self-signed") {
			foundSelfSignedWarning = true
			break
		}
	}
	if !foundSelfSignedWarning {
		t.Error("Expected self-signed certificate warning")
	}
}

func TestGenerateDiagnostics_ExpiredCert(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	expiredCert, err := testdata.GenerateExpiredCert()
	if err != nil {
		t.Fatalf("Failed to generate expired cert: %v", err)
	}

	diag, err := validator.GenerateDiagnostics(expiredCert)
	if err != nil {
		t.Fatalf("GenerateDiagnostics failed: %v", err)
	}

	// Should have an expiry validation error
	foundExpiryError := false
	for _, valErr := range diag.ValidationErrors {
		if valErr.Type == "expiry" {
			foundExpiryError = true
			break
		}
	}
	if !foundExpiryError {
		t.Error("Expected expiry validation error")
	}
}

func TestFormatDiagnostics(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	diag, err := validator.GenerateDiagnostics(validCert)
	if err != nil {
		t.Fatalf("GenerateDiagnostics failed: %v", err)
	}

	formatted := FormatDiagnostics(diag)

	if formatted == "" {
		t.Error("Expected formatted output, got empty string")
	}

	// Check for expected sections
	expectedSections := []string{
		"Certificate Diagnostic Report",
		"Basic Information:",
		"Valid Hostnames:",
		"Summary:",
	}

	for _, section := range expectedSections {
		if !strings.Contains(formatted, section) {
			t.Errorf("Expected section '%s' in formatted output", section)
		}
	}
}

func TestFormatDiagnostics_NilDiagnostics(t *testing.T) {
	formatted := FormatDiagnostics(nil)
	if formatted != "No diagnostic data available" {
		t.Errorf("Expected 'No diagnostic data available', got: %s", formatted)
	}
}

func TestValidateAndDiagnose(t *testing.T) {
	mockTrustStore := newMockTrustStoreManager()
	validator, _ := NewValidator(mockTrustStore)

	validCert, err := testdata.GenerateValidCert()
	if err != nil {
		t.Fatalf("Failed to generate valid cert: %v", err)
	}

	// Test with valid hostname
	diag, err := validator.ValidateAndDiagnose(validCert, "gitea.local")
	if err != nil {
		t.Fatalf("ValidateAndDiagnose failed: %v", err)
	}

	if diag == nil {
		t.Fatal("Expected diagnostics, got nil")
	}

	// Should have hostname match warning
	foundHostnameMatch := false
	for _, warning := range diag.Warnings {
		if strings.Contains(warning, "matches certificate") {
			foundHostnameMatch = true
			break
		}
	}
	if !foundHostnameMatch {
		t.Error("Expected hostname match confirmation")
	}

	// Test with invalid hostname
	diag, err = validator.ValidateAndDiagnose(validCert, "wrong.domain")
	if err != nil {
		t.Fatalf("ValidateAndDiagnose failed: %v", err)
	}

	// Should have hostname validation error
	foundHostnameError := false
	for _, valErr := range diag.ValidationErrors {
		if valErr.Type == "hostname" {
			foundHostnameError = true
			break
		}
	}
	if !foundHostnameError {
		t.Error("Expected hostname validation error")
	}
}
