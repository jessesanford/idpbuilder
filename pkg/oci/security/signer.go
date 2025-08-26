package security

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// cosignSigner implements api.Signer using Cosign-compatible signing
type cosignSigner struct {
	privateKey crypto.PrivateKey
	keyID      string
	algorithm  string
	certChain  []*api.Certificate
}

// keylessSigner implements api.Signer for keyless (OIDC-based) signing
type keylessSigner struct {
	provider  string
	identity  string
	algorithm string
	keyID     string
}

// SignerOption defines configuration options for signers
type SignerOption func(interface{})

// NewCosignSigner creates a new Cosign-compatible signer from a private key file
func NewCosignSigner(keyPath string, opts ...SignerOption) (api.Signer, error) {
	// Read the private key file
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file %s: %w", keyPath, err)
	}

	// Parse the PEM-encoded private key
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from key data")
	}

	var privateKey crypto.PrivateKey
	var algorithm string

	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		algorithm = "RS256"
	case "EC PRIVATE KEY":
		privateKey, err = x509.ParseECPrivateKey(block.Bytes)
		algorithm = "ES256"
	case "PRIVATE KEY":
		privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		// Determine algorithm based on key type
		switch privateKey.(type) {
		case *rsa.PrivateKey:
			algorithm = "RS256"
		case *ecdsa.PrivateKey:
			algorithm = "ES256"
		case ed25519.PrivateKey:
			algorithm = "EdDSA"
		default:
			return nil, fmt.Errorf("unsupported private key type")
		}
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	signer := &cosignSigner{
		privateKey: privateKey,
		keyID:      generateKeyID(privateKey),
		algorithm:  algorithm,
		certChain:  []*api.Certificate{},
	}

	// Apply options
	for _, opt := range opts {
		opt(signer)
	}

	return signer, nil
}

// NewKeylessSigner creates a keyless (OIDC-based) signer
func NewKeylessSigner(provider string) (api.Signer, error) {
	// In a real implementation, this would:
	// 1. Configure OIDC provider
	// 2. Get identity token
	// 3. Return keyless signer

	signer := &keylessSigner{
		provider:  provider,
		identity:  "user@example.com", // Would be from OIDC token
		algorithm: "RS256",
		keyID:     fmt.Sprintf("keyless:%s", provider),
	}

	return signer, nil
}

// WithCertificateChain adds certificate chain to the signer
func WithCertificateChain(certPath string) SignerOption {
	return func(s interface{}) {
		if signer, ok := s.(*cosignSigner); ok {
			certs, err := loadCertificateChain(certPath)
			if err == nil {
				signer.certChain = certs
			}
		}
	}
}

// WithKeyID sets a custom key ID for the signer
func WithKeyID(keyID string) SignerOption {
	return func(s interface{}) {
		if signer, ok := s.(*cosignSigner); ok {
			signer.keyID = keyID
		}
		if signer, ok := s.(*keylessSigner); ok {
			signer.keyID = keyID
		}
	}
}

// Sign creates a digital signature for the provided payload
func (cs *cosignSigner) Sign(payload []byte) ([]byte, error) {
	hash := sha256.Sum256(payload)

	switch key := cs.privateKey.(type) {
	case *rsa.PrivateKey:
		return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	case *ecdsa.PrivateKey:
		return ecdsa.SignASN1(rand.Reader, key, hash[:])
	case ed25519.PrivateKey:
		return ed25519.Sign(key, payload), nil
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", key)
	}
}

// KeyID returns the identifier for the signing key
func (cs *cosignSigner) KeyID() string {
	return cs.keyID
}

// Algorithm returns the signature algorithm used
func (cs *cosignSigner) Algorithm() string {
	return cs.algorithm
}

