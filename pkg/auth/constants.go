// Copyright 2024 The IDP Builder Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import "time"

// AuthType represents the type of authentication mechanism
type AuthType string

const (
	// AuthTypeBasic represents HTTP Basic authentication
	AuthTypeBasic AuthType = "basic"
	
	// AuthTypeBearer represents Bearer token authentication
	AuthTypeBearer AuthType = "bearer"
	
	// AuthTypeOAuth2 represents OAuth2 authentication
	AuthTypeOAuth2 AuthType = "oauth2"
	
	// AuthTypeAnonymous represents anonymous access
	AuthTypeAnonymous AuthType = "anonymous"
)

// HTTP Header constants for authentication
const (
	// AuthorizationHeader is the standard HTTP Authorization header
	AuthorizationHeader = "Authorization"
	
	// WWWAuthenticateHeader is the standard WWW-Authenticate response header
	WWWAuthenticateHeader = "WWW-Authenticate"
	
	// ContentTypeHeader for JSON content
	ContentTypeHeader = "Content-Type"
	
	// ApplicationJSON content type for API requests
	ApplicationJSON = "application/json"
	
	// BasicAuthPrefix for Basic authentication scheme
	BasicAuthPrefix = "Basic "
	
	// BearerAuthPrefix for Bearer token scheme
	BearerAuthPrefix = "Bearer "
)

// Default expiry times for tokens and credentials
const (
	// DefaultTokenExpiry is the default expiry time for auth tokens
	DefaultTokenExpiry = 1 * time.Hour
	
	// DefaultRefreshThreshold is when to refresh tokens before expiry
	DefaultRefreshThreshold = 5 * time.Minute
	
	// MaxRetryAttempts for authentication operations
	MaxRetryAttempts = 3
	
	// RetryBackoffDelay between authentication retry attempts
	RetryBackoffDelay = 2 * time.Second
)

// Registry URL patterns and constants
const (
	// DockerHubRegistry is the default Docker Hub registry URL
	DockerHubRegistry = "https://index.docker.io/v1/"
	
	// DockerHubHost is the host for Docker Hub
	DockerHubHost = "index.docker.io"
	
	// LocalRegistryPattern for local registries
	LocalRegistryPattern = "localhost:"
	
	// HTTPScheme for insecure registry connections
	HTTPScheme = "http://"
	
	// HTTPSScheme for secure registry connections
	HTTPSScheme = "https://"
)

// Docker configuration file constants
const (
	// DockerConfigFilename is the standard Docker config file name
	DockerConfigFilename = "config.json"
	
	// DockerConfigDir is the standard Docker config directory
	DockerConfigDir = ".docker"
	
	// DockerCredentialStore key in config.json
	DockerCredentialStore = "credStore"
	
	// DockerCredentialHelpers key in config.json
	DockerCredentialHelpers = "credHelpers"
	
	// DockerAuths key in config.json
	DockerAuths = "auths"
)

// Error messages for authentication failures
const (
	// ErrInvalidCredentials when credentials are malformed or invalid
	ErrInvalidCredentials = "invalid credentials provided"
	
	// ErrTokenExpired when authentication token has expired
	ErrTokenExpired = "authentication token has expired"
	
	// ErrAuthenticationFailed when authentication attempt fails
	ErrAuthenticationFailed = "authentication failed"
	
	// ErrUnsupportedAuthType when auth type is not supported
	ErrUnsupportedAuthType = "unsupported authentication type"
	
	// ErrMissingCredentials when required credentials are not provided
	ErrMissingCredentials = "required credentials are missing"
	
	// ErrRegistryUnreachable when registry cannot be contacted
	ErrRegistryUnreachable = "registry is unreachable"
	
	// ErrPermissionDenied when access is denied
	ErrPermissionDenied = "permission denied"
	
	// ErrInvalidToken when token format is invalid
	ErrInvalidToken = "invalid token format"
)