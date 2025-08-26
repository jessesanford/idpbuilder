package api

import (
	"context"
	"time"
)

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

// SecurityManager defines the main interface for security operations.
type SecurityManager interface {
	// SignImage signs an image using the provided signer.
	SignImage(ctx context.Context, image string, signer Signer) (*Signature, error)
	
	// VerifySignature verifies an image signature using the provided verifier.
	VerifySignature(ctx context.Context, image string, verifier Verifier) error
	
	// GenerateSBOM creates a Software Bill of Materials for the specified image.
	GenerateSBOM(ctx context.Context, image string) (*SBOM, error)
	
	// ScanVulnerabilities performs security vulnerability scanning on the image.
	ScanVulnerabilities(ctx context.Context, image string) (*VulnerabilityReport, error)
	
	// AttachAttestation attaches a security attestation to the image.
	AttachAttestation(ctx context.Context, image string, attestation *Attestation) error
	
	// VerifyAttestation verifies that an attached attestation is valid and trusted.
	VerifyAttestation(ctx context.Context, image string, policy *Policy) error
	
	// GetImageSecurityProfile returns comprehensive security information for an image.
	GetImageSecurityProfile(ctx context.Context, image string) (*SecurityProfile, error)
	
	// ValidatePolicy checks if a security policy is well-formed and applicable.
	ValidatePolicy(policy *Policy) error
	
	// EnforcePolicy applies security policy rules to an image or operation.
	EnforcePolicy(ctx context.Context, image string, policy *Policy) (*PolicyResult, error)
	
	// RotateKeys performs key rotation for all signers and updates verifiers.
	RotateKeys(ctx context.Context) error
	
	// GetTrustChain retrieves the certificate chain for a given key ID.
	GetTrustChain(keyID string) ([]*Certificate, error)
	
	// AddTrustedKey adds a key to the trust store for verification.
	AddTrustedKey(keyID string, certificate *Certificate) error
	
	// RemoveTrustedKey removes a key from the trust store.
	RemoveTrustedKey(keyID string) error
	
	// ValidateTrustChain validates that a certificate chain is properly formed and trusted.
	ValidateTrustChain(chain []*Certificate) error
}

// Signature represents a digital signature with metadata.
type Signature struct {
	Algorithm string           `json:"algorithm" yaml:"algorithm"`
	KeyID     string           `json:"key_id" yaml:"key_id"`
	Signature string           `json:"signature" yaml:"signature"`
	Timestamp time.Time        `json:"timestamp" yaml:"timestamp"`
	Subject   string           `json:"subject" yaml:"subject"`
	Bundle    *SignatureBundle `json:"bundle,omitempty" yaml:"bundle,omitempty"`
}

// SBOM represents a Software Bill of Materials.
type SBOM struct {
	Version    string               `json:"version" yaml:"version"`
	DataFormat string               `json:"data_format" yaml:"data_format"`
	Subject    string               `json:"subject" yaml:"subject"`
	Timestamp  time.Time            `json:"timestamp" yaml:"timestamp"`
	Components []SBOMComponent      `json:"components" yaml:"components"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// SBOMComponent represents a single component in an SBOM.
type SBOMComponent struct {
	Name      string   `json:"name" yaml:"name"`
	Version   string   `json:"version" yaml:"version"`
	Type      string   `json:"type" yaml:"type"`
	Licenses  []string `json:"licenses,omitempty" yaml:"licenses,omitempty"`
	Source    string   `json:"source,omitempty" yaml:"source,omitempty"`
}

// VulnerabilityReport contains vulnerability scan results.
type VulnerabilityReport struct {
	Timestamp       time.Time              `json:"timestamp" yaml:"timestamp"`
	Image           string                 `json:"image" yaml:"image"`
	Scanner         *ScannerInfo           `json:"scanner" yaml:"scanner"`
	Summary         *VulnerabilitySummary  `json:"summary" yaml:"summary"`
	Vulnerabilities []*Vulnerability       `json:"vulnerabilities" yaml:"vulnerabilities"`
}

// ScannerInfo provides information about the vulnerability scanner used.
type ScannerInfo struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	URL     string `json:"url,omitempty" yaml:"url,omitempty"`
}

// VulnerabilitySummary provides a summary of vulnerability counts.
type VulnerabilitySummary struct {
	Total    int `json:"total" yaml:"total"`
	Critical int `json:"critical" yaml:"critical"`
	High     int `json:"high" yaml:"high"`
	Medium   int `json:"medium" yaml:"medium"`
	Low      int `json:"low" yaml:"low"`
}

// Vulnerability represents a single security vulnerability.
type Vulnerability struct {
	ID          string                 `json:"id" yaml:"id"`
	Title       string                 `json:"title" yaml:"title"`
	Description string                 `json:"description" yaml:"description"`
	Severity    string                 `json:"severity" yaml:"severity"`
	CVSS        float64                `json:"cvss,omitempty" yaml:"cvss,omitempty"`
	Links       []string               `json:"links,omitempty" yaml:"links,omitempty"`
	Component   string                 `json:"component" yaml:"component"`
	Version     string                 `json:"version" yaml:"version"`
	FixedIn     string                 `json:"fixed_in,omitempty" yaml:"fixed_in,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Attestation represents a security attestation for an artifact.
type Attestation struct {
	Type        string                 `json:"type" yaml:"type"`
	Subject     string                 `json:"subject" yaml:"subject"`
	Predicate   map[string]interface{} `json:"predicate" yaml:"predicate"`
	Signature   *Signature             `json:"signature,omitempty" yaml:"signature,omitempty"`
	Timestamp   time.Time              `json:"timestamp" yaml:"timestamp"`
	Issuer      string                 `json:"issuer" yaml:"issuer"`
}

// SecurityProfile represents comprehensive security information for an artifact.
type SecurityProfile struct {
	Image               string               `json:"image" yaml:"image"`
	Timestamp           time.Time            `json:"timestamp" yaml:"timestamp"`
	Signatures          []*Signature         `json:"signatures" yaml:"signatures"`
	Attestations        []*Attestation       `json:"attestations" yaml:"attestations"`
	SBOM                *SBOM                `json:"sbom,omitempty" yaml:"sbom,omitempty"`
	VulnerabilityReport *VulnerabilityReport `json:"vulnerability_report,omitempty" yaml:"vulnerability_report,omitempty"`
	SecurityScore       int                  `json:"security_score" yaml:"security_score"`
}

// PolicyResult represents the result of policy enforcement.
type PolicyResult struct {
	Policy     string             `json:"policy" yaml:"policy"`
	Allowed    bool               `json:"allowed" yaml:"allowed"`
	Violations []*PolicyViolation `json:"violations" yaml:"violations"`
	Warnings   []string           `json:"warnings,omitempty" yaml:"warnings,omitempty"`
}

// PolicyViolation represents a single policy violation.
type PolicyViolation struct {
	Rule     string                 `json:"rule" yaml:"rule"`
	Severity string                 `json:"severity" yaml:"severity"`
	Message  string                 `json:"message" yaml:"message"`
	Details  map[string]interface{} `json:"details,omitempty" yaml:"details,omitempty"`
}