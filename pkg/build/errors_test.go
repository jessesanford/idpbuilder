package build

import (
	"errors"
	"testing"
)

func TestBuildError(t *testing.T) {
	originalErr := errors.New("original error")
	buildErr := &BuildError{
		Op:      "test_operation",
		Err:     originalErr,
		Context: "test context",
	}

	expectedMsg := "build error in test_operation: original error (context: test context)"
	if buildErr.Error() != expectedMsg {
		t.Errorf("BuildError.Error() = %v, want %v", buildErr.Error(), expectedMsg)
	}

	if buildErr.Unwrap() != originalErr {
		t.Errorf("BuildError.Unwrap() = %v, want %v", buildErr.Unwrap(), originalErr)
	}
}

func TestBuildErrorWithoutContext(t *testing.T) {
	originalErr := errors.New("original error")
	buildErr := &BuildError{
		Op:  "test_operation",
		Err: originalErr,
	}

	expectedMsg := "build error in test_operation: original error"
	if buildErr.Error() != expectedMsg {
		t.Errorf("BuildError.Error() = %v, want %v", buildErr.Error(), expectedMsg)
	}
}

func TestWrapBuildError(t *testing.T) {
	tests := []struct {
		name     string
		op       string
		err      error
		context  string
		expected string
		wantNil  bool
	}{
		{
			name:     "wrap normal error",
			op:       "build",
			err:      errors.New("test error"),
			context:  "test context",
			expected: "build error in build: test error (context: test context)",
			wantNil:  false,
		},
		{
			name:     "wrap with empty context",
			op:       "finalize",
			err:      errors.New("finalize failed"),
			context:  "",
			expected: "build error in finalize: finalize failed",
			wantNil:  false,
		},
		{
			name:    "wrap nil error",
			op:      "build",
			err:     nil,
			context: "test context",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := WrapBuildError(tt.op, tt.err, tt.context)

			if tt.wantNil {
				if wrapped != nil {
					t.Errorf("WrapBuildError() = %v, want nil", wrapped)
				}
				return
			}

			if wrapped == nil {
				t.Errorf("WrapBuildError() = nil, want error")
				return
			}

			if wrapped.Error() != tt.expected {
				t.Errorf("WrapBuildError().Error() = %v, want %v", wrapped.Error(), tt.expected)
			}

			// Test unwrapping
			if errors.Unwrap(wrapped) != tt.err {
				t.Errorf("errors.Unwrap(WrapBuildError()) = %v, want %v", errors.Unwrap(wrapped), tt.err)
			}
		})
	}
}

func TestIsBuildError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "is build error",
			err:      &BuildError{Op: "test", Err: errors.New("test")},
			expected: true,
		},
		{
			name:     "wrapped build error",
			err:      WrapBuildError("test", errors.New("test"), "context"),
			expected: true,
		},
		{
			name:     "regular error",
			err:      errors.New("regular error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBuildError(tt.err)
			if result != tt.expected {
				t.Errorf("IsBuildError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrInvalidConfig", ErrInvalidConfig, "invalid build configuration"},
		{"ErrBuildFailed", ErrBuildFailed, "build operation failed"},
		{"ErrLayerAddFailed", ErrLayerAddFailed, "failed to add layer"},
		{"ErrFinalizeFailed", ErrFinalizeFailed, "failed to finalize build"},
		{"ErrStorageInit", ErrStorageInit, "failed to initialize storage"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.msg {
				t.Errorf("%s.Error() = %v, want %v", tt.name, tt.err.Error(), tt.msg)
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test that errors.Is works with wrapped BuildErrors
	originalErr := ErrBuildFailed
	wrappedErr := WrapBuildError("test", originalErr, "context")

	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("errors.Is(wrappedErr, originalErr) = false, want true")
	}

	// Test that errors.As works with BuildError
	var buildErr *BuildError
	if !errors.As(wrappedErr, &buildErr) {
		t.Errorf("errors.As(wrappedErr, &buildErr) = false, want true")
	}

	if buildErr.Op != "test" {
		t.Errorf("buildErr.Op = %v, want %v", buildErr.Op, "test")
	}
}