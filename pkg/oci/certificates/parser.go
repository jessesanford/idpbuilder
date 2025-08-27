package certificates

import (
	"crypto/x509"
	"errors"
	"fmt"
	"sort"
	"time"

	v2 "github.com/cnoe-io/idpbuilder/pkg/oci/api/v2"
)

// CertificateParser provides parsing utilities
type CertificateParser struct {
	strictMode    bool
	maxChainDepth int
}

// NewCertificateParser creates a new certificate parser with default settings
func NewCertificateParser() *CertificateParser {
	return &CertificateParser{
		strictMode:    false,
		maxChainDepth: 10,
	}
}

// ParseCertificateChain parses and validates a certificate chain
func (p *CertificateParser) ParseCertificateChain(certs []*x509.Certificate) ([]*x509.Certificate, error) {
	if len(certs) == 0 {
		return nil, &v2.CertificateError{
			Code:    "EMPTY_CHAIN",
			Message: "certificate chain is empty",
		}
	}

	if len(certs) > p.maxChainDepth {
		return nil, &v2.CertificateError{
			Code:    "CHAIN_TOO_DEEP",
			Message: fmt.Sprintf("certificate chain depth %d exceeds maximum %d", len(certs), p.maxChainDepth),
		}
	}

	// Validate each certificate in the chain
	for i, cert := range certs {
		if err := p.ValidateCertificate(cert); err != nil {
			return nil, &v2.CertificateError{
				Code:    "CHAIN_VALIDATION_ERROR",
				Message: fmt.Sprintf("certificate %d in chain failed validation: %v", i, err),
				Cert:    cert,
				Err:     err,
			}
		}
	}

	// Sort certificates in chain order (leaf to root)
	sortedChain := p.SortByHierarchy(certs)

	// Validate the chain structure if in strict mode
	if p.strictMode {
		if err := p.ValidateChain(sortedChain); err != nil {
			return nil, err
		}
	}

	return sortedChain, nil
}

// BuildChain builds a complete certificate chain from a collection
func (p *CertificateParser) BuildChain(leaf *x509.Certificate, intermediates []*x509.Certificate) ([]*x509.Certificate, error) {
	if leaf == nil {
		return nil, &v2.CertificateError{
			Code:    "NULL_LEAF_CERT",
			Message: "leaf certificate is nil",
		}
	}

	chain := []*x509.Certificate{leaf}
	used := make(map[string]bool)
	used[string(leaf.Signature)] = true

	// Build chain by finding issuers
	current := leaf
	for len(chain) < p.maxChainDepth {
		var issuer *x509.Certificate

		// Look for the issuer among intermediates
		for _, intermediate := range intermediates {
			if used[string(intermediate.Signature)] {
				continue // Already used this certificate
			}

			// Check if this intermediate issued the current certificate
			if p.isIssuer(current, intermediate) {
				issuer = intermediate
				break
			}
		}

		if issuer == nil {
			// No more issuers found, chain might be incomplete
			break
		}

		chain = append(chain, issuer)
		used[string(issuer.Signature)] = true
		current = issuer

		// Stop if we found a self-signed certificate (root CA)
		if p.isSelfSigned(current) {
			break
		}
	}

	return chain, nil
}

