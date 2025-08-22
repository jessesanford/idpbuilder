// Copyright 2024 idpbuilder Contributors
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

package registry

import (
	"testing"
	"time"
)

func TestRegistryConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  RegistryConfig
		wantErr bool
	}{
		{
			name: "valid config with token auth",
			config: RegistryConfig{
				URL: "https://registry.example.com",
				Auth: AuthConfig{
					AuthMethod: AuthMethodToken,
					Token:      "token123",
				},
			},
			wantErr: false,
		},
		{
			name: "valid config with basic auth",
			config: RegistryConfig{
				URL: "https://registry.example.com",
				Auth: AuthConfig{
					AuthMethod: AuthMethodBasic,
					Username:   "user",
					Password:   "pass",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid config - empty URL",
			config: RegistryConfig{
				URL: "",
			},
			wantErr: true,
		},
		{
			name: "invalid config - token auth without token",
			config: RegistryConfig{
				URL: "https://registry.example.com",
				Auth: AuthConfig{
					AuthMethod: AuthMethodToken,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config - basic auth without credentials",
			config: RegistryConfig{
				URL: "https://registry.example.com",
				Auth: AuthConfig{
					AuthMethod: AuthMethodBasic,
					Username:   "user",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RegistryConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultRegistryConfig(t *testing.T) {
	config := DefaultRegistryConfig()
	
	if config.Timeout != 30*time.Second {
		t.Errorf("expected timeout 30s, got %v", config.Timeout)
	}
	
	if config.RetryAttempts != 3 {
		t.Errorf("expected 3 retry attempts, got %d", config.RetryAttempts)
	}
	
	if config.Auth.AuthMethod != AuthMethodNone {
		t.Errorf("expected no auth method, got %v", config.Auth.AuthMethod)
	}
}