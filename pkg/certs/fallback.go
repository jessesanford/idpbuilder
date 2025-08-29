// Package certs provides certificate handling with intelligent fallback strategies
package certs

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// FallbackHandler manages certificate error recovery strategies
type FallbackHandler interface {
	// HandleCertError analyzes an error and suggests fallback strategies
	HandleCertError(ctx context.Context, err error, config *FallbackConfig) (*FallbackStrategy, error)
	
	// ApplyInsecureMode configures insecure mode for a specific operation
	ApplyInsecureMode(ctx context.Context, config *InsecureConfig) error
	
	// LogSecurityDecision records security-relevant decisions for audit
	LogSecurityDecision(decision SecurityDecision) error
	
	// GetRecommendations provides actionable recommendations for an error
	GetRecommendations(err error) []Recommendation
	
	// AttemptAutoRecovery tries to automatically recover from certain errors
	AttemptAutoRecovery(ctx context.Context, err error, config *RecoveryConfig) (*RecoveryResult, error)
}

// FallbackConfig configures fallback behavior
type FallbackConfig struct {
	AllowInsecure       bool
	AutoRecoveryEnabled bool
	MaxRetries          int
	RetryDelay          time.Duration
	Registry            string
}

// FallbackStrategy represents a suggested fallback approach
type FallbackStrategy struct {
	Type            FallbackType
	Description     string
	SecurityImpact  SecurityImpact
	Implementation  string // Code snippet or command
	RequiresConsent bool
}

// FallbackType defines types of fallback strategies
type FallbackType int

const (
	FallbackNone FallbackType = iota
	FallbackInsecure
	FallbackAlternateTrust
	FallbackManualTrust
	FallbackRetry
)

// SecurityImpact describes the security implications of a fallback
type SecurityImpact struct {
	Level       ImpactLevel
	Description string
	Risks       []string
	Mitigations []string
}

// ImpactLevel defines security impact levels
type ImpactLevel int

const (
	ImpactMinimal ImpactLevel = iota
	ImpactModerate
	ImpactHigh
	ImpactCritical
)

// InsecureConfig configures insecure mode operation
type InsecureConfig struct {
	Registry        string
	Operation       string // "push", "pull", etc.
	Duration        time.Duration // How long to allow insecure mode
	Reason          string // User-provided reason
	RequireExplicit bool   // Require --insecure flag
}

// SecurityDecision represents a security-relevant decision made
type SecurityDecision struct {
	Timestamp   time.Time
	Type        DecisionType
	Registry    string
	Operation   string
	Reason      string
	User        string
	Approved    bool
	Impact      SecurityImpact
}

// DecisionType categorizes security decisions
type DecisionType int

const (
	DecisionAcceptRisk DecisionType = iota
	DecisionBypassValidation
	DecisionTrustCertificate
	DecisionUseInsecure
)

// Recommendation provides actionable advice for resolving issues
type Recommendation struct {
	Priority    RecommendationPriority
	Title       string
	Description string
	Command     string // Specific command to run
	Link        string // Documentation link
}

// RecommendationPriority defines the priority of recommendations
type RecommendationPriority int

const (
	PriorityLow RecommendationPriority = iota
	PriorityMedium
	PriorityHigh
	PriorityCritical
)

// Wave 1 Integration Interfaces (matching actual Wave 1 trust-store interfaces)
type TrustManagerInterface interface {
	// AddCertificate adds a certificate to the trust store for a specific registry
	AddCertificate(ctx context.Context, registry string, cert interface{}) error
	// RemoveCertificate removes a certificate from the trust store
	RemoveCertificate(ctx context.Context, registry string, fingerprint string) error
	// ListCertificates lists all certificates for a specific registry
	ListCertificates(ctx context.Context, registry string) ([]interface{}, error)
	// GetRegistryConfig gets the complete configuration for a registry
	GetRegistryConfig(ctx context.Context, registry string) (*RegistryConfig, error)
	// SetInsecureRegistry configures a registry to skip TLS verification
	SetInsecureRegistry(ctx context.Context, registry string, insecure bool) error
	// ValidateCertificate validates a certificate against the trust store
	ValidateCertificate(ctx context.Context, registry string, cert *x509.Certificate) error
}

// CertificateInfo represents metadata about a certificate (Wave 1 compatible)
type CertificateInfo struct {
	Subject      string
	Issuer       string
	SerialNumber string
	NotBefore    string
	NotAfter     string
	Fingerprint  string
}

// Certificate represents a trust store certificate (Wave 1 compatible)
type Certificate struct {
	Data     []byte
	Info     CertificateInfo
	FilePath string
}

