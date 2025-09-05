// Package certs provides common types for certificate handling and validation
// This file consolidates types that were previously duplicated across multiple files
package certs

import (
	"context"
	"crypto/x509"
	"net"
	"time"
)

// ValidationInput contains the input data for certificate validation operations
type ValidationInput struct {
	Certificates []*x509.Certificate
	Registry     string
	Operation    string
	Error        error
	Options      map[string]interface{}
}

// ValidationResult contains the result of a certificate validation attempt
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

// RecoveryConfig contains configuration for certificate recovery operations
// This consolidates fields from both fallback.go and recovery.go versions
type RecoveryConfig struct {
	// From recovery.go
	EnableCertRefresh   bool
	EnableTrustUpdate   bool
	EnableChainRepair   bool
	MaxAttempts         int
	Timeout             time.Duration
	CircuitBreakerThreshold int
	CircuitBreakerTimeout   time.Duration
	
	// From fallback.go (additional fields)  
	AllowInsecure      bool
	FallbackRegistries []string
	RecoveryStrategies []string
	TrustStoreUpdate   bool
}

// RecoveryResult contains the result of a recovery operation
// This consolidates fields from both fallback.go and recovery.go versions
type RecoveryResult struct {
	Success          bool
	Method           string    // From recovery.go (was 'Method')
	Strategy         string    // From fallback.go (additional)
	Actions          []string
	NewConfig        interface{}
	Message          string    // From recovery.go
	FailureReason    string    // From fallback.go (additional)
	RecoveredCerts   []*x509.Certificate // From fallback.go (additional)
	SecurityDowngrade bool     // From fallback.go (additional)
	NextRetryAfter   time.Duration // From fallback.go (additional)
}

// FallbackStrategy defines the interface for certificate validation fallback strategies
type FallbackStrategy interface {
	Name() string
	Priority() int
	CanHandle(err error) bool
	Execute(ctx context.Context, input *ValidationInput) (*ValidationResult, error)
}

// TrustStoreManager defines the interface for managing trusted certificates per registry
type TrustStoreManager interface {
	// Certificate management
	AddCertificate(registry string, cert *x509.Certificate) error
	RemoveCertificate(registry string) error
	GetTrustedCerts(registry string) ([]*x509.Certificate, error)
	GetCertPool(registry string) (*x509.CertPool, error)
	
	// Security settings
	SetInsecureRegistry(registry string, insecure bool) error
	IsInsecure(registry string) bool
	
	// Persistence
	LoadFromDisk() error
	SaveToDisk(registry string, cert *x509.Certificate) error
	
	// Transport configuration
	ConfigureTransport(registry string) (interface{}, error)
	ConfigureTransportWithConfig(registry string, config interface{}) (interface{}, error)
	CreateHTTPClient(registry string) (interface{}, error)
	CreateHTTPClientWithConfig(registry string, config interface{}) (interface{}, error)
}

// CertDiagnostics contains comprehensive certificate diagnostic information
type CertDiagnostics struct {
	Subject          string
	Issuer           string
	SerialNumber     string
	NotBefore        time.Time
	NotAfter         time.Time
	DNSNames         []string
	IPAddresses      []net.IP
	ValidationErrors []ValidationError
	Warnings         []string
}

// ValidationError represents a certificate validation error
type ValidationError struct {
	Type    string // e.g., "chain", "expiry", "hostname"
	Message string
	Detail  string
}