package security

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// cosignVerifier implements api.Verifier for signature verification
type cosignVerifier struct {
	publicKeys  map[string]crypto.PublicKey
	trustedCAs  *x509.CertPool
	trustedKeys []string
	policy      *api.Policy
}

// VerifierOption defines configuration options for verifiers
type VerifierOption func(*cosignVerifier)

// NewCosignVerifier creates a new signature verifier
func NewCosignVerifier(pubKeyPath string, opts ...VerifierOption) (api.Verifier, error) {
	verifier := &cosignVerifier{
		publicKeys:  make(map[string]crypto.PublicKey),
		trustedCAs:  x509.NewCertPool(),
		trustedKeys: []string{},
	}

	// Load the public key if path is provided
	if pubKeyPath != "" {
		err := verifier.loadPublicKey(pubKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load public key: %w", err)
		}
	}

	// Apply options
	for _, opt := range opts {
		opt(verifier)
	}

	return verifier, nil
}

// WithTrustedCA adds a trusted CA certificate for certificate chain validation
func WithTrustedCA(caPath string) VerifierOption {
	return func(v *cosignVerifier) {
		caData, err := os.ReadFile(caPath)
		if err != nil {
			return
		}

		v.trustedCAs.AppendCertsFromPEM(caData)
	}
}

// WithTrustedKey adds a trusted public key
func WithTrustedKey(keyID string, pubKey crypto.PublicKey) VerifierOption {
	return func(v *cosignVerifier) {
		v.publicKeys[keyID] = pubKey
		v.trustedKeys = append(v.trustedKeys, keyID)
	}
}

// WithPolicy sets the verification policy
func WithPolicy(policy *api.Policy) VerifierOption {
	return func(v *cosignVerifier) {
		v.policy = policy
	}
}

// Verify validates a signature against the provided payload
func (cv *cosignVerifier) Verify(payload []byte, signature []byte) error {
	if len(cv.publicKeys) == 0 {
		return fmt.Errorf("no trusted public keys configured")
	}

	hash := sha256.Sum256(payload)

	// Try each trusted public key
	for keyID, pubKey := range cv.publicKeys {
		err := cv.verifyWithKey(pubKey, hash[:], signature)
		if err == nil {
			// Verification succeeded with this key
			return nil
		}
	}

	return fmt.Errorf("signature verification failed with all trusted keys")
}

// verifyWithKey performs signature verification with a specific public key
func (cv *cosignVerifier) verifyWithKey(pubKey crypto.PublicKey, hash []byte, signature []byte) error {
	switch key := pubKey.(type) {
	case *rsa.PublicKey:
		return rsa.VerifyPKCS1v15(key, crypto.SHA256, hash, signature)
	case *ecdsa.PublicKey:
		return cv.verifyECDSA(key, hash, signature)
	case ed25519.PublicKey:
		if !ed25519.Verify(key, hash, signature) {
			return fmt.Errorf("ed25519 signature verification failed")
		}
		return nil
	default:
		return fmt.Errorf("unsupported public key type: %T", key)
	}
}

// verifyECDSA performs ECDSA signature verification
func (cv *cosignVerifier) verifyECDSA(pubKey *ecdsa.PublicKey, hash []byte, signature []byte) error {
	if !ecdsa.VerifyASN1(pubKey, hash, signature) {
		return fmt.Errorf("ECDSA signature verification failed")
	}
	return nil
}

// TrustedKeys returns the list of trusted key identifiers
func (cv *cosignVerifier) TrustedKeys() []string {
	return cv.trustedKeys
}

// VerifyPolicy checks if the verifier complies with a security policy
func (cv *cosignVerifier) VerifyPolicy(policy *api.Policy) error {
	if policy == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	// Check if required signatures are present
	for _, reqSig := range policy.RequiredSignatures {
		if reqSig.KeyID == "" {
			continue
		}

		if _, exists := cv.publicKeys[reqSig.KeyID]; !exists {
			return fmt.Errorf("required signature key %s not found in trusted keys", reqSig.KeyID)
		}
	}

	// Validate policy rules that affect verification
	for _, rule := range policy.Rules {
		if rule.Type == "signature_required" && len(cv.publicKeys) == 0 {
			return fmt.Errorf("policy requires signatures but no trusted keys configured")
		}
	}

	return nil
}

