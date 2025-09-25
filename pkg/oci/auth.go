package oci

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Authenticator defines the interface for registry authentication methods.
// Different authentication schemes (basic, token, oauth) implement this interface
// to provide a consistent authentication mechanism for OCI registry operations.
type Authenticator interface {
	// Authenticate adds authentication headers to the HTTP request
	Authenticate(req *http.Request) error

	// GetType returns the authentication type for debugging and logging
	GetType() string

	// IsValid returns true if the authenticator has valid credentials
	IsValid() bool
}

// BasicAuthenticator implements HTTP Basic Authentication for registry access.
// It encodes username and password credentials in the Authorization header.
type BasicAuthenticator struct {
	username string
	password string
	encoded  string
}

// NewBasicAuthenticator creates a new basic authentication handler.
// It validates the provided credentials and pre-encodes them for efficiency.
func NewBasicAuthenticator(username, password string) Authenticator {
	if username == "" || password == "" {
		return nil
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return &BasicAuthenticator{
		username: username,
		password: password,
		encoded:  encoded,
	}
}

// Authenticate adds the Basic Authentication header to the request.
func (a *BasicAuthenticator) Authenticate(req *http.Request) error {
	if a == nil || !a.IsValid() {
		return NewRegistryError(401, "invalid basic authentication credentials")
	}

	req.Header.Set("Authorization", "Basic "+a.encoded)
	return nil
}

// GetType returns the authentication type for this authenticator.
func (a *BasicAuthenticator) GetType() string {
	return "basic"
}

// IsValid returns true if the authenticator has valid credentials.
func (a *BasicAuthenticator) IsValid() bool {
	return a != nil && a.username != "" && a.password != "" && a.encoded != ""
}

// TokenAuthenticator implements Bearer Token Authentication for registry access.
// It uses pre-obtained tokens (like Docker registry tokens) for authentication.
type TokenAuthenticator struct {
	token  string
	scheme string // Usually "Bearer"
}

// NewTokenAuthenticator creates a new token-based authentication handler.
// The token should be a valid bearer token obtained from the registry's auth service.
func NewTokenAuthenticator(token string) Authenticator {
	if token == "" {
		return nil
	}

	return &TokenAuthenticator{
		token:  token,
		scheme: "Bearer",
	}
}

// Authenticate adds the Bearer token header to the request.
func (a *TokenAuthenticator) Authenticate(req *http.Request) error {
	if a == nil || !a.IsValid() {
		return NewRegistryError(401, "invalid token authentication credentials")
	}

	req.Header.Set("Authorization", a.scheme+" "+a.token)
	return nil
}

// GetType returns the authentication type for this authenticator.
func (a *TokenAuthenticator) GetType() string {
	return "token"
}

// IsValid returns true if the authenticator has a valid token.
func (a *TokenAuthenticator) IsValid() bool {
	return a != nil && a.token != ""
}

// AuthDetectionResponse represents the response from a registry auth detection request.
type AuthDetectionResponse struct {
	Scheme      string            `json:"scheme"`
	Realm       string            `json:"realm,omitempty"`
	Service     string            `json:"service,omitempty"`
	Scope       string            `json:"scope,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty"`
	StatusCode  int               `json:"status_code"`
	Supported   bool              `json:"supported"`
}

// DetectAuthScheme probes a registry endpoint to determine the authentication scheme.
// It makes an unauthenticated request and analyzes the WWW-Authenticate header
// to determine what authentication method the registry supports.
func DetectAuthScheme(endpoint string) (string, error) {
	if endpoint == "" {
		return "", NewRegistryError(400, "endpoint cannot be empty")
	}

	// Ensure endpoint has a scheme
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		// Default to https for security
		useHTTPS := true
		if insecure := os.Getenv("REGISTRY_INSECURE"); insecure == "true" {
			useHTTPS = false
		}

		if useHTTPS {
			endpoint = "https://" + endpoint
		} else {
			endpoint = "http://" + endpoint
		}
	}

	// Parse the endpoint URL
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return "", NewRegistryError(400, fmt.Sprintf("invalid endpoint URL: %v", err))
	}

	// Build the auth detection URL (usually /v2/ endpoint)
	authURL := fmt.Sprintf("%s://%s/v2/", parsedURL.Scheme, parsedURL.Host)

	// Create HTTP client with timeout
	timeout := 10 * time.Second
	if timeoutStr := os.Getenv("REGISTRY_TIMEOUT"); timeoutStr != "" {
		if parsedTimeout, parseErr := time.ParseDuration(timeoutStr); parseErr == nil {
			timeout = parsedTimeout
		}
	}

	client := &http.Client{
		Timeout: timeout,
	}

	// Make an unauthenticated request to trigger auth challenge
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return "", NewRegistryError(500, fmt.Sprintf("failed to create auth detection request: %v", err))
	}

	req.Header.Set("User-Agent", "idpbuilder-push/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", NewRegistryError(500, fmt.Sprintf("failed to probe registry auth: %v", err))
	}
	defer resp.Body.Close()

	// Check for authentication challenge
	if resp.StatusCode == http.StatusUnauthorized {
		wwwAuth := resp.Header.Get("WWW-Authenticate")
		if wwwAuth == "" {
			return "", NewRegistryError(500, "registry requires authentication but no WWW-Authenticate header found")
		}

		// Parse the WWW-Authenticate header to determine scheme
		if strings.HasPrefix(strings.ToLower(wwwAuth), "basic") {
			return "basic", nil
		} else if strings.HasPrefix(strings.ToLower(wwwAuth), "bearer") {
			return "bearer", nil
		} else {
			return "unknown", NewRegistryError(500, fmt.Sprintf("unsupported authentication scheme: %s", wwwAuth))
		}
	}

	// If we get here, the registry might not require authentication
	// or allows anonymous access
	if resp.StatusCode == http.StatusOK {
		return "none", nil
	}

	// Other status codes indicate different issues
	return "", NewRegistryError(resp.StatusCode, fmt.Sprintf("unexpected response from registry: %s", resp.Status))
}

// CreateAuthenticatorFromConfig creates an appropriate Authenticator based on the provided RegistryAuth.
// It determines the best authentication method based on available credentials.
func CreateAuthenticatorFromConfig(auth *RegistryAuth) (Authenticator, error) {
	if auth == nil || auth.IsEmpty() {
		return nil, NewRegistryError(401, "no authentication credentials provided")
	}

	// Prefer token authentication if available
	if auth.HasTokenAuth() {
		return NewTokenAuthenticator(auth.Token), nil
	}

	// Fall back to basic authentication
	if auth.HasBasicAuth() {
		return NewBasicAuthenticator(auth.Username, auth.Password), nil
	}

	return nil, NewRegistryError(401, "no valid authentication credentials found")
}

// ValidateAuthenticator checks if an authenticator is properly configured and functional.
func ValidateAuthenticator(auth Authenticator) error {
	if auth == nil {
		return NewRegistryError(401, "authenticator cannot be nil")
	}

	if !auth.IsValid() {
		return NewRegistryError(401, fmt.Sprintf("invalid %s authenticator", auth.GetType()))
	}

	return nil
}