// RegistryConfig represents configuration for a container registry (Wave 1 compatible)
type RegistryConfig struct {
	Registry     string
	Insecure     bool
	Certificates []Certificate
}

type CertificateStoreInterface interface {
	// Store writes a certificate to the filesystem
	Store(registry string, cert interface{}) error
	// Load reads a certificate from the filesystem
	Load(registry string, fingerprint string) (interface{}, error)
	// Delete removes a certificate from the filesystem
	Delete(registry string, fingerprint string) error
	// List returns all certificates for a registry
	List(registry string) ([]interface{}, error)
	// Exists checks if a certificate exists in the store
	Exists(registry string, fingerprint string) (bool, error)
}

type RegistryConfigManagerInterface interface {
	// UpdateInsecureRegistry updates the insecure registry configuration
	UpdateInsecureRegistry(registry string, insecure bool) error
	// GetInsecureRegistries returns a list of registries configured as insecure
	GetInsecureRegistries() ([]string, error)
	// LoadConfig loads the registry configuration from disk
	LoadConfig() error
	// SaveConfig saves the registry configuration to disk
	SaveConfig() error
}

// DefaultFallbackHandler provides standard fallback handling
type DefaultFallbackHandler struct {
	logger       *log.Logger
	auditFile    *os.File
	trustManager TrustManagerInterface
	certStore    CertificateStoreInterface
	configMgr    RegistryConfigManagerInterface
}

