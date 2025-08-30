package cert_integration

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CertificateLoader defines the interface for loading certificates from various sources
type CertificateLoader interface {
	// LoadFromExtraction loads certificates from Phase 1 cert-extraction component
	LoadFromExtraction(source string) (*CertificateSet, error)
	
	// LoadFromPath loads certificates from a filesystem path
	LoadFromPath(path string) (*CertificateSet, error)
	
	// LoadFromKindCluster loads certificates from a Kind cluster
	LoadFromKindCluster(clusterName string) (*CertificateSet, error)
}

// certificateLoader implements the CertificateLoader interface
type certificateLoader struct {
	config *ManagerConfig
}

// NewCertificateLoader creates a new certificate loader with the given configuration
func NewCertificateLoader(config *ManagerConfig) CertificateLoader {
	return &certificateLoader{
		config: config,
	}
}

// LoadFromExtraction loads certificates from Phase 1 cert-extraction component
func (cl *certificateLoader) LoadFromExtraction(source string) (*CertificateSet, error) {
	if source == "" {
		return nil, fmt.Errorf("extraction source cannot be empty")
	}
	
	// This would integrate with Phase 1 cert-extraction.Client
	// For now, we'll implement a file-based approach that mimics the extraction output
	extractionPath := filepath.Join(cl.config.DefaultCertPath, "extraction", source)
	
	if _, err := os.Stat(extractionPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("extraction source path does not exist: %s", extractionPath)
	}
	
	return cl.loadCertificatesFromDirectory(extractionPath)
}

// LoadFromPath loads certificates from a filesystem path
func (cl *certificateLoader) LoadFromPath(path string) (*CertificateSet, error) {
	if path == "" {
		return nil, fmt.Errorf("certificate path cannot be empty")
	}
	
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to access certificate path %s: %w", path, err)
	}
	
	if info.IsDir() {
		return cl.loadCertificatesFromDirectory(path)
	} else {
		return cl.loadCertificatesFromFile(path)
	}
}

// LoadFromKindCluster loads certificates from a Kind cluster
func (cl *certificateLoader) LoadFromKindCluster(clusterName string) (*CertificateSet, error) {
	if clusterName == "" {
		return nil, fmt.Errorf("cluster name cannot be empty")
	}
	
	// This would integrate with Phase 1 cert-extraction client to extract from Kind
	// For now, we'll look for certificates in the standard Kind location
	kindPath := filepath.Join(cl.config.DefaultCertPath, "kind", clusterName)
	
	if _, err := os.Stat(kindPath); os.IsNotExist(err) {
		// Try alternative locations
		homeDir, _ := os.UserHomeDir()
		kindPath = filepath.Join(homeDir, ".kind", clusterName, "certs")
		
		if _, err := os.Stat(kindPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("no certificates found for Kind cluster: %s", clusterName)
		}
	}
	
	return cl.loadCertificatesFromDirectory(kindPath)
}

// loadCertificatesFromDirectory loads all certificates from a directory
func (cl *certificateLoader) loadCertificatesFromDirectory(dirPath string) (*CertificateSet, error) {
	certSet := &CertificateSet{
		ClientCerts: make(map[string]*x509.Certificate),
		Metadata:    make(map[string]interface{}),
	}
	
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() {
			return nil
		}
		
		// Only process certificate files
		if !isCertificateFile(path) {
			return nil
		}
		
		cert, certType, err := cl.loadCertificate(path)
		if err != nil {
			return fmt.Errorf("failed to load certificate from %s: %w", path, err)
		}
		
		if cert == nil {
			return nil // Skip non-certificate files
		}
		
		// Categorize the certificate based on its type and usage
		if err := cl.categorizeCertificate(cert, certType, path, certSet); err != nil {
			return fmt.Errorf("failed to categorize certificate from %s: %w", path, err)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to load certificates from directory %s: %w", dirPath, err)
	}
	
	// Build the complete certificate set
	if err := cl.buildCertificateSet(certSet); err != nil {
		return nil, fmt.Errorf("failed to build certificate set: %w", err)
	}
	
	return certSet, nil
}

