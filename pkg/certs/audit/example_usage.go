package audit

import (
	"crypto/x509"
	"fmt"
	"log"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// ExampleBasicAuditLogging demonstrates basic audit logging functionality.
func ExampleBasicAuditLogging() {
	// Create audit logger with default configuration
	config := DefaultConfig()
	config.LogPath = "/tmp/example-audit.log"

	auditLogger, err := NewFileAuditLogger(config)
	if err != nil {
		log.Fatalf("Failed to create audit logger: %v", err)
	}
	defer auditLogger.Close()

	// Log a successful validation
	validationResult := &certs.ValidationResult{
		Valid:          true,
		Chain:          []*x509.Certificate{}, // Would contain actual certificates
		Errors:         []certs.ValidationError{},
		ValidationTime: time.Now(),
		Metadata: certs.ValidationMetadata{
			ChainLength:     3,
			ValidationFlags: certs.ValidatedSignature | certs.ValidatedExpiry,
		},
	}

	metadata := map[string]interface{}{
		"request_id":     "req-12345",
		"client_ip":      "192.168.1.100",
		"hostname":       "example.com",
		"correlation_id": "corr-67890",
	}

	err = auditLogger.LogValidation(validationResult, metadata)
	if err != nil {
		log.Printf("Failed to log validation: %v", err)
	}

	// Log a validation failure
	failedResult := &certs.ValidationResult{
		Valid: false,
		Chain: []*x509.Certificate{}, // Would contain actual certificates
		Errors: []certs.ValidationError{
			{Type: certs.ErrorTypeExpired, Message: "Certificate expired on 2023-12-01"},
		},
		ValidationTime: time.Now(),
		Metadata: certs.ValidationMetadata{
			ChainLength:     2,
			ValidationFlags: certs.ValidatedSignature,
		},
	}

	failureMetadata := map[string]interface{}{
		"request_id":     "req-12346",
		"hostname":       "example.com",
		"correlation_id": "corr-67891",
	}

	err = auditLogger.LogValidation(failedResult, failureMetadata)
	if err != nil {
		log.Printf("Failed to log validation failure: %v", err)
	}

	fmt.Println("Basic audit logging completed successfully")
}

// ExampleIntegratedChainValidation demonstrates integrating audit logging with chain validation.
func ExampleIntegratedChainValidation() {
	// Set up audit logger
	config := AuditLoggerConfig{
		LogPath:       "/tmp/integrated-audit.log",
		BufferSize:    50,
		AutoFlush:     true,
		FlushInterval: 30 * time.Second,
	}

	auditLogger, err := NewFileAuditLogger(config)
	if err != nil {
		log.Fatalf("Failed to create audit logger: %v", err)
	}
	defer auditLogger.Close()

	correlationID := fmt.Sprintf("chain-validation-%d", time.Now().Unix())
	
	// Log start of validation process
	err = auditLogger.LogEvent("validation_started", "Certificate chain validation initiated", map[string]interface{}{
		"correlation_id": correlationID,
		"chain_length":   3,
		"hostname":       "api.example.com",
	})
	if err != nil {
		log.Printf("Failed to log validation start: %v", err)
	}

	// Create validation result
	result := &certs.ValidationResult{
		Valid: true,
		Chain: []*x509.Certificate{}, // Would contain actual certificates
		Errors: []certs.ValidationError{},
		ValidationTime: time.Now(),
		Metadata: certs.ValidationMetadata{
			ChainLength:     3,
			ValidationFlags: certs.ValidatedSignature | certs.ValidatedExpiry | certs.ValidatedHostname,
		},
	}

	// Log the validation result with context
	validationMetadata := map[string]interface{}{
		"correlation_id":    correlationID,
		"hostname":          "api.example.com",
		"validation_method": "pkix",
		"processing_time_ms": 95,
	}

	err = auditLogger.LogValidation(result, validationMetadata)
	if err != nil {
		log.Printf("Failed to log validation result: %v", err)
	}

	fmt.Println("Integrated chain validation audit completed successfully")
}

// ExampleCustomAuditLogger demonstrates implementing a custom audit logger.
type ConsoleAuditLogger struct{}

// NewConsoleAuditLogger creates a simple console-based audit logger for testing.
func NewConsoleAuditLogger() *ConsoleAuditLogger {
	return &ConsoleAuditLogger{}
}

func (c *ConsoleAuditLogger) LogValidation(result *certs.ValidationResult, metadata map[string]interface{}) error {
	status := "SUCCESS"
	if !result.Valid {
		status = "FAILURE"
	}
	
	fmt.Printf("[AUDIT] VALIDATION %s - Chain: %d certs, Errors: %d\n", 
		status, len(result.Chain), len(result.Errors))
	return nil
}

func (c *ConsoleAuditLogger) LogFailure(err error, context map[string]interface{}) error {
	fmt.Printf("[AUDIT] FAILURE - Error: %s\n", err.Error())
	return nil
}

func (c *ConsoleAuditLogger) LogEvent(eventType string, message string, metadata map[string]interface{}) error {
	fmt.Printf("[AUDIT] EVENT [%s] - %s\n", eventType, message)
	return nil
}

func (c *ConsoleAuditLogger) Flush() error {
	return nil
}

func (c *ConsoleAuditLogger) Close() error {
	return nil
}

// ExampleCustomAuditLoggerUsage demonstrates using the custom console audit logger.
func ExampleCustomAuditLoggerUsage() {
	auditLogger := NewConsoleAuditLogger()
	defer auditLogger.Close()

	result := &certs.ValidationResult{
		Valid:          true,
		Chain:          []*x509.Certificate{},
		Errors:         []certs.ValidationError{},
		ValidationTime: time.Now(),
	}

	metadata := map[string]interface{}{
		"hostname": "console-test.example.com",
	}

	auditLogger.LogValidation(result, metadata)
	auditLogger.LogFailure(fmt.Errorf("demo error"), map[string]interface{}{"demo": true})
	auditLogger.LogEvent("demo_event", "This is a demo system event", map[string]interface{}{"demo": true})

	fmt.Println("Custom audit logger demonstration completed")
}