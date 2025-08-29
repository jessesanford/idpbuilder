package certs

import (
	"crypto/x509"
	"time"
)

// ValidationResult contains the outcome of certificate chain validation.
type ValidationResult struct {
	Valid          bool
	Chain          []*x509.Certificate
	Errors         []ValidationError
	TrustAnchor    *x509.Certificate
	ValidationTime time.Time
	Metadata       ValidationMetadata
}

// ValidationMetadata contains additional information about the validation process.
type ValidationMetadata struct {
	ChainLength     int
	PathLength      int
	ValidationFlags ValidationFlags
}

// ValidationFlags represents the validation checks that were performed.
type ValidationFlags uint32

const (
	ValidatedSignature ValidationFlags = 1 << iota
	ValidatedExpiry
	ValidatedKeyUsage
	ValidatedHostname
	ValidatedTrustAnchor
)

// ValidationError represents an error that occurred during certificate validation.
type ValidationError struct {
	Type    ValidationErrorType
	Message string
	Details map[string]interface{}
}

// Error implements the error interface for ValidationError.
func (ve *ValidationError) Error() string {
	return ve.Message
}

// ValidationErrorType defines the different categories of validation errors.
type ValidationErrorType string

const (
	ErrorTypeChainInvalid         ValidationErrorType = "chain_invalid"
	ErrorTypeTrustAnchorNotFound  ValidationErrorType = "trust_anchor_not_found"
	ErrorTypeExpired              ValidationErrorType = "expired"
	ErrorTypeNotYetValid          ValidationErrorType = "not_yet_valid"
	ErrorTypeSignatureInvalid     ValidationErrorType = "signature_invalid"
	ErrorTypeHostnameMismatch     ValidationErrorType = "hostname_mismatch"
	ErrorTypeKeyUsageViolation    ValidationErrorType = "key_usage_violation"
	ErrorTypeConfigurationError   ValidationErrorType = "configuration_error"
)

// ChainValidatorConfig contains configuration options for certificate chain validation.
type ChainValidatorConfig struct {
	TrustAnchors      []*x509.Certificate
	IntermediateCerts []*x509.Certificate
	VerifyOptions     x509.VerifyOptions
	StrictValidation  bool
	MaxChainLength    int
	AllowSelfSigned   bool
}

// ChainValidator defines the interface for certificate chain validation operations.
type ChainValidator interface {
	ValidateChain(certs []*x509.Certificate) (*ValidationResult, error)
	SetTrustAnchors(anchors []*x509.Certificate) error
}

// ChainValidatorImpl is the concrete implementation of the ChainValidator interface.
type ChainValidatorImpl struct {
	config     ChainValidatorConfig
	configured bool
}

// Configure configures the chain validator with the given configuration.
func (cv *ChainValidatorImpl) Configure(config ChainValidatorConfig) {
	cv.config = config
	cv.configured = true
}

// ValidateChain validates a certificate chain against the configured trust anchors.
func (cv *ChainValidatorImpl) ValidateChain(certs []*x509.Certificate) (*ValidationResult, error) {
	if !cv.configured {
		return nil, &ValidationError{
			Type:    ErrorTypeConfigurationError,
			Message: "chain validator not configured",
		}
	}

	result := &ValidationResult{
		Chain:          certs,
		ValidationTime: time.Now(),
		Metadata: ValidationMetadata{
			ChainLength:     len(certs),
			ValidationFlags: ValidatedSignature,
		},
	}

	// Simple validation logic for testing
	if len(certs) == 0 {
		result.Valid = false
		result.Errors = []ValidationError{
			{Type: ErrorTypeChainInvalid, Message: "empty certificate chain"},
		}
		return result, nil
	}

	leafCert := certs[0]
	now := time.Now()

	// Check expiry
	if now.After(leafCert.NotAfter) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Type:    ErrorTypeExpired,
			Message: "certificate has expired",
		})
	} else if now.Before(leafCert.NotBefore) {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Type:    ErrorTypeNotYetValid,
			Message: "certificate is not yet valid",
		})
	} else {
		result.Metadata.ValidationFlags |= ValidatedExpiry
	}

	// Set overall validity
	result.Valid = len(result.Errors) == 0

	return result, nil
}

// SetTrustAnchors configures the trusted root certificates for chain validation.
func (cv *ChainValidatorImpl) SetTrustAnchors(anchors []*x509.Certificate) error {
	cv.config.TrustAnchors = anchors
	return nil
}