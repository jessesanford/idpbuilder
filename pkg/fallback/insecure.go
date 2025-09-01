// Package fallback provides insecure mode implementation for development environments
package fallback

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// InsecureConfig manages insecure mode settings and audit trail
type InsecureConfig struct {
	Enabled      bool     // Whether insecure mode is currently enabled
	WarningShown bool     // Whether warning has been displayed to user
	AuditLog     []string // Audit trail of insecure mode usage
	logger       *log.Logger
}

// NewInsecureConfig creates a new insecure configuration with logging
func NewInsecureConfig() *InsecureConfig {
	return &InsecureConfig{
		Enabled:      false,
		WarningShown: false,
		AuditLog:     make([]string, 0),
		logger:       log.New(os.Stderr, "[INSECURE] ", log.LstdFlags),
	}
}

// CreateInsecureTLSConfig creates a TLS configuration that skips all verification
func CreateInsecureTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

// CreateInsecureTLSConfigWithSNI creates insecure TLS config with Server Name Indication
func CreateInsecureTLSConfigWithSNI(serverName string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         serverName,
	}
}

// ShowInsecureWarning displays prominent security warning to the user
func ShowInsecureWarning() {
	warning := `
===============================================
🚨 WARNING: TLS VERIFICATION DISABLED 🚨
===============================================

You are using the --insecure flag which disables
TLS certificate verification. This means:

⚠️  DANGEROUS in production environments
⚠️  Only suitable for development/testing  
⚠️  Vulnerable to man-in-the-middle attacks
⚠️  All TLS security guarantees are void

This should ONLY be used with:
✓ Local Kind clusters
✓ Development registries
✓ Testing environments with self-signed certs

NEVER use --insecure with production registries!

===============================================`
	
	fmt.Fprintln(os.Stderr, warning)
}

// ShowInsecureWarningWithRegistry displays warning with specific registry context
func ShowInsecureWarningWithRegistry(registryURL string) {
	warning := fmt.Sprintf(`
===============================================
🚨 WARNING: TLS VERIFICATION DISABLED 🚨
===============================================

TLS certificate verification has been disabled
for registry: %s

This connection is NOT SECURE and is vulnerable
to man-in-the-middle attacks.

⚠️  All data transmitted can be intercepted
⚠️  Registry identity cannot be verified
⚠️  NEVER use with sensitive data

Use --insecure only for:
✓ Local development with Kind clusters
✓ Testing with self-signed certificates
✓ Non-production environments

===============================================`, registryURL)
	
	fmt.Fprintln(os.Stderr, warning)
}

// ApplyInsecureFlag configures the system based on the --insecure flag value
func (ic *InsecureConfig) ApplyInsecureFlag(flagValue bool, registryURL string) error {
	if flagValue {
		ic.Enabled = true
		
		// Show warning if not shown before
		if !ic.WarningShown {
			if registryURL != "" {
				ShowInsecureWarningWithRegistry(registryURL)
			} else {
				ShowInsecureWarning()
			}
			ic.WarningShown = true
		}
		
		// Add to audit log
		auditEntry := fmt.Sprintf("Insecure mode enabled at %v for registry: %s", 
			time.Now().Format(time.RFC3339), registryURL)
		ic.AuditLog = append(ic.AuditLog, auditEntry)
		
		// Log to system
		ic.logger.Printf("SECURITY WARNING: Insecure mode enabled for %s", registryURL)
		
		return nil
	}
	
	// Flag is false - disable insecure mode
	if ic.Enabled {
		ic.Enabled = false
		auditEntry := fmt.Sprintf("Insecure mode disabled at %v", time.Now().Format(time.RFC3339))
		ic.AuditLog = append(ic.AuditLog, auditEntry)
		ic.logger.Printf("Insecure mode disabled")
	}
	
	return nil
}

// WrapTransportWithInsecure wraps HTTP transport to use insecure TLS when enabled
func WrapTransportWithInsecure(transport http.RoundTripper, insecure bool) http.RoundTripper {
	if !insecure {
		return transport // Return unchanged if not insecure
	}
	
	// Handle different transport types
	switch t := transport.(type) {
	case *http.Transport:
		// Clone the transport to avoid modifying the original
		clonedTransport := t.Clone()
		clonedTransport.TLSClientConfig = CreateInsecureTLSConfig()
		return clonedTransport
		
	default:
		// For unknown transport types, wrap with a new Transport
		return &http.Transport{
			TLSClientConfig: CreateInsecureTLSConfig(),
		}
	}
}

