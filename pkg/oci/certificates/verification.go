// Package certificates provides verification mode management for certificate validation.
// This file implements dynamic verification mode switching, custom CA pool management,
// and mode-specific validation strategies with fallback support.
package certificates

import (
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// VerificationManager handles verification mode management and mode-specific validation
type VerificationManager struct {
	// Current configuration
	currentMode  v2.VerificationMode // Active verification mode
	fallbackMode v2.VerificationMode // Fallback mode on verification failure

	// Certificate pools
	strictPool   *x509.CertPool // System-only certificate pool
	customCAPool *x509.CertPool // System + custom CA pool
	skipPool     *x509.CertPool // Minimal pool for skip mode

	// Mode transition tracking
	modeHistory    []ModeTransition // History of mode changes
	lastModeSwitch time.Time        // Timestamp of last mode switch
	switchCount    int              // Number of mode switches

	// Thread safety
	mu sync.RWMutex
}

// ModeTransition represents a verification mode change event
type ModeTransition struct {
	FromMode  v2.VerificationMode `json:"from_mode"`
	ToMode    v2.VerificationMode `json:"to_mode"`
	Timestamp time.Time           `json:"timestamp"`
	Reason    string              `json:"reason"`
	Success   bool                `json:"success"`
}

// ValidationStrategy defines the validation approach for each mode
type ValidationStrategy struct {
	Mode            v2.VerificationMode
	RequireChain    bool
	AllowSelfSigned bool
	CheckHostname   bool
	CheckExpiry     bool
}

// NewVerificationManager creates a new verification manager with default settings
func NewVerificationManager(systemPool, customPool *x509.CertPool) (*VerificationManager, error) {
	if systemPool == nil {
		return nil, fmt.Errorf("system certificate pool is required")
	}

	// Create strict pool (system certificates only)
	strictPool := systemPool.Clone()

	// Create custom CA pool (system + custom certificates)
	customCAPool := systemPool.Clone()
	if customPool != nil {
		// Note: x509.CertPool doesn't provide a way to iterate and copy certificates
		// In a real implementation, we would need to track certificates separately
		customCAPool = customPool.Clone()
	}

	// Create minimal skip pool
	skipPool := x509.NewCertPool()

	manager := &VerificationManager{
		currentMode:    v2.VerificationModeStrict,
		fallbackMode:   v2.VerificationModeCustomCA,
		strictPool:     strictPool,
		customCAPool:   customCAPool,
		skipPool:       skipPool,
		modeHistory:    make([]ModeTransition, 0),
		lastModeSwitch: time.Now(),
		switchCount:    0,
	}

	return manager, nil
}

// SwitchVerificationMode atomically switches to a new verification mode
func (v *VerificationManager) SwitchVerificationMode(newMode v2.VerificationMode) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	oldMode := v.currentMode

	// Validate mode transition
	if !v.isValidModeTransition(oldMode, newMode) {
		return fmt.Errorf("invalid mode transition from %s to %s", oldMode, newMode)
	}

	// Record transition attempt
	transition := ModeTransition{
		FromMode:  oldMode,
		ToMode:    newMode,
		Timestamp: time.Now(),
		Reason:    "manual mode switch",
		Success:   false,
	}

	// Perform mode switch
	v.currentMode = newMode
	v.lastModeSwitch = time.Now()
	v.switchCount++

	// Mark transition as successful
	transition.Success = true
	v.modeHistory = append(v.modeHistory, transition)

	// Trim history if too large
	if len(v.modeHistory) > 100 {
		v.modeHistory = v.modeHistory[len(v.modeHistory)-100:]
	}

	return nil
}

