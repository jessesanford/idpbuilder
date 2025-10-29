package auth_test

import (
	"testing"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/stretchr/testify/assert"
)

// MockProvider is a mock implementation of auth.Provider for testing
type MockProvider struct {
	authenticator authn.Authenticator
	validateErr   error
}

func (m *MockProvider) GetAuthenticator() (authn.Authenticator, error) {
	if m.authenticator == nil {
		return authn.Anonymous, nil
	}
	return m.authenticator, nil
}

func (m *MockProvider) ValidateCredentials() error {
	return m.validateErr
}

func TestMockProvider(t *testing.T) {
	t.Run("mock provider implements Provider interface", func(t *testing.T) {
		var _ auth.Provider = &MockProvider{}
	})

	t.Run("mock provider returns authenticator", func(t *testing.T) {
		mock := &MockProvider{
			authenticator: authn.Anonymous,
		}

		authenticator, err := mock.GetAuthenticator()
		assert.NoError(t, err)
		assert.Equal(t, authn.Anonymous, authenticator)
	})

	t.Run("mock provider validates credentials", func(t *testing.T) {
		mock := &MockProvider{
			validateErr: nil,
		}

		err := mock.ValidateCredentials()
		assert.NoError(t, err)
	})

	t.Run("mock provider returns validation error", func(t *testing.T) {
		validationErr := &auth.CredentialValidationError{
			Field:  "username",
			Reason: "cannot be empty",
		}
		mock := &MockProvider{
			validateErr: validationErr,
		}

		err := mock.ValidateCredentials()
		assert.Error(t, err)
		assert.Equal(t, validationErr, err)
	})
}

func TestCredentialValidationError(t *testing.T) {
	t.Run("error message formatting", func(t *testing.T) {
		err := &auth.CredentialValidationError{
			Field:  "password",
			Reason: "too short",
		}

		assert.Equal(t, "credential validation failed (password): too short", err.Error())
	})
}
