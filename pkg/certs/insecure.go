package certs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// InsecureMode manages the --insecure flag functionality with proper warnings
type InsecureMode struct {
	logger         *log.Logger
	auditLogger    *log.Logger
	warningsIssued map[string]bool
}

// NewInsecureMode creates a new insecure mode handler
func NewInsecureMode() *InsecureMode {
	return &InsecureMode{
		logger:         log.New(os.Stderr, "[INSECURE] ", log.LstdFlags),
		auditLogger:    log.New(log.Writer(), "[AUDIT] ", log.LstdFlags),
		warningsIssued: make(map[string]bool),
	}
}

// IsInsecureModeRequested checks if the --insecure flag is present in command line arguments
func (im *InsecureMode) IsInsecureModeRequested() bool {
	for _, arg := range os.Args {
		if arg == "--insecure" || arg == "-k" {
			return true
		}
	}
	return false
}

// ApplyInsecureMode enables insecure mode with comprehensive warnings
func (im *InsecureMode) ApplyInsecureMode(ctx context.Context, config *InsecureConfig) error {
	if config.RequireExplicit && !im.IsInsecureModeRequested() {
		return fmt.Errorf("insecure mode requires explicit --insecure flag")
	}

	// Issue warnings for this registry/operation combination
	warningKey := fmt.Sprintf("%s:%s", config.Registry, config.Operation)
	if !im.warningsIssued[warningKey] {
		im.issueSecurityWarnings(config)
		im.warningsIssued[warningKey] = true
	}

	// Log the insecure mode activation
	im.auditLogger.Printf("INSECURE MODE ACTIVATED: registry=%s operation=%s duration=%v reason=%s",
		config.Registry, config.Operation, config.Duration, config.Reason)

	return nil
}

// issueSecurityWarnings displays comprehensive security warnings
func (im *InsecureMode) issueSecurityWarnings(config *InsecureConfig) {
	im.logger.Println("════════════════════════════════════════════════════════")
	im.logger.Println("🚨🚨🚨 SECURITY WARNING: INSECURE MODE ENABLED 🚨🚨🚨")
	im.logger.Println("════════════════════════════════════════════════════════")
	im.logger.Printf("Registry: %s", config.Registry)
	im.logger.Printf("Operation: %s", config.Operation)
	im.logger.Printf("Duration: %v", config.Duration)
	im.logger.Printf("Reason: %s", config.Reason)
	im.logger.Println("")
	im.logger.Println("⚠️  CERTIFICATE VALIDATION IS DISABLED")
	im.logger.Println("⚠️  YOUR CONNECTION MAY NOT BE SECURE")
	im.logger.Println("⚠️  DATA MAY BE INTERCEPTED OR MODIFIED")
	im.logger.Println("")
	im.logger.Println("RISKS INCLUDE:")
	im.logger.Println("• Man-in-the-middle attacks")
	im.logger.Println("• Data interception and modification")
	im.logger.Println("• Malicious registry impersonation")
	im.logger.Println("• Compromised image integrity")
	im.logger.Println("")
	im.logger.Println("RECOMMENDATIONS:")
	im.logger.Println("• Only use in trusted networks (e.g., development environments)")
	im.logger.Println("• Verify registry certificates manually when possible")
	im.logger.Println("• Use secure mode in production environments")
	im.logger.Println("• Consider adding registry certificate to trust store")
	im.logger.Println("")
	im.logger.Println("To resolve permanently:")
	im.logger.Printf("  idpbuilder trust add-registry %s", config.Registry)
	im.logger.Println("")
	im.logger.Println("════════════════════════════════════════════════════════")
}

// PromptUserConsent asks for explicit user confirmation (simulated)
func (im *InsecureMode) PromptUserConsent(config *InsecureConfig) (bool, error) {
	// In a real implementation, this would prompt the user for input
	// For simulation purposes, we'll assume consent if --insecure flag is present
	if im.IsInsecureModeRequested() {
		im.logger.Printf("User consent assumed from --insecure flag")
		return true, nil
	}

	return false, fmt.Errorf("user consent required for insecure mode")
}

