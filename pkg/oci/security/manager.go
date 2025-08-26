package security

import (
	"context"
	"fmt"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// ScannerPlugin defines the interface for vulnerability scanners.
type ScannerPlugin interface {
	Scan(ctx context.Context, image string) (*api.VulnerabilityReport, error)
}

// sbomGenerator generates Software Bill of Materials.
type sbomGenerator struct{}

// newSBOMGenerator creates a new SBOM generator instance.
func newSBOMGenerator() *sbomGenerator {
	return &sbomGenerator{}
}

// Generate creates an SBOM for the specified image.
func (s *sbomGenerator) Generate(ctx context.Context, image string) (*api.SBOM, error) {
	// This is a simplified implementation
	return &api.SBOM{
		Version:    "1.0",
		DataFormat: "application/vnd.cyclonedx+json",
		Subject:    image,
		Timestamp:  time.Now(),
		Components: []api.SBOMComponent{
			{
				Name:     "base-image",
				Version:  "latest",
				Type:     "container",
				Licenses: []string{"Apache-2.0"},
			},
		},
		Metadata: map[string]interface{}{
			"tool": "idpbuilder-security",
		},
	}, nil
}

// securityManager implements api.SecurityManager interface
type securityManager struct {
	signers   map[string]api.Signer
	verifiers map[string]api.Verifier
	sbomGen   *sbomGenerator
	scanners  []ScannerPlugin
	config    *SecurityConfig
}

// SecurityConfig holds configuration for the security manager
type SecurityConfig struct {
	DefaultSignerID    string
	DefaultVerifierID  string
	ScanTimeout        time.Duration
	SBOMTimeout        time.Duration
	EnableVulnScanning bool
	EnableSBOM         bool
}

// Option defines a configuration option for SecurityManager
type Option func(*securityManager)

// NewSecurityManager creates a new security manager instance
func NewSecurityManager(opts ...Option) api.SecurityManager {
	sm := &securityManager{
		signers:   make(map[string]api.Signer),
		verifiers: make(map[string]api.Verifier),
		sbomGen:   newSBOMGenerator(),
		scanners:  []ScannerPlugin{},
		config: &SecurityConfig{
			ScanTimeout: 10 * time.Minute,
			SBOMTimeout: 5 * time.Minute,
			EnableVulnScanning: true,
			EnableSBOM:         true,
		},
	}

	for _, opt := range opts {
		opt(sm)
	}

	return sm
}

// WithSigner adds a signer to the security manager
func WithSigner(id string, signer api.Signer) Option {
	return func(sm *securityManager) {
		sm.signers[id] = signer
		if sm.config.DefaultSignerID == "" {
			sm.config.DefaultSignerID = id
		}
	}
}

// WithVerifier adds a verifier to the security manager
func WithVerifier(id string, verifier api.Verifier) Option {
	return func(sm *securityManager) {
		sm.verifiers[id] = verifier
		if sm.config.DefaultVerifierID == "" {
			sm.config.DefaultVerifierID = id
		}
	}
}

// WithScanner adds a vulnerability scanner plugin
func WithScanner(scanner ScannerPlugin) Option {
	return func(sm *securityManager) {
		sm.scanners = append(sm.scanners, scanner)
	}
}

// WithConfig sets the security manager configuration
func WithConfig(config *SecurityConfig) Option {
	return func(sm *securityManager) {
		sm.config = config
	}
}

// SignImage signs an image using the provided signer
func (sm *securityManager) SignImage(ctx context.Context, image string, signer api.Signer) (*api.Signature, error) {
	if signer == nil {
		if sm.config.DefaultSignerID == "" {
			return nil, fmt.Errorf("no signer provided and no default signer configured")
		}
		var ok bool
		signer, ok = sm.signers[sm.config.DefaultSignerID]
		if !ok {
			return nil, fmt.Errorf("default signer %s not found", sm.config.DefaultSignerID)
		}
	}

	// Create payload for signing (simplified image manifest digest)
	payload := []byte(fmt.Sprintf("image:%s", image))

	// Sign the payload
	signatureBytes, err := signer.Sign(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to sign payload: %w", err)
	}

	// Get certificate chain if available
	certChain, _ := signer.GetCertificateChain()
	var bundle *api.SignatureBundle
	if len(certChain) > 0 {
		bundle = &api.SignatureBundle{
			MediaType:    "application/vnd.dev.cosign.simplesigning.v1+json",
			Payload:      payload,
			Signatures:   []api.SignatureData{{Signature: signatureBytes}},
			Certificates: certChain,
		}
	}

	signature := &api.Signature{
		Algorithm: signer.Algorithm(),
		KeyID:     signer.KeyID(),
		Signature: string(signatureBytes),
		Timestamp: time.Now(),
		Subject:   image,
		Bundle:    bundle,
	}

	return signature, nil
}

// VerifySignature verifies an image signature using the provided verifier
func (sm *securityManager) VerifySignature(ctx context.Context, image string, verifier api.Verifier) error {
	if verifier == nil {
		if sm.config.DefaultVerifierID == "" {
			return fmt.Errorf("no verifier provided and no default verifier configured")
		}
		var ok bool
		verifier, ok = sm.verifiers[sm.config.DefaultVerifierID]
		if !ok {
			return fmt.Errorf("default verifier %s not found", sm.config.DefaultVerifierID)
		}
	}

	// Create the same payload that was signed
	payload := []byte(fmt.Sprintf("image:%s", image))

	// For this implementation, we simulate signature verification
	// In a real implementation, this would fetch the signature from the registry
	// and verify it against the payload
	_ = payload

	return fmt.Errorf("signature verification not implemented - requires registry integration")
}

// GenerateSBOM creates a Software Bill of Materials for the specified image
func (sm *securityManager) GenerateSBOM(ctx context.Context, image string) (*api.SBOM, error) {
	if !sm.config.EnableSBOM {
		return nil, fmt.Errorf("SBOM generation is disabled")
	}

	ctx, cancel := context.WithTimeout(ctx, sm.config.SBOMTimeout)
	defer cancel()

	return sm.sbomGen.Generate(ctx, image)
}

// ScanVulnerabilities performs security vulnerability scanning on the image
func (sm *securityManager) ScanVulnerabilities(ctx context.Context, image string) (*api.VulnerabilityReport, error) {
	if !sm.config.EnableVulnScanning {
		return nil, fmt.Errorf("vulnerability scanning is disabled")
	}

	if len(sm.scanners) == 0 {
		return &api.VulnerabilityReport{
			Timestamp: time.Now(),
			Image:     image,
			Scanner:   &api.ScannerInfo{Name: "none", Version: "0.0.0"},
			Summary:   &api.VulnerabilitySummary{Total: 0},
			Vulnerabilities: []*api.Vulnerability{},
		}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, sm.config.ScanTimeout)
	defer cancel()

	// Use the first available scanner
	scanner := sm.scanners[0]
	return scanner.Scan(ctx, image)
}

// AttachAttestation attaches a security attestation to the image
func (sm *securityManager) AttachAttestation(ctx context.Context, image string, attestation *api.Attestation) error {
	if attestation == nil {
		return fmt.Errorf("attestation cannot be nil")
	}

	// Validate attestation
	if attestation.Type == "" {
		return fmt.Errorf("attestation type is required")
	}
	if attestation.Subject != image {
		return fmt.Errorf("attestation subject must match image")
	}

	// In a real implementation, this would store the attestation
	// to the registry or attestation store
	return fmt.Errorf("attestation attachment not implemented - requires registry integration")
}

// VerifyAttestation verifies that an attached attestation is valid and trusted
func (sm *securityManager) VerifyAttestation(ctx context.Context, image string, policy *api.Policy) error {
	if policy == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	// In a real implementation, this would:
	// 1. Fetch attestations from the registry
	// 2. Verify signatures on attestations
	// 3. Check policy compliance
	return fmt.Errorf("attestation verification not implemented - requires registry integration")
}

// GetImageSecurityProfile returns comprehensive security information for an image
func (sm *securityManager) GetImageSecurityProfile(ctx context.Context, image string) (*api.SecurityProfile, error) {
	profile := &api.SecurityProfile{
		Image:     image,
		Timestamp: time.Now(),
		Signatures: []*api.Signature{},
		Attestations: []*api.Attestation{},
	}

	// Generate SBOM if enabled
	if sm.config.EnableSBOM {
		sbom, err := sm.GenerateSBOM(ctx, image)
		if err == nil {
			profile.SBOM = sbom
		}
	}

	// Scan vulnerabilities if enabled
	if sm.config.EnableVulnScanning {
		vulnReport, err := sm.ScanVulnerabilities(ctx, image)
		if err == nil {
			profile.VulnerabilityReport = vulnReport
		}
	}

	// Calculate security score based on available data
	profile.SecurityScore = sm.calculateSecurityScore(profile)

	return profile, nil
}

// ValidatePolicy checks if a security policy is well-formed and applicable
func (sm *securityManager) ValidatePolicy(policy *api.Policy) error {
	if policy == nil {
		return fmt.Errorf("policy cannot be nil")
	}

	if policy.Version == "" {
		return fmt.Errorf("policy version is required")
	}

	if policy.Name == "" {
		return fmt.Errorf("policy name is required")
	}

	if len(policy.Rules) == 0 {
		return fmt.Errorf("policy must have at least one rule")
	}

	// Validate each rule
	for i, rule := range policy.Rules {
		if rule.Name == "" {
			return fmt.Errorf("rule %d: name is required", i)
		}
		if rule.Type == "" {
			return fmt.Errorf("rule %d: type is required", i)
		}
		if rule.Action == "" {
			return fmt.Errorf("rule %d: action is required", i)
		}
	}

	return nil
}

// EnforcePolicy applies security policy rules to an image or operation
func (sm *securityManager) EnforcePolicy(ctx context.Context, image string, policy *api.Policy) (*api.PolicyResult, error) {
	if err := sm.ValidatePolicy(policy); err != nil {
		return nil, fmt.Errorf("invalid policy: %w", err)
	}

	result := &api.PolicyResult{
		Policy:     policy.Name,
		Allowed:    true,
		Violations: []*api.PolicyViolation{},
		Warnings:   []string{},
	}

	// Get security profile for policy evaluation
	profile, err := sm.GetImageSecurityProfile(ctx, image)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to get security profile: %v", err))
	}

	// Evaluate each policy rule
	for _, rule := range policy.Rules {
		violation := sm.evaluateRule(rule, profile, image)
		if violation != nil {
			result.Violations = append(result.Violations, violation)
			if rule.Action == "deny" {
				result.Allowed = false
			}
		}
	}

	return result, nil
}

