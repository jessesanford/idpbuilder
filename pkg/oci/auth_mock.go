package oci

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MockAuthenticator provides configurable mock for testing
type MockAuthenticator struct {
	AuthFunc    func(ctx context.Context, reg string) error
	Credentials map[string]*MockCredential
	CallCount   map[string]int
	ErrorOnCall map[string]error
	mu          sync.RWMutex
}

// MockCredential represents test credential data
type MockCredential struct {
	Username   string
	Token      string
	ValidUntil time.Time
}

// NewMockAuthenticator creates a new mock authenticator
func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{
		Credentials: make(map[string]*MockCredential),
		CallCount:   make(map[string]int),
		ErrorOnCall: make(map[string]error),
	}
}

// Authenticate implements authentication with configurable behavior
func (m *MockAuthenticator) Authenticate(ctx context.Context, registry string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CallCount["Authenticate"]++

	if err, exists := m.ErrorOnCall["Authenticate"]; exists {
		return err
	}
	if m.AuthFunc != nil {
		return m.AuthFunc(ctx, registry)
	}
	return nil
}

// GetCredential retrieves a mock credential by key
func (m *MockAuthenticator) GetCredential(key string) (*MockCredential, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.CallCount["GetCredential"]++

	if err, exists := m.ErrorOnCall["GetCredential"]; exists {
		return nil, err
	}
	if cred, exists := m.Credentials[key]; exists {
		return cred, nil
	}
	return nil, fmt.Errorf("credential not found: %s", key)
}

// SetError configures error injection for testing
func (m *MockAuthenticator) SetError(method string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorOnCall[method] = err
}

// Reset clears state for test isolation
func (m *MockAuthenticator) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CallCount = make(map[string]int)
	m.ErrorOnCall = make(map[string]error)
}

// GetCallCount returns method call count
func (m *MockAuthenticator) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.CallCount[method]
}

// AddCredential adds test credential
func (m *MockAuthenticator) AddCredential(key string, cred *MockCredential) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Credentials[key] = cred
}

// MockSecretManager provides basic secret storage for testing
type MockSecretManager struct {
	Store map[string][]byte
	mu    sync.RWMutex
}

// NewMockSecretManager creates new mock secret manager
func NewMockSecretManager() *MockSecretManager {
	return &MockSecretManager{
		Store: make(map[string][]byte),
	}
}

// GetSecret retrieves secret by key
func (m *MockSecretManager) GetSecret(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if value, exists := m.Store[key]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("secret not found: %s", key)
}

// StoreSecret stores secret with key
func (m *MockSecretManager) StoreSecret(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Store[key] = value
	return nil
}