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

func (i *Injector) InjectBuildArgs(ctx context.Context, args map[string]string, secretIDs []string) (map[string]string, error) {
	if args == nil {
		args = make(map[string]string)
	}
	result := make(map[string]string)
	for k, v := range args {
		result[k] = v
	}
	
	for _, id := range secretIDs {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}
		
		secret, err := i.vault.Retrieve(id)
		if err != nil {
			return nil, fmt.Errorf("retrieve secret %s: %w", id, err)
		}
		if secret.Type != SecretTypeBuildArg {
			return nil, fmt.Errorf("secret %s not build-arg type", id)
		}
		if secret.Target == "" {
			return nil, fmt.Errorf("secret %s missing target", id)
		}
		
		i.sanitizer.RegisterSecret(id, secret.Value)
		result[secret.Target] = string(secret.Value)
		log.Printf("Injected build arg: %s", secret.Target)
		clearBytes(secret.Value)
	}
	return result, nil
}

func (i *Injector) PrepareSecretMount(ctx context.Context, secretID string) (string, func(), error) {
	secret, err := i.vault.Retrieve(secretID)
	if err != nil {
		return "", nil, fmt.Errorf("retrieve secret %s: %w", secretID, err)
	}
	if secret.Type != SecretTypeMount {
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("secret %s not mount type", secretID)
	}
	
	tmpDir, err := os.MkdirTemp("", "oci-secret-mount-*")
	if err != nil {
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("create temp dir: %w", err)
	}
	
	if err := os.Chmod(tmpDir, 0700); err != nil {
		os.RemoveAll(tmpDir)
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("secure temp dir: %w", err)
	}
	
	tmpFile := filepath.Join(tmpDir, "secret")
	if secret.Target != "" {
		tmpFile = filepath.Join(tmpDir, filepath.Base(secret.Target))
	}
	
	mode := secret.Mode
	if mode == 0 {
		mode = 0600
	}
	
	if err := os.WriteFile(tmpFile, secret.Value, mode); err != nil {
		os.RemoveAll(tmpDir)
		clearBytes(secret.Value)
		return "", nil, fmt.Errorf("write secret file: %w", err)
	}
	
	if secret.UID > 0 || secret.GID > 0 {
		os.Chown(tmpFile, secret.UID, secret.GID)
	}
	
	i.sanitizer.RegisterSecret(secretID, secret.Value)
	clearBytes(secret.Value)
	
	cleanup := func() {
		if data, err := os.ReadFile(tmpFile); err == nil {
			clearBytes(data)
			os.WriteFile(tmpFile, data, mode)
		}
		os.RemoveAll(tmpDir)
		i.sanitizer.UnregisterSecret(secretID)
		log.Printf("Cleaned up secret mount: %s", secretID)
	}
	
	log.Printf("Created secret mount: %s -> %s", secretID, tmpFile)
	return tmpFile, cleanup, nil
}

func (i *Injector) InjectEnvironmentSecrets(ctx context.Context, env []string, secretIDs []string) ([]string, error) {
	if env == nil {
		env = []string{}
	}
	result := make([]string, len(env))
	copy(result, env)
	
	for _, id := range secretIDs {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
		default:
		}
		
		secret, err := i.vault.Retrieve(id)
		if err != nil {
			return nil, fmt.Errorf("retrieve secret %s: %w", id, err)
		}
		if secret.Type != SecretTypeEnv {
			clearBytes(secret.Value)
			return nil, fmt.Errorf("secret %s not env type", id)
		}
		if secret.Target == "" {
			clearBytes(secret.Value)
			return nil, fmt.Errorf("secret %s missing target", id)
		}
		
		i.sanitizer.RegisterSecret(id, secret.Value)
		result = append(result, fmt.Sprintf("%s=%s", secret.Target, string(secret.Value)))
		log.Printf("Injected env var: %s", secret.Target)
		clearBytes(secret.Value)
	}
	return result, nil
}

func (i *Injector) ValidateSecretType(secretID string, expectedType SecretType) error {
	secret, err := i.vault.Retrieve(secretID)
	if err != nil {
		return fmt.Errorf("retrieve secret %s: %w", secretID, err)
	}
	defer clearBytes(secret.Value)
	if secret.Type != expectedType {
		return fmt.Errorf("secret %s type mismatch", secretID)
	}
	return nil
}