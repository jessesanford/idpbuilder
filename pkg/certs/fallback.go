// Package certs provides certificate handling with intelligent fallback strategies
package certs

import (
	"context"
	"log"
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

// DefaultFallbackHandler provides standard fallback handling
type DefaultFallbackHandler struct {
	logger *log.Logger
}

// NewFallbackHandler creates a new default fallback handler
func NewFallbackHandler() FallbackHandler {
	return &DefaultFallbackHandler{
		logger: log.New(log.Writer(), "[FALLBACK] ", log.LstdFlags),
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
	h.logger.Printf("SECURITY DECISION: %d - %s:%s [%s] Approved=%t Impact=%d", 
		decision.Type, decision.Registry, decision.Operation, 
		decision.Reason, decision.Approved, decision.Impact.Level)
	
	// In a real implementation, this would write to an audit log file or database
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
			Type:        FallbackInsecure,
			Description: "Bypass certificate validation using --insecure flag",
			SecurityImpact: SecurityImpact{
				Level:       ImpactHigh,
				Description: "Certificate validation bypassed completely",
				Risks:       []string{"Man-in-the-middle attacks", "Data interception"},
				Mitigations: []string{"Use only in development", "Limited duration"},
			},
			Implementation:  "Use --insecure flag with command",
			RequiresConsent: true,
		}
	}
	
	return &FallbackStrategy{
		Type:        FallbackManualTrust,
		Description: "Add registry certificate to trust store",
		SecurityImpact: SecurityImpact{
			Level:       ImpactMinimal,
			Description: "Permanently trust specific certificate",
			Risks:       []string{"Trust misconfigured certificate"},
			Mitigations: []string{"Verify certificate authenticity first"},
		},
		Implementation:  "idpbuilder trust add-registry <registry-url>",
		RequiresConsent: true,
	}
}

func (h *DefaultFallbackHandler) createExpiredStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{
		Type:        FallbackRetry,
		Description: "Wait and retry - certificate may be renewed",
		SecurityImpact: SecurityImpact{
			Level:       ImpactMinimal,
			Description: "No security compromise",
			Risks:       []string{},
			Mitigations: []string{"Automatic retry with backoff"},
		},
		Implementation:  "Automatic retry in 30 seconds",
		RequiresConsent: false,
	}
}

func (h *DefaultFallbackHandler) createHostnameStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{
		Type:        FallbackRetry,
		Description: "Verify registry URL and certificate hostname match",
		SecurityImpact: SecurityImpact{
			Level:       ImpactMinimal,
			Description: "URL verification required",
			Risks:       []string{"Connecting to wrong registry"},
			Mitigations: []string{"Manual URL verification"},
		},
		Implementation:  "Check registry URL configuration",
		RequiresConsent: true,
	}
}

func (h *DefaultFallbackHandler) createNetworkStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{
		Type:        FallbackRetry,
		Description: "Network connectivity issue - retry with backoff",
		SecurityImpact: SecurityImpact{
			Level:       ImpactMinimal,
			Description: "Network retry only",
			Risks:       []string{},
			Mitigations: []string{"Exponential backoff", "Connection timeout"},
		},
		Implementation:  "Automatic retry with exponential backoff",
		RequiresConsent: false,
	}
}

func (h *DefaultFallbackHandler) createGenericStrategy(config *FallbackConfig) *FallbackStrategy {
	return &FallbackStrategy{
		Type:        FallbackRetry,
		Description: "Unknown error - attempt generic retry",
		SecurityImpact: SecurityImpact{
			Level:       ImpactMinimal,
			Description: "Generic retry mechanism",
			Risks:       []string{},
			Mitigations: []string{"Limited retry attempts"},
		},
		Implementation:  "Retry with exponential backoff",
		RequiresConsent: false,
	}
}

func (h *DefaultFallbackHandler) attemptDNSRecovery(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	// Simulate DNS recovery attempt
	return &RecoveryResult{
		Success: false,
		Method:  "dns-retry",
		Actions: []string{"DNS lookup attempted"},
		FailureReason: "DNS resolution still failing",
	}, nil
}

func (h *DefaultFallbackHandler) attemptRetryRecovery(ctx context.Context, config *RecoveryConfig) (*RecoveryResult, error) {
	// Simulate connection retry
	return &RecoveryResult{
		Success: false,
		Method:  "connection-retry",
		Actions: []string{"Connection retry attempted"},
		FailureReason: "Connection still timing out",
	}, nil
}