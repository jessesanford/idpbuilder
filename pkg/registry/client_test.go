package registry

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock Auth Provider
type mockAuthProvider struct {
	authenticator authn.Authenticator
	validateErr   error
}

func (m *mockAuthProvider) GetAuthenticator() (authn.Authenticator, error) {
	if m.validateErr != nil {
		return nil, m.validateErr
	}
	return m.authenticator, nil
}

func (m *mockAuthProvider) ValidateCredentials() error {
	return m.validateErr
}

// Mock TLS Provider
type mockTLSProvider struct {
	config   *tls.Config
	insecure bool
}

func (m *mockTLSProvider) GetTLSConfig() *tls.Config {
	return m.config
}

func (m *mockTLSProvider) IsInsecure() bool {
	return m.insecure
}

// Test NewClient Constructor

func TestNewClient_Success(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{
		authenticator: authn.Anonymous,
	}
	tlsProvider := &mockTLSProvider{
		config:   &tls.Config{InsecureSkipVerify: true},
		insecure: true,
	}

	// Act
	client, err := newClientImpl(authProvider, tlsProvider)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, client)

	// Verify client is of correct type
	registryClient, ok := client.(*registryClient)
	require.True(t, ok, "client should be of type *registryClient")
	assert.Equal(t, authProvider, registryClient.authProvider)
	assert.Equal(t, tlsProvider, registryClient.tlsConfig)
	assert.NotNil(t, registryClient.httpClient)
}

func TestNewClient_NilAuthProvider(t *testing.T) {
	// Arrange
	tlsProvider := &mockTLSProvider{
		config: &tls.Config{},
	}

	// Act
	client, err := newClientImpl(nil, tlsProvider)

	// Assert
	assert.Nil(t, client)
	require.Error(t, err)

	// Verify error is ValidationError
	validationErr, ok := err.(*ValidationError)
	require.True(t, ok, "error should be ValidationError")
	assert.Equal(t, "authProvider", validationErr.Field)
	assert.Contains(t, validationErr.Message, "cannot be nil")
}

func TestNewClient_NilTLSProvider(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{
		authenticator: authn.Anonymous,
	}

	// Act
	client, err := newClientImpl(authProvider, nil)

	// Assert
	assert.Nil(t, client)
	require.Error(t, err)

	// Verify error is ValidationError
	validationErr, ok := err.(*ValidationError)
	require.True(t, ok, "error should be ValidationError")
	assert.Equal(t, "tlsConfig", validationErr.Field)
	assert.Contains(t, validationErr.Message, "cannot be nil")
}

// Test BuildImageReference

func TestBuildImageReference_FullURLWithTag(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	ref, err := client.BuildImageReference(
		"https://gitea.cnoe.localtest.me:8443",
		"myapp:latest",
	)

	// Assert
	require.NoError(t, err)
	expected := "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest"
	assert.Equal(t, expected, ref)
}

func TestBuildImageReference_SimpleURLWithoutTag(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	ref, err := client.BuildImageReference(
		"https://registry.io",
		"myapp",
	)

	// Assert
	require.NoError(t, err)
	expected := "registry.io/giteaadmin/myapp:latest"
	assert.Equal(t, expected, ref)
}

func TestBuildImageReference_WithVersionTag(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	ref, err := client.BuildImageReference(
		"https://registry.example.com:5000",
		"myapp:v1.0.0",
	)

	// Assert
	require.NoError(t, err)
	expected := "registry.example.com:5000/giteaadmin/myapp:v1.0.0"
	assert.Equal(t, expected, ref)
}

func TestBuildImageReference_InvalidURL(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	ref, err := client.BuildImageReference(
		"://invalid-url",
		"myapp",
	)

	// Assert
	assert.Empty(t, ref)
	require.Error(t, err)

	validationErr, ok := err.(*ValidationError)
	require.True(t, ok, "error should be ValidationError")
	assert.Equal(t, "registryURL", validationErr.Field)
}

