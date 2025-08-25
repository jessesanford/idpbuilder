package errors

import (
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	code := ErrorNotFound
	message := "test resource not found"
	
	err := New(code, message)
	
	if err == nil {
		t.Fatal("Expected error to be created")
	}
	
	if err.Error() != message {
		t.Errorf("Expected error message %q, got %q", message, err.Error())
	}
	
	if err.Code() != code {
		t.Errorf("Expected error code %+v, got %+v", code, err.Code())
	}
	
	if !err.IsPermanent() {
		t.Error("Expected error to be permanent")
	}
	
	if err.IsTransient() {
		t.Error("Expected error to not be transient")
	}
}

func TestNewf(t *testing.T) {
	code := ErrorInvalidInput
	format := "invalid input for operation %s: %d"
	operation := "create"
	value := 123
	
	err := Newf(code, format, operation, value)
	
	expectedMessage := "invalid input for operation create: 123"
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, err.Error())
	}
	
	if err.Code() != code {
		t.Errorf("Expected error code %+v, got %+v", code, err.Code())
	}
}

func TestWrap(t *testing.T) {
	original := errors.New("original error")
	code := ErrorNetworkTimeout
	message := "network operation failed"
	
	wrapped := Wrap(code, message, original)
	
	expectedMessage := "network operation failed: original error"
	if wrapped.Error() != expectedMessage {
		t.Errorf("Expected wrapped error message %q, got %q", expectedMessage, wrapped.Error())
	}
	
	if wrapped.Code() != code {
		t.Errorf("Expected error code %+v, got %+v", code, wrapped.Code())
	}
	
	if !wrapped.IsTransient() {
		t.Error("Expected error to be transient")
	}
	
	if wrapped.Unwrap() != original {
		t.Error("Expected unwrapped error to match original")
	}
	
	// Test wrapping another OCIError
	firstErr := New(ErrorInvalidConfig, "first error")
	secondErr := Wrap(ErrorNetworkTimeout, "second error", firstErr)
	
	if secondErr.Stack().Depth != 1 {
		t.Errorf("Expected stack depth 1, got %d", secondErr.Stack().Depth)
	}
}

func TestWrapf(t *testing.T) {
	original := errors.New("connection refused")
	code := ErrorNetworkUnavailable
	format := "failed to connect to %s:%d"
	host := "localhost"
	port := 8080
	
	wrapped := Wrapf(code, original, format, host, port)
	
	expectedMessage := "failed to connect to localhost:8080: connection refused"
	if wrapped.Error() != expectedMessage {
		t.Errorf("Expected wrapped error message %q, got %q", expectedMessage, wrapped.Error())
	}
}

func TestBaseError_Context(t *testing.T) {
	err := New(ErrorNotFound, "resource not found").(*BaseError)
	
	context := err.Context()
	if context.Timestamp.IsZero() {
		t.Error("Expected context to have timestamp")
	}
	
	if context.Custom == nil {
		context.Custom = make(map[string]string)
	}
}

func TestBaseError_WithContext(t *testing.T) {
	err := New(ErrorNotFound, "resource not found")
	
	updatedErr := err.WithContext("resource_type", "pod")
	
	context := updatedErr.Context()
	if context.Custom["resource_type"] != "pod" {
		t.Error("Expected custom context to be set")
	}
}

func TestBaseError_WithOperation(t *testing.T) {
	err := New(ErrorNotFound, "resource not found")
	
	updatedErr := err.WithOperation("delete_pod")
	
	context := updatedErr.Context()
	if context.Operation != "delete_pod" {
		t.Errorf("Expected operation to be 'delete_pod', got %q", context.Operation)
	}
}

func TestBaseError_WithResource(t *testing.T) {
	err := New(ErrorNotFound, "resource not found")
	
	updatedErr := err.WithResource("pod/my-pod")
	
	context := updatedErr.Context()
	if context.Resource != "pod/my-pod" {
		t.Errorf("Expected resource to be 'pod/my-pod', got %q", context.Resource)
	}
}

