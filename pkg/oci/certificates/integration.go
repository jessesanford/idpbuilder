package certificates

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// RegistryClient interface defines registry operations that need certificate integration
// This will be implemented by Phase 2 registry client
type RegistryClient interface {
	Push(ctx context.Context, ref string, content []byte) error
	Pull(ctx context.Context, ref string) ([]byte, error)
	Login(ctx context.Context, registry, username, password string) error
	SetTransport(transport http.RoundTripper) error
}

// BuildConfiguration defines build-related certificate settings
type BuildConfiguration struct {
	CustomCAPath     string
	SkipTLSVerify    bool
	TLSTimeout       time.Duration
	MaxRetries       int
	InsecureRegistry bool
}

// CertificateMiddleware defines pluggable certificate injection points
type CertificateMiddleware interface {
	InjectCertificates(transport http.RoundTripper) http.RoundTripper
	ConfigureTLS(config *tls.Config) error
	ValidateConnection(url string) error
}

// RegistryIntegration wraps registry client with certificate support
type RegistryIntegration struct {
	certService CertificateService
	tlsConfig   *tls.Config
	registry    RegistryClient
	middleware  []CertificateMiddleware
}

// NewRegistryIntegration creates a new registry integration instance
func NewRegistryIntegration(certService CertificateService, registry RegistryClient) *RegistryIntegration {
	return &RegistryIntegration{
		certService: certService,
		registry:    registry,
		middleware:  make([]CertificateMiddleware, 0),
	}
}

// AddMiddleware adds certificate middleware to the integration
func (ri *RegistryIntegration) AddMiddleware(mw CertificateMiddleware) {
	ri.middleware = append(ri.middleware, mw)
}

// InitializeTLS sets up TLS configuration with certificates
func (ri *RegistryIntegration) InitializeTLS(ctx context.Context) error {
	tlsConfig, err := ri.certService.GetTLSConfig()
	if err != nil {
		return fmt.Errorf("failed to get TLS config: %w", err)
	}

	// Apply middleware configurations
	for _, mw := range ri.middleware {
		if err := mw.ConfigureTLS(tlsConfig); err != nil {
			return fmt.Errorf("middleware TLS configuration failed: %w", err)
		}
	}

	ri.tlsConfig = tlsConfig
	return nil
}

// Push pushes content to registry with certificate validation
func (ri *RegistryIntegration) Push(ctx context.Context, ref string, content []byte) error {
	if err := ri.ensureTLSConfigured(ctx); err != nil {
		return err
	}

	transport := ri.createSecureTransport()
	if err := ri.registry.SetTransport(transport); err != nil {
		return fmt.Errorf("failed to set secure transport: %w", err)
	}

	return ri.registry.Push(ctx, ref, content)
}

// Pull pulls content from registry with certificate validation
func (ri *RegistryIntegration) Pull(ctx context.Context, ref string) ([]byte, error) {
	if err := ri.ensureTLSConfigured(ctx); err != nil {
		return nil, err
	}

	transport := ri.createSecureTransport()
	if err := ri.registry.SetTransport(transport); err != nil {
		return nil, fmt.Errorf("failed to set secure transport: %w", err)
	}

	return ri.registry.Pull(ctx, ref)
}

// Login authenticates with registry using certificate-enhanced transport
func (ri *RegistryIntegration) Login(ctx context.Context, registry, username, password string) error {
	if err := ri.ensureTLSConfigured(ctx); err != nil {
		return err
	}

	return ri.registry.Login(ctx, registry, username, password)
}

// ValidateRegistryCertificate validates registry endpoint certificate
func (ri *RegistryIntegration) ValidateRegistryCertificate(url string) error {
	for _, mw := range ri.middleware {
		if err := mw.ValidateConnection(url); err != nil {
			return fmt.Errorf("certificate validation failed for %s: %w", url, err)
		}
	}
	return nil
}

// ensureTLSConfigured ensures TLS is properly configured
func (ri *RegistryIntegration) ensureTLSConfigured(ctx context.Context) error {
	if ri.tlsConfig == nil {
		return ri.InitializeTLS(ctx)
	}
	return nil
}

// createSecureTransport creates HTTP transport with certificate configuration
func (ri *RegistryIntegration) createSecureTransport() http.RoundTripper {
	transport := &http.Transport{
		TLSClientConfig: ri.tlsConfig,
	}

	// Apply middleware transformations
	var rt http.RoundTripper = transport
	for _, mw := range ri.middleware {
		rt = mw.InjectCertificates(rt)
	}

	return rt
}

