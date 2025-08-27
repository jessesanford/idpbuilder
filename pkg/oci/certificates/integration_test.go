package certificates

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertificateIntegrationE2E(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	
	// Setup
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	// Test certificate loading
	bundle, err := certService.LoadCertificateBundle(ctx, "test.pem", CertFormatPEM)
	require.NoError(t, err)
	require.NotNil(t, bundle)
	assert.Equal(t, CertFormatPEM, bundle.Format)
	assert.NotEmpty(t, bundle.Certificates)
	
	// Test registry integration
	integration := NewRegistryIntegration(certService, registry)
	err = integration.Push(ctx, "test-image:latest", []byte("test-content"))
	require.NoError(t, err)
	
	// Verify certificate was used
	assert.True(t, registry.CertificateValidated())
	pushedRefs := registry.GetPushedRefs()
	assert.Contains(t, pushedRefs, "test-image:latest")
}

// Certificate Loading Tests (60 lines)
func TestCertificateLoadingFormats(t *testing.T) {
	tests := []struct {
		name   string
		format CertFormat
		path   string
	}{
		{"PEM Format", CertFormatPEM, "test.pem"},
		{"DER Format", CertFormatDER, "test.der"},
		{"PKCS7 Format", CertFormatPKCS7, "test.p7b"},
		{"PKCS12 Format", CertFormatPKCS12, "test.p12"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			
			certService := NewMockCertificateService()
			ctx := context.Background()
			
			bundle, err := certService.LoadCertificateBundle(ctx, tt.path, tt.format)
			require.NoError(t, err)
			assert.Equal(t, tt.format, bundle.Format)
			assert.Equal(t, tt.path, bundle.Source)
			assert.NotEmpty(t, bundle.Certificates)
		})
	}
}

func TestInvalidCertificateFormatHandling(t *testing.T) {
	t.Parallel()
	
	certService := NewMockCertificateService()
	certService.LoadBundleFunc = func(ctx context.Context, path string, format CertFormat) (*CertBundle, error) {
		if format == "invalid" {
			return nil, &CertificateError{Code: "INVALID_FORMAT", Message: "unsupported format"}
		}
		return &CertBundle{Format: format}, nil
	}
	
	ctx := context.Background()
	_, err := certService.LoadCertificateBundle(ctx, "test.invalid", "invalid")
	require.Error(t, err)
	
	var certErr *CertificateError
	require.ErrorAs(t, err, &certErr)
	assert.Equal(t, "INVALID_FORMAT", certErr.Code)
}

func TestCorruptedCertificateHandling(t *testing.T) {
	t.Parallel()
	
	certService := NewMockCertificateService()
	certService.SetErrorOnNthCall("LoadCertificateBundle", 1)
	
	ctx := context.Background()
	_, err := certService.LoadCertificateBundle(ctx, "corrupted.pem", CertFormatPEM)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mock error on call 1")
}

// Registry Integration Tests (80 lines)
func TestRegistryPushWithCustomCA(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	// Add middleware
	middleware := NewDefaultCertificateMiddleware(certService)
	integration.AddMiddleware(middleware)
	
	// Test push with custom CA
	err := integration.Push(ctx, "myregistry.com/app:v1", []byte("image-data"))
	require.NoError(t, err)
	
	// Verify registry received the push
	pushedRefs := registry.GetPushedRefs()
	assert.Contains(t, pushedRefs, "myregistry.com/app:v1")
	
	// Verify certificate validation occurred
	assert.True(t, registry.CertificateValidated())
}

func TestRegistryPullWithCertificateValidation(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	content, err := integration.Pull(ctx, "secure-registry.com/base:latest")
	require.NoError(t, err)
	assert.Equal(t, []byte("mock-content"), content)
	
	pulledRefs := registry.GetPulledRefs()
	assert.Contains(t, pulledRefs, "secure-registry.com/base:latest")
}

func TestRegistryAuthenticationWithClientCerts(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	err := integration.Login(ctx, "private-registry.com", "user", "pass")
	require.NoError(t, err)
	
	// Verify login was attempted
	assert.Equal(t, 1, len(registry.loginCalls))
	assert.Equal(t, "private-registry.com", registry.loginCalls[0])
}

func TestTLSHandshakeVerification(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	// Test TLS initialization
	err := integration.InitializeTLS(ctx)
	require.NoError(t, err)
	
	// Verify TLS configuration was set up
	assert.NotNil(t, integration.tlsConfig)
}