// ValidateCertificate performs comprehensive certificate validation
func (p *CertificateParser) ValidateCertificate(cert *x509.Certificate) error {
	if cert == nil {
		return &v2.CertificateError{
			Code:    "NULL_CERTIFICATE",
			Message: "certificate is nil",
		}
	}

	now := time.Now()

	// Check if certificate is expired
	if cert.NotAfter.Before(now) {
		if p.strictMode {
			return &v2.CertificateError{
				Code:    "CERT_EXPIRED",
				Message: fmt.Sprintf("certificate expired at %v", cert.NotAfter),
				Cert:    cert,
			}
		}
	}

	// Check if certificate is not yet valid
	if cert.NotBefore.After(now) {
		if p.strictMode {
			return &v2.CertificateError{
				Code:    "CERT_NOT_YET_VALID",
				Message: fmt.Sprintf("certificate not valid until %v", cert.NotBefore),
				Cert:    cert,
			}
		}
	}

	// Validate certificate structure
	if cert.Subject.String() == "" {
		return &v2.CertificateError{
			Code:    "INVALID_SUBJECT",
			Message: "certificate has empty subject",
			Cert:    cert,
		}
	}

	// Validate key usage for non-CA certificates
	if !cert.IsCA {
		if cert.KeyUsage == 0 {
			return &v2.CertificateError{
				Code:    "MISSING_KEY_USAGE",
				Message: "certificate has no key usage defined",
				Cert:    cert,
			}
		}
	}

	// Validate CA certificates have proper constraints
	if cert.IsCA {
		if cert.KeyUsage&x509.KeyUsageCertSign == 0 {
			return &v2.CertificateError{
				Code:    "INVALID_CA_KEY_USAGE",
				Message: "CA certificate missing cert sign key usage",
				Cert:    cert,
			}
		}
	}

	return nil
}

// ExtractCAs separates CA certificates from end-entity certificates
func (p *CertificateParser) ExtractCAs(certs []*x509.Certificate) ([]*x509.Certificate, []*x509.Certificate) {
	var cas []*x509.Certificate
	var endEntity []*x509.Certificate

	for _, cert := range certs {
		if cert.IsCA {
			cas = append(cas, cert)
		} else {
			endEntity = append(endEntity, cert)
		}
	}

	return cas, endEntity
}

// SortByHierarchy sorts certificates in chain order (leaf to root)
func (p *CertificateParser) SortByHierarchy(certs []*x509.Certificate) []*x509.Certificate {
	if len(certs) <= 1 {
		return certs
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]*x509.Certificate, len(certs))
	copy(sorted, certs)

	// Sort by chain hierarchy
	sort.Slice(sorted, func(i, j int) bool {
		certI := sorted[i]
		certJ := sorted[j]

		// If one is the issuer of the other, put issuer later
		if p.isIssuer(certI, certJ) {
			return true // certI is issued by certJ, so certI comes first
		}
		if p.isIssuer(certJ, certI) {
			return false // certJ is issued by certI, so certJ comes first
		}

		// If both are CAs or both are end-entity, sort by NotBefore (newer first)
		return certI.NotBefore.After(certJ.NotBefore)
	})

	return sorted
}

// ValidateChain validates a complete certificate chain
func (p *CertificateParser) ValidateChain(chain []*x509.Certificate) error {
	if len(chain) == 0 {
		return &v2.CertificateError{
			Code:    "EMPTY_CHAIN",
			Message: "certificate chain is empty",
		}
	}

	// Validate each certificate in the chain
	for i, cert := range chain {
		if err := p.ValidateCertificate(cert); err != nil {
			return &v2.CertificateError{
				Code:    "CHAIN_CERT_INVALID",
				Message: fmt.Sprintf("certificate %d in chain is invalid: %v", i, err),
				Cert:    cert,
				Err:     err,
			}
		}
	}

	// Validate chain links (each certificate is issued by the next)
	for i := 0; i < len(chain)-1; i++ {
		current := chain[i]
		next := chain[i+1]

		if !p.isIssuer(current, next) {
			return &v2.CertificateError{
				Code:    "BROKEN_CHAIN",
				Message: fmt.Sprintf("certificate %d is not issued by certificate %d", i, i+1),
				Cert:    current,
			}
		}
	}

	// Check if the root is self-signed
	root := chain[len(chain)-1]
	if !p.isSelfSigned(root) && p.strictMode {
		return &v2.CertificateError{
			Code:    "INCOMPLETE_CHAIN",
			Message: "certificate chain does not end with a self-signed root CA",
			Cert:    root,
		}
	}

	return nil
}

