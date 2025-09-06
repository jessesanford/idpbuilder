// Package fallback provides insecure mode implementation for development environments
package fallback

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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

// ShowInsecureWarning displays prominent security warning to the user
func ShowInsecureWarning() {
	warning := `===============================================
🚨 WARNING: TLS VERIFICATION DISABLED 🚨
===============================================
⚠️  DANGEROUS in production environments
⚠️  Only suitable for development/testing  
⚠️  Vulnerable to man-in-the-middle attacks
⚠️  All TLS security guarantees are void

Use --insecure only for local Kind clusters!
===============================================`
	fmt.Fprintln(os.Stderr, warning)
}

// ShowInsecureWarningWithRegistry displays warning with specific registry context
func ShowInsecureWarningWithRegistry(registryURL string) {
	fmt.Fprintf(os.Stderr, `===============================================
🚨 WARNING: TLS VERIFICATION DISABLED 🚨
===============================================
Registry: %s
This connection is NOT SECURE!
⚠️  All data transmitted can be intercepted
⚠️  Registry identity cannot be verified
===============================================
`, registryURL)
}

// ApplyInsecureFlag configures the system based on the --insecure flag value
func (ic *InsecureConfig) ApplyInsecureFlag(flagValue bool, registryURL string) error {
	if flagValue {
		ic.Enabled = true
		if !ic.WarningShown {
			if registryURL != "" {
				ShowInsecureWarningWithRegistry(registryURL)
			} else {
				ShowInsecureWarning()
			}
			ic.WarningShown = true
		}

		auditEntry := fmt.Sprintf("Insecure mode enabled at %v for registry: %s",
			time.Now().Format(time.RFC3339), registryURL)
		ic.AuditLog = append(ic.AuditLog, auditEntry)
		ic.logger.Printf("SECURITY WARNING: Insecure mode enabled for %s", registryURL)
	} else {
		// Flag is false - disable insecure mode if it was enabled
		if ic.Enabled {
			ic.Enabled = false
			auditEntry := fmt.Sprintf("Insecure mode disabled at %v", time.Now().Format(time.RFC3339))
			ic.AuditLog = append(ic.AuditLog, auditEntry)
			ic.logger.Printf("Insecure mode disabled")
		}
	}
	return nil
}

// WrapTransportWithInsecure wraps HTTP transport to use insecure TLS when enabled
func WrapTransportWithInsecure(transport http.RoundTripper, insecure bool) http.RoundTripper {
	if !insecure {
		return transport
	}

	switch t := transport.(type) {
	case *http.Transport:
		clonedTransport := t.Clone()
		clonedTransport.TLSClientConfig = CreateInsecureTLSConfig()
		return clonedTransport
	default:
		return &http.Transport{
			TLSClientConfig: CreateInsecureTLSConfig(),
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

// IsInsecureModeEnabled checks if insecure mode is currently enabled
func (ic *InsecureConfig) IsInsecureModeEnabled() bool {
	return ic.Enabled
}

// LogInsecureConnection logs when an insecure connection is established
func (ic *InsecureConfig) LogInsecureConnection(registryURL string) {
	auditEntry := fmt.Sprintf("Insecure connection established to %s at %v",
		registryURL, time.Now().Format(time.RFC3339))
	ic.AuditLog = append(ic.AuditLog, auditEntry)
	ic.logger.Printf("INSECURE CONNECTION: %s", registryURL)
}

// ValidateInsecureUsage checks if insecure mode is being used appropriately
func (ic *InsecureConfig) ValidateInsecureUsage(registryURL string) error {
	if !ic.Enabled {
		return nil
	}

	if registryURL == "" {
		return fmt.Errorf("cannot validate insecure usage: registry URL is empty")
	}

	// Warn about suspicious patterns
	if containsProductionIndicators(registryURL) {
		ic.logger.Printf("WARNING: Insecure mode with production-like registry: %s", registryURL)
		fmt.Fprintf(os.Stderr, "🚨 CRITICAL WARNING: Using --insecure with production-like registry: %s\n", registryURL)
	}
	return nil
}

// containsProductionIndicators checks if a URL looks like a production registry
func containsProductionIndicators(url string) bool {
	productionPatterns := []string{"prod", "production", "live", ".com", ".net", "docker.io", "gcr.io", "quay.io"}
	urlLower := strings.ToLower(url)
	for _, pattern := range productionPatterns {
		if strings.Contains(urlLower, pattern) {
			return true
		}
	}
	return false
}
