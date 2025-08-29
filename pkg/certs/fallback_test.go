package certs

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestFallbackChain(t *testing.T) {
	chain := NewFallbackChain()
	
	// Verify chain is not nil
	if chain == nil {
		t.Fatal("NewFallbackChain returned nil")
	}
	
	// Verify default strategies are added
	strategies := chain.GetStrategies()
	if len(strategies) != 3 {
		t.Errorf("Expected 3 default strategies, got %d", len(strategies))
	}
	
	// Verify strategies are in priority order
	expectedNames := []string{"primary-strict", "secondary-relaxed", "tertiary-minimal"}
	for i, strategy := range strategies {
		if strategy.Name() != expectedNames[i] {
			t.Errorf("Strategy %d: expected %s, got %s", i, expectedNames[i], strategy.Name())
		}
	}
}

func TestFallbackChainExecution(t *testing.T) {
	chain := NewFallbackChain()
	ctx := context.Background()
	
	// Test successful validation (no error)
	input := &ValidationInput{
		Registry:  "test.registry.com",
		Operation: "pull",
		Error:     nil,
	}
	
	result, err := chain.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	if !result.Success {
		t.Error("Expected successful validation for no error")
	}
	
	if result.Strategy != "primary-strict" {
		t.Errorf("Expected primary-strict strategy, got %s", result.Strategy)
	}
}

func TestFallbackChainWithErrors(t *testing.T) {
	chain := NewFallbackChain()
	ctx := context.Background()
	
	testCases := []struct {
		name            string
		error           error
		expectedSuccess bool
		expectedStrategy string
	}{
		{
			name:            "Connection timeout - Tertiary handles",
			error:           errors.New("connection timeout"),
			expectedSuccess: true, // Tertiary accepts all errors
			expectedStrategy: "tertiary-minimal",
		},
		{
			name:            "Expired certificate - Tertiary handles", 
			error:           errors.New("certificate has expired"),
			expectedSuccess: true, // Tertiary accepts all errors
			expectedStrategy: "tertiary-minimal",
		},
		{
			name:            "Self-signed certificate - Tertiary handles",
			error:           errors.New("self signed certificate"),
			expectedSuccess: true, // Tertiary accepts with warning
			expectedStrategy: "tertiary-minimal",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &ValidationInput{
				Registry:  "test.registry.com",
				Operation: "pull",
				Error:     tc.error,
			}
			
			result, err := chain.Execute(ctx, input)
			if err != nil && tc.expectedSuccess {
				t.Fatalf("Execute failed unexpectedly: %v", err)
			}
			
			if result.Success != tc.expectedSuccess {
				t.Errorf("Expected success=%t, got %t", tc.expectedSuccess, result.Success)
			}
			
			if result.Strategy != tc.expectedStrategy {
				t.Errorf("Expected strategy %s, got %s", tc.expectedStrategy, result.Strategy)
			}
		})
	}
}

func TestPrimaryStrategy(t *testing.T) {
	strategy := NewPrimaryStrategy()
	ctx := context.Background()
	
	// Test strategy properties
	if strategy.Name() != "primary-strict" {
		t.Errorf("Expected name 'primary-strict', got %s", strategy.Name())
	}
	
	if strategy.Priority() != 100 {
		t.Errorf("Expected priority 100, got %d", strategy.Priority())
	}
	
	// Test error handling capabilities
	testCases := []struct {
		error     error
		canHandle bool
	}{
		{nil, true},
		{errors.New("connection timeout"), true},
		{errors.New("temporary failure"), true},
		{errors.New("self signed certificate"), false},
		{errors.New("certificate has expired"), false},
	}
	
	for _, tc := range testCases {
		if strategy.CanHandle(tc.error) != tc.canHandle {
			t.Errorf("CanHandle(%v) = %t, expected %t", tc.error, strategy.CanHandle(tc.error), tc.canHandle)
		}
	}
	
	// Test execution
	input := &ValidationInput{
		Registry:  "test.com",
		Operation: "pull",
		Error:     nil,
	}
	
	result, err := strategy.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	if !result.Success {
		t.Error("Expected successful execution for no error")
	}
	
	if result.SecurityLevel != SecurityHigh {
		t.Errorf("Expected SecurityHigh, got %v", result.SecurityLevel)
	}
}

func TestSecondaryStrategy(t *testing.T) {
	strategy := NewSecondaryStrategy()
	ctx := context.Background()
	
	// Test strategy properties
	if strategy.Name() != "secondary-relaxed" {
		t.Errorf("Expected name 'secondary-relaxed', got %s", strategy.Name())
	}
	
	if strategy.Priority() != 50 {
		t.Errorf("Expected priority 50, got %d", strategy.Priority())
	}
	
	// Test expired certificate handling
	input := &ValidationInput{
		Registry:  "test.com",
		Operation: "pull",
		Error:     errors.New("certificate has expired"),
	}
	
	result, err := strategy.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	if result.Success {
		t.Error("Expected failure for expired certificate")
	}
	
	if result.SecurityLevel != SecurityMedium {
		t.Errorf("Expected SecurityMedium, got %v", result.SecurityLevel)
	}
	
	if !strings.Contains(result.Message, "Expired certificate") {
		t.Errorf("Expected message about expired certificate, got: %s", result.Message)
	}
}

