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

type authHandler struct {
	config       api.AuthConfig
	tokens       map[string]*tokenResponse
	httpClient   *http.Client
	insecure     bool
}

type tokenResponse struct {
	Token       string    `json:"token"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int       `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
	Scope       string    `json:"scope"`
}

type authChallenge struct {
	Scheme string
	Realm  string
	Service string
	Scope  string
}

func newAuthHandler(config api.AuthConfig) *authHandler {
	return &authHandler{
		config: config,
		tokens: make(map[string]*tokenResponse),
		insecure: config.AllowInsecureRegistry || strings.Contains(config.ServerAddress, "gitea.cnoe.localtest.me"),
	}
}

func (ah *authHandler) setHTTPClient(client *http.Client) {
	ah.httpClient = client
}

func (ah *authHandler) authenticate(req *http.Request, challenge *authChallenge) error {
	if strings.Contains(ah.config.ServerAddress, "gitea.cnoe.localtest.me") {
		return ah.handleGiteaAuth(req)
	}
	if challenge != nil {
		return ah.handleAuthChallenge(req, challenge)
	}
	if ah.config.RegistryToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ah.config.RegistryToken))
		return nil
	}
	if ah.config.Username != "" && ah.config.Password != "" {
		return ah.setBasicAuth(req, ah.config.Username, ah.config.Password)
	}
	return nil
}

func (ah *authHandler) handleGiteaAuth(req *http.Request) error {
	if ah.config.Username != "" && ah.config.Password != "" {
		return ah.setBasicAuth(req, ah.config.Username, ah.config.Password)
	}
	if ah.config.RegistryToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ah.config.RegistryToken))
		return nil
	}
	return nil
}

func (ah *authHandler) setBasicAuth(req *http.Request, username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password required for basic auth")
	}
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
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

func (ah *authHandler) obtainBearerToken(challenge *authChallenge) (*tokenResponse, error) {
	if challenge.Realm == "" {
		return nil, fmt.Errorf("no realm specified in bearer challenge")
	}
	tokenURL, err := url.Parse(challenge.Realm)
	if err != nil {
		return nil, fmt.Errorf("invalid realm URL: %w", err)
	}
	params := tokenURL.Query()
	if challenge.Service != "" {
		params.Set("service", challenge.Service)
	}
	if challenge.Scope != "" {
		params.Set("scope", challenge.Scope)
	}
	tokenURL.RawQuery = params.Encode()
	tokenReq, err := http.NewRequest("GET", tokenURL.String(), nil)
	if err != nil {
		return nil, err
	}
	if ah.config.Username != "" && ah.config.Password != "" {
		ah.setBasicAuth(tokenReq, ah.config.Username, ah.config.Password)
	}
	client := ah.httpClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(tokenReq)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status: %s", resp.Status)
	}
	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	tokenResp.IssuedAt = time.Now()
	return &tokenResp, nil
}

func (ah *authHandler) isTokenExpired(token *tokenResponse) bool {
	if token.ExpiresIn <= 0 {
		return false
	}
	expiry := token.IssuedAt.Add(time.Duration(token.ExpiresIn) * time.Second)
	return expiry.Before(time.Now().Add(5 * time.Minute))
}

func parseAuthChallenge(challengeHeader string) (*authChallenge, error) {
	if challengeHeader == "" {
		return nil, fmt.Errorf("empty challenge header")
	}
	parts := strings.SplitN(challengeHeader, " ", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid challenge format")
	}
	challenge := &authChallenge{Scheme: strings.TrimSpace(parts[0])}
	pairs := strings.Split(parts[1], ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
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

func (ah *authHandler) clearTokenCache() {
	ah.tokens = make(map[string]*tokenResponse)
}