// Package certs provides insecure mode certificate handling for development and testing
package certs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// InsecureMode defines the interface for insecure certificate bypass operations
type InsecureMode interface {
	IsEnabled() bool
	ShouldBypass(err error) bool
	LogSecurityWarning(context string)
	RequiresConsent() bool
	GetWarningMessage(operation string) string
}

// InsecureConfig holds configuration for insecure mode operations
type InsecureConfig struct {
	Enabled          bool
	LogWarnings      bool
	AllowProduction  bool
	RequireConsent   bool
	MaxWarnings      int
	WarningThreshold time.Duration
	AllowedRegistries []string
}

// DefaultInsecureConfig returns sensible default configuration for insecure mode
func DefaultInsecureConfig() *InsecureConfig {
	return &InsecureConfig{
		Enabled:          false, // Disabled by default for security
		LogWarnings:      true,
		AllowProduction:  false, // Never allow in production by default
		RequireConsent:   true,
		MaxWarnings:      10,
		WarningThreshold: 5 * time.Minute,
		AllowedRegistries: []string{}, // Empty whitelist by default
	}
}

// InsecureModeImpl provides concrete implementation of InsecureMode interface
type InsecureModeImpl struct {
	config       *InsecureConfig
	logger       *log.Logger
	warningCount int
	lastWarning  time.Time
	mutex        sync.RWMutex
	environment  string
	consentGiven bool
}

// NewInsecureMode creates a new insecure mode handler with the given configuration
func NewInsecureMode(config *InsecureConfig) InsecureMode {
	if config == nil {
		config = DefaultInsecureConfig()
	}

	mode := &InsecureModeImpl{
		config:      config,
		logger:      log.New(os.Stderr, "[INSECURE-MODE] ", log.LstdFlags),
		environment: detectEnvironment(),
		consentGiven: false,
	}

	// Log initialization
	if config.Enabled {
		mode.logger.Printf("⚠️  INSECURE MODE INITIALIZED - Environment: %s", mode.environment)
		if mode.environment == "production" && !config.AllowProduction {
			mode.logger.Printf("🚨 CRITICAL: Insecure mode disabled in production environment")
			mode.config.Enabled = false
		}
	}

	return mode
}

// IsEnabled returns true if insecure mode is currently enabled
func (im *InsecureModeImpl) IsEnabled() bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	// Never enable in production unless explicitly allowed
	if im.environment == "production" && !im.config.AllowProduction {
		return false
	}

	return im.config.Enabled
}

// ShouldBypass determines if the given error should be bypassed in insecure mode
func (im *InsecureModeImpl) ShouldBypass(err error) bool {
	if !im.IsEnabled() || err == nil {
		return false
	}

	errorStr := strings.ToLower(err.Error())
	
	// Define which errors can be bypassed in insecure mode
	bypassableErrors := []string{
		"certificate signed by unknown authority",
		"self signed certificate",
		"certificate has expired", 
		"certificate name does not match",
		"x509: certificate signed by unknown authority",
		"tls: failed to verify certificate",
		"certificate verify failed",
	}

	for _, bypassableError := range bypassableErrors {
		if strings.Contains(errorStr, bypassableError) {
			im.LogSecurityWarning(fmt.Sprintf("Bypassing certificate error: %s", bypassableError))
			return true
		}
	}

	return false
}

// LogSecurityWarning logs a security warning with throttling
func (im *InsecureModeImpl) LogSecurityWarning(context string) {
	if !im.config.LogWarnings {
		return
	}

	im.mutex.Lock()
	defer im.mutex.Unlock()

	now := time.Now()
	
	// Throttle warnings to prevent log spam
	if now.Sub(im.lastWarning) < im.config.WarningThreshold && im.warningCount >= im.config.MaxWarnings {
		return
	}

	im.warningCount++
	im.lastWarning = now

	im.logger.Printf("🚨 SECURITY WARNING [%d]: %s", im.warningCount, context)
	im.logger.Printf("   Environment: %s | Consent Required: %t | Consent Given: %t", 
		im.environment, im.config.RequireConsent, im.consentGiven)
}

// RequiresConsent returns true if user consent is required for insecure operations
func (im *InsecureModeImpl) RequiresConsent() bool {
	return im.config.RequireConsent && !im.consentGiven
}

// GetWarningMessage returns an appropriate warning message for the given operation
func (im *InsecureModeImpl) GetWarningMessage(operation string) string {
	return fmt.Sprintf(`
⚠️  SECURITY WARNING: Insecure mode is bypassing certificate validation for operation: %s
   
   This reduces security and should only be used in development/testing environments.
   Environment detected: %s
   
   Risks include:
   - Man-in-the-middle attacks
   - Data interception
   - Compromised authenticity
   
   Continue only if you understand these risks.
`, operation, im.environment)
}

// GiveConsent allows user to provide explicit consent for insecure operations
func (im *InsecureModeImpl) GiveConsent() {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	
	im.consentGiven = true
	im.logger.Printf("⚠️  User consent given for insecure mode operations")
}

