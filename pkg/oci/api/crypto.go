package api

import "time"

// Signer defines the interface for digital signing operations.
type Signer interface {
	// Sign creates a digital signature for the provided payload.
	Sign(payload []byte) ([]byte, error)

	// KeyID returns the identifier for the signing key.
	KeyID() string

	// Algorithm returns the signature algorithm used (e.g., "RS256", "ES384").
	Algorithm() string

	// PublicKey returns the public key for signature verification.
	PublicKey() ([]byte, error)

	// GetCertificateChain returns the certificate chain if applicable.
	GetCertificateChain() ([]*Certificate, error)
}

// Verifier defines the interface for signature verification operations.
type Verifier interface {
	// Verify validates a signature against the provided payload.
	Verify(payload []byte, signature []byte) error

	// TrustedKeys returns the list of trusted key identifiers.
	TrustedKeys() []string

	// VerifyPolicy checks if the verifier complies with a security policy.
	VerifyPolicy(policy *Policy) error

	// GetTrustedRoots returns trusted root certificates or keys.
	GetTrustedRoots() ([]*Certificate, error)
}

// Certificate represents an X.509 certificate for cryptographic operations.
type Certificate struct {
	Data        []byte    `json:"data" yaml:"data"`
	Subject     string    `json:"subject" yaml:"subject"`
	Issuer      string    `json:"issuer" yaml:"issuer"`
	Fingerprint string    `json:"fingerprint" yaml:"fingerprint"`
	ValidFrom   time.Time `json:"valid_from" yaml:"valid_from"`
	ValidTo     time.Time `json:"valid_to" yaml:"valid_to"`
}

// Policy defines security policies for signature verification.
type Policy struct {
	Version              string                   `json:"version" yaml:"version"`
	Name                 string                   `json:"name" yaml:"name"`
	Rules                []*PolicyRule            `json:"rules" yaml:"rules"`
	RequiredSignatures   []*SignatureRequirement `json:"required_signatures,omitempty" yaml:"required_signatures,omitempty"`
}

// PolicyRule represents a single policy rule.
type PolicyRule struct {
	Name      string `json:"name" yaml:"name"`
	Type      string `json:"type" yaml:"type"`
	Condition string `json:"condition" yaml:"condition"`
	Action    string `json:"action" yaml:"action"`
	Severity  string `json:"severity" yaml:"severity"`
}

// SignatureRequirement specifies signature verification requirements.
type SignatureRequirement struct {
	KeyID   string `json:"key_id" yaml:"key_id"`
	Issuer  string `json:"issuer,omitempty" yaml:"issuer,omitempty"`
	Subject string `json:"subject,omitempty" yaml:"subject,omitempty"`
}

// SignatureBundle contains all data needed for signature verification.
type SignatureBundle struct {
	MediaType    string            `json:"media_type" yaml:"media_type"`
	Payload      []byte            `json:"payload" yaml:"payload"`
	Signatures   []SignatureData   `json:"signatures" yaml:"signatures"`
	Certificates []*Certificate    `json:"certificates,omitempty" yaml:"certificates,omitempty"`
}

// SignatureData represents a single signature within a bundle.
type SignatureData struct {
	Signature []byte                 `json:"signature" yaml:"signature"`
	Protected string                 `json:"protected" yaml:"protected"`
	Header    map[string]interface{} `json:"header,omitempty" yaml:"header,omitempty"`
}