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

type memoryVault struct {
	mu      sync.RWMutex
	secrets map[string]*encryptedSecret
	key     []byte
	gcm     cipher.AEAD
}

type encryptedSecret struct {
	encrypted []byte
	metadata  secretMetadata
}

func NewVault() (SecretVault, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("generate vault key: %w", err)
	}
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}
	
	vault := &memoryVault{
		secrets: make(map[string]*encryptedSecret),
		key:     key,
		gcm:     gcm,
	}
	
	runtime.SetFinalizer(vault, (*memoryVault).Clear)
	log.Println("Secure vault initialized with AES-256 encryption")
	return vault, nil
}

func (v *memoryVault) Store(secret *Secret) error {
	if secret == nil {
		return fmt.Errorf("secret cannot be nil")
	}
	if secret.ID == "" {
		return fmt.Errorf("secret ID cannot be empty")
	}
	
	v.mu.Lock()
	defer v.mu.Unlock()
	
	if _, exists := v.secrets[secret.ID]; exists {
		return fmt.Errorf("secret ID %s already exists", secret.ID)
	}
	
	encrypted, err := v.encrypt(secret.Value)
	if err != nil {
		return fmt.Errorf("encrypt secret %s: %w", secret.ID, err)
	}
	
	clearBytes(secret.Value)
	
	v.secrets[secret.ID] = &encryptedSecret{
		encrypted: encrypted,
		metadata: secretMetadata{
			Type: secret.Type, Source: secret.Source, Target: secret.Target,
			Mode: secret.Mode, UID: secret.UID, GID: secret.GID,
		},
	}
	
	log.Printf("Secret stored: id=%s type=%s", secret.ID, secret.Type)
	return nil
}

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
	
	value, err := v.decrypt(encSecret.encrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt secret %s: %w", id, err)
	}
	
	secret := &Secret{
		ID: id, Type: encSecret.metadata.Type, Value: value,
		Source: encSecret.metadata.Source, Target: encSecret.metadata.Target,
		Mode: encSecret.metadata.Mode, UID: encSecret.metadata.UID, GID: encSecret.metadata.GID,
	}
	
	log.Printf("Secret retrieved: id=%s", id)
	return secret, nil
}

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
	
	clearBytes(encSecret.encrypted)
	delete(v.secrets, id)
	log.Printf("Secret deleted: id=%s", id)
	return nil
}

func (v *memoryVault) Clear() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	
	for id, encSecret := range v.secrets {
		clearBytes(encSecret.encrypted)
		delete(v.secrets, id)
	}
	clearBytes(v.key)
	log.Println("Vault cleared securely - all secrets removed")
	return nil
}

func (v *memoryVault) encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, v.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	return v.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

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

func clearBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}