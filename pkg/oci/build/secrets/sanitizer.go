package secrets

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
)

type Sanitizer struct {
	patterns map[string]*regexp.Regexp
	masks    map[string]string
	mu       sync.RWMutex
}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		patterns: make(map[string]*regexp.Regexp),
		masks:    make(map[string]string),
	}
}

func (s *Sanitizer) RegisterSecret(id string, value []byte) error {
	if id == "" || len(value) == 0 {
		return fmt.Errorf("empty secret ID or value")
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	pattern := regexp.QuoteMeta(string(value))
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("compile pattern %s: %w", id, err)
	}
	
	s.patterns[id] = regex
	s.masks[id] = "***" + strings.ToUpper(id) + "***"
	return nil
}

func (s *Sanitizer) UnregisterSecret(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.patterns, id)
	delete(s.masks, id)
}

func (s *Sanitizer) Sanitize(input string) string {
	if input == "" {
		return input
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	output := input
	for id, pattern := range s.patterns {
		output = pattern.ReplaceAllString(output, s.masks[id])
	}
	return s.sanitizeCommonPatterns(output)
}

func (s *Sanitizer) sanitizeCommonPatterns(input string) string {
	patterns := []struct {
		regex *regexp.Regexp
		mask  string
	}{
		{regexp.MustCompile(`(?i)(password|passwd|pwd)\s*[=:]\s*\S+`), "$1=***REDACTED***"},
		{regexp.MustCompile(`(?i)(token|api[_-]?key)\s*[=:]\s*\S+`), "$1=***REDACTED***"},
		{regexp.MustCompile(`(?i)(secret)\s*[=:]\s*\S+`), "$1=***REDACTED***"},
		{regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9\-\._~\+\/]+=*`), "Bearer ***TOKEN***"},
		{regexp.MustCompile(`(?i)(mongodb|postgres|mysql)://[^:]+:[^@]+@`), "$1://***USER***:***PASS***@"},
		{regexp.MustCompile(`-----BEGIN [A-Z ]+PRIVATE KEY-----[\s\S]*?-----END [A-Z ]+PRIVATE KEY-----`), "***PRIVATE_KEY***"},
	}
	
	output := input
	for _, p := range patterns {
		output = p.regex.ReplaceAllString(output, p.mask)
	}
	return output
}

type SanitizeWriter struct {
	sanitizer *Sanitizer
	writer    io.Writer
}

func NewSanitizeWriter(sanitizer *Sanitizer, writer io.Writer) *SanitizeWriter {
	return &SanitizeWriter{sanitizer: sanitizer, writer: writer}
}

func (w *SanitizeWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	
	sanitized := w.sanitizer.Sanitize(string(p))
	_, err = w.writer.Write([]byte(sanitized))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *Sanitizer) SanitizeBytes(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	return []byte(s.Sanitize(string(input)))
}

func (s *Sanitizer) GetRegisteredSecrets() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	secrets := make([]string, 0, len(s.patterns))
	for id := range s.patterns {
		secrets = append(secrets, id)
	}
	return secrets
}