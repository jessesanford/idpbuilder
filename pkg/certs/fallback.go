// Package certs provides certificate handling with intelligent fallback strategies
package certs

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// FallbackStrategy defines the interface for certificate validation fallback strategies
type FallbackStrategy interface {
	Name() string
	Priority() int
	CanHandle(err error) bool
	Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
}

// FallbackChain manages multiple fallback strategies with chain of responsibility pattern
type FallbackChain interface {
	AddStrategy(strategy FallbackStrategy)
	RemoveStrategy(name string)
	Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
	GetStrategies() []FallbackStrategy
}

// ValidationInput contains the input data for fallback validation
type ValidationInput struct {
	Certificates []*x509.Certificate
	Registry     string
	Operation    string
	Error        error
	Options      map[string]interface{}
}

// ValidationResult contains the result of a fallback validation attempt
type ValidationResult struct {
	Success       bool
	Strategy      string
	Message       string
	SecurityLevel SecurityLevel
	Actions       []string
	NewConfig     map[string]interface{}
}

// SecurityLevel defines the security level of a validation result
type SecurityLevel int

const (
	SecurityHigh SecurityLevel = iota
	SecurityMedium
	SecurityLow
	SecurityNone
)

// FallbackConfig configures fallback behavior and policies
type FallbackConfig struct {
	AllowInsecure  bool
	MaxRetries     int
	RetryDelay     time.Duration
	Registry       string
	RequireConsent bool
}

// RecoveryConfig configures automatic recovery attempts
type RecoveryConfig struct {
	EnableCertRefresh bool
	EnableTrustUpdate bool
	MaxAttempts       int
	Timeout           time.Duration
}

// RecoveryResult describes the outcome of a recovery attempt
type RecoveryResult struct {
	Success       bool
	Method        string
	Actions       []string
	NewConfig     interface{}
	FailureReason string
}

// DefaultFallbackChain implements the FallbackChain interface
type DefaultFallbackChain struct {
	strategies []FallbackStrategy
	mutex      sync.RWMutex
	logger     *log.Logger
}

// NewFallbackChain creates a new fallback chain with default strategies
func NewFallbackChain() FallbackChain {
	chain := &DefaultFallbackChain{
		strategies: make([]FallbackStrategy, 0),
		logger:     log.New(os.Stdout, "[FALLBACK] ", log.LstdFlags),
	}
	
	// Add default strategies in priority order
	chain.AddStrategy(NewPrimaryStrategy())
	chain.AddStrategy(NewSecondaryStrategy())
	chain.AddStrategy(NewTertiaryStrategy())
	
	return chain
}

// NewFallbackChainWithInsecure creates a fallback chain including insecure mode strategy
func NewFallbackChainWithInsecure(insecureConfig *InsecureConfig) FallbackChain {
	chain := &DefaultFallbackChain{
		strategies: make([]FallbackStrategy, 0),
		logger:     log.New(os.Stdout, "[FALLBACK] ", log.LstdFlags),
	}
	
	// Add all strategies in priority order (higher priority first)
	chain.AddStrategy(NewPrimaryStrategy())       // Priority: 100
	chain.AddStrategy(NewSecondaryStrategy())     // Priority: 50  
	
	// Add insecure strategy before tertiary for certificate-specific bypassing
	insecureMode := NewInsecureMode(insecureConfig)
	chain.AddStrategy(NewInsecureStrategy(insecureMode)) // Priority: 15
	
	chain.AddStrategy(NewTertiaryStrategy())      // Priority: 10
	
	return chain
}

// AddStrategy adds a strategy to the chain, maintaining priority order
func (c *DefaultFallbackChain) AddStrategy(strategy FallbackStrategy) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Insert strategy in priority order (higher priority first)
	inserted := false
	for i, existing := range c.strategies {
		if strategy.Priority() > existing.Priority() {
			c.strategies = append(c.strategies[:i], append([]FallbackStrategy{strategy}, c.strategies[i:]...)...)
			inserted = true
			break
		}
	}
	
	if !inserted {
		c.strategies = append(c.strategies, strategy)
	}
	
	c.logger.Printf("Added strategy: %s (priority: %d)", strategy.Name(), strategy.Priority())
}

