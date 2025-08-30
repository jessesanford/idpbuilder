package errors_test

import (
	"strings"
	"testing"
	testing_time "time"

	"github.com/stretchr/testify/assert"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
)

func TestFormatConsole(t *testing.T) {
	err := errors.NewCertNotFound("/tmp/test.crt")
	formatted := errors.FormatConsole(err)

	assert.Contains(t, formatted, "CERTIFICATE ERROR")
	assert.Contains(t, formatted, "CERT_NOT_FOUND")
	assert.Contains(t, formatted, "/tmp/test.crt")
	assert.Contains(t, formatted, "━") // Unicode separator
	assert.Contains(t, formatted, "📍") // Details icon
	assert.Contains(t, formatted, "💡") // Resolution icon
}

func TestFormatConsoleSeveritySymbols(t *testing.T) {
	tests := []struct {
		name     string
		severity errors.Severity
		symbol   string
	}{
		{"warning", errors.SeverityWarning, "⚠️"},
		{"error", errors.SeverityError, "❌"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.NewCertificateError(errors.ErrCertNotFound, "test message")
			err.WithSeverity(tt.severity)
			
			formatted := errors.FormatConsole(err)
			assert.Contains(t, formatted, tt.symbol)
		})
	}
}

func TestFormat(t *testing.T) {
	err := errors.NewCertExpired(testing_time.Now().AddDate(0, 0, -15))
	formatted := errors.Format(err)

	assert.Contains(t, formatted, "CERT_EXPIRED")
	assert.Contains(t, formatted, "15 days ago")
	assert.Contains(t, formatted, "Resolution Steps:")
}

func TestFormatMultiple(t *testing.T) {
	errors_list := []*errors.CertificateError{
		errors.NewCertNotFound("/tmp/cert1.crt"),
		errors.NewCertExpired(testing_time.Now().AddDate(0, 0, -5)),
	}

	formatted := errors.FormatMultiple(errors_list)
	
	// Should contain both errors
	assert.Contains(t, formatted, "CERT_NOT_FOUND")
	assert.Contains(t, formatted, "CERT_EXPIRED")
	assert.Contains(t, formatted, "/tmp/cert1.crt")
	assert.Contains(t, formatted, "5 days ago")
	
	// Should be separated by double newlines
	parts := strings.Split(formatted, "\n\n")
	assert.True(t, len(parts) >= 2)
}

func TestFormatMultipleEmpty(t *testing.T) {
	formatted := errors.FormatMultiple([]*errors.CertificateError{})
	assert.Empty(t, formatted)
}

func TestIsImportantDetail(t *testing.T) {
	err := errors.NewCertMismatch("example.com", "registry.local")
	errors.EnrichError(err, "test-component", "test-operation")
	
	formatted := errors.FormatConsole(err)
	
	// Should show important details
	assert.Contains(t, formatted, "Cn: example.com")
	assert.Contains(t, formatted, "Registry: registry.local")
	
	// Should not show internal details
	assert.NotContains(t, formatted, "component:")
	assert.NotContains(t, formatted, "timestamp:")
}

