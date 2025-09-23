package oci

import (
	"fmt"
	"sync"
	"time"
)

// AccessEntry represents a secret access log entry
type AccessEntry struct {
	Key       string
	Operation string
	Timestamp time.Time
	Success   bool
}

// MockSecretManager provides a test double for secret storage
type MockSecretManager struct {
	Store      map[string][]byte
	AccessLog  []AccessEntry
	FailOnKeys []string
	mu         sync.RWMutex
}

// NewMockSecretManager creates a new mock secret manager
func NewMockSecretManager() *MockSecretManager {
	return &MockSecretManager{
		Store:      make(map[string][]byte),
		AccessLog:  make([]AccessEntry, 0),
		FailOnKeys: make([]string, 0),
	}
}

// GetSecret retrieves a secret by key
func (m *MockSecretManager) GetSecret(key string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := AccessEntry{
		Key:       key,
		Operation: "GET",
		Timestamp: time.Now(),
		Success:   false,
	}

	// Check if this key should fail
	for _, failKey := range m.FailOnKeys {
		if failKey == key {
			m.AccessLog = append(m.AccessLog, entry)
			return nil, fmt.Errorf("simulated failure for key: %s", key)
		}
	}

	value, exists := m.Store[key]
	if !exists {
		m.AccessLog = append(m.AccessLog, entry)
		return nil, fmt.Errorf("secret not found: %s", key)
	}

	entry.Success = true
	m.AccessLog = append(m.AccessLog, entry)
	return value, nil
}

// StoreSecret stores a secret with the given key
func (m *MockSecretManager) StoreSecret(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := AccessEntry{
		Key:       key,
		Operation: "STORE",
		Timestamp: time.Now(),
		Success:   true,
	}

	m.Store[key] = value
	m.AccessLog = append(m.AccessLog, entry)
	return nil
}

// DeleteSecret removes a secret by key
func (m *MockSecretManager) DeleteSecret(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry := AccessEntry{
		Key:       key,
		Operation: "DELETE",
		Timestamp: time.Now(),
		Success:   false,
	}

	if _, exists := m.Store[key]; !exists {
		m.AccessLog = append(m.AccessLog, entry)
		return fmt.Errorf("secret not found: %s", key)
	}

	delete(m.Store, key)
	entry.Success = true
	m.AccessLog = append(m.AccessLog, entry)
	return nil
}

// ListSecrets returns all secret keys
func (m *MockSecretManager) ListSecrets() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.Store))
	for key := range m.Store {
		keys = append(keys, key)
	}
	return keys
}

// SimulateFailure configures the mock to fail on specific keys
func (m *MockSecretManager) SimulateFailure(keys ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FailOnKeys = append(m.FailOnKeys, keys...)
}

// GetAccessLog returns the access log for verification
func (m *MockSecretManager) GetAccessLog() []AccessEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	log := make([]AccessEntry, len(m.AccessLog))
	copy(log, m.AccessLog)
	return log
}

// Reset clears all stored secrets and access logs
func (m *MockSecretManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Store = make(map[string][]byte)
	m.AccessLog = make([]AccessEntry, 0)
	m.FailOnKeys = make([]string, 0)
}