// loadCertificatesFromFile loads certificates from a single file
func (cl *certificateLoader) loadCertificatesFromFile(filePath string) (*CertificateSet, error) {
	cert, certType, err := cl.loadCertificate(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate from %s: %w", filePath, err)
	}
	
	if cert == nil {
		return nil, fmt.Errorf("no valid certificate found in file: %s", filePath)
	}
	
	certSet := &CertificateSet{
		ClientCerts: make(map[string]*x509.Certificate),
		Metadata:    make(map[string]interface{}),
	}
	
	// Categorize the single certificate
	if err := cl.categorizeCertificate(cert, certType, filePath, certSet); err != nil {
		return nil, fmt.Errorf("failed to categorize certificate: %w", err)
	}
	
	// Build the certificate set
	if err := cl.buildCertificateSet(certSet); err != nil {
		return nil, fmt.Errorf("failed to build certificate set: %w", err)
	}
	
	return certSet, nil
}

// loadCertificate loads a single certificate from a file
func (cl *certificateLoader) loadCertificate(filePath string) (*x509.Certificate, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}
	
	// Try PEM format first
	block, _ := pem.Decode(data)
	if block != nil {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, "", fmt.Errorf("failed to parse PEM certificate: %w", err)
		}
		return cert, "pem", nil
	}
	
	// Try DER format
	cert, err := x509.ParseCertificate(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse DER certificate: %w", err)
	}
	
	return cert, "der", nil
}

// categorizeCertificate determines the type and role of a certificate
func (cl *certificateLoader) categorizeCertificate(cert *x509.Certificate, certType, filePath string, certSet *CertificateSet) error {
	filename := filepath.Base(filePath)
	
	// Determine certificate role based on file name and certificate properties
	if cert.IsCA {
		if strings.Contains(filename, "root") || strings.Contains(filename, "ca") {
			if certSet.RootCA == nil {
				certSet.RootCA = cert
			} else {
				// Additional root CAs go to intermediates
				certSet.Intermediates = append(certSet.Intermediates, cert)
			}
		} else {
			certSet.Intermediates = append(certSet.Intermediates, cert)
		}
	} else {
		// Non-CA certificates
		if strings.Contains(filename, "server") || strings.Contains(filename, "tls") {
			certSet.ServerCert = cert
		} else if strings.Contains(filename, "client") {
			// Extract client name from filename
			clientName := strings.TrimSuffix(filename, filepath.Ext(filename))
			certSet.ClientCerts[clientName] = cert
		} else {
			// Default to server certificate if unclear
			if certSet.ServerCert == nil {
				certSet.ServerCert = cert
			} else {
				// Additional certificates become client certificates
				certSet.ClientCerts[filename] = cert
			}
		}
	}
	
	// Store metadata
	certSet.Metadata[filePath] = map[string]interface{}{
		"format":      certType,
		"subject":     cert.Subject.String(),
		"issuer":      cert.Issuer.String(),
		"not_before":  cert.NotBefore,
		"not_after":   cert.NotAfter,
		"is_ca":       cert.IsCA,
		"serial":      cert.SerialNumber.String(),
	}
	
	return nil
}

// buildCertificateSet finalizes the certificate set by building the trust bundle
func (cl *certificateLoader) buildCertificateSet(certSet *CertificateSet) error {
	// Build trust bundle from root CA and intermediates
	var trustBundle []byte
	
	if certSet.RootCA != nil {
		rootPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certSet.RootCA.Raw,
		})
		trustBundle = append(trustBundle, rootPEM...)
	}
	
	for _, intermediate := range certSet.Intermediates {
		intermediatePEM := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: intermediate.Raw,
		})
		trustBundle = append(trustBundle, intermediatePEM...)
	}
	
	certSet.TrustBundle = trustBundle
	
	return nil
}

// isCertificateFile determines if a file is likely a certificate file based on extension
func isCertificateFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	certificateExtensions := []string{".crt", ".cer", ".pem", ".der", ".cert", ".ca", ".csr"}
	
	for _, certExt := range certificateExtensions {
		if ext == certExt {
			return true
		}
	}
	
	return false
}