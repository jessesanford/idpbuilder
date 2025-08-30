package errors_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
)

func TestGetResolutionSteps(t *testing.T) {
	tests := []struct {
		name      string
		errorType errors.ErrorType
		expected  int // expected number of steps
	}{
		{"cert not found", errors.ErrCertNotFound, 3},
		{"cert expired", errors.ErrCertExpired, 3},
		{"cert untrusted", errors.ErrCertUntrusted, 3},
		{"cert mismatch", errors.ErrCertMismatch, 3},
		{"cert permission", errors.ErrCertPermission, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			steps := errors.GetResolutionSteps(tt.errorType)
			assert.Len(t, steps, tt.expected)
			
			// Verify steps contain actionable content
			for _, step := range steps {
				assert.NotEmpty(t, step)
				assert.True(t, len(step) > 10) // Meaningful content
			}
		})
	}
}

func TestFormatResolution(t *testing.T) {
	resolution := errors.FormatResolution(errors.ErrCertNotFound)
	
	assert.Contains(t, resolution, "Resolution Steps:")
	assert.Contains(t, resolution, "1.")
	assert.Contains(t, resolution, "kubectl")
}

func TestAddResolutionToError(t *testing.T) {
	err := errors.NewCertNotFound("/tmp/test.crt")
	errors.AddResolutionToError(err)

	resolution := err.Details["resolution"]
	assert.NotEmpty(t, resolution)
	assert.Contains(t, resolution, "Resolution Steps:")
	assert.Contains(t, resolution, "kubectl")
}

func TestGetResolutionStepsUnknownType(t *testing.T) {
	steps := errors.GetResolutionSteps("UNKNOWN_ERROR")
	
	assert.Len(t, steps, 1)
	assert.Contains(t, strings.ToLower(steps[0]), "no resolution steps available")
}