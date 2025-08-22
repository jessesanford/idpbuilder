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
	"errors"
	"testing"
)

func TestRegistryError(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewAuthError("auth failed", cause)
	
	if err.Type != ErrorTypeAuth {
		t.Errorf("expected auth error type, got %v", err.Type)
	}
	
	if err.IsRetryable() {
		t.Error("auth errors should not be retryable")
	}
	
	if err.Unwrap() != cause {
		t.Error("should unwrap to underlying error")
	}
	
	expected := "authentication error: auth failed (caused by: underlying error)"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorType
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "registry error",
			err:      NewAuthError("test", nil),
			expected: ErrorTypeAuth,
		},
		{
			name:     "generic error",
			err:      errors.New("generic"),
			expected: ErrorTypeRegistry,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClassifyError(tt.err)
			if tt.err == nil {
				if result != nil {
					t.Error("expected nil for nil error")
				}
				return
			}
			
			if result.Type != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.Type)
			}
		})
	}
}