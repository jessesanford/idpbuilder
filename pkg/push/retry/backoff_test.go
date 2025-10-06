package retry

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

func TestExponentialBackoff_NextDelay(t *testing.T) {
	backoff := &ExponentialBackoff{
		BaseDelay:      100 * time.Millisecond,
		MaxDelay:       10 * time.Second,
		Multiplier:     2.0,
		JitterFraction: 0, // No jitter for deterministic testing
		rng:            rand.New(rand.NewSource(1)),
	}

	tests := []struct {
		attempt      int
		expectedMin  time.Duration
		expectedMax  time.Duration
		description  string
	}{
		{0, 100 * time.Millisecond, 100 * time.Millisecond, "first attempt"},
		{1, 200 * time.Millisecond, 200 * time.Millisecond, "second attempt"},
		{2, 400 * time.Millisecond, 400 * time.Millisecond, "third attempt"},
		{3, 800 * time.Millisecond, 800 * time.Millisecond, "fourth attempt"},
		{10, 10 * time.Second, 10 * time.Second, "capped at max"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			delay := backoff.NextDelay(tt.attempt)
			if delay < tt.expectedMin || delay > tt.expectedMax {
				t.Errorf("attempt %d: expected delay between %v and %v, got %v",
					tt.attempt, tt.expectedMin, tt.expectedMax, delay)
			}
		})
	}
}

func TestExponentialBackoff_WithJitter(t *testing.T) {
	backoff := NewExponentialBackoff()
	backoff.BaseDelay = 1 * time.Second
	backoff.MaxDelay = 30 * time.Second
	backoff.Multiplier = 2.0
	backoff.JitterFraction = 0.1

	// Test that jitter adds variability
	attempt := 2
	delays := make([]time.Duration, 10)
	for i := 0; i < 10; i++ {
		delays[i] = backoff.NextDelay(attempt)
	}

	// Check that not all delays are identical (jitter is working)
	allSame := true
	for i := 1; i < len(delays); i++ {
		if delays[i] != delays[0] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("expected jitter to produce different delays, but all were the same")
	}

	// Check that all delays are within expected bounds
	expectedBase := 4 * time.Second // 1s * 2^2
	maxWithJitter := expectedBase + time.Duration(float64(expectedBase)*0.1)
	for i, delay := range delays {
		if delay < expectedBase {
			t.Errorf("delay %d (%v) is less than base delay %v", i, delay, expectedBase)
		}
		if delay > maxWithJitter {
			t.Errorf("delay %d (%v) exceeds max with jitter %v", i, delay, maxWithJitter)
		}
	}
}

func TestExponentialBackoff_MaxDelayEnforced(t *testing.T) {
	backoff := &ExponentialBackoff{
		BaseDelay:      1 * time.Second,
		MaxDelay:       5 * time.Second,
		Multiplier:     2.0,
		JitterFraction: 0,
		rng:            rand.New(rand.NewSource(1)),
	}

	// After enough attempts, delay should be capped at MaxDelay
	for attempt := 0; attempt < 10; attempt++ {
		delay := backoff.NextDelay(attempt)
		if delay > backoff.MaxDelay {
			t.Errorf("attempt %d: delay %v exceeds max delay %v", attempt, delay, backoff.MaxDelay)
		}
	}
}

func TestExponentialBackoff_NegativeAttempt(t *testing.T) {
	backoff := NewExponentialBackoff()

	// Negative attempts should be treated as 0
	delay := backoff.NextDelay(-1)
	expected := backoff.BaseDelay
	if delay < expected {
		t.Errorf("expected delay >= %v for negative attempt, got %v", expected, delay)
	}
}

func TestWait_CompletesSuccessfully(t *testing.T) {
	ctx := context.Background()
	delay := 10 * time.Millisecond

	start := time.Now()
	err := Wait(ctx, delay)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// Allow some tolerance for timer precision
	if elapsed < delay || elapsed > delay+50*time.Millisecond {
		t.Errorf("expected elapsed time ~%v, got %v", delay, elapsed)
	}
}

func TestWait_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	delay := 1 * time.Second

	// Cancel context immediately
	cancel()

	start := time.Now()
	err := Wait(ctx, delay)
	elapsed := time.Since(start)

	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}

	// Should return quickly, not wait full delay
	if elapsed > 100*time.Millisecond {
		t.Errorf("expected quick return on cancel, but waited %v", elapsed)
	}
}

func TestWait_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	delay := 1 * time.Second

	start := time.Now()
	err := Wait(ctx, delay)
	elapsed := time.Since(start)

	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded error, got: %v", err)
	}

	// Should timeout around 50ms, not wait full 1s
	if elapsed > 150*time.Millisecond {
		t.Errorf("expected timeout around 50ms, but waited %v", elapsed)
	}
}

func TestConstantBackoff_NextDelay(t *testing.T) {
	delay := 500 * time.Millisecond
	backoff := NewConstantBackoff(delay)

	// All attempts should return the same delay
	for attempt := 0; attempt < 10; attempt++ {
		actual := backoff.NextDelay(attempt)
		if actual != delay {
			t.Errorf("attempt %d: expected constant delay %v, got %v", attempt, delay, actual)
		}
	}
}

func TestConstantBackoff_Reset(t *testing.T) {
	backoff := NewConstantBackoff(100 * time.Millisecond)

	// Reset should not panic or change behavior
	backoff.Reset()

	delay := backoff.NextDelay(5)
	expected := 100 * time.Millisecond
	if delay != expected {
		t.Errorf("expected %v after reset, got %v", expected, delay)
	}
}

func TestNewExponentialBackoff_Defaults(t *testing.T) {
	backoff := NewExponentialBackoff()

	if backoff.BaseDelay != 100*time.Millisecond {
		t.Errorf("expected default BaseDelay 100ms, got %v", backoff.BaseDelay)
	}
	if backoff.MaxDelay != 30*time.Second {
		t.Errorf("expected default MaxDelay 30s, got %v", backoff.MaxDelay)
	}
	if backoff.Multiplier != 2.0 {
		t.Errorf("expected default Multiplier 2.0, got %v", backoff.Multiplier)
	}
	if backoff.JitterFraction != 0.1 {
		t.Errorf("expected default JitterFraction 0.1, got %v", backoff.JitterFraction)
	}
	if backoff.rng == nil {
		t.Error("expected rng to be initialized")
	}
}

func TestNewExponentialBackoffWithConfig(t *testing.T) {
	baseDelay := 200 * time.Millisecond
	maxDelay := 5 * time.Second
	multiplier := 3.0
	jitter := 0.2

	backoff := NewExponentialBackoffWithConfig(baseDelay, maxDelay, multiplier, jitter)

	if backoff.BaseDelay != baseDelay {
		t.Errorf("expected BaseDelay %v, got %v", baseDelay, backoff.BaseDelay)
	}
	if backoff.MaxDelay != maxDelay {
		t.Errorf("expected MaxDelay %v, got %v", maxDelay, backoff.MaxDelay)
	}
	if backoff.Multiplier != multiplier {
		t.Errorf("expected Multiplier %v, got %v", multiplier, backoff.Multiplier)
	}
	if backoff.JitterFraction != jitter {
		t.Errorf("expected JitterFraction %v, got %v", jitter, backoff.JitterFraction)
	}
}