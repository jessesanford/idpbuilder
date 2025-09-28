package errors

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestStructuredError(t *testing.T) {
	t.Run("NewStructuredError creates valid error", func(t *testing.T) {
		cause := fmt.Errorf("original error")
		err := NewStructuredError(BuildFailed, "test_op", "test message", cause)

		if err.Code != BuildFailed {
			t.Errorf("Expected code %s, got %s", BuildFailed, err.Code)
		}
		if err.Op != "test_op" {
			t.Errorf("Expected op 'test_op', got '%s'", err.Op)
		}
		if err.Message != "test message" {
			t.Errorf("Expected message 'test message', got '%s'", err.Message)
		}
		if err.Cause != cause {
			t.Errorf("Expected cause to be original error")
		}
		if err.Timestamp.IsZero() {
			t.Error("Expected timestamp to be set")
		}
	})

	t.Run("Error() formats message correctly with cause", func(t *testing.T) {
		cause := fmt.Errorf("original error")
		err := NewStructuredError(BuildFailed, "test_op", "test message", cause)
		expected := "[BUILD_FAILED] test_op: test message (caused by: original error)"
		if err.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("Error() formats message correctly without cause", func(t *testing.T) {
		err := NewStructuredError(BuildFailed, "test_op", "test message", nil)
		expected := "[BUILD_FAILED] test_op: test message"
		if err.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("Unwrap returns cause", func(t *testing.T) {
		cause := fmt.Errorf("original error")
		err := NewStructuredError(BuildFailed, "test_op", "test message", cause)
		if err.Unwrap() != cause {
			t.Error("Unwrap should return the cause")
		}
	})
}

func TestErrorCode(t *testing.T) {
	t.Run("IsRetryable returns correct values", func(t *testing.T) {
		retryableCodes := []ErrorCode{RegistryUnreachable, BuildTimeout}
		nonRetryableCodes := []ErrorCode{BuildFailed, RegistryAuthFailed, CertificateInvalid}

		for _, code := range retryableCodes {
			if !code.IsRetryable() {
				t.Errorf("Expected %s to be retryable", code)
			}
		}

		for _, code := range nonRetryableCodes {
			if code.IsRetryable() {
				t.Errorf("Expected %s to not be retryable", code)
			}
		}
	})

	t.Run("IsRecoverable returns correct values", func(t *testing.T) {
		recoverableCodes := []ErrorCode{ConfigurationError, ValidationFailed}
		nonRecoverableCodes := []ErrorCode{BuildFailed, RegistryAuthFailed, CertificateInvalid}

		for _, code := range recoverableCodes {
			if !code.IsRecoverable() {
				t.Errorf("Expected %s to be recoverable", code)
			}
		}

		for _, code := range nonRecoverableCodes {
			if code.IsRecoverable() {
				t.Errorf("Expected %s to not be recoverable", code)
			}
		}
	})
}

func TestExponentialBackoff(t *testing.T) {
	t.Run("NewExponentialBackoff creates valid strategy", func(t *testing.T) {
		base := 100 * time.Millisecond
		max := 5 * time.Second
		maxAttempts := 3

		strategy := NewExponentialBackoff(base, max, maxAttempts)

		if strategy.baseDelay != base {
			t.Errorf("Expected base delay %v, got %v", base, strategy.baseDelay)
		}
		if strategy.maxDelay != max {
			t.Errorf("Expected max delay %v, got %v", max, strategy.maxDelay)
		}
		if strategy.maxAttempts != maxAttempts {
			t.Errorf("Expected max attempts %d, got %d", maxAttempts, strategy.maxAttempts)
		}
	})

	t.Run("ShouldRetry works correctly", func(t *testing.T) {
		strategy := NewExponentialBackoff(100*time.Millisecond, 5*time.Second, 3)

		// Test retryable error
		retryableErr := NewStructuredError(RegistryUnreachable, "test", "test", nil)
		if !strategy.ShouldRetry(retryableErr, 0) {
			t.Error("Should retry RegistryUnreachable error on attempt 0")
		}
		if !strategy.ShouldRetry(retryableErr, 1) {
			t.Error("Should retry RegistryUnreachable error on attempt 1")
		}
		if strategy.ShouldRetry(retryableErr, 3) {
			t.Error("Should not retry RegistryUnreachable error on attempt 3 (max attempts)")
		}

		// Test non-retryable error
		nonRetryableErr := NewStructuredError(BuildFailed, "test", "test", nil)
		if strategy.ShouldRetry(nonRetryableErr, 0) {
			t.Error("Should not retry BuildFailed error")
		}
	})

	t.Run("NextDelay calculates exponential backoff", func(t *testing.T) {
		base := 100 * time.Millisecond
		max := 1 * time.Second
		strategy := NewExponentialBackoff(base, max, 5)

		delay0 := strategy.NextDelay(0)
		delay1 := strategy.NextDelay(1)
		delay2 := strategy.NextDelay(2)
		delay3 := strategy.NextDelay(3)

		if delay0 != base {
			t.Errorf("Expected delay0 to be %v, got %v", base, delay0)
		}
		if delay1 != 2*base {
			t.Errorf("Expected delay1 to be %v, got %v", 2*base, delay1)
		}
		if delay2 != 4*base {
			t.Errorf("Expected delay2 to be %v, got %v", 4*base, delay2)
		}
		if delay3 != 8*base {
			t.Errorf("Expected delay3 to be %v, got %v", 8*base, delay3)
		}

		// Test max delay cap
		delay10 := strategy.NextDelay(10)
		if delay10 != max {
			t.Errorf("Expected delay10 to be capped at %v, got %v", max, delay10)
		}
	})

	t.Run("NewExponentialBackoffFromEnv respects environment variables", func(t *testing.T) {
		// Set environment variables
		os.Setenv("RETRY_BASE_DELAY_MS", "200")
		os.Setenv("RETRY_MAX_DELAY_MS", "5000")
		os.Setenv("RETRY_MAX_ATTEMPTS", "5")
		defer func() {
			os.Unsetenv("RETRY_BASE_DELAY_MS")
			os.Unsetenv("RETRY_MAX_DELAY_MS")
			os.Unsetenv("RETRY_MAX_ATTEMPTS")
		}()

		strategy := NewExponentialBackoffFromEnv()

		if strategy.baseDelay != 200*time.Millisecond {
			t.Errorf("Expected base delay 200ms, got %v", strategy.baseDelay)
		}
		if strategy.maxDelay != 5*time.Second {
			t.Errorf("Expected max delay 5s, got %v", strategy.maxDelay)
		}
		if strategy.maxAttempts != 5 {
			t.Errorf("Expected max attempts 5, got %d", strategy.maxAttempts)
		}
	})
}

func TestDefaultRecovery(t *testing.T) {
	t.Run("NewDefaultRecovery creates valid handler", func(t *testing.T) {
		recovery := NewDefaultRecovery()
		if recovery.recoveryStrategies == nil {
			t.Error("Expected recovery strategies to be initialized")
		}
	})

	t.Run("CanRecover identifies recoverable errors", func(t *testing.T) {
		recovery := NewDefaultRecovery()

		// Test recoverable errors
		configErr := NewStructuredError(ConfigurationError, "test", "test", nil)
		if !recovery.CanRecover(configErr) {
			t.Error("Should be able to recover from ConfigurationError")
		}

		validationErr := NewStructuredError(ValidationFailed, "test", "test", nil)
		if !recovery.CanRecover(validationErr) {
			t.Error("Should be able to recover from ValidationFailed")
		}

		// Test non-recoverable error
		buildErr := NewStructuredError(BuildFailed, "test", "test", nil)
		if recovery.CanRecover(buildErr) {
			t.Error("Should not be able to recover from BuildFailed")
		}
	})

	t.Run("Recover handles recoverable errors", func(t *testing.T) {
		recovery := NewDefaultRecovery()

		configErr := NewStructuredError(ConfigurationError, "test", "test", nil)
		recovered := recovery.Recover(configErr)

		if recovered == configErr {
			t.Error("Expected recovered error to be different from original")
		}

		if structuredRecovered, ok := recovered.(*StructuredError); ok {
			if structuredRecovered.Code != ConfigurationError {
				t.Errorf("Expected recovered error to have same code")
			}
		} else {
			t.Error("Expected recovered error to be StructuredError")
		}
	})

	t.Run("RecoverWithContext handles context cancellation", func(t *testing.T) {
		recovery := NewDefaultRecovery()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		configErr := NewStructuredError(ConfigurationError, "test", "test", nil)
		recovered := recovery.RecoverWithContext(ctx, configErr)

		if structuredRecovered, ok := recovered.(*StructuredError); ok {
			if structuredRecovered.Code != InternalError {
				t.Error("Expected context cancellation to return InternalError")
			}
		} else {
			t.Error("Expected recovered error to be StructuredError")
		}
	})
}

func TestErrorHandlerImpl(t *testing.T) {
	t.Run("NewErrorHandler creates valid handler", func(t *testing.T) {
		handler := NewErrorHandler()
		if handler.retry == nil {
			t.Error("Expected retry strategy to be set")
		}
		if handler.recovery == nil {
			t.Error("Expected recovery handler to be set")
		}
	})

	t.Run("Handle processes nil error", func(t *testing.T) {
		handler := NewErrorHandler()
		result := handler.Handle(nil)
		if result != nil {
			t.Error("Expected nil error to return nil")
		}
	})

	t.Run("Handle processes structured error", func(t *testing.T) {
		handler := NewErrorHandler()
		originalErr := NewStructuredError(BuildFailed, "test", "test", nil)
		result := handler.Handle(originalErr)

		if result != originalErr {
			t.Error("Expected structured error to be returned as-is after processing")
		}
	})

	t.Run("Handle classifies generic error", func(t *testing.T) {
		handler := NewErrorHandler()
		genericErr := fmt.Errorf("build failed during compilation")
		result := handler.Handle(genericErr)

		if structured, ok := result.(*StructuredError); ok {
			if structured.Code != BuildFailed {
				t.Errorf("Expected BuildFailed, got %s", structured.Code)
			}
		} else {
			t.Error("Expected result to be StructuredError")
		}
	})

	t.Run("HandleWithContext handles context cancellation", func(t *testing.T) {
		handler := NewErrorHandler()
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		genericErr := fmt.Errorf("test error")
		result := handler.HandleWithContext(ctx, genericErr)

		if structured, ok := result.(*StructuredError); ok {
			if structured.Code != InternalError {
				t.Error("Expected context cancellation to return InternalError")
			}
		} else {
			t.Error("Expected result to be StructuredError")
		}
	})

	t.Run("WithRetry updates retry strategy", func(t *testing.T) {
		handler := NewErrorHandler()
		newStrategy := NewExponentialBackoff(50*time.Millisecond, 1*time.Second, 5)

		updatedHandler := handler.WithRetry(newStrategy)
		if updatedHandler != handler {
			t.Error("Expected WithRetry to return the same handler instance")
		}
		if handler.retry != newStrategy {
			t.Error("Expected retry strategy to be updated")
		}
	})

	t.Run("WithRecovery updates recovery handler", func(t *testing.T) {
		handler := NewErrorHandler()
		newRecovery := NewDefaultRecovery()

		updatedHandler := handler.WithRecovery(newRecovery)
		if updatedHandler != handler {
			t.Error("Expected WithRecovery to return the same handler instance")
		}
		if handler.recovery != newRecovery {
			t.Error("Expected recovery handler to be updated")
		}
	})

	t.Run("ExecuteWithRetry succeeds on first attempt", func(t *testing.T) {
		handler := NewErrorHandler()
		ctx := context.Background()

		callCount := 0
		operation := func() error {
			callCount++
			return nil // Success
		}

		err := handler.ExecuteWithRetry(ctx, operation)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if callCount != 1 {
			t.Errorf("Expected operation to be called once, got %d", callCount)
		}
	})

	t.Run("ExecuteWithRetry retries on retryable error", func(t *testing.T) {
		handler := NewErrorHandler()
		ctx := context.Background()

		callCount := 0
		operation := func() error {
			callCount++
			if callCount < 3 {
				return fmt.Errorf("registry unreachable error")
			}
			return nil // Success on third attempt
		}

		err := handler.ExecuteWithRetry(ctx, operation)
		if err != nil {
			t.Errorf("Expected no error after retries, got %v", err)
		}
		if callCount != 3 {
			t.Errorf("Expected operation to be called 3 times, got %d", callCount)
		}
	})

	t.Run("GetRetryableErrorCodes returns correct codes", func(t *testing.T) {
		handler := NewErrorHandler()
		codes := handler.GetRetryableErrorCodes()

		expectedCodes := []ErrorCode{RegistryUnreachable, BuildTimeout}
		if len(codes) != len(expectedCodes) {
			t.Errorf("Expected %d retryable codes, got %d", len(expectedCodes), len(codes))
		}

		codeMap := make(map[ErrorCode]bool)
		for _, code := range codes {
			codeMap[code] = true
		}

		for _, expected := range expectedCodes {
			if !codeMap[expected] {
				t.Errorf("Expected %s to be in retryable codes", expected)
			}
		}
	})

	t.Run("GetRecoverableErrorCodes returns correct codes", func(t *testing.T) {
		handler := NewErrorHandler()
		codes := handler.GetRecoverableErrorCodes()

		expectedCodes := []ErrorCode{ConfigurationError, ValidationFailed}
		if len(codes) != len(expectedCodes) {
			t.Errorf("Expected %d recoverable codes, got %d", len(expectedCodes), len(codes))
		}

		codeMap := make(map[ErrorCode]bool)
		for _, code := range codes {
			codeMap[code] = true
		}

		for _, expected := range expectedCodes {
			if !codeMap[expected] {
				t.Errorf("Expected %s to be in recoverable codes", expected)
			}
		}
	})
}

func TestErrorClassification(t *testing.T) {
	handler := NewErrorHandler()

	testCases := []struct {
		errorMsg     string
		expectedCode ErrorCode
	}{
		{"build failed during compilation", BuildFailed},
		{"build timeout exceeded", BuildTimeout},
		{"registry unreachable", RegistryUnreachable},
		{"registry auth failed", RegistryAuthFailed},
		{"registry push failed", RegistryPushFailed},
		{"certificate invalid", CertificateInvalid},
		{"certificate expired", CertificateExpired},
		{"validation error occurred", ValidationFailed},
		{"config file not found", ConfigurationError},
		{"unknown error", InternalError},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("classifies '%s' correctly", tc.errorMsg), func(t *testing.T) {
			err := fmt.Errorf(tc.errorMsg)
			result := handler.Handle(err)

			if structured, ok := result.(*StructuredError); ok {
				if structured.Code != tc.expectedCode {
					t.Errorf("Expected %s, got %s for error: %s", tc.expectedCode, structured.Code, tc.errorMsg)
				}
			} else {
				t.Error("Expected result to be StructuredError")
			}
		})
	}
}