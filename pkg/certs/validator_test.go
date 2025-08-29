package certs

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

// Helper to generate test certificate
func generateTestCert(cn string, expired bool) (*x509.Certificate, error) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	now := time.Now()
	notBefore := now.Add(-24 * time.Hour)
	notAfter := now.Add(365 * 24 * time.Hour)
	if expired {
		notAfter = now.Add(-24 * time.Hour)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: cn},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{cn},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(certDER)
}

func TestNewValidator(t *testing.T) {
	if validator := NewValidator(nil); validator == nil || !validator.allowSelfSigned {
		t.Error("NewValidator with nil config should create validator with defaults")
	}
	if validator := NewValidator(&ValidatorConfig{AllowSelfSigned: false}); validator.allowSelfSigned {
		t.Error("NewValidator should use provided config")
	}
}

func TestDefaultValidatorConfig(t *testing.T) {
	config := DefaultValidatorConfig()
	if !config.AllowSelfSigned || len(config.RequiredKeyUsages) != 2 || len(config.RequiredExtKeyUsages) != 1 {
		t.Error("Default config should allow self-signed certs with 2 key usages and 1 ext usage")
	}
}

func TestValidator_ValidateCertificate(t *testing.T) {
	validator := NewValidator(DefaultValidatorConfig())
	
	// Test nil certificate
	if _, err := validator.ValidateCertificate(nil); err == nil {
		t.Error("Expected error for nil certificate")
	}
	
	// Test valid certificate
	if validCert, _ := generateTestCert("test.com", false); validCert != nil {
		if result, _ := validator.ValidateCertificate(validCert); !result.Valid {
			t.Error("Valid certificate should pass validation")
		}
	}
	
	// Test expired certificate  
	if expiredCert, _ := generateTestCert("expired.com", true); expiredCert != nil {
		if result, _ := validator.ValidateCertificate(expiredCert); result.Valid {
			t.Error("Expired certificate should fail validation")
		}
	}
}

func TestValidator_CheckExpiry(t *testing.T) {
	validator := NewValidator(DefaultValidatorConfig())
	
	// Test nil certificate
	if _, err := validator.CheckExpiry(nil, 30); err == nil {
		t.Error("Expected error for nil certificate")
	}
	
	// Test expired certificate
	if expiredCert, _ := generateTestCert("expired.com", true); expiredCert != nil {
		if result, _ := validator.CheckExpiry(expiredCert, 30); !result.Expired {
			t.Error("Expired certificate should be detected as expired")
		}
	}
	
	// Test valid certificate
	if validCert, _ := generateTestCert("valid.com", false); validCert != nil {
		if result, _ := validator.CheckExpiry(validCert, 30); result.Expired {
			t.Error("Valid certificate should not be expired")
		}
	}
}

func TestValidator_ValidateChain(t *testing.T) {
	validator := NewValidator(DefaultValidatorConfig())
	
	// Test empty chain
	if _, err := validator.ValidateChain([]*x509.Certificate{}); err == nil {
		t.Error("Empty chain should return error")
	}
	
	// Test single certificate
	if cert, _ := generateTestCert("test.com", false); cert != nil {
		if result, _ := validator.ValidateChain([]*x509.Certificate{cert}); !result.Valid {
			t.Error("Single valid certificate should pass chain validation")
		}
	}
}

func TestValidator_VerifyHostname(t *testing.T) {
	validator := NewValidator(DefaultValidatorConfig())
	
	// Test nil certificate
	if err := validator.VerifyHostname(nil, "test.com"); err == nil {
		t.Error("Expected error for nil certificate")
	}
	
	// Test empty hostname
	if cert, _ := generateTestCert("test.com", false); cert != nil {
		if err := validator.VerifyHostname(cert, ""); err == nil {
			t.Error("Expected error for empty hostname")
		}
	}
	
	// Test matching hostname
	if cert, _ := generateTestCert("test.com", false); cert != nil {
		if err := validator.VerifyHostname(cert, "test.com"); err != nil {
			t.Error("Matching hostname should not return error")
		}
	}
}

func TestValidator_GenerateDiagnostics(t *testing.T) {
	validator := NewValidator(DefaultValidatorConfig())
	
	// Test nil certificate
	diagnostics := validator.GenerateDiagnostics(nil)
	if diagnostics["error"] == nil {
		t.Error("Expected error field for nil certificate")
	}
	
	// Test valid certificate
	if cert, _ := generateTestCert("test.com", false); cert != nil {
		diagnostics := validator.GenerateDiagnostics(cert)
		expectedFields := []string{"subject", "issuer", "status", "is_self_signed"}
		for _, field := range expectedFields {
			if _, exists := diagnostics[field]; !exists {
				t.Errorf("Missing diagnostic field: %s", field)
			}
		}
	}
}

func TestValidator_KeyUsageStrings(t *testing.T) {
	validator := NewValidator(DefaultValidatorConfig())
	
	result := validator.keyUsageStrings(x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment)
	if len(result) != 2 {
		t.Errorf("Expected 2 key usage strings, got %d", len(result))
	}
}

func TestValidator_ExtKeyUsageStrings(t *testing.T) {
	validator := NewValidator(DefaultValidatorConfig())
	
	usages := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
	result := validator.extKeyUsageStrings(usages)
	if len(result) != 2 {
		t.Errorf("Expected 2 ext key usage strings, got %d", len(result))
	}
}