package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// FileAuditLogger implements the AuditLogger interface with file-based persistence.
type FileAuditLogger struct {
	mu          sync.Mutex
	logPath     string
	file        *os.File
	writer      *bufio.Writer
	flushBuffer []AuditRecord
	bufferSize  int
	autoFlush   bool
	flushTicker *time.Ticker
	stopCh      chan struct{}
}

// AuditLoggerConfig contains configuration options for the audit logger.
type AuditLoggerConfig struct {
	LogPath       string
	BufferSize    int
	AutoFlush     bool
	FlushInterval time.Duration
}

// DefaultConfig returns a sensible default configuration for the audit logger.
func DefaultConfig() AuditLoggerConfig {
	return AuditLoggerConfig{
		LogPath:       "/tmp/cert-audit.log",
		BufferSize:    100,
		AutoFlush:     true,
		FlushInterval: 30 * time.Second,
	}
}

// NewFileAuditLogger creates a new file-based audit logger with the given configuration.
func NewFileAuditLogger(config AuditLoggerConfig) (*FileAuditLogger, error) {
	if config.LogPath == "" {
		return nil, fmt.Errorf("log path cannot be empty")
	}
	if config.BufferSize <= 0 {
		config.BufferSize = 100
	}
	if config.FlushInterval <= 0 {
		config.FlushInterval = 30 * time.Second
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(config.LogPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open or create log file
	file, err := os.OpenFile(config.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	al := &FileAuditLogger{
		logPath:     config.LogPath,
		file:        file,
		writer:      bufio.NewWriter(file),
		flushBuffer: make([]AuditRecord, 0, config.BufferSize),
		bufferSize:  config.BufferSize,
		autoFlush:   config.AutoFlush,
		stopCh:      make(chan struct{}),
	}

	// Start auto-flush if enabled
	if config.AutoFlush {
		al.flushTicker = time.NewTicker(config.FlushInterval)
		go al.autoFlushRoutine()
	}

	return al, nil
}

// LogValidation logs a certificate validation operation.
func (al *FileAuditLogger) LogValidation(result *certs.ValidationResult, metadata map[string]interface{}) error {
	level := LevelInfo
	message := "Certificate validation completed successfully"

	if !result.Valid || len(result.Errors) > 0 {
		level = LevelWarn
		message = fmt.Sprintf("Certificate validation failed with %d errors", len(result.Errors))
	}

	record := AuditRecord{
		ID:               uuid.New().String(),
		Timestamp:        time.Now().UTC(),
		EventType:        EventTypeValidation,
		Level:            level,
		Message:          message,
		ValidationResult: result,
		Metadata:         metadata,
	}

	// Add correlation ID if present in metadata
	if corrID, ok := metadata["correlation_id"].(string); ok {
		record.CorrelationID = corrID
	}

	return al.writeRecord(record)
}

// LogFailure logs a system failure during certificate operations.
func (al *FileAuditLogger) LogFailure(err error, context map[string]interface{}) error {
	auditErr := &AuditError{
		Type:    "system_error",
		Message: err.Error(),
		Details: context,
	}

	record := AuditRecord{
		ID:        uuid.New().String(),
		Timestamp: time.Now().UTC(),
		EventType: EventTypeFailure,
		Level:     LevelError,
		Message:   fmt.Sprintf("System failure: %s", err.Error()),
		Error:     auditErr,
		Metadata:  context,
	}

	// Add correlation ID if present in context
	if corrID, ok := context["correlation_id"].(string); ok {
		record.CorrelationID = corrID
	}

	return al.writeRecord(record)
}

// LogEvent logs a general audit event.
func (al *FileAuditLogger) LogEvent(eventType string, message string, metadata map[string]interface{}) error {
	record := AuditRecord{
		ID:        uuid.New().String(),
		Timestamp: time.Now().UTC(),
		EventType: EventTypeSystemEvent,
		Level:     LevelInfo,
		Message:   message,
		Metadata:  metadata,
	}

	// Add correlation ID if present in metadata
	if corrID, ok := metadata["correlation_id"].(string); ok {
		record.CorrelationID = corrID
	}

	return al.writeRecord(record)
}

// writeRecord writes a single audit record to the buffer or file.
func (al *FileAuditLogger) writeRecord(record AuditRecord) error {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Add to buffer
	al.flushBuffer = append(al.flushBuffer, record)

	// Check if we need to flush
	if len(al.flushBuffer) >= al.bufferSize {
		return al.flushUnsafe()
	}

	return nil
}

// Flush writes all buffered audit records to persistent storage.
func (al *FileAuditLogger) Flush() error {
	al.mu.Lock()
	defer al.mu.Unlock()
	return al.flushUnsafe()
}

// flushUnsafe performs the actual flush operation without locking.
func (al *FileAuditLogger) flushUnsafe() error {
	if len(al.flushBuffer) == 0 {
		return nil
	}

	// Write all buffered records
	for _, record := range al.flushBuffer {
		jsonBytes, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal audit record: %w", err)
		}

		if _, err := al.writer.Write(jsonBytes); err != nil {
			return fmt.Errorf("failed to write audit record: %w", err)
		}

		if _, err := al.writer.WriteString("\n"); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	// Flush the writer
	if err := al.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	// Sync to disk
	if err := al.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	// Clear the buffer
	al.flushBuffer = al.flushBuffer[:0]

	return nil
}

// autoFlushRoutine runs in a goroutine to automatically flush the buffer at regular intervals.
func (al *FileAuditLogger) autoFlushRoutine() {
	for {
		select {
		case <-al.flushTicker.C:
			al.Flush() // Ignore errors in auto-flush
		case <-al.stopCh:
			return
		}
	}
}

// Close gracefully shuts down the audit logger.
func (al *FileAuditLogger) Close() error {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Stop auto-flush
	if al.flushTicker != nil {
		al.flushTicker.Stop()
		close(al.stopCh)
	}

	// Flush remaining records
	if err := al.flushUnsafe(); err != nil {
		return fmt.Errorf("final flush failed: %w", err)
	}

	// Close file
	if err := al.file.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %w", err)
	}

	return nil
}