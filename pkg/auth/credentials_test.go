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
	"strings"
	"testing"
	"time"
)

func TestCredentialsValidation(t *testing.T) {
	tests := []struct {
		name    string
		creds   *Credentials
		wantErr bool
	}{
		{
			"Valid username/password",
			&Credentials{Username: "user", Password: "pass"},
			false,
		},
		{
			"Valid token",
			&Credentials{Token: "token123"},
			false,
		},
		{
			"Valid access token",
			&Credentials{AccessToken: "access123"},
			false,
		},
		{
			"Username without password",
			&Credentials{Username: "user"},
			true,
		},
		{
			"Empty credentials",
			&Credentials{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.creds.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCredentialsExpiration(t *testing.T) {
	// Test non-expiring credentials
	creds := &Credentials{Username: "user", Password: "pass"}
	if creds.IsExpired() {
		t.Error("Credentials without expiry should not be expired")
	}

	// Test expired credentials
	past := time.Now().Add(-time.Hour)
	expiredCreds := &Credentials{
		Token:     "token",
		ExpiresAt: &past,
	}
	if !expiredCreds.IsExpired() {
		t.Error("Credentials with past expiry should be expired")
	}

	// Test valid credentials
	future := time.Now().Add(time.Hour)
	validCreds := &Credentials{
		Token:     "token",
		ExpiresAt: &future,
	}
	if validCreds.IsExpired() {
		t.Error("Credentials with future expiry should not be expired")
	}
}

func TestCredentialsNearExpiry(t *testing.T) {
	threshold := 10 * time.Minute
	
	// Test credentials expiring soon
	nearFuture := time.Now().Add(5 * time.Minute)
	creds := &Credentials{
		Token:     "token",
		ExpiresAt: &nearFuture,
	}
	
	if !creds.IsNearExpiry(threshold) {
		t.Error("Credentials should be near expiry")
	}

	// Test credentials not expiring soon
	farFuture := time.Now().Add(time.Hour)
	creds.ExpiresAt = &farFuture
	
	if creds.IsNearExpiry(threshold) {
		t.Error("Credentials should not be near expiry")
	}

	// Test credentials without expiry
	creds.ExpiresAt = nil
	if creds.IsNearExpiry(threshold) {
		t.Error("Credentials without expiry should not be near expiry")
	}
}

func TestCredentialsToBasicAuth(t *testing.T) {
	creds := &Credentials{
		Username: "user",
		Password: "pass",
	}

	basicAuth := creds.ToBasicAuth()
	expected := "Basic dXNlcjpwYXNz" // base64 encoded "user:pass"
	
	if basicAuth != expected {
		t.Errorf("Expected %s, got %s", expected, basicAuth)
	}

	// Test empty credentials
	emptyCreds := &Credentials{}
	if emptyCreds.ToBasicAuth() != "" {
		t.Error("Empty credentials should return empty basic auth")
	}
}

func TestCredentialsToBearerAuth(t *testing.T) {
	// Test with Token field
	creds := &Credentials{Token: "token123"}
	bearerAuth := creds.ToBearerAuth()
	expected := "Bearer token123"
	
	if bearerAuth != expected {
		t.Errorf("Expected %s, got %s", expected, bearerAuth)
	}

	// Test with AccessToken field
	creds = &Credentials{AccessToken: "access123"}
	bearerAuth = creds.ToBearerAuth()
	expected = "Bearer access123"
	
	if bearerAuth != expected {
		t.Errorf("Expected %s, got %s", expected, bearerAuth)
	}

	// Test empty credentials
	emptyCreds := &Credentials{}
	if emptyCreds.ToBearerAuth() != "" {
		t.Error("Empty credentials should return empty bearer auth")
	}
}

func TestCredentialsToAuthHeader(t *testing.T) {
	// Test preference for bearer token
	creds := &Credentials{
		Username:    "user",
		Password:    "pass",
		AccessToken: "access123",
	}

	authHeader := creds.ToAuthHeader()
	expected := "Bearer access123"
	
	if authHeader != expected {
		t.Errorf("Expected bearer auth preference, got %s", authHeader)
	}

	// Test fallback to basic auth
	basicCreds := &Credentials{
		Username: "user",
		Password: "pass",
	}

	authHeader = basicCreds.ToAuthHeader()
	if !strings.HasPrefix(authHeader, "Basic ") {
		t.Errorf("Expected basic auth fallback, got %s", authHeader)
	}
}

func TestCredentialsClone(t *testing.T) {
	expiry := time.Now().Add(time.Hour)
	original := &Credentials{
		Username:     "user",
		Password:     "pass",
		Token:        "token",
		RefreshToken: "refresh",
		AccessToken:  "access",
		TokenType:    "Bearer",
		Scope:        "read",
		ExpiresAt:    &expiry,
	}

	clone := original.Clone()

	// Verify all fields are copied
	if clone.Username != original.Username {
		t.Error("Username not cloned correctly")
	}
	if clone.Password != original.Password {
		t.Error("Password not cloned correctly")
	}
	if clone.Token != original.Token {
		t.Error("Token not cloned correctly")
	}
	if clone.RefreshToken != original.RefreshToken {
		t.Error("RefreshToken not cloned correctly")
	}
	if clone.AccessToken != original.AccessToken {
		t.Error("AccessToken not cloned correctly")
	}
	if clone.TokenType != original.TokenType {
		t.Error("TokenType not cloned correctly")
	}
	if clone.Scope != original.Scope {
		t.Error("Scope not cloned correctly")
	}

	// Verify expiry is deep copied
	if clone.ExpiresAt == original.ExpiresAt {
		t.Error("ExpiresAt should be deep copied, not same pointer")
	}
	if !clone.ExpiresAt.Equal(*original.ExpiresAt) {
		t.Error("ExpiresAt values should be equal")
	}

	// Modify clone to ensure independence
	clone.Username = "modified"
	if original.Username == "modified" {
		t.Error("Modifying clone should not affect original")
	}
}

func TestCredentialsRedacted(t *testing.T) {
	creds := &Credentials{
		Username:     "user",
		Password:     "secret",
		Token:        "token123",
		RefreshToken: "refresh123",
		AccessToken:  "access123",
	}

	redacted := creds.Redacted()

	// Username should not be redacted
	if redacted.Username != "user" {
		t.Error("Username should not be redacted")
	}

	// Sensitive fields should be redacted
	if redacted.Password != "[REDACTED]" {
		t.Errorf("Expected [REDACTED], got %s", redacted.Password)
	}
	if redacted.Token != "[REDACTED]" {
		t.Errorf("Expected [REDACTED], got %s", redacted.Token)
	}
	if redacted.RefreshToken != "[REDACTED]" {
		t.Errorf("Expected [REDACTED], got %s", redacted.RefreshToken)
	}
	if redacted.AccessToken != "[REDACTED]" {
		t.Errorf("Expected [REDACTED], got %s", redacted.AccessToken)
	}
}

func TestCredentialsString(t *testing.T) {
	creds := &Credentials{
		Username:    "user",
		Password:    "secret",
		Token:       "token123",
		TokenType:   "Bearer",
		Scope:       "read",
	}

	str := creds.String()

	// Check that sensitive data is redacted in string representation
	if strings.Contains(str, "secret") || strings.Contains(str, "token123") {
		t.Error("String representation should redact sensitive data")
	}

	// Check that non-sensitive data is included
	if !strings.Contains(str, "user") {
		t.Error("String should contain username")
	}
	if !strings.Contains(str, "Bearer") {
		t.Error("String should contain token type")
	}
	if !strings.Contains(str, "read") {
		t.Error("String should contain scope")
	}
}

func TestCredentialsMarshalJSON(t *testing.T) {
	creds := &Credentials{
		Username: "user",
		Password: "secret",
		Token:    "token123",
	}

	data, err := json.Marshal(creds)
	if err != nil {
		t.Fatalf("JSON marshal error: %v", err)
	}

	// Check that sensitive data is redacted in JSON
	jsonStr := string(data)
	if strings.Contains(jsonStr, "secret") || strings.Contains(jsonStr, "token123") {
		t.Error("JSON should redact sensitive data")
	}
	if !strings.Contains(jsonStr, "[REDACTED]") {
		t.Error("JSON should contain redacted markers")
	}
}

func TestFromBasicAuth(t *testing.T) {
	// Test valid basic auth header
	authHeader := "Basic dXNlcjpwYXNz" // base64 "user:pass"
	creds, err := FromBasicAuth(authHeader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if creds.Username != "user" || creds.Password != "pass" {
		t.Errorf("Expected user:pass, got %s:%s", creds.Username, creds.Password)
	}

	// Test invalid header prefix
	_, err = FromBasicAuth("Bearer token")
	if err == nil {
		t.Error("Expected error for bearer token")
	}

	// Test invalid base64
	_, err = FromBasicAuth("Basic invalid@#$")
	if err == nil {
		t.Error("Expected error for invalid base64")
	}

	// Test invalid format (no colon)
	_, err = FromBasicAuth("Basic dXNlcg==") // base64 "user"
	if err == nil {
		t.Error("Expected error for missing colon")
	}
}

func TestFromBearerAuth(t *testing.T) {
	// Test valid bearer auth header
	authHeader := "Bearer token123"
	creds, err := FromBearerAuth(authHeader)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if creds.Token != "token123" {
		t.Errorf("Expected token123, got %s", creds.Token)
	}
	if creds.TokenType != "Bearer" {
		t.Errorf("Expected Bearer, got %s", creds.TokenType)
	}

	// Test invalid header prefix
	_, err = FromBearerAuth("Basic token")
	if err == nil {
		t.Error("Expected error for basic auth")
	}

	// Test empty token
	_, err = FromBearerAuth("Bearer ")
	if err == nil {
		t.Error("Expected error for empty token")
	}
}

func TestCredentialStore(t *testing.T) {
	store := NewCredentialStore()

	registry := "registry.example.com"
	creds := &Credentials{
		Username: "user",
		Password: "pass",
	}

	// Test Set
	err := store.Set(registry, creds)
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Test Get
	retrieved, err := store.Get(registry)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if retrieved.Username != creds.Username {
		t.Error("Retrieved credentials don't match stored credentials")
	}

	// Test List
	registries, err := store.List()
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(registries) != 1 || registries[0] != registry {
		t.Error("List should return stored registry")
	}

	// Test Delete
	err = store.Delete(registry)
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	// Verify deletion
	_, err = store.Get(registry)
	if err == nil {
		t.Error("Expected error after deletion")
	}

	// Test Clear
	store.Set("reg1", creds)
	store.Set("reg2", creds)
	err = store.Clear()
	if err != nil {
		t.Fatalf("Clear error: %v", err)
	}

	registries, _ = store.List()
	if len(registries) != 0 {
		t.Error("Store should be empty after clear")
	}
}

func TestCredentialStoreValidation(t *testing.T) {
	store := NewCredentialStore()

	// Test empty registry
	err := store.Set("", &Credentials{Username: "user", Password: "pass"})
	if err == nil {
		t.Error("Expected error for empty registry")
	}

	// Test nil credentials
	err = store.Set("registry.com", nil)
	if err == nil {
		t.Error("Expected error for nil credentials")
	}

	// Test invalid credentials
	err = store.Set("registry.com", &Credentials{}) // empty credentials
	if err == nil {
		t.Error("Expected error for invalid credentials")
	}
}

func TestCredentialStoreExpiredCredentials(t *testing.T) {
	store := NewCredentialStore()
	registry := "registry.example.com"
	
	// Store expired credentials
	past := time.Now().Add(-time.Hour)
	expiredCreds := &Credentials{
		Token:     "token",
		ExpiresAt: &past,
	}

	err := store.Set(registry, expiredCreds)
	if err != nil {
		t.Fatalf("Set error: %v", err)
	}

	// Try to retrieve expired credentials
	_, err = store.Get(registry)
	if err == nil {
		t.Error("Expected error for expired credentials")
	}
}

func TestTokenResponse(t *testing.T) {
	response := &TokenResponse{
		AccessToken:  "access123",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: "refresh123",
		Scope:        "read write",
		IssuedAt:     time.Now(),
	}

	creds := response.ToCredentials()

	if creds.AccessToken != response.AccessToken {
		t.Error("Access token not converted correctly")
	}
	if creds.TokenType != response.TokenType {
		t.Error("Token type not converted correctly")
	}
	if creds.RefreshToken != response.RefreshToken {
		t.Error("Refresh token not converted correctly")
	}
	if creds.Scope != response.Scope {
		t.Error("Scope not converted correctly")
	}

	// Check expiry calculation
	if creds.ExpiresAt == nil {
		t.Error("ExpiresAt should be set")
	}
	expectedExpiry := response.IssuedAt.Add(time.Duration(response.ExpiresIn) * time.Second)
	if !creds.ExpiresAt.Equal(expectedExpiry) {
		t.Error("ExpiresAt calculation is incorrect")
	}
}