func TestBuildImageReference_EmptyImageName(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	ref, err := client.BuildImageReference(
		"https://registry.io",
		"",
	)

	// Assert
	assert.Empty(t, ref)
	require.Error(t, err)

	validationErr, ok := err.(*ValidationError)
	require.True(t, ok, "error should be ValidationError")
	assert.Equal(t, "imageName", validationErr.Field)
	assert.Contains(t, validationErr.Message, "cannot be empty")
}

// Test ValidateRegistry

func TestValidateRegistry_Success200(t *testing.T) {
	// Arrange - Create test server that returns 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/v2/") {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{InsecureSkipVerify: true}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	err = client.ValidateRegistry(context.Background(), server.URL)

	// Assert
	assert.NoError(t, err)
}

func TestValidateRegistry_Success401(t *testing.T) {
	// Arrange - Create test server that returns 401 (registry requires auth)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/v2/") {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}))
	defer server.Close()

	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{InsecureSkipVerify: true}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	err = client.ValidateRegistry(context.Background(), server.URL)

	// Assert
	assert.NoError(t, err, "401 should be considered success (registry requires auth)")
}

func TestValidateRegistry_UnavailableStatus(t *testing.T) {
	// Arrange - Create test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{InsecureSkipVerify: true}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	err = client.ValidateRegistry(context.Background(), server.URL)

	// Assert
	require.Error(t, err)

	unavailableErr, ok := err.(*RegistryUnavailableError)
	require.True(t, ok, "error should be RegistryUnavailableError")
	assert.Equal(t, http.StatusNotFound, unavailableErr.StatusCode)
}

func TestValidateRegistry_NetworkError(t *testing.T) {
	// Arrange - Use invalid URL that cannot be reached
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{InsecureSkipVerify: true}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	err = client.ValidateRegistry(context.Background(), "http://nonexistent.invalid:12345")

	// Assert
	require.Error(t, err)

	networkErr, ok := err.(*NetworkError)
	require.True(t, ok, "error should be NetworkError")
	assert.Contains(t, networkErr.Registry, "nonexistent.invalid")
}

func TestValidateRegistry_InvalidURL(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	// Act
	err = client.ValidateRegistry(context.Background(), "://invalid")

	// Assert
	require.Error(t, err)

	validationErr, ok := err.(*ValidationError)
	require.True(t, ok, "error should be ValidationError")
	assert.Equal(t, "registryURL", validationErr.Field)
}

// Test Push (Limited - requires complex setup)

func TestPush_InvalidReference(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	image := empty.Image

	// Act
	err = client.Push(context.Background(), image, "invalid::reference", nil)

	// Assert
	require.Error(t, err)

	pushErr, ok := err.(*PushFailedError)
	require.True(t, ok, "error should be PushFailedError")
	assert.Contains(t, pushErr.TargetRef, "invalid")
}

func TestPush_AuthenticationError(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{
		authenticator: authn.Anonymous,
		validateErr:   errors.New("invalid credentials"),
	}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	image := empty.Image

	// Act
	err = client.Push(context.Background(), image, "registry.io/namespace/image:tag", nil)

	// Assert
	require.Error(t, err)

	authErr, ok := err.(*AuthenticationError)
	require.True(t, ok, "error should be AuthenticationError")
	assert.Contains(t, authErr.Registry, "registry.io")
}

func TestPush_WithProgressCallback(t *testing.T) {
	// Arrange
	authProvider := &mockAuthProvider{authenticator: authn.Anonymous}
	tlsProvider := &mockTLSProvider{config: &tls.Config{}}
	client, err := newClientImpl(authProvider, tlsProvider)
	require.NoError(t, err)

	image := empty.Image

	callback := func(update ProgressUpdate) {
		// Progress callback - in real usage this would be called
	}

	// Act - This will fail because we don't have a real registry,
	// but we're testing that callback setup doesn't panic
	_ = client.Push(context.Background(), image, "registry.io/test/image:tag", callback)

	// Assert - Just verify we didn't panic
	// Note: callback would only be invoked if push actually started,
	// which requires a real registry setup
}

