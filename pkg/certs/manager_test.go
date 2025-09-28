package certs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"testing"
	"time"
)

func TestNewX509Manager(t *testing.T) {
	tests := []struct {
		name      string
		validator Validator
		store     Store
		wantNil   bool
	}{
		{
			name:      "with validator and store",
			validator: NewDefaultValidator(),
			store:     NewMemoryStore(),
			wantNil:   false,
		},
		{
			name:      "with nil validator and store",
			validator: nil,
			store:     nil,
			wantNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewX509Manager(tt.validator, tt.store)
			if (manager == nil) != tt.wantNil {
				t.Errorf("NewX509Manager() = %v, want nil: %v", manager, tt.wantNil)
			}
			if manager != nil {
				if manager.GetValidator() == nil {
					t.Error("NewX509Manager() validator should not be nil")
				}
				if manager.GetStore() == nil {
					t.Error("NewX509Manager() store should not be nil")
				}
			}
		})
	}
}

func TestNewDefaultX509Manager(t *testing.T) {
	manager := NewDefaultX509Manager()
	if manager == nil {
		t.Fatal("NewDefaultX509Manager() should not return nil")
	}
	if manager.GetValidator() == nil {
		t.Error("Default manager should have validator")
	}
	if manager.GetStore() == nil {
		t.Error("Default manager should have store")
	}
}

func TestX509Manager_LoadSystemCerts(t *testing.T) {
	manager := NewDefaultX509Manager()
	ctx := context.Background()

	pool, err := manager.LoadSystemCerts(ctx)
	if err != nil {
		t.Fatalf("LoadSystemCerts() error = %v", err)
	}
	if pool == nil {
		t.Error("LoadSystemCerts() returned nil pool")
	}
}

func TestX509Manager_LoadSystemCerts_CancelledContext(t *testing.T) {
	manager := NewDefaultX509Manager()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := manager.LoadSystemCerts(ctx)
	if err == nil {
		t.Error("LoadSystemCerts() should fail with cancelled context")
	}
	if err != context.Canceled {
		t.Errorf("LoadSystemCerts() error = %v, want %v", err, context.Canceled)
	}
}

func TestX509Manager_CreateTLSConfig_Insecure(t *testing.T) {
	manager := NewDefaultX509Manager()
	ctx := context.Background()

	config, err := manager.CreateTLSConfig(ctx, true)
	if err != nil {
		t.Fatalf("CreateTLSConfig(insecure=true) error = %v", err)
	}
	if config == nil {
		t.Fatal("CreateTLSConfig() returned nil config")
	}
	if !config.InsecureSkipVerify {
		t.Error("Insecure config should have InsecureSkipVerify=true")
	}
}

func TestX509Manager_CreateTLSConfig_Secure(t *testing.T) {
	manager := NewDefaultX509Manager()
	ctx := context.Background()

	config, err := manager.CreateTLSConfig(ctx, false)
	if err != nil {
		t.Fatalf("CreateTLSConfig(insecure=false) error = %v", err)
	}
	if config == nil {
		t.Fatal("CreateTLSConfig() returned nil config")
	}
	if config.InsecureSkipVerify {
		t.Error("Secure config should have InsecureSkipVerify=false")
	}
	if config.RootCAs == nil {
		t.Error("Secure config should have RootCAs set")
	}
	if config.MinVersion != tls.VersionTLS12 {
		t.Errorf("Secure config MinVersion = %v, want %v", config.MinVersion, tls.VersionTLS12)
	}
}

func TestX509Manager_ValidateCertificate_NilCert(t *testing.T) {
	manager := NewDefaultX509Manager()
	ctx := context.Background()

	err := manager.ValidateCertificate(ctx, nil)
	if err == nil {
		t.Error("ValidateCertificate() should fail with nil certificate")
	}

	certErr, ok := err.(*CertError)
	if !ok {
		t.Errorf("ValidateCertificate() error type = %T, want *CertError", err)
	} else if certErr.Code != ErrInvalidCert {
		t.Errorf("ValidateCertificate() error code = %v, want %v", certErr.Code, ErrInvalidCert)
	}
}

func TestX509Manager_AddTrustedCert_NilCert(t *testing.T) {
	manager := NewDefaultX509Manager()
	ctx := context.Background()

	err := manager.AddTrustedCert(ctx, nil)
	if err == nil {
		t.Error("AddTrustedCert() should fail with nil certificate")
	}

	certErr, ok := err.(*CertError)
	if !ok {
		t.Errorf("AddTrustedCert() error type = %T, want *CertError", err)
	} else if certErr.Code != ErrInvalidCert {
		t.Errorf("AddTrustedCert() error code = %v, want %v", certErr.Code, ErrInvalidCert)
	}
}

// createTestCertificate creates a minimal test certificate for testing purposes
func createTestCertificate(t *testing.T) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: nil,
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
}