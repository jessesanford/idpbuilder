// Package certs provides interfaces for certificate validation (minimal version for fallback strategies)
package certs

import (
	"crypto/x509"
	"net"
	"time"
)

// CertValidator provides comprehensive X.509 certificate validation
type CertValidator interface {
	// ValidateChain verifies the certificate chain against trusted roots
	ValidateChain(cert *x509.Certificate) error
	
	// CheckExpiry checks if certificate is expired or expiring soon
	CheckExpiry(cert *x509.Certificate) (*time.Duration, error)
	
	// VerifyHostname checks if the certificate is valid for the given hostname
	VerifyHostname(cert *x509.Certificate, hostname string) error
	
	// GenerateDiagnostics creates a detailed diagnostic report for the certificate
	GenerateDiagnostics(cert *x509.Certificate) (*CertDiagnostics, error)
}

// CertDiagnostics contains detailed information about certificate validation
type CertDiagnostics struct {
	Subject         string
	Issuer          string
	SerialNumber    string
	NotBefore       time.Time
	NotAfter        time.Time
	DNSNames        []string
	IPAddresses     []net.IP
	ValidationErrors []ValidationError
	Warnings        []string
}

// ValidationError represents a specific validation failure
type ValidationError struct {
	Type    string // "chain", "expiry", "hostname", etc.
	Message string
	Detail  string
}