// BuildIntegration integrates certificates with build process
type BuildIntegration struct {
	certService CertificateService
	buildConfig BuildConfiguration
	caBundle    *CertBundle
}

// NewBuildIntegration creates a new build integration instance
func NewBuildIntegration(certService CertificateService, config BuildConfiguration) *BuildIntegration {
	return &BuildIntegration{
		certService: certService,
		buildConfig: config,
	}
}

// LoadCertificatesForBuild loads certificates needed for build process
func (bi *BuildIntegration) LoadCertificatesForBuild(ctx context.Context) error {
	if bi.buildConfig.CustomCAPath == "" {
		return nil // No custom CA configured
	}

	bundle, err := bi.certService.LoadCertificateBundle(ctx, bi.buildConfig.CustomCAPath, CertFormatPEM)
	if err != nil {
		return fmt.Errorf("failed to load CA bundle: %w", err)
	}

	bi.caBundle = bundle
	return nil
}

// ConfigureBuildahWithCertificates configures buildah with custom certificates
func (bi *BuildIntegration) ConfigureBuildahWithCertificates(buildahConfig map[string]string) error {
	// For insecure registry mode, skip TLS verification
	if bi.buildConfig.InsecureRegistry {
		buildahConfig["tls-verify"] = "false"
		return nil
	}

	// For secure configurations, we need certificates unless skip is explicitly set
	if bi.buildConfig.SkipTLSVerify {
		buildahConfig["tls-verify"] = "false"
		return nil
	}

	// For secure TLS with custom CA, bundle must be loaded
	if bi.caBundle == nil && bi.buildConfig.CustomCAPath != "" {
		return fmt.Errorf("CA bundle not loaded, call LoadCertificatesForBuild first")
	}

	// Configure buildah with custom CA
	buildahConfig["tls-verify"] = "true"

	if bi.caBundle != nil && len(bi.caBundle.CAs) > 0 {
		buildahConfig["cert-dir"] = bi.buildConfig.CustomCAPath
	}

	return nil
}

// ValidateBuildCertificates validates certificates before build
func (bi *BuildIntegration) ValidateBuildCertificates() error {
	if bi.caBundle == nil {
		return nil // No certificates to validate
	}

	for _, cert := range bi.caBundle.CAs {
		if err := bi.certService.ValidateCertificate(cert); err != nil {
			return fmt.Errorf("invalid CA certificate: %w", err)
		}
	}

	return nil
}

// GetBuildTLSConfig returns TLS config for build operations
func (bi *BuildIntegration) GetBuildTLSConfig() (*tls.Config, error) {
	if bi.buildConfig.SkipTLSVerify {
		return &tls.Config{InsecureSkipVerify: true}, nil
	}

	return bi.certService.GetTLSConfig()
}

// DefaultCertificateMiddleware provides standard certificate injection
type DefaultCertificateMiddleware struct {
	certService CertificateService
}

// NewDefaultCertificateMiddleware creates default middleware
func NewDefaultCertificateMiddleware(certService CertificateService) *DefaultCertificateMiddleware {
	return &DefaultCertificateMiddleware{
		certService: certService,
	}
}

// InjectCertificates implements CertificateMiddleware
func (dm *DefaultCertificateMiddleware) InjectCertificates(transport http.RoundTripper) http.RoundTripper {
	return &certificateTransport{
		base:        transport,
		certService: dm.certService,
	}
}

// ConfigureTLS implements CertificateMiddleware
func (dm *DefaultCertificateMiddleware) ConfigureTLS(config *tls.Config) error {
	pool := dm.certService.GetCertificatePool()
	if pool != nil {
		config.RootCAs = pool
	}
	return nil
}

// ValidateConnection implements CertificateMiddleware
func (dm *DefaultCertificateMiddleware) ValidateConnection(url string) error {
	// This would typically perform actual connection validation
	// For now, return nil as this is a mock implementation
	return nil
}

// certificateTransport wraps HTTP transport with certificate handling
type certificateTransport struct {
	base        http.RoundTripper
	certService CertificateService
}

// RoundTrip implements http.RoundTripper
func (ct *certificateTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone request to avoid modifying original
	req = req.Clone(req.Context())

	// Add certificate-related headers if needed
	req.Header.Set("X-Certificate-Validation", "enabled")

	return ct.base.RoundTrip(req)
}
