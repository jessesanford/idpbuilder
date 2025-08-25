package auth

import (
	"testing"
	"time"
)

func TestCredentials_IsExpired(t *testing.T) {
	past := time.Now().Add(-time.Hour)
	future := time.Now().Add(time.Hour)
	
	if (&Credentials{ExpiresAt: &past}).IsExpired() != true {
		t.Error("Expected expired credentials to return true")
	}
	if (&Credentials{ExpiresAt: &future}).IsExpired() != false {
		t.Error("Expected valid credentials to return false")
	}
}

func TestToken_IsValid(t *testing.T) {
	future := time.Now().Add(time.Hour)
	past := time.Now().Add(-time.Hour)
	
	valid := &Token{AccessToken: "test", ExpiresAt: future}
	expired := &Token{AccessToken: "test", ExpiresAt: past}
	empty := &Token{AccessToken: "", ExpiresAt: future}
	
	if !valid.IsValid() {
		t.Error("Expected valid token to return true")
	}
	if expired.IsValid() {
		t.Error("Expected expired token to return false")
	}
	if empty.IsValid() {
		t.Error("Expected empty token to return false")
	}
}