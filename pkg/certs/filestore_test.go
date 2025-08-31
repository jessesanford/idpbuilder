package certs

import (
	"os"
	"testing"
)

func TestFileStore_Store(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}

	store := newFileStore(config)

	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}

	// Parse certificate info first
	mgr := &manager{}
	if err := mgr.parseCertificateInfo(cert); err != nil {
		t.Fatalf("Failed to parse certificate info: %v", err)
	}

	err := store.Store("test-registry", cert)
	if err != nil {
		t.Fatalf("Failed to store certificate: %v", err)
	}

	// Verify file was created
	if cert.FilePath == "" {
		t.Error("Expected file path to be set")
	}

	if _, err := os.Stat(cert.FilePath); err != nil {
		t.Errorf("Certificate file does not exist: %v", err)
	}
}

func TestFileStore_Load(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}

	store := newFileStore(config)

	// Store a certificate first
	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}

	mgr := &manager{}
	if err := mgr.parseCertificateInfo(cert); err != nil {
		t.Fatalf("Failed to parse certificate info: %v", err)
	}

	err := store.Store("test-registry", cert)
	if err != nil {
		t.Fatalf("Failed to store certificate: %v", err)
	}

	// Load the certificate
	loadedCert, err := store.Load("test-registry", cert.Info.Fingerprint)
	if err != nil {
		t.Fatalf("Failed to load certificate: %v", err)
	}

	if loadedCert.Info.Fingerprint != cert.Info.Fingerprint {
		t.Error("Loaded certificate fingerprint does not match")
	}
}

func TestFileStore_Delete(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}

	store := newFileStore(config)

	// Store a certificate first
	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}

	mgr := &manager{}
	if err := mgr.parseCertificateInfo(cert); err != nil {
		t.Fatalf("Failed to parse certificate info: %v", err)
	}

	err := store.Store("test-registry", cert)
	if err != nil {
		t.Fatalf("Failed to store certificate: %v", err)
	}

	// Delete the certificate
	err = store.Delete("test-registry", cert.Info.Fingerprint)
	if err != nil {
		t.Fatalf("Failed to delete certificate: %v", err)
	}

	// Verify it's gone
	_, err = store.Load("test-registry", cert.Info.Fingerprint)
	if err == nil {
		t.Error("Expected error when loading deleted certificate")
	}
}

func TestFileStore_List(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}

	store := newFileStore(config)

	// Store multiple certificates
	cert1 := &Certificate{Data: createTestCertificate(t, "test1.example.com")}
	cert2 := &Certificate{Data: createTestCertificate(t, "test2.example.com")}

	mgr := &manager{}

	if err := mgr.parseCertificateInfo(cert1); err != nil {
		t.Fatalf("Failed to parse certificate 1 info: %v", err)
	}
	if err := mgr.parseCertificateInfo(cert2); err != nil {
		t.Fatalf("Failed to parse certificate 2 info: %v", err)
	}

	if err := store.Store("test-registry", cert1); err != nil {
		t.Fatalf("Failed to store certificate 1: %v", err)
	}
	if err := store.Store("test-registry", cert2); err != nil {
		t.Fatalf("Failed to store certificate 2: %v", err)
	}

	// List certificates
	certs, err := store.List("test-registry")
	if err != nil {
		t.Fatalf("Failed to list certificates: %v", err)
	}

	if len(certs) != 2 {
		t.Errorf("Expected 2 certificates, got %d", len(certs))
	}
}

func TestFileStore_Exists(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}

	store := newFileStore(config)

	// Test non-existent certificate
	exists, err := store.Exists("test-registry", "nonexistent")
	if err != nil {
		t.Fatalf("Unexpected error checking existence: %v", err)
	}
	if exists {
		t.Error("Expected certificate to not exist")
	}

	// Store a certificate
	testCert := createTestCertificate(t, "test.example.com")
	cert := &Certificate{Data: testCert}

	mgr := &manager{}
	if err := mgr.parseCertificateInfo(cert); err != nil {
		t.Fatalf("Failed to parse certificate info: %v", err)
	}

	if err := store.Store("test-registry", cert); err != nil {
		t.Fatalf("Failed to store certificate: %v", err)
	}

	// Test existing certificate
	exists, err = store.Exists("test-registry", cert.Info.Fingerprint)
	if err != nil {
		t.Fatalf("Unexpected error checking existence: %v", err)
	}
	if !exists {
		t.Error("Expected certificate to exist")
	}
}

func TestFileStore_ListEmptyRegistry(t *testing.T) {
	tempDir := t.TempDir()
	config := TrustStoreConfig{
		Location:        UserTrustStore,
		BaseDir:         tempDir,
		DirPermissions:  0755,
		FilePermissions: 0600,
	}

	store := newFileStore(config)

	// List certificates for a registry with no certificates
	certs, err := store.List("empty-registry")
	if err != nil {
		t.Fatalf("Failed to list certificates for empty registry: %v", err)
	}

	if len(certs) != 0 {
		t.Errorf("Expected 0 certificates for empty registry, got %d", len(certs))
	}
}
