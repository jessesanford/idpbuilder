package secrets

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MountSecret creates a secret mount for Buildah
func (sm *SecretManager) MountSecret(mount *SecretMount) (string, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	switch mount.Type {
	case SecretTypeFile:
		return sm.mountFileSecret(mount)
	case SecretTypeKubernetes:
		return sm.mountKubernetesSecret(mount)
	case SecretTypeEnv:
		return sm.mountEnvSecret(mount)
	case SecretTypeVault:
		return sm.mountVaultSecret(mount)
	default:
		return "", fmt.Errorf("unsupported secret type: %s", mount.Type)
	}
}

func (sm *SecretManager) mountFileSecret(mount *SecretMount) (string, error) {
	if _, err := os.Stat(mount.Source); err != nil {
		return "", fmt.Errorf("source file not accessible: %w", err)
	}
	secretID := generateSecretID()
	targetPath := filepath.Join(sm.tempDir, fmt.Sprintf("mount-%s", secretID))
	data, err := os.ReadFile(mount.Source)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}
	if err := os.WriteFile(targetPath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to create temp secret file: %w", err)
	}
	sm.secretFiles[secretID] = &secretFile{Path: targetPath, CreatedAt: time.Now(), TTL: 30 * time.Minute}
	sm.logger.WithFields(logrus.Fields{"secret_id": secretID, "source": redactPath(mount.Source), "target": mount.Target}).Debug("File secret mounted")
	return targetPath, nil
}

func (sm *SecretManager) mountKubernetesSecret(mount *SecretMount) (string, error) {
	if sm.k8sClient == nil {
		return "", fmt.Errorf("kubernetes client not configured")
	}
	parts := strings.Split(mount.Source, "/")
	secretName, key := parts[0], ""
	if len(parts) > 1 {
		key = parts[1]
	}
	secret, err := sm.k8sClient.CoreV1().Secrets(sm.namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get secret %s: %w", secretName, err)
	}
	var data []byte
	if key != "" {
		var exists bool
		if data, exists = secret.Data[key]; !exists {
			return "", fmt.Errorf("key %s not found in secret %s", key, secretName)
		}
	} else {
		if len(secret.Data) == 1 {
			for _, v := range secret.Data {
				data = v
				break
			}
		} else if len(secret.Data) > 1 {
			return "", fmt.Errorf("secret %s has multiple keys, specify key with secretname/key", secretName)
		} else {
			return "", fmt.Errorf("secret %s has no data", secretName)
		}
	}
	secretID := generateSecretID()
	targetPath := filepath.Join(sm.tempDir, fmt.Sprintf("k8s-%s", secretID))
	if err := os.WriteFile(targetPath, data, 0600); err != nil {
		return "", fmt.Errorf("failed to create temp secret file: %w", err)
	}
	sm.secretFiles[secretID] = &secretFile{Path: targetPath, CreatedAt: time.Now(), TTL: 30 * time.Minute}
	sm.logger.WithFields(logrus.Fields{"secret_id": secretID, "k8s_secret": secretName, "key": key, "target": mount.Target}).Debug("Kubernetes secret mounted")
	return targetPath, nil
}

func (sm *SecretManager) mountEnvSecret(mount *SecretMount) (string, error) {
	value := os.Getenv(mount.Source)
	if value == "" {
		return "", fmt.Errorf("environment variable %s not found or empty", mount.Source)
	}
	secretID := generateSecretID()
	targetPath := filepath.Join(sm.tempDir, fmt.Sprintf("env-%s", secretID))
	if err := os.WriteFile(targetPath, []byte(value), 0600); err != nil {
		return "", fmt.Errorf("failed to create temp secret file: %w", err)
	}
	sm.secretFiles[secretID] = &secretFile{Path: targetPath, CreatedAt: time.Now(), TTL: 30 * time.Minute}
	sm.logger.WithFields(logrus.Fields{"secret_id": secretID, "env_var": mount.Source, "target": mount.Target}).Debug("Environment secret mounted")
	return targetPath, nil
}

func (sm *SecretManager) mountVaultSecret(mount *SecretMount) (string, error) {
	return "", fmt.Errorf("vault secret mounting not yet implemented - requires external vault client integration")
}

// CleanupSecret manually removes a specific secret
func (sm *SecretManager) CleanupSecret(secretID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	secretFile, exists := sm.secretFiles[secretID]
	if !exists {
		return nil
	}
	if err := os.Remove(secretFile.Path); err != nil {
		sm.logger.WithFields(logrus.Fields{"secret_id": secretID, "error": err}).Warn("Failed to remove secret file")
		return err
	}
	delete(sm.secretFiles, secretID)
	sm.logger.WithField("secret_id", secretID).Debug("Manually cleaned up secret")
	return nil
}

// Cleanup removes all secret files and stops the manager
func (sm *SecretManager) Cleanup() error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.cleanupTicker.Stop()
	sm.cancel()
	var lastErr error
	for secretID, secretFile := range sm.secretFiles {
		if err := os.Remove(secretFile.Path); err != nil {
			sm.logger.WithFields(logrus.Fields{"secret_id": secretID, "error": err}).Warn("Failed to remove secret file during cleanup")
			lastErr = err
		}
	}
	sm.secretFiles = make(map[string]*secretFile)
	if err := os.RemoveAll(sm.tempDir); err != nil {
		sm.logger.WithFields(logrus.Fields{"temp_dir": sm.tempDir, "error": err}).Warn("Failed to remove temp directory")
		lastErr = err
	}
	sm.logger.Debug("Secret manager cleanup completed")
	return lastErr
}

// redactPath removes sensitive information from file paths for logging
func redactPath(path string) string {
	for _, pattern := range secretPatterns {
		path = pattern.ReplaceAllString(path, "${1}[REDACTED]")
	}
	return path
}