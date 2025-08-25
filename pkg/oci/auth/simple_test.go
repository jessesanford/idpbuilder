package auth

import (
	"testing"
	"time"
)

func TestSimpleValidation(t *testing.T) {
	// Test basic credential validation without complex dependencies
	now := time.Now()
	future := now.Add(1 * time.Hour)
	
	// Test basic auth
	creds := &Credentials{
		Registry:   "registry.example.com",
		Username:   "testuser",
		Password:   "testpass",
		AuthMethod: AuthMethodBasic,
		CreatedAt:  now,
	}
	
	err := ValidateCredentials(creds)
	if err != nil {
		t.Errorf("valid basic credentials should not return error: %v", err)
	}
	
	// Test token auth
	tokenCreds := &Credentials{
		Registry:   "registry.example.com",
		AuthMethod: AuthMethodToken,
		Token: &Token{
			Value:     "test-token",
			Type:      TokenTypeBearer,
			IssuedAt:  now,
			ExpiresAt: future,
		},
		CreatedAt: now,
	}
	
	err = ValidateCredentials(tokenCreds)
	if err != nil {
		t.Errorf("valid token credentials should not return error: %v", err)
	}
	
	// Test invalid registry
	err = ValidateRegistryURL("invalid registry!")
	if err == nil {
		t.Error("invalid registry should return error")
	}
	
	// Test valid registry
	err = ValidateRegistryURL("registry.example.com:5000")
	if err != nil {
		t.Errorf("valid registry should not return error: %v", err)
	}
}

func TestTypeConstants(t *testing.T) {
	// Test that our constants are properly defined
	if AuthMethodBasic != "basic" {
		t.Errorf("expected AuthMethodBasic to be 'basic', got %s", AuthMethodBasic)
	}
	
	if TokenTypeBearer != "Bearer" {
		t.Errorf("expected TokenTypeBearer to be 'Bearer', got %s", TokenTypeBearer)
	}
	
	if StoreTypeFile != "file" {
		t.Errorf("expected StoreTypeFile to be 'file', got %s", StoreTypeFile)
	}
}