package auth_test

import (
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
)

// T1.1.3-001: AuthProvider interface compiles
func TestAuthProviderInterfaceCompiles(t *testing.T) {
	var _ auth.AuthProvider = nil
}

// T1.1.3-002: InvalidCredentialsError implements error
func TestInvalidCredentialsError_ImplementsError(t *testing.T) {
	err := &auth.InvalidCredentialsError{Reason: "password too short"}
	var _ error = err

	expected := "invalid credentials: password too short"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

// T1.1.3-003: MissingCredentialsError implements error
func TestMissingCredentialsError_ImplementsError(t *testing.T) {
	err := &auth.MissingCredentialsError{Field: "username"}
	var _ error = err

	expected := "missing required credential: username"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

// T1.1.3-004: NewAuthProvider constructor signature valid
func TestNewAuthProvider_SignatureValid(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "not implemented" {
				t.Errorf("Expected panic 'not implemented', got %v", r)
			}
		} else {
			t.Error("Expected NewAuthProvider to panic (not implemented)")
		}
	}()

	_, _ = auth.NewAuthProvider("", "")
}