// detectEnvironment attempts to detect the current environment
func detectEnvironment() string {
	// Check common environment variables
	envVars := []string{"NODE_ENV", "GO_ENV", "ENVIRONMENT", "ENV"}
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			normalized := strings.ToLower(value)
			if normalized == "production" || normalized == "prod" {
				return "production"
			}
			if normalized == "staging" || normalized == "stage" {
				return "staging"
			}
			if normalized == "development" || normalized == "dev" || normalized == "local" {
				return "development"
			}
		}
	}

	// Check for common CI/development indicators
	ciIndicators := []string{"CI", "CONTINUOUS_INTEGRATION", "GITHUB_ACTIONS", "JENKINS_URL"}
	for _, indicator := range ciIndicators {
		if os.Getenv(indicator) != "" {
			return "ci"
		}
	}

	// Check for development indicators
	devIndicators := []string{"KUBECONFIG", "MINIKUBE_ACTIVE_DOCKERD"}
	for _, indicator := range devIndicators {
		if os.Getenv(indicator) != "" {
			return "development"
		}
	}

	return "unknown"
}

// InsecureStrategy implements FallbackStrategy interface for insecure operations
type InsecureStrategy struct {
	name         string
	priority     int
	insecureMode InsecureMode
	logger       *log.Logger
}

// NewInsecureStrategy creates a new insecure fallback strategy
func NewInsecureStrategy(insecureMode InsecureMode) FallbackStrategy {
	if insecureMode == nil {
		insecureMode = NewInsecureMode(DefaultInsecureConfig())
	}

	return &InsecureStrategy{
		name:         "insecure-bypass",
		priority:     15, // Higher than tertiary (10), lower than secondary (50)
		insecureMode: insecureMode,
		logger:       log.New(os.Stderr, "[INSECURE-STRATEGY] ", log.LstdFlags),
	}
}

// Name returns the strategy name
func (is *InsecureStrategy) Name() string {
	return is.name
}

// Priority returns the strategy priority (0 = lowest, used as last resort)
func (is *InsecureStrategy) Priority() int {
	return is.priority
}

// CanHandle determines if this strategy can handle the given error
func (is *InsecureStrategy) CanHandle(err error) bool {
	if !is.insecureMode.IsEnabled() || err == nil {
		return false
	}

	return is.insecureMode.ShouldBypass(err)
}

// Execute performs insecure strategy validation (bypasses certificate checks)
func (is *InsecureStrategy) Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	if !is.insecureMode.IsEnabled() {
		return &ValidationResult{
			Success:       false,
			Strategy:      is.Name(),
			Message:       "Insecure mode is disabled",
			SecurityLevel: SecurityHigh,
			Actions:       []string{"Enable insecure mode or fix certificate issues"},
		}, nil
	}

	if input.Error == nil {
		return &ValidationResult{
			Success:       true,
			Strategy:      is.Name(),
			Message:       "No error - validation passed in insecure mode",
			SecurityLevel: SecurityNone,
			Actions:       []string{"Consider enabling secure mode for production"},
		}, nil
	}

	// Check if consent is required and not given
	if is.insecureMode.RequiresConsent() {
		warning := is.insecureMode.GetWarningMessage(input.Operation)
		return &ValidationResult{
			Success:       false,
			Strategy:      is.Name(),
			Message:       "User consent required for insecure operation",
			SecurityLevel: SecurityNone,
			Actions: []string{
				"Review security warning",
				"Provide explicit consent if acceptable",
				warning,
			},
		}, nil
	}

	// Check if this error can be bypassed
	if !is.insecureMode.ShouldBypass(input.Error) {
		return &ValidationResult{
			Success:       false,
			Strategy:      is.Name(),
			Message:       fmt.Sprintf("Error cannot be bypassed in insecure mode: %v", input.Error),
			SecurityLevel: SecurityNone,
			Actions:       []string{"Fix the underlying issue or contact administrator"},
		}, nil
	}

	// Bypass the certificate error
	is.logger.Printf("🚨 BYPASSING certificate validation for registry: %s", input.Registry)
	is.insecureMode.LogSecurityWarning(fmt.Sprintf("Certificate bypass for %s operation on %s", input.Operation, input.Registry))

	return &ValidationResult{
		Success:       true,
		Strategy:      is.Name(),
		Message:       fmt.Sprintf("Certificate validation bypassed for %s", input.Registry),
		SecurityLevel: SecurityNone,
		Actions: []string{
			fmt.Sprintf("🚨 SECURITY RISK: Certificate validation bypassed"),
			fmt.Sprintf("Registry: %s", input.Registry),
			fmt.Sprintf("Operation: %s", input.Operation),
			fmt.Sprintf("Original error: %v", input.Error),
			"Audit this operation",
			"Fix certificate issues for production use",
		},
		NewConfig: map[string]interface{}{
			"insecure_mode":       true,
			"bypass_certificates": true,
			"security_level":      "none",
			"audit_required":      true,
			"original_error":      input.Error.Error(),
			"registry":           input.Registry,
		},
	}, nil
}

// IsRegistryAllowed checks if a registry is in the insecure mode allowlist
func (is *InsecureStrategy) IsRegistryAllowed(registry string) bool {
	// Get config through type assertion (unsafe but controlled)
	if impl, ok := is.insecureMode.(*InsecureModeImpl); ok {
		impl.mutex.RLock()
		defer impl.mutex.RUnlock()
		
		// If allowlist is empty, all registries are allowed (risky but explicit)
		if len(impl.config.AllowedRegistries) == 0 {
			return true
		}
		
		// Check exact matches and wildcard patterns
		for _, allowed := range impl.config.AllowedRegistries {
			if allowed == registry {
				return true
			}
			// Simple wildcard support (*.domain.com)
			if strings.HasPrefix(allowed, "*.") {
				domain := allowed[2:]
				if strings.HasSuffix(registry, domain) {
					return true
				}
			}
		}
		return false
	}
	
	// If we can't access config, allow all (fallback)
	return true
}