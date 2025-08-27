package certificates

import (
	"context"
	"crypto/x509"
	"strings"
	"testing"

	v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

const testPEMCert = `-----BEGIN CERTIFICATE-----
MIIDQTCCAimgAwIBAgITBmyfz5m/jAo54vB4ikPmljZbyjANBgkqhkiG9w0BAQsF
ADA5MQswCQYDVQQGEwJVUzEPMA0GA1UEChMGQW1hem9uMRkwFwYDVQQDExBBbWF6
b24gUm9vdCBDQSAxMB4XDTE1MDUyNjAwMDAwMFoXDTM4MDExNzAwMDAwMFowOTEL
MAkGA1UEBhMCVVMxDzANBgNVBAoTBkFtYXpvbjEZMBcGA1UEAxMQQW1hem9uIFJv
b3QgQ0EgMTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALJ4gHHKeNXj
ca9HgFB0fW7Y14h29Jlo91ghYPl0hAEvrAIthtOgQ3pOsqTQNroBvo3bSMgHFzZM
9O6II8c+6zf1tRn4SWiw3te5djgdYZ6k/oI2peVKVuRF4fn9tBb6dNqcmzU5L/qw
IFAGbHrQgLKm+a/sRxmPUDgH3KKHOVj4utWp+UhnMJbulHheb4mjUcAwhmahRWa6
VOujw5H5SNz/0egwLX0tdHA114gk957EWW67c4cX8jJGKLhD+rcdqsq08p8kDi1L
93FcXmn/6pUCyziKrlA4b9v7LWIbxcceVOF34GfID5yHI9Y/QCB/IIDEgEw+OyQm
jgSubJrIqg0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMC
AYYwHQYDVR0OBBYEFIQYzIU07LwMlJQuCFmcx7IQTgoIMA0GCSqGSIb3DQEBCwUA
A4IBAQCY8jdaQZChGsV2USggNiMOruYou6r4lK5IpDB/G/wkjUu0yKGX9rbxenDI
U5PMCCjjmCXPI6T53iHTfIuJruydjsw2hUwsHyx/3GGhU4eq0Xm3LjQgvw7p7Sks
Lc6TxWB6EIvs7ggx4HWBYXQPgNJCZEi1x5gPgpWj5ZZ9PzBhUDL1A7jg4wjIhm8k
bWZZD8LwWIaZ6Ox9m1CmGhZcNzE8G4D2t9sAjOCf3v3Z9j6GdYBc/yMm4IjCp1j4
jZg1s7fKGAfQyA3LJKa4CWqLJ7Q6s6VLTd0KQJJCgA6JVoLCQ4jVZJJj8sVy2SN8
KHXLqDZT5TBlQN7rVEd8KUrG3m4X
-----END CERTIFICATE-----`

func TestBundleLoader_LoadFromReader(t *testing.T) {
	loader := NewBundleLoader()
	ctx := context.Background()

	reader := strings.NewReader(testPEMCert)
	bundle, err := loader.LoadFromReader(ctx, reader)
	if err != nil {
		t.Fatalf("LoadFromReader failed: %v", err)
	}

	if bundle == nil {
		t.Fatal("LoadFromReader returned nil bundle")
	}

	if bundle.Format != v2.CertFormatPEM {
		t.Errorf("Expected format PEM, got %s", bundle.Format)
	}

	// Test with nil reader
	_, err = loader.LoadFromReader(ctx, nil)
	if err == nil {
		t.Error("Expected error for nil reader")
	}
}

func TestBundleLoader_LoadFromData(t *testing.T) {
	loader := NewBundleLoader()
	ctx := context.Background()

	bundle, err := loader.LoadFromData(ctx, []byte(testPEMCert))
	if err != nil {
		t.Fatalf("LoadFromData failed: %v", err)
	}

	if bundle == nil {
		t.Fatal("LoadFromData returned nil bundle")
	}

	// Test with empty data
	_, err = loader.LoadFromData(ctx, []byte{})
	if err == nil {
		t.Error("Expected error for empty data")
	}
}

func TestBundleLoader_DetectFormat(t *testing.T) {
	loader := NewBundleLoader()

	format, err := loader.DetectFormat([]byte(testPEMCert))
	if err != nil {
		t.Fatalf("DetectFormat failed for PEM: %v", err)
	}
	if format != v2.CertFormatPEM {
		t.Errorf("Expected PEM format, got %s", format)
	}

	// Test with empty data
	_, err = loader.DetectFormat([]byte{})
	if err == nil {
		t.Error("Expected error for empty data")
	}
}

func TestValidateCertificate(t *testing.T) {
	err := ValidateCertificate(nil)
	if err == nil {
		t.Error("Expected error for nil certificate")
	}
}

func TestValidateChain(t *testing.T) {
	err := ValidateChain([]*x509.Certificate{})
	if err == nil {
		t.Error("Expected error for empty chain")
	}
}

func BenchmarkBundleLoader_DetectFormat(b *testing.B) {
	loader := NewBundleLoader()
	data := []byte(testPEMCert)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := loader.DetectFormat(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}