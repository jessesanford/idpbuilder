// Package fallback provides detailed certificate chain logging for debugging
package fallback

import (
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"strings"
)

// LogCertificateChain logs comprehensive certificate chain information for debugging
func LogCertificateChain(certs []*x509.Certificate, logger *log.Logger) {
	if logger == nil || len(certs) == 0 {
		return
	}

	logger.Printf("Certificate Chain Analysis (%d certificates):", len(certs))
	logger.Printf("================================================")

	for i, cert := range certs {
		logger.Printf("")
		logger.Printf("Certificate %d of %d:", i+1, len(certs))
		logger.Printf("  Subject: %s", cert.Subject)
		logger.Printf("  Issuer:  %s", cert.Issuer)
		logger.Printf("  Serial Number: %s", cert.SerialNumber.String())
		logger.Printf("  Not Before: %s", cert.NotBefore.Format("2006-01-02 15:04:05 MST"))
		logger.Printf("  Not After:  %s", cert.NotAfter.Format("2006-01-02 15:04:05 MST"))
		
		// Log Subject Alternative Names
		if len(cert.DNSNames) > 0 {
			logger.Printf("  DNS Names: %s", strings.Join(cert.DNSNames, ", "))
		}
		
		if len(cert.IPAddresses) > 0 {
			ipStrings := make([]string, len(cert.IPAddresses))
			for j, ip := range cert.IPAddresses {
				ipStrings[j] = ip.String()
			}
			logger.Printf("  IP Addresses: %s", strings.Join(ipStrings, ", "))
		}
		
		// Certificate type and usage information
		logger.Printf("  Is CA: %v", cert.IsCA)
		if cert.MaxPathLen >= 0 {
			logger.Printf("  Max Path Length: %d", cert.MaxPathLen)
		}
		
		// Key usage information
		if cert.KeyUsage != 0 {
			usages := getKeyUsageStrings(cert.KeyUsage)
			logger.Printf("  Key Usage: %s", strings.Join(usages, ", "))
		}
		
		// Extended key usage
		if len(cert.ExtKeyUsage) > 0 {
			extUsages := getExtKeyUsageStrings(cert.ExtKeyUsage)
			logger.Printf("  Extended Key Usage: %s", strings.Join(extUsages, ", "))
		}
		
		// Signature algorithm
		logger.Printf("  Signature Algorithm: %s", cert.SignatureAlgorithm.String())
		
		// Public key information
		logger.Printf("  Public Key Algorithm: %s", cert.PublicKeyAlgorithm.String())
		
		// Certificate extensions
		if len(cert.Extensions) > 0 {
			logger.Printf("  Extensions: %d total", len(cert.Extensions))
		}
		
		// Self-signed check
		if cert.Subject.String() == cert.Issuer.String() {
			logger.Printf("  Note: Self-signed certificate")
		}
	}
	
	logger.Printf("================================================")
}

// LogValidationError logs detailed validation error information
func LogValidationError(err error, cert *x509.Certificate, logger *log.Logger) {
	if logger == nil || err == nil {
		return
	}

	logger.Printf("Certificate Validation Error:")
	logger.Printf("=============================")
	logger.Printf("Error Type: %T", err)
	logger.Printf("Error Message: %v", err)
	
	if cert != nil {
		logger.Printf("Certificate Details:")
		logger.Printf("  Subject: %s", cert.Subject)
		logger.Printf("  Issuer: %s", cert.Issuer)
		logger.Printf("  Serial: %s", cert.SerialNumber.String())
		logger.Printf("  Valid: %s to %s", 
			cert.NotBefore.Format("2006-01-02 15:04:05"),
			cert.NotAfter.Format("2006-01-02 15:04:05"))
		
		if len(cert.DNSNames) > 0 {
			logger.Printf("  Valid Hostnames: %s", strings.Join(cert.DNSNames, ", "))
		}
		
		if cert.Subject.CommonName != "" {
			logger.Printf("  Common Name: %s", cert.Subject.CommonName)
		}
	}
	
	logger.Printf("=============================")
}

// LogCertificateProblem logs detailed information about detected certificate problems
func LogCertificateProblem(problem *CertProblem, logger *log.Logger) {
	if logger == nil || problem == nil {
		return
	}

	logger.Printf("Certificate Problem Detected:")
	logger.Printf("============================")
	logger.Printf("Problem Type: %s", problem.Type)
	logger.Printf("Summary: %s", problem.GetProblemSummary())
	logger.Printf("Detailed Description: %s", problem.GetDetailedDescription())
	
	if problem.Error != nil {
		logger.Printf("Original Error: %v", problem.Error)
	}
	
	// Log problem details
	if len(problem.Details) > 0 {
		logger.Printf("Problem Details:")
		for key, value := range problem.Details {
			switch v := value.(type) {
			case []string:
				logger.Printf("  %s: %s", key, strings.Join(v, ", "))
			case []net.IP:
				ips := make([]string, len(v))
				for i, ip := range v {
					ips[i] = ip.String()
				}
				logger.Printf("  %s: %s", key, strings.Join(ips, ", "))
			default:
				logger.Printf("  %s: %v", key, value)
			}
		}
	}
	
	// Log suggestions
	if len(problem.Suggestions) > 0 {
		logger.Printf("Suggestions:")
		for i, suggestion := range problem.Suggestions {
			logger.Printf("  %d. %s", i+1, suggestion)
		}
	}
	
	logger.Printf("============================")
}

// LogRegistryConnection logs registry connection attempts and results
func LogRegistryConnection(registryURL string, insecure bool, success bool, err error, logger *log.Logger) {
	if logger == nil {
		return
	}

	status := "SUCCESS"
	if !success {
		status = "FAILED"
	}
	
	securityMode := "SECURE"
	if insecure {
		securityMode = "INSECURE"
	}
	
	logger.Printf("Registry Connection [%s/%s]: %s", status, securityMode, registryURL)
	
	if err != nil {
		logger.Printf("  Error: %v", err)
	}
	
	if insecure {
		logger.Printf("  WARNING: TLS verification disabled")
	}
}

// getKeyUsageStrings converts x509.KeyUsage bitmask to human-readable strings
func getKeyUsageStrings(usage x509.KeyUsage) []string {
	var usages []string
	
	if usage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "Digital Signature")
	}
	if usage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "Content Commitment")
	}
	if usage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "Key Encipherment")
	}
	if usage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "Data Encipherment")
	}
	if usage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "Key Agreement")
	}
	if usage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "Certificate Signing")
	}
	if usage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRL Signing")
	}
	if usage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "Encipher Only")
	}
	if usage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "Decipher Only")
	}
	
	return usages
}

// getExtKeyUsageStrings converts extended key usage to human-readable strings
func getExtKeyUsageStrings(extUsages []x509.ExtKeyUsage) []string {
	var usages []string
	
	for _, extUsage := range extUsages {
		switch extUsage {
		case x509.ExtKeyUsageServerAuth:
			usages = append(usages, "TLS Server Authentication")
		case x509.ExtKeyUsageClientAuth:
			usages = append(usages, "TLS Client Authentication")
		case x509.ExtKeyUsageCodeSigning:
			usages = append(usages, "Code Signing")
		case x509.ExtKeyUsageEmailProtection:
			usages = append(usages, "Email Protection")
		case x509.ExtKeyUsageTimeStamping:
			usages = append(usages, "Time Stamping")
		case x509.ExtKeyUsageOCSPSigning:
			usages = append(usages, "OCSP Signing")
		default:
			usages = append(usages, fmt.Sprintf("Unknown (%d)", extUsage))
		}
	}
	
	return usages
}