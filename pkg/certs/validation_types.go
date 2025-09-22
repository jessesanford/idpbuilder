// Copyright 2024 The IDP Builder Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certs

import (
	"crypto/x509"
	"errors"
	"time"
)

// Common errors for certificate operations
var (
	ErrFeatureDisabled = errors.New("feature is disabled")
)

// CertificateValidator defines the comprehensive interface for certificate validation operations
// This is the primary interface from the cert-validation effort
type CertificateValidator interface {
	// ValidateChain validates a complete certificate chain from leaf to root
	ValidateChain(certs []*x509.Certificate) error

	// ValidateCertificate validates a single certificate
	ValidateCertificate(cert *x509.Certificate) error

	// VerifyHostname verifies that a certificate is valid for a given hostname
	VerifyHostname(cert *x509.Certificate, hostname string) error

	// GenerateDiagnostics creates diagnostic information for troubleshooting
	GenerateDiagnostics() (*CertDiagnostics, error)

	// SetValidationMode changes the validation strictness
	SetValidationMode(mode ValidationMode)

	// GetValidationMode returns the current validation mode
	GetValidationMode() ValidationMode
}

// BasicValidator defines a simplified interface for basic certificate validation
// This is the renamed interface from the fallback-strategies effort
type BasicValidator interface {
	// Validate checks if the certificate is valid
	Validate(cert *Certificate) error

	// ValidateChain checks if the certificate chain is valid
	ValidateChain(chain []*Certificate) error

	// IsExpired checks if the certificate has expired
	IsExpired(cert *Certificate) bool

	// WillExpireSoon checks if certificate will expire within threshold
	WillExpireSoon(cert *Certificate, threshold time.Duration) bool
}

// ValidationResult represents the unified result of certificate validation
// This combines fields from both implementations for maximum compatibility
type ValidationResult struct {
	// Core fields
	Valid       bool      `json:"valid"`
	ValidatedAt time.Time `json:"validated_at"`

	// Detailed tracking (from validator.go implementation)
	Errors      []error             `json:"errors,omitempty"`
	Warnings    []string            `json:"warnings,omitempty"`
	Certificate *x509.Certificate   `json:"-"`
	Chain       []*x509.Certificate `json:"-"`

	// Simple message and actions (from utilities.go implementation)
	Message string   `json:"message,omitempty"`
	Actions []string `json:"actions,omitempty"`
}

// IsValid provides backward compatibility for code expecting the old field name
func (v *ValidationResult) IsValid() bool {
	return v.Valid
}

// NewValidationResult creates a new validation result with timestamp
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		ValidatedAt: time.Now(),
		Errors:      make([]error, 0),
		Warnings:    make([]string, 0),
		Actions:     make([]string, 0),
	}
}

// AddError adds an error to the validation result
func (v *ValidationResult) AddError(err error) {
	v.Errors = append(v.Errors, err)
	v.Valid = false
}

// AddWarning adds a warning to the validation result
func (v *ValidationResult) AddWarning(warning string) {
	v.Warnings = append(v.Warnings, warning)
}

// AddAction adds a suggested action to the validation result
func (v *ValidationResult) AddAction(action string) {
	v.Actions = append(v.Actions, action)
}
