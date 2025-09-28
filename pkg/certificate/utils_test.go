package certificate

import (
	"crypto/x509"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCertificateFromPEM(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("parse valid certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.Subject = "test.example.com"

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		x509Cert, err := ParseCertificateFromPEM(cert.CertPEM)
		require.NoError(t, err)
		assert.NotNil(t, x509Cert)
		assert.Equal(t, "test.example.com", x509Cert.Subject.CommonName)
	})

	t.Run("parse invalid PEM data", func(t *testing.T) {
		invalidPEM := []byte("not a valid PEM")

		_, err := ParseCertificateFromPEM(invalidPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode PEM block")
	})

	t.Run("parse wrong PEM block type", func(t *testing.T) {
		// Create a valid PEM block but with wrong type
		wrongTypePEM := []byte(`-----BEGIN PRIVATE KEY-----
dGVzdCBkYXRh
-----END PRIVATE KEY-----`)

		_, err := ParseCertificateFromPEM(wrongTypePEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid PEM block type")
	})
}

func TestValidateCertificate(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("validate current certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.ValidFor = time.Hour * 24

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		err = ValidateCertificate(cert.CertPEM)
		assert.NoError(t, err)
	})

	t.Run("validate expired certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.ValidFor = -time.Hour // Already expired

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		err = ValidateCertificate(cert.CertPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "has expired")
	})

	t.Run("validate invalid PEM", func(t *testing.T) {
		invalidPEM := []byte("invalid PEM data")

		err := ValidateCertificate(invalidPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parsing certificate")
	})
}

func TestGetCertificateInfo(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("get info from valid certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.Subject = "test.example.com"
		opts.DNSNames = []string{"test.example.com", "localhost"}
		opts.Organization = "Test Organization"

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		info, err := GetCertificateInfo(cert.CertPEM)
		require.NoError(t, err)
		assert.Equal(t, "test.example.com", info.Subject)
		assert.Equal(t, []string{"test.example.com", "localhost"}, info.DNSNames)
		assert.False(t, info.IsCA)
		assert.Contains(t, info.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	})

	t.Run("get info from CA certificate", func(t *testing.T) {
		opts := CreateGenerationOptionsForCA("Test CA", time.Hour*24*365)

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		info, err := GetCertificateInfo(cert.CertPEM)
		require.NoError(t, err)
		assert.True(t, info.IsCA)
		assert.Contains(t, info.Subject, "Test CA Certificate Authority")
	})

	t.Run("get info from invalid PEM", func(t *testing.T) {
		invalidPEM := []byte("invalid PEM")

		_, err := GetCertificateInfo(invalidPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parsing certificate")
	})
}

func TestIsCertificateExpiringSoon(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("certificate expiring soon", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.ValidFor = time.Hour * 2 // Expires in 2 hours

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		// Check if expiring within 3 hours
		expiringSoon, err := IsCertificateExpiringSoon(cert.CertPEM, time.Hour*3)
		require.NoError(t, err)
		assert.True(t, expiringSoon)
	})

	t.Run("certificate not expiring soon", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.ValidFor = time.Hour * 48 // Expires in 48 hours

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		// Check if expiring within 3 hours
		expiringSoon, err := IsCertificateExpiringSoon(cert.CertPEM, time.Hour*3)
		require.NoError(t, err)
		assert.False(t, expiringSoon)
	})

	t.Run("invalid certificate", func(t *testing.T) {
		invalidPEM := []byte("invalid PEM")

		_, err := IsCertificateExpiringSoon(invalidPEM, time.Hour)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parsing certificate")
	})
}

func TestGetCertificateFingerprint(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("get fingerprint from valid certificate", func(t *testing.T) {
		opts := DefaultGenerationOptions()
		opts.Subject = "test.example.com"

		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		fingerprint, err := GetCertificateFingerprint(cert.CertPEM)
		require.NoError(t, err)
		assert.NotEmpty(t, fingerprint)
		assert.Regexp(t, "^[0-9a-f]+$", fingerprint) // Should be hex string
	})

	t.Run("fingerprints should be different for different certificates", func(t *testing.T) {
		opts1 := DefaultGenerationOptions()
		opts1.Subject = "test1.example.com"

		cert1, err := manager.GenerateSelfSigned(opts1)
		require.NoError(t, err)

		opts2 := DefaultGenerationOptions()
		opts2.Subject = "test2.example.com"

		cert2, err := manager.GenerateSelfSigned(opts2)
		require.NoError(t, err)

		fingerprint1, err := GetCertificateFingerprint(cert1.CertPEM)
		require.NoError(t, err)

		fingerprint2, err := GetCertificateFingerprint(cert2.CertPEM)
		require.NoError(t, err)

		assert.NotEqual(t, fingerprint1, fingerprint2)
	})

	t.Run("get fingerprint from invalid certificate", func(t *testing.T) {
		invalidPEM := []byte("invalid PEM")

		_, err := GetCertificateFingerprint(invalidPEM)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parsing certificate")
	})
}

