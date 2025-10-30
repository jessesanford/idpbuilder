// Package registry provides OCI registry push operations.
package registry

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/cnoe-io/idpbuilder/pkg/auth"
	tlspkg "github.com/cnoe-io/idpbuilder/pkg/tls"
)

// registryClient implements the Client interface using go-containerregistry.
type registryClient struct {
	authProvider auth.Provider
	tlsConfig    tlspkg.ConfigProvider
	httpClient   *http.Client
}

// newClientImpl creates a new registry client with authentication and TLS configuration.
//
// The client is configured with:
//   - Authentication provider for registry credentials
//   - TLS configuration for secure/insecure mode
//   - HTTP transport with proper timeouts
//
// Parameters:
//   - authProvider: Authentication provider (from pkg/auth)
//   - tlsConfig: TLS configuration provider (from pkg/tls)
//
// Returns:
//   - Client: Registry client interface implementation
//   - error: ValidationError if providers are invalid
//
// Example:
//   authProvider := auth.NewBasicAuthProvider("giteaadmin", "password")
//   tlsProvider := tls.NewConfigProvider(insecure)
//   client, err := registry.NewClient(authProvider, tlsProvider)
//   if err != nil {
//       return fmt.Errorf("failed to create registry client: %w", err)
//   }
func newClientImpl(authProvider auth.Provider, tlsConfig tlspkg.ConfigProvider) (Client, error) {
	if authProvider == nil {
		return nil, &ValidationError{
			Field:   "authProvider",
			Message: "authentication provider cannot be nil",
		}
	}
	if tlsConfig == nil {
		return nil, &ValidationError{
			Field:   "tlsConfig",
			Message: "TLS config provider cannot be nil",
		}
	}

	// Create HTTP client with TLS config
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig.GetTLSConfig(),
		},
	}

	return &registryClient{
		authProvider: authProvider,
		tlsConfig:    tlsConfig,
		httpClient:   httpClient,
	}, nil
}

// Push pushes an OCI image to the specified registry with optional progress reporting.
//
// This method:
//   1. Parses the target reference
//   2. Gets authenticator from auth provider
//   3. Configures remote options (auth, transport, progress)
//   4. Calls go-containerregistry's remote.Write()
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - image: OCI v1.Image to push (from Docker client)
//   - targetRef: Fully qualified image reference
//                Format: "registry-host/namespace/repository:tag"
//   - progressCallback: Optional callback for progress updates (can be nil)
//
// Returns:
//   - error: AuthenticationError if credentials invalid (401/403),
//            NetworkError if registry unreachable,
//            PushFailedError if layer upload or manifest push fails
//
// Example:
//   err := client.Push(ctx, image, "registry.io/repo/myapp:latest", func(update ProgressUpdate) {
//       fmt.Printf("Layer %s: %d/%d bytes\n", update.LayerDigest, update.BytesPushed, update.LayerSize)
//   })
func (c *registryClient) Push(ctx context.Context, image v1.Image, targetRef string, progressCallback ProgressCallback) error {
	// Parse target reference
	ref, err := name.ParseReference(targetRef)
	if err != nil {
		return &PushFailedError{
			TargetRef: targetRef,
			Cause:     fmt.Errorf("invalid reference: %w", err),
		}
	}

	// Get authenticator
	authenticator, err := c.authProvider.GetAuthenticator()
	if err != nil {
		return &AuthenticationError{
			Registry: targetRef,
			Cause:    err,
		}
	}

	// Configure remote options
	options := []remote.Option{
		remote.WithAuth(authenticator),
		remote.WithTransport(c.httpClient.Transport),
		remote.WithContext(ctx),
	}

	// Add progress callback if provided
	if progressCallback != nil {
		options = append(options, remote.WithProgress(createProgressHandler(progressCallback)))
	}

	// Push image
	err = remote.Write(ref, image, options...)
	if err != nil {
		// Classify error type
		if isAuthError(err) {
			return &AuthenticationError{
				Registry: targetRef,
				Cause:    err,
			}
		}
		if isNetworkError(err) {
			return &NetworkError{
				Registry: targetRef,
				Cause:    err,
			}
		}
		return &PushFailedError{
			TargetRef: targetRef,
			Cause:     err,
		}
	}

	return nil
}