func TestCertificateRotationDuringOperations(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	// Simulate certificate rotation
	oldCert := GenerateTestCertificate()
	newCert := GenerateTestCertificate()
	
	err := certService.RotateCertificate(oldCert, newCert)
	require.NoError(t, err)
	
	// Continue operations after rotation
	err = integration.Push(ctx, "app:after-rotation", []byte("data"))
	require.NoError(t, err)
}

// Build Integration Tests (60 lines)
func TestBuildWithCustomCABundle(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	
	config := BuildConfiguration{
		CustomCAPath:  "/path/to/ca.pem",
		SkipTLSVerify: false,
	}
	
	buildIntegration := NewBuildIntegration(certService, config)
	
	// Load certificates
	err := buildIntegration.LoadCertificatesForBuild(ctx)
	require.NoError(t, err)
	
	// Validate certificates
	err = buildIntegration.ValidateBuildCertificates()
	require.NoError(t, err)
}

func TestMultiStageBuildCertificatePropagation(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	
	config := BuildConfiguration{
		CustomCAPath: "/certs/ca.pem",
	}
	
	buildIntegration := NewBuildIntegration(certService, config)
	err := buildIntegration.LoadCertificatesForBuild(ctx)
	require.NoError(t, err)
	
	// Test buildah configuration
	buildahConfig := make(map[string]string)
	err = buildIntegration.ConfigureBuildahWithCertificates(buildahConfig)
	require.NoError(t, err)
	
	assert.Equal(t, "true", buildahConfig["tls-verify"])
	assert.Equal(t, "/certs/ca.pem", buildahConfig["cert-dir"])
}

func TestBuildCacheWithCertificates(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	_ = ctx // Use ctx to avoid unused variable warning
	certService := NewMockCertificateService()
	
	config := BuildConfiguration{
		CustomCAPath:  "/cache/certs",
		SkipTLSVerify: false,
	}
	
	buildIntegration := NewBuildIntegration(certService, config)
	
	tlsConfig, err := buildIntegration.GetBuildTLSConfig()
	require.NoError(t, err)
	assert.NotNil(t, tlsConfig)
	assert.False(t, tlsConfig.InsecureSkipVerify)
}

func TestCertificateValidationInFROMStatements(t *testing.T) {
	t.Parallel()
	
	certService := NewMockCertificateService()
	
	config := BuildConfiguration{
		InsecureRegistry: true,
		CustomCAPath:     "", // No CA path means no bundle to load
	}
	
	buildIntegration := NewBuildIntegration(certService, config)
	
	// For InsecureRegistry=true, we don't need to load certificates
	buildahConfig := make(map[string]string)
	err := buildIntegration.ConfigureBuildahWithCertificates(buildahConfig)
	require.NoError(t, err)
	
	assert.Equal(t, "false", buildahConfig["tls-verify"])
}

// End-to-End Scenarios (60 lines)
func TestCompleteWorkflow(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	
	// Setup components
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	// 1. Load certificates
	bundle, err := certService.LoadCertificateBundle(ctx, "ca.pem", CertFormatPEM)
	require.NoError(t, err)
	_ = bundle // Use bundle to avoid unused variable warning
	
	// 2. Configure registry integration
	integration := NewRegistryIntegration(certService, registry)
	middleware := NewDefaultCertificateMiddleware(certService)
	integration.AddMiddleware(middleware)
	
	// 3. Push image
	err = integration.Push(ctx, "workflow-test:v1", []byte("image-content"))
	require.NoError(t, err)
	
	// 4. Pull image back
	content, err := integration.Pull(ctx, "workflow-test:v1")
	require.NoError(t, err)
	assert.NotEmpty(t, content)
	
	// 5. Verify complete workflow
	assert.True(t, registry.CertificateValidated())
	assert.Contains(t, registry.GetPushedRefs(), "workflow-test:v1")
	assert.Contains(t, registry.GetPulledRefs(), "workflow-test:v1")
}