func TestBaseError_Severity(t *testing.T) {
	err := New(ErrorResourceCorrupted, "data corruption detected")
	
	if err.Severity() != SeverityCritical {
		t.Errorf("Expected severity to be %q, got %q", SeverityCritical, err.Severity())
	}
}

func TestBaseError_Stack(t *testing.T) {
	originalErr := New(ErrorInvalidConfig, "config error")
	wrappedErr := Wrap(ErrorNetworkTimeout, "network error", originalErr)
	
	stack := wrappedErr.Stack()
	if stack.Depth != 1 {
		t.Errorf("Expected stack depth 1, got %d", stack.Depth)
	}
	
	if len(stack.Errors) != 1 {
		t.Errorf("Expected 1 stack frame, got %d", len(stack.Errors))
	}
	
	frame := stack.Errors[0]
	if frame.Error != "network error" {
		t.Errorf("Expected frame error 'network error', got %q", frame.Error)
	}
	
	if frame.Code != ErrorNetworkTimeout {
		t.Errorf("Expected frame code %+v, got %+v", ErrorNetworkTimeout, frame.Code)
	}
}

func TestIsOCIError(t *testing.T) {
	ociErr := New(ErrorNotFound, "not found")
	stdErr := errors.New("standard error")
	
	if !IsOCIError(ociErr) {
		t.Error("Expected OCIError to be identified")
	}
	
	if IsOCIError(stdErr) {
		t.Error("Expected standard error to not be identified as OCIError")
	}
}

func TestAsOCIError(t *testing.T) {
	ociErr := New(ErrorNotFound, "not found")
	stdErr := errors.New("standard error")
	
	convertedErr, ok := AsOCIError(ociErr)
	if !ok {
		t.Error("Expected OCIError to be convertible")
	}
	if convertedErr != ociErr {
		t.Error("Expected converted error to match original")
	}
	
	_, ok = AsOCIError(stdErr)
	if ok {
		t.Error("Expected standard error to not be convertible to OCIError")
	}
}

func TestGetRootCause(t *testing.T) {
	root := errors.New("root cause")
	middle := Wrap(ErrorNetworkTimeout, "middle error", root)
	top := Wrap(ErrorServiceBusy, "top error", middle)
	
	rootCause := GetRootCause(top)
	if rootCause != root {
		t.Error("Expected root cause to be found")
	}
	
	// Test with non-wrapping error
	single := errors.New("single error")
	singleRoot := GetRootCause(single)
	if singleRoot != single {
		t.Error("Expected single error to be its own root cause")
	}
}

func TestUnwrap(t *testing.T) {
	original := errors.New("original")
	wrapped := Wrap(ErrorNotFound, "wrapped", original)
	
	unwrapped := Unwrap(wrapped)
	if unwrapped != original {
		t.Error("Expected unwrapped error to match original")
	}
	
	// Test unwrapping non-wrapped error
	unwrappedNil := Unwrap(original)
	if unwrappedNil != nil {
		t.Error("Expected nil when unwrapping non-wrapped error")
	}
}

func TestErrorChaining(t *testing.T) {
	// Create a chain of errors
	root := errors.New("database connection failed")
	dbErr := Wrap(ErrorNetworkTimeout, "database timeout", root)
	serviceErr := Wrap(ErrorServiceBusy, "service unavailable", dbErr)
	apiErr := Wrap(ErrorRateLimited, "API rate limit exceeded", serviceErr)
	
	// Test that each level can be unwrapped
	level3 := Unwrap(apiErr)
	if level3 != serviceErr {
		t.Error("Expected level 3 to unwrap to service error")
	}
	
	level2 := Unwrap(level3)
	if level2 != dbErr {
		t.Error("Expected level 2 to unwrap to database error")
	}
	
	level1 := Unwrap(level2)
	if level1 != root {
		t.Error("Expected level 1 to unwrap to root error")
	}
	
	// Test root cause detection
	rootCause := GetRootCause(apiErr)
	if rootCause != root {
		t.Error("Expected root cause to be original error")
	}
}

