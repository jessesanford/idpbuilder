package audit

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, "/tmp/cert-audit.log", config.LogPath)
	assert.Equal(t, 100, config.BufferSize)
	assert.True(t, config.AutoFlush)
	assert.Equal(t, 30*time.Second, config.FlushInterval)
}

func TestNewFileAuditLogger(t *testing.T) {
	tempDir := t.TempDir()
	config := AuditLoggerConfig{
		LogPath:    filepath.Join(tempDir, "test-audit.log"),
		BufferSize: 10,
		AutoFlush:  false,
	}

	logger, err := NewFileAuditLogger(config)
	require.NoError(t, err)
	assert.NotNil(t, logger)
	
	// Verify config applied
	assert.Equal(t, 10, logger.bufferSize)

	logger.Close()
}

func TestLogValidation(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test-validation.log")

	config := AuditLoggerConfig{
		LogPath:     logPath,
		BufferSize:  1, // Force immediate flush
		AutoFlush:   false,
	}

	logger, err := NewFileAuditLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	// Test successful validation
	result := &certs.ValidationResult{
		Valid:          true,
		Errors:         []certs.ValidationError{},
		ValidationTime: time.Now(),
	}

	metadata := map[string]interface{}{
		"request_id":     "test-123",
		"correlation_id": "corr-456",
	}

	err = logger.LogValidation(result, metadata)
	assert.NoError(t, err)

	// Read and verify the log entry
	records := readLogRecords(t, logPath)
	require.Len(t, records, 1)

	record := records[0]
	assert.Equal(t, EventTypeValidation, record.EventType)
	assert.Equal(t, LevelInfo, record.Level)
	assert.Equal(t, "corr-456", record.CorrelationID)
}

func TestLogFailure(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test-failure.log")

	config := AuditLoggerConfig{
		LogPath:     logPath,
		BufferSize:  1,
		AutoFlush:   false,
	}

	logger, err := NewFileAuditLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	testErr := errors.New("test system error")
	context := map[string]interface{}{
		"operation":      "certificate_load",
		"correlation_id": "error-123",
	}

	err = logger.LogFailure(testErr, context)
	assert.NoError(t, err)

	// Read and verify the log entry
	records := readLogRecords(t, logPath)
	require.Len(t, records, 1)

	record := records[0]
	assert.Equal(t, EventTypeFailure, record.EventType)
	assert.Equal(t, LevelError, record.Level)
	assert.Contains(t, record.Message, "test system error")
	assert.Equal(t, "error-123", record.CorrelationID)
}

func TestLogEvent(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test-event.log")

	config := AuditLoggerConfig{
		LogPath:     logPath,
		BufferSize:  1,
		AutoFlush:   false,
	}

	logger, err := NewFileAuditLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	eventType := "trust_anchor_update"
	message := "Trust anchor configuration updated"
	metadata := map[string]interface{}{
		"admin_user":     "test@example.com",
		"correlation_id": "event-789",
	}

	err = logger.LogEvent(eventType, message, metadata)
	assert.NoError(t, err)

	// Read and verify the log entry
	records := readLogRecords(t, logPath)
	require.Len(t, records, 1)

	record := records[0]
	assert.Equal(t, EventTypeSystemEvent, record.EventType)
	assert.Equal(t, LevelInfo, record.Level)
	assert.Equal(t, message, record.Message)
	assert.Equal(t, "event-789", record.CorrelationID)
}

func TestBuffering(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test-buffering.log")

	config := AuditLoggerConfig{
		LogPath:     logPath,
		BufferSize:  3, // Buffer 3 records before flushing
		AutoFlush:   false,
	}

	logger, err := NewFileAuditLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	// Write 2 records - should be buffered
	logger.LogEvent("event1", "Test event 1", map[string]interface{}{})
	logger.LogEvent("event2", "Test event 2", map[string]interface{}{})

	// File should be empty (records are buffered)
	records := readLogRecords(t, logPath)
	assert.Len(t, records, 0)

	// Write 3rd record - should trigger flush
	logger.LogEvent("event3", "Test event 3", map[string]interface{}{})

	// Now file should contain all 3 records
	records = readLogRecords(t, logPath)
	assert.Len(t, records, 3)
}

func TestManualFlush(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test-flush.log")

	config := AuditLoggerConfig{
		LogPath:     logPath,
		BufferSize:  10,
		AutoFlush:   false,
	}

	logger, err := NewFileAuditLogger(config)
	require.NoError(t, err)
	defer logger.Close()

	// Write records
	logger.LogEvent("event1", "Test event 1", map[string]interface{}{})
	logger.LogEvent("event2", "Test event 2", map[string]interface{}{})

	// File should be empty (records are buffered)
	records := readLogRecords(t, logPath)
	assert.Len(t, records, 0)

	// Manual flush
	err = logger.Flush()
	assert.NoError(t, err)

	// Now file should contain records
	records = readLogRecords(t, logPath)
	assert.Len(t, records, 2)
}

func TestClose(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "test-close.log")

	config := AuditLoggerConfig{
		LogPath:     logPath,
		BufferSize:  10,
		AutoFlush:   true,
		FlushInterval: 1 * time.Second,
	}

	logger, err := NewFileAuditLogger(config)
	require.NoError(t, err)

	// Write some records
	for i := 0; i < 3; i++ {
		logger.LogEvent("close_test", "Record", map[string]interface{}{})
	}

	// Close the logger
	err = logger.Close()
	assert.NoError(t, err)

	// Verify all records were flushed
	records := readLogRecords(t, logPath)
	assert.Len(t, records, 3)
}

// Helper function to read and parse log records from a file
func readLogRecords(t *testing.T, logPath string) []AuditRecord {
	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []AuditRecord{}
		}
		t.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	var records []AuditRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var record AuditRecord
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("Failed to unmarshal log record: %v", err)
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading log file: %v", err)
	}

	return records
}