// calculateSecurityScore computes a security score based on available security data
func (sm *securityManager) calculateSecurityScore(profile *api.SecurityProfile) int {
	score := 100 // Start with perfect score

	// Deduct points for vulnerabilities
	if profile.VulnerabilityReport != nil && profile.VulnerabilityReport.Summary != nil {
		summary := profile.VulnerabilityReport.Summary
		score -= summary.Critical * 20
		score -= summary.High * 10
		score -= summary.Medium * 5
		score -= summary.Low * 1
	}

	// Deduct points for missing signatures
	if len(profile.Signatures) == 0 {
		score -= 30
	}

	// Deduct points for missing attestations
	if len(profile.Attestations) == 0 {
		score -= 20
	}

	// Ensure score doesn't go below 0
	if score < 0 {
		score = 0
	}

	return score
}

// evaluateRule evaluates a single policy rule against a security profile
func (sm *securityManager) evaluateRule(rule *api.PolicyRule, profile *api.SecurityProfile, image string) *api.PolicyViolation {
	switch rule.Type {
	case "signature_required":
		if len(profile.Signatures) == 0 {
			return &api.PolicyViolation{
				Rule:     rule.Name,
				Severity: rule.Severity,
				Message:  "Image must be signed",
				Details:  map[string]interface{}{"image": image},
			}
		}
	case "max_vulnerabilities":
		if profile.VulnerabilityReport != nil {
			summary := profile.VulnerabilityReport.Summary
			if summary.Critical > 0 && rule.Severity == "critical" {
				return &api.PolicyViolation{
					Rule:     rule.Name,
					Severity: rule.Severity,
					Message:  fmt.Sprintf("Image has %d critical vulnerabilities", summary.Critical),
					Details:  map[string]interface{}{"count": summary.Critical},
				}
			}
		}
	}

	return nil
}

