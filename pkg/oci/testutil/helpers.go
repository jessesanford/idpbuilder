package testutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci"
)

// Option configures a test authenticator
type Option func(*oci.MockAuthenticator)

// NewTestAuthenticator creates configured mock authenticator
func NewTestAuthenticator(opts ...Option) *oci.MockAuthenticator {
	auth := oci.NewMockAuthenticator()
	for _, opt := range opts {
		opt(auth)
	}
	return auth
}

// WithCredentials configures mock with test credentials
func WithCredentials(creds map[string]*oci.MockCredential) Option {
	return func(auth *oci.MockAuthenticator) {
		for key, cred := range creds {
			auth.AddCredential(key, cred)
		}
	}
}

// WithError configures mock to return error for method
func WithError(method string, err error) Option {
	return func(auth *oci.MockAuthenticator) {
		auth.SetError(method, err)
	}
}

// AssertAuthenticated verifies successful authentication
func AssertAuthenticated(t *testing.T, auth *oci.MockAuthenticator, expectedUser string) {
	t.Helper()
	if auth.GetCallCount("Authenticate") == 0 {
		t.Error("Expected Authenticate to be called")
	}
	if cred, err := auth.GetCredential("default"); err != nil {
		t.Errorf("Expected credential to exist: %v", err)
	} else if cred.Username != expectedUser {
		t.Errorf("Expected user %s, got %s", expectedUser, cred.Username)
	}
}

// AssertCredentialValid verifies credential validity
func AssertCredentialValid(t *testing.T, cred *oci.MockCredential) {
	t.Helper()
	if cred == nil || time.Now().After(cred.ValidUntil) || cred.Token == "" {
		t.Error("Invalid credential")
	}
}

// GenerateTestCredentials creates n test credentials
func GenerateTestCredentials(n int) []*oci.MockCredential {
	creds := make([]*oci.MockCredential, n)
	for i := 0; i < n; i++ {
		creds[i] = &oci.MockCredential{
			Username:   fmt.Sprintf("user%d", i),
			Token:      fmt.Sprintf("token%d", i),
			ValidUntil: time.Now().Add(24 * time.Hour)}
	}
	return creds
}

// GenerateExpiredCredential creates expired test credential
func GenerateExpiredCredential() *oci.MockCredential {
	return &oci.MockCredential{
		Username:   "expired-user",
		Token:      "expired-token",
		ValidUntil: time.Now().Add(-1 * time.Hour),
	}
}