// RemoveStrategy removes a strategy by name
func (c *DefaultFallbackChain) RemoveStrategy(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	for i, strategy := range c.strategies {
		if strategy.Name() == name {
			c.strategies = append(c.strategies[:i], c.strategies[i+1:]...)
			c.logger.Printf("Removed strategy: %s", name)
			return
		}
	}
}

// Execute runs the fallback chain, trying each strategy in priority order
func (c *DefaultFallbackChain) Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	if len(c.strategies) == 0 {
		return &ValidationResult{
			Success: false,
			Message: "No fallback strategies available",
		}, fmt.Errorf("no fallback strategies configured")
	}
	
	c.logger.Printf("Executing fallback chain for registry %s", input.Registry)
	
	var lastError error
	var attempts []string
	
	for _, strategy := range c.strategies {
		if !strategy.CanHandle(input.Error) {
			c.logger.Printf("Strategy %s cannot handle error: %v", strategy.Name(), input.Error)
			continue
		}
		
		c.logger.Printf("Trying strategy: %s", strategy.Name())
		attempts = append(attempts, strategy.Name())
		
		result, err := strategy.Execute(ctx, input)
		if err != nil {
			c.logger.Printf("Strategy %s failed: %v", strategy.Name(), err)
			lastError = err
			continue
		}
		
		if result.Success {
			c.logger.Printf("Strategy %s succeeded", strategy.Name())
			result.Actions = append(result.Actions, fmt.Sprintf("Attempted strategies: %v", attempts))
			return result, nil
		}
		
		lastError = fmt.Errorf("strategy %s failed: %s", strategy.Name(), result.Message)
	}
	
	return &ValidationResult{
		Success: false,
		Message: fmt.Sprintf("All fallback strategies failed. Last error: %v", lastError),
		Actions: attempts,
	}, lastError
}

// GetStrategies returns a copy of all strategies in the chain
func (c *DefaultFallbackChain) GetStrategies() []FallbackStrategy {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	strategies := make([]FallbackStrategy, len(c.strategies))
	copy(strategies, c.strategies)
	return strategies
}

// PrimaryStrategy implements strict certificate validation with minimal fallback
type PrimaryStrategy struct {
	name     string
	priority int
}

// NewPrimaryStrategy creates a new primary (strict) validation strategy
func NewPrimaryStrategy() FallbackStrategy {
	return &PrimaryStrategy{
		name:     "primary-strict",
		priority: 100,
	}
}

// Name returns the strategy name
func (s *PrimaryStrategy) Name() string {
	return s.name
}

// Priority returns the strategy priority (higher values execute first)
func (s *PrimaryStrategy) Priority() int {
	return s.priority
}

// CanHandle determines if this strategy can handle the given error
func (s *PrimaryStrategy) CanHandle(err error) bool {
	if err == nil {
		return true
	}
	
	errorStr := strings.ToLower(err.Error())
	// Primary strategy only handles specific, low-risk errors
	return strings.Contains(errorStr, "connection timeout") ||
		strings.Contains(errorStr, "temporary failure")
}

// Execute performs primary strategy validation
func (s *PrimaryStrategy) Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	if input.Error == nil {
		return &ValidationResult{
			Success:       true,
			Strategy:      s.Name(),
			Message:       "No error - validation passed",
			SecurityLevel: SecurityHigh,
			Actions:       []string{"Standard validation completed"},
		}, nil
	}
	
	errorStr := strings.ToLower(input.Error.Error())
	if strings.Contains(errorStr, "connection timeout") || strings.Contains(errorStr, "temporary failure") {
		return &ValidationResult{
			Success:       false,
			Strategy:      s.Name(),
			Message:       "Temporary network issue detected - recommend retry",
			SecurityLevel: SecurityHigh,
			Actions:       []string{"Network retry recommended"},
			NewConfig:     map[string]interface{}{"retry_delay": "30s", "max_retries": 3},
		}, nil
	}
	
	return &ValidationResult{
		Success:       false,
		Strategy:      s.Name(),
		Message:       "Error cannot be handled by primary strategy",
		SecurityLevel: SecurityHigh,
		Actions:       []string{"Escalating to secondary strategy"},
	}, nil
}

