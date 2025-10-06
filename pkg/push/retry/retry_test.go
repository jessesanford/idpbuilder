package retry

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestWithRetry_Success(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 3

	callCount := 0
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		return nil // Succeed immediately
	}

	err := WithRetry(ctx, config, fn)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestWithRetry_SuccessAfterRetries(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 3
	config.BackoffStrategy = NewConstantBackoff(10 * time.Millisecond)

	callCount := 0
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		if callCount < 3 {
			return &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
		}
		return nil // Succeed on third call
	}

	err := WithRetry(ctx, config, fn)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}
}

func TestWithRetry_MaxRetriesExceeded(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 3
	config.BackoffStrategy = NewConstantBackoff(10 * time.Millisecond)

	callCount := 0
	testErr := &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		return testErr
	}

	err := WithRetry(ctx, config, fn)

	var maxErr *MaxRetriesExceededError
	if !errors.As(err, &maxErr) {
		t.Errorf("expected MaxRetriesExceededError, got: %v", err)
	}

	if callCount != 3 {
		t.Errorf("expected 3 calls, got %d", callCount)
	}

	if maxErr.Attempts != 3 {
		t.Errorf("expected 3 attempts in error, got %d", maxErr.Attempts)
	}

	if !errors.Is(err, testErr) {
		t.Error("expected wrapped error to be testErr")
	}
}

func TestWithRetry_NonRetryableError(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 3

	callCount := 0
	testErr := errors.New("permanent error")
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		return testErr
	}

	err := WithRetry(ctx, config, fn)

	if err != testErr {
		t.Errorf("expected testErr, got: %v", err)
	}

	if callCount != 1 {
		t.Errorf("expected 1 call (no retries for non-retryable error), got %d", callCount)
	}
}

func TestWithRetry_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := DefaultConfig()
	config.MaxAttempts = 5
	config.BackoffStrategy = NewConstantBackoff(100 * time.Millisecond)

	callCount := 0
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		if callCount == 2 {
			cancel() // Cancel after second attempt
		}
		return &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
	}

	err := WithRetry(ctx, config, fn)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}

	// Should stop after cancellation
	if callCount > 3 {
		t.Errorf("expected <= 3 calls after cancel, got %d", callCount)
	}
}

func TestWithRetry_InvalidMaxAttempts(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 0

	fn := func(ctx context.Context, attempt int) error {
		return nil
	}

	err := WithRetry(ctx, config, fn)
	if err == nil || !strings.Contains(err.Error(), "max attempts must be at least 1") {
		t.Errorf("expected validation error for invalid max attempts, got: %v", err)
	}
}

func TestWithRetry_NilConfig(t *testing.T) {
	ctx := context.Background()

	callCount := 0
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		return nil
	}

	err := WithRetry(ctx, nil, fn)
	if err != nil {
		t.Errorf("expected no error with nil config (should use defaults), got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestWithRetry_CustomShouldRetry(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 3
	config.BackoffStrategy = NewConstantBackoff(10 * time.Millisecond)

	// Custom retry logic: only retry errors containing "transient"
	config.ShouldRetry = func(err error) bool {
		return strings.Contains(err.Error(), "transient")
	}

	callCount := 0
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		if callCount < 2 {
			return errors.New("transient error")
		}
		return errors.New("permanent error")
	}

	err := WithRetry(ctx, config, fn)

	if !strings.Contains(err.Error(), "permanent error") {
		t.Errorf("expected 'permanent error', got: %v", err)
	}

	// Should call twice: once successful retry, then fail on non-retryable
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestIsRetryable_ContextErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"context canceled", context.Canceled, false},
		{"context deadline exceeded", context.DeadlineExceeded, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for error: %v", tt.expected, result, tt.err)
			}
		})
	}
}

func TestIsRetryable_NetworkErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			"connection refused",
			&net.OpError{Op: "dial", Err: syscall.ECONNREFUSED},
			true,
		},
		{
			"connection reset",
			&net.OpError{Op: "read", Err: syscall.ECONNRESET},
			true,
		},
		{
			"broken pipe",
			&net.OpError{Op: "write", Err: syscall.EPIPE},
			true,
		},
		{
			"timeout",
			&net.OpError{Op: "dial", Err: syscall.ETIMEDOUT},
			true,
		},
		{
			"error message with connection refused",
			errors.New("dial tcp: connection refused"),
			true,
		},
		{
			"error message with i/o timeout",
			errors.New("read tcp: i/o timeout"),
			true,
		},
		{
			"error message with EOF",
			errors.New("unexpected EOF"),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for error: %v", tt.expected, result, tt.err)
			}
		})
	}
}

func TestIsRetryable_HTTPErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"HTTP 408", errors.New("server returned 408"), true},
		{"HTTP 429", errors.New("server returned 429 too many requests"), true},
		{"HTTP 500", errors.New("server returned 500 internal server error"), true},
		{"HTTP 502", errors.New("server returned 502 bad gateway"), true},
		{"HTTP 503", errors.New("server returned 503 service unavailable"), true},
		{"HTTP 504", errors.New("server returned 504 gateway timeout"), true},
		{"HTTP 200", errors.New("server returned 200 ok"), false},
		{"HTTP 404", errors.New("server returned 404 not found"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for error: %v", tt.expected, result, tt.err)
			}
		})
	}
}

