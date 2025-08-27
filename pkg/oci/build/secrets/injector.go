package secrets

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Injector safely injects secrets into build processes while preventing leakage
type Injector struct {
	vault     SecretVault
	sanitizer *Sanitizer
}

// NewInjector creates a new secret injector with vault and sanitizer
func NewInjector(vault SecretVault, sanitizer *Sanitizer) *Injector {
	return &Injector{
		vault:     vault,
		sanitizer: sanitizer,
	}
}

// InjectBuildArgs safely injects secrets into build arguments
func (i *Injector) InjectBuildArgs(ctx context.Context, args map[string]string, secretIDs []string) (map[string]string, error) {
	if args == nil {
		args = make(map[string]string)
	}
	
	result := make(map[string]string)
	
	// Copy existing non-secret args
	for k, v := range args {
		result[k] = v
	}
	
	// Process each secret ID for injection
	for _, id := range secretIDs {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while injecting secret %s: %w", id, ctx.Err())
		default:
		}
		
		// Retrieve secret from vault
		secret, err := i.vault.Retrieve(id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve secret %s for build arg injection: %w", id, err)
		}
		
		// Verify it's a build arg type secret
		if secret.Type != SecretTypeBuildArg {
			return nil, fmt.Errorf("secret %s is type %s, not build-arg", id, secret.Type)
		}
		
		if secret.Target == "" {
			return nil, fmt.Errorf("secret %s has no target build arg name", id)
		}
		
		// Register with sanitizer to prevent leakage in logs
		if err := i.sanitizer.RegisterSecret(id, secret.Value); err != nil {
			log.Printf("Warning: failed to register secret %s with sanitizer: %v", id, err)
		}
		
		// Add to build args
		result[secret.Target] = string(secret.Value)
		
		// Audit log (without exposing the value)
		log.Printf("Injected build arg: %s (from secret: %s)", secret.Target, id)
		
		// Clear secret value from memory
		clearBytes(secret.Value)
	}
	
	return result, nil
}

// PrepareSecretMount creates a temporary secure mount for secret files
func (i *Injector) PrepareSecretMount(ctx context.Context, secretID string) (string, func(), error) {
	// Retrieve secret from vault
	secret, err := i.vault.Retrieve(secretID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to retrieve secret %s for mount: %w", secretID, err)
	}
	
	// Verify it's a mount type secret
	if secret.Type != SecretTypeMount {
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("secret %s is type %s, not mount", secretID, secret.Type)
	}
	
	// Create secure temporary directory
	tmpDir, err := os.MkdirTemp("", "oci-secret-mount-*")
	if err != nil {
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("failed to create temporary directory for secret mount: %w", err)
	}
	
	// Set restrictive permissions on temp directory (owner only)
	if err := os.Chmod(tmpDir, 0700); err != nil {
		os.RemoveAll(tmpDir)
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("failed to secure temporary directory: %w", err)
	}
	
	// Create target file path
	var tmpFile string
	if secret.Target != "" {
		tmpFile = filepath.Join(tmpDir, filepath.Base(secret.Target))
	} else {
		tmpFile = filepath.Join(tmpDir, "secret")
	}
	
	// Set default permissions if not specified
	mode := secret.Mode
	if mode == 0 {
		mode = 0600 // Owner read/write only by default
	}
	
	// Write secret to temporary file
	if err := os.WriteFile(tmpFile, secret.Value, mode); err != nil {
		os.RemoveAll(tmpDir)
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("failed to write secret to temporary file: %w", err)
	}
	
	// Set ownership if specified (requires appropriate privileges)
	if secret.UID > 0 || secret.GID > 0 {
		if err := os.Chown(tmpFile, secret.UID, secret.GID); err != nil {
			log.Printf("Warning: failed to set ownership for secret mount %s: %v", secretID, err)
			// Don't fail the operation for ownership issues
		}
	}
	
	// Register with sanitizer to prevent content leakage
	if err := i.sanitizer.RegisterSecret(secretID, secret.Value); err != nil {
		log.Printf("Warning: failed to register secret %s with sanitizer: %v", secretID, err)
	}
	
	// Clear secret from memory now that it's written to file
	clearBytes(secret.Value)
	
	// Create secure cleanup function
	cleanup := func() {
		// Read file back to securely overwrite it
		if data, err := os.ReadFile(tmpFile); err == nil {
			// Overwrite with zeros
			clearBytes(data)
			os.WriteFile(tmpFile, data, mode)
		}
		
		// Remove temporary directory and contents
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("Warning: failed to clean up secret mount directory %s: %v", tmpDir, err)
		}
		
		// Unregister from sanitizer
		i.sanitizer.UnregisterSecret(secretID)
		
		log.Printf("Cleaned up secret mount: %s", secretID)
	}
	
	// Audit log
	log.Printf("Created secret mount: %s -> %s (permissions: %o)", secretID, tmpFile, mode)
	
	return tmpFile, cleanup, nil
}

// InjectEnvironmentSecrets prepares environment variables with secrets
func (i *Injector) InjectEnvironmentSecrets(ctx context.Context, env []string, secretIDs []string) ([]string, error) {
	if env == nil {
		env = []string{}
	}
	
	result := make([]string, len(env))
	copy(result, env)
	
	// Process each secret ID for environment injection
	for _, id := range secretIDs {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled while injecting env secret %s: %w", id, ctx.Err())
		default:
		}
		
		// Retrieve secret from vault
		secret, err := i.vault.Retrieve(id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve secret %s for env injection: %w", id, err)
		}
		
		// Verify it's an environment type secret
		if secret.Type != SecretTypeEnv {
			clearBytes(secret.Value)
			return nil, fmt.Errorf("secret %s is type %s, not env", id, secret.Type)
		}
		
		if secret.Target == "" {
			clearBytes(secret.Value)
			return nil, fmt.Errorf("secret %s has no target environment variable name", id)
		}
		
		// Register with sanitizer
		if err := i.sanitizer.RegisterSecret(id, secret.Value); err != nil {
			log.Printf("Warning: failed to register secret %s with sanitizer: %v", id, err)
		}
		
		// Add environment variable
		envVar := fmt.Sprintf("%s=%s", secret.Target, string(secret.Value))
		result = append(result, envVar)
		
		// Audit log
		log.Printf("Injected environment variable: %s (from secret: %s)", secret.Target, id)
		
		// Clear secret value from memory
		clearBytes(secret.Value)
	}
	
	return result, nil
}

// ValidateSecretType verifies that a secret is of the expected type
func (i *Injector) ValidateSecretType(secretID string, expectedType SecretType) error {
	secret, err := i.vault.Retrieve(secretID)
	if err != nil {
		return fmt.Errorf("failed to retrieve secret %s for validation: %w", secretID, err)
	}
	
	defer clearBytes(secret.Value)
	
	if secret.Type != expectedType {
		return fmt.Errorf("secret %s is type %s, expected %s", secretID, secret.Type, expectedType)
	}
	
	return nil
}