// InitVerificationMode initializes the verification mode with custom settings
func (v *VerificationManager) InitVerificationMode(mode v2.VerificationMode, fallback v2.VerificationMode) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Validate modes
	if !v.isValidMode(mode) {
		return fmt.Errorf("invalid verification mode: %s", mode)
	}
	if !v.isValidMode(fallback) {
		return fmt.Errorf("invalid fallback mode: %s", fallback)
	}

	v.currentMode = mode
	v.fallbackMode = fallback

	// Record initialization
	transition := ModeTransition{
		FromMode:  "",
		ToMode:    mode,
		Timestamp: time.Now(),
		Reason:    "initialization",
		Success:   true,
	}
	v.modeHistory = append(v.modeHistory, transition)

	return nil
}

// CreateCustomCAPool creates a custom certificate pool with additional certificates
func (v *VerificationManager) CreateCustomCAPool(additionalCerts []*x509.Certificate) (*x509.CertPool, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// Start with system pool
	pool := v.strictPool.Clone()

	// Add custom certificates
	for _, cert := range additionalCerts {
		if cert == nil {
			continue
		}

		// Validate certificate before adding
		if err := v.validateCertificateForPool(cert); err != nil {
			return nil, fmt.Errorf("invalid certificate for pool: %w", err)
		}

		pool.AddCert(cert)
	}

	return pool, nil
}

// ValidateWithMode performs validation using the specified mode's strategy
func (v *VerificationManager) ValidateWithMode(cert *x509.Certificate, mode v2.VerificationMode) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	strategy := v.getValidationStrategy(mode)

	// Check certificate expiry
	if strategy.CheckExpiry {
		now := time.Now()
		if cert.NotAfter.Before(now) {
			return fmt.Errorf("certificate has expired")
		}
		if cert.NotBefore.After(now) {
			return fmt.Errorf("certificate is not yet valid")
		}
	}

	// Mode-specific validation
	switch mode {
	case v2.VerificationModeStrict:
		return v.validateStrictMode(cert, strategy)
	case v2.VerificationModeCustomCA:
		return v.validateCustomCAMode(cert, strategy)
	case v2.VerificationModeSkip:
		return v.validateSkipMode(cert, strategy)
	default:
		return fmt.Errorf("unknown verification mode: %s", mode)
	}
}

// HandleVerificationFallback handles fallback when verification fails
func (v *VerificationManager) HandleVerificationFallback(cert *x509.Certificate, originalError error) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Check if fallback is different from current mode
	if v.currentMode == v.fallbackMode {
		return fmt.Errorf("fallback failed: %w", originalError)
	}

	// Record fallback attempt
	transition := ModeTransition{
		FromMode:  v.currentMode,
		ToMode:    v.fallbackMode,
		Timestamp: time.Now(),
		Reason:    fmt.Sprintf("fallback due to: %v", originalError),
		Success:   false,
	}

	// Try validation with fallback mode
	originalMode := v.currentMode
	v.currentMode = v.fallbackMode

	err := v.ValidateWithMode(cert, v.fallbackMode)
	if err != nil {
		// Restore original mode if fallback fails
		v.currentMode = originalMode
		transition.Success = false
		v.modeHistory = append(v.modeHistory, transition)
		return fmt.Errorf("fallback validation failed: %w", err)
	}

	// Fallback successful
	transition.Success = true
	v.modeHistory = append(v.modeHistory, transition)
	v.lastModeSwitch = time.Now()
	v.switchCount++

	return nil
}

// getValidationStrategy returns the validation strategy for a mode
func (v *VerificationManager) getValidationStrategy(mode v2.VerificationMode) ValidationStrategy {
	switch mode {
	case v2.VerificationModeStrict:
		return ValidationStrategy{
			Mode:            mode,
			RequireChain:    true,
			AllowSelfSigned: false,
			CheckHostname:   true,
			CheckExpiry:     true,
		}
	case v2.VerificationModeCustomCA:
		return ValidationStrategy{
			Mode:            mode,
			RequireChain:    true,
			AllowSelfSigned: true,
			CheckHostname:   true,
			CheckExpiry:     true,
		}
	case v2.VerificationModeSkip:
		return ValidationStrategy{
			Mode:            mode,
			RequireChain:    false,
			AllowSelfSigned: true,
			CheckHostname:   false,
			CheckExpiry:     false,
		}
	default:
		// Default to strict mode
		return ValidationStrategy{
			Mode:            v2.VerificationModeStrict,
			RequireChain:    true,
			AllowSelfSigned: false,
			CheckHostname:   true,
			CheckExpiry:     true,
		}
	}
}