// BuildImageReference constructs a fully qualified registry image reference.
//
// This method:
//   1. Parses the registry URL to extract host:port
//   2. Parses image name to extract repository and tag
//   3. Constructs full reference: registry/namespace/repository:tag
//   4. Uses "giteaadmin" as default namespace for Gitea registries
//
// Parameters:
//   - registryURL: Base registry URL
//                  Examples: "https://gitea.cnoe.localtest.me:8443"
//                           "https://registry.io"
//   - imageName: Image name with optional tag
//                Examples: "myapp:latest", "myapp", "myapp:v1.0.0"
//
// Returns:
//   - string: Fully qualified image reference
//   - error: ValidationError if registry URL or image name is invalid
//
// Example:
//   ref, err := client.BuildImageReference(
//       "https://gitea.cnoe.localtest.me:8443",
//       "myapp:latest",
//   )
//   // ref = "gitea.cnoe.localtest.me:8443/giteaadmin/myapp:latest"
func (c *registryClient) BuildImageReference(registryURL, imageName string) (string, error) {
	// Parse registry URL
	parsedURL, err := url.Parse(registryURL)
	if err != nil {
		return "", &ValidationError{
			Field:   "registryURL",
			Message: fmt.Sprintf("invalid registry URL: %v", err),
		}
	}

	// Extract host:port
	registryHost := parsedURL.Host
	if registryHost == "" {
		return "", &ValidationError{
			Field:   "registryURL",
			Message: "registry URL must include host",
		}
	}

	// Parse image name
	repository, tag := parseImageName(imageName)
	if repository == "" {
		return "", &ValidationError{
			Field:   "imageName",
			Message: "image name cannot be empty",
		}
	}

	// Default tag if not specified
	if tag == "" {
		tag = "latest"
	}

	// Build full reference with "giteaadmin" namespace
	namespace := "giteaadmin"
	fullRef := fmt.Sprintf("%s/%s/%s:%s", registryHost, namespace, repository, tag)
	return fullRef, nil
}

// ValidateRegistry checks if the registry is reachable by pinging the /v2/ endpoint.
//
// This method performs a GET request to the registry's /v2/ endpoint to verify:
//   - Registry is accessible (network connectivity)
//   - Registry responds (service is running)
//   - Registry speaks OCI protocol (returns 200 or 401)
//
// A 401 (Unauthorized) response is considered success because it indicates
// the registry is accessible and requires authentication.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - registryURL: Registry base URL to validate
//
// Returns:
//   - error: NetworkError if unreachable,
//            RegistryUnavailableError if invalid response,
//            ValidationError if URL is malformed
//
// Example:
//   if err := client.ValidateRegistry(ctx, "https://registry.io"); err != nil {
//       return fmt.Errorf("registry validation failed: %w", err)
//   }
func (c *registryClient) ValidateRegistry(ctx context.Context, registryURL string) error {
	// Parse registry URL
	parsedURL, err := url.Parse(registryURL)
	if err != nil {
		return &ValidationError{
			Field:   "registryURL",
			Message: fmt.Sprintf("invalid registry URL: %v", err),
		}
	}

	// Build /v2/ endpoint URL
	v2URL := fmt.Sprintf("%s://%s/v2/", parsedURL.Scheme, parsedURL.Host)

	// Create HTTP GET request
	req, err := http.NewRequestWithContext(ctx, "GET", v2URL, nil)
	if err != nil {
		return &NetworkError{
			Registry: registryURL,
			Cause:    err,
		}
	}

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &NetworkError{
			Registry: registryURL,
			Cause:    err,
		}
	}
	defer resp.Body.Close()

	// Check response status
	// 200 OK = registry accessible and doesn't require auth
	// 401 Unauthorized = registry accessible and requires auth (success!)
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized {
		return nil
	}

	return &RegistryUnavailableError{
		Registry:   registryURL,
		StatusCode: resp.StatusCode,
	}
}

// parseImageName extracts repository and tag from image name.
//
// Examples:
//   - "myapp:v1.0" → ("myapp", "v1.0")
//   - "myapp" → ("myapp", "")
//   - "repo/myapp:latest" → ("repo/myapp", "latest")
func parseImageName(imageName string) (repository, tag string) {
	parts := strings.Split(imageName, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

// createProgressHandler converts ProgressCallback to v1.Update channel for go-containerregistry.
//
// This function:
//   1. Creates a buffered channel (100 capacity for smooth progress)
//   2. Starts a goroutine to convert v1.Update → ProgressUpdate
//   3. Returns the channel for remote.Write() to send updates
//
// The goroutine processes updates and calls the callback until the channel is closed.
// Note: LayerDigest is set to empty string as v1.Update doesn't provide digest information.
func createProgressHandler(callback ProgressCallback) chan v1.Update {
	updates := make(chan v1.Update, 100)
	go func() {
		for update := range updates {
			status := "uploading"
			if update.Complete == update.Total {
				status = "complete"
			}
			callback(ProgressUpdate{
				LayerDigest: "", // v1.Update doesn't include digest
				LayerSize:   update.Total,
				BytesPushed: update.Complete,
				Status:      status,
			})
		}
	}()
	return updates
}

// isAuthError classifies authentication failures from error messages.
//
// Detection patterns:
//   - HTTP 401 status code
//   - HTTP 403 status code
//   - Error message contains "unauthorized"
//   - Error message contains "forbidden"
func isAuthError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "401") || strings.Contains(errStr, "403") ||
		strings.Contains(errStr, "unauthorized") || strings.Contains(errStr, "forbidden")
}

// isNetworkError classifies network connectivity failures from error messages.
//
// Detection patterns:
//   - Error message contains "connection"
//   - Error message contains "timeout"
//   - Error message contains "network"
func isNetworkError(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "connection") || strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "network")
}