func TestErrorContext_Timestamp(t *testing.T) {
	before := time.Now()
	err := New(ErrorNotFound, "test error")
	after := time.Now()
	
	context := err.Context()
	if context.Timestamp.Before(before) || context.Timestamp.After(after) {
		t.Error("Expected timestamp to be within test execution window")
	}
}

func TestErrorContext_CustomFields(t *testing.T) {
	err := New(ErrorNotFound, "test error")
	
	// Test multiple context additions
	err.WithContext("key1", "value1")
	err.WithContext("key2", "value2")
	err.WithOperation("test_operation")
	err.WithResource("test/resource")
	
	context := err.Context()
	
	if context.Custom["key1"] != "value1" {
		t.Error("Expected custom context key1 to be set")
	}
	
	if context.Custom["key2"] != "value2" {
		t.Error("Expected custom context key2 to be set")
	}
	
	if context.Operation != "test_operation" {
		t.Error("Expected operation context to be set")
	}
	
	if context.Resource != "test/resource" {
		t.Error("Expected resource context to be set")
	}
}

func TestStackFrame_Fields(t *testing.T) {
	originalErr := New(ErrorInvalidInput, "input validation failed")
	wrappedErr := Wrap(ErrorServiceBusy, "service overwhelmed", originalErr)
	
	stack := wrappedErr.Stack()
	if len(stack.Errors) != 1 {
		t.Fatalf("Expected 1 stack frame, got %d", len(stack.Errors))
	}
	
	frame := stack.Errors[0]
	
	if frame.Error != "service overwhelmed" {
		t.Errorf("Expected frame error 'service overwhelmed', got %q", frame.Error)
	}
	
	if frame.Code != ErrorServiceBusy {
		t.Errorf("Expected frame code %+v, got %+v", ErrorServiceBusy, frame.Code)
	}
	
	if frame.Timestamp.IsZero() {
		t.Error("Expected frame to have timestamp")
	}
}

func TestComplexErrorScenario(t *testing.T) {
	// Simulate a complex error scenario with multiple wrapping levels
	// and various context information
	
	// Start with a low-level system error
	systemErr := errors.New("connection reset by peer")
	
	// Wrap with network error
	networkErr := Wrap(ErrorNetworkTimeout, "TCP connection failed", systemErr).
		WithOperation("connect").
		WithContext("host", "registry.example.com").
		WithContext("port", "443")
	
	// Wrap with authentication error
	authErr := Wrap(ErrorUnauthenticated, "registry authentication failed", networkErr).
		WithOperation("authenticate").
		WithResource("registry/repository")
	
	// Wrap with high-level operation error
	pullErr := Wrap(ErrorNotFound, "image pull failed", authErr).
		WithOperation("pull_image").
		WithResource("myapp:latest").
		WithContext("registry", "registry.example.com/myapp")
	
	// Verify error message includes full chain
	errorMsg := pullErr.Error()
	expectedParts := []string{"image pull failed", "registry authentication failed", "TCP connection failed", "connection reset by peer"}
	for _, part := range expectedParts {
		if !contains(errorMsg, part) {
			t.Errorf("Expected error message to contain %q, got %q", part, errorMsg)
		}
	}
	
	// Verify context is preserved
	context := pullErr.Context()
	if context.Operation != "pull_image" {
		t.Errorf("Expected operation 'pull_image', got %q", context.Operation)
	}
	
	if context.Resource != "myapp:latest" {
		t.Errorf("Expected resource 'myapp:latest', got %q", context.Resource)
	}
	
	if context.Custom["registry"] != "registry.example.com/myapp" {
		t.Error("Expected registry context to be preserved")
	}
	
	// Verify error stack
	stack := pullErr.Stack()
	if stack.Depth < 1 {
		t.Errorf("Expected stack depth >= 1, got %d", stack.Depth)
	}
	
	// Verify root cause
	rootCause := GetRootCause(pullErr)
	if rootCause != systemErr {
		t.Error("Expected root cause to be system error")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	if s[:len(substr)] == substr {
		return true
	}
	return containsHelper(s[1:], substr)
}