package secrets

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
)

// NewSecretManager creates a new SecretManager instance
func NewSecretManager(k8sClient kubernetes.Interface, namespace string, logger *logrus.Logger) (*SecretManager, error) {
	if logger == nil {
		logger = logrus.New()
	}
	tempDir, err := os.MkdirTemp("", "buildah-secrets-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	if err := os.Chmod(tempDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to set temp directory permissions: %w", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	sm := &SecretManager{
		k8sClient: k8sClient, namespace: namespace, tempDir: tempDir, logger: logger,
		secretFiles: make(map[string]*secretFile), ctx: ctx, cancel: cancel,
	}
	sm.cleanupTicker = time.NewTicker(5 * time.Minute)
	go sm.cleanupRoutine()
	return sm, nil
}

// SanitizeBuildArgs processes build arguments and separates secrets
func (sm *SecretManager) SanitizeBuildArgs(args map[string]string) (*BuildArgs, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	buildArgs := &BuildArgs{Args: make(map[string]string), SecretIDs: make([]string, 0)}
	for key, value := range args {
		if sm.isSecretArg(key, value) {
			secretID, err := sm.createSecretFile(key, value)
			if err != nil {
				return nil, fmt.Errorf("failed to create secret for arg %s: %w", key, err)
			}
			buildArgs.SecretIDs = append(buildArgs.SecretIDs, secretID)
			sm.logger.WithField("arg", key).Debug("Build argument marked as secret")
		} else {
			buildArgs.Args[key] = value
		}
	}
	return buildArgs, nil
}

// isSecretArg determines if a build argument should be treated as a secret
func (sm *SecretManager) isSecretArg(key, value string) bool {
	lowerKey := strings.ToLower(key)
	secretKeys := []string{"password", "passwd", "secret", "token", "credential", "api_key", "access_token"}
	for _, secretKey := range secretKeys {
		if strings.Contains(lowerKey, secretKey) {
			return true
		}
	}
	if strings.Contains(lowerKey, "key") && (strings.Contains(lowerKey, "api") || strings.Contains(lowerKey, "secret") || strings.Contains(lowerKey, "access")) {
		return true
	}
	if len(value) > 20 && (strings.Contains(value, ".") || isBase64Like(value)) {
		return true
	}
	return false
}

// createSecretFile creates a temporary file for a secret value
func (sm *SecretManager) createSecretFile(name, value string) (string, error) {
	secretID := generateSecretID()
	filename := fmt.Sprintf("secret-%s", secretID)
	filepath := filepath.Join(sm.tempDir, filename)
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return "", fmt.Errorf("failed to create secret file: %w", err)
	}
	defer file.Close()
	if _, err := file.WriteString(value); err != nil {
		return "", fmt.Errorf("failed to write secret: %w", err)
	}
	sm.secretFiles[secretID] = &secretFile{Path: filepath, CreatedAt: time.Now(), TTL: 30 * time.Minute}
	return secretID, nil
}

// GetSecretPath returns the path to a mounted secret
func (sm *SecretManager) GetSecretPath(secretID string) (string, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	secretFile, exists := sm.secretFiles[secretID]
	if !exists {
		return "", fmt.Errorf("secret %s not found", secretID)
	}
	if _, err := os.Stat(secretFile.Path); err != nil {
		delete(sm.secretFiles, secretID)
		return "", fmt.Errorf("secret file no longer exists")
	}
	return secretFile.Path, nil
}

// generateSecretID creates a cryptographically secure random ID
func generateSecretID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}

// isBase64Like checks if a string looks like base64
func isBase64Like(s string) bool {
	if len(s)%4 != 0 {
		return false
	}
	base64Pattern := regexp.MustCompile(`^[A-Za-z0-9+/]*=*$`)
	return base64Pattern.MatchString(s) && len(s) > 20
}