// GetTrustedRoots returns trusted root certificates or keys
func (cv *cosignVerifier) GetTrustedRoots() ([]*api.Certificate, error) {
	var certificates []*api.Certificate

	// Convert trusted CA certificates to API format
	for _, cert := range cv.trustedCAs.Subjects() {
		// Note: x509.CertPool doesn't provide direct access to certificates
		// In a real implementation, we'd maintain our own certificate store
		apiCert := &api.Certificate{
			Subject:     string(cert),
			Issuer:      "trusted-ca",
			Fingerprint: fmt.Sprintf("sha256:%x", sha256.Sum256(cert)),
			ValidFrom:   time.Now(),
			ValidTo:     time.Now().AddDate(1, 0, 0),
		}
		certificates = append(certificates, apiCert)
	}

	// Add public keys as certificates (for key-based verification)
	for keyID, pubKey := range cv.publicKeys {
		pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
		if err != nil {
			continue
		}

		apiCert := &api.Certificate{
			Data:        pubKeyBytes,
			Subject:     keyID,
			Issuer:      "self",
			Fingerprint: fmt.Sprintf("sha256:%x", sha256.Sum256(pubKeyBytes)),
			ValidFrom:   time.Now(),
			ValidTo:     time.Now().AddDate(10, 0, 0), // Long validity for public keys
		}
		certificates = append(certificates, apiCert)
	}

	return certificates, nil
}

// VerifyCertificateChain validates a certificate chain against trusted CAs
func (cv *cosignVerifier) VerifyCertificateChain(certs []*api.Certificate) error {
	if len(certs) == 0 {
		return fmt.Errorf("no certificates provided")
	}

	// Parse the leaf certificate
	leafCert, err := x509.ParseCertificate(certs[0].Data)
	if err != nil {
		return fmt.Errorf("failed to parse leaf certificate: %w", err)
	}

	// Build intermediate certificate pool
	intermediates := x509.NewCertPool()
	for i := 1; i < len(certs); i++ {
		cert, err := x509.ParseCertificate(certs[i].Data)
		if err != nil {
			continue
		}
		intermediates.AddCert(cert)
	}

	// Verify certificate chain
	opts := x509.VerifyOptions{
		Roots:         cv.trustedCAs,
		Intermediates: intermediates,
		CurrentTime:   time.Now(),
	}

	_, err = leafCert.Verify(opts)
	if err != nil {
		return fmt.Errorf("certificate chain verification failed: %w", err)
	}

	return nil
}

// ValidateSignatureBundle validates a complete signature bundle
func (cv *cosignVerifier) ValidateSignatureBundle(bundle *api.SignatureBundle, payload []byte) error {
	if bundle == nil {
		return fmt.Errorf("signature bundle cannot be nil")
	}

	if len(bundle.Signatures) == 0 {
		return fmt.Errorf("signature bundle must contain at least one signature")
	}

	// Verify certificate chain if present
	if len(bundle.Certificates) > 0 {
		if err := cv.VerifyCertificateChain(bundle.Certificates); err != nil {
			return fmt.Errorf("certificate chain validation failed: %w", err)
		}
	}

	// Verify each signature in the bundle
	for i, sigData := range bundle.Signatures {
		err := cv.Verify(payload, sigData.Signature)
		if err != nil {
			return fmt.Errorf("signature %d verification failed: %w", i, err)
		}
	}

	return nil
}

// loadPublicKey loads a public key from a PEM file
func (cv *cosignVerifier) loadPublicKey(pubKeyPath string) error {
	keyData, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read public key file %s: %w", pubKeyPath, err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block from key data")
	}

	var publicKey crypto.PublicKey

	switch block.Type {
	case "PUBLIC KEY":
		publicKey, err = x509.ParsePKIXPublicKey(block.Bytes)
	case "RSA PUBLIC KEY":
		publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	case "CERTIFICATE":
		cert, certErr := x509.ParseCertificate(block.Bytes)
		if certErr != nil {
			return fmt.Errorf("failed to parse certificate: %w", certErr)
		}
		publicKey = cert.PublicKey
	default:
		return fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	// Generate key ID and store the public key
	keyID := generatePublicKeyID(publicKey)
	cv.publicKeys[keyID] = publicKey
	cv.trustedKeys = append(cv.trustedKeys, keyID)

	return nil
}

// generatePublicKeyID creates a key identifier from a public key
func generatePublicKeyID(publicKey crypto.PublicKey) string {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "unknown"
	}

	hash := sha256.Sum256(publicKeyBytes)
	return fmt.Sprintf("sha256:%x", hash[:8]) // Use first 8 bytes as key ID
}

// Ensure we implement the interface
var _ api.Verifier = (*cosignVerifier)(nil)