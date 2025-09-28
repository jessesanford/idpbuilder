package certs

import (
	"context"
	"crypto/x509"
	"sync"
)

// Store provides an interface for certificate storage operations.
// This abstraction allows for different storage backends and enables
// easy testing through mock implementations.
type Store interface {
	// GetPool returns the certificate pool containing all stored certificates.
	// The context can be used for cancellation during long-running operations.
	GetPool(ctx context.Context) (*x509.CertPool, error)

	// AddCert adds a certificate to the store.
	// The context can be used for cancellation during long-running operations.
	AddCert(ctx context.Context, cert *x509.Certificate) error
}

// MemoryStore implements the Store interface using in-memory storage.
// It provides thread-safe certificate storage using a mutex for concurrent access.
// This implementation is suitable for runtime certificate management where
// persistence across application restarts is not required.
type MemoryStore struct {
	mu    sync.RWMutex
	pool  *x509.CertPool
	certs []*x509.Certificate
}

// NewMemoryStore creates a new MemoryStore instance with an empty certificate pool.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		pool:  x509.NewCertPool(),
		certs: make([]*x509.Certificate, 0),
	}
}

// GetPool returns the current certificate pool containing all stored certificates.
// This method is thread-safe and can be called concurrently.
func (s *MemoryStore) GetPool(ctx context.Context) (*x509.CertPool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Create a new pool and add all stored certificates
	// This ensures the returned pool is independent of internal state
	newPool := x509.NewCertPool()
	for _, cert := range s.certs {
		newPool.AddCert(cert)
	}

	return newPool, nil
}

// AddCert adds a certificate to the store.
// This method is thread-safe and can be called concurrently.
// If the certificate is nil, it returns an error without modifying the store.
func (s *MemoryStore) AddCert(ctx context.Context, cert *x509.Certificate) error {
	if cert == nil {
		return NewCertError(ErrInvalidCert, "cannot add nil certificate", nil)
	}

	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Add to both the pool and our slice for consistency
	s.pool.AddCert(cert)
	s.certs = append(s.certs, cert)

	return nil
}

// Size returns the number of certificates currently stored.
// This method is thread-safe and can be called concurrently.
func (s *MemoryStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.certs)
}

// Clear removes all certificates from the store.
// This method is thread-safe and can be called concurrently.
func (s *MemoryStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pool = x509.NewCertPool()
	s.certs = s.certs[:0] // Clear slice while preserving capacity
}