// Copyright 2024 The IDP Builder Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBasicAuth(t *testing.T) {
	registry := "registry.example.com"
	username := "testuser"
	password := "testpass"

	auth := NewBasicAuth(registry, username, password)

	// Test Type
	if auth.Type() != AuthTypeBasic {
		t.Errorf("Expected type %s, got %s", AuthTypeBasic, auth.Type())
	}

	// Test GetCredentials with correct registry
	creds, err := auth.GetCredentials(registry)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if creds.Username != username || creds.Password != password {
		t.Errorf("Expected credentials %s:%s, got %s:%s", username, password, creds.Username, creds.Password)
	}

	// Test GetCredentials with wrong registry
	_, err = auth.GetCredentials("wrong.registry.com")
	if err == nil {
		t.Error("Expected error for wrong registry")
	}

	// Test Validate with valid credentials
	if err := auth.Validate(); err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}

	// Test IsExpired (should always be false for basic auth)
	if auth.IsExpired() {
		t.Error("Basic auth should never expire")
	}

	// Test Refresh (should always succeed for basic auth)
	if err := auth.Refresh(); err != nil {
		t.Errorf("Expected refresh to succeed, got %v", err)
	}
}

func TestBasicAuthValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{"Valid credentials", "user", "pass", false},
		{"Empty username", "", "pass", true},
		{"Empty password", "user", "", true},
		{"Both empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewBasicAuth("registry.com", tt.username, tt.password)
			err := auth.Validate()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBearerAuth(t *testing.T) {
	registry := "registry.example.com"
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	future := time.Now().Add(time.Hour)

	auth := NewBearerAuth(registry, token, &future)

	// Test Type
	if auth.Type() != AuthTypeBearer {
		t.Errorf("Expected type %s, got %s", AuthTypeBearer, auth.Type())
	}

	// Test GetCredentials with correct registry
	creds, err := auth.GetCredentials(registry)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if creds.Token != token {
		t.Errorf("Expected token %s, got %s", token, creds.Token)
	}

	// Test GetCredentials with wrong registry
	_, err = auth.GetCredentials("wrong.registry.com")
	if err == nil {
		t.Error("Expected error for wrong registry")
	}

	// Test Validate with valid token
	if err := auth.Validate(); err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}

	// Test IsExpired with future expiry
	if auth.IsExpired() {
		t.Error("Token should not be expired")
	}

	// Test Refresh (should fail for bearer auth)
	if err := auth.Refresh(); err == nil {
		t.Error("Expected refresh to fail for bearer auth")
	}
}

func TestBearerAuthExpiration(t *testing.T) {
	registry := "registry.example.com"
	token := "expired-token"
	past := time.Now().Add(-time.Hour)

	auth := NewBearerAuth(registry, token, &past)

	// Test IsExpired with past expiry
	if !auth.IsExpired() {
		t.Error("Token should be expired")
	}

	// Test GetCredentials with expired token
	_, err := auth.GetCredentials(registry)
	if err == nil {
		t.Error("Expected error for expired token")
	}
}

func TestBearerAuthValidation(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"Valid token", "valid-token", false},
		{"Empty token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewBearerAuth("registry.com", tt.token, nil)
			err := auth.Validate()
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthError(t *testing.T) {
	// Test AuthError without cause
	err := &AuthError{
		Type:    ErrInvalidCredentials,
		Message: "test error",
	}

	expected := "test error"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}

	// Test AuthError with cause
	cause := &AuthError{Type: ErrTokenExpired, Message: "token expired"}
	err = &AuthError{
		Type:    ErrAuthenticationFailed,
		Message: "auth failed",
		Cause:   cause,
	}

	expected = "auth failed: token expired"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}

	// Test Unwrap
	if err.Unwrap() != cause {
		t.Error("Expected cause to be unwrapped")
	}
}

func TestAuthTypeJSON(t *testing.T) {
	// Test MarshalJSON
	authType := AuthTypeBasic
	data, err := json.Marshal(authType)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	expected := `"basic"`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}

	// Test UnmarshalJSON
	var unmarshaled AuthType
	err = json.Unmarshal([]byte(`"bearer"`), &unmarshaled)
	if err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if unmarshaled != AuthTypeBearer {
		t.Errorf("Expected %s, got %s", AuthTypeBearer, unmarshaled)
	}
}

func TestDockerConfig(t *testing.T) {
	config := DockerConfig{
		Auths: map[string]DockerAuthEntry{
			"registry.example.com": {
				Username: "user",
				Password: "pass",
				Auth:     "dXNlcjpwYXNz",
				Email:    "user@example.com",
			},
		},
		CredStore: "desktop",
		CredHelpers: map[string]string{
			"gcr.io": "gcr",
		},
		HttpHeaders: map[string]string{
			"User-Agent": "docker/1.0",
		},
	}

	// Test that config structure is properly defined
	if len(config.Auths) != 1 {
		t.Error("Expected 1 auth entry")
	}

	entry := config.Auths["registry.example.com"]
	if entry.Username != "user" {
		t.Errorf("Expected username 'user', got %s", entry.Username)
	}

	if config.CredStore != "desktop" {
		t.Errorf("Expected credStore 'desktop', got %s", config.CredStore)
	}

	if len(config.CredHelpers) != 1 {
		t.Error("Expected 1 cred helper")
	}
}

func TestAuthConfigValidation(t *testing.T) {
	// Test valid config
	config := AuthConfig{
		Registry: "registry.example.com",
		Type:     AuthTypeBasic,
		Username: "user",
		Password: "pass",
	}

	// Basic structure validation (would need actual Validate method)
	if config.Registry == "" {
		t.Error("Registry should not be empty")
	}

	if !config.Type.IsValid() {
		t.Error("Auth type should be valid")
	}
}