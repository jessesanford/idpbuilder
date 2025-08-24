package auth

import (
	"errors"
	"time"
)

// HTTP header names for authentication and authorization.
const (
	// HeaderAuthorization is the standard HTTP Authorization header.
	HeaderAuthorization = "Authorization"

	// HeaderWWWAuthenticate is used in HTTP 401 responses to indicate required auth methods.
	HeaderWWWAuthenticate = "WWW-Authenticate"

	// HeaderDockerContentDigest is used for Docker registry content verification.
	HeaderDockerContentDigest = "Docker-Content-Digest"

	// HeaderUserAgent is used to identify the client making requests.
	HeaderUserAgent = "User-Agent"
)

// Authentication scheme prefixes for Authorization header values.
const (
	// SchemeBasic indicates HTTP Basic Authentication.
	SchemeBasic = "Basic"

	// SchemeBearer indicates Bearer token authentication.
	SchemeBearer = "Bearer"

	// SchemeOAuth indicates OAuth-style authentication.
	SchemeOAuth = "OAuth"
)

// Default token expiry durations for different authentication types.
const (
	// DefaultTokenExpiry is the default expiration time for bearer tokens.
	DefaultTokenExpiry = 1 * time.Hour

	// DefaultRefreshTokenExpiry is the default expiration time for refresh tokens.
	DefaultRefreshTokenExpiry = 24 * time.Hour

	// DefaultOAuth2TokenExpiry is the default expiration time for OAuth2 tokens.
	DefaultOAuth2TokenExpiry = 30 * time.Minute

	// TokenExpiryBuffer is the time buffer before token expiry to trigger refresh.
	TokenExpiryBuffer = 5 * time.Minute
)

// Registry URL patterns and default values.
const (
	// DockerHubRegistry is the default Docker Hub registry URL.
	DockerHubRegistry = "https://index.docker.io/v1/"

	// DockerHubRegistryV2 is the Docker Hub API v2 registry URL.
	DockerHubRegistryV2 = "https://registry-1.docker.io"

	// DefaultDockerConfigPath is the default path to Docker configuration.
	DefaultDockerConfigPath = "~/.docker/config.json"

	// DefaultCredentialStore is the default credential store name.
	DefaultCredentialStore = "desktop"
)

// OAuth2 related constants.
const (
	// OAuth2GrantTypeClientCredentials is the client credentials grant type.
	OAuth2GrantTypeClientCredentials = "client_credentials"

	// OAuth2GrantTypeRefreshToken is the refresh token grant type.
	OAuth2GrantTypeRefreshToken = "refresh_token"

	// OAuth2ScopeRepository is the scope for repository access.
	OAuth2ScopeRepository = "repository"

	// OAuth2ScopeRegistryCatalog is the scope for registry catalog access.
	OAuth2ScopeRegistryCatalog = "registry:catalog:*"
)

// Common error variables for authentication operations.
var (
	// ErrInvalidCredentials indicates that the provided credentials are invalid or incomplete.
	ErrInvalidCredentials = errors.New("invalid or incomplete credentials")

	// ErrCredentialsNotFound indicates that no credentials were found for the specified registry.
	ErrCredentialsNotFound = errors.New("credentials not found for registry")

	// ErrAuthenticationFailed indicates that authentication with the registry failed.
	ErrAuthenticationFailed = errors.New("authentication failed")

	// ErrTokenExpired indicates that the authentication token has expired.
	ErrTokenExpired = errors.New("authentication token has expired")

	// ErrRefreshNotSupported indicates that token refresh is not supported for this auth type.
	ErrRefreshNotSupported = errors.New("token refresh not supported for this authentication type")

	// ErrInvalidRegistry indicates that the registry URL is invalid or malformed.
	ErrInvalidRegistry = errors.New("invalid or malformed registry URL")

	// ErrCredentialStoreUnavailable indicates that the credential store is not available.
	ErrCredentialStoreUnavailable = errors.New("credential store is not available")

	// ErrUnsupportedAuthType indicates that the authentication type is not supported.
	ErrUnsupportedAuthType = errors.New("unsupported authentication type")
)