// SecondaryStrategy implements relaxed validation with moderate security risk
type SecondaryStrategy struct {
	name     string
	priority int
}

// NewSecondaryStrategy creates a new secondary (relaxed) validation strategy
func NewSecondaryStrategy() FallbackStrategy {
	return &SecondaryStrategy{
		name:     "secondary-relaxed",
		priority: 50,
	}
}

// Name returns the strategy name
func (s *SecondaryStrategy) Name() string {
	return s.name
}

// Priority returns the strategy priority
func (s *SecondaryStrategy) Priority() int {
	return s.priority
}

// CanHandle determines if this strategy can handle the given error
func (s *SecondaryStrategy) CanHandle(err error) bool {
	if err == nil {
		return true
	}
	
	errorStr := strings.ToLower(err.Error())
	// Secondary strategy handles more error types
	return strings.Contains(errorStr, "certificate has expired") ||
		strings.Contains(errorStr, "certificate name does not match") ||
		strings.Contains(errorStr, "connection timeout") ||
		strings.Contains(errorStr, "temporary failure")
}

// Execute performs secondary strategy validation
func (s *SecondaryStrategy) Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	if input.Error == nil {
		return &ValidationResult{
			Success:       true,
			Strategy:      s.Name(),
			Message:       "No error - validation passed",
			SecurityLevel: SecurityMedium,
			Actions:       []string{"Relaxed validation completed"},
		}, nil
	}
	
	errorStr := strings.ToLower(input.Error.Error())
	switch {
	case strings.Contains(errorStr, "certificate has expired"):
		return &ValidationResult{
			Success:       false,
			Strategy:      s.Name(),
			Message:       "Expired certificate detected - manual renewal required",
			SecurityLevel: SecurityMedium,
			Actions:       []string{"Contact registry administrator"},
			NewConfig:     map[string]interface{}{"require_renewal": true},
		}, nil
	case strings.Contains(errorStr, "certificate name does not match"):
		return &ValidationResult{
			Success:       false,
			Strategy:      s.Name(),
			Message:       "Hostname mismatch detected - verify registry URL",
			SecurityLevel: SecurityMedium,
			Actions:       []string{"Verify registry URL"},
			NewConfig:     map[string]interface{}{"verify_hostname": true},
		}, nil
	case strings.Contains(errorStr, "connection timeout") || strings.Contains(errorStr, "temporary failure"):
		return &ValidationResult{
			Success:       false,
			Strategy:      s.Name(),
			Message:       "Network issue - retry with backoff",
			SecurityLevel: SecurityMedium,
			Actions:       []string{"Exponential backoff retry"},
			NewConfig:     map[string]interface{}{"retry_backoff": "exponential"},
		}, nil
	}
	
	return &ValidationResult{
		Success:       false,
		Strategy:      s.Name(),
		Message:       "Error cannot be handled by secondary strategy",
		SecurityLevel: SecurityMedium,
		Actions:       []string{"Escalating to tertiary strategy"},
	}, nil
}

// TertiaryStrategy implements minimal validation with high security risk
type TertiaryStrategy struct {
	name     string
	priority int
}

// NewTertiaryStrategy creates a new tertiary (minimal) validation strategy
func NewTertiaryStrategy() FallbackStrategy {
	return &TertiaryStrategy{
		name:     "tertiary-minimal",
		priority: 10,
	}
}

// Name returns the strategy name
func (s *TertiaryStrategy) Name() string {
	return s.name
}

// Priority returns the strategy priority
func (s *TertiaryStrategy) Priority() int {
	return s.priority
}

// CanHandle determines if this strategy can handle the given error
func (s *TertiaryStrategy) CanHandle(err error) bool {
	// Tertiary strategy can handle any error as last resort
	return true
}