// WrapTransportWithInsecureAndSNI wraps transport with insecure TLS and specific SNI
func WrapTransportWithInsecureAndSNI(transport http.RoundTripper, insecure bool, serverName string) http.RoundTripper {
	if !insecure {
		return transport
	}
	
	switch t := transport.(type) {
	case *http.Transport:
		clonedTransport := t.Clone()
		clonedTransport.TLSClientConfig = CreateInsecureTLSConfigWithSNI(serverName)
		return clonedTransport
		
	default:
		return &http.Transport{
			TLSClientConfig: CreateInsecureTLSConfigWithSNI(serverName),
		}
	}
}

// CreateInsecureHTTPClient creates an HTTP client with insecure TLS configuration
func CreateInsecureHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: CreateInsecureTLSConfig(),
		},
		Timeout: 30 * time.Second,
	}
}

// CreateInsecureHTTPClientWithTimeout creates insecure HTTP client with custom timeout
func CreateInsecureHTTPClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: CreateInsecureTLSConfig(),
		},
		Timeout: timeout,
	}
}

// IsInsecureModeEnabled checks if insecure mode is currently enabled
func (ic *InsecureConfig) IsInsecureModeEnabled() bool {
	return ic.Enabled
}

// GetAuditLog returns the complete audit log of insecure mode usage
func (ic *InsecureConfig) GetAuditLog() []string {
	// Return a copy to prevent modification
	log := make([]string, len(ic.AuditLog))
	copy(log, ic.AuditLog)
	return log
}

// LogInsecureConnection logs when an insecure connection is established
func (ic *InsecureConfig) LogInsecureConnection(registryURL string) {
	auditEntry := fmt.Sprintf("Insecure connection established to %s at %v", 
		registryURL, time.Now().Format(time.RFC3339))
	ic.AuditLog = append(ic.AuditLog, auditEntry)
	ic.logger.Printf("INSECURE CONNECTION: %s", registryURL)
}

// PrintSecurityReminder prints a final security reminder when operations complete
func PrintSecurityReminder() {
	reminder := `
🔒 SECURITY REMINDER:
Remember to use proper TLS certificates in production.
Consider these secure alternatives:
• Add self-signed certificates to your trust store
• Use certificates from a trusted Certificate Authority
• Configure proper certificate validation
`
	fmt.Fprintln(os.Stderr, reminder)
}

// ValidateInsecureUsage checks if insecure mode is being used appropriately
func (ic *InsecureConfig) ValidateInsecureUsage(registryURL string) error {
	if !ic.Enabled {
		return nil // Not using insecure mode
	}
	
	// Check for obviously problematic usage
	if registryURL == "" {
		return fmt.Errorf("cannot validate insecure usage: registry URL is empty")
	}
	
	// Warn about suspicious patterns
	if containsProductionIndicators(registryURL) {
		ic.logger.Printf("WARNING: Insecure mode with production-like registry: %s", registryURL)
		fmt.Fprintf(os.Stderr, "\n🚨 CRITICAL WARNING: Using --insecure with what appears to be a production registry!\n")
		fmt.Fprintf(os.Stderr, "Registry: %s\n", registryURL)
		fmt.Fprintf(os.Stderr, "This is EXTREMELY DANGEROUS and should be avoided.\n\n")
	}
	
	return nil
}

// containsProductionIndicators checks if a URL looks like a production registry
func containsProductionIndicators(url string) bool {
	productionPatterns := []string{
		"prod", "production", "live", "release",
		".com", ".net", ".org", // Public domains
		"docker.io", "gcr.io", "quay.io", // Well-known registries
	}
	
	urlLower := fmt.Sprintf("%s", url) // Convert to lowercase for comparison
	for _, pattern := range productionPatterns {
		if fmt.Sprintf("%s", pattern) != pattern {
			// Simple contains check - avoiding regex for simplicity
			if len(urlLower) >= len(pattern) {
				for i := 0; i <= len(urlLower)-len(pattern); i++ {
					if urlLower[i:i+len(pattern)] == pattern {
						return true
					}
				}
			}
		}
	}
	
	return false
}