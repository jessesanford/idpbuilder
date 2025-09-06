// Package fallback provides detailed certificate chain logging for debugging
package fallback

import (
	"crypto/x509"
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
	for i, cert := range certs {
		logger.Printf("Certificate %d: Subject=%s, Issuer=%s", i+1, cert.Subject, cert.Issuer)
		logger.Printf("  Valid: %s to %s", cert.NotBefore.Format("2006-01-02"), cert.NotAfter.Format("2006-01-02"))

		if len(cert.DNSNames) > 0 {
			logger.Printf("  DNS Names: %s", strings.Join(cert.DNSNames, ", "))
		}

		if len(cert.IPAddresses) > 0 {
			ips := make([]string, len(cert.IPAddresses))
			for j, ip := range cert.IPAddresses {
				ips[j] = ip.String()
			}
			logger.Printf("  IP Addresses: %s", strings.Join(ips, ", "))
		}

		if cert.IsCA {
			logger.Printf("  Type: Certificate Authority")
		}

		if cert.Subject.String() == cert.Issuer.String() {
			logger.Printf("  Note: Self-signed certificate")
		}
	}
}

// LogValidationError logs detailed validation error information
func LogValidationError(err error, cert *x509.Certificate, logger *log.Logger) {
	if logger == nil || err == nil {
		return
	}

	logger.Printf("Certificate Validation Error: %v", err)
	if cert != nil {
		logger.Printf("  Subject: %s", cert.Subject)
		logger.Printf("  Valid: %s to %s",
			cert.NotBefore.Format("2006-01-02"), cert.NotAfter.Format("2006-01-02"))
		if len(cert.DNSNames) > 0 {
			logger.Printf("  Valid Hostnames: %s", strings.Join(cert.DNSNames, ", "))
		}
	}
}

// LogCertificateProblem logs detailed information about detected certificate problems
func LogCertificateProblem(problem *CertProblem, logger *log.Logger) {
	if logger == nil || problem == nil {
		return
	}

	logger.Printf("Certificate Problem: %s - %s", problem.Type, problem.GetProblemSummary())
	if problem.Error != nil {
		logger.Printf("  Original Error: %v", problem.Error)
	}

	// Log key problem details
	for key, value := range problem.Details {
		switch v := value.(type) {
		case []string:
			if len(v) > 0 {
				logger.Printf("  %s: %s", key, strings.Join(v, ", "))
			}
		case []net.IP:
			if len(v) > 0 {
				ips := make([]string, len(v))
				for i, ip := range v {
					ips[i] = ip.String()
				}
				logger.Printf("  %s: %s", key, strings.Join(ips, ", "))
			}
		default:
			logger.Printf("  %s: %v", key, value)
		}
	}
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