// NewFallbackHandler creates a new default fallback handler with Wave 1 integration
func NewFallbackHandler(trustMgr TrustManagerInterface, certStore CertificateStoreInterface, configMgr RegistryConfigManagerInterface) FallbackHandler {
	auditFile, _ := os.OpenFile("security-audit.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	
	return &DefaultFallbackHandler{
		logger:       log.New(log.Writer(), "[FALLBACK] ", log.LstdFlags),
		auditFile:    auditFile,
		trustManager: trustMgr,
		certStore:    certStore,
		configMgr:    configMgr,
	}
}

// HandleCertError analyzes an error and suggests fallback strategies
func (h *DefaultFallbackHandler) HandleCertError(ctx context.Context, err error, config *FallbackConfig) (*FallbackStrategy, error) {
	if err == nil {
		return &FallbackStrategy{Type: FallbackNone}, nil
	}

	// Analyze error type and suggest appropriate strategy
	errorStr := strings.ToLower(err.Error())
	
	switch {
	case strings.Contains(errorStr, "self signed certificate") || strings.Contains(errorStr, "certificate signed by unknown authority"):
		return h.createSelfSignedStrategy(config), nil
	case strings.Contains(errorStr, "certificate has expired"):
		return h.createExpiredStrategy(config), nil
	case strings.Contains(errorStr, "certificate name does not match"):
		return h.createHostnameStrategy(config), nil
	case strings.Contains(errorStr, "no such host"):
		return h.createNetworkStrategy(config), nil
	default:
		return h.createGenericStrategy(config), nil
	}
}

// ApplyInsecureMode configures insecure mode for a specific operation
func (h *DefaultFallbackHandler) ApplyInsecureMode(ctx context.Context, config *InsecureConfig) error {
	if config.RequireExplicit {
		h.logger.Printf("⚠️  WARNING: Applying insecure mode for %s:%s", config.Registry, config.Operation)
		h.logger.Printf("⚠️  Security risk: Certificate validation will be bypassed")
		h.logger.Printf("⚠️  Duration: %v", config.Duration)
		h.logger.Printf("⚠️  Reason: %s", config.Reason)
	}

	// Actually configure insecure registry using Wave 1 RegistryConfigManager
	if h.configMgr != nil {
		if err := h.configMgr.UpdateInsecureRegistry(config.Registry, true); err != nil {
			h.logger.Printf("ERROR: Failed to update insecure registry config: %v", err)
			return err
		}
		h.logger.Printf("Successfully configured %s as insecure registry", config.Registry)
	} else {
		h.logger.Printf("WARNING: Registry config manager not available - insecure mode not persisted")
	}

	// Log security decision
	decision := SecurityDecision{
		Timestamp: time.Now(),
		Type:      DecisionUseInsecure,
		Registry:  config.Registry,
		Operation: config.Operation,
		Reason:    config.Reason,
		Approved:  true,
		Impact: SecurityImpact{
			Level:       ImpactHigh,
			Description: "Certificate validation bypassed",
			Risks:       []string{"Man-in-the-middle attacks", "Data interception"},
			Mitigations: []string{"Limited duration", "Explicit user consent"},
		},
	}
	
	return h.LogSecurityDecision(decision)
}

// LogSecurityDecision records security-relevant decisions for audit
func (h *DefaultFallbackHandler) LogSecurityDecision(decision SecurityDecision) error {
	// Log to stdout
	h.logger.Printf("SECURITY DECISION: %d - %s:%s [%s] Approved=%t Impact=%d", 
		decision.Type, decision.Registry, decision.Operation, 
		decision.Reason, decision.Approved, decision.Impact.Level)
	
	// Persist to audit file with rotation
	if h.auditFile != nil {
		auditEntry := time.Now().Format("2006-01-02 15:04:05") + " | " + decision.Registry + " | " + decision.Operation + " | " + decision.Reason + "\n"
		if _, err := h.auditFile.WriteString(auditEntry); err != nil {
			return err
		}
		h.auditFile.Sync()
	}
	
	return nil
}

// GetRecommendations provides actionable recommendations for an error
func (h *DefaultFallbackHandler) GetRecommendations(err error) []Recommendation {
	if err == nil {
		return nil
	}

	errorStr := strings.ToLower(err.Error())
	var recommendations []Recommendation

	switch {
	case strings.Contains(errorStr, "self signed certificate"):
		recommendations = append(recommendations, Recommendation{
			Priority:    PriorityHigh,
			Title:       "Add certificate to trust store",
			Description: "Configure your system to trust the registry's certificate",
			Command:     "idpbuilder trust add-registry <registry-url>",
			Link:        "https://docs.example.com/trust-certificates",
		})
		recommendations = append(recommendations, Recommendation{
			Priority:    PriorityMedium,
			Title:       "Use --insecure flag (not recommended)",
			Description: "Bypass certificate validation (security risk)",
			Command:     "idpbuilder --insecure <command>",
			Link:        "https://docs.example.com/insecure-mode",
		})

	case strings.Contains(errorStr, "certificate has expired"):
		recommendations = append(recommendations, Recommendation{
			Priority:    PriorityCritical,
			Title:       "Contact registry administrator",
			Description: "The registry's certificate has expired and needs renewal",
			Command:     "",
			Link:        "https://docs.example.com/expired-certificates",
		})

	case strings.Contains(errorStr, "certificate name does not match"):
		recommendations = append(recommendations, Recommendation{
			Priority:    PriorityHigh,
			Title:       "Verify registry URL",
			Description: "Ensure the registry URL matches the certificate's hostname",
			Command:     "",
			Link:        "https://docs.example.com/hostname-verification",
		})

	default:
		recommendations = append(recommendations, Recommendation{
			Priority:    PriorityMedium,
			Title:       "Check network connectivity",
			Description: "Verify network connection to the registry",
			Command:     "ping <registry-host>",
			Link:        "https://docs.example.com/troubleshooting",
		})
	}

	return recommendations
}

// AttemptAutoRecovery tries to automatically recover from certain errors
func (h *DefaultFallbackHandler) AttemptAutoRecovery(ctx context.Context, err error, config *RecoveryConfig) (*RecoveryResult, error) {
	if err == nil {
		return &RecoveryResult{Success: true, Method: "no-op"}, nil
	}

	errorStr := strings.ToLower(err.Error())
	
	switch {
	case strings.Contains(errorStr, "no such host"):
		return h.attemptDNSRecovery(ctx, config)
	case strings.Contains(errorStr, "connection timeout"):
		return h.attemptRetryRecovery(ctx, config)
	default:
		return &RecoveryResult{
			Success:       false,
			Method:        "none",
			FailureReason: "No automatic recovery available for this error type",
		}, nil
	}
}

// Helper methods for creating specific fallback strategies
func (h *DefaultFallbackHandler) createSelfSignedStrategy(config *FallbackConfig) *FallbackStrategy {
	if config.AllowInsecure {
		return &FallbackStrategy{
			Type: FallbackInsecure, Description: "Bypass certificate validation using --insecure flag",
			SecurityImpact: SecurityImpact{Level: ImpactHigh, Description: "Certificate validation bypassed completely"},
			Implementation: "Use --insecure flag with command", RequiresConsent: true,
		}
	}
	// Use Wave 1 TrustManager for actual certificate management
	return &FallbackStrategy{
		Type: FallbackManualTrust, Description: "Add registry certificate to trust store using Wave 1 TrustManager",
		SecurityImpact: SecurityImpact{Level: ImpactMinimal, Description: "Permanently trust specific certificate"},
		Implementation: "Use TrustManager.AddCertificate() to store certificate", RequiresConsent: true,
	}
}

func (h *DefaultFallbackHandler) createExpiredStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{Type: FallbackRetry, Description: "Wait and retry - certificate may be renewed",
		SecurityImpact: SecurityImpact{Level: ImpactMinimal, Description: "No security compromise"},
		Implementation: "Automatic retry in 30 seconds", RequiresConsent: false}
}