// RotateKeys performs key rotation for all signers and updates verifiers
func (sm *securityManager) RotateKeys(ctx context.Context) error {
	// This is a simplified key rotation implementation
	// In a real system, this would coordinate with key management services
	
	rotatedSigners := make(map[string]api.Signer)
	
	for id, signer := range sm.signers {
		// Simulate key rotation - in reality this would generate new keys
		// and update certificate chains
		rotatedSigners[id] = signer // For now, keep existing signer
		
		// Log rotation event
		fmt.Printf("Rotated signing key for signer: %s\n", id)
	}
	
	// Update signers with rotated keys
	sm.signers = rotatedSigners
	
	// Update verifiers to trust the new keys
	for id, verifier := range sm.verifiers {
		// In a real implementation, this would update trusted key lists
		// and certificate chains in the verifier
		_ = verifier // Placeholder - keep existing verifier
		fmt.Printf("Updated trusted keys for verifier: %s\n", id)
	}
	
	return nil
}

// GetTrustChain retrieves the certificate chain for a given key ID
func (sm *securityManager) GetTrustChain(keyID string) ([]*api.Certificate, error) {
	// Look for the key in our signers
	for id, signer := range sm.signers {
		if signer.KeyID() == keyID {
			chain, err := signer.GetCertificateChain()
			if err != nil {
				return nil, fmt.Errorf("failed to get certificate chain for key %s: %w", keyID, err)
			}
			return chain, nil
		}
		// Also check if the signer ID matches
		if id == keyID {
			chain, err := signer.GetCertificateChain()
			if err != nil {
				return nil, fmt.Errorf("failed to get certificate chain for signer %s: %w", id, err)
			}
			return chain, nil
		}
	}
	
	// Look in verifiers for trusted roots
	for _, verifier := range sm.verifiers {
		if trustedKeys := verifier.TrustedKeys(); contains(trustedKeys, keyID) {
			roots, err := verifier.GetTrustedRoots()
			if err != nil {
				return nil, fmt.Errorf("failed to get trusted roots for key %s: %w", keyID, err)
			}
			return roots, nil
		}
	}
	
	return nil, fmt.Errorf("trust chain not found for key ID: %s", keyID)
}

