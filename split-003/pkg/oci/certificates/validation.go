package certificates

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"
	"time"
)

// Validator implements the CertificateValidator interface.
type Validator struct {
	rules map[string]*ValidationRule
	mu    sync.RWMutex
}

// NewValidator creates a new certificate validator with default rules.
func NewValidator() *Validator {
	v := &Validator{
		rules: make(map[string]*ValidationRule),
	}

	// Add default validation rules
	defaultRules := map[string]bool{
		"expiry":     true,
		"not_before": true,
		"pem_format": true,
	}

	for name, critical := range defaultRules {
		rule := &ValidationRule{
			Name:        name,
			Description: fmt.Sprintf("Default %s validation", name),
			Enabled:     true,
			Critical:    critical,
			Parameters:  make(map[string]interface{}),
		}
		v.rules[name] = rule
	}

	return v
}

// ValidateCertificate validates a single certificate against all enabled rules.
func (v *Validator) ValidateCertificate(ctx context.Context, cert *Certificate) (*ValidationResult, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	v.mu.RLock()
	defer v.mu.RUnlock()

	result := &ValidationResult{
		Valid:         true,
		CertificateID: cert.ID,
		ValidatedAt:   time.Now(),
		Errors:        make([]ValidationError, 0),
		Warnings:      make([]ValidationWarning, 0),
	}

	// Apply validation rules
	errors := v.validateExpiry(cert)
	errors = append(errors, v.validateNotBefore(cert)...)
	errors = append(errors, v.validatePEMFormat(cert)...)

	for _, err := range errors {
		err.CertificateID = cert.ID
		if err.Code == ErrCodeCertExpired || err.Code == ErrCodeInvalidPEM {
			result.Errors = append(result.Errors, err)
			result.Valid = false
		} else {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Rule:    "validation",
				Message: err.Message,
				Code:    err.Code,
			})
		}
	}

	return result, nil
}

// ValidatePEM validates PEM-encoded certificate data.
func (v *Validator) ValidatePEM(ctx context.Context, pemData []byte) (*ValidationResult, error) {
	if len(pemData) == 0 {
		return &ValidationResult{
			Valid:       false,
			ValidatedAt: time.Now(),
			Errors: []ValidationError{
				{Code: ErrCodeInvalidPEM, Message: "PEM data cannot be empty"},
			},
		}, nil
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return &ValidationResult{
			Valid:       false,
			ValidatedAt: time.Now(),
			Errors: []ValidationError{
				{Code: ErrCodeInvalidPEM, Message: "invalid PEM format"},
			},
		}, nil
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return &ValidationResult{
			Valid:       false,
			ValidatedAt: time.Now(),
			Errors: []ValidationError{
				{Code: ErrCodeInvalidPEM, Message: fmt.Sprintf("failed to parse certificate: %v", err)},
			},
		}, nil
	}

	cert := &Certificate{
		ID:          "temp",
		Certificate: x509Cert,
		PEM:         pemData,
		ValidFrom:   x509Cert.NotBefore,
		ValidTo:     x509Cert.NotAfter,
		Issuer:      x509Cert.Issuer.String(),
		Subject:     x509Cert.Subject.String(),
	}

	return v.ValidateCertificate(ctx, cert)
}

// ValidateChain validates a chain of certificates.
func (v *Validator) ValidateChain(ctx context.Context, chain []*Certificate) (*ValidationResult, error) {
	if len(chain) == 0 {
		return nil, fmt.Errorf("certificate chain cannot be empty")
	}

	result := &ValidationResult{Valid: true, ValidatedAt: time.Now()}
	for _, cert := range chain {
		certResult, err := v.ValidateCertificate(ctx, cert)
		if err != nil || !certResult.Valid {
			result.Valid = false
		}
		if certResult != nil {
			result.Errors = append(result.Errors, certResult.Errors...)
			result.Warnings = append(result.Warnings, certResult.Warnings...)
		}
	}
	return result, nil
}

// AddValidationRule adds a custom validation rule.
func (v *Validator) AddValidationRule(ctx context.Context, rule *ValidationRule) error {
	if rule == nil || rule.Name == "" {
		return fmt.Errorf("invalid validation rule")
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if _, exists := v.rules[rule.Name]; exists {
		return NewValidationError(ErrCodeValidationRuleExists, fmt.Sprintf("rule '%s' already exists", rule.Name))
	}
	v.rules[rule.Name] = rule
	return nil
}

// RemoveValidationRule removes a validation rule.
func (v *Validator) RemoveValidationRule(ctx context.Context, ruleName string) error {
	if ruleName == "" {
		return fmt.Errorf("rule name cannot be empty")
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	builtinRules := map[string]bool{"expiry": true, "not_before": true, "pem_format": true}
	if builtinRules[ruleName] {
		return fmt.Errorf("cannot remove builtin rule '%s'", ruleName)
	}
	delete(v.rules, ruleName)
	return nil
}

// ListValidationRules returns all validation rules.
func (v *Validator) ListValidationRules(ctx context.Context) ([]*ValidationRule, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	rules := make([]*ValidationRule, 0, len(v.rules))
	for _, rule := range v.rules {
		rules = append(rules, rule)
	}
	return rules, nil
}

// Built-in validation functions

func (v *Validator) validateExpiry(cert *Certificate) []ValidationError {
	var errors []ValidationError
	if cert.Certificate == nil {
		return errors
	}

	now := time.Now()
	if now.After(cert.Certificate.NotAfter) {
		errors = append(errors, ValidationError{
			Code:    ErrCodeCertExpired,
			Message: fmt.Sprintf("certificate expired on %v", cert.Certificate.NotAfter),
		})
	}
	return errors
}

func (v *Validator) validateNotBefore(cert *Certificate) []ValidationError {
	var errors []ValidationError
	if cert.Certificate == nil {
		return errors
	}

	now := time.Now()
	if now.Before(cert.Certificate.NotBefore) {
		errors = append(errors, ValidationError{
			Code:    "CERT_NOT_YET_VALID",
			Message: fmt.Sprintf("certificate not valid until %v", cert.Certificate.NotBefore),
		})
	}
	return errors
}

func (v *Validator) validatePEMFormat(cert *Certificate) []ValidationError {
	var errors []ValidationError
	if len(cert.PEM) == 0 {
		errors = append(errors, ValidationError{
			Code:    ErrCodeInvalidPEM,
			Message: "certificate PEM data is empty",
		})
	}
	return errors
}