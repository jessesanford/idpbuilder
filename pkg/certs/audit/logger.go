package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AuditLogger provides persistent audit logging for certificate validation decisions
type AuditLogger struct {
	file    *os.File
	maxSize int64
	mu      sync.Mutex
}

// AuditEntry represents a single audit log entry
type AuditEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Action      string    `json:"action"`
	Certificate string    `json:"certificate"`
	Decision    string    `json:"decision"`
	Reason      string    `json:"reason"`
	UserID      string    `json:"user_id"`
	Hostname    string    `json:"hostname,omitempty"`
	Details     string    `json:"details,omitempty"`
}

// AuditLogConfig holds configuration for audit logging
type AuditLogConfig struct {
	LogPath    string
	MaxSizeMB  int64
	CreateDirs bool
}

// DefaultConfig returns a default audit log configuration
func DefaultConfig() *AuditLogConfig {
	return &AuditLogConfig{
		LogPath:    "/var/log/idpbuilder/certificate-audit.log",
		MaxSizeMB:  10, // 10MB max file size
		CreateDirs: true,
	}
}

// NewAuditLogger creates a new audit logger with the specified configuration
func NewAuditLogger(config *AuditLogConfig) (*AuditLogger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Create directory if needed
	if config.CreateDirs {
		dir := filepath.Dir(config.LogPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create audit log directory %s: %w", dir, err)
		}
	}

	// Open or create log file
	file, err := os.OpenFile(config.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file %s: %w", config.LogPath, err)
	}

	return &AuditLogger{
		file:    file,
		maxSize: config.MaxSizeMB * 1024 * 1024, // Convert MB to bytes
	}, nil
}

// LogCertificateValidation logs certificate validation decisions
func (a *AuditLogger) LogCertificateValidation(cert, hostname, decision, reason, userID string) error {
	entry := AuditEntry{
		Timestamp:   time.Now(),
		Action:      "CERTIFICATE_VALIDATION",
		Certificate: cert,
		Decision:    decision,
		Reason:      reason,
		UserID:      userID,
		Hostname:    hostname,
	}
	
	return a.writeEntry(entry)
}

// LogTrustDecision logs trust anchor decisions
func (a *AuditLogger) LogTrustDecision(cert, decision, reason, userID string) error {
	entry := AuditEntry{
		Timestamp:   time.Now(),
		Action:      "TRUST_DECISION",
		Certificate: cert,
		Decision:    decision,
		Reason:      reason,
		UserID:      userID,
	}
	
	return a.writeEntry(entry)
}

// LogFallbackActivation logs when fallback mechanisms are activated
func (a *AuditLogger) LogFallbackActivation(cert, fallbackType, reason, userID string) error {
	entry := AuditEntry{
		Timestamp:   time.Now(),
		Action:      "FALLBACK_ACTIVATION",
		Certificate: cert,
		Decision:    fallbackType,
		Reason:      reason,
		UserID:      userID,
	}
	
	return a.writeEntry(entry)
}

// LogSecurityOverride logs when security policies are overridden
func (a *AuditLogger) LogSecurityOverride(cert, overrideType, justification, userID string) error {
	entry := AuditEntry{
		Timestamp:   time.Now(),
		Action:      "SECURITY_OVERRIDE",
		Certificate: cert,
		Decision:    overrideType,
		Reason:      justification,
		UserID:      userID,
	}
	
	return a.writeEntry(entry)
}

// LogError logs certificate validation errors
func (a *AuditLogger) LogError(cert, errorType, details, userID string) error {
	entry := AuditEntry{
		Timestamp:   time.Now(),
		Action:      "VALIDATION_ERROR",
		Certificate: cert,
		Decision:    "ERROR",
		Reason:      errorType,
		UserID:      userID,
		Details:     details,
	}
	
	return a.writeEntry(entry)
}

// writeEntry writes an audit entry to the log file with rotation check
func (a *AuditLogger) writeEntry(entry AuditEntry) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if rotation is needed
	if err := a.checkAndRotate(); err != nil {
		return fmt.Errorf("failed to rotate audit log: %w", err)
	}

	// Marshal entry to JSON
	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal audit entry: %w", err)
	}

	// Write entry with newline
	_, err = a.file.WriteString(string(jsonData) + "\n")
	if err != nil {
		return fmt.Errorf("failed to write audit entry: %w", err)
	}

	// Ensure data is written to disk
	return a.file.Sync()
}

// checkAndRotate checks file size and rotates if necessary
func (a *AuditLogger) checkAndRotate() error {
	// Get current file info
	info, err := a.file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Check if rotation is needed
	if info.Size() < a.maxSize {
		return nil // No rotation needed
	}

	// Get current file path
	currentPath := a.file.Name()
	
	// Close current file
	if err := a.file.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	// Create rotated filename with timestamp
	rotatedPath := fmt.Sprintf("%s.%s", currentPath, time.Now().Format("20060102-150405"))
	
	// Rename current file
	if err := os.Rename(currentPath, rotatedPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Open new file
	newFile, err := os.OpenFile(currentPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	a.file = newFile
	return nil
}

// Close closes the audit log file
func (a *AuditLogger) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	
	if a.file != nil {
		return a.file.Close()
	}
	return nil
}

// ReadEntries reads audit entries from the log file
func (a *AuditLogger) ReadEntries(limit int) ([]AuditEntry, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Get current file path
	filePath := a.file.Name()
	
	// Open file for reading
	readFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log for reading: %w", err)
	}
	defer readFile.Close()

	var entries []AuditEntry
	decoder := json.NewDecoder(readFile)
	
	// Read entries line by line
	for {
		var entry AuditEntry
		if err := decoder.Decode(&entry); err != nil {
			break // End of file or invalid JSON
		}
		
		entries = append(entries, entry)
		
		// Apply limit if specified
		if limit > 0 && len(entries) >= limit {
			break
		}
	}

	// Return entries in reverse order (most recent first)
	for i := 0; i < len(entries)/2; i++ {
		entries[i], entries[len(entries)-1-i] = entries[len(entries)-1-i], entries[i]
	}

	return entries, nil
}

// GetRecentEntries returns recent audit entries within the specified duration
func (a *AuditLogger) GetRecentEntries(since time.Duration) ([]AuditEntry, error) {
	allEntries, err := a.ReadEntries(0) // Read all entries
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-since)
	var recentEntries []AuditEntry
	
	for _, entry := range allEntries {
		if entry.Timestamp.After(cutoff) {
			recentEntries = append(recentEntries, entry)
		}
	}

	return recentEntries, nil
}