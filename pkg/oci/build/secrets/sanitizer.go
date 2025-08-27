package secrets

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

// Sanitizer prevents secret leakage in logs and outputs by detecting and masking secrets
type Sanitizer struct {
	patterns map[string]*regexp.Regexp  // Compiled regex patterns for each secret
	masks    map[string]string          // Replacement masks for each secret
	mu       sync.RWMutex               // Protects concurrent access to patterns
}

// NewSanitizer creates a new log sanitizer
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		patterns: make(map[string]*regexp.Regexp),
		masks:    make(map[string]string),
	}
}

// RegisterSecret registers a secret value to be sanitized from outputs
func (s *Sanitizer) RegisterSecret(id string, value []byte) error {
	if id == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}
	
	if len(value) == 0 {
		return fmt.Errorf("secret value cannot be empty")
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Create regex pattern for exact matches (escape special characters)
	pattern := regexp.QuoteMeta(string(value))
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("failed to compile pattern for secret %s: %w", id, err)
	}
	
	s.patterns[id] = regex
	s.masks[id] = "***" + strings.ToUpper(id) + "***"
	
	return nil
}

// UnregisterSecret removes a secret from sanitization (e.g., after cleanup)
func (s *Sanitizer) UnregisterSecret(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delete(s.patterns, id)
	delete(s.masks, id)
}

// Sanitize removes all registered secrets from the input string
func (s *Sanitizer) Sanitize(input string) string {
	if input == "" {
		return input
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	output := input
	
	// Replace registered secrets
	for id, pattern := range s.patterns {
		output = pattern.ReplaceAllString(output, s.masks[id])
	}
	
	// Also sanitize common secret patterns
	output = s.sanitizeCommonPatterns(output)
	
	return output
}

// sanitizeCommonPatterns removes common secret patterns even if not explicitly registered
func (s *Sanitizer) sanitizeCommonPatterns(input string) string {
	// Define common patterns that might contain secrets
	patterns := []struct {
		regex *regexp.Regexp
		mask  string
	}{
		// Password patterns
		{regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[=:]\s*\S+`), "$1=***REDACTED***"},
		{regexp.MustCompile(`(?i)--password\s+\S+`), "--password ***REDACTED***"},
		
		// Token and API key patterns
		{regexp.MustCompile(`(?i)(token|api[_-]?key)\s*[=:]\s*\S+`), "$1=***REDACTED***"},
		{regexp.MustCompile(`(?i)--api-key\s+\S+`), "--api-key ***REDACTED***"},
		
		// Generic secret patterns
		{regexp.MustCompile(`(?i)(secret)\s*[=:]\s*\S+`), "$1=***REDACTED***"},
		{regexp.MustCompile(`(?i)--secret\s+\S+`), "--secret ***REDACTED***"},
		
		// Bearer token patterns
		{regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9\-\._~\+\/]+=*`), "Bearer ***TOKEN***"},
		{regexp.MustCompile(`(?i)Authorization:\s*Bearer\s+[A-Za-z0-9\-\._~\+\/]+=*`), "Authorization: Bearer ***TOKEN***"},
		
		// Connection string patterns
		{regexp.MustCompile(`(?i)(mongodb|postgres|mysql)://[^:]+:[^@]+@`), "$1://***USER***:***PASS***@"},
		
		// Private key patterns
		{regexp.MustCompile(`-----BEGIN [A-Z ]+PRIVATE KEY-----[\s\S]*?-----END [A-Z ]+PRIVATE KEY-----`), "***PRIVATE_KEY***"},
		
		// Certificate patterns
		{regexp.MustCompile(`-----BEGIN CERTIFICATE-----[\s\S]*?-----END CERTIFICATE-----`), "***CERTIFICATE***"},
	}
	
	output := input
	for _, p := range patterns {
		output = p.regex.ReplaceAllString(output, p.mask)
	}
	
	return output
}

// SanitizeWriter wraps an io.Writer to automatically sanitize all written content
type SanitizeWriter struct {
	sanitizer *Sanitizer
	writer    io.Writer
}

// NewSanitizeWriter creates a new sanitizing writer wrapper
func NewSanitizeWriter(sanitizer *Sanitizer, writer io.Writer) *SanitizeWriter {
	return &SanitizeWriter{
		sanitizer: sanitizer,
		writer:    writer,
	}
}

// Write implements io.Writer by sanitizing content before writing
func (w *SanitizeWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	
	// Sanitize the content
	sanitized := w.sanitizer.Sanitize(string(p))
	sanitizedBytes := []byte(sanitized)
	
	// Write sanitized content
	_, err = w.writer.Write(sanitizedBytes)
	if err != nil {
		return 0, err
	}
	
	// Return original length to maintain write semantics
	return len(p), nil
}

// SanitizeBytes is a convenience function to sanitize byte slices
func (s *Sanitizer) SanitizeBytes(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	
	sanitized := s.Sanitize(string(input))
	return []byte(sanitized)
}

// GetRegisteredSecrets returns the IDs of all registered secrets (for debugging/monitoring)
func (s *Sanitizer) GetRegisteredSecrets() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	secrets := make([]string, 0, len(s.patterns))
	for id := range s.patterns {
		secrets = append(secrets, id)
	}
	
	return secrets
}