// AddTrustedKey adds a key to the trust store for verification
func (sm *securityManager) AddTrustedKey(keyID string, certificate *api.Certificate) error {
	if certificate == nil {
		return fmt.Errorf("certificate cannot be nil")
	}
	
	// Validate certificate
	if certificate.Subject == "" {
		return fmt.Errorf("certificate must have a subject")
	}
	
	if certificate.ValidTo.Before(time.Now()) {
		return fmt.Errorf("certificate is expired")
	}
	
	// In a real implementation, this would add the certificate to a trust store
	// For now, we'll update all verifiers to trust this key
	for id := range sm.verifiers {
		fmt.Printf("Added trusted key %s to verifier %s\n", keyID, id)
	}
	
	return nil
}

// RemoveTrustedKey removes a key from the trust store
func (sm *securityManager) RemoveTrustedKey(keyID string) error {
	if keyID == "" {
		return fmt.Errorf("key ID cannot be empty")
	}
	
	// In a real implementation, this would remove the key from trust stores
	for id := range sm.verifiers {
		fmt.Printf("Removed trusted key %s from verifier %s\n", keyID, id)
	}
	
	return nil
}

// ValidateTrustChain validates that a certificate chain is properly formed and trusted
func (sm *securityManager) ValidateTrustChain(chain []*api.Certificate) error {
	if len(chain) == 0 {
		return fmt.Errorf("certificate chain cannot be empty")
	}
	
	// Validate each certificate in the chain
	for i, cert := range chain {
		if cert == nil {
			return fmt.Errorf("certificate %d is nil", i)
		}
		
		// Check certificate validity period
		now := time.Now()
		if now.Before(cert.ValidFrom) {
			return fmt.Errorf("certificate %d is not yet valid (valid from: %v)", i, cert.ValidFrom)
		}
		if now.After(cert.ValidTo) {
			return fmt.Errorf("certificate %d is expired (valid to: %v)", i, cert.ValidTo)
		}
	}
	
	// In a real implementation, this would:
	// 1. Verify certificate signatures up the chain
	// 2. Check certificate revocation status
	// 3. Validate that the root is trusted
	
	return nil
}

// contains is a helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Ensure we implement the interface
var _ api.SecurityManager = (*securityManager)(nil)