func TestIsRetryable_NonRetryableErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"generic error", errors.New("something went wrong")},
		{"validation error", errors.New("invalid input")},
		{"not found", errors.New("resource not found")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryable(tt.err)
			if result {
				t.Errorf("expected false for non-retryable error: %v", tt.err)
			}
		})
	}
}

func TestWithRetrySimple_Success(t *testing.T) {
	ctx := context.Background()

	callCount := 0
	fn := func() error {
		callCount++
		return nil
	}

	err := WithRetrySimple(ctx, fn)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}
}

func TestWithRetryN_CustomAttempts(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 5

	callCount := 0
	fn := func(ctx context.Context, attempt int) error {
		callCount++
		return &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
	}

	err := WithRetryN(ctx, maxAttempts, fn)

	var maxErr *MaxRetriesExceededError
	if !errors.As(err, &maxErr) {
		t.Errorf("expected MaxRetriesExceededError, got: %v", err)
	}

	if callCount != maxAttempts {
		t.Errorf("expected %d calls, got %d", maxAttempts, callCount)
	}
}

func TestRetryIfTransient_OnlyRetriesNetworkErrors(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 3

	t.Run("retries network error", func(t *testing.T) {
		callCount := 0
		fn := func() error {
			callCount++
			if callCount < 2 {
				return &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
			}
			return nil
		}

		err := RetryIfTransient(ctx, maxAttempts, fn)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if callCount != 2 {
			t.Errorf("expected 2 calls, got %d", callCount)
		}
	})

	t.Run("does not retry non-network error", func(t *testing.T) {
		callCount := 0
		testErr := errors.New("not a network error")
		fn := func() error {
			callCount++
			return testErr
		}

		err := RetryIfTransient(ctx, maxAttempts, fn)
		if err != testErr {
			t.Errorf("expected testErr, got: %v", err)
		}
		if callCount != 1 {
			t.Errorf("expected 1 call (no retry), got %d", callCount)
		}
	})
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.MaxAttempts != 3 {
		t.Errorf("expected default MaxAttempts 3, got %d", config.MaxAttempts)
	}
	if config.BackoffStrategy == nil {
		t.Error("expected BackoffStrategy to be initialized")
	}
	if config.ShouldRetry != nil {
		t.Error("expected ShouldRetry to be nil (use default IsRetryable)")
	}
}

// mockTemporaryError simulates a net.Error with temporary status
type mockTemporaryError struct {
	temporary bool
	timeout   bool
}

func (m *mockTemporaryError) Error() string   { return "mock error" }
func (m *mockTemporaryError) Temporary() bool { return m.temporary }
func (m *mockTemporaryError) Timeout() bool   { return m.timeout }

func TestIsRetryable_TemporaryNetError(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		expected  bool
	}{
		{
			"temporary error",
			&mockTemporaryError{temporary: true, timeout: false},
			true,
		},
		{
			"timeout error",
			&mockTemporaryError{temporary: false, timeout: true},
			true,
		},
		{
			"both temporary and timeout",
			&mockTemporaryError{temporary: true, timeout: true},
			true,
		},
		{
			"neither temporary nor timeout",
			&mockTemporaryError{temporary: false, timeout: false},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRetryable(tt.err)
			if result != tt.expected {
				t.Errorf("expected %v, got %v for error: %v", tt.expected, result, tt.err)
			}
		})
	}
}

func TestWithRetry_AttemptNumber(t *testing.T) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 3
	config.BackoffStrategy = NewConstantBackoff(10 * time.Millisecond)

	attempts := []int{}
	fn := func(ctx context.Context, attempt int) error {
		attempts = append(attempts, attempt)
		return &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
	}

	_ = WithRetry(ctx, config, fn)

	expected := []int{0, 1, 2}
	if len(attempts) != len(expected) {
		t.Errorf("expected %d attempts, got %d", len(expected), len(attempts))
	}

	for i, exp := range expected {
		if attempts[i] != exp {
			t.Errorf("attempt %d: expected %d, got %d", i, exp, attempts[i])
		}
	}
}

func BenchmarkWithRetry_NoRetries(b *testing.B) {
	ctx := context.Background()
	config := DefaultConfig()

	fn := func(ctx context.Context, attempt int) error {
		return nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WithRetry(ctx, config, fn)
	}
}

func BenchmarkWithRetry_WithRetries(b *testing.B) {
	ctx := context.Background()
	config := DefaultConfig()
	config.MaxAttempts = 3
	config.BackoffStrategy = NewConstantBackoff(1 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		callCount := 0
		fn := func(ctx context.Context, attempt int) error {
			callCount++
			if callCount < 3 {
				return &net.OpError{Op: "dial", Err: syscall.ECONNREFUSED}
			}
			return nil
		}
		_ = WithRetry(ctx, config, fn)
	}
}

func ExampleWithRetry() {
	ctx := context.Background()
	config := DefaultConfig()

	err := WithRetry(ctx, config, func(ctx context.Context, attempt int) error {
		fmt.Printf("Attempt %d\n", attempt+1)
		// Simulate work that might fail
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	}
	// Output: Attempt 1
}

func ExampleWithRetrySimple() {
	ctx := context.Background()

	err := WithRetrySimple(ctx, func() error {
		// Simulate operation
		return nil
	})

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
	}
}