package oci

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MockAuthenticator provides a configurable mock implementation for testing
type MockAuthenticator struct {
	AuthFunc      func(ctx context.Context, reg string) (Authenticator, error)
	ValidateFunc  func() error
	Credentials   map[string]*MockCredential
	CallCount     map[string]int
	ErrorOnCall   map[string]error
	mu            sync.RWMutex
}

// MockCredential represents test credential data
type MockCredential struct {
	Username   string
	Password   string
	Token      string
	ValidUntil time.Time
	Scopes     []string
}

// Authenticator interface that our mock implements
type Authenticator interface {
	Authenticate(ctx context.Context, registry string) error
	GetCredential(key string) (*MockCredential, error)
	Validate() error
}

// NewMockAuthenticator creates a new mock authenticator
func NewMockAuthenticator() *MockAuthenticator {
	return &MockAuthenticator{
		Credentials: make(map[string]*MockCredential),
		CallCount:   make(map[string]int),
		ErrorOnCall: make(map[string]error),
	}
}

// Authenticate implements the Authenticator interface
func (m *MockAuthenticator) Authenticate(ctx context.Context, registry string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallCount["Authenticate"]++

	if err, exists := m.ErrorOnCall["Authenticate"]; exists {
		return err
	}

	if m.AuthFunc != nil {
		_, err := m.AuthFunc(ctx, registry)
		return err
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

	cred, exists := m.Credentials[key]
	if !exists {
		return nil, fmt.Errorf("credential not found: %s", key)
	}

	return cred, nil
}

// Validate implements the Authenticator interface
func (m *MockAuthenticator) Validate() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CallCount["Validate"]++

	if err, exists := m.ErrorOnCall["Validate"]; exists {
		return err
	}

	if m.ValidateFunc != nil {
		return m.ValidateFunc()
	}

	return nil
}

// SetError configures the mock to return an error for a specific method
func (m *MockAuthenticator) SetError(method string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorOnCall[method] = err
}

// Reset clears all call history and errors for test isolation
func (m *MockAuthenticator) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CallCount = make(map[string]int)
	m.ErrorOnCall = make(map[string]error)
}

// GetCallCount returns the number of times a method was called
func (m *MockAuthenticator) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.CallCount[method]
}

// AddCredential adds a test credential to the mock
func (m *MockAuthenticator) AddCredential(key string, cred *MockCredential) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Credentials[key] = cred
}