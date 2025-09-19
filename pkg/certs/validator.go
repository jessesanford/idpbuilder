package certs

import (
	"crypto/x509"
	"errors"
)

// Validator interface defines methods for certificate validation
type Validator interface {
	Validate(cert *x509.Certificate) error
	ValidateChain(certs []*x509.Certificate) error
}

// BasicValidator provides basic certificate validation
type BasicValidator struct{}

// Validate implements basic certificate validation
func (v *BasicValidator) Validate(cert *x509.Certificate) error {
	if cert == nil {
		return errors.New("certificate is nil")
	}
	return nil
}

// ValidateChain implements basic chain validation
func (v *BasicValidator) ValidateChain(certs []*x509.Certificate) error {
	if len(certs) == 0 {
		return errors.New("certificate chain is empty")
	}
	for _, cert := range certs {
		if err := v.Validate(cert); err != nil {
			return err
		}
	}
	return nil
}
