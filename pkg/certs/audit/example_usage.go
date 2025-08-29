package audit

import (
	"fmt"
	"log"
	"time"
)

// ExampleAuditUsage demonstrates how to use the audit logger
func ExampleAuditUsage() {
	// Create audit logger with custom configuration
	config := &AuditLogConfig{
		LogPath:    "/tmp/certificate-audit.log",
		MaxSizeMB:  5, // 5MB max file size before rotation
		CreateDirs: true,
	}

	auditLogger, err := NewAuditLogger(config)
	if err != nil {
		log.Fatalf("Failed to create audit logger: %v", err)
	}
	defer auditLogger.Close()

	// Example 1: Log certificate validation decision
	err = auditLogger.LogCertificateValidation(
		"CN=api.example.com,O=Example Corp",
		"api.example.com",
		"VALID",
		"certificate_chain_verified",
		"user123",
	)
	if err != nil {
		log.Printf("Failed to log certificate validation: %v", err)
	}

	// Example 2: Log trust decision
	err = auditLogger.LogTrustDecision(
		"CN=Root CA,O=Example Corp",
		"TRUSTED",
		"found_in_system_trust_store",
		"admin",
	)
	if err != nil {
		log.Printf("Failed to log trust decision: %v", err)
	}

	// Example 3: Log fallback activation
	err = auditLogger.LogFallbackActivation(
		"CN=expired.example.com",
		"INSECURE_MODE",
		"certificate_expired_user_override",
		"user123",
	)
	if err != nil {
		log.Printf("Failed to log fallback activation: %v", err)
	}

	// Example 4: Log security override
	err = auditLogger.LogSecurityOverride(
		"CN=self-signed.example.com",
		"ACCEPT_SELF_SIGNED",
		"development_environment_exception",
		"developer",
	)
	if err != nil {
		log.Printf("Failed to log security override: %v", err)
	}

	// Example 5: Log error condition
	err = auditLogger.LogError(
		"CN=invalid.example.com",
		"HOSTNAME_MISMATCH",
		"Certificate CN does not match requested hostname",
		"user123",
	)
	if err != nil {
		log.Printf("Failed to log error: %v", err)
	}

	// Example 6: Read recent audit entries
	recentEntries, err := auditLogger.GetRecentEntries(1 * time.Hour)
	if err != nil {
		log.Printf("Failed to get recent entries: %v", err)
	}

	fmt.Printf("Found %d recent audit entries:\n", len(recentEntries))
	for i, entry := range recentEntries {
		fmt.Printf("Entry %d: %s - %s (%s) - %s\n",
			i+1,
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Action,
			entry.Decision,
			entry.Certificate,
		)
	}

	// Example 7: Read all entries with limit
	allEntries, err := auditLogger.ReadEntries(10) // Limit to 10 most recent
	if err != nil {
		log.Printf("Failed to read entries: %v", err)
	}

	fmt.Printf("Total audit entries (last 10): %d\n", len(allEntries))
}

// ExampleChainValidatorWithAudit shows how to integrate audit logging with certificate validation
func ExampleChainValidatorWithAudit() {
	// This would be used in actual implementation:
	
	/* 
	auditLogger, err := NewAuditLogger(DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create audit logger: %v", err)
	}

	config := &ChainValidatorConfig{
		BasicValidator:  myBasicValidator,
		TrustManager:    myTrustManager,
		AllowSelfSigned: false,
		AuditLogger:     auditLogger, // This enables audit logging
	}

	validator := NewChainValidator(config)

	// Now all validation operations will be automatically audited
	result, err := validator.ValidateChain(ctx, certificate, intermediates)
	// Audit log entry is automatically created with validation result

	err = validator.VerifyHostname(certificate, "api.example.com")
	// Audit log entry is automatically created with hostname verification result
	*/

	fmt.Println("See the ChainValidator implementation for integration example")
}