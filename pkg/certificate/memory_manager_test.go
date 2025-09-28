package certificate

import (
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryManager(t *testing.T) {
	manager := NewMemoryManager()

	require.NotNil(t, manager)
	assert.NotNil(t, manager.certificates)
	assert.Equal(t, 0, len(manager.certificates))
}

func TestMemoryManager_GenerateSelfSigned(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("generate with default options", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.Subject = "test.example.com"
		opts.DNSNames = []string{"test.example.com", "localhost"}

		cert, err := manager.GenerateSelfSigned(opts)

		require.NoError(t, err)
		require.NotNil(t, cert)
		assert.NotEmpty(t, cert.CertPEM)
		assert.NotEmpty(t, cert.KeyPEM)
		assert.NotNil(t, cert.Metadata)
		assert.Equal(t, "test.example.com", cert.Metadata.Subject)
		assert.Equal(t, []string{"test.example.com", "localhost"}, cert.Metadata.DNSNames)
		assert.False(t, cert.Metadata.IsCA)
	})

	t.Run("generate CA certificate", func(t *testing.T) {
		opts := CreateGenerationOptionsForCA("Test Organization", time.Hour*24*365)

		cert, err := manager.GenerateSelfSigned(opts)

		require.NoError(t, err)
		require.NotNil(t, cert)
		assert.True(t, cert.Metadata.IsCA)
		assert.Contains(t, cert.Metadata.Subject, "Test Organization Certificate Authority")
		assert.Equal(t, x509.KeyUsageDigitalSignature|x509.KeyUsageCertSign|x509.KeyUsageCRLSign, cert.Metadata.KeyUsage)
	})

	t.Run("generate TLS server certificate", func(t *testing.T) {
		dnsNames := []string{"example.com", "www.example.com"}
		opts := CreateGenerationOptionsForTLS(dnsNames, time.Hour*24*90)

		cert, err := manager.GenerateSelfSigned(opts)

		require.NoError(t, err)
		require.NotNil(t, cert)
		assert.False(t, cert.Metadata.IsCA)
		assert.Equal(t, dnsNames, cert.Metadata.DNSNames)
		assert.Contains(t, cert.Metadata.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	})
}

func TestMemoryManager_StoreAndRetrieve(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("store and retrieve certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.Subject = "test.example.com"

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		// Store the certificate
		err = manager.Store("test-cert", cert)
		require.NoError(t, err)

		// Retrieve the certificate
		retrievedCert, err := manager.Retrieve("test-cert")
		require.NoError(t, err)
		assert.Equal(t, cert.CertPEM, retrievedCert.CertPEM)
		assert.Equal(t, cert.KeyPEM, retrievedCert.KeyPEM)
		assert.Equal(t, cert.Metadata.Subject, retrievedCert.Metadata.Subject)
	})

	t.Run("store with empty key should fail", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		err = manager.Store("", cert)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate key cannot be empty")
	})

	t.Run("store nil certificate should fail", func(t *testing.T) {
		err := manager.Store("test-key", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate cannot be nil")
	})

	t.Run("retrieve non-existent certificate should fail", func(t *testing.T) {
		_, err := manager.Retrieve("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("retrieve with empty key should fail", func(t *testing.T) {
		_, err := manager.Retrieve("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate key cannot be empty")
	})
}

func TestMemoryManager_IsValid(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("valid certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.ValidFor = time.Hour * 24 // Valid for 1 day

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		isValid := manager.IsValid(cert)
		assert.True(t, isValid)
	})

	t.Run("expired certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.ValidFor = -time.Hour // Already expired

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		isValid := manager.IsValid(cert)
		assert.False(t, isValid)
	})

	t.Run("nil certificate", func(t *testing.T) {
		isValid := manager.IsValid(nil)
		assert.False(t, isValid)
	})

	t.Run("certificate with nil metadata", func(t *testing.T) {
		cert := &Certificate{
			CertPEM:  []byte("cert"),
			KeyPEM:   []byte("key"),
			Metadata: nil,
		}

		isValid := manager.IsValid(cert)
		assert.False(t, isValid)
	})
}

func TestMemoryManager_GetExpiration(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("get expiration from valid certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.ValidFor = time.Hour * 24

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		expiration, err := manager.GetExpiration(cert)
		require.NoError(t, err)

		// Should be approximately 24 hours from now
		expectedExpiration := time.Now().Add(time.Hour * 24)
		assert.WithinDuration(t, expectedExpiration, expiration, time.Minute)
	})

	t.Run("get expiration from nil certificate", func(t *testing.T) {
		_, err := manager.GetExpiration(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate or metadata is nil")
	})

	t.Run("get expiration from certificate with nil metadata", func(t *testing.T) {
		cert := &Certificate{
			CertPEM:  []byte("cert"),
			KeyPEM:   []byte("key"),
			Metadata: nil,
		}

		_, err := manager.GetExpiration(cert)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate or metadata is nil")
	})
}

func TestMemoryManager_ListAndDelete(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("list empty certificates", func(t *testing.T) {
		keys, err := manager.List()
		require.NoError(t, err)
		assert.Empty(t, keys)
	})

	t.Run("list and delete certificates", func(t *testing.T) {
		opts := DefaultGenerationOptions()

		// Store multiple certificates
		for i := 0; i < 3; i++ {
			cert, err := manager.GenerateSelfSigned(opts)
			require.NoError(t, err)

			key := fmt.Sprintf("cert-%d", i)
			err = manager.Store(key, cert)
			require.NoError(t, err)
		}

		// List certificates
		keys, err := manager.List()
		require.NoError(t, err)
		assert.Len(t, keys, 3)

		// Verify all keys are present
		expectedKeys := []string{"cert-0", "cert-1", "cert-2"}
		for _, expectedKey := range expectedKeys {
			assert.Contains(t, keys, expectedKey)
		}

		// Delete one certificate
		err = manager.Delete("cert-1")
		require.NoError(t, err)

		// Verify it was deleted
		keys, err = manager.List()
		require.NoError(t, err)
		assert.Len(t, keys, 2)
		assert.NotContains(t, keys, "cert-1")

		// Try to retrieve deleted certificate
		_, err = manager.Retrieve("cert-1")
		assert.Error(t, err)
	})

	t.Run("delete with empty key should fail", func(t *testing.T) {
		err := manager.Delete("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "certificate key cannot be empty")
	})

	t.Run("delete non-existent certificate should fail", func(t *testing.T) {
		err := manager.Delete("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestMemoryManager_ConcurrentAccess(t *testing.T) {
	manager := NewMemoryManager()

	// Test concurrent access doesn't cause race conditions
	t.Run("concurrent store and retrieve", func(t *testing.T) {
		opts := DefaultGenerationOptions()

		done := make(chan bool, 10)

		// Start multiple goroutines that store certificates
		for i := 0; i < 5; i++ {
			go func(index int) {
				defer func() { done <- true }()

				cert, err := manager.GenerateSelfSigned(opts)
				require.NoError(t, err)

				key := fmt.Sprintf("concurrent-cert-%d", index)
				err = manager.Store(key, cert)
				require.NoError(t, err)
			}(i)
		}

		// Start multiple goroutines that list certificates
		for i := 0; i < 5; i++ {
			go func() {
				defer func() { done <- true }()

				_, err := manager.List()
				require.NoError(t, err)
			}()
		}

		// Wait for all goroutines to complete
		for i := 0; i < 10; i++ {
			<-done
		}

		// Verify all certificates were stored
		keys, err := manager.List()
		require.NoError(t, err)
		assert.Len(t, keys, 5)
	})
}

