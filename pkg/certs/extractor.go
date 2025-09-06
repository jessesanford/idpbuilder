// pkg/certs/extractor.go
package certs

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DefaultExtractor implements KindCertExtractor interface
type DefaultExtractor struct {
	clusterName string
	namespace   string
}

// NewDefaultExtractor creates a new DefaultExtractor instance
func NewDefaultExtractor(clusterName string) *DefaultExtractor {
	return &DefaultExtractor{
		clusterName: clusterName,
		namespace:   "gitea", // Default namespace for Gitea in idpbuilder
	}
}

// ExtractGiteaCert extracts the Gitea certificate from the Kind cluster
func (e *DefaultExtractor) ExtractGiteaCert(ctx context.Context) (*x509.Certificate, error) {
	// Step 1: Verify Kind cluster exists
	_, err := e.GetClusterName()
	if err != nil {
		return nil, err
	}

	// Step 2: Find Gitea pod
	podName, err := e.findGiteaPod(ctx)
	if err != nil {
		return nil, err
	}

	// Step 3: Extract certificate from pod
	certData, err := e.extractCertFromPod(ctx, podName)
	if err != nil {
		return nil, err
	}

	// Step 4: Parse PEM certificate
	cert, err := e.parsePEMCertificate(certData)
	if err != nil {
		return nil, CertificateInvalidError{Reason: fmt.Sprintf("failed to parse certificate: %v", err)}
	}

	// Step 5: Validate certificate
	if err := e.ValidateCertificate(cert); err != nil {
		return nil, err
	}

	return cert, nil
}

// GetClusterName returns the name of the Kind cluster
func (e *DefaultExtractor) GetClusterName() (string, error) {
	if e.clusterName != "" {
		// Verify the specified cluster exists
		return e.verifyClusterExists(e.clusterName)
	}

	// Auto-detect Kind cluster
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Kind clusters: %v", err)
	}

	clusters := strings.Fields(strings.TrimSpace(string(output)))
	if len(clusters) == 0 {
		return "", ClusterNotFoundError{ClusterName: "no clusters found"}
	}

	// Use the first available cluster
	return clusters[0], nil
}

// ValidateCertificate performs basic validation on the extracted certificate
func (e *DefaultExtractor) ValidateCertificate(cert *x509.Certificate) error {
	if cert == nil {
		return CertificateInvalidError{Reason: "certificate is nil"}
	}

	// Check if certificate is expired
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return CertificateInvalidError{
			Reason: fmt.Sprintf("certificate not yet valid (valid from %v)", cert.NotBefore),
		}
	}

	if now.After(cert.NotAfter) {
		return CertificateInvalidError{
			Reason: fmt.Sprintf("certificate expired on %v", cert.NotAfter),
		}
	}

	// Check if certificate has the expected subject or DNS names
	hasGiteaIdentity := false
	for _, dnsName := range cert.DNSNames {
		if strings.Contains(strings.ToLower(dnsName), "gitea") {
			hasGiteaIdentity = true
			break
		}
	}

	if !hasGiteaIdentity && !strings.Contains(strings.ToLower(cert.Subject.CommonName), "gitea") {
		return CertificateInvalidError{
			Reason: "certificate does not appear to be for Gitea service",
		}
	}

	return nil
}

// SaveCertificate saves the certificate to the local trust store
func (e *DefaultExtractor) SaveCertificate(cert *x509.Certificate, path string) error {
	if cert == nil {
		return CertificateInvalidError{Reason: "certificate is nil"}
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return PermissionError{Path: dir, Action: "create directory"}
	}

	// Convert certificate to PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})

	// Write to file with restrictive permissions
	if err := os.WriteFile(path, certPEM, 0600); err != nil {
		return PermissionError{Path: path, Action: "write"}
	}

	return nil
}

// findGiteaPod finds the Gitea pod in the cluster
func (e *DefaultExtractor) findGiteaPod(ctx context.Context) (string, error) {
	// Use kubectl to find Gitea pod
	cmd := exec.CommandContext(ctx, "kubectl", "get", "pods", "-n", e.namespace, "-l", "app=gitea", "-o", "jsonpath={.items[0].metadata.name}")
	output, err := cmd.Output()
	if err != nil {
		// Try with different label selector
		cmd = exec.CommandContext(ctx, "kubectl", "get", "pods", "-n", e.namespace, "--field-selector=status.phase=Running", "-o", "jsonpath={.items[?(@.metadata.name=~'gitea.*')].metadata.name}")
		output, err = cmd.Output()
		if err != nil {
			return "", PodNotFoundError{PodName: "gitea", Namespace: e.namespace}
		}
	}

	podName := strings.TrimSpace(string(output))
	if podName == "" {
		return "", PodNotFoundError{PodName: "gitea", Namespace: e.namespace}
	}

	return podName, nil
}

// extractCertFromPod extracts the certificate file from the Gitea pod
func (e *DefaultExtractor) extractCertFromPod(ctx context.Context, podName string) ([]byte, error) {
	// Try multiple possible certificate paths
	certPaths := []string{
		"/data/gitea/https/cert.pem",
		"/data/git/https/cert.pem",
		"/etc/ssl/certs/gitea.crt",
		"/app/cert/tls.crt",
	}

	for _, certPath := range certPaths {
		cmd := exec.CommandContext(ctx, "kubectl", "exec", "-n", e.namespace, podName, "--", "cat", certPath)
		output, err := cmd.Output()
		if err != nil {
			// Try next path
			continue
		}

		// Check if output looks like a certificate
		if strings.Contains(string(output), "BEGIN CERTIFICATE") {
			return output, nil
		}
	}

	return nil, fmt.Errorf("could not find certificate in pod %s", podName)
}

// parsePEMCertificate parses a PEM-encoded certificate
func (e *DefaultExtractor) parsePEMCertificate(certData []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM certificate")
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("PEM block is not a certificate (got %s)", block.Type)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X.509 certificate: %v", err)
	}

	return cert, nil
}

// verifyClusterExists verifies that the specified cluster exists
func (e *DefaultExtractor) verifyClusterExists(clusterName string) (string, error) {
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Kind clusters: %v", err)
	}

	clusters := strings.Fields(strings.TrimSpace(string(output)))
	for _, cluster := range clusters {
		if cluster == clusterName {
			return clusterName, nil
		}
	}

	return "", ClusterNotFoundError{ClusterName: clusterName}
}

// GetCertificateInfo extracts metadata from a certificate
func GetCertificateInfo(cert *x509.Certificate) CertificateInfo {
	if cert == nil {
		return CertificateInfo{}
	}

	return CertificateInfo{
		Subject:   cert.Subject.String(),
		Issuer:    cert.Issuer.String(),
		NotBefore: cert.NotBefore,
		NotAfter:  cert.NotAfter,
		IsCA:      cert.IsCA,
		DNSNames:  cert.DNSNames,
	}
}

// ExtractAndSave is a convenience function that extracts and saves the certificate
func ExtractAndSave(ctx context.Context, clusterName, savePath string) error {
	extractor := NewDefaultExtractor(clusterName)

	cert, err := extractor.ExtractGiteaCert(ctx)
	if err != nil {
		return fmt.Errorf("failed to extract certificate: %v", err)
	}

	if err := extractor.SaveCertificate(cert, savePath); err != nil {
		return fmt.Errorf("failed to save certificate: %v", err)
	}

	return nil
}