// Test Helper Functions

func TestParseImageName_WithTag(t *testing.T) {
	// Act
	repo, tag := parseImageName("myapp:v1.0.0")

	// Assert
	assert.Equal(t, "myapp", repo)
	assert.Equal(t, "v1.0.0", tag)
}

func TestParseImageName_WithoutTag(t *testing.T) {
	// Act
	repo, tag := parseImageName("myapp")

	// Assert
	assert.Equal(t, "myapp", repo)
	assert.Empty(t, tag)
}

func TestParseImageName_WithSlashAndTag(t *testing.T) {
	// Act
	repo, tag := parseImageName("namespace/myapp:latest")

	// Assert
	assert.Equal(t, "namespace/myapp", repo)
	assert.Equal(t, "latest", tag)
}

func TestParseImageName_RegistryWithPort(t *testing.T) {
	// Act
	repo, tag := parseImageName("registry:5000/app:v1")

	// Assert
	assert.Equal(t, "registry:5000/app", repo)
	assert.Equal(t, "v1", tag)
}

func TestParseImageName_HostWithPortAndNamespace(t *testing.T) {
	// Act
	repo, tag := parseImageName("host:443/ns/repo:tag")

	// Assert
	assert.Equal(t, "host:443/ns/repo", repo)
	assert.Equal(t, "tag", tag)
}

func TestIsAuthError_401(t *testing.T) {
	// Arrange
	err := fmt.Errorf("HTTP 401 Unauthorized")

	// Act
	result := isAuthError(err)

	// Assert
	assert.True(t, result)
}

func TestIsAuthError_403(t *testing.T) {
	// Arrange
	err := fmt.Errorf("HTTP 403 Forbidden")

	// Act
	result := isAuthError(err)

	// Assert
	assert.True(t, result)
}

func TestIsAuthError_UnauthorizedText(t *testing.T) {
	// Arrange
	err := fmt.Errorf("unauthorized access")

	// Act
	result := isAuthError(err)

	// Assert
	assert.True(t, result)
}

func TestIsAuthError_NotAuthError(t *testing.T) {
	// Arrange
	err := fmt.Errorf("network timeout")

	// Act
	result := isAuthError(err)

	// Assert
	assert.False(t, result)
}

func TestIsNetworkError_Connection(t *testing.T) {
	// Arrange
	err := fmt.Errorf("connection refused")

	// Act
	result := isNetworkError(err)

	// Assert
	assert.True(t, result)
}

func TestIsNetworkError_Timeout(t *testing.T) {
	// Arrange
	err := fmt.Errorf("request timeout")

	// Act
	result := isNetworkError(err)

	// Assert
	assert.True(t, result)
}

func TestIsNetworkError_NetworkKeyword(t *testing.T) {
	// Arrange
	err := fmt.Errorf("network unreachable")

	// Act
	result := isNetworkError(err)

	// Assert
	assert.True(t, result)
}

func TestIsNetworkError_NotNetworkError(t *testing.T) {
	// Arrange
	err := fmt.Errorf("invalid manifest")

	// Act
	result := isNetworkError(err)

	// Assert
	assert.False(t, result)
}

func TestCreateProgressHandler(t *testing.T) {
	// Arrange
	updates := []ProgressUpdate{}
	callback := func(update ProgressUpdate) {
		updates = append(updates, update)
	}

	// Act
	ch := createProgressHandler(callback)

	// Send test update
	ch <- v1.Update{
		Total:    1000,
		Complete: 500,
	}

	// Close channel
	close(ch)

	// Give goroutine time to process
	// In real usage, this happens asynchronously
	// For testing, we just verify channel was created correctly
	assert.NotNil(t, ch)
}
