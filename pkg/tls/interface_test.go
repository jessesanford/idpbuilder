package tls_test

import (
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/tls"
)

// T1.1.3-005: TLSProvider interface compiles
func TestTLSProviderInterfaceCompiles(t *testing.T) {
	var _ tls.TLSProvider = nil
}

// T1.1.3-006: NewTLSProvider constructor signature valid
func TestNewTLSProvider_SignatureValid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "not implemented" {
				t.Errorf("Expected panic 'not implemented', got %v", r)
			}
		} else {
			t.Error("Expected NewTLSProvider to panic (not implemented)")
		}
	}()

	_, _ = tls.NewTLSProvider(false)
}
