package api

import (
	"context"
	"time"
)

// SecurityManager handles image security operations including signing, verification,
// vulnerability scanning, and attestation management for OCI images.
type SecurityManager interface {
	// SignImage signs an image using the provided signer implementation.
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
}

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

// Signature contains metadata about a digital signature applied to an image.
type Signature struct {
	Algorithm              string           `json:"algorithm" yaml:"algorithm"`
	KeyID                  string           `json:"key_id" yaml:"key_id"`
	Signature              string           `json:"signature" yaml:"signature"`
	Timestamp              time.Time        `json:"timestamp" yaml:"timestamp"`
	Subject                string           `json:"subject,omitempty" yaml:"subject,omitempty"`
	Issuer                 string           `json:"issuer,omitempty" yaml:"issuer,omitempty"`
	CertificateFingerprint string           `json:"certificate_fingerprint,omitempty" yaml:"certificate_fingerprint,omitempty"`
	Bundle                 *SignatureBundle `json:"bundle,omitempty" yaml:"bundle,omitempty"`
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

// SBOM represents a Software Bill of Materials for an image.
type SBOM struct {
	Version      string                 `json:"version" yaml:"version"`
	Timestamp    time.Time              `json:"timestamp" yaml:"timestamp"`
	Image        string                 `json:"image" yaml:"image"`
	Components   []*Component           `json:"components" yaml:"components"`
	Dependencies map[string][]string    `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Tools        []*Tool                `json:"tools,omitempty" yaml:"tools,omitempty"`
	Metadata     map[string]string      `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Component represents a software component in an SBOM.
type Component struct {
	Name       string            `json:"name" yaml:"name"`
	Version    string            `json:"version" yaml:"version"`
	Type       string            `json:"type" yaml:"type"`
	Namespace  string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Licenses   []string          `json:"licenses,omitempty" yaml:"licenses,omitempty"`
	PackageURL string            `json:"package_url,omitempty" yaml:"package_url,omitempty"`
	Supplier   string            `json:"supplier,omitempty" yaml:"supplier,omitempty"`
	Hash       map[string]string `json:"hash,omitempty" yaml:"hash,omitempty"`
}

// Tool represents a tool used in SBOM generation.
type Tool struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Vendor  string `json:"vendor,omitempty" yaml:"vendor,omitempty"`
}

// VulnerabilityReport contains the results of a security vulnerability scan.
type VulnerabilityReport struct {
	Timestamp          time.Time               `json:"timestamp" yaml:"timestamp"`
	Image              string                  `json:"image" yaml:"image"`
	Scanner            *ScannerInfo            `json:"scanner" yaml:"scanner"`
	Summary            *VulnerabilitySummary   `json:"summary" yaml:"summary"`
	Vulnerabilities    []*Vulnerability        `json:"vulnerabilities" yaml:"vulnerabilities"`
	FixRecommendations []*FixRecommendation    `json:"fix_recommendations,omitempty" yaml:"fix_recommendations,omitempty"`
}

// ScannerInfo contains information about the vulnerability scanner.
type ScannerInfo struct {
	Name            string `json:"name" yaml:"name"`
	Version         string `json:"version" yaml:"version"`
	DatabaseVersion string `json:"database_version" yaml:"database_version"`
}

// VulnerabilitySummary provides aggregate vulnerability statistics.
type VulnerabilitySummary struct {
	Total    int `json:"total" yaml:"total"`
	Critical int `json:"critical" yaml:"critical"`
	High     int `json:"high" yaml:"high"`
	Medium   int `json:"medium" yaml:"medium"`
	Low      int `json:"low" yaml:"low"`
	Unknown  int `json:"unknown" yaml:"unknown"`
}

// Vulnerability represents a single security vulnerability.
type Vulnerability struct {
	ID           string   `json:"id" yaml:"id"`
	Severity     string   `json:"severity" yaml:"severity"`
	Score        float64  `json:"score,omitempty" yaml:"score,omitempty"`
	Title        string   `json:"title" yaml:"title"`
	Description  string   `json:"description,omitempty" yaml:"description,omitempty"`
	Package      string   `json:"package" yaml:"package"`
	Version      string   `json:"version" yaml:"version"`
	FixedVersion string   `json:"fixed_version,omitempty" yaml:"fixed_version,omitempty"`
	References   []string `json:"references,omitempty" yaml:"references,omitempty"`
}

// FixRecommendation suggests actions to remediate vulnerabilities.
type FixRecommendation struct {
	Type           string `json:"type" yaml:"type"`
	Target         string `json:"target" yaml:"target"`
	Recommendation string `json:"recommendation" yaml:"recommendation"`
	Priority       string `json:"priority" yaml:"priority"`
}

// Attestation represents a verifiable claim about an image.
type Attestation struct {
	Type      string                 `json:"type" yaml:"type"`
	Predicate map[string]interface{} `json:"predicate" yaml:"predicate"`
	Subject   string                 `json:"subject" yaml:"subject"`
	Timestamp time.Time              `json:"timestamp" yaml:"timestamp"`
	Issuer    string                 `json:"issuer" yaml:"issuer"`
	Signature *Signature             `json:"signature,omitempty" yaml:"signature,omitempty"`
}

// Policy defines security policies for images and operations.
type Policy struct {
	Version              string                   `json:"version" yaml:"version"`
	Name                 string                   `json:"name" yaml:"name"`
	Rules                []*PolicyRule            `json:"rules" yaml:"rules"`
	RequiredSignatures   []*SignatureRequirement `json:"required_signatures,omitempty" yaml:"required_signatures,omitempty"`
	MaxVulnerabilities   map[string]int           `json:"max_vulnerabilities,omitempty" yaml:"max_vulnerabilities,omitempty"`
	RequiredAttestations []string                 `json:"required_attestations,omitempty" yaml:"required_attestations,omitempty"`
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

// SecurityProfile contains comprehensive security information for an image.
type SecurityProfile struct {
	Image               string               `json:"image" yaml:"image"`
	Timestamp           time.Time            `json:"timestamp" yaml:"timestamp"`
	Signatures          []*Signature         `json:"signatures" yaml:"signatures"`
	Attestations        []*Attestation       `json:"attestations" yaml:"attestations"`
	VulnerabilityReport *VulnerabilityReport `json:"vulnerability_report,omitempty" yaml:"vulnerability_report,omitempty"`
	SBOM                *SBOM                `json:"sbom,omitempty" yaml:"sbom,omitempty"`
	SecurityScore       int                  `json:"security_score" yaml:"security_score"`
	Recommendations     []string             `json:"recommendations,omitempty" yaml:"recommendations,omitempty"`
}

// PolicyResult contains the outcome of policy enforcement.
type PolicyResult struct {
	Policy     string             `json:"policy" yaml:"policy"`
	Allowed    bool               `json:"allowed" yaml:"allowed"`
	Violations []*PolicyViolation `json:"violations,omitempty" yaml:"violations,omitempty"`
	Warnings   []string           `json:"warnings,omitempty" yaml:"warnings,omitempty"`
}

// PolicyViolation represents a security policy violation.
type PolicyViolation struct {
	Rule     string                 `json:"rule" yaml:"rule"`
	Severity string                 `json:"severity" yaml:"severity"`
	Message  string                 `json:"message" yaml:"message"`
	Details  map[string]interface{} `json:"details,omitempty" yaml:"details,omitempty"`
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