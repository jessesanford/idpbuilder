package auth

import (
	"fmt"
	"os"
)

// CredentialSource represents where credentials were sourced from
type CredentialSource int

const (
	// SourceNone indicates no credentials were found
	SourceNone CredentialSource = iota
	// SourceEnv indicates credentials came from environment variables
	SourceEnv
	// SourceFlags indicates credentials came from command line flags
	SourceFlags
	// SourceKeychain indicates credentials came from Docker keychain
	SourceKeychain
)

// String returns a string representation of the credential source
func (s CredentialSource) String() string {
	switch s {
	case SourceNone:
		return "none"
	case SourceEnv:
		return "environment"
	case SourceFlags:
		return "flags"
	case SourceKeychain:
		return "keychain"
	default:
		return "unknown"
	}
}

// Credentials represents authentication credentials with their source
type Credentials struct {
	Username string
	Password string
	Source   CredentialSource
}

// IsEmpty returns true if both username and password are empty
func (c *Credentials) IsEmpty() bool {
	return c.Username == "" && c.Password == ""
}

// Validate checks if the credentials are valid
func (c *Credentials) Validate() error {
	if c.Username == "" && c.Password != "" {
		return fmt.Errorf("username is required when password is provided")
	}
	if c.Username != "" && c.Password == "" {
		return fmt.Errorf("password is required when username is provided")
	}
	return nil
}

// GetCredentials returns credentials with proper precedence order:
// 1. Command line flags (highest priority)
// 2. Environment variables (fallback)
// 3. Docker config/keychain (will be handled by go-containerregistry DefaultKeychain)
func GetCredentials(flagUser, flagPass string) (*Credentials, error) {
	// First priority: command line flags
	if flagUser != "" || flagPass != "" {
		creds := &Credentials{
			Username: flagUser,
			Password: flagPass,
			Source:   SourceFlags,
		}

		if err := creds.Validate(); err != nil {
			return nil, fmt.Errorf("invalid flag credentials: %w", err)
		}

		return creds, nil
	}

	// Second priority: environment variables
	envUser := os.Getenv("IDPBUILDER_REGISTRY_USER")
	envPass := os.Getenv("IDPBUILDER_REGISTRY_PASSWORD")

	if envUser != "" || envPass != "" {
		creds := &Credentials{
			Username: envUser,
			Password: envPass,
			Source:   SourceEnv,
		}

		if err := creds.Validate(); err != nil {
			return nil, fmt.Errorf("invalid environment credentials: %w", err)
		}

		return creds, nil
	}

	// Third priority: Docker config/keychain
	// Return nil to indicate that DefaultKeychain should be used
	// The go-containerregistry library will handle Docker config automatically
	return nil, nil
}

// GetCredentialsFromEnv retrieves credentials from environment variables only
func GetCredentialsFromEnv() (*Credentials, error) {
	envUser := os.Getenv("IDPBUILDER_REGISTRY_USER")
	envPass := os.Getenv("IDPBUILDER_REGISTRY_PASSWORD")

	if envUser == "" && envPass == "" {
		return &Credentials{Source: SourceNone}, nil
	}

	creds := &Credentials{
		Username: envUser,
		Password: envPass,
		Source:   SourceEnv,
	}

	if err := creds.Validate(); err != nil {
		return nil, fmt.Errorf("invalid environment credentials: %w", err)
	}

	return creds, nil
}

// GetCredentialsFromFlags creates credentials from command line flags
func GetCredentialsFromFlags(username, password string) (*Credentials, error) {
	if username == "" && password == "" {
		return &Credentials{Source: SourceNone}, nil
	}

	creds := &Credentials{
		Username: username,
		Password: password,
		Source:   SourceFlags,
	}

	if err := creds.Validate(); err != nil {
		return nil, fmt.Errorf("invalid flag credentials: %w", err)
	}

	return creds, nil
}

// MergeCredentials merges credentials with flags taking precedence over environment
func MergeCredentials(flagCreds, envCreds *Credentials) *Credentials {
	if flagCreds != nil && !flagCreds.IsEmpty() {
		return flagCreds
	}

	if envCreds != nil && !envCreds.IsEmpty() {
		return envCreds
	}

	return &Credentials{Source: SourceNone}
}
