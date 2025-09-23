package testutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

// Option configures a test authenticator
type Option func(*oci.MockAuthenticator)

// TestRegistry represents a test registry configuration
type TestRegistry struct {
	URL      string
	Username string
	Password string
	cleanup  func()
}

// NewTestAuthenticator creates a configured mock authenticator
func NewTestAuthenticator(opts ...Option) *oci.MockAuthenticator {
	auth := oci.NewMockAuthenticator()
	for _, opt := range opts {
		opt(auth)
	}
	return auth
}

// WithCredentials configures the mock with test credentials
func WithCredentials(creds map[string]*oci.MockCredential) Option {
	return func(auth *oci.MockAuthenticator) {
		for key, cred := range creds {
			auth.AddCredential(key, cred)
		}
	}
}

// WithError configures the mock to return an error for a method
func WithError(method string, err error) Option {
	return func(auth *oci.MockAuthenticator) {
		auth.SetError(method, err)
	}
}

// WithValidation configures custom validation behavior
func WithValidation(valid bool) Option {
	return func(auth *oci.MockAuthenticator) {
		auth.ValidateFunc = func() error {
			if !valid {
				return fmt.Errorf("validation failed")
			}
			return nil
		}
	}
}

// AssertAuthenticated verifies successful authentication
func AssertAuthenticated(t *testing.T, auth *oci.MockAuthenticator, expectedUser string) {
	t.Helper()
	if auth.GetCallCount("Authenticate") == 0 {
		t.Error("Expected Authenticate to be called")
	}

	cred, err := auth.GetCredential("default")
	if err != nil {
		t.Errorf("Expected credential to exist: %v", err)
		return
	}

	if cred.Username != expectedUser {
		t.Errorf("Expected user %s, got %s", expectedUser, cred.Username)
	}
}

// AssertAuthFailed verifies authentication failure
func AssertAuthFailed(t *testing.T, auth *oci.MockAuthenticator, expectedError string) {
	t.Helper()
	if auth.GetCallCount("Authenticate") == 0 {
		t.Error("Expected Authenticate to be called")
	}
}

// AssertCredentialValid verifies credential validity
func AssertCredentialValid(t *testing.T, cred *oci.MockCredential) {
	t.Helper()
	if cred == nil {
		t.Error("Credential is nil")
		return
	}

	if time.Now().After(cred.ValidUntil) {
		t.Error("Credential has expired")
	}

	if cred.Token == "" {
		t.Error("Credential token is empty")
	}
}

// GenerateTestCredentials creates test credentials
func GenerateTestCredentials(n int) []*oci.MockCredential {
	creds := make([]*oci.MockCredential, n)
	for i := 0; i < n; i++ {
		creds[i] = &oci.MockCredential{
			Username:   fmt.Sprintf("user%d", i),
			Password:   fmt.Sprintf("pass%d", i),
			Token:      fmt.Sprintf("token%d", i),
			ValidUntil: time.Now().Add(24 * time.Hour),
			Scopes:     []string{"read", "write"},
		}
	}
	return creds
}

// GenerateExpiredCredential creates an expired test credential
func GenerateExpiredCredential() *oci.MockCredential {
	return &oci.MockCredential{
		Username:   "expired-user",
		Password:   "expired-pass",
		Token:      "expired-token",
		ValidUntil: time.Now().Add(-1 * time.Hour),
		Scopes:     []string{"read"},
	}
}

// GenerateValidToken creates a valid test token
func GenerateValidToken(scopes ...string) string {
	if len(scopes) == 0 {
		scopes = []string{"read"}
	}
	return fmt.Sprintf("valid-token-%d", time.Now().Unix())
}

// SetupTestRegistry creates a test registry configuration
func SetupTestRegistry(t *testing.T) *TestRegistry {
	t.Helper()
	return &TestRegistry{
		URL:      "test-registry.local",
		Username: "test-user",
		Password: "test-pass",
		cleanup:  func() { /* cleanup logic */ },
	}
}

// CleanupTestResources performs test cleanup
func CleanupTestResources(t *testing.T, resources *TestRegistry) {
	t.Helper()
	if resources != nil && resources.cleanup != nil {
		resources.cleanup()
	}
}