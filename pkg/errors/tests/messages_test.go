package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
)

func TestBuildMessage(t *testing.T) {
	tests := []struct {
		name         string
		errorType    errors.ErrorType
		args         []interface{}
		expectedPart string
	}{
		{
			name:         "cert not found",
			errorType:    errors.ErrCertNotFound,
			args:         []interface{}{"/tmp/test.crt"},
			expectedPart: "/tmp/test.crt",
		},
		{
			name:         "cert expired",
			errorType:    errors.ErrCertExpired,
			args:         []interface{}{"30"},
			expectedPart: "30 days ago",
		},
		{
			name:         "cert untrusted",
			errorType:    errors.ErrCertUntrusted,
			args:         []interface{}{"Unknown CA"},
			expectedPart: "Unknown CA",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := errors.BuildMessage(tt.errorType, tt.args...)
			assert.Contains(t, message, tt.expectedPart)
		})
	}
}

func TestUpdateErrorMessage(t *testing.T) {
	err := errors.NewCertNotFound("/tmp/missing.crt")
	errors.UpdateErrorMessage(err)

	assert.Contains(t, err.Message, "/tmp/missing.crt")
}

func TestUpdateErrorMessageWithExpiredCert(t *testing.T) {
	err := errors.NewCertificateError(errors.ErrCertExpired, "old message")
	err.WithDetail("days_ago", "15")
	
	errors.UpdateErrorMessage(err)

	assert.Contains(t, err.Message, "15 days ago")
}

func TestUpdateErrorMessageWithMismatch(t *testing.T) {
	err := errors.NewCertificateError(errors.ErrCertMismatch, "old message")
	err.WithDetail("cn", "example.com").WithDetail("registry", "registry.local")
	
	errors.UpdateErrorMessage(err)

	assert.Contains(t, err.Message, "example.com")
	assert.Contains(t, err.Message, "registry.local")
}