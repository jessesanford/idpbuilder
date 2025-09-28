package errors

import (
	"context"
	"testing"
	"time"
)

func TestNewRetryHandler(t *testing.T) {
	config := RetryInfo{
		MaxAttempts: 5,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  1.5,
		Jitter:      false,
	}

	handler := NewRetryHandler(config)

	if handler == nil {
		t.Fatal("NewRetryHandler() returned nil")
	}

	if handler.config != config {
		t.Errorf("NewRetryHandler() config = %v, want %v", handler.config, config)
	}
}

func TestNewDefaultRetryHandler(t *testing.T) {
	handler := NewDefaultRetryHandler()

	if handler == nil {
		t.Fatal("NewDefaultRetryHandler() returned nil")
	}

	expected := DefaultRetryInfo()
	if handler.config != expected {
		t.Errorf("NewDefaultRetryHandler() config = %v, want %v", handler.config, expected)
	}
}

func TestRetryHandler_Retry_Success(t *testing.T) {
	handler := NewDefaultRetryHandler()
	ctx := context.Background()

	callCount := 0
	fn := func() error {
		callCount++
		return nil // Success on first attempt
	}

	err := handler.Retry(ctx, "test_operation", fn)

	if err != nil {
		t.Errorf("Retry() = %v, want nil", err)
	}

	if callCount != 1 {
		t.Errorf("Function called %d times, want 1", callCount)
	}
}

func TestRetryHandler_Retry_SuccessAfterRetries(t *testing.T) {
	handler := NewRetryHandler(RetryInfo{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    10 * time.Millisecond,
		Multiplier:  2.0,
		Jitter:      false,
	})
	ctx := context.Background()

	callCount := 0
	fn := func() error {
		callCount++
		if callCount < 3 {
			return ErrTimeout // Retryable error
		}
		return nil // Success on third attempt
	}

	err := handler.Retry(ctx, "test_operation", fn)

	if err != nil {
		t.Errorf("Retry() = %v, want nil", err)
	}

	if callCount != 3 {
		t.Errorf("Function called %d times, want 3", callCount)
	}
}

func TestRetryHandler_Retry_MaxAttemptsExceeded(t *testing.T) {
	handler := NewRetryHandler(RetryInfo{
		MaxAttempts: 2,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    10 * time.Millisecond,
		Multiplier:  2.0,
		Jitter:      false,
	})
	ctx := context.Background()

	callCount := 0
	expectedErr := ErrTimeout
	fn := func() error {
		callCount++
		return expectedErr // Always fail with retryable error
	}

	err := handler.Retry(ctx, "test_operation", fn)

	if err != expectedErr {
		t.Errorf("Retry() = %v, want %v", err, expectedErr)
	}

	if callCount != 2 {
		t.Errorf("Function called %d times, want 2", callCount)
	}
}

func TestRetryHandler_Retry_NonRetryableError(t *testing.T) {
	handler := NewDefaultRetryHandler()
	ctx := context.Background()

	callCount := 0
	expectedErr := ErrUnauthorized // Non-retryable error
	fn := func() error {
		callCount++
		return expectedErr
	}

	err := handler.Retry(ctx, "test_operation", fn)

	if err != expectedErr {
		t.Errorf("Retry() = %v, want %v", err, expectedErr)
	}

	if callCount != 1 {
		t.Errorf("Function called %d times, want 1", callCount)
	}
}

func TestRetryHandler_Retry_ContextCancellation(t *testing.T) {
	handler := NewRetryHandler(RetryInfo{
		MaxAttempts: 5,
		BaseDelay:   100 * time.Millisecond, // Long delay to ensure cancellation
		MaxDelay:    1 * time.Second,
		Multiplier:  2.0,
		Jitter:      false,
	})

	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	fn := func() error {
		callCount++
		if callCount == 1 {
			// Cancel context after first attempt
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()
		}
		return ErrTimeout // Retryable error
	}

	err := handler.Retry(ctx, "test_operation", fn)

	if err != context.Canceled {
		t.Errorf("Retry() = %v, want %v", err, context.Canceled)
	}

	// Should have been called at least once, but not all attempts due to cancellation
	if callCount == 0 {
		t.Errorf("Function called %d times, want at least 1", callCount)
	}
}

