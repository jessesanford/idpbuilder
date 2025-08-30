package errors_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
)

func TestNewCertNotFound(t *testing.T) {
	path := "/tmp/certs/missing.crt"
	err := errors.NewCertNotFound(path)

	assert.Equal(t, errors.ErrCertNotFound, err.Type)
	assert.Contains(t, err.Message, path)
	assert.Equal(t, path, err.Details["path"])
	assert.Equal(t, errors.SeverityError, err.Severity)
}

func TestNewCertExpired(t *testing.T) {
	expiry := time.Now().AddDate(0, 0, -30) // 30 days ago
	err := errors.NewCertExpired(expiry)

	assert.Equal(t, errors.ErrCertExpired, err.Type)
	assert.Contains(t, err.Message, "30 days ago")
	assert.Equal(t, "30", err.Details["days_ago"])
	assert.Equal(t, errors.SeverityError, err.Severity)
}

func TestNewCertUntrusted(t *testing.T) {
	issuer := "Unknown CA"
	err := errors.NewCertUntrusted(issuer)

	assert.Equal(t, errors.ErrCertUntrusted, err.Type)
	assert.Contains(t, err.Message, issuer)
	assert.Equal(t, issuer, err.Details["issuer"])
	assert.Equal(t, errors.SeverityWarning, err.Severity)
}

func TestNewCertMismatch(t *testing.T) {
	cn := "example.com"
	registry := "registry.local"
	err := errors.NewCertMismatch(cn, registry)

	assert.Equal(t, errors.ErrCertMismatch, err.Type)
	assert.Contains(t, err.Message, cn)
	assert.Contains(t, err.Message, registry)
	assert.Equal(t, cn, err.Details["cn"])
	assert.Equal(t, registry, err.Details["registry"])
}

func TestNewCertPermission(t *testing.T) {
	path := "/tmp/certs/restricted.crt"
	err := errors.NewCertPermission(path)

	assert.Equal(t, errors.ErrCertPermission, err.Type)
	assert.Contains(t, err.Message, path)
	assert.Equal(t, path, err.Details["path"])
}

func TestWithDetail(t *testing.T) {
	err := errors.NewCertificateError(errors.ErrCertNotFound, "test message")
	err.WithDetail("custom_key", "custom_value")

	assert.Equal(t, "custom_value", err.Details["custom_key"])
}

func TestWithSeverity(t *testing.T) {
	err := errors.NewCertificateError(errors.ErrCertNotFound, "test message")
	err.WithSeverity(errors.SeverityWarning)

	assert.Equal(t, errors.SeverityWarning, err.Severity)
}

func TestErrorInterface(t *testing.T) {
	err := errors.NewCertNotFound("/tmp/test.crt")
	errorString := err.Error()

	assert.Contains(t, errorString, string(errors.ErrCertNotFound))
	assert.Contains(t, errorString, "Certificate not found at: /tmp/test.crt")
}