// ConvertToBundle creates a CertBundle from parsed certificates
func (p *CertificateParser) ConvertToBundle(certs []*x509.Certificate, format v2.CertFormat, source string) *v2.CertBundle {
	if len(certs) == 0 {
		return &v2.CertBundle{
			Certificates: []*x509.Certificate{},
			CAs:          []*x509.Certificate{},
			Format:       format,
			LoadedAt:     time.Now(),
			Source:       source,
		}
	}

	// Separate CAs from end-entity certificates
	cas, endEntity := p.ExtractCAs(certs)

	return &v2.CertBundle{
		Certificates: endEntity,
		CAs:          cas,
		Format:       format,
		LoadedAt:     time.Now(),
		Source:       source,
	}
}

// isIssuer checks if the issuer certificate issued the subject certificate
func (p *CertificateParser) isIssuer(subject, issuer *x509.Certificate) bool {
	// Check if issuer's subject matches subject's issuer
	if !subject.Issuer.Equal(issuer.Subject) {
		return false
	}

	// Verify signature
	err := subject.CheckSignatureFrom(issuer)
	return err == nil
}

// isSelfSigned checks if a certificate is self-signed
func (p *CertificateParser) isSelfSigned(cert *x509.Certificate) bool {
	// Check if subject equals issuer
	if !cert.Subject.Equal(cert.Issuer) {
		return false
	}

	// Verify self-signature
	err := cert.CheckSignatureFrom(cert)
	return err == nil
}

// GetCertificateInfo returns human-readable information about a certificate
func (p *CertificateParser) GetCertificateInfo(cert *x509.Certificate) map[string]interface{} {
	if cert == nil {
		return map[string]interface{}{"error": "certificate is nil"}
	}

	info := map[string]interface{}{
		"subject":            cert.Subject.String(),
		"issuer":             cert.Issuer.String(),
		"serial_number":      cert.SerialNumber.String(),
		"not_before":         cert.NotBefore,
		"not_after":          cert.NotAfter,
		"is_ca":              cert.IsCA,
		"key_usage":          p.keyUsageToString(cert.KeyUsage),
		"extended_key_usage": p.extKeyUsageToString(cert.ExtKeyUsage),
		"dns_names":          cert.DNSNames,
		"ip_addresses":       cert.IPAddresses,
		"email_addresses":    cert.EmailAddresses,
		"is_self_signed":     p.isSelfSigned(cert),
	}

	// Add expiry status
	now := time.Now()
	if cert.NotAfter.Before(now) {
		info["status"] = "expired"
	} else if cert.NotBefore.After(now) {
		info["status"] = "not_yet_valid"
	} else {
		info["status"] = "valid"
	}

	return info
}

// keyUsageToString converts key usage flags to strings
func (p *CertificateParser) keyUsageToString(usage x509.KeyUsage) []string {
	var usages []string

	if usage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "digital_signature")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "content_commitment")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "key_encipherment")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "data_encipherment")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "key_agreement")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "cert_sign")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "crl_sign")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "encipher_only")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "decipher_only")
	}

	return usages
}

// extKeyUsageToString converts extended key usage to strings
func (p *CertificateParser) extKeyUsageToString(usage []x509.ExtKeyUsage) []string {
	var usages []string

	for _, u := range usage {
		switch u {
		case x509.ExtKeyUsageServerAuth:
			usages = append(usages, "server_auth")
		case x509.ExtKeyUsageClientAuth:
			usages = append(usages, "client_auth")
		case x509.ExtKeyUsageCodeSigning:
			usages = append(usages, "code_signing")
		case x509.ExtKeyUsageEmailProtection:
			usages = append(usages, "email_protection")
		case x509.ExtKeyUsageTimeStamping:
			usages = append(usages, "time_stamping")
		case x509.ExtKeyUsageOCSPSigning:
			usages = append(usages, "ocsp_signing")
		}
	}

	return usages
}