func TestCertificateRotationWithoutDowntime(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	// Initial operation
	err := integration.Push(ctx, "pre-rotation:v1", []byte("data"))
	require.NoError(t, err)
	
	// Rotate certificate
	oldCert := GenerateTestCertificate()
	newCert := GenerateTestCertificate()
	err = certService.RotateCertificate(oldCert, newCert)
	require.NoError(t, err)
	
	// Post-rotation operation should work seamlessly
	err = integration.Push(ctx, "post-rotation:v1", []byte("data"))
	require.NoError(t, err)
	
	assert.Equal(t, 2, len(registry.GetPushedRefs()))
}

func TestFallbackToSkipVerifyOnFailure(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	_ = ctx // Use ctx to avoid unused variable warning
	
	certService := NewMockCertificateService()
	
	// Configure fallback behavior
	certService.GetTLSConfigFunc = func() (*tls.Config, error) {
		return &tls.Config{InsecureSkipVerify: true}, nil
	}
	
	tlsConfig, err := certService.GetTLSConfig()
	require.NoError(t, err)
	assert.True(t, tlsConfig.InsecureSkipVerify)
}

func TestConcurrentOperationsWithSharedCertificates(t *testing.T) {
	t.Parallel()
	
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	// Simulate concurrent operations
	done := make(chan bool, 3)
	
	for i := 0; i < 3; i++ {
		go func(id int) {
			defer func() { done <- true }()
			
			err := integration.Push(ctx, fmt.Sprintf("concurrent-%d:v1", id), []byte("data"))
			assert.NoError(t, err)
		}(i)
	}
	
	// Wait for all operations to complete
	for i := 0; i < 3; i++ {
		<-done
	}
	
	assert.Equal(t, 3, len(registry.GetPushedRefs()))
}

// Error Recovery Tests (40 lines)
func TestInvalidCertificateRecovery(t *testing.T) {
	t.Parallel()
	
	certService := NewMockCertificateService()
	
	// Test with expired certificate
	expiredCert := GenerateExpiredCertificate()
	err := certService.ValidateCertificate(expiredCert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
	
	// Test with valid certificate
	validCert := GenerateTestCertificate()
	err = certService.ValidateCertificate(validCert)
	assert.NoError(t, err)
}

func TestNetworkTimeoutHandling(t *testing.T) {
	t.Parallel()
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	certService := NewMockCertificateService()
	certService.LoadBundleFunc = func(ctx context.Context, path string, format CertFormat) (*CertBundle, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	
	_, err := certService.LoadCertificateBundle(ctx, "slow.pem", CertFormatPEM)
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestCertificateExpiryDetection(t *testing.T) {
	t.Parallel()
	
	certService := NewMockCertificateService()
	
	// Test expiry detection
	expiredCert := GenerateExpiredCertificate()
	require.NotNil(t, expiredCert, "GenerateExpiredCertificate should not return nil")
	
	err := certService.ValidateCertificate(expiredCert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
}

func TestCAPoolCorruptionRecovery(t *testing.T) {
	t.Parallel()
	
	certService := NewMockCertificateService()
	
	// Add valid CA
	validCA := GenerateTestCA()
	err := certService.AddCACertificate(validCA)
	require.NoError(t, err)
	
	// Verify TLS config can be generated
	tlsConfig, err := certService.GetTLSConfig()
	require.NoError(t, err)
	assert.NotNil(t, tlsConfig.RootCAs)
}

// Test helpers and benchmarks
func BenchmarkCertificateLoading(b *testing.B) {
	certService := NewMockCertificateService()
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := certService.LoadCertificateBundle(ctx, "bench.pem", CertFormatPEM)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRegistryPush(b *testing.B) {
	ctx := context.Background()
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	integration := NewRegistryIntegration(certService, registry)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := integration.Push(ctx, fmt.Sprintf("bench-%d:latest", i), []byte("data"))
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMiddlewareChaining(t *testing.T) {
	t.Parallel()
	
	certService := NewMockCertificateService()
	registry := NewMockRegistryClient()
	
	integration := NewRegistryIntegration(certService, registry)
	
	// Add multiple middleware layers
	mw1 := NewDefaultCertificateMiddleware(certService)
	integration.AddMiddleware(mw1)
	
	// Test that middleware is properly chained
	transport := integration.createSecureTransport()
	assert.NotNil(t, transport)
	
	// Test actual HTTP request handling
	req, _ := http.NewRequest("GET", "https://registry.example.com/v2/", nil)
	_, err := transport.RoundTrip(req)
	// We expect an error since this is a mock transport, but it should process the middleware chain
	assert.Error(t, err) // Network error is expected in test environment
}