package certificates

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"testing"
	"time"
)

// TestFilesystemStore_Basic tests basic CRUD operations.
func TestFilesystemStore_Basic(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "cert-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	store, err := NewFilesystemStore(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	cert := createTestCert(t, "test-1")
	ctx := context.Background()

	if err := store.AddCertificate(ctx, cert); err != nil {
		t.Fatal(err)
	}

	retrieved, err := store.GetCertificate(ctx, cert.ID)
	if err != nil {
		t.Fatal(err)
	}

	if retrieved.ID != cert.ID {
		t.Errorf("Expected ID %s, got %s", cert.ID, retrieved.ID)
	}
}

func createTestCert(t *testing.T, id string) *Certificate {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test.com"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment,
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
	}
	certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	x509Cert, _ := x509.ParseCertificate(certDER)
	pemData := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	return &Certificate{
		ID:          id,
		Name:        "Test Certificate",
		Certificate: x509Cert,
		PEM:         pemData,
		Status:      CertificateStatusActive,
		ValidFrom:   x509Cert.NotBefore,
		ValidTo:     x509Cert.NotAfter,
		Issuer:      x509Cert.Issuer.String(),
		Subject:     x509Cert.Subject.String(),
		Tags:        make(map[string]string),
	}
}