func TestTertiaryStrategy(t *testing.T) {
	strategy := NewTertiaryStrategy()
	ctx := context.Background()
	
	// Test strategy properties
	if strategy.Name() != "tertiary-minimal" {
		t.Errorf("Expected name 'tertiary-minimal', got %s", strategy.Name())
	}
	
	if strategy.Priority() != 10 {
		t.Errorf("Expected priority 10, got %d", strategy.Priority())
	}
	
	// Tertiary strategy should handle any error
	testErrors := []error{
		errors.New("self signed certificate"),
		errors.New("certificate has expired"),
		errors.New("no such host"),
		errors.New("unknown error"),
	}
	
	for _, err := range testErrors {
		if !strategy.CanHandle(err) {
			t.Errorf("Tertiary strategy should handle all errors, failed for: %v", err)
		}
	}
	
	// Test self-signed certificate handling
	input := &ValidationInput{
		Registry:  "test.com",
		Operation: "pull",
		Error:     errors.New("self signed certificate"),
	}
	
	result, err := strategy.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	if !result.Success {
		t.Error("Expected success for self-signed certificate in tertiary strategy")
	}
	
	if result.SecurityLevel != SecurityLow {
		t.Errorf("Expected SecurityLow, got %v", result.SecurityLevel)
	}
}

func TestFallbackChainAddRemoveStrategy(t *testing.T) {
	chain := NewFallbackChain()
	
	// Create a custom strategy
	customStrategy := &customTestStrategy{
		name:     "custom-test",
		priority: 75,
	}
	
	// Add custom strategy
	chain.AddStrategy(customStrategy)
	
	strategies := chain.GetStrategies()
	if len(strategies) != 4 {
		t.Errorf("Expected 4 strategies after adding custom, got %d", len(strategies))
	}
	
	// Verify priority ordering
	if strategies[1].Name() != "custom-test" {
		t.Errorf("Custom strategy should be second (priority 75), but got: %s", strategies[1].Name())
	}
	
	// Remove custom strategy
	chain.RemoveStrategy("custom-test")
	
	strategies = chain.GetStrategies()
	if len(strategies) != 3 {
		t.Errorf("Expected 3 strategies after removing custom, got %d", len(strategies))
	}
	
	// Verify custom strategy is gone
	for _, strategy := range strategies {
		if strategy.Name() == "custom-test" {
			t.Error("Custom strategy should have been removed")
		}
	}
}

func TestErrorAnalyzer(t *testing.T) {
	analyzer := NewErrorAnalyzer()
	
	testCases := []struct {
		error             error
		expectedCategory  string
		expectedSeverity  string
		expectedRecoverable bool
	}{
		{
			error:             nil,
			expectedCategory:  "none",
			expectedSeverity:  "info",
			expectedRecoverable: true,
		},
		{
			error:             errors.New("self signed certificate"),
			expectedCategory:  "self-signed",
			expectedSeverity:  "high",
			expectedRecoverable: true,
		},
		{
			error:             errors.New("certificate has expired"),
			expectedCategory:  "expired",
			expectedSeverity:  "critical",
			expectedRecoverable: false,
		},
		{
			error:             errors.New("certificate name does not match"),
			expectedCategory:  "hostname-mismatch",
			expectedSeverity:  "medium",
			expectedRecoverable: true,
		},
		{
			error:             errors.New("no such host"),
			expectedCategory:  "network",
			expectedSeverity:  "medium",
			expectedRecoverable: true,
		},
		{
			error:             errors.New("unknown error type"),
			expectedCategory:  "unknown",
			expectedSeverity:  "medium",
			expectedRecoverable: true,
		},
	}
	
	for _, tc := range testCases {
		analysis := analyzer.AnalyzeError(tc.error)
		
		if analysis.Category != tc.expectedCategory {
			t.Errorf("Error %v: expected category %s, got %s", tc.error, tc.expectedCategory, analysis.Category)
		}
		
		if analysis.Severity != tc.expectedSeverity {
			t.Errorf("Error %v: expected severity %s, got %s", tc.error, tc.expectedSeverity, analysis.Severity)
		}
		
		if analysis.Recoverable != tc.expectedRecoverable {
			t.Errorf("Error %v: expected recoverable %t, got %t", tc.error, tc.expectedRecoverable, analysis.Recoverable)
		}
		
		if len(analysis.Actions) == 0 {
			t.Errorf("Error %v: expected non-empty actions list", tc.error)
		}
	}
}

