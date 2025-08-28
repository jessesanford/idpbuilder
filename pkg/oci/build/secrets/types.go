package secrets

import (
	"context"
	"os"
)

// SecretType represents the type of secret being handled
type SecretType string

const (
	// SecretTypeBuildArg represents build arguments passed to the build process
	SecretTypeBuildArg SecretType = "build-arg"
	// SecretTypeMount represents secrets mounted as files during build
	SecretTypeMount SecretType = "mount"
	// SecretTypeEnv represents environment variable secrets
	SecretTypeEnv SecretType = "env"
	// SecretTypeSSH represents SSH authentication secrets
	SecretTypeSSH SecretType = "ssh"
)

// Secret represents a secret with its metadata and secure handling properties
type Secret struct {
	ID     string      // Unique identifier for the secret
	Type   SecretType  // Type of secret (build-arg, mount, env, ssh)
	Value  []byte      // Secret value stored as bytes for secure handling
	Source string      // Source location (file path, env var name, etc.)
	Target string      // Target destination (mount path, arg name, etc.)
	Mode   os.FileMode // File permissions for mounted secrets
	UID    int         // User ID for mounted secrets
	GID    int         // Group ID for mounted secrets
}

// SecretProvider defines the interface for external secret providers
type SecretProvider interface {
	// GetSecret retrieves a secret by its identifier
	GetSecret(ctx context.Context, id string) (*Secret, error)
	// ListSecrets returns all available secret identifiers
	ListSecrets(ctx context.Context) ([]string, error)
}

// SecretVault defines the interface for secure in-memory secret storage
type SecretVault interface {
	// Store securely stores a secret in the vault
	Store(secret *Secret) error
	// Retrieve retrieves a secret from the vault by ID
	Retrieve(id string) (*Secret, error)
	// Delete removes a secret from the vault
	Delete(id string) error
	// Clear securely removes all secrets from the vault
	Clear() error
}

// secretMetadata holds non-sensitive metadata about secrets
type secretMetadata struct {
	Type   SecretType
	Source string
	Target string
	Mode   os.FileMode
	UID    int
	GID    int
}