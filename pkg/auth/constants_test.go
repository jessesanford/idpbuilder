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
	"testing"
	"time"
)

func TestAuthTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		authType AuthType
		expected string
	}{
		{"Basic auth type", AuthTypeBasic, "basic"},
		{"Bearer auth type", AuthTypeBearer, "bearer"},
		{"OAuth2 auth type", AuthTypeOAuth2, "oauth2"},
		{"Anonymous auth type", AuthTypeAnonymous, "anonymous"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.authType) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(tt.authType))
			}
		})
	}
}

func TestHTTPHeaderConstants(t *testing.T) {
	if AuthorizationHeader != "Authorization" {
		t.Errorf("Expected Authorization, got %s", AuthorizationHeader)
	}

	if WWWAuthenticateHeader != "WWW-Authenticate" {
		t.Errorf("Expected WWW-Authenticate, got %s", WWWAuthenticateHeader)
	}

	if BasicAuthPrefix != "Basic " {
		t.Errorf("Expected 'Basic ', got '%s'", BasicAuthPrefix)
	}

	if BearerAuthPrefix != "Bearer " {
		t.Errorf("Expected 'Bearer ', got '%s'", BearerAuthPrefix)
	}
}

func TestDefaultTimeConstants(t *testing.T) {
	if DefaultTokenExpiry != time.Hour {
		t.Errorf("Expected 1 hour, got %v", DefaultTokenExpiry)
	}

	if DefaultRefreshThreshold != 5*time.Minute {
		t.Errorf("Expected 5 minutes, got %v", DefaultRefreshThreshold)
	}

	if MaxRetryAttempts != 3 {
		t.Errorf("Expected 3, got %d", MaxRetryAttempts)
	}
}

func TestRegistryConstants(t *testing.T) {
	if DockerHubRegistry != "https://index.docker.io/v1/" {
		t.Errorf("Expected Docker Hub registry URL, got %s", DockerHubRegistry)
	}

	if DockerHubHost != "index.docker.io" {
		t.Errorf("Expected Docker Hub host, got %s", DockerHubHost)
	}

	if HTTPScheme != "http://" {
		t.Errorf("Expected http://, got %s", HTTPScheme)
	}

	if HTTPSScheme != "https://" {
		t.Errorf("Expected https://, got %s", HTTPSScheme)
	}
}

func TestDockerConfigConstants(t *testing.T) {
	constants := []struct {
		name     string
		value    string
		expected string
	}{
		{"DockerConfigFilename", DockerConfigFilename, "config.json"},
		{"DockerConfigDir", DockerConfigDir, ".docker"},
		{"DockerCredentialStore", DockerCredentialStore, "credStore"},
		{"DockerCredentialHelpers", DockerCredentialHelpers, "credHelpers"},
		{"DockerAuths", DockerAuths, "auths"},
	}

	for _, c := range constants {
		t.Run(c.name, func(t *testing.T) {
			if c.value != c.expected {
				t.Errorf("Expected %q, got %q", c.expected, c.value)
			}
		})
	}
}

func TestErrorMessageConstants(t *testing.T) {
	errorMessages := []string{
		ErrInvalidCredentials,
		ErrTokenExpired,
		ErrAuthenticationFailed,
		ErrUnsupportedAuthType,
		ErrMissingCredentials,
		ErrRegistryUnreachable,
		ErrPermissionDenied,
		ErrInvalidToken,
	}

	for _, msg := range errorMessages {
		if msg == "" {
			t.Error("Error message constant should not be empty")
		}
	}
}

func TestAuthTypeString(t *testing.T) {
	tests := []struct {
		authType AuthType
		expected string
	}{
		{AuthTypeBasic, "basic"},
		{AuthTypeBearer, "bearer"},
		{AuthTypeOAuth2, "oauth2"},
		{AuthTypeAnonymous, "anonymous"},
	}

	for _, tt := range tests {
		if tt.authType.String() != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, tt.authType.String())
		}
	}
}

func TestAuthTypeIsValid(t *testing.T) {
	validTypes := []AuthType{
		AuthTypeBasic,
		AuthTypeBearer,
		AuthTypeOAuth2,
		AuthTypeAnonymous,
	}

	for _, authType := range validTypes {
		if !authType.IsValid() {
			t.Errorf("Expected %q to be valid", authType)
		}
	}

	invalidTypes := []AuthType{
		AuthType("invalid"),
		AuthType(""),
		AuthType("BASIC"), // case sensitive
	}

	for _, authType := range invalidTypes {
		if authType.IsValid() {
			t.Errorf("Expected %q to be invalid", authType)
		}
	}
}