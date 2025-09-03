// Package certs provides certificate handling and TLS configuration for OCI registry operations
package certs

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// InsecureTransport provides an HTTP transport configured for insecure TLS connections
// This is intended for development environments where self-signed certificates are common
type InsecureTransport struct {
	// Base transport to wrap - defaults to http.DefaultTransport
	Base http.RoundTripper

	// Configuration for insecure behavior
	config *InsecureConfig

	// Mutex for thread safety
	mu sync.RWMutex

	// Audit logger for security tracking
	auditLogger *log.Logger
}

// InsecureConfig holds configuration for insecure transport behavior
type InsecureConfig struct {
	// Whether insecure mode is enabled
	Enabled bool

	// Whether to show warnings to user
	ShowWarnings bool

	// Whether to log all insecure connections
	AuditConnections bool

	// Registry URL patterns that are allowed to be insecure
	AllowedRegistries []string

	// Maximum number of insecure connections to track
	MaxAuditEntries int

	// Current audit entries
	auditEntries []AuditEntry

	// Mutex for thread-safe audit operations
	auditMu sync.Mutex
}

// AuditEntry represents a single insecure connection audit record
type AuditEntry struct {
	Timestamp   time.Time
	RegistryURL string
	Method      string
	Status      string
	Warning     string
}

// NewInsecureTransport creates a new InsecureTransport with default configuration
func NewInsecureTransport() *InsecureTransport {
	return &InsecureTransport{
		Base: http.DefaultTransport,
		config: &InsecureConfig{
			Enabled:           false,
			ShowWarnings:      true,
			AuditConnections:  true,
			AllowedRegistries: []string{},
			MaxAuditEntries:   100,
			auditEntries:      make([]AuditEntry, 0),
		},
		auditLogger: log.New(os.Stderr, "[INSECURE-TRANSPORT] ", log.LstdFlags),
	}
}

// NewInsecureTransportWithConfig creates an InsecureTransport with custom configuration
func NewInsecureTransportWithConfig(config *InsecureConfig) *InsecureTransport {
	if config == nil {
		config = &InsecureConfig{
			Enabled:           false,
			ShowWarnings:      true,
			AuditConnections:  true,
			AllowedRegistries: []string{},
			MaxAuditEntries:   100,
			auditEntries:      make([]AuditEntry, 0),
		}
	}

	return &InsecureTransport{
		Base:        http.DefaultTransport,
		config:      config,
		auditLogger: log.New(os.Stderr, "[INSECURE-TRANSPORT] ", log.LstdFlags),
	}
}

// EnableInsecure enables insecure mode with optional registry whitelist
func (it *InsecureTransport) EnableInsecure(allowedRegistries []string) {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.config.Enabled = true
	it.config.AllowedRegistries = allowedRegistries

	if it.config.ShowWarnings {
		it.showInsecureWarning()
	}

	it.auditLogger.Printf("Insecure mode enabled with %d allowed registries", len(allowedRegistries))
}

// DisableInsecure disables insecure mode
func (it *InsecureTransport) DisableInsecure() {
	it.mu.Lock()
	defer it.mu.Unlock()

	it.config.Enabled = false
	it.auditLogger.Printf("Insecure mode disabled")
}

// IsInsecureEnabled returns whether insecure mode is currently enabled
func (it *InsecureTransport) IsInsecureEnabled() bool {
	it.mu.RLock()
	defer it.mu.RUnlock()
	return it.config.Enabled
}