func TestErrorAggregation(t *testing.T) {
	analyzer := NewErrorAnalyzer()
	
	errors := []error{
		errors.New("self signed certificate"),
		errors.New("certificate has expired"),
		errors.New("self signed certificate"),
		errors.New("no such host"),
	}
	
	aggregation := analyzer.AggregateErrors(errors)
	
	if aggregation.Count != 4 {
		t.Errorf("Expected count 4, got %d", aggregation.Count)
	}
	
	if aggregation.Categories["self-signed"] != 2 {
		t.Errorf("Expected 2 self-signed errors, got %d", aggregation.Categories["self-signed"])
	}
	
	if aggregation.Categories["expired"] != 1 {
		t.Errorf("Expected 1 expired error, got %d", aggregation.Categories["expired"])
	}
	
	if aggregation.Categories["network"] != 1 {
		t.Errorf("Expected 1 network error, got %d", aggregation.Categories["network"])
	}
	
	if len(aggregation.Actions) == 0 {
		t.Error("Expected non-empty actions list")
	}
	
	// Test empty error list
	emptyAggregation := analyzer.AggregateErrors([]error{})
	if emptyAggregation.Count != 0 {
		t.Errorf("Expected count 0 for empty list, got %d", emptyAggregation.Count)
	}
}

func TestSecurityLevels(t *testing.T) {
	// Test that security levels have expected values
	if int(SecurityHigh) != 0 {
		t.Errorf("SecurityHigh should be 0, got %d", int(SecurityHigh))
	}
	
	if int(SecurityMedium) != 1 {
		t.Errorf("SecurityMedium should be 1, got %d", int(SecurityMedium))
	}
	
	if int(SecurityLow) != 2 {
		t.Errorf("SecurityLow should be 2, got %d", int(SecurityLow))
	}
	
	if int(SecurityNone) != 3 {
		t.Errorf("SecurityNone should be 3, got %d", int(SecurityNone))
	}
}

// Custom test strategy for testing add/remove functionality
type customTestStrategy struct {
	name     string
	priority int
}

func (s *customTestStrategy) Name() string {
	return s.name
}

func (s *customTestStrategy) Priority() int {
	return s.priority
}

func (s *customTestStrategy) CanHandle(err error) bool {
	return err != nil && strings.Contains(err.Error(), "custom")
}

func (s *customTestStrategy) Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error) {
	return &ValidationResult{
		Success:  true,
		Strategy: s.Name(),
		Message:  "Custom strategy executed",
		SecurityLevel: SecurityMedium,
	}, nil
}

// Integration tests for insecure mode functionality

func TestNewFallbackChainWithInsecure(t *testing.T) {
	config := &InsecureConfig{
		Enabled:        true,
		RequireConsent: false,
		LogWarnings:    false, // Disable for testing
	}
	
	chain := NewFallbackChainWithInsecure(config)
	if chain == nil {
		t.Fatal("NewFallbackChainWithInsecure returned nil")
	}
	
	strategies := chain.GetStrategies()
	
	// Should have 4 strategies (3 default + 1 insecure)
	if len(strategies) != 4 {
		t.Errorf("Expected 4 strategies, got %d", len(strategies))
	}
	
	// Verify strategies are in correct priority order
	expectedNames := []string{"primary-strict", "secondary-relaxed", "insecure-bypass", "tertiary-minimal"}
	for i, strategy := range strategies {
		if strategy.Name() != expectedNames[i] {
			t.Errorf("Strategy %d: expected %s, got %s", i, expectedNames[i], strategy.Name())
		}
	}
	
	// Verify insecure strategy has priority 15
	insecureStrategy := strategies[2]
	if insecureStrategy.Priority() != 15 {
		t.Errorf("Expected insecure strategy priority 15, got %d", insecureStrategy.Priority())
	}
}

