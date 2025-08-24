package errors

import (
	"errors"
	"testing"
	"time"
)

func TestBaseError_Error(t *testing.T) {
	tests := []struct {
		name     string
		baseErr  *BaseError
		expected string
	}{
		{
			name: "error without cause",
			baseErr: &BaseError{
				ErrorCode:     CodeRegistryUnavailable,
				ErrorCategory: CategoryTransient,
				Message:       "registry is unavailable",
			},
			expected: "registry is unavailable",
		},
		{
			name: "error with cause",
			baseErr: &BaseError{
				ErrorCode:     CodeRegistryUnavailable,
				ErrorCategory: CategoryTransient,
				Message:       "registry is unavailable",
				Cause:         errors.New("connection refused"),
			},
			expected: "registry is unavailable: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.baseErr.Error()
			if got != tt.expected {
				t.Errorf("BaseError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBaseError_Code(t *testing.T) {
	err := &BaseError{
		ErrorCode:     CodeAuthenticationFailed,
		ErrorCategory: CategoryAuthentication,
		Message:       "auth failed",
	}

	if got := err.Code(); got != CodeAuthenticationFailed {
		t.Errorf("BaseError.Code() = %v, want %v", got, CodeAuthenticationFailed)
	}
}

func TestBaseError_Category(t *testing.T) {
	err := &BaseError{
		ErrorCode:     CodeAuthenticationFailed,
		ErrorCategory: CategoryAuthentication,
		Message:       "auth failed",
	}

	if got := err.Category(); got != CategoryAuthentication {
		t.Errorf("BaseError.Category() = %v, want %v", got, CategoryAuthentication)
	}
}

func TestBaseError_Wrap(t *testing.T) {
	baseErr := &BaseError{
		ErrorCode:     CodeRegistryUnavailable,
		ErrorCategory: CategoryTransient,
		Message:       "registry error",
		ErrorContext:  map[string]interface{}{"registry": "docker.io"},
	}

	originalErr := errors.New("connection timeout")
	wrappedErr := baseErr.Wrap(originalErr)

	// Check that wrapped error maintains properties
	if wrappedErr.Code() != CodeRegistryUnavailable {
		t.Errorf("wrapped error code = %v, want %v", wrappedErr.Code(), CodeRegistryUnavailable)
	}

	if wrappedErr.Category() != CategoryTransient {
		t.Errorf("wrapped error category = %v, want %v", wrappedErr.Category(), CategoryTransient)
	}

	if wrappedErr.Unwrap() != originalErr {
		t.Errorf("wrapped error cause = %v, want %v", wrappedErr.Unwrap(), originalErr)
	}

	// Check that context is preserved
	context := wrappedErr.Context()
	if registry, ok := context["registry"]; !ok || registry != "docker.io" {
		t.Errorf("wrapped error context = %v, want registry=docker.io", context)
	}
}

func TestBaseError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	baseErr := &BaseError{
		ErrorCode:     CodeNetworkTimeout,
		ErrorCategory: CategoryNetwork,
		Message:       "timeout occurred",
		Cause:         originalErr,
	}

	if got := baseErr.Unwrap(); got != originalErr {
		t.Errorf("BaseError.Unwrap() = %v, want %v", got, originalErr)
	}
}

func TestBaseError_Context(t *testing.T) {
	context := map[string]interface{}{
		"registry":   "docker.io",
		"repository": "library/nginx",
		"tag":        "latest",
	}

	baseErr := &BaseError{
		ErrorCode:     CodeManifestNotFound,
		ErrorCategory: CategoryPermanent,
		Message:       "manifest not found",
		ErrorContext:  context,
	}

	got := baseErr.Context()
	if len(got) != len(context) {
		t.Errorf("BaseError.Context() length = %v, want %v", len(got), len(context))
	}

	for k, v := range context {
		if got[k] != v {
			t.Errorf("BaseError.Context()[%s] = %v, want %v", k, got[k], v)
		}
	}
}

func TestBaseError_Timestamp(t *testing.T) {
	now := time.Now().UTC()
	baseErr := &BaseError{
		ErrorCode:      CodeInvalidInput,
		ErrorCategory:  CategoryValidation,
		Message:        "validation failed",
		ErrorTimestamp: now,
	}

	if got := baseErr.Timestamp(); !got.Equal(now) {
		t.Errorf("BaseError.Timestamp() = %v, want %v", got, now)
	}
}

func TestBaseError_String(t *testing.T) {
	timestamp := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	baseErr := &BaseError{
		ErrorCode:      CodeRegistryUnavailable,
		ErrorCategory:  CategoryTransient,
		Message:        "registry error",
		ErrorTimestamp: timestamp,
	}

	expected := "OCIError{code=2001, category=transient, message=\"registry error\", timestamp=2023-01-01T12:00:00Z}"
	if got := baseErr.String(); got != expected {
		t.Errorf("BaseError.String() = %v, want %v", got, expected)
	}
}

func TestErrorStack_NewErrorStack(t *testing.T) {
	err := &BaseError{
		ErrorCode:     CodeNetworkTimeout,
		ErrorCategory: CategoryNetwork,
		Message:       "timeout",
	}

	stack := NewErrorStack(err, 5)

	if stack.Depth() != 1 {
		t.Errorf("NewErrorStack depth = %v, want %v", stack.Depth(), 1)
	}

	if stack.MaxDepth != 5 {
		t.Errorf("NewErrorStack max depth = %v, want %v", stack.MaxDepth, 5)
	}

	if stack.Latest() != err {
		t.Errorf("NewErrorStack latest = %v, want %v", stack.Latest(), err)
	}

	if stack.Root() != err {
		t.Errorf("NewErrorStack root = %v, want %v", stack.Root(), err)
	}
}

func TestErrorStack_Push(t *testing.T) {
	err1 := &BaseError{ErrorCode: CodeNetworkTimeout, Message: "timeout"}
	err2 := &BaseError{ErrorCode: CodeRegistryUnavailable, Message: "unavailable"}
	err3 := &BaseError{ErrorCode: CodeManifestNotFound, Message: "not found"}

	stack := NewErrorStack(err1, 2)
	stack.Push(err2)

	if stack.Depth() != 2 {
		t.Errorf("stack depth after push = %v, want %v", stack.Depth(), 2)
	}

	if stack.Latest() != err2 {
		t.Errorf("stack latest after push = %v, want %v", stack.Latest(), err2)
	}

	if stack.Root() != err1 {
		t.Errorf("stack root after push = %v, want %v", stack.Root(), err1)
	}

	// Push third error - should exceed max depth and remove oldest
	stack.Push(err3)

	if stack.Depth() != 2 {
		t.Errorf("stack depth after overflow = %v, want %v", stack.Depth(), 2)
	}

	if stack.Latest().Error() != err3.Error() {
		t.Errorf("stack latest after overflow = %v, want %v", stack.Latest().Error(), err3.Error())
	}

	if stack.Root().Error() != err2.Error() {
		t.Errorf("stack root after overflow = %v, want %v", stack.Root().Error(), err2.Error())
	}
}

func TestErrorStack_Error(t *testing.T) {
	err := &BaseError{
		ErrorCode: CodeInvalidInput,
		Message:   "validation failed",
	}

	stack := NewErrorStack(err, 5)
	if got := stack.Error(); got != err.Error() {
		t.Errorf("ErrorStack.Error() = %v, want %v", got, err.Error())
	}

	// Test empty stack
	emptyStack := &ErrorStack{MaxDepth: 5}
	if got := emptyStack.Error(); got != "empty error stack" {
		t.Errorf("empty ErrorStack.Error() = %v, want %v", got, "empty error stack")
	}
}

func TestErrorContext(t *testing.T) {
	ctx := ErrorContext{
		Operation:    "pull",
		Component:    "registry-client",
		ResourceType: "image",
		ResourceID:   "nginx:latest",
		Registry:     "docker.io",
		Repository:   "library/nginx",
		Tag:          "latest",
		Additional: map[string]interface{}{
			"attempt": 3,
		},
	}

	// Test that all fields are accessible
	if ctx.Operation != "pull" {
		t.Errorf("ErrorContext.Operation = %v, want %v", ctx.Operation, "pull")
	}

	if ctx.Component != "registry-client" {
		t.Errorf("ErrorContext.Component = %v, want %v", ctx.Component, "registry-client")
	}

	if attempt, ok := ctx.Additional["attempt"]; !ok || attempt != 3 {
		t.Errorf("ErrorContext.Additional[attempt] = %v, want %v", attempt, 3)
	}
}