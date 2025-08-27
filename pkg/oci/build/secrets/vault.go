package secrets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
)

// memoryVault implements SecretVault with encrypted in-memory storage
type memoryVault struct {
	mu      sync.RWMutex
	secrets map[string]*encryptedSecret
	key     []byte           // AES-256 encryption key
	gcm     cipher.AEAD      // Galois/Counter Mode for authenticated encryption
}

// encryptedSecret represents an encrypted secret with its metadata
type encryptedSecret struct {
	encrypted []byte          // Encrypted secret value
	metadata  secretMetadata  // Non-sensitive metadata
}

// NewVault creates a new secure in-memory vault with AES-256 encryption
func NewVault() (SecretVault, error) {
	// Generate 32-byte key for AES-256
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate vault encryption key: %w", err)
	}
	
	// Create AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	
	// Create GCM mode for authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM mode: %w", err)
	}
	
	vault := &memoryVault{
		secrets: make(map[string]*encryptedSecret),
		key:     key,
		gcm:     gcm,
	}
	
	// Register finalizer for secure cleanup
	runtime.SetFinalizer(vault, (*memoryVault).Clear)
	
	log.Println("Secure vault initialized with AES-256 encryption")
	return vault, nil
}

// Store securely stores a secret in the vault with encryption
func (v *memoryVault) Store(secret *Secret) error {
	if secret == nil {
		return fmt.Errorf("secret cannot be nil")
	}
	
	if secret.ID == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}
	
	v.mu.Lock()
	defer v.mu.Unlock()
	
	// Check if secret already exists
	if _, exists := v.secrets[secret.ID]; exists {
		return fmt.Errorf("secret with ID %s already exists", secret.ID)
	}
	
	// Encrypt secret value
	encrypted, err := v.encrypt(secret.Value)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret %s: %w", secret.ID, err)
	}
	
	// Securely clear original value from memory
	clearBytes(secret.Value)
	
	// Store encrypted secret with metadata
	v.secrets[secret.ID] = &encryptedSecret{
		encrypted: encrypted,
		metadata: secretMetadata{
			Type:   secret.Type,
			Source: secret.Source,
			Target: secret.Target,
			Mode:   secret.Mode,
			UID:    secret.UID,
			GID:    secret.GID,
		},
	}
	
	// Audit log (without exposing the value)
	log.Printf("Secret stored: id=%s type=%s target=%s", secret.ID, secret.Type, secret.Target)
	
	return nil
}

// Retrieve decrypts and returns a secret from the vault
func (v *memoryVault) Retrieve(id string) (*Secret, error) {
	if id == "" {
		return nil, fmt.Errorf("secret ID cannot be empty")
	}
	
	v.mu.RLock()
	encSecret, exists := v.secrets[id]
	v.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("secret %s not found", id)
	}
	
	// Decrypt secret value
	value, err := v.decrypt(encSecret.encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret %s: %w", id, err)
	}
	
	// Reconstruct secret
	secret := &Secret{
		ID:     id,
		Type:   encSecret.metadata.Type,
		Value:  value,
		Source: encSecret.metadata.Source,
		Target: encSecret.metadata.Target,
		Mode:   encSecret.metadata.Mode,
		UID:    encSecret.metadata.UID,
		GID:    encSecret.metadata.GID,
	}
	
	// Audit log
	log.Printf("Secret retrieved: id=%s type=%s", id, secret.Type)
	
	return secret, nil
}

// Delete removes a secret from the vault
func (v *memoryVault) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}
	
	v.mu.Lock()
	defer v.mu.Unlock()
	
	encSecret, exists := v.secrets[id]
	if !exists {
		return fmt.Errorf("secret %s not found", id)
	}
	
	// Securely clear encrypted data
	clearBytes(encSecret.encrypted)
	
	// Remove from map
	delete(v.secrets, id)
	
	log.Printf("Secret deleted: id=%s", id)
	return nil
}

// Clear securely removes all secrets from the vault
func (v *memoryVault) Clear() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	// Securely clear all encrypted secrets
	for id, encSecret := range v.secrets {
		clearBytes(encSecret.encrypted)
		delete(v.secrets, id)
	}
	
	// Clear encryption key
	clearBytes(v.key)
	
	log.Println("Vault cleared securely - all secrets removed")
	return nil
}

// encrypt encrypts data using AES-256-GCM
func (v *memoryVault) encrypt(plaintext []byte) ([]byte, error) {
	// Generate random nonce
	nonce := make([]byte, v.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	
	// Encrypt and authenticate
	ciphertext := v.gcm.Seal(nonce, nonce, plaintext, nil)
	
	return ciphertext, nil
}

// decrypt decrypts data using AES-256-GCM
func (v *memoryVault) decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := v.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	plaintext, err := v.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}
	
	return plaintext, nil
}

// clearBytes securely zeros out a byte slice in memory
func clearBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}