func TestRetryHandler_RetryWithResource(t *testing.T) {
	handler := NewRetryHandler(RetryInfo{
		MaxAttempts: 2,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    10 * time.Millisecond,
		Multiplier:  2.0,
		Jitter:      false,
	})
	ctx := context.Background()

	callCount := 0
	fn := func() error {
		callCount++
		return ErrTimeout // Retryable error
	}

	err := handler.RetryWithResource(ctx, "push", "registry.example.com/repo:tag", fn)

	// Should return wrapped OperationError
	if opErr, ok := err.(*OperationError); !ok {
		t.Errorf("RetryWithResource() returned %T, want *OperationError", err)
	} else {
		if opErr.Operation != "push" {
			t.Errorf("OperationError.Operation = %v, want push", opErr.Operation)
		}
		if opErr.Resource != "registry.example.com/repo:tag" {
			t.Errorf("OperationError.Resource = %v, want registry.example.com/repo:tag", opErr.Resource)
		}
		if opErr.Attempt != 2 {
			t.Errorf("OperationError.Attempt = %v, want 2", opErr.Attempt)
		}
	}

	if callCount != 2 {
		t.Errorf("Function called %d times, want 2", callCount)
	}
}

func TestRetryHandler_calculateDelay(t *testing.T) {
	tests := []struct {
		name     string
		config   RetryInfo
		attempt  int
		minDelay time.Duration
		maxDelay time.Duration
	}{
		{
			name: "first retry",
			config: RetryInfo{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 2.0,
				MaxDelay:   10 * time.Second,
				Jitter:     false,
			},
			attempt:  1,
			minDelay: 100 * time.Millisecond,
			maxDelay: 100 * time.Millisecond,
		},
		{
			name: "second retry",
			config: RetryInfo{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 2.0,
				MaxDelay:   10 * time.Second,
				Jitter:     false,
			},
			attempt:  2,
			minDelay: 200 * time.Millisecond,
			maxDelay: 200 * time.Millisecond,
		},
		{
			name: "max delay reached",
			config: RetryInfo{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 2.0,
				MaxDelay:   300 * time.Millisecond,
				Jitter:     false,
			},
			attempt:  5, // Would be 1600ms, but capped at 300ms
			minDelay: 300 * time.Millisecond,
			maxDelay: 300 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewRetryHandler(tt.config)
			delay := handler.calculateDelay(tt.attempt)

			if delay < tt.minDelay || delay > tt.maxDelay {
				t.Errorf("calculateDelay(%d) = %v, want between %v and %v",
					tt.attempt, delay, tt.minDelay, tt.maxDelay)
			}
		})
	}
}

func TestRetryHandler_ShouldRetry(t *testing.T) {
	handler := NewRetryHandler(RetryInfo{MaxAttempts: 3})

	tests := []struct {
		name     string
		err      error
		attempt  int
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			attempt:  1,
			expected: false,
		},
		{
			name:     "retryable error within attempts",
			err:      ErrTimeout,
			attempt:  2,
			expected: true,
		},
		{
			name:     "retryable error at max attempts",
			err:      ErrTimeout,
			attempt:  3,
			expected: false,
		},
		{
			name:     "non-retryable error",
			err:      ErrUnauthorized,
			attempt:  1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler.ShouldRetry(tt.err, tt.attempt); got != tt.expected {
				t.Errorf("ShouldRetry() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRetryWithBackoff(t *testing.T) {
	ctx := context.Background()
	callCount := 0
	fn := func() error {
		callCount++
		if callCount < 2 {
			return ErrTimeout
		}
		return nil
	}

	err := RetryWithBackoff(ctx, "test", 3, 1*time.Millisecond, fn)

	if err != nil {
		t.Errorf("RetryWithBackoff() = %v, want nil", err)
	}

	if callCount != 2 {
		t.Errorf("Function called %d times, want 2", callCount)
	}
}

func TestRetryConfig(t *testing.T) {
	config := NewRetryConfig().
		WithMaxAttempts(5).
		WithBaseDelay(50 * time.Millisecond).
		WithMaxDelay(5 * time.Second).
		WithMultiplier(1.5).
		WithJitter(false).
		Build()

	expected := RetryInfo{
		MaxAttempts: 5,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  1.5,
		Jitter:      false,
	}

	if config != expected {
		t.Errorf("RetryConfig.Build() = %v, want %v", config, expected)
	}
}