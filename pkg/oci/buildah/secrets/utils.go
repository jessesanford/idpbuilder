package secrets

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// cleanupRoutine runs periodically to clean up expired secret files
func (sm *SecretManager) cleanupRoutine() {
	for {
		select {
		case <-sm.cleanupTicker.C:
			sm.cleanupExpiredSecrets()
		case <-sm.ctx.Done():
			sm.logger.Debug("Secret manager cleanup routine stopping")
			return
		}
	}
}

// cleanupExpiredSecrets removes expired temporary secret files
func (sm *SecretManager) cleanupExpiredSecrets() {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	now := time.Now()
	expiredSecrets := make([]string, 0)
	for secretID, secretFile := range sm.secretFiles {
		if now.Sub(secretFile.CreatedAt) > secretFile.TTL {
			if err := os.Remove(secretFile.Path); err != nil {
				sm.logger.WithFields(logrus.Fields{"secret_id": secretID, "path": redactPath(secretFile.Path), "error": err}).Warn("Failed to remove expired secret file")
			} else {
				sm.logger.WithFields(logrus.Fields{"secret_id": secretID, "age": now.Sub(secretFile.CreatedAt).String()}).Debug("Cleaned up expired secret file")
			}
			expiredSecrets = append(expiredSecrets, secretID)
		}
	}
	for _, secretID := range expiredSecrets {
		delete(sm.secretFiles, secretID)
	}
	if len(expiredSecrets) > 0 {
		sm.logger.WithField("count", len(expiredSecrets)).Debug("Cleaned up expired secrets")
	}
}

// RedactLogMessage sanitizes log messages to remove secrets
func (sm *SecretManager) RedactLogMessage(message string) string {
	redacted := message
	for _, pattern := range secretPatterns {
		redacted = pattern.ReplaceAllString(redacted, "${1}[REDACTED]")
	}
	commonSecrets := []string{"password=", "passwd=", "secret=", "token=", "key=", "Password=", "Secret=", "Token=", "Key=", "PASSWORD=", "SECRET=", "TOKEN=", "KEY="}
	for _, indicator := range commonSecrets {
		if idx := strings.Index(strings.ToLower(redacted), strings.ToLower(indicator)); idx != -1 {
			start := idx + len(indicator)
			end := start
			for end < len(redacted) && redacted[end] != ' ' && redacted[end] != '"' && redacted[end] != '\'' && redacted[end] != '\n' && redacted[end] != '\t' {
				end++
			}
			if end > start {
				redacted = redacted[:start] + "[REDACTED]" + redacted[end:]
			}
		}
	}
	return redacted
}

// SanitizeEnvVars removes secrets from environment variable maps
func (sm *SecretManager) SanitizeEnvVars(envVars map[string]string) map[string]string {
	sanitized := make(map[string]string)
	for key, value := range envVars {
		if sm.isSecretArg(key, value) {
			sanitized[key] = "[REDACTED]"
		} else {
			sanitized[key] = value
		}
	}
	return sanitized
}

// GetStats returns statistics about the secret manager
func (sm *SecretManager) GetStats() map[string]interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	stats := make(map[string]interface{})
	stats["active_secrets"] = len(sm.secretFiles)
	stats["temp_directory"] = redactPath(sm.tempDir)
	now := time.Now()
	var young, old int
	for _, secretFile := range sm.secretFiles {
		age := now.Sub(secretFile.CreatedAt)
		if age < 5*time.Minute {
			young++
		} else {
			old++
		}
	}
	stats["secrets_under_5min"] = young
	stats["secrets_over_5min"] = old
	return stats
}