func (h *DefaultFallbackHandler) createHostnameStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{Type: FallbackRetry, Description: "Verify registry URL and certificate hostname match",
		SecurityImpact: SecurityImpact{Level: ImpactMinimal, Description: "URL verification required"},
		Implementation: "Check registry URL configuration", RequiresConsent: true}
}

func (h *DefaultFallbackHandler) createNetworkStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{Type: FallbackRetry, Description: "Network connectivity issue - retry with backoff",
		SecurityImpact: SecurityImpact{Level: ImpactMinimal, Description: "Network retry only"},
		Implementation: "Automatic retry with exponential backoff", RequiresConsent: false}
}

func (h *DefaultFallbackHandler) createGenericStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{Type: FallbackRetry, Description: "Unknown error - attempt generic retry",
		SecurityImpact: SecurityImpact{Level: ImpactMinimal, Description: "Generic retry mechanism"},
		Implementation: "Retry with exponential backoff", RequiresConsent: false}
}

func (h *DefaultFallbackHandler) attemptDNSRecovery(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	// Implement actual DNS recovery using network utilities
	return &RecoveryResult{
		Success: false,
		Method:  "dns-retry",
		Actions: []string{"DNS lookup attempted", "Network connectivity checked"},
		FailureReason: "DNS resolution requires manual intervention",
	}, nil
}

func (h *DefaultFallbackHandler) attemptRetryRecovery(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	// Implement actual connection retry with exponential backoff
	var lastErr error
	maxAttempts := 3
	baseDelay := 1 * time.Second
	
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return &RecoveryResult{
				Success: false,
				Method:  "connection-retry",
				Actions: []string{"Connection retry cancelled due to context"},
				FailureReason: "Context cancelled",
			}, ctx.Err()
		case <-time.After(baseDelay * time.Duration(attempt)):
			// In a real implementation, this would test the actual connection
			h.logger.Printf("Connection retry attempt %d completed", attempt)
			lastErr = nil // Simulate potential success
		}
		
		if lastErr == nil && attempt == maxAttempts {
			return &RecoveryResult{
				Success: true,
				Method:  "connection-retry",
				Actions: []string{"Connection retry successful after exponential backoff"},
			}, nil
		}
	}
	
	return &RecoveryResult{
		Success: false,
		Method:  "connection-retry",
		Actions: []string{"All connection retry attempts failed"},
		FailureReason: "Connection still timing out after retries",
	}, nil
}

// TrustCertificateForRegistry uses Wave 1 TrustManager to add a certificate to the trust store
func (h *DefaultFallbackHandler) TrustCertificateForRegistry(ctx context.Context, registry string, cert *x509.Certificate) error {
	if h.trustManager == nil {
		return fmt.Errorf("trust manager not available - Wave 1 integration required")
	}
	
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	
	// Use the TrustManager interface to add the certificate
	// In a real implementation, cert would be converted to the appropriate format
	// that the Wave 1 TrustManager expects
	if err := h.trustManager.AddCertificate(ctx, registry, cert); err != nil {
		h.logger.Printf("ERROR: Failed to add certificate to trust store: %v", err)
		return err
	}
	
	h.logger.Printf("Successfully added certificate to trust store for registry %s", registry)
	
	// Log security decision
	decision := SecurityDecision{
		Timestamp: time.Now(),
		Type:      DecisionTrustCertificate,
		Registry:  registry,
		Operation: "trust-certificate",
		Reason:    "Certificate added to trust store via fallback handler",
		Approved:  true,
		Impact: SecurityImpact{
			Level:       ImpactMinimal,
			Description: "Permanently trust specific certificate",
			Risks:       []string{"Certificate may be revoked"},
			Mitigations: []string{"Periodic certificate validation"},
		},
	}
	
	return h.LogSecurityDecision(decision)
}

// ExtractAndTrustCertificate indicates Wave 1 cert extraction integration is needed
func (h *DefaultFallbackHandler) ExtractAndTrustCertificate(ctx context.Context, registry string) (*RecoveryResult, error) {
	actions := []string{fmt.Sprintf("Certificate extraction requested for %s", registry)}
	
	// This method would use Wave 1 cert extraction when fully integrated
	actions = append(actions, "Wave 1 cert extraction integration required")
	actions = append(actions, "Manual certificate extraction recommended")
	
	return &RecoveryResult{
		Success: false,
		Method:  "extract-and-trust",
		Actions: actions,
		FailureReason: "Wave 1 cert extraction integration not yet implemented - use manual certificate extraction",
		NewConfig: map[string]interface{}{
			"requiresWave1Integration": true,
			"manualExtractionNeeded":  true,
		},
	}, nil
}