// CreateTimeLimitedConfig creates a configuration with time limits
func (im *InsecureMode) CreateTimeLimitedConfig(registry, operation, reason string, duration time.Duration) *InsecureConfig {
	return &InsecureConfig{
		Registry:        registry,
		Operation:       operation,
		Duration:        duration,
		Reason:          reason,
		RequireExplicit: true,
	}
}

// ValidateInsecureConfig validates an insecure configuration
func (im *InsecureMode) ValidateInsecureConfig(config *InsecureConfig) error {
	if config.Registry == "" {
		return fmt.Errorf("registry must be specified for insecure mode")
	}
	
	if config.Operation == "" {
		return fmt.Errorf("operation must be specified for insecure mode")
	}
	
	if config.Duration <= 0 {
		config.Duration = 5 * time.Minute // Default to 5 minutes
	}
	
	if config.Duration > 24*time.Hour {
		return fmt.Errorf("insecure mode duration cannot exceed 24 hours")
	}
	
	if config.Reason == "" {
		config.Reason = "No reason provided"
	}
	
	return nil
}

// GenerateSecurityWarning creates a formatted security warning message
func (im *InsecureMode) GenerateSecurityWarning(config *InsecureConfig) string {
	warning := fmt.Sprintf(`
🚨 SECURITY WARNING: Insecure Mode Active

Registry: %s
Operation: %s
Duration: %v
Reason: %s

Certificate validation is disabled for this operation.
Your connection may not be secure and data may be intercepted.

This should only be used in development environments or with
trusted registries where certificate issues are expected.

To resolve permanently, add the registry certificate to your trust store:
  idpbuilder trust add-registry %s

`, config.Registry, config.Operation, config.Duration, config.Reason, config.Registry)
	
	return warning
}

// IsInsecureAllowed checks if insecure mode is allowed for a given registry
func (im *InsecureMode) IsInsecureAllowed(registry string) bool {
	// In a real implementation, this might check against a whitelist
	// of registries where insecure mode is permitted (e.g., internal dev registries)
	
	// For development/Kind environments, allow common local registries
	allowedPatterns := []string{
		"localhost",
		"127.0.0.1",
		"kind-registry",
		"gitea-http",
		".local",
		".dev",
	}
	
	for _, pattern := range allowedPatterns {
		if len(registry) >= len(pattern) {
			for i := 0; i <= len(registry)-len(pattern); i++ {
				if registry[i:i+len(pattern)] == pattern {
					return true
				}
			}
		}
	}
	
	return false
}

// GetInsecureRecommendations provides recommendations for secure alternatives
func (im *InsecureMode) GetInsecureRecommendations(registry string) []Recommendation {
	return []Recommendation{
		{
			Priority:    PriorityHigh,
			Title:       "Add certificate to trust store",
			Description: fmt.Sprintf("Permanently trust %s's certificate", registry),
			Command:     fmt.Sprintf("idpbuilder trust add-registry %s", registry),
			Link:        "https://docs.example.com/certificate-trust",
		},
		{
			Priority:    PriorityMedium,
			Title:       "Verify certificate manually",
			Description: "Check the certificate details before trusting",
			Command:     fmt.Sprintf("openssl s_client -connect %s:443 -servername %s", registry, registry),
			Link:        "https://docs.example.com/verify-certificates",
		},
		{
			Priority:    PriorityMedium,
			Title:       "Use HTTPS with proper certificates",
			Description: "Configure the registry with valid TLS certificates",
			Command:     "",
			Link:        "https://docs.example.com/registry-tls",
		},
		{
			Priority:    PriorityLow,
			Title:       "Use insecure mode temporarily",
			Description: "Bypass certificate validation (development only)",
			Command:     "idpbuilder --insecure <command>",
			Link:        "https://docs.example.com/insecure-mode",
		},
	}
}