// validateStrictMode validates certificate using strict mode (system CA only)
func (v *VerificationManager) validateStrictMode(cert *x509.Certificate, strategy ValidationStrategy) error {
	opts := x509.VerifyOptions{
		Roots: v.strictPool,
	}

	chains, err := cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("strict mode validation failed: %w", err)
	}

	if strategy.RequireChain && len(chains) == 0 {
		return fmt.Errorf("no valid certificate chains found")
	}

	return nil
}

// validateCustomCAMode validates certificate using custom CA mode
func (v *VerificationManager) validateCustomCAMode(cert *x509.Certificate, strategy ValidationStrategy) error {
	opts := x509.VerifyOptions{
		Roots: v.customCAPool,
	}

	chains, err := cert.Verify(opts)
	if err != nil {
		// If validation fails and self-signed is allowed, check for self-signed
		if strategy.AllowSelfSigned && v.isSelfSigned(cert) {
			return nil
		}
		return fmt.Errorf("custom CA mode validation failed: %w", err)
	}

	if strategy.RequireChain && len(chains) == 0 {
		return fmt.Errorf("no valid certificate chains found")
	}

	return nil
}

// validateSkipMode validates certificate using skip mode (minimal checks)
func (v *VerificationManager) validateSkipMode(cert *x509.Certificate, strategy ValidationStrategy) error {
	// Skip mode only performs basic certificate structure validation
	if cert.Raw == nil || len(cert.Raw) == 0 {
		return fmt.Errorf("invalid certificate data")
	}

	// Basic structure checks
	if cert.Subject.String() == "" {
		return fmt.Errorf("certificate has empty subject")
	}

	return nil
}

// isValidMode checks if a verification mode is valid
func (v *VerificationManager) isValidMode(mode v2.VerificationMode) bool {
	switch mode {
	case v2.VerificationModeStrict, v2.VerificationModeCustomCA, v2.VerificationModeSkip:
		return true
	default:
		return false
	}
}

// isValidModeTransition checks if a mode transition is allowed
func (v *VerificationManager) isValidModeTransition(from, to v2.VerificationMode) bool {
	// All transitions are allowed
	return v.isValidMode(from) && v.isValidMode(to)
}

// isSelfSigned checks if a certificate is self-signed
func (v *VerificationManager) isSelfSigned(cert *x509.Certificate) bool {
	return cert.CheckSignatureFrom(cert) == nil
}

// validateCertificateForPool validates that a certificate is suitable for adding to a pool
func (v *VerificationManager) validateCertificateForPool(cert *x509.Certificate) error {
	if cert == nil {
		return fmt.Errorf("certificate cannot be nil")
	}

	if cert.Raw == nil || len(cert.Raw) == 0 {
		return fmt.Errorf("certificate has no raw data")
	}

	if cert.Subject.String() == "" {
		return fmt.Errorf("certificate has empty subject")
	}

	return nil
}

// GetCurrentMode returns the current verification mode (thread-safe)
func (v *VerificationManager) GetCurrentMode() v2.VerificationMode {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.currentMode
}

// GetModeHistory returns the mode transition history (thread-safe)
func (v *VerificationManager) GetModeHistory() []ModeTransition {
	v.mu.RLock()
	defer v.mu.RUnlock()

	// Return copy to prevent modification
	history := make([]ModeTransition, len(v.modeHistory))
	copy(history, v.modeHistory)
	return history
}
