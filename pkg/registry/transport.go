package registry

import (
	"net/http"
	"os"
)

// AuthConfig represents authentication configuration for registry access.
// This mirrors the configuration structure used by the registry configuration.
type AuthConfig struct {
	Type     string // "basic", "token", "anonymous"
	Username string
	Password string
	Token    string
}

// Transport defines the HTTP transport interface for registry operations.
// This abstraction allows for authentication injection and testing with mock transports.
type Transport interface {
	// RoundTrip executes a single HTTP transaction, returning a Response for the provided Request.
	RoundTrip(req *http.Request) (*http.Response, error)

	// WithAuth returns a new Transport configured with the provided authentication.
	WithAuth(auth *AuthConfig) Transport
}

// httpTransport implements Transport with authentication injection capability.
type httpTransport struct {
	base http.RoundTripper
	auth *AuthConfig
}

// NewTransport creates a new HTTP transport with optional authentication.
// If auth is nil, anonymous access is used.
func NewTransport(auth *AuthConfig) Transport {
	return &httpTransport{
		base: http.DefaultTransport,
		auth: auth,
	}
}

// RoundTrip implements the Transport interface by executing the HTTP request
// with authentication headers injected based on the configured auth type.
func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	clonedReq := req.Clone(req.Context())

	// Inject authentication based on configuration
	if t.auth != nil {
		switch t.auth.Type {
		case "basic":
			if t.auth.Username != "" && t.auth.Password != "" {
				clonedReq.SetBasicAuth(t.auth.Username, t.auth.Password)
			}
		case "token":
			if t.auth.Token != "" {
				clonedReq.Header.Set("Authorization", "Bearer "+t.auth.Token)
			}
		default:
			// Anonymous authentication - no headers added
		}
	}

	// Execute the request using the base transport
	return t.base.RoundTrip(clonedReq)
}

// WithAuth returns a new Transport instance configured with the provided authentication.
// This allows for immutable auth configuration changes.
func (t *httpTransport) WithAuth(auth *AuthConfig) Transport {
	return &httpTransport{
		base: t.base,
		auth: auth,
	}
}

// getAuthFromEnv creates an AuthConfig from environment variables.
// This supports common environment-based authentication patterns.
func getAuthFromEnv() *AuthConfig {
	username := os.Getenv("REGISTRY_USERNAME")
	password := os.Getenv("REGISTRY_PASSWORD")
	token := os.Getenv("REGISTRY_TOKEN")

	if token != "" {
		return &AuthConfig{
			Type:  "token",
			Token: token,
		}
	}

	if username != "" && password != "" {
		return &AuthConfig{
			Type:     "basic",
			Username: username,
			Password: password,
		}
	}

	return &AuthConfig{
		Type: "anonymous",
	}
}