package fallback

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// mockValidator implements CertValidator for testing
type mockValidator struct {
	validateChainError error
}

func (m *mockValidator) ValidateChain(cert *x509.Certificate) error {
	return m.validateChainError
}

func (m *mockValidator) CheckExpiry(cert *x509.Certificate) (*time.Duration, error) {
	return nil, nil
}

func (m *mockValidator) VerifyHostname(cert *x509.Certificate, hostname string) error {
	return nil
}

func (m *mockValidator) GenerateDiagnostics(cert *x509.Certificate) (*certs.CertDiagnostics, error) {
	return &certs.CertDiagnostics{}, nil
}

func TestNewDetector(t *testing.T) {
	validator := &mockValidator{}
	detector := NewDetector(validator)

	if detector == nil {
		t.Fatal("NewDetector returned nil")
	}

	if detector.validator != validator {
		t.Error("Detector validator not set correctly")
	}
}

func TestDetectProblem_NilError(t *testing.T) {
	detector := NewDetector(&mockValidator{})
	cert := &x509.Certificate{}

	problem, err := detector.DetectProblem(nil, cert)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if problem != nil {
		t.Error("Expected nil problem for nil error")
	}
}

func TestDetectProblem_NilCert(t *testing.T) {
	detector := NewDetector(&mockValidator{})
	err := errors.New("test error")

	problem, detectErr := detector.DetectProblem(err, nil)

	if detectErr == nil {
		t.Error("Expected error for nil certificate")
	}

	if problem != nil {
		t.Error("Expected nil problem for nil certificate")
	}
}

func TestDetectProblem_SelfSigned(t *testing.T) {
	detector := NewDetector(&mockValidator{})
	cert := &x509.Certificate{
		Subject: pkix.Name{CommonName: "test"},
		Issuer:  pkix.Name{CommonName: "test"}, // Same as subject = self-signed
	}
	err := errors.New("certificate signed by unknown authority")

	problem, detectErr := detector.DetectProblem(err, cert)

	if detectErr != nil {
		t.Errorf("Unexpected error: %v", detectErr)
	}

	if problem == nil {
		t.Fatal("Expected problem to be detected")
	}

	if problem.Type != ProblemSelfSigned {
		t.Errorf("Expected ProblemSelfSigned, got: %s", problem.Type)
	}
}

func TestDetectProblem_Expired(t *testing.T) {
	detector := NewDetector(&mockValidator{})
	cert := &x509.Certificate{}
	err := errors.New("certificate has expired")

	problem, detectErr := detector.DetectProblem(err, cert)

	if detectErr != nil {
		t.Errorf("Unexpected error: %v", detectErr)
	}

	if problem == nil {
		t.Fatal("Expected problem to be detected")
	}

	if problem.Type != ProblemExpired {
		t.Errorf("Expected ProblemExpired, got: %s", problem.Type)
	}
}

func TestDetectProblem_HostnameMismatch(t *testing.T) {
	detector := NewDetector(&mockValidator{})
	cert := &x509.Certificate{
		DNSNames: []string{"example.com"},
	}
	err := errors.New("hostname 'test.com' doesn't match certificate")

	problem, detectErr := detector.DetectProblem(err, cert)

	if detectErr != nil {
		t.Errorf("Unexpected error: %v", detectErr)
	}

	if problem == nil {
		t.Fatal("Expected problem to be detected")
	}

	if problem.Type != ProblemHostnameMismatch {
		t.Errorf("Expected ProblemHostnameMismatch, got: %s", problem.Type)
	}
}

func TestGetProblemSummary(t *testing.T) {
	tests := []struct {
		problemType ProblemType
		expected    string
	}{
		{ProblemSelfSigned, "Certificate is self-signed and not trusted by system"},
		{ProblemExpired, "Certificate has expired"},
		{ProblemHostnameMismatch, "Certificate hostname does not match requested hostname"},
		{ProblemUnknown, "Unknown certificate validation problem"},
	}

	for _, test := range tests {
		problem := &CertProblem{Type: test.problemType}
		summary := problem.GetProblemSummary()

		if summary != test.expected {
			t.Errorf("For %s, expected: %s, got: %s", test.problemType, test.expected, summary)
		}
	}
}