func TestFallbackChainWithInsecureExecution(t *testing.T) {
	config := &InsecureConfig{
		Enabled:        true,
		RequireConsent: false,
		LogWarnings:    false, // Disable for testing
	}
	
	chain := NewFallbackChainWithInsecure(config)
	ctx := context.Background()
	
	testCases := []struct {
		name            string
		error           error
		expectedSuccess bool
		expectedStrategy string
	}{
		{
			name:            "No error - handled by primary",
			error:           nil,
			expectedSuccess: true,
			expectedStrategy: "primary-strict",
		},
		{
			name:            "Self-signed cert - handled by insecure mode",
			error:           errors.New("self signed certificate"),
			expectedSuccess: true,
			expectedStrategy: "insecure-bypass",
		},
		{
			name:            "Unknown authority - handled by insecure mode",
			error:           errors.New("certificate signed by unknown authority"),
			expectedSuccess: true,
			expectedStrategy: "insecure-bypass",
		},
		{
			name:            "Expired cert - handled by insecure mode", 
			error:           errors.New("certificate has expired"),
			expectedSuccess: true,
			expectedStrategy: "insecure-bypass",
		},
		{
			name:            "Network error - handled by tertiary",
			error:           errors.New("network connection refused"),
			expectedSuccess: true,
			expectedStrategy: "tertiary-minimal",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := &ValidationInput{
				Registry:  "test.registry.com",
				Operation: "pull",
				Error:     tc.error,
			}
			
			result, err := chain.Execute(ctx, input)
			if err != nil && tc.expectedSuccess {
				t.Fatalf("Execute failed unexpectedly: %v", err)
			}
			
			if result.Success != tc.expectedSuccess {
				t.Errorf("Expected success=%t, got %t", tc.expectedSuccess, result.Success)
			}
			
			if result.Strategy != tc.expectedStrategy {
				t.Errorf("Expected strategy %s, got %s", tc.expectedStrategy, result.Strategy)
			}
		})
	}
}

func TestInsecureChainWithConsentRequired(t *testing.T) {
	config := &InsecureConfig{
		Enabled:        true,
		RequireConsent: true, // Consent required
		LogWarnings:    false,
	}
	
	chain := NewFallbackChainWithInsecure(config)
	ctx := context.Background()
	
	input := &ValidationInput{
		Registry:  "test.registry.com",
		Operation: "pull",
		Error:     errors.New("self signed certificate"),
	}
	
	result, err := chain.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	// Should fail because consent is required but not given
	if result.Success {
		t.Error("Expected failure when consent is required but not given")
	}
	
	if !strings.Contains(result.Message, "consent") {
		t.Error("Expected message to mention consent requirement")
	}
}

func TestInsecureChainDisabled(t *testing.T) {
	config := &InsecureConfig{
		Enabled: false, // Disabled
	}
	
	chain := NewFallbackChainWithInsecure(config)
	ctx := context.Background()
	
	input := &ValidationInput{
		Registry:  "test.registry.com", 
		Operation: "pull",
		Error:     errors.New("self signed certificate"),
	}
	
	result, err := chain.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	// Should fall back to tertiary strategy since insecure is disabled
	if result.Strategy != "tertiary-minimal" {
		t.Errorf("Expected tertiary-minimal strategy when insecure disabled, got %s", result.Strategy)
	}
	
	if result.SecurityLevel != SecurityLow {
		t.Errorf("Expected SecurityLow from tertiary strategy, got %v", result.SecurityLevel)
	}
}

func TestInsecureModeSecurityLevels(t *testing.T) {
	config := &InsecureConfig{
		Enabled:        true,
		RequireConsent: false,
		LogWarnings:    false,
	}
	
	chain := NewFallbackChainWithInsecure(config)
	ctx := context.Background()
	
	// Test that insecure mode returns SecurityNone
	input := &ValidationInput{
		Registry:  "test.registry.com",
		Operation: "pull", 
		Error:     errors.New("self signed certificate"),
	}
	
	result, err := chain.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	if result.SecurityLevel != SecurityNone {
		t.Errorf("Expected SecurityNone from insecure mode, got %v", result.SecurityLevel)
	}
	
	// Verify result contains security warnings
	found := false
	for _, action := range result.Actions {
		if strings.Contains(action, "SECURITY RISK") {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("Expected security warning in result actions")
	}
}

func TestInsecureIntegrationWithExistingStrategies(t *testing.T) {
	// Test that insecure mode doesn't interfere with existing strategies
	config := &InsecureConfig{
		Enabled:        true,
		RequireConsent: false,
		LogWarnings:    false,
	}
	
	chain := NewFallbackChainWithInsecure(config)
	ctx := context.Background()
	
	// Test that primary strategy still handles its errors first
	input := &ValidationInput{
		Registry:  "test.registry.com",
		Operation: "pull",
		Error:     errors.New("connection timeout"), // Primary handles this
	}
	
	result, err := chain.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	// Should use tertiary (which also handles timeouts) since primary doesn't succeed
	if result.Strategy != "tertiary-minimal" {
		t.Errorf("Expected tertiary-minimal strategy for timeout, got %s", result.Strategy)
	}
	
	// Test that secondary strategy handles its errors
	input = &ValidationInput{
		Registry:  "test.registry.com",
		Operation: "pull",
		Error:     errors.New("network connection refused"), // No strategy specifically handles this well
	}
	
	result, err = chain.Execute(ctx, input)
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	
	// Should fall to tertiary which handles all errors
	if result.Strategy != "tertiary-minimal" {
		t.Errorf("Expected tertiary-minimal strategy for network error, got %s", result.Strategy)
	}
}