// PublicKey returns the public key for signature verification
func (cs *cosignSigner) PublicKey() ([]byte, error) {
	var publicKey crypto.PublicKey

	switch key := cs.privateKey.(type) {
	case *rsa.PrivateKey:
		publicKey = &key.PublicKey
	case *ecdsa.PrivateKey:
		publicKey = &key.PublicKey
	case ed25519.PrivateKey:
		publicKey = key.Public()
	default:
		return nil, fmt.Errorf("unsupported private key type: %T", key)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return publicKeyPEM, nil
}

// GetCertificateChain returns the certificate chain if applicable
func (cs *cosignSigner) GetCertificateChain() ([]*api.Certificate, error) {
	return cs.certChain, nil
}

// Keyless signer implementations
func (ks *keylessSigner) Sign(payload []byte) ([]byte, error) {
	// In a real implementation, this would:
	// 1. Get ephemeral signing key from OIDC provider
	// 2. Sign the payload
	// 3. Return signature and certificate

	// For now, return a mock signature
	hash := sha256.Sum256(payload)
	return []byte(base64.StdEncoding.EncodeToString(hash[:])), nil
}

func (ks *keylessSigner) KeyID() string {
	return ks.keyID
}

func (ks *keylessSigner) Algorithm() string {
	return ks.algorithm
}

func (ks *keylessSigner) PublicKey() ([]byte, error) {
	// In keyless signing, the public key comes from the certificate
	// For now, return a placeholder
	return []byte("keyless-public-key-from-certificate"), nil
}

func (ks *keylessSigner) GetCertificateChain() ([]*api.Certificate, error) {
	// In keyless signing, certificate is issued by the provider
	return []*api.Certificate{}, nil
}

// generateKeyID creates a key identifier from a private key
func generateKeyID(privateKey crypto.PrivateKey) string {
	var publicKey crypto.PublicKey

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		publicKey = &key.PublicKey
	case *ecdsa.PrivateKey:
		publicKey = &key.PublicKey
	case ed25519.PrivateKey:
		publicKey = key.Public()
	default:
		return "unknown"
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "unknown"
	}

	hash := sha256.Sum256(publicKeyBytes)
	return fmt.Sprintf("sha256:%x", hash[:8]) // Use first 8 bytes as key ID
}

// loadCertificateChain loads a certificate chain from a PEM file
func loadCertificateChain(certPath string) ([]*api.Certificate, error) {
	certData, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file %s: %w", certPath, err)
	}

	var certificates []*api.Certificate
	block, rest := pem.Decode(certData)

	for block != nil {
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate: %w", err)
			}

			apiCert := &api.Certificate{
				Data:        block.Bytes,
				Subject:     cert.Subject.String(),
				Issuer:      cert.Issuer.String(),
				Fingerprint: fmt.Sprintf("sha256:%x", sha256.Sum256(block.Bytes)),
				ValidFrom:   cert.NotBefore,
				ValidTo:     cert.NotAfter,
			}

			certificates = append(certificates, apiCert)
		}

		block, rest = pem.Decode(rest)
	}

	if len(certificates) == 0 {
		return nil, fmt.Errorf("no certificates found in file %s", certPath)
	}

	return certificates, nil
}

// GenerateTestKey generates a test RSA key pair for testing purposes
func GenerateTestKey() (crypto.PrivateKey, crypto.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// WriteKeyToPEM writes a private key to PEM format
func WriteKeyToPEM(key crypto.PrivateKey, writer io.Writer) error {
	var keyBytes []byte
	var keyType string
	var err error

	switch k := key.(type) {
	case *rsa.PrivateKey:
		keyBytes = x509.MarshalPKCS1PrivateKey(k)
		keyType = "RSA PRIVATE KEY"
	case *ecdsa.PrivateKey:
		keyBytes, err = x509.MarshalECPrivateKey(k)
		keyType = "EC PRIVATE KEY"
	case ed25519.PrivateKey:
		keyBytes, err = x509.MarshalPKCS8PrivateKey(k)
		keyType = "PRIVATE KEY"
	default:
		return fmt.Errorf("unsupported key type: %T", k)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	return pem.Encode(writer, &pem.Block{
		Type:  keyType,
		Bytes: keyBytes,
	})
}

// Ensure we implement the interface
var _ api.Signer = (*cosignSigner)(nil)
var _ api.Signer = (*keylessSigner)(nil)