// RoundTrip implements the http.RoundTripper interface
func (it *InsecureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Check if insecure mode is enabled
	it.mu.RLock()
	enabled := it.config.Enabled
	showWarnings := it.config.ShowWarnings
	auditConnections := it.config.AuditConnections
	allowedRegistries := make([]string, len(it.config.AllowedRegistries))
	copy(allowedRegistries, it.config.AllowedRegistries)
	it.mu.RUnlock()

	var transport http.RoundTripper = it.Base

	if enabled {
		registryURL := fmt.Sprintf("%s://%s", req.URL.Scheme, req.URL.Host)

		// Check if registry is in allowed list (if list is provided)
		if len(allowedRegistries) > 0 && !it.isRegistryAllowed(registryURL, allowedRegistries) {
			return nil, fmt.Errorf("registry %s is not in the allowed insecure registry list", registryURL)
		}

		// Configure insecure transport
		transport = it.createInsecureTransport()

		// Show warnings if enabled
		if showWarnings {
			it.showRegistryWarning(registryURL)
		}

		// Audit the connection if enabled
		if auditConnections {
			it.auditConnection(registryURL, req.Method, "ATTEMPTING")
		}

		it.auditLogger.Printf("Making insecure connection to %s", registryURL)
	}

	// Make the request
	resp, err := transport.RoundTrip(clonedReq)

	// Audit the result if insecure and auditing is enabled
	if enabled && it.config.AuditConnections {
		status := "SUCCESS"
		if err != nil {
			status = "ERROR: " + err.Error()
		} else if resp != nil {
			status = fmt.Sprintf("HTTP_%d", resp.StatusCode)
		}
		it.auditConnection(fmt.Sprintf("%s://%s", req.URL.Scheme, req.URL.Host), req.Method, status)
	}

	return resp, err
}

// createInsecureTransport creates an HTTP transport with insecure TLS configuration
func (it *InsecureTransport) createInsecureTransport() http.RoundTripper {
	// Clone the base transport if it's an *http.Transport
	if baseTransport, ok := it.Base.(*http.Transport); ok {
		transport := baseTransport.Clone()
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		} else {
			// Clone the TLS config to avoid modifying the original
			tlsConfig := transport.TLSClientConfig.Clone()
			transport.TLSClientConfig = tlsConfig
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
		return transport
	}

	// Fallback to a new transport if base isn't *http.Transport
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

// isRegistryAllowed checks if a registry URL is in the allowed list
func (it *InsecureTransport) isRegistryAllowed(registryURL string, allowedRegistries []string) bool {
	if len(allowedRegistries) == 0 {
		return true // No restrictions if list is empty
	}

	registryLower := strings.ToLower(registryURL)
	for _, allowed := range allowedRegistries {
		allowedLower := strings.ToLower(allowed)
		if strings.Contains(registryLower, allowedLower) {
			return true
		}
	}
	return false
}

// auditConnection records an insecure connection attempt
func (it *InsecureTransport) auditConnection(registryURL, method, status string) {
	it.config.auditMu.Lock()
	defer it.config.auditMu.Unlock()

	entry := AuditEntry{
		Timestamp:   time.Now(),
		RegistryURL: registryURL,
		Method:      method,
		Status:      status,
		Warning:     "INSECURE CONNECTION - TLS verification disabled",
	}

	// Add to audit log
	it.config.auditEntries = append(it.config.auditEntries, entry)

	// Trim audit log if it exceeds max entries
	if len(it.config.auditEntries) > it.config.MaxAuditEntries {
		// Keep the most recent entries
		start := len(it.config.auditEntries) - it.config.MaxAuditEntries
		it.config.auditEntries = it.config.auditEntries[start:]
	}
}

// GetAuditEntries returns a copy of current audit entries
func (it *InsecureTransport) GetAuditEntries() []AuditEntry {
	it.config.auditMu.Lock()
	defer it.config.auditMu.Unlock()

	// Return a copy to prevent external modification
	entries := make([]AuditEntry, len(it.config.auditEntries))
	copy(entries, it.config.auditEntries)
	return entries
}

// ClearAuditEntries removes all audit entries
func (it *InsecureTransport) ClearAuditEntries() {
	it.config.auditMu.Lock()
	defer it.config.auditMu.Unlock()

	it.config.auditEntries = it.config.auditEntries[:0]
	it.auditLogger.Printf("Audit entries cleared")
}

// showInsecureWarning displays a general insecure mode warning
func (it *InsecureTransport) showInsecureWarning() {
	warning := `
===============================================
🚨 WARNING: INSECURE TLS TRANSPORT ENABLED 🚨
===============================================
⚠️  TLS certificate verification is DISABLED
⚠️  Connections are vulnerable to attacks
⚠️  Only use in development/testing
⚠️  NOT suitable for production
===============================================`
	fmt.Fprintln(os.Stderr, warning)
}

// showRegistryWarning displays a registry-specific warning
func (it *InsecureTransport) showRegistryWarning(registryURL string) {
	warning := fmt.Sprintf(`
🚨 INSECURE CONNECTION WARNING 🚨
Registry: %s
This connection bypasses TLS verification!
⚠️  Certificate validation disabled
⚠️  Man-in-the-middle attacks possible
⚠️  Use only for development/testing
`, registryURL)
	fmt.Fprintln(os.Stderr, warning)
}

// SetBaseTransport allows changing the underlying transport
func (it *InsecureTransport) SetBaseTransport(transport http.RoundTripper) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.Base = transport
}