// Execute performs tertiary strategy validation
func (s *TertiaryStrategy) Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	if input.Error == nil {
		return &ValidationResult{
			Success:       true,
			Strategy:      s.Name(),
			Message:       "No error - validation passed",
			SecurityLevel: SecurityLow,
			Actions:       []string{"Minimal validation completed"},
		}, nil
	}
	
	errorStr := strings.ToLower(input.Error.Error())
	if strings.Contains(errorStr, "self signed certificate") || strings.Contains(errorStr, "certificate signed by unknown authority") {
		return &ValidationResult{
			Success:       true,
			Strategy:      s.Name(),
			Message:       "Self-signed certificate accepted with security warning",
			SecurityLevel: SecurityLow,
			Actions:       []string{"WARNING: Self-signed certificate trusted"},
			NewConfig:     map[string]interface{}{"trust_self_signed": true, "audit_required": true},
		}, nil
	}
	
	if strings.Contains(errorStr, "certificate") {
		return &ValidationResult{
			Success:       true,
			Strategy:      s.Name(),
			Message:       "Certificate error bypassed with minimal validation",
			SecurityLevel: SecurityNone,
			Actions:       []string{"WARNING: Certificate validation bypassed"},
			NewConfig:     map[string]interface{}{"bypass_validation": true, "high_risk": true},
		}, nil
	}
	
	return &ValidationResult{
		Success:       true,
		Strategy:      s.Name(),
		Message:       "Network error bypassed - operation allowed",
		SecurityLevel: SecurityLow,
		Actions:       []string{"Network validation bypassed"},
		NewConfig:     map[string]interface{}{"bypass_network": true},
	}, nil
}

// ErrorAnalyzer provides utilities for analyzing certificate errors
type ErrorAnalyzer struct{}

// NewErrorAnalyzer creates a new error analyzer
func NewErrorAnalyzer() *ErrorAnalyzer {
	return &ErrorAnalyzer{}
}

// AnalyzeError categorizes an error and returns analysis results
func (ea *ErrorAnalyzer) AnalyzeError(err error) *ErrorAnalysis {
	if err == nil {
		return &ErrorAnalysis{Category: "none", Severity: "info", Recoverable: true, Actions: []string{"No action required"}}
	}
	
	errorStr := strings.ToLower(err.Error())
	switch {
	case strings.Contains(errorStr, "self signed certificate"):
		return &ErrorAnalysis{Category: "self-signed", Severity: "high", Recoverable: true, 
			Actions: []string{"Add certificate to trust store", "Use --insecure flag"}}
	case strings.Contains(errorStr, "certificate has expired"):
		return &ErrorAnalysis{Category: "expired", Severity: "critical", Recoverable: false,
			Actions: []string{"Contact administrator for certificate renewal"}}
	case strings.Contains(errorStr, "certificate name does not match"):
		return &ErrorAnalysis{Category: "hostname-mismatch", Severity: "medium", Recoverable: true,
			Actions: []string{"Verify registry URL", "Check DNS configuration"}}
	case strings.Contains(errorStr, "no such host"):
		return &ErrorAnalysis{Category: "network", Severity: "medium", Recoverable: true,
			Actions: []string{"Check network connectivity", "Verify hostname"}}
	default:
		return &ErrorAnalysis{Category: "unknown", Severity: "medium", Recoverable: true,
			Actions: []string{"Check logs", "Retry operation"}}
	}
}

// ErrorAnalysis contains the results of error analysis
type ErrorAnalysis struct {
	Category    string
	Severity    string
	Recoverable bool
	Actions     []string
}

// AggregateErrors combines multiple errors into a single error report
func (ea *ErrorAnalyzer) AggregateErrors(errors []error) *ErrorAggregation {
	if len(errors) == 0 {
		return &ErrorAggregation{
			Count:      0,
			Categories: make(map[string]int),
			Actions:    []string{"No errors to aggregate"},
		}
	}
	
	aggregation := &ErrorAggregation{
		Count:      len(errors),
		Categories: make(map[string]int),
		Actions:    make([]string, 0),
	}
	
	actionSet := make(map[string]bool)
	
	for _, err := range errors {
		analysis := ea.AnalyzeError(err)
		aggregation.Categories[analysis.Category]++
		
		// Collect unique actions
		for _, action := range analysis.Actions {
			if !actionSet[action] {
				actionSet[action] = true
				aggregation.Actions = append(aggregation.Actions, action)
			}
		}
	}
	
	return aggregation
}

// ErrorAggregation contains aggregated error information
type ErrorAggregation struct {
	Count      int
	Categories map[string]int
	Actions    []string
}