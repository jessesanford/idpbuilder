package certs

import (
	"context"
	"crypto/x509"
	"time"
)

// Wave 1 interfaces and types for integration
// These would normally be imported from Wave 1, but for isolated development we define them here

// CertValidator defines the interface for basic certificate validation (from Wave 1)
type CertValidator interface {
	// ValidateCertificate performs basic validation checks on a certificate
	ValidateCertificate(cert *x509.Certificate) (*ValidationResult, error)
	
	// CheckExpiry checks if a certificate is expired or expiring soon
	CheckExpiry(cert *x509.Certificate, warnDays int) (*ExpiryResult, error)
}

// TrustManager defines the interface for managing trust store certificates (from Wave 1)
type TrustManager interface {
	// ValidateCertificate validates a certificate against the trust store
	ValidateCertificate(ctx context.Context, registry string, cert *x509.Certificate) error
}

// ValidationResult contains the results of basic certificate validation (from Wave 1)
type ValidationResult struct {
	// Valid indicates if the certificate passed all validation checks
	Valid bool
	
	// Issues contains any validation issues found
	Issues []string
	
	// Errors contains any validation errors found
	Errors []string
	
	// Warnings contains any validation warnings
	Warnings []string
	
	// Subject is the certificate subject
	Subject string
	
	// Issuer is the certificate issuer
	Issuer string
	
	// NotBefore is the certificate validity start time
	NotBefore time.Time
	
	// NotAfter is the certificate validity end time
	NotAfter time.Time
	
	// IsCA indicates if this is a CA certificate
	IsCA bool
	
	// IsSelfSigned indicates if the certificate is self-signed
	IsSelfSigned bool
	
	// DNSNames contains the DNS names in the certificate
	DNSNames []string
	
	// IPAddresses contains the IP addresses in the certificate
	IPAddresses []string
}

// ExpiryResult contains certificate expiry information (from Wave 1)
type ExpiryResult struct {
	// Expired indicates if the certificate has already expired
	Expired bool
	
	// ExpiringSoon indicates if the certificate will expire within the warning period
	ExpiringSoon bool
	
	// DaysUntilExpiry is the number of days until the certificate expires
	DaysUntilExpiry int
	
	// ExpiryDate is the certificate expiry date
	ExpiryDate time.Time
}