// GetConfig returns a copy of the current configuration
func (it *InsecureTransport) GetConfig() InsecureConfig {
	it.mu.RLock()
	defer it.mu.RUnlock()

	// Return a copy to prevent external modification
	config := *it.config
	config.AllowedRegistries = make([]string, len(it.config.AllowedRegistries))
	copy(config.AllowedRegistries, it.config.AllowedRegistries)

	// Don't expose the audit entries slice directly
	return config
}

// UpdateConfig updates the transport configuration
func (it *InsecureTransport) UpdateConfig(config *InsecureConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	it.mu.Lock()
	defer it.mu.Unlock()

	// Validate configuration
	if config.MaxAuditEntries < 0 {
		return fmt.Errorf("MaxAuditEntries cannot be negative")
	}

	// Update configuration
	it.config.Enabled = config.Enabled
	it.config.ShowWarnings = config.ShowWarnings
	it.config.AuditConnections = config.AuditConnections
	it.config.MaxAuditEntries = config.MaxAuditEntries

	// Update allowed registries
	it.config.AllowedRegistries = make([]string, len(config.AllowedRegistries))
	copy(it.config.AllowedRegistries, config.AllowedRegistries)

	it.auditLogger.Printf("Configuration updated - Enabled: %v, ShowWarnings: %v, AuditConnections: %v",
		config.Enabled, config.ShowWarnings, config.AuditConnections)

	return nil
}

// CreateInsecureHTTPClient creates an HTTP client using this insecure transport
func (it *InsecureTransport) CreateInsecureHTTPClient(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &http.Client{
		Transport: it,
		Timeout:   timeout,
	}
}

// IsProductionLikeRegistry checks if a registry URL appears to be production
func IsProductionLikeRegistry(registryURL string) bool {
	productionPatterns := []string{
		"docker.io", "gcr.io", "quay.io", "registry-1.docker.io",
		"ghcr.io", "mcr.microsoft.com", "amazonaws.com", "azurecr.io",
		"prod", "production", "live", ".com", ".net", ".org", ".gov",
	}

	urlLower := strings.ToLower(registryURL)
	for _, pattern := range productionPatterns {
		if strings.Contains(urlLower, pattern) {
			return true
		}
	}
	return false
}

// ValidateInsecureUsage validates that insecure mode is being used appropriately
func ValidateInsecureUsage(registryURL string, allowProduction bool) error {
	if registryURL == "" {
		return fmt.Errorf("registry URL cannot be empty")
	}

	if !allowProduction && IsProductionLikeRegistry(registryURL) {
		return fmt.Errorf("insecure mode not recommended for production-like registry: %s", registryURL)
	}

	return nil
}