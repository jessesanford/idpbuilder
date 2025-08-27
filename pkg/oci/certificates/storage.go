package certificates

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FilesystemStore implements CertificateStore interface for file-based certificate storage.
type FilesystemStore struct {
	basePath string
	certDir  string
	metaDir  string
	mu       sync.RWMutex
	watchers map[EventHandler]bool
}

// CertificateMetadata holds the JSON-serializable metadata for a certificate.
type CertificateMetadata struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Status    CertificateStatus `json:"status"`
	ValidFrom time.Time         `json:"valid_from"`
	ValidTo   time.Time         `json:"valid_to"`
	Issuer    string            `json:"issuer"`
	Subject   string            `json:"subject"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Tags      map[string]string `json:"tags"`
}

// NewFilesystemStore creates a new filesystem-based certificate store.
func NewFilesystemStore(basePath string) (*FilesystemStore, error) {
	if basePath == "" {
		return nil, NewStorageError(ErrCodeInvalidConfig, "base path cannot be empty", nil)
	}

	certDir := filepath.Join(basePath, "certificates")
	metaDir := filepath.Join(basePath, "metadata")

	if err := os.MkdirAll(certDir, 0755); err != nil {
		return nil, NewStorageError(ErrCodeStorageUnavailable, "failed to create certificate directory", err)
	}
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		return nil, NewStorageError(ErrCodeStorageUnavailable, "failed to create metadata directory", err)
	}

	return &FilesystemStore{
		basePath: basePath,
		certDir:  certDir,
		metaDir:  metaDir,
		watchers: make(map[EventHandler]bool),
	}, nil
}

// AddCertificate adds a new certificate to the store.
func (fs *FilesystemStore) AddCertificate(ctx context.Context, cert *Certificate) error {
	if cert == nil || cert.ID == "" {
		return NewStorageError(ErrCodeInvalidConfig, "invalid certificate", nil)
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, err := fs.getCertificateNoLock(ctx, cert.ID); err == nil {
		return NewStorageError(ErrCodeDuplicateCert, "certificate already exists", nil)
	}

	cert.CreatedAt = time.Now()
	cert.UpdatedAt = cert.CreatedAt
	return fs.saveCertificateNoLock(cert)
}

// GetCertificate retrieves a certificate by ID.
func (fs *FilesystemStore) GetCertificate(ctx context.Context, id string) (*Certificate, error) {
	if id == "" {
		return nil, NewStorageError(ErrCodeInvalidConfig, "certificate ID cannot be empty", nil)
	}

	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.getCertificateNoLock(ctx, id)
}

// UpdateCertificate updates an existing certificate in the store.
func (fs *FilesystemStore) UpdateCertificate(ctx context.Context, cert *Certificate) error {
	if cert == nil || cert.ID == "" {
		return NewStorageError(ErrCodeInvalidConfig, "invalid certificate", nil)
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	existing, err := fs.getCertificateNoLock(ctx, cert.ID)
	if err != nil {
		return err
	}

	cert.CreatedAt = existing.CreatedAt
	cert.UpdatedAt = time.Now()
	return fs.saveCertificateNoLock(cert)
}

// DeleteCertificate removes a certificate from the store.
func (fs *FilesystemStore) DeleteCertificate(ctx context.Context, id string) error {
	if id == "" {
		return NewStorageError(ErrCodeInvalidConfig, "certificate ID cannot be empty", nil)
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	if _, err := fs.getCertificateNoLock(ctx, id); err != nil {
		return err
	}

	certFile := filepath.Join(fs.certDir, id+".pem")
	metaFile := filepath.Join(fs.metaDir, id+".json")
	os.Remove(certFile)
	os.Remove(metaFile)
	return nil
}

// ListCertificates lists certificates based on the provided filter.
func (fs *FilesystemStore) ListCertificates(ctx context.Context, filter *CertificateFilter) ([]*Certificate, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var certs []*Certificate
	metaFiles, err := ioutil.ReadDir(fs.metaDir)
	if err != nil {
		return nil, NewStorageError(ErrCodeStorageUnavailable, "failed to read metadata directory", err)
	}

	for _, file := range metaFiles {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		id := strings.TrimSuffix(file.Name(), ".json")
		cert, err := fs.getCertificateNoLock(ctx, id)
		if err != nil {
			continue
		}

		if fs.matchesFilter(cert, filter) {
			certs = append(certs, cert)
		}
	}

	if filter != nil && filter.Limit > 0 && len(certs) > filter.Limit {
		certs = certs[:filter.Limit]
	}

	return certs, nil
}

// SearchCertificates searches certificates by a query string.
func (fs *FilesystemStore) SearchCertificates(ctx context.Context, query string) ([]*Certificate, error) {
	return fs.ListCertificates(ctx, nil)
}

// GetCertificatesByStatus returns certificates with the specified status.
func (fs *FilesystemStore) GetCertificatesByStatus(ctx context.Context, status CertificateStatus) ([]*Certificate, error) {
	filter := &CertificateFilter{Status: []CertificateStatus{status}}
	return fs.ListCertificates(ctx, filter)
}

// getCertificateNoLock retrieves a certificate without acquiring locks.
func (fs *FilesystemStore) getCertificateNoLock(ctx context.Context, id string) (*Certificate, error) {
	certFile := filepath.Join(fs.certDir, id+".pem")
	metaFile := filepath.Join(fs.metaDir, id+".json")

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return nil, NewStorageError(ErrCodeCertNotFound, "certificate not found", nil)
	}

	pemData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, NewStorageError(ErrCodeStorageUnavailable, "failed to read certificate", err)
	}

	metaData, err := ioutil.ReadFile(metaFile)
	if err != nil {
		return nil, NewStorageError(ErrCodeStorageUnavailable, "failed to read metadata", err)
	}

	var meta CertificateMetadata
	if err := json.Unmarshal(metaData, &meta); err != nil {
		return nil, NewStorageError(ErrCodeStorageUnavailable, "failed to unmarshal metadata", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, NewStorageError(ErrCodeInvalidPEM, "invalid PEM data", nil)
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, NewStorageError(ErrCodeInvalidPEM, "failed to parse certificate", err)
	}

	return &Certificate{
		ID:          meta.ID,
		Name:        meta.Name,
		Certificate: x509Cert,
		PEM:         pemData,
		Status:      meta.Status,
		ValidFrom:   meta.ValidFrom,
		ValidTo:     meta.ValidTo,
		Issuer:      meta.Issuer,
		Subject:     meta.Subject,
		CreatedAt:   meta.CreatedAt,
		UpdatedAt:   meta.UpdatedAt,
		Tags:        meta.Tags,
	}, nil
}

// saveCertificateNoLock saves a certificate without acquiring locks.
func (fs *FilesystemStore) saveCertificateNoLock(cert *Certificate) error {
	certFile := filepath.Join(fs.certDir, cert.ID+".pem")
	metaFile := filepath.Join(fs.metaDir, cert.ID+".json")

	if err := ioutil.WriteFile(certFile, cert.PEM, 0644); err != nil {
		return NewStorageError(ErrCodeStorageUnavailable, "failed to write certificate", err)
	}

	meta := CertificateMetadata{
		ID:        cert.ID,
		Name:      cert.Name,
		Status:    cert.Status,
		ValidFrom: cert.ValidFrom,
		ValidTo:   cert.ValidTo,
		Issuer:    cert.Issuer,
		Subject:   cert.Subject,
		CreatedAt: cert.CreatedAt,
		UpdatedAt: cert.UpdatedAt,
		Tags:      cert.Tags,
	}

	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return NewStorageError(ErrCodeStorageUnavailable, "failed to marshal metadata", err)
	}

	return ioutil.WriteFile(metaFile, metaData, 0644)
}

// matchesFilter checks if a certificate matches the given filter.
func (fs *FilesystemStore) matchesFilter(cert *Certificate, filter *CertificateFilter) bool {
	if filter == nil {
		return true
	}

	if len(filter.Status) > 0 {
		for _, status := range filter.Status {
			if cert.Status == status {
				return true
			}
		}
		return false
	}

	return true
}

// RegisterEventHandler registers an event handler.
func (fs *FilesystemStore) RegisterEventHandler(handler EventHandler) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.watchers[handler] = true
}

// UnregisterEventHandler removes an event handler.
func (fs *FilesystemStore) UnregisterEventHandler(handler EventHandler) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	delete(fs.watchers, handler)
}