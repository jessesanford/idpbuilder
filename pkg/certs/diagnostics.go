package certs

import (
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"
)

// GenerateDiagnostics creates comprehensive diagnostic report
func (v *DefaultValidator) GenerateDiagnostics(cert *x509.Certificate) (*CertDiagnostics, error) {
	if cert == nil {
		return nil, fmt.Errorf("certificate cannot be nil")
	}

	diag := &CertDiagnostics{
		Subject:          cert.Subject.String(),
		Issuer:           cert.Issuer.String(),
		SerialNumber:     cert.SerialNumber.String(),
		NotBefore:        cert.NotBefore,
		NotAfter:         cert.NotAfter,
		DNSNames:         make([]string, len(cert.DNSNames)),
		IPAddresses:      make([]net.IP, len(cert.IPAddresses)),
		ValidationErrors: make([]ValidationError, 0),
		Warnings:         make([]string, 0),
	}

	// Copy DNS names and IP addresses
	copy(diag.DNSNames, cert.DNSNames)
	copy(diag.IPAddresses, cert.IPAddresses)

	// Run all validations and collect errors and warnings

	// 1. Chain validation
	if err := v.ValidateChain(cert); err != nil {
		diag.ValidationErrors = append(diag.ValidationErrors, ValidationError{
			Type:    "chain",
			Message: "Certificate chain validation failed",
			Detail:  err.Error(),
		})
	}

	// 2. Expiry validation
	duration, err := v.CheckExpiry(cert)
	if err != nil {
		// Distinguish between expired and expiring soon
		if duration == nil {
			// Certificate is expired or not yet valid
			diag.ValidationErrors = append(diag.ValidationErrors, ValidationError{
				Type:    "expiry",
				Message: "Certificate expiry validation failed",
				Detail:  err.Error(),
			})
		} else {
			// Certificate is expiring soon - this is a warning, not an error
			diag.Warnings = append(diag.Warnings, err.Error())
		}
	} else if duration != nil {
		// Certificate is valid - add info about remaining time
		diag.Warnings = append(diag.Warnings, fmt.Sprintf("Certificate expires in %v (%s)",
			*duration, cert.NotAfter.Format(time.RFC3339)))
	}

	// 3. Add basic certificate information as warnings/notes
	if cert.IsCA {
		diag.Warnings = append(diag.Warnings, "This is a Certificate Authority (CA) certificate")
	}

	if len(cert.DNSNames) == 0 && cert.Subject.CommonName != "" {
		diag.Warnings = append(diag.Warnings, "Certificate uses Common Name instead of Subject Alternative Names (deprecated)")
	}

	// 4. Check for self-signed certificate
	if cert.Subject.String() == cert.Issuer.String() {
		diag.Warnings = append(diag.Warnings, "Certificate is self-signed")
	}

	// 5. Check key usage
	if cert.KeyUsage == 0 {
		diag.Warnings = append(diag.Warnings, "Certificate has no key usage specified")
	}

	return diag, nil
}

// FormatDiagnostics returns human-readable diagnostic output
func FormatDiagnostics(diag *CertDiagnostics) string {
	if diag == nil {
		return "No diagnostic data available"
	}

	var sb strings.Builder

	// Header
	sb.WriteString("Certificate Diagnostic Report\n")
	sb.WriteString("============================\n\n")

	// Basic Information
	sb.WriteString("Basic Information:\n")
	sb.WriteString(fmt.Sprintf("  Subject: %s\n", diag.Subject))
	sb.WriteString(fmt.Sprintf("  Issuer:  %s\n", diag.Issuer))
	sb.WriteString(fmt.Sprintf("  Serial:  %s\n", diag.SerialNumber))
	sb.WriteString(fmt.Sprintf("  Valid:   %s to %s\n",
		diag.NotBefore.Format(time.RFC3339),
		diag.NotAfter.Format(time.RFC3339)))

	// Validity period
	now := time.Now()
	if now.Before(diag.NotBefore) {
		sb.WriteString(fmt.Sprintf("  Status:  Not yet valid (valid in %v)\n", diag.NotBefore.Sub(now)))
	} else if now.After(diag.NotAfter) {
		sb.WriteString(fmt.Sprintf("  Status:  EXPIRED (%v ago)\n", now.Sub(diag.NotAfter)))
	} else {
		sb.WriteString(fmt.Sprintf("  Status:  Valid (expires in %v)\n", diag.NotAfter.Sub(now)))
	}

	sb.WriteString("\n")

	// Hostnames
	if len(diag.DNSNames) > 0 {
		sb.WriteString("Valid Hostnames:\n")
		for _, name := range diag.DNSNames {
			sb.WriteString(fmt.Sprintf("  - %s\n", name))
		}
		sb.WriteString("\n")
	}

	// IP Addresses
	if len(diag.IPAddresses) > 0 {
		sb.WriteString("Valid IP Addresses:\n")
		for _, ip := range diag.IPAddresses {
			sb.WriteString(fmt.Sprintf("  - %s\n", ip.String()))
		}
		sb.WriteString("\n")
	}

	// Validation Errors
	if len(diag.ValidationErrors) > 0 {
		sb.WriteString("Validation Errors:\n")
		for i, err := range diag.ValidationErrors {
			sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, strings.ToUpper(err.Type), err.Message))
			if err.Detail != "" {
				// Indent detail lines
				detail := strings.ReplaceAll(err.Detail, "\n", "\n      ")
				sb.WriteString(fmt.Sprintf("      %s\n", detail))
			}
		}
		sb.WriteString("\n")
	}

	// Warnings
	if len(diag.Warnings) > 0 {
		sb.WriteString("Warnings and Notes:\n")
		for i, warning := range diag.Warnings {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, warning))
		}
		sb.WriteString("\n")
	}

	// Summary
	sb.WriteString("Summary:\n")
	if len(diag.ValidationErrors) == 0 {
		sb.WriteString("  ✓ Certificate passed all validations\n")
	} else {
		sb.WriteString(fmt.Sprintf("  ✗ Certificate has %d validation error(s)\n", len(diag.ValidationErrors)))
	}

	if len(diag.Warnings) > 0 {
		sb.WriteString(fmt.Sprintf("  ⚠ %d warning(s) or informational note(s)\n", len(diag.Warnings)))
	}

	return sb.String()
}

// ValidateAndDiagnose is a convenience method that performs validation and returns diagnostics
func (v *DefaultValidator) ValidateAndDiagnose(cert *x509.Certificate, hostname string) (*CertDiagnostics, error) {
	// Generate base diagnostics
	diag, err := v.GenerateDiagnostics(cert)
	if err != nil {
		return nil, fmt.Errorf("failed to generate diagnostics: %w", err)
	}

	// Add hostname validation if provided
	if hostname != "" {
		if err := v.VerifyHostname(cert, hostname); err != nil {
			diag.ValidationErrors = append(diag.ValidationErrors, ValidationError{
				Type:    "hostname",
				Message: fmt.Sprintf("Hostname '%s' validation failed", hostname),
				Detail:  err.Error(),
			})
		} else {
			diag.Warnings = append(diag.Warnings, fmt.Sprintf("Hostname '%s' matches certificate", hostname))
		}
	}

	return diag, nil
}
