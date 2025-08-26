package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/cnoe-io/idpbuilder/pkg/oci/api"
)

// authHandler manages registry authentication including special handling
// for gitea.cnoe.localtest.me and other self-signed certificate registries
type authHandler struct {
	config       api.AuthConfig
	tokens       map[string]*tokenResponse
	httpClient   *http.Client
	insecure     bool
}

// tokenResponse represents the OAuth2/Docker Registry V2 token response
type tokenResponse struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
	Scope       string    `json:"scope"`
}

// authChallenge represents a WWW-Authenticate challenge from the registry
type authChallenge struct {
	Scheme string
	Realm  string
	Service string
	Scope  string
}

// newAuthHandler creates a new authentication handler
func newAuthHandler(config api.AuthConfig) *authHandler {
	return &authHandler{
		config: config,
		tokens: make(map[string]*tokenResponse),
		insecure: config.AllowInsecureRegistry || 
			strings.Contains(config.ServerAddress, "gitea.cnoe.localtest.me"),
	}
}

// setHTTPClient configures the HTTP client for authentication requests
func (ah *authHandler) setHTTPClient(client *http.Client) {
	ah.httpClient = client
}

// authenticate adds appropriate authentication headers to the request
// Handles basic auth, bearer tokens, and OAuth2 flows
func (ah *authHandler) authenticate(req *http.Request, challenge *authChallenge) error {
	// Special handling for gitea.cnoe.localtest.me
	if strings.Contains(ah.config.ServerAddress, "gitea.cnoe.localtest.me") {
		return ah.handleGiteaAuth(req)
	}

	// Handle authentication challenge from registry
	if challenge != nil {
		return ah.handleAuthChallenge(req, challenge)
	}

	// Try bearer token if available
	if ah.config.RegistryToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ah.config.RegistryToken))
		return nil
	}

	// Fall back to basic authentication
	if ah.config.Username != "" && ah.config.Password != "" {
		return ah.setBasicAuth(req, ah.config.Username, ah.config.Password)
	}

	return nil
}

// handleGiteaAuth implements gitea-specific authentication
// Gitea supports both basic auth and token auth, prioritizing basic auth for simplicity
func (ah *authHandler) handleGiteaAuth(req *http.Request) error {
	// Prefer basic authentication for gitea
	if ah.config.Username != "" && ah.config.Password != "" {
		return ah.setBasicAuth(req, ah.config.Username, ah.config.Password)
	}

	// Try token authentication if available
	if ah.config.RegistryToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ah.config.RegistryToken))
		return nil
	}

	// No credentials available - proceed without auth (anonymous pulls may work)
	return nil
}

// setBasicAuth sets basic authentication header
func (ah *authHandler) setBasicAuth(req *http.Request, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password required for basic auth")
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", username, password)),
	)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	return nil
}

// handleAuthChallenge processes WWW-Authenticate challenge from registry
// Implements Docker Registry V2 authentication protocol
func (ah *authHandler) handleAuthChallenge(req *http.Request, challenge *authChallenge) error {
	switch strings.ToLower(challenge.Scheme) {
	case "basic":
		return ah.handleBasicChallenge(req, challenge)
	case "bearer":
		return ah.handleBearerChallenge(req, challenge)
	default:
		return fmt.Errorf("unsupported authentication scheme: %s", challenge.Scheme)
	}
}

// handleBasicChallenge processes basic authentication challenge
func (ah *authHandler) handleBasicChallenge(req *http.Request, challenge *authChallenge) error {
	if ah.config.Username != "" && ah.config.Password != "" {
		return ah.setBasicAuth(req, ah.config.Username, ah.config.Password)
	}
	return fmt.Errorf("basic authentication required but no credentials provided")
}

// handleBearerChallenge processes bearer token challenge
// Implements OAuth2/Docker Registry V2 token exchange
func (ah *authHandler) handleBearerChallenge(req *http.Request, challenge *authChallenge) error {
	// Check if we have a cached valid token
	tokenKey := fmt.Sprintf("%s:%s", challenge.Service, challenge.Scope)
	if token, exists := ah.tokens[tokenKey]; exists {
		if !ah.isTokenExpired(token) {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
			return nil
		}
		// Token expired, remove from cache
		delete(ah.tokens, tokenKey)
	}

	// Obtain new token
	token, err := ah.obtainBearerToken(challenge)
	if err != nil {
		return fmt.Errorf("failed to obtain bearer token: %w", err)
	}

	// Cache the token
	ah.tokens[tokenKey] = token

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Token))
	return nil
}

// obtainBearerToken performs OAuth2 token exchange with the registry
func (ah *authHandler) obtainBearerToken(challenge *authChallenge) (*tokenResponse, error) {
	if challenge.Realm == "" {
		return nil, fmt.Errorf("no realm specified in bearer challenge")
	}

	// Build token request URL
	tokenURL, err := url.Parse(challenge.Realm)
	if err != nil {
		return nil, fmt.Errorf("invalid realm URL: %w", err)
	}

	// Add query parameters
	params := tokenURL.Query()
	if challenge.Service != "" {
		params.Set("service", challenge.Service)
	}
	if challenge.Scope != "" {
		params.Set("scope", challenge.Scope)
	}
	tokenURL.RawQuery = params.Encode()

	// Create token request
	tokenReq, err := http.NewRequest("GET", tokenURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// Add basic auth if credentials available
	if ah.config.Username != "" && ah.config.Password != "" {
		ah.setBasicAuth(tokenReq, ah.config.Username, ah.config.Password)
	}

	// Use the same HTTP client to respect TLS settings
	client := ah.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	// Perform token request
	resp, err := client.Do(tokenReq)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status: %s", resp.Status)
	}

	// Parse token response
	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Set issued time for expiration tracking
	tokenResp.IssuedAt = time.Now()

	return &tokenResp, nil
}

// isTokenExpired checks if the token has expired
func (ah *authHandler) isTokenExpired(token *tokenResponse) bool {
	if token.ExpiresIn <= 0 {
		return false // No expiration specified
	}

	expiry := token.IssuedAt.Add(time.Duration(token.ExpiresIn) * time.Second)
	// Add 5 minute buffer to avoid using tokens that expire very soon
	buffer := time.Now().Add(5 * time.Minute)
	
	return expiry.Before(buffer)
}

// parseAuthChallenge parses WWW-Authenticate header
func parseAuthChallenge(challengeHeader string) (*authChallenge, error) {
	if challengeHeader == "" {
		return nil, fmt.Errorf("empty challenge header")
	}

	// Split scheme from parameters
	parts := strings.SplitN(challengeHeader, " ", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid challenge format")
	}

	challenge := &authChallenge{
		Scheme: strings.TrimSpace(parts[0]),
	}

	// Parse parameters (key="value", key="value")
	params := parts[1]
	pairs := strings.Split(params, ",")
	
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(strings.TrimSpace(kv[1]), `"`)
		
		switch strings.ToLower(key) {
		case "realm":
			challenge.Realm = value
		case "service":
			challenge.Service = value
		case "scope":
			challenge.Scope = value
		}
	}

	return challenge, nil
}

// clearTokenCache clears the token cache for re-authentication
func (ah *authHandler) clearTokenCache() {
	ah.tokens = make(map[string]*tokenResponse)
}