package audit

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cnoe-io/idpbuilder/pkg/certs"
)

// TestChainValidatorWithAuditIntegration tests the integration between chain validation and audit logging.
func TestChainValidatorWithAuditIntegration(t *testing.T) {
	tempDir := t.TempDir()
	auditLogPath := filepath.Join(tempDir, "integration-audit.log")

	// Set up audit logger
	auditConfig := AuditLoggerConfig{
		LogPath:     auditLogPath,
		BufferSize:  1, // Force immediate flush
		AutoFlush:   false,
	}

	auditLogger, err := NewFileAuditLogger(auditConfig)
	require.NoError(t, err)
	defer auditLogger.Close()

	// Create test certificates
	rootCert, rootKey := createTestRootCA(t)
	leafCert := createTestLeafCert(t, rootCert, rootKey, "example.com")

	// Create chain validator
	validator := &certs.ChainValidatorImpl{}
	config := certs.ChainValidatorConfig{
		TrustAnchors:     []*x509.Certificate{rootCert},
		StrictValidation: true,
		MaxChainLength:   10,
	}
	validator.Configure(config)

	// Test successful validation with audit logging
	t.Run("successful_validation_with_audit", func(t *testing.T) {
		chain := []*x509.Certificate{leafCert, rootCert}
		result, err := validator.ValidateChain(chain)
		require.NoError(t, err)
		assert.True(t, result.Valid)

		// Log the validation result
		metadata := map[string]interface{}{
			"test_name":      "successful_validation_with_audit",
			"hostname":       "example.com",
			"correlation_id": "test-001",
		}

		err = auditLogger.LogValidation(result, metadata)
		assert.NoError(t, err)

		// Verify audit record
		auditRecords := readAuditRecords(t, auditLogPath)
		assert.Len(t, auditRecords, 1)

		record := auditRecords[0]
		assert.Equal(t, EventTypeValidation, record.EventType)
		assert.Equal(t, LevelInfo, record.Level)
		assert.Equal(t, "test-001", record.CorrelationID)
	})

	// Test validation failure with audit logging
	t.Run("validation_failure_with_audit", func(t *testing.T) {
		// Create an expired certificate
		expiredCert := createExpiredTestCert(t, rootCert, rootKey, "expired.com")
		chain := []*x509.Certificate{expiredCert, rootCert}

		result, err := validator.ValidateChain(chain)
		require.NoError(t, err)
		assert.False(t, result.Valid)
		assert.Greater(t, len(result.Errors), 0)

		// Log the validation failure
		metadata := map[string]interface{}{
			"test_name":      "validation_failure_with_audit",
			"correlation_id": "test-002",
		}

		err = auditLogger.LogValidation(result, metadata)
		assert.NoError(t, err)

		// Verify audit record
		auditRecords := readAuditRecords(t, auditLogPath)
		assert.GreaterOrEqual(t, len(auditRecords), 1)

		// Find the failure record
		var failureRecord *AuditRecord
		for i := len(auditRecords) - 1; i >= 0; i-- {
			if auditRecords[i].CorrelationID == "test-002" {
				failureRecord = &auditRecords[i]
				break
			}
		}

		require.NotNil(t, failureRecord)
		assert.Equal(t, LevelWarn, failureRecord.Level)
		assert.False(t, failureRecord.ValidationResult.Valid)
	})
}

// TestAuditLoggerPerformanceWithValidation tests audit logging performance under load.
func TestAuditLoggerPerformanceWithValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	tempDir := t.TempDir()
	auditLogPath := filepath.Join(tempDir, "performance-audit.log")

	auditConfig := AuditLoggerConfig{
		LogPath:    auditLogPath,
		BufferSize: 50,
		AutoFlush:  false,
	}

	auditLogger, err := NewFileAuditLogger(auditConfig)
	require.NoError(t, err)
	defer auditLogger.Close()

	// Create test certificates
	rootCert, rootKey := createTestRootCA(t)
	leafCert := createTestLeafCert(t, rootCert, rootKey, "perf-test.com")

	// Create chain validator
	validator := &certs.ChainValidatorImpl{}
	config := certs.ChainValidatorConfig{
		TrustAnchors:   []*x509.Certificate{rootCert},
		MaxChainLength: 5,
	}
	validator.Configure(config)

	// Run performance test
	numValidations := 100
	chain := []*x509.Certificate{leafCert, rootCert}

	startTime := time.Now()

	for i := 0; i < numValidations; i++ {
		result, err := validator.ValidateChain(chain)
		require.NoError(t, err)

		metadata := map[string]interface{}{
			"iteration":      i,
			"correlation_id": fmt.Sprintf("perf-%d", i),
		}

		err = auditLogger.LogValidation(result, metadata)
		assert.NoError(t, err)
	}

	auditLogger.Flush()
	elapsedTime := time.Since(startTime)

	t.Logf("Performance test: %d validations in %v", numValidations, elapsedTime)

	// Verify all records were written
	auditRecords := readAuditRecords(t, auditLogPath)
	assert.Len(t, auditRecords, numValidations)
}

// Helper functions for test certificate creation

func createTestRootCA(t *testing.T) (*x509.Certificate, *rsa.PrivateKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test CA"},
			Country:      []string{"US"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	return cert, privateKey
}

func createTestLeafCert(t *testing.T, caCert *x509.Certificate, caKey *rsa.PrivateKey, hostname string) *x509.Certificate {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Test Leaf"},
			Country:      []string{"US"},
		},
		DNSNames:    []string{hostname},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(90 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, &privateKey.PublicKey, caKey)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	return cert
}

func createExpiredTestCert(t *testing.T, caCert *x509.Certificate, caKey *rsa.PrivateKey, hostname string) *x509.Certificate {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			Organization: []string{"Test Expired"},
		},
		DNSNames:  []string{hostname},
		NotBefore: time.Now().Add(-365 * 24 * time.Hour), // Expired a year ago
		NotAfter:  time.Now().Add(-30 * 24 * time.Hour),  // Expired a month ago
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, caCert, &privateKey.PublicKey, caKey)
	require.NoError(t, err)

	cert, err := x509.ParseCertificate(certDER)
	require.NoError(t, err)

	return cert
}

// readAuditRecords reads and parses audit records from a log file.
func readAuditRecords(t *testing.T, logPath string) []AuditRecord {
	file, err := os.Open(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []AuditRecord{}
		}
		t.Fatalf("Failed to open audit log file: %v", err)
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
			t.Fatalf("Failed to unmarshal audit record: %v", err)
		}
		records = append(records, record)
	}

	return records
}