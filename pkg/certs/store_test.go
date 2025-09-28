package certs

import (
	"context"
	"crypto/x509"
	"sync"
	"testing"
	"time"
)

func TestNewMemoryStore(t *testing.T) {
	store := NewMemoryStore()
	if store == nil {
		t.Error("NewMemoryStore() should not return nil")
	}
	if store.Size() != 0 {
		t.Errorf("NewMemoryStore() size = %d, want 0", store.Size())
	}
}

func TestMemoryStore_AddCert_NilCertificate(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	err := store.AddCert(ctx, nil)
	if err == nil {
		t.Error("AddCert() should fail with nil certificate")
	}

	certErr, ok := err.(*CertError)
	if !ok {
		t.Errorf("AddCert() error type = %T, want *CertError", err)
	} else if certErr.Code != ErrInvalidCert {
		t.Errorf("AddCert() error code = %v, want %v", certErr.Code, ErrInvalidCert)
	}
}

func TestMemoryStore_AddCert_ValidCertificate(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	cert := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}

	err := store.AddCert(ctx, cert)
	if err != nil {
		t.Errorf("AddCert() error = %v", err)
	}

	if store.Size() != 1 {
		t.Errorf("AddCert() size = %d, want 1", store.Size())
	}
}

func TestMemoryStore_GetPool(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Initially empty
	pool, err := store.GetPool(ctx)
	if err != nil {
		t.Errorf("GetPool() error = %v", err)
	}
	if pool == nil {
		t.Error("GetPool() should not return nil pool")
	}

	// Add a certificate
	cert := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}

	err = store.AddCert(ctx, cert)
	if err != nil {
		t.Errorf("AddCert() error = %v", err)
	}

	// Get pool again
	pool, err = store.GetPool(ctx)
	if err != nil {
		t.Errorf("GetPool() error = %v", err)
	}
	if pool == nil {
		t.Error("GetPool() should not return nil pool")
	}
}

func TestMemoryStore_Clear(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Add a certificate
	cert := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}

	err := store.AddCert(ctx, cert)
	if err != nil {
		t.Errorf("AddCert() error = %v", err)
	}

	if store.Size() != 1 {
		t.Errorf("Size() = %d, want 1", store.Size())
	}

	// Clear the store
	store.Clear()

	if store.Size() != 0 {
		t.Errorf("Size() after Clear() = %d, want 0", store.Size())
	}
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()
	ctx := context.Background()

	// Create certificates
	cert1 := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}
	cert2 := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	// Concurrent adds
	wg.Add(2)
	go func() {
		defer wg.Done()
		errChan <- store.AddCert(ctx, cert1)
	}()
	go func() {
		defer wg.Done()
		errChan <- store.AddCert(ctx, cert2)
	}()

	// Concurrent reads
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := store.GetPool(ctx)
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := store.GetPool(ctx)
		errChan <- err
	}()

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			t.Errorf("Concurrent operation error: %v", err)
		}
	}

	// Verify final state
	if store.Size() != 2 {
		t.Errorf("Final size = %d, want 2", store.Size())
	}
}

func TestMemoryStore_CancelledContext(t *testing.T) {
	store := NewMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	cert := &x509.Certificate{
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  time.Now().Add(time.Hour),
	}

	// Test AddCert with cancelled context
	err := store.AddCert(ctx, cert)
	if err == nil {
		t.Error("AddCert() should fail with cancelled context")
	}
	if err != context.Canceled {
		t.Errorf("AddCert() error = %v, want %v", err, context.Canceled)
	}

	// Test GetPool with cancelled context
	_, err = store.GetPool(ctx)
	if err == nil {
		t.Error("GetPool() should fail with cancelled context")
	}
	if err != context.Canceled {
		t.Errorf("GetPool() error = %v, want %v", err, context.Canceled)
	}
}