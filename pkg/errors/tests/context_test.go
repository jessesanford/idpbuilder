package errors_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cnoe-io/idpbuilder/pkg/errors"
)

func TestCaptureBasicContext(t *testing.T) {
	component := "certificate-manager"
	operation := "validate-cert"
	
	ctx := errors.CaptureBasicContext(component, operation)
	
	assert.Equal(t, component, ctx.Component)
	assert.Equal(t, operation, ctx.Operation)
	assert.NotEmpty(t, ctx.OS)
	assert.WithinDuration(t, time.Now(), ctx.Timestamp, time.Second)
	assert.NotEmpty(t, ctx.User)
}

func TestBasicContextToString(t *testing.T) {
	ctx := errors.CaptureBasicContext("test-component", "test-operation")
	
	str := ctx.ToString()
	
	assert.Contains(t, str, "test-component")
	assert.Contains(t, str, "test-operation")
	assert.Contains(t, str, "Component:")
	assert.Contains(t, str, "Operation:")
	assert.Contains(t, str, "OS:")
	assert.Contains(t, str, "User:")
	assert.Contains(t, str, "Time:")
}

func TestEnrichError(t *testing.T) {
	err := errors.NewCertNotFound("/tmp/test.crt")
	
	enriched := errors.EnrichError(err, "cert-validator", "load-certificate")
	
	assert.Equal(t, "cert-validator", enriched.Details["component"])
	assert.Equal(t, "load-certificate", enriched.Details["operation"])
	assert.NotEmpty(t, enriched.Details["os"])
	assert.NotEmpty(t, enriched.Details["user"])
	assert.NotEmpty(t, enriched.Details["timestamp"])
}

func TestEnrichErrorPreservesExistingDetails(t *testing.T) {
	err := errors.NewCertNotFound("/tmp/test.crt")
	err.WithDetail("custom_field", "custom_value")
	
	enriched := errors.EnrichError(err, "test-component", "test-operation")
	
	// Should preserve existing details
	assert.Equal(t, "custom_value", enriched.Details["custom_field"])
	assert.Equal(t, "/tmp/test.crt", enriched.Details["path"])
	
	// Should add context details
	assert.Equal(t, "test-component", enriched.Details["component"])
	assert.Equal(t, "test-operation", enriched.Details["operation"])
}