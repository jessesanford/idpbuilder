// Package certs provides basic types for certificate fallback strategies
package certs

import (
	"context"
	"crypto/x509"
)

// FallbackStrategy defines the interface for certificate validation fallback strategies
type FallbackStrategy interface {
	Name() string
	Priority() int
	CanHandle(err error) bool
	Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
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