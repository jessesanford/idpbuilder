package certs

import (
	"context"
	"crypto/x509"
)

// ChainValidator extends basic certificate validation with comprehensive chain verification
type ChainValidator interface {
	// ValidateChain verifies the complete certificate chain from leaf to root
	ValidateChain(ctx context.Context, cert *x509.Certificate, intermediates []*x509.Certificate) (*ChainValidationResult, error)
	
	// VerifyHostname checks if the certificate is valid for the given hostname
	VerifyHostname(cert *x509.Certificate, hostname string) error
	
	// CheckChainExpiry verifies no certificates in the chain are expired or expiring soon
	CheckChainExpiry(chain []*x509.Certificate, warnDays int) (*ChainExpiryResult, error)
	
	// GenerateDiagnostics creates a comprehensive diagnostic report
	GenerateDiagnostics(ctx context.Context, cert *x509.Certificate, hostname string) (*CertDiagnosticsReport, error)
}