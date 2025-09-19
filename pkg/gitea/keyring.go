package gitea

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "idpbuilder"
	keyringUser    = "gitea"
)

// KeyringProvider provides credentials from system keyring
type KeyringProvider struct {
	service string
	user    string
}

func NewKeyringProvider() *KeyringProvider {
	return &KeyringProvider{
		service: keyringService,
		user:    keyringUser,
	}
}

func (k *KeyringProvider) GetUsername() (string, error) {
	username, err := keyring.Get(k.service, k.user+"_username")
	if err != nil {
		return "", fmt.Errorf("failed to get username from keyring: %w", err)
	}
	return username, nil
}

func (k *KeyringProvider) GetPassword() (string, error) {
	password, err := keyring.Get(k.service, k.user+"_password")
	if err != nil {
		return "", fmt.Errorf("failed to get password from keyring: %w", err)
	}
	return password, nil
}

func (k *KeyringProvider) IsAvailable() bool {
	// Check if keyring is accessible
	_, err := keyring.Get(k.service, k.user+"_username")
	return err == nil || err == keyring.ErrNotFound
}

func (k *KeyringProvider) Priority() int {
	return 4 // Lowest priority
}

// SetCredentials stores credentials in the keyring
func (k *KeyringProvider) SetCredentials(username, password string) error {
	if err := keyring.Set(k.service, k.user+"_username", username); err != nil {
		return fmt.Errorf("failed to store username: %w", err)
	}
	if err := keyring.Set(k.service, k.user+"_password", password); err != nil {
		return fmt.Errorf("failed to store password: %w", err)
	}
	return nil
}

// DeleteCredentials removes credentials from the keyring
func (k *KeyringProvider) DeleteCredentials() error {
	keyring.Delete(k.service, k.user+"_username")
	keyring.Delete(k.service, k.user+"_password")
	return nil
}
