package certificates

import (
	"context"
	"crypto/x509"
	"fmt"
	"sync"
	"time"
)

// PoolManager implements the CertPoolManager interface for managing certificate pools.
type PoolManager struct {
	store     CertificateStore
	pools     map[string]*CertificatePool
	systemCAs *x509.CertPool
	mu        sync.RWMutex
	validator CertificateValidator
}

// CertificatePool represents a named collection of certificates.
type CertificatePool struct {
	Name          string
	CertificateIDs []string
	x509Pool      *x509.CertPool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewPoolManager creates a new certificate pool manager.
func NewPoolManager(store CertificateStore) (*PoolManager, error) {
	if store == nil {
		return nil, fmt.Errorf("certificate store cannot be nil")
	}

	systemCAs, _ := x509.SystemCertPool()
	if systemCAs == nil {
		systemCAs = x509.NewCertPool()
	}

	return &PoolManager{
		store:     store,
		pools:     make(map[string]*CertificatePool),
		systemCAs: systemCAs,
	}, nil
}

// CreatePool creates a new named certificate pool.
func (pm *PoolManager) CreatePool(ctx context.Context, name string) error {
	if name == "" {
		return NewStorageError(ErrCodeInvalidConfig, "pool name cannot be empty", nil)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.pools[name]; exists {
		return NewStorageError(ErrCodePoolExists, fmt.Sprintf("pool '%s' already exists", name), nil)
	}

	pool := &CertificatePool{
		Name:          name,
		CertificateIDs: make([]string, 0),
		x509Pool:      x509.NewCertPool(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	pm.pools[name] = pool
	return nil
}

// DeletePool removes a certificate pool.
func (pm *PoolManager) DeletePool(ctx context.Context, name string) error {
	if name == "" {
		return NewStorageError(ErrCodeInvalidConfig, "pool name cannot be empty", nil)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if _, exists := pm.pools[name]; !exists {
		return NewStorageError(ErrCodePoolNotFound, fmt.Sprintf("pool '%s' not found", name), nil)
	}

	delete(pm.pools, name)
	return nil
}

// AddCertificateToPool adds a certificate to a named pool.
func (pm *PoolManager) AddCertificateToPool(ctx context.Context, poolName, certificateID string) error {
	if poolName == "" || certificateID == "" {
		return NewStorageError(ErrCodeInvalidConfig, "pool name and certificate ID cannot be empty", nil)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	pool, exists := pm.pools[poolName]
	if !exists {
		return NewStorageError(ErrCodePoolNotFound, fmt.Sprintf("pool '%s' not found", poolName), nil)
	}

	// Retrieve and validate certificate
	cert, err := pm.store.GetCertificate(ctx, certificateID)
	if err != nil {
		return fmt.Errorf("failed to get certificate: %w", err)
	}

	if pm.validator != nil {
		result, err := pm.validator.ValidateCertificate(ctx, cert)
		if err != nil || !result.Valid {
			return NewValidationError(ErrCodeInvalidCertChain, "certificate validation failed")
		}
	}

	// Check for duplicates
	for _, id := range pool.CertificateIDs {
		if id == certificateID {
			return NewStorageError(ErrCodeDuplicateCert, "certificate already in pool", nil)
		}
	}

	// Add to pools
	if cert.Certificate != nil {
		pool.x509Pool.AddCert(cert.Certificate)
	}
	pool.CertificateIDs = append(pool.CertificateIDs, certificateID)
	pool.UpdatedAt = time.Now()

	return nil
}

// RemoveCertificateFromPool removes a certificate from a named pool.
func (pm *PoolManager) RemoveCertificateFromPool(ctx context.Context, poolName, certificateID string) error {
	if poolName == "" || certificateID == "" {
		return NewStorageError(ErrCodeInvalidConfig, "pool name and certificate ID cannot be empty", nil)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	pool, exists := pm.pools[poolName]
	if !exists {
		return NewStorageError(ErrCodePoolNotFound, fmt.Sprintf("pool '%s' not found", poolName), nil)
	}

	// Find and remove certificate ID
	found := false
	for i, id := range pool.CertificateIDs {
		if id == certificateID {
			pool.CertificateIDs = append(pool.CertificateIDs[:i], pool.CertificateIDs[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return NewStorageError(ErrCodeCertNotFound, "certificate not found in pool", nil)
	}

	// Rebuild x509 pool
	pool.x509Pool = x509.NewCertPool()
	for _, id := range pool.CertificateIDs {
		if cert, err := pm.store.GetCertificate(ctx, id); err == nil && cert.Certificate != nil {
			pool.x509Pool.AddCert(cert.Certificate)
		}
	}

	pool.UpdatedAt = time.Now()
	return nil
}

// GetPool retrieves certificates in a named pool.
func (pm *PoolManager) GetPool(ctx context.Context, name string) ([]*Certificate, error) {
	if name == "" {
		return nil, NewStorageError(ErrCodeInvalidConfig, "pool name cannot be empty", nil)
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	pool, exists := pm.pools[name]
	if !exists {
		return nil, NewStorageError(ErrCodePoolNotFound, fmt.Sprintf("pool '%s' not found", name), nil)
	}

	certificates := make([]*Certificate, 0, len(pool.CertificateIDs))
	for _, id := range pool.CertificateIDs {
		if cert, err := pm.store.GetCertificate(ctx, id); err == nil {
			certificates = append(certificates, cert)
		}
	}

	return certificates, nil
}

// GetPoolNames returns the names of all pools.
func (pm *PoolManager) GetPoolNames(ctx context.Context) ([]string, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	names := make([]string, 0, len(pm.pools))
	for name := range pm.pools {
		names = append(names, name)
	}
	return names, nil
}

// GetX509Pool returns the x509.CertPool for a named pool.
func (pm *PoolManager) GetX509Pool(ctx context.Context, name string) (*x509.CertPool, error) {
	if name == "" {
		return nil, NewStorageError(ErrCodeInvalidConfig, "pool name cannot be empty", nil)
	}

	pm.mu.RLock()
	defer pm.mu.RUnlock()

	pool, exists := pm.pools[name]
	if !exists {
		return nil, NewStorageError(ErrCodePoolNotFound, fmt.Sprintf("pool '%s' not found", name), nil)
	}

	return pool.x509Pool, nil
}