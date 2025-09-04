package testdata

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"
)

// GenerateTestCert creates a test certificate with the given parameters
func GenerateTestCert(template *x509.Certificate, parent *x509.Certificate, publicKey, privateKey interface{}) (*x509.Certificate, error) {
	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	
	return x509.ParseCertificate(certDER)
}

// GenerateValidCert creates a valid test certificate for gitea.local
func GenerateValidCert() (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:    "gitea.local",
			Organization:  []string{"Test Org"},
			Country:       []string{"US"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year
		DNSNames:              []string{"gitea.local", "registry.local"},
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	
	return GenerateTestCert(template, template, &privateKey.PublicKey, privateKey)
}

// GenerateExpiredCert creates an expired test certificate
func GenerateExpiredCert() (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	
	template := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: "expired.local",
		},
		NotBefore: time.Now().Add(-365 * 24 * time.Hour), // Started 1 year ago
		NotAfter:  time.Now().Add(-24 * time.Hour),       // Expired 1 day ago
		DNSNames:  []string{"expired.local"},
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	
	return GenerateTestCert(template, template, &privateKey.PublicKey, privateKey)
}

// GenerateExpiringSoonCert creates a certificate expiring in 15 days
func GenerateExpiringSoonCert() (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	
	template := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName: "expiring.local",
		},
		NotBefore: time.Now().Add(-30 * 24 * time.Hour), // Started 30 days ago
		NotAfter:  time.Now().Add(15 * 24 * time.Hour),  // Expires in 15 days
		DNSNames:  []string{"expiring.local"},
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	
	return GenerateTestCert(template, template, &privateKey.PublicKey, privateKey)
}

// GenerateWildcardCert creates a wildcard certificate
func GenerateWildcardCert() (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	
	template := &x509.Certificate{
		SerialNumber: big.NewInt(4),
		Subject: pkix.Name{
			CommonName: "*.example.local",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		DNSNames:  []string{"*.example.local", "example.local"},
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	
	return GenerateTestCert(template, template, &privateKey.PublicKey, privateKey)
}

// GenerateCertWithIPSAN creates a certificate with IP address SAN
func GenerateCertWithIPSAN() (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	
	template := &x509.Certificate{
		SerialNumber: big.NewInt(5),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),
		DNSNames:  []string{"localhost"},
		IPAddresses: []net.IP{
			net.IPv4(127, 0, 0, 1),
			net.ParseIP("::1"),
		},
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	
	return GenerateTestCert(template, template, &privateKey.PublicKey, privateKey)
}

// GenerateNotYetValidCert creates a certificate that is not yet valid
func GenerateNotYetValidCert() (*x509.Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	
	template := &x509.Certificate{
		SerialNumber: big.NewInt(6),
		Subject: pkix.Name{
			CommonName: "future.local",
		},
		NotBefore: time.Now().Add(24 * time.Hour),        // Valid starting tomorrow
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),  // Valid for 1 year from tomorrow
		DNSNames:  []string{"future.local"},
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	
	return GenerateTestCert(template, template, &privateKey.PublicKey, privateKey)
}