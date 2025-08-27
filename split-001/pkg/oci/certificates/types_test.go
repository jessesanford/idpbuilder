package certificates

import (
	"crypto/x509"
	"testing"
)

func TestCertFormat_Constants(t *testing.T) {
	if CertFormatPEM != "PEM" {
		t.Errorf("Expected PEM, got %s", CertFormatPEM)
	}
	if CertFormatDER != "DER" {
		t.Errorf("Expected DER, got %s", CertFormatDER)
	}
	if CertFormatPKCS7 != "PKCS7" {
		t.Errorf("Expected PKCS7, got %s", CertFormatPKCS7)
	}
	if CertFormatPKCS12 != "PKCS12" {
		t.Errorf("Expected PKCS12, got %s", CertFormatPKCS12)
	}
}

func TestValidationStatus_Constants(t *testing.T) {
	if ValidationStatusValid != "valid" {
		t.Errorf("Expected valid, got %s", ValidationStatusValid)
	}
	if ValidationStatusExpired != "expired" {
		t.Errorf("Expected expired, got %s", ValidationStatusExpired)
	}
}

func TestNewCertBundle(t *testing.T) {
	bundle := NewCertBundle(CertFormatPEM)
	if bundle == nil {
		t.Fatal("Expected bundle to be non-nil")
	}
	if bundle.Format != CertFormatPEM {
		t.Errorf("Expected Format to be PEM, got %s", bundle.Format)
	}
	if !bundle.IsEmpty() {
		t.Error("Expected new bundle to be empty")
	}
}

func TestCertBundle_AddCertificate(t *testing.T) {
	bundle := NewCertBundle(CertFormatPEM)
	
	// Test adding a non-CA certificate
	endEntityCert := &x509.Certificate{IsCA: false}
	bundle.AddCertificate(endEntityCert)
	
	if len(bundle.Certificates) != 1 {
		t.Errorf("Expected 1 certificate, got %d", len(bundle.Certificates))
	}
	if len(bundle.CAs) != 0 {
		t.Errorf("Expected 0 CAs, got %d", len(bundle.CAs))
	}
	
	// Test adding a CA certificate
	caCert := &x509.Certificate{IsCA: true}
	bundle.AddCertificate(caCert)
	
	if len(bundle.CAs) != 1 {
		t.Errorf("Expected 1 CA, got %d", len(bundle.CAs))
	}
}

func TestCertBundle_IsEmpty(t *testing.T) {
	bundle := NewCertBundle(CertFormatPEM)
	
	if !bundle.IsEmpty() {
		t.Error("Expected empty bundle to return true for IsEmpty()")
	}
	
	bundle.Certificates = append(bundle.Certificates, &x509.Certificate{})
	if bundle.IsEmpty() {
		t.Error("Expected non-empty bundle to return false for IsEmpty()")
	}
}

func TestCertificateError(t *testing.T) {
	err := NewCertificateError("TEST_ERROR", "test message")
	if err.Code != "TEST_ERROR" {
		t.Errorf("Expected code TEST_ERROR, got %s", err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", err.Message)
	}
	if err.Error() != "TEST_ERROR: test message" {
		t.Errorf("Unexpected error string: %s", err.Error())
	}
}