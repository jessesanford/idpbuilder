package api

import (
	"context"
	"testing"
	"time"
)

// MockSecurityManager implements SecurityManager for testing
type MockSecurityManager struct {
	SignImageFunc                func(context.Context, string, Signer) (*Signature, error)
	VerifySignatureFunc          func(context.Context, string, Verifier) error
	GenerateSBOMFunc             func(context.Context, string) (*SBOM, error)
	ScanVulnerabilitiesFunc      func(context.Context, string) (*VulnerabilityReport, error)
	AttachAttestationFunc        func(context.Context, string, *Attestation) error
	VerifyAttestationFunc        func(context.Context, string, *Policy) error
	GetImageSecurityProfileFunc  func(context.Context, string) (*SecurityProfile, error)
	ValidatePolicyFunc           func(*Policy) error
	EnforcePolicyFunc            func(context.Context, string, *Policy) (*PolicyResult, error)
}

func (m *MockSecurityManager) SignImage(ctx context.Context, image string, signer Signer) (*Signature, error) {
	if m.SignImageFunc != nil {
		return m.SignImageFunc(ctx, image, signer)
	}
	return &Signature{}, nil
}

func (m *MockSecurityManager) VerifySignature(ctx context.Context, image string, verifier Verifier) error {
	if m.VerifySignatureFunc != nil {
		return m.VerifySignatureFunc(ctx, image, verifier)
	}
	return nil
}

func (m *MockSecurityManager) GenerateSBOM(ctx context.Context, image string) (*SBOM, error) {
	if m.GenerateSBOMFunc != nil {
		return m.GenerateSBOMFunc(ctx, image)
	}
	return &SBOM{}, nil
}

func (m *MockSecurityManager) ScanVulnerabilities(ctx context.Context, image string) (*VulnerabilityReport, error) {
	if m.ScanVulnerabilitiesFunc != nil {
		return m.ScanVulnerabilitiesFunc(ctx, image)
	}
	return &VulnerabilityReport{}, nil
}

func (m *MockSecurityManager) AttachAttestation(ctx context.Context, image string, attestation *Attestation) error {
	if m.AttachAttestationFunc != nil {
		return m.AttachAttestationFunc(ctx, image, attestation)
	}
	return nil
}

func (m *MockSecurityManager) VerifyAttestation(ctx context.Context, image string, policy *Policy) error {
	if m.VerifyAttestationFunc != nil {
		return m.VerifyAttestationFunc(ctx, image, policy)
	}
	return nil
}

func (m *MockSecurityManager) GetImageSecurityProfile(ctx context.Context, image string) (*SecurityProfile, error) {
	if m.GetImageSecurityProfileFunc != nil {
		return m.GetImageSecurityProfileFunc(ctx, image)
	}
	return &SecurityProfile{}, nil
}

func (m *MockSecurityManager) ValidatePolicy(policy *Policy) error {
	if m.ValidatePolicyFunc != nil {
		return m.ValidatePolicyFunc(policy)
	}
	return nil
}

func (m *MockSecurityManager) EnforcePolicy(ctx context.Context, image string, policy *Policy) (*PolicyResult, error) {
	if m.EnforcePolicyFunc != nil {
		return m.EnforcePolicyFunc(ctx, image, policy)
	}
	return &PolicyResult{}, nil
}

// MockSigner implements Signer for testing
type MockSigner struct {
	SignFunc                func([]byte) ([]byte, error)
	KeyIDFunc               func() string
	AlgorithmFunc           func() string
	PublicKeyFunc           func() ([]byte, error)
	GetCertificateChainFunc func() ([]*Certificate, error)
}

func (m *MockSigner) Sign(payload []byte) ([]byte, error) {
	if m.SignFunc != nil {
		return m.SignFunc(payload)
	}
	return []byte("mock-signature"), nil
}

func (m *MockSigner) KeyID() string {
	if m.KeyIDFunc != nil {
		return m.KeyIDFunc()
	}
	return "mock-key-id"
}

func (m *MockSigner) Algorithm() string {
	if m.AlgorithmFunc != nil {
		return m.AlgorithmFunc()
	}
	return "RS256"
}

func (m *MockSigner) PublicKey() ([]byte, error) {
	if m.PublicKeyFunc != nil {
		return m.PublicKeyFunc()
	}
	return []byte("mock-public-key"), nil
}

func (m *MockSigner) GetCertificateChain() ([]*Certificate, error) {
	if m.GetCertificateChainFunc != nil {
		return m.GetCertificateChainFunc()
	}
	return []*Certificate{}, nil
}

// Test interface compliance
func TestSecurityManagerInterface(t *testing.T) {
	var _ SecurityManager = &MockSecurityManager{}
}

func TestSignerInterface(t *testing.T) {
	var _ Signer = &MockSigner{}
}

func TestVulnerability(t *testing.T) {
	vuln := &Vulnerability{
		ID:           "CVE-2023-1234",
		Severity:     "HIGH",
		Score:        8.5,
		Title:        "Test Vulnerability",
		Description:  "A test vulnerability for demonstration",
		Package:      "libssl",
		Version:      "1.0.0",
		FixedVersion: "1.0.1",
		References:   []string{"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2023-1234"},
	}

	if vuln.ID != "CVE-2023-1234" {
		t.Errorf("Expected ID 'CVE-2023-1234', got '%s'", vuln.ID)
	}
	if vuln.Severity != "HIGH" {
		t.Errorf("Expected severity 'HIGH', got '%s'", vuln.Severity)
	}
	if vuln.Score != 8.5 {
		t.Errorf("Expected score 8.5, got %f", vuln.Score)
	}
}