func TestCreateGenerationOptionsForTLS(t *testing.T) {
	t.Run("create TLS options with custom parameters", func(t *testing.T) {
		dnsNames := []string{"example.com", "www.example.com", "api.example.com"}
		validFor := time.Hour * 24 * 90

		opts := CreateGenerationOptionsForTLS(dnsNames, validFor)

		assert.Equal(t, dnsNames, opts.DNSNames)
		assert.Equal(t, validFor, opts.ValidFor)
		assert.Equal(t, "TLS Server Certificate", opts.Subject)
		assert.False(t, opts.IsCA)
		assert.Equal(t, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment, opts.KeyUsage)
		assert.Contains(t, opts.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	})

	t.Run("create TLS options with empty DNS names", func(t *testing.T) {
		opts := CreateGenerationOptionsForTLS([]string{}, time.Hour*24)

		assert.Empty(t, opts.DNSNames)
		assert.Equal(t, "TLS Server Certificate", opts.Subject)
	})
}

func TestCreateGenerationOptionsForCA(t *testing.T) {
	t.Run("create CA options with custom parameters", func(t *testing.T) {
		organization := "My Test Organization"
		validFor := time.Hour * 24 * 365 * 10 // 10 years

		opts := CreateGenerationOptionsForCA(organization, validFor)

		assert.Equal(t, organization, opts.Organization)
		assert.Equal(t, validFor, opts.ValidFor)
		assert.Equal(t, "My Test Organization Certificate Authority", opts.Subject)
		assert.True(t, opts.IsCA)
		assert.Equal(t, x509.KeyUsageDigitalSignature|x509.KeyUsageCertSign|x509.KeyUsageCRLSign, opts.KeyUsage)
		assert.Empty(t, opts.ExtKeyUsage) // CAs typically don't have extended key usage
	})
}

func TestDefaultGenerationOptions(t *testing.T) {
	t.Run("verify default options", func(t *testing.T) {
		opts := DefaultGenerationOptions()

		assert.Equal(t, "IDPBuilder", opts.Organization)
		assert.Equal(t, time.Hour*24*365, opts.ValidFor) // 1 year
		assert.False(t, opts.IsCA)
		assert.Equal(t, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment, opts.KeyUsage)
		assert.Contains(t, opts.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	})
}

func TestUtilsIntegration(t *testing.T) {
	manager := NewMemoryManager()

	t.Run("complete certificate lifecycle", func(t *testing.T) {
		// Generate certificate
		opts := CreateGenerationOptionsForTLS([]string{"test.example.com"}, time.Hour*24)
		cert, err := manager.GenerateSelfSigned(opts)
		require.NoError(t, err)

		// Validate certificate
		err = ValidateCertificate(cert.CertPEM)
		assert.NoError(t, err)

		// Get certificate info
		info, err := GetCertificateInfo(cert.CertPEM)
		require.NoError(t, err)
		assert.Equal(t, "TLS Server Certificate", info.Subject)

		// Check if expiring soon (shouldn't be)
		expiringSoon, err := IsCertificateExpiringSoon(cert.CertPEM, time.Hour)
		require.NoError(t, err)
		assert.False(t, expiringSoon)

		// Get fingerprint
		fingerprint, err := GetCertificateFingerprint(cert.CertPEM)
		require.NoError(t, err)
		assert.NotEmpty(t, fingerprint)

		// Parse certificate directly
		x509Cert, err := ParseCertificateFromPEM(cert.CertPEM)
		require.NoError(t, err)
		assert.Equal(t, "TLS Server Certificate", x509Cert.Subject.CommonName)
	})
}